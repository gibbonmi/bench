# Single-source prefix projection

Blocked by: 04-close-classifier-edge-gaps.md
Writes: internal/shellcommand, internal/gitguard, internal/benchguard, internal/systemtest

## What to build

Give the shared shell-command module one command-word projection and one routine-prefix
description. Preserve the distinct Git-guard and Bench-guard policy results while
closing value-taking option gaps and command query false positives. Complete the
wrapper outer-syntax process matrix.

## Acceptance

- [ ] Command-word projection and the routine-prefix set have one production source.
- [ ] `env -u X`, `timeout -s KILL`, and `xargs -n 1` cannot hide a Bench call.
- [ ] `command -V bench` and equivalent query-only forms remain allowed.
- [ ] Wrapper-contained Bench calls followed by outer `|`, `&&`, or redirection are refused.
- [ ] The destructive-Git classification matrix is unchanged.
