package preflight

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// TestTicketTokenScanTrailingNewlineParity is PF17's ticket-token-scan
// half, made real at the gatherer seam itself. A ticket file whose only
// citation sits on an unterminated last line scans to the same token
// slice as its terminated form.
func TestTicketTokenScanTrailingNewlineParity(t *testing.T) {
	terminated := "Ticket citing PF1.\n"
	unterminated := strings.TrimSuffix(terminated, "\n")
	if unterminated == terminated {
		t.Fatal("fixture invalid: terminated fixture already lacks a trailing newline")
	}

	termDir := t.TempDir()
	mustWriteFile(t, filepath.Join(termDir, "one.md"), terminated)
	wantTokens, _, wantErr := gatherTicketTokens(termDir, "review")
	if wantErr != nil {
		t.Fatalf("gatherTicketTokens(terminated) error = %+v, want nil", wantErr)
	}

	unterminatedDir := t.TempDir()
	mustWriteFile(t, filepath.Join(unterminatedDir, "one.md"), unterminated)
	gotTokens, _, gotErr := gatherTicketTokens(unterminatedDir, "review")
	if gotErr != nil {
		t.Fatalf("gatherTicketTokens(unterminated) error = %+v, want nil", gotErr)
	}

	if !reflect.DeepEqual(gotTokens, wantTokens) {
		t.Errorf("gatherTicketTokens(unterminated) = %#v, want %#v (same as terminated)", gotTokens, wantTokens)
	}
	if !reflect.DeepEqual(wantTokens, []string{"PF1"}) {
		t.Fatalf("fixture invalid: terminated form scanned to %#v, want [\"PF1\"]", wantTokens)
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
