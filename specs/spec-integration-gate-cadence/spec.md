# Spec integration branch and gate cadence

Status: staged

Decision source: compiled map at
`specs/spec-integration-gate-cadence/decisions/spec-integration-gate-cadence.md`.
Its four structured sources were re-read against the current tree on 2026-07-31:
`.agents/commands/bench-implement-spec.md`,
`.agents/skills/bench-craft-delegate/SKILL.md`,
`decisions/gate-critical-path.md`, and the FT162, FT169, and FT171 rows in
`ROADMAP.md`. The current gate still owns exact whole-subject verdict reuse, the
worktree package still owns assignments and recovery, and `bench spec` remains
the existing command seam. No source conflict or unreadable source was found.

## Problem

Spec-backed implementation pays the whole-project gate at every ticket landing,
then pays it again over the composed build. That cadence keeps each commit green,
but it serializes independent work around the most expensive resource and repeats
the same proof before semantic review has seen the composition. It also leaves
the harness reconstructing a stale-base, worktree, commit, integration, and
recovery protocol from prose.

The obvious shortcut is unsafe. Focused checks do not prove that a composed
candidate is shippable, the repository-wide latest gate cache can be overwritten
by another worktree, and a green for one tree cannot authorize a changed tree or
oracle. Bench therefore needs a durable distinction between provisional ticket
evidence and project-green evidence, plus one lifecycle owner that can fail
closed through interruption and re-entry.

## Solution

Every spec-backed build gets one durable run rooted at the active working branch
and one provisional integration ref. Independent, ownership-safe frontier
tickets receive separate assignment worktrees and may checkpoint after their
focused evidence is validated. Bench integrates those attributed checkpoints
serially without claiming that the candidate is green.

The exact composition receives fresh Standards, Spec, and Coverage review.
`bench spec build promote` then constructs the deterministic
`Status: implemented` prospective tree, runs or reuses the whole-project gate
for that exact closed subject, and only on green publishes one spec-level squash
commit into the working branch and its branch-scoped project-green ref. One
`internal/specbuild` state machine owns every transition; the harness supplies
judgment and evidence through the public command family but never composes Git
plumbing.

## User stories

1. As a coordinator, I want `bench spec build start <slug>` to create or resume
   one run only from an exactly project-green working-branch tip, so that
   provisional work starts from an authorized base and a missing marker has one
   safe bootstrap route. Line: gpt-5.6-terra / high. The state and ref bootstrap
   are correctness-critical, and interruption postures are only partly visible
   through the current gate.

2. As a coordinator, I want `assign` to cut ownership-tracked ticket worktrees
   from the current candidate and allow independent siblings to share that base,
   so that the eligible frontier can fill the harness concurrency limit without
   sharing a checkout. Line: gpt-5.6-terra / medium. The existing worktree owner
   supplies the shape, while spec-scoped identity and ticket binding need
   contract coverage.

3. As a coordinator, I want `checkpoint` to accept only a canonical receipt that
   binds the charged rows, focused checks, independent probe, assignment tree,
   and ownership fence, so that an attributed provisional commit can exist
   without being mistaken for a whole-project green. Line: gpt-5.6-terra / high.
   Receipt validation is a new authority seam where a permissive mistake would
   silently promote unproved work.

4. As a coordinator, I want `integrate` to compare-and-swap the candidate and
   cleanly replay an unchanged disjoint sibling patch onto a newer candidate,
   so that parallel ticket completion does not turn into routine delegate
   rework. Line: gpt-5.6-terra / high. Patch identity, ownership, assumptions,
   and concurrent ref movement must agree before mutation.

5. As a reviewer, I want semantic-review evidence bound to the exact candidate
   and every accepted repair to follow the same checkpoint path, so that the
   whole-project gate runs only after the reviewed composition stops changing.
   Line: gpt-5.6-terra / medium. The judgment remains human-owned, while exact
   subject binding and invalidation are gate-observable.

6. As an operator, I want `status`, idempotent re-entry, and plan/apply
   abandonment to report one durable next action and preserve recoverable work,
   so that an interruption never requires reconstructing lifecycle state from
   refs or conversation history. Line: gpt-5.6-terra / high. Recovery spans
   multiple Git and filesystem transitions whose wrong answer can lose work.

7. As a reviewer, I want `promote` to gate the exact prospective implemented
   tree and publish a byte-identical squash plus branch-scoped green marker only
   on green, so that the active working branch never acquires an ungated tree.
   Line: gpt-5.6-terra / high. This is the sole project-green transition and
   therefore the feature's highest-risk executable seam.

8. As a coordinator, I want candidate, review, gate, and working-base drift to
   invalidate only the evidence whose exact subject changed, with red gates
   routing to attributed repairs or diagnosis, so that reuse stays sound and a
   red never restores per-ticket whole gates. Line: gpt-5.6-terra / high. The
   invalidation matrix is subtle and an over-broad reuse would pass cheaply.

9. As an agent author, I want the implementation workflow to mandate frontier
   fan-out, composed semantic review, sanctioned porcelain, and the
   spec-backed-only scope, so that harnesses use the lifecycle consistently
   instead of regenerating the old cadence. Line: gpt-5.6-sol / high. This edits
   kit guidance that steers future sessions, so `craft-line`'s leverage override
   applies and implementation pauses for reviewer approval before this story.

## Implementation decisions

- Add one deep `internal/specbuild` module. Its small public surface implements
  `Start`, `Assign`, `Checkpoint`, `Integrate`, `Review`, `Status`, `Promote`,
  and `Abandon`; it owns run-state transitions, lifecycle locking, ref naming,
  exact-old-value updates, next-action derivation, and recovery. CLI dispatch in
  `internal/spec` remains a thin adapter. No shell shim parses build state or
  invokes Git plumbing.
- Durable state lives in the Git common directory, not a worktree Git directory,
  so every linked worktree sees one run. A spec-scoped lock serializes mutations.
  State is written by durable replace with a versioned schema and records opaque
  run and assignment identities, the working branch and expected base, candidate
  and retained checkpoint refs, ticket and receipt digests, current review
  subject, promotion subject, and terminal/recovery state. The private file
  layout and schema fields remain implementation discretion; `status` is their
  only public projection.
- A slug is resolved through the existing `internal/spec` resolver before it
  identifies a run. Ref components are derived from an opaque digest rather than
  interpolating caller text. The literal project-green ref is
  `refs/bench/green/<working-branch>`, as decided by the map; Git validates the
  working ref, and Bench never manufactures a fallback branch name.
- `start` records the checked-out working branch and exact tip. It asks the gate
  owner whether retained green evidence matches that complete current subject.
  The gate owner validates and compare-and-swap establishes a missing
  branch-scoped green ref; `specbuild` does not parse evidence or write that ref
  around the owner. If validation fails, `start` changes nothing and reports the
  single recovery route: run one baseline `bench gate`, then retry `start`. It
  never treats `bench-last-gate` as authority by itself and never runs the gate.
- Repeating `start` for the same slug and compatible working branch resumes the
  existing run and returns its next action. A second incompatible run, a
  detached HEAD, an unresolved working branch, a dirty working checkout, an
  unexpected branch tip, or incomplete prior state fails closed without
  creating another ref or worktree.
- `assign <slug> --ticket <ticket> --request <id>` resolves a real implementation
  ticket under that spec, snapshots its charged coverage rows, ownership fence,
  and declared assumptions, and asks the existing worktree owner to create one
  request-idempotent assignment rooted at the current candidate tip. Different
  eligible tickets may share the same base. A reused request must name the same
  run and ticket or it refuses.
- The harness, not Bench porcelain, decides which tickets are ready and whether
  their focused evidence is adequate. It records that result in one bounded,
  versioned receipt assembled by the coordinator outside the assignment
  worktree. Delegate evidence and coordinator evidence are separate sections.
  The receipt identifies the run, assignment, assignment base, resulting tree,
  ticket digest, charged row outcomes, relevant package, type, lint, or static
  checks, the coordinator's independent probe command, exit, output digest,
  subject tree, and producer request, ownership inventory, and declared
  assumption digests. The probe producer differs from the assignment owner and
  its subject equals the live resulting tree; copying a delegate-owned pass into
  the coordinator section is invalid. Rows reviewed as already covered or not
  TDD-able retain that classification rather than inventing a red.
- `checkpoint <slug> --assignment <id> --evidence <receipt>` reads only an
  ordinary, size-bounded regular file and validates the receipt against durable
  assignment state and live Git facts. It proves that every charged row has one
  admissible outcome, required seam checks and the independent probe passed, the
  coordinator probe is fresh and independently produced for this assignment
  tree, the worktree has no unexplained paths, and the complete patch stays
  inside the ownership fence. A receipt inside the assignment worktree, carrying
  its owner as probe producer, or bound to another assignment or tree refuses.
  Only then does it create one attributed checkpoint commit and bind the receipt
  to that commit. It does not mutate the candidate and does not call the gate.
- Checkpoint attribution includes opaque run, assignment, and ticket identities
  in state and stable commit metadata. Human-supplied ticket labels or paths are
  display text, never identity. The checkpoint ref and receipt relationship are
  retained after integration and promotion so the provisional history remains
  inspectable without entering the working branch ancestry.
- `integrate` accepts only a verified checkpoint. When its assignment base is
  the candidate tip, it applies the recorded patch and advances the candidate
  with an exact-old-tip compare-and-swap. For an older sibling base, it first
  reconstructs the complete patch, verifies byte identity and the recorded
  assumptions, and constructs the prospective replay result without mutation.
  A clean disjoint replay gets a new attributed candidate commit while retaining
  its checkpoint relationship. Patch overlap, conflict, ownership drift, or an
  invalid assumption leaves the candidate unchanged and reports that assignment
  as the next delegate action. A concurrent tip move may retry only while all
  those facts remain true.
- After a successful candidate compare-and-swap, `integrate` records the
  candidate/checkpoint relationship before it asks the existing worktree owner
  to release the now-landed assignment. Interruption between integration and
  release resumes cleanup without replaying the patch. A retained or failed
  release is visible as the next action; successful integration never leaks an
  ordinary completed assignment worktree.
- `review <slug> --evidence <receipt>` accepts a bounded semantic-review receipt
  for Standards, Spec, and Coverage over the exact candidate tip. The receipt
  records every finding's disposition. Accepted repairs are new ownership-fenced
  assignments and pass through checkpoint and integrate; no review command
  writes project files. Any candidate mutation makes the prior review stale, and
  a repaired composition receives a fresh three-axis review.
- `status` is the AXI query face of the family. Its default is one compact TOON
  row containing `slug`, `state`, `subject`, and `next`. `status <slug> --full`
  is the mandatory retained-evidence inspection projection: bounded assignment
  rows resolve ticket and source attribution, assignment base, checkpoint and
  integrated commits, receipt digest, and cleanup state; one review row resolves
  the reviewed candidate, three axes, finding dispositions, and review-receipt
  digest. It never leaks raw receipt bodies. Both views print a definitive
  terminal or empty result. Errors use structured stdout with exit 1, usage is
  exit 2, and no command prompts. Mutations return the same post-transition state
  projection, so their own write cannot make the reported next action stale.
- Each mutating subcommand is request-idempotent. Exact replays return the
  already-recorded result at exit 0; a reused identity with different inputs is
  a conflict at exit 1. State checkpoints surround every external mutation so
  re-entry can finish, roll forward, or report a recovery action without
  guessing whether an operation happened.
- One shared precondition owner runs before `start`, `assign`, `checkpoint`,
  `integrate`, `review`, `promote`, and `abandon --apply`. It distinguishes the
  expected dirty assignment worktree from the active working checkout and
  refuses before mutation on a dirty or unrecognized moved working branch, stale
  candidate/spec tip, missing prerequisite evidence, or ownership mismatch.
  A recognized clean working-branch advance is not silently accepted: only
  `promote` may begin the specified prospective recomposition, and every other
  mutator reports that recomposition as the next action without changing state.
- `abandon <slug>` is a read-only plan. It inventories active worktrees,
  provisional refs, unintegrated checkpoints, and recovery refs and returns a
  fingerprint. `abandon <slug> --apply <fingerprint>` revalidates the complete
  plan, preserves every Git-visible unlanded change through the existing
  worktree recovery owner, releases only verifiably owned worktrees, marks the
  run terminal, and retains the evidence history. Plan drift refuses without a
  partial cleanup.
- Extend `internal/gate`, rather than teaching `specbuild` the cache schema. The
  gate owner exposes subject-addressed inspection and execution for an
  unpublished prospective tree, retains exact green evidence by subject digest,
  and keeps `bench-last-gate` as the latest-run projection. The closed subject
  remains the candidate tree, resolved gate and launcher closure, declared tools
  and environment, execution-policy version, and freshness window. Reds are
  never reusable.
- `promote` first verifies a current three-axis review and a clean expected
  working base. It applies the canonical staged-to-implemented transformation to
  an unpublished prospective tree and asks the gate owner to run or reuse the
  whole-project gate for that exact subject. The gate executes through a
  checkout of the prospective tree so the real wrapper, binary resolution, and
  project gate observe the bytes being authorized. `promote` is the only
  `bench spec build` subcommand allowed to call that API.
- A red or operational gate failure publishes nothing: candidate ref, staged
  spec, working branch, and project-green ref remain unchanged. On green,
  promotion creates one squash commit whose parent is the expected working tip
  and whose tree is byte-identical to the authorized prospective tree. It
  advances the checked-out working branch through one checkout-aware,
  exact-base transition, then updates its project-green ref. The candidate and
  checkpoint refs remain retained evidence and never enter working ancestry.
- Promotion also requires every assignment to be integrated and released. After
  the branch and green-ref transition completes, it marks the run terminal while
  retaining the candidate, checkpoint, receipt, and review relationships exposed
  by `status --full`. Re-entry after terminal promotion is an exit-0 projection,
  not a second gate or cleanup pass.
- A crash after the working branch advances but before the green ref moves
  leaves the branch visibly unmarked and therefore not project-green. Re-entry
  asks the gate owner to validate the retained exact subject, then finishes the
  marker update; it never reruns or trusts a partial state merely because the
  new commit exists. Updating the marker first is forbidden.
- A working-branch advance first triggers recomposition onto the new base.
  `specbuild` constructs the result before mutation and applies the same
  unchanged-patch, ownership, and assumption rules as sibling replay. Conflict
  leaves both refs unchanged. A changed candidate tip invalidates its review.
  Gate evidence is reusable only when the complete closed subject is provably
  identical; commit identity or history alone neither invalidates nor authorizes
  a verdict.
- Candidate-owned gate reds become focused repair assignments on the same run.
  Inherited or infrastructure reds follow existing diagnosis and retry posture
  without manufacturing code work. After any repair integration, the whole
  composition is reviewed again and then `promote` pays one new gate. The
  implementation stage's declared iteration cap remains the stop; the lifecycle
  never falls back to per-ticket whole-project gates.
- This is a wide build. Implementation tickets use `craft-tickets`; ownership
  fences follow package ownership, with `internal/gate`, `internal/specbuild`,
  `internal/spec`, runtime/conformance contract packages, and workflow guidance
  assigned to distinct writers. A shared gate-evidence primitive is not split
  across its consumers. Every slice owns every file it must edit and no two
  active writers share a file.
- Update the following enumerated guidance and inventory artifacts in one
  consistent cut:
  - `.agents/commands/bench-implement-spec.md` owns the eight-command lifecycle
    order, frontier dispatch while other delegates remain active, focused
    evidence and coordinator probe, composed review, repair re-entry, and
    promotion; it forbids ticket `bench commit` and harness-authored Git
    lifecycle plumbing on this path.
  - `.agents/skills/bench-craft-delegate/SKILL.md` owns the spec-backed exception
    to per-ticket whole gates and generic release: the coordinator assembles the
    independent-probe receipt, checkpoints, integrates, and lets the lifecycle
    release the assignment; it forbids a delegate done-claim from becoming
    project-green evidence.
  - `.bench/BENCH.md`'s Workflow and command inventory own the
    spec-backed-versus-light-path split and all eight stable subcommands; they
    forbid extending provisional cadence to `bench shift` or ordinary
    `bench commit`.
  - `.bench/BENCH-reference.md` owns the complete compiled command inventory and
    routes the implementation adapter to the canonical phase.
  - `bin/bench.sh` help owns copy-paste grammar for all eight subcommands,
    including evidence flags and abandonment plan/apply.
  - `projects/benchkit.md` owns the spec-build line routing, gate attachment,
    hostile-input attachment, and dogfood expectations.
  - Conformance mutates every required anchor above independently and forbids the
    old per-ticket whole-gate, raw-Git integration, review-after-gate, and
    provisional-is-done claims.
  The harness must dispatch every ownership-safe frontier concurrently up to the
  smaller of frontier size and harness capacity and continue coordinating new
  eligible work while earlier delegates run. It names a dependency, overlapping
  fence, unavailable slot, or measured resource constraint for each unused slot.
  It must use the public command family, perform fresh composed review before
  promotion, and never claim provisional work is done.
- The cadence is limited to spec-backed builds with reviewed source, tickets,
  ownership fences, focused rows, and composed review. Light-path work,
  `bench shift`, and ordinary `bench commit` keep their current commit-on-green
  contract. The existing gate remains the only oracle in both paths.

## Testing decisions

- The primary seam is the exported `internal/specbuild` lifecycle driven against
  real temporary Git repositories with narrow fake gate and worktree
  collaborators. Tests invoke whole transitions and observe state projections,
  refs, commits, worktrees, and next actions rather than private record fields.
  The state machine must be deep enough that the CLI adapter only parses flags
  and renders results.
- The gate-evidence seam stays in `internal/gate`. Focused tests construct two
  subjects that vary one closure input at a time, exercise prospective-tree gate
  execution, and verify retention, freshness, red non-reuse, and latest-cache
  projection. Existing prior art is `internal/gate/gate.go`,
  `verdict_reuse_test.go`, `fault_engine_test.go`, and story 3/4 proof fixtures.
- Worktree creation, ownership, recovery refs, and release remain tested through
  the existing `internal/worktree` interface. `specbuild` contract tests use the
  real owner rather than duplicate its marker or cleanup fixtures. Prior art is
  `internal/worktree/ownership.go`, `lifecycle_test.go`, and
  `resume_test.go`.
- Black-box runtime contracts invoke the real `bench spec build` family from a
  fixture repository, including the repo-local linked wrapper. They grade TOON
  output, exit codes, idempotency, ref/ancestry results, and the exact recovery
  route. Conformance owns command registration, help/inventory currency, every
  enumerated guidance anchor, coverage-map validity, and the rule that no
  harness surface reconstructs the Git lifecycle.
- Fault tests interrupt each externally visible mutation phase: state prepare,
  worktree create, checkpoint commit/ref bind, candidate compare-and-swap,
  review record, prospective gate completion, squash creation, working-branch
  advance, green-ref update, and abandonment cleanup. Every fault is followed by
  `status` and the documented re-entry command, then compared with an
  uninterrupted result.
- One canary family plants mutations that bypass receipt validation, marks green
  before branch advancement, or lets a non-`promote` subcommand call the
  whole-project gate. Each mutation must make `bench gate` red with a targeted
  diagnostic; the ordinary contract suite supplies the broader behavior.
- The feature gate is `bench gate`. Focused package and runtime contracts precede
  it, while promotion tests use a deliberately small fixture gate and never run
  this repository's full gate recursively.

### Seam diagram: spec-build lifecycle

```text
trigger: bench spec build <start|assign|checkpoint|integrate|review|status|promote|abandon>
    |
    v
argv + spec + receipts --> [ internal/specbuild state machine ] --> TOON state + refs/commits/worktrees
                               |             |
                               |             +--> internal/worktree ownership and recovery
                               +----------------> internal/gate subject evidence and execution
                                  ^ tests attach here: real temp repo, controlled collaborators
```

### Seam diagram: subject-bound promotion

```text
trigger: promote after current semantic review
    |
    v
candidate + status flip --> [ prospective tree ] --> [ internal/gate ] --> retained red/green evidence
                                  |                       ^
                                  |                       tests vary one closed-subject input
                                  v
                           green only: squash child --> working branch --> refs/bench/green/<branch>
```

### Seam diagram: harness workflow and oracle attachment

```text
trigger: /bench-implement-spec
    |
    v
ticket graph --> concurrent assignments --> checkpoints --> serialized integration
                                                           |
                                                           v
                                                three-axis review + repairs
                                                           |
                                                           v
                                                promote --> bench gate
    tests attach: runtime command contracts + conformance guidance mutations + canary bites
```

### Acceptance coverage map

The command family is absent today. The observed entry red is
`bench spec build status spec-integration-gate-cadence`, which exits 2 with
`usage: bench spec (unknown argument: build)`. Rows below marked TDD-able start
with a failing focused case at the named seam before implementation. Existing
owners are named where coverage is already present; no provisional row is
misclassified as a whole-project green.

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `start` on a clean branch whose tip has retained exact green evidence creates one run, candidate ref, and branch-scoped green ref at that tip | spec-build lifecycle | TDD-able after the observed missing-family red: the start contract fails until the lifecycle exists. | An initializer that trusts HEAD without gate evidence or omits either ref cannot satisfy all three observed identities. |
| 1 | a missing green ref with no exact reusable evidence refuses before build state and names `bench gate` then retry as the only route | spec-build lifecycle | TDD-able: a fixture with no gate evidence must remain ref-identical and fails until bootstrap validation exists. | A trust-on-first-use initializer creates state and changes the ref inventory, failing the no-mutation assertion. |
| 1 | `start` never calls the whole-project gate and imports only evidence the gate owner validates | gate-evidence seam | TDD-able: a counting fake gate fails if start executes it or if an invalid record is accepted. | It distinguishes bootstrap from a hidden second gate caller. |
| 1 | repeating `start` resumes the same compatible run, while a different branch, detached HEAD, dirty tree, unexpected tip, or conflicting run refuses without a second run | spec-build lifecycle | TDD-able as a six-case table over one initialized fixture. | An always-create implementation leaks a second identity; a branch-blind resume accepts at least one conflict case. |
| 1 | project-green is true only while the working tip equals `refs/bench/green/<working-branch>` and retained evidence still matches the current subject | gate-evidence seam | TDD-able: independently move the branch, expire evidence, and vary the oracle in three subcases. | A ref-only marker passes the first state and fails all three invalidations. |
| 2 | `assign` resolves a real ticket, records its row set, fence, and assumptions, and creates one owned worktree at the current candidate | spec-build lifecycle plus worktree owner | TDD-able: the real-worktree fixture fails until spec-scoped assignment state is bound. | A generic worktree create that ignores the ticket cannot report or later validate the recorded contract. |
| 2 | two independent sibling assignments may share the same candidate base and occupy separate owned worktrees | spec-build lifecycle plus worktree owner | TDD-able: issue two request ids before either integrates and assert equal bases plus unequal paths. | A serial-only allocator advances or blocks after the first assignment and fails the sibling assertion. |
| 2 | reusing a request for the same run and ticket is idempotent, while reusing it for another ticket or run refuses | spec-build lifecycle | TDD-able as one replay and two conflict cases. | Request hashing without scope accepts one of the conflicts or creates duplicate worktrees. |
| 2 | a missing, outside-spec, malformed, special-file, or symlinked ticket refuses before worktree creation | spec-build lifecycle | TDD-able per input class with the worktree fake asserting zero calls. | It proves ticket classification precedes side effects and prevents a path-shaped label from becoming authority. |
| 3 | a receipt binding every charged row, required seam checks, independent probe, exact resulting tree, and ownership inventory creates one attributed checkpoint without moving the candidate or running the gate | spec-build lifecycle | TDD-able: the valid receipt case fails until checkpoint exists; counting fakes pin zero candidate and gate mutations. | An implementation that silently integrates or gates the checkpoint changes one of the protected counters. |
| 3 | the coordinator probe section is produced outside the assignment worktree by a producer distinct from the delegate and is fresh for this exact assignment tree | spec-build lifecycle | TDD-able: delegate-produced, stale, other-assignment, other-tree, and inside-worktree probe sections each refuse before checkpoint mutation. | A delegate-substituted `probe=pass` receipt satisfies field presence but fails provenance, location, freshness, and subject binding. |
| 3 | already-covered and reviewed not-TDD-able row classifications remain admissible without fabricated red-to-green output | spec-build lifecycle | TDD-able: one receipt of each classification must checkpoint successfully. | A red-only validator rejects honest rows and encourages false evidence. |
| 3 | omission or failure of any charged row, seam check, or independent probe refuses and leaves state, refs, and worktree bytes unchanged | spec-build lifecycle | TDD-able as per-field mutation cases from one valid receipt. | An always-green or presence-only validator accepts at least one missing or failed result. |
| 3 | a changed assignment base, mismatched tree, altered ticket digest, outside-fence patch, unexplained path, or invalid assumption refuses before commit creation | spec-build lifecycle | TDD-able as six single-defect fixtures. | A receipt parser that does not re-read live Git facts passes signed-looking but stale evidence. |
| 3 | an empty, oversized, malformed, unreadable, FIFO, device, socket, regular symlink, or dangling-symlink receipt is rejected before reading unbounded or blocking input | spec-build lifecycle | TDD-able per hostile file class, with the FIFO case required to return without a writer. | It proves metadata and bounds checks happen before content consumption. |
| 4 | `integrate` consumes only a verified checkpoint and advances the exact candidate tip with one attributed commit | spec-build lifecycle | TDD-able: unverified and verified cases fail until checkpoint state and compare-and-swap are enforced. | A raw-commit or receipt-only integration accepts the negative control. |
| 4 | a later disjoint sibling with an older base replays its byte-identical complete patch onto the current candidate | spec-build lifecycle | TDD-able: two siblings edit separate owned paths and integrate in reverse completion order. | A strict-base-only implementation rejects the second sibling; a stale overwrite loses the first sibling's bytes. |
| 4 | overlap, patch drift, conflict, ownership drift, or changed assumptions leave the candidate unchanged and route the same assignment back to its delegate | spec-build lifecycle | TDD-able as five prospective-replay failures with before-and-after ref equality. | A best-effort replay mutates at least one failing case or loses attribution. |
| 4 | a concurrent candidate move retries only when the unchanged patch, fence, and assumptions still hold, using exact-old-tip compare-and-swap | spec-build lifecycle | TDD-able with a hook that moves the candidate immediately before the first update. | An unconditional retry can apply stale work; a non-CAS write can discard the injected move. |
| 4 | successful integration records the checkpoint relationship then releases the landed assignment; interruption resumes release without applying the patch twice | spec-build lifecycle plus worktree owner | TDD-able with a fault after candidate compare-and-swap and before release, followed by integrate re-entry. | An integration-only implementation leaks the worktree, while a transition without a recorded checkpoint relationship duplicates the commit on replay. |
| 5 | `review` accepts Standards, Spec, and Coverage evidence with every finding disposition only for the exact candidate tip | spec-build lifecycle | TDD-able: omit each axis and disposition, then change the candidate tip before submission. | A prose-presence check or branch-name binding accepts an incomplete or stale review. |
| 5 | integrating any implementation or repair checkpoint invalidates prior review evidence | spec-build lifecycle | TDD-able: record review, integrate one disjoint checkpoint, and assert `next=review`. | A sticky reviewed flag would allow the changed composition to enter the gate. |
| 5 | accepted repairs use ordinary assign, checkpoint, and integrate transitions and the repaired composition requires a fresh full review | harness workflow seam | TDD-able through a runtime fixture whose first review receipt accepts one finding. | A special unreceipted repair path or delta-only review skips one required transition. |
| 6 | `status` renders one definitive TOON row with `slug`, `state`, `subject`, and `next` from durable state, including terminal state | runtime CLI | TDD-able after the observed missing-family red: exact shape and exit tests fail today. | A silent or prose-only status cannot supply the next command to a fresh agent. |
| 6 | re-entering every interrupted mutation reports or performs the one durable next action without duplicating a run, worktree, commit, or ref | spec-build lifecycle | TDD-able with one fault point per external mutation followed by status and replay. | A state record written only at the end cannot distinguish whether the side effect already occurred. |
| 6 | `abandon` plans without mutation and apply succeeds only for the same fingerprint, preserving every unlanded Git-visible change before owned cleanup | spec-build lifecycle plus worktree owner | TDD-able: plan/apply drift and a dirty assignment fail until recovery and fingerprinting are wired. | Direct cleanup loses the dirty fixture; a stale plan mutates after an injected ref change. |
| 6 | exact repeats are exit-0 no-ops and reused identities with different inputs are exit-1 conflicts; syntax errors exit 2 and no command prompts | runtime CLI | TDD-able as a family-wide command table over all mutators. | A command-specific parser or interactive fallback fails at least one shared posture. |
| 6 | `start`, `assign`, `checkpoint`, `integrate`, `review`, `promote`, and `abandon --apply` each refuse before state, ref, commit, or worktree mutation when the active working checkout is dirty or unrecognized-moved, the candidate/spec tip is stale, prerequisite evidence is missing, or ownership mismatches | spec-build lifecycle | TDD-able as a seven-command by four-precondition matrix; expected dirt in the named assignment worktree remains a positive control. | Protecting only start and promote leaves an intermediate mutator able to advance durable state after the working subject stopped being authoritative. |
| 7 | `promote` constructs the canonical implemented prospective tree while the candidate and working spec remain staged | subject-bound promotion | TDD-able: inspect all three trees before the gate fake returns. | An in-place flip dirties or mutates a published tree before authorization. |
| 7 | the real gate wrapper executes against the prospective tree and its complete launcher closure, with `promote` the only build subcommand that can call it | gate-evidence seam plus runtime CLI | TDD-able with a fixture gate that records cwd and sentinel bytes plus per-subcommand call counters. | A gate over the candidate misses the status flip; a hidden caller increments another counter. |
| 7 | a red or operational gate outcome changes no candidate, working branch, staged spec, or project-green ref | subject-bound promotion | TDD-able with red, timeout, and persistence-fault cases compared byte for byte. | Any publish-before-green ordering changes at least one protected identity. |
| 7 | green promotion creates one squash child of the expected working tip whose tree equals the authorized prospective tree, excluding provisional commits from working ancestry while retaining their refs | subject-bound promotion | TDD-able with two checkpoint commits and ancestry plus tree assertions. | A merge-based landing exposes provisional ancestry; a reconstructed squash can drift from the gated tree. |
| 7 | after green promotion `status --full` resolves each retained checkpoint to its ticket/source attribution, assignment base, checkpoint and integrated commits, receipt digest, cleanup state, and the retained review subject, axes, dispositions, and receipt digest | runtime CLI plus spec-build lifecycle | TDD-able with two checkpoints and one review inspected before and after the squash. | Retaining refs while deleting or orphaning receipt and review relationships passes ancestry checks but fails the supported inspection projection. |
| 7 | promotion advances the checked-out working branch before its green ref; an injected crash between them fails closed and resumes from retained exact evidence | subject-bound promotion | TDD-able with faults on both sides of branch advancement and a working-tree cleanliness assertion. | Marker-first leaves an ungated branch labeled green, while no recovery leaves the authorized commit permanently ambiguous. |
| 7 | promotion refuses while any assignment is unintegrated or unreleased, then marks the fully promoted run terminal without deleting its evidence | subject-bound promotion | TDD-able with one active, one integrated-but-release-pending, and one fully released assignment state. | A tree-only promoter can strand owned worktrees or discard lifecycle evidence after the branch moves. |
| 8 | changing candidate tree, oracle or launcher closure, declared tools or environment, execution policy, or freshness invalidates green evidence; changing commit history alone with an identical subject does not | gate-evidence seam | TDD-able as one mutation per enumerated subject component plus one history-only positive control. | A tree-only key accepts five unsafe drifts, while a commit-keyed cache rejects the positive control. |
| 8 | a working-branch advance recomposes the unchanged candidate patch before mutation; conflicts leave both refs unchanged and a changed candidate requires review again | spec-build lifecycle | TDD-able with clean-disjoint and conflict base advances. | A stale-base promotion overwrites branch work; a silent rebase reuses stale review. |
| 8 | a candidate-owned red creates focused repair next actions, while inherited and infrastructure reds produce diagnosis or retry actions without code repair | spec-build lifecycle | TDD-able with three classified gate receipts over the same candidate. | A generic red-to-ticket rule manufactures code changes for machine or inherited failures. |
| 8 | after repair integration the whole composition is reviewed again and one new promote gate runs; no ticket transition runs a whole-project gate | harness workflow seam | TDD-able with a counting fixture spanning one red, one repair, and final green. | Reverting to per-ticket gates or skipping re-review changes the expected call sequence. |
| 8 | exhausting the declared implementation-stage cap leaves the run provisional and reports the stop instead of retrying or escalating silently | harness workflow seam | TDD-able through a bounded orchestrator fixture at the cap. | An unbounded retry loop performs one extra transition and fails the terminal next-action assertion. |
| 9 | every ownership-safe frontier dispatches concurrently up to the smaller of frontier size and harness capacity | guidance conformance | TDD-able with a static guidance mutation and a three-ticket dogfood trace. | Advisory wording permits the exact sequential cadence the decision rejects. |
| 9 | as a running delegate frees a slot or integration unlocks another ownership-safe ticket, the coordinator dispatches that eligible work while other delegates remain active | guidance conformance | TDD-able with a four-ticket trace whose second frontier opens before the first frontier drains. | A coordinator that fan-outs once and then waits satisfies initial concurrency while leaving later capacity idle. |
| 9 | each unused slot names a dependency, overlapping ownership fence, unavailable harness slot, or measured resource constraint | guidance conformance | TDD-able as four allowed reasons plus one generic-reason rejection. | A vague `not parallelizable` explanation passes advisory prose but fails the closed reason set. |
| 9 | the harness invokes the eight public subcommands and never creates commits, edits lifecycle refs, or composes worktree Git plumbing itself | guidance conformance | TDD-able with command-inventory anchors and forbidden plumbing mutations. | Naming start and promote alone still lets intermediate transitions drift in prose. |
| 9 | exact composed Standards, Spec, and Coverage review precedes promote, and accepted repairs re-enter the same provisional path | guidance conformance | TDD-able by deleting or reordering each required phase in a command fixture. | A review-after-gate or delta-only repair review pays or trusts the wrong subject. |
| 9 | `.agents/commands/bench-implement-spec.md` independently pins lifecycle order, continued frontier dispatch, evidence plus coordinator probe, composed review, repair re-entry, promotion, and the no-ticket-gate/no-raw-Git forbids | guidance conformance | TDD-able with one mutation per enumerated required or forbidden anchor in that file. | Updating only the command list leaves the old cadence or a harness-plumbing escape alive in the phase that executes it. |
| 9 | `.agents/skills/bench-craft-delegate/SKILL.md` independently pins the spec-backed checkpoint/integrate/release exception, coordinator-owned probe receipt, and provisional-not-green forbid | guidance conformance | TDD-able with one mutation per enumerated required or forbidden anchor in that file. | Generic delegate guidance can otherwise keep requiring `bench commit` and manual release for every ticket. |
| 9 | `.bench/BENCH.md` independently pins all eight subcommands and the spec-backed-versus-light-path split, while `.bench/BENCH-reference.md` independently carries the complete compiled command inventory and implementation adapter route | guidance conformance | TDD-able by deleting each command and each split/route anchor from otherwise clean copies. | One updated inventory can mask a stale shared workflow or reference surface. |
| 9 | `bin/bench.sh` help independently exposes copy-paste grammar for all eight subcommands, evidence flags, and abandon plan/apply | guidance conformance | TDD-able by deleting each grammar entry or moving a flag value into the positional slot. | A dispatcher can work while agents are taught an incomplete or unsafe invocation. |
| 9 | `projects/benchkit.md` independently pins spec-build routing, gate attachment, hostile-input attachment, and behavior-changing dogfood | guidance conformance | TDD-able with one mutation per profile anchor. | A generic workflow check cannot catch project-specific line or gate drift. |
| 9 | light-path changes, `bench shift`, and ordinary `bench commit` retain commit-on-green behavior | existing shift and commit runtime contracts | Already covered by shift/commit gate-call contracts; the new dogfood trace runs them unchanged as a regression control. | It prevents a broad router change from applying provisional semantics outside reviewed spec builds. |
| edge of 1 | a slug or ticket path containing spaces or glob characters resolves literally and produces opaque safe ref identities | runtime CLI plus spec-build lifecycle | TDD-able with quoted arguments whose neighboring similarly prefixed spec must remain untouched. | Shell interpolation or glob expansion resolves the wrong spec or creates an invalid ref. |
| edge of 2 | each value for `--ticket`, `--request`, `--assignment`, `--evidence`, and `--apply` is consumed only by its flag; the same tokens after `--` are literal path text and never become the slug or subcommand | runtime CLI | TDD-able as one value-only and one `--` case per flag, with no explicit positional decoy. | A first-non-dash parser passes ordinary invocations but treats a flag value as the positional the guard meant to protect. |
| edge of 3 | a receipt whose final line has no newline is either accepted as complete JSON or rejected with one stable malformed-record error, never truncated or silently empty | spec-build lifecycle | TDD-able with identical bytes differing only by the final newline; implementation must choose and pin one result. | It closes the hand-edited-record ambiguity without inventing a partial successful parse. |
| edge of 6 | control bytes and sink-permitted tab, newline, or return in display text never split TOON rows or diagnostics | runtime CLI | TDD-able with Git-sourced subjects and ticket labels carrying every representability class Git permits. | Sanitizing only ESC and BEL still lets a permitted newline corrupt the line-oriented projection. |
| edge of 6 | invocation from a nested repository directory and through the real kit and linked-repo wrappers resolves the same common-directory run | runtime CLI | TDD-able as three invocation surfaces over one run. | A cwd-local or worktree-Git-dir store forks lifecycle state by caller. |
| edge of 6 | SIGINT during gate, checkpoint, integration, promotion, or abandonment leaves recoverable state and no surviving child process | spec-build lifecycle plus gate runner | TDD-able with one signal case per long-running phase and process-group liveness probes. | Context-only cancellation can leave a child mutating after the state machine reports interruption. |
| edge of 6 | a missing `git`, unresolved gate, or missing real binary fails with an actionable next step and no state mutation beyond a recoverable prepared record | runtime CLI | TDD-able with PATH and launcher fixtures at each resolution point. | Dependency stack traces or partial refs leave an agent unable to resume safely. |
| edge of 7 | the prospective status transformation is deterministic and a second construction produces the same tree object | subject-bound promotion | TDD-able by building it twice from one candidate and comparing tree ids. | A time- or worktree-dependent flip cannot be authorized once and published byte-identically. |

Degenerate implementations are pinned per story. Trusting `HEAD` as green (1)
fails the absent-evidence row. A serial assignment allocator (2) fails the equal
sibling-base case. A receipt-presence stub or delegate-authored `probe=pass`
(3) fails every single-field mutation, probe-origin, and outside-fence case. A
strict-base-only or force-update integrator (4) fails opposite halves of the
sibling replay cases, while a no-cleanup integrator fails release re-entry. A
sticky `reviewed=true` bit (5) fails candidate invalidation. End-only persistence,
start-only preconditions, and direct deletion (6) fail the interruption,
all-mutator, and dirty-abandonment rows. An in-place status flip, merge
promotion, or refs-only evidence retention (7) fails the unpublished-tree,
ancestry, and full-inspection assertions. A tree-only evidence key or generic red
repair (8) fails the enumerated drift and classification cases. Advisory
parallelism, one-shot fan-out, or partial command documentation (9) fails the
per-artifact conformance and dogfood traces.

### Edge inventory

- Error path — resolved by the bootstrap refusal, receipt mutations, replay
  failures, red classifications, tool-resolution failures, and structured CLI
  errors.
- Empty or absent input — resolved by missing evidence, missing tickets, empty
  receipts, absent run status, and missing branch-scoped green refs.
- Boundary values — resolved by receipt size limits, full frontier capacity,
  exact freshness expiry, stage-cap exhaustion, and final-newline cases.
- Malformed input — resolved by malformed tickets, receipt schemas, flags,
  identities, TOON-hostile display text, and invalid ref sources.
- Interrupted or partial state — resolved by every lifecycle fault point, the
  branch-before-marker crash, partial matrices of refs/state/worktrees, and
  re-entry rows.
- Re-run idempotency — resolved by start, request reuse, every mutation replay,
  status, and abandon plan/apply rows.
- Hostile environment — resolved by spaces/globs, special files, symlinks,
  nested cwd, wrapper parity, missing tools, signals, and dirty or concurrently
  moving Git state.
- A command whose write changes a fact it reports — resolved by rendering the
  post-transition durable state from the same result each mutation records.
- Flag values mistaken for positionals — resolved by the family-wide usage table;
  `--ticket`, `--request`, `--assignment`, `--evidence`, and `--apply` consume
  exactly one value and the concrete flag-value matrix proves `--` handling.
- Destructive worktree state — resolved by ownership validation, unexplained-path
  refusal, recovery-before-release, and fingerprinted abandonment.
- Non-TTY stdin — **Won't handle**: every input is an argument or bounded file and
  no command prompts, so stdin mode cannot amputate an in-scope caller.
- Invocation through an arbitrary symlink outside the shipped kit or linked-repo
  wrappers — **Won't handle**: the existing binary resolver owns arbitrary
  symlink discovery; this feature proves both supported surfaces and their
  missing-binary posture.
- Invocation through hooks or model adapters — **Won't handle**: those shipped
  surfaces do not expose or invoke `bench spec build`; adding a lifecycle caller
  there would create a second coordinator rather than preserve a real calling
  convention. The real kit CLI and linked-repo by-path CLI are the complete
  in-scope invocation set for this family.
- Host-backed filesystem I/O pressure — **Won't handle**: lifecycle durability
  requires file and directory synchronization, and timing out after the kernel
  may have accepted a write makes completion ambiguous. Re-entry and state
  inspection, not an unsafe fsync timeout, are the recovery contract.

## Out of scope

- Applying provisional integration cadence to light-path work or `bench shift`
  is a separate workflow capability gated on measured evidence from this one —
  12 edits, 3 gate runs.
- General-purpose `bench worktree land` outside spec builds remains FT169's
  distinct lifecycle capability; this spec uses the worktree owner but does not
  replace that public surface — 10 edits, 3 gate runs.
- Build-economics telemetry, wall-clock comparison, gate-count history, and
  retained per-spec metrics remain FT138's instrumentation capability — 10
  edits, 3 gate runs.
- Bounding outer gate-phase concurrency remains FT171's measured scheduling
  decision and does not change which subject or cadence the gate proves — 9
  edits, 3 gate runs.
- Ship-tier release verification and publication authority remain separate from
  dev project-green promotion — 6 edits, 2 gate runs.
- Evidence-retention expiry, pruning policy, export, or a reviewer-facing history
  UI is a separate evidence-management capability; this build retains evidence
  without automatic deletion — 8 edits, 2 gate runs.
- Diff-scoped gating, per-package verdict reuse, file-to-test maps, weakened or
  skipped checks, per-ticket whole-project gates, and provisional commits in
  working-branch ancestry are rejected by the reviewed source, not deferred
  portions of this capability.
