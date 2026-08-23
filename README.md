# Bench

A command-first workflow for agent work in small, gated phases. You invoke the
Bench phase for the work in front of you; the agent runs the CLI substrate
underneath.

---

## Reviewer quick start

In Claude Code, run Bench as slash commands:

```text
/bench-setup-repo
/bench-shape-idea
/bench-write-spec
/bench-debug
/bench-implement-spec
/bench-review-implementation
/bench-final-check
```

In Codex, invoke the matching adapter skills:

```text
$bench-setup-repo
$bench-shape-idea
$bench-write-spec
$bench-debug
$bench-implement-spec
$bench-review-implementation
$bench-final-check
```

Maintenance commands follow the same pattern: `/bench-update-kit` and
`/bench-drain` in Claude Code, or `$bench-update-kit` and
`$bench-drain` in Codex.

Other AGENTS.md harnesses read the matching file under `.agents/commands/` when
they do not expose a native command or skill surface.

For a new repo, ask the agent to run `/bench-setup-repo` or `$bench-setup-repo`.
That phase runs `bench setup` to converge the repo. It then walks you through
the project-specific gate, profile, lines, and an optional `CONTEXT.md`.

For feature work, use the command path:

```text
clear reviewed intent -----------------> /bench-write-spec -> /bench-implement-spec -> /bench-review-implementation -> /bench-final-check
multi-session decision fog -> /bench-shape-idea ---^
```

Decision maps are situational: `/bench-shape-idea` uses decision tickets only
when reviewer choices form a multi-session dependency tree. It compiles a
ready map beside its spec, under `specs/<slug>/decisions/`. A clear idea may
instead authorize `/bench-write-spec` through the reviewer-confirmed
current conversation or a named reviewed artifact. Spec authoring records
exactly one `Decision source:` line and owns engineering seams and coverage.
Implementation then derives independently-green implementation tickets. For
bugs, use `/bench-debug`; it builds the repro loop first.

Each command orients you at entry. It then hands you off at exit with what
changed, the current artifact or gate state, and the single next command it
recommends. The CLI commands below are the worker and maintainer substrate,
not the reviewer's first operating surface.

---

## Why Bench exists

Bench fuses **Matt Pocock's planning pipeline** with **kunchenguid's operational
substrate**. Four invariants hold the two together.

Pocock gives the brain: the method that turns a loose idea into a sequenced
plan, a spec that names its seams up front, and a build that tests the right
places. Kun Chen gives the body: isolated worktrees, an autonomous gated loop,
and the ergonomics that make tools cheap for agents to drive. The invariants
are the connective tissue — the rules that decide who has authority when the
two disagree.

The name is the coaching bench. You run line changes from it; you don't touch
the puck. You deploy the right **line** (model plus effort) for each **shift**
— a bounded unit of agent work in a clean worktree. A shift counts only when
it **backchecks clean**: it passes an external gate you never let the agent
grade.

---

## Operating guide

The shared working agreement is canonical in `.bench/BENCH.md`: roles,
invariant authority, and workflow proportionality. It also covers
communication rules and how the gate, hooks, skills, commands, and CLI fit
together. README is only the onboarding surface.

Lookup material lives on demand in `.bench/BENCH-reference.md`: the file map,
the generated skills index, and harness invocation forms. It also holds the
CLI command list, the shift adapter contract, and the hook layers.

---

## Layout

```
bench/
├── AGENTS.md                 # canonical working agreement — every harness reads this
├── CLAUDE.md                 # imports AGENTS.md + .bench/BENCH.md (for Claude Code)
├── .agents/
│   ├── commands/             # portable workflow phases + maintenance commands
│   │   ├── bench-setup-repo.md
│   │   ├── bench-update-kit.md
│   │   ├── bench-drain.md
│   │   ├── bench-shape-idea.md
│   │   ├── bench-write-spec.md
│   │   ├── bench-debug.md
│   │   ├── bench-implement-spec.md
│   │   ├── bench-review-implementation.md
│   │   └── bench-final-check.md
│   └── skills/               # portable generation-shaping skills; generated index in .bench/BENCH-reference.md
├── .bench/
│   ├── BENCH.md              # full Bench operating guide installed into projects
│   ├── BENCH-reference.md    # on-demand lookup: file map, skills index, adapter + hook detail
│   ├── gate.sh               # the project gate (oracle); kept out of the shipped npm package
│   ├── lines.env             # machine-readable tier->model binding read by hooks + adapters
│   ├── adapters/             # reference BENCH_AGENT adapters for bench shift
│   │   ├── claude
│   │   ├── codex
│   │   └── opencode
│   ├── hooks/                # shared hook scripts used by harness adapters
│   │   ├── session-start.sh
│   │   ├── stop.sh
│   │   ├── block-dangerous-git.sh
│   │   └── check-agent-line.sh   # denies Agent delegations off the bound line
│   ├── lib/                  # shared shell helpers the hooks and gate source
│   └── bin/                  # link-installed local CLI set (consumer repos; not in the kit)
├── .claude/
│   ├── README.md             # explains Claude adapter paths -> .agents and .bench/hooks
│   ├── settings.json         # Claude Code adapter pointing at .bench/hooks
│   ├── commands -> ../.agents/commands
│   ├── skills -> ../.agents/skills
│   └── hooks -> ../.bench/hooks
├── .codex/
│   └── hooks.json            # Codex adapter pointing at .bench/hooks
├── bin/
│   ├── bench.sh              # strangler router to the Go core + one-glance run_gate adapter
│   └── bench-postinstall.sh  # best-effort global shim installer
├── cmd/bench/                # the compiled core's main: dispatch + version
├── dist/                     # gitignored dev build; the gate rebuilds bench here (never committed)
├── internal/                 # the Go core: AXI query surface, status, structure + plumbing
│   ├── adopt/                # link/init/doctor adoption mutators
│   ├── toon/                 # the shared flat-table TOON emitter (one escaping rule)
│   ├── learnings/            # bench learnings parser
│   ├── maps/                 # bench maps engine (decision tickets + close-readiness + count)
│   ├── guards/               # bench guards aggregation
│   ├── diff/                 # bench diff review-base resolution
│   ├── coverage/             # bench coverage extraction + --check validation
│   ├── outline/              # bench outline candidate-seam indexer (file:line as TOON)
│   ├── status/               # bench status renderer + merged-spec retirement counter
│   ├── dashboard/            # bench dashboard: the dashboard page (self-contained static HTML)
│   ├── structure/            # structure checker + budgets parser (whole-tree + touched)
│   ├── worktree/             # worktree pool-path + lease-file conventions
│   ├── models/               # bench models advisory discovery inventory
│   ├── roadmap/              # capture/IDEAS.md + ROADMAP.md index + roadmap/ detail owner
│   └── git/                  # shared git subprocess helpers + gate tree-hash
└── projects/
    ├── benchkit.md           # seams, gate, lines for this kit
    └── gl-axi.md             # seams, gate, AXI conformance for gl-axi
```

---

## Worker and maintainer CLI

The reviewer-facing setup path is the setup command above. The worker-facing
mechanics underneath are the `bench` CLI commands here.

**Prerequisites.** Bench runs on macOS or Linux; Windows is unsupported, so
use WSL2. Until the packages are published, the install path runs straight
from the git repo. So you need **all three**: access to the `gibbonmi/bench`
repo, **Node** (to run `npx`), and **Go**. npx builds the compiled core on
your machine at install time. If you install under a node version manager,
mind the PATH-shim caveat noted with the durable-install steps below.

The fastest way for the worker to wire a repo today is one `npx` command
straight from git — nothing to clone, no global install. This is the
git-dependency form: it still requires a Go toolchain on your machine (npx
builds the compiled core at install time), and it's what actually runs until
Bench publishes to npm. Pin the ref so npx's git cache serves the build you
expect:

```sh
cd ~/src/your-project
npx github:gibbonmi/bench#main setup   # inspect, preview, and converge the repo
```

Once Bench has a first npm publish, the same one command becomes a pinned,
Go-toolchain-free install: `npx redbench@<version> setup`. That form doesn't
work yet — it's not published — so use the git-dependency form above until it
lands.

Run from `npx`, `setup` copies the kit in (the npx cache is ephemeral, so it
won't leave dangling symlinks). Prefer to install once and get a durable
`bench` command? Clone the repo, build the core, and symlink the launcher:

That build command is also the one you rerun after you change the Go
sources. Only `scripts/go-build.sh` refreshes the dev build and reseals it, so
a plain `go build` leaves an executable Bench will not trust. In a repository
that declares Go build inputs — the kit itself — `bench worktree land` proves
its own executable before it enforces any landing contract. It refuses with
the exact rebuild command rather than enforcing a contract that has since
been retired. A linked repository declares no build inputs and never pays
that proof.

```sh
git clone https://github.com/gibbonmi/bench ~/src/bench
bash ~/src/bench/scripts/go-build.sh ~/src/bench ~/src/bench/dist/bench
ln -s ~/src/bench/bin/bench.sh ~/.local/bin/bench
cd ~/src/your-project
bench link             # copies Bench-owned assets in safely
bench link symlink     # optional dogfood mode: point installed assets at this kit
bench init
```

A global install under a node version manager (nvm, asdf, fnm, volta) also
drops a plain-shell `bench` shim on a stable PATH dir. So login shells and
`bash -c` resolve `bench` the same as an interactive shell; `bench doctor`
reports its health and `bench doctor --fix` repairs it.

To uninstall, start with the per-repo footprint: `bench unlink` consumes the
link manifest and reverses the install. It removes the managed files whose
fingerprints still match (including a `CLAUDE.md` that `bench link` itself
created), and prunes emptied managed directories. It also strips the managed
AGENTS.md block while keeping your prose, and removes the bench-managed
pre-push hook.
A file you edited since linking is left in place, as is a `CLAUDE.md` that
predates link (link never records one, even a present-but-empty file). So are
your own artifacts: ROADMAP.md, `roadmap/`, capture/IDEAS.md, CONTEXT.md,
`specs/`, `decisions/`, `capture/learnings.md`, and `.bench/gate.sh`.
Rehearse it first with `bench unlink --dry-run`, which prints the exact plan
and changes nothing:

```sh
cd ~/src/your-project
bench unlink --dry-run   # rehearse: print the removal plan, touch nothing
bench unlink             # remove the per-repo Bench footprint
```

Then remove the global tool — the package and the shim:

```sh
npm uninstall -g redbench && rm -f "$(command -v bench)"
# after npm's own symlink is gone, command -v bench resolves to the shim;
# `bench doctor` prints the machine-exact removal pair while bench still resolves.
```

A repo linked before the manifest existed has nothing for `bench unlink` to
consume, so it exits 1. Remove that footprint by hand (the managed
`AGENTS.md` block, the `.bench/`, `.agents/`, `.claude/`, and `.codex/`
assets, and the pre-push hook).

The reviewer action is the setup phase, not those CLI calls:

```
/bench-setup-repo
# or, in Codex:
$bench-setup-repo
```

Setup is two halves, the same split as Pocock's `setup-matt-pocock-skills`.
The **mechanical** half is the CLI: `link` wires the kit into the repo for
every harness, and `init` scaffolds an empty `.bench/gate.sh`. Both are
deterministic and idempotent. `/bench-setup-repo` confirms whether those steps
already happened, and runs or reports the worker-facing step that is still
needed. It then continues into the **project-specific** half.

That half
explores the repo and walks the reviewer through the gate (the load-bearing
choice) and the profile (seams + lines + design-source path). It also covers
an optional `CONTEXT.md`, one decision at a time, and writes them. That
second half is the part that
can't be hardcoded — the gate command, the seams, and the lines differ in
every repo. So it's an interview, not a script.

`bench link` is idempotent and harness-neutral. It preserves project-owned
files, adds or updates only the managed Bench block in `AGENTS.md`, and
installs the full guide at `.bench/BENCH.md`. It also copies portable skills
and commands into `.agents/`, and installs Claude and Codex hook adapters
that call shared `.bench/hooks/` scripts. It installs a local hook CLI set
under `.bench/bin/`, and installs a git `pre-push` guard.

Copy mode is the
default;
use `bench link symlink` only when you intentionally want a dogfood repo to
follow live edits in a central kit checkout. If a project already owns a same-named
skill, command, or pre-push hook, `bench link` fails with a conflict report
instead of overwriting it. If you copy `bench` somewhere by hand, set
`BENCH_KIT=/path/to/kit`.

## Migrating from Matt Pocock's skills

If a repo is already set up with Pocock's engineering skills, you're most of
the way there. Bench builds on the same substrate, so adopt what's there
rather than restarting. `bench link` won't clobber an existing `CONTEXT.md`,
`docs/adr/`, or project-owned `AGENTS.md`. It appends or replaces only the
managed Bench block, and adds Bench's portable skills and commands alongside
non-conflicting project assets. It also writes the `CLAUDE.md` imports when
absent (retrofitting only the exact file an older link wrote); an edited
`CLAUDE.md` is project-owned and stays untouched.

The pieces line up
directly: a Pocock `CONTEXT.md` is read as-is (Bench's cold-session
convention is the same file). And `docs/adr/` is exactly where Bench's
`craft-adr` skill already writes, so your decision records carry over
untouched. Where Pocock's `setup-matt-pocock-skills` recorded an issue
tracker and domain layout under `docs/agents/`, `/bench-setup-repo` reads
those if present and won't re-ask.

What Bench adds on top is the part Pocock's skills leave to you. It adds the
**gate** as an external oracle (`.bench/gate.sh`) and the **gated shift
loop** that commits only on green. It also adds the **declared line**
(model + effort) per run, and the **profile** (`projects/<name>.md`) that
names the seams.

So migration is: ask the agent to run
`/bench-setup-repo`. It confirms the
link/init mechanics, detects the existing Pocock structure, and reuses it.
It only asks for the things Bench introduces (the gate command, the seams,
the lines). Nothing about Pocock's planning flow is replaced; Bench wraps it in
enforcement it didn't have.

## Keeping Bench current

Both upstream repos move. `/bench-update-kit` re-runs the synthesis against
their latest state. It pulls Pocock's skills and kunchenguid's tooling, and
diffs them against what Bench already incorporates (the provenance table is
the current record). It then proposes adoptions and runs three quality loops
before anything ships.

The first loop is anti-sediment (`craft-skills`: does
the change earn its place or just enlarge the kit?). The second is a
consistency audit (re-grep for stale references, invariant drift,
app-specific leakage). The third is the dogfood loop — a real shift on a
real repo with the changed kit. It is the only loop with the authority to
actually accept a change. It respects closed decisions: something Bench
already rejected isn't re-litigated unless the upstream version materially
changed. It proposes; you own the merge.

## Switching harnesses

You switch harnesses with no extra step, by design. After `bench link`, the same repo
is wired for all of them at once:

- **Claude Code** reads `CLAUDE.md` (which imports `AGENTS.md` and
  `.bench/BENCH.md`), auto-loads skills through `.claude/skills/`, runs
  `/bench-shape-idea`, `/bench-write-spec`, … as slash commands, and fires
  the Stop and PreToolUse hooks through `.claude/settings.json`.
- **Codex** reads `AGENTS.md`, finds skills in `.agents/skills/`, and uses
  the `$bench-*` phase adapters documented in `.bench/BENCH.md`.
  Model-invoked Bench guidance uses visible `craft-*` names so `$bench` stays
  phase-only.
- **OpenCode / other AGENTS.md harnesses** read `AGENTS.md` natively. They
  run commands by reading the file in `.agents/commands/` when they do not
  expose a native command or skill invocation surface.

There is one portable content surface: `.agents/{skills,commands}`. Claude
Code uses adapter paths under `.claude/`; Codex uses `.codex/hooks.json`;
both hook adapters call the same shared scripts in `.bench/hooks/`. The
enforcement that matters remains harness-independent too: the `bench shift`
loop runs the gate after every iteration and commits only on green. The
git `pre-push` hook protects the branch it resolves no matter which agent
(or human) pushes.

Env knobs: `BENCH_AGENT` (required for `bench shift` — a harness adapter
executable that takes the prompt as `$1`; reference adapters ship in
`.bench/adapters/`), `BENCH_MAX_ITERS`, `BENCH_GATE` (a gate command if you'd
rather not ship `.bench/gate.sh`).

---

## Workflow

The workflow contract is canonical in `.bench/BENCH.md`. Use the quick-start
path above as the reviewer surface. Then read the guide for when to shape,
spec, implement, review, final-check, debug, or run an autonomous shift.

---

## The design system as a visual oracle

A UI project's design system lives in its own repo and plugs in as the
**visual oracle**. That's the third gate axis, after tests (behavior) and
the screenshot loop (interaction). The project consumes it as a submodule,
package, or pinned path. The `craft-design-system` skill makes the agent
consume it rather than reinvent it: every value references a token, and
every component composes from the inventory. The design-conformance check
fails the build on raw hex, hardcoded spacing, or a duplicated component.

For a UI shift the gate's green suite is necessary but not sufficient. Route
UI shifts to a mid line; the screenshot loop, not raw model strength, is what
catches the failures there.

The handoff is repo-to-repo, which is what makes it harness-agnostic. When a
shift needs a token or variant that doesn't exist, you add it in the design
repo. Do that via Claude Design when you're in a Claude session, or by
editing the repo directly under Codex or any other harness. Then commit,
re-pin, and build against it. Nothing in the UI workflow depends on which design tool or which
agent you're using; it depends only on the committed artifacts.

## Integrating it into gl-axi

`projects/gl-axi.md` is the profile, and the `craft-cli` skill is the design
spec. The key move: **AXI conformance is a gate check**, so an external oracle holds what you build to its own standard.

```sh
cd ~/src/gl-axi
# the gate runs three oracles in order of authority:
#   pytest                 — behavior at the output boundary + glab adapter seams
#   axi-conformance        — TOON stdout, minimal schemas, structured errors, exit codes
#   bench-glab-delta       — your paired per-task harness vs raw glab, deterministic asserts
#
# add a new command wrapper — cheap line, it's mechanical once the boundary exists
bench shift "add 'mr list' wrapper emitting TOON per the craft-cli skill, with a conformance test"
```

The conformance check and your paired-delta harness are deterministic
assertions (a TOON-shape check is a parser, not a model). Because of that,
the loop can't pass by fooling a judge. A shift that makes gl-axi worse than `glab` on any task
fails the gate and never commits. That's invariant #1 pointed directly at
the tool you're building.

---

## Where each piece came from

| Bench piece | Pocock | Kun Chen | Your discovery |
| --- | --- | --- | --- |
| `/bench-shape-idea`, `/bench-write-spec`, `/bench-implement-spec`, `/bench-review-implementation`, `/bench-final-check` | decision-mapping, research, to-prd, implement, review | — | — |
| Bench craft skills | codebase-design, tdd, grilling, writing-great-skills | AXI spec | generated skills index in `.bench/BENCH-reference.md`; stateless-reader docs; effort, review, delegate, gate, and design guidance |
| `bench worktree` | — | treehouse | — |
| `bench shift` (gated loop) | — | gnhf + no-mistakes | gate-on-green, not self-graded |
| `/bench-setup-repo` (configure a repo) | setup-matt-pocock-skills | — | gate + profile + lines, interviewed |
| `/bench-update-kit` (sync upstream) | — | — | re-run the synthesis vs upstream, 3 loops |
| `/bench-drain` (reconcile roadmap, drain capture) | — | — | the kit learns from its own use, one reviewed batch diff |
| `/bench-debug` (bug path) | diagnosing-bugs | — | repro loop as the bug's gate |
| `/bench-deepen` (deepening survey) | improve-codebase-architecture | — | scopes from `ASSESSMENT.md` findings; vocabulary charged from `craft-seams`, grilling from `craft-grill`; mid-tier read-only delegate |
| design-it-twice in `craft-seams` | codebase-design | — | high-effort line at the uncertain seam |
| `bench shift` notes.md | — | gnhf (iteration context) | — |
| `block-dangerous-git.sh` | git-guardrails | — | agent has no destructive authority |
| Stop hook + `.bench/gate.sh` | — | no-mistakes (external gate) | the gate is the oracle |
| The line declaration | — | — | "suggest model and effort" |

The combination is more than the parts. Pocock's pipeline gives the shift a
real target (seams and stories chosen up front). Kun Chen's substrate
gives the target real teeth (an isolated, gated, autonomous loop). Your
invariants decide that when the agent's judgment and the gate disagree, the
gate wins. That single rule makes an autonomous loop safe to leave running.
</content>
