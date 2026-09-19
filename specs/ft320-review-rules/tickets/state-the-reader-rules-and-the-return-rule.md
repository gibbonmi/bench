# State the finding reader rules and the axis return rule

Blocked by: none
Writes: .agents/skills/bench-craft-review/references/finding-discipline.md, internal/anchors/registry_review_rules.go (new), internal/anchors/registry_review_rules_test.go (new), internal/anchors/registry_data.go, ROADMAP.md, roadmap/FT320.md
Covers: none

## What to build

This ticket is the FT320 kit edit. Each missing rule produced a wrong finding, a red
gate, or a partial result in the bounded-charge-evidence, calibrated-decisions, or
slicing-closure reviews. The `craft-review` finding discipline gets four rules:

1. A review of gate-anchored prose names the anchor state of each sentence that it proposes to change.
2. A Coverage finding describes the tree before the probes of the axis.
3. An axis reads the seam cell of a row before it judges a review-owned row unmet.
4. An axis return lists the evidence cursors that it fetched.

Each rule gets one anchor row and one independent test expectation. The same source
retires FT320 and refreshes the recommended sequence.

## Acceptance

- [ ] The finding discipline states the four rules, and each rule has an anchor row.
- [ ] If you delete the axis-return rule, `bench test --check docs-currency-workflow` goes red.
- [ ] FT320 leaves `ROADMAP.md` and `roadmap/`, and the sequence no longer names it.
- [ ] The prose lane and the gate pass.
