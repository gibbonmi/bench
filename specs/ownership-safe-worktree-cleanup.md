# Ownership-safe worktree cleanup

Status: implemented

## Problem

Bench currently decides that a linked worktree is removable from its location and
Git cleanliness. `resume-clean` removes every clean, unlocked worktree outside the
pool, and `bench worktree clean` is a repository-wide sweep rather than an exact
target operation. Neither surface can prove that Bench created the registration or
that the checkout still belongs to the assignment Bench thinks it does.

The fallback for dirty attached work is also unsafe as a lifecycle primitive: it
commits the changes onto the checked-out branch. Dirty detached work is refused, and
clean detached work can be removed without first giving a unique commit a durable
name. The result is asymmetric preservation, branch mutation hidden inside cleanup,
and no durable recovery contract.

The same gap crosses harnesses. The shipped Claude Code settings have no
`WorktreeCreate` or `WorktreeRemove` adapter, so Claude's native lifecycle bypasses
Bench ownership and assignment state. SessionStart invokes cleanup but suppresses
its failures. A harness, path, branch prefix, elapsed time, or observed process is not
proof of ownership; asking an LLM to classify those signals would make destructive
cleanup less reproducible, not more.

Corrected throwaway-repository probes reproduce all four primary reds in the current
tree:

- `bench resume-clean` deletes both an ordinary foreign branch worktree and a clean
  detached worktree with a unique commit;
- `bench worktree "FT77 probe"` returns a detached, unlocked checkout with no
  ownership marker;
- `bench worktree clean <exact-path>` exits 2 with the old no-argument usage; and
- the no-argument dirty cleanup moves `worktree-recovery-probe`, removes its
  checkout, and creates no `refs/bench/recovery/...` ref.

## Solution

Make one deterministic worktree lifecycle core the only authority for creation and
cleanup. A Bench-controlled registration receives an immutable random ownership
identity in its Git-private administration directory. Each use receives a separate
versioned assignment record, exact generated branch, and Bench-identifiable Git
lock. Creation returns no usable path unless that bundle is complete.

Cleanup becomes a state machine over verified evidence. Automatic cleanup considers
only a matching cleanup-pending assignment with no live lease, an expected Bench
lock, a branch proven landed by the existing `LandedInDefault` rule, and no ignored
or unrepresentable residual. Git-visible dirty state is captured without changing
the assignment branch or real index, then anchored under a durable Bench recovery
ref before Bench unlocks anything. Foreign, active, unmerged, malformed, nested-dirty,
and uncertain states are retained and reported.

Replace broad explicit cleanup with an exact-path, non-interactive plan/apply
protocol. The dry run is pure and emits a versioned fingerprint of every
safety-relevant input. Apply accepts only that fingerprint and re-plans before any
mutation. Ignored deletion is a separately fingerprinted action with a bounded,
metadata-only inventory. Recovery refs have their own exact plan/apply retirement
surface and never leave automatically.

Expose the same core through the interactive Bench worktree command, a
harness-neutral create/release protocol, SessionStart/resume, and thin Claude Code
command-hook adapters. Marker creation, JSON parsing, ownership classification,
landedness, recovery, and acknowledgement are code; no prompt, agent hook, or LLM
participates.

## User stories

1. As a developer or harness creating a Bench worktree, I want creation to bind one
   immutable ownership identity to one explicit assignment and preserve every
   Git-visible state before cleanup, so that no cleanup path has to infer ownership
   or rewrite my assignment branch. Line: `gpt-5.6-luna` / low. The closed decision
   map fixes the lifecycle seam and the gate can observe the Git state, refs, records,
   and failure ordering in throwaway repositories.

2. As a developer resuming work or intentionally cleaning one checkout, I want
   unattended cleanup to act only on proven Bench assignments and explicit cleanup
   to apply only the exact state I inspected, so that foreign, changed, ignored, or
   ambiguous work survives with a concise explanation. Line: `gpt-5.6-luna` / low.
   The command grammar, state matrix, and runtime oracle are exact and fully
   executable without model judgment.

3. As a kit maintainer supporting multiple agent harnesses, I want every harness to
   route worktree lifecycle through the same deterministic core and a canary to prove
   the safety checks bite, so that Claude Code compatibility does not create a second
   cleanup policy. Line: `gpt-5.6-terra` / medium. Simulated hook contracts are
   gate-observable, but the real Claude lifecycle cadence is weakly gated and requires
   one fresh-session dogfood pass, so `craft-line` bumps this story one tier.

## Implementation decisions

### Ownership and assignment identity

The ownership marker lives in the linked worktree's private Git administration
directory resolved by Git, never in the checkout. It is a regular mode-0600 JSON file
with exactly these semantic fields: schema `bench-owner/v1`, a 128-bit
cryptographically random lowercase-hex owner ID, and the canonical absolute
worktree path. Bench writes it only while creating a new registration. It never
adopts an observed worktree and never edits a marker in place. Validation checks the
file type and mode, complete JSON shape, supported schema, ID encoding, canonical
path, current Git registration, and private admin location. Missing, empty,
malformed, broader-mode, symlinked, wrong-type, path-reused, or mismatched markers
classify as unowned or uncertain and cannot authorize automatic mutation.

The existing common-Git-dir intent store becomes the one persisted assignment
source. Its next schema accepts legacy records for reporting and safe compaction but
only a new assignment record can authorize this lifecycle. Each new record owns a
random assignment ID, owner ID, stable caller request ID digest, exact work-item
label, starting commit, exact full branch ref, state, and associated recovery
metadata. Raw labels and request IDs are never interpolated into paths or refs. The
generated branch namespace is
`refs/heads/bench/assign/<owner-id>/<assignment-id>`; the branch name is a locator,
not ownership evidence.

The assignment state machine is:

```
create complete ──▶ active ──release/remove event──▶ cleanup-pending
                                                    │
                         no recovery ref ◀──────────┤──▶ complete, then compactable
                                                    │
                         recovery ref(s) ◀──────────┘──▶ recovered
                                                               │
                                         explicit last-ref retirement
                                                               ▼
                                                        compactable
```

Automatic cleanup considers `cleanup-pending`, never `active`. A crashed active
assignment remains locked and visible until an exact explicit operation or the
future FT58 lease-reclamation capability resolves it. This is deliberate: age,
branch name, directory name, and an observed harness process do not transition
state. A live existing lease also blocks cleanup. Branch or worktree removal failure
leaves `cleanup-pending`; a recovery ref makes the terminal state `recovered` until
its last associated ref leaves explicitly.

The Git lock reason carries the owner and current assignment IDs. A matching Bench
lock is the protective backstop that the core may remove only inside a verified
cleanup transaction. A missing lock or a foreign/mismatched lock is uncertainty, not
permission. The repository's primary checkout is permanently excluded, whether or
not files resembling Bench metadata exist there.

Creation is one idempotent operation keyed by the caller's stable request ID. It
creates or reuses only its own fully verified result: exact branch, registration,
marker, assignment record, and lock. Git creates the registration locked from the
start. Any later failure returns no path; rollback removes only artifacts created by
that request, and an incomplete rollback leaves the registration locked with an
attributable record. Repeating the same request returns the same fully validated
active checkout; a changed label or conflicting state fails rather than silently
creating a second assignment.

### Eligibility and transaction ordering

All surfaces consume one typed lifecycle classification and result. No renderer,
hook, or status caller independently checks path prefixes, branch names, Git status,
or landedness.

For unattended cleanup, removal requires all of the following:

1. the target is a non-primary registration in the ambient repository;
2. marker, current registration, expected Bench lock, and assignment IDs match;
3. the assignment is `cleanup-pending` and no live lease exists;
4. the exact assignment branch is proven landed by `git.LandedInDefault`; and
5. Git-visible state is either clean or successfully recovered, ignored residuals
   are absent, and no dirty nested repository or submodule is present.

`LandedInDefault` remains the only landedness implementation: ancestry or complete
non-merge patch equivalence is sufficient. Squash ambiguity, merge-only history,
missing default, query failure, and incomplete patch containment retain the target.
The old branch-prefix sweep disappears. Cleanup deletes only the exact branch named
by a matching assignment record, after a fresh landedness proof and after the
worktree no longer has it checked out. A foreign explicit cleanup never deletes its
branch.

The mutation order is fixed: acquire the per-registration cleanup transaction,
revalidate the full plan, persist `cleanup-pending`, create and verify all required
recovery objects and refs, remove only the matching Bench lock, remove the exact
worktree, delete the exact landed assignment branch, then persist the terminal
record/result. If removal fails after unlock, Bench re-locks the surviving
registration before returning. Failure to re-lock is a distinct high-severity
residual. Cancellation and SIGINT use the same compensation path. No step reports
success from an intended later step.

### Lossless Git-visible recovery

Recovery refs use only encoded identities:
`refs/bench/recovery/<owner-id>/<assignment-id>/<ordinal>`. Exact foreign cleanup
uses a `foreign` owner segment derived from the repository identity, private
registration identity, and canonical path; creating that recovery metadata does not
adopt the worktree or grant it automatic-cleanup eligibility. Ordinals are chosen
deterministically from existing associated refs, so a pure re-plan is stable and a
partially completed apply can resume.

An isolated-index recovery envelope preserves the real worktree and index without
moving the checked-out branch or writing the real index. Its designated payload
commits are ordinary single-parent commits based on the pre-cleanup HEAD. Together
they preserve:

- the working-copy tree plus every untracked non-ignored path;
- a distinct staged tree when the index differs from both HEAD and the working tree;
  and
- the base, ours, and theirs index stages needed to reconstruct conflicted paths.

A root recovery object anchors every payload and its machine-readable layer
manifest; the recovery ref points at that root. The assignment record names the root
and every payload OID, so status and retirement do not parse commit prose. Identical
layers collapse to one payload. Deletions are represented by absence from the
payload tree, modes and symlinks retain Git semantics, and no recovery object embeds
ignored file contents.

Before unlock, Bench verifies that each object exists, the ref resolves to the
expected root, every payload is reachable from it, and the assignment record names
the same set. A clean detached target receives a recovery ref pointing to or
anchoring its HEAD before explicit removal. Dirty nested repositories and dirty
submodules are not representable by the parent recovery envelope and therefore
retain the entire worktree.

SessionStart and `resume-clean` never delete recovery refs. The exact retirement
surface shows the ref, root, payloads, and landed verdict, then requires a matching
fingerprint. Every designated payload must pass `LandedInDefault`; ambiguity retains
the ref. Apply deletes only that ref, updates the associated record, and compacts the
record only after its last recovery ref is gone.

### Exact plans, ignored residuals, and output

The harness-neutral lifecycle protocol is:

```
bench worktree create --request <opaque-id> --label <work-item>
bench worktree release --request <opaque-id> <path>
```

Create emits `worktree_create[1]{path,assignment,state}` on stdout. Release is the
trusted lifecycle transition used by harness adapters: it verifies the request,
marker, assignment, and exact path, changes `active` to `cleanup-pending`, and calls
the same automatic cleanup operation. The existing interactive
`bench worktree <objective>` wraps this protocol around the user's shell. Future
harnesses use the protocol rather than implementing marker or cleanup logic.

The explicit worktree grammar is:

```
bench worktree clean [--discard-ignored] [--full] <path>
bench worktree clean [--discard-ignored] [--full] <path> --apply <fingerprint>
bench worktree recovery <ref>
bench worktree recovery <ref> --apply <fingerprint>
```

`--` disambiguates a leading-dash path. No-argument cleanup and extra positional
targets are usage errors; there is no broad compatibility mode and no interactive
prompt. A dry run performs no Git, filesystem, marker, intent, lock, or ref mutation.
Apply computes a fresh plan under the transaction and accepts only an exact
fingerprint match. A malformed fingerprint is usage exit 2; a well-formed stale
fingerprint is an unsatisfied intent at exit 1 and returns the current non-mutating
plan.

The fingerprint is lowercase SHA-256 over a versioned canonical binary plan, not
over rendered TOON. It binds the ambient common Git directory and default ref/OID;
private registration identity and canonical path; marker bytes, type, mode, and
digest; assignment bytes, state, start, exact branch, and IDs; HEAD/ref/index and
content identities; lease and lock identities; landedness inputs and verdict;
nested-repository state; ignored inventory and `--discard-ignored`; recovery object,
ref, and payload actions; and the final ordered mutation list. Any change to those
inputs invalidates apply. Output-only `--full` does not change the fingerprint.

Explicit worktree dry-run and outcome use one stable flat schema:

```
worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:
```

`action` is a stable enum such as `retain`, `remove`, `recover-remove`,
`discard-remove`, `removed`, or `error`; `tracked` is `clean`, `dirty`,
`conflicted`, `nested-dirty`, or `unknown`; `ignored` carries count, apparent byte
total, shown count, and truncation; `recovery` is an exact ref or `none`; and
`detail` is a stable reason plus an actionable Bench command. Classification and
apply errors use the same one-row schema with `action=error`, not an unrelated prose
shape. Every block ends with a newline. A target containing a TOON-unrepresentable
control byte is rendered as an opaque SHA-256 label and retained; raw control bytes
never reach output.

Recovery-ref dry-run and outcome use the corresponding stable schema:

```
recovery_cleanup[1]{ref,root,payloads,landed,action,fingerprint,detail}:
```

The default ignored preview is a second `ignored_paths` block capped at 20 safely
escaped paths. `--full` raises only the display cap to the destructive inventory cap.
The inventory reads metadata with `lstat`, never file contents or symlink targets,
never follows a path outside the canonical worktree, and stops destructiveness at
1,000 entries or 1 GiB of apparent bytes. Crossing either bound, encountering a
stat/enumeration race, or finding an unrepresentable path refuses deletion. Above the
entry bound, the count is reported as `at-least=1001`; detail is never silently
omitted.

Ignored residuals always retain an unattended target. Explicit dry-run without
`--discard-ignored` plans `retain`. With that flag it may plan `discard-remove` only
inside the bounds; apply repeats the flag, re-inventories, compares the fingerprint,
and unlinks only the enumerated entries without following symlinks. Bench never reads
or archives ignored contents.

Exit 0 means success, an idempotent replay, or an intentional policy-preserved
result. Exit 1 means the requested classification, preservation, ref update, removal,
or revalidation could not be completed. Exit 2 is invocation error. Successful
applies leave a bounded cleanup receipt keyed by fingerprint in Git-private state;
the last 256 receipts are sufficient for exact retries after assignment compaction
and cannot authorize a new target. Partial transactions derive idempotency from the
assignment/ref state, not from repeating a destructive Git command blindly.

`resume-clean` and SessionStart render one compact operational summary from the same
typed results, with counts for removed, recovered, retained-by-reason, failed, and
open assignments. Unsafe path bytes are counted, not echoed. SessionStart remains
non-blocking, but it captures a nonzero resume result and prints an attributable
preservation warning instead of discarding stderr and claiming success.

### Harness adapters and structure

Claude Code receives two thin command-hook adapters in project settings. The
`WorktreeCreate` adapter parses the official JSON `session_id`, `cwd`, and `name`
fields in the Bench binary, derives a stable request ID, calls the harness-neutral
create operation, and prints only the absolute path on stdout as Claude requires.
Diagnostics go to stderr and any failure returns nonzero, so Claude does not receive
an incomplete path. The `WorktreeRemove` adapter parses `session_id` and
`worktree_path`, calls release, and contains no `git worktree remove`, unlock,
classification, or fallback deletion. Claude's removal event cannot block, so a
failed adapter deliberately leaves the matching Git lock in place; ordinary Git
removal then fails closed. These fields and output roles follow the
[official Claude Code hooks contract](https://code.claude.com/docs/en/hooks).

Hook commands resolve the repository-local Bench launcher from the shipped project
root and do not require a globally installed `bench`, `jq`, `readlink`, an LLM, or a
network call. Codex/OpenCode SessionStart and the real/by-path Bench CLIs remain
separate adapters over the same Go core. Linked kit copies ship the same settings,
scripts, runtime contract, and help/inventory text through the existing manifest
mechanism.

The implementation does not silently raise structure grants. `internal/worktree` is
already at its 12-file directory cap, so replacing the broad classifier/clean/resume
responsibilities must consolidate obsolete responsibility before adding a file or
stop for a reviewed grant. The already-granted large runtime worktree contract file
must replace unsafe expectations or earn a responsibility split; it must not merely
append another accepted exception.

## Testing decisions

- TDD applies only at the three named seams below. Each acceptance row goes red and
  green as one slice; tests are not batch-authored ahead of the behavior.
- Core lifecycle tests use real throwaway Git repositories and registrations. Fault
  injection attaches at the lifecycle's filesystem/Git execution boundary only to
  force ordering failures; tests assert returned state and reread Git, refs, locks,
  marker, index, branch, and records rather than mocking internal collaborators.
- Built-runtime contracts invoke the compiled real CLI, linked by-path launcher,
  SessionStart script, and hook commands. Existing tests that currently require
  broad removal or branch-mutating salvage are rewritten to assert the approved
  safety contract; they are not weakened or simply deleted.
- `internal/git` remains the landedness seam. Its existing ancestry, patch-equivalent,
  unique-patch, evil-merge, missing-default, and query-failure cases are reused rather
  than cloned in lifecycle code.
- A behavior-owned canary plants an ownership-bypass defect at the lifecycle family
  boundary. The existing SessionStart call-presence canary is not sufficient: the new
  fixture must turn green only if a foreign worktree can be deleted, recovery can
  occur after unlock, or the fail-closed lock can be bypassed, and the fixture-bite
  inventory must require it.
- After the automated gate is green, one fresh Claude Code `--worktree` run verifies
  the real create/remove cadence, exact stdout path handling, live marker/assignment/
  lock state, safe cleanup, and failure retention. This is a recorded dogfood verdict,
  not a substitute for runtime contracts.
- Gate command: `.bench/gate.sh`.

### Live Claude Code dogfood verdict

Live compatibility passed on 2026-07-12 with Claude Code v2.1.207 in fresh
`--worktree` sessions. The success lifecycle created a mode-0600 ownership marker,
matching active assignment, dedicated assignment branch, and exact Bench lock; normal
exit removed the checkout and branch, compacted the assignment, and retained complete
cleanup receipts. In the failure lifecycle, blocked recovery-ref creation moved the
assignment to `cleanup-pending` while preserving the dirty checkout and exact lock;
Claude retained the linked worktree and ordinary `git worktree remove` exited 128.
The same preservation trace is gate-owned through clean-drift and post-removal retry
cases, which re-anchor the recorded recovery root before removal and finalize one
`recovered` assignment with one terminal receipt.

### Seam diagrams

Core lifecycle seam:

```
trigger: interactive, automatic, explicit, or harness router
    │
    ▼
create request / cleanup request + real Git state
    ──▶ [ worktree lifecycle: create, plan, apply, automatic ]
    ──▶ registration + marker + assignment + refs + typed result
          ◀ tests attach here: drive a throwaway repo, inject boundary failure,
            then reread every externally durable Git/filesystem state
```

Built CLI and SessionStart seam:

```
trigger: agent command, resume, or fresh session
    │
    ▼
argv + cwd + linked launcher + registered fixture worktrees
    ──▶ [ compiled Bench CLI and SessionStart adapter ]
    ──▶ TOON / compact summary + exit + observable Git state
          ◀ tests attach here: invoke real and by-path surfaces from root/deep cwd
            and assert stdout bytes, exit code, refs, branches, and surviving paths
```

Multi-harness lifecycle seam:

```
trigger: generic harness request or Claude WorktreeCreate / WorktreeRemove event
    │
    ▼
generic request or official hook JSON
    ──▶ [ thin adapter ──▶ shared lifecycle core ]
    ──▶ exact path / cleanup result + locked or safely removed registration
          ◀ tests attach here: feed event JSON to shipped hooks; fresh Claude
            dogfood verifies the harness's real event cadence
```

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A Bench-controlled create returns only a dedicated assignment branch, immutable matching marker, assignment record, and matching Git lock; a retry with the same request is identical. | core lifecycle plus built create command | Observed red: in a throwaway repo, `BENCH_HOME=<tmp> SHELL=/bin/pwd bin/bench.sh worktree "FT77 probe"` exits 0 but inspection reports `marker=no`, `branch=detached`, and `locked=no`. Codify this probe as the first runtime red. | It reads the actual registration after the public command, so a returned path without any member of the creation bundle fails even if the command prints success. |
| 1 | Marker validation rejects absent, empty, malformed, wrong-type, wrong-mode, symlinked, unsupported-schema, invalid-ID, path-mismatched, registration-mismatched, reused-path, and legacy-only metadata without mutation. | core lifecycle | Red-first implementation: table-drive each enumerated marker fixture through plan before adding marker acceptance; each case must classify retained and leave the registration and refs byte-for-byte unchanged. | The cases cover every parser and identity boundary; a permissive parser, path-only classifier, or post-hoc adoption makes at least one unsafe fixture actionable. |
| 1 | Assignment matching requires the exact owner ID, assignment ID, request digest, start commit, full generated branch ref, state, and current registration; branch names and old schema records never substitute. | core lifecycle | Red-first implementation: pair one valid marker with records carrying each enumerated mismatch and run the planner; every mismatch must retain while the complete pair advances. | Varying one field at a time proves the join is real and prevents a branch prefix or partially matching record from masquerading as current ownership. |
| 1 | Automatic eligibility uses the existing ancestry or complete non-merge patch-equivalence proof and retains unique patch, evil merge, squash ambiguity, missing default, and Git query failure. | core lifecycle with existing Git landedness collaborator | Already covered at the Git/runtime seam by the current merged, content-landed, unique-patch, evil-merge, and unresolvable-default cases; red-first lifecycle integration must wrap those cases with valid ownership and assert only the first two become eligible. | Reusing the established black-box matrix detects a duplicated ancestry-only implementation and proves uncertainty does not become permission at the new caller. |
| 1 | Dirty recovery preserves staged-only, unstaged-only, staged-plus-different-working, deleted, untracked, renamed, symlink, executable-mode, and conflicted states in verified reachable payloads without changing the real index or assignment branch. | core lifecycle | Observed red: current no-argument cleanup of a dirty throwaway worktree exits 0 with `branch_moved=yes`, removes the checkout, and leaves `recovery_refs=none`. Implement the enumerated layers one red-green slice at a time. | Branch-tip, index-byte, payload-tree, and ref-reachability assertions reject both the current branch commit and a degenerate snapshot that keeps only the final working file. |
| 1 | A detached unique commit receives a durable recovery ref before explicit removal; a detached owned assignment mismatch is retained automatically. | core lifecycle plus built explicit command | Observed red: corrected `bench resume-clean` removes a clean foreign detached worktree with a unique commit and leaves only the primary registration. Add separate owned-drift and exact-foreign reds. | The unique commit is the only datum at risk; requiring a ref before path disappearance directly catches remove-first or detached-refusal asymmetry. |
| 1 | Dirty nested repositories and dirty submodules retain the parent worktree; clean gitlinks and ordinary parent files remain classifiable. | core lifecycle | Red-first implementation: create dirty nested-repo and dirty-submodule fixtures beside a clean gitlink control, then plan automatic cleanup before adding nested detection. | A parent-only tree can look clean enough while omitting nested content; survival plus no recovery ref prevents a false preservation claim. |
| 1 | Creation and cleanup interruption at every persistence boundary leaves either no returned path or a locked, attributable residual; recovery refs verify before unlock, and a surviving checkout is re-locked after removal failure. | core lifecycle fault boundary | Red-first implementation: inject failure after registration, marker, record, ref, unlock, removal, and branch-removal steps plus SIGINT at unlock; assert the enumerated durable state after each. | Step-specific rereads catch reordered writes, false success, unlocked partials, and cleanup that reports a later state it never reached. |
| 1 | Successful cleanup removes only the exact landed assignment branch, compacts complete records, and retains recovered context until the last associated ref is explicitly gone. | core lifecycle | Red-first implementation: seed two similar assignment branches and two recovery refs; complete one transaction and assert the sibling branch/ref survives and the record follows `cleanup-pending`, `recovered`, then compactable. | Similar siblings are the cheapest wrong broad sweep; the state assertions catch premature context loss and branch-prefix deletion. |
| 2 | `resume-clean` and SessionStart preserve ordinary branch and detached-unique foreign worktrees and report them as foreign. | built CLI and SessionStart | Observed red: corrected throwaway probe runs `BENCH_HOME=<tmp> bin/bench.sh resume-clean`, exits 0, reports `cleaned 2`, and both foreign directories disappear. Run the same fixture through SessionStart before changing the core. | It exercises the accused unattended commands and observes filesystem survival, so adding a marker only to new creation cannot make the red pass. |
| 2 | Automatic cleanup removes a matching cleanup-pending, non-live, landed clean assignment and recovers then removes the dirty equivalent; active, live-leased, unmerged, ignored, malformed, uncertain, and unexpectedly locked assignments survive by named reason. | built CLI and SessionStart | Red-first implementation: one runtime matrix enumerates the two eligible rows and seven retained rows with valid registrations, then invokes `resume-clean`. | The paired positive and negative matrix prevents the safe-looking no-op implementation while making every fail-closed reason observable. |
| 2 | No-argument cleanup is usage exit 2 and never sweeps worktrees or `worktree-*` branches; branch-prefix naming alone has no lifecycle effect. | built CLI | Observed red: current no-argument `bench worktree clean` moves the dirty `worktree-recovery-probe` branch and removes its checkout. Keep a same-prefix sibling in the first parser/runtime red. | The old command itself is the degenerate broad path; requiring a target before discovery makes ambient deletion impossible and the sibling catches a hidden compatibility sweep. |
| 2 | Exact worktree dry-run is pure and emits the stable TOON schema, deterministic fingerprint, bounded detail, trailing newline, and no prompt. | built CLI | Observed red: `bin/bench.sh worktree clean <fixture-path>` exits 2 with the old two-line usage while the target survives; after scaffolding, first assert the literal TOON bytes before apply exists. | Literal stdout plus before/after Git snapshots catches both the missing surface and a dry run that mutates while merely printing a plausible plan. |
| 2 | Apply revalidates marker, assignment, canonical registration, HEAD, branch, index/content, lease, lock, default OID, landedness, nested state, ignored inventory, recovery actions, and discard flag; drift in any one returns the new plan without mutation. | core lifecycle plus built CLI | Red-first implementation: generate a plan, change each enumerated bound input one at a time, then apply the old fingerprint and require exit 1 plus unchanged target state. | One-field drift cases prove the fingerprint covers the full safety boundary instead of only path and HEAD. |
| 2 | Exact foreign apply may remove only its registered non-primary target, never its branch; dirty or detached foreign state gets recovery first, while primary, cross-repository, inside-path, and unregistered targets are refused. | built CLI | Red-first implementation: plan two foreign sibling worktrees and apply one fingerprint, then repeat with dirty, detached, primary, other-repo, child-path, and unregistered fixtures. | Exact sibling survival catches broad discovery; branch/ref assertions distinguish safe checkout removal from adoption or data deletion. |
| 2 | Ignored residuals always block unattended cleanup; explicit discard reports metadata only, caps default detail at 20, caps destructive inventory at 1,000 entries and 1 GiB, binds the flag, and refuses uncertainty or drift. | core lifecycle plus built CLI | The current `TestResumeCleanKeepsIgnoredOnlyOutOfPoolWorktree` already covers the unattended keep. Red-first explicit slices exercise 0, 1, 20, 21, 1,000, and 1,001 entries; just-below and just-above byte limits; symlinks; stat failure; and changed inventory. | Boundary pairs catch off-by-one and unbounded output/deletion, while file-content read traps and symlink escape fixtures enforce the secret-safe inventory contract. |
| 2 | Recovery retirement shows the exact root and payloads, requires every payload landed plus an exact fingerprint, deletes one named ref only, and compacts context only after the last ref. | core lifecycle plus built recovery command | Red-first implementation: seed two refs with payload combinations of ancestor, patch-equivalent, unique, merge-ambiguous, and query-error states; apply one exact plan. | The sibling ref and mixed payload set catch automatic/broad deletion and any proof that checks only the envelope name or one convenient commit. |
| 2 | Explicit success, policy preservation, idempotent replay, failed intent, and usage return 0, 0, 0, 1, and 2 respectively through the same structured schema; resume summaries and SessionStart warnings are derived from typed results. | built CLI and SessionStart | Red-first implementation: table-drive literal stdout, stderr, newline, and exit for the five classes, then inject a resume Git failure through SessionStart. The current stderr-discarding, unconditional-success wrapper must fail the warning assertion. | Byte-level output and exit checks catch false success, stderr-only agent errors, prose drift, missing empty states, and the current swallowed failure. |
| 2 | Paths with spaces, glob characters, Unicode, newline, and leading dash; symlink invocation; deep cwd; quoted multiword labels; and missing global tools all reach the same exact target without injection or raw unsafe controls. | built real/by-path CLI and SessionStart | Red-first implementation: reuse the hostile-input fixture library and invoke each shipped surface from root and deep cwd with an empty PATH containing only fixture shims. | Real argv and output bytes expose quoting, cwd, resolver, and TOON escaping defects that unit-level canonicalization cannot see. |
| 2 | Concurrent plan/apply and interrupted retries create at most one recovery envelope, remove at most one registration/branch, and return the recorded result for an identical completed fingerprint. | core lifecycle plus built CLI | Red-first implementation: synchronize two applies at the transaction boundary and replay the winner's fingerprint after record compaction; run under `go test -race`. | A shared exact target makes duplicate refs, double removal, lost receipts, and check-then-act races observable rather than probabilistic. |
| 3 | The interactive command, harness-neutral create/release protocol, real CLI, and linked by-path CLI all produce and consume the same marker, assignment, lock, branch, and cleanup result. | multi-harness plus built CLI | Red-first implementation: run one lifecycle through each surface and compare semantic Git/private state after normalizing random IDs. The current public creation result is detached, unmarked, and unlocked. | State equivalence catches a thin-looking adapter that secretly creates its own branch, marker, or cleanup policy. |
| 3 | Claude `WorktreeCreate` accepts official event JSON, returns exactly one absolute path, and leaves a marked, assigned, dedicated-branch, locked checkout; `WorktreeRemove` routes that same path through release. | simulated multi-harness hook | Observed red: separate `rg -n 'WorktreeCreate' .claude .bench internal cmd bin` and `rg -n 'WorktreeRemove' .claude .bench internal cmd bin` probes both exit 1. Feed documented create/remove JSON to the shipped commands before adding settings wiring. | The event-level contract catches missing settings, wrong field names, extra stdout that corrupts Claude's path, and an adapter that bypasses the core. |
| 3 | Claude adapter failure before preservation leaves the checkout and matching Git lock intact, so ordinary `git worktree remove` cannot erase it. | simulated multi-harness hook plus Git substrate | Compatibility probe already shows ordinary removal of a locked fixture exits 128 and requires double force. Red-first adapter test injects preservation failure, invokes `WorktreeRemove`, then asserts the hook is nonzero, the lock remains, and ordinary removal still fails. | It tests the substrate backstop necessitated by Claude's non-blockable removal event; a script containing unconditional unlock or remove immediately loses the checkout. |
| 3 | A behavior-owned canary proves foreign preservation, recovery-before-unlock, and fail-closed locking cannot all be bypassed while the gate stays green. | runtime canary registry and fixture-bite inventory | Observed red: repository search finds only the call-presence canary `session-start-resume-cleanup-dropped` and no ownership canary. Register the planted defect first and require its EXPECT/bite before implementing the check. | A fixture that deliberately bypasses the family boundary turns green if enforcement is removed or weakened, so structural hook presence cannot masquerade as safety. |
| 3 | A real fresh Claude Code worktree session follows the documented create/remove cadence and retains a locked checkout when preservation is deliberately made to fail. | live Claude Code dogfood | Not TDD-able: Claude owns the external event cadence and `WorktreeRemove` has no blocking decision output. Record one successful lifecycle and one injected-failure lifecycle after the gate is green. | The live run catches settings loading, parent process, stdout routing, and fallback cleanup behavior that simulated JSON cannot establish. |

**Degenerate-implementation check.** The cheapest wrong change writes a marker after
observing any worktree, keeps branch-prefix cleanup, hashes only path and HEAD, commits
dirty state on the assignment branch, or lets the Claude hook run `git worktree
remove` directly. The marker mismatch, sibling, drift, recovery, and hook-failure rows
reject those versions. The cheapest wrong test only checks that SessionStart or a
hook mentions the new command; the foreign-survival runtime rows and behavior-owned
canary require the destructive family to bite. A no-op implementation is rejected by
the positive clean and dirty eligible-removal rows.

### Edge inventory

- **Generic error paths:** covered by the lifecycle ordering, drift, output/exit, and
  hook-failure rows. Git discovery, default resolution, marker/record read or write,
  object creation, ref update, lock/unlock, worktree removal, branch deletion,
  receipt write, ignored enumeration, and relock failures each retain attributable
  state and never claim a later step.
- **Empty and absent state:** covered by marker validation, assignment matching,
  automatic matrix, and explicit target rows. Absent versus empty marker, record,
  default ref, lease, lock, ignored inventory, recovery list, worktree registration,
  and stdout result are distinct cases; a definitive zero-result summary is not
  silence.
- **Boundary values:** covered by the ignored-inventory boundary pairs, primary and
  foreign sibling targets, zero/one/multiple recovery payloads, zero/one/multiple
  associated refs, one-character and multiword labels, and first/last receipt slots.
- **Malformed or hostile persisted data:** covered by marker, assignment, drift, and
  output rows. Wrong JSON types, unknown and old schemas, duplicate/trailing fields,
  missing trailing newline, invalid hex IDs, unsafe file modes, symlinks, truncated
  writes, unrepresentable controls, and marker/record identity disagreement all fail
  closed. Legacy intent entries remain reportable but never authorize removal.
- **Partial and interrupted operations:** covered by creation/apply fault injection
  and concurrent retry rows. SIGINT and failures before/after marker, record, recovery
  ref, unlock, worktree removal, and branch deletion leave the state machine at the
  last durable truth and compensate the protective lock when the path survives.
- **Re-run and concurrency:** covered by request-id creation, plan stability, stale
  apply, cleanup receipt, and synchronized apply rows. Identical create/apply/release
  is idempotent; a changed request or fingerprint is not conflated with a retry.
- **Path and argv hostility:** covered by the built-surface hostile row. Spaces, glob
  bytes, Unicode, newline, leading dash with `--`, relative and canonical absolute
  forms, symlinked launchers, deep cwd, multiword arguments, and traversal/cross-repo
  attempts cannot widen the exact target.
- **Output hostility:** covered by exact schema and hostile rows. Tab/newline/return
  use TOON escaping; ESC, BEL, form feed, and other unrepresentable controls become an
  opaque digest and a retained/error result. All stdout forms end in a newline, and
  no raw dependency diagnostic or secret content is emitted.
- **Environment and portability:** covered by real/by-path and hook rows. Missing
  global `bench`, `jq`, and `readlink`, a minimal PATH, linked kit invocation, every
  shipped SessionStart surface, and Claude's project-root setting use repo-local
  deterministic code. No runtime service or install-time fetch is introduced.
- **Nested repositories and submodules:** covered by the dedicated recovery row.
  Dirty nested state is a preservation failure; a clean submodule/gitlink remains
  ordinary Git-visible parent state.
- **Hostile local administrator:** **Won't handle** — a user who can rewrite the Git
  private directory, object database, binary, and records can forge any local
  lifecycle state; the contract detects accidental/cross-tool drift, not a malicious
  owner of the repository.
- **Non-Git Claude worktrees:** **Won't handle** — Bench's ownership, landedness, and
  recovery oracle is intentionally Git-specific; adding another VCS is a separate
  lifecycle backend rather than an edge of this implementation.

## Out of scope

- **Identity-safe stale-lease reclamation and private pool-root permissions (FT58)**
  — a separate concurrency and filesystem-permission capability; FT77 preserves an
  active or ambiguous lease rather than guessing. Estimated later cost: `~7 edits,
  4 gate runs`.
- **Shift-wide failure recovery and result-state semantics (FT79)** — a separate
  orchestration capability spanning iteration rollback and shift outcomes; this spec
  owns only per-worktree lifecycle state. Estimated later cost: `~10 edits, 5 gate
  runs`.
- **Removing the repository's primary checkout** — a separate destructive lifecycle
  with different invariants and no protective linked-registration boundary; FT77
  rejects it unconditionally. Estimated reversal cost: `~8 edits, 4 gate runs`.
- **LLM-based ownership, landedness, preservation, acknowledgement, or marker
  creation** — a separate probabilistic policy service that is deliberately rejected
  for this safety boundary; all required classifications are codified here.
  Estimated reversal cost: `~9 edits, 5 gate runs`.
