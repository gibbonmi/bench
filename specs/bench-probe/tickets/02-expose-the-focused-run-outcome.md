# Expose the focused run as a typed outcome

Blocked by: none
Writes: internal/testreport/command.go, internal/testreport/testreport.go, internal/testreport/outcome.go (new), internal/testreport/outcome_test.go (new)
Covers: PB27, PB28

## What to build

Verify the premise first. Read `Command`, `parseFocusedRequest`, `runFocusedRequest`,
`runNamedCheck`, and `runGoTest` in internal/testreport/command.go. Read `report`,
`decode`, `markNonzeroFailures`, and `render` in internal/testreport/testreport.go.
Read `writeCheckGo` in internal/testreport/check_test.go and `focusedTestModule` in
internal/testreport/testreport_test.go.

Split `Command` into two exported steps that keep its output byte-identical:

- `Prepare(root string, args []string) (Request, string, int)` parses the selection
  with the grammar and returns the usage line and code verbatim. `Request` wraps the
  parsed request with unexported fields.
- `Execute(root string, request Request) (Outcome, string, int)` runs the request and
  returns the outcome beside the rendered report and the exit code.
- `Command` is `Prepare` then `Execute`, and it drops the outcome.

`Outcome` carries `Kind` and `FailedTests`. The kinds are `passed`, `failed`,
`build-failed`, `no-test-run`, `refused`, and `interrupted`. Derive each kind from
the report:

- a failing test row gives `failed` with the count
- a run with at least one test and no failure gives `passed`
- a package status `fail` with no failing test gives `build-failed`
- a run that ran no test, or the `run pattern matched no tests` refusal, gives `no-test-run`
- the `go test interrupted` refusal gives `interrupted`
- every other refusal before or after the child gives `refused`

Before the refactor, run the base commit's `Command` over the four canned event sets in
the package form and the run form. Record each exact output and exit as a golden in
`TestCommandKeepsItsBaseOutput`. Then refactor, and keep those goldens green.

Drive the kinds through a stub `go` on `PATH` that emits one canned `-json` event set
per kind. Drive the `refused` kind through `installTestSelectionFactory` with a `Build`
that returns an error. Assert `Command` equals `Prepare` then `Execute` for the four
sets. Keep every existing test in the package unchanged.

Self-probe: make `Execute` answer `failed` for a build failure and show
`TestExecuteClassifiesTheOutcome` red.

## Acceptance

- [ ] `TestExecuteClassifiesTheOutcome` answers the six kinds and the failed count, with `refused` from a failed selection.
- [ ] `TestCommandKeepsItsBaseOutput` holds the base commit's exact output and exit for each canned set.
- [ ] `TestCommandIsThePrepareExecuteProjection` shows `Command` equal to `Prepare` then `Execute` for the four sets.
- [ ] `go test ./internal/testreport` passes with the existing assertions unchanged.
