package manager

import (
	"context"
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spiffe/spire/pkg/common/telemetry"
	"github.com/spiffe/spire/pkg/common/x509util"
	"github.com/spiffe/spire/pkg/server/ca"
	"github.com/spiffe/spire/pkg/server/catalog"
	"github.com/spiffe/spire/pkg/server/datastore"
	"github.com/spiffe/spire/proto/private/server/journal"
	"google.golang.org/protobuf/proto"
)

const (
	// journalCap is the maximum number of entries per type that we'll
	// hold onto.
	journalCap = 10
)

type journalConfig struct {
	cat catalog.Catalog
	log logrus.FieldLogger
}

// Journal stores X509 CAs and JWT keys on disk as they are rotated by the
// manager. The data format is a PEM encoded protocol buffer.
type Journal struct {
	config *journalConfig

	mu                    sync.RWMutex
	activeX509AuthorityID string
	caJournalID           uint
	entries               *journal.Entries
}

func LoadJournal(ctx context.Context, config *journalConfig) (*Journal, error) {
	// Look for the CA journal of this server in the datastore.
	j, err := loadJournalFromDS(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to load journal from datastore: %w", err)
	}
	return j, nil
}

func (j *Journal) getEntries() *journal.Entries {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return proto.Clone(j.entries).(*journal.Entries)
}

func (j *Journal) AppendX509CA(ctx context.Context, slotID string, issuedAt time.Time, x509CA *ca.X509CA) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	backup := j.entries.X509CAs
	j.entries.X509CAs = append(j.entries.X509CAs, &journal.X509CAEntry{
		SlotId:              slotID,
		IssuedAt:            issuedAt.Unix(),
		NotAfter:            x509CA.Certificate.NotAfter.Unix(),
		Certificate:         x509CA.Certificate.Raw,
		UpstreamChain:       chainDER(x509CA.UpstreamChain),
		Status:              journal.Status_PREPARED,
		AuthorityId:         x509util.SubjectKeyIDToString(x509CA.Certificate.SubjectKeyId),
		UpstreamAuthorityId: x509util.SubjectKeyIDToString(x509CA.Certificate.AuthorityKeyId),
	})

	exceeded := len(j.entries.X509CAs) - journalCap
	if exceeded > 0 {
		// make a new slice so we keep growing the backing array to drop the first
		x509CAs := make([]*journal.X509CAEntry, journalCap)
		copy(x509CAs, j.entries.X509CAs[exceeded:])
		j.entries.X509CAs = x509CAs
	}

	if err := j.save(ctx); err != nil {
		j.entries.X509CAs = backup
		return err
	}

	return nil
}

// UpdateX509CAStatus updates a stored X509CA entry to have the given status,
// updating the CA journal.
func (j *Journal) UpdateX509CAStatus(ctx context.Context, authorityID string, status journal.Status) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	backup := j.entries.X509CAs
	var found bool
	for i := len(j.entries.X509CAs) - 1; i >= 0; i-- {
		entry := j.entries.X509CAs[i]
		if authorityID == entry.AuthorityId {
			found = true
			entry.Status = status
			if status == journal.Status_ACTIVE {
				j.activeX509AuthorityID = entry.AuthorityId
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("no journal entry found with authority ID %q", authorityID)
	}

	if err := j.save(ctx); err != nil {
		j.entries.X509CAs = backup
		return err
	}

	return nil
}

func (j *Journal) AppendJWTKey(ctx context.Context, slotID string, issuedAt time.Time, jwtKey *ca.JWTKey) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	pkixBytes, err := x509.MarshalPKIXPublicKey(jwtKey.Signer.Public())
	if err != nil {
		return err
	}

	backup := j.entries.JwtKeys
	j.entries.JwtKeys = append(j.entries.JwtKeys, &journal.JWTKeyEntry{
		SlotId:      slotID,
		IssuedAt:    issuedAt.Unix(),
		Kid:         jwtKey.Kid,
		PublicKey:   pkixBytes,
		NotAfter:    jwtKey.NotAfter.Unix(),
		Status:      journal.Status_PREPARED,
		AuthorityId: jwtKey.Kid,
	})

	exceeded := len(j.entries.JwtKeys) - journalCap
	if exceeded > 0 {
		// make a new slice so we keep growing the backing array to drop the first
		jwtKeys := make([]*journal.JWTKeyEntry, journalCap)
		copy(jwtKeys, j.entries.JwtKeys[exceeded:])
		j.entries.JwtKeys = jwtKeys
	}

	if err := j.save(ctx); err != nil {
		j.entries.JwtKeys = backup
		return err
	}

	return nil
}

// UpdateJWTKeyStatus updates a stored JWTKey entry to have the given status,
// updating the CA journal.
func (j *Journal) UpdateJWTKeyStatus(ctx context.Context, authorityID string, status journal.Status) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	backup := j.entries.JwtKeys

	var found bool
	for i := len(j.entries.JwtKeys) - 1; i >= 0; i-- {
		entry := j.entries.JwtKeys[i]
		if authorityID == entry.AuthorityId {
			found = true
			entry.Status = status
			break
		}
	}

	if !found {
		return fmt.Errorf("no journal entry found with authority ID %q", authorityID)
	}

	if err := j.save(ctx); err != nil {
		j.entries.JwtKeys = backup
		return err
	}

	return nil
}

func (j *Journal) setEntries(entries *journal.Entries) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.entries = entries
}

// saveInDatastore saves the provided marshaled entries in the datastore.
// If caJournalID has not been defined yet (it's value is 0), it first finds
// the CA journal records that corresponds to this server. In case that there is
// no CA record for this server, it creates one.
// The ID of the CA journal record that was saved is returned, in addition to
// the error (if any) of the operation.
func (j *Journal) saveInDatastore(ctx context.Context, entriesBytes []byte) (caJournalID uint, err error) {
	// Check if we already identified what's the CA journal for this server in
	// the datastore. If not, log that we are creating a new CA journal entry.
	if j.caJournalID == 0 {
		if j.activeX509AuthorityID == "" {
			j.config.log.Debug("There is no active X.509 authority yet. Can't save CA journal in the datastore")
			return 0, nil
		}
		j.config.log.Info("Creating a new CA journal entry")
	}

	ds := j.config.cat.GetDataStore()
	caJournal, err := ds.SetCAJournal(ctx, &datastore.CAJournal{
		ID:                    j.caJournalID,
		Data:                  entriesBytes,
		ActiveX509AuthorityID: j.activeX509AuthorityID,
	})
	if err != nil {
		return 0, err
	}

	j.config.log.WithFields(logrus.Fields{
		telemetry.CAJournalID:      caJournal.ID,
		telemetry.LocalAuthorityID: j.activeX509AuthorityID,
	}).Debug("Successfully stored CA journal entry in datastore")

	return caJournal.ID, nil
}

// findCAJournal finds the corresponding CA journal record in the datastore for
// this server. It does that by retrieving all the public keys managed by the
// KeyManager and trying to get a match with a record which last active
// X509 authority ID correspond to one of the keys.
func (j *Journal) findCAJournal(ctx context.Context) (*datastore.CAJournal, error) {
	km := j.config.cat.GetKeyManager()
	ds := j.config.cat.GetDataStore()

	// Get all the public keys managed by the KeyManager.
	kmKeys, err := km.GetKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from key manager: %w", err)
	}

	for _, k := range kmKeys {
		subjectKeyID, err := x509util.GetSubjectKeyID(k.Public())
		if err != nil {
			return nil, fmt.Errorf("failed to calculate the subject key identifier for public key with ID %q", k.ID())
		}

		authorityID := x509util.SubjectKeyIDToString(subjectKeyID)
		caJournal, err := ds.FetchCAJournal(ctx, authorityID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch CA journal from datastore: %w", err)
		}
		if caJournal != nil {
			// There is a CA journal record that has an active X509 authority
			// ID that matches with one of the public keys of this server. This
			// means that this record belongs to this server.
			j.config.log.WithFields(logrus.Fields{
				telemetry.CAJournalID:      caJournal.ID,
				telemetry.LocalAuthorityID: authorityID,
			}).Debug("Found a CA journal record that matches with a local X509 authority ID")

			return caJournal, nil
		}
	}

	return nil, nil
}

// save saves the CA journal both on disk and in the datastore.
// TODO: stop saving the CA journal on disk in v1.10.
func (j *Journal) save(ctx context.Context) error {
	entriesBytes, err := proto.Marshal(j.entries)
	if err != nil {
		return err
	}

	caJournalID, err := j.saveInDatastore(ctx, entriesBytes)
	if err != nil {
		return fmt.Errorf("could not save CA journal in the datastore: %w", err)
	}
	j.caJournalID = caJournalID

	return nil
}

func chainDER(chain []*x509.Certificate) [][]byte {
	var der [][]byte
	for _, cert := range chain {
		der = append(der, cert.Raw)
	}
	return der
}

// loadJournalFromDS loads the CA journal from the datastore.
// It does that by looking for a CA journal record that matches with one of the
// public keys of this server.
func loadJournalFromDS(ctx context.Context, config *journalConfig) (*Journal, error) {
	config.log.Debug("Loading journal from datastore")

	j := &Journal{
		config:  config,
		entries: new(journal.Entries),
	}

	caJournal, err := j.findCAJournal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find CA journal record: %w", err)
	}
	if caJournal == nil {
		j.config.log.Info("There is not a CA journal record that matches any of the local X509 authority IDs")
		return j, nil
	}

	j.caJournalID = caJournal.ID
	if err := proto.Unmarshal(caJournal.Data, j.entries); err != nil {
		return nil, fmt.Errorf("unable to unmarshal entries from CA journal record: %w", err)
	}
	return j, nil
}
