package protection

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"unicode"

	"github.com/gibbonmi/bench/internal/jsonfile"
)

// Env is the reserved environment entry that carries the descriptor vector to a nested CLI.
const Env = "BENCH_PROCESS_RESOURCES"

// Kind is one resource kind of schema 1.
type Kind string

const (
	KindCache           Kind = "cache"
	KindProspective     Kind = "prospective"
	KindBinary          Kind = "binary"
	KindKitTest         Kind = "kit-test"
	KindWorktree        Kind = "worktree"
	KindReleaseEvidence Kind = "release-evidence"
)

var kinds = []Kind{KindCache, KindProspective, KindBinary, KindKitTest, KindWorktree, KindReleaseEvidence}

// Status is one user fact: process, streams, or cleanup.
type Status string

const (
	Pending    Status = "pending"
	Complete   Status = "complete"
	Incomplete Status = "incomplete"
)

// The obligations that protection itself records. A lifetime adds its own unresolved names.
const (
	ObligationIdentity = "identity"
	ObligationProcess  = "process"
	ObligationGroup    = "group"
)

const (
	schema       = 1
	recordMode   = 0o600
	dirMode      = 0o700
	resourceFile = "resource.json"
	retiredFile  = "retired.json"
	lockFile     = "lock"
	usersDir     = "users"
	userFile     = "user.json"

	admissionOpen   = "open"
	admissionSealed = "sealed"

	identityPending = "pending"
	identityKnown   = "known"
	identityNoStart = "no-start"
)

// Descriptor names one resource generation and its store. It grants registration authority
// only: no signal, deletion, or executable authority.
type Descriptor struct {
	Store        string `json:"store"`
	ResourceID   string `json:"resource_id"`
	Kind         Kind   `json:"kind"`
	ResourcePath string `json:"resource_path"`
	Binding      string `json:"binding"`
}

// resourceRecord is resource.json, the generation of one protected resource.
type resourceRecord struct {
	Schema       int    `json:"schema"`
	ResourceID   string `json:"resource_id"`
	Kind         Kind   `json:"kind"`
	ResourcePath string `json:"resource_path"`
	Binding      string `json:"binding"`
	OwnerPID     int    `json:"owner_pid"`
	Admission    string `json:"admission"`
	UsersDir     string `json:"users_dir"`
}

// retiredRecord is retired.json, the resource of the store's last retired generation. It
// outlives that generation, so a reader can still find the resource of a store with no
// generation.
type retiredRecord struct {
	Schema       int    `json:"schema"`
	Kind         Kind   `json:"kind"`
	ResourcePath string `json:"resource_path"`
	Binding      string `json:"binding"`
}

// userRecord is one launch's user.json in one resource store.
type userRecord struct {
	Schema      int      `json:"schema"`
	ResourceID  string   `json:"resource_id"`
	UserID      string   `json:"user_id"`
	ActorPID    int      `json:"actor_pid"`
	Identity    string   `json:"identity"`
	PID         int      `json:"pid,omitempty"`
	PGID        int      `json:"pgid,omitempty"`
	Process     Status   `json:"process"`
	Streams     Status   `json:"streams"`
	Cleanup     Status   `json:"cleanup"`
	Obligations []string `json:"obligations"`
}

// ParseVector strictly decodes a descriptor vector. An empty value names no resource.
func ParseVector(value string) ([]Descriptor, error) {
	if value == "" {
		return nil, nil
	}
	var vector []Descriptor
	if err := jsonfile.DecodeExactDocument([]byte(value), &vector); err != nil {
		return nil, fmt.Errorf("%s: %w", Env, err)
	}
	if vector == nil {
		return nil, fmt.Errorf("%s: the vector is null", Env)
	}
	for _, d := range vector {
		if err := d.check(); err != nil {
			return nil, fmt.Errorf("%s: %w", Env, err)
		}
	}
	return vector, nil
}

// Inherited parses the descriptor vector that this process inherited.
func Inherited() ([]Descriptor, error) { return ParseVector(os.Getenv(Env)) }

// InheritedFrom parses the descriptor vector that environ passes to a child: the last value
// that environ assigns to the vector entry.
func InheritedFrom(environ []string) ([]Descriptor, error) {
	value := ""
	for _, entry := range environ {
		if after, ok := strings.CutPrefix(entry, Env+"="); ok {
			value = after
		}
	}
	return ParseVector(value)
}

// Environ returns env with the vector entry replaced by vector, or removed when it is empty.
func Environ(env []string, vector []Descriptor) []string {
	out := slices.DeleteFunc(slices.Clone(env), func(entry string) bool { return strings.HasPrefix(entry, Env+"=") })
	if len(vector) == 0 {
		return out
	}
	// A descriptor holds strings alone, so its encoding cannot fail.
	data, _ := json.Marshal(vector)
	return append(out, Env+"="+string(data))
}

// check refuses a descriptor whose fields cannot name one generation's store.
func (d Descriptor) check() error {
	if err := d.checkOwner(); err != nil {
		return err
	}
	return checkID("resource", d.ResourceID)
}

// checkOwner checks every field that an owner supplies before its generation exists.
func (d Descriptor) checkOwner() error {
	for _, path := range []string{d.Store, d.ResourcePath} {
		if err := checkPath(path); err != nil {
			return err
		}
	}
	if !slices.Contains(kinds, d.Kind) {
		return fmt.Errorf("descriptor kind %q is not a schema %d kind", d.Kind, schema)
	}
	if d.Binding == "" || strings.ContainsFunc(d.Binding, unicode.IsControl) {
		return fmt.Errorf("descriptor binding %q is empty or holds a control byte", d.Binding)
	}
	return nil
}

func checkPath(path string) error {
	if strings.ContainsFunc(path, unicode.IsControl) {
		return fmt.Errorf("descriptor path %q holds a control byte", path)
	}
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("descriptor path %q is not clean and absolute", path)
	}
	return nil
}

func checkID(name, id string) error {
	if raw, err := hex.DecodeString(id); err != nil || len(raw) != 16 || hex.EncodeToString(raw) != id {
		return fmt.Errorf("%s identifier %q is not 32 lowercase hex digits", name, id)
	}
	return nil
}

func newID() string {
	raw := make([]byte, 16)
	_, _ = rand.Read(raw)
	return hex.EncodeToString(raw)
}

// valid checks every schema value of a resource record.
func (r resourceRecord) valid() error {
	d := Descriptor{Store: "/", ResourceID: r.ResourceID, Kind: r.Kind, ResourcePath: r.ResourcePath, Binding: r.Binding}
	switch {
	case r.Schema != schema:
		return fmt.Errorf("resource schema %d is unsupported", r.Schema)
	case r.OwnerPID < 1:
		return fmt.Errorf("resource owner PID %d is not positive", r.OwnerPID)
	case r.Admission != admissionOpen && r.Admission != admissionSealed:
		return fmt.Errorf("resource admission %q is unknown", r.Admission)
	case r.UsersDir != usersDir:
		return fmt.Errorf("resource users directory %q is not %q", r.UsersDir, usersDir)
	}
	return d.check()
}

// resource names the kind, the path, and the binding of the retired resource, with no store
// and no generation.
func (r retiredRecord) resource() Descriptor {
	return Descriptor{Kind: r.Kind, ResourcePath: r.ResourcePath, Binding: r.Binding}
}

// valid checks every schema value of a retired record.
func (r retiredRecord) valid() error {
	if r.Schema != schema {
		return fmt.Errorf("retired schema %d is unsupported", r.Schema)
	}
	d := r.resource()
	d.Store = "/"
	return d.checkOwner()
}

// readRetired strictly reads and checks the retired record of store.
func readRetired(store string, limit int64) (retiredRecord, error) {
	var r retiredRecord
	if err := readRecord(filepath.Join(store, retiredFile), limit, &r); err != nil {
		return r, err
	}
	return r, r.valid()
}

// matches refuses a descriptor that names another generation, kind, path, or binding.
func (r resourceRecord) matches(d Descriptor) error {
	if r.ResourceID != d.ResourceID || r.Kind != d.Kind || r.ResourcePath != d.ResourcePath || r.Binding != d.Binding {
		return fmt.Errorf("descriptor for %s does not match the published generation in %s", d.ResourcePath, d.Store)
	}
	return nil
}

// valid checks every schema value of a user record read from the directory named dir.
func (u userRecord) valid(dir string) error {
	known, workers := u.Identity == identityKnown, []Status{Pending, Complete, Incomplete}
	switch {
	case u.Schema != schema:
		return fmt.Errorf("user schema %d is unsupported", u.Schema)
	case u.UserID != dir:
		return fmt.Errorf("user record in %s claims user %q", dir, u.UserID)
	case u.ActorPID < 1:
		return fmt.Errorf("user actor PID %d is not positive", u.ActorPID)
	case u.Identity != identityPending && !known && u.Identity != identityNoStart:
		return fmt.Errorf("user identity %q is unknown", u.Identity)
	case known != (u.PID > 0 && u.PGID > 0) || !known && (u.PID != 0 || u.PGID != 0):
		return fmt.Errorf("user identity %q disagrees with PID %d and PGID %d", u.Identity, u.PID, u.PGID)
	case !slices.Contains([]Status{Pending, Complete}, u.Process):
		return fmt.Errorf("user process status %q is unknown", u.Process)
	case !slices.Contains(workers, u.Streams) || !slices.Contains(workers, u.Cleanup):
		return fmt.Errorf("user stream or cleanup status %q, %q is unknown", u.Streams, u.Cleanup)
	case u.Obligations == nil || slices.Contains(u.Obligations, ""):
		return errors.New("user obligations are absent or hold an empty name")
	case len(u.Obligations) == 0 && u.Identity != identityNoStart && !u.proved():
		return errors.New("user names no obligation without proved completion")
	}
	if err := checkID("resource", u.ResourceID); err != nil {
		return err
	}
	return checkID("user", u.UserID)
}

// proved reports that the live owner recorded complete process, stream, and cleanup facts.
func (u userRecord) proved() bool {
	return u.Process == Complete && u.Streams == Complete && u.Cleanup == Complete
}

// readRecord reads one private record without following a link or opening a special file.
func readRecord(path string, limit int64, value any) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != recordMode {
		return fmt.Errorf("%s is %v, want a regular file at %v", path, info.Mode(), os.FileMode(recordMode))
	}
	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	if opened, err := file.Stat(); err != nil || !os.SameFile(info, opened) {
		return fmt.Errorf("%s changed while it was opened", path)
	}
	if limit < 1 {
		return fmt.Errorf("%s has no read limit", path)
	}
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return err
	}
	if int64(len(data)) > limit {
		return fmt.Errorf("%s exceeds the %d-byte record limit", path, limit)
	}
	if err := jsonfile.DecodeExactDocument(data, value); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

// checkDir refuses a store directory that is a link, a file, or not private.
func checkDir(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode().Perm() != dirMode {
		return fmt.Errorf("%s is %v, want a directory at %v", path, info.Mode(), os.FileMode(dirMode))
	}
	return nil
}
