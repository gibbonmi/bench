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

## 2026-07-07 — gate contract test starved under a concurrent delegate build  [open]
- **What happened:** During the FT batch, `bench gate` in the main checkout went
  red on `bench_worktree_concurrent-acquire_contract` ("second acquire did not
  record within a minute") while a write-delegate was compiling and testing in a
  sibling worktree. The same test passed in 0.15s in isolation immediately after.
- **Right behavior:** Treat a single-writer gate run as the verdict for one diff
  *and* one machine-load context: when parallel delegates share the box, either
  serialize full-gate runs or re-run a red timing-shaped contract in isolation
  before believing it.
- **Proposed rule change:** consider raising/removing the fixed 60s overlap
  window in the concurrent-acquire contract (event-barrier instead of wall
  clock), so machine load cannot fake a red.

## 2026-07-07 — specs authored from orchestrator-resolved grills under batch approval  [open]
- **What happened:** FT5 (`bench outline`) and FT7 (`bench dashboard`) needed
  their grills, but the reviewer pre-approved all recommendations for the batch
  and went AFK. The orchestrating session resolved the grill decisions itself
  and handed spec delegates a closed decision handoff, skipping the interactive
  `/bench-shape-idea` interview; specs stay in `specs/` as veto surface.
- **Right behavior:** Under an explicit batch approval this is the "build on
  rather than stall" rule working as intended, but each self-resolved grill
  must be flagged for post-hoc veto (done in the exit report).
- **Proposed rule change:** none.

## 2026-07-07 — FT7 dashboard built data-minimal; taste half deferred  [open]
- **What happened:** The FT7 idea wanted ui_examples-inspired styling and
  animated characters; the referenced image is not in the repo and taste is a
  reviewer call. v1 shipped as a neutral, self-contained HTML snapshot.
- **Right behavior:** Ship the data-faithful core, park the visual identity as
  its own reviewable item instead of guessing taste AFK.
- **Proposed rule change:** none.
