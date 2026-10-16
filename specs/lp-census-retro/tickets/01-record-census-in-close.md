# Record the census in the close and the retro

Blocked by: none
Writes: .agents/commands/bench-final-check.md, .agents/skills/bench-craft-delegate/SKILL.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, CHANGELOG.md
Covers: none

## What to build

The reviewer decided on 2026-10-16 that a landing census writes no `bench learning` entry.
A census entry in the journal needs no drain verdict, so each such entry adds work to the drain and no value.
The `/bench-final-check` post-merge tail reads `census=<n>` and keeps the per-verb account in the phase close.
For a nonzero count, the close names the Bench form for each verb head and proposes the verb or the help change.
A spec's implementation retro records that account as one `### Bench CLI` item with its `Feeds:` line.
The `craft-delegate` charge folds the delegate's Bench CLI improvements into the same account of the close.

The anchor rows that pinned the old learning duty pin the new sentences.
The `CHANGELOG.md` unreleased census entry states the new close duty.

## Acceptance

- [ ] `/bench-final-check` tells the phase close to write no `bench learning` entry for a census.
- [ ] `/bench-final-check` tells a spec retro to record the census account under `### Bench CLI`.
- [ ] `craft-delegate` folds delegate CLI improvements into the census account of the phase close.
- [ ] The anchor gate reds when the final-check text restores a `bench learning` entry for a census.
