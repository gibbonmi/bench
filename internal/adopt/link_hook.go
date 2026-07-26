package adopt

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// prePushTemplate is the one source for the managed pre-push hook body. It ships as
// lintable shell (the shellcheck gate phase lints prepush.sh at its kit-tree path);
// installGitHook substitutes the default-branch token to install a byte-identical hook.
//
//go:embed prepush.sh
var prePushTemplate string

// prePushBranchToken is the placeholder the asset carries in place of the repo's
// default branch, substituted at install time.
const prePushBranchToken = "__BENCH_DEFAULT_BRANCH__"

// PrePushMarker is the fingerprint of a bench-managed pre-push hook: the marker line the
// template carries. It is the one source installGitHook (conflict check), ClassifyPrePush
// (doctor/status backstop verification), and unlink's removal gate match by substring —
// not byte-identity, which would false-red across default-branch token substitution and
// benign template evolution. internal/guards also reads it, to detect a bench-managed
// pre-push hook without duplicating the marker literal.
const PrePushMarker = "bench:managed-pre-push"

// PrePushState classifies a repo's pre-push hook against the bench-managed template.
type PrePushState int

const (
	PrePushManaged  PrePushState = iota // present and carries the managed marker
	PrePushAbsent                       // no hook at the default hooks dir (a fresh clone drops it)
	PrePushForeign                      // present in the default hooks dir but not bench-managed
	PrePushDiverted                     // core.hooksPath points elsewhere, with no managed hook there
)

// PrePushStatus is the classifier result: the state and the resolved pre-push path git
// would use for the repo, honoring core.hooksPath — so a red row can name where the hook belongs.
type PrePushStatus struct {
	State PrePushState
	Path  string
}

// ClassifyPrePush resolves the hooks directory git will use for root — honoring
// core.hooksPath — and classifies the pre-push there against the managed marker. It never
// writes: a fresh clone legitimately has no hook (git does not clone hooks), which is the
// absent state doctor and status surface rather than repair, and a foreign hook is reported,
// never overwritten.
func ClassifyPrePush(root string) PrePushStatus {
	path := filepath.Join(hooksDir(root), "pre-push")
	content, err := os.ReadFile(path)
	if err == nil && strings.Contains(string(content), PrePushMarker) {
		return PrePushStatus{State: PrePushManaged, Path: path}
	}
	if hooksPathConfigured(root) {
		return PrePushStatus{State: PrePushDiverted, Path: path}
	}
	if err != nil {
		return PrePushStatus{State: PrePushAbsent, Path: path}
	}
	return PrePushStatus{State: PrePushForeign, Path: path}
}

// hooksPathConfigured reports whether the repo sets a non-empty core.hooksPath, diverting
// hooks away from the default .git/hooks — the signal that distinguishes a diverted hooks
// directory from a plain foreign or absent hook in the default location.
func hooksPathConfigured(root string) bool {
	v, err := git.Output("-C", root, "config", "--get", "core.hooksPath")
	return err == nil && strings.TrimSpace(v) != ""
}

func hooksDir(root string) string {
	out, err := git.Output("-C", root, "rev-parse", "--git-path", "hooks")
	if err != nil || out == "" {
		return filepath.Join(root, ".git", "hooks")
	}
	if filepath.IsAbs(out) {
		return out
	}
	return filepath.Join(root, out)
}

// fallbackProtectedBranch is baked into the pre-push hook when the repository has no
// resolvable default branch. It is the guard's fail-safe rather than a claim about the
// repository: the installed hook re-resolves origin/HEAD live on every push and reaches
// the baked token only when that lookup is empty, and a guard protecting nothing is the
// worse failure. It is not git's object-database probe candidate, which is discarded when
// the repository cannot confirm it; this one is kept precisely because nothing confirmed.
const fallbackProtectedBranch = "main"

// protectedBranch names the branch the installed hook refuses a direct push to.
func protectedBranch(root string) string {
	if def, ok := git.ResolvedDefault(root); ok {
		return def
	}
	return fallbackProtectedBranch
}

func installGitHook(root string, stderr io.Writer) error {
	if gitOK("-C", root, "remote", "get-url", "origin") &&
		!gitOK("-C", root, "symbolic-ref", "--quiet", "refs/remotes/origin/HEAD") {
		_ = exec.Command("git", "-C", root, "remote", "set-head", "origin", "--auto").Run()
	}
	def := protectedBranch(root)
	hooks := hooksDir(root)
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		fmt.Fprintln(stderr, err)
		return err
	}
	prepush := filepath.Join(hooks, "pre-push")
	if content, err := os.ReadFile(prepush); err == nil && !strings.Contains(string(content), PrePushMarker) {
		fmt.Fprintf(stderr, "conflict: %s exists and is not Bench-managed\n", prepush)
		return fmt.Errorf("foreign pre-push")
	}
	text := strings.ReplaceAll(prePushTemplate, prePushBranchToken, def)
	return os.WriteFile(prepush, []byte(text), 0o755)
}

func gitOK(args ...string) bool {
	return exec.Command("git", args...).Run() == nil
}
