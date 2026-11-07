# Slicing checks

Charged from `craft-tickets` when the slicer writes each ticket's `Writes:` line. The
first section states the slicing rules that no check enforces, each from a repair round
that a retrospective recorded. The second section lists the evidence classes that the
research of each ticket chunk returns, and the rule that keeps each research path.

## Slicing rules

- A lane-check ticket proves its check through the real lane over a composed tree. A
  stand-in harness can pass while the lane refuses the same tree.
- A posture change that reds a fixture helper makes the slicer list every call site of that helper before the map locks. A posture change for new output also names each stream that the output reaches. An unlisted call site or stream reds at the first gate run of the build.
- A combined behavior row belongs to the ticket that completes its final consumer. An
  earlier ticket cannot prove a behavior whose consumer does not exist yet.
- A retirement pass asks for a Forbid row only for a stale or wrong form. Each sentence that grants the retired behavior is such a form, and it gets its own Forbid row and red-capable check. One Forbid row for a family of sentences leaves the others free to return. The pass adds no Require row for the new rule wording, because review owns rule wording.
- A cited verifier row names the exact checks it performed. A reviewer can then repeat
  the checks and compare the results.
- The slicer runs `bench preflight build <slug>` after each change to a ticket or an acceptance row.
  The spec and its tickets go to sign-off and to their landing only on a green preflight of the current ticket graph.
  Sign-off then grades the final fence, not an earlier one.
- Each row that the `rows-owned` red names goes to the ticket that writes the seam path of that row.
- A plan commit that renames a spec symbol runs one `rg` for the old symbol across every ticket, and it corrects each hit.
- Each acceptance claim names an available producer and an executable check at its ticket checkpoint. A claim whose producer lands in a later ticket has no evidence at that checkpoint.
- A ticket that preserves or moves a ref names a check for each of the head, the tip, and the checked-out-ref shapes. A check on one shape leaves the other shapes free to fail.
- A rule for one member of an enumerated family goes on its declaration row, and `Writes:` names the declaration table and its dispatch test.
- A ticket that repeats an existing rule, such as an adapter, a validator, or a predicate, names its source in `What to build`.
  `Writes:` names the package of that source, so one exported rule serves both callers. A pasted copy is a second source for one fact.
- A ticket that changes a rendered output shape runs one `rg` for the old text, and each hit joins that ticket's own `Writes:` line.
  A rendered output shape is a block name, a row kind, a table header, a line format, or a help line.
  The search covers the tests, the help inventories, the generated references, and their anchor registry files.
  The ticket union masks a hit that is absent from one ticket's `Writes:` line, so the union alone does not close the gap.
- A ticket that moves a test expectation names the file that holds the old expectation in its own `Writes:` line. Otherwise the old expectation stays behind as a second copy.
- Exactly one ticket's `Writes:` line names `.bench/commitment.json`. The spec-stage close commits the binding inside the range that the plan-only preflight grades.

## Fence-expansion evidence

`/bench-write-spec` charges read-only research for each ticket chunk. Each research return lists every likely fence expansion in one of these classes:

- a caller file that `bench consumers` names;
- a fixture or a canary that pins a path;
- an anchor registry file that names a guidance path;
- a count pin, such as a package test-count pin;
- a file over its structure budget, with its headroom;
- a shared file that more than one ticket writes, such as `CHANGELOG.md`;
- a test file that goes red in a scratch-patch test run.

A scratch-patch test run applies to a ticket that adds a refusal or retires a rule, a fixture shape, or a row. Before the slicer locks that ticket, the research delegate applies the refusal or the retirement as a scratch patch in its own Bench worktree. The delegate runs the reached packages once with one `bench test` call, and it returns each file that goes red. The delegate then discards the patch and lands nothing. `bench worktree release` refuses that unlanded worktree, so the delegate plans the discard with `bench worktree clean --discard-branch [--discard-ignored] <path>`. The delegate then runs the same command with `--apply <fingerprint>` and the fingerprint that the plan prints.

Each entry cites its evidence as a `path:line` or as a command and its result.

Each path that a research return names goes into the ticket's `Writes:` line. The slicer drops such a path only with `path:line` evidence that the path cannot change under any design. `Writes:` is the expected write set of the ticket, not an upper bound or a task list. `/bench-implement-spec` owns how the author narrows it, and `.bench/BENCH.md`'s plan-expansion policy owns how the build widens it.
