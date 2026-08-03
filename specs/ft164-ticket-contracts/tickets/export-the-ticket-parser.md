# Export the ticket parser as one cross-package entry point

Blocked by: none
Ownership fence: `internal/specbuild`
Assumptions: the resolution lives in `internal/specbuild/assign.go` as the unexported `ticket` struct and the `resolveTicket` function; its three internal callers are `assign.go` plus `checkpoint.go` plus `integrate.go`; the `expandRows` and `unique` and `listValue` helpers stay unexported. Re-derive from the tree at pickup.

## What to build

The ticket resolution that today only `internal/specbuild` can call becomes
callable from another package, so the FT164 example-agreement check grades the
skill's Good example with the same parse that assignment runs — one source for
the accepted ticket shape instead of a second reader that drifts.

This is an export, not a behavior change. Rename the `ticket` struct to `Ticket`
and `resolveTicket` to `ParseTicket`, and move the three internal callers to the
exported names, rather than adding a wrapper around a duplicated struct: two
declarations of the same record are the knowledge duplication the code standard
names as a defect. The grammar, the field parsing, the silent `internal/<pkg>`
fence fallback, the range expansion, and every refusal stay byte-for-byte the
behavior they are now, and the existing lifecycle tests keep passing unedited.

A doc comment on the exported function states that the conformance
example-agreement check is its cross-package consumer, so a later reader knows
why the symbol is exported and does not narrow it back.

## Acceptance

- [ ] [EP1] the exported entry parses a well-formed ticket file to its title, digest, acceptance rows, fence, and assumptions.
- [ ] [EP2] the exported entry refuses a ticket with zero acceptance rows, and one with neither a fence line nor any `internal/<pkg>` mention, with today's messages.
- [ ] [EP3] the exported entry refuses an absent path, a non-regular file, an empty argument, an absolute path, and a path escaping its spec.
- [ ] [EP4] `R`-prefixed range expansion, non-`R` identifiers left literal, and the comma-split backtick-trimmed list values behave exactly as before the export.
- [ ] [EP5] every internal caller in assign, checkpoint, and integrate resolves tickets through the exported entry, and no second ticket record type exists in the package.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EP1 | drop assumptions collection from the exported path | `TestParseTicketMatchesAssignPath` | write one ticket fixture carrying every field, call the exported parse, compare each field against authored literals |
| EP2 | return a zero-row or fenceless ticket instead of refusing | `TestAssignRejectsHostileTicketsBeforeWorktreeCreation` | `go test ./internal/specbuild -run '^TestAssignRejectsHostileTicketsBeforeWorktreeCreation$'` over its `malformed.md` and `no-fence.md` fixtures, expecting both refusals |
| EP3 | accept a `../` argument after cleaning | `TestAssignRejectsHostileTicketsBeforeWorktreeCreation` | run the same focused test over its missing, `../spec.md`, directory, fifo, and symlink arguments, expecting a refusal for each and zero worktree owner calls |
| EP4 | expand any prefixed identifier as a range, not only `R` | `TestParseTicketRangeExpansionUnchanged` | parse a fixture carrying `[R10-R12]` and `[EP1]`, assert three expanded ids and one literal |
| EP5 | leave one caller on a retained unexported duplicate | `review` | read the package for a second ticket record declaration — a retained duplicate still compiles and still passes every test, so review is the honest owner |
