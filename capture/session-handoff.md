# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `a0fb98e`, 2 dirty paths, 0 unpushed commits
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/bench-preflight/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `d57aa5d` — stale, work tree `8b00a1b`

## State

Spec B (`specs/bench-preflight/spec.md`) is staged and reviewer-signed
(2026-08-11): `bench preflight <build|review> <slug>`, the start-oracle over
artifacts-vs-reality from the Pocock map's ticket #7. Five stories (1–3
opus/medium, 4 fable/high, 5 opus/medium), 23-row coverage map, check-4
predicate reviewer-confirmed as row-ID scan with tag-scoped set equality. A
Codex falsification pass (`gpt-5.6-sol`/high) ran; its accepted findings are
folded into the spec. Spec A (lifecycle removal) is implemented and landed;
`axi-coherent-diff`, `axi-query-disclosure`, and `single-build-serial-gate`
are parked pre-reshape specs awaiting re-rank, not active work.

Discovered while falsifying Spec B: `bench prep-release` is red on today's
main — the subcommand-routing conformance check parses dispatch surfaces
(`commands` map, `run()`) that the `commandRegistry` refactor replaced, so
root conformance reports every command as no-longer-dispatched. The dev gate
never runs root conformance (only prep-release sets the env), so dev green is
honest. The checker repair is Spec B's story 5. Co-discovered, not taken:
stale decision-map citations (`decisions/gate-budget.md`,
`decisions/spec-build-review-gate-cadence.md` — retirement already queued for
the drain) and an injected-port derivation red for `internal/canary` (parked
in `capture/IDEAS.md`).

Open reviewer decisions, none blocking:
- `bench worktree clean --apply` still preserves dirty work into
  `refs/bench/recovery/`, which the next resume sweeps — recommend a
  follow-up making explicit clean refuse a dirty removal instead.
- `projects/benchkit.md` ~311–317 still describes lifecycle cadence and
  `capture/audits/injected-interface-composition.md` cites deleted symbols;
  both fold into Spec C.

## Next command

`/bench-implement-spec` for `bench-preflight` in a fresh mid-tier session.
Then Spec C (doctrine adoption, map #4/#5/#6/#8/#9/#10) from
`specs/remove-spec-build-lifecycle/decisions/pocock-alignment.md`, re-reading
that map's external sources first. Then one `/bench-what-next` drain over the
finished tree (roadmap re-verdicts incl. FT128/FT173/FT184/FT185, mooted and
malformed learnings, parked ideas incl. the injected-port red, retirement of
the `spec-build-review-gate-cadence` and `parallel-session-landings` maps,
and re-rank of the three parked specs).
