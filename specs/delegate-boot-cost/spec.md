# Cut the delegate boot cost

Status: staged

Decision source: reviewer-confirmed conversation, 2026-09-11. The reviewer closed two forks that day. The Codex side closes as a measured no-op. Every phase stays model-invocable, the descriptions are trimmed, and a length check grades them.

Verification log: 1 iteration(s) — the opus/medium round returned BLOCK with two blocking and ten fold findings. The author folded all twelve, and reviewer sign-off is pending.

## Problem

A Claude Code delegate pays a fixed boot cost before it reads one line of the tree.
A general-purpose delegate that replies "ok" with no tool use reports 38,430 tokens.
An Explore delegate that does the same reports 8,834 tokens.
The gap is the tool schemas and the system prompt that the full tool set brings.
Bench runs every review axis and every diagnostic consultation on the full set.

The main session also pays for the skill listing on every turn.
The 31 shipped skill descriptions total 7,631 collapsed characters, and the longest is 377.
Six of the 13 command descriptions also exceed 250 characters, and the longest is 404.
No budget bounds a description, so the listing grows with every skill edit.

Codex exec has no boot problem. One `codex exec` that replies "ok" reports 7,258 tokens.
The count is the same with the docs MCP server disabled.

## Solution

Ship two Bench agent types in the Claude adapter: a reviewer and a writer.
Each declares the tools it needs and nothing else, and neither declares a model.
The charge still passes the tier token, so the agent-line guard keeps its verdict.
The delegate skill routes every review, diagnostic, and user-directed write delegation to these types.
A conformance check pins the agent files, and the package allowlist ships them.

Add a reviewer-owned description budget to the project profile.
A conformance check grades every shipped skill and command description against it.
Trim every description over the budget and keep the leading words that trigger invocation.

Record the Codex measurement and change nothing on that side.

## User stories

Line: opus / high.
Implementation-line reason: the hardest material chunk is the description trim, because guidance prose steers every later session, and the leverage override routes it mid + high. The check seams are exact, with a canary precedent per check, and the gate covers them. Seam uncertainty is low.
Harder chunks: C2.

Cheap Claude delegates

1. As a coordinator, I want a Bench reviewer agent type, so that a review delegate boots without the tool schemas it never uses.
2. As a coordinator, I want a Bench writer agent type, so that a user-directed write delegate keeps the edit tools and drops the rest.
3. As a coordinator, I want no model on an agent type, so that the charge passes the tier token and the guard verdict holds.
4. As a reviewer, I want each agent type to declare its tools, so that no agent file inherits the full tool set by omission.
5. As a reviewer, I want the reviewer type without the spawn, publish, and question tools, so that a delegate cannot spawn, publish, or block.
6. As a consumer, I want the link verb to ship the agent types, so that a linked repo gets the same delegate surface.
7. As a release owner, I want the agent files in the package allowlist, so that the npm tarball carries them.
8. As a coordinator, I want the delegate skill to name the agent types, so that a charge uses them instead of the general-purpose type.
9. As a reviewer, I want the boot cost of each agent type measured and recorded, so that the cut rests on an observed number.
10. As a coordinator, I want a fork reserved for parent-context work, so that a fork's inherited boot is a choice.
11. As a reviewer, I want the agent-line guard unchanged, so that the tier binding remains the one enforcement of the line.
12. As a teammate, I want the Claude adapter README to describe the agents directory, so that a cold reader finds the surface.

Lean skill listing

13. As a reviewer, I want a description budget in the project profile, so that the listing's size has one owner.
14. As an agent, I want each skill description within its budget, so that the listing costs fewer tokens on every turn.
15. As an agent, I want each command description within its budget, so that the phase entries cost fewer tokens on every turn.
16. As a reviewer, I want a missing description reported, so that an unlisted skill cannot pass as lean.
17. As a reviewer, I want an over-budget description named with its count, so that the repair names the file and the excess.
18. As a reviewer, I want a folded description reported, so that the count always reads one value line.
19. As a reviewer, I want a symlinked or special subject refused unread, so that the check reads only regular files.
20. As a reviewer, I want a broken budget table to fail closed, so that a parse fault cannot report a clean tree.
21. As a reviewer, I want a profile edit to run the check in the lane, so that a raised budget grades in the same commit.
22. As an agent, I want every trimmed description to keep its leading words, so that invocation still fires on the same triggers.
23. As a reviewer, I want every phase to stay model-invocable, so that a phase chain can call the next phase.
24. As a reviewer, I want the craft-skills skill to name the budget's owner, so that an author finds the number where it lives.

Codex

25. As a reviewer, I want the Codex exec boot measured, so that a profile change rests on evidence.
26. As a reviewer, I want the Codex side closed with that measurement, so that no unmeasured profile ships.

## Implementation decisions

The Claude adapter gains a real agents directory beside its README and settings.
Codex has no agents surface, so the directory is not a symlink into the portable tree.
The consumer payload carries the directory as a tree row, and the package files list mirrors it.
The required-assets census names both agent files.
The link plan needs no code change, because the link destination already accepts the adapter prefix.

Each agent file carries three frontmatter keys: the name, the description, and the tool list.
The reviewer lists the read tools, the search tools, and the shell tool.
The writer lists the same tools plus the two edit tools.
Neither file declares a model, a skill preload, or a memory directory.
A charge names the agent type and passes the tier token, exactly as it does today.

The `claude-agent-definitions` check grades every adapter agent file whose basename starts with `bench-`.
It reports a name that differs from the basename, an absent or empty tool list, and a declared model.
It also reports a spawning, publishing, or blocking tool, and a tool list without the read tool or the shell tool.
It also reads the delegate skill in both directions: every graded agent is named there, and every `bench-` agent named there exists.
A consumer's own agent files are outside the graded set.

The delegate skill states the routing rule.
A review axis or a diagnostic consultation runs as the reviewer type.
A user-directed write delegation runs as the writer type.
The general-purpose type is not a Bench delegate type.
A fork runs only for work that needs the parent's context.

The project profile gains a `Skill description budgets` section with a subject-and-limit table.
The table has two glob rows: one for every skill file and one for every command file, at 250 characters each.

The `skill-description-budgets` check parses that section with the profile-section and markdown-row helpers the prose budget check uses.
It reads the value through the skills index frontmatter reader, which returns the first value line only.
A description that continues on an indented second line therefore reports as folded, and no second frontmatter parser appears.
It counts the runes of the value after it collapses every run of Unicode white space to one space.
The limit comes only from the table, so a lowered cell lowers the limit the diagnostic reports.

It reports a file with no description key.
It reports a description over its limit, with the count and the limit.
It reports a subject that is a symlink or a special file.

A missing heading or a malformed row returns diagnostics and no policy.
The check joins the benchkit-profile input, so the lane runs it on a profile edit.

The invocation policy keeps every current row.
The description trim edits fourteen skill files and six command files, and no other text.

The Codex side records the measurement in this spec and ships nothing.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| C1 / 1-ship-bench-agent-types.md, 2-pin-bench-agent-definitions.md | Claude delegates run on the two Bench agent types with the tool sets the check pins, shipped by link and package | DB1, DB2, DB3, DB4, DB5, DB6, DB7, DB8, DB9, DB10, DB11, DB12, DB13, DB27, DB28 | `bench test --check claude-agent-definitions`, `bench test --check package-shipped-surface`, `go test ./internal/conformance ./internal/adopt ./internal/lines ./internal/packagesurface` | no |
| C2 / 3-trim-skill-descriptions.md, 4-grade-description-budgets.md | Every shipped description sits within a profile-owned budget that the gate grades | DB14, DB15, DB16, DB17, DB18, DB19, DB20, DB21, DB22, DB23, DB24, DB25, DB26, DB29, DB30 | `bench test --check skill-description-budgets`, `bench test --check conformance-canary-families`, `go test ./internal/gate ./internal/conformance` | yes |

## Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"C1","tickets":["1-ship-bench-agent-types.md","2-pin-bench-agent-definitions.md"],"verification":[{"id":"tests","command":"go test ./internal/adopt ./internal/lines ./internal/packagesurface ./internal/conformance","probe":"omit the model-declared diagnostic in the check"},{"id":"checks","command":"bench test --check claude-agent-definitions"}]},{"id":"C2","tickets":["3-trim-skill-descriptions.md","4-grade-description-budgets.md"],"verification":[{"id":"tests","command":"go test ./internal/gate ./internal/conformance","probe":"omit the white-space collapse before the rune count"},{"id":"checks","command":"bench test --check skill-description-budgets"}]}],"final_verification":[{"id":"acceptance","command":"bench gate"},{"id":"integration","command":"bench test --check system"}]}
```

## Testing decisions

- A good test drives one check over a small fixture tree and asserts the exact diagnostic text.
- The two new checks receive canary fixtures, with the prose budget family as prior art.
- The package allowlist rows receive the shipped-surface check and the required-assets census, with the follow-on guard row as prior art.
- The guard row extends the existing subagent-type test in the lines package.
- The gate observes the feature through `bench test --check` for each new check and through the conformance meta checks.
- The boot numbers are review-owned, because the gate cannot spawn a live delegate.

### Seam diagram

    trigger: bench gate, or bench test --check <name>
        │
        ▼
    kit tree  ──▶  [ claude-agent-definitions ]  ──▶  diagnostics
                      ◀ tests attach here: a fixture tree with one fault per canary
    kit tree  ──▶  [ skill-description-budgets ]  ──▶  diagnostics
                      ◀ tests attach here: a fixture profile table plus one over-budget file

    trigger: Agent tool call with subagent_type bench-reviewer
        │
        ▼
    envelope  ──▶  [ AgentLineVerdict ]  ──▶  allow or deny
                      ◀ tests attach here: the subagent-type table in the lines package

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| DB1 | 1, 2 | The check reports a Bench agent file whose name differs from its basename. | `claude-agent-definitions` check, canary `name-mismatch`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A stale name registers the file under the wrong type. |
| DB2 | 4 | The check reports a Bench agent file with no tool list or an empty one. | canary `tools-absent`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An absent list inherits every tool, which is the cost this spec removes. |
| DB3 | 3 | The check reports a Bench agent file that declares a model key. | canary `model-declared`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A frontmatter model invites a charge with no model field, and the guard denies that charge. |
| DB4 | 5 | The check reports a Bench agent file whose tools include the Agent, Artifact, or AskUserQuestion tool. | canary `spawning-tool`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Those schemas are the boot cost. |
| DB5 | 4 | The check reports a Bench agent file whose tools omit Read or Bash. | canary `shell-tool-absent`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A delegate without the shell tool cannot run a worktree exec review. |
| DB6 | 6, 7 | The payload carries the agents tree row, and the package files list mirrors it. | `internal/conformance/package_shipped_surface_test.go` (`TestAllowlistSourceExists`) as precedent | The link plan and the tarball derive from these rows, so a missing row ships no agents. |
| DB27 | 7 | The required-assets census names both agent files, and the pack census reds a tarball without them. | `internal/conformance/package_core_checks_test.go` (`checkNpmPackAssets`) over `RequiredPackAssets`, with `internal/packagesurface/assets_test.go` (`TestRequiredPackAssetsIncludeFollowOnGuard`) as precedent | The shipped-surface check derives files entries only for kit-only rows, so the pack census is the seam that reds a missing agent file. |
| DB7 | 6 | The link plan for a fixture kit holds the reviewer agent file as a file entry at its adapter path. | new `TestLinkPlanShipsClaudeAgents` in the adopt package over `buildLinkPlan`, which no adopt test calls today | The plan is the only path from payload row to linked repo. |
| DB8 | 8 | The check reports a graded agent basename that the delegate skill never names. | canary `agent-unnamed-in-skill`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An agent nobody is told to use is dead weight. |
| DB9 | 8 | The check reports a `bench-` agent named in the delegate skill that has no file. | canary `skill-names-missing-agent`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A rule that names a missing type sends the charge back to the general-purpose type. |
| DB10 | 9 | A reviewer-type delegate that replies "ok" with no tool use reports under 15,000 tokens on its usage line. | review-owned, recorded in the retro | The whole item exists for this number. |
| DB11 | 9 | A writer-type delegate that replies "ok" with no tool use reports under 16,000 tokens on its usage line. | review-owned, recorded in the retro | The writer keeps two more schemas and must stay near the reviewer. |
| DB12 | 11 | The diff under the lines package touches only the subagent-type test table, which gains both Bench basenames. | review-owned, with `internal/lines/lines_agentline_test.go` (`TestSubagentTypeNeverImpersonatesAFork`) as the table | The verdict treats every non-fork type alike, so a production edit there is a second enforcement. |
| DB13 | 10 | The delegate skill states the fork rule. | review-owned | Prose with no reader is not a rule. |
| DB28 | 12 | The Claude README describes the agents directory. | review-owned | A cold reader finds the surface through the README alone. |
| DB14 | 13, 20 | The check reports a profile with no `Skill description budgets` heading and returns no policy. | `skill-description-budgets` check, canary `budget-table-missing`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An unparsed policy would grade a clean tree. |
| DB15 | 14, 17 | The check reports a skill file whose description exceeds its limit, with the file, the count, and the limit. | canary `over-budget-description`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The diagnostic is the repair instruction. |
| DB16 | 15 | The check reports a command file whose description exceeds its limit. | canary `over-budget-command`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Commands are half the listing. |
| DB17 | 16 | The check reports a graded file with no description key. | canary `description-missing`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An absent description is not a lean one. |
| DB18 | 18 | A description that continues on an indented second line reports a one-line diagnostic. | canary `description-folded`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The frontmatter reader returns the first line, so a silent count would understate a folded value. |
| DB19 | 19 | A symlinked skill directory returns a refusal diagnostic and no count. | Go test with `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetRefusesASymlinkedSkillDirectory`) as precedent | The check must not follow a link out of the tree. |
| DB29 | 19 | A special file at a subject returns a refusal diagnostic and no count. | Go test with `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetRefusesNonRegularSubjects`) as precedent | A FIFO at a subject would block the gate in open. |
| DB30 | 13 | A lowered budget cell lowers the limit the diagnostic reports. | Go test with `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsComeFromTheProfileTable`) as precedent | A hard-coded 250 beside a heading probe passes every other row. |
| DB20 | 21 | The benchkit-profile lane selects the new check beside the prose budget check. | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) case PL32 | Without it a raised budget grades only at the next full gate. |
| DB21 | 14, 15 | The live tree passes the check with every description at or under 250 characters. | `bench test --check skill-description-budgets` over the kit root, through `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | This row is the trim itself. |
| DB22 | 22, 24 | Each trimmed description keeps the leading words its body repeats, and craft-skills names the profile table. | review-owned | A trim that drops a trigger word breaks invocation silently. |
| DB23 | 23 | The invocation policy keeps every current row, so no command gains the disable key. | existing `skills-index-command-adapters` check, canary `command-invocation-disabled-against-policy`, through `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A flip reds the existing check. |
| DB24 | 13 | Each new canary family has at least one fixture directory. | `conformance-canary-families` meta check, through `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | A check with no family is not registered. |
| DB25 | 17 | A description with multibyte characters counts runes and not bytes. | new `TestSkillDescriptionBudgetCountsRunes` in the check's file | A byte count reds a correct description. |
| DB26 | 25, 26 | The Further notes record the two Codex boot runs and their equal counts. | review-owned | The no-op rests on this record. |

### Edge inventory

- **Won't handle** a Codex exec profile — the boot measures 7,258 tokens with and without the MCP server, and the reviewer recipe keeps its current form.
- **Won't handle** the built-in Explore and Plan agents — they boot at 8,834 tokens and skip the instruction files, and the fan-out search charge keeps them.
- **Won't handle** a skill preload key on an agent file — a charge names its inputs by path, and a preload adds tokens to every boot.
- **Won't handle** a memory key on an agent file — a delegate has no conversation memory by design, and the charge carries what it needs.
- **Won't handle** an instruction-file opt-out for a custom agent — the harness has no such key, and only the Explore type skips them.
- **Won't handle** a consumer's own agent files — the check grades only basenames that start with `bench-`, and a consumer's `.claude/agents/` keeps its other files.
- **Won't handle** a `when_to_use` key — no shipped file carries it, and the check counts the description key alone.
- **Won't handle** an invocation flip on any phase — the reviewer kept every phase model-invocable on 2026-09-11, and DB23 pins it.
- **Won't handle** the memory index prune — the index is outside the repo.
- **Won't handle** an unterminated frontmatter block — the frontmatter reader returns no value, and DB17 reports the missing description.
- **Won't handle** a description key with an empty value — DB17 reports it as missing.
- **Won't handle** an absent agents directory — the delegate skill names both agents, so DB9 reports each as missing.
- **Won't handle** an empty agents directory — the same DB9 reports each named agent as missing.

## Ownership fences

- `.claude/agents/bench-reviewer.md` (new)
- `.claude/agents/bench-writer.md` (new)
- `.claude/README.md`
- `.bench/consumer-payload.json`
- `package.json`
- `.bench/structure.budgets`
- `internal/packagesurface/assets.go`
- `internal/adopt/link_plan_test.go` (new)
- `internal/lines/lines_agentline_test.go`
- `internal/conformance/claude_agent_definitions_test.go` (new)
- `internal/conformance/skill_description_budgets_test.go` (new)
- `internal/conformance/checks_test.go`
- `internal/conformance/registry/registry.go`
- `internal/conformance/registry_test.go`
- `internal/conformance/tier_test.go`
- `internal/gate/lane_select_test.go`
- `projects/benchkit.md`
- `tests/canary/claude-agent-definitions/` (new)
- `tests/canary/skill-description-budgets/` (new)
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/package-core-guard/`
- `tests/canary/load-validity-metadata/`
- `tests/canary/guidance-prose-budgets/`
- `tests/canary/line-routing/`
- `tests/canary/row-next-grammar/`
- `tests/canary/skills-index-command-adapters/`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-skills/SKILL.md`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/skills/bench-craft-grill/SKILL.md`
- `.agents/skills/bench-craft-synthesis/SKILL.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `.agents/skills/bench-craft-research/SKILL.md`
- `.agents/skills/bench-craft-adr/SKILL.md`
- `.agents/skills/bench-craft-domain/SKILL.md`
- `.agents/skills/bench-craft-gate/SKILL.md`
- `.agents/skills/bench-craft-seams/SKILL.md`
- `.agents/skills/bench-craft-cli/SKILL.md`
- `.agents/skills/bench-craft-review/SKILL.md`
- `.agents/skills/bench-craft-line/SKILL.md`
- `.agents/skills/prototype/SKILL.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-assess.md`
- `.agents/commands/bench-update-kit.md`
- `.agents/commands/bench-drain.md`
- `.agents/commands/bench-deepen.md`
- `.agents/commands/bench-debug.md`
- `specs/delegate-boot-cost/`
- `reviews/delegate-boot-cost.md`
- `capture/retros/delegate-boot-cost.md` (new)

The trim ticket edits only the description line of each listed skill and command file.
The canary prefixes enter the fence because the preflight names each fixture that pins a touched path, and the tickets carry the fixture names.
The delegate skill has ten lines of headroom under its budget, and the craft-skills skill has eight.

## Out of scope

- A Bench verb that records a delegate's usage line as evidence — 6 edits, 2 gate runs. It needs its own spec, because the usage line has no stable producer in the tree.
- A description budget on the reference files under each skill — 3 edits, 1 gate run. The references never enter the listing.
- A Codex exec profile — 2 edits, 1 gate run. The measurement shows no gain.
- A line budget for the delegate references — 2 edits, 1 gate run. The prose budget table owns it.

## Further notes

Measurements on 2026-09-11, each a delegate or exec child that replied "ok" with no tool use:

| subject | tokens |
|---|---|
| Claude general-purpose delegate | 38,430 |
| Claude Explore delegate | 8,834 |
| Codex exec, default config | 7,258 |
| Codex exec, docs MCP server disabled | 7,258 |

Closed decisions, 2026-09-11: the Codex side closes as a measured no-op. Every phase stays model-invocable. The listing item trims the descriptions and adds a length check. The budget of 250 characters is a proposal for reviewer veto; the trim lands under it.

Flagged additions beyond the decision source: DB10 and DB11 set the thresholds at 15,000 and 16,000 tokens, where the source said "a third". DB17 reports a missing description, and DB18 reports a folded one. DB25 counts runes, DB30 pins the limit to the table, and DB5 requires the read tool and the shell tool.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| A custom agent with a tools list cuts the boot by a third. | DB1, DB2, DB4, DB5, DB10, DB11 |
| Never use fork for a cheap delegate. | DB13 |
| A custom agent that binds its own model may arrive with no model in the call. | DB3, DB12 |
| The consumer payload must ship a new agents tree. | DB6, DB7, DB27 |
| The delegate skill must name the new agent types as the required surface. | DB8, DB9, DB28 |
| Trim the descriptions and add a length check. | DB14 to DB22, DB24, DB25, DB29, DB30 |
| Those will need to be invoked by the model at some point. | DB23 |
| Codex exec boots at 7,258 tokens with and without the docs MCP server. | DB26 |

Pre-review proof checklist:

- `Cited symbols`: `AgentLineVerdict`, `buildLinkPlan`, `linkDestination`, `RequiredPackAssets`, `profileSection`, `markdownRow`, `bounds.ClassifyNoFollow`, and `skillsindex.FrontmatterField` resolve in the tree.
- `Import edges`: none.
- `Source-row clauses and occurrences`: the source is this conversation, quoted in the table above, each sentence with one occurrence.
- `Promised field labels`: `name`, `description`, `tools`, and `Skill description budgets`.
- `Changed-function callers`: no production function changes, and the subagent-type test table gains two rows.
- `Copy survival`: none.
