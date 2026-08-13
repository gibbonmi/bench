# Join bench roadmap to the approved AXI query set

Blocked by: add-context-envelope.md, add-row-selector.md, rebuild-bare-board.md, land-index-doctrine.md
Writes: cmd/bench/main.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, .agents/skills/bench-craft-cli/SKILL.md, .claude/skills/bench-craft-cli/SKILL.md, projects/benchkit.md, internal/roadmap

## What to build

`roadmap` becomes an approved member: the registry entry flips to
`axiApprovedRoot`, one `axiEnvelopeCase` binds the member
(`successMarker: "roadmap["`, empty setup = absent `ROADMAP.md`), the
independently authored `approvedAXIQueries` expectation gains
`"roadmap": nil`, and one craft-cli guidance table row plus one
`projects/benchkit.md` seam command satisfy the equality checks. The
disclosure contract is wording, not membership: index-mode success discloses
the exact row-selector command (and `--full`), the craft-cli disclosure cell
carries the row-selector command, and that phrase joins `checkAXIGuidance`'s
required-phrase list with a `TestAXIGuidanceContractBites` mutation case —
the equality check alone would pass a vague cell. This is the terminal
ticket: it blocks on land-index-doctrine.md because index mode must already
omit bodies for the disclosure row to mean what PI19 says — a disclosure
naming "request the complete value" over a surface that still returns
complete values grades nothing. The disclosure-cell sentence is guidance
prose: a collapsed fable-line edit flagged per craft-line, ridden here
because the equality check makes the row inseparable from the flip.

## Acceptance

- [ ] `bench roadmap` behaves as an approved member on all six graded
      behaviors: structured TOON success, definitive zero-row empty,
      structured stdout refusal outside a repo at exit 1, usage at exit 2 on
      an unknown flag, usage at exit 0 on each help spelling, and identical
      output from a deep cwd — with guidance/registry/profile membership in
      agreement (covers PI15).
- [ ] Index-mode success `help` names the row-selector command; the
      craft-cli cell phrase is pinned by the required-phrase list and its
      bite mutation fails when the phrase is removed (covers PI19).
