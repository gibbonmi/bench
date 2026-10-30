## Outcome

The `markdown-block-reader` build (FT358) landed at `561fdc88` on a green whole-project gate. The new `internal/markdown` package owns the four block rules: frontmatter, fenced blocks, HTML comments, and H2 headings. The prose, spec, coverage, roadmap, maps, tickets, handoff, learnings, retros, skills-index, review-record, anchors, and adopt readers now read their lines from it.

The new `markdown-block-owner` conformance check refuses a block-rule string literal in production Go outside the reader. The spec closes 106 coverage rows. The build added stories 78 to 85 and rows MB102 to MB109 under the reviewer's advance approval.

The reviewer ordered this build before FT125, because the FT125 chunk ER-C2 consumes the block reader. FT125 stays blocked until its next start.

## Cost

- landing: 561fdc88613e693e06b302b906ef11cb4873c89e, range 561fdc88613e693e06b302b906ef11cb4873c89e^1..561fdc88613e693e06b302b906ef11cb4873c89e^2
- wall: 438 min
- commits: 102 total, 79 bookkeeping, 23 change
- graph: 11 tickets, chain 3, 5 chunks, 3 serial checkpoints
- dispatch: 11 authors, 11 repairs
- trace: commit 112 spans in 1050734 ms, lane 0 spans in 0 ms, landing 1 spans in 322763 ms
- wall per ticket: 39.9 min against 60 min
- bookkeeping per ticket: 7.2 against 7.1, over budget

Provider usage and charges are unknown. The run used 79 sessions: one orchestrator, six staleness readers, one amendment writer, 11 authors, 11 repairs, and 48 review axes.

## Gate-stage timings

- landing: commit 561fdc88613e693e06b302b906ef11cb4873c89e, trace 69706da3bfb3f7b5d390c0bbb7367118
- gofmt: 175 ms
- vet: 1605 ms
- test: 217522 ms
- race: 9065 ms
- system: 43795 ms
- shellcheck: 501 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| line declaration | The build needs about 3 repair rounds in total. | refuted: the five chunks used 11 repair sessions | 6 | 0 | claude-opus-5-5 / medium / orchestrator |

Brier mean 0.36 over 1 pair, 0 abstained.

Every author met its ticket boundary, and no author stopped on a fence gap. Each chunk review still found edge gaps between readers, so four of five chunks used both repair cycles or more. MB-C3 and MB-C4 each took one overrun under the advance approval.

## Coordinator catches

- The staleness pass found 55 drifted paths. A retired `FenceTokens`, a retired step walk, removed closures, and moved sites made three rows obsolete and one row a restatement.
- Two plan probes named live-tree tests that a `--check` run never selects, and two other probes could only remove diagnostics. The orchestrator amended each probe before verification.
- A probe excerpt that held an HTML comment opener hid the rest of the review record. Ticket 7 then made the record writer encode the less-than sign (MB107).
- The first author charge omitted the excerpt sequence, so the ticket 1 probes ran twice.
- The spec named the anchor section walk a section-body reader, but every caller passed a whole document. The MB-C4 review caught it, and MB94 got a new statement.
- A spec-prose sync after the last freeze made the MB-C5 review stale, so the chunk took a refreeze and three short reaffirmations before the completion entry.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-add-block-reader.md | 2 | delegate-error, delegate-error |
| 10-read-agents-marker-blocks.md | 0 | none |
| 11-forbid-block-rule-copies.md | 1 | one-source |
| 2-read-spec-and-coverage-blocks.md | 1 | delegate-error |
| 3-read-roadmap-blocks.md | 2 | edge-found, one-source |
| 4-read-field-scan-blocks.md | 1 | edge-found |
| 5-read-handoff-blocks.md | 3 | edge-found, one-source, delegate-error |
| 6-read-journal-blocks.md | 1 | delegate-error |
| 7-read-review-record-blocks.md | 1 | edge-found |
| 8-read-anchor-blocks.md | 3 | spec-row, edge-found, delegate-error |
| 9-read-skills-frontmatter.md | 0 | none |

## Agent-experience improvements

### Bench CLI

- Count `rg` calls apart from the other raw heads in the landing census, because 268 of 300 raw calls were sanctioned reads that hid the 32 others.
  Feeds: new
- Give the review-record writer a test that every excerpt survives the block reader, so that no future field can reopen the MB107 defect.
  Feeds: none
- Add a base-versus-tip parity form to `bench test`, because authors ran 16 raw `cp`, `diff`, `cmp`, and `sha256sum` calls to compare a parser's output across two builds.
  Feeds: new
- Make `bench record ticket` name each owed requirement of a sibling ticket that a repair delta reaches, because three repairs recorded only their own ticket's requirements.
  Feeds: new

### Skills

- Add to `craft-tickets` that a plan probe on a `--check` requirement must add a diagnostic on the live tree; a weakening mutation runs as a package-form unit test.
  Feeds: new
- Add to `craft-spec` that a rule-delta table names each site's heading trim predicate and each caller premise, because MB-C3 and MB-C4 each broke on one.
  Feeds: new

### Process

- Seed the common author charge with the review lessons of earlier chunks before the first dispatch, because each later chunk repeated an edge class that an earlier review found.
  Feeds: new
- Batch every spec-prose sync before the last chunk freeze, because a sync after it forces a refreeze and a reaffirmation round.
  Feeds: none