# Name the root conformance check

Blocked by: none
Writes: internal/gate/phases.go, internal/testreport/check_test.go, internal/testreport/checks_face_test.go, internal/testreport/command.go, internal/testreport/environment.go, internal/testreport/named_check.go, internal/testreport/root_conformance_test.go (new)
Covers: TF1, TF2, TF3, TF4, TF5

## What to build

`bench test --check root-conformance` runs the conformance entry test with the dev tier and an empty scope. An empty scope runs every registered dev-tier check. So this check is the one exact command for the unscoped root conformance pass.

The gate exports the predicate `gate.RootConformanceRuns(root, kit string) bool` in `internal/gate/phases.go`. It holds the rule that `rootConformanceEnv` applies today: the root is the kit, and the kit declares the entry test. An empty kit means the graded root, through the existing gate resolution `kitRootAt`. `rootConformanceEnv` and the named check call this predicate, and no other package keeps a copy of the rule.

The kit that each caller passes is `gate.KitRoot(root)`, the existing export at `internal/gate/tag_census.go:231`. It answers `BENCH_KIT` when that variable is set, and the graded root otherwise. `testBenchSource` in `internal/testreport/command.go` repeats that rule today. This ticket makes `testBenchSource` return `gate.KitRoot(root)`, so `internal/testreport` keeps no private copy of the kit rule.

Where the predicate answers false, the check exits 1 with a refusal. The refusal goes into `runFocusedRequest` in `internal/testreport/command.go`, beside the system check refusal. It comes before the `selectRunBinary` call, so it starts no Go child.

The name `root-conformance` joins the named-check inventory with the kind `conformance`. The help list, `bench test --checks`, and the unknown-check inventory read that one inventory. The named-check environment passes the empty scope for this name only, because the entry test refuses a scope that names no registered check.

`internal/testreport/check_test.go` holds 498 lines and has no structure grant. The lane `structure` row reds a file over 400 lines that grows. So an edit in that file keeps it at 498 lines or fewer, and each new test goes into `root_conformance_test.go`.

At the ticket commit, rewrite the `planned` seam cell of each covered row into the citation form.

## Acceptance

- [ ] TF1: `bench test --check root-conformance` runs `TestRootConformance` with the dev tier and an empty conformance scope.
- [ ] TF2: A planted conformance defect in a kit copy makes the check exit 1 and print the diagnostic of that conformance check.
- [ ] TF3: Over a root where the gate predicate answers false, the check exits 1 with a refusal and starts no Go child.
- [ ] TF4: `bench test --checks` lists `root-conformance` with the kind `conformance`.
- [ ] TF5: The check runs exactly where the gate predicate answers true, over a kit root, a kit without the entry test, and a foreign root.
- [ ] `testBenchSource` returns `gate.KitRoot(root)`, and `internal/testreport/command.go` holds no read of `BENCH_KIT`.
