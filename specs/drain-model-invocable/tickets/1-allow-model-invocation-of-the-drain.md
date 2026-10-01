# Allow the Claude model to start the drain

Blocked by: none
Writes: internal/conformance/skills_index_checks_test.go, .agents/commands/bench-drain.md, .bench/BENCH-reference.md, CHANGELOG.md
Covers: none

## What to build

The reviewer decided that the Claude model may start `/bench-drain` on its own.
Before this ticket, the drain command carried `disable-model-invocation: true`,
so only a typed command could start it. A phase close that pointed at the drain
then had to stop and hand the command back to the reviewer.

The invocation-policy table marks `bench-drain` as model-invocable on Claude.
The command frontmatter drops the `disable-model-invocation` key to agree with
that row. The harness-invocation prose names the drain as the one maintenance
exception. The reason is that the drain lands nothing until the reviewer
approves its batch diff.

The Codex adapter keeps `allow_implicit_invocation: false`. The renamed alias
`bench-what-next` stays user-invoked.

## Acceptance

- [x] The `phaseInvocationPolicy` row for `bench-drain` sets `claudeModelInvocable: true` and leaves `codexImplicit` false.
- [x] `.agents/commands/bench-drain.md` frontmatter has no `disable-model-invocation` key.
- [x] `.bench/BENCH-reference.md` names `/bench-drain` as the maintenance exception on Claude.
- [x] The existing invocation-policy check turns red when the frontmatter and the policy row disagree.

## Verification

The probe restored the `disable-model-invocation: true` line in the command
frontmatter and ran the skills-index command-adapter check. The check was red
with the policy-mismatch diagnostic for `bench-drain`.
