# Benign Stale Gate Status

Status: staged

## Problem

`bench status` reads the last gate cache and compares it to the current working
tree. Any mismatch currently renders the same stale-gate action:
`re-run the gate`. That is correct for code, docs, config, and mixed drift, but
it is too strong when the only drift is capture-only work in `ROADMAP.md` or
`.bench-notes.md`. Those files record parked ideas and shift scratch notes; they
do not change what the gate just proved.

## Solution

Keep the existing stale-gate detection, but classify the path diff between the
cached gated tree and the current tree. If every changed path is exactly
`ROADMAP.md` or `.bench-notes.md`, render the gate row as signal `gate`, detail
`stale (capture-only drift)`, action `re-run when convenient`. If any other path
changed, if the diff is mixed, or if the cache/tree comparison cannot be trusted,
keep the existing stronger stale row and action.

## User stories

1. As an agent starting a session after only `ROADMAP.md` changed, I want
   `bench status` to show `gate`, `stale (capture-only drift)`, and
   `re-run when convenient`, so parked ideas do not look like code drift.
   Line: claude-sonnet-5 / low. This is a known status seam with exact runtime
   contract coverage.
2. As an agent starting a session after only `.bench-notes.md` changed, I want the
   same capture-only stale row, so shift scratch notes get the same softer signal.
   Line: claude-sonnet-5 / low. The behavior shares the same helper and fixture
   shape as story 1.
3. As the gate owner, I want changes to code, docs, config, tests, or any path
   outside the fixed allowlist to keep the existing stale row and
   `re-run the gate`, so real unproven drift remains loud.
   Line: claude-sonnet-5 / low. The status renderer already owns the signal; the
   change is a branch under the same contract.
4. As the gate owner, I want mixed drift involving a capture-only path plus any
   other path to keep the stronger stale row, so benign capture cannot mask real
   drift.
   Line: claude-sonnet-5 / low. This is an edge case under the same path classifier.
5. As the gate owner, I want malformed cache data, missing tree objects, or diff
   failures to fall back to the stronger stale row, so an untrusted comparison
   never weakens the gate signal.
   Line: claude-opus-4-8 / medium. This is the fail-closed branch of the oracle
   path and deserves gate-level judgment.
6. As the maintainer, I want the path comparison rooted in the same tree hashes the
   gate cache already uses, so the status classification has one source of truth
   for what changed.
   Line: claude-opus-4-8 / medium. This adds or exposes a git helper at the cache
   key seam rather than deriving a second tree model.

## Implementation decisions

- `internal/status` remains the owner of dashboard classification and wording.
  `appendGate` still reads `<git-dir>/bench-last-gate`, computes
  `git.TreeHash(root)`, and compares the cached tree with the current tree.
- Add a small `internal/git` helper for tree-to-tree changed paths, for example
  `ChangedPathsBetweenTrees(root, fromTree, toTree) ([]string, bool)`. It should
  shell out to git with root-relative tree objects, return sorted or git-ordered
  relative paths, and return `ok=false` when either tree is invalid or the diff
  command fails.
- The capture-only allowlist is fixed and exact:
  `ROADMAP.md` and `.bench-notes.md`. No directory prefix, suffix, extension, or
  markdown class counts as benign.
- A stale gate row is benign only when the changed-path set is non-empty and every
  changed path is in the allowlist. Added, modified, and deleted allowlisted files
  all count as capture-only because the tree diff path is the relevant fact.
- Any mixed diff or untrusted diff keeps the current detail shape:
  `stale (gated tree <short>, work tree <short>)` and action `re-run the gate`.
  Do not change the gate cache key or stop-hook write format.
- Keep the stale-gate severity at the existing gate-stale severity so the
  dashboard ordering stays stable; only the detail/action text changes for the
  benign subset.

## Testing decisions

- The highest seam is the `bench status` runtime contract in
  `internal/contract/runtime_test.go`, because the SessionStart hook consumes that
  exact renderer. Add fixture cases that write a green cache for the base tree,
  mutate paths, run `bench status`, and assert the exact row text.
- Add lower-level tests for the git helper only if the implementation needs
  meaningful parsing beyond a direct `git diff --name-only` call. The runtime
  contract is required either way because exact status wording is the product.
- Current red evidence for story 1:

  ```text
  ROADMAP.md-only drift currently renders:
    gate       stale (gated tree ..., work tree ...) -> re-run the gate
  ```

  The desired behavior is the capture-only detail/action instead.

### Seam diagram

Seam 1 - status renderer (user-visible contract):

    trigger: bench status
        |
        v
    bench-last-gate cached tree + current git.TreeHash(root)
        -> [ internal/status.appendGate path classifier ]
        -> dashboard row
              tests attach here with runtime fixtures asserting exact text

Seam 2 - git tree diff helper (single source for changed paths):

    trigger: appendGate sees cached tree != current tree
        |
        v
    cached tree object + current tree object
        -> [ internal/git tree-to-tree changed-path helper ]
        -> root-relative changed paths or ok=false
              tests attach here only if helper parsing grows past direct git output

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `ROADMAP.md`-only drift renders signal `gate`, detail `stale (capture-only drift)`, action `re-run when convenient` | `bench status` runtime contract | already red: a throwaway repo with only `ROADMAP.md` drift renders `stale (gated tree ..., work tree ...)` and `re-run the gate` | exact output proves the current status branch cannot distinguish capture-only roadmap drift |
| 2 | `.bench-notes.md`-only drift renders the same capture-only row | `bench status` runtime contract | new test before implementation mutates only `.bench-notes.md` and expects the capture-only row | the existing implementation treats every tree mismatch the same, so this expectation fails until classification exists |
| 3 | any non-allowlisted path keeps `stale (gated tree <short>, work tree <short>)` and `re-run the gate` | `bench status` runtime contract | new test mutates `docs/x.md` and asserts the strong row | a too-broad allowlist or markdown-class shortcut would soften real drift and fail this assertion |
| 4 | mixed allowlisted plus non-allowlisted drift keeps the stronger stale row | `bench status` runtime contract | new test mutates `ROADMAP.md` and another file and asserts `re-run the gate` | the classifier must require every path to be allowlisted, not any path |
| 5 | malformed cache, missing tree object, or diff failure falls back to the stronger stale row | `bench status` runtime contract | existing deadbeef stale-gate test stays strong; add malformed-cache coverage if implementation changes parsing | an uncertain comparison cannot produce the softer action |
| 6 | status classifies paths from the same tree hashes used by the gate cache | status plus git helper | helper/unit test or runtime fixture changes an untracked allowlisted file after a cached tree | `git.TreeHash` includes tracked plus untracked unignored files, so the classifier must inspect that tree rather than only `git diff HEAD` |

### Edge inventory

- Added `ROADMAP.md` after a green cache: benign if it is the only changed path.
  Covered by story 1.
- Modified `ROADMAP.md`: benign if it is the only changed path. Covered by story 1.
- Deleted `ROADMAP.md`: benign if it is the only changed path. Covered by story 1.
- Added, modified, or deleted `.bench-notes.md`: benign if it is the only changed
  path. Covered by story 2.
- `docs/ROADMAP.md`, `ROADMAP.md.bak`, `.bench/notes.md`, or any nested/similar
  path: not benign. Covered by story 3's non-allowlisted path test.
- Mixed capture-only and real drift: not benign. Covered by story 4.
- Cache line with missing fields or `ctree == "none"`: keep the stronger stale row.
  Covered by story 5.
- Cached tree object pruned or otherwise unavailable: keep the stronger stale row.
  Covered by story 5.
- Current `git.TreeHash(root)` returns `"none"`: keep the stronger stale row.
  Covered by story 5.
- ROADMAP footer still renders independently when parked ideas exist; the gate row
  must not become the footer or hide it. Existing footer contract remains in force.
- Worktree dirty/unpushed row still reports separately. Capture-only stale status
  does not suppress the git signal.

## Out of scope

- **Configurable or broader capture-only allowlist.** The decision map explicitly
  chose the fixed exact paths. Expanding it would need a new decision because
  `.bench/learnings.md`, specs, docs, and other markdown can affect active workflow
  or gate-observed behavior.
- **Changing the gate cache key or stop-hook format.** This spec only changes how
  status explains a stale cached verdict after comparing the same trees it already
  compares.
