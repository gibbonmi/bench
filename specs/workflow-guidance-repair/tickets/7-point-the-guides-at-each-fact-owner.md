# Point the guides at each fact owner

Blocked by: 6-bind-each-ticket-author-to-the-declared-line.md
Writes: .bench/BENCH.md, AGENTS.md, .bench/BENCH-reference.md, projects/benchkit.md, .claude/README.md, cmd/bench/testdata/anchors/fixture-repo/AGENTS.md, cmd/bench/anchor_help_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, tests/canary/workflow-guidance-anchors/, tests/canary/docs-currency-token-diet/, tests/canary/load-validity-metadata/, tests/canary/skills-index-command-adapters/, tests/canary/line-routing/, tests/canary/skill-description-budgets/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/guidance-prose-budgets/
Covers: GR76, GR77, GR78, GR79, GR80, GR81, GR82, GR83, GR84, GR85, GR86, GR87, GR88, GR89, GR90, GR91, GR107, GR108, GR114, GR121, GR123, GR124

## What to build

Each pointer in the operating guide, the working agreement, the reference, the profile, and the Claude README names the file that holds its fact. The reference owns the kit facts that a linked repository needs.

Before this ticket starts, the reviewer grants a one-time permission rule for `.bench/BENCH.md`, `AGENTS.md`, and `.claude/README.md`. `.bench/BENCH.md` is at its budget, so its edit stays line-neutral. The `AGENTS.md` edit adds no marker phrase of `.bench/BENCH.md`.

Make these changes:

- In `.bench/BENCH.md`, name the command registry as the owner of the plumbing subcommands.
- In `AGENTS.md`, name `.bench/BENCH-reference.md` as the holder of how the pieces fit and the skills index.
- In the reference's "Plumbing subcommands" section, name `bench gate-prose` as the one internal verb that a session runs, with its forms under Command Notes.
- Move the gate output account to the reference. Replace the profile's copy, the green-run and stream-log sentences included, with a pointer. Move the green-run tuple in `TestBoundedGateOutputAnchorTuples` to the reference.
- Remove the profile's census sentence and its Require row, and update `TestCensusDutyAnchorsRedOnRemoval`.
- Move the handoff verb facts from `AGENTS.md` to the reference. The reference sentence is "`bench handoff` rewrites only the calling worktree's assignment section."
- `AGENTS.md` keeps its phase-close rule and points to the reference for the verb behavior.
- Move that needle's registry row to the reference, retarget the canary `agents-handoff-section-rule`, and drop the moved line from the anchors fixture.
- In the reference, name `bench handoff`, `bench worktree reset`, and `bench retro` without their grammar.
- State in the reference that the plan amendment makes the authored version 1 fence a version 2 plan before the first dispatch.
- Replace the reference's skill-link claim with a pointer to `.claude/README.md`. That file states that `.claude/skills/` links every `.agents/skills/` skill that has no same-named command.
- In `.claude/README.md`, remove "links only the `bench-craft-*` skills". Replace "`bench-writer` runs a user-directed write delegation" with the three write roles that ticket 1 names.
- In `projects/benchkit.md`, remove the claim that `bench commit` works on any branch. `bench commit` refuses the primary checkout. Add a Forbid row for the removed claim.
- In the light-path cell of `.bench/BENCH.md`, replace "gate and commit on green" with the commit and the landing route. Invariant 4 owns the lane pass and the landing gate, so the cell does not restate them. Keep the edit line-neutral, and update the row that pins the cell in place.

`bench anchors AGENTS.md` prints every `AGENTS.md` registry row, a Forbid row included. Change `TestAnchorsReportsNeedleLines` and `TestAnchorsReportsAbsentNeedles` so that each expected row takes its kind from the registry. A Forbid row reads line 0 and adds no missing-file diagnostic. Plant the GR79 needle in the `AGENTS.md` fixture.

Add each planned Require and Forbid row in the registry.

## Acceptance

- [ ] The operating guide names the command registry as the plumbing owner (GR77).
- [ ] The working agreement names the reference as the holder of the pieces and the skills index (GR79).
- [ ] The reference names `bench gate-prose` as the internal verb that a session runs (GR80).
- [ ] The reference owns the green-run gate output sentence (GR82) and the handoff verb sentence (GR85).
- [ ] The reference states the version 2 amendment (GR89).
- [ ] The Claude README states the skill-link rule (GR91).
- [ ] The guides contain none of the sentences that GR76, GR78, GR81, GR83, GR84, GR86, GR87, GR88, and GR90 name.
- [ ] `TestAnchorsReportsNeedleLines` prints each `AGENTS.md` row with its registry kind (GR107).
- [ ] `TestAnchorsReportsAbsentNeedles` gives a missing-file diagnostic for each Require row only (GR108).
- [ ] The README and the profile contain none of the sentences that GR114, GR121, GR123, and GR124 name.
- [ ] The profile does not claim that `bench commit` works on any branch.
- [ ] The light-path cell names the commit and the landing route, and it does not restate invariant 4.
