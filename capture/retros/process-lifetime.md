## Outcome

FT362 landed as `c2022288ea28609addc7f6c48de7a0ca091d85a7` from source tip `0fba401df` on review base `4115afad7`. The build ran 14 tickets in 5 chunks as a `--full --auto-approve` run with an author limit of 1. Authors ran on Opus at high effort, repairs on Opus at medium effort, and review axes on Opus at medium effort.

The build puts every production child process under one subprocess Lifetime owner with bounded shutdown and durable process-user protection in tmpfs stores. A start gate holds each registered child until its users are recorded. Builds, kit tests, caches, prospective bundles, release evidence and worktrees retain their resources while a registered user remains. A registered conformance check holds the census of eleven callers.

The reviewer deferred the native Darwin witness for PL90 to a follow-up, because the build host is Linux; a parked idea carries it. The final verification passed 5 of 5 requirements, and the completion reconciled 130 coverage rows, with PL90 deferred and PL103 owned by review.

## Gate-stage timings

- landing: commit c2022288ea28609addc7f6c48de7a0ca091d85a7, trace 1cdcfc2a4506fbd8c018ce6913bf3951
- gofmt: 208 ms
- vet: 1680 ms
- test: 196216 ms
- race: 2661 ms
- system: 48068 ms
- shellcheck: 626 ms

The FT402 landing gate took a test phase of 169 s and a system phase of 35 s. FT362 adds about 27 s of unit tests and 13 s of system tests, for 265 more system tests than FT402 left. The first PL-C3 checkpoint gate went red with a 395 s system phase, because every launch serialized on durable writes under each inherited store's lock. After the tmpfs store amendment, the PL-C3 checkpoint system phase took 40 s. The two PL-C4 checkpoint gates took about 275 s and 270 s.

## Ticket-versus-spec-slice and delegate performance

Each ticket fit its slice, and no author crossed a chunk boundary. Three tickets needed successor sessions: ticket 11's author neared its context limit, and ticket 13 split into an audit session and two matrix sessions. Ticket 13's first author stopped before its first edit, because the closure audit had to be a registered conformance check whose paths the fence did not hold. PL-C3 took four review rounds and one amendment, PL-C4 took five rounds, and PL-C5 took two.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| PL-C5 r1 Standards | S1 the reclaim scan re-derives the store-name prefix | finding | 8 | true | opus / medium / reviewer |
| PL-C5 r1 Spec | P1 the matrix compares author expectations, not a baseline | finding | 6 | true | opus / medium / reviewer |
| PL-C5 r1 Spec | P2 the protection diagnostic prints a nil cause | finding | 6 | true | opus / medium / reviewer |
| PL-C5 r1 Spec | P3 the shift deadline path prints no shift_result | finding | 5 | true | opus / medium / reviewer |
| PL-C5 r1 Spec | P4 start-failure and bounds timeout shapes are missing | finding | 5 | true | opus / medium / reviewer |
| PL-C5 r1 Spec | P5 PL118 skips the test directory and the worktree store | finding | 5 | true | opus / medium / reviewer |
| PL-C5 r1 Coverage | C5-COV-1 the fail-closed reclaim branches have no test | finding | 5 | true | opus / medium / reviewer |
| PL-C4 r1 Standards | S1 the sealed-reopen policy has three copies | finding | 5 | true | opus / medium / reviewer |
| PL-C4 r1 Standards | S2 the retired-store predicate has two copies | finding | 5 | true | opus / medium / reviewer |
| PL-C4 r1 Standards | S3 the environ vector read has a third copy | finding | 4 | true | opus / medium / reviewer |
| PL-C4 r1 Standards | S4 the inherited registration has four copies | finding | 4 | true | opus / medium / reviewer |
| PL-C4 r1 Standards | S5 exec has two exit-code mappings | finding | 3 | true | opus / medium / reviewer |
| PL-C4 r1 Standards | S6 the release child policies repeat | finding | 2 | false | opus / medium / reviewer |
| PL-C4 r1 Standards | S7 the store-discard test helper repeats | finding | 2 | false | opus / medium / reviewer |
| PL-C4 r1 Standards | S8 the observer row repeats a validator | finding | 1 | false | opus / medium / reviewer |
| PL-C4 r1 Spec | P1 reclaim cannot derive the store of an absent repository | finding | 5 | true | opus / medium / reviewer |
| PL-C4 r1 Spec | P2 held runs and the probe set drop Close failures | finding | 6 | true | opus / medium / reviewer |
| PL-C4 r1 Spec | P3 a kit-run Close failure replaces the child result | finding | 4 | true | opus / medium / reviewer |
| PL-C4 r1 Spec | P4 reset does not consult the protected-user predicate | finding | 2 | false | opus / medium / reviewer |
| PL-C4 r1 Spec | P5 exec exit 1 after the wait delay | finding | 1 | false | opus / medium / reviewer |
| PL-C4 r1 Spec | P6 the Windows fallback removal | finding | 1 | false | opus / medium / reviewer |
| PL-C4 r2 Coverage | the probe set can drop a held Close failure with no red | finding | 5 | true | opus / medium / reviewer |
| PL-C4 r4 Coverage | a failed stream with child zero has no biting test | finding | 5 | true | opus / medium / reviewer |

Brier mean 0.195 over 23 pairs; 0 abstentions. The table holds the PL-C4 and PL-C5 review claims, because the context compaction removed the PL-C1 to PL-C3 claim tables.

## Coordinator catches

- The coordinator traced the red PL-C3 checkpoint to durable writes under every inherited store's lock, and the reviewer approved the tmpfs store and start-gate amendment that cut the system phase from 395 s to 40 s.
- The coordinator found that ticket 11's author misread "no successor caller" and left the protected-user predicate to ticket 12, and gave the wiring to a successor session.
- The coordinator caught that ticket 11's plan had no system check, after the PL-C4 checkpoint went red on BO69, and added the check before the repair.
- The coordinator found that the shift deadline defect predates FT362, because the response bound of FT336 holds output that `os.Exit` skips, and recorded it as a learning instead of a build repair.
- The coordinator recorded the nine matrix reds in the build log after the Spec axis found that they lived only in a session report.
- The coordinator's `bench commitment block` and `unblock` cycle dropped the outcome's delivery binding, and a new `bench commitment start` with the same token restored it.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-own-lifetime-through-bounds.md | 1 | delegate-error |
| 02-protect-nested-bounds-users.md | 1 | delegate-error |
| 03-migrate-gate-process-policy.md | 1 | delegate-error |
| 04-retain-selected-build-binaries.md | 6 | check-gap, one-source, delegate-error, spec-row, edge-found, one-source |
| 05-protect-private-test-directories.md | 2 | check-gap, delegate-error |
| 06-refuse-clean-of-protected-cache.md | 0 | none |
| 07-retain-prospective-gate-artifacts.md | 2 | edge-found, spec-row |
| 08-preserve-release-step-evidence.md | 0 | none |
| 09-bound-preflight-and-scanner-shutdown.md | 0 | none |
| 10-preserve-focused-test-completion.md | 2 | edge-found, delegate-error |
| 11-preserve-worktree-exec-and-shell.md | 3 | one-source, check-gap, delegate-error |
| 12-retain-shift-scratch-through-recovery.md | 1 | one-source |
| 13-enforce-final-lifetime-closure.md | 1 | spec-row |
| 14-gate-registered-starts.md | 2 | delegate-error, one-source |

## Agent-experience improvements

### Bench CLI

- Make `bench spec widen` add the `(new)` marker to a path that does not exist yet, because three widenings needed a raw `sed` edit and a red `writes-resolve` commit (owned gap).
  Feeds: new
- Make `bench spec widen` append the floor checks that a widened package needs, because six raw `python3` plan edits followed six red `check-floor` commits (owned gap).
  Feeds: new
- Add a read for a spilled Bench response by its spill path, because 15 raw `cat` calls read spilled response files whose format Bench owns (owned gap).
  Feeds: new
- Count raw `rg` calls on Bench-owned files apart from source searches in the census, because 68 sanctioned `rg` calls hide the owned-gap reads of the plan and the review record.
  Feeds: new
- Keep the delivery binding across `bench commitment block` and `bench commitment unblock`, because the unblock left a binding gap that refused the next plan commit.
  Feeds: new
- Make `bench record review` accept a review of the last chunk when a deferred row leaves the chunk unfrozen, or name the deferral form in its refusal, because PL-C5's reviews waited in the build log until the freeze.
  Feeds: new

### Skills

- Make the implement-spec author charge require the system suite in a ticket's checks when the ticket changes a CLI verb's exit or output, because ticket 11 broke BO69 and only the checkpoint gate caught it.
  Feeds: new
- Make craft-tickets list the census pin file and the conformance registry paths in `Writes:` for a ticket that adds tests to a pinned package or a live-tree audit, because tickets 11 and 13 stopped on both.
  Feeds: new
- Make craft-spec flag a native-platform row that the build host cannot run, at spec time, because PL90 stopped the build at the last chunk.
  Feeds: new

### Process

- Build the approved repair-cycle redesign: no restage, one run per repair at its commit, the test-only confirmation form, one plan commit, the gate started at the repair commit, run-binary reuse, and the test-only delta checkpoint.
  Feeds: new
- Slice chunks at four tickets or fewer, because PL-C3's six tickets made each restage cost 25 to 30 minutes.
  Feeds: new
- Fix `bench shift` so that a checkpoint exit unwinds instead of calling `os.Exit` under the response bound, because every deadline and interrupt exit drops its `shift_result`.
  Feeds: new