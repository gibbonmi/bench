# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-integrate-learnings` reviews the open entries, promotes the
generalizable ones into the kit with sign-off, and prunes them: a resolved entry
leaves this file, and its verdict (promoted or dismissed, one line of why) is
recorded in the integration commit and CHANGELOG. The journal holds open entries
only; history lives in git. Never rewrite a kit rule yourself — that is the whole
point of capturing here instead.

Format per entry:

## <date> — <short title>  [open]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry leaves this file only via /bench-integrate-learnings.

<!-- entries below -->

## 2026-07-03 — session-start stale gate confuses without a benign/real split  [open]
- **What happened:** Reviewer flagged that the gate "is almost always stale in a
  new session," which reads as alarming. Diagnosis: the verdict is content-addressed,
  so it goes stale the instant the tree moves past the last green — and sessions
  routinely end with a change after the last gate run (a manual commit that wasn't
  re-gated, or a `bench idea` park that dirties ROADMAP.md). So new sessions almost
  always open stale, but the drift is often benign (capture-scratch like ROADMAP.md /
  .bench-notes.md, which no gate check reads) rather than unverified code.
- **Right behavior:** At session start, when the gate reads stale, tell the user
  *why* — split "benign drift only (e.g. a parked idea) → just a reminder to re-run
  the gate" from "committed code moved since the last green → real, re-run before you
  trust it." The bare word "stale → re-run the gate" hides that distinction and reads
  as an error even when it's harmless.
- **Proposed rule change:** Consider having `bench status` classify a stale verdict:
  if the diff from the gated tree is confined to capture-scratch paths, word it as a
  benign reminder; otherwise flag it as real. (Distinct from, and a lighter-weight
  alternative to, the parked idea of carving capture-scratch out of `gate_tree_hash`
  entirely — that changes the oracle's key and is sensitive on the tripwire branch.)

## 2026-07-03 — dogfooding gotcha preserved from a retired map: mid-session edits hit the next session  [open]
- **What happened:** Story 11 of spec-handoff-lifecycle retired `decisions/dogfooding.md`.
  Its self-host+canary-guard decision was already realized (the canary gate layer
  shipped; craft-synthesis carries the dogfood loop), so the promotion read ended in
  deletion. But one operating gotcha in that map was not recorded anywhere: the kit
  loads skills/commands at session start, so a mid-session edit to a skill or command
  lands in the *next* session, not the current one. The egg bites in three places —
  the gate (break it, lose the oracle), skills/commands (break a trigger, next session
  misfires), and the CLI (break `bench shift`, lose the loop).
- **Right behavior:** Per the map-Handoff rule (item 9), operating lessons go through
  `.bench/learnings.md` and `/bench-integrate-learnings`, not per-spec notes or a
  unilateral skill edit — so this gotcha is captured here rather than promoted into a
  skill on the worker's own authority.
- **Proposed rule change:** Consider a one-line note in `bench-craft-synthesis`'s
  dogfood loop: when a candidate change touches a skill or command trigger, the
  dogfood shift must be a *fresh* session, because the edit does not take effect in
  the session that made it. (Reviewer's call — this is the deferred/contestable
  promotion flagged from story 11.)

## 2026-07-03 — review-implementation: verify a returned finding in the live session, not a worktree  [open]
- **What happened:** During `/bench-review-implementation` on spec-sourcing, the
  Coverage axis returned a real finding (the negative anchor's line-oriented grep
  misses a hard-wrapped bypass reintroduction). I reproduced and fixed it directly
  in the live session's working tree — one repro, one fix, one gate run — rather than
  spinning up a separate worktree for the fix.
- **Right behavior:** When a review delegate returns with an issue, the invoking
  session tests and fixes the failure in the live session, not in a separate
  worktree. The worktree-isolation rule is for *write-delegations* (a subagent that
  edits files, run isolated so stray edits can't land in reviewer-owned files); the
  review axes themselves are read-only, so their findings come back to the one live
  session that owns the diff, which verifies and fixes them in place against the gate.
- **Proposed rule change:** Consider a one-line note in `/bench-review-implementation`
  (or `craft-delegate`): review delegates are read-only and return findings; the
  live session verifies and fixes a returned finding in place — worktree isolation
  is for write-delegations, not for reproducing a review finding.

## 2026-07-03 — routing a hook-invoked command to Go broke the linked-repo by-path CLI  [open]

What happened: go-axi-query-surface (slice 2) flipped the router so learnings|maps|guards|diff|coverage exec the Go binary via route_binary. The Go surface is a verified byte-identical drop-in and all AXI/status contracts pass when driven through the REAL kit CLI ($root/bin/bench.sh, kit_dir=$root, dist/bench present). But the linked-repo by-path CLI (.bench/bin/bench.sh, kit_dir=.bench) has NO route_binary candidate that reaches the binary in any deployment — not the gate's copy-mode fixture, not an npm consumer (its <repo>/node_modules/@benchkit/<pkg> path is not among the candidates .bench/dist/bench, .bench/node_modules/<pkg>, .bench/../<pkg>). So session-start's `guards --brief` (the failing contract) and status's learnings/maps adapters silently break for every linked/consumer repo the moment those commands route to Go. Slice 1 never exposed this because only `version` routed and nothing invoked it by-path.

Right behavior: the spec (and the map's Handoff) should have flagged that porting any command session-start/status invoke through the by-path CLI requires the linked-repo binary-distribution to be solved first — a slice-6 (link/doctor) dependency that the dependency order put last. The router flip's consumer-correctness depends on it.

Proposed rule change: in /bench-write-spec's edge inventory, add a check — "for every ported command, which surfaces invoke it, and through which CLI (real kit vs linked by-path)?" A by-path invocation of a routed command is a distribution dependency to resolve or explicitly defer, not assume-works. The "missing-binary path already proven by slice 1" reasoning hid this because slice 1's proof used a fabricated layout with the binary present.

## 2026-07-03 — review findings live only in chat; a mid-implement disconnect loses them  [open]
- **What happened:** `/bench-review-implementation` on go-status-port returned two
  verified Coverage findings. The reviewer asked where findings are stored and the
  answer was nowhere — they exist only in the conversational output. Specs and
  decisions are persisted to disk, and `.bench/learnings.md` holds process captures,
  but the three-axis findings themselves have no durable artifact. If the session
  drops between review and the fix landing, the findings (and the refutation trail
  behind them) are lost and must be regenerated from scratch.
- **Right behavior:** When a review surfaces findings the reviewer intends to act on,
  persist them before starting the fix — so a disconnection mid-`/bench-implement-spec`
  leaves a pick-up surface, not a blank slate. A short findings file (axis, citation,
  confirmed/refuted, worst-issue) is enough; it is throwaway once the fixes land green.
- **Proposed rule change:** Consider having `/bench-review-implementation` write its
  aggregated findings to a durable location (e.g. `reviews/<spec>.md` or appended to
  the spec under a `## Review findings` heading) whenever any axis is non-empty and a
  fix path follows — the same "document for the teammate who just walked in" invariant
  the specs and ADRs already obey, applied to the one workflow output that currently
  evaporates. (Reviewer's call on the location and whether it auto-retires with the spec.)

## 2026-07-03 — hand-rolled a "TOON" emitter instead of surfacing the official library as an option  [open]
- **What happened:** The go-axi-query-surface slice ported the shell's TOON emitter
  as `internal/toon`, faithfully preserving the shell's byte format. The reviewer
  asked why the official Go library (github.com/toon-format/toon-go) wasn't used
  instead. Checking the TOON spec revealed the shell-inherited format was never
  spec-conformant TOON: it doubles quotes CSV-style where the spec uses JSON
  backslash escapes, and it skips the spec's mandatory quoting cases (empty string,
  true/false/null, numeric-looking strings, colons, leading hyphen). The port
  carried a private dialect forward under the TOON name without flagging it.
- **Right behavior:** When a build touches an implementation of a named external
  format, check whether an official library exists and whether the current
  implementation conforms to the published spec — and surface "adopt the library
  vs. keep the dialect" as a reviewer decision, since it trades a dependency
  against byte-stability of contract-pinned output. Byte-identity with the shell
  was the right port constraint, but the nonconformance was discoverable and
  should have been reported then.
- **Proposed rule change:** In /bench-write-spec, when a spec names an external
  format or protocol, add an edge-inventory check: does an official implementation
  exist, and does ours conform? Divergence is a decision to surface, not silently
  preserve.

## 2026-07-03 — a byte-compat research asset needs a runnable probe, not an API table  [open]
- **What happened:** `decisions/assets/go-toon-library.md` (research #5) resolved
  the TOON-library adoption with an API/compatibility table and #6 promised the
  block shape `name[N]{fields}:` was "unchanged". Writing the spec, a live probe of
  `toon-go` surfaced two behaviors the table missed: the library drops `{fields}`
  for an empty array (`name[0]:`) and omits the trailing newline. Both would have
  broken existing empty-table and runtime contracts if the adapter delegated
  naively; the spec added an empty-shim and a newline-append to hold #6's promise
  true. (I also caught that `IsSpace` is not fully superseded — it has a live
  non-emitter consumer.)
- **Right behavior:** The research phase should have run the library against the
  kit's own edge outputs (empty table, trailing newline, numeric cell), not just
  read its API and README. The gaps were cheap to find with one probe program.
- **Proposed rule change:** In the Research ticket type (and /bench-write-spec when
  consuming a research asset that claims byte- or wire-compat with an external
  library), require a runnable probe exercising the caller's own edge cases as the
  evidence for a "compatible/unchanged" claim — an API table alone does not settle
  byte compatibility.

## 2026-07-04 — Codex subagent policy conflicts with review-axis delegation  [open]
- **What happened:** `/bench-implement-spec` reached the required
  `/bench-review-implementation` step, whose phase text says to spawn three
  read-only axis delegates. The available Codex subagent tool's contract says
  not to spawn subagents unless the user explicitly asks for subagents,
  delegation, or parallel agent work, so the review was run inline instead.
- **Right behavior:** Harness-specific tool policy wins at runtime, but the
  Bench phase should not silently assume a delegation surface that a harness may
  forbid.
- **Proposed rule change:** In `/bench-review-implementation`, add a fallback:
  when the harness forbids unsolicited subagents, run the three axes inline,
  state the deviation in the exit report, and keep the same per-axis charges and
  citation standard.
