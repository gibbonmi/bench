# Resolve the line through the harness matrix

Blocked by: none

## What to build

Tier becomes the only shared identity. `internal/lines` parses the closed
`BENCH_<HARNESS>_<TIER>` matrix over `codex`, `claude`, `opencode`, resolves a
tier for a named harness, and renders every allow, deny, and degraded verdict in
the asking harness's own tokens. Every runtime caller — the kit CLI, a linked
repo's CLI, the Claude hook, and the three adapters — passes its own
`--harness` to that one core. The `BENCH_TIER_*` / `BENCH_ALIAS_*` schema is cut
with no dual read, so this repo's `.bench/lines.env`, the profile's `Lines`
section, `checkLineBinding`, and the pinned fixtures migrate in this same green
change.

Covers stories 1, 2, 5, story 3's launch refusal and conformance silence, and
both edge-of-1 and edge-of-2 rows.

## Acceptance

- [ ] `--harness codex` and `--harness claude` each map `BENCH_MODEL` of `top`,
      `mid`, and `cheap` to that harness's own cell — all six.
- [ ] A `BENCH_<HARNESS>_<TIER>`-shaped key naming a harness outside the closed
      set is rejected rather than silently accepted.
- [ ] A `lines.env` carrying only retired `BENCH_TIER_*`/`BENCH_ALIAS_*` keys
      resolves as no binding at all, not as a legacy binding.
- [ ] An unreadable `lines.env` is distinguished from an absent one instead of
      failing open as unrouted.
- [ ] All six runtime surfaces pass their own `--harness` to the core, and each
      adapter's resolved model equals what the core returns for the same harness
      and tier rather than being recomputed in the shim.
- [ ] The retired `--alias` and `--provider-model` flags are rejected; an unknown
      `--harness` value and an unknown or unset tier token each exit 1.
- [ ] The OpenCode adapter refuses to launch while its column is unbound, and an
      absent OpenCode column produces no conformance diagnostic while Codex and
      Claude are complete.
- [ ] A Claude denial lists the Claude column as parsed from the binding, a Codex
      denial lists the Codex column and names no Claude token, and an unbound
      token denies with the active harness's recovery instruction.
- [ ] The Claude hook still fails open with its warning when the core cannot be
      resolved, and each of the three adapters refuses to launch with exit 1 in
      the same condition.
- [ ] Quoted, CRLF, last-wins, and no-final-newline matrix entries parse
      consistently, and an absent `lines.env` stays distinguishable from a
      present-but-empty one.
