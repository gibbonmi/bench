# Learnings — usage journal

## 2026-08-21 — a spec's tickets left the write-spec phase failing bench preflight build [open]

`/bench-implement-spec specs/learnings-dated-line-visibility/spec.md --full`
stopped at its first entry check. `bench preflight build` returned
`rows-owned,red,"declared row(s) cited by no ticket file: DL1 … DL21"`: the one
ticket's acceptance lines are prose, and the exemplar shape
(`specs/land-executable-freshness/tickets/01-*.md`) prefixes each line with its
row ID. The mapping was already one-to-one, so the defect was pure citation.
What happened: the write-spec phase ran a verification round, accepted the
ticket, and exited, but never ran the build phase's own entry check against its
output. The build phase then refused work it does not own the fix for. Right
behavior: `/bench-write-spec` runs `bench preflight build <slug>` before it
exits, so a red lands in the phase that owns the ticket file. Proposed rule
change: `/bench-write-spec`'s exit adds that command as a required check, the
same way the build phase requires it at entry.
