# worktree-cleanup-eligibility

Status: implemented

Decision source: specs/worktree-cleanup-eligibility/decisions/deepening-2026-08.md (compiled ready map, refreshed and reviewer-resolved 2026-08-17; decisions #1 and #2 govern FT216)

Verification log: spec + tickets 2 iteration(s) to accept — one independent `gpt-5.6-terra` high-effort review round; iteration 1 found the unfenced shared `assignmentLanded` string reader and an over-bundled automatic/landed-set migration, both folded; iteration 2 returned ACCEPT with all 33 rows singly owned by the eight-ticket DAG. The required multi-iteration learning is recorded in `capture/learnings.md`.

Reviewer approval: 2026-08-17 — approved the behavior-preserving 33-row spec, the eight-ticket DAG, serial execution of tickets 04 and 05 on the retained integration source, and re-homing FT217/FT218's remaining decisions before FT216 retirement.

## Problem

Worktree cleanup does not have one owner for the answer “is this worktree ours
and safe to remove?” `PlanExplicitWithOptions` gathers ownership, assignment,
lock, lease, tracked-state, nested-state, ignored-residue, landedness, and
recovery facts while repeatedly overwriting `CleanupPlan.Action`, `ReasonCode`,
and `Reason`. The last applicable write is today's effective explicit
precedence. `PlanAutomatic` then reads the mutable plan, assignment fields, and
the formatted `plan.landed` string and makes a second ordered decision. The
`--landed` planner adds another preservation refusal after explicit planning.
`ApplyExplicit`, `ReleaseCommand`, and the landing release path consume those
results without one typed decision they can name.

The behavior is safety-critical and intentional until proved otherwise, but
its precedence is not stated as data and is not characterized as a complete
set of reachable `(Action, ReasonCode)` outcomes. A mechanical extraction can
therefore change which refusal wins while leaving many focused tests green.

## Solution

Add one in-process eligibility module in `internal/worktree` that receives typed
ownership and safety evidence, applies the existing ordered rules, and returns
one eligibility verdict carrying the decided action, reason code, operator
detail, and the evidence later cleanup stages need. Evidence collection may
remain close to Git and filesystem adapters, but no consumer reconstructs the
decision from fields or formatted text.

`PlanExplicitWithOptions` projects the verdict into the existing
operator-facing `CleanupPlan`. `PlanAutomatic` asks the same module for the
stricter automatic reading. The `--landed` set planner, explicit apply,
automatic apply, release, and landing release paths execute or report that
verdict without assigning a new eligibility action or reason. Landedness is a
typed fact, never a string-prefix protocol. `DiscardBranch` remains an
operator assertion applied only after derived landedness and is never visible
to automatic cleanup.

Before production logic moves, characterization tests pin every reachable
eligibility outcome and the collisions that make the present ordering visible.
No row authorizes a behavior correction. If the matrix exposes an anomaly, the
implementation stops that delta and routes it through `/bench-debug`.

## User stories

1. **Characterize explicit eligibility outcomes.** As the reviewer I can read
   one table-driven test that reaches each explicit eligibility tuple exactly
   once: retain with `uncertain`, `foreign`, `malformed`, `unexpected-lock`,
   `live-lease`, or `ignored`, plus remove, recover-remove, and discard-remove
   with no refusal reason. Each case asserts the action, reason code, relevant
   evidence, and no planning mutation. Cases with competing evidence pin the
   current last-write-wins result rather than an idealized precedence.
   Line: `gpt-5.6-terra` / high / ~3 iterations / serial.
2. **Characterize automatic eligibility outcomes.** As the reviewer I can read
   one table-driven test that reaches each automatic eligibility tuple exactly
   once: retain with `uncertain`, `foreign`, `malformed`, `unexpected-lock`,
   `live-lease`, `ignored`, `active`, `landed`, `orphaned`, `unmerged`, or
   `dirty`, plus remove and discard-remove with no refusal reason. The matrix
   pins the automatic overrides of explicit results and proves automatic
   cleanup cannot observe `DiscardBranch`.
   Line: `gpt-5.6-terra` / high / ~3 iterations / serial after story 1.
3. **Decide eligibility once from typed evidence.** As a maintainer I find one
   eligibility module that owns the ordered ownership, assignment, lock, lease,
   landedness, recovery, tracked-state, nested-state, and ignored-residue rules
   and returns one verdict plus evidence. `PlanExplicitWithOptions` gathers
   facts and projects that verdict into the unchanged cleanup-plan/output
   contract. No later explicit-planning statement can overwrite the result.
   Line: `gpt-5.6-terra` / high / ~3 iterations / serial after stories 1–2.
4. **Read the verdict at every cleanup path.** As a maintainer I see
   `ApplyExplicit`, `PlanAutomatic`, the `--landed` set planner,
   `ReleaseCommand`, and `LandCommand` consume the eligibility decision. The
   automatic path applies its stricter policy through typed evidence and has no
   `HasPrefix` or other parsing of landedness text. The apply paths retain the
   exact fingerprint, interruption, recovery, receipt, and branch-deletion
   behavior the current tests describe.
   Line: `gpt-5.6-terra` / high / ~3 iterations / serial after story 3.
5. **Leave one source for eligibility policy.** As a maintainer I can delete the
   old mutable and string-derived decision paths, find no second assignment of
   an eligibility action or refusal outside the eligibility module, and read
   ADR 0005 naming the eligibility verdict while retaining its conjunctive
   ownership and preservation posture. Every pre-existing test passes with test
   logic unchanged.
   Line: `gpt-5.6-terra` / high / ~2 iterations / serial after story 4.

## Implementation decisions

- **One verdict, not a renamed plan.** The verdict is the decided answer to
  “ours and safe to remove” plus the evidence required to execute or explain
  it. `CleanupPlan` remains the command projection and fingerprint carrier. A
  wrapper built only after the current mutable plan has already decided does
  not satisfy this spec.
- **Typed evidence.** Landedness distinguishes landed, unmerged, detached, and
  unknown without encoding proof or error state into a prefix-tested string.
  Ownership, assignment join, lock, lease, recovery agreement, tracked state,
  nested state, and ignored inventory likewise enter the decision as typed
  facts. Exact unexported type and file names are implementer discretion.
- **Two readings, one policy owner.** Explicit cleanup projects the base
  verdict. Automatic cleanup asks the same owner for the stricter reading that
  requires verified ownership, cleanup-pending assignment state, matching
  recovery metadata, landedness, and no work that would require preservation.
  The `--landed` selector's “per-path cleanup is required to preserve work”
  refusal is part of this owner rather than a post-plan action rewrite.
- **Ordered behavior stays exact.** The characterization matrix is written
  against the pre-refactor tree and lands before the move. Its expected tuples
  are not generated from production rules. Production and tests may share
  fixture builders, but they do not share the outcome enumeration.
- **DiscardBranch remains derived-after.** Derived landedness is recorded first.
  The explicit operator assertion may then authorize exact branch deletion.
  It does not alter the evidence automatic cleanup reads and does not make an
  unmerged assignment automatic-cleanup eligible.
- **Execution states are not eligibility verdicts.** `removed`, `error`,
  `release-remove`, and `release-leftover` describe execution or the decayed
  registration lifecycle. They are not added to the eligibility matrix and do
  not become evidence rules. Terminal execution may still project `removed`.
- **Exit-test rule.** Every ticket keeps the pre-existing suite passing with
  test logic unmodified. Mechanical test renames are the only permitted edits
  to an existing test. A changed assertion or expected outcome stops the move
  and routes the delta through a separately reviewed feature or `/bench-debug`.
- **ADR update travels with the contract.** The implementation rewrites
  `docs/adr/0005-worktree-cleanup-requires-verifiable-ownership.md` as
  resulting-state documentation. It names the one eligibility verdict without
  weakening marker + assignment + lock ownership or preservation requirements.
- **No bootstrap authority.** The first production writer gets no temporary
  exemption. Stories 1–2 establish the independent outcome oracle before the
  eligibility module is allowed to own policy.

## Testing decisions

- New characterization attaches at the existing package seam through
  `PlanExplicitWithOptions` and `PlanAutomatic`. The tests use real fixture
  repositories and existing worktree/intent helpers, observe the public plan
  projection, and compare an independently authored expected tuple.
- `TestExplicitEligibilityOutcomeMatrix` in
  `internal/worktree/eligibility_test.go` owns the nine explicit tuple cases.
  `TestAutomaticEligibilityOutcomeMatrix` owns the thirteen automatic tuple
  cases. Each subtest names the final tuple and the competing evidence when a
  later rule is being pinned.
- `TestEligibilityVerdictProjectsWithoutSecondDecision` calls the eligibility
  seam and the explicit and automatic plan projections over the same fixture.
  It asserts the projected action, reason, typed landedness, recovery choice,
  and branch-deletion authority agree with the verdict.
- Existing command tests remain the composition oracle:
  `TestExplicitApplyRejectsContentDriftWithoutMutation`,
  `TestDiscardBranchLeavesTheDerivedClassificationUnchanged`,
  `TestDiscardBranchNeverBypassesARefusal`,
  `TestPlanAutomaticKeepsEarlierRetainReason`,
  `TestReleaseAndPlanExplicitAcceptUnstampedAssignment`,
  `TestReleaseReconcilesCompletedAutomaticCleanup`, and the landed-set apply
  tests. `internal/worktree/landed_test.go` and
  `internal/worktree/list_actions_test.go` continue to cover the shared
  landed-assignment reader and its resume/list consumers. Landing continues to
  cross the same `releaseLandingAssignment` seam.
- The ordinary `test` phase (`go test -count=1 ./...`) is the gate seam. No gate
  or shell edit is required. The hostile-shell checklist is not applicable
  because the build adds no shell boundary or shell argument.

### Seam diagram

    Git + filesystem + intent evidence
                    │
                    ▼
          [ worktree eligibility ]
          [ typed facts → verdict ]
                    │
         ┌──────────┼───────────┐
         ▼          ▼           ▼
    explicit plan  automatic   landed-set policy
         │          policy           │
         └──────┬────┴───────┬───────┘
                ▼            ▼
        apply / release   land release
                │
                ▼
       existing plan rows, receipts, recovery refs, exact removal

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| EX1 | 1 | an explicit primary or unsafe target ends `retain/uncertain` with its current detail and planning mutates no durable state | `TestExplicitEligibilityOutcomeMatrix/retain-uncertain` | First-write red against a stub verdict and later red if the rule returns remove or another reason. Pins the fail-closed tuple and the final unsafe-target override. |
| EX2 | 1 | an unregistered target ends `retain/foreign` | `TestExplicitEligibilityOutcomeMatrix/retain-foreign` | Mutating registration absence to eligible is red. Prevents a path convention from becoming ownership proof. |
| EX3 | 1 | malformed owner or declaration evidence ends `retain/malformed` with the current later-rule winner | `TestExplicitEligibilityOutcomeMatrix/retain-malformed` | Mutating malformed evidence to uncertain or allowing removal is red. Distinguishes malformed known evidence from generic uncertainty. |
| EX4 | 1 | a foreign lock or owned lock mismatch ends `retain/unexpected-lock` | `TestExplicitEligibilityOutcomeMatrix/retain-unexpected-lock` | Ignoring the lock mismatch is red. Pins the protective-lock conjunct from ADR 0005. |
| EX5 | 1 | live or ambiguous lease evidence ends `retain/live-lease` when it is the current final applicable refusal | `TestExplicitEligibilityOutcomeMatrix/retain-live-lease` | Moving lease evaluation without preserving the winner is red. Pins the effective lease position instead of an ideal order. |
| EX6 | 1 | bounded ignored residue without authority or over the destructive limit ends `retain/ignored` | `TestExplicitEligibilityOutcomeMatrix/retain-ignored` | Treating ignored residue as removable is red. Pins preservation of undeclared or excessive ignored bytes. |
| EX7 | 1 | a clean registered removable checkout ends `remove` with an empty reason code | `TestExplicitEligibilityOutcomeMatrix/remove` | Requiring assignment ownership for the explicit path is red. Preserves the explicit command's current registered foreign-clean case. |
| EX8 | 1 | dirty or detached state that is otherwise removable ends `recover-remove` with its current recovery evidence | `TestExplicitEligibilityOutcomeMatrix/recover-remove` | Removing without choosing recovery is red. Pins preservation before explicit removal. |
| EX9 | 1 | authorized bounded ignored residue that is otherwise removable ends `discard-remove` with an empty reason code | `TestExplicitEligibilityOutcomeMatrix/discard-remove` | Retaining despite explicit discard authority or removing without the discard action is red. Pins the separate ignored-disposal authorization. |
| AU1 | 2 | unknown landedness or assignment state ends `retain/uncertain` | `TestAutomaticEligibilityOutcomeMatrix/retain-uncertain` | Interpreting unknown as unmerged or removable is red. Prevents automatic cleanup from converting uncertainty into eligibility. |
| AU2 | 2 | a registration without the verified owner and assignment join ends `retain/foreign` | `TestAutomaticEligibilityOutcomeMatrix/retain-foreign` | Reusing explicit registration-only eligibility is red. Pins automatic cleanup's stricter ownership reading. |
| AU3 | 2 | malformed assignment or mismatched recovery metadata ends `retain/malformed` | `TestAutomaticEligibilityOutcomeMatrix/retain-malformed` | Skipping recovery agreement is red. Pins the recovery conjunct from ADR 0005. |
| AU4 | 2 | an unexpected lock surviving explicit planning remains `retain/unexpected-lock` | `TestAutomaticEligibilityOutcomeMatrix/retain-unexpected-lock` | Replacing the reason with a generic foreign refusal is red. Preserves the operator-visible refusal class. |
| AU5 | 2 | a live lease ends `retain/live-lease` even when another explicit refusal is also present | `TestAutomaticEligibilityOutcomeMatrix/retain-live-lease` | Returning the explicit refusal before the automatic lease check is red. Pins the first automatic override. |
| AU6 | 2 | an aged unlanded assignment with ignored residue ends `retain/ignored` | `TestAutomaticEligibilityOutcomeMatrix/retain-ignored` | Replacing the earlier explicit refusal with orphaned is red. Preserves the outcome already pinned by `TestPlanAutomaticKeepsEarlierRetainReason`. |
| AU7 | 2 | a young active unlanded assignment ends `retain/active` | `TestAutomaticEligibilityOutcomeMatrix/retain-active` | Allowing an active assignment into removal is red. Pins assignment lifecycle gating. |
| AU8 | 2 | an active landed assignment ends `retain/landed` with the current override behavior | `TestAutomaticEligibilityOutcomeMatrix/retain-landed` | Returning active or an earlier explicit refusal is red. Pins the existing landed-active classification without approving cleanup. |
| AU9 | 2 | an aged active unlanded assignment ends `retain/orphaned` | `TestAutomaticEligibilityOutcomeMatrix/retain-orphaned` | Returning active after the age boundary is red. Pins the shared orphan predicate's place in automatic policy. |
| AU10 | 2 | a cleanup-pending assignment whose branch has not landed ends `retain/unmerged` | `TestAutomaticEligibilityOutcomeMatrix/retain-unmerged` | Inferring landing from branch naming or `DiscardBranch` is red. Makes the explicit-only assertion boundary bite. |
| AU11 | 2 | a cleanup-pending landed assignment needing preservation ends `retain/dirty` | `TestAutomaticEligibilityOutcomeMatrix/retain-dirty` | Carrying recover-remove into unattended cleanup is red. Prevents automatic cleanup from authoring recovery refs. |
| AU12 | 2 | a verified cleanup-pending landed clean assignment ends `remove` with an empty reason code | `TestAutomaticEligibilityOutcomeMatrix/remove` | Leaving the successful automatic case retained is red. Proves the stricter reading still admits its one clean path. |
| AU13 | 2 | a verified cleanup-pending landed clean assignment with declared bounded ignored output ends `discard-remove` | `TestAutomaticEligibilityOutcomeMatrix/discard-remove` | Rejecting declared build output or discarding undeclared output is red. Pins the bounded declaration exception. |
| EV1 | 3 | one eligibility call returns the decision and typed evidence for every named ownership and safety fact while `PlanExplicitWithOptions` only projects it | `TestEligibilityVerdictProjectsWithoutSecondDecision` | Wrapping the already-mutated `CleanupPlan` in a verdict at the end is red. Catches a rename that leaves policy scattered. |
| EV2 | 3 | only the eligibility module orders or selects eligibility actions and reasons before execution | source review over `internal/worktree` | Production assignments to eligibility actions or refusals outside the module fail review. Makes the one-source claim structural and reviewable. |
| EV3 | 3 | fingerprints still bind every removal-relevant evidence byte and decision projection | existing stale-fingerprint tests plus `TestEligibilityVerdictProjectsWithoutSecondDecision` | Omitting typed landedness or a verdict field from the fingerprint is red. Prevents a moved policy from weakening plan/apply identity. |
| CO1 | 4 | `PlanAutomatic` and the shared `assignmentLanded` reader obtain their answer from typed verdict evidence and no production code parses `plan.landed` or another formatted decision string | automatic matrix plus source review across automatic, resume, list, and landed-set callers | A surviving `strings.HasPrefix(plan.landed, ...)` or equivalent fails review. Directly catches the leakage FT216 exists to remove. |
| CO2 | 4 | explicit apply and automatic apply replan under the transaction lock and execute only an eligible verdict with the same stale-fingerprint and interruption outcomes | existing apply and lifecycle tests | Bypassing the verdict after lock or executing a retain verdict is red. Pins decision use at the mutation boundary. |
| CO3 | 4 | `--landed` planning obtains its preservation refusal from the eligibility owner and keeps its set fingerprint and no-op rows unchanged | landed-set tests including `TestCleanLandedRetainsUnparseableLeaseAndSkipsUnprovableBranch` | Leaving `plan.Action = ActionRetain` in `planLandedAssignment` fails review. Catches the third decision site after explicit and automatic planning. |
| CO4 | 4 | `ReleaseCommand` and both first-run and resumed `LandCommand` release only through the automatic verdict and retain their exact receipts and diagnostics | release and land command tests | Calling cleanup execution from release without the verdict or changing a release diagnostic is red. Pins the lifecycle consumers named by the map. |
| DB1 | 4 | `DiscardBranch` may authorize exact branch deletion after derived landedness only on explicit cleanup and automatic cleanup still returns `retain/unmerged` for the unprovable branch fixture | `TestDiscardBranchLeavesTheDerivedClassificationUnchanged` and `TestDiscardBranchNeverBypassesARefusal` | Exposing the assertion in automatic evidence or letting it bypass a refusal is red. Preserves the decided derived-after boundary. |
| OS1 | 5 | every pre-existing test passes with test logic unchanged except mechanical renames | ordinary `test` phase plus diff review of existing `_test.go` files | Any existing assertion or expected outcome change fails review. Enforces the shared behavior-preserving exit proof. |
| OS2 | 5 | ADR 0005 names the eligibility verdict and still requires verifiable conjunctive ownership plus preservation before removal | ADR review against EV1 and automatic matrix rows | Removing marker, assignment, lock, landedness, recovery, or preservation from the ADR fails review. Keeps documentation and the executable result aligned. |
| OS3 | 5 | deleting the eligibility module would force ordered decision policy back into at least explicit, automatic, and landed-set consumers | source review | A passive data wrapper that can be deleted without redistributing policy fails review. Applies the deletion test that justifies the seam. |

### Edge inventory

- **Error path:** unreadable worktree inventory and Git evidence remain errors at
  their current command surfaces. Malformed durable ownership or assignment
  evidence remains a retain verdict in EX-M and AU-M.
- **Empty/absent input:** the primary checkout and unregistered paths are pinned
  by EX-U and EX-F. Missing owner markers on registered unlocked worktrees keep
  the current explicit-versus-automatic split through EX-RM and AU-F.
- **Boundary values:** ignored entry and byte limits remain covered by
  `TestIgnoredInventoryEntryAndByteBoundaries`. Assignment age uses the existing
  exact-window tests. No bound changes in this spec.
- **Malformed input:** strict owner-marker shapes and incomplete assignment joins
  remain covered by `TestPlanAutomaticRejectsEveryInvalidMarkerWithoutMutation`
  and `TestPlanAutomaticRequiresCompleteAssignmentJoin`.
- **Interrupted / partial state:** cleanup receipts, relocking, recovery refs,
  and terminal replay remain covered by the existing transaction and fault-step
  tests. The verdict is recomputed before each mutation checkpoint as today.
- **Re-run idempotency:** stale explicit apply, completed release replay, and
  resumed landing retain their current receipts and output.
- **Process-boundary lifecycle:** `LandCommand` crosses into release through the
  existing `releaseLandingAssignment` seam. The verdict remains in-process and
  is not serialized as a new durable schema.
- **Hostile environment:** unsafe control bytes, nested repositories, uncertain
  ignored inventory, hostile private-admin shapes, and unknown leases retain.
- **Won't handle:** correcting any precedence anomaly exposed by the matrix.
  That is a separately reviewed `/bench-debug` repair, estimated 4–8 production
  and test edits plus 2 full gate runs after the exact anomaly is chosen.
- **Won't handle:** changing cleanup TOON fields, reason strings, fingerprints,
  receipt schemas, or command grammar. That is observable CLI/schema work,
  estimated 6–12 edits plus conformance updates and 2–3 full gate runs.
- **Won't handle:** folding decayed registration release (`release-leftover`) into
  eligibility. It deliberately preserves path bytes rather than answering safe
  removal, estimated 5–8 lifecycle/test edits plus 2 full gate runs if reopened.
- **Won't handle:** FT218's named private-admin Git reader. It remains a separate
  light-path ticket from the compiled map, estimated 6 reader/caller edits plus
  2 full gate runs.
- **Won't handle:** FT207's malformed-admin policy decision. Any contradiction
  discovered here is flagged for that decision path rather than silently
  choosing its behavior.

## Ownership fences

One writer at a time. Tickets are serial because each later slice reads the
preceding slice's settled test or production seam. Reviewer disposition: approve,
merge, or split each fence during sign-off.

- Stories 1–2: `internal/worktree/eligibility_test.go` and existing fixture
  helpers under `internal/worktree/` only when a new case cannot use a current
  helper. No production edit.
- Story 3: `internal/worktree/eligibility.go`,
  `internal/worktree/subshell.go`, `internal/worktree/classifier.go`, and
  `internal/worktree/eligibility_test.go`.
- Story 4: `internal/worktree/classifier.go`,
  `internal/worktree/resume.go`, `internal/worktree/clean_landed.go`,
  `internal/worktree/ownership.go`, `internal/worktree/worktree.go`,
  `internal/worktree/land.go`, `internal/worktree/landed.go`,
  `internal/worktree/list.go`, `internal/worktree/lifecycle.go`,
  `internal/worktree/worktree_test.go`,
  `internal/worktree/resume_test.go`,
  `internal/worktree/recovery_retry_test.go`,
  `internal/worktree/clean_branch_test.go`,
  `internal/worktree/orphan_test.go`,
  `internal/worktree/ownership_test.go`,
  `internal/worktree/lifecycle_test.go`,
  `internal/worktree/clean_landed_apply_test.go`, and
  `internal/worktree/land_test.go`, `internal/worktree/landed_test.go`, and
  `internal/worktree/list_actions_test.go`.
- Story 5: `internal/worktree/eligibility.go`,
  `internal/worktree/subshell.go`, `internal/worktree/classifier.go`,
  `internal/worktree/clean_landed.go`, `internal/worktree/resume.go`,
  `internal/worktree/ownership.go`, `internal/worktree/worktree.go`,
  `internal/worktree/land.go`, `internal/worktree/landed.go`,
  `internal/worktree/list.go`, `internal/worktree/lifecycle.go`,
  `docs/adr/0005-worktree-cleanup-requires-verifiable-ownership.md`, and tests
  for mechanical renames only.
- Every ticket: `specs/worktree-cleanup-eligibility/` for ticket checkboxes and
  lifecycle status changes.

## Out of scope

- Any new eligibility outcome, refusal reason, reason wording, or precedence.
- A public or cross-package eligibility API. This remains an in-process
  `internal/worktree` decision.
- A new dry-run, cleanup flag, command, output column, or receipt field.
- Changes to `git.Output`, `git.Raw`, or `git.OK` and FT218's named readers.
- Changes to adoption lifecycle planning (FT217).
- Changes to the worktree enumeration and hostile-admin policy already landed
  under FT189, except consuming its current evidence exactly as it exists.
- FT108's future refactor skill. The exit test is carried explicitly here.
