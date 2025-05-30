package external

import (
	"context"
	"testing"

	"github.com/coredns/coredns/plugin/kubernetes"
	"github.com/coredns/coredns/plugin/kubernetes/object"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	api "k8s.io/api/core/v1"
)

func TestExternal(t *testing.T) {
	k := kubernetes.New([]string{"cluster.local."})
	k.Namespaces = map[string]struct{}{"testns": {}}
	k.APIConn = &external{}

	e := New()
	e.Zones = []string{"example.com.", "in-addr.arpa."}
	e.headless = true
	e.Next = test.NextHandler(dns.RcodeSuccess, nil)
	e.externalFunc = k.External
	e.externalAddrFunc = externalAddress  // internal test function
	e.externalSerialFunc = externalSerial // internal test function

	ctx := context.TODO()
	for i, tc := range tests {
		r := tc.Msg()
		w := dnstest.NewRecorder(&test.ResponseWriter{})

		_, err := e.ServeDNS(ctx, w, r)
		if err != tc.Error {
			t.Errorf("Test %d expected no error, got %v", i, err)
			return
		}
		if tc.Error != nil {
			continue
		}

		resp := w.Msg

		if resp == nil {
			t.Fatalf("Test %d, got nil message and no error for %q", i, r.Question[0].Name)
		}
		if !resp.Authoritative {
			t.Error("Expected authoritative answer")
		}
		if err = test.SortAndCheck(resp, tc); err != nil {
			t.Errorf("Test %d: %v", i, err)
		}
	}
}

var tests = []test.Case{
	// PTR reverse lookup
	{
		Qname: "4.3.2.1.in-addr.arpa.", Qtype: dns.TypePTR, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.PTR("4.3.2.1.in-addr.arpa. 5 IN PTR svc1.testns.example.com."),
		},
	},
	// Bad PTR reverse lookup using existing service name
	{
		Qname: "svc1.testns.example.com.", Qtype: dns.TypePTR, Rcode: dns.RcodeSuccess,
		Ns: []dns.RR{
			test.SOA("example.com.	5	IN	SOA	ns1.dns.example.com. hostmaster.example.com. 1499347823 7200 1800 86400 5"),
		},
	},
	// Bad PTR reverse lookup using non-existing service name
	{
		Qname: "not-existing.testns.example.com.", Qtype: dns.TypePTR, Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.SOA("example.com.	5	IN	SOA	ns1.dns.example.com. hostmaster.example.com. 1499347823 7200 1800 86400 5"),
		},
	},
	// A Service
	{
		Qname: "svc1.testns.example.com.", Qtype: dns.TypeA, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.A("svc1.testns.example.com.	5	IN	A	1.2.3.4"),
		},
	},
	{
		Qname: "svc1.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{test.SRV("svc1.testns.example.com.	5	IN	SRV	0 100 80 svc1.testns.example.com.")},
		Extra:  []dns.RR{test.A("svc1.testns.example.com.  5       IN      A       1.2.3.4")},
	},
	// SRV Service Not udp/tcp
	{
		Qname: "*._not-udp-or-tcp.svc1.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.SOA("example.com.	5	IN	SOA	ns1.dns.example.com. hostmaster.example.com. 1499347823 7200 1800 86400 5"),
		},
	},
	// SRV Service
	{
		Qname: "_http._tcp.svc1.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("_http._tcp.svc1.testns.example.com.	5	IN	SRV	0 100 80 svc1.testns.example.com."),
		},
		Extra: []dns.RR{
			test.A("svc1.testns.example.com.	5	IN	A	1.2.3.4"),
		},
	},
	// AAAA Service (with an existing A record, but no AAAA record)
	{
		Qname: "svc1.testns.example.com.", Qtype: dns.TypeAAAA, Rcode: dns.RcodeSuccess,
		Ns: []dns.RR{
			test.SOA("example.com.	5	IN	SOA	ns1.dns.example.com. hostmaster.example.com. 1499347823 7200 1800 86400 5"),
		},
	},
	// AAAA Service (non-existing service)
	{
		Qname: "svc0.testns.example.com.", Qtype: dns.TypeAAAA, Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.SOA("example.com.	5	IN	SOA	ns1.dns.example.com. hostmaster.example.com. 1499347823 7200 1800 86400 5"),
		},
	},
	// A Service (non-existing service)
	{
		Qname: "svc0.testns.example.com.", Qtype: dns.TypeA, Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.SOA("example.com.	5	IN	SOA	ns1.dns.example.com. hostmaster.example.com. 1499347823 7200 1800 86400 5"),
		},
	},
	// A Service (non-existing namespace)
	{
		Qname: "svc0.svc-nons.example.com.", Qtype: dns.TypeA, Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.SOA("example.com.	5	IN	SOA	ns1.dns.example.com. hostmaster.example.com. 1499347823 7200 1800 86400 5"),
		},
	},
	// AAAA Service
	{
		Qname: "svc6.testns.example.com.", Qtype: dns.TypeAAAA, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.AAAA("svc6.testns.example.com.	5	IN	AAAA	1:2::5"),
		},
	},
	// SRV
	{
		Qname: "_http._tcp.svc6.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("_http._tcp.svc6.testns.example.com.	5	IN	SRV	0 100 80 svc6.testns.example.com."),
		},
		Extra: []dns.RR{
			test.AAAA("svc6.testns.example.com.	5	IN	AAAA	1:2::5"),
		},
	},
	// SRV
	{
		Qname: "svc6.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("svc6.testns.example.com.	5	IN	SRV	0 100 80 svc6.testns.example.com."),
		},
		Extra: []dns.RR{
			test.AAAA("svc6.testns.example.com.	5	IN	AAAA	1:2::5"),
		},
	},
	{
		Qname: "testns.example.com.", Qtype: dns.TypeA, Rcode: dns.RcodeSuccess,
		Ns: []dns.RR{
			test.SOA("example.com.	5	IN	SOA	ns1.dns.example.com. hostmaster.example.com. 1499347823 7200 1800 86400 5"),
		},
	},
	{
		Qname: "testns.example.com.", Qtype: dns.TypeSOA, Rcode: dns.RcodeSuccess,
		Ns: []dns.RR{
			test.SOA("example.com.	5	IN	SOA	ns1.dns.example.com. hostmaster.example.com. 1499347823 7200 1800 86400 5"),
		},
	},
	// svc11
	{
		Qname: "svc11.testns.example.com.", Qtype: dns.TypeA, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.A("svc11.testns.example.com.	5	IN	A	2.3.4.5"),
		},
	},
	{
		Qname: "_http._tcp.svc11.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("_http._tcp.svc11.testns.example.com.	5	IN	SRV	0 100 80 svc11.testns.example.com."),
		},
		Extra: []dns.RR{
			test.A("svc11.testns.example.com.	5	IN	A	2.3.4.5"),
		},
	},
	{
		Qname: "svc11.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("svc11.testns.example.com.	5	IN	SRV	0 100 80 svc11.testns.example.com."),
		},
		Extra: []dns.RR{
			test.A("svc11.testns.example.com.	5	IN	A	2.3.4.5"),
		},
	},
	// svc12
	{
		Qname: "svc12.testns.example.com.", Qtype: dns.TypeA, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.CNAME("svc12.testns.example.com.	5	IN	CNAME	dummy.hostname"),
		},
	},
	{
		Qname: "_http._tcp.svc12.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("_http._tcp.svc12.testns.example.com.	5	IN	SRV	0 100 80 dummy.hostname."),
		},
	},
	{
		Qname: "svc12.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("svc12.testns.example.com.	5	IN	SRV	0 100 80 dummy.hostname."),
		},
	},
	// headless service
	{
		Qname: "svc-headless.testns.example.com.", Qtype: dns.TypeA, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.A("svc-headless.testns.example.com.	5	IN	A	1.2.3.4"),
			test.A("svc-headless.testns.example.com.	5	IN	A	1.2.3.5"),
		},
	},
	{
		Qname: "svc-headless.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("svc-headless.testns.example.com.	5	IN	SRV	0	50	80	endpoint-svc-0.svc-headless.testns.example.com."),
			test.SRV("svc-headless.testns.example.com.	5	IN	SRV	0	50	80	endpoint-svc-1.svc-headless.testns.example.com."),
		},
		Extra: []dns.RR{
			test.A("endpoint-svc-0.svc-headless.testns.example.com.	5	IN	A	1.2.3.4"),
			test.A("endpoint-svc-1.svc-headless.testns.example.com.	5	IN	A	1.2.3.5"),
		},
	},
	{
		Qname: "_http._tcp.svc-headless.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("_http._tcp.svc-headless.testns.example.com.	5	IN	SRV	0	50	80	endpoint-svc-0.svc-headless.testns.example.com."),
			test.SRV("_http._tcp.svc-headless.testns.example.com.	5	IN	SRV	0	50	80	endpoint-svc-1.svc-headless.testns.example.com."),
		},
		Extra: []dns.RR{
			test.A("endpoint-svc-0.svc-headless.testns.example.com.	5	IN	A	1.2.3.4"),
			test.A("endpoint-svc-1.svc-headless.testns.example.com.	5	IN	A	1.2.3.5"),
		},
	},
	{
		Qname: "endpoint-svc-0.svc-headless.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("endpoint-svc-0.svc-headless.testns.example.com.		5	IN	SRV	0	100	80	endpoint-svc-0.svc-headless.testns.example.com."),
		},
		Extra: []dns.RR{
			test.A("endpoint-svc-0.svc-headless.testns.example.com.	5	IN	A	1.2.3.4"),
		},
	},
	{
		Qname: "endpoint-svc-1.svc-headless.testns.example.com.", Qtype: dns.TypeSRV, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.SRV("endpoint-svc-1.svc-headless.testns.example.com.		5	IN	SRV	0	100	80	endpoint-svc-1.svc-headless.testns.example.com."),
		},
		Extra: []dns.RR{
			test.A("endpoint-svc-1.svc-headless.testns.example.com.	5	IN	A	1.2.3.5"),
		},
	},
	{
		Qname: "endpoint-svc-0.svc-headless.testns.example.com.", Qtype: dns.TypeA, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.A("endpoint-svc-0.svc-headless.testns.example.com.	5	IN	A	1.2.3.4"),
		},
	},
	{
		Qname: "endpoint-svc-1.svc-headless.testns.example.com.", Qtype: dns.TypeA, Rcode: dns.RcodeSuccess,
		Answer: []dns.RR{
			test.A("endpoint-svc-1.svc-headless.testns.example.com.	5	IN	A	1.2.3.5"),
		},
	},
}

type external struct{}

func (external) HasSynced() bool                           { return true }
func (external) Run()                                      {}
func (external) Stop() error                               { return nil }
func (external) EpIndexReverse(string) []*object.Endpoints { return nil }
func (external) SvcIndexReverse(string) []*object.Service  { return nil }
func (external) Modified(kubernetes.ModifiedMode) int64    { return 0 }

func (external) SvcImportIndex(s string) []*object.ServiceImport                    { return nil }
func (external) ServiceImportList() []*object.ServiceImport                         { return nil }
func (external) McEpIndex(s string) []*object.MultiClusterEndpoints                 { return nil }
func (external) MultiClusterEndpointsList(s string) []*object.MultiClusterEndpoints { return nil }

func (external) EpIndex(s string) []*object.Endpoints {
	return epIndexExternal[s]
}

func (external) EndpointsList() []*object.Endpoints {
	var eps []*object.Endpoints
	for _, ep := range epIndexExternal {
		eps = append(eps, ep...)
	}
	return eps
}
func (external) GetNodeByName(ctx context.Context, name string) (*api.Node, error) { return nil, nil }
func (external) SvcIndex(s string) []*object.Service                               { return svcIndexExternal[s] }
func (external) PodIndex(string) []*object.Pod                                     { return nil }

func (external) SvcExtIndexReverse(ip string) (result []*object.Service) {
	for _, svcs := range svcIndexExternal {
		for _, svc := range svcs {
			for _, exIp := range svc.ExternalIPs {
				if exIp != ip {
					continue
				}
				result = append(result, svc)
			}
		}
	}
	return result
}

func (external) GetNamespaceByName(name string) (*object.Namespace, error) {
	return &object.Namespace{
		Name: name,
	}, nil
}

var epIndexExternal = map[string][]*object.Endpoints{
	"svc-headless.testns": {
		{
			Name:      "svc-headless",
			Namespace: "testns",
			Index:     "svc-headless.testns",
			Subsets: []object.EndpointSubset{
				{
					Ports: []object.EndpointPort{
						{
							Port:     80,
							Name:     "http",
							Protocol: "TCP",
						},
					},
					Addresses: []object.EndpointAddress{
						{
							IP:            "1.2.3.4",
							Hostname:      "endpoint-svc-0",
							NodeName:      "test-node",
							TargetRefName: "endpoint-svc-0",
						},
						{
							IP:            "1.2.3.5",
							Hostname:      "endpoint-svc-1",
							NodeName:      "test-node",
							TargetRefName: "endpoint-svc-1",
						},
					},
				},
			},
		},
	},
}

var svcIndexExternal = map[string][]*object.Service{
	"svc1.testns": {
		{
			Name:        "svc1",
			Namespace:   "testns",
			Type:        api.ServiceTypeClusterIP,
			ClusterIPs:  []string{"10.0.0.1"},
			ExternalIPs: []string{"1.2.3.4"},
			Ports:       []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}},
		},
	},
	"svc6.testns": {
		{
			Name:        "svc6",
			Namespace:   "testns",
			Type:        api.ServiceTypeClusterIP,
			ClusterIPs:  []string{"10.0.0.3"},
			ExternalIPs: []string{"1:2::5"},
			Ports:       []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}},
		},
	},
	"svc11.testns": {
		{
			Name:        "svc11",
			Namespace:   "testns",
			Type:        api.ServiceTypeLoadBalancer,
			ExternalIPs: []string{"2.3.4.5"},
			ClusterIPs:  []string{"10.0.0.3"},
			Ports:       []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}},
		},
	},
	"svc12.testns": {
		{
			Name:        "svc12",
			Namespace:   "testns",
			Type:        api.ServiceTypeLoadBalancer,
			ClusterIPs:  []string{"10.0.0.3"},
			ExternalIPs: []string{"dummy.hostname"},
			Ports:       []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}},
		},
	},
	"svc-headless.testns": {
		{
			Name:       "svc-headless",
			Namespace:  "testns",
			Type:       api.ServiceTypeClusterIP,
			ClusterIPs: []string{"None"},
			Ports:      []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}},
		},
	},
}

func (external) ServiceList() []*object.Service {
	var svcs []*object.Service
	for _, svc := range svcIndexExternal {
		svcs = append(svcs, svc...)
	}
	return svcs
}

func externalAddress(state request.Request, headless bool) []dns.RR {
	a := test.A("example.org. IN A 127.0.0.1")
	return []dns.RR{a}
}

func externalSerial(string) uint32 {
	return 1499347823
}
