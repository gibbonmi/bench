# Require a row for each in-scope promise

Blocked by: none
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, .agents/skills/bench-craft-spec/SKILL.md, tests/canary/workflow-guidance-anchors/craft-spec-further-notes-bare (new), internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/map-discipline-promise-rows (new), tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller (new), tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions (new), tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table (new), tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition (new), tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows (new), tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace (new), tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state (new), tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SAD10, SAD11, SAD12, SAD13, SAD14, SAD15, SAD16, SAD17, SAD18, SAD19, SAD20, SAD21

## What to build

`.agents/skills/bench-craft-spec/references/map-discipline.md` gains the
coverage-row rules. Nine new normative sentences land there. Each sentence ships
in the FT144 shape: the sentence, one anchor tuple in
`internal/anchors/registry_data.go`, one red-on-removal registry test, and one
canary fixture under `tests/canary/workflow-guidance-anchors/`.

The nine rules are these:

- Each in-scope edge-inventory, source-promise, and fence-closure promise takes
  one red-capable row.
- Each excluded edge takes a Won't handle line with a surviving in-scope caller.
- The flagged-additions list sits under Further notes before the first review
  charge.
- The source-sentence-to-row table sits under Further notes before the first
  review charge.
- The review round demands a row for each listed addition, and it removes each
  unlisted addition.
- An either-side predicate takes two rows, one side per row.
- Each canary or conformance row traces to the executed root before the map
  locks.
- Each named diagnostic state is addable or mutable in a fixture.
- The author quotes each pasted operand in the charge.

The reference file carries no line budget, so the whole contract lands there.
`.agents/skills/bench-craft-spec/SKILL.md` sits at 150 of 152 lines, and its
Further notes template heading stays bare. A one-line addition under that heading
does not cross the budget, so the budget cannot grade the heading. Add a
forbid-in-section tuple over that file's `## Template` section instead, and give
it the fixture `tests/canary/workflow-guidance-anchors/craft-spec-further-notes-bare`.

Two rows grade this ticket's diff rather than a behavior it adds. Row SAD16 rests
on the ownership fence: the fence list omits `internal/coverage/`, so an edit to
the coverage checker falls outside every fence and preflight refuses the landing.
Do not touch `internal/coverage`. Row SAD21 is review-owned on the Standards
axis: the authoritative-inventory rule and the named-test-function rule already
exist in `map-discipline.md`, and each must appear exactly once after your diff.

The five `cmd/bench` and `internal/conformance` entries in `Writes:` are the
registry closure of the `internal/anchors` package. Edit one only when your
change reaches it.

## Acceptance

- [ ] SAD10 — the reference requires one red-capable row for every in-scope promise.
- [ ] SAD11 — the reference requires a Won't handle line with a surviving in-scope caller.
- [ ] SAD12 — the reference requires the flagged-additions list before the first review charge.
- [ ] SAD13 — the reference requires the source-sentence-to-row table before the first review charge.
- [ ] SAD14 — the reference states both arms of the addition disposition.
- [ ] SAD15 — the reference requires two rows for an either-side predicate.
- [ ] SAD16 — the diff touches no file under `internal/coverage/`.
- [ ] SAD17 — the reference requires an executed-root trace for each canary or conformance row.
- [ ] SAD18 — the reference requires a named diagnostic state to be addable or mutable in a fixture.
- [ ] SAD19 — the reference requires the author to quote each pasted operand in the charge.
- [ ] SAD20 — a forbid-in-section tuple keeps the Further notes template heading bare.
- [ ] SAD21 — the authoritative-inventory rule and the named-test-function rule appear once each.
- [ ] Each of the ten new fixtures reds on its mutation and greens on its restore.

## Delegate charge

You work in the Bench repo on the `spec-authoring-discipline` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/spec-authoring-discipline/spec.md` first. Then read
`.agents/skills/bench-craft-spec/references/map-discipline.md` in full. Read the
anchor tuples in `internal/anchors/registry_data.go`. Read
`TestCraftSpecMapDisciplineAnchorsRedOnRemoval` at line 541 and
`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval` at line 1148 in
`internal/anchors/registry_data_test.go`. Read
`tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory/` as
the fixture prior art. That fixture holds `BASE`, `EXPECT`, and `MUTATE.json`.

Write the nine sentences in ASD-STE100 prose. Keep each anchor needle on one
physical line. Give each sentence one anchor tuple and one fixture.

Do not edit `internal/coverage`. `bench coverage --check` stays unchanged by
decision.

Add the forbid-in-section tuple over the `## Template` section of
`.agents/skills/bench-craft-spec/SKILL.md`. Add no contract text under that
file's Further notes heading. That file has two spare lines.

Do not restate the authoritative-inventory rule or the named-test-function rule.
Both already exist in the reference file.

Coverage rows: SAD10 through SAD21. Show each anchored row red before your edit. Show
each row green after. Return the red-to-green log.

Self-probe with an omission mutation. Delete the either-side sentence and report
the observed result. If the mutation returns green, add the missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/anchors/ ./internal/conformance/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
