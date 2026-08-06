# Single-source the fake-go builder stub

Blocked by: none
Ownership fence: `internal/freshness/publication_topology_test.go`, `internal/contract/surface/artifact/posture`
Integration surfaces: freshness-publish token-owner enumeration→`internal/freshness/publication_topology_test.go`; artifact trace and builder-behavior consumers→`internal/contract/surface/artifact/posture`; builder script unchanged→existing `scripts/go-build.sh` + FG2
Contracts: the fake-go stub's `go build`/`-o` argv contract crosses `internal/contract/surface/artifact/posture`→`scripts/go-build.sh`, asserted by FG2 against the real builder script

## What to build

Repair for review finding `standards-02-second-fake-go-fixture-harness`. The
composition authors two independent fake-`go` builder stubs that both fabricate
a bench-like output and re-derive the builder's `-o` argv parsing:
`fakeBuilderEnv`'s goStub in
`internal/contract/surface/artifact/posture/mode_test.go` and
`runArtifactPublicationTrace`'s goStub in
`internal/freshness/publication_topology_test.go`. The spec's testing decision
forbids a second fixture harness. Collapse to one source: relocate the two
artifact command-trace tests (`TestArtifactModeCommandTraceExcludesPublication`,
`TestArtifactModeCommandTraceBitesOnPublishThenDelete`) and their helpers from
the topology file into the posture package, and derive both consumers' stubs
from one definition there. The topology file keeps only the caller-enumeration
contract, which is its actual subject. Do not construct the `freshness-publish`
token dynamically anywhere — that would blind the topology walk's token-owner
enumeration rather than satisfy it; test files stay exempt from the walk by the
existing `_test.go` rule.

## Acceptance

- [ ] [FG1] Exactly one fake-`go` builder stub definition exists across the tree; the artifact command-trace contracts and the builder-behavior tables both drive it.
- [ ] [FG2] The relocated trace contracts keep their teeth: the publish-then-delete mutation of the real builder script still produces the staged-publication-execution red, and the exclusion contract still passes against the unmutated script.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FG1 | reintroduce an independent goStub literal in either consumer | the single-source sweep | run `rg -c 'go:%s' internal/freshness internal/contract/surface/artifact/posture` (the stub's trace-line stem), expect more than one defining site to fail the count the review re-derives |
| FG2 | apply the publish-then-delete mutation to the real builder script | the relocated bites contract | run the relocated `TestArtifactModeCommandTraceBitesOnPublishThenDelete`, expect the staged-publication-execution diagnostic |
