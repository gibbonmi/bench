# go-hooks-port — slice 5 of the Go rewrite, plus TOON library adoption

Child of `decisions/go-rewrite.md` (#6, slice order). Slice 4 already applied the
slice-5 shim pattern to `block-dangerous-git.sh`; this slice ports the rest of the
hook surface: `stop.sh`, `check-agent-line.sh`, the adapters' `_line-guard.sh`,
and the shared `lines-env.sh` tier parser — after which the kit's last runtime
python3 calls are gone. `session-start.sh` is already the target shim shape (a
thin wrapper over `bench status` / `bench guards --brief`); it only swaps to the
shared resolver from #3. Mid-grill the reviewer added a second commitment: replace
the hand-rolled Go TOON emitter with the actual TOON library (the standing
prefer-official-libraries rule), as its own spec sequenced first. The regression
net for the hooks half: `gate-line-contracts.sh` covers both line guards
end-to-end, `gate-runtime-contracts.sh` pins stop.sh's envelope, missing-binary,
and bare-PATH paths, and the guards fragment covers every `--describe` manifest.

## #1: Does `_line-guard.sh` port now, or ride with slice 7's shift-loop port?

Type: Grill

### Question
Its only consumers are the three shift adapters, which slice 7 rewrites anyway —
but it shares `lines-env.sh` with `check-agent-line.sh`, which ports now.

### Answer
**Port it now.** Porting both consumers together kills `lines-env.sh` and gives
the lines.env binding parse one Go source; the adapter edit is a one-call-site
swap (source-a-function → exec a binary subcommand), too small for slice-7 churn
to argue for waiting. Rejected: defer to slice 7 (leaves the shell parser alive
for one consumer).

## #2: Where does the `stop.sh` port boundary sit?

Type: Grill

### Question
Does the binary own the whole verdict path (envelope parse, `stop_hook_active`,
`BENCH_SHIFT`, running the gate, the `bench-last-gate` cache write, the BLOCKED
message), or does the shim keep the gate run and cache write?

### Answer
**Whole verdict path in the binary.** Matches the guard-git shape (shim = pipe
envelope + fail rims), the gate is a subprocess the core can exec through the
same wrapper, and the cache write already depends on Go's tree-hash — splitting
would spread the verdict fact across two languages. The shim keeps `--describe`
and the binary-missing rim: fail open with the loud no-verdict warning, as today.
Rejected: shim-keeps-gate-run split.

## #3: How does the `resolve_wrapper` duplication collapse?

Type: Grill

### Question
The bench-wrapper search (repo `.bench/bin/bench.sh` → `bin/bench.sh` → `bench`
on PATH) is inlined in three hooks, and every shim needs it after this slice. A
blindly-sourced shared lib gives fail-closed guard shims a new fail-open mode
(missing lib → shim errors → non-2 exit → silent grant).

### Answer
**Shared `.bench/lib/resolve-bench.sh`, source-guarded per shim under its own
posture.** Guard shims treat "lib missing" exactly like "core missing" (git guard
refuses git-shaped input; agent-line warns and allows); stop and session-start
warn and fail open. Same pattern `lines-env.sh` set: one source for the search
order, each shim's fail posture stays its own explicit fact. Rejected: keeping
inlined copies; a blind unguarded source.

## #4: Where does the TOON-library swap sit?

Type: Grill

### Question
Ride inside the hooks spec as its first story, or its own small spec ahead of it?

### Answer
**Own spec, sequenced first.** The AXI contract fragments are the swap's
regression net (updated to spec bytes in the same spec, per #6), and it shares
nothing with the hooks work; mixing it into the hooks spec couples two
unrelated seams in one diff. Rejected: one
combined spec; keeping the hand-rolled emitter (overruled by the standing
prefer-official-libraries rule).

## #5: Which TOON library, and how does it enter the repo?

Type: Research

### Question
Confirm the official Go implementation (start at the `github.com/toon-format`
org), its maturity, maintenance state, and license; decide normal `go.mod`
dependency vs vendoring (the kit has no third-party Go deps yet — this sets the
precedent); map what of `internal/`'s emitter survives as a thin adapter vs is
deleted. Output: a short summary asset with the recommended library and inclusion
shape. If no official Go implementation exists in usable state, that is a finding
for the reviewer, not a license to keep the hand-rolled emitter silently.

### Answer
**`github.com/toon-format/toon-go`, as a normal `go.mod` dependency pinned to a
pseudo-version; no vendoring.** It is the Go implementation under the official
spec-owning org: MIT, active, `Marshal` with ordered fields via struct tags,
comma/2-space defaults matching the kit — but pre-v1 with no tagged release, so
the pin is a pseudo-version and bumps are deliberate edits. Consumers never
build Go (prebuilt binaries), so the dependency is kit dev/CI-only and
`go.mod`+`go.sum` already make it reproducible. Precedent set for third-party
Go deps: official-org, MIT-compatible, build-time-only. Research surfaced that
the swap is **not byte-identical** — the hand-rolled emitter deviates from
spec-TOON — which is #6, a reviewer call. Full comparison and rejected
candidates: `decisions/assets/go-toon-library.md`.

## #6: Accept the spec-TOON output change across the AXI surface?

Blocked by: #5
Type: Grill

### Question
`internal/toon` is TOON-shaped but not spec-TOON: doubled-quote escaping vs the
spec's backslash escapes, raw newlines vs `\n`, and far narrower quoting
triggers (spec also quotes empties, `true`/`false`/`null`, numeric-looking
strings, colons, backslashes, brackets/braces, leading `-`). Adopting the real
library therefore changes AXI stdout bytes: every contract fragment asserting
exact TOON lines updates, and row call sites move from `[][]string` to typed
cells so numeric columns stay bare. Accept the output change (spec-compliant
TOON, contracts updated in the same spec), or keep the hand-rolled emitter
(already overruled in #4) — the middle path of shimming the library to emit the
old bytes defeats the point of adopting it.

### Answer
**Accepted — AXI stdout moves to spec-TOON bytes.** The TOON spec updates the
contract fragments in the same diff and moves row call sites to typed cells so
numeric columns stay bare. The block shape, exit codes, and the error/usage
lines are unchanged. Rejected: shimming the library to reproduce the old
non-conformant bytes.

## Handoff

1. **Module boundaries.** New `internal/` package(s) for hook verdict logic —
   per-hook split vs one package is the spec's call, mirroring how `gitguard`
   got its own. `cmd/bench` gains plumbing subcommands (names are the spec's
   call): the stop verdict, the agent-line check, and the adapters' model
   resolution. `.bench/hooks/*.sh` all become or stay thin shims;
   `_line-guard.sh` and `lines-env.sh` are deleted (adapters exec the binary
   subcommand directly); new shared `.bench/lib/resolve-bench.sh` (#3). TOON
   spec: `internal/toon`'s emitter (`Escape`/`Table`/`IsSpace`) is superseded
   by `github.com/toon-format/toon-go` (#5); the AXI error/usage line helpers
   survive — they are the hybrid contract, not TOON.
2. **Contracts.** All observable interfaces carry unchanged. Stop hook: JSON
   envelope on stdin; enforce only when `BENCH_SHIFT=1`; honor
   `stop_hook_active`; exit 0 allow / exit 2 with the BLOCKED message and gate
   tail; cache line `<status> <tree-hash> <iso8601>` in the git dir; a missing
   tree-hash warns and writes nothing — never a forged verdict. Agent-line
   guard: every degraded branch warns on stderr and exits 0; only a present
   model matching no bound tier/alias exits 2 with the DENIED message naming
   the binding. Adapter resolution: unrouted passthrough, incomplete-binding
   warn + passthrough, routed repo requires a bound `BENCH_MODEL` or exit 1 —
   same rules, new shape: resolved model on the subcommand's stdout instead of
   a sourced variable. `--describe` manifests keep their live-binding denies
   clauses (one source per fact). TOON spec: AXI stdout moves to spec-TOON
   bytes (#6); the contract fragments update in the same diff and stay the
   assertion surface — the block shape (`name[N]{fields}:` + indented rows),
   exit codes, and the error/usage lines are unchanged.
3. **Deep vs thin.** The binary is the deep unit: envelope and binding parsing,
   verdicts, gate orchestration, cache write, and every message live behind the
   subcommands. Shims are one-glance pass-throughs owning only `--describe`,
   wrapper resolution, and their posture rim. `resolve-bench.sh` carries the
   search order only — no policy.
4. **Black-box assertables.** `gate-line-contracts.sh` and the stop.sh cases in
   `gate-runtime-contracts.sh` run unchanged against the ported hooks — that is
   the port-parity net; the adapter cases re-point from sourced-function to
   exec shape. New `go test` tables: lines.env parsing edges (quoting,
   whitespace, missing keys, no trailing newline), envelope edges (non-JSON,
   missing fields, `stop_hook_active` variants), and each verdict branch.
5. **Gate attachment.** The unchanged shell gate is the oracle. The
   `lines-env-broken` and `alias-line-broken` canaries guard the shell parser
   and retire with it (guard-slice precedent: the property — a broken binding
   parse degrades loudly, never silently — moves to `go test` tables and
   re-pointed contracts). Gate-blind: real harness wiring end-to-end (Claude
   Code actually invoking the hooks) stays manual smoke, as today.
6. **Hostile-input owners.** Non-JSON stdin / absent envelope fields → the
   binary's parser, fail-open per hook posture, table-tested. Hand-edited
   lines.env (quotes, spaces, no trailing newline) → the Go binding parser
   tables. Binary or lib missing in hook context → each shim's posture rim
   (#3), contract-pinned. Bare `PATH=/usr/bin:/bin` invocation → the existing
   runtime contract case survives the port. Oversized gate output → the
   tail-30 truncation moves into the binary, asserted.
7. **Uncertainty flags.** Codex-side envelope shape for the stop and agent-line hooks: the
   spec verifies from the adapters what `.codex/hooks.json` wires and whether
   its `tool_input`/`stop_hook_active` shapes match Claude Code's; if
   unreadable there, escalate per `craft-line` rather than guess.
8. **Rejected alternatives.** Deferring `_line-guard` to slice 7 (#1); the
   shim-keeps-gate-run split (#2); inlined resolver copies and blind shared
   source (#3); one combined spec and keeping the hand-rolled emitter (#4).
9. **Domain watch-outs.** After this slice every interactive enforcement
   surface depends on the platform binary, and degradation postures differ by
   hook on purpose: guards refuse their shaped input, informational hooks go
   silent, the stop hook allows loudly without a verdict — and a missing
   binary must never forge a gate verdict. `bench guards` aggregates `*.sh`
   manifests only; ported logic must not grow a second `--describe` answerer.
   The adapters are sourced POSIX-style — the exec swap must stay POSIX-clean.

Dependency order: the TOON spec (blocked on #6) lands first — reviewer's
sequencing — then the hooks spec; both precede slice 6 in the parent map.
