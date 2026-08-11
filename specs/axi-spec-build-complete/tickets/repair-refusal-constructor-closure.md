# Close typed refusal construction over the policy axis

Blocked by: none
Ownership fence: `internal/specbuild/disclosure.go`, `internal/specbuild/disclosure_test.go`
Integration surfaces: raw typed refusal construction→`internal/specbuild/disclosure.go`; production refusal policy axis→existing `refusalPolicies` in `internal/specbuild/disclosure.go`; constructor-closure AST and policy-input mutations→`internal/specbuild/disclosure_test.go`; conformance matrix→existing `internal/conformance/axi_disclosure_test.go` exercised unchanged through `DisclosureCells`
Contracts: every public and helper refusal crosses its wrapper→`RefusalForClass` in `internal/specbuild/disclosure.go`, asserted by UC1 against the production AST; declared class membership crosses `refusalPolicies`→typed refusal construction in the same owner, asserted by UC2 and UC3
Closure: UC1/single-raw-constructor-caller, UC2/undeclared-class-refused, UC2/no-undeclared-typed-facts, UC3/policy-membership-authorizes-construction

## What to build

Close the accepted Terra finding S1. `refusalPolicies` is the authorization
source for typed refusal classes: `RefusalForClass` constructs a typed refusal
only for a declared class, while an undeclared class fails closed without typed
facts or an advertised action. Every production wrapper reaches the raw
`refusal` constructor only through `RefusalForClass`, so adding a wrapper name
cannot bypass closure and the test does not maintain a second wrapper registry.

## Acceptance

- [ ] [UC1] (covers SB2) (S1) `RefusalForClass` is the only production caller of the raw typed `refusal` constructor.
- [ ] [UC2] (covers SB2) (S1) an undeclared refusal class fails closed without typed refusal facts or a remedy action.
- [ ] [UC3] (covers SB7) (S1) removing a class from `refusalPolicies` makes construction of that class fail, so policy membership authorizes the constructor and the disclosure axis together.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| UC1/single-raw-constructor-caller | call raw `refusal` from a new or existing wrapper | the production-AST constructor-owner test | parse non-test `internal/specbuild` files, enumerate every direct `refusal` call, and require the sole caller to be `RefusalForClass` |
| UC2/undeclared-class-refused | let `RefusalForClass` fall through its action switch for an unknown class | the pure refusal-construction test | construct a synthetic undeclared class, require an error, then require `RefusalFacts` to report no typed class or actions |
| UC2/no-undeclared-typed-facts | wrap the undeclared-class error in a typed `Refusal` | the pure refusal-construction test | construct the synthetic class and require `errors.As` not to find a typed refusal |
| UC3/policy-membership-authorizes-construction | bypass the `refusalPolicies` membership check for one known class | the policy-input mutation test | remove one declared class from the policy input, construct that class, and require construction to fail until the policy member is restored |
