# Refuse a failed Go build flag read

Blocked by: none
Writes: scripts/go-build.sh, cmd/bench/build_flag_read_test.go (new), cmd/bench/build_subject_mode_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: none

## What to build

The Go builder reads its build flags from the requirements registry through `node`. The builder reads them through a process substitution, and `set -e` does not see a failure there. When `node` fails, the flag list stays empty. Then `go build` runs without `-trimpath`, `-buildvcs=false`, or the version stamp, and the builder exits zero.

Make the builder read the flags through a command substitution, so that it sees the exit status of the read. A failed read stops the build with a message on stderr and a non-zero exit. An empty flag list also stops the build with a message on stderr and a non-zero exit. The build then writes no executable and no seal.

The regression test copies the kit root and puts a shim for node first on the path. The shim fails only the registry read. It runs the real node for all other reads.

## Acceptance

- [ ] When the flag registry read fails, `scripts/go-build.sh` exits non-zero and writes no executable and no seal.
- [ ] When the flag registry read fails, `scripts/go-build.sh` names the registry input on stderr.
- [ ] `TestGoBuildRefusesWhenFlagRegistryReadFails` fails when the builder reads the flags through the old process substitution.
- [ ] A build with a readable flag registry still produces the same executable as before.
