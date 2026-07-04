# go-hooks-port

Status: staged

## Problem

Three of the kit's interactive hooks still carry executable logic in shell, and two of them still shell out to `python3` per invocation. `stop.sh` parses its envelope with inline `python3 -c`, resolves the bench wrapper, runs the gate, truncates the output, and writes the verdict cache — the completion oracle's whole verdict path lives in a bash script. `check-agent-line.sh` sources a shared shell parser (`.bench/lib/lines-env.sh`) and `python3 -c` to decide whether a delegation runs on a bound tier. The three shift adapters (`claude`, `codex`, `opencode`) source `.bench/adapters/_line-guard.sh`, which sources the same shell parser, to resolve `BENCH_MODEL` against the binding. The binding parse (`bench_tier_value`) therefore has one shell home shared by two enforcement surfaces, and the compensating canary layer (`lines-env-broken`, `alias-line-broken`, `adapter-line-broken`) exists only because that parse is shell-untestable.

Slice 4 already proved the shim pattern on `block-dangerous-git.sh`: the logic moved behind `bench guard-git`, the hook became a one-glance fail-closed shim, and the last standalone Python program died. But `stop.sh` and `check-agent-line.sh` still hold `python3 -c`, so python3-the-runtime is not gone kit-wide, and the lines.env binding parse is still shell. This slice ports the remaining hook logic behind binary subcommands, collapses the wrapper-search duplication that slice 4 deliberately inlined for one slice, and removes the kit's last runtime `python3` calls.

## Solution

The verdict logic of the three remaining hooks moves into the Go core behind plumbing subcommands, mirroring how `gitguard` absorbed the analyzer; the `.sh` files stay (the harness invokes hooks by path) but shrink to thin shims owning only `--describe`, wrapper resolution, and their fail-posture rim.

Two new packages carry the ported logic. `internal/lines` owns the lines.env binding parse — the single Go source for what a tier value *is*, replacing `lines-env.sh` — plus the agent-line verdict and the adapter model resolution, so both binding consumers share one parse. `internal/stophook` owns the stop verdict: the `stop_hook_active` envelope parse, the `BENCH_SHIFT` arming check, gate orchestration, the tail-30 output truncation, the verdict-cache write (reusing `git.TreeHash` in-process, never forging a verdict), and the `BLOCKED` message.

`cmd/bench` gains three subcommands (names are this spec's call): `resolve-model` (prints the resolved model to stdout, exit 1 in a routed repo with an unbound line — fits the `commands` map), and two direct cases like `guard-git`/`version` — `check-agent-line` (reads the Agent envelope on stdin, exit 2 DENIED with the message on stderr, plus a `--describe-binding` sub-mode that emits the live binding to stdout without reading stdin) and `stop-verdict` (reads the Stop envelope on stdin, orchestrates the gate, exit 2 BLOCKED with the gate tail on stderr, writes the cache).

The wrapper-search duplication collapses into a new shared `.bench/lib/resolve-bench.sh` carrying only the search order (repo `.bench/bin/bench.sh` → kit `bin/bench.sh` → `bench` on PATH), POSIX-clean because the adapters source it POSIX-style. Each shim source-guards it under its own fail posture: a missing lib is treated exactly like a missing core (the git/agent-line guards refuse or warn-and-allow; `stop` and `session-start` warn and fail open). `session-start.sh` swaps its inlined `bench_cmd()` for the shared resolver and is otherwise unchanged.

`_line-guard.sh` and `lines-env.sh` are deleted; the adapters exec `bench resolve-model` directly. With the shell parser gone, the `gate-line-contracts.sh` shared-parser pin retires and its hostile-input property moves to Go parser tables; the `lines-env-broken` and `alias-line-broken` canaries retire with it; and the `internal/guards` aggregator drops its hardcoded `_line-guard.sh` reference, so the ShiftAdapter row leaves `bench guards` (the enforcement persists in `bench resolve-model`; the advertisement does not). Every enforcement assertion is preserved; only the mechanism and the trigger move from shell/Python to the binary.

## User stories

1. As the two binding-consuming guards, I want the lines.env binding parse ported to `internal/lines` — reading a tier or alias value by key with last-assignment-wins, surrounding double/single-quote stripping, trailing-CR and surrounding-whitespace trimming, an indented key accepted, a missing key read as empty, and a value on a final line with no trailing newline read intact — so that what a tier value *is* has one Go source and `lines-env.sh` can be deleted. Line: claude-sonnet-5 / medium. This is a mechanical port of a small, sharply-specified parser whose every hostile-input edge is pinned by a Go table test, so the gate fully grades it and the cheap tier fits; medium effort because the strip/trim ordering is load-bearing and a silent divergence changes an allow/deny.

2. As the delegation guard, I want the agent-line verdict ported to `bench check-agent-line` and the hook reduced to a thin shim — the binary parses the Agent envelope (`tool_input.resolvedModel` then `tool_input.model`), reads the binding through `internal/lines`, exits 2 with the `DENIED:` message naming the binding only for a present model matching no bound tier or declared alias, and warns-and-allows (exit 0) on every degraded branch (unparseable stdin, no model field, no lines.env, incomplete binding); the shim resolves the wrapper via `resolve-bench.sh`, pipes the envelope through, passes the exit code back, and answers `--describe` by composing name/boundary/why with the live denies clause from `bench check-agent-line --describe-binding` — so that the deep unit owns the verdict and the shim carries no binding logic. Line: claude-opus-4-8 / high. The warn-and-allow-vs-deny posture is the whole guarantee — a fail-open branch that flips to deny bricks delegation, a deny branch that flips to allow silently de-enforces invariant #2 — and the gate cannot grade the message wording, so it takes the mid tier at high effort.

3. As the headless-shift guard, I want the adapter model resolution ported to `bench resolve-model` and the three adapters swapped from sourcing `_line-guard.sh` to exec'ing the subcommand — the binary prints the resolved model to stdout (empty for passthrough) and exits 1 in a routed repo when `BENCH_MODEL` is unset or unbound, applying the same rules (unrouted passthrough with explicit-beats-absent, incomplete-binding warn-and-passthrough, routed-repo requires a bound `BENCH_MODEL`); the adapters stay POSIX-clean, resolve the wrapper via `resolve-bench.sh`, and fail closed (refuse to launch) when the wrapper or core is unreachable in a routed context; `_line-guard.sh` and `lines-env.sh` are deleted — so that both binding consumers now read one Go parse and the adapter carries no sourced logic. Line: claude-opus-4-8 / high. The exec swap must stay POSIX-clean (the adapters are sourced POSIX-style) and the fail-closed rim is safety-critical — an unguarded passthrough in a routed repo is silent de-enforcement the gate can only partly reach — so it takes the mid tier at high effort.

4. As the completion oracle, I want the whole stop verdict path ported to `bench stop-verdict` (in `internal/stophook`) with `stop.sh` reduced to a thin shim — the binary reads the Stop envelope, honors `stop_hook_active` (allow the stop), enforces only when `BENCH_SHIFT=1`, runs the gate as a subprocess through the shim-passed wrapper, truncates the captured output to the last 30 lines for the `BLOCKED` message, computes the tree hash in-process via `git.TreeHash` and writes `<status> <tree-hash> <iso8601>` to the git-dir cache (writing nothing when the hash is unavailable, never a forged verdict), and exits 0 allow / 2 block; the shim resolves the wrapper, passes it to the binary, and keeps `--describe` plus the binary-missing rim (fail open with the loud no-verdict warning, no cache write) — so that the verdict fact lives in one language behind one subcommand. Line: claude-opus-4-8 / high. The fail-open-on-missing-core rim, the no-forged-verdict cache guarantee, and the `BLOCKED` message wording are the safety surface the gate cannot fully grade, and a stop that blocks a green gate or allows a red one is the worst class of oracle bug, so it takes the mid tier at high effort.

5. As every hook shim, I want the wrapper search collapsed into a new POSIX-clean `.bench/lib/resolve-bench.sh` carrying only the search order (repo `.bench/bin/bench.sh` → kit `bin/bench.sh` → `bench` on PATH, no policy), source-guarded by each shim under its own fail posture (lib-missing treated like core-missing: `stop`/`session-start` warn and fail open, `check-agent-line` warns and allows), with `session-start.sh` swapping its inlined `bench_cmd()` for the shared resolver and otherwise unchanged — so that the search order has one source and each shim's fail posture stays its own explicit fact. Line: claude-opus-4-8 / medium. The shared lib buys a new fail-open path if blindly sourced, so the per-shim source-guard posture and POSIX cleanliness (the adapters source it) are the correctness surface the gate can only partly reach, taking the mid tier at medium effort.

6. As the gate, I want the line-routing net re-pointed off the deleted shell parser with no enforcement assertion weakened, and the guards surface reconciled with the deleted `_line-guard.sh` — `gate-line-contracts.sh` drops its `. lib/lines-env.sh` source and its shared-parser pin (`a0`), its binding-shape read (`a`) re-points to `bench check-agent-line --describe-binding`, its hook cases (`c`) re-point the lib-missing sub-case to the shim's `resolve-bench.sh`-missing posture, and its adapter cases (`d`) re-point from sourced-function to exec shape; the `lines-env-broken` and `alias-line-broken` canaries retire with the shell parser; `internal/guards` drops the hardcoded `_line-guard.sh` reference so the AXI guards aggregation contract goes count-5 → count-4 (dropping the `_line-guard` row), `guards --brief` goes six lines → five, and the `guards-aggregation-dropped` and `adapter-line-broken` canary fixtures update to the new shapes — so that the "broken binding degrades loudly / cannot classify → refuse" properties stay gate-locked while their triggers move to the binary. Line: claude-opus-4-8 / high. Touching gate contracts and the canary layer is the worst defect class in this kit (`craft-gate`) — re-pointing triggers without weakening an assertion, and changing an observable AXI count without a rot, takes the mid tier at high effort.

## Implementation decisions

- **Package layout.** New `internal/lines`: the lines.env binding parser (the one Go source for a tier/alias value, replacing `lines-env.sh`), the Agent-envelope parse (`tool_input.resolvedModel` → `model`), the agent-line verdict (bound tier/alias match vs the degraded branches), the resolve-model verdict (unrouted / incomplete / routed-bound), and the `--describe-binding` emitter. New `internal/stophook`: the Stop-envelope parse (`stop_hook_active`), the `BENCH_SHIFT` arming and verdict orchestration, gate subprocess orchestration through the passed wrapper, the tail-30 truncation, the cache write (reusing `internal/git.TreeHash`, not a second hash), and the `BLOCKED` message. The per-hook split (rather than one `internal/hooks`) mirrors `gitguard` getting its own package; the binding parse lives in `internal/lines` specifically so the two binding consumers share one parse — that shared parse is the point of story 1.

- **Dispatch shape.** `resolve-model` is a `commands`-map entry: it writes the resolved model to stdout and returns an exit code, which the map's `func([]string) (string, int)` signature expresses. `check-agent-line` and `stop-verdict` are direct cases in `run()` like `guard-git` and `version`: they read stdin, write their verdict to stderr, and use distinct exits (and `stop-verdict` writes a cache file), none of which the map signature can hold. `--describe-binding` is a sub-mode of the `check-agent-line` case (no stdin, stdout only), exactly as `--describe-classes` is for `guard-git`.

- **The stop wrapper contract — gate resolution stays in one place.** The stop shim resolves the wrapper via `resolve-bench.sh` and passes it to `bench stop-verdict`, and the binary execs `<wrapper> gate` to run the gate. Gate resolution (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect) stays in `bin/bench.sh`, the profile's one source; the binary orchestrates the gate but never re-derives how to find it. The tree hash is computed in-process via `git.TreeHash` (already the one source shared with `gate_record`), not through a `tree-hash` subprocess.

- **The stop verdict/rim table** (the load-bearing mapping, stated as a table because prose blurs it):

  | actor | condition | action |
  |---|---|---|
  | `bench stop-verdict` | `stop_hook_active` true, or `BENCH_SHIFT` ≠ 1 | exit 0, no cache |
  | `bench stop-verdict` | armed, gate green | exit 0, write `green <tree> <iso>` (or nothing if tree unavailable) |
  | `bench stop-verdict` | armed, gate red | `BLOCKED:` + gate tail-30 to **stderr**, exit 2, write `red <tree> <iso>` (or nothing) |
  | shim | wrapper unresolvable **or** core-missing rc (e.g. 127) | warn loudly on stderr, exit 0, no cache (fail open, no forged verdict) |
  | shim | `stop-verdict` rc 0 or 2 | pass through (message already on stderr) |

  The stop posture is fail **open** — the mirror of `guard-git`'s fail-closed — because a missing oracle must degrade to "no verdict," never a false green. This preserves the two observable guarantees the shell hook has today: `stop_hook_active` and unarmed stops always exit 0, and a missing core writes no cache.

- **The agent-line verdict/rim table:**

  | actor | condition | action |
  |---|---|---|
  | `bench check-agent-line` | model matches a bound tier or declared alias | exit 0 |
  | `bench check-agent-line` | present model matching no bound tier/alias | `DENIED:` naming the binding to **stderr**, exit 2 |
  | `bench check-agent-line` | unparseable stdin / no model field / no lines.env / incomplete binding | `WARNING:` to stderr, exit 0 (fail open) |
  | `bench check-agent-line --describe-binding` | (no stdin) | emit `top=… mid=… cheap=…` (+ aliases) or `unrouted` to stdout, exit 0 |
  | shim | wrapper/lib unreachable on the enforcement path | warn, exit 0 (fail open — a broken guard never bricks delegation) |
  | shim | `--describe` with wrapper/lib unreachable | print name/boundary/why with a degraded denies clause, exit 0 |

- **`resolve-bench.sh` is search-order only; the source-guard is per shim.** The lib carries the three-step search and nothing else, POSIX-clean so the adapters can source it. Each shim carries its own `[ -f resolve-bench.sh ] || <posture>` guard: `guard-git`/`check-agent-line` treat a missing lib like a missing core (refuse git-shaped / warn-and-allow), `stop`/`session-start` warn and fail open, and the adapters fail closed (refuse to launch unguarded in a routed repo). No shim blindly sources the lib — the failure mode slice 4 named (missing lib → shim errors before its rims → a non-2 PreToolUse exit silently grants) is why the guard is explicit.

- **`--describe` stays in the shims; no second answerer.** Each `.sh` shim answers `--describe` by composing its fixed name/boundary/why with a denies clause read from the binary's `--describe-*` sub-mode (the same table/binding enforcement reads). The ported Go logic must not grow a hooks-dir `.sh` that answers `--describe` a second time — `bench guards` aggregates `*.sh` manifests and a duplicate would collide, exactly as in the guard slice.

- **The ShiftAdapter guards row drops — a map-forced consequence, flagged for veto.** `internal/guards` today hardcodes `.bench/adapters/_line-guard.sh` as a fifth aggregated manifest. Deleting `_line-guard.sh` (map #1) with the map's rule that `bench guards` aggregates `*.sh` only and grows no second answerer (map #9) forces the ShiftAdapter row out of `bench guards`: the enforcement moves into `bench resolve-model`, but the advertisement has no `.sh` to live in. The aggregation contract's count-5 header, the `guards --brief` line count, and the `guards-aggregation-dropped` canary all update. This changes an observable AXI contract, so it is surfaced for the reviewer's veto; the alternative (keep a manifest by making an adapter answer `--describe`) is rejected as a second answerer the map forbids.

- **Item 7 (Codex envelope) — settled from the adapters, no escalation.** `.codex/hooks.json` wires only two hooks: `Stop` → `stop.sh` and `PreToolUse:Bash` → `block-dangerous-git.sh`. It wires neither `check-agent-line.sh` (`PreToolUse:Agent`) nor `session-start.sh` — those are Claude-Code-only, in `.claude/settings.json`. So the **agent-line** envelope shape is a Claude-Code-only concern (`tool_input.resolvedModel`/`model`); Codex never invokes it, so there is no cross-harness divergence to reconcile, and the Go parse reads the one Claude Code shape. For **stop**, both harnesses invoke the same shim; the only field read is `stop_hook_active`, whose absent/unparseable case is treated as not-active → enforce — fail-toward-the-oracle, which already serves Codex today (Codex is wired now and the shell hook handles an absent field this way). The Go port preserves that default. The uncertainty flag resolves by reading the adapters, as the Handoff allowed; no escalation.

- **Deletion is part of this change, not a follow-up.** `_line-guard.sh` and `lines-env.sh` are deleted; `gate-line-contracts.sh`'s `. lib/lines-env.sh` source and its `a0` shared-parser pin are removed; the `lines-env-broken` and `alias-line-broken` canaries are removed (their guarded shell parser is gone); `internal/guards`' `_line-guard.sh` reference is dropped and the `guards-aggregation-dropped`/`adapter-line-broken` canary fixtures update — all in the same change so the gate load, the canary meta-check, and the docs stale-reference sweep stay green. `bench guard-git` and `block-dangerous-git.sh` (slice 4) are unchanged.

## Testing decisions

- **What a good test is here.** Acceptance drives each hook end-to-end through its shim — pipe a JSON envelope (or invoke an adapter with a prompt and `BENCH_MODEL`) and assert exit code plus message — never Go internals. The existing shell contracts already do exactly this and are the parity net; they run with their assertions intact, re-pointed only where they read the deleted shell parser or a sourced function. Go table tests are additional, at the pure-function seam in `internal/lines` and `internal/stophook`, where the shell-untestability tax on the binding parse and the envelope handling finally retires.

- **Seams.** Three, mirroring the guard slice: the **hook enforcement seam** (the shell contracts driving each shim → binary — `gate-line-contracts.sh` for agent-line and the adapters, the `stop.sh` cases in `gate-runtime-contracts.sh` for stop, including the gate-subprocess sandbox where the binary execs `<wrapper> gate`); the **`go test` unit seam** (table tests beside `internal/lines` and `internal/stophook` — prior art: `internal/gitguard`, `internal/maps`); and the **degradation/manifest seam** (the shim fail-posture rims, the `--describe` manifests, and the `bench guards` aggregation — prior art: `gate-axi-contracts.sh` and `gate-axi-guards-contracts.sh` against the guard-git shim).

- **Gate command:** `bench gate` (the project gate, which already runs `go build`/`go vet`/`go test ./...` and the cross-compile matrix in its Go layer, so the new packages and tests are graded with no new wiring), plus the parent map's explicit per-slice done rule `go build ./... && go vet ./... && go test ./...`. Done per slice = gate green and those three green.

### Seam diagram — hook enforcement seam (shim → binary, incl. the gate subprocess)

    trigger: agent/harness event → hook .sh runs (Claude Stop/PreToolUse:Agent; Codex Stop; bench shift → adapter)
        │
        ▼
    Stop envelope   ──▶ [ stop.sh: source-guard resolve-bench.sh → resolve wrapper → pipe stdin + pass wrapper ]
      on stdin          [ bench stop-verdict → internal/stophook: stop_hook_active? BENCH_SHIFT? ]
      (+ BENCH_SHIFT)   [   exec <wrapper> gate ──▶ capture ──▶ tail-30 ──▶ git.TreeHash ──▶ cache ] ──▶ exit 0 (allow) / exit 2 + BLOCKED+tail on stderr
    Agent envelope  ──▶ [ check-agent-line.sh: resolve wrapper → pipe stdin ]
      on stdin          [ bench check-agent-line → internal/lines: parse model, read binding, verdict ] ──▶ exit 0 (allow/degraded) / exit 2 + DENIED on stderr
    prompt+BENCH_MODEL ▶ [ adapter: resolve wrapper → model="$(bench resolve-model)" || refuse ] ──▶ exec harness (model flag) / exit 1 refuse
        ◀ tests attach here: gate-runtime-contracts.sh pipes {"stop_hook_active":…}/{} with BENCH_SHIFT=1 and asserts
          exit 0/2 + the cache line; gate-line-contracts.sh (c) pipes Agent envelopes and (d) invokes each adapter with a
          stub harness on PATH, asserting allow/deny/passthrough. Assertions unchanged; reads re-point off the deleted parser.

### Seam diagram — `go test` unit seam (binding parse, verdicts, envelopes)

    trigger: gate Go layer runs `go test ./...`
        │
        ▼
    lines.env bytes / ──▶ [ internal/lines.<parse tier value>        ] ──▶ value or ""
    key                   [ internal/lines.<agent-line verdict>      ] ──▶ allow / deny / degraded
    Agent envelope bytes  [ internal/lines.<resolve-model verdict>   ] ──▶ model / "" / exit-1 signal
    Stop envelope bytes   [ internal/stophook.<stop_hook_active>     ] ──▶ active? bool
    gate output string    [ internal/stophook.<tail-30 truncation>   ] ──▶ last ≤30 lines
        ◀ tests attach here: table tests pin the parser's hostile-input edges (quoting, whitespace, CRLF, indent,
          last-wins, missing key, no trailing newline, empty vs absent), each verdict branch, the envelope edges
          (non-JSON, missing fields, stop_hook_active true/false/absent/"True"), and the tail-30 boundary. Red before
          the package exists → does not compile.

### Seam diagram — degradation / manifest seam (shim rims + --describe + guards)

    trigger: `bench guards` runs `bash <hook>.sh --describe`; or a hook fires with the wrapper/core/lib gone
        │
        ▼
    --describe (no stdin) ──▶ [ shim: name/boundary/why + denies from `bench <sub> --describe-binding/-classes` ]
                              [   wrapper/lib unreachable → denies degrades (unrouted / manifest unavailable) ] ──▶ exit 0
    stop, core gone        ─▶ [ stop.sh rim: warn loudly, exit 0, NO cache ]                                    ──▶ exit 0 (no verdict)
    agent-line, lib gone   ─▶ [ check-agent-line.sh rim: warn, exit 0 ]                                          ──▶ exit 0 (fail open)
    adapter, core gone     ─▶ [ adapter rim: refuse to launch in a routed repo ]                                 ──▶ exit 1 (fail closed)
        ◀ tests attach here: gate-runtime-contracts.sh (stop missing-core / missing-bench: no cache, exit 0),
          gate-line-contracts.sh (c: lib-missing fail-open; d: guard/lib-missing fail-closed), gate-axi-guards-contracts.sh
          (aggregation count-4, each surviving --describe row). Re-pointed from the shell/Python triggers to the new ones.

### Acceptance coverage map

Per-item granularity is stated where a behavior quantifies over a set (each parser edge, each verdict branch, each degradation rim).

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | binding parse edges — double/single-quote strip, trailing-CR + surrounding-whitespace trim, indented key, last-assignment-wins, missing key → empty, no-trailing-newline value, empty-value vs absent-key — each edge | go test unit | `internal/lines` parser table test before the parse fn exists → does not compile | inherits the property the retiring `a0` pin held; a mis-ported strip/trim silently changes a tier value and flips an allow/deny |
| 2 | a bound tier id or declared alias allows (exit 0); a present model matching no tier/alias denies (exit 2) with the binding named | hook enforcement + go test unit | already covered (parity net): `gate-line-contracts.sh` (c) asserts allow-bound / allow-alias / deny-undeclared-alias / deny-unbound today and stays green across the port; Go verdict table before it exists → no compile | any drift in the tier/alias match flips one of the four asserted verdicts |
| 2 | each degraded branch — unparseable stdin, no model field, no lines.env, incomplete binding — warns and allows (exit 0), per branch | hook enforcement + go test unit | already covered (parity net): (c) asserts each fail-open branch + its stderr warning today; the lib-missing sub-case re-points to `resolve-bench.sh` missing → red until the shim's source-guard exists | catches a degraded branch that flips to deny (bricks delegation) or drops its warning |
| 2 | `--describe` denies clause carries the live binding (or `unrouted`) from `bench check-agent-line --describe-binding` | degradation/manifest | `gate-axi-guards-contracts.sh` asserts the `check-agent-line` row's denies; re-pointed to the binary sub-mode → red until it emits the binding | one source for enforce and advertise; a divergence between the verdict's binding and the manifest's surfaces here |
| 3 | routed+bound `BENCH_MODEL` → adapter passes model+prompt (incl. a dash-leading prompt); routed+unset/unbound → refuse (exit 1); unrouted → passthrough; unrouted+explicit → passthrough; incomplete binding → warn+passthrough | hook enforcement + go test unit | `gate-line-contracts.sh` (d) asserts all six today; re-pointed from sourced-function to `bench resolve-model` exec shape → red until the adapters swap; Go resolve-model verdict table → no compile | catches an adapter that launches on an unbound model (silent de-enforcement) or refuses a legitimately unrouted run |
| 3 | adapter fails closed (exit 1) when the wrapper or core is unreachable in a routed repo | degradation/manifest | the `adapter-line-broken` canary, re-pointed: a fixture adapter that reaches no core must still refuse in a routed repo → red until the fail-closed rim exists | catches a shim that grants a headless shift when it cannot resolve the binding |
| 3 | `_line-guard.sh` and `lines-env.sh` deleted; no dangling reference in contracts, `internal/guards`, README, or link/package | gate load + docs stale-reference sweep | the sweep against a tree still naming either file → red; a gate load sourcing the deleted `lib/lines-env.sh` → red | the conformance layer fails when a deleted file is still referenced or sourced |
| 4 | armed (BENCH_SHIFT=1) red gate → exit 2 + `BLOCKED:` + gate tail; armed green → exit 0; `stop_hook_active` true → exit 0; unarmed → exit 0; cache line `<status> <tree> <iso8601>` | hook enforcement | already covered (parity net): the `stop_hook_active`, gate-cache-write, and verdict-record contracts in `gate-runtime-contracts.sh` run unchanged and stay green across the port | any regression in arming, the flag, or the cache shape trips one of the unchanged contracts |
| 4 | tree-hash unavailable / core binary missing → write NO cache, fail open exit 0 with a loud warning | hook enforcement | already covered (parity net): the "stop hook missing-core-binary fail-safe" and "missing-bench fail-open" contracts assert no cache and exit 0 today and run unchanged | catches a port that forges a green keyed to a guessed tree, or traps the stop when the core is gone |
| 4 | gate output > 30 lines → only the last 30 appear in the `BLOCKED` message | go test unit | `internal/stophook` tail-30 table test before the truncation fn exists → no compile | pins the boundary the black-box contracts don't exhaust; a mis-ported truncation floods or empties the agent's feedback |
| 4 | `stop_hook_active` variants (true/false/absent/`"True"`/non-bool) and non-JSON stdin resolve to active?→bool with absent→not-active | go test unit | `internal/stophook` envelope table test before it exists → no compile | preserves fail-toward-enforcement on a malformed Stop envelope; also the Codex-shares-stop resolution (absent field → enforce) |
| 5 | wrapper search order (repo `.bench/bin/bench.sh` → kit `bin/bench.sh` → `bench` on PATH) resolves, and finds nothing → each shim's posture fires | hook enforcement + degradation/manifest | already covered (parity net): the stop missing-bench fail-open contract exercises "search finds nothing → fail open"; the gate-cache contract exercises "search finds the wrapper" | catches a broken search order (wrong precedence or a missed candidate) via the shims that depend on it |
| 5 | each shim source-guards `resolve-bench.sh` under its own posture (lib-missing = core-missing) | degradation/manifest | `gate-line-contracts.sh` (c) lib-missing sub-case re-points to `resolve-bench.sh` missing → fail open; (d) re-points to fail closed → red until the source-guards exist | catches a blind source that turns a fail-closed guard fail-open (the slice-4 failure mode) |
| 6 | `gate-line-contracts.sh` re-pointed off the deleted parser (`a0` retired, `a` reads via the binary, `c`/`d` re-pointed) with no assertion weakened; `lines-env-broken`/`alias-line-broken` canaries retired | gate load + canary meta | after deleting `lines-env.sh`, the contract's `. lib/lines-env.sh` source → red until removed; an orphaned canary the registry still names → canary meta red | proves the "broken binding degrades loudly" property stays gate-locked while its trigger moves from the shell parser to the binary |
| 6 | `bench guards` aggregation count-5 → count-4 (`_line-guard` row dropped); `guards --brief` six → five lines; `guards-aggregation-dropped`/`adapter-line-broken` fixtures updated | degradation/manifest + canary meta | `internal/guards` still referencing the deleted `_line-guard.sh` → the aggregation contract's count-5 header + `_line-guard` row assertion go red until updated | catches a guards surface that advertises a manifest for a deleted file or miscounts the deny rows |

### Edge inventory

Walked per behavior against the profile's shell-CLI hostile-input checklist and the map's item-6 owners; each resolved as a coverage row above or a **Won't handle** line here.

- **hand-edited lines.env — quotes, spaces, CRLF, no trailing newline** — covered: story 1's parser table tests (each edge), inheriting the retiring `a0` pin's property.
- **absent lines.env vs present-but-empty-tier** — covered: distinct behaviors both asserted — absent → unrouted/allow (agent-line) or passthrough (adapter); present-with-an-empty-tier → incomplete-binding warn (stories 2 and 3, contracts + Go tables).
- **non-JSON stdin / absent envelope fields** — covered: stop's `stop_hook_active`-absent → enforce and agent-line's no-model → fail open (Go envelope tables + parity contracts).
- **stop_hook_active boundary variants** — covered: the envelope table row (true/false/absent/`"True"`/non-bool).
- **oversized gate output** — covered: story 4's tail-30 boundary table test.
- **required tool missing from PATH (no wrapper, no core binary)** — covered: stop fail-open/no-cache (parity contracts), agent-line fail-open, adapter fail-closed (re-pointed contracts).
- **unquoted multi-word / dash-leading shift prompt** — covered: the adapter dash-leading `--` probe in (d), re-pointed to the exec shape.
- **cwd deeper than repo root** — covered by construction: the binary resolves the root via `git rev-parse`, and the gate runs from the root via `bin/bench.sh`, as today.
- **invocation through the `.claude/hooks` symlink** — safe by construction: each shim resolves `resolve-bench.sh` relative to itself (`pwd -P` / `BASH_SOURCE`), preserving today's symlink-survival; no new test (the search finds the repo-local wrapper regardless of the symlinked entry point).
- **re-run idempotency** — safe by construction: the hooks are stateless; re-running stop re-runs the gate and rewrites the cache for the same tree (identical result); no test needed.
- **binary panic mid-verdict** — covered: the `check-agent-line`/`stop-verdict` direct cases recover a panic to a non-verdict exit (as `guard-git` maps panic→3), so a crash routes to the shim's fail-posture rim (allow / no-verdict) rather than masquerading as an intentional exit-2 block.
- **Won't handle: real end-to-end harness wiring** (Claude Code / Codex actually firing the hooks) — stays manual smoke, gate-blind, as today (map Handoff #5); the contracts drive the shims directly.
- **Won't handle: Codex invoking the agent-line or session-start hook** — `.codex/hooks.json` wires neither, so there is no in-scope Codex caller and no cross-harness envelope divergence for those two hooks; Codex shares only the stop and git-guard shims.
- **Won't handle: the gate subprocess exceeding the harness 30s timeout** — the harness kills the hook, which forges no verdict (fail safe); the timeout is not made configurable here (parity with today).
- **Won't handle: SIGINT mid-cache-write leaving a partial line** — the cache line is a single redirected `printf`, and `bench status` validates the line shape before trusting it, so a partial write reads as no verdict rather than a forged one; not separately tested (parity).
- **Won't handle: a wrapper or repo path containing a literal newline** — parity with today; no in-scope caller produces one, and Go exec passes argv (never a shell string), so spaces and globs in the path are already safe.

## Out of scope

- **Slice 6 — `doctor` + `link` ported to Go** (the highest-stakes mutators). A distinct capability with its own spec by the parent map's slice order; this slice ports only the interactive hooks. Estimate to build later: one spec-sized session (~20–30 edits, ~4–8 gate runs), per the parent map's ~10–15-session budget.

- **Slice 7 — worktree + shift loop ported to Go** (needs its flagged contract backfill first). A separate slice with its own spec and its own uncertainty flag in the parent map. Estimate: one spec-sized session plus the backfill (~25–35 edits, ~5–10 gate runs).

- **Slice 8 — gate fragments → `go test`** (the gate ports last, under the canary layer's watch). A separate slice with its own spec. Estimate: one-plus spec-sized session (~30–40 edits, ~8–12 gate runs).

- **Collapsing the hook shims into pure in-process guards (no `.sh`)** — not possible, not a deferrable cut: the harness invokes a hook by file path (`Stop`/`PreToolUse:Agent` → a script), so a thin shim is structural. The same holds for the shift adapters, which `bench shift` execs by path.
