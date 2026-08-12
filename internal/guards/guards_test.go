package guards

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/capability"
)

func TestCommandAlwaysEmitsCompleteGuardScanMetadata(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	out, code := Command(nil)
	if code != 0 || !strings.Contains(out, "guard_scan[1]{status,inspected,total,omitted,reason}:") || !strings.Contains(out, "complete") {
		t.Fatalf("code/output = %d\n%s", code, out)
	}
}

func TestActionsForRowsDerivesEveryStaleOrUnwiredRowAndDedupes(t *testing.T) {
	rows := [][]string{
		{"stale", "b", "d", "", "", "stale", "claude"},
		{"unwired", "b", "d", "", "", "", "none"},
		{"stale", "b", "d", "", "", "stale", "claude"},
		{"clean", "b", "d", "", "", "current", "claude"},
		{"wired", "b", "d", "", "", "", "claude"},
	}
	help, err := axi.RenderHelp(actionsForRows(rows))
	if err != nil {
		t.Fatal(err)
	}
	want := "help[2]{cmd,why}:\n  bench link,repair stale\n  bench link,repair unwired\n"
	if help != want {
		t.Fatalf("help = %q, want %q", help, want)
	}
}

func TestCommandAppendsHonestEmptyHelpForCleanScan(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Chdir(root)
	out, code := Command(nil)
	if code != 0 || !strings.HasSuffix(out, "help[0]{cmd,why}:\n") || strings.Contains(out, "bench link") {
		t.Fatalf("Command = (exit %d, %q), want honest empty help", code, out)
	}
}

func TestCommandPreservesCheckedInCleanPrimaryResponse(t *testing.T) {
	primary, err := os.ReadFile("testdata/pre-disclosure-clean.stdout")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Chdir(root)
	out, code := Command(nil)
	if code != 0 || out != string(primary)+"help[0]{cmd,why}:\n" {
		t.Fatalf("Command = (%d, %q), want checked-in primary plus exactly one help block", code, out)
	}
}

func TestCommandPreservesCheckedInStaleAndUnwiredPrimaryResponse(t *testing.T) {
	primary, err := os.ReadFile("testdata/pre-disclosure-stale-unwired.stdout")
	if err != nil {
		t.Fatal(err)
	}
	oldEnumerate, oldInspect := enumerateGuards, inspectGuard
	t.Cleanup(func() { enumerateGuards, inspectGuard = oldEnumerate, oldInspect })
	enumerateGuards = func(context.Context, string) ([]candidate, error) {
		return []candidate{{fallback: "stale"}, {fallback: "unwired"}}, nil
	}
	inspectGuard = func(_ context.Context, _ string, c candidate) [][]string {
		if c.fallback == "stale" {
			return [][]string{{"stale", "boundary", "denies", "", "", "stale", "claude"}}
		}
		return [][]string{{"unwired", "boundary", "denies", "", "", "", "none"}}
	}
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Chdir(root)
	out, code := Command(nil)
	const help = "help[2]{cmd,why}:\n  bench link,repair stale\n  bench link,repair unwired\n"
	if code != 0 || out != string(primary)+help {
		t.Fatalf("Command = (%d, %q), want checked-in primary plus exactly one help block", code, out)
	}
}

// TestCommandRendersRealStaleManagedPrePushHookAndRepairAction drives the whole disclosure
// off a real fixture: the kit's own managed pre-push asset, installed at the hook path and
// drifted the way an install left behind by an older kit drifts, so adopt's currency check
// reads it as stale. Nothing is stubbed — the row's cells and the appended `bench link`
// repair are the production scan's own derivation, which is what the stubbed stale/unwired
// tests above cannot show.
func TestCommandRendersRealStaleManagedPrePushHookAndRepairAction(t *testing.T) {
	// The fixture is the shipped hook body itself rather than a hand-copied one, so the
	// only difference from a current install is the trailing drift line.
	canonical, err := os.ReadFile(filepath.Join("..", "adopt", "prepush.sh"))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	// origin/HEAD resolves the protected branch live, the posture a linked repo has.
	if out, err := exec.Command("git", "-C", root, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main").CombinedOutput(); err != nil {
		t.Fatalf("symbolic-ref: %v: %s", err, out)
	}
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	stale := string(canonical) + "\n# left behind by an older bench link\n"
	if err := os.WriteFile(filepath.Join(hooks, "pre-push"), []byte(stale), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	out, code := Command(nil)
	want := "guards[1]{guard,boundary,denies,branch,provenance,currency,wired}:\n" +
		"  pre-push,pre-push,direct push to the protected branch; .bench drift when pinned,main,live,stale,git\n" +
		"guard_scan[1]{status,inspected,total,omitted,reason}:\n" +
		"  complete,\"1\",\"1\",\"0\",none\n" +
		"help[1]{cmd,why}:\n  bench link,repair pre-push\n"
	if code != 0 || out != want {
		t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
	}
}

func TestCommandPreservesCheckedInIncompleteTimeoutPrimaryResponse(t *testing.T) {
	primary, err := os.ReadFile("testdata/pre-disclosure-incomplete-timeout.stdout")
	if err != nil {
		t.Fatal(err)
	}
	oldEnumerate, oldInspect, oldTimeout := enumerateGuards, inspectGuard, guardScanTimeout
	t.Cleanup(func() { enumerateGuards, inspectGuard, guardScanTimeout = oldEnumerate, oldInspect, oldTimeout })
	guardScanTimeout = 10 * time.Millisecond
	enumerateGuards = func(context.Context, string) ([]candidate, error) {
		return []candidate{{fallback: "fast"}, {fallback: "blocked"}}, nil
	}
	inspectGuard = func(ctx context.Context, _ string, c candidate) [][]string {
		if c.fallback == "fast" {
			return [][]string{{"fast", "boundary", "denies", "", "", "current", "claude"}}
		}
		<-ctx.Done()
		return nil
	}
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Chdir(root)
	out, code := Command(nil)
	if code != 0 || out != string(primary)+"help[0]{cmd,why}:\n" {
		t.Fatalf("Command = (%d, %q), want checked-in timeout/incomplete primary plus honest empty help", code, out)
	}
}

func TestCommandPreservesCheckedInEnumerationTimeoutPrimaryResponse(t *testing.T) {
	primary, err := os.ReadFile("testdata/pre-disclosure-enumeration-timeout.stdout")
	if err != nil {
		t.Fatal(err)
	}
	oldEnumerate, oldTimeout := enumerateGuards, guardScanTimeout
	t.Cleanup(func() { enumerateGuards, guardScanTimeout = oldEnumerate, oldTimeout })
	guardScanTimeout = 10 * time.Millisecond
	enumerateGuards = func(ctx context.Context, _ string) ([]candidate, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Chdir(root)
	out, code := Command(nil)
	if code != 0 || out != string(primary)+"help[0]{cmd,why}:\n" {
		t.Fatalf("Command = (%d, %q), want checked-in enumeration-timeout primary plus honest empty help", code, out)
	}
}

func TestScanTimeoutPreservesPartialRowsAndHonestCounts(t *testing.T) {
	oldEnumerate, oldInspect := enumerateGuards, inspectGuard
	t.Cleanup(func() { enumerateGuards, inspectGuard = oldEnumerate, oldInspect })
	enumerateGuards = func(context.Context, string) ([]candidate, error) {
		return []candidate{{fallback: "fast"}, {fallback: "blocked"}}, nil
	}
	blocked := make(chan struct{})
	inspectGuard = func(ctx context.Context, _ string, c candidate) [][]string {
		if c.fallback == "fast" {
			return [][]string{{"fast", "boundary", "denies", "none"}}
		}
		close(blocked)
		<-ctx.Done()
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan ScanResult, 1)
	go func() { result <- Scan(ctx, "/repo") }()
	<-blocked
	cancel()
	got := <-result
	if got.Status != "incomplete" || got.Reason != "timeout" || got.Inspected != "1" || got.Total != "2" || got.Omitted != "1" || len(got.Rows) != 1 {
		t.Fatalf("Scan = %#v", got)
	}
}

func TestScanEnumerationTimeoutUsesUnknownCounts(t *testing.T) {
	old := enumerateGuards
	t.Cleanup(func() { enumerateGuards = old })
	enumerateGuards = func(ctx context.Context, _ string) ([]candidate, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	got := Scan(ctx, "/repo")
	if got.Status != "incomplete" || got.Total != "unknown" || got.Omitted != "unknown" || got.Reason != "timeout" {
		t.Fatalf("Scan = %#v", got)
	}
}

func TestScanWaitsForCancelledWorkerCleanup(t *testing.T) {
	oldEnumerate, oldInspect := enumerateGuards, inspectGuard
	t.Cleanup(func() { enumerateGuards, inspectGuard = oldEnumerate, oldInspect })
	enumerateGuards = func(context.Context, string) ([]candidate, error) {
		return []candidate{{fallback: "one"}}, nil
	}
	started := make(chan struct{})
	cancelled := make(chan struct{})
	release := make(chan struct{})
	cleaned := make(chan struct{})
	inspectGuard = func(ctx context.Context, _ string, _ candidate) [][]string {
		close(started)
		<-ctx.Done()
		close(cancelled)
		<-release
		close(cleaned)
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan ScanResult, 1)
	go func() { result <- Scan(ctx, "/repo") }()
	<-started
	cancel()
	<-cancelled
	select {
	case got := <-result:
		t.Fatalf("Scan returned before worker cleanup: %#v", got)
	default:
	}
	close(release)
	var got ScanResult
	select {
	case got = <-result:
	case <-time.After(time.Second):
		t.Fatal("Scan did not return after worker cleanup")
	}
	<-cleaned
	if got.Status != "incomplete" || got.Inspected != "0" || got.Omitted != "1" || got.Reason != "timeout" {
		t.Fatalf("Scan = %#v", got)
	}
}

func TestManifestFieldReadsLeadingCommentBlock(t *testing.T) {
	header := "#!/usr/bin/env bash\n  \t\n# threat-model prose\n# denies: destructive git operations   \n# name: block-dangerous-git\n# why: protects repository history\n# boundary: PreToolUse Bash\nset -uo pipefail\n# name: too-late\n"
	if got := manifestField(header, "name"); got != "block-dangerous-git" {
		t.Errorf("name = %q", got)
	}
	if got := manifestField(header, "denies"); got != "destructive git operations" {
		t.Errorf("denies = %q", got)
	}
	if got := manifestField(header, "absent"); got != "" {
		t.Errorf("absent = %q, want empty", got)
	}
	if got := manifestField("# name: \n# name: later\n", "name"); got != "" {
		t.Errorf("first empty name lost first-occurrence precedence: %q", got)
	}
}

func TestManifestMissingRequiredUsesCanonicalOrder(t *testing.T) {
	manifest := parseHeader("# why: present\n# name: \n# name: ignored\n")
	if got, want := manifest.MissingRequired(), []string{"name", "boundary", "denies"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("MissingRequired = %v, want %v", got, want)
	}
}

func TestGuardRowReadsStaticHeader(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
		return p
	}
	full := write("full.sh", "#!/usr/bin/env bash\n# why: w\n# name: g\n# denies: d\n# boundary: b\nexit 99\n")
	if r, emit := guardRow(full, "full"); !emit || !reflect.DeepEqual(r, []string{"g", "b", "d"}) {
		t.Errorf("full manifest row = %v emit=%v", r, emit)
	}
	link := filepath.Join(dir, "full-link.sh")
	if err := os.Symlink(full, link); err == nil {
		if r, emit := guardRow(link, "full-link"); !emit || !reflect.DeepEqual(r, []string{"g", "b", "d"}) {
			t.Errorf("regular symlink manifest row = %v emit=%v", r, emit)
		}
	}

	info := write("info.sh", "#!/usr/bin/env bash\n# name: s\n# boundary: SessionStart\n# denies: nothing (informational)\n# why: w\nexit 99\n")
	if _, emit := guardRow(info, "info"); emit {
		t.Errorf("informational guard should be excluded")
	}

	for _, tc := range []struct {
		name string
		body string
	}{
		{"absent", "#!/usr/bin/env bash\nexit 99\n"},
		{"incomplete", "#!/usr/bin/env bash\n# name: partial\n# boundary: b\n# denies: d\nexit 99\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := write(tc.name+".sh", tc.body)
			if r, emit := guardRow(path, tc.name); !emit || !reflect.DeepEqual(r, []string{tc.name, "", "no manifest"}) {
				t.Errorf("row = %v emit=%v, want no manifest", r, emit)
			}
		})
	}
}

func TestGuardRowRejectsFIFOWithoutOpening(t *testing.T) {
	if fifo := os.Getenv("BENCH_TEST_GUARD_FIFO"); fifo != "" {
		row, emit := guardRow(fifo, "fifo")
		if !emit || !reflect.DeepEqual(row, []string{"fifo", "", "no manifest"}) {
			t.Fatalf("FIFO row = %v emit=%v", row, emit)
		}
		return
	}
	if runtime.GOOS == "windows" {
		capability.Capability(t, capability.Fifo, "named pipes use the Unix mkfifo fixture")
	}
	fifo := filepath.Join(t.TempDir(), "special.sh")
	if out, err := exec.Command("mkfifo", fifo).CombinedOutput(); err != nil {
		t.Fatalf("mkfifo: %v: %s", err, out)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGuardRowRejectsFIFOWithoutOpening$")
	cmd.Env = append(os.Environ(), "BENCH_TEST_GUARD_FIFO="+fifo)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatal("guardRow opened a FIFO and blocked")
	}
	if err != nil {
		t.Fatalf("FIFO helper failed: %v\n%s", err, out)
	}
}
