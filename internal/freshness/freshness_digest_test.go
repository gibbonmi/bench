// Tests for build-input discovery and the digest API of package freshness.
package freshness

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
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

// TestDeclaresBuildInputsReadsPresenceRatherThanContent grades the applicability gate
// the landing consults before it proves its own executable. Only a not-exist manifest
// may report absence: every other artifact state routes to Verify, whose reading
// discipline refuses what it cannot trust. An implementation that followed the link —
// Stat rather than Lstat — would read a dangling symlink as an authoritative absence.
// It would then skip the proof in exactly the repository that needs it.
func TestDeclaresBuildInputsReadsPresenceRatherThanContent(t *testing.T) {
	for _, tc := range []struct {
		name    string
		place   func(*testing.T, string)
		declare bool
	}{
		{name: "absent", place: func(*testing.T, string) {}},
		{name: "regular", declare: true, place: func(t *testing.T, path string) {
			if err := os.WriteFile(path, []byte("build_script=scripts/go-build.sh\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "empty", declare: true, place: func(t *testing.T, path string) {
			if err := os.WriteFile(path, nil, 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "live-symlink", declare: true, place: func(t *testing.T, path string) {
			target := filepath.Join(filepath.Dir(path), "elsewhere.inputs")
			if err := os.WriteFile(target, []byte("build_script=scripts/go-build.sh\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, path); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "dangling-symlink", declare: true, place: func(t *testing.T, path string) {
			if err := os.Symlink(filepath.Join(filepath.Dir(path), "no-such-file"), path); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "directory", declare: true, place: func(t *testing.T, path string) {
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
		}},
		// An unusable parent is the case that separates "not-exist" from "errored at all".
		// A predicate that asked only whether Lstat succeeded would read this repository as
		// declaring nothing and skip the proof. That is the one direction that must never
		// happen: an unreadable tree is untrusted, not exempt.
		{name: "parent-is-not-a-directory", declare: true, place: func(t *testing.T, path string) {
			parent := filepath.Dir(path)
			if err := os.Remove(parent); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(parent, []byte("not a directory\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			manifest := filepath.Join(root, filepath.FromSlash(auxiliaryInputsManifest))
			if err := os.MkdirAll(filepath.Dir(manifest), 0o755); err != nil {
				t.Fatal(err)
			}
			tc.place(t, manifest)
			if got := DeclaresBuildInputs(root); got != tc.declare {
				t.Fatalf("DeclaresBuildInputs(%s manifest) = %v, want %v", tc.name, got, tc.declare)
			}
		})
	}
}
