package gate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/spec"
)

type completionSourceKey struct{}

// WithCompletion binds the broker's reviewed source to prospective completion.
func WithCompletion(ctx context.Context, path, sourceTip string) context.Context {
	return context.WithValue(WithCheckpoint(ctx, Checkpoint{Spec: path, Complete: true}), completionSourceKey{}, sourceTip)
}

func (e *gateEvaluation) completionTree(graded *treeGeneration) (string, error) {
	sourceTree, err := benchgit.Output("-C", e.identityRoot, "rev-parse", "--verify", e.completionSource+"^{tree}")
	if err != nil {
		return "", err
	}
	source, err := captureProspectiveTree(e.identityRoot, sourceTree)
	if err != nil {
		return "", err
	}
	path := e.checkpoint.Spec
	original, err := benchgit.ReadTreeFile(e.identityRoot, sourceTree, path)
	if err != nil {
		return "", err
	}
	want, err := spec.Implemented(original)
	if err != nil {
		return "", err
	}
	got, err := benchgit.ReadTreeFile(e.identityRoot, graded.tree, path)
	if err != nil {
		return "", err
	}
	before, _ := source.entry(path)
	after, _ := graded.entry(path)
	if !bytes.Equal(want, got) || strings.Fields(before.Metadata)[0] != strings.Fields(after.Metadata)[0] {
		return "", fmt.Errorf("completion spec %s differs from the exact status transform; review the spec delta", path)
	}
	record, err := reviewrecord.RecordPath(path)
	if err != nil {
		return "", err
	}
	// Only after proving the exact transform can the spec differ between these snapshots.
	for _, pair := range [][2]*treeGeneration{{source, graded}, {graded, source}} {
		for _, entry := range pair[0].snapshot.entries {
			if entry.Path == path || entry.Path == record {
				continue
			}
			other, present := pair[1].entry(entry.Path)
			if !present || other.Metadata != entry.Metadata {
				return "", fmt.Errorf("completion composition changes %s; review the destination delta", entry.Path)
			}
		}
	}
	return sourceTree, nil
}

func validateCompletionContext(e *gateEvaluation) error {
	if e.completionSource == "" && e.checkpoint.Complete && e.prospective {
		return errors.New("prospective completion requires the reviewed source tip")
	}
	return nil
}
