# Add a per-repository opt-out for the protected-branch push clause

Blocked by: none
Writes: internal/adopt/prepush.sh, internal/adopt/link_hook_test.go, .bench/BENCH-reference.md, README.md

## What to build

The managed pre-push hook reads `git config bench.allowProtectedPush`. When
the value is `true`, the hook lets a push to the protected branch through.
The `.bench` drift clause stays armed. The docs name the knob.

## Acceptance

- [ ] Without the config, the hook blocks a push to the protected branch.
- [ ] With `bench.allowProtectedPush=true`, the hook lets a push to the protected branch through.
- [ ] With `bench.allowProtectedPush=false`, the hook blocks a push to the protected branch.
- [ ] The hook reference doc and the README name the knob.
