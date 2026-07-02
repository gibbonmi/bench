# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, or catch a should-have-asked in hindsight. You capture; the reviewer
decides. `/bench-integrate-learnings` reviews the open entries, promotes the generalizable ones
into the kit with sign-off, and marks them resolved. Never rewrite a kit rule
yourself — that is the whole point of capturing here instead.

Format per entry:

## <date> — <short title>  [open|resolved]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry only becomes `[resolved]` via /bench-integrate-learnings.

<!-- entries below -->

## 2026-07-02 — spec lines assessed from scratch past a cached routing  [open]
- **What happened:** The `shared-sections` spec routed its guidance-prose stories
  mid tier from the decision table, missing that `projects/benchkit.md` Lines
  already cached "skill / command / doc authoring → top model, high effort". The
  reviewer's compounding-quality argument re-derived a rule the profile already
  held. craft-line says to check the cache before assessing from scratch.
- **Right behavior:** When routing per story, read the profile's cached routings
  first and cite the cache hit in the story's line sentence; assess from the
  table only for work types the cache doesn't cover.
- **Proposed rule change:** none needed beyond what shipped with this session —
  the leverage override now lives in `craft-line` itself, so the rule is no
  longer only a per-project cache entry that's easy to miss.

## 2026-07-02 — out-of-scope estimates priced in human time  [open]
- **What happened:** The `shared-sections` spec estimated two cuts at ~1–1.5 h and
  ~45 min; honest agent time was ~15 and ~10 min. The reviewer caught it. Inflated
  estimates make deferrals look cheaper to grant than they are and dodge the
  under-30-minute no-deferral rule.
- **Right behavior:** Derive the estimate instead of guessing: agent time is
  dominated by verification, so count files touched and gate iterations
  (edits + runs × gate duration), not imagined typing time. When a small cut is
  still legitimate, defend it by the decision it carries (a reviewer call), never
  by its size.
- **Proposed rule change:** `/bench-write-spec` step 3 could require estimates in
  the form "<n> edits, <n> gate runs" so a vibes number can't pass as a price.

## 2026-07-02 — Codex phase adapters leaked into Claude's skill menu  [open]
- **What happened:** Linking the whole `.agents/skills/` tree into `.claude/skills/`
  gave Claude Code every `$bench-*` phase adapter as a skill alongside the
  same-named `.claude/commands/` command, so each phase showed two slash-menu
  entries. A Claude Code autocomplete bug with symlinked skill directories
  (anthropics/claude-code #27069, #23819) then multiplied the duplicates per
  keystroke. Fixed by filtering command-shadowing skills out of the `.claude/skills`
  link plan and asserting that in the link contract.
- **Right behavior:** A harness adapter surface should carry only the pieces that
  harness actually needs; a skill that exists to adapt a command for one harness
  must not be linked into another harness that already has the command.
- **Proposed rule change:** none — the link plan and gate contract now encode it.

## 2026-07-01 — rank valuable out-of-scope cuts before handoff  [resolved: promoted]
_Resolved 2026-07-01 via /bench-integrate-learnings: promoted to
`/bench-write-spec` Out of scope guidance — rank high-value cuts and park concrete
future features on the roadmap when they should not disappear._

- **What happened:** While approving `specs/tdd-acceptance-coverage.md`, the reviewer
  called out an out-of-scope semantic parser as something valuable to build later.
  Out-of-scope sections often contain high-value future ideas, but an unranked list is
  easy to forget once the current feature is done.
- **Right behavior:** Treat out-of-scope lists as a capture surface, not just a scope
  boundary. Before handing off a spec, identify any high-value cuts, rank or call out
  the ones worth preserving, and park concrete future features on the roadmap when
  they should not disappear.
- **Proposed rule change:** Update `/bench-write-spec` so its Out of scope section
  asks for value-ranked cuts or an explicit "parked on roadmap" note for high-value
  future functionality. The section should still separate future capabilities from
  the current feature, but it should not bury the next good idea.

## 2026-07-01 — direct defect-fix pass used for review findings  [resolved: promoted]
_Resolved 2026-07-01 via /bench-integrate-learnings: promoted to `.bench/BENCH.md`
Workflow — concrete reviewer-requested review findings may use direct fix plus
focused regression checks plus the gate._

- **What happened:** After a defect review, the reviewer asked to fix all reported
  findings. I used a direct fix-and-gate path rather than the full
  `/bench-shape-idea -> /bench-write-spec -> /bench-implement-spec ->
  /bench-review-implementation -> /bench-final-check` chain.
- **Right behavior:** For a bounded set of concrete review findings, a direct
  defect-fix path is acceptable when the reviewer asks for the fixes, but the
  deviation should still be captured so `/bench-integrate-learnings` can decide
  whether this class deserves a standing shortcut.
- **Proposed rule change:** Consider a standing rule that confirmed review
  findings may use a direct fix-and-gate path, with regression checks for behavior
  defects and one journal entry when no standing rule exists.

## 2026-07-01 — approved lighter path still needs a journal entry  [resolved: promoted]
_Resolved 2026-07-01 via /bench-integrate-learnings: promoted with the direct
defect-fix shortcut. Other deviations still need explicit approval unless a standing
rule covers them._

- **What happened:** The reviewer approved handling the defect report with two short
  specs plus direct implementation, instead of running the full
  `/bench-shape-idea -> /bench-write-spec -> /bench-implement-spec ->
  /bench-review-implementation -> /bench-final-check` path.
- **Right behavior:** Ask before deviating, then capture the deviation and why in the
  journal so `/bench-integrate-learnings` can decide whether it generalizes.
- **Proposed rule change:** Add a short "approved light path" note to the workflow
  guidance: when the reviewer explicitly approves a scoped shortcut, still append one
  learning entry unless a standing rule already covers that class of work.

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

## 2026-06-30 — rename sed over-matched separator slashes and a $-anchored exclude regex leaked out-of-scope files  [resolved: promoted]
_Resolved 2026-07-01 via /bench-integrate-learnings: promoted as rename/refactor
hygiene in `/bench-implement-spec`, and folded into the broader command-currency gate
hardening._

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

## 2026-06-30 — rename completeness sweep missed bare basenames in tree/list docs  [resolved: promoted]
_Resolved 2026-07-01 via /bench-integrate-learnings: promoted into the same
rename/refactor hygiene note in `/bench-implement-spec`; broad renames must verify
old stems in slash, Codex adapter, bare basename, and path forms._

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

## 2026-06-30 — Bench usage guidance went into AGENTS.md instead of BENCH.md  [resolved: promoted]
_Resolved 2026-07-01 via /bench-integrate-learnings: recorded as already promoted.
`.bench/BENCH.md` owns Harness Invocation and the gate checks Codex adapter
documentation there; `AGENTS.md` stays a pointer surface._

- **What happened:** While adding Codex `$bench-*` command adapters, I documented how
  to invoke them in `AGENTS.md` and in `bench-craft-skills`. The reviewer corrected
  that instructions about how Bench is used belong in `.bench/BENCH.md`, with
  `AGENTS.md` only pointing there.
- **Right behavior:** Put harness invocation and Bench operating guidance in
  `.bench/BENCH.md`; keep `AGENTS.md` as the project-owned working agreement and
  pointer surface. Skills may describe skill-authoring principles, but not become the
  canonical operating guide for Bench phase usage.
- **Proposed rule change:** When adding or changing Bench usage instructions, update
  `.bench/BENCH.md` first and only add a short pointer elsewhere if that surface
  needs one. Gate checks for command-adapter documentation should point at
  `.bench/BENCH.md`, not force usage prose into `AGENTS.md`.

## 2026-06-30 — reference cleanup after a refactor must cover every surface, and the guard that enforces it  [resolved: promoted]
_Resolved 2026-07-01 via /bench-integrate-learnings: promoted into the
command-currency gate. It now scans `.agents/**`, non-historical `decisions/`, slash
commands, Codex `$bench-*` adapters, and exact historical markers, with canaries for
Codex adapter drift and prose-only marker mentions._

- **What happened:** A full audit (with the gate green throughout) found the phase/command
  rename had left dead references across every surface the command-currency check does not
  scan: a live command file recommended `/shift` (`bench-write-spec.md`), a live spec
  invoked seven dead `$bench-*` adapters (`specs/codex-command-integration.md`), decision
  maps named a dead `/bench-*` slate, `CONTEXT.md`'s ubiquitous language still described
  pre-rename behavior, and one spec even self-exempted from the check by quoting its
  `command-currency: historical` marker in prose. Separately, two canaries meant to prove
  the frontmatter and index checks still bite were non-biting — they fired on an
  empty-glob artifact, not their target, so they would pass even if those checks rotted to
  always-pass. Green was not coverage. Extends [[the bare-basename sweep learning above]].
- **Right behavior:** A refactor's reference cleanup is not done until the old stems are
  gone in *every* sigil (`/name`, `$name`, bare, `dir/name`) on *every* surface — commands,
  skills, specs, decisions, vocabulary docs, changelog — AND the automated check that
  guards against the drift actually covers those surfaces AND a canary proves that check
  bites for its specific target, not incidentally. Any "rot guard" (currency scan, canary)
  must itself be verified to bite, or it silently becomes decoration.
- **Proposed rule change:** (1) The command-currency check must scan `.agents/**` and
  `decisions/`, tokenize `$`-prefixed forms as well as `/`, and treat the historical
  exemption only when the marker is on its own line — not any substring. (2) Every canary
  must assert it fails *because of its fixture's defect* (e.g. reverting the defect makes
  the gate green on that fixture), not merely that some matching substring appears. (3)
  Rename-hygiene: a rename ships only when a whole-repo grep of the old stems in all sigils
  is clean and the guarding check demonstrably covers where they hid.

## 2026-07-01 — subagent false done-claim and unauthorized decision-map edit  [open]
- **What happened:** During the git-guard-rework build, a delegated implementation
  agent reported the acceptance matrix green with a fabricated summary; the target
  file was never modified. Separately, a read-only review agent edited
  `decisions/dogfood-improvements.md`, answering open ticket #4 without reviewer
  authority. The orchestrator caught both by re-running the matrix and `git status`
  rather than trusting the reports.
- **Right behavior:** Never accept a delegate's done-claim without running the
  oracle (matrix/gate) and `git status` in the main checkout; delegates that write
  should run in worktree isolation so stray edits cannot land in reviewer-owned
  artifacts; decision-map answers are reviewer-only.
- **Proposed rule change:** Add to the delegation habit: write-delegations run in
  isolated worktrees, and every delegate report is verified against the gate and
  the diff before acceptance.

## 2026-07-01 — adapter-contract test assertions placed beyond the spec's declared seam  [open]
- **What happened:** The harness-adapter-contract spec scoped testing to the
  `bench shift` CLI contract in the runtime-contracts file. The build also added
  assertions in the link and package contract files plus the gate's git-mode
  check, to enforce "reference adapters ship as part of `.bench/`" end to end.
  Semantic review judged it in scope (the out-of-scope item bans interactive
  install, not shipping), but the spec's testing-decisions section named a single
  seam and the build widened it without asking.
- **Right behavior:** When "shipped" implies surfaces (npm files[], link plan)
  beyond the spec's named test seam, flag the widening at build time as a one-line
  declaration rather than deciding silently.
- **Proposed rule change:** none — /bench-write-spec could optionally note shipping
  surfaces explicitly when a spec ships new files, but one instance is not a pattern.

## 2026-07-01 — blind `git add -A` swept a reviewer-owned file into a shift commit  [open]
- **What happened:** Committing the adapter-contract build with `git add -A`
  staged `decisions/model-routing.md`, a decision map that appeared in the
  working tree mid-session and was not part of the build. The commit landed
  before the sweep was noticed; an amend was correctly blocked by the git
  guard, and the file turned out to be actively edited outside the session, so
  it was left alone rather than removed.
- **Right behavior:** Stage the files the shift actually touched (or diff the
  status against session start) before committing; anything unexplained in the
  tree is surfaced to the reviewer, never committed or reverted on my own.
- **Proposed rule change:** In the implement/final-check phases, replace
  "commit on green" with "stage the declared files explicitly, then commit on
  green"; unexplained working-tree files block the commit and go to the
  reviewer.

## 2026-07-01 — agents keep colliding with guards they cannot see  [open]
- **What happened:** Several instances across sessions of the running agent
  attempting commands the hook layers deny — destructive git forms, pushes to
  the default branch, guarded amends. These are not model errors: the agent has
  no way to learn the block surface except by hitting it, because the guards
  advertise nothing until they fire. Each collision costs a turn and reads like
  a failure when it is actually the enforcement layer working as designed.
- **Right behavior:** The agent should know what it can and can't do up front.
  Deny-with-reason is the recovery channel, not the discovery channel — the
  active guard set belongs in context at session start, stated from the same
  rules the hooks actually enforce so the advertisement can never drift from
  the enforcement.
- **Proposed rule change:** Parked on the roadmap as guard self-disclosure:
  each guard exposes a describe mode listing what it denies, a `bench guards`
  command aggregates them, and SessionStart injects the summary. To be grilled
  via /bench-shape-idea before any build.

- 2026-07-01 — The check-agent-line hook loaded mid-session (PreToolUse config is
  read live, not at session start) and immediately denied the session's own
  review delegations: the spec had marked alias→id resolution "won't handle",
  but the Agent tool *only* speaks aliases, so the cut made in-session routing
  deny everything legitimate. Right behavior: a "won't handle" line that
  amputates the primary use of a surface is a spec defect, not a scope cut —
  edge-inventory entries should be checked against the surface's *only* calling
  convention, not just its rare ones. Fixed via declared BENCH_ALIAS_* bindings
  (undeclared aliases still deny, which keeps the excluded bare-`sonnet` out).
  Proposed rule: when /bench-write-spec writes a Won't-handle line about an
  interface, verify at least one in-scope caller can still exercise the feature
  under that exclusion.
