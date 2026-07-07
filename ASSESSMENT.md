# Bench platform assessment — 2026-07-07

Deep assessment of the workflow commands, skills, enforcement, CLI/Go core,
installation/adoption surface, and live operational state. Produced by six
read-only area sweeps (mid tier) synthesized and spot-verified on the top tier;
claims marked ✓ were independently re-verified against source by the
synthesizer. A live full gate run during the assessment was **green in 29.6s
wall** (3m44s user — the phases parallelize well).

The 2026-07-06 assessment and its drain disposition live in git history
(`72e6a7f`, drain `b2b7dc6`); §1 verifies that drain landed rather than
re-listing it. This file replaces it as the current assessment.

Severity: **high** = an invariant or advertised guarantee is not actually held;
**med** = a real defect or reachable unowned state; **low** = friction, drift
risk, or hygiene.

---

## Executive summary

Yesterday's drain held: 9 of 10 specs verified HELD with gate enforcement, 1
PARTIAL only by its own disclosed posture (§1). The enforcement core survived
adversarial reading, and the oracle is fast. Today's problems cluster in four
themes:

1. **The workflow ends at merge; the tree keeps going.** Every phase hand-off
   up to `/bench-final-check` is owned; nothing after it is. Live evidence: 10
   merged specs unretired, 9 stale roadmap rows, 19 decision maps never swept,
   a stale `HANDOFF.md` giving cold sessions day-old instructions, and 8
   orphaned scratch branches whose recommended remedy the kit's own git guard
   denies to the agent.
2. **Adoption breaks where the home repo never walks.** The kit dogfoods
   itself in a repo where everything is already installed. A fresh
   `npm install` is missing the profile templates `bench init` points at, a
   committed `.bench/dist` binary can silently disable the Stop-hook oracle
   for a different-arch teammate, the scaffolded gate assumes a global `bench`
   on PATH, and there is no repo-level uninstall.
3. **Remaining enforcement gaps are edges — several deliberate but
   unrecorded.** The agent-line guard fails open when a delegation omits the
   model field (the exact silent-escalation path invariant 2 targets);
   pre-push installation is never re-verified after link (fresh clones lose
   the backstop silently); the accepted honor-system residuals live in a
   drained assessment table, not a decision record.
4. **The false-empty defect class outlived its fix.** FT19 fixed
   false-empty-on-git-failure where it was found (`bench diff`); the same
   pattern sits in `bench structure` today. The class was fixed as instances,
   not hunted as a class.

---

## Findings by area

### 1. The 2026-07-06 drain: verified

All ten specs (`Status: implemented`, merged fe5863f..f89aee0) checked
story-by-story against the tree.

| Spec | Verdict | Note |
|---|---|---|
| artifact-lifecycle | HELD | `bench spec retire` + orphan-pickup signal + phase-doc duties, all gate-anchored |
| claude-hook-conformance | HELD | Stop/SessionStart/Bash checks at Codex parity; canary needle bites |
| cli-contract-accuracy | HELD | all 8 stories AXI/anchor-pinned |
| docs-drift | PARTIAL | delivered; stories 1–3, 6 enforced only by reviewer cold-read — the spec's own disclosed posture |
| guards-wiring | HELD | `wired` column derived from harness configs; both pre-push postures pinned |
| non-claude-model-tiers | HELD | safe-token grammar, owner-id binding, multi-source discovery — gate-enforced |
| one-source-collapses | HELD | `toon.NotInRepo` collapse enforced; stories 3–5 review-only |
| review-after-merge | HELD | `bench diff --commit` + step-1 fallback, AXI-pinned |
| roadmap-reconcile | HELD | status signal enforced (and firing today — see §7) |
| shellcheck-coverage | HELD | delivered; bite-proof lives in a Go unit test, not the `tests/canary/` family the acceptance row named — enforcement intact, seam substituted |

Spot-checks: `reviews/ft9.md` deleted as promised (`reviews/` empty ✓); stale
FT4/FT9 roadmap rows gone ✓. Residual: the batch itself now sits in the §2
post-merge gap — specs unretired, rows stale.

### 2. Post-merge lifecycle and workflow commands

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| high | gap | Spec retirement has no reachable post-merge owner. The kit advertises promote-then-delete ("a row leaves when the work ships"), and 10 merged specs sit unretired today ✓. The only prose instructing it is `/bench-write-spec` step 8 — a phase gated on a *new* feature's decision map — plus the `/bench-what-next` backstop; `/bench-final-check`, the last phase a build actually reaches, never mentions it ✓. Retirement happens only if the reviewer notices the status row and hand-runs `bench spec retire` per spec. | `bench-write-spec.md:160-165`; `bench-final-check.md` (no mention ✓); `bench-what-next.md:27` |
| med | gap | **No after-merge phase — the root cause.** Post-merge duties (spec retire, roadmap-row removal, decision-map sweep, handoff refresh) are scattered onto phases reached only for other reasons. `bench status` already computes the duty list; no command consumes its rows as a queue. One status-driven after-merge duty pass closes this finding, the high above, and the two below. | `bench-final-check.md:15-27`; `bench-what-next.md:24-29` |
| med | trap | `HANDOFF.md` is a stale cold-session trap ✓: dated 2026-07-06 pre-drain, it instructs "publish the four local commits" (long since on main ✓), names 2 open learnings that were since resolved (today's 2 are different ✓), and counts 7 structure issues vs today's 4 ✓. No command reads, refreshes, or expires it; it is also shipped to consumers in the npm package (§3). | `HANDOFF.md:14-31,75-94`; `package.json:17` ✓ |
| med | gap | `decisions/` has no lifecycle end ✓: 19 maps accumulate, most for long-shipped work (the go-port series, ft4, ft9, the just-closed non-claude-model-tiers). The status decisions signal counts only *unresolved* maps, so a shipped map is invisible but permanent — asymmetric with specs' promote-then-delete. Whether to sweep or keep closed maps is decided nowhere; "decision map" is absent from CONTEXT.md's core terms ✓. | `decisions/` (19 files ✓); `internal/status/status.go` decisions row; `CONTEXT.md` |
| med | gap | Roadmap rot has no independent trigger: FT13–FT21 rows still read "Staged" for merged work ✓. `/bench-what-next` fires only when ideas/learnings are non-empty; the FT16 signal does fire — but it is sev 9, and today it is exactly the row the five-row budget hides behind `+1 more` ✓ (§7). Downstream of the high finding: retirement removes the rows. | `ROADMAP.md:13-19`; `internal/status/status.go:285-292` ✓ |
| low | drift | `/bench-write-spec`'s entry contract ("refuse without a closed map") is contradicted by the live tree it produced (10 map-less specs under an explicit reviewer batch override). The reconciliation is captured as an open learning with a proposed contract clause — pending drain, not yet a rule. | `bench-write-spec.md:26-30`; `.bench/learnings.md:39-55` ✓ |
| low | friction | Status action strings recommend non-invocable surfaces: specs row → `promote-then-delete (spec-retire)`, decisions row → a skill name — against the platform's own "recommend in the form this harness can invoke" rule. | `internal/status/status.go:217,280` |
| low | drift | The lighter-path threshold for small changes is worded three divergent ways across command docs. | `bench-implement-spec.md:22`; `bench-write-spec.md:2`; `AGENTS.md` |

Clean: the build chain (shape → spec → implement → review → final; debug →
final) is internally coherent — every exit handoff names the next phase and the
red/stop-short routes. The `reviews/` pickup lifecycle is fully owned post-FT13
(create, consume, delete, retire-sweep). `bench-debug`'s retired-spec recovery
correctly owns the consequence of promote-then-delete.

### 3. Installation, packaging, invocation

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | gap | `projects/` is not in `package.json` `files[]` ✓, yet `bench init` prints "see projects/<name>.md in the Bench kit for the profile template" and `/bench-setup-repo` builds from "the example profiles in the kit." On every npm/npx install — the advertised primary path — the templates don't exist; the agent improvises the profile. | `package.json:12-28` ✓; `internal/adopt/init.go:74`; `bench-setup-repo.md:136` |
| med | gap | `bench link` copies the running arch-specific binary to `.bench/dist/bench` with no `.gitignore` entry or guidance. A consumer who commits it hands a different-arch teammate a broken local CLI, and the Stop-hook oracle then fails open — enforcement silently off. | `internal/adopt/link.go:76-78`; `.bench/hooks/stop.sh:58-61` |
| med | gap | The scaffolded `.bench/gate.sh` calls bare `bench canary`, but the gate subprocess env adds nothing to PATH — a machine relying only on the `.bench/bin` local CLI (the documented "global bench optional" story) gets a gate that cannot go green. | `internal/adopt/init.go:101,108-112`; `internal/gate/gate.go:130-139` |
| med | gap | No repo-level uninstall: the link manifest enumerates every installed path, but nothing consumes it to remove the footprint (`.bench/`, `.agents/`, harness adapters, pre-push hook, AGENTS.md block). Only global CLI/shim removal is documented. | `internal/adopt/adopt.go:15-25`; `README.md:187-193` |
| low | waste/trap | npm ships the kit's own `HANDOFF.md` (stale, §2), `CLAUDE.md`, and `AGENTS.md`, all unused by adoption — `link` generates its own from constants. Dead weight, and the stale handoff is exported to every consumer. | `package.json:15-17` ✓; `internal/adopt/marker.go:135-137` |
| low | docs | Undocumented prerequisites: the from-a-clone path needs Go or node to obtain a binary (README doesn't say); Windows is unsupported and unmentioned (no WSL note, `os: [darwin,linux]` errors EBADPLATFORM); nvm/asdf global installs break on `nvm use` beyond what the shim fixes. | `README.md:177,184-187`; `package.json:42-45` |
| low | drift | README's `.bench/` layout omits shipped/managed pieces (`lib/`, `bin/`, `dist/`, `BENCH-reference.md`, `gate.sh`, `lines.env`). | `README.md:113-123` |

Clean: link idempotency and conflict handling are solid (sha256 manifest,
atomic writes, modified-managed refusal, fence-aware AGENTS.md block); npx
ephemeral-cache detection forces copy mode; release mechanics single-source
os/cpu/version from `platforms.json`; postinstall never fails an install;
harness degradation for non-hook harnesses is documented in BENCH-reference.

### 4. Enforcement residuals

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | gap | The agent-line guard fails **open** when the delegation envelope carries no model field ✓ — documented as the intended rim ("every degraded branch is fail-OPEN"), but an omitted model inherits the invoking session's model, which is precisely the silent-escalation path invariant 2 exists to stop. The counter-rule (always pass the bound alias) lives only in `craft-delegate` prose. Deny-on-missing, or record the fail-open as a decision. | `internal/lines/lines.go:113-127` ✓ |
| med | gap | Pre-push has no install verification or self-heal: `bench doctor` checks the PATH shim but never that `.git/hooks/pre-push` exists, is bench-managed, or matches the embedded template ✓ (zero hook mentions in doctor.go). A fresh clone (git doesn't clone hooks), an un-relinked repo, or a global `core.hooksPath` silently loses the harness-independent backstop, and no surface flags it. | `internal/adopt/link_hook.go:26-55`; `internal/adopt/doctor.go` ✓ |
| med | gap | The interactive line guard remains Claude-only (`.codex/hooks.json` wires Stop + git guard but no Agent matcher; OpenCode/plain terminals have no hooks). Known and parked (FT24, pending Codex hook-capability research) — re-confirmed still open, recorded here for completeness. | `.codex/hooks.json`; `ROADMAP.md` FT24 |
| low | gap | Canary coverage is family-level, not per-check: `checkBenchShRoutes` is wired into conformance with no canary fixture — unwiring it escapes the canary layer (the direct-call Go bite test stays green regardless of wiring). Decide whether per-check needles are the standard for wiring-sensitive checks. | `internal/conformance/checks_test.go:34,62,86-106` |
| low | posture | Accepted honor-system residuals (interactive commit-on-red, non-shift done-claims, "declare the line" having no enforcement surface anywhere) are recorded only in the prior assessment's disposition table — now history, not a decision record. One short ADR would make the accepted risk citable. | `docs/adr/` (absent); prior disposition |
| low | perf/posture | The gate verdict cache is write-only ✓ — `gate.Record` writes `<git-dir>/bench-last-gate` but nothing reads it to short-circuit a byte-identical tree, so every Stop re-pays ~30s. Deliberate recompute-always may be the right oracle posture — but it's a decision worth recording, not an accident. (The cache is also forgeable, but only the advisory dashboard trusts it — Stop and shift always re-run — so it cannot force a false done.) | `internal/gate/gate.go:193-209` ✓; `internal/status/status.go:108-150` |

Held under adversarial reading: no-forged-verdict rims (non-hash tree refused,
fail-open paths write no cache); stale cache can't lie green; the Bash git
guard denies all `git push` including `--no-verify`; adapters fail closed in
routed repos; the shift adapter refuses an undeclared or unbound `BENCH_MODEL`;
owner-defined tier ids bind exactly with no provider lookup.

### 5. CLI and Go core

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | defect | `bench structure` converts git failure into false-clean ✓: `git.Output` errors are discarded in both the all-files and `--since` paths, so a bad ref or failed `ls-files` prints "no tracked source files to check" at exit 0. Same class FT19 fixed in `bench diff`; hunt it as a class (audit every `out, _ := git.Output` in porcelain). | `internal/structure/structure.go:44,233` ✓ |
| low | drift | `bench models` ignores its argv — no `-h/--help`, and unknown args emit the full inventory at exit 0. Every sibling porcelain rejects unknown args at exit 2; lone outlier against the AXI usage norm. | `internal/models/models.go:40` |
| low | hygiene | The safe-token corpus omits the newline class — the one input where a `$`-anchored regex would leak under other regex engines. Go RE2's `$` is end-of-text so the grammar is safe today, but the corpus doesn't prove the property it most needs. One corpus row. | `internal/modelid/modelidtest/corpus.go:19-32` |
| low | seam | `stophook.Run`'s I/O path (gate exec, rc→verdict map, the rc==3 no-gate branch that blocks the stop) has no Go-level test — only pure helpers are unit-tested; the armed-shift branch is exercised nowhere cheap. | `internal/stophook/stophook.go:86-125` |
| low | verdicts | Structure-budget verdicts on the 4 flagged files: `axi_wave2_test.go` (639) — genuine split along the command boundary; `line_routing_checks_test.go` (423) — real static-parse vs subprocess-exec seam, minor; `runtime_status_test.go` (402) and `status.go` (415) — budget noise, don't fragment. | `bench structure` output |

Clean: gitguard tokenizer (table-tested against control ops, wrapper recursion,
malformed quotes), lines/modelid (exhaustive edge tables, one shared corpus
across both consuming seams), models (definitive empty states, advisory exit-0
posture), canary (vacuity baseline, unique-fixture check), adopt internals
(atomic writes, fingerprint semantics), subprocess/terminal/packagesurface.
All packages vet-clean and green.

### 6. Skills and always-loaded guidance

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | reachability | The acceptance-coverage-map schema and edge-inventory classes are owned by the `/bench-write-spec` command file, which model-invoked skills cannot auto-load: `craft-tdd` and `craft-review` point at it but an ad-hoc TDD pass or self-review fires without reach to the row schema. Known — parked as FT23 (`craft-spec` skill); this assessment re-confirms it's the structural gap, worth unparking. | `craft-tdd:73-76`; `craft-review:33-40`; `ROADMAP.md` FT23 |
| low | cost | Standing context cost ≈4,300 tokens/turn: always-loaded prose ≈3,300 plus 12 craft-skill descriptions ≈1,020. Shave candidates: the three longest descriptions (synthesis, adr, seams) each carry a redundant third trigger clause; the CLI Inventory in BENCH.md is the largest always-loaded lookup block — its demotion to BENCH-reference is a deliberate-tradeoff reviewer call (cold pickup must not depend on non-loaded files). | `craft-synthesis:3`; `craft-adr:3`; `.bench/BENCH.md` CLI Inventory |
| low | reachability | The gate command and seam list live only in `projects/benchkit.md`, reached only if the session follows BENCH.md's "read the profile" pointer. Mitigated (guards --brief injects the tier binding; the CLI inventory names `bench gate`) but it is the biggest know-to-look dependency for a cold session. | `.bench/BENCH.md:38-40`; `.bench/hooks/session-start.sh:37-40` |
| low | pending | `craft-delegate`'s charge template lacks the stale-base opener ("run `git merge --ff-only main`, verify HEAD, stop if denied") that a real batch build needed and adopted mid-run — captured as an open learning, awaiting drain. | `.bench/learnings.md:25-37` ✓ |
| low | coverage | The platform-assessment judgment itself (this task, second run at this cadence) has no owner — `craft-synthesis` folds one change, no skill or command owns "audit the whole surface, verify the last drain, rank the backlog." If the cadence repeats, a `craft-assess` skill or `/bench-assess` phase earns its place. | (no owner) |

Clean: all 12 `.claude/skills` entries are symlinks to `.agents/skills` ✓;
skills index in sync (`--check` exits 0); no dead references; progressive
disclosure holds (heavy detail behind `references/`); no two skills own the
same fact (line/delegate, seams/tdd/review splits are clean); phase adapters
still thin; trigger phrasing names concrete situations and quoted phrases.

### 7. Live operational state (as of this assessment)

Not defects — the queue the platform itself is flagging, plus what the flags miss:

- **10 merged specs awaiting `bench spec retire`; 9 stale roadmap rows.** The
  reconcile mechanism works; the pass hasn't been run. One maintenance session
  clears both plus the two open learnings.
- **2 open learnings** (stale-base delegate charges; the write-spec batch-drain
  override) — both carry concrete proposed rule changes; drain-ready.
- **8 orphaned `worktree-agent-*` branches, and the remedy is a dead end for
  the agent ✓:** the status row says "delete scratch branch," but no bench
  subcommand deletes them (`bench worktree clean` handles worktrees only;
  `OrphanedDelegateBranches` is detection-only ✓) and the kit's own
  block-dangerous-git hook denies `git branch -D` to the agent ✓. The reviewer
  must hand-delete eight branches; `bench worktree clean` should own the sweep.
- **The dashboard hides its sixth signal ✓:** six signals fire, the five-row
  budget truncates to `+1 more`, and `bench status` has no flag to expand —
  today the hidden row is exactly the FT16 roadmap-reconcile signal built
  yesterday (sev 9 "never displaces" by design). A budget needs an overflow
  valve.
- **1 out-of-pool worktree + 1 stale `.claude/worktrees` agent worktree** (at
  dde7def, 5 commits behind) — `bench worktree clean` covers these; same
  maintenance pass.
- **Gate green, 29.6s wall ✓** — the oracle is not a friction point.

---

## Cross-cutting themes

1. **Close the loop *after* merge.** The 2026-07-06 assessment closed the
   forward artifact loop (reviews/, retirement tooling, reconcile signals);
   what's missing now is the *actor*: no phase reached after a merge consumes
   the duty list the platform already computes. The pattern "signal exists,
   nobody's charged with acting on it" repeats across specs, roadmap,
   decisions/, HANDOFF.md, and scratch branches.
2. **Dogfooding masks adoption defects.** Every §3 med is invisible in the kit
   repo because the kit repo never runs `bench init` cold, never installs from
   npm, and always has a global `bench`. A fresh-install smoke path (a gate
   fragment that links+inits a throwaway repo *from the packed tarball, PATH
   stripped*) would have caught all three.
3. **Fix defect classes, not instances.** False-empty-on-git-failure was fixed
   in `diff` and re-found in `structure`; arg-parsing rigor was standardized
   except `models`. When a drain fixes a pattern, the spec should include the
   repo-wide audit for the pattern's other instances.
4. **Deliberate postures deserve records.** Fail-open rims, recompute-always
   gating, family-level canary granularity, honor-system interactive
   invariants — each is defensible, none is citable. The disposition table of
   a drained (now-replaced) assessment is not where accepted risk should live.

## Ranked improvement backlog

Ordered by platform leverage; sizes are rough.

1. **Run the pending maintenance pass** — `/bench-what-next` to drain the two
   learnings (each with its proposed rule change) and reconcile the roadmap,
   then `bench spec retire` for the ten merged specs, then `bench worktree
   clean`. Zero build cost; clears four status rows and makes the tree honest
   before any new work. (XS, operational)
2. **Charge an owner with the post-merge tail** — extend `/bench-final-check`
   with an exit duty ("on the default branch, check `bench status` and run the
   retirement/reconcile rows it flags") or add a thin after-merge phase that
   consumes status rows as a queue. Include: decision-map sweep policy (decide
   keep-vs-retire for closed maps), and the HANDOFF.md lifecycle (decide:
   generate at session end, or delete the file and let status+roadmap be the
   pickup). Closes the §2 high and three meds. (M)
3. **First-install fixes** — ship `projects/` templates in `files[]` (or embed
   them in `bench init`); stop shipping `HANDOFF.md`; write a `.bench/dist`
   gitignore entry (or arch-check with a loud re-link error) at link time; make
   the scaffolded gate resolve the `.bench/bin` local CLI before assuming a
   global `bench`. Add the fresh-install smoke fragment (theme 2). (S–M)
4. **Close or record the line-guard fail-open** — deny (or warn-to-user
   loudly) when a delegation envelope has no model field, or record fail-open
   as a decision; and add pre-push presence/managed-bytes verification to
   `bench doctor` and a status signal for a missing hook. (S)
5. **`bench worktree clean` sweeps orphaned scratch branches** — the agent
   cannot perform the currently-recommended remedy at all (guard-blocked), so
   the CLI must own it. (S)
6. **Kill the false-empty class repo-wide** — `bench structure` errors loudly
   on git failure; audit every discarded `git.Output` error in porcelain;
   `bench models` rejects unknown args at exit 2. (S)
7. **Status overflow valve + invocable actions** — a flag (or `+N more`
   expansion) to show rows past the five-row budget, and action strings phrased
   as invocable commands/phases. (S)
8. **Record accepted enforcement postures as one decision record** —
   interactive commit-on-red, non-shift done-claims, declare-line
   unenforceability, fail-open rims, recompute-always gating, family-level
   canary granularity. One ADR; makes every future "is this a hole?" question
   answerable by citation. (S)
9. **Unpark FT23 (`craft-spec` skill)** — this assessment independently
   re-confirmed the coverage-map schema is unreachable from the two
   model-invoked skills that need it; second consecutive assessment to flag it.
   (S–M)
10. **Structure splits where genuine** — `axi_wave2_test.go` by command;
    `line_routing_checks_test.go` static-vs-exec; accept the two noise files
    (or set per-file budget notes) so the structure signal stays credible. (S)
11. **Uninstall story** — `bench unlink` consuming the manifest, or a
    documented manual path. Reviewer call on priority; absence is currently
    undocumented. (M)
12. **Cheap hardening batch** — SafeToken newline corpus row; a `stophook.Run`
    seam test; README prerequisites (Go-or-node, Windows/WSL note, nvm
    caveat); README `.bench/` layout refresh; trim the redundant third trigger
    clause on the three longest skill descriptions. (S, mostly mechanical)

Parked/known, re-confirmed still open: FT22 (`bench spec history`), FT24
(Codex line guard), CLI-inventory demotion (reviewer tradeoff), `craft-assess`
owner if the assessment cadence repeats.

## Verification notes

- Synthesizer-verified (✓ above): all ten drain verdict spot-checks sampled
  (reviews/ empty, retire tooling, roadmap signal); HANDOFF.md staleness
  against git history and live tree; retirement prose absence in
  final-check; `files[]` contents; doctor's zero hook coverage;
  `AgentLineVerdict` fail-open rim; `structure.go` discarded errors; gate
  cache write-only; orphaned-branch remedy dead end; the hidden sixth status
  row; live gate green at 29.6s.
- Remaining findings are delegate-cited with file:line and sampled, not
  exhaustively re-checked.
- Known coverage limits (unknowns): `internal/gate` beyond phases/pin and
  `internal/contract` bodies read as outlines only; `git.DefaultBranch`
  behavior in a clone with no `origin/HEAD` (a wrong guess would point
  pre-push at the wrong branch — worth a test); native-Windows and
  offline-registry install behavior not executed; pre-push firing from inside
  a linked worktree reasoned, not exercised; gitguard evasion surface beyond
  push-denial not probed (out of its stated honest-mistake threat model);
  hook runtime behavior read statically, not executed; the two shipped
  example profiles (gl-axi, regroup) assessed for shape only.
