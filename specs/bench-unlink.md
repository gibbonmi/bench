# bench unlink

Status: implemented

## Problem
A maintainer who linked Bench into a repo has no supported way to remove it. `bench link` records every installed path and its fingerprint in `.bench/link-manifest.tsv`, but nothing consumes that manifest to reverse the install. The only documented uninstall is `npm uninstall -g benchkit` plus removing the shim, which strips the global tool while leaving the per-repo footprint — the managed files, the pre-push hook, and the AGENTS.md block — behind. Removing that footprint by hand is error-prone: it is easy to delete a file the user has since edited, or to miss the hook and the fenced block, which are not manifest rows.

## Solution
`bench unlink [--dry-run]` consumes the manifest and reverses the install, sparing everything the user made theirs. It removes each managed file whose fingerprint still matches, removes managed directories that become empty, removes the pre-push hook only when it carries the managed marker, and strips the fenced Bench block from AGENTS.md while leaving surrounding prose. A file the user has edited since link is treated as the user's content now: left in place and reported. User-owned artifacts — ROADMAP.md, IDEAS.md, CONTEXT.md, `.bench/learnings.md`, `specs/`, `decisions/`, `reviews/`, the profile, a drifted `gate.sh` — are never candidates, because unlink only touches manifest rows plus the two bespoke targets. The manifest itself is removed last, and only when nothing was refused, so a partial run stays resumable. `--dry-run` prints the exact plan and touches nothing. With no manifest to consume — or one it cannot read — unlink fails loudly with exit 1 rather than silently no-op'ing on a repo it cannot account for.

## User stories

1. As a maintainer removing Bench, I want `bench unlink` to delete every managed file whose current fingerprint still matches the manifest and then remove any managed directory left empty, so the repo returns to its pre-link state without my running a dozen `rm`s.
   Line: claude-opus-4-8 / medium. This is the deep removal walk where a wrong fingerprint verdict either spares kit cruft or deletes the wrong file, so correctness outranks speed.

2. As a maintainer who has edited a managed file since linking, I want unlink to leave that file untouched and report it as kept-modified, so my edit is treated as my content rather than collateral.
   Line: claude-opus-4-8 / medium. The keep-versus-delete verdict runs on the same destructive walk and protects the user's own work, so it carries the walk's risk.

3. As a maintainer, I want the bench-managed pre-push hook removed while a foreign pre-push hook is left in place, so unlink cleans up its own enforcement without clobbering a hook I installed.
   Line: claude-opus-4-8 / medium. The hook is not a manifest row and gates only on a marker substring, so the bespoke detect-and-delete logic needs the deep-seam care.

4. As a maintainer, I want the fenced Bench block stripped from AGENTS.md while my surrounding prose survives, so my working agreement stays mine after unlink.
   Line: claude-opus-4-8 / medium. Fence-aware block removal with prose preservation is a rewrite the gate only partly observes, so it warrants the mid tier.

5. As a maintainer, I want my user-owned artifacts never touched — ROADMAP.md, IDEAS.md, CONTEXT.md, `.bench/learnings.md`, `specs/`, `decisions/`, `reviews/`, the profile, and a `gate.sh` that has drifted from its scaffold — so unlink removes the kit and not my work.
   Line: claude-opus-4-8 / medium. This is the safety contract of the command, where a single wrong inclusion destroys the user's records, so it earns full care.

6. As a maintainer, I want the manifest removed last and only when nothing was refused, so an interrupted or partial unlink leaves the residual managed state still tracked for a follow-up run.
   Line: claude-opus-4-8 / medium. Resumable-partial correctness is a property of the destructive walk's ordering, not a mechanical afterthought.

7. As a maintainer with a hand-edited manifest, I want unlink to refuse any row whose path escapes the repo root or contains path traversal, and never to follow a symlink out of the repo, so a corrupted manifest can never delete something outside the tree.
   Line: claude-opus-4-8 / medium. Hostile-input safety on a destructive command is where a wrong verdict does damage outside the repo, so it stays on the mid tier.

8. As a maintainer, I want `bench unlink --dry-run` to print the exact removal plan and change nothing on disk, so I can rehearse the destructive default before committing to it.
   Line: claude-opus-4-8 / medium. Dry-run must mirror the real walk's verdicts exactly, so any divergence between the two is a correctness bug rather than a formatting one.

9. As a maintainer, I want unlink to print a report of what it removed and what it kept as modified on both real and dry runs, so I can see the outcome at a glance.
   Line: claude-sonnet-5 / low. The report is plain text whose tokens the gate observes directly by substring, so it routes cheap.

10. As a maintainer running unlink on a repo with no manifest, or a manifest it cannot read, I want it to fail with exit 1 and a clear message rather than exit 0 having done nothing, so an unlink never silently no-ops.
    Line: claude-sonnet-5 / low. This is a pure exit-code guard the gate fully grades, so it routes cheap.

11. As a maintainer, I want `bench unlink` reachable through every shipped surface — the global CLI, the linked-repo `.bench/bin` CLI, and the built binary — routed to one implementation, so the command works from wherever I invoke it while removing Bench.
    Line: claude-sonnet-5 / low. Adding the subcommand name to the shell router and the Go dispatch is mechanical plumbing at the known strangler seam.

12. As a maintainer reading the README, I want the uninstall section to lead with `bench unlink [--dry-run]` as the way to remove the per-repo footprint, keeping the package-and-shim removal as the second step and naming the manual path for pre-manifest installs, so the documented uninstall matches the new command.
    Line: claude-opus-4-8 / medium. This is user-facing prose the gate only sweeps for stale references, not for correctness, so it bumps one tier above cheap.

## Implementation decisions

- **New `Unlink` entry in `internal/adopt`.** Add `Unlink(args, stdout, stderr) int` and dispatch it from `adopt.Run` under `case "unlink"`; `cmd/bench/main.go` reaches it by adding `"unlink"` to the existing `link, init, doctor` group that calls `adopt.Run`. Unlink reuses the manifest reader, the fingerprint helper, the hooks-dir resolver, and the AGENTS marker helpers — no second copy of any of them.
- **Routing.** `bin/bench.sh` routes `unlink` via the plain porcelain path (`route_porcelain`), not the adoption path. Adoption routing refuses to run outside a real kit asset tree, but unlink operates purely on the target repo's manifest and files and needs no kit source assets, so it must succeed from the consumer repo through any shipped surface, including the linked-repo by-path CLI.
- **Manifest-absence guard is explicit.** The shared manifest reader maps an absent file to an empty manifest with no error; unlink must not inherit that false-empty. It checks the manifest's presence and readability first: absent or unreadable is exit 1 with a message; present and readable proceeds, even when the manifest carries zero managed rows (a degenerate but valid manifest reports nothing to remove).
- **Fingerprint verdict per row.** For each managed row, unlink re-fingerprints the current on-disk path with the same `Lstat`-based helper link used. Match removes the path; a fingerprint mismatch keeps it and reports it as kept-modified; an already-absent path needs no action. The helper never follows symlinks, so an adapter row (whose stored fingerprint is the symlink's) matches and the symlink itself is removed, never its target.
- **Path safety before any removal.** Every row's path is resolved against the repo root and rejected unless it stays strictly inside the root — absolute paths, `..` traversal, and any resolved path outside the tree are refused and reported, never removed. A row whose target is a directory or otherwise non-regular cannot be fingerprinted as a managed file and is skipped-and-reported rather than deleted. Each such rejection counts as a refusal.
- **Empty-directory sweep.** After file removals, the parent directories of managed rows are pruned deepest-first with a non-recursive rmdir that only succeeds on an empty directory, so any directory still holding a kept-modified file or a user artifact (for example `.bench/`, which retains learnings.md and gate.sh) survives untouched.
- **Pre-push hook is bespoke.** The hook is not a manifest row. Unlink resolves the effective hooks directory the same way link does (honoring `core.hooksPath`) and removes `pre-push` only when its body contains the managed marker; a hook without the marker is a foreign hook, left in place, and its presence is not a refusal because it was never Bench's.
- **AGENTS.md fence is bespoke.** AGENTS.md is not a manifest row. Unlink strips the fenced managed block using the existing fence-aware marker scan so fenced examples and user prose are preserved. When stripping leaves the file empty or whitespace-only — the case where link created the file with no user content — the file is removed; otherwise the remaining prose is written back. This mirrors link's create-if-absent symmetry.
- **CLAUDE.md is left in place.** It is bench-authored but outside the decision's removal set and is not manifest-tracked, so unlink does not touch it.
- **Manifest removed last, conditionally.** The manifest file is removed only when the run recorded no refusals — no kept-modified files and no rejected rows. When anything was refused, the manifest stays so the residual managed state remains tracked and a later run can finish the job.
- **`--dry-run` is the same walk, a different sink.** It computes every verdict and prints the same plan but performs no filesystem writes and never removes the manifest, so a dry run's on-disk state is byte-identical to before.
- **Report is plain text.** The real run reports removed paths, an empty-directory count, kept-modified rows by path, and any refused rows; the dry run prints the same plan under a dry-run banner. Tokens are stable enough for substring assertions, matching how the adopt family reports today; this is a mutating porcelain command, not part of the AXI query surface, so it emits plain text rather than TOON.
- **README uninstall section rewritten.** It leads with `bench unlink` to remove the per-repo footprint and `bench unlink --dry-run` to rehearse it, retains `npm uninstall -g benchkit` plus shim removal as the step that removes the global tool, and states that a repo linked before manifests must be removed by hand because unlink exits 1 with nothing to consume.

## Testing decisions

- **What a good test is here.** Drive the built `bench` binary through `bin/bench.sh` against a throwaway fixture repo and assert on exit code, on the resulting tree, and on report substrings — never on a reading of the implementation. Link the fixture, mutate it to set up the case, run `unlink` (or `unlink --dry-run`), then check what survived.
- **Seam and prior art.** One seam: the `bench unlink` subcommand exercised through the adopt/link surface contract family in `internal/contract/surface` (neighbor to the `bench link` contract tests). Prior art to mirror: the fresh/relink link test, the modified-managed conflict test, the existing-AGENTS.md preservation test, and the hooksPath test — all of which link a fixture and assert on files, exit codes, and message substrings.
- **Gate command.** The project gate, `.bench/gate.sh`. The new tests run inside its safe-link/surface contract layer.

### Seam diagram

    trigger: maintainer runs `bench unlink [--dry-run]`
             (global CLI · linked-repo .bench/bin CLI · built binary)
        │
        ▼
    .bench/link-manifest.tsv rows ─▶ [ adopt.Unlink                     ] ─▶ stdout: plain-text report
    current on-disk fingerprints ──▶ [  guard: manifest absent → exit 1  ]     (removed / empty-dirs /
    AGENTS.md + pre-push hook ─────▶ [  fingerprint each row → verdict    ]      kept-modified / refused)
    --dry-run flag ────────────────▶ [  strip fence · drop managed hook   ] ─▶ filesystem (real run only):
                                     [  sweep empty dirs · manifest last   ]     managed files/dirs gone,
                                     [  --dry-run: plan only, no writes     ]     fence stripped, hook removed
                                          ◀ tests attach here: link a throwaway repo, mutate it,
                                            run built `bench unlink` via bin/bench.sh, assert on
                                            exit code + surviving tree + report substrings

### Acceptance coverage map

Every row starts red today because `bench unlink` is unrouted (the binary exits 2, `unknown subcommand: "unlink"`); each row additionally names the degenerate implementation it pins against once the subcommand exists.

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | link then unlink leaves only user content; managed files and now-empty managed dirs are gone | unlink contract fixture | link a fixture, `bench unlink`, assert a known managed file and its emptied dir no longer exist and exit 0 — red against a no-op stub that removes nothing | a stub that reports but deletes nothing leaves the managed file on disk, failing the not-exist assertion |
| 2 | a managed file edited after link is kept and reported, not deleted | unlink contract fixture | link, append to a managed file, `bench unlink`, assert the file still exists with the edit and the report contains its path under "modified" — red against an unconditional `rm` of every manifest row | a walk that deletes by row without a fingerprint check removes the edited file, failing the still-exists assertion |
| 3 | the managed pre-push hook is removed; a foreign one is left | unlink contract fixture | after link assert `pre-push` gone post-unlink; in a second fixture plant a non-managed `pre-push`, unlink, assert it survives — red against removing `pre-push` unconditionally, or ignoring it entirely | unconditional removal deletes the foreign hook (survives-assertion fails); ignoring it leaves the managed hook (gone-assertion fails) |
| 4 | the AGENTS.md fenced block is stripped while user prose survives | unlink contract fixture | link over an AGENTS.md carrying "PROJECT RULES", unlink, assert the file still contains "PROJECT RULES" and no `bench:start` marker — red against deleting AGENTS.md wholesale or leaving the block | deleting the file loses the prose; leaving the block leaves the marker — either assertion fails |
| 5 | user artifacts are never removed | unlink contract fixture | link, create ROADMAP.md/IDEAS.md/CONTEXT.md/`.bench/learnings.md`/a drifted `.bench/gate.sh` and a file under `specs/`, unlink, assert every one survives — red against any implementation that recurses `.bench/` or the repo instead of walking manifest rows | a recursive delete of `.bench/` or the tree removes learnings.md or gate.sh, failing a survives-assertion |
| 6 | the manifest is removed only when nothing was refused | unlink contract fixture | clean case: assert manifest gone after unlink; refusal case (one modified-managed file): assert manifest still present after unlink — red against always removing the manifest | always removing the manifest drops it in the refusal case, failing the still-present assertion and breaking resumability |
| 7 | a traversal or root-escaping manifest row is refused, not deleted | unlink contract fixture | hand-edit the manifest to add a `../outside` row pointing at a file created outside the repo, unlink, assert that file survives, the row is reported refused, and the manifest is kept — red against `os.Remove(filepath.Join(root, rel))` with no bounds check | joining root+rel without validation deletes the outside file, failing the survives-assertion |
| 8 | `--dry-run` prints the plan and changes nothing | unlink contract fixture | link, `git add -A && commit`, `bench unlink --dry-run`, assert `git status --porcelain` is empty and the report shows the would-remove plan — red against a dry-run that shares the real removal path | a dry-run that actually mutates leaves a dirty working tree, failing the clean-status assertion |
| 9 | the report lists removed and kept-modified paths on both runs | unlink contract fixture | run the modified-managed case and assert the report contains the kept path and a removed path on the real run, and the same plan on the dry run — red against a silent run that prints nothing | an empty report fails the substring assertions for the kept and removed paths |
| 10 | absent or unreadable manifest exits 1; a second unlink after a full run also exits 1 | unlink contract fixture | in a repo with no manifest run `bench unlink` and assert exit 1 with a "no manifest"-class message; after a clean unlink run it again and assert exit 1 — red against inheriting the reader's absent→empty→exit-0 path | a proceed-on-empty implementation exits 0 on a repo with no manifest, failing the exit-1 assertion |
| 11 | unlink resolves to one implementation from the global CLI and the linked-repo by-path CLI | unlink contract fixture | after link, run `.bench/bin/bench.sh unlink --dry-run` from the fixture and assert it prints the plan and exits 0 (not the adoption-route "real Bench kit" refusal) | routing unlink through the adoption path makes the by-path CLI refuse, failing the exit-0 assertion |
| 12 | the README uninstall section names `bench unlink` | root-conformance stale-reference sweep | grep the README uninstall block for `bench unlink`; the existing stale-command/reference conformance check stays green only if no dangling reference is introduced | a rewrite that drops the command name or leaves a dangling reference is caught by the sweep the gate already runs |

Cheapest wrong implementations checked: a no-op stub (story 1 bites), an unconditional per-row `rm` (stories 2 and 7 bite), a recursive `.bench/`/tree delete (story 5 bites), an always-remove-manifest (story 6 bites), a dry-run that shares the mutation path (story 8 bites), and a reader that treats absent as empty-and-proceed (story 10 bites).

### Edge inventory

Edge classes walked per behavior, each resolved as a coverage row above or a Won't-handle line here.

- **Error path** — absent or unreadable manifest → row (story 10).
- **Empty vs absent input** — a present-but-rowless manifest proceeds and reports nothing, distinct from an absent manifest that exits 1 → edge of story 10 (assert the rowless case exits 0 with an empty plan).
- **Malformed input** — traversal `../`, absolute path, and a symlink row pointing outside the repo → row (story 7); a manifest row whose target is a directory or other non-regular file is skipped-and-reported, never deleted → edge of story 7.
- **Missing trailing newline** — a manifest or AGENTS.md whose last line has no newline is parsed the same way link parses it (line scan drops the newline) → edge of story 4.
- **Re-run idempotency** — a second unlink after a clean run finds the manifest already gone and exits 1 loudly → edge of story 10.
- **Interrupted/partial state** — because the manifest is removed last and only on no-refusal, a run cut short mid-removal leaves the manifest present so a re-run finishes the job → edge of story 6 (the resumability is the tested property).
- **Hostile environment / paths with spaces or glob characters** — Go's path removal is literal with no shell or glob expansion, so a manifest path containing spaces or `*` is removed by exact match → edge of story 1.
- **Invocation surface / symlink invocation** — reachability from the global CLI and the linked-repo by-path CLI → row (story 11); invocation through a symlinked wrapper and a cwd below the repo root inherit the wrapper's existing path resolution and `git.Root()`.
- **SIGINT mid-run atomicity** — Won't handle: unlink is not transactional; the manifest-removed-last ordering makes a partial run resumable by re-running, which is the guarantee offered in place of atomicity.
- **Control bytes in report paths** — Won't handle: the report is plain text, not a TOON table, and its paths come from the kit-controlled manifest, so the TOON control-byte refusal does not apply.
- **CLAUDE.md** — Won't handle: it is bench-authored but outside the decision's removal set and untracked by the manifest, so a bench-created CLAUDE.md persisting after unlink is the accepted cost of not tracking it.
- **Shell-glob expansion of manifest paths** — Won't handle: there is no shell layer between the manifest and removal, so a path with glob characters is never expanded.

## Out of scope

- **`--force` removal of modified-managed files** — a distinct destructive mode with its own confirmation semantics, explicitly rejected by the decision map because deleting user-edited files is the reviewer's hand, not a flag; revive later at ~3 edits, 2 gate runs.
- **Interactive confirmation prompt** — a separate interactive mode that would contradict the non-interactive porcelain rule; if ever wanted it is its own capability at ~2 edits, 1 gate run.
- **Automatic footprint removal for pre-manifest (un-manifested) installs** — a separate heuristic feature that would detect and remove a Bench footprint with no manifest to consume, distinct from the decided exit-1-plus-documented-manual-path behavior; build later at ~4 edits, 3 gate runs.
