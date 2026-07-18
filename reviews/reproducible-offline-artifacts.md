# Review: reproducible offline artifacts

Reviewed `62e50d7..b5f7226` against repository standards, the approved FT83 slice 2 spec, and acceptance coverage.

## Standards

- **P1 — trailing gzip members are accepted:** `internal/releaseevidence/offline_archive.go:23` disables multistream parsing without proving EOF after the first member.
- **P2 — release target cardinality is duplicated:** `scripts/build-offline-archives.sh:61`, `scripts/build-offline-archives.sh:176`, `scripts/aggregate-native-proofs.sh:24`, `scripts/compare-artifacts.sh:19`, and `internal/releaseevidence/artifact_evidence.go:121` independently encode the four targets.
- **P2 — archive inventory is independently re-authored:** `scripts/build-offline-archives.sh:134`, `internal/releaseevidence/offline_archive.go:111`, and `scripts/smoke-offline.sh:56` encode the same inventory.
- **P2 — stale regular files can survive rebuilds:** `scripts/build-offline-archives.sh:38` carries existing output into staging while `scripts/build-offline-archives.sh:175` validates only matching archives.
- **P2 — structural budget exceeded:** `internal/releaseevidence/` has 16 source files against the 12-file budget in `.bench/structure-accept`.
- **P3 — ShellCheck SC1007:** `scripts/smoke-offline.sh:126` uses implicit empty assignments.

## Spec

- **High — reproducibility omits release-bound evidence:** `scripts/compare-artifacts.sh:16`, `scripts/build-artifacts.sh:97`, and `internal/releaseevidence/release_evidence.go:46` compare artifacts before release evidence exists.
- **High — archives lack an archive-level component manifest:** `scripts/build-offline-archives.sh:134` embeds only package manifests and `internal/releaseevidence/artifact_evidence.go:196` substitutes the platform manifest digest.
- **High — clean-room smoke makes approved evidence optional:** `scripts/smoke-offline.sh:72`; neither `.github/workflows/native-runtime.yml:52` nor `.github/workflows/release.yml:56` supplies it.
- **High — denied egress, fallback prevention, and direct cleanup are not proved:** `scripts/smoke-offline.sh:84`, `scripts/smoke-offline.sh:95`, `scripts/smoke-offline.sh:126`, and `scripts/smoke-offline.sh:34` rely on configuration and whole-tree cleanup.
- **High — registry order/integrity is tautological:** `scripts/smoke-offline.sh:107` pre-seeds files and `scripts/smoke-offline.sh:127` compares a self-authored order rather than observed uploads.
- **Medium — native proof is incomplete:** `scripts/native-proof.sh:46` and `scripts/native-proof.sh:61` do not prove all required rebuilt-binary operations and Darwin stripping.
- **Medium — verify evidence is finalized before native proofs:** `internal/releaseevidence/artifact_proofs.go:41`, `internal/releaseevidence/release_evidence.go:149`, and `.github/workflows/native-runtime.yml:82` allow verify evidence without the aggregate.
- **Medium — duplicate npm tar entries are accepted:** `scripts/build-offline-archives.sh:90` does not track normalized names before extraction.
- **Medium — interruption tests do not preserve prior bytes:** `internal/contract/surface/artifact/artifact_concurrency_test.go:52` and `internal/contract/surface/artifact/artifact_fixture_test.go:229` assert counts rather than prior-generation digests.
- **Medium — ambient reproducibility perturbations are incomplete:** `internal/preflight/release_index_test.go:104` omits HOME, cache, temp-root, and build-directory variation.
- **Medium — OFFLINE.md lacks exact runnable instructions:** generated content at `scripts/build-offline-archives.sh:148` describes actions without exact verification/install/seed commands.

## Coverage

- **High — acceptance rows cite inspections rather than runnable red failures:** `specs/reproducible-offline-artifacts.md:168`, `specs/reproducible-offline-artifacts.md:171`, and rows beginning at line 221 need named mutation commands and recorded red output.
- **High — comparator and reproducibility record are not executed negatively:** `internal/contract/surface/artifact/artifact_test.go:33`, `internal/preflight/integration_fixture_test.go:164`, and `internal/conformance/native_workflow_test.go:151` can accept fabricated `Match: true` evidence.
- **High — native proof scripts lack executable failure coverage:** `internal/preflight/integration_fixture_test.go:182` fabricates proof JSON and `internal/conformance/package_core_checks_test.go:313` checks only text anchors.
- **High — registry upload order and stored bytes are unobservable:** `scripts/offline-registry.mjs:47` records GET only, so `scripts/smoke-offline.sh:127` cannot detect upload mutations.
- **High — no negative denied-egress/fallback signal:** `scripts/smoke-offline.sh:84` and `scripts/smoke-offline.sh:95` lack an egress sentinel and missing/corrupt/repair/rebuild mutations.
- **Medium — malformed archive coverage is non-authoritative:** `internal/contract/artifact.go:51` skips special members and overwrites duplicates, and `internal/contract/surface/artifact/artifact_offline_test.go:56` depends on it.
- **Medium — workflow checks are presence-only:** `internal/conformance/package_core_checks_test.go:313` does not verify job/step relationships or evidence upload paths.
- **Medium — missing release evidence has no smoke red signal:** `scripts/smoke-offline.sh:72` treats it as optional.

The repair pass should consolidate overlapping findings around canonical release metadata, fail-closed parsing/smoke behavior, observable native/registry proofs, and mutation-backed acceptance coverage.
