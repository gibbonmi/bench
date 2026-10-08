# Derive the ticket check floor from the Writes line

Blocked by: none
Writes: internal/gate/phases.go, internal/testreport/floor.go (new), internal/testreport/floor_test.go (new), internal/testreport/selection.go, internal/testreport/selection_facts.go, internal/testreport/selection_facts_test.go, internal/tickets/registry_data.go, internal/tickets/registry_data_test.go
Covers: TF7, TF8, TF9, TF10, TF11, TF12, TF13, TF14, TF15, TF19, TF20, TF21, TF22, TF71

## What to build

`internal/testreport` exports the floor function `WritesFloor(ctx context.Context, root string, entries []string) (Floor, error)` in `floor.go`. The `entries` are the raw `Writes:` entries of one ticket, as `tickets.Ticket.Writes` holds them, `(new)` markers included. `Floor` holds `Packages []string`, the directory-relative package patterns in sorted order, and `Embed bool`, the embed flag. The error is the loader error. The spec section "Floor derivation" owns the entry rules, and this ticket implements them unchanged.

The floor and the `--changed` selection share one package index built from `go list -json -test ./...`. The index maps a directory to its package, an embed file to its owner, and a package to its embed patterns. The listed package record gains the embed patterns that `go list` reports. Keep the signature of `selectCurrentPackages`, so `selection_test.go` and `selection_cause_test.go` stay unedited. The file `selection_test.go` holds 604 lines and has no structure grant.

The floor splits each entry with `tickets.WritesPath` and matches a directory entry with `tickets.Covers`. Each element is the pattern that `packagePattern` gives the package directory, so the root package is `./.`. A root without `go.mod` returns an empty floor and calls no loader. The tests inject canned package records through the loader seam, as `resolveChangedPackagesWithLoader` does.

`Request` gains the accessor `func (r Request) Changed() bool` in `selection_facts.go`, which reports the `--changed` flag. The `--changed` form reports the target `./...`, so a reader without this accessor cannot tell the two forms apart. TF-C3 reads the accessor and writes no file in `internal/testreport`.

`internal/tickets` exports the predicate `CommandPackage(path string) bool` in `registry_data.go`. It answers whether a path sits under a bound command package, through `Covers`. The predicate reads the binding rows whose files are the command registry set, so no second list of command packages exists.

At the ticket commit, rewrite the `planned` seam cell of each covered row into the citation form.

## Acceptance

- [ ] TF7: A present Go file entry adds the package of its directory.
- [ ] TF8: A `(new)` Go file entry in a current package directory adds that package.
- [ ] TF9: In a two-package fixture, a directory entry adds both current packages at or below it.
- [ ] TF10: The directory entry `internal/git` adds no package from `internal/git2`.
- [ ] TF11: An embed file entry adds its owner package and sets the embed flag.
- [ ] TF12: A `go.mod` entry adds the pattern `./...`.
- [ ] TF13: A Markdown entry that no package embeds adds no package.
- [ ] TF14: A `(new)` entry under no current package directory adds no package and returns no error.
- [ ] TF15: The tickets predicate answers true for a path under a bound command package, and false for a path under the `internal/toon` row.
- [ ] TF19: A root without `go.mod` gives an empty package floor and calls no loader.
- [ ] TF20: The floor and the `--changed` selection give the same direct package for one canned path.
- [ ] TF21: `TestCauseEmbedOnly` stays green after the index is shared. Its red comes from a mutation probe, not from absence: omit `EmbedFiles` from the shared index's embed map, and the test turns red.
- [ ] TF22: `TestCauseChangedBeatsImports` stays green after the index is shared. Its red comes from a mutation probe, not from absence: let the imports-cause pass overwrite a direct cause, and the test turns red.
- [ ] TF71: An absent `(new)` embed file adds no package and returns no error.
- [ ] The `--changed` accessor answers true for a `bench test --changed` request and false for a `bench test --package` request.
