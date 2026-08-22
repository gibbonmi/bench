# Rewrite comments: small packages, group B, and the root Go files

Blocked by: none
Writes: internal/runbinary/, internal/canary/, internal/models/, internal/modelid/, internal/gittest/, internal/testrepo/, internal/harness/, internal/jsonfile/, internal/releasepreflight/, internal/racetests/, internal/terminal/, internal/sessioninspect/, internal/contract/, internal/capability/, internal/dashboard/, consumer_payload.go, consumer_payload_test.go
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

Every explanatory `//` comment in the fifteen packages, including `internal/modelid/modelidtest/`, and in the two root Go files reads in ASD-STE100 inside the `craft-comments` register. A comment that restates its line is deleted. A provenance tag on an edited line is removed. A Go doc comment keeps its symbol-name opening. Directive comments stay byte-identical.

The delegate changes comment lines only; no code line moves. The orchestrator verifies that every changed line starts with `//` after whitespace and reads each rewritten comment against its code.

## Acceptance

- [ ] `gofmt`, `vet`, `test`, and `race` stay green, and the packages' comments read in STE (covers PD33).
- [ ] Every changed line in the batch starts with `//` after whitespace (covers PD35).
