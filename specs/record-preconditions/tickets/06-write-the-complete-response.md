# Write the complete probe or test response to a caller file

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_help.go (new), cmd/bench/main.go, cmd/bench/census_output.go, cmd/bench/tree_scope.go, cmd/bench/response_bound.go, cmd/bench/out_file.go (new), cmd/bench/out_file_test.go (new), cmd/bench/help_inventory_test.go, cmd/bench/command_help_test.go, internal/responsebound/owner.go, internal/responsebound/retain.go, internal/responsebound/owner_test.go, CHANGELOG.md
Covers: RP33, RP34, RP35, RP36, RP37, RP38, RP39, RP40, RP41

## What to build

Review chunk: RP-C4.

`cmd/bench/command_registry.go` is at its 400-line cap. First move `renderCommandHelp` and its help helpers (`command_registry.go:363-396`) into the new `cmd/bench/command_help.go`, with no behavior change. Then add one capability field to `commandDefinition` for the `--out` option, and set it only on the `probe` and `test` rows of `cmd/bench/main.go:107-108`. The help insertion reads the same field, as the tree-target insertion reads the scope field. Leave `probe.Usage` and `testreport.Usage` unchanged.

The dispatcher reads `--out <file>` as the first argument after the verb, or as the first argument after `--in <label>`. It resolves a relative path against the caller's working directory before the tree-target route, and it passes the absolute path on. Before the verb runs, it refuses at exit 1 when the path exists, when the parent directory is absent, or when the path holds a control byte. After the owner finishes in `runBounded` (`cmd/bench/census_output.go:20-34`), the dispatcher writes the complete output with `O_EXCL` and `O_NOFOLLOW` at mode 0644. Export the complete-output accessor of `internal/responsebound/owner.go:254-271`, and make retention (`retain.go:58`) read the same accessor. When no complete output exists, write no file, print `error: incomplete response — --out wrote no file, because the response has no complete copy` after the response, and exit 1. A write failure after a zero verb exit also exits 1. Stdout keeps its bounded response, so the closed bound-exemption set does not change.

Move the help expectations at `help_inventory_test.go:123-124` and `command_help_test.go:26`. Put new tests in `out_file_test.go`, because `main_test.go` and `command_registry_test.go` are over budget. Add one `CHANGELOG.md` entry under `## [Unreleased]` for RP-C4.

## Acceptance

- [ ] A spilled `bench probe --out <file>` call writes a file that `bench record verification --probe-result` imports at exit 0 (RP33).
- [ ] A spilled `bench test --out <file>` call writes the same bytes as its spill file (RP34).
- [ ] With `--out`, stdout still holds the bounded projection with its `spilled{` frame (RP35).
- [ ] An existing path refuses before the verb runs and keeps its bytes, and a symbolic link refuses and keeps its target bytes (RP36, RP37).
- [ ] Under `--in <label>`, a relative path writes under the caller's working directory (RP38).
- [ ] With no complete output, no file is written, `error: incomplete response` prints, and the exit is 1 (RP39).
- [ ] The `bench help` rows of `probe` and `test` show `[--out <file>]`, and no other row does (RP40).
- [ ] The complete-output accessor returns the same bytes for retention and for `--out` (RP41).

## Verification

Run `bench test --package ./cmd/bench` and `bench test --package ./internal/responsebound` whole. Run `bench test --check docs-currency-workflow`. Observe each planned row red before its production edit, and record the route of each red. Self-probe: write the projected stdout in place of the complete output, and report that RP33 and RP34 go red.
