# Decision-map split: an index map, one file per ticket, per-topic assets

Status: staged

Decision source: the ready compiled map `specs/decision-map-split/decisions/craft-research.md`, tickets #10 through #14, resolved 2026-09-06. Both `## Sources` entries were re-verified on 2026-09-06 before the seams below were chosen. The research asset's local citations were re-resolved: nine hold, seven moved, and four are gone. The gone citations point at FT135 spec files that retired with promote-then-delete. They support only the fan-out precedent in map ticket #3, which stays closed. The assessment asset's four external report links are readable but unverifiable, as its own Evidence status states.

Verification log: 2 iteration(s) to accept — round one rejected on seven blocking findings, the fold split the guidance ticket and restated four rows, and round two verified the folds with three small folds folded in

## Problem

A decision map is one file that loads whole into every shaping session. This map is 471 lines, and the gate-budget map is 1011 lines. A session that resumes one ticket reads every other ticket first. The answers, the index, and the frontier live in one file. No reader can return one ticket as a slice.

Map-owned assets sit in one flat folder keyed by a name prefix. An asset that no live map names has no home rule. Research for the repository cannot live at `docs/research/`, because the shift-scratch ignore rule matches a `research` folder at any depth. `bench maps` omits a ready map from its default view. Nothing projects a stale citation.

## Solution

The map file becomes an index. It keeps the title, Status, Destination, and the four terminal sections. It gains Notes and Decisions so far. Each decision ticket is one file under the map's tickets folder, and the answer lives only there. Map-owned assets live in the map's assets folder. The topic folder moves into a spec as one unit.

`bench maps` reads the ticket files and derives the frontier from `Blocked by` and placeholder answers. It gains an optional map argument, a `path` column, and one `ready` row per ready map. It adds an advisory `stale` row for a Sources path with a commit newer than the map's. One plan-before-apply program migrates every map, moves the assets, and narrows the ignore rule. The commands, the glossary, the README, and one ADR record the shape.

## User stories

### Index map and ticket files

Line: opus / medium. The spec is exact, the field-scan seam has a precedent, and the integrity fixture family covers the parse.

1. As a shaping session, I want the map file to hold only the index sections, so that I read the low-resolution view once.
2. As a shaping session, I want each decision ticket in its own file under the map's tickets folder, so that one ticket is one read.
3. As a shaping session, I want the basename as the ticket number and the title line as its name, so that gists link by number.
4. As the gate, I want an inline `## #n:` heading in a map to red with the ticket file path, so that the old shape dies.
5. As the gate, I want a ticket basename that is not a plain positive integer to red, so that the number stays the id.
6. As the gate, I want a ticket file with a missing or unsupported field to red with today's diagnostics, so that the field rules survive.
7. As the gate, I want a map with no ticket file to red with `missing decision ticket`, so that an index alone cannot pass as a map.
8. As the gate, I want a tickets folder with no index file to red, so that an orphan folder cannot hide unread tickets.
9. As the gate, I want `## Notes` and `## Decisions so far` required in the index, so that every map carries its preferences and its index.
10. As the gate, I want a resolved ticket with no gist line to red, so that the index cannot fall behind the answers.
11. As the gate, I want a gist that links an unresolved ticket to red, so that the index cannot run ahead of the answers.
12. As the gate, I want a gist to a missing file, or a line without the link grammar, to red, so that every gist resolves.
13. As the gate, I want a `Blocked by` that names a ticket with no file to red as dangling, so that the graph stays closed.
14. As the gate, I want the graph, readiness, terminal-list, and Sources rules to grade a split map as before, so that no rule is lost.
15. As the gate, I want the compiled shape under a spec graded the same way, so that a compiled map keeps its tickets.
16. As a shaping session, I want `--template` to print the index skeleton and `--ticket-template` one ticket skeleton, so that both stay paste-ready.
17. As a shaping session, I want the template's asset rule to name the map's assets folder, so that the skeleton teaches the new home.
18. As a shaping session, I want a map whose name holds a space to keep working, so that the folder rule narrows no existing name.

### Projection

Line: opus / medium. The projection seam is the existing command, the golden output fixtures cover it, and the row shape is decided.

19. As a session, I want the default output to keep every frontier, blocked, deferred, fog, and invalid row, so that the frontier is unchanged.
20. As a session, I want one `ready` row per ready active map with a write-spec help action, so that ready work shows too.
21. As a session, I want every row to carry a `path` column, so that the row names the file I read next.
22. As a session, I want `bench maps <map>` to print only that map's rows, so that a resume reads one map.
23. As a session, I want an unknown map name to exit 1 with a refusal naming the operand, so that a typo is not empty.
24. As a session, I want an operand with a slash, `.md`, or a control byte refused, so that no path poses as a name.
25. As a session, I want a `stale` row when a Sources path's last commit is newer than the index's, so that drift is visible.
26. As a session, I want no stale row when either file has no commit, so that a fresh file is not drift.
27. As a session, I want stale rows to leave the exit code at zero and `--count` unchanged, so that freshness stays advisory.
28. As `bench status`, I want the decisions row and its ready action unchanged, so that the board reads the same counts.

### Migration

Line: opus / medium. The program is one-shot, its input set is the enumerated tree, and the differential row grades it.

29. As the reviewer, I want the program to print every target and stop without `--apply`, so that I inspect the plan first.
30. As the reviewer, I want `--apply` to split every active and compiled map in one pass, so that no inline map survives.
31. As the reviewer, I want each split index to keep its sections, gain Notes, and seed its gists, so that nothing is lost.
32. As the reviewer, I want each asset a map's Sources names moved into that map's assets folder with references rewritten, so that no citation breaks.
33. As the reviewer, I want an asset no map names moved to `docs/research/<slug>.md` and the flat folder removed, so that one home rule holds.
34. As the reviewer, I want the ignore rule narrowed to the root shift-scratch folder, so that `docs/research/` is trackable.
35. As the reviewer, I want the migrated tree to pass the map lane with the same unresolved rows as before, so that every state survived.
36. As the reviewer, I want the program deleted once the inline shape is refused, so that no second parser survives.
37. As the gate maintainer, I want the integrity fixture family rebuilt on the split shape with each mutation biting, so that no red is lost.
38. As the gate maintainer, I want one new fixture per new diagnostic, so that each new rule has a proven bite.

### Guidance and record

Line: opus / high. Guidance prose steers every later session, so the leverage override routes it mid plus high.

39. As a shaping session, I want the shape-idea command to teach the index, the ticket file, and both templates, so that cold sessions write right.
40. As a shaping session, I want the command to name a ticket by title with the number beside it, so that numbers read as names.
41. As a shaping session, I want a recommendation that asserts current-code behavior to name the evidence read this session, so that no premise closes unread.
42. As a spec session, I want the write-spec command to move the topic folder as one unit, so that tickets and assets travel together.
43. As a cold reader, I want `CONTEXT.md` to define map index, decision ticket, and gist with Avoid lists, so that each has one name.
44. As a cold reader, I want one ADR to record the ticket-file decision and the index map, so that the split is recorded.
45. As a cold reader, I want the README and the field guide to describe the split shape, so that the public docs match the tree.
46. As the gate maintainer, I want every moved anchor needle updated with its mutation row and fixture in one diff, so that it still bites.
47. As the reviewer, I want the phase-close retro to record the lines read to resume this map before and after, so that FT231 has data.

## Implementation decisions

**One schema owner, two file shapes.** `internal/maps` keeps the decision-map schema. The index parse keeps the title, Status, Destination, and the four terminal sections. It gains `Notes` and `Decisions so far` as required sections. Notes is free prose, and an empty Notes passes. Decisions so far is a Markdown bullet list, and an empty list passes only when no ticket is resolved.

**The gist line.** A gist line is `- [<name>](<topic>/tickets/<n>.md): <gist>`. The link is relative to the map file. The parser resolves it by the `<n>` segment against the ticket files it discovered.

**The ticket file.** A ticket file starts with a `# <title>` line. It carries `Blocked by:`, `Type:`, `### Question`, and `### Answer` with today's field rules. The whole file is one field scope. The basename without `.md` is the ticket id, and it must match `[1-9][0-9]*`. Discovery lists the direct `.md` children of `<map dir>/<topic>/tickets/`, skips dotfiles and README, and reads no nested folder.

**New diagnostics.** An `## #n:` heading in a map index reds with `inline ticket #n: move it to <topic>/tickets/<n>.md`. A tickets folder without its index reds with `<topic>/tickets: no map index at <topic>.md`. Every diagnostic keeps the map index path as its prefix. A ticket-file diagnostic keeps today's `ticket #n:` form.

**Discovery and compiled maps.** The candidate scan stays the direct-child policy for `decisions/` and for each `specs/<slug>/decisions/`. The ticket and asset folders sit one level below the index, so the scan never reads them as maps. The compiled flag and its readiness rule are unchanged.

**Existing rules unchanged.** The graph walk, the readiness rules, the terminal bullet-list rule, and the Sources grammar keep their code and messages. So do the Sources Path validity check and the `## Handoff` refusal.

A missing Sources Path stays a validity diagnostic and a gate red. That check exists today, and the four invariants forbid a weaker check. The map's ticket #13 names a missing path as a stale row. This spec records that deviation for reviewer veto and keeps the existing check.

**Templates.** `DecisionMapTemplate` renders the index skeleton with Notes and Decisions so far. A new `DecisionTicketTemplate` renders one ticket skeleton, and `bench maps --ticket-template` prints it. The template's asset-rule sentence names `decisions/<topic>/assets/`. The two template flags and `--count` are mutually exclusive.

**Projection grammar.** `bench maps [<map>] [--count|--template|--ticket-template]`. The row shape is `maps[N]{map,title,type,state,blockers,path}`. A ticket row's path is its ticket file. A fog, invalid, or ready row's path is the map index file. A ready row carries the map title, type `map`, state `ready`, and the help action `/bench-write-spec <map path>`.

**Stale rows.** A stale row carries the Sources locator as its title, type `source`, state `stale`, an empty blockers cell, and the asset path. The test compares `git log -1 --format=%ct -- <path>` for the asset and for the map index file. It runs through the git output helper with the repository root as the working directory. A URL source never gets a stale row. A stale row changes no exit code and no count.

**The map operand.** The operand filters the scanned rows by map name after the scan. A well-formed name that matches no active map exits 1 with `maps: no active map named "<name>"`. An operand with `/`, a `.md` suffix, or a byte below 0x20 is a grammar error and exits 2. `ActiveCounts` and the status row are unchanged.

**Migration program.** A one-shot Go program at `scripts/split-decision-maps/main.go` runs as `go run ./scripts/split-decision-maps [--apply]`. Without `--apply` it prints every target: each map with its ticket files, each asset with its destination, each reference file, and the ignore line. With `--apply` it writes them. It parses the inline shape through the expand-era parser and writes the split shape. It seeds each gist from the answer's first sentence.

**Migration moves.** The program moves each asset named by a map's Sources into that map's assets folder. It rewrites every reference to the old path in tracked files under `decisions/`, `specs/*/decisions/`, and `docs/`. It never edits `specs/*/spec.md` or `specs/*/tickets/`, because a build may not edit a spec. It then searches the whole tree for the old paths and reports any remaining hit. The template and command needles that carry the old path change in the contract and guidance tickets.

It moves an asset that no map names to `docs/research/<slug>.md`, then removes `decisions/assets/`. It replaces the `research/` ignore line with `/research/` under a comment that names the shift-scratch folder.

No linked repository holds a map on 2026-09-06, a reviewer-supplied premise, so this repository is the program's only audience.

**Expand, migrate, contract.** The parser first accepts both shapes, and the index-drift rules run only for a map that has a tickets folder. The migration then lands the split tree and the rebuilt fixture family. The contract then refuses the inline shape, makes the drift rules universal, renders both templates, and deletes the program with the inline parse. The tree is green at each commit.

**Test seams and helpers.** A test that needs a valid map writes the index from `DecisionMapTemplate` and one ticket from `DecisionTicketTemplate`. No test-only helper crosses a package. The status, conformance, and command-registry tests that write a template map today move to the two-file form.

**Structure.** `internal/maps/schema.go` is over budget. New code lands in `internal/maps/tickets.go` and `internal/maps/freshness.go` with their test files. That brings the package to its twelve-file budget, and a later file needs a split or a grant.

**Record.** ADR 0020 records three decisions: a decision lives in its ticket file, the map is an index, and the topic folder moves as one unit. It names no path. `CONTEXT.md` defines decision map, map index, decision ticket, and gist. The anchor rows that pin `decisions/assets/`, the move sentence, and the retirement sentence change with their prose. The mutation table and the canary fixture that mutate the asset-path constant change with them.

## Testing decisions

- A good test writes an index and ticket files into a temporary tree, runs the tree validation or the command, and asserts the output. It never asserts an internal.
- The parse seam is `maps.ValidateDecisionMapTree`, exercised by the integrity fixture family and the package tests. Prior art: `TestDecisionMapIntegrityCheckValidatesEveryCandidate`.
- The projection seam is `maps.Command` over a `gittest.Repo`. Prior art: `TestCommandAppendsOnlyMapActionsToTheCapturedPrimaryResponse` and the golden files under `internal/maps/testdata/`.
- The status seam is `appendMaps`. Prior art: the ready-map tests in `internal/status/status_signals_test.go`.
- The anchor seam is the registry mutation table and the `workflow-guidance-anchors` fixtures. Prior art: the `decision-map-asset-path` fixture.
- The gate observes the feature through the `decision-map-integrity` check, the `workflow-guidance-anchors` family, and `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.

### Seam diagram

    trigger: bench maps [<map>], bench commit, bench gate
        │
        ▼
    decisions/<topic>.md + <topic>/tickets/*.md  ──▶  [ maps: discover, parse, grade ]  ──▶  diagnostics, rows, counts
                                                          ◀ tests attach here: temp tree in, diagnostics or rows out
        │
        ▼
    rows  ──▶  [ maps.Command: filter, ready, stale, help ]  ──▶  TOON table + help + exit code
                    ◀ tests attach here: gittest.Repo with committed files, output compared to golden bytes

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| DS1 | 1, 2 | An index with Notes, Decisions so far, the four terminal sections, and `tickets/1.md` validates clean | `ValidateDecisionMapTree` on a temp tree | A parser that still wants inline tickets reds the split map |
| DS2 | 3 | The ticket id is the basename `7`, the title is the file's first line, and `Blocked by: #7` resolves | `ValidateDecisionMapTree` on a temp tree | An id read from the title leaves `#7` dangling |
| DS3 | 4 | An index with `## #2: Old` reds with `inline ticket #2: move it to alpha/tickets/2.md` | integrity fixture `inline-ticket-heading` | A parser that accepts the heading passes it |
| DS4 | 5 | `tickets/07.md`, `tickets/0.md`, and `tickets/a.md` each red naming the file | integrity fixture `ticket-basename` | A parser that trims zeros or accepts any name passes them |
| DS5 | 6 | A ticket file without `Type:` reds with `ticket #1: missing Type`, and `Type: Guess` reds as unsupported | `ValidateDecisionMapTree` on a temp tree | A ticket parse that drops the field rules passes both |
| DS6 | 7 | An index with no tickets folder reds with `missing decision ticket`, and an empty folder reds the same | integrity fixtures `tickets-absent` and `tickets-empty` | A parser that treats no tickets as zero rows passes both |
| DS7 | 8 | `decisions/beta/tickets/1.md` with no `decisions/beta.md` reds with `beta/tickets: no map index at beta.md` | integrity fixture `orphan-tickets-folder` | A scan that reads only index files never sees the folder |
| DS8 | 9 | An index without `## Notes` reds with `missing Notes section` | integrity fixture `notes-missing` | An optional Notes passes the omission |
| DS9 | 9 | An index without `## Decisions so far` reds with `missing Decisions so far section` | integrity fixture `decisions-so-far-missing` | An optional index section passes the omission |
| DS10 | 10 | A resolved #3 with no gist reds with `ticket #3: resolved without a gist in Decisions so far` | integrity fixture `gist-missing` | A check that reads only listed gists misses the absent one |
| DS11 | 11 | A gist to `tickets/4.md` while #4's Answer is `— (open)` reds with `Decisions so far links unresolved ticket #4` | integrity fixture `gist-unresolved` | A check that tests only file presence passes it |
| DS12 | 12 | A gist to `tickets/9.md` with no such file reds with `Decisions so far links missing ticket #9` | integrity fixture `gist-missing-file` | A check that trusts the link passes it |
| DS13 | 12 | A Decisions so far line without the link grammar reds with `Decisions so far line has no ticket link` | integrity fixture `gist-malformed` | A lenient list parse accepts prose as a gist |
| DS14 | 13 | `Blocked by: #5` with no `tickets/5.md` reds with the dangling-blocker diagnostic | integrity fixture `graph-dangling` rebuilt on the split shape | A graph built from gists misses the edge |
| DS15 | 14 | Each graph, readiness, source, and terminal-list fixture rebuilt on the split shape reds with its EXPECT | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A rule dropped in the move leaves one fixture green |
| DS16 | 15 | A ready split map under `specs/x/decisions/` validates, and `Status: shaping` there reds with `compiled map must be ready` | `ValidateDecisionMapTree` on a temp tree | A ticket scan that runs only under `decisions/` misses compiled tickets |
| DS17 | 16 | `bench maps --template` output holds `## Notes` and `## Decisions so far` and no `## #1:` line | `maps.Command` | An unchanged template still prints the inline ticket |
| DS18 | 16 | `bench maps --ticket-template` output starts with `# ` and holds `Blocked by: none`, `Type: Research`, `### Question`, and `### Answer` | `maps.Command` | A missing flag exits 2 and prints no skeleton |
| DS19 | 16 | `bench maps --template --ticket-template` exits 2 with the mutual-exclusion help line | `maps.Command` | A parse that takes the first flag exits 0 |
| DS20 | 17 | The index skeleton holds `A map-owned asset stays in the map's assets folder, decisions/<topic>/assets/.` | `maps.Command` and anchor fixture `decision-map-asset-path` | A template with the flat-folder sentence passes an unchanged anchor |
| DS21 | 18 | `decisions/my map.md` with `decisions/my map/tickets/1.md` validates with no diagnostic | `ValidateDecisionMapTree` on a temp tree | A path join that splits on spaces loses the folder |
| DS22 | 19 | Every pre-existing row of the frontier-plus-invalid golden keeps its five cells and gains a sixth cell with its file path, and no row is added, dropped, or reordered | `internal/maps/maps_command_test.go` (`TestCommandAppendsOnlyMapActionsToTheCapturedPrimaryResponse`) against the regenerated golden file | Any dropped or reordered row changes the bytes |
| DS23 | 20 | A ready active map projects `[name, title, map, ready, "", decisions/name.md]` and the help line `/bench-write-spec decisions/name.md` | `maps.Command` on a gittest repo | A projection that skips ready maps prints no row |
| DS24 | 20 | A ready compiled map under `specs/` projects no row | `maps.Command` on a gittest repo | A scan that lists compiled maps adds a row the status count lacks |
| DS25 | 21 | The header is exactly `maps[N]{map,title,type,state,blockers,path}` and every row has six cells | `maps.Command` golden bytes | A five-cell row fails the typed table render |
| DS26 | 22 | `bench maps alpha` with maps alpha and beta prints only alpha's rows and help lines | `maps.Command` on a gittest repo | An unfiltered scan prints beta's rows |
| DS27 | 23 | `bench maps gamma` with no such map exits 1 and prints `maps: no active map named "gamma"` | `maps.Command` on a gittest repo | A filter with zero rows exits 0 with an empty table |
| DS28 | 24 | `bench maps decisions/alpha.md`, `bench maps alpha.md`, and an operand holding `\x1b` each exit 2 | `maps.Command` | A filter that trims accepts the path form |
| DS29 | 25 | With the asset committed after the index, the output holds `[alpha, decisions/alpha/assets/r.md, source, stale, "", decisions/alpha/assets/r.md]` | `maps.Command` on a gittest repo with two commits | A compare on the wrong file or direction prints no row |
| DS30 | 25 | With the index committed after the asset, no stale row prints | `maps.Command` on a gittest repo with two commits | A compare that flags every cited asset prints a row |
| DS31 | 26 | An uncommitted asset and an uncommitted index each print no stale row | `maps.Command` on a gittest repo | A missing date treated as zero flags the asset or the map |
| DS32 | 27 | A stale row leaves the exit code 0, and `--count` prints the same count with and without the stale asset | `maps.Command` on a gittest repo | A stale row counted as unresolved raises the count or the code |
| DS33 | 28 | `appendMaps` returns `1 ready map(s)` with `/bench-write-spec decisions/ready.md` for one ready split map and `1 unresolved map(s)` for one shaping split map | `internal/status/status_signals_test.go` (`TestAppendMapsRoutesReadyOnlyWithoutUnresolvedOrInvalidMaps`) on split maps | A count that reads only inline maps returns zero |
| DS34 | 29 | The dry run prints every map, ticket, asset move, reference file, and the ignore line, and changes no file | review-owned: the migration ticket records the plan and a clean `git status` | A program that applies on a dry run leaves a dirty tree |
| DS35 | 30, 31 | After `--apply`, every index has Notes, one gist per resolved ticket, and no `## #n:` line, and each ticket file validates | `ValidateDecisionMapTree` over the migrated tree through `bench maps` | A migration that leaves one inline ticket reds the contract's parser |
| DS36 | 32 | Every Sources Path names a file under that map's assets folder, and no asset reference under `decisions/`, `specs/*/decisions/`, or `docs/` names a file in `decisions/assets/` | review-owned: the migration ticket records both searches | A rewrite that misses one reference leaves a hit |
| DS37 | 33 | `docs/research/ft191-resolved-reader-research.md` is tracked and `decisions/assets/` does not exist | review-owned: the migration ticket records `git ls-files` | An orphan left in place keeps the flat folder alive |
| DS38 | 34 | `git check-ignore docs/research/x.md` exits 1 and `git check-ignore research/x` exits 0 | review-owned: the migration ticket records both runs | A rule left as `research/` still ignores the docs path |
| DS39 | 35 | The `bench maps` rows before and after the migration are equal once the `path` column is dropped | review-owned: the migration ticket records both outputs and their diff | A ticket whose answer moved wrong changes its state |
| DS40 | 36 | After the contract, `scripts/split-decision-maps/` does not exist and `internal/maps` parses no `## #n:` heading | review-owned: the contract ticket records `ls` and `rg` | A kept program or parse survives the contract |
| DS41 | 37 | Every fixture under `tests/canary/decision-map-integrity/` bites with its EXPECT | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A rebuilt fixture whose mutation targets no file stays green |
| DS42 | 38 | The inventory names each new fixture, and deletion of any one reds the inventory test | `internal/conformance/decision_map_integrity_test.go` (`TestDecisionMapIntegrityFixtureInventoryRejectsDeletion`) | A fixture added without an inventory entry can vanish |
| DS43 | 39, 40, 41 | The shape-idea command carries the eight anchored sentences listed under Further notes, and each mutation-table row bites | `TestDecisionMapSplitAnchorsRedOnRemoval` and the `workflow-guidance-anchors` fixtures | A reworded rule with no anchor passes the docs check |
| DS44 | 42 | The write-spec command carries the topic-folder move sentence, and the old `decisions/assets/` needle is forbidden there | `TestDecisionMapSplitAnchorsRedOnRemoval` | A command that keeps the flat-folder sentence passes an unchanged anchor |
| DS45 | 43 | `CONTEXT.md` carries the four glossary entries with an Avoid list each, and the gist entry has an anchor | `TestDecisionMapSplitAnchorsRedOnRemoval` | A glossary that drops an entry passes without the anchor |
| DS46 | 44 | `docs/adr/0020-a-decision-lives-in-its-ticket-file.md` exists, records the three decisions, and names no path | review-owned: the docs ticket cites the file | An ADR with paths rots |
| DS47 | 45 | `README.md` keeps `Decision maps are situational` and describes the index and ticket files | the README anchors in the `workflow-guidance-anchors` family | A README rewrite that drops the anchored phrase reds |
| DS48 | 46 | Every anchor row changed in this spec has a mutation-table row that bites | `TestDecisionMapSplitAnchorsRedOnRemoval` | A needle without a mutation row is a claim, not a bite |
| DS50 | 18, 21 | `decisions/my map.md` with `decisions/my map/tickets/1.md` projects one row whose path cell is `decisions/my map/tickets/1.md` | `maps.Command` on a gittest repo | A path join that splits on spaces loses the folder |
| DS51 | 12 | A gist to `other/tickets/1.md` on map `split` reds with `Decisions so far links missing ticket #1` | integrity fixture `gist-wrong-folder` | A parser that reads only the number accepts a link to another folder |
| DS52 | 10 | Two gists to `split/tickets/1.md` red with `Decisions so far duplicate gist for ticket #1` | integrity fixture `gist-duplicate` | A set keyed by number drops the second gist in silence |
| DS53 | 45 | The README sentence `each decision lives in one ticket file under the map's tickets folder` reds when removed | `TestReadmeSplitShapeAnchorRedsOnRemoval` and the README anchor row | A README rewrite that drops the split description passes an unchanged anchor |

Not covered: story 47 — the measure is recorded in the phase-close retro by final-check, which no ticket owns.

### Edge inventory

- A tickets folder holds a non-`.md` file, a dotfile, or a README: ignored, as the direct-child policy ignores them today.
- A nested folder under `tickets/`: **Won't handle** — the direct-child scan is the surviving caller, and DS1 covers it.
- A ticket file that is a symlink or unreadable: `bounds.Classify` reports it, and the map is invalid. That path is unchanged.
- A Sources URL source: **Won't handle** a stale row — no commit date exists for a URL, and DS29 covers the Path source.
- A Sources Path that does not exist: **Won't handle** as a stale row — the existing validity check reds it, and DS29 covers the live path.
- A gist whose link text differs from the ticket title: **Won't handle** — the lane checks the link target, and DS12 covers it.
- Section order in the index: **Won't handle** — the parser tolerates any order today, and DS1 covers the valid shape.
- A map name with a space: DS21.
- An operand with a control byte: DS28. A map name with a control byte already fails the typed table render, unchanged.
- An absent `decisions/` folder stays the parsed-empty state. An empty tickets folder reds, per DS6. An absent `specs/` folder stays silent.
- Both audiences, this repository and a linked repository, get the same answer, because the kit ships one parser.
- Two concurrent `bench maps` runs: read-only, no lock needed.
- An untracked map or asset in the stale compare: DS31.
- A `tickets` folder that is a symlink: **Won't handle** — `bounds.ClassifyDir` follows it as the `decisions/` scan does today. The review recorded the edge for reviewer veto.
- A tickets folder plus an inline heading during the expand phase: **Won't handle** — the migration commit lands before the contract, and DS3 closes the window.
- A linked repository with inline maps after the contract: **Won't handle** — no linked repository holds a map on 2026-09-06, and DS40 covers the deletion.
- The FT99 premise sentence: **Won't handle** a mechanical check — `bench preflight build` owns premise checks, and DS43 covers the prose.

## Ownership fences

- `internal/maps/`
- `internal/status/status_signals_test.go`
- `internal/status/status_producible_test.go`
- `internal/status/status_command_test.go`
- `internal/conformance/decision_map_integrity_test.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `scripts/split-decision-maps/`
- `decisions/`
- `specs/decision-map-split/decisions/`
- `docs/research/`
- `.gitignore`
- `tests/canary/decision-map-integrity/`
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/skills-index-command-adapters/`
- `tests/canary/docs-currency-token-diet/`
- `tests/canary/load-validity-metadata/`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `.agents/commands/bench-shape-idea.md`
- `.agents/commands/bench-write-spec.md`
- `CONTEXT.md`
- `README.md`
- `docs/field-guide.html`
- `docs/adr/0020-a-decision-lives-in-its-ticket-file.md`
- `reviews/decision-map-split.md`
- `capture/retros/`
- `internal/axi/action.go`
- `internal/axi/action_test.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/conformance/registry_test.go`

## Build decisions

Each entry below is a call the build made under the batch approval of 2026-09-06. The reviewer can veto any one of them at the review.

- The `graph-duplicate-id` fixture retired at the contract. A ticket id is a file basename, so two tickets cannot share one id, and the duplicate-id branch is dead. DS15 now covers the five remaining graph fixtures.
- `TestMapGraphRejectsInvalidEdges` folded into `TestDecisionMapDiagnosticsGolden`, which asserts the same faults as one ordered slice on the split shape.
- The status helpers for a split map moved into `status_command_test.go`. The two signal test files are over budget, so the growth rule forbids a line there.
- The fence gained `internal/axi/action.go` and its test. The ready row's help action `/bench-write-spec <path>` needs the phase in the harness-phase set with a path argument. `bench status` already renders that action through its own table.
- The decision-map anchor rows and their mutation tables moved into `registry_decision_maps.go` and its test. Both registry files are over budget with no grant, and a grant is the reviewer's decision, so the ticket moved its headroom instead. The family owner set in `internal/conformance/registry_test.go` names the new file too.
- The migration program moved `gate-pipeline-fixture-inventory.md` into the `gate-pipeline` assets folder. Only that map names the asset, in prose rather than in Sources.
- Two seeded gists read thin and wait for a reviewer edit: `gate-critical-path` #1 and `worktree-orphan-retirement` #5.

**Migration evidence for DS34 and DS39, 2026-09-06.** The dry run printed fourteen maps with 128 ticket files, ten asset moves, thirteen reference files, and the ignore line. After the dry run, `git status` showed only the program folder. After `--apply`, `bench maps` printed the same nine rows and three help lines as the coordinator's saved pre-migration output, byte for byte. The `path` column did not exist yet.

The per-map counts of resolved tickets and gists were equal for all fourteen maps. The program is deleted, so the record here is the evidence.

**Review repairs, 2026-09-06.** The review accepted eight repair targets. Rows DS51 to DS53 record the three that add a check. The prose repairs keep every anchored sentence byte-identical. Six items stay open in the pickup for the reviewer.

**Dogfood runs before the review, 2026-09-06.** The coordinator ran the worktree binary over the real tree at the last ticket's tip. `bench maps` printed fourteen rows with six cells, five `ready` rows, and eight help actions. `bench maps software-factory` printed that map's two rows and one action. `bench maps --ticket-template` printed the ticket skeleton.

`bench status` kept the decisions row at eight unresolved maps. No stale row printed, because the migration committed each asset beside its index.

## Out of scope

- A tracker-backed map: 0 edits here, its own decision first.
- A body reader `bench maps show <topic> <n>`: 3 edits, 1 gate run, own spec.
- A gate red on a stale Source: 2 edits, 1 gate run, rejected by the reviewer.
- A kit-shipped migration verb for linked repositories: 4 edits, 2 gate runs, own spec if a linked repository ever holds an inline map.
- The FT304 shared observation view: consumes this projection, own spec.
- The craft-research skill and the report contract: the sibling spec `craft-research-skill`.

## Further notes

**Flagged additions beyond the decision source.**

- `--ticket-template` is a new flag, because the schema owner must render the second file shape.
- The `path` column is the column the reviewer's zoom decision needs.
- The orphan-folder diagnostic (DS7) and the basename rule (DS4) are edge cases the split creates.
- The refusal grammar for the map operand (DS27, DS28) is the hostile-input checklist applied to the new operand.
- The missing-Sources-Path deviation is recorded under Implementation decisions.
- The ignore-rule narrowing moves here from spec 2, because the orphan move needs a trackable `docs/research/`. Map #14 assigned it to spec 2, so this is a reassignment for reviewer veto.
- A `ready` row projects for an active map only, because the status count reads active maps only.
- The linked-repository premise is reviewer-supplied and not checkable in this tree.

**Source sentence to row.**

- Map #11, index and ticket file: DS1, DS2, DS3.
- Map #11, Notes and Decisions so far: DS8, DS9.
- Map #11, worktree lease claim: DS43.
- Map #11, migration of every active map: DS35.
- Map #11, parser accepts only the split shape: DS3, DS40.
- Map #11, lane reds the three drift defects: DS10, DS11, DS12, DS14.
- Map #11, resolved means the Answer is not a placeholder: DS11.
- Map #11, refer by name: DS43.
- Map #11, CONTEXT and ADR: DS45, DS46.
- Map #12, frontier owner and map argument: DS22, DS26.
- Map #12, ready row: DS23.
- Map #12, no body reader: Out of scope.
- Map #12, status consumes the projection: DS33.
- Map #13, assets folder: DS36.
- Map #13, `docs/research/` trackable: DS37, DS38.
- Map #13, orphan move: DS37.
- Map #13, stale row by git dates, advisory: DS29, DS30, DS31, DS32.
- Map #13, a missing Source path as a stale row: deviation recorded under Existing rules unchanged, the validity red stays.
- Map #14, spec 1 contents: every group above.
- Map #14, FT99 rule: DS43.
- Map #14, measurement: the Not covered line for story 47.

**Anchored sentences the commands ticket and the docs ticket register.**

- In `bench-shape-idea.md`: `The map is an index: it lists the decisions made and links the ticket that holds each one.`
- In `bench-shape-idea.md`: `A decision ticket is one file under the map's tickets folder, named by its number.`
- In `bench-shape-idea.md`: `The answer lives only in the ticket file.`
- In `bench-shape-idea.md`: `Decisions so far holds one gist line per resolved ticket, with a link to its file.`
- In `bench-shape-idea.md`: `Notes holds the domain, the skills a session consults, and the standing preferences of that map.`
- In `bench-shape-idea.md`: `The shaping worktree lease is the claim, and no owner field enters a ticket.`
- In `bench-shape-idea.md`: `Name a ticket by its title, with its number beside it.`
- In `bench-shape-idea.md`: `A grill recommendation that asserts current-code behavior names the evidence read in the current session.`
- In `bench-write-spec.md`: `Move the topic folder, its tickets and assets included, into the spec's decisions folder as one unit.`
- In `CONTEXT.md`: `one line in Decisions so far that links a resolved decision ticket`.
- Forbidden after the move, in both commands and the template constant: `decisions/assets/`.

**Reviewer sign-off, 2026-09-06.** The reviewer approved the ignore-rule move into the migration ticket, the missing-Sources-Path validity red, the fences as written, and the seven-ticket graph.

**Cheapest wrong implementation per group.**

- Index and ticket files: a parser that keeps the inline tickets and ignores the folder. DS3 and DS1 red it.
- Projection: a filter that returns an empty table for an unknown name. DS27 reds it.
- Migration: hand-edited maps with no program. DS34 is review-owned, so the migration ticket's return cites the program source file and the dry-run output.
- Guidance: a reworded rule with no needle. DS43 and DS48 red it.

**Pre-review proof checklist.**

- Cited symbols: `ValidateDecisionMapTree`, `ValidateDecisionMap`, `ParseDecisionMap`, `DecisionMapTemplate`, `Command`, `ActiveCounts`, `appendMaps`, `bounds.Classify`, `git.Output`, `FieldScan`, and `GraphWalk` resolve in the tree read on 2026-09-06.
- Import edges: `internal/status` imports `internal/maps`. `internal/tickets` imports `internal/maps` for the generic scan only. `internal/conformance` binds the check in its test table.
- Source-row clauses and occurrences: listed under Source sentence to row.
- Promised field labels: `## Notes`, `## Decisions so far`, `Blocked by:`, `Type:`, `### Question`, `### Answer`, `--ticket-template`, column `path`, states `ready` and `stale`, type `source`.
- Changed-function callers: `DecisionMapTemplate` is called by the `internal/status` tests, `internal/conformance/decision_map_integrity_test.go`, and `cmd/bench/command_registry_test.go`. `Command` is called by `cmd/bench/main.go`. `ActiveCounts` and `DiscoverDecisionMapCandidates` are called by `internal/status/status.go`.
- Copy survival: DS40 reds a surviving inline parse or program.

**Reader sweep.**

- `internal/maps` owns the artifact.
- `internal/status/status.go` reads counts and the ready action.
- `internal/gate/lane_select.go` classifies the path, unchanged.
- `internal/conformance` binds the check and holds the fixture inventory.
- `internal/spec/spec.go` removes the folder at retirement, unchanged.
- `internal/anchors/registry_data.go` holds the needles.
- The two command files, `CONTEXT.md`, `README.md`, and `docs/field-guide.html` describe the shape. The ADR skill's sentence that maps live in `decisions/` stays true and is not edited.
- The canary fixture copies of the two commands repeat the prose.
- `bench handoff` reads maps only through the status board. `bench roadmap` does not read them. No script or workflow reads `decisions/`.
- Shipped-surface claim words: none added.
