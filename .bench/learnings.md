# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-integrate-learnings` reviews the open entries, promotes the
generalizable ones into the kit with sign-off, and prunes them: a resolved entry
leaves this file, and its verdict (promoted or dismissed, one line of why) is
recorded in the integration commit and CHANGELOG. The journal holds open entries
only; history lives in git. Never rewrite a kit rule yourself — that is the whole
point of capturing here instead.

Format per entry:

## <date> — <short title>  [open]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry leaves this file only via /bench-integrate-learnings.

<!-- entries below -->

## 2026-07-02 — review-axis delegate token caps underestimated ~2x  [open]
- What happened: declared the review line as three opus/medium delegates at ~30k tokens each; actuals ran 50–62k each (~166k total vs ~90k declared). No stop-and-report fired because each delegate is a single uncapped run, not an iterating stage.
- Right behavior: estimate ~55–65k per read-heavy review-axis delegate (they read the full diff + standards docs + run verification commands), and say so in the declared line.
- Proposed rule change: none to the kit; a cached routing note in projects/benchkit.md Lines (review-axis delegate ≈ 60k on mid) would make the next declaration honest.
