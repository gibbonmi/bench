# Repair the empty-remedy refusal panic

Blocked by: none
Ownership fence: `internal/specbuild/render.go`, `internal/specbuild/render_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: EP1/refusal-renders-without-a-remedy

## What to build

Close the accepted Spec and Coverage findings (P2, C1) from the Terra/xhigh
review of candidate `34207c44c606797d996f1812a27b9484b22edae6`:
`RenderRefusal` in `internal/specbuild/render.go` panics whenever the typed
refusal it is given carries zero constructible actions —
`if len(actions) == 0 { panic("specbuild: typed public refusal has no
remedy") }`. This is reachable from ordinary caller input, not only from an
internal programming error: `RefusalForClass` and its callers (`internal/specbuild/disclosure.go`)
build most remedies with `fixed(slug)` (also `fixed(ticketArg)`, `fixed(request)`
at the refresh sites), and `ActionFact.action()` refuses construction — by
design, per `specs/axi-spec-build-complete/spec.md` around line 25 — when a
fixed value carries a control byte no single-line command can quote, such as
an embedded newline. `shellAction`/`buildAction` (`internal/specbuild/disclosure.go`
around lines 24-34) collapse that refusal into a zero-value `ActionFact{}`,
which `actions()` then filters out entirely because `fact.action()` also
refuses on the empty `Program` field. A slug (or ticket argument, or request
string) containing a newline — an ordinary hostile argv value, not a crafted
internal state — therefore reaches `RenderRefusal` with zero actions and
crashes the process with a Go panic instead of returning a clean AXI refusal
at exit 1.

Fix: `RenderRefusal` must render the honest "no useful action" terminal form
instead of panicking when `actions` is empty — the same convention the spec
already defines for a successful response with nothing to advertise:
`axi.RenderHelp` renders `help[0]{command}:` for a nil or empty action slice
(see `internal/axi/action.go`'s `RenderHelp`, and its existing callers in this
file, e.g. `RenderRunsHome`, which already pass a possibly-empty `actions`
slice through untouched). Keep the *first* panic in `RenderRefusal`
(`if !ok { panic("specbuild: untyped public refusal reached renderer") }`)
exactly as it is — an untyped error reaching this renderer is a genuine
internal invariant violation, never triggered by caller input, and must stay
a hard failure. Only the zero-actions branch changes: when there is no
constructible action, render the error line with a generic, safe hint (not
built from `actions[0].Command()`, since there is no such element) and append
`axi.RenderHelp(nil)` for the honest empty help block, at exit 1, exactly like
every other refusal.

## Acceptance

- [ ] [EP1] (covers local) (P2, C1) a public operation whose typed refusal
  cannot construct any remedy action — reproduced with a slug or, for
  `assign`/`refresh`, a ticket argument or request string containing an
  embedded newline — renders a clean AXI refusal at exit 1 with an empty
  `help[0]{command}:` block instead of panicking; every refusal that does
  carry a constructible remedy renders exactly as it does today, byte-for-byte.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EP1/refusal-renders-without-a-remedy | restore the unconditional `panic("specbuild: typed public refusal has no remedy")` for the zero-actions branch | focused `RenderRefusal` test | construct a typed refusal whose only candidate remedy embeds a newline-bearing slug (or ticket/request), render it through `RenderRefusal`, and require a clean exit-1 AXI response rather than a panic |
