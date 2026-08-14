# Pin the cross-harness reviewer recipes

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md, internal/anchors/registry_data.go, internal/conformance/fixture_bite_test.go, tests/canary/workflow-guidance-anchors/

## What to build

Create `.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md`
with the exact commands, verified against the installed CLIs:
`claude -p --model <id> --effort <level> "<charge>"` and
`codex exec --sandbox read-only -m <id> -c model_reasoning_effort=<level>
"<charge>"`, plus the fallback rule: a harness with no native subagent surface
uses its own family's CLI. craft-delegate's SKILL.md gains the own-family →
native-agent-surface rule (never that family's CLI) and the pointer to the
reference, displacing equivalent lines — the file sits at its 120-line
budget, and its six cadence-pinned substrings stay byte-identical (else the
byte-exact strings in `internal/conformance/fixture_bite_test.go` update in
this ticket).
New Require anchors pin both recipe commands, the fallback, the rule, and the
pointer; new fixtures bite both halves. No file-missing EXPECT is used — the
harness itself cannot restore a deleted live file, which is the existing
protection. The skill-dir symlink already exposes the reference to Claude
Code; no adapter work. Registry and fixture paths land serially across the
whole spec.

## Acceptance

- [ ] the reference file carries both exact commands and the fallback,
      craft-delegate states own-family→native and points at it, and the new
      fixtures bite both halves (covers WF9)
- [ ] with the rule and pointer landed, craft-delegate is ≤ 120 lines
      (covers WF13)
