# Build unsealed artifacts without execution

Blocked by: publish-one-sealed-local-subject.md
Ownership fence: `scripts/go-build.sh`, `scripts/build-artifacts.sh`, `scripts/native-proof.sh`, `internal/contract/surface/artifact`
Integration surfaces: semantic mode grammar and unsealed promotion→`scripts/go-build.sh`; release artifact caller and seal-deletion removal→`scripts/build-artifacts.sh`; native-proof rebuild and later runner execution→`scripts/native-proof.sh`; artifact cache, reproducibility, staging, and non-host journeys→`internal/contract/surface/artifact`; sealed-publication operation→publish-one-sealed-local-subject.md
Contracts: the exact `artifact` mode token, ordered `<root> <output>` operands, explicit `GOOS`/`GOARCH` target, and absent-mode default cross artifact callers→`scripts/go-build.sh`; AR1, AR2, and AR4 assert the token domain is only `artifact`, mode precedes operands, duplicate or absent values refuse, target variables remain caller-selected only in artifact mode, and absence selects the sealed host-subject behavior

## What to build

Artifact and native-proof producers request an explicit artifact mode that compiles once, validates and atomically promotes the staged executable without running it, and guarantees no destination seal remains. Production callers state that intent once and stop deleting seals after the builder returns; native proof executes the output only later on its native runner.

## Acceptance

- [ ] [AR1] `--mode artifact <root> <output>` is the only added grammar; missing operands, missing or duplicate mode selectors, and unknown modes refuse before invoking Go, staging, or touching an existing output and seal.
- [ ] [AR2] Artifact mode invokes exactly one `go build`, honors the caller's explicit target environment, invokes no `go run` or Bench/helper publication operation, never executes the staged output, atomically promotes it, and leaves no destination seal.
- [ ] [AR3] Compile, validation, promotion, and bounded interruption failures preserve any prior output and prior seal byte-for-byte and remove the staged artifact.
- [ ] [AR4] Release artifact generation and native proof request `--mode artifact` directly and contain no post-build seal deletion; native proof executes the rebuilt output only later on the native runner.
- [ ] [AR5] Subject-to-artifact and artifact-to-artifact reruns converge on an unsealed executable while the existing shared-versus-private build-cache and reproducibility posture remains unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AR1 | treat an unknown, duplicate, or missing mode value as default subject mode | the builder argument table | fingerprint the prior pair, invoke every malformed form with a Go recorder, and expect zero Go calls and unchanged fingerprints |
| AR2 | run the cross-compiled output or publish it sealed and delete the seal afterward | the complete command recorder and non-host execution marker | build a deliberate non-host artifact, inspect every Go and Bench invocation, and expect any execution, sealed-publication reachability, or residual seal to fail |
| AR3 | remove the artifact cleanup trap or rename directly over the prior output before validation completes | the failure and interrupt table | block or fail each stage, signal where applicable, and expect the prior fingerprints and no staged artifact |
| AR4 | omit `--mode artifact` from either production caller or restore caller-side seal deletion | the caller-shape and real artifact/native-proof contracts | run both production journeys and inspect their builder argv and script shapes, expecting the default mode or deletion path to fail |
| AR5 | replace the existing cache selection with `GOCACHE=off` or a new per-build cache in shared mode | the artifact cache-posture and reproducibility contracts | run shared and private artifact journeys plus a two-mode rerun, expecting the established cache selection and reproducible bytes |
