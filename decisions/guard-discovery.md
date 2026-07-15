# Static, bounded guard discovery (FT80)

## Destination

`bench guards` and the SessionStart inspection describe every guard from static
managed metadata without executing any repository shell file, under one
aggregate deadline — closing `RR:R-07` / `RC:H-02` with a sentinel contract
proving non-execution.

## #1: Where does a guard's manifest come from without execution?

Type: Grill

### Question
Descriptions currently come from executing `bash <hook> --describe`. What
static source replaces it — a registry file, kit-embedded metadata, or the
script itself?

### Answer
Static comment-header lines at the top of each hook script (`# name: …`,
`# boundary: …`, `# denies: …`, `# why: …`), parsed as text. The `--describe`
protocol is deleted everywhere. The script stays the single source per fact;
reading text can never execute; project-added hooks self-describe the same way.

## #2: Does anything in the guards path still execute?

Blocked by: #1
Type: Grill

### Question
With headers parsed statically, is execution still required anywhere — in
particular for the managed git pre-push hook, which was the one file already
guarded against foreign execution?

### Answer
No — zero execution everywhere. The pre-push file (bench-written, so it carries
the same comment header) is parsed statically; the existing bench marker still
distinguishes managed from foreign for row wording. FT80's "exact managed
allowlist may execute" collapses to the empty set: `bench guards` spawns no
process.

## #3: Is there a managed-vs-unknown distinction for hooks-dir scripts?

Blocked by: #1
Type: Grill

### Question
Should unknown additions to `.bench/hooks/` be classified against a managed
set, or reported uniformly?

### Answer
Uniform. Every `.bench/hooks/*.sh` gets a row from its parsed header, or a
definitive `no manifest` when the header is absent or incomplete; the derived
`wired` cell already says whether it can fire. No managed-allowlist for the
hooks dir — it would require a registry, reintroducing a second source. Only
pre-push keeps managed/unmanaged wording, via the existing marker.

## #4: Where does the aggregate deadline live?

Blocked by: #2
Type: Grill

### Question
SessionStart's inspection is three calls (`resume-clean`, `status`,
`guards --brief`) orchestrated by `session-start.sh`, and shell `timeout` is
not POSIX. Where does FT80's single deadline attach?

### Answer
One Go plumbing subcommand runs all three under a single context deadline
(~10s; spec pins the value). `session-start.sh` reduces to resolve-and-exec
plus the CLI-location line. On a deadline trip the inspection warns and exits
0 — the session never blocks (existing SessionStart posture).

## #5: What form do the sentinel and the bite-proof take?

Blocked by: #2
Type: Grill

### Question
Closure requires an unwired sentinel that never executes through `bench
guards` or SessionStart, plus proof the non-execution rule bites. Where does
the sentinel live, and is the proof a `tests/canary/` fixture?

### Answer
Sentinel is a test fixture only: a contract test builds a temp repo whose
hooks dir holds a sentinel script that writes an evidence file if executed,
runs `bench guards` and the session-start inspection, and asserts no evidence.
Not a tracked hook in `.bench/hooks/` — the kit tree ships to every linked
repo, so it would propagate as noise. The bite-proof is a recorded mutation
demo (mutate guards back to exec during the TDD red step, observe the sentinel
test go red, record it), not a canary fixture — canary fixtures grade repo
trees, but non-execution is binary behavior a fixture tree cannot violate.

## Not yet specified

- Exact header grammar (key order, tolerance for leading `#!` and blank lines)
  and the deadline value — spec-level detail.
- Name of the Go inspection subcommand.

## Out of scope

- Gate-time execution of trusted project code (the gate itself remains
  intentionally trusted; only the describe/inspection surface goes static).
- Harness-side hook execution at real boundaries (PreToolUse etc.) — the
  harness fires wired hooks regardless of bench; FT80 owns only bench's own
  inspection surface.
- FT81+ release/runtime work.

## Handoff

1. **Module boundaries.** `internal/guards` becomes a static header parser +
   row assembler (no process spawn). Each `.bench/hooks/*.sh` and the
   `internal/adopt/prepush.sh` template carry comment-header manifests and lose
   their `--describe` blocks. A new Go plumbing subcommand owns the SessionStart
   inspection (resume-clean + status + guards --brief) under one context
   deadline; `.bench/hooks/session-start.sh` becomes a thin resolve-and-exec
   wrapper. The conformance manifest check (currently executing `--describe`)
   migrates to static header validation, and the existing
   `tests/canary/package-core-guard/guard-describe-*` fixture family migrates
   to the header convention.
2. **Contracts.** `bench guards` / `--brief`: same table and brief shapes; rows
   from parsed headers; `no manifest` for absent/incomplete headers;
   informational hooks (`denies: nothing (informational)`) still excluded;
   spawns no process. Pre-push row: marker present → parse header (managed);
   no marker → `unmanaged (no manifest)`; absent → `not installed`; never
   executed. Inspection subcommand: runs the three phases under one ~10s
   aggregate deadline; on trip, emits a warning and exits 0; never blocks a
   session.
3. **Deep vs thin.** The header parser is the deep unit (hides the manifest
   grammar; both guards and conformance read through it — one source for the
   grammar). The inspection subcommand and `session-start.sh` are thin
   orchestration/pass-throughs with no seam of their own.
4. **Black-box assertables.** Guards stdout (table and brief) against fixture
   hook trees; sentinel evidence-file absence after running guards and the
   inspection; pre-push row wording per marker state; inspection exit 0 and
   warning on deadline trip.
5. **Gate attachment.** Runtime contract tests against the built `dist/bench`
   own the non-execution sentinel and output contracts; the conformance family
   owns header validity of the kit's shipped hooks. The mutation bite-proof is
   demonstrated at TDD-red and recorded — the gate cannot hold it continuously.
6. **Hostile-input owners.** Hostile header content (huge files, missing
   newline, binary bytes) → the header parser. Foreign pre-push content → the
   marker check + static read. Non-`.sh` entries and subdirectories → existing
   hooks-dir skip logic. Unparseable harness configs → existing
   `wiredHarnesses` not-wired posture.
7. **Uncertainty flags.** Testing the aggregate-deadline trip black-box may
   need an injectable slowness knob or a unit-level context seam — the spec
   writer picks; nothing in the map settles it.
8. **Rejected alternatives.** Separate registry file (second source, drifts);
   kit-embedded metadata (breaks project-added hooks); any execution allowlist
   (empty set won); tracked sentinel in `.bench/hooks/` (propagates to linked
   repos); `tests/canary/` fixture as bite-proof (fixture trees cannot violate
   binary behavior); shell `timeout` watchdog (not POSIX).
9. **Domain watch-outs.** The `.bench` tree ships to every linked repo — any
   file added under `.bench/hooks/` propagates on link. The header convention
   is a contract change for project-added hooks that answered `--describe`;
   they read as `no manifest` until they carry headers. The pre-push hook lives
   in `.git/hooks`, outside the repo tree — its manifest travels in the
   bench-written file, not the working tree.

Dependency order: n/a — single spec.
