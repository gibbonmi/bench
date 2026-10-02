# 2. Write a chunk entry with `bench record chunk`

Blocked by: 1-render-record-fence.md
Writes: internal/reviewrecord/write.go (new), internal/reviewrecord/recordcmd/ (new), internal/reviewrecord/recordtest/fixture.go, cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: RE11, RE12, RE13, RE14, RE15, RE16, RE17, RE18, RE19, RE20, RE21, RE22, RE23, RE24, RE25, RE26, RE27, RE28, RE29, RE30, RE31, RE32, RE33, RE34, RE35, RE37, RE38, RE39, RE40, RE41, RE105, RE106

## What to build

Chunk: RE-C2.

Add the write transaction to `internal/reviewrecord/write.go`. It reads the record through the existing reader, parses it, applies one change, renders with `Render`, and parses the rendered bytes again. Only then does it write a temporary file in `reviews/` and rename it over the record. Only a caller that allows creation may start from an absent file, and the function then creates `reviews/` when it is absent. A refusal, a failed temporary write included, leaves the record bytes unchanged and leaves no temporary file.

Put the chunk change in `write.go` too, so the private plan rules stay in package `reviewrecord`. Keep `write.go` at or below 400 lines.

Add the package `internal/reviewrecord/recordcmd` with `Command(root, args)`. It dispatches the forms, prints the help lines for `--help`, `-h`, and `help`, and applies the refusal order of the spec. `HelpRows` returns the suffix and the description of each implemented form. Add the `chunk` form with the derived fields, the creation rule, the update rule, and the `chunk[1]{id,action,base,tip,source_digest,plan_digest,rows}` output of the spec.

Add `recordtest.NewLinked`. It prepares a fixture in a linked worktree of a new repository and returns the fixture and the primary checkout root.

Add the `record` registry row at help order 40, with the mutation AXI exemption, the tree scope, and the bounded response. The row calls an adapter in `cmd/bench/command_registry.go`, and its help rows project `recordcmd.HelpRows`. Move the inline adapter of the `assessment` row into `cmd/bench/command_registry.go`, so `cmd/bench/main.go` ends at or below 446 lines. Record `record` as routed through `internal/reviewrecord/recordcmd` in the subcommand-routing table. Add the chunk row to the help golden, and add `TestRecordRouteAnswersItsUsage` to `cmd/bench/help_inventory_test.go`.

Tickets 3, 4, and 5 consume the transaction, `Command`, `HelpRows`, the refusal order, and `NewLinked`. `write.go` and `recordcmd/` carry `(new)` because the build preflight reads a tree without them.

## Acceptance

- [ ] The chunk entry holds full commit IDs, the source digest without the record file, the `ReadPlan` digest, and the planned rows.
- [ ] A version 2 plan creates a version 2 record, and a version 1 plan exits 1 and names `version 2`.
- [ ] A new chunk appends, and a recorded chunk updates and keeps its results.
- [ ] Two equal runs leave byte-identical files.
- [ ] An unplanned chunk and an unknown revision exit 1 and leave the record unchanged.
- [ ] The output row names `created`, `added`, or `updated`.
- [ ] Prose around the fence survives, and `git status` lists only the record path.
- [ ] A dangling link, a live link, a FIFO, an empty file, an unterminated fence, and a duplicate JSON key each exit 1 and change nothing.
- [ ] A read-only `reviews/` exits 1 and changes nothing.
- [ ] The primary checkout exits 1 before any revision read.
- [ ] An empty root exits 1 with the not-in-repository line.
- [ ] Each grammar error exits 2 with a `usage: bench record` line.
- [ ] `--help`, `-h`, `help`, and the real dispatcher print the chunk usage line at exit 0.
- [ ] `bench help` prints the chunk row of the spec, and `cmd/bench/main.go` holds at most 446 lines.
