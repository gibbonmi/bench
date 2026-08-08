# Prepare gate artifacts before scheduling

Blocked by: share-artifacts-across-local-processes.md, resolve-target-aware-artifact-identities.md, require-existing-gate-artifact-authority.md, run-promoted-gate-bootstrap.md
Ownership fence: `.bench/gate.sh`, `.bench/gate-prospective.sh`, `scripts/go-build.sh`, `internal/gate/gate_go.go`, `internal/gate/gate_go_test.go`, `internal/gate/phases.go`, `internal/gate/phases_command_test.go`, `internal/conformance/gate_entry_test.go`, `internal/canary/gate_entry_test.go`
Integration surfaces: real artifact-store producer→share-artifacts-across-local-processes.md; selected CLI/verifier records→gate entry, GateGoArgv, build publication, migrate-gate-fixture-artifact-consumers.md, and admit-one-gate-per-common-repository.md; release/install entry consumers→migrate-preprelease-install-proof-consumers.md; remaining raw construction→enforce-the-executable-artifact-contract.md
Contracts: `ArtifactRecord` crosses `internal/artifactstore/`→gate preparation in `internal/gate/phases.go` as CLI or verifier class with matching target identity, records resolve before phase-table construction, missing/stale/malformed records refuse with the rebuild or exact-repair action and never fall back, asserted by GP1-GP2 against the real store and phase table
Closure: GP1/verifier-single-build, GP1/verifier-currency, GP1/refusal-before-cli, GP1/no-verifier-go-run, GP2/prepared-before-schedule, GP2/build-phase-retained, GP2/conformance-reader-edge, GP2/contract-reader-edge, GP2/shellcheck-reader-edge, GP2/canary-reader-edge, GP2/gofmt-no-build-edge, GP2/vet-no-build-edge, GP2/test-no-build-edge, GP2/race-no-build-edge, GP2/conformance-suite-no-build-edge

## What to build

Prepare and verify the CLI and freshness-verifier artifacts before scheduling. Execute gate plumbing through the selected CLI, retain the build phase as the root-local publisher, and preserve every existing reader edge without adding a toolchain-to-build edge.

## Acceptance

- [ ] [GP1] (covers SB7) one verifier artifact is reused by direct and nested entries, verifier source movement invalidates it, stale CLI is refused before `gate-phases`, and no entry invokes `go run ./internal/freshness/check`.
- [ ] [GP2] (covers SB9) story 1 alone prepares artifacts before scheduling, retains the build phase and its reader edges, and runs toolchain plumbing through the selected CLI with no new build edge.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GP1/verifier-single-build | build the verifier inside each gate entry | nested-entry counting-verifier test | execute direct plus nested entries and expect the verifier-build-count-equals-one failure |
| GP1/verifier-currency | omit verifier sources from its identity | verifier-source mutation test | change verifier source after preparation and expect entry to refuse the old identity before CLI execution |
| GP1/refusal-before-cli | execute the stale CLI before verifier success | hostile-CLI gate-entry test | plant a CLI marker and failing verifier, invoke entry, and expect no CLI/phase marker |
| GP1/no-verifier-go-run | restore the assembled `go run ./internal/freshness/check` entry | structural gate-entry test | inspect the parsed entry command and expect the forbidden verifier-construction diagnostic |
| GP2/prepared-before-schedule | resolve the artifact from inside the first phase | phase-construction ordering test | construct the table with a missing record and expect preparation refusal before any phase marker |
| GP2/build-phase-retained | drop the build phase after preparation | real phase-table inspection | resolve a kit-shaped table and expect the missing-build-phase failure |
| GP2/conformance-reader-edge | remove conformance's build need | phase dependency enumeration test | inspect conformance and expect the omitted-edge failure |
| GP2/contract-reader-edge | remove contract's build need | phase dependency enumeration test | inspect contract and expect the omitted-edge failure |
| GP2/shellcheck-reader-edge | remove shellcheck's build need | phase dependency enumeration test | inspect shellcheck and expect the omitted-edge failure |
| GP2/canary-reader-edge | remove canary's build need | phase dependency enumeration test | inspect canary and expect the omitted-edge failure |
| GP2/gofmt-no-build-edge | add the build need to gofmt | toolchain phase inspection | inspect gofmt and expect the extra-edge failure |
| GP2/vet-no-build-edge | add the build need to vet | toolchain phase inspection | inspect vet and expect the extra-edge failure |
| GP2/test-no-build-edge | add the build need to test | toolchain phase inspection | inspect test and expect the extra-edge failure |
| GP2/race-no-build-edge | add the build need to race | toolchain phase inspection | inspect race and expect the extra-edge failure |
| GP2/conformance-suite-no-build-edge | add the build need to conformance-suite | toolchain phase inspection | inspect conformance-suite and expect the extra-edge failure |
