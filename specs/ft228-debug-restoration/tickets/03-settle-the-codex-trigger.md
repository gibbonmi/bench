# Settle the Codex trigger and remove the inert key

Blocked by: none
Writes: internal/conformance/skills_index_checks_test.go, internal/conformance/registry_test.go, .agents/skills/*/SKILL.md (13 phase adapters), .agents/skills/bench-debug/agents/openai.yaml, tests/canary/skills-index-command-adapters/

## What to build

The Codex half of the invocation-policy settle: the reviewed table, its
Codex-side grading, the `bench-debug` flip with a description that can fire,
and the inert-key removal — landing green in one commit because the check and
the yaml must move together.

- Add the `phaseInvocationPolicy` table beside the check in
  `internal/conformance/skills_index_checks_test.go`: per phase, a
  Claude-model-invocable boolean and a Codex-implicit boolean. Values:
  `bench-debug` true/true (the reviewer's 2026-08-19 settle); `bench`,
  `bench-final-check`, `bench-implement-spec`, `bench-review-implementation`,
  `bench-shape-idea`, `bench-write-spec` true/false; `bench-assess`,
  `bench-deepen`, `bench-drain`, `bench-setup-repo`, `bench-update-kit`,
  `bench-what-next` false/false. This ticket wires the Codex-side grading:
  `openai.yaml` spells exactly the value the row demands (a file spelling
  neither, or an empty file, is red as undeclared), the adapter `SKILL.md`
  frontmatter carries no `disable-model-invocation` key, and a Codex-implicit
  row's adapter description carries no explicit-only phrasing (the marker
  `Use only when the reviewer invokes`). The check stays in
  `hostileSkillReaders` with its reads on the hardened producer path.
- Flip `.agents/skills/bench-debug/agents/openai.yaml` to
  `allow_implicit_invocation: true`. Rewrite
  `.agents/skills/bench-debug/SKILL.md`'s description to a symptom-bearing
  trigger mirroring the command file's own — something is broken, throwing,
  failing, or slow — dropping `Explicit` and `Use only when the reviewer
  invokes`. Delete the `disable-model-invocation` line from all thirteen
  phase-adapter `SKILL.md` files; the other twelve descriptions stay.
- Add two `skills-index-command-adapters` fixtures: `bench-debug`'s yaml
  flipped back to `false`, and the inert key reintroduced into an adapter
  `SKILL.md`. Each fixture's `BASE` names the command file, the adapter
  `SKILL.md`, its `agents/openai.yaml`, and the two guide files
  (`.bench/BENCH.md`, `.bench/BENCH-reference.md`). Register the family once in
  `canaryFixtureFamilyRegistry` (`internal/conformance/registry_test.go`)
  naming `skills_index_checks_test.go` and `checks_test.go`; the five existing
  exact rows stay and override.
- Add in-package cases over constructed roots for the undeclared/empty
  `openai.yaml` and the explicit-only description on a Codex-implicit row.

The contract with ticket 04: ticket 04 extends this table and check with the
Claude-side grading and the two table-completeness directions, re-derived from
the tree by review.

## Acceptance

- [ ] The kit gate is green with `bench-debug` at `allow_implicit_invocation: true`, its rewritten description, and no adapter `SKILL.md` carrying the inert key (IP1, IP10).
- [ ] The flipped-yaml fixture and the reintroduced-key fixture red through the registered owner naming their subjects, and each restore re-runs green (IP2, IP6).
- [ ] The in-package cases red the undeclared yaml and the explicit-only description on an implicit row (IP9, IP10).
- [ ] `TestRegisteredSkillReadersRefuseHostileSkillFiles` passes over the modified check (IP7).
- [ ] The CHANGELOG is untouched by this ticket — ticket 04 documents the settle when the policy is complete.
