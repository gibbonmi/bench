# Exclude the planning documents from the prose mechanics check

Blocked by: none
Writes: .bench/prose-exclusions, internal/conformance/prose_mechanics_test.go
Covers: none

## What to build

The prose mechanics check grades every tracked Markdown file. Planning documents
pay the sentence bound and the paragraph bound at each commit. The reviewer decided
on 2026-10-16 that the check does not grade the planning documents.

Add four permanent rows to the prose exclusion file: `specs/`, `roadmap/`, `capture/`,
and `docs/`. Each row uses the existing row grammar and gives a one-clause reason.
Add the same four paths to the approved row set in the conformance test, because
that set is independent of the file and an added row reds the test.

The `.bench/BENCH.md` file stays unchanged, because it ships to every linked repo.
The learning-journal structure check does not read the exclusion file, so it
continues to grade the journal structure.

The `bench learning` verb grades a new entry through the same exclusion file, at the
journal path under `capture/`. In this repo, the verb therefore stops refusing an
entry that is over the prose bounds. This result agrees with the sweep, which does
not grade the ignored journal.

## Acceptance

- [ ] The prose mechanics check reports no finding for a Markdown file under
  `specs/`, `roadmap/`, `capture/`, or `docs/`.
- [ ] The approved-set test passes with the four new rows and reds when a row is
  removed from its approved set.
- [ ] The learning-journal structure test passes without change.
