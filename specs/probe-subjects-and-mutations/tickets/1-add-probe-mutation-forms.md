# Add the probe mutation forms and the closed form matrix

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/main.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/combined_test.go (new), internal/probe/command.go, internal/probe/delete_test.go (new), internal/probe/edit.go (new), internal/probe/edit_test.go (new), internal/probe/mutation_input.go, internal/probe/omit_file_test.go, internal/probe/probe.go, internal/probe/refusal_test.go, internal/probe/subject.go, internal/probe/unwrap_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: PS1, PS2, PS3, PS4, PS5, PS6, PS7, PS8, PS9, PS10, PS11, PS12, PS13, PS14, PS15, PS16, PS17, PS18, PS19, PS20, PS21, PS22, PS23, PS24, PS25, PS76

## What to build

An author gives `bench probe` a multiline swap, an insertion, one occurrence, a deletion, or a combined mutant.
Each form refuses a bad input before any write, and each form restores the subject.

The NRO chunk NR-C3 is on `main`.
It supplies the `probe_execution` table, `probe.ReadResult`, the `bench-probe-result-v1` token, and the canned runner in `internal/probe/probetest`.
Bind that landed shape, and refresh each line count before the charge.

Create `internal/probe/edit.go` as the one owner of the edit list and the form rules.
The mutation input resolves to an ordered list of edits over the start bytes of the subject.
Each edit is a byte range and its replacement bytes.
Move `mutate` into that owner, so every form shares the match, overlap, and empty rules.
`probe.go` keeps the sequence and calls the owner.
`internal/probe/unwrap_test.go:45` and `:49` call `mutate` directly and read its result, so they follow the moved owner.

Extract the current `--omit-file` reader in `mutation_input.go` into one bounded regular-file reader.
Keep its refusal set and its refusal text unchanged.
`--with-file` reads through that reader, and PS2 runs one refusal table over both flags in `omit_file_test.go`.

The forms follow the spec's mutation-input decisions:

- `--insert-before` adds the text immediately before the selected anchor match, and the mutation cell reads `insert`.
- `--occurrence <n>` selects one match for one `--swap`, `--omit`, `--unwrap`, or `--insert-before` edit.
- `--delete` removes the subject for the mutated run, and the mutation cell reads `delete`.
- The restore after `--delete` writes the preserved copy through the current atomic replace and proves the bytes and the mode.
- An omission whose result has zero bytes refuses before any write with the empties line that names `--delete`.
- Repeated `--swap` and `--omit` flags make one combined mutant, and its cell lists the swaps first.

`command.go` exports the usage line as `probe.Usage` and grades the closed form matrix that the spec lists.
An argv names exactly one mutation group, and each partner fault is a usage refusal with exit 2.
`cmd/bench/main.go` derives the probe help suffix from `probe.Usage`, as `testHelpSuffix` derives the `bench test` suffix.
`cmd/bench/main.go` is already over the 400-line structure budget, so keep the derivation inside the existing row and add no line.

`TestProbeUsageNamesOneSelectionAndOneMutation` loses its "two mutations" subtest, because that argv is now a valid combined mutant.
`TestProbeHelpSpellings` builds its expectation from `grammar.Help` at `refusal_test.go:307`, so it follows the new usage line.
The literal expected probe line at `cmd/bench/help_inventory_test.go:120` changes to the new usage line.
The binding registry binds the command registry files to `internal/probe` and `cmd/bench`, so they join this `Writes:` line.
Change a registry file only when the help projection requires it.
The canary unit `tests/canary/package-core-guard/unrouted-subcommand` pins `cmd/bench/main.go`, so it joins this line and keeps its meaning.

Ticket 2 receives the edit owner, the form matrix, and `probe.Usage`.
Ticket 3 adds `--expect` to the same grammar and adds the `first_failure` cell to every expected verdict row.

## Acceptance

- [ ] A `--with-file` swap puts the exact multiline bytes of the file into the mutant and reads `bit`.
- [ ] Each bad `--with-file` or `--omit-file` input refuses with its current reason before any write.
- [ ] `--insert-before` puts the `--with` or `--with-file` text immediately before the one anchor match.
- [ ] An anchor that matches twice refuses as ambiguous unless `--occurrence` selects one match.
- [ ] `--occurrence 2` changes only the second of two identical matches.
- [ ] An out-of-range occurrence refuses with the line that names the match count, and a zero or non-decimal value gives usage exit 2.
- [ ] `--delete` leaves the subject absent for the mutated run, and the subject then holds its start bytes and its start mode.
- [ ] A failed restore after `--delete` reads `restore-failed` and names the preserved copy.
- [ ] An omission that empties the subject refuses with the empties line before any write.
- [ ] Two `--swap` pairs and one `--omit` make one mutant whose bytes hold all three edits, and the cell reads `swap+swap+omit`.
- [ ] Each combined old string matches exactly once in the start bytes, and overlapping edits refuse with the overlap line before any write.
- [ ] Each argv that names two mutation groups, or that breaks a partner rule, prints the usage line with exit 2.
- [ ] A zero-match edit under `--occurrence` or in a combined mutant gives the substring-miss row with the form's cell.
- [ ] The help inventory's probe line names every mutation flag and comes from `probe.Usage`.
- [ ] `cmd/bench/main.go` does not grow, and each new file stays at 400 lines or fewer.
- [ ] Every existing probe assertion outside the retired "two mutations" subtest stays green.
