# Worktree cleanup requires verifiable ownership and preservation

Bench owns worktree lifecycle through an immutable Git-private ownership marker, a separate assignment record, a dedicated branch, and a protective Git lock. Every harness routes creation and release through the same deterministic core.

Three consumers decide whether a worktree is eligible for cleanup: the operator's exact-path command, the unattended automatic sweep, and the landed-set selector. Each resolves that question through one ordered eligibility verdict. The conjunctive ownership and preservation requirements below are therefore decided once, and merely projected onto each command's own operator-facing plan and message.

Automatic cleanup acts only on a matching, landed assignment whose recorded owner is proven dead or absent. A live or unprovable owner retains the tree. Git-visible changes are anchored under durable recovery references before unlock. Ignored residue is discardable only when a valid repository declaration contains the complete bounded inventory; foreign, malformed, undeclared, truncated, or unrepresentable state remains retained.

Explicit cleanup uses an exact-path plan/apply fingerprint, with ignored deletion separately bound. Branch and path conventions and non-blocking harness callbacks cannot prove ownership or preservation.
