# State the lane in the rules and the profile

Line: fable / high.

Blocked by: 03-commit-through-the-lane-authority.md
Writes: .bench/BENCH.md, .bench/BENCH-reference.md, CONTEXT.md, projects/benchkit.md, docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md (new), CHANGELOG.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/conformance/profile_lane_table_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/registry.go

## What to build

A fresh session reads the rules and knows where the one full grade lives.
Invariant 4 in `.bench/BENCH.md` gains one sentence: "green" means the landing's
gate, and a worktree commit requires a lane pass. The spec's guidance decision holds the exact sentence. The reference's landing shape names the
lane at the worktree commit, so the reference and the verb agree.

`CONTEXT.md` defines `fast lane` and `lane record`, each with an Avoid list. The
Avoid list for `fast lane` refuses "mini gate" and "reduced gate". Follow the
glossary entry shape the file already uses.

The profile carries the kit's lane table beside its phase table, with the same
authoritative argv column. Add the conformance check that compares those profile
rows with the kit's built-in lane argv, so a drifted table turns the gate red.
The built-in lane argv is the one the gate ticket landed.

Write ADR 0017 with Status accepted. It records the authority split: the lane
authorizes a worktree commit, and the gate alone authorizes a landing. It names
no file path and no code snippet. Add the anchor registry marker and row for the
new invariant-4 sentence, and the section anchor for the reference's landing
paragraph. Name the fast lane in `CHANGELOG.md` under the unreleased heading.

## Acceptance

- [ ] OG26 shows the anchor registry requires the invariant-4 sentence, and a removal of that sentence turns the anchor check red.
- [ ] OG27 shows the reference's landing section contains `lane`.
- [ ] OG28 shows `CONTEXT.md` defines `fast lane` and `lane record`, each with an Avoid list.
- [ ] OG29 shows the profile's lane table rows equal the kit's built-in lane argv.
- [ ] OG30 shows an ADR with Status accepted records the authority split and names no file path.
- [ ] OG31 shows `CHANGELOG.md` names the fast lane under the unreleased heading.
