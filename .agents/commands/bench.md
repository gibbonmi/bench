---
description: Route from the repository's observed state into the one current Bench action.
---

# /bench

## Entry orientation

Run `bench status --route` and take its one row.

If the command is `git push`, offer the reviewer a choice before you run it.
They can push now or continue with the committed work. For that work, run
`bench roadmap` and take its `commitment_outlook` row. State its outcome, then
run its `command` under the commitment rule in `.bench/BENCH.md`.

If the row's `command` opens with `/bench-` or `$bench-`, take its first token.
Remove the leading `/` or `$`. Read the corresponding
`.agents/commands/<token>.md` completely. Follow that file as the active Bench phase.
Otherwise, run the command exactly. Load nothing else beyond the routed phase or
command.

## Exit handoff

If the routed command fails, continue with `.agents/commands/bench-debug.md` as the active phase.
Read that file completely before you take the next action.

When `command` is empty, report the row's state and stop. Otherwise, report the exact
command's successful result and stop.
