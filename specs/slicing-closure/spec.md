# slicing-closure

Status: staged

Roadmap: FT300

Decision source: reviewer-confirmed conversation, 2026-09-18 (FT300 scope: checks plus prose)

Verification log: 2 iteration(s) to accept — round 1 (fable, medium) rejected the fence-writes seed readers (F1) and raised F2 to F5. Round 2 accepted after one fold.

## Problem

A ticket slicer writes each ticket's `Writes:` line and then the spec's ownership
fence. Build preflight already refuses a `Writes:` entry that omits its fixture
directory or a bound command registry. Two slicing faults still reach the build
without a signal.

First, a ticket edits a guidance file whose sentences the anchor registry pins.
The ticket does not name the registry file that holds the needle, or the test file
that holds its independent expectation. The author finds those files mid-build and
pays a repair round.

Second, the fence and the union of the ticket `Writes:` values drift apart. A fence
entry that no ticket owns grants write authority that nobody planned. A backticked
prose token in the fence section parses as such an entry. A `Writes:` entry that the
fence omits fails at the first commit that touches it.

The slicer also applies several judgment rules that no guidance states. The retros
record repair rounds for each missing rule.

## Solution

Build preflight and review preflight grade two more rows over the parsed tickets.
The `anchor-closure` row reds a ticket that writes an anchored guidance file and
omits a registry file that names it. The `fence-writes` row reds a spec whose fence
tokens differ from the union of its ticket `Writes:` paths. Each red names the
exact entries to add or remove. The `--propose-writes` form lists the missing anchor
registry files beside the fixture and registry closures.

A new `craft-tickets` reference states the preflight-enforced `Writes:` rules and six
judgment rules for the slicer. Anchor rows pin each rule, so a deletion reds the gate.

## User stories

Line: opus / medium for SC-C1 and SC-C2; opus / high for SC-C3.
Implementation-line reason: SC-C1 is the hardest code chunk, because it adds a tree scan, a closure kind, and a proposal source. It also makes headroom in files over budget. The decisions are precise, the closure seam has two precedents, and `Decide` table tests make each row red-capable. SC-C3 is guidance authoring, which the profile routes to high effort.
Harder chunks: SC-C1.

### Anchor closure

1. As a slicer, I want preflight to red an anchored guidance write that omits its registry files, so that I avoid a repair round.
2. As a slicer, I want the red to name the ticket, the entry, and each missing file, so that one pass repairs it.
3. As a slicer, I want a directory entry over an anchored file to take the same closure, so that no file hides.
4. As a slicer, I want every anchors Go file whose string literal names the path in the closure, so that test expectations join too.
5. As a slicer, I want a ticket that names each registry file, or a covering directory, to stay green, so that correct work passes.
6. As a linked-repository user, I want no anchor closure where no anchor registry directory exists, so that my guidance edits pass.
7. As a slicer, I want `--propose-writes` to list each missing anchor file with its source and fence state, so that the proposal repairs it.
8. As a maintainer, I want the scan to refuse a link or special file before a read, so that preflight stays in the tree.
9. As a spec author, I want the closure derived from the graded tree, so that a branch that adds an anchor is graded correctly.
10. As a coordinator, I want the new rows not-applicable in a build without tickets, so that a fresh build reads the same.

### Fence and `Writes:` union

11. As a spec author, I want preflight to red when the fence and the ticket `Writes:` union differ, so that the fence cannot drift.
12. As a spec author, I want the red to name each side's extra entries apart, so that I know which side to repair.
13. As a spec author, I want the review pickup excluded from both sides, so that the required pickup needs no ticket owner.
14. As a slicer, I want entries compared without the `(new)` marker or a trailing slash, so that spelling does not red.
15. As a slicer, I want `Writes:` paths under the spec folder or capture excluded, so that implicit authority stays implicit.
16. As a reviewer, I want a backticked prose token that no ticket writes reported as fence-only, so that it grants no hidden authority.
17. As a slicer, I want `--propose-writes` to run while `fence-writes` is red, so that it can report the fence expansion.
18. As a test author, I want the conformant seeds to stay union-exact, so that every seeded green tree stays green.

### Slicing guidance

19. As a slicer, I want one reference to list every enforced `Writes:` rule, so that I meet each rule before preflight reds.
20. As a slicer, I want a rule that a lane-check ticket proves its check through the real lane, so that no stand-in passes.
21. As a slicer, I want each call site of a fixture helper that a posture change reds listed before map lock, so that none escapes.
22. As a slicer, I want a combined behavior row owned by the ticket that completes its final consumer, so that its owner can prove it.
23. As a slicer, I want each sentence that grants retired behavior to get a forbid row and a red-capable check, so that no grant survives.
24. As a slicer, I want a cited verifier row to name its exact checks, so that a reviewer can repeat them.
25. As a slicer, I want a rule to rerun build preflight after each fence change, so that review grades the final fence.
26. As a maintainer, I want an anchor row and an independent expectation for each new rule, so that a deleted rule reds the gate.
27. As a maintainer, I want the skill to point to the reference within its prose budget, so that the skill stays lean.

### Reviewed exclusions

28. As the FT98 owner, I want the preserve-and-move ref-shape rule left to the preserve-then-discard primitive, so that one owner holds that behavior.
29. As the FT293 owner, I want preflight to keep refusing an omitted closure rather than deriving it, so that the open derivation decision stays open.

## Implementation decisions

- **Anchor reference scan.** The anchors package owns one function that reads the anchor registry directory of a repository root. For each Go string literal there, it returns the sorted repository-relative files that hold it. The scan reads the top-level `.go` files, test files included, and it does not recurse. It tokenizes each file and unquotes every interpreted and raw string literal.
- **Scan refusals.** The scan classifies each entry without following links. A symbolic link, a special file, or an unreadable file refuses the whole scan and names the path. A tokenizer error refuses the scan and names the file. An absent directory is an empty result, not an error.
- **Anchor closure.** A `Writes:` entry requires the referencing files of each literal equal to its path. It also requires the files of each literal under its path at a `/` boundary. The closure grades them against the ticket's owned paths with the existing coverage rule. The `anchor-closure` row renders after `registry-closure`. Its red detail opens with `Writes: entry names an anchored guidance path without naming every anchor registry file: ` and lists `<ticket>: <entry> is anchored by <file>` items.
- **Closure family.** The three closure kinds share one requirement derivation and one message builder. The proposal lists every closure kind. The proposal source for an anchor requirement reads `anchor <entry>`.
- **Proposal-tolerated checks.** One list names the checks whose red does not refuse `--propose-writes`. SC-C1 creates it with the three closure checks. The proposal reads only that list.
- **Fence and `Writes:` union.** The `fence-writes` row renders after `anchor-closure`. It compares the fence tokens with every owned path across the parsed tickets. Each side drops a trailing `/` and the `(new)` marker. Each side removes the review pickup. Each side also removes every path that the implicit entries of `paths-authorized` cover: the spec folder and the capture folder.
- **Union red.** The red detail reads `spec fence and ticket Writes: union differ: fence only: <a, b>; Writes only: <c>`. Each list is sorted, and an empty side is omitted. The review pickup comes from the review record path owner, so preflight does not derive the pickup a second time.
- **Proposal tolerance.** SC-C2 adds `fence-writes` to the proposal-tolerated checks list, because the proposal's fence column already reports the expansion that it requires.
- **Headroom.** `gather.go` and `decision.go` are over their line budgets. SC-C1 moves the per-entry `Writes:` probe and the owned-path helpers into the closure file before it adds lines. No file grows past its budget at any ticket commit.
- **Seeds.** The fixture has two fence tables: the conformant fence and the review fence extras. The seeded ticket's `Writes:` values derive from the fence table that its seed declares. Each derived entry carries the `(new)` marker, because the ticket builder has no tree access. A test that adds fence lines passes the same entries to a writes-aware ticket builder.
- **Seed readers.** Each seed reader that fences an entry no ticket writes moves to the writes-aware builder: the 5000-entry budget seed and the closed-parenthesis bootstrap test. The proposal tests and the anchor closure seed stop rewriting the literal `Writes: specs` text and pass their writes to that builder. The build metadata expectation in `charge_evidence_test.go` takes the derived `Writes:` entries.
- **Shared seeds.** By reviewer decision, three shared seeds outside the preflight package also become union-exact. The review record fixture closes the fence section with a plan heading. It owns the one helper that adds a fence entry above that heading. The landing race ticket writes its fenced files. The version test ticket writes its `outside/` entry.
- **Row readers.** Each new row changes the row count and row list that five test files assert. The legacy charge baseline, the review command's row count, the bare row count in `source_tip_test.go`, and the two `Decide` row-order lists take the new row. `baseFacts` becomes union-exact. The edits in `command_review_test.go` (480 lines), `decision_test.go` (821 lines), and `command_bootstrap_test.go` (395 lines) are net-neutral in line count, so none of them grows.
- **Guidance.** A new reference file under the `craft-tickets` skill holds the enforced `Writes:` rules and the six judgment rules. The skill replaces its parser-rules paragraph with a pointer sentence to that reference, so the skill does not grow. The ticket-slicing anchor group pins each new sentence, and its test file states each needle independently.

## Implementation chunks

One retained implementation session owns the build after approval. Each ticket is one review chunk and one serial green commit checkpoint. After each chunk, freeze its predecessor and current tips for Standards, Spec, and Coverage review. The successor starts after accepted findings have current repair coverage.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| SC-C1 / `1-close-anchor-registry-writes.md` | Preflight grades and proposes anchor closure | SC1, SC2, SC3, SC4, SC5, SC6, SC7, SC8, SC9, SC10, SC11, SC12, SC13 | the anchors and preflight tests the rows name | yes |
| SC-C2 / `2-match-fence-to-writes.md` | Preflight grades the fence against the `Writes:` union | SC14, SC15, SC16, SC17, SC18, SC19, SC20, SC21, SC22, SC23 | the preflight tests the rows name | no |
| SC-C3 / `3-state-slicing-checks.md`, `4-repair-slicing-checks-review.md` | The slicer reads every enforced `Writes:` rule and six judgment rules | SC24, SC25, SC26, SC27, SC28, SC29, SC30, SC31 | `docs-currency-workflow`, `TestTicketSlicingPasses`, and the prose budget check | no |

## Testing decisions

- A good test drives `Decide` over constructed facts, or drives `bench preflight build` over a seeded repository, and reads the rendered row. It never reads an internal field.
- The closure rows follow the fixture-closure and registry-closure tests. The command test follows `TestCommandBuildRendersSixGrammarRows`.
- The anchor scan tests build a temporary directory and call the anchors function directly.
- The guidance rows attach to `docs-currency-workflow`, which grades the live sentences. They also attach to the ticket-slicing anchor test and the prose budget check.

### Seam diagram

    trigger: bench preflight build|review <slug>, and --propose-writes
        │
        ▼
    tree (internal/anchors/*.go, spec fence, tickets)  ──▶  [ gather → Facts → Decide ]  ──▶  check rows
                      ◀ tests attach here: Decide table tests over Facts; one command test over a seeded repo

    trigger: bench gate (anchor harness and prose budget checks)
        │
        ▼
    craft-tickets SKILL.md + references/slicing-checks.md  ──▶  [ anchor registry evaluation ]  ──▶  diagnostics
                      ◀ tests attach here: docs-currency-workflow grades the live tree; TestTicketSlicingPasses states each needle independently

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SC1 | 1 | `Decide` reds `anchor-closure` when a ticket writes `.agents/x/SKILL.md`, the anchor map names `internal/anchors/registry_a.go` for it, and the ticket omits that file | planned TestAnchorClosureRedsUnnamedRegistry in internal/preflight | A row that never reads the anchor map stays green on this input |
| SC2 | 2 | The red detail lists `<ticket>: <entry> is anchored by <file>` once for each missing file | planned TestAnchorClosureNamesEachMissingFile in internal/preflight | A detail that names only the ticket or only the first file fails the exact-string assertion |
| SC3 | 3 | A `Writes:` entry `.agents/x` takes the closure of the anchored file `.agents/x/SKILL.md` | planned TestAnchorClosureCoversDirectoryEntry in internal/preflight | An exact-key lookup, as the fixture pins use, finds no requirement for a directory entry |
| SC4 | 4 | The scan returns a `_test.go` file whose string literal equals the guidance path | planned TestReferencingFilesIncludeTestFiles in internal/anchors | A scan that skips test files drops the independent expectation file |
| SC5 | 4 | The scan returns a file whose literal is a positional value or a constant, not a `File:` key | planned TestReferencingFilesReadEveryStringLiteral in internal/anchors | A scan keyed on the `File:` field misses the constant and positional spellings |
| SC6 | 4 | The scan unquotes a raw string literal and matches its value | planned TestReferencingFilesUnquoteRawLiterals in internal/anchors | A scan that compares quoted source text misses every backquoted literal |
| SC7 | 5 | `anchor-closure` is green when the ticket names each registry file, and green when it names the `internal/anchors` directory | planned TestAnchorClosureGreenWhenNamed in internal/preflight | A closure that ignores directory coverage reds a correct ticket |
| SC8 | 6 | The scan of a root with no `internal/anchors` directory returns an empty result and no error | planned TestReferencingFilesAbsentDirectory in internal/anchors | A scan that treats absence as an error refuses every linked repository |
| SC9 | 7 | `--propose-writes` lists a missing anchor registry file with source `anchor <entry>` and its fence state, and an `anchor-closure` red does not refuse the proposal | planned TestProposeWritesListsAnchorClosure in internal/preflight | A proposal with a hard-coded closure list refuses on the new red or omits its row |
| SC10 | 8 | The scan refuses a symbolic link named `link.go` in the directory and names that path | planned TestReferencingFilesRefuseSymlink in internal/anchors | A scan that follows the link reads bytes outside the tree and returns them |
| SC11 | 8 | The scan refuses a FIFO named `pipe.go` in the directory without a blocking read | planned TestReferencingFilesRefuseFIFO in internal/anchors | A scan that opens the entry before it classifies it blocks on the FIFO |
| SC12 | 9 | `bench preflight build` reds `anchor-closure` for a seeded repository whose own `internal/anchors/extra.go` names the ticket's guidance path | planned TestCommandBuildAnchorClosureReadsTree in internal/preflight | A closure compiled into the binary does not know the seeded file and stays green |
| SC13 | 10 | In build mode with no tickets directory, `anchor-closure` renders not-applicable directly after `registry-closure` | `internal/preflight/decision_test.go` (`TestTicketGatedRowsNotApplicableWithoutTickets`) | A row outside the ticket-row gate renders a verdict for a fresh build |
| SC14 | 11 | `Decide` reds `fence-writes` when the fence holds `a/` and the tickets write `b.go` | planned TestFenceWritesRedsDifferentSets in internal/preflight | A check that tests only one containment direction passes one of the two mismatches |
| SC15 | 12 | The red detail reads `spec fence and ticket Writes: union differ: fence only: a; Writes only: b.go` | planned TestFenceWritesNamesBothSides in internal/preflight | A detail that merges the two sides fails the exact-string assertion |
| SC16 | 11 | `fence-writes` is green when the two sets are equal | planned TestFenceWritesGreenWhenExact in internal/preflight | A check that reds on any fence entry fails the green case |
| SC17 | 13 | The review pickup on either side alone leaves `fence-writes` green | planned TestFenceWritesIgnoresReviewPickup in internal/preflight | A check that keeps the pickup reds every spec, because the coverage grammar requires it in the fence |
| SC18 | 14 | A fence entry `internal/x/` and a `Writes:` entry `internal/x (new)` compare equal | planned TestFenceWritesNormalizesSpelling in internal/preflight | A check that compares raw tokens reds a correct slice |
| SC19 | 15 | A `Writes:` path under the spec folder or under `capture` does not appear as `Writes:`-only | planned TestFenceWritesIgnoresImplicitAuthority in internal/preflight | A check without the implicit entries reds a ticket that amends its own spec |
| SC20 | 16 | `bench preflight build` reports a backticked prose token in the fence section as fence-only | planned TestCommandBuildFenceWritesReportsProseToken in internal/preflight | A check that reads a different fence parser than preflight's facts misses the token that the parser returns |
| SC21 | 17 | A `fence-writes` red does not refuse `--propose-writes` | planned TestProposeWritesToleratesFenceWrites in internal/preflight | A proposal that refuses on every non-closure red fails this case |
| SC22 | 18 | `bench preflight build` over the conformant seed renders `fence-writes,green` | `internal/preflight/command_build_test.go` (`TestCommandBuildResumedTicketsRunForReal`) | A seed whose `Writes:` still names `specs` renders the row red |
| SC23 | 10 | In build mode with no tickets directory, `fence-writes` renders not-applicable directly after `anchor-closure` | `internal/preflight/decision_test.go` (`TestTicketGatedRowsNotApplicableWithoutTickets`) | A row outside the ticket-row gate renders a verdict for a fresh build |
| SC24 | 19, 26 | The reference states the anchor closure rule and the fence union rule as pinned sentences | the `docs-currency-workflow` check and `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingPasses`) | Deleting either sentence reds `docs-currency-workflow`, and a drifted registry row reds `TestTicketSlicingPasses` |
| SC25 | 20, 26 | The reference states the real-lane proof rule as a pinned sentence | the `docs-currency-workflow` check and `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingPasses`) | Deleting the sentence reds `docs-currency-workflow`, and a drifted registry row reds `TestTicketSlicingPasses` |
| SC26 | 21, 26 | The reference states the fixture-helper call-site rule as a pinned sentence | the `docs-currency-workflow` check and `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingPasses`) | Deleting the sentence reds `docs-currency-workflow`, and a drifted registry row reds `TestTicketSlicingPasses` |
| SC27 | 22, 26 | The reference states the final-consumer row rule as a pinned sentence | the `docs-currency-workflow` check and `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingPasses`) | Deleting the sentence reds `docs-currency-workflow`, and a drifted registry row reds `TestTicketSlicingPasses` |
| SC28 | 23, 26 | The reference states the retirement forbid-row rule as a pinned sentence | the `docs-currency-workflow` check and `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingPasses`) | Deleting the sentence reds `docs-currency-workflow`, and a drifted registry row reds `TestTicketSlicingPasses` |
| SC29 | 24, 26 | The reference states the exact-verifier rule as a pinned sentence | the `docs-currency-workflow` check and `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingPasses`) | Deleting the sentence reds `docs-currency-workflow`, and a drifted registry row reds `TestTicketSlicingPasses` |
| SC30 | 25, 26 | The reference states the rerun-preflight rule as a pinned sentence | the `docs-currency-workflow` check and `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingPasses`) | Deleting the sentence reds `docs-currency-workflow`, and a drifted registry row reds `TestTicketSlicingPasses` |
| SC31 | 27 | The skill's pinned pointer sentence names the reference, and the skill stays within 100 lines | the `docs-currency-workflow` check and `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingPasses`) and the `guidance-prose-budgets` check | Deleting the pointer reds `docs-currency-workflow`, a drifted registry row reds `TestTicketSlicingPasses`, and a skill over 100 lines reds the budget check |

Not covered: story 28 — the exclusion moves the clause to FT98 and adds no behavior here.
Not covered: story 29 — the exclusion keeps current refusal behavior, which the existing closure tests already grade.

### Edge inventory

- Paths with spaces: a `Writes:` entry or a literal may hold a space. Both sides compare exact strings, so SC1 and SC18 fixtures use one spaced path each.
- Absent versus empty anchor directory: SC8 grades absence. An empty directory also returns an empty result, in the same test.
- Special files and links in a discovered path: SC10 and SC11.
- A grammar token quoted in prose that parses as live syntax: SC20.
- Trailing slash and the `(new)` marker: SC18.
- **Won't handle:** a concatenated path in an anchors file — each guidance path there is one literal today. The anchor harness still reds a moved needle.
- **Won't handle:** anchored prose pinned outside the anchors package — FT300 names the anchor registry only, and fixture closure covers canary pins.
- **Won't handle:** needle-level closure — preflight cannot know which sentences a ticket will edit, so the rule is file-level.
- **Won't handle:** a deleted-file `Writes:` marker, derived closure files, and a `(new)` test beside an over-budget sibling — FT293 owns those decisions.
- **Won't handle:** a ticket that names an anchors file must also name the five command registries — registry closure binds the anchors package today. Reviewer veto under FT293.
- **Won't handle:** amending the fences of other staged specs — `fence-writes` reports their drift when their build runs preflight.

## Ownership fences

Reviewer disposition: proposed within the confirmed scope. It grants no implementation authority until the implementation phase starts.

- `internal/anchors/references.go`
- `internal/anchors/references_test.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/preflight/closure.go`
- `internal/preflight/decision.go`
- `internal/preflight/gather.go`
- `internal/preflight/proposal.go`
- `internal/preflight/anchor_closure_test.go`
- `internal/preflight/fence_writes.go`
- `internal/preflight/fence_writes_test.go`
- `internal/preflight/decision_test.go`
- `internal/preflight/command_build_test.go`
- `internal/preflight/proposal_test.go`
- `internal/preflight/charge_test.go`
- `internal/preflight/command_review_test.go`
- `internal/preflight/source_tip_test.go`
- `internal/preflight/charge_evidence_test.go`
- `internal/reviewrecord/recordtest/fixture.go`
- `internal/systemtest/owner_landing_fixture_test.go`
- `cmd/bench/preflight_version_test.go`
- `internal/worktree/land_fixtures_test.go`
- `internal/worktree/land_effects_test.go`
- `internal/worktree/land_journey_test.go`
- `internal/preflight/command_bootstrap_test.go`
- `internal/preflight/evidencecmd/evidence_budget_test.go`
- `internal/preflight/preflighttest/fixture.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `.agents/skills/bench-craft-tickets/references/slicing-checks.md`
- `tests/canary/workflow-guidance-anchors/craft-tickets-slice-acceptance-row`
- `tests/canary/workflow-guidance-anchors/delegated-serial-ticket-checkpoint`
- `tests/canary/workflow-guidance-anchors/dg-17`
- `tests/canary/workflow-guidance-anchors/dg-18`
- `tests/canary/workflow-guidance-anchors/dg-19`
- `tests/canary/workflow-guidance-anchors/dg-20`
- `tests/canary/workflow-guidance-anchors/dg-21`
- `tests/canary/workflow-guidance-anchors/dg-22`
- `tests/canary/workflow-guidance-anchors/dg-23`
- `tests/canary/workflow-guidance-anchors/ticket-executable-route-pass`
- `tests/canary/workflow-guidance-anchors/ticket-light-path-carve-out`
- `tests/canary/workflow-guidance-anchors/ticket-lock-passes`
- `tests/canary/workflow-guidance-anchors/ticket-relocation-destinations`
- `tests/canary/workflow-guidance-anchors/ticket-seam-creating-slice`
- `tests/canary/workflow-guidance-anchors/ticket-skill-contract-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-source-clause-pass`
- `tests/canary/workflow-guidance-anchors/ticket-template-anchor`
- `reviews/slicing-closure.md`

## Out of scope

- The preserve-and-move ref-shape rule belongs to FT98's preserve-then-discard primitive. Estimate: 4 edits, 1 gate run.
- A single owner for the review pickup path: the coverage package and the review record package each render it today. Estimate: 3 edits, 1 gate run.
- A single derivation of the row-ID tag: `tickets.TagOf` scans the tag by hand beside `tickets.RowIDPattern`. The live idea inbox holds it. Estimate: 2 edits, 1 gate run.

## Further notes

- Two staged specs parse `.bench/BENCH.md` as a fence token from a prose sentence in the fence section, and no ticket writes it: `specs/session-context-overflow/spec.md:202` and `specs/session-context-queries/spec.md:229`. `fence-writes` reds both when their builds run preflight.
- SC-C3 dogfoods SC-C1: its `Writes:` line names every anchor registry file that names the skill.
- The SC-C3 mutation runs through `docs-currency-workflow`, the check that grades the live anchor group. The anchors package grades only temporary trees, so a deleted live sentence leaves it green.

## Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"SC-C1","tickets":["1-close-anchor-registry-writes.md"],"verification":[{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"preflight","command":"bench test --package ./internal/preflight/..."},{"id":"mutation","command":"bench test --package ./internal/preflight/...","probe":"return no anchor requirement for a directory Writes entry"}]},{"id":"SC-C2","tickets":["2-match-fence-to-writes.md"],"verification":[{"id":"preflight","command":"bench test --package ./internal/preflight/..."},{"id":"mutation","command":"bench test --package ./internal/preflight/...","probe":"keep the review pickup in the fence set"}]},{"id":"SC-C3","tickets":["3-state-slicing-checks.md","4-repair-slicing-checks-review.md"],"verification":[{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"budget","command":"bench test --check guidance-prose-budgets"},{"id":"prose","command":"bench test --check prose"},{"id":"mutation","command":"bench test --check docs-currency-workflow","probe":"delete the rerun-preflight sentence from the reference"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/slicing-closure/spec.md"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"preflight","command":"bench test --package ./internal/preflight/..."},{"id":"budget","command":"bench test --check guidance-prose-budgets"},{"id":"canary","command":"bench test --check canary-fixture-compliance"}]}
```
