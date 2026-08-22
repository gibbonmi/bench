# Craft research — research findings

Three read-only research delegations, 2026-08-02. The coordinator re-read the
primary upstream contract and spot-checked the local claims cited below.

## R1 — current ownership and integration seams

`craft-research` fits Bench as a model-invoked skill, not a reviewer-invoked
phase: every `craft-*` skill is model-invoked, while `/bench-*` adapters remain
explicit phase entry points (`.agents/skills/bench-craft-skills/SKILL.md:13-26`,
`.bench/BENCH-reference.md:61-74`). The portable skill path is
`.agents/skills/bench-craft-research/SKILL.md`; the generated skills index and
the kit checkout's `.claude/skills/` symlink are derived integration surfaces
(`.bench/BENCH-reference.md:36-59`, `projects/benchkit.md:485-492`).

The narrow new owner is factual-reading execution: split a question into
independent source questions, coordinate read-only returns, and verify claims
and their composition. It also preserves conflicts and unknowns, and
synthesizes the selected output. Existing owners remain authoritative:

- `craft-delegate` owns complete charges, one coherent question per delegate,
  read-only isolation, and coordinator verification
  (`.agents/skills/bench-craft-delegate/SKILL.md:38-77,141-145,185-221`).
- `craft-line` owns model, effort, iteration cap, fan-out declaration, and
  escalation (`.agents/skills/bench-craft-line/SKILL.md:13-69,120-139`).
- `/bench-shape-idea` owns Research ticket state, map readiness, source-object
  registration, and reviewer decisions
  (`.agents/commands/bench-shape-idea.md:26-83`).
- `/bench-write-spec` owns source revalidation, compilation and movement of map
  assets, retirement, and engineering decisions
  (`.agents/commands/bench-write-spec.md:23-64,111-125`).
- `internal/maps` owns the only structured `Path|URL`, `Supports`, `Drift`
  source-record grammar (`internal/maps/validation.go:167-235`).

`/bench-assess` is the existing phase-specific realization of most of this
execution contract: it fans out six fixed read-only area sweeps, re-verifies
every load-bearing claim, records unverifiable claims as unknowns, and writes one
cited `ASSESSMENT.md` with verification notes
(`.agents/commands/bench-assess.md:26-75`). It should consume the shared research
discipline while retaining its fixed areas, line choices, severity and backlog
schema, predecessor reconciliation, and replace-in-place lifecycle. Formal
`/bench-review-implementation` remains separate: its fixed three axes and
finding/receipt lifecycle are owned by `craft-review` and its phase
(`.agents/commands/bench-review-implementation.md:44-102`).

The primary duplication risks are restating delegation mechanics or line
routing, inventing a second source manifest, or copying map asset lifecycle into
the skill. The byte/wire compatibility-probe rule is already duplicated between
shaping and `craft-spec`; adding a third copy would deepen that defect
(`.agents/commands/bench-shape-idea.md:50-52`,
`.agents/skills/bench-craft-spec/SKILL.md:51-62`).

## R2 — upstream contract and the useful comparator

Matt Pocock's current skill has four essential moves. It fires when work becomes
factual reading legwork, and delegates that reading to a background agent. It
also restricts evidence to primary sources, and leaves one cited Markdown file
in the repo's existing notes location. It specifies no fan-out and no coordinator
re-verification. It also specifies no citation syntax, drift fields, fixed
path, or completion rubric
([current skill](https://raw.githubusercontent.com/mattpocock/skills/main/skills/engineering/research/SKILL.md),
[first-party documentation](https://raw.githubusercontent.com/mattpocock/skills/main/docs/engineering/research.md),
both retrieved 2026-08-02).
Its durable boundary is that reading is delegated while judgment stays with the
caller; “background agent” and slash-command syntax are harness particulars.

Anthropic's production Research system is the directly relevant primary-source
fan-out comparator. It uses an orchestrator-worker pattern: the lead splits a
question into specialized independent searches, and parallel workers return
findings. The lead then synthesizes, decides whether another round is needed,
and runs a citation pass that attaches sources. It reports that multi-agent
research works best for breadth-first independent directions, and costs
roughly fifteen times a chat interaction in tokens. It also reports that
multi-agent research is a poor fit for dependency-heavy or shared-context
work.

Its learned safeguards are complete non-overlapping charges and effort scaled
to question complexity. They also include evaluation of accuracy, citation
accuracy, completeness, source quality, and tool efficiency
([Anthropic engineering report](https://www.anthropic.com/engineering/multi-agent-research-system),
retrieved 2026-08-02).

This supports conditional fan-out and coordinator synthesis. It does not support
unconditional fan-out, a portable promise that work is literally asynchronous,
or importing Anthropic's production worker counts as Bench policy.

## R3 — committed Bench evidence

FT135 is the one current committed example that explicitly records a research
fan-out: six read-only delegations each owned a concrete factual question, and
the coordinator assembled one R1–R6 asset while re-verifying every claim
(`specs/pre-push-guard-visibility/decisions/assets/pre-push-guard-visibility-research.md:1-6`,
`specs/pre-push-guard-visibility/decisions/pre-push-guard-visibility.md:374-378`).
The useful unit was an independent question or ownership seam, not a number of
files. The map stayed compact, retained one source record with named invalidation
triggers, and the asset moved with the map when the spec compiled.

FT135 also supplies the concrete failure case. One answer found that
`installGitHook` had no call site, but later answers still composed a currency
formula around that dead path. Spec authoring caught the conflict only when it
re-read the research against the current tree
(`specs/pre-push-guard-visibility/spec.md:5-14,108-127`). Claim-by-claim citation
checks are therefore insufficient: synthesis must also check relationships
between returns and identify contradictions.

Other resolved Research tickets show when fan-out is the wrong tool:

- A single bounded inventory of 39 gate fixtures exposed ownership exceptions
  that an arbitrary per-file split could hide
  (`decisions/gate-pipeline.md:246-278`,
  `decisions/assets/gate-pipeline-fixture-inventory.md:1-78`).
- One serialized timing probe falsified a proposed parallelization direction;
  concurrent probes would have distorted the host-load evidence
  (`decisions/cost-follows-project-size.md:33-72`).

The evidence-backed execution shape is therefore: fan out only independent
reading questions, give each one a bounded return schema, and synthesize one
durable asset. It must also verify both claims and cross-return composition,
and preserve contradictions, unknowns, and drift. It must serialize
measurement or compatibility probes whose results interfere with one another.

## Remaining decisions

The sources do not decide the activation breadth, fan-out threshold and cap,
or durable artifact contract. They also do not decide the judgment boundary,
compatibility-probe owner, or the exact clauses migrated from current phase
guidance. Those remain reviewer-owned tickets #4 through #8 in
`decisions/craft-research.md`.
