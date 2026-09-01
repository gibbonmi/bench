# Cross-harness reviewer invocations

Use a family CLI for a cross-family reviewer. A harness without a native
subagent surface falls back to its own family's CLI.

Close stdin on every one of these. A family CLI that receives both a prompt argument and
an open stdin waits for stdin to supply more prompt. So a backgrounded reviewer
parks before it starts and reports no error — it looks launched. Redirect from
`/dev/null` and check the process is alive before treating a fan-out as running.

Claude:

`claude -p --model <id> --effort <level> "<charge>" < /dev/null`

Codex:

`codex exec --sandbox read-only -C <dir> -m <id> -c model_reasoning_effort=<level> -o <file> "<charge>" < /dev/null`

`-C` sets the reviewer's working root, which a backgrounded call cannot inherit
from the caller's shell. `-o` writes the final message alone, so findings are
read without parsing the event stream around them.

Inside a Bench worktree, take the exec form with an empty quoted heredoc,
because the guard refuses any non-heredoc redirection inside an exec span:

`bench worktree exec <target> -- claude -p --model <id> "<charge>" <<'EOF'`

`EOF`
