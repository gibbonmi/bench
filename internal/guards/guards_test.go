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
