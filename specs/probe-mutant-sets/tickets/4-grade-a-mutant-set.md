# Grade a mutant set against one probe baseline

Blocked by: 1-repeat-the-probe-package.md, 2-hold-one-run-binary-across-runs.md, 3-share-the-mutation-flag-grammar.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/main.go, cmd/bench/record_operations_test.go, cmd/bench/test_command.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/baseline.go, internal/probe/combined_test.go, internal/probe/command.go, internal/probe/edit.go, internal/probe/edit_test.go, internal/probe/mutation_input.go, internal/probe/outcome_test.go, internal/probe/probe.go, internal/probe/probe_test.go, internal/probe/probetest/fixture.go, internal/probe/refusal_test.go, internal/probe/replacement_test.go, internal/probe/set.go (new), internal/probe/set_input.go (new), internal/probe/set_input_test.go (new), internal/probe/set_test.go (new), internal/probe/subject.go, internal/probe/subject_check_fixture_test.go, internal/probe/unwrap_test.go, internal/probe/verdict.go, internal/probe/verdict_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: PM1, PM2, PM3, PM4, PM5, PM6, PM7, PM8, PM9, PM10, PM11, PM12, PM13, PM14, PM15, PM16, PM17, PM18, PM19, PM20, PM21, PM22, PM23, PM24, PM25, PM26, PM30, PM32, PM33, PM35, PM36, PM37, PM38, PM39, PM40, PM41, PM43, PM44, PM51, PM67, PM69, PM71, PM72, PM77, PM78

## What to build

`bench probe <file> --mutants <path> <selection>` grades a list of independent mutants of one subject in one call.
The verb runs one probe baseline, and then it preserves, writes, runs, and restores each mutant in turn.
Every run goes through the held runs of ticket 2, so a mutant outside the run binary's closure reuses the baseline's build.
The answer carries the schema token `bench-probe-set-v1`, and the single form stays unchanged.

The spec's Grammar, Mutants file, Sequence, and Envelope sections state every line, cell, and exit.
Build the entry grammar from the shared flag list of ticket 3, and run only the mutation half for each entry.
Call ticket 3's cross-half rule for each entry, so an entry `--expect` under `--check prose` refuses at exit 2.
Read the file through `readMutationFile`, and decode it through `jsonfile.DecodeDocument`.
Refuse a top-level `null` and a `null` entry as malformed, and hold the 64-entry cap in one named constant.

One cancel context from `subprocess.NotifyCancel` spans the set from the baseline to the last restore.
Before each mutant write, examine that context and grade the gate lock again.
Apply the stop cells `set-stopped`, `gate-held`, `gate-unavailable`, `mutation-failed`, and `interrupted` as the spec states.
A restore failure stops the set, keeps the copy, and exits 2.

Keep the set tables in their own registry, apart from `resultTables`, so `FieldSource` never sees a set field.
`internal/probe/probe.go` has little line headroom, so put the loop in `set.go` and the input in `set_input.go`.
`internal/probe/subject.go` keeps its path, because `internal/canonicalpath/canonicalpath_test.go:37` pins that path.

`probetest` gains one installer that answers the first `go test` start as the baseline and each later start in order.
`InstallStarts` keeps its odd-and-even rule, and both installers share one script prelude.
Change the probe description in `cmd/bench/main.go` to state one or several mutants of one file and a restore after each run.
The command registry files join this `Writes:` line through the binding registry.

## Acceptance

- [ ] A two-entry set over `clamp.go` prints two `mutants` rows that both read `bit` at exit 0, after exactly three `go test` starts.
- [ ] During mutant 2's run, the subject holds mutant 2's bytes and none of mutant 1's edit.
- [ ] Each entry refusal names its entry number and answers before any Go start.
- [ ] An unmatched entry reads `substring-miss` while the other entries run, and an all-miss set starts no run and takes no lock.
- [ ] A restore failure after mutant 1 gives row 2 `set-stopped`, a `preserved` table, and exit 2.
- [ ] A signal during mutant 1's fresh build gives row 1 `interrupted`, row 2 `set-stopped`, one Go start in all, and the start bytes.
- [ ] With an owned run binary, a two-entry set over `clamp.go` records one build, and its `builds` cell reads 1.
- [ ] The single form still builds twice, and its second build reads the mutant bytes.
- [ ] The help inventory's probe line names `--mutants <path>`, and its description states one or several mutants.
