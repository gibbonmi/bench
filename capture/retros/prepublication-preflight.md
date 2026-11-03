## Outcome

`bench commit --preflight-build <slug> [--base <commit>]` now runs its chain before publication. The landing writes the candidate commit, runs the worktree build and `preflight.BuildVerdictAt` at that candidate, and moves the ref only on green. A red step refuses as `commit-chain-red` at exit 1, and exit 4 retired. By the reviewer decision of 2026-11-02, a root that declares no Bench build inputs skips the build step and still runs the preflight (story 35, PP35).

The build ran three tickets in three serial chunks as a `--full --delegate --auto-approve` run with `--reviewer opus medium` under minimal probing. Every chunk took one repair round and one confirming round. The final source passed the system suite (288 tests), and the coverage map held 35 valid rows. The landing `a8cce4f60` ran the whole gate green and released the source.

## Cost

- landing: a8cce4f6080ade98e431b0e235234fdb4f4e80ea, range a8cce4f6080ade98e431b0e235234fdb4f4e80ea^1..a8cce4f6080ade98e431b0e235234fdb4f4e80ea^2
- wall: 120 min
- commits: 28 total, 21 bookkeeping, 7 change
- graph: 3 tickets, chain 2, 3 chunks, 2 serial checkpoints
- dispatch: 3 authors, 3 repairs
- trace: commit 28 spans in 275279 ms, lane 0 spans in 0 ms, landing 1 spans in 322684 ms
- wall per ticket: 40.2 min against 60 min
- bookkeeping per ticket: 7.0 against 7.1
- delegates: about 1.69M tokens over 27 sessions, from the harness completion counts; provider charges are unknown
- census: 52 raw calls (rg 40, cat 4, cd 2, ls 2, sed 2, python3 1, wc 1)

## Gate-stage timings

- landing: commit a8cce4f6080ade98e431b0e235234fdb4f4e80ea, trace 60432e77f3687c9d3908ce325fb3399a
- gofmt: 212 ms
- vet: 1955 ms
- test: 234390 ms
- race: 3683 ms
- system: 51958 ms
- shellcheck: 561 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| unknown | unknown | unknown | unknown | unknown | unknown |

The orchestrator declared no expected repair rounds and no confidence, so this build has no labeled pair. Brier mean: unknown, over 0 pairs, with 0 abstentions.

Ticket 01 (Opus/high) met every C1 row in one pass and probed each row red. Its review found four repair targets, and three of them came from spec gaps: two source-tip reads with no varying fixture, and an undefined dirty rule at a non-HEAD commit. Ticket 02 (Opus/medium) met every C2 row in one pass. Its review found a face-routing order that C3 would hit and two duplicated test harnesses. Ticket 03 (Opus/medium) stopped before its first edit on a real fence gap: the chain could read the slug and the base only after publication. After the widen, it met every row, but its review found an undecided state in linked projects.

## Coordinator catches

- The C3 author's fence stop was correct. The orchestrator verified the claim at `cmd/bench/commit_chain.go:44-61` and widened ticket 03 to `internal/commit/grammar.go` with no reviewer stop.
- The C3 Spec axis confirmed that a linked project without `scripts/go-build.sh` would refuse every chained commit. The orchestrator stopped for the reviewer, who chose to skip the build there.
- The orchestrator's own spec sentence for the C1 dirty-rule divergence overclaimed twice. The C1 and C2 Spec axes each caught one overclaim.
- The `docs/known-issues.md` amendment workaround refused before the first record existed. The C1 author reported it, and the orchestrator captured a learning.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-grade-the-candidate-commit.md | 1 | spec-row |
| 02-check-before-publication.md | 1 | edge-found |
| 03-chain-before-publication.md | 1 | edge-found |

## Agent-experience improvements

### Bench CLI

- Make `bench probe --help` state the `--mutants` file format, a JSON array of argv arrays whose entries omit the subject file and take no `--expect`.
  Feeds: new
- Make `bench record amendment` answer "no record yet" with no `next=` route when `reviews/<slug>.md` holds no record.
  Feeds: new
- Make `bench record excerpt` accept a final requirement such as `system`, so the orchestrator does not write a final excerpt by hand.
  Feeds: new
- Count a census `rg` call apart from the other heads, because 40 of the 52 raw calls were sanctioned searches that hide the 4 `cat` and 2 `sed` reads of named files that `bench worktree read` serves.
  Feeds: new

### Skills

- Make `craft-spec` require one varying fixture for each read that a differential row enumerates, or a won't-handle line for that read.
  Feeds: new
- Make the ticket template state the predicate of a "tree form" for both arguments when it asks a ticket to keep one source behind a new argument.
  Feeds: new
- Make the author charge template tell the author to search the package tests for an existing helper before it adds a test harness.
  Feeds: none

### Process

- Read the provider scorecard before the line declaration. This build skipped the read, so it missed unified review for a three-ticket spec and ran three axes in every round.
  Feeds: none
- Make the write-spec edge inventory name each repository class that a chained kit verb reaches, such as a linked project without build inputs.
  Feeds: new
- Correct the `docs/known-issues.md` amendment row so it names the state in which it applies.
  Feeds: none