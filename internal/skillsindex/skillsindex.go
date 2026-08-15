// Package skillsindex is the one reader of the generated skills index in
// .bench/BENCH-reference.md. Source of truth is each skill's frontmatter (`index:`
// trigger, optional `index-note:` suffix) plus the consumer-payload allowlist's
// kit-only rows; the block between the bench:skills-index markers is derived from
// them, never hand-maintained.
//
// The index used to be derived twice — a shell generator and a conformance check that
// re-implemented the same frontmatter read, kit-only marking, line format, and
// diagnostics. Every consumer now goes through this package, so a change to the line
// shape or the fence rule lands once.
package skillsindex

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	kitpayload "github.com/gibbonmi/bench"
)

const (
	startMarker = "<!-- bench:skills-index:start -->"
	endMarker   = "<!-- bench:skills-index:end -->"

	referenceRel = ".bench/BENCH-reference.md"
	allowlistRel = ".bench/consumer-payload.json"

	// regenerateHint names the operator's regenerator on every diagnostic a rerun of
	// that verb would clear. The three diagnostics that carry no hint are the ones a
	// regeneration cannot fix: a skill missing its trigger, and a reference file that
	// is absent or unmarked.
	regenerateHint = "(regenerate: bench skills-index --write)"
)

// ErrAllowlistUnreadable is Write's refusal when the allowlist exists but does not
// parse. Generating anyway would commit a block with every kit-only marker missing
// and attribute the resulting drift to the skills rather than to the reader that
// dropped them, so the write is refused before the file is touched. Check keeps the
// opposite posture — it marks nothing and adds no diagnostic — because the gate must
// not gain a red the index itself did not earn.
var ErrAllowlistUnreadable = errors.New(allowlistRel + " unreadable: kit-only marking unresolved (write refused)")

// Entry is one skill's index row as data. Trigger is empty for a skill that declares
// no `index:` frontmatter; such a skill renders no line but still earns a diagnostic,
// which is why it stays in the enumeration rather than being dropped here.
type Entry struct {
	Name    string
	Trigger string
	Note    string
	KitOnly bool
}

// Indexed reports whether the entry contributes a line to the generated block.
func (e Entry) Indexed() bool { return e.Trigger != "" }

// Line renders the entry's index line.
func (e Entry) Line() string {
	line := fmt.Sprintf("- %s → `.agents/skills/%s/SKILL.md`", e.Trigger, e.Name)
	if e.KitOnly {
		line += " (kit-only)"
	}
	if e.Note != "" {
		line += " + " + e.Note
	}
	return line
}

// Entries enumerates root's skills in alphabetical directory order, skipping a skill
// whose name has a same-named .agents/commands/<name>.md — a command adapter, which
// the slash menu already surfaces. The error is non-nil only when the allowlist is
// present and unparseable; entries come back unmarked in that case so a caller can
// choose between refusing (Write) and today's posture (Check).
func Entries(root string) ([]Entry, error) {
	kitOnly, err := kitOnlySources(root)
	files, _ := filepath.Glob(filepath.Join(root, ".agents", "skills", "*", "SKILL.md"))
	sort.Strings(files)
	var entries []Entry
	for _, file := range files {
		name := filepath.Base(filepath.Dir(file))
		if exists(filepath.Join(root, ".agents", "commands", name+".md")) {
			continue
		}
		entries = append(entries, Entry{
			Name:    name,
			Trigger: FrontmatterField(file, "index"),
			Note:    FrontmatterField(file, "index-note"),
			KitOnly: kitOnly[".agents/skills/"+name],
		})
	}
	return entries, err
}

// FrontmatterField returns the first value of key inside path's leading `---` fence.
// Nothing after the closing fence counts, so a body line that happens to start with
// the key is not frontmatter.
func FrontmatterField(path, key string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	fence := 0
	prefix := key + ":"
	for _, line := range strings.Split(string(data), "\n") {
		if line == "---" {
			fence++
			continue
		}
		if fence == 1 && strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
		if fence > 1 {
			return ""
		}
	}
	return ""
}

// Check grades root's committed block against the generated one and returns the
// attributed diagnostics in skill-alphabetical order, with the unattributed block
// line last. The order is part of the contract: the gate and `bench skills-index`
// print the same sequence.
func Check(root string) []string {
	entries, _ := Entries(root)
	// A skill can earn more than one diagnostic — no trigger and a stale committed
	// entry is the colliding case — so this accumulates per skill rather than
	// keying one string per name, which would silently drop the first of the pair.
	attributed := map[string][]string{}
	var expected []string
	for _, entry := range entries {
		if !entry.Indexed() {
			attributed[entry.Name] = append(attributed[entry.Name], fmt.Sprintf("skill '%s' missing index: frontmatter (the skills index is generated)", entry.Name))
			continue
		}
		expected = append(expected, entry.Line())
	}

	ref := readIfExists(filepath.Join(root, filepath.FromSlash(referenceRel)))
	if ref == "" {
		return append(ordered(attributed), referenceRel+" missing (skills index unverifiable)")
	}
	actual, ok := markerBlock(ref)
	if !ok {
		return append(ordered(attributed), referenceRel+" skills-index markers missing (bench:skills-index)")
	}
	if strings.Join(expected, "\n") == strings.Join(actual, "\n") {
		return ordered(attributed)
	}

	expectedByName := byName(expected)
	actualByName := byName(actual)
	drifted := false
	for name, line := range expectedByName {
		if actualByName[name] == line {
			continue
		}
		if _, ok := actualByName[name]; ok {
			attributed[name] = append(attributed[name], fmt.Sprintf("skills index entry for '%s' drifted from its frontmatter %s", name, regenerateHint))
		} else {
			attributed[name] = append(attributed[name], fmt.Sprintf("skills index missing entry for skill '%s' %s", name, regenerateHint))
		}
		drifted = true
	}
	for name := range actualByName {
		if _, ok := expectedByName[name]; !ok {
			attributed[name] = append(attributed[name], fmt.Sprintf("skills index entry '%s' has no indexed .agents/skills/%s on disk %s", name, name, regenerateHint))
			drifted = true
		}
	}
	diags := ordered(attributed)
	if !drifted {
		diags = append(diags, "skills index block drifted from generated form "+regenerateHint)
	}
	return diags
}

// Write regenerates root's block in place. The reference file is a shipped 0644 asset
// whose mode the release-evidence registry checks, so the temp file's 0600 is
// corrected before the rename replaces it.
func Write(root string) error {
	entries, err := Entries(root)
	if err != nil {
		return ErrAllowlistUnreadable
	}
	path := filepath.Join(root, filepath.FromSlash(referenceRel))
	ref, readErr := os.ReadFile(path)
	if readErr != nil {
		return fmt.Errorf("%s missing (skills index unverifiable)", referenceRel)
	}
	var block []string
	for _, entry := range entries {
		if entry.Indexed() {
			block = append(block, entry.Line())
		}
	}
	rewritten, ok := replaceMarkerBlock(string(ref), block)
	if !ok {
		return fmt.Errorf("%s skills-index markers missing (bench:skills-index)", referenceRel)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".skills-index-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.WriteString(rewritten); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Chmod(name, 0o644); err != nil {
		os.Remove(name)
		return err
	}
	return os.Rename(name, path)
}

// kitOnlySources reads the allowlist's withheld skill sources. The allowlist is the
// one source of who receives an asset, so the index marks the same rows instead of
// naming the kit-only skills a second time; it is JSON, so a row that orders
// "audience" before "source" resolves like any other. A tree with no allowlist has
// nothing withheld, which is also what the generator concludes.
func kitOnlySources(root string) (map[string]bool, error) {
	raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(allowlistRel)))
	if err != nil {
		return nil, nil
	}
	var rows []kitpayload.PayloadRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, ErrAllowlistUnreadable
	}
	out := map[string]bool{}
	for _, source := range kitpayload.PayloadKitOnlyPrefixes(rows) {
		out[source] = true
	}
	return out, nil
}

// ordered flattens the per-skill diagnostics into the contract's one sequence:
// skill-alphabetical, and within a skill the order they were emitted — the missing
// trigger before any block drift, which is what the pre-module check printed.
func ordered(attributed map[string][]string) []string {
	names := make([]string, 0, len(attributed))
	for name := range attributed {
		names = append(names, name)
	}
	sort.Strings(names)
	var out []string
	for _, name := range names {
		out = append(out, attributed[name]...)
	}
	return out
}

var indexLineRe = regexp.MustCompile(`\.agents/skills/([a-z0-9-]+)/SKILL\.md`)

func byName(lines []string) map[string]string {
	out := map[string]string{}
	for _, line := range lines {
		if m := indexLineRe.FindStringSubmatch(line); len(m) == 2 {
			out[m[1]] = line
		}
	}
	return out
}

func markerBlock(text string) ([]string, bool) {
	inBlock := false
	var out []string
	for _, line := range strings.Split(text, "\n") {
		switch {
		case line == endMarker:
			return out, inBlock
		case inBlock:
			out = append(out, line)
		case line == startMarker:
			inBlock = true
		}
	}
	return nil, false
}

func replaceMarkerBlock(text string, block []string) (string, bool) {
	var out []string
	skipping, replaced := false, false
	for _, line := range strings.Split(text, "\n") {
		switch {
		case line == startMarker:
			out = append(out, line)
			out = append(out, block...)
			skipping, replaced = true, true
			continue
		case line == endMarker:
			skipping = false
		}
		if !skipping {
			out = append(out, line)
		}
	}
	if !replaced || skipping {
		return "", false
	}
	return strings.Join(out, "\n"), true
}

func readIfExists(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}
