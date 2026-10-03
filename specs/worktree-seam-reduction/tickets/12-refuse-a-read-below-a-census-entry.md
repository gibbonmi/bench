# 12. Refuse a read below a census entry

Blocked by: 11-lift-each-read-below-a-census-entry.md
Writes: internal/worktree/single_read_census_test.go (new), internal/worktree/single_read_census_cases_test.go (new), internal/worktree/effect_census_test.go, internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS57, WS58, WS59, WS60, WS61, WS62, WS63, WS64, WS65, WS66, WS67, WS82, WS86

## What to build

Chunk: SR-C6.

Add the single-read census in a new test file. A census entry is a function or a method
declaration whose name is exported. The census derives its read set from three sources:

- each function that `effects.go` declares;
- each package-qualified function that an `effects.go` body calls;
- each exported gate function whose body reaches `KitValue` through calls inside the gate package.

A kind is one package-qualified call that the census reaches from a read-set name. One
reference counts as one read of each kind that its name reaches, so the ambient
constructor and a later `currentTime` read `time.Now` twice. The census accepts a read
only in an entry's own body, outside each function literal and loop body. It accepts only
the first read of each kind there.

The census reports each other read with one of the three messages in the spec's
Implementation decisions. One census formatter renders each message. A census test calls
that formatter with literal file, line, declaration, and name inputs, never with the
census's own findings.

Add the live-tree census test. Ticket 11 lifted the known reads, so the live census
reports nothing. If it reports another read, lift that read in this ticket under the
plan-expansion rule.

Raise `worktreeTestCount` by the number of new top-level tests in this commit. The pin
change adds no line to `parallel_census_test.go`.

## Acceptance

- [ ] Each synthetic census case in the spec's rows WS57 to WS67, WS82, and WS86 draws the formatter's report or no report.
- [ ] The census over the live package reports nothing.
- [ ] The package declares exactly `worktreeTestCount` top-level tests at this commit.
