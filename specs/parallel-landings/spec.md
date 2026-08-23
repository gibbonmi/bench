# parallel-landings

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-23 — the reviewer named parallel landings as the goal, parked three ideas (the spec-less landing, the phase-owned file merge rules, all work in a worktree), and authorized immediate action; `roadmap/FT169.md` holds the undecided authority half this spec settles in part

Verification log: 1 iteration(s) to accept — the round (opus / high, read-only) returned REVISE with six blocking findings and two ticket-graph findings; the author folded every finding (the false observed red on the journals, the light-path contradiction, the anchors fence, the base-range and resume rows, the ticket edges) and left the tickets-only close as the flagged reviewer decision; the declared cap was one iteration plus one fold, so no second full pass ran

## Problem

One session can land at a time. A second source that moves `main` under the
first collides on the phase-owned files. Every session writes those files: the
handoff, the learnings journal, the ideas inbox, and the board. Today the
landing settles a conflicted journal by the destination side, so the source
session's entries vanish silently. A board conflict refuses without a named
repair.

The landing verb also needs a staged `spec.md` with an ownership fence. A
drain, a decision map, a bug fix, or a light-path ticket has none, so that work
lands by hand on `main`. Hand landings on `main` are the writes that move
`main` under everyone else.

## Solution

`bench worktree land` accepts a source with no spec. It composes, gates,
publishes, and releases that source exactly as it does a spec-backed one. It
skips only the spec transition and the fence check. A `--spec` that names a
tickets-only folder closes that folder in the published tree, so a light-path
ticket retires through the verb.

The composition settles a conflicted phase-owned file by a rule table. The
source wins the handoff. The two append-only journals compose as the union of
both sides. The destination wins every other capture file. Every other
conflict refuses, names each conflicted path, and names the source repair.

With those two in place, the workflow runs every phase in a bench worktree and
lands it through the verb. `main` receives writes only through landings. Merge
composition stays the primitive; no rebase enters the workflow.

Closed decisions, dated 2026-08-23: merge composition onto the moved `main`
stays the landing primitive, and rebase is rejected because it rewrites the
reviewed tip. `--spec` becomes optional. A tickets-only `--spec` closes its
folder on the landing (the FT224 reviewer decision, taken here; flagged for
veto). The journal union and the destination default are flagged for veto.
The "main is written only by landings" rule ships as guidance, not as a hook.

## User stories

### Group A — a source with no spec lands through the verb

Line: opus / medium. The landing command and the composition owner are the
known seams, the spec is exact, and the gate covers the path; mid because the
surface grants landing authority.

1. As an operator, I want `bench worktree land` to accept a source with no `--spec`, so that spec-less work lands through the verb.
2. As an operator, I want a spec-less landing to compose, gate, publish, mark, reconcile, and release like a spec landing, so that authority holds.
3. As an operator, I want a spec-less landing to skip the spec transition and the fence check, so that a path no fence names lands.
4. As a reviewer, I want a spec-less landing to still require the exact `--base` and `--source-tip` pair, so that the commit pins a reviewed range.
5. As an operator, I want `--resume` to accept a published spec-less landing, so that an interrupted spec-less landing finishes the same way.
6. As an operator, I want a tickets-only `--spec <slug>` to close its folder in the published tree, so that a light-path ticket retires through the verb.
7. As an operator, I want a first-run `--spec` that names a `spec.md` to keep the transition and the fence check, so that spec-backed landings hold.
8. As a script author, I want the `landed{...}` record and the exit codes unchanged for a spec-less landing, so that readers keep their meanings.
9. As a cold-session agent, I want the optional spec stated in the usage line and in `.bench/BENCH-reference.md`, so that the grammar is discoverable.

### Group B — a conflicted phase-owned file composes by rule

Line: opus / medium. The conflict path of the composition owner is the known
seam, and the composition tests cover it; the union merge has one plumbing
uncertainty (`merge-file` over three blobs), so mid rather than cheap.

10. As a landing operator, I want a conflicted journal (`capture/learnings.md`, `capture/IDEAS.md`) composed as the union of both sides, so that no appended entry is lost.
11. As a landing operator, I want the source to keep winning `capture/session-handoff.md`, so that the closing session's state is the handoff.
12. As a landing operator, I want the destination to keep winning every other conflicted capture file, so that the running ledgers on `main` stay authoritative.
13. As a landing operator, I want a union path that one side deleted to take the present side, so that a deletion cannot erase entries.
14. As a landing operator, I want a conflict outside the rules to refuse and name every conflicted path, so that repair starts from the path.
15. As a landing operator, I want that refusal's `next=` to name the source repair and the re-run, so that recovery needs no archaeology.
16. As a reviewer, I want every rule resolution disclosed on the landing's stderr, so that what the merge did not decide is visible.
17. As a reviewer, I want a non-regular conflicted path under `capture/` to refuse, so that the rule table never rewrites an object kind.

### Group C — every phase runs in a worktree and lands through the verb

Line: fable / high. This group edits kit guidance prose, so the leverage
override in `craft-line` applies.

18. As a session, I want every phase to run in a bench worktree and land through the verb, so that only landings write `main`.
19. As a session, I want the phase-close handoff written in the worktree and landed with it, so that handoff and merge rules agree.
20. As a debug session, I want the isolation rule to name the spec-less landing, so that a bug fix isolates and lands like a build.
21. As a reviewer, I want the guidance to keep merge composition as the primitive and to reject rebase, so that the review identity survives.
22. As a reviewer, I want the only-landings-write-`main` rule to be guidance, not a hook, so that `bench commit` keeps working on any branch.

## Implementation decisions

**The spec is optional on both landing grammars.** `--spec` loses `Required`
on the first-run and the resume grammar; an empty value stays a usage error.
Without it, the source proof skips the spec resolve, the transition, and the
fence authorization, and the landing request carries no spec path. Every
other proof runs unchanged: destination cleanliness, assignment, owner marker,
identity, cleanliness of the source, and the range resolution. The resume path
skips only the published-spec authentication when no `--spec` is given.

**A tickets-only `--spec` closes on the landing.** The landing owner already
knows a tickets-only folder and composes its removal for `bench commit`. The
reviewed landing reuses that same removal on the composed tree and removes the
folder from the destination checkout at reconcile. A folder that the
destination already removed composes as a no-op. This is the FT224 reviewer
decision, taken here and flagged for veto.

The close is a precondition of the worktree rule. A light-path ticket lives in
a tickets-only folder, and the folder must close on the landing or it stays
orphaned on `main`. If the reviewer vetoes the close, story 6, WL8, and the
close ticket move to Out of scope. The light path then stays exempt from the
worktree rule.

**The light path joins the worktree rule.** `.bench/BENCH.md`'s right-size
table says a light-path ticket lands inline with no worktree. Under this spec
the light path runs in a worktree too and lands through the verb with its
tickets-only `--spec`. The guidance ticket rewrites that table row, and the
anchor needle that pins it moves with it.

**One rule table for the phase-owned paths.** The composition owner's
`CaptureSide` becomes a rule table with three verbs: `source`, `destination`,
and `union`. `capture/session-handoff.md` takes `source`. `capture/learnings.md`
and `capture/IDEAS.md` take `union`. Every other `capture/` path takes
`destination`. A path outside the table has no rule, so the conflict refuses.

The union is Git's own three-way union over the three stage blobs; with one
stage absent, the present side's blob is the result. Only regular-file stages
qualify; any other object kind under the table refuses with its kind.

**The conflict refusal names the repair.** The refusal keeps its `refused{}`
record and its paths table and gains `next=`. The value names the source
repair in order: merge the destination commit into the source worktree, commit
the repair, review the new range, and re-run the landing. A value that is not
line-safe takes the pointer form, as every `next=` already does. No Bench verb
moves a retained worktree onto the destination yet (FT238), so the merge step
names raw Git; the guidance names that gap.

**Guidance, not enforcement, for the worktree rule.** The workflow text in
`.bench/BENCH.md`, the landing paragraph in `.bench/BENCH-reference.md`, the
phase-close handoff rule in `AGENTS.md`, and the debug skill's isolation rule
change. No hook refuses a commit on the default branch. The `.bench/BENCH.md`
budget is 180 lines and its anchored needles are registered; the guidance
ticket keeps the budget and repairs any needle it moves.

**Bootstrap authority.** The spec-less landing removes a path authorization,
not an executable hop. The executable chain is unchanged. `bin/bench.sh`
resolves `dist/bench`. The landing verifies that executable's seal against the
source digest, or it rebuilds it through the sanctioned build and re-runs.

The gate's private exact-source build grades the composed tree. The
destination CAS publishes. The seal is self-attestation, and the sanctioned
rebuild is the independent root the profile names. A repository that declares
no build inputs skips the seal check, as today; its trust root is the gate's
private exact-source build alone.

## Testing decisions

- A good test drives `bench worktree land` against a fixture repository and
  asserts the output record, the published tree, the refs, and the exit code.
- Seams with prior art: the landing command seam (`land_test.go` in-process
  calls, `land_surface_test.go`, and built-binary journeys through
  `publicLandingFixture`), and the composition seam (`composition_test.go`
  drives `Compose` against real Git conflicts without mutation).
- The gate observes through its `test` phase (`go test -count=1 ./...`); no new
  gate phase.

### Seam diagram

    trigger: the operator or /bench-final-check runs `bench worktree land`
        │
        ▼
    flags + repo state ──▶ [ land command: proofs → composition → gate ] ──▶ refused{…} | landed{…}
                           [ composition owner: merge-tree → rule table ]     stderr: landing source{…},
                                                                              landing composition{resolved=…}
                      ◀ tests attach here: in-process command calls and
                        composition calls drive fixture repositories;
                        assertions read records, trees, refs, and exit codes

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| WL1 | 1 | a land without `--spec`, with a correct request, base, and tip, exits 0 with `worktree=released` | command | today the grammar refuses the missing flag with exit 2 |
| WL2 | 2 | the spec-less published commit has the destination as first parent and the source tip as second parent, and the green marker advances to it | command | reds a shortcut that squashes the source or skips the marker |
| WL3 | 2 | a spec-less landing whose composed tree fails the gate refuses and publishes nothing | command | the gate still owns the spec-less path |
| WL4 | 3 | a spec-less source whose commits touch a path no fence names lands | command | today that source refuses at the fence |
| WL5 | 3 | a spec-less landing publishes every `specs/` file byte-identical to the composition | command | reds an implementation that transitions a staged spec anyway |
| WL6 | 4 | a spec-less land with a `--source-tip` that differs from the worktree HEAD refuses with both identities | command | the identity proofs still run without a spec |
| WL7 | 5 | a spec-less landing interrupted after publication resumes without `--spec` and exits 0 | command | the resume path must not demand a spec the first run did not have |
| WL8 | 6 | a land with `--spec` that names a tickets-only folder publishes a tree without that folder and releases | command | today the verb refuses a tickets-only slug as an unreadable staged spec |
| WL9 | 7 | a land with `--spec` that names a staged `spec.md` still publishes `Status: implemented` | command | freezes the spec-backed transition against the optional flag |
| WL10 | 8 | a spec-less landing's `landed{...}` record carries the same fields and exit 0 | command | reds a record change on the new path |
| WL11 | 10 | a landing whose source and destination both appended distinct lines to `capture/learnings.md` publishes a file that holds both sides' lines | command | today the destination rule silently discards the source's appended lines |
| WL12 | 10 | the same shape on `capture/IDEAS.md` publishes the union | composition | today the destination rule silently discards the source's appended lines, and the table must name both journals |
| WL21 | 7 | a land with `--spec` that names a staged `spec.md` still refuses an out-of-fence path | command | freezes the spec-backed fence check against the optional flag |
| WL22 | 8 | a spec-less refusal still exits 1 | command | reds an exit-code change on the new path |
| WL23 | 4 | a spec-less land whose `--base` is not an ancestor of the source tip refuses before the gate | command | the range proof still runs without a spec |
| WL24 | 4 | a spec-less landing's `landed{source_base=...}` value is the resolved review base | command | reds an implementation that echoes the flag instead of the resolved base |
| WL25 | 5 | a `--resume` without `--spec` on a published spec-backed landing completes the marker and the release without a second publication | command | pins the narrowing the optional resume flag makes |
| WL13 | 11 | a conflicted `capture/session-handoff.md` publishes the source bytes | composition | pins the rule the FT169 fix shipped |
| WL14 | 12 | a conflicted capture file outside the named rules publishes the destination bytes | composition | pins the destination default |
| WL15 | 13 | a union path deleted on one side and appended on the other publishes the appended side's bytes | composition | reds a union that treats a missing side as a refusal or an empty file |
| WL16 | 14 | a conflict on `ROADMAP.md` refuses and names `ROADMAP.md` in the paths table | command | the refusal-with-path contract covers board files |
| WL17 | 14 | a conflict on a code path and a capture path together refuses and names both paths | composition | reds a rule table that settles the capture half and hides the code half |
| WL18 | 15 | the conflict refusal's `next=` names the source repair steps and the re-run invocation | command | today `next=` is absent on a conflict |
| WL19 | 16 | a landing that settled a journal by union prints one `landing composition{resolved=...}` line that names the path and `union` | command | today no resolution names `union`, and silent rule application is the cheapest wrong implementation |
| WL20 | 17 | a symlink or gitlink conflict under `capture/` refuses with the conflict kind and the path | composition | reds a rule table that dereferences or rewrites an object kind |

Not covered: story 9 — usage and reference prose; the review round verifies it.
Not covered: story 18 — workflow guidance prose; the review round verifies it.
Not covered: story 19 — handoff guidance prose; the review round verifies it.
Not covered: story 20 — debug skill prose; the review round verifies it.
Not covered: story 21 — guidance prose; the review round verifies it.
Not covered: story 22 — guidance prose; the review round verifies it.

Cheapest wrong implementation per group, and the row that reds it:

- make the flag optional and refuse later on the missing spec → WL1
- transition some spec anyway → WL5
- accept an unresolved range without a spec → WL23
- echo the `--base` flag as the landed source base → WL24
- demand the spec on resume → WL7
- take one side for the journals → WL11
- settle only the capture half of a mixed conflict → WL17
- settle silently → WL19
- dereference a symlink to make it regular → WL20

### Edge inventory

- A conflict record's path is split on the record's tab, never on spaces, so
  spaces and glob characters survive (under WL11, WL16).
- A conflicted path that carries a control byte renders through the sanitized
  paths table (under WL16). A `next=` that is not line-safe takes the pointer
  form (under WL18).
- Both sides add a journal absent from the merge base: the union of the two
  present stages (under WL11).
- A union path whose content Git's union merge cannot read as text refuses and
  names the path (under WL17's shape).
- A mode conflict on a phase-owned path refuses with kind `mode` (under WL20).
- `--spec ""` stays a usage error with exit 2 (unchanged grammar).
- `--spec` that names a folder absent from the source stays the existing
  unreadable-spec refusal.
- `--spec` that names a tickets-only folder the destination already removed
  composes as a no-op and lands (under WL8).
- `--resume` without `--spec` on a published spec-backed landing authenticates
  the parents and the range only. The same published commit completes either
  way, so no second publication can occur.
- **Won't handle:** a merge rule for `ROADMAP.md` and `roadmap/` — the conflict
  refuses and names the repair, and the drain session repairs by merge.
- **Won't handle:** a hook that refuses `bench commit` on the default branch —
  `bench commit` on `main` survives as the caller until the reviewer decides.
- **Won't handle:** a Bench verb that moves a retained worktree onto the
  destination — FT238 owns it. The named raw `git merge` repair survives as
  the caller.

## Ownership fences

- `internal/worktree/land.go`
- `internal/worktree/land_test.go`
- `internal/worktree/land_surface_test.go`
- `internal/landing/landing.go`
- `internal/landing/close.go`
- `internal/landing/composition_test.go`
- `internal/landing/landing_test.go`
- `internal/usage/worktree.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `.bench/BENCH.md`
- `.bench/BENCH-reference.md`
- `.agents/commands/bench-debug.md`
- `AGENTS.md`

## Out of scope

- A Bench verb that moves a retained worktree onto the destination commit
  (FT238): 2 edits, 2 gate runs.
- A hook that refuses `bench commit` on the default branch: 3 edits, 2 gate
  runs, after the reviewer's enforcement decision.
- A merge rule for `ROADMAP.md` and `roadmap/FT<n>.md`: 2 edits, 1 gate run,
  after a drain-shape decision.
- FT169's interrupted-landing recovery model and the hook pool's one-active
  limit stay on FT169.

## Further notes

The rule table's verbs are the vocabulary for any later phase-owned path: a
path gets `source`, `destination`, `union`, or no rule. The three parked ideas
close by implementation at the next drain.
