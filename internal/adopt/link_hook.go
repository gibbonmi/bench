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

func installGitHook(root string, stderr io.Writer) error {
	if gitOK("-C", root, "remote", "get-url", "origin") &&
		!gitOK("-C", root, "symbolic-ref", "--quiet", "refs/remotes/origin/HEAD") {
		_ = exec.Command("git", "-C", root, "remote", "set-head", "origin", "--auto").Run()
	}
	def := git.DefaultBranch(root)
	hooks := hooksDir(root)
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		fmt.Fprintln(stderr, err)
		return err
	}
	prepush := filepath.Join(hooks, "pre-push")
	if content, err := os.ReadFile(prepush); err == nil && !strings.Contains(string(content), "bench:managed-pre-push") {
		fmt.Fprintf(stderr, "conflict: %s exists and is not Bench-managed\n", prepush)
		return fmt.Errorf("foreign pre-push")
	}
	text := strings.ReplaceAll(prePushTemplate, prePushBranchToken, def)
	return os.WriteFile(prepush, []byte(text), 0o755)
}

func gitOK(args ...string) bool {
	return exec.Command("git", args...).Run() == nil
}
