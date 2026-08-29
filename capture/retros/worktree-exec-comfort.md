# Retro: worktree-exec-comfort

## Outcome

The spec landed at `ee0ed9d5` on `main` through `bench worktree land`. The
review base was `8f7f36af`, the source tip `e8d670e4`, and the assignment
start `37194769`. The gate ran green on the composed tree, and the landing
released the integration worktree with `census=174`. Six tickets, one review
round with two repair tickets, and one repair-scoped re-review carried 48
coverage rows to green.

The landing added `bench worktree show`, the `--env` flag, the stdin help,
the `worktree:` line, and the `next=` line. It also added the missing-tree
recovery route and the guard's exec exception with `segment=` and
`operator=`.

## Gate-stage timings

The landing gate at `ee0ed9d5`: gofmt 83 ms, vet 803 ms, test 47492 ms, race
2301 ms, system 14405 ms, shellcheck 447 ms. Five earlier folds ran the same
gate at 66 to 80 s wall each. Nine `bench commit` lanes ran at a few seconds
each; three of them refused on prose.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized; no charge took a spec slice. Six ticket charges
ran `opus` low (five) and medium (one). Four landed first-pass on behavior.
The guard ticket took three continuations: a fence extension for
`cmd/bench/main.go`, a repair for the named operator under a heredoc, and a
prose fix. The missing-tree ticket corrected the charge's landed predicate
inside its own scope and reported one out-of-fence expectation.

Two repair charges (`opus` low) landed first-pass. Three review axes and one
re-review (`opus` medium) returned 18 raw findings that collapsed to 11 repair
targets. The re-review verified all 8 predicates clean. Every delegate ran its
own mutation probe, and each probe bit.

## Coordinator catches

The coordinator probed every returned diff at a distinct site and kind, and
all nine probes bit. One throwaway probe found an unpinned case: a refusal
under a heredoc plus `2>&1` named the allowed heredoc. The coordinator caught
the guard ticket's `segment=` half that reached only the unit seam. The
coordinator's own spec prose broke the prose lane twice, and the T6 delegate's
prose broke it once.

The second repair worktree was created from a moved `main`. Its fold carried
five unrelated commits into the source, so the landing needed the moved
review base `8f7f36af`.

## Repair attribution

| ticket | repair rounds | causes |
| --- | --- | --- |
| add-child-argv-and-repeatable-flag-attributes | 1 | spec-row |
| pass-argv-stdin-and-env-to-the-exec-child | 0 | none |
| print-the-worktree-path-on-exec-failures | 0 | none |
| refuse-a-missing-tree-and-name-the-next-verb | 0 | none |
| classify-follow-ons-per-bench-segment | 3 | ticket-slicing, spec-row, delegate-error |
| add-the-worktree-show-verb | 1 | spec-row |
| repair-worktree-surfaces | 0 | none |
| repair-guard-and-inventories | 0 | none |

## Agent-experience improvements

### Bench CLI

- Give the exec child an ambient worktree-root variable, so a delegate never
  binds the path in a shell variable. The landing's census holds 127 `W=` and
  20 `cd` heads in the integration worktree alone.
  Feeds: new
- Make `bench test --check <unknown>` name the unknown check and list the
  registry names, and give the prose lane a reachable check name.
  Feeds: new
- Print the tip beside the path in `bench worktree path`, so a delegate's
  preamble is one call.
  Feeds: none

### Skills

- Add to `craft-delegate` a rule for a sibling repair worktree: create it
  before `main` moves, or fold it from the integration source only. A fold
  then never carries unrelated commits into a reviewed source.
  Feeds: none
- Add to `craft-spec` a rule that a row whose behavior names the child also
  names the child-level seam, so a parser-only test cannot close it.
  Feeds: none

### Process

- Run `bench gate-prose` on every prose file before a `bench commit`, which the
  prose-lane learning already records.
  Feeds: none
- Name every file a ticket's collapse can reach in its `Writes:` list,
  including the hook or route file that delivers a new value to stderr.
  Feeds: none
