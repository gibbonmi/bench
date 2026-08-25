# 16. The installed broker owns promotion

Status: accepted (2026-08-25)

## Decision

One installed promotion broker owns the complete landing command. The
installation manifest authenticates the broker by path, version, and
executable digest. The wrapper selects the broker from the installed
distribution only. It refuses inherited routing overrides for this public
command before any repository read, and it ignores repository binaries and
their seals.

The broker resolves the destination and the registered assignment as
separate subjects. It binds the request, the review base, the source tip,
and the source fingerprint before composition. The review base must be the
assignment's recorded start or a descendant of it.

The broker composes the exact prospective tree in private storage and
builds the gate executable from that tree. That executable is the graded
subject, never the publication authority, and the baseline schedule selects
the phases. Green evidence keys to the prospective tree and the baseline
runner identity. Only the broker validates evidence and performs the
destination compare-and-swap; resume stays under the same owner.

## Consequence

Candidate code cannot authorize its own publication, and a landing never
rebuilds and re-runs itself. A broker-changing diff publishes as source
only; the new broker takes authority at the next package install or repair,
and the landing names that step.
