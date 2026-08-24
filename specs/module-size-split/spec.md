# Module size split

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-23 — scope is every source file over 700 lines; the reviewer authorized the landing in the same conversation. Amended by reviewer-confirmed conversation, 2026-08-24 — ticket 13's target already conforms, so the reviewer widened the scope by one file. The widened file is `internal/lines/lines_test.go`.

Verification log: 1 iteration(s) to accept — the round returned one blocking ticket-graph finding and six advisory items. Amendment round, 2026-08-24: 1 iteration(s) to accept — the round returned two blocking consistency findings, folded before staging. The blocking finding: the worktree ticket exceeded one context window. The author split that ticket and the two bundled multi-package tickets. The author also un-bundled row R04, added row R19, and folded the prose fixes before staging.

## Problem

`bench structure` reports 75 issues. Twenty flagged source files hold more than 700 lines each; the worst holds 2210. A twenty-first file over 700 lines, `internal/status/status.go`, stands under a reviewer accept and stays out of scope. A reader cannot hold one of these files in one pass, and a fresh session pays for every line it loads. The debt is advisory today, but it grows with each change.

## Solution

Split the twenty files along responsibility lines so that each result file holds at most 400 lines. Keep every split inside its Go package, so no exported symbol and no call site changes. Where reviewer precedent forbids a split (`internal/conformance/`, FT152), record a dated ratchet grant instead, and remove the duplicated fixture builders in place. No directory gains a crowding violation it does not already have.

## User stories

### Group A — production files split along responsibility

Line: gpt-5.6-terra / medium. The moves are mechanical, but a dropped symbol breaks the build tree-wide, so the mid tier applies.

1. As a maintainer, I want `internal/landing/landing.go` split into composition, spec-tree editing, attribution, and git-primitive files, so that each file reads as one responsibility.
2. As a maintainer, I want `internal/git/git.go` split into process-primitive, worktree-admin, and status files, so that the worktree-admin surface stands alone.
3. As a maintainer, I want `internal/diff/diff.go` split into command-and-render, snapshot-plumbing, and range-and-movement files, so that the identity logic is findable.
4. As a maintainer, I want `internal/freshness/freshness.go` split into seal-core, publish, build-input, and verify files, so that the publish lifecycle stands alone.
5. As a maintainer, I want `internal/worktree/land.go` split into first-run, resume, identity, and refusal files, so that each landing flow reads on its own.
6. As a maintainer, I want every file the split creates or retains to hold at most 400 lines, so that `bench structure` stops flagging them.
7. As a caller of these packages, I want every exported symbol unchanged, so that no call site outside the package needs an edit.

### Group B — test files split by family

Line: gpt-5.6-terra / medium. Test moves risk silent test loss, so the mid tier applies with the test-name-set row as the oracle.

8. As a maintainer, I want `internal/landing/landing_test.go` split into land, reviewed-land, and shared-helper files, so that each flow's tests sit together.
9. As a maintainer, I want `internal/git/git_test.go` split into helper, worktree-admin, and fact-family files, so that the admin matrix stands alone.
10. As a maintainer, I want `internal/status/status_test.go` split by signal family with one shared fixture file, so that each family reads in one pass.
11. As a maintainer, I want `internal/freshness/freshness_test.go` split into digest, publish, and verify files, so that the subprocess helpers sit with the publish tests.
12. As a maintainer, I want `internal/worktree/land_test.go` and `internal/worktree/eligibility_test.go` split by journey family with one shared fixture file, so that the 2210-line file disappears. (Done-by-tree note, 2026-08-24: `eligibility_test.go` already conforms at 209 lines, so no ticket splits it.)
13. As a maintainer, I want `internal/roadmap/context_test.go` split into command, row-selector, and occurrence files, so that the ledger grammar tests stand alone.
14. As a maintainer, I want `internal/maps/maps_test.go` split into command, parse, and graph files, so that each validation family sits together.
15. As a maintainer, I want `internal/systemtest/owner_test.go` split with the harness and `TestMain` kept whole, so that the sibling test files keep compiling.
16. As a maintainer, I want `internal/skillsindex/skillsindex_test.go` split into render, allowlist, and reference files, so that the safety tests stand alone.
17. As a maintainer, I want `internal/coverage/coverage_test.go` split into parse, command, and schema files, so that the fixtures live beside the parse tests.
18. As a maintainer, I want every shared test helper moved once into one file per package, so that no helper is pasted twice.
19. As a reviewer, I want the same test-name set per package before and after each split, so that no test is silently dropped.
27. As a maintainer, I want `internal/lines/lines_test.go` split into parse, verdict-resolution, agent-line, and fixture files, so that each family reads in one pass. (The 2026-08-24 amendment added this story.)

### Group C — structure budgets respected

Line: gpt-5.6-terra / medium. Budget edits are reviewer-owned data; the mid tier applies.

20. As a reviewer, I want no new crowding violation in any directory, so that a line breach is not traded for a directory breach.
21. As a reviewer, I want the total `bench structure` violation count to drop to at most 55, so that the change is a net structural gain.
22. As a reviewer, I want dated ratchet rows for the four conformance files per the FT152 precedent, so that the one-signal family stays whole.
23. As a maintainer, I want the six per-file throwaway-root builders in `internal/conformance/` collapsed into one shared builder, so that the harness knowledge has one source.

### Group D — behavior preserved

Line: gpt-5.6-terra / medium. The oracle work is verification, not construction; the mid tier applies.

24. As a reviewer, I want `go build ./...` and the touched package tests green after every ticket, so that each landing stands alone.
25. As a maintainer, I want the freshness publication-topology test green after the freshness split, so that no new `Publish` caller and no new `package main` appears.
26. As a reviewer, I want `bench gate` green at every ticket commit, so that the repo never holds a red intermediate state.

## Implementation decisions

- Every split is a same-package file move. No symbol changes visibility, name, or signature. No logic changes.
- The directory file budget is 12 unless `.bench/structure.budgets` raises it. Each ticket plans its final file count against that number:
  - `internal/landing/`: the fingerprint group folds into the retained `landing.go`; the directory lands at 12 files.
  - `internal/git/`: the branch-lifecycle group folds into the retained `git.go`; the directory lands at 12 or fewer files.
  - `internal/diff/`: `diff.go` becomes exactly three files; the directory lands at 12 files.
  - `internal/roadmap/`: `context_types.go` folds into `context_parse.go` in the same ticket; the directory lands at 12 files.
  - `internal/worktree/` already carries its one crowding violation at 51 files against a granted 18. The split raises the directory's file count, but the `DIR CROWDED` line count stays at one, and the follow-up subpackage extraction attacks it. This call is flagged for reviewer veto.
- `internal/conformance/` takes no split. Ticket 7 adds one dated `.bench/structure.budgets` row per over-700 file, following the `docs_workflow_helpers_test.go` precedent, and consolidates the duplicated root builders. Only ticket 7 may touch `.bench/structure.budgets`.
- A ticket that cannot meet its directory budget without breaking a 400-line cap exits and reports a material acceptance shortfall. It does not add a budgets row and does not land a partial.
- `internal/systemtest/` files all carry the `//go:build system` tag, and exactly one `TestMain` remains.
- The delegate research reports fix the per-file target layouts. A ticket may regroup within its package when a measured line count differs, under the coverage rows.

## Testing decisions

- The external behavior is the structure report and the unchanged package behavior. A good check drives `bench structure`, `go build`, and `go test` — never a human judgment of "looks smaller".
- The seams are the existing CLIs and the Go toolchain. No new test seam is introduced; a pure move needs no new test.
- One named gap stays human-owned: a mechanical split at an arbitrary line that ignores responsibility passes every automated row. The implementation review, not a row, judges the grouping against the story's named responsibilities.
- The gate seam observes every ticket through `bench gate` on the worktree before its commit.

### Seam diagram

    trigger: ticket lands a split
        │
        ▼
    repo tree  ──▶  [ bench structure scan ]  ──▶  violation list
                      ◀ tests attach here: run `bench structure`, assert the target file is absent
    repo tree  ──▶  [ go build / go test -list ]  ──▶  build result, test-name set
                      ◀ tests attach here: compare the name set before and after the move

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| R01 | 1, 2, 3 | `bench structure` lists none of `internal/landing/landing.go`, `internal/git/git.go`, `internal/diff/diff.go` | bench structure CLI | a rename-only or partial split leaves the file over 400 and listed |
| R02 | 4, 5 | `bench structure` lists neither `internal/freshness/freshness.go` nor `internal/worktree/land.go` | bench structure CLI | a partial split leaves the file listed |
| R03 | 6 | every file a ticket creates counts at most 400 newlines in the structure scan | bench structure CLI | a mechanical halving that leaves one file at 450 stays listed |
| R04 | 7, 24 | `go build ./...` exits zero after each ticket | go toolchain | a dropped or renamed exported symbol breaks a dependent package |
| R05 | 8, 9, 10, 11 | `bench structure` lists none of `landing_test.go`, `git_test.go`, `status_test.go`, `freshness_test.go` | bench structure CLI | a partial test split leaves the file listed |
| R06 | 12, 13 | `bench structure` lists none of `land_test.go`, `eligibility_test.go`, `context_test.go` | bench structure CLI | a partial test split leaves the file listed |
| R07 | 14, 15, 16, 17 | `bench structure` lists none of `maps_test.go`, `owner_test.go`, `skillsindex_test.go`, `coverage_test.go` | bench structure CLI | a partial test split leaves the file listed |
| R08 | 19 | `go test -list '.*'` per touched package emits the same test-name set at the ticket's base and tip | go toolchain | a test deleted to shrink a file vanishes from the list |
| R09 | 15 | `go vet -tags system ./internal/systemtest/...` exits zero after the systemtest split | go toolchain | a new file without the `system` tag or a broken harness fails vet |
| R10 | 15 | exactly one `TestMain` definition remains under `internal/systemtest/` | rg sweep | a duplicated `TestMain` breaks the package under the system tag |
| R11 | 18 | each moved test helper has exactly one definition in its package at the ticket tip | rg sweep | a copy-pasted helper leaves two definitions |
| R12 | 20 | the final `bench structure` output holds no `DIR CROWDED` line for a directory clean at the base commit | bench structure CLI | a split that overfills a directory prints a new crowding line |
| R13 | 21 | the final `bench structure` total is at most 55 issues | bench structure CLI | a batch that fixes fewer files than scoped stays above 55 |
| R14 | 22 | `bench structure` lists none of the four over-700 `internal/conformance/` test files | bench structure CLI | a missing or short grant row leaves the file listed |
| R15 | 22 | `.bench/structure.budgets` holds one dated row naming each of the four conformance files | file content check | an undated or absent row breaks the reviewer-owned grant convention |
| R16 | 23 | the throwaway-root builder has one shared definition in `internal/conformance/` at the ticket tip | rg sweep | a consolidation that keeps a second builder leaves two definitions |
| R17 | 25 | `go test ./internal/freshness/...` exits zero after the freshness split | go toolchain | a new `Publish` caller or `package main` fails the topology test |
| R18 | 26 | `bench gate` exits zero at each ticket commit | bench gate CLI | a red intermediate state cannot commit |
| R19 | 22 | each granted conformance cap is at most 20 newlines above the file's count at the ticket tip | file content check | an inflated cap hides all future growth from the scan |
| R20 | 27 | `bench structure` no longer lists `internal/lines/lines_test.go` | bench structure CLI | a partial test split leaves the file listed |

### Edge inventory

- Cross-file test fixtures: each package's shared helpers move to one fixtures file; R11 observes the single definition.
- `TestMain` uniqueness and the `system` build tag: R09 and R10 observe them.
- The freshness AST topology test: R17 observes it; the split adds no `Publish` caller.
- Directory budgets at exactly 12: a count of 12 passes, 13 fails; the per-directory plans in the implementation decisions pin the counts.
- Hostile-input checklist: the change adds no input surface. The only new parsed content is `.bench/structure.budgets` rows, and `loadBudgets` already drops malformed rows fail-soft. No tuned profile addition is needed.

**Won't handle** — files between 400 and 700 lines — the reviewer scoped this spec to over-700 files, except `internal/lines/lines_test.go` (amended 2026-08-24). The flagged remainder survives for a follow-up spec, and every in-scope caller keeps working.
**Won't handle** — subpackage extraction for `internal/worktree/` eligibility or landing — the extraction must export private symbols and risks an import cycle. The in-package split keeps every caller green.
**Won't handle** — files under standing accepts (`internal/status/status.go`, `internal/releaseevidence/`, `internal/releasepreflight/`) — the accept entries are reviewer decisions that stand.
**Won't handle** — the pre-existing `DIR CROWDED` lines for `internal/adopt/`, `internal/gate/`, `internal/worktree/` — they predate this spec, and R12 pins that no new one appears.

## Ownership fences

- `specs/module-size-split/` — the spec author.
- `internal/landing/` — ticket 01.
- `internal/git/` — ticket 02.
- `internal/diff/` — ticket 03.
- `internal/freshness/` — ticket 04.
- `internal/status/` — ticket 05.
- `internal/worktree/` — ticket 06.
- `internal/lines/` — ticket 13 (re-fenced by the 2026-08-24 amendment).
- `internal/conformance/` — ticket 07.
- `.bench/structure.budgets` — ticket 07 only.
- `internal/roadmap/` — ticket 08.
- `internal/maps/` — ticket 09.
- `internal/systemtest/` — ticket 10.
- `internal/skillsindex/` — ticket 11.
- `internal/coverage/` — ticket 12.

## Out of scope

- Split the 400-to-700-line remainder (about 50 files) — separate capability; roughly 50 edits, 15 gate runs.
- Extract `internal/worktree/` eligibility into a subpackage and shrink the 51-file directory — separate capability; roughly 12 edits, 3 gate runs.
- Wire `bench structure` into the gate oracle — separate capability; roughly 4 edits, 2 gate runs.

## Further notes

- `bench structure` is advisory today; the gate does not run it. The rows therefore bind to the CLI's output, and R18 keeps the real oracle green.
- The seven files over their own ratchet grants (for example `internal/gate/verdict.go` at 673 against 431) sit under 700 lines and stay for the follow-up spec.
