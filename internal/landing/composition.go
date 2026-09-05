package landing

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/landing/settlepolicy"
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
	// The refusal ticket 08 surfaces on Conflict is dropped here; today's message
	// carries the kind alone.
	tree, resolved, _, ok, err := resolveCaptureConflict(r.Root, out, records)
	if err != nil {
		return CompositionResult{}, err
	}
	if ok {
		return CompositionResult{Base: base, Tree: tree, Resolved: resolved}, nil
	}
	return CompositionResult{Base: base, Conflict: conflict}, nil
}

// resolveCaptureConflict applies the settle policy's verdict over a conflict. The
// merge's own written tree carries Git's conflict markers at the conflicted paths; each
// is replaced by the settled object, or removed when the verdict is a removal. A policy
// refusal, or a union verdict whose content Git cannot merge as text, leaves the whole
// conflict for the caller to refuse; the refusal rides back so the caller can name it.
func resolveCaptureConflict(root, mergeOutput string, records []settlepolicy.StageRecord) (string, []string, settlepolicy.Refusal, bool, error) {
	settlement := settlepolicy.Settle(records)
	if settlement.Refusal.Reason != "" || len(settlement.Verdicts) == 0 {
		return "", nil, settlement.Refusal, false, nil
	}
	// settled holds the object each path publishes; a nil entry publishes a removal.
	settled := make([]*settlepolicy.StageRecord, len(settlement.Verdicts))
	for i, verdict := range settlement.Verdicts {
		switch verdict.Kind {
		case settlepolicy.VerdictRemove:
			settled[i] = nil
		case settlepolicy.VerdictUnion:
			record, ok, err := unionStages(root, verdict.Path, verdict.Stages)
			if err != nil {
				return "", nil, settlepolicy.Refusal{}, false, err
			}
			if !ok {
				refusal := settlepolicy.Refuse(settlepolicy.ReasonUnionContentNotText, []string{verdict.Path})
				return "", nil, refusal, false, nil
			}
			settled[i] = record
		default:
			record := verdict.Record
			settled[i] = &record
		}
	}
	baseTree, err := mergeTreeResult(mergeOutput)
	if err != nil {
		return "", nil, settlepolicy.Refusal{}, false, err
	}
	tree, err := editTree(root, baseTree, func(idx string) error {
		for i, verdict := range settlement.Verdicts {
			record := settled[i]
			if record == nil {
				if err := indexRun(root, idx, "update-index", "--force-remove", "--", verdict.Path); err != nil {
					return fmt.Errorf("resolve %q: %w", verdict.Path, err)
				}
				continue
			}
			if err := indexRun(root, idx, "update-index", "--add", "--cacheinfo", record.Mode+","+record.OID+","+verdict.Path); err != nil {
				return fmt.Errorf("resolve %q: %w", verdict.Path, err)
			}
		}
		return nil
	})
	if err != nil {
		return "", nil, settlepolicy.Refusal{}, false, err
	}
	resolved := make([]string, 0, len(settlement.Verdicts))
	for _, verdict := range settlement.Verdicts {
		resolved = append(resolved, verdict.Path+":"+verdict.Side)
	}
	return tree, resolved, settlepolicy.Refusal{}, true, nil
}

// unionStages composes one union path from its merge-tree stages. With both sides
// present, the result is Git's own three-way union merge over the three stage blobs,
// with an absent merge base standing in as empty content. With one side absent, the
// present side's blob is the result, so a deletion cannot erase the other side's
// entries. A false ok is a refusal: Git could not merge the content as text.
func unionStages(root, path string, stages map[int]settlepolicy.StageRecord) (*settlepolicy.StageRecord, bool, error) {
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
			if content, err = blobContent(root, record.OID); err != nil {
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
	return &settlepolicy.StageRecord{Mode: source.Mode, OID: oid, Path: path, Stage: 3}, true, nil
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
func parseConflict(output string) (Conflict, []settlepolicy.StageRecord, error) {
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
	records := make([]settlepolicy.StageRecord, 0, separator-1)
	modes := make([]string, 0, separator-1)
	var paths []string
	seen := map[string]bool{}
	for _, raw := range parts[1:separator] {
		header, path, found := strings.Cut(string(raw), "\t")
		fields := strings.Fields(header)
		if !found || path == "" || len(fields) != 3 || len(fields[2]) != 1 || fields[2][0] < '1' || fields[2][0] > '3' {
			return Conflict{}, nil, errors.New("merge-tree returned malformed conflict record")
		}
		records = append(records, settlepolicy.StageRecord{Mode: fields[0], OID: fields[1], Stage: int(fields[2][0] - '0'), Path: path})
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
			kind = settlepolicy.ConflictKind(modes)
			if kind == "textual" {
				kind = "mode"
			}
		case "CONFLICT (contents)":
			kind = settlepolicy.ConflictKind(modes)
		}
		if kind != "" {
			return Conflict{Kind: kind, Paths: paths}, records, nil
		}
	}
	return Conflict{}, nil, errors.New("merge-tree returned an unrecognized conflict kind")
}
