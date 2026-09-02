---
description: Three-axis semantic review of a branch diff — Standards, Spec, and Coverage. Use after implementation and before the final landing. Advisory, not authoritative.
---

# /bench-review-implementation — the check the gate can't run

## Entry orientation

This is the semantic review phase. It reviews the branch diff against three
separate axes: documented standards, the approved spec, and coverage gaps. It
produces findings the gate cannot see. It claims no authority over done-ness.
A spec-backed review runs from the retained integration source and opens with
`bench preflight review` in explicit-base mode; a red preflight stops the phase.

## Exit handoff

Report the Standards, Spec, and Coverage findings in separate sections. State
the count and the worst issue for each axis. Report two totals apart: the raw
finding count per axis, and the de-duplicated repair-target count after you
collapse findings that name the same fix. Volume and repair work are different
numbers; never report one number where the reviewer asks for the other.

Accepted findings become slim repair tickets with an advisory `Writes:` note,
and they return to `/bench-implement-spec`. Findings that need a later fix pass
use the pickup-file route in step 6. A clean review proceeds to `/bench-final-check`.

The gate is deterministic: it runs the phase table the project profile declares,
and nothing else. Review supplies the semantic judgment that phase table cannot
perform. The gate decides done; the reviewer decides whether a green change ships.

## Review modes

Initial review blocks on the full `frozen-base..reviewed-tip` diff across
Standards, Spec, and Coverage. It is the discovery pass described below.

After you repair accepted findings, a repair-scoped re-review takes the
accepted repair predicates and the prior reviewed tip. It uses the full
`frozen-base..current-tip` diff only as context. Its blocking scope is the
accepted repair predicates plus changes after the prior reviewed tip. This is
checked for repair-induced Standards, Spec, and Coverage problems. A finding
outside both is a non-blocking follow-on and cannot reopen the phase.

The coordinator writes one repair ticket before the repair-scoped re-review,
when accepted repairs amend the coverage map. The ticket records the
accepted repairs, and it cites each amended row in `Covers:`.

A diff that changes kit guidance takes a standing cross-harness
falsification pass. The kit-guidance set is any file under `.agents/` or
the file `.bench/BENCH.md`. Each falsification finding takes one explicit
outcome of accept, merge, or dismiss. An accepted falsification
finding joins the review findings and takes the repair-routing disposition.

Landing proceeds when the repair-scoped result is clean. A repair-induced fix
triggers a new check. It stays scoped to that predicate and repair delta;
it never restarts initial discovery over the original range.

## Process

1. **Pin the diff.** Pull the whole base-relative context with
   `bench diff --full` in explicit-base mode. Supply the same frozen base that
   preflight used. Record the complete reported base and source tip. A dirty
   source, a moved tip, or a pair that differs from preflight stops the review.
   When work already landed, use `bench diff --full --commit <sha>` for that
   historical commit.

2. **Find the sources.** The spec source is `specs/<feature>/spec.md` for this
   work, or the path I give you. The standards sources are `AGENTS.md` and
   `.bench/BENCH.md`, the working agreement and the shared platform rules;
   `CLAUDE.md` holds import pointers only. Also use `projects/<name>.md` and
   any `CONTRIBUTING` or conventions docs in the repo.

3. **Run the blast before axis dispatch.** Run
   `bench consumers --changed --base <b> --source-tip <t> --full` over the same
   frozen pair that preflight and `bench diff` pinned. Attach the returned
   tables and their citation row to the review record. Walk the `touched=false`
   rows first, because a consumer outside the diff's file set is the FT210 class
   an unlisted-consumer miss hides. Then hand each axis the table.

   - Walk a `blast_deleted` row as a deletion whose consumers the tip already edited.
   - A blast refusal stops the review, as a red preflight does.

4. **Spawn the axes in parallel sub-agents.** This isolation keeps one axis's
   derivation from polluting another's context, and stops one axis from seeding
   another's findings. Spawn one delegate per axis — Standards, Spec, and the
   Coverage axis — each under ~400 words. Charge and verify each delegate per
   the `craft-delegate` skill; these are read-only delegations. Each delegate
   re-derives its own facts from its primary source before it compares the
   candidate — see `craft-review`. Give each delegate the diff, the sources for
   its axis, and its charge from the `craft-review` skill
   (`.agents/skills/bench-craft-review/SKILL.md`).

   That skill is the one source
   for what each axis hunts and what a finding must cite; do not restate the
   charges here. A finding cites what its axis read now, not what it recalls.
   A universal claim answers to that citation standard: cite the enumeration, or name itself a sample.
   Procedural inputs per delegate: give Standards the docs from step 2. Give
   Spec the spec file. When the spec carries an acceptance coverage map, also
   give it the rows from `bench coverage <spec>`. The delegate must audit every
   mapped behavior there; that is part of its charge.

   Give Coverage the existing tests and, when one exists, the profile's
   hostile-input checklist.

   If there is no spec, skip the Spec axis and say so. The Coverage axis still
   runs; it needs only the diff and the existing tests.

   If this harness forbids unsolicited sub-agents, run the same axes inline.
   State that fallback in the exit handoff. Keep the same charges and citation
   standard.

5. **Aggregate, don't merge.** Report under `## Standards`, `## Spec`, and
   `## Coverage` headings, and keep the findings separate. Do not rerank across
   axes. Do not pick a single winner. Code can pass one axis and fail another:

   - the right thing, built with the wrong conventions
   - clean conventions, applied to the wrong thing
   - correct on the happy path, but open on the edges

   The separation is the point: when you merge findings, one axis can mask
   another. End with a per-axis count and the worst issue within each axis.
   Give every finding exactly one disposition:

   - `no-op` — the candidate or a cited source refutes the concern; no repair
     target remains
   - `auto-fix` — a deterministic hard rule or an exact spec predicate can be
     repaired inside already-approved scope
   - `ask-user` — the finding needs judgment, scope, authority, or an oracle
     change

   A disposition is a repair-routing label; it is not permission for this
   read-only phase to make the edit itself.

6. **Write and commit the pickup state, in that order, before repair begins.**
   The actionable findings that need a later fix pass go in
   `reviews/<spec-slug>.md`. Keep one section per axis: `## Standards`,
   `## Spec`, and `## Coverage`. Each section carries its
   finding count, its worst issue, and every actionable finding, with its
   disposition and the file or doc citation its axis supplied.
   Keep all three headings, even when only one axis has findings. Replace a
   stale artifact; do not append to it.

   Commit the artifact in the same session that writes it, before any repair edit lands.
   This also applies when another harness picks up this review mid-flight and
   returns findings. Capture and commit its findings the same way, before you
   touch the fix.

   A clean review writes no artifact. So does a review where the reviewer
   accepts every residual risk. The `reviews/` directory means "there is fix
   work to do", not "a review happened". Never commit an empty `reviews/`
   directory or a `.gitkeep`. A no-spec review stays chat-only, unless the
   reviewer supplies an explicit slug. Without a durable spec, an invented
   artifact name would create a second source of feature identity.

   The ordinary artifact is transient pickup state, not a review log. The
   `/bench-implement-spec` session that resolves the findings deletes it in
   the same green fix commit that closes them, so resolved findings cannot
   resurface.

7. **Hand off, don't repair.** This phase makes no fixes and runs no gate.
   Accepted findings become slim repair tickets carrying an advisory `Writes:`
   note and return to `/bench-implement-spec` on the same integration source.
   A spec amendment commits to that same source on the finding cadence. The
   landing publishes the source's spec bytes, so an amendment never routes
   through a hand commit on the destination. A clean review hands its frozen
   base and reviewed tip to `bench worktree land`; `/bench-final-check`
   reports that landing's oracle.
