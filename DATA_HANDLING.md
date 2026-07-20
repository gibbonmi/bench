# Data handling

This is an inventory of every repository-controlled path that operator- or
environment-supplied data reaches when Bench runs a shift, a gate, or a status
render. It describes the current decided state for an auditor with no prior
context: what data travels where, how durable it is, and who can read it.

Bench runs as the invoking user; see `SECURITY.md` for the trust model. Nothing
here is a sandbox. The controls below bound *which* data reaches
repository-controlled subprocess code and *how long* it persists — they are not a
boundary against a malicious process running as the same user.

**No content-based secret detection anywhere.** Sensitivity is handled by
documenting which paths are durable, never by inspecting values. The environment
passlist below filters variable *names*; it never reads a value. A passlisted
variable's value may still be sensitive — an API key under `ANTHROPIC_*`, a
credential under `AWS_*` — and that is a documented fact, not a defect. If a
value must not reach adapter code, do not export it, or narrow the passlist; Bench
will not guess at it.

## Prompt paths

The iteration prompt is the text a shift hands its harness on each cycle.

- **Loop to adapter — stdin, end to end.** The shift loop writes the full
  iteration prompt to the adapter process's standard input and passes no
  positional argument. The prompt therefore never appears in `argv` on any hop
  Bench itself controls, so it is absent from the machine's process listing for
  the loop-to-adapter call.
- **Adapter to harness CLI — stdin where the CLI documents it.** The Claude Code
  and Codex adapters forward stdin to their CLI's documented stdin path
  (`claude -p` and `codex exec` both read a piped prompt), so the prompt stays
  out of `argv` on the final hop too.
- **opencode adapter — residual positional final hop (documented harness
  limitation).** As of 2026-07-20 `opencode run` documents only a positional
  prompt argument. Its adapter reads the prompt from stdin and passes it to the
  CLI positionally after `--`. The prompt is therefore visible in that one final
  process's `argv` on the local machine. This is a known residual to drop when
  upstream documents a stdin path; it is asserted as a deliberate state, not left
  to drift.

## Environment

A shift launches its harness adapter from a documented passlist instead of
inheriting the parent environment verbatim. Only a variable whose *name* matches
a default pattern below, or a pattern added under `[agent]` in a committed
`.bench/env.allow`, reaches adapter code. Values are passed through byte for byte
and are never inspected, truncated, or split.

The project gate is a separate, already-closed subject: `bench gate` launches the
gate script with `PATH` plus only the names declared under `environment` in
`.bench/gate-inputs.json`, which is both the opt-in and part of the gate's verdict
identity. A repo whose gate script needs an extra variable declares it there. The
observable failure when it is missing is the gate turning red with the project's
own error naming the absent variable — not a silent widening.

### Default passlist

Every pattern below is a default the adapter admits. Exact names match exactly; a
`PREFIX*` token matches every name under that prefix. This listing is the single
documented source paired with the exported `internal/env` constants
(`SharedBasics` and `AgentPasslist`); a conformance check fails the gate if a
constant pattern is added without a row here, so the advertisement cannot drift
from the enforcement.

Provider keys for harnesses Bench ships no adapter for (Groq, Gemini, Azure,
Cloudflare) are deliberately absent — each is a one-line `.bench/env.allow`
addition, which is what the opt-in mechanism exists for.

<!-- passlist:begin -->

| Pattern | Family | Official documentation |
|---|---|---|
| `PATH` | Process basics | POSIX environment (IEEE Std 1003.1) |
| `HOME` | Process basics (load-bearing for git config resolution) | POSIX environment (IEEE Std 1003.1) |
| `USER` | Process basics | POSIX environment (IEEE Std 1003.1) |
| `LOGNAME` | Process basics | POSIX environment (IEEE Std 1003.1) |
| `SHELL` | Process basics | POSIX environment (IEEE Std 1003.1) |
| `TMPDIR` | Process basics | POSIX environment (IEEE Std 1003.1) |
| `TERM` | Process basics | POSIX environment (IEEE Std 1003.1) |
| `COLORTERM` | Process basics (terminal rendering) | POSIX-adjacent terminal convention |
| `LANG` | Process basics (locale) | POSIX locale environment (IEEE Std 1003.1) |
| `LC_*` | Process basics (locale; a glob because exact-name matching breaks real-system locales) | POSIX locale environment (IEEE Std 1003.1) |
| `XDG_*` | Process basics (`XDG_CONFIG_HOME` is load-bearing for git config resolution) | XDG Base Directory Specification |
| `BENCH_*` | Bench's own namespace | This repository |
| `ANTHROPIC_*` | Claude Code | https://docs.claude.com/en/docs/claude-code/settings |
| `CLAUDE_CODE_*` | Claude Code | https://docs.claude.com/en/docs/claude-code/settings |
| `CLAUDE_CONFIG_DIR` | Claude Code | https://docs.claude.com/en/docs/claude-code/settings |
| `API_TIMEOUT_MS` | Claude Code | https://docs.claude.com/en/docs/claude-code/settings |
| `CODEX_*` | Codex | https://github.com/openai/codex/blob/main/docs/config.md |
| `RUST_LOG` | Codex | https://github.com/openai/codex/blob/main/docs/config.md |
| `SSL_CERT_FILE` | Codex | https://github.com/openai/codex/blob/main/docs/config.md |
| `OPENCODE_*` | opencode | https://opencode.ai/docs/ |
| `OPENAI_API_KEY` | opencode (documented provider substitution) | https://opencode.ai/docs/ |
| `AWS_*` | Claude Code — Amazon Bedrock routing | https://docs.claude.com/en/docs/claude-code/amazon-bedrock |
| `GOOGLE_*` | Claude Code — Google Vertex routing | https://docs.claude.com/en/docs/claude-code/google-vertex-ai |
| `GCLOUD_PROJECT` | Claude Code — Google Vertex routing | https://docs.claude.com/en/docs/claude-code/google-vertex-ai |
| `CLOUD_ML_REGION` | Claude Code — Google Vertex routing | https://docs.claude.com/en/docs/claude-code/google-vertex-ai |
| `VERTEX_LOCATION` | Claude Code — Google Vertex routing | https://docs.claude.com/en/docs/claude-code/google-vertex-ai |

<!-- passlist:end -->

A default glob never straddles two families: every glob above covers a namespace
one owner controls. A name whose prefix is shared with a foreign family is
enumerated exactly instead — this is why there is no `GO*` glob, which would also
have matched `GOOGLE_APPLICATION_CREDENTIALS`.

### The opt-in: `.bench/env.allow`

A committed `.bench/env.allow` is the only sanctioned way to widen the adapter
passlist. It is optional; absent means defaults only, never an error. It is
line-oriented: `#` comments, blank lines, an `[agent]` section header (the only
known section), and one entry per line that is either an exact name or a
`PREFIX*` glob. A malformed or hostile file fails closed — the launch refuses and
the error names the offending line and reason — rather than degrading to
defaults, because a silently ignored opt-in is indistinguishable from a working
one. A leading UTF-8 byte-order mark is rejected the same way, named explicitly
rather than reported as a stray entry. A gate-side need is declared in
`.bench/gate-inputs.json` instead; `[gate]` is not a valid section here.

## File paths

- **`.bench-objective` (worktree scratch file, mode 0600).** The one place the
  full objective text persists. It is created readable only by the user who
  started the shift, lives inside the shift worktree, and dies with that
  worktree. `bench status` reads it back to show a live shift's objective and
  degrades to the intent key when the file or worktree is gone.
- **Intent ledger (`bench-intent.json`).** The durable ledger references an
  objective by its entry key; it carries no free-text objective field. An
  objective's text does not appear anywhere in this file.
- **Commit subjects.** A shift's commit subject carries the sanitized objective
  text by design. This is reviewer-authored content, and a readable git history
  is a feature, not a leak — the text is intended to persist in the repository's
  history.
- **`.bench/env.allow` and `.bench/gate-inputs.json`.** Committed, reviewable
  declarations of passlist and gate-subject widenings. They carry variable
  *names*, never values.

An objective is capped at 200 runes and rejected at intake — before any ledger
entry, scratch file, or commit — if it is over-long or carries a control byte,
so unbounded or control-bearing text cannot flow into durable state.

## Log and terminal output

Every terminal render of operator-influenced text goes through the single
sanitizer in `internal/sanitize`: control runes are escaped into a readable
`\n`/`\r`/`\t`/`\uXXXX` form so no raw control byte reaches the terminal. This is
the one control-sequence policy; a conformance check pins it as the sole escaper
so a second, drifting copy cannot appear.

TOON table output (`internal/toon`) is deliberately distinct: it *refuses* a
control-bearing cell rather than escaping it, a closed AXI-contract decision, and
is not folded into the sanitizer's escaping policy.

## Network

Bench makes outbound requests in three situations only (see `SECURITY.md`):

- **Model discovery.** `bench models` queries the OpenAI and Anthropic model-list
  APIs when the corresponding API keys are present.
- **Binary self-repair.** Reads package metadata and an artifact from the npm
  registry. Set `BENCH_NO_REPAIR` to any non-empty value to disable this egress.
- **Worktree acquisition.** A best-effort `git fetch origin` before creating a
  pooled worktree.

Ordinary git operations (add, commit, worktree management) talk only to the local
repository and whatever remote the repository is already configured to use. Bench
sends no telemetry and uploads no prompt, objective, or environment data anywhere.

## Cache

Bench keeps no cache of prompt, objective, or environment content. Binary
self-repair uses the npm cache (honoring `NPM_CONFIG_CACHE`) for the artifact it
downloads. Worktrees are pooled and reused, and a pooled worktree's
`.bench-objective` persists for the life of that worktree (see Retention). The git
object store holds commit subjects as ordinary repository history.

## Retention

- **Objective text** persists in exactly two places: the commit subject
  (durable repository history, by design) and `.bench-objective` (mode 0600,
  living only as long as its worktree). A retained or interrupted worktree keeps
  its `.bench-objective` alive at mode 0600 until the worktree is cleaned up —
  this is the one retention edge worth naming, and it stays user-readable-only
  throughout.
- **Intent ledger** retains objective *keys*, not text, for the life of the
  ledger.
- **Environment values** are never written to durable Bench state; they exist
  only in the launched subprocess's memory for the duration of that process.
