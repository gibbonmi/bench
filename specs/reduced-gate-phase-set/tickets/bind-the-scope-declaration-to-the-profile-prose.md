# Bind the scope declaration to the profile prose

Blocked by: Route the ambient staleness signal through the scope declaration, Compute the stripped subject identity, Run excludable phases against a stripped worktree, Confine contract-test capture reads to the subject root, Take the reduced path from bench commit, Refuse a reduced verdict for release evidence

Ownership fence: `internal/conformance/scope_binding_test.go`, `projects/benchkit.md`
Assumptions: every consumer already routes through the declaration, so the prose is the last remaining second derivation of the allowlist and the phase set

## What to build

A conformance check binding the allowlist, the excludable phase set, and the prose
that documents them to their single source, following `checkLineBinding`'s shape —
a prose table cross-checked against its machine-readable source, red on divergence.
Without it a future edit widens the fast path in one place only, and the document
and the oracle disagree with nothing to notice.

The profile's pinned allowlist prose is updated in this same change. It currently
pins the four scattered entries and names expansion "a new decision"; this build is
that decision, so the profile must record the decision that was made rather than
contradict the code.

Bind in both directions. A subset or substring comparison survives a prose-only
addition, so the check has to red for a mutation on either side alone — that is what
rules the weak implementation out, and it is why the bite test drives each direction
separately rather than trusting one.

## Acceptance

- [ ] [R05] The profile's pinned allowlist prose matches the declaration, and widening the declaration without the profile edit produces the binding diagnostic.
- [ ] [R25] The declaration, the excludable phase set, and their prose cannot drift: mutating either side alone produces the binding diagnostic.
- [ ] [R26] The bite test mutates the prose alone and the declaration alone and requires a diagnostic for each.
