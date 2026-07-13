# Oracle-bound gate verdict proof closure

Status: implemented

## Problem

FT78 has executable evidence for every R1–R21 requirement, both behavior-owned
canaries, and a green full gate. Its completeness tests deliberately compare each
executable proof registry with an independently literal expected-ID set, so removing
one registration makes the omission oracle red.

That independence conflicts with the project's one-source-per-fact standard when the
standard is read to forbid every repeated expectation. Deriving the expected IDs from
the executable registry makes the completeness test pass after the same omission it
exists to catch; removing the independent expectation therefore weakens the oracle.
FT78 cannot finish while the approved proof requirement and the project standard give
opposite instructions for the same test evidence.

## Solution

Keep the production semantics, public interfaces, executable registries, independent
expected-ID sets, and demonstrated completeness reds. Clarify the project standard
with one narrow exception: an independently authored test expectation is not
duplicated implementation knowledge when its independence is required for a named
omission or mutation to make the oracle red. The exception applies only to the
expectation-versus-implementation pair with recorded red evidence; it does not permit
duplicate production policy, parsers, fixture harnesses, registries, or derived
counts. Record the same current-state decision in one concise ADR so future reviewers
do not have to reconstruct why the exception exists.

The proof remains six ordered, separately green slices on one shared descendant: R21;
R1–R4; R5–R8; R9–R13; R14–R17; then R18–R20 plus integration. The existing
real-wrapper runtime contracts give every R1–R20 variant a literal assertion. One
private gate-engine seam supplies the fixed clock and scripted filesystem/lock
operations R21 needs. That seam remains limited to failures the shipped CLI cannot
induce deterministically and portably.

Treat the anchor's existing representative tests as positive controls, not as proof
of their whole quantified rows. For each new proof group, demonstrate that its
focused test goes red under the cheapest named sabotage of the anchor, then restore
the production implementation and make it green. Each slice receives a normal
path-scoped green commit. Only the final aggregate commit uses `--spec`, and FT78
becomes implemented only after every R1–R21 row and the full gate pass together.

## User stories

1. As a kit maintainer, I want a private deterministic-fault engine for R21, so that
   every durability and locking operation can be forced independently and its call
   order, result, and typed inspection can be asserted without platform-dependent
   permission tricks. Line: `gpt-5.6-sol` / high. This adds the missing seam inside
   the oracle owner, and the reviewer explicitly requires Sol for that judgment.
2. As a kit maintainer, I want every R1–R4 manifest, subject-mutation, environment-
   passlist, and resolver variant proven at the real wrapper, so that incomplete
   closure or oracle identity cannot authorize reuse. Line: `gpt-5.6-sol` / high.
   False-green identity evidence remains weakly observable until every named
   sabotage is shown to bite, so the reviewer explicitly routes this proof to Sol.
3. As a kit maintainer, I want every R5–R8 record, freshness, secret-safety, and
   hostile-invocation variant proven literally, so that relaxed parsing or unsafe
   rendering cannot inherit the anchor's representative green controls. Line:
   `gpt-5.6-sol` / high. Codec and rendering mistakes can agree with their own tests,
   so the reviewer explicitly routes the mutation proof to Sol.
4. As a kit maintainer, I want every R9–R13 durability, drift, cancellation,
   contention, and crash-recovery outcome proven, so that reordered persistence or
   false ownership cannot expose reusable evidence. Line: `gpt-5.6-sol` / high.
   Cross-process and durability evidence can stay falsely green without exact traces,
   so the reviewer explicitly routes this proof to Sol.
5. As a kit maintainer, I want every R14–R17 commit, shift, Stop, and repeated-partial-
   state outcome proven, so that action consumers cannot add a second authorization
   predicate or erase failed work. Line: `gpt-5.6-sol` / high. Consumer-specific
   fail-closed behavior is weakly observable until each operational result is forced,
   so the reviewer explicitly routes this proof to Sol.
6. As a kit maintainer, I want every R18–R20 typed state projected observationally
   through all three readers and the complete proof integrated with both canaries, so
   that copied parsers, reader repair, or format drift cannot survive the final gate.
   Line: `gpt-5.6-sol` / high. Cross-surface agreement can be self-consistently wrong,
   so the reviewer explicitly routes the mutation proof to Sol.
7. As a kit maintainer, I want independent test expectations distinguished from
   duplicated implementation knowledge when their independence makes a demonstrated
   omission oracle bite, so that one-source discipline cannot accidentally weaken the
   gate it protects, with the settled exception recorded for cold-session reviewers.
   Line: `gpt-5.6-sol` / high. This changes project guidance that steers every future
   review, so the leverage override requires Sol and high effort.

Stories 1–6 are green checkpoints on one descendant rooted at `8f9d2e4`; story 7 is
the terminal standards clarification that lets the completed proof be judged without
weakening its independent oracle. No story may replace the anchor with another FT78
worktree or broaden proof closure into a production redesign.

## Implementation decisions

### Preserve the anchor

`8f9d2e4` is the implementation baseline. Its concrete `Inspect` and `Execute`
interfaces, schema-1 verdict format, closed-subject policy, resolver precedence,
consumer routing, and behavior-owned canaries remain decided. Production edits are
allowed only when a newly executable R1–R21 proof exposes a real defect, and then only
at the existing `internal/gate` deep owner or its thin consumer. No public test port,
second cache parser, second writer, compatibility layer, or alternate implementation
is introduced.

### Delivery checkpoints

The six stories are separately green checkpoints on one shared branch, not six
independent feature definitions. Story 1 lands the private fault engine. Stories 2–5
close their contiguous ledger groups. Story 6 closes the projection group and runs
the aggregate integration proof. Every checkpoint runs its focused family and the
full gate before a normal path-scoped commit; a red checkpoint remains uncommitted.
The spec stays staged across checkpoints. Story 6 alone may make the finishing
`bench commit --spec oracle-bound-gate-verdicts` after all prior coverage rows are
re-verified green together.

The existing canaries remain controls:

- `gate-verdict-oracle-binding-bypassed` must still trigger
  `oracle-bound gate verdict contract failed`;
- `gate-verdict-invalidation-bypassed` must still trigger
  `fail-closed gate verdict persistence contract failed`; and
- both remain registered exactly once as `behavior-owned`, with the empty baseline
  still rejecting vacuous expectation text.

### R1–R20 proof ledger

The original acceptance rows 375–394 are named R1–R20 in order. At the anchor, none
is fully proven. The existing representative assertion in the second column stays as
a positive control; the third column is the exhaustive follow-on breadth.

| requirement | anchor control retained | variants this slice must add |
|---|---|---|
| R1 manifest closure | complete local manifest reuses once | absent; empty; malformed; wrong schema; remote; missing variable, path, and tool; escaped symlink; entry/byte limit; decoded control byte; each asserts reason and run count |
| R2 subject mutation | exact `BENCH_GATE`, declared environment, and ignored-file samples | gate script and interpreter; manifest; declared tool content, mode, and target; declared file and directory; `PATH`; repository root; auto-detected kind; ignored input rerun and refusal |
| R3 gate passlist | declared value retained; AMBIENT, HOME, and `BENCH_GATE` absent | exact `PATH`; `BENCH_KIT`; `BENCH_WRAPPER`; `CI`; `LANG`; `LC_ALL`; `LC_CTYPE`; one arbitrary inherited name; declared-absent variable |
| R4 resolution closure | pure precedence table and narrow real GateSh/BENCH_GATE/npm/no-gate controls | real pnpm, Python, and Cargo paths; launcher/tool mutation for GateSh (`gate.sh` plus its shebang chain), BenchGate (`bash`), Pnpm (`bash`, `pnpm`), Npm (`bash`, `npm`), Pyproject (`bash`, `mypy`, `pytest`, `ruff`), and Cargo (`bash`, `cargo`, `rustc`, `clippy-driver`); full real-wrapper precedence; no-gate exit 3/no record |
| R5 strict inspection | fresh green, ten-minute edge, five malformed samples, broad mode, interrupted and locked pending | absent; ready red; unavailable; zero byte; no-final-newline pair; trailing, duplicate, unknown, wrong type/schema/enum/hash/time; legacy; truncation; 16,384/16,385 bytes; symlink; directory; unreadable |
| R6 freshness | 9:59 fresh and exactly 10:00 stale | after 10:00; future; malformed time; fingerprint/freshness-policy mismatch; readers preserve bytes and metadata for every case |
| R7 secret-safe evidence | aggregate oracle identity exists | sentinel command, environment name/value, manifest path, input content, tool output, gate output, and unsafe control bytes change identity where applicable but never appear in cache or diagnostics |
| R8 hostile invocation | deep cwd and shipped symlink controls | spaces and glob characters in repository and declared paths; manifest without final newline; symlink chain and external target; missing global `bench`; executable-mode change; control-byte-safe output |
| R9 pre-run pending | real wrapper observes pending before gate marker | temp mode 0600 and ordered create/write/file-sync/rename/directory-sync; failure at each pre-run operation returns action 1 and leaves marker absent |
| R10 finalization | ordinary durable green and red | exact original red exit; green and red failures at write, file sync, rename, after-rename directory sync, and subject recheck; resulting bytes and typed state |
| R11 drift and cancellation | tracked-tree drift leaves pending/action 1 | command; manifest; environment; path; tool; launcher; auto-kind drift; cancellation kills the process group, returns interruption, and leaves pending |
| R12 execution ownership | second standalone execution fails behind one owner | standalone gate, commit, shift, and armed Stop each fail immediately without running, staging, or discarding work; a separate Git directory proceeds concurrently |
| R13 crash recovery | lock-free pending is classified interrupted | kill live owner; next run replaces pending and finishes; live-looking/dead-looking PID and young/old age never grant or deny ownership |
| R14 commit authorization | exact green reuse, stale-tree rerun, ordinary red and oracle mismatch preserve HEAD | absent, ready red, stale green, open-subject green, locked pending, interrupted pending, invalid, and unavailable inspections; lock-open, lock-acquire, pending-persistence, final-persistence, subject-build, subject-recheck/drift, cancellation, start-failure, and no-gate execution results; HEAD, index, named paths, and staged spec bytes remain unchanged |
| R15 shift preservation | ordinary red preserves the worktree and registration | lock, persistence, drift, and cancellation results preserve branch, intent, worktree registration, and uncommitted changes |
| R16 Stop ownership | green/red wrapper-owned bytes, no-gate/start failure, unarmed/active guards | real-wrapper invocation count exactly once; lock-open, lock-acquire, pending-persistence, final-persistence, subject-build, subject-recheck/drift, cancellation, start-failure, and no-gate results block; no second verdict write for each named result |
| R17 partial-state idempotency | pending survives a drift failure | repeat interrupted cancellation; invalid zero-byte, malformed, legacy, oversized, wrong-mode, symlink, and directory records; and every R21 operation fault; old green never returns; same-directory temporary files never become reusable evidence |
| R18 typed projections | locked/interrupted pending across all three readers | absent; reusable green; red; stale; invalid; unavailable across status, dashboard page, and roadmap context with literal state assertions |
| R19 read-only purity | pending bytes, mode, mtime unchanged and oracle not run | invalid and legacy as well as pending; no execution lock acquisition, rewrite, permission repair, or timestamp change on any reader |
| R20 presentation preservation | independent status and HTML controls | show-only-on-signal and severity for absent, reusable-green, red, stale, locked-pending, interrupted-pending, invalid, and unavailable; valid self-contained HTML; AXI TOON gate-cache fields `present`, `state`, `pending_status`, `status`, `cached_tree`, `work_tree`, `timestamp`, and `stale`; no raw hostile bytes |

### R21 deterministic-fault engine

The private engine remains inside `internal/gate`. Production `Inspect` and `Execute`
use one concrete implementation. Tests may supply a fixed clock and scripted
operations for exactly this ordered vocabulary:

1. lock open and nonblocking acquire;
2. temporary create, mode establishment, write, file sync, and close;
3. atomic rename, directory open, directory sync, and close; and
4. post-run subject rebuild.

The table forces one failure at each operation, plus the 16,384/16,385 byte boundary
and future-clock evaluation. Each case asserts the complete ordered trace, gate exit,
action exit, durable bytes when observable, and resulting `Inspection`. The seam does
not model repository behavior already deterministic through the real wrapper.

### Proof discipline

Because the anchor implementation already exists, R1–R20 are not honest
test-first-red rows. Before accepting each new proof group, temporarily apply its
named cheapest sabotage, run the focused test and record the red, restore the anchor
behavior, and run green. Representative sabotages are: accept any manifest for R1;
omit one digest component for R2; leak one ambient name for R3; bypass one resolver
kind for R4; relax one codec check for R5; make the freshness edge inclusive for R6;
render a sentinel for R7; normalize hostile paths incorrectly for R8; reorder one
durability call for R9/R10/R21; omit one drift component for R11; use per-command
locking for R12; trust PID/age for R13; rebuild authorization in commit for R14;
discard failed shift state for R15; write again in Stop for R16; retain old green for
R17; split one projection parser for R18; repair on read for R19; or bypass one public
format for R20.

Each story owns one test-side proof registry whose case IDs are the literal
`R<number>/<variant-slug>` expansion of its ledger cells above. The registry is the
only executable case inventory for that story: every entry carries its real driver
and literal assertions, and the group runner executes every registered entry. A
focused completeness test compares the registry with an independently literal
expected-ID set, rejects missing or duplicate IDs and nil drivers, and emits
`FT78 proof ledger completeness contract failed`. Before a checkpoint is accepted,
temporarily remove one expected case registration and record that completeness red,
then restore it. The final integration runs all six completeness tests together.
This bookkeeping signal supplements rather than replaces each requirement's named
behavior sabotage: a registered no-op still has to go red when its production fault
is planted.

The expected-ID set is an independent test oracle, not a second implementation
inventory. The project standard explicitly permits this expectation-versus-registry
repetition because removing the repetition makes the recorded missing-registration
red impossible. The exception is valid only where that mutation red is demonstrated
and named. Failure messages, registry runners, fixture construction, parsing,
production policy, and derived counts remain single-sourced; ordinary duplicated test
harnesses or parallel executable registries remain defects.

## Testing decisions

- R1–R4, R7–R20 attach to the existing real shipped-wrapper seam in throwaway Git
  repositories. Tests assert literal exits, output, run counts, cache bytes/metadata,
  HEAD/index, shift recovery state, and consumer formats as applicable.
- R5, R6, R9, R10, R17, and R21 use the private engine only for the exact fault or
  boundary that the CLI cannot force portably; their ordinary controls stay at the
  real-wrapper seam.
- Existing runtime families remain the harness: gate owns resolver, verdict,
  persistence, contention, crash, and drift; commit owns authorization atomicity;
  shift owns recovery preservation; Stop owns single invocation; dashboard/status/
  roadmap own projections. Do not create a second fixture framework.
- Every quantified variant above is a named table case or named subtest. A single
  broad permission failure, substring-only projection assertion, or source inspection
  does not count as proof.
- Focused mutation-red/green commands precede `.bench/gate.sh`. The full gate is the
  final oracle and FT78 remains staged until it exits zero on the complete descendant.
- The one-source standard names the narrow independent-test-expectation exception in
  project guidance and its settled decision record. A semantic pickup that treats
  these demonstrated completeness expectations as ordinary duplicated implementation
  knowledge is not considered resolved until those documents and the proof discipline
  agree.

### Seam diagram

Real CLI/action/projection seam:

    trigger: focused runtime contract or behavior-owned canary
        │
        ▼
    command + hostile repo state ──▶ [ shipped bench wrapper ──▶ gate owner ] ──▶ exit/cache/action
    verdict + repo state         ──▶ [ Inspect ──▶ typed consumers          ] ──▶ status/page/context
                                             ◀ tests attach here: literal black-box observations

Private deterministic-fault seam:

    trigger: internal gate table test
        │
        ▼
    fixed clock + scripted operation/fault ──▶ [ private gate engine ] ──▶ trace + Result + Inspection
                                                       ◀ tests attach here: one named fault per boundary

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | R21: a fixed clock and scripted lock/filesystem engine forces each enumerated operation independently and asserts full ordered trace, exits, durable floor, and typed inspection. | private deterministic-fault seam | Observed red at `8f9d2e4`: a search for fault, sync, rename, clock, byte-bound, or lock-acquisition tests under `internal/gate` exits 1; after the registry exists, removing one R21 registration must also emit `FT78 proof ledger completeness contract failed`. | Direct OS/time calls cannot force or distinguish the required faults; the trace rejects skipped, reordered, or ignored durability operations, while the registry red rejects a names-only partial table. |
| 2 | R1–R4: every manifest, subject-mutation, passlist, and resolver variant in the proof ledger has a named real-wrapper assertion while its anchor positive control remains green. | real CLI/action/projection seam | Not TDD-able at `8f9d2e4`: require the accept-any-manifest, partial-digest, ambient-leak, and partial-resolver sabotage reds plus a missing-registration `FT78 proof ledger completeness contract failed` red. | The four behavior sabotages reject false-green assertions; the exact registry rejects omission of a quantified variant. |
| 3 | R5–R8: every strict-record, freshness, secret-safe, and hostile-invocation variant in the proof ledger is asserted literally, using the private seam only for clock/byte cases. | real CLI seam plus private deterministic-fault seam | Not TDD-able at the anchor; require one codec, freshness, rendering, and path sabotage red plus a missing-registration completeness red before green. | Behavior mutation rejects relaxed parsing/rendering/path handling; the exact registry rejects a narrowed variant table. |
| 4 | R9–R13: every pending/final durability, drift/cancellation, contention, and crash-recovery variant in the proof ledger produces the exact action result and durable state. | real CLI seam plus private deterministic-fault seam | Not TDD-able at the anchor; require ordered-call, omitted-drift, per-command-lock, and PID/age sabotage reds plus a missing-registration completeness red. | Behavior mutation rejects incorrect ordering/ownership/recovery; the exact registry rejects an omitted fault or consumer case. |
| 5 | R14–R17: commit, shift, Stop, and repeated partial-state cases enumerate every listed state/result and preserve the exact action-owned state. | real CLI/action seam plus private deterministic-fault seam | Not TDD-able at the anchor; require second-predicate, cleanup-on-failure, second-writer, and old-green-resurrection sabotage reds plus a missing-registration completeness red. | Behavior mutation rejects consumer reinterpretation or cleanup; the exact registry rejects an omitted state/result pair. |
| 6 | R18–R20: all eight typed states are literal across status, the dashboard page, and roadmap context; every reader stays observational and each public format keeps its established contract. | real CLI/action/projection seam | Not TDD-able at the anchor; require split-parser, reader-repair, and format-bypass sabotage reds plus a missing-registration completeness red. | Behavior mutation rejects divergent parsers, writes, and formats; the exact registry rejects an omitted state/surface pair. |
| 6 | Integration: both existing behavior-owned canaries, their unique registry ownership, the empty-baseline vacuity check, the gate, commit, shift, Stop, status, dashboard-page, and roadmap-context focused families, all six completeness tests, and the full gate remain green together. | real CLI seam through runtime contracts and canary registry | Already covered at the anchor for the two canary needles; stories 1–6 must additionally retain their named mutation and missing-registration reds plus focused greens. | This preserves the biting independent invariants while preventing separately green checkpoints from passing through wrong routing or ambient expectation text. |
| 7 | Project guidance and its settled decision record distinguish a demonstrated independent test expectation from duplicated implementation knowledge without creating a general duplication escape hatch. | working agreement and decision record plus existing completeness mutation seam | Observed semantic red: the FT78 review pickup remains actionable while the working agreement forbids the independently literal expected-ID sets that the approved proof discipline requires. | The pickup cannot close if the guidance still classifies the required omission oracle as a defect; the existing missing-registration mutation reds ensure the exception preserves rather than weakens enforcement. |

**Degenerate-implementation check.** The cheapest wrong proof adds test names without
variant assertions, loops over cases without observing per-case reason/state, uses one
permission failure for every durability error, or scripts traces that production never
crosses. Named variant tables, literal black-box observations, mutation-red evidence,
and R21's production-backed ordered trace make each of those implementations red.

### Edge inventory

- **Error paths:** every R1–R21 failure named in the ledger is a coverage case; no
  generic error assertion substitutes for exact action exit and durable state.
- **Empty and absent:** R1, R5, R18, and R19 distinguish absent, empty, and valid
  empty collections or output.
- **Boundary values:** R1/R5/R21 cover 16,384/16,385 bytes; R6 covers before/exactly
  at/after ten minutes; R13 covers live/dead-looking PID and young/old age; gate exits
  include 0, 1, 2, 3, 127, interruption, and one larger arbitrary nonzero.
- **Malformed input:** R1, R5, and R8 enumerate schema, framing, UTF-8/control,
  symlink/type/mode, and hostile path cases.
- **Interrupted or partial state:** R9–R17 land every pre-run/finalization fault,
  drift, cancellation, crash, repeated run, and old-green posture as a row above.
- **Re-run idempotency:** R1, R13, and R17 assert run counts, recovery, temporary-file
  cleanup, and non-resurrection.
- **Hostile environment:** R2, R3, R7, and R8 cover PATH, inherited variables,
  secrets, modes, symlinks, missing tools, deep cwd, spaces, and globs.
- **Over-broad standards exception:** story 7 permits only independently authored
  expectations with a named mutation red; duplicated executable registries, fixture
  harnesses, parsers, policies, or derived counts remain defects.
- **Remote equivalence — Won't handle:** remote state remains non-reusable because
  proving service equivalence is a separate trust protocol and the in-scope caller
  still executes it every time.
- **Malicious local repository owner — Won't handle:** the local threat model still
  excludes an owner who replaces the binary, Git directory, lock substrate, and cache.
- **Destructive worktree classification — Won't handle:** ownership and cleanup of
  foreign, identity-mismatched, reused, primary, ignored-residual, or dirty-nested
  worktrees belongs to the separate worktree lifecycle contract; the in-scope gate
  callers still exercise FT78 in ordinary registered worktrees.

## Out of scope

- **Changing FT78 production semantics or public interfaces without a failing R1–R21
  proof** — a separate redesign decision, not proof closure. Estimated later cost:
  `12 edits, 6 gate runs`.
- **Recovering or integrating the retained FT78 experiment** — its committed base is
  byte-identical to the selected anchor, while its uncommitted private-engine and spec
  rewrite is a separate candidate this closure does not absorb. Estimated later cost:
  `10 edits, 5 gate runs`.
- **Remote reusable attestations** — a separate trust/protocol capability; remote
  gates continue to run every time. Estimated later cost: `10 edits, 6 gate runs`.
- **FT82 release-preflight evidence and FT88 platform-wide environment minimization**
  — separate roadmap capabilities. Estimated later cost: `22 edits, 11 gate runs`.
