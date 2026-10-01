# Add the verb runner and move the reset family onto it

Blocked by: none
Writes: internal/worktree/verb_runner_test.go (new), internal/worktree/verb_fixture_test.go (new), internal/worktree/parallel_census_test.go, internal/worktree/reset_apply_test.go, internal/worktree/reset_fingerprint_test.go, internal/worktree/reset_plan_test.go, internal/worktree/reset_refusal_test.go, internal/worktree/reset_repair_test.go, internal/worktree/reset_restore_refusal_test.go, internal/worktree/reset_restore_test.go
Covers: VR1, VR2, VR3, VR4, VR5, VR6, VR7, VR8, VR9, VR10, VR11, VR12, VR13, VR14, VR15, VR16, VR17

## What to build

Add `verb_runner_test.go`. It holds the verb runner, the verb call value, the verb result, and the runner's own tests. Declare one typed verb key for each verb that a worktree test runs today. The keys are `create`, `release`, `clean`, `reclaim`, `reauthorize`, `merge`, `reset`, `land`, `land-resume`, `resume-clean`, `list`, `path`, `show`, `build`, `exec`, `pool`, and `lease-file`. Each key names its verb entry, and a verb with a joins form also names that form.

The verb call value carries the root, the Bench home, the kit root, the clock, stdin, joins, and the arguments. The joins value is optional. Without a joins value, the runner calls the verb entry. With a joins value, the runner calls the joins form with that value unchanged. A key without a joins form refuses a joins value. A call that carries a kit value or a clock value fails, and the message names the seam reduction spec.

The verb result carries the exit code, stdout, stderr, and the exec assignment. The `exec` key calls `ExecCommandResolving`, so the result carries its assignment. The rows reader decodes stdout through `axitest.DecodeDocument` and returns the rows of one table block, or an error. The fingerprint reader returns the agreed `fingerprint` cell of a decoded table, or else the `fingerprint=` cell of the one record line. It returns an error when no value or two different values exist.

Add `verb_fixture_test.go` with the named value that `restoreFixture` returns. Move every reset-family test onto the runner. Delete `runReset`, `runResetWith`, `resetFingerprint`, and `restoreFingerprint`. Keep every test name and every assertion. Raise `worktreeTestCount` by the runner tests this ticket adds. Record the tuple scan program in the review record for VR17.

Contract for later tickets: the verb keys, the call value's fields, the verb result's fields, and both readers stay as this ticket defines them. A later ticket adds a key only for a verb that this list lacks.

## Acceptance

- [ ] Each verb key returns its own verb's usage refusal for a usage-refusing argument list, with the expectation read from the usage constant.
- [ ] The verb result equals a direct call's exit code and both streams for the same input.
- [ ] The rows reader returns the rows of a real `list` table and an error for a table followed by a non-TOON line.
- [ ] The fingerprint reader returns the cell of a real `clean --landed` plan and the record cell of a real `reset --to` plan.
- [ ] The fingerprint reader returns an error for no fingerprint and for two different fingerprints, and it returns a numeric-looking cell exactly.
- [ ] A joins value's stub runs, stdin reaches an exec child, and the exec result carries the resolved assignment.
- [ ] A kit value and a clock value each fail with a message that names the seam reduction spec.
- [ ] No reset-family wrapper or extractor exists, and `restoreFixture` returns one named value.
- [ ] `TestPackageTestCountPin` and `TestSerialSetStaysBelowTheCeiling` pass, and no over-budget file grows.
