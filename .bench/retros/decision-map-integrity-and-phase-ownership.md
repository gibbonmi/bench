## Outcome

Decision maps now have one schema owner, graph/readiness/source validation,
read-only AXI projection, ambient count integration, and a gate-bound
49-fixture mutation family. The tracked corpus uses the canonical format,
shaping and spec authoring have distinct ownership, and the implemented spec
and its compiled provenance have retired from `main`.

## Gate-stage timings

The final landing oracle took 145.76s. Its reported critical paths were:
test/worktree 96.024s, race/worktree 8.109s, conformance-suite 115.891s,
direct conformance 28.631s, and contract/artifact-posture 139.224s. The
post-merge retirement oracle took 131.27s; its reported critical paths were
conformance-suite 100.367s and contract/artifact-posture 127.712s. Build,
gofmt, vet, shellcheck, and canary were green but emitted no independent stage
duration.

## Ticket-versus-spec-slice and delegate performance

The build landed through independently-green schema, validator, corpus,
projection, gate, and guidance tickets, plus one concurrent-corpus migration
and three evidence-driven repair tickets. Fresh write delegates kept each
implementation diff attributable and focused checks cheap; the full-gate
boundary caught integration state without forcing delegates to serialize on
the oracle. Guidance needed the top line because ownership prose crossed
commands, skills, cold-start docs, and mutation anchors.

## Coordinator catches

Coordinator verification caught an initially unregistered canary family, a
deletion-proof fixture rooted at the wrong synthetic topology, a legacy map
committed concurrently after the corpus inventory, and stale source facts in
the compiled map. Fresh-session dogfood found the source drift. Three-axis
review found an ignored second source locator, incomplete cycle diagnostics,
and an empty-fog shaping map whose detailed projection disagreed with its
count. A reported stale CLI defect was an ignored development binary and was
refuted by rebuilding it and rerunning the AXI contracts.

## Agent-experience improvements

### Bench CLI

Allow `bench diff --full` to accept an explicit base/head range so a landed
multi-commit feature can be reviewed without reconstructing the bundle with
raw Git. Revisit the `bench commit --spec` plus immediate default-branch
retirement sequence: it currently requires two full oracle runs for two
short-lived trees, the friction already tracked by FT113.

### Skills

Make review charges state that defaulted decisions remain authoritative when a
narrower acceptance row covers only one representative case. Add a reminder
that diagnostic wording and its exact canary expectation may require one
explicitly justified atomic cross-fence repair.

### Process

Keep fresh-session map-to-spec dogfood before semantic review; it exposed stale
provenance that structural canaries could not. Rebuild the main checkout's
ignored development binary before review so environmental staleness is not
reported as shipped behavior. When a review repair is claimed closed, compare
the claim against both the coverage row and the spec's defaulted-decision
table before ending the terminal repair pass.
