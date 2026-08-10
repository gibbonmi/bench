# AXI coherent diff

Status: staged

Decision source: `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`

## Problem

Routine Git inspection still escapes the AXI surface: orienting on a change takes repeated raw `git status`/`diff --stat`/`diff`/`diff --check`/`rev-parse`/`log`/`show` batches because `bench diff` answers only part of the question. Its default output has no aggregate counts, no untracked paths, no index-versus-worktree distinction, and no whitespace result; `--full` silently omits untracked regular-file bodies; and nothing defends against HEAD, index, or worktree moving mid-read, so one invocation can splice facts from concurrent states.

## Solution

Make `bench diff` the one coherent Git-inspection snapshot in a single atomic output migration under its existing command name, composing the owners that already exist: `internal/git.Facts` for branch, divergence, and XY checkout status, and `internal/diff`'s range owner for base-relative facts. The default response becomes one orientation snapshot; `--full` adds the log and an exact patch including untracked regular files; `--commit` keeps its post-landing meaning. Contextual actions compose the `axi.Action` owner landed by `axi-spec-build-complete`. No `bench git` namespace and no second porcelain parser.

## User stories

1. As an agent orienting on changes, I want one `bench diff` call to return the revision row, pre-computed commit/file/insertion/deletion/staged/unstaged/untracked counts, the changed-file inventory including untracked paths, an index-versus-worktree checkout table, and the whitespace-check result, so orientation replaces the raw status/stat/name-only/revision/check sequence. Line: gpt-5.6-terra / high. The snapshot spans two existing owners and omission of any fact silently re-opens a raw-Git turn.
2. As an agent reviewing bodies, I want at most one `bench diff --full` call to return the same snapshot plus the landed-commit log and an exact patch that does not silently omit untracked regular files, so review needs no follow-up Git calls. Line: gpt-5.6-terra / high. Untracked-body inclusion and binary/hostile-path handling are the classes a plausible implementation quietly drops.
3. As an agent reading during concurrent writes, I want a snapshot whose HEAD, index, and worktree identity is captured around the reads, retried once on movement, and otherwise refused with a structured drift error, so no response ever splices facts from different states. Line: gpt-5.6-terra / high. A torn snapshot is byte-plausible and only an injected mid-read mutation can prove the defense.
4. As an agent, I want each diff response to end with `help[]` derived from the snapshot facts — bounded orientation advertises `bench diff --full`, drift advertises retrying the exact same invocation, and a clean or complete result advertises nothing — so the follow-up call is exact. Line: gpt-5.6-terra / medium. The derivation composes the settled action owner at one seam.

## Implementation decisions

- `internal/diff` remains the single owner. It consumes `git.Facts` for branch, divergence, and porcelain XY entries and its own `diffRange` for base-relative facts; it adds no new porcelain parser and `--commit` resolution is unchanged.
- The default snapshot is TOON: a one-row revision table (branch, base, method, head), a one-row aggregate table (commits, files, insertions, deletions, staged, unstaged, untracked), the `files[N]{status,path}` inventory augmented with untracked paths, a minimal checkout table distinguishing index from worktree state, and a one-row whitespace result derived from `git diff --check`.
- `--full` returns the same snapshot plus the existing `log[N]{sha,subject}` table and the exact patch body. The patch stays raw after the `diff_body:` marker, exactly as today; untracked regular files append after the tracked patch in path-sorted order as standard Git new-file patches — the exact bytes of `git diff --no-index -- /dev/null <path>` per file, including mode headers and Git's binary form — and non-regular untracked entries are named in the inventory rather than read.
- `--commit <sha>` renders the post-landing view of an immutable commit: the revision row identifies the commit and its parent base with `method: commit <sha>` and carries no live branch cell, aggregates derive from that commit's range and file set, and the live-checkout facts — staged/unstaged/untracked counts, the checkout table, and drift capture — do not apply and are omitted rather than zeroed or taken from the unrelated live tree.
- Aggregate scope is exact: `commits` counts the log range, `files` counts inventory rows including untracked entries, `insertions`/`deletions` cover the tracked base-relative patch exactly as `git diff --shortstat`, `staged`/`unstaged` derive from the porcelain XY split, `untracked` counts untracked entries, and the whitespace result covers the tracked patch only, exactly as `git diff --check`. Untracked bodies never inflate insertion counts.
- The file inventory schema is `files[N]{status,path,kind}`: `kind` is blank for regular entries and names a non-regular untracked entry or dangling symlink (`fifo`, `socket`, `device`, `symlink`, `dangling-symlink`), so an omitted body is always explained by its own row.
- Untracked enumeration uses all-files semantics (`--untracked-files=all`): a new directory expands to its individual file entries and never collapses to one `dir/` row, so every untracked regular file — however deeply nested — appears in the inventory, the `untracked` count, and the `--full` patch.
- Snapshot identity for the live views is content-sensitive down to file bytes: the HEAD commit id, a digest of the raw index file bytes, a digest of the complete porcelain status byte stream, and a per-path digest of every reported dirty and untracked entry covering exactly what the patch can observe for its kind — content bytes and file mode for a regular file, the submodule's current HEAD id for a gitlink, the target for a symlink, stat identity for any other kind — are each captured before and after the reads; movement in any one retries once, then returns a structured drift error whose help row carries the exact same invocation. No partial snapshot is ever emitted. A rewrite or chmod of an already-dirty file, and a moved HEAD in a dirty submodule, each change their per-path digest even though HEAD, index, and porcelain bytes stay identical, so no patch-observable state sits outside identity; a write that lands entirely after the closing capture is post-snapshot by definition. Non-regular non-gitlink entries are stat-compared, never read.
- Contextual actions compose `axi.Action` from `axi-spec-build-complete`; a clean tree, a complete `--full` body, and `--commit` results emit the honest empty help block.
- Structured errors and the 0/1/2 taxonomy are unchanged: refusals stay `toon.Errorf` on stdout exit 1, usage stays exit 2, and the existing base-resolution error kinds are preserved.
- This is one atomic migration under the existing name: the reviewed old-to-new fixture set names the exact delta per response mode (new snapshot blocks, untracked inclusion, drift refusal, appended help), and the acceptance target — one orientation call, at most one `--full` call — is graded against the fixture matrix.

## Testing decisions

- TDD attaches at the snapshot composition seam with real repositories built per fixture; every fact is asserted against raw Git output captured in the same fixture, never against a second Bench derivation.
- The paired-delta matrix covers committed, staged, unstaged, untracked, nested-new-directory, rename, deletion, binary, hostile-filename, clean, and mid-read-drift cases; each hostile-state class carries one independent mutation.
- Drift tests inject a mutation between reads through a test seam and assert retry-once, then the structured error; no timing-dependent sleeps.
- Help derivation tests reuse the `axi.Action` counterexample discipline; the package gains real `Command`-level tests (today only the parser is tested).

### Seam diagram

    trigger: bench diff [--full|--commit] on a live checkout
        │
        ▼
    git.Facts + diffRange reads ──▶ [ identity-checked snapshot composer ] ──▶ TOON snapshot (+ log/patch) + help[]
                                          ◀ tests attach here: raw-Git paired fixtures and injected mid-read drift

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| CD1 | 1 | The default response carries the revision row and a pre-computed aggregate row — commits, files, insertions, deletions, staged, unstaged, untracked — matching raw Git for the same fixture. | snapshot composer | observed red: `rg -n -e 'insertions' -e 'deletions' internal/diff` exited 1 | No aggregate derivation exists; each count is asserted against raw Git so a fabricated or visible-row-derived count fails. |
| CD2 | 1 | The `files[N]{status,path,kind}` inventory includes untracked paths with `kind` naming every non-regular entry and dangling symlink (blank for regular files), and the checkout table distinguishes index from worktree state per `git.Facts` XY entries. | facts composition | observed red: `rg -n -e 'ParsePorcelainZ' -e 'git\.Facts' internal/diff` exited 1 | `internal/diff` consumes no checkout owner today, and the kind column makes an omitted body or misclassified special file observable in the row itself. |
| CD3 | 1 | The snapshot ends with a whitespace-check result derived from `git diff --check`, including the clean case. | whitespace probe | observed red: `rg -n 'diff --check' internal/diff` exited 1 | The check exists nowhere in the package; a fixture with a trailing-whitespace hunk must flip the row. |
| CD4 | 2 | `--full` returns the snapshot plus the log and an exact patch whose untracked regular-file bodies append in path-sorted order as the exact per-file bytes of `git diff --no-index -- /dev/null <path>` (mode headers and binary form included), with non-regular untracked entries named instead of read. | full-body assembly | observed red: `rg -n 'untracked' internal/diff` exited 1 | Untracked content is silently absent today, and pinning the `--no-index` bytes per file gives the appended body a raw-Git oracle instead of an arbitrary concatenation. |
| CD5 | 3 | Content-sensitive identity — HEAD commit id, raw-index byte digest, porcelain byte-stream digest, and per-path digests of reported dirty and untracked entries covering each kind's patch-observable state (content and mode, gitlink HEAD, symlink target) — is captured around the reads per dimension; movement in any one retries once, then returns a structured drift error, and no spliced snapshot is ever emitted. | identity capture | observed red: `rg -n 'drift' internal/diff` exited 1 | One injected mutation per dimension distinguishes a torn snapshot from a coherent one, and the dirty-file rewrite, chmod, and submodule-HEAD cases — invisible to HEAD, index, and porcelain alike — only the per-path digests can catch. |
| CD6 | 1,2 | Per hostile-state class — committed, staged, unstaged, untracked, nested new directory, rename, deletion, binary, hostile filename, clean, mid-read drift — the response agrees with raw Git captured in the same fixture, with control-bearing paths refused per the existing TOON rules. | raw-Git paired fixture matrix | observed red: `rg -n 'func TestCommand' internal/diff` exited 1 | The package has no command-level test at all, so every class starts red and each carries its own mutation; the nested-directory class catches a top-level-only enumeration that collapses `dir/`. |
| CD7 | 4 | Bounded orientation appends the `bench diff --full` action, drift appends the exact same invocation, a bounded `--commit` orientation appends `bench diff --full --commit <sha>` with the sha carried, and clean or complete results append the honest empty help block. | help derivation | not TDD-able until `axi-spec-build-complete` lands the action owner; production `help[` emitters are currently zero | Deriving from the snapshot facts means a removed action, a guessed value, or a repeated already-satisfied invocation each fails the focused derivation test. |

### Ticket derivation

Every mapped row becomes ticket acceptance with `(covers <row>)`, atomic `Closure:` tokens, and a subject mutation under the approved fence. CD6 may split per fixture-class batch by repeating its covers ID. Only an unforeseen local behavior may use `(covers local)`.

| row | tracer acceptance and atomic facts | approved fence | subject mutation | independent owner and public operation |
|---|---|---|---|---|
| CD1 | render revision and aggregate rows / commits / files / insertions / deletions / staged / unstaged / untracked | `internal/diff` | derive one count from visible rows instead of raw Git | snapshot tests; run `bench diff` on paired fixtures and compare each count to raw Git |
| CD2 | include untracked paths, kinds, and index/worktree distinction / inventory / kind per special entry / blank kind for regular / XY split | `internal/diff` | collapse staged and unstaged into one state or drop a special entry's kind | facts-composition tests; stage, modify, and add regular and non-regular entries in one fixture and require the split and kinds |
| CD3 | report the whitespace result / dirty case / clean case | `internal/diff` | hardcode the clean result | whitespace tests; fixtures with and without a whitespace defect |
| CD4 | append untracked new-file patches and name non-regular entries / per-file `--no-index` bytes / path-sorted order / nested-directory expansion / mode header / binary form / FIFO refusal / log retained | `internal/diff` | omit one untracked body, drop a file below a new directory, or reorder the appended patches | full-mode tests; compare each appended body to its captured `git diff --no-index` bytes, including files under a new nested directory |
| CD5 | capture identity per dimension, retry once, refuse on drift / HEAD / index / porcelain / dirty-file content / dirty-file mode / dirty-gitlink HEAD / retry / structured error / same-invocation action | `internal/diff` | suppress one dimension's capture at a time and emit the spliced snapshot | drift tests; inject one mid-read mutation per dimension through the test seam, including rewriting an already-dirty file's bytes, chmod-ing a dirty file, and moving a dirty submodule's HEAD without changing porcelain output |
| CD6 | agree with raw Git per hostile class / each of the eleven classes | `internal/diff` | mutate the composer per class (drop a rename, misclassify binary, collapse a nested directory, render a control path) | paired fixture matrix; one real repository per class |
| CD7 | derive help from snapshot facts / bounded / drift / commit / honest empty | `internal/diff` | advertise `--full` on an already-complete response | help derivation tests; drive each response mode and require the exact action set |

### Edge inventory

- Error path — existing base-resolution, log, and diff error kinds are preserved; CD5 adds the drift refusal.
- Empty or absent input — the clean checkout renders definitive zero aggregates and empty inventory (CD6's clean class); outside-repo keeps `toon.NotInRepo`.
- Boundary values — zero, one, and many files; a diff whose only change is untracked; a base equal to HEAD.
- Malformed input — argv handling keeps the existing grammar; `--commit` with an unresolvable or parentless sha keeps its structured errors.
- Interrupted or partial state — CD5 refuses partial snapshots; a rebase-in-progress checkout renders whatever raw Git reports rather than inventing state.
- Post-landing `--commit` view — live-checkout facts are omitted, never zeroed or spliced from the unrelated live tree; CD6's committed class runs both the default and `--commit` modes and asserts the omission.
- Re-run idempotency — repeated invocations on a still tree return identical bytes.
- Process-boundary lifecycle — fixtures run the real command against real repositories; no in-memory state survives between invocations.
- Hostile environment — spaces/globs in paths, control-bearing names refused per `toon.Representable`, deep cwd, and detached HEAD are all fixture classes.
- Command self-observation — the snapshot never writes; identity capture proves reads cannot observe their own effects.
- Special files and dangling symlinks — CD4 names FIFOs and other non-regular untracked entries instead of reading them; CD2's `kind` column carries `dangling-symlink` for a broken link, so the class is a mapped row, not an unstated state.

**Won't handle:** landed-history search, blame, and cross-branch comparison are separate capabilities; `bench diff` stays scoped to the review-base and named-commit views.

## Out of scope

- The spec-build family envelope and the action owner it lands — `axi-spec-build-complete` (prerequisite, ~20 edits, 1 promotion gate).
- Contextual disclosure on the remaining query surfaces and conformance closure — `axi-query-disclosure` (~14 edits, 1 promotion gate).
- A `bench git` namespace, a second porcelain parser, or any new Git-facing command family.
- `--fields`, a legacy mode, a dual renderer, or changing `--commit`'s post-landing meaning.

## Ownership fences

- Snapshot owner: `internal/diff/`.
- Checkout facts: `internal/git/` (additive helpers only; existing consumers' behavior is unchanged).
- Registration and conformance: `cmd/bench/main.go`, `internal/conformance/`, `projects/benchkit.md`.
- `internal/axi` is a consumed input; the spec-build family and every other command remain unchanged.
