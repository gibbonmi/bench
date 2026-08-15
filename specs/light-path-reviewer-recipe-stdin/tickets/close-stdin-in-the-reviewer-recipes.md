# Close stdin in the cross-harness reviewer recipes

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md, internal/anchors/registry_data.go

## What to build

Both pinned reviewer commands gain `< /dev/null`, and the codex recipe gains
`-C <dir>` and `-o <file>`. A family CLI handed both a prompt argument and an
open stdin blocks waiting for stdin to extend the prompt, so a backgrounded
reviewer parks before it starts and exits no error — the fan-out looks
launched and produces nothing. Observed 2026-08-15 dispatching three review
axes for `skills-index-reader`: all three sat at `Reading additional input
from stdin...` until killed. `-C` supplies the working root a backgrounded
call cannot inherit from the caller's shell; `-o` writes the final message
alone so findings are read without parsing the surrounding event stream.

The two command strings are byte-pinned as `Require` needles in
`internal/anchors/registry_data.go` (group `AfterImplementSpec`), so the
recipe edit and the needle edit are one commit or the gate is red. The third
needle on that file — the no-native-surface fallback sentence — is unchanged,
as is the pointer needle in `bench-craft-delegate/SKILL.md`.

The canary fixture `tests/canary/workflow-guidance-anchors/cross-harness-reviewer-recipes/`
carries a stripped file that names no command at all, so it keeps biting on the
new needles with no edit — it stays the independently authored expectation.

## Acceptance

- [ ] both recipes carry `< /dev/null`; the codex recipe carries `-C` and `-o`
- [ ] the two `Require` needles in `registry_data.go` match the new command
      strings byte-for-byte, and the fallback needle is untouched
- [ ] `cross-harness-reviewer-recipes` still bites through its registered owner
      with no fixture edit
