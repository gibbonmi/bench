## Standards

3 findings. Worst issue: release evidence can attest build and install flags that
are independently authored from the commands actually used.

1. **High — hard violation: the supported-target count has three production
   sources.** `AGENTS.md`'s code standard forbids duplicated production policy
   and derived counts. The canonical rows live in `scripts/platforms.json:1-6`,
   while `scripts/build-release-evidence.mjs:16` and
   `internal/releaseevidence/artifact_evidence.go:119-120` independently require
   exactly four. Derive both validations from the matrix or one policy source.
2. **High — hard violation: recorded toolchain flags are independent of executed
   flags.** The same one-source-per-fact rule applies to production policy.
   `internal/releaseevidence/requirements.json:4-6` declares Go/npm build, pack,
   and install flags, while `scripts/go-build.sh:30`,
   `scripts/build-artifacts.sh:49,52`, and `scripts/smoke-artifacts.sh:31`
   hardcode execution separately. `internal/releaseevidence/release_evidence.go:258-270`
   records the registry strings without observing the commands. Derive execution
   from the registry or capture the actual invoked flags.
3. **High — hard violation: release inputs are inventoried twice.** The same rule
   explicitly forbids duplicated executable registries.
   `internal/releaseevidence/release_requirements.go:28-37` lists release inputs,
   while `internal/releaseevidence/registry.json:5-10` separately lists phase
   inputs. `internal/releaseevidence/requirement_inspection.go:67-73` and
   `internal/releaseevidence/evidence_fingerprint.go:15-20` consume only the
   former, so a phase-input addition can leave index binding stale. Establish one
   registry and derive both projections.

## Spec

5 findings. Worst issue: source drift can promote trustworthy-looking evidence
bound to the wrong commit.

1. **High — source drift can promote an index claiming the wrong commit.** The
   spec requires source drift to preserve the prior generation and inputs to be
   revalidated immediately before promotion
   (`specs/governed-release-evidence.md:131-136`). Identity is captured once at
   `internal/preflight/command.go:275-290`; the final fingerprint at
   `internal/releaseevidence/evidence_fingerprint.go:13-35` omits both `HEAD` and
   `RunEvidence.Identity`. If `HEAD` advances from A to a byte-identical commit B
   after capture, promotion succeeds while the index still claims A. Fingerprint
   and re-resolve source identity immediately before promotion.
2. **High — focused diagnostics replace the prior complete indexed generation.**
   Story 5 requires complete generations and says partial or mixed evidence is
   never authoritative (`specs/governed-release-evidence.md:62-67`). Focused
   finalization promotes only one phase record plus a manifest at
   `internal/releaseevidence/release_evidence.go:83-93`; atomic replacement at
   `internal/releaseevidence/evidence_promotion.go:57-95` deletes the prior index
   and checksums. Tests explicitly accept both incomplete focused output and prior
   replacement at `internal/preflight/focused_publish_test.go:36-38` and
   `internal/preflight/preflight_test.go:27-43`. Keep focused diagnostics separate
   or preserve the authoritative full generation.
3. **Medium — the producer-record failure matrix is implemented but not
   acceptance-proven.** Coverage row 220 requires distinct red fixtures for
   unknown versions, duplicate keys, mismatched identity, and digest mismatch
   (`specs/governed-release-evidence.md:220`). Validators exist at
   `internal/releaseevidence/release_requirements.go:145-161`, but producer tests
   at `internal/preflight/preflight_test.go:98-169` cover only missing records.
   Add built-command producer mutations asserting the distinct red status/reason.
4. **Medium — missing Node/npm lacks prior-generation preservation proof.** Row
   221 requires missing tools to fail safely
   (`specs/governed-release-evidence.md:221`), and story 5 requires preservation.
   `internal/preflight/release_index_test.go:52-69` observes only installed
   Go/Node/npm. After a green run, remove Node and npm from `PATH` separately and
   assert prompt failure plus byte-identical prior evidence.
5. **Medium — control bytes in tar member paths lack an acceptance fixture.** Row
   221 requires control bytes to fail safely
   (`specs/governed-release-evidence.md:221`), while the direct hostile archive
   test covers only oversize members at
   `internal/preflight/archive_hostile_test.go:11-31`. Inject a tar member path
   containing ESC or BEL with a consistent manifest, then assert non-blocking red
   and prior-generation preservation.

## Coverage

0 findings. Worst issue: none.

The generic edge inventory and the project hostile-input checklist are fully
declared at `specs/governed-release-evidence.md:237-260`. Live registry/runner
behavior and external FT71/FT87/FT88 semantics are explicit Won't-handle entries
at `specs/governed-release-evidence.md:261-266`. The three initially suspected
test gaps are mapped by rows 220-221 and are therefore Spec findings above, not
unmapped Coverage findings.
