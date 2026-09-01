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

// evaluate writes one minimal tree that carries every needle a conformant tree
// holds, except the broken one, and returns the group's diagnostics. A broken index
// of -1 writes the conformant tree. Rules that share a file join into one body, and
// rules that share a section group under one heading, because a per-rule write or a
// repeated heading would erase or split the earlier needle.
func (h anchorHarness) evaluate(t *testing.T, broken int) []string {
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
	return EvaluateGroup(root, h.group)
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
	for i, r := range h.rules {
		diags := h.evaluate(t, i)
		if !slices.Contains(diags, r.want) {
			t.Errorf("tree broken at %q in %s = %v, want %q", r.needle, r.file, diags, r.want)
		}
		if h.skipCrossTalk {
			continue
		}
		for j, other := range h.rules {
			if j != i && slices.Contains(diags, other.want) {
				t.Errorf("breaking %q in %s also raised %q", r.needle, r.file, other.want)
			}
		}
	}
}

// TestRetroFeedsMarkerAnchorRedsOnRemoval pins RF28. The needle and the diagnostic are
// written here independently of the registry, so a needle edited to match a template that
// dropped the destination marker cannot define itself green.
func TestRetroFeedsMarkerAnchorRedsOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:   ".agents/commands/bench-final-check.md",
				needle: "End the item with one line that reads `Feeds: FT<n>`, `Feeds: new`, or `Feeds: none`.",
				want:   ".agents/commands/bench-final-check.md dropped the implementation-retro improvement-item destination marker",
			},
		},
		templates: map[string]string{
			".agents/commands/bench-final-check.md": "# Final check\n\n## Capture the implementation retro\n\nWrite each improvement item as one list item.%s\n",
		},
	}.check(t)
}

// TestDrainFlowRuleAnchorsRedOnRemoval pins RF26 and RF27. Each drain rule is written here
// independently of the registry, next to the diagnostic that must name it, so the anchor
// cannot be re-derived green from a command file that dropped one rule while the rest
// survived.
func TestDrainFlowRuleAnchorsRedOnRemoval(t *testing.T) {
	const drain = ".agents/commands/bench-drain.md"
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{file: drain, needle: "Run `bench roadmap --flow` once and quote its flow block in the exit.", want: ".agents/commands/bench-drain.md dropped the flow-quote rule: the exit quotes bench roadmap --flow"},
			{file: drain, needle: "An entry feeds a row only when it changes the row's priority, scope, or `Next:`.", want: ".agents/commands/bench-drain.md dropped the feeds-a-row test: priority, scope, or Next:"},
			{file: drain, needle: "Dismiss an occurrence-only entry with one line of why.", want: ".agents/commands/bench-drain.md dropped the occurrence-only dismissal rule"},
			{file: drain, needle: "A new row needs a `Next:` token and a class before it opens.", want: ".agents/commands/bench-drain.md dropped the new-row rule: a Next: token and a class"},
			{file: drain, needle: "When the flow report shows a positive net delta, propose reducing moves in the next batch diff.", want: ".agents/commands/bench-drain.md dropped the positive-delta restructure rule"},
			{file: drain, needle: "build the item in this session (\"implement now\") by default.", want: ".agents/commands/bench-drain.md dropped the build-in-session default for a light-path item"},
			{file: drain, needle: "Open a `ROADMAP.md` row only when the reviewer declines.", want: ".agents/commands/bench-drain.md dropped the roadmap-row fallback for a light-path item the reviewer declines"},
		},
		templates: map[string]string{drain: "# /bench-drain\n%s"},
	}.check(t)
}

// TestWorktreeRuleAnchorsRedOnRemoval pins the parallel-landings guidance. Each needle and
// its diagnostic are written here independently of the registry, so a guide that dropped the
// worktree rule or the optional-spec sentence cannot define itself green.
func TestWorktreeRuleAnchorsRedOnRemoval(t *testing.T) {
	rules := []anchorRule{
		{file: ".bench/BENCH.md", needle: "**Every phase runs in a bench worktree and lands through `bench worktree land`.**", want: ".bench/BENCH.md Workflow section dropped the worktree rule; every phase runs in a bench worktree and lands through bench worktree land"},
		{file: ".bench/BENCH-reference.md", needle: "The spec is optional on the landing and on its resume: a spec-less phase lands with no `--spec`.", want: ".bench/BENCH-reference.md landing paragraph dropped the optional spec; a spec-less phase lands with no --spec"},
		{file: ".bench/BENCH-reference.md", needle: "Each landing refusal face constructs through the registry constructor, which takes the recovery route as a required argument.", want: ".bench/BENCH-reference.md dropped the landing refusal shape; every landing refusal face constructs through the registry constructor, which takes the recovery route as a required argument"},
	}
	// The two guides carry their own headings, so the worktree rule lands under Workflow
	// and the reference needles land in the file body.
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: rules,
		templates: map[string]string{
			".bench/BENCH.md":           "# Bench Operating Guide\n\n## Workflow\n%s\n\n## Capture\n",
			".bench/BENCH-reference.md": "# Bench reference\n%s",
		},
	}.check(t)

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
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: []anchorRule{
			{
				file:   ".bench/BENCH.md",
				needle: "`bench commit` enforces this boundary: it refuses the primary checkout",
				want:   ".bench/BENCH.md Workflow section dropped worktree enforcement; bench commit refuses the primary checkout",
			},
		},
		templates: map[string]string{
			".bench/BENCH.md": "# Bench Operating Guide\n\n## Workflow\n\n" + WorktreeRuleMarker + "\n%s\n\n## Capture\n",
		},
	}.check(t)
}

// TestFastLaneAnchorsRedOnRemoval pins OG26 and OG27. Each needle and its diagnostic are
// written here independently of the registry, so a guide that dropped the invariant-4
// sentence, or a reference whose landing shape names only the gate, cannot define
// itself green.
func TestFastLaneAnchorsRedOnRemoval(t *testing.T) {
	const invariants = "The four invariants (these override convenience, always)"
	rules := []anchorRule{
		{file: ".bench/BENCH.md", needle: "Green is the landing's whole-project gate, and a worktree commit requires a lane pass, not a gate run.", want: ".bench/BENCH.md invariant 4 dropped the fast-lane sentence; green is the landing's whole-project gate, and a worktree commit requires a lane pass, not a gate run"},
		{file: ".bench/BENCH-reference.md", needle: "A worktree `bench commit` runs the fast lane on a private checkout of the composed snapshot, and a lane pass publishes onto the worktree branch.", want: ".bench/BENCH-reference.md landing shape dropped the fast lane; a worktree bench commit runs the lane, and the landing runs the one whole-project gate"},
	}
	if !strings.Contains(rules[1].needle, "lane") {
		t.Fatalf("OG27 needle %q does not name the lane", rules[1].needle)
	}
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: rules,
		templates: map[string]string{
			".bench/BENCH.md":           "# Bench Operating Guide\n\n## " + invariants + "\n\n4. Commit on green, never on red.%s\n\n## Workflow\n\n" + WorktreeRuleMarker + "\n" + WorktreeEnforcementMarker + "\n\n## Capture\n",
			".bench/BENCH-reference.md": "# Bench reference\n\n## Command Notes\n%s\n\nThe spec is optional on the landing and on its resume.\n\n## Hook Layers\n",
		},
	}.check(t)

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
	const finalCheck = ".agents/commands/bench-final-check.md"
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:   finalCheck,
				needle: "Read `census=<n>` from the landed record; for `n > 0`, write exactly one `bench learning --rule` entry for the landing.",
				want:   ".agents/commands/bench-final-check.md post-merge tail dropped the census duty: read census=<n> and write one bench learning --rule entry",
			},
			{
				file:   finalCheck,
				needle: "For `n = 0`, state `census: 0 raw calls` in the close; a nonzero count never blocks a landing.",
				want:   ".agents/commands/bench-final-check.md post-merge tail dropped the zero census close and its advisory rule",
			},
			{
				file:   finalCheck,
				needle: "Its `--what` lists each verb head with its count. Its `--right` names the Bench form per head, or `none`. Its `--rule` proposes the verb or the help change.",
				want:   ".agents/commands/bench-final-check.md post-merge tail dropped the census learning fields: --what, --right, and --rule",
			},
			{
				file:   finalCheck,
				needle: "A spec retro cites the landing's census entry under `### Bench CLI` with its `Feeds:` line.",
				want:   ".agents/commands/bench-final-check.md retro section dropped the census citation under ### Bench CLI with its Feeds: line",
			},
			{
				file:   ".agents/skills/bench-craft-delegate/SKILL.md",
				needle: "Ask the delegate for zero to two Bench CLI improvements derived from its own calls, and fold them into the landing's census entry.",
				want:   ".agents/skills/bench-craft-delegate/SKILL.md charge dropped the delegate's zero to two Bench CLI improvements",
			},
			{
				file:   ".bench/BENCH-reference.md",
				needle: "The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`.",
				want:   ".bench/BENCH-reference.md dropped the census signal account: raw calls per assignment under $BENCH_HOME/census/<repo-key>/",
			},
			{
				file:   "projects/benchkit.md",
				needle: "The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`.",
				want:   "projects/benchkit.md dropped the census signal account: raw calls per assignment under $BENCH_HOME/census/<repo-key>/",
			},
		},
	}.check(t)
}

// TestHarnessRecordPointerAnchorsRedOnRemoval pins HC45 and HC47. Each needle and its
// diagnostic are written here independently of the registry, so a reference whose Hook
// Layers bullet dropped the record pointer, or a changelog that never names the verb,
// cannot define itself green. The changelog needle is a prefix of the reference needle,
// so the two rows do not separate by cross-talk.
func TestHarnessRecordPointerAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: []anchorRule{
			{
				file:   ".bench/BENCH-reference.md",
				needle: "bench harnesses codex",
				want:   ".bench/BENCH-reference.md Hook Layers dropped the pointer to bench harnesses codex; the record holds the Codex delegation-guard verdict with its source and its date",
			},
			{
				file:   "CHANGELOG.md",
				needle: "bench harnesses",
				want:   "CHANGELOG.md dropped the bench harnesses CLI-addition entry",
			},
		},
		skipCrossTalk: true,
	}.check(t)
}

// TestCraftSpecMapDisciplineAnchorsRedOnRemoval holds the spec-authoring rules that
// guidance must keep: a build may not edit its own spec's rows, budgets, or fences, the
// review rubric asks the per-row question, and the coverage map points at its discipline
// reference. Each section, needle, and diagnostic is written here independently of the
// registry, so guidance that drops a rule cannot define itself green.
func TestCraftSpecMapDisciplineAnchorsRedOnRemoval(t *testing.T) {
	const file = ".agents/skills/bench-craft-spec/SKILL.md"
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:      file,
				section:   "Template",
				needle:    "## Further notes\nThe flagged-additions list",
				want:      ".agents/skills/bench-craft-spec/SKILL.md Template put contract text under the Further notes heading; that heading stays bare",
				forbidden: true,
				fenced:    true,
			},
			{
				file:    file,
				section: "Slicing a build for delegates",
				needle:  "A build may not edit its own spec's acceptance rows, budget targets, or ownership fences.",
				want:    ".agents/skills/bench-craft-spec/SKILL.md Slicing a build for delegates dropped the rule that a build may not edit its own spec's acceptance rows, budget targets, or ownership fences",
			},
			{
				file:    file,
				section: "Review rubric",
				needle:  "Per row, does the map name the gate check or test that reds it, or mark the row review-owned?",
				want:    ".agents/skills/bench-craft-spec/SKILL.md Review rubric dropped the per-row question that names the gate check or test which reds the row",
			},
			{
				file:    file,
				section: "The acceptance coverage map",
				needle:  "`references/map-discipline.md` states the rule each row must satisfy",
				want:    ".agents/skills/bench-craft-spec/SKILL.md The acceptance coverage map dropped the pointer to references/map-discipline.md",
			},
			{
				file:   file,
				needle: "This reader sweep includes `.mjs` scripts and workflow files, and `references/map-discipline.md` states its rules.",
				want:   ".agents/skills/bench-craft-spec/SKILL.md dropped the reader-sweep name from the whole-tree sweep step",
			},
		},
	}.check(t)
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
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
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
			{
				file:   writeSpec,
				needle: "Run `craft-spec`'s reader sweep before that lock.",
				want:   ".agents/commands/bench-write-spec.md dropped the reader sweep sequenced before the craft-tickets charge",
			},
			{
				file:   writeSpec,
				needle: "ship on its own gate? Apply `craft-spec`'s named `Bootstrap authority before execution` rule. The ticket graph splits where a consumer branch lands green alone.",
				want:   ".agents/commands/bench-write-spec.md dropped the consumer-branch split arm from the moved ship-test question",
			},
		},
	}.check(t)
}

// TestLoadStopAndQuietAnchorsRedOnRemoval keeps the stop and aggregate-readiness
// conditions independently red-capable. The compact skill pointers lead the
// coordinator to the reference that owns the detailed process.
func TestLoadStopAndQuietAnchorsRedOnRemoval(t *testing.T) {
	const (
		line      = ".agents/skills/bench-craft-line/SKILL.md"
		delegate  = ".agents/skills/bench-craft-delegate/SKILL.md"
		reference = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	)
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
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
		},
	}.check(t)
}

// TestDoneClaimOwnerAnchorsRedOnRemoval keeps owner resolution and repair
// attribution independently red-capable. A done claim needs a tree artifact,
// and an umbrella cannot become a ledger for unrelated findings.
func TestDoneClaimOwnerAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:    ".agents/skills/bench-craft-delegate/SKILL.md",
				section: "Verifying the done-claim",
				needle:  "Resolve every named Red-mutation owner to a real artifact in the tree.",
				want:    ".agents/skills/bench-craft-delegate/SKILL.md Verifying the done-claim dropped Red-mutation owner resolution to a tree artifact",
			},
			{
				file:    ".agents/skills/bench-craft-delegate/references/delegation-discipline.md",
				section: "Before the landing",
				needle:  "Keep an accepted finding on its original ticket when attribution is clear. Use an umbrella repair ticket only for a genuinely shared owner.",
				want:    "delegation-discipline.md Before the landing dropped original-ticket attribution or the genuinely-shared-owner umbrella limit",
			},
		},
	}.check(t)
}

// TestInstalledLaneRepairAnchorsRedOnRemoval pins the installed-lane repair
// constraint. The skill routes the installed-lane case to its canonical rule,
// and the rule keeps snapshot grading
// and the post-landing rebuild coupled to the fallback.
func TestInstalledLaneRepairAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:    ".agents/skills/bench-craft-delegate/SKILL.md",
				section: "Verifying the done-claim",
				needle:  "Installed-lane repair and its post-landing rebuild are in `references/delegation-discipline.md`.",
				want:    ".agents/skills/bench-craft-delegate/SKILL.md Verifying the done-claim dropped the installed-lane repair pointer",
			},
			{
				file:    ".agents/skills/bench-craft-delegate/references/delegation-discipline.md",
				section: "Before the landing",
				needle:  "When an installed lane cannot commit its repair, run the same ordinary commit core from the candidate tree. Grade the composed snapshot, then require the sanctioned rebuild after landing.",
				want:    "delegation-discipline.md Before the landing dropped the candidate commit core, composed-snapshot grade, or sanctioned rebuild",
			},
		},
	}.check(t)
}

// TestRepairChargeTemplateAnchorsRedOnRemoval keeps each repair-charge field
// independently red-capable. A coordinator cannot verify a repair when its charge
// omits the base, fence, effort, focused suite, or independent probe. The template
// section survives each removal, so a red identifies the omitted field instead of a
// missing template section.
func TestRepairChargeTemplateAnchorsRedOnRemoval(t *testing.T) {
	const (
		file    = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
		section = "Repair-charge template"
	)
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{file: file, section: section, needle: "Base commit:", want: "delegation-discipline.md Repair-charge template dropped the base commit field"},
			{file: file, section: section, needle: "Ownership fence:", want: "delegation-discipline.md Repair-charge template dropped the ownership fence field"},
			{file: file, section: section, needle: "Effort:", want: "delegation-discipline.md Repair-charge template dropped the effort field"},
			{file: file, section: section, needle: "Focused suite:", want: "delegation-discipline.md Repair-charge template dropped the focused suite field"},
			{file: file, section: section, needle: "Independent biting probe:", want: "delegation-discipline.md Repair-charge template dropped the independent biting probe field"},
		},
		suffix: " value",
	}.check(t)
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
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:    ".agents/skills/bench-craft-tickets/SKILL.md",
				section: "Draft the breakdown",
				needle:  "A ticket that implements a roadmap row's decided fix first verifies the row's premise against the code. A premise the code contradicts is a reviewer decision, not a fix to implement as written.",
				want:    ".agents/skills/bench-craft-tickets/SKILL.md dropped the roadmap-row premise check from the breakdown method",
			},
		},
	}.check(t)
}

// TestCommentAndReviewRuleAnchorsRedOnRemoval keeps short independent anchors.
// The test does not copy the guidance prose, so the registry cannot make a
// weakened rule self-consistent. The needles are short phrases that share their
// subject sections, so the rows do not separate by cross-talk.
func TestCommentAndReviewRuleAnchorsRedOnRemoval(t *testing.T) {
	const (
		comments = ".agents/skills/bench-craft-comments/SKILL.md"
		review   = ".agents/skills/bench-craft-review/SKILL.md"
	)
	harness := anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{file: comments, section: "The register", needle: "FT<n> story <n>", want: ".agents/skills/bench-craft-comments/SKILL.md dropped the identifier-provenance rule"},
			{file: comments, section: "The register", needle: "State the constraint first", want: ".agents/skills/bench-craft-comments/SKILL.md dropped the constraint-first rule"},
			{file: comments, section: "The register", needle: "One source owns a fact", want: ".agents/skills/bench-craft-comments/SKILL.md dropped the one-source rule"},
			{file: comments, section: "The register", needle: "A sparse file stays sparse", want: ".agents/skills/bench-craft-comments/SKILL.md dropped the sparse-file rule"},
			{file: comments, section: "The register", needle: "The commit or spec owns the red record", want: ".agents/skills/bench-craft-comments/SKILL.md dropped the red-record ownership rule"},
			{file: review, section: "The axes stay separate", needle: "A new `FT<n> story <n>` tag", want: ".agents/skills/bench-craft-review/SKILL.md dropped review rejection of a new story provenance tag"},
		},
		skipCrossTalk: true,
	}
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	live := EvaluateGroup(root, AfterImplementSpec)
	for _, rule := range harness.rules {
		if slices.Contains(live, rule.want) {
			t.Errorf("live guidance raised %q", rule.want)
		}
	}
	if diags := harness.evaluate(t, -1); len(diags) == 0 {
		t.Fatal("minimal LF15 guidance did not exercise the anchor group")
	}
	harness.check(t)
}

// TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval holds the two rules a kit
// spec keeps missing. A kit behavior serves this repository and every linked repository,
// and the two audiences can want different answers for one directory state. A
// transaction-shaped spec loses its middle failure when one row covers all three. Each
// section, needle, and diagnostic is written here independently of the registry.
func TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval(t *testing.T) {
	const file = ".agents/skills/bench-craft-spec/references/map-discipline.md"
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:    file,
				section: "In the edge inventory",
				needle:  "A kit spec names the audience each behavior serves: this repository, or every repository that links the kit. The inventory walks the absent-versus-empty pair for each directory the spec reads",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md In the edge inventory dropped the two-audience prompt with its absent-versus-empty directory pair",
			},
			{
				file:    file,
				section: "Per row",
				needle:  "A transaction-shaped spec gives three rows for its verification failures. The rows are persistence before the oracle runs, interruption inside the oracle, and persistence at the terminal step.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Per row dropped the three transaction verification failure rows",
			},
			{
				file:    file,
				section: "Per row",
				needle:  "Each in-scope edge-inventory promise, source promise, and fence-closure promise takes one red-capable row.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Per row dropped the one-red-capable-row rule for each in-scope promise",
			},
			{
				file:    file,
				section: "In the edge inventory",
				needle:  "Each excluded edge takes a Won't handle line that names a surviving in-scope caller.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md In the edge inventory dropped the surviving-in-scope-caller clause from the Won't handle line",
			},
			{
				file:    file,
				section: "Before the map locks",
				needle:  "The flagged-additions list sits under Further notes before the first review charge.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Before the map locks dropped the flagged-additions list under Further notes before the first review charge",
			},
			{
				file:    file,
				section: "Before the map locks",
				needle:  "The source-sentence-to-row table sits under Further notes before the first review charge.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Before the map locks dropped the source-sentence-to-row table under Further notes before the first review charge",
			},
			{
				file:    file,
				section: "At review",
				needle:  "The review round demands one row for each listed addition, and it removes each unlisted addition.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md At review dropped one arm of the addition disposition",
			},
			{
				file:    file,
				section: "Per row",
				needle:  "An either-side predicate takes two rows, one side per row. One row that names both sides is not sufficient.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Per row dropped the two-row rule for an either-side predicate",
			},
			{
				file:    file,
				section: "Before the map locks",
				needle:  "Each canary row and each conformance row traces to its executed root before the coverage map locks.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Before the map locks dropped the executed-root trace for a canary or conformance row",
			},
			{
				file:    file,
				section: "Per row",
				needle:  "Each named diagnostic state is addable or mutable in a fixture.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Per row dropped the fixture-reachable rule for a named diagnostic state",
			},
			{
				file:    file,
				section: "At ticket slicing",
				needle:  "The author quotes each pasted operand in the delegate charge.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md At ticket slicing dropped the quoted-operand rule for a pasted operand in the delegate charge",
			},
			{
				file:    file,
				section: "Before the map locks",
				needle:  "The reader sweep lists each named consumer of the decision fact.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Before the map locks dropped the named-consumer rule for the reader sweep",
			},
			{
				file:    file,
				section: "Before the map locks",
				needle:  "The reader sweep lists each helper that a named consumer calls directly.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Before the map locks dropped the direct-helper rule for the reader sweep",
			},
			{
				file:    file,
				section: "Before the map locks",
				needle:  "A deeper callee joins the reader sweep only when the callee reads the decision fact.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Before the map locks dropped the depth bound that admits a deeper callee only when the callee reads the decision fact",
			},
			{
				file:    file,
				section: "Before the map locks",
				needle:  "Each shared reader in the reader sweep takes an exact ownership fence.",
				want:    ".agents/skills/bench-craft-spec/references/map-discipline.md Before the map locks dropped the exact-ownership-fence rule for a shared reader",
			},
			{
				file:      file,
				needle:    "reader census",
				want:      ".agents/skills/bench-craft-spec/references/map-discipline.md writes \"reader census\"; the canonical term is reader sweep",
				forbidden: true,
			},
		},
	}.check(t)
}

// TestDecisionMapAuthoringAnchorsRedOnRemoval holds the decision-map authoring steps and
// the asset-path convention. The shaping phase file must send the author to one ready
// map's Sources block, must name both first-skeleton verbs, and must keep the unnamed
// prose-check term out. The two phase files and the template must spell one asset path,
// so a drift in any one of the three files reds its own diagnostic. Each needle and
// diagnostic is written here independently of the registry.
func TestDecisionMapAuthoringAnchorsRedOnRemoval(t *testing.T) {
	const (
		shapeIdea = ".agents/commands/bench-shape-idea.md"
		writeSpec = ".agents/commands/bench-write-spec.md"
		schema    = "internal/maps/schema.go"
	)
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:   shapeIdea,
				needle: "Read one ready decision map's `## Sources` block before the first write.",
				want:   ".agents/commands/bench-shape-idea.md dropped the ready-map Sources read before the first decision-map write",
			},
			{
				file:   shapeIdea,
				needle: "Run `bench maps` and `bench gate-prose` on the first skeleton.",
				want:   ".agents/commands/bench-shape-idea.md dropped the bench-maps and bench-gate-prose checks on the first decision-map skeleton",
			},
			{
				file:      shapeIdea,
				needle:    "prose preflight",
				want:      ".agents/commands/bench-shape-idea.md writes an unnamed prose check; `bench gate-prose` is the one handle",
				forbidden: true,
			},
			{
				file:   shapeIdea,
				needle: "decisions/assets/",
				want:   ".agents/commands/bench-shape-idea.md dropped the decisions/assets/ path for a map-owned asset",
			},
			{
				file:   writeSpec,
				needle: "decisions/assets/",
				want:   ".agents/commands/bench-write-spec.md dropped the decisions/assets/ path for a map-owned asset",
			},
			{
				file:   schema,
				needle: "decisions/assets/",
				want:   "internal/maps/schema.go dropped the decisions/assets/ path for a map-owned asset from the decision-map template",
			},
		},
		templates: map[string]string{schema: "package maps\n%s"},
	}.check(t)
}

// TestContextMapTermAnchorsRedOnRemoval holds the three map-term glossary entries in
// CONTEXT.md. The coverage map and the decision map are two artifacts, and each Avoid
// list names the bare word "map". The reader sweep entry reserves "census". Each needle
// and diagnostic is written here independently of the registry, so a glossary that
// merged the two map entries cannot define itself green.
func TestContextMapTermAnchorsRedOnRemoval(t *testing.T) {
	const file = "CONTEXT.md"
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: []anchorRule{
			{
				file:   file,
				needle: "Not \"map\", not \"traceability matrix\" — coverage map.",
				want:   "CONTEXT.md dropped the coverage-map glossary entry with the Avoid list that names the bare word map",
			},
			{
				file:   file,
				needle: "Not \"PRD\", not \"design doc\", not \"map\" — decision map.",
				want:   "CONTEXT.md decision-map entry dropped the Avoid list that names the bare word map",
			},
			{
				file:   file,
				needle: "Not \"census\", not \"consumer audit\" — reader sweep.",
				want:   "CONTEXT.md dropped the reader-sweep glossary entry with the Avoid list that reserves census",
			},
		},
	}.check(t)
}

// TestCraftGateBothEndsAnchorsRedOnRemoval holds the two rules a gate author reads at the
// two places the author already opens. A check on an indirected value stays green while
// the producer and the consumer disagree, so the check grades both ends, or their binding,
// in one change. A new check that no single edit defeats is a check the author never saw
// bite. Each section, needle, and diagnostic is written here independently of the registry.
func TestCraftGateBothEndsAnchorsRedOnRemoval(t *testing.T) {
	const file = ".agents/skills/bench-craft-gate/SKILL.md"
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:    file,
				section: "Run the real path",
				needle:  "A check on a workflow output, config key, or environment variable grades the producer and the consumer, or their binding, in the same change.",
				want:    ".agents/skills/bench-craft-gate/SKILL.md Run the real path dropped the both-ends rule for a check on an indirected value",
			},
			{
				file:    file,
				section: "Prove it bites",
				needle:  "The author asks which single edit defeats a new check while the gate stays green.",
				want:    ".agents/skills/bench-craft-gate/SKILL.md Prove it bites dropped the single-edit defeat question",
			},
		},
	}.check(t)
}

// TestRepairTicketOwnerAnchorsRedOnRemoval holds the two rules that keep an accepted
// repair on the coverage map. A repair that amends a mapped row leaves no ticket behind,
// so the row loses its owner at the final check. Each section, needle, and diagnostic is
// written here independently of the registry.
func TestRepairTicketOwnerAnchorsRedOnRemoval(t *testing.T) {
	const file = ".agents/commands/bench-review-implementation.md"
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:    file,
				section: "Review modes",
				needle:  "writes one repair ticket before the repair-scoped re-review",
				want:    ".agents/commands/bench-review-implementation.md Review modes dropped the repair ticket the coordinator writes before the repair-scoped re-review",
			},
			{
				file:    file,
				section: "Review modes",
				needle:  "it cites each amended row in `Covers:`.",
				want:    ".agents/commands/bench-review-implementation.md Review modes dropped the repair ticket's amended-row citation in `Covers:`",
			},
		},
	}.check(t)
}

// TestStandingFalsificationAnchorsRedOnRemoval holds the rules that make a second
// harness read every kit-guidance diff. The review phase file states the standing pass,
// names the set by path, gives a falsification finding its three labels, and bridges an
// accepted one to the repair-routing label. The recipes name the exec form the guard
// allows, and the build phase file scopes its ask-before-adding rule and keeps the
// retired offer sentence out. Each section, needle, and diagnostic is written here
// independently of the registry.
func TestStandingFalsificationAnchorsRedOnRemoval(t *testing.T) {
	const (
		review   = ".agents/commands/bench-review-implementation.md"
		recipes  = ".agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md"
		fullSpec = ".agents/commands/bench-implement-spec.md"
	)
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:    review,
				section: "Review modes",
				needle:  "A diff that changes kit guidance takes a standing cross-harness falsification pass.",
				want:    ".agents/commands/bench-review-implementation.md Review modes dropped the standing cross-harness falsification pass for a kit-guidance diff",
			},
			{
				file:    review,
				section: "Review modes",
				needle:  "The kit-guidance set is any file under `.agents/` or the file `.bench/BENCH.md`.",
				want:    ".agents/commands/bench-review-implementation.md Review modes dropped the kit-guidance set of any file under `.agents/` plus the file `.bench/BENCH.md`",
			},
			{
				file:    review,
				section: "Review modes",
				needle:  "Each falsification finding takes one explicit disposition of accept, merge, or dismiss.",
				want:    ".agents/commands/bench-review-implementation.md Review modes dropped the accept, merge, or dismiss disposition a falsification finding takes",
			},
			{
				file:    review,
				section: "Review modes",
				needle:  "An accepted falsification finding joins the review findings and takes the repair-routing disposition.",
				want:    ".agents/commands/bench-review-implementation.md Review modes dropped the bridge that gives an accepted falsification finding the repair-routing disposition",
			},
			{
				file:   recipes,
				needle: "bench worktree exec <target> -- claude -p --model <id> \"<charge>\" <<'EOF'",
				want:   ".agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md dropped the exec reviewer form with the empty quoted heredoc",
			},
			{
				file:    fullSpec,
				section: "`--full <spec>`",
				needle:  "Outside the kit-guidance set, a diff large enough to hide bugs",
				want:    ".agents/commands/bench-implement-spec.md `--full` section dropped the kit-guidance-set scope on the ask-before-adding rule",
			},
			{
				file:      fullSpec,
				needle:    "Both are offers; the command never applies them silently.",
				want:      ".agents/commands/bench-implement-spec.md retained the retired sentence that calls the tier escalation and the falsification pass both offers",
				forbidden: true,
			},
		},
	}.check(t)
}
