package coverage

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/toon"
	toonlib "github.com/toon-format/toon-go"
)

const stories = "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n"
const hdr = "| story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|\n"
const hdr6 = "| row | story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|---|\n"

func spec(body string) parsed { return parse([]byte(body)) }

func TestStateAndRows(t *testing.T) {
	p := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr +
		"| 2–3 | does x \\| y | cli seam | cmd fails, loudly | catches z |\n" +
		"| edge of 1 | edge case | gate | already covered | catches w |\n")
	if State(p) != "mapped" {
		t.Fatalf("state = %q", State(p))
	}
	want := [][]string{{"2–3", "cli seam", "cmd fails, loudly"}, {"edge of 1", "gate", "already covered"}}
	if got := Rows(p); !reflect.DeepEqual(got, want) {
		t.Errorf("Rows = %v, want %v", got, want)
	}

	if State(spec("# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")) != "historical" {
		t.Error("historical state not detected")
	}
	if State(spec("# n\nprose only\n")) != "no-map" {
		t.Error("no-map state not detected")
	}
}

// Every validation phrasing is matched by substring downstream; pin each here.
func TestCheck(t *testing.T) {
	valid := spec("# v\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1, 2–3 | b | s | r | w |\n")
	if v := Check(valid); len(v) != 0 {
		t.Errorf("valid map violations = %v", v)
	}
	// Historical opts out of validation.
	if v := Check(spec("# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")); v != nil {
		t.Errorf("historical Check = %v, want nil", v)
	}
	cases := []struct{ body, want string }{
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n", "coverage map missing the canonical header"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr, "coverage map has no data rows"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | b | s | r |\n", "coverage map row 1 has 4 cells (want 5)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | b |  | r | w |\n", "coverage map row 1 has an empty 'seam' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 9 | b | s | r | w |\n", "references story 9, which the spec does not declare (has: 1, 2, 3)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| x | b | s | r | w |\n", "has an unrecognized story reference 'x'"},
		// A 6-cell map is opt-in; its rows carry a leading row-ID cell that Check
		// validates in addition to the legacy fields.
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| AB1 | 1 | b | s | r | w |\n| AB1 | 1 | b | s | r | w |\n",
			"coverage map row 2 has a duplicate row id 'AB1' (first used at row 1)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "|  | 1 | b | s | r | w |\n",
			"coverage map row 1 has an empty 'row' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| ab-1 | 1 | b | s | r | w |\n",
			"coverage map row 1 has a malformed row id 'ab-1'"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| AB1 | 1 | b | s | r | w |\n| not-an-id | 2 | b | s | r | w |\n",
			"coverage map row 2 has a malformed row id 'not-an-id'"},
	}
	for _, c := range cases {
		v := Check(spec(c.body))
		if len(v) == 0 || !strings.Contains(strings.Join(v, "\n"), c.want) {
			t.Errorf("Check violations %v do not contain %q", v, c.want)
		}
	}
}

// TestCheckOptIn covers the 6-cell canonical header on its own: a valid opt-in map
// has no violations, and its Rows/State behave exactly as a legacy map's.
func TestCheckOptIn(t *testing.T) {
	p := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 +
		"| AB1 | 1 | does x | cli seam | cmd fails | catches z |\n" +
		"| CD2 | 2 | does y | gate | already covered | catches w |\n")
	if v := Check(p); len(v) != 0 {
		t.Errorf("valid opt-in map violations = %v, want none", v)
	}
	if State(p) != "mapped" {
		t.Fatalf("state = %q, want mapped", State(p))
	}
	want := [][]string{{"1", "cli seam", "cmd fails"}, {"2", "gate", "already covered"}}
	if got := Rows(p); !reflect.DeepEqual(got, want) {
		t.Errorf("Rows = %v, want %v", got, want)
	}
}

// TestParseSpecOptIn drives ParseSpec, the package's exported entry point, over a
// 6-cell map on disk: the opt-in verdict, the ordered row IDs, and no violations.
func TestParseSpecOptIn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.md")
	body := "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 +
		"| AB1 | 1 | does x | cli seam | cmd fails | catches z |\n" +
		"| CD2 | 2 | does y | gate | already covered | catches w |\n"
	mustWrite(t, path, body)

	optIn, ids, violations, err := ParseSpec(path)
	if err != nil {
		t.Fatalf("ParseSpec: %v", err)
	}
	if !optIn {
		t.Error("optIn = false, want true for a 6-cell map")
	}
	if want := []string{"AB1", "CD2"}; !reflect.DeepEqual(ids, want) {
		t.Errorf("ids = %v, want %v", ids, want)
	}
	if len(violations) != 0 {
		t.Errorf("violations = %v, want none", violations)
	}

	// A legacy 5-cell map reports not opted in, with nil IDs.
	legacyPath := filepath.Join(dir, "legacy.md")
	mustWrite(t, legacyPath, "# t\n\n"+stories+"\n### Acceptance coverage map\n"+hdr+"| 1 | b | s | r | w |\n")
	optIn, ids, _, err = ParseSpec(legacyPath)
	if err != nil {
		t.Fatalf("ParseSpec: %v", err)
	}
	if optIn {
		t.Error("optIn = true, want false for a 5-cell map")
	}
	if ids != nil {
		t.Errorf("ids = %v, want nil for a legacy map", ids)
	}
}

// mustWrite creates path (and any parent dirs) with content under the current CWD.
func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", dir, err)
		}
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

// wantTable renders the exact TOON output Command produces for a resolved spec label and
// body, via the same State/Rows/toon.Table calls Command itself makes — so the test pins
// the round-trip without re-deriving the TOON format.
func wantTable(t *testing.T, label, body string) string {
	t.Helper()
	p := spec(body)
	tbl, err := toon.Table("rows", []string{"story", "seam", "red_signal"}, Rows(p))
	if err != nil {
		t.Fatalf("toon.Table: %v", err)
	}
	actions := make([]axi.Action, 0, len(Rows(p)))
	for _, row := range Rows(p) {
		actions = append(actions, axi.ExecutableInvocation("check coverage row "+row[0], axi.KnownArgument("coverage"), axi.KnownArgument("--check"), axi.KnownArgument(label)))
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		t.Fatal(err)
	}
	return "spec: " + label + "\n" + "state: " + State(p) + "\n" + tbl + help
}

// TestCommand drives Command through its public (args) -> (output, exit code) interface
// only, per the spec's testing decision — never parse/Check directly.
func TestCommand(t *testing.T) {
	mapped := func(row string) string {
		return "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + row
	}

	t.Run("readable path argument resolves and round-trips", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | r | w |\n")
		mustWrite(t, "spec.md", body)

		out, code := Command([]string{"spec.md"})
		if want := wantTable(t, "spec.md", body); out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0)", out, code, want)
		}
	})

	t.Run("appends one check action per unchecked row and deduplicates exact templates", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | unchecked | catches one |\n| 2 | c | t | already covered | catches two |\n| 1 | b | s | unchecked | catches one |\n")
		mustWrite(t, "path with spaces/spec.md", body)
		out, code := Command([]string{"path with spaces/spec.md"})
		want := "spec: path with spaces/spec.md\nstate: mapped\nrows[3]{story,seam,red_signal}:\n  \"1\",s,unchecked\n  \"2\",t,already covered\n  \"1\",s,unchecked\nhelp[2]{cmd,why}:\n  bench coverage --check 'path with spaces/spec.md',check coverage row 1\n  bench coverage --check 'path with spaces/spec.md',check coverage row 2\n"
		if code != 0 || out != want {
			t.Fatalf("Command = (%d, %q), want (%d, %q)", code, out, 0, want)
		}
	})

	t.Run("repair prose does not alter unchecked classification", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | repair is not evidence | catches one |\n")
		mustWrite(t, "spec.md", body)
		out, code := Command([]string{"spec.md"})
		if code != 0 || !strings.Contains(out, "bench coverage --check spec.md,check coverage row 1") {
			t.Fatalf("Command = (%d, %q), want unchecked check action", code, out)
		}
	})

	t.Run("malformed map gets repair retry without changing extraction exit", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 9 | b | s | red | catches one |\n")
		mustWrite(t, "spec.md", body)
		out, code := Command([]string{"spec.md"})
		want := "spec: spec.md\nstate: mapped\nrows[1]{story,seam,red_signal}:\n  \"9\",s,red\nhelp[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"
		if code != 0 || out != want {
			t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
		}
	})

	t.Run("canonical mapped zero-row table is terminal", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("")
		mustWrite(t, "spec.md", body)
		out, code := Command([]string{"spec.md"})
		if want := wantTable(t, "spec.md", body); out != want || code != 0 {
			t.Fatalf("Command = (%d, %q), want (%d, %q)", code, out, 0, want)
		}
	})

	t.Run("separator-free slug resolves specs/<slug>/spec.md", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 2 | b2 | s2 | r2 | w2 |\n")
		mustWrite(t, "specs/foo/spec.md", body)

		out, code := Command([]string{"foo"})
		if want := wantTable(t, "specs/foo/spec.md", body); out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0)", out, code, want)
		}
	})

	t.Run("slug already ending .md resolves folder spec", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 3 | b3 | s3 | r3 | w3 |\n")
		mustWrite(t, "specs/bar/spec.md", body)

		out, code := Command([]string{"bar.md"})
		if want := wantTable(t, "specs/bar/spec.md", body); out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0)", out, code, want)
		}
	})

	t.Run("slug matching no file names both forms tried, exit 1", func(t *testing.T) {
		t.Chdir(t.TempDir())

		out, code := Command([]string{"missing"})
		if code != 1 || !strings.Contains(out, "spec not found: missing, specs/missing/spec.md") {
			t.Errorf("Command = (%q, %d), want exit 1 naming both 'missing' and 'specs/missing/spec.md' in the not-found message", out, code)
		}
	})

	t.Run("separator-bearing argument gets no fallback", func(t *testing.T) {
		t.Chdir(t.TempDir())

		out, code := Command([]string{"sub/missing.md"})
		if code != 1 || !strings.Contains(out, "sub/missing.md") || strings.Contains(out, "specs/") {
			t.Errorf("Command = (%q, %d), want exit 1 naming only sub/missing.md with no specs/ form", out, code)
		}
	})

	t.Run("slug shadowed by a same-named CWD file resolves path-first", func(t *testing.T) {
		t.Chdir(t.TempDir())
		cwdBody := mapped("| 1 | cwd | cwd-seam | cwd-red | cwd-why |\n")
		specsBody := mapped("| 5 | specs | specs-seam | specs-red | specs-why |\n")
		mustWrite(t, "foo", cwdBody)
		mustWrite(t, "specs/foo/spec.md", specsBody)

		out, code := Command([]string{"foo"})
		if want := wantTable(t, "foo", cwdBody); out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0) — the CWD file should shadow the specs/ fallback", out, code, want)
		}
	})

	t.Run("present-but-unreadable file surfaces the read error, not not-found", func(t *testing.T) {
		if os.Geteuid() == 0 {
			capability.Capability(t, capability.Privilege, "root reads 0000-mode files; permission case unobservable")
		}
		t.Chdir(t.TempDir())
		mustWrite(t, "sub/locked.md", mapped("| 1 | b | s | r | w |\n"))
		if err := os.Chmod("sub/locked.md", 0o000); err != nil {
			t.Fatalf("Chmod: %v", err)
		}

		out, code := Command([]string{"sub/locked.md"})
		if code != 1 || !strings.Contains(out, "spec not readable:") || strings.Contains(out, "spec not found") {
			t.Errorf("Command = (%q, %d), want exit 1 with a 'spec not readable:' error, never 'spec not found'", out, code)
		}
	})

	t.Run("slug matching a directory falls back to specs/<slug>/spec.md", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | r | w |\n")
		mustWrite(t, "dist/keep", "x") // a directory named like the slug
		mustWrite(t, "specs/dist/spec.md", body)

		out, code := Command([]string{"dist"})
		if want := wantTable(t, "specs/dist/spec.md", body); out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0) — a directory is not a spec candidate", out, code, want)
		}
	})

	t.Run("--check resolves a slug and reports violations under the resolved label", func(t *testing.T) {
		t.Chdir(t.TempDir())
		mustWrite(t, "specs/chk/spec.md", mapped("| 1 | | s | r | w |\n")) // empty behavior cell

		out, code := Command([]string{"--check", "chk"})
		if code != 1 || !strings.Contains(out, "error: specs/chk/spec.md coverage map row 1 has an empty 'behavior' cell") {
			t.Errorf("Command = (%q, %d), want exit 1 with the violation under the resolved folder label", out, code)
		}
	})

	t.Run("flag-shaped argument stays a usage error, exit 2", func(t *testing.T) {
		t.Chdir(t.TempDir())

		out, code := Command([]string{"--bogus"})
		if want := toon.Usage("bench coverage", "--bogus") + "\n"; out != want || code != 2 {
			t.Errorf("Command = (%q, %d), want (%q, 2)", out, code, want)
		}
	})

	t.Run("no argument reports missing/required, not unknown argument, exit 2", func(t *testing.T) {
		t.Chdir(t.TempDir())

		out, code := Command(nil)
		if code != 2 {
			t.Errorf("Command = (%q, %d), want exit 2", out, code)
		}
		if !strings.Contains(out, "required") && !strings.Contains(out, "missing") {
			t.Errorf("Command = %q, want it to say the argument is missing/required", out)
		}
		if strings.Contains(out, "unknown argument") {
			t.Errorf("Command = %q, must not use the unknown-argument template for a missing argument", out)
		}
	})
}

func TestCommandPreservesCheckedInPreDisclosureResponses(t *testing.T) {
	mapped := func(row string) string {
		return "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + row
	}
	for _, tc := range []struct {
		name, fixture, body, help string
	}{
		{"mapped actionable", "pre-disclosure-mapped.stdout", mapped("| 1 | b | s | unchecked | catches one |\n"), "help[1]{cmd,why}:\n  bench coverage --check spec.md,check coverage row 1\n"},
		{"repairable malformed", "pre-disclosure-malformed.stdout", mapped("| 9 | b | s | red | catches one |\n"), "help[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"},
		{"mapped zero-row", "pre-disclosure-zero-row.stdout", mapped(""), "help[0]{cmd,why}:\n"},
		{"historical terminal", "pre-disclosure-historical.stdout", "# historical\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n", "help[0]{cmd,why}:\n"},
		{"no-map terminal", "pre-disclosure-no-map.stdout", "# no map\n", "help[0]{cmd,why}:\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			primary, err := os.ReadFile(filepath.Join("testdata", tc.fixture))
			if err != nil {
				t.Fatal(err)
			}
			t.Chdir(t.TempDir())
			mustWrite(t, "spec.md", tc.body)
			out, code := Command([]string{"spec.md"})
			if code != 0 || out != string(primary)+tc.help {
				t.Fatalf("Command = (%d, %q), want checked-in primary plus exactly one help block", code, out)
			}
		})
	}
}

func TestCommandControlBearingSpecPathPreservesPrimaryAndHonestFallback(t *testing.T) {
	mapped := "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | b | s | unchecked | catches one |\n"
	for _, control := range []string{"\t", "\n", "\r", "\x1b"} {
		t.Run("control", func(t *testing.T) {
			t.Chdir(t.TempDir())
			path := "control" + control + "spec.md"
			mustWrite(t, path, mapped)
			out, code := Command([]string{path})
			primary := "spec: " + path + "\nstate: mapped\nrows[1]{story,seam,red_signal}:\n  \"1\",s,unchecked\n"
			if code != 0 || !strings.HasPrefix(out, primary) {
				t.Fatalf("Command = (%d, %q), want primary response and exit 0", code, out)
			}
			if control == "\x1b" {
				if out != primary+"help[0]{cmd,why}:\n" {
					t.Fatalf("Command fallback = %q, want primary plus empty help", out)
				}
				return
			}
			if !strings.Contains(out, "help[1]{cmd,why}:") {
				t.Fatalf("Command = %q, want one action", out)
			}
			decoded, err := toonlib.DecodeString(out[strings.Index(out, "help["):])
			if err != nil {
				t.Fatal(err)
			}
			help := decoded.(map[string]any)["help"].([]any)
			if len(help) != 1 {
				t.Fatalf("decoded help = %#v, want one action", help)
			}
			command := help[0].(map[string]any)["cmd"].(string)
			recovered, err := exec.Command("sh", "-c", "set -- "+command+"; printf %s \"$4\"").Output()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(recovered, []byte(path)) {
				t.Fatalf("shell argv = %q, want %q", recovered, path)
			}
		})
	}
}
