# Rewrite comments: internal/conformance files a to l, and the registry

Blocked by: none
Writes: internal/conformance/
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

Every explanatory `//` comment in the `internal/conformance/` files whose basename starts with `a` through `l`, and in `internal/conformance/registry/`, reads in ASD-STE100 inside the `craft-comments` register. A comment that restates its line is deleted. A provenance tag on an edited line is removed. A Go doc comment keeps its symbol-name opening. Directive comments stay byte-identical.

The delegate changes comment lines only; no code line moves. The orchestrator verifies that every changed line starts with `//` after whitespace and reads each rewritten comment against its code.

## Acceptance

- [ ] `gofmt`, `vet`, `test`, and `race` stay green, and the files' comments read in STE (covers PD33).
- [ ] Every changed line in the batch starts with `//` after whitespace (covers PD35).
