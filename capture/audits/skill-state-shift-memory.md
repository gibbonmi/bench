# Assessment — SKILL.state and shift cross-iteration memory (2026-08-30)

Decision source for the FT231 shift-memory case. The paper is SKILL.state:
Scalable Long-Horizon Agent Skills (arXiv 2608.26263v1). The tree was read at
`b9fc2752`. The paper was read through a fetched summary of its HTML version,
not the PDF.

## Result

EXPERIMENT FIRST, with one precondition. No `bench shift` has run in this
repository. No ref carries a `shift: iteration` commit, the intent ledger holds
zero entries, and no `.bench/done.sh` exists. Every candidate problem is
unobserved. The first deliverable is a baseline, not a schema.

## Doctrine under test

History is evidence. Current state is an interface. The two must not merge.

| Surface | Role |
|---|---|
| repository tree and branch log | world state |
| gate | truth about whether a mutation counts |
| bounded shift state | what the next fresh iteration must remember |
| events, transcripts, gate runs | evidence, never default context |

## Current implementation

`internal/shift/loop.go` validates the objective, persists an intent entry,
and acquires a pooled worktree. It then branches `bench/shift-<ts>` and writes
`.bench-objective` and an empty `.bench-notes.md`. Each iteration runs the
adapter as a fresh process with a constant prompt on stdin, then runs the
gate. On green it stages the touched paths with scratch excluded, commits, and
tests `.bench/done.sh`. A red gate retains and locks the worktree. A green completion
deletes the notes at teardown.

Cross-iteration memory is the tree, the branch log, the objective, and the
whole accumulated notes file. The prompt carries no runtime-derived
observation. The default horizon is 12 iterations; the maximum is 100.

## Gap

The outer loop is not a gap. Three specific gaps remain:

1. No latest observation. The paper's fresh call receives the runtime's
   observation. Bench's prompt is a compile-time constant. The model must
   re-derive the iteration index, the last gate verdict, and the last commit
   from git, or trust its own prior prose.
2. Unbounded, unvalidated, single-author memory. The notes are append-only
   model prose. Nothing bounds them, nothing validates them, and nothing
   separates runtime facts from model claims. A false "gate green" line
   propagates unchecked.
3. Memory destroyed at green. A completed shift keeps no evidence of what its
   iterations learned. This blocks the measurement the decision needs.

## Candidate design

This design exists for criticism. The experiment decides adoption.

| Facet | Choice |
|---|---|
| Owner | `internal/shift` only; no generic state subsystem |
| Shape | One scratch file, git-ignored and scratch-excluded. A `runtime` section that Bench writes: iteration, cap, branch, base, head, committed count, last gate verdict and first red, last touched paths, and a state-rejected reason. An `agent` section that the model writes, bounded: `next_step` (one line), `findings` (at most 8), `open_questions` (at most 4). No constraints field; the objective and the spec own constraints. |
| Transition | Replace, not patch. The model rewrites the `agent` section in full. Bench rewrites `runtime` from git and the gate record before each adapter run. |
| Validation | After the adapter exits and before the gate: decode, check types and bounds, refuse any `runtime` key. On an invalid write, keep the prior `agent` section, set the rejected reason, and continue. The gate grades the tree as today. |
| Recovery | The file lives in the worktree, so a retained red worktree keeps it. If the file is missing or corrupt at iteration start, regenerate `runtime` from git with an empty `agent`. At green teardown, copy the final state into the shift record before the scratch cleanup. |
| Projection | The prompt renders the state as a short block. The notes file goes away rather than survive as a competing source. |
| Commit | Never committed. Evidence goes to the record, not the branch. |

The verdict rule is precise. An invalid state never makes green code red. A
valid state never makes red code count. Only the rejected reason changes, and
the next iteration sees it.

## Risks

- Second source of truth. The `runtime` section duplicates git and the gate
  record. It is a projection regenerated every iteration, never read back as
  truth.
- Dumping ground. Three bounded agent fields is the ceiling. A constraints or
  history field is the first sign of drift.
- False red. Do not retry on an invalid state. Record it and continue.
- State and repo divergence. A finding can go stale after a later commit. The
  last-touched runtime field is the mitigation, not a validator.
- Loss of local reasoning. Not a risk at this boundary. An iteration stays a
  whole harness session. Finer reset boundaries stay out.
- Token savings without task gain. The paper's 25x token ratio is against a
  full-transcript baseline that Bench never had. Bench's expected saving is the
  notes' linear growth, which is small at 12 iterations.
- Closed decision. The 2026-08 capability audit closed L-04: no new canonical
  work record, and the handoff's `## State` gets bounded (action A8). The shift
  state must not become that record.

## Experiment

An FT231 case with three arms, so the result separates the cheap change from
the schema.

| Arm | Memory contract |
|---|---|
| A | Current shift: constant prompt, append-only notes |
| B | Constant prompt plus the rendered `runtime` block; notes unchanged |
| C | `runtime` block plus bounded `agent` state; no notes |

Fixed: task, revision, model, effort, harness, adapter, iteration cap, wall cap,
gate, worktree pool, and done predicate. Tasks run at three horizons: near 3,
12, and 40 iterations. Bench owns iterations, commits, gate verdicts and first
reds, elapsed time, and prompt and state size per iteration. The harness owns
tokens, turns, tool calls, and files read, through FT274. Two shift measures
join: paths re-read that a prior iteration read, and recovery after a kill at
iteration k.

Decision rule: Arm C ships only if it beats Arm B on task success or recovery
at the long horizon at non-inferior cost. Arm B ships alone if it beats Arm A,
because it is derived state with no model-authored schema. Arm A baselines
with Bench-owned measures come first, before FT274 lands.

## Paper facts used

At a 100-step horizon the paper reports 0.78 accuracy at about 2,500 prompt
tokens per call against 0.63 at about 62,000 for a full-transcript baseline.
At a 10-step horizon both reach 1.00. Sliding-window and LLMLingua compression
reach 0.18 and 0.22 at 100 steps. Its schema is static per domain, with five
fields in one benchmark. Its validation checks schema and types and retries on
an invalid patch. It states three failure modes:

- a state structure that cannot be predefined,
- a state that depends on an unrecognized earlier observation,
- a task that targets the history itself.

## Roadmap disposition

| Row | Disposition |
|---|---|
| FT71 | amend: the shift record retains per-iteration memory before the scratch cleanup |
| FT162 | amend: the handoff is a projection plus one bounded reviewer section, per audit A8 |
| FT204 | unchanged: it reads evidence and is not execution memory |
| FT231 | amend: the three-arm shift-memory case at three horizons |
| new row | none until the measurement |

## Answers

1. Observed problem: none; no shift has run here.
2. Measurable: yes, once FT274 lands; Bench-owned measures cover half today.
3. Necessary notes content: the next step and a few repo facts that are not
   cheap to re-derive.
4. Derivable: iteration, cap, branch, base, HEAD, committed count, gate verdict
   and first red, touched paths, objective.
5. Smallest semantic state: `next_step`, `findings`, `open_questions`, bounded.
6. Recovery gain: plausible, unmeasured; an arm measure, not an assumption.
7. Owner: `internal/shift`.
8. Invalid state and green code: the rejection is recorded and never touches
   the verdict.
9. Notes as projection: no; drop the file and render the state into the prompt.
10. Existing owners: FT231, FT71, FT162, FT204. No new row.
