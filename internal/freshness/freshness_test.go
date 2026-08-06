package freshness

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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

func TestPublishRefusesSymlinkedDestinationWithoutChangingPublishedPair(t *testing.T) {
	root, executable := writePublishedFixture(t)
	beforeSeal, err := os.ReadFile(sealPath(executable))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(executable); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "replacement"), executable); err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(root, "staged-replacement")
	if err := os.WriteFile(staged, []byte("replacement executable"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := Publish(root, staged, executable); err == nil {
		t.Fatal("Publish symlinked destination = nil, want refusal")
	}
	if _, err := os.Lstat(staged); err != nil {
		t.Fatalf("Publish removed staged executable on refusal: %v", err)
	}
	if afterSeal, err := os.ReadFile(sealPath(executable)); err != nil || !bytes.Equal(afterSeal, beforeSeal) {
		t.Fatalf("Publish changed old seal on refusal: %v, %q", err, afterSeal)
	}
	if target, err := os.Readlink(executable); err != nil || target != filepath.Join(root, "replacement") {
		t.Fatalf("Publish replaced symlinked destination: target=%q err=%v", target, err)
	}
}

func TestPublishRestoresPriorPairWhenSealPromotionFails(t *testing.T) {
	root, executable := writePublishedFixture(t)
	beforeExecutable, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	beforeSeal, err := os.ReadFile(sealPath(executable))
	if err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(root, "replacement")
	if err := os.WriteFile(staged, []byte("replacement executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	originalRename := replacePublicationFile
	replacePublicationFile = func(old, new string) error {
		if new == sealPath(executable) {
			return errors.New("forced seal promotion failure")
		}
		return originalRename(old, new)
	}
	t.Cleanup(func() { replacePublicationFile = originalRename })

	if err := Publish(root, staged, executable); err == nil {
		t.Fatal("Publish forced seal promotion failure = nil, want refusal")
	}
	afterExecutable, err := os.ReadFile(executable)
	if err != nil || !bytes.Equal(afterExecutable, beforeExecutable) {
		t.Fatalf("Publish changed executable after failed seal promotion: %v, %q", err, afterExecutable)
	}
	afterSeal, err := os.ReadFile(sealPath(executable))
	if err != nil || !bytes.Equal(afterSeal, beforeSeal) {
		t.Fatalf("Publish changed seal after failed seal promotion: %v, %q", err, afterSeal)
	}
	if _, err := os.Lstat(staged); !os.IsNotExist(err) {
		t.Fatalf("Publish left staged executable after restoring the prior pair: %v", err)
	}
	if leftovers := publicationResidue(t, executable); len(leftovers) != 0 {
		t.Fatalf("Publish residue = %v; want none", leftovers)
	}
}

func TestPublishLeavesNoArtifactsWhenFirstSealPromotionFails(t *testing.T) {
	root := writeBuildFixture(t)
	executable := filepath.Join(root, "dist", "bench")
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(root, "first-stage")
	if err := os.WriteFile(staged, []byte("first executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	originalRename := replacePublicationFile
	replacePublicationFile = func(old, new string) error {
		if new == sealPath(executable) {
			return errors.New("forced seal promotion failure")
		}
		return originalRename(old, new)
	}
	t.Cleanup(func() { replacePublicationFile = originalRename })

	if err := Publish(root, staged, executable); err == nil {
		t.Fatal("Publish first forced seal promotion failure = nil, want refusal")
	}
	for _, path := range []string{executable, sealPath(executable), staged} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("Publish first failure left %q: %v", path, err)
		}
	}
	if leftovers := publicationResidue(t, executable); len(leftovers) != 0 {
		t.Fatalf("Publish first-failure residue = %v; want none", leftovers)
	}
}

// The child answers a termination signal only once the publication's own step grace has
// elapsed, so every deadline the parent holds is derived from that bound instead of
// restating a wall-clock guess that could silently fall inside it.
const (
	publicationInterruptDeadline = 5 * publicationStepGrace
	publicationInterruptBlock    = 2 * publicationInterruptDeadline
)

const (
	publicationInterruptRootEnv       = "BENCH_TEST_PUBLICATION_INTERRUPT_ROOT"
	publicationInterruptStagedEnv     = "BENCH_TEST_PUBLICATION_INTERRUPT_STAGED"
	publicationInterruptExecutableEnv = "BENCH_TEST_PUBLICATION_INTERRUPT_EXECUTABLE"
	publicationInterruptReadyEnv      = "BENCH_TEST_PUBLICATION_INTERRUPT_READY"
)

// TestPublishRestoresPriorPairWhenInterruptedBeforeItsSealLands runs the publication in a
// child process: answering a termination signal ends the process that received it, so an
// in-process exercise would take the test binary down with it. The child holds the
// publication open at the one point where the new executable is installed and the seal
// still describes the old one, which is the only state a signal can leave inconsistent.
func TestPublishRestoresPriorPairWhenInterruptedBeforeItsSealLands(t *testing.T) {
	if root := os.Getenv(publicationInterruptRootEnv); root != "" {
		runInterruptedPublicationChild(t, root)
		return
	}
	root, executable := writePublishedFixture(t)
	beforeExecutable, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	beforeSeal, err := os.ReadFile(sealPath(executable))
	if err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(root, "interrupted-replacement")
	if err := os.WriteFile(staged, []byte("interrupted replacement executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	ready := filepath.Join(root, "seal-write-reached")

	cmd, transcript := startPublicationChild(t, "TestPublishRestoresPriorPairWhenInterruptedBeforeItsSealLands", root, staged, executable, ready)
	awaitPublicationMarker(t, ready, cmd)
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	awaitInterruptedExit(t, cmd, transcript)

	afterExecutable, err := os.ReadFile(executable)
	if err != nil || !bytes.Equal(afterExecutable, beforeExecutable) {
		t.Fatalf("interrupted publication changed the executable: %v, %q", err, afterExecutable)
	}
	afterSeal, err := os.ReadFile(sealPath(executable))
	if err != nil || !bytes.Equal(afterSeal, beforeSeal) {
		t.Fatalf("interrupted publication changed the seal: %v, %q", err, afterSeal)
	}
	if leftovers := publicationResidue(t, executable); len(leftovers) != 0 {
		t.Fatalf("interrupted publication residue = %v; want none", leftovers)
	}
}

// runInterruptedPublicationChild publishes with the seal's own promotion held open, so
// the process sits in the window the parent signals. It reaches that window through the
// package's existing publication-file replacement rather than a hook of its own.
func runInterruptedPublicationChild(t *testing.T, root string) {
	t.Helper()
	executable := os.Getenv(publicationInterruptExecutableEnv)
	ready := os.Getenv(publicationInterruptReadyEnv)
	original := replacePublicationFile
	replacePublicationFile = func(old, new string) error {
		if new != sealPath(executable) {
			return original(old, new)
		}
		if err := os.WriteFile(ready, nil, 0o644); err != nil {
			return err
		}
		deadline := time.Now().Add(publicationInterruptBlock)
		for time.Now().Before(deadline) {
			time.Sleep(20 * time.Millisecond)
		}
		return errors.New("held publication was never interrupted")
	}
	t.Cleanup(func() { replacePublicationFile = original })
	if err := Publish(root, os.Getenv(publicationInterruptStagedEnv), executable); err != nil {
		t.Fatalf("held publication returned instead of being interrupted: %v", err)
	}
}

// TestPublishRemovesBackupPairWhenInterruptedAfterItsSealLands covers the other side of a
// handled termination: the seal has landed, so the publication stands and nothing is
// restored, but the backups of the pair it replaced are still on disk. The child drives
// the transaction directly and signals itself, because that window closes as soon as the
// transaction does.
func TestPublishRemovesBackupPairWhenInterruptedAfterItsSealLands(t *testing.T) {
	if root := os.Getenv(publicationInterruptRootEnv); root != "" {
		runSealedPublicationChild(t, root)
		return
	}
	root, executable := writePublishedFixture(t)
	staged := filepath.Join(root, "sealed-replacement")
	replacement := []byte("sealed replacement executable")
	if err := os.WriteFile(staged, replacement, 0o755); err != nil {
		t.Fatal(err)
	}

	cmd, transcript := startPublicationChild(t, "TestPublishRemovesBackupPairWhenInterruptedAfterItsSealLands", root, staged, executable, "")
	awaitInterruptedExit(t, cmd, transcript)

	published, err := os.ReadFile(executable)
	if err != nil || !bytes.Equal(published, replacement) {
		t.Fatalf("interrupted publication undid its landed seal: %v, %q\n%s", err, published, transcript.String())
	}
	if err := Verify(root, executable); err != nil {
		t.Fatalf("interrupted publication left an unverifiable pair: %v\n%s", err, transcript.String())
	}
	if leftovers := publicationResidue(t, executable); len(leftovers) != 0 {
		t.Fatalf("sealed publication residue = %v; want none", leftovers)
	}
}

// runSealedPublicationChild holds the transaction open past its landed seal, so the
// termination it sends itself is answered while the backups of the replaced pair remain.
func runSealedPublicationChild(t *testing.T, root string) {
	t.Helper()
	executable := os.Getenv(publicationInterruptExecutableEnv)
	encoded, err := sealContents(root, os.Getenv(publicationInterruptStagedEnv))
	if err != nil {
		t.Fatal(err)
	}
	transaction, err := beginPublication(executable)
	if err != nil {
		t.Fatal(err)
	}
	if err := transaction.install(os.Getenv(publicationInterruptStagedEnv)); err != nil {
		t.Fatal(err)
	}
	if err := transaction.sealWith(encoded); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Kill(os.Getpid(), syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(publicationInterruptBlock)
	for time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("sealed publication was never interrupted")
}

func startPublicationChild(t *testing.T, name, root, staged, executable, ready string) (*exec.Cmd, *bytes.Buffer) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^"+name+"$")
	cmd.Env = append(os.Environ(),
		publicationInterruptRootEnv+"="+root,
		publicationInterruptStagedEnv+"="+staged,
		publicationInterruptExecutableEnv+"="+executable,
		publicationInterruptReadyEnv+"="+ready,
	)
	var transcript bytes.Buffer
	cmd.Stdout, cmd.Stderr = &transcript, &transcript
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	return cmd, &transcript
}

func awaitInterruptedExit(t *testing.T, cmd *exec.Cmd, transcript *bytes.Buffer) {
	t.Helper()
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	select {
	case err := <-exited:
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 130 {
			t.Fatalf("interrupted publication exit = %v, want code 130\n%s", err, transcript.String())
		}
	case <-time.After(publicationInterruptDeadline):
		_ = cmd.Process.Kill()
		t.Fatalf("interrupted publication did not exit within %s\n%s", publicationInterruptDeadline, transcript.String())
	}
}

func awaitPublicationMarker(t *testing.T, path string, cmd *exec.Cmd) {
	t.Helper()
	deadline := time.Now().Add(publicationInterruptDeadline)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			_ = cmd.Process.Kill()
			t.Fatalf("publication did not reach its seal write within %s", publicationInterruptDeadline)
		}
		time.Sleep(20 * time.Millisecond)
	}
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

// publicationResidue lists every temporary a publication of executable can leave beside
// it, matched by the publisher's own naming rather than a glob restated here, so a new
// temporary family is covered the moment the publisher can create one.
func publicationResidue(t *testing.T, executable string) []string {
	t.Helper()
	var leftovers []string
	for _, pattern := range []string{publicationBackupPattern, sealTemporaryPattern(sealPath(executable))} {
		matches, err := filepath.Glob(filepath.Join(filepath.Dir(executable), pattern))
		if err != nil {
			t.Fatal(err)
		}
		leftovers = append(leftovers, matches...)
	}
	return leftovers
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
