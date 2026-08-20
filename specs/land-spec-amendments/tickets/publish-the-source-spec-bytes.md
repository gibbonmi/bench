# Publish the source spec bytes

Blocked by: authorize-the-spec-folder-implicitly.md
Writes: internal/landing, internal/worktree

## What to build

`LandReviewed` publishes the implemented transition of the source tip's spec
bytes and stops demanding that the destination carry identical bytes. The
staged-spec check keeps its source half — the supplied bytes must be the
source tip's committed spec — and its refusal message names that source-tip
mismatch. The composition treats the spec path as source-owned: the
destination side receives the source's bytes at that path before the merge,
so a destination spec change — stale, amended, overlapping, or absent — can
never become a content conflict, and the transition then rewrites the path.
The resume path's source-only comparison stays as it is.

The end-to-end proof composes this ticket with the implicit authorization from
`authorize-the-spec-folder-implicitly.md`: a source carrying an in-range
spec-amendment commit, with no self-fence entry, lands green through
`bench worktree land`.

Covers: LS1, LS2, LS3, LS4, LS5, LS6, LS12, LS13.

## Acceptance

- [ ] A landing whose source-tip spec differs from the destination's copy publishes the source's transition (LS1).
- [ ] A landing whose supplied spec bytes are not the source tip's committed spec refuses, naming the source-tip mismatch (LS2).
- [ ] A destination whose spec changed after the review base still lands with the source's transition published (LS3).
- [ ] A source-tip spec that does not parse as staged refuses (LS4).
- [ ] A destination that never carried the spec file still lands with the source's transition published (LS5).
- [ ] A resumed landing over a published amended landing verifies and completes in a fresh process (LS6).
- [ ] `bench worktree land` over a source with an in-range spec-amendment commit and no self-fence lands green end-to-end (LS12).
- [ ] A destination spec change overlapping the source's amendment on the same lines still lands with the source's transition published (LS13).
