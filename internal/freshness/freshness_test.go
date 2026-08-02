package freshness

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func TestBuildInputsCoverEveryDigestedPath(t *testing.T) {
	root := writeBuildFixture(t)
	paths, err := BuildInputs(root)
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.New()
	for _, name := range paths {
		contents, err := regularContents(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(hash, "%d:%s%d:", len(name), name, len(contents))
		if _, err := hash.Write(contents); err != nil {
			t.Fatal(err)
		}
	}
	got := hex.EncodeToString(hash.Sum(nil))
	want, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("BuildInputs paths rehashed = %q, want Digest(root) = %q", got, want)
	}
}

func TestBuildInputsErrorOnUnresolvableClosure(t *testing.T) {
	root := t.TempDir()
	paths, err := BuildInputs(root)
	if err == nil {
		t.Fatalf("BuildInputs with no go.mod = %v, %v; want error and no paths", paths, err)
	}
	if paths != nil {
		t.Fatalf("BuildInputs with no go.mod returned paths %v, want nil", paths)
	}
}

func TestSealDigestsReturnsPublishedDigests(t *testing.T) {
	root, executable := writePublishedFixture(t)
	wantSources, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	binary, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	sources, digest, err := SealDigests(executable)
	if err != nil {
		t.Fatal(err)
	}
	if sources != wantSources {
		t.Fatalf("SealDigests sources = %q, want %q", sources, wantSources)
	}
	if digest != digestBytes(binary) {
		t.Fatalf("SealDigests executable digest = %q, want %q", digest, digestBytes(binary))
	}
}

func TestSealDigestsRefuseAnUntrustedSidecar(t *testing.T) {
	t.Run("symlinked sidecar", func(t *testing.T) {
		root, executable := writePublishedFixture(t)
		target := filepath.Join(root, "valid.seal")
		data, err := os.ReadFile(sealPath(executable))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, data, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(sealPath(executable)); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, sealPath(executable)); err != nil {
			t.Fatal(err)
		}
		if _, _, err := SealDigests(executable); err == nil {
			t.Fatal("SealDigests through symlinked sidecar = nil, want refusal")
		}
	})

	t.Run("irregular sidecar", func(t *testing.T) {
		_, executable := writePublishedFixture(t)
		if err := os.Remove(sealPath(executable)); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(sealPath(executable), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, _, err := SealDigests(executable); err == nil {
			t.Fatal("SealDigests on irregular sidecar = nil, want refusal")
		}
	})

	t.Run("malformed contents", func(t *testing.T) {
		_, executable := writePublishedFixture(t)
		if err := os.WriteFile(sealPath(executable), []byte(`{"schema":`), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, _, err := SealDigests(executable); err == nil {
			t.Fatal("SealDigests with malformed seal contents = nil, want refusal")
		}
	})
}

func TestPublishAndVerifyRefuseEmptyExecutables(t *testing.T) {
	t.Run("empty stage", func(t *testing.T) {
		root := writeBuildFixture(t)
		staged := filepath.Join(root, "empty-stage")
		if err := os.WriteFile(staged, nil, 0o755); err != nil {
			t.Fatal(err)
		}
		executable := filepath.Join(root, "dist", "bench")
		if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := Publish(root, staged, executable); err == nil {
			t.Fatal("Publish empty executable stage = nil, want refusal")
		}
	})

	t.Run("sealed empty artifact", func(t *testing.T) {
		root, executable := writePublishedFixture(t)
		if err := os.WriteFile(executable, nil, 0o755); err != nil {
			t.Fatal(err)
		}
		sources, err := Digest(root)
		if err != nil {
			t.Fatal(err)
		}
		encoded, err := json.Marshal(seal{
			Schema:     sealSchema,
			Sources:    sources,
			Executable: digestBytes(nil),
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(sealPath(executable), encoded, 0o644); err != nil {
			t.Fatal(err)
		}
		assertRefusal(t, root, executable)
	})
}

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
	want := RebuildAction(root)
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

func TestDigestIgnoresMalformedAmbientVCSMetadata(t *testing.T) {
	ancestor := t.TempDir()
	if err := os.WriteFile(filepath.Join(ancestor, ".git"), []byte("gitdir: /does/not/exist\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	root := writeBuildFixtureAt(t, filepath.Join(ancestor, "fixture"))

	if _, err := Digest(root); err != nil {
		command := exec.Command("go", "list", "-json", "-deps", "./cmd/bench")
		command.Dir = root
		output, diagnosticErr := command.CombinedOutput()
		t.Fatalf("Digest with malformed ambient VCS metadata: %v\nunprotected go list: %v\n%s", err, diagnosticErr, output)
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

func TestDigestTracksInputsDeclaredByTheBuildOwnerManifest(t *testing.T) {
	root := writeBuildFixture(t)
	unrelated := filepath.Join(root, "declared-later.txt")
	if err := os.WriteFile(unrelated, []byte("one"), 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(unrelated, []byte("two"), 0o644); err != nil {
		t.Fatal(err)
	}
	stillUnrelated, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if before != stillUnrelated {
		t.Fatal("Digest changed before an auxiliary input was declared")
	}
	manifest := filepath.Join(root, filepath.FromSlash(auxiliaryInputsManifest))
	file, err := os.OpenFile(manifest, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("future_input=declared-later.txt\n"); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	declared, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if declared == stillUnrelated {
		t.Fatal("Digest ignored an input added to the build owner manifest")
	}
	if err := os.WriteFile(unrelated, []byte("three"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if changed == declared {
		t.Fatal("Digest ignored a changed declared auxiliary input")
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

func TestSecureContentsRefusesSymlinkedAncestor(t *testing.T) {
	for _, executable := range []bool{false, true} {
		t.Run(map[bool]string{false: "seal", true: "executable"}[executable], func(t *testing.T) {
			root := t.TempDir()
			realDir := filepath.Join(root, "real")
			if err := os.MkdirAll(realDir, 0o755); err != nil {
				t.Fatal(err)
			}
			mode := os.FileMode(0o644)
			if executable {
				mode = 0o755
			}
			if err := os.WriteFile(filepath.Join(realDir, "artifact"), []byte("untrusted"), mode); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(realDir, filepath.Join(root, "linked")); err != nil {
				t.Fatal(err)
			}
			if _, err := secureContents(filepath.Join(root, "linked", "artifact"), executable); err == nil || !strings.Contains(err.Error(), "symbolic link") {
				t.Fatalf("secureContents through symlinked ancestor error = %v, want symbolic-link refusal", err)
			}
		})
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
	return writeBuildFixtureAt(t, t.TempDir())
}

func writeBuildFixtureAt(t *testing.T, root string) string {
	t.Helper()
	files := map[string]string{
		"go.mod":                  "module example.com/freshnessfixture\n\ngo 1.25\n",
		"cmd/bench/main.go":       "package main\n\nimport \"example.com/freshnessfixture/internal/local\"\n\nfunc main() { _ = local.Value }\n",
		"internal/local/local.go": "package local\n\nconst Value = \"value\"\n",
		"scripts/go-build.sh":     "#!/usr/bin/env bash\n",
		"scripts/go-build.inputs": "build_script=scripts/go-build.sh\npackage_version=package.json\ngo_requirements=internal/releaseevidence/requirements.json\n",
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
