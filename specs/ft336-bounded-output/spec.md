# Bounded default output for every Bench verb

Status: staged

Roadmap: FT336

Decision source: reviewer-confirmed current conversation, 2026-09-24, on the named reviewed artifact `roadmap/FT336.md`.

Verification log: 2 iteration(s) to accept — Opus/high iteration 1 returned 4 blocking and 10 non-blocking findings. Iteration 2 closed all 14 and accepted, and the author folded its 3 non-blocking notes. A pre-draft Fable/high consultation changed the help exemption, the spill-failure fallback, and the exec owner rule.

## Problem

A Bench response can print more than 100 lines. The agent pays for every line in its context, and it pays again at each later request. The 2026-09-23 CLI audit found that `bench worktree list` printed about 138 lines, and 90 of them were help rows. `bench preflight build` printed 15 green check rows.

`bench worktree exec` is the only route into a worktree. The follow-on hook refuses a pipe or `tail` after a Bench call. So every byte that an exec child prints reaches the agent context.

The FT71 build session made 1,126 model requests at an average context of about 445k tokens. The cost is the number of requests multiplied by the context. A fix that moves the needed part of a response behind a second call adds requests.

## Solution

One owner holds each public Bench response to 10 lines. A response of 10 lines or fewer prints unchanged. A longer response prints its first 4 lines, one spill line, and its last 5 lines. The spill line names a private file that holds the complete output. The verb keeps its exit code.

The `cmd/bench` dispatcher applies the bound to each public command. `bench worktree exec` passes its child's output through the same bound. The help forms, `bench dashboard --stdout`, `bench worktree shell`, and `bench setup` stay exempt.

Three verbs change their default shape:

- `bench worktree list` prints one inspect action and one exec action with a `<target>` slot, in place of two help rows for each active row.
- `bench preflight build` and `bench preflight review` print one summary line in place of the green and not-applicable rows.
- `bench preflight evidence <id>` prints the manifest summary and the next page command. The new `--to <dir>` form exports each verified source to a file.

`bench commit --preflight-build <slug>` commits, builds the worktree, and runs build preflight in one call. The census records the output lines and bytes of each Bench response in an assignment, and the landing prints the breakdown. A byte bound joins the line bound after the queries ticket-4 budget decision exists.

## User stories

Line: opus / high.
Implementation-line reason: The response owner and the dispatcher posture change are the hardest material chunks. The spec fixes the projection exactly, but stream order, interrupts, and the posture reds across the command tests are uncertain. Owner-level tests are strong, and the dispatcher reds need a gate run to find.
Harder chunks: BO-C1, BO-C2.

### The response bound

1. As an agent, I want each public Bench response held to 10 lines, so that one call cannot flood my context.
2. As an agent, I want a response of 10 lines or fewer printed unchanged, so that a small answer keeps its bytes.
3. As an agent, I want an over-bound response to show its first 4 and last 5 lines, so that the verdict stays visible.
4. As an agent, I want one spill line with the path, the totals, and the omitted count, so that I see the gap.
5. As an agent, I want the spill file to hold the complete output in arrival order, so that I can read every omitted line.
6. As an agent, I want the exit code unchanged by the bound, so that a red verb still reads red.
7. As an agent, I want the help forms exempt, so that the inventory and each grammar stay readable in one call.
8. As an operator, I want `bench dashboard --stdout`, `bench worktree shell`, and `bench setup` exempt, so that a machine artifact and an interactive terminal keep working.
9. As a hook author, I want internal plumbing and hook commands outside the bound, so that the hook envelopes stay complete.
10. As a reviewer, I want each spill file created exclusively in a private directory, so that a hostile name cannot redirect it.
11. As an operator, I want an assignment's spill files removed when the assignment retires, so that spills do not accumulate.
12. As an operator, I want the spills outside an assignment limited to the newest 64 files, so that the primary checkout's spills stay bounded.
13. As an agent, I want a spill that cannot start to print the complete output, so that no byte is lost.
14. As an agent, I want a midway spill failure to keep every byte in the file or the response, so that nothing is lost.
15. As a reviewer, I want one owner for the line value and the projection, so that the verbs and exec cannot drift.
16. As a reviewer, I want every public command to declare its bound disposition, so that a new command cannot skip the bound.
17. As an agent, I want an empty response to print nothing and create no spill, so that silence stays silence.
18. As an agent, I want a last line without a newline counted as a line, so that the count matches the bytes.
19. As an agent, I want binary and invalid UTF-8 bytes kept exactly in the spill file, so that the spill is a byte copy.
20. As an agent, I want an unsafe spill path to print the complete output, so that a control byte cannot forge a line.

### The exec route

21. As an agent, I want `bench worktree exec` to bound its child's output, so that the only worktree route cannot flood my context.
22. As an agent, I want the child's exit code returned unchanged, so that exec keeps its exit contract.
23. As an agent, I want stdin still forwarded to the child, so that the heredoc form keeps working.
24. As an agent, I want an interrupted child's output kept, so that I can read what ran before the interrupt.
25. As an agent, I want a nested Bench verb bounded in its own process, so that each process owns one bound.
26. As an agent, I want exec memory bounded while a child prints many lines, so that a large child cannot exhaust memory.
27. As an agent, I want exec's grammar refusal unchanged, so that I can still tell a refusal from the child's own exit 2.

### The per-verb defaults

28. As an agent, I want `bench worktree list` to print its active-row actions once with a `<target>` slot, so that 47 rows print no 94 help rows.
29. As an agent, I want a cleanup-pending row to keep its row-specific release action, so that its request token and path stay available.
30. As an agent, I want a missing-tree row to keep its row-specific recovery action, so that its recovery stays available.
31. As an agent, I want a foreign row to keep its row-specific clean action, so that its orphan path stays available.
32. As an agent, I want the id cell to stay the address that `bench worktree path` accepts, so that the slot has one known filler.
33. As an agent, I want `bench preflight build` to replace its green rows with one summary line, so that a green verdict costs one line.
34. As an agent, I want `bench preflight review` to print the same summary line, so that a green review verdict costs one line.
35. As an agent, I want each red check row printed in full, so that a red verdict names its check, detail, and next action.
36. As a maintainer, I want the AXI guidance for `bench worktree list` to state the slot rule, so that the list does not print per-row help again.

### The evidence default face

37. As an agent, I want `bench preflight evidence <id>` to print the manifest summary and the next page command, so that the first call costs no source bytes.
38. As a consumer, I want the next page command to return the page that the old default returned, so that retrieval keeps one successor chain.
39. As a consumer, I want each cursor, source, verify, and check-current form unchanged, so that ADR 0022 retrieval keeps working.
40. As an agent, I want `--to <dir>` to export each verified source to its own file, so that I can read the sources with file tools.
41. As an agent, I want `--to` to refuse an existing directory that is not empty, so that the export cannot overwrite files.
42. As an agent, I want `--to` to refuse a symlink at the directory path, so that the export cannot write outside the named directory.
43. As an agent, I want `--to` to write a source only after its digests verify, so that the export holds only verified bytes.
44. As a reviewer, I want `--to` to name each file by source ordinal, so that a hostile source id cannot choose a path.
45. As an agent, I want a failed export to remove what it created, so that a partial export cannot pass for a complete one.

### The commit chain

46. As an agent, I want `bench commit --preflight-build <slug>` to commit, build, and run build preflight, so that one checkpoint costs one request.
47. As an agent, I want a commit that publishes nothing to skip the build and the preflight, so that no follow-on grades an unpublished tree.
48. As an agent, I want a failed build to skip the preflight and still name the published commit, so that I know the commit landed.
49. As an agent, I want the chain's exit code to be the first non-zero step exit, so that a red step stays red.
50. As an agent, I want `--dry-run` with the flag refused at exit 2, so that a dry run cannot build.
51. As an existing caller, I want `bench commit` without the flag unchanged, so that current invocations keep their behavior.

### The census output record

52. As a reviewer, I want the size of each bounded response in an assignment recorded, so that an audit can measure each verb's cost.
53. As a reviewer, I want an exec run from the primary checkout recorded under its target assignment, so that the outer response reaches the census.
54. As a reviewer, I want output records kept apart from raw-call records, so that the raw-call count does not change.
55. As a reviewer, I want the landing to print one output breakdown line, so that the retro reads the context cost per verb.
56. As an operator, I want output records dropped at assignment retirement, so that they follow the raw-call lifecycle.
57. As an agent, I want a failed census write to change neither the response nor the exit code, so that evidence cannot block a verb.

### The byte bound

58. As an agent, I want each bounded response held to the reviewed byte value as well, so that one long line cannot flood my context.
59. As a reviewer, I want the byte value taken from the queries ticket-4 budget decision, so that FT336 invents no number.
60. As a reviewer, I want the byte-bound ticket stopped before product writes while that decision is absent, so that no placeholder value ships.

### Reviewed exclusions and evidence

61. As a reviewer, I want the refusal of a second Bench call in one line unchanged, so that closed decision 1 holds.
62. As a reviewer, I want `bench consumers` to keep one re-query action per candidate row, so that its anchored decision holds.
63. As a reviewer, I want the audit to record each verb's size and acted-on part, so that the fix list comes from evidence.

### Review-round additions

64. As an agent, I want exec to return at its child's exit, so that a background process that holds a pipe cannot hang exec.
65. As a release operator, I want the ship-tier commands exempt, so that the CI log keeps the complete release evidence.
66. As an agent, I want a commit that exits 3 to stop the chain and name its commit, so that I can reconcile the checkout.

## Implementation decisions

### Closed decisions, 2026-09-24

The reviewer closed these decisions on 2026-09-24. They stay closed.

1. Bench call chains: the refusal of a second Bench call in one line stays. The `.bench/BENCH.md` rule and `block-bench-follow-on.sh` do not change. The commit chain removes the need for doubled calls.
2. ADR 0022 is not superseded. `bench preflight evidence` keeps the ADR 0022 pages. The default response prints only the manifest summary and the next page command, with no source bytes. `--to <dir>` is an additive export and does not replace the verified pages.
3. The compound action is a flag on `bench commit`: commit, then the worktree build, then build preflight. It is not a new verb.
4. FT336 waits for the measurement. The 10-line bound stands. FT336 invents no numeric byte value. Each byte bound consumes the value from the measurement report and the reviewer budget decision that `specs/session-context-queries/tickets/4-apply-reviewed-budgets.md` owns. Any FT336 ticket that needs a byte bound waits for that decision.

### The response owner

A new package owns the response bound. The line value, 10, sits in the production policy registry of `internal/bounds` beside the other fixed bounds. No other package states the value.

The owner accepts an ordered stream of writes, each tagged stdout or stderr. Its two tagged writers serialize their writes, because `os/exec` copies the two child streams on two goroutines. It keeps the first 11 lines in memory. When the eleventh line arrives, it creates the spill file, writes every retained byte, and then streams each later byte to the file. It keeps only the head lines and a ring of the last 5 lines in memory after that point.

At the finish, a response of 10 lines or fewer replays each write to its own stream in arrival order. A longer response prints the head, the spill line, and the tail on stdout, in that order. The spill file holds the complete output, both streams, in arrival order. The spill line has this exact form: `spilled{lines=<total>,bytes=<total>,omitted_lines=<n>,path=<absolute path>}`.

A line ends at a newline byte. A final line without a newline counts as one line. The owner counts bytes, not runes, and it copies binary bytes exactly.

### The spill store

The spill store sits under the Bench home at `responses/<repo-key>/<scope>/`. The repo key is the census key of the repository. The scope is the assignment id when the process runs in an assignment worktree, and `primary` otherwise. A retiring verb (`bench worktree release`, `clean`, `reclaim`, or `land`) writes its spill to the `primary` scope. So the retirement cannot remove a spill that is still open. A process outside any repository uses the repo key `none`.

The owner creates each directory with mode 0700 and never follows a symlink. It creates each file with an exclusive create at mode 0600. A generated name selects the file, so no operand or output byte forms a path. The retirement path that drops an assignment's census records also removes that assignment's spill directory. The `primary` and `none` scopes keep the newest 64 files, and the owner removes older files after each new spill.

### Spill failure

If the spill file cannot be created, the owner writes every retained byte and every later byte to the original streams. It then prints one line: `spill-failed{reason=<reason>}`. If a write fails after the spill started, the owner prints the head, then `spill-failed{path=<path>,written_bytes=<n>,reason=<reason>}`, then every byte that the file did not receive. The concatenation of the file and the response then equals the complete output. A spill path that is not line-safe takes the create-failure route.

### The dispatcher and the dispositions

`Command.Run` in `cmd/bench` is the one dispatch site. Each public registry entry declares a bound disposition: bounded, or exempt with a reason. The dispatcher gives a bounded command one owner that both its stdout and its stderr write into. It finishes the owner after the command returns, and it returns the command's own exit code. A hook or internal plumbing entry is outside the bound.

The exempt set is closed:

- The help forms: `bench help`, and a command or a leaf followed by exactly one `--help`, `-h`, or `help` argument. A longer argument list stays bounded.
- `bench dashboard --stdout`, because its stdout is a machine artifact.
- `bench worktree shell` and `bench setup`, because each one drives an interactive terminal.
- The ship-tier commands `release-preflight`, `prep-release`, and `release`, because CI runs them and a spill file on a discarded runner loses the evidence.
- `repair`, because only the wrapper runs it, and it never reaches `Command.Run`.

Each process owns one owner. `bench worktree exec` gives its child the owner's two writers, so Go copies the child's output through pipes into the owner. A Bench verb inside the child bounds its own output in its own process. Exec's outer bound then applies to the child's complete output.

A descendant of the child can keep a pipe open after the child exits. So exec sets a wait delay on the child command, and the policy registry of `internal/bounds` holds that delay. Exec then returns at the child's own exit with the child's exit code, and the owner keeps the output that arrived before the pipes closed. The interrupt path uses the same delay.

A leaf's disposition sits on its row in the leaf family table, because that row is the single declaration of the leaf. The first ticket declares only the `worktree exec` leaf bounded. Every other public entry declares a transitional `pending` disposition. The second ticket bounds every public entry and removes the `pending` value. This order is an expand, then a contract.

### The per-verb defaults

`bench worktree list` renders one `bench worktree path <target>` action and one `bench worktree exec <target> -- <command>` action when one or more active rows with a present tree exist. The cleanup-pending, missing-tree, and foreign rows keep their row-specific actions, because their path operand appears in no cell. The id cell stays the address that `bench worktree path` accepts.

`bench preflight build` and `bench preflight review` keep the `phase`, `spec`, and `source` lines. The check table then becomes one line: `checks{green=<n>,not_applicable=<n>,red=<n>}`. When one or more checks are red, a `checks[<n>]{check,verdict,detail,next}` table follows with the red rows only. The charge forms keep their complete check table in the evidence artifact.

### The evidence default face

`bench preflight evidence <id>`, with no cursor and no source, prints one `evidence_summary` block. Its fields are the evidence identity, the source count, the page count, the manifest bytes, the total source bytes, and `next`. `next` is the exact command that reads the first manifest page. The chargeevidence response schema registers the block, so the encoded-response bound of ADR 0022 applies to it. The `next` of a `--charge` preparation response stays the bare `bench preflight evidence <id>` command, so a consumer meets the summary first.

`bench preflight evidence <id> --to <dir>` exports every declared source. A relative directory resolves against the working directory. The directory must be absent or empty, and it must not be a symlink.

The export verifies each page digest and each source digest before it writes that source. It writes `source-<ordinal>` for each source and one `index.toon` table with the ordinal, the source id, the bytes, and the file name. On any failure, it removes each file it created, and it removes the directory if it created it. It prints one `exported{sources=<n>,bytes=<n>,dir=<absolute path>}` line.

### The commit chain

`bench commit -m <msg> --preflight-build <slug> <path>...` runs three steps in order: the commit, the worktree build of the current assignment, and `bench preflight build <slug>` at the published tip. The composition sits in the `cmd/bench` command layer, so no package import cycle forms. The commit package returns the published commit to that layer, and the chain line takes the sha from that value. Each step prints its own response. The chain then prints one line: `commit-chain{commit=<sha|none>,build=<green|red|skipped>,preflight=<green|red|skipped>}`.

A commit that publishes nothing skips both later steps. A red build skips the preflight. The chain returns the first non-zero step exit, or 0. A commit that exits 3 stops the chain, and the chain line names its published commit. `--dry-run` with `--preflight-build` is a usage refusal at exit 2.

### The census output record

The census writes output records to a second file, `<id>.output`, beside the raw-call file of the same assignment. Each line holds the time, the verb head, the line count, the byte count, and `spilled` or `inline`. The verb head is `bench`, the command name, and the leaf word for a command that has leaves. The raw-call readers read only the `<id>` file, so their counts do not change.

The dispatcher writes the record after the owner finishes. It uses the assignment of the working tree, or the target assignment of `bench worktree exec`. `ExecCommand` reports the assignment that it resolved, so the dispatcher runs no second resolution. A process with neither writes no record. The landing prints `census output{<head>=<calls>/<bytes>,...}` on stderr beside the `census heads` line. `census.Drop` removes both files, so retirement keeps one call site.

### The byte bound

The byte-bound ticket starts only after the queries ticket-4 measurement report and the reviewer budget decision record exist. The ticket then records the approved value, or the per-surface values, in the production policy registry. The owner treats a response as over-bound when its bytes exceed the value, and its projection stays within the value. The ticket chooses no value.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| BO-C1 / `1-bound-exec-output.md` | The response owner exists, and it bounds every `bench worktree exec` child. | BO1, BO2, BO3, BO4, BO5, BO6, BO7, BO14, BO15, BO18, BO19, BO20, BO21, BO22, BO23, BO24, BO25, BO26, BO27, BO28, BO29, BO31, BO68, BO69 | `bench test --package ./internal/responsebound`, `bench test --package ./cmd/bench`, `bench test --check system` | yes |
| BO-C2 / `2-bound-every-public-response.md`, `3-retire-response-spills.md` | Every public response obeys the bound, and the spills follow the assignment lifecycle. | BO8, BO9, BO10, BO11, BO12, BO13, BO16, BO17, BO30, BO66, BO70, BO72 | `bench test --package ./cmd/bench`, `bench test --package ./internal/worktree`, `bench test --package ./internal/responsebound`, `bench test --check system` | yes |
| BO-C3 / `4-slot-worktree-list-actions.md`, `5-summarize-green-preflight.md` | The list prints slot actions, and a green preflight prints one line. | BO32, BO33, BO34, BO35, BO36, BO37, BO38, BO39, BO40, BO41, BO67 | `bench test --package ./internal/worktree`, `bench test --package ./internal/preflight`, `bench test --package ./internal/anchors`, `bench test --package ./internal/consumers` | no |
| BO-C4 / `6-summarize-evidence-default.md`, `7-export-evidence-sources.md` | The evidence default prints a summary, and `--to` exports verified sources. | BO42, BO43, BO44, BO45, BO46, BO47, BO48, BO49, BO50 | `bench test --package ./internal/preflight/evidencecmd`, `bench test --package ./internal/chargeevidence` | no |
| BO-C5 / `8-chain-commit-preflight.md` | One commit call also builds the worktree and runs build preflight. | BO51, BO52, BO53, BO54, BO55, BO56, BO71 | `bench test --package ./cmd/bench`, `bench test --package ./internal/commit` | no |
| BO-C6 / `9-record-response-census.md` | The census records the response sizes, and the landing prints them. | BO57, BO58, BO59, BO60, BO61, BO62 | `bench test --package ./internal/census`, `bench test --package ./internal/worktree`, `bench test --package ./cmd/bench` | no |
| BO-C7 / `10-apply-byte-bound.md` | Each bounded response obeys the reviewed byte value. | BO63, BO64, BO65 | `bench test --package ./internal/responsebound` | no |

Ticket 1 creates the owner that tickets 2, 3, 9, and 10 consume, so BO-C1 stays one small chunk, and its review closes first. Tickets 1, 2, 3, 4, 7, 8, and 9 write `cmd/bench` registry or help files. The orchestrator lands them in ticket-number order inside that shared set. BO-C7 stays blocked by its entry stop until the budget decision exists.

## Testing decisions

- A good owner test writes a canned sequence of tagged writes and compares the printed bytes, the spill file bytes, and the exit code. It does not read the owner's internal state.
- A good dispatcher test calls `Command.Run` with buffers and a test-registered public command, then compares the bytes. `TestCommandRunRoutesRegisteredCommands` style tests in `cmd/bench` are the prior art.
- Exec rows that need a real child process, a signal, or stdin run through `bench test --check system`. `internal/systemtest/bench_follow_on_test.go` is the prior art for a real Bench subprocess.
- Each owner test sets the Bench home to a private temporary directory, so no test writes the fallback home.
- The package tests run in the gate's test phase, and the system rows run in its system phase.

### Posture change: tests that the bound reds

Ticket 2 bounds every public response. So each test that reads more than 10 lines of public output from `Command.Run` or the binary turns red. The ticket author finds these reds with one gate run before the first edit. The author then changes each test to read the spill file or to call the package seam. A test that asserts a verb's complete output keeps its complete assertion through the spill file. The author changes no assertion to a weaker predicate.

### Seam diagram

    trigger: an agent runs a public Bench verb, or a child through bench worktree exec
        │
        ▼
    stdout and stderr writes  ──▶  [ response owner ]  ──▶  bounded response + spill file + census record
                                      ◀ tests attach here: canned tagged writes in, printed bytes and file bytes out

    trigger: an agent runs bench commit --preflight-build <slug>
        │
        ▼
    commit ──▶ worktree build ──▶ build preflight  ──▶  step responses + commit-chain line
                      ◀ tests attach here: Command.Run with injected step seams, exit code and chain line out

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| BO1 | 1 | A test-registered public command that prints 25 lines through `Command.Run` produces exactly 10 stdout lines | planned TestDispatcherBoundsPublicResponse in cmd/bench | A dispatcher that skips the owner prints 25 lines |
| BO2 | 2 | A public command that prints exactly 10 lines produces its exact bytes and creates no spill file | planned TestDispatcherPassesBoundaryResponse in cmd/bench | An off-by-one bound spills the tenth line |
| BO3 | 3 | An 11-line response prints lines 1 to 4, the spill line, and lines 7 to 11 | planned TestOwnerProjectsHeadAndTail in internal/responsebound | A head-only or tail-only projection drops the other end |
| BO4 | 4 | A 25-line response of 300 bytes prints `spilled{lines=25,bytes=300,omitted_lines=16,path=<abs>}` as its fifth line | planned TestOwnerSpillLineFields in internal/responsebound | A wrong count or a relative path fails the exact line match |
| BO5 | 5 | The spill file bytes equal the complete input bytes | planned TestOwnerSpillHoldsCompleteOutput in internal/responsebound | A spill of the omitted lines alone fails the byte comparison |
| BO6 | 5 | Alternate stdout and stderr writes of 12 lines appear in the spill file in write order | planned TestOwnerKeepsArrivalOrder in internal/responsebound | Two separate stream buffers reorder the lines |
| BO7 | 6 | A bounded command that prints 30 lines and exits 3 returns exit 3 | planned TestDispatcherKeepsExitCode in cmd/bench | An owner that maps the exit to 0 or 1 fails the code match |
| BO8 | 7 | `bench help` through `Command.Run` prints every line that `renderCommandHelp` renders | planned TestHelpFormsStayComplete in cmd/bench | A help form under the bound prints 10 lines |
| BO9 | 7 | `bench worktree --help` through `Command.Run` prints its complete grammar | planned TestHelpFormsStayComplete in cmd/bench | A leaf help form under the bound spills the grammar |
| BO10 | 7 | A bounded command given two arguments that end in `--help` stays bounded | planned TestHelpExemptionNeedsOneArgument in cmd/bench | A suffix match exempts an exec child's own `--help` |
| BO11 | 8 | `bench dashboard --stdout` through `Command.Run` prints the complete page | planned TestDashboardStdoutStaysComplete in cmd/bench | A bounded dashboard stream breaks the artifact |
| BO12 | 8 | The exempt set in the registry equals the help forms, `dashboard --stdout`, `worktree shell`, `setup`, the three ship-tier commands, and `repair`, and each exemption names a reason | planned TestBoundExemptionsAreClosed in cmd/bench | An added exemption or a missing reason fails the set comparison |
| BO13 | 9 | The internal `guard-git` command prints a 30-line stderr in full | planned TestPlumbingStaysOutsideBound in cmd/bench | A dispatcher that bounds hook commands truncates the envelope |
| BO14 | 10 | A new spill directory has mode 0700 and a new spill file has mode 0600 | planned TestSpillStorePrivateModes in internal/responsebound | A default umask create leaves the file readable by others |
| BO15 | 10 | A symlink at the spill scope directory makes the owner take the create-failure route | planned TestSpillStoreRefusesSymlink in internal/responsebound | A followed symlink writes outside the Bench home |
| BO16 | 11 | The retirement of an assignment removes its spill directory | planned TestRetirementDropsResponseSpills in internal/worktree | A retirement that drops only the census leaves the spills |
| BO17 | 12 | The 65th spill in the `primary` scope leaves the newest 64 files | planned TestPrimarySpillsKeepNewest in internal/responsebound | An unpruned scope keeps 65 files |
| BO18 | 13 | A regular file at the scope directory path gives the complete output and then `spill-failed{reason=<reason>}` | planned TestOwnerCreateFailureKeepsOutput in internal/responsebound | A fallback that prints the projection loses the omitted lines |
| BO19 | 14 | A write fault after 100 spill bytes gives a file of 100 bytes and a response that holds the rest after the `spill-failed` line | planned TestOwnerMidwayFailureKeepsEveryByte in internal/responsebound, with an injected file writer | A fallback that drops the unwritten bytes fails the concatenation match |
| BO20 | 15 | `bench consumers` for the line value lists only the owner package and the policy registry | review-owned: consumer query at the BO-C1 review | A second package that restates 10 drifts from the owner |
| BO21 | 16 | A public registry entry without a bound disposition fails the registry test | planned TestEveryPublicCommandDeclaresBound in cmd/bench | A new command with no disposition skips the bound silently |
| BO22 | 17 | An empty response prints nothing and creates no spill directory | planned TestOwnerEmptyResponse in internal/responsebound | An owner that always creates the store leaves empty files |
| BO23 | 18 | An 11-line response whose last line has no newline spills and keeps the exact final bytes | planned TestOwnerCountsUnterminatedLine in internal/responsebound | A newline-count rule sees 10 lines and prints the whole response |
| BO24 | 19 | A response that holds a NUL byte and invalid UTF-8 spills with byte-exact file content | planned TestOwnerCopiesBinaryBytes in internal/responsebound | A rune or text conversion changes the bytes |
| BO25 | 20 | A Bench home that holds a newline gives the create-failure route | planned TestOwnerRefusesUnsafeSpillPath in internal/responsebound | A printed raw path forges a second response line |
| BO26 | 21 | `bench worktree exec <target> -- sh -c 'seq 1 40'` prints lines 1 to 4, the spill line, and lines 36 to 40 | planned TestExecBoundsChildOutput in internal/systemtest | An exec that passes the child's streams directly prints 40 lines |
| BO27 | 22 | An exec child that prints 40 lines and exits 7 makes exec exit 7 | planned TestExecKeepsChildExitUnderBound in internal/systemtest | An owner finish that replaces the child code fails the match |
| BO28 | 23 | A heredoc on exec stdin reaches the child, and the child's 12 lines of output spill | planned TestExecForwardsStdinUnderBound in internal/systemtest | A captured stdin leaves the child without its script |
| BO29 | 24 | A SIGINT to exec after the child printed 12 lines gives exit 130 and a spill file that holds those 12 lines | planned TestExecInterruptKeepsOutput in internal/systemtest | An owner that the interrupt path skips loses the printed lines |
| BO30 | 25 | After ticket 2, a nested `bench` child that prints 30 lines gives exec 10 lines with the child's own spill line | planned TestExecNestedBenchBoundsOnce in internal/systemtest | A child outside the bound passes 30 lines through exec |
| BO68 | 26 | After the spill starts, the owner holds only the head lines and a ring of the last 5 lines in memory | review-owned: code reading at the BO-C1 review, with a 64 MiB exec child of short lines in planned TestExecStreamsLargeChild in internal/systemtest | An owner that buffers the complete output grows with the child |
| BO69 | 64 | `sh -c 'sleep 30 & echo up'` under exec returns within the wait delay with the child's exit code and the line `up` | planned TestExecReturnsAtChildExit in internal/systemtest | An exec with no wait delay blocks until the background process exits |
| BO70 | 65 | `release-preflight` through `Command.Run` prints a 30-line response in full | planned TestShipTierStaysComplete in cmd/bench | A bounded ship-tier command spills its evidence on a discarded runner |
| BO72 | 11 | A `release` of the assignment that is its own scope keeps its over-bound spill in the `primary` scope | planned TestRetiringVerbSpillsToPrimary in internal/worktree | A retiring verb that spills to its own scope deletes the spill that it names |
| BO71 | 66 | A commit that exits 3 calls no build and prints `commit-chain{commit=<sha>,build=skipped,preflight=skipped}` at exit 3 | planned TestCommitChainStopsAtRemainder in cmd/bench | A chain that ignores exit 3 builds an unreconciled checkout |
| BO31 | 27 | An exec grammar refusal prints its `usage: bench worktree exec` line unchanged | existing TestWorktreeExecGrammar tests in internal/worktree, run unchanged | An owner that rewrites refusals breaks the documented exit-2 rule |
| BO32 | 28 | A list of 3 active rows prints `help[2]{cmd,why}:` with `bench worktree path <target>` and `bench worktree exec <target> -- <command>` and no active id in a help row | planned TestListActiveRowsUseTargetSlot in internal/worktree | Per-row help prints 6 rows and names each id |
| BO33 | 29 | A cleanup-pending row keeps its `bench worktree release --request <token> <path>` help row | existing TestListActions cleanup-pending case in internal/worktree/list_actions_test.go, run unchanged | A slot rule applied to every state loses the path operand |
| BO34 | 30 | A missing-tree row keeps its recovery help row | existing TestListActions missing-tree case in internal/worktree/list_actions_test.go, run unchanged | A slot rule applied to every state loses the recovery route |
| BO35 | 31 | A foreign row keeps its `bench worktree clean <path>` help row | existing TestListActions foreign case in internal/worktree/list_actions_test.go, run unchanged | A slot rule applied to every state loses the orphan path |
| BO36 | 32 | The id cell of an active row passes `bench worktree path` at exit 0 | rewritten TestListPathActionRunsAsAdvertised in internal/worktree/path_identifier_test.go | A cell that holds the label fails the resolver |
| BO37 | 33 | An all-green build preflight prints `checks{green=13,not_applicable=2,red=0}` and no `checks[` table | planned TestPreflightGreenSummaryLine in internal/preflight | A render that keeps the green rows prints a `checks[15]` table |
| BO38 | 34 | An all-green review preflight prints the same summary line and no `checks[` table | planned TestPreflightGreenSummaryLine review case in internal/preflight | A build-only change leaves the review table unchanged |
| BO39 | 35 | A preflight with one red check prints the summary line and `checks[1]{check,verdict,detail,next}` with the red row only | planned TestPreflightRedRowsOnly in internal/preflight | A summary that hides the red row loses the next action |
| BO40 | 35 | A red preflight exits 1 | existing TestReviewPreflight red cases in internal/preflight/command_review_test.go, updated for the summary line | A summary render that drops the red flag exits 0 |
| BO41 | 36 | The craft-cli `bench worktree list` table row states the target slot rule | planned anchor needle in internal/anchors registry, graded by the gate's anchor check | Guidance that drops the slot rule lets per-row help return |
| BO42 | 37 | `bench preflight evidence <id>` prints one `evidence_summary` block with no content cell | planned TestEvidenceDefaultPrintsSummary in internal/preflight/evidencecmd | The old default prints manifest bytes |
| BO43 | 38 | The `next` command of the summary returns the first manifest page byte-for-byte as the old default did | planned TestEvidenceSummaryNextReadsFirstPage in internal/preflight/evidencecmd | A wrong cursor starts the chain at a source page |
| BO44 | 39 | The cursor, source, verify, and check-current forms keep their current output | existing tests in internal/preflight/evidencecmd/evidence_command_test.go and evidence_modes_test.go, run unchanged | A shared render change alters a page response |
| BO45 | 40 | `--to <dir>` over 3 sources writes `source-1` to `source-3` with the source bytes and one `index.toon` | planned TestEvidenceExportWritesSources in internal/preflight/evidencecmd | A missing file or a changed byte fails the comparison |
| BO46 | 41 | `--to` at a directory that holds one file exits 1 and writes nothing | planned TestEvidenceExportRefusesNonEmptyDir in internal/preflight/evidencecmd | An export that writes beside the file overwrites data |
| BO47 | 42 | `--to` at a symlink exits 1 and writes nothing | planned TestEvidenceExportRefusesSymlink in internal/preflight/evidencecmd | A followed symlink writes outside the directory |
| BO48 | 43 | `--to` over a pack with one corrupt page exits 1 and leaves no directory it created | planned TestEvidenceExportVerifiesBeforeWrite in internal/preflight/evidencecmd | An export that writes before it verifies leaves bad bytes |
| BO49 | 44 | A source id `../escape` exports to `source-1` inside the directory | planned TestEvidenceExportNamesByOrdinal in internal/preflight/evidencecmd | A name taken from the id writes outside the directory |
| BO50 | 45 | A write fault on the second source removes `source-1` and the created directory | planned TestEvidenceExportCleansOnFailure in internal/preflight/evidencecmd, with an injected writer | A partial export stays on the disk |
| BO51 | 46 | A green commit, build, and preflight end with `commit-chain{commit=<sha>,build=green,preflight=green}` at exit 0 | planned TestCommitChainRunsThreeSteps in cmd/bench, with injected build and preflight seams | A chain that stops after the commit prints no build step |
| BO52 | 47 | A lane refusal exits 1, calls no build, and prints `commit-chain{commit=none,build=skipped,preflight=skipped}` | planned TestCommitChainSkipsAfterRefusal in cmd/bench | A chain that ignores the refusal builds an unpublished tree |
| BO53 | 48 | A red build exits with the build code, calls no preflight, and names the published commit | planned TestCommitChainStopsAtRedBuild in cmd/bench | A chain that runs the preflight grades a failed build |
| BO54 | 49 | A red preflight after a green build exits 1 | planned TestCommitChainReportsRedPreflight in cmd/bench | A chain that returns the commit code hides the red |
| BO55 | 50 | `--dry-run` with `--preflight-build` exits 2 with the commit usage line | planned TestCommitChainRefusesDryRun in internal/commit | An accepted combination builds during a dry run |
| BO56 | 51 | `bench commit` without the flag keeps its current output and exit codes | existing tests in internal/commit/commit_test.go and dry_run_test.go, run unchanged | A grammar change that alters the old form fails these tests |
| BO57 | 52 | A bounded verb run in an assignment worktree appends one tab-separated output line with the head, lines, bytes, and `spilled` | planned TestOutputRecordAppendsLine in internal/census | A dispatcher that records nothing leaves the file absent |
| BO58 | 53 | `bench worktree exec <target>` from the primary checkout records under the target assignment with the head `bench worktree exec` | planned TestExecOutputRecordUsesTarget in cmd/bench | A record keyed by the working tree writes no line |
| BO59 | 54 | `census.Counts` for an assignment is unchanged after 3 output records | planned TestOutputRecordsLeaveRawCount in internal/census | A shared file raises the raw-call count |
| BO60 | 55 | The landing prints `census output{bench worktree list=2/28978}` on stderr for two canned records | planned TestLandingPrintsOutputBreakdown in internal/worktree | A landing that reads only raw calls prints no output line |
| BO61 | 56 | `census.Drop` removes both the raw-call file and the output file | planned TestDropRemovesOutputRecord in internal/census | A drop of one file leaves the output record |
| BO62 | 57 | A symlink at the census directory leaves the verb's response and exit code unchanged | planned TestOutputRecordFailureKeepsVerdict in cmd/bench | A write error that reaches the exit code changes the verdict |
| BO63 | 58 | A response above the approved byte value spills, and its projection stays within the value | planned TestOwnerAppliesByteBound in internal/responsebound, after the entry stop clears | A line-only owner prints one oversized line |
| BO64 | 59 | The byte value in the policy registry equals the value in the reviewer budget decision record | review-owned: ticket-entry inspection of the decision record | A value chosen by the build ships an unreviewed budget |
| BO65 | 60 | The byte-bound ticket starts no product write while the decision record is absent | review-owned: ticket-entry inspection by the orchestrator | A placeholder value ships before the decision |
| BO66 | 61 | The diff leaves `.bench/hooks/block-bench-follow-on.sh` and the chain sentence of `.bench/BENCH.md` unchanged | review-owned: final reconciliation of the landed diff | An edit to the chain rule reopens closed decision 1 |
| BO67 | 62 | An ambiguous `bench consumers` name keeps one re-query help row per candidate | existing candidate tests in internal/consumers, run unchanged | A slot rule applied to consumers breaks its anchored decision |

Not covered: story 63 — the verb audit is spec-time evidence in `assets/verb-audit.md`, and no build ticket changes it.

### Edge inventory

The shell CLI hostile-input profile applies to the owner and to `--to`. The walk covers control bytes, symlinks, a file at a directory path, and a full disk. It also covers an interrupt, binary bytes, an empty output, and an unterminated last line. Each of these edges has a row above.

- **Won't handle:** the exact cross-stream order of an exec child's two streams — the owner records arrival order, and every byte survives.
- **Won't handle:** an interactive exec child with no stdin script — its output appears at exit, and the heredoc form stays the agent route.
- **Won't handle:** live progress of a long verb on a human terminal — the response appears at exit, and `bench worktree shell` stays interactive.
- **Won't handle:** spill growth inside one assignment — the assignment's retirement bounds it.
- **Won't handle:** an output record for a verb outside any assignment — no assignment key exists, and the verb audit covers those verbs.
- **Won't handle:** a single line longer than the byte value before ticket 10 lands — the line bound still holds, and ticket 10 closes it.
- **Won't handle:** the memory of a long unterminated line before ticket 10 lands — the owner never reaches line 11, and the byte bound closes it.
- **Won't handle:** a nested verb's output counted in both the child record and the exec record — each head names its own process.

## Ownership fences

- `internal/bounds/bounds.go`
- `tests/canary/package-core-guard/bounds-duplicate-owner`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `internal/responsebound/`
- `cmd/bench/`
- `cmd/bench/main.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/response_bound_test.go`
- `cmd/bench/commit_chain_test.go`
- `cmd/bench/census_output_test.go`
- `cmd/bench/worktree_leaves.go`
- `internal/worktree/exec.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/systemtest/`
- `internal/systemtest/exec_bound_test.go`
- `internal/racetests/racetests.go`
- `internal/worktree/lifecycle.go`
- `internal/worktree/land.go`
- `internal/worktree/list.go`
- `internal/worktree/list_actions_test.go`
- `internal/worktree/path_identifier_test.go`
- `internal/worktree/response_spill_test.go`
- `internal/worktree/land_census_output_test.go`
- `internal/preflight/command.go`
- `internal/preflight/command_review_test.go`
- `internal/preflight/source_tip_test.go`
- `internal/preflight/verdict_summary_test.go`
- `internal/preflight/evidencecmd/`
- `internal/chargeevidence/schema.go`
- `internal/chargeevidence/read.go`
- `internal/chargeevidence/export.go`
- `internal/commit/commit.go`
- `internal/commit/chain_grammar_test.go`
- `internal/census/census.go`
- `internal/census/output.go`
- `internal/census/output_test.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/axi_query_registry_test.go`
- `.agents/skills/bench-craft-cli/SKILL.md`
- `reviews/ft336-bounded-output.md`

## Ticket graph

| ticket | blocked by | chunk |
| --- | --- | --- |
| `1-bound-exec-output.md` | none | BO-C1 |
| `2-bound-every-public-response.md` | `1-bound-exec-output.md` | BO-C2 |
| `3-retire-response-spills.md` | `1-bound-exec-output.md` | BO-C2 |
| `4-slot-worktree-list-actions.md` | none | BO-C3 |
| `5-summarize-green-preflight.md` | none | BO-C3 |
| `6-summarize-evidence-default.md` | none | BO-C4 |
| `7-export-evidence-sources.md` | `6-summarize-evidence-default.md` | BO-C4 |
| `8-chain-commit-preflight.md` | none | BO-C5 |
| `9-record-response-census.md` | `2-bound-every-public-response.md` | BO-C6 |
| `10-apply-byte-bound.md` | `2-bound-every-public-response.md` | BO-C7 |

## Out of scope

- Selected worktree and history views: `specs/session-context-queries` owns them. Estimate: 0 edits here, 0 gate runs.
- Harness-level overflow replacement of a tool result: `specs/session-context-overflow` owns it. Estimate: 0 edits here, 0 gate runs.
- A per-verb `--full` flag on each query that lacks one: the response spill is the detail route. Estimate: 14 edits, 3 gate runs.
- Slot actions for `bench consumers` candidates, which reverses an anchored decision. Estimate: 4 edits, 1 gate run.

## Further notes

### Source trace

| source sentence | rows |
| --- | --- |
| A Bench response that is longer than 10 lines is a defect. | BO1, BO2, BO3 |
| Each verb prints what the agent acts on: the verdict and the next command. | BO3, BO37, BO39 |
| A green table becomes one summary line. | BO37, BO38 |
| Detail stays behind `--full` or in a file that the response names. | BO4, BO5 |
| Every payload has a byte bound and a line bound. | BO1, BO63 |
| Per-row help rows go away. | BO32, BO41 |
| Exec bounds the child output, or it writes the excess to a file and prints its path. | BO26, BO29 |
| The same bound applies to each Bench verb that exec runs. | BO30 |
| The spec audits every verb that an agent runs in a worktree. | story 63, Not covered |
| The census adds the output bytes of each verb. | BO57, BO58, BO60 |
| `bench preflight evidence` prints a small summary by default. | BO42, BO43 |
| Each source stays behind an explicit request or a `--to <dir>` spill. | BO44, BO45 |
| A compound verb commits, runs the worktree build, and runs the build preflight in one call. | BO51 |
| The reviewer decides whether a hook and platform rule allows safe chains of Bench calls. | BO66 |
| Closed decision 4: FT336 waits for the measurement. | BO63, BO64, BO65 |

### Reader sweep and proof checklist

Readers of the public response bytes:

- The hook scripts in `.bench/hooks/` call only internal plumbing commands: `worktree-hook`, `guard-file-write`, `guard-bench-follow-on`, `session-inspect`, `guard-git`, `stop-verdict`, and `check-agent-line`. The bound excludes them.
- `.bench/gate.sh` calls the internal `freshness-check` and `gate-phases` commands. The bound excludes them.
- `internal/gate/lane_select.go` reads only the exit code of `bench test --check`.
- The tests that read public output through `Command.Run` sit in `cmd/bench`. The tests that run the built binary sit in `internal/systemtest`. Both prefixes are in the fence for the posture change.
- `internal/worktree/path_identifier_test.go` extracts an id from a help row. Ticket 4 rewrites it to read the id cell.
- `internal/conformance/axi_query_registry_test.go` requires the phrase `per matching row`. The general AXI sentence keeps that phrase, and only the `bench worktree list` table row changes.
- `internal/preflight/evidencecmd/evidence.go` sets the `next` of a `--charge` response to the bare manifest-first command. After ticket 6 that command returns the summary, which adds one request for each charge read. Closed decision 2 accepts that request.
- `.github/workflows/release.yml` and `.github/workflows/native-runtime.yml` run the ship-tier commands. The ship-tier exemption keeps their output complete.
- `.agents/commands/bench-review-implementation.md` tells a consumer to follow `bench preflight evidence <id>` and each successor command. The summary prints that successor, so the text stays true.

Readers of the census directory: `census.Counts`, `census.HeadBreakdown`, and `census.ReadEvents` read the `<id>` file only. `internal/assessment/collection.go` calls `census.ReadEvents`. `internal/status/status.go` calls `census.Counts`. `census.Drop` has one call site in `internal/worktree/lifecycle.go`, and `TestRequireOneCallSite` in `internal/worktree/worktree_test.go` pins it.

Proof checklist:

- Cited symbols: each symbol below resolves at `a519ed97`.
  - In `cmd/bench`: `Command.Run`, `renderCommandHelp`, `commandRegistry`, and `outputCommand`.
  - In `internal/worktree`: `ExecCommand`, `runWorktreeChild`, `nameWorktree`, `ListCommand`, `actionsForRows`, `recoverMissingTree`, `BuildCommand`, and `printCensusHeads`.
  - In the preflight packages: `verdictCommand`, `evidencecmd.Read`, `Artifact.First`, `Fragment.Encode`, and `ResponseLimit`.
  - In other packages: `commit.Command`, `census.Record`, `census.Counts`, `census.HeadBreakdown`, `census.ReadEvents`, and `census.Drop`.
- Import edges: `cmd/bench` already imports `internal/census`, `internal/worktree`, and `internal/preflight`. The new owner package imports only `internal/bounds`, `internal/benchhome`, `internal/poolkey`, and `internal/sanitize`.
- Source-row clauses and occurrences: the source trace table above.
- Promised field labels: `spilled{lines,bytes,omitted_lines,path}`, `spill-failed{reason}`, `spill-failed{path,written_bytes,reason}`, `checks{green,not_applicable,red}`, `evidence_summary`, `exported{sources,bytes,dir}`, `commit-chain{commit,build,preflight}`, and `census output{...}`.
- Changed-function callers: `Command.Run` has one production caller, `main`. `verdictCommand` serves the build and review modes. `evidencecmd.Read` has one caller, the evidence mode dispatch. `actionsForRows` has one caller, `ListCommand`.
- Copy survival: BO20.

Sources re-read in this session:

- `roadmap/FT336.md` and ADR 0022
- `internal/worktree/exec.go` and `internal/worktree/list.go`
- `internal/census/census.go` and `internal/census/events.go`
- `internal/commit/commit.go` and `internal/preflight/command.go`
- `internal/preflight/evidencecmd/evidence.go`, `internal/preflight/evidencecmd/bound.go`, and `internal/chargeevidence/read.go`
- `internal/axi/action.go`, `cmd/bench/command_registry.go`, and `cmd/bench/main.go`
- `.agents/skills/bench-craft-cli/SKILL.md` and the fences of the two session-context specs

Not re-read: the drained `capture/IDEAS.md` source text, which the drain removed.

### Fence disposition

Reviewer disposition of the ownership fences: open.

The `cmd/bench/` and `internal/systemtest/` prefixes carry the posture change of ticket 2. `cmd/bench/command_registry.go` and `cmd/bench/main.go` are also in the fences of both staged session-context specs. The three specs land serially, and each landing merges the earlier one.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"BO-C1","tickets":["1-bound-exec-output.md"],"verification":[{"id":"owner","command":"bench test --package ./internal/responsebound"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"}]},{"id":"BO-C2","tickets":["2-bound-every-public-response.md","3-retire-response-spills.md"],"verification":[{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"owner","command":"bench test --package ./internal/responsebound"},{"id":"system","command":"bench test --check system"}]},{"id":"BO-C3","tickets":["4-slot-worktree-list-actions.md","5-summarize-green-preflight.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"consumers","command":"bench test --package ./internal/consumers"}]},{"id":"BO-C4","tickets":["6-summarize-evidence-default.md","7-export-evidence-sources.md"],"verification":[{"id":"evidencecmd","command":"bench test --package ./internal/preflight/evidencecmd"},{"id":"chargeevidence","command":"bench test --package ./internal/chargeevidence"}]},{"id":"BO-C5","tickets":["8-chain-commit-preflight.md"],"verification":[{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"commit","command":"bench test --package ./internal/commit"}]},{"id":"BO-C6","tickets":["9-record-response-census.md"],"verification":[{"id":"census","command":"bench test --package ./internal/census"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"cmd","command":"bench test --package ./cmd/bench"}]},{"id":"BO-C7","tickets":["10-apply-byte-bound.md"],"verification":[{"id":"owner","command":"bench test --package ./internal/responsebound"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/ft336-bounded-output/spec.md"},{"id":"owner","command":"bench test --package ./internal/responsebound"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"system","command":"bench test --check system"}]}
```

### Flagged additions

- The help-form, dashboard, shell, and setup exemptions are author choices. The read-only consultation found that a bounded `bench help` is a dead key for the inventory route.
- The spill store, its 64-file limit, and its retirement drop are author choices. No source sentence names a location.
- The `spill-failed` fallback posture comes from the overflow spec's preservation rule.
- `--to` refuses a non-empty directory and names files by ordinal. The source names only the flag.
- The stdout and stderr combination under the bound is an author choice, because the audit measured both streams together.
- The ship-tier and `repair` exemptions, the exec wait delay, and the `primary` scope for a retiring verb came from the review round.

### Flags for reviewer veto

- The ticket parser accepts only sibling basenames in `Blocked by:`. Ticket 10 therefore carries its cross-spec wait as an entry stop in `What to build` and `Acceptance`, as queries ticket 4 does. The orchestrator, not the parser, honors that stop.
- The unbuilt overflow ticket 2 plans a second store for preserved output. The reviewer can direct it to consume this spec's response owner, so that one store holds preserved output.
- `bench consumers` keeps one re-query action per candidate row, because an anchor pins that decision.
- The general AXI sentence "Derive one state-derived action per matching row" stays. Only the `bench worktree list` row states the slot rule.
- The session-context-queries spec pins the bare `bench worktree list` view with QU9. If that spec lands first, ticket 4 adds the QU9 baseline test to its `Writes:` line under the in-scope plan expansion rule.
