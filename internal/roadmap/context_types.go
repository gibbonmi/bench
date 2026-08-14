package roadmap

import (
	"io/fs"
	"path/filepath"
)

type SourceFact struct {
	Source, State string
	Bytes         int
}
type RoadmapRow struct {
	ID, Title, Spec, SpecStatus, Body string
	ExternalTrigger                   bool
	BodyBytes                         int
	OccurrenceKeys                    string
	OccurrenceCount                   int
}
type SequenceRow struct {
	Rank          int
	Text, Command string
}
type IdeaFact struct {
	Date, Text string
	Body       string
	Line       int
	TextBytes  int
}
type LearningFact struct {
	Date, Title, State, Body string
	Line                     int
	BodyBytes                int
}
type RetroFact struct {
	Path, State, Body string
	BodyBytes         int
}
type ParseFailure struct {
	Source, Reason, Raw string
	RawBytes            int
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

type CaptureOccurrence struct {
	Owner, Incident, Source, CaptureUnit, State string
}

type OccurrencePair struct {
	Owner, Incident string
}

type ContextSnapshot struct {
	Full               bool
	SequenceTrusted    bool
	Sources            []SourceFact
	Roadmap            Document
	Ideas              []IdeaFact
	Learnings          []LearningFact
	Retros             []RetroFact
	CaptureOccurrences []CaptureOccurrence
	PendingOccurrences []OccurrencePair
	Structure          [][]any
	Specs              [][]string
	SpecHistory        [][]string
	Git                [][]any
	GitChanges         [][]string
	GateCache          [][]any
	Failures           []ParseFailure
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

func projectBody(s string, full bool) (string, int) {
	n := len([]byte(s))
	if full {
		return s, n
	}
	return "", n
}

func sourcePath(root, label string) string { return filepath.Join(root, filepath.FromSlash(label)) }
