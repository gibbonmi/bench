# Reject duplicate acceptance IDs in ticket parsing

Blocked by: none
Ownership fence: `internal/specbuild`
Assumptions: acceptance IDs are the stable identity for behavioral obligations; duplicates are detected after R-range expansion, where overlap becomes visible; ParseTicket callers relay its error unwrapped, so the diagnostic itself names the ticket

## What to build

`ParseTicket` refuses a ticket whose acceptance rows repeat an ID — literally
or through overlapping R-range expansion — with a diagnostic naming the
duplicate ID and the ticket file, instead of silently collapsing the rows.
Fence and assumption deduplication stay as they are: only acceptance IDs carry
per-obligation identity.

## Acceptance

- [x] [DA1] a ticket repeating an acceptance ID (literal repeat, or an R-range overlapping another row) is refused with an error naming the duplicate ID and the ticket file basename.
- [x] [DA2] a ticket with distinct IDs — including a valid R-range beside literal IDs — still parses with its row order unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DA1 | restore silent `unique()` collapse of parsed rows | the duplicate-ID refusal test in `parse_ticket_test.go` | run `go test ./internal/specbuild`, expect the accepted-duplicate failure |
| DA2 | compare pre-expansion raw tags so a valid range false-positives | the existing range-expansion test | run `go test ./internal/specbuild`, expect the row-set mismatch failure |
