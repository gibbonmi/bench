package spec

import (
	"errors"
	"testing"
)

// selectionCase is one selection of authored bytes. want is the independently authored
// selected text; absent wants an absent selection, and refused wants a selection refusal.
type selectionCase struct {
	content, selector, want string
	absent, refused         bool
}

func gradeSelection(t *testing.T, name string, c selectionCase, selection Selection, err error) {
	t.Helper()
	switch {
	case c.refused:
		if !errors.Is(err, ErrSelection) {
			t.Fatalf("%s: err = %v, want a selection refusal", name, err)
		}
	case err != nil:
		t.Fatalf("%s: err = %v", name, err)
	case c.absent:
		if selection.Present || selection.Start != selection.End {
			t.Fatalf("%s: selection = %+v, want an absent selection", name, selection)
		}
	case !selection.Present || string(selection.Bytes([]byte(c.content))) != c.want:
		t.Fatalf("%s: selection = %+v selecting %q, want present %q", name, selection, selection.Bytes([]byte(c.content)), c.want)
	}
}

// TestExactMarkdownSelections grades the literal visible H2 section: the owning heading
// through the byte before the next visible H2, a heading in a fence or an HTML comment that
// owns nothing, and the refusals of a repeated heading and a Markdown fault.
func TestExactMarkdownSelections(t *testing.T) {
	const doc = "# Title\n\n## Alpha\nalpha body\n```text\n## Fenced\n```\n<!--\n## Commented\n-->\n<!-- ## Inline -->\n## Beta\nbeta\n## Empty\n## Gamma\r\ngamma é"
	for name, c := range map[string]selectionCase{
		"section with hidden headings": {content: doc, selector: "Alpha", want: "## Alpha\nalpha body\n```text\n## Fenced\n```\n<!--\n## Commented\n-->\n<!-- ## Inline -->\n"},
		"middle section":               {content: doc, selector: "Beta", want: "## Beta\nbeta\n"},
		"empty body":                   {content: doc, selector: "Empty", want: "## Empty\n"},
		"last section to end of file":  {content: doc, selector: "Gamma", want: "## Gamma\r\ngamma é"},
		"fenced heading":               {content: doc, selector: "Fenced", absent: true},
		"commented heading":            {content: doc, selector: "Commented", absent: true},
		"inline commented heading":     {content: doc, selector: "Inline", absent: true},
		"H1 title":                     {content: doc, selector: "Title", absent: true},
		"absent heading":               {content: doc, selector: "Missing", absent: true},
		"repeated heading":             {content: "## Dup\na\n## Other\n## Dup\nb\n", selector: "Dup", refused: true},
		"unterminated fence":           {content: "## Alpha\n```\n## Beta\n", selector: "Alpha", refused: true},
		"unterminated comment":         {content: "## Alpha\n<!--\n## Beta\n", selector: "Beta", refused: true},
	} {
		selection, err := Section([]byte(c.content), c.selector)
		gradeSelection(t, name, c, selection, err)
	}
}

// TestExactSpecSelections grades the show shortcuts: the whole file, the actual User stories
// section, and the raw Status: line that the metadata owner selects, where no status differs
// from an empty status, and the refusal of the coverage shortcut that another owner bounds.
func TestExactSpecSelections(t *testing.T) {
	const live = "# Spec\n\nStatus: draft\nStatus: staged\r\n```text\nStatus: fenced\n```\n\n## User stories\n\nLine: opus / medium.\n\n### Group\n\n1. As an agent, I want one story.\n\n## Implementation decisions\n\nnone\n"
	for name, c := range map[string]selectionCase{
		"whole file":   {content: live, want: live},
		"empty file":   {content: "", want: ""},
		"stories":      {content: live, selector: StoriesShortcut, want: "## User stories\n\nLine: opus / medium.\n\n### Group\n\n1. As an agent, I want one story.\n\n"},
		"literal":      {content: live, selector: "Implementation decisions", want: "## Implementation decisions\n\nnone\n"},
		"no stories":   {content: "# Spec\n", selector: StoriesShortcut, absent: true},
		"status":       {content: live, selector: StatusShortcut, want: "Status: staged\r\n"},
		"empty status": {content: "# Spec\nStatus:", selector: StatusShortcut, want: "Status:"},
		"no status":    {content: "# Spec\n```\nStatus: fenced\n```\n", selector: StatusShortcut, absent: true},
		// The coverage owner bounds the coverage shortcut, so a literal Coverage title never
		// answers it here.
		"coverage": {content: "# Spec\n\n## coverage\n\nbody\n", selector: CoverageShortcut, refused: true},
	} {
		selection, err := Show([]byte(c.content), c.selector)
		gradeSelection(t, name, c, selection, err)
	}
	if status, _ := Metadata([]byte(live)); status != "staged" {
		t.Fatalf("metadata status = %q, want the staged value of the selected line", status)
	}
}
