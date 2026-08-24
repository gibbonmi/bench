**FT8 (scheduled, decision required) — Sonnet 5 mid-tier revisit: a risk-based escalation ladder for the default implementer.**

Time-boxed to 2026-08-30, after `module-size-split` runs with Sonnet 5 as the
implementer — the second comparable Go refactor run. This brings the original
2026-09-01 trigger forward.

The proposal: make Sonnet 5 the default implementer, with an escalation
ladder for risk. A ticket carrying an unsound-charge risk (a new seam, real
logic, dedupe) routes to Opus up front. One delegate-attributed catch during
review re-charges that ticket at Opus, low effort. Two catches in one build
move the rest of that build to Opus, medium effort.

Evidence so far: Sonnet 5 ran 53/53 green at diff level on one comparable Go
refactor. Sonnet 5 obeys a wrong charge where Opus corrects it. The ladder's
risk signals must watch spec and ticket attribution, not just catch counts.
Next: decide
Occurrence: 2026-08-23 capture — parked the escalation-ladder proposal, with
one comparable Sonnet 5 run as evidence and the review date brought to
2026-08-30.
