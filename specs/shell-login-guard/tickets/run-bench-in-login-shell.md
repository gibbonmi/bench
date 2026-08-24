# Keep Bench commands in a login shell

Blocked by: none
Writes: AGENTS.md

## What to build

Agents keep the default login shell when they run a `bench` command. The rule
names the user-local binary path and the non-fatal Envman warning condition.

## Acceptance

- [ ] The shell conventions prohibit `login: false` for every `bench` command.
