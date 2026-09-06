# Refuse a file-tool write to a tracked path in the primary checkout

Blocked by: none
Writes: cmd/bench/guards.go (new), internal/conformance/entry_point_parity_bite_test.go (new), internal/conformance/harness_record_fixture_test.go (new), .bench/hooks/block-primary-file-write.sh (new), internal/writeguard/writeguard.go (new), internal/writeguard/writeguard_test.go (new), cmd/bench/main.go, cmd/bench/command_registry_test.go, .claude/settings.json, internal/harnesses/harnesses.go, internal/harnesses/harnesses_test.go, internal/conformance/harness_record_test.go, internal/otelrecord/registry.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/entry_point_parity_test.go, .bench/BENCH-reference.md, README.md
Covers: none

## What to build

The primary checkout takes writes only through a landing. The Bash guards refuse the raw git route there, and `bench commit` refuses the primary checkout. The Claude Code file tools have no guard, so an Edit or a Write to a tracked path in the primary checkout lands silently.

Add one guard to the existing family. The shim `.bench/hooks/block-primary-file-write.sh` carries the four-key manifest header and sources `.bench/lib/resolve-bench.sh`. It resolves the wrapper and pipes the PreToolUse envelope to the new plumbing verb `bench guard-file-write`. The shim warns and allows when the core is missing, stale, or errored, because a file edit outside a repository must stay open. It passes exit 0 and exit 2 through unchanged, as `block-bench-follow-on.sh` does.

The verb reads `tool_input.file_path` from the envelope. The new package `internal/writeguard` owns the envelope read and the verdict. It takes injected functions for the repository root, the primary-checkout test, the tracked test, and the ignored test. The verb wires those functions to `git.RootAt`, `git.IsPrimaryCheckout`, a `git ls-files --error-unmatch` read, and a `git check-ignore -q` read.

The refusal line starts with `BLOCKED:` and names the path, then it appends `usage.PrimaryCheckoutRefusal()`. The verb registers as a hook command with `attachmentDirect`, `axiExempt(axiReasonPlumbing)`, and `internalInventory`. It opens a `hook.guard-file-write` span through `beginHookSpan`.

Wire the shim in `.claude/settings.json` under a PreToolUse group with the matcher `Edit|Write|MultiEdit|NotebookEdit`. Add that group as `PreToolUse:Edit|Write|MultiEdit|NotebookEdit` to the claude row's `HookEvents`. Leave `.codex/hooks.json` unchanged, because Codex has no file-tool hook surface. Add the parity row for the shim and the routing exemption for the verb. Add the registry pin and the otel seam row. Add one bullet under Hook Layers in the reference and one row in the README guard table.

## Acceptance

- [ ] [G1] A tracked path under the primary checkout exits 2 with a `BLOCKED:` line that names the path and the Bench worktree route.
- [ ] [G2] An envelope whose `file_path` is the same relative path under a Bench worktree exits 0.
- [ ] [G3] An envelope whose `file_path` is a git-ignored path under the primary checkout exits 0.
- [ ] [G4] An envelope whose `file_path` is outside every repository exits 0.
- [ ] [G5] An envelope with no readable `file_path` exits 0 with a warning line. The shim exits 0 with a warning when the core is missing.
- [ ] [G6] `bench guards` lists `block-primary-file-write` with the boundary `PreToolUse:Edit|Write|MultiEdit|NotebookEdit` and `claude` in its wired cell.
- [ ] [G7] `bench test --check harness-record` passes with the new event group, and its biting probe names the script when the group is removed.
- [ ] [G8] The entry-point parity check passes for the new shim, and the subcommand routing check passes for the new verb.
- [ ] [G9] `.codex/hooks.json` is byte-identical to its state before the ticket.

## Build decisions

The coordinator amended row G7 before the commit. Its second clause named a hook-events field in `bench harnesses claude`, and that view renders no such field. The row now grades the conformance check and its biting probe. The `Writes:` line gained the two test files whose independent counts the new event group changes.
The growth lane refused three over-budget targets. The guard functions, the parity bite test, and the harness-record fixture builders moved to sibling files in the same change.
