package anchors

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gibbonmi/bench/internal/bounds"
)

// Anchor describes one ordered conformance prose anchor.
type Anchor struct {
	Group   Group
	File    string
	Kind    Kind
	Section string
	// Step names the numbered step of Section that owns the needle. A zero Step leaves
	// the anchor unscoped by step.
	Step       int
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

// Location pairs a registered anchor with its first physical match line.
type Location struct {
	Anchor
	Line int
}

// PathEvaluation reports one classified path's locations and registry diagnostics.
type PathEvaluation struct {
	Locations   []Location
	Diagnostics []string
	State       bounds.FileState
	Reason      string
}

// EvaluatePath evaluates every registered anchor for one repository-relative path.
func EvaluatePath(root, path string) PathEvaluation {
	return evaluate(root, nil, path)
}

// EvaluateGroup checks one ordered registry group against root.
func EvaluateGroup(root string, group Group) []string {
	return evaluateGroup(root, group, "")
}

func evaluateGroup(root string, group Group, subject string) []string {
	return evaluate(root, &group, subject).Diagnostics
}

func evaluate(root string, group *Group, subject string) PathEvaluation {
	var result PathEvaluation
	files := map[string]fileResult{}
	if subject != "" {
		file := read(filepath.Join(root, filepath.FromSlash(subject)), subject)
		files[subject] = file
		result.State, result.Reason = file.classified.State, file.classified.Reason
	}
	sections := map[string]scopeResult{}
	steps := map[string]scopeResult{}
	// The evaluator reports one refusal per file, not one per anchor. A refused file
	// fails every anchor it owns, so the report does not repeat the same repair many
	// times.
	reported := map[string]bool{}
	for _, anchor := range registry {
		if group != nil && anchor.Group != *group || subject != "" && anchor.File != subject {
			continue
		}
		file, loaded := files[anchor.File]
		if !loaded {
			file = read(filepath.Join(root, filepath.FromSlash(anchor.File)), anchor.File)
			files[anchor.File] = file
		}
		result.Locations = append(result.Locations, Location{Anchor: anchor, Line: locate(anchor.Kind, anchor.Section, anchor.Step, anchor.Needle, string(file.classified.Data))})
		if file.refusal != "" {
			if !reported[anchor.File] {
				reported[anchor.File] = true
				result.Diagnostics = append(result.Diagnostics, file.refusal)
			}
			continue
		}
		if anchor.Kind == Require {
			if !file.exists {
				result.Diagnostics = append(result.Diagnostics, "acceptance coverage anchor file missing: "+anchor.File)
				continue
			}
			if !Satisfied(anchor.Kind, file.active, anchor.Needle) {
				result.Diagnostics = append(result.Diagnostics, anchor.Diagnostic)
			}
			continue
		}
		if anchor.Kind == Forbid || anchor.Kind == ForbidCaseFoldedEmphasis {
			if !Satisfied(anchor.Kind, file.active, anchor.Needle) {
				result.Diagnostics = append(result.Diagnostics, anchor.Diagnostic)
			}
			continue
		}
		key := anchor.File + "\x00" + anchor.Section
		section, resolved := sections[key]
		if !resolved {
			section = resolveSection(anchor.File, anchor.Section, file.active, file.exists)
			sections[key] = section
			if section.diagnostic != "" {
				result.Diagnostics = append(result.Diagnostics, section.diagnostic)
			}
		}
		if section.diagnostic != "" {
			continue
		}
		body := section.body
		// The kind decides the step narrowing here, exactly as it does in the locator, so
		// the two cannot disagree about which anchors read a step. A Step on a kind that
		// reads none is an authoring mistake that TestRegistryBindsStepToItsKind refuses.
		if anchor.Kind.stepScoped() {
			stepKey := key + "\x00" + strconv.Itoa(anchor.Step)
			step, stepResolved := steps[stepKey]
			if !stepResolved {
				step = resolveStep(anchor.File, anchor.Section, anchor.Step, section.body)
				steps[stepKey] = step
				if step.diagnostic != "" {
					result.Diagnostics = append(result.Diagnostics, step.diagnostic)
				}
			}
			if step.diagnostic != "" {
				continue
			}
			body = step.body
		}
		if !Satisfied(anchor.Kind, body, anchor.Needle) {
			result.Diagnostics = append(result.Diagnostics, anchor.Diagnostic)
		}
	}
	return result
}

type fileResult struct {
	classified bounds.Classified
	active     string
	exists     bool
	// refusal is set when the path exists but its bytes are untrustworthy. This field
	// stays separate from exists. A missing anchor file tells the reader to write one.
	// A link or a special file at that path needs a different repair.
	refusal string
}

// scopeResult is one resolved narrowing of an anchor's subject: the body the evaluator
// searches, or the diagnostic that says why no single body exists.
type scopeResult struct {
	body       string
	diagnostic string
}

func resolveSection(file, title, active string, exists bool) scopeResult {
	if !exists {
		return scopeResult{diagnostic: "section-scoped anchor file missing: " + file}
	}
	body, count := MarkdownH2Sections(active, title)
	if count == 0 {
		return scopeResult{diagnostic: fmt.Sprintf("%s is missing the %q section that owns a scoped anchor", file, title)}
	}
	if count > 1 {
		return scopeResult{diagnostic: fmt.Sprintf("%s carries %d %q sections; a scoped anchor needs exactly one owning section", file, count, title)}
	}
	return scopeResult{body: body}
}

// resolveStep narrows a resolved section body to one numbered step. Its refusals mirror the
// section's: an unnamed step, no owning step, and more than one line that opens the same step.
func resolveStep(file, title string, step int, section string) scopeResult {
	if step == 0 {
		return scopeResult{diagnostic: fmt.Sprintf("%s carries a step-scoped anchor with no step in the %q section; a step-scoped anchor names one step", file, title)}
	}
	body, count := MarkdownNumberedSteps(section, step)
	if count == 0 {
		return scopeResult{diagnostic: fmt.Sprintf("%s is missing step %d of the %q section that owns a step-scoped anchor", file, step, title)}
	}
	if count > 1 {
		return scopeResult{diagnostic: fmt.Sprintf("%s carries %d lines that open step %d of the %q section; a step-scoped anchor needs exactly one owning step", file, count, step, title)}
	}
	return scopeResult{body: body}
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
	file := fileResult{classified: classified, exists: classified.State != bounds.StateAbsent}
	switch {
	case classified.State == bounds.StateAbsent:
		return file
	case classified.State.Failed():
		file.refusal = fmt.Sprintf("%s%s (%s)", RefusalPrefix, rel, classified.Reason)
	default:
		file.active = StripHTMLComments(string(classified.Data))
	}
	return file
}

// StripHTMLComments removes complete comments and truncates at an unterminated comment.
// It projects stripCommentsMapped, the package's one comment strip, so the evaluator reads
// the text Locate searches.
func StripHTMLComments(text string) string {
	stripped, _ := stripCommentsMapped(text)
	return string(stripped)
}
