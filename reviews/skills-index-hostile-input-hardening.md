# Review pickup — skills-index hostile-input hardening

Frozen base `9c369749` → reviewed tip `e817c382`. Three axes, `opus`/medium, read-only.
Raw findings 13; 6 repair targets accepted and closed; 3 findings remain, all needing
reviewer authority. The closed targets are gone from this file — it is pickup state, not
a review log.

## Standards

Remaining 1. Worst: the missing-tool wording has no single source, and its fix is fenced out.

- `ask-user` — `internal/skillsindex/command.go:73` formats
  `"required tool is missing or not executable: "+launch.Name` inline. `internal/toon/toon.go:139`
  (`NotInRepo`) and `:133` (`MissingArg`) document themselves as the one source for this class of
  phrasing; the literal already exists at `internal/preprelease/preprelease.go:160` and
  `internal/releasepreflight/command.go:38,242,300` and `vulnerability.go:31`.
  **This is one repair with the reference-grammar finding below.** `internal/skillsindex` declares
  `refusedDiagnostic` (`internal/skillsindex/skillsindex.go:387`) "the one shape every
  untrustworthy-producer refusal takes" while `internal/anchors/registry.go:124` authors a
  competing grammar, and `internal/conformance/checks_test.go:649-650` consumes both in one table.
  A repair pass unified them by pointing `anchors` at `skillsindex`; that was reverted, because it
  inverts the dependency (a generic prose-anchor registry importing a domain package for a string)
  and silently changes a user-visible anchors diagnostic no gate check pins. Both findings close
  together with one `toon` constructor — needs `internal/toon/`, outside every approved fence.
- `ask-user` — `internal/skillsindex/skillsindex.go:461-480`: `findBlock` collapses zero,
  start-only, end-only, reversed, duplicate-start, and duplicate-end into `markersDiagnostic`
  (`:49`), while `readReference` at `:491-494` states that states "stay apart because they send an
  operator to different repairs". A duplicate-marker file is told it names no block — true, but it
  does not point at the duplicate. Per-shape diagnostics are a behavior change the spec did not
  decide.

## Spec

Remaining 1. All 14 coverage rows audited row by row; 13 delivered at their mapped seam.

- `ask-user` — HI14 names "Registered `package-core-guard`", but
  `internal/conformance/package_core_checks_test.go:378` calls the guard function directly, so a
  future `checkPackageCoreAndGuards` change that stops reaching `checkNoKitOnlyPackedAssets`
  (`:79`) would not go red here. The registered binding reaches the payload reader only after a
  real `npm pack --dry-run` (`:79` consumes `probe.Stdout`), which a `t.TempDir()` root cannot
  produce, so the registered route would pass vacuously. It is the production function, not a
  replica, and the exception is documented in-comment at `:375-377`. Accept the lower seam, or
  authorize an npm-backed fixture.

## Coverage

Remaining 0. Both open gaps closed in this pass: `internal/packagesurface/contract_documents_test.go`
now pins distinct wordings for absent versus present-but-empty, and
`TestUnterminatedFinalLineRoundTripsByteForByte` drives a reference whose last line lacks a
trailing newline. Recorded for the retro: that second row passed as written — the generator already
preserved the unterminated tail, so it pins existing correct behavior rather than fixing a defect.
