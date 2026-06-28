# Claude Adapter

Claude Code reads this `.claude/` tree.

Bench keeps portable skills and commands in `.agents/`. Files under
`.claude/skills/` and `.claude/commands/` are adapter symlinks to `.agents/`.

Hook config in `.claude/settings.json` points to shared scripts in
`.bench/hooks/`.
