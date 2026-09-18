# Charge evidence is one immutable artifact that a consumer reads in bounded pages

A build author or a review axis needs the complete pinned charge, and one harness tool result cannot carry it.
A charge preparation therefore publishes one immutable evidence artifact and returns only its identity.
The identity is the hash of a canonical manifest, and that manifest binds the exact bytes of each source through ordered page hashes.
A consumer reads the manifest and each source through stateless page commands, and every response stays inside one fixed encoded budget.
The consumer acts only after it verifies the delivery against the identity and binds the artifact to the current assignment and source pair.

## Consequences

- No command returns a complete charge in one response, and no compatibility route exists for that form.
- Verified delivery proves bytes only. Reviewer approval and the task supplement stay separate prerequisites of an action.
- A fresh consumer runs its own retrieval from the trusted identity. A receipt or a final cursor from another consumer delivers nothing to it.
- The store holds artifacts under a default quota, and Bench evicts nothing on its own. An operator removes artifacts with an explicit cleanup that applies one exact fingerprinted plan.
- A review capture is frozen with its provenance. Bench does not make a collector reproducible across machines.
- The write-spec phase keeps its author fork and takes no charge retrieval.
- A reader must be able to open the store lock, so a sandboxed harness needs access to the store.

## Considered options

- A larger single response: rejected, because the harness ceiling is outside Bench control and a truncated charge looks complete.
- A summary with source paths in place of bytes: rejected, because a consumer in another worktree or after a release reads other bytes.
- Automatic eviction: rejected, because an idle artifact can still be the pinned evidence of an open review.
