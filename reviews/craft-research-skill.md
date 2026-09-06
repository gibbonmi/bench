# Review pickup: craft-research-skill

Frozen base `3f29857f7c8a1e52807f48a16c27179586fe305c`, reviewed tip `2dcb445eee9e12a93de762910f8db30dc1182aa9`. Raw findings: 8. Repair targets: 5, all repaired in the repair commit. One item stays open for the reviewer.

## Standards

Count 0 open. S1 and S2 are repaired.

## Spec

Count 0 open. P1 and Codex F3 were no-ops, and Codex F1 and F2 are repaired.

## Coverage

Count 1 open. Worst: C1, a coverage row claims a catch the seam cannot make.

- C1 `ask-user`. Four of the five anchored phrases in `.agents/commands/bench-assess.md` also sit in the frontmatter `description:` line, and the anchor match is a whole-file substring. A body rewrite that drops a phrase stays green while the description keeps it, and the probe proved it. CR15's why-it-catches clause is therefore not true for a body-only change, and the condition predates this diff. Repair shape if accepted: `RequireInSection` rows for those four phrases, or body-unique needles. Not repaired here.
