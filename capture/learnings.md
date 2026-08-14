# Learnings — usage journal

## 2026-08-14 — Read enforcement surfaces by content before speccing against them [open]

What happened: the spec-authoring-and-light-path verification took spec 5 + tickets 3 iterations; blockers in four separate rounds shared one root cause — coverage claims about the anchors registry, canary fixtures, and bespoke conformance tests written from file and fixture names instead of their content (a fixture whose EXPECT pins a different clause than its name implies, a "no live-tree sweep exists" claim the tree falsified, a hardcoded diagnostic-family census, a phantom registry needle). One further miss class was sequencing, not reading: two tickets mutually presupposed each other's not-yet-landed step.

Right behavior: before locking rows against gate-enforced artifacts, dump the enforcement surface exhaustively — every fixture EXPECT and MUTATE string, every bespoke check that greps or counts the target files, the body of any test cited as "already covered" — and derive claims from that dump.

Proposed rule change: landed already via the light path (bench-write-spec.md step 1, two commits, gate-green). A/B probe of the new sentence was inconclusive by instrument flaw: the probe questions themselves demanded content-level verification, so both arms verified; a real test needs authoring-shaped tasks with hidden grading. No further rule change proposed.
