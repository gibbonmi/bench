# Retro: git-admin-readers

## Outcome

The landing `5589e73a` published the spec on 2026-09-05 over the frozen pair `9d1a59e0` to
`9ed89456`. The Git adapter now exports `AdminDir` and `AdminPath`, and `CommonDir` runs
the same bounded, validated path. Twelve directory sites and five file sites outside the
adapter call a reader, and the `git-plumbing-owner` check with its `git-flag-retyped`
fixture guards the seam. The stub harness gained nine modes, and the glossary gained
**checkout administration directory**.

The build found that `--path-format=absolute`
resolves an existing symlink, so the file reader uses Git's default format and joins onto
the root that `canonicalpath.Resolve` answers. `Resolve` now resolves symlinks before it
absolutizes. The spec stays `implemented` as the veto surface for five recorded decisions.

## Gate-stage timings

| stage | landing gate |
| --- | --- |
| gofmt | 95 ms |
| vet | 970 ms |
| test | 68044 ms |
| race | 5553 ms |
| system | 24942 ms |
| shellcheck | 558 ms |

The two fold gates ran the test stage in 65381 ms and 65905 ms. One landing attempt ran
red on the `canonical-path-owner` check and cost one full gate.

## Ticket-versus-spec-slice and delegate performance

Seven ticket charges ran on `sonnet`: five at low and two at medium. Three landed
first-pass on behavior (the adapter, the gate sites, the diff site). The worktree
directory ticket raised the serial ceiling outside its fence and left one guard.

The worktree file ticket reported the symlink seam instead of an edit outside its fence, and
`opus` at low repaired the adapter in one pass. The adopt ticket grew two over-budget
files and widened a comment to 258 characters before it trimmed honestly. The check ticket
grew one over-budget line and ranged its GR27 test over the production flag list. The
review repair on `sonnet` at medium closed four rows but re-derived the canonical path
beside its owner, and the landing gate caught it.

## Coordinator catches

- A GR27 test that ranged over the production flag list stayed green through an omitted
  flag. The coordinator's omission probe was silently green. The test now enumerates
  the four flags itself.
- A comment rewrap to 258-character lines to hold a line budget; the trim was redone at
  the file width.
- A serial ceiling bump outside the ticket fence; accepted as a fence extension.
- A repair that re-derived the canonical path with `filepath.Abs` and `EvalSymlinks`; the
  landing gate reported it, and the fix moved into the owner.
- Two preflight reds from fence spelling: a fixture closure entry with a trailing slash,
  and a deleted pickup path left in a Writes line.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| add-the-named-git-admin-readers | 1 | spec-row |
| migrate-the-worktree-directory-sites | 1 | ticket-slicing |
| migrate-the-worktree-file-sites | 0 | none |
| migrate-the-gate-and-dashboard-sites | 0 | none |
| migrate-the-diff-index-identity | 1 | spec-row |
| refuse-an-unresolved-hooks-directory | 3 | ticket-slicing, delegate-error, spec-row |
| add-the-git-plumbing-owner-check | 2 | ticket-slicing, delegate-error |
| repair-review-round-1 | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- Add a probe verb that copies a file aside, applies one named mutation, runs a focused
  test, and restores the bytes with `cmp`. The census entry `git-admin-readers census 34`
  proposes it.
  Feeds: new
- Give `bench gate-prose --help` a usage example with the `. -- <file>` form, because the
  bare file form errors with `root is not a directory`.
  Feeds: none

### Skills

- Make the spec fence rule in `craft-spec` name the serial ceiling file whenever a row
  binds `PATH` in a package with a pinned serial set.
  Feeds: new
- Make `craft-tickets` check each test file a ticket names against the structure budget,
  so a new test lands in a file with room.
  Feeds: new

### Process

- Run the whole-tree gate after a repair commit and before the landing, because a
  lane-only repair reached the landing gate red.
  Feeds: none
- Make a spec that names a Git flag cite an observed run of that flag over the hostile
  shapes. The absolute path format's symlink resolution was an assumption.
  Feeds: new
