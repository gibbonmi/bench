# After-Merge Owner (FT25)

## #1: Which phase owns the post-merge tail?

Blocked by: —
Type: Grill

### Question
Post-merge duties (spec retirement, roadmap-row removal, decision-map sweep,
handoff refresh) are scattered onto phases reached only for other reasons.
`bench status` computes the duty list; no command consumes it. New phase, or
extend an existing one?

### Answer
Extend `/bench-final-check` with an exit duty rather than adding a phase: after
landing on green, when the session is on the default branch, read `bench
status` and run the housekeeping rows it flags — spec retirement
(`bench spec retire` + the `spec-retire:` commit), worktree/branch cleanup
(`bench worktree clean`), orphaned-pickup disposal. Roadmap and drain rows stay
with `/bench-what-next` (the deep pass). When final-check closes on a topic
branch, the duty defers to the next default-branch session, whose SessionStart
status shows the same rows — the duty text says so, so the tail is owned in
both flows. Rejected: a new after-merge phase (a whole phase for duties status
already enumerates; final-check is the last phase every build reaches).

## #2: What is the decisions/ lifecycle end?

Blocked by: #1
Type: Grill

### Question
19 maps accumulate, most for long-shipped work. The status decisions signal
counts only unresolved maps, so a shipped map is invisible but permanent —
asymmetric with specs' promote-then-delete. Sweep, or keep?

### Answer
Symmetric promote-then-delete: a closed map whose work has shipped is deleted
in the same pass that retires the spec, after promoting any decision not yet
recorded in docs/an ADR/the profile. History lives in git — same rationale as
spec retirement. The FT25 build includes the one-time sweep of the current
shipped maps (promotion check first), and "decision map" joins CONTEXT.md's
core terms. **Flagged for veto:** the alternative — keeping maps as a permanent
decision archive — was rejected because ADRs are the durable decision record
and two archives drift; the one-time sweep deletes ~17 files.

## #3: What happens to HANDOFF.md?

Blocked by: #1
Type: Grill

### Question
HANDOFF.md is a stale cold-session trap: no command reads, refreshes, or
expires it, and BENCH.md already states cold pickup must not depend on it.
Regenerate it at session end, or delete it?

### Answer
Delete the file and don't replace it. Cold pickup is the SessionStart status
line + ROADMAP.md + CONTEXT.md + the profile — all maintained by owned
mechanisms; a hand-written handoff document has no owner and rots by design.
Any prose that is still true and durable moves to CONTEXT.md before deletion.
The npm `files[]` entry is removed by the first-install build ([FT26]).
**Flagged for veto:** deletes a repo-root file the reviewer may be reading by
habit.

## Handoff

1. **Module boundaries.** `.agents/commands/bench-final-check.md` owns the exit
   duty prose; `CONTEXT.md` owns the "decision map" term; `decisions/` sweep is
   a one-time content change; no Go code changes (status rows already exist).
2. **Contracts.** Final-check exit: on green landing on the default branch,
   consume status housekeeping rows (specs sev 7, reviews sev 8, worktree
   sev 2); on a topic branch, state the deferral explicitly. Map retirement:
   promote-then-delete in the spec-retire commit.
3. **Deep vs thin.** All prose; the only depth is the duty wording — it must
   route rows to owners without restating what-next's job.
4. **Black-box assertables.** Conformance doc checks can pin the duty text in
   bench-final-check.md; the map sweep is assertable as `decisions/` containing
   only maps for unshipped work; CONTEXT.md contains the term.
5. **Gate attachment.** Docs conformance (stale-reference scans) sees the
   command edits; the sweep itself is a one-time reviewed diff, not gated
   behavior.
6. **Hostile-input owners.** n/a — prose-only change, no input surface.
7. **Uncertainty flags.** Whether any existing map holds unpromoted durable
   decisions — the sweep pass must check each map before deleting, not bulk-rm.
8. **Rejected alternatives.** New after-merge phase; regenerating HANDOFF.md at
   session end; keeping decisions/ as a permanent archive.
9. **Domain watch-outs.** The retirement and roadmap-reconcile status signals
   fire only on the default branch by design — the duty prose must not promise
   them on topic branches.

Dependency order: n/a — single spec; the decisions sweep rides in the same
build.
