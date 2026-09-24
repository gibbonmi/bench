# Claude Adapter

Claude Code reads this `.claude/` tree.

Bench keeps portable skills and commands in `.agents/`. Files under
`.claude/skills/` and `.claude/commands/` are adapter symlinks to `.agents/`.

`.claude/skills/` links every `.agents/skills/` skill that has no same-named command.
Each `$bench-*` phase adapter skill has a same-named command, so it stays
Codex-only: Claude already gets each phase as a `.claude/commands/` slash
command, and a link to the skill would give the slash menu two entries per phase.

In Claude Code, `/bench-implement-spec` seeds the `TodoWrite` list from
`bench coverage <spec>`, one todo per acceptance-coverage row. That binding is
recorded here because the neutral command file names only the generic
task-list surface.

`.claude/agents/` holds the Bench agent types, one real file each. A Bench
agent file declares its name, its description, and the tools it needs, and it
declares no model. A charge names the agent type and passes the bound tier
token, so the agent-line guard keeps its verdict. `bench-reviewer` runs a
review axis or a diagnostic consultation.
`bench-writer` is the write delegate for a fresh ticket author, a repair session, or a user-directed write delegation.
A consumer's own agent files sit beside them and stay ungraded.

Hook config in `.claude/settings.json` points to shared scripts in
`.bench/hooks/`.
