# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `36049d8`, clean tree, 5 unpushed commits
Spec: `specs/ft148-worktree-orphan-retirement.md` — `Status: implemented`
Gate: green at HEAD — current

## State

- **FT148 is built, reviewed, gate-green, and unpushed. The reviewer owns the
  push.** Three commits: `a1aaafb` (all twelve stories), `f4103e4` (eight review
  findings fixed), `36049d8` (the status flip). Every phase of the gate is green
  including canary; ship-tier verification has not run and is not expected to —
  that is `bench prep-release`, once per release.
- **The feature works on real data, verified end-to-end rather than only by
  tests.** The session-start wall went from twenty-odd lines to five. One genuine
  residue row was compacted, the two worktrees this build cut were reported as
  orphans and then retired through the emitted command, and a freshly stamped
  worktree was correctly left alone. Seventeen `recovered` rows remain and are
  FT98's, not a defect.
- **Orphanhood is age alone, and that is still the load-bearing fact.** No
  liveness signal exists: `bench worktree create` writes no lease, and a lease
  records a pid that dies when the create hook exits. Safety rests on
  `bounds.AssignmentStale` being 7 days, the sweep only ever reporting, and the
  explicit cleanup recovering into a recovery ref before removing. Shortening the
  window or making the sweep reap is a spec deviation, not a refinement.
- **Two spec defects the build found and did not silently paper over.** Story 7's
  "an orphaned record that does hold recovery metadata is preserved" describes an
  unreachable state — `orphaned` requires `active`, and validation rejects an
  `active` record holding recovery metadata — so the sentence wants restating and
  coverage row 20 wants reclassifying. Coverage row 12 was mislabelled red-first;
  nothing faithful could have reddened it, and the build reclassified it with a
  mutation proof. Both are the reviewer's to edit; this session did not touch the
  spec beyond its status line.
- **Three review findings were accepted as residual risk, not fixed.** `lineSafe`
  admits display-hostile non-control runes (U+202E RLO renders the `--apply`
  disclaimer reversed); the sweep grades `residualAssignment` against a snapshot
  taken outside the ledger lock, so a concurrent `clean --apply` promoting a row
  to `recovered` between snapshot and delete would orphan its recovery ref; and
  interrupt-mid-sweep has no fault injection, mitigated by asserted idempotency.
  Each is a design change the reviewer owns, and none is a regression — all three
  are properties the pre-build code shared or did not reach.
- **Two specs sit unretired on purpose.** `ft148-worktree-orphan-retirement`
  must not be retired before the reviewer has merged it — retirement deletes the
  review surface for work they have not seen. `ft91-gate-phase-split` stays
  unretired for the reason it always has: retiring it destroys the veto surface
  on stories 4, 5, and 9.
- **One open learning is parked for the drain.** The `reviews/` pickup artifact
  was written and deleted inside one session without ever being committed, so it
  was never the tracked state the phase requires. Net tree state is correct; the
  rule was not followed. `.bench/learnings.md` carries the entry.
- **One recovery ref still wants a by-hand look:**
  `refs/bench/recovery/incident-20260712-ambient-probe` matches no assignment.

## Next command

Push is the reviewer's: `git push origin main`. After the push lands, the
retirement row fires and the next session runs
`bench spec retire ft148-worktree-orphan-retirement`, promoting durable content
first. The drain row (`1 idea, 1 open learning`) belongs to `/bench-what-next`.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
