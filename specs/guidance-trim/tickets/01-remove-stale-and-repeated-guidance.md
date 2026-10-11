# Remove stale and repeated text from the always-loaded guidance

Blocked by: none
Writes: .bench/BENCH.md, AGENTS.md, CHANGELOG.md, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ft311_review_dispatch.go
Covers: none

## What to build

`.bench/BENCH.md` and `AGENTS.md` load into every session. On 2026-10-11 the reviewer approved a cut of the stale, contradictory, and repeated text in both files. The reviewer chose to delete the harness auto-memory clause, not to reverse it. The relocation of the build orchestration rules and the light-path routes stays with FT100.

The operating guide gets these corrections:

- Workflow step 5 and invariant 4 say that a commit needs a lane pass and that a landing needs green.
- The guide names no fix-and-gate path and no unrecorded size rule as standing approvals.
- The guide drops the harness auto-memory clause and the kit-only claim that the repo ignores both capture inboxes. It also drops the sentence that put editors and raw Git outside the command boundary.
- The guide states the drain route and the list-or-prose rule once each, and it leaves the board layout to the reference file.
- The closed-decisions rule moves to Roles, and the header, Roles, and the CLI contract lose their repeated clauses.

The working agreement gets these corrections:

- The header states once that the shared rules are canonical in the operating guide. The account of how each harness loads the file and the account of the marker check go.
- The phase-close block requires `bench handoff` and no longer offers a copy-paste prompt instead. It keeps what the agent supplies: the spec path and status, the closed decisions, and the next command.
- The retro rules that `/bench-final-check` owns, the account of section dates, and the manual prose run before a commit go.

The structured-phase declaration keeps its clause list, because the structured-phase contract check reads it. No anchored sentence changes, and no forbidden needle returns. The `Writes:` line names each anchor registry file that names either guidance path, and the ticket edits none of them.

## Acceptance

- [ ] Step 5 and invariant 4 of the operating guide say that a commit needs a lane pass and that a landing needs green.
- [ ] The operating guide contains no fix-and-gate approval, size rule, auto-memory clause, kit-only inbox claim, or raw-Git boundary sentence.
- [ ] The phase-close block of the working agreement requires `bench handoff` and offers no alternative.
- [ ] The anchors of both files, the shared-rule check, and the structured-phase contract check pass.
- [ ] `CHANGELOG.md` has one entry for the operating guide changes.
