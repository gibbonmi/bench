# Narrow the pre-push README claim

Blocked by: none
Ownership fence: `README.md`, `internal/conformance/docs_workflow_checks_test.go`
Integration surfaces: pre-push protection claim→`README.md` + P1; `docs-currency-workflow` assertion→`internal/conformance/docs_workflow_checks_test.go` + P1
Contracts: the pre-push protection claim crosses `README.md`→`internal/conformance/docs_workflow_checks_test.go` as a README sentence whose domain is the hook's resolved branch, whose clause order is unconstrained, and whose resolved-branch absence or unqualified default-branch presence is invalid, asserted by P1 against the real README
Closure: P1/resolved-branch-claim, P1/unqualified-default-absence

## What to build

State that the installed pre-push hook protects the branch it resolves, and make the existing `docs-currency-workflow` check reject restoration of the unqualified default-branch guarantee. The README edit and its assertion stay together: an assertion-only cut makes the project gate's `docs-currency-workflow` check red with the missing resolved-branch and unqualified default-branch diagnostics against the current sentence, while a README-only cut leaves restoration of that sentence ungraded, so neither thinner cut is an independently-green behavior-plus-test ticket.

## Acceptance

- [ ] [P1] `README.md` states that pre-push protects the branch the hook resolves and contains no unqualified claim that it protects the default branch.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| P1/resolved-branch-claim | replace the resolved-branch qualification with a vague branch reference while retaining no default-branch promise | the `docs-currency-workflow` check | change the sentence to say only that the hook protects a branch, run `BENCH_CONFORMANCE_CHECK=docs-currency-workflow go test -count=1 ./internal/conformance -run '^TestRootConformance$'`, expect only the missing resolved-branch diagnostic |
| P1/unqualified-default-absence | add a second unqualified default-branch promise elsewhere in `README.md` | the `docs-currency-workflow` check | append the stale promise, run `BENCH_CONFORMANCE_CHECK=docs-currency-workflow go test -count=1 ./internal/conformance -run '^TestRootConformance$'`, expect the unqualified default-branch diagnostic |
