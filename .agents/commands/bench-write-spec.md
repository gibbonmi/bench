---
description: Turn a reviewed decision source into a build spec with user stories, engineering seams, and testing decisions. Use before any build past the lighter-path threshold in .bench/BENCH.md.
---

# /bench-write-spec — lock the seams before the loop runs

## Entry orientation

This is the feature-spec phase: it charges `bench-craft-spec` on one authorized
decision source and produces `specs/<feature>/spec.md` plus its tickets.
`craft-spec` owns the process — explore, seams, synthesize (no interview) —
and the discipline: stories, acceptance coverage map, red signal, edge
inventory, **Won't handle**, hostile-input checklist, fences.

## Exit handoff

Close by showing the approval table. The spec carries `Status: staged`
(staged → implemented at the green gate → promote-then-delete on merge). Stop
for sign-off, then recommend a fresh mid-tier build session on one retained
integration source; review its frozen base and tip, and hand the accepted source
to `bench worktree land`.

## Entry contract

Accept exactly one of three decision sources: a ready compiled map, the
reviewer-confirmed current conversation, or a named reviewed artifact. No
unnamed memory, unreviewed note, or fourth override authorizes a draft.
Top-level `decisions/` holds pre-spec working maps; compiled maps live under
`specs/<slug>/decisions/`.

- **Ready compiled map.** Validate it, move (do not copy) the source map and
  any map-owned assets from top-level `decisions/` into
  `specs/<slug>/decisions/` and update every reference to the moved paths in
  the same green change. A re-run reads the already-compiled spec-local map; it
  never recreates a top-level copy.
- **Reviewer-confirmed current conversation.** Confirm every load-bearing
  product or scope fork is closed here; record the date; do not manufacture a
  map to restate it.
- **Named reviewed artifact.** Name the exact approved artifact and verify it
  holds settled decisions, not unresolved prompts.

Record exactly one `Decision source:` line in the spec. For a map-backed
source, re-read and re-verify every structured `## Sources`
entry before choosing engineering seams, disclosing anything that could not be
re-read, and consume them in place without copying a research manifest into
the spec. Ask at most two late clarification
questions, one at a time, each with a recommended answer; route a dependency
tree or multi-session fog to `$bench-shape-idea`.

## Ownership

Spec authoring owns engineering seams, deep-versus-thin design, tests,
acceptance coverage, hostile-input attachment, and gate attachment; shaping
sources constrain behavior, scope, compatibility, or a reviewer-chosen seam.

## Who runs this phase

The session holding the decision source authors the spec and tickets at whatever
tier it runs. The fresh implementation session starts only after ticket
approval. Top tier remains an explicit escalation under `craft-line`.

## Process

1. **Author.** Charge `bench-craft-spec` (and `bench-craft-domain` for terms)
   on the decision source; read the enforcement surface before locking rows
   that touch anchored clauses, canary fixtures, or counted substrings. Give
   each story group its resolved model and effort from `craft-line`. Write
   `specs/<feature>/spec.md` from `craft-spec`'s template and run
   `bench coverage --check`. The stale-command-reference sweep remains
   fail-closed across staged specs. When no hostile-input checklist class covers
   a surface, quarry the seams library and
   propose a tuned profile addition. Apply `craft-spec`'s named
   `Bootstrap authority before execution` rule.

2. **Retire superseded work by promotion then deletion.** Do not leave a
   `Superseded by` marker: promote durable decisions, delete the old spec under
   a `spec-retire: <name>` commit, repair references. The same
   promote-then-delete commit removes the spec's `ROADMAP.md` row. Whole-folder
   retirement removes the compiled maps and map-owned assets, plus tickets.

3. **Review rounds before sign-off.** Two rounds per artifact, spec then
   tickets, with the reviewer-named model: round 1 hunts; round 2 verifies
   only whether round-1 findings were fixed. Partials left after round 2 are
   folded by the author and named in the verification log. Verdicts stay
   advisory; reviewer sign-off remains the hard stop. A finding blocks only
   when it changes observable behavior, an ownership fence, or the ticket
   graph; a round returning only prose or accounting findings is the
   acceptance round — fold those fixes into the acceptance instead of another
   round. A revision may not add a promise beyond the decision source unless
   a blocking finding demands it; the review flags an unflagged addition for
   removal rather than demanding rows for it.

   `--reviewer <tier-or-model> [effort]` overrides both loop delegates: a tier
   resolves through the invoking harness's own `.bench/lines.env` column
   (`--reviewer mid xhigh` under Codex resolves `BENCH_CODEX_MID` at xhigh), a
   model id must already be bound in `.bench/lines.env` and an unbound id is
   refused, an own-family id runs through the native agent surface, and a
   cross-family id uses the exact recipe in
   `.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md`.

   Round 1 asks the falsification questions: would the cheapest wrong
   implementation pass, does every source behavior have a red-capable row, does
   every line match cached routing, does any behavior, red signal, or decision
   answer name an outcome family instead of an exact predicate, are the source
   and observed reds sound even when the source is same-session, conflicting, or
   mostly not observed, and — where the stories partition into disjoint package
   or fence sets — could a narrower capability ship
   on its own gate? Apply `craft-spec`'s named
   `Bootstrap authority before execution` rule. The degenerate standard is
   the cheapest plausible wrong implementation — a degenerate that needs
   deliberate contrivance is the build's mutation-probe target, never a new
   spec row.

4. **Slice the implementation tickets.** After loop 1 accepts, charge
   `craft-tickets` and write the breakdown under `specs/<slug>/tickets/`. Carry
   its numbered title, `Blocked by:`, and delivered outcome list into the
   approval table so the spec and tickets receive one sign-off, then run loop 2
   on that breakdown. `craft-tickets`' granularity/edges/merge-split quiz is the
   approval round.

    At close, write `Verification log: spec <n> + tickets <m> iteration(s) to
    accept — <note>` into the spec. When either loop takes more than one iteration
    to accept, append one `capture/learnings.md` entry naming the stage that
    missed, what review caught, why it was missed, and the proposed rule change.
