# Name the repair root in the landing refusal

Blocked by: none
Writes: internal/freshness/freshness_verify.go, internal/freshness/freshness_verify_test.go, internal/runbinary/runbinary.go, internal/runbinary/runbinary_test.go, internal/systemtest/owner_test.go, internal/systemtest/owner_stale_seal_test.go (new), internal/gate/prospective.go, internal/gate/prospective_owner_test.go
Covers: BF5, BF6, BF7, BF8, BF9

## What to build

Verify the premise first. `refusal` in internal/freshness/freshness_verify.go
formats `RebuildAction` from the caller's digest root, and the prospective
owner in internal/gate/prospective.go passes the composed tree as that root.
`Factory.validate` in internal/runbinary/runbinary.go verifies an inherited
executable with the seal pair alone.

Then give `refusal` a repair root distinct from the digest root, threaded
through a new `Factory` field. Make the prospective owner pass the kit checkout
as the repair root. The system suite refusal fires inside `TestMain`.

Its test therefore runs the suite binary as a child with a stale seal and reads
the exit code. Make
`Factory.validate` call `freshness.Verify` against a named source root for an
inherited executable. When the caller names no root, keep the seal-pair check. Make the system suite owner in
internal/systemtest/owner_test.go verify `BENCH_RUN_BINARY` against
`BENCH_KIT` at setup. Leave `runbinary.Own` unchanged.

Run the system-tagged tests with `BENCH_KIT` and `BENCH_RUN_BINARY` set, as
the system suite requires.

## Acceptance

- [ ] The prospective owner refusal prints the kit root inside its rebuild command.
- [ ] An inherited executable with a stale source digest and a named source root refuses in `validate`.
- [ ] The system suite owner refuses a binary whose seal mismatches `BENCH_KIT`.
- [ ] `TestVerifyUsesContentRatherThanMtime` and `TestFactoryOwnBuildsOnePrivateAbsoluteSelectionAndCleansIt` still pass.
- [ ] Self-probe: pass the digest root as the repair root again, and report the prospective owner test red.
