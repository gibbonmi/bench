# Rehearse the promotion broker before a broker-changing landing

Blocked by: none
Writes: .agents/commands/bench-final-check.md, internal/anchors, tests/canary/workflow-guidance-anchors
Covers: none

## What to build

The final check tells a phase when a spec changes the promotion broker source.
The signal is a diff that changes a Bench build input, and the landing prints `landing changes the promotion broker source` for that diff.
Before the first landing of such a spec, the phase runs `bench worktree build <target>` and then the built binary's `bench doctor` through `bench worktree exec`.
The rehearsal checks the broker seal and manifest on the candidate binary before the landing depends on them.
One Require anchor row pins the rehearsal sentence, so the `docs-currency-workflow` check reds when the sentence goes.

## Acceptance

- [ ] `.agents/commands/bench-final-check.md` states the rehearsal once, with the landing signal that selects it.
- [ ] `bench test --check docs-currency-workflow` passes on the tree with the rehearsal sentence.
- [ ] `bench test --check docs-currency-workflow` reds when a probe removes the rehearsal sentence.
- [ ] `.bench/BENCH-reference.md` carries no second copy of the rehearsal rule.
