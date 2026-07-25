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

## 2026-07-25 — a write-delegate's `git stash` crosses worktree boundaries  [open]

- **What happened:** During the FT86 build I ran two write-delegates concurrently
  in two separate `bench worktree` checkouts, charged with disjoint packages. Both
  independently reached for `git stash push` / `git stash pop` to restore the
  pre-change file and observe a behavioral red — a technique nothing in the kit
  warns against. Git's stash ref is repository-global, not per-worktree, so the
  two stacks were one stack: each delegate's pop applied the other's in-flight
  changes into its own checkout. Both worktrees ended up holding both slices'
  edits with byte-identical content, and neither diff could be attributed. One
  delegate surfaced it only as a symptom it could not explain ("my file reverted
  to HEAD on its own"). I discarded one slice's work and re-ran it. `main` was
  never touched and no mixed diff reached a gate.
- **Right behavior:** The charge should have banned `git stash` outright and named
  the safe per-worktree substitute — `cp` the working file aside, then
  `git show HEAD:<path> > <path>` to restore the original, test, and copy back.
  That touches only the delegate's own checkout. More generally: every
  repository-global git surface (the stash stack, refs, config, the index of a
  shared gitdir) is a shared-state hazard the worktree isolation does *not* cover,
  and the isolation guarantee is worth stating with its actual boundary rather
  than as "the delegate gets a checkout, not your checkout."
- **Proposed rule change:** `craft-delegate`'s Isolation section claims worktrees
  keep concurrent writers from colliding. That is true for the working tree and
  false for the stash. Add the boundary and the banned command to that section, and
  give the copy-based red-proving technique as the sanctioned alternative — it is
  the thing delegates actually want when a charge asks them to prove a red. The
  `block-dangerous-git` hook already refuses `git checkout <path>`; `git stash` is
  the same class of hazard for concurrent work and is a candidate for the same hook.

## 2026-07-25 — a stale `dist/bench` makes contract rows lie in both directions  [open]

- **What happened:** After rebuilding a salvaged slice into a clean worktree, two
  of its three AXI contract rows failed. The source was correct; the worktree's
  `dist/bench` predated the change, and those tests drive the built binary rather
  than the package under test. Rebuilding with `scripts/go-build.sh` turned all
  three green with no source edit. I nearly re-charged a delegate to fix code that
  was already right.
- **Right behavior:** Rebuild the binary before running any suite that drives it,
  and treat an unexplained contract red in a fresh or salvaged worktree as a stale
  binary until proven otherwise. The dangerous direction is the opposite one: a
  stale binary that happens to satisfy an assertion makes a *broken* change look
  green, which is a false done-claim the gate would catch only later.
- **Proposed rule change:** The build phase tells delegates to run the relevant
  single test file frequently; for the AXI and runtime contract seams that
  instruction is incomplete without the rebuild step. Either say so where the
  phase names those seams, or have the contract helper itself fail loudly when the
  binary under test is older than the package sources it is meant to exercise —
  the second is the single-source fix, since it removes the instruction entirely.
