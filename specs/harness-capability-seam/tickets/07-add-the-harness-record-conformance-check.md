# Add the harness-record conformance check

Blocked by: 01-add-the-harness-record-package.md
Writes: internal/conformance/harness_record_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/registry.go, tests/canary/harness-record/ (new), projects/benchkit.md

## What to build

A new conformance check, `harness-record`, grades the record against the
tree. It runs at the dev tier over the root.

The check enumerates `.bench/adapters/*` and the shipped hook configs from
disk. Each shipped adapter must map to a row, and each shipped hook config
must map to a row. It then grades each row in both directions, and each
mismatch below is its own diagnostic:

- a declared hook event that the config does not wire
- a wired `.bench/hooks/` script that the row omits
- a row that names an absent headless adapter
- a `delegation_guard` cell that contradicts the `check-agent-line` wiring

The check classifies every path with the no-follow classifier before it
reads, and it names a refused path. A FIFO or a live symlink at a config path
is refused, and that row produces no other diagnostic. A dangling symlink at
an adapter path counts as absent. An absent config and an empty config are
two distinct diagnostics. A hook reference through `$CLAUDE_PROJECT_DIR` and
one through `${CLAUDE_PROJECT_DIR}` both count as wired.

The check joins the enumerated family everywhere the family already appears:

- the `conformanceChecks` map in `internal/conformance/checks_test.go`, bound as `checkHarnessRecord` at `registry.Dev` and `registry.SubjectRoot`.
- the check list in `internal/conformance/registry/registry.go`.
- the conformance table in `projects/benchkit.md`.
- one canary family directory, `tests/canary/harness-record/`.

The check reads no `SKILL.md`, so it does not join `hostileSkillReaders`. The
canary family carries a planted red, and the restored fixture turns the check
green.

This ticket writes the profile's conformance table. Ticket 06 writes the
profile's AXI seam bullet, and ticket 08 writes the same conformance table,
so a coordinator serializes this ticket with ticket 08.

## Acceptance

- [ ] A root with `.bench/adapters/cursor` and no `cursor` row yields a diagnostic naming the adapter. (covers HC29)
- [ ] A root with `.cursor/hooks.json` naming a `.bench/hooks/` script and no `cursor` row yields a diagnostic naming the config. (covers HC30)
- [ ] A `claude` config missing the `Stop` wiring yields a diagnostic naming `Stop`. (covers HC31)
- [ ] A `codex` config that wires `check-agent-line.sh` yields a diagnostic naming the script. (covers HC32)
- [ ] A root without `.bench/adapters/codex` yields a diagnostic naming the row's adapter path. (covers HC33)
- [ ] A `claude` config without `check-agent-line.sh` yields a diagnostic naming `delegation_guard`. (covers HC34)
- [ ] A FIFO at `.codex/hooks.json` yields a diagnostic naming that path and no other diagnostic for that row. (covers HC35)
- [ ] The `harness-record` canary fixture turns the check red, and the restored fixture turns it green. (covers HC36)
- [ ] An absent `.codex/hooks.json` and an empty one yield two distinct diagnostics.
- [ ] The check reports no diagnostic over the live root.
