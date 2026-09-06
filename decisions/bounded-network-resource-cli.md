# Bounded network, resource, and CLI behavior (FT87)

Status: shaping

## Destination

One explicit policy bounds every Bench-initiated network attempt, subprocess,
read, and output. `BENCH_OFFLINE=1` provably prevents all of them. Repair is
explicit and manifest-pinned. The CLI has one argument grammar. Sources:
`RR:A-14`, `RR:C-06`, `RR:C-07`, `RR:C-09`–`RR:C-12`; `RC:H-04`, `RC:M-02`.

Closed 2026-07-21 under the reviewer's blanket approval of the worker's
recommendations; contestable calls are marked **[veto]** for post-hoc review.

## Notes

## Decisions so far

- [What does `BENCH_OFFLINE=1` mean, exactly?](bounded-network-resource-cli/tickets/1.md): `BENCH_OFFLINE=1` is the master no-network switch.
- [Where does the one bounds policy live and what are its defaults?](bounded-network-resource-cli/tickets/2.md): A single Go policy package (`internal/bounds` or equivalent) single-sources every Go-side named bound.
- [What is the hardened repair posture?](bounded-network-resource-cli/tickets/3.md): Repair becomes explicit: it never runs as a silent side effect of resolution failure.
- [What replaces the implicit worktree `git fetch origin`?](bounded-network-resource-cli/tickets/4.md): `Acquire` drops the implicit fetch.
- [How is model discovery bounded and parallelized?](bounded-network-resource-cli/tickets/5.md): The three providers query concurrently (goroutines joined by a WaitGroup); `runCommand` moves to `exec.CommandContext` under the policy deadline.
- [What does default `bench outline` emit?](bounded-network-resource-cli/tickets/6.md): Default output is a bounded summary: total file and symbol counts, plus the first N symbol rows (recommended N=200 **[veto.
- [What is the one CLI argument grammar?](bounded-network-resource-cli/tickets/7.md): A small shared parsing helper (hand-rolled, beside the existing `toon.Usage` helpers; no third-party CLI framework) owns the grammar.
- [How do capability skips become evidence, and deadlines decouple?](bounded-network-resource-cli/tickets/8.md): A shared capability-skip helper replaces bare `t.Skip` in security-relevant tests, emitting a recognizable structured skip line.
- [What is the one user-facing identity and complete package metadata?](bounded-network-resource-cli/tickets/9.md): The user-facing identity is **Bench**.

## Not yet specified

- Whether the shift loop should ever pass `--refresh` by default in CI-like
  environments (today: never implicit).
- Whether repair's shipped digest manifest and FT83's release evidence index
  can later share a generator (single-source candidate once both exist).

## Out of scope

- Host firewalls, egress enforcement, IAM, endpoint controls — outside the
  repository-controlled scope.
- Shift failure/evidence semantics (`RC:H-05`, FT71) and objective data
  exposure (`RR:C-08`, shipped).
- Consumer/maintainer capability separation (FT85) and transactional link
  lifecycle (FT84).
- Gate wall-clock reduction (FT91) — the gate timeout here is a ceiling, not
  a speed fix.
- Reopening the npm distribution identity (`redbench`, ADR 0004).

## Spec-writer discretion

## Sources
