# Stable-owner worktree landing

Status: staged

Decision source: reviewer-confirmed conversation on 2026-08-24

Verification log: 2 iterations to accept — three interface reviews rejected candidate-owned publication and selected a stable promotion owner

## Problem

The landing command can run under an executable from the assignment worktree.
That executable can contain behavior that the destination does not yet trust.

The current freshness repair rebuilds the primary executable and repeats the command.
This can replace a reviewed fix with the old behavior that caused the refusal.

The command does not explain that a landing-owner change cannot authorize itself.
Operators can move unrelated residue or invoke different binaries without a useful result.

## Solution

The installed distribution supplies one stable promotion owner for the complete landing.
The wrapper authenticates this broker through the installation manifest.
It selects the broker independently of repository state and inherited run-binary state.

The owner resolves the primary destination and the registered assignment as separate subjects.
It composes an exact prospective tree in private storage.
It runs the baseline gate policy against that tree.
Only the stable owner can publish the destination ref.

Candidate landing code never runs during its own promotion.
A broker source change takes effect only after the release or repair path installs it.
Normal landing can publish the source change but cannot replace its authority.

## User stories

### Select the stable owner

Line: gpt-5.6-sol / high.

This group removes executable selection from the candidate subject.

1. As an operator, I want landing to use the installed promotion broker, so that repository code cannot authorize its own publication.
2. As an operator, I want current-directory changes to leave owner selection unchanged, so that source and destination roots cannot swap.
3. As a reviewer, I want inherited routing overrides refused, so that a caller cannot bypass the promotion owner.
4. As an operator, I want one physical owner path for the complete command, so that landing cannot rebuild and re-execute itself.

### Grade the exact prospective tree

Line: gpt-5.6-sol / high.

This group separates the promotion owner from the tree that the gate grades.

5. As a reviewer, I want the assignment identity bound before composition, so that a changed request cannot enter the prospective tree.
6. As a reviewer, I want the review base bound before composition, so that a changed base cannot enter the prospective tree.
7. As a reviewer, I want the source tip bound before composition, so that a changed tip cannot enter the prospective tree.
8. As a reviewer, I want the source tree fingerprint bound before composition, so that uncommitted source changes cannot enter the prospective tree.
9. As an operator, I want the gate binary built from the prospective tree, so that the gate grades the exact publication subject.
10. As a reviewer, I want baseline gate policy to select the phases, so that candidate policy cannot omit its own checks.

### Publish and recover under one authority

Line: gpt-5.6-sol / high.

This group keeps the deterministic verdict and destination update under one owner.

11. As an operator, I want a failed prospective gate to leave the destination unchanged, so that red evidence cannot publish.
12. As an operator, I want a destination race to fail the final compare-and-swap, so that a stale candidate cannot overwrite new work.
13. As an operator, I want an identical prospective retry to reuse green evidence, so that recovery does not add a full gate.
14. As an operator, I want post-publication recovery to resume without another publication, so that retries are idempotent.
15. As an operator, I want temporary binaries removed after completion, so that landing does not leave executable residue.
16. As an operator, I want a broker-changing diff to name the install step, so that I do not expect source publication to replace authority.
17. As a reviewer, I want repository binary replacement ignored, so that a forged binary and seal cannot become the promotion owner.

## Implementation decisions

The installed wrapper has a dedicated route for `bench worktree land`.
That route selects a broker from the installed distribution, not the repository.
It does not honor `BENCH_KIT`, `BENCH_RUN_BINARY`, or `BENCH_WRAPPER` for this public command.

The installation manifest binds the broker path, version, platform, and executable digest.
The package installation or repair owner publishes the manifest and broker together.
The design trusts that installer and the owner-only installation directory.

The promotion owner resolves the primary checkout through the shared Git directory.
It requires the destination to be the attached default-branch checkout.
It resolves the source only through the registered assignment and request.

The promotion owner performs all assignment identity checks before composition.
It repeats the source fingerprint check before the gate starts.

The promotion owner materializes the exact prospective tree in private temporary storage.
It uses the trusted baseline build recipe and phase schedule.
The prospective tree supplies the code and tests under examination.

The promotion owner builds a gate executable from the prospective tree.
That executable is the gate subject, not the publication authority.
The gate receives no destination update capability from Bench.

The owner keys green evidence to the prospective tree and baseline runner identity.
A candidate gate script cannot replace the baseline phase schedule.

The stable owner validates green evidence and performs the destination compare-and-swap.
It then advances the marker, reconciles the checkout, and releases the assignment.
Resume operations remain under the same owner.

The current broker lands a broker-changing source diff under the current contract.
The new broker becomes active only after the release or repair owner installs it.
This explicit install step is a reviewer-visible bootstrap rule.

## Testing decisions

The highest seam is the public `bench worktree land` process journey.
The journey invokes the command from the primary tree, the source tree, and nested directories.

Mutation probes replace inherited routing variables and candidate gate policy.
The probes also change each assignment identity dimension independently.

The gate journey records both executable identities.
It proves that the stable owner and prospective gate binary have different roles.

### Seam diagram

    installed wrapper + manifest
          │
          ▼
    [ stable promotion owner ] ──▶ private prospective tree ──▶ gate subject
          │                                                    │
          └──────────────── validates evidence ────────────────┘
          │
          ▼
    destination compare-and-swap

    tests attach at the public process and destination ref

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SOL01 | 1 | A candidate landing marker never runs during its own promotion. | public landing process | Candidate-owned publication executes the marker. |
| SOL02 | 2 | Primary, source, nested, and outside-repository invocations select the same owner. | wrapper process | A current-directory resolver selects different binaries. |
| SOL03 | 3 | Each inherited routing override makes public landing refuse before repository reads. | wrapper process | The generic route accepts at least one override. |
| SOL04 | 4 | One owner process identity remains active through publication and release. | public landing process | A rebuild or re-execution changes the recorded identity. |
| SOL05 | 5 | A changed request invalidates the landing before composition. | assignment identity seam | A comparator that omits the request accepts the mutation. |
| SOL06 | 6 | A changed review base invalidates the landing before composition. | assignment identity seam | A comparator that omits the base accepts the mutation. |
| SOL07 | 7 | A changed source tip invalidates the landing before composition. | assignment identity seam | A comparator that omits the tip accepts the mutation. |
| SOL08 | 8 | A changed source fingerprint invalidates the landing before the gate. | assignment identity seam | A commit-only comparison accepts dirty source content. |
| SOL09 | 9 | The gate records an executable built from the exact prospective tree. | gate process | A source-tip binary records a different source digest. |
| SOL10 | 10 | A candidate phase omission does not change the baseline phase set. | prospective gate process | Candidate-owned policy drops the planted phase. |
| SOL11 | 11 | A red prospective gate leaves the destination ref and marker unchanged. | public landing process | A candidate-owned verdict can publish before validation. |
| SOL12 | 12 | A destination change before publication makes the compare-and-swap refuse. | publication seam | A non-atomic update overwrites the competing change. |
| SOL13 | 13 | Evidence reuse requires the same prospective tree and baseline runner identity. | gate evidence seam | A partial key reuses evidence after one subject changes. |
| SOL14 | 14 | Marker, reconcile, and release failures resume without a second publication. | public resume process | A first-run retry publishes the same change twice. |
| SOL15 | 15 | Temporary prospective binaries are absent after success and every refusal. | public landing process | An incomplete cleanup leaves at least one artifact. |
| SOL16 | 16 | A broker-changing landing reports the required release or repair action. | public landing process | A silent install dependency recreates the bootstrap confusion. |
| SOL17 | 17 | Replacing the primary binary and adjacent seal does not execute their planted marker. | public landing process | A repository-path trust root executes the forged owner. |

### Edge inventory

- The invocation starts in the primary root, source root, a nested directory, or outside a repository.
- An inherited routing variable is absent, empty, relative, non-physical, or points to another executable.
- A root path contains spaces, glob characters, a newline, or a control byte.
- The source and destination are equal, swapped, or from different shared Git directories.
- The source is absent, moved, outside the pool, or no longer a registered assignment.
- The request, review base, source tip, or source fingerprint changes independently.
- The destination changes before composition, during the gate, or before publication.
- The gate exits red, receives an interrupt, or leaves partial temporary state.
- Marker, reconcile, and release fail independently after publication.
- The candidate changes its landing code, gate script, or build inputs.

**Won't handle** hostile same-user test code that directly edits Git refs — the existing gate threat model remains the surviving in-scope authority.

**Won't handle** a promotion-owner hot swap during one process — the next command is the surviving in-scope activation point.

## Ownership fences

### Ticket 01

- `bin/bench.sh`
- `bin/bench-postinstall.sh`
- `internal/adopt/`
- `internal/runbinary/`
- `internal/worktree/land.go`
- `internal/worktree/land_freshness_test.go`
- `internal/worktree/land_identity.go`
- `internal/worktree/land_identity_test.go`
- `internal/systemtest/`

### Ticket 02

- `internal/landing/`
- `internal/gate/`
- `internal/worktree/land.go`
- `internal/worktree/land_resume.go`
- `internal/worktree/land_freshness_test.go`
- `internal/systemtest/`
- `CHANGELOG.md`
- `projects/benchkit.md`

The tickets are serial because both tickets edit the public landing owner.

## Out of scope

- The spec-build review and gate cadence is blocked by FT130 and FT173: at least 30 edits and 3 gate runs.
- A changed-package gate is FT215: at least 12 edits and 3 gate runs.
- Legacy cleanup-pending recovery is a separate lifecycle repair: at least 4 edits and 2 gate runs.
- Capture guidance reconciliation is a separate conformance repair: at least 5 edits and 1 gate run.
- Cross-user process isolation is a separate security capability: at least 20 edits and 4 gate runs.

## Further notes

The 2026-08-24 debug run proved the executable loop.
The candidate fix passed the full gate.
The primary invocation rebuilt old behavior and repeated the old refusal.
The final repair used an explicit fast-forward because the old owner could not consume its own policy fix.
