# Retro: worktree-test-floor

## Outcome

The spec landed on `main` as `cc408adf` (tree `b82b44e7`) from the reviewed pair
base `09287798`, source tip `43340746`. The `internal/worktree` suite runs 312 of
its 355 tests in parallel under a census. Four verbs take an explicit root,
and every verb takes an explicit home at its command boundary. Sixteen
tickets, three spec amendments, one review pickup commit, and two repair
tickets landed on the integration source. The review found 6 raw findings
across three axes and 3 repair targets, all closed by ticket 17. The
repair-scoped re-review found one repair-induced gap, which ticket 17b
closed, and the second scoped check was clean.

The package ran 70.0 s on `main` with 335 serial tests. It runs 22.5 s at the
tip with 312 parallel tests and 43 serial tests that sum to 7.6 s. In the
gate, the package wall went from 63–73 s to 21.9 s. It no longer sets the
`test` phase; `freshness`, `runbinary`, and `conformance` do at 12–18 s. A
quiet gate went from about 94 s to about 75 s before the last three tickets.
At the tip, `-parallel 4` equals the default on this package, and
`-parallel 2` costs about 20 s more.

## Gate-stage timings

The landing gate took 44.4 s. Its stages were gofmt 0.1 s, vet 1.3 s, test
24.7 s, race 4.5 s, system 9.7 s, and shellcheck 0.4 s. The landing before
this spec took 102.3 s with a 79.3 s test stage. Every worktree
commit paid a full gate of about 75 s to 105 s, because the installed broker
predates the fast lane. Twenty-two full runs paid for ticket commits, three for spec amendments, one
for the review artifact, and three were red. One red was a prose bound, one a
hand-typed path list, and one a worktree created during the gate.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized. Tickets 05–10, 14, and 15 ran at sonnet / low
and landed first-pass. Tickets 02, 03, 11, 12b, 13b, 16, and 17 ran at
opus / medium; ticket 03 took one repair round for a missing census row.
Ticket 01 ran at opus / low and landed first-pass with a design deviation
the byte pins forced. Tickets 12 and 13c ran at opus / high; ticket 12
stopped once on a seam gap the reviewer re-scoped, and ticket 13c landed
first-pass. Ticket 13 stopped on a stub-swap collision outside its fence,
which became tickets 13b and 13c.

Tickets 01–03 ran in parallel off `main`; tickets 04–10 ran seven wide off
the integration tip. The coordinator folded each diff by patch in
`Blocked by:` order and probed each fold at a site and kind distinct from the
delegate's. Ticket 04 landed no diff: its five tests bind `BENCH_HOME`
in-process, which the reviewer turned into the home seam.

## Coordinator catches

- The measured wall after eleven tickets was 73 s, not the spec's 20 s. The
  93 parallel tests summed to 3.5 s and the 256 serial tests to 72.6 s. The
  spec's premise "the cost is count" was wrong, and the reviewer reopened the
  `BENCH_HOME` seam.
- Ticket 03's census passed a parent whose subtest called `t.Parallel()`;
  the coordinator's closure probe stayed green, and the delegate added the
  row.
- Ticket 09 found the census blind to `journeyStubGit`'s `PATH` bind; ticket
  11 routed it through `bindEnv`.
- Ticket 12 found the gate child's environment is the closure's business;
  row WF16 moved to the exec and subshell children.
- Ticket 13 found 47 tests that swap package-level variables; the census
  gained the assignment edge, then the joins became a per-call value.
- The coordinator caused two reds: a hand-typed path list that omitted three
  files, and a worktree create during a gate that flipped a live-state test.

## Repair attribution

| ticket | repair rounds | cause per round |
|---|---|---|
| 01-take-explicit-root-in-four-verbs | 0 | none |
| 02-guard-harness-state-for-parallel-siblings | 0 | none |
| 03-build-parallel-census-and-pins | 1 | delegate-error |
| 04-mark-five-slowest-tests-parallel | 0 | none |
| 05–10 marks | 0 | none |
| 11-mark-harness-tests-and-close-census | 0 | none |
| 12-take-explicit-home-in-every-verb | 1 | spec-row |
| 12b-serialize-package-global-stub-swaps | 0 | none |
| 13-fixtures-pass-a-private-home | 1 | ticket-slicing |
| 13b-a-stubbing-test-is-serial | 0 | none |
| 13c-the-landing-joins-become-a-per-call-value | 0 | none |
| 14 and 15 direct binds | 0 | none |
| 16-pin-the-serial-ceiling-and-close | 0 | none |
| 17-close-the-review-findings | 1 | other |
| 17b-census-root-finder | 0 | none |

## Agent-experience improvements

### Bench CLI

- Let `bench commit` take its path list from the checkout's status, because a hand-typed list omitted three files and paid a red gate.
  Feeds: new
- Give the gate a `--test-parallel` option and a shared test-thread pool, because delegate test runs beside a gate doubled every package wall.
  Feeds: new
- Let `bench worktree clean --apply` take a fingerprint that survives a sibling removal, because each removal invalidated every other plan.
  Feeds: new

### Skills

- Make the `craft-delegate` charge forbid a background test run, because a delegate that backgrounds a run ends its turn before the result.
  Feeds: none
- Make the `craft-spec` edge inventory ask which tests swap a package-level variable, because the census predicate could not see that class.
  Feeds: none

### Process

- Measure the wall after the first slice of a performance spec, because the premise "the cost is count" failed only under measurement.
  Feeds: none
- Never create or clean a worktree while a gate runs, because a live-state test compares two reads seconds apart.
  Feeds: none
