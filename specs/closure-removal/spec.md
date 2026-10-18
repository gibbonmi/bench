# Remove the ownership closure and widen ticket writes from the gate

Status: staged

Roadmap: FT293

Decision source: roadmap/FT293.md (named reviewed artifact; reviewer decision 2026-10-17, drain d-6c97ed54c629)

Verification log: pending

## Problem

Build preflight predicts the files that each ticket must own. It copies three couplings into a static closure: fixture pins, a binding registry, and anchor registries. Each build then finds a coupling kind that the copy lacks. A compile error, a test, or `bench structure` already reports each such coupling, so the copy is a second source of one fact.

The prediction also blocks work. A closure red stops a build charge. The `--propose-writes` and `--breakdown` forms add two more views of the same prediction. The `charge-premises` row stops a charge until every declared premise has proof.

The spec also holds a second copy of the ticket writes. Its `## Ownership fences` section must equal the union of the ticket `Writes:` lines, and the `fence-writes` row grades the two copies against each other. A gate-found need therefore takes a hand edit of two files.

The delegate guidance contradicts itself. One sentence calls the fence a refusal boundary. The next sentence calls `Writes:` an expectation.

## Solution

Preflight stops the prediction. The closure rows, the `charge-premises` row, `--propose-writes`, and `--breakdown` go. The consistency rows and the landing path check stay.

The tickets own the fence. The union of the ticket `Writes:` entries, plus the implicit entries, authorizes each changed path at the landing. A new spec holds no fence section, and no new code reads one. This spec keeps a live section until its retirement, because the installed broker that lands it still reads the section.

One plan verb widens a ticket. `bench spec widen` adds a gate-found path to one ticket's `Writes:` line when the path is inside the approved scope. The approved scope is the directory of each signed-off `Writes:` entry. A path outside that scope goes to the reviewer, and `--decision <reference>` records the reviewer's admission.

The guidance states one rule. `Writes:` is an expectation, the gate finds each missing need, and the plan verb widens the expectation.

## User stories

Line: `opus` / high.
Implementation-line reason: the hardest material chunk is CR3, which moves the landing authority to the ticket union under a pre-FT293 landing broker and reads ticket writes from git for the staleness audit. The code deletions have exact inventories and strong compile and test coverage. The new verb is a small seam with a known precedent. The leverage override routes the guidance chunk to high effort.
Harder chunks: CR3.

### Preflight predicts no ownership closure

1. As a slicer, I want build preflight to emit no `fixture-closure`, `registry-closure`, or `anchor-closure` row, so that a static copy of a coupling never reds a ticket graph.
2. As an orchestrator, I want a build charge to prepare evidence without a closure verdict, so that a gate-found need drives each expansion.
3. As an orchestrator, I want a build charge to emit no `charge-premises` row, so that preflight checks no charge premises.
4. As a slicer, I want the `--propose-writes` form gone, so that no second proposal of ticket writes competes with the gate.
5. As a slicer, I want the `--breakdown` form gone, so that the ordinary build verdict is the one view of uncited rows.
6. As an agent, I want `bench help` to list neither removed form, so that no session discovers a dead form.
7. As an orchestrator, I want `--audit-spec` to keep its premise, manual-claim, symbol, and test faces, so that the staleness pass keeps its evidence.
8. As an orchestrator, I want the audit to report no closure diagnostic, so that a deleted prediction cannot keep manual review required.
9. As an operator, I want a build verdict that reads no closure inventory, so that a fixture or anchor inventory fault cannot refuse a verdict that does not use it.
10. As a maintainer, I want the binding registry deleted with its only reader, so that no static copy of a registry coupling survives.
11. As a maintainer, I want `writes-resolve`, `kit-pin`, and `rows-owned` to keep their verdicts, so that the removal weakens no kept check.
12. As a maintainer, I want the ticket grammar to refuse an unclean `Writes:` path, so that the clean-path refusal survives the closure that owned it.
13. As a maintainer, I want the preflight row-list pins to name the new row order, so that a returned closure row reds a test.
52. As a maintainer, I want the `ticket-grammar` check to keep its ticket-folder sweep when its binding-table grader goes, so that the gate still grades every staged ticket.

### A gate-found need widens a ticket's expected writes

14. As an orchestrator, I want `bench spec widen` to add an in-scope path to one ticket's `Writes:` line, so that a gate-found need joins the expectation without a reviewer stop.
15. As a reviewer, I want an out-of-scope path refused with my route named, so that a widening past the approved scope comes to me.
16. As an orchestrator, I want `--decision <reference>` to admit an out-of-scope path after the reviewer decides, so that one verb writes every widening.
17. As a reviewer, I want the approved scope read from the tickets that the default-branch tip holds, so that a build cannot widen the scope that grades its own widening.
18. As an orchestrator, I want the directory of each signed-off entry to bound the scope, so that a test or a sibling file beside a planned file needs no stop.
19. As an orchestrator, I want an entry that the ticket already covers to write nothing, so that a repeated call is harmless.
20. As an orchestrator, I want a request with one refused entry to write nothing, so that a partial widening never reaches the tree.
21. As an author, I want every other byte of the ticket kept, so that the widening changes only its `Writes:` value.
22. As an operator, I want the verb to refuse the primary checkout, so that a widening reaches `main` only through a worktree landing.
23. As an operator, I want the verb to refuse a hostile entry or ticket file, so that a bad widening cannot corrupt the `Writes:` grammar.
24. As an agent, I want the verb listed in `bench help` and in `bench spec --help`, so that a session discovers it without a trial call.
25. As an orchestrator, I want the verb to make no commit, so that the plan commit stays with `bench commit`.
26. As a reviewer, I want the verb to refuse when the default-branch tip holds no ticket of the slug, so that an absent scope never reads as an unrestricted one.
27. As an operator, I want the verb to work from a directory below the repository root, so that the working directory never changes which ticket it edits.

### The tickets own the fence

28. As a reviewer, I want the landing to authorize each changed path from the ticket `Writes:` union, so that the spec holds no second copy of the fence.
29. As an orchestrator, I want the review pickup authorized without a `Writes:` entry, so that the derived fence still admits the review record.
30. As a maintainer, I want the `fence-writes` row gone, so that no check compares two copies of one fact.
31. As an operator, I want a spec with no fence section to reach a verdict, and an empty union to authorize only the implicit entries, so that the derived fence stays fail-closed.
32. As an author, I want the charge evidence `fence` block to list the ticket union, so that the charge names the fence that the landing grades.
33. As an orchestrator, I want the staleness audit scope to read the ticket `Writes:` entries at the base and at the tip, so that the audit scope follows the derived fence.
34. As an operator, I want the landing refusal route to name `bench spec widen`, so that the route names the one widening verb.
35. As a spec author, I want `bench coverage --check` to stop grading a fence section, so that a spec needs no fence section.
36. As a maintainer, I want the fence parser deleted, so that no code reads a spec fence.
37. As an operator, I want a leftover fence section to authorize nothing, so that a stale section cannot widen the landing.
38. As a reviewer, I want every other staged spec left unedited, so that no bound deliverable identity changes before its outcome starts.
39. As a reviewer, I want this spec to keep its own fence section equal to its ticket union for the installed broker, so that this landing passes the pre-FT293 path check.

### The guidance states one rule

40. As a write delegate, I want one rule that calls `Writes:` an expectation, so that no sentence calls the fence a refusal boundary.
41. As a slicer, I want the guidance that repeats a preflight check deleted, so that the check stays the one source of its rule.
42. As a slicer, I want no guidance step to run a removed form, so that the slicing procedure runs on the new tree.
43. As a spec author, I want the template and the approval table without a fence section, so that a new spec carries no second copy.
44. As an orchestrator, I want the plan-expansion rule to name `bench spec widen` and the scope rule, so that a gate-found need follows one route.
45. As an orchestrator, I want the implement-spec staleness route to name only kept rows, so that the route never waits on a deleted row.
46. As a maintainer, I want an anchor to forbid each retired command sentence, so that a removed form cannot return to the guidance.
47. As a reader, I want the known issue whose retire condition names the removed proposal deleted, so that `docs/known-issues.md` names no dead route.
48. As a reader, I want the changelog to state the removal and the new verb, so that the unreleased entry matches the tree.
49. As a maintainer, I want each budgeted guidance file to stay inside its line budget, so that the guidance edits keep the budget check green.
50. As a reviewer, I want the slicing research for likely `Writes:` entries kept, so that this removal changes no process that FT293 does not name.
51. As a reviewer, I want plan-commit grading left with FT380, so that this spec does not take another outcome's scope.

## Implementation decisions

### The removal

The closure family goes as one unit. `Decide` stops calling `closureRows`, and the closure kinds, the missing-closure projection, and the closure messages leave the preflight package. The ownership package loses `Requirement`, `Kind`, the collection stages, the closure walk, the grammar-effects chain, the caller and liveness censuses, `OwnedPaths`, `Owned`, `Unordered`, and the uncited projection. The binding registry in `internal/tickets` goes with it, because only the closure reads `BindingRules` and `BoundFiles`. The `ticket-grammar` conformance check keeps its ticket-folder sweep and loses its binding-table grader, because that grader only keeps the closure's table valid (`internal/conformance/ticket_grammar_test.go:121-154`). No canary pins the grader. The canary pin reader `canary.Pins` and the anchor reference reader `anchors.References` also go, because the closure loader is their only production caller (`internal/preflight/ownership/closure.go:91`, `:95`). The canary inventory and the anchor registries keep their formats (CR76).

The premise checks mean the `charge-premises` row and its tests. The premise collectors, the planning-facts manifest, `Snapshot.Symbols`, `ChargeReadiness`, and `SubjectPaths` stay, because `--audit-spec` reads them. `Collect` keeps its premise half and loses its closure half and its closure diagnostics. Only the audit kind gathers the snapshot, because no other verdict row reads it after the charge-premises row goes. This call is flagged for reviewer veto.

The kept helpers stay where the research found their kept readers: the tree probe and `treeHolds`, `deletedBefore`, `buildDeletedPaths`, `systemTagged`, `newSystemTest`, `UncitedIDs`, `grepFiles`, and `localName`. `KitObligations` goes, because only the proposal read it.

The closure also refused a `Writes:` entry whose path escapes the repository or is not clean. That refusal moves into `tickets.ParseWrite`, so `tickets-parse` reds such an entry. The verdict then reds the entry even when its path exists outside the repository. The premise collectors that used the closure's clean-path predicate call the grammar's rule.

`--propose-writes` and `--breakdown` leave the preflight grammar, the evidence operations table, the command dispatch, and the help. A removed flag gets the grammar's usage refusal at exit 2 (`internal/usage/parse.go:164-206`).

### The plan verb

`bench spec widen <slug> --ticket <basename> [--decision <reference>] <entry>...` is a leaf of the existing `bench spec` command. `bench spec retire` is its precedent for the primary-checkout refusal, the repository-root anchor, and a write that makes no commit. The verb edits only the one `Writes:` line of `specs/<slug>/tickets/<basename>` and writes the file through `durablefile`.

The approved scope comes from the default-branch tip. The spec-stage landing puts the signed-off tickets there, and build plan commits reach it only at the final landing. A caller-supplied base would let a build name its own widened commit, so the verb takes no base. `git.ResolvedDefault` owns the branch, and the verb reads the slug's ticket files at that tip through the git adapter.

The scope is the set of directory prefixes of the signed-off `Writes:` entries of every ticket of the slug. The prefix of entry E is E when the default tip holds E as a directory. Otherwise the prefix is the parent directory of E. A root-level file is its own prefix. `tickets.Covers` decides whether a prefix holds a path. The granularity is flagged for reviewer veto.

The verb grades each entry through `tickets.ParseWrite`, so an entry that the grammar refuses never reaches the line. It also refuses an entry that holds the `, ` separator and an entry that TOON cannot render. An entry that the ticket's own `Writes:` already covers reports `current` and writes nothing. Each other in-scope entry reports `added`, and each admitted out-of-scope entry reports `decided`. Any refusal writes nothing for the whole request.

The verb does not consult the implicit entries. The implicit-entry rule lives in the preflight package, and `internal/spec` cannot import that package, because the preflight package imports `internal/spec`. An implicit path needs no widening, so the verb grades it like any other path.

The plan commit stays with `bench commit`. A widened ticket moves the plan digest and the source digest, so a widening after a chunk freeze takes `bench record amendment` under the existing plan-expansion procedure.

### The derived fence

The ownership fence now names the union of the ticket `Writes:` entries plus the implicit entries. The implicit entries are `capture`, the spec folder, and the review pickup `reviews/<slug>.md`, which the review-record owner derives. `paths-authorized`, `bench preflight review`, and the landing's source authorization all read that one set.

The gather no longer refuses an empty fence section. An empty union authorizes only the implicit entries, so it stays fail-closed.

The `fence-writes` row goes, because its fence side no longer exists. The charge evidence `fence` block lists the ticket union in ticket and entry order without duplicates. Each block entry is the path that `tickets.WritesPath` projects, so the block drops a trailing `/` and the `(new)` marker, as the authorizing match does. The block keeps its name and field, and its description and the charge evidence format reference name the ticket union.

The staleness audit scope keeps the spec and its tickets folder. It replaces the fence-section entries with the `Writes:` entries of each ticket at the base and at the tip, read from committed objects. The audit already reads ticket objects at both identities for their `Covers:` tokens.

The coverage parser stops reading the section, and `checkReviewPickup` goes. `spec.FenceTokens` goes with its last reader. New code ignores a leftover section and refuses nothing about it.

The `source-not-fenced` route tells the operator to add the path to a ticket with `bench spec widen`. Its sentence keeps the term "ownership fence", which now names the derived fence.

### Bootstrap authority

ADR 0016 makes the installed broker the owner of every landing. The installation manifest authenticates that broker by path, version, and digest, so the candidate tree cannot replace it. The broker that lands this build therefore runs the pre-FT293 path check.

That check reads this spec's `## Ownership fences` section. It refuses an empty section, applies the old pickup rule, and loads the canary inventory and the anchor registries for the old closure. It also reads the binding registry when the kit seal verifies (`internal/preflight/ownership/closure.go:119-131`), but an absent registry gives an unbound string and no refusal. So this spec keeps its own section equal to its final ticket union plus `reviews/closure-removal.md`. The orchestrator updates that section in each plan expansion. A gap shows as the broker's `source-not-fenced` refusal. Spec retirement removes the section with the spec.

The build keeps the canary inventory and the anchor registry formats that the old closure loader parses. The prospective gate builds from the composed tree, so it runs the new code and ignores the section.

The landing names the broker install step, because this diff changes the broker's build inputs. `/bench-final-check` runs that step before any later spec-backed landing.

### No staged-spec sweep

The build edits no other staged spec. `bench commitment start` validates each bound deliverable against the sha256 identity of its spec bytes (`internal/commitment/repository/sources.go:99-113`). Five staged specs are bound deliverables in `.bench/commitment.json`, so a sweep would refuse each of those outcome starts until a new reviewer approval. A leftover section is inert text, because no new code reads it (CR72). Spec retirement removes each section with its spec. This call is flagged for reviewer veto.

### The guidance

`delegation-discipline.md` loses the refusal-boundary sentence and keeps the expectation sentence. That sentence names the widening route: the orchestrator widens a gate-found need with `bench spec widen` or asks the reviewer.

The slicing checks lose the `Enforced rules` section, because each rule there repeats a check that the ticket parser or build preflight owns. The checks also lose each slicing rule that runs a removed form or edits a spec fence. The `.bench/commitment.json` rule keeps its ticket half: exactly one ticket's `Writes:` line names that file, because only a ticket entry authorizes it under the derived fence. Ticket 05 owns that rewrite. The upper-bound sentence becomes an expectation sentence. `Fence-expansion evidence` stays, because FT293 does not name it.

The craft-spec template, its slicing paragraph, and its approval table drop the fence section. `/bench-write-spec` drops the fence-section sentences and the `--propose-writes` comparison. `/bench-implement-spec` drops the closure and `fence-writes` clauses, names the verb at the plan-expansion sentence, and states the scope rule at the out-of-`Writes:` pause. The operating guide's plan-expansion paragraph names the verb and the scope rule.

Each retired command sentence gets one forbid anchor. Each edited budgeted file holds or shrinks its line count.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| CR1 / 01-delete-proposal-and-breakdown-forms.md, 02-delete-ownership-closure.md | Preflight predicts no closure, checks no charge premises, and offers no proposal or breakdown form | CR01-CR22, CR80-CR85, CR89, CR91, CR93, CR95, CR96 | preflight decision, build, charge, and audit tests, the ticket grammar test, the help inventory, and the ticket-slicing anchors | no |
| CR2 / 03-add-spec-widen-verb.md | `bench spec widen` widens one ticket inside the signed-off scope | CR23-CR54, CR92 | the spec widen tests, the spec help test, and the help inventory | no |
| CR3 / 04-derive-fence-from-tickets.md, 05-retire-spec-fence-section.md | The ticket union is the fence, and no code reads a spec fence section | CR55-CR72, CR75-CR77, CR86, CR87, CR90 | landing, preflight, charge, audit, coverage, and refusal route tests, and the ticket-slicing anchors | yes |
| CR4 / 06-align-writes-guidance.md | The guidance states one `Writes:` rule beside the widening verb | CR78, CR79, CR88, CR94 | the guidance anchors and the budget and prose checks | no |

CR1 deletes before CR3 moves the fence, so CR3 changes no code that CR1 removes. CR2 is its own small chunk, because CR3's refusal route and CR4's guidance name its verb. Its review closes before CR3 starts. CR4 runs last, because it names the verb and the derived fence.

A ticket that removes a rendered shape owns every `rg` hit of the old text, under the slicing checks' rendered-shape rule. Ticket 01 owns the `--propose-writes` and `--breakdown` hits, the `Enforced rules` section that describes them, the known-issues row, and the unreleased changelog entry. Ticket 02 owns the closure row names, `*-closure`, and `charge-premises`. Ticket 03 owns the new help rows and the changelog line for the verb. Ticket 04 owns `fence-writes`, the route text, and the evidence block description. Ticket 05 owns the fence-section hits and the craft-spec template. Ticket 06 owns the delegate contradiction and the plan-expansion rule.

### Completion plan

The authored version 1 plan records the verification obligations. An approved build records its fresh authors in a version 2 amendment before dispatch. System verification uses the sealed binary that `bench test` supplies.

```bench-completion-plan
{"version":1,"chunks":[{"id":"CR1","tickets":["01-delete-proposal-and-breakdown-forms.md","02-delete-ownership-closure.md"],"verification":[{"id":"preflight-packages","command":"bench test --package ./internal/preflight --package ./internal/preflight/ownership --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd"},{"id":"tickets-and-help","command":"bench test --package ./internal/tickets --package ./cmd/bench"},{"id":"guidance-anchors","command":"bench test --check docs-currency-workflow"},{"id":"guidance-budgets","command":"bench test --check guidance-prose-budgets"},{"id":"author-clean-path-omission","command":"bench test --package ./internal/tickets","probe":"remove the clean-path refusal from tickets.ParseWrite; TestParseWriteRefusesAnUncleanPath must fail"}]},{"id":"CR2","tickets":["03-add-spec-widen-verb.md"],"verification":[{"id":"spec-and-help","command":"bench test --package ./internal/spec --package ./cmd/bench"},{"id":"axi-registry","command":"bench test --check axi-query-registry"},{"id":"command-routing","command":"bench test --check subcommand-routing"},{"id":"author-scope-omission","command":"bench test --package ./internal/spec","probe":"make the scope check admit every entry; TestSpecWidenRefusesAnOutOfScopeEntry must fail"}]},{"id":"CR3","tickets":["04-derive-fence-from-tickets.md","05-retire-spec-fence-section.md"],"verification":[{"id":"preflight-packages","command":"bench test --package ./internal/preflight --package ./internal/preflight/ownership --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd"},{"id":"fence-readers","command":"bench test --package ./internal/coverage --package ./internal/spec --package ./internal/chargeevidence --package ./internal/refusalroute"},{"id":"landing","command":"bench test --package ./internal/worktree"},{"id":"system","command":"bench test --check system"},{"id":"guidance-anchors","command":"bench test --check docs-currency-workflow"},{"id":"guidance-budgets","command":"bench test --check guidance-prose-budgets"},{"id":"author-pickup-omission","command":"bench test --package ./internal/worktree","probe":"drop the review pickup from the implicit entries; TestLandCommandAuthorizesTheReviewPickup must fail"}]},{"id":"CR4","tickets":["06-align-writes-guidance.md"],"verification":[{"id":"guidance-anchors","command":"bench test --check docs-currency-workflow"},{"id":"guidance-budgets","command":"bench test --check guidance-prose-budgets"},{"id":"prose","command":"bench test --check prose-mechanics"},{"id":"author-boundary-return","command":"bench test --check docs-currency-workflow","probe":"restore the refusal-boundary sentence in delegation-discipline.md; the forbid anchor must red"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/closure-removal/spec.md"},{"id":"preflight-packages","command":"bench test --package ./internal/preflight --package ./internal/preflight/ownership --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd"},{"id":"fence-readers","command":"bench test --package ./internal/coverage --package ./internal/spec --package ./internal/chargeevidence --package ./internal/refusalroute --package ./internal/tickets --package ./cmd/bench"},{"id":"landing","command":"bench test --package ./internal/worktree"},{"id":"system","command":"bench test --check system"},{"id":"guidance-anchors","command":"bench test --check docs-currency-workflow"},{"id":"guidance-budgets","command":"bench test --check guidance-prose-budgets"},{"id":"prose","command":"bench test --check prose-mechanics"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"coordinator-section-reader","command":"bench test --package ./internal/preflight","probe":"append the leftover fence-section tokens to the authorizing entries; TestCommandBuildIgnoresAFenceSection must fail"}]}
```

## Testing decisions

A good test drives the public seam and observes its output: the preflight verdict rows, the landing result, the charge evidence blocks, the audit tables, and the widen verb's exit code, table, and ticket bytes. A test never asserts that a deleted function is absent. A row that needs a symbol to be absent is review-owned.

The preflight rows keep their existing seams. `Decide` tests pin the row list, and the command tests drive `Command` against a real fixture repository from `preflighttest`. The build-grammar test keeps its seeded fixture pin and bound package, so a returned closure row reds it.

The widen verb's tests live beside the retire tests in `internal/spec`. They build real repositories, because the verb reads the default-branch tip and the primary-checkout state. The help expectation stays hand-written in the help inventory test as an independent expectation.

The landing rows run through the landing surface tests, which already seed ticket `Writes:` equal to the old fence. The audit rows run through the staleness tests and their committed fixture pairs. The guidance rows run through the anchor registry tests and the live `docs-currency-workflow` check.

### Seam diagram

    trigger: an orchestrator after a red gate or lane
        │
        ▼
    bench spec widen <slug> --ticket <t> [--decision <r>] <entry>...
        │
        ▼
    [ scope: directory prefixes of the default-tip tickets ] ──▶ one Writes: line, or a refusal
        ◀ tests attach here: real repository, exit code, widen table, ticket bytes

    trigger: bench worktree land --spec <slug>
        │
        ▼
    [ paths-authorized: ticket Writes: union + capture + spec folder + review pickup ] ──▶ publish or source-not-fenced
        ◀ tests attach here: landing surface fixture, refusal_paths table, published tree

    trigger: bench preflight build <slug>
        │
        ▼
    [ Decide: kept consistency rows only ] ──▶ verdict rows
        ◀ tests attach here: Decide row list and the build command grammar fixture

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CR01 | 1, 13 | The all-green build verdict lists no `fixture-closure`, `registry-closure`, or `anchor-closure` row | `internal/preflight/decision_test.go` (`TestDecideAllGreen`) | The exact row list reds a returned closure row |
| CR02 | 1, 13 | The build verdict renders no closure row for a ticket that writes a fixture-pinned path and a bound package | planned TestCommandBuildRendersFourGrammarRows in internal/preflight/command_build_test.go | The seeded pin and binding red a surviving closure row |
| CR03 | 13 | The build-mode verdict without a tickets folder lists its ticket-gated rows with no closure row | `internal/preflight/decision_test.go` (`TestTicketGatedRowsNotApplicableWithoutTickets`) | The exact ticket-gated list reds a returned closure row |
| CR04 | 2 | A build charge exits 0 for a ticket whose `Writes:` names a fixture-pinned path without its fixture | planned TestChargePreparesWithoutClosure in internal/preflight/charge_test.go | The current closure row refuses this charge |
| CR05 | 3 | A build charge exits 0 for a declared-complete inventory with an unproved current premise | planned TestChargePreparesWithUnprovedPremise in internal/preflight/charge_test.go | The current `charge-premises` row refuses this charge |
| CR06 | 3 | The build charge verdict holds no `charge-premises` row | planned TestChargePreparesWithUnprovedPremise in internal/preflight/charge_test.go | A kept premise row prints in the charge verdict |
| CR07 | 4 | `bench preflight build s --propose-writes --ticket a.md --base B --source-tip T` exits 2 with a usage refusal | planned TestPreflightBuildRefusesRemovedForms in internal/preflight/command_build_test.go | A kept proposal form exits 0 |
| CR08 | 5 | `bench preflight build s --breakdown` exits 2 with a usage refusal | planned TestPreflightBuildRefusesRemovedForms in internal/preflight/command_build_test.go | A kept breakdown form exits 0 |
| CR09 | 6 | `bench help` prints no `--propose-writes` row | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | The complete inventory expectation reds a kept help row |
| CR10 | 6 | `bench help` prints no `--breakdown` row | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | The complete inventory expectation reds a kept help row |
| CR11 | 7 | `--audit-spec` reports an open diagnostic for a partial planning inventory | `internal/preflight/staleness/completeness_test.go` (`TestAuditCompleteness`) | A deleted premise collector drops this diagnostic |
| CR12 | 7 | `--audit-spec` reports each pending native-future claim | `internal/preflight/staleness/completeness_test.go` (`TestAuditCompleteness`) | A deleted manual-claim collector drops these diagnostics |
| CR13 | 8 | `--audit-spec` reports no `kit-registry` diagnostic for a ticket that writes a bound package without a verified seal | planned TestAuditReportsNoClosureDiagnostic in internal/preflight/staleness/completeness_test.go | The current closure adds this diagnostic and keeps manual review required |
| CR14 | 9 | `bench preflight build` returns a verdict when the canary inventory does not parse | planned TestBuildVerdictReadsNoClosureInventory in internal/preflight/command_build_test.go | The current gather refuses with `fixture inventory refused` |
| CR15 | 10 | `internal/tickets` declares no `BindingRules` or `BoundFiles` function | review-owned: the Standards axis confirms the deletion of the binding table and its tests | A surviving registry compiles silently with no reader |
| CR16 | 11 | `writes-resolve` reds a `Writes:` entry that names no tree path | planned TestCommandBuildRendersFourGrammarRows in internal/preflight/command_build_test.go | A deleted shared tree probe drops this red |
| CR17 | 11 | `kit-pin` reds a system-tagged test file that a ticket writes without `BENCH_KIT` | planned TestCommandBuildRendersFourGrammarRows in internal/preflight/command_build_test.go | A deleted system-tag probe drops this red |
| CR18 | 11 | `rows-owned` reds a declared row that no ticket cites | `internal/preflight/decision_test.go` (`TestDecideRowsOwned`) | A deleted `UncitedIDs` drops this red |
| CR19 | 12 | `tickets.ParseWrite` refuses an entry with a `..` segment | planned TestParseWriteRefusesAnUncleanPath in internal/tickets/writes_test.go | The closure's clean-path refusal leaves with the closure |
| CR20 | 12 | `tickets.ParseWrite` refuses an absolute entry | planned TestParseWriteRefusesAnUncleanPath in internal/tickets/writes_test.go | A rule that checks only `..` admits an absolute path |
| CR21 | 12 | `tickets-parse` reds a ticket whose `Writes:` holds `internal/x/../../cmd/y.go` | planned TestCommandBuildRefusesAnUncleanWritesEntry in internal/preflight/command_build_test.go | A grammar rule that the ticket parser skips leaves the entry green |
| CR22 | 12 | `bench preflight build` exits 1 for a ticket whose `Writes:` names `../outside.txt` while that file exists outside the repository | planned TestCommandBuildRefusesAnUncleanWritesEntry in internal/preflight/command_build_test.go | Without the grammar rule the tree probe finds the outside file and every row stays green |
| CR23 | 14 | `bench spec widen s --ticket a.md internal/x/y_test.go` appends `, internal/x/y_test.go` to the `Writes:` value when a signed-off entry names `internal/x/y.go` | planned TestSpecWidenAddsAnInScopeEntry in internal/spec/widen_test.go | A verb that writes nothing leaves the line unchanged |
| CR24 | 14 | The in-scope widening prints one `widen[1]{ticket,entry,state}` row with state `added` at exit 0 | planned TestSpecWidenAddsAnInScopeEntry in internal/spec/widen_test.go | A silent write gives the caller no evidence |
| CR25 | 25 | The in-scope widening leaves `HEAD` unchanged | planned TestSpecWidenAddsAnInScopeEntry in internal/spec/widen_test.go | A committing verb moves `HEAD` |
| CR26 | 15 | An entry outside every directory prefix exits 1 with unchanged ticket bytes | planned TestSpecWidenRefusesAnOutOfScopeEntry in internal/spec/widen_test.go | A verb with no scope check writes the entry |
| CR27 | 15 | The out-of-scope refusal names the refused entry | planned TestSpecWidenRefusesAnOutOfScopeEntry in internal/spec/widen_test.go | A generic refusal hides which path needs a decision |
| CR28 | 15 | The out-of-scope refusal names the `--decision <reference>` rerun | planned TestSpecWidenRefusesAnOutOfScopeEntry in internal/spec/widen_test.go | A refusal without the route strands the caller |
| CR29 | 16 | `--decision r1` admits an out-of-scope entry and prints state `decided` | planned TestSpecWidenAdmitsADecidedEntry in internal/spec/widen_test.go | A verb that ignores the flag keeps refusing |
| CR30 | 16 | `--decision` with an empty value exits 2 | planned TestSpecWidenUsage in internal/spec/widen_test.go | An empty reference admits an unreviewed path |
| CR31 | 17 | A widening under a directory that only a worktree ticket names exits 1 | planned TestSpecWidenReadsTheSignedOffScope in internal/spec/widen_test.go | A worktree-read scope admits a path under a self-widened directory |
| CR32 | 17 | A widening into a new worktree-only ticket inside the signed-off scope exits 0 | planned TestSpecWidenReadsTheSignedOffScope in internal/spec/widen_test.go | A per-ticket scope refuses a ticket that the build added |
| CR33 | 18 | A signed-off file entry `internal/x/y.go` admits `internal/x/sub/z.go` | planned TestSpecWidenDirectoryPrefix in internal/spec/widen_test.go | A file-only rule refuses this sibling subtree |
| CR34 | 18 | A signed-off directory entry `internal/status` refuses `internal/other/x.go` | planned TestSpecWidenDirectoryPrefix in internal/spec/widen_test.go | A parent-of-directory rule admits this path |
| CR35 | 18 | A signed-off root file entry `CHANGELOG.md` refuses `README.md` | planned TestSpecWidenDirectoryPrefix in internal/spec/widen_test.go | A root parent rule admits every path |
| CR36 | 18 | A signed-off `(new)` entry `internal/n/new.go` that the default tip lacks admits `internal/n/other.go` | planned TestSpecWidenDirectoryPrefix in internal/spec/widen_test.go | A tree-only rule drops the absent entry's directory |
| CR37 | 18 | A signed-off entry `internal/git` refuses `internal/git2/x.go` | planned TestSpecWidenDirectoryPrefix in internal/spec/widen_test.go | A bare string prefix admits a sibling directory |
| CR38 | 19 | An entry that the ticket's `Writes:` already covers exits 0 with state `current` and unchanged bytes | planned TestSpecWidenReportsACoveredEntry in internal/spec/widen_test.go | A verb that appends again duplicates the entry |
| CR39 | 20 | A request with one in-scope entry and one out-of-scope entry exits 1 with unchanged ticket bytes | planned TestSpecWidenWritesAllOrNothing in internal/spec/widen_test.go | A per-entry writer lands the in-scope half |
| CR40 | 21 | The widened ticket keeps every byte outside the `Writes:` value | planned TestSpecWidenKeepsOtherBytes in internal/spec/widen_test.go | A re-rendering writer changes other lines |
| CR41 | 21 | A ticket with CRLF line endings keeps the CRLF ending of its widened line | planned TestSpecWidenKeepsOtherBytes in internal/spec/widen_test.go | A split on LF writes the carriage return into the new entry |
| CR42 | 21 | A ticket whose last line lacks a newline keeps that state | planned TestSpecWidenKeepsOtherBytes in internal/spec/widen_test.go | A writer that appends a newline changes the file end |
| CR43 | 22 | The verb exits 1 on the primary checkout with the shared refusal as its prefix | planned TestSpecWidenRefusesThePrimaryCheckout in internal/spec/widen_test.go | A verb without the check writes the primary checkout |
| CR44 | 23 | An entry with a `..` segment exits 1 | planned TestSpecWidenRefusesHostileEntries in internal/spec/widen_test.go | A lexical prefix admits `internal/x/../../cmd/y.go` |
| CR45 | 23 | An entry that holds `, ` exits 1 | planned TestSpecWidenRefusesHostileEntries in internal/spec/widen_test.go | The one split reads the entry as two entries |
| CR46 | 23 | An entry with an ESC byte exits 1 | planned TestSpecWidenRefusesHostileEntries in internal/spec/widen_test.go | The field then holds a byte that TOON refuses |
| CR47 | 23 | An in-scope entry with a space and a bracket is added with exactly its bytes | planned TestSpecWidenRefusesHostileEntries in internal/spec/widen_test.go | An over-strict filter refuses a legal path |
| CR48 | 23 | A ticket file that is a symbolic link exits 1 | planned TestSpecWidenRefusesHostileTickets in internal/spec/widen_test.go | Following the link writes outside the ticket folder |
| CR49 | 23 | A ticket path that is a FIFO exits 1 without a read | planned TestSpecWidenRefusesHostileTickets in internal/spec/widen_test.go | A plain read blocks on the special file |
| CR50 | 23 | A ticket with no `Writes:` line exits 1 | planned TestSpecWidenRefusesHostileTickets in internal/spec/widen_test.go | A writer that adds a line invents a field |
| CR51 | 24 | `bench help` prints the `bench spec widen` row | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | The complete inventory reds a missing row |
| CR52 | 24 | `bench spec --help` lists the `widen` subcommand | `internal/spec/spec_test.go` (`TestSpecHelpPrintsThePublicSubcommandInventory`) | The exact help text reds a missing subcommand |
| CR53 | 26 | A default-branch tip that holds no ticket of the slug exits 1 | planned TestSpecWidenRefusesWithoutSignedOffTickets in internal/spec/widen_test.go | A worktree fallback admits a scope that nobody signed off |
| CR54 | 27 | The verb run from `internal/x` below the root widens `specs/s/tickets/a.md` | planned TestSpecWidenFromADeeperDirectory in internal/spec/widen_test.go | A working-directory-relative path misses the ticket |
| CR55 | 28 | The landing publishes a changed path that a ticket's `Writes:` names and no fence section lists | planned TestLandCommandAuthorizesFromTicketWrites in internal/worktree/land_surface_test.go | A fence-section reader refuses this path |
| CR56 | 28 | The landing refuses a changed path that no ticket `Writes:` entry covers | `internal/worktree/land_surface_test.go` (`TestLandCommandFenceRefusalNamesThePath`) | A widened authority publishes the stray path |
| CR57 | 28 | `bench preflight build` reds `paths-authorized` for a change outside every ticket `Writes:` entry | `internal/preflight/command_build_test.go` (`TestCommandBuildOutOfFenceRed`) | A widened authority greens the out-of-scope change |
| CR58 | 29 | The landing publishes `reviews/<slug>.md` that no ticket names | planned TestLandCommandAuthorizesTheReviewPickup in internal/worktree/land_surface_test.go | Without the implicit pickup the landing refuses the review record |
| CR59 | 28 | The landing publishes a `capture/` path that no ticket names | `internal/worktree/land_surface_test.go` (`TestLandCommandAuthorizesCaptureOutsideTheFence`) | A union that drops the implicit entries refuses the capture path |
| CR60 | 30 | The all-green build verdict lists no `fence-writes` row | `internal/preflight/decision_test.go` (`TestDecideAllGreen`) | The exact row list reds a returned row |
| CR61 | 31 | `bench preflight build` on a spec with no fence section returns a verdict | planned TestCommandSpecWithoutFenceSection in internal/preflight/command_bootstrap_test.go | The current gather refuses with `ownership fences empty` |
| CR62 | 31 | With an empty tickets folder, `paths-authorized` reds a changed production path | planned TestCommandSpecWithoutFenceSection in internal/preflight/command_bootstrap_test.go | An empty union read as unrestricted greens the path |
| CR63 | 32 | The build charge `fence` block lists each `Writes:` entry of every ticket once | planned TestChargeFenceListsTheTicketUnion in internal/preflight/charge_test.go | A fence-section reader prints an empty block |
| CR64 | 32 | The review charge `fence` block lists each `Writes:` entry of every ticket once | planned TestReviewChargeFenceListsTheTicketUnion in internal/preflight/review_charge_test.go | A review path left on the section prints an empty block |
| CR65 | 32 | The charge evidence format reference states the `fence` block description that the schema holds | review-owned: the Standards axis compares the reference row with the schema row | A schema edit without the reference leaves the old description |
| CR66 | 33 | `--audit-spec` reports drift for a path that a base ticket's `Writes:` names | planned TestAuditScopeReadsTicketWrites in internal/preflight/staleness/drift_test.go | A fence-section scope omits this path |
| CR67 | 33 | `--audit-spec` reports drift for a path that only a tip ticket's `Writes:` names | planned TestAuditScopeReadsTicketWrites in internal/preflight/staleness/drift_test.go | A base-only reader omits this path |
| CR68 | 34 | The `source-not-fenced` route names `bench spec widen` | planned TestSourceNotFencedRouteNamesWiden in internal/refusalroute/route_test.go | A stale route names the retired fence section |
| CR69 | 34 | The `source-not-fenced` route holds no `## Ownership fences` text | planned TestSourceNotFencedRouteNamesWiden in internal/refusalroute/route_test.go | A partial edit keeps the old instruction beside the new one |
| CR70 | 35 | `bench coverage --check` exits 0 for a spec whose fence section omits the review pickup | planned TestCommandCheckIgnoresAFenceSection in internal/coverage/citations_test.go | The current pickup rule exits 1 |
| CR71 | 36 | `internal/spec` declares no `FenceTokens` function | review-owned: the Standards axis confirms the deletion of the parser and its tests | A surviving parser invites a new reader |
| CR72 | 37 | `paths-authorized` reds a changed path that only a leftover fence section lists | planned TestCommandBuildIgnoresAFenceSection in internal/preflight/command_bootstrap_test.go | A surviving section reader authorizes the listed path |
| CR75 | 39 | This spec declares a fence section equal to its ticket `Writes:` union plus `reviews/closure-removal.md` at the landing | review-owned: the pre-FT293 broker refuses a gap as `source-not-fenced` | A missing path stops the landing |
| CR76 | 39 | The build keeps the canary inventory and anchor registry formats that the pre-FT293 closure loader parses | review-owned: the old broker refuses an unparsable inventory as `fixture inventory refused` or `anchor registry refused` | A format change stops the landing |
| CR77 | 31 | `bench preflight review` reds `paths-authorized` for a change outside every ticket `Writes:` entry | planned TestCommandReviewAuthorizesFromTicketWrites in internal/preflight/command_review_test.go | A review path left on the section authorizes nothing or everything |
| CR78 | 40 | `delegation-discipline.md` holds no `is a refusal boundary` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored sentence reds the forbid anchor |
| CR79 | 40 | `delegation-discipline.md` names `bench spec widen` beside the expectation rule | review-owned: the Spec axis reads the sentence against story 40 | Review owns rule wording |
| CR80 | 41 | `slicing-checks.md` holds no `## Enforced rules` heading | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored section reds the forbid anchor |
| CR81 | 42 | `bench-write-spec.md` holds no `--propose-writes` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored comparison step reds the forbid anchor |
| CR82 | 42 | `slicing-checks.md` holds no `--propose-writes` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored proposal step reds the forbid anchor |
| CR83 | 42 | `slicing-checks.md` holds no `--breakdown` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored breakdown step reds the forbid anchor |
| CR84 | 42 | `map-only-route.md` holds no `--propose-writes` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored proposal step reds the forbid anchor |
| CR85 | 42 | `map-discipline.md` holds no `breakdown, proposal` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored slicing step reds the forbid anchor |
| CR86 | 43 | The craft-spec skill holds no `## Ownership fences` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored template section reds the forbid anchor |
| CR87 | 43 | `bench-write-spec.md` holds the needle "After the ticket review accepts, `bench preflight build <slug> --plan-only` grades the whole ticket graph." in its Process section | `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A sentence that keeps the fence step fails the exact needle |
| CR88 | 44 | The operating guide names `bench spec widen` in its plan-expansion paragraph | review-owned: the Spec axis reads the paragraph against story 44 | Review owns rule wording, so no require anchor pins it |
| CR89 | 45 | `bench-implement-spec.md` holds no `*-closure` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored closure route reds the forbid anchor |
| CR90 | 45 | `bench-implement-spec.md` holds no `fence-writes` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored fence route reds the forbid anchor |
| CR91 | 47 | `docs/known-issues.md` holds no `writes proposal` text | planned forbid row in `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored row reds the forbid anchor |
| CR92 | 48 | The `[Unreleased]` section of `CHANGELOG.md` names `bench spec widen` | review-owned: the Spec axis reads the entry | Review owns changelog prose |
| CR93 | 48 | The `[Unreleased]` section holds no entry that adds `--breakdown` | review-owned: the Spec axis reads the entry | Review owns changelog prose |
| CR94 | 49 | Each budgeted guidance file that the build edits holds a line count at or below its budget row | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | A net addition to a full file exceeds its limit and reds the live-tree budget test |
| CR95 | 46 | The anchor registry holds no require row for the `--propose-writes` sentence | `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A kept require row reds the live check after the sentence goes |
| CR96 | 52 | The `ticket-grammar` check reds a ticket whose `Blocked by:` names an absent sibling | `internal/conformance/ticket_grammar_test.go` (`TestTicketGrammarSweepRedsDanglingBlocker`) | A deletion of the whole check with its binding grader drops this red |

Not covered: story 50 — reviewed exclusion; the `Fence-expansion evidence` section and the paired slicing readers stay unchanged, and review grades that no ticket touches them.
Not covered: story 51 — reviewed exclusion; FT380 owns plan-commit grading, and no ticket here changes `bench commit --preflight-build`.
Not covered: story 38 — reviewed exclusion; no ticket `Writes:` line names another staged spec, and review grades that the build edits none.

### Edge inventory

The profile's hostile-input checklist applies to the widen verb's entries, its ticket file, its working directory, and the default-branch read. CR44 to CR50 and CR53 to CR54 cover the reached classes. The verb renders its table through `toon.Table`, so a numeric-looking entry takes the encoder's quoting. The derived fence reads only grammar-validated entries, and an unclean entry covers no clean git path.

The kit audience is every repository that links the kit. A linked repository with no `specs/` folder never calls the verb. An absent tickets folder at the default tip and an empty one take the same refusal (CR53), because neither holds a signed-off entry.

- **Won't handle**: an implicitly authorized path in `bench spec widen` — the landing authorizes `capture`, the spec folder, and the review pickup without an entry, and `implicitEntries` in the preflight package stays their one owner.
- **Won't handle**: a light-path tickets-only folder — the default tip holds no signed-off ticket, so CR53 refuses, and the tickets-only landing skips path authorization in `landingSourceRange`.
- **Won't handle**: a hand edit of a `Writes:` line that bypasses the verb — today's plan commit has the same trust, review grades each ticket diff, and CR56 still refuses a path outside the union.
- **Won't handle**: a durable record of each `--decision` reference — the plan-expansion rule's `bench learning` entry keeps that trail.
- **Won't handle**: a recovery-matrix face for the widen refusals — `bench spec retire` has no face either, and `refusalroute` keeps its write-verb set.
- **Won't handle**: a scope that grows after a decided entry — each out-of-scope path takes its own decision until the landing moves the signed-off tickets.
- **Won't handle**: two widen calls at once in one worktree — one worktree has one writer, and `durablefile` makes each write atomic.
- **Won't handle**: a symbolic default-branch ref — `git.ResolvedDefault` owns that resolution for every reader.
- **Won't handle**: non-ASCII whitespace inside an entry — `tickets.ParseWrite` owns the entry trim for every `Writes:` reader.
- **Won't handle**: the bench-writer agent definition keeps "Edit only the paths inside your fence." — its fence is the current `Writes:` line, which the orchestrator widens.
- **Won't handle**: `README.md` keeps "Checks the approved artifacts, dependencies, and write fence before implementation." — the derived fence is still a write fence.
- **Won't handle**: `bench-implement-spec.md` keeps "Verify the fence, dependency completion, and expected source tip before the ticket's fresh author starts." — an anchor and a canary pin it, and its fence now names the ticket union.
- **Won't handle**: the repair-charge template keeps "Ownership fence: <exact repo-relative file or path prefix>" — an anchor pins the label, and the field holds the ticket's `Writes:` line.
- **Won't handle**: older `[Unreleased]` and released changelog entries that name `--propose-writes` — they record the changes of their own entries.
- **Won't handle**: `docs/research` and `docs/audits` mentions of the closure — they are dated evidence.
- **Won't handle**: a leftover `## Ownership fences` section, and prose about a fence or the proposal, in another staged spec — CR72 keeps the section inert, and each build-entry staleness pass corrects its own claims.

## Ownership fences

- `reviews/closure-removal.md`

## Out of scope

- A landing check that compares each `Writes:` entry with the signed-off scope: 6 edits, 2 gate runs.
- A durable decision trail for `--decision` in the review record: 7 edits, 2 gate runs.
- A package-prefix granularity option for the verb: 3 edits, 1 gate run.
- FT380 plan-commit grading through `bench commit --preflight-build`: owned by the FT380 outcome.
- FT375-residual staleness of `writes-resolve` reds from landed deletions: owned by the FT375-residual outcome, and `writes-resolve` stays for it.

These are planning estimates, not approved priorities.

## Further notes

Evidence status: source-backed proposal from eight paired read-only research returns at 77f4fb2ff and targeted reads at 658ba0dbb, with no runtime red claimed during specification.

### Flagged for reviewer veto

1. Premise checks: the build removes only the `charge-premises` row. The premise collectors stay for `--audit-spec`. The alternative deletes the planning-facts manifest and changes the shipped audit.
2. No staged-spec sweep: the build leaves each other staged spec's fence section as inert text, because a sweep changes five bound deliverable identities that `bench commitment start` validates. The alternative sweeps the sections and re-approves those five bindings.
3. Scope granularity: the scope is the directory prefix of each signed-off entry. The FT293 replay summary gives 78 percent of 91 widenings for a directory prefix. The tree holds no replay record and no definition of a package prefix.

### Flagged additions

- The clean-path refusal moves from the closure into `tickets.ParseWrite` (CR19 to CR22). It preserves an existing guarantee that the closure deletion would otherwise remove.
- Only the audit kind gathers the ownership snapshot. No row grades this call, because no remaining verdict reads the snapshot.
- The forbid anchors CR78, CR80 to CR87, and CR89 to CR91 follow the slicing checks' retirement rule for each sentence that grants a removed form.
- The review pickup joins `implicitEntries` (CR58). The derived fence needs it, because no ticket names the pickup. The entry is the exact path from `reviewrecord.RecordPath`, never the `reviews` folder; the build's mutation probe widens the entry to the folder and expects `reviews/<other>.md` to red `paths-authorized`.

### Source-sentence-to-row table

| source | clause | rows |
|---|---|---|
| `roadmap/FT293.md:3-6` | `Writes:` is an expectation, and a gate-found need widens it inside the approved scope | CR23-CR39, CR78-CR79 |
| `roadmap/FT293.md:8-10` | predictive closure cannot be complete, and the copy breaks one source per fact | CR01-CR04, CR13-CR15 |
| `roadmap/FT293.md:13-14` | delete the closure, `--propose-writes`, `--breakdown`, the premise checks, and the guidance that repeats a check | CR01-CR10, CR80-CR86, CR95 |
| `roadmap/FT293.md:15` | keep the consistency checks and the landing check | CR16-CR18, CR56-CR57, CR77 |
| `roadmap/FT293.md:16` | derive the spec scope from the tickets | CR55, CR58-CR64, CR70-CR72, CR75 |
| `roadmap/FT293.md:17-18` | one plan verb adds a gate-found path, and an out-of-scope path goes to the reviewer | CR23-CR54 |
| `roadmap/FT293.md:19-20` | remove the delegate-guidance contradiction | CR78-CR79 |
| `roadmap/FT293.md:22-26` | the spec stage selects the scope granularity from the replay | CR33-CR37 |
| `decisions/architecture-planning/tickets/1.md:12-16` | preflight checks no charge premises | CR05-CR06 |
| `decisions/architecture-planning/tickets/1.md:18` | unknown facts stay unknown or refused, never guessed green | CR53, CR61-CR62 |

Each `Occurrence:` line of `roadmap/FT293.md:31-41` names a closure gap, a proposal cost, or a fence edit. CR01 to CR10 remove the closure and proposal faces, and CR23 to CR39 give each gate-found widening one verb.

### Pre-review proof checklist

- `Cited symbols`: `Decide` (`internal/preflight/decision.go:218`); `closureRows` (`internal/preflight/closure.go:57-74`); `chargePremisesCheck` (`decision.go:378-384`); `ownership.Collect` (`internal/preflight/ownership/ownership.go:232`); `collectOwnership` (`closure.go:116`); `pathsAuthorizedCheck`, `unauthorizedPaths`, and `authorizingEntries` (`internal/preflight/paths_authorized.go:10`, `:24`, `:45`); `implicitEntries`, `fenceWritesCheck`, and `writesResolveCheck` (`internal/preflight/fence_writes.go:72`, `:21`, `:87`); `spec.FenceTokens` (`internal/spec/fences.go:31`); `checkReviewPickup` (`internal/coverage/citations.go:341`); `AuthorizeReviewedSource` (`internal/preflight/gather.go:210`); `landingSourceRange` (`internal/worktree/land_identity.go:244`); `spec.Command` and `retireCommand` (`internal/spec/spec.go:181`, `:236`); `git.ResolvedDefault` (`internal/git/default_branch.go:22`); `git.ListTree` and `git.ReadTreeFile` (`internal/git/tree.go:154`, `:197`); `tickets.ParseWrite` and `tickets.Covers` (`internal/tickets/writes.go:35`, `:76`); `BindingRules` and `BoundFiles` (`internal/tickets/registry_data.go:154`, `:173`); `repositoryPath` (`internal/preflight/ownership/closure.go:85`); `treeHolds` (`closure.go:188`); `usage.PrimaryCheckoutRefusal` (`internal/usage/worktree.go:64`).
- `Import edges`: the verb adds `internal/spec` to `internal/tickets` and `internal/spec` to `internal/durablefile`. `go list -deps` of both targets holds no `internal/spec`, so neither edge makes a cycle. `internal/spec` cannot import `internal/preflight`, because `internal/preflight` imports `internal/spec`.
- `Source-row clauses and occurrences`: the source table above quotes each clause, and `roadmap/FT293.md:31-41` holds the eleven occurrences.
- `Promised field labels`: `widen[N]{ticket,entry,state}`, with the states `added`, `decided`, and `current`.
- `Changed-function callers`: `ownership.Collect` has one production caller at `closure.go:117`; `OwnedPaths`, `Owned`, and `Unordered` have only deleted callers; `UncitedIDs` keeps `ticket_rows.go:44`; `SubjectPaths` keeps `staleness/collect.go:113`; `ParseManifest` keeps `staleness/collect.go:31-32`; `spec.FenceTokens` has `coverage/rows.go:137`, `gather.go:54`, and `staleness/collect.go:79`; `AuthorizeReviewedSource` has `land_identity.go:254`; `BoundFiles` has no production caller; `spec.Command` has `cmd/bench/main.go:133`.
- `Copy survival`: CR15 and CR71 are review-owned absences. CR72 reds a surviving section reader.
- `Rendered-shape readers`: the help rows at `cmd/bench/help_inventory_test.go:112-113`; the spec help at `internal/spec/spec_test.go:102-124`; the closure row names at `decision_test.go:59`, `:507-538`, `:580`, `command_build_test.go:173-174`, `anchor_closure_test.go`, `ownership/authority_test.go:195-218`, and `proposal_test.go:99`; `fence-writes` at `decision.go:231`, `decision_test.go:59`, `:580`, `command_build_test.go:42`, `fence_writes_test.go`, `charge_test.go:320`, `ownership/authority_test.go:218`, `bench-implement-spec.md:19`, `map-only-route.md:24`, and `roadmap/FT375.md:10`; `charge-premises` at `ownership/premises_test.go:341`, `:369`, `:384`, and `:394`; the route text at `internal/refusalroute/faces_land.go:66`; the `fence` block description at `internal/chargeevidence/schema.go:232` and `charge-evidence-format.md:66`; `ownership fences empty` at `gather.go:56` and `command_bootstrap_test.go:207-263`.
- `Pin operators`: `TestDecideAllGreen` compares the row count and each name with `!=`; `TestHelpInventoryIsComplete` compares a hand-written inventory exactly; the budget check fails when a count exceeds its limit (`internal/conformance/prose_budget_test.go:99-100`).
- `Entry reads`: the verb reads the default branch through `git.ResolvedDefault` and the checkout identity through `git.IsPrimaryCheckout`, and its tests drive real repositories as `TestRetireOnPrimaryCheckoutRefusesAndDeletesNothing` does.
- `Derived expectations`: the widen table expectation derives through `toon.Table`; the help expectation stays an independent hand-written list, and its red comes from a deleted registry row.
- `Consolidated rules`: one rule moves. The old site is the closure's clean-path refusal (`ownership/closure.go:53-55` through `repositoryPath`); the new site is `tickets.ParseWrite`. CR19 to CR22 cover the changed cell.
- `Quantified obligations`: "each `Writes:` entry of every ticket" is CR63 and CR64; "each budgeted guidance file" is CR94; "every other staged spec" is an exclusion (story 38).
- `Workflow-step writes`: the verb writes one ticket file. The plan commit then moves the plan digest (`internal/reviewrecord/plan.go:84-138`) and the source digest (`plan.go:201-210`), so a widening after a chunk freeze takes `bench record amendment`. The verb writes nothing under `reviews/`.

### Preflight row-list pins

These files assert the preflight row list or a row name that this build removes: `internal/preflight/decision_test.go` (`TestDecideAllGreen` at :56-61 and `TestTicketGatedRowsNotApplicableWithoutTickets` at :572-576, and the cases at :507-538), `internal/preflight/command_build_test.go` (`TestCommandBuildRendersSixGrammarRows` at :162-190 and the row at :42), `internal/preflight/fence_writes_test.go`, `internal/preflight/anchor_closure_test.go`, `internal/preflight/charge_cells_test.go:55-59`, `internal/preflight/charge_test.go:320`, `internal/preflight/ownership/authority_test.go:195-218`, `internal/preflight/ownership/premises_test.go:317-394`, and `internal/preflight/staleness/drift_test.go:310`.

### Posture-change fixtures

The derived fence reds a fixture whose tickets do not name the paths that its old fence section authorized. These fixture writers carry a fence section today: `internal/preflight/preflighttest/fixture.go:139-150`, `internal/worktree/land_fixtures_test.go:125`, `internal/worktree/delegated_integration_test.go:28`, `internal/systemtest/owner_landing_fixture_test.go:57`, `internal/systemtest/charge_evidence_test.go:61`, `cmd/bench/preflight_version_test.go:25`, `internal/reviewrecord/recordtest/fixture.go:112-113`, and `internal/preflight/ownership/ownershiptest/repo.go:99-116`. `land_fixtures_test.go:126-132` already sets the ticket `Writes:` equal to its fence less the pickup. The slicing research confirms each other writer.

### Reader inventory

- The landing chain: `internal/worktree/land_identity.go:244-266` calls `preflight.AuthorizeReviewedSource` (`gather.go:210-226`), which runs `pathsAuthorizedCheck` and `unauthorizedPaths` (`paths_authorized.go:10-33`).
- The fence parser and its three production callers: `internal/spec/fences.go:31`, `internal/coverage/rows.go:137`, `internal/preflight/gather.go:54`, and `internal/preflight/staleness/collect.go:79`.
- The charge evidence fence: `internal/preflight/charge_pack.go:68`, `internal/preflight/review.go:260`, `internal/chargeevidence/metadata.go:22`, `internal/chargeevidence/schema.go:232`, and `.agents/skills/bench-craft-delegate/references/charge-evidence-format.md:66`.
- The audit snapshot readers: `internal/preflight/staleness/completeness.go:18-35` and `internal/preflight/staleness/collect.go:23-60`.
- The guidance readers: `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:36-39`, `.agents/skills/bench-craft-tickets/references/slicing-checks.md:10-35`, `:48-58`, `:60`, `:70`, and `:85`, `.agents/skills/bench-craft-tickets/SKILL.md:63`, `.agents/skills/bench-craft-spec/SKILL.md:65`, `:71`, `:145-147`, and `:155`, `.agents/skills/bench-craft-spec/references/map-discipline.md:143` and `:168-169`, `.agents/skills/bench-craft-spec/references/map-only-route.md:21-26`, `:41`, and `:48`, `.agents/commands/bench-write-spec.md:67`, `:73`, and `:79`, `.agents/commands/bench-implement-spec.md:19`, `:54`, and `:87`, `.agents/skills/bench-craft-review/SKILL.md:42`, `.agents/commands/bench-review-implementation.md:102`, `.bench/BENCH.md:124`, and `docs/known-issues.md` row 17.
- The anchors: `internal/anchors/registry_ticket_passes.go:8` and `:10` with `registry_ticket_passes_test.go:14` and `:16`. `bench anchors` shows no anchor on `map-only-route.md`, two retained forbids on `docs/known-issues.md`, and an operating-guide forbid on an older plan-expansion sentence that the edit must not restore.
- The non-Go readers: no `.mjs` script, workflow, or `package.json` entry reads these features.

### Budget headroom

At 658ba0dbb, `bench-write-spec.md` holds 87 of 87 lines, `bench-implement-spec.md` holds 93 of 93, the craft-spec skill holds 155 of 155, and the craft-review skill holds 122 of 122. The craft-tickets skill holds 97 of 100. The operating guide gained one paragraph from FT396 after 658ba0dbb, so the slicer reads its count after the fold.

### Validation plan

Each ticket author observes each independent new expectation red before the production edit, or through a named `bench probe`. The author probes in CR1 remove the clean-path refusal, in CR2 admit every entry, in CR3 drop the implicit pickup, and in CR4 restore the refusal-boundary sentence. The coordinator probe appends the leftover section tokens to the authorizing entries. A CR3 review probe widens the pickup entry to the `reviews` folder and expects `reviews/<other>.md` to red `paths-authorized`. Each probe must report a restored subject.

The final landing runs under the pre-FT293 broker. Before the landing, the orchestrator compares this spec's fence section with its ticket union plus the pickup. After the landing, `/bench-final-check` runs the broker install step that the landing names.
