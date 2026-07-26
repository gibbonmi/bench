package roadmap

import (
	"io/fs"
	"path/filepath"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/bounds"
)

// controlRecordLimit is the size bound every roadmap control-record read applies.
// It reuses bounds.OutlineFileLimit rather than declaring a second constant: both
// bound a whole small text file read into memory, and bounds is the one package
// that owns read-size policy.
const controlRecordLimit = bounds.OutlineFileLimit

const contextBodyLimit = 4096

type SourceFact struct {
	Source, State string
	Bytes         int
}
type RoadmapRow struct {
	ID, Title, Spec, SpecStatus, Body string
	ExternalTrigger                   bool
	BodyBytes                         int
	Truncated                         bool
}
type SequenceRow struct {
	Rank          int
	Text, Command string
}
type IdeaFact struct {
	Date, Text string
	TextBytes  int
	Truncated  bool
}
type LearningFact struct {
	Date, Title, State, Body string
	BodyBytes                int
	Truncated                bool
}
type ParseFailure struct {
	Source, Reason, Raw string
	RawBytes            int
	Truncated           bool
}

type GateCacheFact struct {
	Present                                                       bool
	State, PendingStatus, Status, CachedTree, WorkTree, Timestamp string
	Stale                                                         bool
	CacheBytes                                                    int
}

// Document is the typed roadmap projection shared by the human renderer and context.
type Document struct {
	Text         string
	Rows         []RoadmapRow
	Sequence     []SequenceRow
	SequenceText string
}

type ContextSnapshot struct {
	Full        bool
	Sources     []SourceFact
	Roadmap     Document
	Ideas       []IdeaFact
	Learnings   []LearningFact
	Structure   [][]any
	Specs       [][]string
	SpecHistory [][]string
	Git         [][]any
	GitChanges  [][]string
	GateCache   [][]any
	Failures    []ParseFailure
}

// dirBytes sums the regular-file entries a classified directory listing carries —
// the same tally readDirSource used to compute inline, kept as its own function now
// that ClassifyDir owns the read itself.
func dirBytes(entries []fs.DirEntry) int {
	total := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if info, err := e.Info(); err == nil {
			total += int(info.Size())
		}
	}
	return total
}

func limited(s string, full bool) (string, int, bool) {
	n := len([]byte(s))
	if full || n <= contextBodyLimit {
		return s, n, false
	}
	b := []byte(s)[:contextBodyLimit]
	for !utf8.Valid(b) {
		b = b[:len(b)-1]
	}
	return string(b), n, true
}

func sourcePath(root, label string) string { return filepath.Join(root, filepath.FromSlash(label)) }
