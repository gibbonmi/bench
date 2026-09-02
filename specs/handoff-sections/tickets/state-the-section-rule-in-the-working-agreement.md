# State the section rule in the working agreement

Blocked by: keep-the-next-command-and-refuse-a-stale-state.md
Writes: AGENTS.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/agents-handoff-section-rule (new), tests/canary/load-validity-metadata/shared-rule-drift, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: HS25

## What to build

Verify the premise first: the phase-close paragraph in AGENTS.md says to
rewrite `capture/session-handoff.md` in full. Then replace that rule. A phase
close runs `bench handoff` from its own worktree, which rewrites only that
assignment's section, and the primary checkout owns `main`. Ship the sentence
in the five-part precedent shape: the sentence, one anchor tuple, one
red-on-removal registry test, and one live-mirror fixture. Keep the needle on
one physical line.

Remove the rewrite-in-full sentence. No anchor tuple pins it
today, so no needle moves.

## Acceptance

- [ ] The registry test reds a synthetic tree that drops the new sentence and stays silent on the live root.
- [ ] The fixture bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] The rewrite-in-full sentence is gone from AGENTS.md.
- [ ] Self-probe: restore the old sentence, and report which check reds it or that none does.
