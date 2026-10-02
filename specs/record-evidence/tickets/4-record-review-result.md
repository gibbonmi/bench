# 4. Append a review result with `bench record review`

Blocked by: 3-record-verification-result.md
Writes: internal/reviewrecord/write.go (new), internal/reviewrecord/recordcmd/ (new), cmd/bench/help_inventory_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RE72, RE73, RE74, RE75, RE76, RE77, RE78, RE79, RE80, RE81, RE82, RE83

## What to build

Chunk: RE-C3.

Add the `review` form through the transaction of ticket 2 and the excerpt reader of ticket 3. Put its record change in `internal/reviewrecord/write.go`, so `findChunk` stays private. Apply the `--id` check and the control-character rule of ticket 3.

The chunk entry supplies the base, the tip, and the source digest. The role is `independent-review`, and the state is `completed`. The outcome is `pass` with no `--finding`, and `fail` otherwise. The finding IDs keep argv order. `supersedes` holds the ID of the last result of the same axis in that chunk, or no ID.

Print the `review[1]{chunk,id,axis,outcome,supersedes,source_digest,excerpt_digest}` row of the spec. Add the review row to `HelpRows` and to the help golden.

## Acceptance

- [ ] The entry copies the chunk's base, tip, and source digest.
- [ ] The first Standards result supersedes nothing, and a later Standards result supersedes the earlier Standards ID.
- [ ] No finding writes `pass`, and two findings write `fail` with both IDs in order.
- [ ] The entry has the role `independent-review` and lands in the review list.
- [ ] A result for an unrecorded chunk exits 1 and names `bench record chunk`.
- [ ] A version 1 self-review exits 1 with the parser message and changes nothing.
- [ ] A record written only by the verb passes `reviewrecord.Check` for chunk `1`.
- [ ] Each grammar error exits 2, and `bench help` prints the review row.
