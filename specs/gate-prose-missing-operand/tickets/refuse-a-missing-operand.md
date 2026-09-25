# Refuse a missing gate-prose operand

Blocked by: none
Writes: internal/gate/gate_prose.go, internal/gate/gate_prose_test.go, .bench/BENCH-reference.md
Covers: none

## What to build

The path-list form of `bench gate-prose` skips a named path that does not exist. The verb then
prints a `pass` row for that path. A grade on the wrong checkout therefore reads green.

The verb must check each named path before the grade. A path that the verb cannot find refuses
the run. The refusal uses the sentence that the staged form uses for an unreadable subject. The
verb writes one refusal line to stdout for each such path, and it exits 1.

The `--staged` form does not change. The prose package keeps its skip of an absent path. The
fast lane does not change, because it names only the Markdown that the composed tree holds. A
commit that deletes or renames a Markdown file therefore does not give the verb a missing path.

## Acceptance

- [ ] `bench gate-prose <root> -- <missing path>` exits 1 and prints no `pass` table.
- [ ] The verb writes one `refused unreadable subject` line for each missing path, and each line names its path.
- [ ] A missing path beside a clean path still refuses the run.
- [ ] The existing `gate-prose` tests and the lane tests pass with no change.
- [ ] The reference notes for `bench gate-prose` state the refusal of a missing path.
