# go-doctor-link-port — slice 6 of the Go rewrite

Child of `decisions/go-rewrite.md` (#6, slice order): `doctor` + `link`, the
highest-stakes mutators, ported once Go idioms are settled. Bootstrap evidence:
`bench-link.sh` (373 lines: plan/preflight/install, fence-aware AGENTS.md marker
rewrite, pre-push generation, manifest), `bench-doctor.sh` (200: shim
report/fix), `bench-init.sh` (105: gate/canary/learnings scaffold — unassigned
in the parent slice list, adopted here by #1). The npm entry point is
`bin/bench.sh` itself — there is no separate node launcher. The regression net
is strong and black-box on both seams: `gate-link-contracts.sh` (13 cases —
fresh link, relink, conflicts, fenced/malformed markers, worktree, hooksPath,
metachar kit path, default-branch) and `gate-doctor-contracts.sh` (~14 sandboxed
cases — health states, fix idempotency/atomicity, foreign refusal, spaced
target, PATH fallback), all invoking the CLI as a subprocess.

## #1: Does `bench init` ride in this slice?

Type: Grill

### Answer
**Yes.** The parent slice list never assigned `init` a home; it is the same
adoption-mutator family (sourced by the same dispatcher, needs the same kit-dir
crossing), and its contracts live in the same `gate-link-contracts.sh` fragment
this slice updates anyway. Rejected: leaving it as a stranded sourced shell
file for a later slice.

## #2: Does the kit-version stamp land here?

Type: Grill

### Answer
**Yes, both halves this slice owns:** `link` writes the installed kit version
into the link manifest, and `doctor` gains the skew check (stamp vs
`bench version`) — parent #4 committed to this as the skew detector for the
bundled `.bench/dist/bench`. The session-start half stays out: that surface
belongs to the slice-5 hooks spec, which lands before this one. Rejected:
deferring the stamp (consumer repos would accumulate unstamped manifests a
later slice must special-case).

## #3: Kit-dir asset walk, or `go:embed`?

Type: Grill

### Answer
**Kit-dir walk; embed rejected.** Parent #3/#4 closed on "plain files are the
product / auditable surface"; embedding freezes skills and commands into the
binary and creates exactly the binary↔asset version skew #2's stamp exists to
detect. Accepted consequence: `link`/`init`/`doctor` work only through the
kit-installed wrapper (which passes its resolved kit dir down), and the
consumer-planted `.bench/bin/bench.sh` refuses adoption subcommands with a
pointer at the real kit rather than half-working.

## #4: What happens to the planted `.bench/bin/` surface?

Type: Grill

### Answer
**Organic shrink — no deliberate thin-shim conversion this slice.** Link keeps
copying whatever `bin/*.sh` remains in the kit (dispatcher with `run_gate`,
worktree, postinstall) plus the bundled `dist/bench`; the planted dispatcher
still needs the shell gate path for the stop hook until slices 7–8, so a shim
conversion now would build a transitional shape twice. Parent #4's thin-shim
end state arrives for free when slice 8 empties `bin/`. The link contracts drop
their `.bench/bin/bench-link.sh` assertions this slice, and the planted
dispatcher's `link`/`init`/`doctor` gain the refuse-with-pointer from #3.

## Handoff

1. **Module boundaries.** New `internal/` package(s) for adoption logic — the
   split (one package vs link/init/doctor apart) and naming are the spec's
   call, following the `gitguard` precedent. `cmd/bench` gains `link`, `init`,
   `doctor` subcommands; `bin/bench.sh` routes all three through the strangler
   dispatch and deletes the three sourced files (`bench-link.sh`,
   `bench-doctor.sh`, `bench-init.sh`). Still shell after this slice: the
   dispatcher (`run_gate`), `bench-worktree.sh` (slice 7), postinstall and the
   generated pre-push (permanently shell — the pre-push *text* is now emitted
   by Go, unchanged).
2. **Contracts.** All observable behavior carries over — both contract
   fragments run unchanged except the planted-file assertions (#4) and new
   stamp rows. Link: preflight-refuses-before-any-mutation; every conflict
   class (project-owned file, modified managed file, foreign pre-push,
   malformed markers, unclosed fence, symlink parent) exits non-zero naming
   the conflict with nothing rewritten; manifest stays TSV `rel<TAB>sha256`
   with one writer; the stamp row must never parse as a file row (encoding is
   the spec's call); a pre-stamp manifest reads as skew-unknown — warn, never
   fail. Doctor: states healthy/stale/foreign/missing with exit 0 only on
   healthy; fix is idempotent, atomic (temp+rename), refuses foreign files
   byte-identically; report never executes shim contents; the shim continues
   to target the resolved *wrapper* path, never the platform binary. New:
   doctor reports version skew; the planted dispatcher's adoption subcommands
   exit non-zero pointing at the installed kit. Wrapper→binary context (kit
   dir + wrapper's resolved path) crosses explicitly — env vs flags is the
   spec's call; the existing `BENCH_KIT` override keeps working.
3. **Deep vs thin.** The binary is the deep unit: plan build, preflight,
   fence-aware marker rewrite, pre-push emission, manifest read/write/stamp,
   scaffolds, shim template + target-dir selection (manager-owned exclusion
   list). `bin/bench.sh` stays a thin router; postinstall stays a thin
   best-effort caller that degrades loudly-but-zero when the binary is
   missing.
4. **Black-box assertables.** The two fragments are the port-parity net —
   identical stdout/exit/filesystem assertions before and after. New `go test`
   tables: marker/fence parser edges (reversed, fenced examples, unclosed, no
   trailing newline), manifest parse (stamp row, absent file, duplicate rel),
   adapter symlink target computation, doctor dir selection and shim
   content/readback round-trip.
5. **Gate attachment.** The unchanged shell gate is the oracle; the fragments
   plus `go build`/`vet`/`test` cover the slice. Gate-blind, as today: a real
   multi-package `npm i -g` (manual smoke per release — verify the kit-dir
   crossing under the global layout while there), and PATH-probe behavior
   beyond the sandboxed cases.
6. **Hostile-input owners.** Spaced/glob/metachar kit and target paths → Go
   exec argv + the existing metachar and spaced-target contract cases.
   Symlink parent dirs → ported preflight check, table-tested. A hostile
   marker-bearing file on PATH → report's no-execute readback (contract
   exists). SIGINT mid-fix → atomic write carries over. Re-run idempotency →
   relink and second-init contracts. Unset origin/HEAD → default-branch
   contract. Unwritable target dir → doctor's error paths. Binary missing at
   postinstall (`--ignore-scripts`, off-matrix host) → postinstall's
   best-effort rim.
7. **Uncertainty flags.** The wrapper→binary context crossing under a real
   global npm install (hoisted platform package, no repo checkout) is
   gate-blind — the spec schedules the manual smoke and escalates per
   `craft-line` if kit-dir resolution differs from the dev layout. npm
   optional-deps lockfile edge remains open from the parent map (asset, #4).
8. **Rejected alternatives.** `go:embed` (#3); deliberate thin-shim conversion
   now (#4); deferring `init` (#1) or the stamp (#2); pointing the doctor shim
   at the platform binary (recreates the stale-shim failure doctor exists to
   repair — the wrapper owns binary resolution).
9. **Domain watch-outs.** These subcommands mutate reviewer-owned files in
   consumer repos — refuse-before-mutate is the safety property, and a
   half-ported seam must never have two live implementations (the dispatcher
   routes whole subcommands, never shares one). The link manifest is both the
   ownership test and, after this slice, the skew detector — it keeps exactly
   one writer. A doctor shim must always exec the wrapper: anything pinned to
   a version-specific `node_modules` path is the stale-shim failure class.

Dependency order: single spec; lands after the slice-5 hooks spec
(`decisions/go-hooks-port.md`), per the parent map's order.
