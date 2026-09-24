---
description: Three-axis semantic review of a frozen implementation-chunk diff — Standards, Spec, and Coverage. Use after each planned chunk. Advisory, not authoritative.
---

# /bench-review-implementation — the check the gate can't run

## Entry orientation

This is the semantic review phase. It reviews the branch diff against three
separate axes: documented standards, the approved spec, and coverage gaps. It
produces findings the gate cannot see. It claims no authority over done-ness.
A spec-backed review runs from the retained integration source after each planned chunk. It opens with `bench preflight review` in explicit-base mode; a red preflight stops the phase.

Each review round dispatches a fresh session for each axis. The coordinator keeps the manifest, the metadata page, and the current binding, and it gives the evidence identity to each axis.

A narrow axis runs `bench preflight evidence <id> --check-current` one time, to bind the artifact to the assignment and the source pair. The axis reads the delta through one `git diff` of the frozen pair. It reads the spec rows, the tickets, the standards, and the surrounding code with targeted reads. It returns a bounded report with one line for each finding. A resumed axis session keeps every earlier stream in its context, so the phase does not resume an axis for a later round.

This narrow axis shape is provisional. Each narrow round records what each axis read and what it found. One full-retrieval control review of the same diff decides whether the shape becomes the permanent rule.

Review action requires a narrow axis read: act only after this session reads the frozen delta and the targeted sources of its axis. An axis does not retrieve every evidence page. `bench preflight evidence <id> --verify` is the separate artifact-integrity check.

Review action requires available required context: a retrieval receipt, a terminal cursor, or another axis's read never replaces the sources this session reads. Review action requires a current-action binding: act only after `bench preflight evidence <id> --check-current` binds the artifact to the current assignment and source pair. Review action requires reviewer approval, which a generated charge, a staged artifact, or a verified artifact never supplies. Review action requires the complete task supplement, which this session writes. The supplement names the axis, the model, the effort, and the frozen pair, and verified evidence never supplies it.

The coordinator reuses exact available source bytes only after the new manifest verifies their membership, role, and requiredness. A matching body digest alone authorizes no reuse. A fresh axis runs its own retrieval from the trusted evidence identity. A transferred final cursor and another axis's receipt deliver no byte to it.

## Exit handoff

Report the Standards, Spec, and Coverage findings in separate sections. State
the count and the worst issue for each axis. Report two totals apart: the raw
finding count per axis, and the de-duplicated repair-target count after you
collapse findings that name the same fix. Volume and repair work are different
numbers; never report one number where the reviewer asks for the other.

Accepted findings go to a fresh repair author under `.bench/BENCH.md`'s repair rule. Findings that need a fix pass use the pickup-file route in step 6. A clean review returns to the next chunk, or to final reconciliation after the last chunk.

The gate is deterministic: it runs the phase table the project profile declares,
and nothing else. Review supplies the semantic judgment that phase table cannot
perform. The gate decides done; the reviewer decides whether a green change ships.

## Review modes

Before classifying repairs, read [the bounded repair policy](../skills/bench-craft-line/references/bounded-repair-policy.md). Apply its allowance and retained-state rules.
Each planned chunk takes one review across Standards, Spec, and Coverage. The axes read the whole approved spec and focus on the frozen `chunk-base..chunk-tip` delta.

After the fresh repair author repairs accepted findings, current repair coverage closes those predicates. Repeat delegated review only for a later semantic delta or a cross-chunk concern that invalidates prior evidence. A repeated review uses the full chunk diff as context and blocks only on that later delta or named concern.
A chunk that ends on a repair takes one confirming round of all three axes at its final tip.
A confirming round reads only the repair delta, and its charge names the folds to confirm.

The coordinator writes one repair ticket when accepted repairs amend the coverage map. The ticket records the accepted repairs, and it cites each amended row in `Covers:`.

The coordinator records every dogfood run in the spec before repair coverage closes. An unrecorded run is a blocking finding.

A cross-harness falsification pass runs only when the reviewer requests it.
Each falsification finding takes one explicit outcome of accept, merge, or
dismiss. An accepted falsification finding joins the review findings and takes
the repair-routing disposition.

A delegated chunk review starts after every ticket of the chunk reaches the integrated chunk tip. A per-ticket review does not replace that full-chunk review. Each delegated axis excludes the orchestrator and every current and former author of the run.

By default, an authorized review dispatches every prepared axis through the
native agent surface without a second approval. Each axis uses a different
independent session, an isolated read-only venue, and its own source derivation.
A version 2 completion plan can set `execution.review_mode` to `unified` as the
sole exception. One independent session re-derives Standards, Spec, and Coverage
separately, then reports each axis separately.

For each issue or review miss, it states whether the implementation command
contributed. When it did, it names the exact improvement. Otherwise, it states
that no command change is necessary.

Here, findings that prevent progression are unresolved blockers; retain optional advice separately under the policy. The successor chunk starts only after findings and repair coverage close. After the last chunk, the orchestrator reconciles overall acceptance and integration before landing.

## Process

1. **Pin the diff from the prepared evidence.**
   Prepare the shared evidence once with
   `bench preflight review <slug> --charge --base <b> --source-tip <t>`.
   Supply the same frozen base that preflight used. The command publishes one
   immutable artifact and prints its identity. Retrieve every byte with
   `bench preflight evidence <id>` and each exact successor command it prints.
   The metadata names one charge row for each axis, and it binds the diff,
   consumers, and coverage captures to their frozen sources.

   Read the base-relative diff facts from that prepared evidence. Do not collect
   the same pair a second time with `bench diff --full`. Record the complete
   reported base and source tip. A dirty source, a moved tip, or a pair that
   differs from preflight stops the review.

   A historical review keeps `bench diff --full --commit <sha>` for the landed
   commit. A spec-less review keeps `bench diff --full` in explicit-base mode.

   The first chunk base is the `main` tip merged into the source. Each later chunk base is the accepted predecessor tip, so the range holds only that chunk's delta.
   A later plan commit is never a chunk base.

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

4. **Dispatch the review sessions.** In the default mode, spawn the axes in
   parallel sub-agents. This isolation keeps one axis's
   derivation from polluting another's context, and stops one axis from seeding
   another's findings. Spawn one delegate per axis — Standards, Spec, and the
   Coverage axis — each under ~400 words. Charge and verify each delegate per
   the `craft-delegate` skill; these are read-only delegations. Each delegate
   re-derives its own facts from its primary source before it compares the
   candidate — see `craft-review`.

   Give each delegate the diff, the sources for
   its axis, and its charge from the `craft-review` skill
   (`.agents/skills/bench-craft-review/SKILL.md`).

   Resolve every axis through `craft-line`'s conditional review line from the implementation model.

   On one shared tree, only one axis runs tests or probes while the other axes read.
   A probing axis otherwise takes its own worktree.

   Collect every axis return before you accept a finding. A missing or failed
   axis return leaves the review incomplete. It is never a clean finding set.

   In the explicit unified mode, dispatch one independent reviewer with all
   three prepared charges. Require three separate axis returns from that one
   session. Apply the same source derivation, finding, and completion rules to
   each return.

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
   harness's exact native continuation command. The handoff carries the trusted
   expected evidence identity and the exact `bench preflight evidence <id>`
   retrieval command, never the originating checkout path. Substitute no
   same-family CLI launcher, and collect no axis into coordinator context. An
   explicitly requested cross-harness falsification pass keeps its separate
   route.

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
   Each actionable finding line carries its stated confidence.
   Keep all three headings, including axes with zero findings.

   Retain every terminal return in one fenced `bench-review-record` JSON payload.
   The `internal/reviewrecord` types own the schema. Preflight supplies the source and plan digests.
   Record the performer, role, model, effort, frozen base and tip, source, state, and native result.

   Use explicit `unknown` for unavailable model or effort metadata.
   Embed the minimal native excerpt and its SHA-256 digest; local logs are supplemental evidence.
   Keep author verification separate from independent review.

   A request is pending. A failed transport or skipped review remains failed or skipped.
   Completed results with zero findings are positive terminal results.

   Commit the artifact in the same session that writes it, before any repair edit lands.
   A review worktree moves to the record commit, and the frozen pair still names the source tip.
   Append superseding occurrences and retain earlier findings.
   Apply the bounded repair policy's current-evidence rule and its narrow evidence-only exception after a repair.

   Record validation checks occurrence and source coverage. It cannot prove judgment correctness or authenticate an invented transcript.

   Keep the artifact through completion. The existing retirement path removes it after the committed results remain reachable in history.
   A no-spec review stays chat-only unless the reviewer supplies an explicit slug.

7. **Hand off, don't repair.** This phase makes no fixes and runs no gate. Accepted findings return to `/bench-implement-spec` on the same integration source. A spec amendment commits to that source on the finding cadence.

   A clean chunk review hands its frozen pair back to the orchestrator. The orchestrator starts the successor ticket's fresh author or performs final reconciliation. Only the reconciled final source proceeds to `bench worktree land`; `/bench-final-check` reports that landing's oracle.

   The landing base is the `main` tip merged before the first chunk. `bench worktree land --base` takes that `main` tip, not a later chunk base.

## Ordinary assessment evidence

At each chunk review, update the ordinary-work assessment record through `bench assessment record --input <file>`. Follow the collection guidance in `.bench/BENCH-reference.md`. Retain unavailable harness measurements as unknown; they do not block implementation.
