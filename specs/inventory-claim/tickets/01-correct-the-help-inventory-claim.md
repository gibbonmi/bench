# Correct the help inventory claim

Blocked by: none
Writes: .bench/BENCH.md, .bench/BENCH-reference.md, projects/benchkit.md, CHANGELOG.md
Covers: none

## What to build

The operating guide says that `bench help` is the complete executable inventory. This claim is false. The registry classifies `bench gate-prose` as internal inventory, so `bench help` does not show it. A session must run that verb before a commit. A delegate obeyed the claim, did not find the verb, and did not do the prose step.

The operating guide keeps the pinned phrase and limits the claim to the public verbs. The same sentence points to `.bench/BENCH-reference.md` for the one internal verb that a session runs. The reference stays the one source for that verb and its forms.

The Command Notes section of the reference and the project profile also call `bench help` the executable inventory. Each copy limits that claim to the public verbs. No code changes, and `bench gate-prose` stays internal.

## Acceptance

- [ ] `.bench/BENCH.md` says that `bench help` is the complete executable inventory of the public verbs.
- [ ] The same sentence in `.bench/BENCH.md` points to `.bench/BENCH-reference.md` for the one internal verb that a session runs.
- [ ] `.bench/BENCH-reference.md` and `projects/benchkit.md` limit their inventory claim to the public verbs.
- [ ] `bench help` still omits `bench gate-prose`, and the conformance check for the inventory phrase stays unchanged and passes.
- [ ] `CHANGELOG.md` records the corrected guide.
