# Delegation discipline

Charged from `craft-delegate` when the coordinator writes a charge, runs a probe,
or accepts a return. Each rule below settles one question the delegation must
answer. `craft-delegate` keeps the inline allowance, the worktree rule, the
mutation-probe rule, and the done-claim check.

## Before the charge

- A spec-backed ticket that will commit on the integration source runs there from
  the start, in `Blocked by:` order. An independent worktree is for a diff that
  lands on its own.
- Disjoint ownership fences across sibling tickets do not license concurrent
  writers in one tree. The lever for parallelism is separate worktrees. A build
  that wants serial verdicts and parallel delegates says which one it buys.
- A parallel repair batch declares which files its tickets expect to touch in
  common. The coordinator then orders the ports to keep the conflict surface
  small.

## Repair-charge template

Every repair charge uses this template. An improvised charge can omit a field
that the coordinator needs to verify the repair.

```text
Base commit: <commit that the repair changes>
Ownership fence: <exact repo-relative file or path prefix>
Effort: <level and iteration cap>
Focused suite: <exact command>
Independent biting probe: <property, mutation kind, site, and expected red>
```

## In the charge

- A write charge names the root conformance pass and the file's wrap width in its
  verification list.
- A write charge states that its ownership fence is a refusal boundary above any
  mirror-every-registry instruction. The delegate reports an out-of-fence write
  before the delegate edits.
- A charge whose collapse crosses its `Writes:` list inside the spec fence takes a
  fence extension in a continuation. The delegate never writes a second spelling
  to stay in fence.
- A charge that exports from a package outside its fence names the fence amendment
  in its return.
- A charge treats a coverage-row citation in a test doc as a reference to keep.
- A pure test-addition charge names the production mutation each new test is
  proven against. The delegate then cannot return a test that is green by
  construction.
- A repair charge names the property the fixture must make red-capable, and lets
  the delegate derive the shape. A prescribed shape buys a vacuous test when the
  prescribed shape is not red-capable.
- A charge that pins an exact verification command names one the coordinator ran
  this session. Otherwise the charge makes the delegate report the command's wall
  time and prove that the command did not skip. A test binary that exits green in
  milliseconds is a skip, not a pass.
- A charge that adds a `_test.go` file under `internal/gate/`, `internal/git/`,
  `internal/canary/`, or `internal/contract/` names the branch-native census test.
  The charge runs that test in the focused checks.
- A charge that shells out names the branch-native census's allowed process seams.
  The charge never backgrounds a test run.
- A charge that binds `PATH` or the process environment includes the ceiling
  file `internal/worktree/parallel_census_test.go` in its fence.
- A charge lists `bench test --check skip-ownership` in its focused checks
  when the charge adds a test that can skip.
- A conformance registry can name a file by path. A charge that moves such a file
  runs that check with `bench test --check` in the focused checks.
- A charge that moves or reflows an anchor runs the fixture-bite check in the
  focused verification list.
- A charge that adds an anchor names `bench test --check <owning-check>` as
  its probe.
- A grammar charge enumerates two inventories in its fence:
  the shared fixture owners and the exact-record assertion families.
- A charge that adds a live-tree test fences the live-tree inventory file
  `internal/conformance/tier_test.go`.

## Isolation and end of life

- Release is the creating request's default end of life for a delegate worktree.
  `bench worktree clean` and then `bench resume-clean` is the recovery pair.
- A fold between owned worktrees runs through `bench worktree merge`, never
  through a file copy.
- A large uncommitted build that no worktree can hold may run in the main checkout
  under exactly four conditions. The conditions are one writer, a named file
  allowlist, no commit authority, and a `git status` check verified on return.

## Probes

- Probe a tracked file that has pending changes with a copy aside.
  `git checkout --` wipes the whole in-flight diff, not only the probe.
- Before the coordinator reads a probe verdict, the coordinator confirms the
  mutated bytes against the copy aside.
- A probe of a live-tree anchor or marker runs through the gate, not through the
  conformance package alone. Only the gate's entry test grades the live tree.
  `bench test --package ./internal/conformance` is not the root conformance
  pass.
- A coordinator probe writes through the writer the test captures, not through an
  ambient stream the test ignores.
- A delegate's reading of a gate skip or a gate-adjacent signal is a claim, not a
  result. Plant the matching break and run the oracle before you believe either
  reading.

## Retry stops and aggregate readiness

- After the second known-flaky refusal proves green in isolation, stop
  coordination and hand both results to the reviewer.
- Before aggregate grading, wait until returned delegates have no live tests and
  serialize the coordinator-owned resource.

## Before the landing

- A running `bench commit` is active until its process exits. A reported red does
  not authorize a tree edit while a later phase or a cleanup can still run. The
  coordinator waits for the terminal exit before diagnosis or repair.
- The coordinator reads a background commit by its `committed` line.
- The coordinator compares the destination tip with the frozen base. The
  coordinator grades each destination change that a new oracle check reads.
- A hand verification does not close an acceptance row without a named
  red-capable test.
- A ticket that returns without a pre-edit red for each row goes back to the
  delegate for those reds. The coordinator gets the reds before the commit.
- Keep an accepted finding on its original ticket when attribution is clear. Use
  an umbrella repair ticket only for a genuinely shared owner.
- When an installed lane cannot commit its repair, run the same ordinary commit core from
  the candidate tree. Grade the composed snapshot, then require the sanctioned rebuild after landing.
- The checkpoint ticket records a reviewer decision that occurs during a build.
  Later review reads that decision.

## In a review round

- A fan-out is reported as launched only after a liveness check, never after a
  successful spawn.
- A re-review charge names the folds, and asks only for findings above the
  blocking bar.
