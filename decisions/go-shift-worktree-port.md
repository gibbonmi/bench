# go-shift-worktree-port — slice 7 of the Go rewrite

Child of `decisions/go-rewrite.md` (#6, slice order): the worktree pool and the
`bench shift` gated loop. Bootstrap evidence: `bin/bench-worktree.sh` (103 lines —
lease claim/reclaim, acquire scan-or-mint, owner-only release, interactive
subshell; pool-path and lease-path addressing already Go via `internal/worktree`)
and the shift half of `bin/bench.sh` (~180 lines — gated loop, touched-path
staging, refactor phase, `run_gate` resolution + verdict recording; helpers
`structure --since`, `tree-hash` already Go). The regression net is black-box and
strong across three fragments: `gate-runtime-shift-contracts.sh` (green commit +
`benchBase`, staging byproduct/pre-dirt/spaces-globs, red rollback, early-done,
refactor trigger/no-trigger/no-op, SIGINT cleanup, adapter contract),
`gate-runtime-contracts.sh` (lease hardening, reuse/release cleaning, interactive
subshell via fake `$SHELL`, `BENCH_GATE` cwd), plus the gate-cache assertions.

## #1: Does gate running port in this slice?

Type: Grill

### Answer
**Yes.** The resolution chain (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect)
and verdict recording move into the Go core with the loop, their primary caller.
The dispatcher's `run_gate` becomes a one-glance adapter (the `worktree_pool`
precedent), so `bench gate` and the stop hook's `<wrapper> gate` path keep their
shape. Rejected: the Go loop shelling back out to bash `run_gate` — a Go→bash→gate
chain is two live gate-resolution implementations, which the slice-6 map's
watch-out forbids. Gate *content* (the fragments) stays shell until slice 8.

## #2: Does the auto-detect chain survive the port?

Type: Grill

### Answer
**Port as-is, plus a new resolution-order contract** (`.bench/gate.sh` beats
`$BENCH_GATE` beats auto-detect, with one detection case) — the chain currently
has no assertion at all. The strangler's premise is behavior carries unchanged;
bundling a product-behavior removal into a port diff muddies both. Dropping
auto-detect later remains open as its own small decision.

## #3: Is the parent map's contract-backfill precondition satisfied?

Type: Grill

### Answer
**Yes — satisfied by existing coverage; no backfill spec.** The parent flagged
shift/worktree coverage as thin; probing shows the lifecycle is asserted
black-box end to end (see bootstrap evidence above). The only gap worth closing
is #2's resolution-order contract, which rides in this slice. Rejected: a
separate backfill spec re-deriving coverage that already bites.

## #4: Where does the loop's agent-steering prose live?

Type: Grill

### Answer
**Go string constants.** `iteration_prompt` and `refactor_prompt` have always
lived inside the executable (heredocs in `bin/bench.sh`), never as
reviewer-tunable files; parent #2 drew the port line at "executable logic ports,
markdown content stays text," and slice 5 already put the stop hook's
agent-facing messaging in Go. Rejected: `.bench/prompts/*` kit files — a new
content surface plus the binary↔asset skew class the slice-6 map rejected
`go:embed` over, in reverse.

## Handoff

1. **Module boundaries.** `internal/worktree` grows the lifecycle it already
   addresses (lease claim/reclaim, acquire scan-or-mint, owner-only release);
   the loop and the gate runner land in new `internal/` package(s) — the split
   and naming are the spec's call, following the `gitguard` precedent.
   `cmd/bench` gains `shift` and `worktree` subcommands and a gate-runner
   subcommand `bench gate` routes through. `bin/bench-worktree.sh` is deleted;
   `bin/bench.sh` drops the shift/gate/worktree bodies and keeps only routing
   plus the one-glance `run_gate` adapter. Still shell after this slice:
   the dispatcher's routing, postinstall, the generated pre-push, hook entry
   shims (slice 6 separately deletes the link/init/doctor sources).
2. **Contracts.** All observable behavior carries — the three fragments run
   unchanged and are the port-parity net. The loop's stdout lines are contract
   surface (the fragments grep literals: `1 committed iteration(s)`,
   `■ shift done:`, `red gate`, `refactor phase`, `resume or clean up`);
   `branch.<name>.benchBase` config; lease file format `<pid> <utc-time>` and
   owner semantics; exit non-zero + full cleanup on interrupt; adapter invoked
   with the prompt as its single positional argument under `BENCH_SHIFT=1`,
   preflight rejecting missing/non-executable adapters before any agent or
   gate run; gate cache line `<verdict> <tree-hash> <iso8601>` written only
   with a real hash — a missing binary or hash degrades to no verdict, never a
   forged one. Env contract carries: `BENCH_HOME`, `BENCH_AGENT`,
   `BENCH_MAX_ITERS`, `BENCH_REFACTOR_ITERS`, `BENCH_GATE`. New: the
   resolution-order contract (#2).
3. **Deep vs thin.** The binary is the deep unit: lease state machine, pool
   scan/mint, touched-path staging diff, the iteration and refactor loops,
   gate resolution + recording, prompt text (#4). `bin/bench.sh` stays a thin
   router; the `run_gate` adapter and hook shims stay one-glance.
4. **Black-box assertables.** The existing fragments before and after the flip
   — identical stdout/exit/filesystem/git assertions. New `go test` tables:
   lease reclaim decision (live pid, dead pid, non-numeric legacy by mtime,
   fresh-empty mid-claim), touched-path diff (spaces, globs, scratch
   exclusion — porcelain `-z` parsed natively, which retires the shell's
   known newline-in-path misread), pool candidate naming, resolution order as
   a pure function.
5. **Gate attachment.** The unchanged shell gate is the oracle, plus
   `go build`/`vet`/`test`. Gate-blind: real-terminal Ctrl-C delivers SIGINT
   to the process group (the contract approximates with `kill -INT $PPID`) —
   one manual interrupt smoke when the spec lands; auto-detect chains beyond
   the single contracted case stay best-effort as today.
6. **Hostile-input owners.** Spaced/glob touched paths → Go staging via argv
   `:(literal)` pathspecs (contracted). SIGINT/SIGTERM mid-loop → signal
   handling + deferred cleanup, adapter child killed with the group
   (contracted at the observable layer). Zombie, legacy, and mid-claim leases
   → reclaim table tests. Missing `origin` / unset `origin/HEAD` →
   default-branch fallback carries. Unwritable pool home → acquire's loud
   failure. Byproduct-emitting gates → staging snapshot semantics
   (contracted). Repeat acquire after crash → lease hardening contracts.
7. **Uncertainty flags.** Signal forwarding to the adapter child in Go
   (process-group handling, making sure cleanup runs exactly once) is the one
   implementation area with real platform nuance — spec pins the observable
   contract (exit non-zero, scratch gone, lease released, pool reusable) and
   escalates per `craft-line` only if Go's semantics can't reproduce it.
8. **Rejected alternatives.** Shelling back into bash `run_gate` (#1);
   dropping auto-detect inside the port (#2 — open separately); a separate
   contract-backfill spec (#3); kit prompt files (#4).
9. **Domain watch-outs.** The dispatcher routes whole subcommands, never
   shares one — `shift`, `worktree`, and `gate` flip to Go together with their
   shell bodies deleted in the same diff. Lease ownership is the recording
   process's pid: acquire, loop, and release must run in one process, and a
   release by anyone but the recorded live owner is a no-op. The gate cache
   must never contain a forged verdict — no hash, no write, loudly. The
   interactive subshell inherits the user's stdio and releases on exit
   regardless of the shell's exit status.

Dependency order: single spec; lands after `go-doctor-link-port` per the parent
map's slice order.
