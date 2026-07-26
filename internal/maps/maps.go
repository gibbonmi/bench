// Package maps ports `bench maps`: unresolved decision-map tickets plus the
// close-readiness handoff row, as a `maps[N]{map,ticket,type,state}:` TOON table.
//
// One engine (parseFile) feeds both the ticket/handoff row listing and the
// distinct-file count `bench status` consumes, so what "unresolved" and "ready to
// close" mean has a single definition — the two-derivations bug class this slice ends.
// The marker/fence/CRLF/`## Handoff` rules live here once.
package maps

import (
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// DecisionsDir is the one control directory this command reads, repo-relative so every
// error line names the same path an agent would type. It is exported because the status
// board names the same directory in the row it prints when the scan fails.
const DecisionsDir = "decisions"

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:   "bench maps",
	Help:  "usage: bench maps [--count]",
	Flags: []usage.Flag{{Name: "--count"}},
}

var (
	ticketRe  = regexp.MustCompile(`^## #([0-9]+):`)
	handoffRe = regexp.MustCompile(`^## Handoff([ \t]|$)`)
)

// marker classifies an unresolved placeholder anchored at line start: an `— (open` /
// `— (deferred` placeholder or a `GRILL DEFERRED` banner. "" for any normal line.
func marker(line string) string {
	switch {
	case strings.HasPrefix(line, "— (open"):
		return "open"
	case strings.HasPrefix(line, "— (deferred"):
		return "deferred"
	case strings.HasPrefix(line, "GRILL DEFERRED"):
		return "grill-deferred"
	}
	return ""
}

type ticket struct{ num, typ, state string }

// fileResult is one map file scanned: the unresolved tickets, plus the three
// close-readiness signals both consumers read (a marker outside Handoff means open
// work; whether a `## Handoff` heading was seen outside a fence; the placeholder kind
// found inside the Handoff section).
type fileResult struct {
	tickets          []ticket
	isMap            bool
	seenHandoff      bool
	handoffState     string
	preHandoffMarker bool
}

// handoffIncomplete is the one close-readiness predicate both outputs call: a file
// only owes a `## Handoff` if it is a map (carries a `## #<n>:` ticket heading), and
// a map is incomplete when the heading is absent or its section still holds a placeholder.
func (r fileResult) handoffIncomplete() bool {
	return r.isMap && (!r.seenHandoff || r.handoffState != "")
}

// notCloseReady is true when the file keeps a row: open work outside Handoff, or an
// incomplete Handoff on a map. This is exactly the count predicate status consumes.
func (r fileResult) notCloseReady() bool {
	return r.preHandoffMarker || r.handoffIncomplete()
}

func parseFile(content []byte) fileResult {
	var r fileResult
	var inFence, inHandoff bool
	var cur *ticket
	flush := func() {
		if cur != nil && cur.state != "" {
			t := cur.typ
			if t == "" {
				t = "unknown"
			}
			r.tickets = append(r.tickets, ticket{cur.num, t, cur.state})
		}
		cur = nil
	}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSuffix(line, "\r")
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if m := ticketRe.FindStringSubmatch(line); m != nil {
			r.isMap = true
			inHandoff = false // a `## #…:` heading is also a `## ` heading
			flush()
			cur = &ticket{num: m[1]}
			continue
		}
		if handoffRe.MatchString(line) {
			r.seenHandoff = true
			inHandoff = true
			continue
		}
		if strings.HasPrefix(line, "## ") {
			inHandoff = false
		}
		if mk := marker(line); mk != "" {
			if inHandoff {
				if r.handoffState == "" {
					r.handoffState = mk
				}
			} else {
				r.preHandoffMarker = true
			}
		}
		if strings.HasPrefix(line, "Type:") {
			if cur != nil {
				cur.typ = strings.TrimLeft(strings.TrimPrefix(line, "Type:"), " \t")
			}
			continue
		}
		if cur != nil && cur.state == "" && !inHandoff {
			if mk := marker(line); mk != "" {
				cur.state = mk
			}
		}
	}
	flush()
	return r
}

// fileRows renders one file's rows: its unresolved ticket rows, then a close-readiness
// handoff row when the file is a zero-open map with no filled Handoff. The `ticket`
// cell is a genuine int for a real (all-digit) ticket so it emits bare, and the literal
// string `handoff` for a close-readiness row — the column is genuinely mixed.
func fileRows(name string, r fileResult) [][]any {
	var rows [][]any
	for _, t := range r.tickets {
		num, _ := strconv.Atoi(t.num) // t.num is all-digit by ticketRe
		rows = append(rows, []any{name, num, t.typ, t.state})
	}
	if !r.preHandoffMarker && r.handoffIncomplete() {
		if !r.seenHandoff {
			rows = append(rows, []any{name, "handoff", "handoff", "missing"})
		} else {
			rows = append(rows, []any{name, "handoff", "handoff", r.handoffState})
		}
	}
	return rows
}

// issue is one decisions file that did not feed the parser cleanly: unreadable,
// wrong-type, malformed bytes, or — the parser-owned state — bytes that read fine
// but attempt no ticket heading and no marker at all, so no rule in this package
// recognizes the file as either a map or a scoped note. State is the classifier's
// state for every case but the last, where it is bounds.StateUnsupportedSchema.
type issue struct {
	name, reason string
	state        bounds.FileState
}

// scanResult is the one engine Rows, UnresolvedCount, and Command all narrow from:
// every cleanly parsed file plus every file that could not be, so the listing and
// the count can never disagree about what was readable. dirState is decisions/'s
// own readability — StateAbsent (no such directory: the authoritative empty state)
// or StateEmpty leave names and issues both nil; a state FileState.Failed reports
// means the directory itself could not be enumerated, so names and issues stay nil
// even though nothing was actually confirmed absent, and a caller must read
// dirState rather than treat the zero-length results as authoritative.
type scanResult struct {
	names     []string
	results   map[string]fileResult
	issues    []issue
	dirState  bounds.FileState
	dirReason string
}

// isDirectoryDoc reports whether a decisions entry documents the directory rather than
// claiming membership in it. A README is the one universally understood name for that,
// and it is the only exemption: every other ordinary name in `decisions/` is a candidate
// decision map, so a genuinely broken map still earns its row and its exit code. Without
// the exemption a routine index file turns `bench maps` red for no defect.
func isDirectoryDoc(name string) bool {
	return strings.EqualFold(strings.TrimSuffix(name, ".md"), "README")
}

// scan classifies root/decisions and every *.md entry inside it, in sorted order.
func scan(root string) scanResult {
	dir := filepath.Join(root, DecisionsDir)
	cd := bounds.ClassifyDir(dir)
	s := scanResult{dirState: cd.State, dirReason: cd.Reason, results: map[string]fileResult{}}
	if cd.State != bounds.StateParsed {
		return s
	}
	for _, e := range cd.Entries {
		// A leading-dot file is skipped to match the shell glob decisions/*.md, which
		// without dotglob never expands a hidden name — so a `.foo.md` stays invisible to
		// the listing and the count (and cannot falsely nag the dashboard).
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if isDirectoryDoc(e.Name()) {
			continue
		}
		base := strings.TrimSuffix(e.Name(), ".md")
		fc := bounds.Classify(filepath.Join(dir, e.Name()), bounds.ControlRecordLimit)
		if fc.State == bounds.StateParsed {
			r := parseFile(fc.Data)
			if r.isMap || r.preHandoffMarker || r.handoffState != "" {
				s.names = append(s.names, base)
				s.results[base] = r
				continue
			}
			s.issues = append(s.issues, issue{name: base, state: bounds.StateUnsupportedSchema, reason: "no ticket heading or marker found"})
			continue
		}
		s.issues = append(s.issues, issue{name: base, state: fc.State, reason: fc.Reason})
	}
	sort.Strings(s.names)
	sort.Slice(s.issues, func(i, j int) bool { return s.issues[i].name < s.issues[j].name })
	return s
}

// Rows lists every unresolved ticket and close-readiness handoff row under
// root/decisions, in sorted file order, plus one row per file that could not be
// classified — naming it and its state — so a broken file is a visible row rather
// than a silent omission from the listing (story 10).
func Rows(root string) [][]any {
	return rowsFromScan(scan(root))
}

func rowsFromScan(s scanResult) [][]any {
	var rows [][]any
	for _, name := range s.names {
		rows = append(rows, fileRows(name, s.results[name])...)
	}
	for _, iss := range s.issues {
		rows = append(rows, []any{iss.name, "error", string(iss.state), iss.reason})
	}
	return rows
}

// unresolvedCount is the DISTINCT not-close-ready file count: every parsed file
// whose notCloseReady predicate holds, plus every file this package could not
// classify at all — an unreadable or unrecognized file is exactly as unresolved
// as an open ticket, and it also earns a row (story 10), so the count and the
// row listing agree by construction (story 11).
func (s scanResult) unresolvedCount() int {
	n := len(s.issues)
	for _, r := range s.results {
		if r.notCloseReady() {
			n++
		}
	}
	return n
}

// UnresolvedCount is the count status surfaces, paired with the scan's own
// readability state so a failed scan cannot fabricate zero. state is
// bounds.StateParsed whenever n is trustworthy — including the legitimately
// empty cases (no decisions/ directory, or none of its files unresolved) — and
// the decisions/ directory's own failure state when the whole scan could not run
// at all, in which case n is always 0 and must not be read as "nothing
// unresolved."
func UnresolvedCount(root string) (n int, state bounds.FileState) {
	s := scan(root)
	if s.dirState.Failed() {
		return 0, s.dirState
	}
	return s.unresolvedCount(), bounds.StateParsed
}

// Command implements `bench maps`. `--count` is the status adapter's hook: it
// prints the distinct not-close-ready file count as a bare integer, through the
// same scan the listing uses. Absence of root/decisions is the only authoritative
// empty state (exit 0); a decisions/ directory that exists but cannot be
// enumerated exits 1 with a structured `error:` line on *both* forms, because a
// bare integer for a scan that never ran is a count nothing supports and it would
// contradict the listing's own error on the same tree (story 11). Once a directory
// listing exists, a per-file failure is not whole-document: every readable file
// still renders its rows, plus one row per file that failed, and the command exits
// 1 (story 10).
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		if _, count := parsed.Flags["--count"]; count {
			return "0\n", 0
		}
		return toon.NotInRepo() + "\n", 1
	}
	s := scan(root)
	if s.dirState.Failed() {
		return toon.RecordError(DecisionsDir, s.dirState, s.dirReason) + "\n", 1
	}
	if _, count := parsed.Flags["--count"]; count {
		return strconv.Itoa(s.unresolvedCount()) + "\n", 0
	}
	out, err := toon.TableTyped("maps", []string{"map", "ticket", "type", "state"}, rowsFromScan(s))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	if len(s.issues) > 0 {
		return out, 1
	}
	return out, 0
}
