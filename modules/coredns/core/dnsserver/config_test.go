package dnsserver

import (
	"testing"

	"github.com/coredns/caddy"
)

func TestKeyForConfig(t *testing.T) {
	tests := []struct {
		name          string
		blockIndex    int
		blockKeyIndex int
		expected      string
	}{
		{"zero_indices", 0, 0, "0:0"},
		{"positive_indices", 1, 2, "1:2"},
		{"larger_indices", 10, 5, "10:5"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := keyForConfig(tc.blockIndex, tc.blockKeyIndex)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s for blockIndex %d and blockKeyIndex %d",
					tc.expected, result, tc.blockIndex, tc.blockKeyIndex)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	controller := caddy.NewTestController("dns", "")
	initialCtx := controller.Context()
	dnsCtx, ok := initialCtx.(*dnsContext)
	if !ok {
		t.Fatalf("controller.Context() did not return a *dnsContext, got %T", initialCtx)
	}
	if dnsCtx.keysToConfigs == nil {
		t.Fatal("dnsCtx.keysToConfigs is nil; it should have been initialized by newContext")
	}

	t.Run("returns and saves default config when config missing", func(t *testing.T) {
		controller.ServerBlockIndex = 0
		controller.ServerBlockKeyIndex = 0
		key := keyForConfig(controller.ServerBlockIndex, controller.ServerBlockKeyIndex)

		// Ensure config doesn't exist initially for this specific key
		delete(dnsCtx.keysToConfigs, key)

		cfg := GetConfig(controller)
		if cfg == nil {
			t.Fatal("GetConfig returned nil (should create and return a default)")
		}
		if len(cfg.ListenHosts) != 1 || cfg.ListenHosts[0] != "" {
			t.Errorf("Expected default ListenHosts [\"\"] for auto-created config, got %v", cfg.ListenHosts)
		}

		savedCfg, found := dnsCtx.keysToConfigs[key]
		if !found {
			t.Fatal("fallback did not save the default config into the context")
		}
		if savedCfg != cfg {
			t.Fatal("config is not the same instance as the one saved in the context")
		}
	})
}
