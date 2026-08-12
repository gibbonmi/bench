# Repair learnings refusal evidence

Blocked by: none
Writes: `internal/learnings/learnings_test.go`, `internal/learnings/testdata/`

## What to build

Complete the independently green old-to-new evidence for every learnings refusal class. Preserve each pre-disclosure primary response, stream, and exit exactly; append disclosure only for malformed entries discovered among otherwise parsed journal entries, where the query knows the repair. Early refusals that cannot reach a truthful remedy remain byte-identical and do not gain a help block.

## Acceptance

- [ ] [LR1] (covers QD6) the fixed input `# Learnings — usage journal\n\n## 2026-01-01 — first [open]\n## broken\n` has checked-in pre-disclosure stdout `learnings[2]{date,title}:\n  2026-01-01,first\n  line 4,malformed learning heading\n` on empty stderr and exit 1. The candidate differs only by `help[1]{cmd,why}:\n  /bench-what-next,"verdict 2026-01-01: first"\n`.
- [ ] [LR2] (covers QD6) parsed bytes with no heading attempt produce exactly `error: capture/learnings.md is unsupported-schema — no dated heading found\n` on stdout, empty stderr, and exit 1 before and after the change, with no help block.
- [ ] [LR3] (covers QD6) an empty regular file produces exactly `error: capture/learnings.md is empty — \n`; invalid UTF-8 produces exactly `error: capture/learnings.md is malformed — invalid UTF-8\n`; a file exceeding `bounds.ControlRecordLimit` produces exactly `error: capture/learnings.md is unreadable — read limit exceeded\n`; and a directory at the journal path produces exactly `error: capture/learnings.md is wrong-type — not a regular file: d---------\n`. Each is stdout with empty stderr and exit 1 before and after the change, with no help block.
