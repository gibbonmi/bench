# Pre-push guard visibility

Status: ready

## Destination

FT135's three faces decided as one owner: an enforcement control that reports its
own status must report what it actually verified. Face 1 — the installed pre-push
hook's protected branch may be a baked guess and nothing says so. Face 2 — a
managed hook may be stale against the embedded template and every surface reports
it as fine, or as inert, without checking. Face 3 — `bench link`, the sanctioned
repair, refuses in exactly the repository that needs it. The map decides which
surface owns the report, what the currency comparison is, what a mismatch does,
how link is unblocked, and whether the wider class R6 found is in this row or a
separate one. Out at the far end: `/bench-write-spec` from a ready map.

## #1: How does the installed hook determine the branch it protects?

Blocked by: none
Type: Research

### Question
Cover both the live resolution path and the fallback, what happens when live
resolution returns empty, and whether anything records which path was taken.

### Answer
Two resolutions. Install time bakes `git.ResolvedDefault`'s answer, falling back
to the literal `main` (`internal/adopt/link_hook.go:98-106`); the guessed case is
narrow — no resolvable `origin/HEAD`, no commit named `main`, and not exactly one
local branch — but real. Runtime re-resolves `origin/HEAD` on every push and
keeps the baked token only when that lookup is empty
(`internal/adopt/prepush.sh:17-21`), silently: no warning, no degraded posture.
**Nothing anywhere records which path produced the branch** — not the manifest
header, not `ClassifyPrePush`, not `prePushRow`. The branch appears only in the
block message, which names it without saying where it came from. Detail:
`assets/pre-push-guard-visibility-research.md` (R1).

## #2: What does `bench doctor` check on installed git hooks?

Blocked by: none
Type: Research

### Question
Every check, and for each whether it inspects a marker, file content, or
something else.

### Answer
Exactly one hook check exists — the pre-push backstop row
(`internal/adopt/doctor.go:230-232,247-265`) — and its entire managed test is the
substring `bench:managed-pre-push` (`link_hook.go:32,59`), plus hooks-dir
resolution and `core.hooksPath` diversion. No executable bit, no baked branch, no
body beyond the marker; no doctor check of any other git hook or any harness
hook. Substring matching is a documented deliberate choice (`link_hook.go:26-31`)
to avoid false-reds across branch substitution and benign template evolution.
Doctor never repairs the hook (`doctor.go:243-246`). Detail: asset (R2).

## #3: How does `bench guards` determine manifest state and boundary?

Blocked by: none
Type: Research

### Question
The mechanism, the output when a managed hook carries no manifest header, and
where that output is consumed.

### Answer
`parseHeader` (`internal/guards/guards.go:119-133`) reads only the leading
comment block for four required keys (`name`, `boundary`, `denies`, `why`,
`guards.go:90`). The boundary cell is the parsed value with no inference — an
incomplete header yields `""` and a definitive `no manifest` row
(`guards.go:138-153`). Consumed as TOON by `bench guards` and as one line per
guard by SessionStart through `--brief`
(`internal/sessioninspect/sessioninspect.go:84`, wired at
`.claude/settings.json:22` and `.codex/hooks.json:8`). **No gate or conformance
check enforces manifest presence on an installed hook**, and `bench status` uses
a different classifier that emits no row at all when the hook is managed
(`internal/status/status.go:409-413`), so it cannot surface this. Detail: asset
(R3).

## #4: Does `bench link` abort before it can refresh the hook?

Blocked by: none
Type: Research

### Question
Face 3's sequencing claim, recorded from capture and never re-verified: does the
symlink-parent conflict land before the pre-push refresh?

### Answer
**Confirmed.** The symlink-parent abort returns `1, false` at
`internal/adopt/link_transaction.go:73-76`, inside the plan loop; the hook is
staged at `:186-199` and nothing reaches disk until `promoteAll` at `:222`.
Verified in this repo: `.claude/commands` is a symlink to `../.agents/commands`
and `.claude/commands/bench-debug.md` resolves through it, so `bench link` here
cannot reach the hook refresh. Face 2's sanctioned repair is unavailable exactly
where the stale hook lives. Unplanned finding: `installGitHook`
(`link_hook.go:108-126`) has no call site, so the token-substitution knowledge
exists twice with one copy dead. Detail: asset (R4).

## #5: Is there precedent for comparing an installed artifact against a declared source?

Blocked by: none
Type: Research

### Question
How the template is stored, whether anything compares an installed hook against
it, and what comparison mechanisms exist elsewhere in the repo.

### Answer
The template is `//go:embed`-ed at `link_hook.go:19-20`; `prePushTemplate` is
read at exactly two sites and both are writes. Nothing compares an installed hook
against it, and `PrePushState` has no drifted or stale value
(`link_hook.go:37-42`). Precedent exists and is well-established: manifest
SHA-256 fingerprinting for every linked kit file (`internal/adopt/manifest.go:101-131`,
`manifestOwnedClean` at `link.go:252-262`, `upgradePlanCounts` at
`upgrade.go:116-138`), plus release-package digest comparison
(`internal/releaseevidence/package_artifact.go:190`) and source-tree fingerprints
(`evidence_fingerprint.go:31-79`). The hook's exclusion from that mechanism is a
recorded decision — "The hook is bespoke, not a manifest row"
(`internal/adopt/unlink.go:265`) — which #9 must either keep or reopen
explicitly. Detail: asset (R5).

## #6: Does the failure class recur beyond the pre-push hook?

Blocked by: none
Type: Research

### Question
Every place an enforcement control reports its own status, what it asserts, and
what it actually checks.

### Answer
**Yes, at least three further instances.** The doctor shim row reports "healthy"
from a marker plus an executable target, never comparing the body against
`ShimContent` — that comparison exists only in `doctorFix`
(`doctor.go:196,216-228,305-308`). The doctor gate row asserts the gate is
present and executable and checks stat plus an exec bit, never that the gate
asserts anything (`doctor_rows.go:110-125`). The `guards` `wired` cell claims to
report which harness configs "actually reference" a hook and checks
`bytes.Contains` of a path token over the JSON (`guards.go:201-227`), so event
name and matcher go unverified. And `internal/gitguard/gitguard.go:2-5` claims
`denyTable` is the single source for both classification and the guard manifest's
advertisement, but `denyTable` feeds only `denyLabels`: the manifest line at
`.bench/hooks/block-dangerous-git.sh:4` is hand-written prose. That last one is
an enforcement/advertisement duplication the repo's own code standard names as a
defect. Honest self-scoping does exist where written deliberately
(`SECURITY.md:5-9`, `check-agent-line.sh:12-13`), while `README.md:330-332`
makes the stronger claim the marker-only checks cannot support. Detail: asset
(R6).

## #7: Which surface owns the resolved-versus-guessed and currency report?

Blocked by: #1, #2, #3, #6
Type: Grill

### Question
`doctor`, `guards`, the SessionStart banner via `bench status`, `link`'s output,
or more than one. The row named "doctor and link" but explicitly left it open.
The three surfaces already disagree about the hook: doctor renders a row from
`ClassifyPrePush`, status renders a row from the same classifier but only on
failure, and guards runs its own independent parse. Adding the report to all
three multiplies the classifier, which the repo's one-source standard forbids —
so this ticket also decides where the single hook-health fact is computed.

### Answer
`bench doctor` is the sole owner of the hook-health computation (reviewer,
2026-07-31); `bench guards` and `bench status` render from it rather than each
re-deriving. This collapses an existing duplication as well as preventing a new
one: doctor and status share `ClassifyPrePush` today but disagree about when to
show a row, and `guards` runs an entirely independent header parse, so a report
added to all three would be the fourth and fifth derivation of one fact.
Residual: whether the SessionStart banner surfaces the resolution and currency
fields specifically — as opposed to merely rendering from the same computation —
is decided in #8, since that is where the rendering shape is chosen.

## #8: Is a guessed branch a warning, a failure, or a status field?

Blocked by: #7
Type: Grill

### Question
The narrow trigger (#1) argues against reddening doctor for a condition the
runtime path usually corrects; the fact that this repo is living in exactly that
condition argues the other way.

### Answer
A status field, green either way (reviewer, 2026-07-31). Doctor's pre-push row
carries a resolved-versus-guessed field and still exits 0, and `guards` and the
SessionStart banner render that same field from doctor's computation — which
settles #7's residual: the banner does surface the field, not merely the row.
Reddening doctor for a condition the runtime path usually corrects would train
readers to ignore reds, and visibility is this row's objective rather than
enforcement.

## #9: What is the currency comparison, and does the hook become a manifest row?

Blocked by: #5, #7
Type: Grill

### Question
Exact bytes against the substituted template, a version token in the manifest
header, or folding the hook into the SHA-256 manifest the rest of link already
uses. Byte identity was rejected once as too brittle (`link_hook.go:26-31`) and
the manifest-row route was closed once (`unlink.go:265`); reopening either is the
reviewer's call. The evidence that sharpens this: the stale hook here fails on a
*protocol* change (runtime `--describe` versus static header), which a version
token would catch and a marker never can.

### Answer
Exact bytes against the embedded template with the branch token substituted
(reviewer, 2026-07-31): recompute
`strings.ReplaceAll(prePushTemplate, prePushBranchToken, protectedBranch(root))`
and compare against the installed file. It introduces no new fact to maintain,
where a version token only catches drift someone remembered to bump — and the
stale hook here drifted across a protocol change a conscientious author would
have bumped, so a token protects exactly where care was already taken. The
recorded objection at `link_hook.go:26-31` does not survive the specifics:
substitution is handled by recomputing rather than comparing raw, and benign
template evolution self-corrects because `bench upgrade` refreshes the hook
through `transactionalLink` (`upgrade.go:73`). Accepted cost: the comparison is
binary, so it cannot distinguish one release behind from a dead protocol. The
hook does **not** become a manifest row — ruled out on structure, since the
manifest addresses repo-relative paths and the hook lives in `.git/hooks`
outside the worktree, so reuse would need a second addressing scheme.
`unlink.go:265` therefore stands.

## #10: What does a currency mismatch do, and is repair ever silent?

Blocked by: #9
Type: Grill

### Question
Report only, red the doctor row, or offer repair. Doctor's no-self-heal rule is
recorded (`doctor.go:243-246`) and the drain offered the repair rather than
performing it, so the default is report-and-name-the-remedy — but the remedy is
the thing face 3 shows is broken.

### Answer
A stale managed hook reds doctor's pre-push row and names the remedy; repair is
never automatic, but `bench doctor --fix` gains it (reviewer, 2026-07-31). That
path is already explicitly invoked and already compares the shim body against
`ShimContent` (`doctor.go:305-308`), so extending it to the hook is opt-in
rather than silent and the no-self-heal rule at `doctor.go:243-246` stands
unamended. Note the ordering dependency: the remedy is only trustworthy once
#11 lands, since `bench link` currently refuses in exactly the repositories
that need it.

## #11: How is `bench link` unblocked on a symlinked managed directory?

Blocked by: #4
Type: Grill

### Question
Traverse the symlinked directory, skip already-converged files before the
conflict check, or move the hook refresh ahead of the plan loop. The third is the
narrowest fix for face 2's repair route but leaves link still refusing on this
repo for everything else; the first two change link's conflict semantics for all
managed paths. Wider than the row stated: `bench upgrade` calls the same
`transactionalLink` (`internal/adopt/upgrade.go:73`), so this abort blocks both
adoption routes, and upgrade's `changed` count is derived from plan entries only
(`upgrade.go:54`) while the hook is staged outside the plan — so
`bench upgrade --check` never reports the hook refresh it would in fact perform.

### Answer
Skip already-converged files before the symlink-parent abort (reviewer,
2026-07-31): a plan entry whose target already matches the manifest hash needs
no write, so a symlink parent is not a conflict for it. This unblocks both
adoption routes without changing conflict semantics for genuinely new writes,
and it is sufficient here because this repo's `.claude/commands` entries are all
converged. Traversing the symlink was rejected as too permissive — it would
write through a directory the reviewer symlinked deliberately — and reordering
the hook refresh was rejected as leaving the repository un-upgradable for
everything else.

## #12: One row, or a prefactor plus consumers?

Blocked by: #6, #7
Type: Grill

### Question
#6 found the class in at least three more places (shim body, gate row, `wired`
cell) plus one enforcement/advertisement duplication in `gitguard`. Options: keep
FT135 to the pre-push hook and park the rest; make FT135 a prefactor — one
"reports what it verified" seam — with the other sites as consumers; or split the
`gitguard` duplication out as its own row since it is a standards defect rather
than a visibility gap.

### Answer
Prefactor plus consumer rows (reviewer, 2026-07-31). FT135 builds the pre-push
report as a reusable "reports what it verified" seam; the doctor shim row, the
doctor gate row, and the `guards` `wired` cell become their own roadmap rows
that consume that seam. The `gitguard` `denyTable` duplication splits out
separately because it is a one-source-per-fact defect rather than a visibility
gap — enforcement and its advertisement claim a shared source they do not have.
This keeps FT135's build small and gives the class an owner instead of letting
it recur a fourth time. The consumer rows and the `gitguard` row are
`/bench-what-next`'s to create; this map does not edit `ROADMAP.md`.

## #13: Does the local hook get repaired, and when?

Blocked by: none
Type: Task

### Question
Owner: reviewer. `.git/hooks/pre-push` is dated Jul 5 and speaks the old
`--describe` protocol; it has no gate-pin drift clause and hard-codes `main` with
no live re-resolution. The hand-copy repair the row's constraints date to
2026-07-29 is not in the tree. Repairing it destroys the live evidence this row
is built on, so shaping did not touch it — but the repo is meanwhile running a
guard that enforces less than the current template does.

### Answer
Repair after the spec is written (reviewer, 2026-07-31). The stale hook is the
one live instance of the dead-protocol case, so `/bench-write-spec` writes its
coverage against it first and the repair follows. Until then the repo runs a
guard that enforces less than the shipped template does — no gate-pin drift
clause, and a hard-coded `main`. The repair itself is a `bench link` once #11
lands, or a hand-copy before that.

## #14: Does the kit owe older linked repos a migration path?

Blocked by: #9, #10
Type: Grill

### Question
Repos linked at older kit versions may carry hooks speaking the `--describe`
protocol. Under #9's exact-byte comparison every one of them reds, and under
#10 the remedy is `bench doctor --fix` or a relink. Is that sufficient, or does
the kit owe them a recognizing-and-migrating path that names the old protocol?

### Answer
No migration path (reviewer, 2026-07-31). A stale hook reds identically whatever
made it stale, and `bench doctor --fix` repairs it. Protocol-specific
recognition would mean teaching the kit to parse a format it no longer emits —
a second source of truth about hook shape, which is the defect class this row
exists to remove. Auto-repairing dead-protocol hooks without `--fix` was
rejected as silent repair of an enforcement control, which #10 forbids.

## #15: Does `README.md`'s protection claim get narrowed?

Blocked by: #8
Type: Grill

### Question
`README.md:330-332` states that "the git `pre-push` hook protects the default
branch." Under #8 a guessed branch stays green, so the claim can still overstate
what the installed hook does. Narrow the prose to what the checks support, or
leave it and let the new doctor field carry the caveat.

### Answer
Narrow the prose (reviewer, 2026-07-31): the hook protects the branch it
resolves, with the doctor field naming which one that is. This is the one place
in the repo making a stronger operational claim than any check supports, so
closing it is part of the same objective rather than a doc chore.

## Not yet specified

## Spec-writer discretion

- The name and shape of any new `PrePushState` value or status field, provided
  the four existing states keep their meanings.
- Where the single hook-health computation lives once #7 names its owner, so long
  as no second classifier appears.

## Out of scope

- Making guards or hooks an evasion-resistant boundary. `SECURITY.md:5-9` scopes
  them as advisory controls against honest mistakes; this row is about honest
  reporting, not about raising the threat model.
- Doctor self-healing the hook as a default behavior, unless #10 explicitly
  reopens it — the no-self-heal rule at `doctor.go:243-246` is a recorded
  decision.

## Sources

- Path: `decisions/assets/pre-push-guard-visibility-research.md`
  Supports: #1 through #6, and the evidence clauses in #7, #9, #11, #12, #13 — six read-only research delegations run 2026-07-31, each claim re-verified against the tree by the coordinator.
  Drift: invalidated if the pre-push template, `ClassifyPrePush`, the guards header parser, or `bench link`'s conflict ordering changes; re-verify the cited line numbers before relying on it.
