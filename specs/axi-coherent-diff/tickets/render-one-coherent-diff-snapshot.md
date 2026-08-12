# Render one coherent diff snapshot

Blocked by: none
Writes: `internal/diff/`, `internal/git/`, `internal/axi/`, `cmd/bench/main.go`, `internal/conformance/`, `projects/benchkit.md`, `CHANGELOG.md`

Closure: CD1/revision, CD1/commits, CD1/files, CD1/insertions, CD1/deletions, CD1/staged, CD1/unstaged, CD1/untracked, CD2/inventory, CD2/special-kind, CD2/regular-kind, CD2/xy-split, CD2/no-own-porcelain, CD2/legacy-facts-preserved, CD3/dirty, CD3/clean, CD4/untracked-patches, CD4/path-order, CD4/nested-expansion, CD4/mode-header, CD4/binary-form, CD4/nonregular-refusal, CD4/log-retained, CD5/head, CD5/default-tip, CD5/recorded-base, CD5/index, CD5/porcelain, CD5/content, CD5/mode, CD5/gitlink, CD5/retry, CD5/structured-refusal, CD5/same-invocation-help, CD5/single-resolution, CD6/committed, CD6/staged, CD6/unstaged, CD6/untracked, CD6/nested-directory, CD6/rename, CD6/deletion, CD6/binary, CD6/hostile-filename, CD6/clean, CD6/mid-read-drift, CD6/detached-head, CD6/deep-cwd, CD6/idempotency, CD6/base-equals-head, CD7/bounded, CD7/drift, CD7/commit, CD7/honest-empty, CD8/fixed-arguments, CD8/honest-empty, CD8/block-rendering, CD8/no-prose-command, CD9/bare, CD9/full, CD9/commit, CD9/commit-full, CD9/preserved-regions, CD9/named-delta

## What to build

Replace the current partial `bench diff` response with the spec's one coherent snapshot while preserving the existing command name, argv grammar, exit taxonomy, tracked file rows, log rows, tracked patch bytes, and structured error kinds outside the named delta. Capture the reviewed pre-migration responses before changing output, then keep each candidate mode paired with its old response.

The production command must consume the existing `git.Facts` owner for branch and divergence plus an additive diff-specific all-files facts path for porcelain state; existing `git.Facts` callers and their collapsed-directory output remain byte-identical. It derives the revision, aggregate, inventory, checkout, and whitespace blocks from the same attempt; includes path-sorted raw-Git new-file patches for untracked regular files in `--full`; and retains immutable post-landing semantics for `--commit`. Live snapshots capture every patch-observable identity dimension around their reads, retry the whole attempt once on movement, and otherwise refuse with the drifted dimension named.

Introduce `internal/axi`'s typed action and `help[]` renderer only with its first production consumer. Bounded orientation and drift derive exact executable next invocations from live facts; complete or clean results emit the honest zero-row block. Update only the existing registration, conformance, profile, and changelog surfaces whose current claims are changed by this migration.

CD1-CD9 remain one ticket because a thinner landing changes the public response before either the paired compatibility oracle or every-response `help[]` contract is complete; that intermediate public shape is the specific stranded red the approved atomic migration forbids.

## Acceptance

- [ ] [CD1] (covers CD1) default mode renders the exact live `revision` and `aggregate` schemas, with every value checked against raw Git from the same fixture.
- [ ] [CD2] (covers CD2) default mode includes all-files untracked inventory, exact special-entry kinds, and an XY `checkout` table sourced through an additive `internal/git` facts path, while a static probe finds no `git status` invocation in `internal/diff` and a nested-directory fixture keeps existing `git.Facts` consumers byte-identical.
- [ ] [CD3] (covers CD3) live modes render the tracked-patch whitespace verdict and offense count for both clean and offending fixtures.
- [ ] [CD4] (covers CD4) `--full` retains its log and tracked patch bytes, then appends every untracked regular-file patch as path-sorted exact `git diff --no-index -- /dev/null <path>` bytes while naming rather than reading non-regular entries.
- [ ] [CD5] (covers CD5) each live attempt captures HEAD, default tip, consumed recorded base, raw index, raw porcelain, and per-path patch-observable identities before and after reads; one movement retries once and a second refuses with `snapshot drift`, the dimension, and the exact refused invocation.
- [ ] [CD6] (covers CD6) real production-command fixtures for committed, staged, unstaged, untracked, nested-directory, rename, deletion, binary, hostile-filename, clean, mid-read-drift, detached-HEAD, deep-cwd, and base-equals-HEAD states agree with raw Git, including control-path refusal, and repeated invocations on a still tree return identical bytes.
- [ ] [CD7] (covers CD7) the production command derives exact `help[]` actions for bounded live, bounded commit, and drift modes, and emits the honest empty block for clean or complete modes across fixtures with different shas and modes.
- [ ] [CD8] (covers CD8) the typed-action owner retains every known fixed argument, renders executable commands and one-clause reasons, and rejects removed actions, placeholders for known values, guessed unknowns, dropped fixed flags, and prose masquerading as commands.
- [ ] [CD9] (covers CD9) checked-in old-to-new pairs for bare, `--full`, `--commit`, and `--commit --full` differ only by the spec's named delta; all preserved byte regions, error kinds, exits, and argv behavior remain equal.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CD1/raw-aggregate | derive insertion or divergence counts from rendered file rows | real-repository snapshot test | construct a fixture whose row count differs from its raw-Git aggregate, run the production command, and require the literal raw-Git values |
| CD2/facts-owner | run or parse `git status` inside `internal/diff`, or change existing `git.Facts` to all-files semantics | static owner probe plus dual-consumer nested-directory fixture | inspect production sources for the forbidden invocation, then require distinct index/worktree cells and expanded paths from `bench diff` while the existing context consumer retains its collapsed directory row |
| CD3/whitespace | hardcode `clean=true` and zero offenses | paired whitespace fixtures | run clean and trailing-whitespace repositories and require their verdict cells to differ exactly as `git diff --check` |
| CD4/untracked-body | omit or reorder one nested untracked regular-file patch | full-mode raw-Git oracle | capture each per-file `--no-index` result, run `bench diff --full`, and require every exact body in sorted order |
| CD5/identity | suppress one before/after identity dimension | injected drift matrix | mutate each named dimension between reads, require one full retry, mutate it again, and require refusal without a snapshot |
| CD6/hostile-class | collapse one class, resolve from the process cwd instead of the repository root, or introduce stateful output across invocations | production-command fixture matrix | drive every enumerated real repository from root and deep cwd, including detached and base-equals-HEAD states, repeat a still-tree invocation, and compare each observable with raw Git or the existing TOON refusal |
| CD7/help-derivation | hardcode one help block in `internal/diff` | differing-mode production fixtures | run bounded live, bounded commit with distinct shas, drift, clean, and complete modes and require each exact action set |
| CD8/action-fidelity | replace a known argument with a placeholder | action-owner counterexample test | render the action, require the exact executable argv, and independently delete, guess, drop, and prose-wrap one owned fact at a time |
| CD9/preserved-byte | alter one tracked file row, log row, patch byte, error kind, exit, or argv result outside the named delta | old-to-new pairing oracle | compare every captured old response with its candidate by mode and require only the enumerated sections to differ |
