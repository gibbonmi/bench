// Package maps ports `bench maps`: unresolved decision-map tickets plus the
// close-readiness handoff row, as a `maps[N]{map,ticket,type,state}:` TOON table.
//
// One engine (parseFile) feeds both the ticket/handoff row listing and the
// distinct-file count `bench status` consumes, so what "unresolved" and "ready to
// close" mean has a single definition — the two-derivations bug class this slice ends.
// The marker/fence/CRLF/`## Handoff` rules live here once.
package maps

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

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

// scan reads and parses every decisions/*.md under root in sorted order.
func scan(root string) ([]string, map[string]fileResult) {
	dir := filepath.Join(root, "decisions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	var names []string
	results := map[string]fileResult{}
	for _, e := range entries {
		// A leading-dot file is skipped to match the shell glob decisions/*.md, which
		// without dotglob never expands a hidden name — so a `.foo.md` stays invisible to
		// the listing and the count (and cannot falsely nag the dashboard).
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		base := strings.TrimSuffix(e.Name(), ".md")
		names = append(names, base)
		results[base] = parseFile(content)
	}
	sort.Strings(names)
	return names, results
}

// Rows lists every unresolved ticket and close-readiness handoff row under
// root/decisions, in sorted file order.
func Rows(root string) [][]any {
	names, results := scan(root)
	var rows [][]any
	for _, name := range names {
		rows = append(rows, fileRows(name, results[name])...)
	}
	return rows
}

// UnresolvedCount is the DISTINCT not-close-ready map FILE count status surfaces:
// evaluated through the same predicate the listing uses, so the figure and the rows
// cannot drift. Two unresolved tickets in one file count as one.
func UnresolvedCount(root string) int {
	_, results := scan(root)
	n := 0
	for _, r := range results {
		if r.notCloseReady() {
			n++
		}
	}
	return n
}

// Command implements `bench maps`. `--count` is the status adapter's hook: it prints
// the distinct not-close-ready file count as a bare integer. It exists because that
// count is NOT recoverable from the row listing — a file with a file-scope marker but
// no `## #` ticket heading is not-close-ready yet emits no row — so the adapter reads
// the engine's own figure through the same parseFile, keeping one source of the rule.
func Command(args []string) (string, int) {
	switch {
	case len(args) == 0:
	case args[0] == "--count":
		root, err := git.Root()
		if err != nil {
			return "0\n", 0
		}
		return strconv.Itoa(UnresolvedCount(root)) + "\n", 0
	case args[0] == "-h" || args[0] == "--help":
		return "usage: bench maps [--count]\n", 0
	default:
		return toon.Usage("bench maps", args[0]) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	out, err := toon.TableTyped("maps", []string{"map", "ticket", "type", "state"}, Rows(root))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
}
