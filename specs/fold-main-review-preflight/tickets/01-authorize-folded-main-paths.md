# Authorize the paths that a folded main merge brings in

Blocked by: none
Writes: internal/preflight/, internal/worktree/, internal/reviewrecord/coverage.go, internal/gate/review_checkpoint_commits_test.go, .agents/commands/bench-review-implementation.md, .agents/commands/bench-implement-spec.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go
Covers: none

## What to build

A spec build can see `main` move after its first chunk. The completion gate then refuses the landing until the source folds `main`, because the composed tree must equal the reviewed source. After the source merges `main`, the review preflight refuses the merge range, because `paths-authorized` reads each merged path as an unfenced write. No current sequence lands such a build.

Make `paths-authorized` authorize a changed path when its bytes and mode at the source tip equal the bytes and mode at the default-branch tip. Such a path came in unchanged from `main`, so it is not a write of the build. A path that the build changed keeps the fence rule. Compose the existing seams. The landing's folded-base route in the worktree package is prior art for the default-branch tip lookup. Keep one source for that lookup, and do not write a second one.

Correct the review command so that it states the current route. When `main` moves during a build, the source folds `main` before the completion landing. The fold joins the review delta of the last chunk. The review preflight authorizes the paths that the fold brings in unchanged. Move the anchor needle in the chunk-chain registry with the sentence it pins.

The landing's folded-base refusal route predates this rule. It advises a later `--base`, and that narrows the range the landing authorizes. Now that folded paths authorize, retire or correct that advice, and update its test to the new behavior. The implement command's "Land" section and the review chain's gap message also say that a `main` merge lands only before the first chunk. State the fold route there too, and move each pin with its text.

## Acceptance

- [ ] A review preflight over a range that merges `main` authorizes each path whose source-tip bytes equal the default-branch tip bytes.
- [ ] A path that the build changed and that no fence covers stays red, even when the source also merged `main`.
- [ ] A path that `main` changed and that the build then changed again stays red unless a fence covers it.
- [ ] The review command and its anchor state the fold-before-completion route, and the anchor test passes.
- [ ] The landing gives no `--base` advice that narrows the authorized range for a folded source, and its folded-base test grades the new behavior.
- [ ] The implement command's "Land" section, the review chain gap message, and their pins state the fold route.
- [ ] Dogfood: with this binary, `bench preflight review ft290-test-projection --charge --base 6fdd27b3fc3ce791931004ad4202173766c1d66d --source-tip 742624291d04978a09dcabec2042eec394ed47f3` in the `ft290-test-projection` worktree goes green.
