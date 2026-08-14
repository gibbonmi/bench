# Learnings — usage journal

## 2026-08-13 — Keep the map-to-ticket flow in one session [open]

What happened: the suggested fresh-session boundary fell between the reviewed decision map and the spec-to-ticket flow, forcing the next session to reconstruct planning context before the tickets were sliced.

Right behavior: carry the reviewed map through spec authoring and reviewer-approved ticket slicing in the same session. Start the fresh implementation session only after the ticket graph is sliced and recorded.

Proposed rule change: move the suggested fresh-session handoff to the end of ticket slicing, with the continuation resuming at `/bench-implement-spec` against the approved ticket graph.

## 2026-08-14 — Discover landing verbs non-interactively [open]

What happened: mid-landing, `bench worktree land --help` was run bare; the wrapper opened an interactive worktree subshell instead of printing usage, briefly creating and releasing a stray assignment. AGENTS.md already forbids trying bare unknown verbs.

Right behavior: discover verbs via `bench commands --brief` or the grammar in source, and pipe stdin from `/dev/null` on any command that might prompt.

Proposed rule change: none — the existing rule suffices; this entry records the trip so the wrapper's subshell-on-unknown-flags behavior can be weighed by the parked `bench commands` staleness idea.

## 2026-08-14 — Follow the profile's build script for durable artifacts [open]

What happened: `dist/bench` was rebuilt with plain `go build` to drive the landing, violating the profile note that requires `scripts/go-build.sh` so the binary carries the package version; caught and rebuilt during final-check.

Right behavior: before touching a durable artifact, re-read the profile's cold-session notes; the sanctioned build path is `bash scripts/go-build.sh <root> <out>`.

Proposed rule change: none — the note existed; the miss was not reading it before acting.
