# Refuse a pathless or substituting rg call

Status: staged

Roadmap: FT366

Decision source: named reviewed artifact `roadmap/FT366.md` (the drained FT366 row), with reviewer decisions D1 and D2 closed in conversation on 2026-11-02.

Verification log: 4 iteration(s) to accept — the spec round took 2 iterations (it found an exec child with no substitution mark, a mistraced R54 message, and two plan requirements with no probe), and the ticket round took 2 iterations (it found that ticket 05 did not wait for ticket 04, and that ticket 02 named no source for the Bench executable names).

## Problem

An agent runs `rg` with no path operand. When the stdin of the call is an open pipe, `rg` reads that pipe and waits until the pipe closes.
Three recorded delegate calls waited on stdin (`roadmap/FT254.md:84`, `roadmap/FT254.md:99`, `roadmap/FT366.md:18`). Two records give a duration, 48 minutes and about two hours, and the third call waited until the orchestrator stopped it.
Today the main Claude Code Bash tool and the `bench worktree exec` child both get `/dev/null` as stdin, and `rg` then searches the current directory (observed on 2026-11-02).
The hazard therefore depends on the route, and only a delegate route with an inherited pipe shows it.
The only statement of the rule is prose in `AGENTS.md:61-63`, and no hook enforces it.

An agent also writes a backtick span or a `$(...)` inside a double-quoted `rg` pattern. The shell runs that command before `rg` starts.
A backtick span in a double-quoted pattern ran `bench handoff` past every hook (`roadmap/FT366.md:20`).
The shared projection has no command-substitution fact. Quote folding erases the quote style, so `'lit$(x)'` and `"$(x)"` give the same word (`internal/shellcommand/shellcommand.go:23-26`, `:323-349`).
An escaped backtick and a live backtick give the same token (`internal/shellcommand/shellcommand.go:350-356`).
No guard reads the quote style of a word or a substitution inside it, so neither Bash guard sees the span (`internal/gitguard/scan.go:44-71`, `internal/benchguard/benchguard.go:119-149`).

The Bench follow-on guard refuses an unquoted `$(` after a Bench span with the label `operator=(` (`roadmap/FT366.md:19`).
The label names a parenthesis, not the command substitution that the agent wrote.
The label comes from the token after the span (`internal/benchguard/benchguard.go:263-265`) and from the message format (`internal/benchguard/benchguard.go:78`).

## Solution

A new PreToolUse guard reads every Bash call through the shared `shellcommand` projection and inspects each `rg` call.
It refuses an `rg` call with no path operand, unless a pipe into that call or an input redirection on that call supplies stdin.
It knows the ripgrep options that consume a value, and it permits a call when its operand analysis is uncertain.
It refuses every `rg` call that holds a live command substitution, quoted or unquoted, and its refusal names command substitution.

The lexer of the shared projection marks each live command substitution once. The new guard and the Bench follow-on guard read that one mark.
The Bench follow-on guard then names command substitution where it now prints `operator=(`. It refuses nothing new.
The existing Git and Bench guards keep their fail-closed policy for a command that the lexer refuses.

Like every guard, this guard is an honest-mistake check, not an evasion-resistant shell interpreter. Its projection never runs the command.

## User stories

Line: opus / xhigh.
Implementation-line reason: C1 changes the lexer that three guards, census, and worktree exec read, so a wrong mark moves verdicts outside this spec. C2 encodes the ripgrep operand grammar with a permit-on-uncertainty policy that only exact rows can grade. The rows are exact and the unit seams are cheap, so the remaining chunks run at medium.
Harder chunks: C1, C2.

### Shared substitution fact

1. As a guard author, I want the projection to mark a live `$(` inside double quotes, so that a guard sees a substitution that quote folding hides.
2. As a guard author, I want the projection to mark a live backtick span, unquoted or inside double quotes, so that a guard sees the span that ran `bench handoff`.
3. As a guard author, I want the projection to mark the `(` operator that opens an unquoted `$(`, so that a guard can tell a substitution from a subshell.
4. As a guard author, I want a `$(` or a backtick inside single quotes to stay unmarked, so that a literal pattern passes.
5. As a guard author, I want a backslash-escaped `$` or backtick to stay unmarked, so that an escaped literal passes.
6. As a guard author, I want a bare subshell `(` with no `$` before it to stay unmarked, so that a subshell keeps its operator meaning.
7. As a maintainer of the existing guards, I want the token text, the token kinds, the spans, and the `Unlexed` flag to stay unchanged, so that the Git and Bench guard verdicts stay the same.

### The rg guard

8. As an agent, I want an `rg` call with a pattern and no path refused, so that the call does not wait on an open stdin pipe.
9. As an agent, I want an `rg` call with a path operand permitted, so that an ordinary search runs.
10. As an agent, I want a `|` or `|&` into the `rg` call to count as its stdin, so that `cmd | rg foo` runs.
11. As an agent, I want a pipe after the `rg` call not to count as its stdin, so that `rg foo | wc -l` is refused.
12. As an agent, I want an input redirection on the `rg` call to count as its stdin, so that `rg foo < file` runs.
13. As an agent, I want each value option of ripgrep to consume its value in every spelling, so that `rg -g '*.go' foo` reads as a pattern with no path.
14. As an agent, I want a short-option bundle to end at its first value letter, so that `rg -nA3 foo .` reads `.` as a path.
15. As an agent, I want `-e` and `-f` to make every positional a path, so that `rg -e foo src` runs and `rg -e foo` is refused.
16. As an agent, I want a no-pattern mode such as `--files` or `--version` permitted, so that a mode that never reads stdin runs.
17. As an agent, I want `--` to end the options, so that `rg -- -foo .` reads `-foo` as the pattern and `.` as the path.
18. As an agent, I want an option that the table does not know to make the call uncertain, so that the guard permits it.
19. As an agent, I want a command that the lexer refuses to make every `rg` call in it uncertain, so that the guard permits it.
20. As an agent, I want `xargs rg` to make the call uncertain, so that the guard permits an operand that comes from stdin.
21. As an agent, I want a process substitution on the `rg` call to make the call uncertain, so that `rg foo <(cmd)` runs.
22. As an agent, I want an inline `RIPGREP_CONFIG_PATH=` assignment to make the call uncertain, so that a configured pattern does not cause a false refusal.
23. As an agent, I want an `rg` call with no pattern and no path permitted, so that the call that ripgrep rejects at once runs to its own error.
24. As an agent, I want every `rg` call with a live command substitution refused, even with a path, so that the shell does not run a hidden command.
25. As an agent, I want the substitution refusal to name command substitution and to advise single quotes, so that I can repair the call.
26. As an agent, I want a single-quoted or escaped `$(` or backtick in an `rg` call permitted, so that a literal search runs.
27. As an agent, I want the substitution refusal to win over the pathless refusal, so that I repair the cause that splits the call.
28. As an agent, I want the pathless refusal to name the missing path and the `.` repair, so that I can repair the call.
29. As an agent, I want the guard to find `rg` behind a routine prefix such as `env`, `timeout`, `command`, or `nohup`, so that a prefix does not hide the call.
30. As an agent, I want the guard to read one `sh -c`, `bash -c`, or `zsh -c` child, so that a wrapped call gets the same verdict.
31. As an agent, I want the guard to read the child of `bench worktree exec <target> -- rg ...`, so that the delegate route gets the same verdict.
32. As an agent, I want a pipe or an input redirection on the wrapper or exec call to feed the first `rg` of its child, so that `cmd | bench worktree exec X -- rg foo` runs.
33. As an agent, I want each simple command of a chain checked, so that `cd x && rg foo` is refused.
34. As an agent, I want `command -v rg` not to count as an `rg` call, so that a lookup runs.
35. As an agent, I want `-` as a path operand, or `-f -`, to count as an explicit stdin request, so that a call that names stdin on purpose runs.
36. As a maintainer, I want the guard to stop at one wrapper level and at the known routine prefixes, so that it keeps the depth of the existing guards.
37. As a kit maintainer, I want the new guard package under the projection-owner check, so that it cannot decode an envelope or a wrapper privately.
38. As a kit maintainer, I want one package to own the ripgrep option table, so that a later reader reuses the same grammar.
58. As an agent, I want a live substitution in the child of `bench worktree exec <target> -- rg ...` refused, so that the delegate route cannot run a hidden command.

### The hook

39. As an agent in Claude Code, I want the rg guard to run in the Bash PreToolUse group, so that a refusal reaches me before the call.
40. As an agent in Codex, I want the rg guard to run in the Bash PreToolUse group, so that both adopted harnesses refuse the same call.
41. As a guard reader, I want `bench guards` to list the rg guard with its deny text, so that I learn the surface without a collision.
42. As an agent, I want the shim to allow the call when the Bench core or the shared library is missing, so that a broken install does not block every Bash call.
43. As an agent, I want the shim to allow the call when the core does not know the verb, so that a stale core does not block every Bash call.
44. As an agent, I want the verb to allow the call when the core errors or the envelope is unreadable, so that the guard fails open.
45. As an agent, I want the hook to exit 2 with the refusal on stderr, so that the harness blocks the call and shows the repair.
46. As a kit maintainer, I want the new verb registered in the hook seams, the routing table, and the parity table, so that conformance grades the shim.
47. As a route author, I want each refusal-route step graded against the rg guard, so that no route tells an agent to run a refused call.
48. As a teammate, I want the Hook Layers section to name the rg guard, so that the reference describes the live deny surface.
49. As a teammate, I want the `AGENTS.md` rg path bullet retired, so that the rule has one source.
50. As a linked-repo user, I want the shim to ship in the consumer payload, so that a linked repository gets the guard.

### The Bench follow-on refusal

51. As an agent, I want the follow-on refusal at an unquoted `$(` after a Bench span to name command substitution, so that I see the real cause.
52. As an agent, I want the follow-on refusal at an unquoted `$(` before a Bench span to name command substitution, so that `echo $(bench status)` shows the real cause.
53. As an agent, I want a double-quoted `$(...)` in a Bench call to pass as it does now, so that the guard refuses nothing new.
54. As an agent, I want a bare subshell `(` before a Bench span to keep `operator=(`, so that a subshell refusal keeps its meaning.
55. As a maintainer, I want the guard differential golden to stay byte-equal, so that the existing guards keep their verdicts.

### Reviewed exclusions

56. As a maintainer, I want an exported `RIPGREP_CONFIG_PATH` left unread, so that the guard reads only the command text.
57. As a maintainer, I want the digit-only operand before a redirection left as the projection gives it, so that this spec does not change `ProjectCommandWords`.

## Implementation decisions

**The live-substitution mark.** The lexer `lex` returns plain strings today (`internal/shellcommand/shellcommand.go:299-365`), and `Parse` builds each `Token` from them (`:62-97`).
`Token` gains one boolean field, `Substitution` (`:23-26`). `Parse` sets the field from a mark that `lex` records beside each token.
The lexer marks a Word token when the token holds a `$(` or a backtick outside single quotes and not after a backslash.
Inside double quotes the lexer treats `\$` and a backslash before a backtick as escapes for the mark only. The folded text stays as it is now.
The lexer marks a ControlOperator token whose text starts with `(` when an unquoted, unescaped `$` comes directly before it with no space.
The mark adds information only. Token text, token kinds, spans, and `Unlexed` stay as they are now, so no existing verdict changes.
`$((` opens arithmetic expansion. The lexer marks it like `$(`, and the edge inventory records this over-refusal.

**The guard package.** A new package `internal/rgguard` owns the verdict. Its `Classify(command string) Verdict` calls `shellcommand.Parse` and reads each span.
It finds the executable with `ResolveRoutinePrefix` (`internal/shellcommand/shellcommand.go:136-168`) and projects the words with `ProjectCommandWords` (`:100-115`).
A span is an `rg` call when the basename of the resolved executable is `rg` and `Executes` is true (`:46`, `:149-155`).
The guard reads one `-c` child level through `ShellChildren` (`internal/shellcommand/projection.go:23-34`) and the exec child through `WorktreeExec` (`:49-61`).
`WorktreeExec` requires that its caller has already accepted the executable as Bench (`:46-48`). A new exported predicate `shellcommand.IsBenchName` owns the Bench executable names `bench` and `bench.sh`. It replaces the private copies in `internal/gitguard/scan.go:76-78` and in the name half of `isBench` (`internal/benchguard/benchguard.go:309-324`), and the guard calls it before `WorktreeExec`.
The package decodes the envelope through `shellcommand.DecodeEnvelope` (`internal/shellcommand/envelope.go:32-63`). The package name joins `guardProjectionGoverned` (`internal/conformance/guard_projection_owner_test.go:29`).

**The verdict order.** For each `rg` call, the guard decides in this order:

1. If the stream is `Unlexed`, the guard permits every call (`internal/shellcommand/shellcommand.go:64-68`).
2. If a token of the span, or the operator token at the span end, carries the substitution mark, the guard refuses with the substitution reason.
3. If the span is uncertain, the guard permits it. The uncertain kinds are `ViaXargs`, a process-substitution operator (`<(` or `>(`) at the span end, an inline `RIPGREP_CONFIG_PATH=` assignment, and an option word that the table does not list.
4. If a `|` or `|&` operator comes directly before the span, the span holds an input redirection, or the call holds an explicit stdin request, the guard permits the call. An input redirection is a `Redirection` token whose text starts with `<` after an optional descriptor digit (`internal/shellcommand/shellcommand.go:52`, `:72-82`).
5. If the call names a no-pattern mode, the guard permits it.
6. If the call has a pattern and no path operand, the guard refuses with the pathless reason. Otherwise the guard permits it.

The operator before a span is `Tokens[span.Start-1]`, the same position that `leadRefusal` reads (`internal/benchguard/benchguard.go:201-222`).
For a `-c` child, `ShellChildren` parses the operand, so the child stream carries its own marks. The first `rg` span of the child also counts as fed when step 4 feeds the outer span.
For an exec child, `WorktreeExec` returns plain argv with no tokens and no marks (`internal/shellcommand/projection.go:46-60`). The guard does not join or parse that argv again.
It reads the child argv as the words of the `rg` call. Step 2 reads the marks of the outer span tokens after `--` and of the operator at the outer span end. Steps 3 and 4 read the operator at the outer span end, the operator before the outer span, and the outer span redirections.

**The option table.** The table lists every ripgrep 15.1.0 option word with its arity, taken from `rg --help`. The research saved the help text with 1617 lines.
The table includes each boolean option word and each `--no-*` negation that the help names, in a flag heading or in prose, so `rg -n --no-heading foo` is refused as pathless.
The value options are the 14 short letters `e f E m j g d t T A B C M r` with their long pairs, the alias `--maxdepth`, and 21 long-only options. The long-only options are `--pre`, `--pre-glob`, `--dfa-size-limit`, `--engine`, `--regex-size-limit`, `--iglob`, `--ignore-file`, `--max-filesize`, `--type-add`, `--type-clear`, `--color`, `--colors`, `--context-separator`, `--field-context-separator`, `--field-match-separator`, `--hostname-bin`, `--hyperlink-format`, `--path-separator`, `--sort`, `--sortr`, and `--generate`.
A long value option takes its value after `=` or from the next word. Ripgrep has no prefix abbreviation for a long option.
In a short bundle, the first value letter takes the rest of the bundle, or the next word when the bundle ends. A value after `=` loses the `=`.
The modes `--files`, `--type-list`, `-h`, `--help`, `-V`, `--version`, `--pcre2-version`, and `--generate` need no path. A bundle that holds `h` or `V` is a mode.
After `-e` or `-f`, every positional is a path. Otherwise the first positional is the pattern and each later positional is a path. After `--`, every word is a positional.
A `-` path operand and the value `-` of `-f` are explicit stdin requests. Step 4 permits a call that holds one, so the substitution refusal of step 2 still wins.
The table is the one source of ripgrep grammar in the tree. `roadmap/FT478.md` can later read it for census classification. This spec makes no census change.

**The refusal text.** `internal/rgguard` owns both refusal sentences and one message format. The pathless message names the missing path and the `.` repair. The substitution message names command substitution and advises single quotes.
Each message ends with `segment=` and the projected words of the call, like the follow-on format (`internal/benchguard/benchguard.go:78`).

**The verb and the shim.** A new registry row `guard-rg` takes `Hook: true`, `internalInventory`, and `axiExempt(axiReasonPlumbing)`, like the guard rows at `cmd/bench/main.go:123-126`.
Its handler in `cmd/bench/guards.go` follows `guardBenchFollowOn` (`cmd/bench/guards.go:152-189`). It recovers a panic to exit 3, exits 2 with the message on stderr on a refusal, and warns and exits 0 on an unreadable envelope.
The new shim `.bench/hooks/block-unsafe-rg.sh` copies the fail-open posture of `block-bench-follow-on.sh` (`:20`, `:24`, `:42-57`). A missing library, a missing core, a stale core, and a core error each warn and exit 0.
The shim carries the four manifest keys that `internal/guards/guards.go:94` requires. Its `# denies:` line is the shipped advertisement of the rule.
The verb takes a `hook.guard-rg` seam (`internal/otelrecord/registry.go:36-42`), a `subcommandRouting` row, and an `entryPointParity` row with a benign input.

**The wiring.** Both `.claude/settings.json` and `.codex/hooks.json` run the shim in their existing PreToolUse Bash group. The group already exists, so `internal/harnesses` takes no change (`internal/conformance/harness_record_test.go:160-181`).
The Claude wiring check (`internal/conformance/line_routing_static_test.go:503-512`) and the Codex wiring check (`internal/conformance/validity_checks_test.go:211-217`) each gain one needle for the shim.
`checkGuardHeaderManifests` gains the shim (`internal/conformance/package_core_checks_test.go:312-323`). `stepGuardFaults` gains the rg verdict (`internal/conformance/refusal_route_guard_test.go:53-68`).
`.bench/consumer-payload.json:18` ships `.bench/hooks` as a tree, so the payload takes no edit.

**The follow-on message.** `benchguard` gains `SideSubstitution`. `refusal` selects it when the operator that `operatorFor` names is a ControlOperator token with the substitution mark (`internal/benchguard/benchguard.go:245-273`).
The repair sentence for that side names command substitution. The operator cell prints `$(` in place of `(`. The verdict refuses exactly the commands that it refuses now.

**The advertisement.** The shim's `# denies:` line and the refusal text become the one source of the explicit-path rule. The bullet at `AGENTS.md:61-63` retires. No anchor pins it (`bench anchors AGENTS.md`, 2026-11-02).
The Hook Layers section of `.bench/BENCH-reference.md` gains one bullet for the rg guard. The section keeps the `bench harnesses codex` anchor at line 891.

## Implementation chunks

Each ticket is a complete green checkpoint on one retained integration source.
C1 adds the shared mark and is its own chunk. Its review closes before any consumer starts.
C2 and C4 each read the mark. C3 wires the C2 verdict into a hook.
The chunks run in table order, because C3 needs C2 and the staged `shell-argument-ownership` spec shares registry files.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| C1 / 01-mark-live-substitution.md | The projection marks each live command substitution and leaves every existing token fact unchanged | R01, R02, R03, R04, R05, R06, R07, R55 | Substitution-mark table, existing Parse tables, guard differential golden | yes |
| C2 / 02-classify-rg-calls.md | `rgguard.Classify` refuses a pathless or substituting rg call and permits each uncertain call | R08, R09, R10, R11, R12, R13, R14, R15, R16, R17, R18, R19, R20, R21, R22, R23, R24, R25, R26, R27, R28, R29, R30, R31, R32, R33, R34, R35, R36, R37, R38, R58 | rgguard verdict tables, projection-owner check | yes |
| C3 / 03-add-guard-rg-verb.md, 04-wire-rg-guard-hook.md | The `guard-rg` verb and its shim refuse in both harnesses, fail open, and appear in the guard surface | R39, R40, R41, R42, R43, R44, R45, R46, R47, R48, R49, R50, R56, R57 | Hook system test, parity, seams, routing, wiring, manifest checks | no |
| C4 / 05-name-substitution-in-follow-on.md | The follow-on refusal names command substitution at a live `$(` and refuses nothing new | R51, R52, R53, R54 | benchguard message tables | no |

### Completion plan

```bench-completion-plan
{
  "version": 2,
  "execution": {"mode": "delegate", "run_id": "ft366-build-20261102", "orchestrator_session": "session_01JK5CPdJaS9Xn1VAmFL1sNN", "author_limit": 1, "review_mode": "unified", "assignments": {"01-mark-live-substitution.md": [], "02-classify-rg-calls.md": [], "03-add-guard-rg-verb.md": [], "04-wire-rg-guard-hook.md": [], "05-name-substitution-in-follow-on.md": []}},
  "chunks": [
    {
      "id": "C1",
      "tickets": ["01-mark-live-substitution.md"],
      "verification": [
        {"id": "mark-tests", "command": "bench test --package ./internal/shellcommand", "probe": "Drop the double-quote branch of the mark and record the R01 red, the restore, and the green.", "ticket": "01-mark-live-substitution.md"},
        {"id": "guard-preservation", "command": "bench test --package ./internal/gitguard --package ./internal/benchguard --package ./internal/census", "probe": "Mark the bare subshell operator and record the R06 red, the restore, and the green.", "ticket": "01-mark-live-substitution.md"},
        {"id": "guard-golden", "command": "bench test --check system --run TestGuardProjectionDifferentialFamilies", "probe": "Lex a | operator run as a Word token and record the R55 red, the restore, and the green.", "ticket": "01-mark-live-substitution.md"}
      ]
    },
    {
      "id": "C2",
      "tickets": ["02-classify-rg-calls.md"],
      "verification": [
        {"id": "verdict-tests", "command": "bench test --package ./internal/rgguard", "probe": "Remove one value letter from the option table and record the R13 red, the restore, and the green.", "ticket": "02-classify-rg-calls.md"},
        {"id": "projection-owner", "command": "bench test --check package-core-guard", "probe": "Plant a private envelope decode in internal/rgguard and record the R37 red, the restore, and the green.", "ticket": "02-classify-rg-calls.md"}
      ]
    },
    {
      "id": "C3",
      "tickets": ["03-add-guard-rg-verb.md", "04-wire-rg-guard-hook.md"],
      "verification": [
        {"id": "registrations", "command": "bench test --check entry-point-parity", "probe": "Remove the guard-rg parity row and record the R56 red, the restore, and the green.", "ticket": "03-add-guard-rg-verb.md"},
        {"id": "wiring", "command": "bench test --check line-routing", "probe": "Remove the shim from the Claude Bash group and record the R39 red, the restore, and the green.", "ticket": "04-wire-rg-guard-hook.md"},
        {"id": "command-contract", "command": "bench test --package ./cmd/bench --run 'HookSeams|GuardRg|CommandRegistry'", "probe": "Remove the hook.guard-rg seam row and record the R46 red, the restore, and the green.", "ticket": "03-add-guard-rg-verb.md"},
        {"id": "routing", "command": "bench test --check subcommand-routing", "probe": "Remove the guard-rg subcommandRouting row and record the R57 red, the restore, and the green.", "ticket": "03-add-guard-rg-verb.md"}
      ]
    },
    {
      "id": "C4",
      "tickets": ["05-name-substitution-in-follow-on.md"],
      "verification": [
        {"id": "message-tests", "command": "bench test --package ./internal/benchguard", "probe": "Return the plain operator side for a marked opener and record the R51 red, the restore, and the green.", "ticket": "05-name-substitution-in-follow-on.md"}
      ]
    }
  ],
  "final_verification": [
    {"id": "system", "command": "bench test --check system"},
    {"id": "coverage", "command": "bench coverage --check rg-path-guard"}
  ]
}
```

## Testing decisions

A good test drives the public `Parse` or `Classify` function with one exact command string and observes the verdict, the reason, and the message.
The shellcommand table tests (`internal/shellcommand/shellcommand_test.go:8-22`) are the prior art for exact token facts. The benchguard row tables (`internal/benchguard/benchguard_test.go:92-137`) are the prior art for verdict rows keyed by row ID.
The rg verdict table keys each case by its coverage row, so a review reads one row against one case.

The rg option table and its expectations both come from ripgrep. The tests are independently authored expectations.
Each such table needs a recorded red under `AGENTS.md`'s exception. The completion plan names one omission probe for each table: the mark table in C1, the verdict table in C2, and the message table in C4.
The new `Token` field breaks the unkeyed `Token` literals of the existing shellcommand tests (`internal/shellcommand/shellcommand_test.go:11-16`, `internal/shellcommand/unlexed_test.go:32`). C1 adds field names to those literals and changes no expected text, kind, or span, so R07 keeps its witness.

The hook seam is the system test that pipes an envelope through the real shim, like `internal/systemtest/bench_follow_on_test.go`.
The guard differential golden (`internal/systemtest/testdata/guard-projection-family-d.golden`) pins the Git and follow-on tuples. It holds no `$(` command and no rg guard entry, so it stays byte-equal and witnesses that the existing guards do not move.

### Seam diagram

    trigger: a Bash tool call in Claude Code or Codex
        │
        ▼
    PreToolUse envelope ──▶ [ block-unsafe-rg.sh ] ──▶ bench guard-rg
                                                          │
                                                          ▼
                              [ rgguard.Classify over shellcommand.Parse ]
                                                          │
                                       exit 0 allow, or exit 2 + BLOCKED message
            ◀ tests attach here: Parse tables for the mark; Classify tables for each verdict;
              one system test pipes an envelope through the shim

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| R01 | 1 | `rg "$(x)" .` yields a Word token `$(x)` with the substitution mark | planned TestParseMarksLiveSubstitution in internal/shellcommand/substitution_test.go | A lexer that marks only unquoted forms leaves this word unmarked |
| R02 | 2 | `` rg "`x`" . `` and `` rg `x` . `` each yield a marked Word token | planned TestParseMarksLiveSubstitution in internal/shellcommand/substitution_test.go | A lexer without a backtick case leaves both words unmarked |
| R03 | 3 | `rg $(x) .` yields a marked `(` ControlOperator token after the Word `$` | planned TestParseMarksLiveSubstitution in internal/shellcommand/substitution_test.go | A mark only on Word tokens misses the unquoted opener |
| R04 | 4 | `rg 'lit$(x)' .` and `` rg 'a`b' . `` yield unmarked tokens | planned TestParseMarksLiveSubstitution in internal/shellcommand/substitution_test.go | A mark from token text alone marks the single-quoted literal |
| R05 | 5 | `rg \$\(x\) .`, `` rg \`x\` . ``, and `rg "\$(x)" .` yield unmarked tokens | planned TestParseMarksLiveSubstitution in internal/shellcommand/substitution_test.go | A mark that ignores escapes marks the escaped literal |
| R06 | 6 | `(rg x .)` and `a $ (b)` yield an unmarked `(` operator | planned TestParseMarksLiveSubstitution in internal/shellcommand/substitution_test.go | A mark on every `(` operator marks a subshell |
| R07 | 7 | Parse keeps the exact token text, kinds, and spans of its existing table | `internal/shellcommand/shellcommand_test.go` (`TestParsePublishesTokenClassesAndSimpleCommands`) | A lexer change that moves text or spans fails the exact DeepEqual |
| R08 | 8 | `Classify("rg foo")` refuses with the pathless reason | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that never refuses passes every other row |
| R09 | 9 | `Classify("rg foo .")` and `Classify("rg foo a b")` permit | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that refuses every rg call fails this row |
| R10 | 10 | `Classify("cat f \| rg foo")` and `Classify("cat f \|& rg foo")` permit | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that ignores the lead operator refuses the fed call |
| R11 | 11 | `Classify("rg foo \| wc -l")` refuses with the pathless reason | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that counts any pipe in the command permits this call |
| R12 | 12 | `Classify("rg foo < f")`, `Classify("rg foo <<< s")`, and a `<<EOF` heredoc call each permit | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that reads only the pipe refuses the redirected call |
| R13 | 13 | `Classify("rg -g '*.go' foo")`, `Classify("rg --glob=*.go foo")`, and `Classify("rg -A 3 foo")` each refuse | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A table without the value option counts the value as the pattern and permits |
| R14 | 14 | `Classify("rg -nA3 foo")` refuses and `Classify("rg -nA 3 foo .")` permits | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A bundle reader that skips the attached value miscounts the positionals |
| R15 | 15 | `Classify("rg -e foo src")` permits and `Classify("rg -e foo")` refuses | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A reader that keeps the first positional as the pattern after `-e` refuses `src` |
| R16 | 16 | `Classify("rg --files")`, `Classify("rg --version")`, and `Classify("rg -nh")` each permit | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard without the mode list refuses a call that never reads stdin |
| R17 | 17 | `Classify("rg -- -foo .")` permits and `Classify("rg -- -foo")` refuses | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A reader that treats `-foo` as an option after `--` miscounts the positionals |
| R18 | 18 | `Classify("rg --zzz foo")` and `Classify("rg -Zq foo")` each permit | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A reader that treats an unknown option as a no-value option refuses |
| R19 | 19 | `Classify("rg foo 'x")` permits | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that copies the fail-closed unlexed policy refuses |
| R20 | 20 | `Classify("ls \| xargs rg foo")` permits | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that ignores `ViaXargs` refuses the call |
| R21 | 21 | `Classify("rg foo <(cmd)")` permits | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that reads only the span words sees no path and refuses |
| R22 | 22 | `Classify("RIPGREP_CONFIG_PATH=c rg foo")` permits | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that skips the assignment refuses the call |
| R23 | 23 | `Classify("rg")` and `Classify("rg -i")` permit | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that refuses a missing path without a pattern refuses these calls |
| R24 | 24 | `Classify("rg \"$(x)\" .")` and `Classify("rg foo $(git ls-files)")` each refuse with the substitution reason | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that checks substitution only on pathless calls permits both |
| R25 | 25 | The substitution message contains `command substitution` and `single quotes` and ends with `segment=` and the call words | planned TestClassifyRgMessages in internal/rgguard/rgguard_test.go | A message without the cause name fails the substring predicate |
| R26 | 26 | `Classify("rg 'a$(b)' .")` and `` Classify("rg \\`x\\` .") `` permit | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that reads raw text for `$(` refuses the literal |
| R27 | 27 | `Classify("rg $(x)")` refuses with the substitution reason | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that decides the pathless reason first names the wrong cause |
| R28 | 28 | The pathless message contains `no path` and `` `.` `` and ends with `segment=rg foo` | planned TestClassifyRgMessages in internal/rgguard/rgguard_test.go | A message without the repair fails the substring predicate |
| R29 | 29 | `Classify("env X=1 rg foo")` and `Classify("timeout 5 rg foo")` each refuse | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that reads only the first word misses the prefixed call |
| R30 | 30 | `Classify("bash -c 'rg foo'")` refuses | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that skips shell children permits the wrapped call |
| R31 | 31 | `Classify("bench worktree exec X -- rg foo")` refuses | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that skips the exec child permits the delegate route |
| R32 | 32 | `Classify("cat f \| bench worktree exec X -- rg foo")` permits | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that reads only the child stream refuses the fed child |
| R33 | 33 | `Classify("cd x && rg foo")` refuses | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that reads only the first span permits the chained call |
| R34 | 34 | `Classify("command -v rg")` permits | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that ignores `Executes` refuses the lookup |
| R35 | 35 | `Classify("rg foo -")` and `Classify("rg -f -")` each permit | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that ignores the explicit stdin request refuses both calls as pathless |
| R36 | 36 | `Classify("bash -c \"bash -c 'rg foo'\"")` and `Classify("sudo rg foo")` permit | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that recurses past one level refuses the nested call |
| R37 | 37 | A private envelope decode planted in `internal/rgguard` reds the projection-owner check | planned TestGuardProjectionOwnerGovernsRgguard in internal/conformance/guard_projection_rgguard_test.go | A governed list without `internal/rgguard` lets the copy pass |
| R38 | 38 | `internal/rgguard` is the only Go package that declares ripgrep option words | review-owned: the C2 review searches the tree for a second ripgrep option list | A second table drifts from the first without a red |
| R39 | 39 | The Claude wiring check reds when the Bash group does not run `block-unsafe-rg.sh` | planned TestClaudeHookWiringNeedsTheRgGuard in internal/conformance/line_routing_static_cases_test.go | A missing needle lets an unwired shim pass |
| R40 | 40 | The Codex wiring check reds when the Bash group does not run `block-unsafe-rg.sh` | planned TestCodexHooksNeedTheRgGuard in internal/conformance/validity_checks_test.go | A missing needle lets an unwired shim pass |
| R41 | 41 | `checkGuardHeaderManifests` reds when the shim lacks a manifest key | planned TestGuardHeaderManifestsCoverTheRgGuard in internal/conformance/guard_header_manifests_test.go | A list without the shim leaves its header ungraded |
| R42 | 42 | The shim with no resolvable core exits 0 and prints a warning | planned TestRgGuardHookFailsOpen in internal/systemtest/rg_guard_hook_test.go | A fail-closed shim exits 2 and blocks every Bash call |
| R43 | 43 | The shim over a core without `guard-rg` exits 0 | planned TestRgGuardHookFailsOpen in internal/systemtest/rg_guard_hook_test.go | A shim that reads exit 2 as a refusal blocks under a stale core |
| R44 | 44 | `bench guard-rg` with an unreadable envelope exits 0 and warns | planned TestGuardRgVerb in cmd/bench/guard_rg_test.go | A verb that refuses an unreadable envelope blocks the call |
| R45 | 45 | The shim with an `rg foo` envelope exits 2 and prints the pathless message on stderr | planned TestRgGuardHookRefuses in internal/systemtest/rg_guard_hook_test.go | A shim that drops the core exit lets the call run |
| R46 | 46 | The hook-seam test reds when the `guard-rg` row has no `hook.guard-rg` seam | `cmd/bench/otel_hook_seams_test.go` (`TestHookSeamsMatchTheHookDispatchRows`) | A verb with no seam row fails the set equality |
| R47 | 47 | `stepGuardFaults` names an rg refusal for a route step `rg foo` | planned TestStepGuardFaultsNamesTheRgGuard in internal/conformance/refusal_route_guard_test.go | A fault list without the rg verdict passes the refused step |
| R48 | 48 | The Hook Layers section names `guard-rg` | review-owned: the C3 review reads the Hook Layers section | Review-owned prose; the row records the obligation |
| R49 | 49 | `AGENTS.md` holds no bullet that states the explicit-path rule | review-owned: the C3 review searches `AGENTS.md` for the retired sentence | A surviving bullet is a second source of the rule |
| R50 | 50 | The link plan copies `.bench/hooks/block-unsafe-rg.sh` | review-owned: the C3 review confirms the tree row at `.bench/consumer-payload.json:18` | The payload ships the hooks tree, so the obligation is a review read |
| R51 | 51 | `Classify("bench gate $(x)")` refuses with a message that names `command substitution` and ends `operator=$(` | planned TestClassifyNamesCommandSubstitution in internal/benchguard/substitution_test.go | A message that keeps the plain side prints `operator=(` |
| R52 | 52 | `Classify("echo $(bench status)")` refuses with a message that names `command substitution` | planned TestClassifyNamesCommandSubstitution in internal/benchguard/substitution_test.go | A side choice only after the span misses the lead opener |
| R53 | 53 | `Classify("bench learning t --what \"$(cat f)\"")` stays unblocked | planned TestClassifyNamesCommandSubstitution in internal/benchguard/substitution_test.go | A guard that refuses every marked token blocks this call |
| R54 | 54 | `Classify("true; (bench status)")` keeps the message that ends `operator=(` | planned TestClassifyNamesCommandSubstitution in internal/benchguard/substitution_test.go | A side choice on every `(` renames the subshell refusal |
| R55 | 55 | The guard differential run matches its golden byte for byte | `internal/systemtest/guard_projection_differential_test.go` (`TestGuardProjectionDifferentialFamilies`) | A projection change that moves a Git or follow-on tuple fails the golden |
| R56 | 46 | A tree whose hooks directory holds a shim with no parity row reds the entry-point-parity check | `internal/conformance/entry_point_parity_test.go` (`TestEntryPointParityNamesAShimOutsideTheTable`) | A `block-unsafe-rg.sh` with no `entryPointParity` row fails the live check by the same rule |
| R57 | 46 | A registry name with no `subcommandRouting` row reds the routing check | `internal/conformance/subcommand_routing_test.go` (`TestSubcommandRoutingRegistryBites`) | A `guard-rg` row with no routing row fails the live check by the same rule |
| R58 | 58 | `Classify("bench worktree exec X -- rg \"$(x)\" .")` and ``Classify("bench worktree exec X -- rg `x` .")`` each refuse with the substitution reason | planned TestClassifyRgCalls in internal/rgguard/rgguard_test.go | A guard that reads marks only from the child argv sees none and permits both calls |

Not covered: story 56 — an exported `RIPGREP_CONFIG_PATH` is not in the command text, so no guard input reaches it.
Not covered: story 57 — the digit-only operand loss is pre-existing in `ProjectCommandWords` (`internal/shellcommand/shellcommand.go:104`), and this spec keeps that contract.

### Edge inventory

The guard reads one string and returns one verdict, so the cleanup and the concurrency classes do not apply.
The hostile-input walk covers quotes, escapes, substitutions, heredocs, operators, prefixes, wrappers, and malformed input in rows R01-R36.

- **Won't handle:** an exported `RIPGREP_CONFIG_PATH` in the session environment — the command text does not show it, and the inline assignment form stays covered by R22.
- **Won't handle:** a digit-only pattern or path before a redirection, such as `rg 404 > out` — `ProjectCommandWords` drops it today (`internal/shellcommand/shellcommand.go:104`), and R08 covers the ordinary pathless call.
- **Won't handle:** a wrapper below one level, and the unknown prefixes `sudo`, `time`, and `nice` — the existing guards read the same depth, and R30 and R31 cover the one-level routes.
- **Won't handle:** arithmetic expansion `$((...))` — the lexer marks it like `$(`, so an rg call with it in double quotes is refused as a substitution; R24 covers the ordinary form.
- **Won't handle:** a substitution in a command that the lexer refuses — step 1 permits every call in an unlexed stream, and R19 keeps that permit; R24 covers each lexed substitution.
- **Won't handle:** an escaped `\$(` inside a double-quoted `-c` operand, which the outer shell unescapes into a live substitution — the child parse sees an escaped literal, and R30 covers the ordinary `-c` child.
- **Won't handle:** a backtick span that holds a double quote inside a double-quoted word — the lexer ends the word at that quote today, and R02 covers the ordinary span.
- **Won't handle:** a substitution inside a heredoc body — `stripHeredocBodies` removes the body before the lexer runs (`internal/shellcommand/shellcommand.go:369-397`), and R12 covers the heredoc as a stdin supplier.
- **Won't handle:** a ripgrep version other than 15.1.0 — the table records the installed grammar, and R18 permits an option that a later version adds.

## Out of scope

- A census column for sanctioned `rg` reads (`roadmap/FT478.md`). It can read the option table of this spec. Estimate: 6 edits, 2 gate runs.
- A quote-aware heredoc stripper. A quoted `<<X` in an ordinary word removes later lines from every guard's view today. Estimate: 4 edits, 2 gate runs.
- A link from the gitguard deny table to the shim `# denies:` header. The header claims to derive from the table, and no code reads it (`internal/gitguard/gitguard.go:3-4`). Estimate: 5 edits, 2 gate runs.

## Further notes

### Flagged additions

Each item below goes beyond the literal decision source. The reviewer can veto each one.

1. The input-redirection exemption (R12). The source names only the pipe. Ripgrep reads a `<` redirect in place of a pipe, so a refusal would be a false positive.
2. The substitution mark in the shared lexer (R01-R07). D1 and D2 need one fact for two guards.
3. `-` counts as a path (R35). The caller names stdin on purpose.
4. The `AGENTS.md:61-63` bullet retires (R49). The guard header and the refusal become the one source.
5. The Codex wiring and its needle (R40). The source says "a PreToolUse guard" and names no harness.
6. `stepGuardFaults` grades route steps against the rg verdict (R47).
7. The permits for `sudo rg foo` and for a nested `-c` wrapper (R36). The source says only that the guard is not evasion-resistant; these rows pin the one-level depth of the existing guards.
8. The retirement in item 4 also removes the rule for a harness that runs no hook, such as an unadopted OpenCode session, because every harness reads `AGENTS.md`. The reviewer decides whether the bullet stays as a pointer.

### Seams beyond the decision source

- The substitution mark on `Token`: `internal/shellcommand/shellcommand.go:23-26` and `:299-365`. The internals readers found that quote folding erases the quote style and that no guard reads word text.
- The `SideSubstitution` side in benchguard: `internal/benchguard/benchguard.go:245-273`. The internals readers traced `operator=(` to `operatorFor`.
- The `shellcommand.IsBenchName` predicate: `internal/shellcommand/projection.go:46-48`, `internal/gitguard/scan.go:76-78`, and `internal/benchguard/benchguard.go:309-324`. The ticket review found that the exec child needs a Bench test, and that two private copies of the names already exist.
- The `internal/rgguard` package: `internal/conformance/guard_projection_owner_test.go:29`. The enforcement readers found that policy lives in each guard package and that a new package must join the governed list.

### Source-sentence-to-row table

| Source clause | Coverage |
|---|---|
| A new PreToolUse guard uses the shared projection to inspect `rg` calls. | R29, R30, R31, R33, R37, R39, R40, R45 |
| It requires a path operand unless a pipe in the same command supplies stdin. | R08, R09, R10, R11, R32 |
| It recognizes options that consume a value. | R13, R14, R15, R16, R17 |
| It permits a call when that operand analysis is uncertain. | R18, R19, R20, R21, R22 |
| This new guard's uncertainty policy does not weaken the existing dangerous-command guards. | R07, R55 |
| The guard names command substitution in its refusal. | R24, R25, R27 |
| It treats the quoted and unquoted forms the same way. | R01, R02, R03, R24 |
| It is an honest-mistake check, not an evasion-resistant shell interpreter. | R36 |
| Its projection never executes the command. | review-owned in C2: `Classify` takes a string and starts no process |
| D1 (2026-11-02): refuse every live substitution; literal forms pass. | R24, R26, R58 |
| D2 (2026-11-02): the follow-on refusal names command substitution and refuses nothing new. | R51, R52, R53, R54 |
| Occurrence 2026-10-03: a pathless rg waited about two hours. | R08, R31 |
| Occurrence 2026-10-07: `operator=(` for an unquoted `$(` in a Bench call. | R51 |
| Occurrence 2026-10-10: a backtick span in a double-quoted pattern ran `bench handoff`. | R02, R24 |

### Pre-review proof checklist

- Cited symbols: `Parse`, `lex`, `Token`, `ProjectCommandWords`, `ResolveRoutinePrefix`, `ShellChildren`, `WorktreeExec`, `DecodeEnvelope`, `operatorFor`, `refusal`, `Verdict.Message`, `guardBenchFollowOn`, and `stepGuardFaults` were opened at their cited lines.
- Import edges: `internal/rgguard` imports `internal/shellcommand`, and `cmd/bench` imports `internal/rgguard`. `internal/conformance` imports `internal/rgguard` for R47. No cycle exists, because `internal/rgguard` imports no package that imports it.
- Source-row clauses and occurrences: each clause and each occurrence is in the source table.
- Promised field labels: the field `Token.Substitution`, the side `SideSubstitution`, the verb `guard-rg`, and the seam `hook.guard-rg`.
- Changed-function callers: `lex` has one caller, `Parse`. `Parse` callers keep their contract (R07, R55). `operatorFor` has one caller, `refusal`.
- Copy survival: `shellcommand.IsBenchName` replaces the name copies in `gitguard` and `benchguard`; the C2 review searches `internal/gitguard`, `internal/benchguard`, and `internal/rgguard` for a surviving `"bench.sh"` literal (review-owned under R38's review). Path literals in other packages state other facts and stay.
- Rendered-shape readers: the follow-on message changes only for a marked opener. `rg -n 'operator=' internal/ cmd/ tests/` finds no row with `$(` or a `(` opener, so no existing pin moves.
- Pin operators: R25 and R28 use substring containment and a suffix; R51 uses containment and a suffix; R55 uses byte equality.
- Entry reads: `bench guard-rg` reads the envelope on stdin; `Classify` takes the command as a value.
- Derived expectations: every expectation comes from ripgrep 15.1.0 help and probes or from POSIX shell quoting, never from the code under test.
- Consolidated rules: none.
- Quantified obligations: "every rg call" holds for each simple command (R33), each wrapper child (R30), and each exec child for the pathless verdict (R31) and the substitution verdict (R58).
- Workflow-step writes: none.

### Captured learnings

- The heredoc stripper ignores quotes (`internal/shellcommand/shellcommand.go:367-397`), so `rg "<<X" .` followed by a `git reset --hard` line passes both guards. The coordinator records this with `bench learning`.
- The gitguard header says the deny table feeds the guard manifest, and no code reads it (`internal/gitguard/gitguard.go:3-4`).

### Ordering note

The staged spec `shell-argument-ownership` also writes `cmd/bench/command_registry.go`, `internal/guards/guards_test.go`, and `internal/conformance/subcommand_routing_table_test.go`. Its build and this build run in sequence, never at the same time.
