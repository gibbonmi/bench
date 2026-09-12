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
