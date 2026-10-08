# Grade the completion plan against the ticket check floor

Blocked by: 01-name-the-root-conformance-check.md, 02-derive-the-ticket-check-floor.md
Writes: internal/preflight/check_floor.go (new), internal/preflight/check_floor_command_test.go (new), internal/preflight/check_floor_test.go (new), internal/preflight/command_bootstrap_test.go, internal/preflight/command_build_test.go, internal/preflight/command_review_test.go, internal/preflight/decision.go, internal/preflight/decision_test.go, internal/preflight/gather.go, internal/preflight/plan.go, internal/preflight/preflighttest/fixture.go, internal/preflight/ticket_rows.go (new)
Covers: TF6, TF16, TF17, TF18, TF23, TF24, TF25, TF26, TF27, TF28, TF29, TF30, TF31, TF32, TF33, TF34, TF35, TF36, TF37, TF38, TF39, TF40, TF41, TF42, TF43, TF44, TF45, TF46, TF66, TF68, TF69, TF70

## What to build

`Decide` adds the row `check-floor` directly after `completion-plan`, in build mode and in review mode. The row shares the ticket-directory gate of the other ticket rows. It answers not-applicable when the completion plan is invalid. The spec section "The check-floor row" owns the floor composition, the command parse, and the satisfaction rules, and this ticket implements them unchanged.

The gatherer keeps the parsed completion plan beside its digest, and it composes each ticket's floor from three parts:

- the element `root-conformance`, from the gate predicate of ticket 01 with the kit value `gate.KitValue()`;
- each package pattern and the loader error, from the floor function of ticket 02;
- the element `./cmd/bench`, from the tickets predicate of ticket 02 or the embed flag.

Each plan command goes through `shellcommand.Parse` and then `testreport.Prepare`. The `--changed` accessor of ticket 02 tells a `--changed` command from a `./...` package command. A version 1 plan grades the union of a chunk's ticket floors against that chunk's requirements. A delegated version 2 plan grades each ticket's floor against the chunk requirements that name that ticket.

The red detail lists `<ticket>: <element> (add <command>)` for each unsatisfied element, in ticket order. A loader failure reds the row with the loader message. The existing verdict refusal then refuses the build charge and the review charge with `preflight required: check-floor`. The commit lane and its class table stay unchanged.

`decision.go` holds 418 lines and `gather.go` holds 413 lines, and neither file has a structure grant. Neither file may end this ticket longer than at the base. Put the new row and its gather step in `check_floor.go`. Move existing ticket-row checks from `decision.go` into `ticket_rows.go` to make room for the new `Facts` fields. In `gather.go`, wrap the existing `withCompletionPlan` return rather than add lines. The row-list edits in `decision_test.go` (601 lines) and `command_review_test.go` (480 lines) keep each line count.

The system suite drives `bench preflight review` over fixture repositories that hold a Go module. Run `bench test --check system` before the ticket commit, because a loader failure in such a fixture reds `check-floor`.

At the ticket commit, rewrite the `planned` seam cell of each covered row into the citation form.

## Acceptance

- [ ] TF6: A kit root that declares the entry test gives each ticket the `root-conformance` element, also under an ambient `BENCH_KIT` that names that kit.
- [ ] TF16: A command-package entry adds the element `./cmd/bench`.
- [ ] TF17: An entry that selects a package with an embed pattern adds `./cmd/bench`.
- [ ] TF18: An entry outside every command package and every embed package adds no `./cmd/bench`.
- [ ] TF23: An omitted floor element reds `check-floor` in three fixtures: a package, `./cmd/bench` alone, and the second ticket's package in a two-ticket version 1 chunk.
- [ ] TF24: A plan whose chunk satisfies every floor element makes `check-floor` green.
- [ ] TF25: The red detail reads `<ticket>: <element> (add <command>)` with the exact command token.
- [ ] TF26: `bench test --package ./internal/x` satisfies the element `./internal/x`.
- [ ] TF27: `bench test --package ./internal/...` satisfies the element `./internal/x`.
- [ ] TF28: `bench test --package ./...` satisfies every package element.
- [ ] TF29: `bench test --package ./internal/git/...` leaves `./internal/git2` unsatisfied.
- [ ] TF30: A package command with a quoted `--run` pattern leaves its package unsatisfied.
- [ ] TF31: A `bench test --changed` command leaves every element unsatisfied.
- [ ] TF32: A command with `&&` before `bench test` satisfies nothing.
- [ ] TF33: The commands `go test` and `BENCH_KIT=x bench test` satisfy nothing.
- [ ] TF34: A `bench test` command that the grammar refuses satisfies nothing.
- [ ] TF35: `bench test --check root-conformance` satisfies the `root-conformance` element.
- [ ] TF36: `bench test --check docs-currency-workflow` leaves the `root-conformance` element unsatisfied.
- [ ] TF37: `bench test --package ./internal/conformance` leaves the `root-conformance` element unsatisfied.
- [ ] TF38: A requirement in another chunk leaves this chunk's element unsatisfied.
- [ ] TF39: A delegated requirement that names another ticket leaves this ticket's element unsatisfied.
- [ ] TF40: The default `bench preflight build` form prints a red `check-floor` row for an omitted package.
- [ ] TF41: The review form prints a red `check-floor` row for an omitted package.
- [ ] TF42: A red `check-floor` row refuses the build charge with `preflight required: check-floor`.
- [ ] TF43: An invalid completion plan makes `check-floor` not-applicable.
- [ ] TF44: Build mode with no tickets directory makes `check-floor` not-applicable.
- [ ] TF45: A loader failure makes `check-floor` red with the loader message.
- [ ] TF46: `TestBenchkitLaneTable` stays green, so the kit lane table keeps its check list.
- [ ] TF66: A linked root with `go.mod` and no declared entry test gives no `root-conformance` element.
- [ ] TF68: `bench test --package ./internal/git2` leaves the element `./internal/git` unsatisfied.
- [ ] TF69: `bench preflight build <slug> --plan-only` prints a red `check-floor` row for an omitted package.
- [ ] TF70: The rendered `check-floor` red holds no raw control byte inside its row.
