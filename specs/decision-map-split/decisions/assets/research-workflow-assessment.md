# Research workflow assessment

## Recommendation

Use the existing Research map as one planning index for the research portions
of FT99, FT106, FT125, FT304, and FT231. Keep each roadmap row as the owner of
its wider scope. Group report quality and evidence work first, freshness and
focused readers second, and shared observation last. Grouping alone selects no
spec composition. The reviewer selects that composition in map ticket #14.

## Evidence status

This assessment examined four local incident-investigation documents. They are
the [Duo execution report](</home/mgibs/workspace/gitlab-runner-observabiltiy/incident-investigation/docs/research/gitlab-duo-investigation-execution.md:3>),
the [Orbit capability report](</home/mgibs/workspace/gitlab-runner-observabiltiy/incident-investigation/docs/research/gitlab-orbit-capabilities.md:3>),
the [Orbit design report](</home/mgibs/workspace/gitlab-runner-observabiltiy/incident-investigation/docs/research/orbit-incident-investigation-design.md:3>),
and the [incident roadmap](</home/mgibs/workspace/gitlab-runner-observabiltiy/incident-investigation/docs/incident-investigation-roadmap.md:115>).

The documents are untracked at the reported source HEAD
`a50be7d47fadfc9aecde65dd6372ca6b9dda464a`. This assessment has no immutable
content snapshot for them. Their citations support report structure only. They
do not verify vendor facts, tenant behavior, or agent execution. No transcript
was available to establish that an agent ran.

## Reference skill chain

The [Research skill](</home/mgibs/workspace/gitlab-runner-observabiltiy/.agents/skills/research/SKILL.md:6>)
has three substantive rules. It delegates factual reading, uses primary
sources, and leaves one cited Markdown report. The
[Wayfinder skill](</home/mgibs/workspace/gitlab-runner-observabiltiy/.agents/skills/wayfinder/SKILL.md:11>)
adds a planning boundary. Its map is an index, its tickets hold decisions, and
its frontier names unblocked questions.

The [Grilling skill](</home/mgibs/workspace/gitlab-runner-observabiltiy/.agents/skills/grilling/SKILL.md:6>)
uses a dependency tree and asks the settled frontier. The
[domain-modeling skill](</home/mgibs/workspace/gitlab-runner-observabiltiy/.agents/skills/domain-modeling/SKILL.md:48>)
sharpens terms and tests boundaries against code. The
[Prototype skill](</home/mgibs/workspace/gitlab-runner-observabiltiy/.agents/skills/prototype/SKILL.md:8>)
answers one question with disposable code.

## Report contract

Put the recommendation first. Then state scope, evidence status, and the
decision the report can support. Keep evidence, inference, tested results, and
proposals in clear separate sections. Do not require a label on every sentence.

| element | content | consequence |
|---|---|---|
| comparison | capability or option table with a consequence | The reader can compare choices. |
| evidence | nearby primary citation for each material claim | The reader can re-open support. |
| uncertainty | contradictions and unknowns | The report does not invent closure. |
| relation | diagram where prose hides a material flow | The reader can inspect the join. |
| next evidence | validation plan and observable result | The reviewer can choose the next check. |

The Duo report demonstrates a recommendation, separate interface claims,
unknowns, and a validation plan
([report](</home/mgibs/workspace/gitlab-runner-observabiltiy/incident-investigation/docs/research/gitlab-duo-investigation-execution.md:5>)).
The Orbit report demonstrates a recommendation-first capability table and
separates documented facts from proposals
([report](</home/mgibs/workspace/gitlab-runner-observabiltiy/incident-investigation/docs/research/gitlab-orbit-capabilities.md:5>)).
The design report demonstrates a bounded flow diagram and named outcome branches
([report](</home/mgibs/workspace/gitlab-runner-observabiltiy/incident-investigation/docs/research/orbit-incident-investigation-design.md:11>)).

## Bench fit

The settled [Research map](../craft-research.md:25) assigns factual reading to
Craft research. It assigns delegation and line mechanics to their current
owners. A coordinator verifies material claims and the joins between returns.
It preserves conflicts and unknowns
([fan-out](../craft-research.md:116), [artifact](../craft-research.md:152)).

The coordinator does not use automatic recursive duplicate delegation. A new
round needs a named conflict, unknown, or coverage gap. Empirical probes remain
separate Prototype prerequisites. Research can identify a probe, but cannot
close its claim before the probe returns
([fan-out](../craft-research.md:132), [boundary](../craft-research.md:203)).

The current parser keeps answers inline and supports `shaping` and `ready`
states ([schema](../../../../internal/maps/schema.go:116)). `bench maps` projects
unresolved ticket rows ([projection](../../../../internal/maps/maps.go:44)). Status
shows unresolved maps before ready maps
([status](../../../../internal/status/status.go:950)). This supports focused discovery
without treating default-output omission as data loss or adding another state
owner.

## Roadmap composition

[FT99](../../../../roadmap/FT99.md:3) contributes current-tree premise checks.
[FT106](../../../../roadmap/FT106.md:58) contributes document freshness and the
artifact-home conflict. [FT125](../../../../roadmap/FT125.md:3) contributes focused
readers. [FT304](../../../../roadmap/FT304.md:3) contributes shared observation.
[FT231](../../../../roadmap/FT231.md:73) contributes advisory measurement and permits
a minimal baseline before the first approved pass.

The map groups those research portions only. FT89, FT219, and FT292 remain
related roadmap rows. They require a fresh decision before they join this map.
The grouping does not close any row. It authorizes no runtime build, new CLI
syntax, schema migration, or monolithic specification.

## Open choices

Inline storage remains the current default. The reviewer can retain it, or
reopen a migration only when a focused-reader benefit justifies the cost. The
original schema exclusion remains in force until that explicit reviewer action.

Ticket #13 narrows evidence work to FT106's artifact-home and trackability
conflict plus a mechanical freshness projection. It does not reopen the settled
coordinator, Sources, drift, or retirement rules. Ticket #14 selects the future
specification cuts and related measurement.

## Validation plan

1. Re-read the map, roadmap, parser, and status owners before a specification.
2. Compare a report with the contract's evidence and inference separation.
3. Exercise any focused reader against unresolved and ready maps.
4. Measure the approved change with FT231's advisory posture.
5. Ask the reviewer to resolve tickets #11 through #14 before implementation.

## Sources

- [Settled Research map](../craft-research.md:25) owns research policy.
- [FT99](../../../../roadmap/FT99.md:3), [FT106](../../../../roadmap/FT106.md:58),
  [FT125](../../../../roadmap/FT125.md:3), [FT304](../../../../roadmap/FT304.md:3), and
  [FT231](../../../../roadmap/FT231.md:73) own the grouped roadmap scope.
- The four absolute incident-investigation links above are report-structure
  references with the Evidence status limitation.
