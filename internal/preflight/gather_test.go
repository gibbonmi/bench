package preflight

import (
	"bytes"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestFenceTokensTrailingNewlineParity is PF17's fenceTokens half, made
// real at the gatherer seam itself. It replaces the tautological
// Decide(x)==Decide(x) it used to stand in for. A `## Ownership fences`
// section whose final line lacks a trailing newline parses to the same
// token slice as its terminated form.
//
// fenceTokens splits on "\n" directly, so the unterminated last element
// from strings.Split still carries its content. A scanner keyed on a
// trailing "\n" per token instead would drop it.
func TestFenceTokensTrailingNewlineParity(t *testing.T) {
	terminated := []byte("## Ownership fences\n\n- `internal/example/`\n")
	unterminated := bytes.TrimSuffix(terminated, []byte("\n"))
	if bytes.Equal(terminated, unterminated) {
		t.Fatal("fixture invalid: terminated fixture already lacks a trailing newline")
	}

	got := fenceTokens(unterminated)
	want := fenceTokens(terminated)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("fenceTokens(unterminated) = %#v, want %#v (same as terminated)", got, want)
	}
	if !reflect.DeepEqual(want, []string{"internal/example/"}) {
		t.Fatalf("fixture invalid: terminated form parsed to %#v, want [\"internal/example/\"]", want)
	}
}

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
