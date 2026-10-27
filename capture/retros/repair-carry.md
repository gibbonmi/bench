## Outcome

The FT448 repair carry landed at `f6f3f3fcc` from the reviewed source `0dbfa0147..2d0842cfb`. `bench record ticket` now carries an owed requirement that a test-only repair delta cannot reach, and it writes a carried entry through the one verification writer. A refusal names each check that a repair session must run, with its reason. The guidance, ADR 0023, and the changelog state the rule. The landing gate was green, and the final coverage check passed for 51 rows.

## Gate-stage timings

- landing: commit f6f3f3fcc537ff287ea8f83a48e024c1d1d62f4a, trace 5ce7e609889d014ffbb0f9c3ec19331c
- gofmt: 178 ms
- vet: 1598 ms
- test: 367718 ms
- race: 2415 ms
- system: 55726 ms
- shellcheck: 520 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 03 repair 1 | the canonical rule refuses a mid-path wildcard | verified | 9 | held | opus / high / repair |
| ticket 04 repair 1 | the two-reason refusal case bites | verified | abstained | abstained | opus / high / repair |
| ticket 05 | CF49 red then green | verified | 9 | held | opus / high / implementer |

Brier mean 0.01 over 2 pairs; 1 abstained.

Each author stayed inside its ticket's chunk and its fence, and each widening went through `bench spec widen`. Ticket 03 delivered 28 rows in one session at about 165k tokens, which is above the 150k target.

## Coordinator catches

- The CF-C2 Spec and Coverage axes found that the spec's canonical rule accepted `./internal/a...` and `./../x`. The orchestrator tightened spec line 119 and added row CF51 under `--auto-approve`.
- The CF-C3 Spec axis found a full-rerun sentence in ADR 0023 that the spec's reader sweep missed. The orchestrator widened ticket 05 with a decision reference.
- The orchestrator named a plan commit as the first chunk base in the ticket 01 charge. The ticket 02 freeze corrected the base to the `main` tip.
- The orchestrator started the prepared CF-C2 run before the review returned, and the run overlapped the repair session.
- `bench record confirmation` refused the ticket 04 repair, because a fixture helper file is not a Go test file. The repair took a confirming round instead.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-share-the-plan-command-parser.md | 0 | none |
| 02-share-the-pass-predicate.md | 0 | none |
| 03-decide-the-carry.md | 1 | spec-row |
| 04-record-carried-entries.md | 1 | one-source |
| 05-state-the-repair-carry.md | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Make `bench record ticket --help` state that the first ticket of each chunk passes `--base` with the chunk base.
  Feeds: new
- Let each entry of a `bench probe --mutants` set name its own subject file, so that one call covers a witness across two files; this replaces the one `python3` call that wrote two mutant files.
  Feeds: new
- Make the page index of `bench preflight evidence --to` name the repository path of each source, so that an author reads no exported file by hand; this replaces the four `cat` reads of exported evidence.
  Feeds: new
- Route a line-count read to `bench structure` in the author charge template, because the one `wc` call measured a line budget that `bench structure` already reports.
  Feeds: none
- Count the 16 `rg` calls apart in the census as sanctioned reads, so that the one `cd` call of unknown purpose stays visible.
  Feeds: new

### Skills

- Make the build entry of `/bench-implement-spec` name the chunk base rule in each first-ticket charge, and commit each earlier author record before the next charge.
  Feeds: new
- Make the reader sweep of `craft-spec` cover `docs/adr` for each statement of a retired rule.
  Feeds: new

### Process

- Start a prepared chunk run only after the review returns clean, or wait on its PID before a repair dispatch.
  Feeds: new
- Keep the minimal-probing direction for a build whose rows sit in table-driven tests, because one mutant for each decision step bit every row group and no review found a silent row.
  Feeds: none