# Review pickup: ledger-settle-policy

Range: `897f52be..ddfc8a7f` on the integration source `ft302-c34-spec`.
One round, three axes on opus at medium. The accepted repairs run as tickets
09 and 10. The findings below stay open for the reviewer.

## Standards

Count: 8 raw, 8 repair targets. Worst: a doc comment false about the constants it heads (repaired in ticket 10).

- `internal/intent/transaction_test.go:71` — the comment says the expected values were recorded against the pre-migration tree. `craft-comments` keeps a red record in the commit or the spec. The spec's build notes now carry the record. Disposition: ask-user (cut the sentence or keep it as the independence rationale).
- `internal/intent/admissionpolicy/records.go:75` — `CompareAndSwapRequestDigest` is exported with one caller in its own file. Lazy Element. Disposition: ask-user (unexport).
- `internal/intent/ledger_aliases_test.go:33` — the alias set is enumerated by hand, a second copy of the leaf's exported surface. A new leaf name with no alias passes. Disposition: ask-user (accept, or derive the set with `go/types`).
- `internal/landing/composition.go:156` — `unionStages` hardcodes stages 2 and 3, which the settle policy owns. Disposition: ask-user (accept, or export the stage numbers).
- `internal/landing/composition.go:157` — after ticket 10 the policy answers a removal for a side-less union, so the `!hasDestination && !hasSource` arm in `unionStages` is unreachable. Disposition: ask-user (delete the arm, or keep it as a stated guard). Found by the repair-scoped re-review.

## Spec

Count: 4 raw, 2 repair targets. Worst: the two union journeys deleted with no policy twin (repaired in ticket 10).

- `spec.md`, row LS41 — the row promises that the adapter renders no `union content not text` for a side-less union, but it cites only the policy table case. Disposition: ask-user (narrow the row to the policy, or cite an adapter journey). Found by the repair-scoped re-review.
- `spec.md`, the message shape — "the refusing paths, comma-separated in merge-tree order". Every producer supplies one path, and `Settle` returns on the first refusal, so the list never joins. Disposition: ask-user (narrow the promise to one path, or make a producer accumulate).

## Coverage

Count: 5 raw, 3 repair targets. Worst: a side-less union path renders `union content not text` (repaired in ticket 10).

- The tolerant mode has no schema guard: a purge over a ledger at schema 99 rewrites it at schema 2 and drops an unknown field. Pre-existing; the tolerant read is the new home for a guard. Disposition: ask-user.
- `ConflictKind(nil)` answers `textual`, and a conflict output with no record answers a refusal that names no path. Pre-existing and not reached from real `merge-tree` output. Disposition: no-op.
