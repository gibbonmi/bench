package adopt

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/adopt/transaction"
)

// stagedChange identifies a managed replacement, or a deletion when stage is empty.
type stagedChange struct {
	rel, dest, stage string
}

var syncDirectory = syncDir

func stageBytes(dir, name string, data []byte, mode os.FileMode) (string, error) {
	path := filepath.Join(dir, name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return "", err
	}
	if err := writeSyncClose(path, f, data); err != nil {
		return "", err
	}
	return path, nil
}

// writeSyncClose is the one durable staged-write lifecycle. Callers choose where
// their stage file lives, while this owns writing, syncing, closing, and cleanup.
func writeSyncClose(path string, f *os.File, data []byte) error {
	var err error
	if _, err = f.Write(data); err == nil {
		err = f.Sync()
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(path)
		return err
	}
	return nil
}

// changesModifyTree compares planned fingerprints before publication for setup's result.
func changesModifyTree(root string, changes []stagedChange) (bool, error) {
	for _, c := range changes {
		dest, ok := changeDestination(root, c)
		if !ok {
			return false, fmt.Errorf("invalid managed path %s", c.rel)
		}
		before, _ := fingerprintPath(dest)
		after := ""
		if c.stage != "" {
			fp, err := fingerprintPath(c.stage)
			if err != nil {
				return false, err
			}
			after = fp
		}
		if before != after {
			return true, nil
		}
	}
	return false, nil
}

func promoteAll(root string, changes []stagedChange) error {
	return promoteWithLease(root, changes, nil)
}

func promoteWithLease(root string, changes []stagedChange, guard *adoptionGuard) error {
	prepared := make([]transaction.Change, 0, len(changes))
	parents := map[string]bool{}
	destinations := map[string]int{}
	for i, c := range changes {
		dest, ok := changeDestination(root, c)
		if !ok {
			return fmt.Errorf("invalid managed path %s", c.rel)
		}
		change, err := observedChange(dest, c.stage, guard)
		if err != nil {
			return err
		}
		if guard != nil && guard.repair != nil {
			change.Span = guard.repair.spans[dest]
		}
		prepared = append(prepared, change)
		parents[filepath.Dir(dest)] = true
		destinations[dest] = i + 1
	}
	faultName := os.Getenv("BENCH_LINK_FAULT")
	fault, _ := strconv.Atoi(faultName)
	interrupt := 0
	if value, ok := strings.CutPrefix(faultName, "interrupt:"); ok {
		interrupt, _ = strconv.Atoi(value)
	}
	var lease *transaction.Lease
	if guard != nil {
		lease = guard.Lease
	}
	options := transaction.Store{
		Lease: lease,
		SyncDirectory: func(dir string) error {
			if parents[dir] {
				return syncDirectory(dir)
			}
			return syncDir(dir)
		},
		Rename: func(old, new string) error {
			if index := destinations[new]; index > 0 && index == interrupt {
				// Abrupt exit leaves the durable journal for a fresh-process recovery test.
				os.Exit(97)
			}
			if index := destinations[new]; index > 0 && ((fault > 0 && index == fault) || (faultName == "last" && index == len(changes))) {
				return fmt.Errorf("injected link promotion fault %s", faultName)
			}
			return os.Rename(old, new)
		},
	}
	if guard != nil && guard.repair != nil {
		options.Directory = guard.repair.store.Directory
		id, err := options.Apply(prepared)
		guard.repair.id = id
		return err
	}
	return transaction.Publish(prepared, options)
}

func syncDir(dir string) error {
	return transaction.SyncDirectory(dir)
}

func changeDestination(root string, c stagedChange) (string, bool) {
	if c.dest != "" {
		return c.dest, filepath.IsAbs(c.dest)
	}
	return resolveInside(root, c.rel)
}

type adoptionObservation struct {
	change transaction.Change
	err    error
}

type adoptionGuard struct {
	*transaction.Lease
	observations map[string]adoptionObservation
	repair       *repairRun
}

func lockAdoption(root string, plan []planEntry) (*adoptionGuard, error) {
	manifest := filepath.Join(root, ".bench", "link-manifest.tsv")
	observed, err := transaction.Prepare(manifest, "")
	if err != nil {
		return nil, err
	}
	paths, err := adoptionPaths(root, plan)
	if err != nil {
		return nil, err
	}
	guard, err := lockObserved(paths)
	if err != nil {
		return nil, err
	}
	if err := observed.Recheck(); err != nil {
		guard.Close()
		return nil, err
	}
	return guard, nil
}

func lockObserved(paths []string) (*adoptionGuard, error) {
	lease, err := transaction.Lock(paths)
	if err != nil {
		return nil, err
	}
	guard := &adoptionGuard{Lease: lease, observations: map[string]adoptionObservation{}}
	for _, path := range paths {
		change, err := transaction.Prepare(path, "")
		guard.observations[path] = adoptionObservation{change, err}
	}
	return guard, nil
}

func adoptionPaths(root string, plan []planEntry) ([]string, error) {
	paths := []string{filepath.Join(root, ".bench", "link-manifest.tsv"), filepath.Join(root, "AGENTS.md"), filepath.Join(root, "CLAUDE.md")}
	for _, e := range plan {
		paths = append(paths, filepath.Join(root, filepath.FromSlash(e.rel)))
	}
	old, err := ReadManifest(paths[0])
	if err != nil {
		return nil, err
	}
	for _, row := range old.Rows() {
		if dest, ok := resolveInside(root, row.rel); ok {
			paths = append(paths, dest)
		}
	}
	hooks, err := hooksDir(root)
	if err != nil {
		return nil, err
	}
	return append(paths, filepath.Join(hooks, "pre-push")), nil
}

func observedChange(dest, stage string, guard *adoptionGuard) (transaction.Change, error) {
	if guard == nil {
		return transaction.Prepare(dest, stage)
	}
	observed, ok := guard.observations[dest]
	if !ok {
		return transaction.Change{}, fmt.Errorf("adoption destination absent from the inspected plan: %s", dest)
	}
	change := observed.change
	change.Stage = stage
	return change, observed.err
}

func publishObserved(dest, stage string, guard *adoptionGuard) error {
	change, err := observedChange(dest, stage, guard)
	if err != nil {
		return err
	}
	return publishChange(change, guard)
}

func publishChange(change transaction.Change, guard *adoptionGuard) error {
	if guard != nil && guard.repair != nil {
		return retainDoctorChange(change, guard.repair)
	}
	options := transaction.Store{}
	if guard != nil {
		options.Lease = guard.Lease
	}
	return transaction.Publish([]transaction.Change{change}, options)
}
