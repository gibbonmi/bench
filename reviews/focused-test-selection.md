# Focused test selection review

Frozen base: `0dae83df5372727d608c6fb2b29fa118d732562f`

Reviewed tip: `38074d760fc7204a8cf17995d1af69fe527f1bef`

Raw findings: 15. De-duplicated repair targets: 14. The Spec and Coverage
findings for changed-mode `go list` cancellation name the same repair.

## Standards

Count: 1. Worst issue: P2 duplicated fixture-graph knowledge.

- **P2 — auto-fix — single-source the changed-selection fixture graph.**
  `AGENTS.md`'s code standard requires one source per fact and specifically keeps
  fixture harnesses single-sourced. `internal/testreport/selection_test.go:217-234`
  manually defines the production, test, external-test, and embed graph, while
  `internal/testreport/selection_test.go:237-256` encodes the same relationships
  again in generated fixture files. The tests at lines 17-18 and 42-43 inject the
  manual graph; the live-loader test uses a different fixture at lines 143-164.
  Keep only filesystem-shape facts in `changedSelectionModule` and one source for
  the injected graph facts.

## Spec

Count: 3. Worst issue: P1 deleted non-Go subjects violate C06.

- **P1 — auto-fix — make a deleted non-Go-only subject explicitly empty.**
  `specs/focused-test-selection/spec.md:44-49` and C06 at line 249 require a
  proven non-Go subject to exit 0 with zero-row package, failure, and skip tables.
  `internal/testreport/selection.go:66-82` marks a deleted `README.md` absent,
  then lines 127-149 refuse it as outside the current graph. The C06 test at
  `internal/testreport/selection_test.go:94-108` covers only an added README.

- **P1 — auto-fix — scrub conformance controls in every focused mode.**
  `specs/focused-test-selection/spec.md:157-160` requires every mode to remove
  ambient root, tier, scope, selected-set, inherited-set, and capability-log
  controls. `internal/testreport/command.go:145` removes only the capability log
  for default, package, and changed runs; lines 153-172 perform the full scrub
  only for named checks. Cover all ordinary modes and retain the exact named-check
  environment.

- **P2 — auto-fix — cancel and drain the changed-mode `go list` process group.**
  N02 at `specs/focused-test-selection/spec.md:257` and the edge inventory at
  lines 284-286 require process-group cancellation for both Go hops.
  `internal/testreport/selection.go:85-90` uses `exec.CommandContext` without a
  process group or descendant drain. The existing cancellation test exercises
  run-binary construction, not graph loading. This is the same repair target as
  Coverage finding 1.

## Coverage

Count: 11. Worst issue: P1 changed-mode `go list` cancellation can leak a
descendant.

- **P1 — auto-fix — prove changed-mode `go list` descendant drain.** Cancel a
  changed run while its list child owns a descendant; N02 and
  `spec.md:284-286` require the common interruption posture, no survivor, and no
  partial tables. `selection.go:85-90` lacks the group behavior and
  `cancel_test.go:36-121` covers only executable construction. This de-duplicates
  with Spec finding 3.

- **P2 — auto-fix — prove focused `go test` descendant drain.** Cancel the Go
  test hop while it owns a descendant. `internal/testreport/command.go:175-208`
  implements group signalling, but no test exercises that path; N02 requires a
  mutation removing it to turn red.

- **P2 — auto-fix — independently prove the staged-path producer.** C01's
  `TestChangedSubjectDefaultLiveComposition` at
  `internal/diff/explicit_base_test.go:118-136` reuses `tracked.txt` for staged
  and unstaged state. Give committed, staged-only, tracked-worktree, and untracked
  producers distinct paths so omission of the index source turns red.

- **P2 — auto-fix — cover every module-wide metadata name.** C05 applies to
  `go.mod`, `go.sum`, `go.work`, and `go.work.sum` as implemented at
  `selection.go:176-178`; `TestChangedPackageSelectionMetadataAndMixedUnion` at
  lines 41-55 exercises only `go.mod`.

- **P2 — auto-fix — cover both control-byte postures.** The edge inventory at
  `spec.md:266-268` requires ESC and BEL to refuse and tab, newline, and return to
  remain one escaped cell or refuse before execution. The C07 matrix at
  `selection_test.go:57-81` tests only ESC.

- **P2 — auto-fix — cover every named special-node and symlink state.** The edge
  inventory at `spec.md:279-281` names FIFO, socket, device, dangling symlink,
  and live symlink. The C07 matrix covers only FIFO and a live symlink.

- **P2 — auto-fix — exercise space and glob paths end to end.** The edge
  inventory at `spec.md:264-265` requires package and embed fixtures containing
  both; no focused-selection test currently does.

- **P2 — auto-fix — exercise Git's no-renames projection through selection.**
  `spec.md:276-278` requires deletion and addition to feed classification. Generic
  diff rename coverage exists, but no `ResolveChangedSubject` or changed-command
  test traces both halves into focused selection.

- **P2 — auto-fix — prove `--run` reaches the complete changed-package union.**
  C08 at `spec.md:251` requires the unchanged filter on the complete union.
  `TestChangedPackageRunPattern` at `selection_test.go:110-124` selects only one
  package, so partial application would remain green.

- **P2 — auto-fix — cover hostile values and the legacy terminator at the
  focused-command seam.** `spec.md:269-271` requires `--` and flag-looking,
  spaced, bracketed, and shell-operator package/regex values to remain one value.
  Generic parser tests do not authenticate the new focused command surface.

- **P2 — auto-fix — cover missing Go as a structured start refusal.** The edge
  inventory at `spec.md:287-290` requires this focused-mode posture. Existing
  bootstrap tests do not exercise the focused Go hop with `go` absent from
  `PATH`.
