# Name the completion order and the single final package run

Blocked by: none
Writes: .agents/commands/bench-final-check.md, .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-spec/references/completion-plan.md
Covers: none

## What to build

Two retros found guidance gaps.

The completion order has no single owner. `/bench-final-check` owns completion, but it names only the `--complete` checkpoint. `/bench-review-implementation` gives part of the order. In `/bench-final-check`, point the record writes to the order that `bench help` prints, because `reviewrecord.CompletionOrder` owns that order. Then state the two steps that no other source holds: commit the record, and then run the complete checkpoint. Make `/bench-review-implementation` point to `/bench-final-check` for that order.

The completion-plan reference states the probe-only rule for chunks, but it does not say where a package test that names no probe goes. Two builds amended their plans at build entry to put these tests in one `./...` final run. State that rule in the `final_verification` key.

## Acceptance

- [ ] `/bench-final-check` points the record writes to the order that `bench help` prints for `bench record completion`, then names a record commit, then `bench gate --checkpoint <spec-path> --complete`.
- [ ] `/bench-review-implementation` points to `/bench-final-check` for the completion order and states no part of the order itself.
- [ ] The `final_verification` key in the completion-plan reference puts every package test that names no probe in one `bench test --package ./...` requirement.
- [ ] The lane passes.
