# citation-execution-proof

Status: staged

Roadmap: FT133

Decision source: roadmap/FT133.md plus the recommended sequence in ROADMAP.md (reviewed 2026-08-30 drain 2238d8fa). Commit 926af243 delivered the light-path half; this spec covers the remainder.

Verification log: 2 iteration(s) to accept — the author folded the round's named B2 repair after the cap: the census reads the resolved phase table. Amendment round (uncited-row report): 2 iteration(s) to accept — the round caught the missing cited-row exclusion, folded as row CE27.

## Problem

`bench coverage --check` resolves a cited test name to a declared function. It
does not prove that the gate executes that function. A citation into a
stress-tagged file resolves and greens, but the gate never runs that file. A
backticked test path without a name list grades nothing, so a mention poses as
evidence. A renamed subtest keeps its citation green, because only the parent
function resolves.

A map with mixed row-ID tags escapes the preflight
membership check, which scopes to one tag. The shipped-surface check derives a `package.json` exclusion from each
kit-only allowlist row. It never proves that the row's source path exists,
so a misspelled row passes.

## Solution

Extend the one coverage validator so a citation must point at a test the gate
actually executes on this host. Derive the executed tag census from the same
phase table the gate runner executes, so the census cannot drift. Grade a
mention, a stale subtest segment, and a mixed-tag row-ID set as violations.
Make the shipped-surface conformance check stat every allowlist source path.
Prove each new diagnostic with a canary fixture in the
coverage-map-validation family.

Add one informational report to `--check`: name each mapped row whose seam
cell holds no citation, unless the seam cell carries the `review-owned:`
prefix. The report adds no violation, so a staged spec stays green and a
build proves a new row is wired in one call.

## User stories

### Group A — a citation proves an executed test

Line: opus / medium. The scorecard routes gate and conformance logic to Opus at medium effort.

1. As a reviewer, I want a violation when no executed tag set satisfies a cited file's constraints, so that a citation proves an executed test.
2. As a build delegate, I want the violation to name the row, file, and constraint, so that I repair it without a hand census.
3. As a maintainer, I want the census derived from the phase table the runner executes, so that the census and the oracle cannot drift.
4. As a kit developer, I want a citation into the system suite accepted, so that a system-tagged test stays valid acceptance evidence.
5. As a reviewer, I want a citation to a stress-tagged file refused, so that release-only evidence cannot green a staged spec.
6. As a reviewer, I want a cited file with a foreign GOOS or GOARCH constraint refused, so that a foreign-OS test cannot serve as evidence.
7. As a reviewer, I want a cited non-regular path refused without an open, so that a planted FIFO cannot wedge the check.

### Group B — a mention is not a citation

Line: opus / medium. The scorecard routes gate and conformance logic to Opus at medium effort.

8. As a reviewer, I want a seam-cell test token without a name list to red, so that a mention cannot pose as evidence.
9. As a spec author, I want a test path outside the seam cell left ungraded, so that honest prose in other cells does not red.

### Group C — a subtest citation stays resolvable

Line: opus / medium. The scorecard routes gate and conformance logic to Opus at medium effort.

10. As a reviewer, I want a cited subtest segment resolved against the file's `t.Run` literals, so that a renamed subtest turns the check red.
11. As a spec author, I want a file with a non-literal `t.Run` name exempt from segment resolution, so that a table-driven suite cannot false-red.

### Group D — one row-ID tag per map

Line: opus / medium. The scorecard routes gate and conformance logic to Opus at medium effort.

12. As a reviewer, I want a refusal when row IDs carry more than one tag, so that the membership check covers every declared row.
13. As a build delegate, I want that violation to name each tag it found, so that I see which rows disagree.

### Group E — every allowlist source exists

Line: opus / medium. The scorecard routes gate and conformance logic to Opus at medium effort.

14. As a release owner, I want the shipped-surface check to stat every consumer-payload allowlist source path, so that a misspelled row cannot pass vacuously.
15. As a maintainer, I want that red to name the row's source and audience, so that the repair starts at the row.

### Group F — the oracle observes the new checks

Line: opus / medium. The scorecard routes gate and conformance logic to Opus at medium effort.

16. As a reviewer, I want a canary fixture for each new coverage diagnostic, so that the diagnostic reaches a real gate run.
17. As a spec author, I want the historical marker to keep silencing every coverage check, so that the opt-out stays whole.
18. As a preflight consumer, I want the new violations surfaced through the shared parse entry point, so that the review preflight and the gate agree.

### Group G — an uncited row is named

Line: opus / medium. The scorecard routes gate and conformance logic to Opus at medium effort.

19. As a build delegate, I want `--check` to name each row with no seam-cell citation, so that I prove a new row is wired.
20. As a spec author, I want the `review-owned:` prefix to exempt a seam cell from the report, so that a deliberate exclusion adds no noise.
21. As a preflight consumer, I want the uncited report kept out of the violation list, so that a staged spec with planned-test prose stays green.
22. As a build delegate, I want the report to name an opted-in row by its row ID, so that the name survives a row reorder.

## Implementation decisions

- The executed tag census is one exported derivation in the gate package. It
  parses each test phase's argv from the resolved phase table. That table is
  the root's manifest when the root declares one, else the built-in kit
  table. The runner executes the same resolution.
  It collects the `-tags` sets those argvs carry, plus the untagged default
  set. No second census copy exists.
- The census root is the base the
  citations resolve against, and the kit resolves through the gate's own kit
  rule. The cited files and the census come from one tree.
- The coverage package imports the gate package for the census. The dependency
  is acyclic today; the gate package imports neither the coverage package nor
  the spec package.
- Build-constraint evaluation uses the standard `go/build/constraint` package
  over the cited file's `//go:build` line, plus the GOOS and GOARCH filename
  suffix rule. A term satisfies from a census set or from the host GOOS or GOARCH. A
  release tag and a toolchain-implied tag (`unix`, `gc`, `cgo`) also
  satisfy. The release tags come from `go/build.Default.ReleaseTags`. The
  suffix rule strips `_test` before it reads the suffix. A citation passes
  when any executed set satisfies it.
- A graded root with no test phase leaves the execution check inapplicable, so
  a non-Go root keeps its current behavior.
- The mention rule, the subtest rule, and the constraint rule extend the one
  citation grammar in the coverage package. The mixed-tag refusal extends the
  one row-ID validator there. No consumer re-derives map structure.
- The allowlist existence check lives beside the derived-exclusion loop in the
  shipped-surface conformance check. A tree row must name a directory, and a
  file row must name a regular file, on both audiences.
- Cited paths are classified before an open, with `bounds.ClassifyNoFollow`,
  so a FIFO or a link reds instead of blocking.
- The uncited report is one row classification beside the citation grammar. A
  row is cited when its seam cell holds at least one citation. A row is exempt
  when its trimmed seam cell starts with the `review-owned:` prefix. Every
  other mapped row is uncited. That prefix is the one marker grammar, and the
  delivered no-citation test already uses it. The match is literal and
  case-sensitive on the lowercase form.
- The report renders in the `--check` arm of the coverage command, beside the
  pass line. It joins no violation list, and the exit code does not change.
  An opted-in map's report names the row ID; a non-opt-in map's report names
  the row number. The report is one bounded line: the row count plus the row
  names.

## Testing decisions

- A good test drives `bench coverage --check` or the shared parse entry point
  over a written spec fixture and asserts the exact violation text.
- Seams with prior art: the coverage package's citation tests, its command
  tests, the gate package's phase-table tests, the shipped-surface
  conformance test, and the coverage-map-validation canary family.
- The gate observes the feature through the staged-spec sweep in the
  docs-currency conformance family and through the canary meta-gate.

### Seam diagram

    trigger: bench coverage --check <spec>, the staged-spec sweep, preflight bootstrap
        │
        ▼
    spec.md + cited files  ──▶  [ coverage CheckFiles + gate tag census ]  ──▶  violations
                                    ◀ tests attach here: write a spec fixture and a
                                      tagged test file, then assert the violation text

    trigger: gate conformance run
        │
        ▼
    .bench/consumer-payload.json  ──▶  [ shipped-surface check ]  ──▶  per-row absence red
                                    ◀ tests attach here: fixture root with a misspelled row

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CE1 | 1 | a citation to a file whose build constraints no executed tag set satisfies is a violation | new test TestCitationUnexecutedConstraint in package internal/coverage | a resolved but never-run citation greens today |
| CE2 | 2 | the unexecuted-constraint violation names the row number, the cited file, and the constraint | the same coverage test asserts the message fields | an unnamed red forces a hand census |
| CE3 | 3 | the census on the kit root holds the untagged set and the `system` set | new test TestExecutedTagCensus in package internal/gate | a second census copy drifts from the oracle |
| CE4 | 4 | a citation into the system suite passes on the kit root | coverage test with a system-tagged fixture file | a fail-closed census bans real system evidence |
| CE5 | 5 | a citation to a stress-tagged fixture file is a violation | coverage test with a stress-tagged fixture file | release-only evidence greens a staged spec today |
| CE6 | 6 | a cited file with a foreign GOOS filename suffix is a violation | coverage test with a windows-suffixed fixture file | a foreign-OS test never executes on the gating host |
| CE7 | 7 | a cited path that is not a regular file is a violation without an open | coverage test with a planted FIFO | an open on a FIFO wedges the check inside the oracle |
| CE8 | 8 | a backticked test path in a seam cell without a name list is a violation | new test TestMentionIsNotACitation in package internal/coverage | a mention poses as evidence today |
| CE9 | 9 | a backticked test path outside the seam cell is not graded | the same coverage test writes the token into the why cell | an overreaching rule reds honest prose — review-owned anti-regression |
| CE10 | 10 | a cited subtest segment absent from the file's `t.Run` string literals is a violation | new test TestSubtestSegmentResolves in package internal/coverage | a renamed subtest stays green today |
| CE11 | 11 | a file with a non-literal `t.Run` name is exempt from segment resolution | the same coverage test adds a dynamic `t.Run` call | a table-driven suite would false-red — review-owned anti-regression |
| CE12 | 12 | a map whose row IDs carry two tags is a violation | new test TestMixedTagRowIDs in package internal/coverage | tag-scoped membership skips foreign-tag rows today |
| CE13 | 13 | the mixed-tag violation names each tag found | the same test asserts both tags in the message | an unnamed red hides which rows disagree |
| CE14 | 14 | an allowlist row whose source path is absent reds the shipped-surface check | new test TestAllowlistSourceExists in the shipped-surface conformance file | a misspelling mirrored in the payload and `package.json` passes today |
| CE15 | 15 | the absence red names the row's source and audience | the same test asserts the message fields | the repair must start at the row |
| CE16 | 16 | a mixed-tag canary fixture reds a real gate run | new fixture in the coverage-map-validation canary family | a diagnostic that never reaches the oracle proves nothing |
| CE17 | 16 | an unexecuted-tag canary fixture reds a real gate run | new fixture in the coverage-map-validation canary family | the constraint path must bite through the sweep |
| CE18 | 17 | the historical marker silences the new checks | coverage test over a historical spec with a bad citation | a partial opt-out breaks the documented contract — review-owned anti-regression |
| CE20 | 1 | a malformed `//go:build` expression in a cited file is a violation naming the file | the same coverage test asserts the parse-failure arm | a bad expression must red rather than pass silently |
| CE21 | 1 | a graded root with no test phase leaves the execution check inapplicable | coverage test over a fixture root without a test phase | a fail-closed default would red every non-Go root |
| CE22 | 3 | a cited file behind a manifest-declared custom tag passes | new manifest fixture case in TestCitationUnexecutedConstraint | a hardcoded census copy cannot see a manifest tag |
| CE19 | 18 | the shared parse entry point returns the mixed-tag, mention, subtest, and unexecuted-constraint classes | new assertion in the coverage package on the exported entry point | the preflight and the gate would otherwise disagree |
| CE23 | 19 | `--check` on a green non-opt-in map names each uncited row by its row number | new test TestUncitedRowReport in package internal/coverage | a build cannot prove a new row is wired without a hand read |
| CE24 | 20 | a seam cell with the `review-owned:` prefix is absent from the uncited report | the same command test adds a marked row | an unexempt marker turns every deliberate exclusion into noise |
| CE25 | 21 | the uncited report leaves the check's exit code at zero | the same command test asserts the exit code | a fail-closed report would red every staged spec |
| CE26 | 22 | an opted-in map's uncited report names the row by its row ID | the same command test drives an opted-in fixture | a positional number goes stale on a row reorder |
| CE27 | 19 | a row whose seam cell holds a citation is absent from the uncited report | the same command test cites a declared test file | a report that names every mapped row would pass the sibling rows |

### Edge inventory

- A malformed `//go:build` expression in a cited file reds and names the file (row CE20).
- A graded root with no test phase leaves the execution check inapplicable (row CE21).
- A cited file with no build constraint always satisfies the census (CE4 covers the untagged arm).
- An empty parenthesized list stays a non-citation, unchanged (existing behavior, pinned by the delivered tests).
- A spec path with spaces resolves through the existing resolver (delivered behavior, unchanged).
- A historical spec emits no uncited report; the historical pass line stays unchanged (the historical arm of row CE23's test, at the command seam).
- A fully cited map reports nothing beside the pass line (the green arm of row CE23's test).
- A check with violations emits no uncited report, because the violation arm returns first (row CE23 scopes the report to a green map).
- The report's exact wording is build-owned; the rows pin the named fields, and the build's test pins the string.

**Won't handle** — a citation valid only on another operating system — the census grades the gating host, and cross-OS evidence has no oracle here.
**Won't handle** — a dynamic subtest name that a literal scan cannot see — the exemption in story 11 keeps the suite citable by function name.
**Won't handle** — proof that a cited test asserts the row's exact diagnostic — the build's mutation probe owns that judgment.
**Won't handle** — an in-test skip with an executed body — the gate's capability-skip log already reds an environment-class skip.
**Won't handle** — golden and fixture inventories behind an "already covered" claim — the review axes own that judgment per the review discipline.

## Ownership fences

- `specs/citation-execution-proof/`
- `reviews/citation-execution-proof.md`
- `internal/coverage/`
- `internal/gate/phases.go`
- `internal/gate/tag_census.go`
- `internal/gate/tag_census_test.go`
- `internal/conformance/package_shipped_surface_test.go`
- `internal/conformance/registry_test.go`
- `internal/preflight/gather.go`
- `tests/canary/coverage-map-validation/`

## Out of scope

- Mechanical sole-red verification, which runs a mutation against each cited test and requires the red — about 30 edits, 8 gate runs.
- Allowlist mode enforcement, which compares each row's declared mode with the on-disk permission — about 6 edits, 2 gate runs.
- AST-level subtest name resolution for composed names — about 10 edits, 3 gate runs.

## Further notes

Commit 926af243 (2026-08-30) delivered the resolution half: citation
resolution against declared functions, the shared fence parser, and the
review-pickup fence member. The 2026-08-18 roadmap annotation closed the
capability-skip half inside the gate runner. This spec closes the rest of
FT133 and the FT103 absence half that joined it. The 2026-08-31 sign-off
pulled the uncited-row report into scope, so this spec also disposes of the
FT133 line "names the rows that no test names".

The fixture inventory asset under decisions/assets records two
coverage-map-validation fixtures. Three exist on disk, and two more land
here. That snapshot is not gate-enforced, and the drift is noted for a
later drain.
