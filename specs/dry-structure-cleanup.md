# dry-structure-cleanup

## Problem

`bench structure` is red: the AXI gate-contract fragment (486 lines) and the
git guard hook (440 lines) are over the 400-line file budget. Both overruns are
knowledge-duplication debt, now a graded defect under the working agreement's
one-source-per-fact standard: the gate fragment pastes the same fixture harness
(~20 copies of mktemp / subshell / git init / err / cleanup), and the guard
carries a ~390-line Python analyzer as a bash heredoc string — a program no
linter, parser, or unit seam can see.

## Solution

Collapse the fixture harness to one contract runner shared by the AXI
fragments, split the AXI fragment along its existing responsibility banners,
and extract the analyzer to a real Python file the hook invokes and the gate
parses. `bench structure` exits green; every existing contract, verdict, and
canary keeps passing unchanged.

## User stories

1. As a kit maintainer, I want the AXI fixture harness defined once as a
   contract runner that the AXI fragments share, so the pasted
   mktemp/subshell/cleanup blocks collapse to one source per the
   one-source-per-fact standard.
   Line: claude-opus-4-8 / medium. This refactors the oracle's own test
   harness, and the profile routes gate and conformance logic to the mid tier.

2. As a kit maintainer, I want `gate-axi-contracts.sh` split along its section
   banners into per-surface fragments, with the spec-relative "Row N"/"Fix N"
   comment labels replaced by self-contained contract prose, so each file is
   one responsibility under budget and readable cold.
   Line: claude-sonnet-4-6 / low. This is a mechanical move fully observed by
   the existing contracts and canaries; run inline or declare a bump to mid.

3. As a kit maintainer, I want the git-guard analyzer extracted from the bash
   heredoc into a sibling Python file the hook invokes, so the classifier is a
   lintable, parseable single source and the hook shrinks to its ~50 bash
   lines.
   Line: claude-opus-4-8 / medium. The guard is fail-closed enforcement whose
   correctness the profile rates oracle-adjacent.

4. As a reviewer, I want the hook to degrade honestly when the analyzer file
   is missing or empty — `--describe` reports the manifest unavailable and
   exits 0; enforcement blocks (fail closed) — so a broken install can never
   silently grant destructive-git authority.
   Line: claude-opus-4-8 / medium. Failure-path semantics of a guard are
   exactly the class of bug the gate exists to prevent.

5. As a kit maintainer, I want the gate to parse the extracted analyzer
   (py_compile, best-effort when python3 is present, mirroring the shellcheck
   layer) so a syntax error is caught at gate time, not on the next Bash call.
   Line: claude-sonnet-4-6 / low. One gate line plus one canary, fully
   gate-observed.

6. As a reviewer, I want `bench structure` to exit 0 when the build lands, so
   structural debt returns to zero and stays enforced by the existing budget.
   Line: covered by stories 2–3; no separate build work.

## Implementation decisions

- **Contract runner.** A new gate helper file, sourced once by the gate before
  the fragments. Interface (the decision, trimmed):

      contract "<label>" [--space-path] <<'BODY'
        ... assertions; $root and $tmp in scope; runs under set -u in a
        subshell, cwd = fresh git-inited fixture dir ...
      BODY

  The runner owns mktemp, git init, the `|| err "<label> failed"` reporting,
  and cleanup — cleanup runs on failure too (no leaked scratch dirs).
  `--space-path` provisions the fixture under a space-containing parent.
  Blocks that source `bin/bench-query.sh` directly instead of running the CLI
  use the same runner; only their bodies differ.
- **Adoption scope.** The runner is adopted by the AXI-family fragments this
  build touches (the split fragments plus the wave-2 AXI fragment). The five
  non-AXI fragments migrate in a follow-up (see Out of scope).
- **Split shape.** The current AXI fragment's banners are the seam:
  guard-manifest conformance + query/parser contracts stay in the existing
  fragment; guards aggregation + session-start contracts move to a new
  fragment. Each fragment keeps its own CLI-missing early return, and the
  guard-manifest section keeps running before any early return (the minimal
  canary fixture depends on that ordering).
- **Analyzer extraction.** The Python program moves verbatim to a sibling file
  in the hooks directory. The hook resolves it relative to its own location
  (BASH_SOURCE), quoted against spaces. `DENY_LABELS` stays the single source
  for enforcement and advertisement, now living only in the Python file. The
  hooks tree ships whole (package files list and `bench link` both copy the
  directory), and `bench guards` globs `*.sh`, so the `.py` neither breaks
  shipping nor appears as a phantom guard.
- **Fail-closed matrix for the analyzer file.** python3 missing → existing
  degradation unchanged. Analyzer file missing *or empty* → `--describe`
  prints a `manifest unavailable (analyzer missing)` denies line and exits 0;
  the enforcement path emits BLOCKED and exits 2. An empty file is explicitly
  in the deny branch because an empty program exits 0 printing nothing, which
  the allow path would misread as "no verdict".
- **Gate parse layer.** py_compile on the analyzer joins the parse checks,
  conditional on python3 being present (same best-effort posture as
  shellcheck).
- **Directory budget risk.** The split adds two files to the gate directory;
  if that trips the crowded-directory budget, surface it rather than silently
  re-merging — grouping fragments into a subdirectory is a reviewer call.

## Testing decisions

- A good test here runs the real gate or the real hook against a fixture and
  asserts exit code plus output shape — never internals. The refactor stories
  are graded by the existing contracts continuing to pass and the existing
  canaries continuing to bite; the new behavior (analyzer-missing degradation,
  py_compile) gets new red-first contracts.
- Seams: the gate contract (highest), the guard hook CLI (stdin JSON → exit
  code/stderr; `--describe` → manifest), and `bench structure`'s exit code.
  Prior art: the AXI contract fragments, the allow/block verdict matrix in the
  runtime-git fragment, and the canary layer.
- Gate command: `bench gate`, plus `bench structure` exiting 0 as the
  structural acceptance signal.

### Seam diagram

Gate contract (stories 1, 2, 5):

    trigger: bench gate (reviewer, Stop hook, shift loop)
        │
        ▼
    repo tree ──▶ [ gate.sh + fragments (runner, split, py_compile) ] ──▶ exit 0/1 + per-check err lines
    canary fixtures ──▶ [ same gate, inner run ]  ──▶ red + targeted substring
                      ◀ tests attach here: run the gate on this tree (green) and on
                        tests/canary fixtures (each must stay red with its EXPECT)

Guard hook CLI (stories 3, 4):

    trigger: PreToolUse:Bash (harness) · bench guards (--describe)
        │
        ▼
    stdin JSON {command} ──▶ [ block-dangerous-git.sh ──invokes──▶ git-guard.py ] ──▶ exit 0 | exit 2 + BLOCKED
    --describe           ──▶ [                                                  ] ──▶ 4-key manifest, exit 0
                      ◀ tests attach here: pipe fixture JSON / pass --describe,
                        assert exit code and output (verdict matrix + AXI contracts)

`bench structure` (stories 2, 3, 6):

    trigger: bench structure (reviewer, dashboard)
        │
        ▼
    tracked file sizes ──▶ [ structure budgets ] ──▶ exit 0 | exit 1 + FILE TOO LONG rows
                      ◀ tests attach here: run it; red today on both files, green defines done

### Acceptance coverage map

Granularity note (per the quantifier rule): "every existing contract/canary"
below means, concretely: all contracts in the two AXI fragments being
refactored and the wave-2 fragment, the full allow/block verdict matrix in the
runtime-git fragment, and every fixture under `tests/canary/` — each enforced
individually by the gate's own aggregation, not sampled.

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | refactored fragments still enforce every existing AXI contract | gate contract | already covered — behavior-preserving refactor; the existing contracts are the net, and the canary layer proves per-fixture they still bite | a runner that swallows a failure or skips a block turns a contract into an always-pass, which its canary fixture turns red |
| 1 | runner cleans fixture dirs on failure too | gate contract | not TDD-able cheaply (would need an induced failing contract); asserted by inspection during review | leak is cosmetic, not correctness; review checks the runner's cleanup path |
| 2 | both flagged files under 400 lines; no new file over budget | bench structure | observed red: `bench structure` exits 1 today with FILE TOO LONG on both files | goes green only when the split and extraction actually land |
| 2 | guard-manifest check still fires before the CLI early return | gate canary layer | already covered — the minimal guard-manifest canary fixture (no CLI) must stay red on the inner run | if the split reorders the check behind the early return, that canary stops biting and the outer gate goes red |
| 3 | full allow/block verdict matrix unchanged after extraction | guard hook CLI | already covered — the runtime-git fragment asserts every verdict both ways | any classifier drift in the move flips a matrix row |
| 3 | `--describe` manifest still four keys, denies generated from DENY_LABELS; python3-missing degradation intact | guard hook CLI | already covered — AXI guard-manifest and python3-missing contracts | manifest or degradation drift fails those contracts against the extracted layout |
| 4 | analyzer file missing → `--describe` exits 0 with `manifest unavailable (analyzer missing)`; enforcement exits 2 BLOCKED | guard hook CLI | new contract, written first: fails against the pre-extraction hook (embedded analyzer can't be "missing", so the degradation line is absent) | asserts the fail-closed branch exists rather than an implicit allow |
| 4 (edge) | analyzer file present but empty → same fail-closed behavior | guard hook CLI | new contract, written first, red for the same reason | an empty program exits 0 printing nothing; without this branch the allow path reads that as "no verdict" and grants authority |
| 5 | syntax-broken analyzer turns the gate red | gate canary layer | new canary fixture with a broken `.py`, EXPECT the py_compile error substring | proves the parse check bites instead of rotting into an always-pass |
| 6 | structural debt zero | bench structure | same observed red as story 2 | definition of done for the cleanup |

### Edge inventory

Walked against the profile's hostile-input checklist per behavior:

- Paths with spaces — covered: the existing AXI space-path contracts run
  `bench guards` (hence hook `--describe` and its new path resolution) under a
  space-containing fixture.
- Required tool missing from PATH — covered: the existing python3-missing
  contract keeps running against the extracted hook.
- Absent vs empty file — covered: the two story-4 rows assert both, distinctly.
- Interrupted/partial state — covered (weakly, by review): runner cleanup row.
- Re-run idempotency — covered by existing link contracts: the hooks tree
  (now including the `.py`) ships through the same managed-tree copy that
  relink idempotence already exercises.
- cwd deeper than root — covered: the gate cds to root; existing subdir
  contracts unchanged.
- **Won't handle:** hook invoked through a symlink — harness configs and
  adapters reference the hook by its real path; a symlinked install is not a
  supported layout, and BASH_SOURCE-relative resolution under it fails closed,
  not open.
- **Won't handle:** hand-edited analyzer with CRLF / no trailing newline — the
  file is kit-shipped and managed by link; py_compile at gate time catches a
  corrupted copy in this repo.

## Out of scope

- **Migrate the five non-AXI gate fragments to the contract runner** — same
  runner, mechanical adoption across link/runtime/shift/line/docs fragments;
  a distinct hardening pass over files this build otherwise doesn't touch —
  ~25 edits, ~8 gate runs.
