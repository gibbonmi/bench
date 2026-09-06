# Review pickup: decision-map-split

Reviewed pair: base `1019f1f7`, tip `0a9a3b90`, on 2026-09-06. Raw findings: Standards 6, Spec 5, Coverage 3, Falsification 5. Eight repair targets landed in the two repair tickets. The six items below stay open for the reviewer, and none blocks the landing.

## Standards

Count: 2 open. Worst: the six-cell row order is derived at six sites.

- ask-user: the row order `{map,title,type,state,blockers,path}` is built at six sites in `maps.go` and `freshness.go` and read by position. A row struct with one render method would single-source it. Citation: the smell baseline, Primitive Obsession. Open, because it is a seam change across a landed fence.
- no-op: `var registry = append(generalAnchors, decisionMapAnchors...)` is sound today. `slices.Concat` would state the intent. A note only.

## Spec

Count: 3 open. Worst: two migrated gists do not state their decision.

- ask-user: the gist of `gate-critical-path` #1 reads `With lever 1 landed.`, and `worktree-orphan-retirement` #6 reads `Bound it.`. A gist is reviewer prose, so the reviewer edits or accepts them.
- ask-user: the `gate-pipeline-fixture-inventory.md` move into the map's assets folder is authorized by no story. It is recorded under the spec's Build decisions for assent or veto.
- no-op: the proof checklist names `ParseDecisionMap`, which the contract unexported, and DS33 names a test the build superseded. Recorded-state drift; the retro notes it.

## Coverage

Count: 1 open. Worst: a symlinked tickets folder is followed.

- ask-user: a `tickets` folder that is a symlink is followed and validates. The spec's Edge inventory now records it as Won't handle for reviewer veto.
