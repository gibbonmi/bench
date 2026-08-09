# Validate registry completeness and inertness

Blocked by: declare-production-axi-registry.md
Ownership fence: `cmd/bench`, `internal/conformance`, `projects/benchkit.md`
Integration surfaces: production declarations→declare-production-axi-registry.md; validation advertisement→`projects/benchkit.md`; guidance→document-ten-principle-axi.md
Contracts: declaration set and raw command observation cross `cmd/bench`→`internal/conformance`, domain is complete unique classified membership, order follows registry, and missing or parallel advertisements refuse, asserted by RV1
Closure: RV1/missing, RV1/duplicate, RV1/unclassified, RV1/parallel, RV1/inert

## What to build

validation rejects missing, duplicate, unclassified, or parallel members while metadata stays byte-inert.

## Acceptance

- [ ] [RV1] (covers CR6) validation rejects missing, duplicate, unclassified, or parallel members while metadata stays byte-inert.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RV1/missing | remove one declaration | registry conformance test | validate and require missing member |
| RV1/duplicate | duplicate one name | registry conformance test | validate and require duplicate refusal |
| RV1/unclassified | clear disposition | registry conformance test | validate and require classification refusal |
| RV1/parallel | add one independent advertisement | registry conformance test | enumerate and require sole production owner |
| RV1/inert | route metadata into stdout | compatibility oracle | run the affected command and require exact mismatch |
