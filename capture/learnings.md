# Learnings — usage journal

## 2026-08-11 — abandon has no sanctioned route past a dirty staged spec  [open]
Applying `bench spec build abandon` with an uncommitted staged-spec revision hit
a dead end from every side: abandon refuses a dirty working checkout; `git
stash` is guard-blocked as a working-tree set-aside; and committing the revision
satisfied cleanliness but changed the spec blob, so abandon then refused with
`staged spec no longer matches recorded subject`. Recovery took a revert commit,
the abandon, and a revert-of-the-revert — three full-gate cycles (~9 minutes)
for what should be one operation. A superseding spec revision is a normal reason
to abandon a run, so the pin defends nothing the apply doesn't destroy anyway.
Proposed rule change: exempt the staged spec path from abandon's clean-checkout
and spec-identity preconditions in `internal/specbuild/precondition.go` —
mirroring abandon's existing moved-tip exemption — so a dirty or committed-ahead
spec revision survives the abandon untouched instead of forcing the revert
dance. Identity, ownership, and evidence stay binding.

## 2026-08-11 — `bench commit --spec` is promotion semantics, not ticket attribution  [open]
Passed `--spec axi-spec-build-complete` while landing a repair-ticket planning
batch, intending only to associate the commit with the active build. The commit
oracle automatically changed the staged spec's status to `implemented`, making
the active run refuse further lifecycle operations because its recorded spec
blob no longer matched. The right behavior was a normal path-scoped `bench
commit` naming the ticket and learning paths, with no `--spec`; only `bench spec
build promote` owns implemented status. Proposed rule change: make `bench commit
--help` state that `--spec` marks the named spec implemented, or refuse `--spec`
when that slug has an active spec-build run.

## 2026-08-11 — checkpoint receipt ownership is the exact changed-path set, not the ticket fence  [open]
Built a coordinator checkpoint receipt with every path in the assignment's
ownership fence, including `internal/specbuild/refresh.go`, although that path
was unchanged. `bench spec build checkpoint` rejected the otherwise valid
receipt with only the generic invalid-receipt response. Reading
`validateReceipt` showed that receipt `ownership` must equal the sorted unique
paths actually changed between the assignment base and receipt tree; the
ticket fence is only the upper bound. The same validator also requires receipt
rows to equal the assignment rows exactly and binds the coordinator probe to
the assignment id and tree. Removing the unchanged path made the checkpoint
succeed. Proposed rule change: expose or scaffold the complete checkpoint
receipt shape through `bench`, including the derived changed-path set and rows,
so coordinators do not duplicate private validator knowledge or confuse the
authorization fence with observed ownership.

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
