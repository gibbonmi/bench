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
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gittest"
)

const (
	// familySGolden is the reference capture of scan family S.
	familySGolden = "testdata/scan-family-s.golden"
	familySHeader = "=== member: "
	// cancelTick is the scan window of a member whose operation waits for cancellation.
	cancelTick = 10 * time.Millisecond
	// inspectionWindow outlasts any plausible stall, so only a blocked second inspection meets its end.
	inspectionWindow = time.Second
	// operationHold is how long an injected operation ignores cancellation before it returns.
	operationHold = 20 * time.Millisecond
)

// stubSeams replaces the scan seams for one test. Cleanup restores both seams and the scan
// window, so a test can also set the window.
func stubSeams(t *testing.T, enumerate func(context.Context, string) ([]candidate, error), inspect func(context.Context, string, candidate) [][]string) {
	t.Helper()
	oldEnumerate, oldInspect, oldTimeout := enumerateGuards, inspectGuard, guardScanTimeout
	t.Cleanup(func() { enumerateGuards, inspectGuard, guardScanTimeout = oldEnumerate, oldInspect, oldTimeout })
	enumerateGuards, inspectGuard = enumerate, inspect
}

func cannedRow(name, currency, wired string) []string {
	return []string{name, "boundary of " + name, "denies of " + name, "", "", currency, wired}
}

func rowsOf(name string) [][]string { return [][]string{cannedRow(name, "current", "claude")} }

func candidatesNamed(names ...string) (candidates []candidate) {
	for _, name := range names {
		candidates = append(candidates, candidate{fallback: name})
	}
	return candidates
}

func enumerateNames(names ...string) func(context.Context, string) ([]candidate, error) {
	return func(context.Context, string) ([]candidate, error) { return candidatesNamed(names...), nil }
}

// inspectByName answers by name: "zero" has no rows, "multi" a stale and an unwired row, a
// "blocked-" name waits for cancellation, and "blocked-holding" then ignores it for operationHold.
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

// familyMember is one family S event script. arrange returns the working directory where each
// argv in args runs through Command; scan also records Scan over the same window.
type familyMember struct {
	name    string
	args    [][]string
	scan    bool
	arrange func(t *testing.T) string
}

var commandForms = [][]string{nil, {"--brief"}}

// seamMember drives Command through controlled seams and window in a fresh repository.
func seamMember(name string, window time.Duration, enumerate func(context.Context, string) ([]candidate, error)) familyMember {
	return familyMember{name: name, args: commandForms, scan: true, arrange: func(t *testing.T) string {
		stubSeams(t, enumerate, inspectByName)
		guardScanTimeout = window
		return gittest.Repo(t)
	}}
}

// realMember drives Command through the production seams over a repository that plant prepares.
func realMember(name string, plant func(t *testing.T, root string)) familyMember {
	return familyMember{name: name, args: commandForms, scan: true, arrange: func(t *testing.T) string {
		stubSeams(t, enumerateGuards, inspectGuard)
		root := gittest.Repo(t)
		plant(t, root)
		return root
	}}
}

// familyS lists the deterministic scan states; each cancellation waits for the scan window.
// No Command event places cancellation between two candidates, so a direct Scan test owns it.
var familyS = []familyMember{
	seamMember("complete", guardScanTimeout, enumerateNames("zero", "one", "multi")),
	realMember("pre-cancelled", func(*testing.T, string) { guardScanTimeout = 0 }),
	seamMember("enumeration-cancelled-error", cancelTick, func(ctx context.Context, _ string) ([]candidate, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}),
	seamMember("enumeration-cancelled-candidates", cancelTick, func(ctx context.Context, _ string) ([]candidate, error) {
		<-ctx.Done()
		return candidatesNamed("late"), nil
	}),
	seamMember("inspection-cancelled-empty", inspectionWindow, enumerateNames("one", "blocked-empty")),
	seamMember("inspection-cancelled-rows", inspectionWindow, enumerateNames("one", "blocked-rows")),
	seamMember("blocking-operation", cancelTick, enumerateNames("blocked-holding")),
	seamMember("enumeration-error-candidates", guardScanTimeout, func(context.Context, string) ([]candidate, error) {
		return candidatesNamed("one"), errors.New("enumeration failed")
	}),
	seamMember("enumeration-error-empty", guardScanTimeout, func(context.Context, string) ([]candidate, error) {
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
				capability.Unavailable(t, capability.Symlink, err.Error())
			}
		}
	}),
	realMember("fifo-header", func(t *testing.T, root string) {
		if runtime.GOOS == "windows" {
			capability.Unavailable(t, capability.Fifo, "named pipes use the Unix mkfifo fixture")
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

// commandCapturingStderr runs Command and returns its stdout, its exit, and its stderr bytes.
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

// captureMember renders one member's Command tuples and ScanResult, with no elapsed time.
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

// TestScanLoopDifferentialCommandFamilies compares family S Command tuples with the reference.
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

// checkpointContext cancels itself at its cancelAt-th Err read (never when zero). Scan reads
// ctx before and after enumeration and before and after each inspection.
type checkpointContext struct {
	context.Context
	cancel          context.CancelFunc
	reads, cancelAt int
}

func (c *checkpointContext) Err() error {
	c.reads++
	if c.reads == c.cancelAt {
		c.cancel()
	}
	return c.Context.Err()
}

// scanCase scripts one direct Scan over inspectByName candidates. cancelOn names the operation
// that cancels as it runs ("start" cancels before Scan); ops lists the operations Scan starts.
type scanCase struct {
	names, ops []string
	enumErr    error
	cancelOn   string
	cancelAt   int
	want       ScanResult
}

func runScan(t *testing.T, c scanCase) {
	t.Helper()
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()
	var ops []string
	stubSeams(t, func(context.Context, string) ([]candidate, error) {
		ops = append(ops, "enumerate")
		if c.cancelOn == "enumerate" {
			cancel()
		}
		return candidatesNamed(c.names...), c.enumErr
	}, func(ctx context.Context, root string, candidate candidate) [][]string {
		ops = append(ops, "inspect "+candidate.fallback)
		if c.cancelOn == "inspect "+candidate.fallback {
			cancel()
		}
		return inspectByName(ctx, root, candidate)
	})
	if c.cancelOn == "start" {
		cancel()
	}
	got := Scan(&checkpointContext{Context: parent, cancel: cancel, cancelAt: c.cancelAt}, "/repo")
	if !reflect.DeepEqual(got, c.want) || !reflect.DeepEqual(ops, c.ops) {
		t.Fatalf("Scan = %#v after operations %q\nwant %#v after %q", got, ops, c.want, c.ops)
	}
}

var enumerationInterrupted = ScanResult{Status: "incomplete", Inspected: "0", Total: "unknown", Omitted: "unknown", Reason: "timeout"}

func TestScanLoopCompleteOrder(t *testing.T) {
	multi := inspectByName(context.Background(), "", candidate{fallback: "multi"})
	runScan(t, scanCase{names: []string{"a", "multi", "b"}, ops: []string{"enumerate", "inspect a", "inspect multi", "inspect b"},
		want: ScanResult{Rows: slices.Concat(rowsOf("a"), multi, rowsOf("b")), Status: "complete", Inspected: "3", Total: "3", Omitted: "0", Reason: "none"}})
}

func TestScanLoopPreCancelled(t *testing.T) {
	runScan(t, scanCase{names: []string{"a"}, cancelOn: "start", want: enumerationInterrupted})
}

func TestScanLoopEnumerationReturnsAfterCancellation(t *testing.T) {
	runScan(t, scanCase{names: []string{"a", "b"}, cancelOn: "enumerate", ops: []string{"enumerate"}, want: enumerationInterrupted})
}

// The fifth context read is the checkpoint before the second inspection.
func TestScanLoopBetweenCandidates(t *testing.T) {
	runScan(t, scanCase{names: []string{"a", "b"}, cancelAt: 5, ops: []string{"enumerate", "inspect a"},
		want: ScanResult{Rows: rowsOf("a"), Status: "incomplete", Inspected: "1", Total: "2", Omitted: "1", Reason: "timeout"}})
}

func TestScanLoopEmptyInspectionCounts(t *testing.T) {
	runScan(t, scanCase{names: []string{"zero"}, ops: []string{"enumerate", "inspect zero"},
		want: ScanResult{Status: "complete", Inspected: "1", Total: "1", Omitted: "0", Reason: "none"}})
	runScan(t, scanCase{names: []string{"zero", "b"}, cancelOn: "inspect b", ops: []string{"enumerate", "inspect zero", "inspect b"},
		want: ScanResult{Status: "incomplete", Inspected: "1", Total: "2", Omitted: "1", Reason: "timeout"}})
}

func TestScanLoopPostInspectionCancellation(t *testing.T) {
	runScan(t, scanCase{names: []string{"a", "b"}, cancelOn: "inspect b", ops: []string{"enumerate", "inspect a", "inspect b"},
		want: ScanResult{Rows: rowsOf("a"), Status: "incomplete", Inspected: "1", Total: "2", Omitted: "1", Reason: "timeout"}})
}

// Every field of a complete, an enumeration-cancelled, and an inspection-cancelled result.
func TestScanLoopResultMetadata(t *testing.T) {
	runScan(t, scanCase{names: []string{"a"}, ops: []string{"enumerate", "inspect a"},
		want: ScanResult{Rows: rowsOf("a"), Status: "complete", Inspected: "1", Total: "1", Omitted: "0", Reason: "none"}})
	runScan(t, scanCase{names: []string{"a"}, enumErr: context.Canceled, cancelOn: "enumerate", ops: []string{"enumerate"}, want: enumerationInterrupted})
	runScan(t, scanCase{names: []string{"a"}, cancelOn: "inspect a", ops: []string{"enumerate", "inspect a"},
		want: ScanResult{Status: "incomplete", Inspected: "0", Total: "1", Omitted: "1", Reason: "timeout"}})
}

func TestScanLoopEnumerationError(t *testing.T) {
	failed := errors.New("enumeration failed")
	runScan(t, scanCase{names: []string{"a"}, enumErr: failed, ops: []string{"enumerate", "inspect a"},
		want: ScanResult{Rows: rowsOf("a"), Status: "complete", Inspected: "1", Total: "1", Omitted: "0", Reason: "none"}})
	runScan(t, scanCase{enumErr: failed, ops: []string{"enumerate"}, want: ScanResult{Status: "complete", Inspected: "0", Total: "0", Omitted: "0", Reason: "none"}})
}

// Two candidates make six checkpoint reads, so only a seventh read sees cancellation.
func TestScanLoopAfterFullCompletion(t *testing.T) {
	runScan(t, scanCase{names: []string{"a", "b"}, cancelAt: 7, ops: []string{"enumerate", "inspect a", "inspect b"},
		want: ScanResult{Rows: slices.Concat(rowsOf("a"), rowsOf("b")), Status: "complete", Inspected: "2", Total: "2", Omitted: "0", Reason: "none"}})
}

// The inspection ignores the expired window for operationHold; Scan must return after it.
func TestScanLoopBlockingOperationHasNoHardDeadline(t *testing.T) {
	var returned atomic.Bool
	stubSeams(t, enumerateNames("blocked-holding"), func(ctx context.Context, root string, c candidate) [][]string {
		defer returned.Store(true)
		return inspectByName(ctx, root, c)
	})
	ctx, cancel := context.WithTimeout(context.Background(), cancelTick)
	defer cancel()
	got := Scan(ctx, "/repo")
	if !returned.Load() || !reflect.DeepEqual(got, ScanResult{Status: "incomplete", Inspected: "0", Total: "1", Omitted: "1", Reason: "timeout"}) {
		t.Fatalf("Scan = %#v, operation returned first = %v; want an incomplete scan after the operation returns", got, returned.Load())
	}
}

func TestScanLoopRowsAdapter(t *testing.T) {
	stubSeams(t, enumerateNames("a", "multi"), inspectByName)
	want := Scan(context.Background(), "/repo").Rows
	if got := Rows("/repo"); len(want) != 3 || !reflect.DeepEqual(got, want) {
		t.Fatalf("Rows = %q, want Scan(background).Rows %q", got, want)
	}
}
