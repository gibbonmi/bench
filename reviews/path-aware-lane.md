# Review: path-aware-lane

Frozen pair: base `7fa3023614e585d879c283ef5e74c8d87b91f4d0`, reviewed tip
`fc39350fe020119f5410bb35eb8c281fe198a925`. Three axes ran on opus / medium.
The repair ticket is `specs/path-aware-lane/tickets/repair-review-findings.md`.

## Standards

Findings: 6. Worst: two guidance documents advertise a cost the lane does
not deliver.

- S1 `auto-fix` — `docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md`
  and `CHANGELOG.md` say a Markdown commit "pays no Go toolchain". The prose
  check runs through the private run binary, which the Go toolchain builds.
  Repair: say the commit runs no `vet` and no `build` check. Rule: invariant 3.
- S2 `auto-fix` — the file comment in `internal/gate/lane_select.go` argues its
  own placement. Repair: keep the first sentence only. Rule: `craft-comments`.
- S3 `ask-user` — `internal/gate/lane_test.go` grew to 447 lines, and
  `internal/gate/` grew to 45 files. No grant in `.bench/structure-accept`
  names either. The reviewer decides the grant or the split.
- S4 `auto-fix` — the four content classes name lane checks by string with no
  test that binds them to `BenchkitLane`. Repair: one assertion that every
  class's checks are a subset of the lane's names.
- S5 `auto-fix` — `packagesurface.EmbedTargets` errors when `cmd/` or
  `internal/` is absent, so three fixtures carry a placeholder Go source.
  Repair: an absent directory contributes nothing; remove the placeholders.
- S6 `no-op` — `changeClasses` reads a gitlink mode on either side. That is a
  safe superset of the spec's wording.

## Spec

Findings: 6. Worst: the PL5 edge line promises a glob-character case no row
exercises.

- P1 `auto-fix` — the Edge inventory line for PL5 names "a glob character";
  the row and its test carry `café notes.md` only. Repair: drop the clause.
- P2 `ask-user` — the source amended the spec's Ownership fences with
  `internal/conformance/tier_test.go`. The reviewer confirms it at sign-off.
- P3 to P6 `no-op` — the gitlink superset, the `EmbedTargets` error path, the
  literal PL15 expectation, and the empty-selection refusal. Each is refuted
  by the spec text or by the tree. All 46 audited rows are closed.

## Coverage

Findings: 5. Worst: the real `test --check` hop is asserted by argv
comparison and never run by a test.

- F1 `ask-user` — `landing.Owner.Merge` has no tree-equality guard. A merge
  whose composed tree equals the previous tip reaches the authority with an
  empty change list, and `selectLaneChecks` refuses it. The reviewer decides
  pass or refuse. Recommendation: an empty list selects every declared check.
- F2 `ask-user` — no test runs the real run binary through `test --check`.
  The spec's Testing decisions chose the stub factory. The build session ran
  the hop live twice: a dry run over `roadmap/FT215.md` and the ticket's own
  lane commit. The reviewer decides whether a real-build row is worth its cost.
- F3 `auto-fix` — no row asserts a change set that mixes `unknown` with a
  known class. Repair: one `TestSelectLaneByClass` row.
- F4 `auto-fix` — no row pins the prefix boundaries `roadmapx/a.md`,
  `docs/ROADMAP.md`, `specs/decisions/x.md`, and `capture/retros` as a file.
  Repair: four `TestSelectLaneByClass` rows.
- F5 `ask-user` — an unknown path now selects nine checks, and five of them
  are `bench test --check` runs. A shell or JSON commit therefore costs more
  than the old lane. PL16 decided the rule; the cost is unmeasured. Recommendation: accept
  for this landing and measure in the retro.
