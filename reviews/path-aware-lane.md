# Review: path-aware-lane

Frozen pair: base `7fa3023614e585d879c283ef5e74c8d87b91f4d0`, reviewed tip
`fc39350fe020119f5410bb35eb8c281fe198a925`. Three axes ran on opus / medium.
The repair ticket `specs/path-aware-lane/tickets/repair-review-findings.md`
closed S1, S2, S4, S5, P1, F3, and F4. Every finding below needs a reviewer
decision.

## Standards

Findings: 1. Worst: two structure budgets are over with no grant.

- S3 `ask-user` — `internal/gate/lane_test.go` grew to 447 lines, and
  `internal/gate/` grew to 45 files. No grant in `.bench/structure-accept`
  names either. The reviewer decides the grant or the split.

## Spec

Findings: 1. Worst: the source widened the spec's Ownership fences.

- P2 `ask-user` — the source amended the spec's Ownership fences with
  `internal/conformance/tier_test.go`. The reviewer confirms it at sign-off.

## Coverage

Findings: 3. Worst: the real `test --check` hop is asserted by argv
comparison and never run by a test.

- F1 `ask-user` — `landing.Owner.Merge` has no tree-equality guard. A merge
  whose composed tree equals the previous tip reaches the authority with an
  empty change list, and `selectLaneChecks` refuses it. The reviewer decides
  pass or refuse. Recommendation: an empty list selects every declared check.
- F2 `ask-user` — no test runs the real run binary through `test --check`.
  The spec's Testing decisions chose the stub factory. The build session ran
  the hop live twice: a dry run over `roadmap/FT215.md` and the ticket's own
  lane commit. The reviewer decides whether a real-build row is worth its cost.
- F5 `ask-user` — an unknown path now selects nine checks, and five of them
  are `bench test --check` runs. A shell or JSON commit therefore costs more
  than the old lane. PL16 decided the rule; the cost is unmeasured. Recommendation: accept
  for this landing and measure in the retro.
