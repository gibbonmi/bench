# Shared preparation for Bench delegates

Status: staged

Decision source: ready compiled map `specs/shared-delegate-startup/decisions/shared-delegate-startup.md`, confirmed 2026-09-16.

Verification log: 2 iteration(s) to accept — one Sol/high review found two blockers. One author repair pass closed both under coordinator verification.

## Problem

Delegates can repeat the same preparation before they start distinct assignments.
A full conversation fork can preserve useful preparation, but it also inherits irrelevant history and its parent's line.
Review forks can inherit conclusions that impair independent judgment.
Fewer reads alone do not establish lower total cost.

## Solution

Offer one optional workflow across spec authorship, delegated implementation, and review.
Use the existing approved-context spec fork.
Prepare common implementation sources in a fresh session before ticket assignments.
Prepare neutral review sources before each axis derives its findings.
Choose a suitable native fork or the phase's permitted fallback.
Keep current authorization, source verification, isolated worktrees, and retained repair authors.

This change edits guidance only.
It adds no command, flag, collector, assessment schema, or default route.
The workflow can operate without a paid comparison.
It makes no promise about cache hits, billed savings, or latency.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: Review neutrality is the hardest outcome. Decisions are precise, existing seams are known, and semantic guidance requires independent review.
Harder chunks: SP1.

### Activation and authority

1. As a reviewer, I want one optional instruction, so that an approved delegated run and its reviews share preparation.
2. As a reviewer, I want standalone phase activation, so that I can use the workflow without a full run.
3. As a reviewer, I want existing authorization checks, so that shared preparation cannot approve a build.
4. As a coordinator, I want unchanged defaults, so that an absent instruction preserves the normal route.

### Native context and line

5. As a coordinator, I want fork eligibility checks, so that inherited context matches the assignment's requirements.
6. As a coordinator, I want honest line inheritance, so that a model override cannot disguise the actual fork line.
7. As a reviewer, I want a complete cost estimate, so that a different line requires a lower estimated total.
8. As a reviewer, I want unknown costs preserved, so that missing evidence cannot justify a line change.
9. As a coordinator, I want native parent checks, so that unsupported nested forks use the permitted route.
10. As a coordinator, I want phase-specific fallback, so that an unavailable fork cannot relax phase authority.

### Spec authorship

11. As a spec author, I want the existing approved-context fork, so that shared preparation preserves the authorized decision source.
12. As a reviewer, I want one normal spec author, so that the optional workflow adds no preparation session or authorship split.

### Delegated implementation

13. As a ticket author, I want fresh common preparation, so that unrelated planning history does not enter my assignment.
14. As a ticket author, I want my assignment after common preparation, so that sibling preparation remains reusable.
15. As a coordinator, I want current source checks, so that inherited material cannot authorize work against stale bytes.
16. As a coordinator, I want refreshed later batches, so that accepted predecessor changes reach dependent authors.
17. As a coordinator, I want retained repair authors, so that shared preparation cannot replace an existing author for convenience.
18. As a ticket author, I want my own worktree identity, so that inherited paths cannot grant write authority.

### Independent review

19. As a review axis, I want the frozen diff and primary sources, so that common preparation supplies relevant raw material.
20. As a review axis, I want no inherited author rationale, so that implementation reasoning cannot seed my conclusions.
21. As a review axis, I want no inherited reviewer findings, so that another axis cannot seed my conclusions.
22. As a review axis, I want independent source checks, so that common preparation cannot replace my derivation.
23. As a coordinator, I want existing review evidence, so that the optional route does not collect the same frozen pair again.
24. As a coordinator, I want native review authority preserved, so that a missing fork cannot authorize an inline or CLI review.

### Evidence and maintenance

25. As a reviewer, I want separate estimates and actual charges, so that projected savings cannot appear as measured savings.
26. As a reviewer, I want parent and child attribution, so that inherited counters cannot count preparation twice.
27. As a reviewer, I want unknown telemetry preserved, so that unavailable child records do not become measured zero.
28. As a maintainer, I want one owner per rule, so that shared guidance remains consistent within existing budgets.
29. As a maintainer, I want existing validation intact, so that the optional workflow cannot weaken verification.
30. As a reviewer, I want compatibility claims tied to observed evidence, so that documentation does not imply an untested live capability.

## Implementation decisions

### Owners and terms

Fold the optional procedure into `craft-delegate` through one focused reference.
Use `references/shared-preparation.md` for activation, eligibility, common preparation, and fallback.
Keep line selection in `craft-line` and cost evidence in the existing assessment guidance.
Each phase links to the shared owner and states only its phase-specific duties.
Keep the decision source as the authority for scope.

Add these terms to the glossary during implementation:

| term | meaning | avoid |
| --- | --- | --- |
| shared preparation | Common source material prepared before distinct delegate assignments. | shared findings, cached context |
| native fork | A native delegate that inherits the parent conversation under the active harness contract. | worktree copy, guaranteed cache hit |
| fresh delegate | A delegate whose task inputs arrive through its charge rather than inherited conversation history. | context-free delegate |
| suitable fork | A native fork that preserves required context, the selected line, and assignment isolation. | any available fork |

`worktree`, `line`, and `seam` keep their current glossary meanings.
These terms describe guidance, not new serialized fields.

### Activation and selection

One explicit reviewer instruction activates shared preparation for an approved delegated run and its reviews (SP1).
A standalone phase accepts the same instruction (SP2).
The instruction grants no delegated-build approval or ticket-graph change (SP3, SP4).
Without that instruction, keep the existing route (SP5).

Check the active native surface before each batch (SP6).
A compiled harness record cannot prove runtime availability (SP7).
A full fork inherits model and effort (SP8).
Omit model and effort overrides when the active full-fork contract rejects them.
Declare the inherited line in the charge without claiming unavailable runtime metadata.

A different line requires a documented lower total-cost estimate (SP9).
The estimate includes preparation, all delegates, output, and expected repairs (SP10).
Applicable cache-write or other charges belong in that total.
Unknown required inputs preserve the existing line (SP11).
An equal or higher estimated total also preserves it (SP12).
This permission applies before eligible dispatch; it does not replace retained authors during repairs (SP23).

A different parent line is usable only through a native route that actually supports it (SP8).
No label in a charge changes the line of an inherited fork.
Check parent eligibility before designing a nested preparation tree (SP13).
Claude's documented no-nested-fork limit means a forked preparation helper cannot be the portable parent for sibling forks.
Use an eligible parent or the phase's fallback.

If the fork is unsuitable, builds and reviews use fresh native delegates with explicit prepared evidence (SP14, SP15).
Report the fallback (SP16).
If native review delegation itself is unavailable, preserve the existing capable-harness handoff (SP17).
The fresh fallback never substitutes a same-family CLI or coordinator review.

### Spec author

Preserve the approved-source author fork and its inherited line (SP18).
Apply common cost accounting to that existing fork (SP10, SP35).
Add no preparation session or split authorship to the normal spec phase (SP19).
When the mandatory spec fork is unavailable, hand off to a capable session (SP20).
Do not replace the authorized source with an unreviewed summary.

The user's two-author sequence is an explicit experiment for this specification phase.
It does not change the approved product behavior in SP19.
The coordinator first forks a spec author, then forks a ticket author after the spec return.
The user selects one Sol/high review pass over the completed pair.

### Implementation authors

Use a fresh implementation session for common source preparation (SP21).
Keep ticket-specific assignments after that common material (SP22).
The current build preflight remains the mechanical source owner (SP24).
Each author verifies its assignment, source tip, and required bytes before action (SP25).
Regenerate stale charges rather than trusting inherited evidence (SP26).

Refresh changed common material before a later batch (SP27).
Do not dispatch a dependent ticket before its existing prerequisite checkpoint (SP28).
Repairs remain with the recorded author (SP23).
Each child receives its own Bench assignment and source pin (SP29).
Conversation inheritance grants no permission to write the parent's checkout.

### Review axes

Use only neutral preparation for a review fork (SP30).
Include the frozen pair, raw diff, primary sources, and prepared evidence identities.
Exclude implementation reasoning and author rationale from inherited history (SP31).
Exclude other reviewers' findings from inherited history (SP32).
A coordinator with either kind of history cannot serve as that review parent.
A clean final message or a compaction summary does not establish neutral ancestry.

Create a neutral eligible parent when an authorized native route permits one.
Otherwise, use fresh delegates with explicit evidence (SP15).
Check the native nesting contract before choosing a preparation helper (SP13).
Do not deliver one axis's return into a parent that still needs to fork a sibling axis (SP32).

Each axis re-reads current sources and derives its own findings (SP33).
The preflight owner collects shared evidence once per preparation attempt (SP34).
Keep each axis's isolated read-only venue and the current exclusion of authors and coordinator from reviewer roles (SP40).
Keep existing terminal results, repair routing, chunk checkpoints, and final reconciliation (SP41).

### Cost and observed evidence

Keep projections apart from actual charges (SP35).
Use the existing assessment record and provenance rules (SP36).
Reference the preparation session and each child when evidence exists.
Verify child counter baselines before attributing cumulative totals (SP37).
Do not count inherited parent counters again as child preparation.

Unavailable usage, cache counters, or child records remain unknown (SP38).
Unsupported provider charges can use existing referenced other-cost entries.
A partial estimate does not satisfy the lower-total rule.
Fewer reads, a successful spawn, and prompt-cache documentation do not establish measured savings (SP39).
Record observed native context, line, and assignment evidence before claiming live compatibility (SP42).
No paid comparison is required or authorized by this specification.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| SP1 / `1-align-spec-preparation.md`, `2-prepare-implementation-delegates.md`, `3-preserve-neutral-review.md` | One optional workflow across spec authorship, delegated implementation, and review. | SP1–SP44 | Existing workflow, guidance-budget, and prose checks, plus the semantic cases below. | yes |

Use one chunk because activation, fallback, and phase consumers form one adoption decision.
Each ticket commits a complete guidance route and its evidence before the next ticket starts.
Keep these serial green commits on one integration source.
The shared reference and review pickup overlap intentionally; do not assign concurrent writers.
The last ticket completes combined adoption and reconciles the whole guidance invariant.
Review the integrated SP1 delta across Standards, Spec, and Coverage before completion.

| ticket | Blocked by | delivered outcome | owned rows |
| --- | --- | --- | --- |
| `1-align-spec-preparation.md` | none | Existing spec fork gains shared eligibility, line, and cost rules. | SP2, SP3, SP4, SP5, SP6, SP7, SP8, SP9, SP10, SP11, SP12, SP13, SP16, SP18, SP19, SP20, SP35, SP36, SP37, SP38, SP39, SP42 |
| `2-prepare-implementation-delegates.md` | `1-align-spec-preparation.md` | Approved implementation delegates receive fresh preparation and current assignments. | SP14, SP21, SP22, SP23, SP24, SP25, SP26, SP27, SP28, SP29 |
| `3-preserve-neutral-review.md` | `2-prepare-implementation-delegates.md` | Independent review receives neutral material and closes combined adoption. | SP1, SP15, SP17, SP30, SP31, SP32, SP33, SP34, SP40, SP41, SP43, SP44 |

Each row owner writes the guidance seam named in its ticket.
The first ticket owns spec activation, shared eligibility, spec authority, line selection, and cost evidence.
The second owns implementation duties and the build fallback.
The last owns combined activation, review duties, review fallback, and integrated ownership and preservation checks.

## Testing decisions

The changed interface is the guidance that a coordinator and each delegate read.
The acceptance rows below are review-owned unless they name a mechanical check.
A semantic red means an independent reviewer reports the exact counterexample as a blocker.
It is not a claim that a prose check proves native runtime behavior.

For each semantic row, compare the current guidance with the proposed guidance over its named case.
Record the old and new route or duty in the implementation verification evidence.
Preserved duties must give the same result in both readings.
New optional behavior must follow the approved source only when activated.

No new runtime seam, parser, or test registry is needed.
Do not add one anchor per semantic sentence.
Existing checks retain their present scope:

- `docs-currency-workflow`: existing authorization, freshness, native review, and retained-author anchors.
- `guidance-prose-budgets`: existing profile budgets for the edited skills and phases.
- `prose-mechanics`: sentence and paragraph mechanics.
- `decision-map-integrity`: compiled map structure and references during this spec phase.
- `bench coverage --check`: acceptance-row and story grammar during this spec phase.

The existing workflow-anchor and budget fixtures provide precedents for these mechanical seams.
No check or fixture changes are planned.
Keep anchored bytes intact by qualifying fresh-delegate rules and adding owner references.
If a necessary edit cannot preserve a pinned contract, report the exact conflict before expanding the implementation plan.

For prose-only implementation, use `craft-synthesis`'s proportional route.
Read every steered surface and obtain the commit's green verdict.
Use live native observations only when ordinary authorized work supplies them.
A missing live experiment limits compatibility claims and does not create an invented pass.

### Seam diagram

```text
explicit instruction + approved phase source
                    |
                    v
 existing phase -> shared preparation guidance -> native fork or phase fallback
                    ^
                    |
 independent review cases + existing guidance checks
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SP1 | 1 | One explicit instruction covers an approved delegated run and its reviews. | Review-owned: run activation case. | Requiring new activation at its review fails the case. |
| SP2 | 2 | A standalone phase accepts explicit shared preparation. | Review-owned: standalone review case. | Restricting activation to full runs fails the case. |
| SP3 | 3 | Shared preparation alone does not authorize delegated implementation. | Review-owned: instruction without approved delegated run. | Dispatch based only on that instruction fails the case. |
| SP4 | 3 | Shared preparation leaves the approved ticket graph unchanged. | Review-owned: dependency-bypass case. | Treating activation as graph approval fails the case. |
| SP5 | 4 | An absent instruction preserves each phase's existing route. | Review-owned: spec, build, and review without activation. | A new default fork route fails the phase comparison. |
| SP6 | 5 | The coordinator verifies required context, line, and isolation before fork dispatch. | Review-owned: one incompatible condition per case. | Dispatch with any incompatible condition fails the case. |
| SP7 | 5 | A compiled harness record does not prove native availability. | Review-owned: compiled capability with unavailable native tool. | Dispatch justified by the compiled record fails the case. |
| SP8 | 6 | A full native fork retains its parent's model and effort. | Review-owned: proposed different child line. | A conflicting model or effort override fails the case. |
| SP9 | 7 | A different line requires a documented lower estimated total. | Review-owned: changed line without comparison. | A cache-reuse assertion alone fails the case. |
| SP10 | 7 | The estimate includes preparation, delegates, output, and expected repairs. | Review-owned: omit each cost component in turn. | Each omission makes the line-change evidence incomplete. |
| SP11 | 8 | An unknown required cost input preserves the existing line. | Review-owned: unknown repair cost. | Treating the unknown as zero fails the case. |
| SP12 | 7 | An equal or higher estimated total preserves the existing line. | Review-owned: tied and higher totals. | A line change in either case fails the comparison. |
| SP13 | 9 | A parent that cannot fork does not receive a nested-fork plan. | Review-owned: Claude fork as proposed preparation parent. | Planning its child forks contradicts the active contract. |
| SP14 | 10 | An unsuitable build fork uses a fresh delegate with explicit prepared evidence. | Review-owned: build parent on an incompatible line. | An inherited incompatible line fails the route. |
| SP15 | 10 | An unsuitable review fork uses a fresh native delegate with explicit evidence. | Review-owned: reviewer parent with author history. | Forking contaminated history fails the route. |
| SP16 | 10 | The coordinator reports a fresh-delegate fallback. | Review-owned: permitted fallback report. | A silent route change fails the report. |
| SP17 | 24 | An unavailable native review surface retains the capable-harness handoff. | Review-owned: no native review tool. | A same-family CLI or inline axis fails the route. |
| SP18 | 11 | Spec authorship retains its approved-context fork. | Review-owned: each authorized spec-source kind. | A fresh author or unreviewed replacement summary fails the route. |
| SP19 | 12 | The normal spec phase adds no preparation session or authorship split. | Review-owned: optional spec activation. | Requiring another preparer or ticket author fails the route. |
| SP20 | 10, 11 | An unavailable spec fork requires a capable-session handoff. | Review-owned: spec dispatch without a suitable native fork. | A fresh-author substitute fails the mandatory fork contract. |
| SP21 | 13 | Implementation common preparation starts in a fresh implementation session. | Review-owned: proposed planning-history parent. | Reusing that history as the common baseline fails the case. |
| SP22 | 14 | Each ticket assignment follows the common preparation. | Review-owned: sibling ticket assignments. | Putting one ticket's assignment into the shared baseline fails the case. |
| SP23 | 17 | Repairs return to the recorded implementation author. | Review-owned: repair after a fork batch. | Creating a replacement for cache convenience fails the case. |
| SP24 | 15 | Existing build preflight supplies the mechanical charge inputs. | Review-owned: inherited context beside a required charge. | Skipping preflight because context exists fails the case. |
| SP25 | 15 | Each author verifies current assignment identity and required source bytes before action. | Review-owned: mismatched assignment or source. | Treating the parent's verification as sufficient fails the case. |
| SP26 | 15 | Changed source identity or required bytes cause charge regeneration. | Review-owned: changed tip and changed bytes cases. | Reusing the old charge fails either case. |
| SP27 | 16 | A later batch refreshes changed common material. | Review-owned: accepted predecessor changes a shared source. | Dispatch with the old shared source fails the case. |
| SP28 | 16 | A dependent ticket waits for its prerequisite checkpoint. | Review-owned: prepared successor with pending predecessor. | Early dispatch fails despite reusable context. |
| SP29 | 18 | A write delegate acts only in its assigned isolated worktree. | Review-owned: inherited parent checkout path. | Writing that parent checkout fails the case. |
| SP30 | 19 | Review preparation includes the frozen diff and current primary sources. | Review-owned: omitted frozen pair or source. | A review baseline without either input fails the case. |
| SP31 | 20 | Review forks inherit no author rationale or implementation reasoning. | Review-owned: author-history parent with a neutral final prompt. | The neutral final prompt cannot cure inherited reasoning. |
| SP32 | 21 | Review forks inherit no other axis's findings. | Review-owned: first return reaches parent before last sibling fork. | The later sibling inherits findings and fails the case. |
| SP33 | 22 | Each axis independently re-reads primary sources and derives findings. | Review-owned: axis that trusts prepared conclusions. | A declaration-only source check fails the case. |
| SP34 | 23 | Review consumes the existing shared evidence for one preparation attempt. | Review-owned: one frozen pair with multiple axes. | A second diff or consumer collection for that pair fails the case. |
| SP35 | 25 | Reports distinguish projected cost from actual charges. | Review-owned: estimate without a billing record. | Labelling that estimate actual fails the report. |
| SP36 | 25 | Evidence uses the existing assessment owner and provenance rules. | Review-owned: proposed parallel cost store. | A new store or schema fails the scope check. |
| SP37 | 26 | Child cost attribution excludes counters inherited from parent preparation. | Review-owned: cumulative child counter includes parent baseline. | Counting the baseline twice fails the attribution. |
| SP38 | 27 | Unavailable child measurements remain unknown. | Review-owned: missing child record and missing cache counter. | Reporting measured zero fails either case. |
| SP39 | 25, 30 | Fewer reads or a successful spawn do not establish measured savings. | Review-owned: spawn evidence without comparative billing evidence. | A savings claim fails the evidence comparison. |
| SP40 | 22, 29 | Review retains isolated venues and the exclusion of authors and coordinator. | Review-owned: author or coordinator proposed as an axis. | Shared preparation cannot make that reviewer independent. |
| SP41 | 29 | Existing terminal evidence, repairs, checkpoints, and reconciliation remain required. | Review-owned: omit each obligation in turn. | Each omission fails the baseline route comparison. |
| SP42 | 30 | A live compatibility claim cites observed context, line, and assignment evidence. | Review-owned: documentation-only compatibility claim. | A source citation alone fails the live-evidence claim. |
| SP43 | 28 | Each shared rule has one guidance owner within the existing budgets. | Review-owned owner sweep plus guidance-prose-budgets. | A competing rule copy or excess subject length fails validation. |
| SP44 | 28, 29 | Existing pinned guidance and prose checks remain green. | docs-currency-workflow and prose-mechanics. | A removed pinned duty or invalid prose turns its existing check red. |

### Edge inventory

The audience is every repository that links the kit.
The optional instruction is session input, not a filesystem input or CLI grammar change.
The cases above cover absent activation, unsuitable lines, unavailable forks, stale sources, contaminated review history, and unknown evidence.
The authoritative behavior inventory is the resolved decision tickets linked below.

There is no new directory reader.
An absent or empty ticket directory keeps the build phase's existing return to spec writing.
An absent or empty native evidence source keeps unknown measurements.
No test swaps a package variable for this change.
No executable trust chain, parser, cleanup refusal, or hostile-input class changes.
The inherited worktree and native-review safeguards continue to serve those callers.

**Won't handle:** Guaranteed cache hits — all three phases remain usable through their authorized native route.
**Won't handle:** Automatic compaction for cache optimization — the spec author still consumes its approved source.
**Won't handle:** Unsupported nested forks — build and review coordinators retain fresh native delegates.
**Won't handle:** Silent spec fallback — the spec caller retains the capable-session handoff.
**Won't handle:** Live-model comparative trials — ordinary approved work can supply observations without a new benchmark requirement.

## Ownership fences

The fences below are the exact union of ticket Writes paths.
The first group contains planned content changes.
The second contains current fixture pins that preflight requires beside those paths.
Those fixture paths preserve coverage; they do not authorize changed enforcement or planted diagnostics.

- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/shared-preparation.md`
- `.agents/skills/bench-craft-line/SKILL.md`
- `.agents/skills/bench-craft-review/SKILL.md`
- `.bench/BENCH-reference.md`
- `CHANGELOG.md`
- `CONTEXT.md`
- `reviews/shared-delegate-startup.md`
- `specs/shared-delegate-startup/spec.md`

Preflight-required fixture closure:

- `tests/canary/claude-agent-definitions/agent-unnamed-in-skill`
- `tests/canary/claude-agent-definitions/model-declared`
- `tests/canary/claude-agent-definitions/name-mismatch`
- `tests/canary/claude-agent-definitions/shell-tool-absent`
- `tests/canary/claude-agent-definitions/skill-names-missing-agent`
- `tests/canary/claude-agent-definitions/spawning-tool`
- `tests/canary/claude-agent-definitions/tools-absent`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/acceptance-coverage-anchor`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/command-handoff-anchor`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/coverage-axis-anchor`
- `tests/canary/workflow-guidance-anchors/craft-review-coverage-row-projection`
- `tests/canary/workflow-guidance-anchors/decision-map-asset-path`
- `tests/canary/workflow-guidance-anchors/delegate-cap-change-pinning-package`
- `tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap`
- `tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge`
- `tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green`
- `tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer`
- `tests/canary/workflow-guidance-anchors/delegate-exec-only-every-caller`
- `tests/canary/workflow-guidance-anchors/delegate-model-id-escalation`
- `tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface`
- `tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor`
- `tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance`
- `tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents`
- `tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row`
- `tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor`
- `tests/canary/workflow-guidance-anchors/delegated-author-limit`
- `tests/canary/workflow-guidance-anchors/delegated-axis-exclusions`
- `tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review`
- `tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration`
- `tests/canary/workflow-guidance-anchors/delegated-entry-refusals`
- `tests/canary/workflow-guidance-anchors/delegated-mid-review-route`
- `tests/canary/workflow-guidance-anchors/delegated-per-ticket-author`
- `tests/canary/workflow-guidance-anchors/delegated-resumption-contents`
- `tests/canary/workflow-guidance-anchors/delegated-tier-authorization`
- `tests/canary/workflow-guidance-anchors/delegated-unbound-model-stop`
- `tests/canary/workflow-guidance-anchors/edge-inventory-anchor`
- `tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding`
- `tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer`
- `tests/canary/workflow-guidance-anchors/implement-spec-entry-validation`
- `tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness`
- `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception`
- `tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper`
- `tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight`
- `tests/canary/workflow-guidance-anchors/implement-spec-write-delegation`
- `tests/canary/workflow-guidance-anchors/line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/prepared-build-approval`
- `tests/canary/workflow-guidance-anchors/prepared-build-freshness`
- `tests/canary/workflow-guidance-anchors/prepared-review-axis-returns`
- `tests/canary/workflow-guidance-anchors/prepared-review-blast-evidence`
- `tests/canary/workflow-guidance-anchors/prepared-review-capable-handoff`
- `tests/canary/workflow-guidance-anchors/prepared-review-delegate-handoff-route`
- `tests/canary/workflow-guidance-anchors/prepared-review-inline-axis-route`
- `tests/canary/workflow-guidance-anchors/prepared-review-legacy-entry-points`
- `tests/canary/workflow-guidance-anchors/prepared-review-native-dispatch`
- `tests/canary/workflow-guidance-anchors/prepared-review-runtime-capability`
- `tests/canary/workflow-guidance-anchors/prepared-review-shared-evidence`
- `tests/canary/workflow-guidance-anchors/prepared-triage-bounds`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `tests/canary/workflow-guidance-anchors/review-base-merged-main-tip`
- `tests/canary/workflow-guidance-anchors/review-clean-terminal-result`
- `tests/canary/workflow-guidance-anchors/review-cross-harness-opt-in`
- `tests/canary/workflow-guidance-anchors/review-falsification-accept-routing`
- `tests/canary/workflow-guidance-anchors/review-falsification-dispositions`
- `tests/canary/workflow-guidance-anchors/review-finding-discipline-pointer`
- `tests/canary/workflow-guidance-anchors/review-persistence-anchor`
- `tests/canary/workflow-guidance-anchors/review-preflight-explicit-base`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-covers`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-owner`
- `tests/canary/workflow-guidance-anchors/review-standing-falsification`
- `tests/canary/workflow-guidance-anchors/review-universal-claim-bar`
- `tests/canary/workflow-guidance-anchors/shared-worktree-path-pin`
- `tests/canary/workflow-guidance-anchors/slicing-in-write-spec`
- `tests/canary/workflow-guidance-anchors/spec-retire-roadmap-row`
- `tests/canary/workflow-guidance-anchors/staged-command-sweep-anchor`
- `tests/canary/workflow-guidance-anchors/story-line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/ticket-decision-map-lifecycle-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-stage-routing-anchor`
- `tests/canary/workflow-guidance-anchors/write-spec-artifact-authorization`
- `tests/canary/workflow-guidance-anchors/write-spec-authorization-boundary`
- `tests/canary/workflow-guidance-anchors/write-spec-branch-split`
- `tests/canary/workflow-guidance-anchors/write-spec-conversation-authorization`
- `tests/canary/workflow-guidance-anchors/write-spec-conversation-fork`
- `tests/canary/workflow-guidance-anchors/write-spec-decision-source`
- `tests/canary/workflow-guidance-anchors/write-spec-frozen-base-and-tip-review`
- `tests/canary/workflow-guidance-anchors/write-spec-handoff-anchor`
- `tests/canary/workflow-guidance-anchors/write-spec-late-uncertainty`
- `tests/canary/workflow-guidance-anchors/write-spec-map-sources`
- `tests/canary/workflow-guidance-anchors/write-spec-phase-ownership`
- `tests/canary/workflow-guidance-anchors/write-spec-post-slicing-handoff`
- `tests/canary/workflow-guidance-anchors/write-spec-reader-sweep-sequence`
- `tests/canary/workflow-guidance-anchors/write-spec-ready-map-authorization`
- `tests/canary/workflow-guidance-anchors/write-spec-review-made-conditional`
- `tests/canary/workflow-guidance-anchors/write-spec-review-tier-escalated`
- `tests/canary/workflow-guidance-anchors/write-spec-slicing-sign-off`
- `tests/canary/workflow-guidance-anchors/write-spec-ticket-approval-fields`
- `tests/canary/workflow-guidance-anchors/write-spec-ticket-breakdown-charge`
- `tests/canary/workflow-guidance-anchors/write-spec-verification-learning`
- `tests/canary/workflow-guidance-anchors/write-spec-verification-learning-contents`
- `tests/canary/workflow-guidance-anchors/write-spec-verification-log`

No change to budget limits, default tier bindings, Go sources, CLI grammar, hooks, or adapters is planned.
The spec's own fence permits verification records within the approved behavior.
Spec amendments use the existing plan-expansion policy.
Reviewer disposition: pending the completed spec-and-tickets sign-off.

## Out of scope

The build excludes other staged specs and all implementation ticket files.
Only this spec's own verification record remains within its documentation fence.

- A controlled default-change comparison: 4 authored artifacts, 2 gate runs. A plan, trial evidence, assessment, and adoption decision form a separate capability.
- A new native telemetry importer: 6 code or test edits, 2 gate runs. Provider parsing, attribution, validation, and projection require their own spec.
- A new automatic fork launcher: 7 code, adapter, or test edits, 2 gate runs. Native capability resolution and dispatch need a separate interface.

These are planning estimates, not authorization or measured implementation costs.
They cover separate capabilities and omit no remainder of the optional guidance workflow.

## Further notes

### Source verification

The author re-read every structured source from the compiled map on 2026-09-16.
The source baseline for this draft is `478e5c079c36ca3c11233f638ac95d0c7c572545`.
The relocated research report retains the factual source manifest; this spec does not duplicate it.
Both primary upstream pages were accessible.

Claude documents shared fork context and cache, inherited model, and the no-nested-fork limit.
OpenAI describes model-dependent prompt-prefix reuse and distinct usage counters.
Neither page establishes measured savings for this Bench run.

The active native tool contract permits a full-history fork without model or effort overrides.
This spec phase successfully dispatched its first full-history author fork.
That result does not expose cache counters or prove a comparative cost improvement.
No live Claude fork experiment ran.

### Enforcement reads and reader sweep

The author inspected current producers before choosing the guidance seam:

| source read | observed responsibility | disposition |
| --- | --- | --- |
| `internal/preflight/charge.go`, `chargeSources`, `loadChargeSources`, `readChargeSource` | Build inputs and exact source-byte comparison. | Preserve through SP24–SP26; no code edit. |
| `internal/preflight/review.go`, `collectReviewEvidence`, `renderReviewCharge`, `renderReviewPacket` | One shared evidence set and axis charge production. | Preserve through SP34; no new transport. |
| `internal/lines/lines.go`, `AgentLineVerdict` | Fork model-override guard and declared-model rules. | Preserve through SP8; no hook edit. |
| `internal/assessment/README.md`, `types.go`, `cost.go`, `harness.go` | Existing evidence schema, provenance, estimates, and native counters. | Consume through SP35–SP38; no schema edit. |
| `internal/harnesstranscript/codex.go`, counter boundaries and `measures` | Cumulative tokens and unavailable nested observations. | Preserve unknowns through SP37–SP38. |
| `internal/anchors/registry_ft311_preparation.go` | Existing charge, approval, and freshness anchors. | Preserve exact duties through SP24–SP26. |
| `internal/anchors/registry_ft311_review_dispatch.go` | Native review, shared evidence, and fallback anchors. | Preserve exact duties through SP17, SP34, SP40–SP41. |
| `internal/conformance/docs_workflow_checks_test.go`, `checkSpecAuthorizationContract` | The three permitted spec-source kinds. | Preserve through SP18. |
| `internal/conformance/prose_budget_test.go`, `checkGuidanceProseBudgets` | Budget policy and subject enumeration. | Preserve limits through SP43. |
| `internal/conformance/prose_mechanics_test.go`, `checks_test.go`, `registry_test.go` | Named prose and workflow checks in the dev gate. | Reuse through SP44. |
| `.agents/commands/bench-write-spec.md`, `bench-implement-spec.md`, `bench-review-implementation.md` | Phase consumers of preparation and authorization. | Each receives its phase-specific pointer and duties. |
| `.agents/skills/bench-craft-delegate/SKILL.md`, `references/delegation-discipline.md` | Charge, assignment, isolation, retained author, and verification rules. | Shared owner changes; existing discipline remains unchanged. |
| `.agents/skills/bench-craft-line/SKILL.md`, `.agents/skills/bench-craft-review/SKILL.md` | Line selection and independent source derivation. | Exact candidate fences above. |
| `.bench/BENCH.md`, `projects/benchkit.md`, `.bench/BENCH-reference.md` | Shared policy, budgets, line binding, and assessment use. | Only the reference needs cost-collection guidance. |
| `capture/agent-performance/open-ai-models.md` | Current routing and observed review limitations. | Supports Sol/high guidance implementation; no current-task outcome inferred. |

The full hidden-file path sweep found only two repository-relative references to the moved map.
Both moved with the decision source and now use the compiled location.
The map's relative ticket links and research backlink remain valid after the unit move.
The sweep included dot-directories, scripts, workflows, fixtures, and conformance sources.
No executable reader names the old topic path.

Existing staged `delegate-boot-cost`, `delegated-implementation`, and `bounded-charge-evidence` artifacts are reference material, not editable implementation consumers.
The last proposes another charge transport; this spec consumes the current canonical owner and creates no competing transport.
No count, serialized field, command token, or runtime function changes.
The shipped claim is optional workflow guidance for linked repositories.
Keep kit-only assessment implementation paths out of new shipped instructions; link through the existing reference owner.

### Source-sentence-to-row traceability

Ticket paths below are relative to `decisions/shared-delegate-startup/tickets/` in this spec folder.
The quoted clauses identify the canonical source occurrences.

| source occurrence | source clause | rows |
| --- | --- | --- |
| 4.md, Answer 1 | “One explicit instruction covers an approved delegated run and its reviews.” | SP1 |
| 4.md, Answer 2–4 | “Standalone phases can receive the same explicit instruction.” | SP2–SP5 |
| 3.md, Answer 1 | “A documented total-cost estimate permits a different line when the estimated total is lower.” | SP9, SP12 |
| 3.md, Answer 2 | “The estimate covers preparation, all delegates, output, and expected repairs.” | SP10 |
| 3.md, Answer 3 | “Unknown required cost inputs retain the existing line.” | SP11 |
| 3.md, Answer 4–5 | “Actual results remain separate from estimates.” | SP8, SP35 |
| 5.md, Answer 1 | “Preserve the existing approved-context spec-author fork.” | SP18 |
| 5.md, Answer 2 | “Apply the shared preparation and cost-accounting rules without an extra preparation session or split authorship.” | SP10, SP19 |
| 5.md, Answer 3–4 | “The optional workflow does not relax the phase's existing fork requirement.” | SP8, SP18, SP20 |
| 7.md, Answer 1 | “Implementation delegates inherit common preparation from a fresh implementation session.” | SP21 |
| 7.md, Answer 2 | “Each delegate receives its ticket assignment after that preparation.” | SP22 |
| 7.md, Answer 3 | “Each author validates its current source and assignment before action.” | SP24–SP26, SP29 |
| 7.md, Answer 4 | “Later batches refresh changed material.” | SP27–SP28 |
| 7.md, Answer 5 | “Repairs stay with the existing author under the retained-author rules.” | SP23 |
| 8.md, Answer 1 | “Review delegates inherit neutral preparation containing the frozen diff and primary sources.” | SP30, SP34 |
| 8.md, Answer 2 | “That preparation contains no author rationale or other reviewers' findings.” | SP31–SP32 |
| 8.md, Answer 3–4 | “The inherited source material does not replace required independent source checks.” | SP33, SP40 |
| 9.md, Answer 1–2 | “The coordinator reports the fallback.” | SP14–SP16 |
| 9.md, Answer 3 | “A suitable fork preserves the required context, line, and isolation.” | SP6–SP8, SP13 |
| 9.md, Answer 4 | “Spec writing retains its mandatory approved-context fork and uses a capable-session handoff if that fork is unavailable.” | SP20 |
| 9.md, Answer 5 | “The fallback does not waive current source checks, review independence, or retained authorship.” | SP17, SP23–SP26, SP33, SP40–SP41 |
| 6.md, Answer 2–4 | “The workflow uses the existing preparation and assessment tools.” | SP24, SP34–SP39, SP43–SP44 |
| 1.md, Answer 5–6 | “No live probe or savings comparison ran, so measured benefit and native compatibility remain unknown.” | SP7, SP38–SP39, SP42 |
| 2.md, Answer 1–2 | “The reviewer selects an optional workflow first, without a controlled comparison as its entry condition.” | SP1–SP5, SP39 |
| research.md, Q2 and Q3 | Native inheritance, nesting limits, cumulative counters, and missing child attribution. | SP8, SP13, SP37–SP39, SP42 |

### Pre-review proof checklist

- Cited symbols: the enforcement-read table resolves the named current symbols.
- Import edges: none.
- Source-row clauses and occurrences: the traceability table above names canonical decision occurrences.
- Promised field labels: none beyond the existing spec template and ticket grammar.
- Changed-function callers: none.
- Copy survival: no production owner or code copy is replaced; the approved map moves as one unit.

Flagged additions beyond the decision source: none.
The native-tool handoff, retained-author duties, and cost provenance are preserved constraints from current source.
The named owner reference and semantic verification cases are spec-author implementation details.

### Completion plan

The default implementation venue remains one retained session on Sol/high.
This plan does not authorize delegated implementation.
An explicitly approved delegated build records its required author assignments before dispatch under the current completion-plan policy.

At each serial ticket, retain focused checks and the owned semantic cases in the review pickup.
Each case records current and proposed guidance, its exact counterexample, and the observed or unresolved result.
The independent implementation review checks those cases across the complete SP1 delta.
Semantic cases remain review-owned, not executable completion-plan commands.
The commands below grade only their existing mechanical predicates.
No successful command proves cache reuse or native review independence.

After SP1, retain author verification and all three native review results in the review pickup.
Close accepted repairs with current evidence before the SP1 checkpoint and final reconciliation.
The coordinator runs the whole-project gate only through the normal checkpoint and landing workflow.
No implementation or native benchmark ran during this specification phase.

```bench-completion-plan
{"version":1,"chunks":[{"id":"SP1","tickets":["1-align-spec-preparation.md","2-prepare-implementation-delegates.md","3-preserve-neutral-review.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"},{"id":"prose","command":"bench test --check prose-mechanics"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/shared-delegate-startup/spec.md"},{"id":"workflow","command":"bench test --check docs-currency-workflow"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"},{"id":"prose","command":"bench test --check prose-mechanics"}]}
```

### Specification verification

The user selected one independent Sol/high pass over the completed spec and tickets.
The reviewer examined commit 194128673b7dbf71cb6609691f1b19422749b52c against base 478e5c079c36ca3c11233f638ac95d0c7c572545.
The pass found two blockers and no source-fidelity finding.
The ticket author corrected both in one returned pass.
The coordinator verified the corrections without a second independent review.

| finding | correction | coordinator evidence | disposition |
| --- | --- | --- | --- |
| Exclusion tokens entered parsed write authority. | Exclusion prose moved outside Ownership fences. | The production FenceTokens function returns exactly the unique ticket Writes union. | closed |
| The first ticket owned combined activation before both consumers existed. | The final consumer ticket owns SP1 and its combined-route acceptance. | Each row has one owner, and only the final ticket owns SP1. | closed |

The commit lane and committed-source build preflight pass.
The scope and ownership fences still await the user's implementation sign-off.
Successful author dispatches establish no measured cache or cost benefit.
