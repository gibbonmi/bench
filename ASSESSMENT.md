# Bench platform assessment — 2026-07-06

Deep assessment of the skills, commands, enforcement, docs, CLI/Go core, and
artifact lifecycle, hunting inconsistencies and improvement leverage. Produced by
six read-only area sweeps (mid tier) synthesized and spot-verified on the top
tier; claims marked ✓ were independently re-verified against source by the
synthesizer. Findings from the FT9 code review live in `reviews/ft9.md` and are
referenced, not duplicated, here.

Severity: **high** = an invariant or advertised guarantee is not actually held;
**med** = a real defect or reachable unowned state; **low** = friction, drift
risk, or hygiene.

---

## Executive summary

The platform is in good shape where it is densely gated: docs anchors, mirror
parity, the AXI contract suite, and the shift/push enforcement path all held up
under adversarial reading. The problems cluster in four themes:

1. **Lifecycle states with no owner.** Artifacts can reach states no phase or
   command handles: a merged spec awaiting retirement, a `reviews/<slug>.md`
   whose fix pass never runs, a roadmap row for work that already shipped, a
   review requested after the work merged. Each was found live in the tree
   today, not hypothetically.
2. **Advertisement stronger than enforcement.** The kit's own defect class —
   "an enforcement and its advertisement must collapse to one source" — appears
   in its enforcement layer: `bench guards` advertises hooks by file presence
   not wiring, pre-push advertises pin-drift denial it skips when unpinned, and
   the primary harness (Claude) has *less* conformance-gated hook wiring than
   the secondary (Codex).
3. **One-source-per-fact violations in the kit that preaches it.** Five live
   instances: the not-in-repo message (3 phrasings), the coverage-map schema
   (command + skill), the tier binding (env + prose), the `bin/bench.sh`
   header comment roster (stale third copy), and the new `--full` contract
   facts (help const + package doc).
4. **Agents still hand-run what the CLI should own.** The FT9 pattern
   (fold a repeated hand-run git call into the CLI) has more instances waiting:
   retired-spec recovery, spec status listing, and the `bench commit` docs
   that omit its mandatory `-m`.

---

## Findings by area

### 1. Workflow commands (`.agents/commands/`)

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| high | gap | `reviews/<slug>.md` lifecycle is one-sided: review-implementation promises "the `/bench-implement-spec` session that resolves the findings deletes `reviews/<spec-slug>.md`", but `bench-implement-spec.md` never mentions reading, resolving, or deleting it. ✓ (`rg reviews/` hits only the review command) | `bench-review-implementation.md:89-92` vs `bench-implement-spec.md` |
| high | gap | Empty-diff state unhandled: step 1 says "Confirm the diff is non-empty" with no remedy, and the documented happy path (implement-spec commits → recommends review) reaches it every time work lands on the default branch (merge-base == HEAD). ✓ hit live this session reviewing FT9 | `bench-review-implementation.md:30-36`; `internal/diff/diff.go:105-117` |
| med | gap | `reviews/<slug>.md` orphans when the reviewer accepts residual risk and routes to final-check: no phase or reconcile deletes it, so it lingers as exactly the false pickup work the command warns against. | `bench-review-implementation.md:17-19,88`; `bench-final-check.md` (no mention) |
| med | gap | Shipped-spec retirement has no standalone owner: promote-then-delete lives only inside the *next* `/bench-write-spec`; what-next reconcile removes ROADMAP rows, not `specs/` files; there is no `bench spec retire`. | `bench-write-spec.md:167-168`; `bench-what-next.md:26-27` |
| med | inconsistency | Both `bench commit` call sites omit the mandatory `-m <msg>`; an agent following either doc verbatim exits 2. ✓ Bonus: `commit.go:1`'s own synopsis writes `[-m <msg>]` as optional while `:25/:32` require it. | `bench-implement-spec.md:85`, `bench-final-check.md:16-18` vs `internal/commit/commit.go:32` |
| low | improvement | Retired-spec recovery git incantation duplicated across two commands — an FT9-pattern CLI-folding candidate (`bench spec history`). | `bench-debug.md:87-88`; `bench-write-spec.md:165` |

Clean: `.claude/commands` byte-identical (symlink); all referenced subcommands
and flags exist; phase skills are thin adapters with no drift surface.

### 2. Enforcement layer (`.bench/`, hooks, conformance)

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| high | gap | Claude `settings.json` wiring for Stop, PreToolUse:Bash, and SessionStart is not conformance-gated (only the Agent matcher is checked); deleting the Stop or git-guard wiring on the primary harness leaves the gate green. ✓ | `internal/conformance/line_routing_checks_test.go:89-117`; no Claude Stop/Bash check exists |
| high | inconsistency | The Codex equivalents *are* gated (`Stop→stop.sh`, `Bash→block-dangerous-git.sh`), so the secondary harness is better-protected than the primary. ✓ | `internal/conformance/validity_checks_test.go:130-151` |
| med | inconsistency | `bench guards` discovers hooks by file presence, not wiring — on Codex it advertises `check-agent-line` and `session-start` that aren't wired there, overstating the deny surface. Pre-push, by contrast, gets a real install+marker check. | `internal/guards/guards.go:106-123` vs `.codex/hooks.json`; `guards.go:128-149` |
| med | gap | Pre-push advertises "`.bench` drift from bench gate pin" unconditionally but skips drift enforcement (warn-only) when no pin file exists. | `.git/hooks/pre-push` (`if [[ -n "$pin_tree" … ]]`) |
| med | gap | Invariant 2's interactive line guard (`check-agent-line`) is Claude-only; interactive delegation onto an unbound model on Codex/opencode is unguarded outside the shift loop. | `.codex/hooks.json` (no Agent matcher); `.bench/BENCH.md` "extra layer where supported" |
| med | improvement | shellcheck phase omits the adapters (`.bench/adapters/*` — load-bearing enforcement), the gate fragment scripts, and pre-push. | `internal/gate/phases.go:100-122` |
| low | improvement | shellcheck is Optional → the whole shell-lint defense silently "skips" where the binary is absent; EACCES is also treated as skip, masking a present-but-unexecutable binary. | `internal/gate/phases.go:61-64,227-230,264-269` |
| low | gap | Invariant 4's "commit on green, never on red" is enforced only inside `bench shift`; an interactive `git commit` on a red gate is not blocked (guard blocks `--amend`, not plain commit). | `internal/shift/verdict.go:50-52` |

Invariant → mechanism map (from the sweep): 1 solid for shift/push, honor-system
interactively; 2 solid for shift, weak for interactive delegation off-Claude;
3 solid (densely gated); 4 solid for shift/push, weak for interactive
commit-on-red.

### 3. Artifact lifecycle and workflow state

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| high | accuracy | `ROADMAP.md` FT9 row is stale: still "Grilled and mapped… ready" with recommended-sequence step 2 saying run `/bench-write-spec`, though FT9 is implemented and merged (`51f5075`, `specs/ft9.md:3`). Violates the roadmap header's own "a row leaves when the work ships." ✓ | `ROADMAP.md:13,59` |
| med-high | gap | Review-after-merge has no lifecycle representation: `bench diff --full` from the default branch is empty, yet the review phase both sources its diff there and demands non-empty. No documented fallback (e.g. first-parent diff of the merge/landing commit). ✓ hit live | `bench-review-implementation.md:35`; confirmed empty output this session |
| med | gap | `reviews/<slug>.md` orphans if its spec retires first — spec-retire prose never mentions `reviews/`, and the file escapes git-based rot sweeps because it is untracked. | `bench-write-spec.md:155-168`; `git ls-files reviews/` empty |
| med | gap | The review artifact is born untracked and no phase commits it — invisible to git-based checks until a fix session happens to notice; the command assumes a fix commit deletes it but nothing owns getting it into git first. | `bench-review-implementation.md:89-92` |
| med | ergonomics | The transient review artifact flips the gate stale: `bench tree-hash` includes untracked files, so an advisory review dirties the gate verdict. Caveat worth a look: status's "gated tree a6c2c9f" could not be reconciled with the pin file content at `.git/bench-last-gate`. | `internal/status/status.go:121` |
| low | gap | spec-retire is prose-only — status flags the pending count but a human must remember the manual promote-then-delete. | `rg spec-retire` hits prose only |
| low | ergonomics | Out-of-pool/agent worktrees linger until someone runs `bench worktree clean` (status does surface it). | `internal/status/status.go:209` |

Clean: `IDEAS.md` empty, `.bench/learnings.md` has zero open entries — the
capture drains are current.

### 4. CLI and Go core (`bin/bench.sh`, `internal/`)

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | inconsistency | Not-in-repo message derived in 3+ forms: `toon.NotInRepo()` declares itself "one source", `commit.go:38` hand-duplicates its rendered literal, ~8 operational commands emit a bare `not in a git repo`, `gate/phases.go:82` a third form; `git.go:72` claims a uniform posture that doesn't exist. ✓ | `internal/toon/toon.go:111-115`; `internal/commit/commit.go:38`; `internal/gate/gate.go:253` et al. |
| med | inconsistency | `bench roadmap` outside any repo exits 0 with "no ROADMAP.md" (conflating not-in-repo with no-roadmap) while sibling `bench idea` in the same package returns NotInRepo exit 1. ✓ | `internal/roadmap/roadmap.go:83-89` vs `:31-33` |
| low-med | quality | `bench diff` turns git failures into false-empty success in three places (`changedFiles`, `commitLog`, the `--full` body). Overlaps the low finding in `reviews/ft9.md`; blast radius narrow since base resolution is guarded loudly. | `internal/diff/diff.go:45-48,74-79,140` |
| low | inconsistency | `bench coverage --check` success returns `("", 0)` — violates the AXI "definitive empty states" rule; agent can't distinguish pass from no-output. | `internal/coverage/coverage.go:232-234` vs `bench-craft-cli/SKILL.md:53-54` |
| low | inconsistency | coverage `--check` violation lines use a shape divergent from the canonical `error: <kind> — <hint>` that `toon.Errorf` is meant to solely own. | `internal/coverage/coverage.go:237` |
| low | inconsistency | `bin/bench.sh` header comment lists a stale subcommand subset — a third roster of a fact already owned by the usage heredoc and BENCH.md. | `bin/bench.sh:10-17` vs `:242-266` |
| low | quality | `bench maps --count` and `bench guards --brief` are real flags documented only in code (internal adapter hooks; acceptable but undiscoverable). | `internal/maps/maps.go:212`; `internal/guards/guards.go:162` |
| low | quality | `toon.Usage` misattributes the offending argument when a valid flag precedes an unknown one: `bench diff --full bogus` reports `unknown argument: --full`. ✓ (synthesizer finding, from the FT9 diff) | `internal/diff/diff.go` arg switch default (`args[0]`) |

Shell/Go boundary judged clean — `bin/bench.sh` is binary-resolution and routing
that must run before the binary exists. Not inspected in depth by the sweep:
`adopt`, `canary`, `conformance` (partial), `contract` (partial), `gate`
(partial), `gitguard`, `lines`, `models`, `packagesurface`, `stophook`,
`structure`, `subprocess`, `terminal`.

### 5. Skills (`.agents/skills/`)

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | inconsistency | Acceptance-coverage-map schema derived twice — the 5 fields and the red-signal gloss appear in both the write-spec command and `craft-tdd`. The kit's cardinal rule, broken in its two most-used build/review surfaces. | `bench-write-spec.md:102-104` vs `bench-craft-tdd/SKILL.md:71-84` |
| low-med | improvement | `craft-cli`'s trigger fires only for a CLI "whose project declares AXI conformance" — but deciding *whether* a new CLI conforms is exactly when it's needed; the skill governing the decision won't load at the decision point. | `bench-craft-cli/SKILL.md:3` |
| low | inconsistency | `craft-review` and `craft-tdd` call the edge-inventory *step* "the canonical list" of edge classes; the real source is the profile checklist / `hostile-input-library.md`. Points the weakest reader at a procedure, not the source. | `bench-craft-review/SKILL.md:38`; `bench-craft-tdd/SKILL.md:28` |
| low | improvement | No model-invocable skill owns spec-authoring/coverage-map discipline (it lives only in the user-invoked command) — the structural root of the med finding above. | `bench-write-spec.md:229-232` |
| low | inconsistency | `craft-seams` restates the over-fit rationale `craft-tdd` owns — borderline, cross-referenced, mild one-source smell. | `bench-craft-seams/SKILL.md:46-49` vs `bench-craft-tdd/SKILL.md:8-9` |

Clean: phase skills are thin adapters; `.claude/skills` is 12 symlinks; skills
index auto-generated from frontmatter; all referenced subcommands/flags/paths
exist; the `craft-delegate`↔`craft-line` split is a model one-source example.

### 6. Platform documents

| Sev | Kind | Finding | Citation |
|---|---|---|---|
| med | drift | BENCH-reference's Hook Layers omits the shipped pre-push gate-pin verification — lists only "blocks direct pushes to the default branch." | `.bench/BENCH-reference.md:119` vs `internal/gate/pin.go`, `bin/bench.sh:259` |
| med | drift | ADR 0001's title frames tripwire "**rather than** write protection," but push-time pin protection has since shipped — current decided state is tripwire *and* pin. Violates invariant 3 (ADR = current state). | `docs/adr/0001-working-tree-gate-tripwire.md:1` |
| med | duplication | Tier→model binding derived twice — `.bench/lines.env` (machine) and benchkit.md "Lines" prose — coupled only by a manual keep-in-sync note; no gate check couples them. | `projects/benchkit.md:158-169` vs `.bench/lines.env:5-7` |
| low | duplication | benchkit.md seam list re-enumerates the subcommand set, partially (omits `doctor`, `canary`, `version`, `gate pin`) — a second copy that can drift. | `projects/benchkit.md:24-31` |
| low | improvement | ADR 0001 embeds a file-path token, which invariant 3 forbids in ADRs. | `docs/adr/0001…:3,:11` |
| low | improvement | The 9 hook/adapter plumbing subcommands sit in always-loaded BENCH.md (imported every session) though sessions almost never invoke them — pure token cost; BENCH-reference already owns plumbing detail. | `.bench/BENCH.md` CLI Inventory |
| low | improvement | CONTEXT.md defines "skill" by its `.claude/` mirror path rather than the canonical `.agents/skills/` source. | `CONTEXT.md:32` |

Clean: BENCH.md CLI Inventory matches `bin/bench.sh` dispatch exactly (all 30
subcommands); skills-index and Codex adapter lists match the tree 1:1; AGENTS.md
does not restate the invariants and the conformance test enforces that.

---

## Cross-cutting themes

1. **Close the artifact loop.** The pipeline's forward path (idea → map → spec →
   build → review → gate → commit) is well-owned; the *backward/cleanup* path is
   not. Spec retirement, review-artifact resolution, and roadmap reconciliation
   after a merge all depend on someone remembering. Every high/med lifecycle
   finding above is an instance of this one gap.
2. **Enforcement should be measured by wiring, not files.** `bench guards`,
   the pre-push advertisement, and the Claude-vs-Codex conformance asymmetry all
   stem from checking that enforcement *exists* rather than that it is *wired*.
   The pre-push install check already shows the right pattern.
3. **The kit should pass its own one-source audit.** Five concrete duplications
   found; most are cheap collapses. A conformance-style check for the known
   rosters (subcommand lists, tier binding) would keep them from regrowing.
4. **Keep folding hand-run git into the CLI** (the FT9 pattern earned its keep —
   the review phase now needs one command where it needed three; extend to spec
   history/status surfaces).

## Ranked improvement backlog

Ordered by platform leverage; sizes are rough.

1. **Own the artifact lifecycle end-to-end** — make spec-retire a real operation
   (`bench spec retire <slug>` or an explicit what-next duty) that also sweeps
   `reviews/<slug>.md`; teach `bench-implement-spec.md` to read/delete the
   review pickup file; add an orphan check to `bench status`. Closes 5 findings
   across two areas. (M)
2. **Conformance-gate the Claude hook wiring** (Stop, Bash, SessionStart),
   mirroring the existing Codex validity checks — the largest silent
   de-enforcement hole, on the primary harness. (S)
3. **Give review a review-after-merge mode** — documented fallback in
   review-implementation step 1 (first-parent diff of the landing commit) when
   `bench diff --full` is empty; consider `bench diff --commit <sha>`. Hit live
   this session. (S–M)
4. **Reconcile the roadmap when work ships** — a status row ("roadmap row for
   merged work") or a what-next duty triggered by spec `Status: implemented`,
   so the drain cadence can't lag reality. Fix the stale FT9 row now. (S)
5. **Report actual wiring in `bench guards`** and make the pre-push
   advertisement conditional on a pin existing — collapse advertisement to
   enforcement. (M)
6. **Collapse the five duplicated-knowledge instances** — not-in-repo message
   (route everything through `toon.NotInRepo` or a documented two-form split),
   coverage-map schema (one owner, likely a `craft-spec` skill or the command),
   tier binding (derive prose from `lines.env` or gate the pair),
   `bin/bench.sh` header roster (delete it), `--full` help/doc pair (reviewer
   call, see `reviews/ft9.md`). (M, mostly mechanical)
7. **Fix `bench commit` doc/arg mismatch** — show `-m` in both phase docs and
   fix `commit.go:1`'s synopsis. (XS)
8. **Definitive success outputs** — `bench coverage --check` emits a pass line;
   `bench roadmap` distinguishes not-in-repo (exit 1) from no-roadmap (exit 0);
   coverage violations use the canonical error shape. (S)
9. **Extend shellcheck to adapters, gate fragments, and pre-push**, and decide
   whether Optional is acceptable for enforcement code (at minimum, distinguish
   EACCES from absent). (S)
10. **Docs drift pass** — rewrite ADR 0001 to the current tripwire+pin state,
    add the pin layer to Hook Layers, fix CONTEXT.md's skill definition, demote
    plumbing subcommands to BENCH-reference. (S)
11. **Decide the `reviews/` tree-hash question** — either exclude transient
    pickup state from `bench tree-hash` or require the review phase to commit
    the artifact; today an advisory review flips the gate stale. Also
    investigate the unreconciled "gated tree a6c2c9f" vs pin-file caveat. (S,
    but a reviewer decision first)
12. **Tighten `craft-cli`'s trigger** to fire when *deciding* whether a new
    agent-facing CLI conforms, and repoint the "canonical list" edge-inventory
    pointers at the true source. (XS)

## Verification notes

- Synthesizer-verified (✓ above): reviews/-lifecycle one-sidedness, `bench
  commit -m` requirement, not-in-repo triplication, roadmap not-in-repo
  posture, Claude/Codex conformance asymmetry, stale ROADMAP FT9 row, empty
  `bench diff --full` on merged work (hit live), `toon.Usage` misattribution,
  `toon.Table` control-byte refusal (in `reviews/ft9.md`).
- Remaining findings are delegate-cited with file:line and sampled, not
  exhaustively re-checked.
- Known coverage limits: the CLI sweep's not-inspected package list (§4); the
  enforcement sweep read hooks statically and did not execute them; no
  assessment of `projects/gl-axi.md` / `projects/regroup.md` beyond shape.
