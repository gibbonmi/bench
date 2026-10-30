package adopt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/markdown"
)

// TestMarkersSkipTildeFences pins the managed-block markers to the block reader. Each case
// names the coverage row that it holds. A failure names the case and never prints the input,
// because the input holds fence markers and comment openers.
func TestMarkersSkipTildeFences(t *testing.T) {
	const fenceErr = "conflict: AGENTS.md has an unclosed code fence around Bench markers; marker detection cannot be trusted"
	block := BenchAgentsBlock()
	lines := func(l ...string) string { return strings.Join(l, "\n") }
	start, end, bt := benchStartMarker, benchEndMarker, markdown.FenceMarker
	tildeExample := lines("# Project", "~~~", start, "example", end, "~~~", "")
	liveCRLF := strings.ReplaceAll(lines("before", start, "old", end, "after", ""), "\n", "\r\n")
	cases := []struct {
		name, in, want, wantErr string
		strip                   bool
	}{
		{name: "MB57 rewrite keeps a tilde-fenced example",
			in:   tildeExample + lines(start, "old", end, "after", ""),
			want: tildeExample + block + "after\n"},
		{name: "MB58 strip keeps a tilde-fenced example", strip: true,
			in:   tildeExample + lines(start, "old", end, "after", ""),
			want: tildeExample + "after\n"},
		{name: "MB61 an unterminated tilde fence refuses",
			in:      lines("Broken docs:", "~~~", start, end, "KEEP", ""),
			wantErr: fenceErr},
		{name: "MB70 an unterminated comment hides the fence closer",
			in:      lines(bt, start, "example", end, "<!-- draft", bt, ""),
			wantErr: fenceErr},
		{name: "MB70 an unterminated comment follows a hidden fence opener",
			in:      lines("<!-- note", bt, start, "example", end, bt, "<!-- draft", ""),
			wantErr: fenceErr},
		{name: "MB86 an unterminated comment with no Bench text rewrites",
			in:   lines("Project", "<!-- draft", ""),
			want: lines("Project", "<!-- draft", "", block)},
		{name: "MB87 rewrite keeps each carriage return outside the block",
			in:   liveCRLF,
			want: "before\r\n" + block + "after\r\n"},
		{name: "MB91 strip keeps each carriage return outside the block", strip: true,
			in:   liveCRLF,
			want: "before\r\nafter\r\n"},
		{name: "MB97 an unclosed frontmatter marker reads as body",
			in:   lines("---", "intro", start, "old", end, ""),
			want: lines("---", "intro", block)},
		{name: "rewrite of a block that ends the file without a newline",
			in:   lines("before", start, "old", end),
			want: "before\n" + block},
		{name: "rewrite ends text after the block with a newline",
			in:   lines(start, "old", end, "after"),
			want: block + "after\n"},
		{name: "strip of a block that ends the file without a newline", strip: true,
			in:   lines("before", start, "old", end),
			want: "before\n"},
		{name: "strip of only the block is empty", strip: true,
			in:   block,
			want: ""},
	}
	for _, c := range cases {
		rebuild := RewriteAgentsBlock
		if c.strip {
			rebuild = StripAgentsBlock
		}
		got, err := rebuild(c.in)
		switch {
		case c.wantErr != "" && (err == nil || err.Error() != c.wantErr):
			t.Errorf("%s: err = %v, want the unclosed-fence refusal", c.name, err)
		case c.wantErr == "" && err != nil:
			t.Errorf("%s: err = %v", c.name, err)
		case c.wantErr == "" && got != c.want:
			t.Errorf("%s: got %d bytes that differ from the %d bytes wanted", c.name, len(got), len(c.want))
		}
	}

	linked, err := RewriteAgentsBlock("Project\n")
	if err != nil {
		t.Fatalf("link rewrite: %v", err)
	}
	if got, err := StripAgentsBlock(linked); err != nil || got != "Project\n" {
		t.Errorf("strip of a linked file: err = %v, the original bytes did not come back", err)
	}
}

// TestLinkPlanShipsClaudeAgents pins the one path from the payload's agents tree row to a
// linked repo. The link plan is the only such path, so a missing row or a missing
// destination case ships a consumer no agent types and returns every delegate to the
// full tool set. The fixture kit carries the two agent files and nothing else, because
// every other payload row is absent-tolerant.
func TestLinkPlanShipsClaudeAgents(t *testing.T) {
	kit := t.TempDir()
	agents := filepath.Join(kit, ".claude", "agents")
	if err := os.MkdirAll(agents, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"bench-reviewer.md", "bench-writer.md"} {
		writeFixtureFile(t, filepath.Join(agents, name), "---\nname: x\n---\n", 0o644)
	}

	plan, err := buildLinkPlan(kit)
	if err != nil {
		t.Fatalf("buildLinkPlan = %v", err)
	}
	for _, name := range []string{"bench-reviewer.md", "bench-writer.md"} {
		want := ".claude/agents/" + name
		found := false
		for _, entry := range plan {
			if entry.rel == want && entry.kind == "file" {
				found = true
			}
		}
		if !found {
			t.Errorf("link plan has no file entry at %s", want)
		}
	}
}
