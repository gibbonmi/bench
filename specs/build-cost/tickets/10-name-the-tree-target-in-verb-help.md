# Name the tree target in each verb's help

Blocked by: 06-pass-the-chunk-base-to-the-commit-chain.md, 07-stage-the-probe-in-record-ticket.md
Writes: cmd/bench/command_registry.go, cmd/bench/tree_scope.go, cmd/bench/tree_scope_test.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/commit_chain_test.go, cmd/bench/main_test.go, cmd/bench/record_operations_test.go, cmd/bench/selected_queries_test.go, cmd/bench/response_bound_exempt_test.go, cmd/bench/spec_retire_listing_test.go, cmd/bench/build_subject_mode_test.go, cmd/bench/worktree_leaves_test.go
Covers: BC58, BC59, BC60

## What to build

This rule applies to each tree-scoped definition and each tree-scoped family leaf that answers `--help` at exit 0. Its `--help` output names the `--in` flag with the operand that the root help shows. The text comes from `treeTargetInsertion`, which reads the same scope field as the dispatcher, so the root help and the verb help have one source.
`treeTargetInsertion` keeps its guard line `if definition.Scope != scopeTree {`, because the probe-10 command swaps that exact line.
A repository-scoped definition's `--help` names no `--in`.

The test enumerates the command registry, so a new tree-scoped verb gets the text with no test edit. Each test that compares a tree-scoped verb's whole help output follows the new text.

## Acceptance

- [ ] `bench commit --help` names the `--in` flag with its operand (BC58).
- [ ] Each tree-scoped definition and leaf in the registry whose `--help` answers at exit 0 names it in that output (BC59).
- [ ] No repository-scoped definition names it (BC60).
