package landing

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CompositionRequest identifies two immutable commits to merge without checkout state.
type CompositionRequest struct {
	Root, Destination, Source, ReviewBase string
}

// CompositionResult is either a prospective tree or one bounded conflict kind.
// Resolved lists every capture path the composition policy settled, as
// "<path>:<side>", so the landing can disclose what the merge did not decide.
type CompositionResult struct {
	Base, Tree string
	Conflict   Conflict
	Resolved   []string
}

// Conflict describes why Git could not produce one prospective tree, and names
// every path it could not merge.
type Conflict struct {
	Kind  string
	Paths []string
}

// ConflictError is the refusal a conflicted reviewed landing returns. Its message is
// the bounded kind; the paths ride typed so the caller can render them.
type ConflictError struct{ Conflict }

func (e ConflictError) Error() string { return "composition conflict: " + e.Kind }

// stageRecord is one conflicted-file line of `merge-tree -z`: the mode and object of
// one stage (1 base, 2 destination, 3 source) at one path.
type stageRecord struct {
	mode, oid, path string
	stage           int
}

// unionStage is the CaptureSide stage of a path the rule table settles by union
// rather than by taking one side; no merge-tree stage carries that number.
const unionStage = 0

// CaptureSide is the rule table for a conflicted phase-owned path: it names the verb
// that settles the path and, for a take-a-side verb, the merge-tree stage that wins.
// The phase handoff is the source session's closing state, so the source wins it. The
// two append-only journals compose as the union of both sides, so no appended entry is
// lost. Every other capture file is the destination's running ledger, so the
// destination wins. A path outside the table has no rule, and the conflict stays a
// refusal.
func CaptureSide(path string) (stage int, side string, ok bool) {
	switch path {
	case "capture/session-handoff.md":
		return 3, "source", true
	case "capture/learnings.md", "capture/IDEAS.md":
		return unionStage, "union", true
	}
	if strings.HasPrefix(path, "capture/") {
		return 2, "destination", true
	}
	return 0, "", false
}

// Compose performs Git's three-way tree merge using the repository's real merge base.
// ReviewBase is metadata only and is never used as the merge base.
func (o Owner) Compose(r CompositionRequest) (CompositionResult, error) {
	if r.Root == "" || r.Destination == "" || r.Source == "" {
		return CompositionResult{}, errors.New("composition request is incomplete")
	}
	destination, err := compositionCommit(r.Root, r.Destination, "destination")
	if err != nil {
		return CompositionResult{}, err
	}
	source, err := compositionCommit(r.Root, r.Source, "source")
	if err != nil {
		return CompositionResult{}, err
	}
	base, err := output(r.Root, "merge-base", destination, source)
	if err != nil {
		return CompositionResult{}, fmt.Errorf("find merge base: %w", err)
	}
	out, err := mergeTree(r.Root, destination, source)
	if err == nil {
		tree, err := mergeTreeResult(out)
		if err != nil {
			return CompositionResult{}, err
		}
		return CompositionResult{Base: base, Tree: tree}, nil
	}
	conflict, records, parseErr := parseConflict(out)
	if parseErr != nil {
		return CompositionResult{}, parseErr
	}
	tree, resolved, ok, err := resolveCaptureConflict(r.Root, out, records)
	if err != nil {
		return CompositionResult{}, err
	}
	if ok {
		return CompositionResult{Base: base, Tree: tree, Resolved: resolved}, nil
	}
	return CompositionResult{Base: base, Conflict: conflict}, nil
}

// resolveCaptureConflict settles a conflict whose every path is a regular file the
// rule table names, by applying that path's verb. The merge's own written tree carries
// Git's conflict markers at those paths; each is replaced by the settled object, or
// removed when the settled verb takes a side that deleted the file. A conflict touching
// any path the table does not name, any non-regular stage, any mode disagreement, or
// any union content Git cannot merge as text is left for the caller to refuse.
func resolveCaptureConflict(root, mergeOutput string, records []stageRecord) (string, []string, bool, error) {
	stages := map[string]map[int]stageRecord{}
	var order []string
	for _, record := range records {
		if _, _, ok := CaptureSide(record.path); !ok {
			return "", nil, false, nil
		}
		if record.mode != "100644" && record.mode != "100755" {
			return "", nil, false, nil
		}
		if _, seen := stages[record.path]; !seen {
			stages[record.path] = map[int]stageRecord{}
			order = append(order, record.path)
		}
		stages[record.path][record.stage] = record
	}
	if len(order) == 0 {
		return "", nil, false, nil
	}
	// settled holds the object each path publishes; a nil entry publishes a removal.
	settled := map[string]*stageRecord{}
	for _, path := range order {
		// Two sides that disagree on the file mode leave no settled mode to publish,
		// so the conflict stays a refusal under every verb.
		destination, hasDestination := stages[path][2]
		source, hasSource := stages[path][3]
		if hasDestination && hasSource && destination.mode != source.mode {
			return "", nil, false, nil
		}
		if stage, _, _ := CaptureSide(path); stage != unionStage {
			if record, ok := stages[path][stage]; ok {
				settled[path] = &record
			} else {
				settled[path] = nil
			}
			continue
		}
		record, ok, err := unionStages(root, path, stages[path])
		if err != nil {
			return "", nil, false, err
		}
		if !ok {
			return "", nil, false, nil
		}
		settled[path] = record
	}
	baseTree, err := mergeTreeResult(mergeOutput)
	if err != nil {
		return "", nil, false, err
	}
	tree, err := editTree(root, baseTree, func(idx string) error {
		for _, path := range order {
			record := settled[path]
			if record == nil {
				if err := indexRun(root, idx, "update-index", "--force-remove", "--", path); err != nil {
					return fmt.Errorf("resolve %q: %w", path, err)
				}
				continue
			}
			if err := indexRun(root, idx, "update-index", "--add", "--cacheinfo", record.mode+","+record.oid+","+path); err != nil {
				return fmt.Errorf("resolve %q: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return "", nil, false, err
	}
	resolved := make([]string, 0, len(order))
	for _, path := range order {
		_, side, _ := CaptureSide(path)
		resolved = append(resolved, path+":"+side)
	}
	return tree, resolved, true, nil
}

// unionStages composes one union path from its merge-tree stages. With both sides
// present, the result is Git's own three-way union merge over the three stage blobs,
// with an absent merge base standing in as empty content. With one side absent, the
// present side's blob is the result, so a deletion cannot erase the other side's
// entries. A false ok is a refusal: Git could not merge the content as text.
func unionStages(root, path string, stages map[int]stageRecord) (*stageRecord, bool, error) {
	destination, hasDestination := stages[2]
	source, hasSource := stages[3]
	switch {
	case !hasDestination && !hasSource:
		return nil, false, nil
	case !hasSource:
		return &destination, true, nil
	case !hasDestination:
		return &source, true, nil
	}
	dir, err := os.MkdirTemp("", "bench-landing-union-")
	if err != nil {
		return nil, false, err
	}
	defer os.RemoveAll(dir)
	names := map[int]string{}
	for _, stage := range []int{1, 2, 3} {
		file := filepath.Join(dir, fmt.Sprintf("stage%d", stage))
		content := []byte(nil)
		if record, ok := stages[stage]; ok {
			if content, err = blobContent(root, record.oid); err != nil {
				return nil, false, err
			}
		}
		if err := os.WriteFile(file, content, 0o600); err != nil {
			return nil, false, err
		}
		names[stage] = file
	}
	merged, err := outputRaw(root, "merge-file", "-p", "--union", names[2], names[1], names[3])
	if err != nil {
		// merge-file refuses content it cannot read as text, so the union has no
		// result and the whole conflict returns to the caller as a refusal.
		return nil, false, nil
	}
	oid, err := hashBlob(root, merged)
	if err != nil {
		return nil, false, fmt.Errorf("write union of %q: %w", path, err)
	}
	return &stageRecord{mode: source.mode, oid: oid, path: path, stage: 3}, true, nil
}

func compositionCommit(root, value, role string) (string, error) {
	commit, err := output(root, "rev-parse", "--verify", value+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("composition %s is not a commit", role)
	}
	return commit, nil
}

func mergeTree(root, destination, source string) (string, error) {
	// merge-tree finds the merge base itself; this avoids checkout and index state.
	return outputCombined(root, "merge-tree", "--write-tree", "-z", destination, source)
}

func mergeTreeResult(output string) (string, error) {
	tree, _, found := strings.Cut(output, "\x00")
	if !found || len(tree) != 40 {
		return "", errors.New("merge-tree returned no tree")
	}
	return tree, nil
}

// parseConflict reads `merge-tree --write-tree -z` conflict output: the written tree,
// one stage record per conflicted file (`<mode> <object> <stage>\t<path>`), a NUL
// separator, then the informational messages that carry the conflict kind. It
// returns the bounded kind with every conflicted path, and the stage records.
func parseConflict(output string) (Conflict, []stageRecord, error) {
	parts := bytes.Split([]byte(output), []byte{0})
	separator := -1
	for i, part := range parts {
		if len(part) == 0 {
			separator = i
			break
		}
	}
	if separator < 1 {
		return Conflict{}, nil, errors.New("merge-tree returned no conflict records")
	}
	records := make([]stageRecord, 0, separator-1)
	modes := make([]string, 0, separator-1)
	var paths []string
	seen := map[string]bool{}
	for _, raw := range parts[1:separator] {
		header, path, found := strings.Cut(string(raw), "\t")
		fields := strings.Fields(header)
		if !found || path == "" || len(fields) != 3 || len(fields[2]) != 1 || fields[2][0] < '1' || fields[2][0] > '3' {
			return Conflict{}, nil, errors.New("merge-tree returned malformed conflict record")
		}
		records = append(records, stageRecord{mode: fields[0], oid: fields[1], stage: int(fields[2][0] - '0'), path: path})
		modes = append(modes, fields[0])
		if !seen[path] {
			seen[path] = true
			paths = append(paths, path)
		}
	}
	for _, record := range parts[separator+1:] {
		kind := ""
		switch string(record) {
		case "CONFLICT (modify/delete)":
			kind = "modify/delete"
		case "CONFLICT (rename/rename)":
			kind = "rename/rename"
		case "CONFLICT (directory/file)", "CONFLICT (file/directory)":
			kind = "file/directory"
		case "CONFLICT (distinct modes)":
			kind = contentConflictKind(modes)
			if kind == "textual" {
				kind = "mode"
			}
		case "CONFLICT (contents)":
			kind = contentConflictKind(modes)
		}
		if kind != "" {
			return Conflict{Kind: kind, Paths: paths}, records, nil
		}
	}
	return Conflict{}, nil, errors.New("merge-tree returned an unrecognized conflict kind")
}

func contentConflictKind(modes []string) string {
	for _, mode := range modes {
		if mode == "160000" {
			return "gitlink"
		}
		if mode == "120000" {
			return "symlink"
		}
	}
	for _, mode := range modes {
		for _, other := range modes {
			if mode != other {
				return "mode"
			}
		}
	}
	return "textual"
}
