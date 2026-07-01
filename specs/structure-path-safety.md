# structure path safety

## Problem

`bench structure` is an agent-facing CLI signal, so its output must be trustworthy.
The file-length loop reads file paths safely enough for ordinary names, but the
crowded-directory check pipes paths through whitespace-splitting shell tools. A
directory named `space dir` with 13 source files is reported as crowded `./` and
`dir/`, which sends the agent to the wrong module.

## Solution

Deepen the structure module's path handling so file lists are treated as paths, not
words. Directory counts should be computed without `xargs` or any whitespace split,
and touched-scope structure checks should use the same safe path representation as
the full-repo check. The CLI output stays the same shape, but the paths in it are
correct.

## User stories

1. As a user running `bench structure`, I want directories with spaces in their names
   to be counted and reported correctly.
2. As a user reading structural-debt output, I want each reported path to identify
   the actual module that needs attention.
3. As a shift user, I want touched-scope structure checks to use the same safe path
   handling as whole-repo checks, so refactor prompts do not drift because of a path
   spelling edge case.
4. As a kit maintainer, I want a runtime gate check for a path-with-spaces case, so
   future shell rewrites do not reintroduce whitespace splitting.

## Implementation decisions

- **Primary seam:** the `bench structure` CLI contract. Tests create a throwaway git
  repo and assert output.
- Remove whitespace-splitting path handling from directory counting. A shell loop,
  array, or small structured parser is acceptable as long as pathnames remain intact.
- Prefer one internal file-list representation for whole-repo and touched-scope
  structure checks. If full NUL-delimited plumbing is too large for this slice, the
  implementation must at least cover spaces and document any remaining newline-path
  limitation.
- Keep thresholds, output labels, and exit semantics unchanged.

## Testing decisions

- **Good tests here** exercise the real CLI in a throwaway repo. They do not test the
  helper function directly.
- **Seam:** `bench structure`, because users and the ambient dashboard consume its
  output and exit code.
- **Prior art:** `.bench/gate-runtime-contracts.sh` already has a shell-file length
  contract for `bench structure`. Add a crowded-directory case with spaces.
- **Gate command:** `bench gate`.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1, 2 | A crowded directory named `space dir` is reported as `space dir/`, not split into `./` and `dir/`. | `bench structure` CLI | Observed red before implementation: a throwaway repo with 13 `space dir/*.sh` files and `BENCH_MAX_DIR_FILES=12 bench structure` reports `./` and `dir/`. | The current `xargs dirname` pipeline splits on whitespace, so the output names the wrong directories. |
| 3 | Touched-scope structure checks share the safe path handling. | `bench shift` structure phase via structure helper | Not TDD-able before implementation if the chosen internal representation changes; cover with a focused helper or shift contract only if the touched-scope path is modified. | The same module handles both full and touched structure checks, so the fix should not leave the shift path behind. |
| 4 | The project gate fails if path-with-spaces directory counting regresses. | Project gate: runtime contracts | Not TDD-able before implementation beyond the observed red CLI probe. | The committed throwaway-repo case protects the CLI behavior at the seam. |

## Out of scope

- Changing structure thresholds or adding new structural signals. This spec only
  fixes path correctness, ~0 extra minutes.
- Redesigning structure output format. The existing labels and exit codes remain the
  agent-facing interface.
- Supporting every possible newline-containing git path if it requires broad
  plumbing changes. Decide that separately if the implementation cannot keep the
  slice small, ~45 minutes.
