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

## 2026-07-07 — merge-time done-claim verification thinned as the batch wore on  [open]
- **What happened:** During the twelve-spec batch, every delegate merge got
  `git status` on the worktree, a commit-stat scope check, and the full gate on
  the merged tree — but independent spot-checks of the claimed behaviors
  (re-running a red signal, probing the built binary) were done for the early
  merges and skipped for the later ones. Since delegates author the very tests
  that pin their behaviors, gate-green alone can't distinguish a correct build
  from a self-consistent spec-divergent one. The reviewer's question flushed
  the gap; post-hoc probes of every thinly-verified surface (line-guard deny,
  structure loud error, models argv, split-survival of concurrent hunks,
  unlink end-to-end) all passed, so no defect shipped — this time.
- **Right behavior:** Merge-time verification of a write-delegation includes at
  least one independent behavioral probe per spec — exercising a claimed
  behavior through the built binary or a throwaway fixture the delegate did not
  author — alongside the gate, `git status`, and the scope check. Constant
  across a batch, not front-loaded.
- **Proposed rule change:** Add one sentence to `craft-delegate`'s
  "Verifying the done-claim" list: at least one accepted behavior is probed
  independently of the delegate's own tests before the merge is accepted.
