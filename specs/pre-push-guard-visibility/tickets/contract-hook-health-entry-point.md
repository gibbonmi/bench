# Contract the hook-health entry point

Blocked by: none
Ownership fence: `internal/adopt/link_hook.go`, `internal/adopt/link_hook_test.go`
Integration surfaces: hook-health API→`internal/adopt/link_hook.go` + CE1; exported-symbol structural contract→`internal/adopt/link_hook_test.go` + CE1
Contracts: the complete hook-health record crosses `internal/adopt/link_hook.go`→its package structural contract, asserted by CE1 against the real exported symbol set
Closure: CE1/unexported-legacy-classifier

## What to build

Remove the unused exported `ClassifyPrePush` compatibility view so `InspectPrePush` is the only hook-health entry point. Keep the implementation and its structural assertion together: an extraction-only cut strands the package contract red because it must prove that no exported partial classifier remains.

## Acceptance

- [ ] [CE1] `InspectPrePush` is the sole exported hook-health computation and `ClassifyPrePush` is absent from the production API.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CE1/unexported-legacy-classifier | retain or restore exported `ClassifyPrePush` | link-hook package structural test | enumerate exported declarations in `link_hook.go`, expect the legacy classifier assertion red |
