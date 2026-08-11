# Repair the refresh invalid-evidence remedy

Blocked by: none
Ownership fence: `internal/specbuild/refresh.go`, `internal/specbuild/refresh_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: RF1/refresh-invalid-evidence-carries-known-facts

## What to build

Close the accepted Spec and Coverage findings (P1, C1) from the Terra/xhigh
review of candidate `5ae4029d6207fd1bc1f75b85612d2c4492baae68`:
`Service.Refresh` (`internal/specbuild/refresh.go`) has two sites where a
malformed or otherwise-invalid debug receipt refuses with
`operationRefusal(RefusalInvalidEvidence, "assign", slug, err)` — after
`readDebugReceipt` fails, and after `validateDebugReceipt` fails.
`operationRefusal` routes through `RefusalForClass`'s generic
`RefusalInvalidEvidence` case, which is `operationAction("assign", slug)`: a
bare `bench spec build assign <slug> --ticket <ticket> --request <request>`
template with no `--refresh`. Retrying that literal command against an
already-owned assignment just returns the existing assignment
(`internal/specbuild/assign.go` around lines 67-76) — it does not resubmit a
refresh, so the advertised remedy cannot satisfy the observed state, the same
wrong-remedy class already repaired for the stale/spent-refresh cases
(`specs/axi-spec-build-complete/spec.md` around lines 17-18, 58). The caller's
own `ticketArg` and `request` are already in scope at both call sites — the
same values the neighboring `refreshRefusal(RefusalStaleRefresh/RefusalSpentRefresh,
slug, ticketArg, request, ...)` calls in this function already carry forward.
No existing test drives either invalid-evidence site for its rendered remedy:
the existing invalid-receipt tests
(`internal/specbuild/refresh_test.go` around lines 179-235) check only
rejection and unchanged state, and the remedy assertions elsewhere in that
file cover only the stale/spent classes.

Fix: at both call sites, replace `operationRefusal(RefusalInvalidEvidence,
"assign", slug, err)` with `refreshRefusal(RefusalInvalidEvidence, slug,
ticketArg, request, err)` — `refreshRefusal` (`internal/specbuild/disclosure.go`,
already used by the stale/spent-refresh repair) takes any `RefusalClass` and
builds its remedy from `refreshAction(slug, ticketArg, request)`, so this
requires no change to `RefusalForClass`, `operationAction`, or any other
refusal construction. This also leaves the disclosure matrix's
`assign/invalid-evidence-receipt` fixture cell untouched: that cell reaches
`RefusalForClass` directly with no known ticket/request
(`internal/specbuild/disclosure_observation.go`'s `installRefusal`, the `else`
branch), which this change does not touch.

New coverage in `internal/specbuild/refresh_test.go`: drive the real
`Service.Refresh` with an evidence file that fails `readDebugReceipt` (or, for
the second site, one that fails `validateDebugReceipt` — whichever composes
more naturally with this file's existing fixtures), read the resulting typed
refusal through `RefusalFacts`, and require its one action's `Command()`
carries the caller's own slug, ticket argument, and request as fixed values
plus `--refresh <receipt>` with the receipt open — the same shape the
stale/spent-refresh remedy tests already assert.

## Acceptance

- [ ] [RF1] (covers local) (P1, C1) `Service.Refresh` given evidence that
  fails `readDebugReceipt` or `validateDebugReceipt` returns a
  `RefusalInvalidEvidence` refusal whose remedy is `bench spec build assign
  <slug> --ticket <ticket> --request <request> --refresh <receipt>` with the
  caller's own slug, ticket argument, and request as fixed values — driven
  through the real public `Refresh` call, not a hand-constructed refusal;
  every other refusal class's remedy, and the disclosure matrix's
  `assign/invalid-evidence-receipt` fixture, are unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RF1/refresh-invalid-evidence-carries-known-facts | restore either site's `operationRefusal(RefusalInvalidEvidence, "assign", slug, err)` call | focused real-`Refresh` refusal test | drive `Service.Refresh` with evidence that fails validation at each site, read the remedy via `RefusalFacts`, and require the caller's slug, ticket argument, and request as fixed tokens plus `--refresh <receipt>` |
