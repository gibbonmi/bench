# Introduce the typed action carrier

Blocked by: none
Ownership fence: `internal/axi/action.go`, `internal/axi/action_test.go`
Integration surfaces: action API→`internal/axi/action.go`; registry declarations→declare-axi-query-root-metadata.md; compatibility oracle package root→`internal/axi/compatibility` exercised unchanged (the sibling spec owns it; this ticket adds only root-package files beside it)
Contracts: ordered token vector, fixed argument values, open placeholder names, and executable disposition cross caller→`internal/axi/action.go`, domain is a command template or non-invokable prose, order is literal argv order, and absence is zero actions, asserted by AC1 against a real owner-declared action value rather than a stub validator
Closure: AC1/tokens, AC1/fixed, AC1/open, AC1/prose, AC1/absence, AC1/no-help

## What to build

`internal/axi` gains a typed action carrier. An action is either an executable
command template — an ordered token vector whose arguments are each either a
fixed value or a named open placeholder — or non-invokable prose that the carrier
refuses to mark executable. The carrier preserves literal argv order, keeps fixed
and open arguments distinguishable, treats zero actions as zero, and exports no
help-row renderer: this spec defers public `help[]` entirely, so the package must
not grow an exported identifier that renders help rows.

Tree condition at refresh time: this spec follows `axi-compatibility-oracle`, so
`internal/axi/compatibility` already exists as a sibling package directory. This
ticket writes only the two root-package files on its fence line and must not
touch anything under `internal/axi/compatibility`.

## Acceptance

- [ ] [AC1] (covers CR2) typed actions preserve literal tokens, fixed arguments, open placeholders, prose refusal, and zero-action absence while the package renders no help rows.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AC1/tokens | join the token vector into one space-separated string inside `Action.Tokens` | `TestActionPreservesLiteralTokenVector` in `internal/axi` (this ticket authors it) | run `go test ./internal/axi -run TestActionPreservesLiteralTokenVector -timeout 60s`; expect the length assertion `len(Tokens()) = 1, want 4` for the action built from `bench spec build promote`; bound is the `-timeout 60s` binary deadline over in-memory values |
| AC1/fixed | return every argument as an open placeholder, discarding the fixed value | `TestActionPreservesFixedArgumentValue` in `internal/axi` | run `go test ./internal/axi -run TestActionPreservesFixedArgumentValue -timeout 60s`; expect the assertion `Arguments()[0].Fixed = "", want "--brief"`; bound is the `-timeout 60s` binary deadline |
| AC1/open | replace an open placeholder with a fixed argument holding the empty string | `TestActionKeepsOpenPlaceholderOpen` in `internal/axi` | run `go test ./internal/axi -run TestActionKeepsOpenPlaceholderOpen -timeout 60s`; expect the assertion `Arguments()[1].Open = false, want true` and `Arguments()[1].Placeholder = "", want "<slug>"`; bound is the `-timeout 60s` binary deadline |
| AC1/prose | let `Action.Validate` accept a prose action carrying no token vector as executable | `TestActionRejectsProseAsExecutable` in `internal/axi` | run `go test ./internal/axi -run TestActionRejectsProseAsExecutable -timeout 60s`; expect the `errors.Is(err, axi.ErrProseNotExecutable)` assertion to report a nil error for the prose action `"ask the reviewer to merge"`; bound is the `-timeout 60s` binary deadline |
| AC1/absence | return a one-row default action list when the owner supplies none | `TestActionListKeepsZeroActionsEmpty` in `internal/axi` | run `go test ./internal/axi -run TestActionListKeepsZeroActionsEmpty -timeout 60s`; expect the assertion `len(actions) = 1, want 0` for the outcome-free action list; bound is the `-timeout 60s` binary deadline |
| AC1/no-help | add an exported `RenderHelp` method on the action carrier that formats help rows | `TestActionCarrierExportsNoHelpRowRenderer` in `internal/axi` | run `go test ./internal/axi -run TestActionCarrierExportsNoHelpRowRenderer -timeout 60s`; the test parses `internal/axi/action.go` with `go/parser` and expects the failure `exported identifier RenderHelp renders help rows; help[] is deferred by this spec`; bound is the `-timeout 60s` binary deadline over a single-file parse with no process spawn |
