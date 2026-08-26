# Review pickup: exec-census

Base `3d76ff1b`, reviewed tip `859fca4a`. Three axes, one raw finding each; two repair targets after collapse. The reviewer was unreachable, so each disposition below stands for veto.

## Standards

Count: 1. Worst issue: the census sentence appears in both the reference and the profile.

- `.bench/BENCH-reference.md:39` and `projects/benchkit.md:100` carry the same sentence. Citation: `AGENTS.md` code standard, "Two derivations of the same fact must collapse into one source." Disposition: `no-op`. Spec story 34 and row EC33 require the sentence in both files; the spec is the closed decision. Veto surface: name one file canonical in a later drain.

## Spec

Count: 1. Worst issue: row EC27 overclaims its anchor coverage.

- Row EC27 names "the four learning fields" and "the retro citation" under the anchors registry test. `internal/anchors/registry_data.go` anchors only the duty sentence and the zero close. Citation: `specs/exec-census/spec.md`, row EC27, and story 33, "an edit cannot drop the duty silently." Disposition: `auto-fix`. Repair ticket: `specs/exec-census/tickets/08-anchor-the-learning-fields-and-the-retro-citation.md`.

## Coverage

Count: 1. Worst issue: a text that names two assignment ids records only the first.

- Input `cat <pool>/<o>-<a>/x <pool>/<o>-<b>/y` writes one line under `a` and none under `b`; no row or edge line decides it. Citation: `internal/census/census.go`, `assignment`, returns on the first match. Disposition: `auto-fix` of the pin, decision for veto: story 2 counts one call as one record, so the first id owns the record. Repair ticket: `specs/exec-census/tickets/09-pin-the-first-assignment-id-rule.md`.
