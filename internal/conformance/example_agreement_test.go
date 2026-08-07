package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/specbuild"
)

// The craft-tickets Good example sits between these markers, each at column 0 and
// outside the fenced block, so both fence lines fall inside the extracted region
// and neither the template block nor the Bad example can be graded by mistake.
const (
	ticketExampleBegin      = "<!-- ticket-example:begin -->"
	ticketExampleEnd        = "<!-- ticket-example:end -->"
	ticketExampleOpenFence  = "```markdown"
	ticketExampleCloseFence = "```"
)

var ticketExampleDoc = filepath.Join(".agents", "skills", "bench-craft-tickets", "SKILL.md")

// The taught contract, authored here from a reading of the skill rather than read
// back out of the block under test at run time. That independence is the whole
// check: an example that drifts off the parseable shape while its prose stays
// plausible turns these comparisons red, where a self-derived expectation would
// follow the drift and stay green.
var (
	taughtExampleRows    = []string{"RC1", "RC2"}
	taughtExampleFence   = []string{"internal/status", "internal/render/rows.go"}
	taughtExampleCovers  = []string{"CJ1", "CJ2"}
	taughtExampleClosure = []string{"RC1/reason", "RC2/recovery-action"}
)

// The mutations heading is matched at any depth so a ticket that nests the
// example under a parent heading still grades; the Blocked by line has to carry
// a value, since a bare prefix demonstrates nothing to a copying reader.
var (
	ticketExampleMutationsHeading = regexp.MustCompile(`(?m)^#{2,6}\s+Red mutations\s*$`)
	ticketExampleBlockedByLine    = regexp.MustCompile(`(?m)^Blocked by:[ \t]+\S`)
)

const ticketExampleDiag = "craft-tickets example agreement: "

// checkExampleAgreement asserts the craft-tickets Good example is a ticket the
// real assignment parser reads, not just prose that looks like one. It extracts
// the marker-delimited block, materializes it as a ticket file beside a temp spec
// path (the parser takes paths, not bytes), runs specbuild.ParseTicket over it,
// and compares the parse against the independently authored literals above.
//
// Every extraction failure is loud: absent, duplicated, or crossed markers, an
// empty block, and a missing fence line each name their own drift. A silent
// extractor grades nothing forever — or grades the wrong block — while reporting
// green, which is the failure the marked region exists to make impossible.
func checkExampleAgreement(root string) []string {
	block, diags := ticketExampleBlock(readIfExists(filepath.Join(root, ticketExampleDoc)))
	if len(diags) != 0 {
		return diags
	}
	ticket, diags := parseTicketExample(block)
	if len(diags) != 0 {
		return diags
	}
	return gradeTicketExample(ticket, block)
}

// ticketExampleBlock returns the example's ticket body: the marked region with
// its opening and closing fence lines stripped, materialized exactly as found.
// Nothing is appended — a terminating newline added here would normalize away the
// end-of-file boundary hand-edited markdown reaches, so the check would grade a
// shape the tree never holds.
func ticketExampleBlock(doc string) (string, []string) {
	if diags := ticketExampleMarkerDiags(doc); len(diags) != 0 {
		return "", diags
	}
	begin, end := strings.Index(doc, ticketExampleBegin), strings.Index(doc, ticketExampleEnd)
	if end < begin {
		return "", []string{fmt.Sprintf("%send marker %s precedes the begin marker in %s, so the pair marks no block", ticketExampleDiag, ticketExampleEnd, ticketExampleDoc)}
	}
	region := strings.TrimSpace(doc[begin+len(ticketExampleBegin) : end])
	if region == "" {
		return "", []string{fmt.Sprintf("%sthe marked block in %s is empty; the markers have to enclose the Good example", ticketExampleDiag, ticketExampleDoc)}
	}
	lines := strings.Split(region, "\n")
	var diags []string
	if strings.TrimSpace(lines[0]) != ticketExampleOpenFence {
		diags = append(diags, fmt.Sprintf("%sthe marked block does not open with the %s fence line; the begin marker belongs above it", ticketExampleDiag, ticketExampleOpenFence))
	}
	if len(lines) < 3 || strings.TrimSpace(lines[len(lines)-1]) != ticketExampleCloseFence {
		diags = append(diags, fmt.Sprintf("%sthe marked block does not close with the %s fence line; the end marker belongs below it", ticketExampleDiag, ticketExampleCloseFence))
	}
	if len(diags) != 0 {
		return "", diags
	}
	return strings.Join(lines[1:len(lines)-1], "\n"), nil
}

func ticketExampleMarkerDiags(doc string) []string {
	var diags []string
	for _, marker := range []struct{ name, text string }{{"begin", ticketExampleBegin}, {"end", ticketExampleEnd}} {
		switch count := strings.Count(doc, marker.text); {
		case count == 0:
			diags = append(diags, fmt.Sprintf("%s%s marker %s is absent from %s, so the Good example is graded by nothing", ticketExampleDiag, marker.name, marker.text, ticketExampleDoc))
		case count > 1:
			diags = append(diags, fmt.Sprintf("%s%s marker %s appears %d times in %s; exactly one pair may mark the Good example", ticketExampleDiag, marker.name, marker.text, count, ticketExampleDoc))
		}
	}
	return diags
}

// parseTicketExample runs the assignment parser over the block, laid out in the
// tickets/ directory the parser resolves against a spec path. It writes only its
// own temp directory and reports a setup failure as a diagnostic rather than
// failing the process.
func parseTicketExample(block string) (specbuild.Ticket, []string) {
	tmp, err := os.MkdirTemp("", "bench-ticket-example-*")
	if err != nil {
		return specbuild.Ticket{}, []string{ticketExampleDiag + "setup failed: " + err.Error()}
	}
	defer os.RemoveAll(tmp)
	spec := filepath.Join(tmp, "specs", "example", "spec.md")
	if err := os.MkdirAll(filepath.Join(tmp, "specs", "example", "tickets"), 0o755); err != nil {
		return specbuild.Ticket{}, []string{ticketExampleDiag + "setup failed: " + err.Error()}
	}
	if err := os.WriteFile(spec, []byte("# example agreement\n"), 0o644); err != nil {
		return specbuild.Ticket{}, []string{ticketExampleDiag + "setup failed: " + err.Error()}
	}
	if err := os.WriteFile(filepath.Join(tmp, "specs", "example", "tickets", "example.md"), []byte(block), 0o644); err != nil {
		return specbuild.Ticket{}, []string{ticketExampleDiag + "setup failed: " + err.Error()}
	}
	ticket, err := specbuild.ParseTicket(spec, "example.md")
	if err != nil {
		return specbuild.Ticket{}, []string{fmt.Sprintf("%sthe marked example does not parse as a ticket: %v", ticketExampleDiag, err)}
	}
	return ticket, nil
}

func gradeTicketExample(ticket specbuild.Ticket, block string) []string {
	var diags []string
	if len(ticket.Rows) < 2 {
		diags = append(diags, fmt.Sprintf("%sthe marked example parses to %d distinct acceptance ID(s) %q; the taught `- [ ] [ID] <behavior>` grammar has to yield at least two", ticketExampleDiag, len(ticket.Rows), ticket.Rows))
	}
	diags = append(diags, ticketExampleFieldDiag("acceptance row IDs", ticket.Rows, taughtExampleRows)...)
	diags = append(diags, ticketExampleFieldDiag("ownership fence entries", ticket.Fence, taughtExampleFence)...)
	diags = append(diags, ticketExampleFieldDiag("per-row covers annotations", ticket.Covers, taughtExampleCovers)...)
	diags = append(diags, ticketExampleFieldDiag("atomic closure facts", ticket.Closure, taughtExampleClosure)...)
	diags = append(diags, ticketExampleFieldDiag("red-mutation criteria", ticket.Mutations, taughtExampleClosure)...)
	if !ticketExampleBlockedByLine.MatchString(block) {
		diags = append(diags, fmt.Sprintf("%sthe marked example carries no `Blocked by:` line with a value, so the taught blocker field is demonstrated by nothing", ticketExampleDiag))
	}
	return append(diags, ticketExampleMutationDiags(block, ticket.Closure)...)
}

func ticketExampleFieldDiag(field string, parsed, taught []string) []string {
	if slices.Equal(parsed, taught) {
		return nil
	}
	return []string{fmt.Sprintf("%sthe marked example parses to %s %q, but the taught contract states %q", ticketExampleDiag, field, parsed, taught)}
}

// ticketExampleCriterionCellRe matches a mutations-table row whose criterion cell
// — the row's first cell, immediately after the leading `|` — is the given
// atomic closure fact, mirroring the passlist token regex's first-cell anchoring.
// Anchoring to the cell (rather than looking for the ID anywhere in the section)
// keeps one row's prose from standing in for another row: an operation sequence
// that names RC1 while RC1's own row is gone must not read as coverage for RC1.
func ticketExampleCriterionCellRe(id string) *regexp.Regexp {
	return regexp.MustCompile(`(?m)^\|\s*` + regexp.QuoteMeta(id) + `\s*\|`)
}

func ticketExampleMutationDiags(block string, facts []string) []string {
	heading := ticketExampleMutationsHeading.FindStringIndex(block)
	if heading == nil {
		return []string{fmt.Sprintf("%sthe marked example carries no Red mutations section, so its acceptance rows bind no mutation", ticketExampleDiag)}
	}
	section := block[heading[1]:]
	var diags []string
	for _, fact := range facts {
		if !ticketExampleCriterionCellRe(fact).MatchString(section) {
			diags = append(diags, fmt.Sprintf("%sthe marked example's Red mutations section names no mutation for Closure fact %q", ticketExampleDiag, fact))
		}
	}
	return diags
}

// taughtExampleSkill reads the live craft-tickets skill: the bite proofs mutate
// the real document rather than a second hand-written copy of the example, so a
// bite cannot go on proving a red against a shape the skill no longer teaches.
func taughtExampleSkill(t *testing.T) string {
	t.Helper()
	doc := readIfExists(filepath.Join(NewHarness(t).KitRoot, ticketExampleDoc))
	if doc == "" {
		t.Fatalf("read %s: empty or absent", ticketExampleDoc)
	}
	return doc
}

func ticketExampleRoot(t *testing.T, doc string) string {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, ticketExampleDoc)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// replaceOnce rewrites an anchor that must occur exactly once, so a bite whose
// anchor has drifted fails as a broken bite instead of silently mutating nothing.
func replaceOnce(t *testing.T, doc, anchor, replacement string) string {
	t.Helper()
	if count := strings.Count(doc, anchor); count != 1 {
		t.Fatalf("bite anchor %q occurs %d times, want exactly 1", anchor, count)
	}
	return strings.Replace(doc, anchor, replacement, 1)
}

// rewriteExampleRegion replaces the text between the markers, which is how the
// marker-shaped malformations are built without hand-writing a second example.
func rewriteExampleRegion(t *testing.T, doc string, rewrite func(region string) string) string {
	t.Helper()
	begin, end := strings.Index(doc, ticketExampleBegin), strings.Index(doc, ticketExampleEnd)
	if begin < 0 || end < begin {
		t.Fatalf("marked region not found in the taught skill")
	}
	region := strings.TrimSpace(doc[begin+len(ticketExampleBegin) : end])
	return doc[:begin+len(ticketExampleBegin)] + "\n" + rewrite(region) + "\n" + doc[end:]
}

// TestExampleAgreementParsesAuthoredLiterals is the recorded bite proof for the
// agreement half of checkExampleAgreement (per craft-gate): the live example
// grades clean, and each field the check pins — the fence, the mutations table's
// per-ID coverage, the blocked-by line — reds when the example alone changes.
func TestExampleAgreementParsesAuthoredLiterals(t *testing.T) {
	doc := taughtExampleSkill(t)

	if diags := checkExampleAgreement(ticketExampleRoot(t, doc)); len(diags) != 0 {
		t.Fatalf("taught example: want no diagnostics, got %v", diags)
	}

	drifted := replaceOnce(t, doc, "`internal/status`, `internal/render/rows.go`", "`internal/status`")
	diags := checkExampleAgreement(ticketExampleRoot(t, drifted))
	if !containsDiagnostic(diags, "ownership fence entries") {
		t.Fatalf("dropped fence entry: want the fence mismatch diagnostic, got %v", diags)
	}

	stripped := replaceOnce(t, doc, "(covers CJ1) ", "")
	diags = checkExampleAgreement(ticketExampleRoot(t, stripped))
	if !containsDiagnostic(diags, "per-row covers annotations") {
		t.Fatalf("stripped covers annotation: want the covers mismatch diagnostic, got %v", diags)
	}

	unbound := replaceOnce(t, doc, "\n| RC2/recovery-action | return no recovery action for a cancelled record", "\n")
	diags = checkExampleAgreement(ticketExampleRoot(t, unbound))
	if !containsDiagnostic(diags, `names no mutation for Closure fact "RC2/recovery-action"`) {
		t.Fatalf("dropped mutations row: want the unbound-fact diagnostic, got %v", diags)
	}

	unclosed := replaceOnce(t, doc, "Closure: RC1/reason, RC2/recovery-action\n", "")
	diags = checkExampleAgreement(ticketExampleRoot(t, unclosed))
	if !containsDiagnostic(diags, "atomic closure facts") {
		t.Fatalf("dropped Closure field: want the closure mismatch diagnostic, got %v", diags)
	}

	unblocked := replaceOnce(t, doc, "Blocked by: parse-cancelled-job-records.md\n", "")
	diags = checkExampleAgreement(ticketExampleRoot(t, unblocked))
	if !containsDiagnostic(diags, "carries no `Blocked by:` line") {
		t.Fatalf("dropped blocked-by line: want the blocked-by diagnostic, got %v", diags)
	}

	// The template block carries its own mutations heading, so this one is renamed
	// inside the marked region rather than by a document-wide anchor.
	unmutated := rewriteExampleRegion(t, doc, func(region string) string {
		return replaceOnce(t, region, "## Red mutations", "## Notes")
	})
	diags = checkExampleAgreement(ticketExampleRoot(t, unmutated))
	if !containsDiagnostic(diags, "carries no Red mutations section") {
		t.Fatalf("dropped mutations section: want the missing-section diagnostic, got %v", diags)
	}
}

// TestExampleAgreementMarkersFailClosed is the recorded bite proof for the
// extraction half (per craft-gate): every marker and fence malformation fails
// loudly with its own diagnostic, so no malformation can grade nothing in silence.
func TestExampleAgreementMarkersFailClosed(t *testing.T) {
	doc := taughtExampleSkill(t)

	for _, malformation := range []struct {
		name string
		doc  string
		want string
	}{
		{"absent begin", strings.ReplaceAll(doc, ticketExampleBegin, ""), "begin marker " + ticketExampleBegin + " is absent"},
		{"absent end", strings.ReplaceAll(doc, ticketExampleEnd, ""), "end marker " + ticketExampleEnd + " is absent"},
		{"duplicated begin", doc + "\n" + ticketExampleBegin + "\n", "begin marker " + ticketExampleBegin + " appears 2 times"},
		{"duplicated end", doc + "\n" + ticketExampleEnd + "\n", "end marker " + ticketExampleEnd + " appears 2 times"},
		{"crossed markers", crossExampleMarkers(doc), "precedes the begin marker"},
		{"empty block", rewriteExampleRegion(t, doc, func(string) string { return "" }), "the marked block in " + ticketExampleDoc + " is empty"},
		{"absent opening fence", rewriteExampleRegion(t, doc, dropFirstLine), "does not open with the " + ticketExampleOpenFence + " fence line"},
		{"absent closing fence", rewriteExampleRegion(t, doc, dropLastLine), "does not close with the " + ticketExampleCloseFence + " fence line"},
	} {
		t.Run(malformation.name, func(t *testing.T) {
			diags := checkExampleAgreement(ticketExampleRoot(t, malformation.doc))
			if !containsDiagnostic(diags, malformation.want) {
				t.Fatalf("want diagnostic %q, got %v", malformation.want, diags)
			}
		})
	}
}

// crossExampleMarkers swaps the two marker texts, leaving one of each in the
// document with the end marker first.
func crossExampleMarkers(doc string) string {
	const placeholder = "<!-- ticket-example:crossed -->"
	doc = strings.ReplaceAll(doc, ticketExampleBegin, placeholder)
	doc = strings.ReplaceAll(doc, ticketExampleEnd, ticketExampleBegin)
	return strings.ReplaceAll(doc, placeholder, ticketExampleEnd)
}

func dropFirstLine(region string) string {
	_, rest, _ := strings.Cut(region, "\n")
	return rest
}

func dropLastLine(region string) string {
	if cut := strings.LastIndex(region, "\n"); cut >= 0 {
		return region[:cut]
	}
	return ""
}

// TestExampleAgreementRejectsWrappedFields is the recorded bite proof for the
// corpus failure the parser cannot report (per craft-gate): a field wrapped onto
// a second line, or a row written without its ID, is read as a shorter parse
// rather than refused, and only the authored literals turn that red.
func TestExampleAgreementRejectsWrappedFields(t *testing.T) {
	doc := taughtExampleSkill(t)

	wrappedFence := replaceOnce(t, doc, "`internal/status`, `internal/render/rows.go`", "`internal/status`,\n`internal/render/rows.go`")
	diags := checkExampleAgreement(ticketExampleRoot(t, wrappedFence))
	if !containsDiagnostic(diags, "ownership fence entries") {
		t.Fatalf("wrapped fence: want the fence mismatch diagnostic, got %v", diags)
	}

	unlabeled := replaceOnce(t, doc, "- [ ] [RC2] (covers CJ2) status renders", "- [ ] status renders")
	diags = checkExampleAgreement(ticketExampleRoot(t, unlabeled))
	if !containsDiagnostic(diags, "distinct acceptance ID(s)") {
		t.Fatalf("unlabeled row: want the too-few-IDs diagnostic, got %v", diags)
	}
}

// TestExampleAgreementEOFWithoutNewline is the recorded bite proof for the
// boundary hand-edited markdown reaches (per craft-gate): the block the check
// materializes carries the extracted bytes as found, so an example ending at end
// of file reaches the parser with its final line unterminated. Appending a newline
// before materializing would normalize that boundary away and grade a shape the
// tree never holds; the parse is compared against the newline-terminated form, so
// a last-line field cannot vanish from the parse while the check stays green.
func TestExampleAgreementEOFWithoutNewline(t *testing.T) {
	doc := taughtExampleSkill(t)
	end := strings.Index(doc, ticketExampleEnd)
	if end < 0 {
		t.Fatalf("marked region not found in the taught skill")
	}
	atEOF := strings.TrimRight(doc[:end+len(ticketExampleEnd)], "\n")

	eofBlock, diags := ticketExampleBlock(atEOF)
	if len(diags) != 0 {
		t.Fatalf("end-of-file block: %v", diags)
	}
	if strings.HasSuffix(eofBlock, "\n") {
		t.Fatalf("materialized end-of-file block ends with a newline; the extracted bytes were normalized rather than materialized as found")
	}

	first, diags := parseTicketExample(eofBlock + "\n")
	if len(diags) != 0 {
		t.Fatalf("parse newline-terminated block: %v", diags)
	}
	second, diags := parseTicketExample(eofBlock)
	if len(diags) != 0 {
		t.Fatalf("parse end-of-file block: %v", diags)
	}
	if first.Title != second.Title {
		t.Fatalf("title at end of file = %q, want %q", second.Title, first.Title)
	}
	for _, field := range []struct {
		name             string
		parsed, expected []string
	}{
		{"rows", second.Rows, first.Rows},
		{"fence", second.Fence, first.Fence},
	} {
		if !slices.Equal(field.parsed, field.expected) {
			t.Fatalf("%s at end of file = %q, want %q", field.name, field.parsed, field.expected)
		}
	}
	if diags := checkExampleAgreement(ticketExampleRoot(t, atEOF)); len(diags) != 0 {
		t.Fatalf("end-of-file document: want no diagnostics, got %v", diags)
	}
}

// TestExampleAgreementFactsAnchorToCriterionCell is the recorded bite proof for
// the anchoring (per craft-gate): a fact mentioned in another row's prose does
// not stand in for that fact's own mutations row.
func TestExampleAgreementFactsAnchorToCriterionCell(t *testing.T) {
	doc := taughtExampleSkill(t)

	stray := replaceOnce(t, doc, "expect the missing-action failure", "expect the missing-action failure, the same sequence RC1/reason uses")
	deleted := replaceOnce(t, stray, "\n| RC1/reason | render the cancelled row with an empty reason | the cancelled-row render test | blank the field, run `go test ./internal/render`, expect the missing-reason failure |", "")

	if !strings.Contains(deleted, "RC1/reason") {
		t.Fatalf("bite setup: the mutated document no longer mentions RC1/reason at all")
	}
	diags := checkExampleAgreement(ticketExampleRoot(t, deleted))
	if !containsDiagnostic(diags, `names no mutation for Closure fact "RC1/reason"`) {
		t.Fatalf("deleted fact row with a stray fact mention: want the unbound-fact diagnostic, got %v", diags)
	}
}

// TestExampleAgreementFixtureBite is the canary-fixture bite proof: the fixture
// wraps the example's ownership fence onto a second line, and the check fires the
// fence mismatch under a full RunConformance pass.
func TestExampleAgreementFixtureBite(t *testing.T) {
	h := NewHarness(t)
	kitRoot := h.KitRoot
	fixture := "wrapped-example-fence"
	runFixtureBite(t, kitRoot, fixture)
}
