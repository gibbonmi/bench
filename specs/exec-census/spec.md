# The exec census

Status: staged

Decision source: ready compiled map `specs/exec-census/decisions/exec-census.md`, resolved 2026-08-25

Verification log: 2 iterations to accept — nine blocking findings in round one, all unread tree facts. The second iteration folded them, and a scoped re-review accepted.

## Problem

The agent runs multi-step edits against the cached absolute path of a Bench
worktree instead of through `bench worktree exec`. Nothing in the kit counts
those calls. The FT248 fold on 2026-08-25 ran about fifteen such calls, and
the only record is the agent's own recall. A final-check duty that asks the
agent whether it took the verb reads that recall, and recall is not
evidence. The reviewer wants Bench comfortable enough that the agent prefers
the verb, and wants the drift measured before and after the verbs improve.

## Solution

The follow-on guard already reads every Bash tool call's command text. It
gains one duty: when the text names a path under the repository's worktree
pool and no command in the call invokes `bench`, it appends one census record
for that assignment. The record holds the time and the verb head, never the
command text.

`bench status` shows one `census` row with the total across active
assignments, and `--all` expands it per assignment. `bench worktree land`
carries the final count in its landed record as `census=<n>`. A release, a
clean, or the landing's own release drops the assignment's records. The
final-check command reads `census=<n>` and turns a nonzero count into one
`bench learning --rule` entry. A delegate charge asks the delegate for zero
to two CLI improvements, and the coordinator folds them into that entry.

## User stories

### A raw call is recorded

Line: opus / high.

The record rides a guard's plumbing verb and a retirement path owns the
drop, so the cached routing for hook and cleanup-authority work applies.
A ticket re-runs the decision table at charge time.

1. As an agent, I want a Bash call with a pool path and a non-`bench` verb head counted once, so that my habit is measured.
2. As an agent, I want a chain of commands in one call recorded once, so that one decision counts as one raw call.
3. As an agent, I want a call whose only pool path sits inside a `bench worktree exec` command recorded zero times, so that the verb is free.
4. As an agent, I want the wrapper invoked by its absolute path treated as `bench`, so that a wrapper path is not a raw call.
5. As an agent, I want a pool path found only in a heredoc body counted once with the first command's head, so that scripts count.
6. As a reviewer, I want each record to hold the time and the verb head, never the command text, so that no command log exists.
7. As an agent, I want a `git` verb head to include its subcommand (`git add`), so that the learning names the missing verb.
8. As an agent, I want `env`, `timeout`, `xargs`, and leading assignments stepped over, so that the verb head names the real command.
9. As an agent, I want the record write to leave the guard's verdict and exit code unchanged, so that the census never blocks a call.
10. As an agent, I want a failed record write to allow the call with no diagnostic, so that a census defect never stops work.
11. As an agent, I want a `bash -c` string that invokes `bench worktree exec` treated as a `bench` call, so that a wrapped verb is free.
12. As an agent, I want a call that the destructive-git guard blocks still recorded, so that the attempt counts.
13. As an agent, I want a Read, Edit, or Write tool call never recorded, so that a file edit at the path stays right.
14. As an agent, I want a call that names a dead assignment's pool path recorded under that assignment's id, so that no ledger read occurs.
15. As an agent, I want two concurrent delegates' records appended without loss, so that a parallel build counts every call.
16. As an agent, I want the record keyed by the assignment id from the path's last segment, so that the writer never parses the ledger.

### The census is visible

Line: opus / high.

The signal, the landed key, and the drop are signal and cleanup-authority
work, so the cached high routing for that work applies.

17. As a reviewer, I want one `census` row with `<n> raw call(s) across <k> worktree(s)` on the board, so that the count is ambient.
18. As a reviewer, I want no `census` row for a zero count or for an inactive assignment, so that the board shows only signals.
19. As a reviewer, I want the `census` row at severity 3 with no action, so that it ranks beside `guards` and names no false remedy.
20. As a reviewer, I want `bench worktree land` to print `census=<n>` as the last key of its landed record, so that final-check reads evidence.
21. As a reviewer, I want a landing of an assignment with no records to print `census=0`, so that zero is a stated fact.
22. As a reviewer, I want an incomplete landing to keep the records and print the count again on resume, so that no evidence is lost.
23. As a reviewer, I want `bench worktree release` to drop the assignment's records, so that a released assignment leaves no census behind.
24. As a reviewer, I want `bench worktree clean` and the landing's release step to drop records through the same function, so that one owner drops.
25. As a reviewer, I want a control byte in an assignment label sanitized before the board renders it, so that no label splits a row.
26. As a reviewer, I want the census files under `$BENCH_HOME/census/<repo-key>/`, a sibling of the pool, so that `bench worktree reclaim` never reads them.
27. As a reviewer, I want `bench status --all` to expand the `census` row per assignment as `<label> <n> raw call(s)`, so that each worktree shows.

### The duty and the charge

Line: opus / high.

The final-check command, the delegate skill, and the reference docs are
guidance prose, so the leverage override routes them mid and high.

28. As a reviewer, I want final-check to write one `bench learning --rule` entry per landing with `n > 0`, so that findings reach the drain.
29. As a reviewer, I want the entry's title, `--what`, `--right`, and `--rule` fields in the fixed shape below, so that the drain can act.
30. As a reviewer, I want the close to state `census: 0 raw calls` when `n = 0`, so that a zero is reported, not skipped.
31. As a reviewer, I want the duty advisory, so that a nonzero count never blocks or reds a landing.
32. As a reviewer, I want the delegate charge to ask for zero to two CLI improvements, so that delegate friction reaches the drain.
33. As a reviewer, I want the anchors registry to require the census duty's marker in final-check, so that an edit cannot drop the duty silently.
34. As a reviewer, I want the reference, the profile, and the glossary to describe the `census` signal, so that a cold session reads one account.
35. As a reviewer, I want a spec retro to cite the landing's census entry with its `Feeds:` line, so that retro and drain agree.

## Implementation decisions

**The record rides the follow-on guard's plumbing verb.** `bench
guard-bench-follow-on` reads the envelope's command text today and returns
the verdict. It gains one call before the verdict: record the raw call. This
choice adds no hook script, so both harness wirings, the guards manifest, the
conformance registry, and the adopt payload stay unchanged. The record call
never changes the verdict, and its own failure is silent. The verb tests the
text for the `<home>/worktrees/` substring before it resolves any root, so an
ordinary call spawns no git process.

**One new package owns the census.** `internal/census` is the one new seam.
It owns four behaviors:

- the match: the command text names a path under the pool prefix, and no simple command invokes `bench`
- the verb head: the head of the first simple command whose words name a pool path, or the first command's head when only a heredoc body names it
- the record: one appended line per raw call, keyed by the assignment id
- the read and the drop: counts per assignment for a repository, and removal of one assignment's records

The guard verb, `bench status`, and the worktree verbs compose these. The
package reuses `shellcommand.Parse` and `ResolveRoutinePrefix`, and it
reuses the guard's wrapper-aware `bench` test, so a `bash -c` string that
invokes `bench` is a `bench` call. The verb head and the `bench` test have
one source.

**The pool key and the canonical root move to one low package.** The
worktree package derives the pool key from a canonical root today, and the
canonical root from the git common directory. Both derivations move to
`internal/poolkey`, because `internal/worktree` needs the census drop while
`internal/census` needs the key, and a direct pair of imports cycles. The
pool path stays byte-identical. A caller passes any root; the key derivation
canonicalizes it first, so a call from inside a linked worktree keys the
primary repository.

**The records live beside the pool, not inside it.** The files sit at
`$BENCH_HOME/census/<repo-key>/<assignment-id>`, one line per record, opened
for append. The map's discretion note placed them under the pool directory.
They sit beside it instead, because `bench worktree reclaim` reads the pool
parent and a foreign file there changes a plan.

An absent file and an empty file both read as zero. The record holds the time and the verb head. The
assignment label lives in the ledger and joins at render time, so the map's
label field collapses into the file's key.

**The signal.** `bench status` gains `appendCensus`. It sums the counts of
the active assignments and adds one row at severity 3, signal `census`, and
no action. The detail is `<n> raw call(s) across <k> worktree(s)` through
the board's own `Plural`.

`--all` expands the row to one row per assignment,
detail `<label> <n> raw call(s)`. The expansion is a sibling of the `intent`
expander in the same render branch, so the five-row budget holds. The label passes through
`sanitize.Controls`.

**The landed record.** Both landed forms gain `census=<n>` as the last key.
The landing reads the count before its release step. The release step drops
the records, so a resume of an incomplete landing still reads and prints the
count.

**One drop owner.** One function in `internal/census` removes an
assignment's file. The one shared retirement path in `internal/worktree`
calls it, and `release`, `clean`, and the landing's release all reach that
path. The landing already holds the home; the release and clean verbs thread
theirs to that path.

**The duty.** The final-check command's post-merge tail gains the census
duty, the advisory rule, the zero close, and the retro citation. The
learning entry has one fixed shape:

- the title names the assignment label and `n`
- `--what` lists each verb head with its count
- `--right` names the Bench form per head, or `none`
- `--rule` proposes the verb or the help change

The coordinator folds a delegate's zero to two CLI improvements into that
same entry. The `.claude/commands/` path is a symlink to `.agents/commands/`,
so one edit serves both harnesses. The anchors registry gains one `Require`
needle each for the duty, the zero close, the charge line, the reference
sentence, and the profile sentence. The delegate skill sits at its prose
budget, so the charge line absorbs into an existing sentence; the budget
table is reviewer-owned. `CONTEXT.md`'s `**signal**` enumeration gains
`census`, because the signal-vocabulary check reads that enumeration.

## Testing decisions

The highest seam that shows each recording failure is the `internal/census`
unit: a command text, a root, and a temporary home in, a record out. The
guard verb's journey proves the composition with a JSON envelope. The signal
is tested through `appendCensus` with a ledger fixture and a temporary home.
The landed key and the drops are tested through the worktree package's
landing and release journeys. The anchors registry test proves each needle
bites, and the signal-vocabulary check proves the enumeration. The gate's
`test` phase observes all of it.

### Seam diagram

    trigger: a Bash tool call (PreToolUse, both harnesses)
        │
        ▼
    envelope ──▶ [ bench guard-bench-follow-on ] ──▶ verdict (unchanged)
                        │
                        ▼
              [ census.Record(command, root, home, now) ] ──▶ one line in $BENCH_HOME/census/<key>/<id>
                        ◀ tests attach here: command text + temp home in, file out

    trigger: bench status / bench worktree land / release / clean
        │
        ▼
    ledger + census files ──▶ [ census.Counts / census.Drop ] ──▶ census row / census=<n> / no file
                                   ◀ tests attach here: ledger fixture + temp home

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| EC01 | 1 | A command text `sed -i s/a/b/ <pool>/<id>/x` with a temp home appends one record under `<id>`. | census unit | An omitted match writes nothing. |
| EC02 | 2 | A text `cat <pool>/<id>/a && sed -i x <pool>/<id>/b` appends one record. | census unit | A per-command count writes two. |
| EC03 | 3 | A text `bench worktree exec a -- sed -i x <pool>/<id>/b` appends no record. | census unit | A text-only match counts the verb. |
| EC04 | 4 | A text `/tmp/r/bin/bench.sh status <pool>/<id>` appends no record. | census unit with a stub resolver | A `bench`-word-only test counts the wrapper path. |
| EC05 | 5 | A text `python3 - <<EOF\nopen("<pool>/<id>/x")\nEOF` appends one record with head `python3`. | census unit | A token-only match misses the heredoc body. |
| EC06 | 6 | A record line holds a time and a verb head and does not contain the command text. | census unit | A record that stores the text is a command log. |
| EC07 | 7 | A text `git add <pool>/<id>/x` records head `git add`. | census unit | A bare `git` head hides the missing verb. |
| EC08 | 8 | A text `FOO=1 timeout 5 sed -i x <pool>/<id>/y` records head `sed`. | census unit | An unresolved prefix records `timeout` or `FOO=1`. |
| EC09 | 9 | The guard verb given a raw-call envelope with a plain command exits 0 and writes the record. | guard verb journey | A verdict changed by the record is a new block. |
| EC10 | 10 | The guard verb with an unwritable home exits 0 with empty stderr. | guard verb journey | A record error that surfaces blocks or warns. |
| EC11 | 11 | A text `bash -c 'bench worktree exec a -- sed -i x <pool>/<id>/b'` appends no record. | census unit | A top-level-only `bench` test counts the wrapped verb. |
| EC14 | 14 | A text that names `<pool>/<unknown-id>/x` appends one record under `<unknown-id>`. | census unit | A writer that needs the ledger drops a dead path. |
| EC15 | 15 | Two goroutines that record to one assignment leave two intact lines. | census unit under `-race` | A rewrite loses a line. |
| EC16 | 16 | A text naming `<pool>/<owner>-<id>/x` records under `<id>`, not `<owner>-<id>`. | census unit | A key on the whole segment never matches the ledger. |
| EC17 | 17 | Two active assignments with two and one records make `appendCensus` return one row `census` with detail `3 raw calls across 2 worktrees`. | status unit with a ledger fixture | An omitted signal returns no row. |
| EC18 | 18 | An active assignment with no file and a released assignment with two records each produce no row. | status unit with a ledger fixture | A count without the ledger shows a dead assignment. |
| EC19 | 19 | The `census` row has severity 3 and renders no action. | status unit | A wrong severity leads the board. |
| EC20 | 20 | A landing of an assignment with three records prints `census=3` as the last key of `landed{...}`. | worktree landing journey | An omitted key leaves final-check with recall. |
| EC21 | 21 | A landing of an assignment with no file prints `census=0`. | worktree landing journey | An absent file read as an error fails the landing. |
| EC22 | 22 | An incomplete landing keeps the file, and its resume prints `census=<n>` again. | worktree landing journey | A drop before the record loses the count on resume. |
| EC23 | 23 | `bench worktree release` of an assignment with records leaves no file for it. | worktree release journey | A release that keeps the file shows a stale row. |
| EC24 | 24 | `bench worktree clean` of an assignment and a complete landing each leave no file. | worktree journeys | A drop on one path only leaves a stale file on the other. |
| EC25 | 25 | A label with an ESC byte renders in the expanded `census` row with the byte replaced. | status unit | A raw byte splits the fixed-width row. |
| EC26 | 26 | `census.Dir(home, root)` returns a path under `<home>/census/` and never under `<home>/worktrees/`. | census unit | A directory inside the pool changes a reclaim plan. |
| EC27 | 28, 29, 35 | The final-check command carries the duty needle, the four learning fields, the once-per-landing rule, and the retro citation. | anchors registry test | A dropped marker passes a prose-only check. |
| EC29 | 30, 31 | The final-check command carries the zero-close needle with the advisory rule. | anchors registry test | A dropped sentence removes the zero report. |
| EC31 | 32 | The delegate skill's charge section carries the charge needle, and the file stays within its prose budget. | anchors registry test and prose budget check | A charge without the line yields no delegate finding. |
| EC32 | 33 | A copy of the final-check command without the duty needle makes the registry report its diagnostic. | anchors registry test | A needle that never reds proves nothing. |
| EC33 | 34 | The `**signal**` enumeration in `CONTEXT.md` names `census`, and the reference and the profile each carry the census sentence needle. | signal-vocabulary check and anchors registry test | A missing name reds the vocabulary check; a missing sentence reds the needle. |
| EC34 | 27 | `bench status --all` with two active assignments renders two `census` rows, each `<label> <n> raw call(s)`. | status render test | A board that only sums hides which worktree. |

Not covered: story 12 — the destructive-git guard runs as a separate hook process, so the record never depends on its verdict.

Not covered: story 13 — the `PreToolUse` matcher `Bash` is pinned by the existing conformance registry, and the census verb reads only that envelope.

### Edge inventory

- A pool path with a space or a glob character inside quotes folds into one word and matches.
- A pool path that is a prefix of a longer sibling directory (`<pool>x/`) does not match; the match requires the separator after the pool.
- A relative path into the pool never matches; the verb head test reads absolute text only.
- A `cd <pool>/<id> && make` records head `cd`.
- A `bash -c '...'` string that invokes `bench` anywhere makes the call a `bench` call; the wrapper scan goes one level deep, as the guard's does.
- An envelope with no command field takes the guard's existing warning path and records nothing.
- A command with a NUL byte takes the guard's existing refusal path and records nothing.
- A JSON-escaped separator (`&&`) decodes before the match; the standard decoder owns it.
- An unterminated quote in the text takes the tokenizer's fallback split and still matches on the raw text.
- A path whose last segment lacks the owner-and-id shape records nothing.
- A census directory that is a symlink is not followed; the write fails silently.
- A census file that is a FIFO or a device is refused before the read; the count reads as zero.
- A record line without a trailing newline still counts as one record.
- A home the process cannot create records nothing and never blocks.
- A landing that refuses before its gate prints no landed record and drops nothing.
- A `bench status` run inside a linked worktree keys the primary repository and shows the same rows.
- Three active assignments with records still take one board row; `--all` shows the three.
- A text that names two assignment ids records once, under the first id in the text.

**Won't handle** a pool path of another repository — the prefix is this repository's key, and a cross-repository call stays a habit for review.

**Won't handle** a `PostToolUse` outcome record — neither harness wiring carries the event, and the attempt is the fact the census counts.

**Won't handle** the text of a Read, Edit, or Write tool call — the reviewer named the path as the right operand for those tools.

## Ownership fences

- `internal/census/`
- `internal/poolkey/`
- `internal/benchguard/`
- `internal/shellcommand/`
- `internal/racetests/racetests.go`
- `cmd/bench/`
- `internal/status/`
- `internal/worktree/`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `.agents/commands/bench-final-check.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.bench/BENCH-reference.md`
- `projects/benchkit.md`
- `CONTEXT.md`
- `CHANGELOG.md`
- `specs/exec-census/`
- `internal/gitguard/scan.go`
- `internal/systemtest/owner_land_race_test.go`
- `tests/canary/package-core-guard/reintroduced-bare-skip/`

The last three fences were added during the build. The `git` subcommand
finder has one source in the destructive-git guard, so the verb head
reads it there. The `census=<n>` key changes every pinned landed record,
including the one in the system test. The bare-skip canary fixture anchors
on the `cksum` skip, which moved to `internal/poolkey` with the key.

The recording tickets land first on one integration source. The visibility
tickets follow them, and the guidance ticket lands last so its account names
the shipped behavior.

## Out of scope

- The exec-comfort verbs from the same map (stdin help line, `bench worktree resolve`, `bench worktree show`, path on refusal): 14 edits, 3 gate runs. It is the map's second spec.
- A pool-wide match across every repository under `$BENCH_HOME/worktrees`: 2 edits, 1 gate run.
- A `PostToolUse` outcome record in both harness wirings: 6 edits, 2 gate runs.
- The deterministic-step invariant in `.bench/BENCH.md`: 2 edits, 1 gate run, a drain item.
- A raised prose budget for the delegate skill: 1 edit, 1 gate run, a reviewer decision.
- Any change to the destructive-git guard's verdicts.

## Further notes

The compiled map feeds two specs. The exec-comfort spec reads
`specs/exec-census/decisions/exec-census.md` as its decision source. This
folder's retirement waits for that spec's landing, or the second spec copies
the map under its own folder at authoring time.

The reviewer batch-approved the map's rounds 2 to 4. Tickets #9, #11, and
#14 of the map stay open for veto; only #9 touches this spec, through
stories 12, 15, and 32. Two spec choices depart from the map and stand for
veto: the one aggregated board row with `--all` expansion (map #7 said one
row per assignment), and the record's label collapse into the file key (map
#4 named the label as a record field).
