# Forbid the new parent adapters in the policy census

Blocked by: none
Writes: internal/puritycensus/census.go, internal/puritycensus/census_test.go
Covers: LS36

## What to build

Verify the premise first. Read `policyImport`, `leafImport`, `helperImport`,
`ambientEffect`, `Policy`, `PolicyPackage`, `LeafPackage`, `Scan`, `Sources`,
`MustHold`, and `diagnose` in internal/puritycensus/census.go. Read
`TestCensusDiagnosesForbiddenImportAmbientEffectAndParallel` in
internal/puritycensus/census_test.go. Read
internal/worktree/landingpolicy/purity_census_test.go, which is the one caller
of `PolicyPackage` today.

The pattern `policyImport` names `os/exec`, `syscall`,
`github.com/gibbonmi/bench/internal/git`,
`github.com/gibbonmi/bench/internal/bounds`, and
`github.com/gibbonmi/bench/internal/worktree` today. It names one parent adapter
only, so a policy child under a different parent can import that parent.

Add `github.com/gibbonmi/bench/internal/intent` and
`github.com/gibbonmi/bench/internal/landing` to that one pattern. The pattern
stays the one source for the policy boundary, so the two new children and the
three existing children grade against one rule. Keep the leaf pattern
unchanged, because it already forbids every path under `internal/`.

Keep the diagnostic text unchanged. Keep the self-exempt file name and the
comment-stripping rule unchanged.

Add two cases to
`TestCensusDiagnosesForbiddenImportAmbientEffectAndParallel`. Each case drives
`diagnose` with `PolicyPackage()` over one source line. The first case imports
`internal/intent`, and the second case imports `internal/landing`. Each case
expects one `forbidden import` diagnostic that names the import path.

## Acceptance

- [ ] The policy census reports one forbidden-import diagnostic for a source that imports `internal/intent`.
- [ ] The policy census reports one forbidden-import diagnostic for a source that imports `internal/landing`.
- [ ] The policy census keeps its diagnostics for `os/exec`, `syscall`, `internal/git`, `internal/bounds`, and `internal/worktree`.
- [ ] The leaf census reports no new diagnostic, because its pattern is unchanged.
- [ ] `internal/worktree/landingpolicy` passes its census wrapper with its test logic unchanged.
- [ ] The pre-existing `internal/puritycensus` suite passes with its test logic unchanged.
- [ ] Self-probe: remove `internal/landing` from the pattern, and report the new case red.
