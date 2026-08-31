# Read preflight rows from the ticket parser

Blocked by: build-ticket-file-parser.md
Writes: internal/preflight/gather.go, internal/preflight/decision.go, internal/preflight/gather_test.go, internal/preflight/decision_test.go, internal/preflight/command_build_test.go, internal/preflight/command_review_test.go
Covers: TG6, TG12, TG13, TG21, TG22, TG35, TG37, TG38

## What to build

Preflight reads parsed tickets instead of bare tokens. The gatherer parses
every ticket file through `internal/tickets`. `Facts` carries the parsed
tickets. The bare-token scrape deletes, so `tokenRe` and `Facts.TicketTokens`
leave the tree. Decide stays pure: the gatherer collects the `Writes:`
existence bits into `Facts`.

Only `.md` files parse as tickets, and a non-`.md` file is ignored. A basename at two depths yields a
duplicate-identity diagnostic.

Three verdict rows join the table: `tickets-parse`, `blockers-resolve`, and
`writes-resolve`. They render after `paths-authorized` and before `rows-owned`.
One later ticket inserts `fixture-closure`, `registry-closure`, and
`kit-pin` in that same span.

`writes-resolve` reds an entry that names no tree path and carries no `(new)`
marker. An entry marked `(new)` stays green when the tree has no such path.
`rows-owned` and `rows-membership` read the parsed `Covers:` tokens alone, so a
row ID in prose is not evidence.

The existing postures hold. Build mode without a tickets directory keeps the
not-applicable rows. A special file, a dangling symlink, or an unreadable entry
stays a bootstrap failure. A control byte in a diagnostic path is refused
before the verdict table renders.

## Acceptance

- [ ] TG6 — a special file under tickets/ yields a named refusal.
- [ ] TG12 — an absent `Writes:` path without `(new)` reds `writes-resolve` naming it.
- [ ] TG13 — an absent `Writes:` path marked `(new)` stays green.
- [ ] TG21 — a row in no ticket's `Covers:` reds `rows-owned` naming the row.
- [ ] TG22 — a phantom `Covers:` token reds `rows-membership` naming the token.
- [ ] The three new rows each render with their own detail.
- [ ] TG37 — a non-`.md` file under tickets/ is ignored by the grammar.
- [ ] TG38 — one basename at two depths yields a duplicate-identity diagnostic.
- [ ] TG35 — a control byte in a diagnostic path is refused before the render.
- [ ] `tokenRe` and `Facts.TicketTokens` no longer exist in the tree.

## Delegate charge

You work in the Bench repo on the `ticket-grammar` spec. Read
`specs/ticket-grammar/spec.md` first. Then read `internal/preflight/gather.go`,
`internal/preflight/decision.go`, and `internal/tickets/tickets.go` in full.
Read `internal/preflight/decision_test.go` for the table-test shape.

Parse every enumerated `.md` ticket file through `internal/tickets`. Ignore a
non-`.md` file. Report a duplicate-identity diagnostic for one basename at two
depths. Gather the existence bit for each `Writes:` entry in `gather.go`, and
keep `Decide` free of I/O.

Pass the sibling basename set and the spec tag from `Facts`. Carry the parsed tickets on
`Facts`. Delete `tokenRe` and `Facts.TicketTokens`.

Keep the lstat-first classification at every depth. Keep the bootstrap failure
for a special file, a dangling symlink, and an unreadable entry. Keep the
absent-directory posture of build mode.

Add `tickets-parse`, `blockers-resolve`, and `writes-resolve` to `Decide`.
Render them after `paths-authorized` and before `rows-owned`. Give each row its
own detail, so one red never hides another.

Red `writes-resolve` for an entry that names no tree path and carries no
`(new)` marker. Name the entry in the detail. Keep a `(new)` entry green when
the tree has no such path.

Rewrite `rowsOwnedCheck` and `rowsMembershipCheck` to read the parsed `Covers:`
tokens. Keep both delivered detail strings unchanged. Keep the foreign-tag
token ignored at this seam.

Add `TestWritesResolveNamesAbsentPath`, `TestWritesResolveAcceptsNewMarker`,
`TestCoversDrivesRowOwnership`, and `TestProseRowIDIsNotOwnership` in
`internal/preflight`. Add `TestNonMarkdownFileIsIgnored` and `TestDuplicateBasenameAcrossDepths` in
`internal/preflight`. Add a command test that asserts one verdict row per
grammar red. Add a command test for the control-byte refusal.

Run only `bench worktree exec ft174-ticket-grammar -- go test ./internal/preflight/`.
Do not commit. Do not edit the spec.
