# roadmap-flow

Status: staged

Decision source: ready compiled map — `specs/roadmap-flow/decisions/roadmap-flow.md` (Status: ready; its one `## Sources` entry, `specs/roadmap-flow/decisions/assets/roadmap-flow-baseline.md`, was re-read in full during this authoring pass)

Verification log: 1 iteration(s) to accept — accepted with two folds: ticket 04 now blocks on ticket 03 (shared conformance registry files), and ticket 04 dropped a vacuous prose-budget acceptance row

## Problem

The board grows faster than it retires. The baseline measures the rate: open
mass moved from 66 rows to 72 rows in two weeks. A drain opens about 1.5 rows
and feeds about 5 rows, and it retires under 1 row unless the reviewer asks for
a restructure pass. Nothing measures this. A reader who wants the rate counts
commits by hand, and each hand count produces a different number.

Two mechanisms let the growth continue. A drain treats an entry that only adds
an occurrence line as a reason to open or feed a row. A retro writes
improvement items that name no destination, so the next drain turns each item
into a new row. Neither the board nor the gate states what a row's next action
is, so a row with no next action stays on the board forever.

## Solution

Bench measures the flow from one source and enforces two mechanical markers.

`bench roadmap --flow` derives the counts from git history: it counts the
`roadmap/FT<n>.md` files each commit adds, modifies, and deletes. It reports
opened, fed, retired, the net delta, the current open mass, and whether the
target holds. The window spans the last three drain commits. `/bench-drain`
quotes that report in its exit, so every drain states the board's flow.

The gate enforces two markers. Every `ROADMAP.md` row carries a `Next:` line in
its detail file, and the value is one of five tokens. Every retro improvement
item carries a `Feeds:` marker. The judgment behind each marker stays advisory
phase text; only the marker is mechanical. One reviewed migration pass brings
the existing rows under the grammar.

Vocabulary, used consistently below. A **flow event** is a commit that adds or
deletes at least one `roadmap/FT<n>.md` file. A **drain commit** is a flow
event that adds at least one such file. The **window** is the commit range that
spans the last three drain commits. **Opened**, **fed**, and **retired** are
the counts of added, modified, and deleted detail files over that window.
**Open mass** is the number of index rows the current board holds. The **net
delta** is opened minus retired. The **target** holds when the net delta is at
or below zero. A **row token** is one `Next:` value. A **feeds marker** is one
`Feeds:` value on a retro improvement item.

## User stories

### Group A — the board's flow is one measured number

Line: opus / medium. The derivation is new, and it reads git history, so a
wrong window produces a confident wrong number; the seam itself is known.

1. As a drain operator, I want `bench roadmap --flow` to report the rows
   opened, fed, and retired over the window, so that I read the board's flow
   instead of counting commits by hand.
2. As a drain operator, I want the flow report to state the current open mass,
   so that I read the board's size beside its rate.
3. As a drain operator, I want the flow report to state the net delta, so that
   the target needs no arithmetic from me.
4. As a drain operator, I want the flow report to state whether the target
   holds, so that the next drain knows when reducing moves are due.
5. As a drain operator, I want the window to span the last three drain commits
   and to name its two boundary commits, so that two runs on one tree agree
   about what the report counted.
6. As a drain operator, I want each flow event derived from added and deleted
   detail files, so that a commit subject cannot inflate or hide a count.
7. As a drain operator, I want a repository with fewer than three drain commits
   to report the drain count it found and still sum the window, so that a young
   board still reports its flow.
8. As a drain operator, I want a repository with no flow event to report a
   definitive empty flow on exit 0, so that silence never reads as a failed
   query.
9. As an agent, I want `bench roadmap --flow` outside a repository to print a
   structured error and exit 1, so that the surface keeps its AXI contract.
10. As a drain operator, I want a degraded board to report the open mass as
    unknown, so that a broken board never reads as an empty board.
11. As an agent, I want `bench roadmap` and `bench roadmap --context` to keep
    their current output, so that the new flag adds a surface and moves none.

### Group B — every row states its next action

Line: opus / medium. The oracle grows a new fault class, and a wrong predicate
here reds honest boards; the parser seam and its fixtures already exist.

12. As a board reader, I want every row's detail file to carry a `Next:` line,
    so that each row names one next action.
13. As a board reader, I want the `Next:` value restricted to `shape`, `spec`,
    `ticket`, `decide`, or `kit-edit`, so that each token maps to one phase.
14. As a reviewer, I want the gate to red on a row whose detail file carries no
    `Next:` line, so that the grammar is enforced and not advised.
15. As a reviewer, I want the gate to red on a row whose `Next:` value sits
    outside the token set, so that a typo cannot pass as a decision.
16. As a board reader, I want a row under the parked section exempt from the
    `Next:` line, so that a row with no honest next action has a legal home.
17. As a maintainer, I want the gate to compare the drain's token table against
    the parser's token set, so that the grammar and its documentation cannot
    drift apart.
18. As a board reader, I want a `Next:` token inside a fenced code block or
    inside backticks ignored, so that a documented example does not parse as a
    live row grammar.
19. As a board reader, I want a `Next:` line that is indented, that wraps onto
    a second line, or that starts with a non-ASCII space refused with a
    diagnostic, so that a reader sees the marker exactly where the parser does.
20. As an agent, I want each `Next:` diagnostic to name the detail file path and
    the fault, so that the repair needs no search.
21. As a board reader, I want a detail file whose last line is the `Next:` line
    without a trailing newline still accepted, so that a hand edit does not red
    the gate.

### Group C — every retro improvement item names the row it feeds

Line: opus / medium. One new check over hand-edited capture files; the retro
parser and the classifier seams already exist.

22. As a drain operator, I want every item under `## Agent-experience
    improvements` to carry a `Feeds:` line, so that each item states its
    destination.
23. As a drain operator, I want the `Feeds:` value restricted to `FT<n>`,
    `new`, or `none`, so that the disposition stays mechanical.
24. As a reviewer, I want the gate to red on an improvement item that carries no
    `Feeds:` line, so that the retro's change test is enforced.
25. As a drain operator, I want an absent or an empty `capture/retros/`
    directory to leave the check quiet, so that a repository with no pending
    retro stays green.
26. As a reviewer, I want a retro file the reader cannot classify to red rather
    than to be skipped, so that an unreadable retro is never counted compliant.
27. As an agent, I want each `Feeds:` diagnostic to name the retro path and the
    item's line number, so that the repair needs no search.
28. As a drain operator, I want an item outside the improvements section left
    alone, so that the check grades only what the retro promises.

### Group D — the drain and the retro say what the flow requires

Line: fable / high. This is the leverage override in `craft-line`: the phase
prose compounds through every drain session, and the gate can only anchor text
that is written well.

29. As a drain operator, I want `/bench-drain` to quote the flow report in its
    exit, so that every drain reports the board's flow.
30. As a drain operator, I want `/bench-drain` to feed a row only when the entry
    changes that row's priority, scope, or `Next:`, so that occurrence-only
    entries stop growing the board.
31. As a drain operator, I want an occurrence-only entry dismissed with one line
    of why, so that the dismissal stays reviewable.
32. As a drain operator, I want a new row to require a `Next:` token and a
    class, so that a row with no next action never opens.
33. As a drain operator, I want a positive net delta to force reducing moves in
    the next drain's batch diff, so that the target carries a consequence.
34. As a drain operator, I want a drained item that meets the light-path
    observables built in the drain session by default, so that small work leaves
    the board instead of queueing on it.
35. As a retro author, I want `/bench-final-check`'s retro template to carry the
    `Feeds:` marker, so that a retro is written compliant.
36. As a maintainer, I want the gate to anchor each of these rules, so that a
    later edit cannot drop one silently.

### Group E — the existing rows come under the grammar

Line: opus / medium. The pass rewrites 72 rows under a reviewer's approval; the
gate is the oracle for the result.

37. As a reviewer, I want one reviewed pass to give every existing row a `Next:`
    token, so that the gate turns green on the real board.
38. As a reviewer, I want a row with no honest next action moved to the parked
    section in that pass, so that no row receives a dishonest token.
39. As a reviewer, I want the pass proposed as one batch diff, so that I approve
    the whole migration once.

## Implementation decisions

- **One new package for the flow derivation.** `internal/roadmapflow` owns the
  window selection, the counts, and the rendering. `internal/roadmap` holds 11
  source files against a budget of 12, so a flow file and its test would breach
  the directory budget, and that budget is reviewer-owned. The new package
  reads the open mass through `roadmap.LoadTree` and `roadmap.ParseDocument`,
  so the row count has one source.
- **The flow report is a flag, not a schema change.** `bench roadmap --flow`
  dispatches from `roadmapCommand` beside the existing `--context` route.
  `bench roadmap --context` keeps schema 4, because `/bench-drain` accepts only
  schema 4 and a bump would stop every drain. `roadmap` is already an approved
  AXI query, so the registry does not move.
- **The derivation reads file adds, modifies, and deletes.** The command runs
  one `git log --name-status --no-renames` query over `roadmap/`, newest first.
  A path matching `roadmap/FT<n>.md` contributes to opened on `A`, to fed on
  `M`, and to retired on `D`. Commit subjects are never read, so no subject
  grammar enters the count. The window walks from the newest commit to the
  third drain commit inclusive, and the report names that commit and the newest
  commit as its boundaries.
- **The report carries no git-sourced text.** The fields are counts, a boolean,
  and two commit identities. Identities are hexadecimal, so a hostile branch
  name, path, or subject reaches no field.
- **The row token set has one source.** `internal/roadmap` exports the ordered
  token set. The parser validates a row's `Next:` value against it, and a new
  conformance check compares `.agents/commands/bench-drain.md`'s token table
  against the same set. The table is the human-facing form and the check is its
  binding, in the shape `line-routing` and `guidance-prose-budgets` already use.
- **The `Next:` line is position-anchored.** The line matches `^Next: <token>$`
  in the detail file, outside any fenced code block. Leading whitespace is not
  permitted, so a backticked or indented example never parses. A separator that
  `unicode.IsSpace` accepts but a reader cannot see is refused with its own
  diagnostic rather than ignored. A second `Next:` line yields a duplicate
  diagnostic that names the second line; the parser never resolves first-wins.
- **The fault class turns on with the migration.** The parser, the token set,
  and the diagnostics land first. The commit that registers the missing-line
  class as a gate fault is the migration commit, so the tree is never red
  between the check and the board it grades.
- **The parked exemption is a section rule.** A row whose index heading sits
  under a `## ` heading containing `Parked` needs no `Next:` line.
  **Flagged for reviewer veto:** the map names the parked section as the home
  for a row with no honest next action, but it does not state the exemption
  mechanism; a heading substring is the cheapest portable rule the split-board
  parser can apply.
- **The row fault classes join `roadmap-detail-integrity`.** The check is
  already bound to the `roadmap-board` input and already owns every split-board
  fault class, so a second checker would be a second reading of the same tree.
  Four classes arrive: a missing `Next:` line, an unknown token, an unanchored
  line (indented, wrapped, or an unreadable separator), and a duplicate line.
  Each class receives one canary fixture, and the family's independently
  authored fixture inventory grows from 9 entries to 13.
- **The retro marker gets its own check.** `retro-improvement-markers` is a new
  Dev check over a new `capture-retros` input source, implemented in
  `internal/retros`. Its subject is the repository root, and it reuses
  `retros.Facts` for classification and `retros.Recommendations` for the item
  units, so the check and the drain read one inventory.
- **The feeds marker is the item's last line.** A recommendation unit is valid
  when its last line matches `^Feeds: (FT[1-9][0-9]*|new|none)$`. A unit that
  carries the marker in another position is refused, so the grammar stays
  readable to a person repairing a retro.
- **Absent is quiet, unreadable is red.** An absent or empty `capture/retros/`
  directory yields no diagnostic. A present directory whose classification fails,
  and a retro file the classifier refuses, both yield a diagnostic naming the
  path and the state. The check never reads a symbolic link or a special file.
- **The phase prose is anchored, not duplicated.** The drain's verdict rules,
  the flow-quoting rule, and the retro template's marker each receive one
  `internal/anchors` registry entry. The registry needle is the enforcement
  copy of the phase text, in the shape the repair-attribution vocabulary
  already uses.
- **The migration is a build step, not a CLI feature.** The reviewed
  `--restructure` drain writes the tokens. No command sorts, tokenizes, or
  rewrites `ROADMAP.md`.
- **FT172 boundary.** This work absorbs FT172's grammar half. FT172 keeps its
  `roadmap_id` half, and this spec adds no `roadmap_id` behavior.

## Testing decisions

- A good test drives a command or a check against a fixture tree and asserts the
  rendered output, the diagnostic text, and the exit code. It never asserts an
  internal structure alone.
- Seams with prior art: the roadmap command seam
  (`cmd/bench/command_registry_test.go` plus `internal/roadmap/roadmaptest`),
  the split-board parser seam (`internal/roadmap/tree_test.go`, which drives
  `ParseDocument` with in-memory bytes), the canary fixture seam
  (`internal/conformance/roadmap_detail_integrity_test.go`, whose owner check is
  called directly and must go red), the retro parser seam
  (`internal/retros/recommendations_test.go`), and the anchors registry seam
  (`internal/anchors/registry_data_test.go`).
- The one new seam is the flow command seam: an in-process
  `roadmapflow.Command` call over a disposable git repository whose history the
  test builds commit by commit. It is the highest seam that still shows a wrong
  window, because the window is only observable in the rendered counts.
- The gate observes through its `test` phase, `go test -count=1 ./...`, which
  also runs the conformance registry against the live tree. No new gate phase
  arrives.

### Seam diagram

    trigger: /bench-drain (or the operator) runs `bench roadmap --flow`;
             the gate runs its conformance registry
        │
        ▼
    git history  ──▶  [ roadmapflow: window → counts → TOON ]  ──▶  flow block, exit 0/1
    roadmap tree ──▶  [ roadmap.ParseDocument: row faults    ]  ──▶  diagnostics
    capture tree ──▶  [ retros: item units → Feeds markers   ]  ──▶  diagnostics
    kit prose    ──▶  [ anchors + token-table binding        ]  ──▶  diagnostics
                          ◀ tests attach here: in-process command calls over
                            disposable repositories, in-memory parser trees, and
                            canary fixtures whose owner check must go red

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RF1 | 1 | a history whose window holds 4 added, 9 modified, and 6 deleted detail files reports `opened=4`, `fed=9`, `retired=6` | flow command | a report that counts commits instead of files, or that counts the whole history, prints other numbers |
| RF2 | 2 | a board of 7 index rows reports `open_mass=7` while the window counts stay unchanged | flow command | reds an implementation that reports the window's opened count as the board's size |
| RF3 | 3 | the same window reports `net=-2` | flow command | reds a net computed as opened minus fed, the cheapest wrong arithmetic |
| RF4 | 4 | a window whose net delta is 1 reports `target_met=false` and a window whose net delta is 0 reports `target_met=true` | flow command | reds a strict comparison that fails a net delta of exactly zero |
| RF5 | 5, 6 | a history of 5 drain commits reports the third-newest drain commit and the newest commit as its boundaries, and a commit whose subject says `drain` but adds no detail file is not a boundary | flow command | reds a window fixed at a commit count and reds a subject-derived drain classification |
| RF6 | 7 | a history with 2 drain commits reports `drains=2` and sums every flow event it found | flow command | reds an implementation that refuses or returns zero counts below three drains |
| RF7 | 8 | a repository whose history touches no detail file prints the flow block with zero rows of evidence and exits 0 | flow command | silence or exit 1 would read as a failed query |
| RF8 | 9 | `bench roadmap --flow` outside a repository prints the structured not-in-repo error on stdout and exits 1 | flow command | reds a panic, an empty report, or an exit 0 |
| RF9 | 10 | a board whose detail directory cannot be listed reports `open_mass=unknown` | flow command | reds an implementation that renders a failed read as zero rows |
| RF10 | 11 | bare `bench roadmap` and `bench roadmap --context` emit their current blocks unchanged | flow command | freezes the two existing surfaces against the new flag |
| RF11 | 12, 20 | a detail file with no `Next:` line yields one diagnostic naming that file path and the missing line | board parser | today no fault class exists, so the row passes; the row reds on a silent accept |
| RF12 | 13, 15 | a detail file carrying `Next: refactor` yields one diagnostic naming the rejected value | board parser | reds an implementation that checks only that the line is present |
| RF13 | 13 | a detail file carrying each of `shape`, `spec`, `ticket`, `decide`, and `kit-edit` in turn yields no diagnostic | board parser | reds a token set that is narrower than the five decided tokens |
| RF14 | 14 | the `roadmap-detail-integrity` owner check called on a fixture whose only fault is a deleted `Next:` line returns that diagnostic, and returns none once the line is restored | canary | proves the gate bites on the class rather than only that a parser function returns a string |
| RF15 | 16 | a row heading under a `## Parked and scheduled work` section with no `Next:` line yields no diagnostic, while the same row under the features section yields one | board parser | reds both a blanket requirement and a blanket exemption |
| RF16 | 17 | a drain token table missing `kit-edit` yields a drift diagnostic naming that token | grammar binding | reds a documentation copy that the gate never compares |
| RF17 | 18 | a detail file whose only `Next: shape` sits inside a fenced code block yields the missing-line diagnostic | board parser | reds a substring search, which would accept the documented example as a live token |
| RF18 | 19 | a detail file whose `Next:` line begins with U+00A0 yields the unreadable-separator diagnostic rather than passing | board parser | an ASCII-only predicate accepts the invisible separator; the row asserts the refused side |
| RF19 | 21 | a detail file whose final `Next: spec` line has no trailing newline yields no diagnostic | board parser | reds a parser that requires a terminating newline on the last line |
| RF30 | 19 | a detail file whose `Next: spec` line is indented by one ASCII space yields the unanchored-line diagnostic naming the line | board parser | reds a parser that trims the line before it matches, which makes the column-zero anchor meaningless |
| RF31 | 19 | a detail file whose `Next:` line ends after `Next:` and carries `spec` on the next line yields the unanchored-line diagnostic | board parser | reds a parser that joins lines, which accepts a value no reader sees on the marker line |
| RF32 | 12 | a detail file that carries two `Next:` lines yields the duplicate diagnostic naming the second line | board parser | reds a first-wins read, which leaves a stale marker with nothing red |
| RF20 | 22, 24, 27 | a retro whose improvement item carries no `Feeds:` line yields one diagnostic naming the retro path and the item's line number | retro parser | today nothing reads the item, so the retro passes; the row reds on a silent accept |
| RF21 | 23 | items carrying `Feeds: FT12`, `Feeds: new`, and `Feeds: none` pass, and an item carrying `Feeds: maybe` yields a diagnostic | retro parser | reds a presence-only check on the marker |
| RF22 | 25 | an absent `capture/retros/` directory and an empty one each yield no diagnostic | retro parser | absence and emptiness are distinct reads that must both stay quiet |
| RF23 | 26 | a `capture/retros/` entry that is a dangling symbolic link yields a diagnostic naming the path and its state | retro parser | a plain read reports the link as not found, so a reader without a stat call calls it an empty retro |
| RF24 | 28 | an item under a different `## ` heading in the same retro yields no diagnostic | retro parser | reds a check that grades every list item in the file |
| RF25 | 24 | the `retro-improvement-markers` owner check called on a fixture whose only fault is a removed `Feeds:` line returns that diagnostic, and returns none once the line is restored | canary | proves the new check is wired into the gate rather than only defined |
| RF26 | 29, 30, 31, 32 | removing the drain's flow-quote sentence, its feeds-a-row test, its dismissal rule, or its new-row rule each yields the anchor diagnostic naming the dropped rule | anchors registry | reds a later prose edit that drops one rule while the rest survive |
| RF27 | 33, 34 | removing the positive-delta restructure rule or the build-in-session default each yields the anchor diagnostic naming the dropped rule | anchors registry | reds the same silent-drop failure on the two rules the map added last |
| RF28 | 35, 36 | removing the `Feeds:` marker from `/bench-final-check`'s retro template yields the anchor diagnostic | anchors registry | a template without the marker writes non-compliant retros that the Group C check then reds |
| RF29 | 37, 38 | after the migration pass, `roadmap-detail-integrity` returns no `Next:` diagnostic over the live tree | gate | reds a partial migration that leaves rows without tokens |

Not covered: story 39 — the batch-diff proposal is a reviewer interaction that
`/bench-drain` already owns; no mechanical oracle observes how a diff is
presented.

Cheapest wrong implementation per group, and the row that reds it. Group A:
count commits rather than files, red by RF1; read the drain from the commit
subject, red by RF5; refuse a history with fewer than three drains, red by RF6.
Group B: check that a `Next:` line is present without checking its value, red
by RF12; search the file for the substring, red by RF17; require the marker on
every row, red by RF15; trim the line before the match, red by RF30; keep the
first of two markers, red by RF32. Group C: check the marker's presence anywhere in the
file, red by RF21; skip a file the classifier refused, red by RF23; grade every
list item, red by RF24. Group D: edit the phase text and leave the anchors
unchanged, red by RF26. Group E: migrate the rows the reviewer noticed, red by
RF29.

### Edge inventory

- A commit that both adds and deletes detail files counts on both sides and
  stays one flow event.
- A rename inside `roadmap/` reads as one add plus one delete, because the query
  passes `--no-renames`; the net delta is unchanged.
- A file under `roadmap/` whose basename is outside the `FT<n>.md` grammar
  contributes to no count, which matches the board parser's own row-ID rule.
- A detail file whose `Next:` value is empty yields the unknown-token
  diagnostic quoting the empty value, not the missing-line diagnostic.
- A retro improvement item that is a numbered list item and one that is a
  paragraph are graded the same way, because `retros.Recommendations` already
  returns both as units.
- An empty retro file yields no items and therefore no diagnostic.
- A retro file whose last line lacks a trailing newline still yields its final
  item's marker.
- A `Feeds:` value naming a row the board does not hold is accepted; the drain
  resolves the destination, and the gate grades the marker only.
- **Won't handle:** rows introduced only by a merge commit's combined diff —
  this repository's history is linear, and `git log` shows no combined diff by
  default.
- **Won't handle:** a rewritten history, where a rebase changes the commit
  identities the window names — the report describes the current history and
  makes no claim across a rewrite.
- **Won't handle:** a hard cap on the number of open rows — the map places it
  out of scope, and the flow target replaces it.
- **Won't handle:** a shelf life or an automatic expiry on a row — the map
  places it out of scope; a stale row leaves through the parked section.
- **Won't handle:** a limit on `capture/learnings.md` entries per landing — the
  journal keeps its existing dismiss path.
- **Won't handle:** the "changes the row" judgment and the restructure trigger
  as gate predicates — the map keeps both advisory, and only the markers are
  mechanical.
- **Won't handle:** FT172's `roadmap_id` half — FT172 keeps it.
- **Won't handle:** a flow block inside `bench roadmap --context` — the schema-4
  contract binds `/bench-drain`, and a bump would stop every drain.

## Ownership fences

- `specs/roadmap-flow/`
- `internal/roadmapflow/`
- `internal/roadmap/tree.go`
- `internal/roadmap/tree_validation.go`
- `internal/roadmap/tree_test.go`
- `internal/roadmap/tree_helpers_test.go`
- `internal/roadmap/roadmap.go`
- `internal/roadmap/roadmaptest`
- `internal/retros/retros.go`
- `internal/retros/retros_test.go`
- `internal/retros/recommendations.go`
- `internal/retros/recommendations_test.go`
- `internal/conformance/registry/registry.go`
- `internal/conformance/checks.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/registry_test.go`
- `internal/conformance/roadmap_detail_integrity_test.go`
- `internal/conformance/retro_improvement_markers_test.go`
- `internal/conformance/row_next_grammar_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `tests/canary/roadmap-detail-integrity/`
- `tests/canary/retro-improvement-markers/`
- `tests/canary/row-next-grammar/`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `cmd/bench/command_registry_test.go`
- `.agents/commands/bench-drain.md`
- `.agents/commands/bench-final-check.md`
- `projects/benchkit.md`
- `ROADMAP.md`
- `roadmap/`
- `CHANGELOG.md`
- `capture/session-handoff.md`

## Out of scope

- A hard cap on the number of open rows, with its own refusal and its own
  reviewer override — 12 edits, 3 gate runs.
- A shelf life that expires a row after a dated threshold — 15 edits, 3 gate
  runs.
- A `flow` block inside `bench roadmap --context` behind a schema-5 bump, with
  the drain's acceptance moved to it — 10 edits, 3 gate runs.
- A `bench status` severity row that fires when the target fails — 8 edits, 2
  gate runs.
- FT172's `roadmap_id` decision and its column — 10 edits, 2 gate runs.

## Further notes

The baseline asset counts `**FT<n>` heading lines for open mass and reads
per-file adds and deletes only after the 2026-08-17 board split. The window
never reaches back past that split in practice, because three drain commits
span about one week. The asset records its own drift note: re-measure once the
flow report ships, because the CLI then replaces the hand count.
