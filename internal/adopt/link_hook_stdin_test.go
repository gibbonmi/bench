// The pre-push stdin family: the rendered hook is driven out of process over the raw
// stdin bodies git can hand it.
package adopt

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runPrePushHookStdin executes the rendered hook with cwd at the repo root over the exact
// stdin body given. It takes the raw body so a case can drive a shape the four-field line
// cannot express: an empty push, an unterminated last line, or more than one ref.
func runPrePushHookStdin(t *testing.T, root, path, stdin string) (error, string) {
	t.Helper()
	command := exec.Command("bash", path)
	command.Dir = root
	command.Stdin = strings.NewReader(stdin)
	var stderr strings.Builder
	command.Stderr = &stderr
	return command.Run(), stderr.String()
}

// TestPrePushHookReadsEveryStdinShape drives the hook read loop with the bodies git can
// hand it. A quiet shape must stay quiet, and a protected ref must block wherever it sits
// in the body, so a rewritten loop that reads only the first line turns this red.
func TestPrePushHookReadsEveryStdinShape(t *testing.T) {
	root := hookTestRepo(t)
	path := filepath.Join(root, ".git", "hooks", "pre-push")
	writeHook(t, path, "main")

	const blocked = "blocked: direct push to main. Open a PR or merge it yourself.\n"
	oid := strings.Repeat("a", 40)
	zero := strings.Repeat("0", 40)

	for _, testCase := range []struct {
		name  string
		stdin string
		block bool
	}{
		{name: "empty stdin", stdin: ""},
		{name: "unterminated topic line", stdin: "refs/heads/topic " + oid + " refs/heads/topic " + oid},
		{name: "protected deletion", stdin: "(delete) " + zero + " refs/heads/main " + oid + "\n", block: true},
		{name: "protected on the second line", stdin: "refs/heads/topic " + oid + " refs/heads/topic " + oid + "\nrefs/heads/main " + oid + " refs/heads/main " + oid + "\n", block: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			err, out := runPrePushHookStdin(t, root, path, testCase.stdin)
			if !testCase.block {
				if err != nil || out != "" {
					t.Fatalf("push = (err %v, stderr %q), want exit 0 and empty stderr", err, out)
				}
				return
			}
			exit, ok := err.(*exec.ExitError)
			if !ok || exit.ExitCode() != 1 || out != blocked {
				t.Fatalf("push = (err %v, stderr %q), want exit 1 and %q", err, out, blocked)
			}
		})
	}
}
