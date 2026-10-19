package preflight

import (
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/tickets"
)

// probe records the tree facts of one `Writes:` entry: whether its path exists and
// whether it joins the system suite. The policy over both belongs to Decide. A (new) test
// file the tree does not hold yet takes its system tag from its directory, so build
// preflight grades the kit pin before the build creates the file.
func (facts *ticketFacts) probe(root, entry string) {
	path, isNew := tickets.WritesPath(entry)
	facts.writes[entry] = treeHolds(root, path)
	if systemTagged(root, path) || isNew && !facts.writes[entry] && newSystemTest(root, path) {
		facts.systemTag[entry] = true
	}
}

// deletedBefore reports, for each (deleted) entry of the parsed tickets, whether its path
// is a tracked file at base. A deletion holds its authority from that before-state, so the
// tip can hold the file before the implementation and lack it after.
func deletedBefore(root, base string, parsed []tickets.Ticket) map[string]bool {
	before := map[string]bool{}
	for _, ticket := range parsed {
		for _, entry := range ticket.Writes {
			if write, err := tickets.ParseWrite(entry); err == nil && write.Kind == tickets.WriteDeleted && base != "" {
				kind, err := git.Output("-C", root, "cat-file", "-t", base+":"+write.Path)
				before[entry] = err == nil && kind == "blob"
			}
		}
	}
	return before
}

// buildDeletedPaths is the exact D-status subset of the build's landing range, from
// landingBase to tip. The landing base is the review base that `bench diff` resolves, so
// it moves forward with each fold. Only a recorded chunk makes the source a build
// with earlier chunks, so a record that does not read or holds no chunk answers none, and
// so does a range that does not resolve.
func buildDeletedPaths(root, specPath, landingBase, tip string) []string {
	if landingBase == "" || tip == "" {
		return nil
	}
	if record, err := reviewrecord.ReadTree(root, tip, specPath); err != nil || len(record.Chunks) == 0 {
		return nil
	}
	landing, kind, _ := diff.ResolveSourceRange(root, landingBase, tip)
	if kind != "" {
		return nil
	}
	return landing.DeletedPaths
}

// treeHolds reports whether one `Writes:` path names something in the tree.
// Any file type answers yes: the row grades whether the path exists, not
// what sits at it. The lstat never follows a link, so a dangling symlink
// reads as the absent path it is.
func treeHolds(root, path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path)))
	return err == nil
}
