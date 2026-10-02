# 3. Append a verification result with `bench record verification`

Blocked by: 2-record-chunk-entry.md
Writes: internal/reviewrecord/parse.go, internal/reviewrecord/write.go (new), internal/reviewrecord/recordcmd/ (new), cmd/bench/help_inventory_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RE36, RE42, RE43, RE44, RE45, RE46, RE47, RE48, RE49, RE50, RE51, RE52, RE53, RE54, RE55, RE56, RE57, RE58, RE59, RE60, RE61, RE62, RE63, RE64, RE65, RE66, RE67, RE68, RE69, RE70, RE71, RE104, RE107, RE108

## What to build

Chunk: RE-C3.

Add the `verification` form through the transaction of ticket 2. Put its record change in `internal/reviewrecord/write.go`, so `verifier` and `findChunk` stay private. `recordcmd` parses the flags, reads the excerpt, and prints the output.

`--chunk <id>` appends to that chunk's verification list, and `--final` appends to `completion.verification`. The source commit is `--source`, or the chunk entry's tip when `--chunk` has no `--source`. A chunk result whose source digest differs from the chunk entry's refuses and names `bench record chunk`. `--final` requires `--source`.

The plan at the source tree supplies the requirement, its command, and its probe mutation. The plan's `verifier` rule supplies the role. `reviewrecord.SourceDigest` of the source tree supplies the source digest. The outcome is `pass` for exit code 0, and `fail` otherwise. Before the render, refuse an `--id` that the record already holds, and name that ID.

Add the control-character rule for every single-line flag that the spec lists. Add the excerpt reader to `recordcmd`. It reads `--excerpt` through `bounds.ClassifyNoFollow` and refuses each state other than `parsed` with the `toon.RecordError` line. The verb embeds the exact bytes, and `reviewrecord.Digest` of those bytes is the digest. A planned probe takes the three probe flags, and its native result is the entry's native result. Ticket 4 reuses the excerpt reader.

Print the `verification[1]{list,chunk,id,requirement,role,outcome,source_digest,excerpt_digest}` row of the spec. Add the verification row to `HelpRows` and to the help golden.

## Acceptance

- [ ] A chunk result lands in the chunk's verification list, and a `--final` result lands in the completion list.
- [ ] A chunk result with no `--source` takes the chunk tip's digest, and a stale `--source` exits 1.
- [ ] A `--final` result takes the digest of its `--source`, and the command is the planned command.
- [ ] An unplanned requirement exits 1 and names the planned IDs.
- [ ] The roles are `author-verification` and `integration-verification`, and an undispatched ticket exits 1.
- [ ] Exit code 0 writes `pass`, and exit code 3 writes `fail`.
- [ ] A planned probe takes its mutation from the plan, and a probe mismatch exits 1.
- [ ] A result for an unrecorded chunk, and `--final` with no record, name `bench record chunk`.
- [ ] The excerpt bytes and their digest are exact, and each refused excerpt state exits 1.
- [ ] A control character in `--performer`, `--model`, `--effort`, or `--id` exits 1 and names the flag.
- [ ] A recorded `--id` exits 1 and names the ID.
- [ ] A parser refusal of a `../` reference leaves no temporary file.
- [ ] The committed record passes `reviewrecord.Check` for chunk `1` with fixture review results.
- [ ] Each grammar error exits 2, and `bench help` prints the verification row.
