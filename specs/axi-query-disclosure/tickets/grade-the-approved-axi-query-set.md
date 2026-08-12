# Grade the registry-declared AXI query set

Blocked by: disclose-learnings-actions.md, disclose-map-actions.md, render-honest-anchor-help.md, disclose-guard-actions.md, disclose-coverage-actions.md, route-worktree-list-through-shared-disclosure.md
Writes: `internal/conformance/`, `cmd/bench/main.go`, `cmd/bench/command_registry.go`, `cmd/bench/command_registry_test.go`, `projects/benchkit.md`

## What to build

Declare AXI scope on the production registry entries, including nested `worktree list`, and add one registry-derived conformance check over the complete approved set. Keep the approved-set expectation independently authored so membership flips fail in both directions, and grade real command output for every exact predicate rather than treating the family as one sampled case. Update the project seam advertisement from the same production truth.

## Acceptance

- [ ] [AC1] (covers QD3) production registry declarations identify exactly `anchors`, `learnings`, `maps`, `guards`, `diff`, `coverage`, and `worktree list`, with every other member explicitly exempted and nested membership derived rather than guessed.
- [ ] [AC2] (covers QD3) every declared member is graded independently for structured success stdout, definitive zero-row empty, structured refusal on stdout/1, unknown-flag usage/2, identical deep-cwd routing, all three help spellings on exit 0, and an appended or honest-empty help envelope.
- [ ] [AC3] (covers QD3) removing one approved declaration or adding one operational declaration makes conformance red against an independent expectation; package-main tests bind registry declarations to real command envelopes.
- [ ] [AC4] (covers QD3) `projects/benchkit.md` advertises the same approved set and no wider operational family.
