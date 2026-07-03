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

## 2026-07-03 — session-start stale gate confuses without a benign/real split  [open]
- **What happened:** Reviewer flagged that the gate "is almost always stale in a
  new session," which reads as alarming. Diagnosis: the verdict is content-addressed,
  so it goes stale the instant the tree moves past the last green — and sessions
  routinely end with a change after the last gate run (a manual commit that wasn't
  re-gated, or a `bench idea` park that dirties ROADMAP.md). So new sessions almost
  always open stale, but the drift is often benign (capture-scratch like ROADMAP.md /
  .bench-notes.md, which no gate check reads) rather than unverified code.
- **Right behavior:** At session start, when the gate reads stale, tell the user
  *why* — split "benign drift only (e.g. a parked idea) → just a reminder to re-run
  the gate" from "committed code moved since the last green → real, re-run before you
  trust it." The bare word "stale → re-run the gate" hides that distinction and reads
  as an error even when it's harmless.
- **Proposed rule change:** Consider having `bench status` classify a stale verdict:
  if the diff from the gated tree is confined to capture-scratch paths, word it as a
  benign reminder; otherwise flag it as real. (Distinct from, and a lighter-weight
  alternative to, the parked idea of carving capture-scratch out of `gate_tree_hash`
  entirely — that changes the oracle's key and is sensitive on the tripwire branch.)

## 2026-07-03 — dogfooding gotcha preserved from a retired map: mid-session edits hit the next session  [open]
- **What happened:** Story 11 of spec-handoff-lifecycle retired `decisions/dogfooding.md`.
  Its self-host+canary-guard decision was already realized (the canary gate layer
  shipped; craft-synthesis carries the dogfood loop), so the promotion read ended in
  deletion. But one operating gotcha in that map was not recorded anywhere: the kit
  loads skills/commands at session start, so a mid-session edit to a skill or command
  lands in the *next* session, not the current one. The egg bites in three places —
  the gate (break it, lose the oracle), skills/commands (break a trigger, next session
  misfires), and the CLI (break `bench shift`, lose the loop).
- **Right behavior:** Per the map-Handoff rule (item 9), operating lessons go through
  `.bench/learnings.md` and `/bench-integrate-learnings`, not per-spec notes or a
  unilateral skill edit — so this gotcha is captured here rather than promoted into a
  skill on the worker's own authority.
- **Proposed rule change:** Consider a one-line note in `bench-craft-synthesis`'s
  dogfood loop: when a candidate change touches a skill or command trigger, the
  dogfood shift must be a *fresh* session, because the edit does not take effect in
  the session that made it. (Reviewer's call — this is the deferred/contestable
  promotion flagged from story 11.)

## 2026-07-03 — review-implementation: verify a returned finding in the live session, not a worktree  [open]
- **What happened:** During `/bench-review-implementation` on spec-sourcing, the
  Coverage axis returned a real finding (the negative anchor's line-oriented grep
  misses a hard-wrapped bypass reintroduction). I reproduced and fixed it directly
  in the live session's working tree — one repro, one fix, one gate run — rather than
  spinning up a separate worktree for the fix.
- **Right behavior:** When a review delegate returns with an issue, the invoking
  session tests and fixes the failure in the live session, not in a separate
  worktree. The worktree-isolation rule is for *write-delegations* (a subagent that
  edits files, run isolated so stray edits can't land in reviewer-owned files); the
  review axes themselves are read-only, so their findings come back to the one live
  session that owns the diff, which verifies and fixes them in place against the gate.
- **Proposed rule change:** Consider a one-line note in `/bench-review-implementation`
  (or `craft-delegate`): review delegates are read-only and return findings; the
  live session verifies and fixes a returned finding in place — worktree isolation
  is for write-delegations, not for reproducing a review finding.

## 2026-07-03 — implement-spec: the facilitating session must delegate every story, never build inline  [open]
- **What happened:** During `/bench-implement-spec` on shim-autoinstall, the live
  (facilitating) session implemented the stories itself inline — doctor, wrapper,
  postinstall, session-start, README, and the gate fragment — rather than delegating
  each story's build to a sub-agent. It declared the line honestly (mid/medium, one
  notch above the spec's cheap estimate) but did the work in the session's own context.
- **Right behavior:** The session that facilitates a build should never be the one
  that works a story. It must always delegate the story's build to a sub-agent —
  even when the sub-agent runs the exact same model and effort — then continue
  facilitating. The point is role separation and keeping the facilitator's context
  clean for orchestration and verification, not a model/effort saving; the tier
  match is irrelevant to the rule. The facilitator declares the line, spawns the
  delegate (in an isolated worktree for write-work per craft-delegate), verifies the
  returned claim against the gate and git status, then moves to the next story.
- **Proposed rule change:** Consider making delegation mandatory in
  `/bench-implement-spec` and `craft-delegate`: the facilitating session delegates
  every story build to a sub-agent regardless of tier match, and reserves itself for
  declaring the line, sequencing slices, and verifying delegate claims. (Reviewer's
  call — this generalizes the existing "run write-delegations in isolated worktrees"
  invariant into a "always delegate, never self-build" rule.)
