# Test the authorization owner's ancestor refusal

Blocked by: nothing
Ownership fence: `internal/gate/authorization/authorization_test.go`
Assumptions: `Bootstrap` already refuses an expected marker that is not an
ancestor of the tip, already returns nil when the marker equals the tip, and
already fails closed when the marker names an object the store cannot peel to a
commit. This ticket adds no production behavior; it exists because that refusal
carries no test today and is about to become load-bearing. Confirm the refusal
is still unasserted at pickup — if a sibling ticket has already covered it,
stop and say so rather than writing a second copy.

## What to build

`Pass the found green marker as the expected prior tip` makes the lifecycle
start passing a real expectation into `Bootstrap` instead of an empty one. That
turns `Bootstrap`'s ancestor rule from decoration into the thing standing
between a benign stale marker and a divergent one that must never be
fast-forwarded over.

The rule cannot be observed from the lifecycle package: its in-package gate
double is a bare compare-and-swap with no ancestry rule at all, so a test
written there would assert the fixture. This is a deliberate seam split, and
this ticket is the half that runs against the real owner.

Test the refusal directly, at the owner, against a real temporary repository.
The cheap wrong fix in the sibling ticket is to pass the found marker
unconditionally and let `update-ref` sort it out; these criteria are what make
that fix red.

## Acceptance

- [ ] AU1 — `Bootstrap` refuses when the expected marker is a commit on a diverged branch that is not an ancestor of the tip, and the marker ref is left unmoved.
- [ ] AU2 — `Bootstrap` advances the marker when the expected marker is a strict ancestor of the tip, and the marker ends at the tip.
- [ ] AU3 — `Bootstrap` returns nil and writes nothing when the existing marker already equals the tip, whatever the expectation.
- [ ] AU4 — `Bootstrap` refuses when the existing marker names an object absent from the object store, rather than treating it as absent and creating a fresh marker.
- [ ] AU5 — `Bootstrap` accepts a full 40-hex expectation equal to the existing marker, and refuses that same marker expressed as its ref name or as an abbreviated object ID.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AU1 | delete the `merge-base --is-ancestor` guard and rely on `update-ref`'s own compare-and-swap | `TestBootstrapRefusesDivergentExpectedMarker` | build a repo with green evidence at the tip, commit a divergent branch, plant the marker there, call `Bootstrap` with it as expected, assert the error and that the marker still names the divergent commit |
| AU2 | invert the ancestor comparison's operand order | `TestBootstrapFastForwardsAncestorMarker` | plant the marker one commit behind the tip, call `Bootstrap` with it as expected, assert no error and the marker at the tip |
| AU3 | move the marker-equals-tip short circuit below the ancestor check | `TestBootstrapIsNoOpWhenMarkerEqualsTip` | plant the marker at the tip, call `Bootstrap` with an unrelated expectation, assert no error and no ref write |
| AU4 | treat an unpeelable marker as absent | `TestBootstrapFailsClosedOnUnreadableMarker` | plant the marker at a commit, delete that object from the store, call `Bootstrap`, assert an error rather than a fresh marker |
| AU5 | resolve the expectation through `rev-parse` before comparing it, so a ref name and an abbreviation both match | `TestBootstrapExpectationIsAFullObjectID` | plant an ancestor marker, call `Bootstrap` three times with the full OID, the ref name, and a 7-character abbreviation, assert success then refusal then refusal |
