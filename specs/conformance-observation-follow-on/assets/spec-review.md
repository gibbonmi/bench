# Observation specs review record

Status: staged; full spec and ticket graph independently accepted in that order.

## Frozen review source

Repository: /home/mgibs/workspace/bench.
Assignment: maintenance-observations-spec.
Pre-repair commit: 5ce22734fc08ad7e9d0b05e0ca638e8f2ad46def.
The complete pre-repair source remains retained at this immutable commit.

| Spec | Reviewed SHA-256 | Rows |
| --- | --- | --- |
| conformance-observation-follow-on/spec.md | 64c89b84e9df62c23e647a0d554339f266315d75c2e97ae6fb215dea453e7763 | 38 |
| independent-expectation-audit/spec.md | c33b51d1b15082f72310a824046dbaa70e15853cf9ae43956d80c439bb5a9739 | 35 |

## Independent review

The coordinator reported the independent deepen_execution review on 2026-10-07.
Reviewer effort: xhigh.
The model identity was not supplied in that report.
This record attributes acceptance to that independent review, not to the author.

| Spec | Standards | Spec | Coverage | Reported confidence |
| --- | --- | --- | --- | --- |
| independent-expectation-audit | Accepted, 0 findings | Accepted, 0 findings | Accepted, 0 findings | 9 on each axis |
| conformance-observation-follow-on | Held, 1 material finding COF-R1 | Accepted, 0 findings | Accepted, 0 findings | 10 for COF-R1 |

The COF review accepted all 38 behavior rows, the reader bridge, and the strict-first oversized two-read/one-parse witness.
It also accepted the exclusions, complete LTE prerequisite, and other closure boundaries.
The independent expectation audit remains unchanged at its accepted hash.

## COF-R1 and bounded repair

The classifier fence omitted two required existing fixture units.
Both BASE lists name internal/bounds/classify.go on line 2.
Canonical FixturePins derives ownership from BASE dependencies even when the fixture mutates another consumer.
Preflight requires the owning ticket to co-name each inherited unit.

| Required fixture unit | Existing mutation and expectation purpose |
| --- | --- |
| tests/canary/package-core-guard/bounds-classify-limit-restated | Learnings replaces bounds.ControlRecordLimit with 2 << 20 in bounds.Classify, expecting the existing restated-limit diagnostic |
| tests/canary/package-core-guard/bounds-read-limit-restated | Models replaces bounds.ModelReadLimit with 5 << 20 in bounds.Read, expecting the existing restated-limit diagnostic |

Source witnesses: internal/canary/inventory.go:149–153,203–246 and internal/preflight/closure.go:54–73,118–123.
The author read both complete fixture units and these enforcement definitions before repair.
The coordinator reports that the full static census found only these two missing units.

Repair 1 adds both exact units to the COF ownership fence as read-only co-name pins.
It binds them to the classifier-owning COF-C3 slice when that slice writes classify.go.
BASE, MUTATE.json, EXPECT, CHECK, meanings, and existing assertions remain unchanged.
No acceptance predicate, reader policy, cache boundary, dependency, or exclusion changes.

## Learning and verification boundary

The author captured one focused learning through bench learning: Reader ports inherit BASE-pinned canary units.
It records the new reader port's inherited BASE closure, including the two mutations that target other consumers.
It does not duplicate the earlier general process-closure learning.

This repair runs prose validation, coverage validation, fixture digest comparisons, and the required planning commit lane.
It executes no mutation, test suite, manual gate, benchmark, implementation, or landing.
At the repair freeze, independent COF-R1 confirmation was pending.
The confirmation recorded below preceded ticket slicing.

## Acceptance before ticket slicing

The coordinator reported independent retained xhigh deepen_execution acceptance on 2026-10-07.
Independent-expectation-audit was accepted in the original full review, with all three axes at 0 and confidence 9.
COF was accepted after independent confirmation of COF-R1 repair 1, with all three axes at 0 and confidence 9.

The accepted source was clean commit 553a5ce608430f3d8053b78c58d83b99aa353025.
COF SHA-256: 13fe563743b2731b42dc6f29ba69a32599bfe52d0507d68372ecb893d389880f.
IEA SHA-256: c33b51d1b15082f72310a824046dbaa70e15853cf9ae43956d80c439bb5a9739.
The confirmation preserved all 38 COF row lines and eight classifier-fixture artifact files byte-identically; IEA remained unchanged.

## Ticket author and review boundary

Retained gpt-6.1-sol at high effort authored slicing pass 1, with at most two bounded slicing/repair iterations authorized.
Three serial vertical tickets retain COF-C1, COF-C2, and COF-C3 as independent review checkpoints.
The spec records the exact terminal-aggregate row moves; each earlier ticket repeats its first-use preservation and headroom obligations.
All 38 accepted row predicates remain byte-identical in the graph.
The complete accepted spec remains in Git source `553a5ce608430f3d8053b78c58d83b99aa353025` at its canonical spec path.

The source recheck found LTE still staged and sourcefiles absent.
Concrete API binding and landed-tree headroom therefore remain build-admission blockers.
No partial LTE outcome, private fallback, or TS dependency is accepted.
The planning pin retains 77 direct conformance files and the cited file lengths without production changes.

The canonical per-ticket closure proposal must run at the frozen graph commit with equal base and source-tip pins.
At the graph freeze, independent ticket review was pending; acceptance is recorded below.
No future operation count, mutation red, benchmark, or runtime outcome was executed or claimed here.

The authorized planning slice precedes LTE landing.
After complete accepted LTE lands, concrete API binding and the landed-tree headroom recheck must finish before any implementation charge.
This temporal clarification preserves COF01 and the complete-LTE admission boundary.

## Exact temporal clarification for graph review

The coordinator resolved the planning-timing contradiction on 2026-10-07 after reading COF01 and the affected source clauses.
Planning and slicing may use LTE's approved value contracts.
After complete accepted LTE lands, concrete API binding and the landed-tree headroom recheck must finish before any implementation charge.
This includes the first COF ticket; no partial provider or implementation exception is granted.

The complete predecessor graph remains immutable at 8c11acee86e8ffc322a2992128c7570d8117ee1a.
The amendment updates the COF spec, all three COF tickets, and this record.
It preserves all 38 predicates, fences, exclusions, operation counts, and the old-to-new row mapping.
IEA remains unchanged from that graph.
The subsequent independent ticket review accepted this exact admission boundary with the recorded aggregate mappings and first-use obligations.

## Independent ticket acceptance and stage close

The coordinator reported independent xhigh `deepen_execution` acceptance of both observation graphs on 2026-10-07.
The reviewed source was clean `c04ccd1bbac2274fb8197d592ea7aef21beffd94`.
COF spec SHA-256 was `0c3a7dd1616cb5b1caef7a62e51187da0074b3e242fcfe9475ef6e327bec8684`.
IEA spec SHA-256 was `a971cd21813f1f94efb1bee748ca4d846dfac295d5a9c886f881f2dd09104d29`.
Both graphs returned Standards 0, Spec 0, and Coverage 0, with confidence 9 on each axis.

The coordinator confirms the declared review line as `gpt-6.1-sol / xhigh`.
The effective backend identity and usage were not observed.

The review checked 38 COF and 35 IEA row lines and full fences against accepted source `553a5ce608430f3d8053b78c58d83b99aa353025`.
COF partitions 13/10/15 and Writes union 15 were accepted; IEA partitions 13/8/14 and Writes union 13 were accepted.
It accepted canonical closure, first-use freshness and proof obligations, and the complete-LTE timing clarification.
Planning may precede LTE landing, but concrete landed API binding and headroom must finish before any implementation charge.
No partial provider, private fallback, TS dependency, or build authority follows from this planning acceptance.

Retained `gpt-6.1-sol / high` closes only metadata, with at most two metadata/validation iterations authorized per workspace.
All reviewed ticket bytes, behavior rows, fences, and closed contracts stay unchanged.
The complete reviewed graph and all predecessors remain retained in Git.
Redundant full-text spec copies are deleted after verification against the complete accepted source.
The shared observations decision map and resolved topic tickets remain in place for both outcomes.

Stage-close checks cover authored prose, row/fence/ticket preservation, diff consistency, and the required planning lane.
Unchanged broad preflights and runtime checks are not rerun.
No implementation, mutation, test suite, manual gate, benchmark, retro, commitment change, or landing is performed.
The workspace handoff returns to `$bench-mainenance` for authorized planning continuation.

<!-- command-currency: historical -->
