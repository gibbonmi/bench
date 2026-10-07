# FT335: full gate selected for a planning merge

Status: source-verified full-gate selection observed; merge completed green.

The maintenance planning merge selected the full prospective gate from the primary checkout.
The target is a Bench kit assignment with a declared selective lane.
This matches the remaining full-gate-selection face of FT335.
No new benchmark or second reproduction was requested or run.

## Source and executable identity

Repository: `/home/mgibs/workspace/bench`

Primary source: `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68`

Target assignment: `ft362-process-lifetime-map`

Target before merge: `a382f4848d432815968f044820fad593e856e297`

Incoming assignment: `maintenance-existing-specs`

Incoming tip: `b5dde6cacfbae3c9d3f07760c501d70fa798d84b`

Exact command from the primary checkout:

```text
bench worktree merge --from maintenance-existing-specs ft362-process-lifetime-map
```

The running merge process was PID 433776.
Its executable resolved through `/proc/433776/exe` to the primary `dist/bench`.
Its SHA256 matched the current binary and adjacent seal:

```text
e9c18acca938573ac9a25b6f0ea221f63ed0d7054a39c923434d71c8ed37edbb
```

The seal's source digest was:

```text
f9f8bc28c964a4ddca29a6c4f6724a3d2c058933e14ea74854f8265d1271c92d
```

`bench freshness-check /home/mgibs/workspace/bench` returned zero during the run.
The running process environment named the primary checkout as BENCH_KIT.
Only BENCH_KIT and BENCH_RUN_BINARY were inspected; no complete environment was copied.
The binary reports version 0.2.0 on linux/amd64.

## Observed execution

The process tree showed these direct parent-child relationships:

| PID | Parent | Command |
|---|---|---|
| 433776 | invoking app server | primary dist/bench worktree merge |
| 434674 | 433776 | private sealed bench gate-phases on a prospective checkout |
| 437884 | 434674 | go test -trimpath -count=1 ./... |

The ordinary all-package test run was active, including adopt, repairtest, conformance, and gate packages.
The merge was integrating reviewed Markdown and a planning scorecard, with no production source changes.
The merge completed at 60248085580469f76eb86b1d6cb9b63757b85106 with all six phases green.
This observation establishes check selection, not a prediction of performance savings.

## Source mechanism and roadmap overlap

`internal/worktree/merge.go:91` passes the caller's kit and the target worktree to mergeOwner.
`internal/worktree/merge.go:386` resolves LaneForCommitAtKit with that pair.
`internal/gate/kit_source.go:78` distinguishes equal root and kit paths from a linked-repository manifest path.
`internal/landing/lane_owner.go:30` uses the whole-project owner when the resolved lane is nil.

FT335 already records this cause class and requires a source-verified runtime reproduction.
Its caller-root grading defect is recorded as fixed; that separate face is not reopened here.
The current observation supplies the previously missing binary identity and full-gate-selection evidence.
Recommend reassessing FT335's parked status in the maintenance ranking.
Do not change the active FT290 commitment or implement FT335 in this planning batch.

The supported target-local merge invocation should preserve the target kit identity.
Its source route is the existing wrapper's same-repository worktree detection at bin/bench.sh:67.
The next authorized guard-plan integration used the target-local invocation and selected only its prose lane.
It merged at 5351a242ca277c4c1d05689d750f9c27dabdea44.
These integrations carried different Markdown changes, so they are not a controlled timing comparison.
No security or gate policy was relaxed.

## Retained evidence

The source assignment retains these gate records:

- `.logs/gate-20261007T020747.503155310Z-433776.jsonl`
- `.logs/gate-20261007T020747.503155310Z-433776.out`

The original merge's complete bounded response is retained outside the pool:

```text
/home/mgibs/.bench/responses/bench-2826441890/primary/1791339287016940782-b7a67569bc2ce30d.out
```

The record names six phases: gofmt, vet, test, race, system, and shellcheck.
The ordinary all-package and tagged system phases ran despite the planning-only delta.
Eight capability skips were reported: four FIFO and four privilege skips.
Those skips limit runtime qualification and do not affect the lane-selection observation.
