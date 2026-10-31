# Build-entry staleness from landed work

Status: staged

Roadmap: FT375

Decision source: named reviewed artifact — the bound criteria `FT375-residual.map-sources` and `FT375-residual.writes-resolve` in `.bench/commitment.json`, read with `roadmap/FT375.md` at its pinned identity `sha256:837124cb6340ce3d9ae84b287b9ac023cd055677b5622382b1ec3ccd4fcae7b6`, plus the two reviewer decisions of 2026-10-31 in the `/bench-write-spec FT375` conversation (D1 map scope, D2 row split).

Verification log: 0 iteration(s) to accept — pending

## Problem

A build can start from a plan that landed work made stale. Two such cases stop a build today, and the build preflight finds neither case at build entry.

The first case is a decision-map source that the delivery closure deletes.
The final landing of a spec deletes the `roadmap/FT<n>.md` file of each row that the delivery closes (`internal/commitment/repository/closure.go:59`, `Store.Closure`).
A decision map whose `## Sources` section cites that file then fails `decision-map-integrity` on the published tree (`internal/maps/validation.go:249`, `internal/conformance/registry/checks.go:36`).
The FT217 build found this case at its complete checkpoint, and the repoint forced a refreeze, a re-record, and a confirming round (`d656f94ad`).
The FT293 build found it in the top-level map `decisions/architecture-planning.md` (`1b79d041f`).
That map still cites `roadmap/FT375.md` (`decisions/architecture-planning.md:45`), so the closure of this spec repeats the case.
The build preflight reads no decision map today: no production file under `internal/preflight` reads `decisions/` or a `## Sources` section.

The second case is a `Writes:` entry whose path landed work deleted.
The build entry runs a bare `bench preflight build <slug>` and stops on any red row except `completion-plan` (`.agents/commands/bench-implement-spec.md:15`).
At a fresh build entry, the bare form has no explicit range and no recorded chunk, so `DeletedPaths` and `BuildDeletedPaths` are empty (`internal/preflight/gather.go:69-78`, `internal/preflight/tree_probe.go:50`).
`writes-resolve` then reds a path that a landed commit deleted with the same text as a typing error (`internal/preflight/fence_writes.go:71`).
On 2026-10-17, ticket 05 of the first FT375 build named 31 canaries that anchor-admission-prune deleted, and the orchestrator took that red against the written rule (`a0f45a844`).

FT293-closure-removal is delivered, and `writes-resolve` survived it (`internal/preflight/fence_writes.go:50`, `internal/preflight/decision.go:216`).
The condition of the `FT375-residual.writes-resolve` criterion therefore holds.

## Solution

The build preflight gets two new rows, and each row names one kind of staleness from landed work.

`map-sources` grades each `Path:` record of the `## Sources` section of every decision map at the source tip.
A record is red when the published tree does not hold its path as a regular file.
The published tree is the committed tree of the source tip after the delivery closure of the spec's approved deliverable.
The row covers each active map in `decisions/` and each compiled map in `specs/*/decisions/`.

`writes-landed` reds each plain `Writes:` entry that the checkout lacks and whose path a landed commit on the default branch deleted.
`writes-resolve` keeps every other unresolved entry and keeps its text.

The build entry takes a red when every red row is `completion-plan`, `writes-landed`, or `map-sources`.
The staleness pass then repoints the stale citation or removes the stale entry in the enabling plan commit.
Any other red still stops the phase.
The row name is the one classification source, so no reader parses a detail string.

## User stories

Line: opus / medium.
Implementation-line reason: the hardest material chunks are ES-C2 and ES-C3, which move one closure derivation and one Sources grammar behind shared seams across `internal/commitment`, `internal/maps`, and `internal/preflight`. The spec fixes each predicate, each detail text, and each seam, and the seams are existing owners. The command-level tests over real repositories red each wrong branch. The project `Lines` section routes gate and conformance logic to mid effort, because a wrong preflight row is a wrong oracle. Ticket 04 edits guidance prose, so the leverage override of `craft-line` runs its author at high effort.
Harder chunks: ES-C2, ES-C3.

### `map-sources` grades decision-map sources against the published tree

1. As an orchestrator, I want the build preflight to red a decision-map `Path:` source that the delivery closure of my spec deletes, so that I repoint it before the first chunk and not after the complete checkpoint.
2. As an orchestrator, I want the grade to cover each active map in `decisions/`, so that a top-level map such as the FT293 occurrence cannot red my final landing.
3. As an orchestrator, I want the grade to cover the compiled maps of every spec folder, so that the compiled map of another spec cannot red my final landing.
4. As an orchestrator, I want a `Path:` source that the source tip does not hold to red the row, so that a cited file that landed work already removed does not wait for the gate.
5. As an orchestrator, I want a source that the tip holds and the closure keeps to stay green, so that the row reds only real staleness.
6. As an orchestrator, I want the row to grade the committed tip and not the checkout, so that an uncommitted file or map edit cannot pass as published.
7. As a spec author, I want a malformed Sources record to red the row with the shared grammar diagnostic, so that a broken record never becomes an empty green.
8. As a spec author, I want a well-formed `URL:` record to stay green, so that a web source does not red a tree check.
9. As a kit maintainer, I want the Sources grammar to stay in one owner in `internal/maps`, so that the gate check and the preflight row cannot drift.
10. As a kit maintainer, I want the preflight to predict the closure with the same function that the delivery uses, so that the preflight cannot predict a closure that the landing does not make.
11. As an orchestrator, I want a spec without an approved binding to predict no closure deletion, so that the row still grades the tip without a commitment.
12. As an orchestrator, I want a closure prediction that refuses to red the row with its reason, so that a refused prediction never reads as no deletion.
13. As an orchestrator, I want a commitment policy that does not parse at the tip to red the row, so that a broken policy never reads as no closure.
14. As a user of a linked repository, I want the row green when no decision map exists, so that a repository without shaping maps builds as before.
15. As an orchestrator, I want the red detail to name the map, the cited path, and the reason, so that the staleness pass repoints each record without a second read.
16. As an orchestrator, I want the two new rows only in build mode, so that the review preflight keeps its rows.
17. As an orchestrator, I want a `Path:` source whose tip entry is a symbolic link to red, so that the row never follows a link out of the published tree.
18. As an orchestrator, I want a map that the tip holds over the control-record limit to red with its reason, so that an unreadable map never reads as a map without sources.

### `writes-landed` names a `Writes:` entry that landed work deleted

19. As an orchestrator, I want a plain `Writes:` entry whose path a landed commit deleted to red `writes-landed`, so that the build entry tells landed work from a typing error.
20. As an orchestrator, I want the bare build preflight at a fresh entry to make that classification, so that the 2026-10-17 shape gets the new row without a range or a recorded chunk.
21. As an orchestrator, I want a deletion that a landing merge made against its first parent to count as landed, so that a delivery closure deletion gets the new row.
22. As an orchestrator, I want a deletion that a side-branch commit made before its merge into the default branch to count as landed, so that a ticket deletion of an earlier build gets the new row.
23. As an orchestrator, I want a rename on the default branch to count as a deletion of the old path, so that a moved file gets the new row.
24. As an orchestrator, I want an entry whose path the default branch never tracked to stay a `writes-resolve` red, so that a typing error still stops the phase.
25. As an orchestrator, I want an entry whose path the default-branch tip holds again to stay a `writes-resolve` red, so that a deleted and restored file is not landed staleness.
26. As an orchestrator, I want a mixed set of unresolved entries to red both rows, so that one typing error still stops the phase.
27. As an orchestrator, I want an entry that an earlier chunk of my build deleted to stay out of both reds, so that the current earlier-chunk rule keeps its answer.
28. As an orchestrator, I want an entry that the explicit `--base` range deleted to stay out of both reds, so that the current range rule keeps its answer.
29. As an orchestrator, I want a path with a space or a glob character to answer by its literal name, so that a sibling path cannot change the answer.
30. As an orchestrator, I want an unresolved default branch to leave each unresolved entry in `writes-resolve`, so that a missing fact fails closed.
31. As a kit maintainer, I want `writes-resolve` to keep its current detail text for the entries that it keeps, so that each reader of that text stays true.
32. As an orchestrator, I want the `writes-landed` detail to name each ticket, each entry, and the deleting commit, so that the staleness pass can cite the landed work.
33. As an orchestrator, I want the review preflight to keep every unresolved entry in `writes-resolve`, so that a chunk review grades its source as before.
34. As an orchestrator, I want `writes-landed` not applicable in build mode without a tickets folder, so that it obeys the one ticket-row gate.

### The build entry takes a red from landed work

35. As an orchestrator, I want the build entry to take a red when every red row is `completion-plan`, `writes-landed`, or `map-sources`, so that I do not stop for landed work.
36. As an orchestrator, I want any other red row to stop the phase, so that a real plan defect still reaches the reviewer.
37. As a kit maintainer, I want an anchor to refuse the old take-rule sentence, so that a revert to the `completion-plan`-only rule turns the gate red.
38. As an orchestrator, I want the staleness pass to state how it repoints a taken `map-sources` record, so that the repoint goes through the widen rule of `.bench/BENCH.md`.
39. As an orchestrator, I want the audit to require the manual pass while either new row is red, so that nobody can omit the staleness pass over a stale plan.

### This build's own landing stays green

40. As the orchestrator of this build, I want the FT375 record of `decisions/architecture-planning.md` repointed to a source that survives this landing, so that the final landing stays green under `decision-map-integrity`.
41. As the orchestrator of this build, I want that repoint to land in the ticket that adds `map-sources` or earlier, so that no later charge preparation refuses on that record.

### Records

42. As a reader of the release notes, I want one `CHANGELOG.md` entry that names both rows and the take rule, so that the change is visible after the row file retires.

### Reviewed exclusions

43. As a reviewer, I want FT375 not to claim the known-issues row about a spec seam without research, so that the scope stays as the decision source states it.
44. As a reviewer, I want a `(deleted)` entry to keep its current grade, so that this spec changes only the plain-entry classification.
45. As a reviewer, I want no prediction of map citations that the build's own `(deleted)` entries remove, so that the row grades the closure that the source names.

## Implementation decisions

### Row order and applicability

`Decide` adds `writes-landed` directly after `writes-resolve` in build mode only, through `ticketRow` (`internal/preflight/decision.go:268`).
`Decide` adds `map-sources` directly after `binary-seal` in build mode only, so `binary-seal` stays directly after `kit-pin` (`internal/preflight/binary_seal_test.go:38`).
`map-sources` reads no ticket, so it does not go through `ticketRow` and is never not applicable.
Neither row carries `exempt`, so a red or a not-applicable new row keeps `requiredGreen` false (`internal/preflight/decision.go:254`).
That value feeds the audit's `manual_required` through `staleness.Decide` (`internal/preflight/command.go:199`).
The `Decide` doc comment names both build-only rows and their positions.
A bare build verdict then renders 14 rows, and a review verdict keeps 11 rows (`internal/preflight/source_tip_test.go:13`).
Every build-mode preflight form renders both rows: the bare form, `--plan-only`, `--charge`, the current-evidence form, and `--audit-spec`.

### The `writes-landed` fact

The gatherer adds one fact, `LandedDeletions`, a map from a `Writes:` entry to the full identity of its deleting commit.
The gatherer fills it in build mode only, so the review preflight keeps today's `writes-resolve` answer.
It considers a plain entry only when both of these hold:

- `tickets.ParseWrite` accepts the entry, and the kind is plain (`internal/tickets/writes.go:41`).
- The checkout does not hold the path (`internal/preflight/tree_probe.go:64`, `treeHolds`).

The gatherer does not read `DeletedPaths` or `BuildDeletedPaths`. The check alone owns the precedence below, so the precedence has one source.

For each such entry, the gatherer resolves the default branch through `git.ResolvedDefault` (`internal/git/default_branch.go:22`) and its tip through `git.ResolveCommit` (`internal/git/git.go:249`).
The entry is landed when the default-branch tip lacks the path and a first-parent commit of that tip deleted it.
The tip lookup is `git ls-tree <tip> -- :(literal)<path>`, and an empty listing means the tip lacks the path.
The deleting-commit query is `git log -1 --format=%H --diff-filter=D --no-renames --first-parent -m <tip> -- :(literal)<path>`.
A failed query, an unresolved default branch, or an empty answer records no entry, so the entry stays a `writes-resolve` red.

The query flags rest on an observed run over this repository at `a2f793957` with git 2.43.0:

- A ticket commit on a build branch deleted `tests/canary/workflow-guidance-anchors/context-coverage-map-term/BASE`, and the landing merge brought it in. The query answers the landing merge `385464884`. The directory entry of the same canary answers the same commit.
- A delivery closure deleted `roadmap/FT217.md` inside the landing merge. The query answers `8ab514845`. The same query without `--first-parent -m` answers nothing, and the query with `-m` alone answers the side-branch fold merge `5d1119b26`.
- The closure of FT293 deleted `roadmap/FT293.md`, and the query answers the landing merge `4d547d73b`.
- A rename moved `decisions/background-grade.md`, and the query with `--no-renames` answers `3906eb0f0`.
- A path that the default branch never tracked answers nothing, and a path that the tip holds (`roadmap/FT375.md`) answers nothing.
- The history holds no deleted path with a space, no never-tracked path with a glob character, and no deleted and restored path. ES35 and ES31 observe those shapes in fixtures. The `:(literal)` pathspec magic is the tree convention for exact paths (`internal/landing/published/published.go:117`).
- With a path filter, the query answers `3906eb0f0` for the renamed path with or without `--no-renames`, because the filter stops rename pairing. The spec keeps `--no-renames` so that the answer does not depend on that pairing rule.

On git 2.43.0, `--diff-merges=first-parent` prints a patch, so the spec names `--first-parent -m`. On git 2.43.0, `--first-parent` alone also answers `8ab514845`, and `-m` keeps the answer on an older git.
The query reads only first-parent commits of the default branch and the first-parent diff of each merge. A path that only a merged side branch created and deleted is therefore never landed, and it stays a `writes-resolve` red.
The worst-case query over the 2,865 first-parent commits of this repository took 29 ms.

### The `writes-resolve` and `writes-landed` precedence

Both checks walk the same entries in ticket and entry order, and each entry gets exactly one answer:

1. An entry that the grammar refuses, a `(new)` entry, a repeated entry, or an entry whose path `BuildDeletedPaths` holds is out of both rows (`internal/preflight/fence_writes.go:58`).
2. A `(deleted)` entry keeps the current review-base rule of `writes-resolve` (`internal/preflight/fence_writes.go:59`).
3. A plain entry that the checkout holds, or whose path `DeletedPaths` holds, is green in both rows.
4. A plain entry that `LandedDeletions` holds is a `writes-landed` red.
5. Every other plain entry is a `writes-resolve` red.

`writes-resolve` keeps the detail text `Writes: entry names no tree path and carries no (new) marker: <ticket>: <entry>` for the entries of step 5.
`writes-landed` uses the detail text `Writes: entry names a path that a landed commit deleted: <ticket>: <entry> (<commit>)`, joined by `, ` for each entry.
Steps 3 and 4 can both hold for one entry: a `--base` range that holds a fold of a landed deletion puts the path in `DeletedPaths` and in `LandedDeletions`. Step 3 wins, so that entry is green in both rows.
The `<commit>` is the full commit identity, so its length does not depend on the repository.

### One closure prediction

The commitment package gets one pure function that returns the outcome sources that the delivery of an approved deliverable closes.
It holds the `activeDeliverable` lookup, the obligation-free refusal, and the obligation filter of the current `Deliver` body (`internal/commitment/delivery.go:17`, `internal/commitment/delivery.go:43`).
It does not hold the already-delivered refusal (`internal/commitment/delivery.go:26`), because that refusal reads the delivery fact.
`Deliver` calls it for its closed set and keeps its delivery-fact steps, the already-delivered refusal included.

The repository store gets one method that reads the policy at a tree, the way `policyAt` does (`internal/commitment/repository/repository.go:261`).
The method calls the pure function for the deliverable and returns each closed source path that `roadmap.RowOwner` admits (`internal/roadmap/closure.go:17`).
`Store.Closure` uses the same row-owner step for its deletion edits, so the preflight and the landing share one derivation.
A policy that the tree does not hold, or a deliverable that the active milestone does not approve, gives no path.
A refusal or a policy that does not parse returns its error.

### One Sources grammar with two resolvers

`internal/maps` keeps one Sources record grammar: the record split, the field order, the locator rule, the URL rule, and the lexical path rule (`internal/maps/validation.go:143`, `internal/maps/validation.go:238`).
The lexical path rule refuses an empty, absolute, or escaping path before a resolver runs.
The lexical path rule gives the resolver the cleaned slash path, so `./x` and `a/../x` reach the resolver as `x`.
A resolver answers whether the path names a regular file, and it gives its own reason text when it does not.
The working-tree resolver keeps each current reason of `validateSourcePath`, the `escapes the repository` reason for a link that leads outside the repository included (`internal/maps/validation.go:271-273`).

- The working-tree resolver is the current file-system half of `validateSourcePath` (`internal/maps/validation.go:249`). `ValidateDecisionMap` keeps it, so `decision-map-integrity` keeps its messages.
- The published-tree resolver reads the tree entry of the path at a commit. It accepts a mode `100644` or `100755` entry whose path the closure prediction does not hold.
- The published-tree resolver gives the reason `must name an existing regular file` for a missing, link, directory, or submodule entry. It gives the reason `is deleted by the delivery closure of <deliverable>` for a predicted deletion.

`internal/maps` exports one function that takes the map path, the map bytes, and a resolver, and returns the Sources diagnostics with the map path prefix.
`internal/maps` also exports the published-tree resolver, so both reason texts have one owner.
The unexported `sourceDiagnostics` keeps its signature as a wrapper, so its two test callers stay (`internal/maps/maps_graph_test.go:142`, `internal/maps/maps_graph_test.go:170`).

### The candidate maps at the tip

The candidate rule stays one rule in `internal/maps` (`internal/maps/tree_validation.go:24`).
`internal/maps` exports one predicate that admits a repository-relative path as an active or compiled candidate.
The predicate admits `decisions/<name>.md` and `specs/<folder>/decisions/<name>.md`, where the folder and the name do not start with a dot and the name is not a README.
The file-system walk calls the same predicate for each entry that it finds.
The preflight lists the tip tree under `decisions` and `specs` and keeps each path that the predicate admits.
A candidate whose tip entry has a mode other than `100644` or `100755` gives a diagnostic, because a link entry reads as its target string and holds no Sources section.
The preflight reads each candidate with the bounded `git.ReadControlBlob` (`internal/git/tree.go:240`), and an unreadable or over-limit map gives a diagnostic.

### The `map-sources` fact and its red

The gatherer adds one fact, `MapSourceDiagnostics`, an ordered list, and fills it in build mode only.
It reads the source tip `SourceTip`, which is HEAD in the bare form and the snapshot head with `--base`.
It predicts the closure of the deliverable at `SpecPath` from the policy at that tip, and a prediction error becomes one diagnostic.
It grades each candidate map with the published-tree resolver, in candidate path order.
`map-sources` is red with the diagnostics joined by `; ` when the list is not empty, and green otherwise.
Each diagnostic reads `<map>: Sources Path <path> <reason>` or `<map>: <grammar message>`.

### Guidance

The take-rule sentence of `.agents/commands/bench-implement-spec.md:15` becomes `` If every red row is `completion-plan`, `writes-landed`, or `map-sources`, the staleness pass below takes the red. ``
The next sentence stays `` Any other red stops the phase. ``
The first two sentences of that paragraph do not change, because the canary `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight/MUTATE.json` quotes them.
The new text holds no `fence-writes`, `-closure`, or `check-floor`, because Forbid anchors refuse each one (`internal/anchors/registry_ticket_passes.go:24-25`, `internal/anchors/registry_build_cost.go:11`).
The Require needles of `internal/anchors/registry_data.go:92-93` and `internal/anchors/registry_ft311_preparation.go:8-10` stay byte for byte.

The new take-rule sentence and the stop sentence get no Require row, because the admission classes admit no rule wording (`internal/anchors/registry.go:26-41`, `internal/anchors/locate_test.go:93-111`).
`internal/anchors/registry_ticket_passes.go` gets one Forbid row on `.agents/commands/bench-implement-spec.md` with the old exact sentence `` If every red row is `completion-plan`, the staleness pass below takes the red. ``
`TestTicketSlicingAnchors` gets the same Forbid rule, written independently of the registry, as it does for the earlier retired sentence (`internal/anchors/registry_ticket_passes_test.go:31`).
A revert to the `completion-plan`-only rule then turns the gate red.
`.agents/skills/bench-implement-spec/references/staleness-pass.md` gets one sentence under `Orchestrator review`: a taken `map-sources` red repoints each named record to a source that the published tree holds, and a path outside the signed-off `Writes:` lines takes the widen rule of `.bench/BENCH.md`.

### This build's own map record

Ticket 03 repoints the record at `decisions/architecture-planning.md:45-47`.
The new record is `Path: .agents/skills/bench-implement-spec/references/staleness-pass.md`, `Supports: the omission rule and the taken red rows of the build-entry staleness pass.`, and `Drift: a change to its omission rule or its taken red rows.`
That file holds the shipped answer of decision ticket 2 of the same map, and this landing does not delete it.

### Records

The `## [Unreleased]` section of `CHANGELOG.md` gets one entry that names `writes-landed`, `map-sources`, and the new take rule.

## Implementation chunks

ES-C1 ships the `Writes:` classification alone, because its row is green on every current build and needs no other seam.
ES-C2 ships the closure prediction alone, because it is a seam that ES-C3 consumes, and its review closes before a consumer ticket starts.
ES-C3 ships the map row, the take rule, and the dogfood repoint, because the take rule names both rows.
ES-C1 and ES-C2 share no write, so their authors can work at the same time.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| ES-C1 / 01-writes-landed-row.md | The build preflight names a `Writes:` entry that landed work deleted in its own row | ES24, ES25, ES26, ES27, ES28, ES29, ES30, ES31, ES32, ES33, ES34, ES35, ES36, ES37, ES38, ES39, ES40, ES45 | TestWritesLandedBareEntry, TestWritesLandedShapes, TestWritesLandedPrecedence, TestWritesLandedReviewMode, TestWritesLandedKeepsAuditManual, TestTicketGrammarRowsCarryOwnDetails, TestTicketGatedRowsNotApplicableWithoutTickets, TestDecideAllGreen, TestDecideChargeRows, TestSourceTipOmittedKeepsTodaysVerdict, TestSourceTipAcceptedByBothModes | no |
| ES-C2 / 02-closure-prediction.md | The commitment owner predicts the closed row paths of a deliverable without completion evidence, from the derivation that the landing uses | ES13, ES14 | TestCommitmentCloses, TestClosedRowPathsMatchClosure | yes |
| ES-C3 / 03-map-sources-row.md, 04-entry-take-rule.md | The build preflight grades every decision-map source against the published tree, and the build entry takes both new reds | ES01, ES02, ES03, ES04, ES05, ES06, ES07, ES08, ES09, ES10, ES11, ES12, ES15, ES16, ES17, ES18, ES19, ES20, ES21, ES22, ES23, ES41, ES42, ES43, ES44, ES46, ES47, ES48, ES49, ES50, ES51, ES52, ES53 | TestMapSourcesClosureDeletion, TestMapSourcesTipExistence, TestMapSourcesGrammar, TestMapSourcesPrediction, TestMapSourcesBuildOnly, TestMapSourcesKeepsAuditManual, TestMapSourcesCollectIndependentRecordFailures, TestSourceTipOmittedKeepsTodaysVerdict, TestDecideChargeRows, TestTicketSlicingAnchors | yes |

### Completion plan

This version 1 plan records the author verification and the final checks.
The orchestrator amends it to version 2 before the first author dispatch.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "ES-C1",
      "tickets": [
        "01-writes-landed-row.md"
      ],
      "verification": [
        {
          "id": "writes-landed",
          "command": "bench test --package ./internal/preflight --run 'WritesLanded|TicketGrammarRowsCarryOwnDetails|TicketGatedRowsNotApplicableWithoutTickets|DecideAllGreen|DecideChargeRows|SourceTip'",
          "probe": "Drop --first-parent -m from the deleting-commit query, and record that TestWritesLandedShapes turns red on the landing-merge case."
        },
        {
          "id": "precedence",
          "command": "bench test --package ./internal/preflight --run 'WritesLandedPrecedence'",
          "probe": "Grade LandedDeletions before DeletedPaths in the check, and record that TestWritesLandedPrecedence turns red on the folded explicit-range case of ES33."
        }
      ]
    },
    {
      "id": "ES-C2",
      "tickets": [
        "02-closure-prediction.md"
      ],
      "verification": [
        {
          "id": "closure-prediction",
          "command": "bench test --package ./internal/commitment --package ./internal/commitment/repository --run 'CommitmentCloses|ClosedRowPathsMatchClosure|CommitmentDeliver|CommitmentLightClosure|CommitmentClosureRefusesObligationFreeBinding'",
          "probe": "Skip the RowOwner step in the store method, and record that TestClosedRowPathsMatchClosure turns red on the refused obligation source of ES14."
        }
      ]
    },
    {
      "id": "ES-C3",
      "tickets": [
        "03-map-sources-row.md",
        "04-entry-take-rule.md"
      ],
      "verification": [
        {
          "id": "map-sources",
          "command": "bench test --package ./internal/preflight --run 'MapSources|SourceTipOmittedKeepsTodaysVerdict|DecideChargeRows'",
          "probe": "Grade the map sources with the working-tree resolver, and record that TestMapSourcesClosureDeletion turns red."
        },
        {
          "id": "shared-grammar",
          "command": "bench test --package ./internal/maps --package ./internal/preflight --run 'MapSources|ValidateDecisionMapTree'",
          "probe": "Drop the lexical path rule from the shared grammar, and record that TestMapSourcesGrammar turns red on the escaping-path case of ES52."
        },
        {
          "id": "take-rule",
          "command": "bench test --package ./internal/anchors --run 'TicketSlicingAnchors'",
          "probe": "Restore the old completion-plan-only take-rule sentence, and record that TestTicketSlicingAnchors turns red."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "preflight-packages",
      "command": "bench test --package ./internal/preflight/..."
    },
    {
      "id": "owner-packages",
      "command": "bench test --package ./internal/commitment/... --package ./internal/maps --package ./internal/anchors --package ./internal/landing/..."
    },
    {
      "id": "coverage",
      "command": "bench coverage --check preflight-entry-staleness"
    }
  ]
}
```

## Testing decisions

- A good test drives `Command` in package `preflight` over a real fixture repository and reads the rendered row by name. `TestWritesResolveEarlierChunkDeletion` is the precedent (`internal/preflight/fence_writes_test.go:20`).
- The fixture repository comes from `preflighttest.SeedConformant` (`internal/preflight/preflighttest/fixture.go:172`). The closure comes from `commitmenttest.SeedClosure` (`internal/commitment/commitmenttest/repo.go:195`), which binds the spec with obligation `FT1` and writes `roadmap/FT1.md`.
- One shared helper in `internal/preflight/preflighttest` writes a decision map with given Sources records. Both new test files call it, so no test file pastes a map body.
- A test makes default-branch history with `preflighttest.RunGit`: a plain deletion, a side-branch deletion with a no-fast-forward merge, a merge whose resolution deletes the path, a rename, a delete-then-restore, and a path with a space and a glob character.
- Each expected commit identity comes from `git rev-parse` of the commit that the fixture made. Each expected detail is a literal in the test. No expectation comes from the code under test.
- The pure closure function and the store method get unit tests next to their precedents (`internal/commitment/delivery_test.go`, `internal/commitment/repository/closure_test.go:56`).
- The ordinary Go test phase of the gate observes every row except the review-owned rows.

### Seam diagram

    trigger: bench preflight build <slug> [--plan-only | --charge ... | --audit-spec ...]  (build mode)
        │
        ▼
    source tip + checkout ──▶ [ gather: LandedDeletions, MapSourceDiagnostics ] ──▶ [ Decide ] ──▶ rows writes-resolve, writes-landed, map-sources
                                   │                     │
                                   │                     └─ maps: candidate predicate + shared Sources grammar + published-tree resolver
                                   └─ commitment: closure prediction (shared with Deliver and Store.Closure)
                              ◀ tests attach here: Command([]string{"build", slug, ...}) over a preflighttest repository;
                                unit tests at the closure function, the store method, and the maps grammar

    trigger: the build entry of /bench-implement-spec
        │
        ▼
    red rows ──▶ [ take rule in bench-implement-spec.md ] ──▶ staleness pass or stop
                      ◀ tests attach here: TestTicketSlicingAnchors

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| ES01 | 1 | A bare build preflight reds `map-sources` when a compiled map cites `roadmap/FT1.md` and the approved binding closes FT1 | planned TestMapSourcesClosureDeletion in internal/preflight/map_sources_test.go | A row that reads only the tip tree finds `roadmap/FT1.md` and stays green |
| ES02 | 2 | The same red appears when only an active map in `decisions/` cites `roadmap/FT1.md` | planned TestMapSourcesClosureDeletion in internal/preflight/map_sources_test.go | A walk of only `specs/<slug>/decisions/` misses the FT293 shape |
| ES03 | 3 | The same red appears when only the compiled map of another spec folder cites `roadmap/FT1.md` | planned TestMapSourcesClosureDeletion in internal/preflight/map_sources_test.go | A walk of only the active spec folder misses that map |
| ES04 | 4 | `map-sources` reds a `Path:` source that the source tip does not hold | planned TestMapSourcesTipExistence in internal/preflight/map_sources_test.go | A row that grades only predicted deletions passes a file that is already gone |
| ES05 | 5 | `map-sources` is green when each `Path:` source is a regular file at the tip and the closure keeps it | planned TestMapSourcesTipExistence in internal/preflight/map_sources_test.go | A row that reds every `roadmap/` source over-matches |
| ES06 | 6 | `map-sources` reds a `Path:` source that only the checkout holds as an uncommitted file | planned TestMapSourcesTipExistence in internal/preflight/map_sources_test.go | A working-tree lookup passes a file that the landing never publishes |
| ES07 | 6 | `map-sources` keeps the red when the checkout holds an uncommitted repoint of the map | planned TestMapSourcesTipExistence in internal/preflight/map_sources_test.go | A row that reads the checkout map passes a repoint that the tip does not hold |
| ES08 | 7 | `map-sources` reds a Sources record without a `Drift:` field with the detail `Sources <path> missing Drift` | planned TestMapSourcesGrammar in internal/preflight/map_sources_test.go | A parse that collects only `Path:` locators skips the broken record and stays green |
| ES09 | 7 | `map-sources` reds a Sources entry that starts with neither `Path:` nor `URL:` | planned TestMapSourcesGrammar in internal/preflight/map_sources_test.go | A locator filter drops the entry and stays green |
| ES10 | 8 | `map-sources` is green for a well-formed `URL:` record | planned TestMapSourcesGrammar in internal/preflight/map_sources_test.go | A resolver that treats a URL as a tree path reds every web source |
| ES11 | 9 | The working-tree grammar still reports each independent record failure of a Sources body | `internal/maps/maps_graph_test.go` (`TestMapSourcesCollectIndependentRecordFailures`) | A kept-rule witness: the grammar move must keep the gate messages |
| ES12 | 9 | No production file under `internal/preflight` holds a Sources record literal | review-owned: `rg -n '"Path: \|"URL: \|"## Sources' internal/preflight --glob '!*_test.go'` returns no hit | A second splitter in the preflight is a second source for one grammar |
| ES13 | 10 | The pure closure function returns the obligation sources that `Deliver` closes for the same policy and deliverable | planned TestCommitmentCloses in internal/commitment/delivery_test.go | Two derivations of the closed set can disagree on an obligation filter |
| ES14 | 10 | For a binding whose obligations hold a row source and a source whose path `roadmap.RowOwner` refuses, the store method returns only the row path, and `Store.Closure` deletes exactly that path | planned TestClosedRowPathsMatchClosure in internal/commitment/repository/closure_test.go | A store method that skips the row-owner step returns the refused path, which the landing never deletes |
| ES15 | 11 | `map-sources` is green for a map that cites a present `roadmap/` file when no approved binding names the spec | planned TestMapSourcesPrediction in internal/preflight/map_sources_test.go | A prediction that closes every roadmap row reds a source that the landing keeps |
| ES16 | 12 | `map-sources` reds with the text `names no obligation` when the binding of the spec names no obligation | planned TestMapSourcesPrediction in internal/preflight/map_sources_test.go | A prediction that turns its error into no deletion stays green |
| ES17 | 13 | `map-sources` reds when the commitment policy at the tip does not parse | planned TestMapSourcesPrediction in internal/preflight/map_sources_test.go | A prediction that reads a broken policy as absent stays green |
| ES18 | 14 | `map-sources` is green when the tip holds neither `decisions/` nor a `specs/*/decisions/` folder | planned TestMapSourcesTipExistence in internal/preflight/map_sources_test.go | A walk that treats an absent folder as an error reds every linked repository |
| ES19 | 14 | `map-sources` is green when `decisions/` holds only `README.md` and a dot file that cite missing paths | planned TestMapSourcesTipExistence in internal/preflight/map_sources_test.go | A candidate rule apart from the shared predicate grades a file that the gate skips |
| ES20 | 15 | The red detail of ES01 reads `<map>: Sources Path roadmap/FT1.md is deleted by the delivery closure of specs/example/spec.md` | planned TestMapSourcesClosureDeletion in internal/preflight/map_sources_test.go | A detail without the map path makes the staleness pass search every map |
| ES21 | 16 | A review preflight over the ES01 tree prints no `map-sources` row | planned TestMapSourcesBuildOnly in internal/preflight/map_sources_test.go | A row added outside the build-mode branch changes every review verdict |
| ES22 | 18 | `map-sources` reds a map that the tip holds one byte over the control-record limit | planned TestMapSourcesTipExistence in internal/preflight/map_sources_test.go | An unbounded read or a skipped map stays green |
| ES23 | 17 | `map-sources` reds a `Path:` source whose tip entry is a symbolic link to a regular file | planned TestMapSourcesTipExistence in internal/preflight/map_sources_test.go | A resolver that follows the link grades bytes outside the published entry |
| ES24 | 20 | A bare build preflight with no recorded chunk reds `writes-landed` for a plain entry whose path a default-branch commit deleted | planned TestWritesLandedBareEntry in internal/preflight/writes_landed_test.go | Today the bare form reds `writes-resolve` for that entry |
| ES25 | 19 | In the ES24 case, `writes-resolve` is green | planned TestWritesLandedBareEntry in internal/preflight/writes_landed_test.go | A row that copies the entry into both rows still stops the phase |
| ES26 | 21 | `writes-landed` reds an entry whose path only the resolution of a merge into the default branch deleted | planned TestWritesLandedShapes in internal/preflight/writes_landed_test.go | A query without `--first-parent -m` finds no deleting commit |
| ES27 | 22 | `writes-landed` reds an entry whose path a side-branch commit deleted before a merge into the default branch | planned TestWritesLandedShapes in internal/preflight/writes_landed_test.go | A query that walks the first-parent commits without the first-parent diff of each merge misses a merged deletion |
| ES28 | 23 | `writes-landed` reds an entry whose path a default-branch rename moved away | planned TestWritesLandedShapes in internal/preflight/writes_landed_test.go | A classifier that admits only a plain deletion and not the old side of a rename leaves the moved path a `writes-resolve` red |
| ES29 | 24 | `writes-resolve` reds an entry whose path no default-branch commit tracked | planned TestWritesLandedShapes in internal/preflight/writes_landed_test.go | A classifier that calls every absent path landed hides a typing error |
| ES30 | 26 | With one landed entry and one never-tracked entry, `writes-landed` and `writes-resolve` are both red | planned TestWritesLandedPrecedence in internal/preflight/writes_landed_test.go | A check that classifies the whole red by its first entry takes a typing error |
| ES31 | 25 | `writes-resolve` reds an entry whose path the default-branch tip holds again after its deletion | planned TestWritesLandedShapes in internal/preflight/writes_landed_test.go | A query that reads only the deleting commit calls a restored path landed |
| ES32 | 27 | An entry whose path an earlier recorded chunk deleted is green in both rows | planned TestWritesLandedPrecedence in internal/preflight/writes_landed_test.go | A kept-rule witness: the row split must keep the earlier-chunk answer of `BuildDeletedPaths` |
| ES33 | 28 | An entry whose path a default-branch commit deleted, which the build folded inside its explicit `--base` range, is green in both rows | planned TestWritesLandedPrecedence in internal/preflight/writes_landed_test.go | A check that grades `LandedDeletions` before `DeletedPaths` reds an entry that the range already resolves |
| ES34 | 30 | With no resolvable default branch, a landed entry reds `writes-resolve` | planned TestWritesLandedPrecedence in internal/preflight/writes_landed_test.go | A missing fact that reads as landed lets the build take an unknown red |
| ES35 | 29 | `writes-resolve` reds a never-tracked entry `d/gl*ob x` when a default-branch commit deleted its sibling `d/glXob x` | planned TestWritesLandedShapes in internal/preflight/writes_landed_test.go | A query with a glob pathspec matches the deleted sibling and calls the entry landed |
| ES36 | 31 | `writes-resolve` keeps the detail `Writes: entry names no tree path and carries no (new) marker: <ticket>: <entry>` | `internal/preflight/decision_test.go` (`TestTicketGrammarRowsCarryOwnDetails`) | A kept-rule witness: a rewrite of the kept text breaks each reader of it |
| ES37 | 32 | The ES24 detail reads `Writes: entry names a path that a landed commit deleted: one.md: <path> (<commit>)` with the full deleting commit | planned TestWritesLandedBareEntry in internal/preflight/writes_landed_test.go | A detail without the commit makes the staleness pass search the history |
| ES38 | 33 | A review preflight over the ES24 tree prints no `writes-landed` row and reds `writes-resolve` for the entry | planned TestWritesLandedReviewMode in internal/preflight/writes_landed_test.go | A classifier that runs in review mode changes the chunk review verdict |
| ES39 | 34 | In build mode without a tickets folder, `writes-landed` is not applicable between `writes-resolve` and `kit-pin` | `internal/preflight/decision_test.go` (`TestTicketGatedRowsNotApplicableWithoutTickets`) | A row outside `ticketRow` grades a fresh build that has no tickets |
| ES40 | 33 | A review preflight whose explicit `--base` comes after the deletion reds `writes-resolve` for the landed entry | planned TestWritesLandedReviewMode in internal/preflight/writes_landed_test.go | A gatherer that fills the fact in review mode moves the entry out of `writes-resolve` |
| ES41 | 16 | A bare build verdict counts 14 rows and a bare review verdict counts 11 rows | `internal/preflight/source_tip_test.go` (`TestSourceTipOmittedKeepsTodaysVerdict`) | A row added to both modes changes the review count |
| ES42 | 16 | The charge verdict rows hold `writes-landed` after `writes-resolve` and `map-sources` after `binary-seal` | `internal/preflight/charge_cells_test.go` (`TestDecideChargeRows`) | A row in the wrong position moves the first red that the charge refusal names |
| ES43 | 35 | `.agents/commands/bench-implement-spec.md` holds the sentence `` If every red row is `completion-plan`, `writes-landed`, or `map-sources`, the staleness pass below takes the red. `` | review-owned: line 15 of `.agents/commands/bench-implement-spec.md` | A sentence that drops a row name stops the build for landed work |
| ES44 | 36 | `.agents/commands/bench-implement-spec.md` holds the sentence `Any other red stops the phase.` directly after the take-rule sentence | review-owned: line 15 of `.agents/commands/bench-implement-spec.md` | A paragraph that drops the stop takes every red into the staleness pass |
| ES51 | 37 | `.agents/commands/bench-implement-spec.md` holds no sentence `` If every red row is `completion-plan`, the staleness pass below takes the red. `` | `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A revert to the old take rule stops the build for landed work |
| ES45 | 39 | `--audit-spec` reports `manual_required` true while `writes-landed` is red | planned TestWritesLandedKeepsAuditManual in internal/preflight/writes_landed_test.go | An `exempt` flag on the row lets the audit omit the staleness pass |
| ES46 | 39 | `--audit-spec` reports `manual_required` true while `map-sources` is red | planned TestMapSourcesKeepsAuditManual in internal/preflight/map_sources_test.go | An `exempt` flag on the row lets the audit omit the staleness pass |
| ES47 | 38 | The `Orchestrator review` section of `staleness-pass.md` states the repoint of a taken `map-sources` record and names the widen rule of `.bench/BENCH.md` | review-owned: the `Orchestrator review` section of `.agents/skills/bench-implement-spec/references/staleness-pass.md` | A pass without the rule repoints a record outside the fence and the landing refuses it |
| ES48 | 40 | `decisions/architecture-planning.md` holds no `Path:` record that cites `roadmap/FT375.md` | review-owned: `bench preflight build preflight-entry-staleness` after ticket 03 prints no `map-sources` red | The kept record reds this build's own final landing under `decision-map-integrity` |
| ES49 | 41 | The repoint of ES48 lands in ticket 03 or in an earlier ticket | review-owned: the ticket 03 `Writes:` line names `decisions/architecture-planning.md` | A later repoint makes each charge preparation after ticket 03 refuse on `map-sources` |
| ES50 | 42 | The `## [Unreleased]` section of `CHANGELOG.md` holds one entry that names `writes-landed`, `map-sources`, and the take rule | review-owned: the `## [Unreleased]` section of `CHANGELOG.md` | A missing entry leaves the change only in code after the row file retires |
| ES52 | 7 | `map-sources` reds a Sources record `Path: ../outside.md` with the shared escape diagnostic | planned TestMapSourcesGrammar in internal/preflight/map_sources_test.go | A published-tree grade without the shared lexical rule looks up the raw path and reports only a missing file |
| ES53 | 18 | `map-sources` reds a candidate map whose tip entry is a symbolic link | planned TestMapSourcesTipExistence in internal/preflight/map_sources_test.go | A read of the link entry gives its target string, which holds no Sources section and grades green |

Not covered: story 43 — the reviewed exclusion of the known-issues seam-research row takes an `Out of scope` line and no behavior.
Not covered: story 44 — the reviewed exclusion of `(deleted)` entries keeps the current review-base rule, which `TestWritesResolveEarlierChunkDeletion` already grades.
Not covered: story 45 — the reviewed exclusion of build-own deletions takes an `Out of scope` line and no behavior.

### Edge inventory

These behaviors serve the kit and every linked repository, because `bench preflight build` ships.
ES18 and ES19 cover the absent-versus-empty pair of `decisions/` and `specs/*/decisions/`, and both answers are green for both audiences.
ES22 and ES23 cover the over-limit file and the link of the hostile-input checklist for the map read and the source entry.
ES35 covers a path with a space and a glob character, and the grammar admits both (`internal/tickets/writes.go:72`, `uncertain` until ES35 runs the ticket parse).
A path with a newline cannot occur in a `Writes:` entry, because the ticket field is one physical line.
No test swaps a package variable.

Won't handle: a `(deleted)` entry whose path landed work already removed — the review-base rule of `writes-resolve` keeps grading it (`internal/preflight/fence_writes.go:59`), and the audit drift names the path.
Won't handle: a path that only a merged side branch created and deleted — no first-parent tree of the default branch held it, so it stays a `writes-resolve` red and the reviewer sees the true answer.
Won't handle: a shallow clone whose history does not hold the deleting commit — the entry stays a `writes-resolve` red, so the phase stops, and the reviewer sees the true answer.
Won't handle: a deletion on the remote default branch that the local default branch does not hold yet — `base-current` already grades the fold of the local default branch.
Won't handle: a `Path:` source that reaches a regular file only through a link inside the repository — `map-sources` reds it by ES23, and `decision-map-integrity` keeps following the link at the gate. No current map cites such a path.
Won't handle: a candidate map stored as a link — `map-sources` reds it by ES53, and `decision-map-integrity` keeps its own link behavior at the gate.
Won't handle: a map source that the build itself creates later — the tip lacks it, so ES04 reds it, and the staleness pass repoints it.
Won't handle: a change to `bench maps` freshness rows for compiled maps — `decision-map-integrity` still grades each compiled map at the gate.
Won't handle: lane selection of `decision-map-integrity` for a commit that deletes only a roadmap file (`internal/gate/lane_select.go:196`) — the whole-project landing gate and the new row cover that case.
Won't handle: the stale `binary-seal` after a Go commit and the spill under `bench worktree exec` (`docs/known-issues.md:10`, `docs/known-issues.md:12`) — they touch the build process only, and FT406 and FT341 own them.

## Out of scope

- A build preflight check that refuses a spec seam without research is a separate capability (`docs/known-issues.md:15`). Estimate: 6 edits, 2 gate runs.
- A prediction of the map citations that the build's own `(deleted)` entries remove is a separate capability. Estimate: 4 edits, 1 gate run.
- Link resolution inside the published tree for a Sources path is a separate capability. Estimate: 3 edits, 1 gate run.

## Further notes

### Source-sentence-to-row table

| source sentence | rows |
| --- | --- |
| Criterion: "Before the first chunk, the build preflight grades each spec decision map's Sources: against the published tree." | ES01, ES04, ES05, ES06, ES07, ES18, ES19, ES20, ES21, ES22, ES23, ES41, ES42 |
| Row: "A delivery closure can remove a cited roadmap file." | ES01, ES13, ES14, ES15, ES16, ES17 |
| D1 (2026-10-31): every decision map, active and compiled | ES02, ES03, ES19 |
| Criterion: "The build-entry staleness pass takes a writes-resolve red whose every named path a landed commit deleted, if that check survives FT293-closure-removal." | ES24, ES25, ES26, ES27, ES28, ES29, ES30, ES31, ES34, ES35, ES43 |
| Row: "the orchestrator does not stop" | ES43, ES44, ES51 |
| Row: "When FT293 lands, make sure that this face still applies to the remaining consistency checks." | ES36 and the Problem section's current-code fact |
| D2 (2026-10-31): the row split and the take of both new rows | ES25, ES30, ES37, ES38, ES39, ES40, ES43, ES44 |
| Decision ticket 2 of `decisions/architecture-planning.md`: "Missing metadata, unresolved source, or unsupported citation forms require a diagnostic." | ES08, ES09, ES16, ES17, ES22, ES45, ES46 |

### Flagged additions

- ES13 and ES14 go beyond the source sentences. The one-source rule of `AGENTS.md` requires one closure derivation for the preflight and the landing.
- ES11 and ES12 go beyond the source sentences. The one-source rule requires one Sources grammar for the gate check and the preflight row.
- ES47 goes beyond the source sentences. D2 takes a `map-sources` red, and the widen rule decides whether a repoint outside the signed-off `Writes:` lines needs a reviewer decision. A repoint of a top-level map can therefore still stop for that decision. This is the contestable call of this spec.
- ES48 and ES49 go beyond the source sentences. This spec's own closure deletes `roadmap/FT375.md`, which `decisions/architecture-planning.md:45` cites.
- ES51 follows the anchor-registry convention for a retired phase-command sentence. ES43 and ES44 stay review-owned, because the admission classes admit no rule wording.
- ES37 goes beyond D2: the `writes-landed` detail names the deleting commit, so the staleness pass can cite the landed work.
- ES50 follows the tree convention for a user-facing change.
- ES22, ES23, ES35, ES52, and ES53 come from the hostile-input checklist for the new reads.
- ES36, ES39, ES41, ES42, ES45, and ES46 are kept-rule witnesses or row-list readers, not new promises.

### Reader sweep

- Row-list and row-count readers: `internal/preflight/source_tip_test.go:13` (`bareRowCount`, 12 becomes 14, read by `TestSourceTipOmittedKeepsTodaysVerdict` and `TestSourceTipAcceptedByBothModes`), `internal/preflight/charge_cells_test.go:19` (`TestDecideChargeRows`, both rows join), `internal/preflight/decision_test.go:534` (`TestTicketGatedRowsNotApplicableWithoutTickets`, `writes-landed` joins), and `internal/preflight/decision_test.go:100` (the not-applicable count of `TestDecideAllGreen` in build mode without tickets, 9 becomes 10). Ticket 01 moves `bareRowCount` to 13, and ticket 03 moves it to 14. The slicing scratch runs turned exactly these five tests red for ticket 01 and the three `source_tip_test.go` and `charge_cells_test.go` tests red for ticket 03.
- `internal/preflight/decision_test.go` holds 557 lines against the 400-line budget, so its edits keep a net growth of zero lines, and `internal/conformance/ordinary_build_census_test.go:158` forbids a repository fixture in it.
- `internal/preflight/gather.go` holds 393 lines and `internal/preflight/preflighttest/fixture.go` holds 394 lines, so the new gather code and the new fixture helper go in new or other files.
- `.agents/commands/bench-implement-spec.md` holds 97 lines against its 97-line guidance budget (`projects/benchkit.md:553`), so the take-rule edit replaces the sentence in place.
- `internal/maps/maps_graph_test.go:79-87` calls `validateSourcePath` directly, and the `tests/canary/decision-map-integrity/source-*` units pin the exact working-tree messages. The working-tree resolver keeps those messages.
- Readers that keep their answer: `internal/preflight/binary_seal_test.go:38` (`binary-seal` stays after `kit-pin`), `internal/preflight/verdict_summary_test.go:74` (counts derive from `Decide`), and `internal/preflight/charge_test.go:350` (`legacyBaseline` takes its counts from `Decide` and renders red rows only).
- The `testdata/legacy-charge` fixtures hold no row count and no `writes-resolve` text, so they keep their bytes.
- Readers of the `writes-resolve` detail: `internal/preflight/decision_test.go:320-346`, `internal/preflight/decision_test.go:417-424`, `internal/preflight/command_build_test.go:237-240`, and `internal/preflight/command_review_test.go:420-437`. Each fixture names a never-tracked path or runs review mode, so each keeps its text. `uncertain` until ticket 01 runs them.
- Readers of the deletion facts: `internal/preflight/fence_writes_test.go:20` and `internal/preflight/deletion_preflight_test.go:11-55`. The review-mode cases keep their answer. The build-mode `record without a chunk` case names a path that no default-branch commit deleted, so it stays a `writes-resolve` red. `uncertain` until ticket 01 runs it.
- Readers of the take rule: `.agents/commands/bench-implement-spec.md:15` and `:17`, and `.agents/skills/bench-implement-spec/references/staleness-pass.md:22-34`. `.claude/commands` is a link to `.agents/commands`, so no copy exists. `.bench/BENCH.md:63` lists no red row.
- Readers of the Sources grammar: `internal/maps/validation.go:121` (`readinessDiagnostics`), `internal/maps/freshness.go:39` (`sourcePathLocators`, unchanged), and `internal/conformance/decision_map_integrity_test.go:106-141` with `tests/canary/decision-map-integrity/` (messages unchanged).
- Readers of the candidate rule: `internal/maps/maps.go:81` (`scanActiveMap`), `internal/maps/tree_validation.go:59`, and `internal/status/status.go:846`.
- Readers of the closure: `internal/commitment/repository/candidate.go:226`, `internal/gate/completion.go:138`, and `internal/landing/published/published.go:62`. Each keeps its call and answer.
- `.github/workflows` and `scripts/*.mjs` hold no caller of `bench preflight build` or of the maps package.

### Pre-review proof checklist

- Cited symbols: `Decide` (`internal/preflight/decision.go:206`), `requiredGreen` (`internal/preflight/decision.go:254`), `ticketRow` (`internal/preflight/decision.go:268`), `writesResolveCheck` (`internal/preflight/fence_writes.go:50`), `deletedBefore` (`internal/preflight/tree_probe.go:28`), `buildDeletedPaths` (`internal/preflight/tree_probe.go:46`), `treeHolds` (`internal/preflight/tree_probe.go:64`), `gather` (`internal/preflight/gather.go:32`), `chargeVerdictRefusal` (`internal/preflight/charge.go:83`), `tickets.ParseWrite` (`internal/tickets/writes.go:41`), `git.ResolvedDefault` (`internal/git/default_branch.go:22`), `git.ResolveCommit` (`internal/git/git.go:249`), `git.ReadControlBlob` (`internal/git/tree.go:240`), `commitment.Deliver` (`internal/commitment/delivery.go:17`), `obligationFree` (`internal/commitment/delivery.go:43`), `activeDeliverable` (`internal/commitment/delivery.go:258`), `Store.Closure` (`internal/commitment/repository/closure.go:59`), `policyAt` (`internal/commitment/repository/repository.go:261`), `roadmap.RowOwner` (`internal/roadmap/closure.go:17`), `published.Tree` (`internal/landing/published/published.go:25`), `sourceDiagnostics` (`internal/maps/validation.go:143`), `sourceLocator` (`internal/maps/validation.go:238`), `validateSourcePath` (`internal/maps/validation.go:249`), `discoverDirectoryCandidates` (`internal/maps/tree_validation.go:24`), `ValidateDecisionMapTree` (`internal/maps/tree_validation.go:96`), `commitmenttest.SeedClosure` (`internal/commitment/commitmenttest/repo.go:195`), and `preflighttest.SeedConformant` (`internal/preflight/preflighttest/fixture.go:172`).
- Import edges: `internal/preflight` gains a direct edge to `internal/maps`. `internal/preflight/ownership` already imports it (`internal/preflight/ownership/metadata.go:13`), and `go list -deps ./internal/maps` names no preflight or commitment package, so the edge makes no cycle. `internal/maps` already imports `internal/git` (`internal/maps/maps.go:12`). `internal/preflight` already imports `internal/commitment/repository` (`internal/preflight/command.go:9`).
- Source-row clauses and occurrences: the table above quotes each clause. The row records two occurrences, the FT217 complete checkpoint of 2026-10-09 and the FT375 build entry of 2026-10-17 (`roadmap/FT375.md:17-18`). The FT293 repoint (`1b79d041f`) is a third occurrence that the row does not record.
- Promised field labels: none. No intent-ledger or review-record field changes, so no `Persisted fields:` line applies.
- Changed-function callers: `bench consumers` at `a2f793957` names the `commitment.Deliver` production caller `internal/commitment/repository/closure.go:35` and test callers in `internal/commitment/delivery_test.go`, `internal/commitment/continuation_test.go:76`, and `internal/spec/retire_test.go:243`. It names the `Store.Closure` callers at `internal/commitment/repository/candidate.go:226`, `internal/gate/completion.go:138`, `internal/landing/published/published.go:62`, `internal/commitment/repository/closure_test.go`, and `internal/commitment/repository/publication_test.go:115`. It names the `sourceDiagnostics` callers at `internal/maps/validation.go:121` and `internal/maps/maps_graph_test.go:142,170`, the `discoverDirectoryCandidates` callers at `internal/maps/maps.go:81` and `internal/maps/tree_validation.go:59`, and the one `writesResolveCheck` reference at `internal/preflight/decision.go:216`.
- Copy survival: ES12 reds a second Sources splitter in the preflight. ES14 reds a closure derivation that skips the shared row-owner step. ES33 reds a second copy of the precedence in the gatherer, because only the check reads `DeletedPaths`.
- Rendered-shape readers: the `writes-resolve` text keeps its producer at `internal/preflight/fence_writes.go:71`, and its readers are in the reader sweep. The two new detail texts and the closure reason text are new, so no reader exists.
- Pin operators: the new command tests compare the exact red detail of a row by name, as `TestCommandBuildRendersFourGrammarRows` does (`internal/preflight/command_build_test.go:237`). The anchor rows use the Require match of the anchor harness.
- Entry reads: the gatherer reads the checkout, the source tip, and the default branch under the preflight root. The fixtures make each of them, so no ambient read remains.
- Derived expectations: each commit identity comes from `git rev-parse` in the fixture, and each detail text is a literal in the test.
- Consolidated rules: none. The spec moves two rules behind shared seams and does not merge two rules.
- Quantified obligations: "each spec decision map" quantifies over the candidates of both folders (ES02, ES03, ES19). "every named path" quantifies over the unresolved entries of one red (ES30). "every red row" quantifies over the take rule's three row names (ES43).
- Workflow-step writes: none. The rows and the take rule add no write, and the staleness pass keeps its enabling plan commit.
