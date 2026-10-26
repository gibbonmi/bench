package gate

import (
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
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gocache"
)

var envNameRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// identityCollector feeds a closure to both digests through w, and each path to each digest.
type identityCollector struct {
	w, oracle, identity       io.Writer
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
	// h is the oracle digest and id the grade identity. The verdict inputs go to both
	// through w, so the oracle frames keep their order. The policy, the identity root, and
	// the tree go to the oracle only.
	h, id := sha256.New(), sha256.New()
	w := io.MultiWriter(h, id)
	frame(h, policy)
	frame(h, identityRoot)
	frame(h, tree)
	for _, value := range []string{resolutionName(res.Kind), res.Command, pathEnv, manifestIdentity} {
		frame(w, value)
	}
	c := &identityCollector{w: w, oracle: h, identity: id, entryLimit: manifestEntryLimit(), runtimeRoot: root, identityRoot: identityRoot}
	if res.Kind != None {
		if err := c.hashResolution(root, res, pathEnv); err != nil {
			s.open("launcher closure unavailable")
		}
	}
	// res.Kind still names the ordinary script, so a prospective script sets the source kind.
	sourceKind := res.Kind
	prospective, err := hashProspectivePreparation(c, root, identityRoot, pathEnv)
	if err != nil {
		s.open("launcher closure unavailable")
	} else if prospective {
		sourceKind = ProspectiveGateSh
	}
	if m != nil {
		supplied, missing := suppliedEnvironment(m.Environment)
		if len(missing) > 0 {
			s.open("declared environment unavailable: " + strings.Join(missing, ", "))
		}
		for _, entry := range supplied {
			s.Env = append(s.Env, entry[0]+"="+entry[1])
			frame(w, "environment")
			frame(w, entry[0])
			frame(w, entry[1])
		}
		for _, path := range m.Paths {
			frame(w, "path")
			frame(w, path)
			if err := c.hashRepoPath(root, path); err != nil {
				s.open("declared path unavailable")
			}
		}
		for _, tool := range m.Tools {
			frame(w, "tool")
			frame(w, tool)
			if err := c.hashExecutable(root, tool, pathEnv, false, 0); err != nil {
				s.open("declared tool unavailable")
			}
		}
	}
	// The build cache entry rides the closure unhashed. It is derived from the
	// declared HOME the loop above already framed, so hashing it would frame the same
	// fact twice. A closure that declares no HOME carries no entry, because nothing in
	// it names the directory. A closure that declares one the derivation refuses fails
	// here instead: dropping the entry there would hand the oracle an environment its
	// own children never ran with, and no line would say so.
	switch applied, err := gocache.Apply(s.Env); {
	case err == nil:
		s.Env = applied
	case gocache.Declared(s.Env):
		return subject{}, err
	}
	s.Oracle = hex.EncodeToString(h.Sum(nil))
	frameBaselineSchedule(id, identityRoot)
	frame(id, "run binary")
	frame(id, buildInputsIdentity(runBinarySource(sourceKind, root, identityRoot)))
	frame(id, "running executable")
	frame(id, runningExecutableDigest())
	frame(id, requireCapabilitiesEnv)
	frame(id, strconv.FormatBool(strict()))
	s.Identity = hex.EncodeToString(id.Sum(nil))
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
	c.framePath(resolved)
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
		c.framePath(absolute)
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

// framePath frames an executable path in each digest's spelling. The oracle spells a path in
// the runtime root under the identity root. The grade identity spells a path in the runtime
// root, then one in the identity root, relative to that root, so it frames no checkout path.
func (c *identityCollector) framePath(path string) {
	oracle, grade := path, path
	if rel, ok := relativeInside(c.runtimeRoot, path); ok {
		oracle, grade = filepath.Join(c.identityRoot, rel), rel
	} else if rel, ok := relativeInside(c.identityRoot, path); ok {
		grade = rel
	}
	frame(c.oracle, oracle)
	frame(c.identity, grade)
}
func relativeInside(root, path string) (string, bool) {
	rel, err := filepath.Rel(root, path)
	return rel, err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// frameBaselineSchedule frames the phase-table loader's answer at the identity root: present
// bytes as bytes, an absence as one value, and a defect as one value with no path.
func frameBaselineSchedule(w io.Writer, identityRoot string) {
	data, present, err := readManifest(manifestPath(identityRoot))
	switch {
	case err != nil:
		frame(w, "schedule defect")
	case !present:
		frame(w, "schedule absent")
	default:
		frame(w, "schedule")
		frame(w, string(data))
	}
}

// runningExecutableDigest answers the running executable's digest once for each process, or a
// stable value when it does not read. It is a seam, so a test can vary the executable.
var runningExecutableDigest = sync.OnceValue(func() string {
	value, err := os.Executable()
	if err == nil {
		value, err = filepath.EvalSymlinks(value)
	}
	if err == nil {
		value, err = brokermanifest.Digest(value)
	}
	if err != nil {
		return "unreadable"
	}
	return value
})

// canonicalSubjectRoot refuses a root that does not exist. The owner keeps the absolute
// spelling of an absent path rather than refusing, and a subject keyed to an absent root
// would hash as a closed identity over nothing.
func canonicalSubjectRoot(root string) (string, error) {
	resolved, err := canonicalpath.Resolve(root)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(resolved); err != nil {
		return "", err
	}
	return resolved, nil
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
	read := bounds.Read(f, bounds.SubjectFirstLineLimit)
	if read.Status == bounds.ReadFailed {
		return "", read.Err
	}
	return strings.SplitAfterN(string(read.Data), "\n", 2)[0], nil
}
