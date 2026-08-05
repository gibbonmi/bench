package anchors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	for _, anchor := range registry {
		if anchor.Group != group {
			continue
		}
		file, loaded := files[anchor.File]
		if !loaded {
			text, exists := read(filepath.Join(root, filepath.FromSlash(anchor.File)))
			file = fileResult{active: StripHTMLComments(text), exists: exists}
			files[anchor.File] = file
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

func read(path string) (string, bool) {
	data, err := os.ReadFile(path)
	return string(data), err == nil
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
