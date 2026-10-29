package gittest

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

// command is the one git process of this package. It runs git with args in root, and
// root reaches git as its own -C operand, never through a shell. Git leaves the current
// working directory in place for an empty root.
func command(root string, args ...string) *exec.Cmd {
	return exec.Command("git", append([]string{"-C", root}, args...)...)
}

// MustRun runs git with args in root and discards its output. It fails the test when git
// exits non-zero or cannot start. A caller that must handle the failure itself calls Run.
func MustRun(t testing.TB, root string, args ...string) {
	t.Helper()
	_ = Output(t, root, args...)
}

// OutputBytes runs git with args in root and returns its stdout byte for byte, with no
// trim and without stderr. It fails the test when git exits non-zero or cannot start.
func OutputBytes(t testing.TB, root string, args ...string) []byte {
	t.Helper()
	out, err := command(root, args...).Output()
	if err != nil {
		var stderr []byte
		if exit := (*exec.ExitError)(nil); errors.As(err, &exit) {
			stderr = exit.Stderr
		}
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, stderr)
	}
	return out
}

// Commit commits the current index of root with message. Each option reaches git commit
// as its own argument. Commit stages nothing, so the caller keeps its own add paths and
// flags. It fails the test when git refuses, such as for an unchanged index without an
// explicit --allow-empty option.
func Commit(t testing.TB, root, message string, options ...string) {
	t.Helper()
	MustRun(t, root, append([]string{"commit", "--quiet", "-m", message}, options...)...)
}
