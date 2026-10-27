# Bench probe narrows each run to its expect test

Blocked by: none
Writes: internal/probe/probe.go, internal/probe/selection.go, internal/probe/set.go, internal/probe/command.go, internal/probe/narrow_test.go, internal/probe/refusal_test.go, internal/probe/outcome_test.go, internal/testreport/held_runs.go, internal/testreport/held_runs_test.go, docs/known-issues.md, CHANGELOG.md
Covers: none

## What to build

A `bench probe` call without `--run` runs the whole package for its baseline and for each mutant.
One run of a large package takes about 70 seconds, so a probe set takes most of the time of a ticket author.
The `docs/known-issues.md` row for this issue tells the author to give each probe `--run` with the killing test.

The probe now narrows a package selection that has no `--run` to the test that `--expect` names.
The narrowed pattern anchors each slash-separated part of the name, so `TestFoo` does not select `TestFooBar`.
The pattern for `TestFoo/case` is `^TestFoo$/^case$`, so the run selects the parent test and only that subtest.
The narrowed pattern also selects the subtests of the named test.

A single probe uses the narrowed pattern for its baseline run and for its mutated run.
In a mutant set, each entry with an `--expect` test runs its mutant over its own narrowed pattern.
When each matched entry has an `--expect` test, the set baseline runs the union of the entry patterns.
When a matched entry has no `--expect` test, the set baseline runs the whole selection, and that entry runs the whole selection too.
The `run` cell of the selection row shows the pattern of the baseline run.
All runs of a set use the one held run binary, so a set still builds its run binary at most once.

An explicit `--run` governs as before, and a `--check` selection does not change.
A probe without `--expect` does not change.
The verdict words do not change.

Under the narrowed run, only the named test and its subtests run, so another test cannot cause `expect-missed`.
`expect-missed` still occurs when `--expect` names a subtest and the parent test fails outside that subtest.

A misspelled `--expect` test matches no test after this change.
A single probe then reports `invalid` with the cause `baseline-no-test-run`.
A set entry then reports `invalid` with the cause `no-test-run`, or with `baseline-no-test-run` when no entry test matches.
Before this change, the whole package ran, and the probe reported `expect-missed`.

`HeldRuns.RunRequest` runs a given request over the held run binary, and `HeldRuns.Run` calls it with the held request.
The probe help adds one note that states the narrowing.
The change removes the known-issues row, because its retire condition then holds.
A `CHANGELOG.md` entry under `[Unreleased]` records the change.

## Acceptance

- [ ] A single probe with `--expect T` and no `--run` runs only T for the baseline and for the mutant.
- [ ] Each `--mutants` entry with an `--expect` test and no `--run` runs only its test, and the baseline runs the union of the entry tests.
- [ ] An explicit `--run` overrides the narrowing, and a probe without `--expect` runs the whole selection as before.
- [ ] A subtest `--expect` reports `expect-missed` when the parent test fails outside that subtest.
- [ ] A misspelled `--expect` reports `invalid` with the cause `baseline-no-test-run`.
- [ ] `HeldRuns.RunRequest` runs only its own selection over the one held run binary.
- [ ] The known-issues row is removed, `bench test --package ./internal/probe` is green, and the lane passes at `bench commit`.
