package gate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"maps"
	"strings"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/spec"
)

type completionSourceKey struct{}
type completionFolderKey struct{}

// WithCompletion binds the broker's reviewed source to the prospective completion of the
// deliverable at path. A staged spec completes its checkpoint with its completion record.
// A tickets-only folder has no record, so its completion proves the folder closed instead.
func WithCompletion(ctx context.Context, path, sourceTip string) context.Context {
	if spec.IsLiveSpecPath(path) {
		ctx = WithCheckpoint(ctx, Checkpoint{Spec: path, Complete: true})
	} else {
		ctx = context.WithValue(ctx, completionFolderKey{}, path)
	}
	return context.WithValue(ctx, completionSourceKey{}, sourceTip)
}

// completionTree proves that graded is exactly the reviewed source after the broker
// transform, and returns the source tree. The deliverable's own proof comes first: the
// exact spec status with the record allowance, or the absent tickets-only folder. The
// verified closure follows, and every other path must keep its source bytes and mode.
func (e *gateEvaluation) completionTree(graded *treeGeneration) (string, error) {
	sourceTree, err := benchgit.Output("-C", e.identityRoot, "rev-parse", "--verify", e.completionSource+"^{tree}")
	if err != nil {
		return "", err
	}
	source, err := captureProspectiveTree(e.identityRoot, sourceTree)
	if err != nil {
		return "", err
	}
	path, proven := e.checkpoint.Spec, map[string]bool{}
	if e.completionFolder != "" {
		path = e.completionFolder
		err = closedFolder(source, graded, path, proven)
	} else {
		err = e.completedSpec(source, graded, path, proven)
	}
	if err != nil {
		return "", err
	}
	closed, err := e.completionClosure(graded, sourceTree, path)
	if err != nil {
		return "", err
	}
	maps.Copy(proven, closed)
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

// completedSpec proves that graded carries the exact implemented status of the source spec
// at path, and marks the spec and its completion record proven.
func (e *gateEvaluation) completedSpec(source, graded *treeGeneration, path string, proven map[string]bool) error {
	original, err := e.completionFile(source, path)
	if err != nil {
		return err
	}
	want, err := spec.Implemented(original.data)
	if err != nil {
		return err
	}
	got, err := e.completionFile(graded, path)
	if err != nil {
		return err
	}
	if !bytes.Equal(want, got.data) || original.mode != got.mode {
		return fmt.Errorf("completion spec %s differs from the exact status transform; review the spec delta", path)
	}
	record, err := reviewrecord.RecordPath(path)
	if err != nil {
		return err
	}
	proven[path], proven[record] = true, true
	return nil
}

// closedFolder proves that graded keeps no entry beneath the closed tickets-only folder,
// and marks each entry that the source folder holds proven.
func closedFolder(source, graded *treeGeneration, folder string, proven map[string]bool) error {
	under := func(path string) bool { return strings.HasPrefix(path, folder+"/") }
	for _, kept := range graded.snapshot.entries {
		if under(kept.Path) {
			return fmt.Errorf("completion keeps closed %s; review the delivery closure", kept.Path)
		}
	}
	for _, entry := range source.snapshot.entries {
		if under(entry.Path) {
			proven[entry.Path] = true
		}
	}
	return nil
}

// completionClosure proves that graded carries exactly the verified closure that the
// commitment owner derives from the reviewed source tree, and returns each proven path.
// A path list cannot catch a kept sequence entry or an extra row removal.
func (e *gateEvaluation) completionClosure(graded *treeGeneration, sourceTree, path string) (map[string]bool, error) {
	edits, err := commitrepo.Store{Root: e.identityRoot}.Closure(sourceTree, commitrepo.Delivery{Spec: path, Source: e.completionSource})
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
	if e.completionSource == "" && e.prospective && (e.checkpoint.Complete || e.completionFolder != "") {
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
