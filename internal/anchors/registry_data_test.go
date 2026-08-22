package anchors

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// TestFinalCommunicationMarkerTuples keeps the final required tuples and
// retired-marker policy independent of production values, so omissions cannot
// re-derive the expectation green.
func TestFinalCommunicationMarkerTuples(t *testing.T) {
	const finalRoles = "Never assume the reviewer's decisions, and never assume a claim the gate could check instead"
	const retiredRoles = "NEVER assume, always verify"
	if ReviewerDecisionBoundaryMarker != finalRoles {
		t.Errorf("ReviewerDecisionBoundaryMarker = %q; want %q", ReviewerDecisionBoundaryMarker, finalRoles)
	}
	wanted := []Anchor{
		{File: ".bench/BENCH.md", Kind: RequireInSection, Section: "Roles", Needle: finalRoles},
		{File: ".bench/BENCH.md", Kind: RequireInSection, Section: "How to talk to me", Needle: "Clear beats dense"},
		{File: ".bench/BENCH.md", Kind: RequireInSection, Section: "Workflow", Needle: "Right-size the process"},
		{File: ".bench/BENCH.md", Kind: Forbid, Needle: retiredRoles},
	}
	entries := Entries()
	previous := -1
	for _, want := range wanted {
		at := -1
		matches := 0
		for i, entry := range entries {
			if entry.File == want.File && entry.Kind == want.Kind && entry.Section == want.Section && entry.Needle == want.Needle {
				if at < 0 {
					at = i
				}
				matches++
			}
		}
		if matches != 1 {
			t.Errorf("registry has %d rows matching %+v; want exactly one", matches, want)
		}
		if at < 0 {
			continue
		}
		if at <= previous {
			t.Errorf("registry places communication marker row %+v at index %d, out of final policy order", want, at)
		}
		previous = at
	}
}

// TestLandedRetirementAnchorTuples keeps the accepted tuples independent of the
// production registry so a self-consistent omission cannot define itself green.
func TestLandedRetirementAnchorTuples(t *testing.T) {
	wanted := []Anchor{
		{
			Group:      AfterImplementSpec,
			File:       ".agents/skills/bench-craft-delegate/SKILL.md",
			Kind:       RequireInSection,
			Section:    "Verifying the done-claim",
			Needle:     "Acceptance closes an independent worktree after its slice lands: the coordinator runs `bench worktree release --request <opaque-id> <path>` for it.",
			Diagnostic: ".agents/skills/bench-craft-delegate/SKILL.md dropped release-at-acceptance",
		},
		{
			Group:      AfterImplementSpec,
			File:       ".agents/commands/bench-final-check.md",
			Kind:       RequireInSection,
			Section:    "Exit handoff",
			Needle:     "leftover worktrees are retired by `bench worktree clean --landed`: run the plan, apply it, and carry the plan and apply result in the landing report",
			Diagnostic: ".agents/commands/bench-final-check.md post-merge tail dropped the landed worktree sweep step",
		},
		{
			Group:      AfterImplementSpec,
			File:       ".agents/commands/bench-final-check.md",
			Kind:       ForbidInSection,
			Section:    "Exit handoff",
			Needle:     "leftover worktrees and scratch branches go through",
			Diagnostic: ".agents/commands/bench-final-check.md still routes leftover worktrees to a bare per-path clean",
		},
	}
	for _, want := range wanted {
		matches := 0
		for _, entry := range Entries() {
			if entry == want {
				matches++
			}
		}
		if matches != 1 {
			t.Errorf("registry has %d rows matching %+v; want exactly one", matches, want)
		}
	}
}

// TestRetroFeedsMarkerAnchorRedsOnRemoval pins RF28. The needle and the diagnostic are
// written here independently of the registry, so a needle edited to match a template that
// dropped the destination marker cannot define itself green.
func TestRetroFeedsMarkerAnchorRedsOnRemoval(t *testing.T) {
	const (
		clause = "End the item with one line that reads `Feeds: FT<n>`, `Feeds: new`, or `Feeds: none`."
		want   = ".agents/commands/bench-final-check.md dropped the implementation-retro improvement-item destination marker"
	)
	const template = "# Final check\n\n## Capture the implementation retro\n\nWrite each improvement item as one list item. %s\n"

	evaluate := func(t *testing.T, body string) []string {
		t.Helper()
		root := t.TempDir()
		dir := filepath.Join(root, ".agents", "commands")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "bench-final-check.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	// Other rows of the group fire against this minimal tree; only the marker row is
	// this test's subject, so both directions are read by membership.
	if diags := evaluate(t, fmt.Sprintf(template, clause)); slices.Contains(diags, want) {
		t.Fatalf("template carrying the destination-marker clause raised %q", want)
	}
	if diags := evaluate(t, fmt.Sprintf(template, "")); !slices.Contains(diags, want) {
		t.Fatalf("template without the destination-marker clause = %v, want %q", diags, want)
	}
}

// TestDrainFlowRuleAnchorsRedOnRemoval pins RF26 and RF27. Each drain rule is written here
// independently of the registry, next to the diagnostic that must name it, so the anchor
// cannot be re-derived green from a command file that dropped one rule while the rest
// survived.
func TestDrainFlowRuleAnchorsRedOnRemoval(t *testing.T) {
	rules := []struct{ row, rule, want string }{
		{"RF26", "Run `bench roadmap --flow` once and quote its flow block in the exit.", ".agents/commands/bench-drain.md dropped the flow-quote rule: the exit quotes bench roadmap --flow"},
		{"RF26", "An entry feeds a row only when it changes the row's priority, scope, or `Next:`.", ".agents/commands/bench-drain.md dropped the feeds-a-row test: priority, scope, or Next:"},
		{"RF26", "Dismiss an occurrence-only entry with one line of why.", ".agents/commands/bench-drain.md dropped the occurrence-only dismissal rule"},
		{"RF26", "A new row needs a `Next:` token and a class before it opens.", ".agents/commands/bench-drain.md dropped the new-row rule: a Next: token and a class"},
		{"RF27", "When the flow report shows a positive net delta, propose reducing moves in the next batch diff.", ".agents/commands/bench-drain.md dropped the positive-delta restructure rule"},
		{"RF27", "build the item in this session (\"implement now\") by default; open a `ROADMAP.md` row only when the reviewer declines", ".agents/commands/bench-drain.md dropped the build-in-session default for a light-path item"},
	}

	evaluate := func(t *testing.T, body string) []string {
		t.Helper()
		root := t.TempDir()
		dir := filepath.Join(root, ".agents", "commands")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "bench-drain.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}
	command := func(dropped int) string {
		var b strings.Builder
		b.WriteString("# /bench-drain\n\n")
		for i, r := range rules {
			if i == dropped {
				continue
			}
			b.WriteString(r.rule + "\n\n")
		}
		return b.String()
	}

	// Other rows of the group fire against this minimal tree; only the six rule rows are
	// this test's subject, so both directions are read by membership.
	full := evaluate(t, command(-1))
	for _, r := range rules {
		if slices.Contains(full, r.want) {
			t.Errorf("%s: command carrying every rule raised %q", r.row, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, command(i))
		if !slices.Contains(diags, r.want) {
			t.Errorf("%s: command without %q = %v, want %q", r.row, r.rule, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("%s: dropping %q also raised %q", r.row, r.rule, other.want)
			}
		}
	}
}
