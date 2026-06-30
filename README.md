# Bench

A local agent-development workflow that fuses **Matt Pocock's planning pipeline**
with **kunchenguid's operational substrate**, held together by four invariants.

Pocock gives the brain: how to turn a loose idea into a sequenced plan, a spec
with the seams chosen up front, and a build that tests in the right places. Kun
Chen gives the body: isolated worktrees, an autonomous gated loop, and the
ergonomics that make tools cheap for agents to drive. The invariants are the
connective tissue — the rules that decide who has authority when the two
disagree.

The name is the coaching bench. You run line changes from it; you don't touch the
puck. You deploy the right **line** (model + effort) for each **shift** (a bounded
unit of agent work in a clean worktree), and a shift only counts if it
**backchecks clean** — passes an external gate you didn't let the agent grade.

---

## The four invariants

Everything here exists to enforce these. They're canonical in `.bench/BENCH.md` (the
guide `bench link` ships, so they reach every project) and the load-bearing ones are
also wired into hooks so they're not just vibes.

1. **The gate is the oracle — the agent never grades its own work.** "Done" means
   the gate exits zero (tests, types, lint, project conformance), not that the
   diff looks right. A Stop hook physically refuses to let a shift end on a red
   gate. This is the lesson from every benchmark in the field: the agent's own
   green tests are not a valid completion signal; an external, deterministic check
   is.

2. **Declare the line before a long run.** Before any multi-cycle stage, the agent
   states model + effort + token cap with one clause of justification. Cheap model
   for plumbing at a known seam; top model only where the answer is genuinely
   uncertain. No silent escalation.

3. **Document for the teammate who just walked in.** ADRs and docs record the
   current decided state, not the history of how it changed. Every cold agent
   session is that teammate; history lives in git.

4. **One small change at a time, repo stays green.** Smallest diff that advances
   the objective; commit on green, never on red; compose an existing seam before
   inventing one. This is what the shift loop mechanizes — and the agent has no
   authority over the merge or history: a `pre-push` git hook and a PreToolUse
   guard make that enforceable, not aspirational.

---

## The three layers

The kit keeps Pocock's "rules vs commands vs enforcement" separation explicit,
because conflating them is how agent workflows rot:

| Layer | What it is | Authority | Lives in |
| --- | --- | --- | --- |
| **Enforcement** | The gate + hooks | Hard. The agent can't override it. | `.bench/gate.sh`, `.bench/hooks/`, `.claude/`, `.codex/`, `bench` |
| **Generation-shaping** | Skills | Probabilistic. Nudges *how* the agent writes. | `.agents/skills/` |
| **Workflow discipline** | Commands | Canonical phases you invoke by name. | `.agents/commands/` |

If a skill and the gate disagree, the gate wins. That ordering is the whole point.

---

## Layout

```
bench/
├── AGENTS.md                 # canonical working agreement — every harness reads this
├── CLAUDE.md                 # one-line import of AGENTS.md (for Claude Code)
├── .agents/
│   ├── commands/             # portable workflow phases + maintenance commands
│   │   ├── setup.md
│   │   ├── resynthesize.md
│   │   ├── start-ideation.md
│   │   ├── spec.md
│   │   ├── fix-bug.md
│   │   ├── build.md
│   │   ├── prep-shift.md
│   │   └── verify-gate.md
│   └── skills/               # portable generation-shaping skills
│       ├── seams/
│       ├── tdd-at-seams/
│       ├── adr/
│       ├── axi/
│       ├── design-system/
│       ├── writing-great-skills/
│       └── grill/
├── .bench/
│   ├── BENCH.md              # full Bench operating guide installed into projects
│   └── hooks/                # shared hook scripts used by harness adapters
│       ├── stop.sh
│       └── block-dangerous-git.sh
├── .claude/
│   ├── README.md             # explains Claude adapter paths -> .agents and .bench/hooks
│   ├── settings.json         # Claude Code adapter pointing at .bench/hooks
│   ├── commands -> ../.agents/commands
│   ├── skills -> ../.agents/skills
│   └── hooks -> ../.bench/hooks
├── .codex/
│   └── hooks.json            # Codex adapter pointing at .bench/hooks
├── bin/
│   └── bench                 # worktrees (treehouse-lite) + gated loop (gnhf-lite)
└── projects/
    ├── regroup.md            # seams, gate, lines for Regroup
    └── gl-axi.md             # seams, gate, AXI conformance for gl-axi
```

---

## Install

The kit is an npm package, so the fastest way to wire a repo is one command with
`npx` — nothing to clone, no global install:

```sh
cd ~/src/your-project
npx benchkit link      # wire the kit into this repo for every harness
npx benchkit init      # scaffold .bench/gate.sh
```

Run from `npx`, `link` copies the kit in (the npx cache is ephemeral, so it won't
leave dangling symlinks). Prefer to install once and get a durable `bench` command?

```sh
npm i -g benchkit
# or, from a clone:  ln -s ~/src/bench/bin/bench.sh ~/.local/bin/bench
cd ~/src/your-project
bench link             # copies Bench-owned assets in safely
bench link symlink     # optional dogfood mode: point installed assets at this kit
bench init
```

Either way, finish by configuring the repo in your agent (one-time, interactive):

```
/bench-setup
```

Setup is two halves, the same split as Pocock's `setup-matt-pocock-skills`. The
**mechanical** half is the CLI: `link` wires the kit into the repo for every
harness, and `init` scaffolds an empty `.bench/gate.sh` — both deterministic, both
idempotent. The **project-specific** half is `/bench-setup`, a prompt-driven command your
agent runs once: it explores the repo, then walks you through the gate (the
load-bearing choice), the profile (seams + lines + design-source path), and an
optional `CONTEXT.md`, one decision at a time, and writes them. That second half is
the part that can't be hardcoded — the gate command, the seams, and the lines differ
in every repo — so it's an interview, not a script.

`bench link` is idempotent and harness-neutral. It preserves project-owned files,
adds or updates only the managed Bench block in `AGENTS.md`, installs the full guide
at `.bench/BENCH.md`, copies portable skills and commands into `.agents/`, installs
Claude and Codex hook adapters that call shared `.bench/hooks/` scripts, and installs
a git `pre-push` guard. Copy mode is the default; use `bench link symlink` only when
you intentionally want a dogfood repo to follow live edits in a central kit checkout.
If a project already owns a same-named skill, command, or pre-push hook, `bench link`
fails with a conflict report instead of overwriting it. If you copy `bench` somewhere
by hand, set `BENCH_KIT=/path/to/kit`.

## Migrating from Matt Pocock's skills

If a repo is already set up with Pocock's engineering skills, you're most of the way
there — Bench builds on the same substrate, so adopt what's there rather than
restarting. `bench link` won't clobber an existing `CONTEXT.md`, `docs/adr/`, or
project-owned `AGENTS.md`; it appends or replaces only the managed Bench block,
adds Bench's portable skills and commands alongside non-conflicting project assets,
and writes the `CLAUDE.md` import only if one doesn't exist. The pieces line up directly: a Pocock
`CONTEXT.md` is read as-is (Bench's cold-session convention is the same file), and
`docs/adr/` is exactly where Bench's `bench-craft-adr` skill already writes, so your decision
records carry over untouched. Where Pocock's `setup-matt-pocock-skills` recorded an
issue tracker and domain layout under `docs/agents/`, `/bench-setup` reads those if present
and won't re-ask.

What Bench adds on top is the part Pocock's skills leave to you: the **gate** as an
external oracle (`.bench/gate.sh`), the **gated shift loop** that commits only on green,
the **declared line** (model + effort) per run, and the **profile** (`projects/<name>.md`)
that names the seams. So migration is: run `link` and `init`, then run `/bench-setup` — it
detects the existing Pocock structure, reuses it, and only asks for the things Bench
introduces (the gate command, the seams, the lines). Nothing about Pocock's planning
flow is replaced; Bench wraps it in enforcement it didn't have.

## Keeping Bench current

Both upstream repos move. `/resynthesize` re-runs the synthesis against their latest
state: it pulls Pocock's skills and kunchenguid's tooling, diffs them against what
Bench already incorporates (the provenance table and `CHANGELOG.md` are the record),
and proposes adoptions — then runs three quality loops before anything ships. The
first loop is anti-sediment (`bench-craft-skills`: does the change earn its place
or just enlarge the kit?), the second is a consistency audit (re-grep for stale
references, invariant drift, app-specific leakage), and the third is the dogfood loop
— a real shift on a real repo with the changed kit, which is the only loop with the
authority to actually accept a change. It respects closed decisions: something Bench
already rejected isn't re-litigated unless the upstream version materially changed.
It proposes; you own the merge.

## Switching harnesses

Switching harnesses is a no-op, by design. After `bench link`, the same repo is
wired for all of them at once:

- **Claude Code** reads `CLAUDE.md` (which imports `AGENTS.md`), auto-loads skills
  through `.claude/skills/`, runs `/bench-ideate`, `/bench-spec`, … as slash commands,
  and fires the Stop and PreToolUse hooks through `.claude/settings.json`.
- **Codex / OpenCode / other AGENTS.md harnesses** read `AGENTS.md` natively, find
  skills in `.agents/skills/`, and run the commands by reading the file in
  `.agents/commands/` (they're phases, not slash commands, on these harnesses).

There is one portable content surface: `.agents/{skills,commands}`. Claude Code uses
adapter paths under `.claude/`; Codex uses `.codex/hooks.json`; both hook adapters call
the same shared scripts in `.bench/hooks/`. The enforcement that matters remains
harness-independent too: the `bench shift` loop runs the gate after every iteration
and commits only on green, and the git `pre-push` hook protects the default branch no
matter which agent (or human) pushes.

Optional knobs (env): `BENCH_AGENT` (headless agent command, default `claude`),
`BENCH_MAX_ITERS`, `BENCH_MAX_TOKENS`, `BENCH_GATE` (a gate command if you'd rather
not ship `.bench/gate.sh`).

---

## The workflow

Two motions. The **planning motion** is conversational and human-paced; the
**shift motion** is the autonomous loop.

### Planning motion (you + the agent, interactive)

```
loose idea ──/bench-ideate──▶ decision map ──/bench-spec──▶ spec (seams chosen) ──▶ ready to build
   (skip /bench-ideate if there's no real fog — go straight to /bench-spec)
```

- `/bench-ideate` only when the idea needs more than one session of decisions. It uses
  `/bench-craft-grill` to surface them and writes `decisions/<topic>.md`. If there's no fog,
  it tells you to skip ahead.
- `/bench-spec` synthesizes the conversation into `specs/<feature>.md`: user stories
  (the breadth), the **seams chosen before any code exists** (where tests attach),
  and the gate that defines done. This step is what keeps the later loop honest —
  the target is set by you, not invented mid-loop.

There are two ways into the shift motion: the **feature path** above
(`/bench-ideate` → `/bench-spec`) and the **bug path** (`/bench-diagnose`). A bug doesn't get a spec —
it already has one, the thing should work and doesn't. `/bench-diagnose` builds a tight,
red-capable repro loop *first*; that loop becomes the gate the fix shift runs
against, so the fix is done when the loop goes green, not when the agent says so.

### Shift motion (the agent, gated and autonomous)

```
bench shift "<objective>"
   │
   ├─ clean worktree, fresh branch
   ├─ iterate: one small change ▶ run gate   (notes.md carried between iterations)
   │     gate green ▶ commit the iteration
   │     gate red   ▶ roll back, retry  (Stop hook blocks "done" on red)
   └─ stop at: objective met │ iteration cap │ you pull the line (Ctrl-C)
        ▶ you review the branch and own the merge
```

`/bench-build` and `/bench-qa` are the manual equivalents when you want to drive a single
build by hand instead of running the loop. Same gate, same rules.

---

## Integrating it into Regroup

`projects/regroup.md` is the profile. The seams are the phase/game-state machine,
`CoordinateProvider`, and the event store; UI is gated by the screenshot loop, not
unit tests.

A typical feature — say, adding zone-entry events to the phase taxonomy:

```sh
cd ~/src/regroup
# 1. is the domain change settled? if not, map it.
#    /bench-ideate  →  grills the zone-entry decision, writes decisions/zone-entries.md
# 2. spec it — this picks the seam (the state machine) and the tests up front
#    /bench-spec →  specs/zone-entries.md with stories + the transition-test seam
# 3. run the shift. heavy line, because the ontology is the uncertain seam.
BENCH_MAX_ITERS=8 bench shift "add zone-entry events to the phase taxonomy per specs/zone-entries.md"
# each iteration commits only if mypy + pytest + ruff pass. you review the branch.
```

For a UI shift (e.g., a zone-entry marker on the timeline) the gate's green suite
is necessary but not sufficient: the `regroup-ui` skill triggers automatically,
the screenshot loop runs, and the two chasms plus the five interaction states are
the real review. The shuttle slider stays canonical — composed, never regenerated.
Route UI shifts to a mid line; the screenshot loop, not raw model strength, is
what catches the failures there.

### The design system (separate repo)

Your Regroup design system lives in its own repo and plugs in as the **visual
oracle** — the third gate axis after tests (behavior) and the screenshot loop
(interaction). Regroup consumes it as a submodule/package/pinned path; the
`bench-craft-design-system` skill makes the agent consume it rather than reinvent it: every
value references a token, every component composes from the inventory, and the
design-conformance check fails the build on raw hex, hardcoded spacing, or a
duplicated component.

The handoff is repo-to-repo, which is what makes it harness-agnostic. When a shift
needs a token or variant that doesn't exist, you add it in the design repo — via
Claude Design when you're in a Claude session, or by editing the repo directly under
Codex or any other harness — commit, re-pin, then build against it. Nothing in the
UI workflow depends on which design tool or which agent you're using; it depends
only on the committed artifacts.

## Integrating it into gl-axi

`projects/gl-axi.md` is the profile, and the `bench-craft-cli` skill is the design spec. The
key move: **AXI conformance is a gate check**, so the thing you're building is held
to its own standard by an external oracle.

```sh
cd ~/src/gl-axi
# the gate runs three oracles in order of authority:
#   pytest                 — behavior at the output boundary + glab adapter seams
#   axi-conformance        — TOON stdout, minimal schemas, structured errors, exit codes
#   bench-glab-delta       — your paired per-task harness vs raw glab, deterministic asserts
#
# add a new command wrapper — cheap line, it's mechanical once the boundary exists
bench shift "add 'mr list' wrapper emitting TOON per the bench-craft-cli skill, with a conformance test"
```

Because the conformance check and your paired-delta harness are deterministic
assertions (a TOON-shape check is a parser, not a model), the loop can't pass by
fooling a judge. A shift that makes gl-axi worse than `glab` on any task fails the
gate and never commits. That's invariant #1 pointed directly at the tool you're
building.

---

## Where each piece came from

| Bench piece | Pocock | Kun Chen | Your discovery |
| --- | --- | --- | --- |
| `/bench-ideate`, `/bench-spec`, `/bench-build`, `/bench-review`, `/bench-qa` | decision-mapping, to-prd, implement, review | — | — |
| `bench-craft-seams`, `bench-craft-tdd`, `bench-craft-grill`, `bench-craft-adr`, `bench-craft-skills` | codebase-design, tdd, grilling, writing-great-skills | — | stateless-reader docs; effort declaration |
| `bench-craft-cli` skill | — | AXI spec | TOON-first pipeline; conformance-as-gate |
| `bench-craft-design-system` skill + design gate | regroup-ui (canonical components) | — | design system as visual oracle (separate design repo) |
| `bench worktree` | — | treehouse | — |
| `bench shift` (gated loop) | — | gnhf + no-mistakes | gate-on-green, not self-graded |
| `/bench-setup` (configure a repo) | setup-matt-pocock-skills | — | gate + profile + lines, interviewed |
| `/resynthesize` (stay current) | — | — | re-run the synthesis vs upstream, 3 loops |
| `/bench-diagnose` (bug path) | diagnosing-bugs | — | repro loop as the bug's gate |
| design-it-twice in `bench-craft-seams` | codebase-design | — | high-effort line at the uncertain seam |
| `bench shift` notes.md | — | gnhf (iteration context) | — |
| `block-dangerous-git.sh` | git-guardrails | — | agent has no destructive authority |
| Stop hook + `.bench/gate.sh` | — | no-mistakes (external gate) | the gate is the oracle |
| The line declaration | — | — | "suggest model and effort" |

The combination is more than the parts: Pocock's pipeline gives the shift a real
target (seams and stories chosen up front), and Kun Chen's substrate gives the
target real teeth (an isolated, gated, autonomous loop). Your invariants decide
that when the agent's judgment and the gate disagree, the gate wins — which is the
single rule that makes an autonomous loop safe to leave running.
