# FT254 exec-comfort research

Four read-only research delegations ran on 2026-08-29 against commit `a712f84f`.
Each section answers one ticket of `specs/worktree-exec-comfort/decisions/worktree-exec-comfort.md`.
Line citations bind to that commit. Re-resolve them before a spec reads this asset.

## #1 Exec surface as built

### Findings

- Argument parsing. `worktreeExecGrammar` declares no flags
  (`internal/worktree/exec.go:20-26`). It sets `MinArgs: 2`, `MaxArgs: -1`, and
  `ReservedPositionalsBeforeTerminator: 1`. In `usage.Parse` the help, terminator,
  and flag pass is gated on `!endedFlags` (`internal/usage/parse.go:99-139`).
  So a dash-prefixed token after `--` falls through to the positional append at
  `parse.go:156`.
- Argument parsing, the two guards. Two guards sit outside that block and still bite
  past the terminator. They are the empty-positional refusal (`parse.go:144-149`) and the
  `MaxArgs` check (`parse.go:153`, inert at `-1`). So `rg -N '' README.md` is refused
  as `unknown argument: ""`. `ExecCommand` then re-asserts the shape
  (`exec.go:35-38`) and prints the bare usage line on exit 2.
- Exit codes. A child's exit code passes through `childExitCode` (`exec.go:115-127`).
  Its stderr is the child's own fd (`exec.go:50`). The grammar refusal also exits 2
  (`exec.go:33`, `:37`). So only the stderr prefix `usage: bench worktree exec`
  separates a child usage error from an exec refusal.
- stdin. `main.go:75` builds `Command{Stdin: os.Stdin}`. `worktreeCommand` passes it
  to `ExecCommand` (`cmd/bench/main.go:586`), then to `cmd.Stdin` (`exec.go:50`).
  No test covers it. Every `exec_test.go` caller passes `nil` (`exec_test.go:30`),
  as do `identifier_operand_test.go:166,188`.
- Environment. `execEnv` strips `BENCH_KIT` and `BENCH_WRAPPER`
  (`internal/env/wrapper.go:10,15-24`). It strips `BENCH_RUN_BINARY`
  (`internal/runbinary/runbinary.go:22`) and the inherited `BENCH_HOME`
  (`exec.go:112` via `internal/capability/capability.go:106-115`). It sets
  `BENCH_HOME` to the verb-resolved home. When `<worktree>/bin/bench.sh` is a
  regular file, it sets `BENCH_WRAPPER` to that path (`exec.go:102-106`). Every
  other caller variable survives, including a leading shell assignment.
- Error paths, part 1. `usage.Parse` lines print at `exec.go:32`. They are
  `unknown argument: "<tok>"` (`parse.go:115,146-148,154`), `missing argument:
  argument` (`parse.go:163`), and the help line on `--help` (exit 0). The shape
  re-assert at `exec.go:36` exits 2.
- Error paths, part 2. `printTargetRefusal` (`internal/worktree/path.go:129-137`)
  exits 1 with the shape `bench worktree exec: <reason>`. The reasons are
  `target contains control characters` (`path.go:46`), `target is unassigned`
  (`:67`), `target is ambiguous: <ids>` (`:76`), and `assignment <id> is not active`
  (`:56-57`). The creation-bundle reasons (`ownership.go:308-320`) are
  `owner marker does not match assignment <id>`, `assignment branch is not checked out`,
  and `worktree registration does not match assignment <id>`.
- Error paths, part 3. A `cmd.Start` failure at `exec.go:53-55` prints
  `bench worktree exec: <os/exec error>` and exits 1. The cancel path
  (`exec.go:62-71`) prints nothing and returns 130. No error path prints the
  absolute worktree path. No error path names `bench worktree list`.
- Retained assignments. `drain-batch` resolves by label (`path.go:82`) and passes the
  active-state check. `validateCreationBundle` (`ownership.go:308-321`) then fails on
  the missing tree at the marker step. Every release retain funnels through
  `retainedReleaseError` (`ownership.go:440-446`). It prints
  `worktree retained (<code>): <reason>` and a `next=` line (`:448-453`). That line is
  `next=bench worktree release --request <request> '<abs path>'`.
- Retain reasons live in `internal/worktree/lifecyclepolicy/lifecyclepolicy.go`.
  The codes and their lines are:
  - `uncertain`: 331, 355, 358, 362, 365, 367, 373, 416, 535.
  - `live-lease`: 346, 348, 352, 504.
  - `ignored`: 375, 377.
  - `active`: 489.
  - `foreign`: 513.
  - `unmerged`: 538.
  - `dirty`: 554-561.

### Tests that pin exec

- `internal/usage/parse_test.go:157` `TestParseReservedPositionalsPrecedeHelpFlagsAndTerminator`.
- `internal/usage/parse_test.go:190` `TestParseBareDashAndSecondDoubleDash`.
- `internal/worktree/identifier_operand_test.go:151` `TestTargetVerbsNameTheResolverReason`
  pins four resolver refusal strings and exit 1.
- `internal/worktree/identifier_operand_test.go:177` `TestTargetVerbsShareOneRefusalPrinter`.
- `internal/worktree/exec_test.go:86,97,109,119,131` pin `BENCH_HOME`, `BENCH_WRAPPER`,
  the `BENCH_RUN_BINARY` strip, the `BENCH_KIT` strip, and unrelated-variable survival.
- `internal/worktree/exec_test.go:141,169,181,209,223` pin the regular-file marker
  predicate and the exact path bytes.
- `cmd/bench/command_registry_test.go:541` `TestKeptWorktreeOperationsKeepTheirGrammar`.
- `internal/worktree/worktree_test.go:120` `TestCreateCommandPrintsNextHint`.
- `internal/worktree/lifecycle_test.go:408` and `lifecycle_policy_test.go:201,57` pin
  the `worktree retained (<reason>)` strings.
- No test pins stdin pass-through. No test pins a non-zero child exit code through
  `childExitCode`.

### Probes

```
$ bench worktree exec ft254-shape -- rg -N '' README.md
exit 2
usage: bench worktree exec (unknown argument: "")

$ bench worktree exec ft254-shape -- rg --bogus-flag README.md
exit 2
rg: unrecognized flag --bogus-flag

$ bench worktree exec ft254-shape -- sh -c 'exit 2'
exit 2

$ bench worktree exec ft254-shape -- cat <<'EOF' ... EOF
hook: BLOCKED: Bench response is bounded, complete, and self-contained. ...

$ FOO=bar bench worktree exec ft254-shape -- sh -c 'echo FOO=$FOO'
exit 0
FOO=bar

$ bench worktree exec ft254-shape -- sh -c 'echo KIT=${BENCH_KIT-unset} WRAPPER=${BENCH_WRAPPER-unset} RUNBIN=${BENCH_RUN_BINARY-unset} HOME_=${BENCH_HOME-unset}'
exit 0
KIT=unset WRAPPER=<worktree>/bin/bench.sh RUNBIN=unset HOME_=/home/mgibs/.bench

$ bench worktree exec drain-batch -- true
exit 1
bench worktree exec: owner marker does not match assignment 78716f6a3828853ee6444406a5d38bcc

$ bench worktree exec no-such-label -- true
exit 1
bench worktree exec: target is unassigned

$ bench worktree exec ft254-shape -- no-such-binary-xyz
exit 1
bench worktree exec: exec: "no-such-binary-xyz": executable file not found in $PATH
```

### Code versus claim

1. stdin passes through: holds, untested.
2. Every failure after target resolution prints the absolute worktree path: does not
   hold (`exec.go:54`, `:71`).
3. An earlier refusal names `bench worktree list`: does not hold (`path.go:129-137`).
4. A flag after `--` belongs to the child: partly holds. The empty-positional guard
   fires past the terminator (`parse.go:144-149`).
5. The child's usage error and stderr pass through unchanged: partly holds. The exit
   2 collision remains.
6. A retained assignment stays resolvable by exec: does not hold for a missing tree
   (`ownership.go:308-320`).

### Bench CLI improvements from the delegation

- `printTargetRefusal` can carry a `next` field the way `retainedReleaseError` does
  (`ownership.go:445-453`).
- The follow-on guard blocks a heredoc into `bench worktree exec ... -- cat`. So the
  stdin path is unprobeable from an agent shell, and only the source is evidence.

## #2 Follow-on guard classifier

### Findings

- Entry. The hook `.bench/hooks/block-bench-follow-on.sh` pipes the PreToolUse JSON to
  `bench guard-bench-follow-on`. The registry row is `cmd/bench/main.go:136` and the
  handler is `main.go:518-539`. It reads only `.tool_input.command`
  (`internal/benchguard/benchguard.go:28-50`). An empty, non-string, or NUL-bearing
  value errors, and `main.go:529-532` warns and allows.
- Parse. `Classify` (`benchguard.go:53-55`) parses once through `shellcommand.Parse`
  (`internal/shellcommand/shellcommand.go:55-90`). Parse deletes heredoc bodies but
  keeps the `<<` operator (`shellcommand.go:56`, `:353-383`). It folds quoted strings
  and backslash escapes into single words (`:298-336`).
- Tokens. Parse classifies each punctuation run as a redirection (`redirectRe`,
  `:50`, `:70-71`) or a control operator (`:73`). A newline becomes `;` (`:66-69`).
  Segments split on control operators only (`:76-83`).
- Bench-segment predicate. Redirections and their operands drop
  (`shellcommand.go:93-108`). Routine prefixes step over (`:111-143`): assignments,
  `env`, `command`, `nohup`, `timeout`, and `xargs`. The head then goes to `isBench`
  (`benchguard.go:115-149`). A head is Bench when it is the literal `bench`, the
  basename `bench` or `bench.sh`, or a PATH or symlink resolution to those. One
  wrapper level (`sh`, `bash`, `zsh` with a `-c` flag) is walked
  (`benchguard.go:67-74`, `:108-114`).
- Refusal predicate. At the first Bench-headed segment, `scan` returns
  `hasOuterSyntax(stream)` (`benchguard.go:64-65`). That function scans the whole
  token stream for any redirection or control operator (`:150-157`). The operator
  need not be adjacent to, or after, the Bench segment. Exit 2 and the message print
  at `main.go:536-538`.
- Refusal text. The string is a constant (`benchguard.go:158-160`). `span` and
  `words` identify the Bench segment (`:58-64`). `hasOuterSyntax` holds the
  operator's `token.Text` at `:152` before it returns a bare bool. To name the token
  is a return-type widening plus the one call site at `main.go:534`.
- Census. `Record` short-circuits on `benchguard.InvokesBench` before any pool-path
  work (`internal/census/census.go:39-41`). So `bench worktree exec L -- sed ...`
  records nothing. `TestRecordSkipsABenchCall` pins it (`census_test.go:137-159`).
  The "raw call" term binds to `census.go:29-52`. The "verb head" term binds to
  `census.go:79-101` and `resolvedHead` (`:105-120`), which appends the `git`
  subcommand at `:114-119`.

### Tests that pin the classifier

- `internal/benchguard/benchguard_test.go:10` `TestClassifyFollowOns` refuses a
  pipeline, `2>&1`, an assignment prefix, a wrapped `||`, and `exec ... > marker`.
  It allows a bare gate, `rg bench`, a quoted `&&` under exec, and a heredoc body
  that holds `bench gate | tail`.
- `benchguard_test.go:29` `TestClassifyResolvesBareAliasFromProcessPath`.
- `benchguard_test.go:44` `TestCommandFromEnvelope`.
- `benchguard_test.go:57` `TestInvokesBenchWalksOneWrapperLevel`.
- `internal/systemtest/bench_follow_on_test.go:13` `TestBenchFollowOnHookProcess`
  pins the exact refusal string, about 25 refusal cases, and the allow list at
  `:106-114`.
- `internal/systemtest/bench_follow_on_test.go:121` `TestBenchFollowOnHookDegradedRim`
  pins warn-and-allow on degraded input.
- `internal/shellcommand/shellcommand_test.go:8,24,35,42` pin token classes, quote
  folding, redirection projection, and prefix resolution.

### Classification table

Each verdict was confirmed by a run of `dist/bench guard-bench-follow-on` on a
hand-built envelope in the shape of `bench_follow_on_test.go:33`.

| Case | Command | Verdict | Deciding line |
|---|---|---|---|
| i | `bench worktree exec L -- cp a b; cp b a` | refuse | `benchguard.go:65` via `:150-157`; bare `;` is a control operator |
| ii | `bench worktree exec L -- sh -c 'a; b'` | allow | `shellcommand.go:298-308` folds `'a; b'` to one word |
| iii | `cat <<'EOF'` with two non-bench body lines | allow | `benchguard.go:76`; no head is Bench after `shellcommand.go:355-383` strips the body |
| iv | `bench worktree exec L -- cat <<'EOF'` with body `bench gate` | refuse | `shellcommand.go:70-71` keeps `<<` as a redirection |
| v | `bench gate 2>&1` | refuse | `shellcommand.go:50` matches `>&` |
| vi | `x=1; bench gate` | refuse | `;` splits; the prefix-stepped head is Bench (`benchguard.go:63-65`) |
| vii | `cat roadmap/FT254.md; echo ---; bench maps` | refuse | `benchguard.go:150-157` scans the whole stream |
| viii | `bench worktree exec L -- rg -N '' README.md` | allow | no redirection or control operator exists |
| ix | `cp a b && bench worktree exec L -- true` | refuse | `&&` is a control operator; the later Bench head triggers `benchguard.go:65` |

### Code versus claim

1. The guard names the token that caused a refusal: does not hold. The message is
   a constant and `Classify` returns a bool.
2. A heredoc body with no `bench` word is not a follow-on: holds today (case iii,
   pinned by `benchguard_test.go:23`). Note case iv. The `<<` operator itself refuses
   when the call is Bench-headed, whatever the body holds.
3. A non-bench follow-on under `bench worktree exec` is accepted: does not hold.
   Only a quoted follow-on inside a wrapper string passes (case ii).

### Bench CLI improvements from the delegation

- `bench guard-bench-follow-on` cannot be probed directly. A `bench ... < file` or
  `printf ... | bench ...` call trips the guard itself. An `--explain '<command text>'`
  argv form that prints the verdict and the offending token would make the guard
  self-inspectable.
- The verb answers only through an exit code with the message on stderr. A one-line
  verdict on stdout would let a probe read the classification.

## #3 Gate and build forms

### Findings

- The usage trailer. `internal/usage/worktree.go:38` appends
  `bash bin/bench.sh gate --fresh` to the worktree usage string. It is not a member
  of `worktreeCommands` (`:23-35`). The same line is a public help row on `gate` at
  `cmd/bench/main.go:149`.
- The trailer's reason. The reason is the wrapper-owned run binary.
  `.bench/gate.sh:17-18` refuses a wrapper-less entry and names
  `bash bin/bench.sh gate`. `internal/conformance/gate_entry_test.go:78-81` states the
  intent. No ADR or map explains the worktree placement. The only other hit is a
  measurement note at `decisions/gate-critical-path.md:99-100`.
- Wrapper resolution, part 1. `command -v bench` is `~/.local/bin/bench`, a
  generated shim whose target is the main checkout's `bin/bench.sh`. `kit_dir()`
  (`bin/bench.sh:111-125`) returns the CWD's worktree top level when its common dir
  matches. So `BENCH_KIT` becomes the worktree.
- Wrapper resolution, part 2. `bench_binary_path` (`bin/bench.sh:254-275`) tries
  `node_modules`, the `$BENCH_HOME/cache/bin` version, then `<kit>/dist/bench`.
  Each candidate is re-anchored at the main tree through `main_tree_kit`
  (`:259-260`, `:227-237`). The comment at `:225-227` states that a linked worktree
  carries no `dist/`.
- Wrapper resolution, part 3. `BENCH_RUN_BINARY`, when present, pre-empts all of
  that (`:324-339`). `BENCH_WRAPPER` is not a routing input to the wrapper. It is
  forwarded (`:339,343,352`) and refused on the landing route (`:398-401`).
- Consequence. A bare `bench` in a worktree child runs the main checkout's
  `dist/bench`. The gate is different. `internal/gate/run_transaction.go:28-35`
  picks `runbinary.Own` when `BENCH_WRAPPER` is non-empty and `BENCH_RUN_BINARY` is
  absent. `internal/gate/gate.go:322-329` builds it from the graded tree.
- Consequence for the gate. `bench gate` under exec grades the worktree. Only new
  grammar as a bare `bench` gets the main build.
- Build script. `scripts/go-build.sh` builds `./cmd/bench` (`:124`, `:133`) into a
  caller-named output by atomic rename (`:103`, `:131`). Its grammar is
  `go-build.sh [--mode artifact] <module-root> <output-path>` (`:16`). Its inputs
  are files (`scripts/go-build.inputs`), `CGO_ENABLED=0` (`:71`), and `node`
  (`:64`, `:74-77`). It reads no `BENCH_HOME`.
  `internal/runbinary/runbinary.go:253-261` refuses a build environment without an
  absolute home (`runbinary_test.go:222`).
- Gate environment. `.bench/gate.sh:13-14,33` reads and forwards `BENCH_KIT` and
  `BENCH_RUN_BINARY`. `internal/gate/phases.go:145-149` merges both into every
  phase. `phases.go:128-131` adds `BENCH_CONFORMANCE_ROOT` and
  `BENCH_CONFORMANCE_TIER=dev` for the conformance package only.
  `run_transaction.go:144,147` sets `BENCH_KIT` and `BENCH_GATE_BASELINE`.
- Declared environment for a hand run. Ordinary packages need `BENCH_KIT` and
  `BENCH_RUN_BINARY` (`internal/testreport/command.go:188-193`) plus `GOCACHE`
  (`internal/gocache/gocache.go:13,31-33`). Conformance needs
  `BENCH_CONFORMANCE_ROOT` and `BENCH_CONFORMANCE_TIER`
  (`internal/conformance/registry/registry.go:42`). `BENCH_CONFORMANCE_CHECKS` exists
  at `:51`; no singular form exists. The system suite reads `BENCH_RUN_BINARY` and
  `BENCH_KIT` in `internal/systemtest/owner_test.go:70-105`.
- Argv and closure. The test argv comes from `gate.BaseTestArgv`
  (`internal/gate/phases.go:192-204`). `GOFLAGS` appears nowhere in the tree.
  `.bench/gate-inputs.json` declares the gate closure (`BENCH_HOME`, `HOME`, and
  eight tools). `.bench/BENCH-reference.md` has no hand-running section.
  `projects/benchkit.md:315-323` holds the phase table.

### Seam table

| Verb | Grammar | Worktree target today |
|---|---|---|
| `bench gate [--fresh]` | `bench gate [--fresh\|pin]` | CWD only; `bench worktree exec <label> -- bench gate` grades the worktree |
| `bench test` | `bench test [--full] [--package <expr> \| --changed] [--base ...] [--run <re>] \| --check <name>` | CWD only; same exec route; owns `BENCH_RUN_BINARY`, `BENCH_KIT`, `GOCACHE`, conformance pair (`command.go:38,113,142,176-193`); no build-tag form (`:169-171`) |
| `bench worktree exec` | `bench worktree exec <target> -- <command> [args...]` | the targeting seam |
| `gate-run`, `gate-phases`, `gate-go` | plumbing (`main.go:151-153`) | not agent-facing; `gate-phases` needs an inherited `BENCH_RUN_BINARY` (`phases.go:244`) |
| `bench build` | does not exist | `runbinary.Build` (`runbinary.go:196-198`) is exported with no CLI surface |

### Tests that pin the forms

- `cmd/bench/main_test.go:136` `TestHelpInventoryIsComplete` pins the help line text.
- `internal/conformance/gate_entry_test.go:83,111` pin the wrapper-less refusal wording.
- `internal/worktree/exec_test.go:97,109,119` pin the child's wrapper and strips.
- `internal/runbinary/runbinary_test.go:18,131` pin owned and inherited selection.
- `internal/gate/prospective_owner_test.go:45` pins that the gate builds from the
  graded tree.
- `internal/testreport/runbinary_test.go:14,44` and `selection_test.go:355` pin the
  selection and the environment scrub of `bench test`.
- No test pins the full `WorktreeUsage()` string. `parse_test.go:160` pins only the
  exec grammar line.

### Code versus claim

1. An exec-native gate form exists: holds in effect, not as a verb.
2. A declared environment form for focused reproductions: partly holds. No
   system-suite form exists.
3. A worktree's own grammar through `exec <label> -- ./dist/bench`: does not hold.
   No producer exists.
4. The wrapper on PATH serves the main checkout's build: holds, by design
   (`bin/bench.sh:225-227`).

### Bench CLI improvements from the delegation

- `bench test --check` with no value exits 2 with the usage line and never names the
  valid check names (`internal/testreport/command.go:92-97`).
- The follow-on hook refused `bench gate --help; bench test --check` without naming
  the segment. This reproduces the FT254 occurrence.

## #4 Conflict, create, preflight, and provenance

### Findings

- Merge conflict path. `landing.Owner.Compose` runs `git merge-tree --write-tree -z`
  in the root repository (`internal/landing/composition.go:248-250`). A conflict
  returns `ConflictError` (`internal/landing/merge.go:87-89`, message at
  `composition.go:37`). `mergeWith` renders `refused{detail=composition conflict: <kind>}`,
  `paths_total`, and a `refusal_paths` table (`internal/worktree/merge.go:86-91`,
  `classifier.go:29-64`), exit 1. No `MERGE_HEAD` exists afterwards, and HEAD and the
  tree are unchanged. `refusal.next` (`classifier.go:23-25`) stays empty.
- Landing conflict. The landing conflict names the hand route through
  `landingConflictNext` (`internal/worktree/land_refusal.go:61-81`). Tests:
  `internal/worktree/merge_test.go:335-351` (WM9, asserts an empty `git status`),
  `:544-555`, `:573`, `land_surface_test.go:241,299`, and
  `internal/landing/merge_test.go:158-163`.
- Git in-progress state. From a linked worktree, `git rev-parse --git-path` resolves
  `MERGE_HEAD`, `CHERRY_PICK_HEAD`, `REVERT_HEAD`, `rebase-merge`, and `sequencer`
  under `.git/worktrees/<id>/`. The tree reads none of them. `MERGE_HEAD` appears
  only in `roadmap/FT263.md` and `internal/landing/composition_test.go:365`.
- `--git-path` today serves `index`, `hooks`, and the lease. The call sites are
  `internal/worktree/subshell.go:188`, `clean.go:207`, `worktree.go:33`,
  `internal/diff/snapshot.go:272`, `internal/adopt/link_hook.go:201`, and
  `internal/adopt/prepush.sh:10`. No `--diff-filter=U` or `ls-files -u` exists.
- Git docs. `git-merge(1)` "HOW TO RESOLVE CONFLICTS" says `git add`, then
  `git commit` or `git merge --continue`. `--continue` runs only after a conflicted
  merge. `git-cherry-pick(1)` and `git-revert(1)` `--continue` use `.git/sequencer`.
  `git-diff(1)` documents `--diff-filter=U`. `git-ls-files(1)` `-u` prints one line
  per stage.
- Create. The grammar is at `internal/usage/worktree.go:13` and
  `internal/worktree/worktree.go:608-616`. `CreateCommand` (`:618-643`) passes the
  `startRef` of `refreshop.Consume` into `createAt` as `requestedStart`
  (`internal/worktree/ownership.go:139`). Resolution at `:172-181` takes the default
  branch tip, with `HEAD` as fallback.
- Create, the record. `--refresh` is the only non-default source today
  (`internal/worktree/refresh/refresh.go:66-84`). The branch ref comes from two
  random ids (`ownership.go:192-193`). The ledger records `Start` only (`:195-199`).
- Sibling tip. `mergeSiblingTip` (`internal/worktree/merge.go:223-257`) resolves a
  target to a verified tip. It takes a `target intent.Assignment` for the
  self-reference check. `mergeIncoming` (`:180-205`) adds a default-branch ancestry
  lookup that is a merge concern.
- Preflight. The grammar is `review <slug>` and `build <slug>` only.
  `baseCurrentCheck` (`internal/preflight/decision.go:174-188`) reds with
  `default branch tip is not an ancestor of HEAD` (`:185`).
- Preflight, the check shape. `CheckResult` is `Check, Verdict, Detail` (`:83-84`).
  `red()` sets `Detail` only (`:141-142`). The table is
  `checks{check,verdict,detail}` (`internal/preflight/command.go:68-70`). FT162
  (`roadmap/FT162.md:18-21,47,57`) is the sibling row.
- Retired worktree-merge spec. `049d12b0` retired `specs/worktree-merge/spec.md`.
  Its Out of scope (lines 398-406) reads verbatim, one bullet per cut:
  - "A `next=` remedy on the bare `bench preflight` `base-current` red that names
    the merge verb — 2 edits, 1 gate run."
  - "A `bench worktree create --from <target>` that starts a new assignment at a
    sibling's tip in one step — 3 edits, 2 gate runs."
  - "A fold of an uncommitted sibling snapshot through the merge verb — 4 edits, 2
    gate runs."
  - "An interactive conflict-resolution mode — 6 edits, 3 gate runs."
  - "A tagged system journey for the merge verb — 2 edits, 2 gate runs."
- Its "Won't handle" (`:358-361`) reads verbatim: "Conflict resolution inside the
  verb — the refusal names the paths. The hand resolution is raw Git in the
  worktree, which the platform rules leave outside the command boundary." WM27
  (`:309`) pins the landing `next=`.
- Retired census spec. `bf4b1f9e` retired `specs/exec-census/spec.md`. Its Out of
  scope (`:309`) defers the comfort half. It reads: "the exec-comfort verbs from the
  same map (stdin help line, `bench worktree resolve`, `bench worktree show`, path
  on refusal): 14 edits, 3 gate runs. It is the map's second spec." That spec was
  never written.
- Retired census map, #10 stdin. `specs/exec-census/decisions/exec-census.md` at
  that commit records: "the gap is knowledge, not code". The remedy is one help line
  and one acceptance row.
- Retired census map, #11 `resolve`, first half. It records: "stages each named
  path inside the worktree after the caller edits it. It refuses a named path that
  still holds a conflict marker and prints that path. After it stages, it prints the
  remaining unmerged paths."
- Retired census map, #11 `resolve`, second half. It continues: "When no unmerged
  path remains, it continues the in-progress operation (cherry-pick, merge, or
  revert) with the message unchanged and prints the resulting commit. The verb never
  starts the operation; a `bench worktree pick` verb is out of scope. Flagged for
  veto: the continue step makes the verb own a Git state Bench did not create."
- Retired census map, #12 `show`. It records: "prints the blob bytes to stdout
  unchanged and returns Git's exit code. It accepts the one `<rev>:<path>` operand
  form and refuses any other operand with the grammar line."
- Retired census map, #13 error paths. It records: "Every exec failure after the
  target resolves prints one line `worktree: <absolute path>` on stderr before the
  error. A child start failure and a nonzero child exit both print it. A refusal
  before the target resolves has no path; it prints `bench worktree list` as the
  next action instead." No `--env` decision exists.
- `bench worktree list`. `listTree` (`internal/worktree/list.go:126-131`) returns
  `missing` on any stat error. `orphanPath` is set only on foreign rows (`:62-67`).
  `actionsForRows` (`:99-120`) returns for an `active` row before the orphan branch
  (`:115-117`). So a `tree missing` assignment row prints `path` and `exec` actions
  that cannot succeed, and no recovery.
- `bench worktree list`, the offered recovery. `bench worktree clean --landed`
  (`:89-91`) fires on the global `landed` flag, not on this row.
  `bench worktree reclaim` is a pool-key verb (`internal/worktree/pool_reclaim.go`)
  that `list` never names.

### Code versus claim

1. A caller-edited conflict has no Bench form: holds. `land_refusal.go:80` routes to
   raw `git -C <path> merge`.
2. `create --from` and the preflight `next=` were declined cuts: holds, verbatim.
3. The preflight `base-current` red names no remedy: holds. No field exists.

Cost note: `create --from` composes `requestedStart` and `mergeSiblingTip`. The
preflight `next=` widens `CheckResult` and the rendered columns.

### Bench CLI improvements from the delegation

- `bench worktree list` names no recovery for an assignment row with `tree=missing`.
  It prints two actions that fail against it.
- The `bench worktree merge` conflict refusal leaves `next` empty. The landing
  refusal composes the exact remedy for the same conflict.
