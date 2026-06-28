# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, or catch a should-have-asked in hindsight. You capture; the reviewer
decides. `/resynthesize` reviews the open entries, promotes the generalizable ones
into the kit with sign-off, and marks them resolved. Never rewrite a kit rule
yourself — that is the whole point of capturing here instead.

Format per entry:

## <date> — <short title>  [open|resolved]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry only becomes `[resolved]` via /resynthesize.

<!-- entries below -->

## 2026-06-27 — gate was red at HEAD; a commit claimed green without running it  [resolved: dismissed]
_Resolved 2026-06-28 via /resynthesize: dismissed — already governed by invariant 1
and /verify-gate's "never substitute the model's judgment for the gate." A pre-commit
gate run would be a fourth check surface (HANDOFF says prune toward, not add), and the
skipped-reference regression is now caught mechanically by the command↔index
conformance check plus gate checks 1c/1d._

- **What happened:** Fixing the missing `.bench/learnings.md` scaffold, I ran the
  full gate and found it red at clean HEAD — commit `cea2f42` renamed the commands
  and its message asserted "Gate green: the index<->disk conformance check confirms
  no reference was missed," but AGENTS.md still carried the old names
  (`/map`, `/diagnose`, `/review`, `/verify`). The one file the conformance check
  reads was the one the rename missed, which means the gate was never actually run
  after that commit. Same failure class as `d77063c` (a rename that skipped a
  reference).
- **Right behavior:** Never write "gate green" in a commit message from belief; run
  `bench gate` and paste/observe its exit. The oracle is the gate, not the diff.
- **Proposed rule change:** Add a one-line reminder to the build/verify-gate path (or
  a pre-commit nudge) that "gate green" in a message must come from an actual run,
  not inspection. Possibly a Stop-hook check that the gate was run since the last
  edit.

## 2026-06-27 — scaffolded files must be created by init, not just referenced  [resolved: promoted]
_Resolved 2026-06-28 via /resynthesize: promoted to HANDOFF "Discipline carried over"
as a one-line maintainer rule. The executable fix (init scaffold + gate check 1d)
already shipped in 724bf8c. The proposed generalization of check 1d to every `.bench/*`
file was skipped as speculative — only two exist._

- **What happened:** `bench init` scaffolded `.bench/gate.sh` but not
  `.bench/learnings.md`, while AGENTS.md (write side) and `/resynthesize` (read side)
  both depend on that file existing. The self-learning contract pointed at a path
  nothing created.
- **Right behavior:** Any file the kit's prose instructs an agent to read or append
  to must be produced by `init` (guarded, idempotent) and locked by a gate check that
  exercises the real init path.
- **Proposed rule change:** When adding a contract that names a `.bench/*` file,
  the same change must (1) scaffold it in `init()` and (2) add a behavioral gate check.
  Consider generalizing gate check 1d to assert every kit-referenced `.bench/*` file
  is scaffolded.

## 2026-06-28 — questions must carry a recommendation  [open]
- **What happened:** During /start-ideation I asked two AskUserQuestion forks with
  neutral options and no recommended pick. The user corrected me: always put forth a
  recommendation when asking.
- **Right behavior:** Lead every question with the option I'd choose and why — put the
  recommended option first with "(Recommended)" per the AskUserQuestion convention; in
  prose, state the pick and a one-clause reason.
- **Proposed rule change:** Add to the grill skill and AGENTS.md "How to talk to me":
  a question without a recommendation is incomplete. Surface judgment, don't offer a
  blind menu.

## 2026-06-28 — always recommend the proper next action  [open]
- **What happened:** After finishing a step (e.g. writing the canary spec), I offered
  next actions as a neutral menu — "/build or bench shift?" — and the user had to ask
  for a recommendation. Same pattern as the questions learning, applied to hand-off
  between workflow phases.
- **Right behavior:** When handing back at a phase boundary, recommend the proper next
  action, picked from the implementation type and the goal — e.g. /build interactively
  for edits to the oracle or where a design call is still vetoable; bench shift for
  locked-spec mechanical work. State the pick and the one-clause reason; the menu is
  context, not the answer.
- **Proposed rule change:** Generalize the questions rule to all hand-offs — every
  command's exit ("offer to run /build or /shift") should lead with a recommended next
  action keyed to type+goal, not a neutral either/or. Fold into AGENTS.md "How to talk
  to me" alongside the questions rule.
