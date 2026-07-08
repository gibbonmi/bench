# Skills assessment — mattpocock/skills v1.1 vs the Bench kit

**Status (2026-07-08):** adoptions 1–10 applied to the kit as a reviewer-approved
direct fix-and-gate batch; the consciously-not-adopted list stands as decided.

Source: `~/workspace/reference-skill-repos/skills` diff `v1.0.1..v1.1.0` (67 files,
+2088/−448), its CHANGELOG, and the release video transcript. Compared against the
kit's skills (`.agents/skills/`), commands (`.agents/commands/`), and CONTEXT.md.

**Verdict in one line:** bench is already ahead of v1.1 on review structure (three
axes vs two), spec discipline, and delegation — but v1.1 contains two direct bug
fixes bench inherited the bug for (self-grilling, missing enact gate), one
outrageously-cheap quality win (the Fowler smell baseline), and a handful of
smaller adoptions. Nothing in v1.1 contradicts a bench invariant; three upstream
moves should be consciously *rejected* as conflicting with closed bench decisions.

## Recommended adoptions, ranked

| # | Change | Bench target | Priority |
|---|--------|-------------|----------|
| 1 | Facts-vs-decisions split in grilling | `bench-craft-grill` | **High — bug fix** |
| 2 | Confirmation gate before enacting | `bench-craft-grill` | **High — bug fix** |
| 3 | Fowler smell baseline on the Standards axis | `bench-craft-review` (+ dedupe `bench-craft-tdd`) | **High** |
| 4 | Wide-refactor expand–contract sequencing | `bench-craft-spec` | Medium |
| 5 | Map upgrades: Destination + fog/out-of-scope sections | `bench-shape-idea` command | Medium |
| 6 | Negation + Negative Space failure modes | `bench-craft-skills` | Medium |
| 7 | `Task` ticket type (+ optional HITL/AFK tags) | `bench-shape-idea` command | Low-medium |
| 8 | Primary-sources rule on Research tickets | `bench-shape-idea` command | Low |
| 9 | Tautological-test recompute clause | `bench-craft-tdd` | Low |
| 10 | Strip stray "PRD" from kit prose | `bench-craft-grill:26`, `bench-shape-idea.md:144` | Tiny |

### 1+2. Grilling: facts vs decisions, and the enact gate — bug fixes

`bench-craft-grill/SKILL.md:11` still carries the **exact line v1.1 identifies as
the self-grilling bug**: "If a question can be answered by exploring the codebase,
explore the codebase instead." Upstream found that once grilling runs inside a
resolve-the-ticket frame (which is precisely how bench uses it — grill tickets in
`/bench-shape-idea`), that line reads as license to answer *decisions*
autonomously too, and the agent grills itself. The transcript notes this failure
"was especially happening with Fable" — the model this repo runs top-tier.

v1.1's replacement drops into bench's first-person-reviewer voice nearly verbatim:

> If a *fact* can be found by exploring the codebase, look it up rather than
> asking me. The *decisions*, though, are mine — put each one to me and wait for
> my answer.

Second fix, same skill: an explicit stop-gate — "Do not enact the plan until I
confirm we have reached a shared understanding." Upstream added it because
sessions were sliding from the last question straight into implementation. Bench
has workflow-level stops (spec sign-off is a hard stop), but `craft-grill` also
fires standalone ("grill me"), where no downstream gate exists. One line. It does
**not** conflict with the grill-continuously rule: it stops *enactment*, not
questioning — keep carrying the grill ticket to ticket.

Optional micro-add: upstream attached a one-clause *why* to one-question-at-a-time
("asking multiple questions at once is bewildering") because the bare rule kept
being violated. Bench's Discipline bullet states the rule without the why; the
weakest-reader test (`craft-skills`) argues for the clause.

### 3. Fowler smell baseline on the Standards axis

v1.1's biggest quality win per the author ("outrageously useful, really cheap"):
`code-review` now always carries ~12 curated smells from *Refactoring* ch.3
(Mysterious Name, Duplicated Code, Feature Envy, Data Clumps, Primitive
Obsession, Repeated Switches, Shotgun Surgery, Divergent Change, Speculative
Generality, Message Chains, Middle Man, Refused Bequest), each one line of
*what it is → how to fix*. The mechanism is leading words: the smells are deep in
the model's prior, so naming them is enough to recruit detection. Two binding
rules keep it safe: **a documented repo standard overrides the baseline**, and
**every smell is a judgement call, never a hard violation** (and skip anything
tooling already enforces).

Bench's `craft-review` Standards axis currently hunts documented conventions
only, so undocumented-but-classic rot is invisible to it. The baseline slots
cleanly into the Standards section of `bench-craft-review/SKILL.md` — and bench's
delegation model makes it cheaper than upstream's: review delegates are charged
with the skill *by path*, so the smells live in the one charged file instead of
being pasted into every sub-agent prompt. Citation standard already fits: a
baseline finding cites the named smell plus the quoted hunk, filed under
judgement calls, which the axis already separates.

**One-source-per-fact interaction:** `bench-craft-tdd`'s Refactor step
(SKILL.md:58–63) already names a partial smell list (duplication, feature envy,
primitive obsession, shallow modules). If the baseline lands in `craft-review`,
collapse the TDD list to a pointer — two smell lists is exactly the knowledge
duplication the code standard bans.

### 4. Wide refactors: expand–contract

v1.1's `to-tickets` learned that one class of work breaks vertical slicing: a
**wide refactor** — one mechanical change (rename a column, retype a shared
symbol) whose blast radius breaks thousands of call sites at once, so no slice
lands green. The fix is sequencing by **expand–contract**: expand the new form
beside the old; migrate call sites in batches sized by blast radius, green batch
to batch because the old form still exists; contract the old form away once no
caller remains. When even batches can't stay green alone, keep the sequence on a
shared branch and promise green only at a final integrate-and-verify step.

Bench has no counterpart (`rg 'expand.contract|wide refactor'` over the kit is
empty), and invariant 4 ("one small change at a time, repo stays green") plus
`/bench-implement-spec`'s vertical-slice rule fail on exactly this case — the kit
currently offers no legal route through a wide refactor. Proposed home:
`bench-craft-spec`'s "Story sizing and scope cuts" section (sizing is where
slicing is decided), one short paragraph; `/bench-implement-spec` needs no edit
since it already defers slicing rules upstream.

### 5. Decision-map upgrades from the wayfinder reframe

Upstream renamed and reframed `decision-mapping` → `wayfinder`. Bench's
`/bench-shape-idea` is the sibling of the old skill and already has several v1.1
features (no-fog early exit, create-then-wire ticket edges, one-bootstrap-session
rule). Three map-structure ideas are worth porting into the
`decisions/<topic>.md` format in `.agents/commands/bench-shape-idea.md`:

- **`## Destination`** — one or two lines naming what the map is finding its way
  to, written *first*, before any ticket; it fixes scope and every resume session
  orients to it. Bench's map today has only the topic slug; scope lives nowhere
  until the Handoff at close.
- **`## Not yet specified`** — a home for in-scope fog you can't yet phrase as a
  ticket. Bench states the map is "deliberately incomplete beyond the frontier"
  but gives the dim view no place on the page, so it lives in the reviewer's
  head between sessions. Test upstream ships: ticket when the question is sharp
  (even if blocked); fog when it isn't. Don't pre-slice fog into ticket-sized
  pieces.
- **`## Out of scope`** — work ruled beyond the destination, distinct from fog:
  it never graduates. Bench's Handoff has "Rejected alternatives" but only at
  close; mid-map there is no way to record "we looked at this and ruled it out"
  without leaving a ticket that reads as open work.

### 6. Skill-writing failure modes: Negation and Negative Space

v1.1 added two failure modes to `writing-great-skills` that
`bench-craft-skills`'s Failure modes section lacks:

- **Negation** — steering by prohibition drags the forbidden behaviour into
  context and makes it *more* available ("don't think of an elephant"). Cure:
  prompt the positive; a prohibition earns its place only as a hard guardrail on
  behaviour you can't phrase positively, and even then pair it with the positive
  target. Directly relevant to the kit: its prose leans hard on never/don't —
  mostly as legitimate hard guardrails, but the entry gives the pruning pass a
  test for which ones aren't.
- **Negative Space** — every decision a skill declines to make is silently
  delegated to the model's priors, not left neutral. Cure: read a draft for its
  silences and decide each omission deliberately — fill it, or leave it open as
  a real branch.

Both are a bullet each, matching the section's existing shape.

### 7–10. Small adoptions

- **`Task` ticket type** (bench-shape-idea): manual work that must happen before
  a decision can be made — provisioning access, signing up for a service, moving
  data so its shape can be seen. Bench's Research/Prototype/Grill have no slot
  for it; today it would be forced into a Grill ticket it doesn't fit. Upstream's
  framing: the one type that *does* rather than decides, earning its place by
  unblocking a decision. Optionally tag all types **HITL/AFK** (grilling and
  prototype only resolve through live exchange; research is agent-alone) — one
  line, and it reinforces fix #1 ("a grilling agent that answers its own
  questions has broken HITL").
- **Research tickets — primary sources** (bench-shape-idea): upstream's research
  skill requires primary sources ("follow every claim back to the source that
  owns it") and a per-claim citation in the output asset. Bench's Research
  ticket type has the byte/wire-probe rule but no source-quality or citation
  rule; `craft-delegate` spot-checks citations on the verify side but nothing
  requires them on the produce side. One added clause. A standalone research
  skill is *not* needed — bench composes this via `craft-delegate`.
- **Tautological tests** (bench-craft-tdd): bench already bans pasting the
  implementation's output back and has the vacuity check; v1.1 adds the third
  variant — the test *recomputes* the expected value with the same algorithm
  (`expect(calculateTotal(items)).toBe(items.reduce(...))`), passing by
  construction. One clause in the Red step closes it.
- **Vocabulary sweep**: CONTEXT.md bans "PRD" ("Not 'PRD'… — decision map"),
  and v1.1's to-spec rename validates bench's spec-first naming. Two stray
  "PRD"s remain in kit prose (`bench-craft-grill/SKILL.md:26`,
  `.agents/commands/bench-shape-idea.md:144`); make them "spec".

## Consciously not adopted (and why)

- **Tracker-backed collaborative map** (wayfinder's headline move). Bench's
  decision map is a compact git-tracked file loaded whole into every session —
  a decided design (CONTEXT.md's "decision map" entry, promoted-then-deleted
  lifecycle) optimized for a single reviewer, and the map-as-index, claim-by-
  assignment, native-blocking, and refer-by-name features all exist to serve
  the multi-session tracker mode bench doesn't run. Revisit only if the kit
  ever targets team-shared planning.
- **One ticket per session** (wayfinder's work-mode rule). Directly contradicts
  bench's grill-continuously rule — resume mode deliberately carries the grill
  ticket-to-ticket while the reviewer is present and answering. Closed decision;
  keep bench's behavior.
- **TDD reshape: refactor dropped from the loop.** v1.1 made TDD red→green only
  and moved refactoring wholly to code review ("don't overload the
  implementation"). Bench's cycle keeps a Refactor step, and bench's context
  argues for keeping it: inside a `bench shift`, refactoring under a green test
  in the same iteration is cheap and safe, while `/bench-review-implementation`
  is advisory and runs late. Recommendation: keep the step, adopt only the
  dedupe half (smell list moves to `craft-review`, TDD points at it). Flagged
  here as a reviewer call since it's the one place bench now deliberately
  diverges from upstream's loop shape.
- **Everything bench already has or exceeds**: two-axis review with parallel
  sub-agents (bench runs three axes, already parallel, already fail-fast on the
  diff, already skips Spec when no spec exists); the `implement` skill (a
  6-line flow note — `/bench-implement-spec` is a superset); to-spec/to-tickets
  renames (bench never said PRD or issues); seams-as-leading-word in TDD
  (bench's entire TDD skill is seam-anchored); the ask-matt router
  (`.bench/BENCH.md` + the skills index fill that role); triage-for-PRs and
  issue-tracker indirection (bench has no triage surface); teach, loop-me,
  wizard (out of kit scope).

## Housekeeping outside the kit

The machine-level installs at `~/.claude/skills/` and `~/.agents/skills/` are
v1.0 vintage: they still contain `to-prd`, `to-issues`, the pre-fix `grilling`
line inside `grill-me`/`grill-with-docs`, and `write-a-skill`/`caveman`/
`zoom-out`, which v1.0–1.1 renamed or deleted. The installer does not handle
renames, so the video's own migration advice applies: clear the stale ones and
re-run `npx skills add mattpocock/skills`, then sweep the folders for leftovers
(`to-prd`, `to-issues`, `diagnose`, `write-a-skill`, `caveman`, `zoom-out`).
Worth doing so the harness stops offering renamed-away skills alongside the kit.

## Suggested next action

Items 1–3 are small, high-value skill edits; 4–10 are a paragraph or less each.
The whole batch fits one `craft-synthesis`-disciplined pass: one branch, one
commit per item (or one commit for 7–10), each independently vetoable. On your
OK I'd run it as a direct fix-and-gate pass rather than the full spec pipeline —
these are kit skill edits with no runtime surface, and the gate's conformance
sweep plus `craft-skills` review cover them.
