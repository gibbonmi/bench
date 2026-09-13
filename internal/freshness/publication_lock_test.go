package freshness

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/brokermanifest"
)

const publicationProcessEnv = "BENCH_TEST_PUBLICATION_PROCESS"

func TestPublishRefusesInvalidManifestDirectoryBeforeLiveMutation(t *testing.T) {
	for _, kind := range []string{"missing", "regular file", "FIFO", "unreadable"} {
		t.Run(kind, func(t *testing.T) {
			root, executable := writePublishedFixture(t)
			prior := publishedBytes(t, root, executable)
			directory := filepath.Join(root, "invalid-manifest-directory")
			switch kind {
			case "regular file":
				if err := os.WriteFile(directory, []byte("not a directory"), 0o644); err != nil {
					t.Fatal(err)
				}
			case "FIFO":
				if err := syscall.Mkfifo(directory, 0o600); err != nil {
					t.Fatal(err)
				}
			case "unreadable":
				directory = fixtureManifestDir(t, root)
				if err := os.Chmod(directory, 0o333); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = os.Chmod(directory, 0o755) })
			}
			staged := filepath.Join(root, "invalid-directory-stage")
			if err := os.WriteFile(staged, []byte("replacement"), 0o755); err != nil {
				t.Fatal(err)
			}
			original := replacePublicationFile
			replacePublicationFile = func(old, new string) error {
				t.Errorf("invalid manifest directory reached live mutation at %q", new)
				return original(old, new)
			}
			t.Cleanup(func() { replacePublicationFile = original })
			if err := Publish(root, staged, executable, directory, "replacement"); err == nil {
				t.Fatal("invalid manifest directory was accepted")
			}
			if kind == "unreadable" {
				if err := os.Chmod(directory, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			assertPublicationBytes(t, prior)
			if contents, err := os.ReadFile(staged); err != nil || string(contents) != "replacement" {
				t.Fatalf("refusal changed staged executable: %q, %v", contents, err)
			}
			if residue := publicationResidue(t, executable, brokerManifestPath(fixtureManifestDir(t, root))); len(residue) != 0 {
				t.Fatalf("refusal left publication residue: %v", residue)
			}
		})
	}
}

type publicationProcess struct {
	cmd    *exec.Cmd
	input  io.WriteCloser
	events chan string
	exited chan error
	done   chan struct{}
	stderr bytes.Buffer
}

func TestPublishSerializesManifestDirectory(t *testing.T) {
	root, executable := writePublishedFixture(t)
	first := startLockedPublisher(t, root, executable, "first")
	first.expect(t, "mutation")
	second := startLockedPublisher(t, root, executable, "second")
	second.expect(t, "contended")
	first.proceed(t)
	first.expect(t, "seal")
	first.proceed(t)
	first.finish(t, 0)
	second.expect(t, "mutation")
	second.proceed(t)
	second.expect(t, "seal")
	second.proceed(t)
	second.finish(t, 0)
	assertPublishedVersion(t, root, executable, "second")
}

func TestPublishWaiterSeesRestoredTripleAfterInterruption(t *testing.T) {
	root, executable := writePublishedFixture(t)
	prior := publishedBytes(t, root, executable)
	first := startLockedPublisher(t, root, executable, "interrupted")
	first.expect(t, "mutation")
	first.proceed(t)
	first.expect(t, "seal")
	second := startLockedPublisher(t, root, executable, "waiter")
	second.expect(t, "contended")
	if err := first.cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	second.expect(t, "mutation")
	assertPublicationBytes(t, prior)
	first.finish(t, 130)
	second.proceed(t)
	second.expect(t, "seal")
	second.proceed(t)
	second.finish(t, 0)
	assertPublishedVersion(t, root, executable, "waiter")
}

func TestPublicationLockDiesWithProcess(t *testing.T) {
	root, executable := writePublishedFixture(t)
	prior := publishedBytes(t, root, executable)
	first := startLockedPublisher(t, root, executable, "dead")
	first.expect(t, "locked")
	second := startLockedPublisher(t, root, executable, "survivor")
	second.expect(t, "contended")
	if err := first.cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	first.finish(t, -1)
	second.expect(t, "mutation")
	assertPublicationBytes(t, prior)
	second.proceed(t)
	second.expect(t, "seal")
	second.proceed(t)
	second.finish(t, 0)
	assertPublishedVersion(t, root, executable, "survivor")
	entries, err := os.ReadDir(fixtureManifestDir(t, root))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != brokermanifest.Name {
		t.Fatalf("manifest directory has a lock artifact: %v", entries)
	}
}

func TestPublicationProcess(t *testing.T) {
	if os.Getenv(publicationProcessEnv) == "" {
		return
	}
	args := os.Args[len(os.Args)-3:]
	root, executable, version := args[0], args[1], args[2]
	input := bufio.NewScanner(os.Stdin)
	wait := func(event string) {
		fmt.Println(event)
		if !input.Scan() {
			t.Fatalf("%s handshake ended: %v", event, input.Err())
		}
	}
	// The syscall adapter reports actual kernel contention before it waits.
	// Every successful acquisition still comes from the real kernel call.
	publicationFlock = func(fd, operation int) error {
		err := syscall.Flock(fd, operation|syscall.LOCK_NB)
		if errors.Is(err, syscall.EWOULDBLOCK) {
			fmt.Println("contended")
			return syscall.Flock(fd, operation)
		}
		if err == nil && version == "dead" {
			wait("locked")
		}
		return err
	}
	replacePublicationFile = func(old, new string) error {
		if new == executable {
			wait("mutation")
		}
		if new == sealPath(executable) {
			wait("seal")
		}
		return os.Rename(old, new)
	}
	staged := filepath.Join(root, "stage-"+version)
	if err := os.WriteFile(staged, []byte(version), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Publish(root, staged, executable, fixtureManifestDir(t, root), version); err != nil {
		t.Fatal(err)
	}
}

func startLockedPublisher(t *testing.T, root, executable, version string) *publicationProcess {
	t.Helper()
	p := &publicationProcess{events: make(chan string, 16), exited: make(chan error, 1), done: make(chan struct{})}
	p.cmd = exec.Command(os.Args[0], "-test.run=^TestPublicationProcess$", "--", root, executable, version)
	p.cmd.Env = append(os.Environ(), publicationProcessEnv+"=1")
	p.cmd.Stderr = &p.stderr
	var err error
	p.input, err = p.cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	output, err := p.cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := p.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	go func() {
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			p.events <- scanner.Text()
		}
		close(p.events)
		p.exited <- p.cmd.Wait()
		close(p.done)
	}()
	t.Cleanup(func() {
		_ = p.input.Close()
		_ = p.cmd.Process.Kill()
		select {
		case <-p.done:
		case <-time.After(publicationInterruptDeadline):
			t.Error("publisher cleanup did not finish before its deadline")
		}
	})
	return p
}

func (p *publicationProcess) expect(t *testing.T, want string) {
	t.Helper()
	select {
	case got := <-p.events:
		if got != want {
			t.Fatalf("publication event = %q, want %q", got, want)
		}
	case <-time.After(publicationInterruptDeadline):
		t.Fatalf("publication did not reach %q before its deadline", want)
	}
}

func (p *publicationProcess) proceed(t *testing.T) {
	t.Helper()
	if _, err := io.WriteString(p.input, "continue\n"); err != nil {
		t.Fatal(err)
	}
}

func (p *publicationProcess) finish(t *testing.T, code int) {
	t.Helper()
	select {
	case err := <-p.exited:
		got := 0
		if err != nil {
			var exit *exec.ExitError
			if !errors.As(err, &exit) {
				t.Fatal(err)
			}
			got = exit.ExitCode()
		}
		if got != code {
			t.Fatalf("publisher exit = %d, want %d: %s", got, code, p.stderr.String())
		}
	case <-time.After(publicationInterruptDeadline):
		t.Fatal("publisher did not exit before its deadline")
	}
}

func assertPublishedVersion(t *testing.T, root, executable, version string) {
	t.Helper()
	contents, err := os.ReadFile(executable)
	if err != nil || string(contents) != version {
		t.Fatalf("executable = %q, %v; want %q", contents, err, version)
	}
	if err := Verify(root, executable); err != nil {
		t.Fatal(err)
	}
	manifest := brokerManifestPath(fixtureManifestDir(t, root))
	fields, err := brokermanifest.Read(manifest)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := brokermanifest.Digest(executable)
	if err != nil {
		t.Fatal(err)
	}
	if fields["path"] != executable || fields["version"] != version || fields["sha256"] != digest {
		t.Fatalf("manifest = %v; want complete version %q at %q", fields, version, executable)
	}
	if residue := publicationResidue(t, executable, manifest); len(residue) != 0 {
		t.Fatalf("publication residue = %v", residue)
	}
}

func publishedBytes(t *testing.T, root, executable string) map[string][]byte {
	t.Helper()
	prior := make(map[string][]byte)
	for _, path := range []string{executable, sealPath(executable), brokerManifestPath(fixtureManifestDir(t, root))} {
		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		prior[path] = contents
	}
	return prior
}

func assertPublicationBytes(t *testing.T, prior map[string][]byte) {
	t.Helper()
	for path, want := range prior {
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("prior publication changed at %q: %q, %v", path, got, err)
		}
	}
}
