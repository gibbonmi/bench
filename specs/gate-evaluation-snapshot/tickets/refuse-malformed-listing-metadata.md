# Refuse malformed listing metadata at capture

Blocked by: none
Ownership fence: `internal/gate/`
Contracts: none crosses — the parsed-entry validity rule stays inside `internal/gate/tree_snapshot.go`'s generation contract and is asserted by RM1-RM2 through its focused parser and capture tests

## What to build

Close review finding C1-malformed-listing-metadata-authorizes-generation: `parseTreeSnapshot` accepts any tab-separated entry, so a listing entry whose metadata is not git's `<octal mode> <object type> <hex object>` shape still authorizes a generation and only partially surfaces at blob-read time. Validate metadata at parse so a malformed listing refuses the generation in the existing fail-closed direction, while real working-tree and prospective captures keep parsing.

## Acceptance

- [ ] [RM1] A listing entry whose metadata is not exactly an octal mode, a known object type, and a hex object id refuses the parse, so generation capture returns an error and authorizes nothing.
- [ ] [RM2] Real captures over a tree holding regular, executable, and symlink entries still produce a generation, and blob requests retain the existing non-blob refusal.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RM1 | drop the metadata field validation so any tab-separated entry parses again | the malformed-metadata parser test | apply the mutation, run `go test ./internal/gate -run TreeSnapshot`, expect the authorized-generation failure |
| RM2 | reject symlink mode 120000 in the validation | the real-capture parser test | apply the mutation, run `go test ./internal/gate -run TreeSnapshot`, expect the refused-valid-tree failure |
