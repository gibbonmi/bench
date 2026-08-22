# Rewrite the Go-embedded Markdown templates and regenerate the handoff

Blocked by: 01c-register-the-check-and-exclusions.md
Writes: internal/handoff/, internal/adopt/, .bench/prose-exclusions (new in ticket 01c), capture/session-handoff.md
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The Markdown the kit writes into a tree reads in ASD-STE100 and passes both bounds. That Markdown is the handoff shape and scaffold texts in `internal/handoff/text.go`. It is also the AGENTS managed block, the generated `CLAUDE.md`, the learnings preamble, and the profile scaffold in `internal/adopt`. A unit test in each owning package renders each template and grades it through `internal/prose`. The shell scaffold strings belong to ticket 28. Tests that quote the Markdown strings update with them.

The ticket runs `bench handoff` so `capture/session-handoff.md` carries the new shape text in the same commit. The `capture/session-handoff.md` row leaves `.bench/prose-exclusions` in this commit. The orchestrator reads every template sentence.

## Acceptance

- [ ] Each rendered template passes `internal/prose` with no finding, and a planted long sentence in one template reds that test (covers PD42).
- [ ] The handoff shape constant and `capture/session-handoff.md` are byte-equal (covers PD40).
- [ ] The live-tree test passes with the handoff row removed (covers PD27).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
