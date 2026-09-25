# Grade the limit of a bounds read seam in the bounds-policy guard

Blocked by: none
Writes: internal/conformance/bounds_policy_test.go, tests/canary/package-core-guard/bounds-read-limit-restated/BASE, tests/canary/package-core-guard/bounds-read-limit-restated/CHECK, tests/canary/package-core-guard/bounds-read-limit-restated/EXPECT, tests/canary/package-core-guard/bounds-read-limit-restated/MUTATE.json, tests/canary/package-core-guard/bounds-classify-limit-restated/BASE, tests/canary/package-core-guard/bounds-classify-limit-restated/CHECK, tests/canary/package-core-guard/bounds-classify-limit-restated/EXPECT, tests/canary/package-core-guard/bounds-classify-limit-restated/MUTATE.json
Covers: none

## What to build

The bounds-policy guard grades the arguments of three standard-library calls only. A
caller can therefore restate a registry value through a read seam of the bounds package,
for example `bounds.Read(body, 5 << 20)`, and the guard stays green.

The guard grades the limit argument of each read seam of the bounds package. A read seam
is an exported function of that package with a `limit` parameter of type `int64`. The
guard gets the seams and the limit positions from the package source. When the text of a
limit argument is equal to the text of a registry value, the guard reports the caller.
The guard gets the registry texts from the same parse of the registry that it uses now.
A limit that names a registry entry, for example `bounds.ModelReadLimit`, stays green.

## Acceptance

- [ ] The canary `bounds-read-limit-restated` fails the guard when a caller gives `bounds.Read` a limit that restates a registry value, and the live tree passes `bench test --check bounds-policy`.
- [ ] The canary `bounds-classify-limit-restated` fails the guard when a caller gives `bounds.Classify` a limit that restates a registry value.
- [ ] A probe that omits the new read-seam grading turns the `bounds-read-limit-restated` canary green, and the restore returns the guard to its graded state.
