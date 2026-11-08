# Keep the request token in the binding and name it in the landing refusals

Blocked by: none
Writes: internal/intent/ledger/commitment.go, internal/intent/ledger/commitment_test.go (new), internal/intent/assignment.go, internal/commitment/repository/admission.go, internal/commitment/repository/admission_token_test.go (new), internal/commitment/repository/readiness.go, internal/commitment/repository/candidate.go, internal/worktree/worktree.go, internal/worktree/land_identity.go, internal/worktree/land_reauthorization_test.go, internal/worktree/identity_component_test.go, internal/worktree/land_window_test.go, internal/worktree/refusal_route_follow_test.go, internal/worktree/parallel_census_test.go, DATA_HANDLING.md, .bench/commitment.json, CHANGELOG.md
Covers: RP45, RP46, RP47, RP48, RP49, RP50, RP56

## What to build

Review chunk: RP-C6.

Add `request_token` to `DeliveryBinding` (`internal/intent/ledger/commitment.go:46-53`), with `omitempty`. `Store.Start` (`internal/commitment/repository/admission.go:44`) writes it from its `--request` value. `ValidateCommitment` refuses a present token that is not line-safe, or whose `RequestDigest` differs from the binding's request digest. A binding without the field stays valid. `RegisterSibling` (`readiness.go:138`) copies a binding to a sibling with another request digest, so the copy clears the token. The fixture `startOtherBuild` (`land_window_test.go:45-59`) clones a binding and rewrites its request, so it clears the token too. Add a `commitment.bindings.request_token` row to `DATA_HANDLING.md`. The row states that the token is the opaque request label that `bench worktree list` already prints.

Add one function to `internal/intent/assignment.go` that answers the bound request of an assignment from the ledger, with a request filter. It is a function, because `DeliveryBinding` is a type alias in package `intent`. With a token, the sentence is `; the delivery binding of assignment <id> holds the request token <token>`. Without one, it is `; the delivery binding of assignment <id> holds the request digest <digest>, and its token is unknown`.

When `readyFor` (`readiness.go:40-84`) finds a binding of the owner assignment with another request digest, its unbound error adds the sentence and still wraps `errUnbound`. On the light path, `unboundNaming` (`candidate.go:129-147`) puts the production-path clause before the sentence. When the unknown-request refusal of a landing (`internal/worktree/worktree.go:519-566`, through `land_identity.go:148`) names one candidate assignment that holds a binding, the landing caller adds the sentence to the detail. The shared identity component and the release refusal keep their text. `internal/commitment/repository/publication_test.go:68,108` keeps its expectation, because it matches the unbound prefix with `strings.Contains` (`publication_test.go:85,144`). Move the exact expectations at `land_reauthorization_test.go:33` and `identity_component_test.go:44`. A new top-level test in `internal/worktree` moves the test-count pin in `parallel_census_test.go`.

This ticket declares the persisted field of the spec, so run each record verb and checkpoint of the build on an installed candidate, after `bench worktree build --install`. Add one `CHANGELOG.md` entry under `## [Unreleased]` for RP-C6. This ticket owns `.bench/commitment.json` for the spec-stage binding commit.

## Acceptance

- [ ] After `bench commitment start --request <token>`, the binding holds `request_token` equal to the token (RP45).
- [ ] A landing after a reauthorize to a new token refuses with `holds the request token <original token>` (RP46).
- [ ] That refusal still satisfies `errors.Is` with the unbound error (RP47).
- [ ] An unknown-request landing whose one candidate holds a binding names the bound token in its detail (RP48).
- [ ] A binding with no token refuses with `holds the request digest <digest>, and its token is unknown` (RP49).
- [ ] `ValidateCommitment` refuses a token whose digest differs from the request (RP50).
- [ ] `TestCommitmentSiblingBinding` passes, and the sibling binding carries no token (RP56).

## Verification

Run `bench test --package ./internal/intent/...`, `bench test --package ./internal/commitment/...`, and `bench test --package ./internal/worktree` whole. Run `bench test --package ./internal/conformance --run '^TestCommitmentRecord'` for the `DATA_HANDLING.md` row, and `bench test --check system --run Adoption`. Observe each planned row red before its production edit, and record the route of each red. Self-probe: keep the source token in the sibling copy, and report that RP56 goes red.
