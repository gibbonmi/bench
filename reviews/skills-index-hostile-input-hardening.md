# Review pickup — skills-index hostile-input hardening

Frozen base `9c369749` → reviewed tip `e817c382`. Three axes, `opus`/medium, read-only.
Raw findings 13; de-duplicated repair targets 6.

## Standards

Findings 5. Worst: a sixth derivation of the missing-tool wording, where a named
`toon` constructor is the established single source and the diff declined it.

- `ask-user` — `internal/skillsindex/command.go:73` formats
  `"required tool is missing or not executable: "+launch.Name` inline. `internal/toon/toon.go:139`
  (`NotInRepo`) and `:133` (`MissingArg`) document themselves as the one source for this
  class of phrasing; the literal already exists at `internal/preprelease/preprelease.go:160`
  and `internal/releasepreflight/command.go:38,242,300` and `vulnerability.go:31`. The fix is a
  `toon.MissingTool(tool)` constructor, which needs `internal/toon/` — outside every ticket's
  approved write fence.
- `auto-fix` — `internal/skillsindex/skillsindex.go:387` calls `refusedDiagnostic` "the one
  shape every untrustworthy-producer refusal takes", but `internal/anchors/registry.go:124`
  authors a competing grammar and `internal/conformance/checks_test.go:649-650` consumes both
  in one table.
- `ask-user` — `internal/skillsindex/skillsindex.go:461-480`: `findBlock` collapses zero,
  start-only, end-only, reversed, duplicate-start, and duplicate-end into `markersDiagnostic`
  (`:49`), while `readReference` at `:491-494` states that states "stay apart because they send
  an operator to different repairs". A duplicate-marker file is told it names no block.
  Per-shape diagnostics are a behavior change the spec did not decide.
- `auto-fix` — `internal/skillsindex/skillsindex.go:356` narrates prior code ("the one failure
  path cleanup used to miss"), which `craft-comments` forbids.
- `no-op` — `Indexed()` at `:78` re-derives a guard both call sites already enforce.

## Spec

Findings 4. All 14 coverage rows audited row by row; 13 delivered at their mapped seam.
Worst: HI14 is the only row satisfied below its named seam.

- `ask-user` — HI14 names "Registered `package-core-guard`", but
  `internal/conformance/package_core_checks_test.go:378` calls the guard function directly.
  The registered binding reaches the payload reader only after a real `npm pack --dry-run`
  (`:79` consumes `probe.Stdout`), which a `t.TempDir()` root cannot produce, so the registered
  route would pass vacuously. It is the production function, not a replica, and the exception is
  documented in-comment at `:375-377`. Accept the lower seam, or authorize an npm-backed fixture.
- `no-op` — the cut's boundary moved: `guidance_token_sweep_test.go:82` classifies every file
  under `.agents/commands` and `.agents/skills` unconditionally, and `prose_budget_test.go:78`
  and `anchors/registry.go:131` classify `.bench/BENCH.md` and command anchor targets. Net new
  behavior: an oversized or non-UTF-8 command file is now a gate red. The anchors route is
  spec-authorized and a per-path gate inside one read helper would be worse.
- `auto-fix` — `internal/skillsindex/skillsindex_test.go:194`:
  `TestUnparseableAllowlistRefusesWriteAndLeavesCheckAlone` asserts at `:220-223` that Check
  *does* refuse. Mysterious Name; rename.
- `no-op` — `checkSkillFrontmatter` still enumerates by glob
  (`internal/conformance/skills_index_checks_test.go:82`), so under a glob-shaped root it grades
  zero skills. No row maps it and HI2's seam is the skillsindex module; the CHANGELOG claim is
  what reads wider than what landed.

## Coverage

Findings 4, 2 genuinely open. Worst: the packagesurface partition asserts a refusal without
asserting which refusal.

- `auto-fix` — `internal/packagesurface/contract_documents_test.go:18-19` maps both `absent`
  and `empty` to `""` and asserts only that an error naming `consumerPayloadPath` came back
  (`:38`, `:41`). Both hold when a present-but-empty payload is reclassified as absent, because
  `contract_documents.go:53` makes absence its own error naming the same path. The composition
  row does not rescue it: `internal/conformance/package_core_checks_test.go:399-401` asserts only
  `containsDiagnostic(diags, rel)`.
- `auto-fix` — no fixture drives a reference whose last line lacks a trailing newline;
  `reference()` at `internal/skillsindex/skillsindex_test.go:32-34` always terminates with `"\n"`,
  as does every malformed row at `:705-711`. This is a named class in the profile's hostile-input
  checklist.
- `no-op` — FIFO is absent from the `ClassifyNoFollow` table
  (`internal/bounds/classify_test.go:323-420`), but covered through the same `gradeBytes` path by
  `skillsindex_test.go:516-524` and `hostileSkillPlanters`.
- `no-op` — the glob-shaped root is exercised through the library, not the CLI
  (`command_test.go:18`, `:134`); the command delegates to the same functions.

Classes confirmed closed: absent-vs-empty for SKILL.md and the reference, dangling symlink,
special files, the permitted half of the control-byte partition, second-write byte equality,
SIGINT residue, and cross-process re-read.
