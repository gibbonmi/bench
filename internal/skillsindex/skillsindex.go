// Package skillsindex is the one reader of the generated skills index in
// .bench/BENCH-reference.md. Source of truth is each skill's frontmatter (`index:`
// trigger, optional `index-note:` suffix) plus the consumer-payload allowlist's
// kit-only rows; the block between the bench:skills-index markers is derived from
// them, never hand-maintained.
//
// Every consumer reads the index here — the conformance check behind the gate and the
// `bench skills-index` verb — so the line shape, the kit-only marking, and the fence
// rule have one owner.
package skillsindex

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	kitpayload "github.com/gibbonmi/bench"
	"github.com/gibbonmi/bench/internal/bounds"
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

// allowlistDiagnostic is the one wording for a present allowlist the reader cannot
// resolve. Both callers refuse on it: generating would commit a block with every
// kit-only marker missing, and grading against that block would attribute the reader's
// defect to the skills.
const allowlistDiagnostic = allowlistRel + " unreadable: kit-only marking unresolved"

// markersDiagnostic is the one wording for a reference that names no single generated
// block. Check grades against it and Write refuses on it, so the two never drift.
const markersDiagnostic = referenceRel + " skills-index markers missing (bench:skills-index)"

// ErrAllowlistUnreadable is Write's refusal when the allowlist is present but does not
// resolve, so the write is refused before the reference file is touched.
var ErrAllowlistUnreadable = errors.New(allowlistDiagnostic + " (write refused)")

// Entry is one skill's index row as data. Trigger is empty for a skill that declares
// no `index:` frontmatter; such a skill renders no line but still earns a diagnostic,
// which is why it stays in the enumeration rather than being dropped here.
type Entry struct {
	Name    string
	Trigger string
	Note    string
	KitOnly bool
	// Refusal is the classifier's diagnostic for a skill whose SKILL.md is not
	// trustworthy bytes at all — a link, a special file, oversized, or not UTF-8. It is
	// distinct from an empty Trigger because a skill nobody could read has not declined
	// to be indexed; it has to be named as refused so the operator fixes the file rather
	// than adding frontmatter to something that cannot hold it.
	Refusal string
	// Orphan marks a real skill directory with no SKILL.md at all. Absence is not bad
	// bytes, so it carries no Refusal reason, and it is not empty bytes either: those
	// are readable and grade as a missing trigger. A skill source nobody wrote is a
	// producer defect, so it is named rather than dropped.
	Orphan bool
}

// Indexed reports whether the entry contributes a line to the generated block. A refused
// skill renders no line: bytes the classifier would not vouch for cannot author one.
func (e Entry) Indexed() bool { return !e.Orphan && e.Refusal == "" && e.Trigger != "" }

// skillPathFormat is the skill source path an index line carries: Line formats it,
// indexLineRe matches it.
const skillPathFormat = ".agents/skills/%s/SKILL.md"

// Line renders the entry's index line.
func (e Entry) Line() string {
	line := fmt.Sprintf("- %s → `"+skillPathFormat+"`", e.Trigger, e.Name)
	if e.KitOnly {
		line += " (kit-only)"
	}
	if e.Note != "" {
		line += " + " + e.Note
	}
	return line
}

var indexLineRe = regexp.MustCompile(strings.Replace(regexp.QuoteMeta(skillPathFormat), "%s", "([a-z0-9-]+)", 1))

// skillName returns the skill a committed index line names, or "" for a line that is
// not one. Expected lines carry their skill in the Entry that rendered them, so this
// reader is for committed text only — text no Entry can vouch for.
func skillName(line string) string {
	if m := indexLineRe.FindStringSubmatch(line); len(m) == 2 {
		return m[1]
	}
	return ""
}

// Entries enumerates root's skills in alphabetical directory order, skipping a skill
// whose name has a same-named .agents/commands/<name>.md — a command adapter, which
// the slash menu already surfaces. The error is non-nil only when the allowlist is
// present and does not resolve; entries come back unmarked in that case, and both
// callers refuse on the error rather than acting on unmarked rows.
func Entries(root string) ([]Entry, error) {
	kitOnly, err := kitOnlySources(root)
	// Literal child enumeration, not a glob: the root is operator-supplied text, and a
	// pattern walk would read a `[ ] * ?` in it as syntax and report a repository with
	// no skills at all. A child that is not a real directory — a link to one included —
	// is not a skill source, so it never reaches the classifier.
	skillsDir := bounds.ClassifyDir(filepath.Join(root, ".agents", "skills"))
	names := make([]string, 0, len(skillsDir.Entries))
	for _, child := range skillsDir.Entries {
		if child.IsDir() {
			names = append(names, child.Name())
		}
	}
	sort.Strings(names)
	var entries []Entry
	for _, name := range names {
		file := filepath.Join(root, ".agents", "skills", name, "SKILL.md")
		c := bounds.ClassifyNoFollow(file)
		// Classification runs before the command-adapter question so suppression cannot
		// hide a defective producer: suppression drops a row from the rendered block, it
		// does not excuse a directory named like an adapter from being diagnosed when
		// nothing wrote its SKILL.md or the bytes there cannot be vouched for.
		if c.State == bounds.StateAbsent {
			entries = append(entries, Entry{Name: name, Orphan: true})
			continue
		}
		if c.State.Failed() {
			entries = append(entries, Entry{Name: name, Refusal: c.Reason})
			continue
		}
		if exists(filepath.Join(root, ".agents", "commands", name+".md")) {
			continue
		}
		entry := Entry{
			Name:    name,
			Trigger: FrontmatterField(file, "index"),
			Note:    FrontmatterField(file, "index-note"),
			KitOnly: kitOnly[".agents/skills/"+name],
		}
		// A field only becomes a rendered line once it is safe for one; the fence parser
		// publishes the value, this decides whether it may reach the sink.
		if reason := controlRefusal("index", entry.Trigger); reason != "" {
			entry.Refusal = reason
		} else if reason := controlRefusal("index-note", entry.Note); reason != "" {
			entry.Refusal = reason
		}
		entries = append(entries, entry)
	}
	return entries, err
}

// FrontmatterField returns the first value of key inside path's complete leading
// `---` fence. The fence must open at byte zero and close: prose before the opener
// and an unclosed opener are body text, so neither can authorize an indexed field.
// Nothing after the closing fence counts either. This is the `<key>:` prefix contract,
// not YAML — nesting, quoting, and schema validation stay out of scope.
func FrontmatterField(path, key string) string {
	// The producer classifier, not os.ReadFile: this path is attacker-shaped input to
	// generated output, so a link is refused rather than followed and a FIFO cannot
	// block the gate in open(2).
	classified := bounds.ClassifyNoFollow(path)
	if classified.State != bounds.StateParsed {
		return ""
	}
	lines := strings.Split(string(classified.Data), "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return ""
	}
	prefix := key + ":"
	value, found := "", false
	for _, line := range lines[1:] {
		if line == "---" {
			return value
		}
		if !found && strings.HasPrefix(line, prefix) {
			value, found = strings.TrimSpace(strings.TrimPrefix(line, prefix)), true
		}
	}
	// Falling out of the loop means the opener never closed, so the body it swallowed
	// is not frontmatter and the value it appeared to carry is withheld.
	return ""
}

// Check grades root's committed block against the generated one and returns the
// attributed diagnostics in skill-alphabetical order, with the unattributed block
// line last. The order is part of the contract: the gate and `bench skills-index`
// print the same sequence.
func Check(root string) []string {
	entries, allowlistErr := Entries(root)
	// A skill can earn more than one diagnostic — no trigger and a stale committed
	// entry is the colliding case — so this accumulates per skill rather than
	// keying one string per name, which would silently drop the first of the pair.
	attributed := map[string][]string{}
	var expected []string
	expectedByName := map[string]string{}
	for _, entry := range entries {
		if entry.Orphan {
			attributed[entry.Name] = append(attributed[entry.Name], orphanDiagnostic(entry.Name))
			continue
		}
		if entry.Refusal != "" {
			attributed[entry.Name] = append(attributed[entry.Name], refusalDiagnostic(entry.Name, entry.Refusal))
			continue
		}
		if !entry.Indexed() {
			attributed[entry.Name] = append(attributed[entry.Name], fmt.Sprintf("skill '%s' missing index: frontmatter (the skills index is generated)", entry.Name))
			continue
		}
		line := entry.Line()
		expected = append(expected, line)
		expectedByName[entry.Name] = line
	}

	if allowlistErr != nil {
		// Every kit-only marker in the generated block would be missing, so the block
		// comparison below would report drift the skills did not cause.
		return append(ordered(attributed), allowlistDiagnostic)
	}

	ref, refusal := readReference(root)
	if refusal != "" {
		return append(ordered(attributed), refusal)
	}
	span, ok := findBlock(ref)
	if !ok {
		return append(ordered(attributed), markersDiagnostic)
	}
	if strings.Join(expected, "\n") == strings.Join(span.block, "\n") {
		return ordered(attributed)
	}

	actualByName := committedByName(span.block)
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
	ref, refusal := readReference(root)
	if refusal != "" {
		return errors.New(refusal)
	}
	var block []string
	for _, entry := range entries {
		// A producer defect means the generated block would be authored from a tree the
		// reader cannot vouch for, so the reference keeps its bytes until the source is
		// fixed — the same posture Check takes, and the same wording.
		if entry.Orphan {
			return errors.New(orphanDiagnostic(entry.Name))
		}
		if entry.Refusal != "" {
			return errors.New(refusalDiagnostic(entry.Name, entry.Refusal))
		}
		if entry.Indexed() {
			block = append(block, entry.Line())
		}
	}
	span, ok := findBlock(ref)
	if !ok {
		return errors.New(markersDiagnostic)
	}
	rewritten := span.replace(block)
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

// orphanDiagnostic names the one skill-directory-without-SKILL.md wording Check
// reports and Write refuses with.
func orphanDiagnostic(name string) string {
	return fmt.Sprintf(skillPathFormat+" missing (orphan skill directory)", name)
}

// controlRefusal reasons about a parsed field's fitness for the sink. Entry.Line is
// line-structured markdown with no escaping, so any control rune could truncate the
// entry or forge a second one — a CR is enough. The whole Unicode control category is
// refused, C0 and C1 alike, because an ASCII shortlist leaves DEL and the C1 range
// able to reach the same sink.
func controlRefusal(key, value string) string {
	for _, r := range value {
		if unicode.IsControl(r) {
			return fmt.Sprintf("%s: carries control rune U+%04X (an index entry is one line)", key, r)
		}
	}
	return ""
}

// refusalDiagnostic names the one untrustworthy-bytes wording Check reports and Write
// refuses with.
func refusalDiagnostic(name, reason string) string {
	return refusedDiagnostic(fmt.Sprintf(skillPathFormat, name), reason)
}

// refusedDiagnostic is the one shape every untrustworthy-producer refusal takes, so a
// skill file and the reference read the same way to an operator and a consumer can
// recognize either without matching two wordings.
func refusedDiagnostic(path, reason string) string {
	return path + " refused: " + reason
}

// ReferenceRefusalPrefix is what a refused reference diagnostic starts with. It is
// exported because the conformance composition row has to recognize this package's
// refusal through a registered check without restating the wording.
func ReferenceRefusalPrefix() string {
	return refusedDiagnostic(referenceRel, "")
}

// kitOnlySources reads the allowlist's withheld skill sources through the canonical
// payload reader, so the index marks exactly the rows every other destination honours
// instead of naming the kit-only skills — or re-deriving row validity — a second time.
// Only absence is optional: a tree with no allowlist withholds nothing, which is also
// what the generator concludes. Anything present that does not resolve, empty bytes
// included, leaves marking unresolved.
func kitOnlySources(root string) (map[string]bool, error) {
	rows, absent, err := kitpayload.PayloadRowsAt(filepath.Join(root, filepath.FromSlash(allowlistRel)))
	if absent {
		return nil, nil
	}
	if err != nil {
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
// trigger before any block drift.
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

// committedByName keys the committed block's lines by the skill each one names.
// A line naming no skill is unattributable and carries no key.
func committedByName(lines []string) map[string]string {
	out := map[string]string{}
	for _, line := range lines {
		if name := skillName(line); name != "" {
			out[name] = line
		}
	}
	return out
}

// blockSpan is the reference file cut at the generated block: the lines through the
// start marker, the block itself, and the end marker onward. Check reads the block and
// Write substitutes it, so both see the same span of the same file.
type blockSpan struct {
	before []string
	block  []string
	after  []string
}

// findBlock locates the generated block, scanning the whole file rather than stopping
// at the first plausible pair: only one start marker followed by one end marker names a
// block unambiguously. Zero, reversed, unclosed, or repeated markers leave the writer
// with a choice it has no basis to make, so the file carries no block and keeps its bytes.
func findBlock(text string) (blockSpan, bool) {
	lines := strings.Split(text, "\n")
	start, end := -1, -1
	for i, line := range lines {
		switch line {
		case startMarker:
			if start >= 0 {
				return blockSpan{}, false
			}
			start = i
		case endMarker:
			if end >= 0 {
				return blockSpan{}, false
			}
			end = i
		}
	}
	if start < 0 || end < start {
		return blockSpan{}, false
	}
	return blockSpan{before: lines[:start+1], block: lines[start+1 : end], after: lines[end:]}, true
}

func (s blockSpan) replace(block []string) string {
	out := append([]string{}, s.before...)
	out = append(out, block...)
	out = append(out, s.after...)
	return strings.Join(out, "\n")
}

// readReference classifies the reference once for both Check and Write, returning its
// text or the one diagnostic that ends the read. The states stay apart because they
// send an operator to different repairs: absent means write the file, present-and-empty
// means the file lost its content, and anything else means the path is not the producer
// it claims to be — collapsing them to absence is exactly how a false empty source
// would authorize replacing tracked bytes.
func readReference(root string) (string, string) {
	classified := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(referenceRel)))
	switch classified.State {
	case bounds.StateParsed:
		return string(classified.Data), ""
	case bounds.StateAbsent:
		return "", referenceRel + " missing (skills index unverifiable)"
	case bounds.StateEmpty:
		return "", referenceRel + " empty (skills index unverifiable)"
	}
	return "", refusedDiagnostic(referenceRel, classified.Reason)
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}
