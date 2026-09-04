# Name the system-suite route in guidance

Blocked by: none
Writes: AGENTS.md, projects/benchkit.md, internal/anchors/registry_data.go, cmd/bench/testdata/anchors/pre-disclosure-populated.stdout, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/agents-system-suite-route (new), tests/canary/workflow-guidance-anchors/benchkit-system-suite-route (new), tests/canary/load-validity-metadata/shared-rule-drift, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: BF21

## What to build

Verify the premise first: neither AGENTS.md nor the profile names
`bench test --check system` as the hand-run route, and the verb exists in
internal/testreport/command.go. Then add one sentence to the shell conventions
in AGENTS.md and one to the profile's cold-session notes. Each says a hand run
of the system suite goes through `bench test --check system`, which supplies
the sealed run binary and the kit root. Ship both in the five-part precedent
shape: the sentence, one anchor tuple, one red-on-removal registry test, and
one live-mirror fixture each. Keep each needle on one physical line.

## Acceptance

- [ ] The registry test reds a synthetic tree that drops either sentence and stays silent on the live root.
- [ ] Both fixtures bite through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] The prose lane passes on both files.
- [ ] Self-probe: reword one needle, and report the registry test red.
