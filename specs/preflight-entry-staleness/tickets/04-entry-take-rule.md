# Take a red from landed work at build entry

Blocked by: 01-writes-landed-row.md, 03-map-sources-row.md
Writes: .agents/commands/bench-implement-spec.md, .agents/skills/bench-implement-spec/references/staleness-pass.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, CHANGELOG.md, .bench/commitment.json
Covers: ES43, ES44, ES47, ES50, ES51

## What to build

The build entry of `/bench-implement-spec` takes a red into the staleness pass when every red row is `completion-plan`, `writes-landed`, or `map-sources`, and any other red still stops the phase.
Read the spec sections "Guidance" and "Records" and the coverage rows ES43, ES44, ES47, ES50, and ES51 before you start.
This ticket edits guidance prose, so its author runs at high effort under the leverage override of `craft-line`.

Make these changes:

- In `.agents/commands/bench-implement-spec.md` line 15, replace only the take-rule sentence with the sentence that the spec states. Keep the next sentence `Any other red stops the phase.` and every other sentence of the paragraph byte for byte.
- The file holds 97 lines against its 97-line guidance budget (`projects/benchkit.md:553`), so add no line and no sentence.
- The new text must hold no `fence-writes`, `-closure`, or `check-floor`, because Forbid anchors refuse each one. Keep the Require needles of `internal/anchors/registry_data.go:92-93` and `internal/anchors/registry_ft311_preparation.go:8-10`, and the canary `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight/MUTATE.json`, unchanged.
- Add one Forbid row of the old exact sentence to `internal/anchors/registry_ticket_passes.go`. Add the same Forbid rule, written independently of the registry, to `TestTicketSlicingAnchors`, as the earlier retired sentence at `registry_ticket_passes_test.go:31` shows. Give the rule a diagnostic that no other rule of that harness call uses.
- Add one sentence under `Orchestrator review` of `staleness-pass.md`: a taken `map-sources` red repoints each named record to a source that the published tree holds, and a path outside the signed-off `Writes:` lines takes the widen rule of `.bench/BENCH.md`.
- Add one entry to `## [Unreleased]` of `CHANGELOG.md` that names `writes-landed`, `map-sources`, and the new take rule. Keep the two existing Require-anchored entries.
- Run `bench anchors .agents/commands/bench-implement-spec.md` after the edit and compare the result with the anchors before the edit.

`.bench/commitment.json` is on the `Writes:` line for the spec-stage binding commit only. The author does not edit it.

## Acceptance

- [ ] Line 15 of `.agents/commands/bench-implement-spec.md` holds the new take-rule sentence and, directly after it, `Any other red stops the phase.` (ES43, ES44).
- [ ] `TestTicketSlicingAnchors` reds a restore of the old take-rule sentence (ES51).
- [ ] The `Orchestrator review` section of `staleness-pass.md` states the repoint of a taken `map-sources` record and names the widen rule (ES47).
- [ ] The `## [Unreleased]` section of `CHANGELOG.md` names both rows and the take rule (ES50).
