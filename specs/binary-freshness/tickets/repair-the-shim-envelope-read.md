# Repair the shim envelope read

Blocked by: none
Writes: .bench/lib/resolve-bench.sh, .bench/hooks/block-bench-follow-on.sh, .bench/hooks/block-dangerous-git.sh, .bench/hooks/session-start.sh, internal/systemtest/bench_follow_on_test.go, internal/conformance/guard_classifier_table_test.go, tests/canary/guard-classifier-table/word-test-drops-xargs, internal/conformance/registry_test.go
Covers: BF22, BF23, BF24

## What to build

This ticket repairs review findings S1, C2, F5, and S2.

The shim in `.bench/hooks/block-bench-follow-on.sh` defines its own
`envelope_command` with a greedy `sed`. `.bench/hooks/block-dangerous-git.sh`
already defines the same function over the same field. That one reads the
string character by character, and it scopes the read under `"tool_input"`.
Two derivations of one fact are a defect, and the new one is wrong. On a real
envelope that carries `cwd` after `tool_input`, it extracts the tail of the
line and reads a plain `ls` as a Bench call.

Move the correct reader into `.bench/lib/resolve-bench.sh`, which both hooks
already source, and delete both copies. Keep the dangerous-git hook's fail
posture: an unreadable envelope is a refusal there, and it stays open for the
follow-on shim. Do the same for `shell_quote` and `rebuild_action`, which the
shim copies from `.bench/hooks/session-start.sh`.

Then replace the pasted library in the canary fixture with the `BASE` and
`MUTATE.json` shape the two guidance fixtures in this spec already use.

## Acceptance

- [ ] One `envelope_command` exists in the tree, in the shared library.
- [ ] The shim classifies `{"tool_name":"Bash","tool_input":{"command":"ls"},"cwd":"/home/u/bench"}` as a non-Bench call.
- [ ] The shim refuses `{"tool_input":{"command":"ls && bench gate"}}` at exit 2 under a stale core.
- [ ] One `shell_quote` and one `rebuild_action` exist in the shell tree.
- [ ] The canary fixture carries no copy of the shared library.
- [ ] Self-probe: restore the greedy reader, and report which row reds.
