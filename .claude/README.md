# Claude Adapter

Claude Code reads this `.claude/` tree.

Bench keeps portable skills and commands in `.agents/`. Files under
`.claude/skills/` and `.claude/commands/` are adapter symlinks to `.agents/`.

`.claude/skills/` links only the `bench-craft-*` skills. The `$bench-*` phase
adapter skills stay Codex-only: Claude already gets each phase as a
`.claude/commands/` slash command, and linking the same-named skill gives the
slash menu two entries per phase.

Hook config in `.claude/settings.json` points to shared scripts in
`.bench/hooks/`.
