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

## 2026-07-26 — build overrode a spec's fixture location without waiting  [open]

- **What happened:** The FT96 spec's "Implementation decisions" placed both new
  canary fixtures under `tests/canary/docs-currency-token-diet/`, but its own
  "Prior art" line pointed at `tests/canary/workflow-guidance-anchors/`, where
  every other anchor fixture lives. The coordinator put the anchor fixtures in
  the anchors family instead — both families are conformance-owned with the
  same shell source, so only discoverability differed. The override was flagged
  to the reviewer at the time but made without waiting for an answer, so the
  spec text and the tree disagreed until the spec was retired.
- **Right behavior:** Correct-and-flag was right. When a spec contradicts
  itself and the two readings are functionally equivalent, the build should
  follow the reading consistent with the tree's existing convention, flag the
  choice, and continue — stopping spends a reviewer round-trip on a decision
  with no behavioral stake, which is the case flag-for-veto exists for. The
  flag is not optional: a silent override of spec text is always wrong, and
  a divergence with any behavioral difference is always a stop-and-ask.
- **Proposed rule change:** none — the batch-approval clause in `.bench/BENCH.md`
  already grants build-on-with-veto-surface; this entry asks the reviewer to
  confirm it covers internally-inconsistent spec text, not to add a rule.
