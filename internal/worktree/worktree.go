// Package worktree owns pool leases, cleanup, and the subshell.
package worktree

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy"
	refreshop "github.com/gibbonmi/bench/internal/worktree/refresh"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func textDigest(value string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(value))) }

// Pool is the boundary form of poolAt for a caller in another package: it resolves
// the Bench home at the effect boundary. In-package callers below the boundary
// receive home explicitly.
func Pool(root string) string { return poolAt(Home(), root) }

// poolAt derives the repository's pool directory from an explicitly resolved home.
func poolAt(home, root string) string {
	return filepath.Join(home, "worktrees", poolkey.Key(root))
}
func LeaseFile(path string) (string, error) {
	lease, err := git.Output("-C", path, "rev-parse", "--git-path", git.BenchLeaseFilename)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(lease) {
		lease = filepath.Join(path, lease)
	}
	return lease, nil
}
func PoolCommand(home string, args []string) (string, int) {
	var root string
	if len(args) > 0 {
		root = args[0]
	} else {
		r, err := git.Root()
		if err != nil {
			return toon.NotInRepo() + "\n", 1
		}
		root = r
	}
	return poolAt(home, root) + "\n", 0
}
func LeaseFileCommand(args []string) (string, int) {
	if len(args) == 0 {
		return "usage: bench worktree-lease-file <path>\n", 2
	}
	lease, err := LeaseFile(args[0])
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	return lease + "\n", 0
}

type Class string

const (
	ClassRoot      Class = "root"
	ClassPoolWarm  Class = "pool-warm"
	ClassPoolLease Class = "pool-leased"
	ClassOutOfPool Class = "out-of-pool"
)

type Registered struct {
	Path     string
	Class    Class
	Branch   string
	Detached bool
	Locked   bool
}

// ClassifyRegisteredWorktrees is the boundary form of classifyRegisteredWorktreesAt for
// a caller in another package: it resolves the Bench home at the effect boundary.
func ClassifyRegisteredWorktrees(root string) ([]Registered, error) {
	return classifyRegisteredWorktreesAt(root, Home())
}

func classifyRegisteredWorktreesAt(root, home string) ([]Registered, error) {
	facts, err := git.Worktrees(root)
	if err != nil {
		return nil, err
	}
	mainRoot := canonicalRoot(root)
	out := make([]Registered, 0, len(facts))
	for _, fact := range facts {
		out = append(out, Registered{Path: fact.Path, Branch: fact.Branch, Detached: fact.Detached, Locked: fact.Locked})
	}
	pool := poolAt(home, mainRoot)
	for i := range out {
		out[i].Class = classifyPath(mainRoot, pool, out[i].Path)
	}
	return out, nil
}

// canonicalRoot is the in-package spelling of the shared derivation. A path builder
// below receives an already canonical root, so it never repeats the resolution.
func canonicalRoot(root string) string { return poolkey.Canonical(root) }
func classifyPath(root, pool, path string) Class {
	if samePath(path, root) {
		return ClassRoot
	}
	if insidePool(pool, path) {
		shape, err := ClassifyPathShape(path)
		if err != nil || shape != ShapeCheckoutDirectory {
			return ClassPoolWarm
		}
		lease, _ := LeaseFile(path)
		if isRegularFile(lease) {
			return ClassPoolLease
		}
		return ClassPoolWarm
	}
	return ClassOutOfPool
}

// wellFormedFingerprint accepts exactly what fingerprintParts emits: lowercase hex of
// a sha256 digest. The destructive reclaim path keeps this strict form, so a truncated
// paste never authorizes a removal there.
func wellFormedFingerprint(value string) bool {
	return len(value) == sha256.Size*2 && lowercaseHex(value)
}

// wellFormedFingerprintOrPrefix additionally accepts a prefix of at least 8 characters.
// `bench worktree clean --apply` takes this form; its plan carries exactly one digest,
// so a vetted prefix stays unambiguous.
func wellFormedFingerprintOrPrefix(value string) bool {
	return len(value) >= minOperandPrefix && len(value) <= sha256.Size*2 && lowercaseHex(value)
}

func lowercaseHex(value string) bool {
	for _, r := range value {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return false
		}
	}
	return true
}

// matchesFingerprint accepts the full digest or a prefix of at least 8 characters;
// wellFormedFingerprint has already vetted the shape, and a plan carries exactly one
// digest, so a vetted prefix is unambiguous.
func matchesFingerprint(full, given string) bool {
	return given == full || len(given) >= minOperandPrefix && strings.HasPrefix(full, given)
}

func insidePool(pool, path string) bool {
	rel, err := filepath.Rel(pool, path)
	return err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
func samePath(a, b string) bool {
	ar, aerr := filepath.EvalSymlinks(a)
	br, berr := filepath.EvalSymlinks(b)
	if aerr == nil {
		a = ar
	}
	if berr == nil {
		b = br
	}
	return filepath.Clean(a) == filepath.Clean(b)
}
func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// nestedState is the policy package's nested-repository verdict;
// internal/worktree/lifecyclepolicy owns its values.
type nestedState = lifecyclepolicy.NestedState

const (
	nestedClean         = lifecyclepolicy.NestedClean
	nestedDirty         = lifecyclepolicy.NestedDirty
	nestedEmbeddedClean = lifecyclepolicy.NestedEmbeddedClean
	nestedEmbeddedDirty = lifecyclepolicy.NestedEmbeddedDirty
	nestedUnknown       = lifecyclepolicy.NestedUnknown
)

func classifyNestedState(root string) (state nestedState, err error) {
	state = nestedClean
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		if entry.Name() == ".git" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.IsDir() {
			return nil
		}
		if _, err := os.Lstat(filepath.Join(path, ".git")); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		tracked, err := git.Output("-C", root, "ls-files", "--stage", "--", rel)
		if err != nil {
			return err
		}
		embedded := !strings.HasPrefix(tracked, "160000 ")
		out, err := exec.Command("git", "--no-optional-locks", "-C", path, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignore-submodules=none").Output()
		if err != nil {
			return err
		}
		if embedded {
			state = nestedEmbeddedClean
		}
		if len(out) != 0 && embedded {
			state = nestedEmbeddedDirty
		} else if len(out) != 0 {
			state = nestedDirty
		}
		return filepath.SkipDir
	})
	if err != nil {
		return nestedUnknown, err
	}
	return state, nil
}

// inRepository reports whether an explicit root names a git repository. The four
// verbs that receive a root refuse one that does not, in place of the
// working-directory lookup they made before.
func inRepository(root string) bool {
	if root == "" {
		return false
	}
	_, err := git.RootAt(root)
	return err == nil
}

func cleanInvocationError(stdout io.Writer) int {
	_ = renderCleanup(stdout, CleanupPlan{Target: "unknown", Action: ActionError, Tracked: "unknown", ignoredSummary: "unknown", Recovery: "none", Fingerprint: "none", Reason: "invalid invocation; run " + usage.WorktreeClean})
	return 2
}
func CleanCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	return cleanCommandWith(defaultJoins(), root, home, args, stdout, stderr)
}

// cleanCommandWith is CleanCommand with the seam set resolved explicitly at the caller's
// boundary.
func cleanCommandWith(j joins, root, _ string, args []string, stdout, stderr io.Writer) int {
	options := CleanupOptions{}
	target, fingerprint := "", ""
	landed := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--discard-ignored":
			options.DiscardIgnored = true
		case "--discard-branch":
			options.DiscardBranch = true
		case "--full":
			options.Full = true
		case "--landed":
			if landed {
				return cleanInvocationError(stdout)
			}
			landed = true
		case "--apply":
			if i+1 >= len(args) || fingerprint != "" {
				return cleanInvocationError(stdout)
			}
			i++
			fingerprint = args[i]
		case "--":
			if i+1 >= len(args) || target != "" {
				return cleanInvocationError(stdout)
			}
			i++
			target = args[i]
		default:
			if target != "" || strings.HasPrefix(args[i], "-") {
				return cleanInvocationError(stdout)
			}
			target = args[i]
		}
	}
	if target == "" && !landed || target != "" && landed {
		return cleanInvocationError(stdout)
	}
	if fingerprint != "" && !wellFormedFingerprintOrPrefix(fingerprint) {
		return cleanInvocationError(stdout)
	}
	if !inRepository(root) {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	if target != "" {
		target = resolveVerbOperand(root, target)
	}
	if landed {
		set, planErr := planLandedSet(j, root, options)
		if planErr != nil {
			_ = renderCleanup(stdout, CleanupPlan{Target: "unknown", Action: ActionError, Tracked: "unknown", ignoredSummary: "unknown", Recovery: "none", Fingerprint: fingerprint, Reason: planErr.Error()})
			return 1
		}
		if fingerprint != "" && len(set.rows) == 0 {
			return cleanInvocationError(stdout)
		}
		if fingerprint != "" && !matchesFingerprint(set.fingerprint, fingerprint) {
			_ = renderLandedStale(stdout, set, fingerprint)
			return 1
		}
		if fingerprint != "" {
			plans, applyErr := applyLandedSet(j, root, set, options)
			if renderErr := renderCleanups(stdout, plans); renderErr != nil {
				fmt.Fprintf(stderr, "bench worktree clean: %v\n", renderErr)
				return 1
			}
			if applyErr != nil {
				return 1
			}
			return 0
		}
		if err := renderLandedSet(stdout, set, options); err != nil {
			fmt.Fprintf(stderr, "bench worktree clean: %v\n", err)
			return 1
		}
		return 0
	}
	plan, err := planExplicitWith(j, root, target, options)
	if err == nil && fingerprint != "" {
		plan, err = applyExplicitWith(j, root, target, fingerprint, options)
	}
	if errors.Is(err, errStaleFingerprint) {
		_ = renderCleanup(stdout, plan)
		return 1
	}
	if err != nil {
		if plan.Target == "" {
			plan.Target, _ = canonicalPath(target)
		}
		plan.Action, plan.Reason = ActionError, err.Error()
		if plan.Fingerprint == "" {
			plan.Fingerprint = fingerprint
		}
		_ = renderCleanup(stdout, plan)
		return 1
	}
	if err := renderCleanup(stdout, plan); err != nil {
		fmt.Fprintf(stderr, "bench worktree clean: %v\n", err)
		return 1
	}
	if plan.unresolved {
		return 1
	}
	return 0
}
func finishReleaseReceipt(root string, stdout io.Writer, receipt intent.CleanupReceipt) int {
	if assignment, err := assignmentByID(root, receipt.Tracked); err == nil && assignment.State == intent.StateComplete {
		if err := intent.DeleteAssignment(root, assignment.ID); err != nil {
			return 1
		}
	}
	return renderReleaseReceipt(stdout, receipt)
}
func ReleaseCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	return releaseCommandWith(defaultJoins(), root, home, args, stdout, stderr)
}

// releaseCommandWith is ReleaseCommand with the seam set resolved explicitly at the
// caller's boundary.
func releaseCommandWith(j joins, root, _ string, args []string, stdout, stderr io.Writer) int {
	if len(args) != 3 || args[0] != "--request" || args[1] == "" {
		fmt.Fprintln(stderr, "usage: "+usage.WorktreeRelease)
		return 2
	}
	receipt, err := releaseAssignment(j, root, args[1], resolveVerbOperand(root, args[2]))
	if err == nil {
		return finishReleaseReceipt(root, stdout, receipt)
	}
	fmt.Fprintf(stderr, "bench worktree release: %v\n", err)
	return 1
}

func renderResumeSummary(result ResumeResult) string {
	var summary strings.Builder
	fmt.Fprintf(&summary, "bench resume: removed %d, swept refs %d", result.Removed, result.SweptRefs)
	retained := 0
	for _, count := range result.Retained {
		retained += count
	}
	if retained > 0 {
		summary.WriteString("; retained")
		for _, reason := range []CleanupReason{ReasonForeign, ReasonActive, ReasonLanded, ReasonOrphaned, ReasonLiveLease, ReasonUnmerged, ReasonIgnored, ReasonDirty, ReasonMalformed, ReasonUncertain, ReasonUnexpectedLock} {
			if count := result.Retained[reason]; count > 0 {
				fmt.Fprintf(&summary, " %s=%d", reason, count)
			}
		}
	}
	fmt.Fprintf(&summary, "; pruned branches %d; reconciled %d; failed %d; open assignments %d\n", result.PrunedBranches, result.Reconciled, result.Failed, result.Open)
	if result.Retained[ReasonLanded] > 0 {
		summary.WriteString("landed: bench worktree clean --landed (plans only; re-run with --apply <fingerprint> to remove)\n")
	}
	if result.PoolUnreadable != nil {
		fmt.Fprintf(&summary, "pool: not read (%v); bench worktree reclaim reports it properly\n", result.PoolUnreadable)
	} else if result.ReclaimableKeys > 0 {
		fmt.Fprintf(&summary, "pool: %d reclaimable keys; bench worktree reclaim (plans only; re-run with --apply <fingerprint> to remove)\n", result.ReclaimableKeys)
	}
	listCapped(&summary, len(result.Orphans), func(i int) string { return orphanLine(result.Orphans[i]) })
	return summary.String()
}

// resumeListingCap bounds each listing in the resume summary. This prints at every
// session start, so an unbounded backlog would bury the counted line above it and be
// scrolled past rather than read.
const resumeListingCap = 3

// listCapped writes at most resumeListingCap lines and, when the cap bites, one line
// naming both how many it withheld and the true total. Stating the total is what keeps a
// bounded listing from reading as the whole set. Every record stays listable through
// `bench worktree list`.
func listCapped(summary *strings.Builder, count int, line func(int) string) {
	for i := 0; i < count && i < resumeListingCap; i++ {
		summary.WriteString(line(i))
	}
	if count > resumeListingCap {
		fmt.Fprintf(summary, "and %d more (%d total)\n", count-resumeListingCap, count)
	}
}

// orphanLine renders one abandoned assignment's retirement command. The bare
// `bench worktree clean` prints a plan and a fingerprint and removes nothing. The line
// names the apply half rather than reading as one destructive step. It never suggests
// `--discard-ignored`, whose request-less form orphans the assignment. That remedy would
// manufacture the next generation of the residue reported here.
//
// A path this sink cannot carry verbatim is replaced by a pointer rather than escaped or
// digested. Quoting alone is not enough: single quotes make a newline literal but still
// emit the byte, which forges an extra line. An escaped path would name a tree that does
// not exist, so the reader would paste a command that cannot work.
func orphanLine(orphan OrphanCandidate) string {
	if !lineSafe(orphan.Path) {
		return fmt.Sprintf("orphan %s: worktree path holds control bytes; find its id row in bench worktree list\n", orphan.ID)
	}
	return fmt.Sprintf("orphan %s: bench worktree clean %s (plans only; re-run with --apply <fingerprint> to remove)\n", orphan.ID, sanitize.ShellQuote(orphan.Path))
}

// lineSafe is the package-local spelling of the shared line-structure predicate. It is
// deliberately stricter than cleanupOutputSafe: toon.Representable admits tab, newline,
// and return, because the TOON encoder escapes them, while the resume summary writes raw
// lines and escapes nothing.
func lineSafe(value string) bool { return sanitize.LineSafe(value) }

type assignmentRecoveryContext struct {
	target string
	// suffix is the calling verb's own clause after the component sentence. The release
	// verb names its retained checkout there; the landing verbs add nothing.
	suffix string
	base   string
	tip    string
}

// retainedSuffix is the release verb's clause. A release that refuses leaves the
// checkout in place, and the operator has to know that before choosing a next command.
const retainedSuffix = "; checkout retained"

// assignmentForRequest keeps opaque-token resolution in intent. Only an unmatched
// token permits path-derived recovery discovery.
func assignmentForRequest(root, request string, recoveryContext assignmentRecoveryContext) (intent.Assignment, error) {
	assignment, found, err := intent.FindAssignmentForRequest(root, request)
	if err != nil {
		return intent.Assignment{}, err
	}
	if found {
		return assignment, nil
	}
	recovery, _, err := unmatchedRequestRecovery(root, recoveryContext)
	if err != nil {
		return intent.Assignment{}, err
	}
	return intent.Assignment{}, refusalError{recovery}
}

func unmatchedRequestRecovery(root string, recoveryContext assignmentRecoveryContext) (refusal, bool, error) {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return refusal{}, false, err
	}
	recovery := componentRefusal(componentRequest, "", "", "").refusal
	recovery.detail += recoveryContext.suffix
	candidate, count := intent.Assignment{}, 0
	for _, assignment := range assignments {
		if assignment.State == intent.StateActive && assignment.Worktree == recoveryContext.target {
			candidate, count = assignment, count+1
		}
	}
	if count == 1 {
		recovery.observed = "assignment:" + candidate.ID
		recovery.next = reauthorizeRecoveryNext(candidate.ID, recoveryContext.target, recoveryContext.base, recoveryContext.tip)
		return recovery, true, nil
	}
	// Without one active assignment at the target there is no id to reauthorize, so the
	// route is the listing that names which assignment owns which tree.
	recovery.next = "bench worktree list"
	return recovery, false, nil
}

func reauthorizeRecoveryNext(assignment, target, base, tip string) string {
	baseArg := reauthorizeIdentityArg(base, "<full-base-commit>")
	tipArg := reauthorizeIdentityArg(tip, "<full-source-tip-commit>")
	command := "bench worktree reauthorize --assignment " + assignment + " --request <new-request> --base " + baseArg + " --source-tip " + tipArg
	if lineSafe(target) {
		return command + " " + sanitize.ShellQuote(target)
	}
	return "bench worktree exec " + assignment + " -- " + command + " ."
}

func reauthorizeIdentityArg(value, placeholder string) string {
	if !fullCommitIdentity(value) {
		return placeholder
	}
	return sanitize.ShellQuote(value)
}

func fullCommitIdentity(value string) bool { return len(value) == 40 && hexIdentity(value) }

func hexIdentity(value string) bool {
	for _, b := range []byte(value) {
		if !((b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')) {
			return false
		}
	}
	return true
}

var createGrammar = usage.Grammar{
	Cmd:  "bench worktree create",
	Help: "usage: " + usage.WorktreeCreate,
	Flags: []usage.Flag{
		{Name: "--request", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--label", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--refresh", HasValue: false},
	},
}

func CreateCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(createGrammar, args)
	if line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
			return 0
		}
		fmt.Fprintln(stderr, line)
		return code
	}
	_, startRef := refreshop.Consume(root, args, stdout)
	request, label := parsed.Flags["--request"], parsed.Flags["--label"]
	creation, err := createAt(defaultJoins(), root, home, request, label, nil, currentTime(), startRef)
	if err != nil {
		fmt.Fprintf(stderr, "bench worktree create: %v\n", err)
		return 1
	}
	out, err := toon.Table("worktree_create", []string{"path", "assignment", "state"}, [][]string{{creation.Path, creation.Assignment.ID, string(creation.Assignment.State)}})
	if err != nil {
		fmt.Fprintf(stderr, "bench worktree create: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, out)
	fmt.Fprintf(stdout, "next[2]:\n  bench worktree exec \"%s\" -- <command>\n  bench worktree path \"%s\"\n", label, label)
	return 0
}
func Subshell(home string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	args, startRef := refreshop.Consume(root, args, stdout)
	objective := strings.Join(args, " ")
	if objective == "" {
		objective = "interactive worktree"
	}
	request, err := randomID()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	creation, err := createAt(defaultJoins(), root, home, request, objective, nil, currentTime(), startRef)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stderr, "🪵 worktree: %s  (exit to release)\n", creation.Path)
	shell := subshellShell()
	if shell == "" {
		shell = "bash"
	}
	cmd := exec.Command(shell)
	cmd.Dir, cmd.Stdin, cmd.Stdout, cmd.Stderr = creation.Path, stdin, stdout, stderr
	// The interactive shell runs Bench verbs against the worktree this call just
	// created, so it reads the home the call resolved rather than the one its own
	// process carried.
	cmd.Env = withHome(os.Environ(), home)
	_ = cmd.Run()
	return ReleaseCommand(root, home, []string{"--request", request, creation.Path}, io.Discard, stderr)
}
func cleanupOutputSafe(value string) bool { return toon.Representable(value) }
func cleanupOutputValue(value string) string {
	if cleanupOutputSafe(value) {
		return value
	}
	return "sha256:" + textDigest(value)
}
