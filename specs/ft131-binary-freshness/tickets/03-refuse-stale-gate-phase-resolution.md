# Refuse stale gate phase resolution

Blocked by: Publish and verify a content-sealed binary

## What to build

The executable gate entry verifies the exact selected `dist/bench` before
requesting `gate-phases`, normalizes legacy and indeterminate failures to the
stable rebuild action, and preserves every shipped gate launch contract.

## Acceptance

- [ ] The first gate run after an in-scope phase-table or source change refuses
  before old phase output or scheduling appears.
- [ ] After the prescribed successful rebuild, the unchanged next invocation
  resolves and runs the replacement phase table exactly once.
- [ ] Missing executable or seal, unreadable or malformed seal, partial output,
  and a legacy unknown `freshness-check` subcommand all fail closed with one
  stable rebuild action.
- [ ] Root and nested-cwd invocations agree and repeated fresh checks stay green.
- [ ] Direct kit, linked-repository by-path, hook-triggered, and
  adapter-triggered gate launches all reach the same pre-resolution refusal.
