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

## 2026-07-27 — the write-delegate requirement met a harness-level prohibition  [open]

- **What happened:** `/bench-implement-spec`'s "Route the venue" says every
  spec-backed run assigns genuine write work to at least one write subagent before
  the first implementation edit, with no inline threshold of its own. This session
  carries a standing harness instruction forbidding the Agent tool unless the user
  asks for it. The two rules point opposite ways and neither yields to the other, so
  I edited the spec inline and flagged it rather than spawning a delegate against an
  explicit prohibition. The change was four prose edits to an already-implemented
  spec's stories and coverage rows — no code, no seams.
- **Right behavior:** unclear, which is why this is here. The route I took keeps the
  harness instruction (which the reviewer set for this session) above a phase rule,
  and it matches the phase's own lighter-path posture for a few-line change — but
  the phase makes no exception for either, so I took an unsanctioned one.
- **Proposed rule change:** give "Route the venue" a stated precedence clause — when
  the harness cannot or may not spawn a write subagent, name the fallback rather than
  leaving the phase unsatisfiable. `craft-delegate`'s capability-aware policy covers
  *cannot*; it does not cover *may not*. A second candidate: exempt a spec-doc-only
  correction (no code, no seams) from the write-delegate requirement outright, since
  routing prose edits through a worktree-isolated delegate costs more than it catches.
