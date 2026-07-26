# Review — ft91-gate-tier-split

Branch `bench/assign/7c14860ec945e5cb1388ee76fa54ffae/1d9a75cfeb06758f10f0f6102ed3f745`,
base `5b2ccf1`, four commits (`d4cffb5`, `4233456`, `a19ca79`, `7355e60`), 30 files.
Advisory: the gate is green on this branch and stays the oracle.

## Standards

4 findings (2 hard, 2 judgment).

**Worst: the release-only package list in `projects/benchkit.md` is an unanchored
second derivation.** `projects/benchkit.md:231-233` enumerates `internal/preflight`,
`internal/releaseevidence`, `internal/publication`; the enforcement is
`releaseOnlyPackages` at `internal/conformance/package_core_checks_test.go:243`. The
diff added a conformance anchor for `bench-final-check.md` but none for the profile
paragraph, so the list drifts silently on the next retier. AGENTS.md: "two derivations
of the same fact (an enforcement and its advertisement…) must collapse to one source."
Fix by naming the seam rather than the members, or by anchoring the string.

Hard violations:

1. `internal/canary/fixture_tier_test.go:13` — comment cites provenance: "`shipFixtures`
   is *the spec's own statement* of which fixtures follow the release-evidence probe to
   the ship tier." The next reader has no spec. `craft-comments` forbids provenance; the
   rationale that follows ("written out rather than derived so that an implementation
   which ships nothing is red here") is the part that survives merge. Drop the first clause.
2. Comments narrate the change. `internal/contract/surface/preprelease/fixture_test.go:123`
   — "the ~372 s probe *this whole split exists* to keep off the common path"; same register
   leaks into a failure message at `internal/conformance/tier_test.go:25` ("the ~372 s probe
   *the split moves*"). "The split" is diff-talk. State the tier fact, not the migration.

Judgment calls:

3. The profile/enforcement duplication above.
4. `internal/canary/canary.go:141` owns `const checkFileName = "CHECK"`; the same-package
   test hard-codes `"CHECK"` at `fixture_tier_test.go:87,92,116`. AGENTS.md's
   test-independence exception requires a recorded red, and the spec never mentions the
   CHECK file at all — nothing earns the second derivation.

Refuted, not reported: tier membership *is* single-sourced (`registry.Checks` owns names,
order, and tiers; `RunsAt` is the only superset rule; `TestRegistryBindsEveryCheck` closes
both directions). The literal package list at `package_core_checks_test.go:371` is covered
by story 2's recorded red signal.

## Spec

3 findings. All 32 coverage rows audited; every named test exists on the branch, and the
four "already covered" classifications are honest.

**Worst: story 7's ship-canary step is structurally unreachable, and blocks ship green.**
Spec row: "the two probe-derived fixtures still bite under `bench prep-release` — reports
`did not bite` if the ship-tier canary step is skipped." Traced and confirmed against the
tree: `preprelease.Steps` → `canary.SweepShip` → `SweepTier(root, Ship, defaultRunner)` →
`innerEnv()` (`internal/canary/canary.go:381`), which copies `os.Environ()` and never sets
`BENCH_CONFORMANCE_TIER`. The variable is written *per-step* on `conformance-ship` only
(`internal/preprelease/preprelease.go:94`), so the prep-release process env does not carry
it; `BENCH_CONFORMANCE_TIER` appears nowhere in `internal/canary` or `bin/bench.sh`. The
inner `.bench/gate.sh` therefore resolves `entryTier("") == Dev`, `release-evidence-probe`
never runs, and the two EXPECT strings emitted only by `runReleaseEvidenceProbe` cannot
appear. Both new fixtures report "did not bite", so `prep-release` can never exit 0 on the
real tree. Nothing tests it — `internal/contract/surface/preprelease/fixture_test.go:86`
seeds a dev-only fixture deliberately. This is a stronger claim than the handoff's
"unreached because preflight blocks first": fixing the preflight blockers would not close
this row.

2. **The 600 s package-timeout hazard is restaged, not decided.** Spec: "`internal/preflight`
   alone runs 676 s and blows the 600 s default package timeout, which presents as a gate
   hang." `internal/preprelease/preprelease.go:93` runs the ship conformance step with no
   `-timeout`, so the ~372 s probe plus all three release-only suites now share one
   default-bounded `./internal/conformance` run. Dev is fixed; ship inherits the hang.
3. **Missing step (minor — flagged for veto, not a defect).** Seam 2 lists
   `[ release-only go test ]` as its own box and the decisions say `prep-release` "calls
   `go test` over the three release-only packages." No such step exists; it is folded into
   ship-tier `goCoreTestPackages`. Behavior closed, mechanism changed silently.

No findings on two of the three reviewer focus areas. **RunsAt is correct**:
`registry.go:59` `return c.Tier == Dev || tier == Ship` makes ship a strict superset,
`crossCompileMatrix` sits inside the Dev `package-core-guard` check that ship re-runs with
`-tags stress`, and nothing runs at ship while skipping dev. **`prep-release` does refuse
up front** (`preprelease.go:140`, before `gate.Inspect`), naming the tool — which is what
the row's rationale demands. See Coverage #3 for the limit on that refusal.

## Coverage

5 findings.

**Worst: an untiered registry row silently leaves the dev gate, with every existing test
still green.** `registry.Tier` is `type Tier string` (`registry.go:21`), so a row added with
a typo'd or omitted tier gets the zero value `""`. `RunsAt(Dev)` is then false and
`RunsAt(Ship)` true — the check stops running on every commit and only a release rehearsal
would notice. `TestRegistryBindsEveryCheck` (`internal/conformance/tier_test.go:43`)
validates against `Names(registry.Ship)`, which *contains* the untiered check, so it passes;
`TestTierMembership` only asserts dev ⊆ ship; `TestDevTierExecutesExactlyDevChecks` compares
the timing file to `Names(Dev)`, which drops it on both sides. Missing row: *every
`registry.Checks` entry carries a tier in {Dev, Ship}* — one loop, no fixture.

2. **A present-but-empty `CHECK` file.** `fixtureTier` (`internal/canary/canary.go:148`)
   treats absent as dev, but an empty file falls through `TrimSpace` to the not-carried
   branch and errors `names check "" ...`. Absent-vs-present-but-empty is a named class in
   the profile's hostile-input checklist. `TestFixtureTierResolution` covers absent,
   valid-ship, and unknown-name only. Add empty and whitespace-only to that table.
3. **A tool reached only from inside a step — `govulncheck`, absent on this host.**
   `requiredTools` (`preprelease.go:57`) pre-resolves `bash`, `git`, `go`, `node`; the
   vulnerability scanner resolves four levels in, at `internal/preflight/vulnerability.go:30`.
   `TestPrepReleaseMissingTool` exercises only `go` and `node`, and the contract fixture stubs
   the preflight script, so the real phase never runs. A real `prep-release` here burns the
   artifact matrix plus a full ship conformance run, then dies at `release-preflight`. This
   is the concrete limit on the up-front refusal the spec claims. Missing row: PATH without
   the scanner exits 1 with a named `required tool is missing` attribution.
4. **Recovery after an interrupt.** `TestPrepReleaseInterrupt` asserts only that no index
   landed; the orphaned `dist/.preflight.*` staging directory it leaves is never cleaned or
   re-run over. Missing row: a second `prep-release` after an interrupt still reaches exit 0
   and leaves no accumulated staging dirs.
5. **Two conformance runs on one root.** `NewTimingWriter` truncates while another run
   appends (`registry.go:164-194`); nothing serializes them, and `ReadTimingLines` prints
   interleaved or renumbered lines. Print-only, never reds a gate. Needs a row or an explicit
   **Won't handle** line — the spec's concurrency prose covers only the ascending-path case.
