# Bounded repair policy

Status: staged

Decision source: `specs/bounded-repair-policy/decisions/ft232-repair-loop.md` (ready compiled map; policy outcome confirmed by ticket #14).

Verification log: 1 iteration(s) to accept — Sol/high spec findings are folded. The single Astra/high slice review follows slicing.

## Problem

An implementation can keep repairing review findings while each new review introduces more work.
Useful progress does not bound the total repair cost.
Optional improvements can also become blockers, even when required behavior is complete.

## Solution

Give each implementation chunk a default allowance of two repair cycles after its initial review.
A fresh review does not reset that allowance, and progress does not extend it.
Only an explicit reviewer decision extends the allowance.

Required checks, approved acceptance failures, and concrete correctness or safety defects remain blocking.
Optional improvements remain advisory.
At the allowance limit, the author reports remaining blockers and waits for the reviewer.
Pre-review implementation retains its existing progress rules.

This spec delivers workflow guidance and conformance protection for that guidance.
It does not implement a runtime counter or infer a repair loop from telemetry.
The separately approved collection pilot remains future work under the same compiled decision source.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: guidance affects every linked repository, and semantic classification requires review beyond the existing anchor checks.
Harder chunks: BP-C1.

### Bound review repair work

1. As a reviewer, I want a fixed default repair allowance, so that useful progress cannot prolong a review repair loop indefinitely.
2. As an author, I want a repair cycle to include its verification, so that individual tool calls do not consume the allowance.
3. As a reviewer, I want a fresh review to retain the count, so that repeated reviews cannot reset the allowance.
4. As a reviewer, I want sole authority to extend the allowance, so that the author cannot authorize its own additional cycles.
5. As an author, I want the allowance in every implementation mode, so that unattended and delegated work has the same stopping rule.
6. As an author, I want light-path work to count as one chunk, so that a missing spec does not remove the allowance.

### Preserve completion requirements

7. As a reviewer, I want failed required checks to remain blocking, so that the allowance cannot weaken the gate.
8. As a reviewer, I want unmet acceptance requirements to remain blocking, so that the author cannot land a partial implementation.
9. As a reviewer, I want concrete correctness and safety defects to remain blocking, so that a green gate cannot conceal known defects.
10. As an author, I want optional improvements to remain advisory, so that suggestions do not require endless repairs.
11. As a reviewer, I want a bounded stop with remaining blockers, so that I can decide the next action.
12. As an author, I want completion when required work is verified, so that the allowance does not require unnecessary repair cycles.

### Keep existing workflow guarantees

13. As an author, I want existing pre-review progress rules, so that the new allowance does not cap initial implementation work.
14. As a reviewer, I want current review and verification evidence, so that a repair limit cannot waive a completion checkpoint.
15. As an author, I want one policy source and retained repair state, so that a resumed session applies the same allowance.
16. As a reviewer, I want policy conformance checks to detect omissions, so that a later prose edit cannot silently remove the allowance.
17. As a reviewer, I want the collection pilot kept separate, so that this policy does not authorize a detector or new data collection.
18. As a reviewer, I want mandatory standards to remain blocking without automated checks, so that semantic requirements cannot become optional preferences.
19. As a reviewer, I want optional advice retained separately, so that suggestions stay visible without inflating finding totals.

## Implementation decisions

### Policy owner and scope

The line skill owns the allowance through a focused reference, `references/bounded-repair-policy.md`.
The line skill links and charges that reference when implementation reaches its initial review.
The implementation phase, review phase, and review craft skill consume that same owner.
They do not restate the numeric limit or independently define a repair cycle.
The existing glossary remains the definition source for a repair cycle.

The reference applies to retained, full, delegated, unattended, and light-path implementation work.
A light-path run counts as one chunk when it receives review findings.
This change introduces no new mandatory review phase for light-path work.
Initial implementation, its pre-review checks, and the first review do not consume the repair allowance.
A repair attempt followed by verification is one cycle, including when it addresses several findings.

The fixed allowance takes precedence over continuation while progress holds after initial review.
The existing no-progress, cancellation, user-budget, and required-decision stops can still stop work earlier.
An explicitly uncapped implementation line does not implicitly remove the post-review allowance.
A fresh review does not start a new allowance for the same chunk.
Only an explicit reviewer extension permits further cycles after exhaustion.
These instructions implement BP1 through BP7 and BP14 below.

### Classification and current review records

Classify each proposed repair against its cited requirement or concrete defect evidence before accepting it as blocking work.
A style preference with no binding requirement and no concrete defect remains optional advice.
A documented mandatory standard is a requirement, even when its subject is prose.
The cap does not downgrade a required rule to advice.

Optional advice is not a finding or a repair target.
It receives no repair-routing disposition or finding ID and does not enter finding totals.
Retain it in a separate advice section of the native excerpt and the review pickup.

Preserve the existing `no-op`, `auto-fix`, and `ask-user` repair-routing meanings for actionable findings.
Do not invent a fourth routing disposition or label an unrefuted suggestion as `no-op` merely to obtain a pass.
The current native review result may pass with no blocking findings while its accompanying prose retains optional advice.
Earlier native findings and their supersession history remain intact.
No record schema or validator change is required.

Current results or permitted native reaffirmations remain required after repairs.
Repeat independent review only for a later semantic change or a named cross-chunk concern that invalidates earlier evidence.
Such a review does not reset the repair count.
The existing checks still reject unresolved findings, stale source identity, missing axes, and incomplete verification.
These instructions implement BP8 through BP13, BP16, BP20, and BP21 below.

### State and handoff

Retain the consumed cycle count and any explicit reviewer extension in the existing review pickup.
For light-path work without a pickup, retain that state in the existing session handoff.
Use ordinary prose in those artifacts, without a new machine-readable schema or CLI command.
A resumed session reads that state before authorizing another repair cycle.
It reconstructs a missing count from available evidence instead of assuming that the allowance is unused.
If the evidence cannot establish the remaining allowance, return that uncertainty to the reviewer before another repair cycle.

At exhaustion with blockers, report the consumed allowance, completed repairs, remaining blockers, and the needed reviewer decision.
Do not start another repair cycle or advance the dependent chunk while that decision is pending.
If blockers close within the allowance, continue through the existing verification and landing requirements.
The allowance is a maximum, not a target.
These instructions implement BP3, BP12, BP13, and BP15.

### Conformance attachment

Extend the existing implementation-continuation anchor family and its existing conformance test file.
Keep one independently authored expectation for each new policy obligation only where an omission probe demonstrates why independence is necessary.
Record the red from removing the production anchor while retaining its independent test expectation.
Do not add a second production limit constant, policy parser, fixture harness, or count registry.

Preserve the existing review-convergence anchors and tests.
Clarify that progression rules about findings concern unresolved blockers, with optional advice recorded separately.
Do not remove or weaken their requirement for current repair coverage.

The implementation must stay within current guidance and structure budgets.
At the reviewed tip, craft-line uses 128 of 130 physical lines, craft-review uses 122 of 122, and bench-implement-spec uses 79 of 80.
Consumer edits replace or condense existing prose rather than increasing a budget.
The ownership fence grants no budget change.
Use the focused policy reference to avoid expanding the line skill into a second long policy document.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| BP-C1 / pending Sol slicing | All implementation modes apply one bounded repair policy with current completion evidence | BP1 through BP21 | docs-currency-workflow, TestImplementationContinuation, review-owned scenario replay | yes |

The slicer confirms ticket membership and ownership fences after the spec review.
A split cannot land a numeric allowance without its blocker classification and completion safeguards.
Each chunk receives its normal review checkpoint before a successor starts.

## Testing decisions

The existing workflow guidance seam is the product surface.
The gate mechanically checks that required instructions and their consumer links remain present.
Independent scenario review checks whether those instructions produce the approved decisions.
A substring check does not prove that a live agent obeys the policy.

Use `bench test --check docs-currency-workflow` as the live conformance root.
Its route is `checkDocsCurrencyAndWorkflow` to `checkWorkflowAnchors` to `anchors.EvaluateGroup`.
The existing `AfterImplementSpec` group already contains the implementation-continuation family.
Use `TestImplementationContinuation` and `runAnchorBites` as the omission-test precedent.
Reuse that fixture harness rather than adding a parallel one.

Preserve `TestReviewConvergenceContractCurrentDocs` and `TestReviewCheckpointFindingAndReviewIdentity` as existing guard regressions.
Run their focused checks before and after the guidance change over their existing deterministic input families.
The completion record schema and its production validators remain unchanged.
The scenario review must separately check optional-advice classification because those validators cannot judge prose semantics.

### Seam diagram

```text
initial review -> line-owned policy -> repair, verify, continue or reviewer handoff
                         |
                         +-> existing pickup or handoff retains consumed allowance
                         |
                         +-> anchor checks detect missing instructions
                         +-> scenario review checks the required decision
```

### Acceptance coverage map

Each workflow row has a named mechanical attachment and a review-owned scenario.
The implementation adds its policy anchors to `TestImplementationContinuation` and the existing `docs-currency-workflow` root.

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| BP1 | 1 | Guidance stops a third repair cycle after two cycles made progress but left a blocker | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | Removing the fixed limit or allowing progress to extend it fails the anchor and scenario |
| BP2 | 2 | Guidance counts one attempt with verification as one cycle | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | Counting three tool calls as three cycles contradicts the cycle definition |
| BP3 | 3, 15 | A fresh review preserves the same chunk's consumed allowance | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | Restarting the review after two cycles cannot authorize a third cycle |
| BP4 | 4 | An explicit reviewer extension permits the specified additional repair work | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | Ignoring an explicit extension incorrectly stops authorized work |
| BP5 | 4 | An author cannot extend its own exhausted allowance | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | An author declaration of uncapped work cannot substitute for reviewer authority |
| BP6 | 5 | The policy applies to retained, full, delegated, and unattended implementation runs | review-owned: mode census with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | Omitting any enumerated entry route leaves an approved implementation mode outside the policy |
| BP7 | 6 | Light-path review repairs use one chunk's allowance | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | The absence of a spec cannot create an unlimited repair path |
| BP8 | 7 | An unresolved required check remains blocking at exhaustion | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | A red required check cannot become optional because the allowance is spent |
| BP9 | 8 | An unresolved approved acceptance failure remains blocking at exhaustion | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | A partial implementation cannot become complete because its repair allowance is spent |
| BP10 | 9 | A concrete correctness defect remains blocking despite a green gate | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | A green gate cannot dismiss cited evidence of incorrect behavior |
| BP11 | 9 | A concrete safety defect remains blocking despite a green gate | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | A green gate cannot dismiss cited evidence of unsafe behavior |
| BP12 | 11 | Exhaustion with blockers produces a reviewer handoff before dependent work advances | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | A third repair attempt or successor chunk violates the bounded stop |
| BP13 | 10, 12 | A native review containing only optional advice passes with no finding IDs | review-owned: scenario decision with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | A cosmetic preference cannot produce a blocking finding ID or require another repair cycle |
| BP14 | 13 | Pre-review implementation retains its existing continuation rules | `internal/conformance/implementation_continuation_test.go` (`TestImplementationContinuation`), review-owned: scenario decision; evidence in `reviews/bounded-repair-policy.md` | Applying the new numeric limit to initial implementation contradicts the approved scope |
| BP15 | 15 | A resumed author retains the consumed allowance and reviewer extension evidence | review-owned: resume scenario with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | Treating missing or resumed state as a fresh allowance bypasses the same-chunk cap |
| BP16 | 14 | The allowance does not waive current review and verification evidence | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointFindingAndReviewIdentity`), review-owned: scenario decision; evidence in `reviews/bounded-repair-policy.md` | Unresolved findings and stale review identity still fail the existing checkpoint |
| BP17 | 15 | Consumer guidance points to one policy owner without copying the numeric limit | review-owned: source census with docs-currency-workflow instruction checks; evidence in `reviews/bounded-repair-policy.md` | A missing pointer or a second policy definition creates drift between implementation modes |
| BP18 | 16 | Removing a required production policy anchor fails an independent expectation | `internal/conformance/implementation_continuation_test.go` (`TestImplementationContinuation`) | Removing both the instruction and its registry entry cannot define its own green result |
| BP19 | 17 | The policy change introduces neither pilot collection nor detector warnings | review-owned: scope and changed-path inspection; evidence in `reviews/bounded-repair-policy.md` | A telemetry producer or advisory renderer change exceeds this specification |
| BP20 | 18 | A documented mandatory standard violation remains blocking without an automated check | review-owned: mandatory-standard scenario; evidence in `reviews/bounded-repair-policy.md` | Treating a cited required prose rule as a preference would permit a false pass |
| BP21 | 19 | Optional advice is retained outside the finding and repair-target census | review-owned: native excerpt and pickup inspection; evidence in `reviews/bounded-repair-policy.md` | Dropping advice loses its record, while counting it as a finding invents a blocker |

### Edge inventory

The audience is every repository that receives the linked workflow, including this kit.
The mode census is retained, full, delegated, unattended, and light-path implementation.
The classification census includes required checks, acceptance failures, correctness defects, safety defects, mandatory standards, and optional advice.
Each mode and classification attaches to BP6 through BP11, BP13, BP20, and BP21.

A cycle that addresses several findings still consumes one cycle: BP2.
An unchanged rerun without a repair attempt is verification, not another completed repair cycle: BP2.
A reviewer extension and a self-issued extension have separate outcomes: BP4 and BP5.
A new reviewer result on the same chunk retains its count: BP3.

Missing or ambiguous count evidence does not mean zero: BP15.
Earlier user-budget and cancellation stops still apply: BP14 and BP16.
Mandatory prose standards remain requirements even without a check: BP20.
Contrast that case with an otherwise similar nonbinding preference: BP13.
A genuine optional improvement remains advisory even if a reviewer mentions it: BP13.

Won't handle: a new repair-counter file parser — the existing review pickup and handoff remain the callers.
Won't handle: new hostile path, network, or subprocess inputs — existing Bench verbs remain the callers without grammar changes.
Won't handle: changed treatment of absent or empty project directories — the existing conformance readers retain their file-state contracts.
Won't handle: a mechanical proof that an agent obeys instructions — scenario review grades the policy, and current gate checks grade the tree.

## Ownership fences

These provisional fences describe the complete behavior change.
Sol confirms them against the ticket union during slicing.

- `.agents/skills/bench-craft-line/SKILL.md`
- `.agents/skills/bench-craft-line/references/bounded-repair-policy.md`
- `.agents/skills/bench-craft-review/SKILL.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/implementation_continuation_test.go`
- `CHANGELOG.md`
- `reviews/bounded-repair-policy.md`
- `specs/bounded-repair-policy/spec.md`

The implementation owns only the named product files and its ordinary verification artifacts.
The compiled map remains settled provenance.
No build-wide rewrite may change unrelated specs or their tickets.

## Out of scope

The separately approved collection pilot and its report require another specification: estimated 8 file edits, 1 landing gate run.
The estimate covers producer attribution, evidence capture, report output, tests, and documentation at existing owners.
The pilot spec must refine that estimate after selecting its seams.

A detector and warning surface require a later reviewer decision after the pilot: estimated 4 file edits, 1 landing gate run.
That estimate assumes reuse of pilot evidence and an existing output surface.
Neither estimate authorizes those capabilities.

This change does not add a CLI flag, modify the native review schema, or weaken a completion validator.
It does not alter the pre-review continuation policy, default implementation authorship, or model routing.
FT232 remains on the roadmap while its separately approved collection work remains open.

## Further notes

### Source clauses and readers

| Source occurrence | Exact source clause | Coverage |
| --- | --- | --- |
| Ticket #1, Answer | Shape both outcomes under this map and specify them separately. | BP19 |
| Ticket #2, Answer | Required checks, approved acceptance failures, and concrete correctness or safety defects remain blocking. | BP8, BP9, BP10, BP11, BP20 |
| Ticket #2, Answer | Optional improvements remain advisory. | BP13, BP21 |
| Ticket #2, Answer | At the allowance limit, the agent stops and reports unresolved blockers to the reviewer. | BP12 |
| Ticket #3, Answer | Each implementation chunk permits at most two repair cycles after its initial review. | BP1, BP2, BP14 |
| Ticket #3, Answer | Progress does not extend this allowance. | BP1 |
| Ticket #3, Answer | A fresh review does not reset the count. | BP3, BP15 |
| Ticket #3, Answer | At the cap, unresolved blockers return to the reviewer. | BP12 |
| Ticket #4, Answer | The default allowance applies to all implementation runs. | BP6 |
| Ticket #4, Answer | Light-path work counts as one implementation chunk. | BP7 |
| Ticket #4, Answer | Only an explicit reviewer decision extends the allowance. | BP4, BP5 |
| Ticket #7, Answer | The fixed allowance applies only to repairs after initial review. | BP14 |
| Ticket #7, Answer | Pre-review implementation keeps the existing progress and no-progress rules. | BP14 |
| Ticket #8, Answer | Required checks and review evidence still apply. | BP16 |
| Ticket #14, Answer | The first implements the bounded repair policy. | BP1 through BP18, BP20, BP21 |
| Ticket #14, Answer | A detector and any warning remain outside both specifications. | BP19 |

The ticket paths are under `decisions/ft232-repair-loop/tickets/` beside this spec.
These are the authoritative policy-clause occurrences; the compiled index only links their gists.
Tickets #5, #6, and #9 through #13 retain collection evidence and scope for the separate pilot specification.
The working agreement supplies the one-source and independent-expectation rules for BP17 and BP18.

Read the [enforcement evidence](assets/enforcement.md) for current source locations and the complete named reader sweep.
The map's structured sources were reread in this session.
The gate producer still records only the first lane failure, and the research report retains its bounded snapshot and unknowns.
Those facts justify the separation from the pilot, not a detector threshold.
No outside source was required for this policy specification.

### Pre-review proof checklist

- Cited symbols: the enforcement evidence resolves the named conformance and review-record functions.
- Import edges: no new import edge is proposed.
- Source-row clauses and occurrences: the source table and enforcement evidence enumerate the decision clauses and current readers.
- Promised field labels: no new machine-readable field label is promised.
- Changed-function callers: no product function or exported signature changes.
- Copy survival: BP17 compares the canonical reference with every named consumer.
- Restore and copy promises: none in the implementation scope.
- Bootstrap authority before execution: no new executable hop or authority claim.

### Flagged additions and implementation choices

The focused reference, ordinary-prose repair state, and optional-advice placement are implementation choices for the approved behavior.
They add no data-collection service, new denial surface, or native record field.
The mode census includes existing routes and does not introduce a new initial-review requirement for light-path work.
Missing count evidence returns to the reviewer because inventing a fresh allowance would violate the approved cap.
The scenario review remains explicit where a deterministic check cannot establish semantic compliance.
Its durable evidence belongs in `reviews/bounded-repair-policy.md`, not only in an ephemeral delegate return.

### Review and approval protocol

The reviewer directs one Sol/high spec review, then Sol/high ticket slicing, then one Astra/high slice review.
Accepted findings are folded without a second review round.
The Sol/high review found the optional-advice representation gap and missing mandatory-standard coverage.
BP13, BP20, and BP21 carry those repairs.
The exact source table, durable scenario evidence, and headroom constraints incorporate the nonblocking corrections.

The reviewer pre-approves this phase's landing.
This specification phase does not start the implementation or activate the collection pilot.
