# Give the response store root one source

Blocked by: none
Writes: internal/responsebound/store.go, internal/responsebound/retire.go, internal/responsebound/retain_test.go
Covers: none

## What to build

One function in the store file owns the repository directory of the response store. That directory is the absolute Bench home, the `responses` level, and the repo key. A root outside any repository takes the `none` key.

The spill writer, the retained response, and the spill drop all start their paths from that function. The drop does not derive the levels or the key itself.

The drop of an empty root changes. The drop names the `none` key, as the spill writer does. Before, the drop asked git for the repository of the working directory. A drop from inside a repository then removed an assignment scope of that repository. The worktree retirement always passes a repository root, so a retirement removes the same directory as before.

A relative home gives the same directory as before. The drop resolved a relative home against the working directory, and the absolute form names that same directory.

## Acceptance

- [ ] A drop with an empty root keeps the assignment scope of the repository at the working directory.
- [ ] The drop of a retired assignment removes its retained response and its staged excerpts.
- [ ] The spill writer and the drop name the repository directory through one store function.
