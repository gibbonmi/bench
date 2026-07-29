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

- 2026-07-29  Treated a session-opening venue directive ("I want this session to
  build it with a codex reviewer") as a batch approval and entered
  /bench-implement-spec right after emitting the spec approval table, without
  waiting for sign-off. Right behavior: a directive that names the venue and
  reviewer setup for a build is a plan for the session, not approval of a spec
  that does not yet exist — spec sign-off stays a hard stop unless the reviewer
  has approved a batch plan in terms that cover unseen specs ("roll the
  roadmap"). Proposed rule: none needed; BENCH.md already states it — the miss
  was classification, not a rule gap.
