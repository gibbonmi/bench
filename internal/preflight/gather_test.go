package preflight

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// ticketDoc renders a grammar-conformant ticket citing rows. `Writes:` names
// the specs folder every seeded fixture already holds, so writes-resolve
// stays green without a planted file.
func ticketDoc(title string, covers ...string) string {
	return "# " + title + "\n\n" +
		"Blocked by: none\n" +
		"Writes: specs\n" +
		"Covers: " + strings.Join(covers, ", ") + "\n\n" +
		"## What to build\n\nBuild it.\n\n" +
		"## Acceptance\n\n- [ ] It is built.\n"
}

// gatherTicketsAt parses one tickets directory under a root the test owns,
// with the PF tag every fixture spec declares.
func gatherTicketsAt(t *testing.T, root, dir string) (ticketFacts, *BootstrapFailure) {
	t.Helper()
	return gatherTickets(root, dir, "review", "PF")
}

// TestTicketParseTrailingNewlineParity is PF17's ticket half at the gatherer
// seam. A ticket whose last line lacks its newline parses to the same ticket
// as its terminated form, so a hand-edited ending never invents a fault.
func TestTicketParseTrailingNewlineParity(t *testing.T) {
	terminated := ticketDoc("One", "PF1")
	unterminated := strings.TrimSuffix(terminated, "\n")
	if unterminated == terminated {
		t.Fatal("fixture invalid: terminated fixture already lacks a trailing newline")
	}

	termRoot := t.TempDir()
	termDir := filepath.Join(termRoot, "tickets")
	mustWriteFile(t, filepath.Join(termDir, "one.md"), terminated)
	want, wantErr := gatherTicketsAt(t, termRoot, termDir)
	if wantErr != nil {
		t.Fatalf("gatherTickets(terminated) error = %+v, want nil", wantErr)
	}

	unterminatedRoot := t.TempDir()
	unterminatedDir := filepath.Join(unterminatedRoot, "tickets")
	mustWriteFile(t, filepath.Join(unterminatedDir, "one.md"), unterminated)
	got, gotErr := gatherTicketsAt(t, unterminatedRoot, unterminatedDir)
	if gotErr != nil {
		t.Fatalf("gatherTickets(unterminated) error = %+v, want nil", gotErr)
	}

	if !reflect.DeepEqual(got.parsed, want.parsed) {
		t.Errorf("gatherTickets(unterminated) parsed %#v, want %#v (same as terminated)", got.parsed, want.parsed)
	}
	if len(want.diagnostics) > 0 {
		t.Fatalf("fixture invalid: the terminated form reported %#v, want no diagnostic", want.diagnostics)
	}
	if !reflect.DeepEqual(want.parsed[0].Covers, []string{"PF1"}) {
		t.Fatalf("fixture invalid: the terminated form covered %#v, want [\"PF1\"]", want.parsed[0].Covers)
	}
}

// TestNonMarkdownFileIsIgnored is TG37. A stray asset under tickets/ is not a
// half-written ticket, so the grammar passes over it rather than reporting
// every required field as absent.
func TestNonMarkdownFileIsIgnored(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "tickets")
	mustWriteFile(t, filepath.Join(dir, "one.md"), ticketDoc("One", "PF1"))
	mustWriteFile(t, filepath.Join(dir, "diagram.png"), "\xef\xbb\xbfnot a ticket at all\n")
	mustWriteFile(t, filepath.Join(dir, "NOTES"), "scratch\n")

	facts, bootErr := gatherTicketsAt(t, root, dir)
	if bootErr != nil {
		t.Fatalf("gatherTickets = %+v, want facts", bootErr)
	}
	if len(facts.diagnostics) > 0 {
		t.Errorf("diagnostics = %#v, want none (a non-.md file is ignored)", facts.diagnostics)
	}
	if len(facts.parsed) != 1 || facts.parsed[0].Name != "one.md" {
		t.Errorf("parsed = %#v, want the one .md ticket", facts.parsed)
	}
}

// TestDuplicateBasenameAcrossDepths is TG38. A blocker edge resolves by
// basename, so one basename at two depths is an ambiguous identity rather
// than two tickets.
func TestDuplicateBasenameAcrossDepths(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "tickets")
	mustWriteFile(t, filepath.Join(dir, "one.md"), ticketDoc("One", "PF1"))
	mustWriteFile(t, filepath.Join(dir, "sub", "one.md"), ticketDoc("One again", "PF2"))

	facts, bootErr := gatherTicketsAt(t, root, dir)
	if bootErr != nil {
		t.Fatalf("gatherTickets = %+v, want facts", bootErr)
	}
	want := "duplicate ticket basename one.md at one.md and sub/one.md"
	if len(facts.diagnostics) != 1 || facts.diagnostics[0] != want {
		t.Fatalf("diagnostics = %#v, want exactly [%q]", facts.diagnostics, want)
	}
	if len(facts.parsed) != 1 {
		t.Errorf("parsed = %#v, want the first copy alone so the sibling set holds one name per ticket", facts.parsed)
	}
}

// TestSpecialTicketEntryRefused is TG6. A special file, a dangling symlink,
// and an unreadable entry are refused by name at every depth. The extension
// decides what parses, never whether a broken entry may read as green.
func TestSpecialTicketEntryRefused(t *testing.T) {
	cases := []struct {
		name  string
		plant func(t *testing.T, dir string)
		named string
	}{
		{"a FIFO at the top level", func(t *testing.T, dir string) {
			if err := syscall.Mkfifo(filepath.Join(dir, "pipe.md"), 0o600); err != nil {
				t.Fatalf("Mkfifo: %v", err)
			}
		}, "pipe.md"},
		{"a FIFO one level down", func(t *testing.T, dir string) {
			sub := filepath.Join(dir, "sub")
			if err := os.MkdirAll(sub, 0o755); err != nil {
				t.Fatalf("MkdirAll: %v", err)
			}
			if err := syscall.Mkfifo(filepath.Join(sub, "pipe"), 0o600); err != nil {
				t.Fatalf("Mkfifo: %v", err)
			}
		}, "pipe"},
		{"a dangling symlink", func(t *testing.T, dir string) {
			if err := os.Symlink(filepath.Join(dir, "gone.md"), filepath.Join(dir, "dangling.md")); err != nil {
				t.Fatalf("Symlink: %v", err)
			}
		}, "dangling.md"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "tickets")
			mustWriteFile(t, filepath.Join(dir, "one.md"), ticketDoc("One", "PF1"))
			tc.plant(t, dir)

			facts, bootErr := gatherTicketsAt(t, root, dir)
			if bootErr == nil {
				t.Fatalf("gatherTickets = (%+v, nil), want a refusal naming %s", facts, tc.named)
			}
			if !strings.Contains(bootErr.Hint, tc.named) {
				t.Errorf("bootErr = %+v, want the hint to name %s", bootErr, tc.named)
			}
		})
	}
}

// TestGatherSpecStatusOutsideFolderEnumerationNotReadable reaches
// gather.go's spec-status fallback branch. specref.Resolve succeeds (the
// argument is a literal path, tried as-given) over a file that
// specref.Facts' folder-spec enumeration never covers. It is not
// literally named spec.md.
//
// So no typed status is available to trust. Gather answers the "spec
// status not readable" BootstrapFailure naming the enumeration mismatch.
// This is distinct from the specref.Facts-error half of the same Kind.
func TestGatherSpecStatusOutsideFolderEnumerationNotReadable(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	slug := "specs/example/other.md" // valid content, but not specs/<slug>/spec.md
	mustWriteFile(t, slug, specBody("example"))

	facts, bootErr := Gather(root, "review", slug)
	if bootErr == nil {
		t.Fatalf("Gather = (%+v, nil), want a spec-status-not-readable BootstrapFailure", facts)
	}
	if bootErr.Kind != "spec status not readable" {
		t.Fatalf("bootErr.Kind = %q, want %q", bootErr.Kind, "spec status not readable")
	}
	if !strings.Contains(bootErr.Hint, "did not resolve through folder-spec enumeration") {
		t.Errorf("bootErr.Hint = %q, want it to name the folder-spec enumeration mismatch", bootErr.Hint)
	}
}

// activeAssignment registers one active assignment in root's ledger, owning the
// tree at worktree, and returns its id. The record goes in through
// intent.PutAssignment, so the gatherer reads exactly the shape every worktree
// command writes.
func activeAssignment(t *testing.T, root, worktree string) string {
	t.Helper()
	const id = "00000000000000000000000000000001"
	const owner = "00000000000000000000000000000002"
	err := intent.PutAssignment(root, intent.Assignment{
		Schema:   intent.AssignmentRecordSchema,
		ID:       id,
		OwnerID:  owner,
		Request:  intent.RequestDigest("preflight-assignment-target"),
		Label:    "preflight-assignment-target",
		Start:    runGit(t, "rev-parse", "HEAD"),
		Branch:   intent.AssignmentBranchRef(owner, id),
		Worktree: worktree,
		State:    intent.StateActive,
	})
	if err != nil {
		t.Fatalf("PutAssignment: %v", err)
	}
	return id
}

// TestGatherAssignmentTarget covers WF38, WF39, and WF40. The remedy the
// stale-base red prints addresses the assignment the operator is standing in,
// so the gatherer must recognize that tree by identity and no other.
func TestGatherAssignmentTarget(t *testing.T) {
	t.Run("the assignment's own worktree", func(t *testing.T) {
		// WF38: the ledger holds this tree, so its id fills the fact.
		root, slug := seedConformant(t)
		canonical, err := canonicalRoot(root)
		if err != nil {
			t.Fatalf("canonicalRoot(%q): %v", root, err)
		}
		id := activeAssignment(t, root, canonical)

		facts, bootErr := Gather(root, "review", slug)
		if bootErr != nil {
			t.Fatalf("Gather = %+v, want facts", bootErr)
		}
		if facts.AssignmentTarget != id {
			t.Errorf("AssignmentTarget = %q, want %q", facts.AssignmentTarget, id)
		}
	})

	t.Run("a symlink to that worktree", func(t *testing.T) {
		// WF39: the ledger records a resolved path, so a raw string compare
		// against a symlinked root would miss the assignment that owns it.
		root, slug := seedConformant(t)
		canonical, err := canonicalRoot(root)
		if err != nil {
			t.Fatalf("canonicalRoot(%q): %v", root, err)
		}
		id := activeAssignment(t, root, canonical)

		link := filepath.Join(t.TempDir(), "link")
		if err := os.Symlink(canonical, link); err != nil {
			t.Fatalf("Symlink(%q, %q): %v", canonical, link, err)
		}
		if link == canonical {
			t.Fatal("fixture invalid: the symlink names the same path as the tree")
		}

		facts, bootErr := Gather(link, "review", slug)
		if bootErr != nil {
			t.Fatalf("Gather = %+v, want facts", bootErr)
		}
		if facts.AssignmentTarget != id {
			t.Errorf("AssignmentTarget from the symlink = %q, want %q", facts.AssignmentTarget, id)
		}
	})

	t.Run("the primary checkout", func(t *testing.T) {
		// WF40: an active assignment owning some other tree names a stranger, so
		// the primary checkout's remedy keeps the placeholder.
		root, slug := seedConformant(t)
		other, err := canonicalRoot(t.TempDir())
		if err != nil {
			t.Fatalf("canonicalRoot: %v", err)
		}
		activeAssignment(t, root, other)

		facts, bootErr := Gather(root, "review", slug)
		if bootErr != nil {
			t.Fatalf("Gather = %+v, want facts", bootErr)
		}
		if facts.AssignmentTarget != "" {
			t.Errorf("AssignmentTarget in the primary checkout = %q, want empty", facts.AssignmentTarget)
		}
	})
}
