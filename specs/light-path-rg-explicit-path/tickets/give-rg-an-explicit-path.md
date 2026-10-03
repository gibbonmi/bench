# Give every rg call an explicit path

Blocked by: none
Writes: AGENTS.md
Covers: none

## What to build

On 2026-10-03 a subagent ran `rg` with a pattern and no path. Its tool shell
had an open pipe on stdin, so `rg` read stdin instead of the tree. The call
hung for about two hours after the subagent returned its report. The main
session's tool shell has `/dev/null` on stdin, so `rg` searches the tree there.

Add one shell convention to `AGENTS.md`, directly after the `rg` convention.
The convention tells an agent to give every `rg` call an explicit path, unless
a pipe in the same command feeds it. Every harness and every subagent reads
`AGENTS.md`. A PreToolUse guard that enforces the rule is parked in
`capture/IDEAS.md` for the drain.

## Acceptance

- [ ] The shell conventions in `AGENTS.md` tell an agent to give every `rg` call an explicit path, unless a pipe in the same command feeds it.
- [ ] The convention states the cause: without a path, `rg` reads an open stdin pipe and waits until that pipe closes.
- [ ] `bench gate-prose . -- AGENTS.md` passes.
- [ ] The root conformance package stays green.
