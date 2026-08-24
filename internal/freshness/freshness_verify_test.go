// Tests for seal verification, staleness checks, and refusal messaging of package freshness.
package freshness

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestVerifyRefusesSymbolicLinkBeforeReadingExecutable(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	if err := os.WriteFile(target, []byte("not a Bench executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(root, "dist", "bench")
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, executable); err != nil {
		t.Fatal(err)
	}

	err := Verify(root, executable)
	if err == nil || !strings.Contains(err.Error(), "symbolic link") {
		t.Fatalf("Verify symlinked executable error = %v, want symbolic-link refusal", err)
	}
}

func TestVerifyUsesContentRatherThanMtime(t *testing.T) {
	root, executable := writePublishedFixture(t)
	inputs, err := buildInputs(root)
	if err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(root, "internal", "local", "local.go")
	tie := time.Unix(1_700_000_000, 0)
	setMtimes(t, inputs, tie)
	setMtimes(t, []string{executable}, tie)
	if err := Verify(root, executable); err != nil {
		t.Fatalf("Verify equal mtime: %v", err)
	}
	setMtimes(t, inputs, tie)
	setMtimes(t, []string{executable}, tie.Add(time.Second))
	if err := Verify(root, executable); err != nil {
		t.Fatalf("Verify executable newer than every input: %v", err)
	}
	setMtimes(t, inputs, tie.Add(time.Second))
	setMtimes(t, []string{executable}, tie)
	if err := Verify(root, executable); err != nil {
		t.Fatalf("Verify executable older than every input: %v", err)
	}
	setMtimes(t, inputs, tie)
	setMtimes(t, []string{executable}, tie)
	if err := os.WriteFile(source, []byte("package local\n\nconst Value = \"changed\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	setMtimes(t, []string{source}, tie)
	if err := Verify(root, executable); err == nil {
		t.Fatal("Verify changed source with tied mtime = nil, want stale refusal")
	}
}

func TestVerifyRefusesUntrustedArtifactStates(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(t *testing.T, executable string)
	}{
		{
			name: "missing executable",
			mutate: func(t *testing.T, executable string) {
				t.Helper()
				if err := os.Remove(executable); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "missing seal",
			mutate: func(t *testing.T, executable string) {
				t.Helper()
				if err := os.Remove(sealPath(executable)); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "malformed complete seal",
			mutate: func(t *testing.T, executable string) {
				t.Helper()
				if err := os.WriteFile(sealPath(executable), []byte("{}\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "malformed partial seal",
			mutate: func(t *testing.T, executable string) {
				t.Helper()
				if err := os.WriteFile(sealPath(executable), []byte(`{"schema":`), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "unreadable seal",
			mutate: func(t *testing.T, executable string) {
				t.Helper()
				if err := os.Chmod(sealPath(executable), 0); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "executable digest mismatch",
			mutate: func(t *testing.T, executable string) {
				t.Helper()
				if err := os.WriteFile(executable, []byte("altered executable"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "nonexecutable binary",
			mutate: func(t *testing.T, executable string) {
				t.Helper()
				if err := os.Chmod(executable, 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "unreadable executable",
			mutate: func(t *testing.T, executable string) {
				t.Helper()
				if err := os.Chmod(executable, 0o111); err != nil {
					t.Fatal(err)
				}
			},
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			root, executable := writePublishedFixture(t)
			test.mutate(t, executable)
			assertRefusal(t, root, executable)
		})
	}
}

func TestCheckRefusesFailingVerifiedExecutable(t *testing.T) {
	if root := os.Getenv("BENCH_TEST_CHECK_FAILURE_ROOT"); root != "" {
		if err := Check(root, os.Getenv("BENCH_TEST_CHECK_FAILURE_EXECUTABLE")); err == nil {
			t.Fatal("Check failing verified executable = nil, want refusal")
		}
		return
	}

	root := writeBuildFixture(t)
	staged := filepath.Join(root, "staged-bench")
	sentinel := "untrusted child output"
	program := "#!/usr/bin/env bash\nprintf '%s\\n' " + quoteForShell(sentinel) + "\nprintf '%s\\n' " + quoteForShell(sentinel) + " >&2\nexit 2\n"
	if err := os.WriteFile(staged, []byte(program), 0o755); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(root, "dist", "bench")
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Publish(root, staged, executable); err != nil {
		t.Fatal(err)
	}

	err := Check(root, executable)
	if err == nil {
		t.Fatal("Check failing verified executable = nil, want refusal")
	}
	if !strings.Contains(err.Error(), "freshness-check failed") {
		t.Fatalf("Check failing verified executable error = %q, want freshness-check refusal", err)
	}
	if strings.Contains(err.Error(), sentinel) {
		t.Fatalf("Check leaked child output in refusal: %q", err)
	}
	command := exec.Command(os.Args[0], "-test.run", "^TestCheckRefusesFailingVerifiedExecutable$")
	command.Env = append(os.Environ(),
		"BENCH_TEST_CHECK_FAILURE_ROOT="+root,
		"BENCH_TEST_CHECK_FAILURE_EXECUTABLE="+executable,
	)
	output, runErr := command.CombinedOutput()
	if runErr != nil {
		t.Fatalf("subprocess Check failure journey: %v\n%s", runErr, output)
	}
	if strings.Contains(string(output), sentinel) {
		t.Fatalf("Check exposed child output:\n%s", output)
	}
	if got, want := strings.Count(err.Error(), "rebuild with "), 1; got != want {
		t.Fatalf("Check refusal rebuild-action count = %d, want %d: %q", got, want, err)
	}
	if !strings.Contains(err.Error(), RebuildAction(root)) {
		t.Fatalf("Check failing verified executable error = %q, want rebuild action %q", err, RebuildAction(root))
	}
}

func TestVerifyRefusesLiveAndDanglingSymlinkArtifacts(t *testing.T) {
	for _, artifact := range []struct {
		name string
		path func(string) string
	}{
		{"executable", func(executable string) string { return executable }},
		{"seal", sealPath},
	} {
		for _, target := range []struct {
			name string
			live bool
		}{
			{"live", true},
			{"dangling", false},
		} {
			t.Run(artifact.name+" "+target.name, func(t *testing.T) {
				root, executable := writePublishedFixture(t)
				path := artifact.path(executable)
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				linkTarget := filepath.Join(root, "link-target")
				if target.live {
					if err := os.WriteFile(linkTarget, []byte("untrusted"), 0o755); err != nil {
						t.Fatal(err)
					}
				}
				if err := os.Symlink(linkTarget, path); err != nil {
					t.Fatal(err)
				}
				assertRefusal(t, root, executable)
			})
		}
	}
}

func TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths(t *testing.T) {
	root := filepath.Join(t.TempDir(), "root [*] with ' quote")
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	record := filepath.Join(t.TempDir(), "record")
	script := "#!/usr/bin/env bash\nprintf '%s\\n%s\\n%s\\n' \"$PWD\" \"$1\" \"$2\" > " + quoteForShell(record) + "\n"
	if err := os.WriteFile(filepath.Join(root, "scripts", "go-build.sh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	err := refusal(root, filepath.Join(root, "dist", "bench"), fmt.Errorf("fixture"))
	action := strings.TrimPrefix(err.Error()[strings.Index(err.Error(), "; rebuild with "):], "; rebuild with ")
	command := exec.Command("bash", "-c", action)
	command.Dir = t.TempDir()
	if output, runErr := command.CombinedOutput(); runErr != nil {
		t.Fatalf("copy-paste rebuild action failed: %v\n%s", runErr, output)
	}
	want := strings.Join([]string{root, root, filepath.Join(root, "dist", "bench"), ""}, "\n")
	if got, readErr := os.ReadFile(record); readErr != nil || string(got) != want {
		t.Fatalf("rebuild action args = %q, %v; want %q", got, readErr, want)
	}
}

func quoteForShell(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func TestVerifyRefusesSpecialArtifactsBeforeReading(t *testing.T) {
	for _, artifact := range []struct {
		name string
		path func(string) string
	}{
		{"executable", func(executable string) string { return executable }},
		{"seal", sealPath},
	} {
		for _, special := range []struct {
			name   string
			create func(string) (func(), error)
		}{
			{
				name: "FIFO",
				create: func(path string) (func(), error) {
					return func() {}, syscall.Mkfifo(path, 0o600)
				},
			},
			{
				name: "Unix socket",
				create: func(path string) (func(), error) {
					listener, err := net.Listen("unix", path)
					if err != nil {
						return nil, err
					}
					return func() { _ = listener.Close() }, nil
				},
			},
		} {
			t.Run(artifact.name+" "+special.name, func(t *testing.T) {
				root, executable := writePublishedFixture(t)
				path := artifact.path(executable)
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				cleanup, err := special.create(path)
				if err != nil {
					t.Fatal(err)
				}
				defer cleanup()
				assertRefusal(t, root, executable)
			})
		}
	}
}

func TestVerifyRefusesSymbolicLinkSeal(t *testing.T) {
	root, executable := writePublishedFixture(t)
	if err := os.Remove(sealPath(executable)); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "target"), sealPath(executable)); err != nil {
		t.Fatal(err)
	}
	assertRefusal(t, root, executable)
}

func TestVerifyRefusesNamedPipeSealBeforeReading(t *testing.T) {
	root, executable := writePublishedFixture(t)
	if err := os.Remove(sealPath(executable)); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(sealPath(executable), 0o600); err != nil {
		t.Fatal(err)
	}
	assertRefusal(t, root, executable)
}

func writePublishedFixture(t *testing.T) (string, string) {
	t.Helper()
	root := writeBuildFixture(t)
	staged := filepath.Join(root, "staged-bench")
	if err := os.WriteFile(staged, []byte("Bench executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(root, "dist", "bench")
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Publish(root, staged, executable); err != nil {
		t.Fatal(err)
	}
	return root, executable
}

func assertRefusal(t *testing.T, root, executable string) {
	t.Helper()
	err := Verify(root, executable)
	if err == nil {
		t.Fatal("Verify untrusted artifact = nil, want refusal")
	}
	if !strings.Contains(err.Error(), executable) {
		t.Fatalf("Verify error = %q, want binary path %q", err, executable)
	}
	want := RebuildAction(root)
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("Verify error = %q, want rebuild action %q", err, want)
	}
	if got, want := strings.Count(err.Error(), "rebuild with "), 1; got != want {
		t.Fatalf("Verify refusal rebuild-action count = %d, want %d: %q", got, want, err)
	}
}

func setMtimes(t *testing.T, paths []string, timestamp time.Time) {
	t.Helper()
	for _, path := range paths {
		if err := os.Chtimes(path, timestamp, timestamp); err != nil {
			t.Fatal(err)
		}
	}
}
