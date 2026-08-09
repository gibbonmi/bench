# Contract the run-directory lifecycle

Blocked by: introduce-run-scoped-bench-selection.md, own-gate-run-binary.md, own-focused-test-binary.md, propagate-selected-binary-to-nested-gates.md, serialize-phase-tables.md, serialize-primary-stripped-schedule.md, contract-ordinary-build-census.md
Ownership fence: `internal/runbinary/`, `internal/gate/`, `internal/testreport/`
Integration surfaces: selection lifecycle→introduce-run-scoped-bench-selection.md; gate owner→own-gate-run-binary.md; focused-test owner→own-focused-test-binary.md; nested descendants→propagate-selected-binary-to-nested-gates.md; scheduler teardown→serialize-phase-tables.md and serialize-primary-stripped-schedule.md; closed constructor set→contract-ordinary-build-census.md; terminal lifecycle documentation→align-profile-and-changelog.md
Contracts: owner terminal outcome and descendant settlement cross `internal/gate/` or `internal/testreport/`→cleanup in `internal/runbinary/`, membership is canonical-builder failure, green, red, post-owner refusal, timeout, and interrupt, ordering is clean builder failure directly or signal then grace then kill remaining group then reap, followed by removal of executable, seal, and directory, and no terminal path publishes or retains bytes, asserted by CL1 at the real process/lifecycle junction
Closure: CL1/builder-failure-cleanup, CL1/green-cleanup, CL1/red-cleanup, CL1/refusal-cleanup, CL1/timeout-cleanup, CL1/interrupt-cleanup, CL1/descendant-before-remove, CL1/no-executable-survives, CL1/no-seal-survives, CL1/no-later-reuse

## What to build

Finish the contraction with exhaustive lifecycle tests and any teardown repair they expose. Use injected post-owner failures plus real process groups; do not replace process evidence with a defer-only unit assertion. Remove any temporary compatibility fallback left by the expansion and migration sequence.

## Acceptance

- [ ] [CL1] (covers RS8) canonical-builder failure, green, red, post-owner refusal, timeout, and interrupt all reap any descendants before deleting partial executable, seal, and private directory state, and no later process can reuse any run artifact.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CL1/builder-failure-cleanup | return immediately when the canonical builder leaves partial output and fails | runbinary builder-failure test | have the builder create partial executable and seal paths then exit red, and require no child launch plus complete directory removal |
| CL1/green-cleanup | skip cleanup after code zero | lifecycle failure-injection test | complete a green child and require the private parent is absent |
| CL1/red-cleanup | return on ordinary red before cleanup | lifecycle failure-injection test | complete a red child and require the private parent is absent |
| CL1/refusal-cleanup | leak after a post-owner validation refusal | lifecycle failure-injection test | create selection, force pre-child refusal, and require the private parent is absent |
| CL1/timeout-cleanup | return code 124 before cleanup | bounded timeout harness | time out a child and require code 124 plus absent private parent |
| CL1/interrupt-cleanup | return code 130 before cleanup | bounded interrupt harness | interrupt a child and require code 130 plus absent private parent |
| CL1/descendant-before-remove | remove the directory while a descendant still holds or executes the binary | process-group marker harness | leave a descendant alive, cancel, and require its terminal marker precedes directory-removal marker |
| CL1/no-executable-survives | unlink only the directory entry tracked by owner while another run path remains | filesystem inventory test | finish each outcome and require no selected executable inode under recorded run paths |
| CL1/no-seal-survives | remove executable but retain the adjacent seal | filesystem inventory test | finish each outcome and require both recorded paths absent |
| CL1/no-later-reuse | publish the path to cache, dist, or a stable record | consecutive-process test | complete one run, start another process, and require a new build/path with no discoverable prior artifact |
