# Distinguish the published-incomplete exit from a refusal

Blocked by: enrich-refusals-through-one-emitter.md
Writes: internal/worktree/land.go, internal/worktree/land_test.go, internal/systemtest/owner_test.go, .bench/BENCH-reference.md, CHANGELOG.md
Line: opus / high — the exit contract changes at an irreversible boundary.

## What to build

A landing that published its commit and then stopped at the marker, reconcile,
or release step exits 3 — on the first run and on the resume — and its output
carries `next=` with the resume invocation: published commit, `--base`,
`--source-tip`, `--spec`, and path filled exactly, and a placeholder for the
request token (a caller token is never echoed). The `next=` value uses the
emitter ticket's field grammar and the orphan-line quoting precedent.
Refusals keep exit 1, usage errors keep exit 2, a released landing keeps exit
0. `landedIncomplete` and `resumeIncomplete` collapse to one renderer, so the
incomplete record shape has one source. The landing paragraph in
`.bench/BENCH-reference.md` states the four exit meanings.

## Acceptance

- [ ] A first-run landing whose release step fails after publication exits 3 with `worktree=incomplete:release` (covers LR6).
- [ ] A marker failure and a reconcile failure after publication each exit 3 (covers LR7).
- [ ] The incomplete output contains the resume invocation with the published commit, both identity flags, the spec, and the path filled exactly and a placeholder for the request token (covers LR8).
- [ ] A resume that stops at a follow-up step exits 3, not 1 (covers LR9).
- [ ] A pre-publication refusal still exits 1 and a grammar error still exits 2 (covers LR10).
- [ ] A complete landing still exits 0 with `worktree=released` (covers LR11).
- [ ] `.bench/BENCH-reference.md` states the 0/1/2/3 landing exit semantics (story 11, review-verified).
