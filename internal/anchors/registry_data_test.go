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
		{"RF27", "build the item in this session (\"implement now\") by default.", ".agents/commands/bench-drain.md dropped the build-in-session default for a light-path item"},
		{"RF27", "Open a `ROADMAP.md` row only when the reviewer declines.", ".agents/commands/bench-drain.md dropped the roadmap-row fallback for a light-path item the reviewer declines"},
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

	// Other rows of the group fire against this minimal tree; only the seven rule rows are
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

// TestWorktreeRuleAnchorsRedOnRemoval pins the parallel-landings guidance. Each needle and
// its diagnostic are written here independently of the registry, so a guide that dropped the
// worktree rule or the optional-spec sentence cannot define itself green.
func TestWorktreeRuleAnchorsRedOnRemoval(t *testing.T) {
	rules := []struct{ file, needle, want string }{
		{".bench/BENCH.md", "**Every phase runs in a bench worktree and lands through `bench worktree land`.**", ".bench/BENCH.md Workflow section dropped the worktree rule; every phase runs in a bench worktree and lands through bench worktree land"},
		{".bench/BENCH-reference.md", "The spec is optional on the landing and on its resume: a spec-less phase lands with no `--spec`.", ".bench/BENCH-reference.md landing paragraph dropped the optional spec; a spec-less phase lands with no --spec"},
	}
	templates := map[string]string{
		".bench/BENCH.md":           "# Bench Operating Guide\n\n## Workflow\n\n%s\n\n## Capture\n",
		".bench/BENCH-reference.md": "# Bench reference\n\n%s\n",
	}

	evaluate := func(t *testing.T, dropped int) []string {
		t.Helper()
		root := t.TempDir()
		for i, r := range rules {
			needle := r.needle
			if i == dropped {
				needle = ""
			}
			path := filepath.Join(root, filepath.FromSlash(r.file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(fmt.Sprintf(templates[r.file], needle)), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return EvaluateGroup(root, AfterSpecAuthorization)
	}

	// Other rows of the group fire against this minimal tree; only the two needle rows are
	// this test's subject, so both directions are read by membership.
	full := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(full, r.want) {
			t.Errorf("guide carrying %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("guide without %q = %v, want %q", r.needle, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("dropping %q also raised %q", r.needle, other.want)
			}
		}
	}

	// The worktree rule binds to the Workflow section: the same sentence under Capture
	// leaves Workflow without it, so the section-bound row still fires.
	root := t.TempDir()
	misplaced := "# Bench Operating Guide\n\n## Workflow\n\n## Capture\n\n" + rules[0].needle + "\n"
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "BENCH.md"), []byte(misplaced), 0o644); err != nil {
		t.Fatal(err)
	}
	if diags := EvaluateGroup(root, AfterSpecAuthorization); !slices.Contains(diags, rules[0].want) {
		t.Errorf("guide with the worktree rule under Capture = %v, want %q", diags, rules[0].want)
	}
}

func TestWorktreeEnforcementAnchorRedOnRemoval(t *testing.T) {
	const enforcement = "`bench commit` enforces this boundary: it refuses the primary checkout"
	const want = ".bench/BENCH.md Workflow section dropped worktree enforcement; bench commit refuses the primary checkout"
	evaluate := func(t *testing.T, enforcement string) []string {
		t.Helper()
		root := t.TempDir()
		bench := "# Bench Operating Guide\n\n## Workflow\n\n" + WorktreeRuleMarker + "\n" + enforcement + "\n\n## Capture\n"
		path := filepath.Join(root, ".bench", "BENCH.md")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(bench), 0o644); err != nil {
			t.Fatal(err)
		}
		return EvaluateGroup(root, AfterSpecAuthorization)
	}
	if diags := evaluate(t, enforcement); slices.Contains(diags, want) {
		t.Fatalf("guide carrying enforcement raised %q", want)
	}
	if diags := evaluate(t, ""); !slices.Contains(diags, want) {
		t.Fatalf("guide without enforcement = %v, want %q", diags, want)
	}
}

// TestFastLaneAnchorsRedOnRemoval pins OG26 and OG27. Each needle and its diagnostic are
// written here independently of the registry, so a guide that dropped the invariant-4
// sentence, or a reference whose landing shape names only the gate, cannot define
// itself green.
func TestFastLaneAnchorsRedOnRemoval(t *testing.T) {
	const invariants = "The four invariants (these override convenience, always)"
	rules := []struct{ file, needle, want string }{
		{".bench/BENCH.md", "Green is the landing's whole-project gate, and a worktree commit requires a lane pass, not a gate run.", ".bench/BENCH.md invariant 4 dropped the fast-lane sentence; green is the landing's whole-project gate, and a worktree commit requires a lane pass, not a gate run"},
		{".bench/BENCH-reference.md", "A worktree `bench commit` runs the fast lane on a private checkout of the composed snapshot, and a lane pass publishes onto the worktree branch.", ".bench/BENCH-reference.md landing shape dropped the fast lane; a worktree bench commit runs the lane, and the landing runs the one whole-project gate"},
	}
	templates := map[string]string{
		".bench/BENCH.md":           "# Bench Operating Guide\n\n## " + invariants + "\n\n4. Commit on green, never on red. %s\n\n## Workflow\n\n" + WorktreeRuleMarker + "\n" + WorktreeEnforcementMarker + "\n\n## Capture\n",
		".bench/BENCH-reference.md": "# Bench reference\n\n## Command Notes\n\n%s\n\nThe spec is optional on the landing and on its resume.\n\n## Hook Layers\n",
	}
	if !strings.Contains(rules[1].needle, "lane") {
		t.Fatalf("OG27 needle %q does not name the lane", rules[1].needle)
	}

	evaluate := func(t *testing.T, dropped int) []string {
		t.Helper()
		root := t.TempDir()
		for i, r := range rules {
			needle := r.needle
			if i == dropped {
				needle = ""
			}
			path := filepath.Join(root, filepath.FromSlash(r.file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(fmt.Sprintf(templates[r.file], needle)), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return EvaluateGroup(root, AfterSpecAuthorization)
	}

	full := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(full, r.want) {
			t.Errorf("guide carrying %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("guide without %q = %v, want %q", r.needle, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("dropping %q also raised %q", r.needle, other.want)
			}
		}
	}

	// The sentence binds to the invariants section: the same sentence under Workflow
	// leaves the invariants without it, so the section-bound row still fires.
	root := t.TempDir()
	misplaced := "# Bench Operating Guide\n\n## " + invariants + "\n\n4. Commit on green, never on red.\n\n## Workflow\n\n" + rules[0].needle + "\n"
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "BENCH.md"), []byte(misplaced), 0o644); err != nil {
		t.Fatal(err)
	}
	if diags := EvaluateGroup(root, AfterSpecAuthorization); !slices.Contains(diags, rules[0].want) {
		t.Errorf("guide with the fast-lane sentence under Workflow = %v, want %q", diags, rules[0].want)
	}
}
