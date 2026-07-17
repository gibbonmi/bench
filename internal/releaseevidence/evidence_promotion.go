package releaseevidence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var stageOwnerBytes = []byte("bench-preflight-stage-v1\n")

var exchangeEvidenceDirs = atomicExchangeDirs

func SetExchangeForTesting(exchange func(string, string) error) func() {
	previous := exchangeEvidenceDirs
	exchangeEvidenceDirs = exchange
	return func() { exchangeEvidenceDirs = previous }
}

func AtomicExchangeForTesting(left, right string) error { return atomicExchangeDirs(left, right) }

func PromoteEvidence(root string, mode Mode, results []Result, manifest Manifest) error {
	return PromoteEvidenceFiles(root, mode, results, manifest, nil)
}

func PromoteEvidenceFiles(root string, mode Mode, results []Result, manifest Manifest, files map[string][]byte) error {
	dist := filepath.Join(root, "dist")
	if info, err := os.Lstat(dist); err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
		return fmt.Errorf("dist output target is not a real directory")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dist, 0o755); err != nil {
		return err
	}
	if err := cleanupAbandonedStages(dist); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(dist, ".preflight-stage-")
	if err != nil {
		return err
	}
	owner := stage + ".owner"
	if err := writeBytesSync(owner, stageOwnerBytes); err != nil {
		_ = os.RemoveAll(stage)
		return err
	}
	keep := false
	defer func() {
		if !keep {
			_ = os.RemoveAll(stage)
			_ = os.Remove(owner)
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
	for name, data := range files {
		if filepath.Base(name) != name || name == "" {
			return fmt.Errorf("invalid promoted evidence file name: %s", name)
		}
		if err := writeBytesSync(filepath.Join(stage, name), data); err != nil {
			return err
		}
	}
	if err := syncDir(stage); err != nil {
		return err
	}
	target := filepath.Join(dist, "preflight")
	if info, err := os.Lstat(target); err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
		return fmt.Errorf("preflight output target is not a real directory")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if _, err := os.Lstat(target); os.IsNotExist(err) {
		if err := os.Rename(stage, target); err != nil {
			return err
		}
		keep = true
	} else if err != nil {
		return err
	} else {
		if err := exchangeEvidenceDirs(stage, target); err != nil {
			return err
		}
		keep = true
	}
	if err := syncDir(dist); err != nil {
		return err
	}
	if err := os.RemoveAll(stage); err != nil {
		return err
	}
	if err := os.Remove(owner); err != nil {
		return err
	}
	return nil
}

func cleanupAbandonedStages(dist string) error {
	entries, err := os.ReadDir(dist)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, ".preflight-stage-") || !strings.HasSuffix(name, ".owner") {
			continue
		}
		owner := filepath.Join(dist, name)
		info, err := os.Lstat(owner)
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("unsafe preflight stage ownership lease: %s", name)
		}
		data, err := os.ReadFile(owner)
		if err != nil || !bytes.Equal(data, stageOwnerBytes) {
			return fmt.Errorf("invalid preflight stage ownership lease: %s", name)
		}
		stage := strings.TrimSuffix(owner, ".owner")
		stageInfo, err := os.Lstat(stage)
		if err != nil || !stageInfo.IsDir() || stageInfo.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("owned preflight stage is unsafe: %s", filepath.Base(stage))
		}
		if err := os.RemoveAll(stage); err != nil {
			return err
		}
		if err := os.Remove(owner); err != nil {
			return err
		}
	}
	return nil
}

func writeJSONSync(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return writeBytesSync(path, append(data, '\n'))
}

func writeBytesSync(path string, data []byte) error {
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
