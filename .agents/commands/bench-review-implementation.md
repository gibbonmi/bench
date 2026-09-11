---
description: Three-axis semantic review of a frozen implementation-chunk diff — Standards, Spec, and Coverage. Use after each planned chunk. Advisory, not authoritative.
---

# /bench-review-implementation — the check the gate can't run

## Entry orientation

This is the semantic review phase. It reviews the branch diff against three
separate axes: documented standards, the approved spec, and coverage gaps. It
produces findings the gate cannot see. It claims no authority over done-ness.
A spec-backed review runs from the retained integration source after each planned chunk. It opens with `bench preflight review` in explicit-base mode; a red preflight stops the phase.

## Exit handoff

Report the Standards, Spec, and Coverage findings in separate sections. State
the count and the worst issue for each axis. Report two totals apart: the raw
finding count per axis, and the de-duplicated repair-target count after you
collapse findings that name the same fix. Volume and repair work are different
numbers; never report one number where the reviewer asks for the other.

Accepted findings return to the retained `/bench-implement-spec` session. Findings that need a fix pass use the pickup-file route in step 6. A clean review returns to the next chunk, or to final reconciliation after the last chunk.

The gate is deterministic: it runs the phase table the project profile declares,
and nothing else. Review supplies the semantic judgment that phase table cannot
perform. The gate decides done; the reviewer decides whether a green change ships.

## Review modes

Each planned chunk takes one review across Standards, Spec, and Coverage. The axes read the whole approved spec and focus on the frozen `chunk-base..chunk-tip` delta.

After the retained author repairs accepted findings, current repair coverage closes those predicates. Repeat delegated review only for a later semantic delta or a cross-chunk concern that invalidates prior evidence. A repeated review uses the full chunk diff as context and blocks only on that later delta or named concern.

The coordinator writes one repair ticket when accepted repairs amend the coverage map. The ticket records the accepted repairs, and it cites each amended row in `Covers:`.

The coordinator records every dogfood run in the spec before repair coverage closes. An unrecorded run is a blocking finding.

A diff that changes kit guidance takes a standing cross-harness
falsification pass. The kit-guidance set is any file under `.agents/` or
the file `.bench/BENCH.md`. Each falsification finding takes one explicit
outcome of accept, merge, or dismiss. An accepted falsification
finding joins the review findings and takes the repair-routing disposition.

The successor chunk starts only after findings and repair coverage close. After the last chunk, the retained author reconciles overall acceptance and integration before landing.

## Process

1. **Pin the diff from the prepared evidence.**
   Collect the shared evidence once with
   `bench preflight review <spec> --charge --base <b> --source-tip <t> --full`.
   Supply the same frozen base that preflight used. The command returns one
   charge row for each axis. It also returns a `shared_evidence` table with one
   identity each for diff, consumers, and coverage.

   Read the base-relative diff facts from that prepared evidence. Do not collect
   the same pair a second time with `bench diff --full`. Record the complete
   reported base and source tip. A dirty source, a moved tip, or a pair that
   differs from preflight stops the review.

   A historical review keeps `bench diff --full --commit <sha>` for the landed
   commit. A spec-less review keeps `bench diff --full` in explicit-base mode.

   The first chunk base is the `main` tip merged into the source. Each later chunk base is the accepted predecessor tip, so the range holds only that chunk's delta.

2. **Find the sources.** The spec source is `specs/<feature>/spec.md` for this
   work, or the path I give you. The standards sources are `AGENTS.md` and
   `.bench/BENCH.md`, the working agreement and the shared platform rules;
   `CLAUDE.md` holds import pointers only. Also use `projects/<name>.md` and
   any `CONTRIBUTING` or conventions docs in the repo.

3. **Walk the blast before axis dispatch.**
   Take the consumer tables and their citation row from the shared evidence's
   `consumers` identity. Do not run a second `bench consumers --changed`
   collection for the same frozen pair. Attach those tables and that citation row
   to the review record.

   Walk the `touched=false` rows first, because a consumer outside the diff's
   file set is the FT210 class an unlisted-consumer miss hides. Then hand each
   axis the same shared evidence.

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

   Resolve every axis through `craft-line`'s conditional review line from the implementation model.

   An authorized review dispatches every prepared axis through the native agent
   surface. It asks no second approval turn inside that authorization. Each axis
   keeps its own context, its own isolated read-only venue, and its own
   independent source derivation.

   Collect every axis return before you accept a finding. A missing or failed
   axis return leaves the review incomplete. It is never a clean finding set.

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

   Native availability is a fact of the active session. A compiled harness
   record is not runtime proof of a native agent surface.

   When the native tool is unavailable, or when the reviewer prohibits
   delegation, preserve the prepared charges and stop dispatch. Then emit a
   capable-harness handoff. The handoff names the repository, the assignment,
   the frozen pair, the charge inputs, the destination harness, and that
   harness's exact native continuation command. Substitute no same-family CLI
   launcher, and collect no axis into coordinator context. The standing
   cross-harness falsification pass keeps its separate route and its existing
   trigger.

   A historical review and a spec-less review keep their existing preparation
   entry points under this same native authority rule.

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

7. **Hand off, don't repair.** This phase makes no fixes and runs no gate. Accepted findings return to `/bench-implement-spec` on the same integration source. A spec amendment commits to that source on the finding cadence.

   A clean chunk review hands its frozen pair back to the retained author. The author starts the successor or performs final reconciliation. Only the reconciled final source proceeds to `bench worktree land`; `/bench-final-check` reports that landing's oracle.

   The landing base is the `main` tip merged before the first chunk. `bench worktree land --base` takes that `main` tip, not a later chunk base.
