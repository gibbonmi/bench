# Route agent Markdown through per-component gate scoping

Blocked by: none

## What to build

Markdown-only changes anywhere under `.agents/` skip the Go toolchain gate
components through the per-component input declarations, while the contract and
canary components — whose fixtures consume the kit's guidance tree — declare
that Markdown as an input and run. The whole-changeset reduced scope does not
cover `.agents/`: its stripped-worktree enforcement runs the excludable phases
in a tree without the declared paths, and the contract and canary consumers go
red there, which is the refusal that rerouted this ticket (reviewer decision,
2026-08-02).

## Acceptance

- [ ] A changeset confined to `.agents/**/*.md` skips `gofmt`, `vet`, `test`,
      and `race`, and runs `contract` and `canary`.
- [ ] A non-Markdown file under `.agents/` reaches no Markdown derivation: the
      contract component's `agent-markdown` input enumerates `.md` files only.
- [ ] Go edits continue to run Go phases without running canary, and canary
      input edits continue to run canary.
- [ ] `ReducedScope()` does not cover `.agents/`, and the stripped-worktree
      construction strips no `.agents/` path.
