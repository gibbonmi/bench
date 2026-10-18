# Retire the spec fence section

Blocked by: 04-derive-fence-from-tickets.md
Writes: internal/coverage/rows.go, internal/coverage/citations.go, internal/coverage/citations_test.go, internal/coverage/coverage.go, internal/coverage/fence_section_test.go (new), internal/spec/fences.go (deleted), internal/spec/fences_test.go (deleted), internal/preflight/preflighttest/fixture.go, internal/preflight/preflighttest/reviewfiles.go, internal/preflight/deletion_preflight_test.go, internal/preflight/charge_file_kinds_test.go, internal/preflight/command_bootstrap_test.go, internal/preflight/kit_pin_new_test.go, internal/preflight/evidencecmd/evidence_command_test.go, internal/preflight/evidencecmd/evidence_budget_test.go, internal/preflight/ownership/ownershiptest/repo.go, cmd/bench/preflight_version_test.go, internal/systemtest/owner_landing_fixture_test.go, internal/systemtest/charge_evidence_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/delegated_integration_test.go, internal/reviewrecord/recordtest/fixture.go, .agents/skills/bench-craft-spec/SKILL.md, .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-spec/references/map-only-route.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go
Covers: CR70, CR71, CR86, CR87

## What to build

No code reads a spec fence section. The coverage parser stops parsing the section, and `checkReviewPickup` goes, so `bench coverage --check` grades a spec with or without a section. `spec.FenceTokens` and its tests go with their last reader.

The shared spec fixture writer stops writing the section, and each caller that passed extra fence lines passes them as ticket `Writes:` entries or drops them when the spec folder already authorizes them. The fixture writers that only wrote the section drop it: the version, owner-landing, charge evidence, landing, delegated, review-record, and ownership fixtures. The two system-tagged fixtures run with `BENCH_KIT` set through `bench test --check system`. The coverage test for the ignored section goes into the new test file, because the citations test file is over its structure budget.

The guidance drops the spec fence. The craft-spec template, its slicing paragraph, and its approval table hold no fence section, and the approval table names the derived scope. `/bench-write-spec` drops the fence-section sentences. The needle "After the ticket review accepts, `bench preflight build <slug> --plan-only` grades the whole ticket graph." replaces the old fence-step needle in the ticket-slicing anchor registry and its test. The map-only route drops its fence-section steps. The slicing checks drop their spec-fence sentences, and the `.bench/commitment.json` rule keeps its ticket half: exactly one ticket's `Writes:` line names that file. A forbid anchor reds a restored `## Ownership fences` section in the craft-spec skill.

This spec keeps its own fence section for the pre-FT293 landing broker. Each edited budgeted file holds or shrinks its line count.

Contract for ticket 06: the craft-spec skill no longer owns a fence, so the review skill's owner clause and the delegate rule name the ticket `Writes:` lines.

## Acceptance

- [ ] `bench coverage --check` exits 0 for a spec whose fence section omits the review pickup.
- [ ] No production code declares or calls `spec.FenceTokens`.
- [ ] The craft-spec skill holds no `## Ownership fences` text, and a forbid anchor reds it.
- [ ] `bench-write-spec.md` holds the new plan-only needle in its Process section, and the anchor test pins it.
- [ ] The slicing checks keep the ticket-only `.bench/commitment.json` rule.
- [ ] Each edited budgeted guidance file holds a line count at or below its budget row.
