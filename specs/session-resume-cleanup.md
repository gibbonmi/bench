# Session resume cleanup

Status: implemented

## Problem

A cold session can recompute git state, but it cannot recover why a stale worktree
exists, which objective a dead session was pursuing, or which final verification
never ran. The reviewer or the next agent must inspect worktrees, branches, dirty
paths, and upstream state manually before deciding whether anything can be removed,
committed, or pushed.

The current cleanup command is unsuitable for unattended use because it may
salvage-commit dirty WIP before removing a checkout. The ambient dashboard also
reports dirty or unpushed state only for the current checkout and does not join that
evidence to durable work intent. A dead session therefore leaves both clutter and
missing context at the exact point where a cold resume needs a decision-grade answer.

## Solution

Persist machine-written work intent before work starts in one versioned ledger under
the repository's shared git common directory. Bench-owned shift and worktree entry
points write correlated intent; Claude Code's existing Agent PreToolUse path writes
explicitly uncorrelated intent because its pre-spawn envelope cannot name the later
agent worktree. Stop, manual cleanup, and SessionStart refresh the ledger only when
git provides proof that an entry is done; no entry expires by age.

Add a conservative cleanup primitive behind a new plumbing command. SessionStart
runs it before `bench status`, removing only clean, unlocked, out-of-pool worktrees
and orphan `worktree-*` branches proven landed by the existing ancestry or
patch-containment rule. Dirty, detached, locked, leased, unique-patch, and
unclassifiable state is kept. The hook emits at most one plain cleanup line and then
renders post-clean status.

Extend the existing status renderer rather than adding a command or severity rank.
Severity-1 `git` becomes the consolidated landed-state verdict across registered
worktrees and local branches; severity-2 `worktree` carries retained intent. The
default board stays bounded and `--all` expands intent detail. Push remains a manual
reviewer action.

## User stories

1. As an agent starting shift or worktree work, I want my objective persisted in the
   shared repository ledger before meaningful work begins, so that an unclean process
   death does not erase intent. Line: `gpt-5.6-luna` / low. The ledger contract and
   all three writer classes are exact and fully observable through package and runtime
   fixtures.

2. As a reviewer opening a cold session, I want SessionStart to remove only worktrees
   and scratch branches that are provably safe, so that stale clutter disappears
   without committing or destroying recoverable work. Line: `gpt-5.6-luna` / low.
   The destructive direction is tightly specified and the real CLI fixture can prove
   every kept and removed class.

3. As a reviewer resuming interrupted work, I want intent retained until git proves
   it landed or its checkout disappeared, so that neither elapsed time nor a missing
   Stop event silently discards the handoff. Line: `gpt-5.6-luna` / low. The lifecycle
   is a deterministic state transition with package and hook-level red signals.

4. As a reviewer checking whether everything is committed and pushed, I want one
   `bench status` invocation to count dirty paths, unpushed commits, unique local
   branches, and unresolved intent, so that I do not re-investigate each checkout.
   Line: `gpt-5.6-luna` / low. The existing status seam owns ordering and the complete
   aggregate can be asserted from controlled git fixtures.

5. As a user moving between Claude Code, Codex, OpenCode, and headless shifts, I want
   each supported lifecycle surface to capture as much intent as its real envelope
   provides without claiming nonexistent parity, so that portability remains honest.
   Line: `gpt-5.6-luna` / low. Harness asymmetry is closed in the map and wiring plus
   degraded behavior are mechanically checkable.

6. As a kit maintainer, I want package tests, real-path runtime contracts,
   conformance diagnostics, and targeted canaries to guard the ledger, cleanup,
   status, and hook wiring, so that weakening any safety proof turns the gate red.
   Line: `gpt-5.6-terra` / medium. This changes the oracle and follows the project's
   cached mid-tier routing for gate and conformance logic.

## Implementation decisions

`internal/intent` is the new deep module. It owns the ledger address, schema,
locking, atomic replacement, writer upserts, safe objective rendering, and
proof-of-done snapshots. Cleanup and status consume its typed interface; neither
parses its file. The ledger path is `<git-common-dir>/bench-intent.json`, resolved
with the absolute git common-directory query from every checkout. Its companion
lock is private to the module. Both are local git metadata and never tracked.

The ledger is mode `0600` and uses a versioned JSON object with `schema: 1` and an
ordered `entries` array. Each entry contains a stable writer key, writer kind
(`shift`, `worktree`, or `claude-agent`), the objective, an RFC3339 creation time,
and optional worktree and branch identities only when the writer actually knows
them. Claude uses the envelope's `tool_use_id` as its key and records
`tool_input.description`, falling back to a bounded prompt preview only when the
description is empty. Bench-owned writers create a process/time key before acquire
and update that same entry after the worktree or branch becomes known. Repeating an
upsert with the same key changes no bytes.

Writes take a bounded exclusive lock, write the complete next document to a sibling
temporary file, sync and close it, atomically rename it over the ledger, then release
the lock. A missing file means an empty ledger. A present-but-empty, unreadable,
unknown-schema, or malformed file is evidence of failure: preserve it, refuse to
overwrite it, and return an explicit error. A missing final newline is valid JSON and
does not change the result. Bench-owned `shift` and `worktree` entry points fail
before starting work when their initial intent cannot be persisted. Claude's
agent-line writer and lifecycle hooks retain their accepted fail-open posture: warn,
allow the harness action, and never forge or truncate existing ledger state.

Objectives are stored losslessly as JSON strings. Terminal output uses one shared
single-line encoder: preserve printable Unicode, escape newlines and every control
byte, and cap a preview at 120 Unicode code points without splitting UTF-8. A capped
preview carries the original byte count. Branch names and worktree paths pass through
the same encoder before entering SessionStart or status output.

The intent snapshot applies only proof-based transitions:

- A correlated entry stops surfacing when its recorded worktree no longer exists or
  its recorded branch is proven landed by the existing `LandedInDefault` rule.
- A branch that remains unique continues through the git landed-state signal even
  after its checkout disappears; intent need not duplicate a fact git now owns.
- Uncorrelated Claude entries stop surfacing only when no registered or orphan
  `worktree-agent-*` state remains. If any candidate remains, all uncorrelated entries
  remain rather than being joined by timestamp or ordering.
- Stop, manual clean, and automatic clean compact proven-done entries. Status uses the
  same filtered snapshot without mutating the ledger, so a query remains read-only.
- No timestamp expires an entry. An entry with no proof stays visible however old.

`internal/worktree` keeps all cleanup authority. Extend its registered-worktree fact
to retain branch, detached state, and porcelain `locked` state. Refactor the existing
landed-branch sweep to return typed actions so manual and automatic callers share one
classification and one `LandedInDefault` proof. Add a sibling conservative primitive;
do not pass an auto flag through the salvage-capable remover.

Conservative cleanup classifies candidates into mutually exclusive counts in this
order: leased pool entry, locked out-of-pool entry, dirty out-of-pool entry, or clean
unlocked out-of-pool entry. It removes only the final class without `--force`, then
prunes stale registrations. The branch sweep deletes only orphan `worktree-*` refs
proven landed by ancestry or patch containment. Dirty detached state, locked state,
lease state, unique patches, merge-only resolutions, failed git queries, and an
unresolvable default branch all fail toward keeping. The operation creates no commit,
does not change HEAD or the index, and does not alter dirty bytes.

The new operational plumbing command is `bench resume-clean` with no arguments.
Unknown arguments print plain usage to stderr and exit 2; outside a repository or a
classification/persistence failure exits 1 after deleting nothing not already proven
safe. Success exits 0. All-zero success is silent. Any nonzero candidate count renders
exactly one stdout line before post-clean status:

```text
bench resume: cleaned <worktrees> worktree(s), <branches> landed branch(es); kept <dirty> dirty, <locked> locked, <leased> leased; <intent> open intent(s)
```

SessionStart invokes `resume-clean`, preserves its stdout, tolerates its nonzero exit,
then invokes `bench status` and `bench guards --brief` in that order. It remains
non-blocking and silent outside a repository or when no wrapper resolves. Claude's
existing SessionStart wiring remains; Codex gains equivalent SessionStart wiring.
OpenCode and headless adapters gain no invented hook surface.

`bench worktree [objective...]` joins all positional objective words, preserving the
existing reserved `clean` subcommand. With no objective it keeps today's interactive
behavior and records a concise default intent. `bench shift` records before acquire,
then updates the entry with its pooled path and shift branch. Manual
`bench worktree clean` keeps its existing salvage behavior and refreshes intent after
its worktree and branch actions complete.

Claude's existing agent-line core records intent only after its model verdict allows
the Agent call. A denied delegation records nothing. Malformed or incomplete
envelopes keep the accepted agent-line fail-open behavior and warn when intent cannot
be recovered. No PostToolUse, transcript parser, temporal join, or claimed Codex
Agent writer is added.

Stop refresh runs before the current unarmed and `stop_hook_active` early returns, so
both normal and already-continued Stop events are idempotent refresh opportunities.
Refresh never changes the Stop guard's 0/2 verdict, never starts a gate when the shift
is unarmed, and never writes a gate-cache verdict. A refresh error warns and allows
the existing Stop path to continue.

`internal/git` exposes one typed landed-state fact gathered without network access:

- dirty paths are the de-duplicated porcelain entries across every registered
  worktree whose status query succeeds;
- unpushed commits are de-duplicated commit IDs ahead of each local branch's configured
  upstream; a branch without an upstream contributes to the unique-branch fact rather
  than an invented remote verdict; and
- unique branches are every non-default local branch not proven landed in the
  resolved default branch by `LandedInDefault`.

A failed required git query returns an unknown fact; status surfaces `git state
unavailable` at severity 1 instead of treating failure as zero. No fetch, push, or
remote provider call occurs during SessionStart or status.

The severity-1 `git` signal lists the nonzero dirty-path, unpushed-commit, and
unique-branch counts in that order. Its action is composed from the corresponding
canonical actions in that order: `commit on green`, `/bench-final-check`, and `push`.
Severity-2 intent is one aggregate row on the default board, with correlated and
uncorrelated counts plus the oldest safe objective preview. `bench status --all`
expands one severity-2 row per live intent with its known path/branch and safe
objective preview. The HTML dashboard and default status remain compact consumers of
the same typed signals. Gate remains severity 0; guards, drain, structure, decisions,
and housekeeping ranks do not move. The default board still prints at most five rows
and its existing `+N more` line.

The canonical plumbing inventory gains `resume-clean`; the project profile's ambient
dashboard and hook descriptions are updated only where their observable contracts
changed. The guard brief must continue to advertise `check-agent-line` as Claude-only
and Stop as Claude/Codex.

Gate coverage extends existing owners rather than creating duplicate registries:

- `internal/intent` package tests own address, schema, locking, atomicity, encoding,
  concurrency, and lifecycle transitions.
- Existing worktree runtime contracts own conservative removal and branch proof
  cases; existing status runtime contracts own aggregation, ordering, budget, and
  expanded intent rows; existing gate runtime contracts own Stop refresh; the
  existing agent-line executable contract owns Claude intent capture; the existing
  SessionStart contract owns report order and failure posture.
- Root conformance requires Claude Agent/Stop/SessionStart and Codex
  Stop/SessionStart wiring, while preserving the explicit absence of a Codex Agent
  writer.
- Three behavior-owned canaries sabotage SessionStart cleanup invocation,
  common-directory intent addressing, and landed-state status aggregation. Each
  expects its owning contract's distinct message; no canary reimplements the ledger
  parser or landedness rule.

## Testing decisions

- A good test drives the typed intent interface or the real CLI/hook from a throwaway
  git repository and observes ledger bytes, exit code, stdout/stderr, refs,
  registrations, HEAD/index, and dirty-file state. Tests do not inspect private helper
  calls or reproduce landedness in assertions.
- The new package seam covers failure modes too expensive or ambiguous to reach
  through a hook: concurrent writers, interrupted replacement, malformed documents,
  exact objective encoding boundaries, and proof-only snapshot transitions.
- Real-path runtime tests remain the authority for destructive cleanup, writer/hook
  routing, and status bytes. Existing manual-clean, Stop gate-cache, status budget,
  and guard-brief cases are regression assertions, not rewritten copies.
- The gate command is `.bench/gate.sh`.

### Seam diagram

Intent persistence seam:

    trigger: shift/worktree/Claude writer, Stop or cleanup refresh, status read
        │
        ▼
    root + typed entry  ──▶  [ internal/intent ledger ]  ──▶  atomic common-dir bytes
    live git facts      ──▶  [                        ]  ──▶  filtered typed snapshot
                                ◀ tests attach here: drive public writer/snapshot calls
                                  and observe bytes, errors, idempotency, and lifecycle

Conservative cleanup seam:

    trigger: SessionStart invokes bench resume-clean
        │
        ▼
    registered worktrees ──▶  [ worktree conservative cleanup ]  ──▶  kept/removed facts
    local scratch refs   ──▶  [ + shared landed proof          ]  ──▶  one report / exit
                                    ◀ tests attach here: real CLI fixture observes refs,
                                      registrations, HEAD/index, dirty bytes, and stdout

Status seam:

    trigger: SessionStart, reviewer, or HTML dashboard requests status facts
        │
        ▼
    git/worktree facts  ──▶  [ status signal composer ]  ──▶  severity-ordered board
    intent snapshot     ──▶  [                        ]  ──▶  compact / expanded rows
                                ◀ tests attach here: real status command asserts counts,
                                  safe bytes, rank, five-row budget, and --all expansion

Hook and wiring seam:

    trigger: Claude Agent, Claude/Codex Stop, Claude/Codex SessionStart
        │
        ▼
    real hook envelope  ──▶  [ shared hook shims + Go core ]  ──▶  intent/report/verdict
    harness config      ──▶  [ root conformance          ]  ──▶  wiring diagnostics
                                  ◀ tests attach here: execute shims with representative
                                    envelopes and plant unwired/broken configuration

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Main, pooled, and Claude-agent checkouts resolve one `bench-intent.json` in the shared git common directory. | `internal/intent` public address and write/read interface | Observed red: `test -e .git/bench-intent.json` exits 1; red-first package coverage writes from all three checkout classes and must fail until every caller resolves the common directory. | A worktree-local implementation creates different files and the cross-checkout read cannot see the writer's entry. |
| 1 | The schema-1 ledger round-trips all three writer kinds and optional identities; same-key upserts are byte-idempotent. | `internal/intent` writer/snapshot interface | Red-first implementation: add the schema/upsert table test before production types and run `go test -count=1 ./internal/intent`; it must fail on the missing package. | An append-only duplicate writer or lossy schema cannot satisfy exact second-write bytes and typed round-trip assertions. |
| 1 | Shift and worktree persist intent before acquire/work begins, then enrich the same entry; an interrupted writer leaves readable intent. | real `bench shift` and `bench worktree` commands | Red-first runtime cases stop each command immediately after the initial write and assert the common ledger before allowing acquire; run the existing runtime Shift and Worktree contract groups. | Writing only at teardown, only inside the acquired checkout, or creating a second enrichment entry fails the pre-acquire and same-key assertions. |
| 1 | Multiword objectives remain one objective; JSON control bytes are lossless and output previews escape controls at 0/1/120/121-code-point and multibyte boundaries. | intent package plus real worktree/shift CLI | Red-first package/CLI tables feed spaced, globbed, newline, ESC, BEL, and split-rune inputs; `go test -count=1 ./internal/intent ./internal/contract/runtime` must fail before the encoder and argv join exist. | `$1`-only parsing, raw terminal output, byte-splitting, or silent truncation changes the recorded value or emitted preview. |
| edge of 1 | Concurrent writers serialize without a lost entry; a killed writer before rename leaves the prior ledger readable and a reclaimable lock. | intent package persistence interface | Red-first deterministic interleave tests hold the lock and stop at write/sync/rename seams; `go test -count=1 ./internal/intent` must fail before locking and atomic replacement exist. | A read-modify-write race loses one key, while in-place writes expose partial JSON after interruption. |
| edge of 1 | Absent ledger means empty; present empty, malformed, unreadable, or unknown-schema evidence is preserved and fails loudly; no-final-newline valid JSON succeeds. | intent package and status projection | Red-first parse table enumerates the five states and runs `go test -count=1 ./internal/intent`; the empty-file case must fail before the parser distinguishes it from absence. | Treating every read error as empty would overwrite evidence and produce a false-clean resume. |
| 2 | Auto-clean removes exactly clean, unlocked, out-of-pool worktrees and never uses force or creates a commit. | real `bench resume-clean` runtime contract | Observed red: `bash .bench/hooks/session-start.sh` exits 0 without any `bench resume:` line; the red-first runtime fixture plants all worktree classes and reports `session-start conservative cleanup contract failed` until only the safe class disappears. | A no-op leaves the safe checkout; reuse of manual clean commits dirty WIP; a force path deletes a locked or dirty checkout. |
| 2 | The branch sweep deletes ancestry-landed and patch-landed orphan `worktree-*` refs while keeping unique patches, merge-only content, active refs, non-scratch refs, and all refs when default resolution fails. | existing worktree branch-sweep interface through real CLI | Extend the existing runtime sweep table first and run `go test -count=1 ./internal/contract/runtime -run '^TestRuntimeWorktreeContracts$'`; new auto-path cases must fail before typed sweep results are shared. | Enumerating both landed proofs and four keep classes prevents a branch-name-only or `git cherry`-only deletion shortcut. |
| 2 | Cleanup changes neither main/worktree HEAD, index, dirty bytes, nor the commit graph except for deletion of proven-landed scratch refs. | real resume-clean fixture state manifest | Red-first runtime test captures hashes, index tree, porcelain bytes, file bytes, and ref set before/after; it reports the changed artifact under `session-start conservative cleanup contract failed`. | Output-only assertions would miss destructive side effects; the manifest makes every prohibited mutation observable. |
| 2 | SessionStart emits zero lines on all-zero cleanup or exactly one canonical report before post-clean status, then guard brief; errors keep state and never block session open. | real SessionStart hook | Observed red: the current hook has no `bench resume:` line; extend `testSessionStartGuardBriefInjection` with zero/nonzero/error cases before adding the call. | A report after status exposes stale counts, repeated per-item output breaks the token bound, and a propagated error bricks session start. |
| 3 | Correlated entries leave the live snapshot only when their checkout disappears or their branch is proven landed; unique branch state transfers to the git signal. | intent snapshot with shared landed proof | Red-first lifecycle table enumerates worktree-present/absent and ancestry/patch/unique/merge-only branch states; `go test -count=1 ./internal/intent` fails before proof filtering exists. | Time expiry, branch-name inference, and an all-done stub each misclassify at least one enumerated state. |
| 3 | Uncorrelated Claude entries remain while any registered or orphan `worktree-agent-*` candidate exists and all leave only when the candidate set is empty. | intent snapshot uncorrelated-set rule | Red-first table covers zero, one, and several candidates plus several simultaneous intents; `go test -count=1 ./internal/intent` fails before set-level reconciliation exists. | Temporal or one-to-one joins can pass a single-entry case but fail the parallel many-to-many cases. |
| 3 | Stop refresh runs for armed, unarmed, active, Claude-shaped, and Codex-shaped envelopes without changing existing Stop exit or gate-cache behavior; repeat refresh is byte-idempotent. | real Stop hook and `internal/stophook` orchestration | Observed red: an unarmed representative Stop exits 0 but leaves `.git/bench-intent.json` absent; extend `TestRuntimeGateContracts` first and require `stop hook intent refresh contract failed`. | Keeping refresh behind the armed/active early return fails two classes; coupling it to gate verdict recording changes existing cache or exit assertions. |
| 3 | Manual and automatic clean compact only entries their completed actions prove done; malformed ledger state is never overwritten. | real manual-clean and resume-clean commands | Red-first runtime fixture seeds correlated, unrelated, and malformed entries before both commands; `go test -count=1 ./internal/contract/runtime -run '^TestRuntimeWorktreeContracts$'` must fail until refresh is wired. | Clearing the whole ledger after any cleanup loses unrelated intent; skipping refresh leaves already-proven entries forever. |
| 4 | Status emits exact de-duplicated dirty-path, upstream-ahead commit, and unique non-default branch counts across all registered worktrees/local branches without network calls. | real `bench status` command over typed git facts | Observed red: `bench status --all` does not match the three counted fields; red-first runtime case reports `status landed-state contract failed` until all three sets are correct. | Current-checkout-only logic, per-branch double counting, or an always-zero aggregate fails the populated multi-worktree/multi-branch fixture. |
| 4 | Failed worktree/status/upstream/default queries surface `git state unavailable` rather than a zero count, and branches without upstream remain visible through unique-branch state. | typed git fact and status projection | Red-first PATH/git-failure and no-upstream fixture runs `go test -count=1 ./internal/contract/runtime -run '^TestRuntimeStatusContracts$'`. | Swallowed git errors and upstream-only implementations both render a false clean board. |
| 4 | Gate remains severity 0, landed-state git remains 1, all worktree/intent state remains 2, and default output never exceeds five rows with the existing overflow line. | status signal composer and renderer | Extend the existing ladder/budget fixture before code and run `go test -count=1 ./internal/contract/runtime -run '^TestRuntimeStatusContracts$'`; wrong order or six rows must fail. | Adding a new rank or one default row per intent would displace higher-severity signals or exceed the fixed budget. |
| 4 | Default status aggregates intent with the oldest safe preview; `--all` emits one safely encoded row per live correlated/uncorrelated entry with known identities only. | real status compact/expanded rendering | Red-first runtime fixture seeds 0, 1, and 6 hostile entries and asserts compact versus expanded row counts/bytes. | A count-only implementation loses decision context; raw or always-expanded output leaks control bytes and defeats the token budget. |
| 4 | A second unchanged status invocation is byte-identical and mutates no ledger, refs, index, gate cache, or files. | real status command read-only/idempotency manifest | Red-first runtime test compares stdout and a before/after state manifest around two calls. | Reconciliation hidden inside status or nondeterministic map/time order changes bytes or repository metadata. |
| 5 | Allowed Claude Agent calls record one uncorrelated description; denied calls record none; missing/malformed intent fields warn and preserve the existing allow/deny posture. | real `check-agent-line.sh` executable contract | Extend `TestCheckAgentLineExec` before the core writer and run `go test -count=1 ./internal/conformance -run '^TestCheckAgentLineExec$'`; it must report `check-agent-line intent capture contract failed`. | Recording before verdict leaves denied intent, while making capture failure fatal silently strengthens an accepted fail-open rim. |
| 5 | Claude wires Agent, Stop, and SessionStart; Codex wires Stop and SessionStart but not Agent; OpenCode/headless use only shift/worktree writers. | root conformance over real harness configs | Observed red: `rg -q '"SessionStart"' .codex/hooks.json` exits 1; add the distinct Codex SessionStart diagnostic before editing config and run root conformance red-first. | A broad “hooks exist” check can pass with the wrong event and can falsely advertise Codex Agent parity. |
| 5 | Source-kit, linked-repo by-path, deep-CWD, symlinked wrapper, Claude hook, and Codex hook routes reach the same resume-clean/status implementation; missing wrappers never block SessionStart. | existing wrapper/link and hook runtime seams | Red-first surface/runtime table invokes the six enumerated routes and a missing-core case before adding `resume-clean` dispatch. | Testing only the source wrapper misses linked binary drift, cwd assumptions, shell quoting, and fail-open rims. |
| 6 | Package/runtime/conformance checks produce distinct attributable failures for intent common-dir behavior, conservative SessionStart cleanup, landed-state status, Stop refresh, Claude capture, and Codex SessionStart wiring. | owning gate packages on real paths | Red-first implementation adds the named assertions before behavior and runs `.bench/gate.sh`; each of the six enumerated classes must fail with its own message. | One generic failure or a source scan cannot identify which safety proof disappeared and may remain green through a broken real path. |
| 6 | Three behavior-owned canaries permanently sabotage SessionStart cleanup invocation, common-dir intent addressing, and landed-state status aggregation. | canary registry, fixture inventory, and real inner gate | Observed red: `test -e tests/canary/behavior-owned/session-start-resume-cleanup-dropped/EXPECT` exits 1; register all three fixtures before their behavior and run the owning conformance/contract tests. | The three planted failures prove the new contract families still bite without duplicating their parsers or git derivations. |

**Degenerate-implementation check.** The cheapest wrong ledger writes a worktree-local
JSON file at Stop time only; rows 1–6 reject its address, timing, concurrency, and
malformed-input posture. The cheapest wrong cleanup calls existing manual clean;
rows 7–10 catch salvage commits, unsafe removals, side effects, and report order. The
cheapest wrong status counts only the current checkout and prints one row per entry;
rows 15–19 catch incomplete sets, false empties, rank/budget drift, unsafe bytes, and
query mutation. The cheapest wrong hook patch claims all harnesses are equivalent;
rows 20–22 reject denied-call records, missing Codex SessionStart, invented Codex
Agent capture, and wrapper drift. Rows 23–24 keep the oracle from becoming an
always-green assertion.

### Edge inventory

- **Error path — ledger lock/read/write failure:** covered by stories 1, 3, and 5;
  bench-owned work fails before starting, lifecycle hooks warn/fail open, status
  surfaces the unreadable state, and existing bytes remain untouched.
- **Error path — git classification/default/upstream failure:** covered by stories 2
  and 4; cleanup keeps all unproven state and status emits an unavailable signal.
- **Empty/absent input:** covered by the absent-versus-empty ledger row, all-zero
  cleanup suppression, no-objective worktree default, no-upstream branch posture, and
  zero-intent compact/expanded status cases.
- **Boundary values:** covered at 0/1/many entries and candidates, 0/1/5/6 dashboard
  rows, 0/1/120/121 objective code points, a multibyte truncation boundary, zero/one/
  several local branches, and zero/one/several configured upstreams.
- **Malformed input:** covered for invalid/unknown-schema JSON, missing final newline,
  malformed Claude envelopes, malformed lock ownership, invalid git output, and
  unresolvable default refs.
- **Interrupted or partial state:** covered by a writer stopped before rename, a stale
  lock owner, shift/worktree death after initial intent but before enrichment, absent
  Stop, and SessionStart rerun after partial cleanup.
- **Re-run idempotency:** covered for same-key upsert, Stop refresh, manual/automatic
  cleanup, SessionStart, status bytes, and repeated link/wiring installation.
- **Hostile environment — spaces and glob characters:** covered in objective text,
  branch names, common/worktree paths, deep CWD, and wrapper paths.
- **Hostile environment — control bytes:** covered for objective, branch, and path
  rendering with newline, ESC, and BEL; JSON storage stays lossless and terminal
  context stays one line.
- **Hostile environment — no trailing newline and absent versus empty:** covered by
  the ledger parser table; the machine writer emits canonical newline-terminated JSON
  but accepts valid hand-edited JSON without the newline.
- **Hostile environment — unquoted multiword arguments:** covered by
  `bench worktree [objective...]` and the existing `bench shift` argv join; `clean`
  remains the one reserved subcommand.
- **Hostile environment — required tool missing from PATH:** covered by missing git,
  missing repo-local/global wrapper, and existing hook fail-open tests. No new `jq`,
  provider CLI, or network tool becomes a runtime dependency.
- **Hostile environment — symlink and every shipped surface:** covered across the
  source wrapper, global/repo-local wrapper, linked worktree, Claude hooks, Codex
  hooks, and headless adapters. OpenCode has no invented lifecycle hook.
- **Hostile environment — SIGINT mid-operation:** covered by lock cleanup and atomic
  replacement at the intent seam; cleanup creates no scratch worktree or commit to
  recover.
- **Hostile environment — cwd below root:** covered by common-dir, resume-clean,
  status, and hook invocations from a nested directory.
- **Remote freshness:** **Won't handle** — SessionStart remains offline; unpushed state
  is computed from local configured upstream refs, while branches without upstreams
  remain visible as unique local work.
- **Native non-git worktree providers:** **Won't handle** — this repository's declared
  shell-CLI domain is git; Claude's custom non-git WorktreeCreate/Remove protocol is a
  separate VCS integration capability.

## Out of scope

- **Automatic or hook-initiated push** — a separate remote-mutation capability that
  would supersede the accepted reviewer-owned push posture. Estimated later cost:
  `~8 edits, 4 gate runs` plus a new reviewer decision.
- **Codex/OpenCode delegation-objective parity** — a separate harness-integration
  capability blocked until those harnesses expose both objective and stable identity
  on an auditable lifecycle event. Estimated later cost after that prerequisite:
  `~7 edits, 3 gate runs`.
- **Warm-pool sizing or tighter reclamation** — a separate resource-management
  capability; this feature only distinguishes leased pool state from unattended
  cleanup candidates. Estimated later cost: `~5 edits, 2 gate runs`.
- **Model-written prose handoffs** — a separate narrative handoff capability; the
  ledger intentionally stores machine intent rather than generated summaries.
  Estimated later cost: `~6 edits, 2 gate runs`.
- **Interactive cleanup confirmations** — a separate interactive CLI capability;
  SessionStart remains deterministic and non-interactive. Estimated later cost:
  `~4 edits, 2 gate runs`.
