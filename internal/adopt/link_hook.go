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
// lintable shell; the shellcheck gate phase lints prepush.sh at its kit-tree path.
// renderPrePush substitutes the default-branch token before a managed hook is installed.
//
//go:embed prepush.sh
var prePushTemplate string

// prePushBranchToken is the placeholder the asset carries in place of the repo's
// default branch, substituted at install time.
const prePushBranchToken = "__BENCH_DEFAULT_BRANCH__"

// PrePushMarker is the fingerprint of a bench-managed pre-push hook: the marker line the
// template carries. InspectPrePush's state classification, installGitHook's conflict
// check, and unlink's removal gate all match this one source by substring, not by
// byte-identity. Byte-identity would false-red across default-branch token substitution
// and benign template evolution.
const PrePushMarker = "bench:managed-pre-push"

// PrePushState classifies a repo's pre-push hook against the bench-managed template.
type PrePushState int

const (
	PrePushManaged  PrePushState = iota // present and carries the managed marker
	PrePushAbsent                       // no hook at the default hooks dir (a fresh clone drops it)
	PrePushForeign                      // present in the default hooks dir but not bench-managed
	PrePushDiverted                     // core.hooksPath points elsewhere, with no managed hook there
)

// PrePushProvenance identifies whether the protected branch came from origin/HEAD or the hook.
type PrePushProvenance string

const (
	PrePushLive  PrePushProvenance = "live"
	PrePushBaked PrePushProvenance = "baked"
)

// PrePushCurrency reports whether a managed hook matches the embedded template.
type PrePushCurrency string

const (
	PrePushCurrent       PrePushCurrency = "current"
	PrePushStale         PrePushCurrency = "stale"
	PrePushNotApplicable PrePushCurrency = "not-applicable"
)

// PrePushHealth is the complete local record for the effective pre-push hook.
type PrePushHealth struct {
	State      PrePushState
	Path       string
	Branch     string
	Provenance PrePushProvenance
	Fallback   bool
	Currency   PrePushCurrency
}

// InspectPrePush reads the effective hook without contacting a remote and returns its health.
func InspectPrePush(root string) PrePushHealth {
	path := filepath.Join(hooksDir(root), "pre-push")
	if _, err := os.Lstat(path); err != nil {
		return noPrePushHealth(root, path)
	}
	// The mode gate resolves through symlinks before any read. os.ReadFile follows a link,
	// and opening a writerless FIFO never returns. A link pointing at one would wedge every
	// ambient surface that asks for hook health. A link resolving to nothing is foreign on
	// the same terms as a path that exists but cannot be read.
	target, err := os.Stat(path)
	if err != nil || !target.Mode().IsRegular() {
		return PrePushHealth{State: PrePushForeign, Path: path, Currency: PrePushNotApplicable}
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return PrePushHealth{State: PrePushForeign, Path: path, Currency: PrePushNotApplicable}
	}
	if !strings.Contains(string(content), PrePushMarker) {
		return unmanagedPrePushHealth(root, path)
	}
	baked, parsed := prePushBakedBranch(string(content))
	branch, live := prePushLiveBranch(root)
	if !live {
		branch = baked
	}
	currency := PrePushStale
	if parsed && string(content) == renderPrePushBranch(baked) {
		currency = PrePushCurrent
	}
	provenance := PrePushBaked
	if live {
		provenance = PrePushLive
	}
	return PrePushHealth{
		State:      PrePushManaged,
		Path:       path,
		Branch:     branch,
		Provenance: provenance,
		Fallback:   !live && baked == fallbackProtectedBranch,
		Currency:   currency,
	}
}

// prePushRefreshEligible reports whether transactional linking would stage different
// effective pre-push bytes during an upgrade. A configured hooksPath remains eligible
// only when its effective path is absent. Foreign, present-diverted, dangling, and
// special-file paths retain transactionalLink's refusal.
func prePushRefreshEligible(root string) bool {
	health := InspectPrePush(root)
	if health.State == PrePushAbsent {
		return true
	}
	if health.State == PrePushManaged {
		installed, err := os.ReadFile(health.Path)
		return err == nil && string(installed) != string(renderedPrePush(root))
	}
	if health.State != PrePushDiverted {
		return false
	}
	_, err := os.Lstat(health.Path)
	return os.IsNotExist(err)
}

func noPrePushHealth(root, path string) PrePushHealth {
	if hooksPathConfigured(root) {
		return PrePushHealth{State: PrePushDiverted, Path: path, Currency: PrePushNotApplicable}
	}
	return PrePushHealth{State: PrePushAbsent, Path: path, Currency: PrePushNotApplicable}
}

func unmanagedPrePushHealth(root, path string) PrePushHealth {
	if hooksPathConfigured(root) {
		return PrePushHealth{State: PrePushDiverted, Path: path, Currency: PrePushNotApplicable}
	}
	return PrePushHealth{State: PrePushForeign, Path: path, Currency: PrePushNotApplicable}
}

func prePushLiveBranch(root string) (string, bool) {
	out, err := git.Output("-C", root, "symbolic-ref", "--quiet", "--short", "refs/remotes/origin/HEAD")
	if err != nil || out == "" {
		return "", false
	}
	return strings.TrimPrefix(out, "origin/"), true
}

func prePushBakedBranch(content string) (string, bool) {
	index := strings.LastIndex(prePushTemplate, prePushBranchToken)
	if index < 0 {
		return "", false
	}
	anchor := `protected="`
	anchorIndex := strings.LastIndex(prePushTemplate[:index], anchor)
	if anchorIndex < 0 {
		return "", false
	}
	contentAnchor := strings.Index(content, anchor)
	if contentAnchor < 0 {
		return "", false
	}
	branchStart := contentAnchor + len(anchor)
	branchEnd := strings.Index(content[branchStart:], `"`)
	if branchEnd < 0 {
		return "", false
	}
	baked := content[branchStart : branchStart+branchEnd]
	prefix := strings.ReplaceAll(prePushTemplate[:index], prePushBranchToken, baked)
	suffix := prePushTemplate[index+len(prePushBranchToken):]
	if !strings.HasPrefix(content, prefix) {
		return "", false
	}
	remainder := strings.TrimPrefix(content, prefix)
	end := strings.Index(remainder, suffix)
	if end < 0 {
		return "", false
	}
	if remainder[:end] != baked {
		return "", false
	}
	return baked, true
}

// hooksPathConfigured reports whether the repo sets a non-empty core.hooksPath, diverting
// hooks away from the default .git/hooks. This is the signal that distinguishes a
// diverted hooks directory from a plain foreign or absent hook in the default location.
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
// repository. The installed hook re-resolves origin/HEAD live on every push and reaches
// the baked token only when that lookup is empty. A guard protecting nothing is the worse
// failure. It shares its spelling with git's default-branch probe candidate and nothing
// else. That candidate is discarded when unconfirmed; this constant stays regardless, so
// collapsing the two would make a change to either fact silently change the other.
const fallbackProtectedBranch = "main"

// protectedBranch names the branch the installed hook refuses a direct push to.
func protectedBranch(root string) string {
	if def, ok := git.ResolvedDefault(root); ok {
		return def
	}
	return fallbackProtectedBranch
}

func installGitHook(root string, stderr io.Writer) error {
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
	return os.WriteFile(prepush, []byte(renderPrePush(root)), 0o755)
}

func populateOriginHead(root string) {
	if gitOK("-C", root, "remote", "get-url", "origin") &&
		!gitOK("-C", root, "symbolic-ref", "--quiet", "refs/remotes/origin/HEAD") {
		_ = exec.Command("git", "-C", root, "remote", "set-head", "origin", "--auto").Run()
	}
}

func gitOK(args ...string) bool {
	return exec.Command("git", args...).Run() == nil
}
