package roadmap

import (
	"io/fs"
	"path/filepath"
	"unicode/utf8"
)

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
	OccurrenceKeys                    string
	OccurrenceCount                   int
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
type RetroFact struct {
	Path, State, Body string
	BodyBytes         int
	Truncated         bool
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
	Text                    string
	Rows                    []RoadmapRow
	Sequence                []SequenceRow
	SequenceText            string
	OccurrenceDiscrepancies []OccurrenceDiscrepancy
}

type OccurrenceDiscrepancy struct {
	Source, CaptureUnit, Kind, Owner, Incident string
	Structural                                 bool
}

type ContextSnapshot struct {
	Full            bool
	SequenceTrusted bool
	Sources         []SourceFact
	Roadmap         Document
	Ideas           []IdeaFact
	Learnings       []LearningFact
	Retros          []RetroFact
	Structure       [][]any
	Specs           [][]string
	SpecHistory     [][]string
	Git             [][]any
	GitChanges      [][]string
	GateCache       [][]any
	Failures        []ParseFailure
}

// dirBytes sums the sizes of the regular-file entries in a classified directory
// listing. A subdirectory contributes nothing: the figure reports the bytes of the
// records the listing itself names, not a recursive tree size. An entry whose Info
// call fails is skipped rather than failing the tally, because the directory's own
// state already carries whether the listing could be trusted.
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
