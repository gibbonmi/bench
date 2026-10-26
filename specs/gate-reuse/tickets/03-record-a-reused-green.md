# Record a reused green in the gate cache, and land on another run's green

Blocked by: 02-key-the-evidence-store.md
Writes: internal/gate/run_transaction.go, internal/gate/engine.go, internal/gate/authorization_key_test.go (new), internal/gate/durable_replacement_test.go, internal/worktree/land_reuse_test.go (new), internal/worktree/land_fixtures_test.go, internal/worktree/parallel_census_test.go
Covers: GR22, GR23, GR24, GR44, GR47, GR25, GR26, GR27, GR28

## What to build

The reuse predicate is the `reusableEvidence` decision over the key of ticket 02. The run evaluates it at its two existing sites in `executeSubjectOwning`:

- Before the lock: when the predicate holds and the checkout's gate cache already describes this subject's fresh green, the run answers with no write, as today.
- At the lock: the run takes the gate lock with no wait, as today. When another run holds the lock and the predicate held before the lock, the run answers the reuse with no write and no demotion.
- Under the lock: when the predicate holds, the run writes the gate cache unless it already describes this subject's fresh green, and then it answers.

Today's contention path demotes a reusable green and refuses, so the second site changes that path for a reuse alone.
The written record holds the tree, the checkout's oracle digest, `green`, and the recorded time of the store record. The write uses the lock that the run already holds, because a second acquisition fails against `executionLockOwners`.
A failed write returns an operational result with no reusable green.

The landing needs no production change. `authorization.AuthorizeWithWriters` reads the store through `InspectTreeContext` and attributes a reused green as `green`, and the landing advances the green marker on that kind.
The new landing tests go in `internal/worktree/land_reuse_test.go` and use the landing fixture of `land_fixtures_test.go`, whose prospective gate appends one tally run to the common Git directory:

- GR25 and GR26: move `main`, fold it into the source with `bench worktree merge --from main <target>`, and land the merge tip spec-less. The fixture declares no lane, so the merge runs the whole gate once. Compose the merge call with `mergeSetAt` or the verb runner.
- GR27: run `bench gate --checkpoint specs/x/spec.md --complete` in the source worktree through `gate.RunCommand`, and then land with `--spec x`.
- GR28: the same flow over a base that is a minimal Go module with `cmd/bench` and a build-input manifest (`commitLandingBuildInputs`), with no `scripts/go-build.sh`. The source changes a file of `cmd/bench`.

Without a Go module, both checkouts digest the runner to `unreadable`, and GR28 cannot go red. With no `scripts/go-build.sh`, the prospective gate keeps the inherited run binary.

Each new worktree test calls `t.Parallel()`. Raise `worktreeTestCount` in `parallel_census_test.go` by the count of new tests, with no new line, because that file holds 961 lines over its limit.
`run_transaction.go` holds 370 of 400 lines. Put a new helper in `engine.go` or in a new file of `internal/gate`.

## Acceptance

- [ ] A primary run with no gate cache that reuses a linked-worktree green writes the cache with the tree, `green`, and the recorded time of the worktree run (GR22).
- [ ] A primary run whose gate cache holds a red verdict of the same tree replaces it with the reused green and the worktree run's recorded time (GR47).
- [ ] `TestGateRunGreenReusesRecordedVerdict` stays green: a reuse that the cache already describes leaves the cache bytes unchanged (GR23).
- [ ] A reuse whose gate-cache write fails returns an operational result with no reusable green (GR24).
- [ ] A reuse while another run holds the checkout's gate lock, with a cache of another tree, prints the reuse line and leaves the cache bytes unchanged (GR44).
- [ ] After a whole-gate merge of `main` in a fixture with no lane, a spec-less landing of the merge tip leaves the tally at one run (GR25) and advances the green marker to its published commit (GR26).
- [ ] After a green complete checkpoint in the source worktree, the spec landing leaves the tally at one run (GR27).
- [ ] When the source changes a file that the declared build inputs name, the spec landing after the complete checkpoint leaves the tally at two runs (GR28).
