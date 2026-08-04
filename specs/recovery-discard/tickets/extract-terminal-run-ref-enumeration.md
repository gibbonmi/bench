# Classify a terminal run's dead refs by assignment identity

Blocked by: none
Ownership fence: `internal/specbuild`
Contracts: the classified ref inventory crosses the new enumeration→`reclaimProvisionalRefs` inside `internal/specbuild/state.go`, asserted by TR1 and TR5 against real run records written through the package's own fixture harness, over the whole enumerated classification set rather than one sampled class
Assumptions: run records and their enumeration stay unexported in `internal/specbuild`, which already imports `internal/worktree`, so the pass cannot live in `internal/worktree` without an import cycle; terminality is the record's own flag, never a name pattern or a merged-looking ref; claims re-derived from the tree at pickup

## What to build

`reclaimProvisionalRefs` is a delete loop, not a producer of a ref list, and it
consults only the `Branch` field a recent commit added — so it reclaims nothing
at all for any record written before that commit. Extract the one function that
answers "which refs does this terminal run no longer need", and rewrite
`reclaimProvisionalRefs` to consume it. Two independent lists is the
duplicated-knowledge defect the code standard forbids, and the field-only
derivation is exactly the degenerate this ticket exists to kill.

Location is by assignment identity. Every record persists the assignment `ID` in
both the pre- and post-`Branch` shapes; the branch is located by matching that ID
against refs in the assignment namespace. The stored `Branch` is a fast path when
present, never the only path. Location is not classification: terminality still
comes from the record.

The inventory classifies rather than just lists, and three classes are retained,
never deleted:

- refs whose owning record is terminal — reclaimable;
- refs whose owning record exists but is not terminal — reported and retained,
  because deleting a live build's candidate destroys in-flight work;
- refs with no owning record at all — reported unclassified and retained,
  because an absent record cannot prove the work is dead;
- an assignment ID matching more than one ref in the namespace — reported
  ambiguous and retained, because the record persists the ID but not the owner
  half of the ref path, so both halves of an ambiguous match are unclaimed.

This ticket lands the enumeration and its consumer. The maintainer-facing verb
that plans and applies over it is the next ticket.

## Acceptance

- [ ] [TR1] the enumeration locates a terminal run's assignment branch for a record that persists no branch name, by matching the assignment ID against the assignment namespace.
- [ ] [TR2] the enumeration returns a terminal run's assignment branches, candidate ref, and checkpoint refs as one classified inventory.
- [ ] [TR3] refs belonging to a run whose record is not terminal are reported in the inventory and classified retained.
- [ ] [TR4] refs whose owning record is absent are reported unclassified and classified retained.
- [ ] [TR5] an assignment ID matching more than one ref in the assignment namespace is reported ambiguous and classified retained, and every non-reclaimable class the inventory can emit is asserted, not sampled.
- [ ] [TR6] `reclaimProvisionalRefs` consumes this enumeration and deletes exactly its reclaimable class, so promotion's existing reclamation keeps working and now also reaches records that persist no branch name.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TR1 | derive the branch from the stored `Branch` field alone | the no-branch-field enumeration test | drop the namespace match, run `go test ./internal/specbuild -run Reclaim -timeout 180s`, expect the located-branch assertion to fail on the pre-fix record |
| TR2 | omit the candidate ref from the returned inventory | the terminal-inventory test | drop the candidate append, run `go test ./internal/specbuild -run Reclaim -timeout 180s`, expect the inventory-membership assertion to fail |
| TR3 | classify every located ref as reclaimable regardless of the record's terminal flag | the non-terminal retention test | drop the terminal check, run `go test ./internal/specbuild -run Reclaim -timeout 180s`, expect the retained assertion to fail |
| TR4 | classify a ref with no owning record as reclaimable | the unclassified retention test | default the missing-record branch to reclaimable, run `go test ./internal/specbuild -run Reclaim -timeout 180s`, expect the retained assertion to fail |
| TR5 | take the first match when an ID resolves to several namespace refs | the ambiguous-match test | replace the multi-match guard with a first-match pick, run `go test ./internal/specbuild -run Reclaim -timeout 180s`, expect the ambiguous-and-retained assertion to fail |
| TR6 | leave `reclaimProvisionalRefs` on its own field-only loop beside the new enumeration | the promotion-reclamation regression test | restore the original loop body, run `go test ./internal/specbuild -timeout 180s`, expect the post-promotion empty-namespace assertion to fail for the pre-fix record |
