# Advertise the enforced grammar in the craft-tickets skill

Blocked by: build-ticket-file-parser.md
Writes: .agents/skills/bench-craft-tickets/SKILL.md, tests/canary/workflow-guidance-anchors, internal/conformance/fixture_bite_test.go, internal/tickets/example_test.go (new)
Covers: TG30, TG31

## What to build

The `craft-tickets` skill advertises the enforced grammar. The template gains
the `Covers:` field beside `Blocked by:` and `Writes:`. The prose states the
three enforced rules: the dependency rule, the ownership rule, and the mutation
rule. It also states that a `Writes:` path exists in the tree or carries the
`(new)` marker.

The marked example parses clean through the live parser. One ordinary test
reads the example between the `ticket-example` markers and feeds it to
`internal/tickets`. The test asserts zero diagnostics.

The skill stays within its guidance-prose budget row. The budget row keeps its
current value, so the edit trims elsewhere to make room. The anchors registry
pins several needles in this file, and every needle survives the edit.

The fixture that pins this file rides in the ticket. A moved or reflowed
clause needs its `workflow-guidance-anchors` fixture updated in the same diff.
`internal/conformance/fixture_bite_test.go` holds hardcoded needle rows pinned
to this file, and a reflow updates them in the same diff. The example test is
incidental Go under this prose line.

## Acceptance

- [ ] TG30 — the marked example parses through the live parser with no diagnostic.
- [ ] TG31 — the skill stays within its guidance-prose budget row.
- [ ] The template shows `Covers:` beside `Blocked by:` and `Writes:`.
- [ ] Every anchor needle the registry pins in this file still resolves.

## Delegate charge

You work in the Bench repo on the `ticket-grammar` spec. Read
`specs/ticket-grammar/spec.md` first. Then read
`.agents/skills/bench-craft-tickets/SKILL.md` in full. Read
`internal/anchors/registry_data.go` and note every needle it pins in that file.
Read `internal/tickets/tickets.go`, which a sibling ticket landed. Read the
hardcoded needle rows in `internal/conformance/fixture_bite_test.go` that pin
this file.

Read the `bench-craft-tickets` budget row in `projects/benchkit.md`.

Add the `Covers:` field to the template and to the marked example. Put it after
`Writes:`. State the dependency rule, the ownership rule, and the mutation
rule in the prose. State that a `Writes:` path exists in the tree or carries
the `(new)` marker.

Keep the file within its current budget row. Do not raise the row in
`projects/benchkit.md`. Trim redundant prose instead, and keep every anchor
needle intact.

Update each `tests/canary/workflow-guidance-anchors` fixture whose subject
line moved or reflowed. Keep each `EXPECT` line in agreement with the emitted
diagnostic. Update the hardcoded needle and replacement rows in
`internal/conformance/fixture_bite_test.go` when a pinned clause reflows.

Add `TestCraftTicketsExampleParsesClean` in `internal/tickets`. Read the skill
file at `../../.agents/skills/bench-craft-tickets/SKILL.md`. Extract the block
between the `ticket-example` markers. Feed it to the parser and assert zero
diagnostics.

Run `bench worktree exec ft174-ticket-grammar -- go test ./internal/tickets/`.
Then run `bench worktree exec ft174-ticket-grammar -- bench gate-prose . -- .agents/skills/bench-craft-tickets/SKILL.md`.
Then run `bench worktree exec ft174-ticket-grammar -- bench test --check guidance-prose-budgets`.
Then run `bench worktree exec ft174-ticket-grammar -- go test ./internal/conformance/...`.
Do not commit. Do not edit the spec.
