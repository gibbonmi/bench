package spec

import (
	"strings"
	"testing"
)

func TestSpecWidenReportsACoveredEntry(t *testing.T) {
	root := widenRepo(t, "internal/x/y.go, internal/status", "internal/x/y.go", "internal/status/s.go")
	before := readTicket(t, root)
	out, code := runWiden(t, root, "internal/x/y.go", "internal/status/t.go")
	requireWidened(t, root, out, code, widenTable(t, []string{widenTicketPath, "internal/x/y.go", "current"}, []string{widenTicketPath, "internal/status/t.go", "current"}), before)
	out, code = runWiden(t, root, "internal/x/z.go", "internal/x/z.go")
	want := widenTable(t, []string{widenTicketPath, "internal/x/z.go", "added"}, []string{widenTicketPath, "internal/x/z.go", "current"})
	requireWidened(t, root, out, code, want, widenTicket("internal/x/y.go, internal/status, internal/x/z.go"))

	// A marker or a trailing slash on either side names the same path, so it adds no entry.
	for _, tc := range []struct {
		name, writes string
		entries      []string
	}{
		{"a marked or slashed ticket entry", "internal/x/y.go (new), internal/status/, internal/d/old.go (deleted)", []string{"internal/x/y.go", "internal/status/t.go", "internal/d/old.go"}},
		{"a marked or slashed request entry", "internal/x/y.go, internal/status, internal/d/old.go", []string{"internal/x/y.go (new)", "internal/status/", "internal/d/old.go (deleted)"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := widenRepo(t, tc.writes, "internal/x/y.go", "internal/status/s.go", "internal/d/old.go")
			before := readTicket(t, root)
			rows := make([][]string, 0, len(tc.entries))
			for _, entry := range tc.entries {
				rows = append(rows, []string{widenTicketPath, entry, "current"})
			}
			out, code := runWiden(t, root, tc.entries...)
			requireWidened(t, root, out, code, widenTable(t, rows...), before)
		})
	}
	root = widenRepo(t, "internal/x/y.go", "internal/x/y.go")
	out, code = runWiden(t, root, "internal/x/sub/ (new)", "internal/x/sub/z.go")
	want = widenTable(t, []string{widenTicketPath, "internal/x/sub/ (new)", "added"}, []string{widenTicketPath, "internal/x/sub/z.go", "current"})
	requireWidened(t, root, out, code, want, widenTicket("internal/x/y.go, internal/x/sub/ (new)"))
}

func TestSpecWidenKeepsOtherBytes(t *testing.T) {
	for _, tc := range []struct {
		name, before, after string
	}{
		{
			name:   "a fenced Writes: line and the text around the field stay",
			before: "# A\n\n```\nWrites: decoy\n```\nBlocked by: none\nWrites: internal/x/y.go  \nCovers: CR1\n\n## What to build\n\nx\n\n## Acceptance\n\n- [ ] x\n",
			after:  "# A\n\n```\nWrites: decoy\n```\nBlocked by: none\nWrites: internal/x/y.go, internal/x/z.go  \nCovers: CR1\n\n## What to build\n\nx\n\n## Acceptance\n\n- [ ] x\n",
		},
		{
			name:   "CRLF line endings stay",
			before: strings.ReplaceAll(widenTicket("internal/x/y.go"), "\n", "\r\n"),
			after:  strings.ReplaceAll(widenTicket("internal/x/y.go, internal/x/z.go"), "\n", "\r\n"),
		},
		{
			name:   "a last line without a newline stays without one",
			before: "# A\n\nBlocked by: none\nCovers: CR1\n\n## What to build\n\nx\n\n## Acceptance\n\n- [ ] x\nWrites: internal/x/y.go",
			after:  "# A\n\nBlocked by: none\nCovers: CR1\n\n## What to build\n\nx\n\n## Acceptance\n\n- [ ] x\nWrites: internal/x/y.go, internal/x/z.go",
		},
		{
			name:   "the none sentinel gives way to the first entry",
			before: widenTicket("none"),
			after:  widenTicket("internal/x/z.go"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// b.md signs off the scope, so a.md can hold any Writes: value.
			files := map[string]string{widenTicketPath: tc.before, "specs/s/tickets/b.md": widenTicket("internal/x/y.go")}
			root := retireWorktree(t, widenPrimary(t, files, "internal/x/y.go"))
			out, code := runWiden(t, root, "internal/x/z.go")
			requireWidened(t, root, out, code, "", tc.after)
		})
	}
}
