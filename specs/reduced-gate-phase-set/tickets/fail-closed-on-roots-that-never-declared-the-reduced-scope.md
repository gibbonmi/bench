# Fail closed on roots that never declared the reduced scope

Blocked by: none

Ownership fence: `internal/gate/gate.go`, `internal/gate/reduced_run_test.go`, `internal/contract/runtime/runtime_gate_reduced_test.go`
Assumptions: the stripped construction, identity, and verdict class are settled and unchanged

## What to build

Two fail-open holes in the one guard standing between this feature and a repository that
never asked for it.

**The reviewer's ruling (2026-08-01): reduction and stripping apply only to a root that
declares its own scope, or to the kit itself.** Today `ReducedScope()` is applied to every
root the binary gates, so a linked repo materializes a stripped worktree on every full gate
and gets reduced runs against an allowlist it never declared, with no conformance binding of
its own. The spec cut per-repo configurability as out of scope; applying the kit's list to
foreign roots is a third thing neither the story nor the cut chose. Foreign roots keep
today's behavior exactly.

**`phaseTableGate` fails open on manifest presence.** Its comment claims the gate must
*provably* route through the phase table, but a declared `.bench/phases.json` returns true
regardless of what `.bench/gate.sh` actually does — so a hand-written gate that never execs
`gate-phases` gets its oracle swapped for the built-in table. Let the manifest decide *which*
table, never *whether* the gate routes through one.

## Acceptance

- [ ] [RB1] A root that declares no reduced scope of its own runs unsplit and unreduced, exactly as before this feature.
- [ ] [RB2] A root carrying a phase manifest but a gate script that never hands off to `gate-phases` pays a full run.
- [ ] [RB3] The `gate-phases` hand-off branch — the one governing this repository's own reductions — is exercised by a test rather than reached only incidentally.
