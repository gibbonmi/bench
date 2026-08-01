# Pre-push guard visibility — research findings

Six read-only research delegations, 2026-07-31. Every claim below was
re-verified against the tree by the coordinator before landing here.

## R1 — how the hook picks its protected branch

Two independent resolutions, install-time and runtime.

**Install time.** `protectedBranch` (`internal/adopt/link_hook.go:101-106`)
calls `git.ResolvedDefault` and falls back to `fallbackProtectedBranch = "main"`
(`link_hook.go:98`). `git.ResolvedDefault`
(`internal/git/default_branch.go:16-29`) returns a branch only when one resolves
to a commit: `origin/HEAD`'s short name, else the `main` candidate, else a sole
local branch, else `("", false)`. The baked-guess case is therefore narrow — no
resolvable `origin/HEAD`, no commit named `main`, and either zero or two-plus
local branches — but real. The resolved token is substituted for
`__BENCH_DEFAULT_BRANCH__` (`link_hook.go:24`) at install.

**Runtime.** The installed hook re-resolves on every push
(`internal/adopt/prepush.sh:17-21`): it seeds `protected` from the baked token,
then overrides it from `git symbolic-ref --short refs/remotes/origin/HEAD` when
that is non-empty. When the live lookup is empty the guard silently keeps the
baked value — no warning, no exit, no degraded posture. The only stderr in the
script is the gate-pin notice (`prepush.sh:22-24`) and the two block messages
(`prepush.sh:31,35,43`).

**Nothing records which path was taken.** Not in the manifest header
(`prepush.sh:2-6` — `boundary: pre-push` is a literal; the branch appears only
in prose at line 7), not in `ClassifyPrePush` (whose result carries `State` and
`Path` and no branch field, `link_hook.go:44-49`), not in `prePushRow`
(`internal/guards/guards.go:241-262`). The branch surfaces at runtime only in
the block message itself (`prepush.sh:43`), which names it without saying which
path produced it.

The fallback's deliberateness is recorded in
`docs/adr/0010-absence-is-the-only-authoritative-empty-state.md` ("the pre-push
hook is the one deliberate exception").

## R2 — what `bench doctor` checks on hooks

Exactly one hook check: the pre-push backstop row, called from
`internal/adopt/doctor.go:230-232`, rendered by `reportPrePush`
(`doctor.go:247-265`), classified by `ClassifyPrePush`
(`link_hook.go:56-69`). No doctor check of any other git hook and none of the
harness hooks. The row registry `doctorRows`
(`internal/adopt/doctor_rows.go:26-35`) carries no other hook row.

The classification inspects: git-worktree presence (`doctor.go:248-251`), the
hooks dir git will use (`hooksDir`, `link_hook.go:79-88`, honoring
`core.hooksPath`), file readability (`link_hook.go:58`), the **marker substring**
`bench:managed-pre-push` (`link_hook.go:32,59`), and `core.hooksPath` diversion
(`hooksPathConfigured`, `link_hook.go:74-77`). Nothing else — not the executable
bit, not the baked branch, not the body beyond the marker.

Substring matching is deliberate, documented at `link_hook.go:26-31`: byte
identity "would false-red across default-branch token substitution and benign
template evolution."

Doctor never repairs the hook (`doctor.go:243-246`: self-heal rejected for least
surprise; the remedy named is `bench link`). The same classifier is reused by
`bench status` (`internal/status/status.go:409-424`).

## R3 — how `bench guards` reads manifest state

`enumerateGuards` (`internal/guards/guards.go:57-75`) walks `.bench/hooks/*.sh`
and appends one synthetic pre-push candidate (line 73). `HeaderFields`
(`guards.go:104-117`) hands content to `parseHeader` (`guards.go:119-133`), which
reads only the **leading comment block** — it breaks at the first non-blank,
non-`#` line (line 122) — and takes the first `# <key>: ` occurrence of each
required key. Required keys: `name`, `boundary`, `denies`, `why`
(`guards.go:90`); an empty value counts as missing (`guards.go:92-100`).

The boundary cell is purely the parsed value, with no fallback or inference: an
incomplete header yields `""` (`guards.go:143,147`).

`prePushRow` (`guards.go:241-262`): absent → `not installed`; present without the
marker → `unmanaged (no manifest)` (`guards.go:255-256`); marker present → falls
through to `guardRow`, which returns `no manifest` under the fallback name when
the header is unreadable or incomplete (`guards.go:138-153`).

Consumed by `Command` (`guards.go:271-303`) as TOON plus a `guard_scan` meta row,
and by SessionStart through `--brief`
(`internal/sessioninspect/sessioninspect.go:84`, phase list line 20), wired at
`.claude/settings.json:22` and `.codex/hooks.json:8`.

**No gate or conformance check enforces manifest-header presence on an installed
hook.** `bench status`'s guards row uses `ClassifyPrePush` instead and emits no
row at all when the hook is managed (`internal/status/status.go:409-413`), so it
cannot surface this condition.

## R4 — `bench link` order of operations (settles face 3)

Link runs as wrapper preflight → arg/repo/kit checks → plan build → one staged
transaction whose only write step is `promoteAll`.

Aborts, in order: kit-asset tree present (`bin/bench.sh:256-263`); mode valid
(`internal/adopt/link.go:32-35`); `git.Root()` (`link.go:36-44`);
`buildLinkPlan` (`link.go:50-54`, refusal at `link.go:184-190`); manifest read
(`internal/adopt/link_transaction.go:20-24`); AGENTS.md invalid managed block
(`link_transaction.go:32-39`); **foreign pre-push conflict**
(`link_transaction.go:44-48`); then the per-entry plan loop with aborts for
missing kit asset (`:67-72`), **symlink parent** (`:57-60`, `:73-76`),
non-directory parent (`:77-81`); staging errors (`:92-118`); dropped-row
reconciliation (`:120-163`).

The hook is staged at `link_transaction.go:186-199` and lands only at
`promoteAll` (`:222`, `internal/adopt/transaction.go:78-103`).

**Face 3 confirmed.** The symlink-parent abort at `link_transaction.go:73-76`
returns `1, false` roughly 110 lines before the hook is staged, and nothing
reaches disk without `promoteAll`. Verified in this repo: `.claude/commands` is a
symlink to `../.agents/commands`, and `.claude/commands/bench-debug.md` resolves
through it. So `bench link` here genuinely cannot reach the hook refresh — the
sanctioned repair for face 2 is unavailable exactly where the stale hook lives.

**Coordinator extension (2026-07-31).** `bench upgrade` shares this path:
`Upgrade` builds the same plan and calls `transactionalLink`
(`internal/adopt/upgrade.go:50,73`), so every abort above — including the
symlink-parent one — blocks `bench upgrade` as well. Face 3's blast radius is
therefore both adoption routes, not just `link`. Separately, upgrade's reported
`changed` count comes from `upgradePlanCounts(manifest, plan)`
(`upgrade.go:54`), which walks plan entries only; the hook is staged outside the
plan, so `bench upgrade --check` can report a plan whose counts never mention a
hook refresh it will in fact perform — another instance of the R6 class.

**Unplanned finding.** `installGitHook` (`link_hook.go:108-126`) has no call site
anywhere in the tree; link's only hook path is the staged one. The
token-substitution knowledge exists twice (`link_hook.go:124` and
`link_transaction.go:187`), one copy dead.

## R5 — template storage and currency-comparison precedent

`internal/adopt/prepush.sh` is embedded at `link_hook.go:19-20` via `//go:embed`,
shipped as lintable shell (shellcheck phase registers it at
`internal/gate/phases.go:248`, asserted at
`internal/conformance/package_core_checks_test.go:278`).

**Nothing compares an installed hook against it.** `prePushTemplate` is read at
exactly two sites, both writes (`link_hook.go:124`, `link_transaction.go:187`).
Every consumer of installed-hook identity matches the marker substring instead:
`link_hook.go:59`, `link_hook.go:120`, `link_transaction.go:45`,
`internal/adopt/unlink.go:269`, `internal/adopt/doctor.go:252-262`,
`internal/status/status.go:409-425`, `internal/guards/guards.go:254-256`.
`PrePushState` (`link_hook.go:37-42`) has four values — managed, absent, foreign,
diverted — and no drifted or stale state.

**Precedents that do compare against a declared source.**

- Manifest SHA-256 for linked kit files: `fingerprintPath`
  (`internal/adopt/manifest.go:101-131`), consumed by `manifestOwnedClean`
  (`link.go:252-262`), `upgradePlanCounts` (`internal/adopt/upgrade.go:116-138`),
  `link_transaction.go:110,124,151,181`, `internal/adopt/transaction.go:62-65`,
  `unlink.go:103`. The hook is deliberately excluded: "The hook is bespoke, not a
  manifest row" (`unlink.go:265`).
- Release package inventory vs declared manifest — size, mode, digest
  (`internal/releaseevidence/package_artifact.go:190`).
- Source-tree fingerprint before promotion
  (`internal/releaseevidence/evidence_fingerprint.go:31-79`).
- Embed-target existence, not content
  (`internal/packagesurface/assets.go:106-109`).

## R6 — the class: enforcement controls reporting their own status

| Site | Asserts | Actually checks |
|---|---|---|
| `doctor.go:216-228` shim row | "healthy bench shim" | marker substring `bench-shim v1` (`doctor.go:196`, `:161-164`) plus `# bench-target:` points at an executable (`:206-211`). Body comparison against `ShimContent` exists only in `doctorFix` (`:305-308`), so a hand-edited marked shim reports healthy. |
| `doctor.go:247-265` pre-push row | "bench-managed pre-push at <path>" | marker substring only. Neither content currency nor which branch is protected. |
| `doctor_rows.go:110-125` gate row | "gate present and executable" | stat + any exec bit (`:211-223`), plus not still carrying `benchSentinelMarker` (`:121-122`). No check that the gate runs or asserts anything. |
| `doctor_rows.go:139-149` profile row | a profile exists | any `projects/*.md` filename; no content check. |
| `doctor_rows.go:198-202` binding columns | green | returns `true` while listing unbound harness columns. |
| `status.go:396-427` guards signal | pre-push missing / not managed / diverted | same marker-only classifier; additionally silent unless the checkout is primary (`:403-405`) and `.bench/lines.env` exists (`:406-408`). |
| `guards.go:104-153` deny surface | "every deny-capable guard's manifest … so the block surface is learnable" | the script's own leading comment block, read as data. The `denies` cell is never derived from executable behavior. |
| `guards.go:201-227` `wired` cell | "which harness configs actually reference a hook script" | `bytes.Contains` of a path token over `.claude/settings.json` and `.codex/hooks.json` (`:219`). A token anywhere in valid JSON counts as wired; event name and matcher unverified. |
| `internal/gitguard/gitguard.go:2-5` | `denyTable` is "the single source that both classifies … and enumerates the deny classes for the guard manifest, so enforcement and advertisement cannot drift" | `denyTable` (`gitguard.go:27-45`) feeds only `denyLabels` → `internal/gitguard/verdict.go`. Nothing generates the manifest line; `.bench/hooks/block-dangerous-git.sh:4` is hand-written prose. Verified by grep over non-test Go. |

Honest self-scoping does exist where it was written deliberately:
`block-dangerous-git.sh:11-16` and `gitguard.go:7-10` ("honest-mistake layer, not
an evasion-resistant boundary"), `check-agent-line.sh:12-13,22-24` (a broken
guard warns and allows), and `SECURITY.md:5-9` ("advisory controls, not a
security boundary"). `README.md:330-332` makes the stronger operational claim
that "the git `pre-push` hook protects the default branch" — the assertion the
marker-only checks cannot support.

## Coordinator observation — this repo's installed hook

Read-only, not repaired (shaping constraint). `.git/hooks/pre-push` is 588 bytes
dated Jul 5. It carries the marker at line 2 but **no static manifest header**;
instead it answers a `--describe` flag that prints `name`/`boundary`/`denies`/
`why` at runtime. So face 2 is a manifest **protocol change** — runtime-query to
static-header — not merely a header that predates the reader.

The same file is functionally stale in two further ways: it has no `.bench`
gate-pin drift clause, and it does not re-resolve `origin/HEAD`, so it hard-codes
`main` — face 1's guessed-branch condition, live here.

The file's mtime is Jul 5, so the hand-copy repair dated 2026-07-29 in the row's
constraints is not present in the tree.
