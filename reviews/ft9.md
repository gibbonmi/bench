# Review findings — ft9 (bench diff --full)

Reviewed commit: 51f5075 against the ft9 spec (retired; recover with
`git log --grep='spec-retire: ft9'`). Advisory; the gate and the reviewer own
done-ness.

## Standards

1 finding. Worst: the help/doc duplication below.

- **low — duplicated knowledge, reviewer call.** The `--full` contract facts
  (two-dot log vs three-dot diff; raw-not-TOON passthrough) are derived twice:
  the runtime help const `fullHelp` (`internal/diff/diff.go:82-86`) and the
  package doc comment (`internal/diff/diff.go:9-15`). AGENTS.md "one source per
  fact" names an enforcement and its advertisement as the defect class.
  Borderline: help text and doc comment serve different audiences and resist a
  shared source, so it may qualify as allowed honest repetition. Decide, and
  either collapse or accept.

## Spec

0 findings. All 7 acceptance-coverage rows verified covered with their promised
red signals (`internal/contract/axi/axi_wave2_test.go`,
`internal/conformance/docs_workflow_helpers_test.go:38`); all implementation
decisions honored; no scope creep. Mirror parity holds structurally —
`.claude/commands` is a symlink to `../.agents/commands`.

## Coverage

2 findings. Worst: the control-char subject below.

- **med — control-byte commit subject hard-fails the whole command.** A commit
  subject containing a control byte other than tab/newline/return (e.g. ESC
  `\x1b`) survives `git log --format=%s`, reaches `toon.Table`, and the library
  refuses (`internal/toon/toon.go:26-29`), so `bench diff --full` returns the
  AXI error at exit 1 — files table, log, and diff body all swallowed even with
  a resolvable base. No test exercises this class (the contract test covers only
  ASCII comma/quote), and the retired ft9 spec's edge-inventory line "log subjects
  escape once" is factually wrong for control bytes. Note: exit-1-on-refusal is
  the same posture a control-char *path* already produces in bare `bench diff`,
  so one fix shape is a test pinning that posture plus a corrected spec line;
  the other is graceful degradation. Reviewer's call which.

- **low — post-resolution git failure yields false-empty sections at exit 0.**
  If `git log base..HEAD` or `git diff base...HEAD` fails after the base
  resolves, `commitLog` returns nil on error (`internal/diff/diff.go:74-77`) and
  the body call drops its error (`internal/diff/diff.go:140`), so the output
  reads as "no commits / no diff" instead of an error. Consistent with the
  file's existing convention for `changedFiles`, and hard to trigger once the
  base is ancestor-validated — flagged so the emptiness-as-success tradeoff is
  a decision, not an accident.
