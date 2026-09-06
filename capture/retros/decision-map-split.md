# Retro: decision-map-split

## Outcome

The spec landed on `main` at `21de019a` on 2026-09-06 from the reviewed pair `1019f1f7` to `fa868074`, with `Status: implemented`. Fourteen maps are split into an index plus ticket files, ten assets moved into per-map assets folders, and one orphan moved to `docs/research/`. `bench maps` takes a map operand, carries a `path` column, and projects ready and stale rows. The two commands, the glossary, the README, the field guide, and ADR 0020 record the shape.

Story 47 asked for the resume measure. Before the split, a session that resumed one map read the whole file. The fourteen maps held 4389 lines, and the largest held 1011. After the split, the same resume reads the index plus one ticket. The fourteen indexes hold 1117 lines, and a frontier ticket holds at most 20.

## Gate-stage timings

The landing gate at `21de019a` ran green in these stages.

| stage | elapsed |
| --- | --- |
| gofmt | 101 ms |
| vet | 923 ms |
| test | 68.1 s |
| race | 2.4 s |
| system | 24.2 s |
| shellcheck | 515 ms |

The whole-tree gate on the source at `fa868074` ran the same six stages green, with the test stage at 66.6 s. Six fold gates ran green before that.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized, and no delegate took a spec slice. Seven ticket charges and two repair charges ran on opus. Five ran at medium and four at high, for the two guidance tickets and the prose repair. Every charge landed first-pass on behavior, and every delegate probe bit. Four charges needed one continuation each for a budget or fence placement, and none for behavior. The three review axes and the repair-scoped re-review ran on opus, and the falsification pass ran on gpt-5.6-terra through `codex exec`.

## Coordinator catches

- The parse ticket grew an over-budget status test file; the continuation moved the test into a fenced file under budget.
- The contract ticket wrote a helper file outside the spec fence; the continuation moved the helpers into a fenced file.
- The projection ticket found the ready row's help action unreachable inside the fence. The coordinator extended the fence to the axi action seam and recorded it.
- The commands ticket grew both over-budget registry files with no grant; the coordinator routed the rows into a second registry file and recorded it.
- The commands ticket left the family owner set without the new file; the continuation added it.
- The coordinator ran a fold gate beside a live delegate test run, and one handoff test went red; the retry after the return was green.
- The coordinator rebuilt the primary binary before the landing and broke the broker digest; `bench doctor --fix` repaired it.

## Repair attribution

| ticket | repair rounds | causes |
| --- | --- | --- |
| parse-ticket-files-beside-the-inline-shape | 0 | none |
| split-every-map-with-the-migration-program | 0 | none |
| rebuild-the-integrity-fixture-family-on-the-split-shape | 0 | none |
| refuse-the-inline-shape-and-render-both-templates | 0 | none |
| project-ready-stale-and-path-rows-with-the-map-argument | 1 | spec-row |
| rewrite-the-two-commands-and-register-the-anchors | 1 | tree-drift |
| record-the-split-in-context-the-adr-and-the-public-docs | 1 | spec-row |
| repair-the-gist-target-the-duplicate-gist-and-the-review-standards | 0 | none |
| repair-the-guidance-prose-the-falsification-pass-found | 0 | none |

The projection round covered the axi fence extension. The commands round covered the registry headroom. The docs round covered the README anchor the review added as DS53.

## Agent-experience improvements

### Bench CLI

- Record the landing census entry `decision-map-split census: 14 raw calls` from `capture/learnings.md`, which proposes `bench worktree reset <label>` and a `bench commit --dry-run --changed` form.
  Feeds: new
- Make `bench probe --check <name>` name the tests it selected and echo the check name, because a silent verdict there reads as a coverage hole.
  Feeds: new
- Let `bench gate-prose` take a file operand and print the sentence starts of a long paragraph, because every prose red here cost a hand count.
  Feeds: new

### Skills

- Add to `craft-delegate` that a charge names the headroom route for an over-budget file, because two tickets stopped on the growth lane.
  Feeds: new
- Add to `craft-review` that a review-owned row names the artifact that stores its evidence, because the migration evidence lived only in a delegate return.
  Feeds: new

### Process

- Run no fold and no whole-tree gate while a delegate test run is live; the scorecard decision exists and the coordinator broke it once here.
  Feeds: none
- Let the landing rebuild the broker; a hand rebuild of the primary binary before the landing breaks the manifest digest.
  Feeds: none
- Keep spec 1's retirement until `craft-research-skill` lands, because both specs share the compiled decision folder.
  Feeds: none
