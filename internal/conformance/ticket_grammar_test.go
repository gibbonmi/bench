package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/tickets"
)

// ticketGrammarPrefix opens every diagnostic this check emits, so a gate reader
// attributes the red to the ticket grammar and not to the preflight rows that
// grade the same files at phase entry.
const ticketGrammarPrefix = "ticket-grammar: "

// ticketGrammarOwners are the binding rows the table must always carry: the
// dispatcher, the renderer, and the terminal-lifecycle owner. They are
// established owners that recent builds reconstructed by hand, so their absence
// is a fault rather than a choice.
var ticketGrammarOwners = []string{"cmd/bench", "internal/toon", "internal/terminal"}

// checkTicketGrammar grades every staged spec's ticket files and the binding
// table the ownership closure reads. The gate owns this sweep, so a malformed
// ticket cannot land green behind a phase that never re-ran preflight.
func checkTicketGrammar(root string) []string {
	diags := sweepTicketFolders(root)
	return append(diags, gradeTicketBindings(root, tickets.Bindings(), approvedAXIQueries)...)
}

// sweepTicketFolders grades specs/<slug>/tickets/ for every spec folder. An
// absent specs tree or an absent tickets directory is silent, because a staged
// spec before slicing is a legal tree state. Anything present that will not
// read is a diagnostic, so the sweep never reads a refusal as green.
func sweepTicketFolders(root string) []string {
	specsDir := filepath.Join(root, "specs")
	classified := bounds.ClassifyDir(specsDir)
	switch classified.State {
	case bounds.StateAbsent, bounds.StateEmpty:
		return nil
	case bounds.StateParsed:
	default:
		return []string{ticketGrammarUnreadable("specs", classified.State, classified.Reason)}
	}
	var slugs []string
	for _, entry := range classified.Entries {
		if entry.IsDir() {
			slugs = append(slugs, entry.Name())
		}
	}
	sort.Strings(slugs)
	var diags []string
	for _, slug := range slugs {
		diags = append(diags, gradeTicketFolder(root, slug)...)
	}
	return diags
}

// gradeTicketFolder grades one spec folder's tickets. The spec tag scopes the
// citation grammar. A folder with no spec.md is a light-path tickets folder:
// the empty tag skips the Covers row checks, and every other grammar stays in
// force there.
func gradeTicketFolder(root, slug string) []string {
	rel := "specs/" + slug + "/tickets"
	dir := filepath.Join(root, "specs", slug, "tickets")
	classified := bounds.ClassifyDir(dir)
	switch classified.State {
	case bounds.StateAbsent, bounds.StateEmpty:
		return nil
	case bounds.StateParsed:
	default:
		return []string{ticketGrammarUnreadable(rel, classified.State, classified.Reason)}
	}
	files, refusal := ticketGrammarFiles(dir, classified)
	if refusal != "" {
		return []string{ticketGrammarUnreadable(rel, bounds.StateUnreadable, refusal)}
	}
	tag := ticketGrammarSpecTag(filepath.Join(root, "specs", slug, "spec.md"))
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	var diags []string
	parsed := make([]tickets.Ticket, 0, len(names))
	for _, name := range names {
		ticket, ticketDiags := tickets.ParseTicket(name, files[name], names, tag)
		parsed = append(parsed, ticket)
		for _, diag := range ticketDiags {
			diags = append(diags, ticketGrammarPrefix+rel+": "+diag)
		}
	}
	for _, cycle := range tickets.Cycles(parsed) {
		diags = append(diags, ticketGrammarPrefix+rel+": "+cycle)
	}
	return diags
}

// ticketGrammarFiles collects every `.md` ticket under one already-classified
// directory, keyed by basename, recursing with the same lstat-first
// classification at every depth. A non-`.md` entry is an asset the grammar
// ignores. A refusal anywhere below returns the reason, so the folder reds
// rather than grading a partial listing.
func ticketGrammarFiles(dir string, classified bounds.ClassifiedDir) (map[string][]byte, string) {
	files := map[string][]byte{}
	for _, entry := range classified.Entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			below := bounds.ClassifyDir(path)
			switch below.State {
			case bounds.StateEmpty:
				continue
			case bounds.StateParsed:
			default:
				return nil, entry.Name() + " is " + string(below.State) + ": " + below.Reason
			}
			nested, refusal := ticketGrammarFiles(path, below)
			if refusal != "" {
				return nil, refusal
			}
			for name, data := range nested {
				files[name] = data
			}
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		read := bounds.ClassifyNoFollow(path)
		switch read.State {
		case bounds.StateEmpty:
			files[entry.Name()] = nil
		case bounds.StateParsed:
			files[entry.Name()] = read.Data
		default:
			return nil, entry.Name() + " is " + string(read.State) + ": " + read.Reason
		}
	}
	return files, ""
}

// ticketGrammarSpecTag is the alphabetic prefix the spec's declared row IDs
// share. A folder with no readable coverage map answers the empty tag, which is
// the tickets-only posture the parser already understands.
func ticketGrammarSpecTag(specPath string) string {
	_, ids, _, err := coverage.ParseSpec(specPath)
	if err != nil || len(ids) == 0 {
		return ""
	}
	id := ids[0]
	for i := 0; i < len(id); i++ {
		if id[i] >= '0' && id[i] <= '9' {
			return id[:i]
		}
	}
	return id
}

// gradeTicketBindings grades the command-to-registry binding table the
// ownership closure reads. It proves three things: every bound file exists,
// every approved AXI query command's package carries a row, and the three owner
// rows are present. The approved query set is the AXI check's own inventory, so
// a new verb cannot ship with no registry closure at all.
func gradeTicketBindings(root string, rows []tickets.BindingRow, approved map[string][]string) []string {
	bound := make(map[string]bool, len(rows))
	var diags []string
	for _, row := range rows {
		bound[row.Prefix] = true
		for _, file := range row.Files {
			if !exists(filepath.Join(root, filepath.FromSlash(file))) {
				diags = append(diags, fmt.Sprintf("%sbinding row %s names %s, which the tree does not hold", ticketGrammarPrefix, row.Prefix, file))
			}
		}
	}
	commands := make([]string, 0, len(approved))
	for command := range approved {
		commands = append(commands, command)
	}
	sort.Strings(commands)
	for _, command := range commands {
		pkg := "internal/" + command
		if !bound[pkg] {
			diags = append(diags, fmt.Sprintf("%sAXI query command %q has no binding row for %s", ticketGrammarPrefix, command, pkg))
		}
	}
	for _, owner := range ticketGrammarOwners {
		if !bound[owner] {
			diags = append(diags, fmt.Sprintf("%sbinding table has no owner row for %s", ticketGrammarPrefix, owner))
		}
	}
	return diags
}

func ticketGrammarUnreadable(rel string, state bounds.FileState, reason string) string {
	return fmt.Sprintf("%s%s is %s: %s", ticketGrammarPrefix, rel, state, reason)
}

// ticketGrammarRoot builds a synthetic root that carries every file the binding
// table names, so a sweep test observes the sweep's diagnostics alone.
func ticketGrammarRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, row := range tickets.Bindings() {
		for _, file := range row.Files {
			writeTicketGrammarFile(t, root, file, "package placeholder\n")
		}
	}
	return root
}

func writeTicketGrammarFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ticketBody renders one well-formed ticket, so a test mutates the single field
// it means to grade and nothing else.
func ticketBody(title, blockers, writes, covers string) string {
	return "# " + title + "\n\nBlocked by: " + blockers + "\nWrites: " + writes + "\nCovers: " + covers +
		"\n\n## What to build\n\nOne thing.\n\n## Acceptance\n\n- [ ] it works.\n"
}

// TestTicketGrammarSweepRedsDanglingBlocker pins TG27.
func TestTicketGrammarSweepRedsDanglingBlocker(t *testing.T) {
	root := ticketGrammarRoot(t)
	writeTicketGrammarFile(t, root, "specs/demo/spec.md", "# Demo\n\n| row | story |\n|---|---|\n| DM1 | 1 |\n")
	writeTicketGrammarFile(t, root, "specs/demo/tickets/one.md", ticketBody("One", "gone.md", "a.go (new)", "DM1"))
	want := "ticket-grammar: specs/demo/tickets: one.md: dangling blocker gone.md"
	diags := checkTicketGrammar(root)
	if len(diags) != 1 || diags[0] != want {
		t.Fatalf("dangling blocker = %v, want exactly [%q]", diags, want)
	}
}

// TestTicketGrammarSweepToleratesAbsentTicketsDir pins TG28.
func TestTicketGrammarSweepToleratesAbsentTicketsDir(t *testing.T) {
	root := ticketGrammarRoot(t)
	writeTicketGrammarFile(t, root, "specs/demo/spec.md", "# Demo\n\n| row | story |\n|---|---|\n| DM1 | 1 |\n")
	if diags := checkTicketGrammar(root); len(diags) != 0 {
		t.Fatalf("staged spec with no tickets directory = %v, want no diagnostic", diags)
	}
}

// TestTicketGrammarSweepSkipsCoversWithoutSpec pins TG25. A tickets folder with
// no spec.md declares no rows, so a foreign-looking token is not graded against
// a tag. Every other grammar still bites there.
func TestTicketGrammarSweepSkipsCoversWithoutSpec(t *testing.T) {
	root := ticketGrammarRoot(t)
	writeTicketGrammarFile(t, root, "specs/light/tickets/one.md", ticketBody("One", "none", "a.go (new)", "XY7"))
	if diags := checkTicketGrammar(root); len(diags) != 0 {
		t.Fatalf("tickets-only folder = %v, want the Covers checks skipped", diags)
	}
	writeTicketGrammarFile(t, root, "specs/light/tickets/two.md", ticketBody("Two", "gone.md", "b.go (new)", "XY8"))
	want := "ticket-grammar: specs/light/tickets: two.md: dangling blocker gone.md"
	if diags := checkTicketGrammar(root); len(diags) != 1 || diags[0] != want {
		t.Fatalf("tickets-only folder with a dangling blocker = %v, want exactly [%q]", diags, want)
	}
}

// TestTicketGrammarRedsMissingRegistryFile pins TG17.
func TestTicketGrammarRedsMissingRegistryFile(t *testing.T) {
	root := ticketGrammarRoot(t)
	if err := os.Remove(filepath.Join(root, filepath.FromSlash("internal/toon/toon_test.go"))); err != nil {
		t.Fatal(err)
	}
	want := "ticket-grammar: binding row internal/toon names internal/toon/toon_test.go, which the tree does not hold"
	if diags := checkTicketGrammar(root); len(diags) != 1 || diags[0] != want {
		t.Fatalf("absent bound file = %v, want exactly [%q]", diags, want)
	}
}

// TestTicketGrammarRedsUnboundQueryPackage pins TG18.
func TestTicketGrammarRedsUnboundQueryPackage(t *testing.T) {
	root := ticketGrammarRoot(t)
	var rows []tickets.BindingRow
	for _, row := range tickets.Bindings() {
		if row.Prefix != "internal/anchors" {
			rows = append(rows, row)
		}
	}
	want := `ticket-grammar: AXI query command "anchors" has no binding row for internal/anchors`
	diags := gradeTicketBindings(root, rows, approvedAXIQueries)
	if len(diags) != 1 || diags[0] != want {
		t.Fatalf("dropped anchors row = %v, want exactly [%q]", diags, want)
	}
}

// TestTicketGrammarRedsMissingOwnerRow pins TG36. The terminal row binds no AXI
// query command, so only the owner rule can name its absence.
func TestTicketGrammarRedsMissingOwnerRow(t *testing.T) {
	root := ticketGrammarRoot(t)
	var rows []tickets.BindingRow
	for _, row := range tickets.Bindings() {
		if row.Prefix != "internal/terminal" {
			rows = append(rows, row)
		}
	}
	want := "ticket-grammar: binding table has no owner row for internal/terminal"
	diags := gradeTicketBindings(root, rows, approvedAXIQueries)
	if len(diags) != 1 || diags[0] != want {
		t.Fatalf("dropped terminal row = %v, want exactly [%q]", diags, want)
	}
}
