# Bench platform assessment — 2026-07-08

Deep assessment of adoption/packaging, the workflow commands and skills, the
enforcement layer, the CLI/Go core, the test suite and gate authority, and live
operational state. Produced by six read-only area sweeps (mid tier) synthesized
on the top tier; claims marked ✓ were independently re-verified against source
or live output by the synthesizer. A live full gate run during the assessment
was **green in 31.7s wall** (235s user — the four phases parallelize well).

The 2026-07-07 assessment and its build-out live in git history (`89e03f7`,
built out through `34e870f`); §1 verifies that batch landed rather than
re-listing it. This file replaces it as the current assessment.

Severity: **high** = an invariant or advertised guarantee is not actually held;
**med** = a real defect or reachable unowned state; **low** = friction, drift
risk, or hygiene.

---

## Executive summary

Yesterday's batch held: all 12 backlog items verified against the tree — 11
FIXED, 1 PARTIAL — and every prior high/med finding is closed except the
upstream-blocked FT24 (§1). The gate is green and fast, `go vet`/`go test` are
clean, shellcheck passes all 16 shell files, and the historical gate-flake
thread is empirically closed (119 stress runs, zero flakes, ~300× timing
margin). The kit's *inside* is in the best state of its three assessments.
Today's problems cluster in four themes:

1. **The front door is broken; everything behind it works.** `npm view
   benchkit` resolves to a stranger's unrelated package, every `@benchkit/*`
   platform binary 404s, and the README's clone fallback never says how to
   build the binary — so the advertised primary adoption path installs the
   wrong software and there is currently **no install route a stranger could
   follow**. Meanwhile the hands-on by-path adoption walk graded solid:
   link/init/doctor/gate/unlink all work, including the no-global-`bench`
   story under a fully stripped PATH.
2. **Harness asymmetry is the open enforcement flank.** Codex's Stop hook has
   a 30s timeout while the gate takes ~32s on the reference box *today* — the
   completion oracle gets killed mid-run exactly where it is armed; Claude
   sets no timeout. Codex hooks also resolve paths cwd-fragilely where Claude
   uses `$CLAUDE_PROJECT_DIR`, and the line guard remains Claude-only (parked,
   upstream-blocked). The "identical behavior under every harness" promise
   degrades in that order.
3. **Decision records fell behind the code within a day.** ADR 0002 posture 5
   asserts "nothing short-circuits a gate run on a cache hit" — `bench commit`
   now does exactly that (soundly, but the record is false). The CHANGELOG
   missed both 2026-07-07 learnings-sourced promotions despite its own append
   rule, and no gate check protects it.
4. **Two self-inflicted loops in the worktree layer.** `bench status`
   perpetually recommends `bench worktree clean` for unmerged salvage branches
   that `clean` deliberately keeps — an action that provably no-ops, live on
   the dashboard today. And the lease reclaim path has a real (narrow) race
   where two crash-recovery claimers can both win the same worktree,
   contradicting the code's own advertised guarantee.

---

## Findings by area

### 1. The 2026-07-07 batch: verified

All 12 ranked backlog items checked story-by-story against the current tree
(not against commit messages).

| # | Item | Verdict | Evidence |
|---|---|---|---|
| 1 | Maintenance pass | FIXED | specs/, decisions/, reviews/ empty; IDEAS.md empty; roadmap current |
| 2 | After-merge owner | FIXED | `bench-final-check.md:29-41` post-merge exit duty; HANDOFF.md deleted |
| 3 | First-install fixes | FIXED | `projects/` ships ✓; dist gitignored; scaffolded gate resolves `.bench/bin` first; verified hands-on in throwaway repos |
| 4 | Line-guard fail-open + pre-push verify | FIXED | `lines.go:130-146` denies missing-model; `doctor.go:215-236` pre-push rows; residual rims recorded (ADR 0002 §4) |
| 5 | Worktree branch sweep | PARTIAL | merged orphans swept (`clean.go:39-75`); unmerged-salvage case is §2's loop |
| 6 | False-empty class + models argv | FIXED | `structure.go` propagates git errors; `models.go:43-46` exits 2 |
| 7 | Status valve + invocable actions | FIXED | `status.go` `--all`; actions now name runnable commands |
| 8 | Enforcement-postures ADR | FIXED | `docs/adr/0002` records all six postures (but see §5 — posture 5 since rotted) |
| 9 | craft-spec skill (FT23) | FIXED | `bench-craft-spec` single-sources the coverage-map schema; tdd/review point at it |
| 10 | Structure splits | FIXED | oversized test files split; accept mechanism with reasons |
| 11 | Uninstall story | FIXED | `unlink.go` consumes the manifest; round-trip verified hands-on (but see §2 — CLAUDE.md residual) |
| 12 | Cheap hardening batch | FIXED | newline corpus, stophook seam test, README prereqs/layout |

Prior §2–§7 highs/meds: all FIXED except FT24 (Codex line guard — parked
pending upstream, verdict recorded in BENCH-reference Hook Layers; correctly
tracked, not re-filed). The gate-flake thread (FT39) closed properly: the
event-keyed overlap deadline survived 119 stress runs including 8-way parallel
contention and a live conformance co-load, worst case 0.2s against a 60s
deadline.

### 2. Adoption, packaging, first hour

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| high | publication | **The advertised install path installs someone else's package ✓.** `npm view benchkit` → `benchkit@1.4.3`, "Helper library for verifying and benchmarking algorithmic alternatives" (Maximilianos/benchkit) — the name is taken by an unrelated package. Every binary package 404s (`@benchkit/linux-x64` etc. ✓). README's "fastest way" (`npx benchkit link`, `npm i -g benchkit`) hands a new user the wrong software; even with the right name, the binary fetch would fail. May be a known pre-publication state — but the shipped README advertises it as live today. Reviewer decision: publish under an owned scope, rename, or de-advertise npm and lead with the clone path. | `npm view` output ✓; `package.json:2,44-48`; `README.md:180-192` |
| med | dead-end | **The clone fallback never says how to build ✓.** README names Go as a prerequisite but contains no build command (`scripts/go-build.sh` produces `dist/bench`); a clone user who runs `bin/bench.sh link` unbuilt gets an error whose remedy points at the 404 npm package. | `README.md` (no build cmd ✓); `bin/bench.sh:178` |
| med | residual | **`bench unlink` leaves a dangling CLAUDE.md ✓.** `installClaudeMD` writes CLAUDE.md outside `installPlan`, so it is never manifest-recorded; unlink removes `.bench/BENCH.md` but leaves the CLAUDE.md that imports it. The Claude-first entry file is left broken, and it's absent from README's documented leave-behind list. | `internal/adopt/link.go:286-293` ✓ vs `:295-308`; `README.md:206-207` |
| low | trap | `bench init`'s profile hint ("see projects/<name>.md in the Bench kit") doesn't say where the kit is; nothing copies `projects/` into the consumer repo. | `internal/adopt/init.go` hint |
| low | wording | Message nits: `bench coverage` (no arg) says "unknown argument" for a *missing* one; `bench link` in a non-git dir advises "run inside a Bench-linked repo" — link's job is to create the linkage. | observed output |

Clean: tarball shape is correct (66 files; no HANDOFF/CLAUDE/AGENTS/Go-source
leakage; `projects/` ships); relink idempotent; conflict handling atomic with
zero partial state; pre-push bites; space-in-path, no-commit repo,
subdirectory invocation all behave; no-global-`bench` story holds end-to-end
under a PATH stripped of both `bench` and `node`.

### 3. Enforcement and harness parity

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | fail-open | **Codex Stop hook `timeout: 30` can kill the completion oracle ✓.** An armed session (`BENCH_SHIFT=1`) runs the full gate in the Stop hook; the gate measured 31.7s wall on this box, so Codex's 30s timeout kills it — almost certainly fail-open, letting an agent stop on red. Claude's Stop block sets no timeout. Invariant 1's interactive layer silently degrades under exactly one harness. (Codex kill semantics assumed, not executed.) | `.codex/hooks.json:9` ✓ vs `.claude/settings.json`; `internal/stophook/stophook.go:86-98` |
| med | fail-open | **Pre-push can guard a fabricated branch ✓.** `git.DefaultBranch` falls back to `"main"` whenever `origin/HEAD` is unresolvable ✓; `installGitHook` bakes that answer into the hook at link time. A repo whose real default is `master`/`trunk`, linked before a remote exists, gets a backstop that never fires — silently. | `internal/git/git.go:104-110` ✓; `internal/adopt/link_hook.go:90-105`; `internal/adopt/prepush.sh:43` |
| low | fragility | Codex hooks resolve via `$(git rev-parse --show-toplevel)` — empty outside a repo → hook not found → fail-open; Claude uses `$CLAUDE_PROJECT_DIR`, robust to cwd. | `.codex/hooks.json:8,21` ✓ |
| low | hardening | `prepush.sh`'s read loop drops a newline-less final ref line (no `|| [ -n "$line" ]` tail). Unreachable via git today, which always LF-terminates — defense-in-depth for a security-critical loop. | `internal/adopt/prepush.sh:27` |
| low | drift-risk | The git guard deliberately inlines its own wrapper search order (sourcing would add a fail-open mode) — but the two copies have no sync test. | `.bench/hooks/block-dangerous-git.sh:34-43` vs `.bench/lib/resolve-bench.sh:16-27` |
| low | residual | Doctor verifies the pre-push hook, but no `bench status` signal fires when it's missing — a fresh clone loses the backstop until someone runs doctor. | `internal/adopt/doctor.go:215-236` |

Held under adversarial reading: phase aggregation fails closed (no path
returns green on start-failure or signal death); the canary is vacuity-guarded
and catches unwiring of any fixture-backed check; the git guard, adapters, and
shift model-guard postures all match their advertisements; shellcheck clean
across all 16 shell files with zero real positives at `-S style -o all`.
Accepted postures re-confirmed and correctly cited rather than re-filed:
family-level canary granularity and single-check unwiring (ADR 0002 §6, also
parked in FT6), line-guard residual rims (§4), interactive commit-on-red (§1).

### 4. CLI and Go core

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | race | **Two concurrent lease reclaimers can both win a worktree ✓.** `Claim`'s takeover rename is unconditional on lease identity: reclaimer A can complete rename→remove→re-create, then stalled reclaimer B renames A's *fresh* lease and also returns true — falsifying the comment's "two concurrent reclaimers cannot both win" ✓. Both then `reset --hard` and write into the same worktree. Reachable only on crash-recovery (dead-pid lease, two simultaneous `Acquire`s), window is narrow, and the concurrent-acquire contract covers fresh-mint only — this interleaving is untested. | `internal/worktree/lifecycle.go:87-105` ✓ (guarantee at `:82-86`), reached via `:136-148` |
| med | loop | **Status recommends a provably no-op action for salvage orphans ✓.** Unmerged `worktree-*` branches are deliberately kept by the sweep ("inspect or delete by hand" ✓) yet the status/worktree row perpetually says `→ bench worktree clean`. Live today: the 2 flagged branches are salvage commits *created by* `clean` itself; a reviewer who follows the action sees nothing change, forever. FT28 closed the merged-orphan case; this is the loop that outlived it. | `internal/worktree/clean.go:63-66` ✓; `internal/status/status.go:287` |
| low | 1-source | `bench guards` hardcodes the pre-push marker literal that `adopt.prePushMarker` documents itself as the one source of ✓ — change the marker and guards silently reports a managed hook as unmanaged. No test ties the two. | `internal/guards/guards.go:191` ✓ vs `internal/adopt/link_hook.go:31` ✓ |
| low | 1-source | The gate-cache filename `bench-last-gate` is a copied literal in writer and reader ✓ (different packages, no shared const) — rename one side and the protocol breaks silently: no gate row, no commit reuse, no error. | `internal/gate/gate.go:208` ✓; `internal/status/status.go:182` ✓ |
| low | false-empty | `RegisteredWorktrees` discards the `git worktree list` classify error — a git failure reads as "no worktrees" (clean prints nothing-to-clean at exit 0). The last surviving instance of the class FT29 swept. | `internal/worktree/classifier.go:26` |
| low | exit-codes | Unknown subcommand prints help at **exit 0** ✓ — a typo is indistinguishable from success to a wrapping script; `--version` and `--help` fall into the same `*)` case ✓. Every real subcommand rejects bad args at exit 2; the dispatcher is the outlier. | `bin/bench.sh:232-262`; `frobnicate`/`--version` → exit 0 ✓ |
| low | feedback | `bench canary` passing run prints nothing at exit 0 — an oracle command with zero success feedback ("did it run?"); `bench structure` prints "structure ok". | observed output |
| low | drift | Stale comment: `git.DefaultBranch` calls itself "the Go mirror of bench.sh's `default_branch`" — no such shell function exists anymore ✓. | `internal/git/git.go:102-103` ✓ |
| low | hardening | `resolve_script_path`'s hand-rolled symlink chase has no loop cap — a circular symlink to the wrapper hangs the CLI (readlink -f would detect it). | `bin/bench.sh:37-41` |

Clean: `go vet ./...` and `go test ./...` green including contract suites
against a freshly rebuilt `dist/bench`; structure/spec/outline/dashboard/
models/canary/stophook/toon/unlink all read sound; FT29's git-error-honesty
fix is real (errors propagate; the one tolerate-site annotated); outline is
symlink-cycle-safe by construction (`git ls-files` enumeration, `Lstat`
regular-file filter); dashboard HTML is sanitized and atomically written.

### 5. Gate authority, tests, records

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | stale record | **ADR 0002 posture 5 is now false ✓.** It asserts "nothing short-circuits a gate run on a cache hit, so a stale or forged cache cannot lie a tree green where it counts" — `bench commit` (d874629) now reuses a fresh green verdict for an identical tree, skipping the gate. The reuse itself is sound (exact tree-hash key, single writer refusing non-hash trees, both directions regression-tested) — but the decision record misdescribes shipped behavior, which is precisely what invariant 3 forbids. Amend the posture, don't revert the feature. | `docs/adr/0002:42-49` ✓ vs `internal/commit/commit.go:105-114` ✓ |
| med | drift/duty | **CHANGELOG missed both 2026-07-07 learnings-sourced promotions ✓.** Its header mandates one entry per promotion; `453599a` (FT35/36 rule edits) and `2a72310` (verification-probe rule) have none — newest entry is 2026-07-06 ✓. The CHANGELOG is `/bench-update-kit`'s re-synthesis baseline, so the next upstream sync diffs against a record missing two adopted rules. No gate check reads CHANGELOG — nothing protects this duty. | `CHANGELOG.md:9-11` ✓; commits `453599a`, `2a72310` |
| low | coverage | Commit cache-reuse is pinned for normal-file staleness but not for the capture-only allowlist paths (`ROADMAP.md` etc.) — current code is correct (the allowlist softens only the dashboard message, never `Stale`), but a future "optimization" wiring the allowlist into reuse would go uncaught. | `internal/status/status.go:202,228-238`; `internal/contract/runtime/runtime_commit_test.go:78-92` |
| low | accuracy | The concurrent-acquire test's comment claims "a test-owned barrier, not a timed poll" — release is barriered, but overlap *detection* is still an event-keyed 60s deadline, and the shell's self-timeout (~60s) equals that window. Empirically fine (0.2s observed); the comment overstates and the coupling is undocumented. | `internal/contract/runtime/runtime_worktree_test.go:304,315,341,366-368` |
| low | platform | Outline symlink-safety and the git-guard hook tests `t.Skipf` on filesystems lacking the capability — the gate passes on such platforms without exercising the hostile-input class. Covered on the Linux gate box. | `internal/contract/axi/axi_outline_test.go:142` |

Seam coverage: every seam in `projects/benchkit.md` has its contract tested,
including all seven features shipped since the last assessment (spec-history
exact-token cut, dashboard, outline, status valve, unlink reversal, commit
verdict reuse, orphan-branch sweep). No high or med gate finding: **no
constructed path makes the gate lie green.**

### 6. Docs, skills, vocabulary

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| low | drift | CONTEXT.md's canonical "signal" definition names 6 signal types; `bench status` emits 10 (guards, drain, specs, reviews, roadmap added; "uncommitted" is labeled `git`). The file that pins the ubiquitous language teaches half the board's vocabulary. Unprotected by any gate check. | `CONTEXT.md:57-58` vs `internal/status/status.go:261-429` |
| low | drift | "Dashboard" now names two things: CONTEXT defines *ambient dashboard* = `bench status`, but the `bench dashboard` HTML artifact has no canonical term — and FT38 (open, next in the recommended sequence) is about exactly that artifact. | `CONTEXT.md:54-56`; `ROADMAP.md:13-17` |
| low | drift | README's curated `internal/` layout omits `dashboard/` and `outline/` — the packages behind two commands its own CLI list names. | `README.md:142-155` |
| low | drift | The lighter-path threshold is still worded three ways (BENCH.md "a few-line change" / write-spec "more than a trivial change" / implement-spec "the seam is obvious"). Same low as the last two assessments; shrinking, not closed. | `.bench/BENCH.md:150`; command docs |
| low | shipping | `projects/benchkit.md` — this repo's internal dogfood profile, naming its Go seams — ships to every consumer via the `projects/` glob alongside the two intended examples. | `package.json:12` |
| low | orphan | `research/unit_testing.pdf` is referenced nowhere and isn't shipped — unowned scratch. (Contrast `ui_example/`, alive as FT38's design reference.) | `research/` |
| low | recurring | Third consecutive assessment with no owner for the assessment itself — no `craft-assess` skill or `/bench-assess` phase names the drill (verify last drain, sweep areas, rank backlog). The cadence is now established fact, and each run re-derives the method from the prior file. | (no owner) |

Always-loaded cost: **≈4,400 tokens/turn**, flat vs prior ≈4,300 (the 13th
skill added ~100). Only marginal shaves exist (~200 tokens across five
descriptions with redundant closers); BENCH.md's dominance is load-bearing by
recorded reviewer decision. Near-irreducible under current design.

Clean: skills index in sync (`--check` exits 0); all 13 craft + 10 adapter
skills valid, `.claude/skills` all symlinks; no dead command/flag/path
references anywhere in commands or skills (verified against dispatch); ADRs
path- and snippet-free; specs/decisions/reviews all at zero; write-spec's
batch-drain override recorded; craft-spec single-sources the coverage-map
schema with no remaining overlap.

### 7. Live operational state

- **Gate green, 31.7s wall ✓** — the oracle is not a friction point.
- **1 open learning** (FT41 dogfood shortfall — a flagged judgment call, no
  proposed rule) awaiting `/bench-what-next`; IDEAS.md empty.
- **2 orphaned salvage branches** (`worktree-agent-a137764f…`, `…a5ce49b2…`),
  both unmerged, both flagged with the no-op action of §4's loop finding. The
  guard now permits `git branch -D worktree-agent-*`, so after inspection they
  are hand-deletable.
- **Roadmap honest**: FT38 (next: `/bench-shape-idea` grill), FT6 (parked
  pending evidence), FT24 (parked pending upstream), FT8 (time-boxed) — no
  stale rows.

---

## Cross-cutting themes

1. **The kit is now better than its front door.** Everything a user meets
   *after* installation graded solid; the npm identity, the binary packages,
   and the missing build docs mean nobody but the author can reach that
   experience. Yesterday's theme — dogfooding masks adoption defects — closed
   one level down (fresh link/init now work) and re-appears one level up: the
   kit repo never installs itself *from the registry*, so the registry story
   was never exercised.
2. **Records rot faster than code, and only gate-checked records hold.** ADR
   0002 posture 5 and the CHANGELOG both fell behind within a day of being
   written, while every gate-anchored doc (skills index, CLI inventory, tier
   binding, shared-rule single-sourcing) stayed true. The lesson is already in
   the repo's own architecture: a duty without an enforcement anchor or a
   phase-owned checklist slot is a hope.
3. **Signals must recommend actions that work.** The status board's authority
   depends on its actions being executable and effective; a row whose remedy
   provably no-ops (salvage orphans → `worktree clean`) trains the reviewer
   to ignore the board. Same class as last assessment's dead-end remedy — the
   fix pattern (make the CLI own the remedy, or change the recommendation)
   should be applied whenever a keep-branch is introduced.
4. **Advertised guarantees need adversarial tests, not comments.** The reclaim
   race hides directly beneath a comment asserting it can't happen; the
   overlap "barrier" comment overstates what the test does. Where code
   *advertises* a concurrency guarantee, a contract test should try to break
   it — the concurrent-acquire suite is the natural home for a two-reclaimer
   stress case.

## Ranked improvement backlog

Ordered by platform leverage; sizes are rough.

1. **Resolve the npm identity** — reviewer decision: publish under an owned
   scope (e.g. `@gibbonmi/benchkit`), rename, or de-advertise npm and make
   the clone path primary. Whichever way: add the build command
   (`scripts/go-build.sh`) to README's clone path, and fix the broken-binary
   error message that points at the 404 package. Until this, every external
   adoption fails at step 1. (decision + S)
2. **Fix the Codex Stop timeout** — raise or drop `timeout: 30` in
   `.codex/hooks.json` (kit copy and the linked template) so the armed-shift
   oracle can't be killed mid-gate; consider a margin rule (≥ 10× gate wall).
   (XS)
3. **Close the salvage-orphan loop** — either `bench worktree clean` gains an
   explicit salvage disposition (e.g. a listed `--purge-salvage` after
   inspect) or the status/worktree row distinguishes unmerged salvage and
   recommends the inspect-then-delete path. Also clears today's two live
   branches. (S)
4. **Fix the reclaim race** — make `Claim`'s takeover verify identity (e.g.
   re-check the renamed stale file matches the lease content it judged
   reclaimable, or take a per-entry `O_EXCL` claim lock), and add a
   two-reclaimers-one-stale-lease stress case to the concurrent-acquire
   contract. (S–M)
5. **Amend ADR 0002 posture 5** to record the commit verdict-reuse decision
   (exact-tree-hash-keyed, fresh-only, single writer), and pin the
   capture-only-allowlist reuse regression while there. (S)
6. **Record CLAUDE.md in the link manifest** (only when link created it) so
   unlink removes it; add it to README's leave-behind/removal list. (S)
7. **Backfill the two missing CHANGELOG entries and give the duty an owner** —
   the cheapest honest anchor is a step in `/bench-what-next`'s drain
   checklist ("a promoted rule appends its CHANGELOG entry in the same
   diff"); a conformance check is the stronger option if drift recurs. (S)
8. **Pre-push default-branch honesty** — don't bake a fabricated `main`:
   resolve the branch inside the hook at push time, or warn loudly at link
   when `origin/HEAD` is unresolvable and add a doctor row comparing the
   baked branch to the live default. (S)
9. **One-source collapse batch** — export/import the pre-push marker const in
   `guards`; a shared const for `bench-last-gate` (natural home: `internal/
   git`); a sync test for the guard's inlined wrapper search order; delete
   the stale `default_branch` mirror comment. (S)
10. **CLI hygiene batch** — unknown subcommand → usage at exit 2; route
    `--version`/`--help`; one-line `bench canary` success output; fix the
    `coverage`/`link` message nits; decide the harness-form posture for CLI
    strings that print `/bench-*` (accepted Claude-form vs harness-neutral
    phrasing) and record it. (S)
11. **Docs batch** — CONTEXT.md: full signal list + a canonical term for the
    `bench dashboard` artifact (feeds FT38's grill); README `internal/`
    layout + build command; collapse the lighter-path threshold to one
    wording; decide `projects/benchkit.md` shipping; disposition
    `research/unit_testing.pdf`. (S)
12. **Test/hardening batch** — prepush read-loop newline guard; symlink-loop
    cap in `resolve_script_path`; propagate the `RegisteredWorktrees`
    classify error; fix the overlap-comment overstatement; a status signal
    for a missing pre-push hook. (S, mechanical)

Operational (zero build cost): run `/bench-what-next` to drain the FT41
learning; inspect and hand-delete the two salvage branches (guard now permits
it). Recurring, third flag: if the assessment cadence continues, give it an
owner (`craft-assess` or `/bench-assess`) so the method stops being re-derived
from the prior file.

Parked/known, re-confirmed and correctly not re-filed: FT24 (Codex line
guard, upstream-blocked), FT6 (per-anchor canary needles — also ADR 0002 §6),
FT8 (tier revisit), FT38 (dashboard identity — the CONTEXT term gap above
feeds it).

## Verification notes

- Synthesizer-verified (✓ above): `npm view` outputs for `benchkit` and
  `@benchkit/linux-x64`; README build-command absence; `frobnicate`/
  `--version` exit 0; the `Claim` race interleaving read directly from
  `lifecycle.go`; `.codex/hooks.json` timeout vs Claude settings;
  `DefaultBranch` fallback; `clean.go` kept-branch vs the status action; ADR
  0002 posture 5 text vs `commit.go`; `installClaudeMD` outside the manifest
  plan; both one-source literals; CHANGELOG head vs the two promotion
  commits; live ROADMAP/learnings/status state.
- The gate run, flake stress (119 runs), `go vet`/`go test`, shellcheck
  sweep, and all hands-on adoption flows were executed by delegates this
  session, not merely read.
- Known coverage limits (unknowns): Codex's kill-on-timeout semantics assumed
  fail-open, not executed; the reclaim race proven by construction, not
  stress-reproduced; `guards.prePushRow` under a configured `core.hooksPath`
  (ClassifyPrePush compensates; guards may not) reasoned only; live-harness
  hook runtime (real Stop/PreToolUse events) read statically; native-Windows
  and registry-install behavior not executable (no published package);
  canary EXPECT cross-check collision judged implausible, not proven.
