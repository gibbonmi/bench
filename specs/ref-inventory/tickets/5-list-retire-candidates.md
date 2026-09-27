# 5. List the superseded candidates and the unique count on spec retire

Blocked by: 2-route-status-to-the-plan.md, 4-discard-a-unique-ref-by-target.md
Writes: cmd/bench/main.go, cmd/bench/spec_retire_listing.go (new), cmd/bench/spec_retire_listing_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand, internal/spec/spec.go, internal/spec/spec_test.go, internal/spec/history.go, internal/worktree/clean_discard.go
Covers: RI48, RI49, RI50, RI51, RI52, RI69, RI70, RI71, RI77, RI79, RI80

## What to build

Chunk: RI-C2b.
The two blockers are shared-write edges: this ticket reads the counts function that ticket 2 exports and spells the command that ticket 4 fixes.

Compose the listing at the `cmd/bench` dispatch of `bench spec retire`, because `internal/worktree` already imports `internal/spec`.
The retire command in `internal/spec` returns as today.
`bench spec retire` takes `<spec.md | slug>` and derives the slug with an unexported function.
Export that slug derivation as one function that the retire command and the dispatch wrapper both call.

The wrapper appends the candidate lines and the count line only after a code-0 retire that printed a `next:` line.
It inserts them before that line, and it appends nothing after `--help` or after a refusal.
One line per active or cleanup-pending assignment whose label or request token contains the slug reads `superseded candidate: <assignment id> <label> — bench worktree clean --discard-branch --target <assignment id>`.
A complete record is not listed, and the calling worktree's own assignment is listed.

One count line reads `unique refs: <n> — bench worktree clean --discard-branch --unclaimed`, where `<n>` comes from the counts function of ticket 2.
The command text in that line comes from the exported plan command spelling of the worktree package, not from a second literal.
The candidate line's discard command comes from one exported function of the worktree package that takes the target selector.
The slug rename updates every caller, including the history command.
When the ledger read or the planner fails, the count line reads `unique refs: unavailable — <error>`.
The listing discards nothing and changes no exit code.

## Acceptance

- [ ] Retire of slug `s` with one active assignment labelled `s-build` prints one candidate line with its id and the exact `--target` command.
- [ ] Retire with the operand `specs/s/spec.md` prints the same candidate lines as the operand `s`.
- [ ] Retire of slug `s` with one active assignment whose request token is `s-run` and whose label is `other` prints one candidate line.
- [ ] Retire of slug `s` with one complete assignment labelled `s-build` prints no candidate line.
- [ ] Retire with two unique unrecorded refs prints `unique refs: 2 — bench worktree clean --discard-branch --unclaimed`.
- [ ] Retire with no matching assignment prints `unique refs: 0` and no candidate line.
- [ ] The candidate lines and the count line print after the `retired:` lines and before the `next:` line.
- [ ] `bench spec retire --help` prints no candidate line and no count line.
- [ ] Retire with an unreadable ledger prints `unique refs: unavailable — <error>` and exits 0.
- [ ] Retire with a candidate exits 0 and every branch ref survives.
