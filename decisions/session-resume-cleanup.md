## Destination

Make a cold session resume cost nothing: provably-safe stale worktrees are cleaned
automatically at SessionStart, work intent survives session death via a
deterministic ledger, and one `bench status` verdict answers "is everything
committed and pushed?" — no agent re-investigation, no manual cleanup ritual.

## #1: Where does the cleanup fire?

Blocked by: none
Type: Grill

### Question

Auto at SessionStart, prevention at Stop, or a sharper manual prompt — where does
the burden actually disappear?

### Answer

Auto at SessionStart. The session-start hook runs the safe subset of
`bench worktree clean` before rendering `bench status`, so the dashboard reflects
post-clean state, and injects a one-line report of what was cleaned and what was
kept. The work is fully codified in the CLI — the hook adds no agent tokens beyond
the report line. Stop-time prevention alone was rejected: the failing case is the
session that dies uncleanly, where no end-of-session step ever runs.

## #2: What captures the handoff before the session goes cold?

Blocked by: #1
Type: Grill

### Question

When a session dies before its final verification, git state is recomputable at
resume but intent is lost — which worktree was doing what, and what step never ran.
What carries that across?

### Answer

A deterministic ledger, written by the tooling at the *start* of work (so it
survives unclean death) and refreshed by the Stop hook. It lives in the shared git
common directory — machine-written, never tracked, never model-authored. The gate
cache is the file-lifecycle precedent, but its worktree-local git-dir address is
not reused: every linked worktree must reach the same ledger. At resume,
`bench status` joins ledger intent with live git state into decision-grade rows
("worktree X: '<objective>', landed — cleaned" / "branch Y: unique commits, next:
/bench-final-check"). Model-written prose handoffs were rejected: they cost tokens
every stop and are lost exactly in the unclean-death case.

## #3: How aggressive is the unattended clean?

Blocked by: #1
Type: Grill

### Question

Manual `bench worktree clean` salvage-commits dirty WIP and removes the checkout.
May an unattended SessionStart hook do the same?

### Answer

No — conservative auto. The auto path removes only *clean* out-of-pool worktrees
and sweeps `worktree-*` branches proven landed in the default branch (ancestry or
patch containment, the existing `LandedInDefault` rule). Dirty, locked (live
agent), and leased worktrees are never touched unattended — they surface as ledger
rows with a next action. Salvage-commit stays exclusive to the manual
`bench worktree clean`. An unattended hook destroys nothing and commits nothing.

## #4: Which surfaces write intent into the ledger?

Blocked by: #2
Type: Grill

### Question

The pain-source worktrees (`.claude/worktrees/agent-*`) are created by the
harness's Agent tool — bench never sees them born. What are the ledger's writers?

### Answer

Three entry paths: `bench shift` records its objective (it already has one),
`bench worktree` gains an optional objective argument, and Claude Code's existing
`check-agent-line` PreToolUse hook records each Agent delegation's objective. The
Claude envelope cannot correlate that pre-spawn objective to the later agent id or
worktree, so it is explicitly uncorrelated intent until live git state can resolve
it. Codex does not expose the delegation objective on a deny-capable hook event;
OpenCode and headless adapters have no equivalent project hook, so those harnesses
retain only the two bench-owned writers. Recording commits, spec associations, or
gate verdicts in the ledger was rejected: the ledger carries only what git cannot
(intent), per the one-source rule.

## #5: What does the resume surface do about committed/pushed?

Blocked by: #1
Type: Grill

### Question

Half the resume burden is verifying everything landed. Push is reviewer-owned by
guard posture — how far does the surface go?

### Answer

A consolidated landed-state verdict in `bench status` (the single renderer — no
new command): unpushed commit count, branches still holding unique commits, dirty
paths — one glance, one remaining manual action (the push). Auto-push was
rejected: it would reopen the closed guard decision that pushing is the
reviewer's act.

## #6: Which modules and envelopes own each fact?

Blocked by: #3, #4, #5
Type: Research

### Question

Trace the exact seams before the spec locks: (a) what the `check-agent-line`
PreToolUse envelope actually carries — can a delegation be correlated with the
`worktree-agent-*` branch/worktree it later creates, or only recorded as
uncorrelated intent; (b) whether the Stop-hook envelope supports the ledger
refresh; (c) which Go package owns the ledger (new vs `internal/worktree` vs
`internal/status`) and where the conservative-clean subset lives relative to
`CleanCommand`; (d) how the new status rows fit the five-row budget and severity
ladder contract; (e) gate attachment for hook-driven behavior, with the observable
red signal per seam. Cite each claim to the owning source; probe the hook
envelopes against the live harness.

### Answer

The ledger is a new deep `internal/intent` module: it owns the common-git-dir
address, durable format, atomic writes, live/proven-done reconciliation, and safe
single-line objective encoding. `internal/worktree` keeps all cleanup authority;
its already-separated landed-branch sweep is reused and its out-of-pool remover
gains a clean-only, unlocked conservative path alongside (not as flags threaded
through) the salvage-capable manual path. `internal/status` remains a composer: it
reads typed intent and git/worktree facts, folds landed-state counts into the
existing severity-1 `git` signal, and folds unresolved intent into severity-2
`worktree` without creating a new rank. The SessionStart clean report is plain
hook context, not TOON, and does not consume the five-row budget. Exact envelope,
package, status, and gate findings are recorded in
[session-resume-cleanup-seams.md](session-resume-cleanup-seams.md).

## #7: When does a ledger row leave?

Blocked by: #2
Type: Grill

### Question

What removes a row from the ledger, and what happens to one that can never be
proven done?

### Answer

Proof-of-done only. A row leaves when its work is proven landed (the existing
ancestry/patch-containment rule) or its worktree is removed — by the auto-clean
or the manual command. A row that cannot be proven done never expires silently:
it surfaces in the status verdict as stale intent with a next action, however
old. Time-based expiry was rejected — silently dropping old rows would re-create
the lost-intent problem inside the tool built to prevent it.

## Not yet specified

- None. The SessionStart report is at most one deterministic plain-text context
  line (suppressed when every count is zero); status keeps its existing plain
  signal-row contract. Claude supplies delegation intent; Codex, OpenCode, and
  headless runs degrade to the bench-owned writers by design.

## Out of scope

- Auto-push or any hook-initiated push — push stays the reviewer's act.
- Model-written prose handoffs at session end.
- Salvage-committing WIP from the unattended path.
- Pool sizing or tightening changes — the warm pool is self-managing.
- Interactive confirmation flows in hooks — the auto path acts or surfaces, never asks.

## Handoff

1. **Module boundaries.** New `internal/intent` owns the shared-git-dir ledger,
   atomic persistence, safe objective encoding, and live/proven-done snapshots.
   `internal/worktree` owns conservative auto-clean and the existing manual salvage
   clean. `internal/git` owns typed dirty/ahead/unique-branch facts.
   `internal/status` only composes severity-sorted signals. `internal/shift`, the
   `bench worktree` entry, Claude's agent-line core path, and `internal/stophook`
   are thin ledger writers/refreshers. SessionStart remains a thin shell adapter.
2. **Contracts.** Ledger entries record only writer kind, objective, creation
   identity/time, and any worktree/branch identity actually known; the file is
   shared through `git rev-parse --path-format=absolute --git-common-dir` and is
   atomically replaced. Conservative clean removes only clean, unlocked,
   out-of-pool worktrees and landed orphan `worktree-*` branches; it never commits,
   forces, or touches dirty/locked/leased state. Its hook report is at most one
   plain line and is silent on all-zero. `bench status` keeps exit 0 and its bounded
   plain board: landed-state counts stay at severity 1, actionable intent/worktree
   state at severity 2, with `--all` preserving the overflow escape hatch.
3. **Deep vs thin.** Intent and worktree cleanup are deep modules: callers do not
   know file syntax, git-dir topology, landedness, locking, or reconciliation.
   Status, CLI dispatch, and hook scripts are thin projections/adapters and own no
   parser, landedness rule, or second ledger derivation.
4. **Black-box assertables.** Main, pool, and harness worktrees resolve one ledger;
   start-time writes survive killed sessions; denied Claude delegations do not
   create intent; Stop refresh is idempotent; auto-clean deletes only proven-safe
   state and changes neither HEAD nor dirty bytes; the report precedes post-clean
   status; status shows exact dirty/unpushed/unique-branch and open-intent counts,
   ranks gate < git < worktree, emits at most five default rows, and expands under
   `--all`.
5. **Gate attachment.** Package tests pin ledger addressing/atomicity/parsing and
   conservative classification. Runtime contracts drive the real CLI and hook
   scripts for Claude/Codex Stop envelopes, Claude Agent envelopes, auto-clean,
   status ordering/budget, and linked-worktree common-dir identity. Conformance
   pins Claude agent-line plus Claude/Codex SessionStart/Stop wiring. Targeted
   canaries remove the SessionStart clean call, common-dir addressing, and the
   landed-state status aggregation; each must produce its distinct owning-contract
   message. Live harness availability remains a manual compatibility probe, not
   the oracle.
6. **Hostile-input owners.** `internal/intent` owns spaces/globs/newlines/control
   bytes in objectives, missing final newline, absent/empty/malformed ledgers,
   concurrent writers, and interrupted atomic replacement. `internal/worktree`
   owns spaced/globbed paths, dirty/detached/locked/leased/stale registrations,
   missing git/default refs, symlink/deep-CWD invocation, and repeat cleanup.
   Hook/surface contracts own multiword objectives, source/linked/harness routing,
   missing wrappers, and plain-output safety.
7. **Uncertainty flags.** None. Claude's pre-spawn objective is intentionally
   uncorrelated; Codex/OpenCode/headless delegation intent is intentionally absent
   until a future harness exposes both objective and stable identity.
8. **Rejected alternatives.** Do not store the ledger in a caller's worktree-local
   git dir, parse transcripts, time-join concurrent hook events, put persistence in
   status or cleanup in intent, add a second landedness rule, add a new status
   command/rank, emit TOON from SessionStart, expire intent by age, auto-push, or
   salvage-commit from unattended cleanup.
9. **Domain watch-outs.** Stop hooks do not run on process death, so every durable
   intent write precedes work. Claude locks live agent worktrees; lock state is
   safety evidence, not clutter. A git worktree's absolute git dir is private while
   its common dir is shared. Objectives are hostile display bytes and must never be
   printed raw into terminal context.

Dependency order: build the intent/common-dir seam and its writers first; add the
conservative worktree primitive and SessionStart report second; extend typed git
facts and status composition third; finish Stop refresh, cross-harness wiring, and
canaries last. These are recommended green slices, not separate capabilities.
