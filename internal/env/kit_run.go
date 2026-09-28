package env

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// KitTestRun owns the private home and temporary directories for one kit test run.
// A plain go test opens no run and receives none of this policy.
type KitTestRun struct {
	dir     string
	entries []string
}

// OpenKitTestRun preserves Go's resolved settings before it replaces HOME.
// The caller derives its build cache from the base environment before it merges Entries.
func OpenKitTestRun(base []string) (*KitTestRun, error) {
	if home := kitEnvValue(base, "HOME"); !filepath.IsAbs(home) {
		return nil, errors.New("kit test run: HOME must be present and absolute")
	}
	settingNames := []string{"GOMODCACHE", "GOPATH", "GOENV"}
	cmd := exec.Command("go", append([]string{"env", "-json"}, settingNames...)...)
	cmd.Env = base
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("kit test run: resolve Go environment: %w", err)
	}
	var settings map[string]string
	if err := json.Unmarshal(output, &settings); err != nil {
		return nil, fmt.Errorf("kit test run: read Go environment: %w", err)
	}
	tmp := kitEnvValue(base, "TMPDIR")
	if tmp == "" {
		tmp = os.TempDir()
	}
	dir, err := kitRunDirectory(tmp)
	if err != nil {
		return nil, fmt.Errorf("kit test run: TMPDIR: %w", err)
	}
	run := &KitTestRun{dir: dir}
	for _, name := range []string{"h", "t"} {
		if err := os.Mkdir(filepath.Join(dir, name), 0o700); err != nil {
			return nil, errors.Join(fmt.Errorf("kit test run: TMPDIR: %w", err), run.Close())
		}
	}
	run.entries = append([]string{
		"HOME=" + filepath.Join(dir, "h"),
		"TMPDIR=" + filepath.Join(dir, "t"),
		"GIT_CONFIG_GLOBAL=" + os.DevNull,
		"GIT_CONFIG_NOSYSTEM=1",
	}, GitTestConfig()...)
	for _, name := range settingNames {
		run.entries = append(run.entries, name+"="+settings[name])
	}
	return run, nil
}

func kitRunDirectory(tmp string) (string, error) {
	// Short names leave room for Unix socket paths in descendant tests.
	for {
		dir := filepath.Join(tmp, rand.Text()[:6])
		if err := os.Mkdir(dir, 0o700); !errors.Is(err, os.ErrExist) {
			return dir, err
		}
	}
}

// Entries returns the environment overrides for children of this run.
func (r *KitTestRun) Entries() []string {
	return append([]string(nil), r.entries...)
}

// Close restores directory access before removal, because tests can leave unreadable directories.
// It does not follow symbolic links outside the run.
func (r *KitTestRun) Close() error {
	err := filepath.WalkDir(r.dir, func(path string, entry fs.DirEntry, err error) error {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.Chmod(path, 0o700)
		}
		return nil
	})
	if err = errors.Join(err, os.RemoveAll(r.dir)); err != nil {
		return fmt.Errorf("kit test run: remove TMPDIR run: %w", err)
	}
	return nil
}

func kitEnvValue(base []string, name string) string {
	for i := len(base) - 1; i >= 0; i-- {
		if value, ok := strings.CutPrefix(base[i], name+"="); ok {
			return value
		}
	}
	return ""
}
