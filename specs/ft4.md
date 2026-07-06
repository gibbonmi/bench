# ft4 — harness task list in /bench-implement-spec

Status: implemented

## Problem

When `/bench-implement-spec` runs a multi-row build, the reviewer has no
at-a-glance view of progress against the spec's fixed target. The acceptance
coverage map defines the breadth, but during the build it lives only in the spec
file and in `bench coverage` output — the harness's own native task list (Claude
todos, Codex plan) sits empty or holds whatever ad-hoc items the agent invents.
The reviewer can't watch which coverage rows are done without asking.

## Solution

The `/bench-implement-spec` phase instructs the agent to seed its harness's native
task list from the acceptance coverage map — one task per coverage row, marked
in-progress as the vertical slice turns it red-to-green. The list mirrors
`bench coverage <spec>` (a single source, no drift). The instruction is
harness-neutral prose in the canonical command file; the concrete Claude surface
(the TodoWrite tool) is named only in the Claude adapter README, where
Claude-adapter specifics already live. No hook, no per-harness adapter file, no CLI
change.

## User stories

1. As a reviewer, I want the `/bench-implement-spec` "Then build" step to tell the
   agent to mirror each acceptance-coverage row into its harness's native task list
   (seeded from `bench coverage <spec>`), so I can watch the build's progress
   against the spec's fixed target without asking. The instruction extends the
   existing coverage-row-naming bullet rather than adding a duplicate one, and it
   stays under that bullet's existing "if the spec has an acceptance coverage map"
   guard so a small no-map change is unaffected. The wording degrades gracefully for
   a harness with no native list ("if it has one") and names the row source as
   `bench coverage`, never implying the agent hand-lists rows.
   Line: claude-fable-5 / high. The command file is a leverage artifact loaded by
   every implement-spec session, so a dedupe or drift defect multiplies through
   every build — the profile's skill/command/doc-authoring override routes it to the
   top tier despite being one sentence.

2. As a reviewer, I want `.claude/README.md` to note that `/bench-implement-spec`
   seeds the Claude TodoWrite list from `bench coverage`, so the concrete Claude
   surface is documented where Claude-adapter specifics live rather than leaking a
   Claude-only tool name into the harness-neutral command file.
   Line: claude-fable-5 / high. Same leverage class — adapter documentation the
   agent reads when orienting to the Claude tree — and it must name the source token
   (`bench coverage`) for reviewer cold-read; `.claude/` is not gate-swept, so the
   token is reviewer-checked, not gate-verified.

## Implementation decisions

- **Two files, prose only, no new units, no CLI change.** `bench coverage <spec>`
  already emits the rows and its slug fallback already resolves a spec argument
  (shipped under `implement-spec-lean`). Nothing in `internal/` or `bin/` is
  touched.
- **Extend, don't duplicate (story 1).** The "Then build" section of
  `.agents/commands/bench-implement-spec.md` already carries the conditional bullet:
  "If the spec has an acceptance coverage map, each vertical slice names the
  coverage row it is turning red-to-green before editing that slice." The
  task-list instruction extends *that sentence* — the coverage-row-naming discipline
  has one home; the new clause adds where the naming is surfaced (the harness task
  list), not a second bullet restating the discipline. The load-bearing docs anchor
  phrase "turning red-to-green" must survive the edit verbatim.
- **Harness-neutral canonical file (story 1).** The command file is read by Claude,
  Codex, and other AGENTS.md harnesses. The new clause names surfaces generically —
  "your harness's native task list (Codex plan, Claude todos)" — and must not assume
  the list exists (a plain AGENTS.md harness may have none). No Claude-only tool name
  (`TodoWrite`) appears in this file.
- **Claude surface named only in the adapter README (story 2).** `.claude/README.md`
  gains one short line: `/bench-implement-spec` seeds the TodoWrite list from
  `bench coverage`. This is Claude-adapter-local knowledge that cannot live in the
  neutral command file.
- **No Codex-side file.** `.codex/` holds only `hooks.json` (no README), and the
  Codex adapter skill `.agents/skills/bench-implement-spec/SKILL.md` is deliberately
  content-free (points at the canonical file; adding phase content there violates
  one-source-per-fact). The generic line's own "Codex plan" naming carries the Codex
  case; no new file is created. (Map decision #3, veto at spec.)
- **No new gate anchor.** The task-list clause is ergonomic, not a correctness
  contract, and the story already rides top + high; anchoring it over-fits the gate.
  (Map decision #5, veto at spec.)

## Testing decisions

- **What a good test is here:** there is almost none to write. The core behavior —
  the agent actually seeding its task list from the coverage rows — is agent
  behavior at runtime, not a command output, so it is **not gate-observable**. This
  spec is a prose leverage change; the gate guards *regressions* (existing docs
  anchors stay green, and a mistyped `bench <cmd>` in the command file fails the
  gate's CLI-reference check — the `.claude/README.md` token is reviewer/cold-read,
  `.claude/` being unswept), and **the reviewer enforces the change itself** at
  review and cold-read.
  This is the same posture the `implement-spec-lean` spec recorded for its
  deduped-prose legibility.
- **Cheapest wrong implementation:** a vague sentence that omits `bench coverage` and
  the coverage rows. No gate row goes red on it — the command name stays valid so
  the CLI-reference check passes and the anchors are untouched. That is the honest disclosure: the
  behavior cannot be gate-pinned, which is exactly why the venue is
  `/bench-implement-spec` interactive and **not** `bench shift` (the coverage map is
  not fully gate-observable, failing `craft-line`'s venue-routing test).
- **Seams tested:** the only gate-observable surface is the existing gate over the
  command file — `gate-docs-contracts` anchors and the CLI-reference check over
  `.agents/commands/*.md`. The `.claude/README.md` edit has no gate-observable
  surface (`.claude/` is unswept). No new seam, no new test file.
- **Gate command:** the project gate, `bench gate`.

### Seam diagram

The single gate-observable seam is the existing docs/reference gate over the edited
command file. There is no new unit and no new test-attach point — the seam already
exists; the build must keep it green.

    trigger: `bench gate`  (docs-contracts phase + CLI-reference check)
        │
        ▼
    edited bench-implement-spec.md  ──▶  [ gate-docs-contracts ]  ──▶  anchor "turning
    (`.claude/README.md` has no      ──▶  [ + CLI-reference chk ]      red-to-green" present?
     gate surface: `.claude/` is     ──▶  [ over .agents/cmds   ]  ──▶  `bench <cmd>` tokens
     unswept, reviewer-checked)      ──▶  [                     ]      name a real command?
                      ◀ tests attach here: run `bench gate`; the docs anchor and the
                        CLI-reference check stay green after the command-file edit.
                        (The task-list *behavior* itself has no attach point — it is
                         reviewer/cold-read enforced, recorded not-TDD-able below.)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Agent seeds its harness task list, one task per coverage row, from `bench coverage <spec>` | agent runtime behavior (no command output) | not TDD-able — runtime agent behavior, not gate-observable; enforced at `/bench-review-implementation` and cold-read | no gate signal exists; the review axis and reviewer own it (stated openly, not disguised as TDD) |
| 1 | The `bench coverage` token added to the command file resolves | reverse CLI-reference check in `bench gate` (scans `.agents/commands/*.md` for unknown `bench <cmd>` tokens) | already covered — the check already runs in `bench gate`; a mistyped command name would fail it | a mistyped `bench <cmd>` in the new clause is flagged as an unknown command by the existing CLI-reference check |
| 1 | The docs anchor "turning red-to-green" survives the bullet edit | `gate-docs-contracts` anchor on the command file | already covered — the anchor is checked today; run `bench gate` after the edit | if the edit drops the anchor phrase, the docs-contract check goes red |
| 1 | Clause stays under the "if the spec has an acceptance coverage map" guard (no-map change unaffected) | reviewer/cold-read of the edited bullet | not TDD-able — structural placement, reviewer-checked | a clause hoisted out of the guard would instruct a task list on maps that don't exist; caught at review |
| 1 | Wording degrades for a harness with no native task list | reviewer/cold-read | not TDD-able — prose quality, reviewer-checked | wording that assumes TodoWrite/plan exists would misdirect a plain AGENTS.md harness; caught at review |
| 2 | `.claude/README.md` gains the TodoWrite note and its `bench coverage` token resolves | reviewer/cold-read of the README | not TDD-able — `.claude/` is walked by neither the stale-reference sweep nor the CLI-reference check, so the token is reviewer/cold-read enforced, not gate-observable | a dangling reference in the README is caught at review and cold-read, not by the gate |
| 2 | The Claude-only tool name `TodoWrite` does not leak into the neutral command file | reviewer/cold-read of the command file | not TDD-able — no gate check distinguishes tool names in prose | a harness-specific name in the canonical file breaks neutrality for Codex/other harnesses; caught at review |

### Edge inventory

Edge classes walked per behavior; code-path classes are n/a because this spec
touches no code path.

- **empty/absent input (spec has no coverage map)** — resolved as a coverage row
  (story 1): the clause stays under the existing "if the spec has an acceptance
  coverage map" guard, so a no-map change gets no task-list instruction.
- **hostile environment (harness with no native task list)** — resolved as a
  coverage row (story 1): wording degrades with "if it has one" and never assumes
  the list exists.
- **invocation through every shipped surface** — Won't handle as a distinct case:
  the canonical command file reaches Claude (symlink), Codex (adapter skill), and
  other harnesses (native AGENTS.md read) unchanged by the existing link/symlink
  mechanism; a prose edit introduces no new dispatch.
- **paths/globs in names, trailing newline, tool missing from PATH, symlink
  invocation, SIGINT mid-loop, re-run idempotency, cwd depth** — Won't handle: all
  are shell-CLI code-path classes; this spec adds no code path, so none has a
  surface here.
- **docs anchor sitting inside a rewritten sentence** — resolved as a coverage row
  (story 1): the anchor check is a substring match, so the build re-runs the
  docs-contract fragment after the prose edit to confirm "turning red-to-green"
  still matches.

## Out of scope

- **A gate anchor pinning the task-list clause** — a separate capability (a new
  enforcement contract with its own bite-proof), not the rest of this feature. The
  map recommends against it (ergonomic clause, over-fits the gate); if the reviewer
  later wants the clause gate-protected, it is a small follow-up: `1 edit, 1 gate
  run` to add the anchor plus `1 edit, 1 gate run` to prove it bites (remove the
  clause → red, restore).
- **A `.codex/README.md` naming the Codex plan surface** — a separate documentation
  artifact (the Codex adapter has no README today). The generic line's "Codex plan"
  naming covers the case now; creating the file later is `1 edit, 1 gate run` if
  Codex-adapter docs ever grow enough to warrant a README.
