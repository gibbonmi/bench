package anchors

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// coreInvariants is the heading of the operating guide's four-invariants section.
const coreInvariants = "The four invariants (these override convenience, always)"

// TestCoreInvariantAnchorTuples keeps the core rows on the operating guide and the retired
// Roles rule independent of production values, so an omission cannot re-derive the
// expectation green. The shared-rule check owns the Roles marker, so no row may hold it.
func TestCoreInvariantAnchorTuples(t *testing.T) {
	const bench = ".bench/BENCH.md"
	const finalRoles = "Never assume the reviewer's decisions, and never assume a claim the gate could check instead"
	if ReviewerDecisionBoundaryMarker != finalRoles {
		t.Errorf("ReviewerDecisionBoundaryMarker = %q; want %q", ReviewerDecisionBoundaryMarker, finalRoles)
	}
	wanted := []Anchor{
		{File: bench, Kind: RequireInSection, Section: "The Bench CLI contract", Needle: "Each Bench response is bounded and is the complete evidence."},
		{File: bench, Kind: Require, Needle: "`bench gate` is valid. `bench gate 2>&1 | tail -20` is not valid."},
		{File: bench, Kind: Require, Needle: "The command registry owns the plumbing subcommands that hooks and adapters drive."},
		{File: bench, Kind: RequireInSection, Section: coreInvariants, Needle: "Green is the landing's whole-project gate, and a worktree commit requires a lane pass, not a gate run."},
		{File: bench, Kind: Require, Needle: "do not call an API or a function before you read its definition in this session"},
		{File: bench, Kind: Require, Needle: "follow the current tree convention and flag it for reviewer veto"},
		{File: bench, Kind: Require, Needle: "the build exits and reports. It does not land a silent partial"},
		{File: bench, Kind: Require, Needle: "only diff-owned reds count toward fix-loop convergence"},
		{File: bench, Kind: Forbid, Needle: "NEVER assume, always verify"},
	}
	for _, want := range wanted {
		matches := 0
		for _, entry := range Entries() {
			if entry.File == want.File && entry.Kind == want.Kind && entry.Section == want.Section && entry.Needle == want.Needle {
				matches++
			}
		}
		if matches != 1 {
			t.Errorf("registry has %d rows matching %+v; want exactly one", matches, want)
		}
	}
	for _, entry := range Entries() {
		if entry.Needle == finalRoles {
			t.Errorf("registry holds the Roles marker row %+v; the shared-rule check owns it", entry)
		}
	}
}

// TestInvariantClauseAnchorsRedOnRemoval pins AP06, AP07, and AP08. The key clause of
// invariants 1, 2, and 3 and its diagnostic are written here independently of the
// registry, so a guide that dropped one clause cannot define itself green.
func TestInvariantClauseAnchorsRedOnRemoval(t *testing.T) {
	const bench = ".bench/BENCH.md"
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: []anchorRule{
			{file: bench, section: coreInvariants, needle: "Never edit, skip, weaken, or delete a test or a check to make it pass.", want: ".bench/BENCH.md invariant 1 dropped the rule that no test or check is edited, skipped, weakened, or deleted to make it pass"},
			{file: bench, section: coreInvariants, needle: "Before any multi-cycle stage, state the model, the effort, the iteration cap, and a one-clause justification in one line.", want: ".bench/BENCH.md invariant 2 dropped the one-line declaration of the model, the effort, the iteration cap, and the justification before a multi-cycle stage"},
			{file: bench, section: coreInvariants, needle: "Put no file paths and no code snippets in an ADR; they rot.", want: ".bench/BENCH.md invariant 3 dropped the rule that an ADR holds no file paths and no code snippets"},
		},
	}.check(t)
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
			Needle:     "the landing retires the sibling worktrees whose work it carries, and `bench worktree clean --landed` retires every other landed assignment",
			Diagnostic: ".agents/commands/bench-final-check.md post-merge tail dropped the landing's own sibling cleanup",
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

// TestDrainFlowRuleAnchorsRedOnRemoval pins RF26. The command line and its diagnostic are
// written here independently of the registry, so a drain that dropped the flow quote cannot
// define itself green.
func TestDrainFlowRuleAnchorsRedOnRemoval(t *testing.T) {
	const drain = ".agents/commands/bench-drain.md"
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{file: drain, needle: "Run `bench roadmap --flow` once and quote its flow block in the exit.", want: ".agents/commands/bench-drain.md dropped the flow-quote rule: the exit quotes bench roadmap --flow"},
		},
		templates: map[string]string{drain: "# /bench-drain\n%s"},
	}.check(t)
}

// TestFastLaneAnchorsRedOnRemoval pins OG26. The needle and its diagnostic are written here
// independently of the registry, so a guide that dropped the invariant-4 sentence cannot
// define itself green.
func TestFastLaneAnchorsRedOnRemoval(t *testing.T) {
	rules := []anchorRule{
		{file: ".bench/BENCH.md", needle: "Green is the landing's whole-project gate, and a worktree commit requires a lane pass, not a gate run.", want: ".bench/BENCH.md invariant 4 dropped the fast-lane sentence; green is the landing's whole-project gate, and a worktree commit requires a lane pass, not a gate run"},
	}
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: rules,
		templates: map[string]string{
			".bench/BENCH.md": "# Bench Operating Guide\n\n## " + coreInvariants + "\n\n4. Commit on green, never on red.%s\n\n## Workflow\n",
		},
	}.check(t)

	// The sentence binds to the invariants section: the same sentence under Workflow
	// leaves the invariants without it, so the section-bound row still fires.
	root := t.TempDir()
	misplaced := "# Bench Operating Guide\n\n## " + coreInvariants + "\n\n4. Commit on green, never on red.\n\n## Workflow\n\n" + rules[0].needle + "\n"
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

// TestCensusDutyAnchorsRedOnRemoval pins GR83. The needle and its diagnostic are written
// here independently of the registry. The reference owns the census signal account, so a
// profile that restores its copy raises the profile's diagnostic.
func TestCensusDutyAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:      "projects/benchkit.md",
				needle:    "The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`.",
				want:      "projects/benchkit.md restored its copy of the census signal account that .bench/BENCH-reference.md owns",
				forbidden: true,
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

// TestCraftSpecMapDisciplineAnchorsRedOnRemoval keeps contract text out from under the
// template's Further notes heading. The section, needle, and diagnostic are written here
// independently of the registry, so a template that restored the text cannot define
// itself green.
func TestCraftSpecMapDisciplineAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:      ".agents/skills/bench-craft-spec/SKILL.md",
				section:   "Template",
				needle:    "## Further notes\nThe flagged-additions list",
				want:      ".agents/skills/bench-craft-spec/SKILL.md Template put contract text under the Further notes heading; that heading stays bare",
				forbidden: true,
				fenced:    true,
			},
		},
	}.check(t)
}

// TestCraftDelegateDisciplineAnchorsRedOnRemoval keeps two retired delegation rules out:
// the cmp restore against a copy aside, and the no-worktree rule for a read-only
// delegation. Each needle and diagnostic is written here independently of the registry,
// so guidance that restored a rule cannot define itself green.
func TestCraftDelegateDisciplineAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:      ".agents/skills/bench-craft-delegate/references/delegation-discipline.md",
				forbidden: true,
				needle:    "restore with `cmp` against the copy aside",
				want:      "delegation-discipline.md restored the cmp restore against a copy aside; bench probe reports the restore",
			},
			{
				file:      ".agents/skills/bench-craft-delegate/SKILL.md",
				needle:    "Read-only delegations need no worktree",
				want:      ".agents/skills/bench-craft-delegate/SKILL.md retained the retired rule that a read-only delegation needs no worktree",
				forbidden: true,
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

// TestMapDisciplineAnchorsRedOnRemoval keeps the retired term "reader census" out of the
// map discipline. The needle and the diagnostic are written here independently of the
// registry.
func TestMapDisciplineAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:      ".agents/skills/bench-craft-spec/references/map-discipline.md",
				needle:    "reader census",
				want:      ".agents/skills/bench-craft-spec/references/map-discipline.md writes \"reader census\"; the canonical term is reader sweep",
				forbidden: true,
			},
		},
	}.check(t)
}

// TestRepairTicketOwnerAnchorsRedOnRemoval keeps out the single repair ticket for a
// coverage-map amendment, because each affected ticket's `Covers:` line takes the
// amendment. The needle and the diagnostic are written here independently of the registry.
func TestRepairTicketOwnerAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:      ".agents/commands/bench-review-implementation.md",
				forbidden: true,
				needle:    "writes one repair ticket when accepted repairs amend the coverage map",
				want:      ".agents/commands/bench-review-implementation.md restored the single repair ticket for coverage-map amendments; each affected ticket's `Covers:` line takes the amendment",
			},
		},
	}.check(t)
}

// TestCrossHarnessFalsificationAnchorsRedOnRemoval keeps out a standing cross-harness pass,
// and it holds the exec form that the recipes name. Each needle and diagnostic is written
// here independently of the registry.
func TestCrossHarnessFalsificationAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:      ".agents/commands/bench-review-implementation.md",
				needle:    "standing cross-harness",
				want:      ".agents/commands/bench-review-implementation.md restored a standing cross-harness falsification pass",
				forbidden: true,
			},
			{
				file:   ".agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md",
				needle: "bench worktree exec <target> -- claude -p --model <id> --effort <level> <<'EOF'",
				want:   ".agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md dropped the exec reviewer form that carries the charge in a quoted heredoc",
			},
		},
	}.check(t)
}

// TestChargeProbeOracleAnchorsRedOnRemoval holds the three command lines that a charge and
// a probe name: the owning check, the skip-ownership check, and the conformance package
// run. Each section, needle, and diagnostic is written here independently of the registry.
func TestChargeProbeOracleAnchorsRedOnRemoval(t *testing.T) {
	const (
		file   = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
		charge = "In the charge"
	)
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:    file,
				section: charge,
				needle:  "adds an anchor names `bench test --check <owning-check>`",
				want:    "delegation-discipline.md In the charge dropped the owning-check probe for an anchor-adding charge",
			},
			{
				file:    file,
				section: charge,
				needle:  "lists `bench test --check skip-ownership` in its focused checks",
				want:    "delegation-discipline.md In the charge dropped the skip-ownership check from a skippable test's focused checks",
			},
			{
				file:    file,
				section: "Probes",
				needle:  "`bench test --package ./internal/conformance` is not the root conformance",
				want:    "delegation-discipline.md Probes dropped the rule that the conformance package run is not the root conformance pass",
			},
		},
	}.check(t)
}

// TestDelegateExecOnlyAndCapChangeAnchorsRedOnRemoval holds the one command form for every
// caller into an assignment worktree. The section, needle, and diagnostic are written here
// independently of the registry.
func TestDelegateExecOnlyAndCapChangeAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:    ".agents/skills/bench-craft-delegate/SKILL.md",
				section: "Isolation",
				needle:  "`bench worktree exec \"<label>\" -- <command>` is the one command form for every caller into an assignment worktree.",
				want:    ".agents/skills/bench-craft-delegate/SKILL.md Isolation binds the exec-only command form to the charge instead of every caller",
			},
		},
	}.check(t)
}

// TestSystemSuiteRouteAnchorsRedOnRemoval pins BF21. The working agreement and the profile
// each name `bench test --check system` as the hand-run route for the system suite. The
// needles and the diagnostics are written here independently of the registry, so guidance
// that dropped the route cannot define itself green.
func TestSystemSuiteRouteAnchorsRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: []anchorRule{
			{
				file:   "AGENTS.md",
				needle: "Run the system suite by hand through `bench test --check system`.",
				want:   "AGENTS.md dropped the system-suite hand-run route; a hand run goes through bench test --check system, which supplies the sealed run binary and the kit root",
			},
			{
				file:   "projects/benchkit.md",
				needle: "A hand run of the system suite goes through `bench test --check system`.",
				want:   "projects/benchkit.md dropped the system-suite hand-run route; a hand run goes through bench test --check system, which supplies the sealed run binary and the kit root",
			},
		},
	}.check(t)
}

// TestHandoffSectionRuleAnchorRedOnRemoval holds the working agreement to no copy
// of the reference's handoff section rule. The needle and its diagnostic are written here
// independently of the registry, so an agreement that restored its copy cannot define
// itself green.
func TestHandoffSectionRuleAnchorRedOnRemoval(t *testing.T) {
	const rule = "`bench handoff` rewrites only the calling worktree's assignment section."
	anchorHarness{group: AfterSpecAuthorization, rules: []anchorRule{
		{file: "AGENTS.md", needle: rule, want: "AGENTS.md restored its copy of the handoff section rule that .bench/BENCH-reference.md owns", forbidden: true},
	}}.check(t)
}
