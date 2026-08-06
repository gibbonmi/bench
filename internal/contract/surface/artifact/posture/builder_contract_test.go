package posture

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestArtifactModeCommandTrace(t *testing.T) {
	root := contract.SubjectRoot(t)
	dir := t.TempDir()
	out := filepath.Join(dir, "artifact")
	goLog := filepath.Join(dir, "go.log")
	execLog := filepath.Join(dir, "exec.log")
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, out)
	cmd.Env = fakeBuilderEnv(t, "complete", "GOOS=plan9", "GOARCH=arm64", "BENCH_TEST_GO_LOG="+goLog, "BENCH_TEST_EXEC_LOG="+execLog)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("artifact command trace failed: %v\n%s", err, output)
	}
	trace, err := os.ReadFile(goLog)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(trace)), "\n")
	if len(lines) != 1 || !strings.HasPrefix(lines[0], "target=plan9/arm64 go <build>") {
		t.Fatalf("artifact Go trace = %q, want one target-preserving go build", trace)
	}
	if strings.Contains(string(trace), "<run>") || strings.Contains(string(trace), "freshness-publish") {
		t.Fatalf("artifact trace reached a publisher/helper: %q", trace)
	}
	if _, err := os.Stat(execLog); !os.IsNotExist(err) {
		t.Fatalf("artifact build executed its staged output: %v", err)
	}
	if info, err := os.Stat(out); err != nil || !info.Mode().IsRegular() || info.Mode()&0o111 == 0 {
		t.Fatalf("artifact output is not executable: %v, %v", info, err)
	}
	if _, err := os.Stat(out + ".seal"); !os.IsNotExist(err) {
		t.Fatalf("artifact command trace left a seal: %v", err)
	}
}

func TestArtifactModeFailureTablePreservesPriorPair(t *testing.T) {
	for _, test := range []struct {
		name     string
		behavior string
		extra    []string
	}{
		{name: "compile", behavior: "fail"},
		{name: "validation", behavior: "invalid"},
		{name: "promotion", behavior: "complete", extra: []string{"BENCH_TEST_FAIL_PROMOTION=1"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := contract.SubjectRoot(t)
			dir := t.TempDir()
			out := filepath.Join(dir, "prior output")
			writePriorPair(t, out)
			cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, out)
			cmd.Env = fakeBuilderEnv(t, test.behavior, test.extra...)
			if output, err := cmd.CombinedOutput(); err == nil {
				t.Fatalf("%s failure seam succeeded: %q", test.name, output)
			}
			assertPriorPair(t, out)
			if residue, err := filepath.Glob(filepath.Join(dir, ".bench.*")); err != nil || len(residue) != 0 {
				t.Fatalf("%s failure left staging residue: %v, %v", test.name, residue, err)
			}
		})
	}
}

func TestMalformedBuilderGrammarInvokesNoGoAndPreservesPriorPair(t *testing.T) {
	root := contract.SubjectRoot(t)
	script := filepath.Join(root, "scripts", "go-build.sh")
	for _, test := range []struct {
		name string
		args func(string) []string
	}{
		{name: "missing-all-operands", args: func(string) []string { return nil }},
		{name: "missing-default-output", args: func(string) []string { return []string{root} }},
		{name: "missing-mode-value", args: func(out string) []string { return []string{"--mode", root, out} }},
		{name: "missing-artifact-output", args: func(string) []string { return []string{"--mode", "artifact", root} }},
		{name: "unknown-mode", args: func(out string) []string { return []string{"--mode", "unknown", root, out} }},
		{name: "duplicate-mode", args: func(out string) []string { return []string{"--mode", "artifact", "--mode", "artifact", root, out} }},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			out := filepath.Join(dir, "prior output")
			writePriorPair(t, out)
			goLog := filepath.Join(dir, "go.log")
			cmd := exec.Command("bash", append([]string{script}, test.args(out)...)...)
			cmd.Env = fakeBuilderEnv(t, "complete", "BENCH_TEST_GO_LOG="+goLog)
			if output, err := cmd.CombinedOutput(); err == nil {
				t.Fatalf("malformed grammar succeeded: %q", output)
			}
			if data, err := os.ReadFile(goLog); err == nil && len(data) != 0 {
				t.Fatalf("malformed grammar invoked Go: %q", data)
			} else if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			assertPriorPair(t, out)
		})
	}
}

func TestArtifactModeRerunsConvergeUnsealed(t *testing.T) {
	root := contract.SubjectRoot(t)
	out := filepath.Join(t.TempDir(), "rerun output")
	run := func(mode []string, body string) []byte {
		t.Helper()
		args := append([]string{filepath.Join(root, "scripts", "go-build.sh")}, mode...)
		args = append(args, root, out)
		cmd := exec.Command("bash", args...)
		cmd.Env = fakeBuilderEnv(t, "complete", "BENCH_TEST_BUILD_BODY="+body)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s build failed: %v\n%s", body, err, output)
		}
		data, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	subject := run(nil, "subject-v1")
	if _, err := os.Stat(out + ".seal"); err != nil {
		t.Fatalf("subject build did not establish prior seal: %v", err)
	}
	artifact := run([]string{"--mode", "artifact"}, "artifact-v2")
	if bytes.Equal(subject, artifact) || !bytes.Contains(artifact, []byte("artifact-v2")) {
		t.Fatal("subject-to-artifact rerun did not replace output bytes")
	}
	if _, err := os.Stat(out + ".seal"); !os.IsNotExist(err) {
		t.Fatalf("subject-to-artifact rerun left stale seal: %v", err)
	}
	second := run([]string{"--mode", "artifact"}, "artifact-v3")
	if bytes.Equal(artifact, second) || !bytes.Contains(second, []byte("artifact-v3")) {
		t.Fatal("artifact-to-artifact rerun did not replace output bytes")
	}
	if _, err := os.Stat(out + ".seal"); !os.IsNotExist(err) {
		t.Fatalf("artifact-to-artifact rerun left a seal: %v", err)
	}
}

func writePriorPair(t *testing.T, out string) {
	t.Helper()
	if err := os.WriteFile(out, []byte("prior executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(out+".seal", []byte("prior seal"), 0o644); err != nil {
		t.Fatal(err)
	}
}
