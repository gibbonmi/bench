# Per-path structure budgets

## Problem
`bench structure` has two global caps (`BENCH_MAX_LINES`, `BENCH_MAX_DIR_FILES`).
A file that is legitimately one deep module can exceed the line cap; today the
only outs are raising the *global* cap (weakening the check everywhere) or
parking a permanent learnings entry. `craft-seams` promises "no per-file budget"
precisely because there was no reviewer-owned mechanism for one.

## Solution
A committed, reviewer-owned exceptions file: `.bench/structure.budgets`. One
override per line — `<path> <budget>`; a trailing `/` on the path makes it a
directory file-count budget, otherwise it is a file line budget. `bench
structure` (and the shift refactor trigger, which routes through the same
check) uses the override where present, the global cap elsewhere. Granting an
override is the reviewer's call; agents propose, never edit.

## User stories
1. As a reviewer, I want to grant a named file a higher line budget in
   `.bench/structure.budgets`, so one deep module stops paging the structure
   signal without weakening the cap for the whole repo.
   Line: claude-fable-5 (inline, session model) / low. Plumbing at the existing
   structure_check seam, fully gate-observable.
2. As a reviewer, I want directory overrides (`path/ <max-files>`) with the
   same semantics, so a legitimately crowded directory can be granted too.
   Line: claude-fable-5 (inline, session model) / low. Same seam, same shape.
3. As an agent, I want `craft-seams` to route me to *proposing* a
   `.bench/structure.budgets` line instead of "there is no per-file budget",
   so the escape hatch is the documented path and stays reviewer-owned.
   Line: claude-fable-5 (inline, session model) / medium. Guidance prose —
   leverage override applies.

## Implementation decisions
- Lookup lives inside `structure_check` so every consumer (CLI `structure`,
  `structure_touched_since`, the shift refactor phase, `bench status`) inherits
  it — one seam, no new entry points.
- File format: `<path> <integer>` per line; `#` comments and blank lines
  ignored; paths are repo-relative, matched exactly (no globs — an override is
  a named, deliberate grant, and glob grants re-open the hole the global cap
  closes). Malformed lines (missing/non-integer budget) are ignored with a
  one-line stderr warning naming the line — a typo must not silently drop the
  global cap or crash the dashboard.
- An override *lowers* as well as raises (the value replaces the cap for that
  path); the report line names the budget actually applied so a granted file
  reads differently from a default one.
- No gate change in this repo: `bench structure` is not part of benchkit's
  gate; the contract tests live with the other structure runtime contracts.

## Testing decisions
- Good test = plant a throwaway repo, run `bash bin/bench.sh structure`,
  assert exit code and message — the existing structure contract pattern.
- Seam: the `bench` CLI subcommand surface (`structure`), prior art: the
  "structure shell-file contract" in the runtime contracts.
- Gate command: `.bench/gate.sh`.

### Seam diagram

    trigger: bench structure / bench status / shift refactor check
        │
        ▼
    tracked source files ──▶ [ structure_check          ] ──▶ violations + exit code
    global caps (env)    ──▶ [  + .bench/structure.budgets ]
    overrides file       ──▶ [    per-path lookup        ]
                                  ◀ tests attach here: throwaway repo with planted
                                    files + budgets file, assert verdict/message

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | 401-line file + `file 500` override → structure green | CLI contract | contract fails until lookup exists | override not consulted → FILE TOO LONG persists |
| 1 | same file, no override → still red | CLI contract | already covered (shell-file contract) | proves the global cap survives the feature |
| 2 | 13-file dir + `dir/ 20` override → green | CLI contract | contract fails until dir lookup exists | dir overrides are a distinct code path |
| edge of 1 | override *below* global (e.g. `file 100` on a 200-line file) → red | CLI contract | contract fails until value replaces (not maxes) the cap | "override = replace" is the decided semantic |
| edge of 1 | malformed line (`file abc`) → warning, global cap applies, no crash | CLI contract | contract fails while a malformed line crashes or silently grants | hand-edited file is the hostile input here |
| 3 | craft-seams names the budgets file as the escape hatch | gate anchor | anchor red until prose updated | the old "no per-file budget" sentence contradicts the feature |

### Edge inventory
- absent vs present-but-empty budgets file → both mean "global caps only"; folded
  into story-1's no-override row (absent) — **Won't handle:** asserting empty
  separately; same code path as absent by construction (lookup finds nothing).
- malformed input → coverage row (warning + fallback).
- boundary → override value exactly equal to the line count is green (`>` stays
  the comparison, as today); exercised by the 401/500 row's arithmetic.
- paths with spaces → **Won't handle:** the lookup splits on whitespace;
  space-containing paths can't be granted (they can still pass via the global
  cap). One-line limitation noted in the file's header comment.
- missing trailing newline on the last budgets line → covered implicitly: the
  read loop must process the final unterminated line (assert in the malformed
  contract fixture).
- re-run idempotency → read-only feature; nothing to assert.

## Out of scope
- Gate-enforced "agents may not edit .bench/structure.budgets" — a separate
  enforcement capability (guard hook over reviewer-owned files, applies to more
  than this file); ~10 edits, ~6 gate runs.
