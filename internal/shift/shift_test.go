package shift

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent"
)

// rec builds one NUL-terminated porcelain -z record: two status chars, a space, the
// path, then the NUL delimiter git writes.
func rec(status, path string) []byte {
	return append([]byte(status+" "+path), 0)
}

// TestParseDirtyPaths pins the native porcelain -z parse. NUL-delimited paths keep
// spaces, glob characters, and a literal newline, the case the shell's `printf | sort`
// pipeline misread by splitting it across two lines. The scratch files are excluded and
// the result is sorted.
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

// shiftCollisionFixture builds a bare repo plus a passing gate and agent, and points
// timeNow at a fixed instant so the derived branch name is deterministic. preExisting
// names additional branches, relative to the base bench/shift-<ts> name (e.g. "-2"), for
// the fixture to pre-create. This exercises the loop's collision retry.
func shiftCollisionFixture(t *testing.T, preExisting ...string) (tmp, baseBranch string) {
	t.Helper()
	tmp = t.TempDir()
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
	baseBranch = "bench/shift-" + fixed.Format("20060102-150405")
	runGit("branch", baseBranch)
	for _, suffix := range preExisting {
		runGit("branch", baseBranch+suffix)
	}
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
	return tmp, baseBranch
}

// TestLoopRetriesBranchCreationOnCollision covers spec row 18. When the derived
// bench/shift-<ts> branch already exists, because two shifts land in the same second,
// the loop retries with a disambiguating "-2" suffix. It does not fail the shift. The
// recovery ref path, built from s.branch, follows the resolved, suffixed name.
func TestLoopRetriesBranchCreationOnCollision(t *testing.T) {
	_, baseBranch := shiftCollisionFixture(t)

	var stdout, stderr bytes.Buffer
	code := Loop("branch collision", &stdout, &stderr)
	if code != 4 { // no-op adapter (exit 0, no commit) reads as no-op/4
		t.Fatalf("Loop = %d, want 4 (no-op); stdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	wantBranch := baseBranch + "-2"
	if !contains(stdout.String(), wantBranch) {
		t.Fatalf("stdout did not name the suffixed branch %s:\n%s", wantBranch, stdout.String())
	}
	if contains(stderr.String(), "could not create shift branch") {
		t.Fatalf("stderr reported a branch creation failure despite the retry:\n%s", stderr.String())
	}
}

// TestLoopReportsBranchCreationFailureAfterExhaustingRetries covers the bound on row
// 18's retry. Once every suffix through -10 is already taken, the loop gives up and
// reports the failure exactly as it always has for an unresolvable collision.
func TestLoopReportsBranchCreationFailureAfterExhaustingRetries(t *testing.T) {
	var taken []string
	for i := 2; i <= 10; i++ {
		taken = append(taken, fmt.Sprintf("-%d", i))
	}
	_, baseBranch := shiftCollisionFixture(t, taken...)

	var stdout, stderr bytes.Buffer
	if code := Loop("branch collision", &stdout, &stderr); code == 0 {
		t.Fatalf("Loop returned success despite exhausted collision retries; stdout:\n%s\nstderr:\n%s", stdout.String(), stderr.String())
	}
	if !contains(stderr.String(), "could not create shift branch") {
		t.Fatalf("stderr did not report branch creation failure:\n%s", stderr.String())
	}
	if contains(stdout.String(), "shift done") {
		t.Fatalf("shift reported completion after exhausting branch creation retries:\n%s", stdout.String())
	}
	_ = baseBranch
}

func TestLoopPersistsIntentBeforeAcquireFailure(t *testing.T) {
	root := t.TempDir()
	gitCmd := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	gitCmd("init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(root, "tracked"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitCmd("add", ".")
	gitCmd("-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "init")
	old, _ := os.Getwd()
	_ = os.Chdir(root)
	t.Cleanup(func() { _ = os.Chdir(old) })
	t.Setenv("BENCH_AGENT", "/bin/true")
	blockedHome := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(blockedHome, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_HOME", blockedHome)
	var stdout, stderr bytes.Buffer
	if code := Loop("multi word objective", &stdout, &stderr); code == 0 {
		t.Fatal("Loop unexpectedly acquired worktree")
	}
	ledger, err := intent.Read(root)
	if err != nil || len(ledger.Entries) != 1 || ledger.Entries[0].Kind != intent.KindShift || ledger.Entries[0].Worktree != "" {
		t.Fatalf("pre-acquire intent = %#v, %v", ledger.Entries, err)
	}
}

// TestRunGateReportsRed pins runGate's contract: it reports a red gate's exit code
// straight through and does nothing to the session itself. Preservation, whether
// snapshot-and-release or retain-and-lock on a snapshot failure, is the caller's job.
// It happens once, explicitly, via preserveAndRecover at the return site, never
// implied by a flag runGate sets on its way out.
func TestRunGateReportsRed(t *testing.T) {
	root := t.TempDir()
	gitCmd := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	gitCmd("init", "-q")
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte("#!/usr/bin/env bash\nexit 23\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(`{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	gitCmd("add", "-A")
	gitCmd("-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "init")
	s := &session{root: root, stdout: io.Discard, stderr: io.Discard}
	if rc := s.runGate(); rc == 0 {
		t.Fatal("red gate returned zero")
	}
	if s.preserve.Load() {
		t.Fatal("runGate must not itself mark the session preserved — that is preserveAndRecover's job")
	}
}

// TestCheckpointOutcomeDeadlineWinsOverInterrupt pins the decided precedence when a
// wall deadline and a pulled line race and both flags land before the next checkpoint.
// The deadline outcome, incomplete/3, wins over interrupted/130. checkpoint's decision
// is exercised through checkpointOutcome directly, since checkpoint itself exits the
// process on a hit.
func TestCheckpointOutcomeDeadlineWinsOverInterrupt(t *testing.T) {
	s := &session{}
	s.deadline.Store(true)
	s.interrupted.Store(true)

	outcome, detail, ok := s.checkpointOutcome()

	if !ok {
		t.Fatal("checkpointOutcome reported no hit with both flags set")
	}
	if outcome != OutcomeIncomplete {
		t.Fatalf("outcome = %q, want %q (deadline must win over interrupted)", outcome, OutcomeIncomplete)
	}
	if detail != "wall deadline exceeded" {
		t.Fatalf("detail = %q, want the deadline detail", detail)
	}
}

// TestObjectiveBannerEscapesControls covers story 14's banner row: a control sequence
// in the objective renders escaped in the shift-start banner rather than raw. This
// happens because the banner routes through the shared sanitizer. The banner never
// interpolates the raw objective. The " — objective:" delimiter the branch parser
// depends on must survive.
func TestObjectiveBannerEscapesControls(t *testing.T) {
	esc := string(rune(0x1b))
	got := objectiveBanner("bench/shift-x", "paint it "+esc+"[31mred")
	if strings.ContainsRune(got, 0x1b) {
		t.Fatalf("banner leaked a raw ESC byte: %q", got)
	}
	if !strings.Contains(got, `[31mred`) {
		t.Fatalf("banner did not escape the control sequence: %q", got)
	}
	if !strings.Contains(got, "bench/shift-x — objective: ") {
		t.Fatalf("banner dropped the branch delimiter the parser needs: %q", got)
	}
}

func TestValidateObjective(t *testing.T) {
	cases := []struct {
		name      string
		objective string
		wantErr   bool
	}{
		{"empty", "", true},
		{"whitespace only spaces", "   ", true},
		{"whitespace only tab", "\t\t", true},
		{"plain", "improve the parser", false},
		{"esc byte", "bad\x1bobjective", true},
		{"tab byte", "bad\tobjective", true},
		{"del byte", "bad\x7fobjective", true},
		{"newline byte", "bad\nobjective", true},
		{"at cap", strings.Repeat("a", objectiveMaxRunes), false},
		{"over cap", strings.Repeat("a", objectiveMaxRunes+1), true},
		{"multibyte at cap", strings.Repeat("é", objectiveMaxRunes), false},
		{"multibyte over cap", strings.Repeat("é", objectiveMaxRunes+1), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateObjective(tc.objective)
			if tc.wantErr && err == nil {
				t.Errorf("validateObjective(%q) = nil, want error", tc.objective)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("validateObjective(%q) = %v, want nil", tc.objective, err)
			}
		})
	}
}

func TestParseBoundedInt(t *testing.T) {
	const name = "BENCH_TEST_ITERS"
	cases := []struct {
		name    string
		env     string // "" means unset
		unset   bool
		want    int
		wantErr bool
	}{
		{"unset", "", true, 12, false},
		{"empty string", "", false, 12, false},
		{"valid", "5", false, 5, false},
		{"min boundary", "1", false, 1, false},
		{"max boundary", "100", false, 100, false},
		{"zero", "0", false, 0, true},
		{"negative", "-1", false, 0, true},
		{"over max", "101", false, 0, true},
		{"non-integer", "abc", false, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.unset {
				os.Unsetenv(name)
			} else {
				t.Setenv(name, tc.env)
			}
			got, err := parseBoundedInt(name, 12, 1, 100)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseBoundedInt(%q) = %d, nil, want error", tc.env, got)
				}
				if !contains(err.Error(), name) || !contains(err.Error(), "[1,100]") {
					t.Errorf("error %q does not name the variable and range", err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("parseBoundedInt(%q) = %v, want nil", tc.env, err)
			}
			if got != tc.want {
				t.Errorf("parseBoundedInt(%q) = %d, want %d", tc.env, got, tc.want)
			}
		})
	}
}

func TestParseWallDuration(t *testing.T) {
	const name = "BENCH_TEST_WALL"
	cases := []struct {
		name    string
		env     string
		unset   bool
		want    time.Duration
		wantErr bool
	}{
		{"unset", "", true, 0, false},
		{"empty string", "", false, 0, false},
		{"valid", "30m", false, 30 * time.Minute, false},
		{"max boundary", "24h", false, 24 * time.Hour, false},
		{"zero", "0s", false, 0, true},
		{"negative", "-1h", false, 0, true},
		{"over max", "48h", false, 0, true},
		{"unparseable", "soon", false, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.unset {
				os.Unsetenv(name)
			} else {
				t.Setenv(name, tc.env)
			}
			got, err := parseWallDuration(name, 0, 24*time.Hour)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseWallDuration(%q) = %v, nil, want error", tc.env, got)
				}
				if !contains(err.Error(), name) {
					t.Errorf("error %q does not name the variable", err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("parseWallDuration(%q) = %v, want nil", tc.env, err)
			}
			if got != tc.want {
				t.Errorf("parseWallDuration(%q) = %v, want %v", tc.env, got, tc.want)
			}
		})
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
