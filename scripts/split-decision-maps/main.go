// Command split-decision-maps migrates every inline decision map to the split
// shape: an index plus one file per ticket, with each map-owned asset in the map's
// own topic folder. It is a one-shot migration; the contract ticket deletes it.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/maps"
)

const orphanDir = "docs/research"

type migration struct {
	root string
	// files holds the post-split content of every in-scope Markdown file, keyed by
	// repository-relative slash path. The plan is computed against this virtual tree
	// so a reference that moves into a ticket file is reported at its new home.
	files    map[string]string
	written  []string
	splits   map[string][]string
	moves    map[string]string
	rewrites map[string][][2]string
	ignore   string
}

func main() {
	apply := flag.Bool("apply", false, "write the migration instead of printing the plan")
	flag.Parse()
	root, err := os.Getwd()
	if err != nil {
		fail(err)
	}
	m := &migration{root: root, files: map[string]string{}, splits: map[string][]string{},
		moves: map[string]string{}, rewrites: map[string][][2]string{}}
	if err := m.plan(); err != nil {
		fail(err)
	}
	m.report()
	if !*apply {
		return
	}
	if err := m.apply(); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "split-decision-maps:", err)
	os.Exit(1)
}

func (m *migration) plan() error {
	scope, err := m.scope()
	if err != nil {
		return err
	}
	for _, file := range scope {
		content, err := os.ReadFile(filepath.Join(m.root, filepath.FromSlash(file)))
		if err != nil {
			return err
		}
		m.files[file] = string(content)
	}
	candidates, err := maps.DiscoverDecisionMapCandidates(m.root)
	if err != nil {
		return err
	}
	var indexes []string
	for _, candidate := range candidates {
		topic := strings.TrimSuffix(candidate.Path, ".md")
		if info, err := os.Stat(filepath.Join(m.root, filepath.FromSlash(topic), "tickets")); err == nil && info.IsDir() {
			continue
		}
		if err := m.split(candidate.Path); err != nil {
			return err
		}
		indexes = append(indexes, candidate.Path)
	}
	if err := m.planAssets(indexes); err != nil {
		return err
	}
	m.planRewrites()
	return m.planIgnore()
}

// scope is every tracked Markdown file the migration may read or rewrite. A spec's
// own spec.md and tickets/ stay out: a build may not edit a spec.
func (m *migration) scope() ([]string, error) {
	out, err := exec.Command("git", "-C", m.root, "ls-files", "-z", "*.md").Output()
	if err != nil {
		return nil, err
	}
	var scope []string
	for _, file := range strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00") {
		dir := path.Dir(file)
		inSpec := strings.HasPrefix(file, "specs/") && strings.Contains(file, "/decisions/")
		if dir == "docs" || strings.HasPrefix(file, "docs/") || strings.HasPrefix(file, "decisions/") || inSpec {
			scope = append(scope, file)
		}
	}
	sort.Strings(scope)
	return scope, nil
}

func (m *migration) split(index string) error {
	doc, err := cutDoc(index, []byte(m.files[index]))
	if err != nil {
		return err
	}
	topic := strings.TrimSuffix(path.Base(index), ".md")
	m.files[index] = renderIndex(topic, doc)
	for _, ticket := range doc.tickets {
		file := path.Join(strings.TrimSuffix(index, ".md"), "tickets", ticket.id+".md")
		m.files[file] = renderTicket(ticket)
		m.splits[index] = append(m.splits[index], file)
	}
	return nil
}

// planAssets assigns each asset in a flat assets folder to the one map in that
// folder's directory that names it. A map that names an asset only in prose owns it
// as surely as one that names it in Sources, so ownership reads the whole map.
func (m *migration) planAssets(indexes []string) error {
	dirs := map[string][]string{}
	for _, index := range indexes {
		dirs[path.Dir(index)] = append(dirs[path.Dir(index)], index)
	}
	for dir, dirIndexes := range dirs {
		assets := filepath.Join(m.root, filepath.FromSlash(dir), "assets")
		entries, err := os.ReadDir(assets)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			old := path.Join(dir, "assets", entry.Name())
			var owners []string
			for _, index := range dirIndexes {
				if m.names(index, old) {
					owners = append(owners, index)
				}
			}
			switch len(owners) {
			case 0:
				m.moves[old] = path.Join(orphanDir, entry.Name())
			case 1:
				m.moves[old] = path.Join(strings.TrimSuffix(owners[0], ".md"), "assets", entry.Name())
			default:
				return fmt.Errorf("%s: named by %s", old, strings.Join(owners, ", "))
			}
		}
	}
	return nil
}

// names reports whether the map at index, or any file the split cut out of it, holds
// the path.
func (m *migration) names(index, target string) bool {
	if strings.Contains(m.files[index], target) {
		return true
	}
	for _, file := range m.splits[index] {
		if strings.Contains(m.files[file], target) {
			return true
		}
	}
	return false
}

func (m *migration) planRewrites() {
	for file, content := range m.files {
		for old, dest := range m.moves {
			if strings.Contains(content, old) {
				m.rewrites[file] = append(m.rewrites[file], [2]string{old, dest})
				content = strings.ReplaceAll(content, old, dest)
			}
		}
		m.files[file] = content
	}
}

func (m *migration) planIgnore() error {
	file := filepath.Join(m.root, ".gitignore")
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	const old = "\nresearch/\n"
	const anchored = "\n# Anchor the rule to the root shift-scratch folder, so docs/research/ stays tracked.\n/research/\n"
	if !strings.Contains(string(content), old) {
		return fmt.Errorf(".gitignore: no bare research/ rule to replace")
	}
	m.ignore = strings.Replace(string(content), old, anchored, 1)
	return nil
}

func (m *migration) report() {
	fmt.Println("maps:")
	for _, index := range sortedKeys(m.splits) {
		fmt.Printf("  %s -> index + %d ticket files\n", index, len(m.splits[index]))
		for _, file := range m.splits[index] {
			fmt.Printf("    %s\n", file)
		}
	}
	fmt.Println("assets:")
	for _, old := range sortedKeys(m.moves) {
		fmt.Printf("  %s -> %s\n", old, m.moves[old])
	}
	fmt.Println("rewrites:")
	for _, file := range sortedKeys(m.rewrites) {
		fmt.Printf("  %s\n", file)
		for _, pair := range m.rewrites[file] {
			fmt.Printf("    %s -> %s\n", pair[0], pair[1])
		}
	}
	fmt.Println("ignore:")
	fmt.Println("  .gitignore: research/ -> /research/ under a shift-scratch comment")
}

func sortedKeys[V any](in map[string]V) []string {
	keys := make([]string, 0, len(in))
	for key := range in {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// apply writes content before it moves any asset, so every rewrite lands on the file
// at the path the plan computed.
func (m *migration) apply() error {
	for _, index := range sortedKeys(m.splits) {
		for _, file := range append([]string{index}, m.splits[index]...) {
			if err := m.write(file); err != nil {
				return err
			}
		}
	}
	for _, file := range sortedKeys(m.rewrites) {
		if err := m.write(file); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(m.root, ".gitignore"), []byte(m.ignore), 0o644); err != nil {
		return err
	}
	for _, old := range sortedKeys(m.moves) {
		dest := m.moves[old]
		if err := os.MkdirAll(filepath.Join(m.root, filepath.FromSlash(path.Dir(dest))), 0o755); err != nil {
			return err
		}
		if out, err := exec.Command("git", "-C", m.root, "mv", old, dest).CombinedOutput(); err != nil {
			return fmt.Errorf("git mv %s: %s", old, strings.TrimSpace(string(out)))
		}
	}
	for _, old := range sortedKeys(m.moves) {
		os.Remove(filepath.Join(m.root, filepath.FromSlash(path.Dir(old))))
	}
	return m.stage()
}

func (m *migration) write(file string) error {
	full := filepath.Join(m.root, filepath.FromSlash(file))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(full, []byte(m.files[file]), 0o644); err != nil {
		return err
	}
	m.written = append(m.written, file)
	return nil
}

// stage adds the new ticket files, so git status reports the migration as one
// reviewable change rather than a wall of untracked paths. A rewritten asset was
// written at its old path and then moved, so stage names where it landed.
func (m *migration) stage() error {
	args := []string{"-C", m.root, "add", "--"}
	for _, file := range m.written {
		if dest, moved := m.moves[file]; moved {
			file = dest
		}
		args = append(args, file)
	}
	if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("git add: %s", strings.TrimSpace(string(out)))
	}
	return nil
}
