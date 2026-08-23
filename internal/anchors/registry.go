package anchors

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// Anchor describes one ordered conformance prose anchor.
type Anchor struct {
	Group      Group
	File       string
	Kind       Kind
	Section    string
	Needle     string
	Diagnostic string
}

// Group places registry evaluations around bespoke conformance checks.
type Group uint8

const (
	BeforeStructured Group = iota
	AfterStructured
	AfterRoadmapContext
	AfterImplementSpec
	AfterSpecAuthorization
)

// Entries returns the ordered anchor registry.
func Entries() []Anchor {
	return append([]Anchor(nil), registry...)
}

// EvaluateGroup checks one ordered registry group against root.
func EvaluateGroup(root string, group Group) []string {
	var diagnostics []string
	files := map[string]fileResult{}
	sections := map[string]sectionResult{}
	// The evaluator reports one refusal per file, not one per anchor. A refused file
	// fails every anchor it owns, so the report does not repeat the same repair many
	// times.
	reported := map[string]bool{}
	for _, anchor := range registry {
		if anchor.Group != group {
			continue
		}
		file, loaded := files[anchor.File]
		if !loaded {
			file = read(filepath.Join(root, filepath.FromSlash(anchor.File)), anchor.File)
			files[anchor.File] = file
		}
		if file.refusal != "" {
			if !reported[anchor.File] {
				reported[anchor.File] = true
				diagnostics = append(diagnostics, file.refusal)
			}
			continue
		}
		if anchor.Kind == Require {
			if !file.exists {
				diagnostics = append(diagnostics, "acceptance coverage anchor file missing: "+anchor.File)
				continue
			}
			if !Satisfied(anchor.Kind, file.active, anchor.Needle) {
				diagnostics = append(diagnostics, anchor.Diagnostic)
			}
			continue
		}
		if anchor.Kind == Forbid {
			if !Satisfied(anchor.Kind, file.active, anchor.Needle) {
				diagnostics = append(diagnostics, anchor.Diagnostic)
			}
			continue
		}
		key := anchor.File + "\x00" + anchor.Section
		section, resolved := sections[key]
		if !resolved {
			section = resolveSection(anchor.File, anchor.Section, file.active, file.exists)
			sections[key] = section
			if section.diagnostic != "" {
				diagnostics = append(diagnostics, section.diagnostic)
			}
		}
		if section.diagnostic == "" && !Satisfied(anchor.Kind, section.body, anchor.Needle) {
			diagnostics = append(diagnostics, anchor.Diagnostic)
		}
	}
	return diagnostics
}

type fileResult struct {
	active string
	exists bool
	// refusal is set when the path exists but its bytes are untrustworthy. This field
	// stays separate from exists. A missing anchor file tells the reader to write one.
	// A link or a special file at that path needs a different repair.
	refusal string
}

type sectionResult struct {
	body       string
	diagnostic string
}

func resolveSection(file, title, active string, exists bool) sectionResult {
	if !exists {
		return sectionResult{diagnostic: "section-scoped anchor file missing: " + file}
	}
	body, count := MarkdownH2Sections(active, title)
	if count == 0 {
		return sectionResult{diagnostic: fmt.Sprintf("%s is missing the %q section that owns a scoped anchor", file, title)}
	}
	if count > 1 {
		return sectionResult{diagnostic: fmt.Sprintf("%s carries %d %q sections; a scoped anchor needs exactly one owning section", file, count, title)}
	}
	return sectionResult{body: body}
}

// RefusalPrefix opens every refused-anchor-file diagnostic. A consumer that composes
// this registry with its own checks can use this prefix to recognize the registry's
// refusal. The consumer does not need to restate the wording.
const RefusalPrefix = "acceptance coverage anchor file refused: "

// read classifies an anchor file before it opens one. Registry targets include skill
// and reference producer files. A link is refused rather than followed, so a FIFO
// cannot block the gate in open(2).
func read(path, rel string) fileResult {
	classified := bounds.ClassifyNoFollow(path)
	switch {
	case classified.State == bounds.StateAbsent:
		return fileResult{}
	case classified.State.Failed():
		return fileResult{exists: true, refusal: fmt.Sprintf("%s%s (%s)", RefusalPrefix, rel, classified.Reason)}
	}
	return fileResult{active: StripHTMLComments(string(classified.Data)), exists: true}
}

// StripHTMLComments removes complete comments and truncates at an unterminated comment.
func StripHTMLComments(text string) string {
	for {
		start := strings.Index(text, "<!--")
		if start < 0 {
			return text
		}
		end := strings.Index(text[start+4:], "-->")
		if end < 0 {
			return text[:start]
		}
		text = text[:start] + text[start+4+end+3:]
	}
}
