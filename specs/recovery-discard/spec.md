# Recovery discard

Status: staged

Decision source: reviewer-confirmed conversation 2026-08-03 — the `/bench-debug` session that diagnosed provisional-ref and preserved-orphan accumulation; the reviewer chose per-ref `--discard` with an exact fingerprint over age-based expiry and over reporting-only suppression.

## Problem

`bench worktree release` preserves any dirty checkout it cannot prove landed, writing a recovery ref and leaving the assignment row open. That is correct: the payload's content genuinely differs from the default branch, so dropping it silently would destroy work.

What is missing is the other half. `internal/worktree/resume.go` describes preserved orphans as awaiting "a deliberate recover-or-retire", and retire exists only for payloads that pass the landedness proof. An operator who has inspected a preserved payload and decided they do not want it has no way to say so, so the preserved set only grows. Enumerate the current backlog with `bench worktree recovery` over each ref under `refs/bench/recovery/`; the spec deliberately names that enumeration rather than restating its count, which drifts.

The backlog holds three classes, not one, and they need different handling:

- payloads the landedness proof accepts, which `--apply` already retires;
- payloads the proof refuses, whose assignment row still exists;
- recovery refs whose assignment row is already gone, which `PlanRecovery` cannot classify at all because it resolves a ref through its owning row.

A fourth residue class has the same shape but a different cause. Promotion did not reclaim its provisional refs until `4ee6e5f`, so runs that promoted before it left assignment branches, candidate refs, and checkpoint refs behind. Those runs are terminal and the residue is dead, but nothing enumerates it: release compacted the intent rows that named the branches, and the run records those runs persisted predate the `Branch` field `4ee6e5f` introduced.

## Solution

Add one deliberate discard to the existing recovery command, mirroring the `--apply` contract exactly: plan first, exact fingerprint required, one ref per invocation. Discarding deletes the recovery ref and closes the assignment row without asserting the payload landed; the Git objects remain reachable through the reflog until garbage collection, so the act is auditable and briefly reversible. A recovery ref with no owning row is discardable on the strength of the ref alone, because there is no row left to consult and no other command can reach it.

The plan surface grows enough detail that the operator is not choosing blind — it reports what the payload changes relative to its base, so a one-file leftover and a large abandoned build are distinguishable before either is dropped.

Separately, a maintainer can reclaim the provisional residue of terminal spec-build runs, applying retroactively the same policy the promotion path now applies going forward.

## User stories

1. **An operator can discard a preserved payload they have inspected.** `bench worktree recovery <ref> --discard <fingerprint>` retires a recovery ref whose payload is not proven landed. It requires the exact fingerprint the plan just reported, refuses a stale one, deletes the ref, and closes the assignment row by the same compaction `--apply` uses when the last recovery ref leaves an assignment. A ref whose plan action is `retire` — one the proof already accepts — is refused, because `--apply` is its route and the two claims must stay distinguishable in the receipt.
   `Line: opus / medium.` The grammar mirrors existing prior art, but the operation destroys preserved work and every refusal has to fail closed.

2. **A recovery ref with no owning assignment row is reachable.** Planning such a ref today yields the same "no recovered assignment" retain verdict that a fully-discarded ref yields, so the two states are indistinguishable and neither is actionable. The plan distinguishes them: a ref that exists with no owning row is reported as orphaned and is discardable; a ref that does not exist is reported as absent and is a no-op success. Discarding an orphaned ref deletes the ref and touches no intent record.
   `Line: opus / medium.` This is the class the first draft of this spec missed entirely, and it is the class with no other route out.

3. **The plan says what a discard would drop.** Planning a recovery ref reports the payload's change summary against its recorded base — how many paths it touches — so the operator can tell a trivial leftover from real work before discarding. The field is present for every plan, including refs whose action is `retain` or `retire`, and it degrades to a definite unknown rather than an error when the base or payload no longer resolves.
   `Line: opus / medium.` The value is judgment support, so the honest-unknown path matters as much as the happy path.

4. **A maintainer reclaims provisional residue from runs that promoted before reclamation existed.** A maintainer-run pass enumerates assignment branches, candidate refs, and checkpoint refs belonging to spec-build runs whose record is terminal, reports them as a plan, and deletes them on an exact fingerprint. Refs belonging to a non-terminal run, or to no run record at all, are reported and left intact rather than guessed at. The pass reaches records written before `4ee6e5f`, which persist no branch name.
   `Line: opus / medium.` The enumeration is mechanical, but a wrong classification deletes a live build's working state.

5. **Recovery planning and application are proven through the real producer.** The tests that grade discard and retire drive a recovery ref produced by an actual `bench worktree release`, not a hand-built branch. The existing landedness fixtures stay as unit coverage of the comparison itself; the new coverage composes release and recovery so a plan path that production can never reach cannot pass as covered.
   `Line: opus / medium.` This is the junction defect class `ROADMAP.md` FT190 names, and the session that wrote this spec twice reasoned wrongly about the retire path because only fixture-shaped inputs were covered.

6. **An interrupted retire or discard converges on re-run.** Both verbs delete the recovery ref before closing the assignment row, so a process that dies between the two leaves a row naming a ref that no longer exists. Re-running either verb against that state closes the row rather than refusing, and the fault is injected at the existing named-step seam rather than simulated by hand-editing the record. This holds for `--apply` and `--discard` alike.
   `Line: opus / medium.` The ordering is pre-existing and untested for either verb; the coverage is what makes the crash window survivable rather than latent.

## Implementation decisions

**Discard is an action on the existing plan, not a second planner.** `PlanRecovery` remains the one classifier. It already computes the fingerprint over the ledger, ref, root, payloads, landedness, action, and detail; discard consumes that same plan and the same fingerprint. Adding a parallel discard planner would duplicate the classification knowledge the fingerprint is derived from.

**The plan distinguishes orphaned from absent.** Today both a row-less ref and a nonexistent ref reach the same retain verdict. The classifier separates them, because one is actionable and the other is already done. Resolving a ref that exists without an owning row is a ref read, not a record read, so the classifier gains a ref-existence probe it does not have today.

**The action vocabulary gains one terminal value.** The plan's action set today is retain, retire, retired, and error. Discard introduces one further terminal outcome distinguishing "the operator dropped unproven work" from "the tool proved this landed and retired it", because the two are different claims about the same disappearance and the receipt is the only durable record of which happened. Whether the discard-eligible state appears as a distinct planned action or as a modifier on `retain` is implementation discretion, provided the emitted terminal value distinguishes the two.

**The fingerprint's domain tag and effect string change with its authority.** The existing fingerprint commits to a retire-specific domain tag and effect list. Once the same fingerprint can authorize a destructive discard, both must name that authority, so a fingerprint planned under the old semantics cannot authorize the new operation. Relatedly, `ApplyRecovery` currently returns silent success when a matching fingerprint accompanies a retain-action plan. That branch refuses for **both** verbs, not only discard: a caller who supplied a fingerprint asked for an action, and reporting exit-zero success when nothing happened reads as "the work is gone" when it is not. This is a deliberate behaviour change to `--apply`; no existing test depends on the silent success.

**Discard never widens to a sweep.** The command accepts exactly one ref per invocation and no glob, no `--all`, and no predicate that selects a set. This is the reviewer's chosen shape: retiring unproven work stays a typed, per-ref act.

**Discard does not weaken the landedness proof.** `LandedInDefault` and the retire path are unchanged. Discard is refused for a payload the proof accepts. The proof is the oracle for "did this land"; discard is an operator assertion of "I do not want this", and the two must not be conflated in the receipt.

**Both recovery verbs share one argument parse.** The existing fingerprint-format validation — length, hex, lowercase — lives inside the `--apply` arm of `RecoveryCommand`. Discard does not get a second parse; the arm generalizes over the verb so the format control is reached by construction rather than by a second copy of the same knowledge.

**The change summary is derived, never stored.** The plan computes the payload's path count against its recorded base at plan time. Persisting it in the intent record would create a second derivation of a fact Git already holds, and it would go stale the moment the base is rewritten.

**Provisional residue is classified from the run records, by assignment identity.** A `bench/assign/*` branch is dead because the spec-build record that created it is terminal — not because its name matches a pattern or looks merged. The link from record to branch is the assignment ID, which every record persists in both the pre- and post-`4ee6e5f` shapes; the branch is located by matching that ID against refs in the assignment namespace. Location is not classification: terminality still comes from the record. The `Branch` field `4ee6e5f` added is a fast path when present, never the only path, because no pre-fix record carries it. An ID matching more than one ref in the namespace is reported ambiguous and retained, never deleted. Refs whose owning record is absent are reported unclassified and retained, because an absent record cannot prove the work is dead.

**Reclamation is per-run and typed, like every other destructive verb here.** The surface is `bench spec build reclaim <slug>` to plan and `--apply <fingerprint>` to execute, matching the existing `abandon` plan/apply shape rather than inventing a new grammar. It takes one slug per invocation for the same reason discard takes one ref: a maintainer reclaiming residue is making a judgment per run, not authorizing a sweep.

**Reclamation and promotion share one enumeration.** The retroactive pass and `reclaimProvisionalRefs` must consume one function answering "which refs does this terminal run no longer need", derived by assignment identity as above. Two independent lists is the duplicated-knowledge defect the code standard forbids — and the first draft of this spec proved the risk concretely by specifying an enumeration that silently skipped every pre-fix record.

**The reclamation pass lives with the records it reads.** Run records and their enumeration are unexported in `internal/specbuild`, and `internal/specbuild` already imports `internal/worktree`, so the pass and its tests cannot live in `internal/worktree` without an import cycle.

## Testing decisions

- A good test drives the real commands — release to produce a preserved orphan, recovery to plan and discard it — and asserts the observable end state: the ref is gone, the assignment row is closed, and the receipt names discard rather than retire.
- Seams receiving tests: the recovery plan/apply surface in `internal/worktree` (prior art: the existing recovery and cleanup fixtures), the reclamation enumeration in `internal/specbuild` (prior art: the abandon-plan fixtures, which already grade ref inventories), and the runtime command surface in `internal/contract/runtime` (prior art: the worktree runtime fixtures that drive `bench worktree` end to end).
- The gate observes this through the existing contract, worktree, and specbuild packages; no new gate phase.

### Seam diagram

    trigger: operator runs `bench worktree recovery <ref> [--discard <fp>]`
        │
        ▼
    recovery ref  ──▶  [ PlanRecovery: landedness + orphan/absent + change summary ]  ──▶  plan + fingerprint
        │                                                                                      │
        │                                                                                      ▼
        └──────────────▶  [ ApplyRecovery / discard: ref delete + row close ]  ──▶  receipt
                              ◀ tests attach here: a real `bench worktree release`
                                produces the ref; the test drives plan then apply
                                and asserts refs, intent rows, and receipt action

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | discarding an unlanded payload deletes the ref and closes the row | internal/worktree | not yet observed — no discard path exists to fail | the ref and row survive a discard attempt, which is the whole defect |
| 1 | a stale fingerprint refuses and changes nothing | internal/worktree | not yet observed | a discard that ignores the fingerprint would drop work the operator never saw planned |
| 1 | `--discard` on a ref whose plan action is retire refuses | internal/worktree | not yet observed | otherwise discard silently becomes a second route for proof-accepted work and the receipt lies about which claim was made |
| 1 | a fingerprint that is not 64 lowercase hex characters refuses before planning, via `--discard` | internal/worktree | not yet observed | a second parse for the new verb would skip the format control entirely |
| 2 | a ref with no owning row plans as orphaned and is discardable | internal/worktree | not yet observed | this class has no other route out and the first draft left it stranded |
| 2 | a ref that does not exist plans as absent and discard is a no-op success | internal/worktree | not yet observed | conflating absent with orphaned makes re-runs either fail or delete blind |
| 3 | the plan reports the payload's changed-path count | internal/worktree | not yet observed | without it the operator discards blind, which is how real work gets dropped |
| 3 | an unresolvable base or payload yields a definite unknown, not an error | internal/worktree | not yet observed | an error here would make an otherwise-discardable ref unplannable and re-strand it |
| 4 | residue of a terminal run is planned and deleted on exact fingerprint | internal/specbuild | not yet observed | the pre-fix branches stay stranded if the enumeration misses them |
| 4 | the enumeration locates branches for records that persist no branch name | internal/specbuild | not yet observed | matching on the stored branch field alone yields empty for every pre-fix record, silently reclaiming nothing |
| 4 | residue of a non-terminal run is reported and retained | internal/specbuild | not yet observed | deleting a live build's candidate destroys in-flight work |
| 4 | refs with no owning record are reported unclassified and retained | internal/specbuild | not yet observed | an absent record cannot prove death; guessing here is unrecoverable |
| 4 | an assignment ID matching more than one namespace ref is reported ambiguous and retained | internal/specbuild | not yet observed | deleting both halves of an ambiguous match destroys a ref no record claimed |
| 4 | `bench spec build reclaim <slug>` plans and `--apply` executes, per run | internal/contract/runtime | not yet observed | without a command row the CLI shape is a delegate's guess and the package rows cannot see it |
| 5 | plan and discard are graded against a ref produced by a real release | internal/contract/runtime | not yet observed | fixture-shaped refs hid the real behaviour of the retire path twice in one session |
| 6 | a discard interrupted after ref deletion closes the row on re-run | internal/worktree | not yet observed | the crash window leaves a row naming a deleted ref, which no other command can close |
| 6 | an `--apply` interrupted after ref deletion closes the row on re-run | internal/worktree | not yet observed | the pre-existing verb has the identical window and was never covered; leaving it uncovered keeps the latent defect |

**Degenerate implementation.** The cheapest wrong build makes `--discard` an unconditional `update-ref -d` that ignores the fingerprint and the row; story 1's second row goes red on it. A build that accepts any ref regardless of plan action is caught by story 1's third row. The cheapest wrong build of story 4 deletes every `bench/assign/*` ref by name pattern; story 4's non-terminal and no-record rows go red on it. The cheapest wrong build that reuses only the stored branch field reclaims nothing at all and is caught by story 4's second row.

**Composition degenerate.** Because stories 1 and 5 span `internal/worktree` and `internal/contract/runtime`, the composition degenerate is a discard that passes every unit row against hand-built refs while the real `bench worktree release` produces a ref shape the plan rejects. Story 5's row, driven through the real producer, is the only row that goes red on it — which is precisely the defect this session hit.

### Edge inventory

- **Error path** — stale fingerprint, retire-action ref, malformed fingerprint: rows above.
- **Empty/absent input** — refs with no owning record (story 4) and refs with no owning row (story 2): rows above.
- **Boundary values** — an assignment holding several recovery refs: discarding one leaves the row open at `recovered` and only the last closes it, matching `ApplyRecovery`'s existing behaviour; covered by story 1's first row via the shared compaction.
- **Malformed input** — covered by story 1's fourth row, which asserts the shared parse rather than assuming it.
- **Interrupted or partial state** — the ref is deleted before the row is closed; a re-run finds no ref and must still close the row. Covered for **both** verbs by story 6's rows. An earlier draft excluded this on the grounds that `--apply` had no such coverage either; the reviewer overruled that, correctly — parity is reachable by adding the missing coverage rather than withholding the new coverage, and interrupted-state handling of the verb being added is the rest of this capability, not a separate one.
- **Re-run idempotency** — discarding an already-discarded ref reaches the absent verdict in story 2's second row, which is now distinct from the orphaned verdict rather than sharing it.
- **Process-boundary lifecycle** — the intent record is reloaded by a fresh process between plan and apply: the fingerprint covers the ledger bytes, so a concurrent mutation invalidates it; covered by story 1's stale-fingerprint row.
- **Hostile environment** — a recovery ref naming a payload outside the repository, or a ref name carrying control characters: the existing `verifyRecovery` envelope check and the TOON control-escaper own both, and both run before the new action is reachable. For the orphaned-ref path in story 2 there is no row to envelope-check, so that path relies on the ref-name control escaper alone; story 2's rows drive real ref names through the renderer.

## Out of scope

- **Automatic expiry of preserved orphans.** The reviewer chose the per-ref typed act over a timer. A future spec could add age-based reporting without deletion. Estimate if reconsidered: 4 edits, 3 gate runs.
- **Retiring this repository's existing backlog.** That is an operational pass run under stories 1, 2, and 4 once they ship, not a code change, and it needs the reviewer at the keyboard for each judgment. Every class in the current backlog is reachable by those three stories; the first draft of this spec claimed the same thing while leaving the row-less class stranded, which is why story 2 exists.
- **Changing what `LandedInDefault` proves.** The comparison is correct — some refs in the current backlog pass it today. Widening it to tree-equivalence is a separate capability with its own correctness argument. Estimate: 6 edits, 4 gate runs.
- **A general injected-interface real-producer audit.** Story 5 covers this one junction. `ROADMAP.md` FT190 owns the tree-wide sweep.
