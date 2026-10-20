# Hold the post-landing refresh for an installed candidate

Blocked by: 1-declare-persisted-fields.md, 2-install-the-candidate.md
Writes: internal/worktree/land_effects.go, internal/worktree/land_effects_test.go, internal/worktree/land_broker_notice_test.go, internal/worktree/land.go, internal/worktree/land_resume.go, internal/worktree/joins.go, internal/worktree/parallel_census_test.go, internal/commitment/repository/window.go, internal/commitment/repository/window_test.go, .bench/BENCH-reference.md
Covers: CB51, CB52, CB53, CB59

## What to build

Deliver the hold half of CB-C2: while a declared build is open and the manifest binds the primary executable, the post-landing refresh writes nothing and reports `held`.

Export one rule from `internal/commitment/repository/window.go` over the unexported open-build derivation. The rule answers whether a declared build is open, and it names the build's integration target: the first bound assignment of the declared build. Store that assignment in the open-build value when the derivation selects it. A declared build is an open build whose default-branch bound spec declares a persisted field; use the CB-C1 header scan for the declaration. CB-C3 and CB-C4 consume this rule, so it is the one source of the open declared build and its target.

In `refreshBroker`, read the held-install predicate before the seal verification. The predicate holds when a declared build is open and the primary manifest digest equals the primary executable's digest. It also holds when the ledger does not read. When it holds, the refresh writes nothing, and the effect result is a new `held` constant beside `complete`, `failed`, `pending`, and `not-applicable`. A `held` result exits 0. When the predicate does not hold, the refresh runs as today.

Reach the predicate through a `joins.go` seam so `landedAfterEffects` keeps its signature. `land.go` has 3 lines of headroom, so the author keeps any change there inside that headroom. `brokerChangeNotice` keeps its current behavior for a `held` refresh.

The effects result list in `.bench/BENCH-reference.md` is the one guidance copy of the result set. Add `held` to that list and one sentence that says when the refresh holds (CB59, review-owned).

## Acceptance

- [ ] `TestRefreshHoldsAnInstalledCandidate`: with a declared build open and equal digests, the refresh reports `held` and writes no file (CB51).
- [ ] With a declared build open and different digests, the refresh rebuilds and reports `complete` (CB52).
- [ ] With no declared build open and an unverified seal, the refresh rebuilds and reports `complete` (CB53).
- [ ] With an unreadable ledger, the refresh reports `held`.
- [ ] Every existing effects-row test keeps its current result.
- [ ] `.bench/BENCH-reference.md` lists `held` and states the hold condition (CB59).

## Checkpoint verification

Run the `worktree` command of CB-C2's verification plan. Probe the predicate with a mutant that ignores the open build and with one that ignores the digest, and record each red.
