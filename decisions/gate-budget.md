# Gate budget (FT171)

Status: ready

## Destination

One machine-wide resource budget the whole gate draws from, at every nesting
depth, so the run cannot oversubscribe the box. A second effect follows: less
duplicated work inside the sweep it schedules. Two levers under one
destination: bound the demand, and reduce it.

The outer gate saturates the machine on its own: load average ~123 on 16 cores
at the 2026-07-22 baseline with nothing else running. There is no quiet box
during a run, so every symptom here is self-contention and no external load is
needed to reproduce it. Reserved interactive headroom is a first-class
requirement, not a tuning preference. Every width is a formula over the
resolved budget, never a literal, because core counts differ per machine.

Supersedes `gate-concurrency`'s decision #1 (budget model) and #3 (budget =
`runtime.GOMAXPROCS(0)`, the whole box). That map's landed canary arm stands as
built; what changes is that its arithmetic stops computing from the box.

Measured state has moved twice since the #20 census. The #23 concurrency route
cut the focused `internal/gate` median from 150.85 s to 56.72 s. Then the
single-build serial gate and the 2026-08-09 branch-native rebuild (`3701c4a0`)
replaced the fixture-driven workload wholesale. It runs one host binary per
top-level run, one phase process at a time, direct mutation-to-check
canaries, and no stripped-subject reruns. The target of a full gate under 2 minutes is now
exceeded threefold: #26's census on `a3b599ea` measured 38 s wall, ~51 s
CPU, and 25 peak descendants, with no Bench-owned fan-out left. #27 closed
the map 2026-08-13: the destination is met by other means, no spec follows,
and `ready` here marks terminal closure rather than spec-readiness.

## Notes

## Decisions so far

- [What is bounded — processes, or cores?](gate-budget/tickets/1.md): Bound total live demand, denominated in cores rather than processes.
- [How does a token cross the process boundary?](gate-budget/tickets/2.md): An inherited file descriptor holding tokens, jobserver-style.
- [What shape does the reserved-headroom formula take?](gate-budget/tickets/3.md): Proportional and floored: `budget = max(1, resolved * (1 - r))`, with `r` a single named constant priced by #8.
- [What decides whether a process needs a pool?](gate-budget/tickets/4.md): Fan-out owns a pool.
- [How are tokens recovered when a process dies holding them?](gate-budget/tickets/5.md): The fan-out owner reclaims on child exit.
- [What rule picks `r` from the prototype's numbers, and what counts as a symptom?](gate-budget/tickets/6.md): Mechanical: the smallest `r` with zero contention symptoms.
- [What is measured, on which trees and under what conditions?](gate-budget/tickets/7.md): Both trees: today's tree.
- [What are `r` and the span-inflation threshold?](gate-budget/tickets/8.md): Retired unpriced by #27's ruling: no saturating class remains to certify `r` or a grant split against.
- [Can behavior-owned baselines narrow soundly?](gate-budget/tickets/9.md): It carries.
- [Does the Go toolchain participate in an inherited token pool?](gate-budget/tickets/10.md): No.
- [A spawn is granted weight `w`. How is `w` split between `-p` and `GOMAXPROCS`?](gate-budget/tickets/11.md): Priced by #8, against exactly two candidates:
- [Do the cheap phases draw from the pool?](gate-budget/tickets/12.md): Every gate-launched phase acquires; the trivial ones take a single token.
- [What actually makes the expensive packages expensive?](gate-budget/tickets/13.md): The expensive rows do not describe one dev-gate workload.
- [Which workload class is allowed to certify `r`?](gate-budget/tickets/14.md): Yes.
- [Where does one gate run duplicate expensive work?](gate-budget/tickets/15.md): The public dev workflow does not launch the same whole gate twice by accident.
- [Must demand reduction land before the pool is priced?](gate-budget/tickets/16.md): Yes.
- [Does one gate evaluation own one snapshot generation?](gate-budget/tickets/17.md): Resolved and landed.
- [Does the exhaustive gate matrix run at the decision seam?](gate-budget/tickets/18.md): Resolved and landed.
- [Does each conformance fixture run only its registered check?](gate-budget/tickets/19.md): Resolved and landed.
- [What workload remains after demand reduction?](gate-budget/tickets/20.md): Resolved on exact commit `eb6845f` (tree `cd2ece9`), 2026-08-07 local time.
- [Which remaining demand must land before outer-width pricing?](gate-budget/tickets/21.md): Resolved on #20's exact commit `eb6845f`, 2026-08-07.
- [Does the dev gate adopt intra-package test concurrency, and through which mechanism?](gate-budget/tickets/22.md): Route one, plus the specbuild lever.
- [Land the gate kit-root seam and `t.Parallel` adoption](gate-budget/tickets/23.md): Resolved and landed.
- [Land specbuild test-only `t.Parallel`](gate-budget/tickets/24.md): Retired as moot.
- [Land the three sized serial cuts](gate-budget/tickets/25.md): Moot as sized, flagged for reviewer veto against the roadmap's "land #25's cuts".
- [What workload remains after the concurrency route lands?](gate-budget/tickets/26.md): Resolved 2026-08-13 on exact commit `a3b599ea`, same 12-online-CPU host.
- [Does the pool destination survive the serial baseline?](gate-budget/tickets/27.md): Retire and close, reviewer-accepted 2026-08-13.

## Not yet specified

## Spec-writer discretion

- The token encoding in the pipe, and how the acquisition deadline in #5 is
  expressed, provided the fail-toward-less-concurrency posture holds.
- Where the pool owner's bookkeeping lives, provided reclaim-on-child-exit is
  observable from one place.
- #8's sweep order. Comparing #11's two splits at one fixed `r` before
  sweeping `r` keeps the matrix from multiplying. This holds provided the
  chosen split is re-checked at the selected `r`.

## Out of scope

- Scoped or diff-based gating of any kind. FT91 ruled it unsound here, since
  contract and canary are behavior contracts with no file-to-test map, and
  that ruling stands. Scope is not a speed lever.
- Weakening any check to buy wall-clock. Green keeps meaning the same thing.
- Merging same-check fixture runs: `gate-pipeline` #5 rejected it because each
  fixture is a distinct mutated tree and must be graded alone. #9 asks about
  baselines only.
- A cache quota, automatic eviction, or the Go build-cache footprint itself —
  owned by `go-build-cache-footprint`.
- Ship-tier duplicate-proof cleanup, which needs separate shaping because it
  changes publication evidence rather than the dev gate's implementation.
- Publication timeout behavior, which remains owned by FT87.

## Sources

- Path: `decisions/gate-critical-path.md`
  Supports: the 89.91 s entry state and 85.415 s vs 50.917 s `posture` span behind #6's choice of span inflation.
  Drift: historical trigger evidence on the artifact-split landed tree; #20 is the current workload account.
- Path: `decisions/gate-concurrency.md`
  Supports: the landed canary arm, its ~123 load baseline, and decisions #1 and #3 this map supersedes.
  Drift: describes shipped code; re-read if the canary arithmetic changes.
- Path: `decisions/gate-pipeline.md`
  Supports: #5's scoped-baseline soundness argument that #9 asks to extend.
  Drift: none expected while that decision stands.
- Path: `ROADMAP.md`
  Supports: FT171's contention evidence. This includes the 12-core sample where both `TestSetupConflictContracts` FIFO cases exhausted 15 s deadlines under overlap and then passed 3/3 focused at ~0.43 s.
  Drift: a working prioritization document; the row moves as the work lands.
- Path: `decisions/gate-budget/assets/gate-budget-cpu-wall-census.md`
  Supports: #13's three-shape finding, cache A/B, preflight ownership and per-test timing, #20's post-reduction 71-package census, focused repetitions, exact-subject gate span, process fan-out, and cache ruling. It also supports #21's serial-chain finding, per-test attribution, deterministic Git-spawn histograms, and concurrency constraints.
  Drift: #13–#21's figures are pre-rebuild and historical. The decision #26 section is current as of `a3b599ea`, and #27 closed the map with no width pricing to follow.
- Path: `decisions/gate-budget/assets/gate-budget-memory-profile.md`
  Supports: #1–#5's machine-wide process-boundary budget and the non-recursive primary/stripped/canary overlap. It also supports the current operator-width plumbing gap, and #20's required memory/process/I/O observables.
  Drift: one green run on `6607236` and one 12-core host before #18–#19; mechanism evidence only, never authority to price #8.
- Path: `decisions/gate-budget/assets/ft171-shared-fixture-staged-binary.md`
  Supports: #23's measured residual, the shared immutable fixture binary. It also supports the narrowed setup-only work behind the 111.67 s / 78.79 s width figures and the flat ~700k output-block volume.
  Drift: measured on the pre-rebuild fixture-driven suite; landed-route evidence, never the current workload.
- Path: `internal/conformance/registry/packages.go`
  Supports: #13's finding that `internal/preflight` is excluded from the dev package set and owned by the ship tier.
  Drift: code-derived; re-verify if tier ownership changes.
- Path: `internal/gate/evaluation.go`
  Supports: #17's current evaluation-owned accepted pre generation, distinct post generation, and common working/prospective source contract.
  Drift: code-derived; re-verify with the operation ceilings in `internal/gate/evaluation_test.go` if evaluation or snapshot ownership moves.
- Path: `internal/preprelease/preprelease.go`
  Supports: #15's ship step order, repeated core package enumeration, first artifact build, conformance probe, preflight, and ship canary.
  Drift: code-derived; re-verify if the ship sequence changes.
- Path: `internal/releaseevidence/registry.json`
  Supports: #15's release-preflight gate, full-race, vet, artifact, and smoke commands.
  Drift: code-derived; re-verify if release evidence phase ownership changes.
- Path: `internal/gate/runner.go`
  Supports: the current serialized phase table behind #26's baseline and the context-injectable `processGroupGrace` behind #25's mootness.
  Drift: code-derived. The pre-rebuild `canary.go`/`check_slots_test.go`/artifact-fixture sources behind #1, #9, #13, and #18's findings were deleted by `3701c4a0`. Those answers now rest on the census and memory-profile assets, plus spec-retirement history.
