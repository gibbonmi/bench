package gate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// TestResolvePrecedence pins the ordered chain as a pure function: `.bench/gate.sh`
// beats `$BENCH_GATE` beats the auto-detect lockfiles in their fixed order, and an
// absent gate resolves to None. A reordered chain would silently run the wrong oracle;
// no black-box assertion pins this cheaply, so the table is the genuinely-red signal.
func TestResolvePrecedence(t *testing.T) {
	// present names the set of paths the injected probes report as existing/executable.
	fs := func(present ...string) FS {
		set := map[string]bool{}
		for _, p := range present {
			set[p] = true
		}
		has := func(p string) bool { return set[p] }
		return FS{Executable: has, Exists: has}
	}
	const gateSh = "/r/.bench/gate.sh"

	cases := []struct {
		name      string
		benchGate string
		fs        FS
		want      Kind
	}{
		{"gate.sh beats BENCH_GATE and lockfiles", "echo hi", fs(gateSh, "/r/package.json"), GateSh},
		{"BENCH_GATE beats lockfiles", "echo hi", fs("/r/package.json"), BenchGate},
		{"pnpm beats npm", "", fs("/r/pnpm-lock.yaml", "/r/package.json"), Pnpm},
		{"package.json picks npm", "", fs("/r/package.json"), Npm},
		{"pyproject after npm", "", fs("/r/pyproject.toml"), Pyproject},
		{"cargo last", "", fs("/r/Cargo.toml"), Cargo},
		{"nothing resolves to None", "", fs(), None},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Resolve("/r", tc.benchGate, tc.fs)
			if got.Kind != tc.want {
				t.Errorf("Resolve = %v, want %v", got.Kind, tc.want)
			}
			if tc.want == BenchGate && got.Command != tc.benchGate {
				t.Errorf("BenchGate command = %q, want %q", got.Command, tc.benchGate)
			}
		})
	}
}

func TestGateEnvStripsWrapperRoutingInternals(t *testing.T) {
	t.Setenv("BENCH_KIT", "/wrong/kit")
	t.Setenv("BENCH_WRAPPER", "/wrong/wrapper")
	t.Setenv("BENCH_GATE", "echo ok")

	env := gateEnv()
	sawBenchGate := false
	for _, kv := range env {
		if strings.HasPrefix(kv, "BENCH_KIT=") || strings.HasPrefix(kv, "BENCH_WRAPPER=") {
			t.Fatalf("gateEnv leaked wrapper-routing internal %q", kv)
		}
		if kv == "BENCH_GATE=echo ok" {
			sawBenchGate = true
		}
	}
	if !sawBenchGate {
		t.Fatal("gateEnv stripped BENCH_GATE, which is part of the project gate contract")
	}
}

func writeGateTestFile(t *testing.T, root, name, body string, mode os.FileMode) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), mode); err != nil {
		t.Fatal(err)
	}
}

func gateTestRepo(t *testing.T, script, manifest string) string {
	t.Helper()
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, output)
	}
	writeGateTestFile(t, root, ".bench/gate.sh", script, 0o755)
	if manifest != "" {
		writeGateTestFile(t, root, ".bench/gate-inputs.json", manifest, 0o644)
	}
	return root
}

func TestClosedSubjectReuseAndMutations(t *testing.T) {
	const manifest = `{"schema":1,"closure":"local","environment":["ORACLE_INPUT"],"paths":["ignored-input"],"tools":[]}`
	root := gateTestRepo(t, "#!/usr/bin/env bash\necho run >> .git/runs\n[ \"$ORACLE_INPUT\" = green ] && [ \"$(cat ignored-input)\" = green ]\n", manifest)
	writeGateTestFile(t, root, ".gitignore", "ignored-input\n", 0o644)
	writeGateTestFile(t, root, "ignored-input", "green\n", 0o644)
	t.Setenv("ORACLE_INPUT", "green")
	var stderr bytes.Buffer
	first := Execute(context.Background(), root, io.Discard, &stderr)
	if first.ActionExit != 0 || !Inspect(root).ReusableGreen {
		t.Fatalf("closed green did not become reusable: result=%+v inspect=%+v stderr=%s", first, Inspect(root), stderr.String())
	}

	t.Setenv("ORACLE_INPUT", "red")
	if got := Inspect(root); got.ReusableGreen || got.Reason != "oracle changed" {
		t.Fatalf("environment mutation reused verdict: %+v", got)
	}
	second := Execute(context.Background(), root, io.Discard, io.Discard)
	if second.GateExit == 0 || second.ActionExit == 0 {
		t.Fatalf("red mutated oracle normalized to success: %+v", second)
	}

	t.Setenv("ORACLE_INPUT", "green")
	writeGateTestFile(t, root, "ignored-input", "red\n", 0o644)
	if got := Inspect(root); got.ReusableGreen {
		t.Fatalf("declared ignored-input mutation reused verdict: %+v", got)
	}
}

func TestGateEnvironmentIsPasslisted(t *testing.T) {
	const manifest = `{"schema":1,"closure":"local","environment":["DECLARED"],"paths":[],"tools":[]}`
	root := gateTestRepo(t, "#!/usr/bin/env bash\nprintf '%s:%s:%s:%s\\n' \"${DECLARED-unset}\" \"${AMBIENT-unset}\" \"${HOME-unset}\" \"${BENCH_GATE-unset}\" > .git/environment\n", manifest)
	t.Setenv("DECLARED", "kept")
	t.Setenv("AMBIENT", "secret")
	t.Setenv("HOME", "/secret/home")
	t.Setenv("BENCH_GATE", "should-not-win-because-gate-script-does")
	result := Execute(context.Background(), root, io.Discard, io.Discard)
	if result.ActionExit != 0 {
		t.Fatalf("passlist gate failed: %+v", result)
	}
	data, err := os.ReadFile(filepath.Join(root, ".git", "environment"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "kept:unset:unset:unset\n" {
		t.Fatalf("gate environment leaked ambient values: %q", data)
	}
}

func TestStrictVerdictInspection(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	s, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	gitdir, err := benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		t.Fatal(err)
	}
	cache := filepath.Join(gitdir, benchgit.GateCacheFile)
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	ready := fmt.Sprintf(`{"schema":1,"state":"ready","status":"green","tree":%q,"oracle":%q,"recorded_at":"2026-07-13T11:50:01Z"}`, s.Tree, s.Oracle)
	write := func(body string, mode os.FileMode) {
		t.Helper()
		if err := os.RemoveAll(cache); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(cache, []byte(body), mode); err != nil {
			t.Fatal(err)
		}
	}
	write(ready, 0o600)
	if got := inspectAt(root, now); !got.ReusableGreen {
		t.Fatalf("fresh exact ready not reusable: %+v", got)
	}
	if got := inspectAt(root, now.Add(time.Second)); got.ReusableGreen || got.Reason != "verdict expired" {
		t.Fatalf("exact ten-minute boundary reused: %+v", got)
	}
	for _, invalid := range []string{
		"green " + s.Tree + " 2026-07-13T11:50:01Z\n",
		strings.Replace(ready, `"schema":1`, `"schema":1,"schema":1`, 1),
		strings.Replace(ready, `"status":"green",`, "", 1),
		strings.TrimSuffix(ready, "}") + `,"unknown":1}`,
		ready + `{}`,
	} {
		write(invalid, 0o600)
		if got := inspectAt(root, now); got.State != Invalid {
			t.Fatalf("invalid record classified %q: %+v", invalid, got)
		}
	}
	write(ready, 0o644)
	if got := inspectAt(root, now); got.State != Invalid {
		t.Fatalf("broad cache mode classified %+v", got)
	}
	pending := fmt.Sprintf(`{"schema":1,"state":"pending","tree":%q,"oracle":%q,"started_at":"2026-07-13T11:59:00Z","owner_pid":123}`, s.Tree, s.Oracle)
	write(pending, 0o600)
	if got := inspectAt(root, now); got.State != Pending || got.PendingStatus != "interrupted-pending" {
		t.Fatalf("released pending classified %+v", got)
	}
}

func TestExecutionLockAndDriftFailClosed(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\ntouch .git/started\nwhile [ ! -f .git/release ]; do sleep .01; done\nprintf drift > drift.txt\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	done := make(chan Result, 1)
	go func() { done <- Execute(context.Background(), root, io.Discard, io.Discard) }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(root, ".git", "started")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("gate did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if got := Inspect(root); got.State != Pending || got.PendingStatus != "locked-pending" {
		t.Fatalf("live owner not projected: %+v", got)
	}
	if second := Execute(context.Background(), root, io.Discard, io.Discard); second.ActionExit != 1 {
		t.Fatalf("concurrent execution did not fail immediately: %+v", second)
	}
	writeGateTestFile(t, root, ".git/release", "\n", 0o600)
	first := <-done
	if first.ActionExit != 1 || first.GateExit != 0 || first.Inspection.State != Pending {
		t.Fatalf("mid-run drift did not leave pending operational failure: %+v", first)
	}
}

func TestNoGateReturnsThreeWithoutRecord(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", "")
	if err := os.Remove(filepath.Join(root, ".bench", "gate.sh")); err != nil {
		t.Fatal(err)
	}
	result := Execute(context.Background(), root, io.Discard, io.Discard)
	if result.GateExit != 3 || result.ActionExit != 3 {
		t.Fatalf("no gate result = %+v", result)
	}
	if _, err := os.Stat(filepath.Join(root, ".git", benchgit.GateCacheFile)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("no gate wrote verdict: %v", err)
	}
}
