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

## 2026-07-29 — no sanctioned way to land two unrelated dirty changes separately  [open]

- **What happened:** The main tree held two unrelated uncommitted changes: a
  capture-only `IDEAS.md` line parked earlier in the session, and a kit edit to
  `.bench/BENCH.md` plus its `CHANGELOG.md` entry. `bench commit` refuses to run
  while any dirty file sits outside the named set, and `block-dangerous-git`
  correctly refuses `git stash`, so there was no route to two commits. I bundled
  all three paths into one commit, led the message with the kit edit, named the
  ride-along explicitly, and offered the reviewer a redo. It landed at `5fd3789`.
- **Right behavior:** Genuinely unclear, which is why this is captured rather than
  fixed. Bundling with disclosure is defensible and cost one gate run instead of
  two, but it is a convention I invented in the moment, and it does bend "one small
  change at a time." The alternative was to stop and ask over a one-line capture
  file, which spends reviewer attention on something close to trivial.
- **Proposed rule change:** Probably not new machinery. The block-check exists so a
  green verdict describes the whole tree, and a "gate this subset, leave the rest
  dirty" flag would hollow it out — the same hazard the guard message now names.
  The cheaper resolution is to write the convention down: when the block-check and
  the git guard leave no separate path, bundle, lead the message with the
  substantive change, name the ride-along, and flag it for veto. Worth deciding
  alongside the capture-only commit-path idea parked in `IDEAS.md` the same day,
  since both are about what a capture file costs to land.
