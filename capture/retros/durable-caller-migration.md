## Outcome

The build landed at 7619c5ff on a green landing gate. Every migrated caller of a durable file replacement now composes `internal/durablefile`. Each caller classifies a failure by its `Published` field, and the probe, intent, and approval callers act on that class. The build ran 13 one-ticket chunks as a `--full --delegate` run with an author limit of 3. Opus/medium wrote every ticket, and Opus/high ran every review axis and consultation.

## Gate-stage timings

- landing: commit 7619c5ff331dec23c0c8aeab1532d121ddb632e2, trace 82d7f042fb1f8d2b7d3ddfba736c956e
- gofmt: 151 ms
- vet: 1527 ms
- test: 174564 ms
- race: 2400 ms
- system: 61185 ms
- shellcheck: 573 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| line declaration | expected repair cycles per chunk | abstained | unknown | unknown | opus / high / orchestrator |

The orchestrator declared no labeled claim in this build. The Brier mean is unknown over 0 pairs, with 1 abstention.

Every author stayed inside its planned ticket and chunk boundary. Eleven chunks used one repair cycle or none, and D-B2, D-B4, and D-D used both cycles. Most first reviews found a one-source defect or a comment that restated a fact, not a behavior defect. Tickets 10, 12, and 13 authored in parallel worktrees while earlier chunks were in review, and each merge gate was green.

## Coordinator catches

- Two authors recorded probe entries with exit codes that the checkpoint refuses, and the first correction was also wrong. The orchestrator read the checker and corrected the charge.
- The chain check refused a chunk base after a comment-only correction. The orchestrator refroze the successor on the previous tip and confirmed the wider pair with the three axes.
- Two checkpoint reds came from flakes that predate this build. A consultation proved each one unowned before a retry.
- The orchestrator assigned a worktree with the integration assignment ID by mistake, and restored the spec before the commit.
- The D-D census found refusal and pre-rename gaps that the closed chunks left, and a plan expansion closed them in one chunk.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 02-preserve-gate-failures.md | 1 | other |
| 03-migrate-owner-records.md | 2 | one-source, one-source |
| 04-migrate-capture-documents.md | 1 | other |
| 05-migrate-handoff-documents.md | 2 | one-source, other |
| 06-preserve-ledger-effects.md | 1 | other |
| 07-recover-approval-stage.md | 1 | one-source |
| 08-migrate-publication-records.md | 1 | spec-row |
| 09-migrate-broker-manifests.md | 1 | one-source |
| 10-migrate-dashboard-artifacts.md | 1 | one-source |
| 11-preserve-probe-recovery.md | 1 | spec-row |
| 12-migrate-assessment-records.md | 0 | none |
| 13-migrate-repair-documents.md | 1 | one-source |
| 14-close-replacement-ownership.md | 2 | one-source, one-source |

## Agent-experience improvements

### Bench CLI

- Make `bench record verification` refuse a probe entry that the checkpoint would refuse, such as an outcome of `bit` with a probe exit of 0.
  Feeds: new
- Make the chunk chain check accept a base that differs from the previous chunk digest only by a comment gap.
  Feeds: new
- Add a read-only verb that prints one chunk's plan rows, which the 154-call landing census entry asks for.
  Feeds: new
- Make `bench probe` print the swap text it applied, so that a probe excerpt is complete without a hand copy.
  Feeds: new

### Skills

- State in the implement-spec skill that a derived expectation needs no red and a repair replays only the reds it touches.
  Feeds: new
- Tell the first caller ticket of a multi-caller migration to name the shared refusal assertion and the pre-rename row in the shared test package.
  Feeds: new

### Process

- Reorder the delegated plan to its merge order before the first concurrent dispatch, because a reorder after a chunk tip blocks that chunk checkpoint.
  Feeds: none
- Diagnose the ownership of every checkpoint red before a retry, because two of the reds in this build predate it.
  Feeds: none