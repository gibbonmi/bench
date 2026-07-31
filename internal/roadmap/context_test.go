package roadmap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestContextBodyLimitBoundaries(t *testing.T) {
	for _, n := range []int{4095, 4096, 4097} {
		got, full, truncated := limited(strings.Repeat("x", n), false)
		if full != n || truncated != (n > 4096) || len(got) != min(n, 4096) {
			t.Fatalf("n=%d got bytes=%d full=%d truncated=%v", n, len(got), full, truncated)
		}
	}
	s := strings.Repeat("x", 4095) + "€"
	got, full, truncated := limited(s, false)
	if full != 4098 || !truncated || !utf8.ValidString(got) || len(got) != 4095 {
		t.Fatalf("UTF-8 boundary got bytes=%d full=%d truncated=%v valid=%v", len(got), full, truncated, utf8.ValidString(got))
	}
	got, full, truncated = limited(s, true)
	if got != s || full != 4098 || truncated {
		t.Fatal("--full did not preserve body")
	}
}

func TestBuildContextCarriesRetrosAndDegradedEvidence(t *testing.T) {
	root := newRepo(t)
	dir := filepath.Join(root, ".bench", "retros")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.md"), []byte("second"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("first"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing.md", filepath.Join(dir, "bad.md")); err != nil {
		t.Fatal(err)
	}
	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Retros) != 3 || s.Retros[0].Path != ".bench/retros/a.md" || s.Retros[1].Path != ".bench/retros/b.md" || s.Retros[2].State != "unreadable" {
		t.Fatalf("retros = %#v", s.Retros)
	}
	if len(s.Sources) < 2 || s.Sources[len(s.Sources)-2].Source != ".bench/retros/" {
		t.Fatalf("sources = %#v", s.Sources)
	}
	if got, err := renderContext(s); err != nil || !strings.Contains(got, "retros[3]{path,state,body,body_bytes,truncated}:") || !strings.Contains(got, ".bench/retros/bad.md,unreadable") || !strings.Contains(got, "parse_failures[1]{") {
		t.Fatalf("context = %q, %v", got, err)
	}
}

func TestParseDocumentOccurrenceLedgers(t *testing.T) {
	content := []byte("**FT1 — one.** Body\nOccurrences: alpha-1, beta-2\n\n**FT2 — two.** Body\nOccurrences: bad, bad\n")
	doc, failures := ParseDocument(content, nil, false)
	if got := doc.Rows[0]; got.OccurrenceKeys != "alpha-1, beta-2" || got.OccurrenceCount != 2 {
		t.Fatalf("valid ledger = %#v", got)
	}
	if len(failures) != 1 || failures[0].Reason != "malformed-ledger" {
		t.Fatalf("failures = %#v", failures)
	}
}

func TestOccurrenceIncidentGrammar(t *testing.T) {
	for _, key := range []string{"a", strings.Repeat("a", 64), "a-1"} {
		if !ValidOccurrenceIncident(key) {
			t.Fatalf("valid incident rejected: %q", key)
		}
	}
	for _, key := range []string{"", strings.Repeat("a", 65), "A", "é", "a_b", "-a", "a-", "a\tb", "a\rb", "a\nb"} {
		if ValidOccurrenceIncident(key) {
			t.Fatalf("invalid incident accepted: %q", key)
		}
	}
}

func TestOccurrenceLedgerMalformedAndLineEndings(t *testing.T) {
	valid, failures := ParseDocument([]byte("**FT1 — one.**\r\nOccurrences: alpha-1, beta-2"), nil, false)
	if len(failures) != 0 || valid.Rows[0].OccurrenceCount != 2 {
		t.Fatalf("CRLF newline-less ledger = %#v, %#v", valid, failures)
	}
	for _, ledger := range []string{"Occurrences:", "Occurrences: alpha_1", "Occurrences: beta, alpha", "Occurrences: alpha, alpha", "Occurrences: alpha\nOccurrences: beta"} {
		doc, got := ParseDocument([]byte("**FT1 — one.**\n"+ledger+"\n"), nil, false)
		if len(got) != 1 || got[0].Reason != "malformed-ledger" || len(doc.OccurrenceDiscrepancies) != 1 {
			t.Fatalf("ledger %q accepted: %#v, %#v", ledger, doc, got)
		}
	}
}
