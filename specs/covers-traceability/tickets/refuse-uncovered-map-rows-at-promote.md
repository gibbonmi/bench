# Refuse uncovered map rows at promote before the gate spends

Blocked by: refuse-covers-violations-at-assign.md
Ownership fence: `internal/specbuild`
Integration surfaces: assignment Ticket/TicketDigest record→existing internal/specbuild/state.go fields exercised by PR3; recording gate double→existing promotionGate fixture exercised by PR1; refusal-without-spend helper→existing internal/specbuild/promotion_recompose_test.go requirePromotionRefusal
Contracts: the union of non-local covers across the integrated assignments' re-parsed tickets crossing the ticket parse→promote totality within `internal/specbuild`, asserted by PR1 and PR4 against real ticket files on disk

## What to build

For an opted-in spec, `bench spec build promote` refuses a composition that
leaves any map row ID with zero non-`local` covers, naming the uncovered IDs,
before the prospective gate owner executes.

Totality runs in `Service.Promote` immediately after the integrated-and-released
loop over `run.Assignments` and before `owner.Execute`: re-parse exactly the
tickets the integrated assignments record (`assignment.Ticket`), refuse on a
`Digest` mismatch against `assignment.TicketDigest` (the integrate.go re-parse
precedent), and union their non-`local` covers against the coverage export's
map IDs. The tickets directory is never globbed, so an unassigned decoy file
contributes nothing. Over-coverage is legal; only a zero-cover ID refuses. A
non-opted-in spec promotes exactly as today. Refusal tests go through
`requirePromotionRefusal`, which already asserts zero gate executions and
unchanged state; `reviewedPromotionFixture` and `changedTicket` are the prior
art for the composed run and the post-integration edit.

## Acceptance

- [ ] [PR1] promote refuses when any map ID has zero non-local covers across the integrated assignments' tickets, naming the uncovered IDs, with the recording gate at zero executions and state unchanged
- [ ] [PR2] an unassigned decoy ticket file in the tickets directory contributes nothing to totality
- [ ] [PR3] a post-integration covers edit to an assigned ticket refuses at promote on the digest mismatch
- [ ] [PR4] full coverage promotes, over-coverage does not refuse, and a map ID covered only by local rows refuses

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PR1 | move the totality check after the gate owner executes | the uncovered-ID refusal test | apply, run `go test ./internal/specbuild -run TestPromote`, expect the zero-executions assertion failure |
| PR2 | compute totality over a tickets-directory glob | the decoy-ticket test | apply, run `go test ./internal/specbuild -run TestPromote`, expect the decoy-satisfies failure |
| PR3 | trust the assignment's recorded rows without re-parsing the ticket | the digest-mismatch test | apply, run `go test ./internal/specbuild -run TestPromote`, expect the accepted-edit failure |
| PR4 | refuse when two rows cover one map ID | the over-coverage test | apply, run `go test ./internal/specbuild -run TestPromote`, expect the refused-over-coverage failure |
