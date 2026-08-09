# Repair-ticket reslicing

Status: implemented

Decision source: reviewer-confirmed conversation on 2026-08-08

## Problem

The debug-repair guidance says a validated receipt “takes a repair ticket.” A
coordinator can read that singular literally and force every required repair
path into one wide ticket, even when producer and consumer outcomes can land
green independently. That defeats the ordinary ticket-sizing rule and delays
the blocked assignment behind work that could have landed as an ordered repair
chain.

## Solution

Treat the receipt's required fence as the maximum envelope for repair
reslicing. `craft-tickets` applies its ordinary independently-green rule inside
that envelope and may produce either the common one-ticket repair or a
reciprocal ordered chain whose combined ownership never escapes the receipt.
`/bench-implement-spec` points to that owner and orchestrates every resulting
ticket through the ordinary lifecycle, refreshing the original blocked
assignment with its original validated receipt only after the terminal repair
ticket lands.

## User stories

1. As a coordinator, I want `craft-tickets` to reslice a debug receipt inside
   its required fence, so independently-green repair outcomes can land as one
   ticket or a reciprocal ordered chain without widening repair ownership.
   Line: gpt-5.6-sol / high. Skill prose has the kit leverage override, and the
   independently-green split remains semantic beyond the anchor checks.
2. As a coordinator, I want `/bench-implement-spec` to orchestrate every ticket
   in a receipt-derived repair chain normally and wait for the terminal ticket
   before refreshing the blocked assignment, so partial repair state cannot
   authorize resumed implementation.
   Line: gpt-5.6-sol / high. Command prose has the same leverage override and
   steers every future repair lifecycle.

## Implementation decisions

- A validated debug receipt continues to fix three facts: confirmed cause,
  maximum repair ownership fence, and the condition under which the blocked
  assignment may proceed. The receipt does not choose the number of repair
  tickets.
- `craft-tickets` is the single owner of repair reslicing. It applies the
  existing independently-green split rule inside the receipt's maximum
  envelope. The result may be one repair ticket or a reciprocal ordered repair
  chain, and the union of every repair-ticket ownership fence must remain
  inside the receipt's required fence.
- `/bench-implement-spec` owns lifecycle orchestration only. It points to
  `craft-tickets` instead of restating the split rule; says one repair ticket is
  still the common case; sends every repair-chain ticket through ordinary
  assignment, checkpoint, and integration; treats the terminal repair ticket
  as the receipt's proceed condition; and refreshes the original blocked
  assignment with the original validated receipt only after the complete chain
  has landed.
- The repair chain is an ordered producer-to-consumer sequence when a split
  survives contract discovery. Its reciprocal edge uses the existing
  `Integration surfaces:` and `Blocked by:` forms; this spec introduces no new
  ticket schema or lifecycle operation.
- The prose change is covered at the existing workflow-guidance conformance and
  canary seam. Extend the literal-mutation table in
  `internal/conformance/fixture_bite_test.go`, which already reads the real
  command and skill, applies one subject replacement, and checks the exact
  workflow-anchor diagnostic. Do not add four fixture directories or copy a
  command, skill, fixture inventory, or anchor registry.
- Positive policy uses section-scoped `RequireInSection` anchors and omission
  mutations. Contradictory policy uses section-scoped `ForbidInSection` anchors
  and additive mutations that leave the positive prose intact. In particular,
  the singular command mandate, an escaped repair fence, an early refresh, and
  chain-local raw Git plumbing are additive contradictions; removing the
  receipt-envelope rule, the skill's one-or-chain result, the command's
  `craft-tickets` pointer, or the one-ticket common-case clause are omissions. A
  require-only implementation is insufficient because it cannot reject correct
  prose followed by a contradictory permission.
- The two story fences form a sequential chain rather than a concurrent
  frontier because both add their own conformance requirements to
  `internal/anchors/registry_data.go`. The first fence owns
  `.agents/skills/bench-craft-tickets/SKILL.md`, its three repair-reslicing
  literal-mutation table rows, and its registry rows. The second fence, blocked by the
  first, owns `.agents/commands/bench-implement-spec.md`, its two lifecycle
  mutation groups, and its registry rows. Reusing
  `internal/anchors/registry_data.go` and
  `internal/conformance/fixture_bite_test.go` after the first ticket lands is
  deliberate sequential ownership, not concurrent dual ownership.
- Apply `craft-synthesis` and `craft-skills` during implementation. The
  consistency loop re-runs exact-literal and integration-surface discovery
  after coverage lands. The prose-only dogfood substitution is the green
  landing gate plus a read of every surface this prose steers; the implementation
  also applies `craft-tickets` and Matt Pocock's `to-tickets` tracer-bullet
  review from
  `/home/devuser/workspace/skills/skills/engineering/to-tickets/SKILL.md` to its
  own final ticket breakdown before lifecycle start. The external skill is read
  as source precedent; its tracker-publication step is not imported.

## Testing decisions

- Test the completed candidate at the kit-content seam through the registered
  `docs-currency-workflow` conformance owner. The focused in-process test invokes
  that production owner directly against each mutated guidance subject and emits
  the diagnostic owned by the missing fact. Exercise finished-candidate omissions
  and additive contradictions through the established `fixture_bite_test.go`
  mutation table instead of the retired selected-root subprocess route.
- Prove the oracle bites with the existing workflow-anchor literal-mutation
  harness in `internal/conformance/fixture_bite_test.go`. Its table reads the
  canonical guidance source, applies one deletion, swap, or additive
  contradiction, and requires the exact diagnostic before restoring the source
  for the next case.
- Use the nearby `ticket-breakdown-step-anchor`,
  `implement-spec-status-flip-anchor`, and `terminal-repair-bound-anchor`
  fixtures as mutation and diagnostic exemplars. The existing in-process table
  is the narrower owner for these exact command and skill files, so extend it
  instead of creating another fixture family.
- During implementation, add the conformance expectation before its prose
  satisfies it, record the focused red, restore the intended prose, and record
  green. Then apply each completed literal mutation to the finished candidate,
  run its focused literal-mutation test, record the attributable red, restore
  it, and rerun the same check green.
- The project gate remains `.bench/gate.sh`. Focused conformance and canary
  commands are iteration evidence; the spec-build `promote` operation remains
  the sole full-gate and commit boundary.

### Seam diagram

    trigger: focused registered-owner literal-mutation test
        │
        ▼
    canonical kit tree ──▶ [ workflow-anchor conformance behavior ] ──▶ exact diagnostic / exit
    or one table-driven          ▲
    subject mutation             │ tests attach here: mutate one guidance fact,
                                 │ run the owning conformance behavior, restore,
                                 │ and require green

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| RS1 | 1 | The receipt's required fence is the maximum envelope in which ordinary independently-green repair reslicing runs. | workflow-guidance conformance seam | red at build: omit the envelope clause from the finished skill while preserving its plurality and union clauses; its `RequireInSection` literal-mutation row must emit the missing-envelope diagnostic | A singular fence transfer can satisfy the old prose while failing to authorize any inner split; the finished-candidate deletion proves the new envelope rule bites independently. |
| RS2 | 1 | Repair reslicing may yield one repair ticket or a reciprocal ordered chain. | workflow-guidance conformance seam | red at build: omit the one-or-chain result from the finished skill while preserving its envelope and union clauses; the skill-side literal-mutation row must emit the missing-result diagnostic | This red lives inside story 1's fence, so a skill that states only the envelope and union bound cannot pass before the command ticket exists. |
| RS3 | 1 | The union of every repair-ticket ownership fence remains inside the receipt's required fence. | workflow-guidance conformance seam | red at build: add permission for one chain ticket to escape the receipt fence while leaving the positive union clause intact; its `ForbidInSection` mutation row must emit the escaped-union diagnostic | The additive contradiction preserves plural ticketing and the positive sentence while violating only the maximum ownership bound, so require-only coverage cannot pass it. |
| RO1 | 2 | `/bench-implement-spec` points to `craft-tickets` as the repair-reslicing owner without restating the algorithm. | workflow-guidance conformance seam | red at build: omit the repair-specific `craft-tickets` pointer from the finished command while preserving its lifecycle prose; its `RequireInSection` mutation row must emit the missing-owner diagnostic | The command cannot silently become a second slicing owner or leave the algorithm owner unnamed. |
| RO2 | 2 | `/bench-implement-spec` says one repair ticket remains the common case. | workflow-guidance conformance seam | red at build: omit the common-case clause while preserving the skill pointer and plural lifecycle; its `RequireInSection` mutation row must emit the missing-common-case diagnostic | This pins the requested default without making it the only allowed result. |
| RO3 | 2 | `/bench-implement-spec` does not mandate exactly one repair ticket. | workflow-guidance conformance seam | red at build: add the old singular one-ticket mandate while leaving the positive command prose and skill intact; its `ForbidInSection` mutation row must emit the singular-mandate diagnostic | Intact positive prose cannot hide a contradictory command-local override of the skill's plural result. |
| RO4 | 2 | Every ticket in a repair chain is assigned, checkpointed, and integrated through the ordinary lifecycle, never through synthesized Git plumbing. | workflow-guidance conformance seam | red at build: three literal-mutation rows independently omit `assign`, `checkpoint`, and `integrate`; a fourth adds a chain-local raw-Git checkpoint route while preserving all public operation tokens; each must emit its exact missing-member or synthesized-plumbing diagnostic | Exercising every lifecycle member plus the token-preserving hostile route prevents “landed” from meaning an ad hoc commit or partial transition. |
| RO5 | 2 | The terminal repair ticket is the receipt's proceed condition; a landed non-terminal prefix cannot authorize refresh. | workflow-guidance conformance seam | red at build: add permission to refresh after the first repair ticket lands while leaving the terminal-only clause intact; its `ForbidInSection` mutation row must emit the premature-refresh diagnostic | The additive contradiction leaves a valid chain, receipt, and positive terminal clause in place but exposes the partial-state error this behavior must forbid. |
| RO6 | 2 | After the complete repair chain lands, the original blocked assignment refreshes with the original validated debug receipt. | workflow-guidance conformance seam | red at build: replace the original receipt or original assignment in the completed refresh prose, then run the focused literal-mutation test through the registered `docs-currency-workflow` owner for each independently failing identity | This pins both lifecycle identities so a new receipt or replacement assignment cannot silently become the authorization source. |

### Edge inventory

- **Error path:** a chain whose combined fence escapes the receipt → RS3.
- **Empty or absent input:** a validated receipt that yields no repair ticket →
  **Won't handle:** `craft-tickets` already requires a ticket breakdown and the
  decision source permits one or more repair tickets, never zero.
- **Boundary values:** exactly one repair ticket → RS2/RO2/RO3; two or more
  ordered repair tickets → RS2/RO4/RO5.
- **Malformed input:** a non-reciprocal producer/consumer edge → already covered
  by the existing `Integration surfaces:`/`Blocked by:` assignment refusal;
  this change composes that contract without changing it.
- **Interrupted or partial state:** a proper prefix of the chain has landed →
  RO5.
- **Re-run idempotency:** an interrupted refresh re-entry → already covered by
  the existing refresh-convergence guidance and lifecycle behavior; this change
  moves only the precondition for the first refresh.
- **Process-boundary lifecycle:** the original validated receipt survives the
  intervening repair assignments and is used for the eventual refresh → RO6.
- **Hostile environment:** lifecycle guidance that preserves every public
  operation token but routes one repair step through raw Git → RO4's new
  chain-section `ForbidInSection` anchor and additive chain-local mutation.
- **Composition degenerate:** the skill permits a chain but the command still
  mandates one ticket → RS2/RO3.
- **One-source drift:** the command restates the skill's reslicing algorithm →
  RO1 plus semantic review under the repository's one-source standard.

## Out of scope

- Enforcing repair-chain fences or terminality in new `bench spec build` code —
  separate runtime-policy capability, approximately 8 edits and 2 gate runs.
- Changing debug receipt validation, preservation refs, refresh replay,
  recomposition, or the eight-operation lifecycle — separate lifecycle
  capability, approximately 10 edits and 3 gate runs.
- Adopting the reviewer-supplied ticket-shaping dogfood proposal in
  `craft-synthesis` — separate learnings-sourced synthesis capability,
  approximately 4 edits and 1 gate run, pending `/bench-what-next` review.
- Draining or closing learning-journal entries — roadmap maintenance belongs to
  `/bench-what-next`, not this spec.
