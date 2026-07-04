package shift

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// rec builds one NUL-terminated porcelain -z record: two status chars, a space, the
// path, then the NUL delimiter git writes.
func rec(status, path string) []byte {
	return append([]byte(status+" "+path), 0)
}

// TestParseDirtyPaths pins the native porcelain -z parse: NUL-delimited paths keep
// spaces, glob characters, and — the case the shell's `printf | sort` pipeline misread
// by splitting it across two lines — a literal newline, while the scratch files are
// excluded and the result is sorted.
func TestParseDirtyPaths(t *testing.T) {
	var raw []byte
	raw = append(raw, rec(" M", "step 2 [a].txt")...)
	raw = append(raw, rec("??", "a\nb.txt")...) // literal newline inside one path
	raw = append(raw, rec(" M", ".bench-notes.md")...)
	raw = append(raw, rec("??", ".bench-objective")...)
	raw = append(raw, rec(" M", "step 1 [a].txt")...)

	got := parseDirtyPaths(raw)
	want := []string{"a\nb.txt", "step 1 [a].txt", "step 2 [a].txt"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseDirtyPaths = %q, want %q", got, want)
	}
}

// TestParseDirtyPathsEmpty covers a clean tree (no records, only the possible trailing
// NUL) and an all-scratch tree — both yield no touched paths.
func TestParseDirtyPathsEmpty(t *testing.T) {
	if got := parseDirtyPaths(nil); len(got) != 0 {
		t.Errorf("empty input = %q, want none", got)
	}
	scratch := append(rec(" M", ".bench-objective"), rec(" M", ".bench-notes.md")...)
	if got := parseDirtyPaths(scratch); len(got) != 0 {
		t.Errorf("all-scratch input = %q, want none", got)
	}
}

func TestRequireAdapter(t *testing.T) {
	cases := []struct {
		name    string
		agent   string
		wantErr string // substring the error must contain, "" = no error
	}{
		{"unset", "", "BENCH_AGENT"},
		{"missing path", "/no/such/adapter", "not executable"},
		{"shell keyword", "if", "not executable"},
		{"real executable", "true", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := requireAdapter(tc.agent)
			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("requireAdapter(%q) = %v, want nil", tc.agent, err)
				}
				return
			}
			if err == nil || !contains(err.Error(), tc.wantErr) {
				t.Errorf("requireAdapter(%q) = %v, want error containing %q", tc.agent, err, tc.wantErr)
			}
		})
	}
}

func TestLoopReportsBranchCreationFailure(t *testing.T) {
	tmp := t.TempDir()
	runGit := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = tmp
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}
	runGit("init", "-q")
	if err := os.Mkdir(filepath.Join(tmp, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	gatePath := filepath.Join(tmp, ".bench", "gate.sh")
	if err := os.WriteFile(gatePath, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	agentPath := filepath.Join(tmp, "agent")
	if err := os.WriteFile(agentPath, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	runGit("-c", "user.email=bench@local", "-c", "user.name=bench", "add", "-A")
	runGit("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "init")

	fixed := time.Date(2026, 7, 4, 9, 30, 0, 0, time.UTC)
	branch := "bench/shift-" + fixed.Format("20060102-150405")
	runGit("branch", branch)
	oldNow := timeNow
	timeNow = func() time.Time { return fixed }
	t.Cleanup(func() { timeNow = oldNow })

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	t.Setenv("BENCH_AGENT", agentPath)
	t.Setenv("BENCH_HOME", filepath.Join(tmp, "bench-home"))
	t.Setenv("BENCH_MAX_ITERS", "1")

	var stdout, stderr bytes.Buffer
	if code := Loop("branch collision", bytes.NewReader(nil), &stdout, &stderr); code == 0 {
		t.Fatalf("Loop returned success on branch collision; stdout:\n%s\nstderr:\n%s", stdout.String(), stderr.String())
	}
	if !contains(stderr.String(), "could not create shift branch "+branch) {
		t.Fatalf("stderr did not report branch creation failure:\n%s", stderr.String())
	}
	if contains(stdout.String(), "shift done") {
		t.Fatalf("shift reported completion after branch creation failure:\n%s", stdout.String())
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
