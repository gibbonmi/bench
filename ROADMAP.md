# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

## Features, in priority order

Rows FT25–FT36 are the 2026-07-07 assessment drain (`ASSESSMENT.md` carries
the finding-level detail). Each gets a decision map and a spec before build;
rows gain their spec paths as the specs are staged.

**FT25 (HIGH) — charge an owner with the post-merge tail.** A
`/bench-final-check` exit duty (or thin after-merge pass) that consumes the
`bench status` duty rows: spec retirement, roadmap row removal, a decision-map
sweep policy, and the HANDOFF.md lifecycle decision.

**FT26 (MED) — first-install fixes.** Ship the `projects/` profile templates
in the npm package; stop shipping the kit's own HANDOFF.md; ignore-or-arch-check
`.bench/dist` at link time; make the scaffolded gate resolve the `.bench/bin`
local CLI before assuming a global `bench`; add fresh-install smoke coverage.

**FT27 (MED) — enforcement verification.** Agent-line guard stops failing open
when the delegation envelope omits the model field (deny, or record the
posture); `bench doctor` verifies the pre-push hook is present and
bench-managed, with a status signal when it is missing.

**FT28 (MED) — `bench worktree clean` sweeps orphaned scratch branches.** The
currently-recommended remedy is guard-blocked for the agent, so the CLI must
own deletion of fully-merged `worktree-agent-*` branches.

**FT29 (MED) — kill the false-empty class.** `bench structure` errors loudly
on git failure; audit every discarded `git.Output` error in porcelain;
`bench models` rejects unknown args at exit 2.

**FT30 (MED-LOW) — status overflow valve + invocable actions.** A way to see
rows past the five-row budget, and action strings phrased as commands or
phases the reading harness can invoke.

**FT33 (MED-LOW) — `bench unlink`.** Consume the link manifest to remove the
repo footprint; document the uninstall story.

**FT23 (LOW, unparked) — model-invocable spec-authoring skill.** A
`craft-spec` skill owning the coverage-map schema; `/bench-write-spec`
composes it. Second consecutive assessment flagged the reachability gap.

**FT31 (LOW) — enforcement-postures ADR.** One decision record for the
accepted honor-system residuals: interactive commit-on-red, non-shift
done-claims, declare-line unenforceability, fail-open rims, recompute-always
gating, family-level canary granularity.

**FT32 (LOW) — structure splits where genuine.** `axi_wave2_test.go` by
command; `line_routing_checks_test.go` static-vs-exec; record the two noise
files as accepted so the structure signal stays credible.

**FT34 (LOW) — cheap hardening batch.** SafeToken newline corpus row;
`stophook.Run` seam test; README prerequisites (Go-or-node, Windows/WSL,
nvm caveat) and `.bench/` layout refresh; trim the redundant third trigger
clause on the three longest skill descriptions.

**FT35 (LOW, rule) — craft-delegate stale-base opener.** From the 2026-07-06
learning: every worktree-isolated delegate charge opens with `git merge
--ff-only main`, verifies HEAD, and stops if denied. Kit edit under
`craft-synthesis`.

**FT36 (LOW, rule) — write-spec batch-drain override clause.** From the
2026-07-06 learning: an explicit reviewer-directed batch drain may substitute
for per-spec maps, with every defaulted decision flagged for veto; absent that
instruction the map gate stands. Kit edit under `craft-synthesis`.

**FT12 (LOW, kit discipline) — repro a defect claim through the accused command
before draining it.** FT11 was minted from a learning that quoted a raw `git add`
run by hand; the real `bench commit` path already staged deletions, so the row
described a defect that did not exist. Tighten `/bench-what-next` step 3 (and
`bench-debug`'s repro discipline) so a defect-shaped learning becomes a roadmap
row only after its red signal reproduces through the sanctioned command, not a
lookalike. Built later under the `craft-synthesis` discipline.

**FT10 (LOW) — doctor installs the kit repo's pre-push guard.** `bench guards`
already reports the missing guard; `bench doctor` should detect it on the kit
repo itself and offer the install (consumer repos get it via `bench link`).
Overlaps FT27's doctor verification — build together.

**FT5 (LOW) — `bench outline`.** Marginal for this repo, real as a kit
affordance for large/polyglot linked repos. Needs its grill (languages,
on-demand vs committed, prose anchors).

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check, and a per-anchor
bite-proof meta-test (canaries prove one needle per family today; graduate on
observed anchor rot). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search.

**FT7 (LOW) — dashboard.** Low priority by declaration.

**FT22 (LOW, parked) — `bench spec history <slug>`.** Fold the duplicated
`git log --grep=spec-retire` recovery incantation into the CLI (FT9 pattern).
Parked from the artifact-lifecycle build's out-of-scope list.

**FT24 (LOW, parked) — Codex agent-line guard parity.** `check-agent-line` on
the secondary harness, pending research on whether Codex hooks support an
Agent matcher. Parked from the claude-hook-conformance build.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

## Watch

- `bench worktree concurrent-acquire` contract test failed once under
  full-gate load, then passed 3/3 in isolation and on rerun — likely a timing
  flake surfaced by gate phase concurrency. Journal it if it recurs.

## Recommended sequence

1. FT25 after-merge owner — `/bench-shape-idea` then `/bench-write-spec`
2. FT26 first-install fixes — `/bench-shape-idea` then `/bench-write-spec`
3. FT27 enforcement verification — `/bench-shape-idea` then `/bench-write-spec`
