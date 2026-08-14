# Require reading enforcement-surface content before locking rows

Blocked by: none
Writes: .agents/commands/bench-write-spec.md

## What to build

`/bench-write-spec` step 1 gains one short paragraph: when the spec will edit
gate-enforced prose — anchored clauses, canary fixtures, counted or
byte-pinned substrings — the author reads the enforcement surface's *content*
before locking rows (every fixture `EXPECT` and mutation string, every
bespoke check that greps or counts the target files), because a claim built
from file and fixture names alone ships wrong fixtures. Sourced from the
spec-authoring-and-light-path verification loops, where this single miss
class produced blockers in four separate rounds. No anchor or fixture: the
sentence is guidance, not a workflow contract clause; anchoring is a cheap
later addition if it proves load-bearing. The edit stays clear of the
sections the staged spec-authoring-and-light-path tickets rewrite (step 9,
Exit handoff, "Who runs this phase", Template).

## Acceptance

- [ ] bench-write-spec.md step 1 names EXPECT files, mutation strings, and
      bespoke grep/count checks as required pre-row reading for specs that
      edit gate-enforced prose
- [ ] the file's anchored needles and its two `Bootstrap authority before
      execution` occurrences are unchanged
