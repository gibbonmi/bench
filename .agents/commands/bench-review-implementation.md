---
description: Three-axis semantic review of a shift branch — Standards (does the diff follow this repo's conventions), Spec (does it match what the spec asked for), and Coverage (what breaking inputs does nothing exercise). Use after /bench-implement-spec, before merge, to catch what the gate can't see. Advisory, not authoritative.
---

# /bench-review-implementation — the check the gate can't run

## Entry orientation

This is the semantic review phase. It reviews the branch diff against three
separate axes: documented standards, the approved spec, and coverage gaps. It
produces findings the gate cannot see, without claiming authority over done-ness.

## Exit handoff

Close by reporting Standards, Spec, and Coverage findings separately, with counts
and the worst issue in each axis. When findings need a later fix pass, they are
also persisted at `reviews/<spec-slug>.md` (step 5). The recommended next command
is `/bench-implement-spec` when findings need fixes, or `/bench-final-check` when
the review is clean or the reviewer accepts the residual risk.

The gate is deterministic: tests, types, lint, conformance. It catches regressions
and rule violations. It cannot tell whether you built the *right* thing the *right*
way. `/bench-review-implementation` is the semantic pass that can — and it's advisory: it surfaces
findings for you, it has no authority to call anything done. The gate and you do.

Run it on the branch diff against its true base, on three axes that stay separate.

## Process

1. **Pin the diff.** Pull the whole base-relative review context with
   `bench diff --full`: it prefers the branch's recorded pre-shift base and falls
   back to merge-base with the default branch, and its `method:` line says which
   happened — a shift stacked on unmerged work reviews its own commits, not the
   feature's. The changed-file list, the `log[N]{sha,subject}` commit table, and
   the raw diff body arrive in one output. Confirm the diff is non-empty before
   going further. When it is empty because the work already landed on the
   default branch (the documented happy path commits before review), review the
   landing commit instead: `bench diff --full --commit <sha>` bounds the same
   bundle to exactly what that commit landed.

2. **Find the sources.** Spec: `specs/<feature>.md` for this work (or the path I
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

5. **Persist actionable findings — the pickup file.** When the review surfaces
   actionable findings that need a later fix pass, write `reviews/<spec-slug>.md`
   before closing, where `<spec-slug>` is the basename of the reviewed
   `specs/<spec-slug>.md` (write the file directly rather than through shell
   snippets, so paths with spaces survive). Mirror the handoff's shape: one
   section per axis — `## Standards`, `## Spec`, `## Coverage` — each with its
   finding count, its worst issue, and every actionable finding with the file or
   doc citation its axis supplied. Keep all three axis headings even when only
   one axis has findings (`0 findings` is a fine section body), so the reader
   sees the full disposition. If a stale artifact already exists for the slug,
   replace it with the current findings — never append a second log. Commit the
   artifact in the same session that writes it — the pickup is tracked state at
   birth, never untracked drift that flips the gate verdict stale.

   A clean review, or one where the reviewer accepts every residual risk,
   writes no artifact: the `reviews/` directory means "there is fix work to
   do", not "a review happened". Never commit an empty `reviews/` directory or
   a `.gitkeep`. A no-spec review stays chat-only unless the reviewer supplies
   an explicit slug — without a durable spec, inventing an artifact name would
   create a second source of feature identity.

   The artifact is transient pickup state, not a review log: the
   `/bench-implement-spec` session that resolves the findings deletes
   `reviews/<spec-slug>.md` in the same green fix commit that closes them, so
   resolved findings can never resurface as false pickup work.

## Where it sits

`/bench-review-implementation` is generation-shaping, not enforcement: run it, read the findings, decide
what to fix. Then the deterministic gate (`/bench-final-check`) still has to be green, and the
merge is still yours. Review tells you whether it's *good*; the gate tells you
whether it's *done*; you decide whether it ships.
