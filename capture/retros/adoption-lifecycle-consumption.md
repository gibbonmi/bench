## Outcome

The build landed `adoption-lifecycle-consumption` on `main` at 8ab514845e56c40f6219dc672c71f6368eef4768 from source tip 66c5f91c. `PlanLifecycle` now owns every managed lifecycle decision for Link, Setup, Upgrade, Unlink, and compatibility repair, and each caller observes facts and applies the selected effects. Three chunks closed: ALC-OWNER, ALC-LINK, and ALC-UNLINK, with four tickets and all 38 acceptance rows reconciled. The run id is `ft217-build-20261009`.

The reviewer decided three deltas during the build. An owned dangling-symlink CLAUDE.md now exits 0 instead of reporting a preserved conflict. ALC34 no longer asks for an Unlink final-manifest interruption, because the fault seam fires only on renames. The decision map now cites the tracked FT217 review, because the delivery closure removes `roadmap/FT217.md`.

## Gate-stage timings

- landing: commit 8ab514845e56c40f6219dc672c71f6368eef4768, trace 5452db566463d92753d16c7adb292852
- gofmt: 152 ms
- vet: 1561 ms
- test: 465875 ms
- race: 2735 ms
- system: 134709 ms
- shellcheck: 604 ms

The build ran from 07:19 to 18:29, 11 hours 10 minutes with a 22-minute pause. The whole-project gate ran 11 times: four `main` folds, five checkpoints, one timed-out checkpoint, and the landing. Each run took 8 to 13 minutes, about 2 hours in total. The `internal/adopt` package grew from about 180 seconds to 468 seconds inside the gate, and it hit the 600-second test timeout once.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ALC-OWNER ticket 01 | ALC02-ALC05, ALC31, ALC32 verified | verified | 9 | held | opus / high / implementer |
| ALC-LINK ticket 02 | ALC06-ALC15, ALC19, ALC25, ALC33, ALC38 verified | verified | 9 | held | opus / high / implementer |
| ALC-LINK ticket 03 | ALC16-ALC18, ALC21-ALC23 verified | verified | 9 | held | opus / high / implementer |
| ALC-UNLINK ticket 04 | ALC37 census claimed | claimed | 6 | refuted | opus / high / implementer |
| ALC-OWNER r1 Standards | S1 missing recorded reds | auto-fix | 6 | held | opus / medium / reviewer |
| ALC-OWNER r1 Standards | S2 copied table header | auto-fix | 5 | held | opus / medium / reviewer |
| ALC-OWNER r1 Standards | S3 reason literals repeat | no-op | 5 | refuted | opus / medium / reviewer |
| ALC-OWNER r1 Standards | S4 owner repeats caller rules | no-op | 6 | refuted | opus / medium / reviewer |
| ALC-OWNER r1 Standards | S5 comment argues correctness | auto-fix | 4 | held | opus / medium / reviewer |
| ALC-OWNER r1 Standards | S6 placeholder prose stale | auto-fix | 5 | held | opus / medium / reviewer |
| ALC-OWNER r1 Spec | P1 refusal text differs | auto-fix | 5 | held | opus / medium / reviewer |
| ALC-OWNER r1 Spec | P2 ALC08 seam cell stale | auto-fix | 7 | held | opus / medium / reviewer |
| ALC-OWNER r1 Spec | P3 probe prose stale | auto-fix | 6 | held | opus / medium / reviewer |
| ALC-OWNER r1 Coverage | C1 symlink mode guard silent | auto-fix | 8 | held | opus / medium / reviewer |
| ALC-OWNER r1 Coverage | C2 owned special CLAUDE.md silent | auto-fix | 8 | held | opus / medium / reviewer |
| ALC-OWNER r1 Coverage | C3 strict non-regular silent | auto-fix | 7 | held | opus / medium / reviewer |
| ALC-OWNER r1 Coverage | C4 invalid Preserve silent | auto-fix | 8 | held | opus / medium / reviewer |
| ALC-OWNER r2 Coverage | C5 dropped symlink silent | auto-fix | 9 | held | opus / medium / reviewer |
| ALC-LINK r1 Standards | H1 reason vocabulary two sources | auto-fix | 9 | held | opus / medium / reviewer |
| ALC-LINK r1 Standards | H2 ALC38 asset drifted | auto-fix | 8 | held | opus / medium / reviewer |
| ALC-LINK r1 Standards | H3 file-kind literals | auto-fix | 6 | held | opus / medium / reviewer |
| ALC-LINK r1 Standards | J1 duplicate Link runner | no-op | 6 | refuted | opus / medium / reviewer |
| ALC-LINK r1 Standards | J2 snapshot block pasted | auto-fix | 4 | held | opus / medium / reviewer |
| ALC-LINK r1 Standards | J4 caller routes partitions | no-op | 3 | refuted | opus / medium / reviewer |
| ALC-LINK r1 Spec | P1 ALC38 witness stale | auto-fix | 9 | held | opus / medium / reviewer |
| ALC-LINK r1 Spec | P2 dangling CLAUDE.md delta | ask-user | 8 | held | opus / medium / reviewer |
| ALC-LINK r1 Coverage | C1 ALC38 probe invalid | auto-fix | 10 | held | opus / medium / reviewer |
| ALC-LINK r1 Coverage | C2 absent seed symlink parent | auto-fix | 8 | held | opus / medium / reviewer |
| ALC-LINK r1 Coverage | C3 dropped symlink parent | auto-fix | 8 | held | opus / medium / reviewer |
| ALC-LINK r1 Coverage | C4 dangling CLAUDE.md untested | auto-fix | 9 | held | opus / medium / reviewer |
| ALC-LINK r1 Coverage | C5 spaced glob root | auto-fix | 7 | held | opus / medium / reviewer |
| ALC-LINK r2 Coverage | C5 dropped symlink row | auto-fix | 9 | held | opus / medium / reviewer |
| ALC-LINK r2 Standards | J2r snapshot remnant | auto-fix | 4 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Standards | N1 inline launch | auto-fix | 6 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Standards | N2 owner call duplicated | auto-fix | 5 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Standards | N3 observers repeat | no-op | 4 | refuted | opus / medium / reviewer |
| ALC-UNLINK r1 Standards | N4 dead refusal branch | no-op | 4 | refuted | opus / medium / reviewer |
| ALC-UNLINK r1 Standards | N5 reason round trip | no-op | 3 | refuted | opus / medium / reviewer |
| ALC-UNLINK r1 Standards | N6 env pair repeated | no-op | 4 | refuted | opus / medium / reviewer |
| ALC-UNLINK r1 Standards | N7 prose fixture repeated | no-op | 3 | refuted | opus / medium / reviewer |
| ALC-UNLINK r1 Standards | N8 foreign-hook text twice | ask-user | 7 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Spec | U1 bespoke choices bypass owner | auto-fix | 7 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Spec | U2 Unlink interruptions partial | auto-fix | 5 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Spec | U3 ALC36 evidence unrecorded | auto-fix | 5 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Spec | U4 headroom file absent | auto-fix | 4 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Coverage | F1 absent-row sweep silent | auto-fix | 9 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Coverage | F2 root-only fixture | auto-fix | 8 | held | opus / medium / reviewer |
| ALC-UNLINK r1 Coverage | F3 unreachable branch | no-op | 9 | refuted | opus / medium / reviewer |
| ALC-UNLINK r1 Coverage | F4 absent hook not covered | auto-fix | 7 | held | opus / medium / reviewer |
| lock-leak review Spec | old-binary window | ask-user | 8 | held | sonnet / high / reviewer |
| lock-leak review Spec | Close removes by name | ask-user | 4 | held | sonnet / high / reviewer |
| lock-leak review Coverage | Lstat error branch silent | ask-user | 6 | held | sonnet / high / reviewer |
| lock-leak review Coverage | untested edge states | no-op | 5 | refuted | sonnet / high / reviewer |
| dangling CLAUDE.md consultation | recommend accepting exit 0 | held by reviewer | 7 | held | sonnet / high / consultant |
| orchestrator | expected repair rounds not declared | abstained | abstained | abstained | fable then opus / unknown / orchestrator |

Brier mean 0.154 over 54 pairs, 1 abstained.

Each ticket author stayed inside its fence and committed green with its named probes biting. The ticket 01 author ran 63 probes, one per expectation, which cost about 28 minutes beyond a 13-minute implementation. After the reviewer capped witnesses, the ticket 02 to 04 authors ran 6 to 14 probes each. The ticket 02 author was stopped by a pause, and a replacement continued from its preserved test files. The ticket 04 author first left Unlink's AGENTS.md and hook choices outside the owner, which the Spec axis caught.

## Coordinator catches

- The ticket 01 author ran one probe per expectation. The orchestrator capped later witnesses at one per coverage row plus one per new branch.
- The orchestrator's own coordinator probe ran on every ticket and repair at a different site and kind, and each one bit.
- The ticket 03 move of the inline refusal block made the ALC38 probe asset stale; the review caught it, and the plan retargeted the witness.
- The recorder owes each ticket's verification to its own author at the final chunk source. Each repair to one ticket needed a verification session for its sibling.
- The orchestrator scoped the ALC-LINK repair to one ticket's fence and left a ticket 02 test copy, which cost one extra repair and review round.
- The orchestrator wrote integration records while the round 4 axes ran, which broke their evidence binding and cost one rebind.
- The orchestrator skipped the candidate doctor rehearsal before the landing; the landing refresh still completed.
- The repair session for the timeout found 2.18 million leaked adoption lock files; the reviewer directed a separate fix.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-expose-complete-decisions-through-upgrade-preview.md | 2 | delegate-error, delegate-error |
| 02-consume-decisions-for-desired-link-assets.md | 1 | one-source |
| 03-consume-decisions-for-dropped-and-bespoke-link-assets.md | 1 | one-source |
| 04-consume-unlink-decisions-and-close-entry-proof.md | 2 | delegate-error, check-gap |

## Agent-experience improvements

### Bench CLI

- Let `bench probe` take several mutations in one private build, because each call links the package test binary twice and a per-expectation loop cost 17 minutes.
  Feeds: FT168
- Let `bench test` report the slowest tests and warn when a package nears the gate's go test timeout. A test-time regression then shows before a checkpoint fails.
  Feeds: new
- Let `bench record verification` take its excerpt from the last Bench response of the worktree, because authors copied excerpts and spill files by hand.
  Feeds: new
- Let `bench preflight evidence --check-current` report the dirty paths instead of refusing, so a reviewer can bind beside an uncommitted record write.
  Feeds: new
- The landing census recorded 17 raw calls: cp 5, grep 3, sed 3, ls 2, rg 2, cd 1, wc 1. The excerpt form above would own the cp calls.
  Feeds: new

### Skills

- Name a package time budget in each author charge and require the floor's elapsed time in the return. FT217 doubled the adopt package before review saw it.
  Feeds: FT369
- Make the build preflight grade each spec decision map's Sources against the published tree before the first chunk. The delivery closure removed a cited roadmap file.
  Feeds: new
- Add the Skill tool to the Bench writer and reviewer agent definitions. Otherwise, state in the charge template that a delegate reads the phase command file.
  Feeds: new

### Process

- Give the review record an evidence-only confirmation form that a coordinator probe can satisfy, because the checkpoint required three fresh axes after every test-only commit.
  Feeds: new
- Run a chunk checkpoint gate only before a successor chunk; close the last chunk through the complete checkpoint, which grades every chunk again.
  Feeds: new
- Let one repair session record verification for every ticket of its chunk, because sibling re-verification sessions cost about 20 minutes in ALC-LINK.
  Feeds: new
- Remove each adoption lock file on release and sync the transaction record once per step, because leaked locks and per-entry syncs slowed every Link.
  Feeds: new
- Hold other landings while a build runs, because four `main` folds each cost one whole-gate run.
  Feeds: none