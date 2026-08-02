# Honor slots and attestation across a process boundary

Blocked by: Skip the build phase on its attested seal
Ownership fence: `internal/contract/runtime/runtime_gate_component_boundary_test.go`
Assumptions: the runtime contract suite execs the built CLI against throwaway
fixture repos; `runtime_gate_reduced_test.go` is the prior art for the
CLI-observable announcement lines. Re-derive from the tree at pickup.

## What to build

Every piece of evidence this build adds is serialized and re-read by a later
process, and unit-level green hides serialization defects that appear only on
reload. This ticket proves the boundary: one CLI process authors the slots and
the attestation, a second, fresh CLI process reads them back and honors them,
and a third refuses forged ones.

The observable is the CLI's own output and exit code — the announcement lines
naming each skipped component and its evidence, and the executed set the second
run reaches — not an in-process return value.

**Evidence authorship.** `bench gate` is the producer in both runs; the test
asserts that the second process consumed exactly what the first published, and
that a `--fresh` run re-authors rather than inheriting.

## Acceptance

- [ ] [PC21a] slots and an attestation authored by one CLI process are honored by a fresh CLI process, which skips the covered components and announces each with its evidence.
- [ ] [PC21b] a forged slot and a forged attestation written between the two runs are refused by the second process, which runs the affected components.
- [ ] [PS34] the second process's announcement lines name the same ancestor identities the first process authored.
- [ ] [PS35] a `--fresh` third run executes everything and re-authors every slot and the attestation, observable from the CLI.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC21a | hold the resolved partition in a process-lifetime cache instead of reading the store | `TestPartialEvidenceSurvivesAProcessBoundary` | exec the CLI gate twice against one fixture repo with a capture-only edit between, assert the second run's skip lines |
| PC21b | validate the slot only at authorship | `TestForgedEvidenceIsRefusedOnReload` | exec once, overwrite a slot with a foreign-component record, exec again, assert the component ran |
| PS34 | print the current identity rather than the recorded one | `TestAnnouncedAncestorsMatchWhatWasAuthored` | exec twice, parse the second run's skip lines, compare against the store's slot contents |
| PS35 | let `--fresh` reuse the attestation | `TestFreshReauthorsAcrossProcesses` | exec twice, record store bytes, exec with `--fresh`, assert every slot and the attestation moved |
