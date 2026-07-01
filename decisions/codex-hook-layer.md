# Decision map — Codex hook layer

`.codex/hooks.json` ships Claude Code's hook schema (Stop/PreToolUse events,
matcher "Bash") at a Codex path, and nothing establishes that Codex reads it.
If it is inert, the destructive-git guard and stop-gate layers silently do not
exist under Codex, while `.bench/BENCH.md` claims both adapters call the shared
scripts. The harness-independent backstops (git pre-push hook, `bench shift`
loop) are unaffected either way.

## #1: What hook/automation surface does current Codex actually support?

Type: Research

### Question
Does the current stable Codex CLI read any hook configuration at all — and if
so, which events (session start, pre-tool, stop/turn-end), what schema, and
from which path? `.codex/hooks.json`'s shape is Claude Code's; we need the real
surface (or confirmation none exists) before choosing between building an
adapter and dropping the claim. Deliverable: a short summary of the supported
surface with version numbers and doc references.

### Answer
— (open)

## #2: Build a real Codex adapter, drop the layer, or partially map it?

Blocked by: #1
Type: Grill

### Question
If Codex exposes an equivalent surface, an adapter keeps the interactive safety
layer parity BENCH.md promises. If it exposes only part (e.g. lifecycle but not
pre-tool), a partial adapter needs honest documentation of the gap. If it
exposes nothing, the kit should drop `.codex/hooks.json` and lean on the
backstops. Reviewer call once #1 lands.

### Answer
— (open)

## #3: What happens to the shipped file and the BENCH.md claim?

Blocked by: #2
Type: Grill

### Question
`.codex/hooks.json` is in `package.json` files[] and installed by `bench link`;
BENCH.md's "Hook Layers" section claims Codex adapter coverage. Depending on
#2: replace the file with a real adapter, or remove it from the plan/manifest
and reword BENCH.md to name the git hook + shift loop as the only Codex-side
enforcement. Either way the docs must state the current decided truth.

### Answer
— (open)
