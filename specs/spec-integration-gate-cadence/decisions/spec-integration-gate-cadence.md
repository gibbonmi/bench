# Spec integration branch and gate cadence

Status: ready

## Destination

Spec-backed implementation gets one integration branch rooted at the active
working branch. Independent frontier tickets run concurrently in separate
worktrees and return attributed work to that branch. Focused acceptance checks
grade ticket work while it is provisional; one whole-project gate over the
composed candidate authorizes promotion into the working branch.

This map decides what provisional work may claim, how ticket results integrate,
where semantic review and repairs occur, what invalidates a composed green
verdict, how strongly harnesses must fill the parallel frontier, and whether
the new cadence applies outside spec-backed implementation. It does not split,
scope, or weaken the gate itself.

## #1: Where does the project-green invariant apply?

Blocked by: none
Type: Grill

### Question

Must every ticket commit remain project-green, or may the spec integration
branch carry provisional commits that have passed only their focused acceptance
checks while the active working branch remains green?

### Answer

The project-green invariant applies at promotion into the active working branch,
not at every provisional checkpoint on the spec integration branch. A ticket
checkpoint may be committed after its focused acceptance evidence passes, but
it is explicitly provisional: it is not done, shippable, or eligible for the
working branch until the exact composed candidate passes the whole-project
gate.

The working branch remains green. The gate remains the sole authority for
promotion, and this decision changes its frequency and subject boundary rather
than splitting, scoping, or weakening what it proves.

## #2: What evidence makes a ticket eligible for integration?

Blocked by: #1
Type: Grill

### Question

If ticket commits no longer run the whole-project gate, which focused evidence
must the delegate return before the coordinator may integrate the work?

### Answer

Every charged coverage row supplies its focused red-to-green signal before the
ticket may enter the spec branch. Rows already covered or not TDD-able retain
their reviewed classification and evidence instead of manufacturing a red.
The delegate also runs the relevant package, type, lint, or static checks named
for the touched seam.

The coordinator verifies that the diff stays inside the ticket's ownership
fence, that no unexplained worktree files remain, and that one accepted behavior
passes through an independent probe the delegate did not author. These checks
authorize a provisional ticket checkpoint only. They do not run the
whole-project gate and do not support a project-green, done, or shippable claim.

## #3: How do parallel ticket results enter the spec branch?

Blocked by: #1
Type: Grill

### Question

Should each delegate create an attributed commit on its own assignment branch
for serialized integration, or return an uncommitted diff for the coordinator
to commit on the spec branch?

### Answer

Each delegate leaves one assignment worktree ready after its focused checks and
reports the evidence needed by ticket #2. After the coordinator performs the
independent probe, it invokes the sanctioned checkpoint porcelain with a
canonical evidence receipt. That command verifies the ownership fence and
creates one provisional commit attributed to the assignment without running the
whole-project gate. No delegate checks out or writes the spec integration branch.
The coordinator is its sole mutation owner and integrates eligible checkpoint
commits one at a time through sanctioned Bench porcelain.

Integration uses the current spec tip as a compare-and-swap precondition.
An assignment base made older only by a successfully integrated sibling is
eligible for the clean replay decided in #16. An overlapping patch, merge
conflict, ownership drift, or changed assumption fails before the spec branch
changes and routes the ticket back to its owning delegate. A successfully
integrated commit retains its source attribution and recovery relationship, and
the coordinator may then release that ticket's worktree. The later history
decision still determines whether those ticket commits remain visible after
promotion or are squashed.

## #4: Does semantic review precede the composed gate?

Blocked by: #2, #3
Type: Grill

### Question

Should Standards, Spec, and Coverage review run over the composed provisional
branch before its whole-project gate, so accepted repairs join the candidate
before the gate is paid?

### Answer

Yes. Standards, Spec, and Coverage review runs over the exact composed
provisional spec branch before its whole-project gate. The review is a fresh,
read-only charge and makes no project-green claim.

Every actionable finding receives a disposition before gate entry. Accepted
repairs become attributed provisional commits and pass the same focused
evidence and coordinator-integration rules as implementation tickets. The gate
runs only after the reviewed candidate has stopped changing, so a predictable
review repair does not consume and invalidate an earlier whole-project run.

## #5: What invalidates the composed green verdict?

Blocked by: #3, #4
Type: Grill

### Question

Which changes require another whole-project gate before promotion: any candidate
tree change, a working-branch advance, a history-only integration change, or
some smaller set?

### Answer

Green evidence is bound to the exact closed subject: the candidate tree, gate
and launcher closure, declared tools and environment, execution-policy version,
and freshness window. Any change to that subject invalidates the evidence and
requires the whole-project gate before promotion.

A working-branch advance first recomposes the candidate on the new base. A
different resulting closed subject requires a gate; a provably identical
subject may reuse its current green evidence. Commit identity and history alone
are not gate inputs, so a rebase or history-only integration that produces the
same tree and oracle does not manufacture another run. Promotion still moves
the branch-scoped green ref to the exact commit that entered the working branch.

## #6: What parallelism is mandatory?

Blocked by: #2, #3
Type: Grill

### Question

When two or more ownership-safe frontier tickets are independent, must the
coordinator fill the harness concurrency limit and report why any slot stayed
unused, or is parallel fan-out advisory?

### Answer

When two or more ownership-safe frontier tickets are independent, the
coordinator must dispatch them concurrently up to the smaller of the frontier
size and the harness concurrency limit. It continues coordinating while they
run instead of waiting on each result before starting the next eligible charge.

Any unused capacity is reported with its concrete reason: a dependency edge,
an overlapping ownership fence, an unavailable harness slot, or a measured
resource constraint. Sequential execution without one of those reasons is a
phase-contract failure, not coordinator discretion. Whole-project gates remain
a serialized resource and do not run in ticket worktrees.

## #7: What happens after the composed gate goes red?

Blocked by: #4, #5
Type: Grill

### Question

Do red phases become focused repair tickets on the same spec branch, with the
whole gate rerun only after the candidate changes, or does a red restore
per-ticket whole-gate cadence?

### Answer

A red composed gate never promotes and is never reused. The spec branch remains
provisional while the coordinator attributes each red to the candidate,
inherited state, or infrastructure.

Candidate-owned failures become focused repair tickets. Ownership-safe repairs
fan out in parallel, pass ticket #2's evidence, and integrate through ticket
#3's coordinator-owned path. A fresh semantic review grades the entire repaired
composition, with the repair delta called out, before one whole-project gate
reruns. This cadence repeats until green or the stage exhausts its declared
bound; it never restores per-ticket whole gates.

An inherited or infrastructure red follows the existing diagnosis and retry
posture without manufacturing a code repair. It still cannot authorize
promotion, and any later tree change produces a new closed subject.

## #8: Does the cadence apply outside spec-backed builds?

Blocked by: #1
Type: Grill

### Question

Should light-path changes and `bench shift` keep their current per-commit green
contract, or move to the same provisional-branch model?

### Answer

No. The provisional integration-branch cadence applies only to spec-backed
builds with a reviewed decision source, ownership-fenced ticket graph, focused
acceptance evidence, and composed semantic review.

Light-path changes and `bench shift` retain their current commit-on-green
contract. They may be reconsidered only after the spec-backed workflow supplies
measured wall-clock, gate-count, red-attribution, and recovery evidence; this
map grants no standing extrapolation from the new cadence.

## #9: Where does the literal project-green marker live?

Blocked by: #1
Type: Grill

### Question

The repository-wide last-gate record is transient candidate evidence: a red or
green run on another worktree replaces it even when the active working branch
has not moved. Should project-green be a Git-common-directory ref scoped to the
working branch, a branch-scoped data record, or remain derived only from the
latest gate cache?

### Answer

Use a custom Git ref in the repository's common directory:
`refs/bench/green/<working-branch>`. The ref points to the exact working-branch
commit most recently authorized by the whole-project gate. Git natively owns
the ref's durable storage, locking, compare-and-swap update, and visibility
across every linked worktree; Bench owns the meaning.

The marker is authoritative only when the working branch tip equals the green
ref and subject-addressed gate evidence for that commit still matches the
current tree, oracle closure, and freshness policy. The existing
`bench-last-gate` record remains the transient latest-run view, but it cannot be
the durable project-green proof because another worktree's run overwrites it.

Promotion verifies a clean expected working-branch base, a current prospective
tree, and green evidence for that exact closed subject. It creates the
spec-level squash as a child of the expected base, advances the working branch
to that commit, and then updates the green ref. A crash between the branch and
green-ref updates leaves the branch advanced but unmarked, which fails closed
and is repairable from the retained subject evidence; it never leaves an
unverified branch marked green. Any later branch advance makes the two refs
diverge and removes project-green status automatically.

## #10: Does promotion preserve provisional ticket commits?

Blocked by: #3
Type: Grill

### Question

When a green spec candidate enters the working branch, should its attributed
ticket commits remain visible, or should promotion squash them into one
spec-level commit?

### Answer

Promotion squashes the reviewed, gate-green candidate into one spec-level
commit whose parent is the expected working-branch tip. The working branch
therefore receives only a commit whose complete tree the whole-project gate
authorized; intentionally provisional ticket commits never enter its ancestry
or become bisect targets.

The spec integration ref and its attributed ticket commits remain retained as
spec evidence under the lifecycle owner chosen in #11. Their coverage receipts,
review relationship, and source attribution remain inspectable after promotion.
Because commit identity is not part of the closed gate subject, the prospective
status-updated tree's green evidence authorizes the byte-identical squash tree.
Promotion then points the working branch and its project-green ref at the new
squash commit.

## #11: Which Bench porcelain owns the spec-branch lifecycle?

Blocked by: #8, #10
Type: Grill

### Question

Should one spec-scoped command family create or resume the integration branch,
cut ticket worktrees from its expected tip, integrate provisional commits, run
the composed gate, and promote the green candidate, while the harness remains
responsible only for agent dispatch and review?

### Answer

One `bench spec` porcelain family owns the complete repository-state lifecycle:
durable spec-run state; integration and evidence refs; idempotent start and
resume; ticket worktrees rooted at the expected spec tip; compare-and-swap
integration; review-head recording; composed gate execution; squash promotion;
project-green marking; retained history; and explicit recovery or abandonment.

The harness owns ticket derivation, agent dispatch, focused evidence, finding
disposition, and semantic review. It invokes the porcelain with those results
but never reconstructs Git transitions in prose. Every mutating operation
fails closed on a dirty or unexpected working branch, stale spec tip, missing
evidence, or ownership mismatch, and re-entry reports the durable next action
instead of creating a second run.

## #12: What is the spec-build command vocabulary?

Blocked by: #11
Type: Grill

### Question

Which stable `bench spec` subcommands expose start/resume, parallel assignment,
provisional integration, semantic-review evidence, atomic gate-and-promotion,
status, and safe abandonment without making harnesses compose plumbing calls?

### Answer

The stable family is:

- `bench spec build start <slug>` creates or resumes the integration branch.
- `bench spec build assign <slug> --ticket <ticket> --request <id>` cuts an
  ownership-tracked ticket worktree at the current spec tip.
- `bench spec build checkpoint <slug> --assignment <id> --evidence <receipt>`
  validates subject-bound focused evidence and the ownership fence, then creates
  the attributed provisional assignment commit without running the gate.
- `bench spec build integrate <slug> --assignment <id>` consumes a verified
  checkpoint and compare-and-swap integrates its provisional commit.
- `bench spec build review <slug> --evidence <receipt>` records
  semantic-review evidence against the exact candidate tip.
- `bench spec build status <slug>` reports durable state and the next action.
- `bench spec build promote <slug>` verifies current review evidence, runs or
  reuses the composed whole-project gate, creates the working-branch squash
  commit, and updates its project-green ref.
- `bench spec build abandon <slug>` provides plan/apply safe abandonment while
  preserving recoverable work.

`start` is idempotent and is also the resume route. `promote` is the only command
in this family that runs the whole-project gate. Harnesses do not compose lower
level Git or worktree plumbing to emulate any of these operations.

## #13: How does focused evidence enter durable build state?

Blocked by: #2, #3, #11, #12
Type: Grill

### Question

The harness owns focused checks, the independent probe, and semantic review, but
the accepted porcelain has no evidence input and current `bench commit` always
runs the whole-project gate. Should the family add an explicit provisional
checkpoint that ingests a canonical evidence receipt and creates the attributed
assignment commit, with `integrate` consuming only that verified checkpoint?

### Answer

Yes. Add
`bench spec build checkpoint <slug> --assignment <id> --evidence <receipt>`.
The harness gathers the delegate's focused red-to-green results and the
coordinator's independent probe in one canonical receipt. `checkpoint`
validates that the receipt, assignment base, resulting tree, charged coverage
rows, seam checks, and ownership fence agree; creates the one attributed
provisional commit; and durably binds the receipt to it without running the
whole-project gate.

`integrate` accepts only a verified checkpoint and remains a narrow,
compare-and-swap candidate mutation. The same transport applies to
`review <slug> --evidence <receipt>` so semantic judgment remains with the
harness while subject binding and durable state remain inside the porcelain.
The harness never creates a provisional commit or edits lifecycle refs itself.

## #14: When does the spec become implemented?

Blocked by: #4, #10, #12
Type: Grill

### Question

The promoted squash must be byte-identical to the tree authorized by the gate,
but the current workflow flips the spec to `Status: implemented` only while
creating the landing commit. Should promotion construct the canonical status
transition as a prospective tree, gate that exact tree, and publish it only on
green so a red run leaves both the candidate and working branch unchanged?

### Answer

Yes. After current semantic-review evidence is accepted, `promote` applies the
canonical spec-status transformation to an unpublished prospective tree and
runs or reuses the whole-project gate only for that exact tree. A red gate or
operational failure publishes nothing: the provisional candidate, its spec
status, the working branch, and the project-green ref remain unchanged.

On green, the spec-level squash commit uses the prospective tree byte for byte.
The working branch therefore never observes `Status: implemented` without the
same complete tree already being gate-authorized. The status transformation is
a sanctioned, deterministic promotion step; it does not require a second gate
after the landing commit.

## #15: How is the first project-green marker bootstrapped?

Blocked by: #5, #9, #12
Type: Grill

### Question

When a working branch predates `refs/bench/green/<working-branch>`, should
`start` initialize the marker only from retained green evidence for the exact
current subject, otherwise refuse with a route to one baseline `bench gate`
rather than silently trusting the branch or making `start` another gate caller?

### Answer

Yes. A missing branch-scoped green ref may be initialized without another gate
only when the gate owner can recover retained green evidence whose complete
closed subject exactly matches the current working-branch tip. `start` asks that
owner to validate and atomically establish the marker; it never interprets the
latest-run cache or synthesizes authorization itself.

If exact reusable evidence is absent, stale, or unavailable, `start` refuses
before creating build state and gives the single recovery route: run one
baseline `bench gate`, then retry `start`. The next `start` imports that exact
authorization and creates the marker. This is a migration/bootstrap cost only;
`start` does not become a second full-gate caller.

## #16: How do parallel sibling checkpoints integrate?

Blocked by: #3, #6, #13
Type: Grill

### Question

Parallel assignments may share the same candidate base, which necessarily
becomes older after the first sibling integrates. Should `integrate`
automatically replay each later checkpoint onto the current candidate when its
ownership-fenced patch applies cleanly, reserving delegate rework for overlap,
conflict, or changed assumptions rather than rejecting all older bases?

### Answer

Yes. A checkpoint remains permanently bound to its assignment base, checkpoint
tree, attributed patch, and focused evidence, while `integrate` targets the
current candidate tip. The focused receipt authorizes that assignment-scoped
patch and seam; it does not claim to have checked unrelated sibling work. The
porcelain constructs the prospective replay result before mutation and accepts
an older sibling base when the complete ownership-fenced patch is unchanged,
applies cleanly, and its recorded assumptions remain valid.

It then advances the candidate with an exact-old-tip compare-and-swap. A
concurrent candidate move retries from the new tip only while the same
preconditions remain true. Patch overlap, conflict, ownership drift, or an
invalidated assumption leaves the candidate unchanged and routes the assignment
back to its delegate. Candidate advancement by a disjoint sibling is not itself
an error.

## Not yet specified

## Spec-writer discretion

- The private durable-state schema, receipt encoding and versioning, ref-name
  escaping, and CLI rendering, provided their identities remain opaque and
  re-entry remains fail-closed.
- The internal implementation shape for the lifecycle owner. Prefer one deep
  spec-build module with narrow gate-evidence and worktree/CAS collaborators;
  do not expose Git plumbing through the command adapters.
- The clean-replay mechanism used by `integrate`, provided it proves the full
  assignment patch is unchanged, constructs the result before mutation, and
  preserves the attribution and evidence relationships decided above.

## Out of scope

- Diff-scoped gating, per-package verdict reuse, or any file-to-test map.
- Removing, weakening, or skipping a gate check for wall-clock savings.
- Gate-phase scheduling and outer concurrency; those remain separate decisions.
- Ship-tier release verification.

## Sources

- Path: `.agents/commands/bench-implement-spec.md`
  Supports: the current full-gated commit cadence for tickets and accepted repair tickets, plus the final composed check.
  Drift: mutable kit guidance; re-read before resuming this map or compiling its spec.
- Path: `.agents/skills/bench-craft-delegate/SKILL.md`
  Supports: concurrent independent ticket worktrees and serialized whole-project gates.
  Drift: mutable kit guidance; re-read before resuming this map or compiling its spec.
- Path: `decisions/gate-critical-path.md`
  Supports: whole-subject verdict reuse and the rejection of diff-scoped gating.
  Drift: active shaping map; re-read its closed rulings and open frontier before compiling this map.
- Path: `ROADMAP.md`
  Supports: FT162's authoritative full-run subject, FT169's worktree landing primitive, and FT171's separate inside-one-gate contention question.
  Drift: working prioritization document; re-read the named rows before roadmap reconciliation or spec compilation.
