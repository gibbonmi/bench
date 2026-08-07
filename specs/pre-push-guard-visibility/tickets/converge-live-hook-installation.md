# Converge live hook installation

Blocked by: expose-hook-health-record.md
Ownership fence: `internal/adopt/link_hook.go`, `internal/adopt/link_stage.go`, `internal/adopt/link_transaction.go`, `internal/contract/surface/link_test.go`, `internal/contract/surface/link_lifecycle_test.go`
Integration surfaces: hook-health classifier and currency→expose-hook-health-record.md; live branch substitution→`internal/adopt/link_stage.go` and `internal/adopt/link_transaction.go` + S1/S2; hook-path safety→`internal/adopt/link_hook.go` and `internal/adopt/link_transaction.go` + S3/S4; lifecycle contracts→`internal/contract/surface/link_test.go` and `internal/contract/surface/link_lifecycle_test.go` + S1/S2/S3/S4
Contracts: hook classification and expected bytes cross `internal/adopt/link_hook.go`→`internal/adopt/link_transaction.go`, asserted by S3 against the real hook-health producer
Closure: S1/single-substitution, S2/origin-head-probe, S2/probe-failure-open, S3/dangling-foreign, S3/dangling-refusal, S4/markerless-refusal

## What to build

Remove the dead installer only after moving its live `origin/HEAD` probe and branch substitution into the active transaction, then classify dangling paths as foreign before link can overwrite them. Keep that live-installation move together: deleting the dead installer alone strands the link lifecycle contract red when an unset `origin/HEAD` is never populated; moving path safety without the transaction refusal strands the dangling-symlink refusal red.

## Acceptance

- [ ] [S1] Exactly one live site substitutes the branch token into the pre-push template.
- [ ] [S2] Link populates an unset `origin/HEAD` from a remote when possible and probe failure does not fail installation.
- [ ] [S3] A dangling symlink at the hook path is foreign and link refuses it without overwrite.
- [ ] [S4] Link continues to refuse a marker-less foreign pre-push hook.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| S1/single-substitution | retain `installGitHook`’s duplicate replacement | focused source call-site test | enumerate `ReplaceAll` calls over the template/token pair, expect exactly one live site |
| S2/origin-head-probe | delete the dead installer without moving its probe | lifecycle surface contract | link a repo with remote HEAD and no local origin/HEAD, expect the ref to appear |
| S2/probe-failure-open | propagate `remote set-head` failure | lifecycle surface contract | link with an unreachable remote, expect success and an installed hook |
| S3/dangling-foreign | read a dangling hook before `Lstat` | hook-health and lifecycle contracts | plant a dangling symlink, inspect health, expect foreign |
| S3/dangling-refusal | treat the dangling symlink as absent | lifecycle surface contract | run link over the dangling hook, expect refusal and unchanged symlink |
| S4/markerless-refusal | weaken the foreign-hook marker check | existing lifecycle surface contract | plant a marker-less executable, run link, expect refusal and unchanged bytes |
