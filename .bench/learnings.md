# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-integrate-learnings` reviews the open entries, promotes the
generalizable ones into the kit with sign-off, and prunes them: a resolved entry
leaves this file, and its verdict (promoted or dismissed, one line of why) is
recorded in the integration commit and CHANGELOG. The journal holds open entries
only; history lives in git. Never rewrite a kit rule yourself — that is the whole
point of capturing here instead.

Format per entry:

## <date> — <short title>  [open]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry leaves this file only via /bench-integrate-learnings.

<!-- entries below -->

- 2026-07-05  Shared-tree contention: while another process implemented binary-auto-repair in the main working tree, I ran retire/map doc work in the same tree. First gate run was green (committed the retire on it); by the map commit the gate was red from the other process's in-flight bin/bench.sh edit (shellcheck + 7 canaries not biting), so the doc commit is held even though the red is not attributable to the doc diff. Right behavior was probably to do side-work in a bench worktree whenever another writer owns the main tree, making every gate verdict attributable to one diff. Proposed rule: when `git status` shows another writer's in-flight edits, side-work goes to a worktree — or waits.
