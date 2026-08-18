# Spec-build review and gate cadence

Status: shaping

## Destination

A successor spec-build workflow in which each ownership-fenced ticket earns
focused executable evidence and a durable orchestrator review before provisional
integration, the exact composed candidate is pinned for one fresh semantic review,
and promotion alone pays the whole-project gate and authors the landing. The
successor consumes FT173's full AXI lifecycle contract and FT130's capture
accounting contract after both have shipped; it does not reopen the terminal
single-build serial-gate work.

## #1: What evidence admits a ticket into the composed candidate?

Blocked by: none
Type: Grill

### Question

Should ticket-local Standards, Spec, and Coverage results remain harness-only,
or become durable checkpoint evidence and lifecycle preconditions?

### Answer

Each assignment worktree and its assignment branch are one provisional
construct: neither runs the whole-project gate nor authors an ordinary commit.
Before checkpoint, the delegate supplies focused executable evidence and its
mutation probe; the orchestrator independently probes and reviews the exact
assignment tree on separate Standards, Spec, and Coverage axes against the
ticket's requirements, ownership fence, coverage rows, and declared integration
surfaces.

The checkpoint receipt durably binds those axis results and every finding's
reviewer disposition, alongside the focused checks and independent probe, to the
run, assignment, ticket digest, base, and tree. The lifecycle validates presence,
shape, exact binding, and that no repair-required finding remains. It does not
judge whether the semantic review is correct or treat checkpointing as project
green. Only a verified checkpoint may integrate provisionally.

## #2: What may ticket-local checks compile and execute?

Blocked by: none
Type: Grill

### Question

Which focused checks may use direct package tests, and which ticket behaviors
genuinely require an exact Bench executable?

### Answer

Ticket scope is defined by its acceptance rows and declared integration surfaces,
not by a file-only compilation boundary. A ticket with no Go change and no Go
process seam compiles nothing. Ordinary Go behavior uses direct package or unit
tests, which may compile and run Go's package-specific test executable and the
unchanged dependencies required to exercise the declared seam.

An exact Bench executable is required only when an acceptance row observes an
operating-system process seam: wrapper routing, executable identity or freshness,
environment propagation, signals or process teardown, or installed or stripped
behavior. No ticket-local check runs the whole-project gate.

The promotion performance contract remains one selected Bench CLI build per
non-reused promotion gate and one serial phase schedule, not one compiled binary
total. Independent package-test, race, system, changed-source, cross-target, and
other compiler-observing proofs may require their own artifacts.

## #3: How are ticket-local conformance and canary checks selected?

Blocked by: #2
Type: Grill

### Question

Should changed paths alone select conformance and canary checks, or should their
registered owners and input derivations participate?

### Answer

Ticket-local conformance runs only the exact checks that own the ticket's changed
inputs or declared integration surfaces. Ticket-local canary runs only the exact
affected fixture-owner mutations. A canary is affected when the ticket changes
its fixture, registered owner, routing, or input derivation; an artificial fixture
edit is not required to trigger an owner check. Unrelated conformance and canary
inventory stays out of ticket evidence and remains covered by promotion's complete
gate.

## #4: What transition pins the composed candidate for final review?

Blocked by: #1
Type: Grill

### Question

How does the lifecycle prevent candidate movement between the start of semantic
review and the promotion gate?

### Answer

Final review is a two-step lifecycle transition. `bench spec build review <slug>
--begin` atomically pins the current candidate only after every assignment is
integrated and released, and returns the exact run and candidate review subject.
Candidate-changing assign, checkpoint, integrate, refresh, and recomposition
operations refuse while that subject is pinned. The orchestrator submits the exact-candidate
three-axis result through `review --evidence`.

A repair-required result releases the candidate into ownership-fenced repair
assignment, checkpoint, and integration. A clean or explicitly risk-accepted
result stays pinned through promotion. If the working branch moves, promotion
recomposes before any gate begins, invalidates the pin and receipt, and returns
the entire changed composition for a fresh review.

## #5: Which harness phase owns final review and promotion?

Blocked by: #4
Type: Grill

### Question

Should the final wrapper remain inside `$bench-implement-spec --full`, or should
Bench add a separate finalize-spec phase?

### Answer

`$bench-implement-spec --full` remains the end-to-end wrapper and composes the
standalone `$bench-review-implementation` phase immediately before promotion.
There is no separate finalize-spec phase. A clean or risk-accepted exact-candidate
receipt proceeds to `bench spec build promote`; an accepted defect pays no gate,
becomes a fenced repair ticket, and returns through the ordinary provisional
lifecycle before the complete changed composition is pinned and reviewed again.
That composed review specifically hunts cross-ticket knowledge duplication, dual
ownership, unresolved integration surfaces, incomplete composition, and coverage
gaps that no ticket-local review can see.

`promote` alone constructs the prospective implemented tree, runs the sole
whole-project spec-build gate, and authors the landing and implemented state.
Ticket-local and composed semantic reviews remain advisory and never substitute
for that deterministic authority.

## #6: What does review, final-check, commit, and promote each own?

Blocked by: #5
Type: Grill

### Question

Which workflow surface spends model tokens, which spends whole-project compute,
and which one authors a landing?

### Answer

`$bench-review-implementation` asks whether the exact diff is good on Standards,
Spec, and Coverage. It is the token-intensive semantic surface, runs no
whole-project gate, and for a spec build submits only the candidate-bound review
receipt. `bench commit` is the deterministic ordinary and light-path gate-and-
landing owner. `bench spec build promote` is the deterministic spec-build gate-
and-landing owner. Those two commands spend whole-project compute and no model
tokens of their own.

`$bench-final-check` is uniformly post-landing and report/capture-only. It reports
retained evidence and records required closure capture; it never runs a gate or
authors a commit. Ordinary work therefore reviews, lands through `bench commit`,
then reports through final-check. Spec builds review, land through `promote`, then
report through final-check.

## #7: Which owner supplies AXI for the lifecycle family?

Blocked by: none
Type: Grill

### Question

Should full AXI apply only to the new review transition, or to the complete spec-
build lifecycle family, and which roadmap item owns it?

### Answer

FT173 is a separate required predecessor and the only owner of all ten AXI
principles for every `bench spec build` operation. In particular, lifecycle
results and refusals use the shared structured stdout and error helpers, honest
exit codes, idempotent retries, and contextual `help[]` actions populated with
the known slug, run, candidate, assignment, and exact next command. The shared
AXI helpers and schema are not re-derived in the cadence successor.

FT185 continues to own the reusable gate-phase result payload consumed by gate,
ordinary commit, and promotion. FT173 composes that payload into the lifecycle
surface instead of defining a second gate-result schema.

## #8: Which owner accounts for capture after landing?

Blocked by: #6
Type: Grill

### Question

Should post-promotion retro and learning capture live in the cadence successor,
or join the existing capture-lifecycle roadmap owner?

### Answer

FT130 is a separate required predecessor and the one owner of capture queuing,
pending-capture visibility, and eventual reviewed drainage. Report-only final-
check hands post-landing retros and learnings to that contract without paying an
immediate follow-up whole-project gate or misrepresenting project-green
authority. The cadence successor names the handoff and does not implement a
second capture classification or landing path.

## #9: How is the revised cadence accounted for?

Blocked by: #5, #7, #8
Type: Grill

### Question

Does the revised workflow amend the single-build serial-gate spec, require only
guidance edits, or receive a successor spec?

### Answer

The revised cadence receives a successor spec after FT173 and FT130 ship. It
changes durable checkpoint evidence, lifecycle preconditions and transitions,
agent-facing CLI behavior, and final-check responsibility, so guidance alone is
insufficient. The terminal single-build serial-gate spec and lifecycle remain
closed; their selected-executable and serial-schedule performance contract is an
input to the successor, not an amendment target.

## #10: Are the predecessor contracts present for spec authoring?

Blocked by: #7, #8, #9
Type: Task

### Question

Verify that FT173 and FT130 have shipped their decided AXI and capture contracts,
then re-read the lifecycle, workflow guidance, and project profile against this
map before handing it to spec authoring.

### Answer

— (open)

## Not yet specified

## Spec-writer discretion

- Internal type placement and naming for the pinned review generation, provided
  the public `review --begin` and `review --evidence` behavior and invalidation
  predicates remain exact.
- Exact focused Go test filters and package groupings within each ticket's
  acceptance rows and declared integration surfaces, provided they do not widen
  into a whole-project gate or omit an affected registered owner.
- Receipt field and TOON table layout only where FT173 leaves reversible freedom;
  the cadence spec must consume, not fork, FT173's shared AXI contract.

## Out of scope

- Reopening, repairing, or amending the terminal single-build serial-gate run or
  its spec.
- A whole-project gate or ordinary commit in an assignment worktree.
- Treating semantic review as deterministic gate authority or project-green
  evidence.
- Replacing FT173's AXI foundation, FT185's gate-result payload, or FT130's
  capture-accounting contract inside the successor.
- Durable cross-run Bench executable storage, later-process executable reuse,
  release-tier proofs, or a broader gate-scheduler rewrite.

## Sources

- Path: `ROADMAP.md`
  Supports: #7 and #8 existing FT173 and FT130 ownership.
  Drift: re-read after either roadmap item is reshaped, staged, or retired.
- Path: `CHANGELOG.md`
  Supports: #1's checkpoint receipt and #4/#5's candidate-bound review receipt, whose provisional lifecycle was removed wholesale and now survives only as this removal record.
  Drift: re-read before specifying the successor; its seams re-derive from the surviving landing path, not the deleted lifecycle.
- Path: `internal/landing/landing.go`
  Supports: #4 and #5's surviving composition-and-publication operation, now owned by `bench worktree land`.
  Drift: re-read if composition, gating, or publication ordering changes.
- Path: `projects/benchkit.md`
  Supports: #2 and #3 current process, conformance, and canary ownership.
  Drift: re-read if the gate architecture, process-seam inventory, or conformance registry contract changes.
- Path: `.agents/commands/bench-final-check.md`
  Supports: #6 current workflow-dependent final-check behavior.
  Drift: re-read if final-check or either landing route changes before spec authoring.
