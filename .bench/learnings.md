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

## 2026-07-10 — Emitted CLI usage is behavior  [open]
- **What happened:** A `bench commit` usage-output change was classified as
  prose/usage-only, so the synthesis batch substituted one green gate for the
  required dogfood run.
- **Right behavior:** Treat emitted CLI output as behavior and run the real
  `craft-synthesis` dogfood loop whenever a CLI is touched.
- **Proposed rule change:** none — `craft-synthesis` already requires dogfood for
  every CLI or behavior change.

## 2026-07-10 — Review diff omitted the working tree  [open]
- **What happened:** Exact `bin/bench.sh diff --full` returned `files[0]` for an
  uncommitted default-branch batch while `git diff HEAD` showed its 11-file review
  bundle.
- **Right behavior:** When the canonical review diff is empty beside a dirty
  working tree, explicitly pin the bundle with `git diff --no-ext-diff HEAD`.
- **Proposed rule change:** Add that fallback to review guidance, or make `bench
  diff` include working-tree changes on the default branch.

## 2026-07-10 — Bound the terminal repair pass  [open]
- **What happened:** The review/fix cycle recursively opened new review and repair
  rounds, became unbounded, and created a poor user experience.
- **Right behavior:** Declare one terminal repair pass: integrate accepted
  findings, run focused checks, run one final gate, then stop and report. Open a
  new review round only when that gate fails or the user requests one.
- **Proposed rule change:** Add the terminal repair-pass bound to review guidance.
