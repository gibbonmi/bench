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

## 2026-07-23 — Reproduce load-sensitive gate failures with the real oracle  [open]

- **What happened:** Only a real `bench gate` under representative host-side load reproduced the marker stalls; guest-side CPU saturation, parallel contract loaders, inert memory ballast, and `fsync` hammers all stayed green.
- **Right behavior:** Diagnose load-sensitive gate failures with the real gate under the host load that exposes them; synthetic guest stress can narrow a hypothesis, but a green approximation cannot clear the real oracle.
- **Proposed rule change:** Add this reproduction-economics rule to gate-debugging guidance so a fresh session does not spend cycles rebuilding disproven guest-only load shapes.
