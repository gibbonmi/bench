# Isolate the descendant row of the focused differential

Blocked by: none
Writes: internal/systemtest/process_lifetime_callers_test.go
Covers: none

## What to build

`TestProcessLifetimeFocusedDifferential` goes red under load in a subtest that is not the descendant row. The `surviving_normal_descendant` row leaves a live child on the shared suite HOME. The release of its CLI then seals the cache generation, and the CLI of a sibling row fails with `the resource generation is sealed` before its child starts. The test calls `groupOf` before it examines the result, so it reports a 20 s mark timeout in place of that refusal.

Give the descendant row its own HOME, so that its sealed generation cannot reach a sibling row. Before each `groupOf` call, examine the CLI result and stat the mark, as `TestProcessLifetimeBuilderDifferential` does. Then a missing mark fails at once and prints the CLI output.

Change no deadline, no production file, and no other test. The seal rule of the cache release stays as it is.

## Acceptance

- [ ] The descendant row of `TestProcessLifetimeFocusedDifferential` runs on a HOME that no sibling row uses.
- [ ] When the mark is absent after the CLI returns, the test fails at once with the CLI exit code and output, and it does not wait for the mark.
- [ ] Under concurrent `bench test --check system` load, 20 focused runs of `TestProcessLifetimeFocusedDifferential` are green.
- [ ] `bench test --check system` is green.
