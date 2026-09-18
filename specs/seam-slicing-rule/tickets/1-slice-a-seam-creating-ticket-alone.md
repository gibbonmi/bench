# Slice a seam-creating ticket on its own

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/SKILL.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, tests/canary/workflow-guidance-anchors/ticket-seam-creating-slice (new)
Covers: none

## What to build

The ticket slicing skill states one more breakdown rule, and the gate holds that rule.

A ticket that creates a seam its sibling tickets consume is its own slice, and that slice stays small.
Its chunk review closes before any consumer ticket starts.
A seam the implementation introduces counts the same as a seam the spec declares.

The skill already sizes a ticket by its context window and by the lines its author must read.
Neither rule reaches a small ticket whose output is machinery other tickets build on.
The new rule names that case, so a breakdown separates it before the consumers reach it.

The rule sits in the `Draft the breakdown` section, beside the existing sizing prose.
Add the sentences; do not reword a sentence the anchor registry already pins.
Twenty-five anchor rows pin this file, and thirteen canary fixtures quote its sentences.
Resolve every candidate sentence against `internal/anchors` and `tests/canary` before you write.

One anchor row holds the new rule, in the `ticket passes` registry beside its siblings.
Its independent expectation joins the table in that registry's test.
One canary fixture deletes the pinned sentence and expects the row's diagnostic.

## Acceptance

- [ ] The `Draft the breakdown` section states that a seam-creating ticket is its own slice, and that its review closes before a consumer starts.
- [ ] The section states that a seam the implementation introduces counts the same as a seam the spec declares.
- [ ] A `RequireInSection` anchor row pins the new rule to the `Draft the breakdown` section, with its own diagnostic.
- [ ] The registry test's expectation table names the new row's file, section, needle, and diagnostic.
- [ ] A canary fixture deletes the pinned sentence, and the fixture bites with the row's diagnostic.
- [ ] `bench test --check docs-currency-workflow` passes, and every existing anchor row and canary fixture still bites.
