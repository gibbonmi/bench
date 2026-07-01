---
description: Two-axis semantic review of a shift branch — Standards (does the diff follow this repo's conventions) and Spec (does it match what the spec asked for). Use after /bench-implement-spec, before merge, to catch what the gate can't see. Advisory, not authoritative.
---

# /bench-review-implementation — the check the gate can't run

## Entry orientation

This is the semantic review phase. It reviews the branch diff against two separate
axes: documented standards and the approved spec. It produces findings the gate
cannot see, without claiming authority over done-ness.

## Exit handoff

Close by reporting Standards findings and Spec findings separately, with counts
and the worst issue in each axis. The recommended next command is
`/bench-implement-spec` when findings need fixes, or `/bench-final-check` when the
review is clean or the reviewer accepts the residual risk.

The gate is deterministic: tests, types, lint, conformance. It catches regressions
and rule violations. It cannot tell whether you built the *right* thing the *right*
way. `/bench-review-implementation` is the semantic pass that can — and it's advisory: it surfaces
findings for you, it has no authority to call anything done. The gate and you do.

Run it on the branch diff against its merge-base, on two axes that stay separate.

## Process

1. **Pin the diff.** `git diff <base>...HEAD` (three-dot, against merge-base) and
   `git log <base>..HEAD --oneline`. Default base is the project's default branch.
   Confirm the ref resolves and the diff is non-empty before going further.

2. **Find the sources.** Spec: `specs/<feature>.md` for this work (or the path I
   give you). Standards: `CLAUDE.md`, `projects/<name>.md`, and any
   `CONTRIBUTING`/conventions docs in the repo.

3. **Spawn both axes in parallel sub-agents** (so they don't pollute each other's
   context), each under ~400 words:

   - **Standards** — every place the diff violates a documented convention. Cite
     the rule. Separate hard violations from judgment calls. Skip anything the gate
     already enforces — no point double-reporting what tooling caught.
   - **Spec** — (a) requirements the spec asked for that are missing or partial;
     (b) behavior in the diff that wasn't asked for (scope creep); (c) requirements
     that look implemented but wrong. If the spec has an acceptance coverage map,
     also audit each coverage row: missing, partial, falsely-classified, or
     unclosed mapped behavior is a Spec finding. Quote the spec line for each
     finding.

   If there's no spec, skip the Spec axis and say so.

4. **Aggregate, don't merge.** Report under `## Standards` and `## Spec` headings,
   findings kept separate. Do not rerank across axes or pick a single winner — the
   separation is the point: code can pass one axis and fail the other (right thing,
   wrong conventions, or clean conventions, wrong thing), and merging them lets one
   mask the other. End with a per-axis count and the worst issue within each axis.

## Where it sits

`/bench-review-implementation` is generation-shaping, not enforcement: run it, read the findings, decide
what to fix. Then the deterministic gate (`/bench-final-check`) still has to be green, and the
merge is still yours. Review tells you whether it's *good*; the gate tells you
whether it's *done*; you decide whether it ships.
