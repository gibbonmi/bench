package gate

// The reader facts of an inspection: whether the current tree holds a reader record, and how
// many reader landings follow the green marker. A consumer reads them to name a reader green
// that the verdict cache cannot show.

import (
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// ReaderLandingBound is the landing plan's bound of first-parent commits after the green
// marker, which a reader-green count is read against.
const ReaderLandingBound = readerChainBound

// withReaderFacts adds the reader facts of root to gi. A store that cannot be read holds no
// reader record.
func withReaderFacts(root string, gi Inspection) Inspection {
	if recorded, err := ReaderRecorded(root, gi.CurrentTree); err != nil || !recorded {
		return gi
	}
	gi.ReaderRecorded = true
	gi.ReaderLandings = readerLandings(root)
	return gi
}

// readerLandings counts the first-parent commits after the green marker of root's default
// branch, up to HEAD, whose tree holds a reader record. A marker or a history that cannot be
// read counts none.
func readerLandings(root string) int {
	branch, ok := benchgit.ResolvedDefault(root)
	if !ok {
		return 0
	}
	marker, present, err := greenmarker.Read(root, branch)
	if err != nil || !present {
		return 0
	}
	chain, err := benchgit.FirstParentsAfter(root, marker, "HEAD")
	if err != nil {
		return 0
	}
	count := 0
	for _, commit := range chain {
		tree, err := commitTree(root, commit)
		if err != nil {
			continue
		}
		if recorded, err := ReaderRecorded(root, tree); err == nil && recorded {
			count++
		}
	}
	return count
}
