# Seed the gate input manifest at setup

Blocked by: none
Writes: internal/adopt/setup.go, internal/adopt/setup_test.go, internal/adopt/setup_prompt_test.go, .bench/BENCH-reference.md

## What to build

`bench setup` seeds `.bench/gate-inputs.json` through the same transaction that
writes the gate and the profile. One Go function, `scaffoldGateInputs()`, beside
`scaffoldGate()` and `scaffoldProfile()`, returns the exact bytes: schema 1,
`local` closure, `environment` of `BENCH_HOME` and `HOME`, empty `paths`, and
`tools` of `bash`, `basename`, `dirname`, `git`, `readlink`, `uname`, indented
with a trailing newline. The plan entry uses the existing `seed` kind, so the
file is written only when absent, never touched when present, and never
recorded in `link-manifest.tsv`. The seed lands on every `setup` — zero-signal
and detected-ecosystem alike.

`inspectRepo` records whether the file exists and `renderSetupPreview` prints one
line in the profile line's voice: absent → will be seeded declaring `BENCH_HOME`
and `HOME`; present → left as-is. `.bench/BENCH-reference.md`'s sentence on
`.bench/gate-inputs.json` gains the clause that `bench setup` seeds it with the
names the installed wrapper needs.

The tests attach to the existing setup fixture (`setupPromptTestRepo`, which
already binds `BENCH_KIT` and a temporary git repository carrying `go.mod`) and
to a second, zero-signal fixture with no build-system file at all. The seeded
bytes are asserted as an independently authored expectation — the one place the
content is spelled outside the function — because dropping `BENCH_HOME` reds
nothing else in the tree (mutation probe d) and an invalid shape must red here.

## Acceptance

- [ ] after `setup --yes` in the fixture repository, `.bench/gate-inputs.json` holds exactly the seeded bytes (SD1).
- [ ] the preview names the file as about to be seeded when absent and as left as-is when present (SD2).
- [ ] an operator-authored `.bench/gate-inputs.json` survives `setup --yes` byte-identical and is absent from `link-manifest.tsv` (SD3).
- [ ] a second `setup --yes` leaves the seeded file byte-identical (SD4).
- [ ] `setup --yes` in a zero-signal repository — no `go.mod`, `Makefile`, `package.json`, or `Cargo.toml` — also writes the seeded bytes (SD5).
- [ ] `.bench/BENCH-reference.md` says `bench setup` seeds the file.
