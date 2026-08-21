# Allow runtime-root residue in every lifecycle allowance

Blocked by: none
Writes: internal/worktree/build_outputs.go, internal/worktree/subshell.go, internal/worktree/eligibility.go, internal/worktree/eligibility_test.go, internal/worktree/land_test.go, internal/landing/landing.go, internal/landing/landing_test.go
Line: opus / high — the allowance feeds destructive eligibility.

## What to build

The gate's `.logs/` runtime records stop counting as blocking ignored residue
anywhere in the worktree lifecycle. `ignoredWithinBuildOutputs` — the variant
whose *additional* predicate is always false — is deleted; every
declared-outputs allowance composes `landing.RuntimeIgnoredPath`, which stays
the one source of the runtime-root fact. A release of a source whose only
ignored residue is gate records completes; any residue outside the runtime
root and the declared build outputs still retains fail-closed; the allowance
holds when `.bench/build-outputs.json` is absent or empty. Flagged
consequence carried by this ticket (reviewer-approved at spec sign-off):
`bench worktree clean` then plans `discard-remove` over runtime-root records
without `--discard-ignored`; the plan/apply fingerprint gate still stands.

## Acceptance

- [ ] A release of a source whose only ignored residue is runtime-root gate records completes and removes the worktree (covers LR12).
- [ ] A release of a source with one ignored path outside the runtime root and outside declared outputs retains fail-closed (covers LR13).
- [ ] The residue set that passes the destination check also passes release eligibility and clean eligibility (covers LR14).
- [ ] With `.bench/build-outputs.json` absent, a runtime-root-only source still releases; the empty-declaration variant behaves the same (covers LR15).
