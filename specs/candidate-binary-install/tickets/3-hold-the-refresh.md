# Hold the post-landing refresh and land through the installed candidate

Blocked by: 1-declare-persisted-fields.md, 2-install-the-candidate.md
Writes: internal/worktree/land_effects.go, internal/worktree/land_effects_test.go, internal/worktree/land_broker_notice_test.go, internal/worktree/land.go, internal/worktree/land_resume.go, internal/worktree/joins.go, internal/worktree/parallel_census_test.go, internal/commitment/repository/window.go, internal/commitment/repository/window_test.go, .bench/BENCH-reference.md, internal/systemtest/candidate_install_test.go, internal/systemtest/owner_test.go
Covers: CB24, CB25, CB50, CB51, CB52, CB53, CB59

## What to build

Deliver the hold half of CB-C2 and its end-to-end proof: while a declared build is open and the manifest binds the primary executable, the post-landing refresh writes nothing and reports `held`, and the installed candidate lands its own new fields.

Export two answers from `internal/commitment/repository/window.go` over the unexported open-build derivation. The first answer reports whether any open build exists, of any outcome. The second answer reports whether a declared build is open, and it names the build's integration target: the first bound assignment of the declared build. Store that assignment in the open-build value when the derivation selects it. A declared build is an open build whose default-branch bound spec declares a persisted field; use the ticket 1 scan for the declaration. Derive the second answer from the first, so the open-build derivation keeps one source. Tickets 5 and 6 consume these answers.

In `refreshBroker`, read the held-install predicate before the seal verification. The predicate holds when a declared build is open and the primary manifest digest equals the primary executable's digest. It also holds when the ledger does not read. When it holds, the refresh writes nothing, and the effect result is a new `held` constant beside `complete`, `failed`, `pending`, and `not-applicable`. A `held` result exits 0. When the predicate does not hold, the refresh runs as today.

Reach the predicate through a `joins.go` seam so `landedAfterEffects` keeps its signature. `land.go` has 3 lines of headroom, so the author keeps any change there inside that headroom. `brokerChangeNotice` keeps its current behavior for a `held` refresh.

Add `held` to the effects result list in `.bench/BENCH-reference.md`, which is the one guidance copy of the result set, with one sentence that says when the refresh holds (CB59, review-owned).

Extend ticket 2's system test file with two landing proofs. `TestCandidateInstallLandsANewLedgerField` writes the new ledger key through the candidate, lands the build through the unchanged land route, and reads the manifest and the primary executable after the landing. `TestCandidateInstallSurvivesAGapLanding` opens a chunk gap, lands an unrelated change through the broker, and reads the primary digest and `bench status` after that landing. Run the suite only through `bench test --check system --run TestCandidateInstall`, which sets `BENCH_KIT` to the tree under test.

## Acceptance

- [ ] `TestRefreshHoldsAnInstalledCandidate`: with a declared build open and equal digests, the refresh reports `held` and writes no file (CB51).
- [ ] With a declared build open and different digests, the refresh rebuilds and reports `complete` (CB52).
- [ ] With no declared build open and an unverified seal, the refresh rebuilds and reports `complete` (CB53).
- [ ] With an unreadable ledger, the refresh reports `held`.
- [ ] The installed candidate lands a build whose published tree holds a review record with a `folds` key (CB24).
- [ ] After that landing, the primary manifest digest equals the primary `dist/bench` digest (CB25).
- [ ] After an unrelated landing during an open gap, the primary digest is unchanged and `bench status` exits 0 (CB50).
- [ ] Every existing effects-row test keeps its current result.
- [ ] `.bench/BENCH-reference.md` lists `held` and states the hold condition (CB59).

## Checkpoint verification

Run the `worktree`, `system`, `internal-commitment-repository`, and `root-conformance` commands of CB-C2's verification plan, and run `bench test --check prose` for the reference edit. Probe the predicate with a mutant that ignores the open build and with one that ignores the digest, and record each red. The orchestrator then freezes the CB-C2 delta for independent review.
