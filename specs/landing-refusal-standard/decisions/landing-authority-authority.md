# Landing authority split (FT169 research)

Source: read-only Terra delegate, 2026-09-01. Coordinator spot checks:
`internal/usage/worktree.go:54`, `internal/worktree/ownership.go:166`,
`internal/landing/gitexec.go:20` — all resolve.

## Summary

Five surfaces, three authorities. `bench worktree land` is the only verb
that writes `main`. `bench commit` writes a branch tip and refuses the
primary checkout. `bench spec retire` deletes files and commits nothing.
No pool cap on active worktrees exists.

## bench worktree land

Entry `internal/worktree/land.go:155`; `--resume` branches into
`land_resume.go:29`. The mutations run in order. The landing writes the
two-parent merge commit and moves the ref by compare-and-swap
(`internal/landing/landing.go:272-276`). It flips the spec in the
published tree only (`landing.go:241-249`). It closes the tickets-only
folder (`landing.go:253-257`).

It advances the green marker
(`land.go:297`). It resets the destination with `git reset --merge`
(`land_identity.go:213-224`). It releases the worktree (`land.go:304`),
and the census record drops inside the release step
(`lifecycle.go:437-439`). It writes no
board row and no handoff; the handoff joins only as a merge rule where the
source side wins (`internal/landing/composition.go:59`).

One preflight collects every refusal before the first output line
(`land.go:230-258`). The refusal families are destination proofs
(`land_identity.go:35-62`), the six identity components
(`identity_component.go:33-38`), and source proofs
(`land_identity.go:91-160`). Range proofs (`:172`, `:179`), base proofs
(`land.go:144`, `:251`), and owner proofs
(`internal/landing/landing.go:201-270`) follow. The composition conflict
(`composition.go:37`) and the compare-and-swap loss (`gitexec.go:20`)
end the set.
The landing never invokes `bench commit` or `bench spec retire`.

## bench commit

Entry `internal/commit/commit.go:64`. It mutates one thing: the current
branch ref (`commit.go:109-111`, `:173`), after a `gofmt` rewrite of named
Go files (`commit.go:136`). It refuses the primary checkout
(`commit.go:102`). A declared fast lane replaces the whole-project gate
for this commit (`commit.go:130`, `:151`). Exit 3 reports a published
commit whose checkout did not reconcile (`commit.go:196`).

## bench worktree release

Entry `internal/worktree/worktree.go:444`. It moves the assignment from
active to cleanup-pending and writes receipts. It then retires the
checkout and branch, drops the census record, and deletes a complete
assignment (`ownership.go:342-437`, `lifecycle.go:432-438`). A retained worktree
prints `worktree retained (<reason-code>)` with the re-run command
(`ownership.go:463-466`). The landing, the landing resume, and the Claude
WorktreeRemove hook invoke it.

## bench spec retire

Entry `internal/spec/spec.go:322`. It mutates the filesystem only: the
review pickup, the tickets folder, the spec file, the folder
(`spec.go:413`). It never commits and never runs the gate. It refuses the
primary checkout (`spec.go:334`). The board remainder stays human work:
it prints the `next:` line naming the ROADMAP row and the
`spec-retire: <slug>` commit (`spec.go:391`).

## Primary-checkout refusal

One source: `internal/usage/worktree.go:50-55`, rendered as
`error: primary checkout is read-only for Bench phases — run bench
worktree create ...`. Three callers: `commit.go:102`, `spec.go:334`,
`roadmap/learning.go:99`. The landing does the opposite: it requires a
clean primary checkout attached to the default branch
(`land_identity.go:33-51`). The landing is the one verb that runs on the
primary checkout and writes `main`.

## The one-active limit

No code counts active worktrees or landings. Four mechanisms give the
effect:

- One assignment per request digest; a replay returns the existing
  assignment, and a disagreeing replay refuses
  (`internal/worktree/ownership.go:162-172`).
- An exclusive flock on `bench-create.lock` over creation
  (`ownership.go:113-128`).
- The Claude hook derives the request id from `session_id`, so one session
  replays one worktree (`internal/harness/worktree.go:112`).
- One landing at a time by compare-and-swap, not by lock; a loser
  recomposes (`internal/landing/landing.go:263-278`,
  `gitexec.go:19-21`).

## Read and not read

The delegate read the worktree land/release/identity code, commit, spec,
usage, harness-hook, and parts of the landing package and
`.bench/BENCH-reference.md:168-208`. It did not read the test files, the
gate, status, or handoff internals, or the landingpolicy and
lifecyclepolicy bodies. The exit-code and capture-merge claims rest on
call sites and the reference doc.
