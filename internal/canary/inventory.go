// Package canary binds mutation fixtures to the production checks that grade them.
package canary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

const absentHarnessMessage = "canary fixture inventory is empty"

const (
	checkFileName = "CHECK"
)

// PhaseManifestPath is where a graded root declares its phase table.
const PhaseManifestPath = ".bench/phases.json"

const (
	PhaseGofmt = "gofmt"
	PhaseVet   = "vet"
	PhaseTest  = "test"
	PhaseRace  = "race"
)

// IsConformanceFamily reports whether dir is a family rather than a flat fixture.
func IsConformanceFamily(dir string) bool {
	holds, err := holdsExpect(dir)
	return err == nil && !holds
}

// UnboundConformanceFamilies reports family directories without a registry owner.
func UnboundConformanceFamilies(kitRoot string) []string {
	entries, err := os.ReadDir(filepath.Join(kitRoot, "tests", "canary"))
	if err != nil {
		return nil
	}
	var diagnostics []string
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() {
			continue
		}
		familyDir := filepath.Join(kitRoot, "tests", "canary", name)
		holds, err := holdsExpect(familyDir)
		if err != nil {
			diagnostics = append(diagnostics, err.Error())
			continue
		}
		if holds {
			continue
		}
		if _, found := registry.FamilyCheck(name); !found {
			diagnostics = append(diagnostics, fmt.Sprintf("canary conformance family %q is bound to no conformance check; add it to the registry family table so its fixtures resolve a conformance-check binding", name))
		}
	}
	return diagnostics
}

// Fixture is one immutable mutation input and its owning production check.
type Fixture struct {
	Dir, Family, Check string
}

type fixtureRecord struct {
	dir, family string
}

// Fixtures returns the complete fixture inventory keyed by globally unique base name.
func Fixtures(canaryDir string) (map[string]Fixture, error) {
	records, err := discoverFixtures(canaryDir)
	if err != nil {
		return nil, err
	}
	result := make(map[string]Fixture, len(records))
	for _, record := range records {
		_, check, err := fixtureCheck(record.dir)
		if err != nil {
			return nil, err
		}
		result[filepath.Base(record.dir)] = Fixture{
			Dir: record.dir, Family: record.family, Check: fixtureScope(record.family, check),
		}
	}
	return result, nil
}

var grammar = usage.Grammar{Cmd: "bench canary", Help: "usage: bench canary [root]", MaxArgs: 1}

// Run validates the production owner inventory for every fixture.
func Run(args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
		} else {
			fmt.Fprintln(stderr, line)
		}
		return code
	}
	root := ""
	if len(parsed.Positionals) == 1 {
		root = parsed.Positionals[0]
	} else {
		var err error
		root, err = git.Root()
		if err != nil {
			fmt.Fprintln(stderr, toon.NotInRepo())
			return 1
		}
	}
	result, err := Inventory(root)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "canary inventory ok (%d fixture bindings)\n", len(result.Owners))
	return 0
}

// Inventory returns the validated fixture inventory decision for root.
func Inventory(root string) (Selection, error) {
	found, err := Fixtures(filepath.Join(root, "tests", "canary"))
	if err != nil {
		return Selection{}, err
	}
	fixtures := make([]Fixture, 0, len(found))
	for _, fixture := range found {
		fixtures = append(fixtures, fixture)
	}
	result := Select(fixtures)
	if !result.Accepted {
		return result, errors.New(strings.Join(result.Diagnostics, "\n"))
	}
	return result, nil
}

// FixturePins maps every repository path the live fixture inventory pins to the
// repo-relative fixture directories that pin it. A fixture pins a path when its
// BASE list names it, when its files/ overlay carries it, or when its
// MUTATE.json anchors a mutation in it. The BASE walk follows @ includes, so an
// included list pins for the fixture that includes it.
//
// The enumeration reads the inventory rather than a copied list, so a fixture
// that moves carries its pins with it. A root with no fixture inventory pins
// nothing; that is an answer, not a fault.
func FixturePins(root string) (map[string][]string, error) {
	records, err := discoverFixtures(filepath.Join(root, "tests", "canary"))
	if err != nil {
		if err.Error() == absentHarnessMessage {
			return map[string][]string{}, nil
		}
		return nil, err
	}
	pins := map[string][]string{}
	for _, record := range records {
		dir, relErr := filepath.Rel(root, record.dir)
		if relErr != nil {
			return nil, relErr
		}
		dir = filepath.ToSlash(dir)
		paths, pinErr := pinnedPaths(root, record.dir)
		if pinErr != nil {
			return nil, pinErr
		}
		for _, path := range paths {
			if !holdsString(pins[path], dir) {
				pins[path] = append(pins[path], dir)
			}
		}
	}
	for path := range pins {
		sort.Strings(pins[path])
	}
	return pins, nil
}

// pinnedPaths is every repository path one fixture directory pins, in BASE,
// overlay, and mutation order. It reuses the BASE reader the materializer
// drives, so the pinned set and the materialized set cannot disagree.
func pinnedPaths(root, fixture string) ([]string, error) {
	var paths []string
	basePath := filepath.Join(fixture, "BASE")
	if regularFile(basePath) {
		rels, err := basePaths(root, basePath)
		if err != nil {
			return nil, err
		}
		paths = append(paths, rels...)
	}
	filesDir := filepath.Join(fixture, filesDirName)
	if info, err := os.Stat(filesDir); err == nil && info.IsDir() {
		if err := filepath.WalkDir(filesDir, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			rel, relErr := filepath.Rel(filesDir, path)
			if relErr != nil {
				return relErr
			}
			paths = append(paths, filepath.ToSlash(restoredFixturePath(rel)))
			return nil
		}); err != nil {
			return nil, err
		}
	}
	mutatePath := filepath.Join(fixture, "MUTATE.json")
	if !regularFile(mutatePath) {
		return paths, nil
	}
	data, err := os.ReadFile(mutatePath)
	if err != nil {
		return nil, err
	}
	var mutations []fixtureMutation
	if err := json.Unmarshal(data, &mutations); err != nil {
		return nil, err
	}
	for _, mutation := range mutations {
		paths = append(paths, filepath.ToSlash(filepath.Clean(mutation.Path)))
	}
	return paths, nil
}

func holdsString(list []string, want string) bool {
	for _, value := range list {
		if value == want {
			return true
		}
	}
	return false
}

func discoverFixtures(dir string) ([]fixtureRecord, error) {
	families, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New(absentHarnessMessage)
	}
	var result []fixtureRecord
	seen := map[string]bool{}
	add := func(record fixtureRecord) error {
		name := filepath.Base(record.dir)
		if seen[name] {
			return fmt.Errorf("canary fixture name %q appears in multiple families; base names must be globally unique", name)
		}
		seen[name] = true
		result = append(result, record)
		return nil
	}
	for _, family := range families {
		if !family.IsDir() {
			continue
		}
		name := family.Name()
		familyDir := filepath.Join(dir, name)
		holds, err := holdsExpect(familyDir)
		if err != nil {
			return nil, err
		}
		if holds {
			if err := add(fixtureRecord{dir: familyDir}); err != nil {
				return nil, err
			}
			continue
		}
		entries, err := os.ReadDir(familyDir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				if err := add(fixtureRecord{dir: filepath.Join(familyDir, entry.Name()), family: name}); err != nil {
					return nil, err
				}
			}
		}
	}
	if len(result) == 0 {
		return nil, errors.New(absentHarnessMessage)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].dir < result[j].dir })
	return result, nil
}

func fixtureCheck(dir string) (registry.Tier, string, error) {
	name, present, err := readMarker(dir, "CHECK")
	if err != nil {
		return "", "", err
	}
	if !present {
		return registry.Dev, "", nil
	}
	if name == "" {
		return "", "", fmt.Errorf("canary fixture %q has an empty CHECK file", filepath.Base(dir))
	}
	check, found := registry.Find(name)
	if !found {
		return "", "", fmt.Errorf("canary fixture %q names unknown check %q", filepath.Base(dir), name)
	}
	return check.Tier, check.Name, nil
}

func fixtureScope(family, check string) string {
	if check != "" {
		return check
	}
	if owner, found := registry.FamilyCheck(family); found {
		return owner
	}
	if family != "" {
		return "project:" + family
	}
	return ""
}

func readMarker(dir, marker string) (string, bool, error) {
	markerPath := filepath.Join(dir, marker)
	read := bounds.Classify(markerPath, bounds.ControlRecordLimit)
	switch read.State {
	case bounds.StateAbsent:
		return "", false, nil
	case bounds.StateEmpty, bounds.StateParsed:
		return strings.TrimSpace(string(read.Data)), true, nil
	default:
		return "", false, fmt.Errorf("canary fixture marker %s cannot be read: %s", markerPath, read.Reason)
	}
}

func holdsExpect(dir string) (bool, error) {
	markerPath := filepath.Join(dir, "EXPECT")
	classified := bounds.Classify(markerPath, bounds.ControlRecordLimit)
	switch classified.State {
	case bounds.StateAbsent:
		return false, nil
	case bounds.StateEmpty, bounds.StateParsed:
		return true, nil
	default:
		return false, fmt.Errorf("canary fixture marker %s cannot be read: %s", markerPath, classified.Reason)
	}
}

// MaterializeFixture copies one fixture mutation into dst and restores encoded dot paths.
func MaterializeFixture(src, dst string) error {
	if err := copyTree(src, dst); err != nil {
		return err
	}
	return restoreDotSegments(dst)
}

func materialize(src, dst string) error {
	return MaterializeFixture(src, dst)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(source string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(src, source)
		if err != nil || relative == "." {
			return err
		}
		target := filepath.Join(dst, relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		data, err := os.ReadFile(source)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}

func restoreDotSegments(root string) error {
	var directories []string
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "dot-") {
			directories = append(directories, path)
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Slice(directories, func(i, j int) bool { return len(directories[i]) > len(directories[j]) })
	for _, old := range directories {
		if err := os.Rename(old, filepath.Join(filepath.Dir(old), "."+strings.TrimPrefix(filepath.Base(old), "dot-"))); err != nil {
			return err
		}
	}
	return nil
}
