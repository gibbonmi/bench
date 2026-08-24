// Tests for the atomic publication lifecycle of package freshness.
package freshness

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

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
// elapsed. Every deadline the parent holds derives from that bound instead of restating
// a wall-clock guess that could silently fall inside it.
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
// child process. Answering a termination signal ends the process that received it, so an
// in-process exercise would take the test binary down with it. The child holds the
// publication open at the one point where the new executable is installed and the seal
// still describes the old one. That is the only state a signal can leave inconsistent.
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
// restored. The backups of the pair it replaced are still on disk. The child drives the
// transaction directly and signals itself, because that window closes as soon as the
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

// publicationResidue lists every temporary a publication of executable can leave beside
// it, matched by the publisher's own naming rather than a glob restated here. A new
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
