package git

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

// newLinkedCheckout returns a fresh repository and a linked worktree of it.
func newLinkedCheckout(t *testing.T) (root, linked string) {
	t.Helper()
	root = newRepo(t)
	linked = filepath.Join(t.TempDir(), "linked")
	runGit(t, root, "worktree", "add", "-q", "-b", "linked", linked)
	return root, linked
}

// mergeHeadPath answers the MERGE_HEAD path that Git resolves for checkout.
func mergeHeadPath(t *testing.T, checkout string) string {
	t.Helper()
	return strings.TrimSpace(runGit(t, checkout, "rev-parse", "--path-format=absolute", "--git-path", "MERGE_HEAD"))
}

// pendingMergeWithin runs PendingMerge under the test deadline, so a reader that blocks
// fails the test instead of stalling it.
func pendingMergeWithin(t *testing.T, root string) ([]string, error) {
	t.Helper()
	type result struct {
		heads []string
		err   error
	}
	done := make(chan result, 1)
	go func() {
		heads, err := PendingMerge(root)
		done <- result{heads: heads, err: err}
	}()
	select {
	case result := <-done:
		return result.heads, result.err
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal(bounds.TestTimeoutVerdict("git.PendingMerge", bounds.TestDeadline(0)))
		return nil, nil
	}
}

func requirePendingMergeRefusal(t *testing.T, root string) {
	t.Helper()
	heads, err := pendingMergeWithin(t, root)
	if err == nil || !strings.HasPrefix(err.Error(), "read pending merge:") || len(heads) != 0 {
		t.Fatalf("PendingMerge = %q, %v; want no merge head and a read pending merge refusal", heads, err)
	}
}

func TestIsPrimaryCheckoutSeparatesLinkedWorktree(t *testing.T) {
	root, linked := newLinkedCheckout(t)

	primary, err := IsPrimaryCheckout(root)
	if err != nil || !primary {
		t.Fatalf("IsPrimaryCheckout(primary) = %t, %v; want true, nil", primary, err)
	}
	primary, err = IsPrimaryCheckout(linked)
	if err != nil || primary {
		t.Fatalf("IsPrimaryCheckout(linked) = %t, %v; want false, nil", primary, err)
	}
}

// TestPendingMergeReadsEveryHead proves the reader answers the merge heads that Git
// writes in a linked worktree in file order, and no merge head without a MERGE_HEAD.
// (Coverage rows MP01-MP04.)
func TestPendingMergeReadsEveryHead(t *testing.T) {
	root, linked := newLinkedCheckout(t)
	base := strings.TrimSpace(runGit(t, linked, "rev-parse", "HEAD"))
	var heads []string
	for _, name := range []string{"one", "two"} {
		runGit(t, root, "checkout", "-q", "-b", name, base)
		if err := os.WriteFile(filepath.Join(root, name+".txt"), []byte(name+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		runGit(t, root, "add", name+".txt")
		runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", name)
		heads = append(heads, strings.TrimSpace(runGit(t, root, "rev-parse", "HEAD")))
	}
	merge := func(t *testing.T, commits ...string) {
		t.Helper()
		runGit(t, linked, append([]string{"-c", "user.email=bench@local", "-c", "user.name=bench", "merge", "-q", "--no-commit", "--no-ff"}, commits...)...)
	}
	// Each step starts from the merge state that the step before it leaves.
	for _, step := range []struct {
		name    string
		prepare func(t *testing.T)
		want    []string
	}{
		{"absent", func(*testing.T) {}, nil},
		{"merge", func(t *testing.T) { merge(t, heads[0]) }, heads[:1]},
		{"octopus", func(t *testing.T) { runGit(t, linked, "merge", "--abort"); merge(t, heads...) }, heads},
		{"unterminated last line", func(t *testing.T) {
			runGit(t, linked, "merge", "--abort")
			if err := os.WriteFile(mergeHeadPath(t, linked), []byte(heads[0]+"\n"+heads[1]), 0o644); err != nil {
				t.Fatal(err)
			}
		}, heads},
	} {
		t.Run(step.name, func(t *testing.T) {
			step.prepare(t)
			got, err := PendingMerge(linked)
			if err != nil || !slices.Equal(got, step.want) {
				t.Fatalf("PendingMerge = %q, %v; want %q, nil", got, err, step.want)
			}
		})
	}
}

// TestPendingMergeRefusesUnreadableState proves the reader refuses a root outside every
// repository and each MERGE_HEAD that is not a regular file, a link to a valid merge
// record included. (Coverage rows MP05-MP09.)
func TestPendingMergeRefusesUnreadableState(t *testing.T) {
	_, linked := newLinkedCheckout(t)
	mergeHead := mergeHeadPath(t, linked)
	record := filepath.Join(t.TempDir(), "record")
	if err := os.WriteFile(record, []byte(strings.TrimSpace(runGit(t, linked, "rev-parse", "HEAD"))+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		root  string
		plant func(t *testing.T) error
	}{
		{"outside every repository", t.TempDir(), func(*testing.T) error { return nil }},
		{"directory", linked, func(*testing.T) error { return os.Mkdir(mergeHead, 0o755) }},
		{"fifo", linked, func(t *testing.T) error {
			if err := syscall.Mkfifo(mergeHead, 0o600); err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
			}
			return nil
		}},
		{"dangling symbolic link", linked, func(t *testing.T) error {
			return os.Symlink(filepath.Join(t.TempDir(), "missing"), mergeHead)
		}},
		{"live symbolic link", linked, func(*testing.T) error { return os.Symlink(record, mergeHead) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() { _ = os.Remove(mergeHead) })
			if err := tc.plant(t); err != nil {
				t.Fatal(err)
			}
			requirePendingMergeRefusal(t, tc.root)
		})
	}
}

// TestPendingMergeRefusesMalformedState proves the reader refuses each MERGE_HEAD whose
// content is not a list of full commit IDs. (Coverage rows MP10-MP14.)
func TestPendingMergeRefusesMalformedState(t *testing.T) {
	root, linked := newLinkedCheckout(t)
	mergeHead := mergeHeadPath(t, linked)
	head := strings.TrimSpace(runGit(t, linked, "rev-parse", "HEAD"))
	runGit(t, root, "update-ref", "refs/heads/main", head)
	for _, tc := range []struct{ name, content string }{
		{"empty", ""},
		{"abbreviated commit ID", head[:7] + "\n"},
		{"branch name", "main\n"},
		{"unknown object", strings.Repeat("e", len(head)) + "\n"},
		{"tree", strings.TrimSpace(runGit(t, linked, "rev-parse", "HEAD^{tree}")) + "\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(mergeHead, []byte(tc.content), 0o644); err != nil {
				t.Fatal(err)
			}
			requirePendingMergeRefusal(t, linked)
		})
	}
}
