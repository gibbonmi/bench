# Rebind the workflow fixture implementation sources

Blocked by: migrate-anchor-table.md
Ownership fence: `internal/conformance/registry_test.go`
Contracts: `internal/conformance/registry_test.go` binds the `workflow-guidance-anchors` fixture directory to the actual `internal/anchors` implementation and residual conformance adapter, asserted by AR2 against registry classification metadata

## What to build

Update the workflow-guidance fixture-family registration so its Go source inventory names the new matcher, registry, and registry-data implementation files alongside the residual conformance adapter.

## Acceptance

- [ ] [AR2] The fixture-family registry points readers and classification checks at every real Go implementation source that grades workflow-guidance anchors.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AR2 | drop one `internal/anchors` implementation source from the family row | coordinator registry metadata probe | drop → run classification/source audit → reject incomplete seam → restore source |
