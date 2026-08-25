package gate

// This file loads the project-owned phase table: <graded root>/.bench/phases.json. A
// repo that ships one declares its whole gate as data; a repo that ships none keeps the
// built-in kit table. Those are the only two states the loader accepts.
//
// Everything between them is a manifest whose author meant something the loader
// cannot know. Examples: a truncated write, a broken link, a typo'd key, an edge to
// nowhere. The loader reds before any phase runs, rather than grade the tree with the
// wrong oracle.

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/canary"
)

func manifestPath(root string) string {
	return filepath.Join(root, filepath.FromSlash(canary.PhaseManifestPath))
}

// phasePolicyRoot is the tree whose phase manifest selects the schedule. It is the
// graded root for an ordinary run. A prospective landing names the baseline instead, so
// the candidate tree under grade cannot omit the checks that grade it. A named baseline
// that is not a usable directory refuses; falling back would reinstate the omission the
// naming exists to prevent.
func phasePolicyRoot(root string) (string, error) {
	baseline := os.Getenv(baselinePolicyEnv)
	if baseline == "" {
		return root, nil
	}
	info, err := os.Stat(baseline)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("gate: baseline phase schedule root %s is unusable", baseline)
	}
	return baseline, nil
}

type manifestDoc struct {
	Phases []manifestPhase `json:"phases"`
	// Lane is the project's declared fast lane, which a worktree commit runs in place
	// of the whole-project gate. It carries the same entry schema as Phases, because a
	// lane check and a gate phase are the same kind of thing run by the same runner.
	// An absent array declares no lane, which keeps the full-gate commit.
	Lane []manifestPhase `json:"lane"`
}

type manifestPhase struct {
	Name     string            `json:"name"`
	Argv     []string          `json:"argv"`
	Env      map[string]string `json:"env"`
	Needs    []string          `json:"needs"`
	Optional bool              `json:"optional"`
	Dir      string            `json:"dir"`
}

// phaseTable resolves the table gate-phases runs: the graded root's manifest when it
// declares one, the built-in kit table when it declares none. root is the tree under
// grade, resolved through git by every caller and so absolute — a manifest phase's
// working directory is anchored to it.
func phaseTable(root, kit string) ([]Phase, error) {
	policy, err := phasePolicyRoot(root)
	if err != nil {
		return nil, err
	}
	path := manifestPath(policy)
	data, present, err := readManifest(path)
	if err != nil {
		return nil, err
	}
	if !present {
		if _, err := goModuleToolchain(root); err != nil {
			return nil, err
		}
		return benchkitPhasesForCommand(root, kit), nil
	}
	return parseManifest(path, root, data)
}

// readManifest reports the manifest bytes, or that no manifest exists at all. Statting
// before opening is load-bearing: a FIFO left at the path blocks an open forever, so
// the mode must be known before anything opens it.
//
// os.Stat follows the link, so a symlink whose target is gone answers ErrNotExist
// exactly like an honest absence. os.Lstat tells the two apart, and only the honest
// absence may fall back.
func readManifest(path string) ([]byte, bool, error) {
	if _, err := os.Lstat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, defect(path, "unreadable manifest", err.Error())
	}
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			target, _ := os.Readlink(path)
			return nil, false, defect(path, "dangling symlink", target)
		}
		return nil, false, defect(path, "unreadable manifest", err.Error())
	}
	if !info.Mode().IsRegular() {
		return nil, false, defect(path, "not a regular file", info.Mode().String())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false, defect(path, "unreadable manifest", err.Error())
	}
	return data, true, nil
}

// decodeManifest turns manifest bytes into the declaration document, or into the one
// diagnostic shape a defect reports through. The phase table and the lane both read the
// same document, so a defect in the file reads the same whichever array asked for it.
func decodeManifest(path string, data []byte) (manifestDoc, error) {
	var doc manifestDoc
	if err := strictJSON(data, &doc); err != nil {
		if key, ok := unknownField(err); ok {
			return manifestDoc{}, defect(path, "unknown field", key)
		}
		if errors.Is(err, errTrailingJSON) {
			// A second document is the classic half-written or double-appended
			// manifest, and encoding/json alone would grade it as if only the first
			// existed.
			return manifestDoc{}, defect(path, "parse error", "unexpected trailing content")
		}
		return manifestDoc{}, defect(path, "parse error", err.Error())
	}
	return doc, nil
}

func parseManifest(path, root string, data []byte) ([]Phase, error) {
	if strings.TrimSpace(string(data)) == "" {
		return nil, defect(path, "empty manifest", "no phases declared")
	}
	doc, err := decodeManifest(path, data)
	if err != nil {
		return nil, err
	}
	if len(doc.Phases) == 0 {
		return nil, defect(path, "empty manifest", "no phases declared")
	}
	return validateManifest(path, root, doc.Phases)
}

// unknownField lifts the rejected key out of encoding/json's strict-decoding error.
// This class is worth separating from a syntax error, because it is the one a reader
// will not see. "need" for "needs" is valid JSON that silently drops an edge.
func unknownField(err error) (string, bool) {
	const marker = "json: unknown field "
	msg := err.Error()
	idx := strings.Index(msg, marker)
	if idx < 0 {
		return "", false
	}
	return strings.Trim(msg[idx+len(marker):], `"`), true
}

// LaneFor resolves the fast lane a root declares, and nil when it declares none. The
// kit root carries the built-in lane, because the kit ships no manifest of its own.
// Every other root declares its lane in the phase manifest. A root with no manifest,
// and a manifest with no lane array, declare no lane and keep the full-gate commit.
func LaneFor(root, kit string) ([]Phase, error) {
	if sameDirectory(root, kit) {
		return BenchkitLane(root, kit), nil
	}
	path := manifestPath(root)
	data, present, err := readManifest(path)
	if err != nil || !present {
		return nil, err
	}
	doc, err := decodeManifest(path, data)
	if err != nil {
		return nil, err
	}
	if len(doc.Lane) == 0 {
		return nil, nil
	}
	return validateManifest(path, root, doc.Lane)
}

// validateManifest converts declarations to phases, refusing the first defect it finds
// in declaration order so one manifest always produces one diagnostic. root is the tree
// under grade, which is where a declared dir is anchored.
func validateManifest(path, root string, declared []manifestPhase) ([]Phase, error) {
	phases := make([]Phase, 0, len(declared))
	declaredNames := make(map[string]bool, len(declared))
	for _, decl := range declared {
		if !validPhaseName(decl.Name) {
			return nil, defect(path, "invalid phase name", strconv.Quote(decl.Name))
		}
		if declaredNames[decl.Name] {
			return nil, defect(path, "duplicate phase name", decl.Name)
		}
		declaredNames[decl.Name] = true
		if len(decl.Argv) == 0 || decl.Argv[0] == "" {
			return nil, defect(path, "empty argv", decl.Name)
		}
		dir, contained := containedDir(decl.Dir)
		if !contained {
			return nil, defect(path, "escaping dir", decl.Dir)
		}
		phases = append(phases, Phase{
			Name:     decl.Name,
			Argv:     decl.Argv,
			Env:      envEntries(decl.Env),
			Optional: decl.Optional,
			Needs:    dedupe(decl.Needs),
			// The runner's own root is the kit checkout, which in a linked repo is a
			// different tree from the one under grade. Anchoring here, for an undeclared
			// dir too, keeps a manifest phase in the graded tree, where the directories
			// it names actually exist.
			Dir: filepath.Join(root, dir),
		})
	}
	// The scheduler reads a need naming no phase in the table as already satisfied.
	// Only the loader can tell a legitimately filtered edge apart from a typo.
	for _, phase := range phases {
		for _, need := range phase.Needs {
			if !declaredNames[need] {
				return nil, defect(path, "dangling needs edge", need)
			}
		}
	}
	if node := cycleNode(phases); node != "" {
		return nil, defect(path, "cyclic needs edge", node)
	}
	return phases, nil
}

// validPhaseName holds a name to bytes that survive the contracts it must remain
// addressable through. Those contracts are the "phase <name>:" summary lines and the
// "[name] " output prefixes. Whitespace or a control byte splits at least one of
// them.
func validPhaseName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if unicode.IsSpace(r) || r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

// containedDir reduces a declared dir to its cleaned root-relative form and reports
// whether it stays inside the root. Containment is lexical: the manifest already runs
// arbitrary argv from the graded tree, so this catches a mistake rather than an
// attacker. Resolving symlinks would buy nothing against either.
func containedDir(dir string) (string, bool) {
	if dir == "" {
		return "", true
	}
	if filepath.IsAbs(dir) {
		return "", false
	}
	clean := filepath.Clean(filepath.FromSlash(dir))
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", false
	}
	if clean == "." {
		return "", true
	}
	return clean, true
}

// envEntries renders the declared env object into the KEY=VALUE form Phase.Env carries.
// Map iteration order is random, so sorting is what keeps one manifest handing its phase
// the same environment on every run.
func envEntries(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	keys := make([]string, 0, len(env))
	for key := range env {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	entries := make([]string, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, key+"="+env[key])
	}
	return entries
}

func dedupe(names []string) []string {
	if len(names) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(names))
	unique := make([]string, 0, len(names))
	for _, name := range names {
		if seen[name] {
			continue
		}
		seen[name] = true
		unique = append(unique, name)
	}
	return unique
}

// cycleNode names a phase on the first cycle reachable in declaration order, or "" when
// the graph is acyclic. A cycle never reaches the scheduler. There it would settle as
// phases that silently never launch, which reports as skipped rather than as the
// manifest defect it is.
func cycleNode(phases []Phase) string {
	needs := make(map[string][]string, len(phases))
	for _, phase := range phases {
		needs[phase.Name] = phase.Needs
	}
	// An unrecorded phase is unvisited; open means it sits on the walk's current path,
	// which is what a back edge lands on.
	const (
		open = iota + 1
		closed
	)
	state := make(map[string]int, len(phases))
	var walk func(string) string
	walk = func(name string) string {
		switch state[name] {
		case open:
			return name
		case closed:
			return ""
		}
		state[name] = open
		for _, need := range needs[name] {
			if node := walk(need); node != "" {
				return node
			}
		}
		state[name] = closed
		return ""
	}
	for _, phase := range phases {
		if node := walk(phase.Name); node != "" {
			return node
		}
	}
	return ""
}

// defect is the one diagnostic shape every loader red uses: the manifest path, the
// defect class, and the element carrying the defect. A reader needs all three to fix it,
// and a generic message would leave them hunting through the file.
func defect(path, class, element string) error {
	return fmt.Errorf("gate: %s: %s: %s", path, class, element)
}
