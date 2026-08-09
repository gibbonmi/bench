## Outcome

`branch-native-build-test-architecture` published as `3701c4a0420d0af0ac39b5a023cf09cb45972a9e` through the prospective six-phase gate. The lifecycle record was intentionally removed during the gate rebuild, so this exceptional retro relies on the retained prospective run log and final-evidence record rather than `bench spec build status`.

## Gate-stage timings

Prospective run `20260809T083946.876087055Z-2757604` was green in 73.850s: gofmt 0.059s, vet 1.593s, test 66.302s, race 3.427s, system 1.187s, and shellcheck 0.291s. The ordinary test driver is 89.8% of this run, so it is the only material optimization target. A pre-rebuild green run took 354.073s, but it used a different subject and the retired architecture; the apparent 79% reduction is directional evidence, not a controlled performance claim.

## Ticket-versus-spec-slice and delegate performance

The six tickets together replaced the fixture-driven ordinary suite with direct decision tests, one whole-package driver, a bounded system owner, direct canary mutations, and an architecture census. Per-ticket delegate receipts were removed with the rebuilt lifecycle record, so their duration and repair-round attribution are unavailable.

## Coordinator catches

The retained final evidence caught that the prior whole-repository run was not a valid final result and required two repairs before the prospective run. Final verification also found that the rebuilt lifecycle format no longer exposes the prior terminal receipt; this retro records the limitation instead of reconstructing promotion state.

## Repair attribution

| ticket | repair rounds | cause per round |
|---|---:|---|
| 01-expose-branch-native-command-decisions | unavailable | lifecycle record removed during gate rebuild |
| 02-assemble-branch-native-gate-drivers | unavailable | lifecycle record removed during gate rebuild |
| 03-own-bounded-system-journeys | unavailable | lifecycle record removed during gate rebuild |
| 04-move-canary-and-stripped-proofs | unavailable | lifecycle record removed during gate rebuild |
| 05-enforce-ordinary-adapter-budgets | unavailable | lifecycle record removed during gate rebuild |
| 06-contract-legacy-fixtures-and-enforce-census | unavailable | lifecycle record removed during gate rebuild |

## Agent-experience improvements

### Bench CLI

Keep a durable, schema-migratable terminal summary outside the replaceable gate implementation so final-check can discover a promoted run after a gate rebuild.

### Skills

Document an exceptional final-check route for a deliberate lifecycle-format reset: it should require a named retained prospective log, exact published commit, and an explicit limitation statement.

### Process

Record a before-and-after timing pair against the same subject and cache posture before claiming a gate-speed improvement. The retained 73.850s run identifies the ordinary `go test -count=1 ./...` driver as the current dominant cost.
