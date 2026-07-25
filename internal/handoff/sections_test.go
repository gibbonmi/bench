package handoff

import (
	"strings"
	"testing"
)

// TestSplitAmbiguousState walks the splitter's edge set. The refusals are the point: a
// first-match splitter passes every well-formed body below and silently discards prose on
// the ambiguous ones, which is the data-loss path the command must never take.
func TestSplitAmbiguousState(t *testing.T) {
	refusals := []struct {
		name, content string
	}{
		{"no heading at all", "# Session handoff\n\n## Next command\n\n`bench status`\n"},
		{"two unfenced headings", "## State\n\nfirst\n\n## State\n\nsecond\n"},
		{"the only heading sits inside a fence", "# Session handoff\n\n```\n## State\n\nnot a section\n```\n"},
		{"one unfenced heading and one fenced", "## State\n\nreal\n\n## Notes\n\n~~~\n## State\n~~~\n"},
	}
	for _, tc := range refusals {
		t.Run(tc.name, func(t *testing.T) {
			body, err := splitState([]byte(tc.content))
			if err == nil {
				t.Fatalf("split accepted an ambiguous document, body = %q", body)
			}
			if !strings.Contains(err.Error(), stateHeading) {
				t.Fatalf("refusal does not name the section: %q", err.Error())
			}
		})
	}

	accepted := []struct {
		name, content, want string
	}{
		{"body runs to EOF", "## State\n\nkeep me\n", "keep me"},
		{"body stops at the next generated heading", "## State\n\nkeep me\n\n## Next command\n\n`bench status`\n", "keep me"},
		{"a reviewer's own heading does not end the body", "## State\n\nkeep me\n\n## Notes\n\nkeep me too\n\n## Next command\n\n`bench status`\n", "keep me\n\n## Notes\n\nkeep me too"},
		{"a reviewer's heading survives with Shape as the terminator", "## State\n\nkeep me\n\n## Notes\n\nkeep me too\n\n## Shape\n\nshape text\n", "keep me\n\n## Notes\n\nkeep me too"},
		{"a fenced heading inside the body does not end it", "## State\n\nkeep me\n\n```\n## Next command\n```\n\ntail\n\n## Next command\n\n`x`\n", "keep me\n\n```\n## Next command\n```\n\ntail"},
		{"an empty body is not a refusal", "## State\n\n## Next command\n\n`x`\n", ""},
		{"interior bytes survive verbatim", "## State\n\n  indented\ttab\n\n- list\n", "  indented\ttab\n\n- list"},
	}
	for _, tc := range accepted {
		t.Run(tc.name, func(t *testing.T) {
			body, err := splitState([]byte(tc.content))
			if err != nil {
				t.Fatalf("split refused a well-formed document: %v", err)
			}
			if body != tc.want {
				t.Fatalf("body = %q, want %q", body, tc.want)
			}
		})
	}
}
