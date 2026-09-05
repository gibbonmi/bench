# Structural refactor pass

Status: staged

Decision source: the reviewer's prompt of 2026-09-04, "structural refactor pass: six findings, verified against source on 2026-09-04", plus the reviewer's same-day conversation that reopened the count decision and chose a growth ratchet

Verification log: 1 iteration(s) to accept — the sonnet/xhigh round returned one blocking Coverage finding (the missing-branch scan's active filter), folded as SR60 and decision (p); the preflight closure paths and eleven long sentences were folded before the round

## Problem

The 2026-09-04 architecture survey found six places where one fact has two
owners. The `bench worktree` dispatch repeats its root resolution nine times
and states each leaf in three files. Two verbs build the lane authority by
hand. The landing verb reaches up into the adoption surface for one
predicate.

Worktree operand resolution scans the ledger beside the documented lookup.
The spec-or-close classification has one derivation for the first run and
another for the resume. The refresh rule sits under one of its two consumers,
and the other consumer restates it. Each duplicate drifts on its own
schedule. A slice that touches a worktree leaf edits seven files today, and
the survey priced that as the first cost this pass repays.

The reviewer also wants a soft length limit. `bench structure` reports 104
files over their budget today, and nothing stops an over-budget file from
growing. A whole-tree split program is a multi-session cost the reviewer
declined on 2026-09-04.

## Solution

Part 1 is six behavior-preserving refactors, one ticket each. Every ticket
keeps the pre-existing suite green with its test logic unchanged. A mechanical
rename is permitted. A ticket that needs an assertion or an expected output to
change stops and reports, because that is a behavior delta and a reviewer
decision.

Part 2 adds a growth ratchet. `bench structure --growth <base>` reds an
over-budget file that gained lines since the base, and the fast lane runs it
on every commit and merge. Existing debt stays a soft signal, and the census
of the 104 files becomes the restructure backlog. Part 2 lands after Part 1,
and it splits the three over-budget files it grows, so its own commits pass
the check they add.

The pass is incremental. It adds no store, no strangler, and no new canonical
owner. A file splits only where a finding names a responsibility. The module
deepening rows FT217, FT247, and FT218 stay out of this pass.

## User stories

### Group A — the worktree leaf table

Line: opus / low. The grammar is exact, the seam is known, and fourteen byte
pins cover it.

1. As a maintainer, I want each worktree leaf declared once with its name, grammar, root need, and handler, so that one row adds a leaf.
2. As a maintainer, I want the repository root resolved at one site for every leaf that needs it, so that one producer prints the refusal.
3. As an operator, I want the bare, `help`, `--help`, and `-h` forms to print the same bytes and exit codes, so that my scripts survive.
4. As an operator, I want an unknown leaf or flag refused as today before any assignment exists, so that a typo creates nothing.
5. As an operator, I want each leaf's `--help` to answer its own grammar at exit 0, so that the leaf grammar stays the leaf's.
6. As an operator, I want a root-taking leaf outside a repository to print the same refusal at exit 1, so that one line survives.
7. As an operator, I want `list`, `clean`, and `reclaim` to answer identical bytes outside a repository and from a subdirectory, so that boundary roots hold.
8. As an operator, I want `shell` to keep creating and releasing its assignment with no root argument, so that the subshell path does not change.
9. As an operator, I want `land` to keep receiving the invoked executable and never consulting it, so that the landing's trust posture holds.
10. As a reviewer, I want `bench help` to print byte-identical rows, so that the drifted rows and the missing rows stay a recorded decision.
11. As a reviewer, I want both registry parsers to keep reading the worktree entry, so that the two conformance checks stay green without an edit.
12. As a reviewer, I want a differential run over an enumerated argv family at the base and the tip, so that unpinned bytes show.

### Group B — the lane-to-owner constructor

Line: opus / low. Two callers, one mapping, and the lane tests cover both
verbs.

13. As a maintainer, I want one landing constructor that turns a resolved lane and a base into an owner, so that one source maps it.
14. As a maintainer, I want `bench commit` and `bench worktree merge` to each call that constructor once, so that neither builds the authority by hand.
15. As an operator, I want a root with no lane to keep the gate owner in both verbs, so that the nil-lane branch is unchanged.
16. As a reviewer, I want the reviewed landing to keep refusing a lane pass, so that the accepted-kinds split is untouched.

### Group C — the kit-source predicate moves down

Line: opus / low. The move is mechanical, and the adopt tests pin every
answer.

17. As a maintainer, I want `internal/worktree` to import nothing from `internal/adopt`, so that the landing verb does not depend on the adoption surface.
18. As a maintainer, I want the kit-source predicate and the kit-directory derivation beside the gate's kit-root derivation, so that both consumers point down.
19. As an operator, I want `bench link`, the doctor rows, and the broker-change notice to answer exactly as today, so that the move is invisible.
20. As a reviewer, I want the symlink-resolving compare kept without an Abs step, so that the FT217 equivalence question stays open and recorded.
21. As a reviewer, I want the two kit fallbacks recorded as distinct, so that their unification stays my call in the deepening pass.

### Group D — operand ownership

Line: opus / medium. The change crosses into the intent package, and the
census decides which sites move.

22. As a maintainer, I want the canonical-path match rule owned by the documented lookup, so that operand resolution and the handoff lookup cannot disagree.
23. As an operator, I want an id, a label, a prefix, and a path to keep resolving as today, so that every address keeps working.
24. As an operator, I want a path that names a non-active assignment refused with the state named, so that the reason survives.
25. As an operator, I want a path whose tree is gone refused with the recovery line, so that a dead record has a next step.
26. As an operator, I want an ambiguous prefix refused naming every colliding id, so that the collision stays resolvable.
27. As a maintainer, I want the fifteen ledger reads classified in a census, with only the path matches moved, so that a projection stays put.
28. As a reviewer, I want the exact-string identity matches kept unmoved, so that a symlinked pool home does not change their answer.

### Group E — the spec-or-close classification

Line: opus / medium. The predicate gains a second reader, and the resume is
the riskiest seam in the pass.

29. As a maintainer, I want the tickets-only predicate stated once over a tree reader, so that both readers share the rule.
30. As an operator, I want a first-run landing with a tickets-only folder to keep closing it, so that the light path is unchanged.
31. As an operator, I want a resume of an interrupted close to keep authenticating the folder's absence, so that the resume completes as today.
32. As an operator, I want a resume of a spec-backed landing to keep authenticating the transition, so that the status flip is unchanged.
33. As a reviewer, I want `Owner.Land` and `Owner.LandReviewed` to keep their accepted kinds, so that the deliberate split holds.
34. As a reviewer, I want the probe result recorded, so that the "duplicated ownership" call rests on evidence.

### Group F — refresh sited beside its consumers

Line: opus / low. A package move plus one extracted rule, pinned by the
existing shift and worktree tests.

35. As a maintainer, I want the refresh package at `internal/refresh`, so that shift and worktree are its peers.
36. As a maintainer, I want one entry point that owns the refreshed-start-ref rule, so that shift does not restate it.
37. As an operator, I want `bench shift --refresh` to select the fetched start ref and print the same table, so that the loop is unchanged.
38. As an operator, I want `create --refresh` and `shell --refresh` unchanged, so that the two arg-consuming callers see no difference.
39. As a reviewer, I want the bounds-policy registry row to follow the file, so that the timeout bound keeps one named consumer.

### Group G — measurement

Line: none. The main session records these facts, and no delegate runs.

40. As a reviewer, I want the baseline row and the two hypotheses recorded before the build, so that the retro measures against them.

### Group H — the growth ratchet

Line: opus / medium. The group adds CLI output and a lane check, and the
profile check advertises both.

41. As a maintainer, I want `bench structure --growth <base>` to red only an over-budget file that gained lines since the base, so that existing debt stays soft.
42. As a maintainer, I want an over-budget file that lost lines or held its count to pass, so that a split is never punished.
43. As a maintainer, I want a file at or under its budget to pass whatever it gained, so that ordinary edits stay free.
44. As a maintainer, I want a new file over its budget to red, so that a fresh oversized file is flagged.
45. As a reviewer, I want a structure-accept row to exempt its file from the growth check, so that a granted exception holds.
46. As a maintainer, I want the growth row to name the file, both counts, and the limit, so that the restructure target is clear.
47. As a maintainer, I want the fast lane to run the growth check on each commit and merge, so that growth is flagged early.
48. As a maintainer, I want the profile lane table to carry the new row, so that the advertisement matches the lane.
49. As an operator, I want bare `bench structure` and `--since` unchanged, so that the status count and the shift refactor gate keep their meaning.
50. As a maintainer, I want one observed red against a planted growth before the check lands, so that the check is proven to bite.
51. As a reviewer, I want the over-budget census recorded as the restructure backlog, so that a flagged file has a proposed disposition.

## Implementation decisions

- **Leaf table.** The command registry gains a nested-leaf shape. A leaf row names the leaf, its grammar constants, its root need, and its handler. The root need is required, boundary, or none. The dispatcher resolves the root once per call from that value.
- **Root answers.** A required root that fails prints the not-in-repo line at exit 1. A boundary root passes the empty string, and a leaf with no root need receives none.
- **Leaf family answers.** The bare, help, unknown-leaf, and unknown-flag answers stay generic to the family. The help match still needs exactly one argument. The worktree entry keeps its literal AXI child call and its help rows, because two conformance parsers read those literals. The three hand-typed help rows keep their current bytes.
- **Leaf table placement.** The usage text stays in the usage package, and the leaf table references its grammar constants. The table may sit in a new file under `cmd/bench`, and the dispatch file does not grow.
- **Lane constructor.** The landing package gains one constructor that takes a resolved lane and a base. A nil lane answers the gate owner. A lane answers the lane owner with the lane's Checks, Kit, and Selective values and the caller's base. A pure function owns the mapping, and the constructor composes it.
- **Lane constructor callers.** Both verbs call the constructor once. The two accepted-kind lists are unchanged.
- **Kit-source predicate.** The predicate, its symlink-only path compare, and the kit-directory derivation move to the gate package beside the kit-root derivation. They land in a new gate file. The adopt callers and the landing's default joins name the moved symbols.
- **Kit-source imports.** The landing verb drops its adopt import, and adopt gains a gate import. The gate's dependency closure holds no adopt, so no cycle forms. The two fallbacks stay distinct, and one comment beside them names the difference.
- **Operand ownership.** The intent package exports the canonical-path match over a given slice of assignments. It answers every row whose canonical worktree equals the canonical path, in ledger order, in any state. The documented lookup composes it and keeps its active-only answer. The exported match lands in a new intent file, so the assignment file does not grow.
- **Operand callers.** The selector's path arm calls the exported match over the slice it already holds. It then applies its own ambiguity, state, and missing-tree refusals. The missing-branch scan calls the same match and keeps its active-state filter, and the four exact-string identity scans stay as they are.
- **Classification.** The landing package states the tickets-only predicate once over a tree reader. The reader answers two questions: is the folder a directory, and is the spec file absent. The working-tree reader keeps the current stat calls. The git-object reader keeps the current `cat-file -e` and `show` calls over the source commit.
- **Classification callers.** The first run and the resume each call the one predicate through their reader. The name check applies to both readers.
- **Refresh.** The package moves to `internal/refresh` with its tests. It gains one entry point that takes the root, a requested flag, and the output writer. When requested, the entry point runs the refresh and renders the table. It returns the new start ref only when the refresh succeeded.
- **Refresh callers.** The arg-consuming function and the shift loop call the entry point, and the loop maps an empty answer to `HEAD`. The bounds-policy registry row names the new path.
- **Growth mode.** `bench structure --growth <base>` lists the source files that changed between the base and HEAD, with NUL framing and exact-rename pairing. It reads each tip count from the working tree and each base count from the base commit's blob. A file absent at the base counts zero, and an exact rename reads the old path's blob. A file reds when its tip count exceeds both its limit and its base count.
- **Growth limits.** The limit and the accept list come from the engine the all scan uses. The bare scan and `--since` keep their bytes.
- **Growth output.** A red prints one `FILE GREW` row per file with the tip count, the base count, the limit, and the path. One summary line follows, and the command exits 1. No red prints one ok line at exit 0. An unresolvable base prints the git error on stderr at exit 1, as `--since` does.
- **Lane check.** The built-in lane gains a Bench-owned `structure` check after `build`. Its argv names the growth flag and a base token, the lane request carries the base, and the resolver replaces the token. The go-source class selects the check, and a path outside every class already selects every check.
- **Lane advertisement.** The profile lane table gains the row, and the profile check renders the base token as `<base>`. The landing gate is unchanged.
- **Part 2 splits.** The structure command file is one line over its budget, and the lane file and its test file are over theirs. Part 2 adds lines to all three, so each splits along the responsibility the backlog names. The budget and accept loaders move to `internal/structure/budgets.go`. The lane diagnostics and the tap writer move to `internal/gate/lane_output.go`.
- **Part 2 new files.** The lane execution rows move to `internal/gate/lane_run_test.go`. The new growth tests land in a new file, and every other Part 2 line lands in a file under its budget.
- **Backlog.** `capture/restructure-backlog.md` records the census: one row per over-budget file with its disposition, the symbols that move, and the path pins. A reviewer exception lands in `.bench/structure-accept` when the reviewer grants it, not in this pass.

## Testing decisions

- A good test drives the public verb and compares bytes and exit codes with today's values.
- The leaf table takes the fourteen named byte pins plus two new tests. One covers every grammar-bearing leaf's `--help`, and one covers the nine root-required leaves outside a repository.
- The lane constructor takes one new equality test on the pure mapping, and the existing commit and merge lane tests.
- The kit-source move takes the existing adopt and landing tests, and the canonical-path-owner check.
- The operand match takes the existing intent lookup tests, the operand tests, and one new any-state test in intent.
- The classification takes the existing landing and worktree tests, and one new test for the git-object reader.
- The refresh move takes the three moved tests, the shift refresh test, the worktree refresh tests, and one new test on the entry point.
- The growth mode takes new unit tests in a new file beside the since-mode test. The cases: grew, shrank, held, under budget, added, accepted, budgeted, renamed, non-ASCII name, unreadable accept file, and unresolvable base.
- The lane check takes the extended lane table pins and one new commit lane test beside the prose lane test.
- Copy-survival rows are review-owned and run `bench consumers` over the replaced symbol.
- The gate seam is the whole-project gate at the landing, and each worktree commit runs the fast lane.

### Seam diagram

    bench worktree <leaf> [args]
        │
        ▼
    argv  ──▶  [ registry leaf table: root need, handler ]  ──▶  leaf command in internal/worktree
                   ◀ tests attach here: the fourteen byte pins, the two new table tests

    bench commit | bench worktree merge
        │
        ▼
    gate.LaneForCommit  ──▶  [ landing constructor: lane + base → Owner ]  ──▶  lane or gate authority
                                 ◀ tests attach here: the mapping equality test, the lane tests

    bench worktree land (first run) | --resume
        │
        ▼
    tree reader (stat | git objects)  ──▶  [ landing tickets-only predicate ]  ──▶  close or transition
                                              ◀ tests attach here: TestTicketsOnlyFolder, the git-reader test

    bench commit | bench worktree merge (fast lane)
        │
        ▼
    composed changes + base  ──▶  [ structure check: bench structure --growth <base> ]  ──▶  pass or FILE GREW
                                        ◀ tests attach here: the growth unit tests, the commit lane test

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SR1 | 1, 5 | for every leaf row that names a grammar, `bench worktree <leaf> --help` prints `usage: ` plus that grammar at exit 0 | new table test in `cmd/bench` over the leaf table | a row paired with the wrong handler answers another leaf's grammar |
| SR2 | 2, 6 | `path`, `exec`, `show`, `build`, `create`, `release`, `reauthorize`, `merge`, and `land` run in a non-repository directory print `toon.NotInRepo()` on stderr at exit 1 | new table test in `cmd/bench` over the nine root-required leaves | a leaf whose root need was lost runs its handler with an empty root and answers a different refusal |
| SR3 | 2 | `bench consumers git.Root` lists one call inside the worktree dispatch after the build, where it lists nine today | review-owned differential over the consumers output | a surviving copy keeps a second refusal producer |
| SR4 | 3 | bare `bench worktree` prints `usage.WorktreeUsage()` on stdout at exit 2 and acquires nothing | `TestBareWorktreeRefusesBeforeItAcquiresAssignment` | a table that treats an empty argv as a leaf lookup answers the unknown-argument line |
| SR5 | 3 | `bench worktree --help` exits 0, names every kept grammar and the clean grammar, ends with the exec-gate trailer, and names no `bin/bench.sh` or `recovery` | `TestKeptWorktreeOperationsKeepTheirGrammar`, `TestWorktreeHelpNamesCleanGrammar`, `TestWorktreeHelpNamesTheExecGateForm`, `TestRemovedGrammarsRefuseThroughTheirFamily` | a usage text rendered from the leaf table drops the trailer or reorders a grammar |
| SR6 | 3 | `bench worktree help extra` prints `toon.Usage("bench worktree", "help")` on stderr at exit 2 | new test in `cmd/bench` | a table that matches `help` by name alone answers usage at exit 0 for a two-word call |
| SR7 | 4 | `bench worktree unknown` and `bench worktree --unknown` print the family usage line on stderr at exit 2 and leave the ledger unchanged | `TestUnknownWorktreeSubcommandRefusesBeforeItAcquiresAssignment`, `TestUnknownWorktreeFlagRefusesBeforeItAcquiresAssignment` | a fallthrough to the subshell creates a worktree named for the typo |
| SR8 | 5 | `bench worktree create --unknown` refuses through the create grammar, and `bench worktree clean --help` prints exactly the clean grammar | `TestWorktreeCreateKeepsParserFirstDispatch`, `TestWorktreeCleanHelpUsesItsGrammar` | a dispatcher that parses flags before the leaf steals the leaf's grammar |
| SR9 | 7 | `clean --help` and `list` answer identical bytes from the root and from a subdirectory, and `clean --help` answers in a non-repository directory | `TestWorktreeRoutesKeepTheirBytesFromASubdirectory`, `TestWorktreeCleanHelpUsesItsGrammar` | a boundary leaf given a required root refuses outside a repository |
| SR10 | 8 | `bench worktree shell` with `SHELL=true` creates and releases its assignment at exit 0 | `TestWorktreeShellRunsAndReleasesItsAssignment` | a shell leaf given a root need refuses or double-resolves |
| SR11 | 9 | `bench worktree land` refuses with repository proofs and never names or runs the invoked executable | `TestWorktreeLandNeverConsultsTheInvokedExecutable` | a leaf signature that drops the executable field breaks the land route |
| SR12 | 10 | `bench help` prints the byte-identical public inventory, and the wrapper and the binary agree | `TestHelpInventoryIsComplete`, `TestRootAndHelpAlignWrapperAndBinary` | a help row derived from a grammar constant changes the `reauthorize` or `merge` bytes |
| SR13 | 11 | the AXI query registry check and the subcommand routing check stay green over the edited dispatch file | `bench test --check axi-query-registry`, `bench test --check subcommand-routing` | a worktree entry that moves its AXI call or its Name into the leaf table reds a parser |
| SR14 | 12 | every argv row in the differential family below answers identical stdout, stderr, and exit code from the base binary and the tip binary | review-owned differential run with `cmp` over the family | a byte the pins do not enumerate changes in silence |
| SR15 | 13 | the pure mapping from a lane and a base yields a lane authority with the lane's Checks, Kit, and Selective values and that base | new equality test in `internal/landing` | a mapping that drops Selective or the base grades the wrong tree |
| SR16 | 14 | `bench commit` under a declared lane runs the lane and publishes, and `bench worktree merge` runs the declared lane on the composed tree | `TestLanePassStatesItsOutcomeAndPublishesTheComposedSnapshot`, `TestMergeRunsTheDeclaredLaneOnTheComposedTree` | a caller that keeps its hand-built authority skips the constructor |
| SR17 | 14 | `bench consumers landing.NewLane` lists one caller, inside `internal/landing` | review-owned differential over the consumers output | a surviving hand-built authority keeps a second mapping |
| SR18 | 15 | a root with no lane keeps the gate commit, and a merge with no lane publishes the merge tree | `TestNoDeclaredLaneKeepsTheGateCommit`, `TestMergePublishesTheMergeTreeAndMovesTheCheckout` | a constructor that never answers the gate owner grades every commit by a nil lane |
| SR19 | 16, 33 | the reviewed landing refuses a lane pass and names its kind | `TestReviewedLandingRefusesALanePass`, `TestReviewedLandingRefusalNamesTheLanePassKind` | a constructor that widens the reviewed kinds lets a lane pass reach main |
| SR20 | 17 | `go list -deps ./internal/worktree` names no `internal/adopt`, and `bench consumers gate.KitSourceCheckout` lists the landing's default joins | review-owned differential over both outputs | a moved symbol with the old import left behind keeps the edge |
| SR21 | 18, 19 | `bench link` refuses in the kit checkout, and the doctor rows for the kit checkout read green and route the absent pre-push to the fix | `TestLinkInKitSourceCheckoutRefuses`, `TestDoctorKitSourceCheckoutRowsAreGreen`, `TestDoctorKitSourceCheckoutAbsentPrePushRoutesToFix` | a predicate that compares the root with itself answers true for every consumer root, and the link and doctor tests red on it |
| SR22 | 19 | the landing ignores a forged primary executable and seal, with the kit-source seam swapped | `TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal` | a default joins that lost the predicate prints no broker-change notice |
| SR23 | 20 | the moved compare calls no `filepath.Abs`, and the canonical-path-owner check stays green | `bench test --check canonical-path-owner` | an Abs beside EvalSymlinks re-derives the canonical path outside its owner |
| SR24 | 22 | the documented lookup matches through a symlink, ignores a retired assignment, and answers no owner outside the pool | `TestAssignmentForWorktreeMatchesThroughASymlink`, `TestAssignmentForWorktreeIgnoresARetiredAssignment`, `TestAssignmentForWorktreeAnswersNoOwnerOutsideThePool` | a lookup rebuilt over the exported match loses the state filter |
| SR25 | 22 | the exported match answers every row whose canonical worktree equals the canonical path, in ledger order, in any state, including a retired row and a symlinked spelling | new unit test in `internal/intent` with two rows at one path | a match that filters state loses the row the refusal-with-reason names |
| SR26 | 23, 26 | every path-taking verb resolves the label, the id, and an 8-12 character prefix, and an ambiguous prefix names every colliding id | `TestVerbsResolveIdentifierOperands`, `TestPrefixOperandRefusals`, `TestTargetVerbsNameTheResolverReason`, `TestMergeTargetOperandTakesEveryAddress` | a selector that routes the id arm through the path match resolves nothing |
| SR27 | 24 | a `--from` naming a sibling in a non-active state refuses with the state component named | `TestCreateFromRefusesAFailedSiblingIdentityComponent` | a path arm that drops non-active rows answers unassigned instead of the state |
| SR28 | 25 | `list` names one clean landed row for an assignment whose tree is missing, and a target verb names the missing-tree reason | `TestListCommandNamesOneCleanLandedRowForAMissingTree`, `TestTargetVerbsNameTheResolverReason` | a lookup that reads the ledger at the missing path answers no owner |
| SR29 | 27, 28 | `bench consumers intent.Assignments` after the build lists the fifteen production sites minus none, and the census table's two moved sites call the exported match | review-owned differential over the consumers output and a read of the two sites | a build that moves an exact-string identity scan changes its answer under a symlinked pool home |
| SR30 | 29 | the predicate answers true for a folder with no spec file and false for a malformed name, an absent folder, and a present spec file, through the working-tree reader | `TestTicketsOnlyFolder` | a reader that skips the name check answers a nested slug |
| SR31 | 29 | the predicate answers the same four cases through the git-object reader over a fixture commit | new test in `internal/landing` | a git reader that answers true on an unreadable spec file misclassifies a permission error |
| SR32 | 30 | a `--spec` naming a tickets-only folder closes it, an already-removed folder still lands, and an absent folder keeps the unreadable refusal | `TestLandCommandTicketsOnlySpecClosesTheFolder`, `TestLandCommandTicketsOnlySpecLandsWhenTheDestinationAlreadyRemovedTheFolder`, `TestLandCommandAbsentSpecFolderKeepsTheUnreadableRefusal` | a first run that reads the git reader over HEAD misses the working tree |
| SR33 | 31 | a resume of a close interrupted at the marker step authenticates the folder's absence and releases | `TestResumeLandCommandTicketsOnlySpecCompletesAnInterruptedClose` | a resume that never reaches the source commit's objects has nothing to authenticate |
| SR61 | 31 | a resume of a close interrupted at the release step, after the reconcile removed the folder from the destination checkout, authenticates the close and releases | new test beside `TestResumeLandCommandTicketsOnlySpecCompletesAnInterruptedClose` | a resume that reads the working tree finds no folder and refuses the close as a missing transition |
| SR34 | 32 | a resume completes a spec-backed landing and accepts the slug and the path spelling | `TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding`, `TestResumeLandCommandAcceptsSpecSlugAndPath`, `TestResumeLandCommandSpecLessCompletesAnInterruptedLanding` | a predicate that answers close for a folder with a spec file skips the transition proof |
| SR35 | 29 | `bench consumers landing.TicketsOnlyFolder` lists the first-run site and the folders sweep, and the resume file holds no `cat-file -e` and `show` pair over the spec folder | review-owned differential over the consumers output and a read of the resume file | a surviving second derivation keeps two owners |
| SR36 | 35, 39 | the bounds-policy check names `internal/refresh/refresh.go` as the consumer of the refresh timeout and stays green | `bench test --check bounds-policy` | a row that still names the old path reds "does not consume" |
| SR37 | 36, 37 | `bench shift --refresh` selects the fetched remote head as the start ref and prints the refresh table | `TestLoopRefreshResolvesStartAfterFetch` | a loop that maps a failed refresh to the remote head starts from the wrong commit |
| SR38 | 36 | the entry point with the flag off writes nothing and answers empty, and with the flag on under `BENCH_OFFLINE=1` writes the offline row and answers empty | new unit test in `internal/refresh` | an entry point that answers a ref without a successful refresh changes shift's base |
| SR39 | 38 | `create --help` performs no refresh, and `create --from` with `--refresh` refuses before the refresh | `TestCreateCommandHelpPerformsNoRefresh`, `TestCreateFromWithRefreshRefusesBeforeTheRefresh` | an arg consumer that runs the refresh before the grammar fetches on a refusal |
| SR40 | 35 | the three moved refresh tests pass at the new path | `TestRefreshUsesBoundedNoninteractiveFetch`, `TestRefreshFailureAndTimeoutAreNonfatalAndDetailed`, `TestRefreshOfflineStartsNoGitAndNamesFlag` | a move that drops a test file loses the timeout pin |
| SR41 | 36 | `bench consumers refresh.RefreshedStartRef` lists one caller, the entry point | review-owned differential over the consumers output | a shift that keeps its own call restates the rule |
| SR42 | 41, 46 | `bench structure --growth <base>` over a repository where one over-budget file gained lines prints one `FILE GREW` row with the tip count, the base count, the limit, and the path, then a summary line, at exit 1 | new test in `internal/structure` beside `TestCommandSince` | a mode that reuses the touched rule reds an over-budget file that did not grow |
| SR43 | 42 | an over-budget file that lost lines or kept its count since the base prints no row, and the command exits 0 | new test in `internal/structure` | the since rule reds every touched over-budget file |
| SR44 | 43 | a file at or under its limit that gained lines prints no row | new test in `internal/structure` | a growth rule without the limit reds ordinary edits |
| SR45 | 44 | a file absent at the base and over its limit at the tip prints a `FILE GREW` row with a base count of 0 | new test in `internal/structure` | a rule that skips added files lets a new oversized file in |
| SR46 | 45 | an over-budget file with a structure-accept row that gained lines prints no row, and the command exits 0 | new test in `internal/structure` | a growth rule that ignores the accept list voids every grant |
| SR47 | 41 | a file with a `structure.budgets` row grows within that budget and prints no row, and grows past it and prints a row with that budget as the limit | new test in `internal/structure` | a mode with a hard-coded 400 flags a granted file under its grant |
| SR48 | 41 | an exact rename of an over-budget file with no content change prints no row | new test in `internal/structure` | a rename read as an addition reds a pure move |
| SR49 | 41 | a changed source file whose name carries a byte above ASCII is counted, and its row prints the name's own bytes | new test in `internal/structure` | a newline-framed diff C-quotes the name, and the extension filter drops the file |
| SR50 | 41 | `--growth` with a present-but-unreadable accept file prints the loud accept line and exits 1 | new test in `internal/structure` beside `TestAcceptUnreadableIsLoud` | a growth loop that skips the accept read observes an empty list |
| SR51 | 41 | `--growth` with a base that does not resolve prints `git diff failed:` on stderr and exits 1 | new test in `internal/structure` | a swallowed git error passes the lane on a broken base |
| SR52 | 47 | a `bench commit` that raises an over-budget Go file's count fails the lane with `lane{outcome=fail,check=structure}` and the `FILE GREW` line, and a commit that lowers it passes | new lane test in `internal/commit` beside `TestLaneProseGradesOnlyTheNamedMarkdown` | a lane without the check publishes the growth |
| SR53 | 47 | the built-in lane carries a Bench-owned `structure` check after `build` whose argv names the growth flag and the base token, and the go-source class selects it | `TestBenchkitLaneTable` and `TestLaneClassesNameOnlyDeclaredChecks`, extended for the new row | a check outside the class table never runs selectively |
| SR54 | 47 | the lane resolver replaces the base token with the request's base | `TestResolveLane`, extended | an unreplaced token passes the literal placeholder to git |
| SR55 | 48 | the profile lane table carries the `structure` row with `bench structure --growth <base>` selected by go-source, and the profile-lane-table check stays green | `bench test --check profile-lane-table` | a lane row the profile does not advertise reds the check |
| SR56 | 49 | bare `bench structure` and `bench structure --since <base>` print the same bytes as at the base commit over the fixture repository | `TestCommand`, `TestCommandSince`, `TestFileTooLongAndBudget`, `TestAcceptExcludesAndPrints`, and the differential family | a mode change that leaks into the all scan changes the status count |
| SR57 | 50 | `bench structure --growth <base>` over a planted growth in a throwaway repository prints the `FILE GREW` row, and the same command after the revert prints the ok line | review-owned, the build records both outputs in the ticket evidence | a check never seen red is a hope |
| SR58 | 51 | `capture/restructure-backlog.md` holds one row for every file `bench structure` reports over budget at 8eea2d15, each with a disposition | review-owned, the review compares the rows to the structure output | a backlog that misses a file leaves a flagged file without a route |
| SR59 | 41, 47 | `bench structure --growth <part-1 tip>` over the Part 2 diff prints the ok line at exit 0 | review-owned, the build records the output before the Part 2 fold | a pass that grows an over-budget file fails its own ratchet at the next commit |
| SR60 | 22, 28 | the missing-branch scan answers no assignment for a lone retired row whose branch is gone, and answers the active row when a retired row shares its canonical path | new unit test in `internal/worktree` beside the landed classifier tests | a scan rebuilt over the any-state match flags a retired assignment's missing branch |

Not covered: story 21 — the two kit fallbacks are a decision line under Further notes, not a behavior.
Not covered: story 34 — the probe result is a record under Further notes, not a behavior.
Not covered: story 40 — the baseline and the hypotheses are a record under Further notes, not a behavior.

### Edge inventory

- Error paths: a required root that fails answers the not-in-repo line (SR2). An unreadable ledger keeps its typed error through the exported match, because the match takes the slice and reads nothing. A git reader whose `show` fails for any reason answers spec-absent, as today. A refresh that fails or times out answers no start ref (SR38). An unresolvable growth base is loud (SR51), and an unreadable accept file is loud (SR50).
- Empty input: bare `bench worktree` (SR4). An empty argv after a leaf name stays the leaf's grammar refusal. A growth run with no changed source file prints the ok line at exit 0 (SR43).
- Boundary values: `help` with exactly one argument answers usage, and with two it answers unknown (SR6). An 8-character prefix resolves and a 7-character prefix does not (SR26). A file at exactly its limit passes, and one line more reds (SR44, SR47).
- Interrupted state: an interrupted close resumes through the git reader (SR33).
- Re-run idempotency: the differential family runs twice per binary and answers the same bytes (SR14, SR56).
- Hostile paths: a worktree reached by a symlink matches through the canonical compare (SR24, SR25). The kit checkout reached by a symlink spelling matches through the symlink-only compare (SR21). A slug with a separator or `..` is not tickets-only through either reader (SR30, SR31). A source name with a byte above ASCII survives the growth query (SR49).
- Absent versus empty: an absent accept file is an empty grant list, and a present-but-unreadable one is loud (SR50). A file absent at the base counts zero (SR45).
- Partial implementation: a leaf table that keeps one hand-resolved root reds SR3. A caller that keeps a hand-built authority reds SR17. A resume that keeps its own `cat-file` pair reds SR35. A lane without the structure row reds SR53.
- Package variables swapped in tests: `land_freshness_test.go` swaps the kit-source seam field on the joins value (SR22). No test swaps a refresh variable, and the refresh timeout stays the package variable the moved tests keep. The growth tests set `BENCH_MAX_LINES` through the environment as the since-mode tests do.
- Audience: every behavior serves this repository and every linked repository, because the one binary ships the dispatch, the landing, the refresh, and the structure command. The lane check runs only where the kit lane runs, which is this repository. No new directory is read.

**Won't handle** — a `bench help` row for the five unlisted leaves — the default keeps the bytes, and `bench worktree --help` lists every grammar.

**Won't handle** — the `reauthorize`, `merge`, and `--help` row drift — the default keeps the bytes, and the leaf table shows the constants beside them.

**Won't handle** — a usage text rendered from the leaf table — the bare-call pin compares against `usage.WorktreeUsage()`, and the usage package cannot import the registry.

**Won't handle** — one kit fallback for the gate and the adopt surface — the two differ today, and the deepening pass owns the behavior change.

**Won't handle** — an Abs step in the kit-source compare — FT217 records the question, and the canonical-path-owner check would red the pair.

**Won't handle** — the canonical match at the four exact-string identity scans — a symlinked `BENCH_HOME` changes their answer, and the sites keep their exact compare.

**Won't handle** — a root-and-path signature on the documented lookup — the exported match takes a slice, so the three intent tests stay as written.

**Won't handle** — a persisted first-run classification the resume reads back — the resume runs in a fresh process, and the source commit is its durable form.

**Won't handle** — the ignored-file split the probe found — two refusals stop the landing before publication, and the deepening pass may weigh an `--ignored` read.

**Won't handle** — the worktree-to-preflight edge — the fence authorization is real behavior, and the deepening pass owns it.

**Won't handle** — the worktree hub's 91 files against a grant of 18 — the count is diagnostic, and the three policy children stay.

**Won't handle** — a split of the 104 over-budget files — the reviewer chose the ratchet on 2026-09-04, and the backlog records each file's disposition.

**Won't handle** — the growth check in the landing gate — every landing composes lane-passed commits, and the deepening pass may extend the phase table.

**Won't handle** — a growth rule for crowded directories — the reviewer's rule names lines added to a file, and `bench status` keeps the directory signal.

**Won't handle** — NUL framing in the existing `--since` query — the shift refactor gate reads it today, and the change is parked as an idea.

## Ownership fences

- `specs/structural-refactor-pass/`
- `reviews/structural-refactor-pass.md`
- `capture/restructure-backlog.md`
- `cmd/bench/`
- `internal/usage/worktree.go`
- `internal/landing/`
- `internal/commit/commit.go`
- `internal/commit/lane_test.go`
- `internal/worktree/merge.go`
- `internal/worktree/land.go`
- `internal/worktree/land_identity.go`
- `internal/worktree/land_resume.go`
- `internal/worktree/land_tickets_only_test.go`
- `internal/worktree/landed.go`
- `internal/worktree/landed_test.go`
- `internal/worktree/path.go`
- `internal/worktree/worktree.go`
- `internal/worktree/subshell.go`
- `internal/worktree/refresh/`
- `internal/refresh/`
- `internal/gate/kit_source.go`
- `internal/gate/kit_source_test.go`
- `internal/gate/lane.go`
- `internal/gate/lane_select.go`
- `internal/gate/lane_output.go`
- `internal/gate/lane_test.go`
- `internal/gate/lane_run_test.go`
- `internal/gate/lane_select_test.go`
- `internal/gate/authorization/authorization.go`
- `internal/adopt/adopt.go`
- `internal/adopt/doctor.go`
- `internal/adopt/doctor_rows.go`
- `internal/adopt/link.go`
- `internal/adopt/init.go`
- `internal/adopt/setup.go`
- `internal/adopt/upgrade.go`
- `internal/intent/assignment.go`
- `internal/intent/worktree_owner.go`
- `internal/intent/worktree_owner_test.go`
- `internal/intent/assignment_lookup_test.go`
- `internal/shift/loop.go`
- `internal/structure/`
- `internal/conformance/bounds_policy_test.go`
- `internal/conformance/profile_lane_table_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_test.go`
- `projects/benchkit.md`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route`

The fence is the union of the write lists of the eight tickets. The two
registry parser tests, the three command registry files, and the eight canary
fixtures are closure entries. The preflight names them for a write under
`cmd/bench/`, `internal/worktree/`, or the profile. Only the leaf-table
ticket edits the command registry files.

The bounds-policy row is a path rename
that follows the moved file. It is the one test edit Part 1 makes, and the
exit proof permits a mechanical rename. Part 2 extends three pins on purpose:
the lane table pin, the class pin, and the structure usage line.

## Out of scope

- A `bench help` row for the five unlisted worktree leaves: 2 edits, 1 gate run.
- Closing the two help-row grammar drifts and the `--help` row's missing `land`: 3 edits, 1 gate run.
- One kit-directory derivation for the gate and the adopt surface: 4 edits, 1 gate run, behavioral.
- An `--ignored` status read in the dirty-source refusal: 2 edits, 1 gate run, behavioral.
- Canonical compare at the four exact-string identity scans: 4 edits, 1 gate run, behavioral under a symlinked pool home.
- A split of the 104 over-budget files with reviewer exceptions: at least 104 edits and about 20 gate runs, routed to `/bench-shape-idea`.
- The growth check as a landing gate phase: 2 edits (the gate script and the profile phase table), 1 gate run, a `craft-gate` change.
- NUL framing in `bench structure --since`: 1 edit, 1 gate run, behavioral for a non-ASCII name.

## Further notes

Decisions recorded for the reviewer's veto:

- (a) The `reauthorize` and `merge` help rows keep their shorter hand-typed grammar, and the `--help` row keeps its description without `land`. Default: keep the current bytes.
- (b) `create`, `release`, `clean`, `reclaim`, and `land` keep no `bench help` row. Default: keep the current bytes.
- (c) The finding-3 premise is corrected. The gate's kit root falls back to the graded root. The adopt kit directory falls back to the executable's parent, then the current directory. The two are not one derivation. The predicate and the kit-directory derivation move to the gate package beside the kit root, and the two fallbacks stay distinct. The alternative, a new leaf package for both, is rejected as a new seam.
- (d) The kit-source compare keeps its symlink-only resolution with no Abs step. FT217's equivalence question stays open.
- (e) The intent package exports the canonical-path match over a slice, and the documented lookup composes it. The alternative is a root-and-path lookup signature. It is rejected, because a missing tree makes the ledger read at the path fail. That failure would turn the missing-tree refusal into an unassigned answer.
- (f) The four exact-string identity scans stay unmoved (see the census).
- (g) The finding-5 probe. An ignored `specs/t/spec.md` inside a tickets-only folder is invisible to the dirty-source check, which reads no `--ignored`. The first run then classifies the landing as a spec transition, and a minimal spec body refuses at the fence proof. A fence-valid body refuses at the source proof, because the staged bytes are not the source tip's committed spec. No publication occurs, so the resume disagreement is unreachable. The two derivations are duplicated ownership, not a defect.
- (h) The refresh entry point returns an empty ref when no refresh succeeded, and the shift loop maps empty to `HEAD`. The worktree callers already treat empty as the default.
- (i) The bounds-policy registry row is edited to the new path as a mechanical rename.
- (j) On 2026-09-04 the reviewer reopened the prompt's count decision. At 8eea2d15, 105 files exceed 400 lines, 73 exceed 450, and 60 exceed 500. The reviewer weighed a split of all of them and chose a soft limit with a growth ratchet instead. Part 2 is that ratchet, and the census becomes the restructure backlog.
- (k) The growth check runs in the fast lane only, on `bench commit` and `bench worktree merge`. The landing gate is unchanged.
- (l) The growth query pairs exact renames, so a pure move of an over-budget file passes. A rename with an edit reads as a deletion and an addition, and it reds when the new path is over its limit.
- (m) The growth limit is the engine's limit for the path: the `structure.budgets` row when one exists, else `BENCH_MAX_LINES` or 400. A structure-accept row exempts the path.
- (n) Part 2 changes three pins on purpose: `TestBenchkitLaneTable` and `TestLaneClassesNameOnlyDeclaredChecks` gain the structure row, and the structure usage line gains the growth flag.
- (p) The review round's one blocking finding. The missing-branch scan filters on the active state today, and the any-state match would drop that filter in silence. The scan keeps the filter, and SR60 pins it with a new test. The author folded the finding at acceptance without a second round.
- (q) Two leaves answer `--help` with no grammar today. `path` reads it as a target operand and refuses at exit 1. `reclaim` refuses it as an unknown argument at exit 2. Their leaf rows name no grammar constant, so SR1 covers the ten grammar-bearing leaves. The differential family records both answers as they are. A grammar-first `--help` arm in either leaf is a behavior change for a later decision.
- (r) Every adopt fixture sets `BENCH_KIT`, so the suite never reaches the kit-directory fallback branch, before or after the move. The build's swap of that fallback for the graded root stayed green. The self-compare swap redded the link and doctor tests, and SR21 names that kill. The fallback stays untested, as today.
- (o) Part 2 is the ratchet's first dogfood. It splits the three over-budget files it grows, along the responsibilities the backlog names, and it adds no line to any other over-budget file. Part 1 puts its two additions in new files for the same reason: the kit-source predicate and the exported operand match. The per-file grant rows those three files hold today stay until the reviewer removes them.

Census of the fifteen `intent.Assignments` sites in `internal/worktree` (source: `bench consumers intent.Assignments` at 8eea2d15):

| site | enclosing | class | disposition |
|---|---|---|---|
| path.go:60 | resolveAssignment | operand path match through the selector | move the selector's path arm |
| merge.go:67 | mergeAttributed | operand path match through the selector (target and `--from`) | moves with the selector |
| worktree.go:648 | createSiblingStart | operand path match through the selector (`--from`) | moves with the selector |
| landed.go:29 | activeAssignmentWithMissingBranch | active path match by symlink-resolved compare | move to the exported match |
| worktree.go:580 | unmatchedRequestRecovery | exact-string path plus active state, counted | keep (identity scan) |
| resume.go:294 | finishInterruptedExplicit | exact-string path, any state, ambiguity | keep (identity scan) |
| subshell.go:209 | planExplicitWith | exact-string path plus owner id | keep (identity scan) |
| classifier.go:365 | foreignRecoveryAssignment | exact-string path plus label, request, recovery | keep (identity scan) |
| resume.go:60 | assignmentByID | id match | keep |
| clean_landed.go:282 | replanLandedCleanupRow | id match | keep |
| resume.go:463 | sweepOrphanAssignments | list projection | keep |
| resume.go:505 | resumeCleanCommandWith | list projection (count) | keep |
| clean_landed.go:37 | planLandedSet | list projection | keep |
| list.go:46 | ListCommand | list projection | keep |
| clean_unclaimed.go:23 | planUnclaimedAssignmentSet | list projection | keep |

Differential argv family for SR14 and SR56. Run each row from a fixture
repository root unless the row says otherwise, and record stdout, stderr, and
the exit code:

- `worktree`, `worktree --help`, `worktree help`, `worktree -h`, `worktree help extra`
- `worktree unknown`, `worktree --unknown`
- `worktree <leaf> --help` for `list`, `path`, `exec`, `show`, `build`, `create`, `release`, `clean`, `reclaim`, `reauthorize`, `merge`, `land`
- `worktree list` and `worktree clean --help` from a subdirectory
- `worktree clean --help`, `worktree reclaim --help`, `worktree list`, and `worktree path x` from a non-repository directory
- `worktree shell` with `SHELL=true`
- `structure` and `structure --since <base>` over a fixture with one over-budget file

Flagged additions beyond the decision source:

- The exported canonical-path match over a slice in the intent package (decision e).
- The refresh entry point beside the arg consumer (decision h). The source said "have shift call the one entry point", and the existing consumer takes an argv the shift loop does not have.
- The kit-directory derivation moves with the predicate (decision c). The source named the predicate alone.
- The bounds-policy row edit (decision i).
- The lane request's base field and the base token (decision k). The reviewer named the rule, and the lane has no base today.
- Exact-rename pairing and NUL framing in the growth query (decision l). The reviewer's rule is silent on renames and on names.
- The restructure backlog file (story 51). The reviewer asked that flagged files be restructured, and the census gives each a route.
- Three file splits inside Part 2 (decision o). The ratchet would red its own commits without them.
- New tests: two table tests in `cmd/bench` (SR1, SR2) and one intent test (SR25). Also one mapping test and one git-reader test in landing (SR15, SR31). Also one refresh test (SR38), ten growth tests (SR42 to SR51), and one commit lane test (SR52). Each pins a behavior no existing test enumerates.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| each leaf declares its name, grammar constant, root need, and handler once | SR1, SR2 |
| root resolution collapses to one site | SR2, SR3 |
| byte pins that must pass unchanged | SR4 to SR12 |
| no change under `internal/worktree` (finding 1) | SR13, SR14 |
| landing gains the constructor both callers already want; both call it once | SR15, SR16, SR17, SR18 |
| move the kit-source predicate down beside the kit-root derivation | SR20, SR21, SR22 |
| keep the no-Abs behavior | SR23 |
| worktree keeps operand spelling and the refusal-with-reason; ownership comes from the documented lookup | SR24, SR25, SR26, SR27, SR28, SR60 |
| census the fifteen sites; move only the path matches | SR29 |
| the first run produces the classification and the resume re-reads it | SR30, SR31, SR32, SR33, SR34, SR35 |
| do not change the accepted authorization kinds | SR19 |
| move the refresh package up one level; shift calls the one entry point | SR36, SR37, SR38, SR39, SR40, SR41 |
| a soft limit of 400 in the project; flagged for restructure if more lines are added | SR42, SR43, SR44, SR45, SR47, SR48, SR49, SR50, SR51, SR52, SR53, SR54, SR55, SR56, SR57 |
| allow for exceptions | SR46 |
| the census tables become the restructure backlog | SR58 |
| the pass passes its own ratchet | SR59 |

Pre-review proof checklist:

- `Cited symbols`: every symbol in the coverage map resolves at 8eea2d15, and the test names were read in this session.
- `Import edges`: `internal/landing` already depends on `internal/gate` through the authorization package. `internal/gate` depends on no `internal/adopt`, and `internal/adopt` depends on no `internal/gate`. `internal/shift` depended on `internal/worktree/refresh` and `internal/worktree` at 8eea2d15, and on `internal/refresh` after ticket 6. The lane runs `bench structure` as a child, so `internal/gate` gains no import.
- `Source-row clauses and occurrences`: the table above quotes each clause once, and no clause recurs in the source.
- `Promised field labels`: the `FILE GREW` row label and the `lane{outcome=fail,check=structure}` line.
- `Changed-function callers`: the table below.
- `Copy survival`: SR3, SR17, SR20, SR29, SR35, SR41.

| changed function | callers at 8eea2d15 |
|---|---|
| `worktreeCommand` | one registry entry |
| `landing.NewLane` | commit.go:151, merge.go:347 |
| `adopt.KitSourceCheckout` | doctor.go:279, doctor.go:439, doctor_rows.go:73, doctor_rows.go:161, link.go:50, land.go:144 |
| `adopt.kitDir` | adopt.go, doctor.go, init.go, link.go, setup.go, upgrade.go |
| `selectAssignment` | path.go:63, merge.go:181, merge.go:273 |
| `landing.TicketsOnlyFolder` | close.go:57, land_identity.go:142 |
| `refreshop.Consume` | worktree.go:701, subshell.go:48 |
| `refreshop.RefreshedStartRef` | loop.go:154, refresh.go:80 |
| `gate.BenchkitLane` | lane_select.go (LaneFor), profile_lane_table_test.go:105, lane_test.go |
| `gate.resolveLane` | lane.go:127, lane_test.go:149 |
| `gate.LaneRequest` | authorization.go:217 |
| `structure.Command` | cmd/bench/main.go:94 |

Reader sweep: the decision facts this pass changes are the leaf set, the lane
mapping, the kit-source owner, the path-match rule, the tickets-only rule,
the refresh rule, the lane check list, and the structure grammar. Their
readers are the callers above plus five checks: `axi-query-registry`,
`subcommand-routing`, `bounds-policy`, `canonical-path-owner`, and
`profile-lane-table`. Each reader has a row or a fence entry. `bench status`
reads the structure violation count and is unchanged (SR56). No `.mjs` script
or workflow file reads any of them. No shipped-surface claim word changes.

Measurement (FT231 owns it, and this pass adds no store). The baseline row is
the `handoff-sections` landing dae9a77e:

| measure | baseline |
|---|---|
| tickets first-pass | 7 of 9 |
| repair rounds | 6 (4 spec-row, 1 delegate, 1 slicing) |
| review findings to repair targets | 14 to 6 |
| coordinator probes that bit | 13 of 15 |
| red gates in folds | 1 in 11 |
| landing gate elapsed | 104 s, 65 s in test |
| census raw calls | 5 at the landing, 36 across 15 worktrees |
| commits, folds, worktrees | 23, 11, 15 |
| session tokens | unavailable |
| survey cost | one Opus delegate, 141,686 tokens, 41 tool uses, 378 s |
| census cost (this spec) | four Opus medium delegates, totals in the backlog file |

Hypotheses. Accuracy: after ticket 1 lands, one tracer-bullet change through
a worktree leaf touches three files rather than seven, and it lands
first-pass. Cost: the survey plus this pass is repaid when two later
leaf-touching slices each cost fewer tokens than the baseline slice.

Every write delegate runs `opus` at low or medium effort. The review round
runs `sonnet` at xhigh effort by the reviewer's `--reviewer sonnet xhigh`
override of 2026-09-04.

Approval: on 2026-09-04 the reviewer accepted the spec, the eight tickets, and
decisions (a) to (p) as recorded, with the words "I accept all decisions".
