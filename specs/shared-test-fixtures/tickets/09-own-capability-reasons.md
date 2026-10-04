# Own unavailable-capability reasons

Blocked by: 08-share-system-fixture-identity.md
Writes: internal/bounds/classify_test.go, internal/capability/capability_test.go, internal/capability/unavailable.go (new), internal/capability/unavailable_test.go (new), internal/gittest/gittest.go
Covers: GF21, GF22, GF23, GF30

## What to build

Add Unavailable to the capability owner and migrate its first real consumers in bounds and gittest. Keep the existing low-level API, class vocabulary, line format, and strict policy. The later caller batches consume this accepted helper.

Supply a canned unavailable-operation detail to the new helper. Observe the canonical reason and the original capability class before Skip ends the test.

## Acceptance

- [ ] Unavailable derives the canonical prefix from the class (GF21).
- [ ] Migrated FIFO and symlink sites preserve kind and class (GF22).
- [ ] Unavailable writes its record before it skips (GF23).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/bounds`
- `bench test --package ./internal/capability`
- `bench test --package ./internal/gittest`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
