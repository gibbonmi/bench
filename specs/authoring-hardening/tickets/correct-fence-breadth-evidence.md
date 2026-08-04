# Correct the fence-breadth paragraph's false evidence clause

Blocked by: teach-size-split-signal.md
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`
Contracts: none crosses

## What to build

The landed fence-breadth paragraph closes on a claim the evidence contradicts:
that wide fences are where repair rounds have clustered "while narrow ones
landed first time". Measured across the recovery-discard ticket set, the wide
end holds — both four-directory tickets drew repairs — and the narrow end
fails. A one-directory fence drew two repairs, and those two were the most
severe defects in the build: an enumeration trusting an unvalidated persisted
ref, and a residue scan that read only the live record. Two further tickets at
breadth one and two also drew repairs.

Keep the rule, the threshold, the counting convention, and the justification
price. Replace only the evidence clause with one that survives the measurement:
breadth carries some risk signal at the wide end and carries none at the narrow
end, so a narrow fence is not evidence that a ticket is sound. No sentence in
the paragraph may leave a reader believing that staying under the threshold
makes a ticket safe — that belief is what the measured counterexample punishes.

## Acceptance

- [ ] [FB1] the paragraph no longer claims narrow fences landed first time, and no sentence in it implies that staying under the threshold makes a ticket sound.
- [ ] [FB2] the split-signal rule, its more-than-two-directories threshold, the file-entry-counts-its-parent counting convention, and the one-line justification price survive the correction unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FB1 | restore the landed-first-time clause | the reviewer, at review | reinstate the original evidence clause, re-read the paragraph against the measured recovery-discard breadth distribution, expect it to assert a correlation the one-directory counterexample falsifies |
| FB2 | drop the threshold along with the evidence clause | the reviewer, at review | remove the more-than-two-directories trigger while correcting the evidence, re-read against `specs/authoring-hardening/spec.md` story 2, expect the split signal to survive with nothing left to fire it |
