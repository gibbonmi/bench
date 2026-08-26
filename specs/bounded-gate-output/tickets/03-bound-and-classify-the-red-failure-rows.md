# Bound and classify the red failure rows

Blocked by: 02-buffer-each-phase-stream-and-print-the-red-table.md
Writes: internal/gate/report.go (new in 02), internal/gate/report_test.go (new in 02), internal/gate/runner.go, internal/gate/capability_skips.go, internal/gate/capability_skips_test.go, internal/sanitize/sanitize.go, internal/sanitize/sanitize_test.go

## What to build

The red table stays bounded and names only real failures.

Each phase prints at most twenty rows. A phase with more rows prints twenty
rows, then one row `+<k> more lines: <path>`. The path names the complete
stream file. A phase with exactly twenty rows prints no extra row. The cap is
twenty so a one-phase red fits the stop hook's thirty-line tail.

This contract crosses into ticket 05:

- `(*phaseStreams).path() string` returns the complete stream file path, or an empty string.

Ticket 05 makes the path real. Until then the path is empty. An empty path
renders the row `+<k> more lines (stream unavailable)`.

A phase that a red need skipped adds no row, because a consequence is not a
failure. Each cell passes a new control-byte filter in `internal/sanitize`
that removes every byte below 0x20 except tab and escapes nothing. The
existing `sanitize.Controls` escapes a backslash, and the TOON encoder
escapes it again, so a path would reach the reader with four. A long line
renders whole, because the cap counts rows and not bytes. The package
comment in `internal/sanitize/sanitize.go` names two duties today; add the
third.

`reportCapabilitySkips` stops writing its red diagnoses to stderr. It hands
the report one row per diagnosis under the phase name `capability`:

- each environment skip, with its test name and reason
- a strict-mode capability skip, with its count
- an unreadable skip log, with its error

Read one tree fact first: `environmentFailure` reds an environment skip on every run, not
only under strict mode. Keep that posture.

## Acceptance

- [ ] A red phase with fifty failure lines prints twenty rows and one `+30 more lines: <path>` row. (BG07)
- [ ] A phase skipped because its need went red adds no row. (BG09)
- [ ] An environment skip prints a row under phase `capability`. (BG10)
- [ ] A strict-mode capability skip and an unreadable skip log each print a row under phase `capability`. (BG33)
- [ ] A failure line with an ESC byte renders with the byte removed. (BG14)
- [ ] A failure line with one backslash renders with one backslash after TOON decoding.
- [ ] A red phase with exactly twenty lines prints twenty rows and no more-row.
- [ ] A capability-class skip outside strict mode enters no row.
