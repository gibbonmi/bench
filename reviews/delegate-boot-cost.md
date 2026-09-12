# Delegate boot cost review

## Chunk C1

The frozen pair is base `347b5cfff016e2992a43c9d8045ee7c6abf27712` and tip `af405b7d5a16912aacd5c632a9495648c5a3337b`.
Three sonnet axes ran on 2026-09-11 at high effort, each in its own native context, and each read the integration source read-only.
The reviewer selected the cheap tier for the axes at phase entry.
Raw findings: Standards 3, Spec 2, Coverage 5.
De-duplicated repair targets: 7, and every one is auto-fix.

Author verification before the review: the C1 tests ran green, and the completion plan's C1 probe bit.
The probe omitted the model-declared diagnostic, and `model-declared` stopped biting through its owner.

## Standards

Findings: 3. Worst issue: medium.

- ST1 (medium, auto-fix): the `bench-` naming convention was derived twice, once as the file filter's constant and once as a literal inside the skill-name pattern. A change to the convention would reach one site only. The repair builds the pattern from the constant. Citations at tip af405b7d: internal/conformance/claude_agent_definitions_test.go:24, :43.
- ST2 (low, auto-fix): the check's doc comment ran 30 words, over the 25-word ASD-STE100 bound, and joined two facts with "and". The repair splits it into three sentences. Citations at tip af405b7d: internal/conformance/claude_agent_definitions_test.go:45-47.
- ST3 (low, auto-fix): the per-file helper's doc comment ran 27 words. The repair splits it. Citations at tip af405b7d: internal/conformance/claude_agent_definitions_test.go:66-68.

## Spec

Findings: 2. Worst issue: high.

- SPEC-1 (high, auto-fix): the routing rule assigned a fan-out search to `bench-reviewer`. The spec's edge inventory reserves that charge for the built-in Explore agent as a Won't handle. The routing decision names only a review axis and a diagnostic consultation. The repair drops the fan-out clause from three prose sites. Citations at tip af405b7d: .agents/skills/bench-craft-delegate/SKILL.md:24; specs/delegate-boot-cost/spec.md:203.
- SPEC-2 (medium, auto-fix): the README and the agent description repeated the same fan-out claim. The contradiction therefore reached the file a cold reader trusts. SPEC-1's repair closes it. Citations at tip af405b7d: .claude/README.md:22; .claude/agents/bench-reviewer.md:3.

The axis confirmed DB1 to DB9, DB12, DB13, DB27, and DB28 as implemented, and recorded DB10 and DB11 as an open shortfall.
No path in the delta falls outside the ownership fences.

## Coverage

Findings: 5. Worst issue: high.

- COV-1 (high, auto-fix): the forbidden and required tool lists were each proved by one member, so `Artifact` and `AskUserQuestion` were advertised and ungraded. A trim of the forbidden list to `{"Agent"}` compiles, is behavioral, and left every C1 test green. The author's first repair ranged over the production slice and stayed silent under the same trim. The second repair names both lists independently in the test and asserts equality, and the trim then reds. AGENTS.md permits that independent expectation, because the named mutation's red is what the independence buys.
  Citations at tip af405b7d: internal/conformance/claude_agent_definitions_test.go:34, :38.
- COV-2 (medium, auto-fix): the per-file refusal branch reached no test. The only symlink test graded the directory branch instead. A link at an agent file skips every other diagnostic for that file. Citations at tip af405b7d: internal/conformance/claude_agent_definitions_test.go:76-77.
- COV-3 (medium, auto-fix): the "agents path exists as a plain file" branch reached no test. An unrefused one enumerates nothing and reports a clean adapter. Citations at tip af405b7d: internal/conformance/claude_agent_definitions_test.go:123-124.
- COV-4 (medium, auto-fix): the C1 row's tests column named the two checks and three packages, but not `go test ./internal/conformance`. The named check runs the root conformance test alone, so the row's own recipe never reaches the fixture-bite seam that seven of its rows cite. The repair adds the package to the column. Citations at tip af405b7d: specs/delegate-boot-cost/spec.md:131.
- COV-5 (low, no-op as a separate target): the DB4 and DB5 row texts promise more tool names than their canary fixtures carry. COV-1's repair covers both. Citations at tip af405b7d: specs/delegate-boot-cost/spec.md:172, :173.

The axis ran the four C1 packages one at a time and reported every one green.

## Chunk C2

The frozen pair is base `458f72e067dd5ad4695dfd0ea007503d95fe0fac` and tip `0f4dd32fb79828d294aeabd5afd5bfb91d8a9623`.
Three sonnet axes ran on 2026-09-12 at high effort, each in its own native context, and each read the integration source read-only.
Raw findings: Standards 5, Spec 7, Coverage 8.
De-duplicated repair targets: 13, and every one is auto-fix.

The reviewer directed delegated authorship for this chunk's two tickets.
Both tickets ran on the mid tier at low effort, and the coordinator verified each done-claim against the tree.
Author verification before the review: the packages ran green uncached, and the author's probe bit.
The probe omitted the white-space collapse before the rune count.
The coordinator's own probe swapped the glob-row limit for a constant. That site and that kind differ from the author's, and it bit on two tests.

## Standards at C2

Findings: 5. Worst issue: high.

- ST4 (high, auto-fix): the new policy parser copied its sibling but dropped the duplicate-subject guard. Two profile rows naming one subject therefore let the later row win with no diagnostic. The repair restores the guard and its own fail-closed case. Citations at tip 0f4dd32f: internal/conformance/skill_description_budgets_test.go:115-165; internal/conformance/prose_budget_test.go:131.
- ST5 (medium, auto-fix): the fixture table renderer was byte-identical to the prose budget's but for one constant. AGENTS.md names a pasted fixture harness as the defect. One `profileBudgetTable` helper now serves both. Citations at tip 0f4dd32f: internal/conformance/skill_description_budgets_test.go:239-245.
- ST6 (medium, auto-fix): the subject enumerator and its refusal ladder were a second implementation of the prose budget's own. The repair makes the generalized `listingTreeEntries` the one enumerator, and the older check now calls it. Every prose-budget diagnostic string stays byte-identical. Citations at tip 0f4dd32f: internal/conformance/skill_description_budgets_test.go:176-236.
- ST7 (low, no-op): five canary fixtures ship an identical profile baseline that the canary package's `@` include could share. The fixtures stay independent, because each one's mutation anchors in its own copy. Citations at tip 0f4dd32f: tests/canary/skill-description-budgets/.
- ST8 (low, auto-fix): two new profile sentences used the passive voice with no agent. The repair names the check as the subject. Citations at tip 0f4dd32f: projects/benchkit.md:518, :525.

## Spec at C2

Findings: 7. Worst issue: high.

- SPEC3 (high, auto-fix): DB22 was false in six of the twenty trimmed descriptions. Each dropped a phrase its own body repeats. `craft-research` lost a clause its body states almost verbatim, and `craft-tickets` lost its own heading. `craft-synthesis` lost the names of its three loops, and `craft-spec` lost the term Won't handle. `craft-domain` lost its invocation contexts, `craft-line` lost its escalation trigger, and the repair restores every phrase within 250 runes.
  Citations at tip 0f4dd32f: .agents/skills/bench-craft-research/SKILL.md:3; .agents/skills/bench-craft-tickets/SKILL.md:3; .agents/skills/bench-craft-synthesis/SKILL.md:3; .agents/skills/bench-craft-spec/SKILL.md:3; .agents/skills/bench-craft-domain/SKILL.md:3; .agents/skills/bench-craft-line/SKILL.md:3.
- SPEC4 (low, auto-fix): the DB25 row cited a test name the tree does not carry. The repair points the row at `TestSkillDescriptionBudgetCountsCollapsedRunes`. Citations at tip 0f4dd32f: specs/delegate-boot-cost/spec.md:197.

The axis confirmed every other C2 row as implemented, and it verified DB26's two equal Codex counts.
No path in the delta falls outside the ownership fences, and no Won't handle edge gained an implementation.

## Coverage at C2

Findings: 8. Worst issue: high.

- COV6 (high, auto-fix): seven diagnostics shipped with no red-capable evidence. The axis named the glob-validity guard as a mutation that leaves the whole suite green. A typo'd glob would then enforce a budget on nothing instead of failing closed. The repair adds one test per diagnostic, and the guard's omission now reds.
  Citations at tip 0f4dd32f: internal/conformance/skill_description_budgets_test.go:57, :145, :149-151, :205, :225, :227, :231.
- COV7 (low, no-op): the canary fixtures grade synthetic limits and not the shipped 250. The live-tree assertion covers the real numbers on the clean side, and DB21 names that seam. Citations at tip 0f4dd32f: internal/conformance/skill_description_budgets_test.go:41-75.
- COV8 (low, no-op): a future fixture that reuses a real subject path would have its restore backed by live content. Every current fixture path is fictitious. The reviewer keeps this as a note for the next fixture author. Citations at tip 0f4dd32f: internal/canary/mutation.go:132-135.

After the repair the coordinator probed the newly shared enumerator, a site neither earlier probe touched.
The swap turned three tests red, one of them the older prose budget's own, so the collapse did not weaken the check it absorbed.

## Final reconciliation

The coverage map holds 30 rows, and 28 of them close at the landing tip.
Every gate-reachable row grades through its named seam, and both new checks pass over the live tree through the built binary.

Two rows stay open: DB10 and DB11, the observed boot cost of each agent type.
Claude Code reads the adapter's agents directory when a session starts, so the implementation session cannot spawn a type its own start never saw.
The coordinator proved this: a spawn of `bench-reviewer` returned "Agent type not found".
The landing is therefore the step that makes the measurement possible.

The reviewer decided on 2026-09-12 to land with these two rows open.
The next session takes both numbers first, with one "ok" spawn per type and no tool use, and records each usage line in the retro.
If a number misses its threshold, the agent file's tool list is the repair surface, and the spec does not reopen.
