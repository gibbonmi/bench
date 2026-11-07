# Refuse a hand-built inherited registration outside the owner

Blocked by: 04-run-inherited-under-local-protocol.md
Writes: internal/conformance/process_lifetime_closure_test.go, internal/conformance/process_lifetime_test.go, internal/conformance/process_lifetime_copy_test.go (new)
Covers: RS18, RS29

## What to build

The process-lifetime closure check adds a hand-built registration rule to `processLifetimeFileDiags` in `process_lifetime_closure_test.go`.
The rule counts three readers of the process vector: `protection.Inherited()`, `protection.InheritedFrom(os.Environ())`, and `protection.ParseVector(os.Getenv(protection.Env))`.
It refuses a production function outside `internal/subprocess` that calls a counted reader and builds a `subprocess.Registration` literal, in its own body or in a same-package function that it calls.
It does not count `InheritedFrom` of another environment, `ParseVector` of another value, or a reader with no `Registration`.
The closure prefilter parses only files that spell an owner token, so the rule adds its own reader tokens to that prefilter.
Reuse the closure's existing parser, import resolution, and same-package call walk. Do not write a second AST scanner.
`process_lifetime_test.go` changes only if the rule needs a census constant there.

The synthetic fixtures go in the new `process_lifetime_copy_test.go`: one per counted reader, one with the `Registration` in a same-package helper, and one uncounted reader that stays green.

Scratch-patch evidence (slicing check, 2026-11-07, at 92cfa43e8): a rough form of the rule, which counted the three readers and the literal in one function body, ran with `bench test --in ft446-scratch-05 --package ./internal/conformance --run 'ProcessLifetime'`. It named exactly `internal/bounds/bounds.go#run` and `internal/runbinary/runbinary.go#runBuildScript`. It named no other production function, and `internal/worktree/subshell.go` stayed green. The scratch worktree was discarded. After ticket 04 lands, the live tree has no hit.

## Acceptance

- [ ] The rule names a synthetic hand-built copy for each counted reader, and one whose `Registration` sits in a same-package helper.
- [ ] The rule stays green on a synthetic uncounted reader.
- [ ] The live `process-lifetime-closure` check stays green after ticket 04, with both runner census entries.
- [ ] A surviving hand-built copy in either runner turns the live check red.
