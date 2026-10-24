# State one plan commit and one build preflight for each dispatch

Blocked by: 04-grade-the-floor-per-chunk.md
Writes: .agents/commands/bench-implement-spec.md, .bench/BENCH-reference.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_data.go
Covers: RC31, RC32, RC33, RC84

## What to build

Review chunk: RC-C3.

The implement phase tells the orchestrator to write one plan commit for each repair and to run one build preflight for each author or repair dispatch.

State the repair order in `.agents/commands/bench-implement-spec.md`. `bench record assignment` comes first, because it refuses a dirty spec or ticket (RC30, ticket 04). Then `bench spec widen` and any floor edit follow, and one plain `bench commit` carries them all. A repair runs each requirement that its ticket owns once, at its own commit.

Name the charge preparation (`bench preflight build <slug> --charge`) as the one build preflight of an author or repair dispatch. A plan commit that a dispatch follows uses `bench commit` without `--preflight-build`. A plan commit that no dispatch follows keeps `bench commit --preflight-build <slug>`. Change the line 62 clause "Each plan commit uses `bench commit --preflight-build <slug>`, and no dispatch starts until its `commit-chain` line reads green." The commit chain code does not change. Line 52 runs a build preflight after each ticket commit; word the new rule so that line 52 does not read as a second preflight of a dispatch.

The file is at 95 of 95 lines, and each paragraph is one line. Change sentences inside the existing paragraphs, or remove a paragraph for each paragraph you add. Keep these sentences verbatim, because anchors pin them: line 21 (`internal/anchors/registry_ft311_preparation.go:8`), line 52 and line 64 (`internal/anchors/registry_chunk_chain.go:11,14`), and the line 68 sentence "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`." (`registry_chunk_chain.go:23`). Lines 17 and 19 (`internal/anchors/registry_data.go:92-93`) do not change. New wording must avoid the Forbid needles of `registry_chunk_chain.go:12-13,15-22,24-33` and `internal/anchors/registry_ticket_passes.go:23-24`.

Put the command detail that does not fit in `.bench/BENCH-reference.md` and in the bounded repair policy, which have no line budget.

`.bench/commitment.json:950-951` (criterion `FT380-plan-commit.preflight-build`) still requires `--preflight-build` for each plan commit. That ledger text changes only through `bench commitment plan`, so this ticket does not edit it.

## Acceptance

- [ ] RC31: the implement phase names the charge preparation as the one build preflight of a dispatch.
- [ ] RC32: the implement phase keeps `bench commit --preflight-build <slug>` for a plan commit that no dispatch follows.
- [ ] RC33: each edited guidance file stays inside its line budget (`bench test --check guidance-prose-budgets`).
- [ ] RC84: the implement phase orders the repair plan edits: `bench record assignment`, then `bench spec widen` and any floor edit, then one `bench commit`.
- [ ] `bench test --check root-conformance` stays green, so each anchored sentence still matches.
