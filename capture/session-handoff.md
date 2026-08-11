# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `35e887d4`, clean of the retired-spec noise
Spec: none active — `decisions/pocock-alignment.md` (Status: ready) is the
build program; all seven previously staged specs are dispositioned by its #2

## State

The Pocock-alignment shaping closed 2026-08-11 with all 13 decision tickets
resolved and the map `ready`. The program: full removal of the provisional
spec-build lifecycle, doctrine-leaf adoption, and an artifacts-vs-reality
preflight. Closed decisions (do not reopen): full lifecycle removal, nothing
surviving (#1); retire `axi-spec-build-complete`, `checkpoint-scoped-review`,
`ticket-bundle-refusal`, `axi-compatibility-oracle` and park
`axi-coherent-diff`, `axi-query-disclosure`, `single-build-serial-gate` (#2);
only `Blocked by:` stays parsed in tickets (#3); `bench-craft-domain`
companion skill (#4); four doctrine-leaf adoptions including frontier-rounds
grilling (#5); review re-derives from primary sources (#6); `bench preflight`
at phase entry, not in the gate (#7); seam-confirmation and breakdown-quiz
HITL gates replace the delegate breakdown review (#8); hard line budgets
gate-enforced (#9); craft-line and craft-delegate kept slimmed (#10); CLI
shrink per #11; zero backwards compatibility (#12); three specs in order
A lifecycle-removal → B preflight → C doctrine, after housekeeping (#13).

Housekeeping executed (`35e887d4`): the four retired specs deleted with their
tickets and dirty edits (`bench spec retire` refused — it only takes
implemented specs — so deletion into git history was the route). The three
parked specs stay staged in `specs/`; their roadmap re-rank belongs to the
next `/bench-what-next` drain. The abandoned run's ~20 recovery refs and open
assignments are lifecycle debris that Spec A's removal deletes wholesale.
Reviewer authorized direct-to-main commits without the gate for this program.

## Next command

`/bench-write-spec` for Spec A (lifecycle removal) from
`decisions/pocock-alignment.md` tickets #1, #3, #11, #12.

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
