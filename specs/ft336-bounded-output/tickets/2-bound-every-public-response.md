# 2. Bound every public response and close the exemptions

Blocked by: 1-bound-exec-output.md
Writes: cmd/bench/, internal/systemtest/, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: BO8, BO9, BO10, BO11, BO12, BO13, BO30, BO66, BO70

## What to build

Declare every public registry entry bounded, except the closed exempt set. The exempt set holds these members, and each exemption names its reason:

- the help forms
- `bench dashboard --stdout`, `bench worktree shell`, and `bench setup`
- the ship-tier commands `release-preflight`, `prep-release`, and `release`
- the wrapper-only `repair`

A help form is `bench help`, or a command or a leaf followed by exactly one `--help`, `-h`, or `help` argument. Remove the `pending` disposition value.

After the disposition change, run the gate once and list each test that the bound turns red. Change each red test to read the spill file or to call its package seam. Keep each complete-output assertion complete. Change no assertion to a weaker predicate.

Do not edit `.bench/hooks/block-bench-follow-on.sh` or the chain sentence of `.bench/BENCH.md`.

## Acceptance

- [ ] `bench help` through `Command.Run` prints every line that `renderCommandHelp` renders.
- [ ] `bench worktree --help` through `Command.Run` prints its complete grammar.
- [ ] A bounded command given two arguments that end in `--help` stays bounded.
- [ ] `bench dashboard --stdout` through `Command.Run` prints the complete page.
- [ ] The registry's exempt set equals the closed members, and each exemption names a reason.
- [ ] `release-preflight` through `Command.Run` prints a 30-line response in full.
- [ ] The internal `guard-git` command prints a 30-line stderr in full.
- [ ] A nested `bench` child that prints 30 lines gives exec 10 lines with the child's own spill line.
- [ ] The diff leaves the follow-on hook and the chain sentence unchanged.
