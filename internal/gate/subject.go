package gate

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

var envNameRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type identityCollector struct {
	w                         io.Writer
	entries                   int
	entryLimit                int
	bytes                     int64
	runtimeRoot, identityRoot string
}

func buildSubject(root string) (subject, error) { return buildSubjectFor(root, root) }

func buildSubjectFor(root, identityRoot string) (subject, error) {
	return buildSubjectForPolicy(root, identityRoot, policyVersion)
}

func buildSubjectForGeneration(root, identityRoot string, generation *treeGeneration) (subject, error) {
	return buildSubjectForTree(root, identityRoot, policyVersion, generation.tree)
}

func buildSubjectForPolicy(root, identityRoot, policy string) (subject, error) {
	return buildSubjectOverTree(root, identityRoot, policy, benchgit.TreeHash)
}

func buildSubjectOverTree(root, identityRoot, policy string, treeHash func(string) string) (subject, error) {
	return buildSubjectForTree(root, identityRoot, policy, treeHash(root))
}

func buildSubjectForTree(root, identityRoot, policy, tree string) (subject, error) {
	root, err := canonicalSubjectRoot(root)
	if err != nil {
		return subject{}, err
	}
	identityRoot, err = canonicalSubjectRoot(identityRoot)
	if err != nil {
		return subject{}, err
	}
	if !treeHashRE.MatchString(tree) {
		return subject{}, errors.New("tree unavailable")
	}
	pathEnv := os.Getenv("PATH")
	res := Resolve(root, os.Getenv("BENCH_GATE"), RealFS())
	s := subject{Tree: tree, Resolution: res, Closed: true, Env: []string{"PATH=" + pathEnv}}
	m, manifestIdentity, reason := loadManifest(root)
	if reason != "" {
		s.Closed, s.Reason = false, reason
	}
	h := sha256.New()
	for _, value := range []string{policy, identityRoot, tree, resolutionName(res.Kind), res.Command, pathEnv, manifestIdentity} {
		frame(h, value)
	}
	c := &identityCollector{w: h, entryLimit: manifestEntryLimit(), runtimeRoot: root, identityRoot: identityRoot}
	if res.Kind != None {
		if err := c.hashResolution(root, res, pathEnv); err != nil {
			s.open("launcher closure unavailable")
		}
	}
	if err := hashProspectivePreparation(c, h, root, pathEnv); err != nil {
		s.open("launcher closure unavailable")
	}
	if m != nil {
		for _, name := range m.Environment {
			value, ok := os.LookupEnv(name)
			if !ok {
				s.open("declared environment unavailable")
				continue
			}
			s.Env = append(s.Env, name+"="+value)
			frame(h, "environment")
			frame(h, name)
			frame(h, value)
		}
		for _, path := range m.Paths {
			frame(h, "path")
			frame(h, path)
			if err := c.hashRepoPath(root, path); err != nil {
				s.open("declared path unavailable")
			}
		}
		for _, tool := range m.Tools {
			frame(h, "tool")
			frame(h, tool)
			if err := c.hashExecutable(root, tool, pathEnv, false, 0); err != nil {
				s.open("declared tool unavailable")
			}
		}
	}
	s.Oracle = hex.EncodeToString(h.Sum(nil))
	return s, nil
}

func (s *subject) open(reason string) {
	if s.Closed {
		s.Closed = false
		s.Reason = reason
	}
}
func loadManifest(root string) (*manifest, string, string) {
	path := filepath.Join(root, ".bench", "gate-inputs.json")
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, "", "gate input manifest absent"
	}
	if err != nil {
		return nil, "", "gate input manifest unavailable"
	}
	data, readErr := io.ReadAll(io.LimitReader(f, manifestLimit+1))
	closeErr := f.Close()
	if readErr != nil || closeErr != nil || len(data) == 0 || len(data) > manifestLimit || !utf8.Valid(data) {
		return nil, "", "gate input manifest invalid"
	}
	var m manifest
	if strictJSON(data, &m) != nil || requireObjectFields(data, []string{"closure", "environment", "paths", "schema", "tools"}) != nil || m.Schema != 1 || (m.Closure != "local" && m.Closure != "remote") || m.Environment == nil || m.Paths == nil || m.Tools == nil {
		return nil, "", "gate input manifest invalid"
	}
	for _, name := range m.Environment {
		if !envNameRE.MatchString(name) || hasUnsafeText(name) {
			return nil, "", "gate input manifest invalid"
		}
	}
	for _, path := range m.Paths {
		if path == "" || path == "." || filepath.IsAbs(path) || filepath.Clean(filepath.FromSlash(path)) != filepath.FromSlash(path) || filepath.ToSlash(filepath.FromSlash(path)) != path || path == ".." || strings.HasPrefix(path, "../") || hasUnsafeText(path) {
			return nil, "", "gate input manifest invalid"
		}
	}
	for _, tool := range m.Tools {
		if tool == "" || hasUnsafeText(tool) {
			return nil, "", "gate input manifest invalid"
		}
		if !filepath.IsAbs(tool) && strings.Contains(tool, "/") && (filepath.Clean(filepath.FromSlash(tool)) != filepath.FromSlash(tool) || tool == ".." || strings.HasPrefix(tool, "../")) {
			return nil, "", "gate input manifest invalid"
		}
	}
	sort.Strings(m.Environment)
	m.Environment = slices.Compact(m.Environment)
	sort.Strings(m.Paths)
	m.Paths = slices.Compact(m.Paths)
	sort.Strings(m.Tools)
	m.Tools = slices.Compact(m.Tools)
	canonical, _ := json.Marshal(m)
	reason := ""
	if m.Closure == "remote" {
		reason = "remote oracle"
	}
	sum := sha256.Sum256(canonical)
	return &m, hex.EncodeToString(sum[:]), reason
}
func hasUnsafeText(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return !utf8.ValidString(s)
}
func frame(w io.Writer, value string) {
	_, _ = fmt.Fprintf(w, "%d:", len(value))
	_, _ = io.WriteString(w, value)
}
func (c *identityCollector) addEntry() error {
	c.entries++
	if c.entries > c.entryLimit {
		return errors.New("entry limit")
	}
	return nil
}
func (c *identityCollector) copyFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	n, copyErr := io.Copy(c.w, io.LimitReader(f, (1<<30)-c.bytes+1))
	closeErr := f.Close()
	c.bytes += n
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if c.bytes > 1<<30 {
		return errors.New("byte limit")
	}
	return nil
}
func (c *identityCollector) hashResolution(root string, resolution Resolution, pathEnv string) error {
	tools := []string{}
	switch resolution.Kind {
	case GateSh:
		return c.hashExecutable(root, filepath.Join(root, ".bench", "gate.sh"), pathEnv, true, 0)
	case BenchGate:
		tools = []string{"bash"}
	case Pnpm:
		tools = []string{"bash", "pnpm"}
	case Npm:
		tools = []string{"bash", "npm"}
	case Pyproject:
		tools = []string{"bash", "mypy", "pytest", "ruff"}
	case Cargo:
		tools = []string{"bash", "cargo", "rustc", "clippy-driver"}
	}
	for _, tool := range tools {
		if err := c.hashExecutable(root, tool, pathEnv, false, 0); err != nil {
			return err
		}
	}
	return nil
}
func (c *identityCollector) hashRepoPath(root, rel string) error {
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := confinedPath(root, path); err != nil {
		return err
	}
	return c.hashTree(path, root, 0)
}
func (c *identityCollector) hashTree(path, root string, depth int) error {
	if depth > 64 {
		return errors.New("declared path symlink depth limit")
	}
	return filepath.WalkDir(path, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := c.addEntry(); err != nil {
			return err
		}
		info, err := os.Lstat(current)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(path, current)
		frame(c.w, filepath.ToSlash(rel))
		frame(c.w, info.Mode().String())
		switch {
		case info.Mode().IsRegular():
			return c.copyFile(current)
		case info.Mode()&os.ModeSymlink != 0:
			if err := confinedPath(root, current); err != nil {
				return err
			}
			if err := c.hashLinkChain(current, 0); err != nil {
				return err
			}
			resolved, err := filepath.EvalSymlinks(current)
			if err != nil {
				return err
			}
			return c.hashTree(resolved, root, depth+1)
		case info.IsDir():
			return nil
		default:
			return errors.New("unsupported file type")
		}
	})
}
func confinedPath(root, path string) error {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(root, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return errors.New("escaped path")
	}
	return nil
}
func (c *identityCollector) hashExecutable(root, name, pathEnv string, confined bool, depth int) error {
	if depth > 64 {
		return errors.New("launcher hop limit")
	}
	path, err := resolveTool(root, name, pathEnv)
	if err != nil {
		return err
	}
	if confined {
		if err := confinedPath(root, path); err != nil {
			return err
		}
	}
	if err := c.hashLinkChain(path, depth); err != nil {
		return err
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&0o111 == 0 {
		return errors.New("tool unavailable")
	}
	if err := c.addEntry(); err != nil {
		return err
	}
	frame(c.w, c.identityPath(resolved))
	frame(c.w, info.Mode().String())
	if err := c.copyFile(resolved); err != nil {
		return err
	}
	line, err := firstLine(resolved)
	if err != nil || !strings.HasPrefix(line, "#!") {
		return err
	}
	words := strings.Fields(strings.TrimSpace(strings.TrimPrefix(line, "#!")))
	if len(words) == 0 {
		return errors.New("invalid shebang")
	}
	if filepath.Base(words[0]) == "env" {
		if err := c.hashExecutable(root, words[0], pathEnv, false, depth+1); err != nil {
			return err
		}
		selected := ""
		for _, word := range words[1:] {
			if !strings.HasPrefix(word, "-") && !strings.Contains(word, "=") {
				selected = word
				break
			}
		}
		if selected == "" {
			return errors.New("unresolved env shebang")
		}
		return c.hashExecutable(root, selected, pathEnv, false, depth+1)
	}
	return c.hashExecutable(root, words[0], pathEnv, false, depth+1)
}
func (c *identityCollector) hashLinkChain(path string, depth int) error {
	seen := map[string]bool{}
	for hop := depth; hop <= 64; hop++ {
		absolute, err := filepath.Abs(path)
		if err != nil || seen[absolute] {
			return errors.New("cyclic symlink")
		}
		seen[absolute] = true
		info, err := os.Lstat(absolute)
		if err != nil {
			return err
		}
		if err := c.addEntry(); err != nil {
			return err
		}
		frame(c.w, c.identityPath(absolute))
		frame(c.w, info.Mode().String())
		if info.Mode()&os.ModeSymlink == 0 {
			return nil
		}
		target, err := os.Readlink(absolute)
		if err != nil {
			return err
		}
		frame(c.w, target)
		if filepath.IsAbs(target) {
			path = target
		} else {
			path = filepath.Join(filepath.Dir(absolute), target)
		}
	}
	return errors.New("symlink hop limit")
}
func (c *identityCollector) identityPath(path string) string {
	rel, err := filepath.Rel(c.runtimeRoot, path)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return filepath.Join(c.identityRoot, rel)
	}
	return path
}

func canonicalSubjectRoot(root string) (string, error) {
	resolved, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	return filepath.Abs(resolved)
}

func resolveTool(root, name, pathEnv string) (string, error) {
	if filepath.IsAbs(name) {
		return name, nil
	}
	if strings.ContainsRune(name, os.PathSeparator) {
		return filepath.Join(root, filepath.FromSlash(name)), nil
	}
	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == "" {
			dir = "."
		}
		candidate := filepath.Join(dir, name)
		info, err := os.Stat(candidate)
		if err == nil && info.Mode().IsRegular() && info.Mode()&0o111 != 0 {
			return candidate, nil
		}
	}
	return "", errors.New("tool unavailable")
}

func firstLine(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	line, err := bufio.NewReader(io.LimitReader(f, 4096)).ReadString('\n')
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return line, err
}
