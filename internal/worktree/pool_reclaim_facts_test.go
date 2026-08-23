package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/worktree/reclaimpolicy"
)

// The tests in this file are the FA1 focused fact-adapter proofs for the
// reclaim policy: real pool state — registration-created worktrees, dead and
// live pointers, hostile shapes, denied reads — must translate into the exact
// typed facts internal/worktree/reclaimpolicy consumes.

// TestReclaimFactAdapterTranslatesDeadAndLivePointers is the pointer and path
// fact group. A planted dead child and a live child must translate into a
// regular pointer with its exact body and a proven Absent or Present target,
// and the production classifier must reach the child policy's verdict over
// exactly those facts.
func TestReclaimFactAdapterTranslatesDeadAndLivePointers(t *testing.T) {
	pool, _ := newReclaimPool(t)
	deadTarget := plantDeadChild(t, pool, "dead-key", "wt")
	plantLiveChild(t, pool, "live-key", "wt")

	dead := gatherPoolKeyFacts(filepath.Join(pool, "dead-key"), "dead-key")
	requireTest(t, dead.Shape == reclaimpolicy.ShapeDir && len(dead.Children) == 1, "dead-key facts = %+v, want one directory child", dead)
	pointer := dead.Children[0].Pointer
	requireTest(t, dead.Children[0].Name == "wt" && dead.Children[0].Shape == reclaimpolicy.ShapeDir,
		"dead child facts = %+v, want a directory named wt", dead.Children[0])
	requireTest(t, pointer.Shape == reclaimpolicy.ShapeRegular && pointer.Body == "gitdir: "+deadTarget+"\n",
		"dead pointer facts = %+v, want the exact regular-file body", pointer)
	requireTest(t, pointer.TargetExistence == reclaimpolicy.ExistenceAbsent && pointer.TargetErr == "",
		"dead pointer existence = %+v, want a proven absence", pointer)

	live := gatherPoolKeyFacts(filepath.Join(pool, "live-key"), "live-key")
	requireTest(t, live.Children[0].Pointer.TargetExistence == reclaimpolicy.ExistencePresent,
		"live pointer facts = %+v, want a present target", live.Children[0].Pointer)

	requireTest(t, reclaimpolicy.ClassifyKey(dead).Reclaimable() && !reclaimpolicy.ClassifyKey(live).Reclaimable(),
		"policy verdicts over the translated facts diverge from the pool's truth")
	produced := classifyPoolKey(filepath.Join(pool, "dead-key"), "dead-key")
	requireTest(t, produced.Verdict == poolVerdictReclaim && produced.Targets[0] == deadTarget,
		"production classifier = %#v, want the policy verdict over these facts", produced)
}

// TestReclaimFactAdapterTranslatesARealRegistrationKey is the registration
// fact group. A worktree git itself registered under a real repository's pool
// key must translate into a present gitdir: target, and the same key must
// read provably absent after the source repository is deleted.
func TestReclaimFactAdapterTranslatesARealRegistrationKey(t *testing.T) {
	pool, _ := newReclaimPool(t)
	source := newWorktreeRepo(t)
	key := filepath.Base(Pool(canonicalRoot(source)))
	created := mustCreate(t, source, "fa-registration", "fareg")
	requireTest(t, filepath.Dir(created.Path) == filepath.Join(pool, key), "created worktree %q is not under key %q", created.Path, key)

	registered := gatherPoolKeyFacts(filepath.Join(pool, key), key)
	requireTest(t, len(registered.Children) == 1 && registered.Children[0].Pointer.TargetExistence == reclaimpolicy.ExistencePresent,
		"registered facts = %+v, want the real registration's target present", registered)
	target, ok := reclaimpolicy.GitdirTarget(registered.Children[0].Pointer.Body)
	requireTest(t, ok && filepath.IsAbs(target) && strings.Contains(target, filepath.Join(".git", "worktrees")),
		"registration pointer body %q parsed to %q, want git's admin directory", registered.Children[0].Pointer.Body, target)

	mustNoError(t, os.RemoveAll(source))
	orphaned := gatherPoolKeyFacts(filepath.Join(pool, key), key)
	requireTest(t, orphaned.Children[0].Pointer.TargetExistence == reclaimpolicy.ExistenceAbsent,
		"orphaned facts = %+v, want the deleted repository's target absent", orphaned)
}

// TestReclaimFactAdapterTranslatesHostileShapesUnopened is the shape fact
// group. A symlink key must translate to its shape with no descent, and a
// FIFO where a pointer belongs must translate to a non-regular shape with no
// body — proof the adapter never opened it.
func TestReclaimFactAdapterTranslatesHostileShapesUnopened(t *testing.T) {
	pool, _ := newReclaimPool(t)

	mustNoError(t, os.Symlink(t.TempDir(), filepath.Join(pool, "symlinked-key")))
	linked := gatherPoolKeyFacts(filepath.Join(pool, "symlinked-key"), "symlinked-key")
	requireTest(t, linked.Shape == reclaimpolicy.ShapeSymlink && len(linked.Children) == 0,
		"symlink facts = %+v, want an undescended symlink shape", linked)

	child := filepath.Join(pool, "fifo-git", "wt")
	mustMkdirAll(t, child, 0o755)
	if err := syscall.Mkfifo(filepath.Join(child, ".git"), 0o644); err != nil {
		capability.Capability(t, capability.Fifo, "FIFOs unavailable on this filesystem: "+err.Error())
	}
	fifo := gatherPoolKeyFacts(filepath.Join(pool, "fifo-git"), "fifo-git")
	pointer := fifo.Children[0].Pointer
	requireTest(t, pointer.Shape == reclaimpolicy.ShapeOther && pointer.Body == "" && pointer.ReadErr == "",
		"fifo pointer facts = %+v, want an unopened non-regular shape", pointer)
}

// TestReclaimFactAdapterTranslatesDeniedReads is the uncertainty fact group.
// A permission-denied child and a permission-denied target probe must
// translate into unreadable shapes and unknown existence carrying the exact
// error text, never into proof of absence.
func TestReclaimFactAdapterTranslatesDeniedReads(t *testing.T) {
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses directory permissions; cannot deny the stat that leaves existence unknown")
	}
	pool, _ := newReclaimPool(t)

	plantDeadChild(t, pool, "unreadable-child", "wt")
	sealed := filepath.Join(pool, "unreadable-child", "wt")
	t.Cleanup(func() { _ = os.Chmod(sealed, 0o700) })
	mustChmod(t, sealed, 0o000)
	denied := gatherPoolKeyFacts(filepath.Join(pool, "unreadable-child"), "unreadable-child")
	requireTest(t, denied.Children[0].Pointer.Shape == reclaimpolicy.ShapeUnreadable &&
		strings.Contains(denied.Children[0].Pointer.ShapeErr, "permission denied"),
		"denied pointer facts = %+v, want an unreadable shape naming the denial", denied.Children[0].Pointer)

	deniedDir := filepath.Join(t.TempDir(), "denied")
	mustMkdirAll(t, deniedDir, 0o700)
	plantChild(t, pool, "unstattable-target", "wt", filepath.Join(deniedDir, "gone", ".git"))
	t.Cleanup(func() { _ = os.Chmod(deniedDir, 0o700) })
	mustChmod(t, deniedDir, 0o000)
	unknown := gatherPoolKeyFacts(filepath.Join(pool, "unstattable-target"), "unstattable-target")
	pointer := unknown.Children[0].Pointer
	requireTest(t, pointer.TargetExistence == reclaimpolicy.ExistenceUnknown && strings.Contains(pointer.TargetErr, "permission denied"),
		"denied target facts = %+v, want unknown existence naming the denial", pointer)
}
