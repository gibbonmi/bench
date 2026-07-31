# Migrate line binding through the harness matrix

Blocked by: Capture authoritative Claude fork evidence

## What to build

Codex and Claude resolve peer tier columns through `internal/lines`; the CLI,
hook, adapters, denial, describe output, profile, and line guidance all consume
the hard-cut matrix without a canonical model family.

## Acceptance

- [ ] Codex and Claude tier tokens resolve to their own complete matrix columns.
- [ ] OpenCode remains unbound and fails closed at launch.
- [ ] Every runtime caller passes its harness explicitly and preserves its existing fail-open rim.
- [ ] Claude denial and describe advice lead with Claude-native tokens.
- [ ] A captured context fork denies only when its declared and inherited tiers differ.
