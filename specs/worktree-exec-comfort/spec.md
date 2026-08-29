# Worktree exec comfort

Status: staged

Decision source: ready compiled map `specs/worktree-exec-comfort/decisions/worktree-exec-comfort.md`, status `ready` with thirteen tickets resolved 2026-08-29. The late heredoc question closed 2026-08-29 in this session; see Further notes.

Verification log: 1 iteration(s) to accept — round 1 (opus / medium) returned 5 blocking and 8 non-blocking findings; all 13 folded, see Further notes.

## Problem

An agent that works in an assignment worktree still reaches for a raw command on
the worktree path. `bench worktree exec` refuses a legitimate empty child argument.
No refusal prints the worktree path or names the next verb. An assignment whose tree
is gone still advertises `exec` in `bench worktree list`. Exec then refuses on a
component the operator cannot act on.

The follow-on guard refuses a plain `cp` after an exec child. It also refuses the
heredoc that would feed the child's stdin, and its refusal names no token. No verb
reads a blob at a worktree revision.

## Solution

Slice 1 of the FT254 map makes `bench worktree exec` the comfortable path for
multi-step work. Every argument after `--` reaches the child. The exec help shows the
stdin heredoc form and states the exit-2 rule. `--env KEY=VALUE` sets a child-only
variable. Every failure after target resolution prints `worktree: <absolute path>`.
Every refusal before it prints `next=bench worktree list`.

A missing tree is refused by name with its recovery verb, and `list` stops
advertising actions that cannot succeed on it. The follow-on guard accepts a
non-bench step after an exec child and a heredoc into it. It refuses everything else
as before. Its refusal names the Bench segment and the operator that caused it.
`bench worktree show <target> <rev>:<path>` prints blob bytes unchanged.

## User stories

### Group A — arguments, stdin, and environment

Line: opus / low. The grammar seam exists, the parser tests are exact, and the
ordinary test phase covers the package.

1. As an agent, I want every argument after `--` to reach the child unchanged, so that `''` runs the child instead of a grammar refusal.
2. As an agent, I want every other grammar to keep the empty-positional refusal, so that `bench commit -m x -- ''` still refuses.
3. As an agent, I want `bench worktree exec --help` to show a heredoc stdin example, so that I know stdin reaches the child.
4. As an agent, I want a heredoc on stdin to reach the child byte for byte, so that a script runs as `-- python3 -`.
5. As an agent, I want the child's stderr and exit code passed through unchanged, so that a child's usage error reads as the child's.
6. As an agent, I want the exec help to state the exit-2 prefix rule, so that both exit-2 cases read apart.
7. As an agent, I want `--env KEY=VALUE` after the target to set `KEY` in the child only, so that a worktree-local input needs no caller variable.
8. As an agent, I want `--env` to repeat, so that two inputs travel in one call.
9. As an agent, I want a malformed `--env` value refused at exit 2 with the value named, so that a typo never reaches the child.
10. As a reviewer, I want the verb's own routing variables to win over `--env`, so that `--env BENCH_HOME=/x` cannot repoint the child's pool.
11. As an agent, I want `--env` after `--` to belong to the child, so that a child that takes `--env` still receives it.

### Group B — failure surfaces and the missing tree

Line: opus / low. Each surface is one predicate at the command seam, and the
refusal fixtures exist.

12. As an agent, I want `worktree: <absolute path>` on stderr when the child cannot start, so that the raw fallback is never cheaper than the verb.
13. As an agent, I want the same line after a nonzero child exit, so that a failing child leaves me the path.
14. As an agent, I want the same line when a signal cancels the child, so that an interrupted run leaves me the path.
15. As an agent, I want a successful child to print no `worktree:` line, so that exec adds nothing to a green run.
16. As an agent, I want every pre-resolution refusal except a missing tree to print `next=bench worktree list`, so that an unassigned target names the lookup.
17. As an agent, I want `bench worktree path` to print the same `next=` line, so that both target verbs describe one failure one way.
18. As an agent, I want a missing tree refused with `worktree tree is missing`, so that the refusal names a fact I can act on.
19. As an agent, I want that refusal's `next=` to name the recovery verb, so that the retained record has one route out.
20. As an operator, I want no `path` or `exec` action on a `tree missing` assignment row, so that no advertised action fails.
21. As an operator, I want exactly one help row that names that row's recovery verb, so that the record has the route exec names.
22. As an operator, I want a present-tree active row to keep its `path` and `exec` actions, so that the change touches only missing trees.

### Group C — the follow-on guard

Line: opus / medium. The classifier changes from a stream-wide scan to a
span-scoped rule, and the seam has more than one shape.

23. As an agent, I want a non-bench step after an exec child allowed, so that a cp-aside probe costs one call.
24. As an agent, I want a non-bench segment before a Bench call still refused, so that a prefix never shapes a Bench call.
25. As an agent, I want a pipe, `||`, `&`, or a plain redirection on the exec segment refused, so that the output stays complete.
26. As an agent, I want a second Bench segment after an exec child refused, so that two Bench responses never share one call.
27. As an agent, I want a heredoc on the exec segment allowed, so that the stdin example in the help is usable.
28. As an agent, I want a heredoc on any other Bench call still refused, so that FOG38 holds.
29. As an agent, I want a heredoc body with no `bench` word to stay allowed, so that a plain `cat <<'EOF'` never trips the guard.
30. As an agent, I want the refusal to name the Bench segment's words and the adjacent operator, so that I can fix the call.
31. As an agent, I want the adjacent operator chosen by one fixed precedence, so that the named token is the cause.
32. As a reviewer, I want the fixed refusal sentence unchanged, so that the hook's pin and the system suite stay green.
33. As a reviewer, I want the census to keep skipping Bench-headed calls, so that an exec follow-on never counts as a raw call.

### Group D — `bench worktree show`

Line: opus / low. The verb mirrors `path`, and the target fixture exists.

34. As an agent, I want `bench worktree show <target> <rev>:<path>` to print the blob bytes unchanged, so that I read a revision without the path.
35. As an agent, I want a NUL byte in the blob to arrive unchanged, so that a binary blob reads as Git stores it.
36. As an agent, I want Git's exit code and stderr when the object is missing, so that a bad revision reads as Git's own error.
37. As an agent, I want an operand with no `:` refused at exit 2, so that a bare path never reads as a revision.
38. As an agent, I want an operand that starts with `-` refused at exit 2, so that no Git option reaches the child.
39. As an agent, I want `show` to share the target refusal printer, so that an unassigned target names `bench worktree list`.
40. As a cold agent, I want `show` listed in `bench worktree --help` and `bench help`, so that discovery needs no guesswork.
41. As a cold agent, I want the help inventories to show the `--env` grammar, so that help and parser agree.

### Group E — guidance

Line: opus / medium. The group edits shared platform prose and the glossary. The
leverage rule routes those to the mid tier. The reviewer's 2026-08-26 direction caps
every subagent at medium.

42. As a reader, I want the operating guide's contract sentence to name the exec exception, so that the rule and the guard agree.
43. As a reader, I want the `CONTEXT.md` `follow-on` term to state the exec exception, so that the glossary matches the guard.

## Implementation decisions

- The exec grammar declares that every token after `--` is child argv. The parser
  gains one grammar attribute for that declaration. A grammar without it keeps the
  empty-positional refusal after `--`.
- `--env` is a repeatable value flag. The parser gains one flag attribute for
  repetition. A repeated non-repeatable flag stays a usage error.
- `--env` sits after the reserved target and before `--`. The grammar is
  `bench worktree exec <target> [--env KEY=VALUE]... -- <command> [args...]`.
- `KEY` matches `[A-Za-z_][A-Za-z0-9_]*`. A value with no `=`, an empty `KEY`, or a
  bad `KEY` is a usage error that names the value. An empty `VALUE` is allowed.
- The child environment applies `--env` after the caller's environment and before
  the verb's own routing variables. So `BENCH_HOME`, `BENCH_WRAPPER`, `BENCH_KIT`,
  and `BENCH_RUN_BINARY` keep the verb's values.
- The exec grammar's help text carries three lines: the grammar line, the stdin
  heredoc example, and the exit-2 rule. `bench worktree --help` and `bench help`
  keep the one-line grammar.
- `runWorktreeChild` writes `worktree: <absolute path>` to stderr once. It does so
  after a nonzero child exit, after a start failure's error line, and on the cancel
  path. A zero exit prints nothing.
- The shared target refusal printer prints its detail line and then one line
  `next=<verb>`. It takes the `next` per reason. Every pre-resolution reason names
  `bench worktree list`, except `worktree tree is missing`.
- One recovery producer serves the missing-tree `next` for exec and for `list`. It
  names `bench worktree clean --landed` for a landed assignment. It names
  `bench worktree release --request <request> <abs path>` otherwise, with the path
  shell-quoted as `axi.ShellQuote` renders it.
- A landed missing-tree row emits the producer's `bench worktree clean --landed`
  action. The help renderer collapses it with the global landed action, so
  `list` prints that verb once. A not-landed missing-tree row adds exactly one
  action, the release form. The landed fact is the branch's landing in the default
  branch, because a missing tree proves nothing about the assignment.
- The target resolver checks the tree after the active-state check and before the
  creation bundle. A not-exist stat error is `worktree tree is missing`. Any other
  stat error keeps today's path.
- `list` builds an assignment row's actions from the tree cell and the landed cell.
  A `missing` tree yields no `path` or `exec` action.
- The guard classifies per Bench segment. The exec exception applies when the first
  Bench-headed simple command is `bench worktree exec`. Under it, a heredoc
  redirection inside the span is allowed. A `;` or `&&` after the span is allowed
  when every later simple command is non-Bench. Everything else refuses.
- Any other Bench-headed command keeps the stream-wide rule. `Classify` returns the
  verdict, the Bench segment's projected words, and the adjacent operator token.
- The refusal message is the fixed sentence, one space, then
  `segment=<words> operator=<token>`. `<words>` are the segment's projected words
  joined by single spaces.
- `<token>` follows one precedence: a redirection inside the span when one exists,
  else the control operator before the span, else the one after. A redirection
  prints as its source text: the fd digits, the operator, and the operand joined
  without spaces, so `2>&1` prints as `2>&1`.
- The exec allowance covers a heredoc (`<<` and `<<-`) only. A here-string `<<<`
  is a plain redirection and refuses.
- `bench worktree show <target> <rev>:<path>` resolves the target through the shared
  resolver. It runs `git cat-file blob <rev>:<path>` in the worktree with stdout,
  stderr, and the exit code passed through.
- The `show` operand must hold a `:` and must not start with `-`. `show` joins the
  worktree command list, the `bench help` inventory, and the kept routes.

## Testing decisions

- A good test drives the public command function with an argv and a real child or a
  real Git repository. It reads stdout, stderr, and the exit code. It never asserts
  an internal call.
- Rows X4, F3, F10, F11, F12, and F13 drive `runWorktreeChild` and `actionsForRows`
  directly, by the precedent of `TestActionsForRowsEnumeratesActiveAndOrphanRows`.
  The public seam cannot inject a reader or a row.
- The grammar rows attach to `internal/usage/parse_test.go`.
- The stdin, exit, environment, and `worktree:` rows attach to
  `internal/worktree/exec_test.go`, which already runs a real `sh` child through
  `runWorktreeChild`.
- The refusal-line rows attach to `internal/worktree/identifier_operand_test.go`
  beside `TestTargetVerbsNameTheResolverReason` and
  `TestTargetVerbsShareOneRefusalPrinter`.
- The action rows attach to `internal/worktree/list_actions_test.go` beside
  `TestActionsForRowsEnumeratesActiveAndOrphanRows`.
- The classifier rows attach to `internal/benchguard/benchguard_test.go` beside
  `TestClassifyFollowOns`. The hook-process rows attach to
  `internal/systemtest/bench_follow_on_test.go` inside `TestBenchFollowOnHookProcess`.
- The `show` rows attach to a new `internal/worktree/show_test.go` beside the `path`
  tests.
- The inventory rows attach to `cmd/bench/main_test.go`
  (`TestHelpInventoryIsComplete`) and `cmd/bench/command_registry_test.go`
  (`TestKeptWorktreeOperationsKeepTheirGrammar`).
- The gate observes the feature in the `test` phase for every package above, and in
  the `system` phase for the hook-process rows.

### Seam diagram

    trigger: an agent's Bash call, or a delegate's charge
        │
        ▼
    argv  ──▶  [ usage.Parse: exec grammar ]  ──▶  Result{target, env, child argv}
                  ◀ tests attach here: parse_test drives Parse with an argv table
        │
        ▼
    target  ──▶  [ resolveWorktree + printTargetRefusal ]  ──▶  path, or refusal + next=
                  ◀ tests attach here: identifier_operand_test reads stderr lines
        │
        ▼
    child argv, stdin, env  ──▶  [ runWorktreeChild ]  ──▶  child stdout, stderr, exit, worktree: line
                  ◀ tests attach here: exec_test runs a real sh child
        │
        ▼
    hook JSON  ──▶  [ benchguard.Classify ]  ──▶  allow, or BLOCKED sentence + segment= operator=
                  ◀ tests attach here: benchguard_test table; systemtest runs the hook script

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| X1 | 1 | `bench worktree exec <target> -- rg -N '' README.md` parses with child argv `rg`, `-N`, `''`, `README.md` and exit 0 from the parser | `internal/usage/parse_test.go` | the empty-positional guard past `--` refuses the `''` today |
| X2 | 2 | `bench commit -m x -- ''` still returns the usage line naming `""` at exit 2 | `internal/usage/parse_test.go` | a parser that drops the guard for every grammar lets an unset variable widen a path verb |
| X3 | 3 | `bench worktree exec --help` prints the grammar line, a line that shows `<<'EOF'` feeding `-- python3 -`, and the exit-2 rule line | `internal/worktree/exec_test.go` | a one-line help never shows stdin |
| X4 | 4 | a `runWorktreeChild` call with an `io.Reader` of three lines ending in a NUL byte makes `-- cat` emit those bytes unchanged on stdout | `internal/worktree/exec_test.go` | no test drives stdin today, so a dropped `cmd.Stdin` stays green |
| X5 | 5 | a child `sh -c 'echo child-usage >&2; exit 2'` yields exit 2 and stderr equal to `child-usage` plus the `worktree:` line | `internal/worktree/exec_test.go` | a wrapper that rewrites the child's stderr or maps its code breaks the pass-through |
| X6 | 6 | the exec help's third line names `usage: bench worktree exec` as the prefix of an exit-2 grammar refusal | `internal/worktree/exec_test.go` | a help without the rule leaves the two exit-2 cases indistinguishable |
| X7 | 7 | `bench worktree exec <target> --env FOO=bar -- sh -c 'echo $FOO'` prints `bar`, and the caller's process has no `FOO` | `internal/worktree/exec_test.go` | a flag the parser accepts but the child never receives |
| X8 | 8 | `--env A=1 --env B=2` sets both `A` and `B` in the child | `internal/usage/parse_test.go`, `internal/worktree/exec_test.go` | the repeated-flag usage error refuses the second `--env` today |
| X9 | 8 | `bench commit -m a -m b` still returns the usage line naming `-m` at exit 2 | `internal/usage/parse_test.go` | a parser that makes every flag repeatable hides a mistyped invocation |
| X10 | 9 | `--env FOO` and `--env 1X=y` each return the usage line naming the value at exit 2 and start no child | `internal/worktree/exec_test.go` | a value with no `=` would otherwise reach `execEnv` as a malformed entry |
| X11 | 10 | `--env BENCH_HOME=/x -- sh -c 'echo $BENCH_HOME'` prints the verb's resolved home | `internal/worktree/exec_test.go` | an `--env` applied last repoints the child's pool |
| X12 | 11 | `bench worktree exec <target> -- printf '%s' --env` prints `--env` | `internal/usage/parse_test.go` | a flag pass that runs past `--` eats the child's argument |
| X13 | 10 | `--env BENCH_KIT=/x -- sh -c 'echo ${BENCH_KIT-unset}'` prints `unset` | `internal/worktree/exec_test.go` | an `--env` applied after the routing strip survives it |
| F1 | 12 | `-- no-such-binary-xyz` exits 1 with stderr `bench worktree exec: exec: ...` followed by `worktree: <absolute path>` | `internal/worktree/exec_test.go` | the start failure prints only the os/exec error today |
| F2 | 13 | `-- sh -c 'exit 3'` exits 3 and stderr ends with the `worktree: <absolute path>` line | `internal/worktree/exec_test.go` | a nonzero child exit prints no path today |
| F3 | 14 | a child killed by SIGINT through the cancel path exits 130 and stderr ends with the `worktree:` line | `internal/worktree/exec_test.go` | the cancel path prints nothing today |
| F4 | 15 | `-- true` exits 0 with empty stderr | `internal/worktree/exec_test.go` | a line printed on every exit adds noise to a green run |
| F5 | 16 | `bench worktree exec no-such-label -- true` prints `bench worktree exec: target is unassigned` then `next=bench worktree list` on stderr at exit 1 | `internal/worktree/identifier_operand_test.go` (`TestTargetVerbsNameTheResolverReason`) | the refusal ends at the reason today |
| F6 | 17 | `bench worktree path no-such-label` prints the same two lines with the `bench worktree path` prefix | `internal/worktree/identifier_operand_test.go` (`TestTargetVerbsShareOneRefusalPrinter`) | a second printer describes one failure two ways |
| F7 | 18 | an active assignment whose worktree directory is removed refuses exec with `bench worktree exec: worktree tree is missing` at exit 1 | `internal/worktree/identifier_operand_test.go` (`TestTargetVerbsNameTheResolverReason`) | the creation bundle names the owner marker today |
| F8 | 19 | that refusal's second line is `next=bench worktree clean --landed` when the assignment branch has landed | `internal/worktree/identifier_operand_test.go` (`TestTargetVerbsNameTheResolverReason`) | a refusal with no route leaves the record immortal |
| F9 | 19 | that refusal's second line is `next=bench worktree release --request <request> <abs path>`, shell-quoted as `axi.ShellQuote` renders it, when the branch has not landed | `internal/worktree/identifier_operand_test.go` (`TestTargetVerbsNameTheResolverReason`) | a producer that only knows the landed case prints a wrong verb |
| F10 | 20 | `actionsForRows` yields no `path` and no `exec` action for an active assignment row whose tree cell is `missing` | `internal/worktree/list_actions_test.go` | the active branch returns before it reads the tree cell today |
| F11 | 21 | a not-landed active row whose tree cell is `missing` yields exactly one action, `bench worktree release --request <request> <abs path>`, shell-quoted as `axi.ShellQuote` renders it | `internal/worktree/list_actions_test.go` | a list that drops the actions and adds nothing loses the route |
| F12 | 22 | an active row whose tree cell is `present` keeps its `path` and `exec` actions unchanged | `internal/worktree/list_actions_test.go` (`TestActionsForRowsEnumeratesActiveAndOrphanRows`) | a change that keys on state alone strips every active row |
| F13 | 21 | `bench worktree list` with one landed active assignment whose tree is missing prints exactly one `bench worktree clean --landed` help row | `internal/worktree/list_actions_test.go` | a per-row action beside the global landed action prints the verb twice |
| G1 | 23 | `bench worktree exec L -- cp a b; cp b a` and `bench worktree exec L -- cp a b && rg -n x b` are allowed | `internal/benchguard/benchguard_test.go` | the stream-wide scan refuses the `;` today |
| G2 | 24 | `cp a b && bench worktree exec L -- true` is refused | `internal/benchguard/benchguard_test.go` | an exec exception applied to both sides lets a prefix through |
| G3 | 25 | `bench worktree exec L -- true \| cat`, `... \|\| echo x`, `... &`, `... > out`, and `... <<< x` are each refused | `internal/benchguard/benchguard_test.go` | an exception keyed on the exec head alone lets a pipe through |
| G4 | 26 | `bench worktree exec L -- true; bench maps` is refused | `internal/benchguard/benchguard_test.go` | a follow-on rule that never looks at the later heads lets two Bench calls share a response |
| G5 | 27 | `bench worktree exec L -- cat <<'EOF'` with a body line `bench gate` is allowed | `internal/benchguard/benchguard_test.go` | the `<<` redirection refuses today (case iv) |
| G6 | 28 | `bench gate <<'EOF'` with a body line `input` is refused | `internal/systemtest/bench_follow_on_test.go` (FOG38) | a heredoc allowance on every verb breaks FOG38 |
| G7 | 29 | `cat <<'EOF'` with two non-bench body lines is allowed | `internal/benchguard/benchguard_test.go` | a classifier that stops stripping bodies reads the body as a call |
| G8 | 30 | the refusal for `cat a && echo x; bench maps` ends with `segment=bench maps operator=;` | `internal/benchguard/benchguard_test.go` | a rule that names the first stream operator prints `&&` |
| G9 | 31 | the refusal for `bench gate 2>&1` ends with `segment=bench gate operator=2>&1` | `internal/benchguard/benchguard_test.go` | a rule that picks the first stream operator names nothing for an in-span redirection |
| G10 | 31 | the refusal for `cp a b && bench worktree exec L -- true \| cat` names `operator=&&` | `internal/benchguard/benchguard_test.go` | the pipe after the span and the prefix operator compete, and the prefix wins |
| G11 | 32 | every refusal's stderr still holds `BLOCKED: Bench response is bounded, complete, and self-contained` | `internal/systemtest/bench_follow_on_test.go` (`TestBenchFollowOnHookProcess`) | the hook's pin and the system suite read that sentence |
| G14 | 23 | `bench gate; cp a b` is refused | `internal/benchguard/benchguard_test.go` | an allowance keyed on any Bench head lets a follow-on through on every verb |
| G12 | 23, 27 | the hook process allows `bench worktree exec label -- cp a b; cp b a` and `bench worktree exec label -- cat <<'EOF'` with a body | `internal/systemtest/bench_follow_on_test.go` | a unit-level allowance the shell hook never reaches |
| G13 | 33 | `bench worktree exec L -- sed -i x <pool path>; cp a b` records no census entry | `internal/census/census_test.go` (`TestRecordSkipsABenchCall`) | a census that re-classifies the follow-on counts the child's head |
| S1 | 34 | `bench worktree show <target> HEAD:tracked.txt` prints the committed bytes of `tracked.txt` on stdout at exit 0 | `internal/worktree/show_test.go` | no verb exists |
| S2 | 35 | a blob that holds `a\x00b\n` prints exactly those four bytes | `internal/worktree/show_test.go` | a line-oriented writer drops or alters the NUL |
| S3 | 36 | `HEAD:no-such-file` exits 128 with stderr byte-equal to a direct `git cat-file blob HEAD:no-such-file` in the same worktree, and stdout is empty | `internal/worktree/show_test.go` | a verb that maps Git's error to its own hides the revision |
| S4 | 37 | the operand `tracked.txt` returns the grammar line at exit 2 and runs no Git | `internal/worktree/show_test.go` | an operand without `:` would read as a pathspec |
| S5 | 38 | the operand `--output=/tmp/x:tracked.txt` returns the grammar line at exit 2 and runs no Git | `internal/worktree/show_test.go` | a dash operand reaches Git as an option |
| S6 | 39 | `bench worktree show no-such-label HEAD:x` prints `bench worktree show: target is unassigned` then `next=bench worktree list` | `internal/worktree/identifier_operand_test.go` (`TestTargetVerbsNameTheResolverReason`) | a third printer describes one failure a third way |
| S7 | 40, 41 | `bench help` prints the `show` row and the exec row with `[--env KEY=VALUE]...`, byte-equal to the fixture | `cmd/bench/main_test.go` (`TestHelpInventoryIsComplete`) | the fixture is exact, so a missing row reds |
| S8 | 40 | `bench worktree --help` names `bench worktree show <target> <rev>:<path>` | `cmd/bench/command_registry_test.go` (`TestKeptWorktreeOperationsKeepTheirGrammar`) | the kept-routes list names each verb, so a missing `show` reds |

Not covered: story 42 — the operating guide's contract sentence is prose. The review
round grades it, and the guard rows above pin the behavior it describes.
Not covered: story 43 — the glossary entry is prose. The review round grades it
against G1 through G5.

### Edge inventory

- An empty `''` after `--`: X1. A `-` alone after `--`: child argv, by the existing
  parser rule.
- `--env FOO=` sets `FOO` to the empty string. `--env ""` is a usage error that
  names `--env ""`, by the existing `NoEmptyValue` rule.
- `--env` that names a routing variable: X11.
- A target with control bytes, an unassigned target, an ambiguous target, or an
  inactive assignment. Each keeps its existing refusal and gains the `next=` line
  (F5, F6, S6).
- A missing tree versus an unreadable tree. F7 covers not-exist. A permission error
  keeps today's creation-bundle path.
- A child killed by a signal other than SIGINT. The exit is `128+signal`, by the
  existing rule, and the `worktree:` line prints under F2's rule.
- Guard edges. A prefix operator, an in-span pipe, an in-span redirection, a
  heredoc on exec, a heredoc elsewhere, a second Bench segment, and a wrapper
  string. G1 through G10 and the existing wrapper rows cover them.
- Show edges. NUL bytes (S2), a missing object (S3), no `:` (S4), a dash operand
  (S5), and control bytes in the operand (the existing `lineSafe` refusal). A
  `HEAD:` that names a tree gets Git's own error under S3's rule.
- **Won't handle** the release form's behavior on a missing tree — the `next=`
  names it as decided, and the landed row's clean route survives.
- **Won't handle** `--env` before the target — the reserved slot takes it as the
  target and refuses it, and the help shows the order.
- **Won't handle** a heredoc on a wrapper string such as `bash -lc '...' <<EOF` —
  the wrapper's outer syntax rule still refuses, and the direct exec form survives.
- **Won't handle** a `worktree:` line on a pre-resolution refusal — no path exists
  yet, and `next=bench worktree list` survives.
- **Won't handle** `show` on a commit or tree object — `cat-file blob` refuses with
  Git's own error, and a blob path survives.
- **Won't handle** the `show` exit code under an interrupt — Git's own code for a
  killed child survives, and `exec`'s 130 rule stays exec's.

## Ownership fences

- `internal/usage/`
- `internal/worktree/exec.go`
- `internal/worktree/path.go`
- `internal/worktree/list.go`
- `internal/worktree/show.go`
- `internal/worktree/exec_test.go`
- `internal/worktree/identifier_operand_test.go`
- `internal/worktree/identity_component_test.go`
- `internal/worktree/list_actions_test.go`
- `internal/worktree/show_test.go`
- `internal/worktree/testdata/`
- `internal/benchguard/`
- `internal/shellcommand/`
- `internal/census/census_test.go`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `cmd/bench/command_registry_test.go`
- `internal/systemtest/bench_follow_on_test.go`
- `.bench/BENCH.md`
- `CONTEXT.md`
- `specs/worktree-exec-comfort/`
- `decisions/worktree-exec-comfort.md`
- `decisions/assets/ft254-exec-comfort-research.md`
- `CHANGELOG.md`

## Out of scope

- Slice 2 of the map — 12 edits, 4 gate runs. It holds `bench worktree build
  <label>`, the usage trailer replaced by the exec form, `bench test --check
  system`, `bench worktree create --from <target>`, and the preflight `next` column.
- Slice 3 of the map — 6 edits, 2 gate runs. It holds `bench worktree resolve
  <target> <path>...` after FT263 lands.
- A `bench guard-bench-follow-on --explain '<command>'` probe form — 3 edits, 1 gate
  run. It prints the verdict and the token for a literal string.
- A `worktree:` line or a `next=` on `bench worktree list` failures — 2 edits, 1
  gate run.

## Further notes

The late heredoc question closed 2026-08-29 in the spec session. The guard allows a
heredoc redirection on the `bench worktree exec` segment only. Every other
redirection on that segment still refuses, and every heredoc on another Bench call
still refuses. The map's #10 answer carries that clause.

The map compiled here holds three slices, and this spec is slice 1. Slices 2 and 3
read this folder's `decisions/` as settled provenance. FT254's roadmap row stays in
place until slice 3 retires it.

The `Roadmap:` line is absent on purpose. `bench spec retire` names the row's detail
file, and this slice does not retire FT254.

Review round 1 folded these findings:

- The missing-tree `next=` is the one pre-resolution exception (story 16, F9).
- G8 now uses a case whose two operators differ.
- F9, F11, and the new F13 pin the literal recovery strings and the single landed
  action.
- The new X13 pins the `BENCH_KIT` strip.
- The `2>&1` token rule and the `<<<` refusal are stated.
- S8 drops its `--env` half, and S3 pins Git's bytes.
- Every existing seam names its function, and the direct-call precedent is stated.
- Story 41 moved to Group D.
- The new G14 keeps the allowance exec-only.

The build recorded three decisions on 2026-08-29, open to reviewer veto. The
missing-tree ticket found that the landed plan a present tree needs cannot see a
missing tree. So the landed fact for the recovery producer is the branch's landing
in the default branch. The global landed action does not fire on a missing tree.
F13 holds because the landed missing-tree row emits the same action, and the
help renderer collapses the pair.

The new `next=` line reds one expectation in
`internal/worktree/identity_component_test.go`, so the fence gained that file. The
guard ticket's `segment=` and `operator=` half reaches the hook's stderr only
through `cmd/bench/main.go`. That file is inside the fence, so the ticket's
`Writes:` list gained it in a continuation.

The review's own question found a clean partition. The guard set shares no file,
row, or seam with the other five tickets. That set is `internal/benchguard`,
`internal/shellcommand`, `internal/systemtest`, `internal/census`, `.bench/BENCH.md`,
and `CONTEXT.md`. It could ship as its own spec. That split is the reviewer's call
at sign-off, and the guard ticket already runs on the frontier concurrently.
