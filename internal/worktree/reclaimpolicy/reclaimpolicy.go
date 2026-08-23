// Package reclaimpolicy owns the pure pool-reclaim decisions of the worktree
// commands: current-repository key protection, key and child shape
// classification, pointer parsing, target-existence reclaimability, plan
// drift, removal bounding, requalification, removal outcomes, and the apply's
// exit verdict. The parent package translates filesystem, registration, and
// path state into these typed facts once at its effect boundary. The decisions
// here read only the supplied facts and return verdicts the parent projects;
// the package performs no effects, reads no ambient process state, and starts
// no descendants. The source census test enforces that boundary.
package reclaimpolicy

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// FingerprintVersion binds the fingerprint material to this plan's derivation.
// A later change to what the plan carries changes the digest rather than
// silently reusing a value an apply would still accept.
const FingerprintVersion = "bench-pool-reclaim/v1"

const (
	VerdictReclaim = "reclaim"
	VerdictRetain  = "retain"
	// The apply's verdicts are past tense: the plan says what may happen, the
	// apply says what did. Sharing the plan's two words would leave an operator
	// unable to tell a printed intention from a completed removal.
	VerdictRemoved  = "removed"
	VerdictRetained = "retained"
)

// KeyVerdict is one key's classification. Targets holds the gitdir: pointers
// that proved the key dead, in child order. It is empty for every retained key
// and for an empty one. It also feeds the fingerprint material, so a key whose
// children changed identity invalidates a plan that counted them.
type KeyVerdict struct {
	Key     string
	Verdict string
	Reason  string
	Targets []string
}

// Reclaimable reports whether the plan may remove this key.
func (v KeyVerdict) Reclaimable() bool { return v.Verdict == VerdictReclaim }

// ReclaimableCount counts the keys a plan over these verdicts would remove.
func ReclaimableCount(verdicts []KeyVerdict) int {
	count := 0
	for _, verdict := range verdicts {
		if verdict.Reclaimable() {
			count++
		}
	}
	return count
}

// EntryShape is the typed lstat verdict over one filesystem entry. The parent
// translates each Lstat once; the decisions here never touch the filesystem.
type EntryShape string

const (
	ShapeMissing    EntryShape = "missing"
	ShapeUnreadable EntryShape = "unreadable"
	ShapeSymlink    EntryShape = "symlink"
	ShapeDir        EntryShape = "dir"
	ShapeRegular    EntryShape = "regular"
	ShapeOther      EntryShape = "other"
)

// Existence is the typed probe verdict over a pointer's gitdir: target. Only
// ExistenceAbsent is proof of deadness. The zero value is not Absent, so an
// unfilled fact fails closed as unproven rather than as permission to remove.
type Existence string

const (
	ExistenceUnknown Existence = ""
	ExistencePresent Existence = "present"
	ExistenceAbsent  Existence = "absent"
)

// PointerFacts is the translated state of one child's `.git` entry. Body is
// filled only for a regular file the parent read; the parent never opens any
// other shape, so a FIFO where a pointer belongs cannot block the boundary.
// TargetErr carries the probe error text when existence stayed unknown.
type PointerFacts struct {
	Shape           EntryShape
	ShapeErr        string
	Body            string
	ReadErr         string
	TargetExistence Existence
	TargetErr       string
}

// ChildFacts is the translated state of one top-level pool-key entry.
type ChildFacts struct {
	Name     string
	Shape    EntryShape
	ShapeErr string
	Pointer  PointerFacts
}

// KeyFacts is the translated state of one pool key. Current marks the key the
// running repository owns; a current key carries no other facts because the
// protection decides before any filesystem question is asked. ListErr carries
// the enumeration error text when the key's contents could not be listed.
type KeyFacts struct {
	Name     string
	Current  bool
	Shape    EntryShape
	ShapeErr string
	ListErr  string
	Children []ChildFacts
}

// ClassifyKey is the one reclaimability predicate in the tree. A key is
// reclaimable when it holds nothing at top level. It is also reclaimable when
// every top-level entry is a real directory holding a regular `.git` file
// whose `gitdir:` target is provably absent. Everything else retains and says
// what protected it.
//
// Absence is proven only by ExistenceAbsent. Every other existence leaves the
// question open, and unknown retains — treating a permission failure as
// absence is the one direction that destroys work. A symlink where a key, a
// child, or a `.git` belongs is retained unfollowed. The pool keeps bounding
// what the apply can ever remove.
func ClassifyKey(facts KeyFacts) KeyVerdict {
	retain := func(format string, args ...any) KeyVerdict {
		return KeyVerdict{Key: facts.Name, Verdict: VerdictRetain, Reason: fmt.Sprintf(format, args...)}
	}
	// The current repository's key is excluded before the predicate runs. A
	// session between acquiring its pool directory and its first checkout holds
	// an empty key, which the empty-key clause would otherwise claim.
	if facts.Current {
		return retain("key belongs to the current repository")
	}
	switch {
	case facts.Shape == ShapeMissing || facts.Shape == ShapeUnreadable:
		return retain("key cannot be read: %s", facts.ShapeErr)
	case facts.Shape == ShapeSymlink:
		return retain("key is a symlink")
	case facts.Shape != ShapeDir:
		return retain("key is not a directory")
	}
	if facts.ListErr != "" {
		return retain("key contents cannot be listed: %s", facts.ListErr)
	}
	if len(facts.Children) == 0 {
		return KeyVerdict{Key: facts.Name, Verdict: VerdictReclaim, Reason: "key holds nothing"}
	}
	targets := make([]string, 0, len(facts.Children))
	for _, child := range facts.Children {
		target, verdict := classifyChild(child)
		if verdict != "" {
			return retain("%s", verdict)
		}
		targets = append(targets, target)
	}
	return KeyVerdict{Key: facts.Name, Verdict: VerdictReclaim, Reason: "every child points at an absent repository", Targets: targets}
}

// classifyChild decides one top-level entry. It returns the child's dangling
// gitdir: target, or the reason that entry protects the whole key. A key
// mixing one live and one dead pointer is retained whole, because partial
// reclamation would amputate the live half.
func classifyChild(child ChildFacts) (target, retain string) {
	switch {
	case child.Shape == ShapeMissing || child.Shape == ShapeUnreadable:
		return "", fmt.Sprintf("entry %s cannot be read: %s", child.Name, child.ShapeErr)
	case child.Shape == ShapeSymlink:
		return "", fmt.Sprintf("entry %s is a symlink", child.Name)
	case child.Shape != ShapeDir:
		return "", fmt.Sprintf("entry %s is not a directory", child.Name)
	}
	pointer := child.Pointer
	switch {
	case pointer.Shape == ShapeMissing:
		return "", fmt.Sprintf("child %s holds no .git entry", child.Name)
	case pointer.Shape == ShapeUnreadable:
		return "", fmt.Sprintf("child %s .git cannot be read: %s", child.Name, pointer.ShapeErr)
	case pointer.Shape == ShapeSymlink:
		return "", fmt.Sprintf("child %s .git is a symlink", child.Name)
	case pointer.Shape == ShapeDir:
		return "", fmt.Sprintf("child %s .git is a repository directory", child.Name)
	case pointer.Shape != ShapeRegular:
		return "", fmt.Sprintf("child %s .git is not a regular file", child.Name)
	}
	if pointer.ReadErr != "" {
		return "", fmt.Sprintf("child %s .git cannot be read: %s", child.Name, pointer.ReadErr)
	}
	target, ok := GitdirTarget(pointer.Body)
	if !ok {
		return "", fmt.Sprintf("child %s .git carries no gitdir: target", child.Name)
	}
	// A relative target is resolved by git against the child directory, not
	// against whatever directory the probing process happened to be in. The
	// probe would have answered a question about the wrong path, and an absent
	// answer would be proof of nothing.
	if !filepath.IsAbs(target) {
		return "", fmt.Sprintf("child %s gitdir: target %q is not absolute", child.Name, target)
	}
	switch pointer.TargetExistence {
	case ExistenceAbsent:
		return target, ""
	case ExistencePresent:
		return "", fmt.Sprintf("child %s gitdir: target exists", child.Name)
	default:
		return "", fmt.Sprintf("child %s gitdir: target cannot be read: %s", child.Name, pointer.TargetErr)
	}
}

// GitdirTarget reads the pointer a git worktree's `.git` file carries. A file
// with no `gitdir:` line, or one whose value is blank, reports no target: an
// unparseable pointer is never proof that anything is absent.
func GitdirTarget(body string) (string, bool) {
	for line := range strings.SplitSeq(body, "\n") {
		// Only the line terminator and the separator git writes after the colon
		// are stripped. A repository path may legitimately end in a space, and
		// trimming that away would make a live worktree read as absent — proof
		// of nothing.
		rest, found := strings.CutPrefix(strings.TrimLeft(strings.TrimRight(line, "\r"), " \t"), "gitdir:")
		if !found {
			continue
		}
		if target := strings.TrimLeft(rest, " \t"); strings.TrimSpace(target) != "" {
			return target, true
		}
	}
	return "", false
}

// FingerprintMaterial derives exactly what an apply would remove: the
// fingerprint version, the reclaimable key names in pool order, and each one's
// child gitdir: targets. A change elsewhere in the pool leaves it alone, so an
// operator is not sent back to re-plan by an unrelated key. A change to a
// target inside a counted key does move it. The parent digests the material
// with its shared fingerprint owner.
func FingerprintMaterial(verdicts []KeyVerdict) [][]byte {
	parts := [][]byte{[]byte(FingerprintVersion)}
	for _, verdict := range verdicts {
		if !verdict.Reclaimable() {
			continue
		}
		parts = append(parts, []byte("key"), []byte(verdict.Key), []byte(strconv.Itoa(len(verdict.Targets))))
		for _, target := range verdict.Targets {
			parts = append(parts, []byte("target"), []byte(target))
		}
	}
	return parts
}

// PlanDrift reports whether a supplied fingerprint no longer matches the pool
// the apply just re-read. A drifted plan is a reading that has stopped being
// true; nothing is removed on the strength of it.
func PlanDrift(supplied, current string) bool { return supplied != current }

// RetainedOnApply carries a plan's retention into the apply's past-tense
// report unchanged: the key was never a target, and its protecting reason
// still says why.
func RetainedOnApply(verdict KeyVerdict) KeyVerdict {
	return KeyVerdict{Key: verdict.Key, Verdict: VerdictRetained, Reason: verdict.Reason}
}

// RemovalBounds resolves and bounds one removal target. The pool is what
// bounds the apply's blast radius, so a key name whose joined path is not a
// direct child of the pool is refused without touching the bytes it points at.
func RemovalBounds(pool, key string) (target string, refusal KeyVerdict, ok bool) {
	target = filepath.Join(pool, key)
	if filepath.Dir(target) != pool {
		return "", KeyVerdict{Key: key, Verdict: VerdictRetained,
			Reason: fmt.Sprintf("removal target is not a direct child of %s", pool)}, false
	}
	return target, KeyVerdict{}, true
}

// RemovalRequalified re-judges one key at the instant of removal. The
// fingerprint spoke for the plan as a whole; only this re-check speaks for
// this one key now. A key that stopped qualifying survives with the reason
// that protects it.
func RemovalRequalified(key string, current KeyVerdict) (refusal KeyVerdict, ok bool) {
	if current.Reclaimable() {
		return KeyVerdict{}, true
	}
	return KeyVerdict{Key: key, Verdict: VerdictRetained,
		Reason: fmt.Sprintf("key stopped qualifying before removal: %s", current.Reason)}, false
}

// RemovalOutcome is the apply's per-key report over the one destructive
// effect. An empty error text is a completed removal; anything else retains
// with the failure, so a script-read exit can refuse to call the run complete.
func RemovalOutcome(key, removeErr string) KeyVerdict {
	if removeErr != "" {
		return KeyVerdict{Key: key, Verdict: VerdictRetained, Reason: fmt.Sprintf("removal failed: %s", removeErr)}
	}
	return KeyVerdict{Key: key, Verdict: VerdictRemoved, Reason: "key removed"}
}

// ApplyIncomplete reports whether a key the plan named survived the apply.
// That includes a failed removal and a key that stopped qualifying between the
// re-plan and its own re-check. A key the plan retained was never a target, so
// it does not make a clean run look failed. The rows say which and why; this
// verdict is what the exit code carries to a script.
func ApplyIncomplete(planned, applied []KeyVerdict) bool {
	targets := map[string]bool{}
	for _, verdict := range planned {
		if verdict.Reclaimable() {
			targets[verdict.Key] = true
		}
	}
	for _, verdict := range applied {
		if targets[verdict.Key] && verdict.Verdict != VerdictRemoved {
			return true
		}
	}
	return false
}
