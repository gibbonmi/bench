# Review findings — governed-npm-publication (FT83 slice 3)

Reviewed range `ae27e76..27e3832`. Advisory; the reviewer decides what gets fixed.

## Standards

3 findings (1 hard, 2 judgment). Worst: the canary fixture-harness paste (hard).

1. **Hard — fixture harness pasted 4×.** The ~90-line `canaryRoot` / `canaryArtifacts` /
   `canaryApprovedSet` helper block is byte-identical across all four canary test files
   (`tests/canary/behavior-owned/publication-order-bypass/files/internal/publication/order_bypass_canary_test.go:66-148`
   and the same region in `publication-unpublish-attempt`, `premature-wrapper-promotion`,
   `integrity-mismatch-acceptance`). AGENTS.md: "a fixture harness pasted N times … must
   collapse to one source"; the independent-expectation exception explicitly excludes
   fixture harnesses. Collapse route exists — the four already share
   `publication-canary-base.txt` via BASE include, so a common helper file can ride the
   same mechanism. The `canaryRegistry` bodies legitimately diverge 3-of-4; only the
   three helpers are pure paste.
2. **Judgment — `next_action` ordering policy derived twice.**
   `statusNextActionForStaged` (`internal/publication/command.go:292-316`) reconstructs
   the platform→wrapper→promote sequence from record transitions; `RunStagedPublication`
   (`internal/publication/statemachine.go:300-360`) derives the same policy procedurally,
   returning the same constants. A pure predicate both consume would keep status
   registry-free while collapsing to one source.
3. **Judgment — wrapper identity minted three ways in Go.** Status keys the wrapper off
   the literal `p.Package == "redbench"` (`command.go:301`); the state machine uses
   `Kind == "wrapper"` (`statemachine.go:33`, `plan.go:42`); the name is minted in
   `approved.go:76`. `Provenance` (`record.go:30-33`) carries no `Kind`, forcing the
   literal. No shared constant.

## Spec

1 finding (partial requirement). Worst: first-publication promotion is unreachable.

1. **First-publication path can never move `latest`; `bench release promote` silently
   no-ops.** Spec Solution: first publication "verifies each live integrity … then
   promotes platform `latest` tags first and the wrapper last." `RunFirstPublication`
   sets `record.Result = "success"` after candidate-tag publish
   (`internal/publication/statemachine.go:160`) and never issues `TagAdd("latest")`;
   `RunPromotion` early-returns on `Result == "success"` (`statemachine.go:384`), so a
   first-path `promote` exits 0 with `next_action=release-complete` and no `latest` tag
   ever moves. Only the staged path (which leaves `in_progress`) reaches real promotion
   (`publication_interrupt_test.go:239` comment confirms). The acceptance map carries no
   first-path→latest row, so the gap is Solution-vs-map; reviewer call whether first-path
   promotion is in-slice — if it is, it's broken.

All 13 coverage-map rows otherwise closed. The four build-session flags refuted clean on
this axis: submit-only authority holds via the record/digest chain; registry-free status
meets the CLI contract; drift check covers all four ops; go-test canaries follow the
existing `bench_canary` precedent and satisfy story 8.

## Coverage

4 findings (plus 1 minor). Worst: adapter divergence on empty registry integrity.

1. **Empty registry integrity classified oppositely by the two adapters; untested.**
   A `200 {"integrity":""}` makes `FixtureRegistry.Integrity` return live=true with
   empty integrity (`internal/publication/fixture_registry.go:125-131`, no empty guard)
   → terminal mismatch; `NPMCLIRegistry.Integrity` maps empty → live=false
   (`internal/publication/npm_registry.go:182-184`) → republish attempt. Same registry
   state, contradictory behavior. The edge-inventory "empty registry response
   distinguishes absent from present-empty" row is unclosed — no publication test
   exercises empty or whitespace-only integrity. Missing: a hostile-registry row
   asserting one classification.
2. **Cross-index-digest resume guard has zero coverage.** The "resume is unsafe" guard
   (`statemachine.go:74-76`, `:282-283`, `:381-382`) refuses a record whose
   `release_index_sha256` differs from the current index, but no test patches the
   on-disk record's digest and asserts the refusal — deleting the guard stays green.
3. **Corrupt/truncated `publication-record.json` unexercised.** `LoadRecord` returns
   "publication record is malformed" (`internal/publication/record.go:64-66`); no test
   writes `{garbage` and asserts exit 1 + attributed message on `status`/`submit`
   (profile hostile-input "hand-edited files" class).
4. **Concurrent second invocation is an undecided edge.** No lock anywhere; two
   simultaneous `submit`s both load an empty record and the last atomic `SaveRecord`
   wins, dropping the other's transitions. Idempotency rows cover only sequential
   reruns. Needs either a lock + test or an explicit won't-handle in the spec.
5. *Minor:* `RunRollback` checks neither digest binding nor `VerifyPublishAuthority`
   (the `profile` param is dead in promote/rollback), and a promote-after-success rerun
   never re-verifies drifted `latest` tags.
