package maps

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

// parseFile is the shared engine; these pin the marker/fence/handoff edges at the
// pure seam the two-derivations bug class breeds.
func TestParseFileMarkers(t *testing.T) {
	r := parseFile([]byte("## #1: a?\nType: Grill\n### Answer\n— (open)\n"))
	if len(r.tickets) != 1 || r.tickets[0] != (ticket{"1", "Grill", "open"}) {
		t.Fatalf("open ticket = %+v", r.tickets)
	}
	if !r.preHandoffMarker || !r.notCloseReady() {
		t.Errorf("open ticket file should be not-close-ready")
	}

	// A mid-line GRILL DEFERRED mention and a fenced placeholder are not markers.
	r = parseFile([]byte("## #1: a?\nType: Grill\n### Answer\nDecided: mid-line GRILL DEFERRED mention.\n\n```\n— (open)\n```\n"))
	if len(r.tickets) != 0 {
		t.Errorf("over-match: unexpected tickets %+v", r.tickets)
	}

	// No Type line → unknown.
	r = parseFile([]byte("## #1: a?\n### Answer\n— (open)\n"))
	if r.tickets[0].typ != "unknown" {
		t.Errorf("typeless ticket type = %q, want unknown", r.tickets[0].typ)
	}
}

func TestParseFileHandoff(t *testing.T) {
	cases := []struct {
		name         string
		in           string
		wantRows     [][]any
		wantNotReady bool
	}{
		{"missing handoff", "## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n",
			[][]any{{"m", "handoff", "handoff", "missing"}}, true},
		{"filled handoff silent", "## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n## Handoff\n1. a. n/a\n2. b. n/a\n",
			nil, false},
		{"placeholder in handoff", "## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n## Handoff\n1. a.\n— (open)\n",
			[][]any{{"m", "handoff", "handoff", "open"}}, true},
		{"fenced handoff is missing", "## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n```\n## Handoff\n1. a.\n```\n",
			[][]any{{"m", "handoff", "handoff", "missing"}}, true},
		{"open ticket no handoff row", "## #1: q?\nType: Grill\n### Answer\n— (open)\n",
			[][]any{{"m", 1, "Grill", "open"}}, true},
		{"non-map never nagged", "# Index\nprose, not a map.\n", nil, false},
	}
	for _, c := range cases {
		r := parseFile([]byte(c.in))
		if got := fileRows("m", r); !reflect.DeepEqual(got, c.wantRows) {
			t.Errorf("%s: fileRows = %v, want %v", c.name, got, c.wantRows)
		}
		if r.notCloseReady() != c.wantNotReady {
			t.Errorf("%s: notCloseReady = %v, want %v", c.name, r.notCloseReady(), c.wantNotReady)
		}
	}
}

// UnresolvedCount is DISTINCT not-close-ready files — re-homes the shell
// "AXI maps_unresolved_count distinct-file contract" and the close-readiness count
// tail, both of which used to source bin/bench-query.sh.
func TestUnresolvedCount(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Two unresolved tickets in one file → 3 ticket rows but distinct-file count of the
	// two files below is 2.
	write("multi.md", "## #1: a?\nType: Grill\n### Answer\n— (open)\n\n## #2: b?\nType: Grill\n### Answer\n— (deferred)\n")
	write("solo.md", "## #1: c?\nType: Grill\n### Answer\n— (open)\n")
	// A dotfile the shell glob would never expand must stay invisible to rows and count.
	write(".hidden.md", "## #1: h?\nType: Grill\n### Answer\n— (open)\n")
	if got := len(Rows(root)); got != 3 {
		t.Errorf("Rows count = %d, want 3", got)
	}
	if got, state := UnresolvedCount(root); got != 2 || state != bounds.StateParsed {
		t.Errorf("UnresolvedCount = (%d, %s), want (2, %s)", got, state, bounds.StateParsed)
	}
}

// The close-readiness aggregate re-homes the shell contract's count tail: of the six
// files, hm/hx/hp/ho are not-close-ready, hf is ready, and README documents the
// directory rather than claiming to be a map, so the tally never sees it → 4.
func TestUnresolvedCountCloseReadiness(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"hm.md":     "# HM\n## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n",
		"hf.md":     "# HF\n## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n## Handoff\n1. a. n/a\n",
		"ho.md":     "# HO\n## #1: q?\nType: Grill\n### Answer\n— (open)\n",
		"hx.md":     "# HX\n## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n```\n## Handoff\n1. a.\n```\n",
		"hp.md":     "# HP\n## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n## Handoff\n1. a.\n— (open)\n",
		"README.md": "# Index\nnot a map.\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got, _ := UnresolvedCount(root); got != 4 {
		t.Errorf("close-readiness UnresolvedCount = %d, want 4", got)
	}
	// A file-scope marker without any ticket heading is a recognized shape (a marker),
	// not unsupported-schema: no listed row, but counted via preHandoffMarker.
	if err := os.WriteFile(filepath.Join(dir, "scope.md"), []byte("### Answer\n— (deferred)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := len(fileRows("scope", parseFile([]byte("### Answer\n— (deferred)\n")))); got != 0 {
		t.Errorf("file-scope marker emitted %d rows, want 0", got)
	}
	if got, _ := UnresolvedCount(root); got != 5 {
		t.Errorf("UnresolvedCount with file-scope marker = %d, want 5", got)
	}
}
