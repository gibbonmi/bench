# Own the focused test binary

Blocked by: introduce-run-scoped-bench-selection.md
Ownership fence: `internal/testreport/`
Integration surfaces: selection producer→introduce-run-scoped-bench-selection.md; selected environment for unchanged-host suites→migrate-contract-preflight-helpers.md; focused-owner teardown→contract-run-directory-lifecycle.md
Contracts: `Selection` crosses `internal/testreport/`→the one `go test -count=1` child, membership is one owner-authored value when the command starts without inheritance and the exact inherited value otherwise, ordering is validate usage then author-or-reuse then launch and reap the child, absence in top-level mode authors once while invalid inheritance refuses, asserted by TR1 against the real command boundary
Closure: TR1/top-level-single-build, TR1/inherited-zero-build, TR1/child-path-propagation, TR1/next-run-new-selection, TR1/interrupt-reap

## What to build

Make `bench test` the focused-test run owner. Its existing one-child process-group lifecycle stays intact; the new selection is passed to that child and cleaned after the child is reaped. An already selected parent is reuse inside the same run, never a cross-run cache hit.

## Acceptance

- [ ] [TR1] (covers RS3) `bench test` builds once without inherited selection, builds zero and preserves the exact path with valid inheritance, passes that path to its child, creates a new path for a later top-level snapshot, and reaps on interrupt.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TR1/top-level-single-build | omit ownership or invoke the builder twice | testreport builder-count test | run one public focused command and require one trace |
| TR1/inherited-zero-build | treat valid inheritance as permission to author again | testreport inherited-selection test | pass a valid selected path, run the command, and require zero traces |
| TR1/child-path-propagation | drop or rewrite `BENCH_RUN_BINARY` in child environment | child marker test | run a package that prints its selection and require exact equality with the owner path |
| TR1/next-run-new-selection | retain a stable worktree or cache path | consecutive-run test | run red then green source snapshots and require different private paths with matching source seals |
| TR1/interrupt-reap | remove the run directory before the test process group exits | bounded child teardown test | interrupt a child with a lingering descendant and require reap before cleanup marker |
