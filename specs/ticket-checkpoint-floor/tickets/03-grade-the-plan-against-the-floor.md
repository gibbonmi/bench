# Grade the completion plan against the ticket check floor

Blocked by: 01-name-the-root-conformance-check.md, 02-derive-the-ticket-check-floor.md
Writes: internal/preflight/check_floor.go (new), internal/preflight/check_floor_command_test.go (new), internal/preflight/check_floor_test.go (new), internal/preflight/command_bootstrap_test.go, internal/preflight/command_build_test.go, internal/preflight/command_review_test.go, internal/preflight/decision.go, internal/preflight/decision_test.go, internal/preflight/gather.go, internal/preflight/plan.go, internal/preflight/preflighttest/fixture.go, internal/preflight/ticket_rows.go (new), internal/testreport/check_test.go, internal/testreport/command.go, internal/testreport/environment.go, internal/testreport/named_check.go, internal/testreport/root_conformance_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/preflight_version_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/preflight/source_tip_test.go, internal/preflight/preflighttest/testdata/review-documents.toon, internal/preflight/preflighttest/testdata/review-empty.toon, internal/preflight/preflighttest/testdata/review-shapes-no-renames.toon, internal/preflight/preflighttest/testdata/review-shapes-renames.toon, internal/testreport/selection_facts.go, internal/testreport/selection_facts_test.go
Covers: TF6, TF16, TF17, TF18, TF23, TF24, TF25, TF26, TF27, TF28, TF29, TF30, TF31, TF32, TF33, TF34, TF35, TF36, TF37, TF38, TF39, TF40, TF41, TF42, TF43, TF44, TF45, TF46, TF66, TF68, TF69, TF70, TF72, TF73, TF74

## What to build

`Decide` adds the row `check-floor` directly after `completion-plan`, in build mode and in review mode. The row shares the ticket-directory gate of the other ticket rows. It answers not-applicable when the completion plan is invalid. The spec section "The check-floor row" owns the floor composition, the command parse, and the satisfaction rules, and this ticket implements them unchanged.

The gatherer keeps the parsed completion plan beside its digest, and it composes each ticket's floor from three parts:

- the element `root-conformance`, when `gate.RootConformanceRuns(root, gate.KitRoot(root))` from ticket 01 answers true;
- each pattern of `Floor.Packages`, and the loader error, from `testreport.WritesFloor(ctx, root, ticket.Writes)` of ticket 02;
- the element `./cmd/bench`, when `tickets.CommandPackage` answers true for an entry path or for the directory of a floor package, or `Floor.Embed` is true.

The kit resolver `gate.KitRoot` is the one that ticket 01 makes `testBenchSource` call, so the preflight and the named check read one kit. Each plan command goes through `shellcommand.Parse` and then `testreport.Prepare`. The accessor `Request.Changed()` of ticket 02 tells a `--changed` command from a `./...` package command. A version 1 plan grades the union of a chunk's ticket floors against that chunk's requirements. A delegated version 2 plan grades each ticket's floor against the chunk requirements that name that ticket.

The red detail lists `<ticket>: <element> (add <command>)` for each unsatisfied element, in ticket order. A loader failure reds the row with the loader message. The existing verdict refusal then refuses the build charge and the review charge with `preflight required: check-floor`. The commit lane and its class table stay unchanged.

A `check-floor` red also refuses `--propose-writes`, on purpose. The proposal tolerates only the closure rows and `fence-writes`, so `completion-plan` and `rows-owned` reds refuse it the same way. Keep `check-floor` out of `proposalToleratedChecks`.

`decision.go` holds 418 lines and `gather.go` holds 413 lines, and neither file has a structure grant. Neither file may end this ticket longer than at the base. Put the new row in `check_floor.go`. `withCompletionPlan` lives in `plan.go`, not in `gather.go`, so put the floor gather step in `plan.go` or `check_floor.go`. Expect no edit in `gather.go`.

Move `kitPinCheck`, `coversTokens`, `rowsOwnedCheck`, and `rowsMembershipCheck` from `decision.go` into `ticket_rows.go`. That move frees about 65 lines, which covers the new `Decide` line and the new `Facts` fields. If the fields need more room, also move `ticketsParseCheck`, `blockersResolveCheck`, and `writesResolveCheck`.

The structure lane reds a file over 400 lines that grows. `command_bootstrap_test.go` holds 396 lines, and `preflighttest/fixture.go` holds 389 lines. Keep each of them at 400 lines or fewer, and put new fixture code into `check_floor_command_test.go`. The row-list edits in `decision_test.go` (601 lines) and `command_review_test.go` (480 lines) keep each line count.

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
- [ ] TF27: `bench test --package ./internal/...` satisfies the element `./internal/x`, and `bench test --package ./internal/x/...` satisfies the element `./internal/x` itself.
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
- [ ] TF46: `TestBenchkitLaneTable` stays green, so the kit lane table keeps its check list. Its red comes from a mutation probe, not from absence: insert one check before `structure` in `BenchkitLane`, and the test turns red.
- [ ] TF66: A linked root with `go.mod` and no declared entry test gives no `root-conformance` element.
- [ ] TF68: `bench test --package ./internal/git2` leaves the element `./internal/git` unsatisfied.
- [ ] TF69: `bench preflight build <slug> --plan-only` prints a red `check-floor` row for an omitted package.
- [ ] TF70: The rendered `check-floor` red holds no raw control byte inside its row.
- [ ] TF72: `bench test --checks` and `bench test --check root-conformance --fixtures` satisfy nothing.
- [ ] TF73: A ticket file that no plan chunk names makes `check-floor` red for each element of its floor.
- [ ] TF74: A directory entry above a bound command package adds the element `./cmd/bench`.
