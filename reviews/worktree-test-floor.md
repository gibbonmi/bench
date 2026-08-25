# Review pickup: worktree-test-floor

Base `09287798`, reviewed tip `31f3af13`, on the retained source `wtf-int`.
Three axes ran read-only over the pinned diff. Each section carries its
count, its worst issue, and every actionable finding with its disposition.
Raw findings: 6. Repair targets after de-duplication: 3, all accepted by the
reviewer on 2026-08-25 into ticket 17.

## Standards

Count: 0 hard violations, 1 judgment call. Worst: six verbs carry an unused
`_ string` home so every verb shares one calling convention. The site is
`internal/worktree/land.go` and five siblings, and `LandCommand`'s executable
parameter is the precedent. Disposition: no-op. The axis examined the census
walk and predicates, the `joins` value, `effects.go`, `exec.go`, and the
boundary resolvers in `cmd/bench/main.go`. It also examined comment prose,
the changelog, and the `StubGit` split. It refuted each as a duplication or
comment finding.

## Spec

Count: 1. Worst: row WF06 names the seam "selector unit under `-race`", but
`TestParallelJourneysRecordEverySelection` (`internal/worktree/test_run_test.go`)
is not in `internal/racetests/racetests.go`. So the gate's `race` phase never
runs it. Disposition: auto-fix — add the registry entry. Every other row maps
to a test that proves its clause. The byte pins for story 9 and the `Home()`
export for story 13 are spec-authorized.

## Coverage

Count: 4; two are no-op after refutation. Worst: a test that assigns an
imported package's variable through a selector (`cryptorand.Reader = …` in
`setRandomReader`, `internal/worktree/list_actions_test.go`) is invisible to
the census, which matches only a bare identifier on the left-hand side. The
two callers are serial today only because they also bind `BENCH_HOME`.
Disposition: ask-user, accepted — the census treats `pkg.Ident = …` for an
imported package as a serial edge, with a synthetic test.

`TestWorktreeRoutesKeepTheirBytesFromASubdirectory`
(`cmd/bench/command_registry_test.go`) compares two live `worktree list`
runs of the real repository, so any worktree change between the runs flips
it. Disposition: ask-user, accepted — the test runs against a fixture
repository.

No-op: the one-off `prospective authorization refused: infrastructure` under
parallel load shares no file across tests (the gate cache is per test git
dir); it is resource pressure, which the parked test-thread pool owns. No-op:
a helper that moves into a non-test file cannot call `bindEnv` or `chdir`,
so the serial set cannot drift that way.
