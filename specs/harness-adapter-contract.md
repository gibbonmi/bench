# Harness adapter contract

Source: `decisions/dogfood-improvements.md` ticket #4 (ratified 2026-07-01). Closes
the `BENCH_AGENT` multi-word/adapter thread the review-findings remediation spec
deferred here.

## Problem

`bench shift` invokes the agent as `"$AGENT" -p "$prompt"` with `AGENT` defaulting to
`claude` (bench.sh:24, :150, :178). Two problems follow. First, the `-p` flag is
Claude-specific: a Codex or OpenCode user cannot run the loop, yet the kit claims to
be harness-neutral. Second, the silent `claude` default means a non-Claude install
with no configuration gets a cryptic `claude: command not found` swallowed by the
loop's `|| true`, then a confusing gate run over an empty change — instead of a clear
"configure your harness adapter" error.

## Solution

`bench shift` calls a configured **adapter executable**, passing the generated prompt
as its single positional argument and arming `BENCH_SHIFT=1`. The adapter is a thin
wrapper that maps that one argument to its harness's noninteractive command. Bench
ships reference adapters for Claude Code, Codex, and OpenCode, but picks **no default**
— an unconfigured or unusable adapter fails fast with a legible error, never a silent
wrong-CLI call. The `-p` flag lives in the Claude adapter, not the loop.

## User stories

1. As a non-Claude user, I want `bench shift` with no `BENCH_AGENT` configured to exit
   nonzero with a clear "configure your harness adapter" message *before* any agent or
   gate runs, so I am told what to fix instead of watching a cryptic CLI failure get
   swallowed.
2. As any harness user, I want the prompt delivered to my adapter as its single
   positional argument (`$1`), not behind a Claude-specific `-p`, so my wrapper maps
   one argument to my harness's own flags.
3. As the harness operator, I want `BENCH_SHIFT=1` still set on the adapter call, so
   the Stop hook stays armed exactly as before (the arming must survive the contract
   change).
4. As a user who misconfigures `BENCH_AGENT` to a missing or non-executable path, I
   want a clear "adapter not executable" error before the loop, not a cryptic exec
   failure buried in iteration output.
5. As any user, I want a prompt containing spaces, newlines, or shell metacharacters
   to reach my adapter intact as one argument, so the multi-line iteration prompt is
   never word-split or re-tokenized.
6. As a Claude Code user, I want a shipped reference adapter that maps the prompt to
   `claude -p "$1"`, so I configure `BENCH_AGENT` at it and the loop works unchanged.
7. As a Codex user, I want a shipped reference adapter that maps the prompt to
   `codex exec "$1"`.
8. As an OpenCode user, I want a shipped reference adapter that maps the prompt to
   `opencode run "$1"`.
9. As the harness operator, I want each reference adapter to `exec` its harness so the
   harness's own exit code is the adapter's, keeping the adapter a transparent
   pass-through.
10. As both invocation paths (iteration and refactor), I want the same adapter
    contract, so neither loop site keeps a hardcoded `-p`.
11. As a reader of the operating guide, I want `.bench/BENCH.md` to state the adapter
    contract and how to configure `BENCH_AGENT`, and the in-code comment to stop
    claiming a `claude` default.

## Implementation decisions

- **No default in resolution.** `AGENT="${BENCH_AGENT:-claude}"` becomes an
  unconditional read of `BENCH_AGENT`; empty/unset is the unconfigured case.
- **Preflight, once, before the loop.** Before the first iteration, resolve the
  adapter: if `BENCH_AGENT` is empty → "configure your harness adapter" error, exit
  nonzero; if set but not an executable command (neither an executable file nor
  resolvable on `PATH`) → "adapter not executable" error, exit nonzero. Both exit
  before the worktree loop and before any gate run.
- **Single-arg invocation.** Both sites change from `"$AGENT" -p "$(prompt)"` to
  `"$AGENT" "$(prompt)"`. `"$AGENT"` stays quoted (one executable, never split), and
  the prompt stays a single quoted argument — this is what keeps a multi-line prompt
  intact and is why a configured value carrying its own arguments is unsupported
  (see Won't handle).
- **Reference adapters** live in `.bench/adapters/{claude,codex,opencode}`, each a
  one-line `exec <harness command> "$1"` wrapper, executable, shipped as part of
  `.bench/`. They are references to point `BENCH_AGENT` at, not an auto-selected
  default.
- **Recording-adapter migration.** Every contract recording adapter that captures the
  prompt reads it from `$2` today (because of `-p`); each moves to `$1`. This adapts
  the test harness to the new interface — it is not weakening an assertion.
- **Docs.** `.bench/BENCH.md` gains a short adapter-contract note under the shift
  material; the bench.sh comment on the `AGENT` line drops "defaults to claude."

## Testing decisions

- A good test runs `bash bin/bench.sh shift <objective>` in a throwaway repo with a
  recording or trivial adapter and asserts observable behavior: exit code, stderr
  message, and the argument the adapter received. Prior art: every `bench shift` block
  in `.bench/gate-runtime-contracts.sh`.
- Seam: the `bench shift` CLI contract (single seam). Reference-adapter content is
  asserted by grep against the shipped files in the same contract file.
- Gate: `bench gate`. Done = green with all new contract cases present.
- Line (declared): shell/contract + trivial adapter scripts + one doc note at **mid
  tier, medium effort** (binding in projects/benchkit.md: Sonnet 4.6; in-session
  delegates run Opus 4.8 per the alias caveat there). Combined cap ~150k tokens.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | unconfigured `BENCH_AGENT` → nonzero + legible error, no agent/gate run | shift contract | observed 2026-07-01: `env -u BENCH_AGENT ... shift probe` prints no "configure adapter"/`BENCH_AGENT` text | greps stderr for the guidance the silent default currently omits |
| 2 | prompt arrives as `$1` containing the objective, not `-p` | shift contract (recording adapter on `$1`) | observed 2026-07-01: recorder on `$1` gets `-p`, objective token absent | the recorded arg is the behavior; the row greps it |
| 3 | `BENCH_SHIFT=1` present on the adapter call | shift contract (recording adapter echoes `$BENCH_SHIFT`) | not TDD-able red-first (already armed today); asserted post-change so the contract change can't silently drop it | recorder asserts the env var each call |
| 4 | non-executable `BENCH_AGENT` path → nonzero + "not executable" error before loop | shift contract | red-first (run in implement): `BENCH_AGENT=/no/such/adapter shift x` currently exec-fails inside the loop, no preflight message | greps stderr for the preflight error absent today |
| 5 | multi-line/metachar prompt reaches adapter as one intact `$1` | shift contract (recording adapter dumps `$1`) | red-first: with `$1` holding `-p` today the assertion that `$1` equals the full prompt fails | asserts no word-split / single argument |
| 6–8 | each reference adapter maps prompt to its harness command | shift contract (grep shipped files) | red-first: `.bench/adapters/*` do not exist yet, grep fails | file-content assertion per adapter |
| 9 | reference adapters `exec` the harness (pass-through exit) | shift contract (grep shipped files) | red-first: files absent; post-build grep for `exec ` | asserts the transparent-exit shape |
| 10 | refactor-phase invocation uses the adapter contract too | shift contract (existing refactor-prompt-scope block, recorder moved to `$1`) | red-first: that block's recorder reads `$2`; after the loop change `$2` is empty | the refactor recorder captures the second site |
| 11 | BENCH.md states the contract; comment drops the default | review | prose — verified at review against the cited edits | — |

### Edge inventory

Walked per behavior:

- **Error path** — unconfigured (row 1), non-executable adapter (row 4). Covered.
- **Empty/absent input** — `BENCH_AGENT=""` is the unconfigured case (row 1). Covered.
- **Malformed input** — prompt with metacharacters/newlines (row 5). Covered.
- **Boundary** — the two distinct invocation sites, iteration and refactor (rows 2, 10). Covered.
- **Re-run idempotency** — adapter resolution is stateless per shift; no carryover. **Won't handle** as a row — nothing to assert.
- **Hostile environment** — `BENCH_AGENT` as a multi-word command string
  (`"claude -p"`): **Won't handle** — by decision #4, arguments belong in a wrapper
  script, not the config value; a multi-word value resolves as one non-existent
  executable and correctly hits the row-4 error. This is the deliberate closure of the
  remediation spec's deferred `BENCH_AGENT` multi-word thread.
- **Relative `BENCH_AGENT` path** — resolves against the worktree root where the
  adapter runs, not the caller's cwd. **Won't handle** as a row — documented in
  BENCH.md as "use an absolute path or an on-`PATH` name"; a one-line doc note, not a
  code branch.

## Out of scope

- **Config-file adapter binding** (adapter named in `projects/<name>.md` or a
  `.bench/config` instead of `BENCH_AGENT`) — a separate configuration-surface
  capability; the env var is the whole contract here. ~30 min if ever wanted.
- **`bench link` installing/selecting an adapter interactively** — a separate
  install-flow capability; the reference adapters ship as files under `.bench/` and
  the user points `BENCH_AGENT` by hand. ~45 min.
