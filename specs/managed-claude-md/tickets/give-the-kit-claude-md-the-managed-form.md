# Give the kit CLAUDE.md the managed form

Blocked by: none
Writes: CLAUDE.md, tests/canary/docs-currency-token-diet/benchref-imported
Covers: none

## What to build

Claude Code reads the root `CLAUDE.md` of the kit at the start of each session. The file has nine lines, and most of them repeat the header of `AGENTS.md`. FT100 decision 8 in `decisions/guidance-weight/tickets/8.md` gives this file the managed form that `bench link` writes for a linked repository.

The root `CLAUDE.md` holds the exact bytes that `benchClaudeMD` in `internal/adopt/marker.go` returns. That form holds one heading, one pointer sentence, and the imports of `AGENTS.md` and `.bench/BENCH.md`. The two imports stay, so Claude Code loads the same guidance as before. The file does not import `.bench/BENCH-reference.md`.

The `Writes:` line names the fixture `tests/canary/docs-currency-token-diet/benchref-imported`, because its overlay pins `CLAUDE.md`. The ticket does not edit that fixture. No anchor registry names `CLAUDE.md`.

The ticket adds no check. A separate decision owns a check that holds the kit `CLAUDE.md` to the managed form. The rule "Do not edit `CLAUDE.md`" in `AGENTS.md` stays, because decision 8 is the reviewer's direction for this edit. The kit `CLAUDE.md` does not ship, so `CHANGELOG.md` gets no entry.

## Acceptance

- [ ] The root `CLAUDE.md` holds the same bytes that `benchClaudeMD` returns.
- [ ] Claude Code imports `AGENTS.md` and `.bench/BENCH.md` through the root `CLAUDE.md`, and the file does not import `.bench/BENCH-reference.md`.
- [ ] `bench doctor` in the kit reports that `CLAUDE.md` carries the Bench import lines.
- [ ] Each sentence of the root `CLAUDE.md` stays inside the ASD-STE100 sentence and paragraph bounds.
