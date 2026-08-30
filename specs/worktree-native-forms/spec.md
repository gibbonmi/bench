# Worktree native forms

Status: implemented

Decision source: named reviewed artifact `roadmap/FT254.md`, the slice 2 paragraph the 2026-08-29 drain settled at `fd52f3df`. Its provenance is the retired map `fd52f3df^:specs/worktree-exec-comfort/decisions/worktree-exec-comfort.md` (#4, #7, #8, #12, #13).

Verification log: 2 iteration(s) to accept — round 1 (codex gpt-5.6-sol / high) returned 6 blocking and 5 non-blocking findings. All folded, one variant on the seal, and round 2 verified the folds; see Further notes.

## Problem

An agent in an assignment worktree still leaves the Bench command boundary for
four steps. The worktree help names `bash bin/bench.sh gate --fresh` as the gate
form. A worktree's new grammar has no build form, so a session hand-builds
`dist/bench` with `scripts/go-build.sh` and `BENCH_HOME`. The system suite runs
only inside the gate or under a hand-built environment. A dependent ticket starts
at a sibling's tip only through a raw checkout, and the preflight `base-current`
red names no remedy.

## Solution

Slice 2 of FT254 gives each step one Bench form. `bench worktree build <target>`
builds the worktree's tree into `<worktree>/dist/bench` with its seal. Its next
action is `bench worktree exec '<label>' -- ./dist/bench <verb>`. Both help
inventories name `bench worktree exec <target> -- bench gate` instead of the raw
gate line. `bench test --check system` runs the gate's system phase with the
owned run binary and `BENCH_KIT`.

`bench worktree create --from <target>` starts a new assignment at a sibling's
committed tip. Every preflight check row carries a `next` cell, and the
`base-current` red names `bench worktree merge --from <default-branch> <target>`.

## User stories

### Group A — `bench worktree build`

Line: opus / low. The verb composes the shared target resolver, the shared
refusal printer, and the exported `runbinary.Build` seam. The seams have exact
tests, and the ordinary test phase covers the package.

1. As an agent, I want `bench worktree build <target>` to build the worktree's tree into `<worktree>/dist/bench`, so that the worktree's own grammar has an executable.
2. As an agent, I want the build to run the sanctioned build script through `runbinary.Build`, so that the executable carries its seal and version stamp.
3. As an agent, I want the success output to name the assignment and the executable path, so that I can read what was built.
4. As an agent, I want the success output to end with `bench worktree exec '<label>' -- ./dist/bench <verb>`, so that the next command is copy-paste.
5. As an agent, I want a quoted, globbed, or multi-line label rendered paste-safe, so that the next line never breaks the shell.
6. As an agent, I want a build failure to print the builder's error and then `worktree: <absolute path>`, so that I never need the raw path.
7. As an agent, I want a missing `go` on `PATH` refused with the builder's own sentence, so that the refusal names the tool.
8. As an agent, I want an interrupted build to exit 130 with the `worktree:` line, so that an interrupted run leaves me the path.
9. As an agent, I want an unassigned target refused with `next=bench worktree list`, so that `build` shares one refusal with `path`, `exec`, and `show`.
10. As an agent, I want a second build to replace the first executable, so that a rebuild after an edit needs no clean step.
11. As an agent, I want `bench worktree exec <target> -- ./dist/bench <verb>` to run the built file, so that the exec form reaches the artifact.
12. As a reviewer, I want the build to write nowhere but `dist/`, so that the declared build output keeps the landing green.
13. As a cold agent, I want `build` listed in `bench worktree --help` and `bench help`, so that discovery needs no guesswork.

### Group B — the gate form in the help

Line: opus / low. The change is two inventory strings with exact fixtures.

14. As a cold agent, I want `bench worktree --help` to end with `bench worktree exec <target> -- bench gate`, so that the worktree help names no raw gate path.
15. As a cold agent, I want `bench help` to carry that same form under the worktree rows, so that the two inventories agree.
16. As a reviewer, I want the `.bench/gate.sh` wrapper-less refusal unchanged, so that the gate entry keeps its own contract.

### Group C — `bench test --check system`

Line: opus / medium. The check crosses the gate's phase producer and the focused
run's environment, and a wrong environment reads as a suite red.

17. As an agent, I want `bench test --check system` to run `go test -trimpath -count=1 -json -tags=system ./internal/systemtest`, so that a hand run needs no hand-built environment.
18. As an agent, I want the check to carry `BENCH_RUN_BINARY`, `BENCH_KIT`, and `BENCH_SYSTEM_ROOT` the way the gate sets them, so that the suite's owner starts.
19. As a reviewer, I want the phase table and the named check to read one system-phase producer, so that the two argv lists cannot drift.
20. As a reviewer, I want the check to carry no conformance scope variable, so that the system suite never reads as a conformance run.
21. As an agent, I want `bench test --check system` in a linked repo refused with a structured error, so that the kit-only suite never grades a foreign root.
22. As an agent, I want the check to render the focused report at exit 0 or 1, so that it reads like every focused run.
23. As a reviewer, I want the check to write no gate-owned record, so that a focused run never moves a green marker.

### Group D — `bench worktree create --from <target>`

Line: opus / low. The flag composes the sibling lookup the merge verb owns and the
start operand `createAt` already takes.

24. As an agent, I want `bench worktree create --from <target>` to start the new worktree at the sibling's committed tip, so that a dependent ticket starts in one step.
25. As an agent, I want the ledger `start` to record that tip, so that the landing's reviewed range starts there.
26. As an agent, I want the new worktree on its own assignment branch, so that the sibling's branch stays untouched.
27. As an agent, I want a `--from` that names no active assignment refused with `next=bench worktree list`, so that a typo names the lookup.
28. As an agent, I want a dirty sibling refused with `next=bench worktree exec <id> -- bench commit`, so that uncommitted work is never dropped.
29. As an agent, I want a sibling whose checkout is not at its branch tip refused, so that the start is always a committed tip.
30. As an agent, I want `--from` with `--refresh` refused at exit 2 before any refresh runs, so that the two start rules never compete.
31. As an agent, I want an ambiguous `--from` prefix refused with the colliding ids, so that I can retype the address.
32. As an agent, I want a `--from` with control characters refused before any lookup, so that the ledger read never sees it.
33. As a cold agent, I want `[--from <target>]` in the create grammar and `bench worktree --help`, so that help and parser agree.
34. As a reviewer, I want the sibling's creation bundle proved before its tip seeds a start, so that a foreign registration never starts a worktree.

### Group E — the preflight `next` column

Line: opus / low. `Decide` is a pure table function with exact tests, and the
gatherer gains two facts at a known seam.

35. As an agent, I want every preflight check row to carry a `next` cell, so that a red names its remedy in the row.
36. As an agent, I want the `base-current` red `default branch tip is not an ancestor of HEAD` to carry `bench worktree merge --from <default-branch> <target>`, so that the remedy is copy-paste.
37. As an agent, I want `<default-branch>` filled with the resolved default branch name, so that the remedy needs no lookup.
38. As an agent, I want `<target>` filled with the assignment id when the preflight root is its worktree, so that the remedy names my worktree.
39. As an agent, I want the literal `<target>` kept when the root is no active assignment worktree, so that the remedy stays honest.
40. As a reviewer, I want every green and not-applicable row to carry an empty `next`, so that the column adds nothing to a green run.
41. As a reviewer, I want every other red to keep an empty `next`, so that this slice decides one remedy only.
42. As a reviewer, I want a symlinked preflight root to match its assignment, so that macOS temp roots resolve the same target.

### Group F — inventories and guidance

Line: opus / medium. The group edits the changelog and the glossary. The leverage
rule routes guidance prose to the mid tier, and the reviewer's 2026-08-26 direction
caps every subagent at medium.

43. As a reader, I want `CHANGELOG.md` to record the four added forms and the replaced gate line, so that the release notes stay current.
44. As a reader, I want `CONTEXT.md` to define **worktree build** and to name the system check, so that the glossary matches the verbs.

## Implementation decisions

- `bench worktree build <target>` is a worktree family verb. Its grammar is
  `usage.WorktreeBuild`, it joins `worktreeCommands`, the `bench help` inventory,
  the kept-routes list, and the family dispatch.
- The verb resolves the target through `resolveWorktree` and refuses through
  `printTargetRefusal` with the verb name `bench worktree build`.
- The verb builds through a `build` join on the worktree `joins` set whose default
  is `runbinary.Build`. The call is
  `build(ctx, <worktree>, <worktree>/dist/bench)` under `subprocess.NotifyCancel`.
  The script writes the executable and its seal, so the verb writes no seal itself.
- Success prints one `worktree_build[1]{worktree,executable}` table with the
  assignment id and the absolute executable path. Then it prints `next[1]:` and
  one line `bench worktree exec <label> -- ./dist/bench <verb>`. A line-safe label
  renders through `sanitize.ShellQuote`. A label with a control byte gives way to
  the assignment id, by the `mergeReconcileNext` precedent.
- A build error prints `bench worktree build: <error>` and then the `worktree:`
  line through `nameWorktree`, at exit 1. A cancelled context exits 130 through
  the same line.
- The wrapper's resolution order is unchanged. `kit_dir` resolves a cwd inside a
  worktree of the same repository to that worktree, and `bench_binary_path` tries
  `<kit>/dist/bench` first. So a bare `bench` whose cwd is a built worktree serves
  that worktree's build, and a bare `bench` elsewhere serves the main checkout's.
  The exec form `./dist/bench` is the cwd-independent spelling. See Further notes.
- The worktree usage trailer becomes `bench worktree exec <target> -- bench gate`.
  The `bench help` row `bash bin/bench.sh gate --fresh` moves under the worktree
  family as `bench worktree exec <target> -- bench gate  run one active owned
  worktree's gate`. The gate family keeps its two other rows.
- `internal/gate` exports one system-phase producer: the operands
  `-tags=system ./internal/systemtest` and the environment
  `BENCH_SYSTEM_ROOT=<root>`. `BenchkitPhases` and the `system` named check both
  consume it.
- `bench test --check <name>` accepts the dev-tier conformance names and `system`.
  The `system` check runs `focusedTestArgv` over the producer's operands in the
  kit, under `selectedRunEnvironment` plus the producer's environment. It sets no
  conformance root, tier, or scope variable.
- The `system` check requires the graded root to be the kit. A root that differs
  from `BENCH_KIT` returns `toon.Errorf("system check unavailable", "the system
  suite grades the kit checkout only")` at exit 1 and starts no Go.
- `bench worktree create` gains `--from <target>`. The value resolves through the
  sibling lookup the merge verb owns, with no target to exclude. The lookup
  accepts an active assignment on its branch, with a valid creation bundle, a
  clean checkout, and `HEAD` at the branch tip. The resolved tip is the
  `requestedStart` operand of `createAt`.
- `--from` composes no commit lookup. A spelling that names no active assignment
  refuses with `--from names no active assignment`.
- Every `--from` refusal prints through `printTargetRefusal` with the verb name
  `bench worktree create`. So the default `next=` is `bench worktree list`, and a
  dirty sibling carries `bench worktree exec <id> -- bench commit`.
- `--from` with `--refresh` is a usage refusal that names `--from`, at exit 2,
  before `refreshop.Consume` runs.
- `CheckResult` gains `Next`. The preflight table is
  `checks[N]{check,verdict,detail,next}`. `green` and `notApplicable` rows carry an
  empty `Next`.
- `Facts` gains `DefaultBranch` and `AssignmentTarget`. The gatherer fills
  `DefaultBranch` from `git.ResolvedDefault`. It fills `AssignmentTarget` with
  the id of the active assignment whose canonical worktree path equals the
  canonical preflight root, and leaves it empty otherwise.
- The one red that carries a `Next` is `base-current` with the detail
  `default branch tip is not an ancestor of HEAD`. Its `Next` is
  `bench worktree merge --from <DefaultBranch> <AssignmentTarget>`, with the
  literal `<target>` when `AssignmentTarget` is empty.
- The `next` cell passes through the same representability refusal the detail
  cell passes through, because `toon.Table` refuses a control byte.

## Bootstrap authority before execution

The build verb executes the target tree's `scripts/go-build.sh`. The chain is the
installed wrapper, the Bench executable it selects, the assignment ledger, and the
creation bundle. The bundle authenticates the worktree path before the script runs.
The gate authors its private run binary from the graded tree the same way. The
built `dist/bench` is a development artifact that no gate, lane, or `bench test`
run reuses, because each owns a private exact-source build.

The bundle authenticates ownership of the tree, not the script's content. So an
active assignment authorizes the execution of a candidate-controlled build script,
and the reviewer accepts that assumption at sign-off. The seal the script writes
is self-attestation, as `projects/benchkit.md` records for every hand build.

## Testing decisions

- A good test drives the public command function with an argv and a real Git
  repository or a real child. It reads stdout, stderr, and the exit code.
- The build rows drive `BuildCommand` through the `joins` seam by the precedent of
  `mergeWith`. A fake `build` join records its arguments and writes a marker file.
  One row drives the production join against a throwaway worktree whose
  `scripts/go-build.sh` is a stub that writes its second argument.
- The exec-form row drives `runWorktreeChild` with a stub `dist/bench` in the
  worktree, by the precedent of `childOutput` in `internal/worktree/exec_test.go`.
- The refusal rows attach to `internal/worktree/identifier_operand_test.go` beside
  `TestTargetVerbsShareOneRefusalPrinter`.
- The inventory rows attach to `cmd/bench/main_test.go`
  (`TestHelpInventoryIsComplete`) and `cmd/bench/command_registry_test.go`
  (`TestKeptWorktreeOperationsKeepTheirGrammar`).
- The system check rows attach to `internal/testreport/check_test.go` beside
  `TestNamedCheckOwnsConformanceEnvironment`, which drives a fake `go` that records
  its argv and environment. The record row attaches to
  `internal/testreport/testreport_test.go` beside
  `TestExplicitFocusedRunsWriteNoGateOwnedRecords`. The phase-table row attaches to
  `internal/gate/phases_test.go`.
- The `--from` rows attach to `internal/worktree/worktree_test.go` beside
  `TestCreateCommandPrintsNextHint`, with the sibling fixtures
  `internal/worktree/merge_test.go` builds for `TestMergeFoldsASiblingBranchTip`
  and `TestMergeRefusesADirtyOrDetachedSibling`.
- The preflight rows attach to `internal/preflight/decision_test.go` beside
  `TestDecideBaseCurrent`, to `internal/preflight/command_review_test.go` beside
  `TestCommandStaleBase`, and to `internal/preflight/gather_test.go`.
- The gate observes the feature in the `test` phase for every package above.

### Seam diagram

    trigger: an agent's Bash call, or a delegate's charge
        │
        ▼
    argv  ──▶  [ usage.Parse: build, create, test, preflight grammars ]  ──▶  Result
                  ◀ tests attach here: each command test drives Parse through its verb
        │
        ▼
    target  ──▶  [ resolveWorktree + printTargetRefusal ]  ──▶  path, or refusal + next=
                  ◀ tests attach here: identifier_operand_test reads stderr lines
        │
        ▼
    path  ──▶  [ BuildCommand: joins.build → runbinary.Build → scripts/go-build.sh ]  ──▶  dist/bench, table, next[1]
                  ◀ tests attach here: build_test drives the join, one row runs a stub script
        │
        ▼
    --check system  ──▶  [ testreport: gate system producer + selectedRunEnvironment ]  ──▶  go test -json child
                  ◀ tests attach here: check_test reads the fake go's argv and env
        │
        ▼
    --from  ──▶  [ sibling lookup → createAt(requestedStart) ]  ──▶  worktree at the sibling tip
                  ◀ tests attach here: worktree_test reads HEAD, the ledger, and stderr
        │
        ▼
    Facts  ──▶  [ preflight.Decide ]  ──▶  checks{check,verdict,detail,next}
                  ◀ tests attach here: decision_test drives Decide, command tests read the table

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| WF1 | 1, 2 | `bench worktree build <label>` calls the `build` join once with the canonical worktree path and `<worktree>/dist/bench` | `internal/worktree/build_test.go` (`TestBuildCallsTheJoinWithTheWorktreeAndOutput`) | no verb exists, and a hand `go build` would skip the script |
| WF2 | 2 | the production `build` join, against a worktree whose `scripts/go-build.sh` stub writes a marker to its second argument, leaves that marker byte-equal at `dist/bench` and exits 0 | `internal/worktree/build_test.go` (`TestBuildRunsTheWorktreeBuildScript`) | a verb that runs `go build` itself, or a default join that never reaches the tree's script, leaves no marker |
| WF3 | 3 | success prints `worktree_build[1]{worktree,executable}:` with the assignment id and the absolute `dist/bench` path | `internal/worktree/build_test.go` (`TestBuildPrintsTheTable`) | a bare exit 0 leaves the reader with no path |
| WF4 | 4, 5 | for the label `it's a*b` success ends with `next[1]:` and `  bench worktree exec 'it'\''s a*b' -- ./dist/bench <verb>`, and for a label that holds a newline that line carries the assignment id instead | `internal/worktree/build_test.go` (`TestBuildNamesTheExecFormForTheLabel`) | a label pasted inside literal quotes breaks on a quote, a glob, or a newline |
| WF5 | 6 | a `build` join that returns an error prints `bench worktree build: <error>` then `worktree: <absolute path>` on stderr at exit 1 | `internal/worktree/build_test.go` (`TestBuildFailureNamesTheWorktree`) | a failure without the path sends the agent to a raw command |
| WF6 | 7 | with a `PATH` that holds no `go`, the production join refuses at exit 1 with stderr that holds `Go is absent from PATH` and the `worktree:` line | `internal/worktree/build_test.go` (`TestBuildRefusesWithoutGoOnPath`) | a verb that shells to `go` directly prints an `exec` error instead |
| WF7 | 8 | a `build` join that returns `context.Canceled` exits 130 and stderr ends with the `worktree:` line | `internal/worktree/build_test.go` (`TestBuildCancelExitsOneHundredThirty`) | a cancel mapped to 1 hides the interrupt |
| WF8 | 9 | `bench worktree build no-such-label` prints `bench worktree build: target is unassigned` then `next=bench worktree list` at exit 1 | `internal/worktree/identifier_operand_test.go` (`TestTargetVerbsShareOneRefusalPrinter`) | a fourth printer describes one failure a fourth way |
| WF9 | 10 | two builds in a row both exit 0 and `dist/bench` holds the second stub's bytes | `internal/worktree/build_test.go` (`TestBuildReplacesAPriorExecutable`) | a verb that refuses an existing output needs a clean step |
| WF10 | 11 | `runWorktreeChild` with argv `./dist/bench version` in a worktree whose stub `dist/bench` prints `worktree-grammar` emits that word on stdout at exit 0 | `internal/worktree/exec_test.go` (`TestExecRunsTheWorktreeDistBench`) | a child started outside the worktree cannot find the relative path |
| WF11 | 12 | after the production-join build every untracked path in the worktree sits under `dist/` | `internal/worktree/build_test.go` (`TestBuildWritesOnlyUnderDist`) | an output outside the declared `dist/` makes the landing refuse residue |
| WF12 | 13 | `bench worktree --help` names `bench worktree build <target>` | `cmd/bench/command_registry_test.go` (`TestKeptWorktreeOperationsKeepTheirGrammar`) | the kept-routes list names each verb, so a missing `build` reds |
| WF13 | 13 | `bench help` prints the `build` row, byte-equal to the fixture | `cmd/bench/main_test.go` (`TestHelpInventoryIsComplete`) | the fixture is exact, so a registry row without its fixture row reds |
| WF14 | 14 | `bench worktree --help` ends with `bench worktree exec <target> -- bench gate` and holds no `bin/bench.sh` | `cmd/bench/command_registry_test.go` (`TestWorktreeHelpNamesTheExecGateForm`) | an appended line beside the old trailer keeps the raw path alive |
| WF15 | 16 | `.bench/gate.sh` still refuses a wrapper-less entry and names `bash bin/bench.sh gate` | `internal/conformance/gate_entry_test.go` (`gateEntryWrapperAction`) | a sweep that deletes every `bin/bench.sh gate` string breaks the gate entry |
| WF16 | 17 | `bench test --check system` runs the fake `go` with `test -trimpath -count=1 -json -tags=system ./internal/systemtest` | `internal/testreport/check_test.go` (`TestSystemCheckOwnsTheGateEnvironment`) | the registry lookup refuses `system` today |
| WF17 | 18 | that run's environment holds `BENCH_RUN_BINARY=<selected>`, `BENCH_KIT=<root>`, and `BENCH_SYSTEM_ROOT=<root>` | `internal/testreport/check_test.go` (`TestSystemCheckOwnsTheGateEnvironment`) | the suite's owner refuses without the three names |
| WF18 | 19 | `BenchkitPhases(kit, kit)` holds a `system` phase whose argv ends with the producer's operands and whose env holds `BENCH_SYSTEM_ROOT=<kit>` | `internal/gate/phases_test.go` (`TestBenchkitPhasesSystemPhaseReadsTheProducer`), plus review-owned source inspection that both consumers call the exported producer | a second literal in the phase table drifts from the check, and only the Standards axis sees a duplicated literal |
| WF19 | 20 | that run's environment holds no `BENCH_CONFORMANCE_SCOPE`, `BENCH_CONFORMANCE_ROOT`, or `BENCH_CONFORMANCE_TIER` | `internal/testreport/check_test.go` (`TestSystemCheckOwnsTheGateEnvironment`) | a check routed through `conformanceEnvironment` scopes the suite |
| WF20 | 21 | with `BENCH_KIT` naming a directory other than the root, `bench test --check system` prints `system check unavailable` at exit 1 and runs no `go` | `internal/testreport/check_test.go` (`TestSystemCheckRefusesAForeignRoot`) | a linked repo would run the kit's suite against a foreign root |
| WF21 | 22 | a fake `go` that emits one `fail` event for the package `checkfixture` makes `bench test --check system` exit 1 with a `packages[1]{package,status}` table whose row is `checkfixture,fail` | `internal/testreport/check_test.go` (`TestSystemCheckReportsAFailingSuite`) | a check that ignores the child's status reads red as green |
| WF22 | 23 | after `bench test --check system` the gate cache and the lane record are unchanged | `internal/testreport/testreport_test.go` (`TestExplicitFocusedRunsWriteNoGateOwnedRecords`) | a run that writes a verdict makes a focused run a gate |
| WF23 | 24, 25, 26 | `create --from <sibling>` yields a worktree whose `HEAD` equals the sibling's branch tip, whose ledger `start` equals it, and whose checked-out branch equals `intent.AssignmentBranchRef(owner, id)` for the new record | `internal/worktree/worktree_test.go` (`TestCreateFromStartsAtTheSiblingTip`) | a create that ignores the flag starts at the default tip |
| WF24 | 27 | `--from no-such-label` prints `bench worktree create: --from names no active assignment` then `next=bench worktree list` at exit 1, and the ledger gains no record | `internal/worktree/worktree_test.go` (`TestCreateFromRefusesAnUnknownSibling`) | a fallthrough to `HEAD` creates a worktree at the wrong start |
| WF25 | 28 | a sibling with an uncommitted edit refuses with `sibling checkout is not clean` and `next=bench worktree exec <id> -- bench commit` | `internal/worktree/worktree_test.go` (`TestCreateFromRefusesADirtyOrDetachedSibling`) | a start at the tip silently drops the sibling's edit |
| WF26 | 29 | a sibling whose checkout is detached refuses with `sibling is not on its assignment branch` | `internal/worktree/worktree_test.go` (`TestCreateFromRefusesADirtyOrDetachedSibling`) | a lookup by branch alone misses a detached checkout |
| WF27 | 30 | `--refresh --from <sibling>` returns the usage line naming `--from` at exit 2, and `refreshop.Consume` runs nothing | `internal/worktree/worktree_test.go` (beside `TestCreateCommandHelpPerformsNoRefresh`) | a refresh before the refusal moves the default branch |
| WF28 | 31 | two siblings whose labels share a 9-character prefix make `--from <prefix>` refuse with `target is ambiguous: ` and both ids | `internal/worktree/worktree_test.go` (`TestCreateFromRefusesAnAmbiguousPrefix`) | a first-match lookup starts at the wrong sibling |
| WF29 | 32 | `--from $'a\x01b'` refuses with `--from contains control characters` and reads no ledger, proved by a malformed ledger file that would otherwise refuse first | `internal/worktree/worktree_test.go` (`TestCreateFromRefusesControlBytes`) | an unchecked value reaches the selector |
| WF30 | 27 | `--from <landed-sibling>` whose state is not active refuses with `--from names no active assignment` | `internal/worktree/worktree_test.go` (`TestCreateFromRefusesARetiredSibling`) | a lookup over every state starts at a retired tip |
| WF31 | 33 | `bench worktree create --help` prints `usage: bench worktree create [--refresh] --request <opaque-id> --label <work-item> [--from <target>]` at exit 0 | `internal/worktree/worktree_test.go` (`TestCreateCommandAnswersHelpSpellings`) | the help fixture is exact |
| WF32 | 35 | `bench preflight review <slug>` on a conformant tree prints `checks[5]{check,verdict,detail,next}` | `internal/preflight/command_review_test.go` (`TestCommandConformantTree`) | the old header survives a `Next` field the renderer never emits |
| WF33 | 36, 37, 38 | `Decide` with `DefaultBranchCurrent=false`, `DefaultBranch=main`, and `AssignmentTarget=<id>` yields the `base-current` red with `Next` equal to `bench worktree merge --from main <id>` | `internal/preflight/decision_test.go` (`TestDecideBaseCurrent`) | a red without a remedy is today's behavior |
| WF34 | 39 | the same red with an empty `AssignmentTarget` yields `Next` equal to `bench worktree merge --from main <target>` | `internal/preflight/decision_test.go` (`TestDecideBaseCurrent`) | a renderer that prints an empty id prints a broken command |
| WF35 | 40 | every green row and every not-applicable row from `Decide` carries an empty `Next` | `internal/preflight/decision_test.go` (`TestDecideAllGreen`) | a default filled on every row advertises a remedy for nothing |
| WF36 | 41 | every other red `Decide` produces carries an empty `Next`: the `base-current` reds `default branch does not resolve` and `source base does not resolve`, both `tip-current` reds, both `paths-authorized` reds, the `rows-owned` red, the `rows-membership` red, and both `diff-nonempty` reds | `internal/preflight/decision_test.go` (`TestDecideOtherRedsCarryNoNext`) | a remedy copied onto every red names a merge for an unresolved branch or a fence miss |
| WF37 | 36 | `bench preflight review <slug>` behind the default branch prints the row `base-current,red,default branch tip is not an ancestor of HEAD,bench worktree merge --from main <target>` at exit 1 | `internal/preflight/command_review_test.go` (`TestCommandStaleBase`) | the pure decision can pass while the gatherer never fills the facts |
| WF38 | 38 | `Gather` from the worktree path of an active assignment written through `intent.PutAssignment` fills `AssignmentTarget` with that assignment's id | `internal/preflight/gather_test.go` (`TestGatherAssignmentTarget`) | a gatherer that never reads the ledger leaves the placeholder |
| WF39 | 42 | `Gather` from a symlink to that worktree path fills the same id | `internal/preflight/gather_test.go` (`TestGatherAssignmentTarget`) | a raw string compare misses a resolved registration |
| WF40 | 39 | `Gather` from the primary checkout leaves `AssignmentTarget` empty | `internal/preflight/gather_test.go` (`TestGatherAssignmentTarget`) | a gatherer that picks the first active assignment names a stranger |
| WF41 | 35 | a preflight whose `--source-tip` is pinned prints `checks[6]{check,verdict,detail,next}` | `internal/preflight/source_tip_test.go` (`TestSourceTipAcceptedByBothModes`) | the pinned-tip fixture holds the old header |
| WF42 | 33 | `bench worktree --help` names `bench worktree create [--refresh] --request <opaque-id> --label <work-item> [--from <target>]` | `cmd/bench/command_registry_test.go` (`TestKeptWorktreeOperationsKeepTheirGrammar`) | the kept list holds the bare verb today, so a grammar line without the flag passes |
| WF43 | 34 | a mutated owner marker or lock on the sibling makes `--from <sibling>` refuse with that component's named detail, a mutated assignment state refuses with `--from names no active assignment` because the active filter runs before the bundle, and the ledger gains no record in each case | `internal/worktree/worktree_test.go` (`TestCreateFromRefusesAFailedSiblingIdentityComponent`) | a lookup that checks state, cleanliness, and the branch alone skips the bundle |
| WF44 | 15 | `bench help` prints `bench worktree exec <target> -- bench gate  run one active owned worktree's gate` and no `bash bin/bench.sh` row, byte-equal to the fixture | `cmd/bench/main_test.go` (`TestHelpInventoryIsComplete`) | an appended row beside the old gate row keeps the raw path alive |
| WF45 | 24 | a `create --from <sibling>` replay with the same `--request` after the sibling takes an uncommitted edit exits 0 with the existing record, and the ledger gains no second record | `internal/worktree/worktree_test.go` (`TestCreateFromReplayReturnsTheRecord`) | a `--from` resolved before the replay lookup refuses a no-op replay |
| WF46 | 17 | the conformance registry names no check equal to `gate.SystemPhaseName` | `internal/testreport/check_test.go` (`TestSystemCheckNameIsReserved`) | a registry check named `system` would be shadowed with no refusal |

Not covered: story 43 — the changelog entry is prose. The review round grades it
against WF1, WF14, WF16, WF23, and WF32.
Not covered: story 44 — the glossary entry is prose. The review round grades it
against WF1 and WF16.

### Edge inventory

- A target with control bytes, an unassigned target, an ambiguous target, an
  inactive assignment, or a missing tree. Each keeps its existing refusal and
  `next=` line through the shared printer (WF8).
- A worktree path with a space or a glob character. The `worktree:` line prints it
  raw by the existing rule, and the table cell renders it because `toon.Table`
  refuses control bytes only.
- A `dist` entry that is a symlink, a directory-shaped file, or unwritable. The
  script's `refuse_output` owns that refusal, and WF5 proves the verb prints the
  builder's error.
- The seal beside `dist/bench`. The script's subject mode writes it through
  `freshness-publish`, and `runbinary.Own` verifies that pair in every gate run.
  WF2 pins the script as the one producer, so no row asserts the seal bytes.
- A build that a signal interrupts: WF7. The script's own traps remove its staged
  file.
- `--check system` with `--run`, `--package`, or `--changed`: the existing `--check`
  conflict rule refuses at exit 2.
- `--check system --full`: the full report renders by the existing rule.
- A missing `go` for the system check: the existing structured start error.
- `--from ''`: the existing `NoEmptyValue` rule refuses at exit 2.
- `--from` that names the request's own existing assignment on a replay. The
  idempotent replay returns the existing record, because the ledger records only
  `start`.
- A `--from` value that Git can peel to a commit but that names no assignment:
  WF24, because the flag composes no commit lookup.
- A preflight root whose default branch name or assignment id holds a control
  byte: the `next` cell meets the same representability refusal as the detail cell.
- A preflight run with `--base`: the explicit range grades `base-current`, and its
  red keeps an empty `Next` (WF36).
- **Won't handle** the wrapper's resolution order after a build — `kit_dir` and
  `bench_binary_path` are unchanged, and the `./dist/bench` exec form survives.
- **Won't handle** a `build` action row in `bench worktree list` — `list` keeps its
  `path` and `exec` actions, and WF12 gives the verb its discovery.
- **Won't handle** `create --from <commit>` — `--refresh` starts at the fresh default
  tip, and the sibling form survives.
- **Won't handle** a landing of a `--from` worktree before its sibling lands — the
  landing's base rule stands, and a fold after the sibling lands survives.
- **Won't handle** the system check under `bench test --changed` — the changed
  selection is a package list, and `--check system` survives.
- **Won't handle** a `next` for the other preflight reds — the map decided one
  remedy, and each red keeps its detail.
- **Won't handle** the `bench test --help` check inventory and the unknown-check
  sentence — FT270 owns them, and `--check system` survives.

## Ownership fences

- `internal/usage/worktree.go`
- `internal/worktree/build.go`
- `internal/worktree/build_test.go`
- `internal/worktree/land.go`
- `internal/worktree/worktree.go`
- `internal/worktree/worktree_test.go`
- `internal/worktree/merge.go`
- `internal/worktree/merge_test.go`
- `internal/worktree/exec_test.go`
- `internal/worktree/identifier_operand_test.go`
- `internal/worktree/path.go`
- `internal/worktree/ownership.go`
- `internal/worktree/parallel_census_test.go`
- `internal/gate/phases.go`
- `internal/gate/phases_test.go`
- `internal/testreport/command.go`
- `internal/testreport/check_test.go`
- `internal/testreport/testreport_test.go`
- `internal/preflight/`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `cmd/bench/command_registry_test.go`
- `CHANGELOG.md`
- `CONTEXT.md`
- `specs/worktree-native-forms/`

## Out of scope

- Slice 3 of FT254, `bench worktree resolve <target> <path>...` after FT263 lands —
  6 edits, 2 gate runs.
- An ambient worktree-root or label-expanding variable for the exec child — 3
  edits, 1 gate run. The FT254 row states it as an `or`, so it needs a reviewer
  decision before a spec.
- `create --from <commit>` through the merge verb's commit lookup — 2 edits, 1
  gate run.
- A `build` action row in `bench worktree list` — 2 edits, 1 gate run.
- A `next` for the remaining preflight reds — 2 edits, 1 gate run.
- The `bench test --help` check inventory, the unknown-check sentence, and the
  prose lane check — FT270 owns them, 4 edits, 1 gate run.

## Further notes

The retired map's #7 states that the wrapper on `PATH` keeps serving the main
checkout's build for a bare `bench`. That sentence held while no worktree carried
a `dist/bench`. The wrapper's `kit_dir` resolves a cwd inside a worktree of the
same repository to that worktree, and `bench_binary_path` tries `<kit>/dist/bench`
before the main tree's. After `bench worktree build`, a bare `bench` whose cwd is
that worktree therefore serves the worktree's build.

This spec keeps the wrapper
unchanged and records the observed rule. The reviewer can veto that call at
sign-off. A wrapper change that pins the main build is a separate capability.

FT270 names a `--system` flag for the same suite. The map's #12 chose the named
check `bench test --check system` with no new flag, and this spec follows #12.
The next drain reconciles the FT270 sentence.

Review round 1 (codex `gpt-5.6-sol` / high) folded these findings:

- The seal finding took a variant: the script is the seal's one producer, and WF2
  now pins that the script produced `dist/bench`. The edge inventory records why
  no row asserts the seal bytes.
- The next line shell-quotes a line-safe label and falls back to the id (WF4).
- The new WF43 proves the sibling bundle. WF36 enumerates every other red.
- WF13 split from the new WF44, so the build ticket owns its fixture row.
- The three worktree tickets carry no `Blocked by:` edge; the coordinator
  serializes them for their shared files.
- The bootstrap section records the candidate-script assumption for the reviewer.
- WF18 adds review-owned source inspection. WF21, WF23, and WF37 pin exact
  outputs. Thirty rows name their planned test function.

The `Roadmap:` line is absent on purpose. FT254 keeps slice 3, so this slice does
not retire the row.
