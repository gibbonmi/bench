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

## 2026-07-14 — FT93(b)/(c) reconcile-vs-preserve is a reviewer decision  [open]

- **What happened:** FT93 shipped part (a) (retain verdict) but I stopped before
  (b) release-reconcile and (c) the orphan sweep. The row assumed the open
  assignment records were residue to sweep. Inspecting the live ledger, all 8
  tree-gone records (7 `recovered` + 1 `cleanup-pending`) each hold a live
  recovery ref under `refs/bench/recovery/` pointing at preserved unmerged work.
  A session-start auto-delete — the plain reading of "sweep" — would sever the
  only ledger pointer to that work.
- **Right behavior:** Treat reconcile/sweep of a tree-gone assignment as
  conditional on preserved work: compact only when no recovery ref exists;
  otherwise surface the exact `bench worktree recovery <ref>` command and let the
  reviewer recover-or-retire. Deleting records that hold preserved work is the
  reviewer's call, not an automatic step. Left (b)/(c) for a reviewer decision
  rather than building the destructive default.
- **Proposed rule change:** none (row-scoped; the recover-vs-discard fork belongs
  in the FT93 follow-up, not a kit rule).
