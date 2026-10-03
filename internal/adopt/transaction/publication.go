package transaction

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
)

type publication struct {
	path, temporary string
	value           image
}

func stagePublication(path string, value image, persist func(string) error) (publication, error) {
	p := publication{path: path, value: value}
	parent := filepath.Dir(path)
	if err := bounds.RefuseLinks(parent); err != nil {
		return p, err
	}
	if err := ensureDirectory(parent, 0o755, persist); err != nil {
		return p, err
	}
	if value.Kind == "absent" {
		return p, nil
	}
	f, err := os.CreateTemp(parent, ".bench-adopt-")
	if err != nil {
		return p, err
	}
	p.temporary = f.Name()
	failed := true
	defer func() {
		if failed {
			p.close()
		}
	}()
	if value.Kind == "symlink" {
		if err := f.Close(); err != nil {
			return p, err
		}
		if err := os.Remove(p.temporary); err != nil {
			return p, err
		}
		if err := os.Symlink(value.Link, p.temporary); err != nil {
			return p, err
		}
	} else if value.Kind == "file" {
		_, err = f.Write(value.Data)
		if err == nil {
			err = f.Chmod(value.Mode)
		}
		if err == nil {
			err = f.Sync()
		}
		if err = errors.Join(err, f.Close()); err != nil {
			return p, err
		}
	} else {
		_ = f.Close()
		return p, fmt.Errorf("unknown image kind %q", value.Kind)
	}
	p.value, err = inspect(p.temporary)
	if err != nil {
		return p, err
	}
	if p.value.Identity == nil {
		return p, fmt.Errorf("replacement identity unavailable: %s", path)
	}
	failed = false
	return p, nil
}

func (p publication) close() {
	if p.temporary != "" {
		_ = os.Remove(p.temporary)
	}
}

func (p publication) publish(s Store, expected *image) error {
	if expected != nil {
		if err := bounds.RefuseLinks(filepath.Dir(p.path)); err != nil {
			return err
		}
		current, err := inspect(p.path)
		if err != nil {
			return err
		}
		if !sameObserved(current, *expected) {
			return fmt.Errorf("repair destination changed: %s", p.path)
		}
	}
	if p.value.Kind == "absent" {
		if err := os.Remove(p.path); err != nil && !os.IsNotExist(err) {
			return err
		}
	} else {
		rename := s.Rename
		if rename == nil {
			rename = os.Rename
		}
		if err := rename(p.temporary, p.path); err != nil {
			return err
		}
	}
	return s.syncDirectory(filepath.Dir(p.path))
}

func (s Store) publish(path string, value image) error {
	staged, err := stagePublication(path, value, s.syncDirectory)
	if err != nil {
		return err
	}
	defer staged.close()
	return staged.publish(s, nil)
}
