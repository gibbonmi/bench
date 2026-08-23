---
description: Route from the repository's observed state into the one current Bench action.
---

# /bench

## Entry orientation

Run `bench status --route` and take its one row.

If the row's `command` opens with `/bench-` or `$bench-`, take its first token.
Remove the leading `/` or `$`. Read the corresponding
`.agents/commands/<token>.md` completely. Follow that file as the active Bench phase.
Otherwise, run the command exactly. Load nothing else beyond the routed phase or
command.

## Exit handoff

Report the exact command's result. When `command` is empty, report the row's state.
Stop.
