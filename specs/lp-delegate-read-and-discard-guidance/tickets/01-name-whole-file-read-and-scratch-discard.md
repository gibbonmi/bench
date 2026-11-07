# Name the whole-file read and the scratch-worktree discard in delegate guidance

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md
Covers: none

## What to build

A research delegate that reads a spec through `bench worktree read` gets only the first page, and no guidance names a whole-file form. The delegation discipline names the whole-file forms that a charge gives its delegate.

A scratch-patch research delegate cannot release its unlanded worktree, and the slicing checks name no discard command. The slicing checks name the `bench worktree clean` discard plan and its fingerprint apply.

## Acceptance

- [ ] The "In the charge" section of the delegation discipline names the file reader on the worktree path, `bench worktree show <target> HEAD:<path>`, and each `bench worktree read` continuation as the whole-file forms.
- [ ] The scratch-patch paragraph of the slicing checks names `bench worktree clean --discard-branch [--discard-ignored] <path>` and the rerun with `--apply <fingerprint>`.
- [ ] Each addition is one or two sentences, uses ASD-STE100 prose, and matches the grammar that `bench worktree --help` prints.
