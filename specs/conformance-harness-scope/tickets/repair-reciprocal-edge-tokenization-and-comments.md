# Repair reciprocal-edge tokenization and comments

Blocked by: none
Ownership fence: `internal/specbuild/assign.go`, `internal/specbuild/dependency_edge_repro_test.go`, `internal/specbuild/refresh_repin_test.go`
Integration surfaces: Integration-surfaces basename grammar→`internal/specbuild/assign.go` + RE1; assign and refresh dependency fixtures→`internal/specbuild/dependency_edge_repro_test.go` + RE1/RE2; re-pin contract fixture→`internal/specbuild/refresh_repin_test.go` + RE2
Contracts: parsed sibling surface text crosses `internal/specbuild/assign.go`→exact dependent-basename classification, asserted by RE1 in `internal/specbuild/dependency_edge_repro_test.go`; durable invariant prose crosses the same production rule→`internal/specbuild/dependency_edge_repro_test.go` and `internal/specbuild/refresh_repin_test.go`, asserted by RE2 review
Closure: RE1/exact-basename-token, RE1/near-name-non-edge, RE1/prose-non-edge, RE2/timeless-comments

## What to build

Recognize a dependent ticket only when its basename appears as an exact
Integration-surfaces token, not as a substring of another basename or incidental
prose. Retain the exact dependent-edge refusal and repaired-edge success controls.
Rewrite the new dependency and re-pin test comments as timeless current
constraints without naming the incident, a miss, a defect, or review provenance.

## Acceptance

- [ ] [RE1] (covers local) reciprocal-edge validation recognizes an exact dependent ticket basename but ignores a near-name basename such as `stone.md` when assigning `one.md` and ignores an incidental prose mention that is not a dependent target.
- [ ] [RE2] (covers local) dependency and re-pin test comments explain only the enduring rule and failure mode, with no change-history or review-provenance narration.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RE1/exact-basename-token | stop recognizing the exact `consumer.md` dependent target | the reciprocal-edge assign control | apply the omission, assign the unblocked consumer, require the missing-edge refusal to disappear incorrectly, then restore exact matching |
| RE1/near-name-non-edge | match a ticket basename by raw substring | the near-name control | name `stone.md` in a sibling surface, assign `one.md`, require an erroneous blocker refusal, then restore token matching |
| RE1/prose-non-edge | treat an incidental prose mention as a dependent target | the prose control | mention `one.md` outside the dependent-target position, assign it, require an erroneous blocker refusal, then restore grammar-aware matching |
| RE2/timeless-comments | restore incident or defect provenance in one repaired comment | the comment residue audit | apply the stale phrase, search the owned test comments for incident narration, require the residue audit to fail, then restore current-state prose |
