# RepoResident distillation — drain notes

Research note behind six `ROADMAP.md` rows: FT101 (§8), FT106 (§1 and §5), FT107 (§3),
FT108 (§2), FT109 (§4), and FT110 (§6). It was written to let the 2026-07-23 drain
verdict eight parked `capture/IDEAS.md` lines without re-reading the source repo, and it now
carries the design reasoning and rejected alternatives those rows point at rather than
inline. Delete it when the last referencing row ships. Section 7 has no row: the
`.bench/BENCH.md` Capture section already routes any tangent worth not losing to
`bench idea`, which was the gap §7 asked about.

Source: `github.com/ychamel/RepoResident` @ `10ff7bd` (2026-07-23), MIT. The whole repo
is 840 lines of Markdown with no runtime: `CLAUDE.md` (operating manual, always loaded),
`AGENTS.md` (3-line pointer), and `.agent/` holding `STATE.md`, `MAP.md`, `PROJECT.md`,
`DECISIONS.md`, `ISSUES.md`, six workflow files, a design template, an area template,
and an append-only journal. It has no oracle, no enforcement, no isolation, no model
routing, and no install lifecycle — everything is instruction the model may ignore. The
ideas below are the parts that survive being ported onto a kit that *does* have a gate.

Rejected outright, so no drain time is spent on them: `MAP.md` + `areas/` as
hand-maintained navigation (bench computes that with `bench outline` / `bench status`,
and computed can't rot); the L0–L4 context taxonomy as such (`.bench/BENCH-reference.md`
already does lazy-by-path); the role list (bench has roles in `.bench/BENCH.md`); the
README's ">50% lower token usage" claim (self-reported, no methodology).

---

## 1. Doc staleness probes

**What RepoResident does.** `.agent/workflows/maintain.md` runs every 10th session and
on demand. Two steps matter. First, a hard budget table (`wc -l` per file, with a named
remedy per overflow). Second, the probe: *pick 2 area docs and 2 map lines at random,
open the code they describe, fix the lies — and if you found any, probe 2 more of each.*
The escalating sample is the whole trick: a clean tree costs four file reads, a rotten
one keeps paying until it stops finding rot.

**The bench gap.** `bench structure` budgets code and the gate enforces it. Nothing ever
re-verifies a *doc* against the tree. `CONTEXT.md`, `projects/<name>.md`, and `docs/adr/`
are asserted once and trusted forever, which sits badly with invariant 3 — we tell every
session to write for the teammate who just walked in, and never check that what that
teammate would read is still true. `/bench-drain` already reconciles `ROADMAP.md`
rows against the tree; this is the same move aimed at a wider target.

**How I'd land it.** A step in `/bench-drain`, not a new phase — the drain is
already the scheduled maintenance surface and already has a batch-diff verdict
mechanism to hang corrections on. Sample two claims from the doc set
(`CONTEXT.md` terms, `projects/<name>.md` seams, ADR statements), verify each against
the code, correct in the batch diff, and escalate the sample on any hit.

**How I'd improve it.** Three changes:

- **Don't sample randomly — sample by staleness.** RepoResident has no way to know which
  docs are suspect, so it rolls dice. We do: `git log` gives last-touched dates for both
  the doc and the code it describes. A doc older than the code it claims to describe is
  the candidate set; random selection within that set. Strictly better and costs one
  `git log` call.
- **Verdict, don't silently fix.** A probe hit is evidence about our own process, so it
  should land in the batch diff as a visible correction with a one-line why, the same
  shape as a learnings verdict. Silent doc edits inside a maintenance pass are how a
  doc's authority erodes without the reviewer noticing.
- **Keep the escalation rule verbatim.** It is the one part that makes the cost
  self-scaling, and it is cheap to state.

**Open question for the drain.** Whether the probe target list is hardcoded (`CONTEXT.md`,
`projects/*.md`, `docs/adr/`) or reviewer-declared alongside `.bench/structure.budgets`.
Recommend hardcoded first — a declaration file is a second thing to keep current.

---

## 2. Refactor lane with a mechanical exit test

**What RepoResident does.** `.agent/workflows/refactor.md`, for structure change with
zero behavior change. Four rules: (1) tests covering the affected behavior are green
*before* any move — if uncovered, write characterization tests first that assert what
the code currently does, "even where that's ugly"; (2) plan an ordered list of mechanical
moves, each leaving the repo green; (3) execute one move at a time, enumerating every
caller by grep, never from recall; (4) **the exit test** — the full suite must pass with
test logic unmodified, mechanical renames inside tests being the only permitted test
edit. *"Had to change an assertion? You changed behavior → stop, revert that move,
reroute to the feature path."* It also forbids bundling a bug fix into a refactor diff.

**The bench gap.** There is no refactor path at all. `rg -i refactor` across the kit
returns one line in `bench-implement-spec.md` about dry-running broad renames. Today a
pure restructure either gets forced through spec → implement → review (wrong shape: there
are no stories and no new observable behavior, so the coverage map is a fiction), or it
goes through the direct fix-and-gate path with nothing but the gate protecting behavior.

**How I'd land it.** A `craft-refactor` skill, not a phase. Phases are reviewer-chosen
entry points; a refactor is usually a shape the work turns out to have, and skills are
where generation-shaping guidance lives. It composes `craft-seams` (a refactor's whole
purpose is usually reaching a better seam) and `craft-tdd` (characterization tests are a
TDD variant with the red signal running backwards).

**How I'd improve it.** The assertion rule is stated as a discipline; on bench it can be
close to mechanical. A refactor commit's diff over test files should contain no changed
assertion lines — that is checkable, and it is exactly the class of thing
`bench structure` already does (parse the tree, apply a declared budget, fail the gate).
Worth exploring as `bench diff --assertions` or a gate fragment, so "no behavior change"
becomes an oracle verdict instead of a promise. Start with the skill; propose the check
as a follow-on row only if the skill gets used and the rule gets violated.

Also worth keeping: the no-bundling rule (found a bug mid-refactor → park it, fix
separately). That is invariant 4's smallest-diff rule pointed at a specific failure mode
we don't currently name.

---

## 3. Named lighter lanes with stated escalation triggers

**What RepoResident does.** `CLAUDE.md` holds a routing table: request shape → one
workflow file, which is the only workflow the session opens. The lanes are feature,
patch, debug, refactor, review, maintain. Each lane states its own boundary in
mechanical terms rather than by judgment. `patch.md`: cause known, no interface change,
≤2 files, and it escalates *the moment* an interface or schema must change, the fix wants
more than 2 files, or ~5 files have been read without naming the cause with evidence
(that last one reroutes to debug, not feature). `feature.md` carries the mirror test:
"single known site, no interface change, obvious verification → downgrade to patch." The
manual adds a tiebreak: when on the fence, start with the lighter workflow.

**The bench gap.** Right-sizing is prose in `.bench/BENCH.md`: a few-line change doesn't
need the full pipeline, you may propose a lighter path, but you must get an explicit OK
before skipping canonical steps. So every small change costs a round trip, and the
reviewer answers the same question repeatedly with the same answer. The guide already
anticipates the fix — *"If I give you a standing rule for changes of a given size, follow
it and stop asking"* — but no standing rules have ever been written down.

**How I'd land it.** Write the standing rules as a small table in `.bench/BENCH.md`
beside the proportionality paragraph: change shape → the path it takes → the trigger that
forces escalation back to the full pipeline. This creates no new phases and no new files;
it converts a permission ask into a rule the reviewer set once, which is what the
existing clause already licenses.

**How I'd improve it.** Two changes:

- **Bound the lane by blast radius, not file count.** "≤2 files" is a proxy that fails on
  a monorepo and on generated code. bench has the better vocabulary already: a change
  that crosses a declared seam takes the full pipeline; a change strictly inside one is a
  candidate for the light path. That also makes the trigger legible to the profile —
  `projects/<name>.md` already names the seams.
- **Keep the "5 reads without a named cause" trigger verbatim.** It is the best line in
  the file: it catches the specific failure where a session is patching by hypothesis and
  should have opened `/bench-debug`, and it triggers on an observable (reads spent), not
  on self-assessed confidence.

**Do not adopt** RepoResident's own lane set wholesale — bench's phases are the lanes.
This is about writing the escalation boundaries, not importing six workflow files.

---

## 4. Handoff discipline: rewritten in full, capped, git wins

**What RepoResident does.** `.agent/STATE.md` is the single source of "now", capped at
40 lines, with fixed sections (Session, Focus, Active, Next, Blocked, ≤5 watch-outs, ≤3
recently-shipped) and a template embedded in the file itself. Three rules give it teeth:
it is **rewritten in full** at every session close, never appended to; the header says
*"prune, never accrete: this file is read by every session, so every stale line here is a
tax on all future work"*; and the conflict rule — *if this file contradicts `git log` or
the journal, a session died before close: trust git, rebuild this file from the last
journal entry plus `git log -5`, and note the crash.* Detail and history live in the
append-only journal, which is grepped, never loaded. Mid-session, only the one-line
`Active:` field is updated, explicitly as crash insurance.

**The bench gap.** `capture/session-handoff.md` is currently a 66-line free-form narrative with
no template, no cap, no rewrite rule, and no conflict rule. The `AGENTS.md` phase-close
rule says the closing message must emit a continuation prompt *or* update the handoff —
so the file's shape is entirely up to whoever writes it, and nothing says what happens
when it disagrees with the tree.

**How I'd land it.** Put the template inside `capture/session-handoff.md` itself (RepoResident's
trick — the template can't drift from the artifact if it *is* the artifact) and add the
rewrite-in-full plus conflict rules to the phase-close handoff paragraph in `AGENTS.md`.

**How I'd improve it.** The conflict rule can be made computed rather than remembered.
bench already resolves a gated tree hash (`internal/git`, `bench gate pin`) and reports
staleness in `bench status`. A handoff that records the commit it was written at can be
checked against `HEAD` automatically — `bench status` reports "handoff written at
`<sha>`, HEAD is `<sha>`: stale" and points at the rebuild. That turns "trust git" from a
rule the next session must recall into an ambient fact the SessionStart hook already
prints. Cheap: the handoff gains one machine-readable line, `session-inspect` gains one
comparison.

Keep the cap. The reason RepoResident gives is the right one and applies harder to bench,
where the handoff is read cold by a fresh session that is paying for every line.

---

## 5. The `(?)` unverified marker

**What RepoResident does.** A one-character convention for facts written from inference
rather than verification. `bootstrap.md` (adopting an existing repo) marks every inferred
map line `(?)`, and instructs that anything unverifiable be marked `(unverified)` rather
than stated confidently — the workflow's stated reason is that *wrong facts written at
bootstrap poison every future session*. The marks are discharged lazily: the next session
to visit that code verifies and strips the mark, and `maintain.md` treats any remaining
`(?)` as a violation to clear ("no `(?)` remaining" is a checklist item).

**The bench gap.** `/bench-setup-repo` interviews the reviewer for the load-bearing
choices (gate, seams, lines), which is exactly right for those. But an adopting session
still writes plenty into `projects/<name>.md` and `CONTEXT.md` from inspection, and once
written those read identically to reviewer-confirmed facts. Nothing distinguishes "the
reviewer told me this seam is canonical" from "I inferred this seam from the import
graph."

**How I'd land it.** A convention line in `craft-adr` (which owns doc-writing discipline)
plus a marking instruction in `/bench-setup-repo`'s exploration half. Two markers, not
one: `(?)` for inferred-and-checkable, `(unverified)` for asserted-and-not-currently-
checkable.

**How I'd improve it.** Wire it to idea 1 — a `(?)` is a self-declared probe target, so
the staleness probe should drain marked claims before sampling unmarked ones. Together
they close a loop RepoResident leaves open: it marks uncertainty and hopes a session
wanders past, where we would have a scheduled pass that seeks the marks out. That
pairing is the reason both are worth taking; either alone is weaker.

Also: a `(?)` surviving several drains is itself a signal — either the claim doesn't
matter (delete it) or nobody can verify it (escalate to the reviewer). Consider a count.

---

## 6. Two generation-shaping one-liners

**What RepoResident does.** Prime directive 4: *"Never call an API/function you haven't
seen defined this session — read its definition first. Unsure how something behaves?
Verify in source, not from memory."* Prime directive 3, closing clause: *"Never read
directories wholesale. Lost after ~3 file reads? Re-scope from MAP or the area doc
instead of reading more code."*

**The bench gap.** Invariant 4 says "read the surrounding code before you write," which
is true but unfalsifiable — surrounding is undefined and a session always believes it read
enough. The first line names the exact artifact to read and the exact moment. The second
gives an exploration budget, which the kit has nowhere: `craft-seams` says where to attach
a test, nothing says when to stop searching and re-scope.

**How I'd land it.** The API line sharpens invariant 4 in `.bench/BENCH.md` — a clause,
not a new invariant. The exploration budget belongs in `craft-seams` beside the seam-
finding guidance.

**How I'd improve it.** RepoResident re-scopes to a hand-maintained map, which is its
weak link; we re-scope to `bench outline` and `bench status`, which are computed and
can't be stale. So the port is *stronger* on bench than in the original: "three reads
without traction → stop reading, run `bench outline`, pick the seam from the index." That
also creates real pull for `outline`, which is currently underused. Do not port the "3"
as a hard number — state it as a small fixed budget the session declares it has spent, so
the behavior is observable in the transcript rather than counted internally.

---

## 7. Where a technical out-of-scope observation goes

**What RepoResident does.** The journal has an optional `flag:` line — *"debt, risk, or
cleanup noticed but out of scope"* — appended at session close beside the outcome. It is
explicitly the informal inbox: `maintain.md` collects flags from the last two months,
promotes the still-true ones into `ISSUES.md` (urgent ones also get a state watch-out),
and drops the dead ones during rollup. `ISSUES.md` states the reason bluntly: *"if 'we
should fix that later' isn't a line here, it will be forgotten."*

**The bench gap — verify before building.** `capture/learnings.md` is scoped to *process*
(deviations, judgment calls, should-have-askeds, codification candidates). `capture/IDEAS.md` is
the reviewer-facing idea inbox. Neither obviously owns "I noticed this unrelated path is
racy while working on something else." `bench idea` is probably the intended home and may
be entirely sufficient — in which case this drains as a one-line clarification in
`.bench/BENCH.md`'s Capture section, not a new mechanism.

**How I'd land it, if the gap is real.** Say so in Capture: technical observations go to
`bench idea` with a `debt:` prefix, and the drain routes prefixed lines to the roadmap's
debt tier rather than the feature sequence. No new file. bench already has the strictly
better half of RepoResident's loop — `/bench-drain` verdicts every entry in a
reviewed batch diff, where `maintain.md` lets the agent promote and drop flags on its own
authority, which violates our rule that the agent captures and the reviewer decides.

**Recommend checking first.** This is the weakest of the eight and may be a no-op.

---

## 8. Bounded context, scoped for a monorepo — merged into FT101 (2026-07-23)

**Not a drain candidate.** The reviewer merged this into `ROADMAP.md` **FT101** on
2026-07-23 and dropped its `capture/IDEAS.md` line, so this section is background for that row.
FT101 widened from "multi-context domain docs for monorepos" to "per-context scope:
domain docs and profile" — the `CONTEXT.md` half it already had, plus the profile half
below. Read on for the reasoning, and for the gate-scoping pushback recorded at the end,
which is the constraint the row now carries.

**What RepoResident does.** Its stated principle is *"context cost follows the task, not
project age"*: hot files are capped and rewritten, history is greppable but never loaded,
and deep knowledge lives in `areas/<name>.md` (cap 60 lines) that a session reads
*instead of* re-deriving a module from source. `MAP.md` (cap 120 lines) collapses a
subtree into a single pointer line at its area doc when it overflows. An area doc is
created only on a specific signal — *a session had to re-derive something non-obvious* —
and its sections are chosen for what re-derivation costs: mechanism, key files,
invariants, gotchas ("why the obvious approach fails here"), and the sanctioned pattern
for the 2–3 most likely changes.

**Why the layering itself isn't the import.** bench already defers cost by path:
`.bench/BENCH-reference.md` is referenced, never imported, and the ambient surfaces are
computed on demand. Porting an L0–L4 taxonomy on top would add vocabulary without adding
behavior, and porting `areas/` as hand-written module docs would import exactly the rot
problem idea 1 exists to patch.

**What *is* worth importing is the scoping unit — and a monorepo is where it bites.**
bench's per-repo knowledge unit is `projects/<name>.md`: seams, gate command, line
assignments, cold-session notes. That shape assumes one profile per repository. In a
monorepo every one of those facts is per-package: package A's seams are not package B's,
the meaningful gate for a change inside A is A's suite rather than the whole tree, and A's
line (a mechanical CLI wrapper, cheap) is not B's (an uncertain domain model, top). One
profile for the tree forces the union — every shift pays the whole-tree gate and reads a
seam list mostly about code it will never touch. The computed surfaces degrade the same
way: whole-tree `bench outline` and `bench structure` get less useful as the tree grows,
which is the precise failure RepoResident's principle names.

**How I'd land it.** Extend the existing profile concept rather than adding a parallel
docs tree:

- **Path-scoped profiles.** Allow `projects/<name>.md` to declare the paths it owns, and
  resolve the active profile from the paths a change touches. One profile keeps working
  unchanged, so this is additive.
- **Scoped gate invocation.** The gate stays the oracle and stays reviewer-owned — this
  narrows *which* gate runs for a scoped change, it never weakens one. Precedent exists:
  `bench structure` already distinguishes whole-tree from touched, and `bench diff`
  already resolves a review base. The hard case is a change that touches two profiles;
  the safe default is the whole-tree gate, and it should be stated rather than inferred.
- **Scoped ambient context.** `bench outline` and `bench status` take the resolved
  profile's paths, so a cold session in a monorepo gets its package's map rather than the
  tree's.

**How I'd improve on RepoResident here.** Its area docs are hand-written and capped by an
advisory line count that only a maintenance sweep checks. Ours would be a profile —
already reviewer-owned, already the thing `/bench-setup-repo` interviews for — with the
budget enforced by `bench structure` in the gate, which is where bench is genuinely
stronger than the source. And where RepoResident *describes* a module in prose, we point
at computed surfaces for the same information, so the scoped context can't rot; the
profile carries only what can't be computed (the seams, the gate, the line, the gotchas).

The one prose element worth copying verbatim is the **creation trigger**: write the doc
only when a session actually had to re-derive something non-obvious. That is what keeps
this from becoming a documentation chore, and it is the discipline `craft-skills`
(anti-sediment) already applies to skills.

**Merge, as executed.** FT101 already covered the `CONTEXT.md` half — a root
`CONTEXT-MAP.md` pointing at per-context `CONTEXT.md` files, plus a single-versus-multi
question in `/bench-setup-repo` Section C. This idea was the same problem from the profile
side (seams, gate, and line rather than vocabulary), and the two wanted one resolution.
The merged row keeps both halves, adds the additive-resolution property (a single-profile
repo resolves exactly as today), and carries the gate-scoping question as a grill rather
than as settled work.

## 8a. The gate-scoping pushback — the constraint FT101 carries

The reviewer's motivation for scoping was gate wall-clock. I disagreed on that framing,
and FT101 is written to hold the disagreement rather than bury it. Three reasons, in
descending strength:

**Scoping does nothing for this repo.** Bench is not a monorepo. FT91's own post-host-only
measurement (2026-07-22) puts the 10–15 minute gate down to *contention*, not tree size:
the canary phase nests whole gate runs, oversubscribing a 16-core box to load average
~123, and `internal/gate/phases.go` hardcodes `-count=1`, disabling Go's test cache
unconditionally. Those are the levers. Path scoping would leave both untouched and the
wall clock roughly where it is, so adopting it as a speed fix here would spend a build and
measure no delta.

**Diff-scoped gating is already closed, and this must not reopen it.** FT91 ruled it
unsound for this tree because contract and canary are behavior contracts with no file→test
map — there is no sound rule mapping a changed file to the subset of contracts that can
still fail. Nothing in the monorepo case changes that. What *is* different is the source
of the boundary: a profile's declared paths are a reviewer decision, where a diff is an
inference. So a scoped gate is legitimate only where the reviewer declared the boundary,
and a change touching two profiles takes the whole-tree gate. That distinction is the
entire load-bearing difference between this row and the thing FT91 rejected, and it is
worth restating whenever the row is discussed.

**Speed as the justification is how the oracle erodes.** Invariant 1 gives the gate its
authority precisely because nothing is allowed to narrow it for convenience. Once
wall-clock is an accepted reason to run less of it, every future narrowing arrives with
the same argument and no principle to stop at. Scope has to earn its place as a
*correctness and context* feature — the right seams, the right vocabulary, and ambient
surfaces that stay useful as a tree grows — with any wall-clock win treated as a monorepo
side effect. FT91's own constraint says the same thing from the other side: green must
keep meaning what it meant, and a scoped verdict is explicit evidence, never a silent
skip.

If speed is the actual goal, the honest row is FT91, and its arms are already named.

**Verdict.** Worth building as monorepo profile scoping — a real capability gap with a
clear seam and a precedent (`bench structure` already splits whole-tree from touched,
`bench diff` already resolves a review base) — not as a context-layer taxonomy, and not as
a gate-speed lever. Sequence after idea 1; the staleness probe is what keeps scoped
profiles honest once there are several of them.
