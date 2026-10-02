package transaction

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
)

// Change names a prepared replacement. An empty Stage removes the destination.
type Change struct {
	Span        *Span
	Destination string
	Stage       string
	before      *image
}

// Prepare observes a destination before the adoption owner selects its replacement.
func Prepare(destination, stage string) (Change, error) {
	path, err := canonicalDestination(destination)
	if err != nil {
		return Change{}, err
	}
	if err := bounds.RefuseLinks(filepath.Dir(path)); err != nil {
		return Change{}, err
	}
	before, err := inspect(path)
	if err != nil {
		return Change{}, err
	}
	return Change{Destination: path, Stage: stage, before: &before}, nil
}

// Recheck refuses a destination whose observed identity no longer matches the plan.
func (c Change) Recheck() error {
	if c.before == nil {
		return fmt.Errorf("adoption change has no observation: %s", c.Destination)
	}
	if err := bounds.RefuseLinks(filepath.Dir(c.Destination)); err != nil {
		return err
	}
	current, err := inspect(c.Destination)
	if err != nil {
		return err
	}
	if !sameObserved(current, *c.before) {
		return fmt.Errorf("repair destination changed after inspection: %s", c.Destination)
	}
	return nil
}

// Store retains recovery records for one repository.
type Store struct {
	Directory     string
	Lease         *Lease
	SyncDirectory func(string) error
	Rename        func(string, string) error
}

// Apply publishes prepared managed replacements and retains their recovery state.
func (s Store) Apply(changes []Change) (id string, resultErr error) {
	paths := make([]string, len(changes))
	for i, change := range changes {
		paths[i] = change.Destination
	}
	lease := s.Lease
	if lease == nil {
		var err error
		lease, err = Lock(paths)
		if err != nil {
			return "", err
		}
		defer func() { resultErr = errors.Join(resultErr, lease.Close()) }()
	}
	if err := lease.covers(paths); err != nil {
		return "", err
	}
	entries := make([]entry, 0, len(changes))
	for _, change := range changes {
		path, err := canonicalDestination(change.Destination)
		if err != nil {
			return "", err
		}
		if err := bounds.RefuseLinks(filepath.Dir(path)); err != nil {
			return "", err
		}
		before, err := inspect(path)
		if err != nil {
			return "", err
		}
		if change.before != nil && !sameObserved(before, *change.before) {
			return "", fmt.Errorf("repair destination changed after inspection: %s", path)
		}

		after := image{Kind: "absent"}
		if change.Stage != "" {
			after, err = inspect(change.Stage)
			if err != nil {
				return "", err
			}
			if after.Kind == "absent" {
				return "", fmt.Errorf("missing staged replacement: %s", change.Stage)
			}
		}
		if same(before, after) {
			continue
		}
		span, err := managedSpan(change.Span, before, after)
		if err != nil {
			return "", err
		}
		entries = append(entries, entry{Destination: path, Before: before, After: after, Span: span})
	}
	if len(entries) == 0 {
		return "", nil
	}
	id, r, err := s.prepare(entries)
	if err != nil {
		return id, fmt.Errorf("repair backup publication failed: %w", err)
	}
	for _, e := range entries {
		current, err := inspect(e.Destination)
		if err == nil {
			err = bounds.RefuseLinks(filepath.Dir(e.Destination))
		}
		if err == nil && !sameObserved(current, e.Before) {
			err = fmt.Errorf("repair destination changed: %s", e.Destination)
		}
		if err == nil {
			err = s.publish(e.Destination, e.After)
		}
		if err != nil {
			saved, loadErr := s.load(id)
			if loadErr != nil {
				return id, errors.Join(err, loadErr)
			}
			return id, errors.Join(err, s.undo(id, saved))
		}
	}
	r.State = stateApplied
	if err := s.save(id, r); err != nil {
		return id, fmt.Errorf("repair incomplete; retain %s: %w", id, err)
	}
	return id, nil
}

// Undo restores the retained preimages of a completed or interrupted repair.
func (s Store) Undo(id string) (resultErr error) {
	r, err := s.load(id)
	if err != nil {
		return err
	}
	paths := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		paths[i] = e.Destination
	}
	lease, err := Lock(paths)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, lease.Close()) }()
	r, err = s.load(id)
	if err != nil {
		return err
	}
	return s.undo(id, r)
}

func (s Store) undo(id string, r record) error {
	for _, e := range r.Entries {
		if err := bounds.RefuseLinks(filepath.Dir(e.Destination)); err != nil {
			return err
		}
		current, err := inspect(e.Destination)
		if err != nil {
			return err
		}
		allowed := same(current, e.After)
		if r.State == statePrepared || r.State == stateUndoing || r.State == stateUndone {
			allowed = allowed || same(current, e.Before)
		}
		if !allowed {
			return fmt.Errorf("undo conflict: %s", e.Destination)
		}
		if same(current, e.After) {
			if _, err := restoredImage(e, current); err != nil {
				return err
			}
		}
	}
	r.State = stateUndoing
	if err := s.save(id, r); err != nil {
		return err
	}
	var failures []error
	for i := len(r.Entries) - 1; i >= 0; i-- {
		e := r.Entries[i]
		current, err := inspect(e.Destination)
		if err == nil && same(current, e.Before) {
			continue
		}
		if err != nil || !same(current, e.After) {
			failures = append(failures, fmt.Errorf("restore destination changed: %s", e.Destination))
			continue
		}
		if err := bounds.RefuseLinks(filepath.Dir(e.Destination)); err != nil {
			failures = append(failures, err)
			continue
		}
		restored, err := restoredImage(e, current)
		if err != nil {
			failures = append(failures, err)
			continue
		}
		if err := s.publish(e.Destination, restored); err != nil {
			failures = append(failures, fmt.Errorf("restore %s: %w", e.Destination, err))
		}
	}
	if err := errors.Join(failures...); err != nil {
		return fmt.Errorf("repair recovery incomplete; retain %s: %w", id, err)
	}
	r.State = stateUndone
	return s.save(id, r)
}

func (s Store) syncDirectory(path string) error {
	if s.SyncDirectory != nil {
		return s.SyncDirectory(path)
	}
	return SyncDirectory(path)
}

func (s Store) publish(path string, value image) error {
	rename := s.Rename
	if rename == nil {
		rename = os.Rename
	}
	return publish(path, value, s.syncDirectory, rename)
}

// Publish commits an adoption transaction without retaining successful preimages.
func Publish(changes []Change, options Store) error {
	dir, err := os.MkdirTemp("", "bench-adoption-")
	if err != nil {
		return err
	}
	options.Directory, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return err
	}
	id, err := options.Apply(changes)
	if err != nil {
		if id == "" {
			return errors.Join(err, os.RemoveAll(dir))
		}
		return fmt.Errorf("adoption incomplete; recovery directory %s, repair %s: %w", dir, id, err)
	}
	return os.RemoveAll(dir)
}
