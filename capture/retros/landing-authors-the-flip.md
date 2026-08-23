# landing-authors-the-flip — implementation retro

## Outcome

The FT113 build landed on `main` as `148f3a68` through `bench worktree land
--spec landing-authors-the-flip`. The source was the `ft113-build` worktree on
base `0e17d428`, reviewed tip `6d2d2c12`, repair tip `e673e2d5`. The landing
published `Status: implemented`. Five tickets landed in five commits, two merge
commits, and one repair commit. The review round returned three raw findings
and one repair target.

The first landing attempt refused on a moved destination checkout. The second
attempt refused on a dirty destination. The third attempt published after the
reviewer's files were committed.

## Gate-stage timings

One prospective gate run on the integration source, measured from the gate log:

| stage | elapsed |
|---|---|
| gofmt | 0.1 s |
| vet | 1.6 s |
| test | 78.8 s |
| race | 4.8 s |
| system | 9.4 s |
| shellcheck | 0.4 s |

Each ticket commit paid one full run. The build paid eleven gate runs. Two were
red on the coordinator's own prose bounds. One was red on a `sed` path error.
One was red on the reviewer's decision-map prose.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized and carried its rows, its fence, and its probe.
Group C (`sonnet` / medium) landed first-pass with two attributed mutation
probes. Group D (`fable` / high) landed first-pass within the 152-line budget.
It trimmed two wrapped lines instead of a budget raise.

The commit-route ticket (`opus` / medium) landed first-pass and surfaced one
judgment call on a fixture it simplified. The spec-verb ticket (`opus` /
medium) landed first-pass. It reported a stale comment outside its fence
instead of a silent edit. The exit-3 ticket (`opus` / medium) landed
first-pass. It exported `sanitize.LineSafe` as the one predicate instead of a
second copy, and it reported the out-of-fence edit. It also reported a latent
paste defect in the landing verb's `next=` sanitizer.

Two tickets ran in parallel worktrees on disjoint files and merged without
conflict.

## Coordinator catches

- The guidance and exit-3 tickets ran in side worktrees. The coordinator merged
  their commits into the integration source with raw `git merge`, because no
  Bench verb moves a sibling assignment onto the retained source.
- `bench idea` from the repo root wrote to `main`'s checkout. The coordinator
  removed the line and re-ran the verb inside the worktree.
- `bench worktree path` prints a `~` prefix that `sed` does not expand. One
  command failed before the gate ran.
- The spec fence list did not cover the one-source seam the build needed. The
  coordinator extended it in the source and flagged the edit for veto.
- After the landing, `bench spec retire` printed the old generic `next:` line.
  The wrapper ran a `dist/bench` built before the landing. A rebuild through
  `scripts/go-build.sh` showed the new surfaces, and FC1 had passed in the gate.
- The first landing attempt refused because another session wrote two files on
  `main` during the gate. The reviewer owned them. The coordinator split eleven
  over-bound sentences in the decision map so that their commit could gate.

## Repair attribution

| ticket | repair rounds | causes |
|---|---|---|
| retire-names-the-board-remainder | 0 | none |
| guidance-names-one-author | 0 | none |
| retire-the-commit-route-flip-and-close | 1 | other |
| retire-bench-spec-implemented | 0 | none |
| commit-exit-3-names-the-remainder | 0 | none |

The one round was a test comment that narrated the prior behavior. The
Standards axis caught it.

## Agent-experience improvements

### Bench CLI

- Give `bench worktree` a verb that merges one assignment's commits onto another retained assignment, so a side ticket folds without raw Git.
  Feeds: FT238
- Print an absolute path from `bench worktree path`, so a shell substitution works without `~` expansion.
  Feeds: new
- Name the moved paths in the `landing destination checkout changed` refusal, so the caller sees the foreign writer at once.
  Feeds: new
- Rebuild the repo-local `dist/bench` in the landing's publication step, so the next verb on the destination runs the landed code.
  Feeds: new

### Skills

- Add to `craft-delegate` the rule that a delegate drives a worktree through `bench worktree exec <label> --`, so its process note is not a deviation.
  Feeds: none

### Process

- Run `bench idea` inside the active worktree during a build, so `main`'s checkout receives no write outside a landing.
  Feeds: none
