# Enforcement Verification (FT27)

## #1: What happens when a delegation envelope omits the model field?

Blocked by: —
Type: Grill

### Question
The agent-line guard fails open when `tool_input` carries no model field —
documented as the intended rim, but an omitted model inherits the invoking
session's model, which is precisely the silent-escalation path invariant 2
exists to stop. The counter-rule lives only in craft-delegate prose.

### Answer
Deny (exit 2) on a missing model field when the repo is routed (`.bench/
lines.env` present and parseable), with a message naming the fix: pass the
bound alias from lines.env. Unrouted repos keep the fail-open rim — there is
no binding to enforce. Other degraded branches (malformed hook JSON, absent
envelope) stay fail-open as today; only the specific "routed repo, Agent call,
no model" branch tightens, because it is the one degraded branch that is also
the attack path the guard exists for. **Flagged for veto:** this changes a
documented posture; the residual fail-open branches get recorded in the
postures ADR ([FT31]).

## #2: How does pre-push installation get verified after link?

Blocked by: —
Type: Grill

### Question
`bench doctor` checks the PATH shim but never that `.git/hooks/pre-push`
exists, is bench-managed, or is not diverted by `core.hooksPath`. A fresh
clone silently loses the harness-independent backstop.

### Answer
Doctor gains a pre-push row: red when the hook is absent, not bench-managed
(fingerprint mismatch against the embedded template), or diverted by a
`core.hooksPath` that doesn't contain it; the remedy line is `bench link`
(re-link reinstalls). `bench status` gains a matching low-noise signal
(guards family, severity just below the worktree rows) so the gap is ambient,
not doctor-only. FT10 folds in: the same doctor row fires on the kit repo
itself and names the install path. Rejected: silent self-heal on any bench
invocation — installing hooks as a side effect of unrelated commands violates
least surprise.

## Handoff

1. **Module boundaries.** `internal/lines` owns the verdict change;
   `internal/adopt` (doctor) owns the pre-push row; `internal/status` owns the
   new signal; `.claude/hooks` wiring unchanged (same hook binary path).
2. **Contracts.** AgentLineVerdict: routed + Agent + no model → exit 2 with
   alias guidance; unrouted → exit 0 warn. Doctor: pre-push row green/red with
   remedy; exit code follows doctor's existing red posture. Status: one row,
   fires only when the hook is missing/unmanaged in a routed repo.
3. **Deep vs thin.** lines.go verdict logic is the deep seam (documented rim
   changes); doctor/status rows are thin composition of existing
   fingerprint/manifest helpers.
4. **Black-box assertables.** Hook exit codes and stderr for
   present/absent/malformed model fields in routed and unrouted fixtures;
   doctor stdout rows against temp repos with absent/foreign/managed hooks;
   status row presence.
5. **Gate attachment.** All three seams have existing test families
   (line_routing_checks, doctor tests, runtime_status) — extend, don't invent.
6. **Hostile-input owners.** lines owns malformed JSON envelopes, empty
   strings, whitespace model values; doctor owns foreign (user-authored)
   pre-push hooks — never overwrite or call them managed.
7. **Uncertainty flags.** Where the embedded pre-push template fingerprint
   lives today (link_hook.go) and whether doctor can reuse it without an
   import cycle — implementer verifies before wiring.
8. **Rejected alternatives.** Fail-open on missing model in routed repos;
   self-healing hook installs; a hard gate check for hook presence (the gate
   must stay runnable in worktrees and CI clones where hooks are legitimately
   absent).
9. **Domain watch-outs.** git does not clone hooks — absence after a fresh
   clone is the expected state doctor exists to catch, not an error in link.
   The status signal must fire only on the default branch's primary checkout;
   pool worktrees share the main .git and must not double-report.

Dependency order: n/a — single spec.
