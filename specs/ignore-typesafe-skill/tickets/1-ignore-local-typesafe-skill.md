# Ignore the local TypeSafe skill install

Blocked by: none
Writes: .gitignore
Covers: none

## What to build

Git ignores the locally installed `typesafe-ai` skill, its Claude Code link, and the `skills-lock.json` file that records its source. This repo does not maintain the skill. The ignore file also excludes a local `.env` file.

## Acceptance

- [ ] `git status` does not list `.agents/skills/typesafe-ai/`, `.claude/skills/typesafe-ai`, or `skills-lock.json`.
- [ ] Git ignores a `.env` file at the repo root.
