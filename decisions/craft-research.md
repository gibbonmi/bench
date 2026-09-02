# Craft research

Status: ready

## Destination

Define a model-invoked `craft-research` skill. It fans out primary-source
reading to read-only subagents when the questions are mutually independent,
then returns one verified, cited Markdown artifact. Bench shaping,
specification, diagnosis, assessment, and implementation work can consume
that artifact without duplicating research policy.

## #1: Where does research policy live in Bench today, and what must remain single-sourced?

Blocked by: none
Type: Research

### Question

Trace every current research instruction, artifact convention, delegation rule,
and conformance surface in the kit. Identify the exact ownership seams a new
model-invoked skill must compose with rather than restate. Asset: a short cited
inventory of current owners and duplication hazards.

### Answer

Resolved 2026-08-02. `craft-research` should own portable factual-reading
execution: independent source questions, read-only return coordination, claim
and cross-return verification, conflict and unknown preservation, and synthesis.
It must compose with, not restate, `craft-delegate`, `craft-line`, and map
Sources. It must also compose with, not restate, shaping/spec lifecycle and
phase-specific assessment, diagnosis, implementation, and review ownership.
Detail: research asset R1.

## #2: What contract does the upstream research skill actually establish?

Blocked by: none
Type: Research

### Question

Read Matt Pocock's current research skill and its first-party documentation,
plus only directly relevant primary-source skill implementations when they test
a materially different design. Separate its essential contract from incidental
harness behavior. Asset: a short cited comparison that can support #4 through
#7.

### Answer

Resolved 2026-08-02. The upstream contract is intentionally small: model-invoked
factual reading, primary sources only, delegated legwork, and one cited Markdown
file in the repository's existing notes location. It does not specify fan-out,
coordinator verification, a citation schema, drift, or fixed placement.
Anthropic's primary-source multi-agent report supports conditional fan-out for
independent breadth and warns against dependency-heavy work. Detail: research
asset R2.

## #3: What has made Bench's existing research fan-outs trustworthy or costly?

Blocked by: none
Type: Research

### Question

Inspect representative resolved Research tickets and their assets. Identify
repeatable strengths, observed failure modes, verification burden, useful
fan-out shapes, and artifact-placement conventions. Asset: a short cited
evidence summary based on committed examples rather than policy alone.

### Answer

Resolved 2026-08-02. FT135 demonstrates six independent factual delegations
synthesized into one coordinator-owned asset. It also demonstrates the key
failure: individually verified claims composed into a false later conclusion.
Other tickets show that corpus-wide inventories and load-sensitive probes often
need one serialized owner. Fan out by independent question, then verify both
claims and relationships between returns. Detail: research asset R3.

## #4: When should `craft-research` fire?

Blocked by: #1, #2, #3
Type: Grill

### Question

Should it be a reach-for-it-anytime model-invoked skill for factual reading
legwork, with `/bench-shape-idea` using it for Research decision tickets, or
remain confined to shaping?

### Answer

Reach-for-it-anytime model-invoked guidance (reviewer, 2026-08-02).
`craft-research` fires when work becomes factual reading legwork in shaping,
specification, diagnosis, assessment, or implementation. Each calling phase
retains authority over its decisions, artifacts, and completion contract;
`/bench-shape-idea` invokes the skill to execute Research decision tickets. It
continues to own their state, dependencies, Sources registration, and readiness.
Formal `/bench-review-implementation` is explicitly not a caller: its fixed-axis
judgment, source set, delegation, and persistence contract remain wholly owned by
`craft-review` and the review phase. A separately identified factual question
leaves review and may become a research run under the phase that owns its answer.

## #5: When and how widely should research fan out?

Blocked by: #1, #2, #3
Type: Grill

### Question

Choose the threshold for delegation, the maximum useful parallelism, and the
unit assigned to each subagent. Fan-out should reduce wall clock without
producing an unverifiable pile of overlapping summaries.

### Answer

Adaptive, round-based fan-out (reviewer delegated to the recommended design,
2026-08-02). Start by drawing the factual question graph; it is complete when
every load-bearing claim the destination needs traces to one question node.
Resolve a trivial lookup inline. Give one bounded question to one read-only
delegate when its read-set would displace the coordinator's context or the
coordinator has useful concurrent work. Fan out only when at least two frontier
questions are mutually independent: neither needs another's answer, shared
mutable state, or a concurrency-sensitive measurement.

There is no portable numeric maximum. A round's width is the number of mutually
independent frontier questions, capped by what one synthesis pass can verify.
The coordinator declares the exact number through `craft-line`'s fan-out
clause before dispatch, and reports any widening like a ladder move.
`craft-research` states no tier or effort. Each charge follows `craft-delegate`
and adds one research-specific field: the question's primary-source boundary.

Synthesize a completed round before opening another. A new round cannot
be charged precisely until synthesis names the conflict, unknown, or coverage
gap it resolves. Dependent questions and concurrency-sensitive measurements run
serially. Rounds are bounded by the calling phase's declared iteration cap;
exhausting it stops with residual unknowns rather than opening another round.
Compatibility and mutation probes follow #7, outside research execution.

## #6: What artifact should a research run leave behind?

Blocked by: #1, #2, #3
Type: Grill

### Question

Choose whether the coordinator synthesizes one cited Markdown asset from all
subagent returns, and where that asset lives. Choose what source, drift,
uncertainty, and verification fields are mandatory.

### Answer

One coordinator-authored durable Markdown output per research run, keyed by
topic (reviewer delegated to the recommended design, 2026-08-02). A run
answering several Research tickets records one section per question, and the
map's single Sources entry names every ticket it supports. Delegate returns are
ephemeral inputs; read-only delegates never write the durable output.

A calling phase with a durable output contract uses that artifact and creates
no second research file: `/bench-assess`, for example, synthesizes directly into
`ASSESSMENT.md`. Otherwise colocate the research asset with the artifact that
will consume it. Shaping uses `decisions/assets/<topic>-research.md`, which moves
with its map. When the repository rather than another artifact is the consumer,
use `docs/research/<topic>.md`. Before dispatch, the caller names the destination,
the consuming phase or artifact, and the condition that retires or refreshes it.

The output records the question and scope, and one synthesized section per
question. It cites every load-bearing claim inline, to an exact local path and
line or to a direct primary-source URL. It also records conflicts, residual
unknowns, and a verification record.

An asset outside a phase-owned output also states `Consumed by`, `Drift`, and
`Retire when` metadata. External citations carry the retrieval date; mutable
evidence names its invalidation trigger. The artifact under study and first-party
upstream documentation or APIs are primary sources. Commentary, summaries,
secondary write-ups, and unopened recall are not. A secondary source may locate
evidence but cannot warrant a finding; when no primary source is available,
record an unknown rather than a verified answer.

Beyond `craft-delegate`'s verification duties, the coordinator re-opens every
source supporting a load-bearing conclusion and independently checks every join
between returns. The named failure to prevent is a synthesis whose individual
claims are true but whose combined conclusion is false. Completion requires
every in-scope question to be answered or explicitly unknown, every material
claim to be traceable, and every contradiction to be reconciled or retained.
The durable output must also state what could not be verified. A shaping map
registers only that asset through the existing `Path|URL`, `Supports`, `Drift`
source schema. Neither the asset nor a later spec copies a second structured
manifest.

## #7: Where is the boundary between research and reviewer-owned judgment?

Blocked by: #1, #2, #3
Type: Grill

### Question

Decide what `craft-research` may conclude, and what it must leave as an
explicit decision. Decide whether prototypes or compatibility probes belong
inside the skill or remain separate evidence-producing work.

### Answer

Read-side research only (reviewer, 2026-08-02). `craft-research` may establish
source-backed facts, contradictions, unknowns, and implications. It never owns a
write delegate, a done-claim, a reviewer decision, or a prototype. Anything with
write access or a done-claim routes through `craft-delegate`; research delegation
also composes with that skill for charge, isolation, and verification mechanics.

Compatibility or mutation probes remain separate evidence-producing work under
their calling phase; research may identify the need but does not absorb the
probe. In a decision map, the separate evidence work is a Prototype ticket, and
the Research ticket names it in `Blocked by:`. The research conclusion stays
unverified until the probe returns.

## #8: Which current instructions should become references to `craft-research`?

Blocked by: #4, #5, #6, #7
Type: Grill

### Question

Choose the single-source migration boundary: which research clauses move into
the new skill, and which phase-specific clauses remain at their current owners
as integration instructions.

### Answer

Single-source migration (reviewer delegated to the recommended design,
2026-08-02):

- Add model-invoked `.agents/skills/bench-craft-research/SKILL.md`. It owns the
  factual-reading trigger, question decomposition, primary-source standard,
  adaptive fan-out, coordinator synthesis and verification, artifact contract,
  and residual-unknown posture. It also owns the rule that byte or wire
  compatibility cannot become a verified research claim without a separate
  runnable probe.
- Keep `/bench-shape-idea` as owner of decision-ticket types, dependencies,
  reviewer exchange, map state, asset lifecycle, Sources registration, and
  readiness. Its Research definition and execution paragraphs point to
  `craft-research` instead of restating source, artifact, probe, or delegation
  rules. A required compatibility probe becomes a separate Prototype ticket;
  the Research ticket names it in `Blocked by:` and stays unresolved until that
  evidence returns.
- Leave `craft-delegate` unchanged, including its “fan-out search” frontmatter
  trigger. The overlap is deliberate composition. `craft-research` owns research
  semantics. `craft-delegate` owns every delegate's charge, isolation, and
  verification mechanics, including read-only delegates and anything with write
  access or a done-claim.
- Leave `craft-line` unchanged and point to it for every model, effort,
  iteration-cap, fan-out declaration, and escalation decision. Leave
  `craft-tickets` untouched.
- Replace `craft-spec`'s duplicated map/research compatibility-probe wording
  with a pointer to `craft-research`'s compatibility-evidence rule; retain
  `craft-spec`'s spec-side duty to reject an unsupported compatibility promise.
  Leave `/bench-write-spec`'s source revalidation and artifact lifecycle clauses
  unchanged.
- Make `/bench-assess` a caller of `craft-research` for area fan-out,
  coordinator re-verification, citation quality, and unknown posture. Assessment
  retains its six fixed areas, previous-backlog reconciliation, phase-specific
  line choices, severity grammar, ranked backlog, verification marks, and
  replace-in-place `ASSESSMENT.md` lifecycle.
- Add the generated skills-index and `.claude/skills/` integration surfaces the
  kit requires. Add no phase command, Codex explicit phase adapter, bespoke payload
  row, parser, or source schema.

Shared file with FT164: `.agents/skills/bench-craft-spec/SKILL.md`. FT164 edits
its slicing section and edge walk; this map replaces its compatibility-probe
wording. The sentences and enforcement anchors are distinct, so either order
lands. The builds do not run concurrently in one worktree, and the second
re-reads the landed wording before editing. FT164's edit set stays owned by its
spec. Its current “disjoint” statement is stale, because the shared file makes
the changes non-disjoint even though ordering remains free.

## Not yet specified

## Spec-writer discretion

- Exact skill headings and the compact contrastive example required by
  `craft-skills`, provided the example demonstrates independent fan-out versus a
  dependent question that stays serial.
- Exact research-asset heading names and filename slug normalization, provided
  the decided content and location precedence remain intact.
- Exact citation notation for local and external primary sources, provided a
  cold reader can resolve every material claim to the cited evidence.

## Out of scope

- Changing the four decision-ticket types or the decision-map schema.
- Implementing a general-purpose knowledge base, citation database, or web-search CLI.
- Replacing `craft-delegate`, `craft-line`, or harness-native subagent controls.
- Folding formal `/bench-review-implementation` axis review into generic research.

## Sources

- Path: `decisions/assets/craft-research-research.md`
  Supports: #1 through #3 and the factual premises for #4 through #8. Three read-only research delegations ran 2026-08-02, with upstream sources re-read and local claims spot-checked by the coordinator.
  Drift: re-verify if research, delegation, line-routing, map-source, skill-index, assessment, or artifact-lifecycle guidance changes, or if the cited upstream research contracts move. Re-resolve the asset's line citations before `/bench-write-spec` reads this map if FT164 has landed.
