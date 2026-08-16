# Harden every skill-file reader

Blocked by: none
Writes: internal/bounds/classify.go, internal/bounds/classify_test.go, internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go, internal/conformance/checks_test.go, internal/conformance/skills_index_checks_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/docs_workflow_helpers_test.go, internal/conformance/guidance_token_sweep_test.go, internal/conformance/prose_budget_test.go, internal/anchors/registry.go, internal/anchors/registry_test.go

## What to build

Add the permanent no-follow control-record classification beside the existing
follow-live-symlink API, then make every gate reader of
`.agents/skills/*/SKILL.md` consume it before opening bytes. The shared result retains
the existing state/data/stream/reason vocabulary, applies `ControlRecordLimit`, and
refuses both live and dangling links without changing existing callers that require
follow semantics. This ticket publishes that exact contract for later reference and
payload tickets; a generic conformance helper may select it only for the spec's three
producer classes.

Migrate `load-validity-metadata`, `skills-index-command-adapters`,
`docs-currency-workflow` including `anchors.EvaluateGroup`, `line-routing`,
`guidance-prose-budgets`, and `axi-query-registry` with the classifier. This is one
expand-and-migrate commit: a classifier-only cut is horizontal, while splitting the
SKILL migrations strands HI12 red because one hostile fixture reaches every registered
reader and any remaining direct opener can hang the full gate.

## Acceptance

- [ ] `(covers HI1)` The shared no-follow table classifies absent, empty, exact-limit,
  oversized, invalid-UTF-8, live/dangling symlink, FIFO, device, socket, directory,
  and regular-file cases exactly, without changing the existing follow API.
- [ ] `(covers HI12)` Every named registered conformance path, including AXI guidance
  and anchor evaluation, completes on hostile `SKILL.md` with an attributed refusal;
  no later reader opens the FIFO or follows either symlink form.

