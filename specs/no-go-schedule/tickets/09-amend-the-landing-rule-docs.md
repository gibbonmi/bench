# Amend the landing rule docs for the reader set

Blocked by: 06-publish-a-reader-landing.md
Writes: docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md, projects/benchkit.md, .bench/BENCH-reference.md, CONTEXT.md, .bench/commitment.json
Covers: NG68, NG69, NG70

## What to build

ADR 0017 states the current decided rule. A kit landing whose every changed path sits in a covered family runs the reader set instead of the whole gate. Every other landing runs the whole gate. Amend `docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md:17-20` and `:30` to that state, with no history and no file path.

`projects/benchkit.md` states the reader set as the second reduced schedule. Amend `:351-353`, `:393-397`, and `:476-479`. Name the reader record as the second reduced record class, and keep the rule that a partial or component record is a legacy class.

`.bench/BENCH-reference.md` amends `:548`, which calls the landing's prospective execution a full boundary. Its "Gate output" section (`:624-658`) gains the reader-set lines:

- `gate: green (reader set)`
- `gate: whole gate (reader set refused: <case>)`
- the refusal suffix `; reader landings since the green marker: <commits>`

`CONTEXT.md` gains the **covered family** and the **reader hold**, each with an Avoid list, beside the reader terms. The `reader set` entry names the race phase rule.

Run `bench anchors` on each edited path before the edit. Keep every anchored line and every line that a retained canary fixture quotes unchanged.

## Acceptance

- [ ] ADR 0017 states that a kit landing inside the reader table runs the reader set, and that every other landing runs the whole gate.
- [ ] The profile names the reader set, the reader record, and the reader hold, and it keeps the legacy-class rule.
- [ ] The reference's gate output section names the three reader-set lines.
- [ ] `CONTEXT.md` defines the covered family and the reader hold, each with an Avoid list.
- [ ] `bench test --check docs-currency-workflow` passes on the edited tree.
