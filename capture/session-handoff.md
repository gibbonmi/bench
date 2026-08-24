# Session handoff

Repository: `/home/mgibs/workspace/bench` (origin `https://github.com/gibbonmi/bench.git`)
Branch: `main` — HEAD `600d82a7` at drain start; this drain lands one commit on top
Spec: `specs/module-size-split/spec.md` (Status: staged)
Gate: green (fresh, drain worktree)

## State

This drain reconciled the roadmap against the tree: no shipped rows, nothing
to retire, no learnings, no retros. The two parked ideas opened two rows.
FT248 (`spec`) is the Bash guard that denies a shell follow-on appended to a
`bench` command. FT249 (`decide`) moves the idea inbox to a shared git ref.
The capture inbox and journal are both empty.

The reviewer moved the structure fix to the front of the sequence on
2026-08-23; the earlier deferral of `module-size-split` is lifted. The flow
report still shows a positive net delta. Two restructure candidates stand for
a future `/bench-drain --restructure`: fold FT200, FT207, FT178, FT162, and
FT173 into FT169's landing-authority theme, and fold FT213, FT214, FT236, and
FT237 into the craft-visit batch beside FT117, FT179, and FT94.

FT238's `Next:` is `ticket`, and no `tickets/FT238*` file exists yet.

## Next command

`/bench-implement-spec specs/module-size-split/spec.md`

## Shape

Rewritten in full at every phase close, pruned rather than accreted. A fresh
session pays for every line it reads cold; drop anything it would not act on.

Operational gotchas are placed by lifetime, not copied here. One that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build
belongs instead in that spec's coverage rows.

This file names at most when you'll hit one, never the command — a second copy
drifts from the source.

Keep the three sections above. **State** holds what is true now, including anything
uncommitted. **Next command** holds the exact harness-native invocation, not a
description of it. This section is the third.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
