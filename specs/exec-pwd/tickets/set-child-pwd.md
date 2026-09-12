# Set the exec child's PWD to its worktree

Blocked by: none
Writes: internal/worktree/exec.go, internal/worktree/exec_pwd_test.go (new), internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, CHANGELOG.md
Covers: none

## What to build

Give each exec child one PWD assignment that names its actual working directory.
The resolved assignment supplies that directory.
Inherited PWD values and repeated --env PWD values cannot change it.
Apply this rule with and without a worktree wrapper.

Keep the current child arguments, stdin, stdout, stderr, and exit behavior.
Keep the current BENCH_HOME and wrapper routing rules.
Keep unrelated inherited variables and explicit environment values.

Put the regression tests in exec_pwd_test.go.
Keep the existing serial ceiling and structure budgets.
Add a Fixed entry under a dedicated Exec child PWD changelog heading.

## Acceptance

- [ ] A child receives its worktree as PWD when the caller supplies a different inherited PWD.
- [ ] A child receives its worktree as PWD when inherited PWD is absent.
- [ ] Repeated --env PWD values cannot replace the child's worktree path.
- [ ] The child environment contains exactly one PWD assignment.
- [ ] These PWD guarantees hold with and without a worktree wrapper.
