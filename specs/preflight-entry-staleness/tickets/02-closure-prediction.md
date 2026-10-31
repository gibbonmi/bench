# Predict the closed row paths of a deliverable from the landing derivation

Blocked by: none
Writes: internal/commitment/delivery.go, internal/commitment/delivery_test.go, internal/commitment/authority.go, internal/commitment/repository/closure.go, internal/commitment/repository/repository.go, internal/commitment/repository/candidate.go, internal/commitment/repository/closure_test.go, internal/commitment/commitmenttest/repo.go
Covers: ES13, ES14

## What to build

A caller can ask which roadmap row files the delivery of an approved deliverable deletes, before any completion evidence exists, and get the answer that the landing uses.
Read the spec section "One closure prediction" and the coverage rows ES13 and ES14 before you start.

Make these changes:

- Add one pure function to `internal/commitment` that returns the outcome sources that the delivery of an approved deliverable closes. Move into it the `activeDeliverable` lookup, the obligation-free refusal, and the obligation filter of `Deliver`. Do not move the already-delivered refusal, because it reads the delivery fact.
- Make `Deliver` call the new function for its closed set. Keep every delivery-fact step of `Deliver`, the already-delivered refusal included. `Deliver` keeps its signature, so its callers in `delivery_test.go`, `continuation_test.go`, and `internal/spec/retire_test.go` need no edit.
- `obligationFree` keeps its name, because `internal/commitment/authority.go:128` calls it. If a move of that helper is needed, update that caller.
- Add one store method to `internal/commitment/repository` that reads the policy at a tree the way `policyAt` does, calls the pure function, and returns each closed source path that `roadmap.RowOwner` admits. A policy that the tree does not hold, or a deliverable that the active milestone does not approve, gives no path. A refusal or a policy that does not parse returns its error.
- Make `Store.Closure` use the same row-owner step for its deletion edits, so the store method and `Store.Closure` share one derivation. `Store.Closure` keeps its signature and its edits, so its callers in `candidate.go:226`, `internal/gate/completion.go:138`, and `internal/landing/published/published.go:62` need no edit.
- `protectedCandidate` holds a third `RowOwner` filter (`internal/commitment/repository/candidate.go:277`). Leave it unless the shared step is the same rule; if you fold it, keep its answer.
- Add `TestCommitmentCloses` to `delivery_test.go` and `TestClosedRowPathsMatchClosure` to `closure_test.go`. The ES14 fixture binds an obligation source whose path `roadmap.RowOwner` refuses beside a row source. If `commitmenttest` needs a helper to seed that source, add it to `repo.go`.

## Acceptance

- [ ] The pure function returns the obligation sources that `Deliver` closes for the same policy and deliverable (ES13).
- [ ] For a binding with a row source and a source that `roadmap.RowOwner` refuses, the store method returns only the row path, and `Store.Closure` deletes exactly that path (ES14).
- [ ] The existing `Deliver` and `Store.Closure` tests pass unchanged, the obligation-free refusal and the already-delivered refusal included.
