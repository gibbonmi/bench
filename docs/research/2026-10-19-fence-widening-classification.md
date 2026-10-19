# Fence widening classification: what the October `Writes:` widenings are

Recommendation: stop adding slicing prose rules for the mechanical classes, because the classes already exist and the slicer still misses them. Find mechanical callers after the edit, with `bench consumers --changed`, and widen them in one step. Record the trigger of each widening, so that a retro can count slicing misses apart from review scope. Treat review-created scope as a normal build cost, not as a slicing defect.

Most of the October widening volume is not a slicing miss. Pre-FT293 closure bookkeeping and one composition fold make up about 75 percent of the 1689 added paths. Of the remaining build-time widenings, about 41 percent came from review findings, about 25 percent were mechanical misses, and about 14 percent were design choices of the author.

Scope: every commit from 2026-10-01 to 2026-10-19 that added entries to the `Writes:` line of an existing ticket file
Evidence status: one month of history, classified by two independent readers; per-path classes are partly inferred
Source repository: /home/mgibs/workspace/bench
Source commit: b29c4df2d
Retrieved: 2026-10-19
Consumed by: the reviewer's decision on FT369, on the slicing checks in `craft-tickets`, and on the trigger record of `bench spec widen`
Drift: refresh when `bench spec widen`, the slicing checks, or the build fence change, or after three post-FT293 builds land
Retire when: an approved decision replaces these proposals and a post-FT293 sample confirms or refutes the closure-volume result

This report changes no command, skill, or check. It proposes changes for the reviewer to decide.

## Question graph

1. In which phase did each widening occur: spec stage, build-entry staleness pass, or build?
2. What triggered each widening?
3. What class is each build-time path?
4. Did a slicing rule already cover each mechanical miss on the date of the miss?
5. What are the totals by spec and by week?
6. What should Bench change?

Questions 1 to 5 hold the facts. Question 6 holds the proposals.

## Method

The coordinator extracted every commit since 2026-10-01 that modified an existing `specs/*/tickets/*.md` file and added entries to its `Writes:` line. The extraction compared the removed and the added `Writes:` values of each ticket in each commit, with the `(new)` and `(deleted)` markers stripped. The result is 232 events (one event is one commit and one ticket) in 180 commits and 31 specs, with 1689 added paths.

Two read-only readers on opus at high effort received the same charge and the same dataset. Neither reader saw the other's return. Each reader placed every event in its spec's history, classified the trigger from the commit text, and opened a sample of the build commits that edited each path. Reader A opened more than 70 paths in 13 specs. Reader B opened 55 paths in 16 specs. Both readers also ran a sweep over all 724 build-time paths to find the paths that no build commit edited. The coordinator reopened the load-bearing commits and closed each disagreement with a targeted read.

The tables below give both readers' counts where they differ. The difference is the classification spread between two readers.

## 1. Phase

### Facts

| Phase | Events | Paths |
|---|---|---|
| Spec stage | 58 to 60 | 922 to 925 |
| Build-entry staleness pass | 13 to 15 | 40 to 43 |
| Build | 159 | 724 |

Reader B placed every spec-stage event before its sign-off commit. For example, `3f9e3b1bc` (17:30) comes before `19dbdca85` (18:35) in closure-removal. The build starts at the first commit that records a ticket author assignment, for example `fad0dbbb6` for closure-removal and `d914ceffc` for refusal-route-registry. Three placements are judgment calls: `9885b6eca`, `ac53d5da1`, and the inline start of roadmap-delivery-commitment at `2c01b5d23`.

## 2. Trigger

### Facts

| Phase and trigger | Events | Paths |
|---|---|---|
| Spec stage: closure proposals from build preflight | 38 | 880 |
| Spec stage: slicing review finding | 19 to 20 | 41 to 42 |
| Build: implementation review finding | 77 to 86 | 158 to 186 |
| Build: orchestrator plan change before dispatch, or composition | 37 to 46 | 69 to 430 |
| Build: gate or preflight red | 17 to 20 | 61 to 419 |
| Build: ticket author report | 16 to 19 | 50 to 75 |

The wide path ranges come from one event pair. Reader A counted the cli-desktop composition fold as an orchestrator plan change, and reader B counted it as a gate trigger.

Examples of each build trigger:

- **Review finding:** `b2ae501e9` ("CR2 finding S1"), `76fdea60b` ("BW-E review found S1-S3").
- **Author report:** `38646f795` ("The ticket 05 author found eleven fixture callers"), `e6ea15f07`.
- **Gate or preflight red:** `b30976514` ("The pre-FT293 fixture-closure and registry-closure rows refuse the ticket 02 charge"), and `818204589`, a decision map that cites a roadmap file that the landing deletes.
- **Orchestrator plan change:** `6d37c0238`, which adds a test-count pin to the fences of five tickets before dispatch.

## 3. Path class of build-time widenings

### Facts

Two events in cli-desktop-consistency, `827af59ba` (220 paths) and `b3c08d28e` (134 paths), name paths inherited through a composition with main. The spec text states that their contents equal the destination commit, and no build commit edited any of them. The coordinator confirmed the path counts of both events.

| Class | All build paths | Without the cli-desktop composition (about 360) | Share without it |
|---|---|---|---|
| Review scope: new behavior or files that a review finding added | 146 to 155 | 146 to 152 | about 41% |
| Mechanical caller: a consumer of a changed or deleted symbol or helper | 47 to 56 | 47 to 56 | about 14% |
| Mechanical pin: a fixture, golden, count pin, help allowlist, or citation that pins changed text | 37 to 42 | 37 to 42 | about 11% |
| Design-conditional: a structure-budget split, a fold into an owner, or an import-cycle choice | 50 to 52 | 48 to 52 | about 14% |
| Forced closure: an entry that only a closure rule required | 402 to 411 | 44 to 53 | about 13% |
| Named but never edited, or other | 20 to 31 | 20 to 31 | about 7% |

The sweep found 416 to 425 of the 724 build paths with no edit by the build. 354 to 357 of them are the cli-desktop composition paths. 44 to 49 are registry, canary, or anchor closure entries, for example `ea408aec1` (0 of 5 edited) and `cfe82163c` (0 of 8 edited). 13 to 19 are paths that the orchestrator named and the author did not use, for example 4 of the 9 paths in `a54fe8fca`.

Of the 564 spec-stage closure-proposal paths that reader B could resolve against a landing, 534 (95 percent) were never edited by that landing. Another 316 spec-stage paths belong to two specs that are not built yet (guidance-weight and markdown-block-reader).

Examples verified against the editing commit:

- **Mechanical caller:** `38646f795` widened ticket 05 by eleven fixture callers. The build commit `828e39ab8` changes one `SpecBody` call or one `Fence:` field in each file.
- **Mechanical pin:** the worktree test-count pin moved in `29eabd309`, `f99f7bcd6`, and `45fa836e6`. The count in `source_tip_test.go` moved in `1e8846bb0`.
- **Design-conditional:** `a688f6844` split the face files for the structure budget. `aa7a35c45` states that `widen_test.go` is at its structure budget.
- **Review scope:** `d723f0c38` led to `2fd8872cc`, which adds `RequireAbsent`, and `aee777a13` led to the CR3 repair `ea2cc4659`.

### Tested result

At the parent of `38646f795`, a reference search for `SpecBody` and the `Fence` field found all eleven widened callers. A plain call-text search found eight of them, and the other three use the struct-literal key. The registry-order tests that `273a28a49` and `5cdf305cc` added depend on data order, not on a symbol, so a symbol search does not find them.

### Inference

Slicing can predict only the mechanical classes, which are about 25 percent of the real build widenings, or about 90 paths in one month. Slicing cannot predict review scope, because review creates it after the code exists. Design-conditional paths depend on choices that a fresh author makes later.

## 4. Rule coverage of mechanical misses

### Facts

`f191167a7` (2026-10-10 08:48) added the fence-expansion evidence classes to `.agents/skills/bench-craft-tickets/references/slicing-checks.md`: callers from `bench consumers`, fixture and canary pins, anchor registries, count pins, structure-budget files, and shared files. `47eed422b` (12:42 the same day) added the keep rule. Before that date, the file had a fixture-helper call-site rule and a rendered-output search rule, but no count-pin rule and no general caller rule.

The two readers checked 31 sampled misses against the version of the file on the date of each miss.

| Verdict | Reader A (14 misses) | Reader B (17 misses) |
|---|---|---|
| A rule covered the miss, and the slicer missed it | 8 | 8 |
| No rule covered the miss | 5 | 9 |
| Partial or not applicable | 2 | 0 |

Before 2026-10-10, most misses were count pins and callers, and no rule covered them. After 2026-10-10, most misses fall in a class that the rule names. The closure-removal spec was sliced on 2026-10-18 under the full rule set, and its build still widened by about 27 mechanical paths.

Four miss shapes recur and have no rule:

- **Deleted roadmap citations:** a decision map cites a `roadmap/FT<n>.md` file that the landing deletes (`afc3ddec2`, `818204589`). The gate catches this one.
- **Help allowlists:** the help-row projection allowlist (`4ad69dd9b`, `1e9a2f97f`).
- **Registry-order tests:** tests that depend on the order of a registry (`273a28a49`, `5cdf305cc`).
- **Stale comments:** comments and wording that name a removed concept (`5fb3be66a`, `87c2b16ab`).

### Inference

More prose does not close the covered classes. The 2026-10-10 report reached the same result for the second-caller rule. FT369 now holds about 30 occurrences and keeps adding rules of this kind.

## 5. Totals by spec and by week

### Facts

| Spec | Paths (spec stage / build) | Main build classes |
|---|---|---|
| cli-desktop-consistency | 89 / 364 | 354 to 358 inherited composition paths |
| guidance-weight | 311 / 0 | not built |
| closure-removal | 157 / 41 | 26 caller, 5 to 12 pin, 5 to 8 review |
| refusal-route-registry | 50 / 90 | 41 to 43 review, 13 to 22 design, 10 to 14 pin, 6 to 11 caller |
| spec-stage-binding | 102 / 0 | spec stage only |
| repair-confirmation | 66 / 3 | spec stage only |
| build-landing-window | 19 / 28 | 19 review |
| roadmap-delivery-commitment | 8 / 25 | 13 review |
| preflight-spec-staleness | 0 / 26, plus 7 staleness | 11 review |

| Week | Build events / paths | Spec-stage events / paths |
|---|---|---|
| 2026-10-01 to 10-04 | 38 / 413 (55 without the composition) | 18 / 152 |
| 2026-10-05 to 10-11 | 63 / 139 | 24 to 26 / 508 to 511 |
| 2026-10-12 to 10-18 | 50 / 143 | 16 / 262 |
| 2026-10-19 | 8 / 29 | 0 |

### Inference

Without the composition outlier, build widening holds at about 140 paths and 50 to 63 events each week, with no downward trend. The readers did not normalize these counts by the number of builds.

## 6. Proposals

### Inference: what the reviewer perceives as slicing issues

The widening stream mixes four sources, and only one of them is a slicing miss.

| Source | Share of all 1689 paths | Status after FT293 |
|---|---|---|
| Closure bookkeeping (spec-stage proposals and forced build entries) | about 55% | FT293 removed the closure rows on 2026-10-19 |
| One composition fold (cli-desktop) | about 21% | not checked |
| Review-created scope | about 9% | continues, and slicing cannot predict it |
| Mechanical misses, design choices, spec-stage review and staleness fixes, and unused entries | about 15% | continues |

Each source arrives as the same plan commit, which says "widen ticket N". That shared form makes all four sources look like one slicing problem.

### Proposals

- **P1. Use the tool at edit time.** Stop adding slicing prose for the classes that `slicing-checks.md` already names. Instead, after an author's edit and before the ticket commit, `bench consumers --changed` lists the consumers of each touched declaration. Each consumer file outside the ticket's `Writes:` line becomes one `bench spec widen` call. A path inside the directory of a signed-off entry needs no stop. The composition of these two verbs is a new step in the build, so it needs a reviewer decision and probably a spec.
- **P2. Record the trigger of each widening.** `bench spec widen` records the trigger as a field: review finding with its ID, author report, gate red with its check, or orchestrator plan change. Today the trigger exists only in free commit text, so a retro cannot count slicing misses apart from review scope. FT231 owns measurement and is the nearest owner.
- **P3. Count review scope as review scope.** A widening that a review finding causes is not a slicing defect. The retro and the drain report it apart from the slicing misses.
- **P4. Close FT369's covered classes.** At the next drain, merge FT369's occurrences that fall in a class that the slicing checks already name into P1, and do not add a new rule for them. Keep the four uncovered shapes from section 4 as candidates for a mechanical check or one rule each.
- **P5. Confirm the FT293 effect.** Repeat this classification on the next three landed builds, and confirm that the closure volume is gone.

## Threats to validity

- **Classification spread:** the per-path classes rest on a sample of about 125 opened paths. The readers inferred the rest from the commit text and the path shape. The two readers' class counts differ by up to 10 paths in each class.
- **Unrecorded triggers:** the trigger comes from free commit text, and up to 46 events have no named trigger.
- **Double counting:** a path that one commit adds to several tickets counts once for each ticket. The 1689 rows hold 1565 distinct commit-and-path pairs.
- **Landing-scope sweep:** an edit that the build reverted, or an edit by a different ticket of the same landing, can hide in the landing diff.
- **One month:** the period holds about 30 specs, and the rules changed during it (on 2026-10-10 and on 2026-10-19).
- **`bench consumers` not run:** the readers tested findability with `git grep` and the editing diffs, not with the verb itself.

## What was read and what was not

The coordinator and the readers read these sources:

- the git history of every spec folder in the dataset;
- the editing commits of the sampled paths;
- `slicing-checks.md` at `a4604f8a1`, `f191167a7`, `47eed422b`, and HEAD;
- `bench consumers --help`;
- the 2026-10-10 slicing comparison in this folder.

Nobody read these sources:

- the ticket bodies beyond their `Writes:` lines;
- the review records of each build;
- the plan JSON of the 37 to 46 events with no named trigger.

Nobody ran `bench consumers`.

## Verification record

- [x] The output opens with the recommendation, the scope, and the evidence status.
- [x] It keeps facts, inferences, tested results, and proposals in separate sections.
- [x] It holds one section for each question.
- [x] It uses tables where the comparison matters.
- [x] It keeps each contradiction and each residual unknown.
- [x] It states what could not be verified.
- [x] It ends with a validation plan.
- [x] Each load-bearing claim names a commit, a path, or a recorded measurement.

## Validation plan

1. **Run the tool on a known miss.** At the parent of `828e39ab8`, apply the ticket 05 edit and run `bench consumers --changed`. Record whether it names all eleven callers that `38646f795` added.
2. **Replay one build.** Replay the author step of closure-removal ticket 04 with P1. Count the widening round trips with P1 and without it.
3. **Measure after FT293.** After three post-FT293 builds land, repeat the extraction and the classification. Record the closure share, the review share, and the mechanical share.
4. **Keep only useful shapes.** For each uncovered shape in section 4, record whether it recurs in those three builds before anyone writes a rule for it.
