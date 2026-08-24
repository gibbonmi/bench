# Make the session handoff a local, git-ignored file

Blocked by: none
Writes: .gitignore, AGENTS.md, internal/git/localnote.go (new), internal/roadmap/roadmap.go, internal/handoff/handoff.go, internal/handoff/text.go, internal/status/handoff.go, internal/conformance/handoff_single_source_test.go

## What to build

Git ignores `capture/session-handoff.md` in this repo. The file stays on disk
and stays the cold-start pin. A shared helper resolves the primary checkout
for an ignored note file; `bench idea` and `bench handoff` both write through
it. `bench status` computes the handoff's age from the file's own write time
when the file is ignored. The conformance Shape check grades the artifact
only when git tracks it, so a fresh worktree without the file stays green.

## Acceptance

- [ ] `git check-ignore capture/session-handoff.md` exits zero, and the file is not in the index.
- [ ] With the handoff ignored, `bench handoff` in a linked worktree writes the primary checkout's copy.
- [ ] With the handoff ignored and commits after its last write, `bench status` still reports the handoff row with the commit distance.
- [ ] With the handoff untracked, the Shape single-source check emits no artifact diagnostic; a tracked drifted copy still bites.
