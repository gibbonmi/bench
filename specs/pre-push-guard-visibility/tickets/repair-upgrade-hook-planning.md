# Repair upgrade hook planning

Blocked by: none
Ownership fence: `internal/adopt/link_hook.go`, `internal/adopt/upgrade.go`, `internal/contract/surface/link_test.go`, `internal/contract/surface/upgrade_test.go`
Integration surfaces: hook eligibility and prospective bytes→`internal/adopt/link_hook.go` + RP1/RP2/RP3; upgrade plan consumer→`internal/adopt/upgrade.go` + RP1/RP2/RP3; single-substitution enumeration→`internal/contract/surface/link_test.go` + RP4; upgrade runtime contract→`internal/contract/surface/upgrade_test.go` + RP1/RP2/RP3
Contracts: hook eligibility, including absent, current, stale, foreign, dangling, and special-file states, crosses `internal/adopt/link_hook.go`→`internal/adopt/upgrade.go`, asserted by RP1/RP2/RP3 against the real producer; substitution-site enumeration crosses `internal/adopt/link_hook.go`→`internal/contract/surface/link_test.go`, asserted by RP4 against the real source files
Closure: RP1/shared-plan-eligibility, RP2/fifo-no-block, RP3/dangling-no-promise, RP4/enumerated-single-substitution

## What to build

Make upgrade planning consume the authoritative hook-health and prospective-write eligibility seam, so its count describes only a transaction that can run. Keep extraction, plan use, and the three contract cases together: an extraction-only cut strands the upgrade runtime contract red because it must prove the stale/current count, FIFO return, and dangling no-promise against the real shared producer.

## Acceptance

- [ ] [RP1] An unequal-version `bench upgrade --check` counts only a shared-seam-eligible stale or absent hook refresh and omits a current or refused hook.
- [ ] [RP2] A writerless FIFO at the effective pre-push path returns within the contract deadline and is never counted as a refresh.
- [ ] [RP3] A dangling pre-push symlink is not counted by `--check`; plan mode and apply retain the same refusal posture.
- [ ] [RP4] The single-substitution contract enumerates every production consumer of the template and reports exactly one branch-token replacement site.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RP1/shared-plan-eligibility | restore an upgrade-local marker/read/template derivation | upgrade surface contract | run unequal-version stale/current/refused fixtures, expect the stale/current or refused count mismatch |
| RP2/fifo-no-block | read the hook before the producer’s special-file classification | deadline-backed upgrade surface contract | place a writerless FIFO at pre-push, run `upgrade --check`, expect the bounded refusal |
| RP3/dangling-no-promise | treat any read error as a prospective refresh | upgrade surface contract | replace pre-push with a dangling symlink, run `--check` then apply, expect no counted write and matching refusal |
| RP4/enumerated-single-substitution | add a second `strings.ReplaceAll` template substitution in upgrade planning | single-substitution surface contract | run the source enumeration, expect its replacement-site count red |
