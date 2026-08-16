# Enumerate skill directories literally

Blocked by: harden-every-skill-file-reader.md
Writes: internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go

## What to build

Replace pattern-based `*/SKILL.md` discovery with literal child-directory enumeration
under `.agents/skills`, using the published no-follow classifier for each expected
file. A repository root containing spaces and `[ ] * ?` still finds its skill, and
check/write retain the generated row. Keep current missing-file handling temporarily;
the dependent orphan-policy ticket owns that newly observable state, so this tracer can
land fully green without conflating path discovery with diagnostic policy.

## Acceptance

- [ ] `(covers HI2)` A hostile literal root enumerates its one skill through Check and
  Write, and Write cannot erase the row through glob interpretation.

