# Adopt the proportionate ticket-evidence contract and re-anchor workflow canaries

Blocked by: permanent-optional-ticket-inventory.md
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `.agents/skills/bench-craft-delegate/SKILL.md`, `.agents/commands/bench-implement-spec.md`, `internal/anchors/registry_data.go`, `internal/conformance/example_agreement_test.go`, `internal/conformance/fixture_bite_test.go`, `tests/canary/workflow-guidance-anchors/`
Integration surfaces: craft-tickets taught template→`internal/conformance/example_agreement_test.go`; write-delegate evidence contract→`.agents/commands/bench-implement-spec.md`; permanent guidance markers→`internal/anchors/registry_data.go` and `internal/conformance/fixture_bite_test.go`; retained evidence-routing sentences in craft-tickets→`tests/canary/workflow-guidance-anchors/` mutations; each canary mutation→its registry expectation in `internal/conformance/fixture_bite_test.go`
Contracts: the taught ticket example crosses `.agents/skills/bench-craft-tickets/SKILL.md`→`internal/conformance/example_agreement_test.go` with blockers, fence, integration surfaces, contracts, acceptance rows, and covers mapping but no mandatory Closure or Red-mutations section, asserted by MG1; focused delegate checks and coordinator-owned public behavior verification cross `.agents/skills/bench-craft-delegate/SKILL.md`→`.agents/commands/bench-implement-spec.md`, `internal/anchors/registry_data.go`, and `internal/conformance/fixture_bite_test.go` without self-mutation or different-kind probe duties, asserted by MG2; each retained evidence-route clause crosses `.agents/skills/bench-craft-tickets/SKILL.md`→its `tests/canary/workflow-guidance-anchors/` MUTATE.json→its registry expectation in `internal/conformance/fixture_bite_test.go`, asserted by CF1-CF3 against the real materializer

## What to build

Make the permanent lifecycle default honest in its guidance, self-checks, and
canaries as one landing. This ticket is the combined SB9/SB10 repair: the spec
records that the previous two-ticket split cannot land independently — the CF2
fixture edit strands `TestSpecTicketHandoffWorkflowFixturesAreComplete` red in
`internal/conformance/fixture_bite_test.go` while the guidance edit owns the
sentence the fixture materializes. That observed mirror red (debug receipt,
2026-08-11, abandoned run on candidate 9146b5ed) is the specific project-gate
red any thinner cut strands; review re-derives it by attempting the split.

Remove `Closure:` and `## Red mutations` from the taught ticket template and
from the example-agreement requirements. Replace every mandatory-mutation
directive in craft-tickets: discovery traceability, per-spec-row evidence
derivation, the taught template, compound-claim explanation, field/reference
prose, the Good example explanation, and breakdown-review obligations. The
retained policy is acceptance coverage, concrete coverage-map red signals at
genuine TDD seams, focused checks, optional validation for a declared legacy
graph, blocker/fence/contracts grammar, and review-time semantic falsification.

Remove the write delegate's hand-applied self-mutation and the coordinator's
different-kind/different-site mutation duties from craft-delegate, its charge
example, and bench-implement-spec's checkpoint and repair-ticket prose.
Delegates still run focused ticket checks and coordinators still independently
verify public behavior in the exact returned tree before checkpoint. Replace
the `AfterImplementSpec` registry requirements for the template mutation table,
the already-covered subject-mutation route, delegate self-probe, different probe
site, mutation-kind vocabulary, bench-implement-spec's coordinator-owned
different-kind checkpoint, and craft-delegate's global different-kind rule.
Replace fixture-bite cases `probe kind`, `template mutations header`, `delegate
self-probe`, `probe site differs`, `probe kind vocabulary`, and `raw git route`
with bites that go red if focused delegate checks or independent coordinator
public-behavior verification disappear, or if a hand-applied self/different-kind
mutation duty returns. Rebase `raw git route` on the replacement checkpoint
wording while preserving its independent refusal of synthesized lifecycle Git.

In the same landing, re-anchor the three workflow-guidance canaries on the
retained evidence routes only. Narrow each `old` value to the stable
evidence-route clause that remains policy: observed-red rows carry their
failing public operation into ticket acceptance; already-covered rows retain
their named control; not-TDD-able rows map to the first ticket where their seam
exists. Each mutation's replacement drops that whole retained route, its
`EXPECT` names the same omission with the exact retained-route suffix — never a
category-only label — and its registry expectation in
`internal/conformance/fixture_bite_test.go` changes in this same fence.
Preserve the fixture directories, auto-discovery shape, and one exact bite per
fixture. Do not remove canary mutation fixtures, TDD red-green practice, or
review-time falsification.

## Acceptance

- [ ] [MG1] (covers SB9) every craft-tickets discovery, derivation, template, explanation, Good-example, field-reference, and breakdown-review clause plus the example-agreement oracle teaches independently-green tickets through blockers, exact fence, integration surfaces, contracts, acceptance/covers rows, TDD-seam red signals, focused checks, and optional declared-legacy validation without mandatory Closure or Red-mutations fields.
- [ ] [MG2] (covers SB9) craft-delegate and bench-implement-spec require focused delegate checks plus independent coordinator public-behavior verification; the seven named `AfterImplementSpec` records and six named bite cases protect those duties, the already-covered TDD-seam route, and the raw-Git lifecycle refusal while rejecting reintroduced self/different-kind mutation ceremony; canary fixtures, TDD red-green, and review falsification stay unchanged.
- [ ] [CF1] (covers SB10) the observed-red fixture materializes against current guidance by removing the retained public-operation-to-ticket route, and `go test ./internal/conformance -run TestSpecTicketHandoffWorkflowFixturesAreComplete -count=1` proves the exact `dropped the observed-red public-operation route` bite with its registry expectation changed in this same fence.
- [ ] [CF2] (covers SB10) the already-covered fixture materializes against current guidance by removing the retained named-control route and proves the exact `dropped the already-covered named-control route` bite with its registry expectation changed in this same fence; the ambiguous category-only `dropped the already-covered` suffix is red.
- [ ] [CF3] (covers SB10) the not-TDD-able fixture materializes against current guidance by removing the retained blocker-to-first-seam-ticket route and proves the exact `dropped the not-TDD-able first-seam route` bite with its registry expectation changed in this same fence.
