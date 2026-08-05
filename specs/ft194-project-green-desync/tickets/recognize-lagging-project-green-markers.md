# Recognize lagging project-green markers

Blocked by: none
Ownership fence: `internal/gate/authorization`, `internal/specbuild`, `cmd/bench/specbuild.go`, `cmd/bench/wiring_pins.go`
Contracts: the marker-advance request `(context, root, branch, destination full commit ID, expected lineage full commit ID or empty)` crosses `internal/gate/authorization`→`internal/specbuild`→`cmd/bench/specbuild.go`, asserted by MA1-MA6 against the real owner and production compile pin; actual-marker reads precede ancestry checks and the compare-and-swap, and absence preserves the existing empty-expectation behavior

## What to build

The gate-authorization owner exposes one focused marker-advance operation that recognizes an actual marker lagging behind the caller's expected lineage, validates both ancestry relationships, and compare-and-swaps from the actual marker position. `Bootstrap` composes it after its existing branch-stability and composed-green checks. The lifecycle port, production adapter, and test adapters expose the same operation without re-deriving marker ancestry.

## Acceptance

- [ ] [MA1] When the actual marker differs from the expected lineage but is an ancestor of both that lineage and the destination, the owner advances it to the destination by compare-and-swap from its actual position.
- [ ] [MA2] A marker divergent from the expected lineage refuses and remains untouched.
- [ ] [MA3] A marker strictly between the expected lineage and the destination refuses and remains untouched.
- [ ] [MA4] A marker whose ancestry cannot be decided refuses and remains untouched.
- [ ] [MA5] A marker already at the destination stays an idempotent success, while absent-marker and empty-expectation behavior remain unchanged.
- [ ] [MA6] `Bootstrap`, the lifecycle `GateOwner` port, and the production adapter all route marker advancement through the one authorization-owned operation; lifecycle code contains no marker-ancestry derivation.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MA1 | restore the equality-only conflict check before advancing | authorization owner unit test | plant the marker at the grandparent, pass its child as the expected lineage and the green tip as destination, invoke the owner, run `go test -count=1 ./internal/gate/authorization -run TestAdvance`, expect the conflict failure |
| MA2 | accept every actual marker that is an ancestor of the destination | authorization owner unit test | plant a sibling marker, invoke the owner, run the focused authorization test, expect refusal and an unchanged marker |
| MA3 | omit the actual-marker-to-lineage ancestry check | authorization owner unit test | plant the marker between lineage and destination, invoke the owner, run the focused authorization test, expect refusal and an unchanged marker |
| MA4 | treat an ancestry-probe error as recognition | authorization owner unit test | remove the planted marker's object after it peels, invoke the owner, run the focused authorization test, expect refusal and the marker file unchanged |
| MA5 | remove the destination-equality early accept or change empty expectation to a nonzero compare value | existing idempotent and absent-marker authorization controls | run the focused Bootstrap controls, expect replay or fresh bootstrap to fail |
| MA6 | bypass the focused owner operation in `Bootstrap` or the production adapter | authorization composition tests and production compile pin | run `go test -count=1 ./internal/gate/authorization ./internal/specbuild ./cmd/bench`, expect the lagging-marker composition or interface assertion to fail |
