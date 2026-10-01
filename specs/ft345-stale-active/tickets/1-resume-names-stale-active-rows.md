# The session-start resume names aged active rows stale-active

Blocked by: none
Writes: internal/worktree/lifecyclepolicy/lifecyclepolicy.go, internal/worktree/lifecyclepolicy/lifecyclepolicy_test.go, internal/worktree/classifier.go, internal/worktree/worktree.go, internal/worktree/landed_test.go, internal/worktree/orphan_render_test.go, internal/worktree/orphan_test.go, CHANGELOG.md
Covers: none

## What to build

This ticket fixes the orphan face of FT345. The session-start resume called an
active, unlanded assignment older than the stale window "orphaned". At the same
time, `bench worktree list` showed the ledger state `active` for the same row. The
two surfaces seemed to disagree.

The resume now uses the word `stale-active` for these rows. The retain reason
`stale-active` replaces `orphaned` in the retained counts. The line for each row
starts with `stale-active <id>:`. The line states that the row is active and
unlanded past the stale window, and it keeps the plan-only
`bench worktree clean <path>` command. The `bench worktree list` output, the
ledger, and the drain's orphan definition do not change.

## Acceptance

- [x] An active, aged, unlanded assignment gives `retained landed=1 stale-active=1` and a `stale-active <id>:` line in the resume summary.
- [x] A landed assignment gives `retained landed=1` and no `stale-active` line.
- [x] The control-byte fallback line starts with `stale-active <id>:`.
- [x] The automatic decision table gives the `stale-active` reason for an aged active row.

## Verification

The `internal/worktree/...` packages pass. Before the production edit, the
resume summary test and the fallback test were red on the `orphan` wording, and
the decision table did not compile. A `bench probe` run put the `orphaned` value
back on the reason constant, and the resume summary test and the fallback test
were red. The probe restored the file.
