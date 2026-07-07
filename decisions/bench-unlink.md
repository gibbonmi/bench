# Bench Unlink (FT33)

## #1: What does `bench unlink` remove, and what does it spare?

Blocked by: —
Type: Grill

### Question
The link manifest enumerates every installed path with fingerprints, but
nothing consumes it to remove the footprint. What are the removal semantics?

### Answer
`bench unlink` consumes the manifest: it removes managed files whose
fingerprint still matches (unmodified), removes now-empty managed
directories, removes the pre-push hook only when bench-managed, and removes
the fenced bench block from AGENTS.md while leaving surrounding user content.
Modified-managed files are left in place and reported by path — a user edit is
user content now. User-owned artifacts are never touched: ROADMAP.md,
IDEAS.md, CONTEXT.md, `.bench/learnings.md`, `specs/`, `decisions/`,
`reviews/`, the profile, and a `gate.sh` that no longer matches its scaffold
fingerprint. Exit report lists removed / kept-modified / never-managed. The
manifest itself is removed last, only if nothing else refused. Rejected:
`--force` full removal (deleting user-modified files is the reviewer's hand,
not a flag); interactive prompts (bench porcelain is non-interactive).

## #2: Does it need a rehearsal mode?

Blocked by: #1
Type: Grill

### Question
Unlink is the kit's only multi-file destructive command. Is the
report-after-removal contract enough?

### Answer
Yes, with `--dry-run` as the rehearsal: prints the exact removal plan
(remove / keep-modified rows) and exits 0 without touching anything. Cheap to
build (same walk, no rm), and it makes the destructive default defensible.
The README documents both forms plus the manual path for un-manifested
repos (pre-manifest installs). **Flagged for veto:** default-destructive with
opt-in dry-run, versus default-dry-run with `--apply` — chose the former to
match every other bench command doing what its name says.

## Handoff

1. **Module boundaries.** `internal/adopt` owns unlink (reuses manifest,
   fingerprint, AGENTS-fence helpers from link); `bin/bench.sh` routes the
   subcommand; README owns the story.
2. **Contracts.** `bench unlink [--dry-run]`: exit 0 on success (including
   kept-modified files); exit 1 when the manifest is absent or unreadable
   (nothing to consume — the false-empty rule: absence is loud, not a silent
   no-op); report format stable for tests.
3. **Deep vs thin.** The removal walk with fingerprint verdicts is the deep
   piece; dry-run is the same walk with a different sink.
4. **Black-box assertables.** Temp repo: link → unlink leaves only user
   content; link → modify a managed file → unlink keeps it and reports;
   dry-run changes nothing (tree hash identical); AGENTS.md user prose
   survives fence removal; absent manifest → exit 1.
5. **Gate attachment.** Runtime contract tests against the built binary
   (adopt family).
6. **Hostile-input owners.** adopt owns a hand-edited manifest (unknown
   paths → skip and report, never delete outside the repo root),
   path-traversal rows (`../` → refuse), and a manifest naming a directory.
7. **Uncertainty flags.** Whether the manifest records the pre-push hook and
   AGENTS fence as entries or they need bespoke handling — implementer reads
   the manifest writer first and mirrors it.
8. **Rejected alternatives.** `--force`; interactive confirmation;
   uninstall-via-docs-only.
9. **Domain watch-outs.** Unlink must never follow symlinks out of the repo;
   every removal path is verified inside the repo root before rm.

Dependency order: n/a — single spec.
