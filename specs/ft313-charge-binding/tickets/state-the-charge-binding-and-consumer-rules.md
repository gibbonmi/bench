# State the charge binding rules and the transport consumer protocol

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-delegate/references/charge-evidence-format.md, internal/chargeevidence/reference.go, internal/anchors/registry_charge_binding.go (new), internal/anchors/registry_charge_binding_test.go (new), internal/anchors/registry_data.go, ROADMAP.md, roadmap/FT313.md, roadmap/FT318.md
Covers: none

## What to build

This ticket is the FT313 kit edit, with the FT321 rules that the drain folded into it. The
delegation discipline already states the baseline artifact, the row-to-test binding, the
pre-change replay, the attempt count, and the skip-ownership check. This ticket adds the
missing rules:

- The coordinator ticks the charge list against the ticket `Writes:` line.
- A root-privilege guard routes through the capability seam.
- A repair fence is the chunk union plus the review-named paths.
- A repair on a frozen sibling uses an integration assignment from `main`.
- A completion-plan probe names the check that detects its mutation.
- A stated confidence freezes at return time.
- The final repair attempt checks every axis for duplicated facts, and the cap requests an evidence-scoped extension.
- A rebase repeats verification, and an adoption repair covers the contradiction class.
- Delegation and independent review keep separate costs.

The generated charge evidence reference gets a `Consumers` section through its generator,
because the shipped file must equal the generator output. The section names the decoder
rule, the hook constraint, and the sandbox constraint.

The record-content items move to FT318, which owns the review record writer. Each new
sentence gets one anchor row and one independent test expectation.

## Acceptance

- [ ] The delegation discipline states each added rule, and each rule has an anchor row.
- [ ] `FormatReference` renders the `Consumers` section, and the shipped reference equals its output.
- [ ] If you delete the sandbox sentence, `bench test --check docs-currency-workflow` goes red.
- [ ] FT313 leaves the roadmap, and FT318 carries the record-content residual.
- [ ] The prose lane and the gate pass.
