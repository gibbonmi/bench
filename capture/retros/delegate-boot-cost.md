# Retrospective: delegate-boot-cost spec phase

## Outcome

The phase stages one spec with 30 coverage rows and four tickets in two chunks.
The decision source is the reviewer-confirmed conversation of 2026-09-11.
One Opus/medium review round returned BLOCK with two blocking and ten fold findings.
The author verified every finding against the tree and folded all twelve.
The reviewer has not yet signed off the spec or the ticket graph.
No product implementation ran during this phase.

## Gate-stage timings

The scaffold named the source merge that brought main 244a8acb into the worktree.
That merge ran the whole-project gate before it composed.

- landing: commit 244a8acb725b0e2aaeabcd8f4e06033ce5f09264, trace b4fee53d64658d28c75b4c52ed04fdb0
- gofmt: 139 ms
- vet: 1376 ms
- test: 111947 ms
- race: 3505 ms
- system: 36558 ms
- shellcheck: 708 ms

All six phases passed, and the gate reported six capability skips.

## Ticket-versus-spec-slice and delegate performance

The authoring session retained the spec and the tickets.
One Sonnet docs delegate verified the three harness questions against the subagent and skill pages.
Two Sonnet "ok" delegates measured the general-purpose boot at 38,430 tokens and the Explore boot at 8,834 tokens.
Two Codex exec children measured the Codex boot at 7,258 tokens with and without the docs MCP server.
One Opus/medium review delegate read 19 enforcement files and returned twelve verified findings in one iteration.
Token use of the review delegate was 179,688 by the harness usage line.

| Ticket | Slice | State |
| --- | --- | --- |
| 1-ship-bench-agent-types.md | Agent files, payload, census, guard test, delegate rule | Staged for approval |
| 2-pin-bench-agent-definitions.md | Conformance check and canaries | Staged after ticket 1 |
| 3-trim-skill-descriptions.md | Fourteen skill and six command descriptions | Staged after ticket 1 |
| 4-grade-description-budgets.md | Budget table, check, canaries, lane row | Staged after tickets 2 and 3 |

## Coordinator catches

The coverage check refused three row-id tags in one map, and the author renumbered every row under one tag.
The renumber collided three ids, and the author repaired the collision by line range.
The build preflight refused the spec without a completion-plan fence, and the author added the fence from the record fixture's shape.
The preflight named 76 canary fixtures that pin the touched guidance files, and the tickets now carry each one.
The review round found six command descriptions over the budget that the author's census had skipped.
The review round found that a hard-coded limit would pass every budget row, and DB30 now pins the limit to the table.

## Repair attribution

| ticket | rounds | causes |
| --- | --- | --- |
| 1-ship-bench-agent-types.md | 1 review fold | spec: census row split, README row split, threshold flags |
| 2-pin-bench-agent-definitions.md | 0 | none |
| 3-trim-skill-descriptions.md | 1 review fold | spec: six command files missed by the census |
| 4-grade-description-budgets.md | 1 review fold | spec: single-source row, folded-value row, special-file row |

## Agent-experience improvements

### Bench CLI

The exec wrapper returned a bare exit code when a child search found nothing.
The reviewer could not tell a no-match search from an exec fault.
Feeds: none

The coverage check listed review-owned rows beside uncited rows in one line.
An author must read each seam cell to find the rows that still need a citation.
Feeds: none

The preflight proposal names every pinning fixture from the tree, and the ticket then restates the list.
The parked idea proposes derivation at charge time.
Feeds: none

### Skills

The spec discipline measured one listing surface and skipped the second.
The captured learning proposes a census of every glob-matched subject before a budget row locks.
Feeds: none

### Process

The Codex boot measurement closed one planned ticket before it was written.
A ten-second measurement replaced a config change with no measured gain.
Feeds: none
