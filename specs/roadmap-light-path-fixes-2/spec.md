# Roadmap light-path fixes, batch 2

Status: staged

Roadmap: FT286, FT272, FT288, FT289, FT266, FT282, FT291

Decision source: the reviewer-confirmed conversation of 2026-09-02, over the eight named roadmap rows, checked against the live tree

Verification log: 2 iteration(s) to accept — the first round blocked on two unordered `Writes:` overlaps, one undisposed FT115 clause, and one unflagged addition; the second pass confirmed the fold and the acyclic graph

## Problem

Eight roadmap rows describe bounded fixes that each fit one ticket. Three are
code defects in the oracle and its readers. The prose bound resets at any
label-shaped line, so a long paragraph passes. `bench status` routes a shipped
tickets-only folder to a grammar `bench commit` refuses. Fourteen test waits
carry literal deadlines, so machine load produces reds the diff does not own.

Five rows are kit-guidance sentences that the last two cycles earned. Each row
holds one owner file, and each new sentence must ship with its anchor, its
registry test, and its live-mirror fixture. Shipped one at a time, every row
pays a full gate. Shipped together under one landing, they pay one.

## Solution

Keep one reviewed umbrella for scheduling and retirement. Split every row at
its executable seam. Each ticket owns one small behavior, focused evidence, and
an exact write fence. The three code fixes run in parallel. The guidance
tickets share the anchor registry and the fixture directory, so they land in
one serial chain on the retained integration source.

Preserve behavior already present in the live tree. Tickets implement only the
remaining delta. Each ticket first verifies its row's premise against the code
it names.

## User stories

### Group A — the prose oracle

Line: opus / low. The seam is one function with a table test, and the gate
covers it.

1. As a guidance author, I want a label-shaped prose sentence to stay in its paragraph, so that the six-sentence bound bites.
2. As a roadmap author, I want an `Occurrence:` line to stay a field line although it ends with a period, so that a ledger passes.
3. As a ticket author, I want a `Writes:` or `Blocked by:` line to stay ungraded, so that a template does not fail the prose lane.
4. As a reviewer, I want the STE prose reference to state the label rule the code applies, so that the doc and the oracle agree.
5. As a gate operator, I want a label-shaped line with only a mid-line terminator to count as prose, so that the cheapest wrong fix reds.

### Group B — the status route

Line: opus / low. The change is one definition string and two pinned tests.

6. As a coordinator, I want `bench status` to route a shipped tickets-only folder to `bench worktree land --spec <slug>`, so that the printed command is one `bench commit --help` does not refuse.
7. As a coordinator, I want the route to keep the `<slug>` argument, so that the row still names the folder it closes.
8. As a maintainer, I want every action definition to render and parse the same command, so that a new route cannot break the round trip.

### Group C — wait deadlines

Line: opus / low for the migration, opus / medium for the sweep. The migration
is exact and covered; the sweep generalizes an existing check at a known seam.

9. As a test author, I want every outer wait deadline derived from `bounds.TestDeadline`, so that a loaded machine does not red a correct diff.
10. As a gate operator, I want an expired wait to print its name and window, so that a timeout is not an assertion.
11. As a maintainer, I want a dev-tier check to red a literal duration in a test wait deadline, so that no helper reintroduces one.
12. As a maintainer, I want the new check to reuse the marker-wait literal scanner, so that one scanner owns literal detection.
13. As a test author, I want a poll interval to stay a literal, so that the check does not red a tick value.

### Group D — kit guidance

Line: opus / medium. Guidance prose steers every later session, and the gate
covers only the anchor bytes, so the line rises one step above the code tickets.

14. As a spec author, I want `craft-spec` to derive the fence from the tickets' `Writes:` lines after the slice, so that no registry is omitted.
15. As a spec author, I want a Won't handle over an anchored sentence to quote the kept bytes, so that the kept sentence is read.
16. As a spec author, I want the map discipline to name the shipped-surface claim words, so that a `tests/` claim is caught before the gate.
17. As a coordinator, I want `craft-delegate` to bind the exec-only form to every caller, so that a raw `git -C` read is a named bypass.
18. As a coordinator, I want a shell loop inside the pool path named as the same bypass, so that a for-loop cannot escape.
19. As a charge author, I want a cap-change charge to name the closest pinning package, so that an outside pin is found before the gate.
20. As a delegate, I want `craft-tdd` to say a re-exec helper returns silently and never skips, so that no helper skip reads as environmental.
21. As a coordinator, I want the phase close to read the census before the landing removes the record, so that the per-verb breakdown survives.
22. As a coordinator, I want a light-path fix to land only under an untouched `CHANGELOG.md` heading, so that two entries do not conflict at composition.
23. As a reviewer, I want the review base named as the merged `main` tip, so that the reviewed range holds the spec diff alone.
24. As a gate operator, I want every new sentence pinned by an anchor, a registry test, and a fixture, so that a dropped sentence reds.
25. As a kit maintainer, I want each edited skill inside its prose budget row, so that the budget check does not red the landing.

## Implementation decisions

- A label line has a one-to-four-word prefix that ends at the first colon. It is a label when the line carries no sentence terminator, or when the prefix is a template field name. The field-name list is a closed constant beside the abbreviation list, and the STE prose reference states the same two-part rule.
- The status route definition keeps the one-word argument shape and changes only its command to `bench worktree land --spec`. Status prints no request, base, or tip, because the primary checkout does not know the in-flight worktree.
- `bounds` gains one timeout-verdict renderer beside `TestDeadline`, with no `testing` import. Every migrated wait derives its window from `TestDeadline` and prints that verdict on expiry.
- The literal-deadline sweep generalizes the marker-wait check: same walk, same literal scanner, keyed on a test wait's deadline argument, with a poll-interval allowlist. It registers at the dev tier beside the existing row.
- Every guidance sentence ships as five things. They are the sentence, one anchor tuple, one red-on-removal registry test, one live-mirror fixture, and a guard-table row where Go guards the behavior. The `.claude` tree is a symlink, so a ticket names `.agents` paths only.
- The guidance tickets form one serial chain, because they share the anchor registry, its test file, and the fixture directory. The ticket graph carries every `Writes:` overlap on bytes a ticket edits as a blocker edge. A fixture-closure-only `Writes:` entry, where the ticket edits no bytes the fixture pins, needs no such edge to a sibling naming the same entry. The deadline ticket follows the status route ticket over `status_producible_test.go`, and the guidance chain follows the deadline ticket over the five registry closure files.
- FT115 keeps its conformance-phase headroom clause. The drain reduces the row to that clause after this spec lands, so the `Roadmap:` line omits FT115.
- The CHANGELOG-heading constraint is a rule, not a caution. A fix whose entry shares a heading with the spec waits for the landing.
- FT282 and FT291 share `bench-final-check.md`, so one ticket carries both faces plus the review-base sentence.

## Testing decisions

- Prose behavior uses the canonical grader's table test and the live-tree prose check.
- The status route uses the two pinned status tests and the render-parse round trip.
- Wait deadlines use the migrated helpers' own tests and one new dev-tier conformance check with a bite test.
- Guidance changes use the anchor registry test, the fixture-bite sweep, and the guidance token sweep.
- Existing behavior receives preservation checks only where its caller could regress.

### Seam diagram

    roadmap row
        |
        v
    focused ticket -> package or guidance seam -> focused evidence
        |
        v
    retained integration source -> serial gate -> reviewed landing

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LP1 | 1, 5 | a label-shaped line with a mid-line terminator and no trailing terminator stays in its paragraph and the eight-sentence run reds | `TestFindings` in internal/prose | an end-of-line test keeps the split and returns no finding |
| LP2 | 2 | twenty `Occurrence:` lines that end with a period produce no finding | `TestFindings` in internal/prose | a naive no-terminator rule merges the ledger into one paragraph |
| LP3 | 3 | a `Writes:` or `Blocked by:` line with no terminator produces no finding | `TestFindings` in internal/prose | a rule that drops the label branch grades every field line |
| LP4 | 4 | the STE prose reference states the field-name clause beside the no-terminator clause | anchor registry test on `ste-prose.md` | a doc that keeps only one clause contradicts the code |
| LP5 | 1 | the live tree passes the prose check under the new rule | `TestProseMechanicsHoldsOnTheLiveTree` | a rule that reds an `Occurrence:` ledger reds ten roadmap files |
| LP6 | 6, 7 | the tickets-only residue row prints `bench worktree land --spec <slug>` | `TestAllProducibleBoardActionsAreInvocableOrEmpty` tickets-only case | the old literal or a missing `<slug>` fails the exact signal |
| LP7 | 6 | two tickets-only folders order after the earlier residue rows with the new route | `TestTicketsOnlyResidueRowCountsAndRanksBelowItsBand` | the old literal fails the struct equality |
| LP8 | 8 | the new definition renders and parses to the same command | `TestActionDefinitionsRenderAndParseTheSameCommand` | a command that breaks parse symmetry reds |
| LP9 | 9 | each of the fourteen named wait sites derives its window from `bounds.TestDeadline` | the literal-deadline sweep over the live tree | one surviving literal reds the sweep |
| LP10 | 10 | an expired wait prints the timeout verdict with the wait name and the window | the migrated helper's own test with a zero inner window | a bare `t.Fatal` message carries no window |
| LP11 | 11 | a numeric literal in a test wait deadline argument reds the dev-tier check | bite test on a synthetic Go file | a check keyed on one helper name misses the fourteen sites |
| LP12 | 12 | the new check calls the marker-wait literal scanner | review-owned: the check body names `containsNumericLiteral` | a second scanner is a one-source defect |
| LP13 | 13 | a poll interval literal inside the wait loop does not red the check | bite test with a `10ms` tick beside a derived deadline | a check with no allowlist reds every migrated helper |
| LP14 | 14, 15 | `craft-spec` carries the fence-after-slice sentence and the quoted-bytes sentence | anchor registry test and fixture `craft-spec-fence-after-slice` | a dropped sentence reds the registry test |
| LP15 | 16 | the map discipline names the shipped-surface claim words | anchor registry test and fixture `map-discipline-shipped-surface-claim` | a sweep without the words does not catch the merge-gate red |
| LP16 | 17, 18 | `craft-delegate` binds the exec-only form to every caller and names a shell loop as the same bypass | anchor registry test and fixture `delegate-exec-only-every-caller` | a sentence that binds only the charge survives the old bytes |
| LP17 | 19 | `craft-delegate` says a cap-change charge names the closest pinning package | anchor registry test and fixture `delegate-cap-change-pinning-package` | an omitted sentence reds the registry test |
| LP18 | 20 | `craft-tdd` says a re-exec helper returns silently outside its role and never skips | anchor registry test and fixture `craft-tdd-helper-returns-not-skips` | an omitted sentence reds the registry test |
| LP19 | 21 | `bench-final-check` places the census read before the landing step | anchor registry test and fixture `final-check-census-read-before-land` | a sentence after the landed record reads a deleted file |
| LP20 | 22 | `bench-final-check` states the CHANGELOG-heading rule as a rule | anchor registry test and fixture `final-check-light-path-changelog-heading` | a caution wording survives the needle |
| LP21 | 23 | `bench-review-implementation` names the merged `main` tip as the review base | anchor registry test and fixture `review-base-merged-main-tip` | an unnamed base leaves the range to the reviewer's memory |
| LP22 | 24 | every new fixture bites through its registered owner | `TestEveryRetainedFixtureBitesThroughRegisteredOwner` | a fixture whose swap is not registered passes silently |
| LP23 | 25 | each edited skill stays inside its budget row | the prose budget check | a skill over budget reds the landing |

### Edge inventory

- Error paths: a label line inside a code span, a URL scheme colon, and a five-word prefix stay prose under the current tests.
- Empty input: a status board with zero tickets-only folders prints no residue row, as today.
- Boundary values: a four-word prefix is a label and a five-word prefix is prose; a zero inner window derives the twenty-second floor.
- Interrupted state: an expired wait prints its verdict before the helper returns.
- Re-run idempotency: the sweep over the migrated tree stays green on a second run.
- Hostile paths: a needle quotes the physical line bytes, so a wrapped sentence cannot pin.
- Partial implementation: a code ticket that edits only the definition or only the tests reds its pinned pair.

**Won't handle** — the preflight face of FT272 — `bench preflight build` on a tickets-only folder stays a reviewer decision on the row; the status route survives alone.

**Won't handle** — the `bench probe` primitive that would replace the cp-aside probe — FT168 owns it; the FT289 sentence names the bypass only.

**Won't handle** — a `Sources:` line with a terminator — the field-name list covers it, and the two live rows stay green under LP5.

**Won't handle** — literal durations that are subjects under test — the sweep keys on a wait's deadline argument, so the guards tests keep their values.

**Won't handle** — the `.claude` twin of each guidance file — the tree links `.claude` to `.agents`, so no second write exists.

**Won't handle** — explicit headroom for long conformance phases — the gate runner keeps its phase deadlines, and FT115 survives with that clause.

**Won't handle** — a faster concurrency-failure path in `internal/models/models_test.go` — the derived wait trades failure-path speed for the anti-hang property; it still catches a real defect, slower.

## Ownership fences

- `specs/roadmap-light-path-fixes-2/`
- `reviews/roadmap-light-path-fixes-2.md`
- `internal/prose/parse.go`
- `internal/prose/parse_test.go`
- `internal/status/status.go`
- `internal/status/status_counters_test.go`
- `internal/status/status_producible_test.go`
- `internal/bounds/bounds.go`
- `internal/bounds/bounds_test.go`
- `internal/runbinary/runbinary_test.go`
- `internal/systemtest/owner_teardown_test.go`
- `internal/systemtest/owner_test.go`
- `internal/systemtest/otel_crash_test.go`
- `internal/systemtest/owner_artifact_recovery_test.go`
- `internal/worktree/classifier_shape_test.go`
- `internal/worktree/subshell_test.go`
- `internal/otelrecord/writer_test.go`
- `internal/gocache/lock_test.go`
- `internal/gate/prospective_owner_test.go`
- `internal/gate/run_failure_outcomes_test.go`
- `internal/conformance/marker_wait_deadline_test.go`
- `internal/conformance/wait_deadline_literal_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/registry/registry.go`
- `internal/conformance/guidance_token_sweep_test.go`
- `internal/models/models_test.go`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/skills/bench-craft-spec/references/map-discipline.md`
- `.agents/skills/bench-craft-spec/references/ste-prose.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-tdd/SKILL.md`
- `.agents/commands/bench-final-check.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/skills/bench-craft-review/SKILL.md`
- `DATA_HANDLING.md`
- `decisions/cost-follows-project-size.md`
- `decisions/craft-research.md`
- `decisions/diff-visual.md`
- `decisions/gate-budget.md`
- `decisions/gate-critical-path.md`
- `docs/reporesident-distillation.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/benchguard/benchguard_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/docs-currency-token-diet/`
- `tests/canary/data-handling-derivation/`
- `tests/canary/package-core-guard/bounds-duplicate-owner`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/main_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_test.go`
- `projects/benchkit.md`

The fence is the union of the ten tickets' `Writes:` lines, closed by
`bench preflight build` over the fixture and registry pins.

## Out of scope

- The preflight face of FT272 stays on the row: 2 edits, 1 gate run.
- Explicit headroom for long conformance phases stays on FT115: 2 edits, 1 gate run.
- A `bench probe` primitive stays FT168: a spec, 3 gate runs.
- The remaining wait helpers under `internal/guards` are subjects under test, not waits: 0 edits.
- A general status `next=` grammar is not required: 0 edits.

## Further notes

Flagged additions beyond the decision source:

- The closed field-name list in the label rule. The conversation named the no-terminator rule alone; research showed that rule reds ten roadmap files.
- The STE prose reference sentence that states the two-part rule. Without it the doc contradicts the code.
- The timeout-verdict renderer in `bounds`. The row asks to distinguish a timeout from an assertion; this is the one-source form.
- The dev-tier literal-deadline sweep and its registry row. The row asks for derived waits; the sweep is what keeps them derived, and it generalizes an existing check.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| FT286: the six-sentence bound is defeatable by an ordinary prose line | LP1, LP5 |
| FT286: a template field stays a label | LP2, LP3, LP4 |
| FT272: the route names the landing form that closes a tickets-only folder | LP6, LP7, LP8 |
| FT115: derive the waits from `bounds.TestDeadline` | LP9, LP11, LP12, LP13 |
| FT115: distinguish timeout verdicts from assertions | LP10 |
| FT288 face one and face two | LP14, LP15 |
| FT289 face one and face two | LP16, LP17 |
| FT266 | LP18 |
| FT282 | LP19 |
| FT291 face one and face two | LP20, LP21 |
| the precedent's five-part guidance shape | LP22, LP23 |

Reviewer decisions closed in conversation on 2026-09-02:

- FT286 takes the two-part label rule; the sign-off is the veto surface.
- FT272 ships the status route only.
- FT291's CHANGELOG constraint is a rule.
- Every subagent runs `opus` at low or medium effort, per the 2026-08-26 routing memory.

Reviewer decision closed during implementation, 2026-09-02: the label-line
ticket's own run of `TestProseMechanicsHoldsOnTheLiveTree` found ten
over-length sites outside its fence. This is a material LP5 shortfall the
spec's research did not anticipate. The reviewer chose a ninth ticket
(`trim-over-length-prose-in-the-live-tree`) over widening the label-line
ticket's fence or narrowing LP5. The ticket is Blocked by the label-line
ticket and carries the ten sites; the Ownership fences list above carries its
`Writes:` union.

Fence widening accepted during implementation, 2026-09-02: the sweep ticket's
own live-tree run found one more literal wait deadline. It sits outside the
fourteen the deadline ticket named, at `internal/models/models_test.go:59`.
Unlike LP5, this is a one-line fix with a one-line import. It lands inside
the sweep ticket rather than a new one, per the standing "fix, don't park"
rule.
The Ownership fences list above carries the added file.
