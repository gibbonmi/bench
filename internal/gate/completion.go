package gate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
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
	original, err := e.completionFile(source, path)
	if err != nil {
		return "", err
	}
	want, err := spec.Implemented(original.data)
	if err != nil {
		return "", err
	}
	got, err := e.completionFile(graded, path)
	if err != nil {
		return "", err
	}
	if !bytes.Equal(want, got.data) || original.mode != got.mode {
		return "", fmt.Errorf("completion spec %s differs from the exact status transform; review the spec delta", path)
	}
	record, err := reviewrecord.RecordPath(path)
	if err != nil {
		return "", err
	}
	proven, err := e.completionClosure(graded, sourceTree, path)
	if err != nil {
		return "", err
	}
	proven[path], proven[record] = true, true
	// Only after proving the exact transform can a proven path differ between these snapshots.
	for _, pair := range [][2]*treeGeneration{{source, graded}, {graded, source}} {
		for _, entry := range pair[0].snapshot.entries {
			if proven[entry.Path] {
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

// completionClosure proves that graded carries exactly the verified closure that the
// commitment owner derives from the reviewed source tree, and returns each proven path.
// A path list cannot catch a kept sequence entry or an extra row removal.
func (e *gateEvaluation) completionClosure(graded *treeGeneration, sourceTree, spec string) (map[string]bool, error) {
	edits, err := commitrepo.Store{Root: e.identityRoot}.Closure(sourceTree, commitrepo.Delivery{Spec: spec, Source: e.completionSource})
	if err != nil {
		return nil, fmt.Errorf("completion closure: %w", err)
	}
	proven := map[string]bool{}
	for _, edit := range edits {
		proven[edit.Path] = true
		if edit.Delete {
			if _, present := graded.entry(edit.Path); present {
				return nil, fmt.Errorf("completion keeps closed %s; review the delivery closure", edit.Path)
			}
			continue
		}
		got, err := e.completionFile(graded, edit.Path)
		if err != nil || got.mode != edit.Mode || !bytes.Equal(got.data, edit.Data) {
			return nil, fmt.Errorf("completion %s differs from the exact closure transform; review the delivery closure", edit.Path)
		}
	}
	return proven, nil
}

func validateCompletionContext(e *gateEvaluation) error {
	if e.completionSource == "" && e.checkpoint.Complete && e.prospective {
		return errors.New("prospective completion requires the reviewed source tip")
	}
	return nil
}

type completionFile struct {
	data []byte
	mode string
}

// The captured entry supplies both mode and object identity for the bounded read.
func (e *gateEvaluation) completionFile(generation *treeGeneration, path string) (completionFile, error) {
	entry, present := generation.entry(path)
	fields := strings.Fields(entry.Metadata)
	if !present || !(benchgit.IndexEntry{Mode: fields[0]}).IsRegularFile() {
		return completionFile{}, fmt.Errorf("missing or nonregular completion file %s", path)
	}
	data, err := benchgit.ReadControlBlob(e.identityRoot, fields[2])
	return completionFile{data: data, mode: fields[0]}, err
}
