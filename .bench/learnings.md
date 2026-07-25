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

## 2026-07-24 — `bench commit --spec` always leaves the gate reporting strong-stale  [open]

- **What happened:** Landing FT91's first arm with `bench commit --spec
  gate-concurrency-budget` produced a green gate and a clean commit, but
  `bench status` immediately reported `gate stale (gated tree …, work tree …) →
  re-run the gate`. The drift was `specs/gate-concurrency-budget.md`: `--spec`
  flips `Status: staged` to `implemented` as part of the commit, necessarily
  *after* the gate has passed, so the gated tree can never contain the flip.
  `specs/*.md` is not in the capture-only allowlist, so the row fails closed to
  the strong stale wording rather than softening. This is structural, not
  specific to this spec — every `--spec` landing ends this way, and the row tells
  the next session to re-run a ~6-minute gate that would find nothing.
- **Right behavior:** Unclear, which is why this is captured rather than fixed. A
  single-line spec status flip is exactly the "capture-only drift" case the
  softening was built for, but widening the allowlist to `specs/*.md` would admit
  arbitrary spec edits, which the fixed exact-path allowlist deliberately refuses
  (the profile calls expanding it a new decision). A narrower option is for the
  gate cache to record the post-flip tree when `--spec` performs the flip, since
  the flip is the tool's own write and its content is known.
- **Proposed rule change:** none — this is a tooling behavior, and the fix is a
  reviewer decision between widening the allowlist, teaching the cache about the
  `--spec` write, or accepting the row as cosmetic and documenting it.
