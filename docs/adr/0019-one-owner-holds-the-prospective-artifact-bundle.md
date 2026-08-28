# 19. One owner holds the prospective artifact bundle

Status: accepted (2026-08-28)

## Decision

One bundle owner holds the private prospective checkout and every
owner-authored run binary. Both resources sit under one dedicated prefix in
the operating-system temporary root, so the bundle root is the one source for
the teardown scope. The owner publishes a strict record before it creates the
checkout. That record carries the schema, the owner process identifier, and
the canonical Git common directory of the repository. The owner writes the
record as a regular file in the private 0600 mode, and an unrecorded tree
therefore has no owner.

Every prospective producer sweeps the recognized dead bundles of its own
repository before it creates a bundle. The full gate, the evidence inspection,
and the fast lane are those producers. The sweep accepts only a canonical
bundle root whose valid record names the current repository. Only the
operating system's absent-process result proves that an owner is dead. An
answering probe, a permission-refused probe, and any other failed probe
retain the bundle. The sweep removes the Git registration before it removes
the bundle root, and a removal failure refuses the new operation.

The owner resolves the temporary root through its symbolic-link components,
because Git registers the resolved spelling of a worktree. An inherited run
binary stays outside the bundle and keeps its existing lifetime. The sweep
executes no recovered bytes, so residue never becomes cleanup authority.

## Consequence

A killed prospective run leaves no durable residue, because the next run on
the same repository removes the dead owner's bundle first. Four states stay
outside this recovery and leak safely instead. These states are an unrecorded
legacy tree, a changed temporary root, a reused owner process identifier, and
a descendant that survives the killed owner.
