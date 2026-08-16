# Clean up SIGINT before replacement

Blocked by: clean-up-failed-replacements.md
Writes: internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go, internal/skillsindex/command.go, internal/skillsindex/command_test.go

## What to build

Thread one cancellation context through the real write/command path and the replacement
lifetime published by the failed-replacement ticket. At the production-reached
pre-replacement seam, a child operation publishes an exact inherited-pipe marker after
the real temp is created and written, then blocks on that context. The parent waits for
the marker, sends SIGINT, and observes cleanup before nonzero exit without directory
polling or sleep.

This ticket owns the entire deterministic signal handshake. The failed-replacement
ticket is its only dependency because that ticket publishes the cleanup lifetime and
replacement-operation seam that cancellation extends; no parser, payload, orphan, or
repository-diagnostic behavior is needed to make HI10 independently green.

## Acceptance

- [ ] `(covers HI10)` The child publishes the exact pre-replacement marker, blocks on
  context, receives SIGINT only after the parent observes the marker, exits nonzero,
  preserves original bytes, and leaves no residue.
