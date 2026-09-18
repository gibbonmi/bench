# Close anchor registry files into ticket Writes

Blocked by: none
Writes: internal/anchors/references.go (new), internal/anchors/references_test.go (new), internal/preflight/closure.go, internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/proposal.go, internal/preflight/anchor_closure_test.go (new), internal/preflight/decision_test.go, internal/preflight/command_build_test.go, internal/preflight/proposal_test.go, internal/preflight/charge_test.go, internal/preflight/command_review_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: SC1, SC2, SC3, SC4, SC5, SC6, SC7, SC8, SC9, SC10, SC11, SC12, SC13

## What to build

A ticket slicer runs `bench preflight build <slug>` and sees an `anchor-closure` row after `registry-closure`. The row reds a ticket that writes an anchored guidance path and omits an anchor registry file that names it. The red lists `<ticket>: <entry> is anchored by <file>` for each missing file.

The anchors package gets one function that scans the top-level `.go` files of the anchor registry directory, test files included. It returns each string literal with the sorted files that hold it. It classifies each entry without following links. It refuses a link, a special file, an unreadable file, or a tokenizer error, and it names the path. An absent directory is an empty result.

The gatherer calls that scan once. For each `Writes:` entry, it records the files of every literal equal to the entry path or under it at a `/` boundary. The three closure kinds then share one derivation and one message builder. This ticket creates the proposal-tolerated checks list with the three closure checks, and `--propose-writes` reads only that list. The proposal lists a missing anchor file with source `anchor <entry>`.

Before this ticket adds lines, it moves the per-entry `Writes:` probe out of `gather.go` and the owned-path helpers out of `decision.go` into the closure file. The command registry files are named only because registry closure binds the anchors package; this ticket does not edit them.

The new row changes the rows that three test files assert. This ticket adds `anchor-closure` to the legacy charge baseline, the review command's row count, and the two `Decide` row-order lists. The edits in `command_review_test.go` and `decision_test.go` are net-neutral in line count.

## Acceptance

- [ ] `Decide` reds `anchor-closure` with the exact detail for an omitted registry file, and it is green when the ticket names the file or its directory.
- [ ] A directory `Writes:` entry takes the closure of an anchored file under it.
- [ ] The scan returns test files, constant and positional literals, and raw literals, and it returns an empty result for an absent directory.
- [ ] The scan refuses a symbolic link and a FIFO by path, without a blocking read.
- [ ] `bench preflight build` over a seeded repository with its own anchors file reds the row.
- [ ] `--propose-writes` lists the missing anchor file and is not refused by the `anchor-closure` red.
- [ ] In build mode with no tickets directory, the row renders not-applicable after `registry-closure`.
- [ ] The structure lane passes at the ticket commit.
