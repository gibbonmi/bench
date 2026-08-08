# Run the promoted gate bootstrap

Blocked by: promote-exact-green-bootstrap.md
Ownership fence: `bin/bench.sh`, `cmd/bench/`, `internal/freshness/`, `.bench/BENCH-reference.md`, `internal/conformance/subcommand_routing_test.go`
Integration surfaces: promoted bootstrap producer→promote-exact-green-bootstrap.md; installed-only shell route and plumbing inventory→`bin/bench.sh`, `cmd/bench/`, `.bench/BENCH-reference.md`, and `internal/conformance/subcommand_routing_test.go`; strict CLI/verifier authority→`internal/freshness/`; authenticated direct-entry consumer→prepare-gate-artifacts-before-scheduling.md
Contracts: exactly one exact-green promoted launcher crosses its fixed non-ambient trust-root slot in `bin/bench.sh`→the `cmd/bench/` bootstrap entry, then exactly one current CLI `ArtifactRecord` and one current verifier `ArtifactRecord` cross strict authority in `internal/freshness/`→that trusted entry, membership is the closed CLI and freshness-verifier artifact classes for the current host target, ordering is reject an absent or hostile slot then launch the promoted trust root then derive both current identities then validate both existing records and authority then execute the verifier then execute the CLI, and absent, malformed, corrupt, stale-CLI, stale-verifier, or failed-verifier state refuses before builder, Go, released-cache, selected CLI, phase, repair, network, or store writes, asserted by RB1-RB3
Closure: RB1/missing-slot-refusal, RB1/slot-type-refusal, RB1/malformed-record-refusal, RB1/corrupt-record-refusal, RB1/no-released-cache-fallback, RB1/no-dist-fallback, RB1/no-go-fallback, RB2/current-cli-identity, RB2/current-verifier-identity, RB3/authority-before-exec, RB3/verifier-before-cli, RB3/no-authoring-fallback

## What to build

Expose the promoted bootstrap through one plumbing command and an installed-only shell resolver. The resolver selects only the distinct bootstrap trust-root slot from the blocker ticket, rejects a missing, symbolic-link, special, or non-executable entry, and never tries ordinary released-cache, repair, repository, Go, or network routes. Exact-green promotion is the authentication boundary for that launcher; hostile replacement of the trust root after promotion requires package signing or an OS integrity mechanism and is explicitly outside this in-process contract. Once launched, the trusted executable strictly loads its record and the current two-record authority, validates the selected CLI and verifier without authoring, runs the verifier before the CLI, and remains dormant until GP1 changes the ordinary gate entry.

## Acceptance

- [ ] [RB1] (covers local) bootstrap execution resolves only the fixed promoted trust-root slot; missing, symbolic-link, special, non-executable, malformed-record, or corrupt-record state refuses without any ordinary binary, Go, repair, or network fallback.
- [ ] [RB2] (covers local) bootstrap execution separately derives and pins the current CLI identity and current freshness-verifier identity before either executable runs.
- [ ] [RB3] (covers local) bootstrap execution strictly loads authority before either executable, runs a successful validated verifier before the CLI, and never invokes an authoring builder on refusal.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RB1/missing-slot-refusal | fall through when the promoted slot is absent | installed-only resolver test | invoke bootstrap execution with ordinary binaries planted and expect missing-slot refusal with no executable marker |
| RB1/slot-type-refusal | accept a symbolic-link, FIFO, or non-executable trust-root slot | installed-only resolver hostile-slot test | replace the fixed slot with each hostile type, plant an executable symlink referent where applicable, and expect bounded refusal with no referent or fallback marker |
| RB1/malformed-record-refusal | accept an unknown, missing, or contradictory bootstrap record field | installed-only resolver record table | mutate each field class independently and expect exact-schema refusal before execution |
| RB1/corrupt-record-refusal | accept corrupt bootstrap record bytes | installed-only resolver test | replace the record with non-record bytes and expect refusal with no executable marker |
| RB1/no-released-cache-fallback | append the ordinary version-target cache to bootstrap resolution | released-cache sentinel fixture | remove the promoted slot, plant a released-cache marker, and expect no marker |
| RB1/no-dist-fallback | append repository `dist/bench` to bootstrap resolution | hostile candidate CLI fixture | remove the promoted slot, plant a candidate CLI marker, and expect no marker |
| RB1/no-go-fallback | reconstruct the verifier with `go` when bootstrap state refuses | hostile Go fixture | remove the promoted slot, plant a Go marker, and expect no marker |
| RB2/current-cli-identity | reuse authority without deriving the current CLI identity | CLI source-movement fixture | move selected CLI source after authority publication and expect stale-CLI refusal before execution |
| RB2/current-verifier-identity | reuse authority without deriving the current verifier identity | verifier source-movement fixture | move selected verifier source after authority publication and expect stale-verifier refusal before execution |
| RB3/authority-before-exec | execute either selected artifact before strict authority loading | hostile authority fixture | corrupt the authority record with verifier and CLI markers planted and expect neither marker |
| RB3/verifier-before-cli | execute the CLI before verifier success | failing-verifier fixture | make the validated verifier fail with a CLI marker planted and expect no CLI or phase marker |
| RB3/no-authoring-fallback | call the authoring resolver after strict lookup refuses | counting builder fixture | remove one selected record, invoke bootstrap execution, and expect zero builder calls and no store writes |
