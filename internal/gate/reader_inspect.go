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
// reader record. Without a reader landing count, gi keeps no reader facts, so its stale row
// still prints.
func withReaderFacts(root string, gi Inspection) Inspection {
	if recorded, err := ReaderRecorded(root, gi.CurrentTree); err != nil || !recorded {
		return gi
	}
	landings, counted := readerLandings(root)
	if !counted {
		return gi
	}
	gi.ReaderRecorded, gi.ReaderLandings = true, landings
	return gi
}

// readerLandings counts the reader landings of the landing plan's chain from the green marker
// of root's default branch up to HEAD. It answers false when the marker is absent, when the
// marker is not an ancestor of HEAD, or when a read fails.
func readerLandings(root string) (int, bool) {
	branch, ok := benchgit.ResolvedDefault(root)
	if !ok {
		return 0, false
	}
	marker, present, err := greenmarker.Read(root, branch)
	if err != nil || !present {
		return 0, false
	}
	if ancestor, err := benchgit.IsAncestor(root, marker, "HEAD"); err != nil || !ancestor {
		return 0, false
	}
	_, landings, err := readerChain(root, landingValue{destination: "HEAD", marker: marker})
	if err != nil {
		return 0, false
	}
	return len(landings), true
}
