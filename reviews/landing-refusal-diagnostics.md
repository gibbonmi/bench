# landing-refusal-diagnostics review pickup

Frozen base: `850fd2677c4fd56ff2a0e8f2f1c5ef1698ba1148`

Reviewed tip: `5884f5f9ba7c87c111cb475de366c9e2963dc289`

## Standards

Finding count: 0. Worst issue: none.

## Spec

Finding count: 0. Worst issue: none.

## Coverage

Finding count: 2. Worst issue: P2 omission gap.

- **auto-fix — a genuine 64-hex caller token has no command-seam control.**
  The edge inventory requires a real caller token that happens to be 64 hex to
  digest and authenticate normally (`specs/landing-refusal-diagnostics/spec.md:318`).
  `TestLandCommandRefusesStoredRequestDigest` exercises only the stored digest
  (`internal/worktree/land_test.go:1208`), so an overbroad implementation that
  rejects every digest-looking token would stay green. Add a successful land
  control whose original request token is exactly 64 hexadecimal characters.

- **auto-fix — LR19 does not prove the escaped path remains named.**
  LR19 requires the hostile offending path to be carried without raw controls
  (`specs/landing-refusal-diagnostics/spec.md:285`).
  `TestLandCommandRefusalKeepsControlBearingPathInOneTableRow` checks safety,
  the table header, and row count (`internal/worktree/land_test.go:1414`) but
  not the escaped path cell. A safe constant could replace the path and pass.
  Assert the exact escaped newline, ESC, and comma representation.
