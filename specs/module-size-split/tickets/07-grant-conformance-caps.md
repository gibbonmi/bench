# Grant conformance line caps and consolidate builders

Blocked by: none
Writes: internal/conformance, .bench/structure.budgets

## What to build

Consolidate the six throwaway-root builders in `internal/conformance/` into one shared builder, shrinking the files in place. The six: `writeHostileSkillRoot`, `writeHostileReferenceRoot`, `writeProseBudgetRoot`, `writeLinesRoot`, `recordedScopeRoot`, and `gitInitedRoot`.

Then add one dated `.bench/structure.budgets` row per remaining over-700 conformance test file. The four files: `tier_test.go`, `fixture_bite_test.go`, `docs_workflow_checks_test.go`, and `checks_test.go`. Set each cap at the post-consolidation line count. The comment names this spec and the FT152 precedent. No conformance file splits; the directory file count does not change.

## Acceptance

- [ ] R16: the throwaway-root builder has one shared definition in `internal/conformance/`.
- [ ] R14: `bench structure` lists none of the four conformance test files.
- [ ] R15: `.bench/structure.budgets` holds one dated row naming each of the four files.
- [ ] R19: each granted cap is at most 20 newlines above the file's count at the ticket tip.
- [ ] R08: `go test -list '.*' ./internal/conformance/` emits the same test-name set at base and tip.
- [ ] R18: `bench gate` exits zero before the commit.
