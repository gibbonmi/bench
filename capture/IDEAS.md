
- 2026-08-07  relax the link ownership gate for converged unowned adapter destinations: the kit repo (and any repo satisfying adapters with directory symlinks, never self-linked) has no manifest, so owned=="" skips convergedFingerprint and the symlink-parent rule hard-aborts bench link; deciding whether bench may claim ownership of files it never wrote is a policy change needing its own ticket
