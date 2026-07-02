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

## 2026-07-02 — canary promise read as per-class, review caught the gap  [open]
- **What happened:** the state-surface spec promised "canary fixtures proving each
  check bites"; the build slice added one canary for the guard-manifest class and
  cited precedent; the three-axis review flagged ~9 CLI-shape checks with no
  canary; grouped canaries were added in the fix round.
- **Right behavior:** read "each" in an oracle promise literally at build time, or
  surface the per-class interpretation before building instead of after review.
- **Proposed rule change:** none — but spec authors should name canary granularity
  explicitly (per check vs per class) so the build can't misread it.
