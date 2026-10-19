package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/toon"
)

// TestCommandTicketsAbsent is C12's absent half: review mode with tickets/ absent is a
// structured red.
func TestCommandTicketsAbsent(t *testing.T) {
	slug := "example"
	preflighttest.StartRepo(t)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug))
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error:") {
		t.Errorf("output = %q, want a structured error", out)
	}
}

// TestCommandSpecialFileInTickets is C12's refusal half: a FIFO inside tickets/ is
// refused before reading, named in the error, without blocking. The goroutine/timeout
// wrapper is the test-level bound the ticket requires: a hang mutation fails the test
// instead of the suite.
func TestCommandSpecialFileInTickets(t *testing.T) {
	_, slug := preflighttest.SeedConformant(t)
	fifoPath := filepath.Join("specs", slug, "tickets", "fifo")
	if err := syscall.Mkfifo(fifoPath, 0o644); err != nil {
		t.Fatalf("Mkfifo: %v", err)
	}

	type result struct {
		out  string
		code int
	}
	done := make(chan result, 1)
	go func() {
		out, code := Command([]string{"review", slug})
		done <- result{out, code}
	}()
	select {
	case r := <-done:
		if r.code != 1 {
			t.Fatalf("Command exit = %d, want 1; output:\n%s", r.code, r.out)
		}
		if !strings.Contains(r.out, "fifo") {
			t.Errorf("output must name the refused fifo:\n%s", r.out)
		}
	case <-time.After(bounds.TestDeadline(5 * time.Second)):
		t.Fatal("Command blocked reading a FIFO instead of refusing it by mode")
	}
}

// TestCommandNotInRepo is C13's not-in-repo branch: outside a git repository the
// command answers the standard not-in-repo error.
func TestCommandNotInRepo(t *testing.T) {
	t.Chdir(t.TempDir())
	out, code := Command([]string{"review", "example"})
	if want := toon.NotInRepo() + "\n"; out != want || code != 1 {
		t.Errorf("Command = (%q, %d), want (%q, 1)", out, code, want)
	}
}

// TestCommandUsageBranches is C13's four usage branches: missing mode, unknown mode,
// missing slug, and unknown flag each exit 2. `build` is a real accepted mode now (see
// TestCommandBuildFresh and its siblings), so it is no longer one of these branches.
func TestCommandUsageBranches(t *testing.T) {
	preflighttest.StartRepo(t)
	preflighttest.RunGit(t, "commit", "-q", "--allow-empty", "-m", "c0")

	cases := []struct {
		name string
		args []string
	}{
		{"missing mode", nil},
		{"unknown mode", []string{"frob", "example"}},
		{"missing slug", []string{"review"}},
		{"unknown flag", []string{"--frob"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, code := Command(c.args)
			if code != 2 {
				t.Errorf("Command(%v) = (%q, %d), want exit 2", c.args, out, code)
			}
		})
	}
}

// TestCommandMissingSpecNamesPath is H1's missing half (mutation
// H1/missing-spec-names-path): no spec at all answers "spec not found" naming the
// tried path, exit 1.
func TestCommandMissingSpecNamesPath(t *testing.T) {
	slug := "example"
	preflighttest.StartRepo(t)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: spec not found") {
		t.Errorf("output = %q, want the spec-not-found error", out)
	}
	if !strings.Contains(out, filepath.Join("specs", slug, "spec.md")) {
		t.Errorf("output must name the tried spec path:\n%s", out)
	}
}

// TestCommandDanglingSymlinkClassifiedBroken is H1's symlink half (mutation
// H1/dangling-symlink-classified): a dangling symlink where the spec should
// be is classified broken, not an authoritative empty state. This is a
// distinct error from the missing-spec case, naming the spec path.
func TestCommandDanglingSymlinkClassifiedBroken(t *testing.T) {
	slug := "example"
	preflighttest.StartRepo(t)
	specPath := filepath.Join("specs", slug, "spec.md")
	if err := os.MkdirAll(filepath.Dir(specPath), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.Symlink("missing-target.md", specPath); err != nil {
		t.Fatalf("Symlink: %v", err)
	}

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if strings.HasPrefix(out, "error: spec not found") {
		t.Errorf("a dangling symlink must not be classified as authoritative absence:\n%s", out)
	}
	if !strings.HasPrefix(out, "error: spec not readable") {
		t.Errorf("output = %q, want the spec-not-readable error", out)
	}
	if !strings.Contains(out, specPath) {
		t.Errorf("output must name the spec path:\n%s", out)
	}
}

// TestCommandInvalidCoverageMapCarriesMessage is H2's invalid-map half (mutation
// H2/validator-message-carried): a coverage map that fails validation surfaces the
// validator's own message, exit 1.
func TestCommandInvalidCoverageMapCarriesMessage(t *testing.T) {
	slug := "example"
	preflighttest.StartRepo(t)
	body := "# " + slug + "\n\nStatus: staged\n\n## User stories\n1. As a, I want b, so c.\n\n" +
		"### Acceptance coverage map\n" +
		"| row | story | behavior | seam | why it catches the failure |\n" +
		"|---|---|---|---|---|\n" +
		"| PF1 | 1 | does x | cli seam |\n" // 4 cells where 5 are wanted
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", body)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: coverage map invalid") {
		t.Errorf("output = %q, want the coverage-map-invalid error", out)
	}
	if !strings.Contains(out, "row 1 has 4 cells (want 5)") {
		t.Errorf("output must carry the validator's own message:\n%s", out)
	}
}

// TestCommandNoRowIDMapNamesOptIn is H2's no-row-ID half (mutation
// H2/optin-hint-named): a valid reduced 4-cell map (no `row` header
// column) is refused. The error names the row-ID opt-in, exit 1.
func TestCommandNoRowIDMapNamesOptIn(t *testing.T) {
	slug := "example"
	preflighttest.StartRepo(t)
	body := "# " + slug + "\n\nStatus: staged\n\n## User stories\n1. As a, I want b, so c.\n\n" +
		"### Acceptance coverage map\n" +
		"| story | behavior | seam | why it catches the failure |\n" +
		"|---|---|---|---|\n" +
		"| 1 | does x | cli seam | catches z |\n"
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", body)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: coverage map not opted into row IDs") {
		t.Errorf("output = %q, want the opt-in error naming the row-ID convention", out)
	}
	if !strings.Contains(out, "`row`") {
		t.Errorf("output must name the leading `row` column opt-in:\n%s", out)
	}
}

// TestCommandNonStagedNamesFoundStatus is H4 (mutation H4/found-status-named): a
// spec whose Status: is anything but staged answers a structured error naming the
// found status, exit 1.
func TestCommandNonStagedNamesFoundStatus(t *testing.T) {
	slug := "example"
	preflighttest.StartRepo(t)
	body := strings.Replace(preflighttest.SpecBody(slug), "Status: staged", "Status: implemented", 1)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", body)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: spec not staged") {
		t.Errorf("output = %q, want the spec-not-staged error", out)
	}
	if !strings.Contains(out, "Status: implemented") {
		t.Errorf("output must name the found status:\n%s", out)
	}
}

// TestCommandTrailingNewlineParity is C14. A spec whose last line lacks a
// trailing newline, and a ticket file whose last line lacks one, parse
// identically to their terminated forms.
func TestCommandTrailingNewlineParity(t *testing.T) {
	slug := "example"
	terminated := preflighttest.SpecBody(slug)
	unterminated := strings.TrimSuffix(terminated, "\n")
	if unterminated == terminated {
		t.Fatal("fixture invalid: SpecBody already lacks a trailing newline")
	}

	for _, tc := range []struct {
		name      string
		specBody  string
		ticketDoc string
	}{
		{"terminated", terminated, preflighttest.TicketDoc("One", "PF1", "PF2")},
		{"unterminated spec", unterminated, preflighttest.TicketDoc("One", "PF1", "PF2")},
		{"unterminated ticket", terminated, strings.TrimSuffix(preflighttest.TicketDoc("One", "PF1", "PF2"), "\n")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			preflighttest.StartRepo(t)
			preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", tc.specBody)
			preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", tc.ticketDoc)
			preflighttest.RunGit(t, "add", ".")
			preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
			preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
			preflighttest.MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
			preflighttest.RunGit(t, "add", "internal/"+slug+"/foo.go")
			preflighttest.RunGit(t, "commit", "-q", "-m", "c1")

			out, code := Command([]string{"review", slug})
			if code != 0 {
				t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
			}
			verdicts := renderedVerdicts(t, out, modeReview, slug)
			for _, name := range []string{"rows-owned", "rows-membership"} {
				if verdicts[name] != verdictGreen {
					t.Errorf("%s: %s is %q (trailing-newline handling dropped the last citation):\n%s", tc.name, name, verdicts[name], out)
				}
			}
		})
	}
}
