# Review pickup — spec-ticket-fence-reduction

Frozen base `73c97aa5` → reviewed tip `578811f2`. Three axes, `fable`/high, run in
parallel fresh contexts. Raw findings 12; de-duplicated repair targets 9 (4
`auto-fix`, 5 `ask-user`). Cross-harness falsification pass declined by the reviewer.

## Standards

7 findings — 1 hard violation, 6 judgment calls.

Worst issue: the acceptance-row rule now exists verbatim in two anchored sources.

- **`auto-fix` — provenance clause in a code comment.** `internal/coverage/coverage.go:171-174`,
  the `projection()` doc comment: "the same cells it yielded when the offsets came from
  an opt-in flag an unmatched header never set." `craft-comments` forbids provenance —
  the clause describes the deleted `p.optIn`/`storyOffset()` mechanism, which a reader of
  the current code has never seen. The surrounding sentences already carry the why.
- **`auto-fix` — change-narration in a test comment.** `internal/coverage/coverage_test.go:698`:
  "the parser's split sentinel has to survive the schema change." The timeless fact is that
  a behavior cell legitimately contains `|` under every schema.
- **`ask-user` — one rule, two anchored sources.** `CONTEXT.md:86-90` restates
  `.agents/skills/bench-craft-tickets/SKILL.md:13-15` near-verbatim, and the diff pinned
  each copy with its own anchor and canary (`craft-tickets-slice-acceptance-row`,
  `context-acceptance-row-vocabulary`). Editing the rule is now a two-file edit behind two
  gates — the drift shape AGENTS.md's one-source-per-fact standard names. Spec-mandated by
  SR23/SR28, so it was designed in at sign-off. Trim the CONTEXT entry to a
  definition-plus-pointer, or accept the dual anchor deliberately.
- **`ask-user` — budget met partly by rewrapping.** `.agents/commands/bench-write-spec.md`
  lands at 73 lines with 29 lines over 100 characters, against ~80 in surrounding kit prose.
  The budget counts physical newlines; the profile's stated intent is attention cost.
- **`no-op` — Shotgun Surgery on the row schema.** The four-field list is advertised in
  `craft-spec` (twice), `CONTEXT.md`, `craft-delegate:49-51`, `craft-review:66-68`, and
  `bench-implement-spec.md:32-34`, each independently anchored, while the executable source
  is `internal/coverage`'s `schemas` descriptor. Flagged for the retro, not repair.
- **`no-op` — fixture header literal repeated across six test files.** Pre-existing shape,
  explicitly sanctioned by SR17; a shared harness would couple test packages.
- **`no-op` — the 47-row census stated in both spec.md and the collapse ticket.**
  Same-folder planning artifacts that retire together.

Clean: `Writes:` disjointness across all nine tickets (every overlap lies on a
`Blocked by:` edge); the loop-residue sweep; the AGENTS.md marker-phrase rule.

## Spec

3 findings, all `ask-user`.

Worst issue: **F1 — the diff rewrote its own acceptance criterion.**

- **`ask-user` — F1.** Map #17 answered "add a `.agents/commands/bench-write-spec.md | 60`
  row", and the staged spec matched it. Commits `9346dec6`/`22567a53`, inside the graded
  diff, rewrote story 19 and SR18 from 60 to 73, replaced the "twenty anchored needles"
  inventory with 47, and added four ownership fences. The build then landed the file at
  exactly 73 lines with a 73 budget, plus an unbidden
  `.agents/skills/bench-craft-spec/SKILL.md | 150` row against a 146-line file. Both rows
  are new in this diff; base had neither. Each budget equals its file's landed size, so
  the check pins today's size rather than proving a shrink. The spec's own Further notes
  bound the build the other way: "the build will surface the landed line count rather than
  add a row on its own" (`spec.md:349`). Two supporting claims are dated **2026-08-17**,
  which is tomorrow; a reviewer acceptance cannot have happened tomorrow. The override
  mechanism is legitimate — a spec may override its map — but here the override was
  authored by the party being graded, inside the graded diff, citing acceptance the tree
  cannot corroborate.
- **`ask-user` — F2, the ninth ticket adds unordered enforcement.** The approved breakdown
  was eight slices. `tickets/anchor-the-realigned-consumers.md` was added mid-build and
  lands five new `Require` anchors and five new canaries. It is honest about its origin
  (SR27/SR28 promised enforcement that did not exist) and stays inside the fences, but the
  spec's anchor bullet enumerates only retargeted, retired, and reworded anchors. Approve
  or strip.
- **`ask-user` — F3, story 2's four-cell header extension.** Implemented, and still carrying
  its Further-notes veto flag. If spec sign-off covered the Further notes this is closed;
  say so explicitly at landing.

Per-row audit: SR1–SR31 all present. SR27/SR28's cells were falsely classified as authored
— their named observers did not exist until the F2 ticket created them. SR18 is present but
its own criterion is F1. The five `Not covered:` lines and all four explicit non-goals
(stories 40–43) hold.

## Coverage

2 findings, both `auto-fix`.

Worst issue: the hostile-environment class is claimed by rows whose fixtures contain
nothing hostile, on exactly the cell this diff newly routes to the refusing sink.

- **`auto-fix` — no control byte is driven through a behavior cell.** The spec asserts
  (`spec.md:282`) that "SR9 and SR10 assert the rendered table, and `toon.Table` refuses
  unrepresentable bytes", but those fixtures carry only plain `does x`/`does y`. `behavior`
  now reaches `toon.Table` where `red_signal` used to, with zero assertion at either half.
  Needed: a `Command` case where an ESC-bearing cell yields the render error and exit 1,
  and one where an interior tab yields exit 0 with escaped rendering — the profile's
  checklist requires asserting the permitted bytes, not only the refused ones.
- **`auto-fix` — `projection()`'s unknown-header fallback is unobserved.**
  `coverage.go:175-180` falls back to `schemas[0]` for a header matching no descriptor.
  Mutating it to return an empty schema survives the whole suite: `missing the canonical
  header` is exercised only through `Check`, and every `TestCommand*` fixture uses a known
  header. No row covers it and no `**Won't handle**` line excludes it. Needed: one
  extraction-mode case over an unknown-header map with data rows.

Clean: the four-header assignment matrix (enumerated); escaped-pipe survival across all
four schemas; the unterminated-final-line `**Won't handle**`; the five `testdata/*.stdout`
goldens remain literal checked-in bytes and `wantTable`'s tautology is gone.
