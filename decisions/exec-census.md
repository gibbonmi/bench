# Exec census and exec comfort

Status: ready

## Destination

Bench becomes comfortable enough that the agent prefers `bench worktree exec`
over the cached absolute worktree path. Two specs deliver this, in order.
The **census spec** records every raw call into a Bench worktree path. It
shows the count on the ambient dashboard and folds it into the landing
record. A final-check exit duty turns a nonzero count into one
`bench learning --rule` entry. The **exec-comfort spec** closes the four verb
gaps the FT248 fold exposed:

- a stdin example in the exec help
- a `bench worktree resolve` verb
- a `bench worktree show` verb
- the worktree path on every exec refusal

The census lands first so the second spec has a measured baseline: the
FT248-shaped task counts raw calls before, and zero after.

## Provenance

The reviewer gave six answers on 2026-08-25 before this map existed:

- scope: every landing, as a final-check exit duty
- mechanism: a hook-recorded census read from `bench status`, never model recall
- delegates: no command log; a done-claim recommends zero to two CLI improvements
- disposition: `bench learning --rule`; the drain decides whether a CLI change opens
- sizing: shape first, then spec
- goal: the agent prefers the verb over habit

Round 1 (#1 to #4) closed by reviewer answer on 2026-08-25. The reviewer
then gave a batch approval: every later recommendation stands as the
answer. Rounds 2 to 4 (#5 to #15) therefore record the author's
recommendation as the decision. Tickets #9, #11, and #14 are contestable
calls, flagged for veto in each answer.

## #1: One spec or two?

Blocked by: none
Type: Grill

### Question

The shaping bundles two independently useful behaviors: the census and the
exec-comfort verbs. Does one map feed one spec, or two specs in sequence?

### Answer

Two specs, census first. One map compiles beside each spec. The census spec
lands before the exec-comfort spec, so the verb spec has a measured
baseline. Closed by reviewer 2026-08-25.

## #2: What is one census unit?

Blocked by: none
Type: Grill

### Question

What exactly counts as one raw call into a Bench worktree path?

### Answer

One Bash tool call counts once when its command text names a path under the
repository's pool directory and its verb head is not `bench`. A chain of
commands in one call counts once. A `bench worktree exec` call whose inner
command names the path counts zero. A Read, Edit, or Write tool call never
counts; the path is the right operand for a file read or edit. Closed by
reviewer 2026-08-25.

## #3: Where does the census live, and when does it reset?

Blocked by: none
Type: Grill

### Question

The final-check duty reads the count after `bench worktree land`, so the
count must survive the landing. What is the key and the lifetime?

### Answer

Per assignment, under `$BENCH_HOME`. The hook appends one record keyed by
the assignment whose path the command names. `bench status` shows the live
count per assignment. `bench worktree land` carries the final count in its
landing record as a `census=<n>` key, so the duty reads the count from the
evidence it already holds. `bench worktree release` and `bench worktree
clean` drop the assignment's records. Closed by reviewer 2026-08-25.

## #4: What does one census record hold?

Blocked by: none
Type: Grill

### Question

The reviewer bans a command log for delegates. What does the main session's
record hold?

### Answer

The timestamp, the assignment label, and the command's verb head after
prefix resolution (`git add`, `python3`, `sed`). Never the full command
text. The record names the missing verb without a log a scraper can read.
Closed by reviewer 2026-08-25.

## #5: How does the hook recognize a worktree path?

Blocked by: #2
Type: Research

### Question

Does the census match live assignment paths from the ledger, or the pool
directory as a prefix?

### Answer

The pool directory as a prefix: any absolute path under
`$BENCH_HOME/worktrees/<repo-base>-<crc32>/` counts, live or dead. The
prefix is one string the core already derives (finding 3 in the source), so
the census reads no ledger. The record's assignment label comes from the
path's `<ownerID>-<assignmentID>` segment when the ledger knows it, and from
the raw segment otherwise.

## #6: Which hook event records, and which harnesses does it cover?

Blocked by: #2
Type: Research

### Question

Does the record ride on `PreToolUse` or on a post-execution event, and does
each harness fire it?

### Answer

`PreToolUse` with the `Bash` matcher, the one event both Claude Code and
Codex already wire (finding 5). The kit uses no `PostToolUse`, so the census
records the attempt, not the outcome: a raw call another guard then blocks
still counts. The record call rides the same envelope and the same
`shellcommand.Parse` tokenizer the two guards use (finding 8), so the verb
head has one source.

## #7: What does the ambient dashboard show?

Blocked by: #3
Type: Grill

### Question

What is the census signal's name, severity, detail, and action?

### Answer

Signal name `census`, severity 3, beside `guards` and below the gate, git,
and worktree rows. The detail reads `<label> <n> raw call(s)`, one row per
assignment with `n > 0`. The row fires only for an assignment the ledger
holds as active; a released assignment has no records. The row carries no
action, because the remedy is a habit, not a command.

## #8: What is the final-check exit duty?

Blocked by: #3, #4
Type: Grill

### Question

After the landing, what does `/bench-final-check` do with the count?

### Answer

The duty reads `census=<n>` from the landing record. If `n = 0`, the close
states `census: 0 raw calls`. If `n > 0`, the duty writes exactly one
`bench learning --rule` entry for the landing. The duty is advisory: a
nonzero count never blocks or reds the landing.

The entry has a fixed shape. The title names the assignment label and `n`.
The `--what` text lists each verb head with its count. The `--right` text
names the Bench form for each head, or `none` when no verb exists. The
`--rule` text proposes the verb or the help change. A spec retro cites the
entry under `Agent-experience improvements / Bench CLI` with its `Feeds:`
line.

## #9: How do delegate calls enter the census?

Blocked by: #4, #8
Type: Grill

### Question

The hook cannot tell a delegate's Bash call from the session's. Does a
delegate's call count, and where do its zero to two CLI recommendations go?

### Answer

A delegate's call counts against the assignment whose path it names. A
delegate's own worktree therefore gets its own count. The coordinator reads
that count from `bench status` when it accepts the done-claim. The
reviewer's no-command-log answer holds because a record carries a verb head,
never command text.

The delegate charge template asks for zero to two CLI
improvements derived from the delegate's own calls. The coordinator folds
them into the landing's one learning entry from #8. **Flagged for veto:**
this reads "no command log" as "no command text", and counts delegate calls.

## #10: Does exec need stdin support?

Blocked by: #1
Type: Research

### Question

The parked idea says exec must accept stdin so a heredoc script runs as
`bench worktree exec <label> -- python3 -`. Does it?

### Answer

It does today: `cmd/bench/main.go:69` passes `os.Stdin`, and the exec verb
wires it to the child (finding 1). The gap is knowledge, not code. The
exec-comfort spec adds one stdin example line to `bench worktree exec --help`
and one acceptance row that proves a heredoc on stdin reaches the child.

## #11: What does `bench worktree resolve` do?

Blocked by: #1
Type: Grill

### Question

No Bench verb starts a cherry-pick (finding 4); the agent's raw
`git cherry-pick` leaves the conflicted index. What does the resolve verb
own?

### Answer

`bench worktree resolve <target> <path>...` stages each named path inside
the worktree after the caller edits it. It refuses a named path that still
holds a conflict marker and prints that path. After it stages, it prints the
remaining unmerged paths. When no unmerged path remains, it continues the
in-progress operation (cherry-pick, merge, or revert) with the message
unchanged and prints the resulting commit. The verb never starts the
operation; a `bench worktree pick` verb is out of scope. **Flagged for
veto:** the continue step makes the verb own a Git state Bench did not
create.

## #12: What does `bench worktree show` do?

Blocked by: #1
Type: Grill

### Question

What does the read verb accept and print?

### Answer

`bench worktree show <target> <rev>:<path>` prints the blob bytes to stdout
unchanged and returns Git's exit code. It accepts the one `<rev>:<path>`
operand form and refuses any other operand with the grammar line. It
replaces `git show HEAD:<path>` for a base-side copy.

## #13: Which exec refusals print the worktree path?

Blocked by: #1
Type: Grill

### Question

The parked idea asks for the path on every refusal, so the fallback is never
cheaper than the verb. Which refusals can print it?

### Answer

Every exec failure after the target resolves prints one line
`worktree: <absolute path>` on stderr before the error. A child start
failure and a nonzero child exit both print it. A refusal before the target
resolves has no path; it prints `bench worktree list` as the next action
instead.

## #14: Does the deterministic-step invariant ride in either spec?

Blocked by: #1
Type: Grill

### Question

A parked idea proposes a new `.bench/BENCH.md` invariant: a deterministic
step goes in the CLI, and a judgment step goes to the agent. Does it ride
here?

### Answer

No. The invariant text is a kit-rule change under `craft-synthesis`, and
its first application is the stale-broker repair in `bench worktree land`,
not this work. It stays a parked idea for `/bench-drain`. **Flagged for
veto:** the reviewer named this topic the "deterministic pattern" duty.

## #15: Which terms enter the glossary?

Blocked by: #2, #4
Type: Grill

### Question

Which canonical terms does this work add to `CONTEXT.md`?

### Answer

Three terms, written in this pass:

- **census** — the hook-recorded count of raw calls per assignment; not
  "log", not "audit"
- **raw call** — one Bash tool call that names a pool path with a
  non-`bench` verb head; not "path leak", not "bypass"
- **verb head** — the command's first word after prefix resolution; not
  "command name"

## Not yet specified

## Spec-writer discretion

- Which hook script carries the record call: a third `PreToolUse` shim or a
  branch in an existing one, provided `shellcommand.Parse` decides the verb
  head.
- The census file's exact name and format under the pool directory, provided
  the key is the assignment.
- The TOON key names inside one record.

## Out of scope

- Full command text in any record (#4).
- A separate delegate command log (#9).
- A `PostToolUse` event in either harness (#6).
- Counts for Read, Edit, or Write tool calls (#2).
- A `bench worktree pick` verb that starts a cherry-pick (#11).
- The deterministic-step invariant in `.bench/BENCH.md` (#14).
- Changes to the destructive-git guard's verdicts.

## Sources

- Path: `decisions/assets/exec-census-tree-facts.md`
  Supports: #5, #6, #10, #11, and #13 through the eight cited findings.
  Drift: re-read after any change to `internal/worktree/exec.go`, `internal/worktree/path.go`, or the hook wiring files.
