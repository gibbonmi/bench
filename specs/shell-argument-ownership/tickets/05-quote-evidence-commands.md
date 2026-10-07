# Quote evidence pagination command arguments

Blocked by: 01-quote-tree-build-arguments.md
Writes: internal/preflight/evidencecmd/evidence.go, internal/preflight/evidencecmd/evidence_command_test.go, internal/preflight/evidencecmd/evidence_modes_test.go, internal/preflight/evidencecmd/evidence_summary_test.go, internal/preflight/evidencecmd/evidence_review_test.go, internal/preflight/evidencecmd/evidence_cleanup_test.go, internal/preflight/evidencecmd/shell_arguments_test.go (new)
Covers: S17

## What to build

Migrate evidenceInvocation to sanitize.ShellQuote.
Keep identity, source, cursor, pagination, argument order, and existing sink controls.
Update every approved next-command fixture with the same always-quoted argument policy.

Review chunk: S-B4.
The predecessor supplies the accepted contract named above.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Decode a real next command and recover its original identity, source, and cursor through a POSIX shell.
Exact next-command fixtures change only argument quote bytes.

- [ ] Evidence commands preserve identity, source, and cursor arguments through the surviving owner. (S17).

## Checkpoint verification

- `bench test --package ./internal/preflight/evidencecmd`

## Required omission evidence

Return a bare cursor or bypass the owner at evidenceInvocation.
The exact producer fixture and argv oracle must detect each omission.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
