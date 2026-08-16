# Plan the landed set under one fingerprint

Blocked by: count-and-advertise-landed-assignments.md
Writes: internal/usage/worktree.go, internal/worktree/worktree.go, internal/worktree/clean_landed.go, internal/worktree/clean_landed_test.go, internal/worktree/clean_landed_hostile_test.go, cmd/bench/command_registry_test.go, .bench/BENCH.md

## What to build

`usage.WorktreeClean` becomes `bench worktree clean [--discard-ignored] [--discard-branch]
[--full] (<path> | --landed) [--apply <fingerprint>]` (`--apply` takes exactly 64
lowercase hex characters); `bench worktree --help` names the full string and the
`.bench/BENCH.md` inventory sentence follows it. Bare `bench worktree clean --landed`
selects every ledger assignment the shared classifier admits, repository-wide whatever
request, label, or session created it, sorted by assignment id; shape-classifies each
path and retains any non-checkout shape (FIFO, socket, device, dangling symlink,
absent, decayed) as `uncertain` with the shape named and no git invoked; runs the
existing per-path explicit planner on the rest with the given options; retains any row
whose per-path action would preserve as `dirty` with its per-path remedy; and prints the
`worktree_cleanup` table with one row per selected assignment, every row's `fingerprint`
equal to one set fingerprint — a digest over, per row in id order, (assignment id,
target, planned action, HEAD OID, tracked state, ignored count, lease probe state), plus
the three option flags and a version tag — followed by a help block: the `--apply <fp>`
invocation when any row removes, and one quoted `bench worktree clean <path>` per
retained row (paths with spaces/globs pasteable, control bytes replaced by the pointer
form). Zero selected rows print the empty table, no help, exit 0; `--landed` with a path
operand in either order, a malformed `--apply` value, and `--apply` on an empty set are
the existing invocation error (exit 2). The bare plan mutates nothing. Demo: run the bare
verb over a mixed pool.

## Acceptance

- [ ] `(covers LR7)` Three landed assignments under three requests/labels (two clean, one dirty) plan as three id-ordered rows sharing one fingerprint, actions `remove`, `remove`, `retain` (`dirty`), no recovery ref, help with the apply invocation and one `clean <path>` for the dirty row, nothing removed.
- [ ] `(covers LR11)` Zero landed rows: empty table, no help, exit 0 on repeated runs; `--apply <64-hex>` on that pool exits 2.
- [ ] `(covers LR12)` `--landed <path>` and `<path> --landed` exit 2 naming the new usage; short, long, non-hex, and uppercase `--apply` values exit 2 with `--landed`; `bench worktree --help` contains the full new grammar.
- [ ] `(covers LR17)` A path with a space and `*` plans `remove`; its dirty variant's remedy line and the apply invocation render quoted and pasteable; apply (next ticket) removes the clean one.
- [ ] `(covers LR18)` A path with ESC plans `retain` (`uncertain`), pointer target, no byte in any help line, exit 0, present after apply; a path with a tab renders escaped as one row.
- [ ] `(covers LR19)` FIFO, dangling symlink, unix socket (`internal/bounds` `requireSocket` pattern, capability-guarded), and a `/dev/null`-resolving path each plan `retain` (`uncertain`, shape named); the command completes without opening the special file; apply removes none.
- [ ] Selector partition: a landed row with an unparseable lease is selected; a live-lease row, a non-true-proof row, and a `cleanup-pending`/`recovered`/`complete` record are not.
