# Repair the preflight facts and the comment record

Blocked by: none
Writes: internal/freshness/freshness_verify.go, internal/worktree/build.go, internal/freshness/freshness_verify_test.go, internal/preflight/gather.go, internal/preflight/decision_test.go, internal/adopt/doctor.go, internal/adopt/doctor_rows.go, internal/adopt/link.go, internal/adopt/broker_test.go, internal/brokermanifest, internal/systemtest/owner_stale_seal_test.go, cmd/bench/build_subject_mode_test.go, internal/preflight/source_tip_test.go, specs/binary-freshness/spec.md
Covers: BF26, BF27

## What to build

This ticket repairs review findings C3, F4, S4 to S9, F2, and F7.

`binarySealFacts` in `internal/preflight/gather.go` calls `os.Stat`, which
follows a symbolic link. A dangling `dist/bench` link therefore reports the
row not applicable, while `evalBinarySealRow` in `internal/adopt` uses
`os.Lstat` and reds the same state. The `--full` build preflight is the one
that fails open. Change the call to `os.Lstat` and add the decision row for
the non-regular and dangling state.

`binarySealFacts` also has no test at all. Both preflight rows set the facts
by hand, so the path spelling and the verifier call are ungraded. Add the row
that grades the gatherer against a real root.

Then fold the record fixes:

- Delete the `kitSourceCheckout` alias in `internal/adopt/doctor.go` and
  rename its five call sites.
- Remove the retained red record from the two test comments, and the ordinal
  "refusal three" that no code names.
- Correct the stale nine-row sentence in `internal/adopt/broker_test.go` and
  the five-row sentence in `internal/preflight/source_tip_test.go`.
- Rename `staleySealedCopy`.
- Unexport `brokermanifest.Fields`; no reader outside its own file exists.
- Correct the spec's Implementation decisions line that says
  `freshness-publish` gains two arguments, and the line that says
  `commands --brief` becomes a root-taking handler.

## Acceptance

- [ ] A dangling `dist/bench` symlink reds `binary-seal` in build mode.
- [ ] A new decision test grades `binarySealFacts` against a real root.
- [ ] No `kitSourceCheckout` alias remains.
- [ ] The two Implementation decisions lines match the tree.
- [ ] Self-probe: restore `os.Stat`, and report which row reds.
