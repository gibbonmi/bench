# Declare production AXI registry metadata

Blocked by: introduce-outcome-carrier.md, introduce-action-carrier.md, introduce-projection-carrier.md, introduce-aggregate-empty-carriers.md
Ownership fence: `cmd/bench/main.go`, `cmd/bench/command_registry.go`, `internal/usage`, `internal/spec`, `internal/publication`, `internal/preflight`, `internal/gate`, `internal/harness`
Integration surfaces: carrier types→introduce-outcome-carrier.md, introduce-action-carrier.md, introduce-projection-carrier.md, introduce-aggregate-empty-carriers.md; validation→validate-axi-registry.md; disposition→assert-approved-axi-disposition.md
Contracts: member name, attachment, AXI disposition, grammar/help owner, renderer family, empty class, detail modes, and route set cross executable registries→`cmd/bench/command_registry.go`, membership is every root/nested production member, order is production order, and no field may be absent, asserted by RD1
Closure: RD1/root, RD1/nested, RD1/attachment, RD1/disposition, RD1/grammar, RD1/renderer, RD1/empty, RD1/detail, RD1/routes

## What to build

one production declaration enumerates every root and nested member with complete AXI metadata.

## Acceptance

- [ ] [RD1] (covers CR5) one production declaration enumerates every root and nested member with complete AXI metadata.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RD1/root | omit one root registry member | registry enumeration test | derive and require all 48 roots |
| RD1/nested | omit one nested operation | registry enumeration test | derive and require each production grammar member |
| RD1/attachment | clear one attachment | registry validation test | require missing attachment refusal |
| RD1/disposition | default an empty disposition | registry validation test | require explicit classification |
| RD1/grammar | clear one grammar owner | registry validation test | require missing owner refusal |
| RD1/renderer | clear one renderer family | registry validation test | require missing renderer refusal |
| RD1/empty | clear one empty class | registry validation test | require explicit empty/no-empty |
| RD1/detail | omit one applicable detail mode | registry validation test | require mode completeness |
| RD1/routes | omit one applicable shared route | registry validation test | require route completeness |

