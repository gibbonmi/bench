# Register the prose mechanics check, the exclusion file, and its fixtures

Blocked by: 01-add-prose-mechanics-check.md
Writes: internal/conformance/prose_mechanics_test.go (new), internal/conformance/registry/registry.go, internal/conformance/registry_test.go, internal/conformance/checks_test.go, internal/conformance/tier_test.go, tests/canary/prose-mechanics/ (new), .bench/prose-exclusions (new), projects/benchkit.md
Line: `opus` / medium under Claude Code, `gpt-5.6-terra` / medium under Codex — registry joins and fixtures at known seams.

## What to build

The conformance check `prose-mechanics` is the registered wrapper over `internal/prose`, with subject `root` and input source `catch-all`. The check joins the seven registries the spec lists, and the profile's conformance table gains its row in this commit. The hostile-skill root writer plants an empty `.bench/prose-exclusions`, so the check reaches the skill file and names it.

`.bench/prose-exclusions` starts with the four permanent rows and one temporary row per later ticket's scope. The temporary rows, each with the reason `migrates in ticket NN`, cover these prefixes:

- ticket 02b: `.bench/BENCH.md`, `.bench/BENCH-reference.md`
- ticket 03: `.agents/skills/bench-craft-spec/references/ste-prose.md`, `.agents/skills/bench-craft-skills/`, `.agents/skills/bench-craft-comments/`
- ticket 04: `.agents/skills/bench-craft-spec/SKILL.md`, `.agents/skills/bench-craft-spec/references/bootstrap-authority.md`, `.agents/skills/bench-craft-seams/`, `.agents/skills/bench-craft-line/`, `.agents/skills/bench-craft-domain/`, `.agents/skills/bench-craft-gate/`
- ticket 05: `.agents/skills/bench-craft-tickets/`, `.agents/skills/bench-craft-tdd/`, `.agents/skills/bench-craft-delegate/`, `.agents/skills/bench-craft-grill/`, `.agents/skills/bench-craft-adr/`
- ticket 06: `.agents/skills/bench-craft-review/`, `.agents/skills/bench-craft-synthesis/`
- ticket 07: `.agents/skills/bench-craft-cli/`, `.agents/skills/bench-craft-design-system/`, `.agents/skills/prototype/`
- ticket 08: `.agents/commands/bench.md`, `.agents/commands/bench-write-spec.md`, `.agents/commands/bench-implement-spec.md`, `.agents/commands/bench-review-implementation.md`, `.agents/commands/bench-final-check.md`, `.agents/commands/bench-debug.md`, `.agents/commands/bench-what-next.md`
- ticket 09: `.agents/commands/bench-drain.md`, `.agents/commands/bench-setup-repo.md`, `.agents/commands/bench-shape-idea.md`, `.agents/commands/bench-deepen.md`, `.agents/commands/bench-assess.md`, `.agents/commands/bench-update-kit.md`
- ticket 10: the thirteen `.agents/skills/bench*/` adapter directories that are not `bench-craft-*`
- ticket 11: `AGENTS.md`, `CLAUDE.md`, `CONTEXT.md`, `.claude/`, `DATA_HANDLING.md`, `SECURITY.md`, `projects/gl-axi.md`, `docs/greenfield-build-sequence.md`, `docs/release-runbook.md`, `docs/reporesident-distillation.md`
- ticket 12: `README.md`, `projects/benchkit.md`
- ticket 13: `ASSESSMENT.md`, `skills-assessment.md`, `docs/adr/`
- ticket 14: `ROADMAP.md`, `roadmap/`
- tickets 15 to 18: each `decisions/` map file, each `decisions/assets/` file, and `decisions/byte-preserving-axi-foundation/`, as those tickets list them
- ticket 19: `specs/inherited-toolchain-environment/`, `tickets/`, `capture/agent-performance/`, `capture/audits/`, `capture/FIXES.md`, `capture/parallel-session-friction.md`, `capture/learnings.md`
- ticket 28b: `capture/session-handoff.md`

A live-tree subset test pins that initial row set as the approved set and asserts the live rows are a subset of it. A live-tree test grades the kit root and is the delegates' focused seam. `specs/asd-ste100-progressive-disclosure/` has no row and passes from this commit.

The new check's canary family holds one `files/`-form fixture per planted red. Each fixture plants at least one `*.md` subject, so a restore that removes it clears the red. The planted reds are:

- a long sentence
- a long paragraph
- an unterminated fence
- a stale row
- a duplicate row
- a glob row
- a missing exclusion file

## Acceptance

- [ ] An excluded prefix is not graded, and a stale, duplicate, or glob row reds with its own message (covers PD18, PD19).
- [ ] Every fixture in `tests/canary/prose-mechanics/` bites through the registered owner and clears on restore (covers PD23).
- [ ] The check over a hostile skill root names the refused path (covers PD24).
- [ ] The live-tree test reds when a temporary row is removed early and is green over the committed tree (covers PD25).
- [ ] The subset test reds when a row outside the approved set is added and is green over the committed file (covers PD32).
- [ ] The registry gains exactly one check, and the profile table names it (covers PD36).
