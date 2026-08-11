---
description: Three-axis semantic review of a branch diff — Standards, Spec, and Coverage. Use after implementation and before the final landing. Advisory, not authoritative.
---

# /bench-review-implementation — the check the gate can't run

## Entry orientation

This is the semantic review phase. It reviews the branch diff against three
separate axes: documented standards, the approved spec, and coverage gaps. It
produces findings the gate cannot see, without claiming authority over done-ness.
A spec-backed review opens with `bench preflight review <slug>` before pinning
the diff; a red preflight stops the phase.

## Exit handoff

Close by reporting Standards, Spec, and Coverage findings separately, with counts
and the worst issue in each axis. Accepted
findings become ownership-fenced repair tickets and return to
`/bench-implement-spec`; findings that need a later fix pass use the pickup-file route
in step 5, and a clean review proceeds to `/bench-final-check`.

The gate is deterministic: tests, types, lint, conformance. It catches regressions
and rule violations. It cannot tell whether you built the *right* thing the *right*
way. `/bench-review-implementation` is the semantic pass that can — and it's advisory: it surfaces
findings for you, it has no authority to call anything done. The gate and you do.

Run it on the branch diff against its
true base, on three axes that stay separate.

## Process

1. **Pin the diff.** Pull the whole base-relative context with
   `bench diff --full`: it prefers the
   branch's recorded pre-shift base and falls back to merge-base with the default
   branch, and its `method:` line says which happened. When that diff is empty
   because the work already landed, use `bench diff --full --commit <sha>`
   to review exactly that landing commit.

2. **Find the sources.** Spec: `specs/<feature>/spec.md` for this work (or the path I
   give you). Standards: `AGENTS.md` and `.bench/BENCH.md` — the working agreement
   and shared platform rules; `CLAUDE.md` is import pointers only — plus
   `projects/<name>.md` and any `CONTRIBUTING`/conventions docs in the repo.

3. **Spawn the axes in parallel sub-agents** (so they don't pollute each other's
   context), one per axis — Standards, Spec, and the Coverage axis — each under
   ~400 words, charged and verified per the `craft-delegate` skill (these are
   read-only delegations). Each delegate takes the diff, the sources for its
   axis, and its charge from the `craft-review` skill
   (`.agents/skills/bench-craft-review/SKILL.md`) — the one source for what each
   axis hunts and what a finding must cite; don't restate the charges here.
   A finding cites what its axis read now, not what it recalls; the bar for a
   universal claim — cite the enumeration or name itself a sample — is that
   citation standard's.
   Procedural inputs per delegate: give Standards the docs from step 2; give
   Spec the spec file plus, when it carries an acceptance coverage map, the rows
   from `bench coverage <spec>` — auditing every mapped behavior there is part
   of its charge; give Coverage the existing tests and the profile's
   hostile-input checklist when one exists.

   If there's no spec, skip the Spec axis and say so; the Coverage axis still
   runs — it needs only the diff and the existing tests.

   If this harness forbids unsolicited sub-agents, run the same axes inline,
   state that fallback in the exit handoff, and keep the same charges and
   citation standard.

4. **Aggregate, don't merge.** Report under `## Standards`, `## Spec`, and
   `## Coverage` headings, findings kept separate. Do not rerank across axes or
   pick a single winner — the separation is the point: code can pass one axis and
   fail another (right thing, wrong conventions; clean conventions, wrong thing;
   correct on the happy path, open on the edges), and merging them lets one mask
   the other. End with a per-axis count and the worst issue within each axis.

5. **Persist the right review state.** The actionable findings that need a later
   fix pass go in `reviews/<spec-slug>.md`. Keep one section per axis: `## Standards`,
   `## Spec`, and `## Coverage`. Each carries its finding count, its worst issue,
   and every actionable finding with the file or doc citation its axis supplied.
   Keep all three headings even when only one axis has findings.
   Replace a stale artifact rather than appending, and commit the artifact in the
   same session that writes it so pickup state is tracked at birth.

   A clean review, or one where the reviewer accepts every residual risk,
   writes no artifact: the `reviews/` directory means "there is fix work to
   do", not "a review happened". Never commit an empty `reviews/` directory or
   a `.gitkeep`. A no-spec review stays chat-only unless the reviewer supplies
   an explicit slug — without a durable spec, inventing an artifact name would
   create a second source of feature identity.

   The ordinary artifact is transient pickup state, not a review log: the
   `/bench-implement-spec` session that resolves the findings deletes it in the
   same green fix commit that closes them, so resolved findings cannot resurface.

6. **Hand off, don't repair.** This phase makes no fixes and runs no gate.
   Accepted findings become
   ownership-fenced repair tickets and return to `/bench-implement-spec`; a
   clean review proceeds to `/bench-final-check` for its fresh oracle run.

## Where it sits

`/bench-review-implementation` is generation-shaping, not enforcement: run it,
read the findings, and decide what to fix. Final-check
runs the gate and commits on green. Review tells you whether it's *good*; the
applicable oracle tells you whether it's *done*; you decide whether it ships.
