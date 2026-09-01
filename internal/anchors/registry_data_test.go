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

// TestBoundedGateOutputAnchorTuples pins BG27. The two needles are written here
// independently of the registry, so a needle edited to match prose that dropped the
// bounded-output account cannot define itself green.
func TestBoundedGateOutputAnchorTuples(t *testing.T) {
	wanted := []Anchor{
		{
			Group:   AfterSpecAuthorization,
			File:    "projects/benchkit.md",
			Kind:    Require,
			Section: "",
			Needle:  "A green run prints one `phases[N]{phase,verdict,elapsed_ms}` table.",
		},
		{
			Group:   AfterSpecAuthorization,
			File:    ".bench/BENCH-reference.md",
			Kind:    Require,
			Section: "",
			Needle:  "A red run prints one `failures[N]{phase,line}` table, and then `gate: red`.",
		},
	}
	for _, want := range wanted {
		matches := 0
		for _, entry := range Entries() {
			if entry.Group == want.Group && entry.File == want.File && entry.Kind == want.Kind && entry.Section == want.Section && entry.Needle == want.Needle {
				matches++
			}
		}
		if matches != 1 {
			t.Errorf("registry has %d rows matching %+v; want exactly one", matches, want)
		}
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

// TestCensusDutyAnchorsRedOnRemoval pins EC27, EC29, EC31, EC32, and EC33. Each needle and
// its diagnostic are written here independently of the registry, so a command, a skill, a
// reference, or a profile that dropped the census account cannot define itself green. The
// reference and the profile carry the same sentence, so a dropped file must raise only its
// own diagnostic.
func TestCensusDutyAnchorsRedOnRemoval(t *testing.T) {
	rules := []struct{ file, needle, want string }{
		{
			".agents/commands/bench-final-check.md",
			"Read `census=<n>` from the landed record; for `n > 0`, write exactly one `bench learning --rule` entry for the landing.",
			".agents/commands/bench-final-check.md post-merge tail dropped the census duty: read census=<n> and write one bench learning --rule entry",
		},
		{
			".agents/commands/bench-final-check.md",
			"For `n = 0`, state `census: 0 raw calls` in the close; a nonzero count never blocks a landing.",
			".agents/commands/bench-final-check.md post-merge tail dropped the zero census close and its advisory rule",
		},
		{
			".agents/commands/bench-final-check.md",
			"Its `--what` lists each verb head with its count. Its `--right` names the Bench form per head, or `none`. Its `--rule` proposes the verb or the help change.",
			".agents/commands/bench-final-check.md post-merge tail dropped the census learning fields: --what, --right, and --rule",
		},
		{
			".agents/commands/bench-final-check.md",
			"A spec retro cites the landing's census entry under `### Bench CLI` with its `Feeds:` line.",
			".agents/commands/bench-final-check.md retro section dropped the census citation under ### Bench CLI with its Feeds: line",
		},
		{
			".agents/skills/bench-craft-delegate/SKILL.md",
			"Ask the delegate for zero to two Bench CLI improvements derived from its own calls, and fold them into the landing's census entry.",
			".agents/skills/bench-craft-delegate/SKILL.md charge dropped the delegate's zero to two Bench CLI improvements",
		},
		{
			".bench/BENCH-reference.md",
			"The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`.",
			".bench/BENCH-reference.md dropped the census signal account: raw calls per assignment under $BENCH_HOME/census/<repo-key>/",
		},
		{
			"projects/benchkit.md",
			"The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`.",
			"projects/benchkit.md dropped the census signal account: raw calls per assignment under $BENCH_HOME/census/<repo-key>/",
		},
	}

	// evaluate writes one minimal tree that carries every needle except the dropped one.
	// A file that owns two needles keeps the other one, so each row is read alone.
	evaluate := func(t *testing.T, dropped int) []string {
		t.Helper()
		root := t.TempDir()
		bodies := map[string]string{}
		for i, r := range rules {
			if i == dropped {
				continue
			}
			bodies[r.file] += r.needle + "\n\n"
		}
		for _, r := range rules {
			path := filepath.Join(root, filepath.FromSlash(r.file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("# subject\n\n"+bodies[r.file]), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	// Other rows of the group fire against this minimal tree; only the five census rows are
	// this test's subject, so both directions are read by membership.
	full := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(full, r.want) {
			t.Errorf("tree carrying %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("tree without %q in %s = %v, want %q", r.needle, r.file, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("dropping %q from %s also raised %q", r.needle, r.file, other.want)
			}
		}
	}
}

// TestHarnessRecordPointerAnchorsRedOnRemoval pins HC45 and HC47. Each needle and its
// diagnostic are written here independently of the registry, so a reference whose Hook
// Layers bullet dropped the record pointer, or a changelog that never names the verb,
// cannot define itself green.
func TestHarnessRecordPointerAnchorsRedOnRemoval(t *testing.T) {
	rules := []struct{ file, needle, want string }{
		{
			".bench/BENCH-reference.md",
			"bench harnesses codex",
			".bench/BENCH-reference.md Hook Layers dropped the pointer to bench harnesses codex; the record holds the Codex delegation-guard verdict with its source and its date",
		},
		{
			"CHANGELOG.md",
			"bench harnesses",
			"CHANGELOG.md dropped the bench harnesses CLI-addition entry",
		},
	}

	// evaluate writes one minimal tree that carries every needle except the dropped one.
	evaluate := func(t *testing.T, dropped int) []string {
		t.Helper()
		root := t.TempDir()
		bodies := map[string]string{}
		for i, r := range rules {
			if i == dropped {
				continue
			}
			bodies[r.file] += r.needle + "\n\n"
		}
		for _, r := range rules {
			path := filepath.Join(root, filepath.FromSlash(r.file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("# subject\n\n"+bodies[r.file]), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return EvaluateGroup(root, AfterSpecAuthorization)
	}

	// Other rows of the group fire against this minimal tree; only these two rows are the
	// test's subject, so both directions are read by membership.
	full := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(full, r.want) {
			t.Errorf("tree carrying %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("tree without %q in %s = %v, want %q", r.needle, r.file, diags, r.want)
		}
	}
}

// TestCraftSpecMapDisciplineAnchorsRedOnRemoval holds the spec-authoring rules that
// guidance must keep: a build may not edit its own spec's rows, budgets, or fences, the
// review rubric asks the per-row question, and the coverage map points at its discipline
// reference. Each section, needle, and diagnostic is written here independently of the
// registry, so guidance that drops a rule cannot define itself green.
func TestCraftSpecMapDisciplineAnchorsRedOnRemoval(t *testing.T) {
	const file = ".agents/skills/bench-craft-spec/SKILL.md"
	rules := []struct{ section, needle, want string }{
		{
			"Slicing a build for delegates",
			"A build may not edit its own spec's acceptance rows, budget targets, or ownership fences.",
			".agents/skills/bench-craft-spec/SKILL.md Slicing a build for delegates dropped the rule that a build may not edit its own spec's acceptance rows, budget targets, or ownership fences",
		},
		{
			"Review rubric",
			"Per row, does the map name the gate check or test that reds it, or mark the row review-owned?",
			".agents/skills/bench-craft-spec/SKILL.md Review rubric dropped the per-row question that names the gate check or test which reds the row",
		},
		{
			"The acceptance coverage map",
			"`references/map-discipline.md` states the rule each row must satisfy",
			".agents/skills/bench-craft-spec/SKILL.md The acceptance coverage map dropped the pointer to references/map-discipline.md",
		},
	}

	// evaluate writes one minimal skill that carries every section, and every needle except
	// the dropped one. The section survives the drop, so a red reports the missing sentence
	// rather than a missing section.
	evaluate := func(t *testing.T, dropped int) []string {
		t.Helper()
		root := t.TempDir()
		body := "# subject\n"
		for i, r := range rules {
			body += "\n## " + r.section + "\n\n"
			if i != dropped {
				body += r.needle + "\n"
			}
		}
		path := filepath.Join(root, filepath.FromSlash(file))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	// Other rows of the group fire against this minimal tree; only these three rows are the
	// test's subject, so both directions are read by membership.
	full := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(full, r.want) {
			t.Errorf("tree carrying %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("tree without %q in %s = %v, want %q", r.needle, file, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("dropping %q also raised %q", r.needle, other.want)
			}
		}
	}
}

// TestCraftDelegateDisciplineAnchorsRedOnRemoval holds the delegation rules that guidance
// must keep: a read-only delegate that reads a graded tree gets its own worktree, a mutation
// probe mutates behavior, the skill points at its discipline reference, and a review round
// declares its iteration cap. It also holds the retired no-worktree sentence out of the
// skill. Each needle and diagnostic is written here independently of the registry, so
// guidance that drops a rule cannot define itself green.
func TestCraftDelegateDisciplineAnchorsRedOnRemoval(t *testing.T) {
	const (
		skill     = ".agents/skills/bench-craft-delegate/SKILL.md"
		writeSpec = ".agents/commands/bench-write-spec.md"
	)
	rules := []struct {
		file, section, needle, want string
		forbidden                   bool
	}{
		{
			file:    skill,
			section: "Isolation",
			needle:  "A read-only delegate that reads a tree the coordinator will grade runs in its own worktree",
			want:    ".agents/skills/bench-craft-delegate/SKILL.md Isolation dropped the own-worktree rule for a read-only delegate that reads a graded tree",
		},
		{
			file:    skill,
			section: "The charge",
			needle:  "A mutation probe requires a behavioral mutation",
			want:    ".agents/skills/bench-craft-delegate/SKILL.md The charge dropped the behavioral-mutation requirement for a mutation probe",
		},
		{
			file:   skill,
			needle: "`references/delegation-discipline.md` holds the rest of the discipline",
			want:   ".agents/skills/bench-craft-delegate/SKILL.md dropped the pointer to references/delegation-discipline.md",
		},
		{
			file:      skill,
			needle:    "Read-only delegations need no worktree",
			want:      ".agents/skills/bench-craft-delegate/SKILL.md retained the retired rule that a read-only delegation needs no worktree",
			forbidden: true,
		},
		{
			file:   writeSpec,
			needle: "The round declares its iteration cap before the first charge",
			want:   ".agents/commands/bench-write-spec.md dropped the review round's iteration-cap declaration",
		},
	}

	// evaluate writes one minimal file per subject. A required needle is present unless it
	// is the broken one; a forbidden needle is present only when it is the broken one. Every
	// section survives the break, so a red reports the sentence rather than a missing section.
	evaluate := func(t *testing.T, broken int) []string {
		t.Helper()
		root := t.TempDir()
		bodies := map[string]string{skill: "# subject\n", writeSpec: "# subject\n"}
		sections := map[string]map[string]bool{}
		for i, r := range rules {
			present := (i == broken) == r.forbidden
			if r.section == "" {
				if present {
					bodies[r.file] += "\n" + r.needle + "\n"
				}
				continue
			}
			if !sections[r.file][r.section] {
				if sections[r.file] == nil {
					sections[r.file] = map[string]bool{}
				}
				sections[r.file][r.section] = true
				bodies[r.file] += "\n## " + r.section + "\n\n"
			}
			if present {
				bodies[r.file] += r.needle + "\n"
			}
		}
		for file, body := range bodies {
			path := filepath.Join(root, filepath.FromSlash(file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	// Other rows of the group fire against this minimal tree; only these rows are the test's
	// subject, so both directions are read by membership.
	conformant := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(conformant, r.want) {
			t.Errorf("tree conformant with %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("tree broken at %q in %s = %v, want %q", r.needle, r.file, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("breaking %q also raised %q", r.needle, other.want)
			}
		}
	}
}

// TestLoadStopAndQuietAnchorsRedOnRemoval keeps LF4's stop and aggregate-readiness
// conditions independently red-capable. The compact skill pointers lead the
// coordinator to the reference that owns the detailed process.
func TestLoadStopAndQuietAnchorsRedOnRemoval(t *testing.T) {
	const (
		line      = ".agents/skills/bench-craft-line/SKILL.md"
		delegate  = ".agents/skills/bench-craft-delegate/SKILL.md"
		reference = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	)
	rules := []struct {
		file, section, needle, want string
	}{
		{
			file:    line,
			section: "Classify reds before the ladder moves",
			needle:  "Known-flaky retry stops are in `craft-delegate`'s delegation discipline.",
			want:    ".agents/skills/bench-craft-line/SKILL.md Classify reds before the ladder moves dropped the known-flaky retry-stop pointer",
		},
		{
			file:    delegate,
			section: "Verifying the done-claim",
			needle:  "Before retry coordination or aggregate grading, load the stopped-retry and quiet-grade rules from `references/delegation-discipline.md`.",
			want:    ".agents/skills/bench-craft-delegate/SKILL.md Verifying the done-claim dropped the stopped-retry and quiet-grade pointer",
		},
		{
			file:    reference,
			section: "Retry stops and aggregate readiness",
			needle:  "After the second known-flaky refusal proves green in isolation, stop coordination and hand both results to the reviewer.",
			want:    "delegation-discipline.md Retry stops and aggregate readiness dropped the second proven flaky-refusal stop",
		},
		{
			file:    reference,
			section: "Retry stops and aggregate readiness",
			needle:  "Before aggregate grading, wait until returned delegates have no live tests and serialize the coordinator-owned resource.",
			want:    "delegation-discipline.md Retry stops and aggregate readiness dropped the quiet-delegate aggregate-grade check",
		},
	}

	evaluate := func(t *testing.T, broken int) []string {
		t.Helper()
		root := t.TempDir()
		bodies := map[string]string{}
		sections := map[string]map[string]bool{}
		for i, r := range rules {
			if sections[r.file] == nil {
				sections[r.file] = map[string]bool{}
				bodies[r.file] = "# subject\n"
			}
			if !sections[r.file][r.section] {
				sections[r.file][r.section] = true
				bodies[r.file] += "\n## " + r.section + "\n\n"
			}
			if i != broken {
				bodies[r.file] += r.needle + "\n"
			}
		}
		for file, body := range bodies {
			path := filepath.Join(root, filepath.FromSlash(file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	conformant := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(conformant, r.want) {
			t.Errorf("tree conformant with %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("tree broken at %q in %s = %v, want %q", r.needle, r.file, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("breaking %q also raised %q", r.needle, other.want)
			}
		}
	}
}

// TestDoneClaimOwnerAnchorsRedOnRemoval keeps owner resolution and repair
// attribution independently red-capable. A done claim needs a tree artifact,
// and an umbrella cannot become a ledger for unrelated findings.
func TestDoneClaimOwnerAnchorsRedOnRemoval(t *testing.T) {
	const (
		skill     = ".agents/skills/bench-craft-delegate/SKILL.md"
		reference = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	)
	rules := []struct {
		file, section, needle, want string
	}{
		{
			file:    skill,
			section: "Verifying the done-claim",
			needle:  "Resolve every named Red-mutation owner to a real artifact in the tree.",
			want:    ".agents/skills/bench-craft-delegate/SKILL.md Verifying the done-claim dropped Red-mutation owner resolution to a tree artifact",
		},
		{
			file:    reference,
			section: "Before the landing",
			needle:  "Keep an accepted finding on its original ticket when attribution is clear. Use an umbrella repair ticket only for a genuinely shared owner.",
			want:    "delegation-discipline.md Before the landing dropped original-ticket attribution or the genuinely-shared-owner umbrella limit",
		},
	}

	evaluate := func(t *testing.T, broken int) []string {
		t.Helper()
		root := t.TempDir()
		for i, r := range rules {
			path := filepath.Join(root, filepath.FromSlash(r.file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			body := "# subject\n\n## " + r.section + "\n\n"
			if i != broken {
				body += r.needle + "\n"
			}
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	conformant := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(conformant, r.want) {
			t.Errorf("tree conformant with %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("tree broken at %q in %s = %v, want %q", r.needle, r.file, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("breaking %q also raised %q", r.needle, other.want)
			}
		}
	}
}

// TestInstalledLaneRepairAnchorsRedOnRemoval pins LF11. The skill routes the
// installed-lane case to its canonical rule, and the rule keeps snapshot grading
// and the post-landing rebuild coupled to the fallback.
func TestInstalledLaneRepairAnchorsRedOnRemoval(t *testing.T) {
	const (
		skill     = ".agents/skills/bench-craft-delegate/SKILL.md"
		reference = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	)
	rules := []struct {
		file, section, needle, want string
	}{
		{
			file:    skill,
			section: "Verifying the done-claim",
			needle:  "Installed-lane repair and its post-landing rebuild are in `references/delegation-discipline.md`.",
			want:    ".agents/skills/bench-craft-delegate/SKILL.md Verifying the done-claim dropped the installed-lane repair pointer",
		},
		{
			file:    reference,
			section: "Before the landing",
			needle:  "When an installed lane cannot commit its repair, run the same ordinary commit core from the candidate tree. Grade the composed snapshot, then require the sanctioned rebuild after landing.",
			want:    "delegation-discipline.md Before the landing dropped the candidate commit core, composed-snapshot grade, or sanctioned rebuild",
		},
	}

	evaluate := func(t *testing.T, broken int) []string {
		t.Helper()
		root := t.TempDir()
		for i, r := range rules {
			path := filepath.Join(root, filepath.FromSlash(r.file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			body := "# subject\n\n## " + r.section + "\n\n"
			if i != broken {
				body += r.needle + "\n"
			}
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	conformant := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(conformant, r.want) {
			t.Errorf("tree conformant with %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("tree broken at %q in %s = %v, want %q", r.needle, r.file, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("breaking %q also raised %q", r.needle, other.want)
			}
		}
	}
}

// TestRepairChargeTemplateAnchorsRedOnRemoval keeps each repair-charge field
// independently red-capable. A coordinator cannot verify a repair when its charge
// omits the base, fence, effort, focused suite, or independent probe.
func TestRepairChargeTemplateAnchorsRedOnRemoval(t *testing.T) {
	const (
		file    = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
		section = "Repair-charge template"
	)
	rules := []struct{ needle, want string }{
		{"Base commit:", "delegation-discipline.md Repair-charge template dropped the base commit field"},
		{"Ownership fence:", "delegation-discipline.md Repair-charge template dropped the ownership fence field"},
		{"Effort:", "delegation-discipline.md Repair-charge template dropped the effort field"},
		{"Focused suite:", "delegation-discipline.md Repair-charge template dropped the focused suite field"},
		{"Independent biting probe:", "delegation-discipline.md Repair-charge template dropped the independent biting probe field"},
	}

	// evaluate preserves the template section while it removes one field. The red
	// therefore identifies the omitted field instead of a missing template section.
	evaluate := func(t *testing.T, broken int) []string {
		t.Helper()
		root := t.TempDir()
		body := "# subject\n\n## " + section + "\n\n"
		for i, r := range rules {
			if i != broken {
				body += r.needle + " value\n"
			}
		}
		path := filepath.Join(root, filepath.FromSlash(file))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	conformant := evaluate(t, -1)
	for _, r := range rules {
		if slices.Contains(conformant, r.want) {
			t.Errorf("tree carrying %q raised %q", r.needle, r.want)
		}
	}
	for i, r := range rules {
		diags := evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("tree without %q = %v, want %q", r.needle, diags, r.want)
		}
		for j, other := range rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("dropping %q also raised %q", r.needle, other.want)
			}
		}
	}
}

// TestReferenceFileAnchorsRedOnAbsence holds the two discipline references in the tree. A
// skill that points at a reference the tree lost leaves the reader with a dead pointer, so
// an absent file raises the missing-file diagnostic. The paths and the lead sentences are
// written here independently of the registry.
func TestReferenceFileAnchorsRedOnAbsence(t *testing.T) {
	references := []struct{ file, lead, want string }{
		{
			".agents/skills/bench-craft-delegate/references/delegation-discipline.md",
			"Charged from `craft-delegate` when the coordinator writes a charge, runs a probe, or accepts a return.",
			".agents/skills/bench-craft-delegate/references/delegation-discipline.md is absent or dropped its charge-time lead",
		},
		{
			".agents/skills/bench-craft-spec/references/map-discipline.md",
			"Charged from `craft-spec` when the author writes or audits an acceptance coverage map.",
			".agents/skills/bench-craft-spec/references/map-discipline.md is absent or dropped its map-time lead",
		},
	}
	for _, reference := range references {
		missing := "acceptance coverage anchor file missing: " + reference.file
		present := EvaluateGroup(writeReferences(t, references, reference.file), AfterImplementSpec)
		if slices.Contains(present, missing) {
			t.Errorf("tree carrying %s raised %q", reference.file, missing)
		}
		if slices.Contains(present, reference.want) {
			t.Errorf("tree carrying %s raised %q", reference.file, reference.want)
		}
		absent := EvaluateGroup(writeReferences(t, references, ""), AfterImplementSpec)
		if !slices.Contains(absent, missing) {
			t.Errorf("tree without %s = %v, want %q", reference.file, absent, missing)
		}
	}
}

// writeReferences builds a tree that carries the named reference file and no other. It
// returns the tree root.
func writeReferences(t *testing.T, references []struct{ file, lead, want string }, keep string) string {
	t.Helper()
	root := t.TempDir()
	for _, reference := range references {
		if reference.file != keep {
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(reference.file))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("# subject\n\n"+reference.lead+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// TestCraftTicketsPremiseCheckAnchorRedsOnRemoval holds the breakdown method's premise
// check. A ticket author who implements a roadmap row's decided fix without reading the
// code builds the row's mistake, so the rule stays in the skill. The sentence and the
// diagnostic are written here independently of the registry.
func TestCraftTicketsPremiseCheckAnchorRedsOnRemoval(t *testing.T) {
	const (
		skill   = ".agents/skills/bench-craft-tickets/SKILL.md"
		section = "Draft the breakdown"
		needle  = "A ticket that implements a roadmap row's decided fix first verifies the row's premise against the code. A premise the code contradicts is a reviewer decision, not a fix to implement as written."
		want    = ".agents/skills/bench-craft-tickets/SKILL.md dropped the roadmap-row premise check from the breakdown method"
	)

	// evaluate writes one minimal subject. The section survives the break, so a red
	// reports the missing sentence rather than a missing section.
	evaluate := func(t *testing.T, present bool) []string {
		t.Helper()
		root := t.TempDir()
		body := "# subject\n\n## " + section + "\n\n"
		if present {
			body += needle + "\n"
		}
		path := filepath.Join(root, filepath.FromSlash(skill))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	if diags := evaluate(t, true); slices.Contains(diags, want) {
		t.Errorf("tree carrying the premise check raised %q", want)
	}
	if diags := evaluate(t, false); !slices.Contains(diags, want) {
		t.Errorf("tree without the premise check = %v, want %q", diags, want)
	}
}

// TestCommentAndReviewRuleAnchorsRedOnRemoval keeps short independent anchors for
// LF15. The test does not copy the guidance prose, so the registry cannot make a
// weakened rule self-consistent.
func TestCommentAndReviewRuleAnchorsRedOnRemoval(t *testing.T) {
	const (
		comments = ".agents/skills/bench-craft-comments/SKILL.md"
		review   = ".agents/skills/bench-craft-review/SKILL.md"
	)
	rules := []struct{ file, section, needle, want string }{
		{comments, "The register", "FT<n> story <n>", ".agents/skills/bench-craft-comments/SKILL.md dropped the identifier-provenance rule"},
		{comments, "The register", "State the constraint first", ".agents/skills/bench-craft-comments/SKILL.md dropped the constraint-first rule"},
		{comments, "The register", "One source owns a fact", ".agents/skills/bench-craft-comments/SKILL.md dropped the one-source rule"},
		{comments, "The register", "A sparse file stays sparse", ".agents/skills/bench-craft-comments/SKILL.md dropped the sparse-file rule"},
		{comments, "The register", "The commit or spec owns the red record", ".agents/skills/bench-craft-comments/SKILL.md dropped the red-record ownership rule"},
		{review, "The axes stay separate", "A new `FT<n> story <n>` tag", ".agents/skills/bench-craft-review/SKILL.md dropped review rejection of a new story provenance tag"},
	}
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	live := EvaluateGroup(root, AfterImplementSpec)
	for _, rule := range rules {
		if slices.Contains(live, rule.want) {
			t.Errorf("live guidance raised %q", rule.want)
		}
	}

	evaluate := func(t *testing.T, dropped int) []string {
		t.Helper()
		root := t.TempDir()
		bodies := map[string]*strings.Builder{
			comments: &strings.Builder{},
			review:   &strings.Builder{},
		}
		fmt.Fprint(bodies[comments], "## The register\n\n")
		fmt.Fprint(bodies[review], "## The axes stay separate\n\n")
		for i, rule := range rules {
			if i == dropped {
				continue
			}
			body := bodies[rule.file]
			fmt.Fprintf(body, "%s\n\n", rule.needle)
		}
		for file, body := range bodies {
			path := filepath.Join(root, file)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(body.String()), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return EvaluateGroup(root, AfterImplementSpec)
	}

	if diags := evaluate(t, -1); len(diags) == 0 {
		t.Fatal("minimal LF15 guidance did not exercise the anchor group")
	} else {
		for _, rule := range rules {
			if slices.Contains(diags, rule.want) {
				t.Errorf("guidance carrying %q raised %q", rule.needle, rule.want)
			}
		}
	}
	for i, rule := range rules {
		if diags := evaluate(t, i); !slices.Contains(diags, rule.want) {
			t.Errorf("guidance without %q = %v, want %q", rule.needle, diags, rule.want)
		}
	}
}
