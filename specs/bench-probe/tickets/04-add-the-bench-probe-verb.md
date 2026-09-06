# Add the bench probe verb

Blocked by: 01-make-headroom-in-the-registry-files.md, 02-expose-the-focused-run-outcome.md, 03-export-the-gate-execution-probe.md
Writes: internal/probe/command.go (new), internal/probe/probe.go (new), internal/probe/subject.go (new), internal/probe/probe_test.go (new), internal/probe/refusal_test.go (new), internal/probe/outcome_test.go (new), cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/help_inventory_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, internal/conformance/subcommand_routing_table_test.go (new), internal/tickets/registry_data.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: PB1, PB2, PB3, PB4, PB5, PB6, PB7, PB8, PB9, PB10, PB11, PB12, PB13, PB14, PB15, PB16, PB17, PB18, PB19, PB20, PB21, PB22, PB23, PB24, PB25, PB26, PB29, PB30, PB31, PB38, PB39, PB40, PB41, PB43, PB44, PB45, PB46

## What to build

Verify the premise first. Read the spec's implementation decisions in full. Read
`testreport.Prepare`, `testreport.Execute`, and `testreport.Outcome`, which ticket 02
adds. Read `gate.ExecutionInProgress`, which ticket 03 adds. Read `anchorsCommand` in
cmd/bench/anchors_command.go for the subject resolution and the root read.

Read `structure.Command` in internal/structure/structure.go for the grammar and root
pattern. Read `otelrecord.Dir` in internal/otelrecord/writer.go for the home and
pool-key composition. Read `toon.Table`, `toon.Errorf`, and `toon.RenderError` in
internal/toon/toon.go. Read the `test` and `outline` rows of `commandRegistry` in
cmd/bench/main.go and the routing table in
internal/conformance/subcommand_routing_table_test.go.

Build the verb end to end in `internal/probe`:

- `Command(args []string) (string, int)` parses the grammar through `usage.Parse`,
  resolves the root, and runs the probe. The spec's decisions fix the grammar, the
  refusal order, the error lines, the copy location, the atomic replacement, the
  verdict mapping, and the output blocks. Spell each error line as the coverage map
  spells it.
- Register `probe` in `commandRegistry` beside `test` with `axiExempt(axiReasonMutation)`
  and the help row the spec fixes. Add `"probe": routed("internal/probe")` to the
  routing table. Add the help golden row after the `bench test` row.

The preserved copy file is the restore source, and the bytes read at the start are
the oracle. Compare the read-back with the start bytes, never with the copy. The
`internal/probe` tests reach no test-only helper of `internal/testreport`; PB7 grades
the parse step in-package against `testreport.Prepare`.

Write the tests the coverage map names. Build one module fixture `probefixture` with
the `Clamp` function and its two tests, with the real `go` on `PATH`. Write one stub
`go` for the outcome rows that writes a marker file when it starts. On request, the
stub truncates the preserved copy or makes the subject's directory read-only. Run the
interrupt row as a child-process helper. Hold the gate lock from a child process
for PB20.

Record two runs over the integration worktree under Further notes in the spec: the
`--check subcommand-routing` run of PB8 and the exec form of PB11. Run them after the
worktree's own `dist/bench` is rebuilt with `bench worktree build`.

Self-probe: remove the byte comparison after the restore and show
`TestProbeReportsAReadBackMismatch` red.

## Acceptance

- [ ] The six verdict and cause rows PB1 to PB5 and PB38 print the exact rows and exit codes.
- [ ] The nine refusal rows PB18 to PB26 print the exact lines, leave the subject and the home untouched, and start no run child.
- [ ] The restore rows PB12 to PB17, PB40, PB41, and PB43 hold, with the copy kept only on `restore-failed`.
- [ ] The write-failure rows PB44 and PB45 and the no-record row PB46 hold.
- [ ] `bench help` prints the probe row, `bench probe --help` exits 0, and the three registry checks pass through `bench test --check`.
- [ ] The two recorded runs sit under Further notes with their stdout.
