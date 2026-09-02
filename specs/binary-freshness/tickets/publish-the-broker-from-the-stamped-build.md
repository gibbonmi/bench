# Publish the broker from the stamped build

Blocked by: add-seal-and-broker-rows-to-doctor.md
Writes: internal/freshness/freshness_publish.go, internal/freshness/freshness_publish_test.go, internal/freshness/publication_topology_test.go, internal/brokermanifest (new), internal/adopt/broker_test.go, scripts/go-build.sh, cmd/bench/main.go, cmd/bench/freshness_publish.go, cmd/bench/freshness_publish_test.go, cmd/bench/build_artifact_mode_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: BF10, BF11, BF12, BF13, BF14

## What to build

Verify the premise first: subject mode in scripts/go-build.sh runs the staged
binary's `freshness-publish`, which publishes the executable and its seal
under one rollback transaction. `freshnessPublish` in
cmd/bench/freshness_publish.go takes exactly two arguments. The manifest
writer binds `os.Executable`, which in subject mode is the staged binary.

Then give `freshness-publish` two more arguments: the published path and the
version. Extend the transaction so it writes the manifest through
`internal/brokermanifest` with that path and version, beside the resolved
wrapper as today. A rollback leaves the old manifest in place. Keep one exec
line in the build script. Artifact mode stays excluded, and a new test in
`cmd/bench` proves it writes no manifest and executes nothing.

This ticket follows the doctor ticket because it imports the leaf package the
doctor ticket creates.

## Acceptance

- [ ] After a subject-mode build, the manifest digest equals the published executable's digest and its version equals the package version.
- [ ] A publication rolled back before the rename leaves the old manifest unchanged.
- [ ] The artifact-mode test still shows no execution and no manifest write.
- [ ] `TestFreshnessPublicationTopology` still passes with only the three allowed callers.
- [ ] Self-probe: write the manifest before the rename outside the transaction, and report the rollback test red.
