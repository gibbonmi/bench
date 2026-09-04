# Handoff sections

Status: staged

Roadmap: FT162

Decision source: roadmap row FT162, carrying its "Decided 2026-09-02" paragraph from the reviewer's grill of that date

Verification log: 2 iteration(s) to accept — the first round blocked on a fence entry that named no file, an undisposed preflight occurrence, and a pin list that dropped the request token

## Problem

The session handoff is one file that describes one phase. Two live phases
rewrite it in full at their boundaries, and the last writer wins in silence.
`bench handoff` regenerates the Next command from the board on every run, so a
mid-build resume invocation vanishes. It keeps the prior State under fresh
pins, so the file pins one commit and describes another.

The verb knows no request. A run from a worktree pins that worktree, and a run from the primary checkout
pins `main`. Nothing says which section the caller owns.

`bench status` dates the whole file by its write time, so a stale sibling reads
current. No linked repo git-ignores the file, so the write lands in the
caller's worktree and dies with it. The reviewer wants concurrent phases as the
normal case, and this file cannot hold two.

## Solution

The handoff becomes one git-ignored file with one section per live assignment,
keyed by request id, plus one `main` section. A leaf package owns the section
grammar, the locked read-modify-write, and the removal. `bench handoff`
resolves the caller's section from the assignment record and rewrites only
that section. It keeps a non-empty Next command, and it refuses a State that
names a commit outside the tip's ancestry.

The shared retirement path removes the assignment's section and ensures
`main`. `bench status` dates each section
against its own assignment branch. `bench init` writes the ignore entries into
a linked repo.

The implementation of this spec starts after the `binary-freshness` landing,
because both edit the working agreement.

## User stories

### Group A — the section grammar and its store

Line: opus / medium. The package is new, and it is the one parser two other
packages import.

1. As a maintainer, I want one leaf package to own the handoff grammar, so that three readers read one shape.
2. As a reviewer, I want each section to open with a level-two heading that carries its key, so that my State headings do not collide.
3. As a reviewer, I want every section field on one label line with no terminator, so that the prose lane skips it.
4. As a coordinator, I want a section to carry the six pins the grill named, so that a fresh session resumes from the file.
5. As a coordinator, I want a section to list every live spec in its worktree, so that a worktree with two staged specs loses neither.
6. As a coordinator, I want two writers on distinct sections to both survive, so that concurrent phases cannot clobber each other.
7. As a maintainer, I want the Shape text to describe the section grammar, so that the single-source check stays true.

### Group B — the verb owns one section

Line: opus / medium. The lookup crosses the intent ledger, and three refusals
need exact predicates.

8. As a coordinator in a worktree, I want `bench handoff` to resolve my section from my path's assignment, so that I never name a request.
9. As a coordinator in the primary checkout, I want `bench handoff` to own the `main` section, so that a phase close with nothing live still writes.
10. As a coordinator, I want `bench handoff` to refuse from a path that owns no section, so that a stray checkout writes nothing.
11. As a coordinator, I want `bench handoff` to re-emit every other section byte for byte, so that a sibling's State survives my close.
12. As a coordinator, I want the header to carry the repository, the path, the `main` HEAD, and the gate verdict, so that the pins resolve.
13. As a coordinator, I want a section's worktree tip read from the assignment's worktree, so that the pin names my tree and not `main`.

### Group C — the Next command and the stale State

Line: opus / medium. Two refusals with exact predicates over free prose.

14. As a coordinator, I want a non-empty Next command kept when I pass no `--next`, so that a mid-build resume invocation survives.
15. As a coordinator, I want the board's leading signal written only into an empty Next command, so that a fresh section still gets a route.
16. As a reviewer, I want `bench handoff` to refuse a State that names a commit outside the tip's ancestry, so that a stale body cannot ship.
17. As a reviewer, I want a hex word that is not a commit ignored by that refusal, so that English never blocks a close.
18. As a reviewer, I want a tree hash in State to refuse, so that the pinned-a-tree failure is caught.

### Group D — the lifecycle and the status row

Line: opus / low for the hooks, opus / medium for the status row. The
retirement seam is one function, and the row joins two data sources.

19. As a coordinator, I want `bench worktree land` and `bench worktree release` to remove the assignment's section, so that a dead request cannot pin the file.
20. As a coordinator, I want `bench worktree clean` to remove the section by the same path, so that no verb leaves a section behind.
21. As a coordinator, I want the file to hold `main` after the last section is removed, so that a fresh session still reads a State.
22. As a session, I want `bench status` to date each section by the branch commits past its recorded tip, so that a stale sibling reads stale.
23. As a session, I want the `main` section dated by the file's write time as today, so that the primary row keeps its meaning.

### Group E — linked repos and the working agreement

Line: opus / low for the ignore entry, opus / medium for the guidance.

24. As a consumer, I want `bench init` to add the three capture files to the ignore file, so that the write lands in the primary checkout.
25. As a consumer, I want that addition idempotent, so that a second `bench init` adds nothing.
26. As a session, I want the working agreement to state the section rule, so that a phase close writes its own section.

### Group F — the review repairs

Line: opus / medium. Each repair changes a refusal or a parser edge.

27. As a coordinator, I want a legacy document read as one `main` section, so that the new verb's first run does not refuse it.
28. As a coordinator, I want an unterminated fence in State refused at write and at parse, so that one section cannot swallow its siblings.
29. As a coordinator, I want a State line that opens a level-two heading refused at write, so that my prose cannot break every later verb.
30. As a coordinator, I want a retirement whose section removal fails to print the error, so that a pinned dead section is never silent.
31. As a reviewer, I want a commit hash inside a fenced block in State left unscanned, so that a quoted example never blocks a close.
32. As a reviewer, I want an ambiguous abbreviation in State refused with its own reason, so that a stale short pin cannot pass as prose.

## Implementation decisions

- A new leaf package `internal/handoffdoc` owns the grammar. The document is a header block, one `## main` section, one `## request <id>` section per live assignment, and the trailing `## Shape`. A section holds label lines for the six pins, one `Next command:` line, and a `### State` body. The package parses, renders, removes a section by key, and ensures `main`. It wraps every read-modify-write in an exclusive flock on a lock file beside the document, with temp-and-rename underneath.
- The section key is the request digest. The `Request token:` line carries the plain token verbatim.
- `intent.AssignmentForWorktree` is promoted from the preflight lookup and exported. `bench handoff` owns `main` when `git.IsPrimaryCheckout` is true, owns the matching active assignment's section otherwise, and refuses when neither holds.
- A section's spec lines come from `spec.Facts` over the assignment's worktree, one pair per live spec. The tip comes from `git rev-parse HEAD` in that worktree, or from the branch when the tree is gone.
- Without `--next`, the verb keeps a non-empty `Next command:` value byte for byte and calls the board route only for an empty one. Whitespace-only and empty backticks count as empty.
- The State scan finds backticked runs of seven to forty hex characters and skips those that name no object. A token that names an object but fails `cat-file -e <token>^{commit}` refuses with a not-a-commit reason, so a tree hash refuses. A commit for which `merge-base --is-ancestor <token> <tip>` is false refuses with an ancestry reason. Each refusal prints the line.
- `executeCleanup` in `internal/worktree` removes the section for the retired assignment and ensures `main`. It is the one call site, guarded by a single-owner test in the census shape. `internal/worktree` imports the leaf package, never `internal/handoff`.
- `bench status` counts `rev-list --count <section tip>..<branch tip>` per section and names the behind section in its row. `main` keeps the file-age rule. Per-section expansion follows the census expansion under `--all`.
- `bench init` appends `capture/session-handoff.md`, `capture/IDEAS.md`, and `capture/learnings.md` to the root ignore file when absent.
- The working agreement's phase-close paragraph states the section rule in the five-part anchor shape.

## Testing decisions

- The grammar uses in-memory round-trip tests, a prose-lane render test, and a two-writer parallel test.
- The verb uses a two-assignment fixture with byte comparison of the untouched section, and three refusal tests.
- The lifecycle uses the census-drop test family and its single-call-site guard.
- The status row uses a two-section fixture with one behind assignment.
- The ignore entry uses the init scaffold tests.
- Guidance uses an anchor tuple and a live-mirror fixture.

### Seam diagram

    bench handoff (worktree or primary)
        │
        ▼
    assignment lookup ──▶ [ handoffdoc: lock, parse, rewrite one section ] ──▶ capture/session-handoff.md
                                   ◀ tests attach here: round trip, parallel writers, refusals

    bench worktree land | release | clean
        │
        ▼
    executeCleanup ──▶ [ handoffdoc.Remove + EnsureMain ] ──▶ the file without that section
                                   ◀ tests attach here: the census-drop family

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| HS1 | 1, 2, 3 | a rendered two-section document parses back to the same sections and passes `prose.Findings` | new round-trip and prose tests in the leaf package | a terminator on a field line grades it as prose |
| HS2 | 4 | a section renders `Request token:` with the plain token byte for byte, `Label:`, `Worktree tip:` equal to the worktree HEAD, `Recorded base:`, `Spec:`, and `Spec status:` | new render test with a commit inside the worktree | a section that renders the digest and omits the token repeats the dead-token occurrence |
| HS3 | 5 | a worktree with two staged specs renders two spec pairs | new render test | a single pair loses one spec |
| HS4 | 6 | two writers on distinct sections, N rewrites each, leave both sections present | new parallel test in the journey-race shape | temp-and-rename alone loses one in silence |
| HS5 | 7 | the Shape text names the section grammar and the single-source check passes | `checkHandoffShape` and `TestHandoffShapeSingleSourcedBites` | a Shape that describes the old shape drifts from the verb |
| HS6 | 8, 11 | `bench handoff` from assignment A's worktree rewrites A's section and leaves B's bytes identical | new two-assignment verb test with byte comparison | a full rewrite drops B |
| HS7 | 9 | `bench handoff` from the primary checkout writes `main` | new verb test | a verb that requires an assignment refuses the phase close |
| HS8 | 10 | `bench handoff` from a path with no active assignment that is not primary exits 1 and writes nothing | new verb test | a match by label or path adopts a stale section |
| HS9 | 12 | the header carries the repository, the path, the `main` HEAD, and the gate verdict | new render test | a header that pins the worktree HEAD repeats the 2026-08-29 occurrence |
| HS10 | 13 | the section tip equals the worktree HEAD | HS2 | covered by HS2 |
| HS11 | 14 | without `--next`, a `Next command:` with flags survives byte for byte on a board whose leading signal differs | new verb test asserting the full line | a `Contains` check passes a truncated rewrite |
| HS12 | 15 | an empty, whitespace-only, or empty-backtick Next command receives the board's leading signal | new verb test with three rows | a rule that treats backticks as content leaves a fresh section routeless |
| HS13 | 16 | a State that names a real commit off the tip's ancestry exits 1 and prints the line | new verb test | an existence check alone passes the off-ancestry commit |
| HS14 | 17 | a State that holds `facade` and `decade` in backticks exits 0 | new verb test | a hex regex without the commit test refuses English |
| HS15 | 18 | a State that names `HEAD^{tree}` exits 1 | new verb test | a check without `^{commit}` passes a tree |
| HS16 | 19 | a landing removes the assignment's section | `TestLandCommandStatesTheCensusCountAndDropsTheRecords` extended | a hook on the release verb alone misses the landing route |
| HS17 | 19 | a release removes the assignment's section | `TestReleaseDropsTheCensusRecords` extended | a hook on the landing alone misses release |
| HS18 | 20 | `clean --apply` removes the section | `TestCleanDropsTheCensusRecords` extended | a hook on `ReleaseCommand` misses clean |
| HS19 | 19 | the removal has one call site in the package | new single-call-site test in the census-drop shape | two sites drift |
| HS20 | 21 | after the last section is removed, the file holds `main` with no further verb run | new lifecycle test | a verb-time `main` leaves the file empty after release |
| HS21 | 22 | with two sections, one three commits behind on its branch, the status row names that section with `3 commits behind` | new status test | a file-level clock reads both current |
| HS22 | 22 | rewriting the fresh section leaves the behind row's count unchanged | new status test | a shared write time resets the sibling |
| HS23 | 23 | `main` is dated by the file write time | `TestAppendHandoffIgnoredUsesFileTime` and `TestAppendHandoffIgnoredAbsentIsSilent` | a per-branch rule has no branch for `main` |
| HS24 | 24, 25 | `bench init` adds the three entries once and a second run adds nothing | new init scaffold test | a blind append duplicates the lines |
| HS25 | 26 | the working agreement states the section rule and the anchor pins it | anchor registry test and fixture `agents-handoff-section-rule` | the rewrite-in-full sentence survives |
| HS26 | 27 | a legacy document with State, Closed decisions, Next command, and Shape reads as one `main` section, and its render parses back | legacy tests in the leaf package | a parser that knows only the section grammar refuses the file it replaces |
| HS27 | 28 | `bench handoff` with a State that opens a fence and never closes it exits 1, prints the line, and writes nothing | new verb test | a write-time check that trusts the parser lets the fence reach disk |
| HS28 | 28 | `Parse` refuses a document whose fence is still open at the end of the file, with the file and line | new leaf test | a per-file fence state absorbs every later section into one State |
| HS29 | 29 | `bench handoff` with a State line that opens a level-two heading outside a fence exits 1 and prints the line | new verb test | a heading the reviewer writes bricks the document for every later verb |
| HS30 | 30 | a retirement whose section removal fails prints the error and keeps its verdict | new lifecycle test with an unparseable document | a discarded error leaves the dead section pinned in silence |
| HS31 | 31 | a real commit off the tip's ancestry inside a fenced block in State exits 0 | new verb test | a scan that splits on newlines alone refuses a quoted example |
| HS32 | 32 | a State that names an ambiguous 7-hex abbreviation exits 1 with an ambiguity reason | new verb test | two failed `cat-file` probes read the token as prose |

### Edge inventory

- Error paths: an ambiguous section heading and a missing `main` each refuse with the file and line. A lock held past the two-second deadline the intent ledger uses refuses with the lock path.
- Empty input: a fresh repo's first `bench handoff` from the primary checkout writes the header and `main`.
- Boundary values: a seven-character hex token is scanned; a six-character one is not.
- Interrupted state: a writer killed between temp and rename leaves the prior document intact.
- Re-run idempotency: a second `bench handoff` with no changes rewrites the same bytes.
- Hostile paths: a worktree path spelled through a symlink matches its assignment through the canonical owner.
- Partial implementation: a section written at create time and never removed reds HS16 to HS18.

**Won't handle** — a section created at `bench worktree create` — creation is a rollback transaction, and the first `bench handoff` in the worktree creates the section lazily.

**Won't handle** — a `Spec` field on the assignment record — the section lists every live spec from the worktree instead.

**Won't handle** — a per-section write timestamp — the age counts commits past the recorded tip, which needs no clock.

**Won't handle** — a tracked handoff in any repo — the decision ignores the file everywhere, and `bench init` writes the entry. The ignore also dissolves the 2026-08-20 occurrence, because no handoff commit can put `main` ahead of a frozen base.

**Won't handle** — `bench preflight build` reading the retained source's frozen base by default — the 2026-08-28 occurrence is preflight work, and FT200 owns the preflight chokepoint.

## Ownership fences

- `specs/handoff-sections/`
- `reviews/handoff-sections.md`
- `internal/handoffdoc/`
- `internal/handoff/`
- `internal/intent/assignment.go`
- `internal/intent/assignment_lookup_test.go`
- `internal/preflight/gather.go`
- `internal/preflight/gather_test.go`
- `internal/worktree/lifecycle.go`
- `internal/worktree/land_journey_test.go`
- `internal/worktree/worktree_test.go`
- `internal/status/handoff.go`
- `internal/status/handoff_test.go`
- `internal/status/status.go`
- `internal/adopt/init.go`
- `internal/adopt/adopt_test.go`
- `internal/conformance/handoff_single_source_test.go`
- `AGENTS.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/load-validity-metadata/shared-rule-drift`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `cmd/bench/testdata/anchors/pre-disclosure-populated.stdout`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/main_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_test.go`

The fence is the union of the eight tickets' `Writes:` lines, closed by
`bench preflight build` over the fixture and registry pins. A closure
headroom file creates no blocker edge; only a file a ticket's `What to build`
names does.

## Out of scope

- A `Spec` field on the assignment record: 3 edits, 1 gate run.
- A `bench handoff --request` flag for a foreign section: 2 edits, 1 gate run.
- A handoff row per section in the default `bench status` view: 1 edit, 1 gate run.

## Further notes

Flagged additions beyond the decision source:

- A section lists every live spec in its worktree. The grill named one spec path and status; the record binds no spec to an assignment, and a worktree can hold two.
- A section's age counts commits past its recorded tip. The grill said "dates each section against its own assignment" and named no data source.
- `bench init` writes the ignore entries. The grill said every linked repo ignores the file, and nothing writes that today.
- The leaf package and its lock. `internal/worktree` cannot import `internal/handoff`, and no capture-file lock exists.
- `bench worktree clean` removes the section. The grill named land and release; clean shares their retirement path, so it comes for free.
- The rule that the verb never rewrites State is existing behavior, covered by the section tests today, and takes no new row.

Build decisions recorded for veto, made under the `--full` batch approval:

- The lock is reclaimed on release with the safe-unlink shape. The alternative, a fourth ignore entry for the lock file, widened the landing predicate and the init list.
- The legacy document is migrated on read (story 27). The pre-existing handoff in this repo refused the new parser on the first run.
- The State scan prints one of two reasons, not-a-commit or off-ancestry, so the commit peel is observable.
- The status row's unresolved-section advisories are cut. Spec line 109 names the behind section alone, and a dead section is residue the retirement removes.
- The lock deadline is mirrored from the intent ledger's unexported literals. The collapse needs an export from `internal/intent`, and stays a reviewer decision.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| one file, one section per live assignment keyed by request id, plus `main` | HS1, HS6, HS7, HS20 |
| a section carries request id, label, tip, base, spec path, spec status | HS2, HS3, HS10 |
| the header carries repository, path, `main` HEAD, gate verdict | HS9 |
| land and release remove the section | HS16, HS17, HS18, HS19 |
| the verb rewrites only the caller's section and refuses another request's | HS6, HS8, HS4 |
| without `--next` keep a non-empty Next; board signal only into an empty one | HS11, HS12 |
| refuse when State names a commit outside the tip's ancestry | HS13, HS14, HS15 |
| `bench status` dates each section against its own assignment | HS21, HS22, HS23 |
| every linked repo git-ignores the file | HS24 |
| the Shape and the working agreement | HS5, HS25 |

Every subagent runs `opus` at low or medium effort.
