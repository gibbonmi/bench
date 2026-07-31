// Package freshness verifies that a Bench executable was built from the current sources.
package freshness

import (
	"crypto/sha256"
	"encoding/hex"
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

const sealSchema = 1
const auxiliaryInputsManifest = "scripts/go-build.inputs"

type seal struct {
	Schema     int    `json:"schema"`
	Sources    string `json:"sources"`
	Executable string `json:"executable"`
}

// Digest returns the deterministic content digest of Bench's local build inputs.
func Digest(root string) (string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	paths, err := buildInputs(root)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	for _, path := range paths {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return "", err
		}
		contents, err := regularContents(path)
		if err != nil {
			return "", fmt.Errorf("read build input %q: %w", rel, err)
		}
		name := filepath.ToSlash(rel)
		fmt.Fprintf(hash, "%d:%s%d:", len(name), name, len(contents))
		if _, err := hash.Write(contents); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Publish replaces executable with staged and writes its matching content seal.
func Publish(root, staged, executable string) error {
	sources, err := Digest(root)
	if err != nil {
		return err
	}
	binary, err := secureContents(staged, true)
	if err != nil {
		return fmt.Errorf("stage executable %q: %w", staged, err)
	}
	encoded, err := json.Marshal(seal{
		Schema:     sealSchema,
		Sources:    sources,
		Executable: digestBytes(binary),
	})
	if err != nil {
		return err
	}
	if err := os.Remove(sealPath(executable)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove old seal: %w", err)
	}
	if err := os.Rename(staged, executable); err != nil {
		return fmt.Errorf("publish executable: %w", err)
	}
	if err := writeSeal(sealPath(executable), encoded); err != nil {
		return fmt.Errorf("publish seal: %w", err)
	}
	return nil
}

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

// Verify reports whether executable has a matching content seal for root.
func Verify(root, executable string) error {
	binary, err := secureContents(executable, true)
	if err != nil {
		return refusal(root, executable, err)
	}
	sealData, err := secureContents(sealPath(executable), false)
	if err != nil {
		return refusal(root, executable, fmt.Errorf("seal %q: %w", sealPath(executable), err))
	}
	stored, err := parseSeal(sealData)
	if err != nil {
		return refusal(root, executable, fmt.Errorf("seal %q: %w", sealPath(executable), err))
	}
	sources, err := Digest(root)
	if err != nil {
		return refusal(root, executable, err)
	}
	if stored.Sources != sources {
		return refusal(root, executable, fmt.Errorf("seal source digest does not match current build inputs"))
	}
	if stored.Executable != digestBytes(binary) {
		return refusal(root, executable, fmt.Errorf("seal executable digest does not match binary contents"))
	}
	return nil
}

// Check verifies an executable from current sources, then requires its freshness subcommand.
func Check(root, executable string) error {
	if err := Verify(root, executable); err != nil {
		return err
	}
	command := exec.Command(executable, "freshness-check", root)
	command.Stdout = io.Discard
	command.Stderr = io.Discard
	if err := command.Run(); err != nil {
		return refusal(root, executable, fmt.Errorf("freshness-check failed"))
	}
	return nil
}

func parseSeal(data []byte) (seal, error) {
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	var stored seal
	if err := decoder.Decode(&stored); err != nil {
		return seal{}, fmt.Errorf("malformed: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return seal{}, fmt.Errorf("malformed trailing data")
	}
	if stored.Schema != sealSchema || !isDigest(stored.Sources) || !isDigest(stored.Executable) {
		return seal{}, fmt.Errorf("malformed contents")
	}
	return stored, nil
}

func writeSeal(path string, data []byte) (err error) {
	temporary, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = os.Remove(temporary.Name())
		}
	}()
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporary.Name(), path)
}

func sealPath(executable string) string { return executable + ".seal" }

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func isDigest(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func refusal(root, executable string, cause error) error {
	return fmt.Errorf("bench binary %q is untrusted: %v; rebuild with %s", executable, cause, RebuildAction(root))
}

// RebuildAction returns the copy-paste command that republishes root's Bench binary.
func RebuildAction(root string) string {
	return fmt.Sprintf("cd %s && bash scripts/go-build.sh %s %s", shellQuote(root), shellQuote(root), shellQuote(filepath.Join(root, "dist", "bench")))
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
