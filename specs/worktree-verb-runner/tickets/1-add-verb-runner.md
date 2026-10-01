# Add the verb runner and its readers

Blocked by: none
Writes: internal/worktree/verb_runner_test.go (new), internal/worktree/verb_runner_check_test.go (new), internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR1, VR2, VR3, VR4, VR5, VR6, VR7, VR8, VR9, VR10, VR11, VR12, VR13, VR14, VR15, VR16, VR17, VR18, VR19, VR20, VR21, VR22, VR59

## What to build

Add `verb_runner_test.go`. It holds the verb runner, the typed verb keys, the verb call value, the verb result, the core readers, and the must-form readers. Declare one key for each verb that a worktree test runs today. The keys are `create`, `release`, `clean`, `reclaim`, `reauthorize`, `merge`, `reset`, `land`, `land-resume`, `resume-clean`, `list`, `path`, `show`, `build`, `exec`, `pool`, and `lease-file`. Each key names its verb entry, and a verb with a joins form also names that form.

The verb call value carries the root, the Bench home, the kit root, the clock, stdin, joins, and the arguments. The joins value is optional. Without a joins value, the runner calls the verb entry. With a joins value, the runner calls the joins form with that value unchanged. A key without a joins form refuses a joins value.

The function `checkVerbCall` returns an error for a call with a kit value or a clock value. The runner fails the test with that error. The kit message is exactly `verb runner: a kit value waits for the seam reduction spec`. The clock message is exactly `verb runner: a clock value waits for the seam reduction spec`.

The verb result carries the exit code, stdout, stderr, and the exec assignment. The `exec` key calls `ExecCommandResolving`, so the result carries its assignment.

The core rows reader, `readVerbRows`, decodes stdout through `axitest.DecodeDocument` and returns the rows of one table block, or an error. The core fingerprint reader, `readVerbFingerprint`, returns the agreed `fingerprint` cell of a decoded table, or else the `fingerprint=` cell of the one record line. It returns an error when no value or two different values exist. When the agreed value equals the package constant `unapplicableFingerprint`, it returns the no-fingerprint error, because an error plan carries that placeholder. Keep both reader names, because no test local in the package uses them.

The must-form methods `mustRows`, `mustFingerprint`, and `mustNoFingerprint` take a `testing.TB` and fail it on a core reader result that the method does not accept. Migrated tests use only the must forms. Each core reader and `checkVerbCall` is a top-level function whose last result is an `error`, so the ticket 11 census can find it. The must forms call only `Helper` and `Fatalf` on the `testing.TB`, and each returns right after its `Fatalf` call.

The must-form tests use a test recorder that embeds a nil `testing.TB`. The recorder overrides `Helper`, `Fatal`, `Fatalf`, and `FailNow`, and its `Fatalf` records the failure and returns.

Add `verb_runner_check_test.go` for the runner's own tests. Each new file stays at or below 400 lines. Raise `worktreeTestCount` by the tests this ticket adds. This ticket moves no existing test.

Contract for later tickets: the verb keys, the call value's fields, the verb result's fields, and the must-form methods stay as this ticket defines them. A later ticket adds a key only for a verb that this list lacks.

## Acceptance

- [ ] Each key with a usage grammar returns its own verb's usage refusal, with the expectation read from the usage constant.
- [ ] The `pool` key returns the pool path of its root argument.
- [ ] The `lease-file` and `resume-clean` keys each match a direct call to their own verb entry.
- [ ] The verb result equals a direct call's exit code and both streams for the same input.
- [ ] The core rows reader returns the rows of a real `list` table and an error for a table followed by a non-TOON line.
- [ ] The core fingerprint reader returns the cell of a real `clean --landed` plan and the record cell of a real `reset --to` plan.
- [ ] The core fingerprint reader returns an error for no fingerprint and for two different fingerprints, and it returns a numeric-looking cell exactly.
- [ ] The core fingerprint reader returns the no-fingerprint error for a real `clean` error plan whose rows carry the `none` placeholder.
- [ ] `mustRows` and `mustFingerprint` fail a test recorder on a core reader error.
- [ ] `mustNoFingerprint` accepts a real error plan with the `none` placeholder and fails a test recorder on stdout with a fingerprint.
- [ ] A joins value's stub runs, stdin reaches an exec child, and the exec result carries the resolved assignment.
- [ ] `checkVerbCall` returns each pinned message for a kit value and for a clock value.
- [ ] `TestPackageTestCountPin` and `TestSerialSetStaysBelowTheCeiling` pass.
