# shift in worktree

## Problem

`CONTEXT.md`, `README.md`, and the `bin/bench.sh` header all promise that a **shift**
runs in an isolated **worktree** "without touching the main checkout." It does not.
`shift_loop` runs in place in the main checkout: it `git switch -c`es the checkout onto
a `bench/shift-<ts>` branch, runs the agent there, and on a red iteration runs
`git reset --hard` + `git clean -fd` **on the main checkout**. The pooled worktree the
docs describe exists only behind the separate, manual `bench worktree` subshell, which
the loop never calls.

The gap has two costs. A user who trusts the promise runs `bench shift` in their main
checkout expecting a sandbox, and instead finds their checkout switched to another
branch with destructive git run against it. And the isolation the loop *should* give —
run autonomous, gate-rollback-heavy work somewhere the human's checkout can't be
disturbed — is absent. `decisions/dogfood-improvements.md` ticket #3 framed the fork
("make shift own the worktree, or narrow the docs"); the reviewer chose to make the
promise true.

## Solution

`bench shift` acquires a warm worktree from the pool, runs the **entire** gated loop
inside it against the main checkout's committed `HEAD`, and leaves the resulting
`bench/shift-<ts>` branch — with its commits — for the reviewer to inspect and merge.
The main checkout's branch, index, and working tree are never switched, reset, or
cleaned. From the user's seat: `bench shift "<objective>"` returns them exactly where
they started, plus a branch to review; `bench worktree` still drops them into an
interactive pooled subshell as before.

## User stories

1. As a user running a shift, I want it to execute in a pooled worktree so my main
   checkout is never switched, reset, or cleaned.
2. As a user, I want my main checkout to stay on its original branch with its working
   tree untouched for the whole shift, including on red-iteration rollback and on Ctrl-C.
3. As a user, I want the shift to build on my current committed state (main `HEAD`), so
   the work is relevant to where I am — not on a stale `origin` default.
4. As a user, I want the finished `bench/shift-<ts>` branch to remain after the shift,
   with its commits, and be checkout-able/mergeable from my main repo (no worktree still
   holding it).
5. As a user pressing Ctrl-C to pull the line, I want the pooled worktree released
   (unleased and cleaned) and no scratch files or partial branch state anywhere in my
   main checkout.
6. As a user, I want every existing shift behavior preserved exactly — clean-tree
   precondition, commit-on-green, rollback-on-red, touched-scope refactor phase, early
   stop via `.bench/done.sh`, iteration cap, truthful commit-count summary — with only
   the *location* of the work changed.
7. As a user of `bench worktree` (interactive), I want it unchanged: the pool
   acquire/lease/reset/release logic is now shared with `bench shift`, but the
   subshell-and-release behavior is identical.
8. As a kit developer, I want the shift CLI contracts updated to assert the isolation
   (main checkout untouched, shift branch present with the expected commits, worktree
   released) and still exercise the real `bin/bench.sh shift` in throwaway repos.

## Implementation decisions

- **Extract the pool as a seam.** Pull the pool acquire/lease/reset logic and the
  release logic out of `worktree()` into two helpers (an acquire that returns a leased,
  clean worktree path; a release that unleases and resets it for reuse). Both
  `worktree()` (interactive) and `shift_loop` (programmatic) call them. `worktree()`'s
  external behavior is unchanged.
- **Run the loop in the worktree.** `shift_loop` captures the main checkout's `HEAD`,
  acquires a pooled worktree, resets it to that `HEAD` sha (not `origin/<default>` — the
  shared object DB makes the sha always resolvable and keeps the shift on the user's
  current state), creates `bench/shift-<ts>` there, and runs the agent, gate, commit,
  rollback, and refactor phase entirely against the worktree. The main checkout is never
  a `git` target for switch/reset/clean.
- **Scratch and touched-scope move with the loop.** `.bench-objective` / `.bench-notes.md`
  live in the worktree; the commit exclude-pathspec, the trap cleanup, and
  `structure_touched_since`'s `<base>..HEAD` diff all operate against the worktree.
- **Preserve the branch on release.** When the loop ends (completion or interrupt),
  detach the worktree's `HEAD` (leaving the `bench/shift-<ts>` ref pointing at the work),
  then release the worktree to the pool. Because no worktree then holds the branch, the
  reviewer can check it out or merge it from the main repo.
- **Precondition retained.** The clean-tree check stays, now on the main checkout only
  (shift builds on committed state; carrying uncommitted work is a separate capability —
  see Out of scope). The check no longer needs the scratch-file exemption in the main
  tree, since scratch lives in the worktree.
- **Docs become true, not edited to retreat.** With the behavior shipped, the existing
  `CONTEXT.md` / `README.md` / header "runs in a worktree" statements are correct;
  tighten them only to say *pooled* worktree and *main checkout untouched*. Record the
  resolution in `decisions/dogfood-improvements.md` #3.

## Testing decisions

- **A good test here** asserts external behavior at the `bench shift` CLI seam: after a
  shift, the main checkout is on its original branch and clean; `bench/shift-<ts>` exists
  with the expected number of commits; the pooled worktree is released (no leftover
  lease, reused on the next shift). It does **not** reach into loop internals.
- **Seam + prior art:** the existing `bench shift` contract blocks in
  `.bench/gate-runtime-contracts.sh` (throwaway repo + controlled `.bench/gate.sh` +
  tiny shell agent) — update them so their assertions read against the shift branch and
  the untouched main checkout instead of the main `HEAD`. The `bench worktree`
  lease/reuse contract must stay green after the acquire/release extraction. Add one
  isolation assertion: the main checkout's branch name and porcelain status are identical
  before and after a shift.
- **Coordination:** these contracts live in `.bench/gate-runtime-contracts.sh`, which the
  parallel testing-improvement work also touches — implementation must rebase onto that
  work rather than clobber it.
- **Gate command:** `bench gate`.

## Out of scope

- **Shift from a dirty main checkout** (snapshot/stash uncommitted work into the shift
  worktree) — a separate capability with its own carry-semantics decision, not the rest
  of this one; ~30–45 min later.
- **Auto-removing merged shift branches and their pooled worktrees** — a separate
  lifecycle capability; ~30 min later.
- **The portable harness command contract** (dropping the hardcoded `-p` flag for
  Codex/OpenCode) — that is ticket #4 in `decisions/dogfood-improvements.md`, its own
  research spec.
