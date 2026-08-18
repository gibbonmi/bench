# Let a staged spec declare the phase commands it introduces

Blocked by: none
Writes: internal/conformance/docs_workflow_checks_test.go, internal/conformance/registry_test.go, tests/canary/docs-currency-token-diet, .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-spec/SKILL.md

## What to build

The stale-command sweep derives its valid `/bench-*` and `$bench-*` tokens from the
files present in `.agents/commands/`, so a staged spec whose deliverable is a new or
renamed phase is red on every mention of it. A staged `specs/<slug>/spec.md` may carry
one header line `Introduces commands:` listing its slash and Codex phase tokens; those tokens are valid
inside `specs/<slug>/` only, only while the spec's `Status:` is `staged`. Every other
file, and every undeclared token, keeps the fail-closed rule (reviewer decision
2026-08-18; recorded in `capture/learnings.md`).

## Acceptance

- [ ] A staged spec naming a declared token in `spec.md` and in `tickets/*.md` produces no diagnostic.
- [ ] The same token in `README.md`, `.bench/BENCH.md`, or another spec's directory is still red.
- [ ] An undeclared `/bench-*` token in the declaring spec's directory is still red (canary fixture asserts it).
- [ ] An implemented spec's declaration grants nothing.
- [ ] `bench-write-spec.md` and the craft-spec template name the header line.
