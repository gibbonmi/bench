package repository

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/roadmap"
)

// stage replaces the policy and projects its sequence onto the board. It returns the
// restore of both files' snapshots, which the ledger transaction runs when the receipt
// write leaves the old ledger. A failed stage restores before it returns, because a
// replacement can fail after its rename has made the new bytes visible.
func (store Store) stage(replacer durablefile.Replacer, plan commitment.Plan) (func() error, error) {
	policyPath := filepath.Join(store.Root, filepath.FromSlash(commitment.PolicyPath))
	roadmapPath := filepath.Join(store.Root, roadmap.RoadmapFile)
	policyBefore, err := snapshotFile(policyPath, true)
	if err != nil {
		return nil, err
	}
	roadmapBefore, err := snapshotFile(roadmapPath, true)
	if err != nil {
		return nil, err
	}
	policy, err := commitment.Parse(plan.Proposed)
	if err != nil {
		return nil, err
	}
	var projected []byte
	if roadmapBefore.present {
		projected, err = roadmap.ProjectSequence(roadmapBefore.data, commitment.Remaining(policy))
		if err != nil {
			return nil, err
		}
	}
	restore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }
	if err := replace(replacer, policyPath, plan.Proposed, 0o644); err != nil {
		return nil, errors.Join(err, restore())
	}
	// A project with no board has no sequence to project, so approval stages the policy alone.
	if !roadmapBefore.present {
		return restore, nil
	}
	if err := replace(replacer, roadmapPath, projected, roadmapBefore.mode); err != nil {
		return nil, errors.Join(err, restore())
	}
	return restore, nil
}

// replace creates the parent directory that the replacement requires.
func replace(replacer durablefile.Replacer, path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return replacer.Replace(path, data, mode)
}

// restoreAll attempts every snapshot, so one failed file does not leave the other staged.
func restoreAll(replacer durablefile.Replacer, snapshots ...fileSnapshot) error {
	var errs []error
	for _, snapshot := range snapshots {
		if err := snapshot.restore(replacer); err != nil {
			errs = append(errs, fmt.Errorf("commitment approval recovery incomplete: restore %s: %w", snapshot.path, err))
		}
	}
	return errors.Join(errs...)
}

type fileSnapshot struct {
	path    string
	data    []byte
	mode    os.FileMode
	present bool
}

func snapshotFile(path string, allowAbsent bool) (fileSnapshot, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) && allowAbsent {
		return fileSnapshot{path: path, mode: 0o644}, nil
	}
	if err != nil {
		return fileSnapshot{}, fmt.Errorf("read commitment output %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fileSnapshot{}, fmt.Errorf("read commitment output %s: not a regular file", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fileSnapshot{}, fmt.Errorf("read commitment output %s: %w", path, err)
	}
	return fileSnapshot{path: path, data: data, mode: info.Mode().Perm(), present: true}, nil
}

// restore counts as complete only when a read-back matches the snapshot, because a
// replacement can return success while another writer changes the file.
func (snapshot fileSnapshot) restore(replacer durablefile.Replacer) error {
	if !snapshot.present {
		if err := os.Remove(snapshot.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else if err := replace(replacer, snapshot.path, snapshot.data, snapshot.mode); err != nil {
		return err
	}
	current, err := snapshotFile(snapshot.path, true)
	if err != nil {
		return err
	}
	if current.present != snapshot.present || !bytes.Equal(current.data, snapshot.data) || current.mode != snapshot.mode {
		return errors.New("the read-back does not match the snapshot")
	}
	return nil
}
