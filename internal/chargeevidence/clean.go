package chargeevidence

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// Cleanup refusal classes. A busy store names an active reader or writer; a stale plan
// names a target inventory that changed after the plan was fingerprinted.
const (
	RefuseBusy      = "store-busy"
	RefuseStalePlan = "stale-plan"
)

// Target kinds. A published target is one complete artifact; an orphan target is one
// temporary object that no live writer holds.
const (
	TargetPublished = "published"
	TargetOrphan    = "orphan"
)

// StepRemove is the publication step name one target deletion reaches, so a test injects a
// deterministic deletion failure the same way it injects a publication failure.
const StepRemove = "remove"

// Target is one exact deletion target. ID is the artifact identity for a published pack and
// the store object name for an orphan temporary. name is the store object the deletion
// removes, and identity is the observed file identity the fingerprint commits to.
type Target struct {
	ID    string
	Kind  string
	Bytes uint64

	name     string
	identity string
}

// CleanupPlan is one store's complete, deterministic target list and the fingerprint that
// commits to it. An absent or empty store plans no target.
type CleanupPlan struct {
	Fingerprint string
	Targets     []Target
	Bytes       uint64
}

// Applied is the terminal disposition of one cleanup apply. A stopped apply reports the
// targets it removed and the targets a fresh plan must cover.
type Applied struct {
	Fingerprint        string
	Removed, Remaining int
	Complete           bool
}

// Plan enumerates every published artifact and every proven orphan temporary object under
// the exclusive evidence lock. An absent store plans nothing and creates nothing. A store
// with no lock file has published nothing, so it also plans nothing.
func (s *Store) Plan() (CleanupPlan, error) {
	dir, operation, err := s.exclusive()
	if err != nil || dir == nil {
		return CleanupPlan{Fingerprint: fingerprintOf(nil)}, err
	}
	defer dir.Close()
	defer operation.release()
	return collect(dir)
}

// Apply deletes the targets of the plan that fingerprint commits to. It reconstructs the
// complete plan under the exclusive lock and refuses a changed fingerprint before it deletes
// any target. It rechecks each target's file identity immediately before its deletion, and
// it stops on the first failure without claiming a rollback.
func (s *Store) Apply(fingerprint string) (Applied, error) {
	dir, operation, err := s.exclusive()
	if err != nil {
		return Applied{}, err
	}
	if dir == nil {
		if empty := fingerprintOf(nil); fingerprint != empty {
			return Applied{}, refuse(RefuseStalePlan, "the store holds no target the plan named")
		}
		return Applied{Fingerprint: fingerprint, Complete: true}, nil
	}
	defer dir.Close()
	defer operation.release()

	plan, err := collect(dir)
	if err != nil {
		return Applied{}, err
	}
	if plan.Fingerprint != fingerprint {
		return Applied{}, refuse(RefuseStalePlan, "the store changed after the plan was fingerprinted")
	}
	applied := Applied{Fingerprint: fingerprint, Remaining: len(plan.Targets)}
	for _, target := range plan.Targets {
		if err := s.remove(dir, target); err != nil {
			return applied, nil
		}
		applied.Removed++
		applied.Remaining--
	}
	applied.Complete = true
	return applied, nil
}

// remove deletes one target after proving the directory still names the regular file the
// plan observed. A changed identity or a failed deletion stops the apply.
func (s *Store) remove(dir *os.Root, target Target) error {
	info, err := dir.Lstat(target.name)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return refuse(RefuseUnsafe, "store object %s is %s", target.name, kindOf(info))
	}
	if identityOf(info) != target.identity {
		return refuse(RefuseReplaced, "store object %s changed after the plan was fingerprinted", target.name)
	}
	if err := s.fault(StepRemove); err != nil {
		return err
	}
	return dir.Remove(target.name)
}

// exclusive opens the store and takes the operation lock exclusively without waiting. Every
// reader and every writer holds that lock for its whole operation, so this one acquisition
// is the writer exclusion a plannable orphan needs; a second writer-lock acquisition would
// add no guarantee. A nil directory reports a store that plans nothing: absent, or present
// with no lock file because nothing has ever published.
func (s *Store) exclusive() (*os.Root, *lock, error) {
	dir, err := s.openDir(false)
	if err != nil {
		if refusalClass(err) == RefuseAbsent {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	operation, err := acquire(dir, OperationLockName, syscall.LOCK_EX|syscall.LOCK_NB, false)
	if err != nil {
		dir.Close()
		if refusalClass(err) == RefuseAbsent {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	return dir, operation, nil
}

// collect reads the complete target inventory of an opened store. A pack-named or
// temporary-named object that is not a regular file makes the store unsafe, so an unsafe
// target refuses the whole plan before any deletion.
func collect(dir *os.Root) (CleanupPlan, error) {
	listing, err := dir.Open(".")
	if err != nil {
		return CleanupPlan{}, refuse(RefuseUnsafe, "store directory is not listable: %v", err)
	}
	names, err := listing.Readdirnames(-1)
	listing.Close()
	if err != nil {
		return CleanupPlan{}, refuse(RefuseUnsafe, "store directory is not listable: %v", err)
	}
	var plan CleanupPlan
	for _, name := range names {
		published, temporary := isPackName(name), isTempName(name)
		if !published && !temporary {
			continue
		}
		info, err := dir.Lstat(name)
		if err != nil {
			return CleanupPlan{}, refuse(RefuseUnsafe, "store object %s is not inspectable: %v", name, err)
		}
		if !info.Mode().IsRegular() {
			return CleanupPlan{}, refuse(RefuseUnsafe, "store object %s is %s", name, kindOf(info))
		}
		target := Target{ID: name, Kind: TargetOrphan, Bytes: uint64(info.Size()), name: name, identity: identityOf(info)}
		if published {
			target.ID = IdentityPrefix + strings.TrimSuffix(name, PackSuffix)
			target.Kind = TargetPublished
		}
		plan.Targets = append(plan.Targets, target)
		plan.Bytes += target.Bytes
	}
	sort.Slice(plan.Targets, func(i, j int) bool {
		if plan.Targets[i].Kind != plan.Targets[j].Kind {
			return plan.Targets[i].Kind < plan.Targets[j].Kind
		}
		return plan.Targets[i].ID < plan.Targets[j].ID
	})
	plan.Fingerprint = fingerprintOf(plan.Targets)
	return plan, nil
}

// fingerprintOf commits to the complete target list: each target's kind, identity, byte
// length, and observed file identity, in plan order. A changed inventory, a changed length,
// and a replaced file each change the fingerprint.
func fingerprintOf(targets []Target) string {
	sum := sha256.New()
	for _, target := range targets {
		for _, field := range []string{target.Kind, target.ID, strconv.FormatUint(target.Bytes, 10), target.identity} {
			sum.Write([]byte(field))
			sum.Write([]byte{0})
		}
	}
	return IdentityPrefix + hex.EncodeToString(sum.Sum(nil))
}

// identityOf renders one observed file identity. A platform that reports no device and
// inode pair returns an empty identity, which the fingerprint still commits to.
func identityOf(info os.FileInfo) string {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return ""
	}
	return strconv.FormatUint(uint64(stat.Dev), 10) + ":" + strconv.FormatUint(uint64(stat.Ino), 10)
}

func refusalClass(err error) string {
	var refusal *Refusal
	if errors.As(err, &refusal) {
		return refusal.Class
	}
	return ""
}

// TargetsPerPage is how many exact targets one bounded cleanup page carries. Each target row
// is far shorter than the response limit divides by this count, so a page fits the bound
// before the final guard applies.
const TargetsPerPage = 100

// Page returns the targets at index and the cursor of the page after them. A terminal page
// returns a nil successor. An index beyond the plan refuses.
func (p CleanupPlan) Page(index int) ([]Target, *Cursor, error) {
	// An empty plan has exactly one page at target zero; every other page starts at a
	// target the plan names.
	if index < 0 || index >= max(len(p.Targets), 1) {
		return nil, nil, refuse(RefuseCursor, "cursor names target %d beyond %d targets", index, len(p.Targets))
	}
	end := min(index+TargetsPerPage, len(p.Targets))
	if end == len(p.Targets) {
		return p.Targets[index:end], nil, nil
	}
	return p.Targets[index:end], &Cursor{Identity: p.Fingerprint, Index: end, Clean: true}, nil
}

// Encode renders one cleanup plan page: the orientation row, then this page's target rows.
// next is the exact successor command, or empty at the end of the plan.
func (p CleanupPlan) Encode(targets []Target, last bool, next string) (string, error) {
	rows := make([][]any, len(targets))
	for i, target := range targets {
		rows[i] = []any{target.ID, target.Kind, target.Bytes}
	}
	out, err := encodeResponse(blockCleanup, []any{p.Fingerprint, len(p.Targets), p.Bytes, true, last, next})
	if err != nil {
		return "", err
	}
	page, err := encodeRows(blockTargets, rows)
	if err != nil {
		return "", err
	}
	return out + page, nil
}

// Encode renders the terminal cleanup disposition. next is the recovery action a stopped
// apply names, and it is empty after a complete apply.
func (a Applied) Encode(next string) (string, error) {
	return encodeResponse(blockApplied, []any{a.Fingerprint, a.Removed, a.Remaining, a.Complete, next})
}
