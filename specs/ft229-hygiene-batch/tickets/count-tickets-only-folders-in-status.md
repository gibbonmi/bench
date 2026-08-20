# Count tickets-only spec folders in the status board

Blocked by: close-the-light-path-ticket-on-landing.md
Writes: internal/status

## What to build

`bench status` gains a housekeeping row reporting how many tickets-only spec
folders `specs/` holds, and naming the command that closes one. The row enters
the housekeeping severity band below the merged-spec retirement and
orphaned-pickup rows, and obeys the board's show-only-on-signal rule, so a count
of zero prints nothing and the five-row budget is untouched.

The tickets-only predicate is the one
`close-the-light-path-ticket-on-landing.md` sets: a direct child of `specs/`
holding no `spec.md`. Read it from that source rather than restating it.

## Acceptance

- [ ] a specs tree holding tickets-only folders renders one row with the count and the closing command.
- [ ] a specs tree with no tickets-only folder renders no such row.
- [ ] the row ranks below the retirement and orphaned-pickup rows.
