# Contract the six-column coverage branch

Blocked by: none
Writes: internal/coverage/coverage.go, internal/coverage/coverage_test.go, tests/canary/coverage-map-validation/broken-coverage-map/

## What to build

`bench coverage` accepts four header shapes. Two are legacy — they carry the
retired `red signal` cell — and two are the reduced schema. The reduced schema
is now the only form any live spec uses, so the legacy pair contracts away and
`internal/coverage` is left with one column list per accepted header.

The deletion in `internal/coverage/coverage.go` is the `fieldRedSignal` const and
the first two entries of the `schemas` slice (the no-row-ID legacy form and its
row-ID variant). The `schema` type, its `header`/`known`/`width`/`optIn`/`index`/
`cell` methods, and `schemaFor` are shared with the reduced schema and all
survive unchanged — the descriptor list is the only edit, which is what the type's
own doc comment promises ("Adding a header is adding a descriptor here, not
editing the checks").

Two consequences are easy to miss, and both are the real work in this ticket:

**Order the surviving descriptors so the unknown-header projection does not
move.** `projection()` falls back to `schemas[0]` for a header matching no
descriptor, and its doc comment currently says rows "project through the legacy
field order". That fallback output was pinned by the reviewer at `44830fa1`
("pin the unknown-header projection"). If the two legacy entries are simply
removed, `schemas[0]` becomes the reduced row-ID form, whose leading `row` cell
shifts every field one place right, and a 4-cell unknown-header row that projects
`"1",b,s` today would start projecting `"b",s,w`. Putting the 4-cell reduced form
`story|behavior|seam|why it catches the failure` first keeps `story`, `behavior`,
and `seam` at offsets 0, 1, and 2, so the projection stays byte-identical.
Reword the `projection()` comment, which will otherwise name a field order that
no longer exists.

**`TestFiveCellHeadersSelectSchemaByName` cannot survive, and that is correct.**
It exists to prove story 4's guarantee — the parser picks a schema by cell
*names*, never by cell count — and it does so by driving one byte-identical row
through the two five-cell headers that differ only in naming `red signal`.
Deleting the legacy pair leaves descriptors of width 5 and 4, so no two accepted
headers share a width and the confusion class the test guards is structurally
gone rather than merely untested. Retire the test with a comment recording that,
so a later reader does not read its absence as a dropped check. Do not invent a
synthetic same-width header to keep it alive.

Nothing outside `internal/coverage` and its own fixtures depends on the six-column
form parsing. Verified against the tree at `a3752447`: no `specs/*/spec.md`
carries a six-column map (`spec-ticket-fence-reduction` was the last one and
retired at `a3752447`); the `red signal` hits in `docs/field-guide.html`,
`docs/reporesident-distillation.md`, `capture/FIXES.md`, `ROADMAP.md`,
`decisions/`, `.agents/skills/bench-craft-tdd/SKILL.md`, and
`.agents/commands/bench-what-next.md` are prose about the retired column;
`internal/status/status.go` is an unrelated "shared signal ordering" phrase; the
one `internal/anchors/registry_data.go` hit is inside a `Diagnostic` string, not
a needle; and the six-column tables under
`tests/canary/workflow-guidance-anchors/*/files/dot-agents/commands/bench-write-spec.md`
sit in fixture copies of a command file that the coverage parser never reads.

**Migrate the canary fixture rather than deleting it.**
`tests/canary/coverage-map-validation/broken-coverage-map/` is the fixture that
proves the legacy branch runs, and its `files/specs/bad-map/spec.md` says so in
prose that must be rewritten. But its `EXPECT` is `coverage map row 1 has an
empty 'seam' cell` — a schema-independent empty-cell diagnostic that no other
fixture in `coverage-map-validation/` covers. Deleting the directory outright
would retire that diagnostic's only canary, which is weakening a check. Move the
fixture's map to the reduced header and keep the `EXPECT` byte-identical.

## Acceptance

- [ ] `bench coverage --check` rejects a spec whose map uses either `red signal` header, naming it as missing the canonical header rather than parsing it.
- [ ] `bench coverage` on a spec whose header matches no descriptor still emits `rows[1]{story,behavior,seam}` reading the first three cells in order, with its repair action unchanged.
- [ ] Every validation the reduced schema already enforces — story references, fan-out bound, duplicate and malformed row IDs, orphan stories, one predicate per behavior, empty cells, zero-row tables, the historical opt-out — behaves as it does today under both surviving headers.
- [ ] The `broken-coverage-map` canary still fails for its own planted empty-`seam` reason, and no canary fixture in the tree asserts a `red signal` column.
- [ ] `internal/coverage` names `red signal` nowhere outside prose describing the retired column.
