# Verify a frozen source tip in preflight

Blocked by: none
Writes: internal/preflight

## What to build

The review phase asks the reviewer for a frozen base and tip, but `bench
preflight` takes only `--base`. Both `review` and `build` gain `--source-tip
<commit>`. The value is verified rather than recorded: preflight already derives
the tip, and a supplied value that disagrees renders a red verdict row. A value
that does not resolve is a grammar-level error, distinct from the mismatch red,
the way `--base` already separates the two. Omitting the flag keeps today's
behavior exactly.

## Acceptance

- [ ] `bench preflight review <slug> --source-tip <commit>` and the `build` form are accepted.
- [ ] a `--source-tip` disagreeing with the derived tip renders a red verdict row.
- [ ] a `--source-tip` that does not resolve returns a grammar error distinct from the mismatch red.
- [ ] omitting `--source-tip` leaves both forms behaving as they do today.
