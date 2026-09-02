# Verify identity in commands brief

Blocked by: publish-the-broker-from-the-stamped-build.md
Writes: cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: BF1, BF2

## What to build

Verify the premise first: `commandsCommand` in cmd/bench/main.go returns a
fixed three-line string and reads no filesystem. Then convert the handler to
the root-taking shape `testCommand` uses. Resolve the running executable
through `os.Executable` and symlink resolution, and call `freshness.Verify`
against the kit root. On a mismatch, refuse with one `toon.Errorf` line whose
detail is `freshness.RebuildAction`. Outside a repository, keep the three-line
answer. Keep the verb on the binary route in bin/bench.sh.

This ticket follows the publish ticket, because both edit `cmd/bench/main.go`.

## Acceptance

- [ ] A behavior test with a mismatched seal receives a non-zero exit and the `RebuildAction` sentence.
- [ ] A behavior test with a matching seal receives the three-line answer at exit 0.
- [ ] The help row test in `main_test.go` still passes.
- [ ] Self-probe: inline the rebuild text instead of calling `RebuildAction`, and report which sweep reds it or that none does.
