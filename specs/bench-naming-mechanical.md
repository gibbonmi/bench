# Bench command & skill renaming — mechanical slice

<!-- command-currency: historical -->

Historical implementation spec. It intentionally names pre-rename commands and skill
paths while documenting the rename that removed them.

Implements the mechanical renames decided in `decisions/bench-naming.md`. The two
behavior changes from that map (the `resynthesize` split, the capture wiring) are
out of scope here — this slice is pure rename + reference fixes.

## Problem

Bench's commands and skills have bare names (`spec`, `build`, `prep-shift`, `adr`,
`axi`, `seams`…) that don't say what they do and don't say where they sit in the
workflow. They also collide visually with the many non-Bench skills a harness loads
(`grill-me`, `tdd`, `prototype`…), so a user can't tell which skills are Bench's.

## Solution

Rename every Bench-owned command and skill into a `bench-*` namespace:

- **Commands → `bench-<verb>`** (an action verb, ordered by the workflow).
- **Skills → `bench-craft-<noun>`** (the `craft` infix marks "auto-firing know-how"
  vs "a phase I invoke").

After the rename, the user types self-describing phase names, and every Bench skill
is unmistakably Bench's. Nothing about *behavior* changes — only names and the
references that point at them.

The slate (from the decision map):

| from (command) | to | | from (skill) | to |
|---|---|---|---|---|
| start-ideation | `bench-ideate` | | adr | `bench-craft-adr` |
| spec | `bench-spec` | | axi | `bench-craft-cli` |
| build | `bench-build` | | design-system | `bench-craft-design-system` |
| fix-bug | `bench-diagnose` | | grill | `bench-craft-grill` |
| prep-shift | `bench-review` | | seams | `bench-craft-seams` |
| verify-gate | `bench-qa` | | tdd-at-seams | `bench-craft-tdd` |
| setup | `bench-setup` | | writing-great-skills | `bench-craft-skills` |

`resynthesize` is deliberately **left as-is** in this slice.

## User stories

1. As a user, I want each of the 7 commands renamed to its `bench-<verb>` form, so
   the phase I invoke is self-describing.
2. As a user, I want each of the 7 skills renamed to its `bench-craft-<noun>` form,
   so I can tell Bench's skills from the other loaded skills.
3. As the agent, I want each renamed skill's frontmatter `name:` to match its new
   directory, so the harness loads it under the new id. (Commands have no `name:` —
   the filename is the id.)
4. As a user, I want the `AGENTS.md` skills index (the `.agents/skills/<name>/SKILL.md`
   pointers) and its `/command` references updated, so the canonical working
   agreement points at files that exist.
5. As a user, I want the shared-rules doc `.bench/BENCH.md` workflow list updated to
   the new `/command` names, so the canonical phase list is accurate.
6. As a user, I want `README.md`, `HANDOFF.md`, `CONTEXT.md`, and `CLAUDE.md`
   command/skill references updated, so no published or onboarding doc names a
   command that no longer exists.
7. As the agent, I want every cross-reference *inside* command and skill bodies
   updated — especially `bench-craft-skills`' two example lists, which name the old
   commands and skills explicitly — so a skill that points at another skill or
   command resolves.
8. As a user, I want `bin/bench.sh`'s skill-pointer prose and `bench status` output
   strings updated to the new names, so the CLI's guidance and ambient-status rows
   point at findable skills.
9. As the agent, I want `.bench/gate.sh`'s hardcoded name-literals moved to the new
   names, so the oracle checks the renamed artifacts — without weakening any
   assertion.
10. As a user, I want only invocation/identifier forms renamed — not the English
    words "spec", "build", or "seam" used as plain concepts — so prose stays
    readable.
11. As a user, I want `resynthesize` untouched in this slice, so the behavior-change
    split lands in its own spec.
12. As a user, I want `bench gate` green and zero surviving old identifier forms on
    the active surface, so "done" is the oracle's verdict plus a clean sweep — not a
    visual scan.

## Implementation decisions

**Scope.** 7 commands + 7 skills. `resynthesize` excluded.

**`.claude/` adapter is untouched.** `.claude/skills` and `.claude/commands` are
*directory symlinks* into `.agents/` (each tracked by git as a single entry; the
adapter README confirms it). Renaming a dir/file inside `.agents/` passes through the
symlink automatically, so the renamed children appear under `.claude/` with no edit
there.

**Use `git mv`** for every dir/file rename, to preserve history.

**Replacement is scoped to identifier/invocation forms only.** A naive global
find-replace would corrupt prose. Rename only:
- `/command` invocations (`/spec` → `/bench-spec`),
- `.agents/{skills,commands}/<name>` paths,
- skill `name:` frontmatter,
- the `AGENTS.md` skills-index entries,
- "the `<skill>` skill" pointers and `split (<skill>)` status pointers in `bin/bench.sh`.

Do **not** rename the bare English words: the noun "spec" (the artifact), the verb
"build", the concept/leading-word "seam(s)". These carry meaning independent of the
command/skill id and must stay.

**The gate and CLI carry hardcoded name-literals that MOVE with the rename — this is
authorized, and it is not weakening a check.** `.bench/gate.sh` and `bin/bench.sh`
name `build.md` and `seams` as literal example fixtures in several places. Moving
those literals to the new names is fixture-following, not assertion-weakening: the
gate's *generic* name-sync checks (5a–5c loop over disk and `AGENTS.md`, so they
auto-adapt) and the canary suite (section 7) still bite, proving no check was gutted.
Enumerated literals to move:
- npm-pack required[] list (`.agents/commands/build.md`, `.agents/skills/seams/SKILL.md`).
- `bench link` contract literals — section 1e (`.agents/` and `.claude/` forms of the
  same two files).
- `bench status` contract — section 1F: `split (seams)` and the `/grill → /spec`
  action string. (`/resynthesize` in 1F stays — it is excluded from this slice.)
- the start-ideation seam check — section 5d (`.agents/commands/start-ideation.md` →
  `.agents/commands/bench-ideate.md`).

**Status string:** `split (seams)` → `split (bench-craft-seams)` — point users at the
*findable* id, not the concept word. Update `bin/bench.sh` and the gate's 1F check in
lockstep.

**Canary fixtures use synthetic names (`orphan`), not the real slate**, so they need
no rename. If any fixture's `EXPECT` references a renamed status string, update it in
lockstep — the canary run reports the mismatch.

## Testing decisions

**Seam: `bench gate` (the project oracle) — the single highest seam, and it already
exists.** No new test is written. The gate verifies the rename through checks it
already runs:
- **Name-sync, both directions** (sections 5a/5b/5c) — every skill dir ↔ `AGENTS.md`,
  every command file ↔ its `/name` reference. Generic loops, so they pass only when
  disk and index agree on the new names.
- **Link + pack integrity** (sections 1e, 4) — the moved literals must resolve to the
  renamed files.
- **Status contract** (section 1F) — the renamed status strings.
- **Canary** (section 7) — the known-broken fixtures must still go red with their
  target substrings, proving the literal-moves didn't turn any check into an
  always-pass.

**A good test here** exercises the kit's external conformance at the gate seam, not
internals — which is exactly what the gate does.

**Completeness sweep (gate-green is necessary, not sufficient).** The gate does *not*
fail on a stale `/spec` buried in README *prose* (5c only checks the index). So the
build's completeness check is a grep-sweep of the active surface for any surviving old
identifier form (`/start-ideation`, `/fix-bug`, `.agents/skills/seams`, …), asserting
zero hits — excluding `decisions/`, `ROADMAP.md`, `CHANGELOG.md`, git history, and the
excluded `resynthesize`. Run it during the build before calling the gate green.

**Gate command:** `bash .bench/gate.sh`.

## Out of scope

- **`resynthesize` → `bench-update` + `bench-learn` split.** A separate capability —
  two new commands reading different inputs (upstream kit vs local learnings journal),
  a behavior change, not a rename. Own spec. ~40 min.
- **Capture wiring** (agent runs `bench idea` on "park this"; no `bench-capture`
  surface). A separate capability — an interaction convention, not a rename. Own spec.
  ~20 min.
- **`bench-craft-skills` good/bad-example requirement** (Pocock-style contrastive
  examples). A change to the skill's *content*, already parked in `ROADMAP.md`. ~30 min.
- **A permanent "no dangling old command/skill name" gate check** (a rename-lint +
  canary fixture). A new conformance capability of its own; this slice uses a one-time
  grep-sweep instead. ~30 min. Optional.
- **HANDOFF.md refresh/retire.** Already a parked roadmap item; its staleness is about
  content, not names. Its `/name` invocation references *are* updated here (story 6).
