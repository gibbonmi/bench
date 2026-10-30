package conformance

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/tickets"
)

// ticketGrammarPrefix opens every diagnostic this check emits, so a gate reader
// attributes the red to the ticket grammar and not to the preflight rows that
// grade the same files at phase entry.
const ticketGrammarPrefix = "ticket-grammar: "

// checkTicketGrammar grades every staged spec's ticket files and the command-package list.
// The gate owns this sweep, so a malformed ticket cannot land green behind a phase that
// never re-ran preflight.
func checkTicketGrammar(root string) []string {
	return append(sweepTicketFolders(root), gradeCommandPackages(approvedAXIQueries, subcommandRouting, tickets.CommandPackage)...)
}

// gradeCommandPackages reds each approved AXI query command whose grammar package the
// command-package predicate does not hold. The approved query set is the AXI check's own
// inventory, and the grammar package is the routing check's own row, so a new verb cannot
// ship outside the command-package list. A nested family dispatches its leaves from the
// dispatch package, so that package holds its grammar. A query with no routing row reds.
func gradeCommandPackages(approved map[string][]string, routing map[string]routingEntry, commandPackage func(string) bool) []string {
	commands := make([]string, 0, len(approved))
	for command := range approved {
		commands = append(commands, command)
	}
	sort.Strings(commands)
	var diags []string
	for _, command := range commands {
		entry, ok := routing[command]
		if !ok {
			diags = append(diags, fmt.Sprintf("%sAXI query command %q has no subcommand-routing row", ticketGrammarPrefix, command))
			continue
		}
		pkg := entry.Pkg
		if pkg == "" {
			pkg = path.Dir(dispatchFile)
		}
		if !commandPackage(pkg) {
			diags = append(diags, fmt.Sprintf("%sAXI query command %q has no command-package entry for %s", ticketGrammarPrefix, command, pkg))
		}
	}
	return diags
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
// force there. The enumeration is the tickets package seam preflight reads
// through, so the two venues cannot disagree about what a ticket folder holds.
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
	files, duplicates, refusal := tickets.Enumerate(dir, classified.Entries)
	if refusal != nil {
		return []string{ticketGrammarPrefix + refusal.Message(rel)}
	}
	tag := ticketGrammarSpecTag(filepath.Join(root, "specs", slug, "spec.md"))
	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name)
	}
	var diags []string
	for _, duplicate := range duplicates {
		diags = append(diags, ticketGrammarPrefix+rel+": "+duplicate)
	}
	parsed := make([]tickets.Ticket, 0, len(files))
	for _, file := range files {
		ticket, ticketDiags := tickets.ParseTicket(file.Name, file.Data, names, tag)
		parsed = append(parsed, ticket)
		if field, value, unrepresentable := tickets.UnrepresentableValue(ticket); unrepresentable {
			diags = append(diags, fmt.Sprintf("%s%s: %s declares a %s: entry %q with a byte spec-TOON cannot represent", ticketGrammarPrefix, rel, file.Rel, field, value))
			continue
		}
		for _, diag := range ticketDiags {
			diags = append(diags, ticketGrammarPrefix+rel+": "+file.Rel+": "+diag)
		}
	}
	for _, cycle := range tickets.Cycles(parsed) {
		diags = append(diags, ticketGrammarPrefix+rel+": "+cycle)
	}
	return diags
}

// ticketGrammarSpecTag is the tag the spec's declared row IDs share. A folder
// with no readable coverage map answers the empty tag, which is the
// tickets-only posture the parser already understands.
func ticketGrammarSpecTag(specPath string) string {
	_, ids, _, err := coverage.ParseSpec(specPath)
	if err != nil || len(ids) == 0 {
		return ""
	}
	return tickets.TagOf(ids[0])
}

func ticketGrammarUnreadable(rel string, state bounds.FileState, reason string) string {
	return fmt.Sprintf("%s%s is %s: %s", ticketGrammarPrefix, rel, state, reason)
}

// ticketGrammarRoot builds an empty synthetic root, so a sweep test observes the
// sweep's diagnostics alone.
func ticketGrammarRoot(t *testing.T) string {
	t.Helper()
	return t.TempDir()
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
	want := "ticket-grammar: specs/demo/tickets: one.md: one.md: dangling blocker gone.md"
	diags := checkTicketGrammar(root)
	if len(diags) != 1 || diags[0] != want {
		t.Fatalf("dangling blocker = %v, want exactly [%q]", diags, want)
	}
}

// TestTicketGrammarSweepRefusesASpecialEntry pins the sweep's refusal shape. A
// FIFO under tickets/ is refused by its repo-relative path, so a gate reader
// finds the entry in the tree rather than a path from the machine that ran the
// check.
func TestTicketGrammarSweepRefusesASpecialEntry(t *testing.T) {
	root := ticketGrammarRoot(t)
	writeTicketGrammarFile(t, root, "specs/demo/spec.md", "# Demo\n\n| row | story |\n|---|---|\n| DM1 | 1 |\n")
	writeTicketGrammarFile(t, root, "specs/demo/tickets/one.md", ticketBody("One", "none", "a.go (new)", "DM1"))
	if err := syscall.Mkfifo(filepath.Join(root, "specs", "demo", "tickets", "pipe.md"), 0o600); err != nil {
		t.Fatalf("Mkfifo: %v", err)
	}
	want := "ticket-grammar: specs/demo/tickets/pipe.md is wrong-type: not a regular file: p---------"
	diags := checkTicketGrammar(root)
	if len(diags) != 1 || diags[0] != want {
		t.Fatalf("special entry = %v, want exactly [%q]", diags, want)
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
	want := "ticket-grammar: specs/light/tickets: two.md: two.md: dangling blocker gone.md"
	if diags := checkTicketGrammar(root); len(diags) != 1 || diags[0] != want {
		t.Fatalf("tickets-only folder with a dangling blocker = %v, want exactly [%q]", diags, want)
	}
}

// TestTicketGrammarRedsUnlistedQueryPackage pins TG18 over the command-package list. A list
// that drops the routed diff package reds the approved diff query command, and a query with
// no routing row reds.
func TestTicketGrammarRedsUnlistedQueryPackage(t *testing.T) {
	dropped := func(path string) bool { return path != "internal/diff" && tickets.CommandPackage(path) }
	want := `ticket-grammar: AXI query command "diff" has no command-package entry for internal/diff`
	if diags := gradeCommandPackages(approvedAXIQueries, subcommandRouting, dropped); len(diags) != 1 || diags[0] != want {
		t.Fatalf("dropped diff package = %v, want exactly [%q]", diags, want)
	}
	unrouted := map[string]routingEntry{}
	for name, entry := range subcommandRouting {
		if name != "read" {
			unrouted[name] = entry
		}
	}
	want = `ticket-grammar: AXI query command "read" has no subcommand-routing row`
	if diags := gradeCommandPackages(approvedAXIQueries, unrouted, tickets.CommandPackage); len(diags) != 1 || diags[0] != want {
		t.Fatalf("unrouted read query = %v, want exactly [%q]", diags, want)
	}
}
