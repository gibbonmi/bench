package anchors

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// anchorRule is one test's own expectation: the subject file, the optional H2
// section that owns the needle, the needle itself, and the diagnostic the registry
// must raise when the needle breaks. A forbidden rule inverts the break: the needle
// arrives instead of leaving.
type anchorRule struct {
	file      string
	section   string
	needle    string
	want      string
	forbidden bool
	// fenced wraps the needle in a backtick fence, so a needle that opens with an H2
	// heading does not close the section that must hold it.
	fenced bool
	// step nests the needle under a numbered step opener inside its section. A zero step
	// leaves the needle in the section body, under no step.
	step int
}

// anchorHarness writes one minimal tree per run and evaluates a registry group
// against it. Each test keeps its own rules, needles, and diagnostics; only this
// loop mechanism is shared.
type anchorHarness struct {
	group Group
	rules []anchorRule
	// templates wraps one file's body. The verb %s takes the body. A file with no
	// entry here gets defaultAnchorTemplate.
	templates map[string]string
	// suffix follows each needle on its line, for a needle that names a field.
	suffix string
	// skipCrossTalk drops the cross-talk assertion, for a rule set whose subjects
	// cannot separate one diagnostic from another.
	skipCrossTalk bool
	// extraSteps writes one bare step opener per entry into each section, after the
	// openers the rules themselves need. A repeated entry gives a section two lines that
	// open the same step.
	extraSteps []int
}

// defaultAnchorTemplate gives a subject file a title and nothing else.
const defaultAnchorTemplate = "# subject\n%s"

func (h anchorHarness) evaluate(t *testing.T, broken int) []string {
	t.Helper()
	return EvaluateGroup(h.write(t, broken), h.group)
}

// write gives each run its own minimal tree. Rules sharing a file or section
// contribute to one body so later rules cannot erase earlier needles.
func (h anchorHarness) write(t *testing.T, broken int) string {
	t.Helper()
	root := t.TempDir()
	var files []string
	for _, r := range h.rules {
		if !slices.Contains(files, r.file) {
			files = append(files, r.file)
		}
	}
	for _, file := range files {
		body := ""
		for i, r := range h.rules {
			if r.file == file && r.section == "" && h.present(i, r, broken) {
				body += "\n" + h.line(r)
			}
		}
		var sections []string
		for _, r := range h.rules {
			if r.file == file && r.section != "" && !slices.Contains(sections, r.section) {
				sections = append(sections, r.section)
			}
		}
		for _, section := range sections {
			body += "\n## " + section + "\n\n"
			for i, r := range h.rules {
				if r.file == file && r.section == section && r.step == 0 && h.present(i, r, broken) {
					body += h.line(r)
				}
			}
			var steps []int
			for _, r := range h.rules {
				if r.file == file && r.section == section && r.step != 0 && !slices.Contains(steps, r.step) {
					steps = append(steps, r.step)
				}
			}
			// The opener stands whether or not this run keeps the needle, so an omission
			// reads as a missing needle and never as a missing step.
			for _, step := range steps {
				body += stepOpenerLine(step)
				for i, r := range h.rules {
					if r.file == file && r.section == section && r.step == step && h.present(i, r, broken) {
						body += "   " + h.line(r)
					}
				}
			}
			for _, step := range h.extraSteps {
				body += stepOpenerLine(step)
			}
		}
		template := defaultAnchorTemplate
		if custom, ok := h.templates[file]; ok {
			template = custom
		}
		path := filepath.Join(root, filepath.FromSlash(file))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(fmt.Sprintf(template, body)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// present reports whether the rule's needle belongs in this run's tree.
func (h anchorHarness) present(i int, r anchorRule, broken int) bool {
	return (i == broken) == r.forbidden
}

// stepOpenerLine writes one numbered step opener at column zero, the shape a reader sees
// in a `## Process` section.
func stepOpenerLine(step int) string {
	return fmt.Sprintf("%d. step %d\n", step, step)
}

func (h anchorHarness) line(r anchorRule) string {
	if r.fenced {
		return "```markdown\n" + r.needle + "\n```\n"
	}
	return r.needle + h.suffix + "\n"
}

// check reads both directions by membership: the conformant tree raises no rule's
// diagnostic, and each broken tree raises its own rule's diagnostic and no other.
// Other rows of the group fire against the minimal tree, so membership, not the
// diagnostic count, is the subject.
func (h anchorHarness) check(t *testing.T) {
	t.Helper()
	conformant := h.evaluate(t, -1)
	for _, r := range h.rules {
		if slices.Contains(conformant, r.want) {
			t.Errorf("tree conformant with %q in %s raised %q", r.needle, r.file, r.want)
		}
	}
	for i := range h.rules {
		root := h.write(t, i)
		h.checkBroken(t, root, i, EvaluateGroup(root, h.group))
	}
}

// TestAnchorHarnessStepRules pins the step-scoped kind against its own expectations: a
// needle inside the registered step is conformant, the same needle under another step
// raises the anchor's diagnostic, and a missing or duplicated step opener raises its own
// diagnostic. The four trees differ only in where the step openers sit, so the step
// narrowing, not the section narrowing, is the subject.
func TestAnchorHarnessStepRules(t *testing.T) {
	anchor := pathTestAnchor(t, RequireInStep)
	missing := fmt.Sprintf("%s is missing step %d of the %q section that owns a step-scoped anchor", anchor.File, anchor.Step, anchor.Section)
	duplicate := fmt.Sprintf("%s carries 2 lines that open step %d of the %q section; a step-scoped anchor needs exactly one owning step", anchor.File, anchor.Step, anchor.Section)
	rule := anchorRule{file: anchor.File, section: anchor.Section, needle: anchor.Needle, want: anchor.Diagnostic, step: anchor.Step}
	for _, tc := range []struct {
		name       string
		step       int
		extraSteps []int
		want       string
	}{
		{name: "in the registered step", step: anchor.Step},
		{name: "moved to another step", step: anchor.Step + 1, extraSteps: []int{anchor.Step}, want: anchor.Diagnostic},
		{name: "no such step", step: 0, want: missing},
		{name: "duplicated step", step: anchor.Step, extraSteps: []int{anchor.Step}, want: duplicate},
	} {
		t.Run(tc.name, func(t *testing.T) {
			subject := rule
			subject.step = tc.step
			h := anchorHarness{group: anchor.Group, rules: []anchorRule{subject}, extraSteps: tc.extraSteps}
			diags := EvaluateGroup(h.write(t, -1), anchor.Group)
			for _, other := range []string{anchor.Diagnostic, missing, duplicate} {
				if got := slices.Contains(diags, other); got != (other == tc.want) {
					t.Fatalf("%s tree = %v, want %q present=%t", tc.name, diags, other, other == tc.want)
				}
			}
		})
	}
}

type anchorReporter interface {
	Helper()
	Errorf(string, ...any)
}

func (h anchorHarness) checkBroken(t anchorReporter, root string, broken int, diags []string) {
	t.Helper()
	r := h.rules[broken]
	if !slices.Contains(diags, r.want) {
		t.Errorf("tree broken at %q in %s = %v, want %q", r.needle, r.file, evaluateGroup(root, h.group, r.file), r.want)
	}
	if h.skipCrossTalk {
		return
	}
	for i, other := range h.rules {
		if i != broken && slices.Contains(diags, other.want) {
			t.Errorf("breaking %q in %s also raised %q", r.needle, r.file, other.want)
		}
	}
}
