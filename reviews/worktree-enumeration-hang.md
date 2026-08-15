## Standards

Finding count: 0. Worst: none.

The frozen candidate has no Standards findings.

## Spec

Finding count: 1. Worst: the scanner accepts a non-regular `bench-lease` despite
the approved direct-entry predicate.

- `ask-user` — `specs/worktree-enumeration-hang/spec.md:61-65` requires every
  direct admin entry to be regular or a directory, but the existing
  `bench-lease` exemption in `internal/git/git.go:233-237` preserves lifecycle
  uncertainty. Removing it requires propagating the typed scan refusal through
  the explicitly out-of-scope mutating ownership path. Decide whether the
  enumeration contract removes the exemption or the spec gains the lifecycle
  exception.

## Coverage

Finding count: 1. Worst: the undocumented `bench-lease` exception.

- `ask-user` — The `bench-lease` symlink/FIFO allowance remains unreviewed;
  it is the same contract conflict recorded by the Spec axis above.
