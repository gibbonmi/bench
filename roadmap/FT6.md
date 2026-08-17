**FT6 (LOW, parked pending evidence — leave parked):**
`bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check, and a per-anchor
bite-proof meta-test (canaries prove one needle per family today; graduate on
observed anchor rot). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search. Also parked here 2026-07-23:
concurrent `bench upgrade` runs, raised as a coverage gap by the FT85 review
and closed by decision rather than left open — `transactionalLink` already
moves tree, manifest rows, and version stamp together, so the damage is
bounded; graduate on an actual report of two upgrades interleaving badly, not
before. Two upstream candidates join the pending-evidence tier rather than
becoming new skills: a generated human-procedure wizard plus last-resort HITL
loop graduates on a Bench workflow that cannot be made agent-operable, and a
third-party questionnaire graduates when a real decision map blocks on someone
other than the reviewer and reviewer-directed grilling misroutes the question.
Source: `upstream(mattpocock/skills@84fdeff)`, drained from
`capture/IDEAS.md` here. Also parked here 2026-08-12: a one-off `bench commit`
refusal — `gate: red` / `prospective authorization refused: inherited` — on a
tree whose immediately following direct `bench gate` was green (the retried
commit landed on the fresh verdict); not reproduced across seven later
landings, and a repro through anything but `bench commit` itself proves
nothing. Workaround on recurrence: run `bench gate` directly and retry.
Graduate on a second reproduced refusal through `bench commit`; the
verdict-class plumbing around `inherited` records is the suspect. Source:
`capture/learnings.md` 2026-08-12, verdicted here. Also parked here
2026-08-16: moving `.bench/BENCH.md`'s light-path threshold (map #16 of the
retired `spec-ticket-fence-reduction` spec) — its own decision #16 deferred the
move pending a measurement of whether the reduced-schema heavy path sees more
uptake, and that observation window has not run yet. Graduate on either the
measurement landing or a reviewer call to move without it; the direction, the
replacement Observable wording, and the anchors/canaries it drags are
unresolved. Source: `capture/IDEAS.md`, drained here.
