# Repair the accepted review findings

Blocked by: route-skills-index-verb-and-retire-script.md
Writes (advisory): internal/skillsindex/skillsindex.go, internal/skillsindex/command.go, internal/skillsindex/skillsindex_test.go, internal/skillsindex/command_test.go, internal/conformance/checks_test.go, internal/conformance/skills_index_checks_test.go, cmd/bench/command_registry_test.go, specs/skills-index-reader/spec.md

## What to build

Six accepted repair targets from the 2026-08-15 three-axis review of base `5b41322a`
against tip `2ca77cb2`. Reviewer scoped the repair to diff-owned findings only; the ten
inherited edges (behavior byte-identical to the pre-collapse implementation) are parked
as a roadmap item and are explicitly out of scope here.

1. `bench skills-index --check --write` currently writes. Both help surfaces advertise the
   modes as exclusive, so the conflicting pair exits 2 with usage.
2. `markerBlock` and `replaceMarkerBlock` each scan for the marker pair independently and
   diverge on a duplicated block. Derive the span once and have check and write consume it.
3. The index-line format is encoded twice: rendered with `fmt.Sprintf`, then re-parsed with
   `indexLineRe` to recover names already held in `Entry`. Carry expected names from
   `Entry`; leave one owner for recognising an actual committed line's skill name.
4. Six comments narrate history or cite provenance and three restate forwarding code.
   Rewrite to timeless present per `craft-comments`.
5. SI3's mandated "exactly the three attributed lines" fixture was extended into a
   four-line case by the collision repair. Keep both: restore the exact three-line
   assertion and add the collision as its own case.
6. SI7 names the `cmd/bench` dispatch seam but the test calls `skillsindex.Command`
   in-process. Add the dispatch-level case so wrapper routing and fresh-process reload
   cannot regress silently.

Also record two reviewer vetoes in `spec.md`: the permitted-edit list widens to cover the
`skillNameFromIndexLine` deletion, and SI4's "already covered" claim is accepted as
inaccurate — the four canaries do not exercise the non-attributable `block drifted from
generated form` case.

## Acceptance

- [ ] `--check --write` exits 2 with usage and writes nothing; covered by a new SI7 case
- [ ] one owner derives the marker span; one owner encodes the index-line format
- [ ] no comment narrates the change, cites provenance, or argues its own correctness
- [ ] SI3 asserts both the exact three-line ordering and the collision case
- [ ] the verb is exercised through `cmd/bench` dispatch
- [ ] spec.md records both vetoes; no other row's behavior changes
