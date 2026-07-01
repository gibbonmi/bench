# Decision map — Codex hook layer

The whole-repo review suspected `.codex/hooks.json` was inert (Claude Code's
schema at a Codex path). Research disproved that for current Codex: the file is
plausibly functional as shipped. Remaining decisions are about verification and
honest documentation, not build-vs-drop. The harness-independent backstops (git
pre-push hook, `bench shift` loop) are unaffected throughout.

## #1: What hook/automation surface does current Codex actually support?

Type: Research

### Answer
Current Codex CLI (hooks shipped in 2026) supports lifecycle hooks that are
schema-compatible with Claude Code's: it reads `<repo>/.codex/hooks.json` (and
inline `[hooks]` in `.codex/config.toml`), with top-level
`{"hooks": {"<Event>": [{"matcher", "hooks": [{"type": "command", "command",
"timeout", "statusMessage"}]}]}}`. Supported events include `PreToolUse`
(matcher `"Bash"` for shell; stdin carries `tool_name` and
`tool_input.command`) and `Stop` (stdin carries `stop_hook_active` and
`turn_id`). Exit 2 + stderr blocks, same as Claude Code. The kit's shipped
`.codex/hooks.json` conforms to this schema field-for-field, and the shared
hook scripts consume exactly these stdin fields (stop.sh honors
`stop_hook_active`; the guard reads `tool_input.command`).

Caveats that drive #2/#3:
- **Trust gating** — project-level hooks load only when the project `.codex/`
  layer is trusted, and each non-managed command hook must be reviewed and
  trusted via `/hooks` (hash-recorded). After `bench link`, a Codex user gets
  the layer only after that one-time trust step.
- **Unverified in practice** — the `command` strings use
  `$(git rev-parse --show-toplevel)` substitution; whether Codex executes hook
  commands through a shell (making that resolve) is not stated in the docs.
  Needs one live dogfood run.
- **Version floor** — hooks are a 2026 Codex feature; older Codex versions
  ignore the file silently. Docs reference: developers.openai.com/codex/hooks.

## #2: How much verification does the layer get before the claim is trusted?

Type: Grill

### Answer
Dogfood plus a gate contract. A live `codex exec` run (codex-cli 0.142.5,
sandboxed) confirmed both hooks fire: PreToolUse on the Bash tool with
`tool_input.command`, and Stop with `stop_hook_active`, and the hook `command`
resolved through a shell so `$(git rev-parse --show-toplevel)` expanded to the
real script path. So the adapter is functional as shipped — no build needed.
Added gate check 2b (`.bench/gate.sh`): `.codex/hooks.json` must run
`stop.sh` on Stop and `block-dangerous-git.sh` under a `PreToolUse` `"Bash"`
matcher, with a `codex-hooks-broken` canary proving the check bites. Drift now
turns the gate red.

## #3: What do BENCH.md and the adapter docs say about the Codex layer?

Blocked by: #2
Type: Grill

### Answer
BENCH.md's "Hook Layers" section now states the two facts the old claim
omitted: project hooks require a one-time `/hooks` trust step in Codex after
`bench link`, and hooks are a 2026-era Codex feature — an older Codex ignores
`.codex/hooks.json` silently and keeps only the harness-independent backstops
(git pre-push hook + `bench shift` loop). The "adapters call the shared
scripts" claim stands, qualified by those two conditions.

---

All tickets resolved. The build (gate check 2b + canary) is done; the doc
edit is the remaining artifact. No spec needed — this was a research-plus-doc
change at an obvious seam.
