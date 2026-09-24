# 1. Bound the exec child output through one response owner

Blocked by: none
Writes: internal/bounds/bounds.go, tests/canary/package-core-guard/bounds-duplicate-owner, internal/responsebound/ (new), cmd/bench/main.go, tests/canary/package-core-guard/unrouted-subcommand, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, cmd/bench/response_bound_test.go (new), internal/systemtest/exec_bound_test.go (new)
Covers: BO1, BO2, BO3, BO4, BO5, BO6, BO7, BO14, BO15, BO18, BO19, BO20, BO21, BO22, BO23, BO24, BO25, BO26, BO27, BO28, BO29, BO30, BO31, BO68

## What to build

Add the response line value, 10, to the production policy registry in `internal/bounds`. Add the response owner package. The owner takes tagged stdout and stderr writes in arrival order. It replays a response of 10 lines or fewer to the original streams. It prints the first 4 lines, the spill line, and the last 5 lines of a longer response on stdout.

The owner creates the spill store under the Bench home as the spec states. It keeps the head and a tail ring in memory after the spill starts. It takes the create-failure route and the midway-failure route that the spec states.

Give each public registry entry in `cmd/bench` a bound disposition. Declare the `worktree exec` leaf bounded. Declare every other public entry `pending`. `Command.Run` gives a bounded command one owner for its stdout and its stderr. It finishes the owner after the command returns, and it returns the command's exit code. A registry test refuses a public entry with no disposition.

`bench worktree exec` then passes its child's streams through the owner, because the dispatcher hands exec the owner's writers. A Bench verb inside the child keeps its own process bound.

## Acceptance

- [ ] A test-registered bounded command that prints 25 lines through `Command.Run` produces exactly 10 stdout lines.
- [ ] A bounded command that prints 10 lines produces its exact bytes and no spill file.
- [ ] An 11-line response prints lines 1 to 4, the spill line, and lines 7 to 11.
- [ ] A 25-line response of 300 bytes prints `spilled{lines=25,bytes=300,omitted_lines=16,path=<abs>}` as its fifth line.
- [ ] The spill file holds the complete input bytes, both streams, in write order.
- [ ] A bounded command that prints 30 lines and exits 3 returns exit 3.
- [ ] Each spill directory has mode 0700, each spill file has mode 0600, and a symlink at the scope directory takes the create-failure route.
- [ ] A file at the scope directory path gives the complete output and then `spill-failed{reason=<reason>}`.
- [ ] A write fault after 100 spill bytes leaves a 100-byte file and a response that holds every later byte after the `spill-failed` line.
- [ ] `bench consumers` for the line value lists only the owner package and the policy registry.
- [ ] A public registry entry without a disposition fails the registry test.
- [ ] An empty response prints nothing, and an unterminated eleventh line spills with its exact final bytes.
- [ ] A NUL byte and invalid UTF-8 reach the spill file byte-exact, and a Bench home with a newline takes the create-failure route.
- [ ] `bench worktree exec <target> -- sh -c 'seq 1 40'` prints lines 1 to 4, the spill line, and lines 36 to 40.
- [ ] An exec child that prints 40 lines and exits 7 makes exec exit 7.
- [ ] A heredoc on exec stdin reaches the child, and its 12 output lines spill.
- [ ] A SIGINT after the child printed 12 lines gives exit 130 and a spill file that holds those lines.
- [ ] A nested `bench` child that prints 30 lines gives exec 10 lines with the child's own spill line.
- [ ] An exec child that prints 64 MiB completes, and the owner holds only the head lines and the tail ring after the spill starts.
- [ ] The exec grammar refusal keeps its `usage: bench worktree exec` line.
