# Expand the verdict record with the reduced shape

Blocked by: Declare the reduced gate scope

Ownership fence: `internal/gate/verdict.go`, `internal/gate/reduced_verdict_test.go`
Assumptions: the verdict cache stays a single slot with strict field validation, and no execution path emits a reduced record yet — this ticket only teaches the record and the loader its shape

## What to build

A reduced verdict is a distinct record class, not the existing record with extra
fields set. Teach the record and its loader that second shape: the reduced marker,
the phases that actually ran, and the ancestor's identity and recorded time. The
existing full shape round-trips unchanged.

Because the slot is single, the reduced record carries the ancestor's identity
forward rather than pointing at a record that may no longer exist. Inherited
evidence does not re-stamp its recorded time — the ancestor's own time is what the
freshness window applies to, so a stale ancestor can still be recognized as stale.

The loader already validates fields strictly and must keep failing closed: a record
carrying an ancestor but no reduced marker, or a reduced marker with no ancestor, is
a malformed hybrid and is rejected rather than guessed at. Expand only here; the
execution path that produces reduced records is a later ticket.

## Acceptance

- [ ] [R15] A reduced record round-trips its marker, its executed phase list, and its ancestor identity and recorded time, and the existing full-shape record round-trips unchanged.
- [ ] [R16] The loader rejects a record mixing the full and reduced shapes in either direction.
