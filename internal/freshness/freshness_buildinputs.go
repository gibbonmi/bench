// Build-input discovery and symlink-hardened file reading for package freshness.
package freshness

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

const auxiliaryInputsManifest = "scripts/go-build.inputs"

type listedPackage struct {
	Dir          string
	GoFiles      []string
	CgoFiles     []string
	CFiles       []string
	CXXFiles     []string
	MFiles       []string
	HFiles       []string
	FFiles       []string
	SFiles       []string
	SwigFiles    []string
	SwigCXXFiles []string
	SysoFiles    []string
	EmbedFiles   []string
}

func buildInputs(root string) ([]string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	command := exec.Command("go", "list", "-buildvcs=false", "-json", "-deps", "./cmd/bench")
	command.Dir = root
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("resolve Bench build inputs: %w", err)
	}
	decoder := json.NewDecoder(strings.NewReader(string(output)))
	paths := map[string]struct{}{}
	for {
		var pkg listedPackage
		if err := decoder.Decode(&pkg); err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decode resolved Bench build inputs: %w", err)
		}
		for _, name := range packageFiles(pkg) {
			path := filepath.Join(pkg.Dir, name)
			if isWithinRoot(root, path) {
				paths[path] = struct{}{}
			}
		}
	}
	for _, name := range []string{
		"go.mod",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if _, err := os.Lstat(path); err != nil {
			return nil, fmt.Errorf("required build input %q: %w", name, err)
		}
		paths[path] = struct{}{}
	}
	auxiliary, err := auxiliaryBuildInputs(root)
	if err != nil {
		return nil, err
	}
	for _, path := range auxiliary {
		paths[path] = struct{}{}
	}
	if path := filepath.Join(root, "go.sum"); exists(path) {
		paths[path] = struct{}{}
	}
	ordered := make([]string, 0, len(paths))
	for path := range paths {
		ordered = append(ordered, path)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return filepath.ToSlash(ordered[i]) < filepath.ToSlash(ordered[j])
	})
	return ordered, nil
}

func auxiliaryBuildInputs(root string) ([]string, error) {
	manifest := filepath.Join(root, filepath.FromSlash(auxiliaryInputsManifest))
	data, err := regularContents(manifest)
	if err != nil {
		return nil, fmt.Errorf("read auxiliary build-input manifest %q: %w", auxiliaryInputsManifest, err)
	}
	paths := []string{manifest}
	keys := map[string]struct{}{}
	for number, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		key, name, ok := strings.Cut(line, "=")
		if !ok || key == "" || name == "" || strings.TrimSpace(key) != key || strings.TrimSpace(name) != name {
			return nil, fmt.Errorf("malformed auxiliary build input at %s:%d", auxiliaryInputsManifest, number+1)
		}
		if _, exists := keys[key]; exists {
			return nil, fmt.Errorf("duplicate auxiliary build input key %q", key)
		}
		keys[key] = struct{}{}
		path := filepath.Join(root, filepath.FromSlash(name))
		if !isWithinRoot(root, path) || filepath.IsAbs(name) {
			return nil, fmt.Errorf("auxiliary build input %q leaves the source root", name)
		}
		if _, err := os.Lstat(path); err != nil {
			return nil, fmt.Errorf("required build input %q: %w", name, err)
		}
		paths = append(paths, path)
	}
	if len(paths) == 1 {
		return nil, fmt.Errorf("auxiliary build-input manifest %q is empty", auxiliaryInputsManifest)
	}
	return paths, nil
}

func packageFiles(pkg listedPackage) []string {
	var files []string
	for _, group := range [][]string{
		pkg.GoFiles, pkg.CgoFiles, pkg.CFiles, pkg.CXXFiles, pkg.MFiles, pkg.HFiles,
		pkg.FFiles, pkg.SFiles, pkg.SwigFiles, pkg.SwigCXXFiles, pkg.SysoFiles, pkg.EmbedFiles,
	} {
		files = append(files, group...)
	}
	return files
}

func isWithinRoot(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func regularContents(path string) ([]byte, error) {
	return secureContents(path, false)
}

func secureContents(path string, executable bool) ([]byte, error) {
	if err := rejectSymlinkComponents(path); err != nil {
		return nil, err
	}
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("is a symbolic link")
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("is not a regular file")
	}
	if executable && info.Mode()&0o111 == 0 {
		return nil, fmt.Errorf("is not executable")
	}
	if executable && info.Size() == 0 {
		return nil, fmt.Errorf("is empty")
	}
	// A path replacement after Lstat stays untrusted instead of redirecting this read.
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_CLOEXEC|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(fd), path)
	defer file.Close()
	opened, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !os.SameFile(info, opened) {
		return nil, fmt.Errorf("changed while opening")
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if executable && len(data) == 0 {
		return nil, fmt.Errorf("is empty")
	}
	return data, nil
}

func rejectSymlinkComponents(path string) error {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	volume := filepath.VolumeName(absolute)
	current := volume + string(filepath.Separator)
	relative := strings.TrimPrefix(absolute, current)
	for _, component := range strings.Split(relative, string(filepath.Separator)) {
		if component == "" {
			continue
		}
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("path component %q is a symbolic link", current)
		}
	}
	return nil
}
