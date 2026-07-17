package preflight

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func PromoteEvidence(root string, mode Mode, results []Result, manifest Manifest) error {
	dist := filepath.Join(root, "dist")
	if info, err := os.Lstat(dist); err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
		return fmt.Errorf("dist output target is not a real directory")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dist, 0o755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(dist, ".preflight-stage-")
	if err != nil {
		return err
	}
	keep := false
	defer func() {
		if !keep {
			_ = os.RemoveAll(stage)
		}
	}()
	for _, result := range results {
		record := Record{SchemaVersion: 1, Phase: result.Name, Mode: mode, Status: result.Status, ExitCode: result.ExitCode, Error: result.Failure}
		if err := writeJSONSync(filepath.Join(stage, result.Name+".json"), record); err != nil {
			return err
		}
	}
	if err := writeJSONSync(filepath.Join(stage, "manifest.json"), manifest); err != nil {
		return err
	}
	if err := syncDir(stage); err != nil {
		return err
	}
	target := filepath.Join(dist, "preflight")
	if info, err := os.Lstat(target); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("preflight output target is a symlink")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	backup := filepath.Join(dist, ".preflight-previous")
	_ = os.RemoveAll(backup)
	if _, err := os.Lstat(target); err == nil {
		if err := os.Rename(target, backup); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(stage, target); err != nil {
		if _, statErr := os.Stat(backup); statErr == nil {
			_ = os.Rename(backup, target)
		}
		return err
	}
	keep = true
	if err := syncDir(dist); err != nil {
		return err
	}
	if err := os.RemoveAll(backup); err != nil {
		return err
	}
	return nil
}

func writeJSONSync(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err = f.Write(data); err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func syncDir(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync %s: %w", path, err)
	}
	return nil
}
