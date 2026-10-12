package reviewrecord

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/commentgap"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// FoldDelta is the part of a fold entry that the caller derives from the two parents of
// the fold. This package cannot import the landing composition or the test package
// index, so the record command derives the delta.
type FoldDelta struct {
	Resolution       []string
	ResolutionDigest string
	TestPackages     []string
}

// FoldParents returns the first parent of commit and the default-branch commit that it
// merges. Any other commit refuses: one that is not a two-parent merge, or one whose second
// parent is outside the history of the default-branch tip. A failed default-branch query
// refuses and names the failure.
func FoldParents(root, commit string) (from, incoming string, err error) {
	line, err := benchgit.Output("-C", root, "rev-list", "--parents", "-n", "1", commit)
	parents := strings.Fields(line)
	if err != nil || len(parents) != 3 {
		return "", "", fmt.Errorf("fold %s is not a two-parent merge", commit)
	}
	contains, err := benchgit.DefaultContains(root, parents[2])
	if err != nil {
		return "", "", fmt.Errorf("fold %s: %w", commit, err)
	}
	if !contains {
		return "", "", fmt.Errorf("fold %s does not merge a default-branch commit", commit)
	}
	return parents[1], parents[2], nil
}

// RecordFold appends the fold entry of commit, a full commit ID, after the last recorded
// chunk. The first parent of commit must follow the covered digest, the source digest of
// the last chunk or of the last fold after it. derive supplies the delta after every
// refusal passes. The answer is the entry as written.
func RecordFold(root, spec, commit string, derive func(from, incoming string) (FoldDelta, error)) (Fold, error) {
	digest, plan, err := atSource(root, spec, commit)
	if err != nil {
		return Fold{}, err
	}
	from, incoming, err := FoldParents(root, commit)
	if err != nil {
		return Fold{}, err
	}
	var entry Fold
	err = write(root, spec, absent(""), func(record *Record) error {
		if err := currentPlan(*record, plan, commit); err != nil {
			return err
		}
		after, covered := "", ""
		if len(record.Chunks) != 0 {
			chunk := record.Chunks[len(record.Chunks)-1]
			after = chunk.ID
			covered, _ = coveredAfter(root, *record, chunk)
		}
		if !follows(root, record.Spec, from, covered) {
			return fmt.Errorf("fold %s does not follow the reviewed source", commit)
		}
		delta, err := derive(from, incoming)
		if err != nil {
			return err
		}
		entry = Fold{ID: foldID(len(record.Folds) + 1), After: after, From: from, Commit: commit, Incoming: incoming, SourceDigest: digest,
			Resolution: append([]string{}, delta.Resolution...), ResolutionDigest: delta.ResolutionDigest, TestPackages: append([]string{}, delta.TestPackages...), Reviews: []Review{}}
		record.Folds = append(record.Folds, entry)
		return nil
	})
	return entry, err
}

// foldID is the ID of the nth fold entry of a record, counted from 1.
func foldID(n int) string { return fmt.Sprintf("F%d", n) }

// coveringFold returns the last fold after chunk when each fold after chunk follows the
// digest before it, the chunk digest first. It returns nil when no fold follows chunk or
// the folds break that chain.
func coveringFold(root string, record Record, chunk Chunk) *Fold {
	var last *Fold
	covered := chunk.SourceDigest
	for i := range record.Folds {
		fold := &record.Folds[i]
		if fold.After != chunk.ID {
			continue
		}
		if !follows(root, record.Spec, fold.From, covered) {
			return nil
		}
		last, covered = fold, fold.SourceDigest
	}
	return last
}

// coveredAfter returns the source digest and the commit that the chain covers after chunk:
// those of the last fold after chunk, or the chunk's own.
func coveredAfter(root string, record Record, chunk Chunk) (digest, commit string) {
	if fold := coveringFold(root, record, chunk); fold != nil {
		return fold.SourceDigest, fold.Commit
	}
	return chunk.SourceDigest, chunk.Tip
}

// foldCoversBase reports whether the last fold after chunk covers base, the next chunk
// base, whose source digest is digest. The fold commit must be an ancestor of base, and
// base must hold the fold's source digest.
func foldCoversBase(root string, record Record, chunk Chunk, base, digest string) bool {
	fold := coveringFold(root, record, chunk)
	return fold != nil && fold.SourceDigest == digest && benchgit.OK("-C", root, "merge-base", "--is-ancestor", fold.Commit, base)
}

// uncoveredFold returns the newest fold of a default-branch commit in the range from..tip,
// or "" when the range holds none.
func uncoveredFold(root, from, tip string) string {
	merges, err := benchgit.Output("-C", root, "rev-list", "--merges", from+".."+tip)
	if err != nil {
		return ""
	}
	for _, merge := range strings.Fields(merges) {
		if _, _, err := FoldParents(root, merge); err == nil {
			return merge
		}
	}
	return ""
}

// follows reports whether the source digest of commit equals covered or differs from it
// only by a comment-only gap.
func follows(root, spec, commit, covered string) bool {
	tree, err := sourceTree(root, commit)
	if err != nil {
		return false
	}
	digest, err := SourceDigest(root, tree, spec)
	return err == nil && commentgap.Prove(root, covered, digest) == nil
}
