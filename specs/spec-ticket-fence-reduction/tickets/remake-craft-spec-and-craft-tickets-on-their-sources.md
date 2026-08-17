# Remake craft-spec and craft-tickets on their sources

Blocked by: collapse-the-review-loop-and-narrow-reviewer.md
Writes: .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-tdd/SKILL.md, .agents/skills/bench-craft-tickets/SKILL.md, internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors

## What to build

No cell asks an author to assert current-code behavior before a test runs, and no
two sizing rules sit in the tree at once.

The four-form red-signal grammar is deleted from `craft-spec`. `craft-tdd`
becomes its one owner, classifying a row as `already covered` or `not TDD-able`
at the moment the row runs. The canonical edge-class walk moves to `craft-tdd`
too, at the seam where the classes are visible; `craft-spec` keeps the profile's
hostile-input checklist attachment and keeps `**Won't handle**` for deliberately
excluded edges, but no longer requires every walked class to produce a row or a
line.

`craft-tickets`' prose is already correct — commit 24cad87d landed the
vertical-slice rule, prefactoring-first, and the retirement of "smallest
independently-green" (the phrase appears nowhere under `.agents/`). What is
missing is enforcement: the acceptance-row clause at `craft-tickets:14` has no
anchor, so deleting it is silent. Add that anchor and a canary that plants the
deletion. Do not re-write the prose that already landed.

`craft-tdd`'s anchored sentence — "The row schema and the red-signal definition
are `bench-craft-spec`'s" — becomes false the moment the grammar leaves
`craft-spec`. Reword the sentence and its needle together so the anchor guards a
true claim. Verify that "independently-green implementation tickets" still reads
in README and CHANGELOG; do not edit those files.

## Acceptance

- [ ] `(covers SR21)` `craft-spec` contains no `observed red:` / `not observed:` / `already covered:` /
      `not TDD-able:` grammar, and `craft-tdd` names `already covered` and
      `not TDD-able`.
- [ ] `(covers SR21)` Deleting the red-signal grammar's replacement clause from `craft-spec`
      turns a new anchor red.
- [ ] `(covers SR22)` The canonical edge-class run appears in `craft-tdd`; **both** craft-spec
      edge-class anchors — the plain `re-run idempotency` Require and the
      `RequireInSection("The edge inventory")` class run — name `craft-tdd`, and a
      new canary plants the moved section.
- [ ] `(covers SR22)` `craft-spec` still requires `**Won't handle**` for reviewed exclusions and
      no longer requires a disposition per walked class.
- [ ] `(covers SR23)` The acceptance-row clause in `craft-tickets` is anchored, and deleting it
      turns that anchor red through a canary that plants the deletion.
- [ ] `craft-tickets`' existing prose is unchanged apart from anything an anchor
      requires.
- [ ] `(covers SR31)` `craft-tdd`'s red-signal-ownership sentence names the new owner and its
      needle matches; the old needle appears nowhere.
- [ ] `(covers SR24)` `independently-green implementation tickets` still reads in README and
      CHANGELOG, with their anchors green and neither file edited.
- [ ] `(covers SR26)` `bench anchors` is green and every touched canary reproduces its own
      `EXPECT` line.
