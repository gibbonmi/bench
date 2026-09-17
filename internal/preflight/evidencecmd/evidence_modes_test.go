package evidencecmd_test

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// This file covers three read modes: the selected source stream, full verification, and the
// current-action binding. The expected headers, cell types, and refusal texts below are
// stated independently of the format and operation registries.

// checkCurrent runs the current-action check over one prepared artifact.
func checkCurrent(t *testing.T, identity string) (string, int) {
	t.Helper()
	return preflight.Command([]string{"evidence", identity, "--check-current"})
}

// craftedDigestArtifact republishes one prepared artifact with a changed source digest, so
// the Git pair still matches while the prepared bytes no longer describe the tree. It
// reads the header length field and the manifest with literal offsets, as an independent
// consumer does.
func craftedDigestArtifact(t *testing.T, root, identity, oldDigest, newDigest string) string {
	t.Helper()
	dir := preflighttest.StoreDir(t, root)
	data, err := os.ReadFile(filepath.Join(dir, strings.TrimPrefix(identity, "sha256:")+".pack"))
	if err != nil {
		t.Fatal(err)
	}
	length := int(binary.LittleEndian.Uint64(data[16:24]))
	manifest := string(data[24 : 24+length])
	if strings.Count(manifest, oldDigest) == 0 {
		t.Fatalf("manifest holds no digest %s", oldDigest)
	}
	crafted := strings.Replace(manifest, oldDigest, newDigest, 1)
	copyOf := append(append(append([]byte{}, data[:24]...), crafted...), data[24+length:]...)
	craftedIdentity := "sha256:" + sha(crafted)
	if err := os.WriteFile(filepath.Join(dir, sha(crafted)+".pack"), copyOf, 0o600); err != nil {
		t.Fatal(err)
	}
	return craftedIdentity
}

// TestEvidenceCurrentBinding is CE57, CE58, CE59, CE123, and CE124.
func TestEvidenceCurrentBinding(t *testing.T) {
	t.Run("CE59 current binding for unchanged pins", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		identity, prepared, _ := prepareEvidence(t, args)
		out, code := checkCurrent(t, identity)
		if code != 0 {
			t.Fatalf("check-current = (%d):\n%s", code, out)
		}
		if got := headerLine(out); got != "current[1]{evidence,assignment,base,source_tip,current,delivery}:" {
			t.Fatalf("current header = %q", got)
		}
		row := preflighttest.TableRows(t, preflighttest.DecodeMap(t, out), "current")[0]
		assertTypedRow(t, row, map[string]string{"evidence": "string", "assignment": "string", "base": "string",
			"source_tip": "string", "current": "boolean", "delivery": "string"})
		if row["evidence"] != identity || row["assignment"] != preflighttest.ChargeFixtureAssignment ||
			row["base"] != prepared["base"] || row["source_tip"] != prepared["source_tip"] ||
			row["current"] != true || row["delivery"] != "unverified" {
			t.Fatalf("current row = %v, want the current binding of %v", row, prepared)
		}
	})
	t.Run("CE57 moved source tip", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
		preflighttest.MustWriteFile(t, "internal/"+slug+"/moved.go", "package example\n")
		preflighttest.RunGit(t, "add", "-A")
		preflighttest.RunGit(t, "commit", "-q", "-m", "move the source tip")
		out, code := checkCurrent(t, identity)
		if code != 1 || !strings.Contains(out, "is not the prepared source tip") || strings.Contains(out, "current[1]") {
			t.Fatalf("moved tip = (%d):\n%s", code, out)
		}
	})
	t.Run("CE58 released assignment", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
		preflighttest.OwnedAssignment(t, root, root, intent.StateComplete)
		out, code := checkCurrent(t, identity)
		if code != 1 || !strings.Contains(out, "assignment required") || strings.Contains(out, "current[1]") {
			t.Fatalf("released assignment = (%d):\n%s", code, out)
		}
	})
	t.Run("CE123 dirty assignment checkout", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
		preflighttest.MustWriteFile(t, "internal/"+slug+"/dirty.go", "package example\n")
		out, code := checkCurrent(t, identity)
		if code != 1 || !strings.Contains(out, "source checkout is dirty") || strings.Contains(out, "current[1]") {
			t.Fatalf("dirty checkout = (%d):\n%s", code, out)
		}
	})
	t.Run("CE124 changed required source bytes", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
		// The crafted artifact keeps the Git pair and every page digest; only the ticket
		// source descriptor claims bytes the tree does not hold.
		ticket := preflighttest.TicketDoc("One", "PF1", "PF2")
		crafted := craftedDigestArtifact(t, root, identity, sha(ticket), strings.Repeat("c", 64))
		if out, code := checkCurrent(t, identity); code != 0 {
			t.Fatalf("prepared artifact = (%d):\n%s", code, out)
		}
		out, code := checkCurrent(t, crafted)
		if code != 1 || !strings.Contains(out, "differs from the prepared evidence") || strings.Contains(out, "current[1]") {
			t.Fatalf("changed required source = (%d):\n%s", code, out)
		}
	})
}

// traverseSource follows one declared source stream from its first page to its end.
func traverseSource(t *testing.T, identity, source string) []evidencePage {
	t.Helper()
	return traverseFrom(t, []string{"evidence", identity, "--source", source})
}

// TestEvidenceSourceNavigation is CE62: an explicit selector returns only that declared
// source stream, ends after it, and never continues into the default traversal. The
// single-page source proves the stream ends at a first page that is also its last.
func TestEvidenceSourceNavigation(t *testing.T) {
	t.Run("paged source", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		ticket := preflighttest.TicketDoc("One", "PF1", "PF2") + strings.Repeat("paged ticket\n", 1000)
		preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticket)
		identity, _, _ := prepareEvidence(t, preflighttest.LegacyCommitted(t, root, slug, "paged source"))
		pages := traverseSource(t, identity, "s2")
		body := ""
		for i, page := range pages {
			if page.row["stream"] != "source" || page.row["source"] != "s2" {
				t.Fatalf("page %d left the selected source: %v", i, page.row)
			}
			if got := page.row["next"].(string); (got == "") != (i == len(pages)-1) {
				t.Fatalf("page %d next = %q", i, got)
			}
			body += page.row["content"].(string)
		}
		if len(pages) < 2 || body != ticket || pages[len(pages)-1].row["stream_end"] != true {
			t.Fatalf("source stream held %d pages and %d of %d bytes", len(pages), len(body), len(ticket))
		}
		// Every page of the selected source appears exactly once, and no other source does.
		if all := traverseEvidence(t, identity); len(all) <= len(pages) {
			t.Fatalf("default traversal held %d pages, want more than the %d of one source", len(all), len(pages))
		}
	})
	t.Run("single-page source", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		ticket := preflighttest.TicketDoc("One", "PF1", "PF2")
		identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
		pages := traverseSource(t, identity, "s2")
		if len(pages) != 1 {
			t.Fatalf("single-page source held %d pages", len(pages))
		}
		row := pages[0].row
		if row["next"] != "" || row["stream_end"] != true || row["index"] != float64(0) ||
			row["source"] != "s2" || row["content"] != ticket {
			t.Fatalf("single page = %v, want the whole %d-byte ticket and the stream end", row, len(ticket))
		}
	})
}

// TestEvidenceVerifyCommand is CE160 at the command seam: verification reports every page
// and source it checked and still claims no delivery.
func TestEvidenceVerifyCommand(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, prepared, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	out, code := preflight.Command([]string{"evidence", identity, "--verify"})
	if code != 0 {
		t.Fatalf("verify = (%d):\n%s", code, out)
	}
	if got := headerLine(out); got != "verified[1]{evidence,manifest_verified,pages_verified,sources_verified,delivery}:" {
		t.Fatalf("verified header = %q", got)
	}
	row := preflighttest.TableRows(t, preflighttest.DecodeMap(t, out), "verified")[0]
	assertTypedRow(t, row, map[string]string{"evidence": "string", "manifest_verified": "boolean",
		"pages_verified": "integer", "sources_verified": "integer", "delivery": "string"})
	if row["evidence"] != identity || row["manifest_verified"] != true || row["delivery"] != "unverified" ||
		row["pages_verified"] != prepared["pages"] || row["sources_verified"] != prepared["sources"] {
		t.Fatalf("verified row = %v, want the prepared page and source counts %v", row, prepared)
	}
}
