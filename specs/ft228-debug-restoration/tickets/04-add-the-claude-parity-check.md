# Add the Claude parity check and complete the policy

Blocked by: 03-settle-the-codex-trigger.md
Writes: internal/conformance/skills_index_checks_test.go, tests/canary/skills-index-command-adapters/, .bench/BENCH-reference.md, CHANGELOG.md

## What to build

The Claude half: grade the command files' invocation frontmatter against the
table ticket 03 landed, close both table-completeness directions, and update
the docs.

- Extend the check: each command file's frontmatter `disable-model-invocation`
  key equals the negation of its row's Claude boolean, parsed from the
  frontmatter block only so a body mention of the key is inert. A command file
  with no table row is red, and a table row with no command file is red — the
  policy lookup runs ahead of the adapter-existence checks, so an undeclared
  phase reds as undeclared even when its adapter is also missing. The
  `bench-what-next` alias keeps its thin-alias body handling; its policy row
  (false/false) matches the key its own command file carries.
- Add two `skills-index-command-adapters` fixtures: a Claude-invocable command
  file gaining the disable key, and a ghost command file the table does not
  name (the fixture supplies only the ghost command file). `BASE` per the
  spec's fixture mechanics.
- Add in-package cases over constructed roots for the stale table row (a
  `BASE`-style root omitting one command file the table names) and the
  body-prose key mention staying green.
- Rewrite `.bench/BENCH-reference.md`'s Harness Invocation paragraph: per-phase
  policy graded by the gate on both surfaces, `$bench-debug` named as the
  implicitly invocable exception with its one-clause why. One CHANGELOG entry
  under Unreleased documents the invocation-policy settle.

## Acceptance

- [ ] The kit gate is green with the full two-sided grading over the real tree.
- [ ] The flipped-frontmatter fixture and the unlisted-phase fixture red through the registered owner naming the phase, and each restore re-runs green (IP3, IP4).
- [ ] The in-package cases red the stale table row and keep the body-prose mention green (IP5, IP8).
- [ ] `.bench/BENCH-reference.md` states the per-phase policy and the exception, the docs-currency sweep is green, and the CHANGELOG carries the settle entry.
