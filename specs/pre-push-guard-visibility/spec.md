# Pre-push guard visibility

Status: implemented

Decision source: compiled map at
`specs/pre-push-guard-visibility/decisions/pre-push-guard-visibility.md`, with its
one structured source, the research asset at
`specs/pre-push-guard-visibility/decisions/assets/pre-push-guard-visibility-research.md`.
Both were re-read against the current tree on 2026-07-31. The map's drift clause
names the pre-push template, `ClassifyPrePush`, the guards header parser, and
`bench link`'s conflict ordering as its invalidation triggers; all four were
re-verified and all four still hold. One source conflict was found and is recorded
under Implementation decisions: the map's #1 and #9 describe `installGitHook` as
the install path, but that function has no call site. Nothing was unreadable.

## Problem

An enforcement control that reports its own status reports things it never
verified. The installed pre-push hook's protected branch may be a baked guess and
no surface says so. A managed hook may be arbitrarily stale against the embedded
template — this repository is running one — and doctor, guards, and the
SessionStart banner all report it as fine or say nothing, because the entire
managed test is a marker substring. And `bench link`, the sanctioned repair,
aborts on a symlinked managed directory before it ever reaches the hook refresh,
so the repair is unavailable in exactly the repository that needs it.

Underneath the three faces sits a structural problem: three surfaces already
answer "how is the pre-push hook doing" from two different derivations, and adding
a report to each would make it four or five. The repository's own code standard
names that as the defect.

## Solution

One computation answers for the pre-push hook. It reports the hook's state, the
branch that hook will actually protect on the next push, whether that branch came
from live repository facts or from the token baked at install time, and whether
the installed bytes are current against the embedded template. `bench doctor` owns
it; `bench guards` and `bench status` render from it instead of deriving anything
themselves, which collapses the existing duplication rather than adding to it.

The record mirrors the hook's own resolution rather than recomputing what a fresh
install would bake, so it reports the guard that exists instead of the guard that
would be installed. A guessed branch is a green status field, because the runtime
path usually corrects it and a red would train readers to ignore reds. A stale
hook reds doctor's row and names its remedy, and `bench doctor --fix` gains the
repair — opt-in, never automatic. The remedy becomes trustworthy because `bench
link` and `bench upgrade` stop aborting on a symlink parent for a plan entry that
already matches the manifest and therefore needs no write at all.

## User stories

Stories 1–7 are **mid** tier and story 8 is **top**; the resolved ids are the
`claude` column of the profile's harness × tier table. A different harness rebinds
from that table rather than from these tokens.

1. As an operator, I want one computation to answer for the installed pre-push
   hook — its state, the branch it will protect, where that branch came from, and
   whether its bytes are current — so that every surface reporting on the hook
   reports the same verified facts. Line: opus / medium. The record's shape is the
   whole feature's seam and the existing classifier constrains it, but no surface
   behavior is novel.

2. As an operator, I want `bench doctor`'s pre-push row to name the branch it
   protects and whether that branch was resolved live or fell back to the baked
   token, staying green either way, so that a guarded repository can tell what its
   guard actually covers. Line: opus / medium. Small rendering change over story
   1's record, with the green-either-way posture the part that must not drift.

3. As an operator, I want a stale managed hook to red doctor's pre-push row and
   name its remedy, and `bench doctor --fix` to perform that repair, so that drift
   is visible and repairable without doctor ever healing itself. Line: opus /
   medium. The `--fix` path already compares the shim body against its source, so
   the extension is at a known seam, but the no-self-heal boundary is
   correctness-critical.

4. As an operator, I want `bench guards` and the SessionStart banner to render the
   branch, provenance, and currency fields from the shared record rather than
   their own parse, and `bench status` to keep its signal budget, so that the
   ambient surfaces stop disagreeing with the diagnostic one. Line: opus / medium.
   Replacing an independent parser with a shared record is mechanical, but the
   TOON contract and the banner's brevity both bind.

5. As an operator of a repository with a symlinked managed directory, I want
   `bench link` and `bench upgrade` to skip a plan entry that already matches the
   manifest before the symlink-parent conflict check, so that the sanctioned
   repair route works where the stale hook lives. Line: opus / medium. The change
   is a few lines in the plan loop, but it moves a conflict boundary that the
   lifecycle gate phase exercises heavily.

6. As an operator, I want `bench upgrade --check` to count the hook refresh it
   will in fact perform, so that the plan I approve matches the writes I get.
   Line: opus / medium. The count derives from plan entries while the hook is
   staged outside the plan; correcting it is small and local.

7. As a maintainer, I want the branch substitution to have exactly one live source
   and the hook path to be classified by file mode before it is read, so that
   "current" cannot mean something different from what the installer writes and a
   broken hook cannot read as an absent one. Line: opus / medium. Deleting the
   dead installer is mechanical; preserving its unique probe and closing the
   dangling-symlink hole are the parts that need care.

8. As a reader of `README.md`, I want the pre-push claim narrowed to what the
   checks support, so that the one place in the repository overstating an
   enforcement guarantee stops doing so. Line: fable / high. The profile routes
   doc authoring to the top tier under `craft-line`'s leverage override; this
   follows that cached routing rather than claiming an exception for it.

## Implementation decisions

**A source conflict the map did not see.** The map's #1 records that install time
bakes `git.ResolvedDefault`'s answer at `installGitHook`, and #9 names the
currency comparison as `strings.ReplaceAll(prePushTemplate, prePushBranchToken,
protectedBranch(root))`. `installGitHook` has no call site — the map's #4 noted
this as an unplanned finding but did not follow it through. The live installer is
the link transaction's own staging call, which resolves the branch with
`hookBranch`: a networked `ls-remote --symref` probe falling back to
`protectedBranch`. So the map's named formula would compare an installed hook
against a branch the installer never used.

**The resolution is to stop recomputing the branch for the comparison at all.**
Currency compares the installed bytes against the template substituted with *the
token already baked into that installed file*, recovered from the file itself.
This makes currency a pure template-drift check that is correct regardless of
which resolver installed the hook, needs no network, and still reds this
repository's dead-protocol hook. It is a narrower reading of #9's "exact bytes
against the substituted embedded template" than the formula the map spelled out,
and it is the reading that survives the tree. **Flagged for reviewer veto.** The
map's accepted cost — a binary comparison that cannot distinguish one release
behind from a dead protocol — is unchanged.

The install-time resolver therefore does not need to collapse into the reporting
path, and `hookBranch` stays as the single live installer resolver. What is
removed is the dead duplicate: `installGitHook` is deleted along with its second
copy of the substitution. Its one unique behavior, the `git remote set-head origin
--auto` probe that populates `origin/HEAD` locally, moves into the live install
path; its hook-directory creation, foreign-hook refusal, substitution, and write
all already exist there. The probe is best-effort and its failure is not an install
failure — it improves later resolution rather than gating it.

**The record mirrors the hook, not the installer.** The map's #1 establishes that
the installed hook re-resolves `origin/HEAD` on every push and reaches the baked
token only when that lookup is empty. The record therefore computes the effective
branch exactly the way the hook does — local `symbolic-ref` on
`refs/remotes/origin/HEAD`, else the token parsed out of the installed file — so
it reports the guard that exists rather than the guard a fresh install would
produce. Recomputing a resolver would answer a different question and would go
wrong whenever repository state changed after install.

Provenance is therefore two-valued and about the live lookup: the branch resolved
from `origin/HEAD`, or the lookup was empty and the baked token is in force. The
guessed case the map cares about is the second one, and the record additionally
reports whether that baked token is the bare literal fallback rather than
something a resolution produced. No network call is made anywhere in this path: it
runs in `bench doctor`, in `bench guards`, and in the SessionStart banner, all of
which must work offline and fast.

**The record's shape.** One exported function in `internal/adopt` returns the
complete record: the existing four-valued state, the resolved hook path, the
effective branch, its provenance, and currency. `ClassifyPrePush` becomes an
internal step of that computation and stops being an exported entry point, so no
caller can reach a partial answer — that unexported-ness is what makes "one
computation" checkable rather than merely intended. The four existing
`PrePushState` values keep their meanings exactly, as the map's discretion clause
requires; the new fields are fields, not states.

Currency is three-valued, because a comparison that cannot be made is not a
mismatch: current, stale, or not-applicable for a hook that is absent, foreign, or
diverted. Comparing a foreign hook's bytes against the kit template is
meaningless, and the foreign state already carries its own report.

**Classify by file mode before reading.** The hook path is `Lstat`-ed before it is
read. A dangling symlink there is **foreign**, not absent: something occupies the
path and it is not a working managed hook, and reporting absent would tell the
operator a fresh clone dropped the hook when in fact a broken one is installed.
The link transaction's foreign-hook refusal gets the same treatment, because its
current read-first check treats a dangling symlink as no-hook and would overwrite
it. Special files at the hook path are refused on their mode rather than read, so
the ambient banner cannot block on a FIFO.

**Rendering, not re-deriving.** `bench doctor` owns the computation.
`internal/guards` drops `prePushRow`'s independent hook read entirely and renders
the shared record; its generic header parser stays for `.bench/hooks/*.sh`, which
is a different question about a different file class. `bench status` renders the
same record and keeps its existing show-only-on-signal budget: a guessed branch is
green and therefore not a signal, so status gains no row for it, while a stale
hook is a signal and gets one. The banner surfaces the fields through `guards
--brief`, which is where the map's #8 puts them. Guards' `denies` column keeps the
install posture it carries today and the `wired` cell stays the constant `git`;
the new facts are additive.

**Repair.** `bench doctor` never writes. `bench doctor --fix` gains the hook repair
on the same terms as its existing shim repair: it rewrites a *managed* stale hook
to the freshly substituted template with the executable bit set, reports an
already-current hook as a no-op, and refuses a foreign hook without touching it.
It does not install an absent hook — that is `bench link`'s transaction, and
doctor's row already names it.

**Link convergence.** The plan loop skips an entry whose destination already
matches the bytes and mode this plan would write before it reaches the
symlink-parent check, because an entry that needs no write cannot conflict with
anything. The predicate is prospective, not historical: a destination that
matches the recorded manifest hash but not the incoming kit bytes needs a write,
so relink and upgrade always propagate changed kit content. Conflict semantics
are unchanged for every entry that does need a write: a new entry and a drifted
entry both still abort on a symlink parent, and those are two independent
behaviors that get two independent fixtures.

**Upgrade counts.** `bench upgrade` returns before `transactionalLink` when the
installed and linked versions match, so no hook refresh happens on that path and
none is counted. On a real upgrade the hook is always refreshed while being staged
outside the plan, so the count adds it when its prospective bytes differ from the
installed ones.

**No manifest row.** The hook stays out of the SHA-256 manifest, as `unlink.go`'s
recorded decision holds and the map's #9 confirms: the manifest addresses
repo-relative paths and the hook lives outside the worktree.

**No migration path** for hooks speaking the dead `--describe` protocol. They red
as stale like any other stale hook and `--fix` repairs them.

**Consumers are not in this build.** The map's #12 makes this a prefactor. The
doctor shim row, the doctor gate row, the `guards` `wired` cell, and the
`gitguard` `denyTable` duplication become their own roadmap rows that
`/bench-what-next` creates; this spec does not edit `ROADMAP.md`.

**The local stale hook is repaired after this spec's coverage lands**, per the
map's #13. It is the only live instance of the dead-protocol case and the currency
fixture is built from it.

## Testing decisions

- The primary seam is the exported hook-health record in `internal/adopt`, driven
  against real temporary git repositories with the hook planted in each state.
  Tests assert the whole record rather than individual predicates, so a surface
  cannot read a field the computation did not fill.
- Single ownership is tested structurally, not by inspection: `ClassifyPrePush`
  becomes unexported, so a second caller outside the package fails to compile and
  a second derivation inside it is a focused call-site assertion.
- Doctor, guards, and status rendering are black-box runtime contracts over
  fixture repositories, in the existing runtime and AXI gate phases. Guards is
  graded as TOON; doctor as its row text and exit code; status as its row budget
  and severity ladder.
- The no-network property is tested by shimming `git` on `PATH` with a recorder
  and asserting no remote-contacting subcommand is invoked, rather than by timing
  an unreachable remote — an unreachable path fails fast and would pass a timing
  assertion while still making the call.
- Link and upgrade convergence ride the existing managed-asset lifecycle gate
  phase, which already builds throwaway repositories for link, relink, and unlink.
  The symlink-parent fixtures are new; the rest of the phase is the regression
  control that conflict semantics did not move.
- Prior art: `internal/adopt/link_hook.go` and the lifecycle fragment for hook
  install, `internal/guards/guards_test.go` for row shape, and
  `internal/status/status_test.go` for the row budget.
- The feature gate is `bench gate`.

### Seam diagram: the hook-health record

```text
trigger: bench doctor | bench guards | bench status | SessionStart banner
    |
    v
repo root --> [ internal/adopt hook-health record ] --> {state, path, branch,
                   |         |                            provenance, currency}
                   |         +--> live origin/HEAD, else the token baked in the file
                   +--> hooks dir (core.hooksPath aware), Lstat then read
                      ^ tests attach here: temp repo, hook planted per state
```

### Seam diagram: currency is template drift only

```text
installed hook --> [ recover baked token ] --> substitute into embedded template
                                                        |
                                                        v
                                          exact byte compare --> current | stale
   ^ tests attach: the dead-protocol hook reds; a fresh install is current,
     whichever branch the installer baked
```

### Seam diagram: link plan convergence

```text
trigger: bench link | bench upgrade on a symlinked managed directory
    |
    v
plan entries --> [ already matches what this plan would write? ] --yes--> skipped, no write, no conflict
                        |
                        no  (new entry, or drifted entry)
                        v
                 [ symlink-parent check ] --> conflict, abort
                    ^ tests attach: lifecycle fixture with a symlinked .claude/commands,
                      three entry classes asserted independently
```

### Acceptance coverage map

The record does not exist today, so story 1's rows start red by construction. The
currency rows are built against this repository's own `.git/hooks/pre-push`, which
is dated Jul 5, speaks the dead `--describe` protocol, hard-codes `main`, and
carries no gate-pin drift clause — the observed live instance the map's #13
preserved for exactly this purpose. Rows naming an existing owner are already
covered and say so.

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | one exported computation returns state, hook path, effective branch, provenance, and currency for a managed hook | hook-health record | TDD-able: the record type does not exist, so the first assertion fails to compile. | A surface-local answer cannot satisfy a single-caller assertion over all five fields. |
| 1 | the effective branch is the live `origin/HEAD` when it resolves, and the token baked in the installed file when it does not | hook-health record | TDD-able as two fixtures sharing one installed hook: set `origin/HEAD`, then unset it, and assert the branch changes without reinstalling. | A record recomputing an installer resolver reports the same branch for both, which is the guard a fresh install would get rather than the one installed. |
| 1 | provenance is live when `origin/HEAD` resolves and baked when it does not, and names when the baked token is the bare literal fallback | hook-health record | TDD-able as three fixtures: resolvable, baked-from-a-real-branch, and baked-from-the-fallback. | Two-valued provenance collapses the last two, and the fallback case is the guess the map exists to surface. |
| 1 | currency compares installed bytes against the template substituted with the token baked in that same file | hook-health record | TDD-able: the dead-protocol hook must report stale, and a hook whose baked branch differs from any current resolution must still report current. | A comparison recomputing the branch reports stale for a perfectly current hook whose repository moved, conflating branch drift with template drift. |
| 1 | currency is not-applicable — never stale — for an absent, foreign, or diverted hook | hook-health record | TDD-able as three fixtures asserting the not-applicable value. | Comparing a foreign hook's bytes yields a mismatch, so a two-valued currency mislabels a foreign hook as stale drift. |
| 1 | the computation invokes no remote-contacting git subcommand | hook-health record | TDD-able with a `PATH` git recorder asserting the invoked subcommand set excludes `ls-remote` and `fetch`. | A resolver reusing the installer's remote probe hangs or degrades the ambient session path, and a timing assertion would not detect the call. |
| 1 | the four existing `PrePushState` values keep their exact meanings and set membership | hook-health record | Already covered by the existing classifier contracts, which run unchanged as the regression control. | It prevents the new fields from being smuggled in as a fifth state, which the map's discretion clause forbids. |
| 2 | doctor's pre-push row names the effective branch and its provenance for a managed hook | doctor runtime contract | TDD-able: assert both tokens in the row for a live-resolved fixture and a baked one. | A row printing only the state cannot tell an operator what the guard covers. |
| 2 | a baked-fallback branch leaves doctor's exit code unchanged from the live-resolved case | doctor runtime contract | Already covered in the sense that both are green today; it runs as the regression control that adding the field did not red the guessed case. | Reddening on a guess is the posture the map explicitly rejected, and only an exit-code comparison across the two fixtures catches it. |
| 3 | a stale managed hook reds doctor's pre-push row and names `bench doctor --fix` | doctor runtime contract | TDD-able: the dead-protocol fixture exits 0 today and must exit non-zero and name the remedy. | A report-only row leaves the operator with a visible problem and no route. |
| 3 | plain `bench doctor` never writes to the hook path | doctor runtime contract | Already covered by doctor's read-only pre-push contract; it runs unchanged against a stale fixture as the regression control for the new `--fix` capability. | Any self-heal shortcut added alongside `--fix` changes the recorded bytes. |
| 3 | `bench doctor --fix` rewrites a stale managed hook to the current substituted template and the record then reports current | doctor runtime contract plus hook-health record | TDD-able: repair the dead-protocol fixture, then re-read the record. | A `--fix` that touches only the shim leaves the re-read stale. |
| 3 | the repaired hook is executable | doctor runtime contract | TDD-able: assert the mode bits after repair. | Byte-correct content written without the executable bit disables the guard entirely while every content assertion passes. |
| 3 | `--fix` reports an already-current hook as an unchanged no-op and refuses a foreign hook without touching it | doctor runtime contract | TDD-able as two fixtures with byte-equality assertions after the run. | An unconditional rewrite destroys a foreign hook and makes repeat runs non-idempotent. |
| 3 | `--fix` does not install an absent hook | doctor runtime contract | Already covered by the current `--fix` contract, which never touches the hook; it runs unchanged as the regression control bounding the new repair. | Silently adopting an unlinked repository exceeds the repair the map authorized. |
| 4 | `bench guards` renders the pre-push row's branch, provenance, and currency from the shared record | guards AXI contract | TDD-able: assert the new cells for a baked, stale fixture. | A row without the fields leaves the ambient surface unable to report what the diagnostic one knows. |
| 4 | no package outside `internal/adopt` derives hook health, and `internal/adopt` has one derivation site | guards AXI contract plus hook-health record | TDD-able: unexporting `ClassifyPrePush` breaks the guards and status call sites at compile time, plus a call-site assertion inside the package. | Adding cells while a private parser or a second exported helper survives leaves the derivation the map's #7 exists to remove. |
| 4 | `bench status` emits no pre-push row for a managed hook with a baked branch, and one for a stale hook | status runtime contract | TDD-able: the stale fixture emits no row today and must emit one; the baked fixture must stay silent. | Rendering the green field as a status row spends the five-row budget on a non-signal. |
| 4 | the SessionStart banner surfaces the branch, provenance, and currency through `guards --brief` | guards AXI contract | TDD-able: assert all three fields in the brief output for a baked, stale fixture. | A brief that drops the fields makes the banner disagree with `bench doctor`. |
| 4 | a linked worktree renders no duplicate pre-push row | status runtime contract | Already covered by the existing primary-checkout guard, which runs unchanged. | Pool and linked worktrees share one `.git`, so a record-based rewrite could reintroduce double reporting. |
| 5 | `bench link` completes in a repository whose managed directory is a symlink when every entry beneath it already matches the manifest | lifecycle gate phase | TDD-able: a fixture mirroring this repository's symlinked `.claude/commands` aborts today. | The abort is the exact failure that makes story 3's remedy untrustworthy. |
| 5 | a *drifted* entry under a symlink parent still aborts | lifecycle gate phase | TDD-able only after the skip lands, as the negative control on its predicate; it is green today and must stay green. | A skip keyed on existence rather than manifest match silently writes through a deliberately symlinked directory. |
| 5 | a *new* entry under a symlink parent still aborts | lifecycle gate phase | TDD-able only after the skip lands, as the negative control on its predicate; it is green today and must stay green. | Sharing one fixture with the drifted case lets an implementation that skips all existing paths pass both; separate fixtures are what force the manifest-match predicate. |
| 5 | `bench upgrade` completes on the converged symlinked fixture and refreshes the hook | lifecycle gate phase | TDD-able: assert hook bytes after upgrade on the symlinked fixture, which aborts today. | Fixing only link leaves the other adoption route blocked by the same abort. |
| 5 | link over the converged symlinked fixture is re-run idempotent | lifecycle gate phase | TDD-able: run link twice and compare the tree and manifest. | A skip that also skips manifest bookkeeping diverges on the second run. |
| 5 | a clean managed asset whose kit bytes changed is refreshed by link and upgrade | lifecycle gate phase | TDD-able: two kits differing in one shared file; assert destination content and manifest hash after relink and after upgrade. | A skip keyed on the old manifest hash silently turns upgrade into a content no-op for every untouched asset. |
| 6 | `bench upgrade --check` counts the hook refresh when the installed hook differs from what upgrade would write, and does not when it matches | lifecycle gate phase | TDD-able as two fixtures at differing versions, differing only in the installed hook's bytes. | A plan-entry-derived count reports the same number for both, so the operator approves a plan that understates its writes. |
| 6 | at equal installed and linked versions the count includes no hook refresh, because that path performs none | lifecycle gate phase | TDD-able: an equal-version fixture with a stale hook must not count it. | An unconditional hook count promises a refresh on the one path that returns before the transaction runs. |
| 6 | `--check` performs no write | lifecycle gate phase | Already covered by the existing upgrade plan-mode contract, which runs unchanged. | It pins that adding the hook to the count did not add it to the plan-mode write set. |
| 7 | exactly one live site substitutes the branch token into the template | install/currency agreement | TDD-able: a focused call-site assertion, red while the dead installer's copy remains. | Two copies drift silently, and the dead one is what misled the decision map. |
| 7 | `bench link` sets `origin/HEAD` from the remote when a remote exists and the ref is unset, and a failing probe does not fail the install | lifecycle gate phase | TDD-able as two fixtures: one where the probe succeeds and the ref appears, one with an unreachable remote where link still exits 0. | Deleting the dead installer without moving its probe silently drops the only thing that populates `origin/HEAD` for later resolution. |
| 7 | a dangling symlink at the hook path is classified foreign, and `bench link` refuses it instead of overwriting it | hook-health record plus lifecycle gate phase | TDD-able: link overwrites it today because its refusal reads before it stats. | A read-first classifier reports a broken guard as absent and the transaction destroys evidence of it. |
| 7 | `bench link` refuses a marker-less foreign pre-push hook | lifecycle gate phase | Already covered by the existing transaction refusal, which runs unchanged. | It is the regression control that moving the stat ahead of the read did not weaken the existing refusal. |
| 8 | `README.md` states that the hook protects the branch it resolves rather than the default branch | conformance docs sweep | TDD-able: assert the narrowed phrasing and the absence of the unqualified claim. | The overstated sentence is the one place prose promises more than any check verifies. |
| edge of 1 | a hook path containing spaces or glob characters is read literally | hook-health record | TDD-able with `core.hooksPath` set to a directory named with a space and a bracket. | An unquoted or glob-expanded path reads the wrong file or none. |
| edge of 1 | a hook file that is a FIFO, device, or socket is rejected on its file mode before any read | hook-health record | TDD-able with a FIFO at the hook path and no writer, asserting the call returns. | A read-first implementation blocks the SessionStart banner forever. |
| edge of 1 | a managed hook whose final line lacks a trailing newline reports stale rather than current | hook-health record | TDD-able with bytes identical to the substituted template except the final newline. | Exact-byte comparison must not be softened to a trimmed compare, which would let real truncation pass. |
| edge of 1 | an empty file at the hook path is foreign, distinct from an absent one | hook-health record | TDD-able as two fixtures asserting different states. | A read-error-only branch collapses the two and misreports a truncated hook. |
| edge of 1 | a managed hook whose bytes do not share the template's pre-token prefix reports stale rather than failing to extract | hook-health record | TDD-able with a marker-bearing file that diverges before the first token position. | Token recovery that assumes a well-formed file panics or reports current on arbitrary marker-bearing content. |
| edge of 1 | a repository with no `git` on `PATH` degrades without a panic | hook-health record | TDD-able with a stripped `PATH` fixture. | An ambient banner that panics takes the session start with it. |
| edge of 3 | `bench doctor --fix` run twice over a stale hook is idempotent | doctor runtime contract | TDD-able: compare bytes, mode, and exit code across two consecutive runs. | A repair that rewrites unconditionally reports a change on the second run and falsifies its own report. |

Degenerate implementations are pinned per story. A five-field record filled by
each surface separately (1) fails the compile-time single-entry-point row. A
marker-substring currency (1, 3) fails the dead-protocol fixture, and one that
recomputes the branch fails the moved-repository row. A record recomputing an
installer resolver (1) fails the unset-`origin/HEAD` row. A row printing the
branch without provenance (2) fails the token assertion, and one that reds on a
baked branch fails the exit-code comparison. A `--fix` that rewrites
unconditionally (3) fails the foreign-refusal and idempotency rows, and one
writing correct bytes without the executable bit fails the mode row. Cells added
to guards while any second derivation survives (4) fails the compile-time row. A
skip keyed on existence rather than manifest match (5) fails the drifted-entry
control. A count left deriving from plan entries (6) fails the two-fixture
comparison, and an unconditional one fails the equal-version row. Leaving the dead
substitution copy (7) fails the call-site assertion, and deleting it without the
probe fails the `origin/HEAD` row.

### Edge inventory

- Error path — resolved by the missing-`git`, FIFO, dangling-symlink,
  malformed-prefix, and foreign-refusal rows.
- Empty or absent input — resolved by the absent-versus-empty rows, the
  not-applicable currency rows, and the unset-`origin/HEAD` row.
- Boundary values — resolved by the trailing-newline row, the already-current
  `--fix` no-op, and the equal-version upgrade row.
- Malformed input — resolved by the empty-file, dead-protocol, and
  malformed-prefix rows.
- Interrupted or partial state — **Won't handle** beyond what the existing link
  transaction already owns: this build adds no new multi-step mutation. `--fix`
  writes one file, and the link transaction's staging and rollback are unchanged
  and remain covered by the lifecycle phase.
- Re-run idempotency — resolved by the `--fix` twice row and the link re-run row.
- Hostile environment — resolved by the spaces/globs, special-file, symlink,
  `core.hooksPath`, linked-worktree, stripped-`PATH`, and no-remote-call rows.
- A command whose write changes a fact it reports — resolved by `--fix` re-reading
  the record after repair rather than asserting its own success, and by the
  idempotency row.
- Control bytes in git-sourced text — **Won't handle**: the only git-sourced value
  this feature renders is a branch name, and git rejects every ASCII control byte,
  space, and glob character in a ref name, so no such value can reach the TOON
  sink. Verified by attempting the branch creations rather than assumed. The hook
  *path* is not ref-constrained and keeps its spaces-and-globs row above.
- Flag values mistaken for positionals — **Won't handle**: this build adds no new
  flag. `--fix` is an existing recognized flag with no value, so there is no value
  to misread as a positional.
- Non-TTY stdin — **Won't handle**: no command on this path prompts, so stdin mode
  cannot amputate a caller.
- Host-backed filesystem I/O pressure — **Won't handle**: every read here is a
  single bounded file and `--fix` writes one; the existing link transaction owns
  durability for the staged route.

## Out of scope

- The doctor shim row, the doctor gate row, and the `guards` `wired` cell
  consuming this seam are the map's #12 consumer rows, each its own future
  capability — 6 edits, 2 gate runs each.
- The `gitguard` `denyTable` enforcement/advertisement duplication is a
  one-source-per-fact defect rather than a visibility gap and is its own roadmap
  row — 4 edits, 2 gate runs.
- Distinguishing one release behind from a dead protocol is refused by the source:
  #9 accepted a binary comparison and rejected the version token that would be
  needed for a finer answer.
- Making guards or hooks evasion-resistant is refused by the source, not deferred:
  `SECURITY.md` scopes them as advisory controls against honest mistakes.
- Doctor self-healing the hook without `--fix`, folding the hook into the SHA-256
  manifest, and a dead-protocol migration path are all rejected by the source
  rather than deferred.
- Repairing this repository's own `.git/hooks/pre-push` follows this spec rather
  than belonging to it, per the map's #13 — 1 edit, 1 gate run.
