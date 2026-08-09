# Own bounded system journeys

Blocked by: 01-expose-branch-native-command-decisions.md, 02-assemble-branch-native-gate-drivers.md
Ownership fence: `internal/systemtest/`
Integration surfaces: inherited selected executable→`internal/systemtest/`; linked/install, lifecycle/reload, and stripped repository families→`internal/systemtest/`; canary and stripped journeys→04-move-canary-and-stripped-proofs.md; command observation producer→`cmd/bench/` plus 01-expose-branch-native-command-decisions.md
Contracts: the cleaned absolute selected executable path, inode, and digest cross the gate owner→`internal/systemtest/` in stable order with absence refused, asserted by SI1 against the inherited producer; repository and process records cross each journey→suite cleanup with owner-only membership and terminal drain semantics, asserted by SB1 and SJ1; the command implementation observation crosses `cmd/bench/`→the real-binary composition journey, asserted by CP1
Closure: SI1/path, SI1/inode, SI1/digest, SB1/three-repositories, SB1/owner-only, SB1/green-cleanup, SB1/red-cleanup, SB1/interrupt-cleanup, SB1/timeout-cleanup, SJ1/wrapper, SJ1/install, SJ1/freshness, SJ1/reload, SJ1/interrupt, SJ1/timeout, SJ1/descendant-drain, CP1/real-dispatch, CP1/same-implementation

## What to build

Create the tagged system-test package and its single owner. It validates one inherited selected executable, owns at most three repository families and every child process, and proves only the irreducible wrapper, installation, freshness, serialized-state, interrupt, timeout, teardown, and command-composition observations.

## Acceptance

- [x] [SI1] (covers SY1) one tagged system owner validates and propagates the same cleaned selected-executable path, inode, and digest to every real-binary journey.
- [x] [SB1] (covers SY2) the complete system run constructs at most three repositories through its owner and removes repositories plus drained descendants after green, red, interrupt, and timeout outcomes.
- [x] [SJ1] (covers SY3) the system owner proves wrapper routing, installation, freshness refusal, serialized-state reload, interrupt, timeout, and process-group descendant teardown exactly once.
- [x] [CP1] (covers CP1) a selected-executable journey reaches the same production command implementation observed by direct command tests while the owner accounts for every repository and process.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SI1/path | substitute a second executable path for one journey | system identity ledger | mutate the journey selection, run the tagged suite, expect path identity mismatch |
| SI1/inode | replace the selected executable between journeys | system identity ledger | mutate the executable inode, run the tagged suite, expect inode identity mismatch |
| SI1/digest | change selected executable bytes between journeys | system identity ledger | mutate the bytes, run the tagged suite, expect digest identity mismatch |
| SB1/three-repositories | request a fourth repository | system materialization budget | mutate the journey plan, run the tagged suite, expect repository budget exceeded |
| SB1/owner-only | construct a repository inside an individual journey | architecture census | add the constructor, run the census, expect non-owner constructor diagnostic |
| SB1/green-cleanup | retain the green journey repository | terminal cleanup test | mutate cleanup, run the tagged suite, expect retained repository diagnostic |
| SB1/red-cleanup | retain the red journey repository | terminal cleanup test | mutate cleanup, run the tagged suite, expect retained repository diagnostic |
| SB1/interrupt-cleanup | retain a repository after interrupt | terminal cleanup test | mutate cleanup, run the tagged suite, expect retained repository diagnostic |
| SB1/timeout-cleanup | retain a repository after timeout | terminal cleanup test | mutate cleanup, run the tagged suite, expect retained repository diagnostic |
| SJ1/wrapper | route the wrapper to the source checkout | wrapper journey | mutate wrapper target, run the tagged suite, expect selected-route mismatch |
| SJ1/install | omit the installed by-path CLI | install journey | mutate installation input, run the tagged suite, expect missing installed route |
| SJ1/freshness | accept a stale selected executable | freshness journey | mutate the seal input, run the tagged suite, expect freshness refusal missing |
| SJ1/reload | reuse in-memory state instead of a fresh process | reload journey | mutate the reader, run the tagged suite, expect serialized-state mismatch |
| SJ1/interrupt | ignore interrupt delivery | interrupt journey | mutate signal forwarding, run the tagged suite under a bound, expect interrupt failure |
| SJ1/timeout | let the child outlive the timeout | timeout journey | mutate timeout teardown, run the tagged suite under a bound, expect descendant retained |
| SJ1/descendant-drain | let the leader exit before an undrained descendant | process-group journey | mutate the helper, run the tagged suite under a bound, expect undrained descendant diagnostic |
| CP1/real-dispatch | route the selected executable around production dispatch | command composition journey | mutate the dispatch target, run the tagged suite, expect observation absent |
| CP1/same-implementation | report a different command implementation ID | command composition journey | mutate the production observation, run direct and tagged checks, expect identity mismatch |
