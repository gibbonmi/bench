# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2bcb4c8` at write time; this pass commits on top of it
Spec: `specs/ft128-agent-line-binding/spec.md` (Status: staged, reviewer-approved 2026-07-31)
Gate: green — full run, 2026-07-31

## State

- **FT128's spec was reshaped from its compiled decision map and is staged for sign-off.**
  It is uncommitted. `bench coverage --check` passes at 37 rows.
- **The previous session's blocker was half wrong, and that changed the spec.**
  `tool_input.subagent_type` *is* a documented Agent input and was captured from a real
  `fork` delegation on 2026-07-31, so the discriminator exists. What does not exist is the
  other half: per the Claude Code hooks reference, only `SessionStart` can receive a `model`
  field and there is no `$CLAUDE_MODEL`, so a `PreToolUse` hook cannot read the session's own
  tier and the map's "declared versus inherited tier" comparison is unimplementable at this
  seam. Story 4 instead denies a fork that declares any model (the harness ignores it) and
  allows one that declares none. **That substitution is the open reviewer decision.**
- A second map premise is false: `check-agent-line --describe` does not exist, so decision
  #7's session-start line report describes no surface. No story re-adds it; the real FT97
  bite is the deny message, which story 5 owns.
- A falsification pass ran on Codex (`gpt-5.6-sol`, high, read-only) and returned
  NEEDS-REWORK. Its findings were applied: per-cell and per-tier enumeration instead of
  sampled rows, matrix-derived denial rows that kill a hard-coded renderer, a no-dual-read
  row, doctor's unbound-column advice, honest TDD-able classifications, and two new rows for
  the hook and adapter missing-core rims — which nothing tests today, contrary to what the
  earlier draft claimed.
- **The four tickets under `specs/ft128-agent-line-binding/tickets/` are committed but
  stale.** They encode the old frontier in which every ticket blocks on capturing fork
  evidence, and that premise is falsified — the discriminator exists. Re-derive them from
  the reshaped spec before building; do not work the dependency order they record.
- A separate defect was found and fixed in this pass: the release scripts called `rg`
  unconditionally, which failed wherever ripgrep is not installed. `scripts/lib/search.sh`
  now selects the tool and falls back to POSIX `grep`. Unrelated to FT128.

## Next command

Sign off on the spec (or veto the story 4 substitution), then:

`/bench-implement-spec specs/ft128-agent-line-binding/spec.md`

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
