# Delegate the light-path fix of a captured learning

Blocked by: none
Writes: .bench/BENCH.md, .agents/commands/bench-drain.md, .agents/skills/bench-craft-tickets/SKILL.md, docs/adr/0023-each-ticket-gets-a-fresh-author.md, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_data.go, CHANGELOG.md
Covers: none

## What to build

The reviewer decided a standing policy on 2026-09-25. A captured learning can have a fix that meets the light-path row and needs no reviewer decision. Then the session dispatches a fresh write delegate on the mid tier at high effort to implement that fix. The policy applies at every point in the workflow.

The operating guide states this rule next to the light-path table, as the one exception to the inline route of that table. The delegate works in its own bench worktree, so the active phase keeps its worktree and verdict. The coordinator verifies the done-claim under `craft-delegate` and lands the ticket as the light-path row states. The learning entry stays in the journal, and the drain closes it by implementation.

The drain, the tickets skill, and the fresh-author ADR each keep the fact through a pointer to the operating guide, not through a second copy. An anchor row pins each new owner sentence and each pointer.

## Acceptance

- [ ] `.bench/BENCH.md` holds the learning delegation rule, and `bench test --check docs-currency-workflow` reds when the rule's dispatch sentence is absent.
- [ ] The light-path table row and the new rule do not contradict, and the table row keeps its pinned text.
- [ ] `.agents/commands/bench-drain.md` routes a learning entry with a light-path fix to the operating guide's write delegate and closes the entry by implementation.
- [ ] `.agents/skills/bench-craft-tickets/SKILL.md` names `.bench/BENCH.md` as the owner of the light-path author, not a copy of the inline route.
- [ ] `.bench/BENCH.md` stays inside its prose budget, and `bench gate-prose` passes on each edited Markdown file.
