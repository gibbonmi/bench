# 4. Print slot actions for the active worktree rows

Blocked by: none
Writes: internal/worktree/list.go, internal/worktree/list_actions_test.go, internal/worktree/path_identifier_test.go, .agents/skills/bench-craft-cli/SKILL.md, internal/anchors/registry_retained_workflow.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: BO32, BO33, BO34, BO35, BO36, BO41, BO67

## What to build

When one or more active rows with a present tree exist, `bench worktree list` prints one `bench worktree path <target>` action and one `bench worktree exec <target> -- <command>` action. No help row names an active row's id. The cleanup-pending, missing-tree, and foreign rows keep their row-specific actions.

Rewrite `TestListPathActionRunsAsAdvertised` so that it reads the id cell and passes it to `bench worktree path`. Change the `bench worktree list` row of the craft-cli table to state the target slot rule, and add one anchor needle that pins that sentence. Keep the general sentence that holds `per matching row`. Do not change `bench consumers`.

## Acceptance

- [ ] A list of 3 active rows prints `help[2]{cmd,why}:` with the two slot actions and no active id in a help row.
- [ ] A cleanup-pending row keeps its release help row.
- [ ] A missing-tree row keeps its recovery help row.
- [ ] A foreign row keeps its clean help row.
- [ ] The id cell of an active row passes `bench worktree path` at exit 0.
- [ ] The anchor check refuses the craft-cli file when the slot sentence is removed.
- [ ] An ambiguous `bench consumers` name keeps one re-query row per candidate.
