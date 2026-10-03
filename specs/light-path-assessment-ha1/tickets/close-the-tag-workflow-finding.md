# Close the tag workflow finding

Blocked by: none
Writes: ASSESSMENT.md, specs/light-path-assessment-ha1/tickets/close-the-tag-workflow-finding.md
Covers: none

## What to build

`ASSESSMENT.md` keeps finding H-A1 open, but its content is fixed. The tag
workflow compiles the publisher from the tag checkout and runs
`bench release submit`. The workflow also keeps the publication record when
the publish step fails. The reviewer closed H-A1.

The assessment changes each claim that depends on H-A1:

- The Decision paragraph names the conditions that still block release.
- The inventory counts 0 high findings.
- Row 6 of the reconciliation table states its true grade and its current
  workflow line ranges.
- The adoption section records the closed finding in the same form as the
  closed planted-reason finding.
- The ranked backlog drops the CI publication routing row.

Deployment stays NO-GO. The Decision paragraph agrees with the
release-readiness status and the reassessment gate in `ROADMAP.md`. This
ticket does not edit `ROADMAP.md` or `roadmap/`.

## Acceptance

- [ ] `ASSESSMENT.md` has no open H-A1 finding and records the publication lifecycle as closed.
- [ ] The `ASSESSMENT.md` inventory counts 0 high findings.
- [ ] Row 6 of the reconciliation table states a grade that matches its evidence.
- [ ] No `ASSESSMENT.md` backlog line lists CI publication routing as open work.
- [ ] The `ASSESSMENT.md` Decision paragraph keeps NO-GO and names the unmet release conditions.
