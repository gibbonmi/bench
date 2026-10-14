# Add the spec review known issues

Blocked by: none
Writes: docs/known-issues.md
Covers: none

## What to build

The FT405 spec review on 2026-10-14 found two gaps that no check catches yet.
The reviewer parked each gap as an idea in `capture/IDEAS.md`, and the reviewer directed one known-issue row for each idea.

Add two rows to the known-issue table in `docs/known-issues.md`:

1. A spec fork can add a seam beyond its decision source with no research. In spec review iteration 1, three of five blockers sat on such seams. The workaround: the spec author cites a `path:line` and a research finding for each added seam before the spec review. The retire condition is `no owner`, because the seam-delta idea is uncommitted intake.
2. A spec fork can drop an anchor registry file, or its `_test.go` twin, that pins a needle in a guidance file. The FT405 spec fork dropped two charged anchor registries. The workaround: run `rg` for each edited guidance path across `internal/anchors` and `internal/conformance`, and name each hit, test twins included. The retire condition is `no owner`, because the anchor-closure idea is uncommitted intake.

Change no other row.

## Acceptance

- [ ] The table holds one row for unresearched seams that a spec fork adds, with its workaround and the retire cell `no owner`.
- [ ] The table holds one row for anchor registry files that a spec fork drops, with its workaround and the retire cell `no owner`.
- [ ] Every other row of `docs/known-issues.md` is unchanged.
