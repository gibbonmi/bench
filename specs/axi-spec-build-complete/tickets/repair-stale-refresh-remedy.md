# Repair the stale and spent refresh remedy

Blocked by: none
Ownership fence: `internal/specbuild/disclosure.go`, `internal/specbuild/refresh.go`, `internal/specbuild/refresh_test.go`, `internal/specbuild/disclosure_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: SR1/refresh-remedy-carries-known-facts, SR2/refresh-remedy-carries-the-refresh-flag

## What to build

Close the accepted Spec finding (P2) from the Terra/xhigh review of candidate
`399dca908c7b1e1a4162eb7625e497dfb6786750`: `RefusalForClass`
(`internal/specbuild/disclosure.go:256-275`) groups the two refresh classes
with the missing-assignment class on line 263 —
`case RefusalMissingAssignment, RefusalStaleRefresh, RefusalSpentRefresh:
action = operationAction("assign", slug)` — and `operationAction`
(`disclosure.go:51-63`) builds `bench spec build assign <slug> --ticket
<ticket> --request <request>`: both known values reopened as placeholders, and
no `--refresh` at all. A retry of that literal command is a fresh assign, not
the refresh the caller was attempting, so the advertised command cannot satisfy
the observed state — the class story 2 exists to close (`spec.md:18`, SB3 at
`spec.md:58`: "the remedy derives from the typed precondition result and
advertises the command that can satisfy that exact class").

Every stale and spent refresh refusal is raised inside `Service.Refresh`
(`internal/specbuild/refresh.go:131, 140, 153, 156, 162, 184, 244`), well after
line 117 has parsed the ticket, so the caller's own `ticketArg` and `request`
are in scope at all seven sites. Give the two refresh classes their own
constructor beside the existing typed ones in `disclosure.go` — taking the
slug, the ticket argument, and the request — that builds `bench spec build
assign <slug> --ticket <ticket> --request <request> --refresh <receipt>` with
slug, ticket, and request as fixed values and only the receipt open, and call
it from all seven refresh sites in place of `operationRefusal`. `--refresh` is
a declared assign flag (`internal/spec/build.go:50`), so the template is
invokable as rendered.

**Carry `ticketArg`, not the resolved path.** `assigned.Ticket` and
`ticket.Path` are absolute (`refresh.go:135` compares them), and `ParseTicket`
refuses an absolute argument outright (`internal/specbuild/assign.go:397`:
`if arg == "" || filepath.IsAbs(arg)`). A remedy carrying the resolved path
would advertise a command that cannot run — the same wrong-remedy class the
finding is about. The caller-supplied `ticketArg` is the value that reproduces
the invocation.

`RefusalMissingAssignment` stays exactly as it is: no ticket or request is
known at that refusal, so its open-placeholder assign template is already the
honest remedy. Split it out of the shared case and leave stale and spent
refresh on their own case, whose no-known-values form — reached only by the
matrix constructor at `disclosure_observation.go:249`, which has no ticket or
request to supply — still adds `--refresh <receipt>` so the ledger's cell shows
the flag shape production emits.

Both existing test files are in the fence, and each carries one row: the
class-exact behavior can only be driven through the real `Refresh` path in
`internal/specbuild/refresh_test.go` (its `TestRefreshRefuses...` cases at
`:179, 266, 303, 370` own these refusals but assert no remedy today), while
`internal/specbuild/disclosure_test.go:87-88` hard-pins the two classes'
template strings in `TestEveryRefusalClassRetainsTypedIdentityAndExactRemedy`
and goes red on the `--refresh` addition. Neither file can stand in for the
other, and leaving `disclosure_test.go` out of the fence would make the ticket
unimplementable.

## Acceptance

- [ ] [SR1] (covers local) (P2) a stale-refresh and a spent-refresh refusal
  raised by the real `Service.Refresh` each carry a remedy of the form
  `bench spec build assign <slug> --ticket <ticket> --request <request>
  --refresh <receipt>` whose slug, ticket, and request are the caller's own
  values rendered as fixed arguments — the ticket exactly as the caller spelled
  it, never the resolved absolute path — with only the receipt open.
- [ ] [SR2] (covers local) (P2) the stale-refresh and spent-refresh classes
  hold their own `RefusalForClass` case, whose template carries `--refresh
  <receipt>`, while `RefusalMissingAssignment` keeps its current
  `assign --ticket <ticket> --request <request>` remedy byte-unchanged and
  every other class's remedy is unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SR1/refresh-remedy-carries-known-facts | restore `operationRefusal(RefusalStaleRefresh/RefusalSpentRefresh, ...)` at the refresh sites, or carry `ticket.Path` instead of `ticketArg` | focused real-`Refresh` refusal test | drive a stale receipt and a checkpointed assignment through `Service.Refresh`, read the remedy from `RefusalFacts`, and require the caller's slug, ticket argument, and request as fixed tokens and a runnable ticket argument |
| SR2/refresh-remedy-carries-the-refresh-flag | drop `--refresh <receipt>` from the refresh classes' case, or fold them back onto `RefusalMissingAssignment` | refusal-class remedy table test | construct each declared refusal class through `RefusalForClass` and require the two refresh classes to advertise the `--refresh` template while missing-assignment and every other class keep their current commands |
