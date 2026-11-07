# Discard gone retired caches at the last cache release

Blocked by: 01-record-retired-store-path.md, 02-sweep-retired-nested-stores.md
Writes: internal/subprocess/protection/directory.go, internal/subprocess/protection/directory_test.go, internal/subprocess/protection/export_test.go, internal/gocache/lock.go, internal/gocache/protection.go, internal/gocache/protection_test.go, internal/conformance/injected_ports_registry_test.go, docs/adr/0030-one-lifetime-owner-and-process-user-protection.md, internal/subprocess/protection/fixture_test.go
Covers: RS9, RS10, RS11, RS12, RS13, RS27, RS28, RS30

## What to build

A new kind-scoped absent-resource check on `Directory` lists the store root and deletes each store of its kind that passes the retired-store predicate from ticket 02.
It keeps both kind tests: the store-name prefix filter and the record kind comparison.
A store that it cannot lock or read stays and adds no error. A root that it cannot list, or a deletion that fails, joins its error.

The cache retire calls the check through one package variable, `discardRetiredCaches`, declared beside its reader in `lock.go`.
It calls the check only after the hold's own `SealAndObserve` returns without an error. A contended release returns before the check, as it does today.
The injected-port registry admits the variable with one `seamVariable` row, placed between the `internal/git` and `internal/guards` rows: package `internal/gocache`, variable `discardRetiredCaches`, owner function `heldLock.retire`, purpose "a recording check shows that only the last uncontended release runs the absent-resource check, after its own retirement", `serialIsolation`, and consumer `internal/gocache/protection_test.go::TestReleaseRunsTheRetiredCheck`.
That consumer replaces the variable with a recorder and does not call `t.Parallel`.

The protection tests drive the check on a private store base through `SetStoreBase`.
The failed-deletion case uses a 0700 store that passes the predicate and holds a non-empty 0500 subdirectory. It skips through the privilege class of `internal/capability` when that subdirectory stays writable, as the gate's write-denied guard does.

The protection ADR consequence states that the last uncontended release of a cache hold deletes each retired cache store whose cache is gone. The ADR holds no path and no code.

## Acceptance

- [ ] The check deletes a retired cache store whose cache directory is gone and keeps one whose cache directory exists.
- [ ] The cache-kind check keeps a retired store of another kind whose resource is gone.
- [ ] An unreadable retired record keeps its store, and the check returns no error for it.
- [ ] A deletion that fails joins the error that the check returns.
- [ ] The last uncontended release records one call of the check after its own retirement, and a contended release records none.
- [ ] The injected-port registry check accepts the new variable only with its row.
- [ ] The ADR consequence names the cache release as the second deleter.
