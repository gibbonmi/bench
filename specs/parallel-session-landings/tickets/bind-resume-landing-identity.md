# Bind resume to the landed source identity

Blocked by: centralize-request-digest.md
Writes: internal/intent, internal/worktree

## What to build

Preserve the released assignment branch and its exact tip in the existing
terminal receipt `Branch` and `BranchOID` fields. Require that durable tip to
equal `--source-tip`; combine it with the exact published commit's destination
ancestry and second-parent proof for active and completed resume.

## Acceptance

- [ ] Active and completed resume authenticate the same source tip and exact published commit relationship.
- [ ] Wrong published commit, wrong source tip, mismatched terminal receipt, and evicted receipt refuse unchanged.
- [ ] No new receipt field, schema, journal, or commit format is introduced.
