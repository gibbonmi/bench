package adopt

import (
	"os"
	"path/filepath"
	"testing"
)

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
