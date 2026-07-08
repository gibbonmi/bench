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

## 2026-07-07 — FT41 dogfood shortfall: no decision map to resume  [open]
- **What happened:** The shape-idea-grill-continuation spec's testing decisions
  call for dogfooding via a real `/bench-shape-idea` resume; the repo has no
  `decisions/` map to resume and the only pending grill (FT38) needs the
  reviewer present, so the build shipped on the bite test plus a green gate
  with a read-through in place of the live resume.
- **Right behavior:** Uncertain — either the synthesis dogfood clause scales to
  "first real use" when no substrate exists, or the build waits for the next
  real grill.
- **Proposed rule change:** none (judgment call flagged for the reviewer).
