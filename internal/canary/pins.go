package canary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"

	"github.com/gibbonmi/bench/internal/bounds"
)

// PinRule names the inventory rule by which one fixture pins a path.
type PinRule string

const (
	// PinBase is a path the fixture's own BASE list names.
	PinBase PinRule = "base"
	// PinInclude is a path an @ include of the BASE list names.
	PinInclude PinRule = "include"
	// PinOverlay is a path the fixture's files/ overlay carries.
	PinOverlay PinRule = "overlay"
	// PinMutation is a path a MUTATE.json entry anchors a mutation in.
	PinMutation PinRule = "mutation"
)

// Pin is one fact by which a fixture pins a repository path: the repository-relative
// fixture directory, the file that names the path, the line there, and the rule that file
// follows. An overlay file pins its path as a whole, so its line is zero.
type Pin struct {
	Path, Fixture, Owner string
	Line                 int
	Rule                 PinRule
}

// Pins is every pin of root's live fixture inventory, in fixture order and then in BASE,
// overlay, and mutation order. The BASE walk follows @ includes, so an included list pins
// for the fixture that includes it, and its pin names the included file.
//
// The inventory directory is classified without a link followed before it is read. An
// absent or empty directory pins nothing; that is an answer, not a fault. A link, a
// special file, or an unreadable directory refuses the enumeration, so a failed live
// inventory never reads as an empty one.
func Pins(root string) ([]Pin, error) {
	inventory := bounds.ClassifyDirNoFollow(Dir(root))
	switch inventory.State {
	case bounds.StateAbsent, bounds.StateEmpty:
		return nil, nil
	case bounds.StateParsed:
	default:
		return nil, fmt.Errorf("canary fixture inventory %s refused: %s %s", Dir(root), inventory.State, inventory.Reason)
	}
	records, err := rootRecords(root)
	if err != nil {
		return nil, err
	}
	var pins []Pin
	for _, record := range records {
		dir, relErr := filepath.Rel(root, record.dir)
		if relErr != nil {
			return nil, relErr
		}
		found, pinErr := fixturePins(root, record.dir, filepath.ToSlash(dir))
		if pinErr != nil {
			return nil, pinErr
		}
		pins = append(pins, found...)
	}
	return pins, nil
}

// FixturePins maps every repository path the live fixture inventory pins to the sorted
// repo-relative fixture directories that pin it. It projects Pins, so it refuses whatever
// Pins refuses.
//
// The enumeration reads the inventory rather than a copied list, so a fixture
// that moves carries its pins with it.
func FixturePins(root string) (map[string][]string, error) {
	found, err := Pins(root)
	if err != nil {
		return nil, err
	}
	pins := map[string][]string{}
	for _, pin := range found {
		if !slices.Contains(pins[pin.Path], pin.Fixture) {
			pins[pin.Path] = append(pins[pin.Path], pin.Fixture)
		}
	}
	for path := range pins {
		sort.Strings(pins[path])
	}
	return pins, nil
}

// PinnedPaths is every repository path one fixture directory pins, in BASE,
// overlay, and mutation order. The fixture directory is an absolute path or a
// path relative to the working directory; the root resolves the BASE includes.
//
// A caller that reads one fixture, such as a failure message that narrows a
// diagnostic list to the pinned lines, uses this accessor. Pins builds on the same
// reader, so the two answers cannot disagree.
func PinnedPaths(root, fixtureDir string) ([]string, error) {
	pins, err := fixturePins(root, fixtureDir, filepath.ToSlash(fixtureDir))
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(pins))
	for i, pin := range pins {
		paths[i] = pin.Path
	}
	return paths, nil
}

// fixturePins is every pin of one fixture directory, in BASE, overlay, and mutation order.
// rel is the directory's repository-relative spelling, which each pin reports. It reuses
// the BASE reader the materializer drives, so the pinned set and the materialized set
// cannot disagree.
func fixturePins(root, fixture, rel string) ([]Pin, error) {
	var pins []Pin
	basePath := filepath.Join(fixture, "BASE")
	if regularFile(basePath) {
		entries, err := baseEntries(root, basePath, path.Join(rel, "BASE"))
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			rule := PinBase
			if entry.included {
				rule = PinInclude
			}
			pins = append(pins, Pin{Path: entry.path, Fixture: rel, Owner: entry.list, Line: entry.line, Rule: rule})
		}
	}
	filesDir := filepath.Join(fixture, filesDirName)
	if info, err := os.Stat(filesDir); err == nil && info.IsDir() {
		if err := filepath.WalkDir(filesDir, func(file string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			overlay, relErr := filepath.Rel(filesDir, file)
			if relErr != nil {
				return relErr
			}
			pins = append(pins, Pin{
				Path:    filepath.ToSlash(restoredFixturePath(overlay)),
				Fixture: rel,
				Owner:   path.Join(rel, filesDirName, filepath.ToSlash(overlay)),
				Rule:    PinOverlay,
			})
			return nil
		}); err != nil {
			return nil, err
		}
	}
	mutatePath := filepath.Join(fixture, "MUTATE.json")
	if !regularFile(mutatePath) {
		return pins, nil
	}
	data, err := os.ReadFile(mutatePath)
	if err != nil {
		return nil, err
	}
	var mutations []fixtureMutation
	if err := json.Unmarshal(data, &mutations); err != nil {
		return nil, err
	}
	lines := elementLines(data, len(mutations))
	for i, mutation := range mutations {
		pins = append(pins, Pin{Path: filepath.ToSlash(filepath.Clean(mutation.Path)), Fixture: rel, Owner: path.Join(rel, "MUTATE.json"), Line: lines[i], Rule: PinMutation})
	}
	return pins, nil
}

// elementLines answers the line of each top-level element of a JSON array that already
// decoded, so a mutation pin cites the entry that names its path.
func elementLines(data []byte, count int) []int {
	lines := make([]int, count)
	decoder := json.NewDecoder(bytes.NewReader(data))
	if _, err := decoder.Token(); err != nil {
		return lines
	}
	for i := range lines {
		start := int(decoder.InputOffset())
		var element json.RawMessage
		if err := decoder.Decode(&element); err != nil {
			return lines
		}
		if at := bytes.Index(data[start:], element); at >= 0 {
			start += at
		}
		lines[i] = bytes.Count(data[:start], []byte("\n")) + 1
	}
	return lines
}
