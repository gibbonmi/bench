# Prove the new diagnostics through the oracle

Blocked by: grade-citation-execution.md, tighten-citation-grammar.md, refuse-mixed-tag-row-ids.md
Writes: tests/canary/coverage-map-validation/, internal/coverage/citations_test.go
Covers: CE16, CE17, CE18, CE19

## What to build

The new coverage diagnostics must reach a real gate run, and the opt-out must stay
whole.

Two fixtures join the coverage-map-validation canary family. The first fixture
declares a map whose row IDs carry two tags. The second fixture cites a file whose
build constraints no executed tag set satisfies. Each fixture carries an `EXPECT`
file with the exact diagnostic. Each fixture reds a real gate run.

Two assertions close the feature. The first asserts that a historical spec with a
bad citation stays silent, so the historical marker keeps silencing every coverage
check. The second asserts that the shared parse entry point returns the new
violation classes, so the review preflight and the gate agree.

## Acceptance

- [ ] CE16 — the mixed-tag canary fixture reds a real gate run with its `EXPECT`
      diagnostic.
- [ ] CE17 — the unexecuted-tag canary fixture reds a real gate run with its
      `EXPECT` diagnostic.
- [ ] CE18 — a historical spec with a bad citation returns no violation.
- [ ] CE19 — the shared parse entry point returns the mixed-tag, mention, subtest,
      and unexecuted-constraint violation classes.
- [ ] `bench gate` stays green.

## Delegate charge

You work in the Bench repo on the `citation-execution-proof` spec. Read
`specs/citation-execution-proof/spec.md` first. Then read the three delivered
fixtures under `tests/canary/coverage-map-validation/`. Copy their shape exactly.
Each fixture holds a `files/` tree with the fixture spec, plus an `EXPECT` file
with the diagnostic text.

Add one fixture whose map declares row IDs with two tags. Add a second fixture
whose seam cell cites a test file the gate does not execute. Ship that cited file
inside the second fixture's `files/` tree. Ship a minimal `go.mod` in the same
tree, so the fixture root carries a test phase and the census is not empty. Take each `EXPECT` string from the
diagnostic the sibling tickets already produce. Do not invent new wording.

Add two assertions in `internal/coverage`. The first drives the check over a
historical spec that holds a bad citation, and asserts an empty result. The
second drives `ParseSpec`. It asserts the mixed-tag, mention, subtest, and
unexecuted-constraint violation classes come back.

`internal/preflight/gather.go` needs no change, because every `ParseSpec`
violation already blocks the preflight bootstrap. Do not edit it.

Run `bench worktree exec <label> -- go test ./internal/coverage/ ./internal/canary/`.
Do not commit. Do not edit the spec.
