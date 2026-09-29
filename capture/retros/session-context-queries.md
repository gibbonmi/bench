## Outcome

The delegated build of `session-context-queries` landed at `7072d08b`. It adds a selected view to `bench worktree list` (`--view paths --target`), a selected view to `bench spec history` (`--spec --limit`), and focused-read guidance in `craft-cli`, `bench-debug`, and `bench-drain`. The build re-authored the 2026-09-22 wave stream `194b7dba` on `main` after a staleness audit and a spec amendment. All 25 acceptance rows are covered.

## Gate-stage timings

- landing: commit 7072d08ba1809e4018d14df7c19cf2c88028b2e8, trace fe8a22f90378bc4ea1fb5ac3ff8c324e
- gofmt: 125 ms
- vet: 1256 ms
- test: 138798 ms
- race: 3061 ms
- system: 41828 ms
- shellcheck: 520 ms

## Ticket-versus-spec-slice and delegate performance

Each of the three fresh Opus/high ticket authors landed its ticket in the first attempt and inside its fence. Each author ported the stream code against the current tree. The Spec axis found zero findings on each first round, so the ticket boundaries held. Five Opus/high repair sessions and one re-verification session each finished in one attempt. The 27 review axes ran on Opus/high over 9 rounds.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| QU-C1 Standards | the selected view copies two refusal texts | accepted | 6 | 1 | Opus / high / reviewer |
| QU-C1 Standards | four test expectations restate their owners | accepted | 5 | 1 | Opus / high / reviewer |
| QU-C1 Standards | a test comment restates the QU26 red record | accepted | 4 | 1 | Opus / high / reviewer |
| QU-C1 Coverage | a `--target`-only request is not graded | accepted | 9 | 1 | Opus / high / reviewer |
| QU-C1 Coverage | the QU18 probe row grades help discovery | accepted | 9 | 1 | Opus / high / reviewer |
| QU-C1 Coverage | four author probes cannot be rerun | accepted | 7 | 1 | Opus / high / reviewer |
| QU-C1 Spec | a grammar test comment overstates one case | confirmed | 8 | 1 | Opus / high / reviewer |
| QU-C2 Standards | the unsafe-operand ordinal has two sources | accepted | 5 | 1 | Opus / high / reviewer |
| QU-C2 Standards | a test repository setup has two sources | accepted | 3 | 1 | Opus / high / reviewer |
| QU-C2 Coverage | a failed history row has no graded detail | accepted | 6 | 1 | Opus / high / reviewer |
| QU-C2 Standards | `initGitRepo` duplicates `gittest.RepoOnBranch` | accepted | 7 | 1 | Opus / high / reviewer |
| QU-C3 Standards | a membership sentence points at the wrong table | accepted | 7 | 1 | Opus / high / reviewer |
| QU-C3 Standards | the drain `detail` pointer is unclear | accepted | 6 | 1 | Opus / high / reviewer |
| QU-C3 Standards | a drain sentence restates the inventory fact | accepted | 5 | 1 | Opus / high / reviewer |
| QU-C3 Standards | one sentence gives two instructions | accepted | 6 | 1 | Opus / high / reviewer |
| QU-C3 Spec | the record overstates the kept debug text | accepted | 6 | 1 | Opus / high / reviewer |
| QU-C3 Coverage | the new table turns a conformance test red | accepted | 10 | 1 | Opus / high / reviewer |
| QU-C3 Standards | a condition follows its instruction | accepted | 5 | 1 | Opus / high / reviewer |
| QU-C3 Standards | a record sentence is passive | accepted | 6 | 1 | Opus / high / reviewer |
| fold Standards | the reconciliation text contradicts the fold | accepted | 9 | 1 | Opus / high / reviewer |

The Brier mean is 0.158 over 20 labeled pairs, with 4 abstentions. Each labeled claim held, so the reviewers were right and under-confident.

## Coordinator catches

- The ticket 1 author recorded the plan commit as the chunk base, and the coordinator corrected the base to the `main` tip.
- The QU-C2 first checkpoint went red on an intermittent land test with an infrastructure refusal, which passed alone and on a rerun.
- The ticket 3 commit passed its Markdown lane with `TestAXIGuidanceContractBites` red, and the Coverage axis caught it.
- The coordinator's own light-path landing moved `main` during the build. The first landing refused on the composition, and a Fable/high consultation found the fold route.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-select-worktrees.md | 1 | one-source |
| 2-select-histories.md | 2 | one-source, one-source |
| 3-guide-relevant-reads.md | 2 | check-gap, other |

## Agent-experience improvements

### Bench CLI

- Print the plan digest and the source digest from `bench preflight review`, so that no author computes them by hand.
  Feeds: new
- Stop the census from counting `rg` where `AGENTS.md` prescribes it. The `scq-build` census was 59: `rg` 49, `ls` 6, `cat` 3, and `cp` 1.
  Feeds: new
- Make the lane classifier run each conformance check that reads an edited guidance file.
  Feeds: new

### Skills

- Make each ticket that repeats a per-target pattern from a sibling chunk name the owner of that pattern.
  Feeds: none

### Process

- Adopt the fold route as the one sanctioned route for a `main` that moves after the last chunk.
  Feeds: FT342
- Do not land an unrelated change on `main` while a chunked build is open.
  Feeds: FT342
- Charge each ticket author with the chunk base and with one record entry for each author probe.
  Feeds: none