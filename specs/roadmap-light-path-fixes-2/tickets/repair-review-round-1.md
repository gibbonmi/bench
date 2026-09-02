# Repair review round 1

Blocked by: state-the-census-read-the-changelog-rule-and-the-review-base.md
Writes: internal/conformance/wait_deadline_literal_test.go, internal/gate/run_failure_outcomes_test.go, internal/prose/parse.go, internal/prose/parse_test.go, .agents/skills/bench-craft-spec/references/ste-prose.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/, specs/roadmap-light-path-fixes-2/spec.md, tests/canary/workflow-guidance-anchors/ste-prose-paragraph-bound, tests/canary/workflow-guidance-anchors/ste-prose-sentence-bound, tests/canary/workflow-guidance-anchors/ste-prose-two-part-label-rule, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LP9, LP11, LP1, LP2, LP3

Accepted repairs from the initial review round
(`reviews/roadmap-light-path-fixes-2.md`), reviewer sign-off 2026-09-02.

## What to build

**Repair A — widen the wait-deadline sweep's spelling set (Coverage finding 1).**
`checkWaitDeadlineLiterals` in `internal/conformance/wait_deadline_literal_test.go`
recognized only `time.After` and `time.Now().Add`. Add `time.NewTimer(d)` as
the same deadline-argument shape. Share `literalDeadline` and the
`derivedDeadline` stop rule. Derive `internal/gate/prospective_owner_test.go:367`
from `bounds.TestDeadline` — the site this diff itself migrated, the
finding's actual headline case.

Narrowed during repair, reviewer sign-off 2026-09-02: `context.WithTimeout`
and `time.AfterFunc` do not widen. The finding's named
`run_failure_outcomes_test.go:34` site spells `context.WithTimeout(ctx,
time.Second)`. `time.Second` is a named constant, not a numeric literal, so
`containsNumericLiteral` never saw it. The finding's characterization of
that site was wrong. `context.WithTimeout`'s other live sites sit in
`internal/guards/guards_test.go`. This spec's own Won't-handle already
covers guards subjects under test.

`time.AfterFunc`'s one live site, `internal/gate/phases_test.go:103`, is a
scheduled-callback delay before a fixture cancels a running phase. It is
not a deadline a wait outlasts. Migrating it broke
`TestFixturePhaseCancelledRunExitsOneHundredThirty`: the derived window let
the phase finish first. `run_failure_outcomes_test.go:34` still migrates to
`bounds.TestDeadline(0)`, a real wall-clock guess, in fence and harmless.

**Repair B — close the `templateFields` gap (Standards finding 3, Coverage
finding 2).** Add `supports` and `drift` to the closed `templateFields`
constant in `internal/prose/parse.go`. Extend the seven-name clause in
`ste-prose.md:31` to nine. Add one guard test that reds if the code list and
the doc's named list diverge. Read the doc's clause with a small parser or a
literal string match, whichever `parse_test.go` already has a precedent
for. Verify the premise first: confirm the 67 `Supports:`/`Drift:` lines
across `decisions/*.md` pass under the widened list.

**Repair C — three spec-only text fixes (Standards findings 1, 2, 4; Spec
findings 1, 2; and the models_test.go disposition).** All in
`specs/roadmap-light-path-fixes-2/spec.md` unless noted:

- Add one sentence to Implementation decisions: a fixture-closure-only
  `Writes:` entry needs no `Blocked by:` edge to a sibling naming the same
  entry.
- Add `tests/canary/data-handling-derivation/` to Ownership fences.
- Change "the union of the eight tickets'" to "the union of the nine
  tickets'".
- Add one Won't-handle line naming `internal/models/models_test.go`'s
  slower concurrency-failure path as an accepted tradeoff of the anti-hang
  property.
- Fix the two anchor diagnostics in `internal/anchors/registry_data.go`
  added by the census ticket: the strings "post-merge tail" and
  "gate-then-commit path" name no real heading. Reword each to start with
  its actual `Section` value, "Exit handoff", matching every other
  `RequireInSection` diagnostic's convention.

## Acceptance

- [ ] `TestWaitDeadlineLiteralsBites` covers all three spellings (two
      original, one new), each with a red and a green case.
- [ ] The full dev-tier live-tree run is green with the migrated site
      included.
- [ ] `TestFindings` covers a terminated `Supports:` and a terminated
      `Drift:` line producing no finding.
- [ ] The new code-vs-doc guard test reds when either list drops a name the
      other keeps, and passes on the live tree.
- [ ] `TestProseMechanicsHoldsOnTheLiveTree` stays green.
- [ ] The two reworded anchor diagnostics still red on removal of their
      sentence and stay silent on the live root.
- [ ] `bench preflight review roadmap-light-path-fixes-2` is green on the
      repaired tip.
- [ ] Self-probe: revert the `time.NewTimer` migration at
      `prospective_owner_test.go:367`, and report the sweep check red with
      that file named.
