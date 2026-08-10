# Migrate status and dashboard aggregates

Blocked by: none
Ownership fence: `internal/status`, `internal/dashboard`
Integration surfaces: ordered aggregate carrier→`internal/axi/aggregate.go` exercised by SD1; signal producer and severity ladder→`internal/status/status.go` exercised by SD1; dashboard composition→`internal/dashboard/dashboard.go` exercised by SD1; dashboard section renderer→`internal/dashboard/render.go` exercised by SD1; empty-class sibling→migrate-status-dashboard-empties.md; legacy carrier contraction→contract-aggregate-empty-routes.md
Contracts: the ordered `[]status.Signal` ladder (severity, name, detail, action) and the composed dashboard section facts cross `internal/status/status.go`→`internal/axi/aggregate.go` and on to `internal/dashboard/dashboard.go`; severity is the ascending integer sort key, name/detail/action are strings, and a failed capture read renders `toon.UnknownCell` rather than `0`; order is the producer's severity order, then the seven `gather` sections gate, signals, roadmap, sequence, ideas, open-learnings, worktrees; a zero count is rendered rather than dropped, asserted by SD1 against the real `status.Signals` producer
Closure: SD1/signal-order, SD1/signal-budget, SD1/zero, SD1/unknown, SD1/section-gate, SD1/section-signals, SD1/section-roadmap, SD1/section-sequence, SD1/section-ideas, SD1/section-learnings, SD1/section-worktrees, SD1/route

## What to build

`bench status` and `bench dashboard` supply their already-derived signal and section facts
to the shared ordered aggregate carrier and render identical bytes. `status.Signals` keeps
ownership of the severity ladder, of the five-row budget and its `+N more` overflow line,
and of the `toon.UnknownCell` posture for a failed capture read. `dashboard.gather` keeps
ownership of section composition and re-parses no source: every section fact comes from the
existing reader that owns it.

The status and dashboard tests currently pin bytes without proving a shared route, so this
ticket adds the shared-route mutation the row lacks: a new
`TestDashboardComposesEveryOwnerSectionThroughTheAggregate` in `internal/dashboard` asserts
each section reached the page through the shared carrier holding the producer's own values,
so a locally recomputed section is red even when the HTML is unchanged.

This ticket is deliberately separate from `migrate-status-dashboard-empties.md`: the ordered
signal/section aggregates and the prose-clean and absent-refusal empty classes land green
independently, and no project-gate red spans the two.

Tree condition that must hold when this ticket is refreshed: `internal/axi/aggregate.go`
exists and declares the exported ordered-aggregate type `Aggregate` with its typed fact
entry `Fact`. If that path or either symbol is absent, stop and report rather than build —
the prerequisite `axi-carriers-and-registry` build has not landed.

## Acceptance

- [ ] [SD1] (covers AE7) status renders its severity-ordered signal facts, and the dashboard composes each of its seven owner sections from those same facts, through the shared aggregate carrier.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SD1/signal-order | in `SignalsWith`, return the rows in collection order instead of the ascending severity sort | `TestRenderDirtyLeadsGitOverDrainRow` (`internal/status`) | run `go test ./internal/status -run TestRenderDirtyLeadsGitOverDrainRow -count=1 -timeout 180s`; the assertion that the lead `▶` line names the git action fails with the drain row leading; the fixture is a `t.TempDir()` git repository whose only git calls are `add`/`commit`/`status`, so nothing waits on a network or a lock |
| SD1/signal-budget | in `render`, print every signal when `all` is false and drop the `+%d more` line | `TestStatusSignalBudgetKeepsTheOverflowCount` (`internal/status`, new) | run `go test ./internal/status -run TestStatusSignalBudgetKeepsTheOverflowCount -count=1 -timeout 180s`; the assertion that a seven-signal repository renders five rows plus `  +2 more (bench status --all)` fails with seven rows and no overflow line; same local `t.TempDir()` git fixture |
| SD1/zero | in `appendDrain`, omit the `%d idea(s)` segment from the degraded-source detail when the count is `0` | `TestStatusDrainRowKeepsExplicitZeroCounts` (`internal/status`, new) | run `go test ./internal/status -run TestStatusDrainRowKeepsExplicitZeroCounts -count=1 -timeout 180s`; the assertion that the drain detail still names `0` parked ideas beside a non-zero learnings count fails with the count absent; the fixture writes two small capture files, each read capped at `bounds.ControlRecordLimit` |
| SD1/unknown | in `appendDrain`, render `%d open learning(s)` for a `bounds.FileState` that `Failed()` instead of `toon.UnknownCell` | `TestStatusDrainRowKeepsExplicitZeroCounts` (`internal/status`, new) | run `go test ./internal/status -run TestStatusDrainRowKeepsExplicitZeroCounts -count=1 -timeout 180s`; the unreadable-journal case fails because the detail reads `0 open` where the unknown cell naming `capture/learnings.md` and its state is required; bounded by `bounds.ControlRecordLimit` |
| SD1/section-gate | in `gather`, leave `Snapshot.Gate` zero-valued instead of taking `status.GateVerdict(root)` | `TestDashboardComposesEveryOwnerSectionThroughTheAggregate` (`internal/dashboard`, new) | run `go test ./internal/dashboard -run TestDashboardComposesEveryOwnerSectionThroughTheAggregate -count=1 -timeout 180s`; the assertion that the gate section reports the verdict owner's state fails with the absent-cache rendering on a fixture that has a gate cache; `Render` is pure over the injected snapshot and the gate read is one bounded file read |
| SD1/section-signals | in `gather`, populate `Snapshot.Signals` from a locally rebuilt ladder instead of `status.Signals(root)` | `TestDashboardComposesEveryOwnerSectionThroughTheAggregate` (`internal/dashboard`, new) | run `go test ./internal/dashboard -run TestDashboardComposesEveryOwnerSectionThroughTheAggregate -count=1 -timeout 180s`; the assertion that the rendered signal rows equal `status.Signals`' own name/detail/action triples in order fails at the first differing row; bounded as above |
| SD1/section-roadmap | in `gather`, set `RoadmapPresent` from `text != ""` instead of the reader's presence flag | `TestDashboardComposesEveryOwnerSectionThroughTheAggregate` (`internal/dashboard`, new) | run `go test ./internal/dashboard -run TestDashboardComposesEveryOwnerSectionThroughTheAggregate -count=1 -timeout 180s`; the present-but-empty `ROADMAP.md` case fails by rendering the `No ROADMAP.md` section against a file that exists; the roadmap read is capped at `bounds.ControlRecordLimit` |
| SD1/section-sequence | in `gather`, drop the `roadmap.RecommendedSequence(text)` call and leave `Sequence` empty | `TestDashboardComposesEveryOwnerSectionThroughTheAggregate` (`internal/dashboard`, new) | run `go test ./internal/dashboard -run TestDashboardComposesEveryOwnerSectionThroughTheAggregate -count=1 -timeout 180s`; the assertion that the sequence section carries the roadmap owner's recommended sequence fails with an empty section; bounded as above |
| SD1/section-ideas | in `gather`, populate `Ideas` from the drain count instead of `roadmap.ParkedIdeas(root)` | `TestDashboardComposesEveryOwnerSectionThroughTheAggregate` (`internal/dashboard`, new) | run `go test ./internal/dashboard -run TestDashboardComposesEveryOwnerSectionThroughTheAggregate -count=1 -timeout 180s`; the assertion that each parked idea line renders verbatim fails with a numeric placeholder; `capture/IDEAS.md` is read once under `bounds.ControlRecordLimit` |
| SD1/section-learnings | in `gather`, set `OpenLearnings` from `len(snap.Ideas)` instead of `drain.OpenLearnings` | `TestDashboardComposesEveryOwnerSectionThroughTheAggregate` (`internal/dashboard`, new) | run `go test ./internal/dashboard -run TestDashboardComposesEveryOwnerSectionThroughTheAggregate -count=1 -timeout 180s`; the assertion that the open-learnings count equals `roadmap.DrainCounts`' own tally fails with the ideas count; bounded as above |
| SD1/section-worktrees | in `gather`, discard `werr` so a classify failure renders as an empty pool | `TestRenderWorktreeClassifyFailureIsVisible` (`internal/dashboard`) | run `go test ./internal/dashboard -run TestRenderWorktreeClassifyFailureIsVisible -count=1 -timeout 180s`; the assertion that the classify error renders as a visible error fails because the section shows the empty-pool message; `Render` is pure over the injected snapshot, so no git call runs |
| SD1/route | keep the pre-migration direct `Snapshot` field reads in `Render` and never construct the shared aggregate | `TestDashboardComposesEveryOwnerSectionThroughTheAggregate` (`internal/dashboard`, new) | run `go test ./internal/dashboard -run TestDashboardComposesEveryOwnerSectionThroughTheAggregate -count=1 -timeout 180s`; the assertion that the seven section facts were carried by `axi.Aggregate` from their owning readers fails with no aggregate observed, even though the rendered page bytes are unchanged; `Render` is pure over the injected snapshot |
