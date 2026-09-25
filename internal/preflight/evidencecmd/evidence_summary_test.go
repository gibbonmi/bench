package evidencecmd_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// readSummary runs the default read of identity and returns its one summary row. The
// response must hold exactly one block.
func readSummary(t *testing.T, identity string) map[string]any {
	t.Helper()
	out, code := preflight.Command([]string{"evidence", identity})
	if code != 0 {
		t.Fatalf("default read = (%d):\n%s", code, out)
	}
	document := preflighttest.DecodeMap(t, out)
	if len(document) != 1 {
		t.Fatalf("default read holds %d blocks, want one:\n%s", len(document), out)
	}
	rows := preflighttest.TableRows(t, document, "evidence_summary")
	if len(rows) != 1 {
		t.Fatalf("evidence_summary rows = %d:\n%s", len(rows), out)
	}
	return rows[0]
}

// commandArgs returns the preflight arguments of one printed next command.
func commandArgs(t *testing.T, next string) []string {
	t.Helper()
	fields := strings.Fields(next)
	if len(fields) < 4 || strings.Join(fields[:2], " ") != "bench preflight" {
		t.Fatalf("next = %q is not a preflight command", next)
	}
	return fields[2:]
}

// TestEvidenceDefaultPrintsSummary is BO42: the default read orients the consumer and
// carries no content bytes.
func TestEvidenceDefaultPrintsSummary(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, prepared, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	row := readSummary(t, identity)
	assertTypedRow(t, row, map[string]string{"evidence": "string", "sources": "integer", "pages": "integer",
		"manifest_bytes": "integer", "source_bytes": "integer", "next": "string"})
	manifest, _ := reconstructEvidence(t, identity, traverseEvidence(t, identity))
	sourceBytes := 0.0
	for _, source := range manifestSources(t, manifest) {
		sourceBytes += source["bytes"].(float64)
	}
	if row["evidence"] != identity || row["sources"] != prepared["sources"] || row["pages"] != prepared["pages"] ||
		row["manifest_bytes"] != prepared["manifest_bytes"] || row["source_bytes"] != sourceBytes {
		t.Fatalf("summary = %v, want the prepared counts %v and %v source bytes", row, prepared, sourceBytes)
	}
}

// TestEvidenceSummaryNextReadsFirstPage is BO43: the summary's successor is the first
// manifest fragment, byte for byte the response of the manifest-start cursor.
func TestEvidenceSummaryNextReadsFirstPage(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	got, code := preflight.Command(commandArgs(t, readSummary(t, identity)["next"].(string)))
	if code != 0 {
		t.Fatalf("summary next = (%d):\n%s", code, got)
	}
	want, code := preflight.Command([]string{"evidence", identity, "--cursor", "v1." + strings.TrimPrefix(identity, "sha256:") + ".m.0.0"})
	if code != 0 || got != want {
		t.Fatalf("summary next read:\n%s\nwant the manifest-start page (%d):\n%s", got, code, want)
	}
}
