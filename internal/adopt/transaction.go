package adopt

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// stagedChange is the journal entry for one managed target.  stage is empty for
// a deletion.  Promotion always renames the old target aside before installing
// the staged replacement, which makes reversing a returned write error local and
// deterministic.
type stagedChange struct {
	rel, dest, stage, backup string
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

// changesModifyTree reports whether promoting changes would actually alter any file's
// content - as opposed to rewriting the same bytes back (transactionalLink always
// restages and repromotes every accepted entry regardless of content, so mtimes/inodes
// churn on every run even when nothing meaningful changed). It must run before
// promoteAll: it reads each destination's current fingerprint, which promoteAll is
// about to overwrite.
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
	created := map[string]bool{}
	done := make([]stagedChange, 0, len(changes))
	faultName := os.Getenv("BENCH_LINK_FAULT")
	fault, _ := strconv.Atoi(faultName)
	for i, c := range changes {
		dest, ok := changeDestination(root, c)
		if !ok {
			return fmt.Errorf("invalid managed path %s", c.rel)
		}
		if err := makeParents(filepath.Dir(dest), root, created); err != nil {
			rollback(root, done, created)
			return err
		}
		if _, err := os.Lstat(dest); err == nil {
			if err := os.Rename(dest, c.backup); err != nil {
				rollback(root, done, created)
				return err
			}
		}
		done = append(done, c)
		if (fault > 0 && i+1 == fault) || (faultName == "last" && i+1 == len(changes)) {
			rollback(root, done, created)
			return fmt.Errorf("injected link promotion fault %s", faultName)
		}
		if c.stage != "" {
			if err := os.Rename(c.stage, dest); err != nil {
				rollback(root, done, created)
				return err
			}
		}
	}
	if err := syncChangeParents(root, done); err != nil {
		rollback(root, done, created)
		syncChangeParentsBestEffort(root, done)
		return err
	}
	for _, c := range done {
		_ = os.Remove(c.backup)
	}
	return nil
}

func syncChangeParents(root string, changes []stagedChange) error {
	dirs, err := changeParentDirs(root, changes)
	if err != nil {
		return err
	}
	for dir := range dirs {
		if err := syncDirectory(dir); err != nil {
			return fmt.Errorf("sync destination directory %s: %w", dir, err)
		}
	}
	return nil
}

func syncChangeParentsBestEffort(root string, changes []stagedChange) {
	dirs, err := changeParentDirs(root, changes)
	if err != nil {
		return
	}
	for dir := range dirs {
		_ = syncDirectory(dir)
	}
}

func changeParentDirs(root string, changes []stagedChange) (map[string]bool, error) {
	dirs := map[string]bool{}
	for _, c := range changes {
		dest, ok := changeDestination(root, c)
		if !ok {
			return nil, fmt.Errorf("invalid managed path %s", c.rel)
		}
		dirs[filepath.Dir(dest)] = true
	}
	return dirs, nil
}

func syncDir(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()
	return f.Sync()
}

func rollback(root string, done []stagedChange, created map[string]bool) {
	for i := len(done) - 1; i >= 0; i-- {
		c := done[i]
		dest, ok := changeDestination(root, c)
		if !ok {
			continue
		}
		_ = os.Remove(dest)
		if _, err := os.Lstat(c.backup); err == nil {
			_ = os.Rename(c.backup, dest)
		}
	}
	dirs := make([]string, 0, len(created))
	for dir := range created {
		dirs = append(dirs, dir)
	}
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) > len(dirs[j]) })
	for _, dir := range dirs {
		_ = os.Remove(dir)
	}
}

func changeDestination(root string, c stagedChange) (string, bool) {
	if c.dest != "" {
		return c.dest, filepath.IsAbs(c.dest)
	}
	return resolveInside(root, c.rel)
}

func makeParents(dir, root string, created map[string]bool) error {
	var missing []string
	for d := dir; ; d = filepath.Dir(d) {
		if _, err := os.Stat(d); err == nil {
			break
		}
		missing = append(missing, d)
		if d == root || filepath.Dir(d) == d {
			break
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, d := range missing {
		created[d] = true
	}
	return nil
}
