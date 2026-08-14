# Cross-harness reviewer invocations

Use a family CLI for a cross-family reviewer. A harness without a native
subagent surface falls back to its own family's CLI.

Claude:

`claude -p --model <id> --effort <level> "<charge>"`

Codex:

`codex exec --sandbox read-only -m <id> -c model_reasoning_effort=<level> "<charge>"`
