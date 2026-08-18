# Add `bench status --route` with the harness flag and the clean-board fallback

Blocked by: 03-setup-and-staged-signals.md
Writes: internal/status, internal/handoff, cmd/bench

## What to build

Grammar `bench status [--all] [--route [--harness claude|codex]]`; `--route` and `--all`
exclusive, `--harness` requires `--route`, harness names from the shared table. Output:
`next[1]{state,why,command}` via the TOON table writer with `state`=signal name,
`why`=detail, `command`=action, then one `also:` line (runners-up in ladder order,
non-empty commands only, `none` otherwise). Empty board → `clean,nothing pending,bench
roadmap` (ROADMAP.md present) or `/bench-drain` (absent). Rows but no command → lead row
with empty command, no fallback. `bench handoff`'s Next consumes the same selection
including the fallback. Outside a repo → the structured not-in-repo error, exit 1.

Covers: R1, R2, R3, R4, R5, R6, R7, R8, R9, R11, R12, R13

## Acceptance

- [ ] Red-gate fixture prints exactly `next[1]{state,why,command}:` + `  gate,red,/bench-debug` + the `also:` line naming git and drain in order.
- [ ] Empty board with/without ROADMAP.md routes `bench roadmap` / `/bench-drain`; locked-pending-only board routes an empty command with `also: none`.
- [ ] `--harness codex` renders `$bench-` in the row and in `also:`.
- [ ] `bench handoff` on the empty-board fixture writes `bench roadmap` as Next; on the drain fixture `/bench-drain`.
- [ ] `--route --all`, `--harness codex`, `--route --harness opencode`, `--route extra` exit 2 with the grammar line; outside a repo exit 1.
- [ ] Contract test: `bench status --route` at this repo's HEAD emits one row.
- [ ] Gate green.
