package guards

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gittest"
)

// familySGolden is the reference capture of scan family S. It was taken on the worker
// pipeline before the synchronous loop replaced it, and it is reviewed apart from the
// candidate code, so the candidate is never its oracle.
const familySGolden = "testdata/scan-family-s.golden"

const familySHeader = "=== member: "

// cancelTick is the scan window of a member whose operation waits for cancellation.
const cancelTick = 10 * time.Millisecond

// operationHold is how long an injected operation ignores cancellation before it returns.
const operationHold = 20 * time.Millisecond

// stubSeams replaces the scan seams for one test; a nil seam keeps the production one.
// Cleanup restores both seams and the scan window, so a test can also set the window.
func stubSeams(t *testing.T, enumerate func(context.Context, string) ([]candidate, error), inspect func(context.Context, string, candidate) [][]string) {
	t.Helper()
	oldEnumerate, oldInspect, oldTimeout := enumerateGuards, inspectGuard, guardScanTimeout
	t.Cleanup(func() { enumerateGuards, inspectGuard, guardScanTimeout = oldEnumerate, oldInspect, oldTimeout })
	if enumerate != nil {
		enumerateGuards = enumerate
	}
	if inspect != nil {
		inspectGuard = inspect
	}
}

func cannedRow(name, currency, wired string) []string {
	return []string{name, "boundary of " + name, "denies of " + name, "", "", currency, wired}
}

func rowsOf(name string) [][]string { return [][]string{cannedRow(name, "current", "claude")} }

func candidatesNamed(names ...string) []candidate {
	candidates := make([]candidate, len(names))
	for i, name := range names {
		candidates[i] = candidate{fallback: name}
	}
	return candidates
}

func enumerateNames(names ...string) func(context.Context, string) ([]candidate, error) {
	return func(context.Context, string) ([]candidate, error) { return candidatesNamed(names...), nil }
}

// inspectByName answers a canned candidate by its name. "zero" has no rows, and "multi" has
// a stale row and an unwired row. A "blocked-" candidate waits for cancellation, and
// "blocked-holding" then ignores it for operationHold. Every other name has one row.
func inspectByName(ctx context.Context, _ string, c candidate) [][]string {
	switch c.fallback {
	case "zero":
		return nil
	case "multi":
		return [][]string{cannedRow("multi-stale", "stale", "claude"), cannedRow("multi-unwired", "", "none")}
	case "blocked-empty":
		<-ctx.Done()
		return nil
	case "blocked-rows":
		<-ctx.Done()
	case "blocked-holding":
		<-ctx.Done()
		time.Sleep(operationHold)
	}
	return rowsOf(c.fallback)
}

func manifest(name string) string {
	return "#!/usr/bin/env bash\n# name: " + name + "\n# boundary: PreToolUse:Bash\n# denies: " + name + " denials\n# why: w\nexit 0\n"
}

// plantFiles writes each slash-relative path under root, with its parent directories.
func plantFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, body := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

// familyMember is one family S event script. arrange prepares the seams, the scan window,
// and a working directory, and returns that directory. Each argv in args runs through
// Command there, and a member with scan set also records Scan over the same window.
type familyMember struct {
	name    string
	args    [][]string
	scan    bool
	arrange func(t *testing.T) string
}

var commandForms = [][]string{nil, {"--brief"}}

// seamMember drives Command through controlled seams in a fresh repository. A nonzero
// window replaces the scan window.
func seamMember(name string, window time.Duration, enumerate func(context.Context, string) ([]candidate, error)) familyMember {
	return familyMember{name: name, args: commandForms, scan: true, arrange: func(t *testing.T) string {
		stubSeams(t, enumerate, inspectByName)
		if window != 0 {
			guardScanTimeout = window
		}
		return gittest.Repo(t)
	}}
}

// realMember drives Command through the production seams over a repository that plant
// prepares.
func realMember(name string, plant func(t *testing.T, root string)) familyMember {
	return familyMember{name: name, args: commandForms, scan: true, arrange: func(t *testing.T) string {
		stubSeams(t, nil, nil)
		root := gittest.Repo(t)
		plant(t, root)
		return root
	}}
}

// familyS lists the deterministic scan states. A cancellation member makes its operation
// wait for the scan window, so the event order is fixed. Cancellation between two
// candidates is not here, because no Command event places it there; a direct Scan test
// owns it.
var familyS = []familyMember{
	seamMember("complete", 0, enumerateNames("zero", "one", "multi")),
	realMember("pre-cancelled", func(*testing.T, string) { guardScanTimeout = 0 }),
	seamMember("enumeration-cancelled-error", cancelTick, func(ctx context.Context, _ string) ([]candidate, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}),
	seamMember("enumeration-cancelled-candidates", cancelTick, func(ctx context.Context, _ string) ([]candidate, error) {
		<-ctx.Done()
		return candidatesNamed("late"), nil
	}),
	seamMember("inspection-cancelled-empty", cancelTick, enumerateNames("one", "blocked-empty")),
	seamMember("inspection-cancelled-rows", cancelTick, enumerateNames("one", "blocked-rows")),
	seamMember("blocking-operation", cancelTick, enumerateNames("blocked-holding")),
	seamMember("enumeration-error-candidates", 0, func(context.Context, string) ([]candidate, error) {
		return candidatesNamed("one"), errors.New("enumeration failed")
	}),
	seamMember("enumeration-error-empty", 0, func(context.Context, string) ([]candidate, error) {
		return nil, errors.New("enumeration failed")
	}),
	realMember("hooks-absent", func(*testing.T, string) {}),
	realMember("hooks-empty", func(t *testing.T, root string) {
		if err := os.MkdirAll(filepath.Join(root, ".bench", "hooks"), 0o755); err != nil {
			t.Fatal(err)
		}
	}),
	realMember("static-headers", func(t *testing.T, root string) {
		plantFiles(t, root, map[string]string{
			".bench/hooks/full.sh":       manifest("full-guard"),
			".bench/hooks/incomplete.sh": "#!/usr/bin/env bash\n# name: partial\n# boundary: b\n# denies: d\nexit 0\n",
			".bench/hooks/headerless.sh": "#!/usr/bin/env bash\nexit 0\n",
			".bench/hooks/info.sh":       "#!/usr/bin/env bash\n# name: info\n# boundary: SessionStart\n# denies: nothing (informational)\n# why: w\nexit 0\n",
			".bench/hooks/notes.txt":     manifest("not-a-hook"),
		})
	}),
	realMember("symlink-and-missing-headers", func(t *testing.T, root string) {
		plantFiles(t, root, map[string]string{".bench/hooks/full.sh": manifest("full-guard")})
		for link, target := range map[string]string{"link.sh": "full.sh", "missing.sh": "absent-target.sh"} {
			if err := os.Symlink(target, filepath.Join(root, ".bench", "hooks", link)); err != nil {
				capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
			}
		}
	}),
	realMember("fifo-header", func(t *testing.T, root string) {
		if runtime.GOOS == "windows" {
			capability.Capability(t, capability.Fifo, "named pipes use the Unix mkfifo fixture")
		}
		hooks := filepath.Join(root, ".bench", "hooks")
		if err := os.MkdirAll(hooks, 0o755); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command("mkfifo", filepath.Join(hooks, "special.sh")).CombinedOutput(); err != nil {
			t.Fatalf("mkfifo: %v: %s", err, out)
		}
	}),
	realMember("wiring-claude-and-codex", func(t *testing.T, root string) {
		plantFiles(t, root, map[string]string{
			".bench/hooks/a-wired.sh":   manifest("a-wired"),
			".bench/hooks/z-unwired.sh": manifest("z-unwired"),
			".claude/settings.json":     `{"hook":".bench/hooks/a-wired.sh"}`,
			".codex/hooks.json":         `{"hook":".bench/hooks/a-wired.sh"}`,
		})
	}),
	realMember("wiring-empty-and-invalid", func(t *testing.T, root string) {
		plantFiles(t, root, map[string]string{
			".bench/hooks/a-wired.sh": manifest("a-wired"),
			".claude/settings.json":   "",
			".codex/hooks.json":       `{"hook":".bench/hooks/a-wired.sh"`,
		})
	}),
	{name: "invalid-arguments", args: [][]string{{"--bogus"}, {"--brief", "extra"}}, arrange: func(t *testing.T) string { return gittest.Repo(t) }},
	{name: "outside-repository", args: commandForms, arrange: func(t *testing.T) string { return t.TempDir() }},
}

// commandCapturingStderr runs Command and returns its stdout, its exit, and the bytes it
// wrote to the process stderr.
func commandCapturingStderr(t *testing.T, args []string) (string, int, string) {
	t.Helper()
	sink, err := os.CreateTemp(t.TempDir(), "stderr")
	if err != nil {
		t.Fatal(err)
	}
	defer sink.Close()
	saved := os.Stderr
	os.Stderr = sink
	out, code := Command(args)
	os.Stderr = saved
	written, err := os.ReadFile(sink.Name())
	if err != nil {
		t.Fatal(err)
	}
	return out, code, string(written)
}

// captureMember renders one member's Command tuples and ScanResult. It records no elapsed
// time.
func captureMember(t *testing.T, member familyMember) string {
	dir := member.arrange(t)
	t.Chdir(dir)
	var b strings.Builder
	for _, args := range member.args {
		out, code, stderr := commandCapturingStderr(t, args)
		fmt.Fprintf(&b, "--- args %q exit %d\n%s--- stderr\n%s", args, code, out, stderr)
	}
	if member.scan {
		ctx, cancel := bounds.Context(context.Background(), guardScanTimeout)
		defer cancel()
		fmt.Fprintf(&b, "--- scan\n%#v\n", Scan(ctx, dir))
	}
	return b.String()
}

// TestScanLoopDifferentialCommandFamilies compares the candidate's family S tuples with the
// reference capture, member by member, through the real Command formatter.
func TestScanLoopDifferentialCommandFamilies(t *testing.T) {
	body, err := os.ReadFile(familySGolden)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{}
	var captured, members []string
	for _, section := range strings.Split(string(body), familySHeader)[1:] {
		name, capture, _ := strings.Cut(section, "\n")
		want[name] = capture
		captured = append(captured, name)
	}
	for _, member := range familyS {
		members = append(members, member.name)
	}
	if !reflect.DeepEqual(captured, members) {
		t.Fatalf("%s holds members %q, want family S %q", familySGolden, captured, members)
	}
	for _, member := range familyS {
		t.Run(member.name, func(t *testing.T) {
			if got := captureMember(t, member); got != want[member.name] {
				t.Errorf("member %s differs from %s\ngot:\n%s\nwant:\n%s", member.name, familySGolden, got, want[member.name])
			}
		})
	}
}
