# Learnings — usage journal

## 2026-08-14 — Read enforcement surfaces by content before speccing against them [open]

What happened: the spec-authoring-and-light-path verification took spec 5 + tickets 3 iterations; blockers in four separate rounds shared one root cause — coverage claims about the anchors registry, canary fixtures, and bespoke conformance tests written from file and fixture names instead of their content (a fixture whose EXPECT pins a different clause than its name implies, a "no live-tree sweep exists" claim the tree falsified, a hardcoded diagnostic-family census, a phantom registry needle). One further miss class was sequencing, not reading: two tickets mutually presupposed each other's not-yet-landed step.

Right behavior: before locking rows against gate-enforced artifacts, dump the enforcement surface exhaustively — every fixture EXPECT and MUTATE string, every bespoke check that greps or counts the target files, the body of any test cited as "already covered" — and derive claims from that dump.

Proposed rule change: landed already via the light path (bench-write-spec.md step 1, two commits, gate-green). A/B probe of the new sentence was inconclusive by instrument flaw: the probe questions themselves demanded content-level verification, so both arms verified; a real test needs authoring-shaped tasks with hidden grading. No further rule change proposed.

## 2026-08-14 — A/B rerun: enforcement-surface paragraph [open]

What happened: five treatment drafts scored 5, 5, 8, 5, and 6 (mean 5.8/10; range 5–8) versus control scores of 3, 5, 4, 7, and 4 (mean 4.6/10; range 3–7), a +1.2-point treatment difference. Files-read lists explicitly named an `EXPECT` or `MUTATE.json` in 2/5 treatment drafts versus 1/5 controls; all ten named `fixture_bite_test.go`. The direction weakly supports the 2026-08-14 learning "Read enforcement surfaces by content before speccing against them", but the evidence is not decisive: n=5 per arm, one repository, one authoring task family, and the same base model in both arms.

Right behavior: keep content-level enforcement verification as the spec-authoring expectation, while treating this rerun as narrow behavioral evidence rather than proof that the paragraph reliably changes every draft.

Proposed rule change: none. Preserve the current rule and let the reviewer decide at the next drain whether the weak directional result warrants any follow-up experiment.
