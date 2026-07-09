package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// signalRowRe matches the emitted signal-name literal in internal/status/status.go's
// row{<severity>, "<name>", …} constructor: the first string after the integer
// severity is the signal name shown on the ambient dashboard.
var signalRowRe = regexp.MustCompile(`row\{\s*\d+\s*,\s*"([a-z][a-z-]*)"`)

// checkSignalVocabulary cross-checks the signal names internal/status/status.go
// emits against CONTEXT.md's signal-definition enumeration, so the file that pins
// the ubiquitous language can't silently drift when a signal is added to status.go
// (the drift the assessment actually observed). Prior art checkLineBinding: read
// the graded tree's machine source, assert the prose names each value.
//
//   - Forward only — every emitted signal name must appear in CONTEXT.md — matching
//     checkLineBinding's one-directional prose check.
//   - Scoped match, not whole-file Contains: signal names are common words (git,
//     specs, reviews) that recur throughout CONTEXT.md, so a whole-file substring
//     test false-passes. Each name is matched as a delimited token within the
//     parenthesized enumeration on the `signal` term, nowhere else.
//   - Absent source no-ops: an unreadable status.go or CONTEXT.md (or a signal term
//     with no parenthesized enumeration) returns nothing, mirroring checkLineBinding
//     skipping an empty profile — the compiled-core build/test check already fails a
//     tree missing status.go.
//   - Zero-extraction floor: a present status.go whose row shape has drifted from
//     signalRowRe (a named-const severity, a helper constructor, …) yields zero
//     matches. That must not silently pass as "nothing to check" — it reds with one
//     diagnostic naming the extraction as stale, so the drift is caught instead of
//     going quiet.
func checkSignalVocabulary(root string) []string {
	status := readIfExists(filepath.Join(root, "internal", "status", "status.go"))
	context := readIfExists(filepath.Join(root, "CONTEXT.md"))
	if status == "" || context == "" {
		return nil
	}
	enum := signalEnumeration(context)
	if enum == "" {
		return nil
	}
	matches := signalRowRe.FindAllStringSubmatch(status, -1)
	if len(matches) == 0 {
		return []string{"internal/status/status.go yielded no signal names — signalRowRe no longer matches the row shape; update the extraction"}
	}
	var diags []string
	seen := map[string]bool{}
	for _, match := range matches {
		name := match[1]
		if seen[name] {
			continue
		}
		seen[name] = true
		if !tokenInEnumeration(enum, name) {
			diags = append(diags, fmt.Sprintf("CONTEXT.md signal enumeration missing signal '%s' (internal/status/status.go emits it)", name))
		}
	}
	return diags
}

// signalEnumeration returns the parenthesized list on CONTEXT.md's `signal` term —
// the first "(...)" group after the "**signal**" marker. Empty when the term or its
// enumeration is absent.
func signalEnumeration(context string) string {
	idx := strings.Index(context, "**signal**")
	if idx < 0 {
		return ""
	}
	rest := context[idx:]
	open := strings.Index(rest, "(")
	if open < 0 {
		return ""
	}
	end := strings.Index(rest[open:], ")")
	if end < 0 {
		return ""
	}
	return rest[open+1 : open+end]
}

// tokenInEnumeration reports whether name appears in enum delimited by non-letters,
// so a common name (git, specs) matches only as a whole word and never inside a
// longer word.
func tokenInEnumeration(enum, name string) bool {
	return regexp.MustCompile(`(^|[^a-zA-Z])` + regexp.QuoteMeta(name) + `([^a-zA-Z]|$)`).MatchString(enum)
}

// TestSignalVocabularyBites is the recorded bite proof (per craft-gate): an intact
// tree whose CONTEXT.md enumeration names every emitted signal passes clean;
// dropping one name from the enumeration fires exactly that signal's diagnostic and
// no sibling; the scoped match refuses a whole-file false-pass (a dropped name that
// still occurs elsewhere in CONTEXT.md outside the enumeration stays red); and an
// absent status.go no-ops.
func TestSignalVocabularyBites(t *testing.T) {
	const statusSrc = "package status\n" +
		"func rows() {\n" +
		"\t_ = row{7, \"gate\", d, a}\n" +
		"\t_ = row{1, \"git\", d, a}\n" +
		"\t_ = row{10, \"roadmap\", d, a}\n" +
		"}\n"
	// enum builds a CONTEXT.md whose signal term names exactly the given signals,
	// with a decoy sentence afterward that mentions git/roadmap outside the
	// enumeration so a whole-file Contains impl would false-pass.
	enum := func(names ...string) string {
		return "# Context\n\n" +
			"- **signal** — one ranked line on the ambient dashboard (" + strings.Join(names, ", ") + "). Not \"check\".\n" +
			"- **gate cache** — the git dir cache the roadmap owner reads.\n"
	}
	write := func(t *testing.T, status, context string) string {
		t.Helper()
		root := t.TempDir()
		if status != "" {
			p := filepath.Join(root, "internal", "status")
			if err := os.MkdirAll(p, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(p, "status.go"), []byte(status), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(root, "CONTEXT.md"), []byte(context), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}

	if diags := checkSignalVocabulary(write(t, statusSrc, enum("gate", "git", "roadmap"))); len(diags) != 0 {
		t.Fatalf("intact enumeration: want no diagnostics, got %v", diags)
	}

	cases := []struct {
		name string
		drop string
		want string
	}{
		{"git dropped", "git", "CONTEXT.md signal enumeration missing signal 'git'"},
		{"roadmap dropped", "roadmap", "CONTEXT.md signal enumeration missing signal 'roadmap'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var kept []string
			for _, n := range []string{"gate", "git", "roadmap"} {
				if n != tc.drop {
					kept = append(kept, n)
				}
			}
			// The decoy line in enum() still mentions git and roadmap outside the
			// parenthesized list, so a whole-file Contains impl would miss this.
			diags := checkSignalVocabulary(write(t, statusSrc, enum(kept...)))
			if !containsDiagnostic(diags, tc.want) {
				t.Fatalf("want %q in diagnostics, got %v", tc.want, diags)
			}
			for _, other := range cases {
				if other.want != tc.want && containsDiagnostic(diags, other.want) {
					t.Fatalf("%s also fired sibling %q: %v", tc.name, other.want, diags)
				}
			}
		})
	}

	if diags := checkSignalVocabulary(write(t, "", enum("gate"))); len(diags) != 0 {
		t.Fatalf("absent status.go: want no diagnostics (skip posture), got %v", diags)
	}

	// Zero-extraction floor: a present status.go whose row shape has drifted from
	// signalRowRe (here a named-const severity, `row{sevGate, ...}`, instead of the
	// integer literal the regex requires) must not silently pass with zero
	// diagnostics — it must red with a diagnostic naming the extraction as stale.
	const driftedStatusSrc = "package status\n" +
		"func rows() {\n" +
		"\t_ = row{sevGate, \"gate\", d, a}\n" +
		"}\n"
	const wantFloorDiag = "internal/status/status.go yielded no signal names — signalRowRe no longer matches the row shape; update the extraction"
	if diags := checkSignalVocabulary(write(t, driftedStatusSrc, enum("gate"))); !containsDiagnostic(diags, wantFloorDiag) {
		t.Fatalf("drifted row shape (zero extraction): want %q in diagnostics, got %v", wantFloorDiag, diags)
	}
}
