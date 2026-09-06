# Add the omission-proof and determinism rules to the spec skill

Blocked by: none
Writes: .agents/skills/bench-craft-spec/SKILL.md
Covers: none

## What to build

The bench-probe spec's first review iteration returned six blocking findings.
Three of them share one cause: the author locked rows before three checks.
`craft-spec` gains the three checks as rules, placed under `## The acceptance
coverage map` or under the `Explore the repo` reads, whichever section the
existing reads sit in. Write each rule as one ASD-STE100 sentence with its
one-clause why.

1. Before the first review charge, walk each restore or copy promise for a
   deterministic omission case. A restore comparison without an omission row
   is green by construction.
2. Cite no test-only helper across a package boundary, because a row that
   names another package's test helper has no seam in the graded package.
3. Use canned events for any row that compares two rendered reports, because
   two real runs differ in their elapsed times.

Keep every anchored sentence the registry in
`internal/anchors/registry_data.go` names for this file byte-exact.

## Acceptance

- [ ] `.agents/skills/bench-craft-spec/SKILL.md` states the deterministic omission walk for each restore or copy promise.
- [ ] The skill states that a row cites no test-only helper across a package boundary.
- [ ] The skill states that a row which compares two rendered reports uses canned events.
- [ ] `bench gate-prose . -- .agents/skills/bench-craft-spec/SKILL.md` exits 0.
- [ ] `bench test --check guidance-prose-budgets` and `bench test --check kit-compliance` exit 0 over the worktree.
