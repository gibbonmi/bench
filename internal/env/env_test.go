package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeAllow writes a .bench/env.allow fixture and returns the repo root
// Build/parseAllowFile should read it from.
func writeAllow(t *testing.T, body string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "env.allow"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestParseAllowFileMissingIsDefaultsOnly covers: absent env.allow is never an
// error and yields no additions.
func TestParseAllowFileMissingIsDefaultsOnly(t *testing.T) {
	root := t.TempDir()
	got, err := parseAllowFile(filepath.Join(root, ".bench", "env.allow"))
	if err != nil {
		t.Fatalf("missing env.allow must not error, got %v", err)
	}
	if len(got.agent) != 0 {
		t.Fatalf("missing env.allow must yield no entries, got %+v", got)
	}
}

// TestParseAllowFileEmptyIsAcceptedAsDefaults covers the edge row: an env.allow
// that is present but empty is accepted and yields the defaults (distinct from
// the absent-file case, per the profile checklist).
func TestParseAllowFileEmptyIsAcceptedAsDefaults(t *testing.T) {
	root := writeAllow(t, "")
	got, err := parseAllowFile(filepath.Join(root, ".bench", "env.allow"))
	if err != nil {
		t.Fatalf("empty env.allow must not error, got %v", err)
	}
	if len(got.agent) != 0 {
		t.Fatalf("empty env.allow must yield no entries, got %+v", got)
	}
}

// TestParseAllowMissingFinalNewlineParsesLastLine covers the edge row: a file
// whose last line lacks a trailing newline still parses that line normally,
// against a naive split that drops the final entry.
func TestParseAllowMissingFinalNewlineParsesLastLine(t *testing.T) {
	got, err := parseAllow("[agent]\nMY_VAR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.agent) != 1 || got.agent[0] != "MY_VAR" {
		t.Fatalf("last entry without trailing newline was dropped, got %+v", got)
	}
}

// TestParseAllowAgentSection covers story 3's scoping: an entry under the
// [agent] section — the only known section — lands in the agent result.
func TestParseAllowAgentSection(t *testing.T) {
	got, err := parseAllow("[agent]\nAGENT_ONLY\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.agent) != 1 || got.agent[0] != "AGENT_ONLY" {
		t.Fatalf("agent section wrong: %+v", got)
	}
}

// TestValidateEntryRejectsBareAndMidNameGlob covers the edge row: a bare *
// entry and a mid-name glob are both rejected, and parseAllow reports both by
// line number.
func TestValidateEntryRejectsBareAndMidNameGlob(t *testing.T) {
	cases := []struct {
		name string
		body string
		line int
	}{
		{"bare wildcard", "[agent]\n*\n", 2},
		{"mid-name glob", "[agent]\nFOO*BAR\n", 2},
		{"double trailing glob", "[agent]\nFOO**\n", 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseAllow(tc.body)
			if err == nil {
				t.Fatal("expected rejection, got nil error")
			}
			wantPrefix := fmt.Sprintf(".bench/env.allow:%d:", tc.line)
			if !strings.HasPrefix(err.Error(), wantPrefix) {
				t.Fatalf("error %q does not name line %d", err, tc.line)
			}
		})
	}
}

// TestParseAllowRejectsMalformedLines walks the remaining rejection reasons
// the grammar names, each pinned to its offending line number: an entry
// before any section header, an unknown section name (including a stale [gate],
// whose opt-in home is the manifest and not this file), an entry containing /
// or =, and a character outside the portable environment-name set.
func TestParseAllowRejectsMalformedLines(t *testing.T) {
	cases := []struct {
		name string
		body string
		line int
	}{
		{"entry before section header", "STRAY_VAR\n[agent]\n", 1},
		{"unknown section name", "[nope]\nFOO\n", 1},
		{"stale gate section", "[gate]\nFOO\n", 1},
		{"entry with slash", "[agent]\nFOO/BAR\n", 2},
		{"entry with equals", "[agent]\nFOO=BAR\n", 2},
		{"entry with invalid character", "[agent]\nFOO-BAR\n", 2},
		{"entry starting with digit", "[agent]\n1FOO\n", 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseAllow(tc.body)
			if err == nil {
				t.Fatal("expected rejection, got nil error")
			}
			wantPrefix := fmt.Sprintf(".bench/env.allow:%d:", tc.line)
			if !strings.HasPrefix(err.Error(), wantPrefix) {
				t.Fatalf("error %q does not name line %d", err, tc.line)
			}
		})
	}
}

// TestParseAllowSkipsCommentLines covers the edge row (review C2): a `#`
// comment line is skipped rather than rejected as a malformed entry.
// DATA_HANDLING.md advertises `#` comments as part of the grammar; a
// regression that dropped the skip branch would fail-close every env.allow
// file that uses one.
func TestParseAllowSkipsCommentLines(t *testing.T) {
	got, err := parseAllow("[agent]\n# a comment\nMY_VAR\n")
	if err != nil {
		t.Fatalf("comment line must not error, got %v", err)
	}
	if len(got.agent) != 1 || got.agent[0] != "MY_VAR" {
		t.Fatalf("comment line was not skipped cleanly, got %+v", got)
	}
}

// TestParseAllowAcceptsCRLFLineEndings covers the edge row (review C3): a file
// with CRLF line endings parses the same as LF-only, because TrimSpace strips
// the trailing \r along with the line's other whitespace.
func TestParseAllowAcceptsCRLFLineEndings(t *testing.T) {
	got, err := parseAllow("[agent]\r\nMY_VAR\r\n")
	if err != nil {
		t.Fatalf("CRLF line endings must not error, got %v", err)
	}
	if len(got.agent) != 1 || got.agent[0] != "MY_VAR" {
		t.Fatalf("CRLF entry was not parsed cleanly, got %+v", got)
	}
}

// TestParseAllowRejectsUTF8BOMByName covers the edge row (review C3): a file
// that opens with a UTF-8 byte-order mark is rejected fail-closed, and the
// error names the BOM directly rather than reporting the garbled first line
// as an "entry before any section header" — the BOM is not Unicode
// whitespace, so a naive TrimSpace-only diagnosis would misname the cause.
func TestParseAllowRejectsUTF8BOMByName(t *testing.T) {
	_, err := parseAllow("\ufeff[agent]\nMY_VAR\n")
	if err == nil {
		t.Fatal("expected rejection, got nil error")
	}
	if !strings.Contains(err.Error(), "byte-order mark") {
		t.Fatalf("error %q does not name the BOM", err)
	}
}

// TestParseAllowDuplicateEntriesAreAppendedAsIs covers the edge row (review
// C4): a duplicate entry is not deduplicated — both copies are appended,
// which is harmless because matchesAny only needs one match. This pins
// current behavior; there is no dedup branch to change.
func TestParseAllowDuplicateEntriesAreAppendedAsIs(t *testing.T) {
	got, err := parseAllow("[agent]\nMY_VAR\nMY_VAR\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.agent) != 2 || got.agent[0] != "MY_VAR" || got.agent[1] != "MY_VAR" {
		t.Fatalf("duplicate entry was not appended as-is, got %+v", got)
	}
}

// TestParseAllowLeadingWhitespaceEntryIsLenientlyAccepted covers the edge row
// (review C4): an entry line with leading whitespace is accepted, not
// rejected — TrimSpace normalizes it before the grammar checks run. This
// pins current lenient behavior as-is.
func TestParseAllowLeadingWhitespaceEntryIsLenientlyAccepted(t *testing.T) {
	got, err := parseAllow("[agent]\n   MY_VAR\n")
	if err != nil {
		t.Fatalf("leading whitespace before an entry must not error, got %v", err)
	}
	if len(got.agent) != 1 || got.agent[0] != "MY_VAR" {
		t.Fatalf("leading-whitespace entry was not accepted leniently, got %+v", got)
	}
}

// TestDefaultGlobsDoNotStraddleFamilies covers the edge row: no default glob in
// the agent passlist matches a name belonging to a different family, checked by
// enumerating each default glob against a fixture of foreign names. The
// GOOGLE_* Vertex-routing glob is the concrete GO* versus GOOGLE_* hazard the
// retired Go-toolchain enumeration draft existed to avoid; it legitimately
// admits GOOGLE_* here, so the fixture keeps foreign names that no default glob
// should match and asserts GOOGLE_APPLICATION_CREDENTIALS is admitted rather
// than treated as foreign.
func TestDefaultGlobsDoNotStraddleFamilies(t *testing.T) {
	foreign := []string{
		"GITLAB_TOKEN",
		"NODEBUG_SOMETHING",
		"AZURE_CLIENT_SECRET",
		"DOCKER_HOST",
	}
	for _, name := range foreign {
		if matchesAny(name, AgentPasslist) {
			t.Errorf("agent passlist matched foreign name %q", name)
		}
	}

	// GOOGLE_* is a documented Vertex-routing default, so
	// GOOGLE_APPLICATION_CREDENTIALS is legitimately admitted here rather than
	// foreign. Assert that directly so the family's treatment is pinned.
	if !matchesAny("GOOGLE_APPLICATION_CREDENTIALS", AgentPasslist) {
		t.Fatal("agent passlist should admit GOOGLE_* for documented Vertex routing")
	}
	// No default glob may be a bare GO* — the straddle that would also match
	// GOOGLE_APPLICATION_CREDENTIALS in a context that does not want it.
	for _, p := range AgentPasslist {
		if p == "GO*" {
			t.Fatal("agent passlist must not contain a GO* glob")
		}
	}
}

// TestBuildPassesMultilineAndLargeValuesUnaltered covers the edge row: a
// passlisted variable holding a multi-line or very large value passes through
// unaltered — the passlist filters names, never values.
func TestBuildPassesMultilineAndLargeValuesUnaltered(t *testing.T) {
	root := t.TempDir()
	multiline := "line one\nline two\nline three"
	large := strings.Repeat("x", 200000)

	t.Setenv("ANTHROPIC_API_KEY", multiline)
	t.Setenv("BENCH_LARGE_TEST_VALUE", large)

	got, err := Build(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantMultiline := "ANTHROPIC_API_KEY=" + multiline
	wantLarge := "BENCH_LARGE_TEST_VALUE=" + large
	foundMultiline, foundLarge := false, false
	for _, kv := range got {
		if kv == wantMultiline {
			foundMultiline = true
		}
		if kv == wantLarge {
			foundLarge = true
		}
	}
	if !foundMultiline {
		t.Fatal("multi-line value was altered or dropped")
	}
	if !foundLarge {
		t.Fatal("large value was altered or dropped")
	}
}

// TestBuildFiltersToPasslistByName pins Build as a filter over parent names,
// not a wholesale inherit and not a wipe: a marker name is absent while a
// passlisted name survives.
func TestBuildFiltersToPasslistByName(t *testing.T) {
	root := t.TempDir()
	t.Setenv("ENV_TEST_MARKER_UNLISTED", "leak")
	t.Setenv("PATH", os.Getenv("PATH"))

	got, err := Build(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sawMarker, sawPath := false, false
	for _, kv := range got {
		if strings.HasPrefix(kv, "ENV_TEST_MARKER_UNLISTED=") {
			sawMarker = true
		}
		if strings.HasPrefix(kv, "PATH=") {
			sawPath = true
		}
	}
	if sawMarker {
		t.Fatal("unlisted marker name leaked into build")
	}
	if !sawPath {
		t.Fatal("PATH, a shared basic, did not survive into build")
	}
}

// TestBuildOptInGlobAdmitsAllMatchingNames covers story 3's glob-admission
// row: a PREFIX* entry in env.allow admits every matching name in the [agent]
// section, not just an exact match.
func TestBuildOptInGlobAdmitsAllMatchingNames(t *testing.T) {
	root := writeAllow(t, "[agent]\nMYCO_*\n")
	t.Setenv("MYCO_TOKEN", "a")
	t.Setenv("MYCO_REGION", "b")
	t.Setenv("OTHERCO_TOKEN", "c")

	got, err := Build(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := map[string]bool{}
	for _, kv := range got {
		name, _, _ := strings.Cut(kv, "=")
		names[name] = true
	}
	if !names["MYCO_TOKEN"] || !names["MYCO_REGION"] {
		t.Fatalf("glob entry did not admit all matching names: %v", got)
	}
	if names["OTHERCO_TOKEN"] {
		t.Fatal("glob entry admitted a non-matching name")
	}
}

// TestBuildMalformedAllowFailsClosed covers story 4: a malformed env.allow
// refuses to build an environment at all, returning an error naming the line
// number rather than degrading to defaults.
func TestBuildMalformedAllowFailsClosed(t *testing.T) {
	root := writeAllow(t, "[agent]\n*\n")
	_, err := Build(root)
	if err == nil {
		t.Fatal("expected a fail-closed error for a bare wildcard entry")
	}
	if !strings.Contains(err.Error(), ".bench/env.allow:2:") {
		t.Fatalf("error %q does not name the offending line", err)
	}
}
