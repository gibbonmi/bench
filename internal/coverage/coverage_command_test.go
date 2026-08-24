// Tests for Command, the (args) -> (output, exit code) surface a caller drives.
package coverage

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/toon"
)

// TestCommandRendersShortRow drives the same shortfall through Command, the surface a
// caller actually reads. A short row still renders as a full three-column TOON record,
// with its absent cells empty. It also renders alongside the repair action its width
// violation earns.
func TestCommandRendersShortRow(t *testing.T) {
	t.Chdir(t.TempDir())
	mustWrite(t, "spec.md", "# t\n\n"+stories+"\n### Acceptance coverage map\n"+hdrReduced+"| 1 | does x |\n")

	out, code := Command([]string{"spec.md"})
	want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",does x,\"\"\nhelp[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"
	if code != 0 || out != want {
		t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
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

// TestCommand drives Command through its public (args) -> (output, exit code)
// interface only, and never calls parse or Check directly.
func TestCommand(t *testing.T) {
	mapped := func(row string) string {
		return "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + row
	}

	t.Run("readable path argument resolves and round-trips", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | w |\n")
		mustWrite(t, "spec.md", body)

		out, code := Command([]string{"spec.md"})
		want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\nhelp[1]{cmd,why}:\n  bench coverage --check spec.md,check coverage for stories 1\n"
		if out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0)", out, code, want)
		}
	})

	t.Run("appends one check action per unchecked row and deduplicates exact templates", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | catches one |\n| 2 | c | t | catches two |\n| 1 | b | s | catches one |\n")
		mustWrite(t, "path with spaces/spec.md", body)
		out, code := Command([]string{"path with spaces/spec.md"})
		want := "spec: path with spaces/spec.md\nstate: mapped\nrows[3]{story,behavior,seam}:\n  \"1\",b,s\n  \"2\",c,t\n  \"1\",b,s\nhelp[2]{cmd,why}:\n  bench coverage --check 'path with spaces/spec.md',check coverage for stories 1\n  bench coverage --check 'path with spaces/spec.md',check coverage for stories 2\n"
		if code != 0 || out != want {
			t.Fatalf("Command = (%d, %q), want (%d, %q)", code, out, 0, want)
		}
	})

	t.Run("repair prose does not alter unchecked classification", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | catches one |\n")
		mustWrite(t, "spec.md", body)
		out, code := Command([]string{"spec.md"})
		if code != 0 || !strings.Contains(out, "bench coverage --check spec.md,check coverage for stories 1") {
			t.Fatalf("Command = (%d, %q), want unchecked check action", code, out)
		}
	})

	t.Run("malformed map gets repair retry without changing extraction exit", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 9 | b | s | catches one |\n")
		mustWrite(t, "spec.md", body)
		out, code := Command([]string{"spec.md"})
		want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"9\",b,s\nhelp[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"
		if code != 0 || out != want {
			t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
		}
	})

	// An unrecognized header has no descriptor of its own. Its rows project through the
	// four-cell reduced schema, the projection fallback. A caller still gets the story,
	// behavior, and seam it seeds tasks from, next to the repair action the unknown
	// header earns. If the rows were dropped, the map's content would stay unreachable
	// until someone fixes the header by hand.
	t.Run("unrecognized header projects rows through the fallback schema", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := "# t\n\n" + stories + "\n### Acceptance coverage map\n" +
			"| story | behavior | seam | outcome |\n|---|---|---|---|\n| 1 | b | s | w |\n"
		mustWrite(t, "spec.md", body)

		out, code := Command([]string{"spec.md"})
		want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\nhelp[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"
		if code != 0 || out != want {
			t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
		}
	})

	t.Run("canonical mapped zero-row table is terminal", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("")
		mustWrite(t, "spec.md", body)
		out, code := Command([]string{"spec.md"})
		want := "spec: spec.md\nstate: mapped\nrows[0]{story,behavior,seam}:\nhelp[0]{cmd,why}:\n"
		if out != want || code != 0 {
			t.Fatalf("Command = (%d, %q), want (%d, %q)", code, out, 0, want)
		}
	})

	// The behavior cell is author-controlled prose that reaches the TOON encoder, so a
	// control byte in it refuses the whole response. The AXI error contract replaces
	// the response rather than rendering a table with one lossy cell in it.
	t.Run("control-bearing behavior cell refuses the whole response", func(t *testing.T) {
		t.Chdir(t.TempDir())
		mustWrite(t, "spec.md", mapped("| 1 | b\x1bx | s | w |\n"))

		out, code := Command([]string{"spec.md"})
		want := "error: unrepresentable TOON cell — toon: unsupported control character U+001B in string\n"
		if code != 1 || out != want {
			t.Fatalf("Command = (%d, %q), want (1, %q)", code, out, want)
		}
	})

	// A comma is the row delimiter, and a tab is an escapable control. A behavior
	// carrying either is the case where an unquoted cell would silently split or
	// corrupt the record. Both stay one three-field row, quoted by the encoder.
	for _, tc := range []struct{ name, behavior, row string }{
		{"comma", "a,b", "  \"1\",\"a,b\",s\n"},
		{"tab", "does\tx", "  \"1\",\"does\\tx\",s\n"},
	} {
		t.Run("delimiter-bearing behavior renders one quoted row: "+tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			mustWrite(t, "spec.md", mapped("| 1 | "+tc.behavior+" | s | w |\n"))

			out, code := Command([]string{"spec.md"})
			want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n" + tc.row +
				"help[1]{cmd,why}:\n  bench coverage --check spec.md,check coverage for stories 1\n"
			if code != 0 || out != want {
				t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
			}
		})
	}

	t.Run("separator-free slug resolves specs/<slug>/spec.md", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b2 | s2 | w2 |\n")
		mustWrite(t, "specs/foo/spec.md", body)

		out, code := Command([]string{"foo"})
		want := "spec: specs/foo/spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b2,s2\nhelp[1]{cmd,why}:\n  bench coverage --check specs/foo/spec.md,check coverage for stories 1\n"
		if out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0)", out, code, want)
		}
	})

	t.Run("slug already ending .md resolves folder spec", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b3 | s3 | w3 |\n")
		mustWrite(t, "specs/bar/spec.md", body)

		out, code := Command([]string{"bar.md"})
		want := "spec: specs/bar/spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b3,s3\nhelp[1]{cmd,why}:\n  bench coverage --check specs/bar/spec.md,check coverage for stories 1\n"
		if out != want || code != 0 {
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
		cwdBody := mapped("| 1 | cwd | cwd-seam | cwd-why |\n")
		specsBody := mapped("| 5 | specs | specs-seam | specs-why |\n")
		mustWrite(t, "foo", cwdBody)
		mustWrite(t, "specs/foo/spec.md", specsBody)

		out, code := Command([]string{"foo"})
		want := "spec: foo\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",cwd,cwd-seam\nhelp[1]{cmd,why}:\n  bench coverage --check foo,check coverage for stories 1\n"
		if out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0) — the CWD file should shadow the specs/ fallback", out, code, want)
		}
	})

	t.Run("present-but-unreadable file surfaces the read error, not not-found", func(t *testing.T) {
		if os.Geteuid() == 0 {
			capability.Capability(t, capability.Privilege, "root reads 0000-mode files; permission case unobservable")
		}
		t.Chdir(t.TempDir())
		mustWrite(t, "sub/locked.md", mapped("| 1 | b | s | w |\n"))
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
		body := mapped("| 1 | b | s | w |\n")
		mustWrite(t, "dist/keep", "x") // a directory named like the slug
		mustWrite(t, "specs/dist/spec.md", body)

		out, code := Command([]string{"dist"})
		want := "spec: specs/dist/spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\nhelp[1]{cmd,why}:\n  bench coverage --check specs/dist/spec.md,check coverage for stories 1\n"
		if out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0) — a directory is not a spec candidate", out, code, want)
		}
	})

	t.Run("--check resolves a slug and reports violations under the resolved label", func(t *testing.T) {
		t.Chdir(t.TempDir())
		mustWrite(t, "specs/chk/spec.md", mapped("| 1 | | s | w |\n")) // empty behavior cell

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
		return "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + row
	}
	for _, tc := range []struct {
		name, fixture, body, help string
	}{
		{"mapped actionable", "pre-disclosure-mapped.stdout", mapped("| 1 | b | s | catches one |\n"), "help[1]{cmd,why}:\n  bench coverage --check spec.md,check coverage for stories 1\n"},
		{"repairable malformed", "pre-disclosure-malformed.stdout", mapped("| 9 | b | s | catches one |\n"), "help[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"},
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
	mapped := "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + "| 1 | b | s | catches one |\n"
	for _, control := range []string{"\t", "\n", "\r", "\x1b"} {
		t.Run("control", func(t *testing.T) {
			t.Chdir(t.TempDir())
			path := "control" + control + "spec.md"
			mustWrite(t, path, mapped)
			out, code := Command([]string{path})
			primary := "spec: " + path + "\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\n"
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
			argv, err := axitest.RecoverHelpCommandArgv(out)
			if err != nil {
				t.Fatal(err)
			}
			want := []string{"bench", "coverage", "--check", path}
			if !slices.Equal(argv, want) {
				t.Fatalf("shell argv = %q, want %q", argv, want)
			}
		})
	}
}

func TestCommandAngleBracketSpecPathPreservesPrimaryAndHonestFallback(t *testing.T) {
	mapped := "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + "| 1 | b | s | catches one |\n"
	for _, marker := range []string{"<", ">"} {
		t.Run(marker, func(t *testing.T) {
			t.Chdir(t.TempDir())
			path := "angle" + marker + "spec.md"
			mustWrite(t, path, mapped)
			out, code := Command([]string{path})
			primary := "spec: " + path + "\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\n"
			want := primary + "help[0]{cmd,why}:\n"
			if code != 0 || out != want {
				t.Fatalf("Command = (%d, %q), want checked primary plus honest empty help", code, out)
			}
		})
	}
}
