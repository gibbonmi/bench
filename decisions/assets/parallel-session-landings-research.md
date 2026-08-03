# Parallel-session landing research

Local-code research performed 2026-08-03 against `main`. This asset answers
decision tickets #1–#3 in `decisions/parallel-session-landings.md`.

## Landing and gate subjects

- Ordinary `bench commit` refuses every dirty or untracked path outside its
  named set, then gates the complete working-tree snapshot before staging only
  the named paths (`internal/commit/commit.go:96-132`). The snapshot is a Git
  tree built through a throwaway index (`internal/git/tree.go:11-43`).
- Gate evidence binds the exact Git tree and the oracle closure. The gate
  rebuilds that subject under its execution lock and after the run, refusing
  drift at either boundary (`internal/gate/gate.go:364-366,443-446`).
- `gate.ExecuteTree` already materializes and gates an unpublished Git tree in
  a detached checkout (`internal/gate/engine.go:98-112`). Prospective subjects
  deliberately run the full gate rather than component narrowing
  (`internal/gate/gate.go:377-390`). There is no first-class subject consisting
  only of a patch plus separately supplied gate inputs.
- Spec-build promotion constructs an unpublished prospective tree from the
  retained candidate plus the staged-to-implemented spec transition, gates that
  exact tree, creates the commit from it, and compare-and-swap advances the
  branch (`internal/specbuild/state.go:264-299`;
  `internal/specbuild/assign.go:357-417`).
- Assignment requires the selected ticket's working copy and index to match
  the digest of `HEAD:<ticket>`; unrelated dirty tickets are permitted
  (`internal/specbuild/precondition.go:171-184`;
  `internal/specbuild/start_test.go:589-646`).
- One current exception needs an explicit disposition: ordinary
  `bench commit --spec` gates the staged spec, then flips it to implemented
  before staging, so its published tree is not byte-identical to the gated tree
  (`internal/commit/commit.go:115-132`).

## Tip drift, fences, and checkout identity

- Every nonterminal run records a base. Any recognized branch advance where
  `run.Base != current tip` produces the same recomposition refusal for every
  mutator except abandon; the predicate does not inspect changed paths or
  assignment fences (`internal/specbuild/precondition.go:68-106`).
- Empty runs have one narrow fast-forward path. Once work exists, promote
  replays the candidate on the moved tip, compare-and-swap moves the candidate
  ref, clears the entire held review and promotion evidence, and persists the
  new base (`internal/specbuild/recompose.go:35-88`).
- Run state already retains ticket digests, charged rows, ownership fences,
  assumptions, checkpoint tree and receipt digests, integration candidate,
  candidate tip, and review candidate/digest. No current owner compares fences
  across runs or proves that a landed change is gate-input-disjoint
  (`internal/specbuild/state.go:22-91`).
- The observed assignment-worktree refusal is not owned by an explicit
  `primary checkout` predicate. CLI routing constructs the service from the
  invocation checkout's Git root, while ticket/spec paths must belong to that
  service root (`cmd/bench/specbuild.go:21-31`;
  `internal/specbuild/precondition.go:163-184`;
  `internal/specbuild/state.go:264-268`). Primary-checkout-only behavior is
  therefore a consequence of root/spec identity in the observed case.

## State, evidence, and promote coupling

- Each spec slug has one exclusive flock at
  `<git-common-dir>/bench/specbuild/<sha256(slug)>.json.lock`; different slugs
  use different locks (`internal/specbuild/state.go:142-158,256-261`). State
  replacement is durable temp-write, fsync, rename, installed-file fsync, and
  directory fsync (`internal/specbuild/state.go:217-253`).
- Per-slug transitions still touch shared branch, candidate, green, Git object,
  and worktree surfaces. The transitions are journaled multi-step operations,
  not one atomic state write; compare-and-swap protects refs but makes one
  concurrent winner's base movement visible as drift to the others
  (`internal/specbuild/integrate.go:28-57,119-191,306-318`;
  `internal/specbuild/assign.go:378-417`).
- A checkpoint receipt supplies version, run, assignment, base, tree, ticket
  digest, charged rows, checks, coordinator probe, ownership, and assumption
  digests. Validation independently recomputes the tree, ticket, assumptions,
  and fence containment (`internal/specbuild/checkpoint.go:93-123,237-346`). A
  review receipt supplies the exact candidate, ordered Standards/Spec/Coverage
  axes, and findings (`internal/specbuild/lifecycle.go:328-426`).
- The public CLI forwards receipt paths but emits neither receipt body nor
  skeleton (`cmd/bench/specbuild.go:61-129`). `bench tree-hash` and ticket
  parsing already own two component derivations (`cmd/bench/main.go:298-306`;
  `internal/specbuild/assign.go:98-151`).
- Recomposition and publication are separate internal branches but one public
  `promote` operation. A moved-tip call runs recomposition and returns without
  gating or publishing; a current-tip call validates review/evidence, gates,
  commits, advances the branch and green ref, and terminates the run
  (`internal/specbuild/assign.go:329-417`). Recomposition is not state-only: it
  requires exact-green bootstrap for the moved base, conflict-free replay,
  candidate-ref compare-and-swap, and review invalidation
  (`internal/specbuild/recompose.go:9-32,77-90`).

## Constraints carried forward

- Keep the exact materialized Git tree and oracle closure as the gate subject;
  a path list alone does not prove the published composition.
- Keep prospective compositions on the full gate unless a separately reviewed
  design changes that oracle claim.
- Fail closed when tree/oracle identity, fence overlap, evidence derivation, or
  compare-and-swap ancestry cannot be established.
- Keep `bench prep-release`'s exact-tree full-green prerequisite unchanged.
