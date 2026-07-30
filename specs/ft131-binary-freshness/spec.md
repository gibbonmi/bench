# FT131 — freshness-seal the built Bench binary

Status: staged

## Problem

AXI and runtime contracts can execute whatever `dist/bench` happens to be on disk instead of the package sources under test. That produces a false red when correct sources are judged by an old binary, or a false green when an old binary masks a broken source change. The gate has the same fault: its shell entry starts a pre-build binary to resolve `gate-phases`, so the first run after a table change can grade the previous table and only an unchanged second run sees the new one.

## Solution

Make one freshness owner compare the current build-input set with a seal written only by a successful build. Contract launch helpers and the gate entry invoke that owner before they credit or use `dist/bench`. A mismatch is actionable red; it never silently falls through to an older executable or phase table.

### Defaulted decisions and veto points

No decision map exists. Every decision below is **[defaulted — post-hoc veto]** under the reviewer-directed batch-drain override.

| decision | default | veto consequence |
|---|---|---|
| stale-binary posture | **Refuse before execution** with a stable diagnostic naming the binary and rebuild command. **[defaulted — post-hoc veto]** | **Rebuild/re-resolve** would make the gate/helper build then resolve the replacement in one invocation, adding write/recovery rules. **Diagnostic-only** would run the stale binary/table and retain the false-green hole. |
| one owner and bootstrap | A Go freshness package owns source-set discovery and seal validation. Go contract helpers call that package directly from current test sources; `.bench/gate.sh` invokes the selected `dist/bench` freshness subcommand before `gate-phases`. An unknown/missing freshness subcommand is a legacy red, never a positive verdict. **[defaulted — post-hoc veto]** | Separate checks require a named drift guard or duplicate policy. A shell reimplementation or trusting an unknown subcommand would recreate the circular trust bypass. |
| comparison | Use a deterministic content seal, not mtime ordering; a tie never credits a binary. **[defaulted — post-hoc veto]** | An mtime policy must define ties and clock skew with extra false-red/false-green proofs. |
| source set | Derive the command's local Go build graph plus build-script/version/flag inputs from the build owner, not a hand-maintained list. **[defaulted — post-hoc veto]** | A narrower set risks undetected staleness; a broader set must prove unrelated edits do not red. |
| lifecycle | Publish executable and seal only after successful build replacement; interrupted or partial output is untrusted. **[defaulted — post-hoc veto]** | Any alternative must still fail closed and define recovery. |
| missing/legacy artifact | Missing executable/seal, unreadable seal, or binary without freshness support reds with the rebuild action. **[defaulted — post-hoc veto]** | Treating one as fresh creates an unverified bypass. |
| gate hand-off | `.bench/gate.sh` verifies before asking `gate-phases`; after a source/table change first run refuses, rebuild then second run resolves the new table. **[defaulted — post-hoc veto]** | Re-resolve must prove no old table is scheduled; diagnostic-only preserves five-then-ten behavior. |

## User stories

1. As a contract-suite author, I want a stale `dist/bench` rejected before AXI or runtime assertions run, so a source change cannot receive a false red or false green. Line: gpt-5.6-terra / high. The oracle-facing seam is known but requires paired fixture proofs.

2. As a gate operator, I want the first gate run after a phase-table/source change to refuse before an old binary resolves phases, so the verdict answers for the current table. Line: gpt-5.6-terra / high. Pre-build ordering and recovery are safety-critical.

3. As a developer in a fresh or salvaged worktree, I want missing, partial, legacy, or indeterminate output to explain how to restore a trusted binary, so I do not mistake a skipped or stale contract for source failure. Line: gpt-5.6-terra / high. Hostile filesystem states need explicit black-box coverage.

4. As a developer whose files share a coarse timestamp, I want deterministic freshness across re-runs and cwd, so clock granularity never makes trust timing luck. Line: gpt-5.6-terra / high. Content-seal and dual entry surface behavior need deeper fixtures.

## Implementation decisions

- The freshness verifier is a deep Go module that owns build-input discovery, seal validation, and diagnostic vocabulary. Contract helpers import it directly from the sources their `go test` process compiled; the built executable exposes only a thin `freshness-check <root>` plumbing subcommand over that same module. Its callers receive only fresh or actionable failure; they do not compare timestamps or rebuild independently. **[defaulted — post-hoc veto]**
- The build owner resolves every regular local file the `./cmd/bench` Go build graph consumes, plus `go.mod`, `go.sum` when present, and build-script/version/flag inputs. It hashes normalized repository-relative names plus bytes. Generated and untracked regular inputs count when the build resolves them; ignored status is not an exemption. **[defaulted — post-hoc veto]**
- A successful build publishes a companion seal only after replacing the executable. Missing, unreadable, non-regular, symlinked, partial, legacy/no-seal, malformed, or mismatched artifacts are not fresh. **[defaulted — post-hoc veto]**
- The contract helper is the sole test-side adapter. It checks before AXI wrapper selection or runtime direct invocation, preserving source-under-test attribution. **[defaulted — post-hoc veto]**
- The gate shell entry resolves the exact executable it would otherwise use for `gate-phases`, invokes `<selected-dist-bench> freshness-check <root>` first, and only then invokes `<selected-dist-bench> gate-phases <root>`. A nonzero result, missing executable, or the exit-2 unknown-subcommand answer from a pre-FT131 executable is normalized to the stable legacy/stale red; no result from that executable counts as fresh until the check returns success. The shell has no source scanner and does not auto-build. Its red instructs `bash scripts/go-build.sh <root> <root>/dist/bench`, then a gate rerun. This is enforced refusal with a remedy, not guidance-only wording. **[defaulted — post-hoc veto]**
- Existing phase dependencies continue to serialize build writes before readers once phase execution begins; FT131 supplies the earlier check needed before old-process table resolution. **[defaulted — post-hoc veto]**

## Testing decisions

- Test external behavior at the highest existing seams: the contract fixture launch helper and executable gate entry. Use a lower verifier seam only for malformed input states the high seams cannot observe. **[defaulted — post-hoc veto]**
- The feature must pass `.bench/gate.sh`; focused tests run before the whole gate during implementation.
- Prior art: `BenchkitPhases` already expresses build-to-reader edges, `phaseTable` fails malformed manifests closed, and runtime contracts run `dist/bench` directly.

### Seam diagram

**Contract launch seam**

```text
trigger: AXI or runtime contract test
    |
    v
sources + dist/bench + seal --> [ freshness owner ] --> fresh / actionable red
                                ^ tests attach: fixture source change, exit, diagnostic
                                |
                                +--> [ wrapper or direct dist/bench ] --> assertion
```

**Gate entry seam**

```text
trigger: .bench/gate.sh from any repository cwd
    |
    v
sources + dist/bench + seal --> [ freshness owner ] --> fresh / actionable red
                                ^ tests attach: table change and emitted phase set
                                |
                                +--> [ dist/bench gate-phases ] --> current table --> runner
```

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A changed source input makes each current AXI suite using `runBenchInDir` fail before an old binary runs. | contract launch seam | Future focused fixture — not observable before implementation. | An old executable that emits expected AXI output cannot be credited. |
| 1 | A stale executable that would make a broken runtime assertion pass and one that would make a correct assertion fail both red with freshness, not assertion output. | contract launch seam | Future paired runtime fixtures — not observable before implementation. | Paired mutations prove both false-green and false-red directions. |
| 2 | First `.bench/gate.sh` run after an in-scope phase-table/source change refuses before old `gate-phases` output or scheduling. | gate entry seam | Future gate fixture — not observable before implementation. | The observed five-phase then ten-phase unchanged rerun cannot pass first invocation. |
| 2 | After prescribed successful rebuild, unchanged second gate invocation resolves/runs new table once. | gate entry seam | Future gate fixture — not observable before implementation. | Proves recoverability and resolution from replacement rather than a cached process. |
| 3 | Missing executable/seal, unreadable seal, legacy executable (including unknown `freshness-check`), and interrupted/partial output fail closed with one rebuild action. | contract launch seam and gate entry seam | Future fixture cases — not observable before implementation. | Each state would otherwise skip or execute an unverified subject; an old executable cannot vouch for itself. |
| 3 | Symlinked or special executable/seal paths are never trusted or followed. | verifier unit seam | Future hostile-filesystem fixture — not observable before implementation. | A link can smuggle bytes; a FIFO can block verification. |
| 4 | Equal mtimes, newer executable, and older executable are decided by content seal only. | contract launch seam | Future coarse-mtime fixture — not observable before implementation. | Equal mtimes with different bytes make an mtime-only implementation red. |
| 4 | Source-set enumeration includes resolved generated/untracked regular build inputs and excludes unrelated non-input files. | verifier unit seam plus contract launch seam | Future build-graph fixture — not observable before implementation. | It catches both a hand-maintained omission and an overbroad false red. |
| 4 | Root and nested-cwd invocations agree; rebuild followed by unchanged repeated checks stays green. | both high seams | Future nested-cwd/re-run fixture — not observable before implementation. | Root resolution and idempotence stay externally observable. |
| 2, 3 | The in-scope shipped launch surfaces — direct kit gate, linked-repository by-path gate, hook-triggered gate, and adapter-triggered gate — all reach the same pre-`gate-phases` freshness refusal and retain their existing invocation contracts. | gate entry seam | Future routed-surface fixture — not observable before implementation. | A surface that bypasses the check can still credit the old table even when the direct gate is safe. |
| 3 | A successful freshness check does not mutate tracked configuration or make its own executable/seal appear stale on an unchanged second check. | contract launch seam | Future tracked fixture — not observable before implementation. | It catches a verifier that changes the source fact it reports or self-invalidates its own seal. |

### Degenerate-implementation check

- Always-green verifier: rows 1–6 and 9–11 red on changed source, missing/malformed artifact, special-file, legacy-subcommand, self-write, or stale-table fixtures.
- Timestamp-only verifier: row 7 reds on equal mtimes with different source bytes.
- Hand-maintained package list: row 8 reds on generated/untracked resolved input and stays green for unrelated input.
- Gate check after table resolution: row 3 reds because old table output would appear first; a legacy binary's unknown subcommand cannot be treated as fresh.
- Guidance-only rebuild text: rows 1–6 red because an unverified binary still executes or skips.

### Quantifier enumeration

- Consumers: (1) AXI tests using `runBenchInDir`; (2) runtime tests directly invoking subject `dist/bench`, currently capability-skip; (3) `.bench/gate.sh` before `gate-phases`; (4) the four routed gate launch surfaces in row 10.
- Build inputs: every regular local file resolved into `./cmd/bench`'s Go graph, plus the build-script/version/flag inputs above; build owner enumerates per build.
- Indeterminate artifacts: missing, non-executable, non-regular, unreadable, partial/interrupted, legacy/no-seal, malformed seal, digest mismatch.

### Edge inventory

| edge class | resolution |
|---|---|
| error path | Rows 1–6 cover stale, missing, legacy, unreadable, malformed, and mismatch refusal. |
| empty/absent input | Row 5: absent executable/seal reds, never skips or infers fresh. |
| boundary values | Row 7: equal/coarse mtime, newer, and stale output use content comparison. |
| malformed input | Rows 5–6: malformed seals and symlink/special paths red. |
| interrupted or partial state | Row 5: interrupted build cannot publish a trusted pair. |
| re-run idempotency | Rows 4 and 9: rebuild then repeated checks resolve current behavior. |
| hostile environment | Rows 6 and 9: symlink/special files, spaces/globs, and nested cwd are fixtures. |
| generated/untracked state | Row 8: resolved generated/untracked inputs count; unrelated files do not. |

### Hostile-input checklist disposition

Every disposition in this table is **[defaulted — post-hoc veto]**; it records
the complete project checklist without treating an unrelated surface as an
implicit FT131 exclusion.

| checklist class | resolution |
|---|---|
| spaces or glob characters in paths | Row 9 exercises root and nested paths with spaces/globs through both high seams. |
| refused control bytes in git-sourced text | **Won't handle:** freshness inputs are build files and seals, not rendered git text; TOON's existing refusal owns this separate renderer capability. |
| permitted tab/newline/return splitting a line sink | **Won't handle:** the seal is binary-safe hashed input, not a markdown/single-line sink; adding a user-visible report format is separate (2 edits, 1 gate run). |
| command self-write changes reported fact | Row 11 proves the verifier excludes its output pair from source inputs and does not mutate tracked configuration. |
| hand-edited file lacks trailing newline | Row 8 hashes raw bytes, so a missing final newline in a resolved build input changes the seal deterministically. |
| absent versus present-but-empty | Row 5 distinguishes absent executable/seal from empty/non-executable or malformed files; neither is fresh. |
| special discovered files | Row 6 covers FIFO/device/socket rejection before open. |
| dangling symlink | Row 6 covers dangling and live symlinks as non-regular, before any read. |
| unquoted multi-word arguments | Row 9 uses spaces/globs in root and artifact paths through shell entry. |
| flag value mistaken for positional | **Won't handle:** `freshness-check` has one required root positional and no flags; parser misuse is impossible at this surface. |
| required tool missing from PATH | **Won't handle:** direct contract helpers use compiled current Go sources and the stale check uses the selected executable; build-tool absence belongs to the explicit rebuild command and existing build failure. |
| invocation through a symlink | Row 6 refuses symlinked artifacts; row 10 exercises the routed gate surfaces rather than assuming real-path invocation. |
| every shipped invocation surface | Row 10 enumerates direct kit, linked-repository by-path, hook, and adapter launches. |
| destructive worktree state | **Won't handle:** worktree registration/recovery is unrelated to reading the current root's artifact pair; existing worktree contract owns it, while row 9 keeps nested worktree cwd callable. |
| SIGINT mid-loop | **Won't handle:** freshness itself performs no write; interrupted build publication is row 5 and gate runner process-group teardown remains its existing contract. |
| re-run idempotency | Rows 4, 9, and 11 cover rebuild/re-run, root/cwd repeat, and self-write-free repeat. |
| cwd deeper than root | Row 9 covers nested cwd resolution. |
| non-TTY stdin | **Won't handle:** freshness and gate entry are non-prompting; no primary caller is removed. |
| host-backed I/O pressure | **Won't handle:** no deterministic fixture exists for host VHDX stalls; atomic publication and refusal on unreadable/partial state are row 5, while latency bounds remain a separate capability (3 edits, 2 gate runs). |

**Won't handle**

- Clock rollback/cross-host timestamp reconciliation — seals have no timestamp semantics; remote attestation is separate (3 edits, 2 gate runs).
- Auto-rebuilding from contract tests/gate entry — separate write-and-recovery capability with concurrency/offline policy (5 edits, 3 gate runs).
- Consumer-installed platform-package freshness — separate release-artifact provenance capability (4 edits, 2 gate runs).

## Out of scope

- Automatic rebuild/re-resolution after a stale result — a distinct in-launch write capability needing ownership, interruption, and offline semantics — 5 edits, 3 gate runs.
- Remote or signed artifact attestation — a distinct release/distribution capability, not local source-to-binary freshness — 4 edits, 2 gate runs.
