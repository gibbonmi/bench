# Pass the found green marker as the expected prior tip

Blocked by: Reach recomposition before the gates it discards
Ownership fence: `internal/specbuild/assign.go` (`startRun` only),
`internal/specbuild/lifecycle.go` (`finishStart` only),
`internal/specbuild/start_test.go`
Assumptions: `Start`'s fresh path calls `startRun` with an empty
`previousGreen`; its restart-after-terminal path in the same function already
passes the prior base, so the asymmetry is visible in one function.
`finishStart` passes an empty expectation at its own bootstrap call. The
authorization owner already refuses a non-ancestor expectation, already
short-circuits when the marker equals the tip, and already fails closed on a
marker naming a missing object — no new ancestor check belongs here. Re-derive
both sites from the tree at pickup; two earlier tickets have edited these files.

## What to build

Once any spec build has completed on a branch, `start` refuses forever. It asks
the gate owner for green evidence with **no** expected prior marker, so the
`refs/bench/green/<branch>` marker its own predecessor left behind reads as
`project-green marker conflicts with another tip` — a conflict, rather than the
benign ancestor it is. Nothing in the kit ever moves or retires that marker, so
the refusal is permanent. `main` in this repository is in exactly that state
today.

Stop passing an empty expectation. Read the current
`refs/bench/green/<branch>` and hand it to the owner as the expected prior tip,
at both sites: the fresh-start bootstrap in `startRun`, which is where the
observed refusal comes from, and the completion helper in `finishStart`, which
carries the same empty expectation and must stop disagreeing with it.

**The empty string is a meaningful value, not a missing one.** The owner treats
`""` as "no marker expected" and substitutes the zero object ID for its
compare-and-swap; it treats a non-empty expectation as a commit that must be an
ancestor of the tip. So a marker that cannot be read must be passed as `""` and
a marker that can be read must be passed verbatim — passing a zero object ID
literal for an absent marker would send it into the ancestor check, which
reports false for it, and turn "no marker yet" into a refusal. Read the marker
with the same `^{commit}` peel the owner uses, so the two sides cannot disagree
about what "unreadable" means.

**Every criterion below runs against the in-package gate double, and the double
is more permissive than the real owner.** `greenGate.Bootstrap` is a bare
`git update-ref <marker> <tip> <expected>`. `update-ref` resolves ref names and
abbreviated object IDs, so it accepts expectation values the real
`authorization.Bootstrap` rejects — the real owner compares `existing !=
expected` as **strings**, against a full 40-hex OID it obtained with a
`^{commit}` peel. Hand it the ref name `refs/bench/green/<branch>` and the
double is green while the real owner refuses with
`project-green marker conflicts with another tip`.

That is the shape that cost the per-component-gate-scoping build a delegate
round: two correct, disjoint fences, each green alone, red composed, because a
fixture accepted what the real counterpart would not. The end-to-end row is
three tickets downstream, so ST6 pins the handover **form** here, at the
junction, and `Test the authorization owner's ancestor refusal` carries the dual
on the real owner.

**Do not teach the double the ancestry rule.** The divergent-marker refusal is
deliberately tested one layer down, against the real owner. Upgrading the double
would make that case assert the fixture instead of the product. This is a
structural constraint on the build, not a behavior — the spec classifies it
`not TDD-able` and it stays that way. Do **not** write a test pinning the
double's behavior; a test that asserts the fixture is the very thing this
constraint exists to prevent. State in your checkpoint that you left the double
alone, and the coordinator will verify it by reading the definition.

## Acceptance

- [ ] ST1 — a fresh `Start` succeeds when `refs/bench/green/<branch>` is a strict ancestor of the tip, and the marker ends at the tip.
- [ ] ST2 — a fresh `Start` still succeeds when no marker exists at all, and the marker is created at the tip.
- [ ] ST3 — a fresh `Start` still succeeds when the marker already equals the tip, and the marker is unchanged.
- [ ] ST4 — a `Start` whose marker names a commit absent from the object store refuses rather than fast-forwarding or creating a fresh marker.
- [ ] ST5 — the completion helper passes the same expectation as the fresh-start path for the same marker state.
- [ ] ST6 — the expectation handed to the gate owner is a full 40-hex object ID obtained with a `^{commit}` peel, or the empty string when the marker is unreadable, and never a ref name or an abbreviated ID.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| ST1 | patch the completion helper only, leaving the fresh-start bootstrap passing `""` | `TestStartFastForwardsAncestorMarker` | plant the marker one commit behind the tip, call `Start` on a fresh run, assert no error and the marker at the tip |
| ST2 | pass a zero object ID literal instead of `""` when the marker is unreadable | `TestStartSucceedsWithNoMarker` | delete the marker, call `Start`, assert no error and the marker created at the tip |
| ST3 | read the marker without the `^{commit}` peel | `TestStartIsNoOpWhenMarkerEqualsTip` | plant the marker at the tip, call `Start`, assert no error and no ref write |
| ST4 | treat an unreadable marker as absent and create a fresh one | `TestStartFailsClosedOnUnreadableMarker` | plant the marker, delete its object, call `Start`, expect a refusal and the marker unmoved |
| ST5 | leave the completion helper's empty expectation in place | `TestCompletionHelperPassesSameExpectation` | drive `Start` to a prepared start operation with an ancestor marker present, resume it, assert the same success and final marker position as ST1 |
| ST6 | pass the marker's ref name, which `update-ref` resolves and the double therefore accepts | `TestExpectationIsAFullObjectID` | capture the expectation the service hands the gate owner, assert it matches `^[0-9a-f]{40}$` for a readable marker and is exactly `""` for an unreadable one |
