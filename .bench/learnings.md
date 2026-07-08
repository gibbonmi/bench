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

## 2026-07-08 — isolation worktree gave a write-delegate a stale base  [open]
- **What happened:** An FT45 fix delegate ran in a harness-created isolated
  worktree cut at a commit five behind main (before the very code it was
  charged to fix had landed). It faithfully rebuilt the missing mechanism from
  the charge's prose; the diff could not apply onto main and was ported by the
  orchestrator inline.
- **Right behavior:** On receiving a write-delegate's worktree, the
  orchestrator checks the worktree's merge-base against the expected tip
  before reading the diff; the charge should also name a sentinel the
  delegate must confirm exists (a function or test added by the commit under
  fix) so a stale snapshot fails fast instead of producing a divergent
  rebuild.
- **Proposed rule change:** add the sentinel-precondition line to the
  craft-delegate charge template for fix-pass delegations.
