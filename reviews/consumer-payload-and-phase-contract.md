# Review — consumer-payload-and-phase-contract

Three-axis semantic pass over `feebe4a..7c9d79f` (4 commits, 50 files, +2039/-152),
already landed on `main`. Advisory: the gate is green on this tree; these are the
findings it cannot see.

## Standards

**2 findings** (plus one judgment call deferred to the reviewer).
Worst: the allowlist parser/matcher is triplicated inside one Go package, in the
diff whose headline claim is "one canonical consumer allowlist".

1. **Duplicated knowledge — allowlist reader re-implemented three times.**
   `internal/conformance/package_core_checks_test.go:110` (`kitOnlyAllowlistPrefixes`)
   and `:131` (`excludedByAllowlistPrefix`) re-derive what `consumer_payload.go:68`
   (`PayloadKitOnlyPrefixes`) and `:81` (`PayloadExcluded`) already export. The
   sibling `internal/conformance/skills_index_checks_test.go:123` proves the import
   works while still reading the graded root's bytes, so the comment's stated reason
   does not require a second implementation — the helpers take rows as an argument.
   A third derivation sits at `internal/conformance/package_shipped_surface_test.go:78`,
   which declares an anonymous row struct instead of using `kitpayload.PayloadRow`.
   AGENTS.md: "a parser and its count … must collapse to one source"; the
   independently-authored-expectation exception explicitly does not cover parsers.

2. **Comment narrates the change and argues its own correctness.**
   `internal/conformance/package_shipped_surface_test.go:71-77` — "Before the wholesale
   `.agents/` entry was narrowed, this failed against every one of these rows."
   `craft-comments` forbids narration, provenance, and "argument for its own
   correctness"; the red proof belongs in the spec/review record, and this sentence
   becomes false the moment the allowlist changes.

*Judgment call (reviewer's, repo-wide, not this diff's):* the `FT<n> story <n>`
provenance tags conflict with `craft-comments` but match ~12 pre-existing sites
(`internal/adopt/setup.go:17` among them).

## Spec

**3 findings** against `specs/consumer-payload-and-phase-contract.md` (12 mapped rows).
Worst: the shipped tarball can no longer build from source.

1. **The npm tarball omits the Go source it now depends on.**
   `consumer_payload.go:9` is a new root package that `internal/adopt/link.go:11`
   imports and that `go:embed`s `.bench/consumer-payload.json`. Neither path is in
   `package.json` `files[]`, and `internal/packagesurface/assets.go`
   (`RequiredBuildPackAssets`) walks only `cmd/` and `internal/`, so root-level `.go`
   files are structurally outside the derived guard. Verified: `npm pack --dry-run
   --json` packs 217 files, neither among them. Story 3's row claims `files[]` carries
   "exactly the allowlist's consumer set", but the grading test
   (`internal/conformance/package_shipped_surface_test.go:71-113`) only asserts
   presence plus the derived `!` exclusions — nothing catches a *missing required* entry.

2. **Story 2's mapped behavior is not implemented; the test was fitted to the code.**
   The row demands the link plan's consumer entries equal the allowlist's `consumer`
   rows, with the red signal "a stray file under `.agents/skills/` is not written".
   `internal/adopt/link.go:107-115` still walks tree-typed rows as directories, so such
   a file *is* written. The landed test
   (`internal/contract/surface/link_lifecycle_test.go:74-88`) places its stray under
   `.agents/skills/bench-assess/` — a kit-only subtree, which is story 1's exclusion
   behavior. Same failure shape as `a6dcec3`: assertion fitted to implementation.

3. **Row 11 is half-closed.** The skills index is marked kit-only
   (`.bench/BENCH-reference.md:56`), but the Codex phase-adapter list at
   `.bench/BENCH-reference.md:87,89` still advertises `$bench-update-kit` and
   `$bench-assess` unmarked — in a file that is itself a consumer row
   (`.bench/consumer-payload.json:17`). Consumers read pointers to commands they
   never receive.

**The `b313d90` withheld-surface fix is real.** `.bench/consumer-payload.json:15`
names `.agents/skills/bench-craft-synthesis`, which exists; `upgrade_test.go:41`
asserts via `ReadFileAbs`, which fatals on a wrong path — non-vacuous by construction.
All allowlist source paths resolve on disk; no sibling row carries the `a6dcec3`
hollowness.

Row verdicts: rows 1, 4–10, 12 exercised; row 2 hollow; rows 3 and 11 partial.

## Coverage

**8 findings** (6 confirmed, 2 suspected). Worst, tied with Spec finding 1 on the
same defect: the kit-only guards go vacuously green when the allowlist's kit-only
rows are deleted — the withdrawal direction of the exact defect this spec existed
to close.

1. **CONFIRMED — tarball cannot build from source.** Same defect as Spec finding 1,
   reached from the packaging side: `go build ./cmd/bench` inside the packed tree
   (the `prepare` script) fails with "no Go files". No test builds from a packed tarball.

2. **CONFIRMED — deleting the kit-only rows makes every kit-only guard pass.**
   `internal/conformance/package_core_checks_test.go:88`
   (`if err != nil || len(prefixes) == 0 { return nil }`) and
   `package_shipped_surface_test.go:100` (`if row.Audience != "kit-only" { continue }`)
   both iterate zero rows silently. Remove the three kit-only rows and the tarball
   reships `bench-assess`, `bench-update-kit`, and `bench-craft-synthesis` with a green
   gate. The one canary (`kit-only-asset-admitted`) tests admission, never withdrawal.
   Missing: an assertion that the allowlist declares a non-empty kit-only set.

3. **CONFIRMED — the source-existence check the learnings entry proposes was not added.**
   `.bench/learnings.md:38-41` names it; nothing stats a row's `source`.
   `internal/adopt/link.go:104` skips a missing non-tree row, `treeEntries` returns
   `nil, nil` on a bad root, and `linkDestination` silently drops any row outside the
   five known prefixes. Re-typo a row to `.agents/skills/craft-synthesis` and the
   `a6dcec3` defect reproduces, still undetectable.

4. **CONFIRMED — prerelease→release upgrade is a silent no-op.**
   `internal/adopt/upgrade.go:145` — `compareKitVersions("1.2.3", "1.2.3-rc1")` strips
   the suffix, finds equal components, returns 0, and `upgrade.go:71` returns before
   `transactionalLink`. A repo pinned at `1.2.3-rc1` gets exit 0, a printed plan, no
   relink, and a manifest permanently stamped `-rc1`.
   `internal/contract/surface/upgrade_test.go:115` covers only identical strings.

5. **CONFIRMED — a symlink inside an allowlisted tree hard-fails link and upgrade.**
   `internal/adopt/link.go:181` refuses any non-regular file by name; a symlink is not
   `Mode().IsRegular()`. Only a FIFO is exercised
   (`link_lifecycle_test.go:106`), only under `.agents/skills`. One symlink anywhere in
   `.bench/lib`, `.bench/hooks`, `.bench/adapters`, or `.agents/commands` makes
   `bench link` and `bench upgrade` exit 1 for every consumer.

6. **CONFIRMED — the shell allowlist parser is format-fragile.**
   `.bench/skills-index.sh:37` uses a single-line `sed` requiring `"source"` before
   `"audience"`; pretty-printing or reordering keys drops every `(kit-only)` marker.
   Space-bearing sources break the unquoted `for source in $kit_only_sources` split at
   line 53. No fixture exercises a reformatted payload.

7. **SUSPECTED — duplicate and traversal rows unexercised.** `PayloadRows`
   (`consumer_payload.go:38`) validates only non-empty source and known audience.
   Duplicate rows yield duplicate `planEntry.rel` into `transactionalLink` and the
   manifest; `.bench/../../x` passes `linkDestination`'s prefix check and resolves
   outside the repo.

8. **SUSPECTED — `bench upgrade` argument and manifest-state edges.** No coverage for
   an unknown flag, `--check --force` together, an unparseable `#kit` header
   (`upgrade.go:151` treats it as an upgrade and relinks), an unreadable manifest, or
   concurrent runs.

**Verified non-vacuous:** all payload `source` paths resolve; the eight absence
assertions in `link_lifecycle_test.go:39-52` name real files; all four new canary
`forbid` needles are present in their fixtures.
