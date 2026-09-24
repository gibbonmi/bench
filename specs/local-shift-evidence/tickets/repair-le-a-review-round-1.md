# Repair the LE-A review round 1 findings

Blocked by: 3-rotate-and-retain-the-record.md
Writes: internal/otelrecord/writer.go, internal/otelrecord/reader.go, internal/otelrecord/encode_test.go, internal/otelrecord/reader_test.go, internal/otelrecord/processor_test.go, internal/otelrecord/retention_test.go, cmd/bench/main.go, tests/canary/package-core-guard/unrouted-subcommand, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, specs/local-shift-evidence/spec.md, specs/local-shift-evidence/tickets/3-rotate-and-retain-the-record.md
Covers: LE94

## What to build

Chunk: LE-A.

Close the accepted LE-A review findings. The reviewer decided that the listing of sealed names is the existence check, so remove the free-name step from the writer. Restate the LE94 reason as a rotation that ignores the present names. Add a Won't handle line for a rotation between a reader's list and its live open.

Make `ReadSelected` name the sealed segment of a problem and count lines within each segment. Run the LE4 record write in a helper process, so the row can fail in a full package run. Replace each planned LE-A seam with its test citation. Remove the duplicated encode steps, the second sealed-name parser, and the history in the test comments.

## Acceptance

- [ ] A rotation that ignores the present names turns `TestARotationSkipsAPlantedSequenceName` red.
- [ ] `TestBeginWritesNoEnvironmentResource` turns red in a full package run when the encoder writes the SDK resource.
- [ ] A malformed line in a sealed segment reads as that segment's name and its line 1.
- [ ] Each LE-A row of the coverage map cites its landed test.
