# Fan out the staleness pass and let the orchestrator amend the spec

Blocked by: none
Writes: .agents/skills/bench-implement-spec/references/staleness-pass.md, .agents/commands/bench-implement-spec.md, .bench/BENCH.md, internal/anchors/registry_chunk_chain.go, CHANGELOG.md
Covers: none

## What to build

The reviewer directed this change on 2026-10-06, after the FT290 staleness audit.
That audit used five cheap-tier read-only delegates. One delegate audited the spec
body and the decision map, and each other delegate audited two or three tickets.
The orchestrator then reviewed the returns and decided each finding. The
build-entry staleness pass of `/bench-implement-spec` now works in the same way.

The staleness-pass reference changes from one mid-tier charge with a 15-call
budget to a fan-out of cheap-tier delegates. Each delegate audits one slice and
returns a verdict for each item and a row for each stale item. The orchestrator
reviews the returns and decides each finding with no reviewer stop.

The phase file drops the delegated review round that confirmed the amendment.
The preflight still reruns green before the version 2 plan amendment.
`.bench/BENCH.md` lets the orchestrator decide a behavioral contradiction in this
pass and flag it for reviewer veto. The plan-expansion limits stay unchanged:
unrelated scope, a material acceptance change, or a weakened guarantee still
stops the phase.

## Acceptance

- [ ] The staleness-pass reference states the cheap-tier fan-out, the slice
      partition, the delegate procedure, the return shape, and the orchestrator
      review.
- [ ] The phase file charges that pass, keeps the preflight rerun, and has no
      delegated confirmation round.
- [ ] The `.bench/BENCH.md` predicate lets the orchestrator decide a behavioral
      contradiction in the staleness pass and flag it for reviewer veto.
- [ ] The chunk-chain anchor pins the new charge line, and the gate is green.
