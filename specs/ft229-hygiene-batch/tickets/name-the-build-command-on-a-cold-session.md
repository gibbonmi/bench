# Name the build command when the core is missing

Blocked by: narrow-the-degraded-guard-rim.md
Writes: .bench/hooks/session-start.sh, internal/sessioninspect, internal/systemtest

## What to build

SessionStart with no core execs into a missing binary and the session opens with
no dashboard and no hint. The hook instead prints one line naming `bash scripts/go-build.sh` with its
arguments — the profile forbids plain `go build`, which produces a binary
carrying no version — and exits zero. It still never blocks a
session, and it still prints nothing outside a repository, so the hint does not
become ambient noise.

Evidence comes from running the hook as a real subprocess with the core absent.

## Acceptance

- [ ] a session start with an unreachable core prints the `scripts/go-build.sh` invocation and exits zero.
- [ ] a session start with a reachable core is unchanged.
- [ ] a session start outside a repository prints nothing.
