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

## 2026-07-07 — pre-ran bench gate before bench commit, paying the 2-minute gate twice  [open]
- **What happened:** The drain close-out ran `bench gate` to green, then `bench
  commit` — which runs the full gate itself, unconditionally. The reviewer waited
  ~5 minutes for a three-file docs commit.
- **Right behavior:** Never pre-run the gate before `bench commit`; commit is the
  gate run. (Separately diagnosed: commit ignores a fresh green verdict for the
  identical tree — candidate fix under review.)
- **Proposed rule change:** none (work-shaped fix pending reviewer decision).

## 2026-07-07 — paused a grill to ask for a re-prompt between tickets  [open]
- **What happened:** After resolving canary-cost #6, the session stopped and told
  the reviewer to say go (or re-invoke `/bench-shape-idea`) before grilling #2 —
  a one-question ticket the same sitting could have carried straight through.
- **Right behavior:** Once a grill is running, keep asking question after
  question until no unresolved question remains; never pause to ask permission
  to continue or require a re-prompt between tickets.
- **Proposed rule change:** In `/bench-shape-idea` resume mode, soften "resolve
  that one ticket, then stop" to allow continuing into newly-unblocked grill
  tickets in the same sitting when the reviewer is present and answering.
