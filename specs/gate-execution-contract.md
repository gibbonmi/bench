# gate execution contract

## Problem

The gate is the oracle, but `bench gate` does not run every oracle command from the
same place. It resolves the repo root, then runs `.bench/gate.sh` and `$BENCH_GATE`
from the caller's current directory, while auto-detected Node/Python/Rust gates run
from the repo root. A project gate that reads repo-relative files can pass from the
root and fail from a subdirectory. That leaks cwd knowledge across the gate seam.

## Solution

Deepen the gate execution module. `bench gate` always runs the selected oracle from
the repo root, regardless of where the user invoked Bench. The same behavior is used
when `bench shift` calls the gate and when the Stop hook shells out to `bench gate`.
Exit codes remain the oracle verdict: 0 means green, non-zero means red or missing
gate.

## User stories

1. As a user running `bench gate` from a subdirectory, I want a repo-local
   `.bench/gate.sh` to run from the repo root, so repo-relative paths inside the
   gate behave the same everywhere.
2. As a user setting `$BENCH_GATE`, I want the command string to run from the repo
   root, so the env override has the same cwd contract as `.bench/gate.sh`.
3. As a user relying on auto-detection, I want the existing root-run behavior
   preserved.
4. As a shift user, I want `bench shift` to keep using the same gate execution
   contract as `bench gate`, so a green shift and a manual green gate mean the same
   thing.
5. As a kit maintainer, I want the gate to exercise root-cwd behavior in throwaway
   repos, so this regression cannot return silently.

## Implementation decisions

- **Primary seam:** the `bench gate` CLI contract. Tests call the real CLI from a
  subdirectory and observe exit code, not helper internals.
- The gate execution module owns cwd. After resolving the repo root, it runs
  `.bench/gate.sh`, `$BENCH_GATE`, and auto-detected gates from that root.
- `.bench/gate.sh` keeps receiving no extra arguments. The interface change is cwd
  only.
- `$BENCH_GATE` remains a command string evaluated by `bash -c`, but it is evaluated
  after entering the repo root.
- Missing-gate behavior and exit code 3 stay unchanged.

## Testing decisions

- **Good tests here** exercise the external CLI seam: create a throwaway repo, add a
  gate that reads a root-relative file, invoke `bench gate` from a nested directory,
  and assert the exit code.
- **Seam:** `bench gate`. `bench shift` and the Stop hook compose that seam, so the
  root-cwd regression is best protected at the highest common interface.
- **Prior art:** `.bench/gate-runtime-contracts.sh` already creates throwaway repos
  for runtime CLI contracts.
- **Gate command:** `bench gate`.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A repo-local `.bench/gate.sh` that checks a root-relative file passes when `bench gate` is invoked from a subdirectory. | `bench gate` CLI | Observed red before implementation: a throwaway repo with `package.json`, `.bench/gate.sh` containing `[ -f package.json ]`, and invocation from `sub/` exits 1. | The current implementation runs `.bench/gate.sh` from `sub/`, so the root-relative file is missing. |
| 2 | `$BENCH_GATE` runs from the repo root when invoked from a subdirectory. | `bench gate` CLI | Observed red before implementation: `BENCH_GATE=./gate-root.sh bench gate` from `sub/` exits 127 because `./gate-root.sh` is resolved under `sub/`. | The command string currently inherits caller cwd; running it at root makes the relative command resolvable. |
| 3 | Auto-detected gates keep their existing root-cwd behavior. | `bench gate` CLI | Already covered by inspection and existing implementation: auto-detect branches already wrap commands in `cd "$root"`. | This row protects against regressing the branch while unifying the others. |
| 4 | `bench shift` continues after a green project gate rather than changing loop control. | `bench shift` CLI | Already covered by existing runtime contracts that assert `shift done` after a green gate. | The gate execution change must not reintroduce the earlier loop-control bug around running the gate in-process. |
| 5 | The runtime contract fails if either `.bench/gate.sh` or `$BENCH_GATE` root-cwd behavior regresses. | Project gate: runtime contracts | Not TDD-able before implementation beyond the two red CLI probes above. | The gate check is the committed protection; its red capability is demonstrated by the failing probes. |

## Out of scope

- Changing gate cache semantics or `bench status` behavior. That is a separate
  ambient-dashboard capability, ~30 minutes.
- Redesigning `$BENCH_GATE` away from a shell command string. That is a separate
  config-interface decision, ~45-60 minutes.
- Enforcing token caps in `bench shift`. That is tracked in
  `decisions/dogfood-improvements.md` #12.
