# Build the exact landing owner

Blocked by: none
Ownership fence: `internal/landing`, `internal/spec/spec.go`, `internal/spec/spec_test.go`
Integration surfaces: exact-tree authorization→existing `internal/gate.ExecuteTree` asserted by EL4; isolated-index composition→existing `internal/git` Git invocation surface asserted by EL1-EL3; staged-to-implemented byte transform→`internal/spec/spec.go` asserted by EL5; landing request/result contract→adopt-exact-landing-in-commit.md
Contracts: repository root, destination ref, expected base commit, message, literal root-relative paths, and optional spec transition cross `internal/commit`→`internal/landing`, asserted by EL1-EL9 against real Git; staged spec bytes cross `internal/spec/spec.go`→`internal/landing`, asserted by EL5 against the real transformer; authorization result crosses `internal/gate`→`internal/landing`, asserted by EL4 against the real prospective gate surface

## What to build

Provide one reusable owner that composes an immutable tree from an expected base plus only literal attributed working-copy paths and an optional staged-spec transition, authorizes that exact tree, publishes a single-parent commit by expected-old ref update, and reconciles only owned checkout state.

## Acceptance

- [ ] [EL1] Real Git composition covers additions, modifications, deletions, staged deletions, rename halves, named and deleted directory descendants, duplicate names, literal hostile names, symlinks without traversal, gitlinks as atomic entries, and whole-segment sibling exclusion without reading unnamed real-index entries.
- [ ] [EL2] Unknown, repository-escaping, unreadable, repository-root, and directly or transitively discovered special-file paths refuse before authorization; an absent path in worktree, index, and expected base refuses too.
- [ ] [EL3] A composed tree equal to the expected base returns `nothing to commit` before authorization even when the real index carries unrelated staged content.
- [ ] [EL4] The owner sends the immutable tree through the exact prospective gate surface and refuses red, unavailable subject or oracle identity, and incomplete outcomes without changing the ref, real index, or working-copy bytes.
- [ ] [EL5] The optional lifecycle transition derives implemented bytes from the real staged-spec transformer inside composition; authorization sees those bytes, and refusal or CAS loss leaves the invoking spec staged.
- [ ] [EL6] The landing commit tree equals the authorized tree, has the expected base as its sole parent, and publishes only by compare-and-swap of the resolved branch ref or detached `HEAD`.
- [ ] [EL7] Injected destination movement after authorization preserves the winning tip and losing attributed work, with no retry, rebase, merge, or project-green update.
- [ ] [EL8] Successful reconciliation cleans only named entries and the optional transitioned spec while preserving every unnamed staged blob and unnamed unstaged, untracked, and ignored byte and classification.
- [ ] [EL9] A reconciliation failure after successful publication reports landed-but-checkout-incomplete and never rolls the destination ref backward.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EL1 | seed composition from the real index or replace whole-segment matching with raw prefix matching | landing owner Git-object tests | stage foreign content and add a prefix-sharing sibling, compose named paths, inspect the tree and expect exclusion |
| EL2 | skip the pre-authorization path classifier for one hostile class | landing owner path table with authorization tally | name the hostile path, invoke landing, expect refusal and zero authorization calls |
| EL3 | determine emptiness from `git diff --cached` on the real index | landing owner empty-composition test | stage only foreign content, name a clean path, invoke landing, expect `nothing to commit` and zero authorization calls |
| EL4 | authorize the ambient checkout or treat a non-exact outcome as green | landing owner authorization table plus prospective gate integration | return each failure outcome or inspect the gate checkout, invoke landing, expect refusal and unchanged fingerprints |
| EL5 | apply the staged-to-implemented rewrite after authorization | landing owner spec-transition test | authorize through a gate double that reads the tree, expect implemented bytes there and staged bytes after injected CAS loss |
| EL6 | create the commit with the current tip as parent or publish with an unconditional update | landing owner Git-object test | move the destination around publication, inspect parent and ref, expect expected-base ancestry and CAS refusal |
| EL7 | replace the expected-old ref update with an unconditional destination update | injected ref-updater test | advance the destination immediately before publication, expect the winner to remain tip and loser bytes to remain available |
| EL8 | reconcile with a whole-index reset or checkout-wide read-tree | landing owner fingerprint test | capture unnamed index blobs and filesystem classifications, publish named work, expect exact pre/post equality for unnamed state |
| EL9 | report a post-CAS reconciliation error as an unpublished failure or roll back the ref | landing owner post-publication fault test | inject reconciliation failure after CAS, expect landed identity in the error and the published commit to remain tip |
