# 19. One owner holds the prospective artifact bundle

Status: accepted (2026-08-28)

## Decision

One bundle owner holds the private prospective checkout and every
owner-authored run binary. Both resources sit under one dedicated prefix in
the operating-system temporary root, so the bundle root is the one source for
the teardown scope. The owner opens the bundle's protection generation, and
then it publishes a strict record before it creates the checkout. That record
carries the schema, the owner process identifier, the canonical Git common
directory of the repository, and a reference to the protection generation.
The owner writes the record as a regular file in the private 0600 mode, and
an unrecorded tree therefore has no owner. Each launch that uses the checkout
registers as a user of that generation, and a nested launch inherits it.

Every prospective producer sweeps the recognized dead bundles of its own
repository before it creates a bundle. The full gate, the evidence inspection,
and the fast lane are those producers. The sweep accepts only a canonical
bundle root whose valid record names the current repository. Only the
operating system's absent-process result proves that an owner is dead. An
answering probe, a permission-refused probe, and any other failed probe
retain the bundle.

Owner death alone does not authorize removal. Removal also requires complete
user proof: every registered user of the referenced generation is resolved.
A live descendant group that is a registered user retains the bundle. A
legacy record that references no generation proves no user set, so it
retains the bundle too. The owner's own close observes the same users, so a
surviving user retains the bundle after a normal exit as well. The close and
the sweep observe the users before they remove the Git registration, then
they remove the bundle root. A removal failure refuses the new operation.

The gate and the fast lane release the bundle before they record an outcome.
A retained bundle permits no green verdict, no retained green evidence, and
no lane pass. A run binary that the owner authored inside the bundle releases
first, and its retention retains the bundle.

The owner resolves the temporary root through its symbolic-link components,
because Git registers the resolved spelling of a worktree. An inherited run
binary stays outside the bundle and keeps its existing lifetime. The sweep
executes no recovered bytes, so residue never becomes cleanup authority.

## Consequence

A killed prospective run leaves no durable residue after its registered users
end, because the next run on the same repository then removes the dead
owner's bundle first. Until those users end, the bundle stays in place. Five
states stay outside this recovery and leak safely instead. These states are an
unrecorded tree, a legacy record, a bundle whose owner had no protection, a
changed temporary root, and a reused owner process identifier.
