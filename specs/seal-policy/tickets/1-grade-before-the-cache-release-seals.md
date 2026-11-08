# Grade every user before the cache release seals

Blocked by: none
Writes: internal/subprocess/protection/recovery.go, internal/subprocess/protection/recovery_test.go, internal/subprocess/protection/fixture_test.go, internal/subprocess/protection/store.go, internal/subprocess/protection/directory.go, internal/subprocess/protection/lifecycle_test.go, internal/gocache/lock.go, internal/gocache/protection_test.go, internal/systemtest/process_lifetime_callers_test.go, internal/systemtest/process_lifetime_release_callers_test.go, internal/systemtest/process_lifetime_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/git_fixture_owner_test.go, docs/adr/0030-one-lifetime-owner-and-process-user-protection.md, CHANGELOG.md
Covers: SP1, SP2, SP3, SP4, SP5, SP6, SP7, SP8, SP9, SP10, SP34

## What to build

Review chunk: SP-C1.

Add a release observation to the protection protocol beside `SealAndObserve` (`internal/subprocess/protection/recovery.go:17`). Under the store's protocol lock, it reads the generation strictly and grades every user with `grade`. A read failure returns uncertain and writes nothing. When any user is protected or uncertain, it returns that disposition with each unresolved user and writes no record. When every user is safe, it runs the existing tail: seal, remove each user, write the retired record, and remove the generation. The two observations share one body whose only difference is the seal order. Keep one grade loop and one tail. The release observation never writes admission open.

Call the release observation from `heldLock.retire` (`internal/gocache/lock.go:172`). Keep `SealAndObserve` for every other caller, and keep its signature. Keep the `discardRetiredCaches` call in `retire` and the test name `TestReleaseRunsTheRetiredCheck`, because `internal/conformance/injected_ports_registry_test.go:188-189` pins both.

Change the comments at `internal/gocache/lock.go:157-160`, `internal/subprocess/protection/store.go:93-95`, and `internal/gocache/protection_test.go:193-194` to the new state. Drop the private HOME of the descendant row at `internal/systemtest/process_lifetime_callers_test.go:222-226`. Add the SP10 system test there, or in `process_lifetime_release_callers_test.go` when the callers file has no headroom (it has 53 lines). The new system test obeys the serial and parallel census rules of `internal/systemtest`.

Add one sentence to the protocol paragraph of ADR 0030: the last release of a cache hold grades every user first, and it seals and retires the generation only when every user is safe. Add one `CHANGELOG.md` entry under `## [Unreleased]`.

## Acceptance

- [ ] After a release beside a live pending user, a new hold of the same cache joins the same generation without an error (SP1).
- [ ] After a release observation beside a protected user, a `Prepare` on the same descriptor succeeds (SP2).
- [ ] The release observation beside a protected user returns the protected state with that user's id (SP3).
- [ ] A release observation whose every user is safe leaves the retired record and no generation (SP4).
- [ ] After a release observation beside a user whose probe answers EPERM, a `Prepare` succeeds (SP5).
- [ ] After a seal-first observation and then a release observation beside a live user, a `Prepare` still returns `ErrSealed` (SP6).
- [ ] A failed seal publish on the safe release path returns uncertain and keeps every user record (SP7).
- [ ] `TestRegistrationSealedAdmission` and `TestProcessLifetimeCacheSecondProcessClean` stay green (SP8, SP9).
- [ ] A second `bench test` on the HOME of a run that leaked a live descendant runs its go test and prints no `the resource generation is sealed` (SP10).
- [ ] ADR 0030 states the release policy (SP34).

## Verification

Run `bench test --package ./internal/subprocess/protection` and `bench test --package ./internal/gocache` whole, not through `--run`. This ticket writes system-tagged test files, so `BENCH_KIT` must name the candidate kit for each system run; `bench test --check system` supplies it. Run the system suite through `bench test --check system --run TestFocusedTestAfterALeakedDescendant`, then `bench test --check system` once. Observe each planned row red before its production edit, and record the route of each red. Self-probe: move the seal publish back before the grade in the release observation, and report that SP1 and SP2 go red.
