# Pair two independent readers on each research slice

Blocked by: none
Writes: .agents/skills/bench-craft-research/SKILL.md, projects/benchkit.md
Covers: none

## What to build

The `craft-research` skill sends one read-only delegate to each research slice. On the FT373 build on 2026-10-08, two fan-outs read the same three slices, one on haiku and one on sonnet. Each tier missed one fact that the other found. Haiku missed the mapsloop Go version gate, and sonnet missed the changelog prose exclusion. Each miss was a file that one reader opened and the other did not. Each miss later returned as a review finding at a higher cost.

Add one rule to the `## Rounds and fan-out` section of the skill, as one or two paragraphs. A research slice goes to two independent readers at the same time. The `craft-line` skill owns the tier of each reader. The two readers get the same charge and do not see each other's return.

The coordinator diffs the two returns by their `path:line` citations. The agreed set is the trusted set. Each disagreement and the union of the two not-verified lists form the frontier. One targeted read closes the frontier before the coordinator writes the consuming charge. A sequential second read that sees the first return is not a pair, because it inherits the first reader's spelling and scope.

Add one rule to the `## The question graph` section, as one sentence or two. A question asks for the fact and for its nearest qualifier. Examples of a qualifier are the exemption list of a bound and the version gate of a rule. With the qualifier, a literal answer does not stop one file early.

Add the closure rule to the `## Sources and verification` section, as one short paragraph. For each fact a reader cites, the reader follows one hop. It opens the file that the fact names and reports the rule there that qualifies the fact. An edge the reader saw and did not follow goes on its not-verified list, so that list is the open set of the search. The coordinator dispatches the next read at that list.

Raise the skill's row in the guidance-prose budget table of `projects/benchkit.md` to the new line count. Write the prose in ASD-STE100: sentences of 25 words or fewer, active voice.

## Acceptance

- [ ] The `## Rounds and fan-out` section states the paired-reader rule. The rule names two independent readers, one charge, a citation diff, the agreed set, and the frontier closed before the consuming charge.
- [ ] The `## The question graph` section states the fact-plus-qualifier rule.
- [ ] The `## Sources and verification` section states the one-hop closure rule and names the not-verified list as the open set.
- [ ] `bench gate-prose . -- .agents/skills/bench-craft-research/SKILL.md` passes, and `bench test --check guidance-prose-budgets` passes with the raised row.
