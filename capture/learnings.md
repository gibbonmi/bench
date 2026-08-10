# Learnings — usage journal

## 2026-08-10 — craft-delegate's worktree-cutting instruction doesn't carve out spec-build assign  [open]
Cut a manual `bench worktree create --request <id> --label <work-item>` ahead of
`bench spec build assign`, following craft-delegate's generic "the coordinator
cuts the worktrees" instruction. `assign` then provisioned its own worktree
under a different assignment id, leaving the manual one stray; caught via
`bench worktree list` and released it before proceeding. The right behavior
was reading `bench spec build assign`'s own contract first — spec-build
assignments self-provision. Proposed rule change: craft-delegate (or
`bench-implement-spec`'s "Route the venue" section) should say explicitly
that `bench spec build assign` provisions its own assignment worktree, so the
coordinator does not cut one ahead of it.

## 2026-08-10 — bench spec build assign --ticket's resolution rule isn't in its usage text  [open]
Guessed the `--ticket` argument shape twice — a bare ticket filename, then the
full `specs/<slug>/tickets/<file>.md` path — before either was accepted; both
failed with the same generic "must name one regular ticket file" error. Only
reading `ParseTicket` in `internal/specbuild/assign.go` revealed the argument
must be the ticket's basename relative to the spec's own `tickets/`
directory. The right behavior would have been not needing the source read.
Proposed rule change: `bench spec build assign`'s usage/error text should
name the exact resolution rule, matching this project's own CLI design
standard (`bench-craft-cli`) for self-explanatory structured errors.

## 2026-08-10 — no CLI affordance to read an assignment's current tree hash for a checkpoint receipt  [open]
`bench spec build checkpoint --evidence`'s receipt schema requires the exact
tree hash of the assignment worktree's current content (tracked + untracked,
computed via a throwaway-index `read-tree HEAD` + `add -A` + `write-tree`),
but no `bench` subcommand exposes it. Assembled it by hand in shell,
reimplementing `internal/git.TreeHash`'s algorithm externally — it worked,
but is exactly the kind of duplicated derivation `.bench/BENCH.md`'s
"one source per fact" standard warns against. Proposed rule change: add a
`bench` affordance (a checkpoint-receipt scaffold command, or a bare
tree-hash query for a given assignment) so assembling a receipt never
requires reimplementing `TreeHash` outside the binary.

## 2026-08-10 — declared a delegate's iteration cap and never checked usage against it on return  [open]
Declared "sonnet/medium/~6 iterations" for a write-delegate implementing the
derive-nested-grammar-membership ticket. On return (97 tool uses, ~25 min,
215k tokens per the completion notification), never compared that usage
against the declared cap or reported on it — the figures were already in
hand and simply went unchecked. This is a personal discipline lapse, not a
missing kit rule: craft-line's declared-cap discipline expects the
coordinator to notice and report an exhausted or exceeded cap, not just
state one at the start. Proposed rule change: none at the kit level;
treating "check usage against the declared cap on delegate return" as a
standing step alongside the other verify-the-done-claim checks going
forward.
