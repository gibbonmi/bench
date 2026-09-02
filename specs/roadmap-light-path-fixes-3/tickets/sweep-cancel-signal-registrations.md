# Sweep cancel-signal registrations

Blocked by: derive-the-canonical-path-in-one-leaf-package.md
Writes: internal/conformance/cancel_signal_registrations_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/registry.go, internal/skillsindex/command.go, internal/worktree/subshell.go, tests/canary/cancel-signal-registrations (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LQ21, LQ22, LQ23

## What to build

Verify the premise first: `signal.NotifyContext` in
internal/skillsindex/command.go registers `os.Interrupt` alone, and
`signal.Notify` in internal/worktree/subshell.go omits SIGHUP. Migrate both to
`subprocess.CancelSignals`. Then add one dev-tier check that mirrors
`checkBoundCallers` in internal/conformance/bounds_policy_test.go: production
files only, `internal/subprocess/` exempt, an AST walk over call expressions
whose selector is `signal.Notify` or `signal.NotifyContext`. The check reds a
call whose variadic argument is not `subprocess.CancelSignals`. Register it
after the canonical-path owner row with Go-source inputs. Mirror the row in
`checks_test.go`. Add a canary fixture whose source holds the token in a call,
a comment, and a string.

The blocker ticket rewrites `internal/worktree/subshell.go` and adds a row to
the check registry, so this ticket edits on top of both changes.

## Acceptance

- [ ] The check passes on the live worktree after the two migrations.
- [ ] The fixture reds on the call and stays green when only the comment and the string hold the token.
- [ ] A `_test.go` file that registers `os.Interrupt` alone does not red the check.
- [ ] `TestCancelSignalsMembershipIsExact` still passes.
- [ ] Self-probe: revert the skillsindex migration, and report the check red with that file named.
