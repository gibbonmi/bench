# Slicing checks

Charged from `craft-tickets` when the slicer writes each ticket's `Writes:` line and
the spec's ownership fence. The first section lists every `Writes:` rule that the ticket
parser and build preflight enforce, and some `Blocked by:` and `Covers:` rules beside them.
The second section states the slicing rules that no check enforces, each from a repair
round that a retrospective recorded. The third section lists the evidence classes that
the research of each ticket chunk returns, and the rule that keeps each research path.

## Enforced rules

- `Blocked by:` holds `none` or sibling ticket file basenames. A basename survives a
  retitle, and `--ticket` already names it.
- Each ticket states `Writes:` exactly once, beside `Blocked by:` and `Covers:`.
- Each `Writes:` path exists in the tree or carries the `(new)` marker.
- Each `Writes:` entry holds only text that spec TOON can render. Preflight refuses the
  ticket set before it renders a verdict.
- A `Writes:` entry that names a fixture-pinned path also names the fixture directory
  that pins it.
- A `Writes:` entry that names a bound package also names every registry file that the
  binding registry binds to that package.
- A `Writes:` entry for an anchored guidance path, or a directory above one, names every anchor registry file that names it, test files included.
- The spec's ownership fence equals the union of the ticket `Writes:` paths, less the review pickup and every path under the spec folder or `capture`.
- A ticket that writes a system-tagged test file states `BENCH_KIT`.
- `Covers:` holds `none` or declared row IDs. Cite each ID in full, because preflight
  reads IDs, not ranges.

The `--propose-writes` form of build preflight lists each missing closure file with its
source and its fence state. A trailing `/` and the `(new)` marker never make two entries
differ.

## Slicing rules

- A lane-check ticket proves its check through the real lane over a composed tree. A
  stand-in harness can pass while the lane refuses the same tree.
- A posture change that reds a fixture helper makes the slicer list every call site of that helper before the map locks. A posture change for new output also names each stream that the output reaches. An unlisted call site or stream reds at the first gate run of the build.
- A combined behavior row belongs to the ticket that completes its final consumer. An
  earlier ticket cannot prove a behavior whose consumer does not exist yet.
- A retirement pass gives each sentence that grants the retired behavior its own forbid row and red-capable check. One forbid row for a family of sentences leaves the others free to return.
- A cited verifier row names the exact checks it performed. A reviewer can then repeat
  the checks and compare the results.
- The slicer runs `bench preflight build <slug>` after each change to a ticket or an acceptance row.
  Before the ticket graph goes to approval, and before each plan-expansion commit, the slicer also runs
  `bench preflight build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>`
  for each affected ticket. The slicer adds each listed closure file to the `Writes:` line of that ticket.
  The spec and its tickets go to sign-off and to their landing only on a green preflight of the current ticket graph.
  Sign-off then grades the final fence, not an earlier one.
- A plan commit that renames a spec symbol runs one `rg` for the old symbol across every ticket, and it corrects each hit.
- A fence expansion reads the fence section of the spec and corrects each sentence that names the added package.
- Each acceptance claim names an available producer and an executable check at its ticket checkpoint. A claim whose producer lands in a later ticket has no evidence at that checkpoint.
- A ticket that preserves or moves a ref names a check for each of the head, the tip, and the checked-out-ref shapes. A check on one shape leaves the other shapes free to fail.
- A rule for one member of an enumerated family goes on its declaration row, and `Writes:` names the declaration table and its dispatch test.
- A ticket that adds a second caller of another package's private rule names that package in `Writes:`, so one exported rule serves both callers.
- A ticket that changes a rendered output shape runs one `rg` for the old text, and each hit joins that ticket's own `Writes:` line.
  A rendered output shape is a block name, a row kind, a table header, a line format, or a help line.
  The search covers the tests, the help inventories, the generated references, and their anchor registry files.
  The spec-wide fence masks a hit that is absent from one ticket's `Writes:` line, so the fence alone does not close the gap.
- A ticket that moves a test expectation names the file that holds the old expectation in its own `Writes:` line. Otherwise the old expectation stays behind as a second copy.
- The ownership fence and exactly one ticket's `Writes:` line name `.bench/commitment.json`. The spec-stage close commits the binding inside the range that the plan-only preflight grades.

## Fence-expansion evidence

`/bench-write-spec` charges read-only research for each ticket chunk. Each research return lists every likely fence expansion in one of these classes:

- a caller file that `bench consumers` names;
- a fixture or a canary that pins a path;
- an anchor registry file that names a guidance path;
- a count pin, such as a package test-count pin;
- a file over its structure budget, with its headroom;
- a shared file that more than one ticket writes, such as `CHANGELOG.md`.

Each entry cites its evidence as a `path:line` or as a command and its result.

Each path that a research return names goes into the ticket's `Writes:` line. The slicer drops such a path only with `path:line` evidence that the path cannot change under any design. `Writes:` is the approved upper bound of the ticket, not a task list. `/bench-implement-spec` owns how the ticket author narrows it.
