# go-shift-worktree-port

Status: staged

## Problem

The warm-worktree pool and the `bench shift` gated loop are the last operational
surface still carrying their logic in sourced shell. `bin/bench-worktree.sh` (103
lines) owns the lease state machine (atomic O_EXCL claim, owner-checked reclaim of
a provably-gone owner, owner-only release), the acquire scan-or-mint over the pool,
and the interactive subshell. The shift half of `bin/bench.sh` (~180 lines) owns
the gated loop: the iteration and refactor phases, touched-path staging, the gate
resolution chain (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect) plus verdict
recording, the adapter preflight, and the SIGINT trap that must tear the whole run
down cleanly. Both are sourced into the dispatcher and dispatched as shell
functions.

The strangler has reached slice 7 of `decisions/go-rewrite.md`. The regression net
is strong and fully black-box — two fragments drive `bench shift`/`worktree`/
`gate` as a subprocess: `gate-runtime-shift-contracts.sh` (green commit + benchBase,
touched-path staging incl. byproduct/pre-dirt/spaces-globs, red rollback, done.sh
early-completion, refactor trigger/no-trigger/no-op, refactor-prompt scope, SIGINT
cleanup, the adapter preflight + single-argument contracts) and
`gate-runtime-contracts.sh` (lease hardening, reuse/release cleaning, concurrent
acquire, the interactive subshell via a fake `$SHELL`, `BENCH_GATE` cwd, gate
repo-root cwd, verdict-record incl. the gate-cache assertions, and both
missing-core fail-safes). But the logic
those contracts guard is shell, so the reclaim decision, the porcelain touched-path
parse, and the resolution precedence are shell-untestable except through a full CLI
round-trip — and the shell `shift_dirty_paths` carries a documented misread of a
path containing a literal newline. One contract gap remains: the resolution chain
has **no** assertion of its precedence order at all.

## Solution

The worktree lifecycle and the loop move into Go behind three subcommands, with the
shell bodies deleted in the same diff — mirroring how slice 6 folded link/init/
doctor into `internal/adopt` and deleted their `.sh` files. `internal/worktree`,
which already owns the pool-path and lease-path addressing, grows the lifecycle it
addresses: the lease claim/reclaim state machine, acquire scan-or-mint, owner-only
release, and the interactive-subshell entry. A new `internal/gate` owns the
resolution chain, the gate run, and the verdict-cache record — the oracle's
selection logic in one Go home, called both by the standalone `bench gate` and,
in-process, by the loop. A new `internal/shift` owns the iteration and refactor
loops, the touched-path staging diff (porcelain `-z` parsed natively, which retires
the shell's newline misread), the adapter preflight, the `.bench/done.sh`
early-completion check, and the iteration/refactor prompt text as Go string
constants (ticket #4 — prompts have always lived inside the executable).

`cmd/bench` gains `shift`, `worktree`, and a `gate-run` plumbing subcommand.
`bin/bench-worktree.sh` is deleted; `bin/bench.sh` drops every shift/gate/worktree
body and keeps only routing plus a one-glance `run_gate` adapter (the
`worktree_pool` precedent) that forwards to `bench gate-run` — so `bench gate` and
the Stop hook's `<wrapper> gate` path keep their shape while gate resolution lives
in exactly one place. The loop calls `internal/gate` directly rather than shelling
back out, because a Go→bash→gate chain would be a second live gate-resolution
implementation, which ticket #1 forbids.

One new behavior lands, owned by this slice per ticket #2: the resolution-order
contract. `.bench/gate.sh` beats `$BENCH_GATE` beats auto-detect, with one
representative detection case and a no-gate case (exit 3, nothing recorded). This
is a product-behavior *assertion* over a chain that carries unchanged — not the
removal of auto-detect, which stays open as its own small decision. The two
fragments run unchanged and are the port-parity net; new `go test` tables cover
what the shell never reached: the four-way reclaim decision, the pool candidate
name, the touched-path diff, and the resolution precedence as a pure function.

## User stories

1. As a reviewer running warm isolated worktrees, I want the worktree lifecycle
   ported to `internal/worktree` — the atomic O_EXCL lease claim recording
   `<pid> <utc-time>`, reclaim only of a provably-gone owner (a recorded pid no
   longer running, or unreadable/legacy content aged out by mtime, never a
   fresh-empty writer mid-claim), acquire's scan-of-clean-released-then-mint over
   the pool, owner-only release that a non-owner's deferred cleanup leaves alone,
   and the `bench worktree` interactive subshell that inherits the user's stdio and
   releases on any exit status — so that the pool the shift loop and I both depend on
   runs from the Go core with its lease semantics intact. Line: claude-sonnet-5 /
   medium. Every reclaim, concurrency, and reuse case is pinned by an existing
   contract that drives the CLI as a subprocess, so the gate fully grades the port
   and the cheap tier fits; medium effort because the atomic claim and the reclaim
   mtime threshold are load-bearing and a silent divergence steals a live lease.

2. As the shift loop, I want the touched-path staging diff ported to
   `internal/shift` — the dirty-path set parsed from `git status --porcelain -z`
   natively (NUL-delimited, so spaces, glob characters, and a literal newline in a
   path all survive), the scratch files excluded, and the post-minus-pre difference
   staged via `:(literal)` pathspecs so a gate byproduct or pre-existing dirt never
   rides into an iteration commit — so that an iteration commits exactly what the
   agent touched and the shell's known newline-in-path misread retires. Line:
   claude-sonnet-5 / medium. This is a mechanical port of a pure parsing function
   whose byproduct/pre-dirt/spaces-glob cases are pinned by an existing contract and
   whose newline case is pinned by a new Go table, so the gate fully grades it;
   medium effort because a mis-parsed diff stages the wrong paths.

3. As the gate, I want the resolution chain and verdict recording ported to
   `internal/gate` — `.bench/gate.sh` beats `$BENCH_GATE` beats auto-detect with the
   no-gate case exiting 3, the gate run from the repo root, and the verdict written
   as `<verdict> <tree-hash> <iso8601>` keyed to `git.TreeHash` (nothing written when
   the binary or hash is unavailable, never a forged verdict) — plus a NEW
   resolution-order contract pinning that precedence, which has no assertion today —
   so that the oracle's selection logic has one Go home that `bench gate` and the
   loop both read. Line: claude-opus-4-8 / medium. Correctness of the oracle matters
   more than speed (the profile routes gate logic to mid effort), and the precedence
   is a freshly-contracted behavior where a reordered chain would silently run the
   wrong oracle, so it takes the mid tier at medium effort.

4. As the shift loop, I want the iteration and refactor orchestration ported to
   `internal/shift` — the bounded iteration loop (`BENCH_MAX_ITERS`) that commits on
   green and rolls back on red while preserving scratch, the `.bench/done.sh`
   early-completion check, the touched-scope refactor phase that triggers only on
   debt this shift touched (never pre-existing), runs within the
   `BENCH_REFACTOR_ITERS` budget (default 4), and stops on a no-op pass, the
   adapter preflight rejecting an unset/
   empty/missing/non-executable/shell-keyword `BENCH_AGENT` before any run, the
   adapter invoked with the generated prompt as its single positional argument under
   `BENCH_SHIFT=1`, and the iteration/refactor prompt text as Go string constants —
   so that the gated loop runs from the binary with every stdout literal, benchBase
   record, and commit-boundary preserved. Line: claude-opus-4-8 / medium. This is
   loop orchestration with many branches, most pinned by the unchanged shift
   contracts, but a regression that commits on red or refactors repo-wide debt is a
   correctness failure the gate grades only at the observable edge, so it takes the
   mid tier at medium effort.

5. As a reviewer who pulls the line mid-run, I want SIGINT/SIGTERM handling ported to
   `internal/shift` with the observable teardown contract pinned — an interrupted
   shift exits non-zero, its `.bench-objective` and `.bench-notes.md` scratch are
   gone, the pool lease is released, the pool is reusable by a follow-up shift, the
   adapter child is killed with the process group, and cleanup runs exactly once — so
   that Ctrl-C tears the whole run down cleanly rather than leaking a leased
   worktree. Line: claude-opus-4-8 / high. Signal forwarding to the adapter child
   (process-group handling, cleanup-exactly-once) is the map's one uncertainty flag
   with real platform nuance, and the real-terminal Ctrl-C delivery is gate-blind
   past the contract's `kill -INT $PPID` approximation, so it takes the mid tier at
   high effort with a manual interrupt smoke when the spec lands.

6. As the strangler dispatch and the gate, I want `bin/bench.sh` to route `shift`,
   `worktree`, and `gate` to the binary — `bench-worktree.sh` deleted, the shift/gate/
   worktree bodies dropped, `run_gate` reduced to a one-glance adapter over
   `bench gate-run`, `cmd/bench` gaining the three subcommands, and every subcommand
   flipping whole in the same diff so no seam has two live implementations — with no
   existing contract assertion weakened and no dangling reference left in contracts,
   the dispatcher, README, or link/package — so that the port leaves the dispatcher,
   the gate load, and the stale-reference sweep green. Line: claude-opus-4-8 / high.
   Touching the dispatcher and the gate path is the worst defect class in this kit
   (`craft-gate`) — deleting a sourced file and re-pointing the gate runner without a
   second live resolver — so it takes the mid tier at high effort.

## Implementation decisions

- **Package split (spec's call — the map left it open, following the `gitguard` /
  `adopt` precedent).** Three homes, each a deep unit split by concern:
  - `internal/worktree` **grows** the lifecycle it already addresses: `Claim`/
    reclaim (the lease state machine), `Acquire` (scan-or-mint), `Release`
    (owner-only), and the interactive-subshell entry. It already owns the pool-path
    cksum and the lease git-path, so the lifecycle belongs beside them.
  - `internal/gate` (new): `Resolve` (the ordered chain as a pure function), `Run`
    (execute from root, capture rc), `Record` (the verdict-cache write via
    `git.TreeHash`, reused in-process — never a second hash). It gets its own package
    because `bench gate` and the Stop hook need the gate runner independently of the
    loop; folding it into `internal/shift` would couple the standalone oracle to the
    loop package.
  - `internal/shift` (new): the iteration and refactor loops, the touched-path
    staging diff, the adapter preflight, the `objective_met`/`.bench/done.sh` check,
    and the prompt constants. It imports `internal/worktree`, `internal/gate`, and
    `internal/structure` (the touched-scope refactor detector, called in-process).
  Filesystem and git truth inject at each package boundary so the reclaim decision,
  the resolution order, and the touched-path diff are unit-testable without mutating
  a real tree — the `gitguard.Checker` seam pattern.

- **Acquire, loop, and release run in one process (watch-out #9).** Lease
  ownership is the recording process's pid, so `bench shift` performs
  acquire→loop→release inside a single binary invocation — never as separate
  plumbing subcommand calls — and a release by anyone but the recorded live owner
  stays a no-op.

- **The env contract carries whole (Handoff item 2).** `BENCH_HOME` (pool
  location), `BENCH_AGENT` (the adapter), `BENCH_MAX_ITERS` (iteration cap),
  `BENCH_REFACTOR_ITERS` (refactor budget, default 4), and `BENCH_GATE`
  (resolution chain) are read by the Go core under today's names and defaults;
  the shift contracts already exercise the iteration and refactor budgets.

- **The loop calls `internal/gate` in-process; the shell `run_gate` forwards to
  `bench gate-run` (ticket #1).** Resolution and recording move fully into Go. The
  loop runs the gate by calling `internal/gate` directly, not by shelling back out —
  a Go→bash→gate chain would be a second live gate-resolution implementation. The
  standalone path keeps its shape: `bin/bench.sh`'s `gate)` case → the one-glance
  `run_gate` adapter → `bench gate-run "$(repo_root)"`; the Stop hook's
  `<wrapper> gate` (slice 5, `internal/stophook`) hits the same `gate)` case
  unchanged. `gate-run` is the plumbing subcommand name (spec's call, following
  `worktree-pool`/`tree-hash`).

- **The resolution-order contract is an assertion, not a behavior change (ticket
  #2).** The chain carries byte-for-byte; this slice only adds the missing precedence
  and no-gate assertions. Dropping auto-detect stays open as its own decision and is
  out of scope here.

- **The prompt text ports to Go string constants (ticket #4).** `iteration_prompt`
  and `refactor_prompt` have always been heredocs inside the executable, never
  reviewer-tunable files; the parent map drew the line at "executable logic ports,
  markdown content stays text." No `.bench/prompts/*` surface is created.

- **Deletion is part of this change, not a follow-up.** `bin/bench-worktree.sh` is
  deleted with its source line removed from `bin/bench.sh`; the shift/gate/worktree
  function bodies (`shift_loop`, `shift_stage_touched`, `shift_dirty_paths`,
  `cleanup_shift_scratch`, `shift_cleanup`, `require_adapter`, `objective_met`,
  `iteration_prompt`, `refactor_prompt`, `gate_record`, `structure_touched_since`,
  `worktree` and the `worktree_*` helpers) are removed; `run_gate` shrinks to the
  adapter. `bin/bench.sh` keeps the dispatch, `run_gate`, `route_binary`,
  `adoption_route`, `bench_binary_path`, and the path/kit resolvers. All in one diff
  so the gate load, the two contract fragments, and the docs stale-reference sweep
  stay green together.

- **Drift absorbed since the map closed (no deviation).** Slice 6 has already
  landed: `internal/adopt` exists, `bench-link/init/doctor.sh` are deleted, and
  `bin/bench.sh` routes those three through `adoption_route`. The map's item-1 line
  "slice 6 separately deletes the link/init/doctor sources" is now history, so this
  slice touches only the shift/gate/worktree bodies and leaves `adoption_route` and
  the adopt routing alone. After this slice, `bin/` keeps only `bench.sh` and
  `bench-postinstall.sh` as shell; the generated pre-push and hook entry shims stay
  shell as the map specified.

## Testing decisions

- **What a good test is here.** Acceptance drives each subcommand end-to-end through
  the CLI as a subprocess — run `bench shift`/`worktree`/`gate` against a throwaway
  repo (with a stub `BENCH_AGENT` or a fake `$SHELL`) and assert stdout literals,
  exit code, filesystem, and git state — never Go internals. The two existing
  fragments already do exactly this and are the port-parity net; they run with their
  assertions intact. Go table tests are additional, at the injected filesystem/git
  boundary, where the shell-untestability tax on the reclaim decision, the
  touched-path parse, and the resolution precedence finally retires.

- **Seams.** Two, the fewest that exercise the real behavior:
  - The **shift/gate contract seam** — the two fragments driving the CLI subprocess
    (`gate-runtime-shift-contracts.sh` and `gate-runtime-contracts.sh`, both the
    parity net), plus the new resolution-order contract. Prior art: the fragments
    themselves, unchanged in shape.
  - The **`go test` unit seam** — table tests beside `internal/worktree`,
    `internal/gate`, and `internal/shift` at the injected boundary. Prior art:
    `internal/gitguard` (Checker-injected tables), `internal/adopt`, `internal/maps`.

- **Gate command:** `bench gate` (the project gate, whose Go layer already runs
  `go build`/`go vet`/`go test ./...` and the cross-compile matrix, so the new
  packages and tests are graded with no new wiring), plus the parent map's per-slice
  done rule `go build ./... && go vet ./... && go test ./...`. Done = gate green and
  those three green.

### Seam diagram — shift/gate contract seam (CLI subprocess → binary)

    trigger: contract fragment runs `bench shift|worktree|gate` (+ a stub BENCH_AGENT / fake $SHELL) against a throwaway repo
        │
        ▼
    objective + BENCH_AGENT ─▶ [ bench shift → internal/shift: acquire pool wt → branch + benchBase ]
      (+ BENCH_MAX_ITERS)      [   iterate: adapter (BENCH_SHIFT=1, prompt as $1) ▸ gate ▸ stage touched ▸ commit-on-green / rollback-on-red ]
                               [   refactor phase at green over touched budget ▸ cleanup ▸ release ] ─▶ stdout literals + isolated shift branch, lease freed
    fake $SHELL ─────────────▶ [ bench worktree → internal/worktree: acquire → subshell (inherit stdio) → release on any exit ] ─▶ reused clean path, lease freed
    gate.sh|BENCH_GATE|pkg ──▶ [ bench gate → run_gate adapter → bench gate-run → internal/gate: resolve chain ▸ run from root ▸ record ] ─▶ exit 0/≠0 + `<verdict> <tree> <iso>` cache
        ◀ tests attach here: gate-runtime-shift-contracts.sh + gate-runtime-contracts.sh drive the CLI as a subprocess and assert
          stdout/exit/filesystem/git before and after the flip unchanged; the NEW resolution-order contract asserts the precedence.

### Seam diagram — `go test` unit seam (reclaim, diff, resolution)

    trigger: gate Go layer runs `go test ./...`
        │
        ▼
    lease bytes + pid state ──────▶ [ internal/worktree.<reclaim decision> ] ─▶ reclaim | respect
    root path ────────────────────▶ [ internal/worktree.<pool candidate>   ] ─▶ pooled candidate name
    porcelain -z bytes ───────────▶ [ internal/shift.<touched-path diff>    ] ─▶ touched paths (spaces/globs/newline kept, scratch excluded)
    gate.sh? BENCH_GATE? lockfiles ▶ [ internal/gate.<Resolve>              ] ─▶ gate.sh | BENCH_GATE | detect(one) | none→exit 3
        ◀ tests attach here: table tests pin the four-way reclaim decision (live pid, dead pid, non-numeric legacy by mtime,
          fresh-empty mid-claim), the pool candidate name, the touched-path diff (incl. the literal-newline path the shell
          misread), and the resolution precedence. Red before the package exists → does not compile.

### Acceptance coverage map

Per-item granularity is stated where a behavior quantifies over a set (each reclaim
case, each resolution branch, each preflight rejection).

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | lease claim/reclaim/reuse — atomic claim; reclaim on dead-pid and aged-out-empty; respect a live-foreign and a fresh-empty lease; concurrent acquires never share; release cleans dirty + ignored and removes the lease; interactive subshell reuses a clean released path | shift/gate contract | already covered — port-parity via the "lease hardening", "lease/reuse", "concurrent-acquire", and interactive-subshell contracts in `gate-runtime-contracts.sh`, which stay green across the flip | any drift in the reclaim decision, the atomic claim, or the owner check trips one of the unchanged lease contracts (`live foreign lease was stolen`, `fresh empty lease was stolen`, `concurrent acquires shared a worktree`) |
| 1 | reclaim decision as a pure function — live pid → respect; dead pid → reclaim; non-numeric legacy aged by mtime → reclaim; fresh-empty mid-claim → respect (each of 4) | go test unit | `internal/worktree` reclaim table before the fn exists → does not compile | pins the four-way decision the black-box contract exercises but can't enumerate cheaply; a mis-ported mtime threshold steals a mid-claim lease |
| 1 | pool candidate naming stays inside the pool dir | go test unit | `internal/worktree` pool-candidate table before the fn exists → does not compile | a wrong candidate name mints outside the pool and breaks warm reuse silently |
| 2 | an iteration commit carries exactly the touched paths — gate byproduct excluded, pre-agent dirt subtracted, space + glob chars preserved | shift/gate contract | already covered — port-parity via the stage-touched contract in `gate-runtime-shift-contracts.sh`, which stays green | a staging regression rides a byproduct into a commit or drops a spaced path → `gate byproduct rode into an iteration commit` / `touched path with space+glob chars was not staged` |
| 2 | touched-path diff pure function — spaces, globs, scratch exclusion, and a path containing a literal newline (retires the shell misread) | go test unit | `internal/shift` touched-path table incl. the newline path before the fn exists → does not compile | the porcelain `-z` native parse must keep NUL-delimited paths; the literal-newline case is the one the shell misread, unreachable by the space+glob contract |
| 3 | gate runs from the repo root; verdict recorded as `<verdict> <tree-hash> <iso8601>` keyed to the tested tree; missing binary/hash → nothing written (no forged verdict); commit-on-green does not stale a fresh verdict | shift/gate contract | already covered — port-parity via the "gate repo-root cwd", "BENCH_GATE cwd", "verdict-record", and both "missing-core fail-safe" contracts in `gate-runtime-contracts.sh` | a resolution/record regression trips an unchanged gate contract (`gate_record forged a verdict with no core binary`, `manual gate run did not record ... the tested tree`) |
| 3 | resolution precedence — `.bench/gate.sh` beats `$BENCH_GATE` beats one detection case; no gate → exit 3, nothing recorded (each branch) | shift/gate contract + go test unit | NEW resolution-order contract: a repo with an exit-0 `.bench/gate.sh` **and** an exit-1 `BENCH_GATE` must exit 0 (gate.sh wins); an exit-0 `BENCH_GATE` in a repo whose auto-detect path would fail must exit 0 (BENCH_GATE wins); `package.json`-only picks the npm path; no gate/no BENCH_GATE/no lockfile → exit 3. Cannot start red black-box — today's shell already resolves in this order, so the contract is written green pre-flip as a parity pin that bites post-flip. The genuinely red signal is the Go `Resolve` precedence table: does not compile before the package exists | no assertion pins the precedence today; a reordered chain would silently run the wrong oracle — the contract catches it at the flip, the Go table at the function |
| 4 | green iteration commits + records `branch.<name>.benchBase` + names the branch + leaves the main checkout untouched; red iteration rolls back with no commit and scratch surviving; `.bench/done.sh` ends the loop early; objective assembled from all positional args | shift/gate contract | already covered — port-parity via the green-commit/benchBase, red-rollback, scratch-survival, and done.sh early-completion contracts in `gate-runtime-shift-contracts.sh` | any loop regression (commit on red, lost benchBase, dirtied main checkout) trips an unchanged shift contract (`shift did not record the pre-shift HEAD in branch.<name>.benchBase`, `red shift preserved rolled-back work`) |
| 4 | refactor phase triggers only on debt this shift touched (never pre-existing), exits early on a no-op pass with no phantom commit, and scopes its prompt to the flagged touched files | shift/gate contract | already covered — port-parity via the touched-scope-structure, refactor no-op, and refactor-prompt scope contracts | a refactor-trigger regression refactors repo-wide debt or loops on a no-op → `pre-existing structural debt triggered refactor phase` / `no-op refactor pass reported a phantom commit` |
| 4 | adapter preflight rejects unset / empty / missing-path / shell-keyword `BENCH_AGENT` before any iteration; adapter invoked with the prompt as its single positional arg under `BENCH_SHIFT=1` (each rejection) | shift/gate contract | already covered — port-parity via the adapter preflight and single-argument contracts in `gate-runtime-shift-contracts.sh` | a preflight regression enters the loop with a broken adapter, or splits the multi-line prompt → `missing adapter still entered the loop` / `prompt was not the adapter's single positional argument` |
| 5 | SIGINT mid-loop → exit non-zero, `.bench-objective` + `.bench-notes.md` gone, lease released, a follow-up shift completes (pool reusable), cleanup runs exactly once | shift/gate contract | already covered — port-parity via the interrupt-cleanup contract in `gate-runtime-shift-contracts.sh` (approximates real Ctrl-C with `kill -INT $PPID`), which stays green | a Go signal / process-group or double-cleanup regression leaves a leased worktree or crashes on second cleanup → `interrupted shift left a leased worktree` / `follow-up shift after interrupt did not complete` |
| 5 | real-terminal Ctrl-C delivers SIGINT to the process group including the adapter child | manual smoke (gate-blind) | not gate-gradable — one manual interrupt smoke when the spec lands (map #5) | a real TTY process-group delivery the `kill -INT $PPID` approximation can't reproduce; escalate per `craft-line` only if Go's signal semantics diverge |
| 6 | `bench-worktree.sh` deleted; shift/gate/worktree bodies dropped from `bin/bench.sh`; `run_gate` is the one-glance adapter; `cmd/bench` gains shift/worktree/gate-run; no dangling reference in contracts, dispatcher, README, or link/package | gate load + docs stale-reference sweep | the sweep against a tree still naming a deleted file/function → red; a dispatcher still sourcing `bench-worktree.sh` → gate load (`bash -n` / source) red; the two unchanged fragments run against the flipped binary and must stay green | the conformance layer fails when a deleted file or function is still referenced or sourced, and a seam left half-shell would fail its port-parity contract (watch-out #9: two live implementations double-resolve) |

### Edge inventory

Walked per behavior against the profile's shell-CLI hostile-input checklist and the
map's item-6 owners; each resolved as a coverage row above or a **Won't handle**
line here.

- **paths/dirs with spaces or globs** — covered: the stage-touched contract
  (story 2) and the touched-path Go table; Go stages via `:(literal)` argv
  pathspecs, never a shell string.
- **path containing a literal newline** — covered and *fixed*: the porcelain `-z`
  native parse keeps NUL-delimited paths, retiring the shell's documented misread
  (story 2 Go table). Not a Won't-handle — it falls out of the port for free.
- **hand-edited / malformed lease file (no trailing newline, legacy content)** —
  covered: the reclaim table reads the pid field regardless of a trailing newline,
  and a non-numeric legacy lease ages out by mtime (story 1).
- **absent vs present-empty** — covered as distinct behaviors: a fresh-empty lease
  (writer mid-claim) is respected while an aged-out empty (crash) is reclaimed
  (story 1); an absent gate (no `.bench/gate.sh`, no `$BENCH_GATE`, no lockfile)
  exits 3 while a present gate runs (story 3).
- **unquoted multi-word args** — covered: the objective is assembled from all
  positional args (story 4, mirroring `$*`); the adapter receives the prompt as a
  single `$1` (story 4 single-argument contract).
- **required tool missing from PATH (no core binary)** — covered by construction:
  `bench shift` never reaches Go when the binary is absent (`route_binary`'s
  exit-127 install-message rim, unchanged); the gate/record skip their write on a
  missing binary (story 3 fail-safes).
- **invocation through a symlink** — safe by construction: `resolve_script_path`
  and `route_binary` already walk the symlink chain (the unchanged symlinked
  kit-dir portability contract); no new test.
- **SIGINT mid-loop: leftover scratch, leases, worktrees** — covered: the interrupt
  cleanup contract plus the manual smoke (story 5).
- **re-run idempotency (reused worktree, second shift)** — covered: the lease/reuse
  contract and the follow-up-shift-after-interrupt assertion (stories 1, 5).
- **cwd deeper than the repo root** — safe by construction: the gate runs from the
  root (story 3 repo-root cwd contract) and the loop resolves the main root via
  `git.Root`, as today.
- **missing `origin` / unset `origin/HEAD`** — safe by construction (parity):
  `git.DefaultBranch` falls back to `main` and acquire falls back from
  `origin/<branch>` to a plain `worktree add`, exactly as the shell does; no new
  test, no in-scope caller removes `origin`.
- **byproduct-emitting gates** — covered: the staging snapshot precedes the gate run
  (story 2 stage-touched contract).
- **repeat acquire after a crash** — covered: the aged-out-empty reclaim case
  (story 1 lease hardening).
- **Won't handle: unwritable pool home** — parity: acquire echoes
  `could not lease a pool worktree` and returns 1 as today; not separately
  gate-tested, since no in-scope caller provisions an unwritable `$BENCH_HOME`.
- **Won't handle: real end-to-end harness wiring** (Claude Code / Codex actually
  driving a shift) — gate-blind, a manual smoke as today (map #5); the contracts
  drive the loop with a stub `BENCH_AGENT`.
- **Won't handle: dropping the auto-detect chain** — open as its own small decision
  in the parent map (ticket #2 / #8), not this port; the chain carries unchanged.

## Out of scope

- **Slice 8 — gate fragments → `go test`** (the gate content ports last, under the
  canary layer's watch). A separate capability with its own spec; this slice ports
  the loop and the resolution chain, not the fragment *content*, which stays shell.
  Estimate to build later: one-plus spec-sized session (~30–40 edits, ~8–12 gate
  runs).

- **Dropping the auto-detect resolution chain** — a distinct product-behavior
  removal the map (ticket #2) deliberately kept out of this port so the strangler
  diff stays behavior-preserving; it needs its own small decision. Estimate if ever
  taken: ~4–6 edits, ~2 gate runs.
