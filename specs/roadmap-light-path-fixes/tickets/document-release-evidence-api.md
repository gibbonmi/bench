# Document release-evidence API contracts

Blocked by: none
Writes: internal/releaseevidence/release_evidence.go, internal/releaseevidence/evidence_promotion.go, internal/releaseevidence/release_index.go, internal/releaseevidence/types.go
Covers: LF16

## What to build

Document FinalizeEvidence, PromoteEvidence, PromoteEvidenceFiles,
TerminalStatus, RequirementsRegistry, PackageEvidenceRegistry, Requirements,
PhaseNames, PhaseDefinitionFor, and PhaseSummaries. State caller obligations,
results, and failure posture without reviewer narration.

## Acceptance

- [ ] All ten named functions state their caller-facing contracts.
- [ ] Promotion comments describe atomicity and failure posture.
- [ ] Registry comments identify the canonical source of requirements.
