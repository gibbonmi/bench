# FT83 slice 2 terminal review pickup

Pinned diff: `390c419..4f19892`

## Standards

4 findings. Worst: duplicated release-plan ownership can make builders and
authorization disagree about the release set.

- **High — release-plan knowledge has multiple owners.** The one-source-per-fact
  rule forbids duplicate parsers and derived counts (`AGENTS.md:24-34`), but
  `scripts/release-plan.mjs:7-44`,
  `internal/releaseevidence/release_plan.go:24-88`, and
  `internal/releaseevidence/artifact_evidence.go:121-126` independently validate
  or derive the same release plan and artifact names.
- **High — the release-bound evidence inventory is triplicated production
  policy.** `scripts/compare-artifacts.sh:44-52`,
  `scripts/compare-artifacts.sh:62-73`, and
  `internal/releaseevidence/artifact_proofs.go:56-69` reconstruct the same
  inventory despite `AGENTS.md:29-34` requiring production policy and executable
  registries to remain single-sourced.
- **Medium — native verification orchestration is duplicated across workflows.**
  `.github/workflows/native-runtime.yml:13-107` and
  `.github/workflows/release.yml:13-118` independently encode the same build,
  native-proof, evidence, and smoke chain, contrary to `AGENTS.md:24-34`.
- **Medium judgment call — archive readers duplicate secure traversal.**
  `internal/releaseevidence/package_artifact.go:16-79` and
  `internal/releaseevidence/offline_archive.go:15-103` repeat gzip/tar traversal,
  limits, special-member rejection, duplicate detection, and bounded reads; this
  is Duplicated Code under `.agents/skills/bench-craft-review/SKILL.md:51-52`.

## Spec

5 findings. Worst: valid release evidence is rejected before any offline journey
runs.

- **High — valid checksums are always rejected.** The spec requires supplied
  index/checksum verification (`specs/reproducible-offline-artifacts.md:111`).
  `internal/releaseevidence/release_index.go:114` emits `digest  name`, while
  `scripts/smoke-offline.sh:92` builds a digest-to-name map and queries it by
  filename.
- **High — self-contained verification requires unavailable artifacts.** A
  consumer authenticates one target archive
  (`specs/reproducible-offline-artifacts.md:37`), but `SHA256SUMS` covers all nine
  deliverables (`:136`) and `scripts/assemble-offline-archive.mjs:34` instructs an
  unfiltered `sha256sum -c SHA256SUMS`.
- **High — public or undeclared egress is not denied.** Only the provisioned
  loopback registry is permitted (`specs/reproducible-offline-artifacts.md:131`),
  but `scripts/smoke-offline.sh:161` supplies no observing sentinel, namespace,
  or firewall that can detect arbitrary outbound connections.
- **Medium — the internal-registry journey omits the second read-only operation.**
  Acceptance requires both operations (`specs/reproducible-offline-artifacts.md:249`),
  while `scripts/smoke-offline.sh:165` runs only `version` before uninstalling.
- **Medium — M7 coverage is falsely closed.** The map promises distinct corrupt
  bytes, wrong target, checksum, fetch/rebuild/repair, and residue mutations
  (`specs/reproducible-offline-artifacts.md:247`), but
  `internal/contract/surface/artifact/reproducibility_test.go:95` only supplies
  missing evidence and `internal/conformance/native_workflow_test.go:161` relies
  on source-string sentinels rather than real-process behavior.

## Coverage

0 findings. Worst: none.

The edge inventory and hostile-input checklist were traced through archive
mutation, reproducibility, hostile-root, native-proof, and workflow-edge tests;
remaining states are mapped or explicitly excluded by the spec.
