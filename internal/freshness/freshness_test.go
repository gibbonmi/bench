package freshness

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestVerifyMissingExecutableRefusesWithRebuild(t *testing.T) {
	root := t.TempDir()
	executable := filepath.Join(root, "dist", "bench")

	err := Verify(root, executable)
	if err == nil {
		t.Fatal("Verify missing executable = nil, want actionable refusal")
	}
	if !strings.Contains(err.Error(), executable) {
		t.Fatalf("Verify missing executable error = %q, want binary path %q", err, executable)
	}
	want := "bash scripts/go-build.sh " + root + " " + filepath.Join(root, "dist", "bench")
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("Verify missing executable error = %q, want rebuild action %q", err, want)
	}
}

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

func TestDigestTracksResolvedUntrackedBuildInput(t *testing.T) {
	root := writeBuildFixture(t)
	beforeUnrelated, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "unrelated.txt"), []byte("outside the build graph"), 0o644); err != nil {
		t.Fatal(err)
	}
	afterUnrelated, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if beforeUnrelated != afterUnrelated {
		t.Fatal("Digest changed for an unrelated file")
	}
	generated := filepath.Join(root, "internal", "local", "generated.go")
	if err := os.WriteFile(generated, []byte("package local\n\nconst Generated = \"one\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(generated, []byte("package local\n\nconst Generated = \"two\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	after, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if before == after {
		t.Fatal("Digest ignored a resolved untracked Go input")
	}
}

func TestDigestTracksBuildOwnerInputs(t *testing.T) {
	cases := []struct {
		name string
		path string
		body string
	}{
		{"module definition", "go.mod", "module example.com/freshnessfixture\n\ngo 1.25\n\n"},
		{"optional module checksum", "go.sum", ""},
		{"build script", "scripts/go-build.sh", "#!/usr/bin/env bash\n# changed\n"},
		{"package version", "package.json", "{\"version\":\"0.0.1\"}\n"},
		{"build flags registry", "internal/releaseevidence/requirements.json", "{\"changed\":true}\n"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			root := writeBuildFixture(t)
			before, err := Digest(root)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, test.path), []byte(test.body), 0o644); err != nil {
				t.Fatal(err)
			}
			after, err := Digest(root)
			if err != nil {
				t.Fatal(err)
			}
			if before == after {
				t.Fatalf("Digest ignored %s", test.name)
			}
		})
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
			name: "missing seal",
			mutate: func(t *testing.T, executable string) {
				t.Helper()
				if err := os.Remove(sealPath(executable)); err != nil {
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
	want := "bash scripts/go-build.sh " + root + " " + filepath.Join(root, "dist", "bench")
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("Verify error = %q, want rebuild action %q", err, want)
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

func writeBuildFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"go.mod":                  "module example.com/freshnessfixture\n\ngo 1.25\n",
		"cmd/bench/main.go":       "package main\n\nimport \"example.com/freshnessfixture/internal/local\"\n\nfunc main() { _ = local.Value }\n",
		"internal/local/local.go": "package local\n\nconst Value = \"value\"\n",
		"scripts/go-build.sh":     "#!/usr/bin/env bash\n",
		"package.json":            "{\"version\":\"0.0.0\"}\n",
		"internal/releaseevidence/requirements.json": "{}\n",
	}
	for name, body := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
