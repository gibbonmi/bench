package structure

import (
	"bytes"
	"errors"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// The repository paths of the two reviewer-owned limit lists.
const (
	budgetsPath = ".bench/structure.budgets"
	acceptPath  = ".bench/structure-accept"
)

// limits is one read of the limit engine's configuration: the file and directory caps, the
// budget rows with their warnings, and the accept grants with their warnings and read fault.
type limits struct {
	maxLines, maxFiles int
	budgets            map[string]int
	budgetWarnings     []string
	accepts            map[string]string
	acceptWarnings     []string
	acceptErr          error
}

// Fact is one typed structure observation shared by query consumers.
type Fact struct {
	Kind, Path    string
	Actual, Limit int
	State, Detail string
}

// Facts returns the whole-tree engine's typed observations.
func Facts(root string) ([]Fact, error) {
	out, err := git.Output("-C", root, "ls-files")
	if err != nil {
		return nil, err
	}
	_, _, facts := evaluate(root, strings.Split(out, "\n"), true)
	return facts, nil
}

// checkAll is the human/count projection of the same whole-tree evaluation.
func checkAll(root string) (report string, violations int, err error) {
	out, err := git.Output("-C", root, "ls-files")
	if err != nil {
		return "", 0, err
	}
	report, violations, _ = evaluate(root, strings.Split(out, "\n"), true)
	return report, violations, nil
}

// Headroom is the typed sizing projection of one directory: the cap that the crowding rule
// applies to it, the source files it holds directly, and the accept reason that exempts it
// when Granted. Warnings holds each malformed budget row, which can hide the directory's own
// row.
type Headroom struct {
	Limit, Current int
	Grant          string
	Granted        bool
	Warnings       []string
}

// DirectoryHeadroom projects dir through the limits that the whole-tree scan reads at root.
// tracked is the caller's tracked path list, so the caller pins the source that it counts.
// An unreadable accept list answers an error, never an empty grant list.
func DirectoryHeadroom(root string, tracked []string, dir string) (Headroom, error) {
	l := readLimits(root)
	if l.acceptErr != nil {
		return Headroom{}, errors.New(acceptUnreadable(l.acceptErr))
	}
	key := dir + "/"
	headroom := Headroom{Limit: budgetFor(l.budgets, key, l.maxFiles), Warnings: l.budgetWarnings}
	for _, f := range filterSources(tracked) {
		if sourceDir(f) == dir {
			headroom.Current++
		}
	}
	headroom.Grant, headroom.Granted = acceptedReason(l.accepts, key)
	return headroom, nil
}

// Fits reports whether the directory takes net more source files without crowding, or a
// grant exempts it from the crowding rule.
func (h Headroom) Fits(net int) bool { return h.Granted || !crowded(h.Current+net, h.Limit) }

// SourceDirectory is the directory whose crowding count holds path, and whether the engine
// counts path as a source file at all.
func SourceDirectory(path string) (string, bool) {
	return sourceDir(path), len(filterSources([]string{path})) == 1
}

// LimitSources is the repository paths of the limit lists that DirectoryHeadroom reads from
// the checkout. A caller that pins a commit keeps them unchanged from that commit.
func LimitSources() []string { return []string{budgetsPath, acceptPath} }

// Lines is the physical line count that every length rule reads: the newline count, with
// wc -l semantics.
func Lines(content []byte) int { return bytes.Count(content, []byte{'\n'}) }

// sourceDir is the directory that the crowding rule counts a source file in: the path up to
// its last '/', or '.' when it has none.
func sourceDir(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[:i]
	}
	return "."
}

// crowded is the crowding rule itself: a source-file count over its directory cap.
func crowded(count, limit int) bool { return count > limit }
