# Keep research paths in Writes and let the author narrow the fence

Blocked by: none
Writes: .agents/commands/bench-write-spec.md, .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-tickets/SKILL.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-review/SKILL.md, projects/benchkit.md, docs/field-guide.html, CHANGELOG.md, internal/anchors/, internal/conformance/, tests/canary/workflow-guidance-anchors/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go
Covers: none

## What to build

The reviewer directed this change on 2026-10-10, after the ticket-slicing comparison in `docs/research/2026-10-10-ticket-slicing-comparison.md`. On FT349, the per-chunk research found the real fence (recall 0.869). The slicer then dropped 10 real paths by design choices that no ticket line bound, and the cheap per-chunk reviews did not catch the drops. Two consultants signed the rules below, and the reviewer approved them, with a ticket review line of mid + medium.

The change moves the design choice from the slicer to the ticket author. `Writes:` becomes the approved upper bound of a ticket, and the author narrows it.

1. **Research once.** After the spec review accepts, one round charges two read-only readers for each chunk, with every chunk at the same time. The invoking session merges each pair into one research table: the agreed set and the frontier. If the spec review changes the chunk table, the research reruns for each changed chunk.
2. **Keep every research path.** Each path that a research return names goes into the ticket's `Writes:` line. The slicer drops such a path only with `path:line` evidence that the path cannot change under any design. `Writes:` is the approved upper bound of the ticket, not a task list. This rule lives in the slicing checks, and the command points to it.
3. **One ticket review.** The invoking session still slices one chunk at a time. After the last chunk, one review round grades the whole ticket set under the `craft-tickets` quiz, with a cap of 2 iterations. It replaces the per-chunk review.
4. **Review line.** `craft-line`'s stage table replaces the chunk ticket review row with "Ticket review (`/bench-write-spec`) | mid + medium". The spec review row stays at cheap + high.
5. **Slicing line.** `craft-line`'s stage table adds "Ticket slicing (`/bench-write-spec`) | mid + high". When the invoking session runs on another model, a fresh writer on the mid binding slices the tickets.
6. **The author narrows the fence.** The author charge says that the author changes only what the acceptance needs. Before the first edit, the author lists every place the change must reach: callers, tests, fixtures, and exports. The author's return names each untouched `Writes:` entry with a one-clause reason.
7. **The review reads the fence the same way.** An untouched `Writes:` entry is not a defect. An edit outside the acceptance's reach is a finding, even when its path is in `Writes:`.
8. **Context checkpoint.** If the slicing session passes about 300k tokens of context, it compacts or restarts at a chunk boundary, never inside a chunk. It carries a ledger. The ledger holds the accepted tickets, the research table, the deferrals, the ordering edges, the spec flags, and the paths that wait for the fence.

Each fact has one owner. The slicing checks own the keep rule. `craft-line` owns both lines. `bench-implement-spec` owns the author and review charges. `bench-write-spec` owns the phase order and the context checkpoint, and it points to the others. The verification-log rule counts the spec round and the one ticket round.

The anchors and canaries that pin a changed sentence change with it. New require anchors pin the keep rule, the one ticket review, the two `craft-line` rows, the author narrowing, and the review reading. Restatements of the phase order in `projects/benchkit.md` and `docs/field-guide.html` change too.

## Acceptance

- [ ] The slicing checks keep each research-named path in `Writes:`, unless `path:line` evidence shows that it cannot change under any design.
- [ ] The slicing checks state that `Writes:` is the approved upper bound of a ticket.
- [ ] `/bench-write-spec` runs the ticket research for all chunks in one round after the spec review.
- [ ] `/bench-write-spec` runs one ticket review after the last chunk, in place of a review for each chunk.
- [ ] `craft-line`'s stage table holds "Ticket review (`/bench-write-spec`) | mid + medium" and "Ticket slicing (`/bench-write-spec`) | mid + high", with no chunk ticket review row.
- [ ] The `/bench-implement-spec` author charge carries the reach list before the first edit and the untouched-entry return.
- [ ] The `/bench-implement-spec` review charge states that an untouched entry is not a defect.
- [ ] `/bench-write-spec` states the context checkpoint and its ledger.
- [ ] The anchor harness goes red when each new key sentence is removed.
- [ ] `bench test --package ./internal/conformance`, `bench test --check root-conformance`, and `bench test --check docs-currency-workflow` pass.
- [ ] `CHANGELOG.md` records the change.
