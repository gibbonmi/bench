# Resolve a parameter grammar through its callers

Blocked by: 6-retire-step-kind.md
Writes: internal/preflight/ownership/producers.go, internal/preflight/ownership/effects_test.go, internal/preflight/ownership/arguments.go (new), internal/preflight/ownership/parameter_grammar_test.go (new)
Covers: AP49

## What to build

The charge-premises grammar-effects check roots the grammar argument of each `usage.Parse` call. A function such as `parseHarness` in `cmd/bench/main.go` receives its grammar as a parameter and passes it on. The producer resolver cannot root a parameter, so it reports the call as unresolved. Today any ticket after a `cmd/bench` change in the same chunk is refused at its charge.

Make a parameter grammar resolve through its callers. When each call site of the function in its package passes a grammar that the resolver already roots, root the parameter to those producers. When one call site passes a grammar that does not resolve, the call stays unresolved, as it does now. Compose the existing producer and call-resolution code in `producers.go`; do not write a second walker.

## Acceptance

- [ ] The ticket 7 charge, `bench preflight build anchor-admission-prune --charge --ticket 7-admit-by-class.md` with the AP-C5 chunk base as `--base`, passes `charge-premises`.
- [ ] A regression test shows that a parameter grammar with declared callers resolves, and that one undeclared caller keeps it unresolved.
- [ ] `bench test --package ./internal/preflight/...` and `bench test --check root-conformance` pass.
