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

Format per entry:

## <date> — <short title>  [open]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry leaves this file only via /bench-what-next.

<!-- entries below -->

## 2026-07-06 — delegate worktrees cut from a stale ref  [open]
- **What happened:** During the FT13–FT21 batch build, Agent-tool worktrees were
  repeatedly created several commits behind main: one delegate stopped and
  reported instead of building, a second was denied its own `git merge
  --ff-only` by the permission layer, and a third built against the stale base
  (its diff still applied, but it couldn't reuse already-landed sources and
  re-flagged fixed behavior).
- **Right behavior:** Every worktree-isolated delegate charge opens with "run
  `git merge --ff-only main`, verify HEAD equals main, stop and report if
  denied" — adopted mid-run and it held; a blocked worktree is fast-forwarded
  by the orchestrating session, which then resumes the same delegate.
- **Proposed rule change:** Add that opening action to `craft-delegate`'s
  charge template so worktree delegates never build on a stale base.

## 2026-07-06 — assessment drained to specs without decision maps  [open]
- **What happened:** The reviewer invoked /bench-write-spec on ASSESSMENT.md
  with an explicit batch approval ("write all of them and commit until the
  assessment is drained", defaults delegated, reviewer AFK). The phase's entry
  contract says refuse without a closed map, and the profile routes ordinary
  spec authoring to the mid tier; nine specs were authored map-less on the
  invoking top-tier session, with every defaulted decision flagged in-spec for
  post-hoc veto.
- **Right behavior:** Honor the explicit reviewer override, keep the veto
  surface (flagged defaults, staged specs, disposition table) — but the
  contract text and the override path shouldn't have to be reconciled ad hoc
  each time.
- **Proposed rule change:** Add one sentence to /bench-write-spec's entry
  contract: an explicit reviewer-directed batch drain (assessment or reviewed
  findings doc → specs) may substitute for per-spec maps, with every defaulted
  decision flagged for veto; absent that explicit instruction, the map gate
  stands.
