# Repair the duplicate partial-reclaim help[] finding

Blocked by: none
Ownership fence: `internal/specbuild/render.go`, `internal/specbuild/reclaim_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: DH1/single-envelope-partial-reclaim

## What to build

Close the accepted Spec and Coverage findings (P1, C1) from the Terra/xhigh
review of candidate `123edbab919d2462c73f206dd6623aba5feccc9f` (receipt
`925a9ec921f57225b529def2885158427ba9858ecacf152b64402751357e250b`):
`RenderReclamationFailure` in `internal/specbuild/render.go` renders two
`help[]` blocks when a reclamation apply fails after already deleting at
least one ref. Its `len(receipt.Refs) != 0` branch calls `RenderReclamation(receipt)`
for the spent-deletion tables — which appends its own trailing
`axi.RenderHelp(plan.Actions())` — then appends `RenderRefusal(err)`, which
appends its own trailing `axi.RenderHelp(actions)`. The response ends up with
two help envelopes, violating the spec's single-response-envelope rule
(`specs/axi-spec-build-complete/spec.md:24`: "one `help[N]{command}` block
appended after the primary response"). The one existing partial-reclaim test
with spent refs (`TestReclaimApplyReportsTheDeletionsItSpentBeforeADriftedRef`
in `internal/specbuild/reclaim_test.go`) stops at the raw `ApplyReclaim`
receipt and error and never renders them, so this never got caught; the
operation/state disclosure matrix's reclaim/spent-fingerprint fixture faults
before any ref is deleted (`internal/specbuild/reclaim.go:92`), so
`receipt.Refs` is empty there and the buggy branch is unreached through that
path too.

Fix: extract the reclaim summary and per-ref table rendering currently
inlined at the top of `RenderReclamation` into an unexported helper that
renders those two tables without a trailing help block. `RenderReclamation`
calls the helper and then appends its own `axi.RenderHelp(plan.Actions())`,
unchanged in its own output bytes. `RenderReclamationFailure`'s
`len(receipt.Refs) != 0` branch calls the same helper (not
`RenderReclamation`) for the spent-deletion portion, then appends the
refusal — so the final response carries exactly the refusal's one help
block, and the spent-deletion tables stay fully present.

## Acceptance

- [ ] [DH1] (covers SB1, SB5) (P1, C1) a reclamation apply that fails after
  deleting at least one ref renders the spent-deletion tables followed by
  exactly one `help[]` block — the refusal's — never two; a reclamation
  apply or plan that succeeds, and a reclamation failure that deleted
  nothing, are byte-unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DH1/single-envelope-partial-reclaim | restore the direct `RenderReclamation(receipt)` call in `RenderReclamationFailure`'s nonempty-`Refs` branch (or otherwise let the spent-deletion portion append its own help block) | focused render test | drive a reclamation apply that deletes at least one ref before hitting a typed refusal, render the returned receipt and error through `RenderReclamationFailure`, and require the output to contain exactly one `help[` block |
