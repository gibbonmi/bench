# Filter consumers by production or test and take several symbols in one query

Blocked by: none
Writes: internal/consumers/command.go, internal/consumers/consumers.go, internal/consumers/query.go (new), internal/consumers/query_test.go (new), internal/outline/outline.go, cmd/bench/main.go, cmd/bench/help_inventory_test.go, CHANGELOG.md
Covers: none

## What to build

`bench consumers` gets two forms from roadmap row FT125.

The first form is a production-or-test filter.
`--production` keeps only the consumer rows whose file is not a `_test.go` file.
`--test` keeps only the consumer rows whose file is a `_test.go` file.
The filter applies before the row cap, so the meta counts and the cap read the filtered rows.
An action that the response offers repeats the filter flag.

The two flags together are a grammar error at exit 2.
A filter flag with `--changed` is also a grammar error at exit 2, because the `--changed` form is out of scope.

The second form takes several `<qualified-symbol>` operands in one query.
Each block then has a leading `symbol` column, and each row names the operand that it answers.
The row cap and the candidates form apply to each symbol separately.
One meta row, one citation row, and one help block close the response.
A query with one operand keeps its current output bytes.

The usage line in `internal/consumers/command.go` and the help inventory row in `cmd/bench/main.go` show the two forms.

## Acceptance

- [ ] `bench consumers <symbol> --test` returns only the rows in `_test.go` files.
- [ ] `bench consumers <symbol> --production` returns only the rows that are not in `_test.go` files.
- [ ] `bench consumers <symbol> --production --test` prints a usage line and exits 2.
- [ ] `bench consumers <symbol> <symbol>` returns the rows of both symbols in one `consumers` block with a leading `symbol` column.
- [ ] One `citation` row closes the answer of several symbols above the help block.
- [ ] `bench help` and `bench consumers --help` show the new usage.
