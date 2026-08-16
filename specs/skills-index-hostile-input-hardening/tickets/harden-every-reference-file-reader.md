# Harden every reference-file reader

Blocked by: harden-every-skill-file-reader.md
Writes: internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go, internal/conformance/checks_test.go, internal/conformance/skills_index_checks_test.go, internal/conformance/docs_workflow_checks_test.go, internal/conformance/docs_workflow_helpers_test.go, internal/anchors/registry.go, internal/anchors/registry_test.go

## What to build

Consume the published no-follow classification contract for every
`.bench/BENCH-reference.md` opener in skills-index, command-adapter/docs conformance,
and `anchors.EvaluateGroup`. Preserve caller-owned diagnostics: missing is
unverifiable, present zero bytes are empty, and oversized, invalid-UTF-8, or wrong-type
states are distinct attributed refusals. Every refusal blocks `Write` without changing
the reference. The generic helper continues its prior behavior for paths outside the
spec's three producer classes.

All reference readers and the hostile composition fixture land together because a
thinner per-consumer migration strands HI13 red: `docs-currency-workflow` evaluates
both bespoke reference checks and registered anchors, so one remaining direct open can
hang the same full-gate run after skills-index itself has refused.

## Acceptance

- [ ] `(covers HI4)` Missing, empty, oversized, invalid-UTF-8, and wrong-type reference
  fixtures produce their distinct attributed refusals, and every refused write
  preserves the original bytes.
- [ ] `(covers HI13)` Registered skills-index/adapter and docs-workflow checks,
  including `anchors.EvaluateGroup`, complete on hostile reference fixtures without a
  direct reopen or FIFO hang.

