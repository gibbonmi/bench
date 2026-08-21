package learnings

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestCommandAppendsTypedDrainActionsAndHonestEmptyHelp(t *testing.T) {
	for _, tc := range []struct {
		name    string
		journal string
		want    string
		code    int
	}{
		{name: "open", journal: "# Learnings — usage journal\n\n## 2026-01-01 — first [open]\n\n## 2026-01-02 — second [open]\n\n## 2026-01-01 — first [open]\n", want: "learnings[3]{date,title}:\n  2026-01-01,first\n  2026-01-02,second\n  2026-01-01,first\nhelp[2]{cmd,why}:\n  /bench-drain,\"verdict 2026-01-01: first\"\n  /bench-drain,\"verdict 2026-01-02: second\"\n", code: 0},
		{name: "drained", journal: "# Learnings — usage journal\n", want: "learnings[0]{date,title}:\nhelp[0]{cmd,why}:\n", code: 0},
		{name: "absent", want: "learnings[0]{date,title}:\nhelp[0]{cmd,why}:\n", code: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := gittest.Repo(t)
			if tc.journal != "" {
				if err := os.MkdirAll(filepath.Join(root, "capture"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, JournalPath), []byte(tc.journal), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			t.Chdir(root)
			got, code := Command(nil)
			if code != tc.code || got != tc.want {
				t.Fatalf("Command = (%d, %q), want (%d, %q)", code, got, tc.code, tc.want)
			}
			if strings.Contains(got, "bench /bench-drain") {
				t.Fatal("harness phase rendered as shell command")
			}
		})
	}
}

func TestCommandPreservesCheckedInPreDisclosureResponses(t *testing.T) {
	for _, tc := range []struct {
		name, fixture, journal, help string
	}{
		{"open", "pre-disclosure-open.stdout", "# Learnings — usage journal\n\n## 2026-01-01 — first [open]\n", "help[1]{cmd,why}:\n  /bench-drain,\"verdict 2026-01-01: first\"\n"},
		{"drained", "pre-disclosure-drained.stdout", "# Learnings — usage journal\n", "help[0]{cmd,why}:\n"},
		{"absent", "pre-disclosure-absent.stdout", "", "help[0]{cmd,why}:\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			primary, err := os.ReadFile(filepath.Join("testdata", tc.fixture))
			if err != nil {
				t.Fatal(err)
			}
			root := gittest.Repo(t)
			if tc.journal != "" {
				if err := os.MkdirAll(filepath.Join(root, "capture"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, JournalPath), []byte(tc.journal), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			t.Chdir(root)
			got, code := Command(nil)
			if code != 0 || got != string(primary)+tc.help {
				t.Fatalf("Command = (%d, %q), want checked-in primary plus exactly one help block", code, got)
			}
		})
	}
}

func TestRefusalEvidencePreservesPrimaryAndDisclosesOnlyRepairableMalformed(t *testing.T) {
	type fixtureCase struct {
		name, pre, candidate, journal string
		setup                         func(*testing.T, string)
		// resolve supplies the run-scoped bytes a fixture cannot check in — an OS error
		// naming an absolute temp path — for substitution into journalPathToken.
		resolve func(*testing.T) string
		code    int
	}
	cases := []fixtureCase{
		{name: "malformed-among-parsed", pre: "pre-disclosure-malformed.stdout", candidate: "candidate-malformed.stdout", journal: "# Learnings — usage journal\n\n## 2026-01-01 — first [open]\n## broken\n", code: 1},
		{name: "unsupported-schema", pre: "pre-disclosure-unsupported.stdout", candidate: "candidate-unsupported.stdout", journal: "not a learnings journal\n", code: 1},
		{name: "empty", pre: "pre-disclosure-empty.stdout", candidate: "candidate-empty.stdout", setup: func(t *testing.T, root string) { writeJournal(t, root, nil) }, code: 1},
		{name: "invalid-utf8", pre: "pre-disclosure-invalid-utf8.stdout", candidate: "candidate-invalid-utf8.stdout", setup: func(t *testing.T, root string) { writeJournal(t, root, []byte{0xff, '\n'}) }, code: 1},
		{name: "oversized", pre: "pre-disclosure-oversized.stdout", candidate: "candidate-oversized.stdout", setup: func(t *testing.T, root string) { writeJournal(t, root, make([]byte, bounds.ControlRecordLimit+1)) }, code: 1},
		{name: "dangling-symlink", pre: "pre-disclosure-dangling-symlink.stdout", candidate: "candidate-dangling-symlink.stdout", setup: func(t *testing.T, root string) {
			if err := os.MkdirAll(filepath.Join(root, "capture"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(filepath.Join(root, "capture", "vanished.md"), filepath.Join(root, JournalPath)); err != nil {
				capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
			}
		}, resolve: journalAbsPath, code: 1},
		{name: "wrong-type", pre: "pre-disclosure-wrong-type.stdout", candidate: "candidate-wrong-type.stdout", setup: func(t *testing.T, root string) {
			if err := os.MkdirAll(filepath.Join(root, JournalPath), 0o755); err != nil {
				t.Fatal(err)
			}
		}, code: 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pre, err := os.ReadFile(filepath.Join("testdata", tc.pre))
			if err != nil {
				t.Fatal(err)
			}
			candidate, err := os.ReadFile(filepath.Join("testdata", tc.candidate))
			if err != nil {
				t.Fatal(err)
			}
			root := gittest.Repo(t)
			if tc.setup != nil {
				tc.setup(t, root)
			} else {
				writeJournal(t, root, []byte(tc.journal))
			}
			t.Chdir(root)
			if tc.resolve != nil {
				runScoped := tc.resolve(t)
				pre = []byte(strings.ReplaceAll(string(pre), journalPathToken, runScoped))
				candidate = []byte(strings.ReplaceAll(string(candidate), journalPathToken, runScoped))
			}
			got, code := Command(nil)
			t.Logf("observed candidate code=%d stdout=%q", code, got)
			if code != tc.code || got != string(candidate) {
				t.Fatalf("candidate = (%d, %q), want (%d, %q)", code, got, tc.code, candidate)
			}
			if tc.name != "malformed-among-parsed" && string(pre) != string(candidate) {
				t.Fatalf("early refusal changed: pre=%q candidate=%q", pre, candidate)
			}
		})
	}
}

// journalPathToken stands in the checked-in fixtures for the absolute journal path an
// OS error string embeds, which is temp-directory-scoped and so differs every run.
const journalPathToken = "{{JOURNAL}}"

// journalAbsPath is the absolute path the command classified, resolved the way the
// command resolves it so a filesystem that reports a different real path than the one
// t.TempDir handed out still substitutes the bytes the error actually names.
func journalAbsPath(t *testing.T) string {
	t.Helper()
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(root, JournalPath)
}

func writeJournal(t *testing.T, root string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, JournalPath), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRows(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want [][]string
	}{
		{"em-dash separator", "## 2026-01-01 — first learning  [open]\n- body\n", [][]string{{"2026-01-01", "first learning"}}},
		{"ascii-hyphen separator", "## 2026-01-01 - ascii title  [open]\n", [][]string{{"2026-01-01", "ascii title"}}},
		{"comma/quote title", `## 2026-03-03 — a, "b"  [open]`, [][]string{{"2026-03-03", `a, "b"`}}},
		{"CRLF heading", "## 2026-04-04 — crlf  [open]\r\n", [][]string{{"2026-04-04", "crlf"}}},
		{"no trailing newline", "## 2026-05-05 — tail  [open]", [][]string{{"2026-05-05", "tail"}}},
		{"template example is not an entry", "## <date> — <short title>  [open]\n", nil},
		{"resolved heading ignored", "## 2026-01-01 — done  [resolved]\n", nil},
		{"empty", "", nil},
	}
	for _, c := range cases {
		if got := Rows([]byte(c.in)); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: Rows = %v, want %v", c.name, got, c.want)
		}
	}
}

// TestParseReportsDatedLineThatMissesHeadingShape covers DL1, DL2, DL3, DL6, DL32,
// DL7, DL8, DL9, DL10, DL11, DL12, DL20, and DL21: a line that leads with a date but
// is not a well-formed dated heading becomes its own malformed record, while the two
// pre-existing `## ` dispositions keep exactly their current single record.
func TestParseReportsDatedLineThatMissesHeadingShape(t *testing.T) {
	const schema = JournalSchemaHeading + "\n\n"
	for _, tc := range []struct {
		name, in string
		want     []Malformed
	}{
		{"DL1,DL2,DL3 bullet under the schema heading", schema + "- 2026-08-21 — first thing\n",
			[]Malformed{{Reason: "dated learning entry is not a heading", Raw: "- 2026-08-21 — first thing", Line: 3}}},
		{"DL6 every markdown marker", "- 2026-01-01 a\n* 2026-01-02 b\n+ 2026-01-03 c\n> 2026-01-04 d\n# 2026-01-05 e\n",
			[]Malformed{
				{Reason: "dated learning entry is not a heading", Raw: "- 2026-01-01 a", Line: 1},
				{Reason: "dated learning entry is not a heading", Raw: "* 2026-01-02 b", Line: 2},
				{Reason: "dated learning entry is not a heading", Raw: "+ 2026-01-03 c", Line: 3},
				{Reason: "dated learning entry is not a heading", Raw: "> 2026-01-04 d", Line: 4},
				{Reason: "dated learning entry is not a heading", Raw: "# 2026-01-05 e", Line: 5},
			}},
		{"DL32 flush at column one", "2026-08-21 — flush\n",
			[]Malformed{{Reason: "dated learning entry is not a heading", Raw: "2026-08-21 — flush", Line: 1}}},
		{"DL7 no-break space separator", "- 2026-08-21 — nbsp\n",
			[]Malformed{{Reason: "dated learning entry is not a heading", Raw: "- 2026-08-21 — nbsp", Line: 1}}},
		{"DL8 ideographic space separator", "-　2026-08-21 — ideographic\n",
			[]Malformed{{Reason: "dated learning entry is not a heading", Raw: "-　2026-08-21 — ideographic", Line: 1}}},
		{"DL8 zero-width space is not a separator", "-​2026-08-21 — zero width\n", nil},
		{"DL9 dated bullet inside an open entry's body", schema + "## 2026-01-01 — first [open]\n- body\n- 2026-08-21 — appended\n",
			[]Malformed{{Reason: "dated learning entry is not a heading", Raw: "- 2026-08-21 — appended", Line: 5}}},
		{"DL10 broken heading keeps its one record", "## broken\n",
			[]Malformed{{Reason: "malformed learning heading", Raw: "## broken", Line: 1}}},
		{"DL11 dated heading without [open] keeps its one record", "## 2026-01-01 — x\n",
			[]Malformed{{Reason: "dated learning heading must end with [open]", Raw: "## 2026-01-01 — x", Line: 1}}},
		{"DL12 digit-shaped non-calendar date", "- 2026-88-88 — x\n",
			[]Malformed{{Reason: "dated learning entry is not a heading", Raw: "- 2026-88-88 — x", Line: 1}}},
		{"DL20 final line with no trailing newline", schema + "- 2026-08-21 — tail",
			[]Malformed{{Reason: "dated learning entry is not a heading", Raw: "- 2026-08-21 — tail", Line: 3}}},
		{"DL21 CRLF line loses its carriage return", "- 2026-08-21 — crlf\r\n",
			[]Malformed{{Reason: "dated learning entry is not a heading", Raw: "- 2026-08-21 — crlf", Line: 1}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, got := Parse([]byte(tc.in))
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Parse malformed = %#v, want %#v", got, tc.want)
			}
		})
	}
}

// TestCommandSurfacesLostDatedLinesAndLeavesTheQuietPosturesAlone covers DL4, DL5,
// DL15, DL16, and DL17: the journal shape that produced the 2026-08-21 drop exits 1
// with one `line <n>` row per lost entry, byte-exact against a checked-in fixture,
// while the freshly scaffolded, drained, and well-formed journals stay green.
func TestCommandSurfacesLostDatedLinesAndLeavesTheQuietPosturesAlone(t *testing.T) {
	scaffold, err := os.ReadFile(filepath.Join("testdata", "scaffold-learnings.md"))
	if err != nil {
		t.Fatal(err)
	}
	lostRows, err := os.ReadFile(filepath.Join("testdata", "candidate-dated-lines.stdout"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, journal, want string
		code                int
	}{
		{name: "DL4,DL5 the two-bullet journal", journal: JournalSchemaHeading + "\n\n- 2026-08-21 — spec anchor drift\n- 2026-08-21 — worktree tip mismatch\n", want: string(lostRows), code: 1},
		{name: "DL15 freshly scaffolded journal", journal: string(scaffold), want: "learnings[0]{date,title}:\nhelp[0]{cmd,why}:\n", code: 0},
		{name: "DL16 schema heading only", journal: JournalSchemaHeading + "\n", want: "learnings[0]{date,title}:\nhelp[0]{cmd,why}:\n", code: 0},
		{name: "DL17 well-formed open entry", journal: JournalSchemaHeading + "\n\n## 2026-01-01 — first [open]\n", want: "learnings[1]{date,title}:\n  2026-01-01,first\nhelp[1]{cmd,why}:\n  /bench-drain,\"verdict 2026-01-01: first\"\n", code: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := gittest.Repo(t)
			writeJournal(t, root, []byte(tc.journal))
			t.Chdir(root)
			got, code := Command(nil)
			if code != tc.code || got != tc.want {
				t.Fatalf("Command = (%d, %q), want (%d, %q)", code, got, tc.code, tc.want)
			}
		})
	}
}
