# implement-spec-full-run

Status: staged

Map: `decisions/implement-spec-full-run.md` (closed, commit 9a5d2e4). Slicing
confirmed with the reviewer 2026-07-28: the map's slices A and B land as one
spec in A-then-B story order, because both are the same ownership fence —
platform and command prose plus the anchors and canary fixtures that observe it.

Stories 1 through 10 carry that A-then-B order. Stories 11 through 15 were added
after approval, in the same ownership fence, and append rather than re-sort — so
story 15 is slice-A material sitting at the end. Build order follows the
dependencies stated in each story, not the numbering: story 4 writes the bounded
`--full` section that stories 5, 6, 7, 8, and 11 through 14 anchor facts inside,
and story 2 registers the markers whose prose stories 1 and 15 write.

Flagged for reviewer veto — decided here, not carried by the map:

- Routing the `--full` review through one fresh-context delegate charged with
  the standalone `/bench-review-implementation` contract, rather than fanning
  its three axes from the orchestrator.
- Discharging the map's Handoff item 7 uncertainty flag with a verified Codex
  invocation instead of the stub the Handoff allowed.
- Resolving the map's open "not yet specified" handoff question this way: the
  orchestrator writes the State section naming the phase reached, then refreshes
  the pin block with `bench handoff --next <command>`. No CLI change.
- Deviating from the map's three *fixed* escalation options by dropping any
  route this harness cannot actually invoke, and by specifying a plain numbered
  list as the fallback where a harness has no structured-prompt surface.
- Carrying the fix-don't-park rule and the Roles assume/verify sentence as
  entries in the existing shared-rule marker list rather than as new standalone
  anchors.
- Scoping the new anchors to their owning markdown section (the existing
  `markdownH2Section` idiom) rather than asserting whole-file presence, and
  pairing each contradiction-prone anchor with a `forbid`.
- Adding one clause to the command file's frontmatter description so `--full`
  is discoverable when the phase fires.
- Four canary fixtures rather than one per new anchor, following the family's
  established fixture-to-assertion ratio.
- Keeping `--full` out of `.bench/BENCH.md`'s Workflow list so the command file
  stays its one source.
- The review-delegate-unavailable stop in story 5, and story 10's mid-tier line
  for a story the profile's leverage override would otherwise route top.
- Fencing the run's diff to the spec's stories as an anchored counterweight to
  fix-don't-park, which the map never opened.
- Requiring the exit report to enumerate the coverage map row by row and to cite
  the record behind each phase claim, rather than leaving both to story 10's
  one-time observation.
- Adding stories 11 through 15's anchors without new canary fixtures, on the
  family's established density rule.
- Hanging a cross-harness falsification pass on story 8's size trigger as a
  separate question rather than a fourth route in that story's fixed menu.
- Routing story 15's warrant rule into this spec directly rather than through a
  `/bench-what-next` drain, on the reviewer's 2026-07-28 instruction. The journal
  entry that sourced it stays open in `.bench/learnings.md` for that drain to
  verdict as shipped.

## Problem

A spec's build today crosses four phase boundaries by hand. The session that
finishes `/bench-implement-spec` has to be told to run
`/bench-review-implementation`, then `/bench-final-check`, and each hand-off is
a place a run stalls or a phase gets skipped. Two habits make it worse: a small
defect discovered mid-build gets parked to `IDEAS.md` instead of fixed, so the
tree keeps known-bad code while the backlog grows; and the "NEVER assume,
always verify" rule sits once in `.bench/BENCH.md` Roles, far from the three
moments where sessions actually violate it — asking about a fact the tree
already answers, implementing against remembered instead of read code, and
reviewing from recall instead of citation. A run that reports on itself compounds
all of it: a diff that widened past the spec, a coverage map left
part-implemented, and a phase that never ran each produce an exit report that
reads exactly like a clean one.

## Solution

An opt-in `--full` mode on `/bench-implement-spec` that carries a spec from
build to push-ready — implement inline, review in a fresh-context delegate with a
cross-harness falsification pass offered when the diff is large, final-check
inline, debug on demand — with `session-handoff.md` rewritten at
every phase boundary so a cold session can resume mid-run. Alongside it, two
shared-rule changes the mode leans on: fix-don't-park, which sends a discovered
defect into the active workflow unless it genuinely needs a reviewer decision;
and point-of-use reinforcement of the assume/verify rule at the three phases
where the failure was observed. A third shared rule rides the same marker list
without belonging to the mode at all: in the communication rules, a claim resting
on a source outside the tree names what was read and what was not. The run's own
diff is fenced to the spec's
stories, and its exit report enumerates the coverage map row by row and cites the
record behind every phase claim — so a run that widened, stopped short, or skipped
a phase is visible without re-deriving it. Every prose change gets an anchor scoped to the
section that owns it, so it cannot drift out, be commented out, or be
contradicted elsewhere in the same file.

## User stories

1. As a session that discovers a small defect mid-work, I read one rule in
   `.bench/BENCH.md`'s Workflow section telling me the fix lands in the active
   workflow as its own commit, and that parking to `IDEAS.md` or
   `.bench/learnings.md` is reserved for a fix needing a reviewer decision, a
   new seam, or spec-level design — so the boundary is a decision test, not a
   size guess. The rule is anchored to the Workflow section specifically, not
   merely present somewhere in the file. It lives in the shared platform file
   only, so linked repos inherit it and `AGENTS.md` never restates it. Line:
   gpt-5.6-sol / high. This is shared platform prose that steers every future
   session's park-versus-fix call, which is the leverage override's cached
   routing in the profile.

2. As the conformance gate, I hold three rules in
   `checkSharedRuleSingleSource`'s marker list — the fix-don't-park rule, the
   Roles sentence "NEVER assume, always verify", and story 15's warrant rule — so
   each one reds when `.bench/BENCH.md` loses it and reds again when `AGENTS.md`
   or `README.md` restates it. This story owns every registration in that list;
   stories 1 and 15 own their rules' prose. One list entry per rule carries both
   directions; no rule gets a second, independently-worded assertion. Line: gpt-5.6-terra / medium.
   Registering a marker is mechanical at a known seam, but it is oracle code —
   a marker that names a string the canonical file does not contain would red
   the gate on a correct tree.

3. As a session entering a grill, an implementation, or a review, I meet a
   one-clause verify hook at the point of use: `/bench-shape-idea` tells me to
   look a fact up in the tree before asking about it, `/bench-implement-spec`
   tells me to verify a claim against the tree rather than memory — and a claim
   over a whole set by enumerating the set, not by extending one measured
   member — and `/bench-review-implementation` tells me a finding cites what I
   read now, not what I recall, and that a universal claim cites its
   enumeration or names itself a sample. The Roles sentence stays where it is;
   each hook points at the moment, not at a restated rule, and the quantifier
   clause's source is `craft-review`'s citation standard. Line: gpt-5.6-sol / high. Three phase-command
   edits are guidance prose whose defect is invisible to the gate and multiplies
   through every session that loads them.

4. As a reviewer invoking the build phase, `--full` is documented as one bounded
   markdown section in `.agents/commands/bench-implement-spec.md`: it is an
   opt-in flag, plain invocation keeps today's implement-only semantics, and
   `/bench-review-implementation`, `/bench-final-check`, and `/bench-debug` stay
   standalone commands for strict phased use and mid-run resumption. `--full`
   invoked with no spec argument, or one naming a path that does not exist,
   refuses at the entry contract and says which it was rather than inferring a
   spec. Every anchor for this feature is scoped to that section, and the
   opt-in and no-inference facts each carry a paired prohibition so a
   contradicting default-on or spec-inferring sentence cannot sit beside them.
   Line: gpt-5.6-sol / high. The invocation contract is the surface every future
   run reads first, a silently-inferred spec would let a full run build the
   wrong target, and a section-scoped anchor is what makes the rest of the
   feature's anchors mean placement rather than mere presence.

5. As a `--full` run, I execute implement inline, then spawn one fresh-context
   delegate charged with the standalone `/bench-review-implementation` contract
   and given the spec and the diff and nothing else, then run final-check
   inline, invoking `/bench-debug` inline when an issue needs deep analysis. The
   context that produced the code carries the assumptions that produced its
   bugs, so inline self-review is closed rather than deprioritized, and the
   input restriction is part of the anchored fact — a context-inheriting
   delegate is the same failure wearing a delegate's name. I verify the
   delegate's done-claim against the gate and `git status` per invariant 1
   before acting on it. When the delegate cannot run or returns nothing, the run
   stops and reports at that boundary rather than proceeding to final-check with
   review unrun. Line: gpt-5.6-sol / high. The fresh-context requirement is the
   quality claim the whole mode rests on, and the phrasing has to close both the
   inline shortcut and the inherited-context shortcut rather than merely prefer
   against them.

6. As a `--full` run holding review findings, I fix concrete defects — bugs,
   spec misses, missing coverage — and re-gate without stopping, and I flag
   contestable design and judgment findings in the exit report for reviewer
   veto rather than applying them. Neither half may be contradicted elsewhere in
   the section: a sentence that stops the run on a concrete defect, or one that
   applies a judgment finding, is prohibited alongside the rule. The repair pass
   is bounded by `/bench-review-implementation`'s existing terminal-repair rule
   and routed through `craft-delegate`'s existing repair allowance; this mode
   adds no second version of either. Line: gpt-5.6-sol / high. The disposition
   rule is what keeps an unattended run from either stalling on every finding or
   silently rewriting the reviewer's design, and one half without the other is a
   cheap and wrong reading.

7. As a session picking up an interrupted `--full` run cold, `session-handoff.md`
   tells me which phase the run reached, because at every phase boundary the run
   wrote that into the State section and then refreshed the pin block with
   `bench handoff --next <command>`, which derives the repository, branch, HEAD,
   spec, and gate facts and preserves State. The phase reached is the one fact
   this mode adds; everything else the handoff carries is already the existing
   handoff contract's, and this section points at it rather than re-enumerating
   it. A re-invoked `--full` resumes from the phase the handoff names instead of
   re-implementing from the top, and where the handoff and the tree disagree the
   tree wins. Line: gpt-5.6-sol / high. Cold resumption is the reason the mode
   persists anything at all, and the resume-don't-restart clause is what stops a
   re-invocation redoing landed work.

8. As the orchestrator judging a diff large enough that the mid binding could
   miss important bugs, I pause and ask the reviewer as a structured decision
   list with a recommendation for this run, offering three routes: continue at
   the mid binding, escalate to the top binding in this harness, or escalate to
   the top binding via the Codex CLI. The three routes are the fixed menu; a
   route is omitted only when this harness cannot invoke it at all, and the
   omission is stated rather than silent. I never escalate without asking.
   Harnesses without a structured-prompt surface ask the same question as a
   plain numbered list. Line: gpt-5.6-sol / high. This is where invariant 2's
   no-silent-escalation rule meets an autonomous loop, and the prose must name
   tiers by binding rather than by model token so `.bench/lines.env` stays the
   one source.

9. As the canary sweep, I prove four distinct anchor classes bite, one fixture
   each under `tests/canary/workflow-guidance-anchors/`: a verify hook removed
   from a phase command, the fresh-context review requirement removed from the
   `--full` section, the phase-boundary handoff persistence removed from it, and
   the ask-before-escalating rule removed from it. Each fixture names the
   diagnostic it expects, and no two fixtures exercise the same class. Fixture
   base names stay globally unique across every family. Line: gpt-5.6-terra /
   medium. Fixture authoring is mechanical against the family's existing
   examples, but a fixture whose EXPECT does not match its emitter reds as
   did-not-bite and blocks the gate.

10. As the reviewer, the first real `/bench-implement-spec --full` run is
    observed independently of its own exit report: the session's delegate
    invocation record shows a review delegate was spawned and what it was given,
    `git log` and the committed `session-handoff.md` show a boundary rewrite
    naming each phase, and the escalation decision is either a recorded question
    to me or an explicit statement that the trigger did not fire. The exit
    report is compared against those records, not accepted as them. Line:
    gpt-5.6-terra / medium. This is the only observation of the gate-invisible
    half of the feature, so it needs enough judgment to notice a report that
    describes a run that did not happen — which is why it does not take the
    cheap tier despite being read-and-compare work.

11. As a `--full` run, I implement the spec's stories and nothing else: work I
    notice outside them — an adjacent refactor, an unrelated improvement, a story
    the spec chose not to take — is recorded for the reviewer rather than built.
    The fence sits in the `--full` section beside the fix-don't-park route it
    counterweights, and a sentence licensing opportunistic improvement while the
    file is open is prohibited alongside it. The park-versus-fix test itself stays
    story 1's; this fact is only about the run's diff staying inside the spec's
    stories and seams. Line: gpt-5.6-sol / high. Fix-don't-park ratchets a run's
    diff in one direction only, and invariant 4's smallest-diff rule sits once in
    `.bench/BENCH.md` — the same distance-from-the-moment failure story 3 exists
    to fix.

12. As the reviewer reading a `--full` exit report, I get one disposition per row
    of the spec's acceptance coverage map — implemented, deferred, or won't-handle
    — named row by row against `bench coverage <spec>`'s enumeration, rather than
    a summary asserting the spec is done. A phrasing that reports completion in
    aggregate is prohibited alongside the rule. When the spec carries no coverage
    map, the report says so and accounts for the user stories instead. Line:
    gpt-5.6-sol / high. Stopping early is invisible in an aggregate claim and
    obvious in an enumeration, and `bench coverage` already emits that enumeration,
    so the rule adds a duty rather than a second source for what the rows are.

13. As the reviewer, every phase claim in a `--full` exit report cites the record
    that proves it — the review delegate's invocation, the commit shas the phase
    landed, the `session-handoff.md` boundary rewrite — rather than asserting the
    phase ran. A phrasing that reports a phase complete without its record is
    prohibited alongside the rule. Line: gpt-5.6-sol / high. Story 10 buys this
    observation once, by hand; making the report cite its own records is what makes
    every run after the first cheap to check, and it applies `craft-review`'s
    existing citation standard to the run's claims about itself.

14. As a `--full` run whose diff trips story 8's size condition, I ask a second
    question at that same trigger: whether to add a cross-harness falsification
    pass over the diff before final-check — the Codex CLI at the top binding,
    charged to refute the claim that the spec was implemented rather than to grade
    it against the three axes. It is its own question, not a fourth route in story
    8's fixed menu, because it changes the review rather than the build line, and
    it never runs standing — absent the trigger it is not offered. When this
    harness cannot invoke Codex the pass is omitted and the omission is stated,
    the same posture story 8 takes. Line: gpt-5.6-sol / high. The three review
    axes are already fresh contexts, so a second same-family review buys little;
    what this adds is a different model and a refutation charge — the pair that
    returned block on FT91's draft after this repo's own axes had cleared it.

15. As any session making a claim that rests on a source outside the tree — a
    reference repo, a vendored kit, an upstream doc — I read one rule in
    `.bench/BENCH.md`'s communication section telling me to name what I read and
    what I did not, so a claim's warrant travels with the claim. Story 2 registers
    its marker alongside the other two; this story owns the prose and a
    section-scoped assertion, so the anchor means placement in the communication
    rules rather than presence anywhere in the file.
    A phrasing asking for thoroughness rather than for disclosure is prohibited
    alongside it: thoroughness is unfalsifiable and a session always believes it
    read enough, while what went unread is checkable by the reviewer at a glance.
    Line: gpt-5.6-sol / high. This is always-loaded prose governing every claim in
    every session, which is the leverage override's cached routing in the profile.

## Implementation decisions

- The whole change is prose plus its observers: `.bench/BENCH.md`, three
  `.agents/commands/*.md` files — `bench-shape-idea.md`,
  `bench-implement-spec.md`, and `bench-review-implementation.md` — with their
  byte-identical `.claude/commands/` mirrors, assertions in `internal/conformance/docs_workflow_helpers_test.go`
  and `internal/conformance/validity_checks_test.go`, and four canary fixtures.
  No `bench` subcommand, no new flag parsing, no Go production code — `--full`
  is read by the session from its own invocation, which is why the map called
  the orchestration thin coordination over existing phases.
- Anchors for this feature are section-scoped, using the `markdownH2Section`
  helper `checkStructuredPhaseContract` already establishes, so an anchor means
  the fact is in the section that owns it. Whole-file `requireCollapsed` does
  not strip HTML comments, so a commented-out or relocated sentence would
  otherwise satisfy it. Where a fact has a cheap contradicting reading — opt-in
  versus default-on, no-inference versus inference, fix-and-continue versus
  stop, flag versus apply — the anchor is a require/`forbid` pair, matching the
  "one owner per workflow agreement" idiom already in the file.
- `--full` lives in `.agents/commands/bench-implement-spec.md` only.
  `.bench/BENCH.md`'s Workflow list keeps naming phase 3 as it does today; a
  second description of the mode there would be a second source for the same
  contract. The command file's frontmatter description gains one clause so the
  mode is discoverable when the phase fires.
- The review delegate is charged with the standalone
  `/bench-review-implementation` contract, not with a restatement of it. How
  that phase fans out to its three axes, what each axis hunts, and when it
  persists `reviews/<spec-slug>.md` all stay that command's and `craft-review`'s
  business. The pickup artifact needs no new rule: the existing rule already
  writes it only for findings needing a *later* fix pass, so a `--full` run that
  resolves findings in-run writes none, and one that stops short writes the
  pickup by the same rule. Done-claim verification is likewise `craft-delegate`'s
  and invariant 1's; the section points at them.
- Handoff persistence needs no CLI change. `bench handoff` derives the pin block
  and preserves the State section rather than writing orchestration facts into
  it, so the phase reached — the one fact this mode adds — is written by the
  orchestrator into State, and `bench handoff --next <command>` then refreshes
  the derived facts around it. The `--next` value is passed through verbatim, so
  the harness-native form is the orchestrator's to supply.
- Tiers are named by binding — "the mid binding", "the top binding" — never by
  model token. `.bench/lines.env` is the machine-readable source and
  `projects/<name>.md` carries the narrative binding; a token in shipped command
  prose would be a third copy and would be wrong in any linked repo. No gate
  check scans command prose for model tokens, so this one is review-graded, not
  gate-graded.
- The Codex route in story 8 is specified at the level of what must be set —
  the `codex exec` entry point, the model, the reasoning effort, the working
  directory, and a non-interactive approval posture — not as a pinned flag
  string, because flag spellings rot and `codex exec --help` is authoritative.
  The map's Handoff item 7 flagged these mechanics as unverified; they were
  verified in the tree while writing this spec, so the route is specified as
  real rather than as a stub reporting unavailability.
- All three shared rules — fix-don't-park, the Roles assume/verify sentence, and
  story 15's warrant rule — ride `checkSharedRuleSingleSource`'s existing marker
  list rather than new `requireCollapsed` calls. One list entry gives the
  must-be-in-BENCH.md and must-not-be-in-AGENTS.md/README.md halves from a
  single source, which is exactly the map's item 9 watch-out; a separate anchor
  would state the same fact twice.
- Fixture density follows the family's established fixture-to-assertion ratio
  rather than one fixture per assertion, so a new assertion inside an
  already-proven check does not automatically earn a fixture. The four fixtures
  cover four distinct classes, and that is the honest floor given the gate is
  currently canary-bound.
- The scope fence, the coverage-map accounting, the record-citation rule, and the
  falsification-pass offer are four more section-scoped require/`forbid` pairs
  inside the already-bounded
  `--full` section, not new checks. They earn no new fixture: the family's density
  rule buys a fixture per class, and story 9's four already prove this feature's
  section-scoped anchor class. `bench coverage <spec>` is the enumeration the
  accounting rule points at — the rows are its output, so the rule adds a duty
  rather than a second source for what they are.
- The `.claude/commands/` mirror is updated in the same change. No gate check
  compares the two command trees for content; that gap is named in Out of scope.

## Testing decisions

- A good test here drives the conformance check over a tree and asserts the
  diagnostics it returns, or drives the canary sweep over a fixture and asserts
  it bites. Neither reads the prose to grade its meaning — the anchor string and
  its section are the contract, and the fixture is what proves the anchor fires.
- Seams: `checkDocsCurrencyAndWorkflow` in `internal/conformance` (prior art:
  the existing `require`/`requireCollapsed`/`forbid` assertions and the
  section-scoped `checkStructuredPhaseContract` in
  `docs_workflow_helpers_test.go`), `checkSharedRuleSingleSource` in the same
  package (prior art: its existing marker list and the `shared-rule-drift` /
  `readme-shared-rule-drift` fixtures), and the canary sweep over
  `tests/canary/workflow-guidance-anchors/` (prior art: every existing fixture
  in that family, each one file under `files/dot-agents/`).
- Gate: `bench gate` (dev tier) must be green. The canary phase within it is
  what proves every new fixture bites and none is vacuous.
- The orchestration behavior itself has no gate seam and is not given a fake
  one. Story 10 is its observation, and it observes records rather than the
  run's own report.

### Seam diagram

Prose-anchor seam (stories 1, 3, 4, 5, 6, 7, 8, 11, 12, 13, 14, 15):

    trigger: gate conformance phase → checkDocsCurrencyAndWorkflow(root, kitRoot)
        │
        ▼
    .agents/commands/*.md  ──▶ [ markdownH2Section → require / forbid ] ──▶ diags
    .bench/BENCH.md        ──▶ [                                      ]
                      ◀ tests attach here: the assertions run against the kit root
                        every gate; canary fixtures run the same check against a
                        tree with one anchored fact removed and match EXPECT

Shared-rule marker seam (stories 2, 15):

    trigger: gate conformance phase → checkSharedRuleSingleSource(root)
        │
        ▼
    .bench/BENCH.md ──▶ [ marker list: required here,     ] ──▶ diags
    AGENTS.md       ──▶ [ forbidden in the other two      ]
    README.md       ──▶ [                                 ]
                      ◀ tests attach here: add the marker before the sentence and
                        the check reds naming it; restate it in AGENTS.md or
                        README.md and the same entry reds the other direction

Canary-bite seam (story 9):

    trigger: gate canary phase → Sweep(tests/canary, runner)
        │
        ▼
    fixture files/ tree ──▶ [ nested gate, conformance scope ] ──▶ output
    fixture EXPECT      ──▶ [ bite + vacuity grading         ] ──▶ errs
                      ◀ tests attach here: the sweep reds did-not-bite when the
                        mutation produces no matching diagnostic, and vacuous when
                        the baseline already emits it

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | the fix-don't-park rule is present in `.bench/BENCH.md` and inside its Workflow section, not merely somewhere in the file | shared-rule marker seam plus a section-scoped assertion | new marker plus the section-scoped assertion added before the prose; conformance reds `shared rule missing from canonical .bench/BENCH.md` and reds again when the sentence sits outside Workflow | the marker alone passes a rule pasted into any section or an HTML comment, since `requireCollapsed` is whole-file and does not strip comments; the section scope is what makes the anchor mean placement |
| 1 | the rule appears in `.bench/BENCH.md` and nowhere else | shared-rule marker seam | already covered — the same marker's `AGENTS.md`/`README.md` half, proven biting by the `shared-rule-drift` and `readme-shared-rule-drift` fixtures | restating a shared rule in `AGENTS.md` is the drift the map's item 9 named; the existing fixtures prove the direction fires for any marker in the list |
| 2 | `NEVER assume, always verify` reds when `.bench/BENCH.md` loses it | shared-rule marker seam | new marker added before any prose change; conformance reds naming the marker on a tree with the sentence removed | the Roles sentence is currently unanchored, so it can drift out silently — the exact failure the map's #6 opened on |
| 2 | all three markers red when `AGENTS.md` or `README.md` restates them | shared-rule marker seam | restate any of the three sentences in `AGENTS.md` on a scratch tree; conformance reds `shared rule duplicated in AGENTS.md` naming it | a presence-only assertion added outside the marker loop would satisfy the removal direction while leaving the duplication direction unguarded — the cheapest wrong implementation of story 2 |
| 3 | `/bench-shape-idea` carries the look-it-up-before-asking hook | prose-anchor seam | new `requireCollapsed` added before the sentence; conformance reds with its diagnostic | without the assertion the hook is ordinary prose that a later trim removes; the diagnostic names which phase lost it |
| 3 | `/bench-implement-spec` carries the verify-against-the-tree-not-memory hook | prose-anchor seam | new `requireCollapsed` added before the sentence; conformance reds with its diagnostic | same class, different owner — one assertion per phase is what makes the diagnostic actionable rather than a generic "a hook is missing" |
| 3 | `/bench-review-implementation` carries the cite-what-you-read hook | prose-anchor seam | new `requireCollapsed` added before the sentence; conformance reds with its diagnostic | a review that recalls instead of citing is the observed failure; the hook must be pinned in the review phase specifically |
| 3 (edge of 3) | the three hooks sit at their phases' points of use rather than anywhere in the file | manual, review-graded | not TDD-able — the three files place their hooks in three different structures, so no single section name scopes all three, and section-scoping each one costs more machinery than the placement risk warrants | recorded as a stated limitation rather than an implied guarantee: the gate proves the hooks exist, and the review axis grades where they sit |
| 4 | `--full` is documented as one bounded section, and every anchor for this feature resolves inside it | prose-anchor seam | the section-scoped helper is wired before the section is written; conformance reds that the `--full` section is absent | this row is what gives the other `--full` rows their meaning — without a bounded section, every fragment below could be scattered through unrelated parts of the file and still pass |
| 4 | `--full` is opt-in, and no default-on or spec-inferring sentence sits beside it | prose-anchor seam | new require/`forbid` pair added before the section is written; conformance reds on the require, then reds again when a default-on sentence is added to a scratch tree | the cheapest wrong implementation adds the opt-in sentence and keeps a contradicting one, which the map rejected in #2; the require alone cannot see the contradiction |
| 4 | `--full` with a missing or unknown spec argument refuses and says which | prose-anchor seam | new require/`forbid` pair on the refusal and on an inference fallback; conformance reds before either is written | an inferred spec would let a full unattended run build the wrong target with every phase green |
| 4 (edge of 3) | the three standalone phase commands keep the specific contracts their existing anchors pin | prose-anchor seam | already covered, narrowly — the existing anchors on `bench-review-implementation.md`, `bench-final-check.md`, and `bench-debug.md` red if those pinned clauses move; anything they do not pin is review-graded | the degenerate implementation folds the standalones into the orchestrator, which necessarily moves pinned clauses; the row claims only what the anchors actually cover |
| 5 | the review runs in a fresh-context delegate given the spec and the diff only, with inline self-review and context inheritance both closed | prose-anchor seam | new require on the fresh-context-and-inputs statement plus a `forbid` on an inline-review fallback; conformance reds before the section is written | the cheapest wrong implementations are reviewing inline to save tokens and passing the implementing context to a delegate; naming the permitted inputs in the anchored sentence is what separates a real fresh context from a delegate in name only |
| 5 | the review delegate's done-claim is verified against the gate and `git status` before it is acted on | prose-anchor seam | new `requireCollapsed` on the verification pointer; conformance reds before the section is written | the map's Handoff item 6 assigns this hostile input to invariant 1; without the pointer an unattended run accepts a claim as a result, which is the failure invariant 1 exists to prevent |
| 5 | a review delegate that cannot run or returns nothing stops the run at that boundary | prose-anchor seam | new `requireCollapsed` on the stop clause; conformance reds before it is written | proceeding to final-check with review unrun would land a spec whose semantic pass never happened while every gate stayed green |
| 6 | concrete defects are fixed and re-gated in-run; design and judgment findings are flagged for veto, not applied | prose-anchor seam | new require on both halves plus a `forbid` on a stop-on-concrete-finding or apply-the-judgment-finding phrasing; conformance reds before the section is written | one half without the other gives either a stalling run or a run that silently redesigns; the forbid is what catches the version that states the rule and contradicts it two sentences later |
| 6 (edge of 3) | the repair bound and the delegate repair allowance are pointed at, not restated | manual, review-graded | not TDD-able — no assertion can distinguish a pointer from a faithful restatement, and a `forbid` on the originals' wording would fire on the originals themselves | recorded honestly rather than claimed as covered: the existing anchors keep the originals in place, and the Standards axis grades the new section for a second derivation |
| 7 | the `--full` section states that each phase boundary writes the phase reached into State and refreshes the pin block via `bench handoff --next <command>` | prose-anchor seam | new `requireCollapsed` naming both the State write and the refresh; conformance reds before the section is written | the degenerate implementation calls `bench handoff --next` alone, which derives repo, branch, spec, and gate but no phase — so a resuming session learns everything except where the run stopped |
| 7 | the section points at the existing handoff contract for every field it does not add | prose-anchor seam | already covered — `handoff_single_source_test.go` and the existing handoff anchors own the contract; a second enumeration here is graded by the Standards axis | re-enumerating the pin block's fields would be a second derivation of one fact, the duplication the code standard names |
| 7 | a re-invoked `--full` resumes from the phase the handoff names rather than re-implementing | prose-anchor seam | new `requireCollapsed` on the resume clause; conformance reds before it is written | re-run idempotency: without it, a resumed run redoes landed work, and the second implement pass would fight the first's commits |
| 8 | escalation pauses and asks, with three named routes and a per-run recommendation, omitting a route only when this harness cannot invoke it and saying so | prose-anchor seam | new require enumerating the three routes and the ask, plus a `forbid` on a silent-escalation phrasing; conformance reds before the section is written | enumerating all three routes is what stops a build shipping a one-route prompt; the forbid catches an escalate-if-obviously-needed sentence that would restore silent escalation |
| 8 | whether a given run's trigger fires, and whether the ask actually happens | manual, dogfood | not TDD-able — the trigger is judgment prose, which the map's Handoff item 7 accepted; an orchestrator that always decides escalation is unnecessary looks identical to the gate | story 10's independent records are the observation: a run whose diff was large and which never recorded a question is the visible signature |
| 8 | the escalation prose names tiers by binding, never by model token | manual, review-graded | not TDD-able against the current oracle — `checkLineBinding` validates `.bench/lines.env` and the profile paragraph and never scans command files | recorded as review-graded rather than claimed covered: a token in shipped command prose is a third copy of the binding and wrong in every linked repo, and only the Standards axis sees it today |
| 9 | four fixtures exist, one per class — verify hook, fresh-context review, boundary handoff, escalation ask — and no two exercise the same class | canary-bite seam plus the fixture inventory | each fixture's EXPECT names a different diagnostic; a duplicate class is visible as two fixtures sharing one EXPECT string | the cheapest wrong migration adds four biting fixtures that all remove verify hooks, satisfying bite and vacuity while leaving three `--full` anchors unproven; enumerating the four classes per row is what forbids it |
| 9 | each fixture bites: removing its anchored fact produces its EXPECT | canary-bite seam | the sweep reds did-not-bite for any fixture whose mutation the check does not diagnose | a fixture that does not bite is a fixture that proves nothing; the sweep enforces this by construction for every new fixture |
| 9 | no new fixture is vacuous | canary-bite seam | already covered — the sweep grades every EXPECT against its scope group's empty-tree baseline | an EXPECT the baseline already emits would pass without the mutation, making the fixture decorative |
| 9 | fixture base names stay globally unique across families | canary sweep | already covered — the existing `fixtures()` uniqueness check | a collision with a `behavior-owned` name would misroute scope resolution once FT91's migration lands |
| 10 | the first real `--full` run is graded against independent records — the delegate invocation, `git log`, and the committed `session-handoff.md` at each boundary — not against its own exit report | manual, dogfood | not TDD-able — orchestration behavior is prose-driven and has no gate seam (map Handoff item 5) | a run that silently self-reviewed or skipped a boundary rewrite can still emit a correct-looking exit report; only records the run did not author separate the two |
| 11 | the `--full` section fences the run's diff to the spec's stories and seams, with work outside them recorded rather than built | prose-anchor seam | new require on the fence plus a `forbid` on an opportunistic-improvement phrasing; conformance reds before the section is written | fix-don't-park is a one-way ratchet, and the cheapest wrong reading discharges a defect by widening the diff into a refactor no story asked for; the require alone passes a section that states the fence and licenses the widening two sentences later |
| 11 (edge of 3) | whether a given run's diff actually stayed inside the spec's seams | manual, review-graded | not TDD-able — the fence is judgment prose, and no gate check compares a diff against a spec's named seams | recorded rather than claimed covered: `/bench-review-implementation`'s Spec axis already grades a diff against its spec, and this row claims only that the anchor is present |
| 12 | the exit report accounts for every coverage-map row by name, with the no-map fallback stated | prose-anchor seam | new require naming the per-row disposition and `bench coverage <spec>` as the enumeration, plus a `forbid` on an aggregate-completion phrasing; conformance reds before the section is written | a run that implemented eight of twelve rows emits an aggregate claim identical to a complete one; the enumeration is what makes the four missing rows visible, and the forbid catches the version that states the rule and then summarizes anyway |
| 13 | the exit report cites the record behind each phase claim | prose-anchor seam | new require on the citation rule plus a `forbid` on an uncited phase-complete phrasing; conformance reds before the section is written | story 10 buys this observation once, by hand; without the anchored rule every run after the first is graded on its own unsupported report, which is the failure story 10's row already names |
| 14 | the falsification pass is offered on story 8's trigger, as its own question rather than a fourth route, and never runs standing | prose-anchor seam | new require naming the trigger, the refutation charge, and the separate question, plus a `forbid` on a standing-pass phrasing; conformance reds before the section is written | the two cheapest wrong readings are folding it into story 8's menu — which breaks that story's fixed-three anchor — and making it standing, which doubles review cost on every three-file run; the forbid catches the second, and naming the separate question closes the first |
| 14 (edge of 3) | whether a given run's trigger fires and the pass is actually offered | manual, dogfood | not TDD-able — it shares story 8's judgment trigger, which the map's Handoff item 7 already accepted as unobservable to the gate | story 10's independent records are the observation, by the same argument story 8's row makes: a large-diff run with no recorded question is the visible signature |
| 15 | the warrant rule is present in `.bench/BENCH.md` and inside its communication section, not merely somewhere in the file | shared-rule marker seam plus a section-scoped assertion | new marker plus the section-scoped assertion added before the prose; conformance reds `shared rule missing from canonical .bench/BENCH.md` and reds again when the sentence sits outside the communication rules | the same argument story 1's first row makes, and it binds harder here: a whole-file marker passes a rule pasted into any section or commented out, and this rule only works if it is loaded at the moment a claim is made |
| 15 | the rule appears in `.bench/BENCH.md` and nowhere else | shared-rule marker seam | already covered — the same marker's `AGENTS.md`/`README.md` half, proven biting by the `shared-rule-drift` and `readme-shared-rule-drift` fixtures | a third entry rides the loop the first two ride; the existing fixtures prove the direction fires for any marker in the list |
| 15 | the rule asks for disclosure of what went unread, not for thoroughness | prose-anchor seam | new `forbid` on a be-thorough phrasing; conformance reds on a scratch tree carrying one | a be-thorough restatement is the cheapest wrong implementation and is worse than nothing: it reads as compliance while staying unfalsifiable, so no reviewer can ever catch a breach of it |
| 15 (edge of 3) | whether a session actually discloses what it left unread | manual, review-graded | not TDD-able — the gate grades the tree and cannot observe a conversational claim; no seam exists and none is invented | recorded as the rule's honest limit rather than an implied guarantee. Its whole design concedes this: it cannot prevent a breach, only make one visible in the message where it happens, which is exactly why the anchored form is disclosure rather than thoroughness |
| 11, 12, 13, 14 | each new `--full` anchor bites when its fact is removed | canary-bite seam | already covered — story 9's `--full` fixtures prove the section-scoped require/`forbid` class in `.agents/commands/`, and the family's density rule earns no fixture for a new assertion inside a proven check | flagged for veto: four more fixtures would restate the proven class rather than extend coverage, so the bite claim for these anchors rests on the class, not on a fixture each |
| 15 (edge of 3) | story 15's section-scoped anchor in `.bench/BENCH.md` bites when the sentence moves out of the communication rules | manual, review-graded | not TDD-able as specified — the `shared-rule-drift` fixtures prove the *marker* directions for any list entry, and story 9's four fixtures prove section scoping only in `.agents/commands/`; no fixture proves a section-scoped anchor over `.bench/BENCH.md` | recorded rather than claimed covered, and it names a gap this spec inherits rather than creates: story 1's Workflow-section half rests on the same untested combination. If the reviewer wants either proven, one fixture removing a section-scoped `.bench/BENCH.md` sentence covers both stories at once |

### Edge inventory

- Missing input (`--full` with no spec) and malformed input (`--full` with an
  unknown path) → coverage rows (story 4).
- Error path (review delegate unavailable or empty-handed) → coverage row
  (story 5).
- Hostile input (a review delegate's done-claim taken as a result) → coverage
  row (story 5), per the map's Handoff item 6.
- Interrupted or partial state (run killed mid-phase) → coverage rows (story 7):
  the boundary rewrite plus the tree-wins rule already in `session-handoff.md`
  are the recovery, and `bench status` reports the handoff's age.
- Re-run idempotency (`--full` re-invoked on a completed or partly-completed
  run) → coverage row (story 7).
- Boundary value (final gate red after the repair pass) → **Won't handle**
  beyond existing behavior: `/bench-review-implementation`'s terminal-repair
  rule and `/bench-implement-spec`'s "When the build stops short" route already
  own the cap and the exit, and a second bound here would be a second source.
- Hostile environment (`bench` not on PATH at a phase boundary) → **Won't
  handle**: `session-handoff.md` is a markdown file any editor writes, so there
  is no CLI dependency to fall back from; `bench handoff --next` is the route
  that derives the pin block, not a prerequisite for writing the file.
- Hostile environment (Codex CLI absent when escalation is offered, or when the
  story 14 falsification pass is) → covered by story 8's omit-and-say-so clause,
  which story 14 adopts rather than restates; no separate row.
- `--full` on a spec already `Status: implemented` → **Won't handle**:
  `/bench-final-check` owns the status transition and already reports an honest
  green for a branch with nothing to commit, so the run terminates correctly
  without a new rule.
- Empty input (a spec with no acceptance coverage map) → coverage row (story 12):
  the exit report states the map's absence and accounts for the user stories
  instead. `/bench-implement-spec` already handles building without a map, and
  `--full` changes nothing about that path.
- Concurrent runs (two `--full` runs on one tree) → **Won't handle**: invariant 1
  already forbids a gate verdict answering for two diffs, and `bench worktree`
  is the existing isolation route.
- Compatibility probe → n/a: this spec names no external format or protocol.
  The one external surface, the Codex CLI, was probed in the tree while writing
  this spec rather than promised.

## Out of scope

- **A command-mirror content check** — a conformance check proving
  `.agents/commands/` and `.claude/commands/` are byte-identical. A separate
  capability: it decides whether the Claude tree is generated-and-verified or
  hand-maintained, which is an architectural call this map never opened, and it
  guards a drift class wider than this feature. Today only `.agents/commands/`
  is anchored, so a mirror that drifts is gate-invisible. ~4 edits, ~2 gate runs.
- **A model-token sweep over shipped command prose** — a check that no
  `.agents/commands/*.md` names a tier's model token, which would turn story 8's
  binding rule from review-graded into gate-graded. A separate capability: it
  polices every command file against `.bench/lines.env`, not just this feature's
  section. ~5 edits, ~2 gate runs.
- **Making orchestration the default invocation** — rejected in map #2; recorded
  here so it stays closed rather than reappearing as a convenience.
- **A `bench` subcommand for the full run** — a separate capability: it would
  move orchestration from prose into the CLI, which changes the enforcement
  layer rather than extending this one, and it needs its own map. ~30 edits,
  ~6 gate runs.
