# Tighten comment and review rules

Blocked by: codify-load-stop-and-quiet-check.md
Writes: .agents/skills/bench-craft-comments/SKILL.md, .agents/skills/bench-craft-review/SKILL.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go
Covers: LF15

## What to build

Add identifier-provenance, one-source, constraint-first, sparse-file, and
red-record ownership rules. Reject new FT story provenance tags without
sweeping untouched sites.

## Acceptance

- [ ] craft-comments contains all five operational rules.
- [ ] craft-review rejects a newly introduced story provenance tag.
- [ ] Anchor tests protect the rules without copying their full prose.
- [ ] Untouched legacy tags remain outside the ticket.
