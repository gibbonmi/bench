# Citation phase package scope

Status: staged

Roadmap: FT281

Decision source: named reviewed artifact `roadmap/FT281.md`

Verification log: 1 iteration(s) to accept — the review accepted behavior, the amended fence, and the ticket graph; the author corrected one check name.

## Problem

The citation validator records the build tags for each Go test phase. It drops
the package operands from the same phase. A cited test can therefore pass when
its file satisfies a tag set but its package never enters a test phase.

Go excludes package directories such as `testdata` and `_private` from a
recursive package pattern. The current validator does not observe that
selection. It can accept evidence that the gate never compiles.

## Solution

Record one execution entry for each Go test phase. Keep its tags, package
operands, and execution directory together. Do not merge entries only because
their tags match.

Use the Go package loader to expand each entry's package operands. Accept a
citation only when one entry selects the package and accepts the cited file.
The kit does not copy Go's package-pattern or ignored-directory rules.

## User stories

### Group A — a citation matches one real phase

Line: `gpt-5.6-terra` / high. The change repairs a false green in the oracle.

1. As a reviewer, I want one phase to select and compile cited evidence, so that separate phase facts cannot combine falsely.
2. As a maintainer, I want equal tag sets to keep distinct package scopes, so that a later phase does not disappear.
3. As a project owner, I want manifest test phases included, so that a custom gate gets the same citation check.

### Group B — Go owns package selection

Line: `gpt-5.6-terra` / high. Package selection must stay consistent with the installed Go toolchain.

4. As a reviewer, I want a recursive pattern to exclude `testdata` packages, so that fixture tests cannot pose as executed evidence.
5. As a reviewer, I want a recursive pattern to exclude underscore-prefixed packages, so that ignored tests cannot pose as evidence.
6. As a project owner, I want an explicit package operand honored, so that a narrow test phase can supply evidence.
7. As a project owner, I want an absent package list to mean the phase directory, so that the validator matches `go test`.

### Group C — failures stay actionable and compatible

Line: `gpt-5.6-terra` / high. A derivation failure must not restore the false green.

8. As a build delegate, I want a scope violation to name its row and file, so that I can repair the citation.
9. As a reviewer, I want package expansion failures to reject the citation, so that unknown scope cannot authorize evidence.
10. As a spec author, I want existing build-constraint checks preserved, so that package scope does not weaken tag checks.
11. As a non-Go project owner, I want no execution check without a test phase, so that unrelated specs keep their behavior.
12. As a gate owner, I want the package-scope diagnostic to bite through the canary, so that the oracle observes the repair.

## Implementation decisions

- Replace the tag-only census result with one test-execution entry per resolved
  Go test phase. Each entry carries the phase name, tag set, package operands,
  and effective execution directory.
- Keep the phase entries separate. Do not deduplicate entries by tag set.
- Parse the package-list position once from each test phase argv. Stop at
  `-args` or at a test-binary argument, as the Go test grammar requires.
- Use a real `go list` child to expand package operands. Pass the entry's tags,
  execution directory, and `-C` setting on the same terms as the test phase.
- Use the package directories that `go list` returns. Do not implement local
  rules for recursive patterns, `testdata`, dot prefixes, or underscore prefixes.
- Treat an absent package list as `.` in the effective execution directory.
- Accept a cited file only when one execution entry selects its package and
  `go/build.Context.MatchFile` accepts the file for that entry's tags.
- Keep package scope and tag scope paired. Never compare the union of all
  package scopes with the union of all tag sets.
- Report a citation violation when package expansion fails. The report names
  the coverage row and cited file. The gate remains the owner of phase failure.
- Keep an empty census inapplicable. A tree with no Go test phase gains no
  citation-execution claim.
- Move the execution-specific coverage code into one focused source file if
  the existing citation file would grow. Keep the citation grammar in its
  current owner.
- Extend the existing unexecuted-citation canary. Do not add a second fixture
  family or a second owner registry.

## Testing decisions

- Drive the gate census over built-in and manifest phase tables. Assert each
  phase's tags, package operands, and effective directory.
- Drive the coverage package over real temporary Go modules. Use the installed
  Go toolchain instead of a package-pattern mock.
- Put competing facts in two phases. One phase selects the package, and the
  other phase accepts the file's tag. The citation must fail.
- Extend the current coverage-map-validation canary with an ignored package.
  Its expected diagnostic proves the live docs check reaches the new rule.
- The highest seam is `bench coverage --check <spec>`. The shared parse entry
  point and the staged-spec gate sweep consume the same violations.

### Seam diagram

    trigger: bench coverage --check <spec>
        │
        ▼
    cited test path ──▶ [ phase test-execution census + Go package loader ] ──▶ violation or accepted evidence
                           ◀ tests attach here: real fixture module and resolved phase table

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PS1 | 1 | a citation fails when package selection and file constraints succeed in different phases | planned `TestCitationPackageScope` in package `internal/coverage` | a cross-product implementation falsely combines unrelated phase facts |
| PS2 | 2 | two phases with equal tags retain their different package operands | existing `TestExecutedTagCensus` in package `internal/gate` | tag-only deduplication drops one executed package scope |
| PS3 | 3 | a manifest-declared Go test phase contributes its package operands | existing `TestExecutedTagCensus` in package `internal/gate` | a built-in-only census ignores a custom gate phase |
| PS4 | 4 | a citation under `testdata` fails for a phase scoped by `./...` | planned `TestCitationPackageScope` in package `internal/coverage` | tag matching alone accepts a package that Go excludes |
| PS5 | 5 | a citation under an underscore-prefixed directory fails for a phase scoped by `./...` | planned `TestCitationPackageScope` in package `internal/coverage` | a patch for only `testdata` leaves the sibling false green |
| PS6 | 6 | a citation passes when an explicit package operand selects its package | planned `TestCitationPackageScope` in package `internal/coverage` | a copied ignored-directory rule rejects a package that Go selects explicitly |
| PS7 | 7 | a phase with no package operands selects its effective execution directory | planned `TestCitationPackageScope` in package `internal/coverage` | a parser that requires one operand loses Go's default package |
| PS8 | 8 | a package-scope violation names the coverage row and cited file | planned `TestCitationPackageScope` in package `internal/coverage` | an unnamed failure requires a hand package census |
| PS9 | 9 | a package expansion error rejects the cited file | planned `TestCitationPackageScopeFailure` in package `internal/coverage` | an inapplicable fallback restores the false green |
| PS10 | 10 | one phase must both select the package and accept the file constraint | planned `TestCitationPackageScope` in package `internal/coverage` | independent unions weaken the existing build-constraint rule |
| PS11 | 11 | a root with no Go test phase leaves citation execution inapplicable | existing `TestCitationUnexecutedConstraint` in package `internal/coverage` | a fail-closed default breaks non-Go roots |
| PS12 | 12 | the unexecuted-citation canary reports the ignored-package citation | coverage-map-validation canary owner in the `docs-currency-workflow` check | a unit-only diagnostic can stay outside the oracle |

### Edge inventory

- Error path: an invalid package operand rejects the citation with row and file context (PS9).
- Empty input: no package operands select `.` from the effective phase directory (PS7).
- Boundary value: two equal tag sets with different scopes remain two entries (PS2).
- Malformed input: `-args` ends package-operand collection and cannot donate a package (PS3).
- Interrupted state: an interrupted `go list` child yields no partial authorization (PS9).
- Re-run idempotency: repeated checks read the same phase table and produce the same scope result (PS1).
- Process lifecycle: the checker waits for the real `go list` child and consumes its terminal result (PS9).
- Hostile environment: a missing Go executable rejects scope derivation when a test phase exists (PS9).
- Paths with spaces stay distinct argv entries and do not enter a shell (PS6).
- Package patterns with glob characters go to `go list` unchanged (PS4).
- A cited link or special file keeps the existing no-follow refusal (PS10).

**Won't handle** — `-run`, `-skip`, and `-bench` name filters — the ordinary unfiltered test phase remains an in-scope caller.

**Won't handle** — a runtime skip inside `TestMain` or a test body — package and file selection remain the in-scope execution claim.

**Won't handle** — proof that the cited test causes the mapped red — the implementation mutation probe remains the in-scope owner.

## Ownership fences

- `specs/citation-phase-package-scope/`
- `reviews/citation-phase-package-scope.md`
- `CHANGELOG.md`
- `internal/gate/tag_census.go`
- `internal/gate/tag_census_test.go`
- `internal/coverage/citations.go`
- `internal/coverage/citations_test.go`
- `internal/coverage/citation_execution.go`
- `internal/coverage/citation_execution_test.go`
- `tests/canary/coverage-map-validation/unexecuted-tag-citation/`

## Out of scope

- Test-name filter proof for `-run`, `-skip`, and `-bench` is a separate capability: about 12 edits and 4 gate runs.
- Runtime skip proof for `TestMain` and test bodies is a separate capability: about 15 edits and 5 gate runs.
- Causal mutation proof for each citation is a separate capability: about 30 edits and 8 gate runs.

## Further notes

The read sweep found one production consumer of `ExecutedTagCensus`: the
coverage citation validator. The gate test is its only direct test consumer.

The local Go 1.25.14 package documentation and package-search source confirm
that recursive patterns skip `testdata`, dot-prefixed, and underscore-prefixed
directories. The design delegates that decision to `go list`.
