# Order the chunk record and the author charge

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md
Covers: none

## What to build

`/bench-implement-spec` does not state the order of the chunk record in a multi-ticket chunk.
In four builds, the last author committed the chunk record before the earlier authors staged their records.
The retros of `guard-scan-loop` and `plan-commit` ask for this order.

The command also does not state when the orchestrator prepares an author charge.
A charge that the orchestrator prepares before the assignment-record commit fails `bench preflight evidence <id> --check-current`.
The retros of `candidate-binary-install` and `merge-target-lane` ask for this order.

State each order once, at its one owner in `.agents/commands/bench-implement-spec.md`.
Point each other sentence that touches the same order to that owner.
Keep every anchored needle of the file intact.

## Acceptance

- [ ] The Land section states that each ticket author stages its verification and probe record at the chunk source with only the record path dirty.
- [ ] The Land section states that the last author of a multi-ticket chunk commits no record before each earlier author stages its record.
- [ ] The Land section states that one record commit follows the last author's record.
- [ ] The venue section states that the orchestrator prepares the `bench preflight build --charge` charge at the assignment-record commit, never before it.
- [ ] The Build section and the `--delegate` option point to these owners and do not restate either order.
- [ ] `bench anchors` reports every anchored needle of the file intact, and the guidance prose budget and prose mechanics checks pass.
