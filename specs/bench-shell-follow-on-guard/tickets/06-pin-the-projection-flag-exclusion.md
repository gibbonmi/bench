# Pin the projection flag exclusion

Blocked by: 02-refuse-bench-shell-follow-ons.md
Writes: cmd/bench

## What to build

Pin the public gate grammar against a `--brief` projection flag.
Keep the existing usage error and add no public command or option.

## Acceptance

- [ ] `bench gate --brief` exits 2 with the gate usage line.
- [ ] A special-case acceptance of `--brief` makes the focused proof red.
