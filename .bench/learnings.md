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

## 2026-07-10 — Reviewer-pre-authorized grill answers  [open]
- **What happened:** /bench-shape-idea for COMPLIACE_ASSESSMENT.md ran with the reviewer's explicit instruction to adopt all my grill recommendations and write the map as if the full grill had run; grill tickets were self-answered, breaking the live-exchange contract by direction.
- **Right behavior:** Honor the instruction, mark provenance in the map header so it reads as pre-authorized recommendations (veto surface), and flag the most contestable answers at close.
- **Proposed rule change:** none — reviewer instruction governs; consider codifying "pre-authorized grill" as a named mode if it recurs.
