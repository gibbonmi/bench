# Record the resource path of a retired store

Blocked by: none
Writes: internal/subprocess/protection/record.go, internal/subprocess/protection/recovery.go, internal/subprocess/protection/export_test.go, internal/subprocess/protection/recovery_test.go, internal/subprocess/protection/lifecycle_test.go, internal/gocache/protection_test.go, internal/worktree/protection_test.go, .bench/commitment.json
Covers: RS1, RS2, RS3, RS25, RS26

## What to build

A safe retirement leaves one retired record in the store, and `resource.json` stays absent as today.
The record owner adds `retired.json` to its store file set with one strict record type: schema, kind, resource path, and binding.
The retired validator follows the `resourceRecord.valid` pattern and reuses the descriptor's owner checks for the kind, the path, and the binding (`Descriptor.checkOwner` in the record owner). It does not paste a second path rule.
`SealAndObserve` publishes the retired record through the store's durable publish after it removes each user and before it removes `resource.json`.
Every caller of `SealAndObserve` gets this behavior: the cache retire, `ObserveRecorded`, `Generation.Release`, and `Recover`.

`Recover`, `NoGeneration`, the store join, and `ObserveRecorded` do not change. A retired store still reads as a store with no generation, and the next open still publishes a fresh generation.
The export test file exports the retired file name for the tests.

## Acceptance

- [ ] After a safe retirement, the store holds `retired.json` with the kind, the resource path, and the binding of the retired generation, and it holds no `resource.json`.
- [ ] A test first requires `retired.json` in a retired store, then requires that `Recover` of that store gives a result that `NoGeneration` accepts.
- [ ] `TestProtectionCompletedHold` opens a fresh generation after the release, and a clean after the release exits zero.
- [ ] `TestCheckoutProtectionPredicate` first requires `retired.json` in the retired checkout store, then requires that the checkout cleanup finds no user.
