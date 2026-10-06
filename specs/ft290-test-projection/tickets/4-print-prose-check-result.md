# 4. Print the prose check result

Blocked by: 3-prove-named-check-ran.md
Writes: internal/testreport/, internal/prose/walk.go, internal/prose/walk_test.go
Covers: TP7, TP8, TP9, TP10, TP11, TP12

## What to build

Chunk: TP-C1b.

Make the prose grader return the graded subject paths beside its findings, and keep `prose.Grade` for its conformance caller.
A graded subject is a `.md` file that `collect` returns and that the grader's exclusion set does not exclude. In a git root, `collect` returns tracked files only.
Reuse `collect` and the grader's exclusion set inside `walk.go`. Do not write a second walk or a second exclusion test.
The prose check prints the `check` row with `tests_run` 0 and the real `subjects` count.

A green result prints only that row. `--full` adds the `subjects[N]{path}` table in sorted order.
A red result prints the row and then each finding line. A grader refusal prints `subjects` 0 and then its diagnostics, at exit 1.

A run with zero subjects prints the row and the title `named check ran nothing`, and exits 1.

## Acceptance

- [ ] A green run over two subjects prints exactly the row `prose,prose,0,2` at exit 0.
- [ ] `--full` lists the two subject paths in sorted order.
- [ ] A tree with zero graded subjects exits 1.
- [ ] A red run keeps each finding after the `check` row.
- [ ] A tree with one `.md` file and no `.bench/prose-exclusions` file prints `subjects` 0 and the grader diagnostic at exit 1, without the title `named check ran nothing`.
- [ ] The grader does not count an excluded file as a subject.
