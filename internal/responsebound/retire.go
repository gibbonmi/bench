package responsebound

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// retiringVerbs are the `bench worktree` leaves that retire an assignment. Each one
// spills to the primary scope, so the retirement cannot remove a spill that is still open.
var retiringVerbs = []string{"release", "clean", "reclaim", "land"}

// retiring reports that argv calls a retiring verb.
func retiring(argv []string) bool {
	return len(argv) > 1 && argv[0] == "worktree" && slices.Contains(retiringVerbs, argv[1])
}

// Drop removes the spill directory of one retired assignment. The retirement path calls
// it beside the census drop, and a retirement that runs twice completes both times. An
// identifier that is not an assignment id is refused rather than composed into a path,
// because the removal is unrecoverable. A store level that is not a real directory holds
// no spill of this owner, so the drop never follows a symlink.
func Drop(home, root, assignment string) error {
	if _, ok := poolkey.SplitAssignmentSegment(poolkey.AssignmentSegment(assignment, assignment)); !ok {
		return fmt.Errorf("response spill assignment id is malformed: %s", sanitize.Controls(assignment))
	}
	dir := home
	for _, component := range []string{storeDirName, poolkey.Key(root)} {
		dir = filepath.Join(dir, component)
		if err := realDir(dir); errors.Is(err, fs.ErrNotExist) {
			return nil
		} else if err != nil {
			return fmt.Errorf("remove response spills: %w", err)
		}
	}
	if err := os.RemoveAll(filepath.Join(dir, assignment)); err != nil {
		return fmt.Errorf("remove response spills: %w", err)
	}
	return nil
}

// storedSpill is one spill file of a scope and the creation time its name carries.
type storedSpill struct {
	created int64
	name    string
}

// prunePrimary removes the oldest spills of one primary scope until primaryRetained
// remain. The removal is best effort: a failure keeps an older file and never touches the
// new spill. An entry that is not a spill file stays.
func prunePrimary(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var spills []storedSpill
	for _, entry := range entries {
		if created, ok := spillCreated(entry.Name()); ok && entry.Type().IsRegular() {
			spills = append(spills, storedSpill{created: created, name: entry.Name()})
		}
	}
	if len(spills) <= primaryRetained {
		return
	}
	slices.SortFunc(spills, func(a, b storedSpill) int {
		return cmp.Or(cmp.Compare(a.created, b.created), strings.Compare(a.name, b.name))
	})
	for _, spill := range spills[:len(spills)-primaryRetained] {
		_ = os.Remove(filepath.Join(dir, spill.name))
	}
}

// spillCreated reads the creation time that spillName puts before the first hyphen.
func spillCreated(name string) (int64, bool) {
	stem, ok := strings.CutSuffix(name, spillSuffix)
	if !ok {
		return 0, false
	}
	prefix, _, ok := strings.Cut(stem, "-")
	if !ok {
		return 0, false
	}
	created, err := strconv.ParseInt(prefix, 10, 64)
	return created, err == nil
}
