# Render Go failures and result postures

Blocked by: Run fresh Go tests and render packages

## What to build

The structured test reader retains actionable direct-test and package diagnostics
and returns the declared AXI output and exit posture for success, test failure,
build failure, malformed streams, missing tools, and unmatched packages.

## Acceptance

- [x] Direct failing leaves retain their first diagnostic, parent-only aggregates
  are suppressed, and package/build failures use an empty test cell.
- [x] Passing runs exit 0; test and build failures render all three trustworthy
  tables and exit 1; usage/help and start, no-package, and malformed-stream
  outcomes use their declared stdout and exit codes with empty stderr.
