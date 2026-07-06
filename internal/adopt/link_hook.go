package adopt

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

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
	text := fmt.Sprintf(`#!/usr/bin/env bash
# bench:managed-pre-push
# Installed by 'bench link'. The merge is the human's; agents don't push %[1]s.
if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: pre-push\n'
  printf 'boundary: pre-push\n'
  printf 'denies: direct push to %[1]s, .bench drift from bench gate pin\n'
  printf 'why: the merge belongs to the reviewer; agents open a PR instead of pushing %[1]s\n'
  exit 0
fi
pin_path="$(git rev-parse --git-path bench-gate-pin 2>/dev/null || true)"
pin_tree=""
if [[ -n "$pin_path" && -f "$pin_path" ]]; then
  IFS= read -r pin_tree < "$pin_path" || true
fi
if [[ -z "$pin_tree" ]]; then
  echo "bench: gate unpinned - run 'bench gate pin' to enable .bench drift checks." >&2
fi
ref_lines=()
while IFS= read -r line; do
  ref_lines+=("$line")
  read -r local_ref local_oid remote_ref remote_oid <<< "$line"
  if [[ -n "$pin_tree" && ! "$local_oid" =~ ^0+$ ]]; then
    if ! bench_tree="$(git rev-parse "$local_oid:.bench" 2>/dev/null)"; then
      echo "blocked: pushed commit has no .bench tree. Review the gate change, then run 'bench gate pin'." >&2
      exit 1
    fi
    if [[ "$bench_tree" != "$pin_tree" ]]; then
      echo "blocked: pushed .bench tree does not match bench gate pin. Review the gate change, then run 'bench gate pin'." >&2
      exit 1
    fi
  fi
done
for line in "${ref_lines[@]}"; do
  read -r local_ref local_oid remote_ref remote_oid <<< "$line"
  if [[ "$remote_ref" == "refs/heads/%[1]s" ]]; then
    echo "blocked: direct push to %[1]s. Open a PR or merge it yourself." >&2
    exit 1
  fi
done
exit 0
`, def)
	return os.WriteFile(prepush, []byte(text), 0o755)
}

func gitOK(args ...string) bool {
	return exec.Command("git", args...).Run() == nil
}
