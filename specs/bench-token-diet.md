# BENCH.md token diet

## Problem

`.bench/BENCH.md` is loaded in full on **every** session of **every** linked repo
(via `CLAUDE.md`'s `@.bench/BENCH.md` import, and natively by AGENTS.md harnesses).
It is ~2,400 words, but roughly 1,028 of them — the **Files** map, the **Skills
index**, the **Harness Invocation** table, the **Commands** (CLI) list, the
**Harness adapter for the shift loop** contract, and **Hook Layers** — are
reference material an agent looks up occasionally, not generation-steering
guidance it needs resident every turn. That's ~1,370 tokens paid unconditionally,
forever, on every repo, to keep a lookup table warm.

The naive fix — cut or shrink the prose — is the wrong lever (the reference
content is already terse and the wording carries meaning). The right lever is
**placement**: move the lookup sections to an on-demand file the agent reads only
when it needs a command name or adapter wiring, and keep the always-loaded guide
to the content that steers every message and every decision.

## Solution

Split `.bench/BENCH.md` into two files:

- **`.bench/BENCH.md`** stays the always-loaded operating guide: Roles, How the
  pieces fit, The four invariants, How to talk to me, Workflow, Capture, and the
  one generation-steering rule currently buried in Harness Invocation ("translate
  every recommendation into the surface *this* harness exposes"). It gains one
  pointer line to the reference file.
- **`.bench/BENCH-reference.md`** (new, on-demand) holds the lookup sections:
  Files, Skills index, Harness Invocation table + Codex adapter list, Commands,
  Harness adapter for the shift loop, Hook Layers. It is *referenced by path*, not
  `@`-imported, so it costs zero tokens until an agent opens it.

The move re-points three gate doc-currency checks (skills index, command `/name`
refs, Codex `$bench-*` adapters) at the reference file so the oracle keeps biting,
and adds one new check that the reference file stays *referenced-not-imported* so
the diet can't silently regress. The reference file is added to `package.json`
`files[]` so `bench link` delivers it to consumers.

## User stories

1. As an agent, I want the Files map, Skills index, Harness Invocation table,
   Commands list, shift-adapter contract, and Hook Layers moved out of the
   always-loaded guide into `.bench/BENCH-reference.md`, so I don't pay ~1,370
   tokens/session to keep lookup material resident that I rarely read.
   Line: claude-fable-5 / high. This edits the shared platform rules, where the
   leverage override in craft-line routes top-and-high because a relocation
   defect is invisible to the gate's prose coverage and multiplies through every
   session that loads the guide.

2. As an agent, I want the always-loaded guide to carry one pointer line naming
   the reference file and what it holds, so I know where to read command names and
   adapter wiring the moment I need them.
   Line: claude-fable-5 / high. The pointer is the one thread connecting the split
   halves, and shared-platform-rule wording routes top-and-high under the leverage
   override.

3. As an agent, I want the single generation-steering rule inside Harness
   Invocation ("recommend the phase in the form this harness can invoke") to stay
   in the always-loaded guide rather than move with the reference table, so
   relocation does not degrade how I translate next-actions to the active harness.
   Line: claude-fable-5 / high. Deciding which sub-paragraphs steer generation
   versus merely inform is a judgment call on the shared rules, exactly the
   leverage-override case where cheap routing is the expensive choice.

4. As the reviewer, I want a gate check that fails if `.bench/BENCH-reference.md`
   is `@`-imported by `CLAUDE.md` (or otherwise pulled into the always-loaded
   set), so the token diet cannot be silently undone by a later edit or a relink.
   Line: claude-fable-5 / high. This is the oracle check that makes the whole
   feature durable rather than a one-time edit, so it is authored at the top line
   per craft-gate and the leverage override.

5. As the kit maintainer, I want the three doc-currency checks that assert content
   in `.bench/BENCH.md` — the generated skills index, the command `/name` refs,
   and the Codex `$bench-*` adapter refs — to resolve against wherever that content
   now lives and still fail on drift, so moving the sections does not blind the
   oracle.
   Line: claude-fable-5 / high. Re-pointing oracle checks is gate-authoring work
   where a silent fail-open would let the moved docs rot undetected.

6. As the kit maintainer, I want no shared content duplicated between the
   always-loaded guide and the reference file, so the split produces one source
   per fact rather than two derivations that drift.
   Line: claude-fable-5 / high. The one-source-per-fact standard is a shared
   platform rule and a duplicated section is precisely the drift class the kit's
   code standard forbids.

7. As a consumer repo, I want `.bench/BENCH-reference.md` listed in `package.json`
   `files[]`, so `bench link` ships it and the pointer in the always-loaded guide
   never dangles on a linked repo.
   Line: claude-sonnet-4-6 / low. This is a one-line mechanical addition to the
   package manifest, fully observed by the existing package-surface gate check.

## Implementation decisions

- **New file `.bench/BENCH-reference.md`.** Holds the six relocated sections
  verbatim (headings preserved). No `@`-import anywhere; discovered by the pointer
  line in `.bench/BENCH.md`. Ships via `package.json` `files[]` (explicit path,
  matching how `.bench/BENCH.md` is listed — `.bench/` is not globbed).

- **What stays hot vs moves — the split.** *Stay* in `.bench/BENCH.md`: Roles,
  How the pieces fit, The four invariants, How to talk to me, Workflow, Capture,
  plus the "match every recommendation to the active harness" steering paragraph
  lifted out of Harness Invocation. *Move* to the reference file: Files, Skills
  index, Harness Invocation (the mechanical invocation table + Codex adapter
  list), Commands, Harness adapter for the shift loop, Hook Layers. The Harness
  Invocation section is split, not moved whole — its steering rule stays, its
  lookup table goes. This split is the reviewer-judged decision surfaced at spec
  approval; the gate cannot grade it.

- **New architectural pattern: referenced-not-imported doc.** The kit currently
  has no on-demand doc; every guide is `@`-imported. This introduces the pattern
  (a file an agent reads by path on need) and a gate check that guards it. Future
  reference material can follow the same pattern.

- **Gate changes (`.bench/gate.sh` + `.bench/skills-index.sh`):**
  - `skills-index.sh` writes/checks the generated block in the reference file
    instead of `.bench/BENCH.md` (one path constant; `--write` and `--check` move
    together).
  - Gate check 5c (command `/name` referenced) and 5d (Codex `$bench-*` adapter
    documented) search the reference file for the moved refs. If a `/name` also
    appears in a hot section (Workflow), the check passes on either — the contract
    is "referenced in the kit's shipped guide set," now two files.
  - New check: `.bench/BENCH-reference.md` exists, is in `files[]`, and is **not**
    `@`-imported by `CLAUDE.md`. Fail-closed (missing file or stray import → red).
  - New check: no `## ` section heading appears in **both** `.bench/BENCH.md` and
    `.bench/BENCH-reference.md` (the single-source guard for story 6).

- **AGENTS.md / CLAUDE.md untouched except the pointer.** `CLAUDE.md` keeps
  `@AGENTS.md` and `@.bench/BENCH.md`; it does **not** import the reference file.
  The pointer to the reference file lives in `.bench/BENCH.md` (so both Claude and
  native AGENTS.md harnesses see it).

## Testing decisions

- **A good test here** feeds `bench gate` a fixture tree — a correctly-split tree
  (green), and each broken variant (red) — and asserts the **exit code** and the
  **attributed error line**, never a reading of the diff. This matches the
  benchkit gate's existing doc-currency fragments.
- **Seam:** the gate (`bash .bench/gate.sh` / `bench gate`), the highest seam and
  the one that already exercises doc-currency. Prior art: gate checks 5a–5e and
  the shared `contract` fixture harness the contract fragments call.
- **Gate command:** the project gate, `bench gate`.

### Seam diagram

    trigger: bench gate  (Stop hook · shift loop · pre-push · manual)
        │
        ▼
    .bench/BENCH.md          ──▶  [ gate.sh:                       ]  ──▶  exit 0
    .bench/BENCH-reference.md ─▶  [   skills-index --check          ]        (green)
    CLAUDE.md                ──▶  [   5c command /name refs         ]  ──▶  exit 1
    package.json files[]     ──▶  [   5d codex $bench-* adapters     ]        + err
    skills/*, commands/*     ──▶  [   NEW: ref not @-imported        ]      line naming
                                  [   NEW: no section in both files  ]      the break
                                        ◀ tests attach here: build a fixture tree
                                          (split / half-moved / re-imported / dropped
                                          ref), run the gate, assert exit + err text

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | The six sections live in the reference file, not the always-loaded guide | gate + file presence | fixture with a moved section still in `.bench/BENCH.md` → the new "no section in both files" check turns `bench gate` red | a section left behind (or copied) is the exact half-done state; the both-files check fails on it |
| 2 | The always-loaded guide names the reference file | gate | fixture dropping the pointer line while the reference file exists → a currency check for the pointer turns gate red | without the pointer the agent can't find the moved lookup content; the check fails when the thread is cut |
| 3 | The "recommend in the active harness's form" rule stays hot | reviewer at spec approval | already covered / **not TDD-able** — which paragraphs steer generation is semantic, graded by the reviewer at approval, not by the gate | the gate can't judge steering-vs-reference; this row is honest that it's a sign-off gate, not a test |
| 4 | Reference file is referenced, not `@`-imported | gate | fixture adding `@.bench/BENCH-reference.md` to `CLAUDE.md` → new placement check turns gate red | an `@`-import silently re-loads the file every session, undoing the diet; the check bites on the import line |
| 5 | Skills index, `/name` refs, `$bench-*` adapters still enforced at the new location | gate | fixture deleting a moved command `/name` ref (and drifting the skills-index block) in the reference file → 5c and skills-index `--check` turn gate red | if the checks still pointed at the old file they'd fail-open on the moved content; red proves they follow the content |
| 6 | No section duplicated across the two files | gate | fixture with one `## ` heading present in both files → single-source check turns gate red | two copies of a section is the drift class the code standard forbids; the check fails on any shared heading |
| 7 | Reference file ships to consumers | gate | fixture removing `.bench/BENCH-reference.md` from `package.json` `files[]` → existing package-surface check (4) plus the new existence check turn gate red | an unshipped reference file makes every linked repo's pointer dangle; the manifest check fails when it's absent |

### Edge inventory

Walked per behavior; each resolved as a row above or a **Won't handle** line here.

- **empty/absent** — reference file missing entirely → covered by rows 4/7 (existence + `files[]` checks, fail-closed).
- **malformed** — moved `/name` ref or adapter ref deleted from the reference file → covered by row 5.
- **interrupted/partial state** — section copied into both files (half-move) → covered by rows 1/6 (both-files single-source check).
- **re-run idempotency** — `skills-index.sh --write` run twice must leave the block byte-identical → covered by row 5 (the `--check` fails if `--write` isn't idempotent at the new path).
- **hostile environment** — a later edit or `bench link` re-adds the `@`-import → covered by row 4.
- **boundary (the split point)** — a section that is part steering, part lookup (Harness Invocation) → row 3; reviewer-judged at approval, not gate-gradable.
- **Won't handle:** enforcing lazy-load on arbitrary third-party AGENTS.md harnesses — the pointer is a documented convention, and only Claude Code's `@`-import path is gate-enforced (it's the only always-load mechanism the kit controls); a non-Claude harness that eagerly slurps every `.bench/*.md` is out of the kit's reach.
- **Won't handle:** correctness of the reference *content itself* (that the CLI list is accurate) — unchanged by this move and already governed by the existing currency checks; relocating text doesn't reopen it.

## Out of scope

- **Skill description-frontmatter trim** (~932 words of skill `description:` YAML
  loaded into the Claude menu every session) — a *separate capability*: it touches
  a different surface (every skill's frontmatter, not BENCH.md) and carries a
  distinct risk (descriptions are probabilistic *triggers*, so trimming trades
  token savings against firing reliability under craft-skills, which the gate
  cannot grade). Estimate to build later: **~19 edits, 2 gate runs** (one
  description edit per craft + phase skill; the gate only checks frontmatter
  presence, so trigger-reliability needs a craft-skills judgment pass on top).
  Ranked highest of the cuts — real always-loaded savings, but needs its own spec
  because of the trigger-reliability judgment. Parked on the roadmap.

- **Wording/prose tightening of on-demand command and skill bodies** — assessed
  and declined as a *token* play: the bodies are execution prompts whose verbosity
  encodes judgment discipline, they're paid only per-invocation, and the polish is
  a quality task under craft-skills judged on clarity, not a token task. Not a
  separate capability worth a spec; recorded here so it doesn't resurface as
  token work.

- **Second-wave parsers** (`bench refs` / `doctor` / `detect`, ASSESSMENT Part 1
  finding B) — a *separate capability* already parked on the roadmap pending
  learnings-funnel evidence; speccing them now would override that reviewer
  decision.
