package roadmap

import (
	"fmt"
	"io/fs"
	"os"
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
type PromotionFact struct {
	Kind, Date, Scope, RoadmapIDs, Body string
	BodyBytes                           int
	Truncated                           bool
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
	Promotions  []PromotionFact
	Failures    []ParseFailure
}

func readSource(path string) ([]byte, string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "absent", nil
		}
		return nil, "", err
	}
	if !utf8.Valid(b) {
		return b, "malformed", fmt.Errorf("invalid UTF-8")
	}
	if len(b) == 0 {
		return b, "empty", nil
	}
	return b, "parsed", nil
}

func readDirSource(path string) ([]fs.DirEntry, string, int, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "absent", 0, nil
		}
		return nil, "", 0, err
	}
	if len(entries) == 0 {
		return entries, "empty", 0, nil
	}
	bytes := 0
	for _, e := range entries {
		if !e.IsDir() {
			if i, err := e.Info(); err == nil {
				bytes += int(i.Size())
			}
		}
	}
	return entries, "parsed", bytes, nil
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
