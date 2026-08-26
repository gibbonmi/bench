package landing

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate/authorization"
)

func TestLandPreAuthorizationRefusalTable(t *testing.T) {
	cases := []struct {
		name    string
		setup   func(*testing.T, string) []string
		wantErr string
	}{
		{"unknown-path", func(*testing.T, string) []string { return []string{"missing"} }, "not found in worktree, index, or expected base"},
		{"repository-escape", func(*testing.T, string) []string { return []string{"../outside"} }, "escapes repository"},
		{"repository-root", func(*testing.T, string) []string { return []string{"."} }, "repository root is not an attributed path"},
		{"unreadable-path", func(t *testing.T, root string) []string {
			path := filepath.Join(root, "unreadable")
			write(t, root, "unreadable", "unreadable")
			if err := os.Chmod(path, 0); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.Chmod(path, 0o600) })
			return []string{"unreadable"}
		}, "unreadable path"},
		{"direct-fifo", func(t *testing.T, root string) []string {
			requireLandingFIFO(t, filepath.Join(root, "fifo"))
			return []string{"fifo"}
		}, "special file"},
		{"direct-socket", func(t *testing.T, root string) []string {
			requireLandingSocket(t, filepath.Join(root, "socket"))
			return []string{"socket"}
		}, "special file"},
		{"direct-device", func(t *testing.T, root string) []string {
			requireLandingDevice(t, filepath.Join(root, "device"))
			return []string{"device"}
		}, "special file"},
		{"descendant-fifo", func(t *testing.T, root string) []string {
			requireLandingFIFO(t, filepath.Join(root, "dir", "fifo"))
			return []string{"dir"}
		}, "special or unreadable descendant"},
		{"descendant-socket", func(t *testing.T, root string) []string {
			requireLandingSocket(t, filepath.Join(root, "dir", "socket"))
			return []string{"dir"}
		}, "special or unreadable descendant"},
		{"descendant-device", func(t *testing.T, root string) []string {
			requireLandingDevice(t, filepath.Join(root, "dir", "device"))
			return []string{"dir"}
		}, "special or unreadable descendant"},
		{"symlink-to-special-without-traversal", func(t *testing.T, root string) []string {
			requireLandingFIFO(t, filepath.Join(root, "fifo-target"))
			if err := os.Symlink("fifo-target", filepath.Join(root, "fifo-link")); err != nil {
				capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
			}
			git(t, root, "add", "--", "fifo-link")
			git(t, root, "commit", "-qm", "tracked symlink")
			return []string{"fifo-link"}
		}, "nothing to commit"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t)
			paths := tc.setup(t, root)
			base := git(t, root, "rev-parse", "HEAD")
			calls := 0
			o := New()
			o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
				calls++
				return authorization.Result{Kind: authorization.Green}
			}
			done := make(chan error, 1)
			go func() {
				_, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: tc.name, Paths: paths})
				done <- err
			}()
			select {
			case err := <-done:
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("Land error = %v, want %q", err, tc.wantErr)
				}
				if calls != 0 {
					t.Fatalf("authorization calls = %d, want 0", calls)
				}
			case <-time.After(bounds.TestDeadline(0)):
				t.Fatal("Land did not return within the bounded refusal deadline")
			}
		})
	}
}

func TestLandPreservesTrackedPathCompositionFailureBeforeAuthorization(t *testing.T) {
	root := fixture(t)
	base := git(t, root, "rev-parse", "HEAD")
	write(t, root, "named", "modified bytes that require a new object")

	objectDirectory := filepath.Join(t.TempDir(), "objects")
	if err := os.Mkdir(objectDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	blob := git(t, root, "hash-object", "--no-filters", "--", "named")
	if err := os.WriteFile(filepath.Join(objectDirectory, blob[:2]), []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_OBJECT_DIRECTORY", objectDirectory)
	t.Setenv("GIT_ALTERNATE_OBJECT_DIRECTORIES", filepath.Join(root, ".git", "objects"))

	calls := 0
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		calls++
		return authorization.Result{Kind: authorization.Green}
	}
	_, err := o.Land(context.Background(), Request{
		Root: root, Destination: "refs/heads/main", Expected: base,
		Message: "composition failure", Paths: []string{"named"},
	})
	if err == nil || !strings.Contains(err.Error(), `compose attributed path "named"`) {
		t.Fatalf("Land error = %v, want attributed composition error", err)
	}
	if strings.Contains(err.Error(), "nothing to commit") {
		t.Fatalf("Land error = %v, must preserve the composition failure", err)
	}
	if calls != 0 {
		t.Fatalf("authorization calls = %d, want 0", calls)
	}
}

func requireLandingFIFO(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
}

func requireLandingSocket(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("unix", path)
	if err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("unix sockets unavailable: %v", err))
	}
	t.Cleanup(func() { _ = listener.Close() })
}

func requireLandingDevice(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	const nullDevice = 1<<8 | 3
	if err := syscall.Mknod(path, syscall.S_IFCHR|0o600, nullDevice); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot create a character device: %v", err))
	}
}

type unnamedFingerprint struct {
	StagedIndexBlob             string
	StagedStatus                string
	StagedPlusUnstagedIndexBlob string
	StagedPlusUnstagedBytes     string
	StagedPlusUnstagedStatus    string
	UnstagedTrackedBytes        string
	UnstagedTrackedStatus       string
	UntrackedBytes              string
	UntrackedStatus             string
	IgnoredBytes                string
	IgnoredClassification       string
	WholeStatus                 string
}

func TestLandPreservesWholeUnnamedFingerprintAndCleansOwnedPaths(t *testing.T) {
	root := fixture(t)
	write(t, root, ".gitignore", "unnamed-ignored\n")
	write(t, root, "unnamed-staged", "base")
	write(t, root, "unnamed-staged-plus-unstaged", "base")
	write(t, root, "unnamed-unstaged", "base")
	git(t, root, "add", "--", ".gitignore", "unnamed-staged", "unnamed-staged-plus-unstaged", "unnamed-unstaged")
	git(t, root, "commit", "-qm", "fingerprint base")

	write(t, root, "named", "named landed bytes")
	write(t, root, "unnamed-staged", "staged bytes")
	git(t, root, "add", "--", "unnamed-staged")
	write(t, root, "unnamed-staged-plus-unstaged", "staged bytes")
	git(t, root, "add", "--", "unnamed-staged-plus-unstaged")
	write(t, root, "unnamed-staged-plus-unstaged", "worktree bytes")
	write(t, root, "unnamed-unstaged", "unstaged bytes")
	write(t, root, "unnamed-untracked", "untracked bytes")
	write(t, root, "unnamed-ignored", "ignored bytes")

	before := fingerprintUnnamedState(t, root)
	base := git(t, root, "rev-parse", "HEAD")
	o := greenOwner()
	got, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "fingerprint", Paths: []string{"named"}})
	if err != nil {
		t.Fatal(err)
	}
	after := fingerprintUnnamedState(t, root)
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("unnamed fingerprint changed\nbefore: %#v\nafter:  %#v", before, after)
	}
	if status := git(t, root, "status", "--porcelain=v1", "--", "named"); status != "" {
		t.Fatalf("owned paths are not clean: %q", status)
	}
	if published := git(t, root, "show", got.Commit+":named"); published != "named landed bytes" {
		t.Fatalf("published named bytes = %q", published)
	}
}

func TestFingerprintStatusRejectsMalformedRecord(t *testing.T) {
	if _, err := fingerprintStatus([]byte("?? path")); err == nil {
		t.Fatal("fingerprintStatus accepted unterminated porcelain record")
	}
}

func TestFingerprintStatusRejectsEmptyPrimaryPath(t *testing.T) {
	if _, err := fingerprintStatus([]byte("?? \x00")); err == nil {
		t.Fatal("fingerprintStatus accepted an empty primary path")
	}
}

func TestFingerprintStatusRejectsEmptyRenameSource(t *testing.T) {
	if _, err := fingerprintStatus([]byte("R  renamed\x00\x00")); err == nil {
		t.Fatal("fingerprintStatus accepted an empty rename source")
	}
}

func TestFingerprintStatusPreservesValidRecord(t *testing.T) {
	got, err := fingerprintStatus([]byte("?? file\x00"))
	if err != nil || string(got) != "?? file\x00" {
		t.Fatalf("fingerprintStatus valid record = %q, %v", got, err)
	}
}

func TestFingerprintStatusFiltersOnlyDeclaredRuntimeResidue(t *testing.T) {
	raw := []byte("!! .logs/build\x00!! cache/output\x00")
	got, err := fingerprintStatus(raw)
	if err != nil || string(got) != "!! cache/output\x00" {
		t.Fatalf("fingerprintStatus filtering = %q, %v", got, err)
	}
}

func fingerprintUnnamedState(t *testing.T, root string) unnamedFingerprint {
	t.Helper()
	return unnamedFingerprint{
		StagedIndexBlob:             git(t, root, "rev-parse", ":unnamed-staged"),
		StagedStatus:                git(t, root, "status", "--porcelain=v1", "--", "unnamed-staged"),
		StagedPlusUnstagedIndexBlob: git(t, root, "rev-parse", ":unnamed-staged-plus-unstaged"),
		StagedPlusUnstagedBytes:     string(mustRead(t, filepath.Join(root, "unnamed-staged-plus-unstaged"))),
		StagedPlusUnstagedStatus:    git(t, root, "status", "--porcelain=v1", "--", "unnamed-staged-plus-unstaged"),
		UnstagedTrackedBytes:        string(mustRead(t, filepath.Join(root, "unnamed-unstaged"))),
		UnstagedTrackedStatus:       git(t, root, "status", "--porcelain=v1", "--", "unnamed-unstaged"),
		UntrackedBytes:              string(mustRead(t, filepath.Join(root, "unnamed-untracked"))),
		UntrackedStatus:             git(t, root, "status", "--porcelain=v1", "--", "unnamed-untracked"),
		IgnoredBytes:                string(mustRead(t, filepath.Join(root, "unnamed-ignored"))),
		IgnoredClassification:       git(t, root, "check-ignore", "-v", "--", "unnamed-ignored"),
		WholeStatus: git(t, root, "status", "--porcelain=v1", "--untracked-files=all", "--ignored", "--",
			"unnamed-staged", "unnamed-staged-plus-unstaged", "unnamed-unstaged", "unnamed-untracked", "unnamed-ignored"),
	}
}

func TestLandReviewedRechecksCheckoutFingerprintsAfterAuthorization(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		mutate     func(*testing.T, string, string)
	}{
		{name: "destination", want: "destination checkout changed", mutate: func(t *testing.T, root, _ string) { write(t, root, "late", "destination dirt") }},
		{name: "source", want: "source checkout changed", mutate: func(t *testing.T, _, source string) { write(t, source, "late", "source dirt") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t)
			write(t, root, "specs/x/spec.md", "Status: staged\n")
			git(t, root, "add", "specs/x/spec.md")
			git(t, root, "commit", "-qm", "stage spec")
			destination := git(t, root, "rev-parse", "HEAD")
			sourceWorktree := filepath.Join(t.TempDir(), "source")
			git(t, root, "worktree", "add", "-qb", "reviewed-source", sourceWorktree, destination)
			write(t, sourceWorktree, "reviewed", "source bytes\n")
			git(t, sourceWorktree, "add", "reviewed")
			git(t, sourceWorktree, "commit", "-qm", "reviewed")
			source := git(t, sourceWorktree, "rev-parse", "HEAD")
			sourceFingerprint, _ := CheckoutFingerprint(sourceWorktree)
			destinationFingerprint, _ := CheckoutFingerprint(root)
			o := New()
			o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
				tc.mutate(t, root, sourceWorktree)
				return authorization.Result{Kind: authorization.Green}
			}
			_, err := o.LandReviewed(context.Background(), ReviewedRequest{
				Root: root, Destination: "refs/heads/main", DestinationBase: destination,
				Source: "refs/heads/reviewed-source", SourceTip: source, ReviewBase: destination,
				SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
				SpecPath: "specs/x/spec.md", SpecBytes: []byte("Status: staged\n"), SpecMode: 0o644, Message: "must refuse",
			})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("fingerprint mutation error = %v, want %q", err, tc.want)
			}
			if got := git(t, root, "rev-parse", "main"); got != destination {
				t.Fatalf("destination moved to %s", got)
			}
		})
	}
}
