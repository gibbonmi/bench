## Outcome

FT366 landed at `905c77f5` on a green landing gate. A PreToolUse guard, `bench guard-rg` with the shim `block-unsafe-rg.sh`, runs in Claude Code and Codex. It refuses a pathless `rg` call and an `rg` call with a live command substitution. The shared `shellcommand` lexer marks each live substitution once, and the follow-on refusal names command substitution. A pipe into a subshell group feeds its first `rg` call, and the `AGENTS.md` rg path bullet retired.

The build was a `--full --delegate --auto-approve` run with `--reviewer opus medium` and unified review. C1 used one repair cycle, C2 used both, and C3 and C4 used none. The C2 review added story 59 and row R59 under `--auto-approve`, and that row waits for reviewer veto.

## Cost

- landing: 905c77f51cfbab82d0ffb6d7557cd767273917db, range 905c77f51cfbab82d0ffb6d7557cd767273917db^1..905c77f51cfbab82d0ffb6d7557cd767273917db^2
- wall: 144 min
- commits: 39 total, 31 bookkeeping, 8 change
- graph: 5 tickets, chain 5, 4 chunks, 4 serial checkpoints
- dispatch: 5 authors, 3 repairs
- trace: commit 47 spans in 403998 ms, lane 0 spans in 0 ms, landing 1 spans in 13386 ms
- wall per ticket: 28.9 min against 60 min
- bookkeeping per ticket: 6.2 against 7.1
- delegates: about 1.39M subagent tokens over 17 sessions. The ticket 02 author used 246k, and the seven unified review rounds used 393k together.
- provider usage and charges: unknown.

## Gate-stage timings

- landing: commit 905c77f51cfbab82d0ffb6d7557cd767273917db, trace d2b3a055fc30dcddb01d54fa16b6c324
- The landing reused the green of the complete checkpoint.
- gofmt: 185 ms
- vet: unknown
- test: unknown
- race: unknown
- system: 51243 ms
- shellcheck: 537 ms
- The chunk checkpoints ran the system phase in 49254 ms, 48902 ms, and 50453 ms.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| C1 ticket 01 | R01 to R06 mark table green | verified | 9 | held | opus / xhigh / implementer |
| C2 ticket 02 | R38 no second option list | claimed | 8 | held | opus / xhigh / implementer |
| C2 repair 1 | R59 lead test is complete | verified | 9 | refuted | opus / xhigh / repair |
| C2 repair 2 | P2 no-space group lead closed | verified | 9 | held | opus / xhigh / repair |
| C3 ticket 03 | R41 to R45 red then green | verified | 9 | held | opus / medium / implementer |
| C3 ticket 04 | R39, R40, R47 red then green | verified | 9 | held | opus / medium / implementer |
| C4 ticket 05 | R51 to R54 red then green | verified | 9 | held | opus / medium / implementer |
| C1 review 1 | C1 opener reset has no test | claimed | 7 | held | opus / medium / reviewer |
| C2 review 1 | P1 two R16 inputs cannot fail | claimed | 9 | held | opus / medium / reviewer |
| C2 review 1 | C3 piped subshell group refused | claimed | 8 | held | opus / medium / reviewer |
| C2 confirm 1 | P2 no-space group refused | claimed | 8 | held | opus / medium / reviewer |
| C2 confirm 1 | S2 group-lead term is dead | claimed | 5 | held | opus / medium / reviewer |
| C3 review 1 | S1 shim runner copied | claimed | 6 | abstained | opus / medium / reviewer |
| C3 review 1 | P1 pointer bullet needed | claimed | 6 | abstained | opus / medium / reviewer |
| build | two repair rounds over the build | claimed | 6 | refuted | opus / medium / orchestrator |

Brier mean 0.130 over 13 pairs; 2 abstained.

## Coordinator catches

- The C1 `guard-preservation` plan probe could not bite inside its own three packages, because the R06 witness is in `internal/shellcommand`. The author added that package, and a learning holds the plan rule.
- One coordinator probe on each returned tree bit at a site that no author probed. The sites were the unquoted backtick mark, the `-e` positional rule, the group feed, the live Codex wiring, and the substitution repair sentence.
- The C2 reviewer proposed `ask-user` for the refused `cat f | (rg foo)`. The orchestrator decided it as a defect against the decision source, added R59, and flagged it.
- A concurrent spec-stage landing released the spec worktree while this session read it. The build moved to a new integration worktree at the landed commit.
- The orchestrator did not read this scorecard before the line declaration, which repeats the FT406 catch.
- Six orchestrator record commits failed their lane once on a sentence over 25 words or a paragraph over six sentences.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-mark-live-substitution.md | 1 | delegate-error |
| 02-classify-rg-calls.md | 2 | edge-found, edge-found |
| 03-add-guard-rg-verb.md | 0 | none |
| 04-wire-rg-guard-hook.md | 0 | none |
| 05-name-substitution-in-follow-on.md | 0 | none |

## Agent-experience improvements

### Bench CLI

- Let `bench test --in` and `bench probe --in` rebuild a stale worktree build once before the run, because each Go-editing author ran `bench worktree build` by hand up to three times.
  Feeds: new
- Add a coverage verb that rewrites a planned seam cell into its citation, because four authors rewrote the cells with `sed` or `python3`.
  Feeds: new
- Print the failing line of each mutant in a `bench probe --mutants` answer, because two authors reran each mutant alone to see which row failed.
  Feeds: new
- Keep a nested spill file until its parent response file goes, because the worktree release removed the complete checkpoint's phase table, so vet, test, and race read unknown.
  Feeds: new
- Census: 40 raw calls. The 29 `rg` calls were sanctioned reads. The `cat` and `python3` calls wrote mutant files and rewrote coverage cells, an owned gap that the two items above remove. The 2 `grep` calls, the `cd` call, and the `ls` call have no known purpose.
  Feeds: none

### Skills

- Let the staleness pass omit its delegates when the reviewed graph commit is the source tip and the audit drift list is empty, because three slices spent 155k tokens to return current.
  Feeds: new
- Require in each author charge that a silent witness triggers a check of every other input of its row, because the R16 author fixed one input and a repair cycle fixed the other two.
  Feeds: new
- Require in `/bench-write-spec` that each plan probe's witness row has its seam inside the requirement command.
  Feeds: new
- Require in a guard spec review a walk of each operator that can come before a span, including a subshell paren after a pipe, because C2 took two repair cycles on that one edge.
  Feeds: new
- Move the fail-open runner of the Bash guard shims into one shared library, because the new shim is the third hand copy and all three fail closed when `mktemp` fails.
  Feeds: new

### Process

- Write a ticket's witness reds in the same plan commit as the next ticket's assignment, because C3 saved one commit that way.
  Feeds: none
- Generate the per-round review pickup prose from the record entries, because the orchestrator restated record facts by hand and failed the lane six times.
  Feeds: new
- Read the scorecard before the line declaration, because two builds in a row skipped it.
  Feeds: none