# Learnings — usage journal

## 2026-08-18 — `bench commit` refused twice with `prospective authorization refused: inherited` [open]

**What happened.** Landing the A1 light-path ticket, `bench commit -m ... <path>` ran its
prospective gate and exited `error: prospective authorization refused: inherited` — twice,
on two separate small commits, each time on a tree whose ordinary `bench gate` was green.
The composed tree had a green verdict available, but as an *inherited* (partial-evidence)
one, and `internal/landing/landing.go` requires `authorization.Green` exactly. Running
`bench gate --fresh` and re-issuing the identical `bench commit` succeeded both times.
The second occurrence also printed `gate: red` above the refusal, which reads as a real
red and sent me to diagnose a failure that did not exist — the following `--fresh` run was
green with no diff between them.

**Right behavior.** Nothing was wrong with the tree, so nothing needed diagnosing. But the
message is unactionable as written: `inherited` names an internal verdict kind, not
something the operator did or should do, and it names no next action. The
`gate: red` line printed alongside it is worse than unhelpful — a refusal to authorize an
inherited verdict is not a red gate, and printing one costs a fresh gate run to disprove.

**Proposed rule change.** `internal/landing`: when authorization refuses on an inherited
verdict, say so in operator terms and name the fix — the composed tree has only partial
evidence, run `bench gate --fresh`. Do not emit `gate: red` for it; a partial verdict and
a failing one are different facts, which is the same distinction this session's A1 change
drew between an environment skip and a green. Alternatively `bench commit` could escalate
to a fresh prospective gate itself rather than refusing, since that is what the operator
does next in every case.
