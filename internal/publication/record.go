package publication

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const RecordSchemaVersion = 1

// Transition is one row of the durable record's audit trail: one action taken
// (or observed) against one package during a publication run.
type Transition struct {
	Package           string `json:"package"`
	Version           string `json:"version"`
	Action            string `json:"action"`    // stage|publish|verify|promote|deprecate|tag-remove
	AuthMode          string `json:"auth_mode"` // oidc-stage|direct|approval
	StageID           string `json:"stage_id,omitempty"`
	LocalIntegrity    string `json:"local_integrity,omitempty"`
	RegistryIntegrity string `json:"registry_integrity,omitempty"`
	TagState          string `json:"tag_state,omitempty"`
	Result            string `json:"result"`
	Timestamp         string `json:"timestamp"`
}

// Provenance binds one package name to the sha256 of the approved local
// tarball the run intended to publish — recorded once up front from
// dist/preflight/SHA256SUMS, independent of whatever the registry later reports.
// Kind ("wrapper"|"platform") is release-plan.mjs's own kind, carried through
// unchanged from ApprovedPackage — the one source every wrapper-vs-platform
// decision downstream of the record (next_action derivation included) reads,
// instead of re-deriving wrapper identity from the package name.
type Provenance struct {
	Package string `json:"package"`
	SHA256  string `json:"sha256"`
	Kind    string `json:"kind"`
}

// Record is the resumable state for one publication run: dist/publication/
// publication-record.json. It references the immutable release-index digest
// and never rewrites dist/preflight/ or package bytes.
type Record struct {
	SchemaVersion      int          `json:"schema_version"`
	ReleaseIndexSHA256 string       `json:"release_index_sha256"`
	Path               string       `json:"path"`    // first|public
	Profile            string       `json:"profile"` // public|bank
	Transitions        []Transition `json:"transitions"`
	Provenance         []Provenance `json:"provenance"`
	Result             string       `json:"result"`
}

// RecordPath is the one source for where the durable record lives.
func RecordPath(root string) string {
	return filepath.Join(root, "dist", "publication", "publication-record.json")
}

// LoadRecord reads the durable record if present. A missing file returns a
// zero Record and no error — the caller treats that as "no prior run".
func LoadRecord(root string) (Record, error) {
	data, err := os.ReadFile(RecordPath(root))
	if err != nil {
		if os.IsNotExist(err) {
			return Record{}, nil
		}
		return Record{}, fmt.Errorf("publication record is unreadable: %w", err)
	}
	var record Record
	if err := json.Unmarshal(data, &record); err != nil {
		return Record{}, fmt.Errorf("publication record is malformed: %w", err)
	}
	return record, nil
}

// SaveRecord writes record atomically (temp file + rename) so a crash mid-write
// never leaves a torn record. It only ever touches dist/publication/ — never
// dist/preflight/ (the immutable release index and SHA256SUMS) or package bytes.
func SaveRecord(root string, record Record) error {
	record.SchemaVersion = RecordSchemaVersion
	path := RecordPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("could not create publication record directory: %w", err)
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode publication record: %w", err)
	}
	data = append(data, '\n')
	temp := path + ".tmp"
	if err := os.WriteFile(temp, data, 0o644); err != nil {
		return fmt.Errorf("could not write publication record: %w", err)
	}
	if err := os.Rename(temp, path); err != nil {
		return fmt.Errorf("could not commit publication record: %w", err)
	}
	return nil
}

// releaseLockPath is the one lock file every serialized release operation
// (submit, promote, rollback) contends for, alongside the durable record it
// protects — so two simultaneous invocations against the same release root
// can never race a load-modify-SaveRecord cycle against each other.
func releaseLockPath(root string) string {
	return filepath.Join(root, "dist", "publication", "release.lock")
}

// AcquireReleaseLock exclusively locks root for one release operation
// (O_CREATE|O_EXCL: create-or-fail, never blocking) and returns a release func
// the caller must defer. A lock already held by another process is reported
// as an attributed error naming the lock path — never a bare "file exists" —
// so a concurrent `bench release submit`/`promote`/`rollback` fails loudly
// instead of silently racing the last SaveRecord to win.
func AcquireReleaseLock(root string) (release func(), err error) {
	path := releaseLockPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("could not create release lock directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("another release operation is already running: lock held at %s; remove it only once you are certain no other bench release process is active", path)
		}
		return nil, fmt.Errorf("could not acquire release lock at %s: %w", path, err)
	}
	_ = file.Close()
	return func() { _ = os.Remove(path) }, nil
}
