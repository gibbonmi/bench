# Debug loop trial

Status: staged

Decision source: Reviewer-confirmed conversation, 2026-09-15. The reviewer confirmed four decisions. The flag selects the arm only, and both arms record through the existing assessment step. The loop reference lives under `craft-tdd`. A delegate runs the full loop inside its fence and Phases 1 to 3 outside it. `/bench-debug` drops the `craft-delegate` fix route, and `/bench-implement-spec` keeps its "final" sentence.

Verification log: 0 iteration(s) to accept — pending the sign-off round.

## Problem

Implementation runs take longer than they should. The reviewer's research shows the agents spend that time on more than one way to implement. The `/bench-debug` phase converges fast because its loop precedes theory, its alternatives are written and capped, and its close is keyed to the loop. The build phases carry only the last half of those moves. Two related defects sit in the debug phase itself. It routes the fix to a fresh delegate, which has no loop, and it forbids a write delegate from running the phase.

## Solution

One reference states the loop discipline, and `/bench-debug` cites it. A `--debug-trial <on|off>` flag on the three build phases selects an arm. The `on` arm applies the loop moves that fit the phase. Both arms write trial provenance into the existing assessment record, so `bench assessment compare` can read them.

The debug session writes its own fix. A write delegate runs the debug phase inside its fence and the first three phases outside it. The flag and its trial reference leave the kit after the compare.

## User stories

Line: opus / high.
Implementation-line reason: The hardest material chunk is DL-C3, because the trial reference must match the compare's binding rules exactly. The spec is precise, the seams are known anchors and fixtures, and the fixture-bite test catches a dropped rule. The leverage override in `craft-line` routes guidance prose to mid + high.
Harder chunks: DL-C3.

An `on`-arm actor runs its phase under `--debug-trial on`, and an `off`-arm actor under `--debug-trial off`.

The loop is one source:

1. As an agent in any phase, I want the loop's four moves in one reference, so that every phase reads one discipline.
2. As an agent in `/bench-debug`, I want Phase 1 to cite that reference for the loop properties, so that the phase carries no second copy.
3. As a reviewer, I want the glossary to define the loop, so that phases and specs use one term.

The debug session owns the fix:

4. As a debug session, I want to write the fix myself, so that the loop stays in the context that built it.
5. As a debug session, I want to delegate only a read-only fan-out search, so that no fresh delegate rebuilds Phase 1.
6. As an in-fence write delegate, I want to run the debug phase, so that I fix the bug against a red loop.
7. As an out-of-fence write delegate, I want to run Phases 1 to 3 and return the loop command, so that the receiver starts red.
8. As a coordinator, I want the fence to decide who writes a fix outside the fence, so that one diff keeps one verdict.

The trial flag selects an arm:

9. As a reviewer, I want `--debug-trial <on|off>` on the three build phases, so that I choose the arm for each run.
10. As an `on`-arm implementer, I want the red signal run once and pasted before the first edit, so that the first action is mechanical.
11. As an `on`-arm implementer, I want a classified row to stand in for a red, so that the loop never forces a false red.
12. As an `on`-arm spec author, I want at most three ranked seam candidates per story group, so that the seam choice is bounded.
13. As an `on`-arm spec author, I want each candidate to name its red check, so that the choice is falsifiable.
14. As a reviewer of an `on`-arm spec, I want the ranking in the approval table, so that I can re-rank it.
15. As an `on`-arm review axis, I want each strong finding to carry its refutation run, so that an unrun finding is at most `ask-user`.
16. As an `off`-arm agent, I want the phase unchanged except for the trial provenance, so that the control arm measures today's phase.
17. As an agent without the flag, I want the phase text to name no loop rule, so that the control arm stays uncontaminated.

Both arms record:

18. As a reviewer, I want both arms recorded through the existing assessment step with trial provenance, so that the compare accepts them.
19. As a reviewer, I want the quality measures named for each phase, so that the plan's tolerance can reference them.
20. As a reviewer, I want an unmeasured token counter to stay unknown, so that no arm gains a manufactured zero.
21. As a reviewer, I want a plan template with two conditions on one revision, so that one capability separates the arms.
22. As a reviewer, I want the exit rule written, so that the flag does not outlive the compare.

The kit stays green:

23. As a maintainer, I want each edited budgeted file to keep its budget in the same ticket, so that the lane stays green.
24. As a maintainer, I want each new rule anchored with a fixture that bites, so that a later reflow cannot drop it silently.

## Implementation decisions

### One rule, one sentence, one anchor

Every rule this spec adds is one sentence inside the STE bound, and each rule takes its own anchor and its own fixture. A sentence that carried two rules would let an author keep the anchored half and drop the rest. The coverage map therefore holds one row per rule.

### The loop reference

A new `craft-tdd` reference, `references/loop.md`, owns the loop discipline. It states four moves in order, one sentence each. First, build one command and run it once before you read code to form a theory, and paste the invocation and the output. Second, when a choice remains after the command exists, write at most three ranked candidates, each with the prediction the command makes about it. Drop a candidate that has no prediction, act on the first, and do not block on the reviewer.

Third, change one variable per run. Fourth, close on the command, not on judgment.

The same reference owns the four properties of the loop command as the checkbox list: red-capable, deterministic, fast, and agent-runnable. `/bench-debug` Phase 1 keeps its one-command rule, its pointer to the loop constructions, and its stop-gate sentence. Phase 1 states that the command is complete when it meets the four properties in the reference. The anchor that pins the checkbox spelling moves its file to the reference.

`craft-tdd` gains one sentence that points at the reference from its cycle. `CONTEXT.md` gains a **loop** entry: the one command that produces a stage's red signal, built before theory and run after each change.

### The fix stays in the loop session

`/bench-debug` states that the session that owns the loop writes the fix. It delegates only a read-only fan-out search. The sentence that routes code authorship through `craft-delegate` leaves the phase, and a Forbid anchor keeps it out. The closing delegation-line sentence narrows to a read-only fan-out search. The words "a scoped fix" leave the phase, and a second Forbid anchor keeps them out. The worktree rule, the declared line, the seam ownership, and the gate-replacement rule stay.

The write-delegate paragraph changes shape. A write delegate runs the phase through Phase 6 on a bug inside its fence. On a bug outside its fence, it runs Phases 1 to 3 and returns the loop command in its blocked report. It stops implementation edits and keeps its in-fence work dirty in its worktree. The report also carries the red output digest, the ranked hypotheses, the failing surface, and the in-fence dirty paths. The coordinator still validates the report and reslices repair tickets under the build phase's "When the build stops short".

`craft-delegate`'s blocked-delegate sentence names the three debug phases the delegate runs before it stops. The edit replaces one sentence, so the skill stays at its line count.

### The trial flag

Each build phase carries one section titled `` `--debug-trial <on|off>` ``, backticked like the `--full` and `--delegate` sections beside it. The section is four lines: a blank line, the heading, a blank line, and one line that holds two sentences. The first sentence states the grammar: `on` or `off` selects the arm, and any other or missing value stops the phase. The second sentence points at the trial reference. No rule text sits in a phase file, and a Forbid anchor keeps `references/loop.md` out of each phase file. A restated rule that omits the path is review-owned.

A new `craft-tdd` reference, `references/loop-trial.md`, owns the trial. It opens with the `off` arm: the phase runs as written today and adds only the trial provenance to its record. It then states each `on` arm, one sentence per rule.

- **`/bench-implement-spec` on.** Before the first edit of each ticket, run the row's red signal once and paste the invocation and the output. An `already covered` or `not TDD-able` row records that classification in place of a red. Off a marked seam, the ticket's focused checks run once green as the baseline, and each later run makes one edit at the ticket's seam. When the named seam refuses the behavior, exit through "When the build stops short" as a wrong spec. The "final" sentence of the Build section stays in force.
- **`/bench-write-spec` on.** The seam sketch writes at most three ranked candidates per story group. Each candidate names the check that goes red at it, and a candidate with no red is dropped. The approval table shows the chosen seam and its ranking. `bench coverage --check` runs after the first map draft and after each row change.
- **`/bench-review-implementation` on.** Each strong finding carries the invocation and the output digest of the run that tried to refute it. A finding with no runnable check takes `ask-user` at most. A no-findings return lists the runs it made.

The trial reference owns the record contract, one sentence per field group. The plan id is `debug-loop-trial`, the conditions are `loop-on` and `loop-off`, and the capability key is `debug-loop` with the value `on` or `off`. Both conditions pin one kit revision, and each run's `source` is that revision. Each run names its condition, its task, its repetition, and `holdout: true` when the plan named the task before the run.

The plan pins the model and effort of every role the phase records. It also pins the `implementation` line in every condition, even when a phase records no implementation attempt. Each condition also pins the harness, the limits, the acceptance list, and the review-axes list, and each run's trial object repeats those four values. A counter the harness does not expose stays unknown.

The trial reference names the quality measures. `/bench-write-spec` records `verification_iterations`. `/bench-implement-spec` records `lane_first_pass` and `mutation_probe_bit` as 1 or 0. `/bench-review-implementation` records `review_findings_raw` and `review_repair_targets`. Wall time comes from the attempt intervals the record already derives.

The plan's quality tolerance carries `max_failure_rate`, which the validator requires. Its `measures` map stays empty unless every trial task records the named measure, because the compare grades each named measure against every run.

The trial reference carries the plan template as one fenced JSON block. The template's purpose is `descriptive`, its variable is `capability`, and its two conditions differ only in the `debug-loop` capability. The reviewer copies the template to `capture/benchmarks/debug-loop-trial/plan.json`, a local ignored path. Before the first run, the reviewer fills the tasks, the repetitions, the revision, the harness, the lines, and the limits. The reviewer also fills the acceptance list, the review-axes list, the tolerance, the approval reference, and the budget. The report's first limit line is `descriptive evidence only; not default-change evidence`, and that line is the expected result, because adoption is a reviewer decision.

The trial reference states the exit in two sentences. After the compare, one landing removes the three flag sections, the trial reference, and their anchors and fixtures. The reviewer then decides whether the loop moves become unconditional phase text or the loop reference is removed. That decision is outside this spec.

### Budgets and anchors

`projects/benchkit.md` raises three prose budgets in the tickets that grow the files. The rows move `bench-implement-spec.md` from 80 to 84, `bench-write-spec.md` from 73 to 77, and `craft-tdd` from 122 to 124. The four-line section shape above is what makes a raise of four exact. `bench-debug.md` and `craft-delegate` stay inside their rows.

A new registry file, `internal/anchors/registry_debug_loop.go`, holds the new anchors, and the ordered registry appends it. The canary family registry lists the new file as an owner source. Each new Require or Forbid anchor has one fixture under `tests/canary/workflow-guidance-anchors`.

## Implementation chunks

The three tickets land in order, because tickets 2 and 3 both extend the registry file, the changelog heading, and the fixture family.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| DL-C1 / `1-write-the-loop-reference.md` | The loop discipline has one source, and the debug phase cites it | DL1, DL2, DL3, DL4, DL5, DL6, DL7, DL38, DL39 | `docs-currency-workflow`, fixture bite, prose budgets | no |
| DL-C2 / `2-keep-the-fix-in-the-loop-session.md` | The debug session writes its fix, and a delegate runs the loop to its fence | DL8, DL9, DL10, DL11, DL12, DL13 | `docs-currency-workflow`, fixture bite | no |
| DL-C3 / `3-add-the-debug-trial-flag.md` | The flag selects an arm, and both arms record for the compare | DL14 to DL37 | `docs-currency-workflow`, fixture bite, prose budgets, assessment compare dogfood | yes |

```bench-completion-plan
{"version":1,"chunks":[{"id":"DL-C1","tickets":["1-write-the-loop-reference.md"],"verification":[{"id":"anchors","command":"bench test --check docs-currency-workflow","probe":"omit the signal-before-theory sentence from loop.md"},{"id":"bite","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"}]},{"id":"DL-C2","tickets":["2-keep-the-fix-in-the-loop-session.md"],"verification":[{"id":"anchors","command":"bench test --check docs-currency-workflow","probe":"omit the loop-session-writes-the-fix sentence from bench-debug.md"},{"id":"bite","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner"}]},{"id":"DL-C3","tickets":["3-add-the-debug-trial-flag.md"],"verification":[{"id":"anchors","command":"bench test --check docs-currency-workflow","probe":"omit the record-ids sentence from loop-trial.md"},{"id":"bite","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"}]}],"final_verification":[{"id":"acceptance","command":"bench test"},{"id":"integration","command":"bench test --check system"}]}
```

## Testing decisions

- A good test drops one rule sentence from its file and observes the `docs-currency-workflow` check go red with that rule's diagnostic.
- The anchor registry and the workflow-guidance fixture family are the seams. Each fixture's `MUTATE.json` is the drop, and its `EXPECT` is the diagnostic. The prior art is `tests/canary/workflow-guidance-anchors/write-spec-conversation-fork`.
- The `guidance-prose-budgets` check observes each budgeted file against `projects/benchkit.md`.
- The record contract has no new parser. The DL-C3 review proves the plan template and two synthetic runs through `bench assessment compare`.
- The gate observes the feature through `docs-currency-workflow`, the fixture-bite test, and `guidance-prose-budgets`, all inside the ordinary `test` phase.

### Seam diagram

    trigger: a reflow, a rewrite, or a deletion in a phase file or a craft reference
        │
        ▼
    edited Markdown  ──▶  [ anchor registry + docs-currency-workflow ]  ──▶  diagnostic or green
                              ◀ tests attach here: a fixture drops the sentence; the owner check must red

    trigger: the reviewer copies the plan template and imports two runs
        │
        ▼
    plan.json + runs  ──▶  [ bench assessment compare ]  ──▶  comparison report
                              ◀ review attaches here: the report binds both runs to their conditions

### Acceptance coverage map

Every anchor row below runs through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`), and each names the canary that drops its sentence.

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| DL1 | 1 | `loop.md` states move 1: one command, run once and pasted, before any theory | anchor, canary `loop-signal-before-theory`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a dropped move 1 reopens theory-first work |
| DL2 | 1 | `loop.md` states move 2: at most three ranked written candidates, each with a prediction, and act on the first | anchor, canary `loop-bounded-candidates`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | an uncapped list is the deliberation the trial measures |
| DL3 | 1 | `loop.md` states move 3: one variable per run | anchor, canary `loop-one-variable`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | two variables per run hide which one moved the signal |
| DL4 | 1 | `loop.md` states move 4: close on the command, not on judgment | anchor, canary `loop-close-on-command`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a judgment close is the agent grading its own work |
| DL5 | 1, 2 | `loop.md` carries the four-property checkbox list and `bench-debug.md` carries none | the moved `- [ ] **red-capable**` anchor targets `loop.md`, a Forbid anchor on `bench-debug.md`, canaries `loop-properties-checklist` and `debug-checklist-copy-forbidden`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a second copy or a lost checklist reds the check |
| DL6 | 2 | `bench-debug.md` Phase 1 names `references/loop.md` as the owner of the properties | anchor, canary `debug-phase1-loop-pointer`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a dropped pointer leaves Phase 1 without its completion rule |
| DL7 | 3 | `CONTEXT.md` defines **loop** | anchor, canary `context-loop-term`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a dropped term lets the phases drift to synonyms |
| DL8 | 4 | `bench-debug.md` no longer routes code authorship through `craft-delegate` | Forbid anchor, canary `debug-fix-route-retired`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a reintroduced route reds |
| DL9 | 4, 5 | `bench-debug.md` states that the session that owns the loop writes the fix and delegates only a read-only search | anchor, canary `debug-loop-session-writes-fix`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a dropped sentence reopens the fresh-delegate route |
| DL10 | 5 | `bench-debug.md` carries no scoped-fix delegation clause | Forbid anchor on `a scoped fix`, canary `debug-scoped-fix-retired`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | the surviving clause would still authorize a delegated fix |
| DL11 | 6 | `bench-debug.md` permits a write delegate to run the phase through Phase 6 on a bug inside its fence | anchor, canary `debug-delegate-in-fence`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a dropped permission returns the delegate to a description-only report |
| DL12 | 7, 8 | `bench-debug.md` limits a write delegate to Phases 1 to 3 outside its fence and puts the loop command in the blocked report | anchor, canary `debug-delegate-out-of-fence`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a dropped limit lets a delegate fix out of fence |
| DL13 | 7 | `craft-delegate` names the three debug phases in its blocked-delegate rule | anchor, canary `delegate-blocked-runs-loop`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a delegate that reads only the skill would stop before the loop exists |
| DL14 | 9 | each build phase carries one `--debug-trial` section with the `on` and `off` grammar that points at `loop-trial.md` | three anchors, canaries `implement-spec-debug-trial-pointer`, `write-spec-debug-trial-pointer`, and `review-debug-trial-pointer`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a phase without the section cannot enter an arm |
| DL15 | 9 | an unrecognized or missing flag value stops the phase with the grammar | review-owned: the grammar sentence in each phase section | a phase that guesses an arm mixes the arms |
| DL16 | 17 | no build phase file names `references/loop.md`, and a restated rule without the path is review-owned | three Forbid anchors, canaries `implement-spec-loop-rule-leak`, `write-spec-loop-rule-leak`, and `review-loop-rule-leak`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a restated rule contaminates the off arm |
| DL17 | 10 | `loop-trial.md` states that the implement-spec arm runs the red signal once and pastes it before the first edit | anchor, canary `loop-trial-implement-red-first`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | without the paste the first action is a theory |
| DL18 | 11 | `loop-trial.md` states that an `already covered` or `not TDD-able` classification stands in for a red | anchor, canary `loop-trial-implement-classification`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | without it the arm forces a false red |
| DL19 | 10 | `loop-trial.md` states that the implement-spec arm makes one edit per run at the ticket's seam after a green baseline | anchor, canary `loop-trial-implement-one-edit`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a batch of edits per run hides which one moved the signal |
| DL20 | 10 | `loop-trial.md` states that a seam that refuses the behavior exits as a wrong spec | anchor, canary `loop-trial-implement-wrong-spec`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | without it the agent searches for a second seam |
| DL21 | 12 | `loop-trial.md` states that the write-spec arm caps the seam sketch at three ranked candidates per story group | anchor, canary `loop-trial-write-spec-cap`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | an uncapped sketch is the deliberation the trial measures |
| DL22 | 13 | `loop-trial.md` states that each candidate names its red check and a candidate with no red is dropped | anchor, canary `loop-trial-write-spec-red-check`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a candidate with no red is a vibe |
| DL23 | 14 | `loop-trial.md` states that the approval table shows the chosen seam and its ranking | anchor, canary `loop-trial-write-spec-ranking`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a hidden ranking cannot be re-ranked |
| DL24 | 15 | `loop-trial.md` states that a strong finding carries the invocation and output digest of its refutation run | anchor, canary `loop-trial-review-run`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | an unrun finding is an opinion |
| DL25 | 15 | `loop-trial.md` states that a finding with no runnable check takes `ask-user` at most | anchor, canary `loop-trial-review-ask-user`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | an unrunnable finding routed as `auto-fix` repairs a vibe |
| DL26 | 16 | `loop-trial.md` states that the off arm runs today's phase and adds only the trial provenance | anchor, canary `loop-trial-off-arm`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | an off arm with loop rules measures nothing |
| DL27 | 18, 21 | `loop-trial.md` names the plan id `debug-loop-trial`, the conditions `loop-on` and `loop-off`, and the capability key `debug-loop` | anchor, canary `loop-trial-record-ids`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a run under another plan or condition id fails the compare's binding |
| DL28 | 18, 21 | `loop-trial.md` states that both conditions pin one kit revision and each run's source is that revision | anchor, canary `loop-trial-record-revision`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a run on another revision is a hard compare error |
| DL29 | 18 | `loop-trial.md` states that each run names its condition, task, repetition, and holdout | anchor, canary `loop-trial-record-run-fields`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a missing repetition or holdout is a reported limit on every run |
| DL30 | 18 | `loop-trial.md` states that the plan pins every recorded role's model and effort and the `implementation` line in every condition | anchor, canary `loop-trial-record-lines`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | an unpinned role is a hard compare error and a missing implementation line fails validation |
| DL31 | 18 | `loop-trial.md` states that each condition pins the harness, the limits, the acceptance list, and the review-axes list, and each run's trial object repeats them | anchor, canary `loop-trial-record-conditions`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a harness or limits mismatch is a hard compare error |
| DL32 | 19 | `loop-trial.md` names the per-phase quality measures | anchor, canary `loop-trial-measures`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | an unnamed measure cannot enter a report |
| DL33 | 19 | `loop-trial.md` states that the tolerance carries `max_failure_rate` and its measures map stays empty unless every task records the measure | anchor, canary `loop-trial-tolerance`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a named measure one phase never records marks every run of that phase unknown |
| DL34 | 20 | `loop-trial.md` states that a counter the harness does not expose stays unknown | anchor, canary `loop-trial-unknown-stays-unknown`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a zero would win an arm |
| DL35 | 18, 21 | the plan template validates and two synthetic runs on one revision compare, with a first limit line that starts `descriptive evidence only` | review-owned: `bench assessment compare --plan <template copy> --runs <two synthetic ids>` at the DL-C3 review | a template the verb refuses records nothing |
| DL36 | 22 | `loop-trial.md` states that one landing removes the flag sections, the trial reference, and their anchors and fixtures after the compare | anchor, canary `loop-trial-exit-landing`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a flag with no exit becomes permanent |
| DL37 | 22 | `loop-trial.md` states that the reviewer then decides whether the loop moves become phase text or the loop reference is removed | anchor, canary `loop-trial-exit-decision`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | without it the exit lands the rule by default |
| DL38 | 23 | every edited budgeted file stays within its `projects/benchkit.md` row | `guidance-prose-budgets` check over the live tree, through `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | an over-budget file reds the lane |
| DL39 | 24 | every new fixture bites its owner and restores | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a fixture that does not red proves nothing |

### Edge inventory

- Won't handle: an unterminated HTML comment in an edited file — the anchor reader truncates at it today, and the fixture-bite test still grades it.
- Won't handle: a symbolic link where a reference belongs — the anchor reader refuses a link, and the `craft-tdd` references directory stays real.
- Won't handle: a plan file with control bytes or unknown fields — `bench assessment` refuses it today, and the reviewer authors the plan locally.
- Won't handle: a Claude-native token counter — no producer exists, so those runs record `unknown`, and the Codex fragment path still records.
- Won't handle: `--debug-trial` on `/bench-debug`, `/bench-final-check`, or `/bench-drain` — the debug phase is the source, not an arm, and the other phases keep today's text.
- Won't handle: an attempt role the plan does not pin — the compare refuses it today, and DL30 tells the reviewer to pin every recorded role.
- Won't handle: an `off`-arm agent that applies an `on` rule — the reference opens with the off arm, and the chunk review grades the rest.
- Won't handle: a restated loop rule that omits the path — DL16 marks that residual review-owned, and the Forbid anchor still catches the path.

## Ownership fences

- `.agents/skills/bench-craft-tdd/SKILL.md`
- `.agents/skills/bench-craft-tdd/references/loop.md`
- `.agents/skills/bench-craft-tdd/references/loop-trial.md`
- `.agents/commands/bench-debug.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `CONTEXT.md`
- `projects/benchkit.md`
- `CHANGELOG.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/conformance/registry_test.go`
- `tests/canary/workflow-guidance-anchors`
- `reviews/debug-loop-trial.md`

The ticket grammar co-names each fixture that pins an edited file and each registry bound to the anchors package. These paths enter the fence for that reason, and no ticket plans an edit there:

- `tests/canary/skills-index-command-adapters`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/skill-description-budgets`
- `tests/canary/claude-agent-definitions`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`

Reviewer disposition: Proposed fences await the spec-and-ticket approval.
No implementation fence authorizes a write during this spec phase.

## Out of scope

- Won't handle: the exit landing — 8 edits, 1 gate run. It removes the flag sections, the trial reference, and their anchors after the compare.
- Won't handle: an ADR for the loop as an unconditional phase rule — 1 edit, 1 gate run. It waits for the exit decision.
- Won't handle: a Claude-native token producer for the assessment harness mapper — 4 edits, 2 gate runs. It is a Go seam with its own spec.
- Won't handle: a field-guide change — 1 edit, 1 gate run. The `/bench-debug` card describes no delegate route, so it stays true.

These exclusions grant no future implementation authority.

## Further notes

The trial's plan purpose is `descriptive`, because the compare requires exactly three named conditions for a kit-causal purpose. A two-arm plan reports descriptive limits by design.

The fixture count is the reviewer's dial. Each row names its canary, so a cut removes the row's anchor and fixture together.

The [enforcement read](assets/enforcement.md) records the readers this spec examined and the line counts at phase entry.
