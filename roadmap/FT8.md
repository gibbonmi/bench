**FT8 (overdue, decision required) — Sonnet 5 mid-tier revisit: a risk-based escalation ladder for the default implementer.**

The 2026-08-30 decision is overdue and awaits the reviewer's verdict.

The proposal: make Sonnet 5 the default implementer, with an escalation
ladder for risk. A ticket carrying an unsound-charge risk (a new seam, real
logic, dedupe) routes to Opus up front. One delegate-attributed catch during
review re-charges that ticket at Opus, low effort. Two catches in one build
move the rest of that build to Opus, medium effort.

Evidence so far: Sonnet 5 ran 53/53 green at diff level on one comparable Go
refactor. Sonnet 5 obeys a wrong charge where Opus corrects it. The ladder's
risk signals must watch spec and ticket attribution, not just catch counts.
Next: decide
Occurrence: 2026-08-23 capture — parked the escalation-ladder proposal after one comparable Sonnet 5 run.
