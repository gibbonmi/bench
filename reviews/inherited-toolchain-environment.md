# inherited-toolchain-environment review pickup

## Standards

Finding count: 1

Worst issue: P2 — the project profile records historical host evidence and duplicates the spec evidence.

- Disposition: auto-fix. Keep the current recovery rule in the project profile. Remove the dated host observation at `projects/benchkit.md:463`. Record the required phase-close evidence only in the spec.

## Spec

Finding count: 0

Worst issue: none.

## Coverage

Finding count: 4

Worst issue: P2 — four explicit edge predicates have no regression test.

- Disposition: auto-fix. Add the empty regular `go.mod` variant promised at `specs/inherited-toolchain-environment/spec.md:198` to the TE1 phase-table test.
- Disposition: auto-fix. Test positive discovery when `ENVMAN_LOAD` is absent and empty, as required at `specs/inherited-toolchain-environment/spec.md:203`.
- Disposition: auto-fix. Test the allowed symlinked discovered executable from `specs/inherited-toolchain-environment/spec.md:212`.
- Disposition: auto-fix. Test nonzero discovery that prints a valid path, as required at `specs/inherited-toolchain-environment/spec.md:215`.
