// Package worktree owns pool leases, cleanup, and the subshell.
package worktree

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	refreshop "github.com/gibbonmi/bench/internal/worktree/refresh"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func textDigest(value string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(value))) }
func cksum(data []byte) uint32 {
	var crc uint32
	step := func(value byte) {
		crc ^= uint32(value) << 24
		for range 8 {
			crc = crc<<1 ^ 0x04C11DB7*(crc>>31)
		}
	}
	for _, value := range data {
		step(value)
	}
	for n := len(data); n > 0; n >>= 8 {
		step(byte(n))
	}
	return ^crc
}
func benchHome() string {
	if h := os.Getenv("BENCH_HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bench")
}
func Pool(root string) string {
	sum := cksum([]byte(root + "\n"))
	key := filepath.Base(root) + "-" + strconv.FormatUint(uint64(sum), 10)
	return filepath.Join(benchHome(), "worktrees", key)
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
func PoolCommand(args []string) (string, int) {
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
	return Pool(root) + "\n", 0
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

func ClassifyRegisteredWorktrees(root string) ([]Registered, error) {
	facts, err := git.Worktrees(root)
	if err != nil {
		return nil, err
	}
	mainRoot := canonicalRoot(root)
	out := make([]Registered, 0, len(facts))
	for _, fact := range facts {
		out = append(out, Registered{Path: fact.Path, Branch: fact.Branch, Detached: fact.Detached, Locked: fact.Locked})
	}
	pool := Pool(mainRoot)
	for i := range out {
		out[i].Class = classifyPath(mainRoot, pool, out[i].Path)
	}
	return out, nil
}
func canonicalRoot(root string) string {
	common, err := git.CommonDir(root)
	if err != nil || filepath.Base(common) != ".git" {
		return root
	}
	return filepath.Dir(common)
}
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

type nestedState string

const nestedClean, nestedDirty, nestedEmbeddedClean, nestedEmbeddedDirty, nestedUnknown nestedState = "clean", "dirty", "embedded-clean", "embedded-dirty", "unknown"

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
func cleanInvocationError(stdout io.Writer) int {
	_ = renderCleanup(stdout, CleanupPlan{Target: "unknown", Action: ActionError, Tracked: "unknown", ignoredSummary: "unknown", Recovery: "none", Fingerprint: "none", Reason: "invalid invocation; run " + usage.WorktreeClean})
	return 2
}
func CleanCommand(args []string, stdout, stderr io.Writer) int {
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
	if fingerprint != "" {
		decoded, err := hex.DecodeString(fingerprint)
		if err != nil || len(decoded) != sha256.Size || fingerprint != strings.ToLower(fingerprint) {
			return cleanInvocationError(stdout)
		}
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	if landed {
		set, planErr := planLandedSet(root, options)
		if planErr != nil {
			_ = renderCleanup(stdout, CleanupPlan{Target: "unknown", Action: ActionError, Tracked: "unknown", ignoredSummary: "unknown", Recovery: "none", Fingerprint: fingerprint, Reason: planErr.Error()})
			return 1
		}
		if fingerprint != "" && len(set.rows) == 0 {
			return cleanInvocationError(stdout)
		}
		if fingerprint != "" && fingerprint != set.fingerprint {
			_ = renderLandedStale(stdout, set, fingerprint)
			return 1
		}
		if fingerprint != "" {
			plans, applyErr := applyLandedSet(root, set, options)
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
	plan, err := PlanExplicitWithOptions(root, target, options)
	if err == nil && fingerprint != "" {
		plan, err = ApplyExplicitWithOptions(root, target, fingerprint, options)
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
func ReleaseCommand(root string, args []string, stdout, stderr io.Writer) int {
	if len(args) != 3 || args[0] != "--request" || args[1] == "" {
		fmt.Fprintln(stderr, "usage: "+usage.WorktreeRelease)
		return 2
	}
	receipt, err := releaseAssignment(root, args[1], args[2])
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
	listCapped(&summary, len(result.Orphans), func(i int) string { return orphanLine(result.Orphans[i]) })
	return summary.String()
}

// resumeListingCap bounds each listing in the resume summary. This prints at every
// session start, so an unbounded backlog would bury the counted line above it and be
// scrolled past rather than read.
const resumeListingCap = 3

// listCapped writes at most resumeListingCap lines and, when the cap bites, one line
// naming both how many it withheld and the true total. Stating the total is what keeps a
// bounded listing from reading as the whole set; every record stays listable through
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
// `bench worktree clean` prints a plan and a fingerprint and removes nothing, so the line
// names the apply half rather than reading as one destructive step. It never suggests
// `--discard-ignored`, whose request-less form orphans the assignment (FT93b) — the
// remedy would manufacture the next generation of the residue reported here.
//
// A path this sink cannot carry verbatim is replaced by a pointer rather than escaped or
// digested. Quoting alone is not enough — single quotes make a newline literal but still
// emit the byte, which forges an extra line — and an escaped path would name a tree that
// does not exist, so the reader would paste a command that cannot work.
func orphanLine(orphan OrphanCandidate) string {
	if !lineSafe(orphan.Path) {
		return fmt.Sprintf("orphan %s: worktree path holds control bytes; find its id row in bench worktree list\n", orphan.ID)
	}
	return fmt.Sprintf("orphan %s: bench worktree clean %s (plans only; re-run with --apply <fingerprint> to remove)\n", orphan.ID, sanitize.ShellQuote(orphan.Path))
}

// lineSafe reports whether a value carries no control rune. It is deliberately stricter
// than cleanupOutputSafe: toon.Representable admits tab, newline, and return because the
// TOON encoder escapes them, while the resume summary writes raw lines and escapes
// nothing, so a newline forges a line and an ESC drives the terminal that prints it.
// Display-hostile runes outside the control categories — a bidi override, U+2028,
// invalid UTF-8 — pass, so this guards the summary's line structure rather than how a
// terminal renders one line.
func lineSafe(value string) bool { return !strings.ContainsFunc(value, unicode.IsControl) }
func CreateCommand(root string, args []string, stdout, stderr io.Writer) int {
	var request, label string
	args, startRef := refreshop.Consume(root, args, stdout)
	for len(args) > 0 {
		if len(args) < 2 || (args[0] != "--request" && args[0] != "--label") {
			fmt.Fprintln(stderr, "usage: "+usage.WorktreeCreate)
			return 2
		}
		if args[0] == "--request" {
			request = args[1]
		} else {
			label = args[1]
		}
		args = args[2:]
	}
	creation, err := Create(root, request, label, nil, startRef)
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
	return 0
}
func Subshell(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
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
	creation, err := Create(root, request, objective, nil, startRef)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stderr, "🪵 worktree: %s  (exit to release)\n", creation.Path)
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}
	cmd := exec.Command(shell)
	cmd.Dir, cmd.Stdin, cmd.Stdout, cmd.Stderr = creation.Path, stdin, stdout, stderr
	_ = cmd.Run()
	return ReleaseCommand(root, []string{"--request", request, creation.Path}, io.Discard, stderr)
}
func cleanupOutputSafe(value string) bool { return toon.Representable(value) }
func cleanupOutputValue(value string) string {
	if cleanupOutputSafe(value) {
		return value
	}
	return "sha256:" + textDigest(value)
}
