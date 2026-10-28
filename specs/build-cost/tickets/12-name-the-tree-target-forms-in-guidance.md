# Name the tree-target forms in the guidance and record the changes

Blocked by: 02-refuse-plain-chunk-requirements.md, 05-add-trace-times-and-the-drain-rule.md, 07-stage-the-probe-in-record-ticket.md, 08-print-a-short-lane-pass.md, 09-print-a-short-bitten-mutant-set.md, 10-name-the-tree-target-in-verb-help.md, 11-reuse-a-matching-sealed-build.md
Writes: docs/known-issues.md, .bench/BENCH-reference.md, .agents/skills/bench-craft-synthesis/SKILL.md, .bench/BENCH.md, internal/anchors/registry_build_cost.go (new), internal/anchors/registry_build_cost_test.go (new), internal/anchors/registry_retained_workflow.go, CHANGELOG.md, .bench/commitment.json
Covers: BC66, BC67, BC68, BC69, BC70, BC73

## What to build

The guidance names the `--in` form for three Bench calls:

- `docs/known-issues.md`: "Run it as `bench preflight --in <target> evidence <id> --check-current`." The row's retire column keeps no owner and names the `--in` form, not the exec form.
- `.bench/BENCH-reference.md`: "For a worktree, run `bench probe --in <target> ...`." The old exec probe sentence goes.
- `.agents/skills/bench-craft-synthesis/SKILL.md`: "Take the green verdict from the whole-project gate. That gate is the landing's gate, or `bench gate --in <target>` for a batch that waits for approval."

`internal/anchors/registry_retained_workflow.go` loses its exec gate Require entry, and `registry_build_cost.go` gains the Require entry of the new craft-synthesis sentence, so the sentence has one anchor.
`.bench/BENCH.md` keeps the learning for every plan or gate expansion except a `Writes:` widen. Insert the text ` except a `Writes:` widen` after "record every plan or gate expansion", because the probe-12 command omits that exact text. The file keeps its 192-line budget.

Add the anchors of BC66, BC67, BC68, BC69, and BC73 to `registry_build_cost.go`. BC69 is a Forbid anchor on the retired text "record every plan or gate expansion with `bench learning`", which the probe-12 omission restores. BC70 is review-owned, because ADR 0012 keeps rule wording out of the anchor registry. Write one `CHANGELOG.md` entry for the whole build. `.bench/commitment.json` is in this fence because the spec-stage binding commit sits in the graded range.

## Acceptance

- [ ] The known-issues row names the `--in` form, and its retire column holds no exec form (BC66, BC73).
- [ ] The command notes name `bench probe --in <target> ...` and hold no exec probe form (BC67).
- [ ] The craft-synthesis skill names `bench gate --in <target>` (BC68).
- [ ] `.bench/BENCH.md` holds the widen exception and still records a learning for every other expansion (BC69, BC70).
