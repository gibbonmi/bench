# Refuse a planned seam cell at the chunk freeze

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/coverage/planned_rows.go (new), internal/coverage/planned_rows_test.go (new), internal/reviewrecord/plan.go, internal/reviewrecord/recordcmd/chunk.go (new), internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/completion.go, internal/reviewrecord/recordcmd/planned_cell.go (new), internal/reviewrecord/recordcmd/planned_cell_test.go (new)
Covers: TF47, TF48, TF49, TF50, TF51, TF52, TF53, TF54, TF55, TF56, TF57, TF58, TF59, TF67

## What to build

`internal/coverage` exports a projection over spec bytes. The projection returns the row IDs whose seam cell matches `plannedRe`, and it reads opt-in maps only. A seam cell that starts with `review-owned:` is never projected. The projection reuses the existing parse and the existing marker expression, and it adds no second spelling. `coverage.go` holds 599 lines with no structure grant, so the projection goes into `planned_rows.go`.

`bench record chunk` reads the spec and the plan from the tip commit before it writes. It refuses when any row of the chunk's planned rows holds the marker. The refusal names each such row and the citation form `` `<path>_test.go` (`<Name>`) ``, and it writes no record. `bench record completion` applies the same refusal to every row of every plan chunk at its source. `bench coverage --check` keeps its current behavior.

`internal/coverage` already depends on `internal/reviewrecord`, so only `recordcmd` imports `internal/coverage`. The unexported `atSource` in `internal/reviewrecord/source.go` returns a digest and a plan, not the spec bytes. So export `ReadSpecAtCommit(root, spec, commit string) ([]byte, Plan, error)` in `internal/reviewrecord/plan.go`. It composes the unexported `sourceTree` of the same package with `benchgit.ReadTreeFile` and `ReadPlan` on that one tree. `source.go` needs no edit.

`recordcmd/command.go` holds 399 lines, and the lane `structure` row reds a file over 400 lines that grows. Move the `chunk` form handler into `chunk.go`, as `completion.go` holds the completion handler, before you add the refusal. Put the shared refusal in `planned_cell.go`.

The `internal/coverage` binding row names the command registry files, so this ticket names them. This ticket changes no public grammar, so it edits none of them.

At the ticket commit, rewrite the `planned` seam cell of each covered row into the citation form.

## Acceptance

- [ ] TF47: `bench record chunk` exits 1 when a chunk row holds a `planned` seam cell at the tip.
- [ ] TF48: The chunk refusal names each planned row ID, in a two-row fixture, and the citation form.
- [ ] TF49: A refused chunk freeze leaves the review record bytes unchanged.
- [ ] TF50: A chunk row whose seam cell uses the citation form lets the freeze record.
- [ ] TF51: A chunk row whose seam cell reads `review-owned: the planned canary stays green` lets the freeze record.
- [ ] TF52: A `planned` row outside the chunk's rows lets the freeze record.
- [ ] TF53: A seam cell that holds `unplanned` is not projected as planned.
- [ ] TF54: `bench record completion` exits 1 when any plan row holds a `planned` seam cell.
- [ ] TF55: `bench record completion` records when every plan row is cited.
- [ ] TF56: `TestSeamCellTestNameOutsideTheCitationForm` stays green, so `bench coverage --check` keeps zero violations for a planned bare-name cell. Its red comes from a mutation probe, not from absence: delete the `plannedRe` early return in `strayTestNames`, and the test turns red.
- [ ] TF57: The projection returns the opt-in row IDs whose seam cell holds the marker.
- [ ] TF58: A non-opt-in map projects no row ID.
- [ ] TF59: A walk of every non-test Go file that `go list ./...` names finds the expression text `\bplanned\b` only in `internal/coverage`. The walk matches that expression text, not the word, so `ReceiptPhasePlanned = "planned"` in `internal/intent/ledger/ledger.go` leaves it green.
- [ ] TF67: A tip that cites the row lets `bench record chunk` record while the working copy still holds `planned`.
