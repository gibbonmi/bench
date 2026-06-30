# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, or catch a should-have-asked in hindsight. You capture; the reviewer
decides. `/bench-learn` reviews the open entries, promotes the generalizable ones
into the kit with sign-off, and marks them resolved. Never rewrite a kit rule
yourself — that is the whole point of capturing here instead.

Format per entry:

## <date> — <short title>  [open|resolved]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry only becomes `[resolved]` via /bench-learn.

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

## 2026-06-28 — questions must carry a recommendation  [resolved: promoted]
_Resolved 2026-06-30 via /resynthesize: promoted to AGENTS.md "How to talk to me"
(mirrored in .bench/BENCH.md) as the "recommend, don't offer a blind menu" rule,
covering both questions and hand-offs._

- **What happened:** During /start-ideation I asked two AskUserQuestion forks with
  neutral options and no recommended pick. The user corrected me: always put forth a
  recommendation when asking.
- **Right behavior:** Lead every question with the option I'd choose and why — put the
  recommended option first with "(Recommended)" per the AskUserQuestion convention; in
  prose, state the pick and a one-clause reason.
- **Proposed rule change:** Add to the grill skill and AGENTS.md "How to talk to me":
  a question without a recommendation is incomplete. Surface judgment, don't offer a
  blind menu.

## 2026-06-28 — always recommend the proper next action  [resolved: promoted]
_Resolved 2026-06-30 via /resynthesize: promoted together with the questions learning
into the single "recommend, don't offer a blind menu" rule (AGENTS.md / BENCH.md) — it
applies to every question and every phase hand-off._

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

## 2026-06-28 — show spec/issue breakdowns as a table for approval before building  [resolved: promoted]
_Resolved 2026-06-30 via /resynthesize: promoted to the /spec exit — emit a
stories/seams/out-of-scope approval table and pause for sign-off before any build._

- **What happened:** I presented the canary spec as prose and moved to /build on a
  verbal OK. The user wants the breakdown shown in table form for explicit approval
  *before* construction continues.
- **Right behavior:** Before building from a spec (or any issue/work breakdown), render
  it as a scannable table — e.g. user stories / seams / out-of-scope, or the issue list
  with its slices — and get sign-off. Only then continue the workflow. The table is the
  approval gate, not a prose aside.
- **Proposed rule change:** /spec (and /to-issues) should emit a table summary at exit
  and pause for approval before /build or /shift proceeds. The full spec file stays as
  written; the table is the at-a-glance veto surface. Fold the "table for approval"
  expectation into /spec and AGENTS.md "How to talk to me".

## 2026-06-28 — clean summary + persistent progress visual at every stage boundary  [resolved: split]
_Resolved 2026-06-30 via /resynthesize: split. The stage-summary half is promoted —
covered by the "pick it up cold (what / where / next)" bullet in AGENTS.md / BENCH.md.
The persistent task-list tracker is dismissed: TaskCreate/TaskUpdate are
Claude-Code-only, so mandating them in harness-shared rules leaks one harness into the
core; it also went unused all session._

- **What happened:** Moving between workflow phases I dumped walls of text without a
  crisp "what just happened / what's next." The user wants every stage transition to
  carry a short summary and, ideally, a persistent visual of where we are in the
  pipeline.
- **Right behavior:** At each phase boundary (map→spec→build→prep-shift→verify-gate,
  etc.) lead with a 2–3 line summary: what finished, the result, what's next. Maintain
  a persistent progress visual — a Claude Code task list (TaskCreate/TaskUpdate) with
  one task per phase, updated as phases complete — so the pipeline state is always
  visible, not buried in prose.
- **Proposed rule change:** Every command exit emits a compact stage summary, and the
  workflow keeps a task-list progress tracker across phases. Fold into AGENTS.md "How
  to talk to me" and each command's exit section. Keep prose minimal; the table/list
  carries the state.

## 2026-06-30 — declared a decision map "resolved" while two tickets were unwritten  [resolved: promoted]
_Resolved 2026-06-30 via /resynthesize: promoted to the /grill and /start-ideation
exits — scan for unfilled "— (open|deferred)" / "GRILL DEFERRED" markers and refuse to
close a map while any remain._

- **What happened:** After grilling decisions/ambient-feedback.md I set the header to
  "GRILLED & RESOLVED" and moved to /spec, but tickets #3 and #5 still held "— (deferred)"
  placeholder answers — the decisions were made in conversation but never written into the
  map. Caught only when /spec re-read the map for feature A. (Ironically, the very
  `bench status` decision-map signal built this session keys on those exact placeholder
  markers and would have flagged it.)
- **Right behavior:** Before marking a map resolved or handing it to /spec, verify every
  ticket's Answer is actually written — scan for the open/deferred placeholder markers and
  confirm none remain. A decision agreed in chat but absent from the artifact is not
  recorded; the artifact is the source of truth, not the conversation.
- **Proposed rule change:** /grill and /start-ideation should, at the exit that declares a
  map closed, scan for unfilled Answer placeholders (`^— \((open|deferred)` / `GRILL
  DEFERRED`) and refuse to close while any remain. One-line guard in each command's exit.

## 2026-06-30 — rename sed over-matched separator slashes and a $-anchored exclude regex leaked out-of-scope files  [open]
- **What happened:** Building the mechanical command/skill rename, I drove the
  reference updates with `sed 's#/spec\b#/bench-spec#g'` (and `/build`, `/setup`,
  `/grill`) across the doc surface. Two bugs: (1) the exclude filter
  `grep -vE '^(decisions/|…)$'` was `$`-anchored, so `decisions/` only excluded a file
  literally named "decisions/" — every `decisions/*.md` (out of scope, historical) got
  edited anyway. (2) `/spec\b` matched **separator slashes** in prose — `map/spec/issue`,
  `Install/setup`, `agent/setup` — corrupting concept words the spec explicitly
  protected. The gate stayed green throughout (it doesn't read prose), so only the
  semantic `/bench-review` pass plus a "word-char before the new identifier" scan caught it.
- **Right behavior:** For a rename sed, (a) build the exclude list with *prefix* matching
  and dry-run it against the real file list before editing; (b) anchor an
  invocation pattern so a leading separator can't read as a command — match `/cmd`
  only when preceded by whitespace/start/punctuation, never bare, since `word/cmd` is a
  path or separator, not an invocation; (c) after editing, scan for any replacement
  preceded by a word-char and eyeball each — the gate will not catch prose corruption.
- **Proposed rule change:** Give the rename path (or a future `bench rename` helper) two
  guards: dry-run the scope filter and show in/out before touching files, and flag every
  replacement where the matched identifier is preceded by a word character for human
  review. Fold a one-line "rename hygiene" note into the build guidance.

## 2026-06-30 — rename completeness sweep missed bare basenames in tree/list docs  [open]
- **What happened:** The merged command/skill rename slice swept for `/cmd`
  invocations and `<dir>/<name>` path forms, but README's file-tree listing names
  commands and skills as *bare basenames* (`setup.md`, `seams/`) with the
  `commands/`/`skills/` prefix on a separate line. Those evaded the sweep and shipped
  stale; the resynthesize-split slice found the entire tree block still pre-rename and
  fixed it there.
- **Right behavior:** A rename completeness sweep must also match bare basenames in
  enumerated contexts (file trees, bulleted inventories), not only the slash- and
  path-anchored forms. Grep the old stems as whole words and eyeball, rather than
  trusting the anchored patterns to find everything.
- **Proposed rule change:** Fold "also sweep bare basenames in trees/lists" into the
  same rename-hygiene note as the separator-slash learning above — they're one rule:
  a rename isn't done until the stems are gone in *every* form, anchored or not.
