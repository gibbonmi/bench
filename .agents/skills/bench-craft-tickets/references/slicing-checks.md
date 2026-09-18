# Slicing checks

Charged from `craft-tickets` when the slicer writes each ticket's `Writes:` line and
the spec's ownership fence. The first section lists the rules that the ticket parser and
build preflight enforce. The second section states the slicing rules that no check
enforces, each from a repair round that a retrospective recorded.

## Enforced rules

- `Blocked by:` holds `none` or sibling ticket file basenames. A basename survives a
  retitle, and `--ticket` already names it.
- Each `Writes:` path exists in the tree or carries the `(new)` marker.
- A `Writes:` entry that names a fixture-pinned path also names the fixture directory
  that pins it.
- A `Writes:` entry that names a bound package also names every registry file that the
  binding registry binds to that package.
- A `Writes:` entry that names an anchored guidance path also names every anchor registry file whose string literal names that path, test files included.
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
- A posture change that reds a fixture helper makes the slicer list every call site of that helper before the map locks. An unlisted call site reds at the first gate run of the build.
- A combined behavior row belongs to the ticket that completes its final consumer. An
  earlier ticket cannot prove a behavior whose consumer does not exist yet.
- A retirement pass gives each sentence that grants the retired behavior its own forbid row and red-capable check. One forbid row for a family of sentences leaves the others free to return.
- A cited verifier row names the exact checks it performed. A reviewer can then repeat
  the checks and compare the results.
- The slicer runs build preflight again after each fence change and before review. Review
  then grades the final fence, not an earlier one.
