# State the one-expression rule of bench test --package in the help

Blocked by: none
Writes: internal/testreport/selection_facts.go, internal/testreport/command.go, internal/testreport/selection_facts_test.go, internal/probe/refusal_test.go, CHANGELOG.md
Covers: none

## What to build

Two authors in drain d-0bca6e72fedd passed two packages to `bench test --package`
in one space-separated value. The run failed with `directory not found`, and the
help did not state the rule.

The run passes the `--package` value to Go as one typed argument. A test already
requires that a value with a space stays one argument, so the run does not split
a list. The help states this rule and the two routes: one call per package, or a
parent `./...` pattern.

One constant holds the sentence. The `bench test` help prints it in a notes block
before the check inventory. The probe notes print the same constant in place of
their old package sentence.

## Acceptance

- [x] `bench test --help` prints the package-expression note before the check inventory.
- [x] `bench probe --help` prints the same note as its first fact.
- [x] A value with a space still reaches Go as one argument.
- [x] Removal of the note from the `bench test` help turns the new test red.

## Verification

The `internal/testreport` and `internal/probe` packages pass. A `bench probe` run
removed the note from the `bench test` help, and
`TestHelpStatesThePackageExpressionRule` was red. The probe restored the file.
