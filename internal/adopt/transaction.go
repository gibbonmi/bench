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

func stageBytes(dir, name string, data []byte, mode os.FileMode) (string, error) {
	path := filepath.Join(dir, name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return "", err
	}
	if _, err = f.Write(data); err == nil {
		err = f.Sync()
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(path)
		return "", err
	}
	return path, nil
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
	for _, c := range done {
		_ = os.Remove(c.backup)
	}
	return nil
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
