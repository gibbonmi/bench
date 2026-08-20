# Narrow the degraded guard rim to the command field

Blocked by: none
Writes: .bench/hooks/block-dangerous-git.sh, internal/systemtest

## What to build

When no core is reachable the shim's rim refuses any envelope whose raw text
contains `git` anywhere, so reading a file under `.github/` is blocked during
exactly the cold session that needs to recover. The rim instead extracts the
envelope's command field and tests whether that command invokes `git`, at the
same one-level wrapper depth the core documents. An envelope with no readable
command field still refuses, so the fail-closed posture survives. The rim stays
an honest-mistake layer and the guard's stated threat model is unchanged.

Evidence comes from running the shim as a real subprocess with the wrapper
unreachable and `bench` off PATH.

## Acceptance

- [ ] with no reachable core, a destructive git command is refused with exit 2 (H20).
- [ ] with no reachable core, a command whose `git` text appears only in a path or an unrelated envelope field is allowed (H21).
- [ ] an envelope carrying no readable command field is refused (H22).
- [ ] the guard's documented honest-mistake threat model is unchanged (H33).
