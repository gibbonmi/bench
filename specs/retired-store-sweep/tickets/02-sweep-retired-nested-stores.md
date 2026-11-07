# Sweep the retired stores of the resources that a removal deleted

Blocked by: 01-record-retired-store-path.md
Writes: internal/subprocess/protection/directory.go, internal/subprocess/protection/store.go, internal/subprocess/protection/record_test.go, internal/subprocess/protection/lifecycle_test.go, internal/subprocess/protection/fixture_test.go, docs/adr/0030-one-lifetime-owner-and-process-user-protection.md, DATA_HANDLING.md
Covers: RS4, RS5, RS6, RS7, RS8, RS14, RS15, RS31

## What to build

`discardNested` deletes a retired store when the one retired-store predicate passes and the recorded path is inside the removed directory.
Its live-generation branch does not change.

The predicate lives in the directory owner, so ticket 03's absent-resource check calls the same predicate.
Before it locks, it skips a store with a `resource.json` entry or with no `retired.json` entry.
Under a non-blocking store lock, it requires all of these:

- `resource.json` is absent;
- the retired record reads through the record owner's strict reader from ticket 01;
- the users directory holds no entry;
- the recorded resource path is absent.

The store owner's lock helper gains a non-blocking form beside its blocking form, and both share one lock routine.
A contended lock, an unreadable retired record, or a failed lock skips the store with no error.

The existing `TestRecord` tests in `record_test.go` run over each record file name through their record planting helpers (`plant` and `resourceRecord`). They do not paste a second harness.

The protection ADR consequence states that a retired store keeps a record of its resource and that a protected removal deletes it with its resource. The sentence that keeps a retired store until the next boot goes. The ADR holds no path and no code.
The process resource protection section of `DATA_HANDLING.md` names the retired record and states that it stays until its resource is gone. Leave its sentence about a store beside its resource as it is.

## Acceptance

- [ ] `TestGenerationRemove` deletes the retired store of a nested resource that the removal deleted, built through a real `SealAndObserve` retirement.
- [ ] A retired nested store whose resource still exists stays, and a retired store of a sibling path that shares the prefix stays.
- [ ] A retired store whose users directory holds an entry stays after a sweep.
- [ ] A retired store whose lock another open file holds stays after a sweep, and the sweep returns no error for it.
- [ ] The record tests refuse a retired record that is malformed, a link, at another mode, has an unknown field, or exceeds the read limit, and a sweep keeps each such store.
- [ ] The ADR consequence and the data-handling section state the decided retired-store lifetime.
