# Repair the root lookup that precedes help and argv dispatch

Blocked by: none
Ownership fence: `cmd/bench/specbuild.go`, `cmd/bench/specbuild_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: RH1/help-catalog-without-a-checkout, RH2/usage-two-without-a-checkout

## What to build

Close the accepted Spec finding (P1) from the Terra/xhigh review of candidate
`399dca908c7b1e1a4162eb7625e497dfb6786750`: `specBuildCommand` in
`cmd/bench/specbuild.go:19-33` resolves `git.Root()` on line 20 — before the
help-spelling check on line 28 and before `spec.ParseBuild` on line 31 — and
returns `buildError(errors.New("Git repository unavailable"), ...)` at exit 1
when that fails. The help catalog and the usage result are static, catalog-only
responses that read no repository: `specbuild.RenderBuildHelp` derives its rows
from `spec.BuildOperations()` alone, and `spec.ParseBuild` grades argv alone.
Outside a git checkout both therefore return error/1 today, against
`specs/axi-spec-build-complete/spec.md:28` ("`bench spec build help`, `--help`,
and `-h` return the operation catalog with exit 0") and `spec.md:33-34`
("Malformed operations and flag errors keep their usage/2 contract"; "The 0/1/2
exit taxonomy is unchanged: success/no-op 0, lifecycle refusal 1, usage 2").

Reorder `specBuildCommand` so the repository is resolved only for the paths
that need it. The help-spelling check runs first over the raw argv, then the
argv parse, and only then `git.Root()` and the `specbuild.New(...)` service
construction that feeds the no-argument family home and real operation
dispatch. Hoisting the help check above the no-argument branch is safe:
`spec.BuildHelpTarget` (`internal/spec/build.go:31-44`) returns false for empty
argv, so the live family home stays the no-argument response and keeps needing
a checkout. Nothing else about the three responses changes — same catalog
bytes, same usage bytes and code, same error/1 for a home or an operation that
genuinely cannot find a repository.

New coverage goes in `cmd/bench/specbuild_test.go`. No test there simulates a
missing repository today: every existing case (`:17-19`, `:64-67`, `:243`)
runs `exec.Command("git", "init", "-q", root)` on its `t.TempDir()` before
`os.Chdir`, so the failure state has to be constructed by chdir-ing into a
`t.TempDir()` with no `git init` and letting `git rev-parse --show-toplevel`
fail, restoring the old working directory through `t.Cleanup` exactly as the
existing cases do.

## Acceptance

- [ ] [RH1] (covers local) (P1) with the process working directory outside any
  git checkout, each of `help`, `--help`, and `-h` at the family root returns
  the nine-operation catalog at exit 0, byte-identical to the same spelling
  inside a checkout, and each per-operation spelling returns its one-row
  catalog at exit 0.
- [ ] [RH2] (covers local) (P1) with the process working directory outside any
  git checkout, a malformed operation and a flag error each return the pinned
  usage bytes at exit 2 rather than a repository error at exit 1; the
  no-argument family home and a well-formed operation still return the
  repository error at exit 1 there, since both genuinely need the checkout.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RH1/help-catalog-without-a-checkout | move the `git.Root()` resolution back above the `spec.BuildHelpTarget` check | focused `specBuildCommand` test | chdir into a `t.TempDir()` that was never `git init`-ed, invoke each of `help`, `--help`, `-h` at the root and one per-operation spelling, and require the catalog bytes at exit 0 |
| RH2/usage-two-without-a-checkout | move the `git.Root()` resolution back above the `spec.ParseBuild` call while leaving the help check first | focused `specBuildCommand` test | chdir into a `t.TempDir()` that was never `git init`-ed, invoke an unknown operation and a known operation with a bad flag, and require the usage bytes at exit 2 while the no-argument home still returns exit 1 |
