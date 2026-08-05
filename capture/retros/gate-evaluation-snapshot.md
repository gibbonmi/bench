# Retro: gate-evaluation-snapshot

## Outcome

Promoted terminal on 2026-08-04: candidate `60b44de` published as squash `c3d8e47` on
`main`, spec flipped to `Status: implemented` by promotion, exact green retained at
`refs/bench/green/main`. Five assignments integrated and released over one run: three
original tickets (expand generation source, migrate evaluation identities, contract
independent captures) built by a prior session, plus two repair tickets from the accepted
review findings — S1 (optimistic exact-green reuse built its subject through root-based
`buildSubject` instead of the evaluation-owned generation; fixed by routing
`ExecuteReusingFreshGreen` through `acceptPre` and making `reusableEvidence` project from
the caller's plan) and C1 (the listing parser accepted any tab-separated metadata; fixed
by refusing entries outside git's mode-type-object shapes at parse). The fresh composed
review of the repaired candidate returned zero findings on all three axes.

## Gate-stage timings

Only session-observed values; promotion's internal gate was not streamed and its stage
split is unmeasured.

- Full `internal/gate` package suite: ~146–153s per run (three full runs this session).
- Focused snapshot/reuse tests: 0.01–0.22s per run — the repair loop iterated here.
- Ticket-staging `bench commit` gate (full phase run incl. conformance): ~2–3 min;
  `package-core-guard` 7.4s and `line-routing` 1.2s were the largest visible
  conformance entries.
- Mutation probes: four probes (two per repair), each a focused-test round trip under 1s
  plus one ~150s full-suite confirmation per ticket.

## Ticket-versus-spec-slice and delegate performance

The repair round ran inline by explicit reviewer instruction (no write delegates), so no
delegate comparison exists for it; the spec's per-story lines (mid-tier / medium) were
collapsed into one session-model inline line, flagged at the time. Ticket-sized charges
still shaped the work: each repair was one ownership-fenced ticket with two acceptance
rows and landed through assign → checkpoint → integrate in one pass, no re-charge. The
serialized single-fence cadence (`internal/gate/` owns everything) cost nothing here
because both repairs were small; a wider repair round would have been throttled by the
one-writer fence.

## Coordinator catches

- The S1 defect was deeper than its finding title: `reusableEvidence` called `inspectAt`,
  which performed a second independent `buildSubject` capture — the reuse hit did two
  root-based materializations, not one. Fixing only the named call site would have left
  the bypass half-closed; the count assertion (exactly one parsed listing per reuse hit)
  is what forced the second site into view.
- The count-recorder fixture mutates PATH, and PATH is part of the oracle identity: a
  recorder installed after the seed run silently kills the reuse it is trying to count.
  The test had to install the recorder before seeding and count the operation delta.
- Two full-suite failures in a fresh assignment worktree were environment, not diff:
  `TestFT78Story4ProofLedger` needs `dist/bench`, which no fresh worktree has until
  `scripts/go-build.sh` runs. Verified by rebuilding and rerunning before attributing.

## Agent-experience improvements

### Bench CLI

- Active assignment worktrees are hard to discover from a fresh session (preserved from
  the prior session's pin): the path only surfaces in the one-time `assign` output or by
  reading durable state JSON under `.git/bench/specbuild/`. `bench spec build status
  --full` lists them but buries the path mid-row; a `bench worktree path <target>`-style
  direct answer for assignments would close it.
- Checkpoint receipt assembly is manual JSON against an undocumented schema: the receipt
  and probe field contracts live only in `internal/specbuild/checkpoint.go`, and the
  coordinator recomputed `TreeHash` by replicating its throwaway-index sequence by hand.
  A `bench spec build receipt --assignment <id>` scaffold (or a documented schema) would
  remove the riskiest hand-derived step.
- A fresh assignment worktree fails two gate-adjacent tests until `dist/bench` exists;
  assign could pre-build it or the fixture could report the missing prerequisite as a
  capability skip instead of a failure.

### Skills

- No friction: craft-tickets' red-mutation table mapped one-to-one onto the checkpoint
  receipt's probe and rows fields, and craft-review's refute-first rule caught nothing to
  cut because the axes were walked against an already-small composition.

### Process

- The repair cadence (findings → tickets → assign/checkpoint/integrate → fresh composed
  review → promote) fit the two-finding round exactly; the only surprise was that
  committing the repair ticket files moved the branch tip and forced a `promote`-driven
  recomposition before `assign` would proceed. That ordering (land tickets, recompose,
  then assign) is worth stating in the implement-spec phase text; discovering it through
  two lifecycle refusals cost one round trip each.
- The on-disk session handoff described a different spec than the reviewer's pinned
  continuation; the tree-wins rule resolved it instantly, which is the rule working as
  designed — but it confirms the handoff must be rewritten at every phase close, not
  only at session end.
