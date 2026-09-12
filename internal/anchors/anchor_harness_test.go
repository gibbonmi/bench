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
				if r.file == file && r.section == section && h.present(i, r, broken) {
					body += h.line(r)
				}
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
