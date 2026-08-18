# Learnings — usage journal

## 2026-08-18 — a spec that introduces a phase is red under the stale-command sweep [open]

**What happened.** Staging `specs/bench-front-door/spec.md` (A3: the `/bench` front door
and the what-next → drain rename), the live-root conformance sweep
(`checkStaleCommandReferences`) reported twenty `stale command reference` diagnostics —
every mention of the phase the spec introduces. The sweep derives its valid `/bench-*`
tokens from the files present in `.agents/commands/`, and it walks `specs/` fail-closed,
so any spec whose deliverable is a new or renamed phase names a token that cannot be valid
until the build lands. `/bench` itself passed only because it lacks the `/bench-` prefix.

**Right behavior.** The spec should name its deliverable exactly; the sweep should be able
to tell "a token this staged spec introduces" from "a token that has been removed". I did
not paraphrase the commands around the check or mark the spec `command-currency:
historical` — both hide the deliverable to appease the oracle — and surfaced the red for
the reviewer instead.

**Proposed rule change.** Give the sweep a spec-local allowance: a `Introduces commands:`
line in a staged spec's header (naming the drain phase's slash and Codex forms) makes those tokens valid within that spec's directory only, for as long as
the spec is staged; an implemented spec loses the allowance so a phase promised but never
shipped goes red at retirement. Prove it bites with a fixture where the spec names an
undeclared future command. This is a gate change (`craft-gate`), so it is the reviewer's
to approve — recommended as a light-path ticket landed ahead of the A3 spec commit.
