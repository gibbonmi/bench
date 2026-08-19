# Learnings — usage journal

- 2026-08-19 — /bench-write-spec (FT228) took 2 iterations to accept. Stage
  that missed: spec authoring. What review caught: (1) the spec flipped
  `$bench-debug` to implicit invocation while the adapter's own description
  still said "Use only when the reviewer invokes" — the settle's observable
  behavior lived in match text the author treated as boilerplate; (2) the new
  fixture family needed a `registry_test.go` registration that sat outside
  every ownership fence. Why missed: the author verified the check and fixture
  mechanics but not the complete landing path for a *new* fixture family, and
  read the policy key without reading every surface the harness matches
  against. Proposed rule change: when a spec flips an invocation or trigger
  policy, enumerate every surface the harness matches (description text, not
  only the policy key); when a spec adds fixtures, verify the family's
  registration exists or fence the registry that grants it.
