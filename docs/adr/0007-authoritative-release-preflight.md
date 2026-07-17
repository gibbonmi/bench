# Release publication requires one authoritative preflight

Bench release verification is owned by one Go-core preflight with an ordered phase registry, atomic machine-readable evidence, and two modes: `verify` runs the complete source and artifact checks, while `publish` is its strict superset with release identity, ancestry, and changelog policy. Pull-request, default-branch, native-runner, and tag workflows are thin callers of that same owner because npm publication is immutable and no workflow-local reimplementation may bypass or reinterpret its verdict.

## Consequences

Only exact stable-version tags whose package, binary, changelog, commit, and patch toolchain identities agree may publish, and the tagged commit must be on the main release line. Vulnerability findings require tracked, expiring, currently used exceptions; focused native smoke records remain diagnostic and cannot authorize publication; complete preflight evidence is replaced old-or-new so interruption never exposes partial authority. Hosted artifact retention and live registry behavior remain external proofs, while the repository gate and canaries enforce phase completeness, workflow dependencies, and planted release-policy failures.
