# 2. Move the Bench home read into a stdlib-only package

Blocked by: none
Line: opus / medium
Rows: none — this ticket is the prefactor that the span tickets need
Writes: internal/benchhome/ (new), internal/benchhome/benchhome.go (new), internal/benchhome/benchhome_test.go (new), internal/worktree/effects.go, internal/worktree/exec.go

## What to build

A new package, `internal/benchhome`, owns the one `BENCH_HOME` read. The
package imports the standard library alone, so any package can import it.

`internal/worktree` keeps its exported `Home()`, and that function delegates to
the new package. No caller changes, and no second read of the variable
survives.

This move is a prefactor. `internal/worktree` imports `internal/gate`, so the
gate cannot import the read where it lives today. Every later ticket resolves
the home through `internal/benchhome` at its verb boundary.

An unset `BENCH_HOME` still falls back to `~/.bench`, and the fallback keeps
its current behavior.

The `homeEnv` name moves with the read. `internal/benchhome` exports the constant, and the child-env injection in `internal/worktree/exec.go` reads it from there, so both halves still name it once.

## Acceptance

- [ ] `benchhome` returns the value of `BENCH_HOME` when the operator sets it.
- [ ] `benchhome` returns the `~/.bench` fallback when `BENCH_HOME` is unset.
- [ ] `worktree.Home()` returns what `benchhome` returns, and it holds no second read.
- [ ] `internal/benchhome` imports the standard library alone.
- [ ] the gate `test` phase stays green for `internal/worktree` and its callers.
