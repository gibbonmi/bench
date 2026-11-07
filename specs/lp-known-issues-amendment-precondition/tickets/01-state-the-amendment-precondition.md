# State the amendment precondition in the known-issues row

Blocked by: none
Writes: docs/known-issues.md
Covers: none

## What to build

The known-issues row for the stale plan digest after a `bench record assignment` call states when its workaround applies. An assignment rewrites the completion plan, so the plan digest moves. When the review record exists, `bench record chunk` refuses until `bench record amendment` records the new digest. Before the first chunk entry, the review record is absent, and the amendment refuses. The first freeze starts the record at the current plan digest. `bench record ticket` writes the amendment itself.

## Acceptance

- [ ] The row tells an orchestrator to run `bench record amendment` before a `bench record chunk` call only when the review record already holds a chunk entry.
- [ ] The row tells an orchestrator not to run the amendment before the first chunk entry or before a `bench record ticket` call.
