# Own the Go build cache path

Blocked by: none

Writes: internal/gocache/ (new), internal/gate/subject.go, internal/gate/gate.go, internal/gate/runner.go, internal/gate/cache_env_test.go (new), internal/runbinary/runbinary.go, internal/runbinary/runbinary_test.go, internal/testreport/testreport.go, internal/testreport/testreport_test.go, CHANGELOG.md

## What to build

Bench derives one build cache directory and hands it to every Go toolchain
child. Add the package `internal/gocache/` as the one owner of that fact. No
other package derives the path.

The package exports two functions. `Dir(env []string) (string, error)` returns
`<HOME>/.cache/bench/go-build` from the slice's `HOME` alone. It never reads
`XDG_CACHE_HOME`, because the gate's closure declares no XDG name and must not
follow a value it cannot see. A relative value counts as absent, and `Dir`
returns an error that names `HOME` when no absolute `HOME` is present.
`Apply(env []string) ([]string, error)` calls `Dir`, removes every existing
`GOCACHE` entry from the slice, and appends the Bench entry. Both functions
read the given slice only, never `go env`.

The four Go child roots call `Apply`. The roots are the oracle closure in
`internal/gate/subject.go`, `gateEnv` in `internal/gate/gate.go`,
`buildEnvironment` in `internal/runbinary/runbinary.go`, and the `bench test`
command env in `internal/testreport/testreport.go`. The phase runner in
`internal/gate/runner.go` composes `gateEnv`, so it reds the phase and prints
`HOME` on stderr when `Apply` returns an error.

The oracle closure derives its entry from its own declared `HOME` and carries
the entry unhashed, because `HOME` is already an identity frame. An ambient
`XDG_CACHE_HOME` therefore never reaches the subject identity. The file
`.bench/gate-inputs.json` stays unchanged, because Bench computes the entry and
does not declare it.

`Dir` and `Apply` are the seam that three later tickets consume.
`04-add-the-bench-cache-verb.md` reads the directory for the footprint walk,
`05-hold-the-cache-lock-and-add-bench-cache-clean.md` opens the lock file
inside it, and `06-print-the-gate-line-and-log-the-footprint.md` reports it.

## Acceptance

- [ ] C01 — The derivation returns `<HOME>/.cache/bench/go-build` for an env with an absolute `HOME`.
- [ ] C02 — The derivation returns the same `HOME` path when the env also carries an absolute `XDG_CACHE_HOME`.
- [ ] C03 — The derivation returns an error that names `HOME` for an env with no absolute `HOME` and no absolute `XDG_CACHE_HOME`.
- [ ] C04 — The apply step replaces an existing `GOCACHE` entry with the Bench entry.
- [ ] C05 — The derivation keeps a space in `HOME` unchanged in the returned path.
- [ ] C06 — The closed oracle env of a gate run carries `GOCACHE` set to the closure `HOME` plus `/.cache/bench/go-build`.
- [ ] C07 — Two gate-run envs that differ only in an ambient `XDG_CACHE_HOME` produce the same closure entry and the same identity.
- [ ] C08 — `gateEnv` carries the Bench `GOCACHE` entry and no other `GOCACHE` entry.
- [ ] C09 — The run-binary builder's env carries the Bench `GOCACHE` entry.
- [ ] C10 — The `bench test` child env carries the Bench `GOCACHE` entry.
- [ ] C11 — A phase run whose env has no absolute `HOME` reds before the child starts with `HOME` on stderr.

Delivered outcome: every Go toolchain child that Bench spawns writes to the one
Bench-owned cache directory, and no undeclared variable steers it.
