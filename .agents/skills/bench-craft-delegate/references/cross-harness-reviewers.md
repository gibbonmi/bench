# Cross-harness reviewer invocations

Use a family CLI for a cross-family reviewer. A harness without a native
subagent surface takes the capable-harness handoff in
`.agents/commands/bench-review-implementation.md`, and it substitutes no
same-family CLI.

Each direct recipe reads its charge from a file on stdin. The shell does not
expand the text of a redirected file, so a backtick or a `$()` in the charge stays
inert. In a double-quoted argument, the coordinator's shell runs each one as a
command substitution. The file also closes stdin at its end. A family CLI with an
open stdin waits for more prompt, so a backgrounded reviewer parks before it starts
and reports no error. Check that the process is alive before you treat a fan-out as
running.

Claude:

`claude -p --model <id> --effort <level> < <charge-file>`

Codex:

`codex exec --sandbox read-only -C <dir> -m <id> -c model_reasoning_effort=<level> -o <file> < <charge-file>`

`-C` sets the reviewer's working root, which a backgrounded call cannot inherit
from the caller's shell. `-o` writes the final message alone, so findings are
read without parsing the event stream around them.

Inside a Bench worktree, the guard refuses a non-heredoc redirection. The exec
form carries the charge in the body of a quoted heredoc. The quoted delimiter
keeps the body inert, and the heredoc closes the child's stdin at `EOF`:

```bash
bench worktree exec <target> -- claude -p --model <id> --effort <level> <<'EOF'
<charge>
EOF
```

```bash
bench worktree exec <target> -- codex exec --sandbox read-only -C . -m <id> -c model_reasoning_effort=<level> -o <file> <<'EOF'
<charge>
EOF
```
