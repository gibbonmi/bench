# Retro — native runtime evidence reduction

## Outcome

The landing published `e1299080` from source tip `4b7877a7` onto base
`cd14fb9c`. The gate ran green across all six phases. The census recorded 110
raw calls in the integration assignment.

The `native runtime` workflow now builds one artifact generation instead of
two, and it publishes one upload instead of four. The `artifacts` job cost
falls from four complete builds to two, because `scripts/build-artifacts.sh`
still runs its own second build. The byte-for-byte artifact comparison
therefore survives. The cross-checkout comparison of finalized release
evidence is retired.

The release plan now states, per target, whether that target carries a native
proof. The two Darwin targets do not. No macOS runner starts for
`native-proof`, and the `smoke` matrix keeps all four shipped targets. Every
proof count derives from the proven list. The proof script refuses a target it
cannot verify.

The landing also closed an open defect. The Darwin strip assertion could never
pass on a loadable Mach-O, and it left the tree with its branch.

## Gate-stage timings

| stage | verdict | elapsed |
| --- | --- | --- |
| gofmt | green | 86 ms |
| vet | green | 1100 ms |
| test | green | 52680 ms |
| race | green | 2574 ms |
| system | green | 11121 ms |
| shellcheck | green | 491 ms |

The ship-tier release-evidence probe ran separately at about 70 seconds. It is
the only check that exercises finalization end to end.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized. No charge received a spec slice, so this
landing offers no contrast between the two shapes.

Eight write charges ran, six from the build and two from the repair pass. Six
landed first-pass on behavior. Two returned work that needed a second round,
and both second rounds came from a defect in my spec rather than from the
delegate.

Two charges did better than their ticket asked. The collapse charge stopped
and reported a fifth consumer of the artifacts upload that my spec never
listed, instead of widening its own scope. The proof-count charge reported an
acceptance row it could not close with real evidence, rather than claiming it.

One charge reported a read-only command it ran through a shell `cd` instead of
`bench worktree exec`. It disclosed the deviation without being asked.

## Coordinator catches

I verified every done-claim against the tree and re-ran the focused checks
myself. Three catches were material.

The first delegate's diff changed a shared export signature. I swept for
callers before folding and confirmed the export had none outside its file.

Folding the third ticket produced a merge conflict in the conformance checks,
because two tickets edited the same function. I resolved it by keeping both
additions rather than taking one side.

I re-read the assembled workflow end to end after the fifth fold, and found a
vestigial upload name that no longer describes anything. I parked it rather
than widen the landing.

The honest entry on the other side: I missed the ungraded matrix producer
twice. I wrote the ticket that asked for the check. I also reviewed the diff that
delivered it. Neither time did I ask what single edit would defeat the check
while leaving it green. Two review axes found it independently.

## Repair attribution

| ticket | rounds | cause per round |
| --- | --- | --- |
| add-proven-target-field | 0 | none |
| bind-the-proof-scripts-to-proven-targets | 1 | spec-row |
| count-release-evidence-against-proven-targets | 0 | none |
| switch-the-proof-matrix-to-proven-targets | 1 | spec-row |
| collapse-the-artifacts-job-to-one-generation | 1 | spec-row |
| restate-the-release-claims-in-the-docs | 0 | none |
| repair-the-ungraded-gate-anchors | 1 | other |
| repair-the-proof-paths-that-fail-open | 1 | spec-row |

Four of five repair rounds trace to a spec row I wrote. Three of those share
one root cause: I traced the producers of a changed fact and not its readers.
The review's first round found two unlisted readers. A delegate found a third.
The fifth round traces to my own accepted finding, which removed a test that
carried coverage nothing else held.

## Agent-experience improvements

### Bench CLI

- Add a verb that folds one owned worktree's diff into another owned worktree, refusing on conflict with the conflicting paths named.
  Feeds: new
- Add a verb that brings an owned worktree onto the current default-branch tip, because no verb does this and the build needed it twice.
  Feeds: new
- Make the `bench commit` lane's `prose` check run the same bound as the gate's live-tree prose test, or stop naming the lane check `prose`.
  Feeds: new

The census entry for this landing is
`native-runtime-evidence-reduction: 110 raw calls in the integration assignment`
in `capture/learnings.md`. It names the same two integration verbs.
Feeds: new

### Skills

- Teach `craft-spec` to sweep for every reader of a derived value before the coverage map locks.
  Three repair rounds here traced to unlisted readers.
  Feeds: new
- Teach `craft-gate` that a check on an indirected value grades both ends, because a check that names an output proves nothing about what fills it.
  Feeds: new
- Teach `craft-delegate` that a charge naming an exact verification command must have run it once.
  One charge here named a ship-tier command that silently skips.
  Feeds: new

### Process

- Ask what single edit would defeat a new check while leaving it green, before accepting the check.
  Feeds: none
- Treat a Standards finding that deletes a test as a coverage question first.
  One such finding here removed the only assertion binding the plan's proven set.
  Feeds: none
- Run the ship-tier release-evidence probe on the integration source before the landing, because the dev gate does not reach finalization.
  Feeds: none
