package roadmap

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/learnings"
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

func TestBuildContextProjectsPendingCaptureOccurrences(t *testing.T) {
	root := newRepo(t)
	for path, body := range map[string]string{
		RoadmapFile:            "**FT1 — one.**\n\n**FT2 — two.**\n",
		IdeasFile:              "- 2026-07-10  later [occurrence:FT2/idea-2]\n- 2026-07-10  first [occurrence:FT1/idea-1]\n",
		learnings.JournalPath:  "## 2026-07-10 — lesson  [open]\nbody [occurrence:FT1/learning-1]\n",
		".bench/retros/one.md": "## Agent-experience improvements\n\nParagraph [occurrence:FT2/retro-1]\n",
	} {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, path)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, path), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := renderContext(s)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"FT1,idea-1,IDEAS.md,line 2,pending",
		"FT1,learning-1,.bench/learnings.md,line 1,pending",
		"FT2,idea-2,IDEAS.md,line 1,pending",
		"FT2,retro-1,.bench/retros/one.md,line 3,pending",
	}
	last := -1
	for _, row := range want {
		at := strings.Index(out, row)
		if at < 0 || at <= last {
			t.Fatalf("capture rows = %s", out)
		}
		last = at
	}
}

func TestBuildContextClassifiesOccurrenceDiscrepancies(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\nOccurrences: recorded\n\n**FT2 — two.**\nOccurrences: bad, bad\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ideas := strings.Join([]string{
		"- 2026-07-10  duplicate [occurrence:FT1/recorded]",
		"- 2026-07-10  malformed [occurrence:FT1/unterminated",
		"- 2026-07-10  unknown [occurrence:FT9/new]",
		"- 2026-07-10  multiple [occurrence:FT1/one] [occurrence:FT1/two]",
	}, "\n")
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte(ideas), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if s.SequenceTrusted {
		t.Fatal("structural discrepancies left sequence trusted")
	}
	out, err := renderContext(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []string{
		"FT1,recorded,IDEAS.md,line 1,already-recorded",
		"IDEAS.md,line 1,already-recorded,FT1,recorded,false",
		"IDEAS.md,line 2,malformed-token,FT1,unterminated,true",
		"IDEAS.md,line 3,unknown-owner,FT9,new,true",
		"IDEAS.md,line 4,multiple-tokens,FT1,two,true",
		"ROADMAP.md,FT2,malformed-ledger,FT2,\"\",true",
	} {
		if !strings.Contains(out, row) {
			t.Fatalf("missing %q in %s", row, out)
		}
	}
}

func TestBuildContextTreatsSeparatedValidTokensAsMultiple(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte("- 2026-07-10  first [occurrence:FT1/one] with prose between [occurrence:FT1/two]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if s.SequenceTrusted || len(s.CaptureOccurrences) != 0 {
		t.Fatalf("multiple tokens = trusted:%v occurrences:%#v", s.SequenceTrusted, s.CaptureOccurrences)
	}
	if got, want := s.Roadmap.OccurrenceDiscrepancies, []OccurrenceDiscrepancy{{Source: IdeasFile, CaptureUnit: "line 1", Kind: "multiple-tokens", Owner: "FT1", Incident: "two", Structural: true}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("discrepancies = %#v, want %#v", got, want)
	}
}

func TestBuildContextKeepsLearningAndRetroOccurrencesInTheirUnits(t *testing.T) {
	root := newRepo(t)
	for path, body := range map[string]string{
		RoadmapFile:           "**FT1 — one.**\n\n**FT2 — two.**\n",
		learnings.JournalPath: "## 2026-07-10 — lesson  [open]\nProse [occurrence:FT2/ignored_] does not end the entry.\nFinal [occurrence:FT1/learning-final]\n",
		".bench/retros/one.md": strings.Join([]string{
			"Outside [occurrence:FT2/ignored]",
			"",
			"## Agent-experience improvements",
			"",
			"Paragraph [occurrence:FT1/retro-paragraph]",
			"",
			"- First item [occurrence:FT1/retro-list]",
			"- Second item [occurrence:FT2/retro-other-owner]",
			"",
			"## Other",
			"Ignored [occurrence:FT2/ignored-too]",
		}, "\n"),
	} {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, path)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, path), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := renderContext(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []string{
		"FT1,learning-final,.bench/learnings.md,line 1,pending",
		"FT1,retro-list,.bench/retros/one.md,line 7,pending",
		"FT1,retro-paragraph,.bench/retros/one.md,line 5,pending",
		"FT2,retro-other-owner,.bench/retros/one.md,line 8,pending",
	} {
		if !strings.Contains(out, row) {
			t.Fatalf("missing %q in %s", row, out)
		}
	}
	for _, incident := range []string{"ignored", "ignored-too"} {
		if strings.Contains(out, "FT2,"+incident+",") {
			t.Fatalf("non-final token %q became an occurrence: %s", incident, out)
		}
	}
}

func TestBuildContextNormalizesPendingOccurrencePairs(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte("- 2026-07-10  first [occurrence:FT1/shared]\n- 2026-07-10  second [occurrence:FT1/shared]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if len(s.CaptureOccurrences) != 2 {
		t.Fatalf("capture occurrences = %#v", s.CaptureOccurrences)
	}
	if got, want := s.PendingOccurrences, []OccurrencePair{{Owner: "FT1", Incident: "shared"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("pending pairs = %#v, want %#v", got, want)
	}
}

func TestBuildContextKeepsSameIncidentForSeparateOwners(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\n\n**FT2 — two.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte("- 2026-07-10  first [occurrence:FT2/shared]\n- 2026-07-10  second [occurrence:FT1/shared]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	want := []OccurrencePair{{Owner: "FT1", Incident: "shared"}, {Owner: "FT2", Incident: "shared"}}
	if !reflect.DeepEqual(s.PendingOccurrences, want) {
		t.Fatalf("pending pairs = %#v, want %#v", s.PendingOccurrences, want)
	}
}

func TestBuildContextProjectsEveryRecordedSourceWithoutPendingPair(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\nOccurrences: recorded\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte("- 2026-07-10  first [occurrence:FT1/recorded]\n- 2026-07-10  second [occurrence:FT1/recorded]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if !s.SequenceTrusted {
		t.Fatal("only advisory discrepancies made the sequence untrusted")
	}
	if len(s.PendingOccurrences) != 0 {
		t.Fatalf("recorded pair remained pending: %#v", s.PendingOccurrences)
	}
	if got, want := s.CaptureOccurrences, []CaptureOccurrence{
		{Owner: "FT1", Incident: "recorded", Source: IdeasFile, CaptureUnit: "line 1", State: "already-recorded"},
		{Owner: "FT1", Incident: "recorded", Source: IdeasFile, CaptureUnit: "line 2", State: "already-recorded"},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("capture occurrences = %#v, want %#v", got, want)
	}
	if got, want := s.Roadmap.OccurrenceDiscrepancies, []OccurrenceDiscrepancy{
		{Source: IdeasFile, CaptureUnit: "line 1", Kind: "already-recorded", Owner: "FT1", Incident: "recorded"},
		{Source: IdeasFile, CaptureUnit: "line 2", Kind: "already-recorded", Owner: "FT1", Incident: "recorded"},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("discrepancies = %#v, want %#v", got, want)
	}
	out, err := renderContext(s)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "3,false,true") {
		t.Fatalf("trusted context header missing: %s", out)
	}
	if !strings.Contains(out, "IDEAS.md,line 1,already-recorded,FT1,recorded,false") {
		t.Fatalf("advisory discrepancy missing from context: %s", out)
	}
}

func TestBuildContextRequiresUsableCaptureSourcesForOccurrenceTrust(t *testing.T) {
	for _, tc := range []struct {
		name    string
		prepare func(t *testing.T, root string)
	}{
		{
			name: "ideas wrong type",
			prepare: func(t *testing.T, root string) {
				t.Helper()
				if err := os.Mkdir(filepath.Join(root, IdeasFile), 0o755); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "learnings malformed",
			prepare: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, learnings.JournalPath)
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte{0xff}, 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "retros wrong type",
			prepare: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, ".bench", "retros")
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("not a directory"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			tc.prepare(t, root)

			s, err := BuildContext(root, false, GateCacheFact{})
			if err != nil {
				t.Fatal(err)
			}
			if s.SequenceTrusted {
				t.Fatal("degraded capture source left recurrence sequence trusted")
			}
			if _, err := renderContext(s); err != nil {
				t.Fatalf("complete context did not render: %v", err)
			}
		})
	}
}
