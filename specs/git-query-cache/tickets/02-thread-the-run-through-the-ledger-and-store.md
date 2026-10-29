# Thread one run through the ledger, the commitment store, and the roadmap reads

Blocked by: 01-cache-full-object-id-reads.md
Writes: internal/commitment/repository/admission.go, internal/commitment/repository/candidate.go, internal/commitment/repository/closure.go, internal/commitment/repository/gap.go, internal/commitment/repository/lease.go, internal/commitment/repository/light_path.go, internal/commitment/repository/outlook.go, internal/commitment/repository/override.go, internal/commitment/repository/publication.go, internal/commitment/repository/readiness.go, internal/commitment/repository/repository.go, internal/commitment/repository/run_test.go (new), internal/commitment/repository/sources.go, internal/commitment/repository/verification.go, internal/commitment/repository/window.go, internal/conformance/bounds_policy_test.go, internal/intent/assignment.go, internal/intent/branches.go, internal/intent/intent.go, internal/intent/repo_test.go (new), internal/intent/replacement_test.go, internal/intent/transaction.go, internal/roadmap/tree_validation.go, internal/roadmap/tree_validation_test.go (new), tests/canary/package-core-guard/bounds-intent-window-fixed, internal/gittest/gittest.go
Covers: GQ27, GQ28, GQ29, GQ30, GQ31, GQ32, GQ33

## What to build

Add the value `intent.Repo{Run, Root}` and its constructor `intent.In(run, root)`. The value's methods are the ledger reads and writes.
`Repo` reads the ledger address through `r.Run.CommonDir(r.Root)`. The completion-plan probe of this chunk swaps that exact line.
Each exported root-taking `intent` function becomes a one-line delegation with a nil run. So no exported caller changes in this ticket.

Collapse each unexported pair in `internal/intent` to one run form, such as `transact` and `reauthorizeAssignment`. Each test calls that form with a run or nil.
The test caller `internal/intent/replacement_test.go:149` changes with `reauthorizeAssignment`.
Keep the line `bounds.VerdictWindow(bounds.IntentLockTimeout)` in `intent.go`. The pin at `internal/conformance/bounds_policy_test.go:46` and the canary `bounds-intent-window-fixed` read it.

Give `repository.Store` a `Run` field. Its ledger, lease, window, and tree reads go through that run, and a store with no run caches nothing.
The 34 production `Store{Root` literals stay as they are, so they stay uncached. Ticket 03 adds the run to the landing's literals.
`Store.policyAt` lists through `Run.ListTree` in place of the raw `git.Output` call at `repository.go:250`. A malformed listing record then refuses.

Give `RevisionIndex` and `RevisionDocument` run forms that the store calls. The root forms delegate with a nil run.
Their raw listings at `internal/roadmap/tree_validation.go:55` and `:73` move onto `Run.ListTree`, with the argv that `ReadTreeFile` uses.
The `ls-tree --name-only` call at `tree_validation.go:78` stays live, as the spec's Won't-handle line states.

This ticket uses the nil-safe run forms, `git.FullObjectID`, and the run's `ListTree` and `ReadTreeFile` from ticket 01.
Each new test counts git starts through `gittest.StubGitDir` and `gittest.StubGitStarts`. These packages meet no serial pin or count pin.

## Acceptance

- [ ] Two ledger reads through one `intent.Repo` with a run start the common-directory query once (GQ27).
- [ ] Two root-taking `intent.Read` calls start the common-directory query twice (GQ28).
- [ ] Two policy reads through one `repository.Store` with a run, at one full tree id, start `ls-tree` once and `cat-file` once (GQ29).
- [ ] Two ledger reads through one `repository.Store` with a run start the common-directory query once (GQ30).
- [ ] Two policy reads through a `repository.Store` with no run start `ls-tree` twice (GQ31).
- [ ] A policy read at a tree that lacks `.bench/commitment.json` reports absence through the run's listing (GQ32).
- [ ] Two `RevisionIndex` reads through one run at one full tree id start `ls-tree` once (GQ33).
- [ ] Each exported `intent`, `repository`, and `roadmap` root form keeps its signature.

## Checkpoint verification

- `bench test --package ./internal/intent`
- `bench test --package ./internal/commitment/repository`
- `bench test --package ./internal/commitment`
- `bench test --package ./internal/roadmap`
- `bench test --package ./internal/conformance`
- `bench test --package ./internal/worktree`

Run `bench consumers` on each changed unexported function before the edit. Record one red for each independent start count before the green.
