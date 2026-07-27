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

## 2026-07-27 — the review pickup file was never tracked state  [open]

- **What happened:** `/bench-review-implementation` step 5 says the
  `reviews/<slug>.md` pickup artifact is "tracked state at birth, never untracked
  drift that flips the gate verdict stale", and must be committed in the same
  session that writes it. On FT148 I wrote it and went straight into the repair
  pass in the same session, so it lived and died untracked — the fix delegate
  deleted it and `git status` never showed it at all. Net tree state is identical,
  but a session that had died mid-repair would have lost the findings entirely,
  which is the exact failure the "tracked at birth" clause exists to prevent.
- **Right behavior:** commit the pickup artifact before starting the repair pass,
  even when the repair is immediate, so the findings survive a session death; the
  green fix commit then deletes it as designed.
- **Proposed rule change:** none to the rule — it already says this and I did not
  follow it. Possible ergonomic fix: have the phase name the commit explicitly
  ("commit the artifact, then repair"), since the current wording puts the
  tracked-at-birth requirement in a subordinate clause at the end of a long
  paragraph and it reads as advisory.
