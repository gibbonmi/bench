package chargeevidence

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/fs"
	"math"
	"os"
	"strings"
	"syscall"
	"time"
)

// Store refusal classes. A store operation reports exactly one class for its first failed
// predicate, the same way a strict pack read does.
const (
	RefuseAbsent   = "absent-artifact"
	RefuseUnsafe   = "unsafe-store"
	RefuseCapacity = "capacity"
	RefuseExisting = "corrupt-existing-artifact"
	RefusePublish  = "publication-failed"
	RefuseReplaced = "replaced-artifact"
	RefuseStorage  = "storage-failed"
)

// Lock and object names inside the store directory. A pack name derives only from its
// validated identity; a temporary pack name derives only from random bytes.
const (
	OperationLockName = "operation.lock"
	WriterLockName    = "writer.lock"
	PackSuffix        = ".pack"
	TempPrefix        = "tmp-"
	TempSuffix        = ".partial"
)

// StoreOptions are the store's local-substitutable dependencies.
type StoreOptions struct {
	// Fault runs before each named publication step. A returned error fails that step, so
	// a test injects deterministic filesystem failures without editing private state.
	Fault func(step string) error
	// Pause runs at each named store stage. The system owner coordinates process
	// interruption through it, and a test replaces store objects at a stage.
	Pause func(stage string)
}

// Publication steps and store stages, in the order a writer reaches them. A reader reaches
// only StageStoreInspected.
const (
	StepWriteTemp        = "write-temp"
	StepSyncTemp         = "sync-temp"
	StepVerifyTemp       = "verify-temp"
	StepLink             = "link"
	StageStoreInspected  = "store-inspected"
	StageWriterLock      = "writer-lock"
	StageCapacity        = "capacity"
	StageStaged          = "staged"
	StageVerifying       = "verifying"
	PauseEnvironment     = "BENCH_EVIDENCE_PAUSE"
	pauseEnvironmentMark = ":"
)

// Store is the repository-common evidence store beneath one verified Git common directory.
type Store struct {
	parent  string
	options StoreOptions
}

// OpenStore addresses the store beneath an already verified common directory. It touches
// no file; only preparation creates an absent store.
func OpenStore(commonDir string, options StoreOptions) *Store {
	return &Store{parent: commonDir, options: options}
}

// PauseFromEnvironment returns the stage pause the system owner requests through
// PauseEnvironment, as `<stage>:<marker path>`. At that stage the writer creates the marker
// file and waits while it exists, so the owner either removes it to resume the writer or
// ends the process. On resumption the writer creates `<marker path>.resumed`, so the owner
// can order a release after the writer has left its pause. An unset value pauses nothing.
func PauseFromEnvironment() func(string) {
	stage, marker, ok := strings.Cut(os.Getenv(PauseEnvironment), pauseEnvironmentMark)
	if !ok || stage == "" || marker == "" {
		return nil
	}
	return func(reached string) {
		if reached != stage {
			return
		}
		if err := os.WriteFile(marker, []byte(reached), 0o600); err != nil {
			return
		}
		for {
			if _, err := os.Lstat(marker); errors.Is(err, fs.ErrNotExist) {
				_ = os.WriteFile(marker+".resumed", nil, 0o600)
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (s *Store) fault(step string) error {
	if s.options.Fault == nil {
		return nil
	}
	return s.options.Fault(step)
}

func (s *Store) pause(stage string) {
	if s.options.Pause != nil {
		s.options.Pause(stage)
	}
}

// openDir opens the store directory without following a replaced or linked store. Only a
// writer passes create; a reader of an absent store refuses and changes nothing.
func (s *Store) openDir(create bool) (*os.Root, error) {
	parent, err := os.OpenRoot(s.parent)
	if err != nil {
		return nil, refuse(RefuseUnsafe, "common directory is not openable: %v", err)
	}
	defer parent.Close()
	info, err := parent.Lstat(StoreName)
	if errors.Is(err, fs.ErrNotExist) && create {
		if err := parent.Mkdir(StoreName, 0o700); err != nil && !errors.Is(err, fs.ErrExist) {
			return nil, refuse(RefuseStorage, "store directory is not creatable: %v", err)
		}
		info, err = parent.Lstat(StoreName)
	}
	if errors.Is(err, fs.ErrNotExist) {
		return nil, refuse(RefuseAbsent, "the evidence store holds no artifacts")
	}
	if err != nil {
		return nil, refuse(RefuseUnsafe, "store directory is not inspectable: %v", err)
	}
	if !info.IsDir() {
		return nil, refuse(RefuseUnsafe, "store path is %s, not a directory", kindOf(info))
	}
	s.pause(StageStoreInspected)
	dir, err := parent.OpenRoot(StoreName)
	if err != nil {
		return nil, refuse(RefuseUnsafe, "store directory is not openable: %v", err)
	}
	opened, err := dir.Stat(".")
	if err != nil || !os.SameFile(info, opened) {
		dir.Close()
		return nil, refuse(RefuseUnsafe, "store directory changed while it was opened")
	}
	return dir, nil
}

// CapacityError reports a store that cannot admit a candidate under the selected quota.
// Required is zero when the required total is not representable.
type CapacityError struct {
	Observed, Candidate, Quota, Required uint64
}

func (e *CapacityError) Error() string {
	return RefuseCapacity + ": the store cannot admit the candidate"
}

// Admit is the store's capacity rule: a store holding observed bytes admits a candidate of
// candidate bytes when their sum fits within quota. It returns nil or the *CapacityError
// that refuses the candidate. There is no per-artifact limit.
func Admit(observed, candidate, quota uint64) error {
	if observed > math.MaxUint64-candidate {
		return &CapacityError{Observed: observed, Candidate: candidate, Quota: quota}
	}
	if observed+candidate > quota {
		return &CapacityError{Observed: observed, Candidate: candidate, Quota: quota, Required: observed + candidate}
	}
	return nil
}

// usage sums the published and temporary pack bytes in the store. A pack-named object that
// is not a regular file makes the store unsafe.
func usage(dir *os.Root) (uint64, error) {
	listing, err := dir.Open(".")
	if err != nil {
		return 0, refuse(RefuseUnsafe, "store directory is not listable: %v", err)
	}
	defer listing.Close()
	names, err := listing.Readdirnames(-1)
	if err != nil {
		return 0, refuse(RefuseUnsafe, "store directory is not listable: %v", err)
	}
	var total uint64
	for _, name := range names {
		if !isPackName(name) && !isTempName(name) {
			continue
		}
		info, err := dir.Lstat(name)
		if err != nil {
			return 0, refuse(RefuseUnsafe, "store object %s is not inspectable: %v", name, err)
		}
		if !info.Mode().IsRegular() {
			return 0, refuse(RefuseUnsafe, "store object %s is %s", name, kindOf(info))
		}
		size := uint64(info.Size())
		if size > math.MaxUint64-total {
			return math.MaxUint64, nil
		}
		total += size
	}
	return total, nil
}

func isPackName(name string) bool {
	hexPart, ok := strings.CutSuffix(name, PackSuffix)
	return ok && validDigest(hexPart)
}

func isTempName(name string) bool {
	return strings.HasPrefix(name, TempPrefix) && strings.HasSuffix(name, TempSuffix)
}

func packName(identity string) string {
	return strings.TrimPrefix(identity, IdentityPrefix) + PackSuffix
}

// Staged is one verified temporary pack that holds the writer lock until it publishes or
// is discarded.
type Staged struct {
	store     *Store
	dir       *os.Root
	operation *lock
	writer    *lock
	pack      *Pack
	temp      string
	attempt   int
	done      bool
}

// Stage writes a verified temporary pack for one preparation attempt. It holds the shared
// operation lock and the writer lock from the capacity calculation until Publish or Discard.
// An identical published artifact is verified and reused instead of written again.
func (s *Store) Stage(pack *Pack, quota uint64, attempt int) (*Staged, error) {
	dir, err := s.openDir(true)
	if err != nil {
		return nil, err
	}
	st := &Staged{store: s, dir: dir, pack: pack, attempt: attempt}
	if st.operation, err = acquire(dir, OperationLockName, syscall.LOCK_SH, true); err != nil {
		st.Discard()
		return nil, err
	}
	s.pause(StageWriterLock)
	if st.writer, err = acquire(dir, WriterLockName, syscall.LOCK_EX, true); err != nil {
		st.Discard()
		return nil, err
	}
	if _, err := dir.Lstat(packName(pack.identity)); err == nil {
		if err := verifyExisting(dir, pack.identity); err != nil {
			st.Discard()
			return nil, err
		}
		return st, nil
	}
	observed, err := usage(dir)
	if err != nil {
		st.Discard()
		return nil, err
	}
	if err := Admit(observed, uint64(len(pack.data)), quota); err != nil {
		st.Discard()
		return nil, err
	}
	s.pause(StageCapacity)
	if err := st.writeTemp(); err != nil {
		st.Discard()
		return nil, err
	}
	s.pause(StageStaged)
	if err := st.verifyTemp(); err != nil {
		st.Discard()
		return nil, err
	}
	return st, nil
}

func (st *Staged) writeTemp() error {
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		return refuse(RefuseStorage, "temporary name is not available: %v", err)
	}
	st.temp = TempPrefix + hex.EncodeToString(random) + TempSuffix
	file, err := st.dir.OpenFile(st.temp, os.O_WRONLY|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW, 0o600)
	if err != nil {
		st.temp = ""
		return refuse(RefuseStorage, "temporary pack is not creatable: %v", err)
	}
	defer file.Close()
	if err := st.store.fault(StepWriteTemp); err != nil {
		return refuse(RefuseStorage, "temporary pack write failed: %v", err)
	}
	if _, err := file.Write(st.pack.data); err != nil {
		return refuse(RefuseStorage, "temporary pack write failed: %v", err)
	}
	if err := st.store.fault(StepSyncTemp); err != nil {
		return refuse(RefuseStorage, "temporary pack sync failed: %v", err)
	}
	if err := file.Sync(); err != nil {
		return refuse(RefuseStorage, "temporary pack sync failed: %v", err)
	}
	return nil
}

// verifyTemp rereads the complete temporary pack and validates every byte against the
// expected identity before publication can expose it.
func (st *Staged) verifyTemp() error {
	if err := st.store.fault(StepVerifyTemp); err != nil {
		return refuse(RefuseStorage, "temporary pack verification failed: %v", err)
	}
	data, err := readRegular(st.dir, st.temp)
	if err != nil {
		return err
	}
	st.store.pause(StageVerifying)
	if _, err := Read(data, st.pack.identity); err != nil {
		return refuse(RefuseStorage, "temporary pack verification failed: %v", err)
	}
	return nil
}

// verifyExisting fully validates an already published artifact. A corrupt artifact is
// never repaired or replaced implicitly.
func verifyExisting(dir *os.Root, identity string) error {
	data, err := readRegular(dir, packName(identity))
	if err != nil {
		return refuse(RefuseExisting, "published artifact is unreadable: %v", err)
	}
	if _, err := Read(data, identity); err != nil {
		return refuse(RefuseExisting, "published artifact fails verification: %v", err)
	}
	return nil
}

// Publish exposes the verified temporary pack under its identity without overwriting an
// existing artifact, then releases the locks. attempt must name the attempt that staged
// the pack, so a stale candidate cannot publish.
func (st *Staged) Publish(attempt int) (string, error) {
	defer st.Discard()
	if st.done || attempt != st.attempt {
		return "", refuse(RefusePublish, "the staged candidate does not belong to the successful attempt")
	}
	if st.temp == "" {
		return st.pack.identity, nil
	}
	final := packName(st.pack.identity)
	if err := st.store.fault(StepLink); err != nil {
		return "", refuse(RefusePublish, "atomic publication failed: %v", err)
	}
	if err := st.dir.Link(st.temp, final); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			return "", refuse(RefusePublish, "atomic publication failed: %v", err)
		}
		if err := verifyExisting(st.dir, st.pack.identity); err != nil {
			return "", err
		}
	}
	if dir, err := st.dir.Open("."); err == nil {
		_ = dir.Sync()
		dir.Close()
	}
	return st.pack.identity, nil
}

// Discard removes the temporary pack, if any, and releases every lock. It is idempotent.
func (st *Staged) Discard() {
	if st == nil || st.done {
		return
	}
	st.done = true
	if st.temp != "" {
		_ = st.dir.Remove(st.temp)
	}
	st.writer.release()
	st.operation.release()
	if st.dir != nil {
		st.dir.Close()
	}
}
