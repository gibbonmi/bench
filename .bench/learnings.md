# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

## 2026-07-31 — a reviewer-approved drain was blocked three times, once by its own defect  [open]

- **What happened:** One `/bench-what-next` batch took three `bench commit`
  attempts. First, `TestHandoffShapeSingleSourced` was already red on `4a756b3`;
  reproducing it at HEAD in a clean worktree showed the drain had not caused it,
  and `bench handoff` was the fix the check itself named. Second, reconciling
  shipped row FT128 out of `ROADMAP.md` turned `checkOccurrenceLedgerMigration`
  red, because its `want` map pins that row's occurrence count and no phase
  instruction mentions the map — this one was the drain's own defect. Third, a
  concurrent session writing the same checkout tripped `bench commit`'s
  whole-tree attribution refusal twice: once on untracked `decisions/` files from
  an in-flight `/bench-shape-idea`, and once on a mid-drain `bench idea` write to
  `IDEAS.md`. The attribution refusal is cheap and fires before the gate; the two
  full gate runs were spent on the first two reds, before the concurrent writer's
  files existed at all.
- **Right behavior:** Attributing each red before reacting was right and is what
  kept a check from being edited to go green. Beyond that, two of the three have
  no settled answer. Removing a shipped row should verdict the occurrence-ledger
  map in the same pass — the check's own bite test already establishes that a
  retired FT leaves the map, so this is maintenance the phase omits rather than a
  judgment call. Landing beside a concurrent writer has no sanctioned sequence at
  all: invariant 1 says take side-work to a worktree, but a drain is not side-work
  and its diff already lives in the main checkout, so the advice does not reach
  this case.
- **Proposed rule change:** (1) `/bench-what-next` step 1 names the
  occurrence-ledger migration map as part of removing a shipped row, so the
  reconcile and the check that pins it stay one pass. (2) The kit names one
  sanctioned sequence for landing a batch beside a concurrent session in the same
  checkout — or states that the drain assumes a single writer and how to detect
  that it does not. FT169 already owns the foreign-dirty face and is the likely
  home; FT168 owns the adjacent question of whether the oracle may answer for less
  than the whole tree, but the blocker here is `bench commit`'s attribution check
  rather than the gate's scope, so they are neighbours and not the same row.
