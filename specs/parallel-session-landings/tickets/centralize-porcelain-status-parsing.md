# Centralize porcelain status parsing

Blocked by: none
Writes: internal/git, internal/landing

## What to build

Make the Git package the sole owner of strict porcelain-v1 `-z` framing while
preserving landing's filtering of declared build-output residue.

## Acceptance

- [ ] Landing consumes the Git-owned parser and retains strict malformed-record refusal.
- [ ] A parser-bypass mutation makes the focused landing test red.
- [ ] Fingerprint filtering removes declared runtime residue but retains undeclared ignored paths.
