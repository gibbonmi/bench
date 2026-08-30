package gate

// The lane's subject: what the composed tree changes against its base, and the lane the
// change list is graded under. The derivation lives here rather than beside the lane run
// because the run's own file is already at its structure budget.

import (
	"fmt"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// composedChangeFields is the field count of one raw-diff entry's metadata: the two
// modes, the two object IDs, and the status letter.
const composedChangeFields = 5

// ComposedChange is one entry of the raw diff between the base tree and the composed
// tree. Status is Git's own status letter. SrcMode and DstMode are the six-digit modes
// of the two sides, and one of them is `000000` when that side holds nothing. Path
// carries the file's own bytes, so a name with a space or a byte above ASCII survives.
type ComposedChange struct {
	Status  string
	SrcMode string
	DstMode string
	Path    string
}

// ComposedChanges lists what the tree changes against the base commit's tree. It is the
// one derivation of a lane's subject: a caller that named a directory reaches the files
// under it, and a caller that named nothing at all reaches an empty list.
//
// Rename detection is off, so a rename arrives as one deletion and one addition and each
// side classifies by its own path. The NUL framing is load-bearing: under the default
// `core.quotepath` a newline-framed name with a byte above ASCII arrives C-quoted, so a
// reader would carry a path no file has.
func ComposedChanges(root, base, tree string) ([]ComposedChange, error) {
	raw, err := benchgit.Raw("-C", root, "diff", "--raw", "--no-renames", "-z", base+"^{tree}", tree)
	if err != nil {
		return nil, fmt.Errorf("gate: composed change list unavailable: %w", err)
	}
	frames := strings.Split(string(raw), "\x00")
	var changes []ComposedChange
	for i := 0; i+1 < len(frames); i += 2 {
		change, err := parseComposedChange(frames[i], frames[i+1])
		if err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}
	return changes, nil
}

// parseComposedChange reads one `--raw -z` entry, whose metadata frame is
// `:<srcmode> <dstmode> <srcsha> <dstsha> <status>` and whose path is the frame after it.
func parseComposedChange(meta, path string) (ComposedChange, error) {
	fields := strings.Fields(strings.TrimPrefix(meta, ":"))
	if !strings.HasPrefix(meta, ":") || len(fields) != composedChangeFields || path == "" {
		return ComposedChange{}, fmt.Errorf("gate: unreadable composed change entry %q", meta)
	}
	return ComposedChange{Status: fields[4], SrcMode: fields[0], DstMode: fields[1], Path: path}, nil
}

// proseSubject answers the prose placeholder's paths: the changed Markdown the composed
// tree holds as a regular file. A deletion leaves the tree no file to grade, so it
// contributes no subject.
func proseSubject(changes []ComposedChange) []string {
	var paths []string
	for _, change := range changes {
		if strings.HasSuffix(change.Path, ".md") && regularFile(change.DstMode) {
			paths = append(paths, change.Path)
		}
	}
	return paths
}

func regularFile(mode string) bool {
	return mode == "100644" || mode == "100755"
}

// Lane is the fast lane a worktree commit runs in place of the whole-project gate.
// Checks is the declared check list. Kit is the source root the Bench-owned checks are
// built from, and it is empty when the graded root is the kit root itself, which selects
// the private checkout of the composed tree: the kit grades with its own code.
type Lane struct {
	Checks []Phase
	Kit    string
}

// LaneForCommit resolves the lane a worktree commit at root runs. It answers a nil lane
// for a root that declares none, so one read tells a caller both whether a lane exists
// and what it is. It applies the gate's own kit-root selection, so a caller outside this
// package asks the lane question once.
func LaneForCommit(root string) (*Lane, error) {
	kit := kitRoot(root)
	checks, err := LaneFor(root, kit)
	if err != nil || checks == nil {
		return nil, err
	}
	if sameDirectory(root, kit) {
		return &Lane{Checks: checks}, nil
	}
	return &Lane{Checks: checks, Kit: kit}, nil
}
