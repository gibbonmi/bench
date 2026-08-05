# Exact prospective landing

Status: implemented

Decision source: named reviewed artifact `decisions/parallel-session-landings.md` at `fd73811` (first ordered scope in decision #10)

## Problem

An ordinary path-scoped `bench commit` cannot land while another session has unrelated dirty or untracked files in the same checkout. It refuses those files because it gates the complete working-tree snapshot and only later stages the named paths. Removing that refusal without changing the gate subject would be unsound: the gate would grade one tree while Git published another.

This turns checkout cleanliness into a repository-wide writer lock. A session with an attributed, independently reviewable change must either wait, mix another author's files into its commit, or manually move foreign work out of the checkout. The 2026-08-03 parallel-session friction capture records this exact refusal during a live build.

## Solution

Introduce one exact prospective landing substrate. It starts from an explicit expected destination base, overlays only the caller-attributed named-path content and any command-owned lifecycle transition, and produces an immutable Git tree. The existing prospective gate grades that complete tree with the authoritative whole-project oracle. A commit whose sole parent is the expected base is then published only by compare-and-swap of the destination ref.

Ordinary `bench commit` adopts this substrate without changing its argument grammar or path attribution rules. Unrelated staged, unstaged, untracked, and ignored checkout content does not enter the prospective tree, does not block the landing, and is not changed by a refusal or a successful landing. A destination-tip movement loses the compare-and-swap and requires the caller to run the command again against the new base.

## User stories

1. **A path owner can compose exactly the change they named.** Starting from the destination ref's resolved commit, `bench commit` overlays the current working-copy state of each literal named path, including additions, modifications, deletions, staged deletions, both halves of a rename, and all changed descendants of a named directory. A directory widens only on whole path segments. Current index entries and working-copy bytes outside that attribution remain absent from the composition.
   `Line: gpt-5.6-terra / medium.` The path grammar is known, but moving it from a cleanliness predicate to immutable tree construction changes a correctness-bearing Git seam.

2. **The oracle grades the exact tree that can land.** The command sends the immutable prospective tree through the existing prospective gate, whose subject identity includes the complete derived oracle closure and whose prospective path runs the whole dev gate. Red, unavailable subject identity, unavailable oracle identity, or incomplete gate outcome refuses publication and leaves the destination ref and checkout state unchanged.
   `Line: gpt-5.6-terra / medium.` The existing gate seam is strong, but the new composition-to-oracle attachment decides whether green means shippable.

3. **A spec completion is inside the graded composition.** With `--spec`, staged-to-implemented is a command-owned transition applied while composing the prospective tree. The gate sees the implemented spec bytes, the published commit contains those same bytes, and a red or lost compare-and-swap leaves the checkout's spec staged. Invalid or already-implemented specs still fail before a gate run.
   `Line: gpt-5.6-terra / medium.` This closes the known gate-then-flip exception and must preserve the lifecycle's fail-fast behavior.

4. **Only the command whose expected base still owns the destination may publish.** The landing commit has the prospective tree and exactly one parent, the expected base. Publication compare-and-swaps the resolved destination ref from that base to that commit. If another session advances the ref after composition or during the gate, the command refuses without rebasing, publishing, consuming the caller's changes, or disturbing the winning commit.
   `Line: gpt-5.6-terra / medium.` The observable contract is compact, but deterministic race coverage and fail-closed ancestry are correctness-critical.

5. **Successful publication preserves concurrent checkout work.** After the compare-and-swap wins, the invoking checkout remains attached to the same destination, every named path is clean against the landed commit, and the optional spec file reflects its implemented bytes. Unnamed staged entries retain their staged bytes; unnamed unstaged and untracked files retain their bytes and classification; ignored files remain ignored. No real-index mutation is used to construct or gate the prospective subject.
   `Line: gpt-5.6-terra / medium.` This is the user-visible concurrency benefit and needs end-to-end Git-state assertions beyond a unit tree hash.

## Implementation decisions

**One deep landing owner composes, authorizes, and publishes.** Add one internal owner for exact prospective landings rather than teaching each caller its own throwaway-index, gate, commit-tree, and update-ref sequence. Its request fixes the repository root, resolved destination ref, expected base commit, commit message, literal attributed paths, and optional lifecycle-owned transition. Its result reports the base, tree, and published commit identity. Reversible type and file names remain implementation discretion.

**The expected base and destination ref are resolved once.** `bench commit` resolves the current checkout's destination ref and commit before reading attributed content. A branch checkout targets that branch ref; a detached checkout targets `HEAD` with the same expected-old requirement. Missing, unborn, ambiguous, or non-commit bases refuse before gating. The composer never silently substitutes a later `HEAD`.

**Composition uses an isolated Git index and the caller's named working-copy content.** The index is seeded from the expected base tree, and literal named paths are overlaid with the same semantics the current command applies before `git commit`: current working-copy bytes win over earlier staged bytes, absent named paths record deletions, and a named deleted directory covers the children tracked beneath it at the expected base. The real index supplies validation for staged deletion/rename shapes but is not the tree being written. Unnamed real-index entries never enter the prospective tree.

**Emptiness belongs to the prospective composition.** The command compares the composed tree with the expected base tree before gating. Equality returns the existing `nothing to commit` refusal even when the real index contains unrelated staged work; ambient index state can neither turn an empty named delta into a commit nor suppress a real named delta.

**Attribution is exact and fail closed.** Every named argument must resolve inside the repository and be attributable as a file, a directory at either the working copy or expected base, or a deletion present at the expected base. Literal pathspecs preserve spaces, glob characters, and dash-led names. Prefix-sharing siblings are excluded. A duplicate name is normalized without widening. A special file named directly or discovered beneath a named directory is rejected before reading; a symlink is recorded as the symlink Git entry and is never traversed as a directory. The repository root is not accepted as an implicit all-changes selector; FT98 owns that explicit capability. An absent path in the working copy, real index, and expected base remains a pre-gate error.

**Lifecycle transitions are composition inputs, not post-gate edits.** The initial transition is the existing staged-to-implemented spec operation. It runs only in the isolated composition, after the current fail-fast staged check, and its resulting bytes join the tree before authorization. The landing owner accepts no arbitrary caller-authored patch or unsealed tree mutation between gate and publish.

**The prospective gate remains the sole oracle.** Authorization uses the existing exact-tree prospective gate entry. It may reuse retained evidence only when both tree and oracle identities match exactly; prospective component narrowing or a diff-scoped verdict cannot authorize publication. The owner verifies that the authorized tree still equals the tree placed in the landing commit.

**Publication is a Git compare-and-swap, not `git commit` over ambient state.** After green, the owner creates a commit directly from the authorized tree with the expected base as its sole parent, then updates the destination ref only if it still equals that base. Losing the compare-and-swap leaves the created object unreachable and returns a recomposition-required refusal. It does not retry, merge, rebase, or advance any project-green ref.

**Checkout reconciliation is path-scoped and post-publication.** A winning command updates only the named real-index entries and, for a lifecycle transition, the transitioned working-copy path so those paths match the landed tree. Ordinary named working-copy bytes already match by construction. Unnamed index and working-copy entries must compare byte-for-byte afterward; the contract tests own those fingerprints rather than adding a second production inventory of Git state. Reconciliation is idempotent against the just-published commit; a reconciliation error is reported as a landed-but-checkout-incomplete operational failure rather than falsely claiming the ref did not advance.

**The ordinary command stays operational, not AXI-query shaped.** `bench commit` keeps its existing arguments, stdout/stderr posture, exit codes, and branch-agnostic role. The behavioral change is that foreign dirty paths are preserved instead of refused. Errors name whether failure occurred during attribution, prospective authorization, destination compare-and-swap, or post-publication reconciliation, and never print raw receipt or internal object schemas.

**The first implementation is wide across two ownership fences.** The landing/Git owner and its focused tests form one fence; the `bench commit` adapter and public runtime contracts form the other. Build-time ticketing follows `craft-tickets`: the reusable primitive lands green before the command adopts it, and no delegate crosses both fences.

## Testing decisions

- **External behavior a good test exercises.** Invoke the real `bench commit` binary in a temporary Git repository containing attributed changes plus foreign staged, unstaged, untracked, and ignored content. Make the gate inspect the materialized prospective checkout, then assert the published commit, destination ref, real index, and working-copy bytes.
- **Primary seam.** The existing runtime commit contract fixture is the highest seam: it already drives hostile path grammar, deletions, directories, spec flipping, gate red, and exit codes through the public command. Extend that family so its old "unexplained file blocks" control becomes preservation-and-success coverage.
- **Lower seam.** The exact landing owner receives focused tests with a controllable gate owner and ref updater. This second seam is justified because a destination movement between gate completion and publication cannot be triggered deterministically at the public CLI seam without a timing race. The lower tests inject that movement at the exact compare-and-swap junction; the public suite still proves one end-to-end concurrent winner/loser journey.
- **Prior art.** Prospective tree materialization and exact tree/oracle evidence live in the gate package. Spec-build promotion supplies the existing `commit-tree` plus expected-old `update-ref` pattern. Throwaway-index construction and literal named-path grammar live in the Git and commit packages. Tests reuse those owners rather than cloning fixture harnesses.
- **Gate seam.** `bench gate`'s runtime/behavior contracts observe the public command and its Git state. Gate-package tests remain the owner of exact prospective evidence and full-gate selection. No new gate policy is introduced.
- **Central mutation probe.** After the build is green, replace the compare-and-swap update with an unconditional destination update. The injected concurrent-advance test must go red because both callers publish or the winner is overwritten.

### Seam diagrams

    trigger: bench commit -m <message> [--spec <slug>] <path>...
        │
        ▼
    expected destination + literal named working-copy paths
        │
        ▼
    [ exact prospective landing owner ]
        │ compose from expected base; apply owned spec transition
        ├──────────────▶ immutable tree ──▶ [ prospective whole-project gate ]
        │                                      │ exact green only
        ▼                                      ▼
    commit(tree, parent=base) ──CAS────────▶ destination ref + scoped checkout reconciliation
        ◀ tests attach here: real CLI fixture observes gate checkout, commit, ref, index, and files

    trigger: gate returns while another writer advances the destination
        │
        ▼
    expected base + authorized tree ──▶ [ publish compare-and-swap ] ──▶ published commit or drift refusal
                                             ▲
                                             └ tests attach here: injected ref updater moves the ref
                                               immediately before the expected-old update

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1, 5 | A named file lands while an unrelated tracked modification and an unrelated untracked file remain byte-identical and uncommitted | public runtime commit fixture | observed 2026-08-03 and pinned by current `testCommitUnexplainedBlocks`: the public command exits 1 with `outside the named set` before the gate | The cheapest wrong implementation retains the global cleanliness refusal; this journey cannot pass until foreign dirt is excluded from the composition rather than rejected or committed |
| 1, 5 | Unrelated staged content remains staged but absent from the landed commit; named staged-plus-unstaged content lands its current working-copy bytes | public runtime commit fixture with index assertions | to be observed at build time — write the success and exact index/tree assertions before changing the command | A composer seeded from the real index silently pulls foreign staged bytes into the tree or lands stale staged bytes for the named path |
| 1, 5 | A clean named path returns `nothing to commit` while unrelated staged content remains staged and no gate or publication occurs | public runtime commit fixture with gate tally, ref, and index assertions | to be observed at build time — the existing empty control has no foreign staged entry | The cheap port retains the real-index `nothingStaged` predicate; foreign staged work makes that predicate false and permits an empty prospective commit |
| 1 | Literal paths preserve spaces, glob characters, dash-led names, additions, deletions, staged deletions, rename halves, named-directory descendants, deleted-directory descendants, and whole-segment sibling exclusion | existing public runtime commit family, one assertion per listed class | already covered for the current command grammar; the same tests must be rerun through prospective composition, with the two outside-path refusal expectations changed only where story 5 requires preservation | Enumerating the classes prevents a new composer from implementing only regular-file modification while claiming current path compatibility |
| 1 | Duplicate named paths do not widen attribution or produce a duplicate change | landing owner plus public CLI | to be observed at build time | Deduplication by literal repository path catches implementations that treat repeated directory/file arguments as a second ambient staging pass |
| 1 | Unknown, repository-escaping, unreadable, repository-root, and directly or transitively discovered special-file inputs refuse before gate execution; symlinks are recorded without traversal | landing owner path table plus gate-run tally | to be observed at build time per class | A fail-open reader can block on a FIFO, escape or widen to the repository, or follow bytes the reviewer did not name; the zero gate tally proves refusal ordering |
| 2 | The gate runs against the exact prospective tree: it sees the named change and base bytes, and cannot see foreign staged, unstaged, untracked, or ignored content | public runtime fixture with a gate script that asserts its checkout | to be observed at build time | The composition degenerate is a command that preserves foreign work in the commit but still gates the ambient checkout; inspecting the gate checkout makes that mismatch red |
| 2 | Gate red, unavailable subject/oracle identity, and incomplete outcome each leave destination ref, real index, and every working-copy byte unchanged | landing owner failure table plus public red-gate control | current public red-gate control covers ref refusal only; full state assertions are to be observed at build time for all three outcome classes | An always-green stub or a pre-gate real-index mutation can otherwise satisfy the happy path while violating commit-on-green |
| 2 | Prospective authorization runs the whole gate or reuses evidence only for the exact tree and oracle; component/reduced evidence cannot authorize a different prospective subject | existing gate prospective tests plus landing-owner integration | already covered by gate prospective selection and evidence-key tests; integration row must show the landing owner calls that surface | This catches a cheaper diff-scoped implementation that weakens the oracle while keeping commit tests green |
| 3 | With `--spec`, the gate observes `Status: implemented`, and the identical implemented bytes land in the commit and invoking checkout | public runtime spec fixture with gate-side file assertion | to be observed at build time — the current command gates `Status: staged` and flips afterward | A post-gate flip produces a green tree different from the published tree; comparing all three byte sources catches it |
| 3 | Invalid or implemented `--spec` fails before the gate; gate red or compare-and-swap loss leaves a valid spec staged | public spec fixture plus gate tally and injected ref movement | invalid/already-implemented fail-fast is already covered; red and drift preservation are to be observed at build time | This prevents lifecycle validation from moving late and prevents an unpublished transition from leaking into the checkout |
| 4 | The landing commit's tree equals the authorized tree and its sole parent equals the captured expected base | landing owner test with real Git objects | to be observed at build time | Creating a normal ambient commit after gating can silently adopt a later parent or foreign index entry; object identity exposes either |
| 4 | A destination advance after green causes expected-old publication to fail; the winner remains the tip and the losing work remains available for a rerun | injected ref-updater test and one two-process public journey | to be observed at build time | The unconditional-update mutation overwrites the winner, while an auto-rebase publishes a tree that never received this gate verdict |
| 4 | A detached checkout publishes by expected-old update of detached `HEAD`, and a concurrent detached-HEAD movement refuses identically to branch-ref drift | landing owner test with real detached checkout | to be observed at build time | Preserving the command's branch-agnostic contract requires the destination resolver and compare-and-swap to work without a symbolic branch ref |
| 4, 5 | Rerunning the losing command against the new tip recomposes, gates, and lands a second commit containing both independently attributed changes without mixed authorship | two-process public journey | to be observed at build time | Testing only the refusal could strand safe work; the rerun proves the ordinary user can make progress on the primitive without scope-two run machinery |
| 5 | After success, named paths are clean; every unnamed staged entry retains its index blob and every unnamed working-copy/untracked/ignored path retains its bytes and classification | public runtime fixture with pre/post index blob and filesystem fingerprints | to be observed at build time | A whole-index reset or checkout-wide read-tree makes the commit correct while deleting or unstaging concurrent work; the fingerprints catch that cheapest destructive implementation |
| edge of 4, 5 | A reconciliation failure after a successful compare-and-swap reports that the commit landed and does not roll the destination backward | landing owner fault injection after ref update | to be observed at build time | Treating the operation as wholly failed invites a retry that creates a duplicate commit; rolling back can overwrite a later writer |

The central degenerate is the current sequence with its cleanliness check removed: gate the ambient checkout, then stage and commit only named paths. It preserves foreign files and may even produce the right commit, but the gate grades the wrong tree. The gate-checkout inspection row goes red on it. The concurrency degenerate is `commit-tree` followed by an unconditional ref update; all serial rows pass, while the injected race row reds.

### Edge inventory

- **Error path** — covered by rows for attribution failure, gate red/unavailable/incomplete, compare-and-swap loss, and post-publication reconciliation failure.
- **Empty or absent input** — existing usage coverage keeps an empty positional and no-path invocation at exit 2; a named path absent from working copy, index, and expected base remains a pre-gate operational error; a named clean path remains `nothing to commit`.
- **Boundary values** — one file, multiple files, a directory with two children, a deleted directory with two tracked children, duplicate names, and a repository-root path argument are covered. The repository-root argument refuses; it never becomes the implicit `--all` flag that FT98 separately owns.
- **Malformed and hostile input** — covered per class: spaces, glob characters, a leading dash after `--`, repository escape, special files, unreadable paths, symlinks, and directory-prefix collisions. Control bytes in Git paths are exercised at the Git-object seam and must produce a bounded escaped diagnostic or refusal, never an unprintable progress line.
- **Interrupted or partial state** — pre-publication failures and signals leave no ref movement and no real-index mutation. After the ref compare-and-swap, the commit is authoritative; a reconciliation failure reports landed-but-incomplete and never rolls the ref backward. Process death in that narrow post-CAS window may leave named checkout entries visibly dirty against an already-landed commit; automatic crash resumption is **Won't handle** here because durable operation journals and recovery authority belong to the downstream multi-coordinator lifecycle spec, while Git preserves both the landed commit and working bytes.
- **Re-run idempotency** — rerunning after a pre-publication refusal recomposes from the current destination; rerunning after a successful landing finds no named delta and returns `nothing to commit`; rerunning a CAS loser after the winner lands is the progress path covered above.
- **Process-boundary lifecycle** — the two-process journey resolves the new destination in a fresh process before recomposition. No in-memory plan or gate result may substitute for re-reading the ref.
- **Hostile environment** — a held real `index.lock` may prevent post-publication reconciliation but cannot affect isolated composition or gate identity; a held gate lock follows the gate owner's existing wait/reuse contract; missing Git or prospective preparation fails before publication. CWD below the repository root keeps current root-relative literal path behavior.
- **Ignored content** — **Won't handle as attributed input**: current `bench commit` semantics do not land ignored files, and making ignored content committable is a separate CLI decision. Ignored foreign content is preserved and excluded from the composition.
- **Submodules** — **Won't handle as traversed directories**: a gitlink is an atomic named entry, not a directory whose nested dirty state can be represented in the superproject tree. Naming the gitlink may land a changed gitlink identity; nested repository content remains outside the superproject composition.

## Out of scope

- **Multi-coordinator spec-build lifecycle.** Run revisions, run-owned planning subjects, CLI-derived AXI receipt helpers, evidence-preserving recomposition, checkout-independent lifecycle operations, recovery authority, and concurrent-run status form the second ordered spec and literally depend on this substrate. Derived breadth: at least six ownership clusters, `30+ edits, 7+ gate runs`.
- **Cross-composition gate-evidence policy beyond exact tree-and-oracle reuse.** Sharing component/check evidence among independently composed trees is decision #12 of the downstream scope. This build retains the existing prospective full-gate rule. `6 edits, 3 gate runs` once the policy is decided.
- **Recoverable set-aside, mutation revert, and `bench commit --all`.** FT98 owns that preserve-then-discard vocabulary. Exact prospective composition removes set-aside as a prerequisite for ordinary path-scoped landing but does not supply the other capability's recovery plans or explicit all-path selector. `12 edits, 3 gate runs`.
- **A generic light-path worktree landing command.** FT169 owns preparing and retiring assignment worktrees and stale-base transfer ergonomics. This substrate can become its publication primitive later; it does not add worktree porcelain. `10 edits, 3 gate runs`.
- **Weakening `bench prep-release`.** The ship tier continues to require a current exact-tree full-green dev verdict. `0 edits, 0 gate runs`.
- **Advancing project-green evidence from ordinary `bench commit`.** This command remains a gated commit, not the reviewed spec-build promotion transition. Adding a second project-green publisher would change lifecycle authority and belongs to a separate reviewed decision. `5 edits, 2 gate runs`.
- **Harness-specific landing behavior or adapters.** The shared compiled `bench` command remains the one implementation reached by every supported harness; this scope does not change the closed portability decision or add a harness-only publish path. `0 edits, 0 gate runs`.
