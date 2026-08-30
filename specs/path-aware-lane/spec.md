# Path-aware lane

Status: staged

Roadmap: FT215

Decision source: named reviewed artifact `roadmap/FT215.md`, the row body the 2026-08-30 drain settled at `09e1f15b`.

Verification log: 2 iteration(s) to accept — round 1 (codex `gpt-5.6-sol` / high) returned 7 blocking and 4 non-blocking findings. Round 2 verified the folds and left one partial. The author folded it as PL47 and two citations; see Further notes.

## Problem

A worktree `bench commit` runs the kit's fast lane as four fixed checks. A
Markdown-only commit pays `go vet` and `go build`, and a Go-only commit pays
the prose check. The lane grades the Markdown among the *named* paths, so a
spec named by its folder reaches no prose check. On 2026-08-28 the lane passed
a spec commit that the landing gate then failed on nine sentences. A roadmap,
decision-map, or retro defect reaches no lane check at all and costs a full
gate run at the landing.

## Solution

`bench commit` derives its change list from the composed tree. The list is the
raw diff between the expected base tree and the composed tree, with no rename
detection. A named directory therefore expands to its changed files, and a
rename becomes one deletion and one addition. The attribution owner still
refuses a special or unreadable named path before any check runs.

Each composed change carries one or more path classes. Each class selects
checks from the kit's declared lane. Go source selects gofmt, vet, and build.
Go metadata or an embed input selects vet and build. Markdown or the prose
policy selects prose.

A document family with a registered check selects that check through
`bench test --check <name>`. Mixed classes select the union in declared
order. An unknown path, a symbolic link, or a gitlink selects every declared
check. The lane line names the selected checks and their classes.

A lane check carries the name of a gate check only when it runs the same bound
over the same bytes. The gate counterpart table in Further notes names each
check's counterpart and the mutation that reds both. A manifest lane in a
linked repo runs as declared, with no selection.

## User stories

### Group A — the composed change list

Line: opus / medium. The derivation replaces the authority's two inputs with
one, and it sits under the commit and the merge. A wrong list grades the wrong
bytes.

1. As an agent, I want the lane to grade the base-to-composed tree diff, so that it grades what the commit composes.
2. As an agent, I want a named directory to expand to its changed files, so that a spec named by folder reaches prose.
3. As an agent, I want a rename to reach the lane as one deletion and one addition, so that only the new path is graded.
4. As an agent, I want a deleted Markdown file to leave the prose subject, so that the check grades only files the tree holds.
5. As an agent, I want a path with a space and a non-ASCII byte carried verbatim, so that the prose check grades the real file.
6. As a reviewer, I want attribution to refuse a special or unreadable path before any check runs, so that a FIFO reaches no check.
7. As a reviewer, I want `bench worktree merge` to derive the same change list from the previous tip, so that one derivation serves both callers.
8. As a reviewer, I want the lane authority to hold no stated path list, so that the change list has one source.

### Group B — path classes and selection

Line: opus / medium. The selection is gate logic, and an over-narrow selection
lets a red commit pass onto the worktree branch.

9. As an agent, I want a Go source change to select gofmt, vet, and build, so that a Go commit pays no prose check.
10. As an agent, I want a `go.mod` or `go.sum` change to select vet and build, so that a dependency edit is compiled and not formatted.
11. As an agent, I want a `//go:embed` target change to select vet and build, so that a changed embedded asset is compiled in.
12. As an agent, I want a Markdown change to select prose only, so that a prose commit pays no Go toolchain.
13. As an agent, I want a `.bench/prose-exclusions` change to select prose, so that a malformed exclusion row reds at the commit.
14. As an agent, I want mixed classes to select the union of their checks in declared order, so that each relevant check runs once.
15. As an agent, I want an unknown path to select every declared check, so that a shell or JSON change keeps today's lane.
16. As a reviewer, I want a symbolic link on either side to select every declared check, so that a link never narrows the lane.
17. As a reviewer, I want a gitlink to select every declared check, so that a submodule pointer never narrows the lane.
18. As an agent, I want the lane line to read `lane{outcome=pass,checks=<names>,classes=<classes>}`, so that I read why each check ran.
19. As an agent, I want a manifest lane to run as declared with no `classes=` cell, so that a project's own lane keeps its meaning.
20. As a reviewer, I want the kit lane marked selective and a manifest lane marked not, so that the selection has one switch.
21. As an agent, I want `bench commit --dry-run` to select the same checks, so that a dry run predicts the real run.
22. As a reviewer, I want the lane record unchanged and the gate cache untouched, so that a selected lane is still never green.
23. As a reviewer, I want a selection to be a declared-order subset of the lane, so that no path adds an undeclared check.
24. As a reviewer, I want one exported embed-target derivation that the pack-asset list also composes, so that the embed rule has one source.
39. As a reviewer, I want the gate's phase table unchanged, so that the landing stays the one full grade.
40. As a reviewer, I want a lane check named after a gate check to run that bound, so that a name never overstates a grade.

### Group C — document families

Line: opus / medium. The rows bind the lane to the registry, and a wrong
binding runs a check on the wrong family or misses one.

25. As an agent, I want a `ROADMAP.md` or `roadmap/` change to select `roadmap-detail-integrity`, so that a broken row reds at the commit.
26. As an agent, I want a `decisions/` or `specs/<slug>/decisions/` change to select `decision-map-integrity`, so that a broken map reds at the commit.
27. As an agent, I want a `capture/retros/` change to select `retro-improvement-markers`, so that a retro paragraph reds at the commit.
28. As an agent, I want a `projects/benchkit.md` change to select `guidance-prose-budgets` and `profile-lane-table`, so that a stale profile reds at the commit.
29. As a reviewer, I want each family's checks derived from the registry's `Inputs` binding, so that the family-to-check fact has one source.
30. As a reviewer, I want a document check to run as `bench test --check <name>` through the lane's run binary, so that the lane builds no second executable.
31. As an agent, I want a roadmap Markdown change to select prose and `roadmap-detail-integrity`, so that both rules run.
32. As an agent, I want a red document check to refuse the commit as `lane{outcome=fail,check=<name>}`, so that the refusal names the rule.

### Group D — advertisement and guidance

Line: opus / medium. The profile table is gate-pinned, and the glossary and
the ADR compound through every cold session.

33. As a cold agent, I want the profile lane table to carry a `selected by` column and the document rows, so that it matches the lane.
34. As a reviewer, I want the `profile-lane-table` check to red a stale `selected by` cell, so that the advertisement cannot drift.
35. As a cold agent, I want the glossary to define `composed change` and `path class`, so that the vocabulary does not drift.
36. As a reviewer, I want ADR 0017 to record that the composed changes select declared checks, so that the decision reads current.
37. As a reviewer, I want one `### Changed` changelog entry, so that the release notes carry the behavior change.
38. As a cold agent, I want the `fast lane` glossary entry to describe the selection, so that the term reads current.

## Implementation decisions

- A **composed change** is one entry of `git diff --raw --no-renames -z <base>^{tree} <tree>`. It carries the status letter, the source mode, the destination mode, and the path. The `-z` framing keeps a path's own bytes.
- `gate.ComposedChanges(root, base, tree)` is the one derivation. The lane authority calls it once, from the base commit and the composed tree, and passes the list in `LaneRequest.Changes`. `LaneRequest` loses `NamedMarkdown`. `LaneAuthority` loses `NamedMarkdown` and `PreviousTip` and gains `Base`. `bench commit` passes its expected base commit, and `bench worktree merge` passes the previous tip.
- `RunLane` resolves the prose placeholder from `Changes`: the Markdown whose destination mode is a regular file. A deletion contributes no subject.
- The attribution owner runs before composition, as today. A special or unreadable named path refuses with the existing message, and no lane run starts.
- `gate.LaneForCommit` returns a `gate.Lane` value: the checks, the source root, and `Selective`. `Selective` is true for the kit's built-in lane and false for a manifest lane.
- The `Lane` type, `LaneForCommit`, `ComposedChanges`, the class table, and `SelectLane` live in one new file beside the lane run. The lane run's own file is already over the structure budget, so it does not grow.
- The worktree `mergeLane` join returns the same `gate.Lane` value, so the merge and the commit read one lane shape.
- A **path class** is one row of a kit-owned table. The table binds each class to a path predicate and to check names. Its rows are `go-source`, `go-build-input`, `markdown`, `prose-policy`, `roadmap-board`, `decision-documents`, `capture-retros`, and `benchkit-profile`. A path matches every class whose predicate holds, and a path that matches none, a symbolic link mode, or a gitlink mode is `unknown`.
- The class predicates are:
  - `go-source`: a `.go` suffix
  - `go-build-input`: `go.mod`, `go.sum`, or an embed target
  - `markdown`: a `.md` suffix
  - `prose-policy`: `.bench/prose-exclusions`
  - `roadmap-board`: `ROADMAP.md` or the `roadmap/` prefix
  - `decision-documents`: the `decisions/` prefix or a `specs/<slug>/decisions/` prefix
  - `capture-retros`: the `capture/retros/` prefix
  - `benchkit-profile`: `projects/benchkit.md`
- The four document classes share their names with the registry's `InputSource` values. Each selects the dev-tier registry checks whose `Inputs` equals that value. The kit lane declares one check per such registry check, named after it, with the argv `<run binary> test --check <name>`.
- The embed targets come from `packagesurface.EmbedTargets(root)`, an exported derivation over the composed checkout's Go sources. `RequiredBuildPackAssets` composes the same function.
- `gate.SelectLane(checks, changes)` returns the declared checks the classes select, in declared order, with `unknown` selecting every declared check. It returns the class names in table order. `RunLane` applies it when the lane is selective, and `LaneResult` carries the selected check names and the class names. The authority prints them in the lane line. The lane record schema does not change.
- The lane line for a selective lane is `lane{outcome=pass,checks=<names>,classes=<classes>}`. A manifest lane keeps `lane{outcome=pass,checks=<names>}`. A fail keeps `lane{outcome=fail,check=<name>}`.
- The profile lane table gains a third column, `selected by`, that lists the classes per check. The `profile-lane-table` check renders the class names from the kit table and reds a stale cell, by the existing argv rule.

## Bootstrap authority before execution

The chain starts as today. The installed wrapper selects the Bench executable,
the commit verb composes the tree, and the lane builds one private run binary
from the composed tree. The attribution owner refuses a special or unreadable
named path before that build, so no check runs on such a path (PL6, PL40).
The selection reads bytes and executes nothing, and the embed derivation reads
Go sources and executes nothing.

A document check adds two executable hops. The lane's run binary starts
`bench test --check <name>`, which reuses that same binary through
`BENCH_RUN_BINARY`. That verb starts `go` from `PATH`, and `go test` builds
and runs the conformance test binary from the composed tree's sources. Neither
hop is new to the lane's trust posture. The `vet` and `build` checks already
start `go` from `PATH`, and the run binary itself is built from the same
candidate tree. The gate's `test` phase runs the same candidate test binary.

The reviewer accepts at sign-off that the lane executes candidate-controlled
test code in the private checkout, as ADR 0002 accepts for the gate.

## Testing decisions

- A good test drives the real derivation against a real Git repository and reads the change list, the lane line, and the exit code.
- The derivation rows attach to a new `internal/gate/lane_select_test.go` with fixtures by the precedent of `laneFixture` and `commitFiles` in `internal/gate/authorization/lane_test.go`.
- The selection rows are a table test over `gate.SelectLane` with literal change entries.
- The end-to-end kit-lane rows run `RunLane` with `BenchkitLane` and a stub run-binary factory, by the precedent of `internal/worktree/test_run_test.go`. The stub records its argv. The fixture holds no `go.mod`, so a selected `vet` or `build` reds the lane and proves the selection.
- The commit rows attach to `internal/commit/lane_test.go` beside `TestLaneProseGradesOnlyTheNamedMarkdown`, with the manifest fixture `laneRepo` builds.
- The merge row is the existing `TestMergeResolvesTheProsePlaceholderToTheIncomingMarkdown`, which someone reads in the build session.
- The profile rows call `checkProfileLaneTable` on a fixture root and on the live tree.
- The gate observes the feature in the `test` phase for every package above, and the `profile-lane-table` check observes the advertisement.

### Seam diagram

    trigger: bench commit, bench commit --dry-run, bench worktree merge
        │
        ▼
    named paths  ──▶  [ landing.attributedPaths: refuse special or unreadable ]  ──▶  fence
                          ◀ tests attach here: commit_test drives a FIFO and reads stderr
        │
        ▼
    base, tree  ──▶  [ gate.ComposedChanges: git diff --raw --no-renames -z ]  ──▶  []ComposedChange
                          ◀ tests attach here: lane_select_test reads the list from a fixture repo
        │
        ▼
    changes  ──▶  [ gate.SelectLane: class table → declared checks ]  ──▶  checks, classes
                          ◀ tests attach here: a table test over literal entries
        │
        ▼
    checks  ──▶  [ gate.RunLane: private checkout, run binary, schedule ]  ──▶  lane line, lane record
                          ◀ tests attach here: RunLane with a stub factory that records argv

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PL1 | 1, 2 | `ComposedChanges` over a base that holds `docs/a.md` and `docs/b.go`, and a tree that changes both, lists exactly `docs/a.md` and `docs/b.go` with status `M` | `internal/gate/lane_select_test.go` (`TestComposedChangesExpandsANamedDirectory`) | a list built from the named paths holds `docs` and no file |
| PL2 | 3 | a tree that renames `docs/old.md` to `docs/new.md` lists `docs/new.md` with status `A` and `docs/old.md` with status `D`, and no `R` entry | `internal/gate/lane_select_test.go` (`TestComposedChangesRepresentsARenameAsDeletionAndAddition`) | rename detection emits one `R` entry that no class reads |
| PL3 | 16 | a tree that adds a symbolic link `link.go` lists it with destination mode `120000` | `internal/gate/lane_select_test.go` (`TestComposedChangesCarriesTheSymlinkMode`) | a name-only reader classifies the link as Go source |
| PL4 | 4 | the authority over a base with `kept.md`, `changed.md`, `gone.md` and a tree that changes `changed.md`, adds `added.md`, and deletes `gone.md` hands the prose check exactly `added.md changed.md` | `internal/gate/authorization/lane_test.go` (`TestLaneAuthorityDerivesTheProseSubjectFromTheBase`) | a subject that carries `gone.md` grades a file the tree lacks |
| PL5 | 5 | a tree that adds `café notes.md` hands the prose check that path's own bytes as one argument | `internal/gate/authorization/lane_test.go` (`TestLaneAuthorityCarriesANonASCIIProsePathVerbatim`) | a newline-framed read C-quotes the name, and a fields split halves it |
| PL6 | 6 | `bench commit -m m fifo` where `fifo` is a FIFO exits 1, prints `special file "fifo" is not attributable`, and the lane's tally file is absent | `internal/commit/lane_test.go` (`TestLaneRefusesASpecialNamedPathBeforeAnyCheckRuns`) | a lane that runs before attribution blocks on the FIFO or grades it |
| PL7 | 7 | a merge whose incoming side changes `incoming.md` hands the prose check exactly the incoming Markdown | `internal/worktree/merge_test.go` (`TestMergeResolvesTheProsePlaceholderToTheIncomingMarkdown`) | a merge with no base hands the check nothing |
| PL8 | 8 | `LaneAuthority` declares `Base` and declares neither `NamedMarkdown` nor `PreviousTip` | review-owned: the Standards axis reads the type | a kept field keeps a second derivation alive |
| PL9 | 9, 23 | `SelectLane` over the kit lane with one `M` entry for `internal/x/y.go` returns `gofmt`, `vet`, `build` in that order and the class `go-source` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a selection that keeps `prose` pays the prose check on Go |
| PL10 | 10 | one `M` entry for `go.sum` returns `vet`, `build` and the class `go-build-input` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a suffix-only rule classifies `go.sum` as unknown and runs gofmt |
| PL11 | 11 | one `M` entry for `internal/adopt/prepush.sh`, with that path among the embed targets, returns `vet`, `build` and the class `go-build-input` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a rule without embed targets runs every check on the asset |
| PL12 | 12 | one `M` entry for `docs/note.md` returns `prose` alone and the class `markdown` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a selection that keeps `vet` pays the toolchain on prose |
| PL13 | 13 | one `M` entry for `.bench/prose-exclusions` returns `prose` alone and the class `prose-policy` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a suffix-only rule treats the policy as unknown |
| PL14 | 13 | `bench gate-prose <root>` with no path and an exclusion row without a reason exits 1 and prints `malformed exclusion row` | `internal/gate/gate_prose_test.go` (`TestGateProseCommandRefusesAMalformedExclusionRowWithNoSubject`) | a grader that loads the policy only per subject passes an empty list |
| PL15 | 14 | entries for `a.go` and `b.md` return `gofmt`, `prose`, `vet`, `build` in declared order and the classes `go-source,markdown` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a per-class concatenation runs a check twice or out of order |
| PL16 | 15 | one `M` entry for `bin/bench.sh` returns every declared check and the class `unknown` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | an unmatched path that selects nothing passes ungraded |
| PL17 | 16 | one `M` entry for `x.go` with source mode `120000`, and one with destination mode `120000`, each return every declared check and the class `unknown` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a rule that reads one side passes the other side as Go source |
| PL18 | 17 | one `M` entry with destination mode `160000` returns every declared check and the class `unknown` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a submodule pointer path with a `.md` suffix selects prose alone |
| PL19 | 18, 12 | `RunLane` with `BenchkitLane`, `Selective` true, a stub run binary that records argv, and a tree that changes `note.md` alone passes, and the stub recorded one `gate-prose` invocation and no `gate-go` invocation | `internal/gate/lane_test.go` (`TestBenchkitLaneRunsOnlyTheSelectedChecks`) | the fixture holds no `go.mod`, so a selected `vet` or `build` reds the lane |
| PL20 | 18 | that run's stdout holds `lane{outcome=pass,checks=prose,classes=markdown}` | `internal/gate/authorization/lane_test.go` (`TestLaneAuthorityNamesTheSelectedChecksAndClasses`) | the old line names four checks and no class |
| PL21 | 19 | the manifest fixture commit of `note.md` alone prints `lane{outcome=pass,checks=check,gofmt,prose,build}` and no `classes=` | `internal/commit/lane_test.go` (`TestManifestLaneRunsAsDeclared`) | a selection applied to a manifest lane drops the project's own checks |
| PL22 | 20 | `LaneForCommit` at a root equal to `BENCH_KIT` returns `Selective` true, and at a manifest root returns `Selective` false | `internal/gate/lane_test.go` (`TestLaneForCommitMarksOnlyTheKitLaneSelective`) | a flag set for every lane selects inside a linked repo |
| PL23 | 21 | `bench commit --dry-run` with the manifest fixture prints the lane line and moves no ref | `internal/commit/lane_test.go` (`TestLaneDryRunStatesTheOutcomeAndPublishesNothing`) | a dry run that skips the authority prints no lane line |
| PL24 | 22 | after a selective `RunLane` the lane record names the tree, the lane, and the outcome, and the gate cache and the evidence store are absent | `internal/gate/lane_record_test.go` (`TestRunLaneRecordsItsOwnFileOnly`) | a record that gains the selection or a verdict cache entry reads as green |
| PL25 | 2 | the manifest fixture commit that names the directory `docs` holding a 27-word sentence in `docs/note.md` exits 1 and prints `docs/note.md:3:` | `internal/commit/lane_test.go` (`TestLaneProseGradesAMarkdownFileUnderANamedDirectory`) | the named-path filter drops the directory and passes the sentence |
| PL26 | 24 | `EmbedTargets` over a source `internal/adopt/link_hook.go` that declares `//go:embed prepush.sh` returns `internal/adopt/prepush.sh` | `internal/packagesurface/assets_test.go` (`TestEmbedTargetsResolveAgainstTheSourceDirectory`) | a target resolved against the root names a path that does not exist |
| PL27 | 24 | `RequiredBuildPackAssets` over that fixture lists `internal/adopt/prepush.sh` | `internal/packagesurface/assets_test.go` (`TestRequiredBuildPackAssetsCarryEmbedTargets`) | a second derivation in the pack list drifts from the exported one |
| PL28 | 25, 26, 27, 28 | `BenchkitLane` holds one `<run binary> test --check <name>` check per entry of `registry.Checks` at the dev tier whose `Inputs` is a document class, named after it, and no other `test --check` row | `internal/gate/lane_test.go` (`TestBenchkitLaneDocumentRowsFollowTheRegistry`), expectation enumerated from `registry.Checks` | a hand-written row list misses a check the registry adds |
| PL29 | 25, 31 | one `M` entry for `roadmap/FT1.md` returns `prose` and `roadmap-detail-integrity` and the classes `markdown,roadmap-board` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a first-match rule drops one of the two |
| PL30 | 26 | one `M` entry for `specs/x/decisions/map.md` returns `prose` and `decision-map-integrity` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a rule that reads only top-level `decisions/` misses the spec-local map |
| PL31 | 27 | one `M` entry for `capture/retros/x.md` returns `prose` and `retro-improvement-markers` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | an unmatched retro path runs the whole lane |
| PL32 | 28 | one `M` entry for `projects/benchkit.md` returns `prose`, `guidance-prose-budgets`, and `profile-lane-table` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a family bound to one check misses the second |
| PL33 | 29 | every document class name is a valid registry `InputSource`, and `registry.Checks` binds at least one dev-tier check to it | `internal/gate/lane_select_test.go` (`TestDocumentClassesAreRegistryInputSources`) | a class spelled apart from the registry binds to no check |
| PL34 | 30 | every document check in `BenchkitLane` carries the run-binary token as `argv[0]` | `internal/gate/lane_test.go` (`TestBenchkitLaneTable`) | a literal `bench` runs an installed executable and builds a second one |
| PL35 | 32 | `RunLane` with the stub run binary that exits 1 on `test --check roadmap-detail-integrity`, and a tree that changes `ROADMAP.md`, fails with `Check` equal to `roadmap-detail-integrity` | `internal/gate/lane_test.go` (`TestBenchkitLaneDocumentCheckFailsTheLane`) | a document check run as optional reads red as pass |
| PL36 | 33, 34 | `checkProfileLaneTable` over a fixture profile whose `gofmt` row reads `selected by` `markdown` returns one diagnostic that names `gofmt` | `internal/conformance/profile_lane_table_test.go` (`TestProfileLaneTableRedsAStaleSelectedByCell`) | a parser that ignores the third column passes any advertisement |
| PL37 | 33 | `checkProfileLaneTable` over the live tree returns no diagnostic | `internal/conformance/profile_lane_table_test.go` (`TestProfileLaneTableHoldsOnTheLiveTree`) | the profile table lacks the document rows and the column |
| PL38 | 23 | `SelectLane` over any entry set returns a subsequence of the declared lane | `internal/gate/lane_select_test.go` (`TestSelectLaneReturnsADeclaredSubsequence`) | a selection that appends a check adds one nobody declared |
| PL39 | 39 | `BenchkitPhases` keeps the `test` phase argv `go test -trimpath -count=1 ./...` | `internal/gate/lane_test.go` (`TestBenchkitPhasesKeepWholeProjectTestArgv`) | a lane that narrows the phase table makes the landing grade less |
| PL40 | 6 | a named path with mode `000` refuses with `unreadable path`, and an empty change list refuses with `nothing to commit`, each before the authority runs | `internal/landing/state_test.go` (`TestLandPreAuthorizationRefusalTable`) | an authority that runs first grades the unreadable or empty tree |
| PL41 | 40 | each lane check that carries a gate check's name runs that check's bound over the composed checkout, per the counterpart table in Further notes | review-owned: the Standards axis reads each argv against the phase table and the registry | a `vet` row with a narrowed package list keeps the gate's name |
| PL42 | 11, 15 | one `M` entry for `templates/a.txt` with the embed target list `templates/*` returns every declared check and the class `unknown` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a glob match narrows the lane on a pattern nobody resolved |
| PL43 | 35 | the `CONTEXT.md` core terms list holds one `composed change` entry and one `path class` entry, and the `path class` entry states that `unknown` selects every declared check | review-owned: the Spec axis reads the glossary | a glossary without the terms lets a synonym drift in |
| PL44 | 38 | the `CONTEXT.md` `fast lane` entry names no fixed four-check list and states that the composed changes select the checks | review-owned: the Spec axis reads the entry | the old entry advertises four checks the lane no longer always runs |
| PL45 | 36 | ADR 0017 holds no sentence that calls the lane a list never derived from the diff, states that the composed changes select declared checks, and carries no file path and no code | review-owned: the Spec axis reads the ADR | the accepted ADR contradicts the shipped behavior |
| PL47 | 9, 11 | a `D` entry for `internal/x/y.go` returns `gofmt`, `vet`, `build`, and a `D` entry for the embed target `internal/adopt/prepush.sh` returns `vet`, `build` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | a rule that skips deletions never compiles the dangling reference |
| PL46 | 37 | `CHANGELOG.md` holds one `### Changed` entry under `## [Unreleased]` that names the `classes=` cell of the lane line | review-owned: the Spec axis reads the changelog | a release ships a changed lane line with no note |


### Edge inventory

- A path with a space or a non-ASCII byte. The `-z` framing carries the bytes, and PL5 proves the prose subject.
- A path with a control byte. The lane line names checks and classes and no path (PL20).
- An empty change list. The landing refuses `nothing to commit` before the authority runs (PL40).
- A deleted Go file, a deleted embed target, and a renamed Go file. Each side classifies by its own path (PL47), so `build` runs and reds a dangling reference.
- A typechange from a file to a link, or a link on the source side. Either mode is `120000`, so the change is `unknown` (PL17).
- A `//go:embed` pattern that names a directory or a glob. The derivation resolves the literal pattern, so a file under it is `unknown` and selects every check (PL42). Over-selection is the safe direction.
- A Markdown file under a prefix that `.bench/prose-exclusions` names. The class selects prose (PL12), and the grader skips the subject (`TestGradeNamedHonorsExclusions`).
- A `.bench/prose-exclusions` change with no Markdown change. The prose check runs on an empty subject, and a malformed row reds (PL14).
- A `_test.go` file. It is Go source, and `vet` compiles it.
- A `.mjs`, `.sh`, `.json`, or `.yml` path. Each is `unknown` and keeps today's lane.
- A `.go` symbolic link under `.claude/skills/`. The mode is `120000`, and the change is `unknown`.
- A private checkout whose `git.Root()` is the checkout itself. `bench test --check` grades that root and compiles from `BENCH_KIT`, which the lane sets to the checkout.
- Two concurrent lanes on one repository. Each holds the shared cache lock (`TestLaneRunHoldsTheCacheLockAcrossItsChecks`).
- A manifest lane whose check names match the kit's. `Selective` is false, so the names select nothing (PL21).
- **Won't handle** a `paths` key on a manifest lane entry — a linked repo keeps its declared lane, and the kit lane survives.
- **Won't handle** a shellcheck lane check — a shell path is `unknown` and selects every declared check, so the gate's optional phase survives.
- **Won't handle** the `go-source` registry checks in the lane — Go source selects gofmt, vet, and build, and the `test` phase survives.
- **Won't handle** the `.agents/` subjects of `guidance-prose-budgets` — the class is the registry's advertised input, and the gate's check survives.
- **Won't handle** the selection in the lane record — the lane line carries it, and FT232 owns the record's shape.
- **Won't handle** an unchanged named Markdown file — the composed tree does not differ there, and the landing gate survives.

## Ownership fences

- `internal/gate/lane.go`
- `internal/gate/lane_select.go` (new)
- `internal/gate/lane_select_test.go` (new)
- `internal/gate/lane_test.go`
- `internal/gate/lane_record_test.go`
- `internal/gate/gate_prose_test.go`
- `internal/gate/authorization/authorization.go`
- `internal/gate/authorization/lane_test.go`
- `internal/commit/commit.go`
- `internal/commit/lane_test.go`
- `internal/worktree/land.go`
- `internal/worktree/merge.go`
- `internal/worktree/merge_test.go`
- `internal/status/status_test.go`
- `internal/packagesurface/assets.go`
- `internal/packagesurface/assets_test.go`
- `internal/conformance/profile_lane_table_test.go`
- `internal/conformance/tier_test.go` (the live-tree test inventory that PL37 joins)
- `projects/benchkit.md`
- `CONTEXT.md`
- `CHANGELOG.md`
- `docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md`
- `specs/path-aware-lane/`

## Out of scope

- A `paths` key on a manifest lane entry, so a linked repo selects by class — 4 edits, 1 gate run.
- A shellcheck lane check for a shell path — 3 edits, 1 gate run.
- `bench test --changed` package selection for the lane's Go checks — FT277 owns it, 6 edits, 2 gate runs.
- The lane record's `check` and `diagnostic` fields — FT232 owns them, 3 edits, 1 gate run.
- A `bench help` row for `bench commit` that names the lane — 2 edits, 1 gate run.

## Further notes

Gate counterpart table. A lane check carries a gate check's name only when it
runs the same bound over the same bytes.

| lane check | gate counterpart | same name | mutation that reds both |
|---|---|---|---|
| `gofmt` | the `gofmt` phase, `bench gate-go gofmt` over the whole tree | yes | an unformatted `.go` file |
| `vet` | the `vet` phase, `go vet -trimpath ./...` | yes | a `fmt.Printf` verb count mismatch |
| `build` | the `test` phase compiles every package | no | a syntax error in a `.go` file |
| `prose` | the `prose-mechanics` check over the whole tree | no, the subject is the changed Markdown | a 27-word sentence in a changed `.md` |
| `roadmap-detail-integrity` | the same registry check over the whole tree | yes | a `roadmap/FT<n>.md` with no index row |
| `decision-map-integrity` | the same registry check | yes | a map ticket with a malformed state |
| `retro-improvement-markers` | the same registry check | yes | a paragraph under a retro improvement heading |
| `guidance-prose-budgets` | the same registry check | yes | a budget row below its subject's line count |
| `profile-lane-table` | the same registry check | yes | a stale argv cell |

Measured on 2026-08-29 with a warm cache: `go vet` costs 0.8 s to 2.5 s,
`go build` 0.7 s, `gofmt -l` 0.2 s, and one `bench test --check` 2.0 s. The
lane's floor is the private run-binary build, which this spec does not change.

The merge caller changes only by the field rename, because the authority
already derived the prose subject from the previous tip. The reviewer can veto
the merge's share of the selection at sign-off.

Additions beyond `roadmap/FT215.md`, flagged for veto at sign-off:

- the `classes=` cell of the lane line (PL20), so a reader learns why a check ran
- dry-run parity (PL23) and the explicit manifest-lane rule (PL19, PL21, PL22)
- the exported embed derivation (PL26, PL27), so the embed fact has one source
- the profile's `selected by` column (PL36, PL37), so the advertisement stays gate-pinned
- the glossary, the ADR, and the changelog (PL43 to PL46), so the guidance reads current

Group D runs opus / medium. The profile's Lines section routes guidance prose
to mid / high, and the scorecard's current decisions limit every subagent to
Opus low or medium (2026-08-26). The later direction wins here, and the
reviewer owns the reconciliation of the two sources.

Review round 1 (codex `gpt-5.6-sol` / high) folded these findings:

- `LaneRequest` carries `Changes`, so the authority derives once and `RunLane` selects.
- PL17 covers the destination-side link. PL39, PL40, and PL41 cover the oracle
  sentence, the unreadable path, and the counterpart rule from the source.
- PL5 carries a space, and PL42 covers the embed pattern. The bootstrap section
  traces the `go test` hops and records the trust assumption.
- Stories 35 to 38 take their own review-owned rows, and the guidance ticket is
  blocked by the selection ticket alone.
- PL28 and PL33 enumerate from `registry.Checks`, and the new file takes the
  lane types so the lane run's file does not grow.

Review round 2 verified those folds and left one partial: three edge lines
still promised behavior without a row. The author folded it. The control-byte
line keeps only its PL20 claim, the deletion line takes PL47, and the excluded
Markdown line cites the prose package's own test. Two non-blocking remainders
stay as reviewer calls. One is the Group D line conflict above. The other is
the package-green invariant row each ticket ends with, by the tree's precedent.
