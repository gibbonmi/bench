# Restrict upgrade hook absence

Blocked by: none
Ownership fence: `internal/adopt/link_hook.go`, `internal/contract/surface/upgrade_test.go`
Integration surfaces: configured-path absence predicate→`internal/adopt/link_hook.go` + RA1; dependent prospective-count ticket→align-upgrade-hook-plan.md; consumer upgrade contract→`internal/contract/surface/upgrade_test.go` + RA1
Contracts: an absent effective configured hook path crosses `internal/adopt/link_hook.go`→the upgrade surface contract, asserted by RA1 against the real path predicate
Closure: RA1/not-found-only-absence

## What to build

Limit configured-hook absence to the precise `os.IsNotExist` result before upgrade planning considers a refresh. A path-predicate-only ticket is independently green: its non-directory-parent fixture reds the current predicate without needing prospective template rendering, and `align-upgrade-hook-plan.md` is blocked on this ticket so its rendered-byte work cannot broaden absence again.

## Acceptance

- [ ] [RA1] With an unequal version and `core.hooksPath` whose effective `pre-push` path returns a non-not-found error, `bench upgrade --check` does not count a hook refresh.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RA1/not-found-only-absence | treat every `Lstat` error as an absent configured hook | upgrade surface contract | configure `core.hooksPath` through a non-directory parent, run unequal-version `upgrade --check`, expect the impossible change count red |
