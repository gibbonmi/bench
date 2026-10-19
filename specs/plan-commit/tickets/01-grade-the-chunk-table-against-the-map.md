# Grade the chunk table against the coverage map in build preflight

Blocked by: none
Writes: internal/coverage/rows.go, internal/coverage/coverage.go, internal/coverage/chunk_table.go (new), internal/coverage/chunk_table_test.go (new), internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/decision.go, internal/preflight/ticket_rows.go, internal/preflight/chunk_rows.go (new), internal/preflight/chunk_rows_test.go (new), internal/preflight/charge_pack.go, internal/preflight/command.go, internal/preflight/decision_test.go, internal/preflight/charge_cells_test.go, internal/preflight/source_tip_test.go, internal/preflight/command_build_test.go, internal/preflight/command_review_test.go, internal/preflight/command_bootstrap_test.go, internal/preflight/preflighttest/fixture.go, internal/preflight/preflighttest/chunk_table.go (new), internal/preflight/preflighttest/reviewfiles.go, internal/preflight/preflighttest/reviewcases.go, internal/preflight/preflighttest/testdata/review-documents.toon, internal/preflight/preflighttest/testdata/review-empty.toon, internal/preflight/preflighttest/testdata/review-shapes-renames.toon, internal/preflight/preflighttest/testdata/review-shapes-no-renames.toon, internal/preflight/staleness/completeness_test.go, internal/preflight/ownership/ownershiptest/repo.go, internal/preflight/ownership/premises_test.go, internal/preflight/ownership/premises_fixture_test.go, internal/preflight/ownership/charge_fixture_test.go, internal/reviewrecord/recordtest/fixture.go, internal/systemtest/charge_evidence_test.go, internal/systemtest/owner_landing_fixture_test.go, cmd/bench/preflight_version_test.go, internal/gate/refusal_route_test.go, .bench/commitment.json
Covers: PC1, PC2, PC3, PC4, PC5, PC6, PC7, PC8, PC9, PC10, PC11, PC12, PC13, PC14, PC15, PC16, PC17, PC18, PC37, PC38

## What to build

Build preflight gains the `chunk-rows` row. It compares the spec's `## Implementation chunks` table with the coverage map in both directions. The spec's decision "The `chunk-rows` row" states the exact predicate, the token rule, the range rule, and the two pinned details. Obey these exact predicates:

- The reader reads the first pipe table under the `## Implementation chunks` heading, before the next `##` or `###` heading.
- It finds the acceptance column by the lowercased header cell `acceptance rows`, at any column position.
- It splits a cell on commas, trims spaces and backticks from each token, and skips an empty token.
- A token `<A>-<B>` or `<A>–<B>` is a range when both endpoints match `tickets.RowIDPattern` with one tag. It names its endpoints and each declared row with that tag whose number lies between them.
- A range with an undeclared endpoint, two tags, or a first number above its last number is a fault that names the token. Any other token that is not a row ID is a fault that names the token.
- An absent table gives the detail `no Implementation chunks table`. An absent column gives the detail `chunk table has no acceptance rows column`.

Extract the pipe-cell splitter inside `processMapLine` (`internal/coverage/rows.go`) into one function. The coverage-map reader and the new chunk-table reader in `internal/coverage/chunk_table.go` both call it. Write no second Markdown splitter. `coverage.ParseSpec` keeps its signature and its own read. The gatherer passes the spec bytes that it already holds (`content`, `internal/preflight/gather.go:37`) to the chunk-table reader, and `Facts` carries the listed row IDs and the faults. The reader never fails `Gather`: an absent table, an absent column, and each token fault reach `Facts` and show as a red row. The landing path runs `Gather` and only `paths-authorized` (`internal/preflight/gather.go:209-224`), so this rule keeps every landing unchanged.

`Decide` appends `ticketRow(f, "chunk-rows", chunkRowsCheck)` directly after `rows-membership`. Put the check in `internal/preflight/chunk_rows.go` beside the pattern of `rowsMembershipCheck`. Add no charge code: the build and review charges already refuse on `verdict.Red`.

Make every spec fixture that reaches a graded verdict render a chunk table. Put the table renderer in `internal/preflight/preflighttest/chunk_table.go`, because `fixture.go` has 4 lines of headroom. The shared `specBody` renders its table from the rows that it declares. Place the table after the coverage map and before the completion-plan fence: `charge_fixture_test.go:38` and `premises_test.go:296-299` pin the map rows at spec lines 11 and 12, and `completion_plan_test.go:13` and `completeness_test.go:124` trim the plan fence as the last block. The table text must not contain `bench test --package ./...`, because `check_floor_command_test.go:24-25` replaces the first such text in `SpecBody`. A test that rewrites a map row ID, such as `command_build_test.go:55` and `staleness/completeness_test.go:31-32`, rewrites the table cell too.

The canonical review spec in `reviewfiles.go:45` renders `SpecBodyNoInventory`, so the four stored `.toon` review baselines change with it. Regenerate each baseline once through its owner, and compare the old and new files in the commit message. The readers `internal/diff/patches_test.go` and `internal/preflight/evidencecmd/evidence_file_reconstruction_test.go` must stay green.

`internal/preflight/decision_test.go` (557 lines) and `internal/preflight/command_review_test.go` (480 lines) are over the 400-line structure ceiling with no grant, so the net line growth of each is zero. `internal/preflight/gather.go` has 8 lines of headroom and `internal/reviewrecord/recordtest/fixture.go` has 3. Put each new test in `chunk_rows_test.go` and extend `baseFacts` and the row-order lists in place. Each inline spec body calls the shared table renderer; paste no table. The plan writer of `internal/reviewrecord/recordtest/fixture.go` stays unchanged, and only the inline bodies of its two system callers gain a table. Update the row-order lists and counts in `decision_test.go` (`TestDecideAllGreen`), `charge_cells_test.go` (`TestDecideChargeRows`), and `source_tip_test.go` (`bareRowCount`). `internal/systemtest/charge_evidence_test.go` has 1 line of headroom, so edit its inline spec at line 61 in place.

This ticket writes the system-tagged files `internal/systemtest/charge_evidence_test.go` and `internal/systemtest/owner_landing_fixture_test.go`, so run them through `bench test --check system`, which sets `BENCH_KIT` to the kit root.

The spec-stage close commits this spec's commitment binding inside the range that the plan-only preflight grades, so this ticket's `Writes:` line names `.bench/commitment.json`. The build itself does not edit that file.

PC18 is review-owned: record one plan-only preflight verdict for each ticketed staged spec under `specs/`, at the ticket tip, and report each `chunk-rows` cell.

Contract for ticket 03: this ticket names the row `chunk-rows`, and the completion-plan reference names it by that exact token.

## Acceptance

- [ ] `bench preflight build <slug>` prints a `chunk-rows` row directly after `rows-membership`.
- [ ] The coverage-map reader and the chunk-table reader call one cell-split function.
- [ ] A declared row that no acceptance cell lists reds `chunk-rows`, and the detail names that row.
- [ ] An acceptance-cell ID that the map does not declare reds `chunk-rows`, and the detail names that ID.
- [ ] `PF2-PF4` and `PF2–PF4` list PF2, PF3, and PF4 when the map declares all three, and they list exactly PF2 and PF4 when the map omits PF3.
- [ ] A range with an undeclared endpoint, two tags, or reversed bounds, and a token that is neither a row ID nor a range, each give a fault that names the token.
- [ ] The cell `` `PF1`, PF2, `` lists PF1 and PF2 with no fault, and a row that two chunk rows list stays green.
- [ ] A header `Acceptance rows` at the fifth position after a `blocked by` column gives the acceptance cells.
- [ ] An absent table and an absent acceptance column red the row with their two pinned details.
- [ ] Build mode without a tickets folder reads `chunk-rows` not-applicable with no detail.
- [ ] A red `chunk-rows` row refuses the build charge with `preflight required: chunk-rows`.
- [ ] The shared fixture builds a green `chunk-rows` row, and every PC-C1 completion-plan command passes.
- [ ] Each ticketed staged spec reads `chunk-rows` green in its plan-only preflight at the ticket tip.
