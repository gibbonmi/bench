# go-status-port — slice 3 of the Go rewrite

Child of `decisions/go-rewrite.md` (#6, slice order). Slice 3 as mapped is "the
`bench status` renderer", but `bin/bench-status.sh` is entangled across three
future slices: `structure_check` feeds `status`, `bench structure`, and the shell
shift loop (`structure_touched_since`); `gate_tree_hash` has three consumers
(status, the shift loop's gate-cache write, and a hand-mirrored copy in
`.bench/hooks/stop.sh` — one of the two live duplication defects the parent map
named); and status's worktree row reads pool-path and lease-file conventions
owned by `bin/bench-worktree.sh` (slice 7). This map settles who owns what
during the slice-3→7 window. The regression net is strong: ~10 black-box status
contracts in the runtime gate fragment drive `bench.sh status` and survive the
route flip unchanged.

## #1: What is slice 3's scope boundary?

Type: Grill

### Question
Status renderer only, status + structure, or the whole of `bin/bench-status.sh`
(adding `models`, `idea`, `roadmap`, the specs-retirement counter,
`gate_tree_hash`)?

### Answer
**Whole file; `bin/bench-status.sh` is deleted at the end of the slice.**
Structure must port anyway — Go status needs its violation count in-process,
and the shift loop can reach it through a thin adapter (the slice-2 pattern).
`models`/`idea`/`roadmap` are trivial riders, and porting `models` drops the
python3+curl dependency. Rejected: status-only (Go status would exec bash for
the structure count; the file survives as a straggler of orphaned helpers) and
status+structure (leaves three small commands for a later sweep nobody owns).

## #2: Who owns gate_tree_hash?

Type: Grill

### Question
The hash moves to Go with status, but the shift loop's gate-cache write and the
mirrored copy in `.bench/hooks/stop.sh` are shell until slices 7 and 5. Retire
the mirror now, or carry the duplication defect two more slices?

### Answer
**Go owns it; the stop.sh mirror dies this slice.** The hash is exposed as a
plumbing subcommand; the shift loop's cache write and stop.sh both call the
binary. Fixes a named live defect now — stop.sh's edit is one call site, not a
slice-5 rewrite. Caveat carried to the Handoff: stop.sh cannot source
`bench.sh`, so it needs its own binary-path resolution; if the binary is
missing the hook must fail safe, never write a verdict keyed to a guessed tree.

## #3: Who owns the worktree pool/lease facts during the slice-3→7 window?

Type: Grill

### Question
Go status must exclude warm pool worktrees, which requires the pool directory
path and lease-file naming currently owned by `bin/bench-worktree.sh`
(slice 7). Duplicate the two facts temporarily, or invert ownership now?

### Answer
**Go owns, shell adapts.** The pool-path and lease-file logic port to Go;
`bin/bench-worktree.sh` reads them back through thin plumbing calls, same
pattern as #2. No duplication window — the repo standard grades one against
diffs — and slice 7 later absorbs the shell callers rather than re-porting the
facts. Rejected: temporary duplication flagged for slice-7 cleanup (a defect
window held open across four slices).

## Handoff

1. **Module boundaries.** New `internal/` packages: the status renderer
   (rows, severity sort, five-row budget, footer), structure (checker +
   budgets parser, whole-tree and touched-scope modes), the specs-retirement
   counter, and the tree-hash + worktree pool/lease path facts (new packages
   or homes in the existing `internal/git` — spec-writer's layout call).
   `cmd/bench` gains `status`, `structure`, `models`, `idea`, `roadmap`, and
   the plumbing subcommands from #2/#3. The shell router flips all of these to
   `route_binary`; `bin/bench-status.sh` is deleted. Remaining shell callers
   (shift loop, `bin/bench-worktree.sh`, `.bench/hooks/stop.sh`) become thin
   adapters over the plumbing.
2. **Contracts.** All existing observable behavior carries unchanged: status
   stdout shape (lead `▶` line, severity order, five-row budget + `+N more`,
   footer), signal set and trigger conditions, `BENCH_LEARNINGS_FLOOR`; gate
   cache format `<status> <tree-hash> <iso8601>` read and written; structure
   output lines (`FILE TOO LONG`/`DIR CROWDED`), stderr debt summary, exit 1
   on violations, `BENCH_MAX_LINES`/`BENCH_MAX_DIR_FILES`,
   `.bench/structure.budgets` format including warn-and-skip on malformed
   lines; `idea` exit 2 on empty text and trailing-newline normalization;
   `roadmap` cat-or-empty; `models` both branches (API list vs no-key
   guidance). New plumbing contracts (exact names/flags are the spec's call):
   tree-hash prints the hash or `none`, exit 0; structure touched-scope mode
   for the shift loop; worktree pool/lease path queries for
   `bin/bench-worktree.sh`. `bench maps --count` remains public — only
   status's shell adapter around it dies.
3. **Deep vs thin.** The binary stays the deep unit — all parsing, counting,
   severity logic, and the tree-hash computation live behind subcommands. The
   shell call sites in `bench.sh`, `bin/bench-worktree.sh`, and
   `.bench/hooks/stop.sh` are one-glance adapters with no logic of their own.
4. **Black-box assertables.** The existing status contract suite (clean,
   footer, stale-gate, fresh-green, decisions, maps-count sourcing, severity
   budget, warm-pool exclusion, gate-cache write format, retirement signal)
   runs unchanged against the ported binary. New `go test` table tests:
   budgets parsing, specs-counter fence/CRLF/no-`Status:` edges, severity
   sort + row budget, gate-cache parsing, pool/lease path derivation.
5. **Gate attachment.** The unchanged shell gate is the oracle; the Go layer
   (`go build`/`vet`/`test`) already runs in its parse phase. Gate-blind spot:
   stop.sh's hook-context behavior (binary missing / not a repo) — cover with
   a contract if the runtime-shift fragment can reach it, else TDD + manual
   verify, flagged in the spec.
6. **Hostile-input owners.** Missing/stale binary in hook context → stop.sh
   adapter fails safe: skip the cache write loudly, never record a verdict for
   a guessed tree (status's own posture on a missing binary — count 0, render
   anyway — stays). No-trailing-newline `ROADMAP.md` and hand-edited specs →
   the idea appender and the specs/roadmap parsers, table-tested. CRLF and
   fenced `Status:` markers → specs counter. Malformed budget lines →
   warn-and-skip, preserved verbatim. Worktree paths with spaces →
   `git worktree list --porcelain` parsing in Go, argv exec. Cwd deeper than
   root → root resolution in the binary, asserted per subcommand.
7. **Uncertainty flags.** stop.sh's binary-path resolution without sourcing
   `bench.sh` (plant a resolver line vs duplicate the lookup vs a planted
   shim) — settle at spec time; if no shape avoids duplicating the resolution
   knowledge, escalate per `craft-line` rather than accept a second copy.
   Everything else is settled.
8. **Rejected alternatives.** Status-renderer-only and status+structure
   scopes (#1); leaving the stop.sh mirror until slice 5 (#2); temporary
   duplication of the worktree facts (#3).
9. **Domain watch-outs.** After this slice the binary sits on the gate-cache
   write path, not just the query path — a missing platform binary degrades
   verdict recording, not only dashboards; the fail-safe posture in item 6 is
   what keeps that honest. The session-start hook consumes status stdout
   verbatim, so row format and severity order are behavior, not cosmetics.
   Warm pooled worktrees (released, no lease) are deliberately not a signal —
   only leased pool entries and out-of-pool worktrees count.

Dependency order: n/a — single spec.
