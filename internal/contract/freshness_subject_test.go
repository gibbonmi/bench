package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/contract/internal/freshnessfixture"
	"github.com/gibbonmi/bench/internal/freshness"
)

const freshnessSubjectModeEnv = "BENCH_FT131_FRESHNESS_SUBJECT_MODE"
const freshnessSubjectCaseEnv = "BENCH_FT131_FRESHNESS_SUBJECT_CASE"

func TestRequireFreshBenchRefusesStaleSubject(t *testing.T) {
	if os.Getenv(freshnessSubjectModeEnv) == "stale" {
		RequireFreshBench(t)
		t.Fatal("stale subject reached the contract assertion")
	}

	root := freshnessfixture.StaleSubject(t, "old subject output")
	command := exec.Command(os.Args[0], "-test.run=^TestRequireFreshBenchRefusesStaleSubject$")
	command.Env = append(os.Environ(), freshnessSubjectModeEnv+"=stale", canary.SubjectRootEnv+"="+root)
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("stale subject reached the contract assertion:\n%s", output)
	}
	for _, want := range []string{"bench binary", "rebuild with " + freshness.RebuildAction(root)} {
		if !strings.Contains(string(output), want) {
			t.Fatalf("stale subject diagnostic missing %q:\n%s", want, output)
		}
	}
	if strings.Contains(string(output), "stale subject reached the contract assertion") {
		t.Fatalf("stale subject reached the assertion instead of freshness:\n%s", output)
	}
}

func TestRequireFreshBenchIsRepeatableAcrossCWD(t *testing.T) {
	if os.Getenv(freshnessSubjectModeEnv) == "repeat" {
		RequireFreshBench(t)
		RequireFreshBench(t)
		return
	}

	root := freshnessfixture.PublishedSubject(t, "fresh subject output")
	nested := filepath.Join(root, "nested cwd")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "-A"}, {"-c", "user.email=contract@example.invalid", "-c", "user.name=contract", "commit", "-qm", "fixture"}} {
		command := exec.Command("git", append([]string{"-C", root}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("initialize tracked subject: %v\n%s", err, output)
		}
	}
	for _, dir := range []string{root, nested} {
		command := exec.Command(os.Args[0], "-test.run=^TestRequireFreshBenchIsRepeatableAcrossCWD$")
		command.Dir = dir
		command.Env = append(os.Environ(), freshnessSubjectModeEnv+"=repeat", canary.SubjectRootEnv+"="+root)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("freshness check from %q: %v\n%s", dir, err, output)
		}
	}
	status := exec.Command("git", "-C", root, "status", "--porcelain")
	output, err := status.Output()
	if err != nil {
		t.Fatalf("read tracked subject status: %v", err)
	}
	if strings.TrimSpace(string(output)) != "" {
		t.Fatalf("freshness checking changed the tracked subject:\n%s", output)
	}
}

func TestRequireFreshBenchRefusesUntrustedArtifacts(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(t *testing.T, root string)
	}{
		{"missing executable", func(t *testing.T, root string) { removeFreshnessArtifact(t, filepath.Join(root, "dist", "bench")) }},
		{"missing legacy seal", func(t *testing.T, root string) { removeFreshnessArtifact(t, filepath.Join(root, "dist", "bench.seal")) }},
		{"unreadable seal", func(t *testing.T, root string) {
			chmodFreshnessArtifact(t, filepath.Join(root, "dist", "bench.seal"), 0)
		}},
		{"malformed seal", func(t *testing.T, root string) {
			writeFreshnessArtifact(t, filepath.Join(root, "dist", "bench.seal"), "{}\n", 0o644)
		}},
		{"partial seal", func(t *testing.T, root string) {
			writeFreshnessArtifact(t, filepath.Join(root, "dist", "bench.seal"), `{"schema":`, 0o644)
		}},
		{"source digest mismatch", func(t *testing.T, root string) {
			writeFreshnessArtifact(t, filepath.Join(root, "cmd", "bench", "main.go"), "package main\n\nfunc main() { println(\"changed\") }\n", 0o644)
		}},
		{"binary digest mismatch", func(t *testing.T, root string) {
			writeFreshnessArtifact(t, filepath.Join(root, "dist", "bench"), "#!/bin/sh\nprintf 'old subject output\\n'\n", 0o755)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if os.Getenv(freshnessSubjectModeEnv) == "artifact" {
				RequireFreshBench(t)
				t.Fatal("untrusted subject reached the contract assertion")
			}

			root := freshnessfixture.PublishedSubject(t, "old subject output")
			tc.mutate(t, root)
			output, err := runFreshnessChild(t, "^TestRequireFreshBenchRefusesUntrustedArtifacts$/"+tc.name+"$", root, "artifact", tc.name, root)
			if err == nil {
				t.Fatalf("untrusted %s subject reached the contract assertion:\n%s", tc.name, output)
			}
			requireFreshnessRefusal(t, root, output, "untrusted subject reached the contract assertion", "old subject output")
		})
	}
}

func TestRequireFreshBenchUsesContentRatherThanMtime(t *testing.T) {
	tie := time.Unix(1_700_000_000, 0)
	for _, tc := range []struct {
		name      string
		prepare   func(t *testing.T, root string)
		wantStale bool
	}{
		{"equal mtimes", func(t *testing.T, root string) { setSubjectMtimes(t, root, tie, tie) }, false},
		{"executable newer", func(t *testing.T, root string) { setSubjectMtimes(t, root, tie, tie.Add(time.Second)) }, false},
		{"executable older", func(t *testing.T, root string) { setSubjectMtimes(t, root, tie.Add(time.Second), tie) }, false},
		{"changed source with tied mtimes", func(t *testing.T, root string) {
			setSubjectMtimes(t, root, tie, tie)
			writeFreshnessArtifact(t, filepath.Join(root, "cmd", "bench", "main.go"), "package main\n\nfunc main() { println(\"changed\") }\n", 0o644)
			setFreshnessMtime(t, filepath.Join(root, "cmd", "bench", "main.go"), tie)
		}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if os.Getenv(freshnessSubjectModeEnv) == "mtime" {
				RequireFreshBench(t)
				if os.Getenv(freshnessSubjectCaseEnv) == "changed source with tied mtimes" {
					t.Fatal("tied stale subject reached the contract assertion")
				}
				return
			}

			root := freshnessfixture.PublishedSubject(t, "old subject output")
			tc.prepare(t, root)
			output, err := runFreshnessChild(t, "^TestRequireFreshBenchUsesContentRatherThanMtime$/"+tc.name+"$", root, "mtime", tc.name, root)
			if tc.wantStale {
				if err == nil {
					t.Fatalf("tied stale subject reached the contract assertion:\n%s", output)
				}
				requireFreshnessRefusal(t, root, output, "tied stale subject reached the contract assertion", "old subject output")
				return
			}
			if err != nil {
				t.Fatalf("fresh subject with %s: %v\n%s", tc.name, err, output)
			}
		})
	}
}

func runFreshnessChild(t *testing.T, pattern, root, mode, name, dir string) ([]byte, error) {
	t.Helper()
	command := exec.Command(os.Args[0], "-test.run="+pattern)
	command.Dir = dir
	command.Env = append(os.Environ(), freshnessSubjectModeEnv+"="+mode, freshnessSubjectCaseEnv+"="+name, canary.SubjectRootEnv+"="+root)
	return command.CombinedOutput()
}

func requireFreshnessRefusal(t *testing.T, root string, output []byte, forbidden ...string) {
	t.Helper()
	if !strings.Contains(string(output), "rebuild with "+freshness.RebuildAction(root)) {
		t.Fatalf("untrusted subject did not report the rebuild action:\n%s", output)
	}
	for _, marker := range forbidden {
		if strings.Contains(string(output), marker) {
			t.Fatalf("untrusted subject reached %q:\n%s", marker, output)
		}
	}
}

func setSubjectMtimes(t *testing.T, root string, source, executable time.Time) {
	t.Helper()
	for _, path := range []string{
		"go.mod",
		"cmd/bench/main.go",
		"scripts/go-build.sh",
		"scripts/go-build.inputs",
		"package.json",
		"internal/releaseevidence/requirements.json",
	} {
		setFreshnessMtime(t, filepath.Join(root, filepath.FromSlash(path)), source)
	}
	setFreshnessMtime(t, filepath.Join(root, "dist", "bench"), executable)
}

func setFreshnessMtime(t *testing.T, path string, timestamp time.Time) {
	t.Helper()
	if err := os.Chtimes(path, timestamp, timestamp); err != nil {
		t.Fatalf("set mtime %s: %v", path, err)
	}
}

func removeFreshnessArtifact(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove %s: %v", path, err)
	}
}

func chmodFreshnessArtifact(t *testing.T, path string, mode os.FileMode) {
	t.Helper()
	if err := os.Chmod(path, mode); err != nil {
		t.Fatalf("chmod %s: %v", path, err)
	}
}

func writeFreshnessArtifact(t *testing.T, path, contents string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), mode); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
