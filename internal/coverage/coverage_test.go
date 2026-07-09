package coverage

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/toon"
)

const stories = "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n"
const hdr = "| story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|\n"

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
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 9 | b | s | r | w |\n", "references story 9 but the spec numbers only 3"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| x | b | s | r | w |\n", "has an unrecognized story reference 'x'"},
	}
	for _, c := range cases {
		v := Check(spec(c.body))
		if len(v) == 0 || !strings.Contains(strings.Join(v, "\n"), c.want) {
			t.Errorf("Check violations %v do not contain %q", v, c.want)
		}
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
	return "spec: " + label + "\n" + "state: " + State(p) + "\n" + tbl
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

	t.Run("separator-free slug resolves specs/<slug>.md", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 2 | b2 | s2 | r2 | w2 |\n")
		mustWrite(t, "specs/foo.md", body)

		out, code := Command([]string{"foo"})
		if want := wantTable(t, "specs/foo.md", body); out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0)", out, code, want)
		}
	})

	t.Run("slug already ending .md is not double-appended", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 3 | b3 | s3 | r3 | w3 |\n")
		mustWrite(t, "specs/bar.md", body)

		out, code := Command([]string{"bar.md"})
		if want := wantTable(t, "specs/bar.md", body); out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0) — a double-append would look up specs/bar.md.md and miss", out, code, want)
		}
	})

	t.Run("slug matching no file names both forms tried, exit 1", func(t *testing.T) {
		t.Chdir(t.TempDir())

		out, code := Command([]string{"missing"})
		if code != 1 || !strings.Contains(out, "spec not found: missing, specs/missing.md") {
			t.Errorf("Command = (%q, %d), want exit 1 naming both 'missing' and 'specs/missing.md' in the not-found message", out, code)
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
		cwdBody := mapped("| 4 | cwd | cwd-seam | cwd-red | cwd-why |\n")
		specsBody := mapped("| 5 | specs | specs-seam | specs-red | specs-why |\n")
		mustWrite(t, "foo", cwdBody)
		mustWrite(t, "specs/foo.md", specsBody)

		out, code := Command([]string{"foo"})
		if want := wantTable(t, "foo", cwdBody); out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0) — the CWD file should shadow the specs/ fallback", out, code, want)
		}
	})

	t.Run("present-but-unreadable file surfaces the read error, not not-found", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("root reads 0000-mode files; permission case unobservable")
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

	t.Run("slug matching a directory falls back to specs/<slug>.md", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | r | w |\n")
		mustWrite(t, "dist/keep", "x") // a directory named like the slug
		mustWrite(t, "specs/dist.md", body)

		out, code := Command([]string{"dist"})
		if want := wantTable(t, "specs/dist.md", body); out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0) — a directory is not a spec candidate", out, code, want)
		}
	})

	t.Run("--check resolves a slug and reports violations under the resolved label", func(t *testing.T) {
		t.Chdir(t.TempDir())
		mustWrite(t, "specs/chk.md", mapped("| 1 | | s | r | w |\n")) // empty behavior cell

		out, code := Command([]string{"--check", "chk"})
		if code != 1 || !strings.Contains(out, "error: specs/chk.md coverage map row 1 has an empty 'behavior' cell") {
			t.Errorf("Command = (%q, %d), want exit 1 with the violation under the resolved specs/chk.md label", out, code)
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
