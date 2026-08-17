# Progressive roadmap

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-17 (`/bench-shape-idea 198`, three grill rounds and a confirmed shared understanding, recorded in `capture/session-handoff.md` at commit `handoff: FT198 shaped`). Supersedes decision #1 of the retired `roadmap-progressive-index` compiled map (2026-08-13, "`ROADMAP.md` remains the only durable owner"): the FT198 row's redrain reopened the detail-owner question and the 2026-08-17 grill closed it as one file per row.

Verification log: 1 iteration to accept — one round over the pair with `opus`/high through the native agent surface. BLOCK on 6 blocking + 8 prose findings, all folded: `docs-currency-workflow` parses the real board's ledgers (parser and migration merged into one ticket, the ledger check and retired-verb sweep move onto the loader), fences and `Writes:` completed (`docs_workflow_checks_test.go`, `tier_test.go`, `command_registry_test.go`, `projects/benchkit.md`, the retro path), `roadmap/` joins the sequence-trust source list with a degraded-file row, PR15 restated to the canary restore contract with the fixture-overlay constraint, per-class row disposition decided, line-neutral prose edits under two zero-headroom budgets, the group-1 effort declared as a bump, the row-ID vocabulary kept, the malformed-ledger source path decided, the dashboard row marked a regression pin, spec-writer additions flagged, and PR14's cells named.

## Problem

`ROADMAP.md` is 93 KB and is read whole every time it is read at all. Every drain,
cold shaping session, and status question loads sixty-seven row bodies to reach the
three rows it needs, and every row edit — an occurrence key, a body sentence, a
retirement — is an edit inside one 1,400-line document. The CLI already projects a
body-less index (`bench roadmap --context`, schema 4, 18 KB) and a per-row detail
fetch, but the file layout underneath still forces whole-document reads and writes on
anyone holding a file tool, and nothing mechanical says a row's detail is present.

## Solution

Split the board into an index and one detail owner per row. `ROADMAP.md` keeps its
section headings, board-level prose, and — per row — exactly one physical heading
line, `**FT<n> (<qualifier>) — <title>.**`, with no body. Each row's body, its
`Occurrence:` ledger, and its `Sources:` line move to `roadmap/FT<n>.md`, whose first
line repeats that heading verbatim. The roadmap package parses only this shape: an
index row with no file, a file with no index row, an inline body under an index line,
a file whose first line differs from its index line, an unrecognized file under
`roadmap/`, a duplicate index ID, and a heading wrapped across physical lines are each
one named diagnostic. A conformance check returns those diagnostics, so the gate is
red on any of them, and `bench roadmap --context` reports them as parse failures with
the roadmap source marked malformed. This repo's board is migrated once by a script
that is never committed; the migration is accepted when the schema-4 rows and
sequence blocks are identical before and after. There is no compatibility path for
the inline-body shape: this is the only repository that carries a roadmap.

## User stories

**Reading and grading the split board** — Line: `opus` / high. Gate and conformance
logic routes to mid effort in the profile; high is a declared bump because this
rewrite replaces the one parser every drain, status read, and gate check grades
against.

1. As a cold session, I want `ROADMAP.md` to carry each row as one physical heading line `**FT<n> (<qualifier>) — <title>.**` and no body, so that the whole board is one small file.
2. As a maintainer editing one row, I want its body, `Occurrence:` ledger, and `Sources:` line in `roadmap/FT<n>.md`, so that a row edit opens one file.
3. As a reader of one row file, I want its first line to be the row's index line verbatim, so that the file reads alone without the index.
4. As a maintainer, I want an index row with no `roadmap/FT<n>.md` reported as `roadmap/FT<n>.md: missing detail owner for ROADMAP.md row FT<n>`, so that a lost detail owner is named, not silent.
5. As a maintainer, I want a `roadmap/FT<n>.md` with no index row reported as an orphan naming the file, so that a stale detail file cannot linger unindexed.
6. As a maintainer, I want any non-blank line between an index line and the next index line, `## ` heading, or end of file reported as an inline body naming the row, so that bodies cannot creep back into the index.
7. As a maintainer, I want a row file whose first line differs byte-for-byte from its index line reported as a heading mismatch naming the file and the row, so that the two copies cannot drift.
8. As a maintainer, I want a file under `roadmap/` whose basename is not `<row ID>.md` under today's row-ID grammar reported as unrecognized, so that stray files are loud.
9. As a maintainer, I want an ID that appears on two index lines reported as duplicated, so that a row cannot hold two positions.
10. As a maintainer, I want an index heading whose `**` close is not on the same physical line reported as malformed, so that the index line grammar is exactly one line.
11. As a maintainer, I want the occurrence ledger parsed from the row-file body under today's rules, so that occurrence keys and counts are unchanged and a malformed ledger is reported against `roadmap/FT<n>.md`.
12. As a maintainer, I want spec slug, spec status, and the external-trigger words derived from the heading plus the row-file body, so that every schema-4 column keeps the value it has today.
13. As an agent, I want `bench roadmap --context` to list `roadmap/` in its sources block with state and byte total, so that the index reports every file it read.
14. As an agent, I want each integrity diagnostic rendered as a `parse_failures` row naming the offending path, with `ROADMAP.md`'s source state `malformed` and `sequence_trusted` false, so that the index cannot look clean over a broken tree.
15. As an agent, I want `bench roadmap --context --row FT<n>` to return the row-file body as the row's `body`, so that the detail fetch keeps its shape.
16. As an agent, I want `bench roadmap` to keep its exit-1 postures (empty file, failed read, unsupported schema) and render its top rows from the split tree, so that the porcelain is unchanged.
17. As a maintainer, I want a Dev-tier `roadmap-detail-integrity` conformance check registered against the repo root that returns every diagnostic in stories 4–10, so that the gate is red on any of them.
18. As a maintainer, I want one canary fixture per diagnostic class, each with `EXPECT` and `MUTATE.json`, so that the check is proven to bite for its own planted reason.
19. As an agent running `bench idea --owner FT<n>`, I want the owner validated against the split tree, so that an untrusted tree or an absent owner is refused exactly as today.
20. As a maintainer, I want `bench status`'s roadmap-reconcile signal to classify spec paths found in row files, so that merged and dangling counts do not fall to zero after the split.
21. As a dashboard reader, I want the roadmap page and its recommended sequence rendered from `ROADMAP.md`'s index text, so that the dashboard keeps working without bodies.
22. As a maintainer of a fresh repo, I want an absent `ROADMAP.md` with no `roadmap/` directory to stay the quiet posture, so that a repo without a board is not red.
23. As a maintainer, I want an absent `ROADMAP.md` beside a non-empty `roadmap/` reported as one orphan per file, so that detail without an index is loud.
24. As a maintainer, I want a row file holding only its heading accepted as a row with an empty body, so that a freshly drained row can start bare.
25. As a maintainer, I want the stripped-distribution system journey to remove `roadmap/` alongside `ROADMAP.md`, so that the stripped subject stays representative.
26. As a maintainer, I want the schema-4 snapshot's schema number unchanged, so that `/bench-what-next`'s entry contract does not move for an additive sources row.
44. As an agent, I want a row file the classifier cannot read (oversized, wrong type, unreadable) to mark `sequence_trusted` false and render a parse failure naming that file, so that an unread ledger is never trusted.
45. As a maintainer, I want the docs-currency occurrence-ledger check and the retired-verb sweep to read the loader and the row files, so that neither gate check loses its corpus when bodies move.
46. As a maintainer, I want every faulted row's disposition fixed per class — index-side faults keep the row, file-side and structural faults drop it — so that `rows_total`, `--row`, and status agree on what a faulted board contains.

**Migrating this board** — Line: `opus` / medium. Mechanical once the parser
exists, but it rewrites the whole board and the acceptance is a differential the
builder must run and record.

27. As a maintainer, I want this repo's sixty-seven rows split by a one-shot script that is never committed, so that no migration code survives to be maintained.
28. As a maintainer, I want every wrapped index heading joined onto one physical line and every inline body moved into its row file, so that the migrated index parses under story 10 and story 6.
29. As a maintainer, I want the migrated tree's `bench roadmap --context --full` `roadmap_rows` and `sequence` blocks byte-identical to the pre-migration output, so that no row, body, ledger, or sequence line is lost or altered.
30. As a maintainer, I want the differential evidence — the two captures and the empty diff of those blocks — recorded in the migration commit message, so that acceptance is auditable after the script is gone.
31. As a maintainer, I want the section headings, section-intro prose, release-readiness, dependencies, and recommended-sequence text kept in `ROADMAP.md` verbatim, so that the board's judgment stays where it was.

**Documenting the split** — Line: `fable` / high. Guidance prose compounds through
every session that loads it while the edit costs few tokens — the profile's
doc-authoring leverage override.

32. As a cold session, I want `CONTEXT.md`'s **roadmap** and **roadmap index / roadmap detail** entries to name the index file and the `roadmap/FT<n>.md` detail owner in the same green change, so that the glossary never describes a tree that does not exist.
33. As a drain session, I want `/bench-what-next` to say that row bodies and ledgers are edited in `roadmap/FT<n>.md`, ordering in `ROADMAP.md`, and that retiring a row deletes both, so that the drain edits the right file.
34. As a spec author, I want `/bench-write-spec`'s promote-then-delete step to remove the row's index line and its `roadmap/FT<n>.md`, so that a shipped row leaves no orphan.
35. As a shaping session, I want `/bench-shape-idea`'s cold entry to read the index and fetch detail per row, so that shaping stops paying for sixty-seven bodies.
36. As a reader of the kit's own docs, I want `.bench/BENCH.md`'s capture paragraph, `projects/benchkit.md`, `.bench/BENCH-reference.md`'s file map, and `README.md`'s tree line to name `roadmap/` as the detail owner, so that the file map is true.
37. As a maintainer, I want the `/bench-what-next` detail-owner sentence anchored with a canary that fails for its own planted reason, so that the drain instruction cannot silently regress.
38. As a maintainer, I want a CHANGELOG entry for the split, so that the change is discoverable.
39. As a maintainer, I want `bench skills-index --check` green after the command prose changes, so that the index and the commands cannot disagree.

Spec-writer additions beyond the decision source, flagged for reviewer veto:
stories 8 and 44 (integrity classes the one-file-per-row layout implies), 25
(stripped journey), 38–39 (CHANGELOG, skills-index), 45–46 (found by review
against the enforcement surface).

**Explicitly not wanted** (see Out of scope)

40. As a maintainer, I do not want the inline-body shape accepted beside the split shape, so that there is one parser and one layout.
41. As a maintainer, I do not want `ROADMAP.md` generated from row-file fields, so that priority order stays a hand-authored judgment.
42. As a maintainer, I do not want a retained `bench roadmap migrate` verb, so that no linked-repo migration path is implied.
43. As a maintainer, I do not want a new `--check` flag, so that the conformance check and the `--context` failures block remain the two surfaces.

## Implementation decisions

- **One tree parse, pure at the seam.** `ParseDocument(content, statuses, full)`
  becomes a tree parse taking the index bytes and the classified `roadmap/`
  directory listing (basename → bytes and classifier state) and returning the same
  `Document`, parse failures, and a new ordered list of integrity diagnostics. The
  filesystem read stays in one loader used by `bench roadmap`, `--context`,
  `bench idea --owner`, `bench status`, the dashboard, and the conformance check, so
  no caller re-derives the layout. `RoadmapFile` gains a sibling `RoadmapDir =
  "roadmap"` constant in the same package.
- **Index line grammar is one physical line.** `^\*\*<ID>[^*]*\*\*\s*$` where `<ID>`
  is today's row-ID grammar (`[A-Za-z]+[0-9]+`, shared with `--row`), unchanged.
  Today's parser joins wrapped headings and accepts trailing inline text; the split
  shape accepts neither, and each is its own diagnostic (stories 10 and 6). Title
  extraction — trim ` —:-` and whitespace between ID and close — is unchanged, so a
  wrapped heading joined by the migration projects the same `title` cell it does
  now. Row files are `roadmap/<ID>.md`; this board's IDs are all `FT<n>`.
- **Faulted-row disposition, per class.** Missing detail owner, heading mismatch,
  and inline body keep the index row (title from the index, body from the file when
  present, inline text discarded) so `rows_total` and status still see it. Wrapped
  heading yields no row; the second of a duplicate ID yields no row; orphan and
  unrecognized files yield no row. A malformed ledger reports its discrepancy and
  parse failure with `Source` = `roadmap/<ID>.md`, capture unit = the ID.
- **`roadmap/` joins the sequence-trust sources.** A row file whose classifier state
  is anything but parsed or empty renders a `parse_failures` row naming it and
  flips `sequence_trusted` false, exactly as a degraded `ROADMAP.md` does.
- **Two docs-currency readers move onto the loader.** `checkOccurrenceLedgerMigration`
  reads the loader's `Document` (its count map unchanged) and
  `checkRemovedVerbSweep` sweeps `README.md`, `ROADMAP.md`, and every
  `roadmap/<ID>.md`. Both land with the parser and the migration in one commit — the
  docs-currency check grades the real board, so parser and migration are one
  indivisible green.
- **Row-file grammar.** Line 1 is the index line byte-for-byte; the body is every
  byte after the first line break, trimmed as today's body is trimmed. The ledger
  parser, `spec.LiveSpecSlugs`, and the external-trigger words read the heading
  plus body, exactly the text they read today, so schema-4 cells do not move.
- **Diagnostics are ordered and path-first.** Missing owner, orphan, inline body,
  heading mismatch, unrecognized file, duplicate ID, wrapped heading — one message
  per finding, index order then directory order, each beginning with the offending
  repo-relative path. The conformance check returns them verbatim; `--context`
  renders each as a `parse_failures` row whose `source` is that path, and marks
  `ROADMAP.md` malformed. Schema stays 4: a new `sources` row and new
  `parse_failures` rows are additive.
- **Gate attachment follows `decision-map-integrity`.** A `func(root string)
  []string` in the roadmap package is bound in the conformance table and the
  registry as `roadmap-detail-integrity` (Dev, SubjectRoot, its own input source
  naming `ROADMAP.md` and `roadmap/`), with a canary family of the same name whose
  fixture inventory the check's own test asserts.
- **`bench idea --owner`, status reconcile, dashboard read the loader.** The owner
  check refuses on any diagnostic (structurally untrusted) as it refuses on any
  parse failure today; status classifies the slugs the loader found in row files;
  the dashboard renders the index text.
- **Migration is a scratchpad script.** It reads `ROADMAP.md`, writes the joined
  heading line into the index and heading + body into `roadmap/FT<n>.md`, and is
  never added to the tree. Acceptance is the differential in story 29, captured with
  the pre-migration binary on the pre-migration tree and the post-migration binary
  on the post-migration tree.

## Testing decisions

- A good test drives the tree parse with in-memory index bytes and a map of row
  files and observes the `Document`, the diagnostics, and the rendered `--context`
  blocks; the conformance check is proven through its canary family, one fixture
  per diagnostic, each mutation making the check red for that reason alone.
- Seams: the tree parse in `internal/roadmap` (prior art: `ParseDocument` tests
  and `context_test.go`'s rendered-block assertions); the conformance binding
  (prior art: `decision_map_integrity_test.go` and its fixture inventory);
  `internal/status`'s reconcile test; the systemtest stripped journey; the
  anchors registry canary for the drain sentence.
- The gate observes the feature through the `test` phase (roadmap, conformance,
  status packages), the `system` phase (stripped journey), and — for this repo's
  own board — the `roadmap-detail-integrity` check grading the migrated tree.

### Seam diagram

    trigger: bench roadmap | --context [--row] | bench idea --owner | bench status | gate conformance
        │
        ▼
    ROADMAP.md bytes + roadmap/ listing  ──▶  [ roadmap tree parse ]  ──▶  Document + diagnostics
                                                   ◀ tests attach here: in-memory index + row-file map,
                                                     assert rows, ledgers, diagnostics, rendered blocks
        │
        ▼
    [ roadmap-detail-integrity check ]  ──▶  []string diagnostics  ──▶  gate red / green
                      ◀ tests attach here: canary fixtures with EXPECT + MUTATE.json

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PR1 | 1, 3, 24 | An index of one heading line per row plus a row file whose first line is that heading and whose body is empty parses to a row with the index title and empty body and no diagnostic | tree parse unit test | A parser that still requires or reads an inline body reports a diagnostic or a wrong title |
| PR2 | 2, 11 | A row file body carrying `Occurrence: 2026-08-14 FT189 …` and an `Occurrences: baseline-01, baseline-02` ledger projects `occurrence_count` 2 and the same keys `ParseDocument` projects today, and a descending ledger yields the same malformed-ledger discrepancy | tree parse unit test | A parser that reads the ledger from the index or skips the file loses the count |
| PR3 | 12 | A row file body naming `specs/foo/spec.md` and the word `scheduled` projects `spec` foo, its status, and `external_trigger` true | tree parse unit test | Deriving these cells from the index line alone yields empty spec and false trigger |
| PR4 | 4 | An index line `**FT7 (LOW) — x.**` with no `roadmap/FT7.md` yields the diagnostic `roadmap/FT7.md: missing detail owner for ROADMAP.md row FT7` | tree parse unit test | A parser that ignores absent files reports nothing |
| PR5 | 5, 23 | `roadmap/FT8.md` present with no `FT8` index line yields an orphan diagnostic naming `roadmap/FT8.md`, and the same file with `ROADMAP.md` absent yields the same diagnostic | tree parse unit test | A parser that only walks the index never sees the file |
| PR6 | 6 | An index line followed by the non-blank line `Body text.` before the next `**` line yields an inline-body diagnostic naming the row | tree parse unit test | Today's parser treats that line as the body and stays quiet |
| PR7 | 7 | A row file whose first line is `**FT7 (LOW) — y.**` under an index line `**FT7 (LOW) — x.**` yields a heading-mismatch diagnostic naming `roadmap/FT7.md` and `FT7` | tree parse unit test | A parser that trusts the filename never compares the two lines |
| PR8 | 8 | A file `roadmap/notes.md` yields an unrecognized-file diagnostic naming it, and `roadmap/FT7.txt` likewise, while `roadmap/AB1.md` under a `**AB1 …**` index line is a valid row | tree parse unit test | Accepting any basename lets a stray file pass, and narrowing to `FT` drops a valid ID |
| PR9 | 9 | Two index lines with ID `FT7` yield one duplicate diagnostic naming `FT7` and the second position | tree parse unit test | Today's parser appends both rows silently |
| PR10 | 10, 46 | An index heading `**FT7 (LOW) — a` on one line closed by `long title.**` on the next yields a wrapped-heading diagnostic and no row, while a missing-owner FT8 keeps its row with empty body and `rows_total` counts it | tree parse unit test | Today's parser joins the two lines into a row, and an undecided disposition lets `rows_total` drift between callers |
| PR11 | 13, 26 | `bench roadmap --context` on a split tree renders a `sources` row `roadmap/,parsed,<bytes>` and its `context` row still reads schema `4` | context command test | A snapshot that never lists the directory or bumps the schema fails the literal block |
| PR12 | 14 | `bench roadmap --context` over a tree with one missing detail owner renders a `parse_failures` row whose source is `roadmap/FT7.md`, `ROADMAP.md`'s sources state `malformed`, and `sequence_trusted` false | context command test | Rendering diagnostics without flipping the source state lets the index look clean |
| PR13 | 15 | `bench roadmap --context --row FT7` renders the row-file body as `body` and its byte count as `body_bytes` | context command test | Reading the body from the index yields an empty cell |
| PR14 | 16 | `bench roadmap` on a split tree renders `title` from the index line and `spec`, `occurrence_count`, and `occurrence_keys` from the row-file body, and an empty `ROADMAP.md` still exits 1 with the record error | roadmap command test | A command that reads only the index renders empty spec and zero occurrences |
| PR15 | 17, 18 | The conformance table, registry, and profile check-input row bind `roadmap-detail-integrity` at Dev tier and SubjectRoot, and every canary fixture in that family emits exactly its `EXPECT` diagnostic when mutated and no longer emits it after `RestoreMutationFixture` — each fixture's `files/roadmap/` carrying a file for every index row it keeps, since restore overlays the live board | conformance registry test and canary fixture inventory | An unregistered check or a fixture whose mutation does not red is the same silent gap the row exists to close |
| PR16 | 19 | `bench idea --owner FT7 text` on a tree whose `roadmap/FT7.md` is missing exits with `ROADMAP.md is structurally untrusted`, and on a tree with FT7 present appends the idea | idea command test | An owner check reading only the index accepts the untrusted tree |
| PR17 | 20 | `bench status`'s roadmap-reconcile row counts a merged spec named only in `roadmap/FT7.md`'s body | status reconcile test | A scan of `ROADMAP.md` alone counts zero |
| PR18 | 21 | The dashboard's roadmap text and recommended sequence render from the split tree's index | dashboard render test | Regression pin: both readers already parse the index only, so no code change is expected |
| PR19 | 22 | A root with neither `ROADMAP.md` nor `roadmap/` yields no diagnostic and the absent posture from `bench roadmap` | tree parse unit test and roadmap command test | A validator that requires the directory reds a fresh repo |
| PR20 | 25 | The stripped-distribution journey's excluded-path list includes `roadmap/` and the stripped subject carries no such directory | systemtest stripped journey | Pin: keeps the excluded-path list complete once the directory exists |
| PR21 | 27, 28, 29, 31 | After migration, `roadmap_rows` and `sequence` blocks from `bench roadmap --context --full` are byte-identical to the pre-migration capture, every index heading is one physical line, no index line is followed by body text, and the board's non-row prose is unchanged | migration differential and the roadmap-detail-integrity check on the migrated tree | Any dropped ledger key, altered body byte, or leftover inline body shows in the diff or reds the check |
| PR22 | 30 | The migration commit message carries the two capture paths and the empty-diff statement | commit message review | An unrecorded differential cannot be audited once the script is gone |
| PR23 | 32, 36 | `CONTEXT.md`'s roadmap entries, `.bench/BENCH.md`, `projects/benchkit.md`, `.bench/BENCH-reference.md`, and `README.md` name `roadmap/FT<n>.md` as the detail owner in the landing commit | docs-currency review read | A glossary or file map still describing one file is a false current-state claim |
| PR24 | 33, 37 | `.agents/commands/bench-what-next.md` carries the anchored sentence naming `roadmap/FT<n>.md` as where row bodies and ledgers are edited and both files as what retirement deletes, and its canary fixture fails for that planted removal | anchors registry canary | An unanchored sentence can be pruned by the next prose diet |
| PR25 | 34, 35 | `bench-write-spec.md`'s promote-then-delete step names both the index line and the row file, and `bench-shape-idea.md`'s cold entry names the index read and per-row fetch | docs-currency review read and existing anchor needle | The existing anchor keeps `removes the spec's ROADMAP.md row` intact; a rewrite that drops it reds today |
| PR27 | 44 | `bench roadmap --context` over a tree whose `roadmap/FT7.md` is a directory renders a `parse_failures` row sourced `roadmap/FT7.md` and `sequence_trusted` false | context command test | A trust list without `roadmap/` claims a trusted sequence over an unread ledger |
| PR28 | 45 | With the ledgers in row files, `docs-currency-workflow`'s occurrence-ledger check reports the same counts it reports today and its retired-verb sweep visits every `roadmap/<ID>.md` | conformance docs-currency tests | A check still reading `ROADMAP.md` bodies counts zero and the sweep loses sixty-seven bodies |
| PR26 | 38, 39 | `CHANGELOG.md` gains the split entry and `bench skills-index --check` exits zero after the command edits | skills-index check and review read | A stale index reds the check |

Not covered: story 40 — reviewed exclusion, no legacy shape is built.
Not covered: story 41 — reviewed exclusion, no generator is built.
Not covered: story 42 — reviewed exclusion, no verb is built.
Not covered: story 43 — reviewed exclusion, no flag is built.

### Edge inventory

- **Won't handle** a row file with a heading whose ID differs from its filename
  (`roadmap/FT5.md` opening `**FT6 …**`) as its own class — PR7's heading mismatch
  against FT5's index line already names the file and the row.
- **Won't handle** a `roadmap/` entry that is a subdirectory or symlink as its own
  class — the classifier's wrong-type state reports it as an unrecognized entry
  through PR8's path, and no in-scope caller creates one.
- **Won't handle** a row file over `bounds.ControlRecordLimit` specially — the
  classifier's existing bounded-read state renders as a parse failure exactly as an
  oversized `ROADMAP.md` does today.
- **Won't handle** a `## ` heading inside a row file — today's body never carries
  one, and the file's body is taken to end of file, so a heading is body text.
- **Won't handle** an index line with a non-`FT` ID as its own class — the row-ID
  vocabulary is unchanged, so `**AB1 …**` with `roadmap/AB1.md` is simply a row.
- **Won't handle** an index line whose qualifier is missing (`**FT7 — x.**`) —
  the title trim already accepts it, and no row in this board lacks one.

## Ownership fences

- `internal/roadmap/`
- `internal/conformance/checks_test.go`
- `internal/conformance/registry_test.go`
- `internal/conformance/docs_workflow_checks_test.go`
- `internal/conformance/tier_test.go`
- `cmd/bench/command_registry_test.go`
- `internal/conformance/roadmap_detail_integrity_test.go`
- `internal/conformance/registry/registry.go`
- `internal/status/status.go`
- `internal/status/status_test.go`
- `internal/dashboard/`
- `internal/systemtest/owner_test.go`
- `internal/anchors/registry_data.go`
- `tests/canary/roadmap-detail-integrity/`
- `tests/canary/workflow-guidance-anchors/`
- `ROADMAP.md`
- `roadmap/`
- `CONTEXT.md`
- `CHANGELOG.md`
- `README.md`
- `.bench/BENCH.md`
- `.bench/BENCH-reference.md`
- `projects/benchkit.md`
- `.agents/commands/bench-what-next.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/commands/bench-shape-idea.md`
- `specs/progressive-roadmap/`
- `capture/session-handoff.md`
- `capture/learnings.md`
- `capture/retros/progressive-roadmap.md`

## Out of scope

- **Legacy inline-body compatibility** — reviewer-excluded (2026-08-17); this is
  the only repo with a board. If a linked repo ever grows one, that is a separate
  capability: 2 edits, 2 gate runs.
- **A generated index** (`ROADMAP.md` rendered from row-file priority fields) —
  reviewer-excluded; ordering stays hand-authored. Separate capability: 4 edits,
  3 gate runs.
- **A retained migration verb** — reviewer-excluded. Separate capability: 3 edits,
  2 gate runs.
- **A `--check` flag on `bench roadmap`** — reviewer left to discretion; not built,
  because the conformance check and the `--context` failures block already expose
  every diagnostic. Separate capability: 2 edits, 2 gate runs.
- **A `bench roadmap --context --row` fetch that reads the row file without parsing
  the index** — the fetch keeps parsing the whole tree; a file-only fast path is a
  later performance capability: 2 edits, 2 gate runs.

## Further notes

- Landing order inside the build: parser and migration land in one commit (the
  `docs-currency-workflow` check parses the real board's ledgers, so neither can
  land green alone); then commands, then status and the stripped journey, then the
  conformance check, then the docs. The docs ticket is blocked by the check ticket
  because both write `projects/benchkit.md`.
- `.bench/BENCH.md` (180/180) and `.agents/commands/bench-write-spec.md` (73/73)
  sit at their prose budgets; the docs ticket's edits to both are line-neutral, or
  a budget change is a reviewer grant landed in that ticket.
- The pre-migration capture is taken with the binary built at the integration
  source's base commit — the last one whose parser accepts the inline shape.
- The `bench` wrapper on this machine resolves an installed 0.2.0 release ahead of
  the dev build, so `bench coverage --check` there predates the reduced coverage
  schema; the check in this phase ran through `go run ./cmd/bench`.
