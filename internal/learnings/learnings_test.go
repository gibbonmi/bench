package learnings

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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
