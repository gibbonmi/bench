package commitment_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
)

// Plan refuses a changed pinned row whose working path is hostile, even when the proposal
// names the identity of the bytes that the path holds or links to, and the refusal writes
// nothing. Each row plants one hostile shape at the working path of FT1.
func TestPinnedRowPlanRefusesAHostileWorkingFile(t *testing.T) {
	oversized := strings.Repeat("a", int(bounds.ControlRecordLimit)+1)
	for _, row := range []struct {
		name, body string
		plant      func(t *testing.T, worktree, path string)
	}{
		{name: "absent", body: pinnedEdit, plant: func(t *testing.T, worktree, path string) {
			if err := os.Remove(filepath.Join(worktree, filepath.FromSlash(path))); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "empty", body: "", plant: func(t *testing.T, worktree, path string) {
			commitmenttest.Write(t, worktree, path, "")
		}},
		{name: "invalid-utf8", body: "\xff\xfe\n", plant: func(t *testing.T, worktree, path string) {
			commitmenttest.Write(t, worktree, path, "\xff\xfe\n")
		}},
		{name: "oversized", body: oversized, plant: func(t *testing.T, worktree, path string) {
			commitmenttest.Write(t, worktree, path, oversized)
		}},
		{name: "live-link", body: pinnedEdit, plant: func(t *testing.T, worktree, path string) {
			target := filepath.Join(t.TempDir(), "FT1.md")
			commitmenttest.Write(t, filepath.Dir(target), filepath.Base(target), pinnedEdit)
			replaceWithLink(t, filepath.Join(worktree, filepath.FromSlash(path)), target)
		}},
		{name: "dangling-link", body: pinnedEdit, plant: func(t *testing.T, worktree, path string) {
			replaceWithLink(t, filepath.Join(worktree, filepath.FromSlash(path)), filepath.Join(t.TempDir(), "gone.md"))
		}},
		{name: "fifo", body: pinnedEdit, plant: func(t *testing.T, worktree, path string) {
			full := filepath.Join(worktree, filepath.FromSlash(path))
			if err := os.Remove(full); err != nil {
				t.Fatal(err)
			}
			if err := syscall.Mkfifo(full, 0o600); err != nil {
				capability.Unavailable(t, capability.Fifo, err.Error())
			}
		}},
		{name: "linked-roadmap", body: pinnedEdit, plant: func(t *testing.T, worktree, path string) {
			outside := t.TempDir()
			commitmenttest.Write(t, outside, filepath.Base(path), pinnedEdit)
			replaceWithLink(t, filepath.Join(worktree, filepath.FromSlash(filepath.Dir(path))), outside)
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, worktree, current := commitmenttest.PinnedRows(t)
			source := current.Milestones[0].Outcomes[0].Sources[0]
			row.plant(t, worktree, source.Path)
			before := commitmenttest.MilestoneState(t, root, worktree)
			_, err := (commitrepo.Store{Root: worktree}).Plan(mustBytes(t, repin(current, 0, row.body)))
			want := "commitment source " + strconv.Quote(source.ID) + " refused: "
			if err == nil || !strings.HasPrefix(err.Error(), want) || commitmenttest.MilestoneState(t, root, worktree) != before {
				t.Fatalf("Plan(a %s working row) = %v, want a refusal naming %q that writes nothing", row.name, err, want)
			}
		})
	}
}

// replaceWithLink removes path and plants a symbolic link to target in its place.
func replaceWithLink(t *testing.T, path, target string) {
	t.Helper()
	if err := os.RemoveAll(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
}
