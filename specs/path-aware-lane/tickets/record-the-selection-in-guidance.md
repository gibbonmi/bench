# Record the selection in the guidance

Blocked by: select-kit-lane-checks.md
Writes: CONTEXT.md, docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md, CHANGELOG.md

## What to build

A cold reader learns the selection from the guidance and not from the code. The
glossary names the two new terms, the ADR reads as the current decision, and the
release notes carry the behavior change.

`CONTEXT.md` gains two core terms. A **composed change** is one entry of the
diff between the expected base tree and the composed tree. It carries a status
letter, a source mode, a destination mode, and a path. A **path class** is one
row of the kit's table that binds a path predicate to check names. A path that
matches no class is `unknown`, and `unknown` selects every declared check. Each
entry names the term it replaces, in the shape the file's other entries use.

The existing `fast lane` entry drops the four fixed checks and states the
selection instead. The kit's lane is a declared check list, and the composed
changes select the checks that run. A lane pass still authorizes the worktree
commit alone.

ADR 0017 records the current decided state. The lane is a declared list, and the
composed changes select from that list. An unknown path selects every declared
check, so the selection never narrows below today's lane. The landing's
whole-project gate stays the one full grade. The ADR carries no file path and no
code, and it records the decision rather than its history. The stale sentence
that calls the lane "never a list derived from the diff" goes away.

`CHANGELOG.md` gains one `### Changed` entry under `## [Unreleased]`. The entry
names the new lane line and the classes cell, in the shape the section's other
entries use.

The review round grades this prose through the four review-owned rows below.

## Acceptance

- [ ] PL43: the `CONTEXT.md` core terms list holds one `composed change` entry and one `path class` entry, and the `path class` entry states that `unknown` selects every declared check.
- [ ] PL44: the `CONTEXT.md` `fast lane` entry names no fixed four-check list and states that the composed changes select the checks.
- [ ] PL45: ADR 0017 states that the composed changes select declared checks, drops the diff-denial sentence, and carries no file path and no code.
- [ ] PL46: `CHANGELOG.md` holds one `### Changed` entry under `## [Unreleased]` that names the `classes=` cell of the lane line.
- [ ] the gate `prose-mechanics` and `docs-currency-workflow` checks stay green over the whole tree.
