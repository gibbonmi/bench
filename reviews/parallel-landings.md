# Review pickup: parallel-landings

Frozen base `1a135f1b`, reviewed tip `9ffb8278`. Loop 1 of 2, opus / medium,
three read-only axes. Raw findings: 15 (Standards 5, Spec 4, Coverage 6), none
blocking. De-duplicated repair targets: 8. The reviewer pre-approved the
worker's judgement on every `ask-user` item; the dispositions below are final.

## Standards

Count 5. Worst: the new workflow rule is advertised as single-sourced, but the
marker list that enforces the rule does not name it.

- S1 `auto-fix` — `capture/session-handoff.md` pins a worktree label, not the
  branch ref and tip (AGENTS.md, phase-close handoff). Repair: rewrite the
  handoff at this phase boundary with the branch ref and the tip.
- S2 `auto-fix` — `WorktreeRuleMarker` (`internal/anchors/registry_data.go`)
  is absent from the marker list in
  `internal/conformance/validity_checks_test.go:224-234`. Repair: add it; the
  spec fence gains that test file by amendment on this source.
- S3 `no-op` — the export of `WorktreeRuleMarker` earns its keep once S2 lands.
- S4 `auto-fix` — `AGENTS.md:60-62` restates the handoff→source rule a third
  time. Repair: keep "it lands with the worktree"; drop the restated rule.
- S5 `no-op` — the `unionStage = 0` sentinel is guarded by `ok` upstream.

## Spec

Count 4. Worst: the debug skill's isolation rule lost its rationale.

- P1 `auto-fix` — `.agents/commands/bench-debug.md:14-15` dropped the
  "dirty in an unattributable way" reason. Repair: restore one short reason
  sentence; the file stays at or under 170 lines.
- P2 `no-op` — the tickets-only close skips the fence; a folder with no
  `spec.md` declares no fence. Drain note: one spec sentence.
- P3 `no-op` — WL5 compares to the source tree; the stray-transition failure is
  pinned separately.
- P4 `no-op` — WL9 is a freeze on pre-existing tests.

## Coverage

Count 6. Worst: the surface union test asserts by containment.

- C1 `auto-fix` — `landingConflictNext`'s unsafe-spec branch (`<spec>`
  placeholder) has no test. Repair: one test with a control-byte spec slug.
- C2 `no-op` — the `<full-destination-commit>` guard is unreachable and harmless.
- C3 `auto-fix` — `TestLandCommandComposesCaptureOntoMovedDestination` asserts
  the union by containment. Repair: assert the exact published bytes.
- C4 `auto-fix` — `TestWorktreeRuleAnchorsRedOnRemoval` never proves the
  `Workflow` section binding. Repair: add a case with the marker under another
  section.
- C5 `auto-fix` — the `capture/` prefix boundary has no refusal row. Repair:
  one row with `capture.md` or `capturex/` refusing.
- C6 `no-op` — both-sides-deleted union is unreachable; a D/F conflict under
  `capture/` is an undecided edge. Drain note.
