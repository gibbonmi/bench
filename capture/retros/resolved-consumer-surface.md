# Retro: resolved-consumer-surface

## Outcome

`bench consumers` landed as commit `41b329f9` on `main` (FT191). The reviewed pair was base `976b26e6` and source tip `cd5beafb`. The landing gate was green in all six phases. The surface resolves real Go reference edges and classifies each row as `call`, `reference`, or `implements`. Each row names its enclosing declaration. A bare ambiguous name answers with a candidates table and one re-query action per row.

Every success response ends with a `citation[1]{sha,state,version,cmd,hash}` row before the help envelope. Three unsound inputs refuse at exit 1. The `--changed` form enumerates the consumers of every declaration a frozen pair touched and marks each row `touched`. It offers one `--full` walk action per symbol with an outside-diff consumer. `bench outline` gained the `helper`, `double`, and `fixture` kinds, and the review command runs the blast step before axis dispatch. The build's own blast at the final tip listed 641 rows over 250 changed declarations, with four symbols outside the diff.

## Gate-stage timings

The landing gate at tree `02b7aca3`: `gofmt` 102 ms, `vet` 833 ms, `test` 49,665 ms, `race` 2,313 ms, `system` 15,842 ms, `shellcheck` 503 ms. The first landing attempt went red in `test` after 10 s on the `bounds-policy` registry. The retry was green. Each worktree commit ran the four-check lane in well under a minute.

## Ticket-versus-spec-slice and delegate performance

Every charge was one ticket at `opus` / low. No charge needed the pre-approved `fable` / low rung. Eight ticket charges, four repair charges, and three continuations returned diff-ready with focused checks green. Seven of the eight ticket charges showed each coverage row red before the edit. The blast ticket wrote its tests and its code in one pass and reported that. Both behavior gaps the review found sat in that ticket.

The three codex review axes at `gpt-5.6-sol` / high returned 28 raw findings. They collapsed to 10 repair targets and 6 reviewer decisions. The repair-scoped re-review found one real gap the first repair had missed: git C-quotes control bytes in patch headers whatever `core.quotePath` says. It also found two stale spec sentences. The final scoped re-check found one unpinned decoder arm, and the coordinator graded its mutation directly.

## Coordinator catches

- The core delegate's `Origin()` identity was unpinned: a swap probe stayed green, and the repair added the generic plant.
- The core loader loaded a package with tests twice, so `Resolve` returned two matches for one declaration. The delegate had flagged it, and the coordinator folded the fix.
- The blast delegate's real run cited a synthesized `_testmain.go` path outside the root. The repair dropped out-of-root files at the loader.
- Three registries surfaced late. The delegate found the AXI envelope cases, preflight found the `internal/axi` fence, and the landing gate found the `bounds-policy` consumer line.
- The wrapper on `PATH` served the main checkout's stale executable, so live demonstrations needed the worktree's own build.

## Repair attribution

| ticket | repair rounds | cause per round |
|---|---|---|
| add-consumers-reference-core | 1 | spec-row |
| extend-outline-inventory | 1 | spec-row |
| add-consumers-command-surface | 1 | spec-row |
| classify-via-and-candidates | 1 | other |
| emit-citation-row | 1 | delegate-error |
| refuse-unsound-inputs | 1 | delegate-error |
| add-blast-mode | 4 | delegate-error, delegate-error, spec-row, delegate-error |
| wire-review-blast-step | 0 | none |
| repair-frozen-pair-refusals | 0 | none |

## Agent-experience improvements

### Bench CLI

- The landing's census entry is `Assignment rcs-integration: 338 raw calls` in `capture/learnings.md`.
  Feeds: new
- Give `bench structure` a path argument, so a delegate confirms its fence added no issue in one line.
  Feeds: new
- Make the follow-on guard name the segment it refused, because a delegate reads the refusal as the `bench` call's own fault.
  Feeds: new
- Add `bench worktree build <label>`, or document `bench worktree exec <label> -- ./dist/bench`, because the wrapper on `PATH` serves the main checkout's grammar.
  Feeds: new
- Keep the census record readable until the tail writes its learning entry, because the landing sweeps the store first.
  Feeds: new

### Skills

- In `craft-delegate`, send a ticket that returns without a pre-edit red per row back for the reds before it commits.
  Feeds: new
- In `craft-tickets`, require a charge that exports from a package outside its fence to name the fence amendment in its return.
  Feeds: none

### Process

- Before a landing, run `bench test --check` for every conformance registry that names a file the build moved.
  Feeds: none
- Run every ticket of one spec on the retained integration source in `Blocked by:` order, because a copied-over diff leaves a dirty registration.
  Feeds: none
- Treat a landing's "changes the promotion broker source" line as a tail duty: rebuild `dist/bench` and run `bench doctor --fix`.
  Feeds: new
