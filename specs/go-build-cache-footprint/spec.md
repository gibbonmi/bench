# go-build-cache-footprint

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-27 — the `/bench-write-spec go-build-cache-footprint` charge with its six numbered requirements, the measured problem, and the npm, uv, and pip question.

Verification log: 1 iteration(s) to accept — one independent `opus` medium round found five blockers and seven folds. The author folded all twelve. The largest folds are the `HOME`-only derivation, the unhashed closure entry, the ship conformance argv, the AXI-exempt verb, and the verb and gate-tail split.

## Problem

On 2026-08-27 the Go build cache on the development box held 53.4 GB across
259 shard directories. Every file was younger than seven days. Go's own trim
ran every day and removed entries unused for five days. The gate wrote 7 to
8 GB per day, so the trim never caught up.

The cause is the cache key, not the number of builds. Without `-trimpath`, Go
hashes the package directory into every compile action ID. Every Bench
worktree, every lane's private checkout, and every landing's composed checkout
therefore writes a complete new set of archives. `-count=1` disables only the
test-result lookup and never reduces those writes.

Two gates in one fresh worktree measured the multiplier on this box. The first
gate grew the cache by 198,839,178 bytes and 2,621 files. A second `--fresh`
gate on the same path grew it by 304,494 bytes and 73 files. The path key is
almost the whole per-gate cost.

Nothing in the development gate sets `GOCACHE`. The closed oracle environment
carries `PATH`, `BENCH_HOME`, and `HOME` only. Go therefore resolves the
directory from the machine's `go env -w` file or from its default. Only the
ship scripts stage a private cache.

An interim `go env -w GOCACHE=/home/mgibs/.cache/bench-go-build` on this box
relocates the writes but bounds nothing. It applies to every Go workload, so
Bench does not own that directory.

## Solution

Bench owns one Go build cache directory and hands it to every Go toolchain
child it spawns. The directory is `$HOME/.cache/bench/go-build`, derived from
`HOME` alone. One derivation serves the gate oracle closure, the phase
children, the run-binary builder, the fast lane, and `bench test`.

Every Bench-owned `go test`, `go vet`, and `go build` argv carries
`-trimpath`, and every test argv keeps `-count=1`. One producer owns the base
test argv, so the flag set cannot drift.

Every gate, lane, and focused test run holds a shared lock inside the
directory for its span. `bench cache clean` takes the exclusive lock without a
wait and runs `go clean -cache` against the directory, or it refuses and names
the holder. No gate evicts anything.

`bench cache` reports the directory, its apparent bytes, its file count, its
last trim time, the bound, and whether the footprint is over the bound. A green
gate prints one `go-build-cache:` line after the phase table, and every run
that reaches a verdict logs one `cache.footprint` event. The bound is 10 GiB. An over-bound
footprint warns and names `bench cache clean`. It never reds the gate.

The last task reverts the interim machine setting on this box and deletes the
orphaned directory. The closing ticket records that revert.

## Terms

- **Bench build cache** — the one directory Bench derives and hands to a Go
  toolchain child as `GOCACHE`. It is Go's build cache at a Bench-owned
  location. Avoid "shared cache", "go cache", and "the cache dir".
- **Ambient cache** — the directory Go resolves on its own when the child
  environment carries no `GOCACHE` entry.
- **Go child root** — an exec site where Bench constructs the environment of a
  Go toolchain child. Four exist: the gate oracle closure, `gateEnv`, the
  run-binary builder, and the `bench test` command.
- **Footprint** — the apparent byte total and the regular-file count of the
  Bench build cache, measured with `lstat` and with no file opened.
- **Bound** — the constant footprint above which the report warns: 10 GiB,
  which is 10,737,418,240 bytes.
- **Cache lock** — the record lock on `<directory>/bench.lock`: shared for a
  holder, exclusive for a clean.
- **Holder** — a Bench process that holds the cache lock in shared mode: a
  gate run, a lane run, or a `bench test` run.
- **Binary cache** — the wrapper's pinned platform binary store under
  `$BENCH_HOME/cache`, which `bench repair --prune` prunes. It is not this
  spec's subject.

## User stories

### Own the cache path

Line: opus / medium. The oracle closure hashes its environment, and four roots
must agree, so the seam is only partly known.

1. As a gate operator, I want each gate Go child to use the Bench build cache,
   so that no write hits the ambient cache.
2. As a gate operator, I want the gate's private run-binary build to write to
   the Bench build cache, so that it shares the phases' archives.
3. As a committer, I want the fast lane's Go children to use the Bench build
   cache, so that a lane checkout shares its worktree's archives.
4. As a developer, I want `bench test` to use the Bench build cache, so that a
   focused run warms the archives a gate reads.
5. As an operator, I want the directory derived as `$HOME/.cache/bench/go-build`
   from the env's `HOME` only, so that the location never obeys machine config.
6. As an operator, I want the derivation to read only the parent process
   environment, so that `go env -w` and an ambient `GOCACHE` never win.
7. As an operator, I want a root with no `HOME` or `XDG_CACHE_HOME` to refuse
   before the child starts, so that Go never uses its default. A relative
   value counts as absent, and the refusal names `HOME`.
8. As a reviewer, I want the closure to derive its entry from its own declared
   `HOME`, so that no undeclared variable steers the gate's cache. The entry is
   carried unhashed, because `HOME` is already hashed.

### Stop keying on the path

Line: opus / low. The argv owners are exact, and verbatim argv tests cover
them.

9. As a gate operator, I want the `test`, `race`, and `system` argvs to carry
   `-trimpath`, so that identical content in two checkouts shares one key.
10. As a gate operator, I want the `vet` argv to carry `-trimpath`, so that
    vet's compile and facts entries share one key across checkouts.
11. As a committer, I want the kit lane's `vet` and `build` argvs to carry
    `-trimpath`, so that a lane checkout writes no path-keyed archive.
12. As a developer, I want the `bench test`, `coreTestStep`, and ship
    conformance argvs to carry `-trimpath`, so that all Bench test forms share
    one rule.
13. As a gate operator, I want `-count=1` to stay on every Bench-owned
    `go test` argv, so that a test result is never reused. This is the
    `decisions/gate-budget.md` story 3 intent.
14. As a maintainer, I want one producer to own the base `go test` argv for
    every Bench test form, so that the flags cannot drift. The forms are test,
    race, system, focused, release, and ship conformance.
15. As a test author, I want the three `runtime.Caller` root helpers to use the
    working directory or git, so that their packages pass under `-trimpath`.

### Bound the footprint without a compile race

Line: opus / medium. The lock crosses processes, and a wrong lock mode fails
silently.

16. As a gate operator, I want each gate run to hold the shared cache lock for
    its span, so that no clean races a compile. The span opens before the
    run-binary build and closes after teardown.
17. As a committer, I want a lane run and a `bench test` run to hold the
    shared lock, so that one rule covers every holder.
18. As an operator, I want `bench cache clean` to take the exclusive lock at
    once and run `go clean -cache`, so that I reclaim disk. It reports the
    bytes and the files it removed.
19. As an operator, I want `bench cache clean` to refuse and name the holder
    while a holder exists, so that no compile loses an entry. The refusal
    exits 1 and names the holder's pid.
20. As an operator, I want `bench cache clean` on an absent directory to
    report zero and exit zero, so that a fresh machine passes.
21. As an operator, I want two concurrent gates to hold the shared lock
    together, so that the bound never serializes gates.
22. As a reviewer, I want no automatic eviction inside a gate, so that a
    gate's span and verdict never depend on another checkout's cache state.

### Observe the footprint

Line: opus / medium. A new verb joins several registries, and the report tail
has a fixed shape.

23. As an operator, I want `bench cache` to print one `go_build_cache` table
    with six columns, so that the size is one command away. The columns are
    the directory, the bytes, the files, the last trim, the bound, and the
    over-bound flag.
24. As an operator, I want `bench cache` to work outside a git repository, so
    that the machine-wide directory is inspectable from anywhere.
25. As an operator, I want `bench cache` on an absent or empty directory to
    report zeros and exit zero, so that a first run passes.
26. As an operator, I want a green gate to print one `go-build-cache:` line
    after the phase table, so that every green run shows the footprint. The
    line carries the bytes, the files, the directory, and the bound.
27. As an operator, I want the line to name `bench cache clean` above the
    bound, so that the next 53 GB is noticed. The line also says `over bound`.
28. As a gate operator, I want an over-bound footprint to leave the verdict
    green, so that disk pressure never grades the tree.
29. As a log reader, I want every run that reaches a verdict to record one
    `cache.footprint` event, so that a red run logs its footprint. The event
    carries the directory, the bytes, the files, and the over-bound flag.
30. As an operator, I want the footprint walk to open no file and follow no
    symlink, so that a FIFO or dangling link is harmless.
31. As an operator, I want a path with a control byte to make `bench cache`
    refuse, so that no table renders it. The refusal names the reason, and the
    gate line prints the path stripped.
32. As a reviewer, I want `bench status` unchanged, so that the dashboard's
    cost stays flat.

### Measure and revert

Line: none — the coordinating session runs these operator steps itself.

33. As a reviewer, I want the per-gate footprint measured after the build as
    before, so that the acceptance number is a before-and-after pair. The
    measurement is `du -sb` around one gate on a fresh worktree path and one
    `--fresh` gate on that path.
34. As an operator, I want the interim `go env -w GOCACHE` setting on this box
    reverted, so that no unowned directory looks Bench-owned. The revert runs
    `go env -u GOCACHE`, confirms `go env GOCACHE` as
    `/home/mgibs/.cache/go-build`, and deletes `/home/mgibs/.cache/bench-go-build`.
35. As a reviewer, I want the revert and the after numbers recorded in the
    closing ticket's checklist, so that the handoff reader sees them done.

## Implementation decisions

- One new deep module owns every cache fact: the directory derivation, the
  environment entry, the footprint walk, the cache lock, the clean, and the
  bound. No other package derives the path or walks the directory.
- The derivation reads `HOME` from the environment it is given, never from
  `go env`, and it does not read `XDG_CACHE_HOME`. The gate's closure carries
  no XDG name, so a value the gate cannot see must not steer its cache. A
  missing absolute `HOME` is an error that names `HOME`.
- The apply step removes every existing `GOCACHE` entry and appends the Bench
  entry. The four Go child roots call this one function. The oracle closure
  derives the entry from its own declared `HOME` and carries it unhashed,
  because `HOME` is already an identity frame. `.bench/gate-inputs.json` stays
  unchanged, because the entry is computed, not declared.
- A root that cannot derive the directory refuses before the child starts and
  names `HOME`. In the phase runner that refusal is a phase red with the
  reason on stderr.
- One exported function in the gate package produces the base test argv
  `go test -trimpath -count=1`. The test, race, system, focused, release, and
  ship conformance forms compose it, and the release package's own copy goes
  away. The vet argv and the kit lane's vet and build argvs take `-trimpath`
  from one shared flag owner in the gate package. That owner stays separate
  from the release build flags in `internal/releaseevidence/requirements.json`,
  because those flags state the shipped binary's evidence contract, not gate
  policy.
- The three `runtime.Caller` root helpers move to a working-directory or
  `git rev-parse --show-toplevel` resolution. The conformance harness already
  holds a git-based resolver to mirror.
- The cache lock is a POSIX record lock on `<directory>/bench.lock`. It
  mirrors the gate execution lock. A holder takes a read lock and keeps its
  descriptor open for the run's span. `bench cache clean` requests a write
  lock with a no-wait set, and `EAGAIN` is the refusal. The refusal names the
  blocking pid that `F_GETLK` reports. `go clean -cache` removes only the
  two-hex subdirectories, so the lock file survives a clean.
- `bench cache clean` composes `go clean -cache` with `GOCACHE` set to the
  directory. It measures the footprint before the removal and reports the
  difference. A missing `go` on `PATH` is a refusal that names `go`.
- The footprint walk uses `lstat` only. It recurses into a `-d` directory,
  because Go stores a cached executable as a directory. It counts regular files
  and never follows a symlink.
- The last trim time comes from `<directory>/trim.txt` as unix seconds after an
  `lstat` regular-file check. The report renders it as UTC RFC 3339, or empty
  when the file is absent, not regular, or unparsable.
- The gate reporter is one more report closure beside the capability-skips
  reporter in the verdict tail. It reads the directory from its own `GOCACHE`
  entry, walks once, and answers no rows, one green line, and never red. It
  logs `cache.footprint` through the inherited run log on every run that
  reaches a verdict. An interrupted run reaches no verdict, so it logs no
  event.
- The green line reads `go-build-cache: <bytes> bytes in <files> files at <dir> (bound <bound> bytes)`.
  Above the bound the parenthesis reads `(over bound <bound> bytes, next: bench cache clean)`.
- `bench cache` and `bench cache clean` join the command registry and the
  subcommand routing map. `bench cache` is an operational report exempt from
  the AXI query registry, like `bench structure`. `clean` is a mutation
  exemption, like `setup`.
- The bound is one constant in the new module. There is no knob.
- Bootstrap authority before execution: the interactive `bench cache clean`
  resolves the fixed literal `go` from the operator-owned `PATH` of its own
  process, which is the trust root of every interactive verb. A gate child
  resolves `go` through the closure's hashed `PATH`, and the closure hashes the
  `go` executable named in `tools`. The clean's delete target derives from
  `HOME` and never from an argument.
- Ship scripts keep their private caches. Consumer gate kinds and the release
  path keep their inherited environments.

## Testing decisions

- A good test drives an environment slice into a root and reads the child
  environment. Or it compares a produced argv verbatim. Or it holds the lock
  in one process and attempts a clean from a second process.
- Seams with prior art:
  - the phase-table argv tests and the ordinary build census
  - the lane table test
  - the `bench test` selection test
  - the oracle closure tests
  - the report tail tests
  - the command registry conformance tests
- The whole-tree gate observes every row through its `test` phase. The
  `system` phase's adoption journey observes the run log event after a real
  gate.

### Seam diagram

    trigger: bench gate | bench commit (lane) | bench test | bench cache [clean]
        │
        ▼
    process env  ──▶  [ derive + apply ]  ──▶  child env with GOCACHE=<dir>
                        ◀ tests attach here: drive an env slice, read the returned env

    phase table | lane table | test command  ──▶  [ base test argv producer ]  ──▶  go test -trimpath -count=1 …
                        ◀ tests attach here: compare the argv verbatim

    holder process     ──▶  [ cache lock: shared ]              ──▶  <dir>/bench.lock held for the span
    bench cache clean  ──▶  [ cache lock: exclusive, no wait ]  ──▶  go clean -cache | refusal naming the holder
                        ◀ tests attach here: hold in one process, attempt from a second

    gate-phases tail   ──▶  [ footprint reporter ]  ──▶  go-build-cache: line + cache.footprint event
                        ◀ tests attach here: inject a footprint and a bound, read stdout and the log

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| C01 | 5 | The derivation returns `<HOME>/.cache/bench/go-build` for an env with an absolute `HOME`. | derivation unit | A derivation that reads `go env` or the Go default returns another path. |
| C02 | 5, 8 | The derivation returns the same `HOME` path when the env also carries an absolute `XDG_CACHE_HOME`. | derivation unit | An XDG-honoring derivation returns a path the closure can never see. |
| C03 | 7 | The derivation returns an error that names `HOME` for an env with no absolute `HOME` and no absolute `XDG_CACHE_HOME`. | derivation unit | A fallback to the Go default returns a path instead of an error. |
| C04 | 6 | The apply step replaces an existing `GOCACHE` entry with the Bench entry. | derivation unit | A set-if-absent implementation keeps the ambient entry. |
| C05 | 5 | The derivation keeps a space in `HOME` unchanged in the returned path. | derivation unit | A word-split derivation returns the first word. |
| C06 | 1, 8 | The closed oracle env of a gate run carries `GOCACHE=<closure HOME>/.cache/bench/go-build`. | oracle closure unit | An omission at the closure leaves the oracle child with `PATH`, `BENCH_HOME`, and `HOME` only. |
| C07 | 8 | Two gate-run envs that differ only in an ambient `XDG_CACHE_HOME` produce the same closure entry and the same identity. | oracle closure unit | A derivation that reads the ambient env makes an undeclared variable an identity input. |
| C08 | 1, 3 | `gateEnv` carries the Bench `GOCACHE` entry and no other `GOCACHE` entry. | gate env unit | An injection at the closure only leaves a lane child with the ambient entry. |
| C09 | 2 | The run-binary builder's env carries the Bench `GOCACHE` entry. | runbinary unit | An untouched builder env passes the ambient entry through. |
| C10 | 4 | The `bench test` child env carries the Bench `GOCACHE` entry. | testreport unit | An untouched test env passes the ambient entry through. |
| C11 | 7 | A phase run whose env has no absolute `HOME` reds before the child starts with `HOME` on stderr. | phase runner unit | A silent fallback starts the child against the ambient cache. |
| C12 | 1, 29 | After one green gate in the adoption journey, the run's `.jsonl` holds one `cache.footprint` event whose `path` equals the derived dir of the journey env. | system journey | A gate that never hands the directory down records no path or another path. |
| T01 | 9, 13, 14 | The `test`, `race`, and `system` phase argvs each begin with `go test -trimpath -count=1`. | ordinary build census literals | A comparison against the producer itself passes, so only the census literals red on a missing flag. |
| T02 | 10 | The `vet` phase argv is `go -C <root> vet -trimpath ./...`. | phase table unit | An untouched vet argv fails the verbatim comparison. |
| T03 | 11 | The kit lane's `vet` and `build` argvs each carry `-trimpath`. | lane table unit | An untouched lane table fails the verbatim comparison. |
| T04 | 12, 13, 14 | The `bench test` argv carries `-trimpath` and `-count=1`. | testreport unit | An untouched focused argv lacks `-trimpath`. |
| T05 | 12, 13, 14 | The release `coreTestStep` argv carries `-trimpath` and `-count=1`. | gate-go unit | An untouched release argv lacks `-trimpath`. |
| T06 | 15 | `go test -trimpath` on `internal/runbinary`, `internal/conformance`, and `internal/preprelease` is green. | package run under `-trimpath` | A `runtime.Caller` root fails with `lstat github.com: no such file or directory`, as observed on 2026-08-27. |
| T07 | 12, 14 | The ship conformance step argv begins with `go -C <kit> test -trimpath -count=1`. | preprelease unit | A second producer in the release package keeps the old flags. |
| L01 | 16, 19 | While a gate run holds the cache lock, `bench cache clean` exits 1. | two-process lock test | A clean that takes no lock proceeds under a live gate. |
| L02 | 17, 19 | While a `bench test` run holds the cache lock, `bench cache clean` exits 1. | two-process lock test | A focused run that takes no lock lets the clean proceed. |
| L03 | 17, 19 | While a lane run holds the cache lock, `bench cache clean` exits 1. | two-process lock test | A lane that takes no lock lets the clean proceed. |
| L04 | 19 | A refused `bench cache clean` removes no file from the directory. | two-process lock test | A clean that removes before it locks empties the directory under a holder. |
| L05 | 21 | Two holders acquire the shared lock at the same time. | two-process lock test | An exclusive-only lock blocks the second gate. |
| L06 | 18 | With no holder, `bench cache clean` removes every two-hex subdirectory. | command test on a fixture directory | A clean that skips a subdirectory leaves bytes behind. |
| L07 | 18 | `bench cache clean` leaves `bench.lock`, `trim.txt`, and `README` in place. | command test on a fixture directory | A removal of the whole directory deletes the lock file a holder opens next. |
| L08 | 20 | `bench cache clean` on an absent directory reports zero removed and exits 0. | command test | A stat error becomes a refusal. |
| L09 | 22, 28 | A gate reporter given a footprint above the bound removes no file and answers not red. | report tail unit | An in-gate eviction removes entries or reds the run. |
| L10 | 18 | `bench cache clean` with no `go` on `PATH` refuses and names `go`. | command test | A bare exec error reaches the operator unnamed. |
| L11 | 19 | The clean refusal names the holder's pid as `pid <n>`. | two-process lock test | A refusal without the pid leaves the operator unable to find the holder. |
| L12 | 18 | With no holder, `bench cache clean` reports the bytes and files removed as measured before the removal. | command test on a fixture directory | A clean that measures after the removal reports zero. |
| R01 | 23 | `bench cache` prints one `go_build_cache` table with `dir`, `bytes`, `files`, `last_trim`, `bound`, and `over_bound`. | command test | A missing column fails the header comparison. |
| R02 | 23, 30 | `bytes` equals the size sum of the regular files under the directory, with a `-d` directory recursed. | footprint unit | A walker that skips a `-d` directory undercounts. |
| R03 | 25 | `bench cache` on an absent directory prints zero bytes and zero files with exit 0. | command test | A stat error becomes a refusal. |
| R04 | 25 | `bench cache` on an empty directory prints a zero-byte, zero-file row with exit 0. | command test | A walker that requires `trim.txt` refuses. |
| R05 | 24 | `bench cache` run in a directory outside any git repository exits 0. | command test | A git-root lookup refuses outside a repository. |
| R06 | 23 | `last_trim` renders the `trim.txt` unix seconds as UTC RFC 3339. | footprint unit | A raw integer or a local-time render fails the comparison. |
| R07 | 30 | `last_trim` is empty when `trim.txt` is absent, a symlink, or a FIFO. | footprint unit | A reader that opens a FIFO blocks. |
| R08 | 26 | A green gate prints `go-build-cache: <bytes> bytes in <files> files at <dir> (bound <bound> bytes)` after the phase table and before `gate: green`. | report tail unit | A missing or misplaced line fails the stdout comparison. |
| R09 | 27 | Above the bound the line's parenthesis reads `(over bound <bound> bytes, next: bench cache clean)`. | report tail unit | A line without the next action leaves the operator without the remedy. |
| R10 | 28 | Above the bound the run still prints `gate: green`. | report tail unit | A reporter that answers red reds the run. |
| R11 | 29 | A red gate run logs one `cache.footprint` event. | report tail unit | A reporter that logs only on green loses the red run's footprint. |
| R12 | 30 | A FIFO inside the directory adds zero bytes and does not block the walk. | footprint unit | A walker that opens files blocks on the FIFO. |
| R13 | 30 | A dangling symlink inside the directory adds zero bytes and no error. | footprint unit | A walker that follows links reports a not-found error. |
| R14 | 31 | A directory path with an ESC byte makes `bench cache` refuse with a named reason. | command test | An unfiltered path reaches the encoder and the table renders a control byte. |
| R15 | 31 | The gate line prints the directory with a control byte stripped. | report tail unit | An unfiltered path puts a control byte on stdout. |
| R16 | 23 | The `cache` verb routes to the new module in the subcommand routing map. | routing conformance test | An unregistered verb fails the routing census. |
| R17 | 29 | A red gate run prints no `go-build-cache:` line. | report tail unit | A reporter that prints on red lands its line ahead of the failure table. |

Not covered: story 32 — a reviewed exclusion, so `bench status` gains no section and no test changes.
Not covered: story 33 — an operator measurement, recorded in the closing ticket's checklist.
Not covered: story 34 — an operator step on one machine, recorded in the closing ticket's checklist.
Not covered: story 35 — the closing ticket's checklist is the record, and the reviewer reads it.

### Edge inventory

The walk applies the shell-CLI hostile-input checklist in `projects/benchkit.md`
to the four seams.

- Paths with spaces: C05 at the derivation, and the exec array carries the
  path with no quoting.
- Control bytes in a path from `HOME`: R14 refuses the table, and R15 strips
  the gate line.
- Absent versus empty directory: R03 and L08 for absent, R04 for empty.
- Special files in a discovered path: R07 and R12 for a FIFO, with no file
  opened by the walk.
- Dangling symlink: R13. A live symlink where `trim.txt` is expected: R07
  treats it as not regular.
- A command whose own write changes a fact it reports: `bench cache clean`
  measures before it removes and reports the difference. The gate reporter
  measures after the phases, so its line states the post-write footprint.
- Hand-edited `trim.txt` with no trailing newline: R06 trims whitespace before
  the parse.
- Required tool missing from `PATH`: L10 names `go`.
- Invocation through every shipped surface: R16 binds the verb to the routing
  map the wrapper and the binary share.
- `cwd` outside the repository: R05.
- Non-TTY stdin: `bench cache clean` prompts for nothing, so no TTY contract
  applies.
- State serialized by one process and read by a fresh one: L01 to L05 drive a
  second process for every lock row.
- Concurrent gates from several worktrees: L05.
- Destructive state under plan-and-apply drift: not applicable, because the
  clean has no fingerprint (see Won't handle).

**Won't handle** — the wrapper's own `dist/bench` rebuild keeps the ambient
cache — a shell-side derivation would be a second source of the path fact. An
operator's hand `scripts/go-build.sh` shares that exclusion, and both write one
`-trimpath` compile set at most.

**Won't handle** — a project-declared phase or lane argv in a linked repository
keeps its own flags — the project owns its manifest. Bench's default tables
carry `-trimpath` for the kit.

**Won't handle** — a Go process outside Bench pointed at the directory by hand
holds no lock — such a process is outside the contract.

**Won't handle** — `bench cache clean` carries no `--apply` fingerprint — the
target set is always the whole directory, and a rebuild recovers every entry.

**Won't handle** — the clean's duration on a very large directory — the removal
is Go's own, and a WSL2 host stall has no Bench-side remedy.

**Won't handle** — a symlinked `HOME` — the derivation returns the path string
as given, and Go resolves it.

**Won't handle** — `GOMODCACHE`, `GOTMPDIR`, and `TMPDIR` — none showed churn,
and the module cache grows only with a dependency change.

**Won't handle** — the `-race` phase's second key set — race builds are a
separate key set by Go's design, and `-trimpath` still dedupes them across
checkouts.

**Won't handle** — an interrupted gate run logs no `cache.footprint` event —
an interrupt reaches no verdict, and the report tail returns before any
reporter runs.

## Ownership fences

- `internal/gocache/` (new)
- `internal/gate/`
- `internal/runbinary/`
- `internal/testreport/`
- `internal/systemtest/`
- `cmd/bench/`
- `internal/conformance/subcommand_routing_test.go`
- `internal/conformance/ordinary_build_census_test.go`
- `internal/conformance/harness_test.go`
- `internal/preprelease/preprelease.go`
- `internal/preprelease/preprelease_test.go`
- `projects/benchkit.md`
- `CHANGELOG.md`
- `specs/go-build-cache-footprint/`

## Out of scope

- Automatic in-gate eviction under the exclusive lock at gate start, decided
  after the story 33 measurement: 3 edits, 2 gate runs.
- A `bench status` cache row with a registered action: 4 edits, 1 gate run.
- Release-path cache ownership for the `prep-release` step env and the
  `go run` form of `GateGoArgv`: 3 edits, 1 gate run, 1 release rehearsal.
- Consumer gate-kind cache pins for `npm_config_cache`, the pnpm store,
  `UV_CACHE_DIR`, `PIP_CACHE_DIR`, and `CARGO_HOME`. None depends on the
  checkout path, and none showed churn on 2026-08-27. A per-project versus
  Bench-owned directory decision comes first: 4 edits, 1 gate run.
- A wrapper-side derivation for the `dist/bench` rebuild through a Go plumbing
  verb: 3 edits, 1 gate run.
- A CI cache path for the Bench build cache in the workflows: 1 edit, 0 gate
  runs.
- `XDG_CACHE_HOME` support, which needs an optional-declaration form in the
  closure manifest so that the gate can see the variable: 3 edits, 1 gate run.
- The `decisions/gate-budget.md` drift at #7, which names a
  `ConformanceSuiteArgv` that does not exist and a `-count=1` addition that
  already landed. The map is reviewer-owned: 1 edit, 0 gate runs.

## Further notes

**Measurements on this box, 2026-08-27, against the interim directory.** The
cache held 368,803,070 bytes in 6,884 files before the first gate. One gate in
a fresh worktree path grew it to 567,642,248 bytes and 9,505 files. One
`--fresh` gate on the same path grew it to 567,946,742 bytes and 9,578 files.
The story 33 measurement repeats both runs on the built tree. The `du -sb`
apparent-byte figure is the acceptance number, and the `bench cache` `bytes`
column measures the same quantity.

**Go 1.25.0 facts, verified in GOROOT source.** The trim runs at most once per
24 hours at the end of a build, test, or vet command. It removes entries unused
for five days. An entry's mtime refreshes on use only when older than one hour.
Without `-trimpath` the compile action ID hashes the package directory.
`go vet` builds real compile actions through the same key.

Test binaries are not stored in the cache, so the churn is compiled archives
and vet facts. `go clean -cache` removes the 256 two-hex subdirectories with no
lock. A concurrent compile that already resolved an entry can fail when that
entry vanishes. `GOCACHE` precedence is the process environment, then the
`go env -w` file, then the default. Go refuses a relative `GOCACHE`.

**npm, pnpm, uv, and pip on this box.** None shows the churn pattern. The one
gate-path npm call is the conformance `npm pack --dry-run` probe. It already
pins `NPM_CONFIG_CACHE` to `$TMPDIR/bench-npm-cache`, where the day's 14 gate
runs wrote 124 MB in 95 files. `~/.npm` holds 81 MB, mostly a July `_npx`
tree, and no `_cacache`.

The pnpm store belongs to three other projects and had zero writes in the
week. uv and pip never run from this repository, and neither cache directory
exists. The consumer gate kinds inherit their caches unpinned, which is the
priced cut above.

**Interim state.** `go env -w GOCACHE=/home/mgibs/.cache/bench-go-build` was
set on 2026-08-27. Story 34 reverts it after stories 1 to 31 are green on the
integration source. The landing's private exact-source binary then already
carries the Bench-owned path.

**Build checkpoint.** Tickets 01 to 03 carry the measured saving on their own.
After ticket 07's measurement the reviewer decides whether tickets 04 to 06
build, because requirement 4 prefers that measurement before the bound's shape
is fixed. Ticket 04 blocks on ticket 07, so the checkpoint holds in the graph.

**Concurrency admission.** The gate execution lock admits one gate per
checkout, not one per machine. Several worktrees therefore gate at once against
the one Bench build cache. That is why the cache lock is shared for holders and
exclusive for a clean, and why no gate evicts.
