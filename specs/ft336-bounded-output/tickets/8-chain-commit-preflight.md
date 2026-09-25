# 8. Chain the worktree build and build preflight after a commit

Blocked by: none
Writes: internal/commit/commit.go, internal/commit/chain_grammar_test.go (new), cmd/bench/main.go, tests/canary/package-core-guard/unrouted-subcommand, cmd/bench/commit_chain_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: BO51, BO52, BO53, BO54, BO55, BO56, BO71

## What to build

Add `--preflight-build <slug>` to the `bench commit` grammar. The `cmd/bench` commit handler composes three steps in order: the commit, the worktree build of the current assignment, and `bench preflight build <slug>` at the published tip. The commit package returns the published commit to the command layer, and the chain line takes its sha from that value. Each step prints its own response. The chain then prints `commit-chain{commit=<sha|none>,build=<green|red|skipped>,preflight=<green|red|skipped>}`.

A commit that publishes nothing skips both later steps. A red build skips the preflight. The chain returns the first non-zero step exit, or 0. A commit that exits 3 stops the chain, and the chain line names its published commit. `--dry-run` with `--preflight-build` is a usage refusal at exit 2. A commit without the flag keeps its current behavior.

## Acceptance

- [ ] A green commit, build, and preflight end with `commit-chain{commit=<sha>,build=green,preflight=green}` at exit 0.
- [ ] A lane refusal exits 1, calls no build, and prints `commit-chain{commit=none,build=skipped,preflight=skipped}`.
- [ ] A red build exits with the build code, calls no preflight, and names the published commit.
- [ ] A red preflight after a green build exits 1.
- [ ] A commit that exits 3 calls no build and prints `commit-chain{commit=<sha>,build=skipped,preflight=skipped}` at exit 3.
- [ ] `--dry-run` with `--preflight-build` exits 2 with the commit usage line.
- [ ] The existing commit tests pass unchanged.
