# landing-authors-the-flip

Status: staged

Roadmap: FT113

Decision source: reviewer-confirmed conversation, 2026-08-23, over the drained board entry `roadmap/FT113.md` — the reviewer closed two forks: the landing verb is the one author of the flip and the close, so `bench commit` drops `--spec` and `bench spec implemented` retires; and a retirement that does not complete names its remainder, with no retry verb

Verification log: 1 iteration(s) to accept — the round (opus / high, read-only) returned REVISE with eight blocking findings and eleven prose findings; the author folded every finding (FA5 and FA6 moved to the owner seam with FA6 marked new, FB9 for the production reconcile, shell quotes and the placeholder in `next=`, FA9 for the spec help rows, real ticket edges and the `locateStaged` removal, the NBSP edge, exit 3 in the final-check guidance, the path order); the declared cap was one iteration plus one fold, so no second full pass ran

## Problem

Three commands author the `Status: staged` → `Status: implemented` flip.
They are `bench worktree land --spec`, `bench commit --spec`, and `bench spec
implemented`. They race. A premature `bench commit --spec` on an early ticket
flips the spec before the build ends. A `bench spec implemented` before a
landing makes the landing refuse. Every phase now lands through `bench
worktree land` (ADR 0014), so the two extra authors have no documented caller.

`bench commit --spec` on a tickets-only folder can publish a correct green
commit and then exit 1 with `landed-but-checkout-incomplete`. The confirmed
cause is a caller who names `specs/<slug>` as a path and as the `--spec` close
in one command. The close removes the folder from the published tree, and the
checkout reconcile then fails on that path. A green landing that reports a
red-looking error reads as a failed commit until the caller inspects `git
log`. The error names no commit, no path, and no repair.

`bench spec retire` deletes the spec folder and names "the ROADMAP row" as the
remaining duty. It does not name the row's `roadmap/FT<n>.md` detail file. On
the inherited-toolchain retirement the orphaned detail file consumed a full
red gate before the caller found it.

## Solution

`bench worktree land --spec` is the one author of the flip and of the
tickets-only close. `bench commit` loses `--spec` and becomes a plain
attributed commit: gate, then publish the named paths on green. `bench spec
implemented` leaves the executable. `bench spec retire` and `bench spec
history` stay.

A `bench commit` that publishes its commit and then cannot reconcile the
checkout exits 3, the landing verb's publication boundary. Its report names
the published commit, the path that did not reconcile, and the exact restore
command as `next=`. A refusal before publication keeps exit 1, and a grammar
error keeps exit 2.

`bench spec retire` names the board remainder it leaves. When the spec carries
a `Roadmap: FT<n>` line, the `next:` line names the `ROADMAP.md` row `FT<n>`
and, when that file exists, `roadmap/FT<n>.md`. Without a valid `Roadmap:`
line, the `next:` line names the row and its detail file generically. The spec
template gains the optional `Roadmap: FT<n>` line.

Closed decisions, dated 2026-08-23:

- the landing is the one author
- the "flip counts as a path" face closes by removal of the flag
- no retry verb ships
- Bench does not remove the board row or the detail file
- the tests of the removed commit-route flip and close are deleted with the behavior

## User stories

### Group A — the landing verb is the one author

Line: opus / medium. The commit command and the landing owner are known
seams, the removal is exact, and the gate covers the path; mid because the
landing owner grants publication authority.

1. As an operator, I want `bench commit` to refuse `--spec` as a grammar error with exit 2, so that the commit route cannot author the flip.
2. As an operator, I want `bench commit --help` and the `bench help` commit row to show no `--spec`, so that the advertised grammar has one author.
3. As an operator, I want `bench spec implemented` to exit 2 as an unknown subcommand and rewrite no file, so that no flip races the landing.
4. As an operator, I want `bench spec retire` and `bench spec history` to keep their grammar and behavior, so that the lifecycle keeps its other verbs.
5. As a reviewer, I want `bench worktree land --spec` to keep the flip and the tickets-only close, so that the one author still authors.
6. As an operator, I want `bench commit -m <msg> <path>...` to keep its gate-then-publish path, so that only the spec half leaves.
7. As a session, I want `bench commit -m <msg> specs/<slug>` to publish a tickets-only folder like any path, so that no named path collides with a close.
8. As a reviewer, I want the removed commit-route flip and close tests deleted with the behavior, so that no test grades a gone surface.

### Group B — a published commit whose checkout did not reconcile names its remainder

Line: opus / medium. The reconcile step and the commit's report are the
seams, and the injected-reconcile test is prior art; mid because the report is
the caller's only evidence after publication.

9. As an operator, I want a `bench commit` that published and then failed to reconcile the checkout to exit 3, so that "published" and "refused" differ.
10. As an operator, I want that report to carry the published commit id, so that I learn the commit landed without a second gate.
11. As an operator, I want that report to name the path whose reconcile failed, from the production reconcile, so that the repair targets that path.
12. As an operator, I want the report's `next=` to name the restore over every named path, shell-quoted, so that recovery is one paste.
13. As a script author, I want a refusal to keep exit 1 and a grammar error exit 2, so that readers keep their meanings.
14. As an operator, I want a control-byte path sanitized in `path=` and replaced by a placeholder in `next=`, so that the record stays pasteable.
15. As a cold-session agent, I want `bench commit --help` to name the exit-3 meaning, so that the code is discoverable without the source.

### Group C — `bench spec retire` names the board remainder

Line: sonnet / medium. The retire command's `next:` line is the known seam
and its tests are prior art; cheap because the change is one rendered line.

16. As an operator, I want retire on a `Roadmap: FT<n>` spec to name row `FT<n>` and `roadmap/FT<n>.md`, so that the detail file is not forgotten.
17. As an operator, I want retire without a `Roadmap:` line to name the row and the detail file generically, so that the duty shows.
18. As an operator, I want the detail file named only when it exists on disk, so that the remainder is what remains.
19. As an operator, I want a `Roadmap:` value other than `FT<digits>` to print the generic line, so that a bad value names no path.
20. As an operator, I want retire's deletions and refusals unchanged, so that only the remainder line changes.

### Group D — guidance names one author

Line: fable / high. This group edits kit guidance prose, so the leverage
override in `craft-line` applies.

21. As a spec author, I want the spec template to offer the optional `Roadmap: FT<n>` line, so that retire can name the exact remainder.
22. As a cold-session agent, I want `/bench-final-check` to name the landing as the only flip author, so that guidance matches the executable.
23. As a cold-session agent, I want `/bench-final-check` to name commit exit 3 as published with a repair, so that a green commit reads green.

## Implementation decisions

**The landing verb is the one author.** `bench commit`'s grammar loses the
`--spec` flag, and the landing request loses its spec field. The landing owner
composes only the attributed paths: the spec transition branch, the close
branch, and the close removal at reconcile leave the commit path. The reviewed
landing keeps its transition and its close unchanged.

The tickets-only predicate, the closed-folder path, and the index-tree removal
stay, because the landing verb and `bench status` read them. The `bench spec`
verb loses `implemented`. The working-tree flip function, the staged-check
function, and their shared locate helper leave the spec package. The
derived-implemented-bytes function stays the one source of the flip for the
landing.

**Exit 3 is the commit's publication boundary.** When the ref update
succeeded and the reconcile failed, the command exits 3 and prints one stdout
record in the landing verb's `name{key=value,...}` grammar. The record carries
`published_commit=<sha>`, `path=<path>` for the path whose reconcile failed,
and `next=<command>`. The `next=` value is one `git restore
--source=<sha> --staged --worktree --` followed by every named path,
shell-quoted, in the owner's sorted and deduplicated order. The restore is
idempotent for the paths already reconciled.

Every return of the reconcile step carries the failed path in a typed error,
so the command renders the path without text parsing. `path=` renders through
the control-byte sanitizer, as the landing verb's records do. When a named
path is not line-safe, `next=` replaces the path list with the placeholder
`<named-paths>`, the landing verb's pointer form. The `error:` prefix stays
for exit 1 only.

**Retire names the board remainder.** The retire command reads the spec's
`Roadmap:` value through the existing metadata parser before it deletes the
folder. A value that matches `FT` followed by digits names the row `FT<n>`;
when `roadmap/FT<n>.md` exists as a regular file, the line names it too. Any
other value, or no line, yields the generic line that names the row and its
`roadmap/FT<n>.md` detail file. The deletion order and every refusal stay as
they are.

**Guidance.** The spec template in `craft-spec` gains the optional
`Roadmap: FT<n>` line. The skill keeps its prose budget: the ticket trims one
line elsewhere, or it raises the budget row in the project profile by one.
`/bench-final-check` drops its `bench spec implemented` sentence, states that
the landing verb is the only author, and names `bench commit`'s exit 3 as
published with a named repair. The `bench help` commit row and the `bench
spec` usage text lose their removed surfaces.

**Bootstrap authority.** No executable hop changes. The commit command still
runs the gate's private exact-source build before publication, and the
landing verb's chain is untouched.

## Testing decisions

- A good test drives the real command in a fixture repository with a stub
  gate. It asserts the exit code, the stdout record, the published tree, and
  the checkout state.
- Seams with prior art: the commit command seam (`internal/commit`'s fixture
  harness and `runCommand`). The landing owner seam (`landing_test.go`'s
  injected `authorize` and `reconcile` functions, and its `LandReviewed`
  calls). Also the spec command seam (`spec_test.go`'s retire tests) and the
  registry help golden (`cmd/bench/main_test.go`, `command_registry_test.go`).
- Two rows are new tests, not pins. FA6 is new because the reviewed close has
  no test today. FB9 is new because the production reconcile must attach the
  path without an injected function.
- The gate observes through its `test` phase; no new gate phase.

### Seam diagram

    trigger: an agent runs `bench commit`, `bench spec retire`, or `bench worktree land`
        │
        ▼
    argv + repo state ──▶ [ commit command: grammar → landing owner: compose → gate → publish → reconcile ] ──▶ committed N path(s) | committed{published_commit=…,path=…,next=…} | error: …
                          [ spec command: resolve → validate → delete → next: line ]
                          [ landing verb: proofs → composition → transition/close → gate → publish ]
                      ◀ tests attach here: in-process command calls in fixture repositories;
                        the owner's reconcile and authorize functions are injected or real;
                        assertions read exit codes, records, trees, and checkout state

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| FA1 | 1 | `bench commit -m m --spec x a.txt` exits 2 with the usage line and moves no ref | command | today the flag parses and a green gate publishes a flip |
| FA2 | 2 | `bench commit --help` prints no `--spec`, and the `bench help` commit row reads without `--spec` | registry | reds a removal that leaves the advertisement |
| FA3 | 3 | `bench spec implemented x` exits 2 as an unknown subcommand and leaves `specs/x/spec.md` byte-identical | command | today the command flips the file in place |
| FA9 | 3 | the `bench help` spec rows show `retire` and `history` only | registry | reds a removal that leaves the help row |
| FA4 | 4 | `bench spec retire` on a merged-implemented spec still deletes the folder and exits 0 | command | pins the surviving verb against the removal |
| FA5 | 5 | `LandReviewed` with a staged spec path still publishes `Status: implemented` | owner | reds a removal that reaches the reviewed transition |
| FA6 | 5 | `LandReviewed` with a close path publishes a tree without that folder | owner | the reviewed close has no test today; the deleted commit-route tests were its only coverage |
| FA7 | 6 | `bench commit -m m a.txt` on a green gate publishes `a.txt` and exits 0 | command | reds a removal that breaks the plain path |
| FA8 | 7 | `bench commit -m m specs/<slug>` on a tickets-only folder publishes the folder's files, leaves the checkout clean, and exits 0 | command | today the same folder named beside `--spec` fails reconcile after publication |
| FB1 | 9 | a `bench commit` whose ref update succeeded and whose reconcile fails exits 3 | command | today the command exits 1 |
| FB2 | 10 | that record's `published_commit=` equals the new HEAD | command | today the error names no commit |
| FB3 | 11 | when the second of two named paths fails reconcile under an injected reconcile, `path=` names the second path | owner | reds a report that names the first path or no path |
| FB9 | 11 | the production reconcile, with no injected function, failing on a named path yields `path=` that names it | owner | reds a typed error on the restore return only, which leaves `path=` empty on the other six returns |
| FB4 | 12 | that record's `next=` is `git restore --source=<sha> --staged --worktree --` followed by every named path, shell-quoted, in sorted order, with one path carrying a space | command | today no `next=` exists, and an unquoted space breaks the paste |
| FB5 | 13 | a red gate exits 1 and moves no ref | command | pins the refusal boundary against the new code |
| FB6 | 13 | a missing `-m` exits 2 | command | pins the grammar boundary against the new code |
| FB7 | 14 | a failed path that carries an ESC byte renders `path=` as the sanitizer's spelling and `next=` with `<named-paths>` | command | reds a raw write of the path and an unpasteable restore |
| FB8 | 15 | `bench commit --help` names exit 3 | command | reds an undocumented exit code |
| FC1 | 16 | retire on a spec with `Roadmap: FT7` and an existing `roadmap/FT7.md` prints a `next:` line that names `FT7` and `roadmap/FT7.md` | command | today the line names only "the ROADMAP row" |
| FC2 | 17 | retire on a spec with no `Roadmap:` line prints a `next:` line that names the row and `roadmap/FT<n>.md` | command | today the detail file is not named |
| FC3 | 18 | retire on `Roadmap: FT7` with no `roadmap/FT7.md` on disk names `FT7` and no detail path | command | reds a line that names a path that does not exist |
| FC4 | 19 | retire on `Roadmap: ft7`, `Roadmap: FT 7`, and `Roadmap:` with no value prints the generic line | command | reds a path built from an unchecked value |
| FC5 | 20 | retire on a staged spec still refuses with exit 1 and deletes nothing | command | pins a refusal against the new code |
| FC6 | 16 | a `Roadmap: FT7` line inside a fenced code block is not read, so the generic line prints | command | the metadata parser skips fences; reds a second parser |

Not covered: story 8 — a deletion of tests; the review round verifies it against the tree.
Not covered: story 21 — template prose; the review round verifies it.
Not covered: story 22 — guidance prose; the review round verifies it.
Not covered: story 23 — guidance prose; the review round verifies it.

Cheapest wrong implementation per group, and the row that reds it:

- drop the flag from the help text and keep the parse → FA1
- drop the flag and keep `bench spec implemented` → FA3
- drop the verb and keep its help row → FA9
- remove the close from the commit path and from the landing verb too → FA6
- exit 3 with the old error text → FB2
- name the first named path on any failure → FB3
- type only the restore return → FB9
- print `next=` without the paths, or without quotes → FB4
- always print the generic retire line → FC1
- build `roadmap/<value>.md` from any value → FC4
- name the detail file without a stat → FC3

### Edge inventory

- A named path or a spec slug with a space or a glob character resolves
  literally. The commit fixtures and the retire fixtures carry one (under FA8,
  FB4, FC1).
- A failed path with a control byte renders through the sanitizer, and
  `next=` takes the placeholder (under FB7).
- The named paths in `next=` follow the owner's sorted and deduplicated order,
  not the argv order (under FB4).
- A `Roadmap:` line as the last line with no trailing newline parses (under
  FC1's shape).
- `Roadmap:` present with an empty value is the generic case, not an error
  (under FC4).
- An NBSP around the `Roadmap:` id trims like a space, so `FT<digits>` still
  matches (under FC1's shape).
- The spec file is read through the classifier before the delete, as today; a
  special file where a spec belongs still refuses unread (unchanged).
- A cwd deeper than the repo root still resolves a bare slug at the root
  (unchanged).
- A reconcile failure on the first path leaves every later path unreconciled.
  The restore in `next=` covers every named path, so the order does not matter
  (under FB4).
- **Won't handle:** a retry verb for an unreconciled checkout — the `next=`
  restore is the repair, and `bench commit` re-run reports nothing to commit.
- **Won't handle:** a Bench write that removes the board row or the detail
  file — the spec-retire commit owns the board.
- **Won't handle:** a cross-check of `Roadmap: FT<n>` against the board index
  — the `roadmap-detail-integrity` gate check reds an orphan detail file.
- **Won't handle:** a hook that refuses `bench commit` on the default branch —
  ADR 0014 keeps the rule as guidance until the reviewer decides.
- **Won't handle:** a `--resume` for `bench commit`'s exit 3 — the published
  commit is HEAD and the checkout repair is one restore.
- **Won't handle:** an exit-3 line in `.bench/BENCH-reference.md` — `bench
  commit --help` is the discovery route, and the reference names the landing
  verb's codes only.

## Ownership fences

- `internal/commit/`
- `internal/landing/landing.go`
- `internal/landing/landing_test.go`
- `internal/landing/state_test.go`
- `internal/landing/close.go`
- `internal/spec/spec.go`
- `internal/spec/spec_test.go`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `cmd/bench/command_registry_test.go`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/commands/bench-final-check.md`
- `projects/benchkit.md`

## Out of scope

- A `bench` verb that removes a retired spec's board row and detail file: 3
  edits, 2 gate runs, after a drain-shape decision.
- A hook that refuses `bench commit` on the default branch: 3 edits, 2 gate
  runs, after the reviewer's enforcement decision.
- A reconcile retry verb for `bench commit`: 2 edits, 1 gate run, only if the
  named restore proves insufficient in use.
- The recoverable set-aside and the commit command's smallest sound contract
  stay on FT98 and FT166.

## Further notes

The `Roadmap:` line already parses in the spec metadata reader, and `bench
status` shows it. This spec only documents it in the template and gives it a
consumer. This spec carries `Roadmap: FT113`, so its own retirement names
`roadmap/FT113.md`.

Group C is disjoint from the other groups: its ticket writes only the spec
package and has no blocker. It lands first on its own green gate, so the
FT242 occurrence closes before the flip-author chain starts.

The anchor canary fixtures under `tests/canary/workflow-guidance-anchors/`
still teach `bench commit --spec` and `bench spec implemented`. They are
fixtures, not shipped guidance, and the gate stays green; the retro notes
them.
