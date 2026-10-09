# Enforce the ticket check floor and the planned-cell rewrite

Status: staged

Roadmap: FT349

Decision source: `roadmap/FT349.md`, the named reviewed artifact from drain `d-bab100a48f42`, with the commitment criteria `FT349.floor` and `FT349.planned` in `.bench/commitment.json`.

Verification log: 2 iteration(s) to accept — paired Sonnet and Haiku reviewers reviewed the spec, then the tickets, with one confirming pass each. The spec fold closed six blocking items: the build-entry red route, four degenerate fixtures, and the exact public-response predicate. The ticket fold closed one blocking item: ticket 01 writes `internal/testreport/command.go` for its refusal before any Go child.

## Problem

A ticket author must run three kinds of focused checks before the ticket commit. The first is root conformance. The second is the tests of each package that the ticket writes. The third is the whole `cmd/bench` package for a public response or embed change.

Only guidance prose states this floor (`.agents/commands/bench-implement-spec.md:44`). No command grades it.
In the roadmap-delivery-commitment build, nine plan commits added package checks that the planned chunk commands omitted. In FT337, a chunk checkpoint went red after a clean review because the ticket checks omitted `cmd/bench`.

A coverage seam cell can hold the marker `planned` for a test that does not exist yet (`internal/coverage/citations.go:147-150`).
No check refuses the marker after the test exists. At FT393 final reconciliation, 55 rows still held `planned` cells for tests that existed.

Root conformance also has no exact command. `bench test --check <name>` runs one conformance scope (`internal/testreport/named_check.go:91-99`). Only the gate runs the unscoped entry test (`internal/gate/phases.go:138-149`).

## Solution

The spec's completion plan already holds the planned verification commands of each chunk. The chunk checkpoint already refuses a planned command that has no passing record (`internal/reviewrecord/check.go:27-56`).
This spec makes the build preflight compare those planned commands with a floor that it derives from each ticket's `Writes:` line. A red `check-floor` row stops the next build charge and the review charge.

`bench record chunk` freezes a chunk after its last ticket commits. This spec makes that freeze refuse a row of the chunk whose seam cell still holds `planned`. `bench record completion` refuses the same cell on any plan row.

`bench test --check root-conformance` becomes the one exact command for the unscoped root conformance pass. The commit lane stays unchanged.

## User stories

Line: opus / high.
Implementation-line reason: The hardest material chunk is TF-C3, a new preflight refusal whose wrong-green result the gate cannot catch. The spec is precise and the seams are known. TF-C5 edits guidance, so the leverage override applies.
Harder chunks: TF-C2, TF-C3.

### Root conformance has one command

1. As an author, I want `bench test --check root-conformance` to run the unscoped entry test, so that one command grades every dev-tier check.
2. As an author, I want a red root-conformance run to print each failing diagnostic, so that the red names its defect.
3. As a linked-repository author, I want the check refused where the gate skips root conformance, so that no run claims a false pass.
4. As an operator, I want `bench test --checks` to list `root-conformance` as a `conformance` check, so that discovery finds the command.
5. As a maintainer, I want the gate phase and the named check to read one applicability rule, so that the two never disagree.

### The floor derives from the Writes line

6. As the preflight, I want each floor to hold root conformance where it applies, so that every plan grades the live tree.
7. As the preflight, I want the floor to hold the package of each `Writes:` Go file, so that the plan tests the edited code.
8. As the preflight, I want a `(new)` Go file in a current package to add that package, so that the package keeps its tests.
9. As the preflight, I want the floor to hold each current package under a `Writes:` directory, so that a directory entry hides no package.
10. As the preflight, I want a `Writes:` embed file to add its owner package, so that a changed asset runs its owner's tests.
11. As the preflight, I want the floor to hold every package for a `go.mod` or `go.sum` entry, so that a module change runs the suite.
12. As the preflight, I want a path that no package embeds to add no package, so that a guidance ticket owes root conformance only.
13. As the preflight, I want the floor to hold no package for a `(new)` entry outside every package, so that plan-only time invents no package.
14. As the preflight, I want a bound command package entry to add `./cmd/bench`, so that a response change runs the command tests.
15. As the preflight, I want an embed package entry to add `./cmd/bench`, so that an embed change runs the command tests.
16. As the preflight, I want the floor to omit `./cmd/bench` when neither condition holds, so that the floor stays narrow.
17. As a linked-repository author, I want an empty package floor in a root without `go.mod`, so that a non-Go repository owes no Go test.
18. As a maintainer, I want the floor to read the `bench test --changed` package index, so that path-to-package knowledge has one source.
19. As an operator, I want `bench test --changed` selection unchanged after the index is shared, so that no focused run changes.

### The preflight grades the plan

20. As an orchestrator, I want `--plan-only` to red `check-floor` when a chunk omits a floor element, so that I fix the plan before dispatch.
21. As an orchestrator, I want the red to name the ticket, the element, and the command to add, so that the repair is mechanical.
22. As an orchestrator, I want `bench test --package ./internal/x` to satisfy `./internal/x`, so that the usual spelling counts.
23. As an orchestrator, I want `./internal/...` to satisfy each package below it, so that one broad command counts.
24. As an orchestrator, I want `./...` to satisfy every package element, so that a whole-suite command counts.
25. As an orchestrator, I want `./internal/git/...` to leave `./internal/git2` unsatisfied, so that a match obeys segment boundaries.
26. As an orchestrator, I want a `--run` command to leave its package unsatisfied, so that one filtered test cannot stand in for a package.
27. As an orchestrator, I want a `--changed` command to satisfy nothing, so that a base-dependent selection cannot stand in for the floor.
28. As an orchestrator, I want a compound or foreign command to satisfy nothing, so that a shell line cannot hide an omission.
29. As an orchestrator, I want a command that the `bench test` grammar refuses to satisfy nothing, so that a malformed command never counts.
30. As an orchestrator, I want only the root-conformance check to satisfy its element, so that one scope cannot stand in for the pass.
31. As an orchestrator, I want each chunk graded against its own requirements, so that no chunk borrows another chunk's checks.
32. As an orchestrator, I want each delegated ticket graded against its own requirements, so that no author borrows another author's checks.
33. As an orchestrator, I want the after-ticket preflight to run `check-floor`, so that a new package reds until the plan names it.
34. As a reviewer, I want the review preflight to run `check-floor`, so that the review charge refuses an incomplete plan.
35. As an orchestrator, I want a red `check-floor` row to refuse the build charge, so that no author starts on an incomplete plan.
36. As an orchestrator, I want `check-floor` not-applicable under an invalid completion plan, so that one defect gives one red.
37. As an orchestrator, I want `check-floor` not-applicable in build mode without tickets, so that a fresh spec stage is not graded.
38. As an orchestrator, I want a package-index failure to red `check-floor`, so that an unreadable package graph never passes.
39. As the reviewer, I want the commit lane unchanged, so that the ADR 0017 fast lane keeps its check list.

### The chunk freeze refuses a planned cell

40. As an orchestrator, I want `bench record chunk` to refuse a chunk row with a `planned` seam cell, so that a frozen chunk cites real tests.
41. As an orchestrator, I want the refusal to name each such row and the citation form, so that the repair is mechanical.
42. As an orchestrator, I want the refusal to write no record, so that no partial chunk entry remains.
43. As an orchestrator, I want a chunk row with a cited seam cell to pass, so that a rewritten cell is accepted.
44. As an orchestrator, I want a chunk row with a `review-owned:` seam cell to pass, so that a review-owned row needs no test.
45. As an orchestrator, I want a `planned` row outside the chunk to pass, so that a later chunk's future test stays valid.
46. As an orchestrator, I want a cell with `unplanned` to pass, so that the marker keeps its word boundary.
47. As an orchestrator, I want `bench record completion` to refuse a `planned` cell on any plan row, so that the landing cites real tests.
48. As a spec author, I want `bench coverage --check` to keep passing a `planned` cell, so that a staged spec can plan future tests.
49. As a maintainer, I want the planned-row projection to read the one coverage marker owner, so that the marker has one spelling.

### Guidance states the floor

50. As an author, I want the implement-spec guidance to state that the plan carries the floor, so that I plan checks before I run them.
51. As an author, I want the implement-spec guidance to state the planned-cell rewrite, so that the chunk freeze accepts my rows.
52. As a spec author, I want map discipline to state when a `planned` marker must go, so that I write the marker correctly.
53. As a delegate, I want delegation discipline to name the root-conformance command, so that a charge names one exact command.
54. As the reviewer, I want the FT89 question about the per-return gate run left open, so that this spec decides nothing there.
55. As an orchestrator, I want the build-entry route to send a `check-floor` red to the staleness pass, so that I amend the plan there.

## Implementation decisions

### Root conformance command

The gate exports one predicate that answers whether root conformance applies to a graded root. The predicate wraps the rule that `rootConformanceEnv` applies today: the root is the kit, and the kit declares the entry test.
The gate test phase and the new named check both read that predicate.

`bench test --check root-conformance` runs the conformance entry test with the dev tier and an empty scope. An empty scope runs every registered dev-tier check (`internal/conformance/checks_test.go:97-120`).
Today the named-check environment sets the scope to the check name (`internal/testreport/environment.go:24`). The entry test refuses a scope that names no registered check (`internal/conformance/checks_test.go:105-109`). So the root-conformance run passes an empty scope, not its own name.

The name joins the named-check inventory with the kind `conformance`. Where the predicate answers false, the check refuses before it starts a Go child. The refusal sits in `runFocusedRequest` beside the system check refusal, before the run-binary selection (`internal/testreport/command.go:134-144`).

The kit that the predicate reads is the answer of the existing export `gate.KitRoot` (`internal/gate/tag_census.go:231`). That answer is `BENCH_KIT` when that variable is set, and the graded root otherwise. Today `testBenchSource` (`internal/testreport/command.go:329-334`) keeps a private copy of that rule. TF-C1 makes it call `gate.KitRoot`, and the build preflight reads the predicate through the same resolver.

### Floor derivation

`internal/testreport` owns the floor's package derivation, because it already maps changed paths to packages for `bench test --changed` (`internal/testreport/selection.go:172-263`).
The selection and the floor share one package index built from `go list -json -test ./...`. The index maps a directory to its package, an embed file to its owner, and a package to its embed patterns.

The floor function takes the split `Writes:` paths of one ticket. It returns the directory-relative patterns of the packages that the ticket writes directly, and whether any entry touches an embed.
It applies these rules:

- A Go file entry, present or absent, adds the package of its directory when that directory holds a current package.
- A directory entry adds each current package whose directory the entry covers at a segment boundary.
- An embed file entry adds its owner package and sets the embed flag.
- An entry that selects a package with a non-empty embed pattern list sets the embed flag.
- A `go.mod` or `go.sum` entry adds the pattern `./...`.
- Any other entry adds nothing.
- A root without `go.mod` returns an empty floor and starts no `go list`.

The repository root is a package, because the module root holds `consumer_payload.go`. Its element is `./.`, which is the target that `packagePattern` gives the operand `.` (`internal/testreport/command.go:314-327`).

An absent `(new)` embed file adds nothing. The package index lists only embed files that exist, and the selection skips an absent non-Go path (`internal/testreport/selection.go:201-205`). The after-ticket preflight adds the owner once the file exists.

The floor splits each entry with `tickets.WritesPath` and matches directories with `tickets.Covers` (`internal/tickets/writes.go:13-25`). `go list -deps` shows that `internal/testreport` already depends on `internal/tickets`. The dependency closure of `internal/tickets` holds seven Bench packages, and none of them reaches `internal/testreport`. So the edge adds no cycle.

The reverse-import closure of `--changed` is not part of the floor. The floor names the packages that the ticket writes, as the guidance does.

`testreport.Request` exposes the form, the target, and the run fact only (`internal/testreport/selection_facts.go:23-42`). TF-C2 adds the accessor that reports the `--changed` flag. TF-C3 exports the root-conformance check name from `internal/testreport`, so the preflight reads the one spelling of that name.

`internal/tickets` exports one predicate that answers whether a path sits under a bound command package. A bound command package is a binding row whose files are the command registry set (`internal/tickets/registry_data.go:21-58`).

The embed flag fires for any edit in a package that declares an embed pattern. Today those packages are `internal/adopt`, `internal/releaseevidence`, `internal/preflight/preflighttest`, and the root package. This choice is conservative, and flagged item 5 records it.

### The check-floor row

The build preflight composes each ticket's floor from three parts:

- the element `root-conformance`, when the gate predicate applies to the graded root;
- each package pattern from the testreport floor;
- the element `./cmd/bench`, when an entry or a floor package is in a bound command package, or the floor's embed flag is set.

The gatherer keeps the parsed completion plan beside its digest. For each requirement, it parses the command with the shell parser in `internal/shellcommand`.
A command satisfies an element only when the stream is lexed, holds one simple command, and its words start with `bench test`. The remaining words then go through `testreport.Prepare`.

The decision applies these satisfaction rules:

- The root-conformance element needs the check form with the target `root-conformance`.
- A package element needs the package form, a run fact of `all`, and no `--changed` flag.
- The package target equals the element, or it ends in `/...` and covers the element at a segment boundary. So `./internal/x/...` covers the element `./internal/x` itself.
- The target `./...` satisfies every package element, the root element `./.` included.
- A command that runs no test, such as the check inventory form or a fixture listing, satisfies nothing.

A ticket that no plan chunk names is graded against an empty requirement set. A version 1 plan grades the union of a chunk's ticket floors against that chunk's requirements. A delegated version 2 plan grades each ticket's floor against the chunk requirements that name that ticket.

`Decide` adds the row `check-floor` directly after `completion-plan`, in build mode and in review mode. The row shares the ticket-directory gate of the other ticket rows. It answers not-applicable when the completion plan is invalid.
Its red detail lists `<ticket>: <element> (add <command>)` for each unsatisfied element, in ticket order. The command is `bench test --check root-conformance` or `bench test --package <pattern>`. A package-index failure reds the row with the loader message.

### The planned-cell refusal

`internal/coverage` exports a projection that returns the row IDs whose seam cell matches the existing marker expression. The projection reads opt-in maps only. It reuses `plannedRe` and adds no second spelling.

A seam cell that starts with `review-owned:` is never projected, even when a later word is `planned`. The review-owned rule is a prefix (`internal/coverage/citations.go:352` and `:369`), while the marker expression matches anywhere (`citations.go:35`). So the prefix wins.

`bench record chunk` reads the spec and the plan from the tip commit before it writes. It does not read the working copy, as `RecordChunk` reads the plan at the tip (`internal/reviewrecord/write.go:125-129`).
It refuses when any row of the chunk's planned rows holds the marker. The refusal names each row and the form `` `<path>_test.go` (`<Name>`) ``.

`bench record completion` applies the same refusal to every row of every plan chunk at its source. The record command imports the coverage package. The checkpoint packages do not, because `internal/coverage` already depends on `internal/gate`.

`bench coverage --check` keeps its current behavior. A staged spec keeps its `planned` cells until the covering ticket commits.

### Guidance

The implement-spec guidance states that the completion plan carries each ticket's floor and that `check-floor` grades it. It also states that the ticket commit rewrites each covered `planned` cell into the citation form.
Map discipline states that the marker lasts until the covering ticket commits. Delegation discipline names `bench test --check root-conformance` as the root conformance pass, and it keeps its anchored sentence.

The build-entry route at `.agents/commands/bench-implement-spec.md:19` sends only `*-closure`, `fence-writes`, and `completion-plan` reds to the staleness pass. TF-C5 adds `check-floor` to that sentence. It updates the anchor needle at `internal/anchors/registry_chunk_chain.go:11` and the `old` text of the `implement-spec-red-preflight-route` canary to match. A staged spec then amends its plan in its enabling plan commit.

`.agents/commands/bench-implement-spec.md` holds 83 lines against its budget of 83 (`projects/benchkit.md:491`). Each new sentence joins an existing paragraph line, so the line count stays 83 and `projects/benchkit.md` stays outside the fence.
The floor sentence replaces the focused-checks sentence on line 44, so that paragraph keeps six sentences. The planned-cell rewrite sentence joins the four-sentence paragraph on line 48.

### Unchanged surfaces

The commit lane, its class table, and ADR 0017 stay unchanged. The chunk checkpoint keeps its record rules. FT89's open question stays open.

## Implementation chunks

Each chunk is one review checkpoint. TF-C1 and TF-C2 create seams that TF-C3 consumes, so each is its own chunk.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| TF-C1 / `01-name-the-root-conformance-check.md` | One exact root conformance command | TF1, TF2, TF3, TF4, TF5 | Completion-plan checks | no |
| TF-C2 / `02-derive-the-ticket-check-floor.md` | The floor derives from the Writes line | TF7, TF8, TF9, TF10, TF11, TF12, TF13, TF14, TF15, TF19, TF20, TF21, TF22, TF71 | Completion-plan checks | yes |
| TF-C3 / `03-grade-the-plan-against-the-floor.md` | The preflight grades the plan | TF6, TF16, TF17, TF18, TF23 to TF46, TF66, TF68, TF69, TF70, TF72, TF73, TF74 | Completion-plan checks | yes |
| TF-C4 / `04-refuse-a-planned-cell-at-the-freeze.md` | The chunk freeze refuses a planned cell | TF47 to TF59, TF67 | Completion-plan checks | no |
| TF-C5 / `05-state-the-floor-in-the-guidance.md` | Guidance states the floor | TF60 to TF65 | Completion-plan checks | no |

### Completion plan

The commands and probes below are future obligations, not executed proof. Each chunk requirement set already satisfies the floor that TF-C3 grades.

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft349-build-20261008","orchestrator_session":"claude:ft349-orchestrator-20261008","author_limit":1,"assignments":{"01-name-the-root-conformance-check.md":[{"session":"claude:ft349_t1","assignment":"deb5b981218a3208327006daaf4697bd","model":"opus","effort":"high","source":"30ce4ef56333d40f14ac1f14527f9ef8897c67ca","native_ref":"claude-agent:ft349_t1"},{"session":"claude:ft349_t1_repair1","assignment":"deb5b981218a3208327006daaf4697bd","model":"opus","effort":"high","source":"2c3dd9d2e79cf29503deb3f32c8fdf80ff4a8524","native_ref":"claude-agent:ft349_t1_repair1","predecessor":"claude:ft349_t1","trigger":"user-directed","stopped":"the ticket 01 author session reported completion of its TF-C1 verification record commit 2bd1b405 and has no live child","preserved":"2c3dd9d2e79cf29503deb3f32c8fdf80ff4a8524"},{"session":"claude:ft349_t1_repair2","assignment":"deb5b981218a3208327006daaf4697bd","model":"opus","effort":"high","source":"7c308b1bd806cd5a94a451c02f90f5e8fd3037f7","native_ref":"claude-agent:ft349_t1_repair2","predecessor":"claude:ft349_t1_repair1","trigger":"user-directed","stopped":"the ticket 01 repair cycle 1 session reported completion of its verification record commit 1b265fab and has no live child","preserved":"7c308b1bd806cd5a94a451c02f90f5e8fd3037f7"}],"02-derive-the-ticket-check-floor.md":[{"session":"claude:ft349_t2","assignment":"deb5b981218a3208327006daaf4697bd","model":"opus","effort":"high","source":"07f1ee469e43e5d45a111cd5788df8c17f4d1cd2","native_ref":"claude-agent:ft349_t2"},{"session":"claude:ft349_t2_repair1","assignment":"deb5b981218a3208327006daaf4697bd","model":"opus","effort":"high","source":"938f04f9da2bea7f81a944e13b8878b1c3ef13c6","native_ref":"claude-agent:ft349_t2_repair1","predecessor":"claude:ft349_t2","trigger":"user-directed","stopped":"the ticket 02 author session reported completion of its TF-C2 verification record commit 130322c0 and has no live child","preserved":"938f04f9da2bea7f81a944e13b8878b1c3ef13c6"},{"session":"claude:ft349_t2_repair2","assignment":"deb5b981218a3208327006daaf4697bd","model":"opus","effort":"high","source":"6892a0522a0b4e238c48b347bdae32a472bd2609","native_ref":"claude-agent:ft349_t2_repair2","predecessor":"claude:ft349_t2_repair1","trigger":"user-directed","stopped":"the ticket 02 repair cycle 1 session reported completion of its verification record commit 90705564 and has no live child","preserved":"6892a0522a0b4e238c48b347bdae32a472bd2609"}],"03-grade-the-plan-against-the-floor.md":[{"session":"claude:ft349_t3","assignment":"deb5b981218a3208327006daaf4697bd","model":"opus","effort":"high","source":"a9ae270abf4bedb99ce56602086cb809c68164e1","native_ref":"claude-agent:ft349_t3"},{"session":"claude:ft349_t3_repair1","assignment":"deb5b981218a3208327006daaf4697bd","model":"opus","effort":"high","source":"d673bd766e420e68c48b36baae6c36d262a2c9b8","native_ref":"claude-agent:ft349_t3_repair1","predecessor":"claude:ft349_t3","trigger":"user-directed","stopped":"the ticket 03 author session reported completion of its TF-C3 verification record commit 2cd283a2 and has no live child","preserved":"d673bd766e420e68c48b36baae6c36d262a2c9b8"}],"04-refuse-a-planned-cell-at-the-freeze.md":[{"session":"claude:ft349_t4","assignment":"deb5b981218a3208327006daaf4697bd","model":"opus","effort":"high","source":"41df830c693ac196a88d66010379bf721608bbf4","native_ref":"claude-agent:ft349_t4"}],"05-state-the-floor-in-the-guidance.md":[]}},"chunks":[{"id":"TF-C1","tickets":["01-name-the-root-conformance-check.md"],"verification":[{"id":"pkg-internal-testreport","command":"bench test --package ./internal/testreport","ticket":"01-name-the-root-conformance-check.md"},{"id":"pkg-internal-gate","command":"bench test --package ./internal/gate","ticket":"01-name-the-root-conformance-check.md"},{"id":"pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"01-name-the-root-conformance-check.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"01-name-the-root-conformance-check.md"},{"id":"root-scope-probe","command":"bench test --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'","probe":"Swap the empty conformance scope for the check name in the root-conformance environment. TestRootConformanceCheckRunsUnscoped turns red.","ticket":"01-name-the-root-conformance-check.md"}]},{"id":"TF-C2","tickets":["02-derive-the-ticket-check-floor.md"],"verification":[{"id":"pkg-internal-testreport","command":"bench test --package ./internal/testreport","ticket":"02-derive-the-ticket-check-floor.md"},{"id":"pkg-internal-tickets","command":"bench test --package ./internal/tickets","ticket":"02-derive-the-ticket-check-floor.md"},{"id":"pkg-internal-gate","command":"bench test --package ./internal/gate","ticket":"02-derive-the-ticket-check-floor.md"},{"id":"pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"02-derive-the-ticket-check-floor.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"02-derive-the-ticket-check-floor.md"},{"id":"floor-probe","command":"bench test --package ./internal/testreport --run '^TestWritesFloor'","probe":"Omit the directory-entry rule from the floor function. TestWritesFloorCoversADirectoryEntry turns red.","ticket":"02-derive-the-ticket-check-floor.md"}]},{"id":"TF-C3","tickets":["03-grade-the-plan-against-the-floor.md"],"verification":[{"id":"pkg-internal-preflight","command":"bench test --package ./internal/preflight/...","ticket":"03-grade-the-plan-against-the-floor.md"},{"id":"pkg-internal-testreport","command":"bench test --package ./internal/testreport","ticket":"03-grade-the-plan-against-the-floor.md"},{"id":"pkg-internal-conformance","command":"bench test --package ./internal/conformance","ticket":"03-grade-the-plan-against-the-floor.md"},{"id":"pkg-internal-diff","command":"bench test --package ./internal/diff","ticket":"03-grade-the-plan-against-the-floor.md"},{"id":"pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"03-grade-the-plan-against-the-floor.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"03-grade-the-plan-against-the-floor.md"},{"id":"system","command":"bench test --check system","ticket":"03-grade-the-plan-against-the-floor.md"},{"id":"check-floor-probe","command":"bench test --package ./internal/preflight --run '^TestCheckFloor'","probe":"Omit the package elements from the check-floor comparison. TestCheckFloorRedsAnOmittedPackage turns red.","ticket":"03-grade-the-plan-against-the-floor.md"}]},{"id":"TF-C4","tickets":["04-refuse-a-planned-cell-at-the-freeze.md"],"verification":[{"id":"pkg-internal-coverage","command":"bench test --package ./internal/coverage","ticket":"04-refuse-a-planned-cell-at-the-freeze.md"},{"id":"pkg-internal-reviewrecord","command":"bench test --package ./internal/reviewrecord/...","ticket":"04-refuse-a-planned-cell-at-the-freeze.md"},{"id":"pkg-internal-conformance","command":"bench test --package ./internal/conformance","ticket":"04-refuse-a-planned-cell-at-the-freeze.md"},{"id":"pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"04-refuse-a-planned-cell-at-the-freeze.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"04-refuse-a-planned-cell-at-the-freeze.md"},{"id":"planned-cell-probe","command":"bench test --package ./internal/reviewrecord/recordcmd --run '^TestRecordChunkRefusesAPlannedCell'","probe":"Omit the planned-row refusal before the chunk write. TestRecordChunkRefusesAPlannedCell turns red.","ticket":"04-refuse-a-planned-cell-at-the-freeze.md"}]},{"id":"TF-C5","tickets":["05-state-the-floor-in-the-guidance.md"],"verification":[{"id":"pkg-internal-anchors","command":"bench test --package ./internal/anchors","ticket":"05-state-the-floor-in-the-guidance.md"},{"id":"pkg-internal-conformance","command":"bench test --package ./internal/conformance","ticket":"05-state-the-floor-in-the-guidance.md"},{"id":"pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"05-state-the-floor-in-the-guidance.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"05-state-the-floor-in-the-guidance.md"},{"id":"guidance-probe","command":"bench test --check docs-currency-workflow","probe":"Delete the new implement-spec floor sentence. bench test --check docs-currency-workflow turns red.","ticket":"05-state-the-floor-in-the-guidance.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/ticket-checkpoint-floor/spec.md"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"pkg-cmd-bench","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

A good test drives the real owner through its public seam with canned inputs. The floor tests feed canned package records to the shared index, as the existing selection cause tests do (`internal/testreport/selection_cause_test.go`).
The check-floor decision tests build `Facts` values directly, as `internal/preflight/decision_test.go` does. One command-level test drives `bench preflight build` over a fixture repository that holds a Go module and a plan.

The record tests drive `bench record chunk` and `bench record completion` over the review-record fixture repositories in `internal/reviewrecord/recordtest`. The guidance rows use anchor registrations, which `bench test --check docs-currency-workflow` grades, with a canary mutation each.
At its ticket commit, each anchor row (TF60 to TF63, TF65) rewrites its cell to cite `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`). That test bites every retained canary through its registered owner, and `docs-currency-workflow` owns this canary family (`internal/conformance/registry/registry.go:131`).

Plan the evidence now, and do not implement the feature to obtain a red. Each ticket observes one meaningful red and green for each of its rows. Each ticket rewrites the `planned` cells of its own rows into the citation form in its ticket commit.

### Seam diagram

    ticket Writes line ──▶ [ testreport floor + tickets command predicate + gate predicate ] ──▶ floor elements
                              ◀ tests attach here: canned package records, binding rows, kit roots
    completion plan ──▶ [ shellcommand + testreport.Prepare ] ──▶ requirement selections
    floor elements + selections ──▶ [ preflight Decide: check-floor ] ──▶ green, red, or not-applicable row
                              ◀ tests attach here: Facts values and one fixture repository
    spec at chunk tip ──▶ [ coverage planned projection ] ──▶ [ bench record chunk / completion ] ──▶ refusal or record
                              ◀ tests attach here: recordtest fixture repositories

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| TF1 | 1 | The root-conformance check runs the entry test with an empty scope | `internal/testreport/root_conformance_test.go` (`TestRootConformanceCheckRunsUnscoped`) | A run that passes the check name as the scope grades one check only. |
| TF2 | 2 | A planted conformance defect makes the root-conformance check exit 1 with that check's diagnostic | `internal/testreport/root_conformance_test.go` (`TestRootConformanceCheckPrintsFailures`) | A run that drops the diagnostic or the exit code hides the defect. |
| TF3 | 3 | The root-conformance check refuses a root where the gate predicate answers false | `internal/testreport/root_conformance_test.go` (`TestRootConformanceCheckRefusesAForeignRoot`) | A run in a linked repository starts a pass that its tests cannot honor. |
| TF4 | 4 | `bench test --checks` lists `root-conformance` with the kind `conformance` | `internal/testreport/root_conformance_test.go` (`TestChecksFaceListsRootConformance`) | An absent inventory row hides the command from discovery. |
| TF5 | 5 | Over a kit root, a root without the entry test, and a foreign root, the named check runs exactly where the gate predicate answers true | `internal/testreport/root_conformance_test.go` (`TestRootConformanceCheckFollowsTheGatePredicate`) | A private copy of the kit rule in testreport disagrees with the gate on one of the three roots. |
| TF6 | 6 | A kit root that declares the entry test gives each ticket the root-conformance element, also when an ambient `BENCH_KIT` names that kit | `internal/preflight/check_floor_test.go` (`TestCheckFloorAddsRootConformanceOnTheKit`) | A floor without the element lets a plan skip the live-tree pass. |
| TF7 | 7 | A present Go file entry adds the package of its directory | `internal/testreport/floor_test.go` (`TestWritesFloorAddsTheFilePackage`) | A floor that drops the package lets the plan omit its tests. |
| TF8 | 8 | A `(new)` Go file in a current package directory adds that package | `internal/testreport/floor_test.go` (`TestWritesFloorAddsTheNewFilePackage`) | An absent-file shortcut drops the package of a new file. |
| TF9 | 9 | A directory entry adds every current package at or below it | `internal/testreport/floor_test.go` (`TestWritesFloorCoversADirectoryEntry`) | A file-only rule lets a directory entry hide its packages. |
| TF10 | 9 | A directory entry `internal/git` adds no package from `internal/git2` | `internal/testreport/floor_test.go` (`TestWritesFloorKeepsTheSegmentBoundary`) | A raw prefix match adds a sibling package. |
| TF11 | 10 | An embed file entry adds its owner package and sets the embed flag | `internal/testreport/floor_test.go` (`TestWritesFloorAddsTheEmbedOwner`) | A Go-only rule drops the owner of an embedded asset. |
| TF12 | 11 | A `go.mod` entry adds the pattern `./...` | `internal/testreport/floor_test.go` (`TestWritesFloorAddsEveryPackageForModuleMetadata`) | A metadata change otherwise owes no package test. |
| TF13 | 12 | A Markdown entry that no package embeds adds no package | `internal/testreport/floor_test.go` (`TestWritesFloorIgnoresAGuidancePath`) | A floor that maps every path to a package demands a test that does not exist. |
| TF14 | 13 | A `(new)` entry under no current package directory adds no package and no error | `internal/testreport/floor_test.go` (`TestWritesFloorSkipsAnUnbornPackage`) | The `--changed` refusal for an unborn package would red every plan-only run. |
| TF15 | 14 | A path under a bound command package answers true from the tickets predicate | `internal/tickets/registry_data_test.go` (`TestCommandPackageReadsTheBindingRows`) | A second list of command packages can drift from the binding rows. |
| TF16 | 14 | A command-package entry adds the element `./cmd/bench` to the preflight floor | `internal/preflight/check_floor_test.go` (`TestCheckFloorAddsCmdBenchForACommandPackage`) | A floor without the element lets a public response change skip the command tests. |
| TF17 | 15 | An entry that selects a package with an embed pattern adds `./cmd/bench` | `internal/preflight/check_floor_test.go` (`TestCheckFloorAddsCmdBenchForAnEmbedChange`) | A floor without the element lets an embed change skip the command tests. |
| TF18 | 16 | An entry outside every command package and every embed package adds no `./cmd/bench` | `internal/preflight/check_floor_test.go` (`TestCheckFloorOmitsCmdBenchOtherwise`) | An unconditional element demands a slow run for every ticket. |
| TF19 | 17 | A root without `go.mod` gives an empty package floor and starts no `go list` | `internal/testreport/floor_test.go` (`TestWritesFloorSkipsANonGoRoot`) | A loader call in a non-Go root reds every linked repository. |
| TF20 | 18 | The floor and the `--changed` selection give one direct package for one canned path | `internal/testreport/floor_test.go` (`TestWritesFloorAgreesWithChangedSelection`) | A second path-to-package map can disagree with the selection. |
| TF21 | 19 | The existing embed cause test stays green after the index is shared | `internal/testreport/selection_cause_test.go` (`TestCauseEmbedOnly`) | A refactor that loses the embed map changes the selection cause. |
| TF22 | 19 | The existing changed-beats-imports cause test stays green after the index is shared | `internal/testreport/selection_cause_test.go` (`TestCauseChangedBeatsImports`) | A refactor that loses the cause precedence changes the selection cause. |
| TF23 | 20 | A chunk that omits one floor element makes `check-floor` red in three fixtures: a package, `./cmd/bench` alone, and the second ticket's package in a two-ticket version 1 chunk | `internal/preflight/check_floor_test.go` (`TestCheckFloorRedsAnOmittedPackage`) | A row that never compares packages, skips the command element, or reads only the first ticket passes the observed omission. |
| TF24 | 20 | A plan whose chunk satisfies every floor element makes `check-floor` green | `internal/preflight/check_floor_test.go` (`TestCheckFloorGreenWhenSatisfied`) | A row that is always red blocks every build. |
| TF25 | 21 | The red detail reads `<ticket>: <element> (add <command>)` with the exact command token | `internal/preflight/check_floor_test.go` (`TestCheckFloorDetailNamesTheCommand`) | A detail without the command leaves the plan repair to guesswork. |
| TF26 | 22 | `bench test --package ./internal/x` satisfies the element `./internal/x` | `internal/preflight/check_floor_test.go` (`TestCheckFloorAcceptsTheExactPackage`) | A strict string rule that rejects the usual spelling reds a correct plan. |
| TF27 | 23 | `bench test --package ./internal/...` satisfies the element `./internal/x`, and `bench test --package ./internal/x/...` satisfies `./internal/x` itself | `internal/preflight/check_floor_test.go` (`TestCheckFloorAcceptsAParentPattern`) | An exact-match rule rejects a broad command that covers the package. |
| TF28 | 24 | `bench test --package ./...` satisfies every package element | `internal/preflight/check_floor_test.go` (`TestCheckFloorAcceptsTheWholeSuite`) | A rule that ignores `./...` rejects the widest command. |
| TF29 | 25 | `bench test --package ./internal/git/...` leaves `./internal/git2` unsatisfied | `internal/preflight/check_floor_test.go` (`TestCheckFloorKeepsTheSegmentBoundary`) | A raw prefix rule accepts a sibling package. |
| TF30 | 26 | A package command with a quoted `--run` pattern leaves its package unsatisfied | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsAFilteredRun`) | A rule that ignores the run fact accepts one test for a whole package. |
| TF31 | 27 | A `bench test --changed` command leaves every element unsatisfied | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsAChangedRun`) | A rule that reads the default target accepts a base-dependent selection. |
| TF32 | 28 | A command with `&&` before `bench test` satisfies nothing | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsACompoundCommand`) | A word scan that skips operators accepts a hidden compound line. |
| TF33 | 28 | A command whose first word is not `bench`, in fixtures `go test` and `BENCH_KIT=x bench test`, satisfies nothing | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsAnotherExecutable`) | A rule that searches for `test` anywhere accepts `go test`. |
| TF34 | 29 | A `bench test` command that the grammar refuses satisfies nothing | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsARefusedCommand`) | A rule that skips the grammar accepts a malformed command. |
| TF35 | 30 | `bench test --check root-conformance` satisfies the root-conformance element | `internal/preflight/check_floor_test.go` (`TestCheckFloorAcceptsTheRootConformanceCheck`) | A rule that rejects the canonical command reds every plan. |
| TF36 | 30 | `bench test --check docs-currency-workflow` leaves the root-conformance element unsatisfied | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsAScopedConformanceCheck`) | A rule that accepts any check form lets one scope stand in for the pass. |
| TF37 | 30 | `bench test --package ./internal/conformance` leaves the root-conformance element unsatisfied | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsTheConformancePackage`) | The package run skips the entry test without its root variable. |
| TF38 | 31 | A requirement in another chunk leaves this chunk's element unsatisfied | `internal/preflight/check_floor_test.go` (`TestCheckFloorReadsOnlyTheOwnChunk`) | A plan-wide union lets one chunk borrow another chunk's checks. |
| TF39 | 32 | A delegated requirement that names another ticket leaves this ticket's element unsatisfied | `internal/preflight/check_floor_test.go` (`TestCheckFloorReadsOnlyTheOwnTicket`) | A chunk-wide rule lets one author borrow another author's checks. |
| TF40 | 33 | The default `bench preflight build` form prints a red `check-floor` row for an omitted package | `internal/preflight/check_floor_command_test.go` (`TestCheckFloorCommandRunsAfterTheTicket`) | A plan-only placement misses a package that the ticket commit created. |
| TF41 | 34 | The review form prints a red `check-floor` row for an omitted package | `internal/preflight/check_floor_command_test.go` (`TestCheckFloorCommandRunsInReview`) | A build-only placement lets the last chunk reach review with an incomplete plan. |
| TF42 | 35 | A red `check-floor` row refuses the build charge with `preflight required: check-floor` | `internal/preflight/check_floor_command_test.go` (`TestCheckFloorRefusesTheBuildCharge`) | A charge that ignores the row starts the next author on an incomplete plan. |
| TF43 | 36 | An invalid completion plan makes `check-floor` not-applicable | `internal/preflight/check_floor_test.go` (`TestCheckFloorDefersToTheCompletionPlan`) | A second red for the same defect hides which row owns it. |
| TF44 | 37 | Build mode with no tickets directory makes `check-floor` not-applicable | `internal/preflight/check_floor_test.go` (`TestCheckFloorNotApplicableWithoutTickets`) | A graded row on a fresh spec stage reds before tickets exist. |
| TF45 | 38 | A package-index failure makes `check-floor` red with the loader message | `internal/preflight/check_floor_test.go` (`TestCheckFloorRedsALoaderFailure`) | A swallowed loader error gives an empty floor and a false green. |
| TF46 | 39 | The kit lane table keeps its current check list | `internal/gate/lane_test.go` (`TestBenchkitLaneTable`) | A lane edit for this outcome changes the pinned table. |
| TF47 | 40 | `bench record chunk` exits 1 when a chunk row holds a `planned` seam cell at the tip | `internal/reviewrecord/recordcmd/planned_cell_test.go` (`TestRecordChunkRefusesAPlannedCell`) | A freeze without the check records a chunk whose rows cite no test. |
| TF48 | 41 | The chunk refusal names each planned row ID and the citation form | `internal/reviewrecord/recordcmd/planned_cell_test.go` (`TestRecordChunkNamesEachPlannedRow`) | A refusal without the row IDs leaves the repair to guesswork. |
| TF49 | 42 | A refused chunk freeze leaves the review record bytes unchanged | `internal/reviewrecord/recordcmd/planned_cell_test.go` (`TestRecordChunkRefusalWritesNothing`) | A refusal after the write leaves a partial chunk entry. |
| TF50 | 43 | A chunk row whose seam cell uses the citation form lets the freeze record | `internal/reviewrecord/recordcmd/planned_cell_test.go` (`TestRecordChunkAcceptsACitedRow`) | A refusal that ignores the marker blocks every freeze. |
| TF51 | 44 | A chunk row whose seam cell reads `review-owned: the planned canary stays green` lets the freeze record | `internal/reviewrecord/recordcmd/planned_cell_test.go` (`TestRecordChunkAcceptsAReviewOwnedRow`) | A refusal that lets the marker beat the prefix blocks a review-owned row. |
| TF52 | 45 | A `planned` row outside the chunk's rows lets the freeze record | `internal/reviewrecord/recordcmd/planned_cell_test.go` (`TestRecordChunkIgnoresAnotherChunksRow`) | A spec-wide refusal blocks a later chunk's future test. |
| TF53 | 46 | A seam cell that holds `unplanned` is not projected as planned | `internal/coverage/planned_rows_test.go` (`TestPlannedRowsKeepsTheWordBoundary`) | A substring match refuses a cell without the marker. |
| TF54 | 47 | `bench record completion` exits 1 when any plan row holds a `planned` seam cell | `internal/reviewrecord/recordcmd/planned_cell_test.go` (`TestRecordCompletionRefusesAPlannedCell`) | A completion without the check lands a spec whose rows cite no test. |
| TF55 | 47 | `bench record completion` records when every plan row is cited | `internal/reviewrecord/recordcmd/planned_cell_test.go` (`TestRecordCompletionAcceptsCitedRows`) | A refusal that ignores the marker blocks every completion. |
| TF56 | 48 | `bench coverage --check` keeps zero violations for a planned bare-name cell | `internal/coverage/citation_form_test.go` (`TestSeamCellTestNameOutsideTheCitationForm`) | A coverage-level refusal reds every staged spec. |
| TF57 | 49 | The projection returns the opt-in row IDs whose seam cell holds the marker | `internal/coverage/planned_rows_test.go` (`TestPlannedRowsReadsTheMarkerOwner`) | A projection that reads another cell or another spelling misses the marker. |
| TF58 | 49 | A non-opt-in map projects no row ID | `internal/coverage/planned_rows_test.go` (`TestPlannedRowsSkipsANonOptInMap`) | A positional ID cannot join a ticket's `Covers:` token. |
| TF59 | 49 | A walk of every non-test Go file that `go list ./...` names finds the expression text `\bplanned\b` only in `internal/coverage` | `internal/coverage/planned_rows_test.go` (`TestPlannedMarkerHasOneOwner`) | A planted copy of the expression in `internal/reviewrecord/recordcmd` turns the walk red, and the bare word in `ReceiptPhasePlanned` does not. |
| TF60 | 50 | The implement-spec guidance holds the anchored floor sentence | planned anchor family in internal/anchors/registry_check_floor.go, joined in registry_data.go, with a workflow-guidance-anchors canary, graded by docs-currency-workflow | A deleted sentence leaves the floor without guidance. |
| TF61 | 51 | The implement-spec guidance holds the anchored planned-cell rewrite sentence | planned anchor family in internal/anchors/registry_check_floor.go, joined in registry_data.go, with a workflow-guidance-anchors canary, graded by docs-currency-workflow | A deleted sentence leaves authors to meet the refusal unwarned. |
| TF62 | 52 | Map discipline holds the anchored marker-lifetime sentence | planned anchor family in internal/anchors/registry_check_floor.go, joined in registry_data.go, with a workflow-guidance-anchors canary, graded by docs-currency-workflow | A deleted sentence leaves spec authors without the marker rule. |
| TF63 | 53 | Delegation discipline holds the anchored root-conformance command sentence | planned anchor family in internal/anchors/registry_check_floor.go, joined in registry_data.go, with a workflow-guidance-anchors canary, graded by docs-currency-workflow | A deleted sentence leaves charges without the exact command. |
| TF64 | 53 | The existing `delegate-root-conformance-pass` canary bites through `docs-currency-workflow` after the guidance edit | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An edit of the anchored needle leaves the canary without a bite. |
| TF65 | 55 | The build-entry route sentence of the implement phase names `check-floor` among the reds that the staleness pass takes | planned updated needle in internal/anchors/registry_chunk_chain.go with the implement-spec-red-preflight-route canary, graded by docs-currency-workflow | An unchanged route stops the phase on a `check-floor` red, so no staged spec can amend its plan. |
| TF66 | 6 | A linked root with `go.mod` and no declared entry test gives no root-conformance element | `internal/preflight/check_floor_test.go` (`TestCheckFloorOmitsRootConformanceOffTheKit`) | A floor that always adds the element demands a pass that the linked root cannot run. |
| TF67 | 40 | A tip that cites the row lets `bench record chunk` record while the working copy still holds `planned` | `internal/reviewrecord/recordcmd/planned_cell_test.go` (`TestRecordChunkReadsTheTip`) | A refusal that reads the working copy grades bytes that the frozen chunk does not hold. |
| TF68 | 25 | `bench test --package ./internal/git2` leaves the element `./internal/git` unsatisfied | `internal/preflight/check_floor_test.go` (`TestCheckFloorKeepsTheReverseBoundary`) | A prefix rule in the other direction accepts a sibling package. |
| TF69 | 20 | `bench preflight build <slug> --plan-only` prints a red `check-floor` row for an omitted package | `internal/preflight/check_floor_command_test.go` (`TestCheckFloorCommandRunsAtPlanOnly`) | A row that runs only after a ticket misses the plan defect before dispatch. |
| TF70 | 21 | The rendered `check-floor` red holds no raw control byte inside its row | `internal/preflight/check_floor_command_test.go` (`TestCheckFloorDetailRendersEscaped`) | A detail written past the renderer can split its own row. |
| TF71 | 13 | An absent `(new)` embed file adds no package and no error | `internal/testreport/floor_test.go` (`TestWritesFloorSkipsAnUnbornEmbedFile`) | The `--changed` refusal for an absent path would red a plan before the asset exists. |
| TF72 | 29 | `bench test --checks` and `bench test --check root-conformance --fixtures` satisfy nothing | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsANonRunningForm`) | A rule that reads only the form and the target accepts a command that runs no test. |
| TF73 | 32 | A ticket file that no plan chunk names makes `check-floor` red for each element of its floor | `internal/preflight/check_floor_test.go` (`TestCheckFloorRedsAnUnplannedTicket`) | A gather that skips an unplanned ticket leaves its whole floor ungraded under a green row. |
| TF74 | 14 | A directory entry above a bound command package adds the element `./cmd/bench` | `internal/preflight/check_floor_test.go` (`TestCheckFloorAddsCmdBenchForADirectoryAboveACommandPackage`) | A rule that reads only the entry path misses a command package that the directory entry writes. |

Not covered: story 54 — FT89 owns the decision about the per-return gate run, and this spec writes no surface that decides it.

### Edge inventory

The edge walk covers both audiences: this repository and every repository that links the kit.

- A linked repository without `go.mod` gets an empty package floor. Row TF19 covers it.
- A linked repository with `go.mod` gets package elements and no root-conformance element. Rows TF3 and TF66 cover the predicate.
- A plan command with quotes goes through the shell parser. Row TF30 uses a quoted `--run` pattern.
- A plan command with an assignment prefix satisfies nothing, because its first word is not `bench`. Row TF33 holds that fixture.
- A red detail renders through the existing verdict renderer. Row TF70 asserts that the rendered row holds no raw control byte.

**Won't handle:**

- A `(new)` directory entry with no current package — the after-ticket preflight (TF40) derives that package once the ticket commit creates it.
- A Go package directory whose name holds a space — Go refuses that import path, and the landing gate's `go test ./...` still grades any such package.
- An executed-check record — the chunk checkpoint (`internal/reviewrecord/check.go:27-56`) already refuses a planned command without a passing record.
- The npm pack check from FT337 — `bench prep-release` keeps that release check, and the floor derives only Go packages and root conformance.
- A non-opt-in coverage map — the build preflight bootstrap (`internal/preflight/gather.go:65-67`) already refuses that map.
- A `planned` cell on a row that no ticket covers — `rows-owned` (`internal/preflight/decision.go:375-387`) already reds that row.
- An absent `(new)` embed file — the after-ticket preflight (TF40) adds its owner once the file exists.

## Ownership fences

This fence is the exact union of the ticket Writes lines, plus the review pickup.

- `.agents/commands/bench-implement-spec.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-spec/references/map-discipline.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/preflight_version_test.go`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_charge_binding.go`
- `internal/anchors/registry_charge_binding_test.go`
- `internal/anchors/registry_check_floor.go`
- `internal/anchors/registry_check_floor_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_spec_trace_test.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/coverage/citations.go`
- `internal/coverage/planned_rows.go`
- `internal/coverage/planned_rows_test.go`
- `internal/gate/gate.go`
- `internal/gate/phases.go`
- `internal/preflight/check_floor.go`
- `internal/preflight/check_floor_command_test.go`
- `internal/preflight/check_floor_test.go`
- `internal/preflight/command_bootstrap_test.go`
- `internal/preflight/command_build_test.go`
- `internal/preflight/command_review_test.go`
- `internal/preflight/decision.go`
- `internal/preflight/decision_test.go`
- `internal/preflight/gather.go`
- `internal/preflight/plan.go`
- `internal/preflight/preflighttest/testdata/review-documents.toon`
- `internal/preflight/preflighttest/testdata/review-empty.toon`
- `internal/preflight/preflighttest/testdata/review-shapes-no-renames.toon`
- `internal/preflight/preflighttest/testdata/review-shapes-renames.toon`
- `internal/preflight/preflighttest/fixture.go`
- `internal/preflight/source_tip_test.go`
- `internal/preflight/ticket_rows.go`
- `internal/reviewrecord/completion.go`
- `internal/reviewrecord/coverage.go`
- `internal/reviewrecord/plan.go`
- `internal/reviewrecord/recordcmd/chunk.go`
- `internal/reviewrecord/recordcmd/command.go`
- `internal/reviewrecord/recordcmd/completion.go`
- `internal/reviewrecord/recordcmd/planned_cell.go`
- `internal/reviewrecord/recordcmd/planned_cell_test.go`
- `internal/reviewrecord/recordtest/fixture.go`
- `internal/testreport/check_test.go`
- `internal/testreport/checks_face_test.go`
- `internal/testreport/command.go`
- `internal/testreport/environment.go`
- `internal/testreport/floor.go`
- `internal/testreport/floor_test.go`
- `internal/testreport/named_check.go`
- `internal/testreport/outcome_test.go`
- `internal/testreport/root_conformance_test.go`
- `internal/testreport/selection.go`
- `internal/testreport/selection_facts.go`
- `internal/testreport/selection_facts_test.go`
- `internal/tickets/registry_data.go`
- `internal/tickets/registry_data_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `reviews/ticket-checkpoint-floor.md`

Reviewer disposition: approved on 2026-10-08 as the exact union above.

## Out of scope

- Lane-time execution of the floor. FT215 and FT380 keep that decision open. Estimate: 9 edits, 3 gate runs.
- A persisted record of executed `bench test` runs, so that exit codes stop being author-stated. Estimate: 12 edits, 4 gate runs.
- A bulk rewrite of the `planned` cells in the 21 staged specs. Each spec's own build rewrites its cells under the new freeze rule. Estimate: 0 edits, 0 gate runs.
- The FT375 cited-tests row. It composes this spec's planned projection when its build starts. Estimate: 0 edits here, 0 gate runs.

## Further notes

### Source verification

The author read these sources at `8bb4b894` on 2026-10-08:

- `internal/testreport/selection.go:60-280`, `outcome.go:36-67`, `command.go:40-125`, `selection_facts.go:1-80`, and `named_check.go:40-110`.
- `internal/preflight/plan.go`, `decision.go:150-250`, `gather.go:1-120`, and `review.go:207-225`.
- `internal/reviewrecord/plan.go:15-170`, `check.go:1-80`, `coverage.go:50-80`, `write.go:125-162`, and `recordcmd/command.go:270-300`.
- `internal/gate/phases.go:130-152`, `internal/conformance/checks_test.go:97-127`, `internal/tickets/registry_data.go:1-92`, and `internal/packagesurface/assets.go:80-125`.
- `internal/shellcommand/shellcommand.go:1-135` and `internal/coverage/coverage.go:505-530`.

The review fold read these sources at `561da614`: `internal/testreport/command.go:127-145` and `:329-334`, `internal/gate/tag_census.go:221-231`, `internal/preflight/proposal.go:1-35`, `internal/reviewrecord/source.go`, and `internal/intent/ledger/ledger.go:164`.

`go list` confirmed the import edges. `internal/preflight` already depends on `internal/gate`. No package in the dependency closure of `internal/gate`, `internal/testreport`, or `internal/tickets` imports `internal/preflight`.
The dependency closure of `internal/coverage` holds `internal/reviewrecord` and not `internal/reviewrecord/recordcmd`. So the record command can import `internal/coverage`, and the checkpoint packages cannot.

The research for this spec is `specs/ticket-checkpoint-floor/assets/research.md`.

### Held evidence

The open learning at `capture/learnings.md:207-210` asks the author to read the held plan commits. A diff of the completion-plan commands across the roadmap-delivery-commitment build history found nine plan commits after build start that added a package check:

- `3ec313b0` added `./internal/commitment/repository` to DC-C5.
- `e3005a6e` added `./internal/commitment/repository` to DC-C6.
- `47b779cd` added `./internal/commitment` to DC-C6.
- `1a9d1c4d` added `./internal/commitment/repository` to DC-C7.
- `6db48cb0` added `./cmd/bench` and `./internal/conformance` to DC-C7.
- `5a341734` added `./internal/gate` to DC-C7.
- `54d2fbc5` added `./internal/commitment`, `./internal/intent`, and `./internal/usage` to DC-C8.
- `d7ad195d` added `./internal/commitment/repository` to DC-C8.
- `9ed6c6f8` added four packages to DC-C9.

The learning names six commits. This census found nine, and the six held commits are a subset of them. Each added check is a package that the ticket writes, or `./cmd/bench`, or the conformance package. Rows TF7, TF9, TF16, and TF23 cover those classes.

### Reader and writer sweep

- The named-check inventory is read by `internal/testreport/fixtures.go:58`, `checks_face_test.go`, and `check_test.go`. TF-C1 updates the inventory tests.
- The `Decide` row list is asserted by `decision_test.go`, `command_build_test.go`, `command_review_test.go`, and `command_bootstrap_test.go` in `internal/preflight`. TF-C3 updates each assertion.
- The `planned` marker is written by authors in 21 staged specs and read only by `internal/coverage/citations.go:35-150`.
- The phrase "root conformance" appears in guidance at `.agents/commands/bench-implement-spec.md:44` and `delegation-discipline.md:33` and `:127`. The anchor at `internal/anchors/registry_data.go:452` pins line 127.
- `bench anchors` shows that `bench-implement-spec.md:44` holds two anchored sentences before the floor sentence. The floor sentence itself carries no anchor. `map-discipline.md:124` is anchored, and line 125 is not.

### Source clauses and coverage

| source clause | coverage |
| --- | --- |
| The ticket checkpoint derives the floor from the ticket's `Writes:` line | TF6 to TF22, TF66, TF71 |
| It refuses a commit whose checks omitted a package in that set | TF23 to TF45, TF68 to TF70. The red stops the next build charge and the review charge, not the published commit (flagged item 1). |
| `FT349.floor`: root conformance | TF1 to TF6, TF35 to TF37, TF66 |
| `FT349.floor`: the whole `cmd/bench` package for a public response or embed change | TF16 to TF18 |
| `FT349.planned`: refuse a planned citation for a row that the committed ticket covers | TF47 to TF59, TF67. The refusal fires at the chunk freeze and at completion (flagged item 1). |
| The spec decides which surface owns each refusal | Implementation decisions, flagged item 1 |
| The build-entry route admits the new red | TF65 |

### Pre-review proof checklist

- Cited symbols: every current symbol was read. New symbols and tests carry the `planned` marker.
- Import edges: `go list` confirmed every new edge in the source verification above.
- Source-row clauses and occurrences: the source-clause table covers the three FT349 occurrences and both criteria.
- Promised field labels: the row name `check-floor` and the check name `root-conformance`.
- Changed-function callers: `selectCurrentPackages` has one caller, `resolveChangedPackagesWithLoader`. `RecordChunk` has one production caller in `recordcmd`. `Decide` has the preflight command callers only.
- Copy survival: TF20 and TF59 fail if a second path map or a second marker expression survives.
- Rendered-shape readers: the `bench test --checks` table and the preflight row list. The sweep above names their tests.
- Pin operators: the row-list tests compare with `len` and element equality.
- Entry reads: the floor's `go list` call takes an injected loader in tests, as `resolveChangedPackagesWithLoader` does.
- Derived expectations: the floor tests use canned package records, not the code under test.
- Consolidated rules: `testBenchSource` calls `gate.KitRoot`, so `internal/testreport` keeps no private copy of the kit rule.
- Quantified obligations: "every package" in TF9 and TF12, and "each row" in TF48, each take a two-member fixture.
- Workflow-step writes: the freeze refusal writes nothing. Row TF49 compares the record bytes.

### Flagged items for reviewer veto

The reviewer approved the spec and the tickets on 2026-10-08. The reviewer accepted the seam choice, the new `root-conformance` check name, and the bound-command predicate with its known misses. The reviewer also kept ticket 03 whole and accepted the new export names. No item below was vetoed.

1. **Seam choice.** The floor refusal lives in the build preflight as `check-floor`. The planned refusal lives in `bench record chunk` and `bench record completion`. The commit lane stays unchanged. The row reds after the ticket commit and blocks the next charge and the review charge. It does not refuse the publication itself.
2. **New command.** `bench test --check root-conformance` is a new public check name. No current command runs the unscoped pass, and `FT349.floor` needs one exact command to name an omission.
3. **Root conformance spelling.** Only `bench test --check root-conformance` satisfies the element. Three staged plans spell it `bench test --check conformance-meta`, and others use `bench test --package ./internal/conformance`. Each staged spec amends its plan at its own build entry, through the route that TF65 opens.
4. **Non-Go and directory entries.** A non-Go path owes no package. A directory entry owes each current package below it. An unborn package owes nothing until the ticket commit creates it.
5. **Public response predicate.** The exact predicate is "an entry under a bound command package" from `internal/tickets/registry_data.go:21-58`. It stands for "a changed public response".
   - Known misses: `internal/testreport`, `internal/preflight`, `internal/gate`, and `internal/reviewrecord/recordcmd` change public responses, but they derive no `./cmd/bench`.
   - The embed flag fires for any Go edit in a package that declares an embed pattern: `internal/adopt`, `internal/releaseevidence`, `internal/preflight/preflighttest`, and the root package. The spec picks this conservative form over a narrower "entry is an embed file" form.
6. **One source.** The floor shares the `--changed` package index and excludes its reverse-import closure.
7. **FT375 composition.** `check-floor` changes the ordinary build preflight, and `specs/preflight-spec-staleness/spec.md:208` says that preflight "remains unchanged". That sentence scopes FT375's own audit rows, so the two do not conflict. FT375 reconciles its premise through its enabling plan commit. The planned refusal does not touch the after-ticket preflight. FT375's cited-tests row composes this spec's planned projection.
8. **Staged-spec migration.** No staged spec changes now. `bench coverage --check` and its pinned test at `internal/coverage/citation_form_test.go:81` stay unchanged, so ADR 0006 needs no new red. The 21 staged specs meet both refusals at their own builds, and the TF65 route lets their enabling plan commit take the `check-floor` red.
9. **Sequencing.** FT293 and FT375 fence `internal/preflight`, `internal/coverage`, `internal/tickets`, and `map-discipline.md`. They pin source censuses that this spec makes stale. They reconcile through their own enabling-plan procedure.
10. **Plan granularity.** A version 1 plan grades each chunk. A delegated version 2 plan grades each ticket.
