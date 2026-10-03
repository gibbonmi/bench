# Add the block reader terms to the glossary

Blocked by: none
Writes: CONTEXT.md
Covers: none

## What to build

The staged FT358 spec, `specs/markdown-block-reader/spec.md`, uses three terms
that `CONTEXT.md` does not define: block reader, fenced block, and unfenced
line. A consultation on 2026-10-03 decided that the terms land on their own
light path. The spec branch then carries no `CONTEXT.md` edit, so its build
preflight has no path outside an ownership fence.

Add the three terms beside the prose exclusion row entry. Each term states
the rule it names and the synonym that the glossary refuses. The wording
matches the folded spec.

## Acceptance

- [ ] `CONTEXT.md` defines block reader, fenced block, and unfenced line.
- [ ] Each definition matches the reader rules in the folded FT358 spec.
- [ ] `bench gate-prose . -- CONTEXT.md` passes.
- [ ] The root conformance package stays green.
