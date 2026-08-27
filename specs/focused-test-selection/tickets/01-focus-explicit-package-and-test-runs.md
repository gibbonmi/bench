# Focus explicit package and test runs

Blocked by: none

Writes: internal/testreport/command.go (new), internal/testreport/testreport.go, internal/testreport/testreport_test.go, internal/testreport/runbinary_test.go, internal/testreport/cancel_test.go, cmd/bench/main.go, cmd/bench/command_registry_test.go, specs/focused-test-selection/

## What to build

Extract the current command mechanics from the renderer and parse one typed
focused request. Preserve the positional package form, add `--package` and
`--run`, and keep one selected executable and cancellation path. Observe Go
test run events so a filtered invocation that matches zero tests refuses.
Advertise only the forms this ticket makes runnable.

## Acceptance checklist

- [x] F01 — default, positional, and `--package` selection preserve package-expression compatibility.
- [x] F03 — `--run` reaches Go as one unchanged argv value for default and explicit-package subjects.
- [x] F04 — zero matched tests refuses while matched skips and failures remain observed runs.
- [x] N01 — default, package, and run-filtered forms write no gate-owned record.

Delivered outcome: agents can run one package or one test pattern through the
supported focused renderer without reconstructing a Go command.
