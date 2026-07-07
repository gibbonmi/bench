# Git-error honesty (FT29)

Status: implemented

## Problem

`bench structure` lies when git itself fails. Both the all-files scan (`git ls-files`) and the `--since <base>` scan (`git diff --name-only`) discard the git error and treat an empty result as "no tracked source files to check" at exit 0. So a corrupt index, an unreadable object store, or a bad `--since` ref reports a clean tree — the exact false-clean bug that was fixed once in `bench diff` (FT19) and re-appeared here because the class was never killed, only the instance.

`bench models` has a smaller version of the same dishonesty about its inputs: it ignores argv entirely. An unknown argument prints the full inventory at exit 0 instead of rejecting the misuse, which every sibling porcelain command already does at exit 2.

Underneath both is a class defect: discarded git-helper errors scattered across the porcelain packages, where a dropped error is indistinguishable from a deliberate decision to tolerate absence.

## Solution

A porcelain command that cannot ask git the question reports the git error and exits 1 — never an empty-but-clean answer at exit 0. Concretely:

- `bench structure` on a git-query failure prints an error naming the failed git operation to **stderr** and exits 1. It never emits "no tracked source files to check" when the query failed; that clean line is reserved for the true-empty case (git succeeded, zero source files).
- `bench models` rejects any argument with a usage line at exit 2, matching the sibling norm. No new flags. The no-arg inventory and its exit-0 tolerance for unreachable *discovery* sources are unchanged — argv rejection and discovery tolerance are different contracts.

This is a class kill, not an instance fix. Every discarded git-helper error across the porcelain packages is audited and given a verdict: propagate (a dropped error can produce false-clean/false-empty output) or tolerate (an advisory surface or a value-signal that must degrade rather than crash). Tolerated drops are made visible as decisions in the code, not left as bare dropped values.

## User stories

1. As a session running `bench structure`, when git cannot answer the tracked-file query — a corrupt index breaks `git ls-files`, or `--since <bad-ref>` breaks `git diff` — I want an error on stderr naming the failed git operation and exit 1, never "no tracked source files to check" at exit 0, so a repo git cannot read is never reported as clean. Line: claude-sonnet-5 / medium. Threading the git error through the one shared query source while preserving the true-empty exit-0 boundary is mechanical plumbing at a known seam, but the boundary care lifts it above trivial.

2. As a session running `bench models` with an argument, I want an unknown argument rejected with a usage line at exit 2, so the command holds the same argv contract every sibling porcelain does instead of silently printing the inventory. Line: claude-sonnet-5 / low. This is a single guard clause mirroring the sibling default case and fully observable at the gate.

3. As the maintainer, I want every discarded git-helper error in the porcelain packages enumerated and given a propagate-or-tolerate verdict with its reason, so the false-clean class is killed across the surface instead of patched at one instance and left to recur. Line: claude-opus-4-8 / medium. The propagate-versus-tolerate call per site is a semantic judgment the gate cannot grade, the deep part the map flagged.

4. As a teammate reading the code later, I want each deliberately-tolerated git-error drop to read as an explicit decision with a one-line rationale rather than a bare `_`, so a future reader can tell a chosen tolerance from a forgotten error. Line: claude-sonnet-5 / low. Once the audit fixes each verdict, annotating the tolerant drops is mechanical.

## Implementation decisions

**`internal/structure` — the loud-error fix (stories 1, 3, 4).**

- The tracked-file query (`git ls-files`, the all-files path) and the touched-file query (`git diff --name-only … base..HEAD`, the `--since` path) each become a single error-returning source. There is exactly one `ls-files` call site and one `diff` call site; both now return their git error rather than discarding it, so the query and its error posture have one source.
- `Command` (the `bench structure` porcelain) propagates: on a query error it writes a message naming the failed git operation and the underlying git error to `os.Stderr` and returns exit 1 with no stdout report. It never renders the "no tracked source files to check" line on a query failure.
- The stderr write follows the established `resolveModel` pattern: map-dispatched commands (the `commands` map in `cmd/bench`) return only stdout, so a command that must reach stderr writes `os.Stderr` directly inside itself. The *message shape* (the "git `<op>` failed: `<err>`" text) is built by a pure helper so a package unit test can pin it; only the write itself is the thin, process-boundary-tested rim.
- Rationale for stderr specifically (decided in the map): the structured errors this repo puts on stdout (`toon.NotInRepo`, `toon.Usage`) are parseable; a raw git failure is unstructured diagnostic text. Keeping it off stdout leaves the parseable report/TOON stream clean while the failure is loud on stderr + exit 1. The pre-existing not-in-repo path (`git.Root` fails → `toon.NotInRepo` on stdout, exit 1) is deliberately unchanged — the map names only the git-query-failure case.
- The internal consumers of the same two queries keep tolerating, now as explicit decisions: `ViolationCount` (the count `bench status` reads) degrades a query failure to zero rather than crashing the SessionStart hook, and the shift refactor-gate's touched-scope reads degrade to zero because the shift loop's own `bench gate` run is the loud oracle for that worktree. Each is an explicit ignore with a one-line rationale, not a bare `_`.

**`internal/models` — argv rejection (story 2).**

- `Command` gains an argv guard at the top: any argument is unknown (the command takes none), so it returns `toon.Usage("bench models", args[0])` at exit 2, mirroring the sibling default case. No `-h`/`--help` or any flag is added (decided). The no-arg inventory path and the per-source exit-0 discovery tolerance (unreachable providers become unavailable rows) are untouched. `models` makes no git calls, so it is not part of the discard audit.

**The audit (story 3) — every discarded git-helper error, with a per-site verdict.**

Enumerated by sweeping `git.Output`, `git.Raw`, and `git.OK` across the porcelain packages. The set of discard sites is twelve; each verdict below is the deliverable the map's Handoff requires.

| # | Package / site | git call | Verdict | Reason |
|---|---|---|---|---|
| 1 | structure (all-files scan) | `ls-files` | **Propagate** | The FT29 bug: a failed query yields an empty list → "no tracked source files" at exit 0 (false-clean). |
| 2 | structure (`--since` scan) | `diff --name-only … base..HEAD` | **Propagate** | Same class: a bad `--since` ref fails the diff → empty list → false-clean at exit 0. |
| 3 | status (appendGit, dirty) | `status --porcelain` | Tolerate | Ambient advisory board the SessionStart hook consumes; a git failure must degrade the git row, not crash the hook. |
| 4 | status (appendGit, ahead) | `log @{u}..HEAD` | Tolerate | Same advisory surface; read only after an `OK`-checked upstream, and degrades to no ahead-count. |
| 5 | status (appendRetirement) | `rev-parse --abbrev-ref HEAD` | Tolerate | A failure reads as "not the default branch" → the housekeeping signal is skipped; advisory, non-fatal. |
| 6 | status (appendRoadmapReconcile) | `rev-parse --abbrev-ref HEAD` | Tolerate | Same as #5. |
| 7 | diff (branch label) | `symbolic-ref --quiet HEAD` | Tolerate | `--quiet` uses nonzero exit to signal detached HEAD — a value, not a failure; handled by the "(detached)" fallback. Already documented. |
| 8 | diff (resolveBase, branch) | `symbolic-ref --quiet HEAD` | Tolerate | Same value-signal; empty → the loud merge-base fallback. Already documented. |
| 9 | diff (resolveBase, key) | `config branch.<n>.benchBase` | Tolerate | An unset config key exits nonzero by design; empty → the loud merge-base fallback (the FT19 posture). Already documented. |
| 10 | shift (loop preflight, dirty) | `status --porcelain` | Tolerate | Immediately followed by a `rev-parse HEAD` whose error fails the loop loudly; a broken repo cannot slip past. |
| 11 | shift (dirtyPaths, staging) | `status --porcelain -z` | Tolerate | Runs inside a shift-created, validated worktree; an empty parse stages nothing → the loop's own commit / gate is the loud oracle. |
| 12 | commit (unexplained, attribution guard) | `status --porcelain -z --untracked-files=all` | Tolerate — **flag for veto** | An empty parse relaxes the guard, but the subsequent real `git commit`/`git add` fails loudly on a broken repo. This is the closest sibling to the structure bug; surfaced as a finding per the map's uncertainty rule, hardening deferred (see Out of scope). |

Already honest, no change (audit completeness): the error-*captured* sites in these packages — `status` absolute-git-dir (returns no gate row on error), `diff` commit-range and merge-base resolution (already structured exit-1 errors, FT19), and the git calls in `gate`, `worktree`, `guards`, `adopt`, and `shift`'s `rev-parse HEAD` — already handle their error. The `OK`-form calls (`status` `@{u}` existence; `diff` `cat-file -e` / `merge-base --is-ancestor`) are the test form where the exit code *is* the answer, so there is no error to drop. Only sites 1–2 change behavior; 3–12 keep tolerating.

**Tolerant-site visibility (story 4).** The tolerated sites the audit identifies read as decisions, not bare drops. The status advisory sites (#3–#6), the shift sites (#10, #11), the commit guard (#12), and the shift refactor-gate reads forced to adapt by story 1's engine change each carry an explicit ignore plus a one-line rationale. The diff sites (#7–#9) are left untouched: their tolerance is already a documented value-signal, so the decision is already visible. These annotations are review-checked, not gate-checked.

## Testing decisions

- **A good test here** exercises the external contract at the `bench` CLI subcommand boundary — the built binary in a throwaway fixture repo — observing exit code and, crucially, **stdout and stderr separately**. The stderr-vs-stdout split is the whole point of story 1, and only the process boundary can observe it (a package unit test cannot capture a `Command`-internal `os.Stderr` write).
- **Seams and prior art.** The primary seam is the runtime-contract fixture harness (`internal/contract/runtime`), already used for structure and status contracts — `f.Bench(...)` runs the binary and exposes `probe.ExitCode`, `probe.Stdout`, `probe.Stderr`. A package unit test in `internal/structure` pins the pure error-message shape (git-op name + underlying error), and a package unit test in `internal/models` can drive the argv guard directly since its usage rides stdout. When hand-running the runtime contracts, rebuild `dist/bench` first (the gate/CI build it).
- **Gate command:** `.bench/gate.sh` (the project gate). New runtime contracts land under its runtime-and-behavior layer; `bench coverage --check` (run at author time and enforced by the conformance layer) validates this map. The audit verdicts (story 3) and the tolerant-site annotations (story 4) are review-checked, not gate-checked, per the map's Handoff.

### Seam diagram

**Seam A — `bench structure` (process boundary).**

    trigger: `bench structure` / `bench structure --since <ref>`
             (real kit CLI, by-path CLI, hooks — one routed Command)
        │
        ▼
    argv ──▶ [ Command → one tracked/touched query (git) → Check ] ──▶ stdout: report, exit 0/1/2
    git  ──▶ [   query error ──────────────▶ os.Stderr + exit 1 ] ──▶ stderr: "git <op> failed: …"
                 ◀ tests attach here: run the built binary in a fixture repo;
                   assert exit code AND stdout vs stderr separately

**Seam B — `bench models` (process boundary).**

    trigger: `bench models [args…]`  (one routed Command)
        │
        ▼
    args>0 ──▶ [ Command: toon.Usage(offender) ]        ──▶ stdout: usage, exit 2
    none   ──▶ [ Command: inventory (discovery-tolerant) ] ──▶ stdout: TOON tables, exit 0
                 ◀ tests attach here: run the built binary; assert exit 2 + usage
                   for any arg, exit 0 inventory for none

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | all-files query failure is loud | Seam A | In a fixture repo, corrupt `.git/index` so `git ls-files` fails, then run `bench structure`: current code exits 0 with "structure: no tracked source files to check". New test asserts exit 1 + stderr names the git op + stdout omits the clean line → red on current code. | A false-clean answer keeps exit 0 and the clean line; the assertion on exit 1 and absent clean line fails unless the error is propagated. |
| 1 | `--since` query failure is loud | Seam A | `bench structure --since deadbeef` (nonexistent ref): current code exits 0 with the no-source line. New test asserts exit 1 + stderr names the git op → red on current code. | The bad-ref diff currently yields an empty touched list read as clean; only propagation flips exit and routes to stderr. |
| 1 | true-empty stays exit 0 (boundary) | Seam A | Not red on current code — guards against over-correction. A valid repo with zero tracked source files still exits 0 with "no tracked source files to check". | A naive fix that errors whenever the file list is empty would fail this row; it pins the git-failed-empty vs really-empty boundary the whole story turns on. |
| 2 | unknown argv rejected | Seam B | `bench models bogus`: current code exits 0 and prints the full inventory. New test asserts exit 2 + usage line → red on current code. | Ignored argv keeps exit 0 and inventory output; asserting exit 2 + usage fails unless the guard is added. |
| 2 | discovery tolerance unchanged | Seam B | Not red — regression guard. `bench models` with no args still exits 0 with the source/model TOON tables even when a provider is unreachable. | Confirms the argv guard did not leak into the no-arg path or turn unreachable-provider tolerance into a failure. |
| 3 | every discard site has a verdict | audit table (Implementation decisions) | Not TDD-able: the enumeration-and-verdict is a spec artifact, review-checked, not gate-checked (map Handoff item 5). The set is enumerated exhaustively above (12 sites). | An unenumerated "each" would let the build fix only structure; the explicit 12-row table is the veto surface that forces the class sweep. |
| 4 | tolerated drops read as decisions | code annotations | Not TDD-able: review-checked. Each tolerated site carries an explicit ignore + rationale; diff's value-signal sites are exempt as already-documented. | The gate cannot distinguish a bare `_` from a decision; only review can confirm the tolerance is stated. |

### Edge inventory

Walked per behavior against the shell-CLI hostile-input checklist in `projects/benchkit.md`.

- **Error path** — covered: structure git-query failure (both paths), models unknown argv.
- **Empty vs absent input** — covered by the true-empty boundary row: git-succeeded-empty stays exit 0; git-failed-empty is exit 1. This is the distinct-behaviors pair the checklist names.
- **Boundary values** — `bench structure --since <valid-ref>` with zero touched files takes the same true-empty path (exit 0). **Won't handle** as a separate case — it is code-identical to the true-empty boundary row already covered.
- **Malformed input** — `bench models` argv containing control bytes, glob characters, or spaces. **Won't handle** special-casing — the guard rejects at exit 2 regardless of arg content; sanitizing the echoed offender is `toon.Usage`'s existing contract, unchanged here. An in-scope caller still reaches the feature: any real argument still gets exit 2.
- **Malformed input (structure)** — a hand-edited `.bench/structure.budgets` with a bad line still warns and continues; unchanged and already covered by the existing budgets contract.
- **Interrupted / partial state** — **Won't handle**: both commands are read-only single-shot reads with no scratch state, leases, or partial writes to interrupt.
- **Re-run idempotency** — **Won't handle**: pure reads, idempotent by construction; a second run produces byte-identical output.
- **Hostile environment (git absent / not a repo)** — structure: `git.Root` fails first → `toon.NotInRepo` on stdout, exit 1 (pre-existing, unchanged, deliberately distinct from the git-query-failure stderr path); models: no git dependency, so the argv guard fires regardless of repo. Noted, no new test.
- **Invocation through every shipped surface** — real kit CLI, linked-repo by-path CLI, hooks, and adapters all route to the one `Command`; the runtime contract exercises the built binary, which is that shared implementation. Noted.
- **Paths with spaces / glob characters** — unchanged; covered by the existing structure path-with-spaces contract.

## Out of scope

- **Hardening the shift/commit write-path status-parse guards (audit sites #10–#12) into loud errors** — a separate write-path-robustness capability, not read-output honesty: these guards degrade safely because a downstream real `git` operation fails loudly on a broken repo, so no false-clean *output* escapes. Site #12 (the commit attribution guard) is the one genuinely-arguable case and is surfaced above for reviewer veto rather than silently cut. Estimate to build later: ~2 edits, 2 gate runs.
- **Adding `-h`/`--help` or any flag to `bench models`** — explicitly rejected in the decision map; a flag is a new capability, not part of this contract.
