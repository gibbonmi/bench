# Charge the delegate-side duties in craft-delegate

Blocked by: add-the-contracts-discovery-step.md
Ownership fence: `.agents/skills/bench-craft-delegate/SKILL.md`, `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/fixture_bite_test.go`
Assumptions: the backup rule extends the existing per-worktree stash-substitute paragraph rather than replacing it; the probe-kind sentence's mutation-table entry in `fixture_bite_test.go` hard-codes that sentence's current line wrapping and asserts an anchor count of exactly 1; needles register through `scopedSection`/`requireInSection`. Re-derive from the tree at pickup.

## What to build

FT164 story 5: the charge sections of `craft-delegate` carry the duties that
catch a silent green and a cross-fence clobber at return time, because charge
prose is the only surface that reaches a low-context delegate.

The charge names the central-property-breaking mutation, and the delegate
applies it to its own finished work, reports the observed result, and adds the
missing row when the mutation comes back green — three fixture-shaped silent
greens were caught only because the charge said *run* the mutation rather than
reason about it. The coordinator's own probe differs in kind **and** in site
from any probe the delegate ran, with the omission/swap probe-kind vocabulary
carried in the charge template: three same-site probes were vacuous on the first
try, and a vacuous probe is indistinguishable from a pass. A charge that extends
an enumerated family names every registry the family already appears in, found
by tracing one existing sibling. Transient backups live inside the delegate's
own worktree under a unique name, and restores name exact files and never globs
— a stale scratchpad swept into a later restore glob clobbered four out-of-fence
files.

The existing probe-kind sentence is an anchor artifact: its mutation-table entry
in `fixture_bite_test.go` hard-codes that sentence's current line wrapping. This
ticket decides the sentence's final location in the diff, and the entry moves
with it in that same diff, holding its count-of-exactly-1 assertion against the
sentence's final bytes and wrapping.

Enforcement is anchor-first: five needles scoped to their owning sections, each
with one byte-exact mutation-table row proving its diagnostic fires.

## Acceptance

- [ ] [DL1] the charge names the central-property mutation and the delegate applies it to its finished work, reports the result, and adds the missing row on a silent green.
- [ ] [DL2] the coordinator probe is required to differ in site as well as in kind from any probe the delegate ran.
- [ ] [DL3] the charge template carries the omission/swap probe-kind vocabulary.
- [ ] [DL4] a family-extending charge names every registry found by tracing an existing sibling.
- [ ] [DL5] transient backups live in the delegate's own worktree under a unique name and restores name exact files, never globs.
- [ ] [DL6] the probe-kind sentence's mutation entry moves with the sentence in this diff and still counts exactly one match against its final bytes and wrapping.
- [ ] [DL7] every needle this ticket registers has a byte-exact entry in `TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting` whose named diagnostic fires.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DL1 | swap "apply the mutation and report the result" for "consider whether the mutation would fail" | the `delegate self-probe` mutation subtest | swap the sentence, run the anchor check, expect the self-probe diagnostic |
| DL2 | delete the site clause, leaving kind alone | the `probe site differs` mutation subtest | delete the clause, run the anchor check, expect the site-differs diagnostic |
| DL3 | delete the omission/swap vocabulary from the charge template | the `probe kind vocabulary` mutation subtest | delete the vocabulary line, run the anchor check, expect the vocabulary diagnostic |
| DL4 | soften registry tracing to "check the obvious registries" | the `registry tracing` mutation subtest | swap the sentence, run the anchor check, expect the registry-tracing diagnostic |
| DL5 | permit a glob restore outside the worktree | the `backup isolation` mutation subtest | swap the sentence to allow a glob, run the anchor check, expect the backup-isolation diagnostic |
| DL6 | rewrap the probe-kind sentence and leave its entry unedited | the `probe kind` mutation subtest | rewrap the sentence, run the subtest, expect the anchor-count failure naming a count other than 1 |
| DL7 | register a needle and land no mutation entry for it | `review` plus the mutation harness | delete one entry, run the harness and watch the remaining rows still pass, then read this ticket's needle list against the table |
