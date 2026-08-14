# Review: parallel-session-landings

Exact reviewed source: `be5ec93e911105678876e300e84c38d4727ed656..0ece09426e1670601daba185940020bbd98990c5`.

Raw findings: Standards 5, Spec 4, Coverage 4. Actionable findings: 11. De-duplicated repair targets: 9.

## Standards

Finding count: 5 raw, 4 actionable. Worst issue: high, closed as `no-op`; worst actionable issue: medium.

- **Medium — `auto-fix`: centralize strict porcelain-v1 `-z` framing.** `internal/landing/landing.go:358-387` parses the same Git framing owned by `internal/git/git.go:437-449`. Extend the Git owner for the strict/filtering case instead of retaining a second parser. Rule: `AGENTS.md:27-37` keeps production parsers single-sourced.
- **Medium — `auto-fix`: derive landing grammar dispatch and required flags from grammar metadata.** `internal/worktree/land.go:25-51` declares the grammars, while `internal/worktree/land.go:73,123-148` repeats flag-value and count knowledge; `internal/worktree/reauthorize.go:50` repeats another required-count fact. Rule: `AGENTS.md:27-37` rejects duplicate parsers and derived counts.
- **Medium — `auto-fix`: use one request-digest owner.** `internal/intent/assignment.go:446`, `internal/worktree/ownership.go:38`, and `internal/worktree/worktree.go:24` independently derive the same SHA-256 request digest. Rule: `AGENTS.md:27-37` requires one source per fact.
- **Medium — `auto-fix`: move resume spec-path derivation to the spec owner.** `internal/worktree/land.go:288-293` reimplements the slug/path rule owned by `internal/spec/spec.go:137-169`. Rule: `AGENTS.md:27-37` keeps production parsers single-sourced.

## Spec

Finding count: 4 raw, 3 actionable. Worst issue: high.

- **High — `auto-fix`: refuse every caller-owned destination state before resume reconciliation.** The invoking checkout must be clean across index and tracked worktree (`specs/parallel-session-landings/spec.md:150-156,215-221`), but `internal/worktree/land.go:206-228` checks only unstaged diff before `internal/worktree/land.go:407-411` runs `reset --hard`. Staged, untracked collision, ignored, and nested-repository states need fail-closed coverage before reconciliation.
- **High — `ask-user`: authenticate resume against the exact prior landed-but-incomplete publication.** The form accepts only the exact reported commit (`specs/parallel-session-landings/spec.md:155-157`) and must correspond to the retained assignment (`specs/parallel-session-landings/spec.md:281-285`). Active resume currently authenticates structure plus assignment but no prior landing authority (`internal/worktree/land.go:160-172,252-285`); completed resume reuses a generic release receipt that is not bound to the published commit or source tip (`internal/worktree/land.go:194-203`). The spec forbids a new journal or receipt schema (`specs/parallel-session-landings/spec.md:111-116`), so the authority source requires reviewer choice.
- **Medium — `auto-fix`: distinguish destination movement from infrastructure ref-lock failure.** The edge inventory requires CAS refusal only when the destination actually moved (`specs/parallel-session-landings/spec.md:459-463`), while `internal/landing/landing.go:312-314` labels every `update-ref` failure a compare-and-swap refusal. Re-read the destination and preserve infrastructure classification when it is unchanged.

## Coverage

Finding count: 4 raw, 4 actionable. Worst issue: high.

- **High — `auto-fix`: add destructive-state resume cases.** PL25/PL30 and the hostile destructive-worktree class have no staged, untracked-collision, ignored, or nested destination case; existing resume reconciliation coverage at `internal/worktree/land_test.go:234-258` is clean-only. This shares the Spec repair target for destination cleanliness.
- **High — `ask-user`: bind completed and active resume to the same landing identity.** A terminal `worktree-release` receipt authenticates request/path/assignment cleanup but carries neither published commit nor source tip (`internal/intent/assignment.go:86-106`; `internal/worktree/lifecycle.go:272-274`). This shares the Spec authority target and needs wrong-published/wrong-source controls.
- **High — `auto-fix`: movement-check explicit-base preflight.** Both explicit surfaces must resolve base and tip once per movement-checked attempt (`specs/parallel-session-landings/spec.md:130-141`), but `internal/preflight/gather.go:79-105` can combine an earlier tip with later status/path reads without revalidation. Add HEAD, index, tracked-worktree, and untracked drift controls.
- **Medium — `auto-fix`: prove PL15/PL16 across two public processes.** PL15 and PL16 require a two-process winner/loser and rerun journey (`specs/parallel-session-landings/spec.md:401-402,456-458`); `internal/landing/landing_test.go:521-570` covers only an injected in-process updater. Add the public process-boundary race and exact gate/ref/marker assertions.
