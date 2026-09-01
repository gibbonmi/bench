# Harden retrospective capture

Blocked by: repair-retro-repeat-preservation.md
Writes: internal/roadmap/retro.go, internal/roadmap/retro_test.go, internal/retros/retros_test.go, internal/retros/testdata/eligible.md (new), internal/retros/testdata/fixture.go (new), cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: none

## What to build

Refuse retrospective capture when any existing destination component is a
symlink. Make public help describe create-only behavior. Give the two test
packages one canonical eligible-retrospective fixture. Record and demonstrate
the mutation that makes the independent fixture expectation bite.

## Acceptance

- [ ] A live or dangling symlink in the retrospective destination is refused
      before an outside-repository write can occur.
- [ ] Public help says capture creates a retrospective and does not promise
      replacement.
- [ ] Retrospective parser and writer tests consume one shared canonical
      eligible body.
- [ ] The ticket records the focused command and red output from a named
      heading-omission mutation.

## Mutation evidence

Mutation: omit `## Coordinator catches` from the canonical eligible fixture.

Focused command:

```text
go test ./internal/retros -run '^TestEligibleFixtureKeepsRequiredHeadings$' -count=1
```

Red output:

```text
--- FAIL: TestEligibleFixtureKeepsRequiredHeadings (0.00s)
    retros_test.go:38: eligible fixture is missing "## Coordinator catches"
FAIL
FAIL github.com/gibbonmi/bench/internal/retros 0.006s
FAIL
```
