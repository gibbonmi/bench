## Standards

4 findings. Worst issue: duplicated production knowledge violates the repository's
one-source-per-fact rule.

1. **High — hard violation: package evidence is authored twice.** `AGENTS.md` says
   that two derivations of one fact, including production policy and executable
   registries, must collapse to one source. The eight governance path/schema pairs
   in `internal/preflight/requirements.json:5-12` are repeated by the requirement
   records at `internal/preflight/requirements.json:28-35`. Make the records own
   package inclusion/mode and derive the package evidence projection.
2. **High — hard violation: registry-owned schema and toolchain facts are
   re-authored in validators.** The registry owns toolchains and component-manifest
   fields at `internal/preflight/requirements.json:14-25`, while
   `internal/preflight/release_requirements.go:49-53,119-123` independently requires
   their counts/names and `scripts/build-release-evidence.mjs:16` independently
   requires every field count. The same `AGENTS.md` rule explicitly forbids derived
   counts outside independent omission/mutation expectations. Consumers should
   validate general shape from the registry; independent tests/canaries should own
   the drift expectations.
3. **Medium — hard violation: the package version has two production sources.**
   `package.json:3` owns version `0.2.0`, but
   `governance/sbom.spdx.json:12` independently commits the same version. The builder
   overwrites it at `scripts/build-release-evidence.mjs:145-152`, while source
   validation calls `validateSPDXDocument(data, "", "")` at
   `internal/preflight/release_requirements.go:187-190`, so a version bump leaves a
   stale governed input. Generate the value from the canonical package version or
   use a non-versioned template sentinel.
4. **Medium — judgment call (Divergent Change): release evidence is an ungrouped
   module.** FT83 grows `internal/preflight/` from 8 to 21 Go files and mixes archive
   inspection (`internal/preflight/artifact_evidence.go:75`), governance schemas
   (`internal/preflight/release_requirements.go:202`), fingerprinting, and OS-specific
   atomic promotion behind the façade at
   `internal/preflight/release_evidence.go:76`. `bench structure` reports the
   directory over its 12-file budget; the seam guidance says a crowded directory is
   an ungrouped module. Move the release-evidence implementation behind a dedicated
   package and leave orchestration at the façade, or seek a reviewer-owned directory
   grant if the package is genuinely cohesive.

## Spec

3 findings. Worst issue: the drift check can authorize an index assembled from one
artifact generation while validating a later generation.

1. **P1 — drift detection can bless mixed artifact generations.** The coverage map
   requires that changed inputs produce drift and never promote an index for mixed
   bytes (`specs/governed-release-evidence.md:222`). Artifact set A is read at
   `internal/preflight/artifact_evidence.go:112`; toolchain commands then run at
   `internal/preflight/release_evidence.go:206`; only afterward is the baseline
   fingerprint taken at `internal/preflight/release_evidence.go:247`. If
   `dist/artifacts` is atomically replaced by set B during toolchain observation,
   both that baseline and the final comparison at
   `internal/preflight/release_evidence.go:107-112` see B, while the index still
   describes A. Bind the fingerprint to the same reads used for assembly and add a
   pre-baseline replacement regression.
2. **P1 — component inventories can authenticate only a 256 MiB prefix.** The map
   requires every manifest to match the size/SHA-256 inventory of bytes actually
   present (`specs/governed-release-evidence.md:211`). Tar members are read through
   an unchecked `io.LimitReader` at
   `internal/preflight/artifact_evidence.go:165`; validation compares manifest size
   and digest only to that truncated buffer at
   `internal/preflight/artifact_evidence.go:251-253`. A larger regular member whose
   manifest describes its first 256 MiB is accepted with trailing bytes unbound.
   Reject oversize members or detect an extra byte, with a hostile-tar regression.
3. **P2 — abandoned promotion stages survive reruns.** The map requires abandoned
   stage reruns to leave exactly one complete generation
   (`specs/governed-release-evidence.md:224`). Promotion creates
   `.preflight-stage-*` at `internal/preflight/evidence_promotion.go:26`; SIGKILL
   bypasses deferred cleanup and reruns do not discover old stages. The SIGKILL test
   checks only the canonical target at `internal/preflight/evidence_test.go:81-90`,
   not the complete `dist` generation set. Add ownership-safe abandoned-stage
   recovery and assert no residual stages after rerun.

## Coverage

1 finding. Worst issue: symlinked release-preflight invocation is broken and
untested.

1. **Symlinked entry-script invocation resolves the wrong repository.** The hostile
   checklist requires invocation through a symlink (`projects/benchkit.md:110`), but
   `scripts/release-preflight.sh:3` derives the root from unresolved `BASH_SOURCE`.
   An external symlink therefore looks for `/tmp/scripts/go-build.sh` and exits 127
   before preflight. The hostile coverage row at
   `specs/governed-release-evidence.md:221` mentions symlink escape but has no
   entry-script-symlink test; `internal/preflight/integration_test.go:71` invokes
   only the real script. Add relative and absolute symlink integration cases that
   require the same exit/evidence as the real path.
