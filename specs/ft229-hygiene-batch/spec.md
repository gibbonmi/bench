# FT229 — hygiene batch: residue, stale prose, and small CLI sharp edges

Status: implemented

Decision source: named reviewed artifact — `roadmap/FT229.md`, whose evidence is the 2026-08 capability audit's reconciliation ledger (`docs/audits/2026-08-bench-capability/results-fable-high/reconciliation-ledger.md`, rows L-17 through L-20, L-29, L-30) and action item A9 (`action-items.yaml`). Two forks were closed by the reviewer on 2026-08-19: the existing tickets-only residue is deleted in this build, and `bench preflight` gains `--source-tip` rather than the review phase losing the word.

Verification log: 2 iteration(s) to accept — one `fable`/high round returned five blocking findings, all folded: the invented zero-count conformance check was dropped and story 9 took a named `Not covered` exception, `internal/sessioninspect` joined the fences, the status-count ticket gained its blocking edge on the close-step ticket, the two rebuild-command rows were pinned to the exact `scripts/go-build.sh` invocation, and the log retention count was pinned at 20. Folded partials: the residue count corrected from 30 to 37, Group B rerouted to `fable`/high under `craft-line`'s leverage override, the two `internal/systemtest` writers serialized, a row added for the non-fatal pruning failure, a vacuous carve-out removed, and `internal/structure` added to the diagnostics ticket's `Writes:`.

## Problem

Seven small defects each cost more to ship alone than they cost to fix, so each
one waits. Together they degrade the parts of Bench a cold session trusts first.

A light-path change writes a ticket folder and never removes it, so `specs/`
holds 37 tickets-only folders that no command owns and no signal counts.
`bench structure` opens with four stale accept rows before its real findings.
The review phase's prose describes what the gate runs from a retired shape.
Bare `bench outline` spends 24,533 bytes on 200 of 5,335 symbols chosen by no
ranking rule. The destructive-git guard allows `git stash drop`, `git stash
clear`, `git filter-branch`, and `git rm -rf .`, and when the core is missing
its degraded rim refuses any command whose envelope text contains the three
letters `git` anywhere. SessionStart with no core exits 127 in silence. `.logs`
accumulates one JSONL record per gate run with no pruning owner. The landing
guard treats an ignored `dist/bench` as removable residue, and removing it in a
self-hosting checkout disables the binary that answers `bench`, the git guard
that binary backs, and the `BENCH_RUN_BINARY` the gate requires.

`bench preflight` takes `--base` but not `--source-tip`, while the review phase
asks the reviewer for a frozen pair.

## Solution

One batched pass under one gate. The light-path landing commit consumes its own
ticket folder, `bench status` counts what is left, and the 37 existing folders
go. The stale accept rows and the retired-shape prose are corrected and pinned.
Bare `bench outline` becomes a summary and rows require a path or `--full`. The
guard denies the four missing spellings, and its degraded rim decides from the
envelope's command field rather than the raw envelope text. SessionStart on a
missing core prints the command that builds one, and the residue guard warns
before removing the binary currently answering `bench`. Gate run logs prune by
count. `bench preflight` accepts `--source-tip` and verifies it against the tip
it derives.

## User stories

### Group A — the light path closes its own ticket

Line: `opus` / medium. The close step deletes a directory on a green landing, so
its refusal predicate needs care, but the surrounding lifecycle already exists.

1. As a builder, I want the green landing commit of a light-path change to delete its tickets-only spec folder, so that a landed ticket leaves no residue behind it.
2. As a builder, I want the close step to act only on a folder carrying no `spec.md`, so that a spec-backed folder is never deleted by the light-path route.
3. As a builder, I want the close step to refuse a slug that names no folder, so that a typo reports instead of landing silently.
4. As a builder, I want the folder deleted only on a green landing, so that a red gate leaves the ticket in place.
5. As a builder, I want the deletion to be part of the landing commit, so that the tree never holds a landed change beside its unconsumed ticket.
6. As a reviewer, I want `bench status` to count tickets-only folders, so that the residue is a visible signal rather than invisible debt.
7. As a reviewer, I want that row to name the command that closes one, so that the row is actionable without a lookup.
8. As a reviewer, I want the row absent at a count of zero, so that the five-row budget is never spent on a clean signal.
9. As a reviewer, I want every existing tickets-only folder deleted in this build — 37 of them at HEAD — so that the new count starts at zero.

### Group B — the diagnostics tell the truth

Line: `fable` / high. `craft-line`'s leverage override routes kit-guidance prose
to the top tier: these passages ship to every linked repo, and the anchors pin
whatever wording lands.

10. As a reviewer, I want `bench structure` to print no stale accept row, so that its first four lines are findings rather than bookkeeping.
11. As a reviewer, I want every surviving accept grant to keep its exact reason, so that dropping the stale rows does not silently drop a live decision.
12. As an agent, I want the review phase's prose to name what the gate actually runs, so that I do not plan around a phase that does not exist.
13. As an agent, I want `.bench/BENCH-reference.md` to describe the current conformance shape, so that the reference and the phase table agree.
14. As a reviewer, I want an anchor pinning each corrected passage, so that the prose cannot drift back without a red gate.

### Group C — the bare read surface is a summary

Line: `opus` / medium. The change is small, but it moves an AXI-conformant
output contract, so the empty state and the row-bearing forms must stay exact.

15. As an agent, I want bare `bench outline` to emit meta and per-directory counts, so that one probe costs a screen rather than 24 KB.
16. As an agent, I want `bench outline <path>` to keep emitting rows, so that the targeted form still locates a seam.
17. As an agent, I want `bench outline --full` to emit rows repository-wide, so that the old behavior stays reachable when I want it.
18. As an agent, I want the bare form's size bounded by the directory count rather than a truncation limit, so that its output is complete rather than a silent prefix.
19. As an agent, I want a repository with no scannable source to yield a definitive empty state, so that empty and broken stay distinguishable.

### Group D — the guard denies what it claims

Line: `opus` / high. This is the enforcement boundary, a wrong verdict is silent,
and the parser meets deliberately awkward spellings.

20. As a reviewer, I want `git stash drop` blocked, so that discarding a stashed change is not an agent's call.
21. As a reviewer, I want `git stash clear` blocked, so that discarding every stashed change is not an agent's call.
22. As a reviewer, I want `git filter-branch` blocked, so that a history rewrite stays mine.
23. As a reviewer, I want `git rm -rf` blocked, so that a recursive forced removal cannot pass as ordinary file work.
24. As a builder, I want `git stash`, `git stash list`, and `git stash pop` still allowed, so that the safe half of the verb stays usable.
25. As a builder, I want `git rm <path>` still allowed, so that an ordinary tracked-file removal does not need the reviewer.
26. As a builder, I want a sought token supplied only as a flag value to reach the same verdict as its positional spelling, so that `git stash -m drop` cannot resolve to an allow.
27. As a reviewer, I want the degraded rim to block a destructive git command when no core is reachable, so that the fail-closed posture survives a missing binary.
28. As a builder, I want the degraded rim to allow a command that only mentions `git` in a path or an argument, so that reading a file under `.github/` is not refused.
29. As a builder, I want the degraded rim to decide from the envelope's command field, so that a match in an unrelated envelope field cannot refuse a shell.
30. As a reviewer, I want the degraded rim to refuse when it cannot read a command field at all, so that an unparseable envelope fails closed rather than open.
31. As a reviewer, I want the degraded rim to stay one honest-mistake layer, so that its narrowing is not read as evasion resistance.

### Group E — a cold session recovers

Line: `opus` / medium. Both faces are recovery paths that only fire when
something is already broken, so their evidence has to be produced deliberately.

32. As an agent, I want SessionStart with no core to print `bash scripts/go-build.sh` with its arguments, so that a cold session recovers without reading the handoff or guessing at `go build`.
33. As an agent, I want SessionStart to keep never blocking a session, so that a broken core costs a hint rather than a session.
34. As an agent, I want SessionStart outside a repository to stay silent, so that the hint does not become ambient noise.
35. As a builder, I want the residue guard to warn before removing a `dist/bench` that is currently answering `bench`, so that clearing residue does not disable the CLI, the guard, and the gate at once.
36. As a builder, I want that warning to name `bash scripts/go-build.sh` with its arguments, so that recovery does not need a search and does not reach for plain `go build`.
37. As a builder, I want a `dist/bench` that is not the resolving binary removed without the warning, so that ordinary residue stays ordinary.

### Group F — runtime records have a pruning owner

Line: `opus` / medium. The profile routes gate logic to mid, and this pruner
runs inside the oracle's own package, where a wrong retention destroys evidence.

38. As a builder, I want gate run logs pruned by count with the newest retained, so that `.logs` stays bounded without a chore.
39. As a builder, I want pruning to happen inside a gate run, so that no separate command has to be remembered.
40. As a builder, I want the run currently being written never pruned, so that a gate never truncates its own evidence.
41. As a builder, I want a file in `.logs` that the gate did not write left alone, so that pruning owns only its own records.

### Group G — preflight's grammar matches the phase

Line: `opus` / medium. One flag on an existing parser, whose value is a verified
pin rather than a recorded one.

42. As a reviewer, I want `bench preflight review <slug> --source-tip <commit>` accepted, so that the pair the phase froze is the pair preflight checks.
43. As a reviewer, I want a `--source-tip` that disagrees with the derived tip to red, so that a stale pin is caught rather than reported.
44. As a reviewer, I want a `--source-tip` that does not resolve to report distinctly from a mismatch, so that a typo and a drift are different diagnoses.
45. As a reviewer, I want omitting `--source-tip` to keep today's behavior, so that the flag is an addition rather than a new requirement.
46. As a reviewer, I want `bench preflight build` to accept the same flag, so that the two phase entries share one grammar.

## Implementation decisions

**The close step lives on `bench commit`, not on a new verb.** `bench commit
--spec <slug>` already owns the spec-state transition that a green landing
publishes. A tickets-only folder takes the same flag and the same green
precondition: on a folder holding no `spec.md`, the flag's effect is deletion of
the folder rather than a status flip. One flag, one landing, one source for
"what a landing does to its spec directory". A folder holding a `spec.md` keeps
today's flip exactly.

**The tickets-only count joins the existing housekeeping band in `bench
status`.** `internal/status` already ranks merged-spec retirement at severity 8
and orphaned pickups at 9. The tickets-only count enters below them and shares
the band's show-only-on-signal rule, so a zero count prints nothing and the
five-row budget is untouched.

**The degraded rim parses the envelope's command field in the shim.** The rim
runs precisely when the core is unreachable, so it cannot delegate. It extracts
`tool_input.command` from the PreToolUse envelope and tests whether the command
invokes `git` — the same one-level wrapper depth `internal/gitguard` documents.
An envelope with no readable command field refuses, keeping the fail-closed
posture. The rim stays an honest-mistake layer; the pre-push hook and pooled
isolation remain the backstops for a misaligned agent, and this spec does not
claim otherwise.

**The four new deny spellings go in `internal/gitguard`'s deny table**, beside
the classes already there, so the shim gains nothing and the core stays the one
classifier. `git rm` denies on the recursive-and-forced combination, matching
how `clean` already reads its short-flag clusters, so an ordinary `git rm
<path>` is untouched.

**Bare `bench outline` changes shape, not contract.** It stays AXI-conformant
TOON on stdout: a meta line and one row per top-level directory carrying its
symbol count. `--full` restores repository-wide symbol rows. A path argument is
unchanged. The 200-row cap disappears with the bare form that needed it, so no
form emits a silent prefix.

**Log pruning is a gate-run side effect retaining the newest 20 records.** The gate
already creates `.logs` and writes one record per run, so it prunes there —
the newest 20 retained, the in-flight run excluded, and only files matching the
gate's own `gate-<timestamp>-<pid>.jsonl` name shape considered. A pruning
failure writes one stderr warning and never changes the gate's verdict or exit
code. Twenty is a reviewer-owned constant, sized to hold a full build session's
runs.

**The `dist/bench` warning attaches to the guard that would remove it**, not to
`bench doctor`. Doctor reports; the residue path acts, and the warning must
reach the caller in the moment the removal is proposed. Both this warning and
the SessionStart hint name `bash scripts/go-build.sh` with its arguments, which
the profile requires over plain `go build` because that script stamps the
version the upgrade and gate contracts read. The predicate is
identity, not path: the guard compares the candidate against the binary the
wrapper's own resolution selects, so a `dist/bench` in an unrelated checkout is
ordinary residue.

**`--source-tip` is verified, not recorded.** Preflight already derives the tip;
the flag supplies the reviewer's frozen value and preflight reds when the two
disagree. An unresolvable value is a grammar-level error, matching how `--base`
already separates "cannot resolve" from a predicate red.

**The inert Codex invocation key is out of scope — it already shipped.** FT228
(`0b956ec3`) removed `disable-model-invocation` from all thirteen phase adapters
and (`bbff89e1`) added the Claude-side parity grading. L-30 is closed and this
spec does not restate it.

## Testing decisions

A good test here drives the real surface: the shim as a subprocess with no
reachable core, the guard's classifier through its envelope, the gate through a
run that writes and prunes, and `bench status` over a real `specs/` tree. The
gate seam is the ordinary test phase, which carries the conformance registry, so
prose anchors and command grammar red inside the same run as the Go units.

### Seam diagram — the light-path close

    trigger: bench commit --spec <slug> on a light-path change
        │
        ▼
    slug + tree  ──▶  [ internal/commit → internal/landing ]  ──▶  green landing commit
                          ◀ tests attach here: build a repo whose specs/<slug>/
                            holds tickets and no spec.md, land, assert the folder
                            is gone from the commit and present after a red gate

### Seam diagram — the guard's two layers

    trigger: agent Bash call (PreToolUse envelope)
        │
        ├─ core reachable ──▶ [ internal/gitguard classify ]  ──▶ allow / BLOCKED
        │                         ◀ tests attach here: envelope table per spelling
        │
        └─ core missing   ──▶ [ .bench/hooks/block-dangerous-git.sh rim ]  ──▶ exit 0 / exit 2
                                  ◀ tests attach here: internal/systemtest runs the
                                    shim with the wrapper unreachable and PATH scrubbed

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| H01 | 1, 5 | a green landing on a tickets-only slug removes the folder in the same commit | internal/commit | a close step that runs outside the commit leaves the folder in the tree the commit published |
| H02 | 2 | a slug whose folder holds spec.md takes the status flip and keeps its folder | internal/commit | the cheapest wrong build deletes on every slug, which destroys a staged spec |
| H03 | 3 | a slug naming no folder returns a structured error and lands nothing | internal/commit | a build that treats an absent folder as already-closed lands a typo silently |
| H04 | 4 | a red gate leaves the tickets-only folder present | internal/commit | deletion before the verdict discards a ticket whose work never landed |
| H05 | 6, 7 | a specs tree with tickets-only folders renders one row carrying the count and the closing command | internal/status | a row without the command sends the reader back to the docs |
| H06 | 8 | a specs tree with no tickets-only folder renders no such row | internal/status | an always-on row spends the five-row budget on a clean signal |
| H08 | 10 | bench structure reports no stale accept row | internal/structure | a build that deletes the wrong rows leaves the stale ones printing |
| H09 | 11 | every surviving accept grant resolves to a scanned source file and keeps its reason text | internal/structure | dropping four rows by hand can take a live grant with them |
| H10 | 12, 13 | the review-phase and reference passages match the phase table the gate runs | internal/conformance | prose corrected without a pin drifts back at the next edit |
| H11 | 14 | removing a corrected passage reds the anchors check | internal/anchors | an anchor that does not bite is decoration |
| H12 | 15 | bare bench outline emits meta and one row per top-level directory with its symbol count | internal/outline | the cheapest wrong build truncates the symbol rows instead of changing shape |
| H13 | 16 | bench outline with a path emits symbol rows for that path | internal/outline | a shape change applied to every form removes the targeted probe |
| H14 | 17 | bench outline --full emits symbol rows repository-wide | internal/outline | without the escape hatch the old capability is gone rather than moved |
| H15 | 18 | the bare form's row count equals the top-level directory count with no cap applied | internal/outline | a surviving 200-row cap makes the summary a silent prefix |
| H16 | 19 | a tree with no scannable source yields the definitive empty state | internal/outline | an empty summary and a broken scan are indistinguishable otherwise |
| H17 | 20, 21, 22, 23 | each of the four spellings returns a block verdict with its class label | internal/gitguard | the audit reproduced all four as allowed at HEAD |
| H18 | 24, 25 | the safe stash forms and an ordinary git rm path return an allow verdict | internal/gitguard | a rule written on the verb alone blocks the working half of the command |
| H19 | 26 | a sought token supplied only as a flag value reaches the destructive verdict | internal/gitguard | a positional resolver that skips flags but not their values lands an allow on the mutation |
| H20 | 27 | the shim with no reachable wrapper and no PATH bench refuses a destructive git command | internal/systemtest | a narrowed rim that stops refusing turns the degraded path fail-open |
| H21 | 28, 29 | the shim allows a command whose git text appears only in a path or an unrelated envelope field | internal/systemtest | the raw-substring rim refuses ordinary reads, which is the defect |
| H22 | 30 | the shim refuses an envelope carrying no readable command field | internal/systemtest | a parser that treats an unreadable envelope as empty grants everything |
| H23 | 32, 33 | a session start with an unreachable core prints the build command and exits zero | internal/systemtest | a hint that blocks is worse than the silence it replaced |
| H24 | 34 | a session start outside a repository prints nothing | internal/systemtest | the hint fired unconditionally becomes ambient noise |
| H25 | 35, 36 | the residue guard on a dist/bench matching the resolved binary emits the warning naming the rebuild command | internal/worktree | the guard removed it silently, which is the reported incident |
| H26 | 37 | the residue guard on a dist/bench that is not the resolved binary removes it without the warning | internal/worktree | a path-shaped predicate warns on every checkout and trains the warning away |
| H27 | 38, 39 | a gate run over a log directory holding more than 20 records leaves exactly the newest 20 | internal/gate | an unpinned retention lets N=1 pass while destroying the history the pruning exists to bound |
| H28 | 40 | the record the current run is writing survives its own pruning | internal/gate | a pruner counting before the write truncates the evidence it is producing |
| H29 | 41 | a file in .logs outside the gate's name shape survives pruning | internal/gate | a directory-wide sweep removes records the gate does not own |
| H34 | 41 | a pruning failure writes one stderr warning and leaves the gate's verdict and exit code unchanged | internal/gate | a pruner that fails the gate turns housekeeping into a false red |
| H30 | 42, 45, 46 | review and build accept --source-tip and behave unchanged when it is omitted | internal/preflight | a required flag breaks every existing invocation |
| H31 | 43 | a --source-tip disagreeing with the derived tip renders a red verdict row | internal/preflight | a flag that is parsed and ignored looks identical to a verified one |
| H32 | 44 | a --source-tip that does not resolve returns a grammar error distinct from the mismatch red | internal/preflight | one message for both diagnoses sends a typo to the drift investigation |
| H33 | 31 | the guard documentation states the honest-mistake threat model unchanged | internal/conformance | a narrowed rim invites reading the guard as evasion-resistant |
| H35 | 29, 30 | the rim decodes an escaped shell operator to the operator, and refuses an escape it cannot decode | internal/systemtest | a placeholder decode welds two commands into one token, so the destructive half leaves command position and is allowed |

Not covered: story 9 — the residue deletion is a one-time tree edit with no durable predicate. A standing zero-count check would red on every future light-path ticket, because the close step requires that folder to survive the landing's own green gate run. H06 observes the resulting zero at the landing commit.

### Edge inventory

The shell-CLI hostile-input checklist in `projects/benchkit.md` walks here.
Live classes and their disposition:

- **A flag's value read as a positional** — H19 covers it on the guard's new spellings, supplying the sought token only as a flag value.
- **Required tool missing from PATH** — H20 is exactly this class: no wrapper on disk, no `bench` on PATH.
- **Absent versus present-but-empty** — H16 for `bench outline`, H03 for an absent spec folder.
- **Paths with spaces or glob characters** — the close step's slug and the residue guard's path both resolve exact paths, asserted at H01 and H26.
- **Control bytes in git-sourced text** — `--source-tip` values reach `toon.Table`, which refuses them, so H32's unresolvable path carries the assertion.
- **Special files and dangling symlinks in a discovered path** — `.logs` pruning stats before removing, asserted at H29.
- **A command whose own write changes a fact it reports** — the close step deletes inside the commit it publishes, asserted at H01 in the tracked configuration.
- **Hand-edited file with no trailing newline** — `.bench/structure-accept` is hand-edited; the surviving-grant assertion at H09 reads the file as written.

**Won't handle** — a tickets-only folder whose tickets are unfinished: the close step consumes the folder the landing names, and a builder landing an unfinished ticket has a workflow problem the CLI should not paper over.
**Won't handle** — a light-path folder that mixes a `spec.md` with tickets: it takes the spec route unchanged, and no in-scope caller creates the mixed shape.
**Won't handle** — deliberate evasion of the degraded rim: the rim is an honest-mistake layer by declared threat model, and the pre-push hook is the backstop.
**Won't handle** — pruning `.logs` by age or by total bytes: retention by count is one rule with one assertion, and no caller has asked for a second.
**Won't handle** — a `--source-tip` on a bare `bench preflight` with no subcommand: the grammar rejects the invocation before any flag is read.
**Won't handle** — restoring a removed `dist/bench` automatically: the warning names the rebuild command and the builder decides, matching how `bench repair` stays explicit.

## Ownership fences

- `internal/commit/`
- `internal/landing/`
- `internal/status/`
- `internal/outline/`
- `internal/bounds/` — the retirement of the row limit the bare form no longer needs
- `cmd/bench/` — the help inventory row that advertises `--full`
- `internal/gitguard/`
- `internal/gate/run_log.go` and `internal/gate/run_log_prune_test.go` — the pruner and the test that grades it
- `internal/preflight/`
- `internal/worktree/`
- `internal/sessioninspect/`
- `internal/systemtest/`
- `internal/anchors/registry_data.go`
- `internal/conformance/`
- `projects/benchkit.md` — the input-source row a new conformance check must carry
- `internal/structure/`
- `.bench/hooks/block-dangerous-git.sh`
- `.bench/hooks/session-start.sh`
- `.bench/structure-accept`
- `.agents/commands/bench-review-implementation.md`
- `.bench/BENCH-reference.md`
- `specs/ft229-hygiene-batch/`
- `specs/` — deletion only, for the tickets-only folders the residue ticket names

## Out of scope

- **The structure-policy decision** — whether `bench structure` gates at a frozen budget or leaves the ambient dashboard. The audit records it as a reviewer decision and this spec only removes the stale rows. `4 edits, 2 gate runs`.
- **An `bench outline` ranking rule** — choosing which symbols matter is a capability, not a shape change, and the audit named it a non-goal. `6 edits, 3 gate runs`.
- **Per-check gate evidence in the run log** — L-22's phase-versus-check granularity question, whose benefit the ledger marks as a hypothesis. `8 edits, 4 gate runs`.
- **A general retention policy for `capture/` and evidence records** — this spec prunes the gate's own logs only. `5 edits, 2 gate runs`.

## Further notes

FT174 depends on this spec's close step: the light-path close dispositions the
orphaned tickets FT174's row grammar would otherwise have to design around.
Landing FT229 unblocks it.
