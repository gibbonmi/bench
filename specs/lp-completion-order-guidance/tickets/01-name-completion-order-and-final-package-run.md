# Name the completion order and the single final package run

Blocked by: none
Writes: .agents/commands/bench-final-check.md, .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-spec/references/completion-plan.md
Covers: none

## What to build

Two retros found guidance gaps.

The completion order has no single owner. `/bench-final-check` owns completion, but it names only the `--complete` checkpoint. `/bench-review-implementation` gives part of the order. Name the full order in `/bench-final-check`: record the final verification results, write the completion entry, commit the record, and then run the complete checkpoint. Make `/bench-review-implementation` point to `/bench-final-check` for that order.

The completion-plan reference states the probe-only rule for chunks, but it does not say where a package test that names no probe goes. Two builds amended their plans at build entry to put these tests in one `./...` final run. State that rule in the `final_verification` key.

## Acceptance

- [ ] `/bench-final-check` names the order `bench record verification --final`, then `bench record completion`, then a record commit, then `bench gate --checkpoint <spec-path> --complete`.
- [ ] `/bench-review-implementation` points to `/bench-final-check` for the completion order and states no part of the order itself.
- [ ] The `final_verification` key in the completion-plan reference puts every package test that names no probe in one `bench test --package ./...` requirement.
- [ ] The lane passes.
