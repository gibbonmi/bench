# Recover a dev manifest in the land route

Blocked by: publish-the-broker-from-the-stamped-build.md
Writes: bin/bench.sh, internal/systemtest/land_route_test.go, tests/canary/docs-currency-token-diet/missing-cli-inventory, tests/canary/docs-currency-token-diet/stale-cli-doc-reference, tests/canary/docs-currency-token-diet/stale-skill-cli-reference, tests/canary/load-validity-metadata/extensionless-gate-ref, tests/canary/package-core-guard/bounds-duplicate-owner, tests/canary/package-core-guard/reintroduced-bare-skip, tests/canary/package-core-guard/unrouted-subcommand
Covers: BF15, BF16, BF17, BF18

## What to build

Verify the premise first: `land_route` in bin/bench.sh refuses a version
mismatch and a digest mismatch through one exit 127 path each, after the
inherited-override refusals. Then split the version branch. When the manifest
version reads `dev`, run `scripts/go-build.sh` at the install root once, re-read
the manifest once, and continue when the version now matches. When it still
reads `dev`, exit 127 with the repair advice. The digest branch keeps its
unconditional exit 127. The recovery reads no repository and honors no
inherited override, so the override refusals stay first.

This ticket follows the publish ticket, because the rebuild it runs must write
a stamped manifest to recover.

Run the system-tagged tests with `BENCH_KIT` and `BENCH_RUN_BINARY` set, as
the system suite requires.

## Acceptance

- [ ] A `dev` manifest with a sound digest lands after one rebuild in the land-route table.
- [ ] A digest mismatch still exits 127 with the repair advice.
- [ ] An inherited `BENCH_KIT` still exits 1 before any rebuild.
- [ ] A manifest still `dev` after the rebuild exits 127.
- [ ] Self-probe: make the version branch rewrite the manifest's version field instead of rebuilding, and report the digest row red.
