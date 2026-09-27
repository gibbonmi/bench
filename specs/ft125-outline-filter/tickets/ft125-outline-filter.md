# Filter outline symbols to production or test files

Blocked by: none
Writes: internal/outline/outline.go, internal/outline/outline_filter_test.go, cmd/bench/main.go, cmd/bench/help_inventory_test.go, specs/ft125-outline-filter/tickets/ft125-outline-filter.md
Covers: none

## What to build

`bench outline [path] [--full]` gets two filter flags. The flag `--production` removes each Go test file (a file name that ends in `_test.go`) from the walk. The flag `--test` keeps only the Go test files. Each flag scopes the walk as the path argument does, so the metadata counts only the files in scope. The filter applies to the symbol form and to the directory summary.

A call that gives the two flags together is a grammar error, and the command refuses it at exit 2 before it reads the repository. A call without a filter flag gives the same bytes as before.

The help text and the `bench help` row show the new flags. This change is one face of FT125: a survey delegate got 91.7 KB from `bench outline internal/worktree --full`, and 60 percent of the symbols were test symbols.

## Acceptance

- [ ] `bench outline --full --production` on a mixed tree shows no row from a `_test.go` file and shows every other row.
- [ ] `bench outline --full --test` on the same tree shows only the rows from `_test.go` files.
- [ ] `bench outline --production --test` prints the usage line with the refusal reason and exits 2.
- [ ] `bench outline --full` and the bare `bench outline` give the same bytes as before on the same tree.
- [ ] `bench outline --help` and `bench help` show `[--production|--test]`.
