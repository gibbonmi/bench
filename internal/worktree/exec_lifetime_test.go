package worktree

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestExecNormalPipe covers PL81. A descendant holds the child's stdout after the child
// exits. Exec returns after its output delay with the output that arrived and the child's
// zero, and it reports the forced stream on stderr (BO69). The descendant survives the
// return, because a normal exit sends no signal to the group.
func TestExecNormalPipe(t *testing.T) {
	t.Parallel()
	pidFile := filepath.Join(t.TempDir(), "descendant")
	var stdout, stderr bytes.Buffer
	code := runWorktreeChild([]string{"sh", "-c", `sleep 30 & echo $! > "$1"; echo early`, "sh", pidFile}, t.TempDir(), t.TempDir(), os.Environ(), nil, nil, &stdout, &stderr)
	raw, err := os.ReadFile(pidFile)
	requireTest(t, err == nil, "read the descendant PID: %v", err)
	pid, err := strconv.Atoi(strings.TrimSpace(string(raw)))
	requireTest(t, err == nil && pid > 0, "descendant PID %q: %v", raw, err)
	alive := syscall.Kill(pid, 0) == nil
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
	requireTest(t, alive, "the descendant %d died at the child's normal exit, so exec added a group kill", pid)
	requireTest(t, stdout.String() == "early\n", "exec stdout = %q, want the output that arrived before the delay", stdout.String())
	requireTest(t, code == 0 && strings.Contains(stderr.String(), string(subprocess.StreamClosed)),
		"exec exit = %d with stderr %q, want the child's 0 and a line naming the forcibly closed stream", code, stderr.String())
}

// refusingWriter fails every write, so the stream that copies into it ends failed.
type refusingWriter struct{}

func (refusingWriter) Write([]byte) (int, error) { return 0, errors.New("caller stdout refused") }

// TestExecUnfinishedChildZero covers ticket 11's rule that exec returns child zero only
// after its cleanup obligations complete. The child exits zero, but the caller's stdout
// refuses its output, so exec's owned stdout stream fails rather than completing. Exec
// returns 1 and names the failed stream on stderr.
func TestExecUnfinishedChildZero(t *testing.T) {
	t.Parallel()
	var stderr bytes.Buffer
	code := runWorktreeChild([]string{"sh", "-c", "printf out"}, t.TempDir(), t.TempDir(), os.Environ(), nil, nil, refusingWriter{}, &stderr)
	requireTest(t, code == 1, "exec exit = %d with stderr %q, want 1 for a child zero with a failed stream", code, stderr.String())
	requireTest(t, strings.Contains(stderr.String(), "process stream "+string(subprocess.StreamFailed)),
		"exec stderr = %q, want the line naming the failed stream", stderr.String())
}

// testVector is a well-formed descriptor vector that names no published store. A test binary
// installs no protocol, so the vector reaches the child unchanged.
func testVector(t *testing.T, name string) []protection.Descriptor {
	t.Helper()
	dir := t.TempDir()
	return []protection.Descriptor{{
		Store: filepath.Join(dir, "store"), ResourceID: strings.Repeat("ab", 16), Kind: protection.KindWorktree,
		ResourcePath: filepath.Join(dir, name), Binding: "1:2",
	}}
}

func vectorEntry(t *testing.T, vector []protection.Descriptor) string {
	t.Helper()
	data, err := json.Marshal(vector)
	requireTest(t, err == nil, "encode the vector: %v", err)
	return protection.Env + "=" + string(data)
}

// childVector runs one exec child under environ and extra and returns the vector it read.
func childVector(t *testing.T, environ, extra []string) []protection.Descriptor {
	t.Helper()
	worktree, _ := worktreeWithWrapper(t)
	var stdout, stderr bytes.Buffer
	code := runWorktreeChild([]string{"sh", "-c", `printf %s "${` + protection.Env + `-}"`}, worktree, t.TempDir(), environ, extra, nil, &stdout, &stderr)
	requireTest(t, code == 0, "exec exit = %d, stderr %q", code, stderr.String())
	vector, err := protection.ParseVector(stdout.String())
	requireTest(t, err == nil, "child vector %q: %v", stdout.String(), err)
	return vector
}

// TestExecDescriptorEnvironment covers PL114. The routing transforms strip the wrapper
// routing, replace the home, and set PWD and the wrapper marker, and the inherited vector
// still reaches the child whole. A malformed vector prevents the launch.
func TestExecDescriptorEnvironment(t *testing.T) {
	t.Parallel()
	vector := testVector(t, "inherited")
	environ := append(os.Environ(), "BENCH_KIT=/caller/kit", "BENCH_WRAPPER=/caller/bin/bench", vectorEntry(t, vector))
	got := childVector(t, environ, nil)
	requireTest(t, reflect.DeepEqual(got, vector), "child vector = %#v, want the inherited %#v", got, vector)

	marker := filepath.Join(t.TempDir(), "started")
	var stdout, stderr bytes.Buffer
	code := runWorktreeChild([]string{"touch", marker}, t.TempDir(), t.TempDir(), append(os.Environ(), protection.Env+"=not-json"), nil, nil, &stdout, &stderr)
	_, statErr := os.Stat(marker)
	requireTest(t, code == 1 && os.IsNotExist(statErr), "a malformed vector exited %d and left marker error %v, want 1 and no start", code, statErr)
}

// TestExecDescriptorOverride covers PL115. An --env value cannot replace the inherited
// vector, and it cannot add one where the caller inherited none.
func TestExecDescriptorOverride(t *testing.T) {
	t.Parallel()
	vector, replacement := testVector(t, "inherited"), testVector(t, "replacement")
	got := childVector(t, append(os.Environ(), vectorEntry(t, vector)), []string{vectorEntry(t, replacement)})
	requireTest(t, reflect.DeepEqual(got, vector), "child vector = %#v, want the inherited %#v", got, vector)
	environ := capability.WithoutEnvironment(os.Environ(), protection.Env)
	got = childVector(t, environ, []string{vectorEntry(t, replacement)})
	requireTest(t, len(got) == 0, "child vector = %#v, want none from an --env value", got)
}

// TestCheckoutDescriptor covers the worktree identity adapter that the CLI table installs.
// The descriptor names the checkout's Git administrative directory, which reset and clean
// never reach, and the validator refuses a binding that names another directory.
func TestCheckoutDescriptor(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	d, err := Descriptor(root)
	requireTest(t, err == nil, "describe the checkout: %v", err)
	admin, err := canonicalPath(filepath.Join(root, ".git"))
	requireTest(t, err == nil && d.Kind == protection.KindWorktree && d.ResourcePath == admin, "descriptor = %#v, want the worktree kind at %s", d, admin)
	requireTest(t, ValidateProtection(d) == nil, "the validator refused its own descriptor %#v", d)
	d.Binding = "0:0"
	requireTest(t, ValidateProtection(d) != nil, "the validator accepted a foreign binding")
}
