package learnings

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestCommandAppendsTypedDrainActionsAndHonestEmptyHelp(t *testing.T) {
	for _, tc := range []struct {
		name    string
		journal string
		want    string
		code    int
	}{
		{name: "open", journal: "# Learnings — usage journal\n\n## 2026-01-01 — first [open]\n\n## 2026-01-02 — second [open]\n\n## 2026-01-01 — first [open]\n", want: "learnings[3]{date,title}:\n  2026-01-01,first\n  2026-01-02,second\n  2026-01-01,first\nhelp[2]{cmd,why}:\n  /bench-what-next,\"verdict 2026-01-01: first\"\n  /bench-what-next,\"verdict 2026-01-02: second\"\n", code: 0},
		{name: "drained", journal: "# Learnings — usage journal\n", want: "learnings[0]{date,title}:\nhelp[0]{cmd,why}:\n", code: 0},
		{name: "absent", want: "learnings[0]{date,title}:\nhelp[0]{cmd,why}:\n", code: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := learningsRepo(t)
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
			if strings.Contains(got, "bench /bench-what-next") {
				t.Fatal("harness phase rendered as shell command")
			}
		})
	}
}

func TestCommandPreservesCheckedInPreDisclosureResponses(t *testing.T) {
	for _, tc := range []struct {
		name, fixture, journal, help string
	}{
		{"open", "pre-disclosure-open.stdout", "# Learnings — usage journal\n\n## 2026-01-01 — first [open]\n", "help[1]{cmd,why}:\n  /bench-what-next,\"verdict 2026-01-01: first\"\n"},
		{"drained", "pre-disclosure-drained.stdout", "# Learnings — usage journal\n", "help[0]{cmd,why}:\n"},
		{"absent", "pre-disclosure-absent.stdout", "", "help[0]{cmd,why}:\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			primary, err := os.ReadFile(filepath.Join("testdata", tc.fixture))
			if err != nil {
				t.Fatal(err)
			}
			root := learningsRepo(t)
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
		code                          int
	}
	cases := []fixtureCase{
		{name: "malformed-among-parsed", pre: "pre-disclosure-malformed.stdout", candidate: "candidate-malformed.stdout", journal: "# Learnings — usage journal\n\n## 2026-01-01 — first [open]\n## broken\n", code: 1},
		{name: "unsupported-schema", pre: "pre-disclosure-unsupported.stdout", candidate: "candidate-unsupported.stdout", journal: "not a learnings journal\n", code: 1},
		{name: "empty", pre: "pre-disclosure-empty.stdout", candidate: "candidate-empty.stdout", setup: func(t *testing.T, root string) { writeJournal(t, root, nil) }, code: 1},
		{name: "invalid-utf8", pre: "pre-disclosure-invalid-utf8.stdout", candidate: "candidate-invalid-utf8.stdout", setup: func(t *testing.T, root string) { writeJournal(t, root, []byte{0xff, '\n'}) }, code: 1},
		{name: "oversized", pre: "pre-disclosure-oversized.stdout", candidate: "candidate-oversized.stdout", setup: func(t *testing.T, root string) { writeJournal(t, root, make([]byte, bounds.ControlRecordLimit+1)) }, code: 1},
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
			root := learningsRepo(t)
			if tc.setup != nil {
				tc.setup(t, root)
			} else {
				writeJournal(t, root, []byte(tc.journal))
			}
			t.Chdir(root)
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

func writeJournal(t *testing.T, root string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, JournalPath), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func learningsRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	return root
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
