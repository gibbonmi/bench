# Craft research: one model-invoked skill for factual reading legwork

Status: staged

Decision source: the named reviewed artifact `specs/decision-map-split/decisions/craft-research.md`, tickets #4 through #9, #13, and #14, resolved 2026-08-02 and 2026-09-06. That folder retires with the sibling spec, so the settled answers this spec rests on are restated under Implementation decisions. The three upstream URLs in the research asset were not re-fetched on 2026-09-06. The `craft-spec` compatibility-probe wording that map ticket #8 replaces was already removed by a later remake. That clause lands as a pointer beside the surviving spec-side duty.

Verification log: 2 iteration(s) to accept — round one accepted with folds, round two verified them, and the 120-line budget stays an open reviewer question

## Problem

Research policy is scattered. The shaping command holds one Research bullet, the delegate skill holds the read-only mechanics, and the assessment command holds its own fan-out and re-verification prose. No surface states the primary-source standard, the coordinator's duty to check the joins between returns, or the report shape a reader can decide from. A research asset lands wherever the session puts it. A research claim about byte or wire compatibility can close a decision without a probe.

## Solution

One model-invoked skill, `craft-research`, owns factual reading legwork. It owns the trigger, the question graph, the primary-source standard, and adaptive round-based fan-out. It owns coordinator synthesis and verification, the artifact contract, the report contract, and the residual-unknown posture. It composes with `craft-delegate` for charge and isolation mechanics and with `craft-line` for every model, effort, cap, and fan-out declaration. The shaping and assessment commands point at it. The mirror symlink and the skills index register it, and the anchor registry pins its load-bearing sentences.

## User stories

### The skill

Line: opus / high. A skill steers every later research run, so the leverage override routes it mid plus high.

1. As a coordinator, I want one model-invoked skill that fires when work becomes factual reading legwork, so that no phase restates research policy.
2. As a coordinator, I want a question graph where every load-bearing claim traces to one question node, so that fan-out has a unit.
3. As a coordinator, I want primary sources only, with a secondary source able to locate but not warrant, so that no summary is evidence.
4. As a coordinator, I want fan-out that dispatches only mutually independent frontier questions, so that dependent reads stay serial.
5. As a coordinator, I want the round width declared through the fan-out clause before dispatch, so that spend is visible.
6. As a coordinator, I want each delegate charge to follow the delegate discipline and add the primary-source boundary, so that mechanics are not restated.
7. As a coordinator, I want a completed round synthesized before the next opens, so that a new round names the gap it closes.
8. As a coordinator, I want to re-open every source behind a load-bearing conclusion and check every join, so that true claims cannot mislead.
9. As a coordinator, I want one coordinator-authored durable output per run keyed by topic, so that delegate returns stay ephemeral.
10. As a coordinator, I want one destination precedence for the output, so that one rule places every asset.
11. As a reader, I want the output to open with recommendation, scope, and evidence status, so that I can decide from the first screen.
12. As a reader, I want facts, inferences, tested results, and proposals kept apart, so that a proposal cannot read as a fact.
13. As a reader, I want option tables with consequence and citation, kept contradictions and unknowns, and a validation plan, so that the report is whole.
14. As a reader, I want a diagram where prose hides a material relation, so that I can inspect the join.
15. As a cold reader, I want every load-bearing claim cited to a local line or a dated primary URL, so that I can re-open it.
16. As a cold reader, I want `Consumed by`, `Drift`, and `Retire when` on an asset outside a phase-owned output, so that staleness has a trigger.
17. As a reviewer, I want a compatibility claim to stay unverified until a separate runnable probe returns, so that research never promises compatibility.
18. As a reviewer, I want research to own no write delegate, done-claim, reviewer decision, or prototype, so that judgment stays with the phase.
19. As a coordinator, I want a contrastive example of independent fan-out versus a dependent serial question, so that the skill teaches by contrast.
20. As a coordinator, I want a report checklist to tick in the verification record, so that the next run measures the contract.
21. As a reviewer, I want the skill inside a named 122-line budget in ASD-STE100, so that it loads cheap and passes the prose lane.

### Integration

Line: opus / high. The pointer sentences steer the shaping and assessment phases, so the leverage override applies.

22. As a shaping session, I want the Research bullet and the execution paragraph to point at the skill, so that shaping restates no research rule.
23. As a shaping session, I want a required compatibility probe as a Prototype ticket named in `Blocked by`, so that the Research ticket waits.
24. As an assessing session, I want the assessment command to charge the skill for fan-out, re-verification, citations, and unknowns, so that one owner holds them.
25. As an assessing session, I want the six areas, the severity grammar, the backlog, and the replace-in-place lifecycle kept, so that assessment keeps its shape.
26. As a spec session, I want the bootstrap-authority reference to point at the skill's probe rule, so that one owner holds that rule.
27. As any session, I want `craft-delegate` and `craft-line` unchanged, so that the overlap stays deliberate composition.
28. As a Claude Code session, I want the mirror symlink to resolve to the skill, so that the skill loads and the mirror check passes.
29. As a cold reader, I want the skills index regenerated with the new trigger line, so that the index check passes.
30. As the gate maintainer, I want the skill's load-bearing sentences anchored with mutation rows, so that a rewrite cannot drop them silently.
31. As the reviewer, I want the next research run's asset to record the checklist hits, so that FT231 gets the report measure.

## Implementation decisions

**The skill file.** `.agents/skills/bench-craft-research/SKILL.md` carries frontmatter with `name: craft-research`, a `description` that names the trigger, and an `index:` line for the skills index. It is model-invoked on both harnesses and needs no phase command, adapter, payload row, parser, or source schema. The body is ASD-STE100 prose within a named 122-line budget row. The skill ticket adds that row to the project profile in the same commit, because a row for an absent file reds the budget check. It carries one contrastive pair: an independent fan-out that runs in parallel against a dependent question that stays serial.

**Trigger and ownership.** The skill fires when work becomes factual reading legwork in shaping, specification, diagnosis, assessment, or implementation. Each calling phase keeps authority over its decisions, artifacts, and completion contract. Formal review is not a caller. Research establishes source-backed facts, contradictions, unknowns, and implications. It never owns a write delegate, a done-claim, a reviewer decision, or a prototype.

**Question graph and fan-out.** The coordinator draws the factual question graph first. It is complete when every load-bearing claim the destination needs traces to one node. A trivial lookup resolves inline. One bounded question goes to one read-only delegate when its read set would displace the coordinator's context.

Fan-out happens only when at least two frontier questions are mutually independent. A round's width is that number, capped by what one synthesis pass can verify, and declared through the fan-out clause before dispatch. Dependent questions and concurrency-sensitive measurements run serially, and the calling phase's iteration cap bounds the rounds.

**Sources and verification.** The artifact under study and first-party upstream documentation or APIs are primary. A secondary source may locate evidence but cannot warrant a finding. Beyond the delegate discipline's verification duties, the coordinator re-opens every source behind a load-bearing conclusion and checks every join between returns. Completion requires every in-scope question answered or explicitly unknown, every material claim traceable, and every contradiction reconciled or retained.

**The durable output.** One coordinator-authored Markdown output per run, keyed by topic, with one section per question. The destination precedence has four steps. A phase-owned output wins when the caller has one. Otherwise the output sits beside the artifact that will consume it. For shaping that is the map's assets folder under the topic, and for the repository it is `docs/research/<topic>.md`.

The sibling spec makes both folders real. Before dispatch the caller names the destination, the consumer, and the retire or refresh condition. An asset outside a phase-owned output states `Consumed by`, `Drift`, and `Retire when`.

**The report contract.** The output opens with the recommendation, the scope, and the evidence status. It keeps facts, inferences, tested results, and proposals in separate sections. It uses a capability or option table with a consequence and a nearby citation when comparison matters. It keeps contradictions and unknowns, draws a diagram where prose hides a material relation, and ends with a validation plan.

Every load-bearing claim cites an exact local path and line or a primary URL with its retrieval date. The skill lists these elements as a checklist the coordinator ticks in the output's verification record.

**Compatibility evidence.** A byte or wire compatibility claim stays unverified until a separate runnable probe returns. In a decision map that probe is a Prototype ticket, and the Research ticket names it in `Blocked by`.

**Callers.** The shaping command's Research bullet and execution paragraph point at the skill. They keep only the ticket type, the dependency rule, and the map state duties. The assessment command charges the skill in its fan-out and synthesis steps. It keeps its five anchored phrases, six areas, previous-backlog reconciliation, severity grammar, and ranked backlog. It also keeps its verification marks and its replace-in-place lifecycle.

The bootstrap-authority reference keeps "compatibility proven, not promised" and gains one pointer to the skill's probe rule. `craft-delegate` and `craft-line` do not change, including the delegate skill's fan-out-search trigger.

**Registration.** `.claude/skills/bench-craft-research` is a relative symlink like its siblings. `bench skills-index --write` regenerates the index block in the reference file. The load-validity, skills-index, and guidance-prose-budget checks grade the result.

**Anchors.** The registry gains one needle per load-bearing sentence listed under Further notes. Each needle gets a mutation row in the registry test and a fixture under the workflow-guidance-anchors family. The pointer sentences in the two commands get needles too.

**Shared files with the sibling spec.** The shaping command, the anchor registry and its test, and the workflow-guidance-anchors fixtures are edited by both specs. This spec lands second, per map ticket #14, and its build re-reads the landed wording before it edits.

## Testing decisions

- A good test loads the tree and runs the conformance check that grades the surface. Content rows that no check can grade are review-owned and say so.
- The registration seam is the `load-validity-metadata` and `skills-index-command-adapters` checks. Prior art: `checkClaudeSkillMirror`, `checkCraftSkillNames`, and `checkSkillsIndex`.
- The budget seam is the `guidance-prose-budgets` check and the prose lane.
- The anchor seam is the registry mutation table and the `workflow-guidance-anchors` fixtures. Prior art: `TestCraftDelegateDisciplineAnchorsRedOnRemoval`.
- The gate observes the feature through those checks and `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.

### Seam diagram

    trigger: a phase reaches factual reading legwork
        │
        ▼
    question graph  ──▶  [ craft-research: rounds, delegates, synthesis, verification ]  ──▶  one cited Markdown output
                                    ◀ tests attach here: the skill file, graded by the conformance checks and the anchor registry

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CR1 | 1 | `.agents/skills/bench-craft-research/SKILL.md` exists with frontmatter `name: craft-research`, a description, and an `index:` line | conformance check `load-validity-metadata` | A wrong name or missing frontmatter reds the check |
| CR2 | 1 | The skill carries no `disable-model-invocation` line and no `agents/openai.yaml` | review-owned: the skill ticket cites the folder listing | An adapter makes a craft skill a phase |
| CR3 | 2, 4, 7 | The skill states the question-graph completeness rule, the independence rule for fan-out, and the synthesize-before-next-round rule | `TestCraftResearchAnchorsRedOnRemoval` | A reworded rule without its needle passes the docs check |
| CR4 | 3, 8 | The skill states the primary-source rule and the re-open-and-check-joins rule | `TestCraftResearchAnchorsRedOnRemoval` | A skill that trusts delegate summaries passes without the needle |
| CR5 | 5, 6 | The skill points at the fan-out clause for width and at the delegate discipline for charge mechanics, and states no tier or effort | review-owned: the skill ticket quotes both pointers | A restated tier table duplicates the line skill |
| CR6 | 9, 10, 16 | The skill states one output per run, the four-step destination precedence, and the three metadata labels | `TestCraftResearchAnchorsRedOnRemoval` | A skill without the precedence lets each session choose a home |
| CR7 | 11, 12, 13, 14 | The skill's report section lists recommendation-first, the four separated sections, the option table, kept unknowns, the diagram rule, the validation plan, one section per question, every question answered or unknown, and a statement of what could not be verified | review-owned: the skill ticket cites the section | A report contract missing one element lands a thinner report |
| CR8 | 15 | The skill states the citation rule: exact local path and line, or a primary URL with retrieval date | `TestCraftResearchAnchorsRedOnRemoval` | A citation rule without a needle can soften to "cite sources" |
| CR9 | 17, 23 | The skill states that a compatibility claim stays unverified until a runnable probe returns, and the shaping command names the Prototype ticket rule | `TestCraftResearchAnchorsRedOnRemoval` | A research run that closes a compatibility claim passes without the needle |
| CR10 | 18 | The skill states the read-side boundary: no write delegate, done-claim, reviewer decision, or prototype | `TestCraftResearchAnchorsRedOnRemoval` | A skill that absorbs probes blurs the phase boundary |
| CR11 | 19 | The skill carries one contrastive pair marked as such | review-owned: the skill ticket quotes the pair | A skill without the pair fails the craft-skills rule |
| CR12 | 20 | The skill's checklist names each report element and says where the tick lands | review-owned: the skill ticket cites the list | A checklist without a landing place is never ticked |
| CR13 | 21 | The skill passes `bench gate-prose` and the `guidance-prose-budgets` check at or under 122 lines with its named row present | conformance check `guidance-prose-budgets` and the prose lane | A longer skill, or a row without its file, reds the check |
| CR14 | 22 | The shaping command's Research bullet names `craft-research`, and the three forbidden fragments listed under Further notes are absent | `TestCraftResearchAnchorsRedOnRemoval` with Forbid rows | A bullet that keeps the old rules is a second copy |
| CR15 | 24, 25 | The assessment command names `craft-research` in its fan-out and synthesis steps and keeps its five anchored phrases | the assessment anchors in the `workflow-guidance-anchors` family | A rewrite that drops one phrase reds |
| CR16 | 26 | The bootstrap-authority reference keeps "compatibility proven, not promised" and points at the skill's probe rule | `TestCraftResearchAnchorsRedOnRemoval` | A reference without the pointer leaves two owners |
| CR17 | 27 | `git diff` for the build touches none of `bench-craft-delegate/SKILL.md`, `bench-craft-line/SKILL.md`, and `bench-craft-tickets/SKILL.md` | review-owned: the review reads the diff file list | An edit there breaks the composition decision |
| CR18 | 28 | `.claude/skills/bench-craft-research` is a relative symlink that resolves to `.agents/skills/bench-craft-research/` | conformance check `load-validity-metadata` | A missing or wrong symlink reds the mirror check |
| CR19 | 29 | The skills index block carries the new trigger line in alphabetical position | conformance check `skills-index-command-adapters` | A stale index reds with the regenerate hint |
| CR20 | 30 | Every needle listed under Further notes has a mutation row that bites | `TestCraftResearchAnchorsRedOnRemoval` | A needle without a mutation row is a claim |

Not covered: story 31 — the checklist hits are recorded by the next research run's phase close, which no ticket here owns.

### Edge inventory

- A phase with a phase-owned output, such as assessment: the skill writes into that output and creates no second file. CR6 covers the precedence.
- A research question that needs a shared mutable state or a load-sensitive measurement: stays serial. CR3 covers the independence rule.
- A round that exhausts the calling phase's cap: stops with residual unknowns. **Won't handle** a skill-owned cap — `craft-line` owns caps, and CR5 covers the pointer.
- A secondary source with no reachable primary: the output records an unknown. CR4 covers the rule.
- A claim of byte or wire compatibility: CR9.
- Formal review as a caller: **Won't handle** — the review phase owns its axes, and CR14 covers the shaping caller.
- A linked repository: the skill ships through the consumer payload as its siblings do, and the same checks grade it there. No kit-only marker is added.
- The negative-space read of the draft: the skill ticket lists each silence it fills and each it leaves open. Review-owned under CR7.

## Ownership fences

- `.agents/skills/bench-craft-research/`
- `.claude/skills/bench-craft-research`
- `.bench/BENCH-reference.md`
- `projects/benchkit.md`
- `.agents/commands/bench-shape-idea.md`
- `.agents/commands/bench-assess.md`
- `.agents/skills/bench-craft-spec/references/bootstrap-authority.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_craft_research.go`
- `internal/anchors/registry_craft_research_test.go`
- `internal/conformance/registry_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/skills-index-command-adapters/`
- `tests/canary/guidance-prose-budgets/`
- `tests/canary/line-routing/`
- `tests/canary/docs-currency-token-diet/`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `reviews/craft-research-skill.md`
- `capture/retros/`

## Build decisions

**Fence extension, 2026-09-06.** The anchor registry file and its test file are over the line budget, and the commit lane refuses growth there. The research anchor rows and their mutation table therefore live in `internal/anchors/registry_craft_research.go` and `internal/anchors/registry_craft_research_test.go`, on the precedent of the decision-map anchor pair. Both files join the ownership fence. The content of the rows is unchanged.

**Fixtures per needle, 2026-09-06.** The registry mutation table is the per-needle proof, and each rule is shown red in both directions. The `workflow-guidance-anchors` family holds one fixture per edited subject file: the skill, the shaping command, and the bootstrap-authority reference. Fourteen fixture copies of the same needles would be a second source for the table.

**Review repair fence, 2026-09-06.** The Coverage axis found that `internal/conformance/registry_test.go` lists the anchor registry files that own the fixture family, and the list lacked the new file. That file joins the ownership fence for the one-line addition.

**Assessment step 2 wording.** The read-only clause of the area sweeps moved to the skill, which owns the read-side boundary. The step keeps its tier, its one-delegate-per-area shape, and its six areas.

## Out of scope

- A general-purpose knowledge base, citation database, or web-search CLI: own spec, 0 edits here.
- A replacement for `craft-delegate`, `craft-line`, or harness-native subagent controls: rejected.
- Folding formal review's axis judgment into research: rejected.
- The decision-map storage, projection, and asset home: the sibling spec `decision-map-split`.
- The ignore-rule narrowing: moved to the sibling spec's migration, because the orphan move needs a trackable `docs/research/`.

## Further notes

**Flagged additions beyond the decision source.**

- The pointer sentence in the bootstrap-authority reference replaces the removed duplicate the map named.
- The report checklist is the measurement surface map ticket #14 asked for.

**Source sentence to row.**

- Map #4, reach-for-it-anytime trigger and caller authority: CR1, CR14, CR15.
- Map #4, review is not a caller: Edge inventory.
- Map #5, question graph, independence, width declared, synthesize before next round: CR3, CR5.
- Map #5, charge follows craft-delegate plus the source boundary: CR5.
- Map #6, one output per run, destination precedence, metadata: CR6.
- Map #6, primary sources, citations, coordinator re-verification: CR4, CR8.
- Map #7, read-side only, probes separate: CR9, CR10.
- Map #8, skill file and callers: CR1, CR14, CR15, CR16, CR17, CR18, CR19.
- Map #9, report contract: CR7, CR12.
- Map #13, asset homes: CR6.
- Map #14, spec 2 contents and checklist measure: CR12 and the Not covered line for story 31.

**Anchored sentences the skill ticket registers.**

- The question graph is complete when every load-bearing claim traces to one question node.
- Fan out only when at least two frontier questions are mutually independent.
- Synthesize a completed round before opening another.
- The artifact under study and first-party upstream documentation or APIs are primary sources.
- The coordinator re-opens every source supporting a load-bearing conclusion and independently checks every join between returns.
- One coordinator-authored durable output per research run, keyed by topic.
- Every load-bearing claim cites an exact local path and line or a primary URL with its retrieval date.
- A byte or wire compatibility claim stays unverified until a separate runnable probe returns.
- Research never owns a write delegate, a done-claim, a reviewer decision, or a prototype.
- In the shaping command: the Research ticket type points at `craft-research`.
- In the bootstrap-authority reference: the pointer to the probe rule.
- Forbidden in the shaping command after the pointer lands: `and produce a short`, `Include a runnable compatibility probe when the answer`, and `as a read-only delegation. Otherwise resolve it inline.`

**Reviewer sign-off, 2026-09-06.** The reviewer approved the named budget row `.agents/skills/bench-craft-research/SKILL.md | 122`, the ignore-rule move to the sibling spec, the fences, and the ticket graph. The row lands in the skill ticket's commit.

**Pre-review proof checklist.**

- Cited symbols: `checkClaudeSkillMirror`, `checkCraftSkillNames`, `checkSkillsIndex`, and `TestEveryRetainedFixtureBitesThroughRegisteredOwner` resolve in the tree read on 2026-09-06.
- Import edges: none change.
- Source-row clauses and occurrences: listed under Source sentence to row.
- Promised field labels: `name: craft-research`, `index:`, `Consumed by`, `Drift`, `Retire when`.
- Changed-function callers: none; no Go function changes.
- Copy survival: CR14's Forbid rows red a Research bullet that keeps any of the three old fragments.

**Reader sweep.**

- The skills index in the reference file reads every skill's `index:` line.
- The mirror check reads every `.claude/skills` entry.
- The prose-budget check reads the glob row in the project profile.
- The anchor check reads the registry and the fixture copies of the two commands.
- The consumer payload lists `.agents/skills` whole, so no payload edit is needed.
- Shipped-surface claim words: none added.
