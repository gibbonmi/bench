# Migrate evaluation identities atomically

Blocked by: expand-generation-source.md
Ownership fence: `internal/gate/`
Contracts: the generation from `internal/gate/tree_snapshot.go` crosses identity derivation→selection and evidence handling inside `internal/gate/`, asserted by MI1-MI6 against the public Execute and ExecuteTree journeys

## What to build

Make one evaluation owner distribute one accepted pre-execution generation to component, conformance-check, conformance-canary, whole-subject, and stripped-subject identity derivation, then capture a distinct post-execution generation for complete drift and evidence handling. Ordinary and prospective execution use their respective adapters while retaining all current policy and public behavior.

## Acceptance

- [ ] [MI1] Component, conformance-check, conformance-canary, whole-subject, and stripped-subject identities consume one accepted pre generation while retaining their literal identity values, partitions, and reduced-scope behavior.
- [ ] [MI2] Under-lock validation and the distinct post generation retain tree, oracle closure, resolution, and passlisted-environment drift refusal; tracked and untracked unignored mid-run edits author no green or reusable evidence.
- [ ] [MI3] ExecuteTree reads the exact unpublished supplied tree through the prospective source, leaves the ordinary checkout untouched, runs the complete applicable inventories, and preserves prospective bootstrap evidence.
- [ ] [MI4] One ordinary real execution performs at most three working-tree materializations and two parsed listings, independent of identity-family breadth.
- [ ] [MI5] One prospective real execution performs zero pre working-tree materializations, at most one post materialization, and at most two parsed listings while each distinct requested blob is read once per generation.
- [ ] [MI6] Source, listing, blob, non-blob, symlink, escaping-target, and unavailable-tree failures retain the existing refusal or complete-inventory widening posture and never authorize a partial identity map.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MI1 | restore an independent capture in the conformance-check family | the evaluation family-enumeration and count test | restore the capture, run the focused evaluation test, expect the listing ceiling to fail |
| MI2 | reuse the pre generation as the post generation | the tracked and untracked drift journey tests | make the child mutate each subject kind, run Execute, expect the action to refuse despite the child's zero exit |
| MI3 | route prospective pre identity through the checkout working-tree source | the exact unpublished-tree and prospective-count tests | run ExecuteTree over an unpublished tree, expect exact-tree or zero-pre-materialization failure |
| MI4 | recapture one ordinary family from the root | the ordinary source-cost recorder | run one stable Execute journey, expect materialization or listing count above the formula |
| MI5 | reread a shared prospective blob for a second identity family | the prospective source-cost recorder | run one stable ExecuteTree journey, expect the per-generation object read count above one |
| MI6 | return a shortened snapshot after one source fault | the adapter fault table | inject each fault, run the evaluation journey, expect refusal or full execution rather than inherited evidence |
