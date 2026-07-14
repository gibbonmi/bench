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

## 2026-07-14 — manual worktree create/release pair refuses release  [open]

- **What happened:** During the FT90 debug close-out, a session-created worktree
  (`bench worktree create --request ft90-fix-a46ad5 --label ft90-gate-subject`)
  could not be released after its branch merged to main: `bench worktree release
  --request ft90-fix-a46ad5 <path>` failed with "terminal receipt missing", and
  after `bench worktree clean --discard-ignored <path> --apply <fingerprint>`
  removed the tree, release still failed with "cleanup receipt does not
  authorize release reconciliation". The worktree itself is gone; the
  assignment record's fate is unknown (session start reported 10 open
  assignments).
- **Right behavior:** Unclear — either the `--request` create/release pair is
  adapter plumbing sessions shouldn't type (then BENCH.md's CLI inventory
  placement is misleading and the bare interactive `bench worktree` was the
  right tool), or a merged-and-cleaned assignment should be releasable and this
  is a defect in release's receipt reconciliation.
- **Proposed rule change:** none until triaged; repro is the exact command pair
  above run outside a shift.
