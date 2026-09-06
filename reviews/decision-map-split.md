# Review pickup: decision-map-split

Reviewed pair: base `1019f1f7`, tip `0a9a3b90`, on 2026-09-06. Raw findings: Standards 6, Spec 5, Coverage 3, Falsification 5. Repair targets after the collapse: 8. Items left open for the reviewer: 6.

## Standards

Count: 6. Worst: the doc comment of `internal/anchors/registry_decision_maps.go` states a false universal about its rows and carries a red record.

- auto-fix: reword that doc comment to what the file holds, and delete the sentence that names the tests. Citation: `craft-comments`, the commit or spec owns the red record.
- auto-fix: `ticketFolder` in `tickets.go` and `ticketPath` in `maps.go` derive the ticket-file layout twice. Citation: AGENTS.md, one source per fact.
- auto-fix: `unresolvedState` in `maps.go` forwards one call and adds no meaning. Citation: the smell baseline, Lazy Element.
- ask-user: the six-cell row order is derived at six sites in `maps.go` and `freshness.go`, and read by position. A row struct with one render method would single-source it. Citation: the smell baseline, Primitive Obsession. Left open, because it is a seam change across a landed fence.
- no-op: the template gist line is copied as a `strings.Replace` needle in two test packages. The failure is loud, and the spec forbids a cross-package test helper.
- no-op: `var registry = append(generalAnchors, decisionMapAnchors...)` is sound today. `slices.Concat` would state the intent. Left as a note.

## Spec

Count: 5. Worst: DS34 and DS39 are review-owned rows whose evidence lived only in the migration delegate's return, and the program is deleted.

- auto-fix: record the migration evidence in the spec's Build decisions, from the delegate's return and the coordinator's saved pre-migration rows. Citation: rows DS34 and DS39.
- ask-user: the gist of `gate-critical-path` #1 reads `With lever 1 landed.` and does not state the decision. `worktree-orphan-retirement` #6 reads `Bound it.`. The build decisions already route these to the reviewer. Left open, because a gist is reviewer prose.
- ask-user: the `gate-pipeline-fixture-inventory.md` move into the map's assets folder is authorized by no story. It is recorded under Build decisions. Left open for assent or veto.
- auto-fix: the README sentence that describes the index and the ticket files carries no anchor, so half of DS47 has no bite. Citation: row DS47.
- no-op: the proof checklist names `ParseDecisionMap`, which the contract unexported, and DS33 names a test the build superseded with `TestAppendMapsCountsSplitMaps`. Recorded-state drift; the retro notes it.

## Coverage

Count: 3. Worst: a gist may link any folder, because the parser reads only the ticket number.

- auto-fix: a gist whose link target is not this map's `<topic>/tickets/<n>.md` reds with the missing-ticket diagnostic. New fixture `gist-wrong-folder` and row DS51. Citation: DS12 intends every gist to resolve, and the edge inventory says the lane checks the link target.
- auto-fix: two gists for one ticket red with `Decisions so far duplicate gist for ticket #n`. New fixture `gist-duplicate` and row DS52. Citation: every other duplicate in the family reds. The falsification pass found the same gap.
- ask-user: a `tickets` folder that is a symlink is followed and validates. The spec decides the symlinked ticket file only. Recorded in the Edge inventory as Won't handle for reviewer veto.

## Falsification

Count: 5. Worst: the prose says the answer lives only in the ticket file, then says the gist states the answer.

- accept: the map index entry in `CONTEXT.md` and ADR 0020 say the index holds no answer. The gist entry says a gist states the answer. Reword the gist as a one-sentence summary and the index as holding no canonical answer. The anchored sentence stays.
- accept: the shape-idea Resume mode says every row names one ticket file. A fog row names the index. Reword the resume step for the fog case.
- accept: the write-spec move sentence names the topic folder, and the index is its sibling. Add that the index moves with the folder, in both commands and the retirement sentence. The anchored sentence stays.
- accept: the decision-map glossary entry assigns exclusions, discretion, and research objects to the decision tickets. The index owns those sections. Reword the entry.
- merge: duplicate gists, repaired under Coverage.
