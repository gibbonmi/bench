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
and the worst issue in each axis. The recommended next command is
`/bench-implement-spec` when findings need fixes, or `/bench-final-check` when the
review is clean or the reviewer accepts the residual risk.

The gate is deterministic: tests, types, lint, conformance. It catches regressions
and rule violations. It cannot tell whether you built the *right* thing the *right*
way. `/bench-review-implementation` is the semantic pass that can — and it's advisory: it surfaces
findings for you, it has no authority to call anything done. The gate and you do.

Run it on the branch diff against its true base, on three axes that stay separate.

## Process

1. **Pin the diff.** Resolve the base and the changed-file list with `bench diff`:
   it prefers the branch's recorded pre-shift base and falls back to merge-base
   with the default branch, and its `method:` line says which happened — a shift
   stacked on unmerged work reviews its own commits, not the feature's. Then read
   the content with `git diff <base>...HEAD` (three-dot) and
   `git log <base>..HEAD --oneline` using the base it printed. Confirm the diff
   is non-empty before going further.

2. **Find the sources.** Spec: `specs/<feature>.md` for this work (or the path I
   give you). Standards: `AGENTS.md` and `.bench/BENCH.md` — the working agreement
   and shared platform rules; `CLAUDE.md` is import pointers only — plus
   `projects/<name>.md` and any `CONTRIBUTING`/conventions docs in the repo.

3. **Spawn the axes in parallel sub-agents** (so they don't pollute each other's
   context), one per axis — Standards, Spec, and the Coverage axis — each under
   ~400 words. Each delegate takes the diff, the sources for its axis, and its
   charge from the `craft-review` skill
   (`.agents/skills/bench-craft-review/SKILL.md`) — the one source for what each
   axis hunts and what a finding must cite; don't restate the charges here.
   Procedural inputs per delegate: give Standards the docs from step 2; give
   Spec the spec file plus, when it carries an acceptance coverage map, the rows
   from `bench coverage <spec>` — auditing them (missing, partial,
   falsely-classified, or unclosed mapped behavior) is part of its charge; give
   Coverage the existing tests and the profile's hostile-input checklist when
   one exists.

   If there's no spec, skip the Spec axis and say so; the Coverage axis still
   runs — it needs only the diff and the existing tests.

4. **Aggregate, don't merge.** Report under `## Standards`, `## Spec`, and
   `## Coverage` headings, findings kept separate. Do not rerank across axes or
   pick a single winner — the separation is the point: code can pass one axis and
   fail another (right thing, wrong conventions; clean conventions, wrong thing;
   correct on the happy path, open on the edges), and merging them lets one mask
   the other. End with a per-axis count and the worst issue within each axis.

## Where it sits

`/bench-review-implementation` is generation-shaping, not enforcement: run it, read the findings, decide
what to fix. Then the deterministic gate (`/bench-final-check`) still has to be green, and the
merge is still yours. Review tells you whether it's *good*; the gate tells you
whether it's *done*; you decide whether it ships.
